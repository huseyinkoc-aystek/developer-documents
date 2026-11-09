# AYSTEK Taşımacılık, Geçiş, Çetele ve Veli/Personel Tahsilat Platformu — Enterprise Ana İş Planı

**Tarih:** 2026-05-09  
**Durum:** Başlangıç iş planı / devam dokümanı  
**Kısa kod:** `KT-ARACGPSCETELETAKIP-YONETIM-VE-VELI`  
**Hedef tenant:** `gursel-turizm`  
**Hedef siteler:** `aractakip.gurseltur.com.tr`, `veli.gurseltur.com.tr`, ileride kurum/üniversite public kanalları  
**Önerilen Core modül namespace’leri:** `transport_core`, `transport_access`, `transport_wallet`, `transport_portal`  

---

## 1. Amaç ve Karar

Bu proje iki büyük iş hattını tek AYSTEK Core stratejisi altında birleştirir:

1. **Kartlı/QR’lı taşımacılık, geçiş, çetele ve araç hareket doğrulama sistemi:** Üniversite, okul ve personel servislerinde Mifare kart, QR, misafir bilet, abonelik, indirim, bakiye, ring, güzergâh, offline geçiş ve GPS doğrulamalı çetele süreçleri.
2. **Veli/personel tahsilat ve self-service ödeme sistemi:** `veli.gurseltur.com.tr` üzerindeki kayıt, kişi bağlama, borç görüntüleme, sözleşme, kan grubu/dinamik alan, peşin/taksit/havale ödeme ve SAP dispatch süreçleri.

Bu iki hat ayrı uygulamalar olarak çoğaltılmayacak; **AYSTEK Core içinde tenant/site-aware bounded context’ler** olarak kurulacaktır.

Karar nettir:

- `aystek-core` tek backend otoritesi kalır.
- `aystek-core-next` operasyon/admin yönetim yüzeyi olur.
- `aractakip.gurseltur.com.tr` ve `veli.gurseltur.com.tr` site kanalları korunur; iş mantığı Core’a taşınır.
- Legacy Yii2/PHP sistem geçiş süresince read/bridge ve doğrulama kaynağı olarak kullanılır.
- Kart/QR geçiş, kullanıcı bazlı cüzdan, güzergâh/ücret, çetele, GPS doğrulama, ödeme, sözleşme, kişi bağlama, servis talebi, SAP dispatch, audit ve duyuru süreçleri Core içinde standartlaştırılır.

Bu plan, `/workspace/aystek-core-dev-docs/README_FOR_AI_ASSISTANTS.md` ve `/workspace/aystek-core-dev-docs/AYSTEK_CORE_CANONICAL_ARCHITECTURE.md` kurallarına bağlıdır. Bu projede her yeni modül için DDD katmanları, test, seed, permission, menu ve 28 dil translation standardı zorunludur.

---

## 2. İncelenen Mevcut Yapı

### 2.1 AYSTEK Core Backend

`/workspace/go/aystek-core` içinde bu projeyi doğrudan destekleyen mevcut yapı taşları vardır:

| Alan | Mevcut konum | Projede kullanım kararı |
|---|---|---|
| Tenant/site context | `internal/platform/middleware`, `contextkeys`, `tenantfilter`, `router` | Tüm yeni repository ve handler yollarında zorunlu scope kaynağı |
| Gürsel tenant seed | `internal/seeds/gursel_turizm_seed.go` | Mevcut `aractakip` site seed’i korunur; yeni `veli` site seed’i eklenir |
| GPS/çetele | `internal/modules/fleet_tally` | Araç/GPS çetele tarafı ayrı bounded context olarak kalır; ödeme portalına yalnızca event/API ile bağlanır |
| Ödeme | `internal/modules/payments`, `payments/intents`, `payments/gateways`, `payments/installments`, `payments/links`, `payments/risk` | Yeni tahsilat akışı bu modülleri reuse eder; borç bazlı quote/orchestration eksikleri yeni context’te tamamlanır |
| Sözleşme | `internal/modules/contracts`, `contracts/content`, `subscription_contracts` | Template/content altyapısı reuse edilir; borç bazlı requirement/signature katmanı eklenir |
| Muhasebe defteri | `internal/modules/ledger`, `billing`, `billing/invoices` | DebtAccount ve ödeme kapanışları ledger ile uyumlu event üretir |
| Entegrasyon | `internal/modules/integrations`, `erp_integrations`, RabbitMQ/eventbridge | SAP, legacy bridge, discount provider ve webhook süreçleri bu omurgaya bağlanır |
| SMS/OTP | `internal/modules/sms`, `notifications`, `push_tokens` | OTP, sözleşme kabul kanıtı, ödeme bildirimi ve operasyon duyuruları için reuse edilir |
| Menu/permission/i18n | `internal/seeds/menu_seed.go`, `permission_seed.go`, `translation_seed.go`, `localization_seed.go` | Her yeni ekran ve metin key’i için zorunlu güncelleme noktası |

### 2.2 AYSTEK Core Next Frontend

`/workspace/node/aystek-core-next` içinde benzer admin yüzeyleri mevcuttur:

- Ödeme ekranları: `app/(dashboard)/admin/payments`, `payment-intents`, `payment-gateways`, `payment-links`, `payment-risk`
- Sözleşme ekranları: `app/(dashboard)/admin/contracts`
- Entegrasyon ekranları: `app/(dashboard)/admin/integrations`, `erp-integrations`, `sms`
- Filo/çetele ekranları: `app/(dashboard)/admin/fleet-tally`, `telematics`
- API client stili: `lib/api/*.ts`
- Tip stili: `lib/types/*.ts`
- i18n kontrol komutu: `npm run check-i18n`

