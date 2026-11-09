# AYSTEK Core — Modül Geliştirme Uygulama Kılavuzu

> **Bu doküman kimler için?**  
> aystek-core'da yeni bir modül veya alt modül geliştiren herkes — insan ya da AI asistan.  
> "Nasıl yapıyoruz?" sorusunun cevabı burada. Mimari kararlar için `AYSTEK_CORE_CANONICAL_ARCHITECTURE.md`'ye, `fleet_tally` özelindeki plana ise `FLEET_TALLY_SYSTEM_MASTER_PLAN_TR.md`'ye bakın.

**Versiyon:** 1.0.0  
**Tarih:** 2026-05-05

---

## İçindekiler

1. [Proje Yapısı ve Konumlar](#1-proje-yapısı-ve-konumlar)
2. [Temel Katman Mimarisi (DDD)](#2-temel-katman-mimarisi-ddd)
3. [Zorunlu Dosya Sözleşmesi](#3-zorunlu-dosya-sözleşmesi)
4. [Platform Modülleri — Kullanım Kılavuzu](#4-platform-modülleri--kullanım-kılavuzu)
   - [4.1 response — HTTP Yanıt Zarfı](#41-response--http-yanıt-zarfı)
   - [4.2 contextkeys — Tenant/Site Bağlamı](#42-contextkeys--tenantsite-bağlamı)
   - [4.3 httpclient — Güvenli Dış HTTP İstekleri](#43-httpclient--güvenli-dış-http-istekleri)
   - [4.4 jobs/runner — Arka Plan İşleri](#44-jobsrunner--arka-plan-işleri)
   - [4.5 jobs/scheduler — Zamanlanmış Döngüsel Görevler](#45-jobsscheduler--zamanlanmış-döngüsel-görevler)
   - [4.6 circuitbreaker — Devre Kesici](#46-circuitbreaker--devre-kesici)
   - [4.7 idempotency — Tekrar Çalışma Koruması](#47-idempotency--tekrar-çalışma-koruması)
   - [4.8 masking — PII / Hassas Veri Maskeleme](#48-masking--pii--hassas-veri-maskeleme)
   - [4.9 pagination — Sayfalama](#49-pagination--sayfalama)
   - [4.10 logger — Loglama](#410-logger--loglama)
   - [4.11 tenantfilter — Tenant Filtresi](#411-tenantfilter--tenant-filtresi)
5. [Seed Sistemi](#5-seed-sistemi)
   - [5.1 Yeni Seed Dosyası Ekleme](#51-yeni-seed-dosyası-ekleme)
   - [5.2 Menu Seed](#52-menu-seed)
   - [5.3 Translation Seed](#53-translation-seed)
6. [İntegrasyon Modülü Altyapısı (Bağlantılı Kullanım)](#6-i̇ntegrasyon-modülü-altyapısı-bağlantılı-kullanım)
7. [HTTP Handler Kalıpları](#7-http-handler-kalıpları)
8. [MongoDB Repository Kalıpları](#8-mongodb-repository-kalıpları)
9. [Background Worker Kalıpları](#9-background-worker-kalıpları)
10. [Test Kalıpları](#10-test-kalıpları)
11. [Yeni Modül Kontrol Listesi](#11-yeni-modül-kontrol-listesi)

---

## 1. Proje Yapısı ve Konumlar

```
go/aystek-core/
├── cmd/api/main.go                       ← Uygulama giriş noktası
├── internal/
│   ├── modules/                          ← Tüm domain modülleri (60+)
│   │   └── <domain>/                     ← Örn: fleet_tally, integrations
│   │       └── <subdomain>/              ← Örn: control_types, gps_providers
│   │           ├── model/
│   │           ├── dto/
│   │           ├── repository/
│   │           ├── service/
│   │           ├── handler/
│   │           └── routes.go
│   ├── platform/                         ← Ortak altyapı (değiştirme)
│   │   ├── circuitbreaker/               ← Devre kesici (CB store + metrics)
│   │   ├── contextkeys/                  ← Tenant/Site context yönetimi
│   │   ├── httpclient/                   ← SSRF korumalı HTTP istemcisi
│   │   ├── idempotency/                  ← Redis tabanlı tekrar koruma
│   │   ├── jobs/                         ← Arka plan iş altyapısı
│   │   │   └── runner/                   ← Queue + Worker + Job interface
│   │   ├── logger/                       ← Merkezi loglama
│   │   ├── masking/                      ← PII maskeleme
│   │   ├── pagination/                   ← Sayfalama
│   │   ├── response/                     ← HTTP yanıt zarfı
│   │   ├── router/                       ← Gin router gruplaması
│   │   ├── tenantfilter/                 ← Tenant scoped query yardımcıları
│   │   └── ...
│   └── seeds/                            ← Seed dosyaları
│       ├── runner.go                     ← SeedAll() — yeni seed buraya eklenir
│       ├── translation_seed.go           ← 28 dil çevirisi
│       ├── menu_seed.go                  ← Admin menü girişleri
│       └── ...
├── go.mod
└── go.sum
```

**Import path:** `github.com/aystek/aystek-core/internal/...`

---

## 2. Temel Katman Mimarisi (DDD)

Her modül şu katmanlar üzerine kurulur:

```
HTTP Request
    │
    ▼  routes.go        ← Endpoint tanımları, Gin grupları
    │
    ▼  handler/*.go     ← HTTP parse, validate, service çağrısı, JSON yanıt
    │
    ▼  dto/*.go         ← Request/Response yapıları (binding tag'leri)
    │
    ▼  service/*.go     ← İş mantığı, orchestration, hata yönetimi
    │
    ▼  repository/      ← interface.go (sözleşme) + mongo.go (uygulama)
    │
    ▼  model/*.go       ← Domain entity, status sabitleri, domain metodları
    │
    ▼  MongoDB
```

**Kural:** Her katman yalnızca bir alt katmanı çağırır. Handler service'i, service repository'yi çağırır. Handler doğrudan MongoDB'ye bağlanmaz.

---

## 3. Zorunlu Dosya Sözleşmesi

Her alt modül için bu 6 dosya grubu zorunludur:

| # | Dosya | Zorunlu | İçerik |
|---|---|---|---|
| 1 | `model/<domain>.go` | ✅ | Entity struct, status type/const, domain metodları |
| 2 | `dto/<domain>.go` | ✅ | Request/Response DTO'ları, binding tag'leri |
| 3 | `repository/interface.go` | ✅ | Repository interface tanımı |
| 4 | `repository/mongo.go` | ✅ | MongoDB implementasyonu, index'ler |
| 5 | `service/service.go` | ✅ | İş mantığı, dependency injection |
| 6 | `service/service_test.go` | ✅ | Service unit testleri (mock repo) |
| 7 | `handler/*.go` | Koşullu | HTTP endpoint varsa zorunlu |
| 8 | `routes.go` | Koşullu | HTTP endpoint varsa zorunlu |

---

## 4. Platform Modülleri — Kullanım Kılavuzu

### 4.1 `response` — HTTP Yanıt Zarfı

**Import:** `github.com/aystek/aystek-core/internal/platform/response`

Tüm HTTP handler'lar bu paketten yanıt döner. Middleware zaten zarflar ama `response.*` kullanmak explicit ve güvenlidir.

```go
// Başarılı tek kayıt
response.OK(c, gin.H{"control_type": result}, "")

// Oluşturma
response.Created(c, gin.H{"id": newRecord.ID.Hex()}, "created")

// Sayfalı liste
response.WithPagination(c, gin.H{"items": list}, page, "")

// Hata yanıtları
response.BadRequest(c, "geçersiz istek", err.Error())
response.Error(c, http.StatusNotFound, "bulunamadı", nil)
```

**Yanıt zarfı yapısı:**
```json
{
  "success": true,
  "message": "",
  "data": { ... },
  "pagination": { "page": 1, "limit": 25, "total": 100, "total_pages": 4 }
}
```

> **Dikkat:** Middleware `response_envelope.go` otomatik sarar ama `response.OK` kullanırsanız double-wrap olmaz — middleware bunu algılar.

---

### 4.2 `contextkeys` — Tenant/Site Bağlamı

**Import:** `github.com/aystek/aystek-core/internal/platform/contextkeys`

Tenant ve site bilgisi context üzerinden taşınır. Repository sorgularında `tenant_id` filtresi için kullanılır.

```go
// Context'e tenant_id ekle (genellikle middleware yapar)
ctx = contextkeys.WithTenantID(ctx, tenantID)
ctx = contextkeys.WithSiteID(ctx, siteID)

// Context'ten oku
tenantID := contextkeys.GetTenantID(ctx)
siteID   := contextkeys.GetSiteID(ctx)

// Scheduler/worker context'i için (auth middleware olmadan)
ctx = contextkeys.WithIsPublic(ctx, true)
ctx = contextkeys.WithTenantID(ctx, targetTenantID)
```

**Önemli sabit:**
```go
contextkeys.PublicSiteTenantID // = "aystek" — platform'un kendi tenant'ı
```

Worker ve scheduler'larda auth middleware geçmeden çalışıyorsanız mutlaka `WithTenantID` set edin.

---

### 4.3 `httpclient` — Güvenli Dış HTTP İstekleri

**Import:** `github.com/aystek/aystek-core/internal/platform/httpclient`

GPS API'leri, SAP SOAP servisleri, dış ERP sistemleri gibi **dış HTTP istekleri** için her zaman `SecurityClient` kullanın. Plain `net/http.Client` kullanmak SSRF riskine yol açar.

```go
// Oluşturma — allowlist tanımlama zorunludur
client := httpclient.NewSecurityClient(
    []httpclient.AllowlistRule{
        {Domain: "telematik.triomobil.com", Scheme: "https"},
        {Domain: "*.sapservice.gurseltur.com.tr", Scheme: "https"},
    },
    10*1024*1024, // max 10MB yanıt
    false,        // redirect izleme
)

// İstek yapma
headers := map[string]string{
    "Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+pass)),
    "Content-Type":  "application/json",
}
resp, err := client.Do(ctx, "POST", targetURL, body, headers)
if err != nil {
    return nil, fmt.Errorf("dış servis hatası: %w", err)
}
defer resp.Body.Close()
```

**Per-connector SSRF politikası:**
```go
// Her connector kendi allowlist kurallarını getirebilir
resp, err := client.DoWithSSRFPolicy(ctx, "GET", url, nil, headers,
    connectorAllowlistRules,
    false, // private IP izin ver mi
)
```

**Kimlik bilgisi güvenliği:**
- URL'ye kullanıcı adı/şifre **asla** gömme
- `Authorization` header kullan
- Credentials daima `ProviderConfig.Credentials` map'inden al — kaynak kodda sabit olmasın

---

### 4.4 `jobs/runner` — Arka Plan İşleri

**Import:** `github.com/aystek/aystek-core/internal/platform/jobs/runner`

Asenkron, retry destekli arka plan görevleri için kullanılır.

**1. Job tanımla:**
```go
// fleet_tally/verifications/jobs/verify_job.go
package jobs

import (
    "context"
    "github.com/aystek/aystek-core/internal/platform/jobs/runner"
    "github.com/google/uuid"
)

type VerifyJobPayload struct {
    TenantID   string `json:"tenant_id"`
    ScheduleID string `json:"schedule_id"`
}

type VerifyJob struct {
    runner.BaseJob
    Payload VerifyJobPayload
    svc     *verifications.Service
}

func NewVerifyJob(scheduleID, tenantID string, svc *verifications.Service) *VerifyJob {
    return &VerifyJob{
        BaseJob: runner.BaseJob{
            JobID:   uuid.NewString(),
            JobType: "fleet_tally.verify",
        },
        Payload: VerifyJobPayload{TenantID: tenantID, ScheduleID: scheduleID},
        svc:     svc,
    }
}

func (j *VerifyJob) Run(ctx context.Context) error {
    ctx = contextkeys.WithTenantID(ctx, j.Payload.TenantID)
    return j.svc.RunSchedule(ctx, j.Payload.ScheduleID)
}
```

**2. Queue ve runner oluştur (main.go / bootstrap'ta):**
```go
queue  := runner.NewQueue(redisClient, "fleet_tally_jobs", jobFactory)
r      := runner.NewRunner(queue, runner.Config{Workers: 3, MaxRetries: 5})
r.Start(ctx)
defer r.Stop()
```

**3. Job kuyruğa ekle:**
```go
job := jobs.NewVerifyJob(scheduleID, tenantID, verifSvc)
if err := jobRunner.Enqueue(ctx, job); err != nil {
    logger.Printf("job enqueue failed: %v", err)
}
```

**Retry mekanizması:** `runner.RetryDelay(base, attempt)` ile exponential backoff — base'den başlar, 30s'de sabitlenir.

---

### 4.5 `jobs/scheduler` — Zamanlanmış Döngüsel Görevler

**Gerçek kod kalıbı** (billing scheduler'dan):

```go
// fleet_tally/schedules/service/scheduler.go
func StartFleetTallyScheduler(
    ctx context.Context,
    logg *logger.Logger,
    scheduleSvc *schedules.Service,
    jobRunner  *runner.Runner,
) {
    schedulerCtx := contextkeys.WithIsPublic(ctx, true)

    go func() {
        runFleetTallyOnce(schedulerCtx, logg, scheduleSvc, jobRunner, time.Now())
        ticker := time.NewTicker(15 * time.Minute)  // her 15 dakikada kontrol et
        defer ticker.Stop()
        for {
            select {
            case now := <-ticker.C:
                runFleetTallyOnce(schedulerCtx, logg, scheduleSvc, jobRunner, now)
            case <-schedulerCtx.Done():
                return
            }
        }
    }()
}

func runFleetTallyOnce(ctx context.Context, logg *logger.Logger,
    svc *schedules.Service, runner *runner.Runner, now time.Time) {

    activeSchedules, err := svc.ListActive(ctx)
    if err != nil {
        logg.Printf("fleet_tally scheduler: list error: %v", err)
        return
    }
    for _, s := range activeSchedules {
        if !s.IsWithinActiveHours(now) {
            continue  // active_hours dışında → atla
        }
        if !s.IsDue(now) {
            continue  // cron henüz gelmedi → atla
        }
        ctx2 := contextkeys.WithTenantID(ctx, s.TenantID.Hex())
        job  := jobs.NewVerifyJob(s.ID.Hex(), s.TenantID.Hex(), verifSvc)
        if err := runner.Enqueue(ctx2, job); err != nil {
            logg.Printf("fleet_tally scheduler: enqueue err: %v", err)
        }
    }
}
```

**main.go'ya kayıt:**
```go
// cmd/api/main.go — bootstrap bölümüne ekle
jobs.StartFleetTallyScheduler(ctx, logger, scheduleSvc, jobRunner)
```

---

### 4.6 `circuitbreaker` — Devre Kesici

**Import:** `github.com/aystek/aystek-core/internal/platform/circuitbreaker`

GPS API'leri veya SAP SOAP gibi **dış servis çağrılarında** circuit breaker kullanın. Bağlantı kopuksa tekrar tekrar denemeyi durdurur.

```go
// Servis init'inde oluştur
cbStore := circuitbreaker.NewInMemoryCircuitBreakerStore(circuitbreaker.DefaultConfig())

// Her dış API çağrısı öncesi
allowed, state, err := cbStore.IsAllowed(ctx, connectorID)
if err != nil || !allowed {
    return nil, fmt.Errorf("circuit breaker %s durumunda: %s", connectorID, state)
}

// Çağrı yap
result, callErr := gpsAdapter.FetchLocations(ctx, cfg, req)

// Sonuca göre kaydet
if callErr != nil {
    cbStore.RecordFailure(ctx, connectorID)
    return nil, callErr
}
cbStore.RecordSuccess(ctx, connectorID)
return result, nil
```

**Varsayılan konfigürasyon:**
```go
circuitbreaker.DefaultConfig()
// FailureThreshold:    5  (5 ardışık hata → OPEN)
// SuccessThreshold:    2  (HALF_OPEN'da 2 başarı → CLOSED)
// OpenTimeout:         30s
// FailureWindow:       60s
```

---

### 4.7 `idempotency` — Tekrar Çalışma Koruması

**Import:** `github.com/aystek/aystek-core/internal/platform/idempotency`

Aynı job'un birden fazla çalışmasını engellemek için (RabbitMQ redelivery, cron çakışması vs.).

```go
idemStore := idempotency.NewRedisIdempotencyStore(redisClient)

// İş başında kontrol
key := fmt.Sprintf("fleet_tally:verify:%s:%s", tenantID, scheduleRunID)
exists, cachedResult, err := idemStore.CheckAndStore(ctx, key, 24*time.Hour)
if err != nil {
    return err
}
if exists {
    // Zaten çalışmış, önceki sonuç döner
    return nil
}

// İşi yap
runResult, runErr := doWork(ctx)

// Sonucu kaydet
_ = idemStore.Store(ctx, key, &idempotency.ExecutionResult{
    Success:    runErr == nil,
    ExecutedAt: time.Now(),
}, 24*time.Hour)

return runErr
```

**Hata hash anahtarı (çakışma önleme):**
```go
key := idempotency.IdempotencyKeyWithHash(tenantID, scheduleID, runDate, 1, payloadBytes)
```

---

### 4.8 `masking` — PII / Hassas Veri Maskeleme

**Import:** `github.com/aystek/aystek-core/internal/platform/masking`

Log'lara düşecek veya API yanıtlarında gösterilecek hassas alanları maskele.

```go
policy := masking.DefaultPIIMaskingPolicy()

// Kimlik bilgisi içeren struct'ı log öncesi maskele
masked, err := masking.MaskStruct(connectorCredentials, policy)
logger.Printf("connector test: %+v", masked)
```

**Varsayılan maskeli alanlar:** `password`, `secret`, `token`, `api_key`, `authorization`, `email`, `phone`, kredi kartı numaraları, IBAN, adres alanları.

**Kural:** GPS sağlayıcı ve SAP connector kimlik bilgilerini log'a yazarken her zaman maskeleme uygula.

---

### 4.9 `pagination` — Sayfalama

**Import:** `github.com/aystek/aystek-core/internal/platform/pagination`

```go
// Handler'da parametre çekme
page, limit := pagination.FromQuery(c)  // ?page=1&limit=25

// Repository sorgusu
opts := options.Find().
    SetSkip(int64((page - 1) * limit)).
    SetLimit(int64(limit))

total, err := repo.Count(ctx, filter)
items, err := repo.List(ctx, filter, opts)

// Yanıt
p := pagination.New(page, limit, total)
response.WithPagination(c, gin.H{"items": items}, p, "")
```

---

### 4.10 `logger` — Loglama

**Import:** `github.com/aystek/aystek-core/internal/platform/logger`

```go
// Oluşturma (genellikle main.go'da oluşturulur, inject edilir)
logg := logger.New()

// Kullanım
logg.Printf("fleet_tally: schedule %s çalışmaya başladı", scheduleID)
logg.Printf("fleet_tally: job hata: %v", err)
```

Prefix `[aystek-core]` otomatik eklenir.

---

### 4.11 `tenantfilter` — Tenant Filtresi

MongoDB sorgularına tenant filtresi eklemek için:

```go
import "github.com/aystek/aystek-core/internal/platform/tenantfilter"

// Repository'de kullanım
func (r *mongoRepository) List(ctx context.Context, opts ListOptions) ([]*model.VerificationJob, error) {
    filter := tenantfilter.Build(ctx)   // { "tenant_id": <ctx tenant_id> }
    if opts.Status != "" {
        filter["overall_status"] = opts.Status
    }
    // ...
}
```

---

## 5. Seed Sistemi

### 5.1 Yeni Seed Dosyası Ekleme

**Adım 1:** `internal/seeds/fleet_tally_seed.go` dosyasını oluştur.

```go
package seeds

import (
    "context"
    "log"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func SeedFleetTally(ctx context.Context, db *mongo.Database) error {
    log.Println("🚀 Seeding fleet_tally...")
    return seedFleetTallyControlTypes(ctx, db)
}

func seedFleetTallyControlTypes(ctx context.Context, db *mongo.Database) error {
    coll := db.Collection("fleet_tally_control_types")
    for _, seed := range defaultControlTypes {
        filter := bson.M{"name": seed.Name}
        update := bson.M{"$setOnInsert": seed}
        _, err := coll.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
        if err != nil {
            return err
        }
    }
    return nil
}
```

**Adım 2:** `internal/seeds/runner.go`'ya ekle:

```go
// runner.go içinde, mevcut seed çağrıları sırasına ekle:
if err := SeedFleetTally(ctx, db); err != nil {
    return err
}
```

> **Sıra önemli:** Translation ve Menu seed'lerinden ÖNCE olması gerekmez ama `SeedLocalization`'dan SONRA olmalı.

### 5.2 Menu Seed

`internal/seeds/menu_seed.go` içinde `seeds` slice'ına ekle:

```go
// Fleet Tally — Filo Doğrulama
{Key: "fleet-tally", TitleKey: "menu.fleet_tally", TitleTR: "Filo Doğrulama",
    TitleEN: "Fleet Tally", Path: "/admin/fleet-tally", Icon: "check-square",
    Order: 15, Roles: []string{"admin"}, Visible: true},

{Key: "fleet-tally-control-types", Parent: "fleet-tally",
    TitleKey: "menu.fleet_tally.control_types", TitleTR: "Kontrol Türleri",
    TitleEN: "Control Types", Path: "/admin/fleet-tally/control-types",
    Icon: "list-checks", Order: 1, Roles: []string{"admin"}, Visible: true},

{Key: "fleet-tally-gps-providers", Parent: "fleet-tally",
    TitleKey: "menu.fleet_tally.gps_providers", TitleTR: "GPS Sağlayıcılar",
    TitleEN: "GPS Providers", Path: "/admin/fleet-tally/gps-providers",
    Icon: "radio-tower", Order: 2, Roles: []string{"admin"}, Visible: true},

{Key: "fleet-tally-connectors", Parent: "fleet-tally",
    TitleKey: "menu.fleet_tally.connectors", TitleTR: "Sistem Bağlayıcıları",
    TitleEN: "Connectors", Path: "/admin/fleet-tally/connectors",
    Icon: "plug", Order: 3, Roles: []string{"admin"}, Visible: true},

{Key: "fleet-tally-schedules", Parent: "fleet-tally",
    TitleKey: "menu.fleet_tally.schedules", TitleTR: "Zamanlamalar",
    TitleEN: "Schedules", Path: "/admin/fleet-tally/schedules",
    Icon: "clock", Order: 4, Roles: []string{"admin"}, Visible: true},

{Key: "fleet-tally-verifications", Parent: "fleet-tally",
    TitleKey: "menu.fleet_tally.verifications", TitleTR: "Doğrulama Sonuçları",
    TitleEN: "Verifications", Path: "/admin/fleet-tally/verifications",
    Icon: "clipboard-check", Order: 5, Roles: []string{"admin"}, Visible: true},

{Key: "fleet-tally-runs", Parent: "fleet-tally",
    TitleKey: "menu.fleet_tally.runs", TitleTR: "Çalışma Geçmişi",
    TitleEN: "Runs", Path: "/admin/fleet-tally/runs",
    Icon: "history", Order: 6, Roles: []string{"admin"}, Visible: true},
```

### 5.3 Translation Seed

`internal/seeds/translation_seed.go` içinde ilgili 28 dil bloğuna ekle:

```go
// fleet_tally namespace — her dil bloğuna ekle
{"key": "menu.fleet_tally",                "tr": "Filo Doğrulama",       "en": "Fleet Tally"},
{"key": "menu.fleet_tally.control_types",  "tr": "Kontrol Türleri",      "en": "Control Types"},
{"key": "menu.fleet_tally.gps_providers",  "tr": "GPS Sağlayıcılar",     "en": "GPS Providers"},
{"key": "menu.fleet_tally.connectors",     "tr": "Sistem Bağlayıcıları", "en": "Connectors"},
{"key": "menu.fleet_tally.schedules",      "tr": "Zamanlamalar",         "en": "Schedules"},
{"key": "menu.fleet_tally.verifications",  "tr": "Doğrulama Sonuçları",  "en": "Verifications"},
{"key": "menu.fleet_tally.runs",           "tr": "Çalışma Geçmişi",      "en": "Runs"},
{"key": "fleet_tally.status.ok",           "tr": "Başarılı",             "en": "OK"},
{"key": "fleet_tally.status.no_device",    "tr": "Cihaz Yok",            "en": "No Device"},
{"key": "fleet_tally.status.no_gsm",       "tr": "GSM Hattı Yok",        "en": "No GSM"},
{"key": "fleet_tally.status.no_location",  "tr": "Konum Bulunamadı",     "en": "No Location"},
{"key": "fleet_tally.status.empty_window", "tr": "Saatler Arası Konum Yok","en": "Empty Time Window"},
{"key": "fleet_tally.status.not_in_zone",  "tr": "Hedef Bölge Dışında",  "en": "Not In Zone"},
{"key": "fleet_tally.status.retry_later",  "tr": "Tekrar Kontrol Edilecek","en": "Retry Later"},
{"key": "fleet_tally.status.error",        "tr": "Teknik Hata",          "en": "Error"},
{"key": "fleet_tally.status.unknown",      "tr": "Bilinmiyor",           "en": "Unknown"},
```

> **28 dil:** Mevcut translation_seed.go'daki dil listesine bakarak eksik dilleri tamamla. TR ve EN zorunludur; diğerleri TR değeriyle kopyalanabilir.

---

## 6. İntegrasyon Modülü Altyapısı (Bağlantılı Kullanım)

`fleet_tally/connectors` modülü `integrations` modülüyle entegrasyon noktaları:

**CircuitBreakerStore:** `integrations/governance` ve `circuitbreaker` paketi. Her connector ID için ayrı CB state.

**Worker pattern:** `integrations/worker/integration_worker.go` — RabbitMQ tabanlı async worker referans implementasyonu. Fleet tally için Redis tabanlı job queue tercih edilir (daha basit).

**Tracing:** `integrations/tracing` paketi — dış servis çağrılarını OpenTelemetry span'larıyla izle.

**Metrics:** `integrations/metrics` paketi — çağrı sayacı, latency histogram referans için bakılabilir.

---

## 7. HTTP Handler Kalıpları

### 7.1 Standart Admin Handler

```go
package handler

import (
    "net/http"

    "github.com/aystek/aystek-core/internal/modules/fleet_tally/control_types/dto"
    "github.com/aystek/aystek-core/internal/modules/fleet_tally/control_types/service"
    "github.com/aystek/aystek-core/internal/platform/pagination"
    "github.com/aystek/aystek-core/internal/platform/response"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminHandler struct {
    svc *service.Service
}

func NewAdminHandler(svc *service.Service) *AdminHandler {
    return &AdminHandler{svc: svc}
}

// POST /admin/fleet-tally/control-types
func (h *AdminHandler) Create(c *gin.Context) {
    var req dto.CreateControlTypeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "geçersiz istek", err.Error())
        return
    }
    result, err := h.svc.Create(c.Request.Context(), req)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "oluşturulamadı", err.Error())
        return
    }
    response.Created(c, gin.H{"control_type": result}, "oluşturuldu")
}

// GET /admin/fleet-tally/control-types
func (h *AdminHandler) List(c *gin.Context) {
    page, limit := pagination.FromQuery(c)
    items, total, err := h.svc.List(c.Request.Context(), page, limit)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "listelenemedi", err.Error())
        return
    }
    p := pagination.New(page, limit, total)
    response.WithPagination(c, gin.H{"items": items}, p, "")
}

// GET /admin/fleet-tally/control-types/:id
func (h *AdminHandler) GetByID(c *gin.Context) {
    id, err := primitive.ObjectIDFromHex(c.Param("id"))
    if err != nil {
        response.BadRequest(c, "geçersiz id", err.Error())
        return
    }
    result, err := h.svc.GetByID(c.Request.Context(), id)
    if err != nil {
        response.Error(c, http.StatusNotFound, "bulunamadı", err.Error())
        return
    }
    response.OK(c, gin.H{"control_type": result}, "")
}
```

### 7.2 Routes Dosyası

```go
package control_types

import (
    "github.com/aystek/aystek-core/internal/modules/fleet_tally/control_types/handler"
    "github.com/aystek/aystek-core/internal/modules/fleet_tally/control_types/service"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(admin *gin.RouterGroup, db *mongo.Database) {
    svc := service.New(/* deps */)
    h   := handler.NewAdminHandler(svc)

    g := admin.Group("/fleet-tally/control-types")
    {
        g.POST("",     h.Create)
        g.GET("",      h.List)
        g.GET("/:id",  h.GetByID)
        g.PUT("/:id",  h.Update)
        g.DELETE("/:id", h.Delete)
    }
}
```

---

## 8. MongoDB Repository Kalıpları

```go
package repository

import (
    "context"
    "time"

    "github.com/aystek/aystek-core/internal/modules/fleet_tally/control_types/model"
    "github.com/aystek/aystek-core/internal/platform/tenantfilter"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
    coll *mongo.Collection
}

// Constructor: *mongo.Database alır, collection adı sabit
func NewMongoRepository(db *mongo.Database) Repository {
    return &mongoRepository{coll: db.Collection("fleet_tally_control_types")}
}

func (r *mongoRepository) Create(ctx context.Context, ct *model.ControlType) (*model.ControlType, error) {
    if ct.ID.IsZero() {
        ct.ID = primitive.NewObjectID()
    }
    ct.CreatedAt = time.Now()
    ct.UpdatedAt = time.Now()
    _, err := r.coll.InsertOne(ctx, ct)
    return ct, err
}

func (r *mongoRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.ControlType, error) {
    var ct model.ControlType
    filter := bson.M{"_id": id}
    err := r.coll.FindOne(ctx, filter).Decode(&ct)
    if err == mongo.ErrNoDocuments {
        return nil, nil  // Bulunamazsa nil, nil döner — hata değil
    }
    return &ct, err
}

func (r *mongoRepository) List(ctx context.Context, page, limit int) ([]*model.ControlType, int64, error) {
    filter := tenantfilter.Build(ctx)
    total, err := r.coll.CountDocuments(ctx, filter)
    if err != nil {
        return nil, 0, err
    }
    opts := options.Find().
        SetSkip(int64((page-1)*limit)).
        SetLimit(int64(limit)).
        SetSort(bson.D{{Key: "created_at", Value: -1}})
    cursor, err := r.coll.Find(ctx, filter, opts)
    if err != nil {
        return nil, 0, err
    }
    defer cursor.Close(ctx)
    var items []*model.ControlType
    return items, total, cursor.All(ctx, &items)
}

func (r *mongoRepository) CreateIndexes(ctx context.Context) error {
    _, err := r.coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
        {
            Keys:    bson.D{{Key: "tenant_id", Value: 1}, {Key: "name", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
        {
            Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "status", Value: 1}},
        },
    })
    return err
}
```

**Kural özetleri:**
- Constructor her zaman `*mongo.Database` alır (collection adı içeride)
- Bulunamazsa `nil, nil` döner (hata değil)
- Timestamps (`CreatedAt`, `UpdatedAt`) repo'da set edilir
- `tenantfilter.Build(ctx)` her List/Count sorgusunda zorunlu
- `CreateIndexes` her modülde uygulanmalı

---

## 9. Background Worker Kalıpları

### RabbitMQ Worker (Ağır iş, dış event)

`integrations/worker/integration_worker.go` referans alınabilir. `amqp091-go` kullanır.

### Redis Queue Worker (Hafif iş, iç job)

`platform/jobs/runner` + `platform/jobs/runner/queue.go` kullanılır. `runner.Job` interface'ini implement et, `runner.Runner.Enqueue()` ile kuyruğa ekle.

### Scheduler (Tekrarlayan)

`platform/jobs/scheduler.go` (billing scheduler) kalıbını takip et:
- `go func() { runOnce(); ticker := time.NewTicker(interval); for { select ... } }()`
- `contextkeys.WithIsPublic` + `contextkeys.WithTenantID` ile context hazırla
- `cmd/api/main.go` bootstrap'ta çağır

---

## 10. Test Kalıpları

### Service Unit Testi (Mock Repository)

```go
package service_test

import (
    "context"
    "testing"
    "time"

    "github.com/aystek/aystek-core/internal/modules/fleet_tally/control_types/model"
    "github.com/aystek/aystek-core/internal/modules/fleet_tally/control_types/service"
)

// Mock Repository — tüm interface metodlarını implement et
type mockRepo struct {
    items map[string]*model.ControlType
}

func (m *mockRepo) Create(ctx context.Context, ct *model.ControlType) (*model.ControlType, error) {
    m.items[ct.ID.Hex()] = ct
    return ct, nil
}
func (m *mockRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*model.ControlType, error) {
    ct, ok := m.items[id.Hex()]
    if !ok { return nil, nil }
    return ct, nil
}
// ... diğer metodlar

func TestCreate_Success(t *testing.T) {
    repo := &mockRepo{items: make(map[string]*model.ControlType)}
    svc  := service.New(repo)

    result, err := svc.Create(context.Background(), dto.CreateControlTypeRequest{
        Name:        "test_control",
        AlgorithmID: "point_in_circle_multi_window",
    })

    if err != nil { t.Fatalf("beklenmedi: %v", err) }
    if result.Name != "test_control" { t.Errorf("isim uyuşmuyor") }
}
```

**Kural:** Repository interface'e yeni metod eklendiğinde tüm test mock'larına aynı metod (no-op stub) eklenmeli.

### Geometri Testi (Haversine)

```go
func TestIsPointInCircle(t *testing.T) {
    tests := []struct {
        name     string
        pointLat, pointLon float64
        centerLat, centerLon float64
        radius   float64
        want     bool
    }{
        {"içinde", 41.0350, 29.2570, 41.0348, 29.2567, 700, true},
        {"dışında", 41.0400, 29.2700, 41.0348, 29.2567, 100, false},
        {"sınırda", 41.0348, 29.2567, 41.0348, 29.2567, 0, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := geo.IsPointInCircle(tt.pointLat, tt.pointLon,
                tt.centerLat, tt.centerLon, tt.radius)
            if got != tt.want {
                t.Errorf("IsPointInCircle() = %v, istenen %v", got, tt.want)
            }
        })
    }
}
```

---

## 11. Yeni Modül Kontrol Listesi

Her wave tamamlanmadan önce bu listeyi kontrol et:

**Backend (Go):**
- [ ] `model/` — Entity struct, status type/const, domain metodları
- [ ] `dto/` — Request/Response DTO'ları, `binding` tag'leri
- [ ] `repository/interface.go` — Repository interface
- [ ] `repository/mongo.go` — Mongo implementasyonu + `CreateIndexes()`
- [ ] `service/service.go` — İş mantığı, dependency injection
- [ ] `service/service_test.go` — Mock repo ile unit testler
- [ ] `handler/admin.go` — Admin HTTP handler'ları
- [ ] `routes.go` — Route kaydı
- [ ] `internal/seeds/<module>_seed.go` — Seed dosyası oluşturuldu
- [ ] `internal/seeds/runner.go` — `Seed<Module>()` çağrısı eklendi
- [ ] `internal/seeds/menu_seed.go` — Menü girişleri eklendi
- [ ] `internal/seeds/translation_seed.go` — 28 dil anahtarları eklendi
- [ ] `cmd/api/main.go` — Scheduler ve worker kaydı yapıldı (varsa)

**Frontend (Next.js):**
- [ ] `app/(dashboard)/admin/<path>/page.tsx` sayfaları oluşturuldu
- [ ] Tüm kullanıcıya yönelik metinler `t(...)` / `useI18n` ile kullanılıyor
- [ ] API çağrıları `lib/api/<module>.ts` üzerinden yapılıyor

**Güvenlik:**
- [ ] Dış HTTP istekleri `httpclient.SecurityClient` ile yapılıyor
- [ ] Kimlik bilgileri kaynak kodda sabit değil — `Credentials` map'inden okunuyor
- [ ] Log'a düşen hassas veriler `masking` uygulanıyor
- [ ] Dış API çağrıları circuit breaker ile sarılıyor
- [ ] Tekrarlayan işler idempotency key ile korunuyor

---

*Bu kılavuz referans dokümanlarla birlikte okunmalıdır:*  
- `AYSTEK_CORE_CANONICAL_ARCHITECTURE.md` — Mimari kararlar  
- `README_FOR_AI_ASSISTANTS.md` — Tüm standartlar  
- `FLEET_TALLY_SYSTEM_ANALYSIS_TR.md` — Domain analizi  
- `FLEET_TALLY_SYSTEM_MASTER_PLAN_TR.md` — İş planı
