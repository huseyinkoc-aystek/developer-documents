# Kapsam Birleştirme, Fark Analizi ve Hedef Mimari

**Tarih:** 2026-05-09  
**Bağlam:** Kullanıcı anlatımı + birleşik enterprise rapor + AYSTEK Core standartları  
**Amaç:** Dağınık ihtiyaçları temiz, uygulanabilir ve mükerrer geliştirme üretmeyen bir enterprise hedef mimariye dönüştürmek.

---

## 1. Yönetici Özeti

Gönderilen yeni metin, önceki birleşik rapordan daha geniş bir ürün vizyonu tarif eder. Önceki raporun merkezi `veli.gurseltur.com.tr` üzerinden **veli/personel ödeme, borç, sözleşme, kan grubu, havale ve SAP dispatch** idi. Yeni metin buna ek olarak daha büyük bir omurga getirir:

- Üniversite, okul ve personel servisleri için **kartlı/QR’lı geçiş sistemi**.
- Kullanıcı bazlı **cüzdan/bakiye/yükleme/iade** yapısı.
- Mifare kart, QR, misafir bilet, abonelik ve indirimli/ücretsiz geçiş.
- Güzergâh, yön, gün/saat, geçiş kuralı, route-based ücret ve tarihsel fiyat yönetimi.
- Şoför/biletçi ekranı, araç sahibi ekranı, hizmet alan firma ekranı, taşınan kişi ekranı.
- GPS sağlayıcılarından gelen verilerle **araç hareket raporu, çetele ve doğrulanmış iş**.
- İş ataması olmadan ring/sefer raporlayabilen yeni nesil akış.
- Offline-first geçiş ve senkronizasyon.
- Üniversite içi yemekhane/fotokopi gibi ek hizmetlerde kullanılabilir genel **campus access/value** altyapısı.
- Araç sahibi onayıyla doğrulanmış iş verisinin mikro servis/API üzerinden hizmet alan firmaya paylaşılması.

Bu nedenle doğru ürün adı artık yalnızca “Veli Tahsilat” değildir. Doğru çatı:

**AYSTEK Taşımacılık, Geçiş, Çetele ve Veli/Personel Tahsilat Platformu**

Bu çatı altında iki ana ürün hattı bulunur:

1. **KT / Mobility Access hattı:** Kart/QR geçiş, cüzdan, bilet, abonelik, çetele, GPS doğrulama, araç hareket raporu, hizmet alan/hizmet veren faturalama.
2. **Veli/Personel Tahsilat hattı:** Borç hesabı, sözleşme, dinamik kişi alanları, ödeme quote, sanal POS, havale, SAP dispatch, public self-service.

Bu iki hat aynı `Person`, `Institution`, `Tenant/Site`, `Payment`, `Ledger`, `Contracts`, `Integrations`, `Telematics` ve `Audit` omurgasını paylaşmalıdır. Ayrı ayrı yazılırsa sistem büyüdükçe birbirini tekrar eder ve AYSTEK Core’un platform avantajı kaybolur.

---

## 2. Birleşik Rapor ile Yeni Kullanıcı Metni Arasındaki Kritik Farklar