Yeni admin yüzeyleri bu pattern ile açılmalıdır:

- `app/(dashboard)/admin/transport-platform/...`
- `lib/api/transport-platform.ts`
- `lib/types/transport-platform.ts`
- `components/transport-platform/...`

### 2.3 Legacy Yii2/PHP Bulguları

İncelenen dosyalar:

- `/workspace/php/GurselVB22/models/ServisOnKayitForm.php` — Ön kayıt formu
- `/workspace/php/GurselVB22/models/ServistenAyrilmaTalebiForm.php` — Servisten ayrılma talebi
- `/workspace/php/GurselVB22/models/OgrenciSMSOnay.php` — SMS ile öğrenci girişi
- `/workspace/php/GurselVB22/models/SurucuHostesBilgileri.php` — Sürücü/hostes bilgileri
- `/workspace/php/GurselVB22/models/KardesIndirimi.php` — Kardeş indirimi (öğrenci TC listesi + dosya yükleme)
- `/workspace/php/GurselVB22/models/HavaleForm.php` — Havale bildirimi
- `/workspace/php/GurselVB22/models/VSozlesmesi.php` — Sözleşme değişken yer tutucuları ([VELI-TC] vb.)
- `/workspace/php/GurselVB22/controllers/SiteController.php` — Banka POS akışları (Garanti, İşbank, YapıKredi, Finans, Akbank)

Önemli geçiş sinyalleri:

- Ön kayıt formunda öğrenci, okul, sınıf, adres, doğum tarihi, kan grubu, anne/baba/diğer kişi bilgileri ve KVKK onayı toplanıyor.
- Servisten ayrılma talebi e-posta ve servis çağrısı mantığıyla ilerliyor.
- Öğrenci girişinde 5-6 haneli SMS onay kodu ve okul bağlamı kullanılıyor.
- Kan grubu için dış servis kontrolü var; bu veri KVKK açısından özel nitelikli veri sınıfına alınmalıdır.
- **Kardeş indirimi:** Birden fazla çocuğu aynı servise giden veliler için öğrenci TC listesi ve dosya kanıtıyla indirim uygulanıyor. `transport_discount_rules` modelinde `rule_type: family_group` desteği gerekir.
- **Multi-banka POS:** Her banka için ayrı action metodları var; yeni sistemde BIN/IIN tabanlı otomatik routing yapılacak. Bkz. ADR_003.
- **Sözleşme değişkenleri:** [VELI-TC], [OGR-TC] gibi yer tutucular sözleşme template'ine gömülüyor; yeni `ContractRequirement` modeli template variable snapshot taşımalı.
- Legacy sistem geçiş süresince doğrulama ve veri kaynağıdır; yeni domain modeli legacy field isimlerine göre şekillendirilmemeli.

---

## 3. Non-Negotiable Standartlar

Bu proje için Definition of Done aşağıdaki maddeleri otomatik içerir:

- Her backend context DDD katmanlarıyla yazılır: `model`, `dto`, `repository`, `service`, `handler`, `routes`.
- Her yeni collection/model için seed veya migration/seed kararı dokümante edilir.
- Her admin/public UI metni translation key olur ve `translation_seed.go` içinde 28 dil karşılığı eklenir.
- Her admin sayfası `menu_seed.go` içinde backend-driven menu entry alır.
- Her yeni module/action `permission_seed.go` içinde permission alır; Admin role seed davranışı korunur.
- Her repository sorgusu tenant/site scope taşır; `tenant_id` indeks dışı bırakılmaz.
- Public endpoint’ler tenant/site context resolver ve ModuleGate/AccessGate davranışını bypass etmez.
- Event-driven her özellik event key kataloğu ve payload sözleşmesi ile dokümante edilir.
- Payment card data AYSTEK tarafından saklanmaz; PAN/CVV Core içine girmez.
- Kan grubu gibi özel nitelikli veri maskelenir, ayrı erişim politikası ve audit izi taşır.
- Installer/deployment/runtime davranışı değişirse installer ve operasyon dokümanları aynı iş kaleminde güncellenir.

---

## 4. Önerilen Bounded Context Haritası

Ana backend namespace seti:

- `internal/modules/transport_core`: Kurum, taşınan kişi, grup, güzergâh, sefer, araç, hizmet alan/hizmet veren ilişki omurgası.
- `internal/modules/transport_access`: Kart/QR/Mifare geçiş, okutma, offline sync, refund, misafir bilet, abonelik ve geçiş kuralları.
- `internal/modules/transport_wallet`: Kullanıcı bazlı bakiye, yükleme, kredi limiti, TOSLA/POS yükleme, fatura sorumlusu ve hesap hareketleri.
- `internal/modules/transport_portal`: Veli/personel self-service borç, sözleşme, dinamik alan, havale, duyuru ve servis talebi portalı.

Bu namespace’ler mevcut `fleet_tally`, `telematics`, `stock/vehicle-assets`, `payments`, `contracts`, `ledger`, `integrations`, `sms` modüllerini tekrar yazmaz; onları reusable platform servisleri olarak kullanır.

