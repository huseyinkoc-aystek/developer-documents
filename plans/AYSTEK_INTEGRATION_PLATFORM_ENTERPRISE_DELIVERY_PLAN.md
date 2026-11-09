# AYSTEK CORE — Enterprise Integration Platform Delivery Plan

## 1) Amaç ve Kapsam
Bu plan, `AYSTEK CORE İçin Kurumsal Düzeyde Esnek Üçüncü Parti Entegrasyon Platformu Tasarımı.md` dokümanını AYSTEK CORE geliştirme standartlarına (Go modül iskeleti, seed zinciri, i18n 28 dil, RBAC, multi-tenant/site, audit, outbox, Next.js admin parity) uygun şekilde ürünleştirme planına dönüştürür.

Hedef:
- Tenant+site scoped, event-driven, güvenli, gözlemlenebilir bir Integration Platform kurmak.
- ERP/CRM/GPS/Finans ve benzeri üçüncü parti servislerde “kod yazmadan konfigürasyon” seviyesine yaklaşmak.
- Mevcut core modüllerle (audit/outbox/queue/ops/middleware) uyumlu çalışmak.

Kapsam dışı (ilk release):
- Tam low-code workflow builder
- Çok adımlı görsel orchestrator
- Tüm sağlayıcılar için hazır connector

---

## 2) Mimari ve Bounded Context Ayrımı

### 2.1 `internal/modules/integrations` (ana bounded context)
Sorumluluklar:
- Connector registry (auth, base_url, tls/mTLS, env config)
- Action katalogu (endpoint/method/request/response şablonları)
- Mapping (DSL v1: JMESPath + template)
- Hook/Trigger yönetimi (pre/post, sync/async)
- Policy (retry/backoff, timeout, circuit-breaker, rate-limit profile)
- Execution log + state + DLQ referansı

Sorumluluk dışı:
- Domain-specific business workflows (orders/payments/products içinde kalmalı)
- Genel queue altyapısı (platform/queue kullanılır)

### 2.2 `internal/modules/integrations_runtime` (opsiyonel alt modül)
Sorumluluklar:
- Async execution worker
- Outbox consumer entegrasyonu
- Retry/dead-letter lifecycle

### 2.3 `internal/modules/integrations_catalog` (opsiyonel alt modül)
Sorumluluklar:
- Capability templates (örn. `gps.latest_position`, `erp.invoice.push`)
- Connector preset’leri

Not: Sprint-1..3 içinde tek modül (`integrations`) altında başlayıp Sprint-4+ ile alt modül ayrımı yapılabilir.

---

## 3) Veri Modeli (Mongo) ve İndeks Stratejisi

### 3.1 Koleksiyonlar
1. `integration_connectors`
2. `integration_actions`
3. `integration_mappings`
4. `integration_hooks`
5. `integration_policies`
6. `integration_executions`
7. `integration_dlq`
8. `integration_sync_state`
9. `integration_capabilities`
10. `integration_secrets_ref` (secret pointer metadata; plaintext yok)

### 3.2 Ortak alan standardı
Tüm koleksiyonlarda:
- `tenant_id`, `site_id`
- `created_at/by`, `updated_at/by`
- `deleted_at/by` (soft delete)
- `status` (modüle göre)
- `version` (optimistic update için)

### 3.3 Kritik indeksler
- Connector unique: `(tenant_id, site_id, code, deleted_at)` partial unique
- Action unique: `(tenant_id, site_id, connector_id, action_key, deleted_at)` partial unique
- Hook query: `(tenant_id, site_id, event_key, status, deleted_at)`
- Execution query: `(tenant_id, site_id, hook_id, started_at desc)`
- Execution status: `(tenant_id, site_id, status, next_retry_at)`
- Sync state: `(tenant_id, site_id, connector_id, action_key)` unique

