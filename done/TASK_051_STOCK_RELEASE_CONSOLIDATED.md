# TASK_051 — STOCK Sprint Consolidated Release (S1–S7 Hardening)

## Kapsam
Bu doküman `TASK_042..TASK_050` aralığındaki stock teslimatlarının tek release görünümüdür.

## Tamamlanan Modüller
- `stock_catalog`
- `stock_inventory`
- `stock_vehicles`
- `stock_vehicle_assets`
- `stock_assignments`
- `stock_sales`
- `stock_returns`
- `stock_subscriptions`
- `stock_integrations`
- `stock_locations`
- `stock_sms`
- `stock_reporting`
- legacy `stock/export` parity hardening

## Sprint Bazlı Özet
- Sprint-1: Catalog + Inventory foundation (`TASK_042`)
- Sprint-2: Assignments + vehicles/assets ilerlemesi (`TASK_043`, `TASK_046`)
- Sprint-3/4/5/6 kapsamı: sales/returns/subscriptions/integrations/locations/sms/reporting tamamlandı (modül bazlı kod ve seed zinciri aktif)
- Sprint-7 hardening:
  - reporting/export async runner + CSV sanitize (`TASK_044`)
  - legacy export parity (`TASK_045`)
  - portal stock endpointleri (`TASK_047`)
  - sms callback HMAC+replay (`TASK_048`)
  - security negative tests expansion (`TASK_049`)
  - HTTP scope hardening (`TASK_050`)

## Endpoint Durumu (Özet)
- Admin `stock/*` endpointleri aktif
- Portal:
  - `GET/POST /portal/stock/sales`
  - `POST /portal/stock/field-reports/assets`
  - `GET /portal/stock/locations/current`

## Security & Hardening
- CSV export formula injection sanitize
- Async export job locking/claim parity
- SMS callback:
  - HMAC signature verification
  - replay protection (nonce)
  - timestamp window validation
- Tenant/site context zorunluluğu:
  - service-level guards
  - admin handler-level guards (fallback kaldırıldı)

## Seed / RBAC / i18n
- `menu_seed.go`, `permission_seed.go`, `role_seed.go`, `translation_seed.go`, `stock_seed.go` zinciri aktif
- Stock altı menüler (inventory -> reports) ve modül izinleri eşlendi

## Test Durumu
- `go test ./internal/modules/stock/... ./cmd/api` yeşil
- Security/handler negatif testler dahil:
  - `stock/security`
  - `stock/handler`
  - `stock/sales/handler`

## Referans Dokümanlar
- `TASK_042_STOCK_SPRINT1.md`
- `TASK_043_STOCK_ASSIGNMENTS.md`
- `TASK_044_STOCK_REPORTING_HARDENING.md`
- `TASK_045_STOCK_LEGACY_EXPORT_PARITY.md`
- `TASK_046_STOCK_VEHICLE_ASSETS.md`
- `TASK_047_STOCK_PORTAL_ENDPOINTS.md`
- `TASK_048_STOCK_SMS_CALLBACK_HARDENING.md`
- `TASK_049_STOCK_SECURITY_NEGATIVE_TESTS.md`
- `TASK_050_STOCK_HTTP_SCOPE_HARDENING.md`

