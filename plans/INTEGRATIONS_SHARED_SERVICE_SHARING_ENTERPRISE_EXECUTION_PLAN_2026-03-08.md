# Integrations Shared Third-Party Service Sharing Enterprise Execution Plan

Date: 2026-03-08
Owner: AI Assistant (Codex)
Scope: aystek-core / aystek-core-next / user-guide
Source:
- `/workspace/aystek-core-dev-docs/plans/AYSTEK CORE Integrations Modülüne Güvenli, İzlenebilir ve Sınırlı Üçüncü Taraf Servis Paylaşımı Ekle..md`

## Objective
Integrations modülüne “Shared Third-Party Service” katmanını ekleyerek, merkezi yönetilen connector/action setlerini tenant-site izolasyonu içinde grant/scope/limit ile güvenli biçimde kullandırmak.

## Guardrails
1. Mevcut hook/action/chain trigger semantiği değişmeyecek.
2. Tenant/site izolasyonu tüm yeni koleksiyonlarda zorunlu.
3. SSRF/allowlist, auth ve audit ilkeleri korunacak.
4. Mevcut execution + retry + DLQ semantiği bozulmayacak.
5. DDD katmanları + seed/menu/translation standartları korunacak.

## Delivery Phases

### Phase-1: Domain + Admin CRUD Foundation
1. Shared service + grant domain modeli (status/protocol/limits/scopes).
2. DTO + repository + service contract genişletmeleri.
3. Mongo implementasyonları + unique/index tasarımı.
4. Admin API endpointleri (list/create/get/update/delete).
5. Menu + translation seed foundation.

### Phase-2: Runtime Policy Enforcement
1. Shared grant çözümleyici (request-time grant lookup).
2. Action execute hattında allowed action + scope + limits enforcement.
3. Usage metering (redis + mongo rollup) ve rate/quota kontrolü.
4. Violation taxonomy + DLQ/audit bağlama.

### Phase-3: Security Hardening
1. Secret isolation (shared secret ref + masked retrieval).
2. Method/path/host guard + payload field guard.
3. Optional mTLS policy ve protocol-specific güvenlik kontrolleri.
4. Compliance metadata (PII class, retention, transfer notes).

### Phase-4: Next Admin UX
1. Shared service list/form sayfaları.
2. Grant management list/form (scope + limit builders).
3. Readable usage/limit dashboard ve warning paneli.
4. i18n standardına tam uyum.

### Phase-5: Ops, QA, Runbook
1. Smoke checklist + failure drills.
2. CI job template (post-deploy verify).
3. User-guide + operations guide güncellemesi.
4. Rollout + rollback runbook.

## Immediate Sprint Scope (Started)
1. Phase-1 backend foundation kodlaması.
2. Index + route + handler entegrasyonu.
3. Seed (menu/translation) başlangıç anahtarlarının eklenmesi.
4. Scoped test/compile doğrulaması.

## Progress Status (Updated: 2026-03-08)

### Completed
1. İş planı oluşturuldu ve roadmap netleştirildi.
2. Phase-1 backend foundation teslim edildi:
   - `model/shared_service.go` ile `SharedService`, `SharedServiceGrant`, status/protocol/limits domain modelleri eklendi.
   - `dto/shared_service.go` ile CRUD/list query DTO sözleşmeleri eklendi.
   - `repository/interface.go` shared service/grant kontratları ve filtre tipleri ile genişletildi.
   - `repository/mongo_shared_service.go` ile Mongo CRUD/list/exists/soft-delete implementasyonları eklendi.
   - `service/shared_service_service.go` ile validasyon + normalize + CRUD servis iş akışları eklendi.
   - `handler/admin_shared_service.go` + `routes.go` ile admin API endpointleri açıldı:
     - `GET/POST /admin/integrations/shared-services`
     - `GET/PUT/DELETE /admin/integrations/shared-services/:id`
     - `GET/POST /admin/integrations/shared-service-grants`
     - `GET/PUT/DELETE /admin/integrations/shared-service-grants/:id`
   - `indexes.go` içine shared service/grant indexleri eklendi.
   - `menu_seed.go` ve `translation_seed.go` içine menü + i18n seed foundation anahtarları eklendi.
3. Scoped derleme doğrulaması başarılı:
   - `go test ./internal/modules/integrations/... -run TestDoesNotExist`
   - `go test ./internal/seeds -run TestDoesNotExist`
