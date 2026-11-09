# TASK_020 — PMS Vendors Registry

## Teslim
- Yeni modül: `internal/modules/pms/vendors/*`
- Domain: vendor registry (tenant+site scoped, soft-delete/archive, audit + outbox)
- Endpointler:
  - `GET /admin/pms/vendors`
  - `POST /admin/pms/vendors`
  - `GET /admin/pms/vendors/:id`
  - `PUT /admin/pms/vendors/:id`
  - `POST /admin/pms/vendors/:id/archive`

## Olaylar
- `pms.vendor.created`
- `pms.vendor.updated`
- `pms.vendor.archived`

## Seed
- Menu: `menu.operations.pms.vendors`
- Permission module: `pms_vendors` (+ `archive` action)
- Translation keys: `pms.vendors.*`
