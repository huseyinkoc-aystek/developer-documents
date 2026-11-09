# TASK_007_PMS

## Scope
Sprint-1 implementation for PMS in AYSTEK CORE:

- `pms_properties`
- `pms_occupancies`
- responsibility resolver in service layer
- admin routes under `/admin/pms/*`
- seed updates (menu/permissions/roles/translations/domain seed)
- Next.js admin pages for properties and occupancies

## Backend

### Module
`go/aystek-core/internal/modules/pms`

### Layers
- model: property + occupancy + outbox model
- dto: admin request payloads
- repository: Mongo implementation with strict tenant/site/deleted filters
- service: validation, overlap guard, responsibility resolution, outbox emission, activity logs
- handler: thin HTTP mapping
- routes: RBAC-guarded admin endpoints
- indexes: compound + unique partial indexes

### Endpoints
- `GET /admin/pms/properties`
- `POST /admin/pms/properties`
- `GET /admin/pms/properties/:id`
- `PUT /admin/pms/properties/:id`
- `DELETE /admin/pms/properties/:id`

- `GET /admin/pms/occupancies`
- `POST /admin/pms/occupancies`
- `GET /admin/pms/occupancies/:id`
- `PUT /admin/pms/occupancies/:id`
- `DELETE /admin/pms/occupancies/:id`
- `GET /admin/pms/occupancies/responsibility?unit_id=<id>&at=<RFC3339>`

### Rules implemented
- Tenant + site required from context in service/repo
- Soft delete on property/occupancy
- Property code unique per tenant+site
- Occupancy overlap prevention for same unit/date range
- Responsibility resolver priority:
  1. financial_responsible=true
  2. tenant fallback
  3. owner fallback
  4. first active occupancy

### Outbox events
- `pms.property.created|updated|archived`
- `pms.occupancy.created|updated|ended`

## Frontend

### Routes
- `/admin/pms/properties`
- `/admin/pms/properties/new`
- `/admin/pms/properties/[id]`
- `/admin/pms/occupancies`
- `/admin/pms/occupancies/new`
- `/admin/pms/occupancies/[id]`

Added placeholders:
- `/admin/pms/units`
- `/admin/pms/billing`
- `/admin/pms/reconciliation`
- `/admin/pms/work-orders`

### UI
- Posts-style gradient page header
- stats strip + filter toolbar + table lists
- dedicated create/edit pages (no modal forms)
- i18n key usage in all new pages

## Seeds
- `menu_seed.go`: operations > pms parent and child entries
- `permission_seed.go`: pms module permissions
- `role_seed.go`: editor read visibility for pms modules
- `translation_seed.go`: pms menu + field + validation keys (28-lang via `allLangs`)
- `pms_seed.go`: starter sample property seed
- `runner.go`: PMS seed runner integration

## Wiring
- `cmd/api/main.go` module import + repo/service/handler wiring + route registration

## Tests
- `internal/modules/pms/service/service_test.go`
  - overlap validation
  - responsibility selection
  - context requirement behavior
