# AYSTEK – Gelişmiş Sistem Gereksinimleri Uyum Master İş Planı

**Hedef doküman:** `/workspace/aystek-core-dev-docs/Aystek E-Ticaret Platformu İçin Gelişmiş Sistem Gereksinimleri.pdf`

**Karşılaştırılan kaynaklar:**
- `/workspace/aystek-core-dev-docs/AYSTEK_FULL_ANALYSIS_TR.md`
- `/workspace/aystek-core-dev-docs/AYSTEK_STORY_ANALYSIS_TR.md`
- `/workspace/aystek-core-dev-docs/README_FOR_AI_ASSISTANTS.md`

**Metodoloji:** PDF gereksinim başlıkları tek tek çıkarıldı ve mevcut repo envanteri (analiz dokümanları) ile eşleştirildi. Her madde için **Mevcut Durum**, **GAP**, **Gereken İş** ve **Öncelik** verildi. Bu plan **master seviyede** ve geriye dönüş gerektirmeyecek şekilde hazırlanmıştır.

---

## 1) GAP Matrisi (PDF Gereksinimleri → Mevcut Durum)

### 1.1 Ürün Yönetimi
- **Esnek Ürün Tipleri (fiziksel/dijital/abonelik/varyant)**
  - Mevcut Durum: `products`, `licenses`, `subscriptions`, `variants` modülleri mevcut.
  - GAP: Dijital teslim/erişim akışları (download/licence) + abonelik ekranları tam doğrulanmamış.
  - Gereken İş: Dijital erişim akışının UI/Backend doğrulanması, lisans teslim senaryoları ve erişim logları.
  - Öncelik: P0

- **Varyant ve Özelleştirme**
  - Mevcut Durum: Variants + attributes + pricing altyapısı var.
  - GAP: Varyant bazlı görsel/SEO ve bundle/aksesuar kurgusu uçtan uca doğrulanmadı.
  - Gereken İş: Varyant medya, bundle/aksesuar ilişkileri, admin UI doğrulama.
  - Öncelik: P1

- **Stok Yönetimi (kritik stok, iade/iptal geri stok, CSV/Excel)**
  - Mevcut Durum: `stock`, `warehouses` modülleri mevcut.
  - GAP: Kritik stok uyarıları + toplu import/export akışı net değil.
  - Gereken İş: Low-stock uyarı mekanizması + bulk import/export taskları.
  - Öncelik: P1

- **Kategori Sistemi (hiyerarşik, SEO)**
  - Mevcut Durum: `taxonomy`, `categories` yapıları var.
  - GAP: Ürün kategorileri içerik kategorilerinden ayrışık mı net değil.
  - Gereken İş: Ürün kategori yönetimi için ayrı akış/sekme ve SEO alanları doğrulama.
  - Öncelik: P1

- **Ürün Bilgi Yönetimi (multi-media, teknik özellik, SEO)**
  - Mevcut Durum: `products`, `product_content`, `media` modülleri var.
  - GAP: Teknik özellik/attribute kullanımının admin + frontend uyumu doğrulanmadı.
  - Gereken İş: Attribute/feature UI ve public ürün sayfası denetimi.
  - Öncelik: P1

### 1.2 Sipariş & Müşteri Yönetimi
- **B2B müşteri profilleri (şirket, vergi)**
  - Mevcut Durum: `b2b` ve `users` var.
  - GAP: B2B müşteri alanları ve onay akışı UI/Backend doğrulanmadı.
  - Gereken İş: B2B profil model + panel doğrulaması.
  - Öncelik: P1

- **Teklif/Taslak Siparişler**
  - Mevcut Durum: Net modül yok.
  - GAP: Quote/Draft order akışı eksik.
  - Gereken İş: `orders` altına draft/quote state + UI akışı.
  - Öncelik: P1

- **Siparişe özel sözleşme**
  - Mevcut Durum: `contracts` ve `subscription_contracts` var.
  - GAP: Order-specific contract ilişkilendirme akışı doğrulanmadı.
  - Gereken İş: Order → Contract mapping + admin UI doğrulama.
  - Öncelik: P1

