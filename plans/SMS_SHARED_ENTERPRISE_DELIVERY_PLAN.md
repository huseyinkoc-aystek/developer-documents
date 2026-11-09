# SMS Shared Enterprise Delivery Plan (AYSTEK-CORE)

- Tarih: 2026-02-24
- Durum: Draft / Uygulama Öncesi
- Sahip: Platform (Principal Engineer)
- Kapsam: Shared `sms` modülü (multi-provider, tenant/site scoped, reusable)

## Takip Panosu
| Alan | Durum | Not |
|---|---|---|
| A) Mimari Tasarım | Hazır | Onay bekliyor |
| B) Veri Modeli (Mongo) | Hazır | Index/TTL kararları net |
| C) API Sözleşmesi | Hazır | DTO/envelope uyumlu |
| D) Permission/Menu/i18n | Hazır | 28 dil stratejisi tanımlı |
| E) Backlog/Sprint Planı | Hazır | Epic/Feature/Story çıktı |
| F) Güvenlik/Observability | Hazır | HMAC, masking, metrics dahil |
| G) Compliance Checklist | Hazır | PR checklist dahil |

### Sprint/Epic Takip Matrisi
| Takip ID | Sprint | Epic | Durum | Çıktı |
|---|---|---|---|---|
| SMS-E1 | Sprint-1 | Foundation (interface + registry + mock) | Planlandı | Provider abstraction çekirdeği |
| SMS-E2 | Sprint-2 | Config/Policy/Send | Planlandı | Secret/config/policy + idempotent send |
| SMS-E3 | Sprint-3 | Async + DLR | Planlandı | Outbox retry + webhook signature/replay |
| SMS-E4 | Sprint-4 | Admin UI + Hardening | Planlandı | Providers/policies/templates/logs/outbox ekranları |

### Jira Issue Bağlantı Matrisi
| Jira ID | Tip | Başlık | Sprint | Epic | Durum | Bağımlılık |
|---|---|---|---|---|---|---|
| SMS-101 | Epic | SMS Foundation & Provider Abstraction | Sprint-1 | SMS-E1 | To Do | - |
| SMS-102 | Story | Provider interface + registry + mock provider | Sprint-1 | SMS-E1 | To Do | SMS-101 |
| SMS-103 | Story | Provider FieldSpec sözleşmesi + dynamic config form contract | Sprint-1 | SMS-E1 | To Do | SMS-102 |
| SMS-104 | Story | Secret storage/masking + config rotate flow | Sprint-1 | SMS-E1 | To Do | SMS-102 |
| SMS-201 | Epic | Policy Routing & Send Pipeline | Sprint-2 | SMS-E2 | To Do | SMS-101 |
| SMS-202 | Story | Tenant/site default provider + message type policy engine | Sprint-2 | SMS-E2 | To Do | SMS-201 |
| SMS-203 | Story | Idempotent send (correlation/idempotency key) | Sprint-2 | SMS-E2 | To Do | SMS-202 |
| SMS-204 | Story | Rate limit/quota enforcement + counters | Sprint-2 | SMS-E2 | To Do | SMS-202 |
| SMS-301 | Epic | Async Outbox + DLR Webhooks | Sprint-3 | SMS-E3 | To Do | SMS-201 |
| SMS-302 | Story | Outbox worker retry/backoff/dead-letter | Sprint-3 | SMS-E3 | To Do | SMS-301 |
| SMS-303 | Story | DLR webhook signature + replay guard + status mapping | Sprint-3 | SMS-E3 | To Do | SMS-301 |
| SMS-304 | Story | Provider failover chain execution | Sprint-3 | SMS-E3 | To Do | SMS-203 |
| SMS-401 | Epic | Admin UI + Hardening & Release | Sprint-4 | SMS-E4 | To Do | SMS-301 |
| SMS-402 | Story | Admin providers/config/policies/templates screens | Sprint-4 | SMS-E4 | To Do | SMS-401 |
| SMS-403 | Story | Logs/outbox screens + retry action + RBAC gating | Sprint-4 | SMS-E4 | To Do | SMS-302 |
| SMS-404 | Story | Security/observability/e2e smoke + release checklist | Sprint-4 | SMS-E4 | To Do | SMS-403 |
| SMS-405 | Story | Provider Credentials & Secrets Center (UI+API+rotate/test) | Sprint-4 | SMS-E4 | To Do | SMS-104 |

### Jira Workflow Kuralları
- Issue durumları: `To Do -> In Progress -> Code Review -> QA -> Done`
- Her Story için zorunlu alt görevler:
  - `BE`: model/dto/repo/service/service_test/routes
  - `SEED`: translation/menu/permission/role/domain
  - `FE`: lib/types + lib/api + admin pages
  - `SEC`: secret/hmac/replay/rate-limit kontrolleri
  - `OBS`: audit + metrics + tracing
