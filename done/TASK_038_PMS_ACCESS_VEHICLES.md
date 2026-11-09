# TASK 038 - PMS Access Vehicles

## Scope
Sprint-13.A araç giriş/çıkış ve PTS ingest hattı tamamlandı.

## Backend
- Model: `pms_vehicle_entries`
- Endpoints:
  - `GET /admin/pms/access/vehicles/entries`
  - `POST /admin/pms/access/vehicles/entries`
  - `POST /admin/pms/access/vehicles/entries/:id/exit`
  - `POST /admin/pms/access/vehicles/entries/ingest`
- Source enum: `manual|pts|import`
- Status enum: `in|out`
- HMAC signature doğrulaması (`X-Aystek-Signature`, `PMS_ACCESS_INGEST_SECRET`)
- Ingest route rate-limit korumalı.
- Events:
  - `pms.vehicle.entry_recorded`
  - `pms.vehicle.exit_recorded`

## Frontend
- Page: `/admin/pms/access/vehicles`
- Features:
  - Filter bar (plate + unit + status)
  - Entry create form (property/unit/plate/source)
  - Ingest action
  - Exit action
  - RBAC gating (`read/create/exit/ingest`)

## Notes
- Plate normalization backend’de uygulanır.
- Tenant/site scope repository sorgularında zorunlu.
