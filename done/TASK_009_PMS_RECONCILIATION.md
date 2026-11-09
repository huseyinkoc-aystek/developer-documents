# TASK_009_PMS_RECONCILIATION

## Scope
Sprint-3 implementation for PMS bank reconciliation in AYSTEK CORE.

## Backend
- Added module: `internal/modules/pms_bank_reconciliation`
  - `model/model.go`
  - `dto/dto.go`
  - `repository/interface.go`
  - `repository/mongo.go`
  - `service/interface.go`
  - `service/service.go`
  - `service/service_test.go`
  - `handler/admin.go`
  - `handler/helpers.go`
  - `routes.go`
  - `indexes.go`
- Main wiring:
  - `cmd/api/main.go`

## Endpoints
- `POST /admin/pms/reconciliation/ingest`
- `GET /admin/pms/reconciliation/queue`
- `POST /admin/pms/reconciliation/:id/match`
- `POST /admin/pms/reconciliation/:id/confirm`
- `POST /admin/pms/reconciliation/:id/reject`

## Domain rules
- Duplicate guard on `external_txn_id` (tenant+site scoped unique index)
- State machine enforced:
  - `new -> matched`
  - `matched -> confirmed`
  - `matched -> rejected`
  - `rejected -> matched`
- Invalid transition returns HTTP 409
- Match proposal MVP:
  - amount closeness + description/counterparty vs party_ref heuristic
- Confirm flow:
  - create payment through core payments service
  - allocate via `pms_billing_bridge` service (FIFO already in bridge)
  - emit outbox events + audit logs

## Outbox events
- `pms.bank_txn.ingested`
- `pms.bank_match.proposed`
- `pms.bank_match.confirmed`
- `pms.bank_match.rejected`

## Frontend
- Updated API client/types:
  - `node/aystek-core-next/lib/api/pms.ts`
  - `node/aystek-core-next/lib/types/pms.ts`
- Replaced reconciliation placeholder:
  - `node/aystek-core-next/app/(dashboard)/admin/pms/reconciliation/page.tsx`
- UI features:
  - ingest form
  - queue list + status filter
  - propose/confirm/reject row actions
  - candidate preview column
  - posts-style header/stats/toolbar

## i18n
- Added `pms.reconciliation.*` keys in:
  - `go/aystek-core/internal/seeds/translation_seed.go`