| Alan | Birleşik raporda ağırlık | Yeni metindeki ek/kritik ihtiyaç | Mimari sonuç |
|---|---|---|---|
| Ürün odağı | Veli/personel borç ödeme ve self-service | KT kart/QR taşımacılık, cüzdan, geçiş, çetele | Tek `transport_portal` yetmez; `transport_core`, `transport_access`, `transport_wallet` gerekir |
| Bakiye | Borç hesabı ve ödeme intent | Kullanıcı bazlı bakiye, eksi bakiye/kredi limiti, kart değişiminde devir yok | Cüzdan aggregate’ı karttan bağımsız olmalı |
| Kart/QR | Payment card BIN/POS yönlendirme | Mifare ulaşım kartı, QR bilet, QR ile sefere kayıt, QR borç geçişi | Payment card ile ulaşım kartı kavramsal olarak ayrılmalı |
| Sefer/çetele | SAP dispatch ve ödeme | Araç hareket raporu, taşınan kişi sayısı, km, route/GPS doğrulama | TripLedger ve MovementReport aggregate’ları gerekir |
| Ring | Yok veya zayıf | İş ataması olmadan her turu raporlama ihtiyacı | Schedule-free trip inference engine gerekir |
| Offline | Dolaylı | Sunucu/internet yokken geçişlerin aksamaması en kritik gereksinim | Offline transaction log, idempotent sync, conflict resolution gerekir |
| Üniversite servisleri | Okul/firma borç bağlamı | Sabancı/Özyeğin gibi kurumlar, gruplar, indirimler, raporlar | Institution + group + fare policy modeli genişlemeli |
| Campus hizmetleri | Yok | Yemekhane/fotokopi gibi ek kullanım alanları | Transport dışı ama aynı wallet/access altyapısını kullanabilen `service_category` modeli |
| Dış veri paylaşımı | SAP ağırlıklı | Araç sahibi onayıyla hizmet alan firmalara doğrulanmış iş API paylaşımı | API sharing + approval + data product yaklaşımı |
| Faturalama | SAP’ye ödeme gönderimi | Gürsel’e veya hizmet alan firmaya/araç sahibine fatura | Merchant/billing responsibility modeli gerekir |
| Taşınan kişi deneyimi | Borç ve sözleşme | Araç yaklaşma bildirimi, konum, biniş/iniş, hizmet takibi | Mobile/public user experience ayrı tasarlanmalı |

---

## 3. Temizlenmiş Hedef Kapsam

### 3.1 Ana Amaç

AYSTEK Core içinde multi-tenant ve multi-site kurallarına uygun, üniversite/okul/personel taşımacılığı için tekil bir platform kurulacaktır. Platform:

- Taşınan kişileri, kurumları, grupları, güzergâhları ve ücret politikalarını yönetecek.
- Kartlı, QR’lı, misafir, ücretsiz, indirimli, abonelikli ve borçlu geçişleri işleyecek.
- Kullanıcı bazlı bakiye/cüzdan modelini yönetecek.
- Araç, şoför, biletçi, hizmet alan firma ve hizmet veren taşıyıcı rollerini ayrıştıracak.
- GPS/telematics verisiyle sefer ve çetele doğrulaması yapacak.
- Veli/personel borç ödeme portalını aynı kişi/kurum/ödeme/sözleşme omurgasıyla çalıştıracak.
- SAP ve diğer müşteri sistemlerine idempotent, izlenebilir, event-driven veri gönderecek.
- Offline durumlarda geçişleri kaybedilmeden çalıştıracak.

### 3.2 Kapsam Dışı Değil, Ama Fazlı Alınacak Konular

Bu maddeler vizyondadır; ancak ilk foundation sprint’ine alınmamalıdır:

- Gerçek POS capture.
- Gerçek SAP üretim dispatch.
- Gerçek Mifare okuyucu firmware veya native driver geliştirme.
- Gerçek mobil uygulama yayınlama.
- Otomatik konumdan biniş/iniş çıkarımı.
- Yemekhane/fotokopi üretim entegrasyonları.

Bunlar için önce domain, event, offline ve API contract zemini kurulmalıdır.

---

## 4. Hedef Domain Modeli

### 4.1 Ortak Çekirdek

