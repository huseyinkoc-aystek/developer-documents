# Servisim Nerede Projesi — Master Analiz ve Teknik Mimari Tasarım

**Durum:** İlk master analiz  
**Tarih:** 2026-07-08  
**Sohbet adı:** Servisim Nerede Projesi  
**Hedef platformlar:** Android + iOS, Flutter  
**Bağlı sistemler:** `aystek-core`, `aystek-core-next`, `aysteknoloji-com`, `aystek-core-installer`  
**Zorunlu referanslar:**
- `/workspace/aystek-core-dev-docs/README_FOR_AI_ASSISTANTS.md`
- `/workspace/aystek-core-dev-docs/AYSTEK_CORE_CANONICAL_ARCHITECTURE.md`
- `/workspace/aystek-core-dev-docs/user-guide/06-modules/FLEET_TALLY_USER_GUIDE_TR.md`
- `/workspace/aystek-core-dev-docs/plans/FLEET_TALLY_SYSTEM_ANALYSIS_TR.md`
- `/workspace/aystek-core-dev-docs/plans/KT-ARACGPSCETELETAKIP-YONETIM-VE-VELI.md`

---

## 1. Yönetici Özeti

Bu proje yalnızca “servis haritada nerede?” ekranı değildir. Doğru ürün tanımı şudur:

> Tenant/site bazlı çalışan; okul ve personel taşımacılığında yolcu, veli/hak sahibi, durak, hat, servis seferi, araç konumu, onay, geri bildirim, sözleşme ve tahsilat akışlarını AYSTEK Core altında yöneten mobil self-service ve operasyon platformu.

Bu nedenle yeni çözüm:

- Mevcut `fleet_tally` modülünü **yeniden yazmamalı**; GPS doğrulama, sağlayıcı ve ERP çetele kontrol kabiliyetini bu modülden beslenerek kullanmalıdır.
- Araç ve GPS cihaz varlıklarını `stock/vehicles` ve `stock/vehicle_assets` üzerinden sahiplenmelidir.
- Ödeme için mevcut `payments/intents`, `payments/gateways`, `payments/installments` altyapısını genişletmelidir.
- Sözleşme için mevcut `contracts` ve gerekiyorsa `subscription_contracts` altyapısını policy/version/signature boyutunda genişletmelidir.
- Yolcu, veli, durak, hat, sefer, kayıt ve onay gibi yeni kavramlar için ayrı ve temiz bir bounded context oluşturmalıdır.
- Flutter tarafında tenant/site bazlı statik yayınlanabilen ama backend tarafından dinamik beslenen bir mobil uygulama olarak tasarlanmalıdır.

Önerilen ilk teknik karar: yeni backend bounded context adı **`passenger_transport`**, mobil ürün çalışma adı **`ServisGO`** veya global alternatif olarak **`RouteMate`** olsun. Nihai marka için domain/marka araştırması ayrıca yapılmalıdır.

---

## 2. Değişmez Mimari İlkeler

AYSTEK standartları gereği bu proje şu kurallardan sapmamalıdır:

1. **Tek backend otoritesi:** İş kuralları `aystek-core` içinde olmalı; Flutter yalnızca kullanıcı deneyimi ve güvenli istemci olmalıdır.
2. **Tenant/site bağlamı:** Her kayıt `tenant_id` ve gerekli yerlerde `site_id` taşır. Mobil app de ilgili tenant/site statik konfigürasyonu ile çalışır.
3. **ModuleGate ve permission:** Modül kapalıysa hem admin menüleri gizlenmeli hem API erişimi reddedilmelidir.
4. **DDD ve OOP:** Model davranışları domain modelinde, orkestrasyon servislerde, veri erişimi repository katmanında kalmalıdır.
5. **Mükerrer iş yok:** GPS provider, araç varlığı, ödeme, sözleşme, onay ve bildirim mevcut modüller varken yeniden oluşturulmamalıdır.
6. **Dinamik davranış:** Hatlar, duraklar, aktif saatler, ödeme sağlayıcıları, sözleşme gereksinimleri ve bildirim şablonları kodda sabitlenmemelidir.
7. **i18n ve seed zorunlu:** Go/Next.js tarafında yeni tüm metinler 28 dil translation seed; Flutter tarafında 11 ARB dili ile gelmelidir.

