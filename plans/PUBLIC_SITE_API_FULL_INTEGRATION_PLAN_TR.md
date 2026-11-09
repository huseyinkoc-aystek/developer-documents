# PUBLIC SITE API FULL INTEGRATION PLAN (TR)

**Tarih:** 2026-04-10  
**Kapsam:** `/workspace/node/aysteknoloji-com` → `/workspace/go/aystek-core` tam public CMS/API entegrasyonu  
**Hedef:** `aystek-core-next` seviyesinde tenant/site bazlı, çok dilli, SEO uyumlu, DB-merkezli public site işletimi

---

## 1) Hedef Tanımı

Public site aşağıdaki içerikleri **doğrudan `aystek-core` public API’lerinden** almalı ve mümkün olan her yerde **SSR + SEO metadata** ile sunmalıdır:

- `posts`
- `pages`
- `sliders`
- `media`
- `translations`
- `languages`
- `menus`

Tüm akışlar şu kurallara uymalıdır:

1. **Tenant/Site Scope Doğruluğu**
   - `X-Tenant-ID`
   - `X-Store-Domain`
2. **Çok Dilli İçerik**
   - aktif diller backend’den okunmalı
   - public içerik locale fallback mantığıyla çözülmeli
3. **SEO Uyumluluğu**
   - `seo_title`, `seo_description`, `seo_keywords`
   - `og_title`, `og_description`, `og_image`
   - `canonical_url`
4. **Static içerik azaltımı**
   - fallback içerikler sadece API kesintisi durumunda devreye girmeli

---

## 2) Mevcut Durum Özeti

### Tamamlananlar
- Public nav/footer menüleri backend CMS’den okunuyor.
- Aktif dil listesi backend’den geliyor.
- İspanyolca (`es`) aktif hale getirildi ve public menülerde doğrulandı.
- `posts`, `pages`, `sliders` için public API client ve Next proxy route katmanı mevcut.
- Tenant/site header enjeksiyonu `lib/api/server.ts` üzerinden çalışıyor.

### Açık Gaps
- Bazı public sayfalar içerikleri client-side yüklediği için **SSR/SEO kazanımı eksik**.
- Blog detayında halen **static fallback veri bağımlılığı** var.
- Media/OG/canonical alanları tüm public rotalarda tutarlı şekilde kullanılmıyor.
- Dil seçimi server tarafına net taşınmadığı için **ilk render locale doğruluğu** eksik kalabiliyor.

---

## 3) Fazlı Uygulama Planı

## Faz 1 — SSR + SEO Temeli (BAŞLADI)
**Amaç:** CMS page/post verisini server-side okuyup initial render ve metadata üretmek.

### İşler
- [x] Mevcut entegrasyon yüzeyinin audit edilmesi
- [x] Bu plan dokümanının oluşturulması
- [x] Public site için request-language resolver (`cookie` + `Accept-Language`) eklendi ve cookie fallback zinciri sertleştirildi
- [x] Ana sayfa için CMS verisinin server preload edilmesi tamamlandı
- [x] Ana sayfa metadata’sının DB kaynaklı SEO alanlarıyla üretilmesi tamamlandı
- [x] Blog detay akışı dil parametresi + locale-aware slug/canonical yönlendirmesi ile backend’den okunuyor

### Çıktı
- Ana sayfa ilk render’da CMS içeriğiyle açılır
- SEO etiketleri DB’den beslenir
- Blog detayları `tr/en/es/...` için daha doğru locale fallback ile çözülür

## Faz 2 — Public Pages Parity
**Amaç:** `hakkimizda`, `cozumler`, `aystek-core`, `iletisim` sayfalarında SSR metadata ve CMS içerik bütünlüğü.

### İşler
- [ ] Sayfa bazlı `generateMetadata` entegrasyonu
- [ ] Static başlık/açıklama bloklarının CMS fallback düzenine taşınması
- [ ] Nested page slug senaryolarının tenant/site bazlı doğrulanması

## Faz 3 — Posts / Blog Tam API Geçişi
**Amaç:** Blog listesi ve detayın static demo verilerden bağımsızlaştırılması.

### İşler
- [x] İlgili yazılar listesi API tabanlı hale getirildi
- [x] Kategori/tag filtreleri dinamik post verisiyle besleniyor
- [x] JSON-LD ve article SEO alanları doğrudan post record’larından üretiliyor

## Faz 4 — Sliders / Media / Asset Pipeline
**Amaç:** Hero ve görsel alanlarda DB-medya tabanlı çalışma.

### İşler
- [ ] Slider görselleri ve CTA alanlarını CMS modelinden normalize etme
- [ ] `featured_image`, `cover_image`, `og_image` kullanımını tek pattern’e indirme
- [ ] Public media fallback kuralını standardize etme

