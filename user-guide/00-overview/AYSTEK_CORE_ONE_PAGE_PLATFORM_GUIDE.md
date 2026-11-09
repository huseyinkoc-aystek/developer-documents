---
Title: AYSTEK Core One-Page Platform Guide
Version: 1.1.0
Status: Enterprise Certified - Production Ready
Last Updated: 2026-03-08
Classification: Internal - Enterprise Documentation
Related Documents:
  - ../README.md
  - ../01-architecture/MODULE_INDEX.md
  - ../02-security/AUTHENTICATION_AND_AUTHORIZATION.md
  - ../04-api/API_REFERENCE.md
  - ../05-operations/OPERATIONS_GUIDE.md
  - ../06-modules/INTEGRATIONS_PLATFORM_GUIDE.md
---

# AYSTEK Core: Tek Sayfada Detaylı Platform Rehberi

Bu doküman, AYSTEK Core platformunu tek sayfada ama detaylı şekilde anlatır.  
Amaç: teknik olmayan paydaşların da sistemi anlaması, teknik ekiplerin de doğru başlangıç referansına sahip olması.

## 1) Platform Özeti

AYSTEK Core; çok kiracılı (tenant), çok site’lı (site), modüler ve entegrasyon merkezli bir kurumsal işletim platformudur.

- **Backend:** Go microservice (Gin), MongoDB, Redis, RabbitMQ
- **Frontend:** Next.js admin panel (App Router, TypeScript, Tailwind)
- **Güvenlik:** JWT/cookie auth, RBAC, module permission, tenant-site isolation, CSRF
- **Gözlemlenebilirlik:** Execution trace, DLQ, SSE, metrics/health/ops
- **Ana değer:** içerik + ticaret + operasyon + entegrasyon süreçlerini tek yönetim yüzeyinde birleştirme

## 2) Mimari Katmanlar (Backend + Panel)

### 2.1 Go Katmanları
1. **Router/Middleware**
- Auth doğrulama
- Permission kontrolü
- Tenant/Site context doğrulama
- CSRF/Maintenance/Rate-limit guard

2. **Module Handler**
- HTTP input binding/validation
- Yetki ve hata sözleşmesi

3. **Module Service**
- İş kuralları
- Workflow, retry, idempotency, orchestration

4. **Repository**
- MongoDB kalıcı veri erişimi
- Sorgu, index, tenant filtre

5. **Platform Servisleri**
- Queue/RabbitMQ
- Redis cache/idempotency
- SSE/ops stream
- Monitoring/metrics/logging

### 2.2 Next Admin Katmanları
1. **Route/Layout guard**
- Public/private route ayrımı
- Session bootstrap doğrulama

2. **Context & Store**
- Auth store
- Tenant/Site selector
- Dil (i18n) store

3. **Page/Feature Screens**
- Modül bazlı list/detail/create/edit/test ekranları

4. **Common UI Kit**
- Table/filter toolbar
- Notice/modal
- JSONBuilder, form pattern, action controls

## 3) Kim Hangi Özelliği Kullanır?

### 3.1 Operasyon ve Yönetim
- Dashboard/performans metrikleri
- İşlem logları, audit izleri
- UAT/GoLive/deployment görünürlüğü

### 3.2 İçerik ve Yayın
- Posts, pages, taxonomy, media, galleries, sliders
- Localized content ve çeviri yönetimi
- Moderation/survey/forum süreçleri

### 3.3 Ticari Süreçler
- Catalog, products, orders, checkout, payments
- Pricing/discount/coupon/tax/fatura (invoices, ledger)
- Returns/shipping/post-purchase operasyonları

### 3.4 Stok ve Saha Operasyonları
- Inventory, vehicles, vehicle-assets, assignments
- Sales/returns/subscriptions/transfers
- Raporlama ve konum/bağlı varlık süreçleri

### 3.5 Entegrasyon ve Otomasyon
- Connectors + actions + hooks
- Chains (çok adımlı akış)
- Scheduled pull (zamanlanmış dış veri çekme)
- DLQ/retry/test-run/validation

## 4) Modül Aileleri ve Kritik Kapsam