| Context | Paket önerisi | Sorumluluk |
|---|---|---|
| Platform Bootstrap | `transport_core/bootstrap` veya seed-only | `gursel-turizm`, `aractakip`, `veli`, menü, permission, default policy |
| Institution Registry | `transport_core/institutions` | Üniversite, okul, firma, hizmet alan kurum, sezon, proje kodu |
| Transported Person Registry | `transport_core/persons` | Taşınan kişi, öğrenci/personel/öğretim görevlisi ayrımı, çoklu grup üyeliği |
| Group & Discount Policy | `transport_core/groups` | Bir kişinin birden fazla gruba dahil olması, kurum/grup bazlı indirim |
| Route & Fare Policy | `transport_core/routes` | Güzergâh, yön, gün/saat, tarih aralıklı ücret, çakışmasız ve boşluksuz fiyat geçmişi |
| Trip & Movement Ledger | `transport_core/trips` | İş atamalı veya iş atamasız sefer, ring, araç hareket raporu, çetele |
| Access Devices | `transport_access/devices` | Mifare/QR okuyucu, şoför/biletçi uygulaması, cihaz-araç eşleşmesi |
| Access Transactions | `transport_access/transactions` | Biniş/iniş, misafir geçişi, kart/QR okutma, vehicle context, offline queue |
| Ticketing & Subscriptions | `transport_access/tickets` | Tek kullanımlık QR, misafir bilet, abonelik, route/general validity |
| Wallet Accounts | `transport_wallet/accounts` | Kullanıcı bazlı bakiye, eksi bakiye/kredi limiti, kart değişimsiz hesap modeli |
| Top-ups & Refunds | `transport_wallet/topups` | Sanal POS/TOSLA yükleme, kurum max yükleme, 15 dk iade, fatura sorumlusu |
| Portal Identity Link | `transport_portal/guardians` | Portal kullanıcısı ile Person arasında ilişki, link request, doğrulama kanıtı |
| Dynamic Attributes | `transport_portal/attributes` | Kan grubu ve benzeri metadata tabanlı alanlar, requirement policy, masked access |
| Service Requests | `transport_portal/service_requests` | Ön kayıt, servisten ayrılma, iletişim/operasyon case workflow |
| Debt Ledger | `transport_portal/debts` | DebtAccount, DebtItem, collection window, discount application |
| Contract Policy | `transport_portal/contracting` | Borç bazlı contract requirement, template variable snapshot, signature evidence |
| Payment Orchestration | `transport_portal/payments` | Quote, BIN routing, payment intent binding, taksit kararları, havale intent |
| SAP Dispatch | `transport_core/sapdispatch` | Canonical event → SAP adapter dispatch state, retry, replay, idempotency. Hem `transport_core` hem `transport_portal` event'lerini tüketir; ortak adapter olarak `transport_core` altındadır. |
| Announcements | `transport_portal/announcements` | Public/portal duyuru, kurum/sezon/audience bazlı yayın |
| Operations | `transport_portal/ops` | Audit read model, reconciliation dashboard, event catalog, health |

Not: `transport_portal` sadece veli/personel ödeme portalı değildir. Yeni kullanıcı kapsamına göre daha üst çatı `transport_core + transport_access + transport_wallet + transport_portal` olmalıdır. Detaylı fark analizi ve hedef kapsam `SCOPE_RECONCILIATION_AND_TARGET_ARCHITECTURE.md` içindedir.

---

## 5. Hedef Veri Modeli

### 5.1 Koleksiyonlar ve Scope

> **Kapsam notu:** SAP dispatch koleksiyonu `transport_core/sapdispatch` altındadır. Hem `transport_core` hem `transport_portal` event'lerini tüketir; ortak bir dispatch adapter olarak konumlandırılmıştır. Gerekçe: ADR_001 Bölüm 5.

**transport_core domain koleksiyonları:**

| Collection | Zorunlu scope | Açıklama |
|---|---|---|
| `transport_institutions` | `tenant_id`, `site_id` | Okul, üniversite, firma master data |
| `transport_seasons` | `tenant_id`, `site_id`, `institution_id` | Eğitim dönemi / hizmet sezonu |
| `transport_persons` | `tenant_id`, `site_id`, `institution_id` | Taşınan kişi kaydı (öğrenci/personel) |
| `transport_groups` | `tenant_id`, `site_id`, `institution_id` | İndirim ve kural grubları |
| `transport_group_memberships` | `tenant_id`, `site_id`, `person_id`, `group_id` | Kişi-grup çoklu üyelik |
| `transport_routes` | `tenant_id`, `site_id`, `institution_id` | Güzergâh tanımı |
| `transport_route_directions` | `tenant_id`, `site_id`, `route_id` | Yönlü güzergâh (gidiş/dönüş ayrı) |
| `transport_fare_policies` | `tenant_id`, `site_id`, `route_id` | Tarih aralıklı, çakışmasız ücret tanımları |
| `transport_trips` | `tenant_id`, `site_id`, `vehicle_id` | Sefer/tur kaydı (iş atamalı veya atamasız) |
| `transport_trip_inferences` | `tenant_id`, `site_id`, `vehicle_id` | GPS + geçiş verisiyle çıkarılan sefer önerileri |
| `transport_movement_reports` | `tenant_id`, `site_id`, `trip_id` | Araç hareket raporu ve çetele |
| `transport_verified_work_units` | `tenant_id`, `site_id`, `trip_id` | Onaylı ve paylaşıma açık iş birimi |
| `transport_sap_dispatches` | `tenant_id`, `site_id`, `event_key` | SAP gönderim state, retry ve DLQ |

**transport_access domain koleksiyonları:**

