# STOCK Subscription Pricing Redesign Plan

Bu plan, abonelik ücretini abonelikten bağımsızlaştırıp stok altında merkezi ve geçerlilik aralıklı bir pricing modeli ile yönetmek için hazırlanmıştır.

## A) Domain Tasarımı
- Yeni bounded model: `inventory_subscription_pricings` (merkezi ücret tablosu).
- Abonelik snapshot alanları:
  - `current_price`, `currency`
  - `pricing_id`, `pricing_version`
  - `pricing_valid_from`, `pricing_valid_to`
- Geçerlilik kuralı:
  - `valid_from` inclusive
  - `valid_to` exclusive
  - `valid_to = null` açık uç
- Timezone:
  - DB’de UTC saklanır
  - UI girişleri local, API’de UTC’ye normalize edilir.
- Overlap kuralı:
  - Aynı `tenant_id + site_id + plan_code + currency` için aralıklar çakışamaz.

## B) Mongo Veri Modeli + İndeksler
- Koleksiyon: `inventory_subscription_pricings`
- Önerilen alanlar:
  - `tenant_id`, `site_id`, `plan_code`, `currency`, `amount`
  - `valid_from`, `valid_to`, `status`, `version`
  - audit + soft delete alanları
- İndeksler:
  - `tenant_id, site_id, plan_code, currency, valid_from, deleted_at`
  - `tenant_id, site_id, plan_code, currency, valid_to, deleted_at`
  - unique: `tenant_id, site_id, plan_code, currency, version`
- Overlap engelleme:
  - DB index + service-level overlap query birlikte kullanılır.

## C) Servis Akışları
- Create pricing:
  - required alan doğrulama
  - overlap kontrolü
  - `version = next(version)` ataması
- Create subscription:
  - `plan_code` (subscription code) + `currency` + `started_at` ile resolve
  - pricing bulunduysa snapshot alanları yazılır
  - bulunmazsa `stock.subscriptions.validation.pricing_required_for_date`
- Update pricing:
  - immutable history yaklaşımı
  - mevcut kaydı archive et, yeni versiyon kayıt aç.

## D) API Sözleşmesi (Admin + Domain)
- Admin pricing:
  - `GET /admin/stock/subscription-pricings`
  - `POST /admin/stock/subscription-pricings`
  - `GET /admin/stock/subscription-pricings/:id`
  - `PUT /admin/stock/subscription-pricings/:id`
  - `DELETE /admin/stock/subscription-pricings/:id` (soft)
  - `POST /admin/stock/subscription-pricings/resolve`
- Abonelik create değişikliği:
  - manuel fiyat zorunlu değil, pricing snapshot service içinde çözülür.
- Hata kodları:
  - `stock.subscription_pricing.validation.overlap`
  - `stock.subscription_pricing.error.not_found_for_date`
  - `stock.subscriptions.validation.pricing_required_for_date`

## E) Admin UI (Next.js)
- Yeni ekran:
  - `/admin/stock/subscription-pricings` (list/filter)
  - `/admin/stock/subscription-pricings/new` (create)
  - `/admin/stock/subscription-pricings/[id]` (edit/version)
- Abonelik form UX:
  - manuel fiyat girişi yerine “merkezi fiyat otomatik çözülür” bilgilendirmesi
  - mevcut snapshot fiyatı read-only gösterim
  - pricing yoksa API hatası kullanıcıya net yansıtılır.

## F) Seed / RBAC / i18n
- Permission:
  - `stock_subscription_pricing.read/create/update/delete`
- Menu:
  - `/admin/stock/subscription-pricings` (inventory/stock parent altında)
- Translation:
  - `stock.subscription_pricing.*`
  - `stock.subscriptions.validation.pricing_required_for_date`
  - `stock.subscriptions.fields.price_resolve_hint`
- 28 dil stratejisi:
  - `allLangs(tr,en)` ile placeholder/fallback standardı.

## G) Migration Planı
- Aşama-0 Pre-migration:
  - yeni koleksiyon + indeks + permission/menu/i18n deploy
- Aşama-1 Backfill:
  - mevcut subscription pricing verilerini merkezi tabloya tenant/site bazlı taşı
  - idempotent job + dry-run raporu
- Aşama-2 Dual-read:
  - feature flag ile subscription create önce merkezi tablodan resolve eder
- Aşama-3 Cutover:
  - eski per-subscription pricing yazımı kapatılır (read-only)
- Aşama-4 Cleanup:
  - legacy endpoint/ekranlar kaldırılır veya deprecated edilir
- Rollback:
  - feature flag kapanınca legacy read yolu aktif kalır (veri kaybı olmadan).

## H) Test Planı
- Unit:
  - overlap reject
  - resolve at date (inclusive/exclusive sınırlar)
  - missing pricing -> expected error
- Integration:
  - migration idempotency
  - tenant/site izolasyon doğrulaması
  - concurrency (aynı aralığa eşzamanlı insert)
- Edge:
  - `valid_to = null`
  - `start_date == valid_to` (kapsam dışı)
  - farklı currency senaryoları
  - timezone local->UTC dönüşüm doğruluğu.
