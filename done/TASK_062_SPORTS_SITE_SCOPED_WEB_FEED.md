# TASK 062 — Site Kapsamlı Spor Sağlayıcıları ve Ücretsiz Web Feed

**Tarih:** 2026-09-29 · **Durum:** Kod ve otomatik testler tamamlandı; uzak sunucu dağıtımı bu
görevin dışında bırakıldı.

## Sonuç

Spor Merkezi'nin daha önce tenant-geneli ve tek sağlayıcılı olan otomatik senkronizasyonu,
tenant altındaki her site için ayrı çalışan çoklu sağlayıcı profillerine dönüştürüldü. Ücretli
API olmadan yayın yapabilmek için güvenli JSON/CSV Web Feed adaptörü ve Google E-Tablo kullanım
akışı eklendi.

## Teslimler

- Site ayar anahtarı: `site.<site_id>.sports.sync`.
- Geriye uyumluluk: site kaydı yoksa eski `sports.sync` yalnızca okunur; ilk kayıt site anahtarına
  yazılır.
- Eş zamanlı sağlayıcı profilleri: API-Sports, TheSportsDB, Web Feed / Google Sheets.
- Backend sağlayıcı kataloğu; frontend registry tekrarı kaldırıldı.
- Her tenant+site için bağımsız zamanlayıcı aralığı.
- Canonical JSON ve CSV ayrıştırma; indirilebilir UTF-8 CSV şablonu.
- Web Feed SSRF koruması, timeout, içerik boyutu ve redirect sınırı.
- Kaynak/atıf alanları ve public görünümde atıf bağlantısı.
- Public API'de dahili `feed_url`, sağlayıcı lig kimliği ve sezon redaksiyonu.
- TheSportsDB resmi ücretsiz anahtarı `123`; geçmiş maç yanıtında sezon uyuşmazlığı koruması.
- Kısmi sync hatalarının doğru metriklenmesi ve yalnızca geçici hatalarda retry.
- Şifre çözme anahtarı eksikken ciphertext'in upstream'e gönderilmesini engelleyen fail-closed
  davranış.
- 9 yeni UI anahtarı için 28 gerçek çeviri ve coverage ratchet.

## Değişen projeler

### aystek-core

- `internal/modules/sports/model/league.go`
- `internal/modules/sports/service/sync.go`
- `internal/modules/sports/handler/{sync_admin,public}.go`
- `internal/modules/sports/provider/{apisports,thesportsdb,webfeed}.go`
- `internal/platform/jobs/sports_sync_scheduler.go`
- `internal/seeds/{translation_seed,translation_coverage_test}.go`
- `cmd/api/main.go`

### aystek-core-next

- `app/(dashboard)/admin/sports/page.tsx`
- `app/(dashboard)/admin/sports/[code]/page.tsx`
- `lib/api/sports.ts`

### aystekhaber-com ortak public kod tabanı

- `components/news/SportsSnapshotPanel.tsx`
- `lib/api/sports.ts`

## Doğrulama

- `go test ./internal/modules/sports/... ./internal/platform/jobs/...`
- `go test ./cmd/api`
- `TestFullyTranslatedPrefixesHaveNoGaps` ve yeni sports/webfeed testleri.
- Core Next değişen dosyalarda ESLint temiz; repo genelindeki önceden var olan TypeScript
  borçlarından sports dosyalarına ait yeni hata yok.
- Ortak public site ESLint ve `tsc --noEmit` temiz.
- `git diff --check` temiz.

## Bilinçli sonraki adımlar

1. Feed'i kaydetmeden önce önizleme ve satır hata raporu.
2. Tek yarışmada öncelikli/yedek kaynak zinciri.
3. Büyük otomatik değişiklikler için editör onay kuyruğu.
4. Kota, gecikme, son başarılı çalışma ve sağlayıcı sağlık paneli.
5. Üretime dağıtım sonrası Kuresel Hat sitesinde gerçek Google E-Tablo ile smoke test.

Kullanıcı rehberi ve sözleşme: `../SPORTS_WEB_FEED_GUIDE_TR.md`.
