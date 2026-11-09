# TASK_034_PMS_PORTFOLIO

## Scope
Sprint-12.C kapsamında tenant-level portfolio görünümü eklendi.

## Completed
- Backend module: `internal/modules/pms/portfolio/*`
  - `GET /admin/pms/portfolio/dashboard?period=&site_id=`
  - `GET /admin/pms/portfolio/reports/collections?period=&site_id=`
- Portfolio aggregation tenant scoped çalışır; `site_id` verilirse site bazlı daraltır.
- Tenant izolasyonu korunur; cross-tenant erişim yok.
- Audit logları portfolio read aksiyonlarına eklendi.

## Admin UI
- `/admin/pms/portfolio/dashboard`
  - period + optional site filter
  - KPI kartları (open due, collection rate, open work orders, reconciliation queue, site count)
  - collections summary

## Seeds / i18n / RBAC
- `menu_seed.go`: portfolio menü girişi eklendi.
- `permission_seed.go`: `pms_portfolio` module.
- `role_seed.go`: admin/editor read scope güncellendi.
- `translation_seed.go`: `menu.operations.pms.portfolio`, `pms.portfolio.*` keyleri eklendi.

## Validation
- `go test ./internal/modules/pms/... ./cmd/api`
- `npx eslint app/(dashboard)/admin/pms/portfolio/dashboard/page.tsx lib/api/pms.ts lib/types/pms.ts`
