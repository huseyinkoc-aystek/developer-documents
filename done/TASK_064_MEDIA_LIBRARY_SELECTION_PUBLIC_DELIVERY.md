# TASK 064 — Medya Kütüphanesi Seçim ve Public Görsel Teslimi

**Tarih:** 2026-10-13 · **Durum:** Kod, yerel doğrulamalar ve GitHub aktarımı tamamlandı;
canlı dağıtım VPN rotası bekliyor.

## Sorun

Yönetici medya kütüphanesine yüklenen görsel, içerik formundaki görsel seçicide listelenmiyor;
CMS kaydına `/static/...` yolu yazıldığında haber public sitesinde görsel görünmüyordu.

## Kök nedenler

1. Panel görsel seçici `GET /admin/media?mime=image/` çağrısı yapıyordu. Core repository bu değeri
   gerçek `mime_type` alanı yerine yalnız `jpg`, `png` gibi uzantı tutan `file_type` alanında
   arıyordu; sonuç kümesi bu nedenle boştu.
2. CMS medya yolları taşınabilir biçimde `/static/...` saklanır. Ortak haber public uygulaması bu
   yolu kendi origin'inde doğrudan kullanıyor, fakat dosya Core upload volume'ünde bulunduğundan
   public sitede 404 oluşuyordu.

## Çözüm

- Core medya filtresi:
  - `image/` ve `image/*` → `mime_type` önek eşleşmesi.
  - `application/pdf` gibi tam MIME → `mime_type` tam eşleşmesi.
  - `png`/`.png` gibi geriye uyumlu uzantı → `file_type` tam eşleşmesi.
- Panel yükleme sonrası ilk dosyayı otomatik seçer, aynı dosyanın tekrar yüklenebilmesi için file
  input'u temizler ve medya eylem düğmelerini açıkça `type="button"` olarak işaretler.
- Ortak haber public uygulamasına site-kapsamlı, same-origin `/media/static/...` görsel proxy'si
  eklendi. Proxy yalnız doğrulanmış görsel MIME türlerini döndürür; path traversal reddedilir,
  `nosniff` ve kısıtlı CSP uygulanır.
- Haber, SEO/JSON-LD, yazar, galeri ve reklam görsel yolları ortak `mediaUrl` adaptöründen geçirilir.

## Değişen projeler

- `aystek-core`: `internal/modules/media/repository/mongo.go` ve repository testleri.
- `aystek-core-next`: `app/components/MediaModal.tsx`.
- `aystekhaber-com`: `lib/media.ts`, `/media/[...path]` route'u ve haber/yazar/galeri/reklam/SEO
  adaptörleri. Bu tek kod tabanı Küresel Hat ve AYSTEK Haber dağıtımlarında kullanılabilir.

## Doğrulama

- Core medya repository testleri başarılı.
- Panel değişen dosya ESLint kontrolü temiz; production build başarılı.
- Panel repo-geneli TypeScript/lint kapıları bu görevden önce var olan teknik borçlar nedeniyle
  kırık; değişen medya modalında yeni hata yok.
- Ortak haber sitesi: 12/12 site sözleşme testi, ESLint ve production build başarılı.
- `git diff --check` temiz.

## Sürümler

- `aystek-core`: `21ef8c6` (`fix(media): filter library by MIME type`)
- `aystek-core-next`: `5c910c2` (`fix(media): select uploaded assets in picker`)
- `aystekhaber-com`: `f841f86` (`fix(media): proxy Core images on public sites`)
- Üç commit de ilgili GitHub `main` dalına aktarıldı.

## Canlı dağıtım engeli

2026-10-13 tarihinde uzak sunucuya üç kez SSH bağlantısı denendi. Çalışma
ortamında VPN tünel arayüzü ve `192.168.200.0/24` rotası bulunmadığı için
`192.168.200.250:22` bağlantıları kimlik doğrulama aşamasına gelmeden zaman aşımına
uğradı. VPN rotası geri geldiğinde aşağıdaki kontrol listesi installer üzerinden
uygulanacak. Upload volume silinmeyecek veya yeniden oluşturulmayacak.

## Canlı kontrol listesi

- Core ve panel benzersiz `git-<sha>` image'larıyla installer üzerinden dağıtılacak.
- Küresel Hat ortak public image'ı yeni haber commit'iyle dağıtılacak.
- Güncelleme öncesi/sonrası Core upload SHA-256 manifestleri karşılaştırılacak.
- `mime=image/` çağrısının yüklenen görseli döndürdüğü, Core static URL'sinin ve public proxy
  URL'sinin 200 + doğru `Content-Type` döndürdüğü doğrulanacak.
- Core/panel/Küresel Hat health ve temel public rotalar kontrol edilecek.
