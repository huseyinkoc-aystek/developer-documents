# TASK_030_PMS_EXPORTS_RUNNER_HARDENING

## Scope
Sprint-11.C kapsamında PMS export job runner claim-lock ve retry davranışıyla güçlendirildi.

## Completed
- Export job modeli lock ve attempt alanlarıyla genişletildi:
  - `attempt_count`
  - `locked_at`, `locked_by`, `lock_expires_at`
- Repository’ye atomic claim eklendi:
  - `ClaimNextJob`
  - queued + stale-running reclaim
- Runner akışı queued işleri claim ederek işler hale getirildi.
- Çift çalıştırma riski claim lock ile azaltıldı.
- Dosya üretim/yazım hatalarında retry eklendi (2 deneme).
- Başarısızlıkta outbox `pms.export_job.failed` event’i korunarak yayınlanır.

## Files
- `go/aystek-core/internal/modules/pms/exports/model/model.go`
- `go/aystek-core/internal/modules/pms/exports/repository/interface.go`
- `go/aystek-core/internal/modules/pms/exports/repository/mongo.go`
- `go/aystek-core/internal/modules/pms/exports/service.go`
- `go/aystek-core/internal/modules/pms/exports/service_test.go`

## Validation
- `go test ./internal/modules/pms/... ./cmd/api`