### 4.1 P0/P1 Çekirdek (Sistem ayakta kalması için kritik)
- auth, users, permissions, settings, system
- dashboard, maintenance, performance
- menus, translations, languages

### 4.2 İçerik Ailesi
- posts, pages, media, taxonomy, references, authors
- content_bots, content_ops_settings, localized_content
- interviews, biographies, news_sources, news_agencies

### 4.3 Ticaret Ailesi
- catalog, products, carts, checkout, orders, payments
- pricing, discounts, shipping, taxes, returns
- installments, payment-links, payment-intents

### 4.4 Operasyon/Kurumsal Aile
- contacts, companies, contracts, approvals
- deployment, uat, golive, reports
- licenses, mailtemplates, activity_logs

### 4.5 Sektörel/Domain Aile
- stock (geniş ekosistem)
- pms (property/site operasyonları)
- crm, b2b, digital_editions

### 4.6 İletişim ve Topluluk
- notifications, sms, chatbot
- forum, comments, surveys

## 5) Uçtan Uca Çalışma Akışları

### 5.1 Admin Kullanıcı Akışı
1. Login -> token/cookie üretimi
2. `/me/context` ile oturum + tenant doğrulama
3. Tenant/site seçimi
4. Yetkili modül menülerinin yüklenmesi
5. Sayfa bazlı CRUD/test işlemleri

### 5.2 İş Akışı Akışı (Örnek)
1. Kullanıcı panelde işlem başlatır
2. API handler request alır
3. Service business rule çalıştırır
4. Repository yazma/okuma yapar
5. Gerekirse event publish edilir
6. Async işlem worker ile tamamlanır
7. UI SSE veya refresh ile sonucu görür

### 5.3 Hata Akışı
1. Geçici hata -> retry/backoff
2. Kalıcı hata -> DLQ
3. Operasyon ekibi DLQ ekranından triage/retry yapar
4. Execution kayıtları ile kök neden analizi yapılır

## 6) Integrations Platform (Detaylı Özet)

### 6.1 Neden Var?
İç sistem event’lerini dış servislerle güvenli, izlenebilir ve modüler biçimde bağlamak için.

### 6.2 Ana Nesneler
- **Connector:** Base/Test URL, auth, timeout, outbound policy
- **Action:** Method/path, header/body, mapping, response schema
- **Hook:** Event key + mode (sync/async) + target (action/chain)
- **Chain:** start/action/transform/branch/end node’ları ile çok adımlı akış
- **Execution:** her çalışma kaydı
- **Step Execution:** chain node bazlı kayıt
- **DLQ:** başarısız ve tekrar denenebilir işler
- **Pull Schedule:** cron/interval dış veri çekme planı
- **Pull Run:** schedule çalıştırma geçmişi + ingestion özeti

### 6.3 Payload Pipeline
1. Mapping
2. Static values
3. Enrichment (secondary lookup)
4. Defaults (boş/null fallback)
5. Build + send

### 6.4 Response Pipeline
- JSON/XML/SOAP/Text normalize
- Response schema koşulları ile success/fail kararı
- Mismatch/fault detayının execution metadata’ya yazılması

### 6.5 Güvenlik
- SSRF koruması
- Outbound allowlist
- Private IP erişim politikasının kontrollü yönetimi
- Auth header/secret yönetimi

### 6.6 Operasyon
- Test paneli (request/response preview)
- Validate endpointleri (mapping/chain)
- Retry/DLQ yönetimi
- Smoke runbook ile post-deploy doğrulama

## 7) Scheduled Pull (Dış Veri Çekme) Kapsamı

### 7.1 Özellikler
- Cron veya interval tetikleme
- Concurrency policy (`forbid/allow/replace`)
- Retry limit, timeout, run başına page/record limit
- Cursor/checkpoint ile incremental çekme
- Canonical mapping/upsert
- Per-run ingestion summary + issue taxonomy
- `publish_events` ile downstream event yayını

### 7.2 UI Kullanımı
1. Pull schedule oluştur
2. Connector + action bağla
3. Strategy’yi formdan doldur (JSON opsiyonel)
4. Dry-run ile test et
5. Pull runs detayından kalite ve hata kontrolü yap
6. Gerekirse run retry et

## 8) Zincir (Chain) Orchestration Kapsamı

