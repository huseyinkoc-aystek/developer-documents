# TASK_011_PMS_UNITS_BRIDGE

## Scope
Sprint-5 implementation for PMS Units Bridge in AYSTEK CORE.

## Backend
- Added module: `internal/modules/pms_units_bridge`
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
  - `indexes.go` (no dedicated collection; read facade)
- Main wiring:
  - `cmd/api/main.go`

## Endpoints
- `GET /admin/pms/units?property_id=&q=&page=&limit=`
- `GET /admin/pms/units/:id`

## Bridge behavior
- No unit model duplication.
- Uses existing collections as read facade:
  - `units`
  - `pms_occupancies`
  - `pms_properties`
  - `pms_charges`
  - `pms_payment_allocations`
  - `pms_work_orders`
- Detail response includes:
  - occupancy snapshot
  - dues summary
  - last payment summary
  - open work orders count
  - recent work orders list

## Security and scope
- Strict tenant + site context required in service/repository.
- Route permission:
  - module: `pms_units`
  - action: `read`

## Seeds / i18n
- Added permission module seed entry:
  - `pms_units` in `permission_seed.go`
- Added translation keys:
  - `pms.units.*` in `translation_seed.go`

## Frontend
- Replaced placeholder page:
  - `node/aystek-core-next/app/(dashboard)/admin/pms/units/page.tsx`
- Added PMS unit bridge API mapping:
  - `listPmsUnits`
  - `getPmsUnit`
- Extended PMS types:
  - `PMSUnitListItem`
  - `PMSUnitDetail`
  - related snapshot/summary types

