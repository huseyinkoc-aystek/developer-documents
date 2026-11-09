# Kişisel Asistan — Enterprise Düzey Ajans, Entegrasyon ve AI Pipeline Planı

**Tarih:** 2026-06-24  
**Kapsam:** aystek-core backend + kisisel_asistan Flutter uygulaması  
**Durum:** Aktif — Kısmi Uygulama Tamamlandı  
**Standart:** AYSTEK Core DDD + OOP, AI Assistant README v1.1.0

---

## 1. Mevcut Durum Özeti

### 1.1 Tamamlanan Düzeltmeler (Bu Oturumda)

| # | Değişiklik | Dosya | Durum |
|---|-----------|-------|-------|
| 1 | `_AgencyLoginSheet` güvenlik notu metin kırpılması düzeltildi | `news_prefs_screen.dart` | ✅ |
| 2 | `SeedNewsAgencies` pers-assit sitesini kapsayacak şekilde genişletildi | `news_agency_seed.go` | ✅ |
| 3 | AA ve ANKA için RSS `ProviderConnector` seed'i eklendi | `news_agency_seed.go` | ✅ |
| 4 | Flutter APK build (güvenlik notu fix dahil) | `app-debug.apk` | ✅ |

### 1.2 Mimari Durum Tespiti

```
Flutter App (kisisel_asistan)
│
├── Ajanslar Sekmesi → FlutterSecureStorage'a kimlik bilgisi kaydı
│   ⚠️ Kaydedilen bilgiler backend'e gönderilmiyor — bağlantı doğrulaması YOK
│
├── PA News Feed → GET /pa/news/feed
│   ⚠️ AA news_source kaydında base_url BOŞ → AA haberleri gelmiyor
│   ✅ PA news_sources seed'i doğru URL'ler içeriyor (seed yeniden çalıştırılmalı)
│
└── AI Sohbet → SQLite'da kalıcı geçmiş (✅ Tamamlandı önceki oturumda)

AYSTEK Core (aystek-core)
│
├── news_agencies modülü
│   ✅ AgencyProvider modeli tam
│   ✅ ProviderConnector modeli tam
│   ⚠️ agency_providers koleksiyonunda pers-assit sitesi için VERİ YOK
│   ⚠️ provider_connectors koleksiyonunda hiç kayıt YOK
│
├── news_sources koleksiyonu
│   ⚠️ 15 kayıt var ama hepsi base_url = "" (boş)
│   ✅ pa_news_sources_seed.go doğru URL'ler içeriyor
│
└── integrations modülü
    ✅ Connector/Action/Hook/Execution mimarisi tam
    ⚠️ news_agencies → integrations bağlantısı HENÜZ YAPILMADI
```

---

## 2. Hedef Mimari (Enterprise Pipeline)

```
┌─────────────────────────────────────────────────────────────────────┐
│                     AYSTEK CORE BACKEND                              │
│                                                                       │
│  ┌─────────────────┐    ┌──────────────────┐    ┌────────────────┐  │
│  │  news_agencies  │    │   integrations   │    │  personal_     │  │
│  │  modülü         │───►│   modülü         │───►│  assistant     │  │
│  │                 │    │                  │    │  modülü        │  │
│  │ AgencyProvider  │    │ Connector        │    │                │  │
│  │ ProviderConn.   │    │ Hook             │    │ /pa/news/feed  │  │
│  │ IngestionRun    │    │ Execution        │    │ AI filtreleme  │  │
│  └────────┬────────┘    └────────┬─────────┘    └───────┬────────┘  │
│           │                     │                       │           │
│           │ RSS/API çekme       │ Event publishing      │ Kişisel.  │
│           ▼                     ▼                       ▼ haber     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                         MongoDB                              │   │
│  │  agency_providers  |  provider_connectors  |  news_sources  │   │
│  │  ingestion_runs    |  pa_news_preferences  |  pa_user_prof. │   │
│  └─────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                                │
              ┌─────────────────┴─────────────────┐
              │                                   │
┌─────────────┴──────────────┐    ┌───────────────┴───────────────┐
│  kisisel_asistan (Flutter) │    │  aystek-core-next (Admin)     │
│                            │    │                               │
│  Haber Akışı (PA feed)     │    │  news_agencies modülü UI      │
│  AI Sohbet (SQLite)        │    │  Ajans yönetimi / connector   │
│  Ajanslar tab (credential) │    │  Ingestion run takibi         │
└────────────────────────────┘    └───────────────────────────────┘
```

