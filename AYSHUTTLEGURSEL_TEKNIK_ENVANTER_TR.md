# AYShuttleGursel - Teknik Envanter

**Tarih:** 2026-10-29  
**Ana belge:** [Proje analizi](./AYSHUTTLEGURSEL_PROJE_ANALIZI_TR.md)  
**Kaynak:** [AYShuttleGursel](../php/AYShuttleGursel/), HEAD `0956c9be6bea96a849e198d9785b35c4f7590608` + inceleme anındaki yerel değişiklikler.

Bu ek dosya adları, inline action bildirimleri ve literal koleksiyon adlarından çıkarılmıştır. Liste bir canlı API testi veya tüm action'ların erişilebilirlik/yetki onayı değildir. Kaynaklar göreli bağlantılarla verilmiştir; kaynak proje aynı workspace yerleşiminde açılmalıdır.

## 1. Dizin rehberi ve sayım kuralı

| Dizin | Sorumluluk |
|---|---|
| [controllers](../php/AYShuttleGursel/controllers/) | Ana giriş, kurum yolcusu kaydı, parola yenileme, yardım |
| [models](../php/AYShuttleGursel/models/) | Sistem kimliği, form/yardımcı modeller |
| [yonetim](../php/AYShuttleGursel/modules/yonetim/) | Temel iş modelleri ve yönetim CRUD |
| [odeme](../php/AYShuttleGursel/modules/odeme/) | Yolcu bakiye, abonelik, rezervasyon, QR, hesap özeti |
| [raporlar](../php/AYShuttleGursel/modules/raporlar/) | Operasyon/mali rapor ve dışa aktarma |
| [services](../php/AYShuttleGursel/modules/services/) | Cihaz, kurum ve SAP servisleri; servis yönetim ekranları |
| [ticket_sales](../php/AYShuttleGursel/modules/ticket_sales/) | Ayrı alıcı/OTP oturumu ve QR bilet ödeme |
| [config](../php/AYShuttleGursel/config/) | Web, console, test ve bağlantı yapılandırmaları |
| [commands](../php/AYShuttleGursel/commands/) | Yalnızca örnek `HelloController` |
| [tests](../php/AYShuttleGursel/tests/) | Başlangıç şablonu testleri |

Sayım: Controller dizinlerindeki `*Controller.php` dosyaları; yorumlar çıkarıldıktan sonra `action` + büyük harfle başlayan metot bildirimleri. `actions()` ile tanımlanan error/captcha gibi harici action'lar inline sayısına dahil değildir. Arama modelleri `*Arama` veya `*Search` son ekine göre sayılır; koleksiyon sayısı olarak kullanılmaz.

Toplam **116 controller**, **633 inline action**, **208 model dosyası**, **637 görünüm**; bu yüzeylerde **968 PHP dosyası**. Modül içindeki component/service/Module dosyaları PHP toplamına dahildir.

## 2. Kök controller'lar

| Controller | Inline action'lar |
|---|---|
| [SiteController](../php/AYShuttleGursel/controllers/SiteController.php) | `Changelang`, `Index`, `Musterigirisi`, `Musterionayi`, `Musterikayit`, `Login`, `Logout`, `Mustericikisi`, `Contact`, `About` |
| [SifremiunuttumController](../php/AYShuttleGursel/controllers/SifremiunuttumController.php) | `Musterisifreyenileme`, `Kullanicisifreyenileme`, `Musteriyenisifre`, `Kullaniciyenisifre` |
| [YardimController](../php/AYShuttleGursel/controllers/YardimController.php) | `Index`, `View`, `Create`, `Update`, `Delete` |

Tablolarda action adlarından yalnızca `action` öneki çıkarılmıştır; bunlar doğrudan URL metni değildir. Örneğin `GetTickets` URL'de `get-tickets` olur.

## 3. Yolcu ve ödeme controller'ları

| Controller | Inline action'lar |
|---|---|
| [DefaultController](../php/AYShuttleGursel/modules/odeme/controllers/DefaultController.php) | `Index`, `Bakiyeyukle`, `Servisnerede`, `Haritadaizle`, `Directions`, `Locations`, `Kartbilgileri`, `Kvkkonayi`, `Tcentry`, `Akbanksanalpos`, `Yapikredi`, `Yapikrediodemeonay`, `Bkmekspress`, `Sartname`, `Gizlilik`, `Iptaliade`, `Musterisozlesmesi`, `Musterihizmetsozlesmesi` |
| [ErrorController](../php/AYShuttleGursel/modules/odeme/controllers/ErrorController.php) | Inline action yok; `actions()` haritası var |
| [HaftalikPlanController](../php/AYShuttleGursel/modules/odeme/controllers/HaftalikPlanController.php) | `Index`, `Data` |
| [HesapozetiController](../php/AYShuttleGursel/modules/odeme/controllers/HesapozetiController.php) | `Karthareketozeti` |
| [KartsizgecisController](../php/AYShuttleGursel/modules/odeme/controllers/KartsizgecisController.php) | `Index`, `Hareketdetayi`, `Bakiyecek` |
| [QrCodeController](../php/AYShuttleGursel/modules/odeme/controllers/QrCodeController.php) | `Generate` |
| [RezervasyonController](../php/AYShuttleGursel/modules/odeme/controllers/RezervasyonController.php) | `Index`, `Iptalet`, `Kisilistesi`, `Delete` |

Form/modeller: [SanalPosForm](../php/AYShuttleGursel/modules/odeme/models/SanalPosForm.php), [RezervasyonForm](../php/AYShuttleGursel/modules/odeme/models/RezervasyonForm.php), [QrCode](../php/AYShuttleGursel/modules/odeme/models/QrCode.php).

## 4. Rapor controller'ları

| Controller | Inline action'lar |
|---|---|
| [DefaultController](../php/AYShuttleGursel/modules/raporlar/controllers/DefaultController.php) | `Index` |
| [AracbazlisultopkazanclariController](../php/AYShuttleGursel/modules/raporlar/controllers/AracbazlisultopkazanclariController.php) | `Index`, `View`, `Create`, `Update`, `Tekrargonder`, `Delete` |
| [ArachareketraporuController](../php/AYShuttleGursel/modules/raporlar/controllers/ArachareketraporuController.php) | `Index`, `View`, `Binenkisilistesi`, `Silmeislemi`, `Silmeiptalislemi`, `Anlikhareketozetraporu`, `Create`, `Update`, `Doluluk`, `Suldoluluk`, `Ortalamasuldoluluk`, `Aracisformu`, `Aracsultoplamkazanc`, `Raporuyenidenolustur`, `Delete` |
| [HarekettoplamraporuController](../php/AYShuttleGursel/modules/raporlar/controllers/HarekettoplamraporuController.php) | `Index` |
| [KarthareketraporuController](../php/AYShuttleGursel/modules/raporlar/controllers/KarthareketraporuController.php) | `Index`, `View`, `Khrolustur`, `Borclukartbilgileri`, `Kartokutmasayisi`, `Karthareketozeti`, `Kisisulkullanimbilgisi`, `Create`, `Update`, `Kisiadsoyadlistesi`, `Delete` |
| [MusteriarackazancraporuController](../php/AYShuttleGursel/modules/raporlar/controllers/MusteriarackazancraporuController.php) | `Index`, `View`, `Create`, `Createlist`, `Update`, `Delete` |

