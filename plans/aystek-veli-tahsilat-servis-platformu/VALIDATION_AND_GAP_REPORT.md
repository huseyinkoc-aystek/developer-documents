# Doğrulama ve Boşluk Raporu — Transport Platform Plan

**Tarih:** 2026-05-11  
**Hazırlayan:** GitHub Copilot (Claude Sonnet 4.6)  
**Kapsam:** Codex tarafından üretilen 5 dokümanın (INDEX, MASTER_WORKPLAN, SCOPE_RECONCILIATION, ADR_001, SPRINT_01) doğrulanması; eksik, hatalı veya tamamlanması gereken kısımların tespiti.  
**Referans kaynak:** `/workspace/aystek-core-dev-docs/README_FOR_AI_ASSISTANTS.md`, legacy PHP projeleri (`/workspace/php/GurselVB22`, `/workspace/php/AYShuttleGursel`), mevcut AYSTEK Core modülleri.

---

## 1. Genel Değerlendirme

Codex çıktıları; proje kapsamını, bounded context sınırlarını ve fazlı yaklaşımı genel olarak doğru kurgulamıştır. Temel mimari kararlar (4 namespace, mevcut modüllerin reuse edilmesi, multi-tenant/site scope, state machine'ler, 28 dil zorunluluğu) AYSTEK Core standartlarıyla uyumludur. Bununla birlikte; legacy sistemden çıkarılan belirli iş kuralları, kritik bir altyapı ADR'ı (offline sync), tüm domain'lerin koleksiyon tablosu ve Faz 0 uygulama planı eksiktir.

**Genel not:** ✅ Onaylandı, ⚠️ Eksik/güncellenmeli, 🔴 Kritik boşluk, ❌ Hatalı/yanıltıcı

---

## 2. Onaylanan Kararlar

| Karar | Doğrulama kaynağı |
|---|---|
| 4 namespace: `transport_core`, `transport_access`, `transport_wallet`, `transport_portal` | ✅ ADR_001; mevcut modül yapısıyla çakışma yok |
| `fleet_tally`, `telematics`, `payments`, `contracts`, `ledger`, `integrations`, `sms` modüllerinin reuse edilmesi | ✅ `/workspace/go/aystek-core/internal/modules` dizini doğruladı |
| `gursel-turizm` tenant zaten var, `aractakip` site seed mevcut | ✅ `gursel_turizm_seed.go` içinde `aractakip.gurseltur.com.tr` seed'i mevcut |
| `veli.gurseltur.com.tr` için yeni site seed eklenmesi | ✅ Mevcut seedde yoktu, doğru tespit |
| State machine modelleri (GuardianLink, ServiceRequest, DebtAccount, ContractRequirement, PaymentQuote, BankTransferNotice, SapDispatch) | ✅ AYSTEK Core DDD standardıyla uyumlu |
| Event kataloğu; `tenant_id`, `site_id`, `correlation_id`, `event_version` zorunlu alanları | ✅ Mevcut integrations/eventbridge yapısıyla uyumlu |
| 28 dil zorunluluğu ve menu_seed kuralları | ✅ README_FOR_AI_ASSISTANTS standartlarıyla uyumlu |
| Banka kartı PAN/CVV AYSTEK'e girmemeli | ✅ Güvenlik standardı ve mevcut payments modülü tasarımıyla uyumlu |
| Kan grubu özel nitelikli veri (KVKK) — maskeleme zorunluluğu | ✅ Doğru tespit; legacy PHP `ServisOnKayitForm.php` içinde kan grubu toplanıyor |
| Public endpoint'lerde `tenant_id` raw input değil, host/domain resolver'dan alınmalı | ✅ Mevcut `middleware` ve site context yapısıyla uyumlu |
| SAP dispatch için idempotent retry + DLQ | ✅ Mevcut `integrations/eventbridge` yapısıyla tutarlı |
| `transport_portal/payments` mevcut `payments/intents` modülünü bağlar, kopyalamaz | ✅ Doğru tasarım |
| BIN/IIN routing ile banka POS seçimi | ✅ Legacy `SiteController.php`'de her banka için ayrı action; bu yaklaşımı BIN tabanlı otomatik routinge taşımak doğru karar |
| Offline-first zorunluluğu | ✅ Kullanıcı anlatımından doğrulandı; ancak ADR henüz yok (bkz. Bölüm 3) |

---

## 3. Kritik Boşluklar

### 3.1 🔴 MASTER_WORKPLAN Koleksiyon Tablosu Eksik

**Sorun:** Bölüm 5.1 yalnızca `transport_portal` domain'ine ait 19 koleksiyonu listeler. `transport_core`, `transport_access` ve `transport_wallet` domain'lerinin koleksiyonları hiç yok.

**Etki:** Sprint planlaması ve MongoDB index tasarımı bu tabloya dayanır. Eksik koleksiyonlar sprint'lerde tutarsız isimlerle önerilir; çakışma riski artar.

**Düzeltme:** MASTER_WORKPLAN bölüm 5.1 güncellendi. Tüm eksik koleksiyonlar eklendi. (bkz. bu raporla birlikte yapılan MASTER_WORKPLAN güncellemesi)

**Eksik koleksiyonlar:**

| Collection | Domain | Zorunlu scope |
|---|---|---|
| `transport_institutions` | transport_core | `tenant_id`, `site_id` |
| `transport_seasons` | transport_core | `tenant_id`, `site_id`, `institution_id` |
| `transport_persons` | transport_core | `tenant_id`, `site_id`, `institution_id` |
| `transport_groups` | transport_core | `tenant_id`, `site_id`, `institution_id` |
| `transport_group_memberships` | transport_core | `tenant_id`, `site_id`, `person_id`, `group_id` |
| `transport_fare_policies` | transport_core | `tenant_id`, `site_id`, `route_id` |
| `transport_routes` | transport_core | `tenant_id`, `site_id`, `institution_id` |
| `transport_route_directions` | transport_core | `tenant_id`, `site_id`, `route_id` |
| `transport_trips` | transport_core | `tenant_id`, `site_id`, `vehicle_id` |
| `transport_trip_inferences` | transport_core | `tenant_id`, `site_id`, `vehicle_id` |
| `transport_movement_reports` | transport_core | `tenant_id`, `site_id`, `trip_id` |
| `transport_verified_work_units` | transport_core | `tenant_id`, `site_id`, `trip_id` |
| `transport_access_devices` | transport_access | `tenant_id`, `site_id`, `vehicle_id` |
| `transport_wallet_instruments` | transport_access | `tenant_id`, `site_id`, `user_id` |
| `transport_access_transactions` | transport_access | `tenant_id`, `site_id`, `device_id`, `trip_id` |
| `transport_offline_batches` | transport_access | `tenant_id`, `site_id`, `device_id` |
| `transport_guest_tickets` | transport_access | `tenant_id`, `site_id`, `route_id` |
| `transport_subscription_passes` | transport_access | `tenant_id`, `site_id`, `user_id` |
| `transport_access_rules` | transport_access | `tenant_id`, `site_id` |
| `transport_refund_windows` | transport_access | `tenant_id`, `site_id` |
| `transport_wallet_accounts` | transport_wallet | `tenant_id`, `site_id`, `user_id` |
| `transport_wallet_transactions` | transport_wallet | `tenant_id`, `site_id`, `wallet_account_id` |
| `transport_topup_intents` | transport_wallet | `tenant_id`, `site_id`, `user_id` |
| `transport_topup_provider_bindings` | transport_wallet | `tenant_id`, `site_id` |
| `transport_credit_policies` | transport_wallet | `tenant_id`, `site_id` |
| `transport_invoice_responsibilities` | transport_wallet | `tenant_id`, `site_id` |

### 3.2 🔴 Offline Sync ADR Yok

**Sorun:** Offline-first zorunluluğu hem kullanıcı anlatımında hem de SCOPE_RECONCILIATION dokümanında "en kritik gereksinim" olarak tanımlanmış. Ancak bunun için mimarı seviyede bir karar dokümanı (ADR) yok. Offline transaction log formatı, cihaz-sunucu senkronizasyon protokolü, çakışma çözümü politikası ve offline risk limitleri tanımlanmamış.

**Etki:** Bu ADR olmadan `transport_access/offline_sync` implementasyonu için sprint planlaması yapılamaz. Cihaz firmware/uygulama ekibi ile sözleşme kurulamaz.

**Düzeltme:** `ADR_002_OFFLINE_FIRST_SYNC.md` oluşturuldu.

### 3.3 🔴 Faz 0 İçin Uygulama Planı Yok

**Sorun:** SPRINT_01 gibi ayrıntılı bir execution plan, Faz 0 (Keşif ve Mimari Kilitleme) için yok. MASTER_WORKPLAN sadece ne yapılacağını listeler, nasıl yapılacağını söylemez. Özellikle:
- Legacy DB şemasına nasıl erişileceği (canlı DB mi, dump mı, Yii2 model analizi mi)
- SAP entegrasyon türünün nasıl saptanacağı
- Banka POS sağlayıcılarının nasıl derleneceği
- KVKK belgelendirmesinin nasıl yapılanacağı

**Düzeltme:** `FAZ_0_DISCOVERY_EXECUTION_PLAN.md` oluşturuldu.

### 3.4 🔴 Çoklu Banka POS Routing ADR Yok

**Sorun:** Legacy `SiteController.php`, her banka (İşbank, Garanti, YapıKredi, Finans, Akbank) için ayrı action metodları içeriyor. Yeni sistemde BIN/IIN tabanlı otomatik routing önerilmekte; ancak:
- BIN veritabanı kaynağı belirlenmemiş
- Banka-merchant mutabakat profili (MID, TID) nasıl yönetilecek
- Kartın bankası bilinmiyorsa fallback POS hangisi
- Taksit tablolarının gerçek zamanlı mı yoksa önceden cache'lenmiş mi alınacağı

belirsiz.

**Düzeltme:** `ADR_003_MULTI_POS_BANK_ROUTING.md` oluşturuldu.

---

## 4. Önemli Boşluklar (Faz 2+ Öncesi Tamamlanmalı)

### 4.1 ⚠️ Kardeş İndirimi (Sibling Discount) Eksik

**Kaynak:** `/workspace/php/GurselVB22/models/KardesIndirimi.php`

**Durum:** Legacy sistemde birden fazla çocuğu aynı okul servisine giden veli için kardeş indirimi uygulanıyor. Bu indirim için öğrenci TC listesi ve dosya yükleme gerekiyor. MASTER_WORKPLAN'daki `transport_discount_rules` koleksiyonu genel indirim politikalarını kapsıyor; ancak kardeş indirimi özel koşullar içeriyor:
- Aynı `guardian_link`'e bağlı birden fazla kişi aynı rotada mı?
- Belge kanıtı gerekiyor mu? (legacy'de `file` alanı var)
- İndirim yüzdesi mi, sabit tutar mı, taksit ücreti mi?

**Aksiyon:** Faz 2 (Institution ve Person Registry) scope'una "sibling/family discount policy" kullanım senaryosu eklenmeli. `transport_discount_rules` modelinde `rule_type: family_group` desteklenmeli. Bu Faz 0 keşif soruları arasına alınmalı.

### 4.2 ⚠️ TOSLA / Üçüncü Taraf Top-up Fatura Sorumluluğu Yeterince Detaylı Değil

**Kaynak:** Kullanıcı anlatımı: "TOSLA gibi sistemlerden yapılan yüklemeler de entegre edilmelidir. Yüklemenin kime faturalandırılacağı (TOSLA mı, kullanıcı mı) bilgisi sisteme girilmelidir."

**Durum:** `InvoiceResponsibility` aggregate MASTER_WORKPLAN'da var. Ancak iki senaryonun nasıl ayrıştırılacağı belirsiz:
- TOSLA kullanıcıya yükleme yaptığında fatura TOSLA'ya mı kesilir?
- Eğer fatura TOSLA'ya kesilecekse bu hangi SAP transaction tipiyle gidecek?
- TOSLA'nın kendi API'si varsa entegrasyon sözleşmesi (contract test) lazım.

**Aksiyon:** TOSLA entegrasyon senaryosu Faz 0 keşif soruları arasına alınmalı. `transport_wallet/topups` modeli için `TopUpProviderBinding` içinde `invoice_to: "provider" | "user"` alanı tanımlanmalı.

### 4.3 ⚠️ Servis Yaklaşma Bildirimi Sprint'te Yok

**Kaynak:** Kullanıcı anlatımı: "servis kendisine yaklaştığında bildirim gönderilmesi"

**Durum:** Bu özellik push notification ve GPS konum entegrasyonundan tetiklenir. Mevcut Core'da `telematics` ve `notifications/push_tokens` modülleri var. Ancak yaklaşma bildiriminin:
- Tetikleme mantığı (geofence yarıçapı? tahmini varış süresi?)
- Bildirim kanalı (push, SMS, her ikisi?)
- Opt-in mekanizması
- Spam koruma (tekrar tekrar tetiklenmemesi için debounce)

henüz tanımlanmamış ve hiçbir sprint'te yer almıyor.

**Aksiyon:** `transport_core/trips` domain'i içinde `TripWaypointAlert` event'i tanımlanmalı. Faz 7 veya Faz 6 içine bir alt başlık olarak eklenebilir.

### 4.4 ⚠️ Araç Sahibi Onayıyla Veri Paylaşımı ADR Yok

**Kaynak:** Kullanıcı anlatımı: "Araç sahipleri, hizmet verdikleri firmalara yaptıkları işin doğrulanmış verisini mikro servisler üzerinden paylaşabilir. Bu veri araç sahibi onayı ile gönderilmelidir."

**Durum:** SCOPE_RECONCILIATION dokümanda bahsedilmiş ("API sharing + approval + data product yaklaşımı"). Ancak mevcut `api_sharing` modülü bunu kısmen destekleyebilir; tam olarak destekleyip desteklemediği net değil.

**Aksiyon:** Bu konu `api_sharing` modülünün tenant-scoped "verified work data product" paylaşım kapasitesiyle örtüşüp örtüşmediği Faz 0'da araştırılmalı. Eğer `api_sharing` yeterliyse, yalnızca `transport_core/movement_reports` üzerinde approval workflow eklenmesi gerekir.

### 4.5 ⚠️ SMS ile Sözleşme İmzası Yeterince Yapılanmış Değil

**Kaynak:** Kullanıcı anlatımı: "sms olursa daha da gerçek kişi imzalamış diyebiliriz"

**Durum:** `ContractSignature` aggregate mevcut ve OTP kanıtı taşıyor. Ancak:
- Sözleşme kabul OTP'si ile kimlik doğrulama OTP'si ayrı mı olmalı?
- OTP SMS kanalı `sms` modülü üzerinden mi yoksa notifications modülü üzerinden mi gönderilmeli?
- Hukuki geçerlilik için SMS kanıtının hangi formatta saklanacağı (timestamp, hash, GSM numarası) belirsiz.

**Aksiyon:** `ContractSignature` modeline `otp_channel`, `otp_phone_hash` ve `otp_sent_at` alanları eklenmeli. Bu Faz 4 (Contract Policy) kapsamına alınabilir.

### 4.6 ⚠️ Tahsilat Dönem Yönetimi (Collection Window) Admin UI'da Detaylı Değil

**Kaynak:** Kullanıcı anlatımı: "Bazen o okuldan ya firmadan ödeme almayacağımız dönemler olabilir. Bu durumu yönetimden açıp kapatabiliriz."

**Durum:** `transport_collection_windows` koleksiyonu ve `CollectionWindow` aggregate tanımlanmış. Ancak hangi sprinte ait olduğu, admin panel deseni ve "toggle" işleminin bir event mi yoksa doğrudan status update mi olacağı belirsiz.

**Aksiyon:** `transport_portal/debts` veya yeni `transport_portal/collection_ops` subpackage'inde admin endpoint olarak modellenmeli. Faz 4 scopuna eklenebilir. `CollectionWindow` state machine: `draft → open → closed → archived`.

### 4.7 ⚠️ Mobil ve Rol-bazlı Uygulama Kapsamı Tanımlanmamış

**Kaynak:** Kullanıcı anlatımı: "araç sahibi (bir firma ya da şahıs olabilir), aracı şoför kullanıyorsa şoför ekranı, hizmet verilen firma ekranı ve personel ya da öğrenci ekranları bulunmalıdır"

**Durum:** Codex dokümanları yalnızca `aystek-core` (backend) ve `aystek-core-next` (admin panel) kapsamlıdır. Şoför ekranı, taşınan kişi uygulaması, araç sahibi ekranı gibi mobil/public kanallar için scope tanımlanmamış.

**Aksiyon:** Bu kanallar `aysteknoloji-com` (public site) veya yeni Flutter/React Native uygulama ile mi yönetilecek? Faz 0 keşif soruları arasına alınmalı. Her kanalın API endpoint erişim modeli belgelenmeli.

---

## 5. Küçük Tutarsızlıklar

### 5.1 ⚠️ SAP Dispatch Namespace Belirsizliği

**Sorun:** MASTER_WORKPLAN Bölüm 4'te "SAP Dispatch | `transport_core/sapdispatch` veya `transport_portal/sapdispatch`" şeklinde belirsiz bırakılmış.

**Karar:** SAP dispatch, `transport_core/sapdispatch` altında kurulmalıdır. Gerekçe:
- SAP dispatch hem `transport_portal` (ödeme, borç) hem `transport_core` (sefer, çetele, doğrulanmış iş) event'lerini iletir.
- `transport_portal` içine alınırsa `transport_core` event'leri için `portal` paketine bağımlılık oluşur — bu tersine bir bağımlılıktır.
- `transport_core/sapdispatch` tüm alt domain'lerden bağımsız, canonical event tabanlı dispatch adapter olarak konumlandırılmalıdır.

**Aksiyon:** MASTER_WORKPLAN güncellendi, ADR_001 güncellendi.

### 5.2 ⚠️ MASTER_WORKPLAN Faz 7 İçeriği Eksik

**Sorun:** Faz 7 başlığı açılmış ancak içerik yazılmamış (doküman kesilmiş).

**Aksiyon:** MASTER_WORKPLAN güncellendi; Faz 7 detayları eklendi.

### 5.3 ⚠️ Route Prefix Tutarsızlığı

**Sorun:** SPRINT_01 dokümanında dashboard endpoint `GET /admin/transport-platform/dashboard` olarak yazılmış. MASTER_WORKPLAN Bölüm 6.1'de prefix `/admin/transport-platform` olarak verilmiş. İkisi uyumlu; ancak bazı public endpoint'ler `/public/transport` ve `/public/veli` şeklinde iki farklı prefix öneriyor.

**Karar:** Public self-service endpoint'ler için tek prefix kullanılmalı: `/public/transport-platform`. Site context (veli vs. aractakip) host/domain resolver'dan okunur, URL prefix'ten değil. Bu AYSTEK Core'un mevcut site context middleware modeliyle uyumludur.

**Aksiyon:** MASTER_WORKPLAN Bölüm 6.2 güncellendi.

---

## 6. Doğrulanan Legacy Veriler

Legacy PHP projelerinin incelenmesinden çıkarılan ve Codex plan'ına yansıtılan/yansıtılması gereken bulgular:

| Bulgu | Kaynak dosya | Plan'da durum |
|---|---|---|
| Servis ön kayıt: öğrenci, okul, sınıf, adres, doğum tarihi, kan grubu, anne/baba, KVKK onayı | `ServisOnKayitForm.php` | ✅ `DynamicFieldDefinition` ve `ServiceRequest` ile karşılanıyor |
| Servisten ayrılma: e-posta + servis çağrısı | `ServistenAyrilmaTalebiForm.php` | ✅ `ServiceRequest` state machine içinde |
| Öğrenci girişi: çocuk TC + okul + SMS kodu | `OgrenciSMSOnay.php` | ✅ `GuardianLink + OTP` akışı ile karşılanıyor |
| Kardeş indirimi: öğrenci TC listesi + dosya yükleme | `KardesIndirimi.php` | ⚠️ Plan'da genel `DiscountRule` var; kardeş kuralı özel modellenmemiş |
| Banka POS: Garanti, İşbank, YapıKredi, Finans, Akbank ayrı akışları | `SiteController.php` | ⚠️ BIN routing genel olarak var; banka profil detayları ADR'a çekildi |
| Havale: banka, tutar, dekont, onay akışı | `HavaleController.php`, `HavaleForm.php` | ✅ `BankTransferNotice` ile karşılanıyor |
| Sürücü/hostes bilgileri | `SurucuHostesBilgileri.php` | ✅ `transport_core` içinde Driver/attendant entity kapsar |
| Sözleşme değişken yer tutucuları ([VELI-TC] gibi) | `VSozlesmesi.php` | ✅ `ContractRequirement` içinde template variable snapshot var |

---

## 7. Oluşturulan / Güncellenen Dokümanlar

Bu rapor kapsamında yapılan değişiklikler:

| Doküman | Eylem | Açıklama |
|---|---|---|
| `MASTER_WORKPLAN.md` | Güncellendi | Tüm 4 domain koleksiyon tablosu, Faz 7 içeriği, SAP dispatch namespace netleştirildi, public route prefix netleştirildi |
| `ADR_001_DOMAIN_BOUNDARIES_AND_NAMESPACE.md` | Güncellendi | SAP dispatch = `transport_core/sapdispatch` kararı netleştirildi |
| `FAZ_0_DISCOVERY_EXECUTION_PLAN.md` | Oluşturuldu | Faz 0 keşif çalışmaları için SPRINT_01 benzeri execution plan |
| `ADR_002_OFFLINE_FIRST_SYNC.md` | Oluşturuldu | Offline-first cihaz senkronizasyon protokolü ADR |
| `ADR_003_MULTI_POS_BANK_ROUTING.md` | Oluşturuldu | Multi-banka POS routing ve BIN yönlendirme ADR |
| `INDEX.md` | Güncellendi | Yeni dokümanlar eklendi |

---

## 8. Kritik Aksiyon Kontrol Listesi

Bu rapor sonrası yapılması gereken işler (öncelik sırasına göre):

- [x] Validation raporu oluşturuldu (bu doküman)
- [x] MASTER_WORKPLAN koleksiyon tablosu tamamlandı
- [x] Faz 7 içeriği tamamlandı
- [x] SAP dispatch namespace netleştirildi
- [x] Public route prefix netleştirildi
- [x] FAZ_0_DISCOVERY_EXECUTION_PLAN.md oluşturuldu
- [x] ADR_002_OFFLINE_FIRST_SYNC.md oluşturuldu
- [x] ADR_003_MULTI_POS_BANK_ROUTING.md oluşturuldu
- [ ] Faz 0 keşif çalışmaları başlatılacak (bkz. FAZ_0_DISCOVERY_EXECUTION_PLAN.md)
- [ ] Kardeş indirimi kuralı Faz 0 sorularına eklendi; Faz 2'de modelleme yapılacak
- [ ] TOSLA entegrasyon sözleşmesi Faz 0'da sorulacak
- [ ] Mobil/role-based uygulama kanalı Faz 0'da kararlaştırılacak
- [ ] SAP sürümü ve bağlantı tipi Faz 0'da netleştirilecek
