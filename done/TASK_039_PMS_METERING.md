# TASK 039 - PMS Metering & Utility Allocation

## Scope
Sprint-13.B kapsamında sayaç ve gider paylaştırma tarafı enterprise formatta hizalandı.

## Backend
- Module: `go/aystek-core/internal/modules/pms/metering`
- Meter type enum: `water|electric|heat|gas`
- Allocation method enum: `equal_split|by_consumption|mixed`
- Endpoints:
  - `GET/POST/PUT/DELETE /admin/pms/metering/meters`
  - `GET /admin/pms/metering/readings`
  - `POST /admin/pms/metering/readings/import`
  - `GET /admin/pms/metering/allocations`
  - `POST /admin/pms/metering/allocations/preview`
  - `POST /admin/pms/metering/allocations/publish`
- Publish side uses billing bridge path (no parallel charge engine).

## Frontend
- `/admin/pms/metering/meters`
  - gas type option added
- `/admin/pms/metering/allocations`
  - method values aligned with backend (`equal_split`, `by_consumption`, `mixed`)
  - preview/publish flow maintained

## Notes
- Legacy alias values (`equal`, `consumption`) backend’de normalize edilerek geriye uyum korunur.
