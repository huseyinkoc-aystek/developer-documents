# Fleet Tally System — AYSTEK Core'a Taşıma Master Planı

**Durum:** Planlanıyor  
**Tarih:** 2026-05-05  
**Versiyon:** 1.0.0  
**Analiz Dokümanı:** `FLEET_TALLY_SYSTEM_ANALYSIS_TR.md`

---

## İçindekiler

1. [Hedef ve Kapsam](#1-hedef-ve-kapsam)
2. [Modül Dağılımı](#2-modül-dağılımı)
3. [Wave Planı](#3-wave-planı)
4. [Wave 1 — Temel Altyapı](#4-wave-1--temel-altyapı)
5. [Wave 2 — GPS Sağlayıcı Adaptörleri](#5-wave-2--gps-sağlayıcı-adaptörleri)
6. [Wave 3 — Dış Sistem Bağlayıcıları](#6-wave-3--dış-sistem-bağlayıcıları)
7. [Wave 4 — Doğrulama Motoru](#7-wave-4--doğrulama-motoru)
8. [Wave 5 — Zamanlayıcı ve Worker](#8-wave-5--zamanlayıcı-ve-worker)
9. [Wave 6 — Admin UI](#9-wave-6--admin-ui)
10. [Wave 7 — Multi-Tenant ve Genişletme](#10-wave-7--multi-tenant-ve-genişletme)
11. [Seed ve i18n Gereksinimleri](#11-seed-ve-i18n-gereksinimleri)
12. [Test Stratejisi](#12-test-stratejisi)
13. [Gürsel SAP Geçiş Protokolü](#13-gürsel-sap-geçiş-protokolü)
14. [Bağımlılık Haritası](#14-bağımlılık-haritası)

---

## 1. Hedef ve Kapsam

### 1.1 Hedef

Mevcut PHP+Python GPS doğrulama sistemini AYSTEK Core'a enterprise standartlarda taşımak:

- Kodda gömülü olmayan, **isimlendirilen ve versiyonlanan kontrol türleri**
- Birden fazla firmaya hizmet verebilen **multi-tenant mimari**
- Her firma için farklı GPS sağlayıcı ve ERP bağlantısı destekleyen **adaptör katmanı**
- Tam **audit trail** ve operasyonel görünürlük
- Güvenli kimlik bilgisi yönetimi (encrypted at rest)
- AYSTEK Core standartlarına %100 uyum (DDD, 6 katman, i18n, seed, test)

### 1.2 Kapsam Dışı (Bu Plana Dahil Değil)

- Gerçek zamanlı harita görüntüleme (telematics/telemetry modülü kapsamı)
- Araç cihaz yönetimi (telematics/devices kapsamı)
- SAP entegrasyon platformu yönetimi (ayrı entegrasyon modülü kapsamı)

---

## 2. Modül Dağılımı

```
fleet_tally/
├── control_types        ← Kontrol türü kayıt defteri
├── gps_providers        ← GPS sağlayıcı konfigürasyonları + adaptörler
├── connectors           ← Dış ERP/SAP bağlayıcıları + adaptörler  
├── verifications        ← Doğrulama motoru (algoritmalar + job yönetimi)
└── schedules            ← Zamanlı çalıştırma + worker
```

**AYSTEK Core altyapısına bağımlılıklar:**
- `internal/platform/jobs` — zamanlanmış görev altyapısı
- `internal/platform/queue` — retry kuyruğu
- `internal/platform/config` — tenant konfigürasyonu
- `internal/seeds/` — seed kayıtları (ControlType, menu, translation)

---

## 3. Wave Planı

| Wave | Başlık | Ön Koşul | Tahmini Efor |
|---|---|---|---|
| **Wave 1** | Temel altyapı: ControlType + domain model | — | 3-4 gün |
| **Wave 2** | GPS Sağlayıcı Adaptörleri | Wave 1 | 3-4 gün |
| **Wave 3** | Dış Sistem Bağlayıcıları (SAP SOAP) | Wave 1 | 3-4 gün |
| **Wave 4** | Doğrulama Motoru (algoritmalar + job işleme) | Wave 1, 2, 3 | 4-5 gün |
| **Wave 5** | Zamanlayıcı + Background Worker | Wave 4 | 2-3 gün |
| **Wave 6** | Admin UI (Next.js sayfaları) | Wave 1-5 backend | 4-5 gün |
| **Wave 7** | Multi-tenant genişletme + yeni provider/connector | Wave 1-6 | Sürekli |

---

## 4. Wave 1 — Temel Altyapı

### Amaç
Domain modellerini, `ControlType` kayıt defterini ve temel CRUD altyapısını kurmak.

### Görevler

#### W1.1 — Domain Modelleri (Go)

**Dosya:** `internal/modules/fleet_tally/control_types/model/control_type.go`

```go
type ControlTypeStatus string

const (
    ControlTypeStatusActive     ControlTypeStatus = "active"
    ControlTypeStatusDeprecated ControlTypeStatus = "deprecated"
    ControlTypeStatusArchived   ControlTypeStatus = "archived"
)

type ControlType struct {
    ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    TenantID         primitive.ObjectID `bson:"tenant_id" json:"tenant_id"`
    Name             string             `bson:"name" json:"name"`                   // "okul_servis_v1"
    DisplayName      string             `bson:"display_name" json:"display_name"`
    AlgorithmID      string             `bson:"algorithm_id" json:"algorithm_id"`   // "point_in_circle_multi_window"
    AlgorithmVersion string             `bson:"algorithm_version" json:"algorithm_version"`
    Schema           bson.Raw           `bson:"schema" json:"schema"`               // algoritma parametreleri
    Status           ControlTypeStatus  `bson:"status" json:"status"`
    CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}
```

**Diğer model dosyaları:**
- `gps_providers/model/gps_provider.go` — `GpsProviderConfig`
- `connectors/model/connector.go` — `ExternalConnector`
- `verifications/model/verification.go` — `VerificationJob`, `VerificationRun`, `TimeWindow`, `WindowResult`
- `schedules/model/schedule.go` — `VerificationSchedule`

**Domain model detayları:** `FLEET_TALLY_SYSTEM_ANALYSIS_TR.md` Bölüm 6'ya bakınız.

#### W1.2 — Repository Katmanı

Her alt modül için standart repository:

**Interface gereksinimleri — ControlType:**
```go
type Repository interface {
    Create(ctx, *model.ControlType) (*model.ControlType, error)
    GetByID(ctx, id) (*model.ControlType, error)
    GetByName(ctx, tenantID, name) (*model.ControlType, error)
    List(ctx, tenantID, status) ([]*model.ControlType, error)
    Update(ctx, *model.ControlType) error
    Delete(ctx, id) error
    CreateIndexes(ctx) error
}
```

**MongoDB Index'leri:**
- `control_types`: `{ tenant_id: 1, name: 1 }` UNIQUE
- `gps_provider_configs`: `{ tenant_id: 1, provider_type: 1, is_active: 1 }`
- `external_connectors`: `{ tenant_id: 1, connector_type: 1, is_active: 1 }`
- `verification_jobs`: `{ tenant_id: 1, schedule_id: 1, job_date: -1, overall_status: 1 }`
- `verification_jobs`: `{ vehicle_plate: 1, job_date: -1 }` (sorgu index)
- `verification_runs`: `{ tenant_id: 1, schedule_id: 1, started_at: -1 }`

#### W1.3 — Service Katmanı (CRUD)

Her alt modül için:
- `Create`, `Update`, `Delete`, `GetByID`, `List` işlemleri
- Business validasyonlar (isim tekrarı, aktif kayıt silme koruması)

#### W1.4 — Handler ve Routes

`/admin/fleet-tally/control-types` group:

```
POST   /admin/fleet-tally/control-types
GET    /admin/fleet-tally/control-types
GET    /admin/fleet-tally/control-types/:id
PUT    /admin/fleet-tally/control-types/:id
DELETE /admin/fleet-tally/control-types/:id
GET    /admin/fleet-tally/control-types/algorithms   ← kayıtlı algoritmaları listele
```

**Aynı pattern** gps-providers ve connectors için de.

#### W1.5 — Seed (ControlType Ön Tanımlı Kayıtlar)

`internal/seeds/fleet_tally_seed.go`:

```go
// Gürsel Turizm için iki hazır kontrol türü seed'i
var defaultControlTypes = []model.ControlType{
    {
        Name:             "personel_servis_v1",
        DisplayName:      "Personel Servis Kontrolü v1",
        AlgorithmID:      "point_in_circle_multi_window",
        AlgorithmVersion: "1.0",
        Schema: bson.Raw(`{
            "window_mode": "arbitrary",
            "radius_source": "from_record",
            "default_radius_meters": 700,
            "success_rule": "all_windows_must_pass",
            "max_retry_attempts": 5,
            "retry_delay_seconds": 10,
            "retry_later_for_today": true
        }`),
        Status: model.ControlTypeStatusActive,
    },
    {
        Name:             "okul_servis_v1",
        DisplayName:      "Okul Servisi Kontrolü v1",
        AlgorithmID:      "point_in_circle_multi_window",
        AlgorithmVersion: "1.0",
        Schema: bson.Raw(`{
            "window_mode": "entry_exit",
            "radius_source": "fixed",
            "fixed_radius_meters": 100,
            "success_rule": "all_windows_must_pass",
            "max_retry_attempts": 5,
            "retry_delay_seconds": 10,
            "retry_later_for_today": true
        }`),
        Status: model.ControlTypeStatusActive,
    },
}
```

#### W1.6 — i18n Çeviriler

`internal/seeds/translation_seed.go` — 28 dil için eklenecek anahtar örnekleri:

```
fleet_tally.module_name
fleet_tally.control_type.create
fleet_tally.control_type.update
fleet_tally.control_type.list
fleet_tally.control_type.field.name
fleet_tally.control_type.field.algorithm_id
fleet_tally.control_type.field.status
fleet_tally.status.ok
fleet_tally.status.no_device
fleet_tally.status.no_gsm
fleet_tally.status.no_location
fleet_tally.status.empty_window
fleet_tally.status.not_in_zone
fleet_tally.status.retry_later
fleet_tally.status.error
fleet_tally.status.unknown
```

#### W1.7 — Menu Seed

`internal/seeds/menu_seed.go`:

```go
// Ana menü: Fleet Tally
{
    Key:      "fleet_tally",
    TitleKey: "fleet_tally.module_name",
    Icon:     "CheckSquare",
    Path:     "/admin/fleet-tally",
    Children: [
        { Key: "fleet_tally.control_types", Path: "/admin/fleet-tally/control-types" },
        { Key: "fleet_tally.gps_providers",  Path: "/admin/fleet-tally/gps-providers" },
        { Key: "fleet_tally.connectors",     Path: "/admin/fleet-tally/connectors" },
        { Key: "fleet_tally.schedules",      Path: "/admin/fleet-tally/schedules" },
        { Key: "fleet_tally.verifications",  Path: "/admin/fleet-tally/verifications" },
        { Key: "fleet_tally.runs",           Path: "/admin/fleet-tally/runs" },
    ],
}
```

#### W1.8 — Wave 1 Tamamlanma Kriteri

- [ ] Tüm 5 alt modül için model dosyaları oluşturuldu
- [ ] Tüm repository interface ve mongo implementasyonları yazıldı
- [ ] `CreateIndexes` tüm collection'lar için uygulandı
- [ ] ControlType CRUD API çalışıyor (Postman ile doğrulandı)
- [ ] GpsProvider + Connector CRUD API çalışıyor
- [ ] Seed dosyaları oluşturuldu ve `seeds/runner.go`'a eklendi
- [ ] 28 dil çevirisi translation_seed'e eklendi
- [ ] Menu seed eklendi
- [ ] Unit testler geçiyor (`go test ./internal/modules/fleet_tally/...`)

---

## 5. Wave 2 — GPS Sağlayıcı Adaptörleri

### Amaç
Her GPS sağlayıcısı için bağımsız, test edilebilir adaptör implementasyonları yazmak.

### Görevler

#### W2.1 — GpsProviderAdapter Interface

`internal/modules/fleet_tally/gps_providers/adapters/interface.go`:

```go
type GpsProviderAdapter interface {
    ProviderType() string
    FetchLocations(ctx context.Context, cfg ProviderConfig, req LocationRequest) ([]GpsPoint, error)
    TestConnection(ctx context.Context, cfg ProviderConfig) error
}

type ProviderConfig struct {
    BaseURL         string
    Credentials     map[string]string  // decrypted at call time
    TimezoneOffset  int                // dakika cinsinden UTC offset
    ExtraParams     map[string]any
}

type LocationRequest struct {
    Plate   string
    StartAt time.Time
    EndAt   time.Time
}

type GpsPoint struct {
    Lat        float64
    Lon        float64
    RecordedAt time.Time
    Speed      *float64
    Heading    *float64
}
```

#### W2.2 — Trio Adaptörü

`adapters/trio/adapter.go`:

- Trio REST API'a HTTP GET isteği
- `telematik.triomobil.com/fleet/v1/...` endpoint
- Token yönetimi (TrioToken cache — TTL'li)
- Yanıtı `[]GpsPoint`'e normalize et
- UTC→yerel saat dönüşümü adaptörde yapılır (artık dışarıda değil)

#### W2.3 — Mobiliz Adaptörü

`adapters/mobiliz/adapter.go`:

- Doğrudan DB sorgusu değil → Mobiliz'in kendi REST API'si üzerinden
- Mevcut `VehicleGpsLocation` PHP modeli → AYSTEK Core'da HTTP call'a dönüşür
- UTC-3 dönüşümü adaptör içinde

#### W2.4 — Arvento Adaptörü

`adapters/arvento/adapter.go`:

- Arvento REST API çağrısı
- `arventoTarihAralikliKonumlarIz` PHP metodundaki mantığı taşı
- Error response (`locationsData.Error`) yönetimi

#### W2.5 — N2 Adaptörü

`adapters/n2/adapter.go`:

- N2 REST API çağrısı
- `data` + `count` response yapısı normalize edilir

#### W2.6 — Provider Registry

`adapters/registry/registry.go`:

```go
type Registry struct {
    adapters map[string]GpsProviderAdapter
}

func NewRegistry() *Registry {
    r := &Registry{adapters: make(map[string]GpsProviderAdapter)}
    r.Register(trio.New())
    r.Register(mobiliz.New())
    r.Register(arvento.New())
    r.Register(n2.New())
    return r
}

func (r *Registry) Get(providerType string) (GpsProviderAdapter, error)
func (r *Registry) Register(adapter GpsProviderAdapter)
func (r *Registry) ListAvailable() []string
```

#### W2.7 — Admin Endpoint: Bağlantı Testi

```
POST /admin/fleet-tally/gps-providers/:id/test-connection
```

→ `adapter.TestConnection()` çağırır, sonucu döner.

#### W2.8 — Wave 2 Tamamlanma Kriteri

- [ ] 4 provider için adaptör implementasyonları tamamlandı
- [ ] Her adaptörün unit testi var (mock HTTP server kullanılarak)
- [ ] Registry factory çalışıyor
- [ ] `TestConnection` endpoint'i çalışıyor
- [ ] Kimlik bilgileri encrypted/decrypted cycle doğrulandı

---

## 6. Wave 3 — Dış Sistem Bağlayıcıları

### Amaç
SAP SOAP ve diğer ERP sistemleriyle veri alışverişini adaptör pattern ile soyutlamak.

### Görevler

#### W3.1 — ExternalConnectorAdapter Interface

`connectors/adapters/interface.go`:

```go
type ExternalConnectorAdapter interface {
    ConnectorType() string
    FetchJobs(ctx context.Context, cfg ConnectorConfig, params FetchParams) ([]RawJobRecord, error)
    PushResults(ctx context.Context, cfg ConnectorConfig, results []JobResult) error
    TestConnection(ctx context.Context, cfg ConnectorConfig) error
}

type ConnectorConfig struct {
    FetchConfig  map[string]any     // { "wsdl_url": ..., "method_name": ... }
    PushConfig   map[string]any
    Credentials  map[string]string  // decrypted
    FieldMapping map[string]string  // normalized field adı → dış alan adı
}

type FetchParams struct {
    FromDate        time.Time
    ToDate          time.Time
    GpsModelFilter  string   // ör. "trio"
    Limit           int
}

type RawJobRecord struct {
    ExternalID    string
    RawPayload    map[string]any  // ham JSON (daha sonra normalize edilecek)
}

type JobResult struct {
    ExternalID    string
    Fields        map[string]any  // normalized sonuç alanları
}
```

#### W3.2 — SAP SOAP Adaptörü (Gürsel Servis Tipi)

`adapters/sap_soap/adapter.go`:

Mevcut Python SOAP çağrılarının Go karşılıkları:

**FetchJobs:**
- SOAP method: `OnSiparisKayitListesi` (personel) veya `okulHakediskayitListesiServis` (okul)
- Parametreler: `ConnectorConfig.FetchConfig`'den okunur (sabit kodlanmaz)
- Yanıt JSON string → `[]RawJobRecord`

**PushResults:**
- SOAP method: `OnSiparisKayitGuncelle` veya `okulHakedisGuncellemeServis`
- `FieldMapping` ile `VerificationStatus` → SAP koduna çevrilir
- Batch gönderim

**Go SOAP kütüphanesi:** `github.com/tiaguinho/gosoap` veya `github.com/hooklift/gowsdl` — hangisi daha güncel ise

#### W3.3 — REST API Adaptörü (Genel Amaçlı)

`adapters/rest_api/adapter.go`:

- `FetchConfig.endpoint`, `FetchConfig.method` (GET/POST)
- Pagination desteği (`FetchConfig.pagination_type`: "offset" | "cursor" | "none")
- JSON response → `[]RawJobRecord`

#### W3.4 — Field Mapping Engine

`connectors/service/mapper.go`:

Ham `RawJobRecord` → normalize edilmiş `VerificationJobSpec`:

```go
type VerificationJobSpec struct {
    VehiclePlate  string
    JobDate       time.Time
    Windows       []WindowSpec
    ExternalRef   map[string]string   // orijinal alan değerleri (geriye dönük raporlama)
}

type WindowSpec struct {
    StartAt      time.Time
    EndAt        time.Time
    TargetLat    float64
    TargetLon    float64
    RadiusMeters *float64   // nil → ControlType.Schema'dan al
    RouteType    RouteType
}
```

`FieldMapping` JSON'u:
```json
{
    "vehicle_plate": "Plaka",
    "job_date": "Tarih",
    "windows_source": "KEZamanAraligi",
    "window_lat": "BasKoor",
    "window_lon": "BitKoor",
    "radius": "YariCap"
}
```

#### W3.5 — Wave 3 Tamamlanma Kriteri

- [ ] SAP SOAP adaptörü: FetchJobs ve PushResults çalışıyor (Gürsel test ortamında doğrulandı)
- [ ] REST adaptörü: generic GET/POST fetch çalışıyor
- [ ] Field mapping engine unit testleri geçiyor
- [ ] Connector Registry çalışıyor
- [ ] `TestConnection` endpoint'i çalışıyor

---

## 7. Wave 4 — Doğrulama Motoru

### Amaç
Kontrol algoritmalarını ve ana doğrulama orkestrasyon mantığını uygulamak.

### Görevler

#### W4.1 — Haversine Geometri Kütüphanesi

`verifications/algorithms/geo/haversine.go`:

```go
const earthRadiusMeters = 6371000.0

// DistanceMeters iki GPS noktası arasındaki mesafeyi metre cinsinden hesaplar.
func DistanceMeters(lat1, lon1, lat2, lon2 float64) float64

// IsPointInCircle bir noktanın çember içinde olup olmadığını kontrol eder.
func IsPointInCircle(pointLat, pointLon, centerLat, centerLon, radiusMeters float64) bool

// IsSegmentIntersectingCircle iki nokta arasındaki doğru parçasının çemberi kesip kesmediğini kontrol eder.
func IsSegmentIntersectingCircle(p1Lat, p1Lon, p2Lat, p2Lon, centerLat, centerLon, radiusMeters float64) bool
```

**Test kapsamı hedefi: %100** (saf matematik fonksiyonlar)

#### W4.2 — ControlAlgorithm Interface ve Registry

`verifications/algorithms/interface.go` + `registry.go`:

```go
type ControlAlgorithm interface {
    AlgorithmID() string
    Evaluate(ctx context.Context, input EvaluationInput) (EvaluationOutput, error)
}

// Global registry — main.go'da initialize edilir
type AlgorithmRegistry struct { ... }
func (r *AlgorithmRegistry) Register(a ControlAlgorithm)
func (r *AlgorithmRegistry) Get(id string) (ControlAlgorithm, error)
```

#### W4.3 — PointInCircleMultiWindow Algoritması

`verifications/algorithms/point_in_circle/algorithm.go`:

Mevcut sistemin temel mantığını Go'ya taşır:

**Personel mod (`window_mode: "arbitrary"`):**
- Her zaman penceresi için: GPS provider'dan konumları çek → nokta-çember kontrolü
- Pencerede konum yoksa → tüm gün kontrol (`all_day_fallback_check`)
- `retry_later_for_today: true` → bugünün kayıtları için `StatusRetryLater` döner
- Max 5 retry → `StatusUnknown`

**Okul mod (`window_mode: "entry_exit"`):**
- Pencere 0: BaslangicX/Y koordinatı (giriş)
- Pencere 1: BitisX/Y koordinatı (çıkış)
- `radius_source: "fixed"` → `fixed_radius_meters: 100`

#### W4.4 — RouteIntersectionCircle Algoritması

`verifications/algorithms/route_intersection/algorithm.go`:

- PHP `isLineIntersectingCircle` mantığını Go'ya taşır
- Ardışık GPS noktaları arasındaki her segment için çember kesişim kontrolü
- Daha hassas — nokta atlamaları bile yakalanır

#### W4.5 — Verification Service (Orkestrasyon)

`verifications/service/service.go`:

```go
type Service struct {
    jobRepo           repository.Repository
    runRepo           RunRepository
    controlTypeRepo   controltypes.Repository
    gpsProviderRepo   gpsproviders.Repository
    gpsProviderReg    *gpsproviders.Registry
    algorithmReg      *algorithms.AlgorithmRegistry
}

// ProcessJob bir iş kalemini kontrol eder ve sonucu kaydeder.
// Bu fonksiyon idempotent'tir — aynı job tekrar çalışırsa sonuç güncellenir.
func (s *Service) ProcessJob(ctx context.Context, job *model.VerificationJob, ctrlType *controltypes.Model, providerCfg *gpsproviders.Model) error

// ProcessBatch bir schedule için birden fazla job'u sırayla işler.
func (s *Service) ProcessBatch(ctx context.Context, jobs []*model.VerificationJob, schedule *schedules.Model) error

// RetryPendingJobs "retry_later" statüsündeki job'ları bugün yeniden kontrol eder.
func (s *Service) RetryPendingJobs(ctx context.Context, tenantID primitive.ObjectID) error
```

**İş akışı (ProcessJob):**

```
1. ControlType.AlgorithmID → AlgorithmRegistry.Get()
2. ControlType.Schema → EvaluationInput.SchemaParams
3. Her TimeWindow için:
   a. GpsProviderRegistry.Get(providerType)
   b. adapter.FetchLocations(plate, window.StartAt, window.EndAt)
   c. algorithm.Evaluate(input)
   d. window.Result = WindowResult{ status, isInZone, count, exp }
4. OverallStatus hesapla (ControlType.Schema.success_rule'a göre)
5. job.ControlStartedAt / ControlEndedAt güncelle
6. Repository.Update(job)
```

#### W4.6 — Wave 4 Tamamlanma Kriteri

- [ ] Haversine fonksiyonları test edildi (%100 coverage)
- [ ] `PointInCircleMultiWindow` algoritması birim testlerden geçiyor (mevcut PHP/Python davranışıyla karşılaştırmalı)
- [ ] `RouteIntersectionCircle` algoritması birim testlerden geçiyor
- [ ] `ProcessJob` entegrasyon testi: mock GPS provider + mock connector ile uçtan uca çalışıyor
- [ ] `ProcessBatch` 100 kayıtlık yükle test edildi (< 5 saniye beklenti)

---

## 8. Wave 5 — Zamanlayıcı ve Worker

### Amaç
Arka planda sürekli çalışan, tenant'a özel zamanlanmış doğrulama görevlerini yönetmek.

### Görevler

#### W5.1 — Schedule Runner

`schedules/service/runner.go`:

```go
type Runner struct {
    scheduleRepo      repository.Repository
    connectorSvc      *connectors.Service
    verificationSvc   *verifications.Service
    runRepo           verifications.RunRepository
}

// RunSchedule bir schedule'ı çalıştırır:
// 1. ExternalConnector üzerinden iş kalemlerini çek
// 2. FieldMapper ile normalize et → []VerificationJob
// 3. GPS provider sırasına göre her job için ProcessBatch çalıştır
// 4. VerificationRun kaydı oluştur (istatistikler)
// 5. Sonuçları ExternalConnector üzerinden geri gönder
func (r *Runner) RunSchedule(ctx context.Context, schedule *model.VerificationSchedule) (*verifications.VerificationRun, error)
```

#### W5.2 — AYSTEK Jobs Entegrasyonu

`internal/platform/jobs` ile entegrasyon:

```go
// Fleet tally job tanımı — main bootstrap'ta register edilir
type FleetTallyJob struct {
    runner *schedules.Runner
}

func (j *FleetTallyJob) Name() string { return "fleet_tally.run_schedule" }
func (j *FleetTallyJob) Run(ctx context.Context, payload jobs.Payload) error {
    scheduleID := payload.MustString("schedule_id")
    // ...
}
```

Her aktif `VerificationSchedule.CronExpr` için job registery'e kayıt.

**Active Hours kontrolü:** Schedule'ın `active_hours.start`/`end` aralığı dışındaysa job skip edilir — `sleepStatus()` Python mantığının yerini alır.

#### W5.3 — Retry Queue

`StatusRetryLater` ve `StatusError` olan job'lar için retry mekanizması:

- `VerificationJob.push_status == "failed"` → retry queue'ya düş
- `RetryPolicy.max_attempts` aşıldıysa `StatusUnknown` set et, log'a yaz
- Exponential backoff (`retry_delay_seconds * 2^attempt`)

#### W5.4 — Manual Trigger API

```
POST /admin/fleet-tally/schedules/:id/run-now
```

→ `runner.RunSchedule()` çağırır, `VerificationRun` ID döner.

#### W5.5 — Wave 5 Tamamlanma Kriteri

- [ ] Schedule runner uçtan uca çalışıyor (test tenant ile)
- [ ] AYSTEK jobs entegrasyonu tamamlandı
- [ ] Active hours kontrolü çalışıyor
- [ ] Retry mekanizması: `retry_later` job'ları ertesi gün yeniden kontrol ediliyor
- [ ] `run-now` endpoint'i çalışıyor
- [ ] VerificationRun istatistikleri doğru kaydediliyor

---

## 9. Wave 6 — Admin UI

### Amaç
Yöneticilerin kontrol türlerini, provider'ları, bağlayıcıları ve doğrulama sonuçlarını yönetip izleyebileceği Next.js sayfaları.

### Görevler

#### W6.1 — Control Types Sayfası

`app/(dashboard)/admin/fleet-tally/control-types/page.tsx`:

- Liste: isim, algoritma, durum, son güncelleme
- Detay/Düzenleme: schema JSON editörü (güvenli, şema doğrulama ile)
- Yeni ControlType formu
- Durum badge: active/deprecated/archived

#### W6.2 — GPS Providers Sayfası

`app/(dashboard)/admin/fleet-tally/gps-providers/page.tsx`:

- Provider listesi + aktif/pasif toggle
- Kimlik bilgisi formu (password alanları masked)
- "Bağlantıyı Test Et" butonu → anlık sonuç

#### W6.3 — Connectors Sayfası

`app/(dashboard)/admin/fleet-tally/connectors/page.tsx`:

- Connector listesi + bağlantı türü badge
- FetchConfig / PushConfig JSON form alanları
- FieldMapping görsel editörü (sol: sistem alanları, sağ: dış alan adları)
- "Bağlantıyı Test Et" butonu

#### W6.4 — Schedules Sayfası

`app/(dashboard)/admin/fleet-tally/schedules/page.tsx`:

- Schedule listesi + son çalışma zamanı + durum
- Yeni schedule: connector + control type + provider seçimi + cron
- "Şimdi Çalıştır" butonu

#### W6.5 — Verifications (Sonuçlar) Sayfası

`app/(dashboard)/admin/fleet-tally/verifications/page.tsx`:

- Tarih aralığı, plaka, durum filtresi
- Tablo: belge no, plaka, tarih, pencere sonuçları, genel durum
- Satır detayı genişletme: her pencere için lokasyon sayısı, açıklama
- Durum renk kodlaması

#### W6.6 — Runs (Çalışma Geçmişi) Sayfası

`app/(dashboard)/admin/fleet-tally/runs/page.tsx`:

- Her VerificationRun satırı: başlangıç, bitiş, toplam/geçen/kalan/hata sayıları
- İlerleme çubuğu (başarı oranı)
- Log drill-down (çalışma içindeki job'lar)

#### W6.7 — Wave 6 Tamamlanma Kriteri

- [ ] Tüm 6 sayfa geliştirme ortamında açılıyor
- [ ] ControlType CRUD UI çalışıyor
- [ ] GPS Provider CRUD + test-connection UI çalışıyor
- [ ] Connector CRUD + field mapping UI çalışıyor
- [ ] Verifications listesi tarih+plaka+durum filtresiyle çalışıyor
- [ ] Runs sayfası istatistikler gösteriyor
- [ ] Tüm metinler i18n anahtarlarından geliyor (`t(...)`)

---

## 10. Wave 7 — Multi-Tenant ve Genişletme

### Amaç
Platformun başka firmalara açılması için gereken genişletme ve dokümantasyon.

### Görevler

#### W7.1 — Yeni GPS Provider Ekleme Rehberi

`aystek-core-dev-docs/user-guide/06-modules/modules/fleet_tally.md` içine:

```
Yeni GPS Provider Ekleme:
1. adapters/<provider_name>/adapter.go oluştur — GpsProviderAdapter interface'ini uygula
2. Registry'e kaydet (adapters/registry/registry.go)
3. Seed: GpsProviderConfig için UI'da seçenek belirir
4. Admin'den tenant konfigürasyon kaydı oluştur
```

#### W7.2 — Yeni External Connector Ekleme Rehberi

```
Yeni Firma Ekleme:
1. Firma tenant'ı oluştur
2. ExternalConnector kaydı: connector_type + kimlik bilgileri + field_mapping
3. GpsProviderConfig kaydı: kullanılan GPS sağlayıcılar
4. ControlType: mevcut algoritmalardan uygun olanı seç veya yeni kayıt oluştur
5. VerificationSchedule: hangi sıklıkta çalışacak
→ Tek satır kod değişikliği gerekmez
```

#### W7.3 — Yeni Kontrol Algoritması Ekleme Rehberi

```
Yeni Algoritma Ekleme:
1. algorithms/<algorithm_name>/algorithm.go oluştur — ControlAlgorithm interface'ini uygula
2. AlgorithmRegistry'e kaydet (registry.go)
3. ControlType kaydında algorithm_id olarak kullan
4. Unit test ekle
→ Mevcut kayıtlar etkilenmez
```

---

## 11. Seed ve i18n Gereksinimleri

### 11.1 Seed Dosyaları

| Dosya | İçerik |
|---|---|
| `internal/seeds/fleet_tally_seed.go` | Varsayılan ControlType kayıtları (personel_servis_v1, okul_servis_v1) |
| `internal/seeds/translation_seed.go` | fleet_tally.* namespace altındaki tüm anahtarlar (28 dil) |
| `internal/seeds/menu_seed.go` | fleet_tally menü girişleri |

### 11.2 Çeviri Anahtarları (Tam Liste)

```
fleet_tally.module_name                     → "Filo Doğrulama" / "Fleet Tally"
fleet_tally.control_types.title             → "Kontrol Türleri"
fleet_tally.gps_providers.title             → "GPS Sağlayıcılar"
fleet_tally.connectors.title                → "Sistem Bağlayıcıları"
fleet_tally.verifications.title             → "Doğrulama Sonuçları"
fleet_tally.schedules.title                 → "Çalışma Zamanlamaları"
fleet_tally.runs.title                      → "Çalışma Geçmişi"
fleet_tally.status.ok                       → "Başarılı"
fleet_tally.status.no_device                → "Cihaz Yok"
fleet_tally.status.no_gsm                   → "GSM Hattı Yok"
fleet_tally.status.no_location              → "Konum Bulunamadı"
fleet_tally.status.empty_window             → "Saatler Arası Konum Yok"
fleet_tally.status.not_in_zone              → "Hedef Bölge Dışında"
fleet_tally.status.retry_later              → "Tekrar Kontrol Edilecek"
fleet_tally.status.error                    → "Teknik Hata"
fleet_tally.status.unknown                  → "Bilinmiyor"
fleet_tally.algorithm.point_in_circle       → "Nokta-Çember Kontrolü"
fleet_tally.algorithm.route_intersection    → "Rota-Çember Kesişim Kontrolü"
fleet_tally.window.entry                    → "Giriş Kontrolü"
fleet_tally.window.exit                     → "Çıkış Kontrolü"
fleet_tally.window.any                      → "Herhangi Zaman"
fleet_tally.action.run_now                  → "Şimdi Çalıştır"
fleet_tally.action.test_connection          → "Bağlantıyı Test Et"
fleet_tally.run.total_jobs                  → "Toplam İş"
fleet_tally.run.passed_jobs                 → "Başarılı"
fleet_tally.run.failed_jobs                 → "Başarısız"
fleet_tally.run.error_jobs                  → "Hatalı"
```

---

## 12. Test Stratejisi

### 12.1 Birim Testler

| Bileşen | Test Yaklaşımı |
|---|---|
| `geo/haversine.go` | %100 coverage, bilinen koordinatlarla doğrulama |
| `algorithms/point_in_circle` | Mock GPS noktaları + bilinen beklenen sonuçlar |
| `algorithms/route_intersection` | Segment kesişim köşe case'leri |
| `connectors/service/mapper.go` | Field mapping dönüşüm testleri |
| Service katmanları | Mock repository ile |

### 12.2 Entegrasyon Testleri

- `verifications/service`: Mock GPS adapter + mock connector ile uçtan uca
- `schedules/service/runner`: Mock tüm bağımlılıklar ile tam batch çalışması

### 12.3 Gerçek Sistem Doğrulama (Geçiş Öncesi)

Mevcut PHP+Python sisteminin üretim verileriyle sonuçlarını arşivle. Yeni sistem aynı girdi için aynı sonucu üretiyor mu?

Karşılaştırma testi:
1. N gün için mevcut sistemden kontrol sonuçlarını kaydet
2. Aynı ham veriyi yeni sisteme ver
3. `overall_status` ve `window.status` sonuçlarını karşılaştır
4. Uyumsuzluklar analiz edilir

---

## 13. Gürsel SAP Geçiş Protokolü

### 13.1 Geçiş Öncesi Hazırlık

1. Gürsel için tenant kaydı oluştur
2. SAP SOAP connector kaydını yapılandır:
   - FetchConfig: WSDL URL + method adları
   - FieldMapping: mevcut Python script JSON şeması → normalized
   - Credentials: encrypted kayıt
3. GPS provider config kayıtları: trio, mobiliz, arvento, n2
4. ControlType'lar: personel_servis_v1, okul_servis_v1 (seed ile otomatik gelir)
5. Schedule'lar: her GPS sağlayıcı için ayrı schedule (mevcut `switchedGpsModelId` mantığı)

### 13.2 Paralel Çalışma Dönemi

- Eski Python script'ler çalışmaya devam eder
- Yeni sistem aynı kayıtlar üzerinde **sadece okuma** yapar (PushResults disabled)
- Sonuçlar karşılaştırılır

### 13.3 Tam Geçiş

1. Eski Python script'ler durdurulur
2. Yeni sistem PushResults aktif edilir
3. İlk 7 gün izlenir

---

## 14. Bağımlılık Haritası

```
Wave 1 (domain + CRUD)
    │
    ├── Wave 2 (GPS adapters)    ←─── Wave 4 (verification engine)
    │                                          │
    └── Wave 3 (connectors)     ←─────────────┤
                                               │
                                          Wave 5 (scheduler)
                                               │
                                          Wave 6 (admin UI)
                                               │
                                          Wave 7 (multi-tenant)
```

**Her wave bağımsız deploy edilebilir.** Wave 6 (UI) Wave 5 tamamlanmadan da kısmen geliştirilebilir (mock data ile).

---

## Durum Takibi

| Wave | Durum | Tamamlanma Tarihi |
|---|---|---|
| Wave 1 | ⬜ Planlandı | — |
| Wave 2 | ⬜ Planlandı | — |
| Wave 3 | ⬜ Planlandı | — |
| Wave 4 | ⬜ Planlandı | — |
| Wave 5 | ⬜ Planlandı | — |
| Wave 6 | ⬜ Planlandı | — |
| Wave 7 | ⬜ Süregelen | — |

---

*Analiz için bkz.:* `FLEET_TALLY_SYSTEM_ANALYSIS_TR.md`

*Platform ortak modülleri ve geliştirme kalıpları için bkz.:*  
`FLEET_TALLY_AYSTEK_CORE_IMPL_GUIDE_TR.md`
