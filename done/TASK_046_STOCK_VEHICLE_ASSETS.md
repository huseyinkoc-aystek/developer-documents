# TASK_046 — STOCK Sprint-2 (Vehicle Assets)

## Kapsam
- `internal/modules/stock/vehicle_assets` bounded context eklendi.
- Araç üstü ürün montaj/söküm yaşam döngüsü tenant+site scoped olarak açıldı.
- Admin UI eklendi: `/admin/stock/vehicle-assets`, `/new`, `/[id]`.
- Seed/RBAC/i18n zinciri tamamlandı.

## Backend
- Modül dosyaları:
  - `model`, `dto`, `repository (interface+mongo)`, `service`, `service_test`, `handler`, `routes`, `indexes`
- API yüzeyi:
  - `GET /admin/stock/vehicle-assets`
  - `POST /admin/stock/vehicle-assets`
  - `GET /admin/stock/vehicle-assets/:id`
  - `PUT /admin/stock/vehicle-assets/:id`
  - `POST /admin/stock/vehicle-assets/:id/unmount`
  - `DELETE /admin/stock/vehicle-assets/:id`
- İş kuralları:
  - `vehicle_id` ve `product_item_id` zorunlu.
  - Aynı ürünün eşzamanlı olarak birden fazla araçta aktif olması engellenir.
  - Aynı araç+ürün eşleşmesinin tekrar aktif açılması engellenir.
  - Unmount sonrası kayıt `unmounted` durumuna alınır.

## Next.js Admin
- Sayfalar:
  - `app/(dashboard)/admin/stock/vehicle-assets/page.tsx`
  - `app/(dashboard)/admin/stock/vehicle-assets/new/page.tsx`
  - `app/(dashboard)/admin/stock/vehicle-assets/[id]/page.tsx`
  - `app/(dashboard)/admin/stock/vehicle-assets/VehicleAssetForm.tsx`
- Typed client ve tipler:
  - `lib/api/stock-ecosystem.ts`
  - `lib/types/stock-ecosystem.ts`

## Seeds / RBAC / i18n
- `menu_seed.go`: `stock-vehicle-assets`
- `permission_seed.go`: `stock_vehicle_assets`
- `role_seed.go`: admin/editor yetki eşleşmesi
- `translation_seed.go`: `stock.vehicle_assets.*` + menu key
- `stock_seed.go`: örnek `stock_vehicle_assets` seed akışı

## Doğrulama
- `go test ./internal/modules/stock/... ./cmd/api` ✅
- `npx eslint "app/(dashboard)/admin/stock/vehicle-assets/**/*.tsx" "lib/api/stock-ecosystem.ts" "lib/types/stock-ecosystem.ts"` ✅

