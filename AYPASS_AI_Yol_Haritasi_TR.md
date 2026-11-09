# AYPASS — Mevcut Uygulamayı Tamamlama ve AI Prompt Paketi

**Hazırlanma:** 29 Ekim 2026  
**Sınıflandırma:** AYSTEK özel proje çalışması — yalnızca yetkilendirilmiş geliştirme ortamında kullanım.  
**Amaç:** AYShuttleGursel davranışlarını ve Rev.06 hedeflerini, hâlihazırda geliştirilmiş AYSTEK Core + Next.js AYPASS uygulaması üzerinde kanıtla karşılaştırmak; eksikleri kontrollü iş paketleriyle tamamlamak. Bu paket yeni bir ürün veya ikinci mimari oluşturma talimatı değildir.

## Kullanım

1. Önce `01_KANIT_VE_BOSLUK_ANALIZI_TR.md` ve `02_YOL_HARITASI_TR.md` dosyalarını oku.
2. Seçtiğin yetkilendirilmiş kodlama aracına `prompts/00_ANA_PROMPT_TR.md` dosyasını ver. İlk oturum **P00 keşif** ile sınırlıdır; üretim kodu değiştirilmez.
3. P00 çıktısındaki gerçek eksiklerden seçilmiş bir iş kaydını, ilgili faz promptuyla uygulat. Bütün fazları tek komutta otomatik yürütme.
4. Her iş için `prompts/90_BAGIMSIZ_DENETIM_TR.md` ile ayrı bir inceleme yaptır. İnceleme aracı yazarın başarı beyanını değil diff, aktif çağrı zinciri ve test kanıtını doğrulasın.
5. Yeni oturumda `prompts/91_DEVAM_TR.md` kullan. Planı baştan yazdırma; son kanıt ve iş kaydından devam et.

Prompts klasöründeki her faz, ana promptun sınırlarıyla birlikte kullanılır. Ayrı bir araç oturumunda ana promptu veya projede erişilebilir kopyasını yeniden yükle. Tek dosya isteyen araçlar için tüm promptların birleşimi ayrıca sunulmuştur.

## İncelenen kaynaklar ve sınırlar

Dört yüklenmiş Markdown belgesi, iki özel deponun aşağıdaki sabit GitHub sürümlerindeki seçilmiş uygulama dosyaları, DTO/API istemcileri ve bağlantı noktaları incelendi:

| Kaynak | Referans |
|---|---|
| Backend | `huseyinkoc-aystek/aystek-core` — `f48748fd672a98d838ce60c98a533e0ed9b81ddb` |
| Panel | `huseyinkoc-aystek/admin-panel-next` — `56907ebc6059566690e4b0d4b6e21b63d4db5163` |
| Legacy analiz | `AYSHUTTLEGURSEL_PROJE_ANALIZI_TR.md`, 29.10.2026 |
| Legacy teknik envanter | `AYSHUTTLEGURSEL_TEKNIK_ENVANTER_TR.md`, 29.10.2026 |
| Tarihsel analiz | `AYSHUTTLEGURSEL_CURRENT_STATE.md`, 09.07.2026; yeni analizdeki düzeltmelerle birlikte |
| Mühendislik kuralları | `README_FOR_AI_ASSISTANTS(20261029-120719).md` |
| Hedef ürün kuralları | Bu görüşmede oluşturulan Rev.06 kapsamı ve kullanıcının ulaşım, offline, POS, market ve hizmet devri talepleri |

Bu inceleme bütün depoların satır satır denetimi değildir. GitHub kaynak okuması; test çalıştırması, cihaz testi, veritabanı kontrolü, banka testi veya üretim kabulü değildir. Bu oturumda uygulama koduna, veritabanına, dağıtıma veya GitHub kayıtlarına değişiklik yapılmadı.

Kullanıcının verdiği alan adı **`aypasss.aysteknolojileri.com`** biçimindedir. Alan adı üzerinden bu oturumda içerik alınamadı; bu durum sitenin kapalı olduğunu kanıtlamaz. Dağıtılan sürüm ile GitHub `main` veya yerel çalışma ağacının aynı olduğu varsayılmadı. Önce deployment kimliği eşleştirilecektir.

İlişkili master plan, kanonik mimari ve Rev.06 yerel iş planı dosyalarına kaynaklarda referans bulundu; bu dosyaların tam içerikleri bu incelemede doğrulanmadı. P00 bunları yerel workspace'te okumalı ve bu paketi mevcut iş takibinin devamı olarak işlemelidir.

## Kanıt etiketleri

