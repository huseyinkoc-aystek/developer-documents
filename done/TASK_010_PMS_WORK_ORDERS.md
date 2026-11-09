# TASK_010_PMS_WORK_ORDERS

## Scope
Sprint-4 implementation for PMS work orders in AYSTEK CORE.

## Backend
- Added module: `internal/modules/pms_work_orders`
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
- `GET /admin/pms/work-orders`
- `POST /admin/pms/work-orders`
- `GET /admin/pms/work-orders/:id`
- `PUT /admin/pms/work-orders/:id`
- `POST /admin/pms/work-orders/:id/assign`
- `POST /admin/pms/work-orders/:id/status`
- `POST /admin/pms/work-orders/:id/comment`
- `POST /admin/pms/work-orders/:id/close`

## Domain rules
- Tenant + site context is mandatory in service and repository layers.
- State machine enforced:
  - `open -> in_progress | blocked | cancelled`
  - `in_progress -> blocked | done | cancelled`
  - `blocked -> in_progress | cancelled`
  - `done/cancelled` terminal (except same-state no-op)
- Invalid status transitions return HTTP 409 from handler.
- Close action sets `closed_at` and emits close events.
- Timeline events are persisted in `pms_work_order_events`.

## Outbox events
- `pms.work_order.created`
- `pms.work_order.assigned`
- `pms.work_order.status_changed`
- `pms.work_order.commented`
- `pms.work_order.closed`

## Indexes
- `pms_work_orders__tenant_site_property_deleted`
- `pms_work_orders__tenant_site_status_priority_created_deleted`
- `pms_work_orders__tenant_site_unit_created_deleted`
- `pms_work_order_events__tenant_site_order_created`

## Frontend
- Updated API client/types:
  - `node/aystek-core-next/lib/api/pms.ts`
  - `node/aystek-core-next/lib/types/pms.ts`
- Replaced work-orders placeholder:
  - `node/aystek-core-next/app/(dashboard)/admin/pms/work-orders/page.tsx`
- UI features:
  - Posts/PMS style page header + stats cards
  - list with filters (search/status/priority/property/unit), pagination and actions
  - inline create card for new work order
  - assign/status/close row actions
  - detail timeline section (last events)

## i18n and Seeds
- Added `pms.work_orders.*` translations to:
  - `go/aystek-core/internal/seeds/translation_seed.go`
- Existing menu and permission seeds already include:
  - `menu.operations.pms.work_orders`
  - `pms_work_orders` module permissions (`read/create/update/delete/assign/close`)