---

## 3. Mevcut Altyapıdan Yeniden Kullanılacak Parçalar

| İhtiyaç | Mevcut sahip | Kullanım kararı |
|---|---|---|
| GPS sağlayıcıları | `fleet_tally/gps_providers` | Trio, Mobiliz, Arvento, N2 gibi sağlayıcı konfigürasyonu yeniden kullanılmalı. |
| GPS çetele doğrulama | `fleet_tally/verifications` | “Servis bu saat/lokasyonda gerçekten geçti mi?” kanıtı için kullanılmalı. |
| Zamanlayıcı ve ERP push | `fleet_tally/schedules`, `fleet_tally/connectors` | SAP/ERP operasyon doğrulama akışında kullanılmalı. |
| Araç kayıtları | `stock/vehicles` | Plaka, araç ana verisi ve operasyonel araç yönetimi burada kalmalı. |
| GPS cihaz/varlık montajı | `stock/vehicle_assets` | Araç-cihaz eşlemesi burada kalmalı; canlı takip sadece okuma yapmalı. |
| Onay süreçleri | `approvals` | Kayıtsız servis seçimi, durak değişikliği, veli/kişi eşleştirme onayı için kullanılmalı. |
| Ödeme altyapısı | `payments/intents`, `payments/gateways`, `payments/installments` | Sanal POS, taksit, 3DS ve ödeme webhooks için genişletilmeli. |
| Sözleşme altyapısı | `contracts`, `subscription_contracts` | Veli/personel sözleşme imzası için versioned requirement/signature katmanı eklenmeli. |
| Bildirim | `notifications`, `sms`, FCM altyapısı | Sefer yaklaştı, gecikme, onaylandı, ödeme/sözleşme hatırlatmaları için kullanılmalı. |
| Entegrasyon | `integrations`, RabbitMQ/EventBridge | Trafik/harita sağlayıcıları, SAP/ERP ve ödeme sağlayıcı event akışları için kullanılmalı. |

---

## 4. Ürün Kapsamı

### 4.1 Mobil Kullanıcı Rolleri

| Rol | Açıklama |
|---|---|
| Personel | Kendi personel servis kaydını, biniş/iniş durağını ve aktif servisini yönetir. |
| Öğrenci | Genellikle veli/hak sahibi tarafından yönetilen taşınan kişi kaydıdır. |
| Veli / Hak sahibi | Bir veya birden fazla öğrenciyi bağlar; sözleşme, borç, ödeme ve bildirimleri yönetir. |
| Sürücü / Hostes | Bu fazda opsiyonel; ileride yoklama, olay bildirimi, rota durum güncelleme yapabilir. |
| Operasyon yöneticisi | Admin panelden hat, sefer, durak, araç, onay ve istisnaları yönetir. |

### 4.2 Temel Mobil Özellikler

1. **Tenant/site bazlı onboarding**
   - Uygulama hangi tenant/site için build edildiyse ilk açılışta o bağlamla çalışır.
   - Logo, renk, uygulama adı, destek bilgisi ve sözleşme linkleri backend config ile gelir.

2. **Kayıt ve kimlik doğrulama**
   - Telefon/e-posta/OTP destekli hesap oluşturma.
   - Veli/personel hesabı ile taşınan kişi kaydı ayrı tutulur.
   - T.C. kimlik, kurum, okul, sınıf, sicil vb. alanlar dinamik alan politikası ile yönetilir.

3. **Kişi ve kurum bağlama**
   - Kullanıcı kendisini, çocuğunu veya yetkili olduğu kişiyi ilgili tenant/site/institution altında bağlar.
   - Otomatik eşleşmeyen kayıtlar onay kuyruğuna düşer.

