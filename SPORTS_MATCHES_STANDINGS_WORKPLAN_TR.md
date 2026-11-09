# Spor: Son Maçlar ve Puan Tablosu — İş Planı

**Sahibi:** Hüseyin Koç · **Oluşturma:** 2026-09-27 · **Kaynak:** Bu konuşma (Claude Sonnet 5, i18n çalışmasının duraklatılmasının ardından)
**Durum:** ✅ FAZ 0-6 TAMAMLANDI (2026-09-27/28). Üç sağlayıcı (API-Sports + TheSportsDB +
Web Feed/Google Sheets) uygulandı. Ücretsiz ve kullanıcı kontrollü JSON/CSV köprüsü sayesinde
basketbol/voleybol dahil desteklenen sporlar abonelik olmadan da beslenebilir. API-Sports ve
TheSportsDB'nin önceki canlı doğrulamaları korunur; Web Feed sözleşmesi birim testleriyle
doğrulanmıştır. Ayrıntılı kullanım: `SPORTS_WEB_FEED_GUIDE_TR.md`.

**MUTLAKA ÖNCE OKUYUN:** `/workspace/aystek-core-dev-docs/README_FOR_AI_ASSISTANTS.md` (bu doküman onun kurallarına tabidir — 6 çekirdek dosya modül yapısı, 28 dil çeviri kuralı, menu/permission/view-profile senkronizasyonu, admin sayfa görsel standardı). Herhangi bir çelişkide README kazanır.

Bu doküman, aystekhaber-com ve kureselhat.com (aynı Next.js kod tabanı, site-farkında çalışma —
bkz. `done/TASK_060_SHARED_NEWS_PUBLIC_SITE.md`) için **voleybol (üst lig + bir alt lig), futbol,
basketbol ve bunların uluslararası ayakları**nı kapsayan "son maçlar" ve "puan tablosu" özelliğini
kapsamlı hale getirmek için yazılmıştır. Farklı AI araçları (Codex, Claude, Copilot, ...) bu planı
yarım bıraksa da bir sonraki oturum kaldığı yerden devam edebilmelidir — her faz bağımsız
doğrulanabilir bir teslim birimidir.

---

## 1. Mevcut Durum (2026-09-27 taraması — GAP ANALİZİ)

**Önemli:** Bu özellik sıfırdan kurulmuyor. `internal/modules/sports/` adında DDD uyumlu bir modül
**zaten var** ve production'da çalışıyor. Görev bunu genişletmek, yeniden yazmak değil.

