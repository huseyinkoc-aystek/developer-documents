# Ortak haber/magazin public-site: Küresel Hat geçiş planı

Son güncelleme: 27 Eylül 2026.
Durum: kod, seed, installer ve yerel production-browser doğrulaması tamamlandı;
editoryal içerik, marka varlıkları ve dış DNS/TLS yayını bekliyor.
Installer'ın GitHub/sunucu-içi Git ilk klonlama ve opsiyonel yerel AI
güncelleme kodu eklendi; gerçek uzak sunucuda Docker/Ollama kabul testi henüz yok.

## Kullanıcı talebi ve sınırlar

İlk yayına alınacak domain `kureselhat.com` olacak.
AYSTEK Haber entegrasyonu silinmeyecek veya Küresel Hat olarak yeniden adlandırılmayacak.
`/workspace/node/aystekhaber-com` aynı tasarımı kullanan haber/magazin siteleri
için ortak local geliştirme kaynağı kalacak. Sunucuda her domain aynı repository
ve sürümden gelen ayrı source checkout, install/config/log dizini, env, port,
container ve Nginx/TLS vhost ile çalışacak. Farklı tasarım localde de ayrı projedir.
Domain başına ayrı ürün kodu, yazar listesi veya içerik sistemi oluşturulmayacak.
Marka, içerik, dil, menü ve yönetilebilir ayarlar Core'un mevcut modellerinde tutulacak.
Tenant: AYSTEK Mühendislik (`aystek`).
Yeni site: Küresel Hat (`kureselhat`).
Ana domain: `kureselhat.com`.
Alias: `www.kureselhat.com`.
Çeviri borcu çalışmaları bu görevde değiştirilmedi.

## Aşama 1 — seed ve geliştirme veritabanı

Tamamlandı:

- Mevcut `stores` modelinde yeni site kaydı oluşturuldu.
- TR/EN lokalizasyonlu, siteye bağlı ana sayfa eklendi.
- Ana sayfa JSON'unda `brand.name`, başlık ve copyright site kimliğinden üretildi.
- 9 public navigasyon kaydı eklendi; mevcut 28 dilli menü başlıkları kullanıldı.
- RSS bağlantısını içeren sosyal ayar ve boş ticker olmak üzere 2 site ayarı eklendi.
- `kureselhat_editorial` görünüm profili ortak haber allowlist'i ile oluşturuldu.
- 8 haber kategorisi tenant düzeyinde tekrar kullanılabilir hale getirildi.
- Mevcut kategoriler, başlıkları ve landing-template seçimleri değiştirilmedi.
- Dil kataloğu tenant'tan miras alınıyor; ek bir sabit dil listesi oluşturulmadı.
- Yeni siteye örnek haber, yazar, galeri, placement veya maç verisi kopyalanmadı.
- Yeni frontend projesi, veri modeli veya public API oluşturulmadı.
- Core'un standart seed zincirine `SeedNewsSites` eklendi.
- Hedefli çalıştırma için mevcut `cmd/seed` komutuna `--news-site` seçeneği eklendi.

### Kod yerleşimi

Core kökü: `/workspace/go/aystek-core`.

| Dosya | Sorumluluk |
| --- | --- |
| `internal/seeds/news_sites_seed.go` | Site başlangıç kimliği, güvenli hedefli seed ve insert-only varsayılanlar |
| `internal/seeds/aystek_haber_seed.go` | Yeniden kullanılan ana sayfa, navigasyon, kategori ve public ayar şablonları |
| `internal/seeds/view_profiles_seed.go` | Ortak haber editoryal profil üretimi |
| `internal/seeds/runner.go` | Standart seed zincirine bağlantı |
| `cmd/seed/main.go` | Reset ve ilgisiz seed'leri çalıştırmayan hedefli komut |
| `internal/seeds/news_sites_seed_test.go` | Kimlik, izolasyon, tekrar çalıştırma ve çakışma testleri |

Site tanımları yalnız seed başlangıç verisidir; runtime domain yönlendirmesi için
ayrı bir sabit domain tablosu veya frontend switch-case mekanizması değildir.
Runtime domain/site/tenant çözümlemesi mevcut Core store-context mekanizmasına aittir.

### Çalıştırma

Core kaynak dizininde, hedef ortama ait Mongo yapılandırmasıyla:

```bash
go run ./cmd/seed --news-site=kureselhat
```

Ön koşul: aktif `aystek` tenant kaydı bulunmalıdır.
Komut mevcut `.env` yükleme davranışını kullanır; ortam değişkenleri önceliklidir.
Bu seçenekle tüm seed zinciri veya `shouldReset` çağrılmaz.
Yeni bir deployment için normal seed zinciri de siteyi oluşturur.
Mevcut sisteme yalnız Küresel Hat eklemek için tam seed çalıştırmak gerekli değildir.
Bu komut, henüz dağıtılmamış eski Core binary/image'larında mevcut sayılmamalıdır.

