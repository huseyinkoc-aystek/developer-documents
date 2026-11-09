# AYSTEK CORE — STOCK Ekosistemi Enterprise Delivery Planı

## Kapsam ve adlandırma kararı
- İş adı: MOBİNAV Stok/Araç takip.
- Teknik modül adı: `stock` (MOBİNAV adı modül kökünde kullanılmayacak).
- Backend kök: `go/aystek-core/internal/modules/stock/*`
- Admin kök: `node/aystek-core-next/app/(dashboard)/admin/stock/*`

Faz 8 kapanış notu:
- `stock_integrations` artık stratejik geliştirme yüzeyi değil, `compatibility-read-only` katmanıdır.
- Final migration/retirement kararı için bkz. [FAZ8_CLOSURE_AND_RETIREMENT_DECISION.md](/workspace/aystek-core-dev-docs/plans/FAZ8_CLOSURE_AND_RETIREMENT_DECISION.md)

---

## 1) Domain decomposition

### Modül listesi ve sorumluluk sınırları
- `stock_catalog`: brand/model/supplier katalogları.
- `stock_inventory`: `product_item` yaşam döngüsü (IN_STOCK/ASSIGNED/SALE_PENDING/SOLD/RETURNED/PASSIVE).
- `stock_vehicles`: araç + plaka geçmişi + proje/güzergah referansları.
- `stock_vehicle_assets`: araç-ürün montaj/söküm ilişkisi.
- `stock_assignments`: zimmet (yalnız yetkili stok kullanıcılarına).
- `stock_sales`: satış talebi/onay/state machine.
- `stock_returns`: kısmi iade + otomatik iade zimmeti.
- `stock_subscriptions`: abonelik + fiyat geçmişi.
- `stock_integrations`: SAP + telematik konnektör/job/retry/idempotency.
- `stock_locations`: anlık/çoklu/tarihsel konum servisleri (connector üstünden).
- `stock_sms_commands`: template + command execution + callback.
- `stock_reporting`: export jobs + rapor read model.

### Existing modüller ile entegrasyon
- Kimlik/izin: `auth`, `users`, `rbac`.
- Operasyonel çapraz kesit: `audit`, `outbox/events`, `queue/jobs`, `rate_limit`, `cache`.
- Finansal genişleme: `billing/payments/ledger` (bridge yaklaşımıyla, bypass yok).
- Ortak UI/i18n: `translation_seed`, `menu_seed`, `permission_seed`, `role_seed`.

---

## 2) Veri modeli (Mongo)

### Koleksiyonlar
- `stock_brands`, `stock_models`, `stock_suppliers`
- `stock_product_items`
- `stock_vehicles`, `stock_vehicle_plate_history`
- `stock_vehicle_assets`
- `stock_assignments`
- `stock_sale_orders`, `stock_sale_order_items`
- `stock_return_orders`, `stock_return_order_items`
- `stock_subscriptions`, `stock_subscription_pricing`
- `stock_integration_connectors`, `stock_integration_jobs`
- `stock_vehicle_locations_cache`
- `stock_sms_command_templates`, `stock_sms_command_executions`
- `stock_report_jobs`

### Zorunlu index pattern
- Tüm koleksiyonlar: `{tenant_id:1, site_id:1, deleted_at:1}`
- Liste index’i: `{tenant_id:1, site_id:1, status:1, updated_at:-1, deleted_at:1}`

### Kritik unique constraints
- `stock_models`: `{tenant_id,site_id,brand_id,name,product_type,deleted_at:null}` unique partial
- `stock_product_items`:
  - `{tenant_id,site_id,serial_no,deleted_at:null}` unique partial
  - `{tenant_id,site_id,imei,deleted_at:null}` unique partial
- `stock_vehicles`: `{tenant_id,site_id,current_plate,deleted_at:null}` unique partial
- `stock_assignments`: `{tenant_id,site_id,product_item_id,status:'ACTIVE',deleted_at:null}` unique partial
- `stock_sale_order_items`: aynı ürün için tek aktif satış talebi (unique partial)
- `stock_integration_jobs`: `{tenant_id,site_id,idempotency_key}` unique
- `stock_sms_command_executions`: `{tenant_id,site_id,provider_message_id}` unique sparse
- `stock_vehicle_locations_cache`: `{tenant_id,site_id,vehicle_id}` unique

### TTL (öneri)
- `stock_integration_jobs`: `completed_at + 180 gün`
- `stock_sms_command_executions`: `last_status_at + 365 gün`
- `stock_report_jobs`: `completed_at + 90 gün` (dosya retention ile uyumlu)