- Definition of Done kapanış etiketi: `DoD-SMS-PASS`

### Uygulama İlerlemesi (Checklist)
- [ ] Sprint-1 tamamlandı (backend foundation + seed zinciri)
- [ ] Sprint-2 tamamlandı (policy + routing + quota/rate limit)
- [ ] Sprint-3 tamamlandı (async worker + DLR callback + retry)
- [ ] Sprint-4 tamamlandı (admin UI + RBAC + e2e smoke)
- [ ] Security review tamamlandı (secrets, signature, replay)
- [ ] Observability review tamamlandı (audit, metrics, tracing)
- [ ] Release readiness tamamlandı (docs, runbook, rollback)

---

## A) Mimari Tasarım
### 1) Bounded Context
- Sorumluluk: provider abstraction, routing/policy, send pipeline (sync+async), DLR webhook, logs/outbox.
- Sorumluluk dışı: modül-spesifik business rule (stock/pms/contracts akış kuralları), UI policy authority.
- Kritik alt-sınır: provider entegrasyonuna özel kimlik bilgileri `SMS Credentials Center` içinde yönetilir; diğer modüller secret görmez, yalnızca `sms service` kullanır.
- Provider source-of-truth kararı:
  - Runtime provider listesi koddan (`providers.Register(...)`) gelir.
  - DB tarafı provider çalıştırma kaynağı değil; tenant/site bazlı config/policy/state taşır.
  - `sms_providers` varsa yalnızca UI katalog/cache amaçlıdır, registry ile tutarlı olmak zorundadır.

### 2) Katman diyagramı
- `internal/modules/sms` -> `handler` -> `service` -> `repository` -> `providers/*`
- Cross-cutting: `audit`, `metrics`, `tracing`, `rate-limit`, `secret-store`, `outbox worker`.
- Secret erişim zinciri: `service -> secret-store adapter -> KMS/encrypted store` (UI/API yalnız masked metadata döndürür).
- DLR güvenlik/adaptasyon zinciri:
  - `handler` -> `providers/<id>/signature.go` (verify)
  - `providers/<id>/dlr_parser.go` (provider payload -> canonical DLR DTO)
  - `service` (idempotent status update + audit/outbox)

### 3) Provider seçimi (policy)
1. Explicit provider override (sadece yetki ile)
2. Message-type policy (`otp|transactional|marketing`)
3. Tenant+Site default provider
4. Tenant global default
5. System fallback (`mock`)

- Failover chain destekli: `primary -> secondary -> mock`.

### 4) Idempotency anahtar standardı
- Caller idempotency/correlation verirse öncelikle o kullanılır.
- Caller vermezse sms service standard key üretir:
  - `tenant_id + site_id + to_e164 + message_type + content_hash + window_bucket`
- OTP için `window_bucket` zorunlu (örn. 60 sn / 5 dk policy’den okunur).
- Aynı key için davranış: ilk başarılı/işlemdeki kayıt dönülür, duplicate send tetiklenmez.

---

## B) Veri Modeli (Mongo)
### Koleksiyonlar
- `sms_providers` (opsiyonel katalog/cache, registry mirror)
- `sms_provider_configs`
- `sms_provider_credentials` (versioned secret metadata)
- `sms_policies`
- `sms_templates`
- `sms_messages`
- `sms_outbox_jobs`
- `sms_webhook_receipts`
- `sms_quota_counters` (opsiyonel read model)
- `sms_message_payloads` (opsiyonel, kısa TTL; ham metin gerekiyorsa)

### Kritik alanlar
- Tüm koleksiyonlarda: `tenant_id`, `site_id`, audit + soft-delete alanları.
- `sms_messages`: `idempotency_key`, `correlation_id`, `provider_message_id`, `status`, `to_masked`.
- `sms_provider_configs`: `config_public`, `config_secret_ref|encrypted`, `is_default`.
- `sms_provider_credentials`: `provider_id`, `version`, `secret_ref`, `masked_snapshot`, `rotated_at`, `rotated_by`, `is_active`.
- `sms_policies`: `message_type`, `primary_provider_id`, `failover_chain`, `rate/quota`.
- `sms_messages`: mümkünse ham içerik tutmaz; `template_ref`, `content_hash`, `segment_count`, `encoding`.
- `sms_message_payloads`: `message_id`, `payload_encrypted`, kısa süreli TTL.

