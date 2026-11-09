# TASK 036 — PMS Legal Docs

## Scope
- Added PMS legal endpoints:
  - `GET /admin/pms/legal/receipt/:payment_id`
  - `GET /admin/pms/legal/ledger?period=YYYY-MM&format=csv`

## Architecture Notes
- Implemented as dedicated module under PMS boundary:
  - `internal/modules/pms/legal/*`
- Uses tenant+site scoped repository access.
- Keeps export behavior lightweight with CSV output for ledger.

## Security & Compliance
- RBAC guarded with `pms_legal:read`.
- Audit logs added:
  - `pms.legal.receipt.view`
  - `pms.legal.ledger.export`

