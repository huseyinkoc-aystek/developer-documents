# balci.aysteknolojileri.com — Rucher du VELAY public sitesi iş planı

Son güncelleme: 12 Ekim 2026.
Durum: **Faz 0 tamamlandı. Faz 1 (Core backend) büyük ölçüde tamamlandı (6 Ekim 2026): storefront kataloğu, misafir/hesap sepeti, sunucu-fiyatlı checkout, site ödeme yöntemleri (kart/PayPal/SEPA), gerçek PayPal adaptörü ve webhook sertleştirmesi yazıldı; gerçek MongoDB + HTTP ile doğrulandı. Faz 2 (Next.js frontend, `/workspace/node/balci-aysteknolojileri-com`) misafir akışı için tamamlandı ve tarayıcıda (Playwright) doğrulandı. Sırada: Stripe/PayPal canlı-sandbox doğrulaması (anahtarlarla), müşteri hesapları, Faz 3 (installer).**
Kapsam netleşti (6 Ekim 2026; ikinci netleştirme: Fransa'daki bir üreticinin ürünlerini pazarlayan tanıtım + e-ticaret sitesi, havale/kredi kartı/PayPal, kayıtlı **ve** kayıtsız alışveriş): Figma çıktısı yalnız **şablon**; hedef Core'a tam entegre yeni bir **tanıtım + e-ticaret** sitesi.

## Kullanıcı talebi

- Figma Make çıktısı `/workspace/node/Rucherduvelayaricilik-main` (Vite + React + Tailwind v4 + MUI + react-i18next,
  mock veri) yalnız tasarım şablonudur; birebir uygulama değil, **yeni proje** açılır.
- Tenant: AYSTEK Mühendislik (`aystek`). Yeni site kodu: `balci`. Domain: `balci.aysteknolojileri.com`.
- `aysteknoloji-com` / `aystekhaber-com` ile aynı model: Next.js public site + Core backend.
- Marka: Rucher du VELAY (Haute-Loire arıcılığı), ana dil Fransızca, çok dilli genişlemeye açık.

## Keşif bulguları (doğrulanmış)

| Konu | Bulgu |
| --- | --- |
| Figma kaynağı | `src/data/homepage.ts` + `App.tsx` içinde sabit içerik (hero, bölgeler, ürünler, blog, değerler, güven rozetleri, konum), Unsplash URL'leri, `fr/en/tr.json` UI çevirileri, sepet metinleri. Hepsi README "Dynamic/Admin-Managed" kuralını ihlal eder → Core'a taşınmalı. |
| Stack farkı | Figma: Vite, Tailwind v4, MUI 7, react-i18next. Balcı üretim uygulaması: Next.js 15.5.27 App Router, React 19.3, Tailwind 3.4, TS 5.3+, kendi i18n bağlamı. Doğrudan kopyalanamaz, **port** edilir. |
| Referans site | `aysteknoloji-com`: `lib/api/server.ts` (X-Store-Domain/X-Tenant-ID), `lib/api/{pages,posts,languages,navmenus,siteProfile,newsletter,contacts,businessHours}.ts`, `components/cms`, `hooks/useCmsPage.ts`, Dockerfile, `scripts/check-env-contract.sh`. |
| Site kimliği | `aystekhaber-com` multi-domain çalışma planı: Host yalnız admission, tenant/site `AYSTEK_STORE_DOMAIN` + `AYSTEK_TENANT_ID` ile runtime kimliğinden gelir. |
| Backend site kaydı | `internal/seeds/news_sites_seed.go` desteği haber sitelerine özel (kategori/gazete/yazar seed'i taşır). Genel bir "public site" seed'i yok → ayrı, haber olmayan bir seed gerekir. |
| **Ürün public API (KRİTİK)** | `products` modülünde **yalnız admin** rotaları var (`/admin/products`, variants, media). `/public/catalog` hizmet/ürün *vitrin kataloğu* (`CatalogItem`: tek fiyat, varyantsız) — sepet ise `variant_id` ister. Yani varyantlı ticari ürünler için **public vitrin (storefront) API yok**. |
| **Sepet kimlik doğrulama (KRİTİK)** | `carts` ve `checkout` `protectedGroup` altında; `mustUserID` zorunlu → **misafir sepeti/ödemesi yok**. Karar: kayıtsız alışveriş **zorunlu** → yeni backend işi (Faz 1): anonim sepet kimliği (imzalı, süreli, site-scoped) + misafir checkout (e-posta + adres, `orders.CustomerEmail` zaten var). Kayıtlı müşteri akışı (`/auth/register`, `login-by-email`) korunur; girişte misafir sepeti kullanıcı sepetine birleştirilir. |
| Ödeme sağlayıcıları | `payments/intents/service/adapters.go` içinde **Stripe** ve **PayPal** adaptörleri var (env'den, `GatewayAdapter` arayüzü: Create/Confirm/Capture/Cancel/Refund/VerifyCallback); Iyzico/PayTR/TR bankaları Fransa için gereksiz. Havale = `MethodBankTransfer` (çevrimdışı yöntem). Mollie/Adyen `ProviderType` olarak tanımlı ama adaptörü yok. **Kart = Stripe, PayPal = PayPal, Havale = banka bilgisi + manuel onay** önerilir; **KRİTİK: `PayPalAdapter` bir stub** (`"mode":"stub"`; PayPal API'sine hiç istek atmaz, ID'yi yerelde üretir, Capture hep `succeeded` döner) → canlıda para tahsil etmez. Gerçek Orders v2 (create/approve/capture/refund + webhook imza doğrulama) adaptörü Faz 1'de yazılacak. Stripe adaptörü gerçek API çağırıyor (Faz 1'de uçtan uca doğrulanacak). |
| Ödeme/sipariş | `payments` public yüzeyi yalnız `GET /payments` (kendi ödemelerim); sipariş yaratma sepet `POST /cart/checkout` üzerinden. Ödeme niyeti (payment intent) için müşteri-yüzü uç ve sağlayıcı (Stripe vb.) bağlantısı Faz 4'te doğrulanacak. |
| Kargo | `shipping` public yüzeyi yalnız `track/:tracking_number`; teslimat seçeneği/teklif ucu admin'de. Müşteri-yüzü teklif ucu gerekebilir. |
| Mevcut `balci` izi | Core, installer ve admin panelde `balci` geçen dosya yok (greenfield). |

## Mimari kararlar (öneri, aksi belirtilmedikçe uygulanır)

1. **Yeni, bağımsız repo:** `/workspace/node/balci-aysteknolojileri-com` (Next.js 15.5.27 + React 19.3 + Tailwind 3.4 + TS). Figma klasörü salt-okunur tasarım kaynağıdır, değiştirilmez.
2. **Tek içerik kaynağı Core:** `homepage.ts`, `BLOG_POSTS_BY_LANG`, `locales/*.json`, Unsplash sabitleri portta **taşınmaz**; Core'daki Pages / Posts / Settings / Languages / Menus / Newsletter / Contacts / Business-hours modellerinden okunur. Ana sayfa, `aystekhaber` gibi, `anasayfa` sayfasının dil bazlı JSON'u olur (hero, hikaye, bölgeler, değerler, güven rozetleri, etiketler).
3. **Diller:** `/public/languages/active` (RULE #1.1). Store `default_language=fr`; site override `fr, en, tr` (Figma'daki 3 dil), tenant kataloğundan. Kodda sabit dil listesi yok.
4. **Kapsam:** tanıtım vitrini + blog + iletişim/bülten + **e-ticaret** (ürün, sepet, ödeme, sipariş). E-ticaret Faz 4'tür ama Core'da eksik public yüzey nedeniyle backend işi Faz 1'de başlar (aşağıdaki kritik bulgu).
5. **Tasarım sadakati:** Figma token'ları (`brand-brown/green/beige`, Playfair Display + Inter) `tailwind.config.ts`'e taşınır; bileşen hiyerarşisi (layout/sections/products/shared/ui) korunur, MUI bağımlılığı alınmaz.
6. **Güvenlik/izolasyon:** Site kapsamı (tenant+site) dışında veri sızması yok; taxonomy tenant-geneli olduğundan "tüm kategorileri listele" ucu eklenmez, menüden türetilir.


## Fransa ödeme ve yasal uyum standardı (6 Ekim 2026 ek istek)

Ödeme sağlayıcıları Core'un mevcut `payments/gateways` (admin panelden yönetilen gateway kaydı) ve
`payments/intents/service` (`GatewayAdapter` arayüzü) yapısına eklenir; ayrı bir ödeme modülü veya
sitede sabit anahtar/kod yazılmaz. Anahtarlar admin panelden/gateway kaydından gelir (README "Dynamic" kuralı).

Seçim ölçütü: Fransa'da en yaygın kullanılanlar ve tek üreticinin kendi ürünlerini satması (marketplace değil).

| Öncelik | Sistem | Kapsadığı yöntemler | Uygulama |
| --- | --- | --- | --- |
| **1 (Faz 1)** | **Stripe** | Cartes Bancaires (CB), Visa, Mastercard, American Express, Apple Pay, Google Pay | Mevcut adaptör gerçek API çağırıyor; Payment Intents + SCA uçtan uca doğrulanır |
| **1 (Faz 1)** | **PayPal** | PayPal bakiyesi/hesabı | Stub → gerçek Orders v2 adaptörü |
| **1 (Faz 1)** | **SEPA Credit Transfer (virement)** | IBAN/BIC havale | Çevrimdışı yöntem: `awaiting_payment`, banka bilgisi site ayarından, benzersiz referans, admin onayı |
| 2 (sonra) | **Mollie**, **PayPlug** | CB odaklı Fransız/Avrupa alternatifleri | Aynı `GatewayAdapter` arayüzüyle yeni adaptör (`ProviderMollie` zaten tanımlı; PayPlug yeni `ProviderType`) |
| 2 (sonra) | SEPA Direct Debit, **Wero** (Avrupa hesaptan hesaba) | Abonelik/tekrarlı ödeme, anlık banka ödemesi | Talep olursa; Wero Fransa'da yeni yaygınlaşıyor, sağlayıcı desteğine bağlı |
| Kapsam dışı | Adyen, Checkout.com, Worldline, Lyra/PayZen | Kurumsal acquirer'lar | Gerekirse Faz 2 adaptörleri |
| Kapsam dışı | Mangopay, Lemonway, Stripe Connect, Adyen for Platforms | Marketplace (çok satıcılı) | Tek üretici modeli için gereksiz |

CB, Visa, Mastercard, Amex, Apple Pay, Google Pay ayrı entegrasyon değildir; Stripe üzerinden yöntem olarak açılır.
Hangi yöntemin sitede göründüğü admin panelde gateway kaydı ve site ayarlarıyla yönetilir (kodda sabit liste yok).

Uyum gereklilikleri (ödeme/yasal):
- **PSD2 / SCA (3-D Secure 2):** kart ödemesi SCA tamamlanmadan `succeeded` sayılmaz; `three_ds` akışı `intents` servisinde zaten modelli.
- **Webhook güvenliği:** Stripe/PayPal imza doğrulaması zorunlu, olay idempotent işlenir (çift tahsilat/çift sipariş yok).
- **Kart verisi sitede tutulmaz** (PCI kapsamı küçük: Stripe Elements / PayPal hosted akış); loglara kart/anahtar yazılmaz.
- **TVA (KDV):** Fransa oranları `taxes` modülünden (gıda %5,5); fiyatlar TTC gösterilir, fatura/sipariş TVA kırılımı içerir.
- **Droit de rétractation (14 gün):** gıda ürünleri istisna olabilir; metin/akış hukuk onayıyla; `returns` modülü.
- **Zorunlu sayfalar:** Mentions légales, CGV (satış koşulları), politique de confidentialité (RGPD), çerez onayı, ödeme/iade/teslimat bilgisi, bal için etiket/menşe bilgisi.
- **Misafir sipariş verisi:** RGPD uyumlu saklama süresi ve silme talebi.

## Aşamalar

### Faz 1 — Core: site kaydı, seed ve storefront public API (backend)
- [x] **`storefront` modülü** (`internal/modules/storefront`): ürün/varyant/medya tenant-geneli master olarak kalır (bu modüller `site_id` taşımaz); siteye yayın **`Listing`** kaydıyla (tenant+site+ürün, tekil indeks) yapılır → kardeş siteler birbirinin ürününü görmez. Yayın kapısı: ürün aktif + ≥1 aktif varyant. Public: `GET /public/storefront/products` (sayfalama, `q`, `featured`, `facet[grup]=değer`), `GET /public/storefront/products/:slug`; stok sayısı sızdırılmaz (`in_stock|low_stock|out_of_stock`). Admin: `/admin/storefront/listings` CRUD. Facet'ler admin tanımlı (bölge/sezon/kategori sabit değil). Servis testleri geçti; `main.go` wiring + indeks; `storefront` izni + editor rolü seed'lendi.
- [x] **storefront admin arayüzü** (`aystek-core-next`: `/admin/storefront` liste + `create` + `[id]/edit`, `components/storefront/*`, `lib/api|types/storefront.ts`; gradient-hero standardı, `AdminDataTable`/`TableToolbar`/`NoticeModal`); `menu_seed` girdisi (`storefront`, Ticaret > Ürünler); 45 anahtar × 28 dil çeviri (`storefront_translations_seed.go`, `storefront.` önekleri `fullyTranslatedPrefixes`'te, kapsam testi geçti); `check-i18n` geçti, lint temiz. **Kalan:** `balci` view profile'ı (RULE #2.2) site seed'iyle birlikte.
- [x] **Site seed'i** (`internal/seeds/public_sites_seed.go`, haber olmayan genel public-site ailesi; `SeedNewsSites` ile karıştırılmadı): `stores` (`aystek/balci`, `balci.aysteknolojileri.com`, `fr`, `EUR`, `Europe/Paris`, `public_site_ready=false`), `anasayfa` sayfası (fr/en/tr, `internal/seeds/data/balci_home_*.json` gömülü şablonlar — Figma metinlerinden, sahte iletişim verisi yok), site dil override'ı (fr/en/tr), boş sosyal bağlantı ayarı, `balci_shop_manager` görünüm profili (`balciVisibleMenuKeys`, anahtarların `allMenuKeys`'te varlığı testli). Hepsi insert-only (admin düzenlemesi/silinmiş kayıt ezilmez); ortak mantık `seedSiteStore`/`seedSiteViewProfile`/`upsertSiteSeedDefaults` olarak haber seed'iyle paylaşılır (kopya yok; haber testleri geçti).
- [x] **Demo katalog** (opt-in, `[DEMO]` etiketli 6 bal × 2 varyant + listing): `go run ./cmd/seed --public-site=balci --demo-catalog`. Normal kurulum: `go run ./cmd/seed --public-site=balci` (veya standart seed zinciri). Ön koşul: aktif `aystek` tenant.
- [x] **Doğrulama:** `tenantfilter` + storefront servis/handler testleri; gerçek MongoDB ile seed idempotans/izolasyon testleri (listing başka sitede görünmez); API'yi izole portta başlatıp `curl` ile `/public/storefront/products` (liste, sayfalama, `facet[region]`, dile göre slug, 404) ve yetkisiz admin (401) kontrolü; rota çakışması yok. **Bulunan ve düzeltilen hata:** MongoDB kısmi indeksi `$exists:false` desteklemiyor → listing soft-delete yerine gerçek silme (`Delete`), tekillik indeksi sade; `main.go` artık indeks hatasını logluyor.
- [ ] Faz 3'e kalan: installer public-site profili ve 3 dokümanın senkronu (`--public-site=balci` seed adımı installer'ın seed sırasına işlenecek).
- [ ] **Misafir alışveriş:** anonim sepet + misafir checkout (müşteri hesabı olmadan sipariş, e-posta ile sipariş takibi, girişte sepet birleştirme); IP/oran sınırlama ve site izolasyonu testleri.
- [ ] **Ödeme yöntemleri (Fransa standardı bölümüne bakın):** Stripe uçtan uca (CB, SCA), gerçek PayPal adaptörü + testleri (httptest ile sahte PayPal/Stripe), webhook imza/idempotans testleri; havale (IBAN bilgisi site ayarından, sipariş `awaiting_payment` → admin onayı), kart (Stripe), PayPal; müşteri-yüzü payment-intent ve callback akışı.
- [ ] Müşteri-yüzü kargo teklifi ve ödeme niyeti uçlarının doğrulanması/eklenmesi (mevcut servisleri sarmalayarak).
- [ ] Müşteri kayıt/giriş akışı site-scoped doğrulama (FindPersonByLinkedUser dersi: kişi-scope).
- [ ] `stores` kaydı: `aystek/balci`, `balci.aysteknolojileri.com`, `fr`, `EUR`, `Europe/Paris`.
- [ ] Haber olmayan genel `public_sites_seed.go` (site tanımı + insert-only varsayılanlar; `cmd/seed --public-site=balci`; `SeedNewsSites` ile karıştırılmaz).
- [ ] `anasayfa` sayfası (fr/en/tr) — Figma metinlerinden; menüler (Miels, Produits, Notre histoire, Blog, Contact); `site.{id}.public.*` ayarları (sosyal, iletişim, konum); site dil override'ı.
- [ ] RULE #2.2: `balci_editorial` view profile + allowlist testi; menü anahtarları `allMenuKeys` ile uyumlu.
- [ ] Ürün yönetimi için Core ürün public API doğrulaması (Açık Karar 1).
- [ ] Testler: kimlik/izolasyon/idempotans (news_sites_seed_test emsali), `go vet ./...`, `go run` ile sunucu başlatma.

### Faz 2 — Frontend portu
- [ ] Proje iskeleti (aysteknoloji-com'dan `lib/api/*`, SEO, env sözleşmesi, Dockerfile).
- [ ] Layout (Header/Footer/Container/Section) ve section'lar Core JSON'una bağlı olarak.
- [ ] Rotalar: `/`, `/miels`, `/produits`, `/histoire`, `/blog`, `/blog/[slug]`, `/contact`, hukuki sayfalar.
- [ ] SEO: metadata, canonical, sitemap (yalnız sitenin dilleri), robots, manifest, JSON-LD (LocalBusiness).
- [ ] Erişilebilirlik (WCAG) ve mobil doğrulama; Chromium masaüstü/mobil testi.

### Faz 3 — Kurulum ve yayın
- [ ] Installer public-site profili (`docs/examples/balci-*.json`, port 3003), env şablonu.
- [ ] 3 dokümanın senkron güncellemesi: installer `KULLANIM_KILAVUZU.md`, `SUNUCU_KURULUM_REHBERI.md`, bu plan.
- [ ] README "Active/Planned Work Trackers" listesine bu plan eklendi (✅ yapıldı).
- [ ] DNS/TLS/Nginx, gerçek içerik, logo/favicon, hukuki metinler (kullanıcıdan).

### Faz 4 — E-ticaret (frontend + entegrasyon)
- [ ] Ürün listeleme/detay, sepet, hesap (kayıt/giriş), checkout, ödeme, sipariş takibi; Core `carts/checkout/orders/payments/shipping/discounts/taxes`.
- [ ] Fransa KDV (taxes), kargo (shipping), iade (returns), kupon (discounts) site ayarlarından.
- [ ] Admin panelde ürün/varyant/stok yönetimi mevcut `/admin/products` ile; yeni admin ekranı gerekmeden.

## Güvenlik bulguları (Core, bu çalışmada tespit edildi)

| # | Bulgu | Durum |
| --- | --- | --- |
| 1 | **Stripe webhook imzası doğrulanmıyordu:** `VerifyCallback` imza sırrı boşsa doğrulamayı atlıyordu; admin panelden tanımlanan Stripe geçitleri `webhook_secret`'i hiç okumuyordu. Alıcı kendi PaymentIntent ID'sini `client_secret`'ten bildiği için sahte bir `payment_intent.succeeded` olayıyla siparişi ödenmiş yapabilirdi. | **Düzeltildi:** her zaman doğrulanır (sır yoksa reddedilir), geçit `webhook_secret` okur, geçit kaydı `webhook_secret` (Stripe) / `webhook_id` (PayPal) olmadan kaydedilemez. Mevcut Stripe/PayPal geçit kayıtlarına bu alanlar **eklenmeli**, yoksa webhook'lar reddedilir. |
| 2 | `PayPalAdapter` stub'tı (PayPal'a hiç istek atmıyor, `Capture` hep `succeeded`). | **Düzeltildi:** Orders v2 (oluştur/onay/yakala/iade/imza doğrulamalı webhook), sahte PayPal sunucusuyla testli. Canlı/sandbox hesapla doğrulanmadı. |
| 3 | `mock` ödeme adaptörü her zaman kayıtlıydı ve sağlayıcı boşsa varsayılandı; genel ödeme linki sağlayıcıyı istemciden alıyor. | **Düzeltildi (kullanıcı onayı, 6 Ekim 2026):** `APP_ENV=production\|prod`'da `mock` hiç kaydedilmez; dev/test/uat/staging değişmedi. **Operasyon notu:** üretimde `APP_ENV=production` olmalı; UAT ortamı `uat`/`staging` kullanmalı. |
| 4 | `/cart/checkout` kargo/indirim/KDV tutarlarını istemciden alabiliyor. | Storefront'ta **kapalı** (sunucu hesabı, istekte tutar alanı yok). Mevcut `/cart/checkout` ucu değiştirilmedi. |
| 5 | Sipariş numarası `ORD-<saniye>`: aynı saniyede çakışır ve tahmin edilebilir. | Storefront rastgele `BALCI-yymmdd-XXXXX` üretir (`CheckoutInput.OrderNumber`). Genel üretici değiştirilmedi. |
| 6 | Core KDV tutarı yuvarlanmıyor (`2.4533…`). | Storefront görünümleri kuruşa yuvarlar; saklanan sipariş değeri yuvarlanmamış kalır (fatura için muhasebeyle netleştirilmeli). |
| 7 | Sepet servisinin `With*` kurucuları (mağaza deposu, fiyatlandırma, B2B bayi) `main.go`'da hiç çağrılmıyor. | Storefront para birimini ilk satırdan alır, mağaza deposuna bağımlı değil. Fiyatlandırma/bayi override'ları şu an pasif. |

## İlerleme (Faz 1 — Core backend)

- [x] `storefront` modülü: yayın (Listing), katalog API, admin API/UI, menü, çeviri.
- [x] **Sepet:** misafir oturumu (`X-Cart-Token`, 32 bayt rastgele, yalnız SHA-256 saklanır, TTL indeksi, site-scoped) + hesap sepeti (`/me/storefront/cart`) + `claim` (girişte birleştirme). Yalnız bu sitede yayınlı ürünlerin varyantı eklenebilir; satır/adet sınırları; ilk satırın para birimi.
- [x] **Checkout:** `POST /public/storefront/checkout` (+ `/me/storefront/checkout`): kargo/KDV/indirim sunucuda, satırlar satın alma anında yeniden doğrulanır (yayın, stok, fiyat), şartlar kabulü siparişe yazılır, oturumsuz sipariş sorgu (numara+e-posta), sipariş yalnız sahibi ve kendi sitesince okunur.
- [x] **Ödeme yöntemleri:** site-scoped `PaymentOption` (card/paypal/bank_transfer; IBAN mod-97 doğrulaması), admin CRUD + UI + 38×28 çeviri + menü (`storefront-payments`), `GET /public/storefront/payment-methods` (kimlik bilgisi yok).
- [x] **Ödeme:** `POST orders/:id/payment` (sağlayıcı/tutar/dönüş URL'si sunucuda; yalnız `client_secret`/yönlendirme tarayıcıya gider), `payment/confirm` (PayPal dönüşü), havale talimatı (IBAN, referans = sipariş no).
- [x] **Seed:** Fransa KDV (ürün %5,5, kargo %20, TTC; **muhasebeci onayı gerekli**), pasif ödeme seçenekleri; demo: kargo seçenekleri ve SEPA demo hesabı (resmi örnek IBAN).
- [ ] Stripe uçtan uca: gerçek/test anahtarlarla PaymentIntent + `payment_intent.succeeded` webhook'u (şu an sahte sunucusuz yalnız adaptör birim testleri).
- [ ] Havale onayı: yönetici siparişi "ödendi" yapar (mevcut `/admin/orders` durum akışı); toplu eşleştirme (banka ekstresi) kapsam dışı.
- [ ] İade/cayma (`returns`), kupon yönetimi UI'sı, sipariş e-postaları (`mailtemplates`) — Faz 4.
- [x] Core `mock` adaptörü üretimde kapatıldı (bulgu 3).

## İlerleme (Faz 2 — frontend, 6 Ekim 2026)

Yeni proje: `/workspace/node/balci-aysteknolojileri-com` (Next 15.5.27, React 19.3, Tailwind 3.4, TS, GitHub `master`).
- [x] Figma şablonu **taşındı** (zinc+amber, Playfair Display + Inter): header/footer, hero, güven çubuğu, öne çıkan ürünler, bölgeler (facet `region` filtresine bağlı), hikâye, değerler, blog, iletişim, bülten. **Sabit içerik yok:** metinler/etiketler Core Pages'taki ana sayfa JSON'undan (`labels.shop.*` dahil, fr/en/tr), ürünler `/public/storefront/*`, blog Core Posts, yasal sayfalar `/p/[slug]` (Core Pages), logo/sosyal Core ayarlarından.
- [x] Katalog (`/miels`, `/produits`: arama + bölge/mevsim filtresi, sayfalama), ürün (varyant seçimi, stok durumu, JSON-LD), sepet, misafir checkout, sipariş sayfası (SEPA talimatı, Stripe Payment Element, PayPal yönlendirme + dönüşte onay/yoklama), sipariş sorgu (`/suivi`), iletişim/bülten.
- [x] Mimari: tarayıcı → `app/api/shop/*` (BFF) → Core. Sepet token'ı `httpOnly` çerez; tarayıcıdan fiyat/kargo/sağlayıcı/dönüş URL'si kabul edilmez; mutasyonlarda aynı-origin kontrolü; bilinmeyen Host 421; runtime kimlik (`SITE_URL`, `AYSTEK_STORE_DOMAIN`, `AYSTEK_TENANT_ID`).
- [x] SEO: metadata + canonical + hreflang, LocalBusiness/Product JSON-LD, sitemap (yalnız sitenin dilleri), robots, manifest.
- [x] Doğrulama: `tsc`, `next lint`, `next build`, 7 birim testi; gerçek Core+MongoDB'ye karşı Playwright (masaüstü + mobil): ürün → sepet → checkout → sipariş (SEPA talimatı) → sorgu, konsol/5xx hatası yok, mobilde yatay taşma yok.
- [ ] **Henüz yok:** Stripe/PayPal sandbox ile gerçek ödeme denemesi, ürün/bölge/hikâye görselleri (Core medya yüklemesi), `balci` için gerçek logo/favicon, yasal sayfa içerikleri, e-posta bildirimleri, erişilebilirlik denetimi (axe) ve Lighthouse, GitHub deposu + CI.

## İlerleme (Müşteri hesapları, 6 Ekim 2026)

**Core**
- [x] Hesap rotaları `/me/storefront/*` altına taşındı: `PermissionMiddleware` `/me` önekini RBAC'tan muaf tutar; eski `/storefront/me/*` yolu müşteri token'ıyla 403 veriyordu. Yeni uçlar: `cart[...]`, `cart/claim`, `shipping-options`, `checkout`, `orders` (liste), `orders/:id[/payment[/confirm]]`, `account` (profil), `account/logout`.
- [x] Herkese açık hesap uçları (`/public/storefront/account/*`: `register`, `login`, `refresh`, `verify-email`, `forgot-password`, `reset-password`): `auth` oran sınırlayıcısı (5/dk) altında.
- [x] **Güvenlik:** genel `/auth/register` yanıtında e-posta doğrulama token'ını döndürüyordu (başkasının adresini "doğrulamak" mümkündü). Storefront kaydı token'ı asla döndürmez, yalnız e-postayla gider. **Genel `/auth/register` da düzeltildi (Core `2843264`):** `APP_ENV=production` sunucuya karşı canlı olarak doğrulandı (yetkisiz kayıt → token → doğrulama → giriş); artık token yalnız e-postayla gider, yerel geliştirmede SMTP kapalıyken (production dışı) hâlâ döner. Aynı kalıp `/auth/forgot-password` ve `send-verification`'da yok.
- [x] Müşteri = `users.SiteID == mevcut site`. Kayıt `SiteID` yazar; giriş/refresh/profil/şifre sıfırlama yalnız bu siteye ait kullanıcıya çalışır (yönetici ve başka site müşterisi giremez; yabancı refresh token'ı iptal edilir). Kullanıcı e-postası tenant genelinde tekildir (başka sitede kayıtlı e-posta `account_exists` verir).
- [x] `auth.loginWithUser`: parola kontrolü, "e-posta doğrulanmadı/onay bekliyor" durumundan ÖNCE yapılır (hesap durumu yalnız parolayı bilene açılır). Bilinmeyen e-posta için sahte bcrypt (zamanlama ayrımı yok).
- [x] E-posta bağlantıları site alan adına gider: `auth/service.WithLinkOverride` (ctx) → `https://<site.Domain>/compte/verification?token=…`, `/compte/mot-de-passe?token=…`. Parola politikası: 10–72 bayt, harf+rakam.
- [x] `orders.Filter` → `CreatedBy`, `SiteID` (sipariş geçmişi sahibi+site ile sınırlı).
- [x] Testler: servis (site kapsamı, token sızıntısı yok, hata eşleme, refresh iptali), handler (anahtarlar, `no-store`), auth bağlantı geçersiz kılma; gerçek Mongo+HTTP ile kayıt→doğrulama→giriş→profil→sipariş.

**Frontend (balci)**
- [x] Tarayıcı → `app/api/account/*` (BFF). Token'lar `httpOnly` çerezlerde (`balci_at` ≈30 dk, `balci_rt` 7 gün); tarayıcı yalnız profili görür. Refresh token Core'da her kullanımda döner → eşzamanlı istekler tek refresh'e birleştirilir (`RefreshCoalescer`, 15 sn). Core geçici hatası çerezleri silmez; reddedilen token siler.
- [x] `callWithSession` oturum açıkken `/public/storefront/*` çağrısını otomatik `/me/storefront/*`'a çevirir (sepet, kargo, checkout, sipariş, ödeme) → mevcut sayfalar değişmeden çalışır. Girişte misafir sepeti hesaba birleştirilir (`claim`).
- [x] Sayfalar: `/compte`, `/compte/connexion` (`?next=` açık-yönlendirme korumalı), `/compte/inscription` (gizlilik onayı), `/compte/verification`, `/compte/mot-de-passe`; header hesap simgesi; checkout giriş yapılmışsa ad/e-postayı doldurur. Etiketler Core Pages JSON'unda `labels.account.*` + `shop.errors/fields` (fr/en/tr).
- [x] **Üretim hatası düzeltildi:** BFF Core'a istemci IP'sini iletmiyordu → tüm alışverişçiler tek oran-sınırı kovasını (public 60/dk, auth 5/dk) paylaşırdı. `lib/api/client-ip.ts` `X-Forwarded-For` iletir. **Nginx `proxy_set_header X-Forwarded-For $remote_addr;` (ekleme değil, üzerine yazma) olmalı** — Core anahtar olarak ilk girdiyi kullanır.
- [x] Doğrulama: `tsc`, `next lint`, 12 birim testi, Playwright (masaüstü+mobil): kayıt (zayıf parola reddi) → doğrulanmamış giriş reddi → e-posta bağlantısı → misafir sepeti → giriş (sepet taşındı, `next` yönlendirmesi) → hazır doldurulmuş checkout → sipariş → `/compte` geçmişi → sipariş sayfası → çıkış (çerezler silindi).
- [x] **Framework güvenlik yükseltmesi (12 Ekim 2026):** Next.js `14.2.33` → `15.5.27`, React/React DOM `18.2` → `19.3`; tüm App Router `params`, `searchParams` ve `headers()` kullanımları asenkron sözleşmeye geçirildi. Next.js'in eski transitif PostCSS sürümü `8.5.29` ile override edildi. `npm audit --omit=dev`: **0 bulgu**; `tsc`, ESLint flat config, 12 birim testi, env'li/env'siz production build ve standalone rota smoke testi geçti. CI'a `audit:prod` kapısı, haftalık npm/GitHub Actions Dependabot kontrolü eklendi. Ayrıntılı devir kaydı: `done/TASK_063_BALCI_NEXT15_SECURITY_UPGRADE.md`.
- [ ] **Açık:** (1) Üretimde SMTP yapılandırması şart — yoksa doğrulama e-postası gitmez, kimse giriş yapamaz (`APP_ENV=production` ve SMTP kontrolü kurulum listesine eklenecek); (2) `labels.account.*` seed'i insert-only: mevcut ortamlarda ana sayfa JSON'u yeniden yüklenmeli/admin'den eklenmeli; (3) şifre sıfırlama mevcut oturumları iptal etmiyor (Core `ResetPassword` kullanıcı ID'si döndürmüyor); (4) hesap silme/veri dışa aktarma (RGPD erişim/silme hakkı) ve profil düzenleme/adres defteri yok; (5) `account_exists` yanıtı e-posta varlığını ifşa eder (auth oran sınırı ile sınırlı).

## İlerleme (Faz 3 — installer, 6 Ekim 2026)

- [x] `aystek-core-installer`: `docs/examples/balci-local-source.json` ve `balci-registry.json` (`public-site`, port **3003**, kendi dizin/image/state'i; haber siteleriyle çakışmaz — sözleşme testi). Env şablonu `SITE_URL`, `SITE_DOMAIN_ALIASES`, `API_BASE_URL`, `AYSTEK_STORE_DOMAIN`, `AYSTEK_TENANT_ID`'i zaten üretiyordu; frontend aynı değişkenleri okur.
- [x] Public-site Nginx vhost'u `X-Forwarded-For`'u `$remote_addr` ile **üzerine yazar** (önceden `$proxy_add_x_forwarded_for` ile ekliyordu → ziyaretçi kendi oran-sınırı anahtarını seçebilirdi; haber siteleri de bu şablonu kullanır, davranış değişikliği yalnız güvenlik yönünde). `templates/` ve gömülü `assets/templates/` ikisi de güncellendi.
- [x] Dokümantasyon: `KULLANIM_KILAVUZU.md` §22 ve `SUNUCU_KURULUM_REHBERI.md` §15.7 (seed `--public-site=balci`, SMTP, `APP_ENV`, Nginx, ödeme anahtarları, içerik listesi).
- [ ] **Açık:** (1) `ghcr.io/huseyinkoc-aystek/balci-aysteknolojileri-com` için GitHub deposu + CI (şu an yalnız yerel git, uzak depo yok); registry profili o zamana kadar kullanılamaz, `local-source` için kaynak sunucuya kopyalanmalı; (2) gerçek sunucuda DNS/TLS/Nginx ve `aystek-seed --public-site=balci` denemesi (bu sandbox'ta yapılamaz; yalnız şablon/profil birim testleri çalıştı); (3) ~~`tls_email` yer tutucusu~~ → `info@aysteknolojileri.com` kullanıldı.

## İlerleme (E-postalar, yasal sayfalar, RGPD, CI — 6 Ekim 2026)

- [x] **GitHub:** `huseyinkoc-aystek/balci-aysteknolojileri-com` (varsayılan dal `master`); `.github/workflows/ci.yml` (tsc, lint, test, build) ve `release.yml` (`v*.*.*` etiketi → `ghcr.io/huseyinkoc-aystek/balci-aysteknolojileri-com:<sürüm>`). Core, installer ve admin panel `main`'e push edildi.
- [x] **Marka/dil sorunu bulundu ve düzeltildi:** Core e-posta şablonları tenant geneli, "AYSTEK CORE" markalı ve Fransızca yoktu → balci müşterisi TR/EN AYSTEK e-postası alırdı. Çözüm: `platform/mail.Branding` (ctx) — mailer önce `storefront_<anahtar>` şablonunu dener (site adı `{{site_name}}`, varsayılan dil `fr`), yoksa eski genel şablona düşer; genel şablonlara dokunulmadı. Seed (insert-only, admin düzenlemesi ezilmez): `storefront_email_verification`, `storefront_password_reset`, `storefront_order_confirmation` × fr/en/tr.
- [x] **Sipariş onay e-postası:** sipariş kaydedilince (`storefront/service/order_mail.go`, arka plan goroutine'i, checkout'u asla bozmaz): numara, satırlar, toplam, havale ise IBAN/BIC/referans, takip bağlantısı. Online ödemede "ödemeyi tamamlamadıysanız takip sayfasından devam edin" notu.
- [x] **Mailer düzeltmeleri (tüm Core için):** Subject RFC 2047 ile kodlanır ve CR/LF temizlenir (başlık enjeksiyonu + ham 8-bit "é" bozulması); gövde `quoted-printable` (UTF-8, eskiden 7bit etiketli); HTML gövdesinde kullanıcı değerleri kaçışlanır (`*_html` değişkenleri güvenilir) — kayıtta `full_name` içine HTML/phishing yerleştirme engellendi; `/auth/register` karşılama adı artık ad-soyad.
- [x] **forgot-password oracle'ı kapatıldı:** mail gönderimi hata verirse gerçek hesapta 500, olmayanda 200 dönüyordu (hesap varlığı sızıyordu); artık her durumda aynı yanıt, hata sunucu günlüğüne.
- [x] **Şifre sıfırlama oturumları sonlandırır** (Core `ResetPassword`: access+refresh token'lar silinir).
- [x] **Yasal sayfalar:** `mentions-legales`, `cgv`, `confidentialite` Fransızca **taslak** (`draft`) olarak seed edilir (insert-only); `[À COMPLÉTER]` işaretli şirket/hosting/arabulucu alanları. Taslak olduğundan public API'de yayınlanmaz (footer bağlantıları 404 verir) — yönetici doldurup hukukçuya onaylatıp yayınlamalı. İçerik: LCEN künyesi, CGV (14 gün cayma, L.221-28 istisnaları, yasal garantiler, arabulucu), RGPD politikası (amaç/hukuki dayanak/süre, alıcılar, haklar, CNIL, çerezler).
- [x] **RGPD hesap hakları:** `GET /me/storefront/account/export` (profil + tüm siparişler, JSON indirme) ve `POST /me/storefront/account/delete` (parola ile yeniden doğrulama → oturumlar kapanır → kullanıcı anonimleştirilir + soft-delete; siparişler muhasebe için saklanır). `/compte` sayfasında "Mes données" bölümü.
- [x] **Core hatası düzeltildi (tüm siteleri etkiler):** `GET /public/pages/:lang/*slug` taslak (`draft`), zamanlanmış ve arşivlenmiş sayfaları da yayınlıyordu (slug bilen herkes okuyabilirdi). Artık yalnız `published` döner, diğerleri 404 (varlıkları da sızmaz). Yasal taslakların yayınlanmaması bu düzeltmeye bağlı; canlıda doğrulandı (draft 404 → yayınla → 200).
- [x] **Erişilebilirlik:** axe-core (WCAG 2.1 A/AA + best-practice) 15 sayfa × masaüstü/mobil + oturum açık durumlar (sepet, checkout, hesap, silme formu): **0 ihlal**. Bulunup düzeltilenler: kontrast (amber-600→700, footer zinc-500/600→400, emerald-600→700), başlık sırası, `/histoire` h1, ürün kartı bağlantı adı, sepet bağlantısı adı, `dl` içinde `p`, CMS'ten gelen soluk bölge rengiyle yazı. Lighthouse (performans/SEO skoru) bu sandbox'ta çalıştırılmadı.
- [ ] **Açık:** (1) E-posta SMTP'si gerçek ortamda denenmedi (sandbox'ta yalnız yerel SMTP alıcısı); (2) ~~"ödeme alındı"/"kargoya verildi" e-postaları~~ → yapıldı: `orders.OrderStatusHook` (`storefront/service/order_status_mail.go`) ödendi/gönderildi/iptal için fr/en/tr e-posta yollar; yalnız vitrin siparişleri, siparişin kendi sitesiyle, durumu kim değiştirirse (ağ geçidi geri çağrısı, yönetici) tetiklenir; (3) ~~profil düzenleme~~ ve ~~oturum açıkken şifre değiştirme~~ yapıldı (`PATCH /me/storefront/account`, `POST …/account/password`; şifre değişince tüm oturumlar kapanır); **adres defteri** yok (ayrı site-scoped model + seed gerekir); (4) yasal sayfa içeriklerinin hukuki onayı; (5) ~~Lighthouse~~ ölçüldü (production build, mobil profil, sandbox ağı): ana sayfa/katalog performans **95–96**, erişilebilirlik **100**, en iyi uygulamalar **100**, SEO **100** (hero görseline `fetchpriority=high` ve site ikonu eklenince; öncesinde ana sayfa performans 77, LCP 6,4 sn). Gerçek sunucuda ve kendi görsellerle yeniden ölçülmeli (görseller şu an Unsplash'tan); (6) Stripe/PayPal sandbox; (7) **gerçek** ürün/üretim fotoğrafları — şimdilik şablonun Unsplash fotoğrafları kullanılıyor (hero, hikâye, 4 bölge, 6 demo ürün; seed'de `balciDemoImages` + ana sayfa JSON'ları). Bunlar hot-link edilen yer tutucudur: lisans/telif ve RGPD (ziyaretçi IP'si Unsplash'a gider) nedeniyle yayından önce kendi fotoğraflarla (Admin → Medya, ürün medyası, ana sayfa içeriği) değiştirilmeli; şablonun bölge fotoğrafları da alakasız (Ardèche=ekmek, Haute-Loire=havuç). Şablondaki Isère fotoğrafı Unsplash'ta silinmişti, yenisiyle değiştirildi.

## İlerleme (SMTP / e-posta hesapları, 6 Ekim 2026)

- [x] **Bulgu:** Yönetim panelindeki *Ayarlar → E-posta hesapları* ekranı (28 dil çevirisiyle) `/settings/email-accounts` ucunu çağırıyordu, ama Core'da bu uç **hiç yoktu** (eski admin-panel kalıntısı); mailer yalnız `smtp_*` ayarları ve `SMTP_*` env'ini okuyordu → panelden SMTP girmek mümkün değildi.
- [x] **Çözüm (Core `emailaccounts` modülü):** CRUD, amaç başına varsayılan hesap, **Test** düğmesi (hesabın kendi adresine gönderir, SMTP'nin gerçek hata nedenini döner), şifre `PAYMENT_GATEWAY_ENCRYPTION_KEY` ile şifreli saklanır ve API'den asla dönmez (düzenlemede boş = koru). Mailer önceliği: **varsayılan e-posta hesabı → env (yedek)**; hiçbir ekrandan düzenlenemeyen eski `smtp_*` ayar yolu README standardı gereği kaldırıldı, etkin eski yapılandırma seed ile varsayılan hesaba taşınır (`SeedEmailAccounts`, insert-only). Hata anahtarları panelde zaten seed'li `error.email_account_*` anahtarlarıdır; test e-postası V1 sözleşmesindeki gibi giriş yapan yöneticiye gider. Panelde yalnız kontrol ekranı sayacı düzeltildi (eski `app_settings.email_accounts` alanı kaldırıldı). Gerçek Mongo + sahte SMTP sunucusuyla zincir doğrulandı.
- [x] **Installer düzeltmesi:** `core.env` yalnız `ENCRYPTION_KEY` yazıyordu, Core `PAYMENT_GATEWAY_ENCRYPTION_KEY` okur → ödeme ağ geçidi sırları da düz metin saklanıyordu. Şablona eklendi (mevcut sunucularda elle eklenmeli).
- Rehber: `user-guide/05-operations/EPOSTA_SMTP_KURULUMU_TR.md` (Gmail/Workspace uygulama şifresi, DKIM/SPF/DMARC, panel alanları).

## İlerleme (Blog, 7 Ekim 2026)

**Core altyapısı (teyit edildi — yeni modül yazılmadı):** `posts` modülü zaten site-scoped: yönetici CRUD (`/admin/posts`, panelde *İçerik → Yazılar*, balci görünüm profilinde `posts` menüsü açık), herkese açık liste/detay (`/public/posts`, `/public/posts/:lang/:slug`; yalnız `published`, `published_at <= şimdi`, tenant+site süzgeçli; zamanlanmış/taslak/başka site 404). Canlı doğrulandı: yönetici API'siyle taslak, zamanlanmış ve yayınlanmış yazı oluşturuldu → yalnız yayınlanan görünüyor.
- [x] **Core hatası düzeltildi:** `status=published` ile ama `published_at` verilmeden kaydedilen yazı hiçbir herkese açık sitede görünmüyordu (panel formundaki yayın tarihi alanı isteğe bağlı ve varsayılan boş). `Create` artık tarihi "şimdi" yapar; `Update` ilk yayında "şimdi", sonraki düzenlemelerde özgün tarihi korur (önceden düzenleme tarihi boşaltabilirdi). Tüm siteleri (aystekhaber, kureselhat…) etkileyen bir düzeltmedir.
- [x] **Site (`/blog`, `/blog/[slug]`):** liste (başlık+giriş, kartlar, sayfalama `?page=`, aralık dışı sayfa 404, 2. sayfa ve sonrası `noindex`), yazı sayfası (tarih, okuma süresi, kategori rozetleri, kapak, içerik stili, mağazaya yönlendiren kutu, "À lire aussi"), JSON-LD `BlogPosting` + `BreadcrumbList`, RSS (`/blog/rss.xml`), site haritasına tüm yayınlı yazılar (site dillerine göre), fr/en/tr etiketler (`content.blog.*`). Erişilebilirlik (axe) liste/yazı × fr/en/tr × masaüstü/mobil: 0 ihlal.
- [x] **Örnek içerik:** `--demo-catalog` ile 3 `[DEMO]` bilgilendirici makale (gezici arıcılık, iyi bal nasıl tanınır, saklama ve tadım) fr/en/tr; yalnız bu siteye ait, insert-only. Yayından önce editör silmeli/yeniden yazmalı. Görseller şablonun Unsplash fotoğraflarıdır (yer tutucu).
- Not: kategori/etiket taksonomisi tenant geneli (site izolasyonu yok) olduğundan demo yazılarda kategori atanmadı; editör panelden seçerse yazı sayfasında rozet olarak görünür. Mevcut ortamlarda ana sayfa JSON'u insert-only olduğundan yeni `blog.*` etiketleri admin'den eklenmeli (geliştirme DB'sine tek seferlik uygulandı).
- [ ] **Açık:** yazı kategorisine göre süzme/arama (kategoriler tenant geneli olduğundan site-scoped menüden türetilmeli), blog yazısı yorumları (`comments` modülü hazır, vitrinde açılmadı), yazı içi görsel yükleme akışının panel üzerinden denenmesi.

## İlerleme (Müşteri yorumları, 7 Ekim 2026)

**Teyit:** Core'daki `comments` modülü ürün yorumu için uygun değil (puan yok, satın alma doğrulaması yok, auth'suz public uç, ürün/sipariş bağı yok). README kuralı gereği ayrı ve eksiksiz bir **`reviews`** modülü yazıldı (paralel uygulama yok; `comments` blog/haber yorumları için kalır).
- [x] **Core `reviews` modülü:** model (`product_reviews`: puan 1-5, başlık, içerik 10–2000 karakter, düz metin — `<`/`>` reddedilir, `pending|approved|rejected`, `featured`, `verified_purchase`, tenant+site), repository (arayüz+Mongo, `order_id` üzerinde kısmi benzersiz indeks, sert silme), servis (`PurchaseVerifier` + `ProductDirectory` portları; yalnız **kendi ödenmiş siparişindeki** ürün için, sipariş+ürün başına bir yorum, her yorum önce `pending`), handler/rotalar: yönetici `GET/PUT/DELETE /admin/reviews…`, public `GET /public/storefront/reviews` (yalnız onaylı; özet: adet/ortalama/dağılım; e-posta ve kişi/sipariş kimliği dönmez, yazar "Claire M."), müşteri `POST /me/storefront/reviews`, `GET /me/storefront/orders/:id/reviewable`. İzin, rol, menü (*Ticaret → Siparişler → Yorumlar*), görünüm profili (balci), 27 anahtar × 28 dil çeviri (tamlık testine bağlı) seed'lendi. Tüm katmanlarda test (gerçek Mongo dahil).
- [x] **Yönetim paneli (`/admin/reviews`):** durum sekmeleri (bekleyen sayacıyla), puan/yorum/müşteri/rozet tablosu, onayla / reddet / ana sayfada öne çıkar / sil (onay penceresiyle). Yalnız onaylı yorum öne çıkarılabilir (`reviews_not_approved`).
- [x] **Site:** ana sayfada "Avis clients" bloğu (yalnız **öne çıkarılmış + onaylı + ≥4 yıldız** yorumlar, ortalama/adet özeti, ürün bağlantısı, "Achat vérifié" rozeti yalnız Core doğruladıysa; yorum yoksa blok hiç görünmez), ürün sayfasında özet + onaylı yorumlar + JSON-LD `AggregateRating`/`Review` (yalnız gerçek onaylı yorumlar), sipariş sayfasında oturum açmış ve **ödenmiş** siparişin sahibine ürün başına yorum formu (BFF: `POST /api/account/reviews`, `GET /api/account/orders/:id/reviewable`; yalnız aynı origin). Etiketler `labels.reviews.*` (fr/en/tr).
- [x] **Canlı doğrulama (geçici Core+Mongo, gerçek DB'ye dokunmadan):** kayıt → doğrulama → giriş → sipariş: ödenmeden form yok → ödendi → form → istemci doğrulaması (puan/çok kısa) → gönderim `pending`+`verified_purchase`, ürün sayfasında görünmez → ikinci yorum 409 `reviews_already_reviewed`, siparişte olmayan ürün 403 `reviews_not_purchased`, HTML 400 → onay+öne çıkarma → ana sayfada kart ("Claire M.", doğrulanmış rozeti) ve ürün sayfasında görünür. axe: ana sayfa (masaüstü+mobil) ve ürün sayfası 0 ihlal.
- [x] **Demo içerik:** `--demo-catalog` ile 7 `[DEMO]` yorum (fr/en/tr, hiçbiri "doğrulanmış" değil — sahte güven rozeti verilmez; 5'i öne çıkarılmış). Yayından önce silinmeli; gerçek müşteri yorumları geldikçe panelden onaylanıp öne çıkarılır.
- Not: mevcut ortamlarda ana sayfa JSON'u insert-only olduğundan `labels.reviews.*` admin'den eklenmeli (geliştirme DB'sine tek seferlik uygulandı). Public yorum listesi vitrin önbelleğinde 10 sn tutulur.
- [ ] **Açık:** yorumlara yönetici/işletme yanıtı, yorum sonrası teşekkür e-postası, bekleyen yorum için yönetici bildirimi, dile göre yorum süzme (şimdi tüm diller birlikte listelenir), gerçek müşteri yorumlarıyla moderasyon denemesi.

## İlerleme (balci yönetici hesabı + rol + görünüm profili, 7 Ekim 2026)

- [x] **Rol `balci_manager`** (Core `public_site_manager_seed.go`, `publicSiteManagerScopes`): yalnız balci'nin kendi ekranlarının yetki modülleri — içerik (pages, posts, media, comments, localized_content, contacts, sliders, newsletter_subscribers), katalog/vitrin (products, categories, tags, storefront, reviews), sipariş/kargo/ödeme/fiyat (orders, returns, shipping, payments, billing, taxes, discounts, coupons, customers, currencies), dil/e-posta şablonu/çalışma saatleri (languages, mailtemplates, business_hours); `settings` ve `stores` yalnız read/create/update (silme yok); `dashboard` ve `notifications` yalnız read. Stok, filo, PMS, kullanıcı/rol/izin, entegrasyon ve platform (system, system_ops) modülleri **verilmez**. `PermissionMiddleware` yolun ilk segmentinden modül çıkarır; `TestBalciManagerCoversVisibleScreens` görünür her menünün modülünün rolde olduğunu doğrular.
- [x] **Menü rolleri:** menülerde `Roles` alanı rol adıyla süzer (siparişler, ödemeler, fiyatlama vb. yalnız `admin` idi → yeni rol göremezdi). `menuSeedRoles` rol adını, balci görünüm profilinde görünür olan her menünün `Roles` listesine ekler (`SeedMenus` her çalışmada günceller); görünür olmayan menülere eklenmez (test).
- [x] **Görünüm profili `balci_shop_manager`:** izinli liste genişletildi (`mail-templates`, `settings-dashboard`, `settings-general`); geri kalan tüm menüler gizli. `SeedPublicSite` mevcut profili korur (admin düzenlemeleri kaybolmasın), bu yüzden yönetici komutu profili **zorla yeniden hesaplar** — profil o güne ait menü kümesinden eskimesin.
- [x] **Hesap:** Gül Ali Delidağ, `delidag48@gmail.com`, kullanıcı adı `delidag48` — rol `balci_manager`, `ContextMemberships.SiteIDs=[balci]` (başka siteye geçemez), e-posta doğrulanmış + onaylı. Komut: `go run ./cmd/seed --public-site=balci --create-manager [--manager-password …]`; parola kaynak kodda **yoktur**: verilmezse (`--manager-password` veya `SEED_MANAGER_PASSWORD`) rastgele 20 karakter üretilir ve yalnız terminale bir kez yazılır (günlüğe girmez). İdempotent: var olan hesabın şifresi değişmez, rol/kapsam/onay onarılır. Ana seed (`SeedRoles`, `SeedMenus`) rolü ve menü rollerini her kurulumda eşitler; hesap oluşturma bilerek ayrı komuttur (kimlik bilgisi gerektirir).
- Test: birim (kapsam/menü rolleri/parola) + gerçek Mongo entegrasyonu (izin→rol→menü→hesap, yeniden çalıştırma, bozulmuş hesabı onarma, açık parola).
- [ ] **Açık:** `settings-general` tenant geneli ayar sayfasını da açar (yalnız `settings` read/create/update) — kullanıcı yalnız site ayarı isterse bu menü anahtarı çıkarılıp ayrı bir site-ayar ekranı gerekir; yönetici girişiyle ekranların gerçek tarayıcıda gezilmesi.

## Açık kararlar / kullanıcıdan gerekenler

1. ~~Misafir ödeme~~ → **karar verildi: gerekli** (Faz 1 backend işi).
2. ~~Ödeme yöntemleri~~ → **karar verildi:** havale, kredi kartı, PayPal. Kart sağlayıcısı olarak mevcut Stripe adaptörü önerilir (Mollie istenirse yeni adaptör yazılır) — **onay gerekli**.
3. **Hesaplar:** Stripe ve PayPal API anahtarları, havale için IBAN/BIC/banka adı (admin panelden girilecek, koda gömülmez).
2. **Gerçek içerik:** adres/telefon/e-posta/çalışma saatleri, logo, gerçek ürün ve üretim fotoğrafları (Unsplash yer tutucuları yayına çıkmayacak).
3. **Dil kapsamı:** fr/en/tr yeterli mi, yoksa tenant'ın 28 dilinden daha fazlası mı açılsın?
4. **Hukuki:** Fransa için mentions légales, RGPD/çerez metinleri (kullanıcı/hukuk onayı).

## Temizlik kuralı (README "Remove Obsolete Work")

Port tamamlanınca Figma çıktısının hiçbir sabit içerik/çeviri/mock dosyası yeni projede bulunmamalı; ölü bileşen (ör. kullanılmayan shadcn `ui/*`, MUI) alınmaz.
