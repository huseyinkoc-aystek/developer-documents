# AYPASS Mevcut Uygulama Tamamlama — P00 Keşif Raporu

**Tarih:** 2026-10-29 · **Oturum türü:** yalnızca okuma/keşif · **Hazırlayan:** Claude (Sonnet 5.5) · **Durum:** Hüseyin Koç P00 özetini ve inceleme sınırlarını kabul etti (WP-1 onayı ile birlikte). Bu rapor **"kapsam eksikleri kayıtlı keşif raporu"**dur: tamamı doğrulanmış veya kapanmış sayılmaz; mali altyapının doğrulandığı ya da uygulamanın yayına hazır olduğu anlamına gelmez. Açık keşif işleri §8'dedir.

> **P00 oturumunda yapılmayanlar:** Üretim uygulama kodu değiştirilmedi; commit, push, PR, deploy, migration, seed, canlı DB/POS/banka/sunucu işlemi, gerçek tahsilat/iade/fatura yapılmadı. P00'ın tek çıktısı bu rapordur. `aystek-core-dev-docs` bir git deposu değildir; bu rapor sürüm kontrolü altında değildir.
>
> **Sonraki adımlar:** WP-1 (panel, iki turda tamamlandı — ilk tur + kullanıcı incelemesinin istediği sertleştirme) → `AYPASS_WP1_ACCESSPOINT_SOZLESME_RAPORU_TR.md`. WP-2 (backend, cihaz–erişim noktası bağ denetimi, ayrı dal `fix/aypass-device-accesspoint-binding` `bb1ec78`, taban `3a3bb86`) **uygulandı, yerelde, push edilmedi** — bkz. bu raporun güncellenen WP-2 satırı ve commit mesajı. WP-2.1 (D-22/D-23 takip paketi, aynı dal, ek commit `51330f5`) **uygulandı, yerelde**. WP-15 (backend güncelleme sözleşmesi, KR-13) için önce karar kaydı yazılmıştı; **2026-10-30 güncellemesi: WP-15 artık kod + hedefli doğrulamayla TAMAMLANDI (merge/yayın kararı bekliyor)** — backend dalı `feat/aypass-accesspoint-patch` (son commit `10ad86e`, taban `51330f5`), panel dalı `feat/aypass-accesspoint-patch` (`1f1e60c`, taban `a9752d6`), ayrıca iki Flutter istemcisinin düzeltmesi: `aypass_erisim_app` (`fix/aypass-erisim-devicebinding-patch` `3fa8d29`) ve `kt_surucu_app` (`fix/kt-surucu-devicebinding-patch` `8bdcfb2`). Hiçbiri `main`'e birleşmedi/push edilmedi. P00'ın bulguları ilk yazıldığı gibidir; D-01/D-02/D-15/D-22/D-23'ün güncel durumu §4.1'de; WP-15'in kapsamına giren yeni bulgular (Create'te cihaz-bağlama izin açığı, Flutter istemci sözleşme uyumsuzluğu) D-24/D-25 olarak eklendi.