---

## 3) API sözleşmesi

### Admin endpointleri (`/admin/stock/*`)
- Catalog:
  - `GET/POST /admin/stock/catalog/brands`
  - `GET/PUT /admin/stock/catalog/brands/:id`
  - `GET/POST /admin/stock/catalog/models`
  - `GET/PUT /admin/stock/catalog/models/:id`
- Inventory:
  - `GET/POST /admin/stock/inventory/items`
  - `GET/PUT/DELETE /admin/stock/inventory/items/:id`
  - `POST /admin/stock/inventory/import/csv`
  - `POST /admin/stock/inventory/bulk/sale-pending`
  - `POST /admin/stock/inventory/bulk/freeze`
- Vehicles/Assets:
  - `GET/POST /admin/stock/vehicles`
  - `GET/PUT /admin/stock/vehicles/:id`
  - `GET/POST /admin/stock/vehicles/:id/assets`
  - `POST /admin/stock/vehicles/:id/plates`
- Assignments:
  - `GET/POST /admin/stock/assignments`
  - `POST /admin/stock/assignments/:id/release`
- Sales/Returns:
  - `GET/POST /admin/stock/sales`
  - `POST /admin/stock/sales/:id/approve`
  - `POST /admin/stock/sales/:id/reject`
  - `GET/POST /admin/stock/returns`
- Subscriptions:
  - `GET/POST /admin/stock/subscriptions`
  - `GET/PUT /admin/stock/subscriptions/:id`
  - `POST /admin/stock/subscriptions/:id/pricing`
- Integrations/Locations:
  - `GET/POST /admin/stock/integrations/connectors`
  - `GET/POST /admin/stock/integrations/jobs`
  - `GET /admin/stock/locations/current`
  - `GET /admin/stock/locations/history`
- SMS:
  - `GET/POST /admin/stock/sms/templates`
  - `POST /admin/stock/sms/commands/send`
  - `POST /admin/stock/sms/callback`
- Reporting:
  - `GET /admin/stock/reports/subscriptions-active`
  - `GET /admin/stock/reports/vehicle-assets`
  - `POST /admin/stock/reports/exports/jobs`
  - `GET /admin/stock/reports/exports/jobs`

### Portal endpointleri
- `GET/POST /portal/stock/sales`
- `POST /portal/stock/field-reports/assets`
- `GET /portal/stock/locations/current`

### DTO standardı
- Request: `CreateXRequest`, `UpdateXRequest`, `ListXQuery`, `BulkXRequest`
- Response: `XResponse`, `ListXResponse(items,page,page_size,total)`
- Hata sözleşmesi: `success/error.code/error.message_key/message`
- Pagination: `page`, `page_size`, `total`, `total_pages`

### Örnek hata kodları
- `TENANT_OR_SITE_MISSING`
- `PERMISSION_DENIED`
- `VALIDATION_FAILED`
- `STATE_TRANSITION_INVALID`
- `CONFLICT_DUPLICATE_IDENTIFIER`
- `RESOURCE_NOT_FOUND`
- `INTEGRATION_UNAVAILABLE`
- `IDEMPOTENCY_CONFLICT`
- `RATE_LIMITED`
- `CSV_INVALID_SCHEMA`
- `CSV_INJECTION_BLOCKED`

---

## 4) Menu + Permissions

### Menu entries (`menu_seed.go`)
- Parent: `menu.operations.stock`
  - `path`: `/admin/stock/inventory/items`
  - `icon`: `boxes`
  - `order`: `340`

Child:
- `/admin/stock/inventory/items`
- `/admin/stock/vehicles`
- `/admin/stock/assignments`
- `/admin/stock/sales`
- `/admin/stock/returns`
- `/admin/stock/subscriptions`
- `/admin/stock/locations/current`
- `/admin/stock/sms/templates`
- `/admin/stock/integrations/jobs`
- `/admin/stock/reports/subscriptions-active`

### Permission list (`permission_seed.go`)
- `stock_catalog.read|create|update|delete`
- `stock_inventory.read|create|update|delete|import|bulk_action`
- `stock_vehicles.read|create|update|delete`
- `stock_vehicle_assets.read|create|remove`
- `stock_assignments.read|create|release`
- `stock_sales.read|create|approve|reject|sync`
- `stock_returns.read|create|approve|sync`
- `stock_subscriptions.read|create|update|pricing_manage`
- `stock_integrations.read|manage|run`
- `stock_locations.read`
- `stock_sms.read|template_manage|send|callback`
- `stock_reports.read|export`
- `stock_audit.read`

