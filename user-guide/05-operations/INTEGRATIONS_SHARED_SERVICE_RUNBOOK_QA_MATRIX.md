---
Title: Integrations Shared Service Runbook & QA Matrix
Version: 1.0.0
Status: Active
Last Updated: 2026-03-08
Related Documents:
  - ./INTEGRATIONS_EVENTBRIDGE_SMOKE_RUNBOOK.md
  - ./INTEGRATIONS_CHAIN_RUNBOOK_QA_MATRIX.md
  - ./INTEGRATIONS_SCHEDULED_PULL_RUNBOOK_QA_MATRIX.md
  - ./OPERATIONS_GUIDE.md
  - ../06-modules/INTEGRATIONS_PLATFORM_GUIDE.md
---

# Integrations Shared Service Runbook & QA Matrix

Bu runbook, Shared Service + Shared Service Grant policy hattının güvenli ve izlenebilir çalıştığını doğrulamak için kullanılır.

## 1. Kapsam
1. Shared service/grant CRUD doğrulaması
2. Runtime policy enforcement (action/scope/method/path)
3. Secret binding + payload guard + mTLS policy davranışı
4. Usage rollup observability (`usage-summary`) ve warning paneli
5. Deny/quota olaylarında execution sonucu ve DLQ etkisi

## 2. Ön Koşullar
1. API, worker, RabbitMQ, Redis ayakta olmalı.
2. Tenant/site context geçerli olmalı.
3. Integrations admin yetkisi aktif olmalı.
4. En az bir connector + action tanımı hazır olmalı.
5. Test senaryosu için bir hook üzerinde `shared_service_grant_id` referansı tanımlanmalı.
6. Historical usage fallback için Mongo `integration_shared_service_usage_agg` koleksiyonu erişilebilir olmalı.

## 3. Smoke Adımları
### 3.1 Shared Service Setup
1. `Integrations > Shared Services` ekranından kayıt oluşturun.
2. `connector_id`, `allowed_action_ids`, `default_scopes`, `allowed_methods`, `allowed_path_rules` girin.
3. `default_limits` içinde en az `rpm` ve `daily_quota` tanımlayın.
4. Durumu `active` yapın.

Beklenen:
1. Kayıt başarılı döner.
2. Formda hardening warning paneli kritik eksik göstermiyorsa politika temel seviyesi sağlanır.

### 3.2 Grant Setup
1. `Integrations > Service Grants` ekranından grant oluşturun.
2. `shared_service_id`, `grant_key`, `allowed_action_ids`, `scopes` alanlarını doldurun.
3. Gerekirse grant seviyesinde `allowed_methods`, `allowed_path_rules`, `limits` override edin.
4. Durumu `active` yapın.

Beklenen:
1. Grant kaydı başarılı döner.
2. Service kapsamı dışındaki action seçimi validasyonda reddedilir.

### 3.3 Runtime Policy Positive Path
1. Hook üzerinde `config.shared_service_grant_id=<grant_id>` tanımlayın.
2. Hook action test veya gerçek execution tetikleyin.

Beklenen:
1. Execution `success`.
2. Response/body metadata içinde shared bağlamı görünür:
   - `shared_service_id`
   - `shared_service_grant_id`

### 3.4 Runtime Policy Negative Path
1. Scope veya method/path kuralına uymayan action tetikleyin.
2. Secret binding eksikliği veya payload max-size ihlali senaryosu çalıştırın.

Beklenen:
1. Execution `failed`.
2. `error_code=shared_service_policy_denied` veya ilgili validation kodu.
3. Hook DLQ açıksa retry sonrası DLQ kaydı oluşur.

### 3.5 Usage Observability
1. Shared service edit ekranında `Usage & Limit Observability` kartını açın.
2. 7/14/30 gün pencerelerini sırayla kontrol edin.
3. Grant edit ekranında grant usage kartını kontrol edin.