---

## 3. Uygulama Planı (Aşamalı)

### Faz 1 — Temel Altyapı (TAMAMLANDI / DEVAM)

#### 1.1 Seed Düzeltmesi (✅ Tamamlandı)
- `SeedNewsAgencies` pers-assit sitesini artık kapsıyor
- AA (güncel + ekonomi) ve ANKA için RSS connector seed'leri eklendi
- **Aksiyon:** Seed'i tekrar çalıştırmak gerekiyor (bkz. Bölüm 5.1)

#### 1.2 PA news_sources `base_url` Düzeltmesi (⏳ Bekliyor)
- `RunPANewsSourcesSeed` doğru URL'leri içeriyor
- MongoDB'de mevcut kayıtlar `$set` ile güncellenecek
- **Aksiyon:** Seed yeniden çalıştırılmalı VEYA admin PUT endpoint kullanılmalı

#### 1.3 Flutter Güvenlik Notu Düzeltmesi (✅ Tamamlandı)
- `_AgencyLoginSheet` artık nav bar yüksekliğini de dikkate alıyor

---

### Faz 2 — news_agencies → personal_assistant Köprüsü

#### 2.1 PA News Feed Servisi Güncellenmesi

**Hedef:** `GetNewsFeed()` fonksiyonu `news_sources` koleksiyonundan değil,  
`provider_connectors` koleksiyonundan (RSS type) URL'leri çeksin.

**Mevcut akış:**
```go
// service.go - GetNewsFeed
sources = userPrefs.SourceIDs  // kullanıcı tercihleri
for each sourceID:
    src = newsSourcesRepo.FindByID(sourceID)
    fetchRSS(src.BaseURL)      // boş URL sorunlu!
```

**Hedef akış (OOP refactor):**
```go
// Yeni yaklaşım — ProviderConnector'dan RSS URL'lerini al
connectors = providerConnectorsRepo.ListByType(ctx, "rss", "active")
for each connector:
    if userHasSubscription(connector.ProviderID):
        fetchRSS(connector.Config["feed_url"])
```

**Gerekli değişiklikler:**
- `personal_assistant/service/service.go`: `GetNewsFeed()` refactor
- `personal_assistant/service/service.go`: `newsAgenciesRepo` bağımlılığı ekleme
- `cmd/api/main.go`: DI wiring güncelleme

#### 2.2 Kullanıcı Aboneliği ve Connector Bağlantısı

**Yeni model:**
```go
// PAUserAgencySubscription - kullanıcı bir ajansa "abone" olur
// (şu an yok — eklenmeli)
type PAUserAgencySubscription struct {
    UserID      primitive.ObjectID
    ProviderID  primitive.ObjectID  // AgencyProvider'a referans
    Plan        string              // "trial" | "standard" | "premium"
    CreatedAt   time.Time
}
```

**Endpoint (yeni):**
```
POST /pa/agencies/subscribe      — kullanıcı ajans aboneliği başlatır
GET  /pa/agencies/subscriptions  — kullanıcının abone olduğu ajanslar
```

---

### Faz 3 — Integrations Modülü Bağlantısı (Enterprise)

#### 3.1 Event Publishing: news_agencies → integrations

`news_agencies` modülü zaten `provider.created`, `ingestion.run.started` gibi olayları  
`provider_outbox` koleksiyonuna yazar. Bu olaylar `integrations` eventbridge'ine bağlanmalı.

**Akış:**
```
ingestion.run.completed (event)
    └── integrations eventbridge listener
            └── Hook: "news_agencies.ingestion.completed"
                    └── Action: POST /pa/news/ingest (veya internal queue)
                            └── PA modülü yeni haberleri kullanıcı feed'lerine ekler
```

**Yeni integrations Hook seed'i:**
```go
// integrations_seed.go'ya eklenecek
Hook{
    EventKey:    "ingestion.run.completed",
    Name:        "PA News Feed Refresh",
    TenantID:    "aystek",
    SiteID:      persAssistSiteID,
    ConnectorID: paInternalConnectorID,
    ActionID:    paNewsIngestActionID,
    Mode:        "async",
    Status:      "active",
}
```

#### 3.2 Ingestion Worker: connector URL → news_sources

**Mevcut boşluk:** `ProviderConnector`'da RSS URL var ama bunu otomatik çeken  
bir worker YOK. `TriggerRun()` çağrısı sadece `IngestionRun` kaydı oluşturuyor,  
gerçek fetch yapmıyor.

