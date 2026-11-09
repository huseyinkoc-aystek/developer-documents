# TASK_043 — STOCK Sprint-2 (Assignments)

## Kapsam
- `internal/modules/stock/assignments` modülü eklendi.
- Inventory köprüsü ile assignment oluşturma/bırakma akışı bağlandı.
- Admin UI: `/admin/stock/assignments` + `/admin/stock/assignments/new`.

## Backend
- Model: `stock_assignments` (active/released)
- API:
  - `GET /admin/stock/assignments`
  - `POST /admin/stock/assignments`
  - `GET /admin/stock/assignments/:id`
  - `POST /admin/stock/assignments/:id/release`
- Kurallar:
  - Aynı `product_item_id` için tek aktif assignment
  - Create sırasında inventory item statüsü kontrolü
  - Release sırasında assignment state transition + inventory `in_stock` dönüşü
- Indexler:
  - Unique partial: tenant+site+product_item_id+status(active)
  - List index: tenant+site+status+updated_at+deleted_at

## Seed/RBAC/i18n
- `menu_seed.go`: `stock-assignments` menü girişi
- `permission_seed.go`: `stock_assignments` modülü + `release` action
- `role_seed.go`: editor/admin mapping’e `stock_assignments` eklendi
- `translation_seed.go`: `menu.commerce.inventory.stock_assignments` + `stock.assignments.*`

## Next Admin
- `app/(dashboard)/admin/stock/assignments/page.tsx`
- `app/(dashboard)/admin/stock/assignments/new/page.tsx`
- `app/(dashboard)/admin/stock/assignments/AssignmentForm.tsx`
- API/Types:
  - `lib/api/stock-ecosystem.ts`
  - `lib/types/stock-ecosystem.ts`

## Doğrulama
- `go test ./internal/modules/stock/assignments/... ./internal/modules/stock/vehicles/... ./internal/modules/stock/catalog/... ./internal/modules/stock/inventory/... ./internal/seeds ./cmd/api` ✅
- `npx eslint app/(dashboard)/admin/stock/assignments/**/*.tsx lib/api/stock-ecosystem.ts lib/types/stock-ecosystem.ts` ✅