4. Phase-2 başlangıç enforcement entegrasyonu tamamlandı:
   - Execution hattına shared grant policy resolver bağlandı (`service/shared_service_enforcement.go`).
   - Grant referansı çözümleme kaynakları:
     - `hook.config.shared_service_grant_id`
     - `overrideConfig.shared_service_grant_id`
     - `request_payload._shared_service_grant_id`
   - Zorlanan kurallar:
     - grant aktif olmalı
     - shared service aktif olmalı
     - connector/action service ve grant allowlist’i içinde olmalı
     - RPM + daily quota limit enforcement (Redis/ratelimit üzerinden)
   - Policy ihlallerinde execution `failed` işaretlenir, error code `shared_service_policy_denied`, DLQ açık ise DLQ’ye taşınır.
   - Direct test-run (`DirectExecuteAction`) için de aynı policy enforcement eklendi.
   - Yeni validation key’leri eklendi:
     - `integrations.validation.shared_service_grant_ref_invalid`
     - `integrations.validation.shared_service_policy_denied`
     - `integrations.validation.shared_service_quota_exceeded`
   - Service testleri eklendi:
     - `TestDirectExecuteActionPolicyDeniedBySharedGrant`
     - `TestProcessDueExecutionsPolicyDeniedByHookGrant`
5. Phase-2 enforcement genişletmesi tamamlandı:
   - Scope semantics eklendi (`sms.send`, `sms.*`, `*`) ve action capability/action_key/connector+action tabanlı eşleşme zorlandı.
   - Method/path guard eklendi:
     - Shared service ve grant seviyesinde `allowed_methods` + `allowed_path_rules` alanları desteklenir.
     - Request action method/path bu kurallara uymuyorsa policy deny edilir.
   - Shared usage observability eklendi:
     - Prometheus metric: `aystek_integrations_shared_service_usage_total`
       - label: `tenant_id`, `site_id`, `shared_service_id`, `grant_id`, `outcome`
     - Redis rollup key: `integrations:shared-service:usage_rollup:{service_id}:{grant_id}:{yyyymmdd}`
       - alanlar: `requests_total`, `outcome_*`, `tenant_id`, `site_id`, `updated_at`
   - Scope/method/path için ek servis testleri eklendi:
     - `TestDirectExecuteActionScopeSemanticsAllowsCapabilityRef`
     - `TestDirectExecuteActionMethodPathGuardDenied`
6. Phase-3 security hardening (backend core) tamamlandı:
   - Shared service/grant modeli genişletildi:
     - `secret_bindings`, `secret_values`, `payload_guard`, `mtls_policy`
     - grant override alanları: `secret_binding_overrides`, `payload_guard`, `mtls_policy`
   - DTO + handler + mongo repository katmanları yeni alanlarla güncellendi.
   - Servis normalize/validate kuralları eklendi:
     - secret binding/value doğrulaması
     - payload guard doğrulaması (`deny_fields`, `max_payload_bytes`)
     - mTLS policy doğrulaması (`required`, cert/key tutarlılığı)
   - Runtime execution enforcement eklendi:
     - Secret binding çözümü: grant override + service binding birleşimi ile effective connector auth config üretilir.
     - Payload guard: deny alanlar payload’dan recursive temizlenir; max payload byte aşımı policy deny eder.
     - mTLS hardening: grant/service policy birleşimi ile connector TLS enforce edilir; `required` durumda cert/key yoksa deny edilir.
   - Yeni error taxonomy eklendi:
     - `integrations.validation.shared_service_secret_binding_invalid`
     - `integrations.validation.shared_service_secret_missing`
     - `integrations.validation.shared_service_payload_guard_invalid`
     - `integrations.validation.shared_service_payload_denied`
     - `integrations.validation.shared_service_mtls_policy_invalid`
     - `integrations.validation.shared_service_mtls_required`
   - Handler error mapping güncellendi (new validation keys bad_request akışında).
   - Ek testler eklendi:
     - `TestEnforceSharedServicePolicyAppliesSecretBindings`
     - `TestEnforceSharedServicePolicyPayloadGuardSanitizesPayload`
     - `TestEnforceSharedServicePolicyRejectsMissingSecret`
     - `TestEnforceSharedServicePolicyRejectsWhenMTLSRequiredButMissing`
   - Doğrulama:
     - `go test ./internal/modules/integrations/service -run "SharedService|DirectExecuteAction|ProcessDueExecutions|TestHook|CreateSharedService|EnforceSharedService"`
     - `go test ./internal/modules/integrations/...`
7. Phase-3 UI/seed management katmanı tamamlandı:
   - Next admin tarafında shared service ekranları eklendi:
     - `/admin/integrations/shared-services` (list)
     - `/admin/integrations/shared-services/new` (create)
     - `/admin/integrations/shared-services/:id` (edit)
   - Next admin tarafında shared service grant ekranları eklendi:
     - `/admin/integrations/shared-service-grants` (list)
     - `/admin/integrations/shared-service-grants/new` (create)
     - `/admin/integrations/shared-service-grants/:id` (edit)
   - Form tarafında Phase-3 policy alanları yönetilebilir hale getirildi:
     - `secret_bindings`, `secret_values`, `payload_guard`, `mtls_policy`
     - grant tarafında `secret_binding_overrides`, `payload_guard`, `mtls_policy`
   - API istemcisi ve type sözleşmeleri shared service/grant endpointleri için genişletildi.
   - Translation seed anahtarları genişletildi:
     - shared service/grant field labels, status/protocol label’ları, form hint’leri
     - yeni Phase-3 validation key’leri (secret/payload/mTLS)
   - Doğrulama:
     - `npm -C /workspace/node/aystek-core-next run -s lint -- --file ...` (yeni dosyalar)
     - `go test ./internal/seeds -run TestDoesNotExist`