### İndeksler
- `sms_messages`: unique partial `(tenant_id, site_id, idempotency_key, deleted_at:null)`
- `sms_provider_configs`: unique partial `(tenant_id, site_id, provider_id, deleted_at:null)`
- `sms_provider_credentials`: unique partial `(tenant_id, site_id, provider_id, version, deleted_at:null)`
- `sms_policies`: unique partial `(tenant_id, site_id, message_type, deleted_at:null)`
- `sms_outbox_jobs`: `(tenant_id, site_id, status, next_run_at)`
- `sms_webhook_receipts`: unique `(tenant_id, site_id, provider_id, webhook_event_id)`
- `sms_message_payloads`: `(tenant_id, site_id, created_at)` TTL index

### TTL/Retention
- `sms_webhook_receipts`: 30 gün
- `sms_outbox_jobs`: done 30 gün, dead 90 gün
- OTP detay retention: 7-14 gün (PII minimal)
- `sms_messages` (audit/read model): 180 gün (PII minimize edilmiş alanlarla)
- `sms_message_payloads` (opsiyonel ham içerik): 14 gün TTL

### Secret yönetimi kararı (platform standardı)
- `SecretStore` interface:
  - `Put(ctx, scope, payload) -> secret_ref`
  - `Get(ctx, secret_ref) -> payload`
  - `Rotate(ctx, secret_ref, payload) -> new_secret_ref`
  - `Delete(ctx, secret_ref)`
- Ortam stratejisi:
  - Prod: `secret_ref` (Vault/KMS/Secrets Manager adapter)
  - Dev/CI fallback: `AES-GCM encrypted payload` (env key ile)
- UI/Logs her koşulda masked veri gösterir (son 4 hariç gizli).

---

## C) API Sözleşmesi
### Admin
- Providers: `GET/POST/PUT /admin/sms/providers`, `.../:id/enable`, `.../:id/disable`, `.../:id/field-spec`, `.../:id/test-send`
- Configs: `GET/POST/PUT /admin/sms/provider-configs`, `.../:id/rotate-secret`
- Credentials Center:
  - `GET /admin/sms/providers/:id/credentials` (masked list + active version)
  - `POST /admin/sms/providers/:id/credentials` (create encrypted credential version)
  - `POST /admin/sms/providers/:id/credentials/:version/activate`
  - `POST /admin/sms/providers/:id/credentials/:version/test`
  - `POST /admin/sms/providers/:id/credentials/:version/revoke`
- Policies: `GET/POST/PUT /admin/sms/policies`, `.../:id/validate`
- Templates: `GET/POST/PUT/DELETE /admin/sms/templates`
- Send/Logs/Outbox:
  - `POST /admin/sms/messages/send`
  - `GET /admin/sms/messages`
  - `POST /admin/sms/outbox/jobs` (bulk enqueue)
  - `GET /admin/sms/outbox/jobs`
  - `POST /admin/sms/outbox/jobs/:id/retry`
- Policy optimistic concurrency:
  - `PUT /admin/sms/policies/:id` çağrısında `version` (veya `If-Match`) zorunlu
  - Version uyuşmazlığında `409 SMS_POLICY_VERSION_CONFLICT`

### Webhook
- `POST /public/sms/webhooks/:provider_id/dlr`
- Signature: `X-Aystek-Signature`
- Replay guard: `X-Webhook-Event-Id`
- Provider adapter zorunlu:
  - `signature.go`: header/timestamp/signature verify
  - `dlr_parser.go`: payload canonicalizasyonu + status mapping

### Envelope / Pagination / Hata
- Standard response envelope + pagination (`items/page/limit/total/total_pages`).
- Hata kodları: `SMS_PROVIDER_NOT_FOUND`, `SMS_IDEMPOTENCY_CONFLICT`, `SMS_WEBHOOK_SIGNATURE_INVALID`, `SMS_RATE_LIMITED`, vb.
- Credentials hata kodları:
  - `SMS_PROVIDER_CREDENTIAL_INVALID`
  - `SMS_PROVIDER_CREDENTIAL_ROTATE_REQUIRED`
  - `SMS_PROVIDER_CREDENTIAL_REVOKED`
  - `SMS_PROVIDER_CREDENTIAL_VERSION_NOT_FOUND`
  - `SMS_POLICY_VERSION_CONFLICT`
  - `SMS_PHONE_INVALID_E164`

---

## D) Permission + Menu + i18n
### Permissions
- `sms:read`, `sms:send`, `sms:send_bulk`
- `sms:provider_manage`, `sms:provider_override`
- `sms:config_manage`, `sms:policy_manage`, `sms:template_manage`
- `sms:credentials_read`, `sms:credentials_manage`, `sms:credentials_rotate`, `sms:credentials_test`
- `sms:logs_read`, `sms:outbox_manage`, `sms:webhook_read`, `sms:quota_manage`