Role:
- Admin: full auto-assign.
- Ops editor: create/update + limited approve.
- Field user: portal-scope only.

---

## 5) Backlog (Epic > Feature > User Story)

### Epic-1: Catalog + Inventory Foundation
- Feature: Catalog yönetimi
  - Story: Brand/model CRUD
  - AC: model ürün tipiyle uyumlu, unique kuralı, soft delete
  - Test: duplicate model, tenant/site scope
  - DoD: Go modül + tests + seeds + Next list/create/edit
- Feature: Product item yaşam döngüsü
  - Story: state transition yönetimi
  - AC: invalid transition 409
  - Test: status matrix + race guard
  - DoD: outbox + audit + i18n

### Epic-2: Vehicles + Assets + Assignments
- Feature: Araç/plaka geçmişi
  - Story: plaka değişim geçmişi
  - AC: aktif plaka tekil
  - Test: overlap validation
  - DoD: API + admin UI + seed
- Feature: Zimmet
  - Story: sadece stok rolüne zimmet
  - AC: rol dışı hedefe atama yasak
  - Test: permission/object scope
  - DoD: audit event + integration test

### Epic-3: Sales + Returns
- Feature: Satış onay akışı
  - Story: saha talep, merkez onay
  - AC: aynı ürün tek aktif satış talebi
  - Test: concurrency/idempotency
  - DoD: SAP job enqueue + izlenebilirlik
- Feature: Kısmi iade
  - Story: item bazlı iade + otomatik iade zimmeti
  - AC: araç üstü kayıt güncellenir
  - Test: partial return consistency
  - DoD: outbox + audit + UI

### Epic-4: Subscriptions + Pricing History
- Story: geçerlilik tarihli fiyat
- AC: tarih aralığı çakışması engeli
- Test: overlap + asOf fiyat çözümleme
- DoD: CRUD + reporting readiness

### Epic-5: Integrations + Locations + SMS
- Story: SAP job orchestration
- Story: telematik current/history fetch
- Story: SMS template/send/callback
- AC: retry/idempotency/signature
- Test: callback replay + rate-limit
- DoD: monitoring + dead-letter görünümü

### Epic-6: Reporting + Export Center
- Story: aktif abonelikli hatlar raporu
- Story: araç üzeri ürün raporu
- AC: CSV injection korunumu
- Test: export sanitizer + job lifecycle
- DoD: export jobs UI + download guard

---

## 6) Sprint planı (2 hafta)

### Sıralı plan
- Sprint 1: Epic-1
- Sprint 2: Epic-2
- Sprint 3: Epic-3 (sales)
- Sprint 4: Epic-3 (returns) + Epic-4
- Sprint 5: Epic-5 (integrations core)
- Sprint 6: Epic-5 (locations + sms)
- Sprint 7: Epic-6 + parity hardening

### Bağımlılık grafiği (özet)
- Catalog -> Inventory -> Assignments -> Sales -> Returns
- Vehicles -> VehicleAssets -> Locations/SMS target resolution
- Integrations -> Sales/Returns sync + Locations connectors
- Reporting tüm domain read-model’lerine bağımlı

### Riskler / mitigasyon
- SAP teknik kanal belirsizliği: connector abstraction + contract test
- Legacy PHP gizli iş kuralı: parity matrix + shadow run
- Hacimli CSV: async import job + chunking + retry
- Authorization riski: object-level policy + negatif test paketi

---

## 7) Observability & Security

### Audit eventleri
- `stock.product_item.status_changed`
- `stock.assignment.created|released`
- `stock.sale.created|approved|rejected|synced`
- `stock.return.created|processed|synced`
- `stock.subscription.pricing_changed`
- `stock.integration.job.created|retry|failed|succeeded`
- `stock.sms.command.sent|status_updated`
- `stock.export.job.created|completed|downloaded`

### Security kontrol listesi
- OWASP API authorization (BOLA/BFLA) zorunlu test.
- CSV import: MIME/type/size kontrolü + schema validation.
- CSV export: formula injection sanitize.
- Callback/webhook: HMAC signature + replay protection.
- Secrets: encrypted-at-rest, masked response.
- Rate limit: import, run job, sms send, bulk endpoints.

---

## 8) AYSTEK CORE Standards Compliance Checklist