Gazete görünümünü test etmek için, site seed edildikten sonra örnek sayıları
ayrı ve açıkça opt-in olarak ekleyin:

```bash
go run ./cmd/seed --news-site=kureselhat --demo-newspapers
```

Bu komut Küresel Hat site scope'una toplam 15 özgün, `[ÖRNEK]` / `[DEMO]`
etiketli dijital sayı ve kapak kaydı ekler; önceki üç örneğe ek olarak 12 sayı
daha sunar. Tam reset yapmaz, AYSTEK Haber içeriğine dokunmaz ve tekrar
çalıştırıldığında editör değişikliklerini ezmez. Ortak
public-site anasayfasındaki gazete rafı, `/gazeteler` arşivi ve
`/gazeteler/{slug}` sayı görünümü mevcut `digital_editions` public API'sini
kullanır. Örnek HTML içerikleri doğrulanmış haber değildir.

### Veri koruma kuralları

- Store ve içerik varsayılanları `$setOnInsert` ile eklenir.
- Tekrarlama ID'leri veya `updated_at` değerlerini değiştirmez.
- Yönetici tarafından düzenlenmiş kayıtlar varsayılanlarla ezilmez.
- Silinmiş ana sayfa/menü/profil yeniden oluşturulmaz veya canlandırılmaz.
- Site pasifse, silinmişse veya domaini değiştirilmişse seed hata verir.
- Ana domain/alias başka tenant veya siteye aitse işlem yazmadan durur.
- Gelen context'teki başka tenant/site, yeni sitenin içerik kapsamına taşınmaz.
- Tenant kategorileri paylaşılır; siteye özel haber içerikleri paylaşılmaz.
- Yeni seed, mevcut AYSTEK Haber demo seed'ini çalıştırmaz.
- Seed çalıştırmak veritabanı migration'ı değildir; sonradan değişen şablon
  metinleri mevcut admin içeriğine otomatik uygulanmaz.

### Geliştirme ortamı uygulama kaydı

Veritabanı: yerel MongoDB üzerindeki `aystek_core`.
Son doğrulamadaki site ID: `6ab7fe6da38e0fee294e27bd`.
Bu ID yalnız bu geliştirme ortamının sonucudur; başka kurulumlarda sabit kullanılmaz.
Seed iki kez çalıştırıldı.

| Yeni siteye ait kayıt | Adet |
| --- | ---: |
| Store | 1 |
| Ana sayfa (TR/EN aynı kayıt) | 1 |
| Public menü | 9 |
| Public ayar | 2 |
| Editoryal profil | 1 |
| Haber | 0 |
| Yazar | 0 |
| Galeri | 0 |

Mevcut tenant kategorilerinin toplamı 19 olarak kaldı; yeni kategori kopyası oluşmadı.
AYSTEK'e ait mevcut store, page, menu, setting, view_profile, category, post,
author ve gallery belgelerinin işlem öncesi/sonrası SHA-256 özetleri karşılaştırıldı;
yeni site kayıtları hariç tüm karşılaştırmalar aynı çıktı.
Mevcut kullanıcı `.env` değişiklikleri korunmuştur.

### Test kanıtları

```bash
NEWS_SITE_TEST_MONGO_URI=mongodb://localhost:27017 \
  go test ./internal/seeds -run 'Test(NewsSite|NewsHome|AysTekHaber)' -count=1 -v

go test ./cmd/seed ./internal/modules/stores/repository \
  ./internal/modules/view_profiles/repository ./internal/platform/db
```

Mongo entegrasyon testleri her çalışmada rastgele isimli ayrı test DB'si oluşturur.
`MONGO_DB` geliştirme veritabanını test hedefi olarak kullanmaz.
Test sonunda yalnız oluşturduğu geçici DB'yi temizler.
`NEWS_SITE_TEST_MONGO_URI` verilmezse Mongo testleri skip olur; birim testleri çalışır.

Kapsanan davranışlar:

- Bilinmeyen seed kodu ve eksik tenant reddi.
- Küresel Hat marka/metadata üretimi; AYSTEK Haber marka sızıntısı olmaması.
- Ana domain ile `www` alias'ının aynı store/tenant'a çözülmesi.
- Başka site context'iyle çağrıda doğru yeni siteye yazılması.
- Mevcut kategori ve kardeş site sayfasının korunması.
- İkinci çalıştırmada yönetici düzenlemelerinin ve ID'lerin korunması.
- Silinmiş varsayılanların yeniden oluşturulmaması.
- Ana domain ve alias çakışmaları için yazmadan hata verme.
- Pasif/silinmiş/domaini değiştirilmiş store'u değiştirmeme.
- Mevcut AYSTEK Haber katalog, navigasyon, yazar ve profil testleri.

