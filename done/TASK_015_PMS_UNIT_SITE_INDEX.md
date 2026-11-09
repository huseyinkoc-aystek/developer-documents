# TASK_015_PMS_UNIT_SITE_INDEX

## Scope
Sprint-7.B kapsamında deterministic site scoping için `pms_unit_site_index` read-model hardening uygulanmıştır.

## Architecture
Yeni read-model koleksiyonu:
- `pms_unit_site_index`
  - `tenant_id`
  - `site_id`
  - `unit_id`
  - `property_id`
  - `last_occupancy_id`
  - `updated_at`

## Implemented
- Index registration:
  - unique `(tenant_id, site_id, unit_id)`
  - `(tenant_id, site_id, property_id, updated_at)`
- Occupancy lifecycle hook:
  - occupancy create/update/delete sonrası `RebuildUnitSiteIndexForUnit` çağrısı
  - unit değişiminde eski ve yeni unit için rebuild
- Units bridge query hardening:
  - `/admin/pms/units` ve `/admin/pms/units/:id` scope önce `pms_unit_site_index` üzerinden
  - occupancy koleksiyonu yalnız snapshot amacıyla okunuyor

## Files
- `go/aystek-core/internal/modules/pms/model/pms.go`
- `go/aystek-core/internal/modules/pms/indexes.go`
- `go/aystek-core/internal/modules/pms/repository/interface.go`
- `go/aystek-core/internal/modules/pms/repository/mongo.go`
- `go/aystek-core/internal/modules/pms/service/service.go`
- `go/aystek-core/internal/modules/pms/service/service_test.go`
- `go/aystek-core/internal/modules/pms/units/bridge/repository/mongo.go`

## Validation
- `go test ./internal/modules/pms/... ./cmd/api`

## Notes
- Bu sprintte amaç deterministic scoping/read-model hardening olduğundan PMS dışı modüllere dokunulmadı.