Kalıcı rapor modelleri aşağıdaki koleksiyon tablosundadır. Koleksiyon bildirmeyen anlık/yardımcı rapor modelleri: [BorcluKartlarRaporu](../php/AYShuttleGursel/modules/raporlar/models/BorcluKartlarRaporu.php), [CihazIletisimHareketleri](../php/AYShuttleGursel/modules/raporlar/models/CihazIletisimHareketleri.php), [HareketToplamRaporu](../php/AYShuttleGursel/modules/raporlar/models/HareketToplamRaporu.php), [KartHareketOzeti](../php/AYShuttleGursel/modules/raporlar/models/KartHareketOzeti.php).

<a id="servis-envanteri"></a>
## 5. Servis controller'ları

Bu modülde **16 controller / 85 inline action** vardır. Yönetim CRUD ekranları da bulunduğundan hepsi cihaz API'si değildir.

| Controller | Inline action'lar |
|---|---|
| [DefaultController](../php/AYShuttleGursel/modules/services/controllers/DefaultController.php) | `Index` |
| [BaglantianahtariController](../php/AYShuttleGursel/modules/services/controllers/BaglantianahtariController.php) | `Index` |
| [KartbakiyeislemleriController](../php/AYShuttleGursel/modules/services/controllers/KartbakiyeislemleriController.php) | `Index`, `Kulluygulamagirisi`, `Kartbilgisi`, `Kartbakiyebilgisi`, `Plakabakiyecekmegirisi`, `Aracisbilgisi`, `Bakiyecekmeislemi`, `Kartbakiyecekmeislemi`, `Qrbakiyecekmeislemi`, `Temazsizbakiyecekmeislemi`, `Arackartcekimhareketleri`, `Karthareketiadesi`, `Cevrimdisicekilenbakiyeler`, `Cevrimdisibakiyeleriaktar`, `Cekilenkartsayisi`, `Araciskapatma`, `Arachareketraporuolustur`, `Karthareketraporuolustur`, `Kartyuklemehareketibilgiyenileme`, `Saatliulasimlistesi`, `Saatliulasimlistesiall`, `Zamanliaraciskapatma`, `Sonhareketleriisle`, `Kartbakiyeyukleme`, `Epostagonder`, `Smsgonder` |
| [ArachareketislemleriController](../php/AYShuttleGursel/modules/services/controllers/ArachareketislemleriController.php) | `Okanuniarachareketraporu`, `Rezervasyonbildirimi`, `Rszyedekkontrolu`, `Rezervasyoncekmislemiduzenleme` |
| [AracisatamaislemleriController](../php/AYShuttleGursel/modules/services/controllers/AracisatamaislemleriController.php) | `Saatliisatamasi`, `Uzaktanaracisatamasi`, `Iskapatma` |
| [CihazbaglantiislemleriController](../php/AYShuttleGursel/modules/services/controllers/CihazbaglantiislemleriController.php) | `Kulluygulamagirisi`, `Cihazbaglantibilgisi`, `Cihazisaldibilgisi`, `Kullaniciyetkisorgusu` |
| [EntegrebakiyeislemleriController](../php/AYShuttleGursel/modules/services/controllers/EntegrebakiyeislemleriController.php) | `Baglantianahtari`, `Kisisorgulama`, `Bakiyeyukle`, `Bakiyeyuklemeiptali`, `Yenikartbildirimi`, `Iadeislemi` |
| [FaturaislemleriController](../php/AYShuttleGursel/modules/services/controllers/FaturaislemleriController.php) | `Faturasizislemlistesi`, `Faturaliislemlistesi`, `Faturagonderimbildirimi`, `Faturayagecmeyenleriaktar` |
| [KarthareketleriController](../php/AYShuttleGursel/modules/services/controllers/KarthareketleriController.php) | `Sonbakiyebilgisi`, `Limitlikarthareketleri`, `Tariharaliklikarthareketleri`, `Kisimezunbildirimi`, `Kurumlar` |
| [SapentegrasyonislemleriController](../php/AYShuttleGursel/modules/services/controllers/SapentegrasyonislemleriController.php) | `Aracsultoplamkazanclarigonder`, `Kurumlistesi`, `Sulistesi`, `Guzergahvardiyabilgilerialma` |
| [AystekvalidatorController](../php/AYShuttleGursel/modules/services/controllers/AystekvalidatorController.php) | `Userlist`, `Kartdurumsilme` |
| [GuvenliiplistesiController](../php/AYShuttleGursel/modules/services/controllers/GuvenliiplistesiController.php) | `Index`, `View`, `Create`, `Update`, `Delete` |
| [HataliisatamasiduzeltmeleriController](../php/AYShuttleGursel/modules/services/controllers/HataliisatamasiduzeltmeleriController.php) | `Index`, `View`, `Create`, `Update`, `Delete` |
| [WebservisbaglantilistesiController](../php/AYShuttleGursel/modules/services/controllers/WebservisbaglantilistesiController.php) | `Index`, `View`, `Create`, `Update`, `Delete` |
| [WebservisislevleriController](../php/AYShuttleGursel/modules/services/controllers/WebservisislevleriController.php) | `Index`, `View`, `Create`, `Update`, `Delete` |
| [WebserviskullanicilariController](../php/AYShuttleGursel/modules/services/controllers/WebserviskullanicilariController.php) | `Index`, `View`, `Create`, `Update`, `Delete` |

### 5.1 Yeni istemci bağlarken kontrol edilecek sözleşmeler