Çalışan Core API üzerinden kontrol:

| Domain | Endpoint | Sonuç |
| --- | --- | --- |
| `kureselhat.com` | `/public/pages/tr/anasayfa` | 200, Küresel Hat, yeni site ID |
| `www.kureselhat.com` | `/public/pages/tr/anasayfa` | 200, aynı site ID |
| `kureselhat.com` | `/public/posts?lang=tr&page=1&limit=1` | 200, 0 haber |
| `aystekhaber.com` | `/public/pages/tr/anasayfa` | 200, AYSTEK Haber, önceki site ID |
| Her iki kayıtlı domain | `/public/languages/active` | 200 |
| Her iki kayıtlı domain | `/public/settings` | 200 |

Testler `X-Store-Domain` header'ıyla yerel API'ye yapıldı; dış DNS yayını kanıtı değildir.

## Aşama 2 — ortak frontend (tamamlandı)

- [x] `lib/api/server.ts` sabit domain varsayılanından runtime deployment kimliğine taşındı.
- [x] Site kimliği `SITE_URL`, `AYSTEK_STORE_DOMAIN`, `AYSTEK_TENANT_ID` ile açıkça tanımlandı.
- [x] Host yalnız admission kontrolü oldu; Host tenant/site seçemez.
- [x] Bilinmeyen ve başka deployment'a ait Host 421 ile reddediliyor.
- [x] Core isteklerinde çağıranın site/store header override'ları siliniyor.
- [x] İstekler shared process cache kullanmadan site-kapsamlı ve `no-store` çalışıyor.
- [x] Marka adı site ana sayfasının lokalize başlığından; logo/favicon yalnız site branding ayarından geliyor.
- [x] Header/footer/logo/erişilebilirlik etiketleri ve generic placeholder dinamikleştirildi.
- [x] Sahte piyasa, hava, sosyal, canlı yayın ve uygulama mağazası fallback'leri kaldırıldı.
- [x] Yazar fallback'i markadan bağımsızlaştırıldı.
- [x] Layout ve tüm liste/detay sayfalarının metadata'sı site kimliğine bağlandı.
- [x] Canonical, OpenGraph, publisher ve JSON-LD ortak typed SEO üreticisine taşındı.
- [x] RSS, sitemap, robots, manifest ve newsletter source runtime site kimliğini kullanıyor.
- [x] Site dil listesi/default seçimi `/public/languages/active` sözleşmesini koruyor.
- [x] 429/invalid deployment/Core arızası “sayfa yok” 404'ü gibi gösterilmiyor.
- [x] Aynı build/image'ın farklı runtime kimlikleriyle çalışması contract testine alındı.
- [x] Site-scoped dijital gazete rafı, arşivi ve sayı detay görünümü eklendi; medya Core static hostundan site kimliğini koruyan proxy ile sunuluyor.
- [x] Küresel Hat için insert-only örnek gazete seed'i ve `--demo-newspapers` hedefli komutu doğrulandı.
- [x] Lint, TypeScript/production build ve Chromium masaüstü/mobil testi geçti.

Frontend sözleşme testi 8/8 geçti. Küresel Hat production build tarayıcı testi;
ana sayfa, yazarlar, galeriler, kategori, responsive menü, canonical/OpenGraph,
RSS, sitemap, robots, manifest, `www` alias'ı, bilinmeyen Host ve çapraz-domain
reddini doğruladı. Tarayıcı console/page error oluşmadı.

Bir operasyon bulgusu da kapatıldı: Core 429 döndüğünde kimlik çözümlemesi artık
404 üretmiyor. Next SSR tek sayfada birden çok Core çağrısı yaptığı için
`PUBLIC_RATE_REQ` değeri gerçek trafik testiyle kapasitelendirilmelidir. Sitemap
üretimi ayrıca dil başına menü/yazı/yazar ve link başına sayfa isteği yapmayacak
şekilde toplulaştırıldı. Varsayılan 60 istek/dakika sınırıyla arka arkaya sitemap,
ana sayfa ve alias kontrolünde Core 429 üretmedi; bu sonuç üretim yük testinin
yerini almaz.

## Aşama 3 — yayın hazırlığı (kısmen tamamlandı)