**Çözüm — Yeni servis: `AgencyRSSWorker`:**
```go
// internal/modules/news_agencies/worker/rss_worker.go (YENİ)
type AgencyRSSWorker struct {
    connectorRepo  repository.Repository
    newsSourceRepo newssourcesrepo.Repository
    httpClient     *http.Client
}

func (w *AgencyRSSWorker) ProcessConnector(ctx context.Context, connectorID primitive.ObjectID) error {
    connector, _ := w.connectorRepo.FindConnectorByID(ctx, connectorID)
    feedURL := connector.Config["feed_url"].(string)
    items, _ := w.fetchRSS(ctx, feedURL)
    for _, item := range items {
        w.newsSourceRepo.Upsert(ctx, newsItem)
    }
    return nil
}
```

---

### Faz 4 — AI Entegrasyonu

#### 4.1 Mevcut Durum
- AI sohbet: Claude/GPT/Gemini API doğrudan Flutter'dan
- `ai_service.dart` → provider bazlı
- Sohbet geçmişi: SQLite (önceki oturumda tamamlandı)

#### 4.2 AI News Filtreleme Pipeline (Backend)

**Hedef:** Kullanıcının haber tercihlerine göre AI skorlaması yapılsın.

```go
// personal_assistant/service/service.go - yeni metot
func (s *Service) FilterNewsWithAI(ctx context.Context, 
    articles []NewsArticle, prefs PANewsPreferences) ([]NewsArticle, error) {
    
    prompt := buildFilterPrompt(articles, prefs)
    
    // Claude Haiku 4.5 ile hızlı skorlama
    response, _ := s.aiClient.Complete(ctx, AIRequest{
        Model:  "claude-haiku-4-5-20251001",
        Prompt: prompt,
    })
    
    return parseFilteredArticles(response, articles), nil
}
```

**Prompt şablonu:**
```
Sen bir haber editörüsün. Aşağıdaki haberleri kullanıcı tercihlerine göre 0-100 arası puanla.

Kullanıcı tercihleri:
- Şehir/Bölge: {{user.region}}
- Kategoriler: {{user.categories}}
- Anahtar kelimeler: {{user.keywords}}

Haberler:
{{articles_json}}

Yanıt formatı: JSON array [{"id": "...", "score": 85, "reason": "..."}]
```

#### 4.3 Flutter Tarafı Değişiklikler
- `AIChatNotifier.sendMessage()`: haber bağlamı enjeksiyonu (seçili haber → sohbete ekle)
- `news_tab_screen.dart`: "AI ile Analiz Et" butonu
- Seçili haber başlık + özet → AI sohbet mesajına dönüştür

---

## 4. OOP Tasarım Prensipleri (AYSTEK Standartları)

### 4.1 Bağımlılık Enjeksiyonu Hiyerarşisi

```
main.go
│
├── newsAgenciesRepo    = news_agencies.NewMongoRepository(db)
├── newsAgenciesSvc     = news_agencies.New(newsAgenciesRepo, logSvc)
├── paRepo              = personal_assistant.NewMongoRepository(db)
├── paSvc               = personal_assistant.New(
│       paRepo,
│       newsSourceRepo,
│       newsAgenciesRepo,  ← YENİ bağımlılık
│       aiClient           ← YENİ bağımlılık
│   )
└── paHandler           = personal_assistant.NewHandler(paSvc)
```

### 4.2 Domain Events (Bildirim Zinciri)

```
news_agencies.IngestionRun.completed
    → integrations.eventbridge.Listen("ingestion.run.completed")
        → integrations.Hook.match(tenantID, siteID, eventKey)
            → integrations.executor.POST /pa/news/ingest
                → personal_assistant.IngestNews(items)
                    → AI filtreleme
                        → pa_news_feed güncelleme
                            → SSE push (kullanıcıya anlık bildirim)
```

### 4.3 Interface Tanımları (Yeni)

```go
// news_agencies/service/interface.go
type AgencyIngestionService interface {
    TriggerRun(ctx context.Context, connectorID primitive.ObjectID, dryRun bool, actorID *primitive.ObjectID) (*model.IngestionRun, error)
    ReportRun(ctx context.Context, connectorID primitive.ObjectID, in RunResultInput) error
    ListRuns(ctx context.Context, connectorID primitive.ObjectID, limit int64) ([]model.IngestionRun, error)
}

// personal_assistant/service/interface.go (güncelleme)
type NewsIngestionPort interface {
    IngestArticles(ctx context.Context, connectorID primitive.ObjectID, articles []model.NewsArticle) error
}
```