4. **Durak ve servis seçimi**
   - Kullanıcı biniş ve iniş durağı seçer.
   - Kendisine atanmış servis varsa doğrudan gösterilir.
   - Kayıtlı servis yoksa kullanıcı tercih bildirir; tercih `approval` akışına girer.

5. **Canlı servis takibi**
   - Yalnızca servisin hareket saatleri ve aktif sefer penceresi içinde gösterilir.
   - Araç konumu GPS cihazından alınır; yolcunun telefonu takip kaynağı yapılmaz.
   - Trafik verisine bağlı ETA, gecikme ve rota durumu gösterilir.

6. **Geri bildirim ve olay bildirimi**
   - Gecikme, şoför/hostes, güvenlik, temizlik, rota, durak ve ödeme/sözleşme kategorileri.
   - SLA, kategori ve öncelik admin tarafından yönetilir.

7. **Sözleşme ve onaylar**
   - Veli/personel sözleşmesi, KVKK aydınlatması, servis kullanım şartları.
   - Borç veya dönem değiştiğinde yeniden imza gereksinimi doğabilir.

8. **Borç ve tahsilat**
   - Veli/personel borcu varsa peşin veya taksitli ödeme.
   - Banka sanal POS’u provider-owned hosted payment page / 3DS akışıyla kullanılmalıdır.
   - Kart verisi uygulamada veya AYSTEK frontendinde tutulmamalıdır.

---

## 5. Önerilen Global Ürün İsimleri

| İsim | Karakter | Not |
|---|---|---|
| **ServisGO** | Türkçe + modern | Türkiye pazarı için güçlü ve anlaşılır. |
| **RouteMate** | Global | Öğrenci/personel ayrımı yapmadan rota refakatçisi hissi verir. |
| **ShuttleMate** | Global | Okul/personel shuttle odağı net. |
| **RidePulse** | Global | Canlı takip ve bildirim hissi güçlü. |
| **Aystek Route** | Kurumsal | AYSTEK markasını öne çıkarır; alt marka olarak kullanılabilir. |
| **YolMate** | Hibrit | Türkçe sıcaklık + global okunabilirlik. |

İlk önerim: Türkiye pazarı ve anlaşılabilirlik için **ServisGO**, global ambisyon için **RouteMate**. Kod adı olarak marka bağımsız **`passenger_transport`** kullanılmalıdır.

---

## 6. Önerilen Bounded Context ve Modül Yapısı

Yeni iş alanı için önerilen backend modülü:

```text
internal/modules/passenger_transport/
├── model/
├── dto/
├── repository/
├── service/
├── handler/
├── routes.go
└── subdomains/
    ├── institutions
    ├── riders
    ├── guardian_links
    ├── stops
    ├── routes
    ├── service_lines
    ├── trips
    ├── enrollments
    ├── feedback
    └── mobile_config
```

Not: Kod tabanındaki nihai paket yapısı mevcut modül standartlarına göre sadeleştirilebilir. Ancak kavramsal sahiplik tek bounded context altında tutulmalıdır.

### 6.1 Yeni Domain Kavramları

| Kavram | Açıklama |
|---|---|
| `Institution` | Okul, şirket, kampüs veya taşımacılık hizmetinin bağlı olduğu organizasyon. |
| `TransportPerson` | Taşınan kişi; öğrenci veya personel. |
| `GuardianLink` | Mobil kullanıcı ile taşınan kişi arasındaki yetkili ilişki. |
| `Stop` | Biniş/iniş noktası; koordinat, adres, güvenli bekleme alanı, aktiflik ve yön bilgisi taşır. |
| `TransportRoute` | Durak dizisi ve rota geometrisi. |
| `ServiceLine` | Belirli kurum, kişi türü, yön, gün ve saat politikasına bağlı servis hattı. |
| `Enrollment` | Kişinin bir servis hattına/durağa kayıt talebi ve aktif üyeliği. |
| `Trip` | Belirli tarih/saatte çalışan gerçek servis seferi. |
| `TripVehicleAssignment` | Sefer ile araç/sürücü/hostes eşleşmesi. |
| `TripLocationSnapshot` | Canlı takipte kullanılan son konum ve ETA özeti. |
| `FeedbackCase` | Mobil geri bildirim/şikayet/öneri kaydı. |