| Kavram | Açıklama |
|---|---|
| `Tenant` | AYSTEK müşterisi veya platform sahibi bağlamı. Örn: `gursel-turizm` |
| `Site` | Domain ve kullanım kanalı. Örn: `aractakip`, `veli`, ileride `sabanci-portal` |
| `Institution` | Hizmet alan okul, üniversite, firma veya kampüs |
| `ServiceProvider` | Gürsel Turizm veya başka taşıma hizmeti sunan firma |
| `VehicleOwner` | Araç sahibi firma veya şahıs |
| `Vehicle` | Taşımayı yapan araç; mevcut stock/vehicle-assets ve telematics ile ilişkilidir |
| `Driver` | Şoför; kullanıcı ve/veya personel kaydıyla ilişkilidir |
| `TransportedPerson` | Öğrenci, personel, öğretim görevlisi veya taşınan kişi |
| `PersonGroup` | İndirim ve geçiş kurallarının bağlandığı çoklu grup |
| `Route` | Güzergâh; yönlü tutulur. Örn: `Sabancı → Kadıköy`, `Kadıköy → Sabancı` |
| `Trip` | Belirli tarih/saat ve araç üzerinde gerçekleşen sefer/tur |
| `AccessTransaction` | Biniş/iniş/okutma/QR kullanım olayı |
| `MovementReport` | Sefer tamamlandığında oluşan araç hareket raporu ve çetele |

### 4.2 Wallet ve Ödeme Çekirdeği

| Kavram | Açıklama |
|---|---|
| `WalletAccount` | Kullanıcı bazlı bakiye hesabı; karttan bağımsızdır |
| `WalletInstrument` | Mifare kart, QR kimlik, mobil token gibi erişim aracı |
| `TopUpIntent` | Kullanıcı/TOSLA/POS ile bakiye yükleme niyeti |
| `WalletTransaction` | Yükleme, harcama, iade, düzeltme, borçlandırma hareketi |
| `CreditPolicy` | Grup/kurum bazlı eksi bakiye limiti |
| `RefundWindowPolicy` | Örn. geçişten sonra 15 dakika içinde iade |
| `InvoiceResponsibility` | Fatura TOSLA’ya mı, kullanıcıya mı, kuruma mı kesilecek |

### 4.3 Veli/Personel Tahsilat Çekirdeği

| Kavram | Açıklama |
|---|---|
| `GuardianUser` / `PortalUser` | Portala kendi bilgileriyle kayıt olan kullanıcı |
| `GuardianLink` | Kullanıcının çocuk/personel/taşınan kişi ile doğrulanmış ilişkisi |
| `DebtAccount` | Kişi + kurum + dönem + hizmet tipi ekseninde borç hesabı |
| `DebtItem` | Borç hesabındaki satırlar |
| `ContractRequirement` | Her yeni borç veya borç versiyonu için sözleşme gereksinimi |
| `ContractSignature` | Checkbox, OTP veya e-imza kanıtı |
| `DynamicPersonAttribute` | Kan grubu gibi kişi bazlı dinamik veri |
| `PaymentQuote` | Kart girmeden önce provisional, kart/BIN sonrası final ödeme teklifi |
| `BankTransferNotice` | Havale bildirimi, dekont, operasyon onayı |

### 4.4 Campus Service Genişleme Noktası

Yemekhane, fotokopi gibi servisler ulaşım domain’ini kirletmemelidir. Bunun yerine geçiş ve wallet altyapısı genel bir kullanım kategorisi taşımalıdır:

- `service_category`: `transport`, `cafeteria`, `copy_center`, `parking`, `other`
- `service_point_id`: Okuyucu/turnike/yemekhane kasası/fotokopi noktası
- `fare_policy_id` veya `price_policy_id`
- `entitlement_policy_id`

Bu sayede KT ulaşım sistemi büyürken uygulama “her şeyi tek modüle dolduran” bir yapıya dönüşmez.

---

## 5. Backend Namespace Kararı

### 5.1 Neden Tek `transport_portal` Yetmez?

`transport_portal` ismi veli ödeme ve self-service için uygundur; fakat:

- Kart/QR okutma,
- offline cihaz senkronizasyonu,
- güzergâh/sefer,
- araç hareket raporu,
- kullanıcı cüzdanı,
- yemekhane/fotokopi gibi erişim noktaları

portal değil, operasyonel çekirdek domain’dir. Bu yüzden namespace genişletilmelidir.

### 5.2 Önerilen Paketler

