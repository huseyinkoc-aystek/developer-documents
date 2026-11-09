# Fleet Tally System — Mevcut Durum Analizi ve Enterprise Tasarım Kararları

**Durum:** Analiz Tamamlandı  
**Tarih:** 2026-05-05  
**Versiyon:** 1.0.0

---

## İçindekiler

1. [Mevcut Sistem Genel Bakış](#1-mevcut-sistem-genel-bakış)
2. [Kod Tabanı Analizi](#2-kod-tabanı-analizi)
3. [Domain Kavramları ve Sözlük](#3-domain-kavramları-ve-sözlük)
4. [Tespit Edilen Sorunlar](#4-tespit-edilen-sorunlar)
5. [Enterprise Yeniden Tasarım Kararları](#5-enterprise-yeniden-tasarım-kararları)
6. [Yeni Domain Modeli](#6-yeni-domain-modeli)
7. [Modül Yapısı (AYSTEK Core'da)](#7-modül-yapısı-aystek-corede)
8. [Algoritma Sözleşmesi](#8-algoritma-sözleşmesi)
9. [Bağlantı Sözleşmeleri](#9-bağlantı-sözleşmeleri)
10. [Durum Kodu Yeniden Tanımı](#10-durum-kodu-yeniden-tanımı)

---

## 1. Mevcut Sistem Genel Bakış

### 1.1 Sistemin Amacı

Araç üzerindeki GPS cihazlarından alınan konum verilerini kullanarak, bir araçlı hizmetin (personel servisi, okul servisi vb.) sözleşmede belirtilen lokasyonda, belirtilen zaman aralığında gerçekleşip gerçekleşmediğini **otomatik olarak doğrulamak** ve sonuçları hizmet alan firmanın ERP sistemine iletmektir.

Türkçe tabiriyle: **"çetele kontrol sistemi"** — araç fiilen gitmesi gereken yere gidip gitmediğinin GPS verisiyle kanıtlanması.

### 1.2 Mevcut Teknoloji Yığını

| Bileşen | Teknoloji | Konum |
|---|---|---|
| GPS Sorgu Servisi | PHP / Yii2 Framework | `php/GurselMobisisKurumsal/modules/services/controllers/GpsserviceController.php` |
| Personel Servis Kontrolcüsü | Python 3 / Zeep (SOAP) | `python/GPS_CONTROLLER/PERSONEL/PersonelStartGPSControl.py` |
| Okul Servis Kontrolcüsü | Python 3 / Zeep (SOAP) | `python/GPS_CONTROLLER/SCHOOL/SchootStartGPSControl.py` |
| Veri Kaynağı (ERP) | SAP / SOAP (WSDL) | `sapservice.gurseltur.com.tr` |
| GPS Sağlayıcılar | Trio, Mobiliz, Arvento, N2 | — |

### 1.3 Veri Akışı (Mevcut)

```
SAP WSDL (Gürsel)
    │
    │ OnSiparisKayitListesi / okulHakediskayitListesiServis
    ▼
Python Script (Personel / Okul)
    │
    │ HTTP POST /searchgpslocations
    ▼
PHP GPS Service (GpsserviceController)
    │
    ├── Trio API   (telematik.triomobil.com)
    ├── Mobiliz DB (yerel VehicleGpsLocation tablosu)
    ├── Arvento API
    └── N2 API
    │
    │ { controlStatus, deviceIsPointInCircle, locationsCount, exp }
    ▼
Python Script → sonuçları birleştirir
    │
    │ OnSiparisKayitGuncelle / okulHakedisGuncellemeServis
    ▼
SAP WSDL (Gürsel) — sonuç iletildi
```

---

## 2. Kod Tabanı Analizi

### 2.1 Personel Servis (`PersonelStartGPSControl.py`)

**İş Akışı:**
- SAP'tan son 3 aya ait, GPS kontrolü bekleyen kayıtları çeker (`statu='G'`)
- Her kayıt için `SecilenSaatler` listesindeki zaman aralıklarını döner
- Her zaman aralığı için PHP servisine GPS sorgusu atar
- Araç belirtilen zaman aralığında hedef koordinat çevresindeki yarıçap (YariCap) içinde mi?
- Tüm zaman aralıkları doğrulandıysa `ServisGpsOkay = "X"` set edilir
- Sonuçlar SAP'a geri iletilir

**Kontrol Parametreleri:**
```json
{
  "ZceteleNo": "iş numarası",
  "ZceteleKalemNo": "kalem no",
  "Plaka": "34HHK838",
  "Tarih": "2024-09-05",
  "BasKoor": "41.034868...",   ← hedef lat
  "BitKoor": "29.256762...",   ← hedef lon
  "SecilenSaatler": ["06:45:00-05:30:00-2"],  ← startTime-endTime-routeType
  "KEZamanAraligi": ["2024.09.05 15:00-2024.09.05 17:45"],
  "YariCap": "700",
  "Zgps": "trio"
}
```

**Sonuç Güncelleme Alanları:**
- `ServisKontrol`: "X" (kontrol edildi)
- `ServisGpsOkay`: "X" (tüm zaman aralıkları geçti)
- `GpsDurumu`: sonuç durum kodu
- `Hata`: hata mesajı
- `ZGPSAciklama`: GPS açıklama

### 2.2 Okul Servis (`SchootStartGPSControl.py`)

**Fark:** Her kayıt için **GİRİŞ** ve **ÇIKIŞ** olmak üzere iki ayrı zaman/lokasyon kontrolü yapılır.

```json
{
  "BelgeNo": "4600011612",
  "Plaka": "34LAC242",
  "Tarih": "2025-11-03",
  "BaslangicX": "41.049734",  ← giriş lokasyonu lat
  "BaslangicY": "28.993662",  ← giriş lokasyonu lon
  "BitisX": "41.049734",      ← çıkış lokasyonu lat (genellikle aynı)
  "BitisY": "28.993662",      ← çıkış lokasyonu lon
  "KontrolGirisSaati": "2025.11.03 06:15-2025.11.03 08:15",
  "KontrolCikisSaati": "2025.11.03 15:45-2025.11.03 17:45",
  "YariCap": null              ← staticRadius = 100m kullanılır
}
```

**Sonuç Güncelleme Alanları:**
- `GpsOkey`: "1"=başarılı, "2"=başarısız
- `GpsDurum`: durum kodu
- `KontrolGirisSaatiDurum`: giriş kontrolü sonucu
- `KontrolCikisSaatiDurum`: çıkış kontrolü sonucu
- `ServisKontrol`: "X" (kontrol edildi)
- `ServisOkey`: "X" (sadece başarısız kayıtlara set edilir — DİKKAT: mantık tersine)

### 2.3 PHP GPS Servisi (`GpsserviceController.php`)

**Sorumlulukları:**
1. GPS sağlayıcısından konum verisi çekme (provider'a göre switch-case)
2. Haversine mesafe hesaplama
3. Nokta-çember içi kontrolü (`isPointInCircle`)
4. Çizgi-çember kesişim kontrolü (`isLineIntersectingCircle`) — rota tabanlı kontrol için
5. Tüm gün konum var ama belirtilen saatlerde yok → farklı durum kodu

**GPS Provider Uyarlamaları:**
| Provider | Yöntem | Özellik |
|---|---|---|
| `trio` | Trio REST API | Konum, ID üzerinden |
| `mobiliz` | Yerel DB (VehicleGpsLocation) | UTC-3 dönüşümü yapılıyor |
| `arvento` | Arvento REST API | JSON decode |
| `n2` | N2 REST API | `data` + `count` yapısı |

---

## 3. Domain Kavramları ve Sözlük

### 3.1 Terim Sözlüğü

| Türkçe Terim | İngilizce Karşılığı | Açıklama |
|---|---|---|
| Çetele Kontrol | TallyVerification | Hizmetin GPS verisiyle doğrulanması |
| Kontrol Türü | ControlType | Doğrulama algoritmasının isimlendirilmiş versiyonu |
| İş Kalemi | VerificationJob | Doğrulanacak tek bir hizmet kaydı |
| Doğrulama Sonucu | VerificationResult | Bir iş kaleminin kontrol çıktısı |
| Zaman Penceresi | TimeWindow | Kontrolün yapıldığı başlangıç-bitiş saat aralığı |
| Hedef Bölge | TargetZone | Araçın bulunması gereken coğrafi alan |
| GPS Sağlayıcı | GpsProvider | Konum verisi sağlayan sistem (Trio, Mobiliz, vs.) |
| Dış Sistem Bağlayıcı | ExternalConnector | İş kalemlerinin alındığı ve sonuçların gönderildiği ERP/SAP |
| Yarıçap | RadiusMeters | Nokta etrafındaki kontrol yarıçapı (metre) |
| Rota Tipi | RouteType | 1=Giriş, 2=Çıkış, 0=Herhangi |

### 3.2 Kontrol Türleri (Mevcut)

**`personel_servis_v1`:**
- N adet zaman penceresi (1..n)
- Tek hedef koordinat (BasKoor / BitKoor)
- Yarıçap: kayıttan gelir (YariCap), eksikse varsayılan
- Rota tipi: her zaman penceresinde farklı olabilir (1/2/0)
- Başarı: TÜM zaman pencereleri geçmeli

**`okul_servis_v1`:**
- Sabit 2 zaman penceresi (giriş + çıkış)
- Giriş: BaslangicX/BaslangicY koordinatı
- Çıkış: BitisX/BitisY koordinatı (genellikle aynı)
- Yarıçap: sabit 100m
- Başarı: her iki pencere de geçmeli

---

## 4. Tespit Edilen Sorunlar

### 4.1 Güvenlik Sorunları (Kritik)

| # | Sorun | Konum | Risk |
|---|---|---|---|
| S1 | Kimlik bilgileri kodda sabit (`ps1234`, `\|g1u2r3s4e5l6Portal.-!\|`) | Python dosyaları | Credential leakage |
| S2 | HTTP Basic Auth şifresi kaynak kodda düz metin | `GpsLocations()` | Credential leakage |
| S3 | WSDL password kaynak kodda düz metin | `GetVehicleJobs()` | Credential leakage |
| S4 | `$enableCsrfValidation = false` | PHP Controller | CSRF riski |
| S5 | `$_REQUEST` ile direkt parametre okuma | PHP Controller | Injection riski |

### 4.2 Tasarım Sorunları

| # | Sorun | Açıklama |
|---|---|---|
| D1 | Durum kodu çakışması | `NoDevice="1"` ve `AllLocationIsPointIn="1"` aynı değer, farklı anlam |
| D2 | Mükerrer Haversine implementasyonu | PHP'de `haversineGreatCircleDistanceInMeters` ve `haversineDistance` — aynı formül iki fonksiyon |
| D3 | Kontrol türü kodda gömülü | `personel` ve `okul` akışı iki ayrı Python script — genişletilemez |
| D4 | Provider switch-case spaghetti | PHP'de 4 provider için 400+ satır switch-case |
| D5 | Hata yönetimi tutarsız | Bazen `success=false`, bazen `exp` alanı, bazen sessiz fail |
| D6 | Saat dilimi dönüşümü kodda gömülü | Mobiliz için `-3 saat` PHP içinde, diğerleri için yok |
| D7 | `sleepStatus()` iş mantığı ile zamanlama karışık | Tek fonksiyonda hem zaman kontrolü hem iş başlatma |
| D8 | `willBeCheckedAgain=9` mantığı karmaşık | "Yarın tekrar bak" mantığı kısmen uygulanmış, kısmen atlanmış |

### 4.3 Operasyonel Sorunlar

| # | Sorun | Açıklama |
|---|---|---|
| O1 | 4 ayrı GPS sağlayıcı için tek servis | Yeni sağlayıcı eklemek PHP kodunu değiştiriyor |
| O2 | Yeni firma → yeni Python script | Gürsel dışı firma için yeni dosya + deployment gerekiyor |
| O3 | Görünürlük yok | Hangi işin kontrol edildiği, başarılı/başarısız sayıları izlenemiyor |
| O4 | Retry mekanizması ilkel | `reController == 5` → 5 kez dene, yoksa geç |
| O5 | Tek kiracılı | Tüm yapılandırma tek firmaya özel |

---

## 5. Enterprise Yeniden Tasarım Kararları

### 5.1 Temel Mimari Karar

> **Karar:** Mevcut PHP+Python sistemi AYSTEK Core'a `fleet_tally` domain modülü olarak taşınır. Bu modül `integrations` platformunun üzerine inşa edilir.

**Gerekçe:**
- İş kalemlerini dış sistemden çekip doğrulayıp geri göndermek → bir **entegrasyon** sürecidir
- GPS verisi çekmek → telematics altyapısıyla örtüşür ama bağımsız adaptör katmanı gerektirir
- Doğrulama algoritmaları → `ControlType` kayıt sistemiyle isimlendirilmiş, versiyonlanmış olmalıdır

### 5.2 "ControlType" Mimarisi

**Ana fikir:** Bir `ControlType`, hem **neyin kontrol edileceğini** (algoritma) hem de **nasıl yorumlanacağını** (parametre şeması) tanımlar.

```
ControlType: "okul_servis_v1"
  AlgorithmID:  "point_in_circle_multi_window"
  Schema:
    windows: 2 (sabit)
    radiusMeters: 100 (sabit)
    windowMode: "entry_exit"
    successRule: "all_windows_must_pass"
  
ControlType: "personel_servis_v1"
  AlgorithmID:  "point_in_circle_multi_window"
  Schema:
    windows: dynamic (kayıttan okunur)
    radiusMeters: from_record (YariCap)
    windowMode: "arbitrary"
    successRule: "all_windows_must_pass"

ControlType: "personel_servis_route_v1"
  AlgorithmID:  "route_intersection_circle"
  Schema:
    windows: dynamic
    radiusMeters: from_record
    successRule: "any_segment_must_intersect"
```

Yeni müşteri geldiğinde → yeni `ControlType` kaydı + mevcut algoritma referansı. Algoritma değiştiğinde → yeni algoritma versiyonu + yeni `ControlType`.

### 5.3 Multi-Tenant GPS Provider

Her tenant (hizmet verilen firma) kendi GPS provider konfigürasyonlarını tanımlar:
- Hangi provider'lar aktif
- Her provider için API kimlik bilgileri (encrypted at rest)
- UTC offset / saat dilimi tanımı

### 5.4 Multi-Tenant External Connector

Her tenant kendi ERP/SAP bağlantısını tanımlar:
- Connector tipi: `sap_soap`, `rest_api`, `webhook`, `csv_upload`
- Bağlantı parametreleri (WSDL URL, kimlik bilgileri — encrypted)
- Çekme metodu adı + gönderme metodu adı
- Veri dönüşüm şeması (kayıt alanları → normalized format mapping)

### 5.5 Bağımsız Worker Mimarisi

Mevcut Python while-loop yerine AYSTEK Core'un `jobs` altyapısı kullanılır:
- Zamanlanmış job: her N dakikada bir çalış
- Her tenant + her control type için bağımsız job
- Başarısız iş kalemleri retry queue'ya düşer
- Tam audit trail (ne zaman kontrol edildi, kaç kayıt, başarı/başarısız oranı)

---

## 6. Yeni Domain Modeli

### 6.1 Ana Varlıklar

```
ControlType
├── id: ObjectID
├── tenant_id: ObjectID
├── name: string              ← "okul_servis_v1"
├── display_name: string      ← "Okul Servisi Kontrol v1"
├── algorithm_id: string      ← "point_in_circle_multi_window"
├── algorithm_version: string ← "1.0"
├── schema: bson.Raw          ← algoritma parametreleri (JSON Schema)
├── status: enum              ← active | deprecated | archived
└── timestamps

GpsProviderConfig
├── id: ObjectID
├── tenant_id: ObjectID
├── provider_type: enum       ← trio | mobiliz | arvento | n2 | custom
├── display_name: string
├── base_url: string
├── credentials: encrypted    ← { username, password, api_key, ... }
├── timezone_offset: int      ← dakika cinsinden UTC offset (örn. -180 = UTC-3)
├── is_active: bool
└── timestamps

ExternalConnector
├── id: ObjectID
├── tenant_id: ObjectID
├── connector_type: enum      ← sap_soap | rest_api | webhook | manual
├── display_name: string
├── fetch_config: bson.Raw    ← { wsdl_url, method_name, ... }
├── push_config: bson.Raw     ← { wsdl_url, method_name, ... }
├── field_mapping: bson.Raw   ← dış format → normalized format
├── credentials: encrypted
├── is_active: bool
└── timestamps

VerificationSchedule
├── id: ObjectID
├── tenant_id: ObjectID
├── connector_id: ObjectID    ← hangi dış sistemden iş çekilecek
├── control_type_id: ObjectID ← hangi kontrol türü uygulanacak
├── gps_provider_ids: []OID   ← hangi GPS sağlayıcılar denenecek (sırayla)
├── cron_expr: string         ← "0 4 * * *" (her gece 04:00)
├── lookback_days: int        ← kaç günlük iş çekilecek
├── max_jobs_per_run: int     ← processLimit
├── retry_policy: object      ← max_attempts, backoff_seconds
├── active_hours: object      ← { start: "03:30", end: "23:30" }
├── is_active: bool
└── timestamps

VerificationJob
├── id: ObjectID
├── tenant_id: ObjectID
├── schedule_id: ObjectID
├── control_type_id: ObjectID
├── external_ref: object      ← { source_id, doc_no, line_no, ... }
├── vehicle_plate: string
├── gps_provider: string      ← hangi provider bu job'u kontrol etti
├── job_date: time.Time
├── windows: []TimeWindow
│   ├── window_index: int
│   ├── start_at: time.Time
│   ├── end_at: time.Time
│   ├── target_lat: float64
│   ├── target_lon: float64
│   ├── radius_meters: float64
│   ├── route_type: enum       ← entry | exit | any
│   └── result: WindowResult
│       ├── status: VerificationStatus
│       ├── is_in_zone: bool
│       ├── locations_count: int
│       ├── all_day_locations_count: int
│       └── explanation: string
├── overall_status: VerificationStatus
├── push_status: enum          ← pending | sent | failed
├── push_attempts: int
├── raw_external_payload: bson.Raw   ← kayıt geldığinde dış sistemden gelen ham JSON
├── control_started_at: time.Time
├── control_ended_at: time.Time
└── timestamps

VerificationRun
├── id: ObjectID
├── tenant_id: ObjectID
├── schedule_id: ObjectID
├── started_at: time.Time
├── ended_at: time.Time
├── total_jobs: int
├── passed_jobs: int
├── failed_jobs: int
├── error_jobs: int
├── pending_retry_jobs: int
└── status: enum               ← running | completed | failed
```

### 6.2 Durum Makinesi

```
VerificationStatus (yeniden tanımlanmış, örtüşmesiz):
  "" (empty)    = ok             → Kontrol başarılı, araç bölgedeydi
  "no_device"   = cihaz_yok     → Araçta aktif GPS cihazı bulunamadı
  "no_gsm"      = hat_yok       → GPS cihazı var ama GSM hattı yok
  "no_location" = konum_yok     → Belirtilen saatte ve tüm gün hiç konum yok
  "empty_window" = pencere_bos  → Belirtilen saatte konum yok ama tüm gün var
  "not_in_zone" = bolge_disi    → Konum var ama hedef bölge dışında
  "error"       = hata          → Teknik hata (provider bağlantı sorunu)
  "retry_later" = tekrar_bak    → Veri henüz yok, yarın tekrar kontrol et
  "unknown"     = bilinmiyor    → Belirsiz durum (max retry sonrası)
```

---

## 7. Modül Yapısı (AYSTEK Core'da)

```
go/aystek-core/internal/modules/fleet_tally/
│
├── control_types/                    # Kontrol türü kayıt defteri
│   ├── model/
│   │   └── control_type.go
│   ├── dto/
│   │   └── control_type.go
│   ├── repository/
│   │   ├── interface.go
│   │   └── mongo.go
│   ├── service/
│   │   ├── service.go
│   │   └── service_test.go
│   ├── handler/
│   │   └── admin.go
│   └── routes.go
│
├── gps_providers/                    # GPS sağlayıcı konfigürasyonları
│   ├── model/
│   │   └── gps_provider.go
│   ├── adapters/                     # Her sağlayıcı için interface + impl
│   │   ├── interface.go              # GpsProviderAdapter interface
│   │   ├── trio/
│   │   │   └── adapter.go
│   │   ├── mobiliz/
│   │   │   └── adapter.go
│   │   ├── arvento/
│   │   │   └── adapter.go
│   │   └── n2/
│   │       └── adapter.go
│   ├── registry/
│   │   └── registry.go               # adapter factory: tip → implementation
│   ├── dto/
│   ├── repository/
│   │   ├── interface.go
│   │   └── mongo.go
│   ├── service/
│   │   ├── service.go
│   │   └── service_test.go
│   ├── handler/
│   │   └── admin.go
│   └── routes.go
│
├── connectors/                       # Dış sistem bağlayıcıları
│   ├── model/
│   │   └── connector.go
│   ├── adapters/
│   │   ├── interface.go              # ExternalConnectorAdapter interface
│   │   ├── sap_soap/
│   │   │   └── adapter.go            # Gürsel SAP WSDL uyarlaması
│   │   ├── rest_api/
│   │   │   └── adapter.go
│   │   └── webhook/
│   │       └── adapter.go
│   ├── registry/
│   │   └── registry.go
│   ├── dto/
│   ├── repository/
│   │   ├── interface.go
│   │   └── mongo.go
│   ├── service/
│   │   ├── service.go
│   │   └── service_test.go
│   ├── handler/
│   │   └── admin.go
│   └── routes.go
│
├── verifications/                    # Doğrulama motoru
│   ├── model/
│   │   └── verification.go           # VerificationJob, VerificationRun
│   ├── algorithms/                   # Kontrol algoritmaları
│   │   ├── interface.go              # ControlAlgorithm interface
│   │   ├── registry.go               # algorithmID → ControlAlgorithm
│   │   ├── geo/
│   │   │   ├── haversine.go          # Pure matematik (test edilebilir)
│   │   │   └── haversine_test.go
│   │   ├── point_in_circle/
│   │   │   ├── algorithm.go
│   │   │   └── algorithm_test.go
│   │   └── route_intersection/
│   │       ├── algorithm.go
│   │       └── algorithm_test.go
│   ├── dto/
│   ├── repository/
│   │   ├── interface.go
│   │   └── mongo.go
│   ├── service/
│   │   ├── service.go                # Doğrulama orkestrasyon mantığı
│   │   └── service_test.go
│   ├── handler/
│   │   ├── admin.go                  # Sonuç görüntüleme, manuel tetikleme
│   │   └── internal.go               # Sadece iç worker'ların kullandığı
│   └── routes.go
│
└── schedules/                        # Zamanlanmış doğrulama görevleri
    ├── model/
    │   └── schedule.go
    ├── dto/
    ├── repository/
    │   ├── interface.go
    │   └── mongo.go
    ├── service/
    │   ├── service.go
    │   ├── runner.go                  # Worker: schedule'ı çalıştıran mantık
    │   └── service_test.go
    ├── handler/
    │   └── admin.go
    └── routes.go
```

### 7.1 Frontend Modülleri

```
node/aystek-core-next/app/(dashboard)/admin/fleet-tally/
├── control-types/page.tsx
├── gps-providers/page.tsx
├── connectors/page.tsx
├── verifications/page.tsx           ← sonuç listesi + detay
├── schedules/page.tsx
└── runs/page.tsx                    ← çalışma geçmişi + istatistikler
```

---

## 8. Algoritma Sözleşmesi

### 8.1 ControlAlgorithm Interface (Go)

```go
// ControlAlgorithm, bir zaman penceresi için GPS konumlarını değerlendiren algoritmayı tanımlar.
type ControlAlgorithm interface {
    // AlgorithmID bu algoritmanın benzersiz kayıt adını döner (ör. "point_in_circle_multi_window")
    AlgorithmID() string

    // Evaluate, GPS konumları ve pencere parametreleriyle doğrulama yapar.
    // Her pencere için WindowResult döner.
    Evaluate(ctx context.Context, input EvaluationInput) (EvaluationOutput, error)
}

type EvaluationInput struct {
    Windows       []WindowDef       // kontrol edilecek zaman pencereleri
    Locations     []GpsPoint        // GPS provider'dan gelen konumlar
    SchemaParams  map[string]any    // ControlType.schema'dan gelen parametreler
}

type WindowDef struct {
    Index        int
    StartAt      time.Time
    EndAt        time.Time
    TargetLat    float64
    TargetLon    float64
    RadiusMeters float64
    RouteType    RouteType
}

type EvaluationOutput struct {
    WindowResults []WindowResult
    OverallStatus VerificationStatus
}
```

### 8.2 Kayıtlı Algoritmalar

| AlgorithmID | Açıklama | Mevcut Karşılık |
|---|---|---|
| `point_in_circle_multi_window` | N pencerede, nokta-çember içi kontrolü | Personel + Okul temel kontrol |
| `route_intersection_circle` | Konum nokta değil, rota üzerindeki segment çemberi kesiyor mu | PHP `isLineIntersectingCircle` |
| `all_day_fallback_check` | Belirtilen saatte yoksa tüm gün bakar, boş/dolu farkını bildirir | PHP `EmptyLocation` mantığı |

---

## 9. Bağlantı Sözleşmeleri

### 9.1 ExternalConnectorAdapter Interface (Go)

```go
type ExternalConnectorAdapter interface {
    // ConnectorType döner (ör. "sap_soap")
    ConnectorType() string

    // FetchJobs, dış sistemden doğrulanacak iş kalemlerini çeker.
    // Dönen kayıtlar normalize edilmiş VerificationJobSpec formatına dönüştürülür.
    FetchJobs(ctx context.Context, cfg ConnectorConfig, params FetchParams) ([]RawJobRecord, error)

    // PushResults, doğrulama sonuçlarını dış sistemin beklediği formata çevirir ve gönderir.
    PushResults(ctx context.Context, cfg ConnectorConfig, results []JobResult) error
}

type FetchParams struct {
    FromDate     time.Time
    ToDate       time.Time
    GpsModelFilter string  // ör. "trio" — sadece bu provider'ın kayıtları
    Limit        int
}
```

### 9.2 GpsProviderAdapter Interface (Go)

```go
type GpsProviderAdapter interface {
    // ProviderType döner (ör. "trio")
    ProviderType() string

    // FetchLocations, belirtilen plaka ve zaman aralığı için GPS konumlarını döner.
    FetchLocations(ctx context.Context, cfg ProviderConfig, req LocationRequest) ([]GpsPoint, error)
}

type LocationRequest struct {
    Plate     string
    StartAt   time.Time
    EndAt     time.Time
}

type GpsPoint struct {
    Lat       float64
    Lon       float64
    RecordedAt time.Time
    Speed     *float64   // opsiyonel
    Heading   *float64   // opsiyonel
}
```

---

## 10. Durum Kodu Yeniden Tanımı

Mevcut sistemde durum kodları string sayılar ve çakışan anlamlar içeriyor. Yeni sistemde Go `const` olarak tanımlanır:

```go
type VerificationStatus string

const (
    StatusOK          VerificationStatus = ""           // başarılı
    StatusNoDevice    VerificationStatus = "no_device"  // cihaz yok
    StatusNoGsm       VerificationStatus = "no_gsm"     // hat yok
    StatusNoLocation  VerificationStatus = "no_location" // hiç konum yok
    StatusEmptyWindow VerificationStatus = "empty_window" // pencerede konum yok
    StatusNotInZone   VerificationStatus = "not_in_zone" // bölge dışı
    StatusError       VerificationStatus = "error"       // teknik hata
    StatusRetryLater  VerificationStatus = "retry_later" // yarın tekrar bak
    StatusUnknown     VerificationStatus = "unknown"     // belirsiz
)
```

**Geriye Dönük Uyumluluk (SAP SOAP iletiminde):**

Her `ExternalConnectorAdapter.PushResults()`, `VerificationStatus` değerlerini kendi formatına çevirir. Gürsel SAP adaptörü için dönüşüm:

| Yeni Kod | Gürsel SAP Kodu |
|---|---|
| `""` (ok) | `""` veya `"5"` (okul) |
| `"no_device"` | `"1"` |
| `"no_gsm"` | `"7"` |
| `"empty_window"` | `"2"` |
| `"not_in_zone"` | `"3"` |
| `"retry_later"` | `"9"` |
| `"error"` / `"unknown"` | `"4"` |

Bu dönüşüm tablosu `connector.field_mapping` içinde yapılandırılır — kod değişikliği gerekmez.

---

*Bu analiz dokümanı tamamlandığında Master Plan dokümanına bakınız:*  
`FLEET_TALLY_SYSTEM_MASTER_PLAN_TR.md`

*aystek-core üzerinde nasıl geliştirme yapılacağı için bkz.:*  
`FLEET_TALLY_AYSTEK_CORE_IMPL_GUIDE_TR.md`
