---
Title: Distributed Zero-Trust Tenant Architecture Specification (TR/EN)
Version: 2.1.0
Status: Enterprise Certified
Last Updated: 2026-02-12
Related Documents:
  - ../01-architecture/Enterprise_Isolation_Whitepaper.md
  - ../02-security/TENANT_SITE_SECURITY_MATRIX.md
  - ../archive/internal/ENTERPRISE_CONTEXT_GUARANTEE.md
  - ../archive/internal/MIDDLEWARE_CHAIN_INVARIANTS.md
  - ../04-api/ERROR_CODE_REFERENCE.md

---
# Table of Contents
1. [Türkçe](#türkçe)
2. [English](#english)

## Türkçe

### 1. Kapsam ve Tasarım Hedefi
Bu doküman, AYSTEK Core'u monolitten dağıtık servislere evriltirken tenant/site izolasyonunu zayıflatmadan koruyan Zero-Trust mimari sözleşmesini tanımlar.

Temel hedefler:
1. Tenant bağlamı olmayan hiçbir işlem veri düzlemine ulaşamaz.
2. Site bağlamı gereken yüzeylerde eksik site ile işlem yapılamaz.
3. İç servis çağrıları "güvenli ağ" varsayımı ile değil, imzalı ve doğrulanmış zarf modeli ile kabul edilir.
4. Ölçek artışı (100k+ tenant) izolasyon kurallarını gevşetmez.

### 2. Zero-Trust İç Servis Modeli
#### 2.1 İmzalı iç çağrı sözleşmesi
Her servisler arası çağrı aşağıdaki zarfı taşır:
- `tenant_id`
- `site_id` (gerekli yüzeylerde zorunlu)
- `actor_id` / `service_id`
- `trace_id`
- `timestamp`
- `signature`

`signature`, canonical payload + paylaşılan/rotasyonlu anahtar ile doğrulanır. İmza doğrulanmazsa çağrı reddedilir.

#### 2.2 Tenant-scope service mesh ilkesi
- mTLS zorunlu.
- Policy: servis A, servis B'ye yalnızca belirli tenant-scope methodları için erişebilir.
- Mesh policy ve uygulama policy birlikte çalışır; biri açık verse bile diğeri ihlali bloklar.

#### 2.3 Envelope doğrulama sırası
1. İmza doğrulama
2. Tenant durumu (active/inactive/tombstone)
3. Site'nin tenant'a aitliği
4. İzin/modül kontrolü
5. İşlem

### 3. Dağıtık Event Bus Modeli (Kafka/RabbitMQ Hazır)
#### 3.1 Zorunlu event envelope
Event payload üstünde immutable metadata:
- `event_id`
- `event_type`
- `tenant_id`
- `site_id`
- `trace_id`
- `occurred_at`
- `producer_signature`

#### 3.2 Immutable tenant binding
Consumer, event içindeki tenant/site değerini runtime context'e "kopyalamaz"; önce doğrular, sonra immutable request context'e kilitler.

#### 3.3 Replay suppression
- `event_id` + `tenant_id` bazlı dedup store
- TTL + state-aware replay guard
- Tombstone tenant için replay kesin reddedilir

ASCII akış:
```text
Producer -> Sign Envelope -> Bus Topic/Queue -> Consumer Verify
                                             -> Tenant State Check
                                             -> Site/Tenant Match
                                             -> Idempotency Check
                                             -> Handler Execute
```

### 4. Multi-Region Dağıtım Modeli
#### 4.1 Tenant partition stratejisi
- Tenant, `region_home` ile atanır.
- Yazma trafiği tenant home region'a sabitlenir.
- Read replica erişimi policy kontrollüdür.

#### 4.2 Cache sharding
Anahtar deseni:
`<entity>:<tenant_id>:<site_id>:<resource_id>:<version>`

Shard seçimi:
- birincil: `hash(tenant_id)`
- ikincil: `hash(site_id)` (site-scope veri için)

#### 4.3 Event partitioning
- Kafka partition key: `tenant_id`
- Site yoğunluğu yüksek tenantlarda: `tenant_id:site_id`

#### 4.4 Read-model partitioning
- Search/read model indexleri tenant namespace ile ayrılır.
- Cross-tenant global index yasaktır.

ASCII:
```text
            +---------------- Control Plane ----------------+
            | Tenant State Authority | Key Mgmt | Policy DB |
            +--------------------+--------------------------+
                                 |
      +--------------------------+--------------------------+
      |                                                     |
 +----v------------------+                         +--------v---------------+
 | Region EU             |                         | Region US              |
 | API + Services + DB   |                         | API + Services + DB    |
 | Cache shard (tenant)  |                         | Cache shard (tenant)   |
 | Event partitions      |                         | Event partitions       |
 +-----------------------+                         +------------------------+
```

### 5. Tenant State Authority Service
Merkezi "Tenant State Authority", tenant yaşam döngüsünün tek doğruluk kaynağıdır.

Durumlar:
- `active`: tüm işlemler policy'ye göre izinli
- `inactive`: yeni yazmalar engelli, kontrollü okuma
- `tombstoned`: event/job/cache erişimi dahil tüm tenant işlemleri reddedilir

Kurallar:
1. Her request/event başlangıcında state check
2. State değişiminde cache invalidation ve job kill-signal
3. Tombstone sonrası replay ve retry kalıcı blok

### 6. Servisler Arası Context Propagation
#### 6.1 Trace sürekliliği
- `trace_id` her hop'ta taşınır
- Geçersiz trace formatında yeni trace üretilir, eski değer güvenlik loguna yazılır

#### 6.2 Tenant integrity doğrulaması
- Downstream servis tenant/site'ı header'dan değil signed envelope/request context'ten okur
- Header kaynaklı tenant override yasaktır

### 7. Failure Isolation Modeli
#### 7.1 Blast radius containment
- Tenant bazlı rate limit
- Tenant bazlı queue isolation
- Tenant bazlı circuit breaker

#### 7.2 Circuit breaker tenant izolasyonu
Bir tenant veya site hata üretiyorsa breaker yalnızca ilgili tenant/site dilimini etkiler, global trafiği düşürmez.

ASCII:
```text
[Tenant A surge] -> [Tenant A breaker opens]
                 -> [Tenant B traffic continues]
```

### 8. 100k+ Tenant Ölçek Modeli
Önerilen ölçekleme eksenleri:
1. Stateless API horizontal scaling
2. Tenant-aware cache partition
3. Tenant-partitioned event throughput
4. Tenant-home-region routing
5. Metadata/control plane için güçlü tutarlılık, data plane için hedefli eventual consistency

Beklenen davranış:
- Tenant başına p99 gecikme izlenebilir ve izole kalır.
- Gürültülü komşu etkisi tenant-scoped limiter + breaker ile sınırlandırılır.

### 9. Monolitten Dağıtıka Geçiş Yol Haritası
1. Monolit içinde envelope ve context sözleşmesini sabitle
2. Event publish/consume hattında immutable tenant binding'i aktif et
3. Cache key standardizasyonunu tek yardımcıya taşı
4. Domain bazlı servis ayrıştırmasını "önce read, sonra write" sırasıyla yap
5. Region-aware routing'i tenant bazında kademeli aç

Geçiş kuralı:
İzolasyon garantisi yeni serviste, monolitteki seviyeden düşük olamaz.

### 10. Asla Bozulmaması Gereken İzolasyon İnvariantları
1. Tenant olmadan tenant-scope işlem yok.
2. Site gereken yüzeyde site olmadan işlem yok.
3. Site-tenant eşleşmesi bozuksa işlem yok.
4. Tombstoned tenant için event/job/cache erişimi yok.
5. Repository sorguları tenant dominance olmadan çalışamaz.
6. Canonical hata yapısı tek tiptir: `{error, message, trace_id}`.

### 11. Operasyonel Doğrulama Stratejisi
1. Build-time: AST tabanlı repository dominance analizi
2. Runtime: context integrity validator + security event emission
3. CI: grep gate + audit gate + integration matrix
4. Production: multi-region smoke + replay + failover testleri

---

## English

### 1. Scope and Design Objective
This specification defines how AYSTEK Core evolves from monolith to distributed services while preserving strict tenant/site isolation under a Zero-Trust model.

Primary objectives:
1. No data-plane operation executes without tenant context.
2. Site-required surfaces cannot execute without site context.
3. Internal calls are accepted only through signed and verified envelopes.
4. Scale growth (100k+ tenants) never relaxes isolation rules.

### 2. Zero-Trust Internal Service Model
#### 2.1 Signed internal call contract
Every service-to-service request carries:
- `tenant_id`
- `site_id` (mandatory on site-scoped surfaces)
- `actor_id` / `service_id`
- `trace_id`
- `timestamp`
- `signature`

`signature` is validated against canonical payload + rotated key material. Invalid signatures are rejected.

#### 2.2 Tenant-scoped service mesh principle
- mTLS is mandatory.
- Policy: service A may call service B only for approved tenant-scoped methods.
- Mesh policy and application policy are both enforced (defense-in-depth).

#### 2.3 Envelope validation order
1. Signature validation
2. Tenant lifecycle state check
3. Site-to-tenant ownership check
4. Permission/module checks
5. Business execution

### 3. Distributed Event Bus Model (Kafka/RabbitMQ Ready)
#### 3.1 Mandatory event envelope
Immutable event metadata:
- `event_id`
- `event_type`
- `tenant_id`
- `site_id`
- `trace_id`
- `occurred_at`
- `producer_signature`

#### 3.2 Immutable tenant binding
Consumers do not trust incoming runtime state; they verify event tenant/site first, then lock validated context as immutable.

#### 3.3 Replay suppression
- Dedup store keyed by `event_id + tenant_id`
- TTL and state-aware replay guard
- Tombstoned tenant replay is permanently blocked

ASCII flow:
```text
Producer -> Sign Envelope -> Bus Topic/Queue -> Consumer Verify
                                             -> Tenant State Check
                                             -> Site/Tenant Match
                                             -> Idempotency Check
                                             -> Handler Execute
```

### 4. Multi-Region Deployment Model
#### 4.1 Tenant partition strategy
- Each tenant is assigned `region_home`.
- Writes are pinned to tenant home region.
- Read replicas are policy-controlled.

#### 4.2 Cache sharding
Key format:
`<entity>:<tenant_id>:<site_id>:<resource_id>:<version>`

Shard strategy:
- primary: `hash(tenant_id)`
- secondary: `hash(site_id)` for site-scoped data

#### 4.3 Event partitioning
- Kafka partition key: `tenant_id`
- High-volume tenants may use: `tenant_id:site_id`

#### 4.4 Read-model partitioning
- Search/read-model indexes are tenant namespaced.
- Cross-tenant global indexes are forbidden.

ASCII:
```text
            +---------------- Control Plane ----------------+
            | Tenant State Authority | Key Mgmt | Policy DB |
            +--------------------+--------------------------+
                                 |
      +--------------------------+--------------------------+
      |                                                     |
 +----v------------------+                         +--------v---------------+
 | Region EU             |                         | Region US              |
 | API + Services + DB   |                         | API + Services + DB    |
 | Cache shard (tenant)  |                         | Cache shard (tenant)   |
 | Event partitions      |                         | Event partitions       |
 +-----------------------+                         +------------------------+
```

### 5. Tenant State Authority Service
Tenant lifecycle is enforced by a central Tenant State Authority.

States:
- `active`: all operations allowed per policy
- `inactive`: no new writes, controlled reads
- `tombstoned`: all event/job/cache/data access blocked

Rules:
1. State check at every request/event entry
2. Cache invalidation and job kill-signal on state transition
3. Permanent replay/retry block after tombstone

### 6. Cross-Service Context Propagation
#### 6.1 Trace continuity
- `trace_id` is propagated on every hop
- Invalid trace format triggers server-generated trace and security logging

#### 6.2 Tenant integrity validation
- Downstream services read tenant/site from verified context only
- Direct header-driven tenant overrides are forbidden

### 7. Failure Isolation Model
#### 7.1 Blast-radius containment
- Tenant-scoped rate limits
- Tenant-scoped queue isolation
- Tenant-scoped circuit breakers

#### 7.2 Circuit-breaker isolation
Failures from one tenant/site must not degrade unrelated tenants.

ASCII:
```text
[Tenant A surge] -> [Tenant A breaker opens]
                 -> [Tenant B traffic continues]
```

### 8. 100k+ Tenant Scale Model
Recommended scale axes:
1. Stateless API horizontal scale
2. Tenant-aware cache partitioning
3. Tenant-partitioned event throughput
4. Tenant-home-region routing
5. Strong consistency for control plane, targeted eventual consistency for data plane

Expected outcome:
- Per-tenant p99 latency remains observable and isolated.
- Noisy-neighbor impact is bounded by scoped limiters and breakers.

### 9. Monolith-to-Distributed Migration Path
1. Freeze envelope and context contract inside monolith first
2. Enable immutable tenant binding on publish/consume paths
3. Consolidate cache key construction into one shared helper
4. Split domains service-by-service (read-first, write-second)
5. Enable region-aware routing tenant-by-tenant

Migration rule:
No new service may provide weaker isolation than the monolith baseline.

### 10. Non-Breakable Isolation Invariants
1. No tenant-scoped operation without tenant.
2. No site-required operation without site.
3. No operation on site/tenant mismatch.
4. No event/job/cache access for tombstoned tenants.
5. No repository query without tenant dominance.
6. Canonical error schema is fixed: `{error, message, trace_id}`.

### 11. Operational Verification Strategy
1. Build-time: AST-based repository dominance analysis
2. Runtime: context integrity validator + security event emission
3. CI: grep gates + audit gates + integration matrix
4. Production: multi-region smoke, replay, and failover validation