| İşlem | Route | Girdi/çıktı notu |
|---|---|---|
| Bağlantı anahtarı | `/services/baglantianahtari/index` | POST; `username`, `pass` → `token`, `processCode`, `error` |
| Kart çözümleme | `/services/kartbakiyeislemleri/kartbilgisi` | POST; `kullanici`, `sifre`, `okunan_kart_num`; okunan kart kaydı da oluşturabilir |
| Kart bakiyesi | `/services/kartbakiyeislemleri/kartbakiyebilgisi` | POST; `kullanici`, `sifre`, `kartNum` → `kartBilgisi` |
| Araç işi | `/services/kartbakiyeislemleri/aracisbilgisi` | POST; kullanıcı bilgileri + `plaka` |
| Normal çekim | `/services/kartbakiyeislemleri/bakiyecekmeislemi` | POST; kart/SUL/iş/plaka/sürücü/konum; kurum özel QR dalları da var |
| Temassız çekim kaydı | `/services/kartbakiyeislemleri/temazsizbakiyecekmeislemi` | POST; ayrıca `islem_tutari`, `indirim_orani`; yükleme/çekim kayıt yolu normal çekimden farklı |
| SUL listesi | `/services/kartbakiyeislemleri/saatliulasimlistesi` | POST; kurum ve gün/saat listeleri içeren `sulBilgileri` |
| Harici yükleme | `/services/entegrebakiyeislemleri/bakiyeyukle` | Token + kişi/kurum/tutar/yükleme türü; ayrıntı ana belge §7 |
| Validator liste | `/services/aystekvalidator/userlist` | POST + Basic Auth; adı `userlist` olsa da plaka güvenlik listesini döndürür |
| Yolcu QR üretimi | `/odeme/qr-code/generate` | Yii kullanıcı kimliğini okur; `qrCode`, `balance`, `expiresAt` |
| Haftalık plan JSON | `/odeme/haftalik-plan/data` | `kurum`, `onlyUpcoming`, `limit`; yerel SUL verisi, 60 saniye cache |

Auth, zorunlu alanlar, mutasyon ve hata sözleşmesi action bazında kontrol edilmelidir. URL listesinden bütün çağrıların idempotent veya salt-okunur olduğu sonucu çıkarılamaz.

## 6. Misafir bilet controller'ları

| Controller | Inline action'lar |
|---|---|
| [AuthController](../php/AYShuttleGursel/modules/ticket_sales/controllers/AuthController.php) | `Login`, `SendPhoneCode`, `SendEmailCode`, `VerifyCode`, `Logout` |
| [CustomerController](../php/AYShuttleGursel/modules/ticket_sales/controllers/CustomerController.php) | `Index`, `GetCustomer`, `PersonalInfo`, `EditPersonalInfo` |
| [PaymentController](../php/AYShuttleGursel/modules/ticket_sales/controllers/PaymentController.php) | `InitiatePayment`, `PaymentResponse`, `BuyTicket`, `GetRoutes`, `GetPrice` |
| [TicketController](../php/AYShuttleGursel/modules/ticket_sales/controllers/TicketController.php) | `Index`, `List`, `GetTickets` |

Ortak bileşen: [components/BaseController](../php/AYShuttleGursel/modules/ticket_sales/components/BaseController.php). OTP servisi: [AuthService](../php/AYShuttleGursel/modules/ticket_sales/services/AuthService.php). `PaymentController` ortak bileşenden türemez; callback için CSRF istisnası kendi `beforeAction` metodundadır.

## 7. Yönetim controller'ları

**80 controller** aşağıda tam olarak listelenmiştir. İkinci sütun yalnızca `Index/View/Create/Update/Delete` dışındaki inline action'ları gösterir. `-`, bu beş ad dışında action bulunmadığı anlamına gelir; beş CRUD işleminin hepsinin bulunduğunu, erişilebilir veya güvenli olduğunu göstermez.