### 6.2 State Machine Önerileri

`EnrollmentStatus`:

```text
draft → pending_approval → approved → active
      ↘ rejected
active → change_requested → pending_approval
active → suspended | cancelled | archived
```

`TripStatus`:

```text
planned → preparing → active → completed
        ↘ cancelled
active → delayed → active
```

`FeedbackStatus`:

```text
open → triaged → in_progress → resolved → closed
     ↘ rejected
```

---

## 7. Canlı Takip ve Trafik Verisi Tasarımı

### 7.1 Canlı Takip İlkesi

Canlı takip, tüm gün açık bir araç takip ekranı olmamalıdır. Gizlilik ve operasyon güvenliği için:

- Sefer yalnızca `TripStatus=active|delayed` durumunda mobil kullanıcıya görünür.
- Kullanıcı sadece kendi bağlı olduğu kişi/enrollment üzerinden ilgili seferi görebilir.
- Konum hassasiyeti tenant/site politikasına göre ayarlanabilir.
- Servis bitince canlı konum yerine “sefer tamamlandı” özeti gösterilir.

### 7.2 Veri Kaynakları

| Kaynak | Kullanım |
|---|---|
| GPS provider geçmiş/canlı konum | `fleet_tally/gps_providers` ve araç varlık eşlemesi üzerinden okunur. |
| Stock vehicle asset | Hangi plakanın hangi GPS external ID ile eşleştiğini sağlar. |
| Trafik/ETA provider | Google/Here/Mapbox/Yandex vb. provider adapter olarak tanımlanmalı; kodda sabitlenmemeli. |
| Route/stop geometry | Backend’de tutulan kanonik rota ve durak koordinatları. |

### 7.3 Teknik Akış

```text
Flutter App
  │ GET /mobile/passenger-transport/trips/active
  ▼
aystek-core passenger_transport service
  │ enrollment + trip + permission kontrolü
  │ vehicle assignment çözümü
  ▼
Vehicle/GPS resolver
  │ stock.vehicle_assets + fleet_tally.gps_provider
  ▼
Location/ETA adapter
  │ son konum + rota + trafik provider
  ▼
TripLocationSnapshot
  │
  ├─ REST polling: MVP için 10-20 sn
  ├─ SSE/WebSocket: sonraki faz
  └─ FCM: önemli gecikme/varış bildirimleri
```

MVP için REST polling yeterlidir. Üretim ölçeklenince SSE veya WebSocket, sadece aktif sefer pencerelerinde açılmalıdır.

---

## 8. Mobil Flutter Mimari Kararı

### 8.1 Yeni App mi, Mevcut `kisisel_asistan` mı?

Önerilen karar: **yeni Flutter uygulaması**.

Gerekçe:

- Bu ürün B2B/B2B2C taşımacılık, veli, ödeme ve sözleşme süreçleri içerir.
- `kisisel_asistan` kişisel haber/AI/eşya/seyahat bağlamında konumlanmış durumda.
- Store listing, ikon, izin açıklamaları, ödeme ekranı, çocuk/veli politikaları ve marka dili farklı olacaktır.
- Yeni uygulama mükerrer kod yazmadan `kisisel_asistan` içindeki HTTP client, Riverpod, GoRouter, l10n ve secure storage pattern’larını örnek alabilir.

Önerilen path:

```text
/workspace/flutter/servis_go/
```

Kod tarafında marka bağımsız feature yapısı:

```text
lib/
├── core/
│   ├── config/
│   ├── router/
│   ├── theme/
│   ├── l10n/
│   └── network/
├── features/
│   ├── auth/
│   ├── tenant_site/
│   ├── riders/
│   ├── stops/
│   ├── enrollments/
│   ├── trips/
│   ├── feedback/
│   ├── contracts/
│   └── payments/
└── shared/
    ├── widgets/
    ├── models/
    └── services/
```

### 8.2 Tenant/Site Static Build Modeli

Flutter uygulama tenant/site bazlı static çalışabilir:

```bash
flutter build apk \
  --dart-define=TENANT_CODE=gursel-turizm \
  --dart-define=SITE_CODE=veli \
  --dart-define=API_BASE_URL=https://api.aystek... \
  --dart-define=HOST_HEADER=veli.gurseltur.com.tr
```

Ancak tenant/site davranışı yine backend’den gelen `mobile_config` ile yönetilmelidir:

- Brand adı ve logo
- Tema renkleri
- Aktif modüller
- Ödeme aktif/pasif
- Sözleşme zorunlulukları
- Harita/ETA provider seçimi
- Destek iletişimleri
- KVKK ve gizlilik linkleri

### 8.3 Flutter Paketleri

Mevcut `kisisel_asistan` standartları korunabilir:

- `flutter_riverpod`
- `go_router`
- `flutter_secure_storage`
- `http`
- `shared_preferences`
- `flutter_local_notifications`
- `firebase_messaging`
- `url_launcher`
- `geolocator` yalnızca kullanıcıya yakın durak önerisi gibi açık izinli senaryolarda

Harita için ayrıca değerlendirme gerekir:

- `google_maps_flutter`
- `mapbox_maps_flutter`
- `flutter_map`

Harita sağlayıcısı ürün ve lisans kararına bağlıdır; kodda provider lock-in yaratılmamalıdır.

---

## 9. Ödeme ve App Store / Play Store Politikası

Bu ürün fiziksel dünyada tüketilen taşımacılık hizmeti için ödeme alır. Bu nedenle dijital içerik/abonelik gibi App Store IAP veya Google Play Billing zorunluluğu doğuran bir ürün gibi modellenmemelidir.

Resmi kaynaklara göre:

- Apple App Review Guidelines 3.1.3(e), uygulama dışında tüketilen fiziksel mal/hizmet ödemelerinde in-app purchase yerine Apple Pay veya geleneksel kart girişi gibi yöntemlerin kullanılacağını belirtir: https://developer.apple.com/app-store/review/guidelines/
- Google Play Billing dokümanı, Google Play billing sisteminin dijital ürünler için olduğunu; fiziksel ürün/hizmet ve dijital olmayan içerikler için Google Pay SDK gibi çözümlere bakılacağını belirtir: https://developer.android.com/google/play/billing
- Google Play Payments Policy açıklamaları, dijital ürün/hizmetlerde Play Billing zorunluluğunu vurgular; fiziksel/dış dünyada tüketilen hizmetler bu sınıfa girmez: https://support.google.com/googleplay/android-developer/answer/10281818

### 9.1 Teknik Öneri

- Mobil uygulama kart verisi toplamaz.
- Ödeme `PaymentIntent` üzerinden backend’de başlar.
- Backend provider-owned hosted payment page veya provider-controlled 3DS sayfası üretir.
- Flutter, ödeme ekranını tercihen sistem tarayıcısı / custom tab / SFSafariViewController ile açar.
- WebView kullanılacaksa provider ve mağaza kuralları ayrıca kontrol edilir; kart alanı AYSTEK HTML/JS ile üretilmez.
- Callback/return URL backend tarafından doğrulanır; mobil app yalnızca sonucu sorgular.

Bu yaklaşım hem PCI kapsamını küçültür hem de mağaza politikalarını “kaçınma” değil, fiziksel hizmet ödeme modeli olarak doğru çerçeveye oturtur.

---