- **KODDA GÖRÜLDÜ:** Belirtilen sabit sürümde incelenen dosya/işlev bunu gösteriyor; çalışma zamanı başarısı değildir.
- **BELGEDE BİLDİRİLDİ:** Kaynak dokümanın tespiti veya kullanıcının ürün gereksinimi.
- **HEDEFLE FARK:** İncelenen uygulama davranışı ile açık ürün gereksinimi farklıdır.
- **RİSK / TEST GEREKLİ:** Koddan türetilen senaryo; çalışma zamanında yeniden üretilmedi.
- **DOĞRULANMADI:** İlgili kaynak veya uçtan uca yol incelenmedi. “Eksik” ya da “yok” olarak sunulamaz.

## Gizlilik ve yetki

Özel kaynakları, kullanıcı/kurum kayıtlarını, QR tokenlarını, ödeme verilerini ve anahtarları arama motorlarına, herkese açık issue/gist/depolara veya onaysız servis/AI araçlarına gönderme. Bu paket özel kaynak kodlarının tam kopyalarını veya erişim anahtarlarını içermez. Bir AI aracına prompt vermek, o aracın kaynakların tamamına veya canlı sisteme erişmesine otomatik onay değildir.

GitHub push/PR/issue oluşturma, canlı sunucuya bağlantı, gerçek tahsilat/iade, e-fatura gönderimi, üretim migration/seed ve deploy ayrıca açık yetki gerektirir. Promptlar bu yetkiyi kendiliğinden vermez.

## Paket düzeni

- `01_KANIT_VE_BOSLUK_ANALIZI_TR.md`: Doğrulanan noktalar, farklar ve belirsizlikler.
- `02_YOL_HARITASI_TR.md`: Bağımlılık sırası ve her fazın çıkış kapısı.
- `03_GEREKSINIM_ESLEME_TR.md`: Ürün gereksinimleri → mevcut sahipler → fazlar.
- `04_KABUL_TEST_MATRISI_TR.md`: Henüz çalıştırılmamış kabul senaryoları.
- `05_KAYNAK_MANIFESTOSU.json`: Dosya kimlikleri, hashler ve kod okuma kapsamı.
- `prompts/`: Ana, faz, bağımsız denetim ve devam promptları.
- `templates/`: İş kaydı, karar ve oturum kapanış şablonları.

**İlk uygulanacak çalışma:** P00 ile sürüm/kanıt eşlemesi; ardından P01 güvenlik ve sözleşme tabanı, P02 mali tutarlılık. Yeni özellik sayısını artırmak, bu kapıları geçmenin alternatifi değildir.


---

# AYPASS — Kanıt ve Boşluk Analizi

29.10.2026 • Statik, seçilmiş kaynak incelemesi • Canlı kabul raporu değildir.

## 1. Kaynakların ne söylediği

29 Ekim legacy analizi açıkça eski Yii2/PHP uygulamasını inceler; AYSTEK Core'un güncel kodunu o analiz kapsamında denetlemediğini belirtir. Bu nedenle AYShuttleGursel'deki bir açık, yeni AYPASS uygulamasına doğrudan atfedilemez. Temmuz tarihli `CURRENT_STATE` tarihsel belgedir; yeni analizde PHP bağımlılıkları, PDF kütüphanesi, salt-eklemeli defter, token süresi, SUL alanları ve scheduler konusunda düzeltmeler vardır. [D1 §1, §10, §12, §14; D2 başlık notu; D3 §8]

Legacy analizinin kaynak referansı `0956c9be6bea96a849e198d9785b35c4f7590608` + dört yerel değişikliktir. Bu, yeni Go deposunun incelenen SHA'sı değildir. Uygulama sayımları canlı API/DB sayısı veya işlevsel başarı kanıtı değildir.

## 2. Mevcut uygulamada görülenler ve öncelikli farklar

