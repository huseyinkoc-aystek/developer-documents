# AYSTEK E‑Ticaret Platformu — Master İş Planı (Revize / Senior Seviye)

**Girdi dokümanları (özet):**
- “Aystek E‑Ticaret Platformu İçin Gelişmiş Sistem Gereksinimleri” (hedef kapsam)
- `AYSTEK_FULL_ANALYSIS_TR.md` (mevcut modül + endpoint envanteri)
- `AYSTEK_STORY_ANALYSIS_TR.md` (mimari akış + bağlam)
- `ECOMMERCE_REQUIREMENTS_MASTER_GAP_WORKPLAN.md` (Codex GAP planı)

**Amaç:** Mevcut AYSTEK Core (Go + Next) üzerinde, hedef gereksinimlere ulaşmak için *geri dönüşsüz*, uygulanabilir, ölçülebilir bir “master execution plan” üretmek.

---

## 1) Kuzey Yıldızı: “Satılabilir Platform” Tanımı (Definition of Ready / Done)

### 1.1 Satılabilir (Go‑Live) için minimum şartlar (P0)
Aşağıdaki 8 madde tamamlanmadan müşteri/tenant onboarding başlatılmamalı:

1. **Ödeme sağlayıcı entegrasyonları (TR)**: iyzico/PayTR prod callback + imza doğrulama + idempotency + retry/DLQ.
2. **Auth/Capture akışı**: provizyon + sonradan capture + iptal/refund (tam yaşam döngüsü).
3. **Ödeme linkleri**: siparişe bağlı ödeme linki üretimi + SMS/Email gönderimi + expire + tekil kullanım.
4. **Sipariş durum makinesi**: `created → awaiting_payment → paid → fulfillable/shipped → completed` + iade/iptal/chargeback dalları.
5. **Faturalama**: B2C/B2B fatura profili + otomatik fatura üretimi + revizyon/iptal senaryoları.
6. **Kargo**: en az 1 taşıyıcıyla gerçek etiket/track akışı + teslimat opsiyonları (kapıda ödeme şart değil).
7. **Multi‑store + Multi‑currency**: store’a göre para birimi + ülke tespiti + FX güncelleme + manuel kur/sabit kur kampanyası.
8. **Gözlemlenebilirlik**: audit log + event log + kritik metrik/alarmlar (payment webhooks, order pipeline, FX health).

### 1.2 Kalite kapıları
- **Contract Tests**: ödeme sağlayıcı webhook sözleşmeleri, marketplace webhooks, ERP sync.
- **Idempotency**: ödeme webhook, fatura üretimi, stok düşümü, iade/refund *mutlaka idempotent*.
- **Data Governance**: tenant isolation, PII masking, retention/rollup (özellikle event/usage metering).

---

## 2) Mevcut Envanter — Neyi “tamam sayıyoruz”? (Gerçekçi baseline)

Mevcut analiz dokümanlarına göre Core’da aşağıdaki çekirdek modül kümeleri var (Next tarafında da API wrapper’ları mevcut):
- Katalog: products, product‑variants, product‑media, categories/tags, catalog
- Ticaret: cart, orders, returns, discounts, pricing, installments
- Ödeme: payments, payment‑intents, payment‑gateways, payment‑links (+ public)
- Mali: invoices, taxes, currencies, units
- Lojistik: shipping, stock, warehouses
- Platform: stores, languages, translations, settings/public‑settings, plugins, permissions, users
- İşletme: contracts, contract‑contents, subscription‑contracts, subscriptions, licenses
- Entegrasyon: erp‑integrations, marketplaces, marketplace‑webhooks
- Destek: activityLogs, notifications, chatbot (intent schema + canned responses)
- Analytics/AI: dashboard, ai

Bu tablo “var” demek değil; **uçtan uca çalışıyor** demek hiç değil. Planın kalanında her epik için “Done” kriterleri veriliyor.

---

## 3) Faz Planı (P0→P3) ve Epikler

### FAZ P0 — Ödeme & Sipariş Omurgası (Go‑Live bloklayıcılar)
**Hedef:** Sipariş → Ödeme → Fatura → (Teslimat/Dijital erişim) hattının üretimde güvenli çalışması.

