# AYSTEK Haber Public Site

## Amaç ve Sınır

AYSTEK Haber, `aystek` tenant'ı altındaki `aystekhaber` site/store bağlamında
çalışan haber portalıdır. Public arayüz `/workspace/node/aystekhaber-com`, içerik
ve yönetim yetkisi AYSTEK Core, kurulum yaşam döngüsü ise
`aystek-core-installer` sorumluluğundadır. Site kendi MongoDB'sini veya paralel
bir CMS modelini oluşturmaz.

## Veri Sahipliği

| Görsel alan | Core sahibi | Yönetim noktası |
|---|---|---|
| Haber ve haber detayı | posts | İçerik > Gönderiler |
| Kategori başlıkları | taxonomy | Ürünler > Kategoriler |
| Ana menü | menus | Ayarlar > Menü |
| Manşet/son dakika/sidebar/çok okunan | post placements | Gönderi aksiyonları |
| Ana sayfa etiketleri, linkler ve seçimler | pages (`anasayfa`) | İçerik > Sayfalar |
| Yazar adı/unvan/biyografi/fotoğraf | authors | İçerik > Yazarlar |
| Ana sayfadaki yazar sırası | page JSON `authorSlugs` | İçerik > Sayfalar |
| Maçlar ve puan durumu | sports snapshot | İçerik > Spor Merkezi |
| E-bülten aboneleri | newsletter | Pazarlama > Bülten Aboneleri |
| Reklam yönetimi | ads | Pazarlama > Reklamlar |

Yazar bilgisi sayfa JSON'unda tekrar tutulmaz. `authorSlugs` yalnız sunumda
hangi yazarların hangi sırada görüneceğini belirler; kart içeriği her istekte
Core author API'sinden alınır.

## Public API Sözleşmesi

Tüm istekler aşağıdaki çözümleme bağlamını taşımalıdır:

```http
X-Tenant-ID: aystek
X-Store-Domain: aystekhaber.com
```

Başlıca uçlar:

- `GET /public/posts?lang=tr&page=1&limit=100`
- `GET /public/posts/{lang}/{slug}`
- `GET /public/post-placements/{position}`
- `GET /public/navigation/main?lang=tr`
- `GET /public/pages/tr/anasayfa`
- `GET /public/authors?lang=tr&page=1&limit=100`
- `GET /public/authors/{slug}?lang=tr`
- `GET /public/sports/leagues/super-lig`
- `POST /public/newsletter/subscribe`

Public authors uçları yalnız `published` ve aktif profilleri, public sports ucu
yalnız aktif snapshot'ı döndürür. Admin uçları pasif kaydı yeniden etkinleştirme
amacıyla okuyabilir.

## Site Rotaları

- `/`: dinamik ana sayfa
- `/haber/{slug}?lang={kod}`: haber detayı
- `/yazar/{slug}?lang={kod}`: yazar profili
- `/rss.xml`: son 50 Türkçe haber
- `/sitemap.xml`: sitenin sunduğu dillerdeki haber ve yazar rotaları
- `/robots.txt`: crawler sözleşmesi
- `/manifest.webmanifest`: PWA/site metadata sözleşmesi
- `/api/newsletter`: browser ile Core arasındaki server-side proxy

## Sports Domain Kuralları

`LeagueSnapshot` tenant ve site bağlamında benzersiz `code` ile saklanır.
Takım sıraları ve takım adları tekrar edemez; puan tablosu sayıları negatif
olamaz; galibiyet, beraberlik ve mağlubiyet toplamı oynanan maç sayısını
aşamaz. Maç statüsü `scheduled`, `live` veya `finished` olmalıdır. Kısa takım
adları ve form değerleri normalize edilir.

Sports katmanları:

1. Model normalize eder ve domain invariant'larını doğrular.
2. Service public/admin okuma politikasını ayırır ve sıralamayı uygular.
3. Repository tenant + site filtresini zorunlu kılar ve idempotent upsert yapar.
4. Handler public envelope veya admin sözleşmesini HTTP'ye çevirir.
5. `/admin/sports` ekranı Türkçe/İngilizce lig adını, tabloyu, maçları ve aktif
   durumunu yönetir.

## Seed ve İdempotency

AYSTEK Haber seed zinciri store, kategoriler, ana sayfa, menü, yazarlar, 70
başlangıç haberi, placement kayıtları, branding, sports snapshot ve görünüm
profilini uzlaştırır. Upsert filtreleri sabit domain kimliklerine dayanır;
tekrar çalıştırma yeni kopya üretmemelidir. Üretimde `SEED_FULL_RESET=false`
kullanılır.

Seed sonrası editoryal görünüm profilinde yalnız gerekli alanlar görünür:
posts, pages, authors, sports, taxonomy, media, marketing/ads/newsletter, menu,
languages/site-languages/translations ve piyasa verileri. Sipariş, ödeme, telematics
ve sistem operasyonları haber editörlerinden gizlidir.

## Diller

Sitenin sunduğu diller kodda sabit değildir; Core'daki site dil seçiminden gelir
(`/admin/languages/site`, seçim yoksa tenant geneli aktif diller). Dil değiştirici,
sitemap, bülten kaydı ve `?lang=` çözümlemesi bu listeyi kullanır; sitede
kapalı bir dil istenirse site varsayılan diline düşülür.

Portal etiketleri ve sayfa metinleri Pages modülünde ana sayfa kaydının dil
yerelleştirmesinden yönetilir. Bir dil için yerelleştirme yoksa önce İngilizce,
sonra Türkçe kullanılır. Yeni bir dili tam çevrilmiş göstermek için ana sayfa
kaydına o dilin yerelleştirmesini ekleyin.

## Yayın ve Doğrulama

```bash
cd /workspace/node/aystekhaber-com
npm run check-env
npm run lint
npm run build

cd /workspace/go/aystek-core
go test ./internal/modules/authors/... ./internal/modules/sports/... ./internal/seeds

cd /workspace/go/aystek-core-installer
go test ./internal/config ./internal/tls ./internal/installer ./internal/deploy
```

Production kurulumu için `user-guide/05-operations/SUNUCU_KURULUM_REHBERI.md`
Bölüm 15 izlenmelidir. Public-site Let's Encrypt config'inde `tls_email`
zorunludur; `public_site_domain_aliases` değerleri DNS, Nginx ve sertifika SAN
listesinde birlikte yönetilir.

## Sorun Giderme

- Ana sayfa boşsa önce Core seed logunu ve store/tenant header'larını kontrol edin.
- Yazarlar görünmüyorsa profile `published + is_active=true` olduğundan ve
  `authorSlugs` seçkisinde slug bulunduğundan emin olun.
- Sports alanı görünmüyorsa snapshot aktifliğini ve site bağlamını doğrulayın.
- Sadece `www` çalışmıyorsa DNS alias'ı ile `public_site_domain_aliases`
  sözleşmesini karşılaştırın ve sertifika SAN kapsamını inceleyin.
- E-bülten hatasında Next.js proxy logunu ve Core newsletter rate-limit/yanıtını
  birlikte inceleyin.
