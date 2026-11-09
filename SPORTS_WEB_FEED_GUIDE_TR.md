# Spor Merkezi — Ücretsiz Web Feed / Google E-Tablo Rehberi

**Durum:** Uygulandı · **Tarih:** 2026-09-29  
**Ana plan:** `SPORTS_MATCHES_STANDINGS_WORKPLAN_TR.md` · **Zorunlu standart:** `README_FOR_AI_ASSISTANTS.md`

Bu rehber, ücretli spor API'si kullanmadan puan tablosu ve maç verisi yayınlamak isteyen site
yöneticileri ile çalışmayı devralacak yapay zekâ araçları içindir. Sistem rastgele Google arama
sonuçlarını veya HTML sayfalarını kazımaz. Bunun yerine lisansı ve kaynağı yönetici tarafından
doğrulanmış, standart JSON/CSV beslemesini okur.

## 1. Site kapsamı

- Yarışmalar zaten `tenant_id + site_id + code` ile ayrılır.
- Senkronizasyon ayarı ve sağlayıcı anahtarları `site.<site_id>.sports.sync` anahtarında saklanır.
- Aynı tenant altındaki iki site farklı sağlayıcı, API anahtarı ve çalışma aralığı kullanabilir.
- Eski `sports.sync` tenant kaydı yalnızca site kaydı henüz yoksa geçiş uyumluluğu için okunur;
  yönetim panelinden ilk kayıtla siteye özel anahtar oluşturulur.

## 2. Yönetici için en kolay ücretsiz akış

1. Google E-Tablolar'da boş bir dosya açın.
2. `templates/sports_web_feed_template.csv` dosyasındaki başlıkları ve örnek satırları içe aktarın.
3. Gerçek takım, puan ve maç satırlarını girin. Kullanmadığınız hücreleri boş bırakabilirsiniz.
4. Google E-Tablolar'da **Dosya → Paylaş → Web'de yayınla** yolunu açın.
5. İlgili sayfayı seçip çıktı biçimini **Virgülle ayrılmış değerler (.csv)** yapın.
6. Core yönetim panelinde **Spor Merkezi → yarışma → Otomatik Senkronizasyon** bölümünü açın.
7. Sağlayıcı olarak **Web Feed / Google Sheets** seçin.
8. CSV adresini, görünen kaynak/atıf adını ve gerçek kaynak sayfasının HTTPS adresini girin.
9. Yarışmayı kaydedin. Spor Merkezi listesinden Web Feed profilini ve genel senkronizasyonu açın;
   ardından **Şimdi Çalıştır** ile doğrulayın.

Google E-Tablo içinde `IMPORTHTML`, `IMPORTXML` veya Apps Script kullanılabilir; ancak dış sitenin
kullanım koşulları, robot politikası, telif/lisansı ve verinin yeniden yayımlanmasına izin verip
vermediği yönetici tarafından doğrulanmalıdır. Core bu hukuki izni varsaymaz.

## 3. CSV sözleşmesi

Tek sayfada iki satır tipi bulunur:

- `record_type=standing`: puan tablosu satırı.
- `record_type=match`: maç satırı.

| Alan | Puan satırı | Maç satırı | Açıklama |
|---|---:|---:|---|
| `record_type` | zorunlu | zorunlu | `standing` veya `match` |
| `rank`, `team` | zorunlu | boş | sıra ve takım adı |
| `played`, `won`, `drawn`, `lost`, `points` | boşsa 0 | boş | negatif olamaz |
| `form` | isteğe bağlı | boş | `W/D/L` veya Türkçe `G/B/M` |
| `home_team`, `away_team` | boş | zorunlu | takım adları |
| `home_short`, `away_short` | boş | isteğe bağlı | kısa ad |
| `home_score`, `away_score`, `minute` | boş | boşsa 0 | negatif olamaz |
| `status` | boş | isteğe bağlı | `scheduled`, `live`, `finished`; boşsa `scheduled` |
| `starts_at` | boş | isteğe bağlı | RFC3339, ör. `2026-09-28T18:00:00Z` |

Başlıklar büyük/küçük harfe duyarlı değildir; boşluk ve tire `_` olarak normalize edilir. En az
bir geçerli `standing` satırı zorunludur. Core bütün belgeyi doğrular; bozuk satır varsa mevcut
yayındaki son doğru veri korunur.

