# STOCK Sprint-1 Execution Plan

## Kapsam
- Sprint: S1 (2 hafta)
- Amaç: `stock_catalog` + `stock_inventory` foundation teslimatı.
- Boundary: sadece `internal/modules/stock/*` ve ilgili seed + Next admin katmanı.

## Sprint hedefleri
- `stock_catalog` modülü:
  - brand/model/supplier CRUD
  - product_type uyumlu model validasyonu
  - tenant/site scoped repository
- `stock_inventory` modülü:
  - `product_item` CRUD
  - durum geçişleri (`IN_STOCK -> ASSIGNED -> SALE_PENDING -> SOLD -> RETURNED -> PASSIVE`)
  - soft delete + audit + outbox
- Admin UI:
  - `/admin/stock/catalog/brands`
  - `/admin/stock/catalog/models`
  - `/admin/stock/inventory/items`

## Story kırılımı (Epic > Feature > Story)

### Epic S1-E1 Catalog
- Feature S1-E1-F1: Brand CRUD
  - Story S1-E1-F1-US1: Brand create/update/list/delete
  - AC:
    - `name` zorunlu
    - aynı tenant+site içinde aktif duplicate yok
    - soft delete uygulanır
  - Test:
    - duplicate create -> 409
    - tenant izolasyon
- Feature S1-E1-F2: Model CRUD
  - Story S1-E1-F2-US1: Model create/update/list/delete
  - AC:
    - `brand_id`, `name`, `product_type` zorunlu
    - unique partial index ile duplicate engeli
  - Test:
    - invalid product_type -> 400
    - duplicate model -> 409

### Epic S1-E2 Inventory
- Feature S1-E2-F1: Product Item CRUD
  - Story S1-E2-F1-US1: Item create/list/detail/update/delete
  - AC:
    - `product_type`, `source_type` zorunlu
    - imei/serial unique kuralları koşullu
    - hard delete yok
  - Test:
    - imei duplicate
    - serial duplicate
- Feature S1-E2-F2: Status transition
  - Story S1-E2-F2-US1: Durum geçiş endpointi
  - AC:
    - geçersiz geçişte 409
    - audit + outbox event üretilir
  - Test:
    - tüm geçerli/geçersiz transition matrisi

## Teknik teslimatlar

### Backend dosya yapısı (her modül)
- `model/*.go`
- `dto/*.go`
- `repository/interface.go`
- `repository/mongo.go`
- `service/service.go`
- `service/service_test.go`
- `handler/admin.go`
- `routes.go`

### Seeds
- `translation_seed.go`: `stock.*` key seti (28 dil)
- `menu_seed.go`: stock parent + S1 child sayfaları
- `permission_seed.go`: `stock_catalog.*`, `stock_inventory.*`
- `role_seed.go`: admin auto-assign
- `stock_seed.go`: brand/model örnek seed

### Next.js
- `lib/types/stock.ts`
- `lib/api/stock.ts`
- `app/(dashboard)/admin/stock/catalog/brands/page.tsx`
- `app/(dashboard)/admin/stock/catalog/models/page.tsx`
- `app/(dashboard)/admin/stock/inventory/items/page.tsx`

## API sözleşmesi (S1 kapsamı)
- `GET/POST /admin/stock/catalog/brands`
- `GET/PUT/DELETE /admin/stock/catalog/brands/:id`
- `GET/POST /admin/stock/catalog/models`
- `GET/PUT/DELETE /admin/stock/catalog/models/:id`
- `GET/POST /admin/stock/inventory/items`
- `GET/PUT/DELETE /admin/stock/inventory/items/:id`
- `POST /admin/stock/inventory/items/:id/status`

## Kabul kriterleri (Sprint-1 Exit)
- Backend:
  - `go test ./internal/modules/stock/... ./cmd/api` ✅
  - tüm endpointler tenant+site scoped ✅
- Frontend:
  - ilgili sayfalarda lint/typecheck ✅
  - hardcoded string yok, i18n key kullanımı ✅
- Governance:
  - seed zinciri tam (menu/permission/role/translation/domain) ✅
  - audit/outbox event üretimi doğrulanmış ✅

## Riskler ve mitigasyon (S1)
- Risk: eski stok verisinde imei/serial kirli olabilir.
  - Mitigasyon: rollout öncesi duplicate raporu + migration script.
- Risk: permission matrix eksikliği UI aksiyonlarını kilitleyebilir.
  - Mitigasyon: permission smoke test checklist.
- Risk: model/product_type uyumsuzluğu.
  - Mitigasyon: service-level strict validation + seed örnekleri.