| ID | Kanıt / gözlem | Sonuç ve sınır | İş paketi |
|---|---|---|---|
| E01 | `access_points`, `wallets`, `ledger`, `route_fares`, `route_stop_fares`, `rider_trips` birlikte kullanılıyor; `cmd/api/main.go` bağlantı kesiti bunları aktif servise bağlıyor. [R01, R02, R05, R06] | **KODDA GÖRÜLDÜ.** Boş proje değildir. Aynı sahipleri yeniden yazmak yerine mevcut akışlar genişletilmeli. | P00 / tüm fazlar |
| E02 | `chargeMetered` ilk okutmayı ücret almadan kaydediyor, yolculuk açıyor; sonraki okutma çıkış ücreti üretiyor. [R02:1330–1625] | **HEDEFLE FARK.** Bu yol A3/A4'ün fiili maksimum tahsilat + fark iadesi akışı değildir. Diğer olası uygulama yolları ayrıca aranmalı; eski metered davranış onaysız değiştirilmemeli. | P04 |
| E03 | Aynı `chargeMetered` çıkışta OD fiyatını arıyor; çözülemezse kilometre tarifesine dönüyor. [R02:1471–1625] | **HEDEFLE FARK.** A4 için tanımsız matris hücresinin mesafeye sessiz dönüşü kabul edilmemeli. Mevcut karma strateji ile hedef A4 ayrı politikalar olabilir. | P04 |
| E04 | Fiyat çözümlemesi `time.Now()` kullanıyor; çıkış OD fiyatı da o anki saati sorguluyor. [R02:1000–1160,1471–1625] | **RİSK / TEST GEREKLİ.** Yolculukta sabitlenmiş tarife ve gecikmiş olay zamanının korunması bu kesitte gösterilmedi. Başka bir snapshot mekanizması varsa bağlanıp kanıtlanmalı. | P03–P05 |
| E05 | Transport modelindeki açıklama ve `resolveEventContext`/test kesitleri plaka ve aktif araç seferi zorunluluğunu gösteriyor. [R01, R03] | **HEDEFLE FARK.** Doğru sefer bağlamı korunarak ücret kararını plakanın varlığına gereksiz bağımlı olmaktan çıkarma hedefi tekrar ele alınmalı. Sefer kimliğini tamamen kaldırmak çözüm değildir. | P03 |
| E06 | `chargeMetered` içindeki `OpenTrip` ve `CloseTrip` hataları `_ =` ile göz ardı ediliyor. [R02:1471–1625] | **KODDA GÖRÜLDÜ + RİSK.** Olay var/yolculuk yok veya ödeme var/kapanış yok senaryoları hata enjeksiyonuyla denenmeli; veri kaybı yaşandığı iddia edilmiyor. | P02, P04 |
| E07 | Wallet `PersonID` ve ledger hesabı ile açılıyor; repository tenant/site filtreli. [R05:1–345; R06:1–170] | **KODDA GÖRÜLDÜ.** Kişi hesabı yaklaşımı zaten var. “Karttan kişiye yeni sistem kur” yerine kart bağları, tarihsel eşleme ve kapsam doğrulanmalı. | P02 |
| E08 | Debit önce `CachedBalance` kontrolü yapıyor, journal oluşturuyor, ayrı cache güncellemesi kullanıyor; `UpdateCachedBalance` hesaplanmış değeri `$set` ediyor. [R05:180–345; R06:1–170; R07] | **RİSK / TEST GEREKLİ, yüksek öncelik.** Farklı idempotency anahtarlı eşzamanlı işlemler ve journal sonrası kesinti incelenmeli. Çağıran katmanda tam transaction/seri hale getirme olup olmadığı tamamlanmadan “atomiklik yok” diye kesin hüküm verilmemeli. | P02 |
| E09 | Credit/Debit tekrarında mevcut idempotency journal'ı erken döndürülebiliyor. [R05:180–345] | **RİSK / TEST GEREKLİ.** Aynı anahtar + farklı wallet/tutar/döviz, payload doğrulaması ve kısmi başarı sonrası tamamlama davranışı test edilmeli. Bir anahtarın varlığı uçtan uca exactly-once garantisi değildir. | P02 / P06 |
| E10 | İncelenen ücret/cüzdan/yanıt modellerinde `float64` tutarlar var. [R01, R02, R05, R08] | **KODDA GÖRÜLDÜ.** Mali hassasiyet/yuvarlama sözleşmesi gerekli; kendiliğinden para kaybı yaşandığı kanıtlanmadı. Mevcut veriyi topluca ve onaysız dönüştürme. | P02 |
| E11 | Wallet kurum varsayılanı `affiliations[0]` üzerinden seçiliyor. [R05:1–160] | **KODDA GÖRÜLDÜ + RİSK.** Çok kurumlu kişi için limit, fon sahibi ve program seçimi açık sözleşmeye bağlanmalı. Tek kurum senaryosundaki davranış bozulmamalı. | P02 / P08 / P10 |
| E12 | Backend AccessPointResponse'da `vehicle_plate` ve `pos_enabled` var; panel `parseAccessPoint` bunları döndürmüyor. [R08:1–115; R09:1–220] | **DOĞRULANAN SÖZLEŞME FARKI.** Bu ayrıştırıcıdan geçen alanlar kayboluyor. Etkilenen form/ekran kapsamı ve aynı kaynakta farklı istemci olup olmadığı P01'de test edilmeli. | P01 |
| E13 | Erişim noktası sayfasında `Ad ara...`, `Tüm Türler`, `Transport Platform / Access Points` gibi doğrudan görünen metinler var. [R10:1–140] | **DOĞRULANAN STANDART BORCU.** README'nin çeviri anahtarı kuralıyla uyumsuz kısımlar var. Görsel standardı tüm panelde keyfî yeniden tasarlama gerekçesi değildir. | P01 / her UI dilimi |
| E14 | `access_points/routes.go` eski cihaz-erteleme notu içerirken `cmd/api/main.go` cihaz validate/events/me yollarını bağlıyor. [R04, R11] | **DOKÜMAN/KOD FARKI.** “Cihaz API'si yok” denemez. Gateway middleware, credential bağları ve yetki zinciri gerçek bağlantıdan izlenmeli. | P01 / P05 |
| E15 | Panel package.json `test: tsx --test tests/*.test.ts`, check-env, check-i18n içeriyor; Playwright bağımlılığı mevcut. [R12] | **KODDA GÖRÜLDÜ.** README'nin test altyapısı yok ifadesi güncel manifest ile aynı değil. Gerçek test sayısı/kapsamı veya başarılı E2E olduğu sonucu çıkmaz. | P00 / tüm fazlar |
| E16 | Mevcut erişim API'sinde ve panel istemcisinde kasa/tür, güzergâh ve kişi rapor yolları var; Wallet Reconcile/ReconcileAll var. [R04, R09, R05:700–850] | **KODDA GÖRÜLDÜ.** Raporları sıfırdan çoğaltma. ReconcileAll'ın ilk 10.000 sınırı/hata atlama davranışı, kapsam ve yeniden yürütme ayrı doğrulanmalı. | P09 |
| E17 | `SupportsSaleCatalog` yalnızca canteen/photocopy için true dönüyor; diğer commerce/market yolları bu incelemede tam taranmadı. [R01] | **DOĞRULANMADI.** Bu satış yolu tek başına Model4'ün üye işyeri, hakediş ve fon ayrımı modelini kanıtlamaz. Model4 bütünüyle yok sonucu da çıkarılamaz. | P10 |
| E18 | Rev.06 hizmet devri talepleri var; mevcut contract/devir uygulamasının tam çağrı zinciri bu incelemede doğrulanmadı. [U1] | **DOĞRULANMADI.** Önce mevcut sahipleri ve iş planındaki bitmiş adımları bul; yeni paralel devir modülü oluşturma. | P08 |
| E19 | Tarife/erişim DTO'sunun incelenen olay girişinde özgün olay zamanı/sürüm alanı görünmüyor; gateway ve cihaz kuyrukları tam incelenmedi. [R08] | **DOĞRULANMADI / RİSK.** Bu, offline sistem hiç yok demek değildir. Sunucu + gerçek Ubuntu/Android istemci protokolü birlikte kanıtlanmalı. | P05 |
| E20 | Repo kök ağacında `.env` ve `.env.production` isimleri görülüyor; içerikleri okunmadı. [R13] | **İNCELEME GEREKLİ.** Gerçek secret içerdiği iddia edilmiyor. Yerel maskeli tarama; doğrulanırsa yetkili rotasyon ve iz temizliği planı. | P01 |