### In Progress
1. Yok (plan kapsamındaki zorunlu maddeler tamamlandı).

### Completed (Phase-4)
1. Shared service/grant observability API endpointleri eklendi:
   - `GET /admin/integrations/shared-services/:id/usage-summary?days=7|14|30`
   - `GET /admin/integrations/shared-service-grants/:id/usage-summary?days=7|14|30`
   - Redis usage rollup anahtarlarından (`integrations:shared-service:usage_rollup:*`) tenant/site scoped özet üretimi.
2. Shared service edit formuna UX hardening + observability kartları eklendi:
   - Policy hardening warning paneli (methods/path/scopes/payload guard/mTLS/secret binding kontrolleri).
   - Window (7/14/30 gün) seçimi ile toplam, today, denied, execution failed metrikleri.
   - Grant bazlı günlük kullanım yüzdesi tablosu ve quota yakınlık vurgusu.
3. Shared service grant edit formuna observability kartı eklendi:
   - Window (7/14/30 gün) seçimi, totals, denied ve daily usage ratio kartları.
   - Gün bazlı request/denied timeline tablosu.
4. Next API/type sözleşmeleri observability endpointleri için genişletildi.
5. Translation seed, Phase-4 UX anahtarları (hardening/usage/warnings) ile güncellendi.

### Next
1. Opsiyonel: placeholder görsellerin gerçek ekran capture'ları ile değiştirilmesi.

### Completed (Phase-5)
1. Runbook/ops checklist tamamlandı:
   - `user-guide/05-operations/INTEGRATIONS_SHARED_SERVICE_RUNBOOK_QA_MATRIX.md` eklendi.
   - `user-guide/05-operations/OPERATIONS_GUIDE.md` içine shared service smoke prosedürü eklendi.
2. User-guide observability bölümü tamamlandı:
   - `user-guide/06-modules/INTEGRATIONS_PLATFORM_GUIDE.md` içinde shared service/grant API+UI+ops referansları güncellendi.
   - `user-guide/06-modules/modules/integrations.md` latest state bölümüne Phase-5 observability/hardening özeti eklendi.
3. Ops script tamamlandı:
   - `user-guide/05-operations/scripts/integrations_shared_service_smoke.sh`
   - Service/grant usage-summary endpoint şema doğrulaması + health kontrolü içerir.
4. CI job template entegrasyonu tamamlandı:
   - `go/aystek-core/.gitlab-ci.yml` içine staging + production shared service smoke jobları eklendi.
   - Pipeline script path’i:
     - `go/aystek-core/tools/ops/integrations_shared_service_smoke.sh`
5. Operasyonel alarm eşiği matrisi tamamlandı:
   - `user-guide/05-operations/INTEGRATIONS_SHARED_SERVICE_ALERT_THRESHOLDS.md`
   - Shared service usage outcomes için PromQL alert rules + severity/tuning önerileri eklendi.
6. Walkthrough dokümantasyonu tamamlandı:
   - `user-guide/06-modules/INTEGRATIONS_SHARED_SERVICE_UI_WALKTHROUGH.md`
   - Shared service/grant ekranları için adım adım kullanıcı akışı ve triage kısa yolları eklendi.
   - Görsel asset paketi eklendi:
     - `user-guide/assets/integrations/shared-service-walkthrough/`
     - replace-ready placeholder görseller ve dosya adlandırma standardı tanımlandı.

## Gap Check (Source Doc vs Current Implementation)
1. Kritik bloklayıcı eksik: Yok.
2. Non-blocking enterprise backlog: Tamamlandı (2026-03-08 ikinci iterasyon).
   - Secret storage için native adapter eklendi:
     - `secret_provider`: `inline|vault|aws_kms`
     - `secret_provider_config` ile provider parametreleri yönetimi
     - runtime secret çözümleme katmanı (`shared_secret_adapter.go`)
   - Shared usage rollup Redis + Mongo historical aggregate hibritine taşındı:
     - yeni koleksiyon: `integration_shared_service_usage_agg`
     - Redis miss durumunda observability summary Mongo fallback okur.
   - Shared service policy contract/chaos testleri CI default profile’a alındı:
     - `.gitlab-ci.yml` içinde `integrations_contract_tests` ve `integrations_chaos_resilience_tests`
   - Grafana dashboard as-code provisioning eklendi:
     - `deployments/grafana/dashboards/integrations-shared-service.json`
     - `deployments/grafana/provisioning/dashboards/integrations-shared-service.yaml`
