# TASK_032_PMS_ACCESS_SUITE

## Scope
Sprint-12.A kapsamında PMS Access suite eklendi: visitors + vehicle entries + PTS ingest endpoint.

## Completed
- Backend module: `internal/modules/pms/access/*`
  - Visitors:
    - `GET /admin/pms/access/visitors`
    - `POST /admin/pms/access/visitors`
    - `POST /admin/pms/access/visitors/:id/check-out`
  - Vehicle entries:
    - `GET /admin/pms/access/vehicles/entries`
    - `POST /admin/pms/access/vehicles/entries`
    - `POST /admin/pms/access/vehicles/entries/:id/exit`
    - `POST /admin/pms/access/vehicles/entries/ingest`
- Tenant+site scope repository seviyesinde zorunlu.
- Audit + outbox eventleri eklendi:
  - `pms.visitor.checked_in`
  - `pms.visitor.checked_out`
  - `pms.vehicle.entry_recorded`
  - `pms.vehicle.exit_recorded`
- Admin UI pages:
  - `/admin/pms/access/visitors`
  - `/admin/pms/access/vehicles`

## Seeds / i18n / RBAC
- `menu_seed.go`: access visitors/vehicles menüleri eklendi.
- `permission_seed.go`: `pms_access_visitors`, `pms_access_vehicles` modülleri + `ingest` aksiyonu.
- `role_seed.go`: editor/admin rol eşleştirmeleri güncellendi.
- `translation_seed.go`: `menu.operations.pms.access_*` ve `pms.access.*` keyleri eklendi.

## Validation
- `go test ./internal/modules/pms/... ./cmd/api`
- `npx eslint app/(dashboard)/admin/pms/access/visitors/page.tsx app/(dashboard)/admin/pms/access/vehicles/page.tsx lib/api/pms.ts lib/types/pms.ts`
