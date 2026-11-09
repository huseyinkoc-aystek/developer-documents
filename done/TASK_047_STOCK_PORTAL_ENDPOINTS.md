# TASK_047 — STOCK Portal Endpointleri

## Kapsam
- Planlanan `portal/stock` yüzeyi backend ve Next typed client tarafında açıldı.

## Backend
- Eklenen route’lar:
  - `GET /portal/stock/sales`
  - `POST /portal/stock/sales`
  - `POST /portal/stock/field-reports/assets`
  - `GET /portal/stock/locations/current`
- Dosyalar:
  - `internal/modules/stock/sales/handler/portal.go`
  - `internal/modules/stock/sales/routes.go`
  - `internal/modules/stock/sales/repository/interface.go`
  - `internal/modules/stock/sales/repository/mongo.go`
  - `internal/modules/stock/locations/handler/portal.go`
  - `internal/modules/stock/locations/routes.go`
  - `cmd/api/main.go` (portal route register)

## Next.js (typed client)
- `lib/api/stock-ecosystem.ts`:
  - `listPortalSales`
  - `createPortalSale`
  - `createPortalFieldReportAsset`
  - `listPortalCurrentLocations`

## Doğrulama
- `go test ./internal/modules/stock/... ./cmd/api` ✅
- `npx eslint lib/api/stock-ecosystem.ts` ✅

