# TASK_042 — STOCK Sprint-1 Backend Foundation

## Kapsam
- Domain foundation: `stock/catalog`, `stock/inventory`
- Tenant + site scoped admin API (context tabanlı)
- Index registration
- Seed zinciri (menu + permission + domain seed)
- Servis testleri (validation + transition)
- Next.js admin list/create/edit iskeleti

## Uygulanan Dosyalar
- `go/aystek-core/internal/modules/stock/catalog/*`
- `go/aystek-core/internal/modules/stock/inventory/*`
- `go/aystek-core/internal/seeds/stock_seed.go`
- `go/aystek-core/internal/seeds/menu_seed.go`
- `go/aystek-core/internal/seeds/permission_seed.go`
- `go/aystek-core/internal/seeds/translation_seed.go`
- `go/aystek-core/internal/seeds/runner.go`
- `go/aystek-core/cmd/api/main.go`
- `node/aystek-core-next/lib/types/stock-ecosystem.ts`
- `node/aystek-core-next/lib/api/stock-ecosystem.ts`
- `node/aystek-core-next/app/(dashboard)/admin/stock/components/StockAdminHeader.tsx`
- `node/aystek-core-next/app/(dashboard)/admin/stock/catalog/brands/*`
- `node/aystek-core-next/app/(dashboard)/admin/stock/catalog/models/*`
- `node/aystek-core-next/app/(dashboard)/admin/stock/inventory/items/*`

## API Yüzeyi (Sprint-1)
- `GET /admin/stock/catalog/brands`
- `POST /admin/stock/catalog/brands`
- `GET /admin/stock/catalog/brands/:id`
- `PUT /admin/stock/catalog/brands/:id`
- `DELETE /admin/stock/catalog/brands/:id`
- `GET /admin/stock/catalog/models`
- `POST /admin/stock/catalog/models`
- `GET /admin/stock/catalog/models/:id`
- `PUT /admin/stock/catalog/models/:id`
- `DELETE /admin/stock/catalog/models/:id`
- `GET /admin/stock/inventory/items`
- `POST /admin/stock/inventory/items`
- `GET /admin/stock/inventory/items/:id`
- `PUT /admin/stock/inventory/items/:id`
- `DELETE /admin/stock/inventory/items/:id`
- `POST /admin/stock/inventory/items/:id/status`

## Validasyon ve Kurallar
- `tenant_id` + `site_id` zorunlu (context kaynaklı)
- Brand/Model duplicate adı engeli
- Product item serial/imei duplicate engeli
- Item status transition matrix guard
- Soft delete filtrelemesi default açık

## Test Durumu
- `go test ./internal/modules/stock/catalog/...` ✅
- `go test ./internal/modules/stock/inventory/...` ✅
- `go test ./internal/seeds ./cmd/api` ✅
- `npx eslint app/(dashboard)/admin/stock/catalog/brands/**/*.tsx app/(dashboard)/admin/stock/catalog/models/**/*.tsx app/(dashboard)/admin/stock/inventory/items/**/*.tsx app/(dashboard)/admin/stock/components/StockAdminHeader.tsx lib/api/stock-ecosystem.ts lib/types/stock-ecosystem.ts` ✅

## Sonraki Adım
- Sprint-2 kapsamında `stock_vehicles`, `stock_assignments`, `stock_sales` modüllerinin aynı root altında (`internal/modules/stock/...`) açılması.
