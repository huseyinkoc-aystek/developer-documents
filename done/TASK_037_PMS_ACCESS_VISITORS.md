# TASK 037 - PMS Access Visitors

## Scope
Sprint-13.A kapsamında PMS Access suite altında ziyaretçi giriş/çıkış akışı enterprise standardında tamamlandı.

## Backend
- Module: `go/aystek-core/internal/modules/pms/access`
- Model: `pms_visitors`
- Endpoints:
  - `GET /admin/pms/access/visitors`
  - `POST /admin/pms/access/visitors`
  - `POST /admin/pms/access/visitors/:id/check-out`
- List filtreleri: `unit_id`, `status`, `q`, pagination.
- Status lifecycle: `in -> out`.
- Tenant/site scope: service + repository seviyesinde zorunlu.
- Audit events:
  - `pms.access.visitor.checkin`
  - `pms.access.visitor.checkout`
- Outbox events:
  - `pms.visitor.checked_in`
  - `pms.visitor.checked_out`

## Frontend
- Page: `/admin/pms/access/visitors`
- Features:
  - Filter bar (search + unit + status)
  - Check-in form (property/unit/name/phone/plate/note)
  - List with status + check-out action
  - RBAC gating (`read/create/checkout`)

## Notes
- No core bypass; only PMS boundary.
- i18n keys seeded under `pms.access.*`.
