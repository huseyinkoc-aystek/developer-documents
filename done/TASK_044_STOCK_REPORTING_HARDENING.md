# TASK_044 — STOCK Sprint-7 Hardening (Reporting/Exports)

## Kapsam
- `stock_reporting` modülünde Sprint-7 hardening tamamlandı:
  - CSV export formula-injection sanitize
  - Export job asenkron claim/lock/retry akışı
  - Reporting cache/read-model (tenant+site scoped)

## Backend Değişiklikleri
- Model:
  - `ExportJob` alanları genişletildi: `attempt_count`, `max_attempts`, `next_run_at`, `locked_by`, `locked_at`, `lock_expires_at`, `query`
- Repository:
  - `ClaimExportJobs` ile atomic claim (queued/running-expired)
  - Cache metodları eklendi:
    - `Get/SetActiveSubscriptionsCache`
    - `Get/SetVehicleAssetsCache`
  - Yeni collection: `stock_report_read_models`
- Service:
  - `CreateExportJob` artık senkron export yerine queued job oluşturuyor
  - Request context’ten tenant/site korunarak async işleme tetikleniyor
  - Retry backoff (attempt bazlı) eklendi
  - CSV yazımı sanitize edildi (`=`, `+`, `-`, `@` önekli hücreler)
  - `q` filtresi job modelinde saklanıp export üretiminde kullanılıyor

## Indexler
- `stock_report_jobs__tenant_site_status_next_run`
- `stock_report_read_models__tenant_site_type_key` (unique)
- `stock_report_read_models__expires_at_ttl` (TTL)

## Testler
- `service_test.go` güncellendi:
  - queued -> claimed -> done lifecycle
  - CSV sanitize kontrolü
- Doğrulama:
  - `go test ./internal/modules/stock/reporting/... ./cmd/api` ✅

## Notlar
- API kontratı geriye uyumlu bırakıldı (`/admin/stock/reports/exports/jobs` yüzeyi değişmedi).
- UI tarafında jobs listesi mevcut flow ile queued/running/done durumlarını okumaya devam eder.
