---
Title: Fleet Tally Kullanım Kılavuzu
Version: 1.1.0
Status: Production Reference
Last Updated: 2026-06-09
Classification: Internal - Operations Guide
Module Path: /internal/modules/fleet_tally
Audience: Operasyon, Filo Yönetimi, Entegrasyon, Destek, QA
Related Documents:
  - ./modules/fleet_tally.md
  - ./INTEGRATIONS_PLATFORM_GUIDE.md
  - ../05-operations/INTEGRATIONS_SCHEDULED_PULL_RUNBOOK_QA_MATRIX.md
  - ../../plans/FLEET_TALLY_SYSTEM_ANALYSIS_TR.md
---

# Fleet Tally Kullanım Kılavuzu

Bu doküman, AYSTEK Core içindeki **Fleet Tally / GPS Çetele Kontrol** modülünün kullanıcılar ve operasyon ekipleri tarafından nasıl kullanılacağını anlatır.

Fleet Tally, dış sistemden gelen servis/çetele kayıtlarını GPS verisiyle karşılaştırır; aracın istenen zamanda, istenen koordinat çevresinde bulunup bulunmadığını doğrular ve sonucu dış sisteme geri gönderir.

> Kısa tanım: “Araç gerçekten gitmesi gereken yere, gitmesi gereken zamanda gitti mi?” sorusunun GPS kanıtıyla cevaplanması.

---

## İçindekiler

