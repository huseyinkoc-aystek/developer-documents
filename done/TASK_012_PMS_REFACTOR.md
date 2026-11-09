# TASK_012_PMS_REFACTOR

## Scope
Sprint-6.A refactor: PMS modules consolidated under a single root module boundary.

## Structural changes
- Single root PMS module boundary kept at:
  - `internal/modules/pms`
- New PMS subpackages added:
  - `internal/modules/pms/properties`
  - `internal/modules/pms/occupancies`
  - `internal/modules/pms/billing/bridge`
  - `internal/modules/pms/reconciliation`
  - `internal/modules/pms/units/bridge`
  - `internal/modules/pms/work_orders`
- Added root files:
  - `internal/modules/pms/doc.go`
  - `internal/modules/pms/module.go`
  - `internal/modules/pms/routes.go`

## Routing and wiring
- `pms.RegisterRoutes` now acts as central PMS admin route registry.
- `cmd/api/main.go` PMS wiring moved to single init:
  - `pms.NewAdminModule(database, actSvc, paySvc)`
- Removed fragmented PMS route registration from `main.go`.

## Compatibility strategy
- Legacy packages preserved (`pms_billing_bridge`, `pms_bank_reconciliation`, `pms_units_bridge`, `pms_work_orders`).
- Added thin deprecated route wrappers in legacy modules that forward to new PMS-root route packages:
  - `internal/modules/pms_billing_bridge/routes.go`
  - `internal/modules/pms_bank_reconciliation/routes.go`
  - `internal/modules/pms_units_bridge/routes.go`
  - `internal/modules/pms_work_orders/routes.go`

## Contract and behavior
- Existing endpoint contracts remain unchanged (`/admin/pms/*`).
- Existing outbox/audit behavior preserved via underlying services.
- Seed/permission/menu chains remain compatible.