> **Kanıt etiketleri** (Ana Prompt'tan): `KODDA_GORULDU` · `BELGEDE_BILDIRILDI` · `HEDEFLE_FARK` · `RISK_TEST_GEREKLI` · `DOGRULANMADI`.
> **Kaynak ayrımı:** eski uygulama belgeleri = eski davranış · Rev.06 + kullanıcı istekleri = hedef davranış · mevcut kod = fiili durum · test sonucu ve canlı doğrulama = ayrı kanıt. Eski uygulamanın kusurları (AYG-xx) yeni uygulamada var sayılmamıştır; yeni uygulamadaki karşılığı §4.4'te ayrıca işaretlidir.
> **Satır numaraları** aksi belirtilmedikçe backend `3a3bb86` (dal `aypass-approval-ui`) içindir.

---

## 0. Yönetici özeti

1. **Temel sahipliğin çoğu doğru ve yeniden kullanılabilir:** Person → Wallet → Ledger zinciri hedefle uyumlu. Bakiye kartta değil cüzdanda; kart/QR yalnızca erişim aracı; para çıkışı koşullu `$inc` ile atomik rezerve ediliyor; ledger'da idempotency için kısmi unique indeks var.
2. **Hedefle en büyük fark ücret motorunda:** A1–A4, binişte azami tutarın fiilen tahsili ve inişte nihai ücret üzerinden fark iadesi kodda **yok** (A3/A4 tanımı için bkz. §2 not 1: yalnızca bloke/provizyon değildir). Mevcut `metered` akışı girişte ücret almıyor, çıkışta mesafeden hesaplıyor (`service.go:1472`). A4 için eksik durak-çifti hücresi **sessizce km ücretine düşüyor** (`service.go:1539`, `:1760`).
3. **Olay zamanı sözleşmesi yok:** `RecordEvent` cihaz zamanı (`occurred_at`) almıyor; tarife/indirim çözümlemesi `time.Now()` ile yapılıyor; `AccessEvent` uygulanan tarife, sağlayıcı/firma veya kart kimliğini saklamıyor.
4. **Cihaz tarafı en zayıf halka:** `aypass_erisim_app`'te kalıcı offline kuyruk, idempotency anahtarı ve HTTP timeout yok; `aypass_kullanici_app`'in `kBaseUrl` değeri `http://localhost:9090`. Gateway'de `:id` ile gelen olay uçları, çağıran cihazın o erişim noktasına bağlı olup olmadığını denetlemiyor.
5. **Panel ↔ backend sözleşmesinde gerçek veri kaybı var:** Panelden her AccessPoint düzenlemesi mevcut **araç plakasını siliyor** ve `pos_enabled`'ı `false` yapıyor. Plaka olmadan `fare_transport` olayı kaydedilemiyor (`ErrTransportVehicleRequired`).
6. **Devir (Model 1) ve Global Market (Model 4) için kodda sahip yok.** `service_assignments`, pocket, merchant/program/fon, settlement: sıfır sonuç. En yakın zemin: `TripContract` (geçerlilik aralıklı) ve `trip_billing`.
7. **Kanıt sınırı (2026-10-30 güncellemesi):** P00 oturumunda bu sandbox'ta MongoDB yoktu; cüzdan/ledger eşzamanlılık testlerinin 12'si o oturumda atlanmıştı. **WP-15 turunda** bu sandbox'ta yerel, sentetik-only, geçici bir gerçek MongoDB kuruldu ve **belirlenen 12 cüzdan/ledger eşzamanlılık senaryosu bu gerçek MongoDB üzerinde çalıştırıldı ve başarısızlık gözlenmedi** (SHA `703f05d`, bkz. §1.4 ve §8.1 O-7). Bundan "tüm mali çekirdeğin doğrulandığı", "A1–A4 algoritmalarının tamamlandığı" veya "offline tahsilatın güvenli olduğu" sonucu **çıkarılamaz** — bu 12 test mevcut bakiye/journal/idempotency davranışını kapsar, A1–A4 ücret motorunu kapsamaz. Dağıtılan sürüm hâlâ doğrulanmadı.
8. **Önerilen ilk küçük iş paketi:** WP-1 — AccessPoint panel↔backend sözleşme düzeltmesi + sözleşme testi (§7).

---

## 1. Kaynak ve sürüm manifesti

### 1.1 Depolar

| Depo | Yerel yol | Remote | Dal / HEAD | Notlar |
|---|---|---|---|---|
| Backend (Go) | `/workspace/go/aystek-core` | `git@github.com:huseyinkoc-aystek/aystek-core.git` | **`aypass-approval-ui` `3a3bb86`** | `main` = `f48748f`; yerel dal main'in **44 commit** ilerisinde. `aypass-rev06-faz0` (`a712aca`) HEAD'in atasıdır. Kirli: `.env` (**izlenen** dosya, `M`, +2/−1; içerik raporlanmadı, stage edilmedi). |
| Panel (Next.js) | `/workspace/node/aystek-core-next` | `git@github.com:huseyinkoc-aystek/admin-panel-next.git` | **`main` `c5e0994`** (origin/main ile aynı) | `package.json name` = `aystek-core-next` (depo adı `admin-panel-next`). `origin/approval-policy-ui` = `7b03fe7` (main'de **yok**: main'de 1, o dalda 5 ek commit). Yerel `approval-policy-ui` = `434617d`, origin'in **1 commit önünde** (haber paneli işi; kullanıcıya ait, dokunulmadı). Kirli: `tsconfig.tsbuildinfo` (M), `incognito_check.mjs` (??). Ek worktree: oturum scratchpad'inde `wt-panel` (`aypass-faz0-plate-reason` `7b03fe7`). |
| AYPASS cihaz uygulaması | `/workspace/flutter/aypass_erisim_app` | `git@github.com:huseyinkoc-aystek/aypass-erisim-app.git` | `main` **`ec6e999`** (2026-09-29) | `1.0.0+1`. Kirli: `android/gradle.properties`. |
| AYPASS kullanıcı uygulaması | `/workspace/flutter/aypass_kullanici_app` | — | **git deposu değil → SHA yok** | `1.0.0+1`. Sürüm izlenemiyor (bkz. R-10). |
| İkincil (kullanıcı yönlendirmesiyle AYPASS cihaz kapsamı dışı) | `kt_surucu_app` `b80a230`, `yolcu_app` `e9fcfc8`, `sofor_app` `0eb7784` | — | — | Yalnızca karşılaştırma notu için tarandı. |
| Installer | — | — | `main` `7d45fb6` (önceki oturumdan; bu oturumda **yeniden doğrulanmadı**) | Önceki oturumda bekleyen doküman düzenlemeleri vardı. |
| Eski PHP (AYShuttleGursel) | — | — | `0956c9b` + 4 yerel değişiklik (**Yol Haritası'ndan alıntı; yeniden doğrulanmadı**) | Yalnızca belge üzerinden okundu. |

**Sürüm farkı uyarısı:** Yol Haritası'ndaki E01–E20 kanıtları backend `f48748f` / panel `56907eb` üzerinde statik türetilmiştir. Bu rapor `3a3bb86`'yı okur. Aradaki 44 commit (Faz 0 sertleştirmeleri) bazı bulguları **eskitir** (§4.3). Panel `main` ile `approval-policy-ui` aynı anda "mevcut" kabul edilemez.

### 1.2 Dağıtım

- Prompt'taki adres: `aypasss.aysteknolojileri.com` (üç "s"). Yalnızca üç AI prompt dosyasında geçiyor. Kodda/seed'de (`internal/seeds/gursel_turizm_seed.go`) geçen adres `aypass.gurseltur.com.tr`. **Hiçbiri değiştirilmedi**; hangisinin canlı olduğu DOĞRULANMADI. Bu oturumda DNS/HTTP denenmedi; önceki oturumda sandbox'tan `aypass.gurseltur.com.tr` çözülememişti.
- Depodaki dağıtım yüzeyleri (içerik incelenmedi): `deployments/Dockerfile`, `deployments/production/{README.md,docker-compose.yml}`, `deployments/systemd`, `deployments/grafana`, `deployments/prometheus/alerts.yml`, `.github/workflows/{ci,deploy,release,security-scan}.yml`, `Makefile`.
- **Canlıda çalışan sürüm: DOĞRULANMADI.** GitHub `main`, yerel çalışma ağacı ve dağıtılan sürüm aynı varsayılmadı.

### 1.3 Okunan belgeler

| Belge | Durum |
|---|---|
| `README_FOR_AI_ASSISTANTS.md` | Var. Prompt'taki zaman damgalı ad (`…(20261029-120719).md`) **bulunamadı**; ad farkı kaydedildi, içerik uydurulmadı. |
| `AYSTEK_CORE_CANONICAL_ARCHITECTURE.md` | Tam okundu. |
| `GURSEL_YENI_NESIL_DONUSUM_MASTER_TR.md` | Durum bölümü okundu (önceki oturum). |
| `AYPASS_GECIS_SISTEMI_REV06_ANALIZ_VE_IS_PLANI_TR.md`, `AYPASS_REV06_EXCEL_ESLEME_TR.md` | Okundu (önceki oturum); satır numaraları `f48748f`'e göre, **eskimiş** (§4.3). |
| `AYSHUTTLEGURSEL_PROJE_ANALIZI_TR.md` | §6.1–6.6 ve §10 bu oturumda okundu. §1–5, 7–9, 11–14 **okunmadı**. |
| `AYSHUTTLEGURSEL_TEKNIK_ENVANTER_TR.md`, `AYSHUTTLEGURSEL_CURRENT_STATE.md` | **Bu oturumda okunmadı** (CURRENT_STATE zaten yalnızca tarihsel). |
| `AYPASS_AI_Ana_Prompt_TR.md`, `AYPASS_AI_Yol_Haritasi_TR.md`, `AYPASS_AI_Tum_Promptlar_TR.md` | Ana Prompt tam; Yol Haritası E01–E20 + REQ tablosu; Tüm Promptlar yalnızca P00 bölümü. |
| Ayrı "kanıt analizi" dosyası | **Yok**; kanıt analizi Yol Haritası içindeki E01–E20 bölümüdür. |

### 1.4 Çalıştırılan doğrulama

| Komut (özet) | Sonuç |
|---|---|
| `go test ./internal/modules/{access_points,wallets,ledger,fleet_tally/rider_trips}/... -count=1 -v` (`AYSTEK_TEST_MONGO_URI` tanımsız) | **137 PASS · 12 SKIP · 0 FAIL.** Atlanan 12 test, gerçek Mongo replica set isteyen testlerdir: `TestDebit_ConcurrentNeverCrossFloor`, `TestDebit_ConcurrentSameKeyAppliesOnce`, `TestDebit_JournalPostFailureLeavesCacheUnchanged`, `TestReconcile_ConcurrentMixedCreditDebitMatchesLedger`, `TestConfirmBankTransfer_ConcurrentConfirmCreditsOnce`, `TestApplyBalanceChange_*` (2), `TestEnsureSystemAccount_*` (3), `TestConfirmPendingJournal_ConcurrentConfirmations`, `TestVerifyModule_LedgerJournalIDIndexIsRegistered`. |
| Mongo erişimi | P00 oturumunda `127.0.0.1:27099` erişilemez, `mongod` yok → eşzamanlılık/atomiklik testleri o oturumda doğrulanmadı. **2026-10-30 güncellemesi (WP-15 turu):** bu sandbox'ta yerel, disposable, sentetik-only bir gerçek MongoDB (7.0, replica set `rs0`, port 27199) kuruldu; `AYSTEK_TEST_MONGO_URI` ile hem WP-15'in kendi gerçek-Mongo testleri (9 `TestMongoIT_*` + 11 `TestRealAuthChain_*`) hem de önceden atlanan **12 cüzdan/ledger eşzamanlılık testinin 12'si de çalıştırıldı — 12/12, başarısızlık gözlenmedi** (SHA `703f05d`). Oturum sonunda MongoDB süreci güvenle durduruldu (yalnız bu oturuma ait PID/dbpath doğrulanarak). |
| Panel / Flutter testleri | **Çalıştırılmadı.** |

---

## 2. Gereksinim → backend → panel/cihaz → test → durum matrisi

Durum sözlüğü: **Uyumlu** (kodda görüldü, hedefle uyumlu) · **Kısmi** · **Fark** (hedefle ayrışıyor) · **Yok** · **Zincir çıkarılmadı** (P00'da route→handler→service→repo→UI→test zinciri çıkarılmadı).

> **Not 1 — A3/A4 gereksinim kaydı (Hüseyin Koç düzeltmesi, 2026-10-29).** A3/A4 yalnızca bloke/provizyon **değildir**:
> 1. **Binişte** maksimum tutar fiilen **tahsil edilir** (P).
> 2. **İnişte** nihai ücret (N) mesafe (A3) ya da yönlü sabit matris (A4) üzerinden hesaplanır; varsa indirim uygulanır.
> 3. Fark (P − N) **iade edilir**.
>
> Bu raporun önceki taslağındaki "azami tutar çekimi (hold)" ifadesi bu anlamda okunmalıdır; "hold" kelimesi yalnızca mevcut kodun `access_point.go:65-67` ve `service.go:1468-1471` yorumlarındaki ifadeyi anar (o yorumlar bloke/provizyon fikrini anlatır, hedef bu değildir). Bu kayıt WP-5 ve WP-6'nın gereksinim metnidir; WP-1 kapsamında mali kod geliştirilmedi.
>
> **Not 2 — Tahsilat ayrımı.** Banka/POS tahsilatı (dış sistem, onay bekleyebilir: `PENDING_EXTERNAL`) ile iç cüzdan tüketimi (ledger'da `Debit`) farklı olaylardır ve birbirine karıştırılmaz. A3/A4'te biniş tahsilatı ve fark iadesi **iç cüzdan** hareketleridir; cüzdanın bakiye yüklemesi ise ödeme onayından sonra ayrı bir olaydır (REQ-19, REQ-21).

| REQ | Konu | Backend (sahip, `dosya:satır`) | Panel / cihaz | Test | Durum · etiket |
|---|---|---|---|---|---|
| 01 | Model 1/2/3 kurum–operatör–sağlayıcı ayrımı | `institutions/model/institution.go:84` (`Type` yalnız rol etiketi :16-20); `AccessPoint.InstitutionID` `access_point.go:272`; `Vehicle.SupplierID` `stock/vehicles/model/model.go:37`; `TripContract` `fleet_tally/trip_contracts/model/trip_contract.go:22-28` | DOĞRULANMADI | — | **Fark** · HEDEFLE_FARK (sağlayıcı sahibi tek yerde değil; `service_assignments` yok) |
| 02 | Person hesabı; kart değişiminde bakiye devri yok | `Person` `persons/model/person.go:51`; `Wallet.PersonID` + "kişi×para birimi" unique `wallets/repository/indexes.go:15`; kart/QR `WalletID` taşır `wallet_instruments/model/card.go:34`, `qr.go:18`; kartta bakiye alanı yok; `ReplacedByInstrumentID` `card.go:41` | Panel wallets/instruments parser'ları eşleşiyor (alt ajan) | `wallets/service/service_test.go` (fake), `wallet_instruments` 1 test dosyası | **Uyumlu** · KODDA_GORULDU |
| 03 | Kayıp kart / geçmiş sahiplik; gecikmiş olay doğru hesaba | `AccessEvent` kart/enstrüman **kimliği saklamıyor**, yalnız `InstrumentUsed` türü (`access_points/model/access_event.go`); `CreatedAt` sunucu zamanı `service.go:1827-1850` | — | — | **Fark** · HEDEFLE_FARK + RISK_TEST_GEREKLI |
| 04 | Tenant/site ile kurum/provider yetkisi ayrı | `tenantfilter` repo seviyesinde: `access_points/repository/mongo.go:236-243`, `:291`; kurum erişimi `institutionfilter` | — | önceki oturum denetim testleri | tenant/site **Uyumlu**; provider yetkisi DOĞRULANMADI |
| 05 | QR/kart banka yetkisinden ayrı; karta yazma yok | QR yalnız hash tutar `qr.go:20` (`json:"-"`), `ExpiresAt`/`ConsumedAt` `:21,:26`; kartta yalnız `CardNumber` `card.go:36`; backend'de kart-yazma fonksiyonu yok | `aypass_erisim_app/lib/screens/nfc_read_screen.dart:50-92` yalnız okur | — | **Uyumlu** · KODDA_GORULDU (K06 yalnız Flutter'da) |
| 06 | Sabit ID; ad/plaka değişiminde geçmiş korunur | `Vehicle.PlateHistory` `model.go:53` + zorunlu neden (`93b4232`, atomik `09a426b`); `AccessEvent.VehiclePlate` olay anı anlık görüntüsü | Panel edit AccessPoint plakasını **siliyor** (D-01) | `stock/vehicles` testleri | **Kısmi** · KODDA_GORULDU |
| 07 | Araç ataması olmadan doğrulanmış sefer/tarife; belirsiz fiyata izin yok | `resolveEventContext` `service.go:440-490`: plaka `ap.VehiclePlate`'ten → `ActiveByVehiclePlate` → `ScheduledRouteID`; plaka boşsa `ErrTransportVehicleRequired` | Panel create formunda plaka alanı yok; edit'te siliniyor (D-01) | `service_test.go` (fake) | **Fark** · HEDEFLE_FARK (tarife kararı plakaya bağlı) |
| 08 | F0: sabit, indirimli, ücretsiz | `PricingModeFree/FixedFee` `access_point.go:51-56`; `resolveFee` `:1161`; `chargeWallet` `:1180`; `access_point_discounts` modülü | Panel fee=0 gönderince backend yok sayar `service.go:710` (D-03) | `service_test.go:2179`, `:3031` (`TestChargeCashlessNeverAppliesDiscount`) | **Uyumlu** (FX1 formülü `:1177`) · KODDA_GORULDU |
| 09 | A1 bilinen kampüs/fabrika (sabah biniş/akşam iniş mesafe) | Algoritma kavramı yok; en yakın `chargeMetered` `:1472` | — | — | **Yok** · HEDEFLE_FARK |
| 10 | A2 güzergâh/yön/durak sabit ücret | `RouteFare.CardFee/CashFee` `route_fares/model/route_fare.go:31-32`; `resolveBaseFee` `:1094` | Panel `fleet-tally.ts:217-229` eşleşiyor | `route_fares` testleri (kapsam taranmadı) | **Kısmi** · algoritma adı/seçimi yok |
| 11 | A3: binişte maksimum tutar fiilen tahsil; inişte mesafeden nihai ücret (+indirim); fark iadesi (bkz. Not 1) | `metered`: giriş ücretsiz, binişte maksimum tutar tahsil edilmiyor (`access_point.go:65-67`, `service.go:1468-1471` "deliberately later phase"); çıkışta `distance×oran` `:1546-1548`; fark iadesi yok | — | — | **Fark** · HEDEFLE_FARK |
| 12 | A4: binişte maksimum tutar fiilen tahsil; inişte yönlü sabit matristen nihai ücret (+indirim); fark iadesi (bkz. Not 1) | `RouteStopFare` `route_stop_fares/model/route_stop_fare.go:21-45` (yönlü çift); `ResolveActiveStopFare` hata/eksik hücrede **mesafeye düşer** `service.go:1539`, `:1760` | Panel parser eşleşiyor | `TestReverseDirectionDoesNotCollide` `route_stop_fares/service/service_test.go:246` (**çalıştırılmadı**) | **Fark** · HEDEFLE_FARK (sessiz fallback yasak) |
| 13 | Algoritma/indirim/risk/tahsilat ayarları güzergâh bazlı | Tarifeler `ScheduledRouteID` bazlı; algoritma/risk/offline ayar alanı yok | — | — | **Kısmi** · HEDEFLE_FARK |
| 14 | Tarihli tarife sürümü; eski olaya yeni fiyat yok | `EffectiveFrom/To` var (`route_fare.go:49-50`); ama çözümleme `time.Now()`: `service.go:1101`, `:1131`, `:1169`, `:1539`, `:1760`; olayda tarife kimliği/sürümü yok | — | — | **Fark** · HEDEFLE_FARK + RISK_TEST_GEREKLI |
| 15 | Online alt bakiye ile offline risk/borç ayrımı | `BalanceFloorOverride`, `FloorExempt` `wallet.go:92,97`; offline risk modeli yok | — | — | **Kısmi** · HEDEFLE_FARK |
| 16 | Kalıcı cihaz kuyruğu, ACK, duplicate/out-of-order güvenliği | Sunucuda `occurred_at`/sıra alanı yok (`RecordEventInput` `service.go:131-149`); event idempotency var (aşağıda) | `aypass_erisim_app`: kuyruk **yok**, idempotency anahtarı **yok**, timeout **yok** (`lib/services/validation_service.dart:183`) | — | **Yok** · HEDEFLE_FARK |
| 17 | Tek mali kaynak; posted hareket değişmez, reversal ile düzeltme | `ledger` `PostJournal`/`ConfirmPendingJournal`; cache hep koşullu `$inc` (`wallets/service/service.go:974-1002`); `RefundEvent` `access_points/service/service.go:324` | — | `balance_concurrency_test.go:192` (**atlandı**) | **Uyumlu** · KODDA_GORULDU + RISK_TEST_GEREKLI (bu oturumda koşulmadı) |
| 18 | Eşzamanlı bakiye ve idempotency payload bütünlüğü | `ApplyBalanceChange` `wallets/repository/mongo.go:123`; `debitToAccount` `wallets/service/service.go:323-359`; ledger kısmi unique `ledger/repository/indexes.go:14-25`; event kısmi unique `access_points/repository/indexes.go:59-66` | — | `balance_concurrency_test.go:126,224` (**atlandı**) | Bakiye tabanı **Uyumlu**; **payload bütünlüğü Yok** (aynı anahtar + farklı tutar → mevcut kayıt döner: `debitToAccount` anahtar bulununca, `RecordEvent` `:945-949`) · HEDEFLE_FARK |
| 19 | Yükleme: ödeme onayı sonrası tek mali etki | `CreditFromSucceededIntent` köprüsü; `Credit` `wallets/service/service.go:245-293`; banka havalesi maker-checker (`wallets:approve`) | `aypass_kullanici_app` yalnız havale **bildirimi** çağırır (`lib/services/aypass_api_service.dart:507-517`); kartla online yükleme çağrısı bulunmadı | `payment_intent_bridge_test.go:69` (fake) | **Kısmi** · KODDA_GORULDU; uçtan uca DOĞRULANMADI |
| 20 | Belge alındı ≠ faturalandı; unknown sonuç mutabakatı | `billing/invoiceable_records` unique `indexes.go:13-20`, `MarkInvoiced` idempotent `service.go:117`; mock e-fatura/ERP adaptörleri prodda kapalı (`33467c8`), PDF üretimi prodda sahte üretmez (`c3dd848`) | — | `billing` 1 idempotency testi | **Kısmi**; unknown-sonuç mutabakatı DOĞRULANMADI |
| 21 | Pi iş kuralları; onaylı merchant POS ayrı | `chargeCashless` `service.go:1276-1296`: yalnız istemcinin beyan ettiği `InstrumentUsed` (`cash_pos`/`card_pos`) → anında journal; POS onay referansı/PENDING_EXTERNAL yok; `AccessPoint.PosEnabled` var | `pos_enabled` panelde yok (D-02) | `TestChargeCashless*` | **Fark** · HEDEFLE_FARK |
| 22 | Koltuk/ek yolcu/yedek liste/ön ödeme/iptal | `service_reservations`, `room_reservations` var (idempotency 2'şer dosya); araç koltuk rezervasyonu | — | — | **Zincir çıkarılmadı** · DOGRULANMADI |
| 23 | Abonelik/dönem hakkı; indirim kullanım limitleri | `chargeSubscriptionCovered` `service.go:1034`; dönem unique `subscriptions/repository/indexes.go:51-57`; `Charge` `idempotencyKey`'i **kullanmıyor** `subscriptions/service/service.go:381`; `SyncPaymentIntent` boş gövde `:441`, çağıranı yok | — | `subscriptions` idempotency testi yok | **Kısmi** · KODDA_GORULDU (D-16) |
| 24 | Kampanyalı yükleme; nakit/promosyon ayrımı | — | — | — | **Zincir çıkarılmadı** · DOGRULANMADI |
| 25 | Misafir OTP ve satın alınan QR hakkı ayrı, tek tüketim | — | — | — | **Zincir çıkarılmadı** · DOGRULANMADI (eski AYG-02…05 yeni uygulamaya taşınmadı) |
| 26 | Kurum/SAP/SMS/GPS entegrasyonu, izlenebilir tekrar | `gursel_sap_sync`, `sms`, `telematics`, `integrations`, `institution_directory` modülleri var | — | — | **Zincir çıkarılmadı** · DOGRULANMADI |
| 27 | Mercedes/LCW tanımlı fabrika noktaları | `AccessPointKindFactoryGate` `access_point.go:19` | — | — | tür **var**; tanımlı noktalar DOGRULANMADI |
| 28 | Kasa, rota ortalaması, araç/personel türü raporları | `routes.go`: `/reports/kind-summary`, `route-summary`, `person-summary`, `corporate-quota-usage` | panel sayfaları DOĞRULANMADI | — | **Kısmi** · rota ortalaması/araç-personel türü DOGRULANMADI |
| 29 | Model 1 sağlayıcı devri, kesim, eski/yeni ilişki | `service_assignment*` sıfır sonuç; en yakın `TripContract` `EffectiveFrom/To` `trip_contract.go:47-48` | — | — | **Yok** · HEDEFLE_FARK |
| 30 | Devirde bakiye/borç/hak/ödeme profili; operasyon ≠ mali kapanış | `Wallet` kişi merkezli (devir bakiyeyi taşımak zorunda değil — doğru); borç/ödeme profili devri yok | — | — | **Kısmi** · DOGRULANMADI |
| 31 | Üye işyeri/şube/terminal/sponsor/program/fon ayrımı | `pocket|merchant_program|FundingAuthority|wallet_pockets` Go'da sıfır (yalnız test not metinleri); `merchant` yalnız banka kimlik anahtarları `payments/gateways/service/service.go:405-465` | — | — | **Yok** · HEDEFLE_FARK |
| 32 | Market iade/komisyon/hakediş/settlement ve dış onay | `settlement` yalnız seyir nakit tahsilatı `rider_trip.go:34-40`; hakediş yalnız `trip_billing` | — | — | **Yok** · HEDEFLE_FARK |
| 33 | Backend DTO ↔ panel istemcisi/tip/form | Bkz. §4.1 D-01…D-06 | Panel `lib/api/access-points.ts`, `wallets.ts`, `persons.ts` | `event_response_test.go` (yalnız Go); çapraz sözleşme testi **yok** | **Fark** · KODDA_GORULDU |
| 34 | 28 dil, menü, permission, view-profile, event dokümantasyonu | Approval menüleri + 28 dil (`3666773`, `931e70f`) önceki oturumda | AYPASS modülleri için tarama **yapılmadı**; `check-i18n` baseline 61 | — | **Zincir çıkarılmadı** · DOGRULANMADI |
| 35 | Tekrar çalıştırılabilir migrasyon; çift tahsilat/fatura yok | `cmd/migrate_gursel` var; tekrar-çalıştırma güvenliği incelenmedi; binary kaldırma kararı açık | — | — | **Zincir çıkarılmadı** · DOGRULANMADI |
| 36 | Test/kanıt/gözlenebilirlik/restore, manuel kontrollü yayın | Workflows + `deployments/` var; Mongo'suz sandbox'ta eşzamanlılık testi koşulamıyor | — | §1.4 | **Kısmi** |

---

## 3. Öncelikli alan incelemesi (kısa karar notları)

| Alan | Bulgu | Etiket |
|---|---|---|
| Person / Wallet / Instrument / Ledger sahipleri | Person → Wallet (kişi×para birimi unique) → kart/QR (`WalletID`). Kartta bakiye yok. Cüzdan önbelleği `CachedBalance` float64 (`wallet.go:75`); para birimi alanları float64 (AD-01 kararı bekliyor). | KODDA_GORULDU |
| Bakiye kontrolü + journal + cache atomikliği | Düşümde: koşullu `$inc` ile rezervasyon (`debitToAccount` `:338-346`) → `PostJournal` → damgalama; journal başarısızsa rezervasyon iade (`settleFailedReservation`). Girişte: journal → idempotent cache uygulaması (`applied_journal_ids`). Tek transaction yok; çökme anında cache ledger'ın **altında** kalır (tasarım kararı, `service.go:974-980` yorumu). | KODDA_GORULDU · testler **atlandı** |
| Aynı anahtar tekrarı / farklı anahtar paralel | Ledger'da kısmi unique indeks; `recoverFromLostRace` kazananın journal'ını döndürür. **Olay seviyesinde:** `AccessEvent` unique indeks var, fakat `Insert` çift anahtarda ham hata döner (`mongo.go:250-251`); `RecordEvent` ön kontrolü yarışı kapatmaz. `rider_trips` indeksleri **unique değil** (`rider_trips/repository/indexes.go:12,23`) → aynı kişinin farklı anahtarlı paralel iki girişi iki açık seyir yaratabilir. | RISK_TEST_GEREKLI |
| Plaka + aktif araç-seyir bağımlılığı | Ücret kararı plakaya bağlı (REQ-07). Hedef: plaka ayrı, doğrulanmış kurum/güzergâh/yön/seyir/uygulanan tarife bağlamı. | HEDEFLE_FARK |
| Mevcut metered vs A1–A4 | Tablo §2 REQ-09…12. Özet: A2'nin sabit-ücret zemini ve A4'ün matris zemini var; A1/A3/A4 algoritması, binişte maksimum tutarın tahsili ve inişte fark iadesi yok. | HEDEFLE_FARK |
| A3/A4 çıkış farkı | Yok. Çıkış `chargeMetered` `:1563` `CloseTrip` hatasını **yutuyor** (`_ =`), giriş `OpenTrip` hatasını da (`:1502`): ücret alındı ama seyir açık kalırsa sonraki dokunuş yeniden "çıkış" sayılıp ikinci kez ücretlenebilir. | RISK_TEST_GEREKLI |
| Olay anındaki tarife/kart sahibi/sağlayıcı | `AccessEvent`: `BaseFee`, `ChargedAmount`, `DiscountRuleID`, `VehiclePlate`, `ScheduledRouteID`, `TripAssignmentID`, `RouteStopID` var; tarife kimliği/sürümü, algoritma, sağlayıcı/firma, kurum, kart kimliği, `occurred_at`, cihaz kimliği **yok**. `RevenueInstitutionID` yalnız kantin satış yolunda dolu (`service.go:985`). | HEDEFLE_FARK |
| Cihaz gateway yetkileri | `/device/*` yolları `apiGatewayGroup` altında: `APIClientAuthRequired` + rate limit + `APIClientRouteAuthorization` (AllowedRoutes + scope) (`cmd/api/main.go:2576-2625`; `auth/middleware/api_client_gateway.go:300`). Yetkiler `access_points:record` vb. (`auth/handler/admin.go:39-49`). **`:id` ile gelen olay/nakit/iptal/walk-on uçları cihazın bağlı erişim noktasını denetlemiyor**; bağlama (`ResolveByDeviceBinding`) yalnız `/me` uçlarında. `routes.go:9-11` yorumu "cihaz RecordEvent ertelendi" diyor, oysa `main.go:2603` bağlı (bayat yorum). | KODDA_GORULDU · etki RISK_TEST_GEREKLI |
| Gerçek offline kalıcı kuyruk | Hiçbir AYPASS uygulamasında yok (bkz. §4.1 D-11). Eski uygulama da bunu kanıtlamıyordu (eski §6.6). | KODDA_GORULDU |
| Backend DTO ↔ panel | §4.1. | KODDA_GORULDU |
| Ödeme/POS onayı ↔ dahili cüzdan tüketimi | Cüzdan düşümü (`Debit`) ödeme adaptöründen bağımsız ve doğru ayrılmış; **POS ödemesi** ise istemci beyanıyla anında "granted" journal (REQ-21). | KODDA_GORULDU / HEDEFLE_FARK |
| Tek mali etki: yükleme, iade, fatura, rezervasyon, abonelik | Yükleme/banka onayı: idempotent köprü + `applied_journal_ids`; iade: reserve → adapter → `in_doubt` (önceki oturum, `93b0fe4`); fatura: unique işaret; abonelik dönemi: unique (subscription, period_start) ama ödeme/cüzdan etkisi bağlı değil; rezervasyon: ayrı yazımlar, koltuk yarışı testi yok. | karışık (REQ-19/20/22/23) |
| Sabancı taşıma devri Gürsel → yeni firma | Kod sahibi yok (§2 REQ-29). `trip_billing` sağlayıcıyı **sefer sonunda aracın güncel `SupplierID`'sinden** okur (`fleet_tally/trip_billing/service/service.go:66-101`), tarifeyi `time.Now()` ile çözer, hataları sessizce yutar. | HEDEFLE_FARK |
| Global Market ayrı yapı | Yok; ulaşım cüzdan yolunda merchant/settlement fonu karışmıyor. Aynı `Wallet` kantin/fotokopi satışında da kullanılıyor (`DebitToAccount` `:312`) — kurum geliri, pocket ayrımı yok. | KODDA_GORULDU / HEDEFLE_FARK |

---

## 4. Doğrulanmış farklar ve doğrulanmamış riskler

### 4.1 Doğrulanmış farklar (kodda görüldü)

| ID | Fark | Kanıt |
|---|---|---|
| D-01 | **Panel düzenleme plakayı siliyor.** `parseAccessPoint` `vehicle_plate` okumuyor; form `res.vehicle_plate ?? ""` ile boş yükleniyor ve `undefined` gönderiyor; backend koşulsuz atıyor. | Panel `lib/api/access-points.ts:43-71`, `app/(dashboard)/admin/access-points/[id]/page.tsx:182,311`; backend `service.go:688` (`current.VehiclePlate = …`). Panel `c5e0994`. **WP-1: DÜZELTİLDİ / MERGE BEKLİYOR** — panel dalı `fix/aypass-accesspoint-contract` `46b3a29` (yerel, `main`'e birleşmedi, push edilmedi). Yalnız panel istemcisi düzeldi; backend WP-15 ile tam-doküman PUT + version + alan-varlıklı PATCH sözleşmesine geçti (aşağıda ve §6 WP-15). |
| D-02 | **`pos_enabled` panelde hiç yok;** her PUT `false` yapar. | `grep` panelde sıfır sonuç; backend `service.go:718`. **WP-1: DÜZELTİLDİ / MERGE BEKLİYOR** (aynı dal, aynı commit): form `pos_enabled`'ı taşır ve her kayıtta açık `true`/`false` gönderir. Panelde POS'u değiştiren bir **arayüz kontrolü hâlâ yok** (eklemek 28 dilli çeviri anahtarı gerektirir; ayrı karar). |
| D-03 | Panel `fee: 0` gönderir, backend `if in.Fee > 0` ile yok sayar → ücret 0'a indirilemez. | `service.go:710-712`. |
| D-04 | `walletsApi.list` ve `personsApi.list` eski zarf biçimini okuyor (`body.total_pages`), oysa zarf artık `data=[items]` + kardeş `pagination`; sayfalama hep 1 sayfa. | Panel `lib/api/wallets.ts:86-94`, `lib/api/persons.ts:94-102`. |
| D-05 | `birth_date`: panel `YYYY-MM-DD` gönderir, backend `time.RFC3339` bekler → 400 olasılığı (tarayıcıda doğrulanmadı). | Panel `persons/create/page.tsx:41`; backend `persons/handler/admin.go:289`. |
| D-06 | `parseAccessEvent` 13 alanı düşürür (plaka, seyir, durak, `ledger_entry_id`…): yalnız gösterim kaybı. | `access-points.ts:89-116`, DTO `dto.go:104-125`. |
| D-07 | A4 eksik hücre/hata → sessiz mesafe fallback'i. | `service.go:1539`, `:1760` (`err == nil && fare != nil` dışı durumda hata yutulur). |
| D-08 | `metered`: giriş ücretsiz, binişte maksimum tutar tahsil edilmiyor, inişte fark iadesi yok. | `access_point.go:65-67`; `service.go:1468-1471`. |
| D-09 | Tarife/indirim `time.Now()` ile; olayda tarife anlık görüntüsü yok. | `service.go:1101,1131,1169,1539,1760`; `access_event.go`. |
| D-10 | `occurred_at` yok; `CreatedAt` sunucu zamanı. | `service.go:1827-1850`. |
| D-11 | `aypass_erisim_app`: kalıcı kuyruk yok (yalnız `flutter_secure_storage`/`http`/`uuid`), olay isteğinde idempotency anahtarı yok, zaman damgası yok, HTTP timeout yok. Çevrimdışı dokunuş **kaybolur**; yanıt kaybolursa operatör yeniden okutur → anahtar olmadığından çift tahsilat riski (kod okumasından çıkarım). | `lib/services/validation_service.dart:66,73,169-183,241,273,374` (alt ajan raporu). Ana oturumda yerel grep ile doğrulanan: `idempotency`/`occurred_at` yokluğu (erisim + kt_surucu `lib/`). Bağımlılık listesi ve timeout yokluğu alt ajan grep'ine dayanır. |
| D-12 | `aypass_kullanici_app/lib/core/api_config.dart:11` → `kBaseUrl = 'http://localhost:9090'` (üretim adresi yorum satırı). Benzer: `kt_surucu_app/lib/core/api_config.dart:8`. Release derlemesinin hangi adrese gideceği belirsiz. | Yerel grep ile doğrulandı. |
| D-13 | POS ödemesinde onay referansı/durumu yok; `cash_pos`/`card_pos` istemci beyanı. | `service.go:1276-1296`. |
| D-14 | `OpenTrip`/`CloseTrip` hataları yutuluyor; açık seyir için unique indeks yok. | `service.go:1502`, `:1563`; `rider_trips/repository/indexes.go:12,23`. |
| D-15 | Gateway `:id` olay uçları cihaz–erişim noktası bağını denetlemiyor. | `main.go:2603,2616,2619,2620`; `ResolveByDeviceBinding` yalnız `/me` uçlarında (`handler/admin.go:157,170,184,210`). **WP-2: DÜZELTİLDİ / MERGE BEKLİYOR** (yerel dal `fix/aypass-device-accesspoint-binding` `bb1ec78`, taban `3a3bb86`, `main`'e birleşmedi, push edilmedi). |
| D-21 | **WP-2 sırasında bulundu, yeni.** `RecordEvent`'in idempotency-anahtarı kısayolu (`FindByIdempotencyKey`) yalnız tenant/site ile sınırlı, erişim noktasıyla değil; aynı anahtarı tahmin eden/tekrar eden farklı bir cihaz, başka bir erişim noktasının olayını (cüzdan/kişi/tutar dâhil) D-15'in bağ denetimini hiç görmeden alabiliyordu. **WP-2: düzeltildi** (aynı commit). | `access_points/service/service.go` `RecordEvent`'in anahtar-eşleşme bloğu; `repository/mongo.go:286-299` `FindByIdempotencyKey`. |
| D-22 | `wallet_instruments.ValidateQr`/`ValidateCard`, istemcinin isteğe bağlı olarak beyan ettiği `access_point_id`'yi, kart/QR önbellekte bulunamazsa hangi "canlı doğrulama" bağlayıcısının sorgulanacağına karar vermek için kullanıyordu — cihaz bağı denetlenmiyordu. Mali/defter/sefer etkisi yok (yalnız dış sorgu seçimi); yanlış kuruma sorgu "var mı" sızıntısı riski vardı. **2026-10-30 düzeltmesi: DÜZELTİLDİ.** Bu satır önceki taslakta "DÜZELTİLMEDİ (kapsam dışı bırakıldı)" diyordu — bu, aynı raporun §6 WP-2 satırındaki "WP-2.1: D-22 düzeltildi" ifadesiyle çelişiyordu; kod okunarak doğrulandı ve tutarsızlık giderildi: `tryLiveVerify`, `callerIsAPIClient` ise `ap.DeviceBindings` içinde `apiClientID`'yi arıyor, bağlı değilse `ErrLiveVerifyAccessPointNotBound` ile reddediyor. İnsan-admin çağrısı etkilenmiyor (tasarım gereği). | `wallet_instruments/handler/admin.go:139,161`; `wallet_instruments/service/service.go:219-242` (`tryLiveVerify`, WP-2.1 commit `51330f5`, yerel dal, `main`'e birleşmedi). |
| D-23 | **WP-2 sırasında bulundu, DÜZELTİLMEDİ (kapsam dışı bırakıldı, WP-2'nin düzeltmesiyle kısmen rastlantısal kapandı).** `CancelBoardingTrip` hiçbir zaman `requireActiveAccessPoint` çağırmıyordu (SettleTripCash/RecordWalkOnCashEntry çağırıyor). WP-2'nin koruması kendi aktiflik kontrolünü içerdiği için **api_client yolunda** rastlantısal olarak kapandı; **insan-admin yolunda hâlâ açık** (bu uç insan-admin için hiç kayıtlı değil, ama ileride eklenirse sorun sürer). | `access_points/service/service.go` `CancelBoardingTrip` (3a3bb86 öncesi hâli). |
| D-24 | **WP-15 turunda bulundu, DÜZELTİLDİ / MERGE BEKLİYOR.** `access_points:manage_devices` izni TEK BAŞINA (genel `create`/`update` izni olmadan) `PermissionMiddleware`'in modül-geniş kapısından geçip Create/Update/Patch'e ulaşabiliyordu — Create'te cihazsız bir nokta bile oluşturabiliyordu, PATCH'te ise sadece `device_bindings`'i değiştiren bir istek `access_points:update` olmadan başarılı oluyordu (eski kabul edilmiş bir test bunu "doğru" sayıyordu). Artık Create `access_points:create`'i, Update/Patch `access_points:update`'i kendi içinde AYRICA ve KOŞULSUZ kontrol ediyor; `manage_devices` hiçbir durumda bunların yerine geçmiyor. Gerçek `router.ProtectedChain` + sentetik rol/token ile 11 test, red/green kanıtıyla doğrulandı. | Backend `feat/aypass-accesspoint-patch` `dbfc453`→`beee15d`→`6ebe853` (taban `51330f5`/WP-2.1, yerel, `main`'e birleşmedi); `handler/admin.go`, `handler/patch.go`, `handler/real_auth_chain_test.go`. |
| D-25 | **WP-15 turunda bulundu, DÜZELTİLDİ / MERGE BEKLİYOR.** İki Flutter istemcisi (`aypass_erisim_app`, `kt_surucu_app`) AccessPoint cihaz-bağlama ekranında backend'in PUT/PATCH sözleşmesiyle hiçbir zaman tam uyumlu değildi: `aypass_erisim_app` tam-doküman PUT gönderiyordu ama `version` alanı hiç yoktu (WP-15'in version-zorunlu backend'ine karşı 400 ile kesintiye girecekti); `kt_surucu_app` hiç var olmayan tekil `device_binding` alanıyla eksik bir PUT gönderiyordu (backend sözleşmesiyle bağımsız olarak hiçbir zaman çalışmadı) VE okuma tarafı da aynı yanlış alanı okuyup mevcut cihaz bağını ekranda hep boş gösteriyordu. İkisi de gerçek, ulaşılabilir ekranlardan (`access_point_detail_screen.dart`) çağrılıyor. Düzeltme: PATCH + sadece `{version, device_bindings}`, birincil cihaz slotu dışındaki bağları koruma, değişiklik yoksa istek yok, version asla 0 varsayılmıyor, 409'da tek seferlik otomatik yeniden deneme, hata asla başarı gibi gösterilmiyor. Gerçek yürütülen Flutter birim testleriyle doğrulandı (`flutter test`, araç zinciri çalışıyor). | `aypass_erisim_app` `fix/aypass-erisim-devicebinding-patch` `3fa8d29` (taban `ec6e999`); `kt_surucu_app` `fix/kt-surucu-devicebinding-patch` `8bdcfb2` (taban `b80a230`); ikisi de ayrı worktree'de, ana checkout'ların kullanıcı WIP'ine dokunulmadı, `main`'e birleşmedi. |
| D-16 | `subscriptions.Charge` idempotency anahtarını kullanmıyor (dönem unique ile korunuyor); `SyncPaymentIntent` boş gövde ve çağıranı yok → abonelik ↔ ödeme/cüzdan etkisi bağlı değil. | `subscriptions/service/service.go:381,441`; `interface.go:56`. |
| D-17 | Aynı idempotency anahtarıyla farklı tutar/yük gelirse mevcut sonuç döner (payload özeti karşılaştırması yok). | `wallets/service/service.go:329-333`; `access_points/service/service.go:945-949`. |
| D-18 | Devir ve Global Market için model/sahip yok. | §2 REQ-29…32. |
| D-19 | `checkInstitutionMatch` yalnız `Status == active` bakıyor; üyelik `StartDate/EndDate` (`persons/model/affiliation.go:101-102`) denetlenmiyor. | `service.go:510`. |
| D-20 | Bayat yorum: `routes.go:9-11` cihaz RecordEvent'in "ertelendiğini" söylüyor; gerçekte bağlı. | `access_points/routes.go`; `main.go:2603`. |

### 4.2 Doğrulanmamış riskler

| ID | Risk | Neden doğrulanamadı / gereken test |
|---|---|---|
| R-01 | Cüzdan düşüm tabanı, aynı anahtar, journal-hata, rapor tutarlılığı testleri geçiyor mu | Mongo yok (§1.4); gerçek replica set ile koşulmalı. |
| R-02 | Aynı anahtarla paralel iki `RecordEvent`: ikinci istek ham duplicate-key hatası mı alıyor | Mongo'lu entegrasyon testi gerekli. |
| R-03 | Aynı kişi paralel iki giriş → iki açık seyir | Mongo'lu test + veri denetimi. |
| R-04 | Ücret alındı, `CloseTrip` başarısız → sonraki dokunuş çift tahsilat | Hata enjeksiyon testi. |
| R-05 | Gateway: bir cihazın başka erişim noktasına olay yazabilmesi | Negatif test (AllowedRoutes + scope tanımlı bir api_client ile farklı `:id`). |
| R-06 | Canlıda plakaların panel düzenlemeleriyle zaten silinmiş olması | Canlı veriye bakmak yasak/yetkisiz; yerel kopya veya salt-okunur denetim onayı gerekli. |
| R-07 | Canlıdaki sürüm/adres (`aypasss…` vs `aypass.gurseltur.com.tr`) | Üretim erişimi yok. |
| R-08 | Tarife çakışma/örtüşme doğrulaması (eski AYG-07'nin yeni karşılığı) | `route_fares` servis kuralları taranmadı. |
| R-09 | Android'de `http://localhost` cleartext politikası ve build pipeline'ında `kBaseUrl` override'ı | `android/` ve CI incelenmedi. |
| R-10 | `aypass_kullanici_app` git deposu değil → sürüm/geri alma izlenemiyor | Kullanıcı kararı (depo oluşturma). |
| R-11 | `.env` izlenen dosya ve yerelde değişmiş (önceki bilinen sızıntı konusu) | İçerik raporlanmadı; **stage/commit edilmemeli.** |
| R-12 | Panel `approval-policy-ui` dalındaki 5 commit `main`'de yok; hangi dalın "canlı" olduğu belirsiz | Birleştirme sırası kullanıcı kararı. |

### 4.3 `f48748f` temelli belgelerin eskiyen maddeleri

Aşağıdakiler Excel eşlemesi/Yol Haritası'nda "var/yok" olarak yazılmış ama `3a3bb86`'da değişmiştir (diğer maddelerin güncelliği **kontrol edilmedi**; 104'lük matris yeniden hesaplanmadı. `f48748f` temel durumu: 3 karşılanıyor / 34 kısmi / 45 yok / 22 çelişkili):

- **MK05** (RecordEvent erişim noktası durumunu denetlemez) → `requireActiveAccessPoint` `service.go:378` artık çağrılıyor.
- **Z07 / Z10 / K19** (oku-hesapla-`$set` yarışı) → bakiye artık koşullu `$inc` (`wallet.go:69-75`, `ApplyBalanceChange`).
- **FX1** satır atıfları (1080-1105, 1151, 1163) → formül artık `service.go:1177`.
- **"Karşılanıyor" 3 madde:** T07 (KODDA_GORULDU, koşulmadı), FX1 (KODDA_GORULDU, satır kaydı), K06 (yalnız Flutter okuma ekranı; backend'de yazma yok).
- Cüzdan eşzamanlılık testleri (`balance_concurrency_test.go`) `main`'de **yok**, yalnız bu dalda.

### 4.4 Eski uygulama kusurlarının yeni uygulamadaki karşılığı

Eski belge (`AYSHUTTLEGURSEL_PROJE_ANALIZI_TR.md` §6, §10) davranışı; yeni uygulamada **var sayılmamıştır**.

| Eski | Eski davranış (BELGEDE_BILDIRILDI) | Yeni uygulamada durum |
|---|---|---|
| §6.1 okuma-kontrol-kayıt atomik değil | Yeni: koşullu `$inc` ile atomik tavan kontrolü (KODDA_GORULDU; test atlandı) |
| AYG-06 kayıt hatası başarıya dönüşür | Yeni: `insertChargedEvent` hatası dışarı döner. Ters risk: **ücret alınıp olay yazılamaz** → cihaz anahtarsız yeniden dener (D-11). |
| AYG-07 tarife çakışma sorgusu yanlış alan | Yeni: DOGRULANMADI (R-08) |
| AYG-02…05 misafir bilet tutar/süre/tek tüketim | Yeni: REQ-25 zinciri çıkarılmadı; taşınmış kabul **edilmedi** |
| §6.2 salt-eklemeli olmayan hareketler | Yeni: ledger append-only + reversal tasarımı (KODDA_GORULDU); silme/güncelleme yollarının tamamı taranmadı |
| §6.6 kalıcı offline kuyruk kanıtlanmamış | Yeni: aynı durum — kuyruk yok (D-11) |
| §6.4 grup tahsilatı tek transaction değil | Yeni: DOGRULANMADI (REQ-22) |
| AYG-08/09 kampanya | Yeni: REQ-24 DOGRULANMADI |

---

## 5. Yeniden kullanılacak modüller ve açık kararlar

### 5.1 Yeniden kullanılacak sahipler (yenisi yazılmayacak)

| Alan | Kullanılacak sahip | Not |
|---|---|---|
| Kimlik / hesap | `persons`, `wallets`, `wallet_instruments` | Person kalıcı hesap; kart/QR araç. |
| Defter | `ledger` (kısmi unique, `ConfirmPendingJournal`, `UpdateJournalStatus`) | Reversal ile düzeltme; yeni defter kurulmayacak. |
| Cüzdan düşümü | `wallets.Debit*` rezervasyon deseni | Biniş tahsilatı (P) ve fark iadesi (P−N) için aynı atomik desen genişletilmeli. |
| Geçiş hattı | `access_points` `RecordEvent` → `resolveEventContext` → `charge*` | A1–A4 bir strateji katmanı olarak buraya. |
| Tarifeler | `fleet_tally/route_fares`, `route_stop_fares`, `access_point_discounts` | `RouteStopFare` A4 matrisi zemini. |
| Seyir | `fleet_tally/rider_trips`, `trip_assignments` | Biniş tahsilatı/fark iadesi için açık-seyir değişmezleri eklenecek. |
| Sağlayıcı/devir zemini | `trip_contracts` (+ `trip_billing`) | Geçerlilik aralığı ve `SupplierID` zemini; `service_assignments` (AD-03) için çıkış noktası. |
| Onay | `approvals` motoru | İade/bakiye kapıları hazır. |
| Ödeme | `payments/intents` (iade rezervasyonu, `in_doubt`, ResolveRefund) | POS onayı için aynı model. |
| Cihaz kimliği | api_client gateway + `AllowedRoutes` + `DeviceBindings` | Bağ denetimi eklenecek. |
| Fatura | `billing/invoiceable_records` | Tek etki zaten var. |

### 5.2 Açık kararlar

| # | Karar | Seçenekler / öneri |
|---|---|---|
| KR-1 | Para tipi (AD-01) | `int64` minör birim (öneri) mi mevcut `float64` mi; etkisi: wallet, ledger, event, tarifeler. Faz 1 ayrı onay bekliyor. |
| KR-2 | Algoritma sahibi | Güzergâh mı, erişim noktası mı, kurum mu? Öneri: `ScheduledRoute` bazlı yapılandırma. |
| KR-3 | A4 eksik hücre | Öneri: açık hata (olay reddedilir/`PENDING`), sessiz fallback yok. |
| KR-4 | Offline risk politikası | Cihaz başına borç tavanı, geçerlilik süresi, imzalı kuyruk; Z07 ile birlikte. |
| KR-5 | Olay zamanı | `occurred_at` (cihaz) + `received_at` (sunucu); tarife `occurred_at`'e göre çözülür; izin verilen saat sapması. |
| KR-6 | `aypass_kullanici_app` için git deposu | Depo oluşturulsun mu (R-10). |
| KR-7 | Canlı adres | `aypasss.aysteknolojileri.com` mı, `aypass.gurseltur.com.tr` mi (R-07). |
| KR-8 | Dal birleştirme sırası | Backend 44 commit / panel `approval-policy-ui` → `main` (kullanıcı tetikler; bu oturumda push/PR yok). |
| KR-9 | `service_assignments` (AD-03) tasarımı | `TripContract` genişletmesi mi, yeni modül mü. |
| KR-10 | Global Market | Ayrı modül/DB koleksiyonları; `wallets`'a karıştırılmayacak (öneri). |
| KR-11 | POS onayı | `PENDING_EXTERNAL` durumu ve onay referansı alanı. |
| KR-12 | `migrate_gursel` ikilisi | Kaldırılsın mı (önceki oturumdan açık). |
| KR-13 | AccessPoint güncelleme semantiği (WP-15) | `PATCH`/işaretçi (öneri: gönderilmeyen alan değişmez) mi, belgelenmiş tam-PUT + eksik alana 400 mü. Karar WP-15 öncesi gerekli. |

---

## 6. Öncelikli iş kayıtları

Önem sırası: veri kaybı/güvenlik → mali doğruluk → hedef özellik. "Kabul testi" = tamamlandı sayılması için geçmesi gereken en küçük kanıt.

| Kayıt | Başlık | Bağımlılık | Kabul testi | Faz |
|---|---|---|---|---|
| **WP-1** | **TAMAMLANDI (yerel dal `fix/aypass-accesspoint-contract`, taban `c5e0994`, birleşmedi — bkz. `AYPASS_WP1_ACCESSPOINT_SOZLESME_RAPORU_TR.md`).** Onaylı kapsam yalnızca: AccessPoint düzenlemesinde `vehicle_plate` ve `pos_enabled`'ın istemeden kaybolmasını önleme + gerçek submit yolunu kapsayan regresyon testleri; id-eşleşme ve güçlendirilmiş yanıt-şekli kontrolü ek tur olarak eklendi. `requires_approval`/`cash_fallback_enabled`/`device_bindings` ayrıca kontrol edildi, sorun bulunmadı. Aşağıdakiler **kapsam dışı bırakıldı ve açık**: create'e plaka alanı, `fee=0`, `wallets`/`persons` sayfalama, `birth_date` biçimi, `parseAccessEvent` alanları (`fee=0` için artık WP-15'te sunucu tarafı karar kaydı var) | — | (1) Plakası olan AccessPoint'i panelden kaydet → plaka korunur; (2) `pos_enabled` gidiş-dönüş; (3) 60+ kayıtlı listede 2. sayfa erişilebilir; (4) `tsx --test` ile DTO fixture ↔ parser testi | P01 |
| **WP-2** | **TAMAMLANDI (yerel dal `fix/aypass-device-accesspoint-binding`, taban `3a3bb86` → `bb1ec78` → `51330f5`, birleşmedi).** Gateway: `events`/`cash-settlement`/`cancel`/`walk-on-cash-entry` uçlarında + `RecordEvent`'in idempotency-tekrar yolunda api_client ↔ erişim noktası bağ denetimi eklendi; insan admin etkilenmedi; bayat `routes.go` yorumu **düzeltilmedi** (kapsam dışı bırakıldı, hâlâ açık). **WP-2.1** (aynı dal, `51330f5`): D-22 (ValidateQr/ValidateCard canlı-doğrulama cihaz bağı) düzeltildi, D-23 (CancelBoardingTrip aktiflik kontrolü) belgelenip kasıtlı test edildi, düzeltilmedi (insan-admin yolu yok). Bu paketler tek başına üretime hazır sonucu doğurmaz (bkz. §8.1) | WP-1'den bağımsız | Bağlı olmayan cihazla farklı `:id`'ye POST → 403; bağlı cihaz → 201; negatif test + AllowedRoutes testi | P01/P05 |
| **WP-3** | Seyir değişmezleri: açık seyir için kısmi unique indeks (veri denetimi sonrası); `OpenTrip/CloseTrip` hatalarını yutmama; çift anahtarlı olay `Insert`'ünün mevcut olayı döndürmesi | WP-13 (Mongo'lu test ortamı) | Paralel iki farklı anahtarlı giriş → tek açık seyir; `CloseTrip` hata enjeksiyonunda ikinci dokunuş çift ücret yok | P02/P04 |
| **WP-4** | Olay anı sözleşmesi: `occurred_at`/`received_at`, uygulanan tarife kimliği+sürümü, algoritma, sağlayıcı/firma, kurum, kart/enstrüman ve cihaz kimliği; tarife çözümü `occurred_at`'e göre | KR-1, KR-2, KR-5 | Geçmiş tarihli olay eski tarifeyle ücretlenir; olay kaydı tarife kimliğini taşır | P03/P05 |
| **WP-5** | Ücret strateji katmanı: A1–A4 sözleşmesi, A4 eksik hücrede açık hata (KR-3), `metered` uyumluluğu | WP-4, KR-2, KR-3 | Tablo-güdümlü test: her algoritma × (bilinen/eksik hücre); `metered` eski davranışı korunur | P04 |
| **WP-6** | A3/A4: binişte maksimum tutarın iç cüzdandan fiilen tahsili (P), inişte nihai ücretin (N) mesafe/matristen hesabı + indirim, farkın (P−N) ledger reversal ile iadesi; atomik, idempotent. Banka/POS tahsilatından ayrı (Not 2) | WP-3, WP-4, WP-5, KR-1 | Yetersiz bakiye binişte reddedilir (P tahsil edilemez); çıkışta N hesaplanır, `P−N` tek etkili iade edilir; tekrar çıkış çift iade yaratmaz; indirimli N doğru | P04 |
| **WP-7** | Cihaz: kalıcı kuyruk (sqflite/drift), tap anında idempotency anahtarı, `occurred_at`, ACK, geri-çekilme, HTTP timeout; sunucuda sıra dışı/duplicate dayanıklılığı; `kBaseUrl` yapılandırması (`--dart-define`) | WP-2, WP-4, KR-4 | Uçak modunda 3 dokunuş → yeniden bağlanınca tek tek ve tek kez işlenir; yanıt kaybı çift tahsilat yaratmaz | P05 |
| **WP-8** | Idempotency payload bütünlüğü (anahtar + özet; farklı yük → 409) | KR-1 | Aynı anahtar farklı tutar → 409 | P02 |
| **WP-9** | POS ödemesi: onay referansı, `PENDING_EXTERNAL`, dahili cüzdan tüketiminden ayrım | KR-11 | Onaysız POS olayı mali etki yaratmaz; onay gelince tek etki | P06 |
| **WP-10** | Abonelik ↔ ödeme/cüzdan bağlantısı (`SyncPaymentIntent`, `Charge` anahtarı) | — | Dönem ücreti tek mali etki; tekrar çağrı yeni etki yok | P07 |
| **WP-11** | Devir (Model 1): `service_assignments` / TripContract genişletmesi, kesim anı, `trip_billing` sağlayıcıyı olay anından çözer | KR-9, WP-4 | Kesim öncesi olay eski firmaya, sonrası yeniye yazılır; kişi/kurum/bakiye değişmez | P08 |
| **WP-12** | Global Market (Model 4) ayrı modül tasarımı | KR-10 | Tasarım kaydı onayı; ulaşım cüzdan yoluna dokunmaz | P10 |
| **WP-13** | Test altyapısı: Mongo replica set'li yerel/CI koşusu; panel↔backend sözleşme testi pipeline'a | — | `AYSTEK_TEST_MONGO_URI` ile 12 atlanan test geçer | P00/P09 |
| **WP-14** | Hijyen: `kBaseUrl` localhost (2 uygulama), `aypass_kullanici_app` git deposu, `.env` izlenmesi | KR-6 | Release derlemesi üretim adresine gider; `.env` izlenmez | P11 |
| **WP-15** | **TAMAMLANDI — kod ve hedefli doğrulama bitti; birleştirme/yayın kararı bekliyor (2026-10-30).** Backend dalı `feat/aypass-accesspoint-patch` (`6317608`→`651f596`→`cc76b84`→`dbfc453`→`beee15d`→`703f05d`→`6ebe853`→`10ad86e`, taban `51330f5`); panel dalı `feat/aypass-accesspoint-patch` (`9ce439a`→`1f1e60c`, taban `a9752d6`); ayrıca `aypass_erisim_app@3fa8d29` ve `kt_surucu_app@8bdcfb2` (D-25). Tamamı yerel, `main`'e birleşmedi, push edilmedi. `PATCH /admin/access-points/:id` eklendi (eski `PUT` dokunulmadan kaldı, bilinçli bir açık geçiş yolu, artık versiyon+yetki kontrolünü PUT da paylaşıyor); alan-varlığı izleyen kısmi güncelleme, `version` tabanlı atomik eşzamanlılık kontrolü (409), eski versiyonsuz Mongo kaydının ilk güncellemesi doğru çalışıyor. **Nihai izin matrisi** (istisnasız): POST cihazsız→`create`; POST cihazlı→`create`+`manage_devices`; PUT/PATCH→HER DURUMDA `update`, cihaz bağı gerçekten değişiyorsa (PUT) veya istekte `device_bindings` varsa (PATCH) ayrıca `manage_devices` (D-24). **Önceki risk kapandı:** backend'in atomik `FindOneAndUpdate`'i artık bu oturumda kurulan gerçek, geçici, sentetik-only MongoDB üzerinde doğrulandı (9 `TestMongoIT_*` + 11 `TestRealAuthChain_*`, hepsi PASS, red/green kanıtlı) — `AYPASS_WP15_ACCESSPOINT_PATCH_KARAR_KAYDI_TR.md` artık yalnızca karar kaydı değil, test edilmiş bir sözleşmenin belgesi. | KR-13 (karar verildi, uygulandı) | §5'teki senaryoların backend+panel+iki Flutter istemcisi karşılıkları gerçek middleware zinciri + gerçek Mongo + gerçek `flutter test` ile doğrulandı; **canlıya hiç deploy edilmedi, merge/push yok** | P01 |

---

## 7. İlk küçük iş paketi — WP-1 (onaylandı ve uygulandı)

> **Durum:** Hüseyin Koç tarafından onaylandı; onaylı kapsam bu bölümdeki önerilenden **daha dardır** (yalnız `vehicle_plate` + `pos_enabled`; sayfalama, `birth_date`, `fee=0`, create formu ve olay parser'ı dışarıda). Uygulama ve sonuçlar `AYPASS_WP1_ACCESSPOINT_SOZLESME_RAPORU_TR.md`'de. Aşağıdaki metin P00 önerisinin özgün hâlidir.

**Neden bu?** (a) Aktif, doğrulanmış veri kaybı: panelden yapılan her düzenleme plakayı siliyor ve plaka olmadan ulaşım olayı kaydedilemiyor. (b) Yalnız panel + küçük sözleşme testi; mali koda, migrasyona, canlıya dokunmaz. (c) Yol Haritası'nın önerdiği ilk iş emriyle örtüşüyor. (d) Sonraki iş kayıtları (WP-2, WP-7) için sözleşme testi zemini hazırlar.

**Kapsam:** D-01…D-06 (panel: `lib/api/access-points.ts`, `lib/types/access-points.ts`, `access-points/[id]/page.tsx`, `create/page.tsx`, `lib/api/wallets.ts`, `lib/api/persons.ts`, `persons/*/page.tsx`). Backend değişikliği gerekirse yalnız D-03 (`fee=0`) için; önce karar notu.

**Ön koşullar:** Panel için hangi dalın temel alınacağı (KR-8: `main` `c5e0994` önerilir; `approval-policy-ui` ayrı), çalışma ağacındaki `tsconfig.tsbuildinfo`/`incognito_check.mjs` değişikliklerine dokunmama.

**Kabul:** §6 WP-1 testleri; `tsc` hata sayısı baseline'ı (328) artmaz; `check-i18n` baseline'ı (61) artmaz; yeni metinler 28 dilde.

**Kapsam dışı (bilinçli):** Finansal akışlar, para tipi, offline kuyruk, ücret algoritmaları.

---

## 8. Sınırlar ve dürüstlük notları

- Her alanın route→middleware→handler→service→repository→event/worker→API client→UI→test zinciri çıkarılmadı; çıkarılmayanlar matriste **Zincir çıkarılmadı** olarak işaretli.
- Panel, Flutter ve devir/Market taraması üç salt-okunur alt ajanla yapıldı. Ana oturumda kodda yeniden doğrulananlar: D-01, D-02, D-03, D-04, D-12, devir/Market yokluğu (`service_assignment` sıfır sonuç), `trip_billing` tedarikçi okuması, idempotency anahtarı yokluğu (cihaz uygulamaları). **Yalnızca alt ajan raporuna dayananlar** (satır numaraları yeniden okunmadı): D-05 ile D-06'nın panel satırları, D-11'in `validation_service.dart` satırları, `kt_surucu_app` satırları, FX1/T07/K06 doğrulamaları ve D-19'un model satır atıfı (`persons/model/affiliation.go:101-102`); D-19'un `service.go:510` tarafı ana oturumda okundu.
- `AndroidManifest`/native kod, CI iş akışlarının içeriği, `docs/FRONTEND_BACKEND_CONTRACT.md` ve `deployments/production/*` **okunmadı**.
- Test koşusu yalnız Mongo'suz birim testlerini kapsar (§1.4). Panel ve Flutter testleri çalıştırılmadı.
- Hiçbir gizli değer, gerçek kişi/kart bilgisi veya ham QR token yazılmadı.

### 8.1 Açık keşif işleri (kapanmış sayılmaz)

P00 bunlarla birlikte "kapsam eksikleri kayıtlı" bir rapordur. İlgisiz oldukları sürece WP-1'i durdurmamışlardır; ilgili pakete başlamadan önce tamamlanmalıdır.

| # | Açık iş | Gereken paket |
|---|---|---|
| O-1 | `AYSHUTTLEGURSEL_TEKNIK_ENVANTER_TR.md` **okunmadı** | P00 devamı / ilgili faz |
| O-2 | `AYSHUTTLEGURSEL_CURRENT_STATE.md` **okunmadı** (tarihsel; düşük öncelik) | P00 devamı |
| O-3 | `AYSHUTTLEGURSEL_PROJE_ANALIZI_TR.md` yalnız §6.1–6.6 ve §10 okundu; **§1–5, 7–9, 11–14 okunmadı** | P00 devamı |
| O-4 | REQ-22 (koltuk/ek yolcu/yedek liste), REQ-24 (kampanyalı yükleme), REQ-25 (misafir OTP/QR), REQ-26 (entegrasyonlar), REQ-34 (28 dil/menü/permission/view-profile), REQ-35 (migrasyon) zincirleri **çıkarılmadı** | WP-10, P07, P11 öncesi |
| O-5 | Ana Prompt P00 görev 5 (eski D1/D3 iş akışlarının tam çıkarımı) yalnız §6/§10 üzerinden yapıldı | O-3 ile birlikte |
| O-6 | Canlı sürüm/adres, Android cleartext/`kBaseUrl` override, `deployments/production/*`, CI iş akışı içerikleri **okunmadı** | KR-7, WP-14 |
| O-7 | P00 oturumunda cüzdan/ledger eşzamanlılık testleri **atlanmıştı** (12 SKIP). **2026-10-30 güncellemesi:** belirlenen 12 cüzdan/ledger eşzamanlılık senaryosu, bu oturumda kurulan sentetik gerçek MongoDB üzerinde çalıştırıldı ve **başarısızlık gözlenmedi** (12/12, SHA `703f05d`). Bu, tüm mali çekirdeğin doğrulandığı, A1–A4'ün tamamlandığı veya offline tahsilatın güvenli olduğu anlamına **gelmez** — kapsam yalnız bu 12 senaryodur. | WP-13 (CI'a kalıcı Mongo entegrasyonu hâlâ ayrı iş) |
| O-8 | Gateway cihaz–erişim noktası bağı (D-15) | **WP-2 ile düzeltildi** (yerelde, push edilmedi) |
| O-9 | WP-2'de kapsam dışı bırakılan D-22 ve D-23. **2026-10-30 güncellemesi:** D-22 (ValidateQr/ValidateCard'ın client-beyanlı `access_point_id`'si) kod okunarak **düzeltildi** doğrulandı (WP-2.1, `51330f5`) — önceki taslaktaki "düzeltilmedi" ifadesi bu raporun kendi WP-2 satırıyla çelişiyordu, giderildi (bkz. D-22). D-23 (CancelBoardingTrip'in eksik aktiflik kontrolü, insan-admin yolu için) **hâlâ açık** — bu round'da dokunulmadı, insan-admin için bu uç hiç kayıtlı olmadığından risk pasif kalıyor. | D-23 için ayrı iş kaydı, henüz açılmadı |
| O-10 | WP-15 (AccessPoint PATCH sözleşmesi) — önce karar kaydı vardı, üretim kodu yoktu. **2026-10-30 güncellemesi: artık üretim kodu VAR, test edildi, TAMAMLANDI (merge/yayın bekliyor)** — bkz. §6 WP-15 satırı, D-24, D-25. Panelin "yalnız değişeni gönder"e geçişi zaten önceki turda yapılmıştı (`1f1e60c`); lost-update riski artık version+izin kontrolüyle kapalı. | `AYPASS_WP15_ACCESSPOINT_PATCH_KARAR_KAYDI_TR.md` |