### 3.4 TTL/retention
- `integration_executions` detay payload TTL: 30-90 gün (policy)
- `integration_dlq` TTL: 30 gün
- Audit log retention platform standardına göre kalır

---

## 4) API Sözleşmesi (Admin/Internal)

### 4.1 Admin API (`/admin/integrations/...`)
- `GET/POST /connectors`
- `GET/PUT/DELETE /connectors/:id`
- `POST /connectors/:id/test`
- `GET/POST /connectors/:id/actions`
- `GET/PUT/DELETE /actions/:id`
- `GET/POST /hooks`
- `GET/PUT/DELETE /hooks/:id`
- `POST /hooks/:id/execute` (manual test-run)
- `GET /executions`
- `GET /executions/:id`
- `POST /executions/:id/retry`
- `GET /dlq`
- `POST /dlq/:id/requeue`
- `GET/PUT /policies/:scopeKey`

### 4.2 Internal API (`/internal/integrations/...`)
- `POST /worker/execute`
- `POST /worker/retry-due`
- `POST /worker/cleanup-retention`

### 4.3 Response & error standard
- AYSTEK response envelope kullanılacak (`success/message/error/data`)
- Örnek error code’lar:
  - `INTEGRATION_CONNECTOR_INVALID`
  - `INTEGRATION_ACTION_NOT_FOUND`
  - `INTEGRATION_MAPPING_INVALID`
  - `INTEGRATION_EXECUTION_RETRY_EXCEEDED`
  - `INTEGRATION_SSRF_BLOCKED`

---

## 5) Seed / RBAC / Menu / i18n Planı

### 5.1 Permission seed (örnek)
- `integrations.read`
- `integrations.connector.manage`
- `integrations.action.manage`
- `integrations.hook.manage`
- `integrations.execute`
- `integrations.retry`
- `integrations.policy.manage`
- `integrations.audit.read`

### 5.2 Menu seed
Parent: `operations` (veya mevcut enterprise parent)
- `/admin/integrations/connectors`
- `/admin/integrations/hooks`
- `/admin/integrations/executions`
- `/admin/integrations/policies`

### 5.3 i18n seed
- `translation_seed.go` içinde tüm yeni key’ler `allLangs(...)` ile 28 dilde
- Hardcode metin yok

### 5.4 Role seed
- Admin auto-assign garanti
- Editor için read + sınırlı execute (opsiyonel)

---

## 6) Go Modül Uygulama Standardı (Zorunlu)
Her yeni Go modülü/alt modülü için en az:
- `model/`
- `dto/`
- `repository/interface.go`
- `repository/mongo.go`
- `service/service.go`
- `service/service_test.go`
- `routes.go`
- `handler/admin.go` (gerekirse `public.go`, `internal.go`)

Ek zorunlular:
- Index registration
- `cmd/api/main.go` wiring
- Seed zinciri: menu + permission + role + translation + domain_seed

---

## 7) Next.js Admin UI Planı

### 7.1 Sayfalar
- `/admin/integrations/connectors`
- `/admin/integrations/connectors/new`
- `/admin/integrations/connectors/[id]`
- `/admin/integrations/hooks`
- `/admin/integrations/executions`
- `/admin/integrations/policies`

### 7.2 Dosya standardı
- `lib/types/integrations.ts`
- `lib/api/integrations.ts` (typed client)
- `app/(dashboard)/admin/integrations/...`

### 7.3 UI standartları
- Posts tasarım çizgisi: modern header, breadcrumb, geri butonu
- Liste: stats + search + filter + pagination
- Formlar modal değil, ayrı sayfa
- Tüm metinler i18n key

---

## 8) Güvenlik, Observability, Reliability

### 8.1 Güvenlik
- Outbound URL allowlist + redirect off
- Secret plaintext yok (secret_ref + encryption)
- OAuth token refresh + rotate
- Webhook signature verify (HMAC/timestamp/replay)
- RBAC + tenant/site scope zorunlu