| Controller | CRUD dışı inline action'lar |
|---|---|
| [Abonelikucretleri](../php/AYShuttleGursel/modules/yonetim/controllers/AbonelikucretleriController.php) | - |
| [Abonelikyuklemeleri](../php/AYShuttleGursel/modules/yonetim/controllers/AbonelikyuklemeleriController.php) | - |
| [Aracabaglicihazlar](../php/AYShuttleGursel/modules/yonetim/controllers/AracabaglicihazlarController.php) | `Cihaziletisimlistesi`, `Cihazbaglantidurumlari` |
| [Aracguvenlikbilgileri](../php/AYShuttleGursel/modules/yonetim/controllers/AracguvenlikbilgileriController.php) | - |
| [Aracisatamasi](../php/AYShuttleGursel/modules/yonetim/controllers/AracisatamasiController.php) | `Viewresult`, `Aracraporunuolustur`, `Aracislistesi`, `Tarihinsaatleri` |
| [Aracisdevirleri](../php/AYShuttleGursel/modules/yonetim/controllers/AracisdevirleriController.php) | `Onaysizaracisdevsayisi` |
| [Aracisformbilgileri](../php/AYShuttleGursel/modules/yonetim/controllers/AracisformbilgileriController.php) | - |
| [Arackapasiteturleri](../php/AYShuttleGursel/modules/yonetim/controllers/ArackapasiteturleriController.php) | - |
| [Araclar](../php/AYShuttleGursel/modules/yonetim/controllers/AraclarController.php) | `Araclistesi` |
| [Aracsoforulasimatamasi](../php/AYShuttleGursel/modules/yonetim/controllers/AracsoforulasimatamasiController.php) | `Gunlukplanliisiolanaraclar` |
| [Baglantilar](../php/AYShuttleGursel/modules/yonetim/controllers/BaglantilarController.php) | `Qrtest` |
| [Biletler](../php/AYShuttleGursel/modules/yonetim/controllers/BiletlerController.php) | - |
| [Cevrimdisikarthareketleri](../php/AYShuttleGursel/modules/yonetim/controllers/CevrimdisikarthareketleriController.php) | `Cdisibakiyeleriaktar` |
| [Default](../php/AYShuttleGursel/modules/yonetim/controllers/DefaultController.php) | - |
| [Duraklar](../php/AYShuttleGursel/modules/yonetim/controllers/DuraklarController.php) | - |
| [Duyurular](../php/AYShuttleGursel/modules/yonetim/controllers/DuyurularController.php) | - |
| [Externalshuttle](../php/AYShuttleGursel/modules/yonetim/controllers/ExternalshuttleController.php) | `Startprocess`, `Endprocess` |
| [Faturaliyuklemeler](../php/AYShuttleGursel/modules/yonetim/controllers/FaturaliyuklemelerController.php) | `Tarihduzeltme` |
| [Grupindirimlerikumesi](../php/AYShuttleGursel/modules/yonetim/controllers/GrupindirimlerikumesiController.php) | - |
| [Gunler](../php/AYShuttleGursel/modules/yonetim/controllers/GunlerController.php) | - |
| [Gunlukkasadekontlari](../php/AYShuttleGursel/modules/yonetim/controllers/GunlukkasadekontlariController.php) | - |
| [Guzergahlar](../php/AYShuttleGursel/modules/yonetim/controllers/GuzergahlarController.php) | `Duraklar` |
| [Guzergahvardiyakapasitelistesi](../php/AYShuttleGursel/modules/yonetim/controllers/GuzergahvardiyakapasitelistesiController.php) | `Verileriguncelleme`, `Saataraliklari` |
| [Guzergahvardiyatipleri](../php/AYShuttleGursel/modules/yonetim/controllers/GuzergahvardiyatipleriController.php) | - |
| [Iletisimyontemleri](../php/AYShuttleGursel/modules/yonetim/controllers/IletisimyontemleriController.php) | - |
| [Isikuniversitesikisilistesi](../php/AYShuttleGursel/modules/yonetim/controllers/IsikuniversitesikisilistesiController.php) | - |
| [Kartcekimhareketleri](../php/AYShuttleGursel/modules/yonetim/controllers/KartcekimhareketleriController.php) | `Bakiyecek`, `Raporuyenidenolustur` |
| [Kartdevirhareketleri](../php/AYShuttleGursel/modules/yonetim/controllers/KartdevirhareketleriController.php) | - |
| [Kartdurumhareketleri](../php/AYShuttleGursel/modules/yonetim/controllers/KartdurumhareketleriController.php) | - |
| [Kartdurumlari](../php/AYShuttleGursel/modules/yonetim/controllers/KartdurumlariController.php) | - |
| [Kartiadehareketleri](../php/AYShuttleGursel/modules/yonetim/controllers/KartiadehareketleriController.php) | `Raporuyenidenolustur` |
| [Kartiadeturleri](../php/AYShuttleGursel/modules/yonetim/controllers/KartiadeturleriController.php) | - |
| [Kartlar](../php/AYShuttleGursel/modules/yonetim/controllers/KartlarController.php) | `Okunankart`, `Okunankartislembitir`, `Kartlistesi` |
| [Kartsatisfiyatlari](../php/AYShuttleGursel/modules/yonetim/controllers/KartsatisfiyatlariController.php) | - |
| [Kartsatislari](../php/AYShuttleGursel/modules/yonetim/controllers/KartsatislariController.php) | - |
| [Kartturleri](../php/AYShuttleGursel/modules/yonetim/controllers/KartturleriController.php) | - |
| [Kartyuklemehareketleri](../php/AYShuttleGursel/modules/yonetim/controllers/KartyuklemehareketleriController.php) | `Kisiadsoyadlistesi`, `Kasabilgisi`, `Raporuyenidenolustur` |
| [Kartyuklemeturleri](../php/AYShuttleGursel/modules/yonetim/controllers/KartyuklemeturleriController.php) | - |
| [Kisidurumlari](../php/AYShuttleGursel/modules/yonetim/controllers/KisidurumlariController.php) | - |
| [Kisigrubugecerliliksureleri](../php/AYShuttleGursel/modules/yonetim/controllers/KisigrubugecerliliksureleriController.php) | - |
| [Kisigrupguzergahindirimi](../php/AYShuttleGursel/modules/yonetim/controllers/KisigrupguzergahindirimiController.php) | - |
| [Kisigruplari](../php/AYShuttleGursel/modules/yonetim/controllers/KisigruplariController.php) | `Gruplar` |
| [Kisikartdurumlari](../php/AYShuttleGursel/modules/yonetim/controllers/KisikartdurumlariController.php) | - |
| [Kisikartlari](../php/AYShuttleGursel/modules/yonetim/controllers/KisikartlariController.php) | - |
| [Kisiler](../php/AYShuttleGursel/modules/yonetim/controllers/KisilerController.php) | `Kisidetaybilgisi`, `Kisilistesi`, `Kisiadlistesi`, `Kisisoyadlistesi`, `Kisiepostalistesi`, `Kurumabaglikisiler`, `Kurumnumaralistesi` |
| [Kisisiniflari](../php/AYShuttleGursel/modules/yonetim/controllers/KisisiniflariController.php) | - |
| [Kullanicidurumlari](../php/AYShuttleGursel/modules/yonetim/controllers/KullanicidurumlariController.php) | - |
| [Kullanicikurumayarlari](../php/AYShuttleGursel/modules/yonetim/controllers/KullanicikurumayarlariController.php) | `Yenikurumayari` |
| [Kullanicilar](../php/AYShuttleGursel/modules/yonetim/controllers/KullanicilarController.php) | `Liste` |
| [Kullanicisifreleri](../php/AYShuttleGursel/modules/yonetim/controllers/KullanicisifreleriController.php) | - |
| [Kullaniciyetkidurumlari](../php/AYShuttleGursel/modules/yonetim/controllers/KullaniciyetkidurumlariController.php) | - |
| [Kullaniciyetkigruplari](../php/AYShuttleGursel/modules/yonetim/controllers/KullaniciyetkigruplariController.php) | - |
| [Kurumlar](../php/AYShuttleGursel/modules/yonetim/controllers/KurumlarController.php) | - |
| [Language](../php/AYShuttleGursel/modules/yonetim/controllers/LanguageController.php) | - |
| [Mobilcekimhareketleri](../php/AYShuttleGursel/modules/yonetim/controllers/MobilcekimhareketleriController.php) | - |
| [Modulgorselbaglantilar](../php/AYShuttleGursel/modules/yonetim/controllers/ModulgorselbaglantilarController.php) | - |
| [Ozyeginuniservicekisilistesi](../php/AYShuttleGursel/modules/yonetim/controllers/OzyeginuniservicekisilistesiController.php) | `Kurumdantumverilerial` |
| [Rezervasyonisatamasi](../php/AYShuttleGursel/modules/yonetim/controllers/RezervasyonisatamasiController.php) | `Koltuksayisi`, `Tarihinsaatleri` |
| [Saataraliklari](../php/AYShuttleGursel/modules/yonetim/controllers/SaataraliklariController.php) | - |
| [Saatler](../php/AYShuttleGursel/modules/yonetim/controllers/SaatlerController.php) | - |
| [Saatliulasimlistesi](../php/AYShuttleGursel/modules/yonetim/controllers/SaatliulasimlistesiController.php) | `Gununsaatleri`, `Rezervasyontarihinsaatleri`, `Tarihinsaatleri` |
| [Sabanciuniservicekisilistesi](../php/AYShuttleGursel/modules/yonetim/controllers/SabanciuniservicekisilistesiController.php) | - |
| [Sanalposlar](../php/AYShuttleGursel/modules/yonetim/controllers/SanalposlarController.php) | - |
| [Seyahatrezervasyonlari](../php/AYShuttleGursel/modules/yonetim/controllers/SeyahatrezervasyonlariController.php) | `Uetdsreport` |
| [Sifresureleri](../php/AYShuttleGursel/modules/yonetim/controllers/SifresureleriController.php) | `Surelistesi` |
| [Site](../php/AYShuttleGursel/modules/yonetim/controllers/SiteController.php) | `Login`, `Logout` |
| [Soforler](../php/AYShuttleGursel/modules/yonetim/controllers/SoforlerController.php) | `Soforlistesi` |
| [Tasinabilircihazturleri](../php/AYShuttleGursel/modules/yonetim/controllers/TasinabilircihazturleriController.php) | `Cihazturlistesi` |
| [Ulasimucretleri](../php/AYShuttleGursel/modules/yonetim/controllers/UlasimucretleriController.php) | - |
| [Uygsurumgruplari](../php/AYShuttleGursel/modules/yonetim/controllers/UygsurumgruplariController.php) | - |
| [Uygsurumleri](../php/AYShuttleGursel/modules/yonetim/controllers/UygsurumleriController.php) | - |
| [Uygulamacihazlari](../php/AYShuttleGursel/modules/yonetim/controllers/UygulamacihazlariController.php) | - |
| [Veriaktarim](../php/AYShuttleGursel/modules/yonetim/controllers/VeriaktarimController.php) | `Sabanciveriaktarim`, `Ozyeginveriaktarim`, `Ozyeginserviceveriaktarim`, `Ozyegintutaryukleme`, `Ozyeginbakiyeaktarma`, `Sabancitutaryukleme` |
| [Yazarkasaposlar](../php/AYShuttleGursel/modules/yonetim/controllers/YazarkasaposlarController.php) | - |
| [Yazarkasaposraporlari](../php/AYShuttleGursel/modules/yonetim/controllers/YazarkasaposraporlariController.php) | - |
| [Yedekrezervasyonlar](../php/AYShuttleGursel/modules/yonetim/controllers/YedekrezervasyonlarController.php) | `Uetdsreport` |
| [Yetkiadlari](../php/AYShuttleGursel/modules/yonetim/controllers/YetkiadlariController.php) | - |
| [Yetkigruplari](../php/AYShuttleGursel/modules/yonetim/controllers/YetkigruplariController.php) | - |
| [Yetkiler](../php/AYShuttleGursel/modules/yonetim/controllers/YetkilerController.php) | - |
| [Yuklemekampanyalari](../php/AYShuttleGursel/modules/yonetim/controllers/YuklemekampanyalariController.php) | - |