| Collection | Zorunlu scope | Açıklama |
|---|---|---|
| `transport_access_devices` | `tenant_id`, `site_id`, `vehicle_id` | Mifare/QR okuyucu, şoför/biletçi cihazı |
| `transport_wallet_instruments` | `tenant_id`, `site_id`, `user_id` | Kullanıcıya bağlı Mifare kart, QR token, mobil token |
| `transport_access_transactions` | `tenant_id`, `site_id`, `device_id`, `trip_id` | Biniş/iniş/okutma/geçiş olayı |
| `transport_offline_batches` | `tenant_id`, `site_id`, `device_id` | Offline mod'da toplanan işlem grubu |
| `transport_guest_tickets` | `tenant_id`, `site_id`, `route_id` | Misafir tek kullanımlık QR bilet |
| `transport_subscription_passes` | `tenant_id`, `site_id`, `user_id` | Güzergâh veya genel abonelik |
| `transport_access_rules` | `tenant_id`, `site_id` | Günlük/haftalık tekrar, saat, gecikme kuralları |
| `transport_refund_windows` | `tenant_id`, `site_id` | İade penceresi politikası (örn. 15 dk) |

**transport_wallet domain koleksiyonları:**

| Collection | Zorunlu scope | Açıklama |
|---|---|---|
| `transport_wallet_accounts` | `tenant_id`, `site_id`, `user_id` | Kullanıcı bazlı bakiye hesabı |
| `transport_wallet_transactions` | `tenant_id`, `site_id`, `wallet_account_id` | Yükleme, harcama, iade, düzeltme hareketleri |
| `transport_topup_intents` | `tenant_id`, `site_id`, `user_id` | POS/TOSLA yükleme niyeti |
| `transport_topup_provider_bindings` | `tenant_id`, `site_id` | TOSLA, banka POS gibi yükleme sağlayıcı yapılandırması |
| `transport_credit_policies` | `tenant_id`, `site_id` | Grup/kurum bazlı eksi bakiye limiti |
| `transport_invoice_responsibilities` | `tenant_id`, `site_id` | Fatura sorumlusu (TOSLA, kullanıcı, kurum) tanımı |

**transport_portal domain koleksiyonları:**

| Collection | Zorunlu scope | Açıklama |
|---|---|---|
| `transport_guardian_links` | `tenant_id`, `site_id`, `person_id`, `user_id` | Kullanıcı-kişi ilişki kanıtı |
| `transport_link_requests` | `tenant_id`, `site_id` | Manuel inceleme kuyruğu |
| `transport_dynamic_field_definitions` | `tenant_id`, optional `site_id` | Alan metadata, KVKK sınıfı, mask policy |
| `transport_field_requirement_policies` | `tenant_id`, `site_id`, optional `institution_id` | Alan zorunluluk kuralları |
| `transport_person_attribute_values` | `tenant_id`, `site_id`, `person_id` | Kişi bazlı dinamik değerler (kan grubu vb.) |
| `transport_service_requests` | `tenant_id`, `site_id`, `institution_id` | Ön kayıt/ayrılma/workflow case |
| `transport_debt_accounts` | `tenant_id`, `site_id`, `institution_id`, `person_id`, `season_id` | Borç hesabı aggregate |
| `transport_debt_items` | `tenant_id`, `site_id`, `debt_account_id` | Borç satırları |
| `transport_discount_rules` | `tenant_id`, `site_id`, `institution_id` | Peşin/taksit/kurum/kişi/kardeş policy |
| `transport_collection_windows` | `tenant_id`, `site_id`, `institution_id`, `season_id` | Tahsilat açık/kapalı dönemleri |
| `transport_contract_requirements` | `tenant_id`, `site_id`, `debt_account_id` | Borç bazlı sözleşme zorunluluğu |
| `transport_contract_signatures` | `tenant_id`, `site_id`, `contract_requirement_id` | Kabul/OTP/SMS kanıtı; `otp_channel`, `otp_phone_hash`, `otp_sent_at` alanlarını içerir |
| `transport_payment_quotes` | `tenant_id`, `site_id`, `guardian_link_id` | Ön/final teklif ve taksit snapshot |
| `transport_bank_transfer_notices` | `tenant_id`, `site_id`, `debt_account_id` | Havale intent, dekont, onay/red |
| `transport_announcements` | `tenant_id`, `site_id` | Audience/institution/season scoped duyuru |

### 5.2 State Machine’ler

| Aggregate | Durumlar |
|---|---|
| `GuardianLink` | `requested → otp_pending → approved → rejected → revoked` |
| `ServiceRequest` | `draft → submitted → in_review → accepted → rejected → exported → closed` |
| `DebtAccount` | `draft → active → suspended → settled → exported → archived` |
| `ContractRequirement` | `required → presented → accepted → invalidated → superseded` |
| `PaymentQuote` | `provisional → final → expired → consumed → cancelled` |
| `PaymentIntent binding` | `quoted → pending_3ds → authorized → captured → failed → cancelled → reconciled` |
| `BankTransferNotice` | `created → submitted → under_review → approved → rejected → reconciled` |
| `SapDispatch` | `pending → processing → succeeded → failed → retry_scheduled → dead_lettered → replayed` |

Her state geçişi domain method ile doğrulanmalı; handler doğrudan status set etmemelidir.

---

## 6. API ve Route Tasarımı

### 6.1 Admin API

Önerilen prefix: `/admin/transport-platform`