## 10. API Taslakları

### 10.1 Mobil API

```text
GET    /mobile/passenger-transport/config
POST   /mobile/passenger-transport/auth/register
POST   /mobile/passenger-transport/auth/login
POST   /mobile/passenger-transport/auth/refresh

GET    /mobile/passenger-transport/me
GET    /mobile/passenger-transport/institutions
POST   /mobile/passenger-transport/guardian-links
GET    /mobile/passenger-transport/riders
POST   /mobile/passenger-transport/riders

GET    /mobile/passenger-transport/stops
GET    /mobile/passenger-transport/service-lines
POST   /mobile/passenger-transport/enrollments
GET    /mobile/passenger-transport/enrollments
POST   /mobile/passenger-transport/enrollments/:id/change-request

GET    /mobile/passenger-transport/trips/active
GET    /mobile/passenger-transport/trips/:id/location
POST   /mobile/passenger-transport/feedback

GET    /mobile/passenger-transport/contracts/requirements
POST   /mobile/passenger-transport/contracts/:id/accept
GET    /mobile/passenger-transport/debts
POST   /mobile/passenger-transport/payments/quote
POST   /mobile/passenger-transport/payments/intents
GET    /mobile/passenger-transport/payments/intents/:id
```

### 10.2 Admin API

```text
GET/POST/PUT/DELETE /admin/passenger-transport/institutions
GET/POST/PUT/DELETE /admin/passenger-transport/riders
GET/POST/PUT/DELETE /admin/passenger-transport/stops
GET/POST/PUT/DELETE /admin/passenger-transport/routes
GET/POST/PUT/DELETE /admin/passenger-transport/service-lines
GET/POST/PUT        /admin/passenger-transport/enrollments
POST                /admin/passenger-transport/enrollments/:id/approve
POST                /admin/passenger-transport/enrollments/:id/reject
GET/POST/PUT        /admin/passenger-transport/trips
POST                /admin/passenger-transport/trips/:id/assign-vehicle
GET                 /admin/passenger-transport/feedback
POST                /admin/passenger-transport/feedback/:id/status
GET/PUT             /admin/passenger-transport/mobile-config
```

---

## 11. Önerilen MVP Kapsamı

### MVP-1: Takip ve Kayıt Çekirdeği

- Tenant/site mobile config
- Kullanıcı kayıt/giriş
- Kişi bağlama
- Durak listesi
- Servis/durak kayıt talebi
- Admin onay akışı
- Aktif sefer ve canlı konum görüntüleme
- Temel geri bildirim

### MVP-2: Operasyonel Yönetim

- Admin: kurum, kişi, durak, hat, sefer yönetimi
- Araç atama
- GPS provider çözümü
- Gecikme/ETA bildirimi
- FCM push bildirimleri

### MVP-3: Sözleşme ve Ödeme

- Veli/personel sözleşme gereksinimi
- Borç hesabı görüntüleme
- Ödeme quote
- Hosted payment page / sanal POS
- Ödeme sonucu ve makbuz
- SAP/ERP dispatch eventleri

### MVP-4: Sertleştirme ve Ölçek

- SSE/WebSocket canlı takip
- Trafik provider adapter
- Rota optimizasyonu
- Offline düşük bağlantı modu
- Gelişmiş raporlama ve audit

---

## 12. Event Modeli

Yeni modül event-driven tasarlanmalıdır. Önerilen event key’ler:

```text
passenger_transport.rider.created
passenger_transport.guardian_link.requested
passenger_transport.guardian_link.approved
passenger_transport.enrollment.requested
passenger_transport.enrollment.approved
passenger_transport.enrollment.rejected
passenger_transport.stop_change.requested
passenger_transport.trip.started
passenger_transport.trip.delayed
passenger_transport.trip.completed
passenger_transport.feedback.created
passenger_transport.contract.accepted
passenger_transport.payment.intent_created
passenger_transport.payment.succeeded
```