<a id="koleksiyonlar"></a>
## 8. MongoDB koleksiyonları

Model dizinlerinde **93 farklı literal koleksiyon adı** bildirilmiştir. Tablo canlı DB'deki varlık, kayıt sayısı, indeks veya veri kalitesi tespiti değildir. Modeldeki `unique` kuralı fiziksel MongoDB unique index'inin kanıtı sayılmamalıdır.

### 8.1 Kök, ödeme, rapor, servis ve bilet koleksiyonları (18)

| Koleksiyon | Bildiren model |
|---|---|
| `iletisim_kutusu` | [IletisimKutusu](../php/AYShuttleGursel/models/IletisimKutusu.php) |
| `kullanici_durumlari` | [KullaniciDurumlari](../php/AYShuttleGursel/models/KullaniciDurumlari.php) |
| `kullanicilar` | [Kullanicilar](../php/AYShuttleGursel/models/Kullanicilar.php) |
| `sifre_sureleri` | [SifreSureleri](../php/AYShuttleGursel/models/SifreSureleri.php) |
| `yardim` | [Yardim](../php/AYShuttleGursel/models/Yardim.php) |
| `qr_codes` | [QrCode](../php/AYShuttleGursel/modules/odeme/models/QrCode.php) |
| `arac_bazli_sul_top_kazanclari` | [AracBazliSulTopKazanclari](../php/AYShuttleGursel/modules/raporlar/models/AracBazliSulTopKazanclari.php) |
| `arac_hareket_raporu` | [AracHareketRaporu](../php/AYShuttleGursel/modules/raporlar/models/AracHareketRaporu.php) |
| `kart_hareket_raporu` | [KartHareketRaporu](../php/AYShuttleGursel/modules/raporlar/models/KartHareketRaporu.php) |
| `kart_hareket_raporu_tmp` | [KartHareketRaporuTmp](../php/AYShuttleGursel/modules/raporlar/models/KartHareketRaporuTmp.php) |
| `musteri_arac_kazanc_raporu` | [MusteriAracKazancRaporu](../php/AYShuttleGursel/modules/raporlar/models/MusteriAracKazancRaporu.php) |
| `guvenli_ip_listesi` | [GuvenliIpListesi](../php/AYShuttleGursel/modules/services/models/GuvenliIpListesi.php) |
| `hatali_is_atamasi_duzeltmeleri` | [HataliIsAtamasiDuzeltmeleri](../php/AYShuttleGursel/modules/services/models/HataliIsAtamasiDuzeltmeleri.php) |
| `web_servis_baglanti_listesi` | [WebServisBaglantiListesi](../php/AYShuttleGursel/modules/services/models/WebServisBaglantiListesi.php) |
| `web_servis_islevleri` | [WebServisIslevleri](../php/AYShuttleGursel/modules/services/models/WebServisIslevleri.php) |
| `web_servis_kullanicilari` | [WebServisKullanicilari](../php/AYShuttleGursel/modules/services/models/WebServisKullanicilari.php) |
| `ticket_buyers` | [TicketBuyers](../php/AYShuttleGursel/modules/ticket_sales/models/TicketBuyers.php) |
| `tickets` | [Tickets](../php/AYShuttleGursel/modules/ticket_sales/models/Tickets.php) |

### 8.2 Yönetim koleksiyonları (75)

