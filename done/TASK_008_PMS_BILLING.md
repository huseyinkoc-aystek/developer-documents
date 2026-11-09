# TASK_008_PMS_BILLING

## Scope
Sprint-2 implementation for PMS billing bridge in AYSTEK CORE.

## Backend
- Added `internal/modules/pms_billing_bridge` with:
  - `model/model.go`
  - `dto/dto.go`
  - `repository/{interface.go,common.go,mongo.go}`
  - `service/{interface.go,service.go,service_test.go}`
  - `handler/{helpers.go,admin.go}`
  - `routes.go`
  - `indexes.go`
- Wired module in `cmd/api/main.go`
- Endpoints:
  - `POST /admin/pms/billing/dues/publish`
  - `POST /admin/pms/billing/payments/allocate`
  - `GET /admin/pms/billing/overview`
- Added idempotency support via `Idempotency-Key` and `pms_billing_idempotency` collection.
- Added outbox events:
  - `pms.dues.published`
  - `pms.charge.created`
  - `pms.payment.allocated`

## Frontend
- Extended PMS types/API:
  - `node/aystek-core-next/lib/types/pms.ts`
  - `node/aystek-core-next/lib/api/pms.ts`
- Replaced billing placeholder with working page:
  - `node/aystek-core-next/app/(dashboard)/admin/pms/billing/page.tsx`
- Features:
  - Period + property selection
  - Publish dues action (idempotency header sent)
  - Manual payment allocation (FIFO)
  - Overview table (unit, responsible, due/paid/remaining)

## i18n
- Added `pms.billing.*` keys and related validation/notice keys in:
  - `go/aystek-core/internal/seeds/translation_seed.go`
