# P1‑E6 Multi‑Currency Pricing & FX — Gap Analizi ve İş Kırılımı

## 1) Mevcut Durum (Kod Taraması)
- **Currency/FX altyapısı**
  - `currencies` modülü: tenant bazlı currency listesi + default + FX rate kayıtları.
  - `fx_sync`: provider (static/http), base/quote listesi, fixed rate opsiyonu, scheduler/RunOnce.
  - `rateResolver`: tenant default currency + latest FX rate snapshot (base→quote).
- **Ürün/variant fiyatı**
  - `products` ve `variants` üzerinde **tek currency** ve **tek fiyat** alanı var.
  - `products/pricing` modülü mevcut (price history, currency, validity), **cart/checkout akışında kullanılmıyor**.
- **Cart → Order**
  - Cart item fiyatı `resolveUnitPrice` ile variant → product fallback.
  - Cart currency `resolveCurrency` ile variant → product → TRY.
  - Order oluştururken `base_currency` ve `fx_rate` set ediliyor (rateResolver ile snapshot).
  - Order totals yalnızca **quote currency** (cart currency) üzerinden tutuluyor.
- **Payment Intent**
  - Intent modelinde `base_currency`, `fx_rate` var.
  - Link/checkout üzerinden intent oluşturulurken FX snapshot **zorunlu değil** (çoğu yerde set edilmiyor).

## 2) Açıklar / Gaps

### G1 — Pricing Override (multi‑currency)
- `products/pricing` tablosu **cart hesaplamasında yok**.
- Per‑currency fiyat seçimi yok (currency bazlı override). 
- Variant‑level override yok (sadece ürün bazlı pricing olsa bile cart bundan habersiz).

### G2 — Base Currency Normalize (raporlama / muhasebe)
- Order totals sadece **quote currency** ile kaydediliyor.
- **Base currency toplamları** (`subtotal_base`, `grand_total_base` gibi) yok.
- İleride rapor/finans konsolidasyonunda normalize veri eksik.

### G3 — Tahsilat Kur Ayrımı (FX snapshot + ödeme)
- Checkout/Payment Link → PaymentIntent akışında **fx snapshot taşınmıyor**.
- “Gösterim kuru” vs “tahsilat kuru” ayrımı net değil.
- FX rounding policy (decimals/rounding) tanımlı değil.

### G4 — Currency Seçim Politikası (display vs charge)
- Cart currency, ürün/variant currency’ye kilitli.
- Store default currency veya müşteri seçimi akışa yansımıyor.
- Display currency değişimi için UI/endpoint yok.

## 3) P1‑E6 Teknik İş Kırılımı (Önerilen)

### E6‑T1: Currency Policy & Selection
- Cart için **currency resolution policy**:
  1) explicit request param/metadata (client seçimi)
  2) store default currency
  3) product/variant currency
  4) fallback TRY
- Cart’a currency set etmek için endpoint/param standardı.

### E6‑T2: Pricing Resolver (override)
- `products/pricing` **aktif fiyatı** (currency + validity) cart’a entegre.
- Öncelik sırası:
  1) active price (product_price) for cart currency
  2) variant price (currency match)
  3) product base price
- Opsiyonel: sale/discount price type desteği.

### E6‑T3: FX Snapshot Propagation
- Order oluşturulurken `base_currency` + `fx_rate` zorunlu snapshot.
- PaymentIntent oluşturulurken order’dan snapshot taşınması.
- PaymentLink → intent akışında base/fx set edilmesi.

### E6‑T4: Base Totals (Normalize)
- Order/Invoice üzerinde **base currency totals** alanları (subtotal/tax/grand).
- Rapor ve ledger için normalize hesap.

### E6‑T5: Rounding & Conversion Rules
- Currency decimals (currency table’dan) ile tutarlı rounding.
- FX conversion helper (amount * fx_rate, round). 

## 4) Öncelik Sırası (P1‑E6 için)
1) **E6‑T1** Currency policy (cart currency selection) + minimal UI/param
2) **E6‑T2** Pricing override (cart hesaplamasında aktif product_price)
3) **E6‑T3** FX snapshot propagation (order → payment intent)
4) **E6‑T4** Base totals normalize
5) **E6‑T5** Rounding/Conversion helpers

---

> Not: Bu adımların tamamı P1‑E6 kapsamındadır. Öncelikle policy + pricing override olmadan diğer adımların doğruluğu düşük kalır.