| Koleksiyon | Bildiren model |
|---|---|
| `abonelik_ucretleri` | [AbonelikUcretleri](../php/AYShuttleGursel/modules/yonetim/models/AbonelikUcretleri.php) |
| `abonelik_yuklemeleri` | [AbonelikYuklemeleri](../php/AYShuttleGursel/modules/yonetim/models/AbonelikYuklemeleri.php) |
| `arac_guvenlik_bilgileri` | [AracGuvenlikBilgileri](../php/AYShuttleGursel/modules/yonetim/models/AracGuvenlikBilgileri.php) |
| `arac_is_atama` | [AracIsAtama](../php/AYShuttleGursel/modules/yonetim/models/AracIsAtama.php) |
| `arac_is_devirleri` | [AracIsDevirleri](../php/AYShuttleGursel/modules/yonetim/models/AracIsDevirleri.php) |
| `arac_is_form_bilgileri` | [AracIsFormBilgileri](../php/AYShuttleGursel/modules/yonetim/models/AracIsFormBilgileri.php) |
| `arac_kapasite_turleri` | [AracKapasiteTurleri](../php/AYShuttleGursel/modules/yonetim/models/AracKapasiteTurleri.php) |
| `arac_sofor_ulasim_atamasi` | [AracSoforUlasimAtamasi](../php/AYShuttleGursel/modules/yonetim/models/AracSoforUlasimAtamasi.php) |
| `araca_bagli_cihazlar` | [AracaBagliCihazlar](../php/AYShuttleGursel/modules/yonetim/models/AracaBagliCihazlar.php) |
| `araclar` | [Araclar](../php/AYShuttleGursel/modules/yonetim/models/Araclar.php) |
| `baglantilar` | [Baglantilar](../php/AYShuttleGursel/modules/yonetim/models/Baglantilar.php) |
| `biletler` | [Biletler](../php/AYShuttleGursel/modules/yonetim/models/Biletler.php) |
| `cevrimdisi_kart_hareketleri` | [CevrimdisiKartHareketleri](../php/AYShuttleGursel/modules/yonetim/models/CevrimdisiKartHareketleri.php) |
| `diger_kisiler` | [DigerKisiler](../php/AYShuttleGursel/modules/yonetim/models/DigerKisiler.php) |
| `duraklar` | [Duraklar](../php/AYShuttleGursel/modules/yonetim/models/Duraklar.php) |
| `duyurular` | [Duyurular](../php/AYShuttleGursel/modules/yonetim/models/Duyurular.php) |
| `external_shuttle` | [ExternalShuttle](../php/AYShuttleGursel/modules/yonetim/models/ExternalShuttle.php) |
| `faturali_yuklemeler` | [FaturaliYuklemeler](../php/AYShuttleGursel/modules/yonetim/models/FaturaliYuklemeler.php) |
| `grup_indirimleri_kumesi` | [GrupIndirimleriKumesi](../php/AYShuttleGursel/modules/yonetim/models/GrupIndirimleriKumesi.php) |
| `gunler` | [Gunler](../php/AYShuttleGursel/modules/yonetim/models/Gunler.php) |
| `gunluk_kasa_dekontlari` | [GunlukKasaDekontlari](../php/AYShuttleGursel/modules/yonetim/models/GunlukKasaDekontlari.php) |
| `guzergah_vardiya_kapasite_listesi` | [GuzergahVardiyaKapasiteListesi](../php/AYShuttleGursel/modules/yonetim/models/GuzergahVardiyaKapasiteListesi.php) |
| `guzergah_vardiya_tipleri` | [GuzergahVardiyaTipleri](../php/AYShuttleGursel/modules/yonetim/models/GuzergahVardiyaTipleri.php) |
| `guzergahlar` | [Guzergahlar](../php/AYShuttleGursel/modules/yonetim/models/Guzergahlar.php) |
| `iletisim_yontemleri` | [IletisimYontemleri](../php/AYShuttleGursel/modules/yonetim/models/IletisimYontemleri.php) |
| `kart_cekim_hareketleri` | [KartCekimHareketleri](../php/AYShuttleGursel/modules/yonetim/models/KartCekimHareketleri.php) |
| `kart_devir_hareketleri` | [KartDevirHareketleri](../php/AYShuttleGursel/modules/yonetim/models/KartDevirHareketleri.php) |
| `kart_durum_hareketleri` | [KartDurumHareketleri](../php/AYShuttleGursel/modules/yonetim/models/KartDurumHareketleri.php) |
| `kart_durumlari` | [KartDurumlari](../php/AYShuttleGursel/modules/yonetim/models/KartDurumlari.php) |
| `kart_iade_hareketleri` | [KartIadeHareketleri](../php/AYShuttleGursel/modules/yonetim/models/KartIadeHareketleri.php) |
| `kart_iade_turleri` | [KartIadeTurleri](../php/AYShuttleGursel/modules/yonetim/models/KartIadeTurleri.php) |
| `kart_satis_fiyatlari` | [KartSatisFiyatlari](../php/AYShuttleGursel/modules/yonetim/models/KartSatisFiyatlari.php) |
| `kart_satislari` | [KartSatislari](../php/AYShuttleGursel/modules/yonetim/models/KartSatislari.php) |
| `kart_turleri` | [KartTurleri](../php/AYShuttleGursel/modules/yonetim/models/KartTurleri.php) |
| `kart_yukleme_hareketleri` | [KartYuklemeHareketleri](../php/AYShuttleGursel/modules/yonetim/models/KartYuklemeHareketleri.php) |
| `kart_yukleme_turleri` | [KartYuklemeTurleri](../php/AYShuttleGursel/modules/yonetim/models/KartYuklemeTurleri.php) |
| `kartlar` | [Kartlar](../php/AYShuttleGursel/modules/yonetim/models/Kartlar.php) |
| `kisi_durumlari` | [KisiDurumlari](../php/AYShuttleGursel/modules/yonetim/models/KisiDurumlari.php) |
| `kisi_grubu_gecerlilik_sureleri` | [KisiGrubuGecerlilikSureleri](../php/AYShuttleGursel/modules/yonetim/models/KisiGrubuGecerlilikSureleri.php) |
| `kisi_grup_guzergah_indirimi` | [KisiGrupGuzergahIndirimi](../php/AYShuttleGursel/modules/yonetim/models/KisiGrupGuzergahIndirimi.php) |
| `kisi_gruplari` | [KisiGruplari](../php/AYShuttleGursel/modules/yonetim/models/KisiGruplari.php) |
| `kisi_kart_durumlari` | [KisiKartDurumlari](../php/AYShuttleGursel/modules/yonetim/models/KisiKartDurumlari.php) |
| `kisi_kartlari` | [KisiKartlari](../php/AYShuttleGursel/modules/yonetim/models/KisiKartlari.php) |
| `kisi_siniflari` | [KisiSiniflari](../php/AYShuttleGursel/modules/yonetim/models/KisiSiniflari.php) |
| `kisiler` | [Kisiler](../php/AYShuttleGursel/modules/yonetim/models/Kisiler.php) |
| `kullanici_kurum_ayarlari` | [KullaniciKurumAyarlari](../php/AYShuttleGursel/modules/yonetim/models/KullaniciKurumAyarlari.php) |
| `kullanici_yetki_durumlari` | [KullaniciYetkiDurumlari](../php/AYShuttleGursel/modules/yonetim/models/KullaniciYetkiDurumlari.php) |
| `kullanici_yetki_gruplari` | [KullaniciYetkiGruplari](../php/AYShuttleGursel/modules/yonetim/models/KullaniciYetkiGruplari.php) |
| `kurumlar` | [Kurumlar](../php/AYShuttleGursel/modules/yonetim/models/Kurumlar.php) |
| `language` | [Language](../php/AYShuttleGursel/modules/yonetim/models/Language.php) |
| `mobil_cekim_hareketleri` | [MobilCekimHareketleri](../php/AYShuttleGursel/modules/yonetim/models/MobilCekimHareketleri.php) |
| `modul_gorsel_baglantilar` | [ModulGorselBaglantilar](../php/AYShuttleGursel/modules/yonetim/models/ModulGorselBaglantilar.php) |
| `okunan_kartlar` | [OkunanKartlar](../php/AYShuttleGursel/modules/yonetim/models/OkunanKartlar.php) |
| `ozyeginuni_service_kisi_listesi` | [OzyeginuniServiceKisiListesi](../php/AYShuttleGursel/modules/yonetim/models/OzyeginuniServiceKisiListesi.php) |
| `payment_transaction` | [PaymentTransaction](../php/AYShuttleGursel/modules/yonetim/models/PaymentTransaction.php) |
| `rezervasyon_is_atamasi` | [RezervasyonIsAtamasi](../php/AYShuttleGursel/modules/yonetim/models/RezervasyonIsAtamasi.php) |
| `saat_araliklari` | [SaatAraliklari](../php/AYShuttleGursel/modules/yonetim/models/SaatAraliklari.php) |
| `saatler` | [Saatler](../php/AYShuttleGursel/modules/yonetim/models/Saatler.php) |
| `saatli_ulasim_listesi` | [SaatliUlasimListesi](../php/AYShuttleGursel/modules/yonetim/models/SaatliUlasimListesi.php) |
| `sabanciuni_service_kisi_listesi` | [SabanciuniServiceKisiListesi](../php/AYShuttleGursel/modules/yonetim/models/SabanciuniServiceKisiListesi.php) |
| `sanalposlar` | [Sanalposlar](../php/AYShuttleGursel/modules/yonetim/models/Sanalposlar.php) |
| `seyahat_rezervasyonlari` | [SeyahatRezervasyonlari](../php/AYShuttleGursel/modules/yonetim/models/SeyahatRezervasyonlari.php) |
| `soforler` | [Soforler](../php/AYShuttleGursel/modules/yonetim/models/Soforler.php) |
| `tasinabilir_cihaz_turleri` | [TasinabilirCihazTurleri](../php/AYShuttleGursel/modules/yonetim/models/TasinabilirCihazTurleri.php) |
| `ulasim_ucretleri` | [UlasimUcretleri](../php/AYShuttleGursel/modules/yonetim/models/UlasimUcretleri.php) |
| `uyg_surum_gruplari` | [UygSurumGruplari](../php/AYShuttleGursel/modules/yonetim/models/UygSurumGruplari.php) |
| `uyg_surumleri` | [UygSurumleri](../php/AYShuttleGursel/modules/yonetim/models/UygSurumleri.php) |
| `uygulama_cihazlari` | [UygulamaCihazlari](../php/AYShuttleGursel/modules/yonetim/models/UygulamaCihazlari.php) |
| `yazar_kasa_pos_raporlari` | [YazarKasaPosRaporlari](../php/AYShuttleGursel/modules/yonetim/models/YazarKasaPosRaporlari.php) |
| `yazar_kasa_poslar` | [YazarKasaPoslar](../php/AYShuttleGursel/modules/yonetim/models/YazarKasaPoslar.php) |
| `yedek_rezervasyonlar` | [YedekRezervasyonlar](../php/AYShuttleGursel/modules/yonetim/models/YedekRezervasyonlar.php) |
| `yetki_adlari` | [YetkiAdlari](../php/AYShuttleGursel/modules/yonetim/models/YetkiAdlari.php) |
| `yetki_gruplari` | [YetkiGruplari](../php/AYShuttleGursel/modules/yonetim/models/YetkiGruplari.php) |
| `yetkiler` | [Yetkiler](../php/AYShuttleGursel/modules/yonetim/models/Yetkiler.php) |
| `yukleme_kampanyalari` | [YuklemeKampanyalari](../php/AYShuttleGursel/modules/yonetim/models/YuklemeKampanyalari.php) |