## 3. Yanlış sonuç çıkarılmaması gereken noktalar

`metered` açıklamalarının bir kısmı eksik tarife durumunda ücretsiz çıkıştan söz ediyor; aynı incelenen dosyanın güncel `resolveMeteredRate` gövdesi ise eksik tarifede `ErrActiveRouteFareRequired` döndürüyor. Bu nedenle yorumdan hareketle “eksik tarifeyi ücretsiz geçiriyor” bulgusu yazılmadı. [R02:1000–1160]

`PosEnabled` bir yapılandırma alanıdır. Bankanın terminali onayladığını, POS cihazına gerçekten komut gönderildiğini, callback doğrulamasını veya tahsilatın kesinleştiğini kanıtlamaz. Benzer şekilde çift taraflı journal ve idempotency alanlarının bulunması bütün para akışının atomik/tekrar güvenli olduğunu tek başına kanıtlamaz.

`site_bridge` arama sonuçlarında belirli kapsamlı veri paylaşımı görülmesi, bütün kurumların veya cüzdanların ortaklaştırılmasına izin vermez. Tenant, site, institution, service-provider ve merchant birbirinin eş anlamlısı değildir.

Aramada görünmeyen bir dosya için “uygulamada yok” denmez. Kodda görülen test dosyası için “test geçti” denmez. GitHub main için “canlı sürüm” denmez.

## 4. Kaynak rehberi

- **D1:** `AYSHUTTLEGURSEL_PROJE_ANALIZI_TR.md`, özellikle §1, §6, §10, §12–14.
- **D2:** `AYSHUTTLEGURSEL_CURRENT_STATE.md`; tarihsel, D1 düzeltmeleriyle.
- **D3:** `AYSHUTTLEGURSEL_TEKNIK_ENVANTER_TR.md`; rota/koleksiyon keşif haritası, fiziksel DB doğrulaması değil.
- **D4:** `README_FOR_AI_ASSISTANTS(20261029-120719).md`; kanonik referans, DDD, tenant/site, 28 dil, menu/permission/view-profile, UI, test ve operasyon kuralları.
- **U1:** Kullanıcının bu görüşmedeki ürün kuralları ve Rev.06 hedef kapsamı; yazılımın hazır olduğunun kanıtı değildir.