| Namespace | Sorumluluk |
|---|---|
| `transport_core` | Kurum, kişi, grup, güzergâh, sefer, çetele, hareket raporu |
| `transport_access` | Kart/QR/Mifare, cihazlar, geçiş işlemleri, offline sync, bilet/abonelik |
| `transport_wallet` | Kullanıcı bazlı bakiye, yükleme, iade, kredi limiti, fatura sorumluluğu |
| `transport_portal` | Veli/personel portalı, borç, sözleşme, dinamik alan, havale, duyuru |

Mevcut platform modülleri:

- `payments`: sanal POS, payment intent, payment gateways, installment rules.
- `contracts`: template/content altyapısı.
- `ledger`: muhasebesel hareket ve defter mantığı.
- `integrations` / `api_sharing`: SAP, müşteri web servisleri, doğrulanmış iş paylaşımı.
- `telematics`: GPS/cihaz/araç konumu.
- `fleet_tally`: çetele/GPS doğrulama deneyiminden reuse edilecek mevcut bilgi.
- `sms` / `notifications`: OTP, yaklaşma bildirimi, operasyon bildirimi.

---

## 6. KT / Mobility Access Akışları

### 6.1 Kayıtlı Kullanıcı Kart/QR Geçişi

1. Kullanıcı `TransportedPerson` olarak kurum/gruplara bağlıdır.
2. Kullanıcının `WalletAccount` kaydı vardır.
3. Mifare kart veya mobil QR `WalletInstrument` olarak hesaba bağlıdır.
4. Araçtaki cihaz offline/online durumda okutmayı alır.
5. Sistem route/trip/fare/group/credit policy çözümlemesi yapar.
6. Bakiye yeterliyse harcama, yetersiz ama kredi limiti uygunsa eksi bakiye oluşur.
7. 15 dakika iade penceresi varsa geçiş iptal/iade edilebilir.
8. Transaction vehicle, driver, route, location ve device context ile saklanır.
9. Trip kapanınca `MovementReport` ve çetele oluşur.

### 6.2 Misafir QR Bilet

1. Misafir route bazlı tek kullanımlık QR bilet satın alır.
2. Biletin süre aşımı olmayabilir; kullanım tek seferle sınırlıdır.
3. Bilet okutulduğunda route/vehicle/device context ile kullanılmış olur.
4. Misafir geçişleri ayrı raporlanır.

### 6.3 Abonelik

1. Abonelik route bazlı veya genel olabilir.
2. Geçerlilik satın alma sonrası gün sayısı veya tarih aralığı olabilir.
3. Geçiş kuralları haftalık/günlük tekrar, saat aralığı ve tekrar okutma limitleriyle çalışır.
4. Abonelik kullanımı wallet hareketinden ayrı ama ledger/audit ile ilişkili tutulur.

### 6.4 İş Atamasız Ring / Sefer Çıkarımı

Yeni nesil hedef, her seferden önce manuel iş ataması olmadan rapor üretebilmektir. Bunun için:

- Araç GPS akışı,
- route geofence’leri,
- ilk/son geçiş zamanı,
- cihaz oturumu,
- şoför session’ı,
- kurum/güzergâh saat penceresi

birlikte değerlendirilerek `TripInference` yapılmalıdır.

İlk sürümde tam otomasyon yerine “önerilen sefer → operatör onayı” modeli tercih edilmelidir. Böylece yanlış sefer çıkarımı finansal rapora doğrudan yansımaz.

### 6.5 Konumdan Biniş/İniş Varsayımı

Personel/okul servisleri için taşınan kişinin mobil konumu ile araç konumu aynı yönde ve yakın ilerliyorsa biniş; ters yönde uzaklaşıyorsa iniş varsayılabilir.

Bu özellik hassas konum verisi içerdiği için:

- Açık kullanıcı izni,
- minimum veri saklama,
- düşük frekanslı ve amaçla sınırlı takip,
- maskelenmiş/özet audit,
- manuel override