- [ ] Editör gerçek haber/yazar/galeri içeriğini yeni site bağlamında oluştursun.
- [ ] Hakkımızda, künye, iletişim ve hukuki sayfalar gerçek bilgilerle hazırlansın.
- [ ] Ana sayfa kurumsal linkleri oluşturulan sayfalarla eşleştirilsin.
- [ ] Site logosu/favicon, sosyal hesaplar ve gerekiyorsa analytics ayarlansın.
- [ ] Tenant'tan miras alınan diller gözden geçirilsin; gerekirse site dili kısıtlansın.
- [x] Site kimliğini koruyan installer domain/alias/port/install-path ayarları doğrulandı.
- [x] Localde tek kaynak, sunucuda domain başına ayrı source/install/config/log dizini belgelendi.
- [x] AYSTEK Haber 3001 ve Küresel Hat 3002 portlarıyla ayrı container/state sözleşmesine alındı.
- [x] Local-source ve registry Küresel Hat installer örnekleri eklendi.
- [x] GitHub ve temiz/commit edilmiş sunucu-içi Git kaynağından ayrı checkout
  oluşturan installer profilleri eklendi; yerel AI değerlendirmeli `auto-update`
  komutu ve varsayılan kapalı systemd timer örneği hazırlandı.
- [ ] Uzak sunucuda Git SSH/depo erişimi, Docker build, Ollama CPU hızı,
  başarısız image geri dönüşü ve gerçek trafik altında davranış kabul testinden geçsin.
- [ ] DNS A/AAAA, Nginx `server_name` ve TLS SAN kapsamı doğrulansın.
- [ ] HTTP/HTTPS, canonical, RSS, sitemap ve robots dış domainlerden test edilsin.
- [ ] Yeni sürüm ve rollback kaydı belgelensin.

Otomatik güncelleme yalnız başarılı sağlık kontrolünden sonra dağıtılmış commit'i
kaydeder; yeni image başarısızsa önceki image'a döner. Genel public-site
`rollback` komutu hâlâ otomatik eski image seçimi yapmaz. Düz dosya upload'u
yarım aktarım riski nedeniyle otomatik izlenmez; sunucu-içi Git profilinde
yükleme ancak tam ve temiz bir commit'ten sonra görünür olur. Mevcut Core
`AystekAIEngine` bilgi/kural motorudur, eğitilmiş LLM değildir; installer'ın
release-advisor'ı ayrı bir yerel Ollama modeli kullanır. Ücretli frontier
modellerle eşdeğer hız/kalite garanti edilmez.

`metadata.public_site_ready=false` hazırlık durumunu belirtir; bakım modu veya
erişim kontrolü değildir. DNS/proxy yayınını kendiliğinden engellemez.
Bu çalışma üretim ortamına deployment, DNS değişikliği veya sertifika kurulumu yapmadı.

## Operasyon ve güvenlik bulguları

- Son yerel `/health` kontrolünde MongoDB, Redis ve RabbitMQ `ok` döndü.
- Tanınmayan explicit public domain için settings/languages/pages/posts artık
  fallback yapmadan 400 ile reddediliyor; site-ID/header override saldırıları
  middleware entegrasyon testinde kapsanıyor.
- Core'u test için başlatırken personal-assistant agency scheduler'ın sekiz haber
  ve altı thread yazdığı görüldü. Yalnız bu test sırasında oluşan 14 kayıt kimlik
  ve zaman aralığıyla silindi; sonrasında aynı aralıkta kalan kayıt 0 doğrulandı.
  Bu kayıtlar geri alınamaz; kaynak sistemlerden tekrar çekilebilir.
- Aynı olayın tekrarlanmaması için `PA_AGENCY_SCHEDULER_ENABLED=false` açık
  opt-out'u ve testi eklendi. Browser doğrulaması bu ayarla çalıştırıldı ve yeni
  agency article/thread oluşmadığı doğrulandı.
- İlk çok sayfalı tarayıcı paketi Core'un varsayılan public limitini tüketti.
  Sitemap'in dil ve kurumsal link bazlı tekrar çağrıları tek site-kapsamlı veri
  setine indirildikten sonra varsayılan 60 istek/dakika ile sitemap + ana sayfa
  + alias smoke testi 429 olmadan geçti. Yine de üretim limiti gerçek eşzamanlı
  trafik ve proxy topolojisine göre kapasitelendirilmelidir.

## Operasyon dokümanları

Bu aşamayla birlikte aşağıdaki dokümanlara veri hazırlığı ve yayın sınırı işlendi:

- Installer `docs/KULLANIM_KILAVUZU.md`, Bölüm 21.1.
- Installer `README.md`, ortak haber uygulaması hazırlık notu.
- `user-guide/05-operations/SUNUCU_KURULUM_REHBERI.md`, Bölüm 15.1.
- `user-guide/05-operations/OPERATIONS_GUIDE.md`.
- `user-guide/05-operations/DEPLOYMENT_SCENARIOS.md`.

Installer runtime env şablonu ve doğrulama testleri güncellendi; mevcut projeden
ayrı Küresel Hat frontend reposu oluşturulmadı. Üretim yayını; editoryal/hukuki
içerik, logo/favicon, sosyal/analytics, DNS, TLS ve gerçek trafik kapasite testi
tamamlandıktan sonra onaylanmalıdır.