| Endpoint grubu | Örnek endpoint’ler |
|---|---|
| Dashboard | `GET /dashboard`, `GET /health`, `GET /event-catalog` |
| Institutions | `GET/POST /institutions`, `GET/PUT/DELETE /institutions/:id`, `POST /institutions/:id/status` |
| Seasons | `GET/POST /institutions/:id/seasons`, `PUT /seasons/:id` |
| Persons | `GET /persons`, `GET /persons/:id`, `POST /persons/import`, `PATCH /persons/:id/status` |
| Link Requests | `GET /link-requests`, `POST /link-requests/:id/approve`, `POST /link-requests/:id/reject` |
| Dynamic Fields | `GET/POST /field-definitions`, `GET/POST /field-policies`, `GET /person-attributes/audit` |
| Service Requests | `GET /service-requests`, `GET /service-requests/:id`, `POST /service-requests/:id/transition` |
| Debt | `GET /debt-accounts`, `POST /debt-accounts/import`, `GET /debt-accounts/:id`, `POST /debt-accounts/:id/reprice` |
| Contracts | `GET /contract-requirements`, `POST /contract-policies/recompute`, `GET /contract-signatures` |
| Payments | `GET /payment-quotes`, `GET /bank-transfer-notices`, `POST /bank-transfer-notices/:id/approve` |
| SAP | `GET /sap-dispatches`, `POST /sap-dispatches/:id/retry`, `POST /sap-dispatches/:id/replay` |
| Announcements | `GET/POST /announcements`, `PUT /announcements/:id`, `POST /announcements/:id/publish` |

### 6.2 Public/Self-Service API

Önerilen prefix: `/public/transport-platform`

> **Not:** Site context (veli vs. aractakip) host/domain resolver middleware'inden okunur, URL prefix'ten değil. Bu AYSTEK Core'un mevcut site context middleware modeliyle uyumludur. Request body'de `site_code` varsa yalnızca cross-check amacıyla kullanılır, scope otoritesi değildir.

| Akış | Endpoint’ler |
|---|---|
| Public bootstrap | `GET /bootstrap`, `GET /institutions`, `GET /announcements/public` |
| Portal hesap | Existing auth reuse; gerektiğinde `POST /portal/register`, `POST /portal/verify` |
| Kişi bağlama | `POST /guardian/link-requests`, `POST /guardian/link-requests/:id/otp/send`, `POST /guardian/link-requests/:id/otp/verify` |
| Profil tamamlama | `GET /persons/:id/required-fields`, `PUT /persons/:id/attributes` |
| Borç görüntüleme | `GET /me/debt-accounts`, `GET /debt-accounts/:id` |
| Sözleşme | `GET /debt-accounts/:id/contracts/required`, `POST /contracts/:requirement_id/accept` |
| Ödeme quote | `POST /payment-quotes`, `POST /payment-quotes/:id/finalize`, `POST /payment-intents/:id/confirm` |
| Havale | `POST /bank-transfer-notices`, `GET /bank-transfer-notices/:id` |
| Servis talebi | `POST /service-requests/pre-registration`, `POST /service-requests/leave-service`, `GET /me/service-requests` |

Public endpoint’lerde raw `tenant_id` güvenilmez input olarak alınmamalı; host/domain/site resolver ile bağlam kurulmalıdır. Eğer request body’de `site_code` taşınırsa yalnızca cross-check amacıyla kullanılmalıdır.

---

## 7. Admin UI İş Planı

Yeni admin menü kökü:

- Key: `transport-platform`
- TitleKey: `menu.operations.transport_platform`
- Path: `/admin/transport-platform/dashboard`
- Icon: `bus-front` veya `school`

Alt sayfalar:

| Menü | Path | Permission |
|---|---|---|
| Dashboard | `/admin/transport-platform/dashboard` | `transport_core:read` |
| Kurumlar | `/admin/transport-platform/institutions` | `transport_institutions:read` |
| Taşınan Kişiler | `/admin/transport-platform/persons` | `transport_persons:read` |
| Gruplar ve İndirimler | `/admin/transport-platform/groups` | `transport_groups:read` |
| Güzergâh ve Ücretler | `/admin/transport-platform/routes` | `transport_routes:read` |
| Sefer ve Çetele | `/admin/transport-platform/trips` | `transport_trips:read` |
| Geçiş Cihazları | `/admin/transport-platform/access/devices` | `transport_access:read` |
| Geçiş İşlemleri | `/admin/transport-platform/access/transactions` | `transport_access:read` |
| Bilet ve Abonelikler | `/admin/transport-platform/access/tickets` | `transport_access:read` |
| Cüzdan ve Yüklemeler | `/admin/transport-platform/wallet` | `transport_wallet:read` |
| Bağlama Talepleri | `/admin/transport-platform/portal/link-requests` | `transport_guardians:review` |
| Dinamik Alanlar | `/admin/transport-platform/portal/dynamic-fields` | `transport_attributes:read` |
| Servis Talepleri | `/admin/transport-platform/portal/service-requests` | `transport_service_requests:read` |
| Borç Hesapları | `/admin/transport-platform/portal/debts` | `transport_debts:read` |
| Sözleşme Gereksinimleri | `/admin/transport-platform/portal/contracts` | `transport_contracting:read` |
| Ödeme Orkestrasyonu | `/admin/transport-platform/portal/payments` | `transport_payments:read` |
| Havale Operasyonları | `/admin/transport-platform/portal/bank-transfers` | `transport_bank_transfers:read` |
| SAP Dispatch | `/admin/transport-platform/sap-dispatches` | `transport_sap_dispatch:read` |
| Duyurular | `/admin/transport-platform/announcements` | `transport_announcements:read` |
| Audit | `/admin/transport-platform/audit` | `transport_ops:read` |

