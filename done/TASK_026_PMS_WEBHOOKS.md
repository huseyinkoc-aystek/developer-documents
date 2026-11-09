# TASK 026 — PMS Webhooks (External Integrations)

## Scope
Implemented PMS webhook subscriptions and event push surface for external integrations.

## Backend Module
- Added package: `internal/modules/pms/webhooks`
  - `model.go`
  - `dto.go`
  - `repository.go`
  - `service.go`
  - `handler.go`
  - `routes.go`
  - `indexes.go`
  - `service_test.go`
- Added dispatcher bridge package:
  - `internal/modules/pms/webhooks/dispatch/dispatch.go`

## Endpoints
- `GET /admin/pms/webhooks`
- `GET /admin/pms/webhooks/:id`
- `POST /admin/pms/webhooks`
- `PUT /admin/pms/webhooks/:id`
- `POST /admin/pms/webhooks/:id/test`

## Security
- HMAC-SHA256 signature header: `X-Aystek-Signature`
- Secret is masked in API responses.
- Tenant+site scoped CRUD/read in repository.

## Delivery Behavior (MVP)
- Subscription `test` endpoint performs real signed HTTP delivery.
- Runtime dispatch method implemented and wired to PMS event emission points.
- Delivery records stored with retry metadata:
  - `retry_count`
  - `next_attempt_at`
  - `last_error`
  - status: `pending|sent|failed|dead_letter`

## Event Coverage (wired)
- `pms.payment.allocated`
- `pms.bank_match.confirmed`
- `pms.work_order.closed`
- `pms.vendor.archived`

## Wiring
- `internal/modules/pms/module.go` now initializes webhook module and dispatcher adapter.
- `internal/modules/pms/routes.go` registers webhook admin routes.

## Seeds
- Permission module added: `pms_webhooks` (with default action set)
- Menu key added: `menu.operations.pms.webhooks`
- Translation keys added: `pms.webhooks.*`