- **Müşteri Hizmetleri Akışı**
  - Mevcut Durum: `activity_logs`, `comments` vb. var.
  - GAP: Müşteri hizmetleri ekranı/iş akışı belirsiz.
  - Gereken İş: Ops/customer service panel tasarımı.
  - Öncelik: P2

### 1.3 Ödeme Altyapısı
- **TR ödeme sağlayıcıları (iyzico/PayTR)**
  - Mevcut Durum: Gateway + intents altyapısı var.
  - GAP: Provider callback mapping + prod entegrasyon doğrulaması.
  - Gereken İş: Provider adapter gerçek API bağlama + webhook doğrulama.
  - Öncelik: P0

- **Yurt dışı sağlayıcılar (Stripe/PayPal/Wise)**
  - Mevcut Durum: Gateway altyapısı var, provider mapping kısmi.
  - GAP: Provider kapsamı net değil.
  - Gereken İş: Sağlayıcı listesi + adapter planı.
  - Öncelik: P2

- **Auth/Capture, Refund, Installment**
  - Mevcut Durum: intents + installments var.
  - GAP: Uçtan uca lifecycle test senaryoları eksik.
  - Gereken İş: Intent lifecycle test seti ve ops panel.
  - Öncelik: P1

- **PCI‑DSS + Tokenization politikası**
  - Mevcut Durum: Dokümantasyon taslağı var.
  - GAP: Uygulama checklisti + enforce adımları.
  - Gereken İş: PCI checklist + tokenization policy + audit log.
  - Öncelik: P1

### 1.4 Faturalama
- **Fatura profilleri (B2C/B2B)**
  - Mevcut Durum: `invoices` var.
  - GAP: UI doğrulaması.
  - Gereken İş: B2B/B2C fatura ekranları test.
  - Öncelik: P1

- **Otomatik fatura**
  - Mevcut Durum: Order sonrası fatura oluşturma var.
  - GAP: İş akışı doğrulama.
  - Gereken İş: Otomatik tetikleyiciler/testler.
  - Öncelik: P1

- **E‑Fatura entegrasyonu**
  - Mevcut Durum: Net değil.
  - GAP: GİB e‑fatura adapter yok.
  - Gereken İş: E‑fatura adapter + provider planı.
  - Öncelik: P2

- **ERP entegrasyonu**
  - Mevcut Durum: `erp_integrations` var.
  - GAP: Sync job + mapping kuralı net değil.
  - Gereken İş: ERP sync job tasarımı + mapping taskları.
  - Öncelik: P1

### 1.5 Kargo & Teslimat
- **Çoklu kargo firması**
  - Mevcut Durum: `shipping`, `shipping_carriers`, `delivery_options` var.
  - GAP: Sağlayıcı entegrasyonları kısmi.
  - Gereken İş: DHL/Yurtiçi/Aras/MNG/PTT/Sürat/Hepsijet entegrasyon planı.
  - Öncelik: P1

- **Tahmini teslimat / pickup / same-day / planlı**
  - Mevcut Durum: Opsiyon alanları var.
  - GAP: Algoritma/iş kuralları belirsiz.
  - Gereken İş: SLA tablosu + kurallar + UI görünürlük.
  - Öncelik: P2

- **Çoklu depo**
  - Mevcut Durum: `warehouses` var.
  - GAP: Bölgesel teslimat kurgusu net değil.
  - Gereken İş: Depo bölge eşlemesi + routing policy.
  - Öncelik: P2

### 1.6 Uluslararasılaşma
- **Çoklu dil**
  - Mevcut Durum: `languages`, `translations` var.
  - GAP: Ürün/katalog içerik doğrulaması.
  - Gereken İş: Çoklu dil içerik QA.
  - Öncelik: P1

- **Çoklu para birimi**
  - Mevcut Durum: FX sync mevcut.
  - GAP: Checkout tahsilat + fiyat gösterimi doğrulaması.
  - Gereken İş: Multi‑currency test ve UI doğrulama.
  - Öncelik: P1

