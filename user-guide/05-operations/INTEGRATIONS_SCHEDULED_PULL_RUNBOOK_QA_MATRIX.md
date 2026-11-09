---
Title: Integrations Scheduled Pull Runbook & QA Matrix
Version: 1.0.0
Status: Active
Last Updated: 2026-03-08
Related Documents:
  - ./INTEGRATIONS_EVENTBRIDGE_SMOKE_RUNBOOK.md
  - ./INTEGRATIONS_CHAIN_RUNBOOK_QA_MATRIX.md
  - ./OPERATIONS_GUIDE.md
  - ../06-modules/INTEGRATIONS_PLATFORM_GUIDE.md
---

# Integrations Scheduled Pull Runbook & QA Matrix

Bu runbook, zamanlı dış veri çekme (scheduled pull) hattını uçtan uca doğrulamak için kullanılır.

## 1. Kapsam
1. Pull schedule oluşturma ve güncelleme
2. Scheduler -> run üretimi
3. Worker yürütmesi, retry/backoff, DLQ
4. Ingestion summary/issues görünürlüğü
5. Cursor/checkpoint success-only commit davranışı

## 2. Ön Koşullar
1. API ve worker ayakta olmalı.
2. RabbitMQ bağlantısı healthy olmalı.
3. Integrations admin yetkisi aktif olmalı.
4. Connector + action tanımı hazır olmalı.

## 3. Smoke Adımları
### 3.1 Schedule Kurulumu
1. `Integrations > Scheduled Pulls` ekranından yeni schedule oluşturun.
2. `connector_id`, `action_id`, `trigger_mode`, `incremental_strategy` alanlarını doldurun.
3. Durumu `active` yapın ve kaydedin.

Beklenen:
1. Kayıt başarılı döner.
2. Liste ekranında schedule görünür.

### 3.2 Test-Run
1. Schedule detay ekranında test payload girin.
2. `dry_run=true` ile test çalıştırın.
3. Sonucu detay panelinde kontrol edin.

Beklenen:
1. Run kaydı `trigger_type=test` ile oluşur.
2. Response metadata dolu döner.

### 3.3 Runtime Scheduler + Worker
1. Scheduler tick beklenir veya internal due endpoint tetiklenir.
2. `Integrations > Pull Runs` ekranından yeni run kayıtları izlenir.
3. Bir run detail ekranına girip ingestion kartlarını kontrol edin.

Beklenen:
1. Run status `success` veya kontrollü senaryoda `failed`.
2. Ingestion summary kartları dolu gelir.
3. Ingestion issues tablosu beklenen issue’ları gösterir.

### 3.4 Retry / DLQ
1. Harici endpoint’i geçici olarak başarısız döndürecek şekilde test edin.
2. Run’ın retry_count ve next_retry_at alanlarını takip edin.
3. Retry limit aşımı sonrası DLQ kaydını kontrol edin.

Beklenen:
1. Backoff artan şekilde uygulanır.
2. Limit sonrası run DLQ semantiğine düşer.

### 3.5 Cursor Success-Only Commit
1. Başarılı bir run sonrası cursor state’i kontrol edin.
2. Başarısız run tetikleyin.
3. Cursor state’in başarısız run’da değişmediğini doğrulayın.

Beklenen:
1. Sadece success run sonrası cursor_after commit edilir.

## 4. QA Matrix
| Senaryo | Girdi | Beklenen |
|---|---|---|
| Cron schedule success | active + cron + valid action | success run, ingestion summary dolu |
| Interval schedule success | active + interval + valid action | success run, records_fetched/upserted artar |
| Dry-run test | test-run dry_run=true | external side effect olmadan preview/result |
| Validation fail | invalid strategy_config | API validation error, kayıt yok |
| Downstream HTTP fail | action endpoint 5xx | retry_count artar, limitte DLQ |
| Ingestion issue | schema mismatch | run kırılmadan issues listesine yazım |
| Cursor commit | success -> fail sequence | cursor sadece success sonrası güncellenir |

## 5. Operasyonel Alarm Önerileri
1. `pull_ingestion_issues_total` ani artış alarmı
2. `pull_ingestion_events_published_total` beklenmeyen düşüş alarmı
3. Pull run `failed` oranı eşik alarmı
4. `queued` run birikimi alarmı

## 6. Incident Triage Kısa Yol
1. Önce run detail ekranından `response_meta`, `warnings`, `ingestion_issues` bakın.
2. Sonra connector/action ve allowlist/SSRF policy kontrolü yapın.
3. Ardından retry/backoff ve DLQ akışını doğrulayın.

## 7. Smoke Script (CI/CD)
Script:
- `user-guide/05-operations/scripts/integrations_pull_smoke.sh`

Örnek kullanım:
```bash
chmod +x user-guide/05-operations/scripts/integrations_pull_smoke.sh
API_BASE_URL="https://<staging-api>" \
AUTH_TOKEN="$INTEGRATIONS_SMOKE_TOKEN" \
SCHEDULE_KEY="MY_PULL_SCHEDULE" \
DRY_RUN="true" \
EXPECTED_FINAL_STATUS="success" \
./user-guide/05-operations/scripts/integrations_pull_smoke.sh
```

GitLab job örneği:
```yaml
integrations_pull_smoke:
  stage: post_deploy
  image: alpine:3.20
  rules:
    - if: '$INTEGRATIONS_SMOKE_ENABLED == "true"'
  before_script:
    - apk add --no-cache bash curl jq
    - chmod +x user-guide/05-operations/scripts/integrations_pull_smoke.sh
  script:
    - API_BASE_URL="$STAGING_API_BASE_URL" AUTH_TOKEN="$INTEGRATIONS_SMOKE_TOKEN" SCHEDULE_KEY="$INTEGRATIONS_PULL_SMOKE_SCHEDULE_KEY" ./user-guide/05-operations/scripts/integrations_pull_smoke.sh
```