Aşağıdaki R kaynakları `05_KAYNAK_MANIFESTOSU.json` içinde depo, SHA ve okunan aralıklarla kayıtlıdır:
R01 access_point modeli; R02 access_points service; R03 resolveEventContext ve testi; R04 erişim route'ları; R05 wallets service; R06 wallets repository; R07 applyBalanceChange kesiti; R08 access_points DTO; R09 panel API istemcisi; R10 panel sayfası; R11 aktif main wiring; R12 panel package.json; R13 backend root tree; R14 route_fares/route_stop_fares arama kesitleri.

**Sonuç:** Öncelik yeni ekran sayısı değil; mevcut uygulamanın mali ve operasyonel davranışını gereksinimlere bağlayan kanıt zinciridir. Bu tespitler kontrollü testlere ve güncel workspace doğrulamasına başlangıç sağlar; üretim arızası yaşandığına ilişkin iddia içermez.


---

# AYPASS — Aşamalı Tamamlama Yol Haritası

**Yaklaşım:** Keşfet → mevcut sahibi doğrula → farkı kanıtla → küçük uygulama dilimi → bağımsız test → mevcut iş planında kapanış.  
**Başlangıç:** Yeni uygulama zaten mevcut; legacy PHP kodu referanstır, yeni mimarinin şablonu değildir.  
**Durum:** Bu belgedeki bütün fazlar plan/tasarım önerisidir; tamamlandı veya test edildi işareti taşımaz.

## 1. Bağımlılık planı

```text
P00 Kanıt/sürüm eşlemesi
  └─ P01 Güvenlik + API/panel sözleşmesi + test tabanı
       └─ P02 Kişi hesabı + mali doğruluk
            ├─ P03 Sefer/araç/firma bağlamı ─ P04 Tarife algoritmaları ─ P05 Offline
            └─ P06 Ödeme/POS/yükleme/fatura
                    └─ P07 Legacy işlev eşdeğerliği ve entegrasyonlar
P03+P04+P05+P06 ─ P08 Hizmet sağlayıcı devri
Önceki fazların kanıtları ─ P09 Raporlama, ölçek ve işletim
P02+P06+P08+P09 + dış onay ─ P10 Global Market (ayrı açılış kapısı)
P00–P09 ─ P11 Ulaşım geçiş provası ve pilot
P10 dahilse ayrıca market kabulü ─ P11 Market açılışı
```

P10'un ticari/finansal ortağı hazır değilken ulaşım pilotu zorunlu olarak beklemez; market özelliği kapalı ve açıkça kapsam dışı tutulabilir. Buna karşılık mali tutarlılık/güvenlik kapıları atlanamaz.

## 2. Fazlar ve çıkış kapıları

| Faz | Hedef | Çıkış kanıtı |
|---|---|---|
| P00 | Workspace, SHA, yerel diff, dağıtım kimliği, kanonik mimari ve mevcut master planı eşleştir. | Kaynak manifesti, gereksinim→backend→panel→test matrisi, doğrulanmamışların listesi. Üretim kodu değişmeden. |
| P01 | Yetki/scope/device zinciri, güvenli test tabanı, DTO-panel uyuşmazlıkları. | Çapraz tenant/site/kurum negatif testleri; cihaz binding testleri; vehicle_plate/pos_enabled sözleşme testi; gerçek test komutları. |
| P02 | Kişi hesabı, kart bağları, kesin para aritmetiği, journal/bakiye/idempotency bütünlüğü. | Eşzamanlı farklı işlemler, kayıp yanıt, journal sonrası kesinti, aynı anahtar/farklı payload, kart değişimi testleri; mutabakat eşitliği. |
| P03 | Ücret bağlamını plaka atamasına gereksiz bağımlılıktan ayır; sefer/sözleşme sürümlerini koru. | Plakasız ama yetkili/kanıtlı sefer geçerli; bağlamsız ücret yok; araç değişimi finansal kaydı değiştirmiyor. |
| P04 | F0/ücretsiz mevcut davranışı koruyup A1–A4'ü güzergâh ayarından kesinleştir. | A1/A2 tek ilgili okutma; A3/A4 maksimum fiili debit ve tam bir kez iade; matris eksikse mesafe fallback yok; sürüm testleri. |
| P05 | Sunucu ve gerçek cihaz için dayanıklı offline protokol. | Kalıcı kuyruk, aynı olayda tek mali etki, ters sıra, eski kimlik/tarife/firma, çoklu cihaz riski ve güç kesintisi kabul kanıtları. |
| P06 | Banka/POS, yükleme, ledger, iade ve belge durumlarını ayrıştır. | Kesin ödeme onayı, tek yükleme, outbox/worker tekrar güvenliği, belirsiz sonuç sorgulaması; sandbox kanıtı ayrı. |
| P07 | Legacy'den gerekli işlevleri kayıpsız tamamla; kusurlu davranışı taşımadan. | Rezervasyon/ek yolcu/yedek liste, abonelik, kampanya, OTP/QR bilet, kurum/SAP/SMS/GPS eşdeğerlik matrisi. |
| P08 | Sabancı ulaşımında Gürsel→yeni firma devrini yönet. | Kesim provası, açık yolculuk/iadeler/offline olay sahipliği, yetki daraltma, onaylı bakiye politikası; operasyon/mali kapanış ayrımı. |
| P09 | Finansal ve operasyonel raporların doğruluğunu, ölçeğini ve izlenebilirliğini doğrula. | Kaynağa inilebilir rapor, gecikmiş olay düzeltmesi, kişi/geçiş ayrımı, mutabakat, ölçülen performans ve hata alarmları. |
| P10 | Mevcut yeteneklerden Model4 üye işyeri alışveriş ağını tamamla. | Program/fon/merchant ayrımı, kısmi iade, hakediş/komisyon, dış sağlayıcı ve onay kapısı. Ulaşım offline borcu otomatik taşınmaz. |
| P11 | Eski-yeni veri ve cihaz geçişini prova et, kontrollü pilot ve sürüm planı hazırla. | Dry-run toplamları, açık işlemler/cihaz kuyrukları, tek mali yazar, geri dönüş sınırları, pilot mutabakatı, ayrı yayın onayı. |