### Menu
- `operations -> messaging -> sms`
- Alt sayfalar: dashboard/providers/credentials/policies/templates/messages/outbox

### i18n
- Namespace: `sms.*`, `menu.operations.messaging*`
- 28 dil stratejisi: `allLangs(tr,en)` fallback yaklaşımı

---

## E) Backlog / İş Planı
### Epic-1: Foundation
- Feature: module skeleton + provider interface + registry + mock provider
- AC: tenant/site scope, unit test, seed zinciri

### Epic-2: Config/Policy/Send
- Feature: config secret mgmt, policy engine, idempotent send, rate/quota
- AC: override yetki kontrolü, deterministic policy selection
- Feature: phone normalization + segment accounting
- AC: E.164 normalize, GSM-7/UCS-2 encoding ve segment hesabı service katmanında deterministik

### Epic-3: Async + DLR
- Feature: outbox worker, retry/backoff/dead-letter, webhook doğrulama+mapping
- AC: replay guard ve status mapping testleri
- Feature: provider parser adapter standardı
- AC: her provider için `signature.go + dlr_parser.go` olmadan provider enabled edilemez

### Epic-4: Admin UI
- Feature: provider yönetimi, FieldSpec tabanlı form, policy/template/log/outbox ekranları
- AC: hardcoded string yok, RBAC UI gating
- Feature: Credentials Center
- AC: API key/secret alanları güvenli girilir, maskeli görünür, rotate+test akışı çalışır, raw secret geri okunamaz.

### Sprint Planı (2 hafta)
- Sprint-1: Epic-1 + Secret/Config çekirdeği
- Sprint-2: Epic-2 + send pipeline
- Sprint-3: Epic-3 + webhook
- Sprint-4: Epic-4 + credentials center + hardening/load tests

### Riskler / Mitigasyon
- Provider API varyasyonu -> adapter sözleşmesi + canonical DLR DTO
- Secret güvenliği -> vault/ref + masking + audit
- Abuse/spam -> Redis rate-limit adapter + per-tenant/site quota, Mongo sadece read model

---

## F) Güvenlik & Observability
### Güvenlik
- Secret plaintext tutulmaz, masked render.
- Secret giriş UX:
  - `write-only` input (kayıttan sonra geri gösterim yok)
  - rotate sırasında zorunlu onay
  - test-send yalnız aktif credential versiyonu ile
- Webhook HMAC + timestamp skew + replay guard zorunlu.
- Phone/PII logları maskeli.
- Endpoint abuse için tenant/site rate limit.
- CSV/Export güvenliği: merkezi export helper ile OWASP uyumlu formula-injection nötralizasyonu (`=,+,-,@` guard).

### Observability
- Audit eventler: provider/config/policy/send/dlr/outbox lifecycle.
- Metrikler:
  - `sms_send_total{provider,status,type}`
  - `sms_send_latency_ms_bucket{provider}`
  - `sms_retry_total{provider}`
  - `sms_webhook_invalid_signature_total`
  - `sms_segment_total{provider,encoding}`
  - `sms_rate_limit_block_total{tenant,site}`

---

## G) AYSTEK-CORE Standards Compliance Checklist
| Kalem | Durum |
|---|---|
| Go core dosya seti (model/dto/repo/service/test/routes) | ☐ |
| Seeds (menu/permission/role/translation/domain) | ☐ |
| 28 dil key seti | ☐ |
| Next typed api + types + admin pages | ☐ |
| Envelope/pagination/hata kodları uyumu | ☐ |
| Tenant+Site isolation testleri | ☐ |
| Security (secret/hmac/replay/rate limit) | ☐ |
| Audit + Metrics + Tracing | ☐ |
| Credentials Center (masked read, rotate, revoke, test) | ☐ |
| Credentials RBAC (read/manage/rotate/test) | ☐ |
| Credentials i18n key seti (`sms.credentials.*`) | ☐ |
| Provider parser adapter (`signature.go` + `dlr_parser.go`) | ☐ |
| E.164 normalization + segment hesaplama testleri | ☐ |
| Policy versioning (optimistic lock) | ☐ |

## PR Checklist
- [ ] Provider abstraction provider bağımlılığını modüllerden ayırıyor.
- [ ] Override sadece yetkiyle çalışıyor.
- [ ] DLR webhook replay engeli aktif.
- [ ] Idempotency conflict davranışı testli.
- [ ] UI metinleri tamamen i18n key.
- [ ] go test + lint/typecheck yeşil.
- [ ] Credentials Center raw secret sızdırmıyor (API/UI/log).
- [ ] Credential rotate/activate/revoke akışları audit event üretiyor.