AYSTEK standardı gereği event publish edilen admin list ekranlarında `EventKeysInfo` ile bu eventlerin payload alanları gösterilmelidir.

---

## 13. Güvenlik, KVKK ve Gizlilik Notları

- Çocuk/öğrenci verisi yüksek hassasiyetli kabul edilmelidir.
- Veli-çocuk ilişkilendirme audit trail ile tutulmalıdır.
- Kan grubu ve sağlık benzeri alanlar özel nitelikli veri olarak işaretlenmeli, varsayılan maskeli görünmelidir.
- Canlı araç konumu sadece yetkili enrollment sahiplerine ve aktif sefer penceresinde gösterilmelidir.
- Kullanıcı telefonu ile yolcu takibi yapılmamalıdır; yolcu konumu ancak açık izinli ve ayrı gerekçeli bir özellik olarak değerlendirilmelidir.
- Ödeme kart verisi AYSTEK sistemlerinde tutulmamalıdır.
- Her ödeme, sözleşme kabulü, ilişkilendirme, onay ve canlı takip erişimi audit edilebilir olmalıdır.

---

## 14. Test Stratejisi

| Katman | Kapsam |
|---|---|
| Unit test | State machine, enrollment policy, trip visibility, ETA fallback, payment quote rules. |
| Repository test | Tenant/site scoped filtreler ve index beklentileri. |
| Service test | Kişi bağlama, servis tercih onayı, aktif sefer erişim kontrolü. |
| Integration test | GPS provider resolver, payment intent callback, contract requirement. |
| Flutter widget test | Login, durak seçimi, aktif sefer, ödeme dönüş ekranı. |
| E2E/UAT | Kayıt → kişi bağlama → servis seçimi → onay → canlı takip → geri bildirim → ödeme. |
| Security test | IDOR, tenant leakage, canlı konum yetki sınırı, rate limit, OTP brute force. |

---

## 15. Definition of Done

Yeni backend modülü için:

- Model, DTO, repository, service, handler, routes yapısı tamam.
- Service unit testleri var.
- Tenant/site filtreleri repository seviyesinde zorunlu.
- Permission seed ve role auto-assign uyumlu.
- Menu seed eklendi.
- Translation seed 28 dil tamamlandı.
- Event publish varsa frontend `EventKeysInfo` eklendi.
- Admin UI sayfaları backend-driven menü ile erişilir.
- User guide ve operations dokümanları güncellendi.

Flutter için:

- Kullanıcıya görünen tüm metinler ARB içinde 11 dil.
- Riverpod provider’ları domain bazlı ayrıldı.
- GoRouter rotaları feature bazlı net.
- Secure storage token akışı var.
- Build flavor / dart-define tenant-site config destekli.
- Ödeme ekranı kart verisi toplamıyor.
- Android/iOS izin açıklamaları gerçek kullanıma uygun.

---

## 16. Tahmini Tamamlanma Süresi

Bu proje, AYSTEK standartlarına göre mevcut modülleri yeniden yazmadan ama eksik bounded context’i doğru inşa ederek ilerlemelidir. Bu nedenle süre tahmini iki seviyede düşünülmelidir:

| Kapsam | İçerik | Tahmini süre |
|---|---|---:|
| **MVP-1** | Kayıt/giriş, kişi bağlama, durak seçimi, servis tercih/onay akışı, aktif sefer ve temel canlı takip | **4–6 hafta** |
| **MVP-2** | Admin kurum/durak/hat/sefer yönetimi, araç atama, GPS resolver, temel bildirimler | **3–5 hafta** |
| **MVP-3** | `contracts` tabanlı sözleşme gereksinimleri, `payments` tabanlı sanal POS/taksit/3DS ödeme akışı | **4–6 hafta** |
| **MVP-4** | Trafik/ETA adapter, SSE/WebSocket, audit/security sertleştirme, UAT, store hazırlıkları | **3–5 hafta** |