## 4. JSON sözleşmesi

```json
{
  "standings": [
    {
      "rank": 1,
      "team": "Örnek Spor",
      "played": 2,
      "won": 2,
      "drawn": 0,
      "lost": 0,
      "points": 6,
      "form": "WW"
    }
  ],
  "matches": [
    {
      "home_team": "Örnek Spor",
      "home_short": "ÖRS",
      "away_team": "Örnek Gençlik",
      "away_short": "ÖRG",
      "home_score": 2,
      "away_score": 1,
      "status": "finished",
      "minute": 90,
      "starts_at": "2026-09-28T18:00:00Z"
    }
  ]
}
```

JSON'da bilinmeyen alanlar hata sayılır. Böylece yazım hatası sessizce veri kaybına dönüşmez.

## 5. Güvenlik ve yayın ilkeleri

- Yalnızca HTTPS besleme adresleri kabul edilir.
- URL içinde kullanıcı adı/parola kabul edilmez.
- Loopback, özel ağ, link-local, multicast ve belirsiz IP adresleri reddedilir.
- Her yönlendirme yeniden doğrulanır; en fazla üç yönlendirme izlenir.
- Ortam proxy'si devre dışıdır; URL'nin proxy üzerinden iç ağa erişmesi engellenir.
- Yanıt süresi 15 saniye, açılmış içerik boyutu 2 MiB ile sınırlıdır.
- Kaynak adı/atıf zorunludur; kaynak sayfası yalnızca HTTPS olabilir.
- Hatalı veya erişilemeyen besleme mevcut son doğru veriyi silmez.

## 6. Sağlayıcı davranışı

| Sağlayıcı | Anahtar | Kapsam | Kullanım |
|---|---|---|---|
| API-Sports | zorunlu | futbol, basketbol, voleybol | kapsamlı ticari/ücretsiz kota |
| TheSportsDB | isteğe bağlı | doğrulanmış ana kapsam futbol | yerleşik resmi ücretsiz anahtar |
| Web Feed / Google Sheets | gerekmez | tüm desteklenen sporlar | kullanıcı kontrollü JSON/CSV |

Sağlayıcı profilleri site bazında aynı anda açık olabilir. Her yarışma kendi `external_ref.provider`
değeriyle doğru profile bağlanır. Ücretli API anahtarı olmayan yarışmalar Web Feed'e geçirilebilir;
diğer yarışmaların sağlayıcısı değişmez.

## 7. Teknik devam notları

Backend:

- `internal/modules/sports/provider/webfeed.go`: güvenli HTTP istemcisi ve JSON/CSV ayrıştırıcı.
- `internal/modules/sports/service/sync.go`: site anahtarı, çoklu sağlayıcı profilleri, katalog.
- `GET /admin/sports/sync/providers`: çalışma zamanındaki sağlayıcı kataloğu.
- `GET /admin/sports/sync/feed-template`: kimlik doğrulamalı örnek CSV indirme.
- `GET|PUT /admin/sports/sync/config`: seçili siteye ait ayarlar.
- `POST /admin/sports/sync/run`: seçili sitede anlık çalıştırma.

Admin:

- `aystek-core-next/app/(dashboard)/admin/sports/page.tsx`: site sağlayıcı profilleri.
- `aystek-core-next/app/(dashboard)/admin/sports/[code]/page.tsx`: yarışma sağlayıcı/URL/atıf alanları.

Public:

- `aystekhaber-com/components/news/SportsSnapshotPanel.tsx`: kaynak atfını güvenli bağlantıyla gösterir.

Sonraki olası fazlar, mevcut iş tamamlanmadan başlanmış sayılmaz:

1. Yönetim panelinde URL'yi kaydetmeden önizleme ve satır-hata raporu.
2. Aynı yarışma için sıralı birincil/yedek kaynak zinciri.
3. Onay kuyruğu: otomatik değişiklikleri yayınlamadan önce editör karşılaştırması.
4. Sağlayıcı sağlık/kota ekranı ve başarısızlık bildirimleri.
5. Yetkili federasyonların farklı CSV biçimlerini eşleyen yönetilebilir dönüşüm şablonları.