### 8.2 Observability
- Prometheus metrikleri:
  - `integration_execution_total{connector,action,status}`
  - `integration_execution_duration_seconds{connector,action}`
  - `integration_retry_total{reason}`
  - `integration_dlq_total`
- OTel traces + correlation ID
- Structured logs (PII masking)

### 8.3 Reliability
- Retry policy: exponential backoff + jitter
- Circuit breaker profile
- Idempotency key desteği
- DLQ + requeue endpoint

---

## 9) Epic > Feature > User Story Backlog

## Epic-1: Foundation (Config Plane)
### Feature 1.1 Connector CRUD
- Story: Tenant/site scoped connector oluşturma/düzenleme/silme
- AC:
  - Scope zorunlu
  - Secret alanları maskeli
  - Soft delete uygulanır
- Test:
  - service validation + repository scope

### Feature 1.2 Action CRUD
- Story: Connector altında action yönetimi
- AC:
  - action_key unique
  - request/response schema validation

## Epic-2: Execution Plane
### Feature 2.1 Hook Engine
- Story: pre/post hook execute
- AC:
  - sync ve async ayrımı
  - timeout + retry policy

### Feature 2.2 Execution Logs
- Story: execution list/detail/retry
- AC:
  - filter + pagination
  - payload sanitize

## Epic-3: Governance & Policy
### Feature 3.1 Policy yönetimi
- Story: retry/backoff/circuit/rate-limit policy ekranı
- AC:
  - optimistic update (`version`)

### Feature 3.2 Capability Catalog
- Story: preset action/hook setleri
- AC:
  - template import/export

## Epic-4: Hardening
### Feature 4.1 Security controls
### Feature 4.2 Observability dashboards
### Feature 4.3 DR/Runbook/Retention

---

## 10) Sprint Planı (2 haftalık)

### Sprint-1 (Foundation)
- Integrations connector + action + mapping CRUD (backend + admin UI)
- Seed zinciri + i18n + RBAC
- Baseline testler

### Sprint-2 (Execution)
- Hook execution (sync/async), execution logs, retry
- Worker internal routes
- DLQ v1

### Sprint-3 (Policy + Security)
- Policy engine (timeout/retry/rate/circuit)
- SSRF allowlist, secret_ref flow
- Webhook signature verify

### Sprint-4 (Catalog + Productization)
- Capability templates
- Bulk import/export (JSON)
- Dashboard + metrics panel

### Sprint-5 (Hardening & GA)
- Performance tuning
- Retention jobs
- DR runbook + alert rules
- UAT + cutover checklist

---

## 11) Riskler ve Mitigasyon
- Risk: Connector çeşitliliği nedeniyle mapping karmaşıklığı
  - Mitigasyon: DSL v1’i sınırlı tut, capability template kullan
- Risk: Güvenlik açıkları (SSRF/secret leak)
  - Mitigasyon: allowlist + secret_ref + audit + secure review gate
- Risk: Runtime yükte kuyruk birikmesi
  - Mitigasyon: worker autoscale + retry budget + DLQ

---

## 12) Definition of Done (DoD)
Bir feature “done” sayılması için:
1. Go modül yapısı eksiksiz
2. Seed zinciri eksiksiz (menu/permission/role/translation/domain)
3. Next typed api + types + page teslim
4. Tenant/site scope doğrulanmış
5. Unit/integration testler çalışır
6. Observability metrikleri ekli
7. Dokümantasyon güncel

---

## 13) Uygulama Sırası (Başlangıç Önerisi)
1. Sprint-1 ile `internal/modules/integrations` çekirdeğini çıkar.
2. Seed ve admin menü/i18n’i aynı sprintte kapat.
3. Sprint-2’de execution + retry/dlq ile gerçek iş akışını aç.
4. Sprint-3’te security hardening’i bitir (SSRF, secrets, webhook verify).
5. Sprint-4/5’te catalog + observability + runbook ile production readiness tamamla.