**Gerçekçi toplam:** üretime alınabilir V1 için **14–22 hafta**.

Bu tahmin, aşağıdaki ekip varsayımıyla verilmiştir:

- 2 backend geliştirici
- 1 Next.js/admin panel geliştirici
- 1 Flutter geliştirici
- 1 QA/automation veya part-time QA
- Ürün sahibi / iş analisti desteği

Daha küçük ekiple süre **20–30 haftaya** uzayabilir. GPS provider, banka sanal POS, SMS/OTP, SAP/ERP ve legacy veri migrasyonu erken netleşmezse takvimde ayrıca **2–6 hafta risk payı** bırakılmalıdır.

Kritik not: Sözleşme ihtiyaçları yeni bir sözleşme modülüyle değil, mevcut **`contracts`** modülü ve gerekiyorsa onun policy/version/signature genişletmeleriyle çözülmelidir. Aynı kural ödeme için **`payments`**, araç/GPS varlıkları için **`stock/vehicle_assets`**, GPS doğrulama için **`fleet_tally`**, onay için **`approvals`** modülleri açısından da geçerlidir.

---

## 17. Açık Kararlar

1. İlk hedef müşteri/tenant/site hangisi olacak?
2. MVP önce okul servisi mi, personel servisi mi, ikisi birlikte mi?
3. İlk GPS provider canlı konum için hangisi olacak?
4. Trafik/ETA sağlayıcısı Google, Here, Mapbox, Yandex veya başka bir provider mı olacak?
5. Banka sanal POS sağlayıcısı ve 3DS yöntemi hangisi?
6. SMS/OTP sağlayıcısı mevcut `sms` modülüyle mi, yeni provider adapter ile mi kullanılacak?
7. Flutter uygulama adı nihai olarak ne olacak?
8. Veli/personel sözleşmelerinin hukuki metinleri ve versiyonlama kuralı hazır mı?
9. Eski veli/servis/borç verisi için migration kaynağı neresi?
10. Sürücü/hostes mobil rolü MVP’ye dahil mi, sonraki faz mı?

---

## 18. İlk Uygulama Planı

1. **Discovery checkpoint**
   - Tenant/site, ilk kurum, GPS provider, ödeme provider ve MVP rol kapsamını netleştir.

2. **Backend foundation**
   - `passenger_transport` domain modeli, repository, service ve admin/mobile handler iskeleti.

3. **Admin foundation**
   - Kurum, durak, hat, sefer ve enrollment onay ekranları.

4. **Flutter MVP**
   - Yeni app scaffold, tenant/site config, auth, kişi bağlama, durak seçimi, aktif sefer takip ekranı.

5. **GPS resolver**
   - Stock vehicle asset + fleet_tally provider üzerinden canlı/semi-live konum okuma.

6. **Feedback + notifications**
   - Geri bildirim case ve temel push bildirimleri.

7. **Contracts + payments**
   - Borç/sözleşme requirement, hosted payment flow, callback ve audit.

---

## 19. Sonuç

Bu proje AYSTEK ekosisteminde ayrı bir “küçük mobil app” olarak değil, **passenger transport self-service platformu** olarak ele alınmalıdır.

En doğru mimari, mevcut `fleet_tally`, `stock`, `payments`, `contracts`, `approvals`, `notifications` ve `integrations` yeteneklerini yeniden kullanıp; yalnızca yolcu/veli/durak/hat/sefer/enrollment domainini yeni bir bounded context olarak eklemektir.

Bu yaklaşım:

- Mükerrer çalışmayı engeller.
- AYSTEK Core’un tenant/site/modül mimarisine uyar.
- Flutter mobil uygulamayı sade ve güvenli tutar.
- GPS, trafik, ödeme ve sözleşme gibi değişken alanları provider/policy tabanlı genişletir.
- İleride veli kaydı, sözleşme, borç tahsilatı, taksit, SAP/ERP entegrasyonu ve raporlama özelliklerini aynı omurgada büyütür.