| Katman | Dosya | Şu anki durum |
|---|---|---|
| Model | `internal/modules/sports/model/league.go` | `LeagueSnapshot{Code, LocalizedName, Season, Week, Standings[], Matches[], Active}` — **tenant+site+code başına TEK belge**, üzerine yazılır (Upsert). Spor türü, lig seviyesi (üst/alt), yurt içi/uluslararası ayrımı **yok**. |
| Repository | `internal/modules/sports/repository/{interface,mongo,indexes}.go` | `FindByCode`, `FindActiveByCode`, `Upsert` — kod bazlı tekil sorgu. **Liste/filtreleme ucu yok** ("bu tenant'ın tüm aktif liglerini getir" diye bir metod yok). |
| Service | `internal/modules/sports/service/service.go` | `Get`, `GetActive`, `Upsert`. Sıralama/validasyon iyi durumda (bkz. `model.Validate()` — rank/takım çakışması, played≥won+drawn+lost vb. zaten kontrol ediliyor, buna dokunmaya gerek yok). |
| Handler/Routes | `internal/modules/sports/handler/{admin,public}.go`, `routes.go` | `GET/PUT /admin/sports/leagues/:code` (izin: `sports:read`/`sports:update` — `permission_seed.go:119`'da `"sports"` modülü zaten kayıtlı), `GET /public/sports/leagues/:code`. **Liste ucu yok.** |
| Menü | `menu_seed.go:159` | `"sports"` anahtarı, parent `content`, path `/admin/sports`, ikon `trophy`, roller `admin+editor` — **zaten mevcut**, dokunmaya gerek yok (view-profile senkronizasyonu da muhtemelen zaten yapılmış, Faz 4'te doğrulanacak). |
| Admin UI | `aystek-core-next/app/(dashboard)/admin/sports/page.tsx` (187 satır) | Tek `code` için tek form — çoklu lig/spor seçici **yok**. |
| Public API client | `aystekhaber-com/lib/api/sports.ts` | `fetchLeagueSnapshot(code = 'super-lig')` — **kod sabit kodlanmış varsayılan**, çoklu lig çağrısı yok. |
| Public widget | `aystekhaber-com/components/news/SportsSnapshotPanel.tsx` | Ana sayfada tek panel: son 4 maç + ilk 6 sıra. "Tüm Maçlar →" / "Tam Tablo →" linkleri **hedefsiz** (henüz bir sayfaya gitmiyor). Spor/lig sekmesi yok. |
| Çeviri anahtarları | `translation_seed.go` (grep: `sports.*`, `news.labels.latest_matches`, `news.labels.standings`, `news.labels.all_matches`, `news.labels.full_table`, `news.labels.week`, `news.labels.final`, `news.labels.scheduled`) | Temel etiketler zaten 28 dilde mevcut — **yeni eklenecek anahtarlar bunları tekrar etmemeli**, sadece eksik olanları (spor türü adları, lig/tier etiketleri, "tüm sporlar" sayfası metinleri) eklemeli. |
| Otomatik veri çekme | — | **Hiç yok.** Bugün her şey admin panelinden elle giriliyor. |

**Sonuç:** Faz 1-2 çoğunlukla *genişletme*; Faz 3 (otomatik senkronizasyon) gerçek bir *yeni* alt sistem.

---

## 2. Mimari Kararlar

### 2.1 Veri kaynağı stratejisi (KULLANICI TARAFINDAN ONAYLANDI — 2026-09-27)

Kullanıcıya AskUserQuestion ile soruldu: "Google üzerinden çekelim" ifadesi tam olarak neyi
kastediyor? **Cevap: gerçek bir spor verisi API'si (önerilen seçenek).** Yani:

- ✅ Birincil veri kaynağı: **admin panelinden manuel giriş** (bugünkü davranış, değişmiyor).
- ✅ Otomatik fallback: admin bir ligi güncellemediyse, **gerçek bir spor verisi API'si** (ör.
  API-Football, TheSportsDB, veya benzeri — API seçimi Faz 3'te teknik+maliyet değerlendirmesiyle
  yapılacak, bkz. §2.4) planlı bir job ile veriyi çeker ve `Source: "auto"` ile işaretler.
- ❌ Google arama sonuçlarını kazımak (scraping) **reddedildi** — hem ToS riski hem kırılganlık
  hem de `done/TASK_060_SUMMARY.md`'de zaten kayıtlı olan ekip felsefesiyle çelişiyor: *"Removed
  fabricated weather, market, social, live and app-store fallbacks"* — yani bu ekip güvenilmez/
  sahte veri yerine dürüst boşluk göstermeyi tercih ediyor. Google kazıma bu ilkeyi ihlal eder.
- Bu karar geri döndürülebilir: eğer Faz 3'e gelindiğinde uygun/ücretsiz bir spor API'si
  bulunamazsa, kullanıcıya tekrar danışılmalı — asla sessizce kazımaya geçilmemeli.

### 2.2 Veri modeli genişletmesi (Faz 1)

`LeagueSnapshot`'a şu alanlar eklenir (geriye dönük uyumlu — mevcut `super-lig` belgesi bozulmaz):

```go
type SportType string
const (
    SportFootball   SportType = "football"
    SportBasketball SportType = "basketball"
    SportVolleyball SportType = "volleyball"
)

type CompetitionTier string
const (
    TierTop           CompetitionTier = "top"            // Süper Lig, BSL, Sultanlar/Efeler Ligi
    TierSecond        CompetitionTier = "second"          // TFF 1. Lig, TBL, 1. Lig
    TierInternational CompetitionTier = "international"   // UEFA/EuroLeague/CEV vb.
)

type DataSource string
const (
    SourceManual DataSource = "manual" // admin panelinden girildi
    SourceAuto   DataSource = "auto"   // Faz 3 sync job'ı tarafından çekildi
)

// LeagueSnapshot'a eklenecek alanlar:
Sport        SportType       `bson:"sport" json:"sport"`
Tier         CompetitionTier `bson:"tier" json:"tier"`
Gender       string          `bson:"gender,omitempty" json:"gender,omitempty"` // voleybolda kadın/erkek ayrımı için, boşsa "male" varsayılır
Source       DataSource      `bson:"source" json:"source"`
SourceSyncedAt *time.Time    `bson:"source_synced_at,omitempty" json:"source_synced_at,omitempty"`
```

`Code` alanı benzersiz anahtar olmaya devam eder (ör. `super-lig`, `tff-1-lig`, `ucl`, `bsl`,
`tbl`, `euroleague`, `sultanlar-ligi`, `efeler-ligi`, `voleybol-1-lig`, `cev-champions-league`).
Bu, mevcut `FindByCode`/`Upsert` sözleşmesini bozmadan çoklu yarışmayı destekler.

**Repository'ye eklenecek:** `ListActive(ctx, filter{Sport, Tier})` — site'nin tüm aktif
yarışmalarını (veya spor/tier'a göre filtrelenmiş) listeleyen yeni bir metod. Bu, hem admin
listeleme sayfası hem de public "Spor" hub sayfası için gerekli.

### 2.3 Otomatik senkronizasyon mimarisi (Faz 3) — mevcut desenlerin tekrar kullanımı

Sıfırdan tasarım yapmaya gerek yok, kod tabanında **birebir uyan** iki emsal var:

1. **Sağlayıcı soyutlaması + admin-yönetilir config:** `internal/modules/currencies/service/fx_sync.go`
   desenini kopyala — `SportsSyncConfig{Enabled, Provider, IntervalMinutes, Endpoint, APIKey, ...}`
   bir `settings` kaydı olarak saklanır (`settingsservice.Service` üzerinden, `currencies.fx_sync`
   anahtarına benzer şekilde `sports.sync` gibi bir anahtar), `providers map[string]SportsProvider`
   provider registry'si (`"manual"` no-op + gerçek API adaptörü). Kimlik bilgisi şifreleme için
   `internal/modules/fleet_tally/credcrypt/crypt.go` (gps_providers'ın kullandığı) tekrar kullanılabilir.
2. **Periyodik job:** `internal/modules/fleet_tally/route_fares/jobs/expiring_fare_scheduler.go`
   desenini kopyala — `StartSportsSyncScheduler(ctx, logger, svc)`, `contextkeys.WithIsPublic`
   ile tag'lenmiş context, `ListDistinctTenantSiteScopes()` ile tüm tenant+site'ları gez (arka plan
   job'ının tek bir request-scoped tenant'ı yok), `time.NewTimer(interval)` döngüsü, `ops.Record*`
   ile metrik. `internal/platform/jobs/fx_sync_scheduler.go` da doğrudan referans alınabilir (retry
   + backoff mantığı oradan).

Bu ikisini birleştirmek: `sports/service/sync.go` (yeni dosya) + `internal/platform/jobs/sports_sync_scheduler.go`
(yeni dosya) — ikisi de var olan iki dosyanın neredeyse doğrudan uyarlanmış hali olacak, yeni bir
mimari icat edilmiyor.

### 2.4 Spor API sağlayıcı seçimi — KARARLAŞTIRILDI (2026-09-28, kullanıcı onayı: "API-Sports ailesiyle ilerle")

**Seçilen sağlayıcı: api-sports.io ailesi** — üç ayrı alt alan adı, aynı hesap/anahtar formatı:
- Futbol: `v3.football.api-sports.io` — `x-apisports-key` header, `/fixtures` + `/standings`,
  `?league={id}&season={year}` parametreleri. Response zarfı: `{get,parameters,errors,results,
  paging,response}`; standings `response[0].league.standings[0]` altında (çoklu grup/tur varsa
  iç içe dizi). **Şekil belgelenmiş ve yüksek güvenle uygulandı.**
- Basketbol: `v1.basketball.api-sports.io`, Voleybol: `v1.volleyball.api-sports.io` — aynı desen,
  `/games` + `/standings`. **Şekil aynı alan adlarıyla uygulandı ama canlı bir anahtarla henüz
  doğrulanmadı** — `internal/modules/sports/provider/apisports.go` içindeki `parseStandingEntry`/
  `parseMatchEntry` toleranslı (birden çok olası alan adını dener: `all`/`games`, `win.total`,
  `goals`/`scores`), beklenmeyen bir şekil gelirse sessizce yanlış veri üretmek yerine
  `ErrAPISportsResponseShape` döndürür.
- Voleybol dahil üç spor da kapsanıyor; Faz 3 kapsam daraltması gerekmedi.

**Uygulanan mimari (kod: bkz. §2.3'ün gerçek karşılıkları):**
- `internal/modules/sports/model/league.go`: `ExternalRef{Provider,LeagueID,Season}` — her
  yarışma kaydına opsiyonel, nil ise "sadece manuel", sync job asla dokunmaz.
- `internal/modules/sports/credcrypt/crypt.go`: AES-256-GCM, `SPORTS_CRED_KEY` env, emsal
  (`fleet_tally/credcrypt`, `personal_assistant/credcrypto`) ile aynı desen — paylaşılan platform
  paketi değil, üçüncü modül-özel kopya (bilinçli, bkz. proje konvansiyonu).
- `internal/modules/sports/service/sync.go`: `SyncConfig` (`settings` anahtarı `sports.sync`),
  `SyncService.RunOnce(ctx)` — sadece `ExternalRef` set edilmiş VE `ExternalRef.Provider` aktif
  sağlayıcıyla eşleşen kayıtları senkronlar, `Source=auto` + `SourceSyncedAt` damgalar (admin
  `Upsert` yolu her zaman `Source=manual` damgalamaya devam ediyor, iki yol asla çakışmıyor).
- `internal/modules/sports/provider/apisports.go`: HTTP adaptörü, `httptest.Server` ile test
  edildi (canlı anahtar gerekmeden).
- `internal/platform/jobs/sports_sync_scheduler.go` + `internal/platform/ops/sports_sync_metrics.go`:
  `route_fares/jobs/expiring_fare_scheduler.go` + `fx_sync_scheduler.go` desenlerinin doğrudan
  uyarlanmış hali — `ListDistinctTenantSiteScopes()` ile gez, retry+backoff, Prometheus metrikleri.
- `internal/modules/sports/handler/sync_admin.go` + admin panel (`app/(dashboard)/admin/sports/
  page.tsx` üst kısmı): açık/kapalı, sağlayıcı, aralık, API anahtarı (asla geri okunmaz, sadece
  "configured" bayrağı döner) + "Şimdi Çalıştır" butonu. `[code]/page.tsx`'e her yarışma için
  ayrı "Otomatik Senkronizasyon" bölümü eklendi (opt-in checkbox + provider/league_id/season).
- `cmd/api/main.go`: scheduler varsayılan **kapalı** — `sports.sync` settings kaydı hiç
  oluşturulmadıysa `LoadConfig` → `ErrSyncConfigMissing` → scheduler tick'i atlar. Admin panelden
  açılana kadar hiçbir dış çağrı yapılmaz.

**AÇIK KALAN:** Kullanıcıdan henüz gerçek bir api-sports.io API anahtarı istenmedi — bu, uçtan uca
canlı doğrulama (özellikle basketbol/voleybol response şekli) yapılacağı zaman gerekli olacak.
README'nin "When to Ask" kuralı gereği, canlı test öncesi kullanıcıya sorulmalı.

---

## 3. Kapsam: Taslak Yarışma Listesi (ONAY BEKLİYOR — düşük risk, plan mimarisini bloklamıyor)

Aşağıdaki liste Türkiye pazarı için makul bir varsayılan taslaktır. Editoryal ekip/kullanıcı
onayı gerekir ama bu onay Faz 1-2'nin mimari çalışmasını **bloklamaz** — `code` alanı serbest metin
olduğu için liste Faz 5'e kadar herhangi bir noktada genişletilip daraltılabilir.

| Spor | Üst lig | Alt lig | Uluslararası |
|---|---|---|---|
| Futbol | Süper Lig (`super-lig`) | TFF 1. Lig (`tff-1-lig`) | UEFA Şampiyonlar Ligi (`ucl`), UEFA Avrupa Ligi (`uel`), UEFA Konferans Ligi (`uecl`) |
| Basketbol | Basketbol Süper Ligi (`bsl`) | TBL Türkiye Basketbol Ligi (`tbl`) | EuroLeague (`euroleague`), EuroCup (`eurocup`) |
| Voleybol | Sultanlar Ligi kadın (`sultanlar-ligi`), Efeler Ligi erkek (`efeler-ligi`) | 1. Lig kadın/erkek (`voleybol-1-lig-kadin`, `voleybol-1-lig-erkek`) | CEV Şampiyonlar Ligi (`cev-champions-league`) |

**Not:** Voleybolda kadın/erkek ayrımı `Gender` alanıyla modellenir (§2.2). Liste kesinleşmeden
Faz 1-2'ye başlanabilir; ilk seed sadece mevcut `super-lig` kaydını yeni alanlarla günceller
(`sport=football, tier=top`), diğer kodlar admin panelinden veya Faz 1 seed'inde boş/placeholder
olarak eklenir.

---

## 4. Fazlar

### Faz 0 — Hazırlık (herhangi bir AI, ~30 dk) — ✅ Tamamlandı (2026-09-27)
- [x] Bu dokümanı ve README_FOR_AI_ASSISTANTS.md'nin ilgili bölümlerini oku (Module Creation
      Checklist, Admin Page Visual Design Standard, RULE #1/#1.1/#2/#2.1/#2.2/#3).
- [x] `internal/modules/sports/` mevcut testlerini çalıştır: `go test ./internal/modules/sports/...`
      — genişletmeye başlamadan önce temiz olduğunu doğrula.
- [x] `aystek-core-next/app/(dashboard)/admin/sports/page.tsx` mevcut halini oku (187 satır).

### Faz 1 — Backend: çoklu yarışma veri modeli — ✅ Tamamlandı (2026-09-27)
- [x] `model/league.go`'ya `Sport`, `Tier`, `Gender`, `Source`, `SourceSyncedAt` alanlarını ekle
      (§2.2). `Validate()`'e `Sport`/`Tier`/`Source` için enum kontrolü eklendi.
- [x] `repository/interface.go` + `mongo.go`'ya `List`/`ListActive(ctx, filter)` eklendi
      (site-scoped, opsiyonel sport/tier filtresi — admin için tümü, public için sadece aktif,
      `FindByCode`/`FindActiveByCode` ikilisiyle aynı desen). `repository/indexes.go`'ya
      `(tenant_id, site_id, active, sport, tier)` bileşik index eklendi.
- [x] `service/service.go`'ya `List`/`ListActive` wrapper'ları eklendi; okuma yollarında
      (`Get`/`GetActive`/`List`/`ListActive`) `Normalize()` çağrısıyla eski belgelerde
      Sport/Tier/Source geriye dönük varsayılanlarla dolduruluyor. `Upsert` artık `Source`'u
      her zaman `manual` olarak damgalıyor (Faz 3'teki otomatik senkron ayrı bir yoldan yazacak).
- [x] `handler/admin.go` + `handler/public.go`'ya liste uçları eklendi: `GET /admin/sports/leagues`
      (opsiyonel `?sport=&tier=` query, tümü), `GET /public/sports/leagues` (sadece aktif).
- [x] `routes.go` güncellendi.
- [x] `internal/seeds/sports_seed.go` güncellendi (`newsTemplateSportsSnapshot` artık
      `Sport: football, Tier: top` set ediyor; `seedNewsDemoSports`'un ham bson `fields` map'i de
      sport/tier/gender/source taşıyor). Geriye dönük uyumluluk **gerçek prod-yansıtan dev
      veritabanında** (`aystek_core`, 2026-09-21 tarihli mevcut `super-lig` kaydı) doğrulandı:
      canlı API `GET /public/sports/leagues/super-lig` çağrısı `sport=football, tier=top,
      source=manual` döndürdü — migration scripti gerekmedi.
- [x] `service/service_test.go` + `handler/public_test.go`'ya yeni alan/uç testleri eklendi
      (`TestUpsertAlwaysStampsManualSource`, `TestListAndListActiveBackfillDefaultsOnOlderDocuments`,
      `TestPublicListReturnsItemsAndTotal`).
- [x] `go build ./...` ve `go test ./internal/modules/sports/... ./internal/seeds/...` **gerçek
      MongoDB'ye karşı** (idempotency/site-izolasyon entegrasyon testleri dahil) yeşil.
- [x] 24 yeni `sports.*` çeviri anahtarı (spor türü/seviye/cinsiyet/liste sayfası metinleri)
      **28 dilde genuine çeviriyle** eklendi (allLangs placeholder değil — bu oturumdaki i18n
      borcu kapatma çalışmasıyla tutarlı, yeni borç yaratılmadı).

### Faz 2 — Admin UI + Public hub sayfası — ✅ Tamamlandı (2026-09-27)
- [x] Admin: `/admin/sports/page.tsx` liste + spor/tier filtre çipleri + "Yeni Yarışma" akışına
      dönüştürüldü (`AdminDataTable` kullanıyor, gradient hero header + `rounded-2xl`/
      `border-slate-100`/`shadow-sm` kart standardına uygun). Düzenleme `[code]/page.tsx`'e
      taşındı (mevcut standings/matches editörü korunarak Sport/Tier/Gender select'leri eklendi,
      "Geri" pill butonu + kaynak rozeti eklendi; yeni kod girilirse boş şablonla açılıyor).
- [x] `aystek-core-next/lib/api/sports.ts`'e `list()` + `SportType`/`CompetitionTier`/
      `DataSource` tipleri eklendi.
- [x] Public: `aystekhaber-com/lib/api/sports.ts`'e `fetchActiveLeagues(filter)` eklendi.
  - [x] Yeni "Spor" hub sayfası: `app/spor/page.tsx` — sport/tier filtre çipleri (`?sport=&tier=`
        query, sunucu taraflı render, client JS gerekmiyor), her aktif yarışma için tam panel.
  - [x] `app/spor/[code]/page.tsx` — tek yarışmanın tam sayfası (`SportsSnapshotPanel`'e eklenen
        `full` prop'uyla, dilimlenmemiş liste).
  - [x] `SportsSnapshotPanel.tsx`'teki ölü "Tüm Maçlar →" / "Tam Tablo →" linkleri artık
        `/spor/[code]`'a yönleniyor (ana sayfa widget'ında doğrulandı).
  - [x] Ana sayfa widget'ı çoklu spora genişletildi (kullanıcı onayı alındı, 2026-09-27:
        "tasarım çizgisine uygun işlerine devam et"). `app/page.tsx` artık
        `fetchActiveLeagues({tier:'top'})` çekiyor; yeni `components/news/SportsWidget.tsx`
        (client) birden fazla üst-lig yarışması varsa `/spor` hub'daki pill-filtre tasarımıyla
        birebir aynı bir spor seçici gösteriyor, veri tek seferde sunucuda çekiliyor. Tek
        yarışma varken (bugünkü gerçek durum) sekme satırı hiç render edilmiyor — önceki
        görünümle birebir aynı. Geçici bir ikinci üst-lig kaydıyla (basketbol) canlı
        doğrulandı, sonra temizlendi.
- [x] Site izolasyonu: yarışma verisi zaten `tenant_id+site_id` taşıyor (mevcut `Upsert`/
      `FindByCode` deseni), yeni `List`/`ListActive` de aynı `tenantfilter.WithTenantAndSite`
      üzerinden filtreleniyor — kureselhat.com ek kod gerektirmedi.
- [x] Doğrulama: `npx tsc --noEmit` her iki frontend'de de sports dosyalarında 0 hata (aystekhaber-com
      tam temiz, aystek-core-next'te sadece önceden var olan ilgisiz hatalar), `npx eslint` her iki
      tarafta da 0 uyarı/hata. Gerçek backend + gerçek dev veritabanına karşı canlı sunucu testi:
      `/public/sports/leagues/super-lig` ve `/public/sports/leagues` (liste) 200 döndü; admin
      `/admin/sports/leagues` kimliksiz istekte doğru şekilde 401 verdi (yeni route'lar permission
      middleware'ine doğru bağlanmış); `/spor` ve `/spor/super-lig` gerçek Next.js dev sunucusunda
      SSR ile render edildi ve gerçek takım/lig verisini (Galatasaray, Fenerbahçe, Süper Lig)
      doğru gösterdi; ana sayfa widget'ındaki linkler doğru hedefe (`/spor/super-lig`) işaret
      ediyor. Tarayıcı aracı bu ortamda mevcut değildi, bu yüzden doğrulama curl+gerçek sunucu
      SSR çıktısı üzerinden yapıldı (README'nin "tarayıcıda kullan" kuralına en yakın karşılık).
- [x] kureselhat.com aynı bileşenleri paylaştığı için ek kod gerektirmedi; **site-scoped içerik
      izolasyonu** korundu (`feedback_taxonomy_not_site_scoped` benzeri bir hata yapılmadı —
      yarışma verisi `site_id`'ye göre filtreleniyor, bkz. yukarıdaki doğrulama maddesi).

### Faz 3 — Otomatik senkronizasyon — ✅ Tamamlandı (2026-09-28)
- [x] §2.4'teki sağlayıcı değerlendirmesi yapıldı, kullanıcıya sunuldu, onay alındı: **API-Sports
      ailesi** ("API-Sports ailesiyle ilerle (Önerilen)" seçildi). Detay: §2.4.
- [x] `sports/service/sync.go`: `SyncConfig` (settings anahtarı `sports.sync`), `SyncProvider`
      arayüzü + registry (`map[string]SyncProvider`, şu an sadece `"apisports"`), `RunOnce(ctx)`
      — §2.3'teki `fx_sync.go` deseni uyarlandı.
  - **Tasarım notu — plandan bilinçli sapma:** orijinal plan "manuel girilmiş güncel veriyi asla
    otomatik veriyle ezme" için bir "X saattir güncellenmedi mi" sezgiselliği öneriyordu. Uygulamada
    bunun yerine daha net bir kural seçildi: `RunOnce` SADECE `ExternalRef` alanı set edilmiş VE
    `ExternalRef.Provider` aktif sağlayıcıyla eşleşen kayıtları senkronlar — `ExternalRef`'i admin
    panelden set etmek, o yarışma için "bundan sonra otomatik veriyle ez" onayının ta kendisidir.
    Saatlik sezgisellik yerine açık opt-in kullanmak hem daha öngörülebilir hem de admin `Upsert`
    yolunun her zaman `Source=manual` damgalamasıyla simetrik (iki yazma yolu asla yarışmaz).
- [x] Kimlik bilgisi şifreleme: emsal desenler (`fleet_tally/credcrypt`, `personal_assistant/
      credcrypto`) tekrar kullanılmadı ama BİREBİR AYNI DESENLE üçüncü bir kopya oluşturuldu —
      `internal/modules/sports/credcrypt/crypt.go` (AES-256-GCM, `SPORTS_CRED_KEY` env, `enc:`
      önekli saklama, anahtarsız ortamda pass-through modu).
- [x] `internal/platform/jobs/sports_sync_scheduler.go`: `route_fares/jobs/expiring_fare_scheduler.go`
      + `fx_sync_scheduler.go` desenlerinin birleşimi — periyodik, `ListDistinctTenantSiteScopes()`
      ile tüm tenant+site'ları gez, retry+backoff (`sportsSyncRetryMax=2`,
      `sportsSyncRetryBackoff=5s`), `internal/platform/ops/sports_sync_metrics.go` ile Prometheus
      metrikleri (`sports_sync_runs_total` vb., `fx_sync_metrics.go` ile aynı isimlendirme deseni).
- [x] `cmd/api/main.go`'da scheduler başlatıldı (diğer scheduler'ların yanına). Varsayılan kapalı:
      `sports.sync` settings kaydı hiç oluşturulmamışsa `LoadConfig` → `ErrSyncConfigMissing` →
      scheduler o tick'i atlar; operatör admin panelden açıp API anahtarını girene kadar hiçbir dış
      HTTP çağrısı yapılmaz.
- [x] Senkronize edilen veride `Source="auto"` + `SourceSyncedAt` işaretleniyor; admin panelde
      (`app/(dashboard)/admin/sports/page.tsx` liste tablosu) "Otomatik/Elle Girildi" rozeti zaten
      Faz 1/2'den beri gösteriliyordu — Faz 3 bu rozetin artık gerçekten `auto` değeri üretebileceği
      anlamına geliyor (şeffaflık ilkesi, TASK_060'daki "no fabricated fallback" felsefesiyle uyumlu).
- [x] Test: `internal/modules/sports/service/sync_test.go` (LoadConfig eksik/devre-dışı/desteklenmeyen
      sağlayıcı hataları, `ExternalRef` olmayan/farklı-sağlayıcılı kayıtların atlandığı, başarılı
      senkronun `Source=auto`+`SourceSyncedAt` damgaladığı, provider hatasının o kaydı "updated:false"
      olarak raporlayıp upsert etmediği) + `internal/modules/sports/provider/apisports_test.go`
      (`httptest.Server` ile gerçek HTTP+parse akışı, futbol response şekli, beklenmeyen şekilde
      `ErrAPISportsResponseShape`, boş-ama-geçerli fixtures listesi hata sayılmıyor).
- [x] Admin UI: `[code]/page.tsx`'e her yarışma için "Otomatik Senkronizasyon" bölümü (opt-in
      checkbox + provider/league_id/season alanları) + `page.tsx` liste sayfasına tenant-geneli
      senkron ayarları paneli (açık/kapalı, sağlayıcı, aralık, API anahtarı — anahtar asla geri
      okunmaz, sadece "configured" bayrağı döner, boş bırakılırsa mevcut anahtar korunur) eklendi.
- [x] 19 yeni `sports.external_sync.*` çeviri anahtarı `allLangs(tr,en)` deseniyle eklendi (bu alt
      önek `fullyTranslatedPrefixes`'te değil, çoğu diğer `sports.*` admin-UI etiketiyle aynı desen).
- [ ] **AÇIK:** Gerçek bir api-sports.io API anahtarıyla uçtan uca canlı doğrulama yapılmadı —
      özellikle basketbol/voleybol `/standings` ve `/games` response şekli sadece dokümantasyona
      dayanarak (toleranslı parse ile) uygulandı, canlı veriyle teyit edilmedi. Kullanıcıdan anahtar
      istenmeli (README "When to Ask" kuralı — üçüncü taraf kimlik bilgisi).

### Faz 4 — i18n + Menü + İzin + View-Profile senkronizasyonu — ✅ Tamamlandı (2026-09-28)
- [x] Yeni UI metinleri (spor türü adları, tier etiketleri, cinsiyet, liste sayfası metinleri vb.)
      `translation_seed.go`'ya **28 dilde** eklendi — Faz 1'de yapıldı (24 anahtar), mevcut
      `sports.*`/`news.labels.*` anahtarları tekrar edilmedi.
  - [x] `TRANSLATION_REPORT_PREFIX=sports. go test ./internal/seeds -run TestTranslationCoverageReport -v`
        ile doğrulandı: 24 yeni anahtarın tamamı 0 boşluk. İlk taramada 16 çift ("Futbol"/
        "Basketbol"/"Voleybol"/"Kod"/"Geri"/"Erkekler" — Azerbaycan/Türkmen/Özbek/Boşnak/Hırvat/
        Lehçe'de Türkçe ile genuinely aynı kelime) same-as-Turkish ratchet'i tarafından
        yakalandı; sahte farklı çeviri uydurmak yerine `identicalToTurkishAllowed` deseniyle
        gerekçelendirildi (bkz. commit `5b891ce`). Kalan sports.* boşlukları (33 anahtar, 396
        boşluk) önceden var olan eski `allLangs`-seeded anahtarlar — bu planın kapsamı dışında,
        ayrı `I18N_TRANSLATION_DEBT_WORKPLAN_TR.md`'de izlenir.
- [x] `menu_seed.go`: `"sports"` anahtarı zaten vardı, yeni bir alt-menü eklenmedi (liste ve
      düzenleme aynı `/admin/sports` rota ağacı altında, ek menü kaydı gerektirmedi).
- [x] `view_profiles_seed.go`: `aysTekHaberVisibleMenuKeys`'te `"sports"` zaten mevcuttu,
      değişiklik gerekmedi. `TestAysTekHaberViewProfileAllowlist` dahil tam seeds test paketi
      yeşil (doğrulandı).
- [x] `permission_seed.go`: `"sports"` modülü zaten kayıtlı, yeni handler'lar mevcut `sports:read`/
      `sports:update` aksiyonlarını kullandı (`routes.go`), yeni aksiyon eklenmedi.

### Faz 5 — Test, Dokümantasyon, Kapanış — ✅ Tamamlandı (2026-09-28)
- [x] Tam test paketi: `go test ./...` (core) — yeşil. Frontend: değişen üç dosya
      (`admin/sports/*`, `admin/posts/*`, `SportsSnapshotPanel.tsx`, `NewsChrome.tsx`,
      `login/page.tsx`) için `npx tsc --noEmit` temiz (repo genelindeki tek hata önceden var olan,
      bu işle ilgisiz `publish_at` uyarısı).
- [x] `/workspace/aystek-core-dev-docs/done/TASK_061_SPORTS_AUTO_SYNC.md` (teknik referans) +
      `TASK_061_SUMMARY.md` (özet) oluşturuldu — SOP gereği, TASK_060'tan devam.
- [x] Bu dokümanın başındaki "Durum" satırı güncellendi.
- [x] Installer/operasyon dokümanları gözden geçirildi: Faz 3 yeni bir systemd timer
      gerektirmiyor (scheduler mevcut `aystek-core` sürecinin içinde bir goroutine).
      `SPORTS_CRED_KEY` zaten `.env.example`'a eklenmişti (aynı diğer `*_CRED_KEY`
      değişkenleri gibi tek başına `SUNUCU_KURULUM_REHBERI.md`'de ayrıca belgelenmiyor,
      mevcut konvansiyonla tutarlı — kontrol edildi).

### Faz 6 — Site izolasyonu + çoklu sağlayıcı + ücretsiz Web Feed — ✅ Tamamlandı (2026-09-29)

- [x] Kritik kapsam düzeltmesi: `sports.sync` ayarları tenant-genelinden
      `site.<site_id>.sports.sync` anahtarına taşındı. Eski tenant anahtarı salt geçiş fallback'i;
      yeni kayıtlar yalnızca site kapsamına yazılır.
- [x] Tek aktif sağlayıcı modeli geriye uyumlu çoklu sağlayıcı profillerine dönüştürüldü. Her
      yarışma kendi `external_ref.provider` değeriyle site profilini seçer; API-Sports,
      TheSportsDB ve Web Feed aynı site içinde eş zamanlı açık olabilir.
- [x] Sağlayıcı kataloğu backend çalışma zamanından sunuluyor
      (`GET /admin/sports/sync/providers`); frontend artık Go registry'sini tekrar eden sabit bir
      liste taşımıyor.
- [x] TheSportsDB'nin resmi ücretsiz anahtarı güncel dokümana göre `123` olarak düzeltildi.
- [x] `webfeed` sağlayıcısı eklendi: canonical JSON ve tek-tab CSV okur; Google E-Tablolar'ın
      web'de yayınlanan CSV çıktısıyla uyumludur.
- [x] Web Feed güvenlik sınırı: HTTPS, URL kimlik bilgisini reddetme, private/loopback/link-local/
      multicast IP engeli, DNS çözümünü dial sırasında yeniden doğrulama, yönlendirme denetimi,
      ortam proxy'sini kapatma, 15 saniye timeout ve 2 MiB açılmış gövde sınırı.
- [x] Kaynak şeffaflığı: `ExternalRef` içine `feed_url`, `source_url`, `attribution` eklendi;
      atıf Web Feed için zorunlu ve public spor panelinde gösteriliyor.
- [x] Siteye özgü zamanlayıcı: sabit tenant-geneli aralık yerine her tenant+site'ın kendi
      `interval_minutes` değeri ayrı takip ediliyor.
- [x] Senkron sonuçları `provider` ve güvenli durum kodu (`updated`, `fetch_failed`,
      `validation_failed`, `persist_failed`, `provider_unsupported`) döndürüyor; gizli anahtar veya
      upstream yanıtı admin cevabına sızdırılmıyor.
- [x] Zamanlayıcı kısmi sağlayıcı hatasını artık başarı saymıyor; geçici fetch/persist hatalarını
      sınırlı retry ile yeniden deniyor, doğrulama hatasını gereksiz yere tekrar çağırmıyor.
- [x] `SPORTS_CRED_KEY` eksikken şifreli anahtarın upstream'e cipher metni olarak gönderilmesi
      engellendi; açık yapılandırma hatası dönüyor. Public API ise feed URL/lig kimliği/sezon gibi
      dahili locator alanlarını redakte edip yalnızca sağlayıcı ve kaynak atfını yayımlıyor.
- [x] Yönetim ekranı sağlayıcı profili kartları, yarışma bazlı dinamik sağlayıcı alanları ve Web
      Feed URL/kaynak/atıf alanlarıyla güncellendi.
- [x] 9 yeni yönetim metni 28 dilde gerçek çeviriyle eklendi; yeni anahtarlar çeviri coverage
      ratchet'ine dahil edildi. Önceden var olan 21 `sports.external_sync.*` allLangs kaydı ayrı
      i18n borç planında kalmaya devam ediyor; bu faz yeni borç üretmedi.
- [x] Birim testleri: JSON parse, CSV parse, güvensiz URL reddi, canonical feed HTTP akışı ve
      site anahtarı önceliği. Ayrıntılı kullanıcı/teknik rehber ve örnek CSV:
      `SPORTS_WEB_FEED_GUIDE_TR.md`, `templates/sports_web_feed_template.csv`.

**Bilinçli olarak sonraki faza bırakılan geliştirmeler:** kaydetmeden feed önizleme; aynı yarışma
için birincil/yedek kaynak sıralaması; otomatik değişiklik onay kuyruğu; sağlayıcı kota/sağlık
paneli. Bunlar mevcut yayın akışını bloklamaz ve rehberin §7 bölümünde devralma sırasıyla yazılıdır.

**Kapanışta bulunan ek düzeltmeler** (sports testi sırasında ortaya çıkan, modülden bağımsız
gerçek hatalar — detay TASK_061 dokümanında): editor rolüne `sports` izni eksikti; view-profile
menü görünürlük listesi (`allMenuKeys`) 76 anahtar geride kalmıştı ve artık `menuSeeds`'ten
türetiliyor; etiketler (`tags`) site bazlı izole edilmedi (tenant genelinde sızıyordu); zamanlayıcı
hiç çalışmıyordu (tenant'sız context ile ayar okuyup "kapalı" varsayıyordu — `StartFXSyncScheduler`
ile aynı "fail-open" desene çekildi); kureselhat.com üretimde iki özellik commit'i eksikti (tek
dosya yaması yapılmış, tam yeniden derleme yapılmamıştı); kapak görseli önizlemeleri ve
galibiyet/beraberlik/mağlubiyet harfleri (W/D/L vs G/B/M) düzeltildi.

---

## 5. Devam Notu (bu planı devralan bir sonraki AI için)

1. Önce §1'i (Mevcut Durum) oku — modülün zaten var olduğunu ve nerede olduğunu unutma, sıfırdan
   kurmaya kalkma.
2. `git log --oneline -- internal/modules/sports/ aystekhaber-com/components/news/SportsSnapshotPanel.tsx aystekhaber-com/lib/api/sports.ts` ile bu plandan sonra ne değiştiğini kontrol et — bu
   doküman yazıldıktan sonra başka bir AI ilerlemiş olabilir, "Durum" satırlarına güven ama
   doğrula.
3. Hangi fazda kalındıysa oradan devam et; bir faz bitmeden sıradakine geçme (bağımlılıklar §4'te
   not edildi — Faz 2 ve Faz 3, Faz 1'e bağımlı ama birbirlerinden bağımsız paralel ilerleyebilir).
4. Herhangi bir belirsizlikte (§2.4'teki API seçimi gibi) kullanıcıya sor — README'nin "When to
   Ask" kuralı gereği, özellikle üçüncü taraf API anahtarı gerektiren kararlarda sessizce
   ilerleme.
5. **Bilinen tuzak — `Normalize()` sadece tekil okumayı, filtreli listeyi KORUMAZ (2026-09-28'de
   üretimde bulundu ve düzeltildi):** `LeagueSnapshot.Normalize()` boş `Sport`/`Tier` alanlarını
   Go tarafında `football`/`top` varsayılanına çeker, ama bu SADECE dokümanı okuduktan SONRA
   uygulanır. `repository.list()`'teki `criteria["tier"]=filter.Tier` filtresi MongoDB sorgusunu
   HAM (normalize edilmemiş) alan üzerinden çalıştırır — yani Faz 1 öncesinden kalma, `tier`/`sport`
   alanı hiç yazılmamış bir kayıt, filtresiz `GET /public/sports/leagues` çağrısında görünür
   (normalize sonrası "top" gibi görünür) ama `?tier=top` filtreli çağrıda (anasayfa widget'ının
   asıl kullandığı çağrı budur) SESSİZCE kaybolur — hata yok, sadece boş liste. Üretimde tam olarak
   bu yüzden anasayfadaki spor widget'ı "hiç değişmemiş" görünüyordu; kök neden production'daki tek
   `super-lig` kaydının Faz 1 öncesi seed'den kalma olup `tier`/`sport` alanlarının hiç set
   edilmemiş olmasıydı — `db.sports_league_snapshots.updateOne(...,{$set:{tier:"top",
   sport:"football",...}})` ile düzeltildi. Yeni/manuel eklenen kayıtlarda bu sorun oluşmaz
   (admin `Upsert` yazmadan önce `Normalize()` çağırıyor), sadece eski/elle-Mongo'ya-yazılmış
   kayıtlarda görülür — ama bir sonraki AI benzer bir "veri var ama filtreli listede yok" belirtisi
   görürse önce bu deseni kontrol etsin.