### 8.3 Ana veri sözlüğü

Aşağıdaki alanlar model `attributes()` bildirimlerinden seçilmiştir; fiziksel veri tipi/indeks şeması değildir. Kimlik/audit alanları birçok modelde `_id`, `kayit_tarihi`, `kayit_eden_kullanici`, `guncelleme_tarihi`, `guncelleyen_kullanici` biçimindedir. Ödeme/bilet/harici servis modellerinde farklı audit adları ve tarih biçimleri de kullanılır.

| Model | Başlıca alanlar | İlişki/yorum |
|---|---|---|
| Kartlar | `kart_num`, `kurum`, `kart_turu`, `kart_durumu` | Saklanan bakiye alanı yok |
| Kisiler | `kurum`, `kurum_numarasi`, `kurum_kayit_num`, `kisi_eposta`, `kisi_gruplari`, `kisi_sinifi`, `online_odeme_kapali` | Kimlik/iletişim ve grup bağlamı |
| KisiKartlari | `kisi`, `kart_num`, `kisi_kart_durumu` | Kişi-kart ilişkisi |
| SaatliUlasimListesi | `saatli_ulasim_listesi_adi`, `guzergah`, `guzergah_kodu`, tarih aralığı, gün alanları, `*_rz`, `mobil_yayinlama`, `rezervasyon_yayinla`, `kimlik_bilgisi_zorunlu` | Saat listeleri ve yayın/rezervasyon kuralları |
| AracIsAtama | `saatli_ulasim_listesi`, `plaka`, `sofor`, `hareket_tarihi`, `hareket_saati`, `arac_is_durumu`, `guvenlik_sorgu_kodu` | Çekimlerde `is_numarasi` ile ilişki |
| UlasimUcretleri | `saatli_ulasim_listesi`, `kartli_ucret`, `kartsiz_ucret`, `kupon_ucret`, `gecerlilik_baslangic_tarihi`, `gecerlilik_bitis_tarihi` | Tarihe bağlı üç ücret |
| GrupIndirimleriKumesi | `grup_indirimi_adi`, `indirim_orani`, `zaman_asimi`, `gunler`, `limit_kontrol_durumu`, `limit`, `saat_araliklari` | Limit ve zaman koşulları |
| KisiGrupGuzergahIndirimi | `kisi_grubu`, `guzergah`, `grup_indirim_kumesi` | Grup-güzergah-indirim bağı |
| KartYuklemeHareketleri | `kart_num`, `yukleme_turu`, `islem_kodu`, `dogrulama_kodu`, `referens_numarasi`, `gelen_yanit_metni`, `islem_tutari`, `fatura_tablosu` | Pozitif bakiye kaynağı; alan adı `referens_numarasi` şeklindedir |
| KartCekimHareketleri | `is_numarasi`, `kart_num`, `cekim_turu`, SUL/plaka/sürücü/zaman/konum, `islem_tutari`, `indirim_orani`, `hesaplanan_tutar`, `iade_durumu` | Bakiyede düşülen alan `hesaplanan_tutar` |
| KartIadeHareketleri | `kart_num`, `iade_turu`, `islem_kodu`, `islem_tutari`, `aciklama` | Bakiye etkisi iade türüne bağlı |
| CevrimdisiKartHareketleri | İş/kart/SUL/plaka, `cekim_zamani`, `aktarim_durumu`, ücret/indirim/tutar | Bekleyen/aktarılmış cihaz işlemi |
| AbonelikUcretleri | `abonelik_ucreti`, `indirim_grubu`, `baslangic_tarihi`, `bitis_tarihi` | Ücret ve hak tanımı |
| AbonelikYuklemeleri | `kisi`, `ad_soyad`, ücret/grup/tarih aralığı | Nakit bakiye yüklemesiyle aynı değil |
| FaturaliYuklemeler | `kart_num`, `siparis_num`, `fatura_islemi`, `fatura_islem_turu`, `fatura_islem_kodu`, `fatura_pos_kodu`, `fatura_tutari`, `fatura_eposta`, `fatura_gonderimi`, `odeme_alma_tarihi` | Banka/mali işlem ile gönderim takibi |
| PaymentTransaction | `user_id`, `order_id`, `vp_order_id`, `transaction_id`, `card_num`, `payment_method`, `ip`, `notification`, `order_note` | Sipariş/aşama kayıtları; bütünlük kısıtı ayrıca doğrulanmalı |
| RezervasyonIsAtamasi | `sul`, `plaka`, `sofor`, `rez_tarihi`, `rez_saati`, `rez_tarih_saat` | Kapasite/rezervasyon iş planı |
| SeyahatRezervasyonlari | `kurum`, `kisi_id`, kimlik/yolcu bilgileri, `sul`, `binecegi_durak`, hareket tarih/saati, `rezervasyon_durumu`, `bakiye_cekme_durumu`, `is_numarasi` | Yolcu ve tahsilat durumu |
| YedekRezervasyonlar | `kurum`, `kisi_grup_id`, yolcu bilgileri, `sul`, hareket tarih/saati, `rezervasyon_gecis_durumu` | Yedekten geçiş durumu |
| YuklemeKampanyalari | `kurum`, `kampanya_adi`, `yuklenecek_asgari_tutar`, `yuklenecek_azami_tutar`, `hediye_puan`, geçerlilik tarihleri | Kurum/tutar aralığına bağlı ek yükleme |
| YazarKasaPosRaporlari | `kurum`, `yazar_kasa_pos`, `rapor_tarihi`, `rapor_numarasi`, `nakit_tutari`, `kredi_kart_tutari`, `ek_bilgi` | Günlük cihaz mali raporu |
| TicketBuyers | `name`, `surname`, `id_number`, `email`, `phone`, `i_number`, `i_name`, `verified`, `created_at` | Sistem kullanıcısından ayrı müşteri |
| Tickets | `customer_id`, `route`, `order_id`, banka referansları, `price`, `qr_code`, `status`, `expires_at`, `created_at` | `active`, `used`, `cancelled`, `expired`; üretim/doğrulama zaman tipi uyumsuz |
| QrCode | `person_id`, `code`, `card_num`, `created_at`, `expires_at` | `01TCK`, milisaniye zamanları |