zorunlu olmalıdır. İlk fazda “kanıt güçlendirme sinyali” olarak kullanılmalı, tek başına ücret/fatura kaynağı yapılmamalıdır.

---

## 7. Veli/Personel Tahsilat Akışı

Bu akış önceki birleşik raporla uyumludur; yeni metinde netleşen detaylarla şu hale getirilmelidir:

1. Veli/personel kendi bilgileriyle portala kayıt olur.
2. Hesap onaylanır; e-posta ve tercihen SMS doğrulaması uygulanır.
3. Kullanıcı çocuk/personel sorgular: TC/pasaport + okul/firma + kendi ilişki bilgisi.
4. Legacy/SAP/Core kaynaklarında eşleşme bulunursa `GuardianLink` oluşur; riskli durumlar manuel incelemeye gider.
5. Kişi için eksik zorunlu dynamic attributes varsa bir kez tamamlatılır. Kan grubu kişi bazlıdır, borç bazlı tekrar istenmez.
6. Dönem bazlı borç hesapları listelenir.
7. Borç hesabı için yeni borç veya versiyon değişikliği varsa sözleşme yeniden imzalatılır.
8. Kullanıcı kart girmeden önce olası taksitleri görür.
9. Kart BIN/IIN girildiğinde nihai POS ve taksit seçenekleri çıkar.
10. 3D Secure ile ödeme alınır veya havale bildirimi alınır.
11. Ödeme/havale onayı SAP’ye idempotent dispatch edilir.
12. Her adım audit, trace, event ve operasyon kayıtlarına düşer.

---

## 8. Offline-First Zorunluluk

Geçiş sistemi sunucu veya internet kesintisinde durmamalıdır.

Minimum tasarım:

- Cihaz lokalinde signed append-only transaction log.
- Son bilinen fare/group/route policy snapshot.
- Offline risk limitleri: maksimum offline işlem, maksimum tutar, maksimum süre.
- Her offline transaction için device sequence number.
- Online olunca idempotent sync.
- Çakışma çözümü: duplicate okutma, iade penceresi, route mismatch, balance conflict.
- Offline fraud raporu ve device health audit.

Offline modda yapılmaması gerekenler:

- Yeni yüksek riskli abonelik tanımlama.
- Limitsiz eksi bakiye açma.
- Eski policy snapshot ile sınırsız çalışma.
- PAN/CVV veya hassas ödeme verisi saklama.

---

## 9. Raporlama ve Faturalama

### 9.1 Profil Bazlı Raporlar

| Profil | Rapor ihtiyacı |
|---|---|
| Gürsel Turizm | Kurum bazlı gelir, genel kasa, geçiş, çetele, SAP dispatch |
| Üniversite/Firma | Kendi öğrencileri/personeli, güzergâh kullanımı, indirim ve yükleme raporu |
| Araç sahibi | Araç bazlı yapılan iş, doğrulanmış sefer, km, taşınan kişi sayısı |
| Şoför/Biletçi | Günlük sefer, okutma, iade, uyarı |
| Taşınan kişi | Bakiye, yükleme, geçiş geçmişi, servis konumu |
| Operasyon | Hata, offline sync, fraud/kaçak işlem, reconciliation |

### 9.2 Fatura Sorumluluğu

Faturalar duruma göre:

- Gürsel Turizm’e,
- hizmet alan firmaya,
- kullanıcıya,
- TOSLA gibi yükleme aracısına,
- araç sahibi/hizmet veren taraflara

bağlanabilir.

Bu yüzden her finansal olayda şu alanlar bulunmalıdır:

- `bill_to_type`
- `bill_to_id`
- `merchant_of_record`
- `service_provider_id`
- `vehicle_owner_id`
- `institution_id`
- `tax_profile_id`
- `external_invoice_ref`

---

## 10. API Sharing ve Doğrulanmış İş Paylaşımı

Araç sahipleri veya taşıma firmaları, hizmet verdikleri firmalara doğrulanmış işi mikro servisler üzerinden paylaşabilmelidir.