## Faz 5 — Doğrulama ve Go-Live Checklist
**Amaç:** tenant/site/language/SEO doğrulaması ve operasyonel netlik.

### İşler
- [x] `npm run build` ile public site doğrulaması alındı
- [x] Go API build doğrulaması alındı
- [x] `/public/posts`, `/public/pages`, `/public/sliders`, `/public/navmenus` smoke testleri alındı
- [x] `tr/en/es` için içerik + metadata / locale çözümleme kontrolü yapıldı
- [x] sitemap / robots / canonical hardening doğrulaması yapıldı

---

## 4) Bu Oturumda Başlatılan Uygulama Dalgası

Bu planın hemen ardından aşağıdaki ilk uygulama dilimi yürütülecektir:

1. **Dil bilgisini server tarafına taşıma**
2. **Ana sayfayı CMS-first SSR + metadata uyumlu hale getirme**
3. **Blog detail API akışını locale-aware hale getirme**

---

## 5) Başarı Kriterleri

Aşağıdaki koşullar sağlandığında public entegrasyon bu faz için başarılı kabul edilir:

- Ana sayfa HTML çıktısı CMS kaynaklı başlık/açıklama ile gelir
- `<title>` / meta description / OG alanları DB kayıtlarıyla eşleşir
- `/api/posts/[slug]` tek bir locale’e bağlı kalmaz; seçili dile göre çözüm yapar
- `tr`, `en`, `es` için public içerikler fallback copy yerine doğru locale verisini gösterir
- Build çıktısında yeni TypeScript hatası oluşmaz

---

## 6) 2026-04-10 Checkpoint Güncellemesi

### Bu Oturumda Kapanan / İlerletilen İşler
- `aystek-core` tarafında `TenantCount()` guard düzeltildi; `tenants` koleksiyonu boş olsa bile `stores` fallback sayımı artık çalışıyor.
- `default` tenant için seed zincirine otomatik `admin-panel` sitesi eklendi:
  - `code: admin-panel`
  - `domain: core.aysteknolojileri.com`
- soft-delete olmuş store kayıtlarının seed ile tekrar ayağa kalkabilmesi için restore akışı eklendi.
- `aystek` public CMS seed’lerinde `site_id` çözümleme ve bağlama akışı netleştirildi.
- `hakkimizda`, `cozumler`, `aystek-core`, `iletisim` sayfaları CMS-first server wrapper modeline taşındı.
- blog liste/detay ve SEO metadata zincirinde statik `head.tsx` bağımlılığı temizlendi; OG title / description ayrımı güçlendirildi.

### Taze Runtime Doğrulaması
Aşağıdaki kanıtlar bu checkpoint içinde doğrulandı:

1. Backend health:
   ```bash
   curl -sf http://localhost:9090/healthz
   ```
   Dönüş:
   ```json
   {"success":true,"data":{"status":"ok"}}
   ```

2. Auth login:
   ```bash
   curl -sf -X POST http://localhost:9090/auth/login-by-username \
     -H "Content-Type: application/json" \
     -H "X-Tenant-ID: default" \
     -d '{"identifier":"huseyinkoc","password":"L5834!"}'
   ```
   Sonuç: access token + permissions listesi dönüyor.

3. Default tenant store visibility:
   ```bash
   curl -sf -H "X-Tenant-ID: default" -H "Authorization: Bearer <token>" \
     http://localhost:9090/admin/stores
   ```
   Sonuç: `admin-panel` ve `default` store kayıtları listelendi.

### Bu Plan İçin Net Kalan İşler
1. **Ana sayfa CMS-first preload kapanışı**
   - home hero / slider / featured blocks için server preload
   - home metadata’sını DB SEO alanlarından üretme

2. **Locale resolver sertleştirmesi**
   - `cookie + Accept-Language` ile ilk render dili
   - `tr/en/es` başlangıç render doğrulaması

3. **Media / canonical / sitemap son sertleştirme**
   - `cover_image`, `featured_image`, `og_image` standardizasyonu
   - canonical / robots / sitemap smoke doğrulaması

4. **Public API smoke matrisi**
   - `/public/posts`
   - `/public/pages`
   - `/public/sliders`
   - `/public/navmenus`
   tenant/site/language kombinasyonlarında test edilmelidir.

### Operasyonel Not
Bu workstream’de backend route prefix’leri şu şekildedir:
- `/auth/...`
- `/admin/...`
- `/public/...`

Bu nedenle public/admin frontend tarafında eski `/api/v1/...` varsayımı varsa alignment ayrıca kontrol edilmelidir.

## 7) Not

Bu plan **canlı bir uygulama planıdır**. Faz 1 ve Faz 2’nin önemli bölümü kod seviyesinde ilerletilmiş, runtime bootstrap doğrulaması alınmış ve kalan işler netleştirilmiştir.