1. [Kimler Kullanır?](#1-kimler-kullanır)
2. [Modülün İş Akışı](#2-modülün-i̇ş-akışı)
3. [Ekranlar ve Ne İçin Kullanılır?](#3-ekranlar-ve-ne-i̇çin-kullanılır)
4. [İlk Kurulum Sırası](#4-i̇lk-kurulum-sırası)
5. [Gürsel Araç Takip İçin Hazır Seed Kayıtları](#5-gürsel-araç-takip-i̇çin-hazır-seed-kayıtları)
6. [GPS Sağlayıcıları](#6-gps-sağlayıcıları)
7. [Dış Sistem Bağlayıcıları](#7-dış-sistem-bağlayıcıları)
8. [Kontrol Tipleri](#8-kontrol-tipleri)
9. [Zamanlayıcılar](#9-zamanlayıcılar)
10. [Manuel Kontrol Başlatma](#10-manuel-kontrol-başlatma)
11. [Çalıştırma Geçmişi Takibi](#11-çalıştırma-geçmişi-takibi)
12. [Doğrulama İşleri ve Sonuç Kodları](#12-doğrulama-i̇şleri-ve-sonuç-kodları)
13. [Araç Eşlemeleri](#13-araç-eşlemeleri)
14. [Konum Verisi Kaynak Önceliği](#14-konum-verisi-kaynak-önceliği)
15. [Sık Karşılaşılan Durumlar](#15-sık-karşılaşılan-durumlar)
16. [Operasyon Kontrol Listesi](#16-operasyon-kontrol-listesi)
17. [Yetki ve Güvenlik Notları](#17-yetki-ve-güvenlik-notları)
18. [Migration Sonrası Test Runbook](#18-migration-sonrası-test-runbook)

---

## 1. Kimler Kullanır?

| Rol | Kullanım Amacı |
|---|---|
| Operasyon yöneticisi | Günlük servis çetele kontrollerini başlatır ve sonuçları takip eder. |
| Filo / araç takip ekibi | GPS sağlayıcılarını, plaka-cihaz eşlemelerini ve konum doğrulamalarını yönetir. |
| Entegrasyon ekibi | SAP/ERP bağlayıcılarını ve push-back sonuçlarını izler. |
| Destek ekibi | Başarısız kayıtları inceler, gerekirse yeniden dener. |
| Denetim / raporlama | Hangi aracın hangi kayıtta başarılı/başarısız olduğunu kanıt olarak inceler. |

---

## 2. Modülün İş Akışı

Fleet Tally uçtan uca şu sırayla çalışır:

```text
Dış Sistem (SAP/ERP)
    │
    │ bekleyen çetele kayıtları çekilir
    ▼
Fleet Tally Zamanlayıcı
    │
    │ her kayıt VerificationJob olur
    ▼
GPS Sağlayıcıları
    │
    │ Trio / N2Mobil / Mobiliz / Arvento konum geçmişi sorgulanır
    ▼
Doğrulama Algoritması
    │
    │ zaman penceresi + hedef koordinat + yarıçap kontrol edilir
    ▼
Doğrulama Sonucu
    │
    │ ok / not_in_zone / no_location / error ...
    ▼
Dış Sisteme Sonuç Gönderimi
```

Her zamanlayıcı çalışması bir **Run** kaydı üretir. Her dış sistem kaydı ise bir **Verification Job** olarak takip edilir.

---

## 3. Ekranlar ve Ne İçin Kullanılır?

Fleet Tally ekranları admin panelde şu adreste bulunur:

`/admin/fleet-tally/...`

| Ekran | Yol | Ne İçin Kullanılır? |
|---|---|---|
| Kontrol Tipleri | `/admin/fleet-tally/control-types` | Personel/okul gibi kontrol algoritması kurallarını tanımlar. |
| GPS Sağlayıcılar | `/admin/fleet-tally/gps-providers` | Trio, N2Mobil, Mobiliz, Arvento bağlantı ayarlarını yönetir. |
| Veri Bağlayıcılar | `/admin/fleet-tally/connectors` | SAP/ERP gibi dış sistemden iş çekme ve sonuç gönderme ayarlarını yönetir. |
| Zamanlayıcılar | `/admin/fleet-tally/schedules` | Otomatik veya manuel doğrulama çalıştırmalarını yönetir. |
| Doğrulama İşleri | `/admin/fleet-tally/verifications` | Plaka/kayıt bazında kontrol sonucunu ve detayını gösterir. |
| Çalışma Geçmişi | `/admin/fleet-tally/runs` | Her çalıştırmanın toplam/başarılı/başarısız sayılarını gösterir. |
| Araç Eşlemeleri | `/admin/fleet-tally/vehicle-mappings` | Plaka ile GPS sağlayıcıdaki araç/cihaz ID eşlemesini yönetir. |
| Kullanım Kılavuzu | `/admin/fleet-tally/user-guide` | Panel içi hızlı kullanım rehberidir. |

---

## 4. İlk Kurulum Sırası

Yeni bir tenant/site için Fleet Tally kurulurken önerilen sıra:

1. **Kontrol tiplerini doğrula**
   - `personel_servis_v1`
   - `okul_servis_v1`
2. **GPS sağlayıcılarını tanımla**
   - Trio
   - N2Mobil / WeTrack
   - Mobiliz
   - Arvento
3. **GPS kimlik bilgilerini gir**
   - Kullanıcı adı, parola, token, PIN gibi bilgiler seed ile boş gelir; operatör/admin tarafından doldurulur.
4. **Dış sistem bağlayıcılarını kontrol et**
   - SAP/ERP endpoint, fetch method, push method, record type.
5. **Araç eşlemelerini kontrol et**
   - Öncelik Stock modülündeki araç `gps_device_ids` ve aktif GPS varlıklarıdır.
   - Manuel eşleme sadece Stock verisi yoksa veya bilinçli override gerekiyorsa tutulmalıdır.
6. **Zamanlayıcıyı pasiften aktife al**
   - Seed kayıtları güvenli başlamak için pasif gelebilir.
7. **Önce manuel çalıştırma yap**
   - İlk denemede küçük lookback/max job ile test edilmesi önerilir.
8. **Run ve job sonuçlarını incele**
   - Başarılı/başarısız ayrımı ve push durumu doğrulanmadan otomatik cron aktif edilmemelidir.

---

## 5. Gürsel Araç Takip İçin Hazır Seed Kayıtları

`aractakip.gurseltur.com.tr` sitesi için seed sonrası temel kayıtlar hazır gelir.

### 5.1 GPS sağlayıcıları

| Sağlayıcı | Durum | Not |
|---|---|---|
| Trio Mobil | Aktif | Kullanıcı adı/parola sonradan girilir. |
| N2Mobil / WeTrack | Aktif | `fleet_id=143` seed ile gelir; kullanıcı adı/parola sonradan girilir. |
| Mobiliz | Aktif | İç MongoDB koleksiyonu üzerinden okunur. Varsayılan koleksiyon: `mobiliz_location_events`. |
| Arvento | Aktif | `username`, `pin1`, `pin2` sonradan girilir. Plakadan cihaz node çözümleme desteklenir. |

### 5.2 Dış sistem bağlayıcıları

| Bağlayıcı | Record Type | Fetch Method | Push Method | Lookback |
|---|---|---|---|---|
| Gürsel SAP Personel Servis | `personel_servis` | `OnSiparisKayitListesi` | `OnSiparisKayitGuncelle` | 90 gün |
| Gürsel SAP Okul Servis | `okul_servis` | `okulHakediskayitListesiServis` | `okulHakedisGuncellemeServis` | 35 gün |

### 5.3 Zamanlayıcılar

| Zamanlayıcı | Varsayılan Cron | Durum |
|---|---:|---|
| Gürsel Personel Servis GPS Çetele Kontrolü | `0 0 * * * *` | Pasif |
| Gürsel Okul Servis GPS Çetele Kontrolü | `0 30 * * * *` | Pasif |

> Seed kayıtları canlı veri göndermeyi yanlışlıkla başlatmamak için zamanlayıcıları pasif oluşturur. Test tamamlanınca aktife alınmalıdır.

---

## 6. GPS Sağlayıcıları

GPS sağlayıcı ekranında her kayıt bir konum veri kaynağını temsil eder.

| Alan | Açıklama |
|---|---|
| Sağlayıcı tipi | `trio`, `n2`, `mobiliz`, `arvento` |
| Görünen ad | Operasyon ekranlarında görünen isim |
| Base URL | Sağlayıcı API adresi; Mobiliz gibi iç kaynaklarda boş olabilir |
| Timezone offset | Türkiye için genellikle `180` dakika |
| Credentials | Kullanıcı adı/parola/token/PIN gibi gizli alanlar |
| Aktif/Pasif | Sağlayıcının zamanlayıcılar tarafından kullanılıp kullanılmayacağı |

### 6.1 Sağlayıcı bazlı notlar

#### Trio Mobil

- Araç konum geçmişi Trio API üzerinden çekilir.
- Authorization header değeri Trio kontratına uygun şekilde token olarak gönderilir.
- Zaman aralığı Türkiye saat dilimi offset’iyle hazırlanır.

#### N2Mobil / WeTrack

- Auth yanıtından `access` token ve `app_url` alınır.
- Rapor sorgusu `app_url` üzerinden yapılır.
- Gürsel için varsayılan `fleet_id` değeri `143` olarak seed edilir.

#### Mobiliz

- Konumlar dış API yerine iç MongoDB koleksiyonundan okunur.
- Eski sistemden gelen saat farkı için +3 saat düzeltmesi uygulanır.

#### Arvento

- SOAP `GeneralReport` çağrısıyla geçmiş konum alınır.
- Cihaz node numarası bilinmiyorsa plaka üzerinden `GetLicensePlateNodeMappings` ile çözülebilir.
- `username`, `pin1`, `pin2` alanları doldurulmalıdır.

---

## 7. Dış Sistem Bağlayıcıları

Dış sistem bağlayıcısı, çetele kontrol edilecek kayıtların nereden alınacağını ve sonucun nereye gönderileceğini tanımlar.

| Alan | Açıklama |
|---|---|
| Connector Type | `sap_soap` veya `rest_api` |
| Base URL | Ana endpoint adresi; SAP SOAP için Gürsel varsayılanı `https://sapservice.gurseltur.com.tr/WebServices/OnSiparis.asmx` |
| Integration Connector Code | Integrations modülündeki üst entegrasyon kaydı; Gürsel için `GURSEL_SAP_FLEET_TALLY` |
| Fetch Config | Bekleyen işleri çekmek için endpoint/method/limit/lookback ayarları |
| Push Config | Sonuçları geri göndermek için endpoint/method ayarları |
| Field Mapping | Kayıt tipi, sabit yarıçap gibi dönüştürme ayarları |
| Credentials | Dış sistem kullanıcı adı/parola gibi gizli bilgiler |

### 7.1 Integrations modülü ile ilişki

`/admin/fleet-tally/connectors` ekranı Fleet Tally’nin operasyonel adapter konfigürasyonunu yönetir. Bu kayıtlar doğrudan `fleet_tally_external_connectors` koleksiyonunda tutulur ve doğrulama sırasında domain’e özel `sap_soap` veya `rest_api` adapter’ı tarafından kullanılır.

Integrations modülü ise aynı dış sistemi enterprise entegrasyon envanteri, allowlist, governance, action kataloğu ve izlenebilirlik açısından temsil eder.

Gürsel SAP için seed sonrası şu yapı beklenir:

| Kayıt | Koleksiyon | Amaç |
|---|---|---|
| `GURSEL_SAP_FLEET_TALLY` | `integration_connectors` | SAP SOAP endpoint, allowlist ve entegrasyon kataloğu |
| `GURSEL_SAP_PERSONEL_FETCH` | `integration_actions` | `OnSiparisKayitListesi` action tanımı |
| `GURSEL_SAP_PERSONEL_PUSH` | `integration_actions` | `OnSiparisKayitGuncelle` action tanımı |
| `GURSEL_SAP_OKUL_FETCH` | `integration_actions` | `okulHakediskayitListesiServis` action tanımı |
| `GURSEL_SAP_OKUL_PUSH` | `integration_actions` | `okulHakedisGuncellemeServis` action tanımı |
| `gursel_sap_personel_v1` | `fleet_tally_external_connectors` | Personel servis Fleet Tally adapter ayarı |
| `gursel_sap_okul_v1` | `fleet_tally_external_connectors` | Okul servis Fleet Tally adapter ayarı |

> Önemli: Fleet Tally SAP doğrulama akışı şu an Integrations executor üzerinden değil, Fleet Tally’nin domain’e özel SOAP adapter’ı üzerinden çalışır. Integrations kaydı aynı servisi kurumsal entegrasyon envanteri ve yönetim standardı içinde görünür hale getirir.

### 7.2 Gürsel SAP seed ayarları

Seed kayıtları aşağıdaki varsayılanları oluşturur:

| Alan | Personel Servis | Okul Servis |
|---|---|---|
| Connector Type | `sap_soap` | `sap_soap` |
| Integration Code | `GURSEL_SAP_FLEET_TALLY` | `GURSEL_SAP_FLEET_TALLY` |
| Fetch Method | `OnSiparisKayitListesi` | `okulHakediskayitListesiServis` |
| Push Method | `OnSiparisKayitGuncelle` | `okulHakedisGuncellemeServis` |
| Record Type | `personel_servis` | `okul_servis` |
| Status Filter | `G` | Yok |
| Lookback | 90 gün | 35 gün |

Credentials alanında Gürsel SAP kullanıcı adı/parolası seed ile gelir. UI üzerinden credential güncellenirken alanlar boş bırakılırsa mevcut değerler korunur.

### 7.3 Test amaçlı iş çekme

Bağlayıcı ekranındaki fetch test/iş çekme özelliğiyle dış sistemden kayıt dönüp dönmediği kontrol edilmelidir.

Beklenen başarılı test:

- Kayıt listesi gelir.
- Plaka alanı doludur.
- Tarih alanı parse edilebilir.
- En az bir zaman penceresi oluşur.
- Hedef koordinat ve yarıçap mantıklıdır.

---

## 8. Kontrol Tipleri

Kontrol tipi, bir işin hangi algoritmayla ve hangi kurallarla doğrulanacağını belirler.

### 8.1 `personel_servis_v1`

Personel servis kayıtları için kullanılır.

| Kural | Değer |
|---|---|
| Zaman penceresi | Kayıttan gelen bir veya daha fazla pencere |
| Yarıçap | Kayıttan gelir; yoksa varsayılan kullanılır |
| Başarı kuralı | Tüm pencereler başarılı olmalıdır |
| Tipik radius | 700 metre |

### 8.2 `okul_servis_v1`

Okul servis hakediş kayıtları için kullanılır.

| Kural | Değer |
|---|---|
| Zaman penceresi | Giriş ve çıkış olmak üzere iki pencere |
| Yarıçap | Sabit 100 metre |
| Başarı kuralı | Giriş ve çıkış kontrolleri başarılı olmalıdır |

---

## 9. Zamanlayıcılar

Zamanlayıcı, hangi bağlayıcıdan iş çekileceğini, hangi kontrol tipiyle doğrulanacağını ve hangi GPS sağlayıcılarının hangi sırayla deneneceğini belirler.

| Alan | Açıklama |
|---|---|
| Bağlayıcı | SAP/ERP gibi dış sistem |
| Kontrol tipi | Personel veya okul servis kuralı |
| GPS sağlayıcıları | Öncelik sırasına göre denenir |
| Cron | Otomatik çalıştırma zamanı |
| Lookback days | Kaç gün geriye bakılacağı |
| Max jobs per run | Tek çalıştırmada işlenecek maksimum kayıt |
| Retry policy | Geçici hatalarda kaç kez yeniden deneneceği |
| Active hours | Zamanlayıcının çalışmasına izin verilen saat aralığı |
| Aktif/Pasif | Cron tarafından tetiklenip tetiklenmeyeceği |

### 9.1 Cron formatı

Sistem 5 veya 6 alanlı cron ifadelerini kabul eder.

```text
5 alanlı örnek:
0 6 * * *        Her gün 06:00

6 alanlı örnek:
0 0 * * * *      Her saat başı
0 30 * * * *     Her saat 30. dakikada
```

> Gürsel seed kayıtlarında 6 alanlı cron kullanılır. İlk alan saniyedir.

### 9.2 GPS sağlayıcı sırası

Bir kayıtta sağlayıcı ipucu varsa sistem bu sağlayıcıyı öne alabilir. Aksi halde zamanlayıcıdaki `gps_provider_ids` sırası izlenir.

Gürsel için önerilen varsayılan sıra:

1. N2Mobil / WeTrack
2. Trio Mobil
3. Mobiliz
4. Arvento

---

## 10. Manuel Kontrol Başlatma

İlk kurulum ve destek süreçlerinde manuel çalıştırma önerilir.

1. `Zamanlayıcılar` ekranına girin.
2. İlgili zamanlayıcıyı bulun.
3. Zamanlayıcı pasif olsa bile manuel test için `Çalıştır` butonunu kullanın.
4. Sistem bir `VerificationRun` oluşturur.
5. Sonucu `Çalışma Geçmişi` ekranından takip edin.

Manuel çalıştırmanın API karşılığı:

```http
POST /admin/fleet-tally/schedules/{schedule_id}/run
```

---

## 11. Çalıştırma Geçmişi Takibi

`Çalışma Geçmişi` ekranı her toplu çalıştırmayı gösterir.

| Alan | Anlamı |
|---|---|
| Status | `running`, `completed`, `partial`, `failed` |
| Total Jobs | Bu çalıştırmada dış sistemden gelen toplam iş sayısı |
| Passed Jobs | Başarılı doğrulanan iş sayısı |
| Failed Jobs | Başarısız doğrulanan iş sayısı |
| Error Jobs | Teknik hata veya belirsiz hata sayısı |
| Started / Ended | Çalıştırmanın başlangıç ve bitiş zamanı |

### 11.1 Run durumları

| Durum | Açıklama | Operasyon Yorumu |
|---|---|---|
| `running` | İşlem devam ediyor | Bekleyin veya canlı/polling ile takip edin. |
| `completed` | Tüm işler başarılı tamamlandı | Normal durum. |
| `partial` | Bazı işler başarılı, bazıları başarısız | Doğrulama işleri ekranında başarısızları inceleyin. |
| `failed` | Hiç başarılı iş yok veya kritik hata var | Bağlayıcı/GPS credential/kayıt formatı kontrol edilmeli. |

---

## 12. Doğrulama İşleri ve Sonuç Kodları

`Doğrulama İşleri` ekranı her plaka/kayıt bazında sonucu gösterir.

| Alan | Açıklama |
|---|---|
| Plaka | Kontrol edilen araç |
| Belge No / Kalem No | Dış sistem kaydı |
| Job Date | Hizmet tarihi |
| Overall Status | İşin genel sonucu |
| Push Status | Sonucun dış sisteme gönderilme durumu |
| GPS Provider | Sonucu üreten GPS sağlayıcı |
| Windows | Her zaman penceresinin ayrı sonucu |

### 12.1 Sonuç kodları

| Kod | Kullanıcıya Anlamı | Ne Yapılmalı? |
|---|---|---|
| `ok` veya boş | Araç hedef bölgede doğrulandı | İşlem tamam. |
| `not_in_zone` | Konum var ama hedef bölge dışında | Koordinat, yarıçap ve araç rotası kontrol edilmeli. |
| `route_passed` | Nokta hedefte değil ama rota hedef bölgeden geçmiş | Operasyon kuralına göre kabul edilebilir; algoritma çıktısı incelenmeli. |
| `no_device` | Araç için GPS cihazı/ID bulunamadı | Araç eşlemesi veya stok/araç varlığı kontrol edilmeli. |
| `no_gsm` | Cihaz var ama veri hattı/GSM problemi var | Cihaz ve hat durumu kontrol edilmeli. |
| `no_location` | Tüm gün için konum yok | GPS sağlayıcı/API/cihaz veri gönderimi kontrol edilmeli. |
| `empty_window` | Gün içinde konum var ama istenen pencerede yok | Saat penceresi, vardiya/servis saati ve GPS gecikmesi kontrol edilmeli. |
| `retry_later` | Veri henüz gelmemiş olabilir | Daha sonra tekrar denenebilir. |
| `error` | Teknik hata oluştu | Hata mesajı, provider credential ve bağlantı kontrol edilmeli. |
| `unknown` | Maksimum deneme sonrası net sonuç alınamadı | Destek/entegrasyon ekibi incelemeli. |

### 12.2 Push status

| Durum | Anlamı |
|---|---|
| `pending` | Sonuç henüz dış sisteme gönderilmedi. |
| `sent` | Sonuç dış sisteme başarıyla gönderildi. |
| `failed` | Sonuç gönderilemedi; connector/push ayarı kontrol edilmeli. |

---

## 13. Araç Eşlemeleri

Araç eşlemeleri, bir plakayı GPS sağlayıcıdaki araç/cihaz ID’sine bağlar. Normal akışta Fleet Tally önce Stock modülündeki araç `gps_device_ids` bilgisini ve araca takılı aktif GPS varlığını kullanır; bu ekran manuel override ve operasyonel kontrol içindir.

| Alan | Açıklama |
|---|---|
| GPS Provider ID | Hangi sağlayıcıya ait eşleme olduğu |
| Plate | Normalize edilmiş plaka |
| External ID | Sağlayıcıdaki araç/cihaz ID |
| Note | Operasyon notu |
| Source | `manual`, `stock_vehicle` veya `stock_asset`; Stock kaynaklı satırlar salt okunurdur |

### 13.1 Ne zaman gerekir?

- Stock kaydında araç için ilgili sağlayıcı ID’si yoksa.
- Sağlayıcıdaki araç/cihaz ID’si geçici olarak Stock’tan farklı kullanılacaksa.
- Destek sürecinde bilinçli bir override gerekiyorsa.

### 13.2 Arvento özel notu

Arvento’da cihaz node numarası bilinmiyorsa sistem plaka üzerinden node mapping çağrısı yapabilir. Yine de kritik araçlarda manuel eşleme tutulması destek süreçlerini hızlandırır.

---

## 14. Konum Verisi Kaynak Önceliği

Fleet Tally konum geçmişini **local-first** mantığıyla kullanır:

1. Önce AYSTEK Core içindeki `stock_vehicle_locations_history` kayıtları aranır.
2. İstenen zaman aralığında veri varsa doğrulama bu lokal kayıtlarla yapılır.
3. Lokal veri yoksa ilgili GPS sağlayıcının web servisi çağrılır.
4. Web servisten gelen konumlar tekrar aynı history cache’e `gps_provider_api` kaynağıyla yazılır.

Bu sayede aynı kontrol tekrarlandığında mümkün olduğunca dış web servis çağrısı yapılmaz.

Bu davranış admin panelde `/admin/fleet-tally/settings` ekranından değiştirilebilir.

| Mod | Davranış | Ne zaman seçilir? |
|---|---|---|
| `local_first` | Önce AYSTEK Core veritabanına bakar, yoksa GPS sağlayıcı web servisine gider. | Varsayılan ve önerilen ayardır. |
| `local_only` | Sadece AYSTEK Core’daki konum history/cache kayıtlarını kullanır. | Sağlayıcı servisine çıkılmaması gereken test, denetim veya maliyet kontrollü dönemlerde. |
| `provider_only` | Lokal history/cache’i atlar, doğrudan GPS sağlayıcı web servisine gider. | Lokal cache şüphesi varsa veya sağlayıcı yanıtı canlı doğrulanacaksa. |

> Not: `provider_only` modunda sağlayıcıdan dönen başarılı konumlar yine `gps_provider_api` kaynağıyla history cache’e yazılır. Böylece moda tekrar `local_first` seçildiğinde son çekilen veriler kullanılabilir.

### 14.1 GPS sağlayıcı push endpoint’i

GPS sağlayıcılar veya AYSTEK’e ait Go microservice konum verisini aşağıdaki endpoint’e gönderebilir:

`POST /public/fleet-tally/gps-locations/ingest`

Zorunlu header’lar:

| Header | Açıklama |
|---|---|
| `X-Tenant-ID` | Tenant kodu/id bilgisi |
| `X-Site-ID` | Site/store ObjectID bilgisi |
| `X-Aystek-Signature` | HMAC-SHA256 imza; `sha256=<hex>` formatı önerilir |

Örnek payload:

```json
{
  "provider_type": "trio",
  "plate": "34ABC123",
  "vehicle_external_id": "862491078392964",
  "latitude": 41.0123,
  "longitude": 29.0456,
  "speed_kph": 35,
  "recorded_at": "2026-06-05T09:15:00Z",
  "raw_payload": {
    "source": "gps-microservice"
  }
}
```

`vehicle_id` gönderilirse doğrudan kullanılır. Gönderilmezse sistem Stock araçlarını plaka veya sağlayıcıya ait `vehicle_external_id` üzerinden çözer.

### 14.2 Kaynak takibi

Konum kayıtlarında kaynak alanı tutulur:

| Source | Anlamı |
|---|---|
| `gps_provider_push` | GPS sağlayıcı/microservice AYSTEK Core’a veri gönderdi |
| `gps_provider_api` | Lokal veri yoktu; Fleet Tally sağlayıcı web servisinden çekip cache’e yazdı |
| `connector` | Genel entegrasyon/telematics kaynaklı kayıt |
| `manual` | Operatör/admin tarafından manuel girilen kayıt |

Kayıtlarda ayrıca `provider_type`, `plate` ve `vehicle_external_id` tutulduğu için hangi sağlayıcının hangi cihaz/plaka için veri gönderdiği izlenebilir.

### 14.3 Eski konum verisi temizliği

`stock_vehicle_locations_history` için periyodik retention işi API açılışında başlar.

| Env | Varsayılan | Açıklama |
|---|---:|---|
| `FLEET_TALLY_GPS_LOCATION_RETENTION_ENABLED` | `true` | Temizleme işini açar/kapatır |
| `FLEET_TALLY_GPS_LOCATION_RETENTION_DAYS` | `30` | Kaç günden eski history silinecek |
| `FLEET_TALLY_GPS_LOCATION_RETENTION_INTERVAL` | `6h` | Temizlik çalışma aralığı |
| `FLEET_TALLY_GPS_LOCATION_INGEST_SECRET` | boş | Push endpoint imza secret’ı; production’da zorunludur |
| `STOCK_LOCATION_INGEST_SECRET` | boş | Geriye uyumlu alternatif imza secret’ı |

Canlı ortamda push endpoint’i secret olmadan istek kabul etmez; local/dev ortamda test kolaylığı için secret boş bırakılabilir.

---

## 15. Sık Karşılaşılan Durumlar

### 15.1 Zamanlayıcı çalışmıyor

Kontrol edin:

- Zamanlayıcı `is_active=true` mi?
- Cron ifadesi geçerli mi?
- Active hours içinde mi?
- API loglarında `Fleet-tally cron worker started` mesajı var mı?
- Tenant/site bağlamı doğru mu?

### 15.2 Manuel çalıştırma kayıt üretmiyor

Kontrol edin:

- Bağlayıcı dış sistemden kayıt çekebiliyor mu?
- Lookback gün sayısı doğru mu?
- Dış sistemde `statu=G` gibi bekleyen kayıt var mı?
- Max jobs per run çok düşük mü?

### 15.3 `no_device` çok fazla

Kontrol edin:

- Plaka normalize edilmiş mi?
- Araç eşlemeleri doğru mu?
- Stok/araç varlığı üzerinde GPS cihaz bilgisi var mı?
- SAP kaydında GPS sağlayıcı ipucu yanlış mı?

### 15.4 `no_location` çok fazla

Kontrol edin:

- GPS sağlayıcı credential doğru mu?
- Araç gerçekten veri gönderiyor mu?
- Sağlayıcıdaki vehicle ID doğru mu?
- Saat dilimi/offset doğru mu?

### 15.5 `empty_window` çok fazla

Kontrol edin:

- Hizmet saati doğru mu?
- Zaman penceresi çok dar mı?
- GPS sağlayıcıda veri gecikmeli mi geliyor?
- Operasyon saatleri dış sistemde yanlış mı?

### 15.6 `push_status=failed`

Kontrol edin:

- Dış sistem push endpoint erişilebilir mi?
- SAP/ERP kullanıcı adı/parolası doğru mu?
- Push method doğru mu?
- Dış sistem aynı kayda tekrar update kabul ediyor mu?

---

## 16. Operasyon Kontrol Listesi

Canlıya almadan önce:

- [ ] Tenant/site doğru seçildi.
- [ ] Gürsel Araç Takip görünüm profilinde fleet-tally ekranları görünür.
- [ ] GPS sağlayıcı credentials girildi.
- [ ] SAP/ERP connector credentials girildi.
- [ ] En az bir GPS sağlayıcı aktif.
- [ ] Kontrol tipi aktif.
- [ ] Zamanlayıcı GPS sağlayıcı listesi dolu.
- [ ] Manuel çalıştırma başarıyla run oluşturdu.
- [ ] En az bir başarılı `VerificationJob` görüldü.
- [ ] Push status `sent` olarak doğrulandı.
- [ ] Başarısız örneklerin hata açıklamaları incelendi.
- [ ] Otomatik cron ancak testten sonra aktif edildi.

Günlük takip:

- [ ] `Çalışma Geçmişi` ekranında son run durumu kontrol edilir.
- [ ] `partial` veya `failed` run varsa job detayları incelenir.
- [ ] `push_status=failed` kayıtlar takip edilir.
- [ ] `no_device` artışı varsa araç eşleme/stok verisi kontrol edilir.
- [ ] `no_location` artışı varsa GPS sağlayıcı veya cihaz durumu kontrol edilir.

---

## 17. Yetki ve Güvenlik Notları

Fleet Tally ekranları server-side RBAC ile korunur.

| Modül Yetkisi | Kapsam |
|---|---|
| `fleet_tally_control_types` | Kontrol tipi yönetimi |
| `fleet_tally_gps_providers` | GPS sağlayıcı yönetimi |
| `fleet_tally_connectors` | Dış sistem bağlayıcı yönetimi |
| `fleet_tally_schedules` | Zamanlayıcı yönetimi ve manuel çalıştırma |
| `fleet_tally_verifications` | Doğrulama job/run izleme ve retry |

Güvenlik ilkeleri:

- Credential alanları seed ile boş gelir; canlı bilgiler UI veya güvenli operasyon akışıyla girilir.
- Seed tekrar çalıştırıldığında mevcut credential değerleri ezilmemelidir.
- Dış sistem ve GPS servis çağrıları allowlist/connector güvenlik politikalarıyla sınırlandırılır.
- Tenant/site izolasyonu zorunludur; farklı site verisi aynı ekranda karışmamalıdır.
- Hata kayıtları incelenirken kullanıcı adı/parola/token gibi hassas değerler paylaşılmamalıdır.

---

## 18. Migration Sonrası Test Runbook

Bu bölüm, `make migrate-gursel-reset` veya `make migrate-gursel` sonrası GPS Çetele Kontrol sistemini uçtan uca test etmek için kullanılır.

### 18.1 Servisleri başlatma

Backend API:

```bash
cd /workspace/go/aystek-core
make dev
```

Varsayılan API adresi `http://localhost:9090` olur. API açılışında seed gerekiyorsa ve ortam `dev` değilse `SEED_ON_START=true` ile başlatılabilir; normal `dev` ortamında seed otomatik çalışır ve mevcut credential alanlarını ezmemelidir.

Frontend:

```bash
cd /workspace/node/aystek-core-next
npm run dev
```

Varsayılan panel adresi `http://localhost:3000` olur. Frontend, local geliştirmede varsayılan olarak `http://localhost:9090` API adresini kullanır.

### 18.2 İlk kontrol sırası

1. `http://localhost:3000` üzerinden admin kullanıcıyla giriş yapın.
2. Tenant/site olarak **Gürsel Turizm / Araç Takip** bağlamını seçin.
3. `/admin/fleet-tally/gps-providers` ekranında `Arvento`, `Mobiliz`, `Trio Mobil`, `N2Mobil / WeTrack` kayıtlarının tekil ve aktif olduğunu doğrulayın.
4. `/admin/fleet-tally/vehicle-mappings` ekranında sağlayıcı rozetlerindeki sayıların Stock modülündeki araç GPS varlıklarından beslendiğini kontrol edin.
5. `/admin/fleet-tally/settings` ekranında konum kaynağı modunu seçin:
   - İlk genel test için önerilen mod: `local_first`
   - Sadece sağlayıcı API yanıtını görmek için: `provider_only`
   - Sadece AYSTEK Core history/cache verisini test etmek için: `local_only`
6. `/admin/fleet-tally/connectors` ekranında Gürsel personel/okul bağlayıcılarının aktif ve credential alanlarının dolu olduğunu kontrol edin.
7. `/admin/fleet-tally/schedules` ekranında ilgili zamanlayıcıyı bulun ve ilk denemede `lookback_days` ile `max_jobs_per_run` değerlerini küçük tutun.

### 18.3 Bağlayıcıdan kayıt dönüyor mu?

Zamanlayıcıyı çalıştırmadan önce dış sistemden çetele kaydı gelip gelmediği şu API ile hızlıca doğrulanabilir:

```http
GET /admin/fleet-tally/connectors/{connector_id}/fetch-jobs?from=2026-06-08T00:00:00Z&to=2026-06-09T23:59:59Z&limit=10
```

Başarılı yanıtta `count > 0` beklenir. Kayıt yoksa önce connector ayarları, tarih aralığı, dış sistemdeki bekleyen kayıt statüsü ve credential bilgileri kontrol edilmelidir.

### 18.4 Manuel run başlatma

Panelden:

1. `/admin/fleet-tally/schedules` ekranına girin.
2. Gürsel personel veya okul zamanlayıcısında `Çalıştır` aksiyonunu kullanın.
3. Sistem bir `VerificationRun` kaydı üretir.

API ile:

```http
POST /admin/fleet-tally/schedules/{schedule_id}/run
```

Run başladıktan sonra sonuç `/admin/fleet-tally/runs` ekranından izlenir. Detaylı job sonuçları `/admin/fleet-tally/verifications` ekranında görünür.

### 18.5 Başarı kriterleri

İlk sağlıklı testte aşağıdaki sonuçlar beklenir:

- `/admin/fleet-tally/runs` ekranında yeni run `completed` veya en azından veri varsa `partial` durumuna gelir.
- `total_jobs` dış sistemden çekilen kayıt sayısıyla uyumludur.
- `/admin/fleet-tally/verifications` ekranında plaka, belge/kalem no, job tarihi ve GPS sağlayıcı bilgisi görünür.
- Araç eşleşmesi Stock kaynaklı ise satır kaynağı `Stock` olarak görünür.
- Başarılı kayıtlar `ok`; başarısız kayıtlar anlamlı şekilde `no_device`, `no_location`, `empty_window`, `not_in_zone` veya `error` olarak ayrışır.
- Dış sisteme sonuç gönderimi açıksa `push_status=sent` beklenir.

### 18.6 Hızlı hata ayrıştırma

| Belirti | İlk bakılacak yer |
|---|---|
| Provider sayıları `0` görünüyor | `Stock` araç varlıklarında GPS kimlikleri ve `vehicle-mappings` ekranındaki Stock kaynaklı kayıtlar |
| Connector kayıt dönmüyor | `/admin/fleet-tally/connectors`, credential, tarih aralığı, dış sistem bekleyen kayıt statüsü |
| Run oluşmuyor | `/admin/fleet-tally/schedules`, zamanlayıcı ID, active hours, API logları |
| Çok sayıda `no_device` var | Stock araç GPS varlıkları, plaka normalize değeri, N2 için sağlayıcı external ID |
| Çok sayıda `no_location` var | `/admin/fleet-tally/settings` kaynak modu, GPS provider credential, cihazın gerçekten veri gönderip göndermediği |
| `push_status=failed` var | Connector push method/endpoint/credential ve dış sistem update kabul durumu |

### 18.7 Tekrar deneme

Başarısız bir kayıt provider veya transient veri gecikmesi nedeniyle tekrar denenmek istenirse `/admin/fleet-tally/verifications` ekranından ilgili job seçilerek yeniden deneme yapılır.

API karşılığı:

```http
POST /admin/fleet-tally/verifications/jobs/{job_id}/retry?gps_provider_config_id={gps_provider_id}
```

---

## Hızlı Özet

Fleet Tally’de günlük operasyonun en pratik sırası:

1. `Zamanlayıcılar` ekranından ilgili akışı `Çalıştır`.
2. `Çalışma Geçmişi` ekranında run durumunu izle.
3. `Doğrulama İşleri` ekranında başarısız kayıtları filtrele.
4. `no_device` için araç eşlemesini, `no_location` için GPS sağlayıcıyı, `push_status=failed` için dış sistem bağlayıcısını kontrol et.
5. Gerekirse job’u doğru GPS sağlayıcıyla tekrar dene.