Frontend standardı:

- Hard-coded user-visible text yok.
- Tüm UI key’leri backend seed ile uyumlu.
- API response envelope normalize edilir.
- Loading, empty, error ve permission-denied state’leri ayrı gösterilir.
- Hassas alanlar varsayılan maskeli render edilir; unmask action audit gerektirir.

---

## 8. Event Kataloğu

| Event key | Üretici | Tüketici | İdempotency key |
|---|---|---|---|
| `transport.guardian_link.requested` | guardians service | notifications, ops | `link_request_id` |
| `transport.guardian_link.approved` | guardians service | SAP/CRM optional, audit | `guardian_link_id` |
| `transport.person.attribute.completed` | attributes service | contract policy, audit | `person_id + field_key + version` |
| `transport.service_request.created` | service_requests | SAP/ops/integrations | `service_request_id` |
| `transport.debt_account.activated` | debts service/import | contracting, notifications | `debt_account_id + version` |
| `transport.debt_item.added` | debts service/import | contracting, quote invalidation | `debt_item_id` |
| `transport.contract_requirement.created` | contracting | notifications | `contract_requirement_id` |
| `transport.contract.accepted` | contracting | SAP/archive/audit | `contract_signature_id` |
| `transport.payment_quote.created` | payments orchestration | audit | `payment_quote_id` |
| `transport.payment.intent.bound` | payments orchestration | payments module | `transport_payment_binding_id` |
| `transport.payment.captured` | payments callback/orchestration | ledger, SAP, notifications | `payment_attempt_id` |
| `transport.bank_transfer.submitted` | bank transfers | ops | `bank_transfer_notice_id` |
| `transport.bank_transfer.approved` | bank transfers | ledger, SAP, notifications | `bank_transfer_notice_id + approved_at` |
| `transport.sap_dispatch.failed` | sapdispatch | ops/DLQ | `sap_dispatch_id + attempt` |

Event payload’ları `tenant_id`, `site_id`, `correlation_id`, `actor`, `subject`, `event_version` ve `occurred_at` alanlarını zorunlu taşımalıdır.

---

## 9. Fazlara Bölünmüş Teslimat Planı

### Faz 0 — Keşif ve Mimari Kilitleme

**Hedef:** Belirsizlikleri kapatmak ve ilk implementation slice’ı risksiz başlatmak.

İşler:

- Legacy DB/schema ve Yii2 controller/service uçlarının envanteri.
- Banka/POS sağlayıcıları, merchant-of-record ve 3DS modeli netleştirme.
- SAP sürümü ve entegrasyon yolu: OData, IDoc, BAPI/RFC, PI/PO veya Integration Suite.
- Kurum/sezon/sözleşme/indirim matrisi çıkarma.
- KVKK/hukuk kararı: kan grubu ve özel nitelikli alanların hukuki dayanağı.
- İlk veri sözlüğü ve migration mapping dokümanı.

Çıktılar:

- `DISCOVERY_FINDINGS.md`
- `LEGACY_MAPPING_MATRIX.md`
- `BANK_AND_POS_DECISION_RECORD.md`
- `SAP_ADAPTER_DECISION_RECORD.md`
- `KVKK_SPECIAL_CATEGORY_DATA_POLICY.md`

Kalite kapısı:

- Belirsizlikler “blocker / assumption / accepted risk” olarak sınıflandırılmış olmalı.

### Faz 1 — Core Foundation ve Site Bootstrap

**Hedef:** `transport_core`, `transport_access`, `transport_wallet` ve `transport_portal` suite temelini Core içinde görünür, seed’li ve yetkili hale getirmek.

Backend:

- `transport_core`, `transport_access`, `transport_wallet`, `transport_portal` root doc ve module constants.
- `gursel-turizm` altında `veli` site seed’i.
- Permission seed: `transport_*` module/action setleri.
- Menu seed: admin navigation.
- Translation seed: menü, statü, hata ve temel ekran key’leri 28 dil.
- Route skeleton ve health/dashboard endpoint’i.
- Index seed veya `CreateIndexes` setleri.

Frontend:

- `/admin/transport-platform/dashboard`
- `lib/api/transport-platform.ts`
- `lib/types/transport-platform.ts`
- Dashboard shell, empty state, permission state.

Test:

- Backend route registration/unit tests.
- Seed idempotency tests.
- Frontend `npm run check-i18n`.

### Faz 2 — Institution, Person, Guardian Link

**Hedef:** Kullanıcı hesabı ile taşınan kişi arasında güvenli ilişki kurmak.

Backend:

- Institution/season CRUD.
- Person registry ve institution membership.
- Guardian link request, OTP send/verify, manual approve/reject.
- Legacy person lookup adapter interface.
- Audit events ve notifications.

Frontend:

- Kurum listesi/detail.
- Kişi listesi/detail.
- Bağlama talepleri inceleme kuyruğu.

Public:

- Institution list.
- Link request create.
- OTP verification.

Test:

- TCKN/passport validation unit tests.
- Link state machine tests.
- Rate limit ve brute-force OTP tests.

### Faz 3 — Dynamic Attributes ve Service Requests

**Hedef:** Kan grubu ve form alanlarını statik formdan policy-driven metadata yapısına taşımak.

Backend:

- Field definition CRUD.
- Requirement policy engine.
- Person attribute value write/read with masking.
- Service request aggregate: pre-registration, leave-service.
- Legacy form alan mapping’i.
- KVKK consent/evidence linkleri.