## 3. İlk iş sırası

P00 tamamlandığında ilk aday işler, güncel workspace'te hâlâ mevcut oldukları kanıtlanırsa şöyledir:

1. Backend AccessPoint DTO ↔ panel parser/type/form sözleşme farkını bir testle görünür kıl; en küçük düzeltme ile kapat. İşlem sırasında yeni yetki açma veya canlı seed yok.
2. Wallet debit/cache/journal için ayrı işlem anahtarlarıyla eşzamanlılık ve adım arası kesinti testleri kur. Mali risk doğrulanırsa yeni ücret algoritmalarından önce çöz.
3. Aktif transport bağlamının plaka bağımlılığını ve gerçek sefer kimliği ihtiyacını kararlaştır; person/kurum/provider ayrımını belgeye bağla.
4. `metered` ile A1–A4 farkını test tablosuna dönüştür. Eski davranışa bağlı istemciler için kontrollü sürüm/özellik geçişi tasarla.

Bu sıra, bütün fazı tek PR'a sığdırma önerisi değildir. Her adım küçük, incelemeye elverişli bir dikey dilim olmalıdır.

## 4. Dikey dilim kuralı

Her iş yalnız backend koduyla tamamlanmış sayılmaz. Etkileniyorsa **model + service/repository + handler/route + API istemcisi/type/form + izin + menu + view profile + çeviri + indeks/seed + test + doküman** aynı işin kapsamındadır. İlgisiz panel sayfalarını veya haber/AI/kişisel asistan alanlarını değiştirme.

README'nin DDD yapısını kullan. Mevcut doğru domain sahibini genişlet. Yeni bounded context gerekiyorsa önce mevcut modülle neden çözülemediğini, bağımlılık yönünü ve veri sahibini kısa mimari karar kaydıyla açıkla. Sırf “enterprise” görünmek için servis/DB/queue çoğaltma.

## 5. Çoklu AI araçlarıyla çalışma

**Analiz sahibi:** Kaynakları/sürümü eşler, gereksinim ve karar defterini günceller. Kod değiştirmeden kanıt üretir.

**Uygulayıcı:** Tek seçilmiş iş kaydı üzerinde mevcut modülü değiştirir. Ayrı branch/worktree kullanır; başka aracın veya kullanıcının yerel değişikliklerini geri almaz. Etkilenen dosyaları iş başında kaydeder.

**Bağımsız denetçi:** Uygulayıcının özetinden bağımsız diff'i, gerçek wiring'i ve test komutlarını okur. Mock'taki başarı ile gerçek repository/DB davranışını ayırır. Doğrulanmamış maddeleri “geçti” yapmaz.

Aynı `wallets`, `access_points/service.go`, `cmd/api/main.go`, `translation_seed` veya ortak tip dosyaları üzerinde eşzamanlı yazıcı çalıştırma. Paralellik yalnız bağımsız dosya sahipliği ve üzerinde uzlaşılmış API sözleşmesi varsa kullanılmalı. Çakışan işler sıraya alınır; önce sözleşme, sonra backend/frontend.

## 6. Done tanımı

Bir işin tamamlanması için gereksinim ve kaynak kimliği; aktif çağrı zinciri; ilgili pozitif/negatif/tekrar/eşzamanlılık testleri; frontend API uyumu; tenant/site/kurum sınırı; menü/izin/view-profile/28 dil; mali işlemler için mutabakat; geçici uyumluluk ve kaldırma planı; değişen operasyon dokümanları birlikte görünür olmalıdır.

“Test dosyası eklendi”, “derlendi” veya “ekranda var” tek başına kabul değildir. Gerekli dış ortam yoksa durum **engelli/doğrulanmadı** kalır; yerel kanıtlar korunur.

