# TASK_033_PMS_METERING

## Scope
Sprint-12.B kapsamında metering + utility allocation modülü eklendi.

## Completed
- Backend module: `internal/modules/pms/metering/*`
  - Meters CRUD:
    - `GET /admin/pms/metering/meters`
    - `POST /admin/pms/metering/meters`
    - `PUT /admin/pms/metering/meters/:id`
    - `DELETE /admin/pms/metering/meters/:id`
  - Readings import:
    - `POST /admin/pms/metering/readings/import`
  - Allocation:
    - `GET /admin/pms/metering/allocations`
    - `POST /admin/pms/metering/allocations/preview`
    - `POST /admin/pms/metering/allocations/publish`
- Publish, billing bridge `CreateManualCharge` ile orkestrasyon yapar (parallel charge sistemi yok).
- Allocation method desteği:
  - `equal`
  - `consumption` (okuma verisi varsa tüketim ağırlıklı)
- Outbox eventleri:
  - `pms.meter.reading_ingested`
  - `pms.utility_allocation.published`

## Admin UI
- `/admin/pms/metering/meters`
- `/admin/pms/metering/allocations`

## Seeds / i18n / RBAC
- `menu_seed.go`: metering menüsü eklendi.
- `permission_seed.go`: `pms_metering` module + `ingest`, `publish` aksiyonları.
- `role_seed.go`: editor/admin kapsamı güncellendi.
- `translation_seed.go`: `pms.metering.*` key seti eklendi.

## Validation
- `go test ./internal/modules/pms/... ./cmd/api`
- `npx eslint app/(dashboard)/admin/pms/metering/meters/page.tsx app/(dashboard)/admin/pms/metering/allocations/page.tsx lib/api/pms.ts lib/types/pms.ts`