| Modül | model | dto | repo(interface+mongo) | service | service_test | routes/handler | translation(28) | menu | permission | role auto-assign | domain_seed | Next typed api/types/pages |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| stock_catalog | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ |
| stock_inventory | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ |
| stock_vehicles | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ |
| stock_vehicle_assets | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ |
| stock_assignments | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ |
| stock_sales | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ |
| stock_returns | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ |
| stock_subscriptions | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ |
| stock_integrations | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ |
| stock_locations | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ |
| stock_sms_commands | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ |
| stock_reporting | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ | ☑ |

> Bu tablo sprint bazında güncellenerek release gate olarak kullanılacaktır.

---

## 9) Adım adım başlangıç planı (Execution Kickoff)

### Sprint-0 (Hazırlık ve tasarım sabitleme, 1 hafta)
- Hedef: Sprint-1 implementasyonuna girilmeden kontrat ve standartların netleşmesi.
- Çıktılar:
  - `stock` modül kök dizin şablonu ve naming convention dokümanı.
  - API contract v0 (OpenAPI taslağı): catalog + inventory minimum endpointleri.
  - Mongo index planı ve unique/partial stratejisi.
  - Seed key matrisi (menu + permission + translation key listesi).
  - Test stratejisi dokümanı (unit/integration/contract).
- Gate (Sprint-0 Done):
  - Mimari onay: DDD sınırları ve existing modül entegrasyonları net.
  - Güvenlik onayı: authorization modeli + CSV güvenlik kuralları net.
  - Operasyon onayı: audit/outbox event isimleri dondurulmuş.

### Sprint-1 (Catalog + Inventory Foundation, 2 hafta)
- Hedef: `stock_catalog` ve `stock_inventory` modüllerini production-ready çıkarmak.
- Backend kapsamı:
  - `stock_catalog`: brand/model/supplier CRUD + validation + tests.
  - `stock_inventory`: `product_item` CRUD + durum geçişleri + soft delete + tests.
  - Route registration + `cmd/api/main.go` wiring.
  - Index registration + seed setleri.
- Frontend kapsamı:
  - `/admin/stock/inventory/items` list/create/edit.
  - `/admin/stock/catalog/brands` ve `/admin/stock/catalog/models` list/create/edit.
  - `lib/api/stock.ts` + `lib/types/stock.ts`.
- Gate (Sprint-1 Done):
  - `go test ./internal/modules/stock/... ./cmd/api` yeşil.
  - Next lint/typecheck ilgili sayfalarda yeşil.
  - i18n/menu/permission/role/domain seed zinciri tamam.

### Sprint-2 (Vehicles + Assets + Assignments, 2 hafta)
- Hedef: araç varlık ilişkisi ve zimmet akışını kapatmak.
- Çıktılar:
  - `stock_vehicles`, `stock_vehicle_assets`, `stock_assignments`.
  - Object-level authorization testleri.
  - Admin list/detail flow + timeline/audit görünümü.

### Sprint-3 (Sales + SAP job orchestration çekirdeği, 2 hafta)
- Hedef: satış talep/onay ve SAP sync job pipeline’ı.
- Çıktılar:
  - `stock_sales` state machine.
  - `stock_integrations` içinde SAP job enqueue/retry/idempotency çekirdeği.
  - Admin onay ekranı + saha portal talep ekranı (MVP).

### Sprint-4 (Returns + Subscriptions, 2 hafta)
- Hedef: kısmi iade, otomatik iade zimmeti, abonelik fiyat geçmişi.
- Çıktılar:
  - `stock_returns`, `stock_subscriptions`.
  - tarih aralığı çakışma engeli testleri.

### Sprint-5 (Locations + SMS + Reporting başlangıcı, 2 hafta)
- Hedef: operasyonel rekabetçi katmanlar.
- Çıktılar:
  - `stock_locations` current/history.
  - `stock_sms_commands` template/send/callback.
  - `stock_reporting` export job v1.

---

## 10) Delivery governance

### Her sprint sonunda zorunlu rapor seti
- `TASK_STOCK_S{N}.md`: teknik detay, dosya listesi, endpoint listesi.
- `TASK_STOCK_S{N}_SUMMARY.md`: yönetici özeti, risk, bir sonraki sprint planı.

### Change control
- Breaking contract değişiklikleri: sadece sprint başında, migration planıyla.
- Event adı/değer değişikliği: versiyonlanmadan yapılmaz.
- Seed key silme: yasak (deprecated işaretlenir).

### Release quality gate (her sprint)
- Functional: AC pass oranı >= %95
- Test: kritik servislerde coverage trend artışta
- Security: authz negatif test paketi çalışmış
- Observability: yeni endpointler metric + audit üretmiş
