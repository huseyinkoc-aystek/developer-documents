# TASK_013_PMS_PORTAL

## Scope
Sprint-6.B portal endpoints implemented under `pms/portal`.

## New package
- `internal/modules/pms/portal`
  - `model.go`
  - `repository.go`
  - `service.go`
  - `handler.go`
  - `routes.go`

## Portal routes
Registered via `pms.RegisterPortalRoutes(protectedGroup, portalModule)`:
- `GET /portal/pms/me/units`
- `GET /portal/pms/me/dues`
- `GET /portal/pms/me/payments`
- `POST /portal/pms/me/work-orders`

## Access model
- Tenant + site context required (derived from middleware context).
- `me` identity resolved from authenticated context (`user_id`/`userID`).
- Work order create requires active occupancy ownership on the target unit:
  - deny with `pms.portal.error.unit_access_denied` when not eligible.

## Integration notes
- No bypass of billing/payments/units core collections.
- Work order creation is delegated to existing work-order service (audit + outbox preserved).
- Added i18n keys for portal validation/error contracts:
  - `pms.portal.*`