**P0‑E1: Payment Provider Adapters (iyzico/PayTR)**
- Webhook doğrulama: signature + timestamp + replay protection
- Idempotency key modeli (provider event id ↔ internal event)
- Webhook delivery log + retry + DLQ
- PaymentIntent lifecycle: authorize/capture/void/refund
**Kabul kriteri:** Aynı webhook 100 kez gelse de tek kez etkiler; tüm adımlar auditlenir.

**P0‑E2: Order State Machine & Idempotent Fulfillment**
- State machine tanımı + geçiş kuralları
- Stok düşümü, fatura üretimi, lisans teslimi “exactly‑once” mantıkta
- Concurrency: aynı order üzerinde çift işlem engeli (locking/optimistic version)
**Kabul kriteri:** Double‑charge, double‑invoice, negative stock gibi vakalar testlerde yakalanamaz hale gelmeli.

**P0‑E3: Payment Links (sipariş bazlı)**
- Order’a bağlı link, expire/one‑time, public endpoint
- Link üzerinden ödeme → order state update
- Email/SMS template + gönderim queue
**Kabul kriteri:** Sipariş sonrası linkle ödeme tamamlanınca otomatik fatura + teslimat tetiklenir.

**P0‑E4: Dijital Ürün / Lisans Teslimi**
- ProductType: digital/license/subscription ayrımı net
- Order(PAID) → License create → access grant
- Refund → revoke (kademeli revoke policy)
**Kabul kriteri:** Dijital ürün satışı tek tıkla çalışır; erişim loglanır.

**P0‑E5: Go‑Live Observability**
- Correlation ID: order_id, payment_intent_id, invoice_id, webhook_id
- Dashboard: payment webhook success rate, order pipeline lag, FX sync
- Alerting: webhook failure spike, capture failure, FX stale
**Kabul kriteri:** Prod sorunları “görünür” ve triage edilebilir.

---

### FAZ P1 — Ticari Özellik Derinleştirme (B2B + stok + vergi + kargo)
**Hedef:** Türkiye’de satış + ihracata hazır temel.

**P1‑E6: Multi‑Currency Pricing & FX**
- Default base currency + ek para birimleri
- Ülkeye göre fiyat gösterimi
- Manuel ürün bazlı döviz fiyatı (override)
- Sabit kur kampanyası (ör: 1 USD=1 TL)
- Checkout tahsilat kur çevrimi standardı
**Kabul kriteri:** Fiyat gösterimi + tahsilat kurları tutarlı; raporlar base currency’ye normalize.

**P1‑E7: Tax Rules (TR + uluslararası)**
- KDV dahil/haric gösterim (B2C vs B2B)
- Ülke/bölge bazlı vergi oranı
- İhracat: KDV muaf senaryo
**Kabul kriteri:** Sepette vergi hesapları doğru; fatura vergi kalemleri uyumlu.

**P1‑E8: B2B/Bayi**
- Company profile: vergi no, adres, yetkili, ödeme vadeleri
- Kredi limiti / onay mekanizması
- B2B pricing (segment/bayi grubu)
**Kabul kriteri:** B2B müşteride “net fiyat + vadeli satış” akışı çalışır.

**P1‑E9: Inventory & Warehouse Operations**
- Çoklu depo stok
- Low‑stock uyarı/notification
- Bulk import/export (CSV/Excel) + job tracking
- İade/iptal: stock rollback
**Kabul kriteri:** Stok hareketleri denetlenebilir; negatif stok engellenir.

**P1‑E10: Shipping Carriers MVP**
- 1‑2 taşıyıcıyla gerçek entegrasyon (etiket + takip)
- Teslimat opsiyonları (standart/ekspres/pickup)
- SLA tablosu (bölgeye göre)
**Kabul kriteri:** Sipariş shipped olunca müşteri track linki alır; ops panelde shipment timeline görünür.

**P1‑E11: Returns & Refunds**
- Return request lifecycle + reason codes
- Refund policy + partial refund
- Stock restock rules
**Kabul kriteri:** İade/refund akışları ödeme sağlayıcıyla uyumlu, muhasebe kayıtları bozulmuyor.

---

### FAZ P2 — Entegrasyonlar & Kurumsal Derinlik (ERP, pazaryeri, e‑Fatura)
**Hedef:** “Platform” satışını mümkün kılan entegrasyon yüzeyi.

**P2‑E12: ERP Integrations (Sync Jobs)**
- Mapping layer (product/order/invoice/customer)
- Delta sync + retry + dead letter
- Reconciliation raporu
**Kabul kriteri:** Bir ERP’ye minimum mapping ile bağlanıp tutarlı veri akıtılır.