Beklenen:
1. Toplam/bugünkü/denied metrikleri dolu görünür.
2. Quota baskısı varsa warning panelinde görünür.
3. Gün bazlı tablo/trend verisi güncellenir.

## 4. QA Matrix
| Senaryo | Girdi | Beklenen |
|---|---|---|
| Service create active | valid connector + action allowlist + limits | kayıt başarılı, hardening paneli kritik uyarısız |
| Grant create active | valid service + scoped actions | kayıt başarılı |
| Action allowlist deny | grant dışı action | `shared_service_action_not_allowed` |
| Scope deny | `scopes` uyuşmuyor | policy denied |
| Method/path deny | `allowed_methods` veya `allowed_path_rules` mismatch | policy denied |
| Secret binding missing | binding var ama secret değeri yok | `shared_service_secret_missing` |
| Payload guard deny | payload boyutu `max_payload_bytes` aşar | `shared_service_payload_denied` |
| mTLS required missing cert/key | `required=true`, cert/key yok | `shared_service_mtls_required` |
| Usage summary | endpoint çağrısı | totals/outcomes/days dolu döner |
| Near quota | daily usage >=80% | warning panelinde quota yakınlık uyarısı |

## 5. API Doğrulama (Observability)
1. Service summary:
   - `GET /admin/integrations/shared-services/:id/usage-summary?days=7`
2. Grant summary:
   - `GET /admin/integrations/shared-service-grants/:id/usage-summary?days=7`

Minimum beklenen alanlar:
1. `totals.requests_total`
2. `totals.denied_total`
3. `outcomes`
4. `days[]`
5. `warnings[]`

## 6. Incident Triage Kısa Yol
1. Önce usage-summary warning kodlarına bakın.
2. Sonra execution error_code ve request preview incelenir.
3. Ardından grant/service policy alanları (scope/method/path/limits/secrets) doğrulanır.
4. Redis rollup anahtarlarının tenant/site eşleşmesi kontrol edilir.

## 7. Smoke Script (CI/CD)
Script:
- `user-guide/05-operations/scripts/integrations_shared_service_smoke.sh`
- `go/aystek-core/tools/ops/integrations_shared_service_smoke.sh` (pipeline runtime path)

Örnek:
```bash
chmod +x user-guide/05-operations/scripts/integrations_shared_service_smoke.sh
API_BASE_URL="https://<staging-api>" \
AUTH_TOKEN="$INTEGRATIONS_SMOKE_TOKEN" \
SHARED_SERVICE_ID="<service_id>" \
SHARED_SERVICE_GRANT_ID="<grant_id>" \
WINDOW_DAYS="7" \
./user-guide/05-operations/scripts/integrations_shared_service_smoke.sh
```

## 8. GitLab Job Örneği
```yaml
integrations_shared_service_smoke:
  stage: post_deploy
  image: alpine:3.20
  rules:
    - if: '$INTEGRATIONS_SMOKE_ENABLED == "true"'
  before_script:
    - apk add --no-cache bash curl jq
    - chmod +x user-guide/05-operations/scripts/integrations_shared_service_smoke.sh
  script:
    - API_BASE_URL="$STAGING_API_BASE_URL" AUTH_TOKEN="$INTEGRATIONS_SMOKE_TOKEN" SHARED_SERVICE_ID="$INTEGRATIONS_SHARED_SERVICE_ID" SHARED_SERVICE_GRANT_ID="$INTEGRATIONS_SHARED_SERVICE_GRANT_ID" ./user-guide/05-operations/scripts/integrations_shared_service_smoke.sh
```

## 9. CI Default Profile (Contract + Chaos)
Pipeline’da default çalışan ek doğrulama job’ları:
1. `integrations_contract_tests`
2. `integrations_chaos_resilience_tests`

Amaç:
1. Policy contract kırılmalarını erken yakalamak.
2. Resilience/chaos davranışını (`retry`, `circuit breaker`, `timeout`, `DLQ`) her merge’de doğrulamak.