Frontend:

- Dynamic field definitions paneli.
- Service request list/detail/transition paneli.
- Masked sensitive fields and unmask audit action.

Public:

- Required fields endpoint.
- Ön kayıt ve servisten ayrılma formları.

Test:

- Required field resolution tests.
- Sensitive data masking tests.
- Service request transition tests.

### Faz 4 — Debt Ledger ve Contract Policy

**Hedef:** Borç hesabı ve borç bazlı sözleşme gereksinimini domain içinde kurmak.

Backend:

- DebtAccount/DebtItem model, import ve reprice.
- CollectionWindow ve DiscountRule.
- ContractRequirement ve ContractSignature.
- Existing `contracts/content` reuse adapter.
- Template variable catalog ve rendered snapshot/hash.
- Events: debt activated, debt item added, contract accepted.

Frontend:

- Borç hesapları paneli.
- Sözleşme gereksinimi ve kabul kanıtı paneli.
- Discount/collection window yönetimi.

Public:

- Borç görüntüleme.
- Sözleşme gösterim ve kabul akışı.

Test:

- Debt state machine tests.
- Contract policy recompute tests.
- Template variable snapshot tests.

### Faz 5 — Payment Orchestration ve Bank Transfer

**Hedef:** Borç bazlı kart/havale tahsilatını iki aşamalı quote + final payment modeliyle kurmak.

Backend:

- PaymentQuote aggregate.
- Provisional installment offers.
- BIN/IIN metadata provider interface.
- Acquirer route policy.
- Existing `payments/intents` binding.
- 3DS return/callback orchestration.
- BankTransferNotice workflow.
- Payment captured/reconciled events.

Frontend:

- Ödeme orkestrasyonu paneli.
- Havale onay/red paneli.
- Quote ve attempt detay ekranı.

Public:

- Quote create.
- BIN ile final quote.
- Payment confirm/return.
- Bank transfer notice submit.

Test:

- Route priority tests.
- Installment eligibility tests.
- 3DS success/fail/timeout contract tests.
- Bank transfer reconciliation tests.

### Faz 6 — SAP Dispatch, Legacy Bridge ve Reconciliation

**Hedef:** Core canonical events üzerinden SAP ve legacy geçiş omurgasını çalıştırmak.

Backend:

- SAP adapter interface: OData/IDoc/BAPI profile placeholders.
- `transport_sap_dispatches` state machine.
- Retry, DLQ, replay, payload diff.
- Legacy read bridge ve dual-run comparison.
- Ledger reconciliation reports.

Frontend:

- SAP dispatch list/detail/replay.
- Reconciliation dashboard.
- Bridge health paneli.

Test:

- Contract tests for SAP payload mapping.
- Idempotent dispatch tests.
- DLQ/replay tests.

### Faz 7 — Hardening, UAT ve Pilot

**Hedef:** Pilot kurumlarla kontrollü canlıya geçiş; güvenlik, performans ve operasyonel olgunluk.

**Backend:**

- Security review: authz boundary, rate limiting, CSP/payment page isolation, sensitive data access controls.
- KVKK uyumluluk son kontrol: kan grubu ve özel nitelikli veri maskeleme, consent log completeness.
- Performans optimizasyon: borç listesi, quote hesaplama, sözleşme render, admin raporları — index ve aggregation planları gözden geçirilmeli.
- Stress test: eş zamanlı geçiş okutması (offline batch sync + online), eş zamanlı payment intent.
- Circuit breaker / fallback: SAP dispatch hata oranı arttığında sistem nasıl tepki verir?
- Reconciliation doğrulaması: SAP gönderilen ve muhasebede kapanan kayıtlar eşleşmeli.

**Frontend:**

- Accessibility kontrol: WCAG 2.1 AA temel.
- Responsive design: mobil veli portalı.
- i18n coverage:  tüm yeni keyler için sıfır eksik.
- Empty / error / permission-denied state tüm ekranlarda mevcut.

**Operasyon:**

- Monitoring dashboards: SAP dispatch sağlığı, offline batch queue derinliği, payment quote conversion, guardian link onay bekleyenler.
- SLO tanımları:  → SAP dispatch  max gecikme.
- Runbook'lar: SAP down senaryosu, offline device stuck, duplicate access transaction, debt reconciliation gap.
- Rollback planı: her faz için bağımsız rollback yeteneği.

**UAT:**

- Pilot kurum: en az bir okul + bir üniversite bağlamında test.
- Test senaryoları: ön kayıt, bağlama, kan grubu, borç görüntüleme, sözleşme, peşin ödeme, taksitli ödeme, havale bildirimi, SAP dispatch, admin rapor, kardeş indirimi.
- Misafir QR bilet satın alma + tek kullanımlık okutma.
- Offline geçiş + batch sync + raporlara yansıma.

**Kalite Kapısı:**

- Pilot kurum için uçtan uca  akışı kanıtlanmış olmalı.
- Hiçbir açık P0/P1 güvenlik bulgusu bulunmamalı.
- i18n coverage eksiksiz.
- Tüm kritik paths için runbook mevcut.

---

## 10. İlk Sprint Önerisi

İlk kod sprint’i Faz 1 ile sınırlı tutulmalıdır. Amaç hızlıca ekrana “bir şeyler koymak” değil, doğru zemini atmaktır.

Sprint 1 scope:

- `transport_core` route skeleton ve `transport_access`/`transport_wallet`/`transport_portal` ADR sınırları.
- `gursel-turizm` için `veli` site seed.
- Permission seed ve menu seed.
- 28 dil translation seed minimum key seti.
- Dashboard health endpoint.
- Next.js dashboard shell.
- Seed idempotency ve route tests.

Sprint 1 out-of-scope:

- Gerçek ödeme almak.
- Gerçek SAP dispatch.
- Legacy veri yazmak.
- Kan grubu gibi hassas veriyi toplamak.

Sprint 1 DoD:

- `go test` ilgili package’larda geçer.
- `go build ./cmd/api` geçer.
- `npm run check-i18n` geçer.
- `npm run build` en azından yeni route kaynaklı hata üretmez.
- Menüler backend seed’den gelir.
- Yeni doküman veya checkpoint `plans/aystek-veli-tahsilat-servis-platformu` altında güncellenir.

---

## 11. Test ve Validation Matrisi

| Katman | Komut / yöntem | Ne yakalar |
|---|---|---|
| Go unit | `go test ./internal/modules/transport_core/... ./internal/modules/transport_access/... ./internal/modules/transport_wallet/... ./internal/modules/transport_portal/...` | Domain state, policy, repository mocks |
| Go integration targeted | `go test ./internal/platform/... ./internal/seeds/...` | Tenant/site context, seed, permission/menu |
| Go build | `go build ./cmd/api` | Wiring/import/compile |
| Frontend i18n | `npm run check-i18n` | UI key eksikleri |
| Frontend build | `npm run build` | Route/component/type hataları |
| API contract | Postman/OpenAPI veya Go handler tests | Envelope, pagination, errors |
| Security | Rate limit/authz/tenant isolation tests | Public abuse ve tenant leak |
| E2E UAT | Browser/manual veya agent-browser | Kullanıcı akışı |

---

## 12. Operasyon ve Kurulum Etkileri

Bu proje ilk fazlarda installer davranışını değiştirmeyebilir. Ancak aşağıdaki değişikliklerden biri yapılırsa installer ve operasyon dokümanları aynı iş kaleminde güncellenmelidir:

- Yeni zorunlu environment variable.
- Yeni worker/binary/cron.
- Yeni RabbitMQ queue/exchange contract.
- Yeni seed/bootstrap sırası.
- Yeni release artifact.
- Yeni public domain/site onboarding prosedürü.
- Yeni backup/restore veya migration adımı.

Güncellenecek minimum dosyalar:

- `/workspace/go/aystek-core-installer/docs/KULLANIM_KILAVUZU.md`
- `/workspace/go/aystek-core-installer/README.md`
- `/workspace/aystek-core-dev-docs/user-guide/05-operations/SUNUCU_KURULUM_REHBERI.md`
- `/workspace/aystek-core-dev-docs/user-guide/05-operations/OPERATIONS_GUIDE.md`
- `/workspace/aystek-core-dev-docs/user-guide/05-operations/DEPLOYMENT_SCENARIOS.md`

---

## 13. Riskler ve Mitigasyon

| Risk | Etki | Mitigasyon |
|---|---|---|
| Legacy veri kalitesi düşük | Kişi/borç eşleşmeleri hatalı olur | Data profiling, confidence score, manual review queue |
| SAP yöntemi geç netleşir | Dispatch gecikir | Canonical event + adapter-first tasarım |
| POS sağlayıcısı değişir | Payment kodu kırılır | Provider adapter + hosted payment page + route policy |
| Kan grubu KVKK riski | Hukuki ve reputasyonel risk | Minimum veri, açık dayanak, maskeleme, audit, özel permission |
| Tenant/site leak | Kritik güvenlik açığı | Repository scope tests, tenantfilter checks, context integrity |
| UI/backend drift | Kullanıcı erişemeyen ekranlar | Menu seed + permission seed + i18n DoD |
| Chargeback ispatı zayıf | Finansal kayıp | Contract snapshot, payment metadata, device/IP hash, audit chain |
| Rollout kesinti yaratır | Operasyon durur | Pilot kurum, dual run, rollback, feature flags |

---

## 14. Devamlılık Kuralları

Bu klasör proje boyunca canlı tutulacaktır.

Her sprint sonunda:

- Bu `MASTER_WORKPLAN.md` güncellenir veya yeni checkpoint dokümanı eklenir.
- Tamamlanan işler `done/` altındaki task summary standardına göre ayrıca kapanır.
- Yeni event key, route, menu, permission ve translation key listesi dokümante edilir.
- Açık riskler ve belirsizlikler güncellenir.
- Bir sonraki sprint scope’u net yazılır.

Önerilen checkpoint adı:

`TRANSPORT_PORTAL_SPRINT_{N}_{KISA_KONU}_CHECKPOINT_YYYY-MM-DD.md`

---

## 15. Hemen Sonraki Adım

Önerilen ilk uygulama işi:

1. `transport_core` foundation aç ve `transport_portal` ile sınırlarını netleştir.
2. `gursel-turizm` için `aractakip` mevcut seed’iyle çakışmadan `veli.gurseltur.com.tr` site seed’i ekle.
3. `transport_core`, `transport_access`, `transport_wallet`, `transport_portal` için permission/menu/translation seed temelini tamamla.
4. Admin dashboard shell, health endpoint’i ve scope selector modelini oluştur.
5. Route/seed/i18n/build doğrulamasını çalıştır.

Bu adım ödeme veya hassas veri toplamadan ilerlediği için düşük risklidir; ama tüm sonraki işlerin platform standardına doğru oturmasını sağlar.