- **Vergi kuralları**
  - Mevcut Durum: `tax_rules` var.
  - GAP: Ülke/bölge bazlı kurallar doğrulaması.
  - Gereken İş: Tax rule test seti + UI.
  - Öncelik: P1

- **Multi‑store**
  - Mevcut Durum: `stores` var.
  - GAP: Domain routing + store‑based FX sync doğrulaması.
  - Gereken İş: Store routing + FX store iterasyonu (kısmen tamamlandı).
  - Öncelik: P1

### 1.7 Entegrasyonlar
- **Pazaryeri**
  - Mevcut Durum: `marketplaces` var.
  - GAP: Provider mapping net değil.
  - Gereken İş: Provider‑specific event mapping + sync job taskları.
  - Öncelik: P1

- **CRM & Pazarlama**
  - Mevcut Durum: `crm`, `marketing` var.
  - GAP: Webhook + sync scheduler kısmi.
  - Gereken İş: CRM sync scheduler + webhook retry policy.
  - Öncelik: P2

- **Chatbot / Canlı destek**
  - Mevcut Durum: Chatbot modülü var.
  - GAP: Admin panel bütünlüğü sağlandı, UX testleri yapılmalı.
  - Gereken İş: Intent schema + canned responses QA.
  - Öncelik: P1

- **Sosyal ticaret**
  - Mevcut Durum: Placeholder.
  - GAP: Connector yok.
  - Gereken İş: MVP connector tasarımı + katalog sync.
  - Öncelik: P2

- **Developer API/Webhook standardı**
  - Mevcut Durum: Standart dokümanları oluşturuldu.
  - GAP: Endpoint envanteri normalize edildi, fakat OpenAPI + delivery log uygulaması yok.
  - Gereken İş: OpenAPI v3 spec + webhook delivery log implementasyonu.
  - Öncelik: P1

### 1.8 AI & Performans
- **AI karar destek / öneri**
  - Mevcut Durum: AI modülü + analytics var.
  - GAP: Tahmin, anomali ve öneri kapsamı açık.
  - Gereken İş: AI roadmap tasks.
  - Öncelik: P2

- **Performans / Ops KPI**
  - Mevcut Durum: KPI metrikleri ve dashboard ops kartları güncellendi.
  - GAP: Alerting thresholds ve SLA hedefleri.
  - Gereken İş: Alert policy + ops panel takibi.
  - Öncelik: P1

---

## 2) Master İş Planı (Uçtan Uca)

### Faz 0 – Bloklayıcılar (P0)
1. Ödeme provider callback mapping + webhook doğrulama (iyzico/PayTR)
2. Dijital teslim/erişim akışı (license/download) doğrulama
3. Dev API/Webhook standardı: idempotency + webhook signature enforcement

### Faz 1 – Core E‑Commerce (P1)
1. Draft/Quote order akışı
2. Order‑contract mapping + UI
3. Stok kritik uyarılar + bulk import/export
4. Multi‑currency checkout doğrulama
5. Tax rules QA + UI doğrulama
6. Marketplace provider mapping + sync job taskları
7. Kargo sağlayıcı entegrasyon planı (DHL/Yurtiçi/Aras/MNG/PTT/Sürat/Hepsijet)
8. Chatbot admin panel QA (intent + canned)
9. OpenAPI v3 spec skeleton

### Faz 2 – Advanced / B2B / Entegrasyonlar (P2)
1. E‑Fatura entegrasyonları (GİB)
2. CRM sync scheduler + webhook retry
3. Sosyal ticaret MVP connector
4. Çoklu depo routing policy
5. AI tahmin/anomali roadmap

### Faz 3 – Performans / Observability
1. Alert thresholds (P95, error rate, SSE drops, FX health)
2. Ops KPI panelleri SLA hedefleri

---

## 3) Çıktı ve Kontrol Listesi
- Her yeni key için **translation_seed.go** (28 dil).
- Her yeni modül/bağlantı için **menu_seed.go**.
- Yeni endpointler için **Developer API standardı** (error schema + idempotency + rate-limit).
- Release öncesi Go build + smoke test.

