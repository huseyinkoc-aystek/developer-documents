# TASK 035 — PMS Portal Card Payments

## Scope
- Added portal payment initiation endpoint inside PMS boundary:
  - `POST /portal/pms/me/payments/initiate_card_payment`
- Added finalize endpoint to bridge success status to PMS allocation flow:
  - `POST /portal/pms/me/payments/intents/:intent_id/finalize`

## Architecture Notes
- Kept business logic under `internal/modules/pms/portal`.
- Reused core `payments/intents` service for payment intent creation.
- Reused PMS billing bridge allocation flow (no parallel charge/payment engine).
- Added outbox event emission for `pms.payment.received` after successful finalize.

## Security & Isolation
- Tenant/site context is required from middleware context.
- Unit access is validated with active occupancy (`me` user + unit guard).
- Idempotency key support via `Idempotency-Key` header.

## Audit
- Added portal action logs:
  - `pms.portal.payments.initiated`
  - `pms.portal.payments.finalized`