## 7. Açık kararlar — otomatik varsayılmayacak

- A1/A3 mesafesinin gerçek kaynağı; güzergâh üzerinde tanımlı mesafe, olay konumu veya başka onaylı yöntem. İki GPS noktası arasındaki uzaklık otomatik olarak gerçek yol mesafesi sayılmaz.
- Borç sınırları, offline süre/tutar/adet ve cache yaşı; −15 ve 120 TL örnekleri evrensel varsayılan değildir.
- Hangi cüzdan/fon sağlayıcısı hangi hizmette kullanılabilir; kişi kimliği ortak olsa da fonlar kendiliğinden ortak değildir.
- Eski `metered` müşterilerinin yeni algoritmaya geçiş koşulları; onaylı mevcut davranış sessiz değiştirilmez.
- İkinci okutma eksikliğinde kapanış/iade kuralı; otomatik ceza/süre uydurulmaz.
- POS bankası/modeli/entegrasyonu, gerçek üye işyeri ve alıcı hesap; boolean flag yeterli değildir.
- Fatura profilinin hizmet/yükleme/market bazında doğrulanması; hukuki sonuç varsayılmaz.
- Devirde fonun, borcun, abonmanın, promosyonun sorumlusu ve onayları.
- Rapor ortalamasının pay/payda tanımı, hedef veri hacmi, performans/operasyon eşikleri.
- Dağıtılan kaynak SHA'sı, gerçek cihaz deposu ve rollout yetkisi.

Eksik karar tüm çalışmayı durdurmaz; etkilediği riskli değişiklik bloke edilir, diğer kanıt ve güvenli yerel işler ilerletilir.


---

# Gereksinim → Mevcut Sahip → İş Paketi Eşlemesi

**Durum notu:** Bu tablo hedef kapsamı belirtir. Sahip sütunu, mevcut kodda görülen alan veya P00'da doğrulanacak adaydır; özelliklerin bütünüyle tamamlandığı anlamına gelmez. REQ kodları bu prompt paketinin izleme kodlarıdır; mevcut master plandaki kimliklerin yerine geçmez.