Model adları bu bölümde §8.1/§8.2'deki kaynak bağlantılarının kısa adlarıdır; formel ilişkisel foreign key bildirimleri değildir.

### 8.4 SQL ve koleksiyon bildirmeyen modeller

Doğrudan SQL ActiveRecord kullanan üç eski entegrasyon modeli:

- [IsikUniversitesiKisiListesi](../php/AYShuttleGursel/modules/yonetim/models/IsikUniversitesiKisiListesi.php)
- [OzyeginKisiListesi](../php/AYShuttleGursel/modules/yonetim/models/OzyeginKisiListesi.php)
- [OzyeginKartBakiyeleri](../php/AYShuttleGursel/modules/yonetim/models/OzyeginKartBakiyeleri.php)

Bağlantı [db config](../php/AYShuttleGursel/config/db.php) tarafından sağlanır; gerçek dış veri tabanlarının bu dosyayla çalıştığı test edilmemiştir.

Model dizinindeki her sınıf koleksiyon değildir. Örneğin [MusteriGirisiForm](../php/AYShuttleGursel/modules/yonetim/models/MusteriGirisiForm.php), [Gunler](../php/AYShuttleGursel/modules/yonetim/models/Gunler.php) gibi adlandırmalar kaynak sınıfa bakılarak yorumlanmalıdır; koleksiyon bildirimi olanlar yukarıdaki listede yer alır. `*Arama`/`*Search` sınıfları çoğunlukla filtre/data provider yüzeyidir; ayrı koleksiyon sayılmaz.

## 9. Kaynak okuma sırası

Bir iş kuralını değiştirmeden önce:

1. İlgili controller'ın access/verb/CSRF ve parametre okuma biçimini kontrol et.
2. Ana modelin `rules()`, `beforeSave()` ve çağrılan iş metodunu birlikte oku.
3. Aynı metodu kullanan yönetim, ödeme ve cihaz akışlarını bul.
4. Hareket koleksiyonu, rapor okuma modeli ve banka/fatura etkisini ayrı değerlendir.
5. Kurum, zaman ve tutar birimlerini doğrula.
6. [Ana belgedeki kabul matrisi](./AYSHUTTLEGURSEL_PROJE_ANALIZI_TR.md#test) ile beklenen sonucu tanımla.

Bu ekin tam envanter kapsamı dosya/action/koleksiyon bildirimleridir. İşlevsel hatalar, runtime uyumsuzlukları ve canlı doğrulama sınırları [ana analizde](./AYSHUTTLEGURSEL_PROJE_ANALIZI_TR.md#bulgular) tutulur.