---

## 5. Acil Aksiyonlar

### 5.1 Seed Yeniden Çalıştırma
```bash
# Aystek Core sunucusunda:
cd /path/to/aystek-core
./seed --mongo mongodb://localhost:27017 --db aystek_core
```

**Ne değişecek:**
- `agency_providers`: pers-assit sitesi için AA, ANKA, DHA, İHA, Reuters, AP kayıtları
- `provider_connectors`: AA (güncel RSS + ekonomi RSS) ve ANKA RSS connector'ları
- `news_sources`: AA `base_url` = `https://www.aa.com.tr/tr/rss/default?cat=guncel`

### 5.2 APK Kurulumu
```bash
# Cihazı bağla (kullanıcıdan port alınacak):
adb connect 192.168.1.2:<PORT>
adb -s 192.168.1.2:<PORT> install -r /workspace/flutter/kisisel_asistan/build/app/outputs/flutter-apk/app-debug.apk
```

**APK'da ne var:**
- `_AgencyLoginSheet` güvenlik notu artık tam görünüyor
- Nav bar yüksekliği doğru hesaplanıyor

### 5.3 AA Admin Panel Test
1. Admin panel → pers-assit sitesi → news_agencies
2. Seed sonrası Anadolu Ajansı ve ANKA görünmeli
3. Her ajans için connector'lar listelenmeli
4. "Test Connection" ile RSS bağlantısı doğrulanmalı

---

## 6. Frontend (aystek-core-next) Gereksinimleri

`news_agencies` modülünün pers-assit sitesinde görünmesi için:

### 6.1 Mevcut Durum
Admin panelde news_agencies modülü menüsü kontrol edilmeli.  
Eğer pers-assit sitesi context'i seçildiğinde ajanslar görünmüyorsa  
seed re-run sonrası görünecek.

### 6.2 Event Keys Dokümantasyonu
README standardı gereği, `news_agencies` modül list page'ine  
`EventKeysInfo` komponenti eklenmeli:

```tsx
const NEWS_AGENCY_EVENT_KEYS: EventKeyDefinition[] = [
  {
    key: "ingestion.run.started",
    description: "Ajans connector RSS/API çekme işlemi başladı",
    payload_fields: [
      { name: "connector_id", type: "string", required: true },
      { name: "dry_run", type: "boolean" },
    ],
  },
  {
    key: "ingestion.run.completed",
    description: "Ajans içerik çekme işlemi tamamlandı",
    payload_fields: [
      { name: "connector_id", type: "string", required: true },
      { name: "fetched_count", type: "number" },
      { name: "imported_count", type: "number" },
      { name: "failed_count", type: "number" },
    ],
  },
  {
    key: "provider.created",
    description: "Yeni ajans sağlayıcısı oluşturuldu",
    payload_fields: [
      { name: "entity_id", type: "string", required: true },
      { name: "name", type: "string" },
    ],
  },
];
```

---

## 7. Teknik Borç ve Riskler

| Risk | Önem | Durum | Çözüm |
|------|------|-------|-------|
| AA public RSS URL değişirse | Orta | Açık | Connector config'den okunur, seed güncellenmeli |
| ANKA RSS'i erişilebilir olmayabilir | Orta | Test edilmedi | `TestConnection` endpoint kullanılmalı |
| PA `GetNewsFeed` boş base_url kontrolü | Yüksek | Geçici fix var | Seed + permanent refactor |
| Connector RSS worker YOK | Yüksek | Açık | Faz 3 kapsamında |
| Flutter kimlik bilgisi backend'e gönderilmiyor | Orta | Tasarım kararı | Faz 2 endpoint'leri ile |

---

## 8. Bağlantılı Dokümanlar

- [Mevcut Temel Plan](KISISEL_ASISTAN_AJANS_AI_HABER_VE_SOHBET_GECMISI_PLAN_2026-06-24.md)
- [AI Assistant README](../README_FOR_AI_ASSISTANTS.md)
- [Integrations Platform Plan](AYSTEK_INTEGRATION_PLATFORM_ENTERPRISE_DELIVERY_PLAN.md)
- [Fleet Tally System Plan](FLEET_TALLY_SYSTEM_MASTER_PLAN_TR.md)