Bu özellik mevcut `integrations/api_sharing` yaklaşımıyla kurulmalıdır:

- Veri paylaşımı owner approval gerektirir.
- Paylaşılacak veri minimum kapsamlı olur.
- Tenant/site/institution scope zorunludur.
- API client, grant, scope, audit ve rate limit uygulanır.
- Verinin kaynağı `MovementReport` ve imzalı/denetlenmiş event olmalıdır.

Örnek data product:

- `transport.movement_reports.read`
- `transport.verified_trips.read`
- `transport.access_transactions.summary.read`
- `transport.invoice_basis.read`

---

## 11. İlk Geliştirme Sırası

Doğru hedefe gitmek için ilk sprint “ödeme ekranı” ile başlamamalıdır. Önce ortak çekirdek kurulmalıdır:

1. `transport_core` foundation.
2. Tenant/site bootstrap: `gursel-turizm`, mevcut `aractakip`, yeni `veli`.
3. Institution, transported person, group, route temel modelleri.
4. Permission/menu/translation seed temeli.
5. Admin dashboard ve scope selector.
6. Event key registry.
7. Offline/access/wallet/portal modülleri için ADR.

Bundan sonra iki paralel hat açılabilir:

- **Hat A:** KT kart/QR geçiş + wallet + trip ledger.
- **Hat B:** Veli/personel borç + sözleşme + ödeme portalı.

Bu ayrım hem ürün hızını artırır hem de `Person`, `Institution`, `Payment`, `Ledger`, `Integrations` gibi ortak alanlarda mükerrerliği engeller.

---

## 12. Açık Kararlar

| Karar | Seçenekler | Öneri |
|---|---|---|
| Mobil uygulamalar | Tek app / ayrı driver-reader app / PWA | Driver/reader offline ihtiyacı nedeniyle ayrı native-capable uygulama veya offline-first PWA değerlendirilmeli |
| Mifare okuyucu | Telefon NFC / harici reader / validator cihazı | Device abstraction; donanıma kilitlenme yok |
| Offline policy | Limitsiz offline / kontrollü offline | Kontrollü offline, signed log ve risk limitleri |
| Ring inference | Tam otomatik / öneri + onay | İlk sürüm öneri + operatör onayı |
| Wallet | Kart bazlı / kullanıcı bazlı | Kullanıcı bazlı; kart/QR sadece instrument |
| Campus services | Ulaşım içine gömülü / genel service category | Genel `service_category` |
| SAP | Domain içine gömülü / adapter-first | Canonical event + adapter-first |

---

## 13. Definition of Done Ekleri

Bu genişletilmiş kapsam için her iş kalemi aşağıdakileri de kontrol etmelidir:

- Offline etkisi var mı?
- Device/reader etkisi var mı?
- Wallet/ledger etkisi var mı?
- Fatura sorumluluğu etkisi var mı?
- Tenant/site/institution/person/group scope doğru mu?
- Geçiş/ödeme olayı audit ve replay edilebilir mi?
- GPS/telematics ile ilişki varsa veri sahipliği net mi?
- Hassas veri veya konum verisi için KVKK etkisi dokümante edildi mi?
- API sharing ile dış kurumlara açılacak veri minimum ve onaylı mı?

---

## 14. Sonuç

Yeni kullanıcı metni, önceki veli ödeme raporunu iptal etmiyor; onu daha büyük bir ulaşım ve kampüs değer platformunun parçası haline getiriyor.

Doğru hedef:

- Veli ödeme portalını yapmak,
- ama onu tek başına bir ödeme uygulaması gibi değil,
- KT kart/QR geçiş, cüzdan, çetele, GPS doğrulama ve kurum/araç sahibi faturalama omurgasıyla aynı `TransportedPerson + Institution + Wallet + Ledger + Payment + Integration` modeline oturtmaktır.

Bu yaklaşım AYSTEK Core’un nesne tabanlı, multi-tenant, multi-site, event-driven ve enterprise platform karakterine en uygun yoldur.