### 8.1 Model
- Hook target: `action` veya `chain`
- Chain graph: DAG, cycle yok
- Node tipleri: start/action/transform/branch/end

### 8.2 Çalışma Semantiği
- Runner, mevcut action executor’u reuse eder
- Async semantik RabbitMQ + execution update + retry/DLQ ile korunur
- Branch kararları predicate ile verilir

### 8.3 İzlenebilirlik
- chain_run üst kayıt
- step_execution alt kayıt
- UI timeline + node odaklama + validation paneli

## 9) Güvenlik ve Uyum (Compliance) Özeti

### 9.1 Güvenlik Katmanları
- Kimlik doğrulama (token/cookie)
- Yetkilendirme (role + permission)
- Tenant/site izolasyonu
- CSRF koruması
- Güvenli outbound çağrı politikası

### 9.2 Uyum İlkeleri
- Audit log ve iz bırakma
- Hata kod standardı
- Çoklu dil metin standartları
- Seed tabanlı menü/izin/çeviri yönetimi

## 10) SLO/SLA ve Operasyonel Beklentiler

- API performansı: p95 hedefi operasyon rehberindeki SLA ile uyumlu izlenir
- Uptime: health/ready/metrics uçları ile sürekli takip
- Event pipeline: publish -> queue -> worker -> execution/dlq zinciri smoke test ile doğrulanır
- Dokümantasyon: her büyük feature sonrası user-guide + changelog güncellenir

## 11) Günlük Operasyon Checklist (Kısa)

1. Health ve metrics kontrolü
2. Kritik modül dashboard’larını kontrol
3. Integrations executions başarısız oranına bak
4. DLQ ve retry bekleyen işleri incele
5. Pull runs ingestion issue oranını kontrol et
6. Gerekirse targeted smoke test çalıştır

## 12) Sık Görülen Sorunlar ve Hızlı Çözüm

### 12.1 “Login olmuş gibi açılıyor” problemi
- Sebep: stale token/cookie + eksik bootstrap doğrulaması
- Çözüm: server-side `/me/context` doğrulaması ve 401/403’te login redirect

### 12.2 Integration test çağrısı beklenen body ile gitmiyor
- Action mapping + static + body template + test data senkronunu kontrol et
- Request preview ile gerçek payload’u doğrula

### 12.3 SOAP/XML response’da başarısızlık yakalanmıyor
- Response schema alanlarını doğru path/expectation ile tanımla
- Raw response ve normalized preview karşılaştırması yap

### 12.4 Scheduled pull beklenen veriyi işlemedi
- records_path/cursor_field/cursor_path kontrol et
- canonical mapping/defaults/enrichment kuralını doğrula
- Pull run detail’de ingestion summary/issue reason kontrol et

## 13) Yeni Ekibe Onboarding (Hızlı Plan)

### Gün 1
- Platform ve modül haritasını oku
- Tenant/site context akışını test et

### Gün 2
- CRUD + permission + translation seed akışını uygula
- Bir modülde küçük iyileştirme yap

### Gün 3
- Integrations action/hook test et
- Execution/DLQ akışını gözlemle

### Gün 4
- Chain ve pull schedule kurup dry-run yap
- Runbook adımlarıyla smoke test tamamla

## 14) Önerilen Doküman Okuma Sırası

1. `00-overview/OVERVIEW.md`
2. `00-overview/AYSTEK_CORE_ONE_PAGE_PLATFORM_GUIDE.md` (bu doküman)
3. `01-architecture/MODULE_INDEX.md`
4. `02-security/AUTHENTICATION_AND_AUTHORIZATION.md`
5. `04-api/API_REFERENCE.md`
6. `06-modules/INTEGRATIONS_PLATFORM_GUIDE.md`
7. `05-operations/OPERATIONS_GUIDE.md`

## 15) Sonuç

AYSTEK Core, “tek yönetim paneli + modüler Go backend + enterprise integrations runtime” yaklaşımıyla içerik, ticaret ve operasyon sistemlerini aynı platformda çalıştırır.  
Bu yapı; güvenlik, izlenebilirlik, geriye uyumluluk ve kullanıcı dostu yönetim prensipleri üzerine kuruludur.