| Kod | Korunacak / tamamlanacak gereksinim | Kaynak | Mevcut sahip veya keşif adayı | Faz |
|---|---|---|---|---|
| REQ-01 | Model 1/2/3 kurum, operatör ve hizmet sağlayıcı ayrımı | U1 Rev.06 | institutions/persons + sözleşme/provider sahipleri doğrulanacak | P00,P03,P08 |
| REQ-02 | Person hesabı; kart değişiminde bakiye/iade devri gerekmemesi | U1 Z03, D1 §5–6 | persons/wallets/wallet_instruments mevcut | P02 |
| REQ-03 | Kayıp/iptal kart ve geçmiş sahiplik; gecikmiş olay doğru hesaba | U1 Z03/Z04 | instruments + event-time ilişki sözleşmesi | P02,P05 |
| REQ-04 | Tenant/site ve kurum/provider yetkileri ayrı | D4, U1 | tenantfilter/contextkeys + module authorization | P01,tümü |
| REQ-05 | QR/kartın banka ödeme yetkisinden ayrılması; karta yazmama | U1 | wallet instruments/device auth/payment adapter | P01,P05,P06 |
| REQ-06 | Sabit ID; ad/plaka değişiminde geçmişin korunması | U1 Z02/Z13,D1 §5 | fleet_tally/persons/institutions ve historical snapshots | P03,P11 |
| REQ-07 | Araç ataması olmadan doğrulanmış sefer/tarife; belirsiz fiyata izin yok | U1 Z01 | access_points + trip_assignments; fark görüldü | P03 |
| REQ-08 | F0 sabit, indirimli ve ücretsiz takip korunmalı | U1 | access_points/discounts mevcut | P04 |
| REQ-09 | A1 bilinen kampüs/fabrika; sabah biniş, akşam iniş mesafe hesabı | U1 A1 | route_fares/access_points genişletme analizi | P04 |
| REQ-10 | A2 güzergâh/yön/durak sabit ücret; OD matrisiyle karıştırma | U1 A2 | route/stop fare sahipleri; strateji farkı doğrulanacak | P04 |
| REQ-11 | A3 fiili maksimum debit, mesafeye göre net ücret ve fark iadesi | U1 A3 | metered mevcut, hedefle farklı | P04 |
| REQ-12 | A4 fiili maksimum debit, yönlü sabit matris ve fark iadesi | U1 A4 | route_stop_fares mevcut, metered akışı farklı | P04 |
| REQ-13 | Algoritma, indirim, risk ve tahsilat ayarları güzergâh bazlı | U1 | route/access_point/institution config sınırları | P03,P04 |
| REQ-14 | Tarihli tarife/indirim sürümü; eski olaya yeni fiyat yok | U1 Z06,Z08 | route_fares + applied context; now kullanımı görüldü | P03,P04,P05 |
| REQ-15 | Online alt bakiye ve offline risk/borç ayrımı | U1 Z07 | wallets floor mevcut; offline zincir doğrulanacak | P02,P05 |
| REQ-16 | Kalıcı cihaz kuyruğu, ACK, duplicate/out-of-order güvenliği | U1 Z05,Z06 | API gateway + gerçek Ubuntu/Android kaynakları gerekli | P05 |
| REQ-17 | Tek mali kaynak; posted hareket değişmez, reversal ile düzeltme | U1 Z09 | ledger/wallets mevcut, kesinti testleri gerekli | P02 |
| REQ-18 | Eşzamanlı bakiye ve idempotency payload bütünlüğü | U1 Z06,Z09,Z10 | Credit/Debit/PostJournal/cache zinciri | P02 |
| REQ-19 | Yükleme ödeme onayı sonrası tek mali etki | D1 §6.3,U1 Z10 | payments + wallets + ledger | P06 |
| REQ-20 | Belge alındı ≠ faturalandı; unknown sonuç mutabakatı | U1 Z11 | billing/integrations/outbox sahipleri doğrulanacak | P06 |
| REQ-21 | Raspberry Pi iş kuralları; onaylı merchant POS ayrı | Kullanıcının POS senaryosu | device gateway/payment adapters/merchant config | P06 |
| REQ-22 | Koltuk/ek yolcu/yedek liste/ön ödeme/iptal | D1 §6.4 | fleet/reservation sahipleri keşfedilecek | P07 |
| REQ-23 | Abonelik/dönem hakkı ve indirim kullanım limitleri | D1 §5–6 | subscriptions/access_point_discounts | P07 |
| REQ-24 | Kampanyalı yükleme, nakit ve promosyon ayrımı | D1 §6.3,§10 | wallet funding + campaign sahipleri | P02,P07 |
| REQ-25 | Misafir OTP ve satın alınan QR hakkı ayrı; tek tüketim | D1 §6.5,§10 | auth/payment/entitlement gerçek sahipleri keşfedilecek | P07 |
| REQ-26 | Kurum/SAP/SMS/GPS entegrasyonu ve izlenebilir tekrar | D1 §7–8,D4 | integrations/telematics ve onaylı provider adaptörleri | P07 |
| REQ-27 | Mercedes/LCW tanımlı fabrika QR/kart noktaları; mevcut ücretsiz kapsam | U1 | factory_gate AccessPoint mevcut | P03,P09 |
| REQ-28 | Kasa, rota ortalaması, araç/personel türü raporları | U1 Z12 | access_points raporları/fleet_tally mevcut | P09 |
| REQ-29 | Model 1 sağlayıcı devir kapsamı, kesim, eski/yeni ilişki | U1 Rev.06 devir | contract/assignment sahibi doğrulanacak | P08 |
| REQ-30 | Devirde bakiye/borç/hak ve ödeme profili; operasyon≠mali kapanış | U1 Rev.06 devir | wallets/ledger/payment profile + handover akışı | P08 |
| REQ-31 | Üye işyeri/şube/terminal/sponsor/program/fon ayrımı | U1 Model4 | mevcut commerce/payment/contract sahipleri keşfedilecek | P10 |
| REQ-32 | Market iade/komisyon/hakediş/settlement ve dış onay | U1 Model4 | ledger/payments/billing/integrations | P10 |
| REQ-33 | Backend DTO ↔ panel istemcisi/type/form uyumu | D4,R08–R10 | access-points DTO/parser/page; somut fark var | P01 |
| REQ-34 | 28 dil, menu, permission ve view-profile, event documentation | D4 | seeds/FE shared components/EventKeysInfo | her faz |
| REQ-35 | Tekrar çalıştırılabilir migrasyon; eski kayıt çift tahsilat/fatura yaratmaz | U1 Z14,D1 §12 | mevcut migration araçları önce keşfedilecek | P11 |
| REQ-36 | Test/kanıt/gözlenebilirlik/restore ve manuel kontrollü yayın | D4,U1 | mevcut test/installer/ops altyapıları | P00,P09,P11 |

## Bütün matriste kullanılacak durum sütunları

P00 bu tabloyu aşağıdaki alanlarla genişletir; var olmayan kanıtı doldurmaz:

`requirement_id | canonical_work_id | backend_path_symbol_sha | frontend_path_symbol_sha | device_dependency | test_name | test_command | test_result | runtime_evidence | status | missing_piece | dependency | owner | rollback_boundary`

Yeni alan/sözleşme önerileri kodda varmış gibi işaretlenmez. Geçmiş master planda tamamlanmış fakat bu sürümde doğrulanamayan iş, önce sürüm/branch/deploy farkı açısından incelenir.
