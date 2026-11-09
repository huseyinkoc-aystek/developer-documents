# TASK 025 — PMS Exports Center (Job-based)

## Scope
Implemented enterprise PMS Exports Center as canonical export surface.

## Backend
- New module surface under existing `internal/modules/pms/exports`:
  - `service.go`
  - `handler.go`
  - `routes.go`
  - `indexes.go`
  - `service_test.go`
- Existing export data layer completed:
  - `model/model.go`
  - `dto/dto.go`
  - `repository/interface.go`
  - `repository/mongo.go`

## API
- `POST /admin/pms/exports/jobs`
- `GET /admin/pms/exports/jobs`
- `GET /admin/pms/exports/jobs/:id`
- `GET /admin/pms/exports/jobs/:id/download`

## Behavior
- Job types:
  - `accounting_csv`
  - `collections_csv`
  - `bank_txns_csv`
- Job statuses:
  - `queued`
  - `running`
  - `done`
  - `failed`
- In-process async runner (MVP) processes created jobs.
- Export file bytes persisted in `pms_export_files` and served by file ref.
- Outbox events emitted:
  - `pms.export_job.created`
  - `pms.export_job.completed`
  - `pms.export_job.failed`
- Audit logs:
  - job create
  - job download

## Multi-tenant & Site Scope
- Repository read/write operations require tenant+site context.
- Download guarded by same tenant/site constraints.

## Indexes
- `pms_export_jobs__tenant_site_status_created`
- `pms_export_jobs__tenant_site_type_period_created`
- `pms_export_files__tenant_site_ref` (unique)

## Frontend
- New canonical page:
  - `/admin/pms/exports`
  - file: `node/aystek-core-next/app/(dashboard)/admin/pms/exports/page.tsx`
- Features:
  - create export job (type + period)
  - list/filter jobs (type + status)
  - download finished jobs
- Dashboard keeps quick button and now links canonical Exports Center.

## Seeds / i18n
- Menu entries:
  - `menu.operations.pms.exports`
- Permissions:
  - `pms_exports:download` added
- Translation keys:
  - `pms.exports.*` (title/subtitle/types/status/actions/validation/notice)

## Validation
- `go test ./internal/modules/pms/... ./cmd/api`
- `npx eslint app/(dashboard)/admin/pms/exports/page.tsx lib/api/pms.ts lib/types/pms.ts`