**P2‑E13: Marketplace Integrations**
- Provider listesi + event mapping
- Catalog sync, order pull, stock push
- Webhook standardı + signature
**Kabul kriteri:** En az 1 pazaryeri ile MVP çalışır; geri kalanlar adapter backlog’ta netleşmiş olur.

**P2‑E14: e‑Fatura/e‑Arşiv Planı**
- Provider seçimi + adapter tasarımı
- GİB uyum süreçleri (prosedürel checklist)
**Kabul kriteri:** Mimari hazır, adapter stub + test harness var (tam prod zorunlu değil, ama yol temiz).

---

### FAZ P3 — AI, Analitik, Büyüme & Ürünleşme
**Hedef:** Rakiplerden ayrışma (operasyonel zeka + otomasyon).

- Demand forecasting, dynamic reorder suggestions
- Fraud/risk scoring (payment-risk)
- Content AI (ürün açıklaması, kategori metni)
- KPI/SLA setleri + otomatik aksiyon (workflow/automation)

---

## 4) Cross‑Cutting Epikler (Her fazda paralel yürür)

### X1: API Standardizasyonu (OpenAPI + Error Schema)
- OpenAPI v3 tek kaynak
- Standard error envelope + trace id
- Rate limit + pagination + filtering standardı
- Backward compatible versioning
**Done:** Her yeni endpoint otomatik dokümante + contract test.

### X2: RBAC + Audit + Compliance
- Permission matrix: ops, finance, catalog, integrations
- Audit log zorunlu: payment, order, invoice, stock, contract/license
- PII masking + retention policy

### X3: Event‑Driven Sağlamlık (Queue Topology)
- Event naming conventions
- Outbox pattern (DB→queue)
- DLQ ve reprocess UI

### X4: Seed & Localization Governance
- menu_seed + translation_seed zorunlu kapı
- Yeni modül: i18n coverage %100 hedef

---

## 5) Codex Planındaki Tipik Eksikler (Bu revizyonda kapatıldı)
Bu master plan, Codex GAP planında genelde “atlanan” ama prod’da can yakan kısımları özellikle ekledi:
- Webhook idempotency + replay protection + delivery log/DLQ
- Sipariş state machine’in stok/fatura/lisansla atomik bağlanması
- Multi‑currency: gösterim kuruyla tahsilat kurunun ayrıştırılması + rapor normalize
- B2B: onay/limit/vade gibi gerçek hayat gereksinimleri
- Observability + alarm seti (go‑live olmazsa olmaz)
- Entegrasyonlar için reconciliation (ERP/marketplace sync “tutarlılık raporu”)

---

## 6) Teslimat Paketleri (Repo/PR Disiplini)
Her epik için PR “paketi” şu üçlüden oluşmalı:
1) **Backend**: model/service/handler + test + migration
2) **Next Admin**: ilgili sayfalar + form validation + optimistic UI
3) **Ops**: dashboard kartı + alarm kuralı + runbook

---

## 7) Uygulama Sırası (Kesin Öneri)
1) P0‑E1 → P0‑E2 → P0‑E3 (ödemeyi sağlamlaştır)
2) P0‑E4 (dijital/lisans satışı)
3) P1‑E6 + P1‑E7 (multi‑currency & tax) — ihracat kapısı
4) P1‑E9 + P1‑E10 + P1‑E11 (operasyon hattı)
5) P1‑E8 (B2B) — satış stratejisine göre erkene çekilebilir
6) P2 entegrasyonları
7) P3 AI

---

## 8) “Bitti” (Go‑Live) Kontrol Listesi
- [ ] Prod ödeme: webhook success rate > %99, retry+DLQ çalışıyor
- [ ] Sipariş lifecycle: 50+ senaryo test seti yeşil
- [ ] Fatura: B2C/B2B doğru; refund/iptal senaryoları kapalı
- [ ] Kargo: shipment tracking uçtan uca
- [ ] Multi‑store: tenant isolation + store routing
- [ ] FX: stale değil; manuel override çalışıyor
- [ ] Audit log: kritik domainlerde zorunlu
- [ ] Runbook: ödeme sorunları, FX sorunları, sipariş stuck durumları

---

> Bu doküman “yapılacaklar listesi” değil; bir **platformu satılabilir hale getirme yürütme planıdır**.
