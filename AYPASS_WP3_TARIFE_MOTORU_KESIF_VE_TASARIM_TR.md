# AYPASS WP-3 — Ulaşım Tarife Motoru / Sefer Bağlamı Keşfi

**Tarih:** 2026-10-30 · **Oturum türü:** yalnızca okuma/keşif+tasarım · **Durum (ilk yazım):** Üretim kodu bu paketinde DEĞİŞTİRİLMEDİ. Bu belge, WP-15'in kapanmasının ardından açılan bir sonraki faz için analiz+tasarım çıktısıdır; A1–A4 üretim kodlaması için Hüseyin Koç'un ayrı onayı gerekir.

**Durum (2026-10-30 güncellemesi):** WP-3.1→WP-3.5 **TAMAMLANDI** (kod + hedefli doğrulama + WP-3.5 kapanış turu dahil); son SHA `2764898`, dal `feat/wp3-fare-algorithm-foundation` (worktree `/workspace/go/aystek-core-wp3`, taban `3a3bb86`). Merge/yayın kararı bekliyor. Ayrıca bu dal, WP-15'in kabul edilmiş AccessPoint sözleşme hattıyla (`feat/aypass-accesspoint-patch@10ad86e`) yeni bir yerel entegrasyon dalında (`integration/aypass-accesspoint-fare-engine`, merge commit `4dc8c7a`) birleştirildi ve güvenlik+mali+regresyon kabul testleriyle doğrulandı — bkz. §9 ve WP-3.5 satırının kanıt listesi. WP-3.6 (A1) için üretim kodlaması onaylandı, bu entegrasyon provasının yeşil sonucuna bağlı olarak.

**Kapsam dışı (bu paket):** firma devri, Global Market, POS entegrasyonu, offline cihaz kuyruğu implementasyonu, raporlama geliştirmesi. Bunlara etkisi varsa yalnız bağımlılık olarak kaydedilmiştir.

**Satır numaraları:** backend worktree `/workspace/go/aystek-core-wp15`, bu paket için okunan HEAD (`10ad86e` + okuma-amaçlı, kod değişikliği yok).

---

## 0. Kanonik tanımlar (değiştirilmedi, aynen kabul edildi)

- **F0** — Mevcut sabit/sabit indirimli geçiş.
- **A1** — Kampüs/fabrika sabit uçtur. Sabah biniş→kampüs/fabrika, akşam kampüs/fabrika→iniş. Mesafeden ücret hesaplanır, indirim uygulanır, net ücret tahsil edilir.
- **A2** — Güzergâh+yön+durak için ÖNCEDEN BELİRLENMİŞ ücret. Mesafeden fiyat türetilmez.
- **A3** — Varış bilinmiyorsa binişte güzergâhın MAKSİMUM ücreti FİİLEN TAHSİL edilir. İnişte ikinci okutmayla mesafeden nihai ücret hesaplanır. İndirim nihai ücrete uygulanır. Fazla tahsilat iade edilir.
- **A4** — A3 ile aynı tahsilat/iade akışı; nihai ücret mesafeden değil `biniş×iniş` ücret matrisinden gelir. Hücre yoksa km tarifesine sessiz fallback YOKTUR.
- Banka/POS tahsilatı ile iç cüzdan tüketimi aynı mali olay DEĞİLDİR (ayrı tutulur).

---

## 1. Mevcut zincir: route → handler → service → repository → wallet/ledger

```
POST /admin/access-points/:id/events  (routes.go)
  → handler.AdminHandler.RecordEvent / handler.PublicHandler (device gateway: /device/access-points/:id/events)
    → service.Service.RecordEvent(ctx, accessPointID, RecordEventInput)   [service.go:1080]
        1. idempotency kontrolü (FindByIdempotencyKey, tenant/site scoped)      [service.go:1081-1096]
        2. AccessPoint yükle, Status/kurum/blok kontrolleri                      [checkAccessBlock :603, checkInstitutionMatch :573]
        3. resolveEventContext(ctx, ap, geo)  → eventContext{ScheduledRouteID, TripAssignmentID, RouteStopID, VehiclePlate}
             ├─ Kind != fare_transport → eventContext{} (plaka/sefer YOK)        [service.go:521-523]
             └─ Kind == fare_transport:
                  ap.VehiclePlate boşsa → ErrTransportVehicleRequired           [service.go:524-527]
                  tripAssignments.ActiveByVehiclePlate(plaka)                   [trip_assignments/service.go:98-99]
                    → TripAssignment{ScheduledRouteID, Status=in_progress}      [trip_assignments/model]
                  ScheduledRouteID, ap.EffectiveScheduledRouteIDs() eşleşmeli    [service.go:542-551]
                  resolveNearestStop → RouteStopID                              [service.go:1506]
        4. ap.Kind + ap.PricingMode'a göre dallan:
             PricingModeFree          → insertEvent (ücretsiz)                  [service.go:1161]
             PricingModeFixedFee      → chargeWallet veya chargeCashless        [service.go:1164-1168]
             Kind=canteen/photocopy   → resolveSaleItems → chargeWalletWithFee  [service.go:1142]
             subscription_covered     → chargeSubscriptionCovered              [service.go:1190]
             Kind=fare_transport+metered → chargeMetered (giriş/çıkış)          [service.go:1628]
        5. chargeWallet/chargeMetered → resolveBaseFee/resolveMeteredRate (RouteFare, time.Now()) → resolveFee (discount, time.Now()) → chargeWalletWithFee
             → wallets.Service.Debit/DebitToAccount(walletID, amount, idemKey)  [wallets/service.go:297,323]
                 → repo.ApplyBalanceChange (koşullu $inc rezervasyon)           [wallets/repository/mongo.go:123]
                 → ledger.Service.PostJournal (çift girişli)                    [ledger/service.go]
        6. insertChargedEvent → AccessEvent yazılır (BaseFee/DiscountRuleID/DiscountPercentage/ChargedAmount/LedgerEntryID) [service.go:2031-2066]
```

Cihaz tarafı aynı zincire `/device/access-points/:id/events` üzerinden `APIClientAuthRequired` + cihaz-bağ denetimi (WP-2) ile giriyor; insan admin tarafı (panel) `/admin/...` ile, `access_points:create`/`update`/`manage_devices` izinleriyle (WP-15).

---

## 2. 12 soruya yanıtlar

**1) Plaka/araç ataması ücret zincirinde nerede zorunlu?**
SADECE `resolveEventContext`'te, SADECE `ap.Kind == fare_transport` ise (`service.go:521-527`). Plaka, `TripAssignmentResolver.ActiveByVehiclePlate` ile "şu an devam eden bir TripAssignment var mı" sorusuna cevap bulmak için kullanılıyor; dönen `TripAssignment`'tan SADECE `ScheduledRouteID` ve kendi `ID`'si (`TripAssignmentID`) alınıyor — ücret seçimi (`resolveBaseFee`/`resolveMeteredRate`/`ResolveActiveStopFare`) SADECE `ScheduledRouteID`'yi kullanıyor, plakayı bir daha hiç görmüyor. Yani plaka, ücret seçiminin kendisi için gerekli değil — "şu an gerçekten devam eden bir sefer var mı" kanıtı için gerekli (çoklu-araçlı aynı güzergâh senaryosunda hangi aracın RiderTrip'lerinin "açık" sayılacağını ayırt etmek için de kullanılıyor, bkz. soru 2).

**2) Araç/plaka ile sefer/tarife bağlamını ayırmak için hangi mevcut modeller kullanılabilir?**
`TripAssignment` (plaka+sefer ilişkisi, `fleet_tally/trip_assignments`) zaten ayrı bir model — "sefer bağlamı" (ScheduledRouteID) ile "hangi araç" (VehiclePlate) ZATEN ayrı alanlar, aynı modelde birlikte duruyorlar. Gerçek ayrım ihtiyacı: `RiderTrip` (`fleet_tally/rider_trips`) ve `AccessEvent` ikisi de `TripAssignmentID`'yi taşıyor ama "bu seferin hangi aracı" sorgusu HER ZAMAN `TripAssignment` üzerinden tek dolaylı sorgu gerektiriyor. Mevcut model DEĞİŞTİRİLMEDEN de ayrım mantıksal olarak zaten var; önerilen: ücret/strateji seçimini (bkz. §5) `ScheduledRouteID`+`RouteFare.FareAlgorithm` üzerine kur, plaka/`TripAssignmentID`'yi SADECE "sefer gerçekten devam ediyor mu" ve "hangi aracın açık RiderTrip'leri" sorularına ayrılmış tut — bu zaten mevcut tasarım, yeni bir model gerekmez.

**3) F0 mevcut sabit ücret davranışı nerede?**
`resolveBaseFee` (`service.go:1250-1253`): transport-dışı `Kind`'lar için `return ap.Fee, nil` — doğrudan AccessPoint'in sabit ücreti. `chargeWallet`→`chargeWalletWithFee`→`resolveFee` (discount, SADECE cüzdan kanalında) →`wallets.Debit`. Nakit kanalı (`chargeCashless`, `service.go:1432-1453`) discount'suz, doğrudan ledger'a yazar.

**4) Mevcut `metered` giriş/çıkış state machine'i nasıl çalışıyor?**
Durum tamamen DB-temelli, istemci bayrağı yok: `riderTrips.FindOpenForPerson(personID)` — açık kayıt yoksa GİRİŞ (`service.go:1647-1680`, ücret YOK, `OpenTrip` best-effort, hata `_ =` ile yutuluyor), açık kayıt varsa ÇIKIŞ (`:1682-1727`, `geo.DistanceMeters` ile mesafe, `ResolveActiveStopFare` önce denenir yoksa km×oran, `CloseTrip` de best-effort `_ =` ile yutuluyor). **Doğrulanmış risk:** `rider_trips` indeksleri UNIQUE DEĞİL (`indexes.go:9-32`) — aynı kişinin paralel iki isteği iki açık sefer yaratabilir; `CloseTrip` sessizce başarısız olursa sefer sonsuza kadar "açık" kalır ve sonraki dokunuş onu yanlışlıkla "çıkış" sayabilir.

**5) A1/A2/A3/A4 için hangi mevcut yapılar yeniden kullanılabilir?**
`RouteFare` (route-seviyeli sabit/km-oranı + `EffectiveFrom/To` versiyonlama) → A1/A2/metered-oran zemini. `RouteStopFare` (yönlü O-D çifti, aynı versiyonlama) → A2'nin durak-bazlı hali VE A4'ün matris zemini (BİREBİR aynı model, sadece "hücre yoksa ne olur" davranışı farklı olacak). `RiderTrip` (açık/kapalı sefer, giriş/çıkış GPS+zaman+mesafe) → A1/A3/A4'ün board/alight iki-dokunuş izi için zemin, ama YENİ alanlar gerekiyor (bkz. §5). `wallets.Debit`/`CreditFromAccount` + `ledger.PostJournal` (koşullu `$inc` rezervasyon + çift-girişli defter + idempotency anahtarı) → A3/A4'ün "binişte tahsil, inişte fark iadesi" ikilisi İÇİN HAZIR, YENİ bir mali mekanizma icat etmeye gerek YOK (bkz. soru 6). `access_point_discounts.ResolveBestDiscount` → her algoritmanın nihai ücretine aynı noktadan uygulanabilir (zaten öyle, bkz. soru 8).

**6) A3/A4 için binişte maksimum tutarın FİİLEN TAHSİLİ ve inişte fark iadesi hangi wallet/ledger fonksiyonlarıyla güvenli yapılabilir?**
Biniş: `wallets.Service.Debit`/`DebitToAccount(walletID, maxFare, idemKey="board:"+riderTripID, ...)` — AYNI mekanizma (koşullu `$inc` rezervasyon → `PostJournal`), bugün `chargeWalletWithFee`'nin zaten kullandığı. Yetersiz bakiyede `ErrInsufficientBalance` döner — **binişte sefer/AccessEvent hiç oluşturulmamalı** (bugünkü `chargeMetered` girişi KOŞULSUZ `insertEvent` çağırıyor; A3/A4 için bu değişmeli, bkz. §5/§6). İniş: `wallets.Service.Credit`/`CreditFromAccount(walletID, collected-final, idemKey="adjust:"+riderTripID, ...)` — `RefundEvent`'in BUGÜN ZATEN yaptığı tam desen (yeni, ters-yönlü journal; eski journal'a reversal değil). İkisi de ayrı, deterministik idempotency anahtarlarıyla (`riderTripID`'den türetilmiş) tekrar-güvenli; ledger'daki kısmi-unique indeks çift-tahsilat/çift-iadeyi engelliyor. **Yeni mali mekanizma gerekmez** — sadece `RiderTrip`'e `CollectedAmount`/`CollectedJournalID`/`AdjustmentJournalID`/`RefundAmount` alanları eklenip bu iki çağrı birbirine bağlanmalı.

**7) A4'te eksik matris hücresinin km tarifesine düşmesini engelleyecek karar noktası nerede olmalı?**
Tam olarak BUGÜNKÜ `chargeMetered` satır 1694-1706'daki `if fare, err := ResolveActiveStopFare(...); err == nil && fare != nil && fare.CardFee != nil { stopFareResolved = true } ; if !stopFareResolved { km×oran }` bloğu — bu, A2/metered için DOĞRU davranış (kasıtlı fallback) ama A4 için YASAK. Karar noktası: A4 stratejisinin alighting adımı, `stopFareResolved == false` olduğunda km-fallback'e DÜŞMEMELİ, bunun yerine olayı `PENDING`/reddedilmiş sayıp (örn. yeni `ErrFareMatrixCellMissing`) sürücü/operatöre açık hata döndürmeli — aynı `ResolveActiveStopFare` çağrısı, sadece strateji-bazlı farklı post-koşul.

**8) İndirim hangi katmanda ve hangi anda uygulanıyor?**
TEK katman, TEK fonksiyon: `resolveFee` (`service.go:1317-1334`), `chargeWalletWithFee`'den (TEK çağrı noktası) çağrılıyor — F0/A2/metered-çıkış'ın ÜÇÜ DE aynı noktadan geçiyor, girdi her seferinde o modun NİHAİ hesaplanmış ücreti (`baseFee`). Nakit kanalı hiç indirim görmüyor (mimari dışlama, PersonID eksikliği değil). Metered'da indirim SADECE ÇIKIŞTA uygulanıyor (giriş ücretsiz grant, ücret hesaplanmıyor) — A3/A4 tasarımında da indirim YALNIZ İNİŞTEKİ nihai ücrete uygulanmalı (kanonik tanımla tutarlı: "indirim nihai ücrete uygulanır").

**9) `occurred_at`, tarife sürümü, sağlayıcı, kart/QR, olay-anı bağlamı bugün nerede kayboluyor?**
HİÇBİR YERDE TUTULMUYOR: `RecordEventInput`/`dto.RecordEventRequest`'te istemci-zaman alanı yok (`CreatedAt` her zaman sunucu `time.Now()`); `resolveBaseFee`/`resolveMeteredRate`/`ResolveActiveStopFare`/`resolveFee`/`hasActiveSubscription`/`checkAccessBlock` TÜMÜ `time.Now()` kullanıyor (repository seviyesinde `at time.Time` parametresi ZATEN var ve tarihe-göre-sorgu DESTEKLENİYOR — sadece hiçbir çağrı site'ı `time.Now()` dışında bir değer geçmiyor). `AccessEvent`'te `RouteFareID`/`RouteStopFareID` FK'sı YOK (sadece sonuç `BaseFee` tutulu yor, hangi tarife dokümanının ürettiği tutulmuyor). `InstrumentUsed` bir ENUM (kart/qr/nakit tipi), kart/QR'nin KİMLİĞİ (hangi fiziksel kart) tutulmuyor. Olayın "kurumu" doğrudan alan değil, `AccessPointID→AccessPoint.InstitutionID` ile HER OKUMADA yeniden türetiliyor (AccessPoint'in kurumu sonra değişirse eski olayların bağlamı sessizce değişir).

**10) Offline gelen olayın sonradan güncel tarifeyle hesaplanmasını engellemek için hangi snapshot/referans alanları gerekli?**
Minimum: (a) `OccurredAt time.Time` — istemci/cihaz zamanı, `CreatedAt`'ten (sunucu alım zamanı) AYRI yeni bir alan; (b) `RouteFareID`/`RouteStopFareID *primitive.ObjectID` + eşleşen `EffectiveFrom`/`EffectiveTo` snapshot'ı; (c) uygulanan tam tarife değerleri (`CardFee`/`CashFee`/`CardRatePerKm` okunan değer, final hesaplanan tutar değil). Mekanizma ZATEN VAR (`ResolveActiveFare(ctx, routeID, at)`, `at` parametresi), sadece `time.Now()` yerine `OccurredAt` geçirilmesi gerekiyor — yeni bir sorgu altyapısı icat etmeye gerek yok.

**11) Açık yolculuk, maksimum tahsilat, nihai ücret ve iade için mevcut mali model yeterli mi; minimum hangi alanlar gerekli?**
Mali MEKANİZMA (Debit/Credit/PostJournal/idempotency) YETERLİ ve DOĞRUDAN yeniden kullanılabilir (soru 6). VERİ MODELİ (`RiderTrip`, `AccessEvent`) YETERSİZ — `RiderTrip.ChargedAmount` TEK bir alan (bugünkü tek-anlı metered ücreti için yeterli), A3/A4'ün "biniş tahsilatı ≠ iniş nihai ücreti ≠ iade tutarı" üçlüsünü AYRI AYRI tutamıyor. Minimum yeni alanlar: `RiderTrip.CollectedAmount`, `CollectedJournalID`, `AdjustmentJournalID`, `RefundAmount`, `FareAlgorithm` (hangi algoritmayla açıldığı, kapanışta değişmemesi için sabitlenmeli); `AccessEvent`'e `OccurredAt`, `RouteFareID`/`RouteStopFareID` (soru 10).

**12) Hangi mevcut API/Flutter/panel akışları etkilenir?**
Backend: `access_points/routes.go` (`RecordEvent`, `SettleTripCash`, `CancelBoardingTrip`, `RecordWalkOnCashEntry`), `dto.go` (`RecordEventRequest`/`AccessEventResponse` — yeni alanlar eklenirse). Panel: `lib/types/access-points.ts` (zaten backend modelinin ALT KÜMESİ, A3/A4 alanları eklemek için genişletilmeli), `lib/api/access-points.ts` (`parseAccessEvent`, rapor parser'ları), `lib/types/fleet-tally.ts`+`lib/api/fleet-tally.ts` (RouteFare/RouteStopFare CRUD — `FareAlgorithm` alanı eklenirse), `lib/api/access-point-discounts.ts` (değişmeyecek, zaten algoritma-bağımsız). Flutter: `kt_surucu_app`/`aypass_erisim_app`'in `validation_service.dart`+`validation_outcome.dart`+`validation_result_screen.dart`+`cash_settlement_screen.dart`+`walk_on_cash_entry_screen.dart` — ÖZELLİKLE A3/A4'ün biniş-anında-maksimum-tutar göstergesi bu ekranlara YENİ bir alan gerektirir (bugün hiç "biniş tahsilatı" göstermiyorlar, sadece nihai `charged_amount`). **Önizleme/tahmin ucu bulunamadı** — rider/sürücü bugün binişte ne kadar tahsil edileceğini GÖREMİYOR; A3/A4 bunu tahsil ANINDA (response'ta) göstermeli, ayrı bir "preview" endpoint'i bu paketin önerdiği minimum kapsamda gerekli DEĞİL.

---

## 3. Doğrulanmış eksikler (bu paket için, P00'ın D-07/D-08/D-09/D-10/D-14'üne ek/çapraz referans)

| # | Eksik | Kanıt | Hangi WP'yi engeller | Durum (2026-10-30) |
|---|---|---|---|---|
| G-1 | Giriş/çıkışta hiçbir mali işlem yok (metered girişi ücretsiz grant) | service.go:1647-1680 | A1/A3/A4 — ikisi de BİNİŞTE tahsilat gerektiriyor | **ÇÖZÜLDÜ A4 için (resolved in WP-3.5, `boardA4`)** — gerçek Debit biniş anında. A1 için WP-3.6'da ele alınacak (A3 henüz aktive edilemiyor). |
| G-2 | `OpenTrip`/`CloseTrip` hataları `_ =` ile yutuluyor, orphan "açık kalan sefer" riski | service.go:1658,1719 | A3/A4 — biniş tahsilatı yapılmış bir seferin kapanamaması mali tutarsızlık yaratır | **KISMEN ÇÖZÜLDÜ A4 için (resolved in WP-3.5)** — `boardA4`/`alightA4` kendi `OpenTripWithCollection`/`CloseTripWithAdjustment` hatalarını ASLA yutmuyor (doğrudan propagate); ek olarak kapanış turu "Debit var, RiderTrip yok" orphan riskini salt-okunur `PreflightOrphanBoardingDebits` ile tespit edilebilir kıldı. F0/A2/metered'ın kendi `_ = s.riderTrips.OpenTrip(...)` deseni BİLİNÇLİ OLARAK değiştirilmedi (sıfır-mali-risk, WP-3.7'nin kapsamı). |
| G-3 | `rider_trips` indeksleri unique değil → paralel çift-açık-sefer | rider_trips/repository/indexes.go:9-32 | A1/A3/A4 | **ÇÖZÜLDÜ (resolved in WP-3.1)** — `fleet_tally_rider_trips__tenant_site_person_open_unique` partial unique index + `PreflightDuplicateOpenTrips` (canlıya almadan önce zorunlu ön kontrol, bkz. R1, hâlâ çalıştırılmadı). |
| G-4 | A4 eksik hücre → sessiz km-fallback | service.go:1694-1706 | A4 (doğrudan kanonik tanıma aykırı) | **ÇÖZÜLDÜ (resolved in WP-3.4)** — `resolveStrictMatrixFare`, fallback yok, `ErrFareMatrixCellMissing`. |
| G-5 | Tarife/indirim/abonelik HER ZAMAN `time.Now()` ile çözülüyor | service.go:1257,1287,1325,1695,1929; rider_trips/service.go:187-227 (`autoCloseOne` de discount'suz, tam-güzergah km ücretiyle) | A1/A3/A4 — offline/gecikmeli olay yanlış tarifeyle fiyatlanır | **ÇÖZÜLDÜ RecordEvent çağrı site'ları için (resolved in WP-3.2)** — `resolvedEventTime`, tüm tarife/indirim çağrılarının tek kaynağı. `autoCloseOne` kısmı hâlâ açık (WP-3.8/G-9'un parçası). |
| G-6 | `AccessEvent`'te `RouteFareID`/`RouteStopFareID`/`OccurredAt` yok | access_event.go (tam alan listesi §A.1'de) | A1–A4 (denetim/itiraz/iade kanıtı) | **KISMEN ÇÖZÜLDÜ, YAKLAŞIM DEĞİŞTİ (resolved in WP-3.1/3.2/3.5):** `OccurredAt` eklendi. `RouteFareID`/`RouteStopFareID` alanları `AccessEvent`'e eklendi ama fiilen DOLDURULMUYOR — gerçek tarihsel kanıt sahibi `RiderTrip.BoardingFareSnapshot`/`AlightingFareSnapshot` oldu (ID+EffectiveFrom'un kaynak kayıt yerinde değişebildiği için yetersiz kalacağı WP-3.5 kapanışında kanıtlandı, bkz. `TestMongoIT_A4_ImmutableSnapshot_SurvivesSourceFareMutation`). |
| G-7 | `RiderTrip`'te biniş-tahsilat/iade alanları yok, tek `ChargedAmount` | rider_trip.go:56-116 | A3/A4 (temel veri modeli boşluğu) | **ÇÖZÜLDÜ A4 için (resolved in WP-3.1/3.5)** — `CollectedAmount`/`CollectedJournalID`/`AdjustmentJournalID`/`RefundAmount`/`BoardingFareSnapshot`/`AlightingFareSnapshot`/idempotency alanları eklendi ve doldurulan. |
| G-8 | Metered giriş, bakiye kontrolü YAPMADAN koşulsuz grant veriyor | service.go:1647-1680 | A3/A4 — yetersiz bakiyede seferin HİÇ açılmaması gerekiyor | **ÇÖZÜLDÜ A4 için (resolved in WP-3.5)** — `boardA4`, Debit başarısız olursa (yetersiz bakiye dahil) sefer hiç açılmıyor, AccessEvent hiç yazılmıyor. |
| G-9 | `autoCloseOne` (TripAssignment tamamlandığında toplu kapama) tam-güzergah km ücretiyle kapatıyor, A2/A4 matrisini hiç kullanmıyor, indirimsiz | rider_trips/service.go:187-227 | A1–A4 — tutarlılık: normal çıkışla auto-close'un aynı algoritmayı kullanması gerekir | **AÇIK, bilinçli olarak bu turda ele alınmadı** — WP-3.8'in kapsamı; model/durum hazırlığı (`RiderTripStatusAutoClosedMissingAlight`) var, davranış implemente değil. |

---

## 4. Yeniden kullanılacak parçalar (yeni yazılmayacak)

- `wallets.Service.Debit/DebitToAccount/Credit/CreditFromAccount` + idempotency-anahtarlı rezervasyon deseni — A3/A4'ün biniş/iniş ikilisi İÇİN HAZIR.
- `ledger.Service.PostJournal`/çift-girişli defter/kısmi-unique idempotency indeksi — reversal "ters journal" deseni zaten `RefundEvent`'te kanıtlanmış, aynen kullanılacak.
- `RouteFare`/`RouteStopFare` modelleri + `EffectiveFrom/To` versiyonlama + `ResolveActiveFare(ctx, id, at)`'in `at`-parametreli tasarımı — SADECE çağrı site'larının `time.Now()` yerine `OccurredAt` geçmesi gerekiyor, model/repository DEĞİŞMEYECEK.
- `RiderTrip` açık/kapalı durum modeli — genişletilecek (yeni alanlar), YENİDEN YAZILMAYACAK.
- `access_point_discounts.ResolveBestDiscount` — algoritma-bağımsız, dokunulmayacak.
- `geo.DistanceMeters` (Haversine) — A1/A3'ün mesafe hesaplaması için aynen kullanılacak.

---

## 5. Önerilen domain modeli

**İlke:** A1–A4'ü TEK fonksiyonda `switch` ile yamamak yerine, güzergâh bazında seçilebilir bir **strateji arayüzü**. Mevcut `chargeMetered`/`chargeWallet` fonksiyonlarının YERİNE geçecek yeni bir "charge orchestrator" değil — onların İÇİNDEKİ algoritma-seçim noktası bir strateji tablosuna çıkarılacak.

```go
// internal/modules/access_points/service/farealgorithm (yeni, küçük paket — mevcut service.go'yu YENİDEN YAZMAZ)

type FareAlgorithm string
const (
    FareAlgorithmF0 FareAlgorithm = "f0"
    FareAlgorithmA1 FareAlgorithm = "a1"
    FareAlgorithmA2 FareAlgorithm = "a2"
    FareAlgorithmA3 FareAlgorithm = "a3"
    FareAlgorithmA4 FareAlgorithm = "a4"
)

// RouteFare'e eklenecek YENİ alan (geriye uyumlu: omitempty, yoksa "a2" varsayılan —
// bugünkü tüm RouteFare kayıtları fiilen A2 davranışı sergiliyor).
// RouteFare.FareAlgorithm FareAlgorithm `bson:"fare_algorithm,omitempty"`

type BoardingDecision struct {
    Grant          bool
    CollectAmount  float64   // A3/A4: maksimum ücret; F0/A1/A2: nihai ücretin ta kendisi (tek-dokunuş)
    RequiresExit   bool      // A1/A3/A4/metered: true (iki dokunuş); F0/A2 tek-dokunuş: false
}

type AlightingDecision struct {
    FinalAmount    float64
    RefundAmount   float64   // CollectAmount - FinalAmount (A3/A4); diğerlerinde 0
    Rejected       bool      // A4: hücre eksik → true, km-fallback YOK
    RejectReason   error
}

type Strategy interface {
    OnBoard(ctx context.Context, in BoardingInput) (BoardingDecision, error)
    OnAlight(ctx context.Context, in AlightingInput) (AlightingDecision, error) // RequiresExit=false ise hiç çağrılmaz
}
```

Dört strateji: `fixedFareStrategy` (F0, tek-dokunuş, mevcut `resolveBaseFee`'nin transport-dışı dalı), `routeFareStrategy` (A2, tek-dokunuş, mevcut `resolveBaseFee`'nin transport dalı), `campusEndpointStrategy` (A1 — **açık tasarım sorusu**, bkz. §9 KR-A1), `distanceOrMatrixStrategy` (A3/A4, iki-dokunuş, parametrik: `FinalFareResolver` olarak mesafe (A3) ya da matris-hücre-zorunlu (A4) enjekte edilir — AYNI strateji kodu, farklı resolver).

**Yeni/değişen alanlar (veri modeli):**
- `RouteFare.FareAlgorithm` (yeni, omitempty, varsayılan A2) — hangi algoritma bu ScheduledRoute+dönem için geçerli.
- `RiderTrip.FareAlgorithm`, `CollectedAmount *float64`, `CollectedJournalID *primitive.ObjectID`, `AdjustmentJournalID *primitive.ObjectID`, `RefundAmount *float64` (yeni).
- `AccessEvent.OccurredAt *time.Time`, `RouteFareID *primitive.ObjectID`, `RouteStopFareID *primitive.ObjectID` (yeni, hepsi omitempty — eski kayıtlar bozulmaz).
- `RecordEventInput`/`dto.RecordEventRequest.OccurredAt *time.Time` (yeni, opsiyonel — yoksa sunucu zamanı kullanılır, DAVRANIŞ DEĞİŞMEZ).

---

## 6. Algoritma bazında durum makinesi

```
F0 / A2  (tek dokunuş):
  [tap] → resolveEventContext → resolveBaseFee(F0:ap.Fee | A2:RouteFare@OccurredAt)
        → resolveFee(discount) → Debit(final) → AccessEvent(granted, final) → [son]

A1  (kampüs sabit uç — iki dokunuş, mesafe tek seferde biniş+iniş ARASINDA hesaplanır):
  [biniş tap] → RiderTrip.Open(entry=biniş noktası) → AccessEvent(granted, ücretsiz grant, "beklemede")
  [kampüs/iniş tap] → distance(entry,exit) → resolveFee(discount) → Debit(final) → RiderTrip.Close + AccessEvent(final)
  NOT: A1 bugünkü `chargeMetered`'ın İSKELETİNİ (open/close) kullanır ama biniş anında TAHSİLAT YAPMAZ (A3/A4'ten farkı budur — kanonik tanımda "net ücret tahsil edilir" tek seferde deniyor, maksimum-tahsil-sonra-iade YOK).

A3 / A4  (iki dokunuş, biniş-tahsil + iniş-ayarlama):
  [biniş tap] → RouteFare'in max ücretini belirle (A3: route'un en uzak O-D'si ya da sabit tavan; A4: matrisin en yüksek hücresi — AÇIK TASARIM SORUSU, bkz. §9 KR-A2)
             → Debit(maxFare, idemKey=board:<riderTripID>)
                 ├─ ErrInsufficientBalance → RiderTrip AÇILMAZ, AccessEvent(denied), [son]
                 └─ başarı → RiderTrip.Open(CollectedAmount=maxFare, CollectedJournalID=...) + AccessEvent(granted, charged=maxFare)
  [iniş tap] → A3: distance(entry,exit)×oran  |  A4: RouteStopFare[entry,exit] (YOKSA → Rejected, km'ye DÜŞME)
             → resolveFee(discount) → finalFee
             → refund = CollectedAmount - finalFee
             → Credit(refund, idemKey=adjust:<riderTripID>)  [refund<=0 ise atla, finalFee>maxFare teorik olarak olmamalı]
             → RiderTrip.Close(AdjustmentJournalID, RefundAmount) + AccessEvent(final, charged=finalFee, refunded=refund)
```

---

## 7. Veri/migration etkisi

Tüm yeni alanlar **omitempty/opsiyonel** — mevcut kayıtlar bozulmaz, geriye dönük migration GEREKMEZ (eski `RouteFare` kayıtları `FareAlgorithm` boşsa A2 varsayılır, eski `AccessEvent`/`RiderTrip` kayıtları yeni alanları hiç taşımaz, okuma kodu nil-safe olmalı). Tek gerçek "migration" kararı: hangi mevcut `ScheduledRoute`'ların A1/A3/A4 olarak işaretleneceği — bu bir VERİ GİRİŞİ kararı (panel/seed), şema migration'ı değil. Index değişikliği: `rider_trips`'e `(tenant_id, site_id, person_id, status)` üzerinde UNIQUE PARTIAL index eklenmesi önerilir (G-3) — bu geriye dönük veri denetimi gerektirir (zaten iki açık kaydı olan kişi varsa index oluşturma başarısız olur); WP paketinde ayrı denetim adımı olarak planlanmalı.

---

## 8. Test matrisi (özet, uygulama paketlerinde detaylandırılacak)

| Algoritma | Senaryo | Beklenen |
|---|---|---|
| F0 | İndirimli/indirimsiz sabit ücret | mevcut testler zaten kapsıyor, DOKUNULMAYACAK |
| A1 | Biniş→kampüs, kampüs→iniş, mesafe×oran, indirim | yeni tablo-güdümlü test |
| A2 | Bilinen O-D sabit ücret | mevcut `RouteFare`/`RouteStopFare` testleri zaten kapsıyor |
| A3 | Biniş-max-tahsil → yetersiz bakiye (sefer açılmaz) → iniş-nihai → iade | yeni, Mongo'lu (atomiklik) |
| A3 | İniş olmadan `CloseOpenTripsForAssignment` (auto-close) tetiklenirse ne olur | yeni — G-9 ile bağlı, karar gerekli |
| A4 | Bilinen hücre → doğru iade; eksik hücre → RED (km'ye düşmez) | yeni, kanonik tanımın doğrudan testi |
| Hepsi | `OccurredAt` geçmişte bir tarih → o tarihteki `RouteFare` versiyonu kullanılır | yeni, tarih-versiyonlama testi |
| Hepsi | İki eşzamanlı biniş isteği aynı kişi için → tek açık sefer | yeni, Mongo'lu eşzamanlılık (G-3 ile bağlı) |
| A3/A4 | `OpenTrip` biniş-tahsilattan SONRA başarısız olursa | yeni — mali tutarlılık: tahsilat yapıldı ama sefer kaydı yok senaryosu AÇIKÇA ele alınmalı (G-2 ile bağlı) |

---

## 9. Küçük, bağımsız implementasyon iş paketleri (önerilen sıra)

| Paket | İçerik | Bağımlılık | Durum (2026-10-30) |
|---|---|---|---|
| **WP-3.1** | Veri modeli: `RiderTrip`/`AccessEvent`/`RouteFare` yeni alanları (hepsi omitempty) + `rider_trips` unique partial index (önce veri denetimi) | — | **TAMAMLANDI** |
| **WP-3.2** | `OccurredAt` ucu: `RecordEventInput`/DTO'ya opsiyonel alan, tüm `ResolveActiveFare`/`ResolveActiveStopFare`/discount/`hasActiveSubscription` çağrı site'larını `time.Now()` yerine `OccurredAt ?? time.Now()` kullanacak şekilde güncelle | WP-3.1 | **TAMAMLANDI** |
| **WP-3.3** | `Strategy` arayüzü + F0/A2 stratejileri (mevcut davranışı AYNEN yeniden ifade eder, fonksiyonel DEĞİŞİKLİK yok — refactor-only, regresyon testiyle kanıtlanır) | WP-3.1 | **TAMAMLANDI** (`resolveFareAlgorithm`, tek algoritma-seçim sahibi; F0/A2/metered davranışı değişmedi) |
| **WP-3.4** | A4 "eksik hücre → red" kararı (G-4) — SADECE A4 stratejisinde, A2/metered'ın mevcut fallback'ine DOKUNMADAN | WP-3.3 | **TAMAMLANDI** (`resolveStrictMatrixFare`, fallback yok) |
| **WP-3.5** | A3/A4 biniş-tahsil/iniş-ayarlama akışı (yetersiz bakiyede sefer açılmaz, G-8) | WP-3.1, WP-3.3, KR-A2 kararı | **TAMAMLANDI** — kod + hedefli doğrulama + kapanış turu tamamlandı; merge/yayın kararı bekliyor. Son SHA **`2764898`** (dal `feat/wp3-fare-algorithm-foundation`). Kanıtlar: immutable mali snapshot testi (`TestMongoIT_A4_ImmutableSnapshot_SurvivesSourceFareMutation`, gerçek Mongo), process-restart alight-retry testi (`TestMongoIT_A4_AlightRetryAfterAccessEventFailure_ResumesAfterRestart`, gerçek Mongo, çift restart), orphan boarding-Debit preflight testi (`TestMongoIT_PreflightOrphanBoardingDebits_*`, gerçek Mongo), A4 maksimum-validasyon testleri (`TestCreateRejectsA4When*`/`TestCreateSucceedsForA4When*`, `route_fares/service/service_test.go`), A3 aktivasyon guard testi (`TestCreateRejectsA3EvenWithValidMaxGrossFare`). Ayrıca WP-15 AccessPoint hattıyla yerel entegrasyon dalında (`integration/aypass-accesspoint-fare-engine`, merge commit `4dc8c7a`) birleştirildi, güvenlik+mali+regresyon kabul testleriyle doğrulandı. |
| **WP-3.6** | A1 kampüs-sabit-uç stratejisi | WP-3.3, KR-A1 kararı | **SIRADAKİ İŞ** — entegrasyon provası yeşil sonuçla onaylandı, üretim kodlamasına geçiliyor |
| **WP-3.7** | `OpenTrip`/`CloseTrip` hata-yutma davranışının A3/A4 için sıkılaştırılması (G-2) — orphan-tahsilat senaryosu | WP-3.5 | **BEKLİYOR** — WP-3.5 kapanışı bu maddenin yalnız BİR dilimini (A4 boarding Debit→RiderTrip orphan'ı, salt-okunur preflight ile) ele aldı; A4'ün kendi `boardA4`/`alightA4` kod yolu zaten hata yutmuyor (G-2 A4 için fiilen giderilmiş durumda), ama genel `OpenTrip`/`CloseTrip` (F0/A2/metered'ın `_ = ...` deseni) hâlâ sıkılaştırılmadı — bu paketin asıl kapsamı budur, dokunulmadı |
| **WP-3.8** | `autoCloseOne`'ın yeni stratejilerle tutarlılığı (G-9) | WP-3.3, WP-3.5 | **BEKLİYOR** — auto-close/late-alight modelde durum (`RiderTripStatusAutoClosedMissingAlight`/`RiderTripStatusAdjustedAfterLateAlight`) hazır ama davranış implemente değil (bilinçli, bu turda kapsam dışı) |
| **WP-3.9** | Panel/Flutter: yeni alanların gösterimi (biniş-tahsil tutarı, iade bildirimi) | WP-3.5, WP-3.6 | **BEKLİYOR** |

---

## 10. Açık kararlar (önce karar, sonra kod)

- **KR-A1**: A1 tek-dokunuş mu (biniş noktası + sabit kampüs konumu arası mesafe önceden bilinir, TEK tap yeterli) yoksa iki-dokunuş mu (metered gibi biniş+iniş tap'i)? Kanonik tanım "sabah biniş→kampüs, akşam kampüs→iniş" diyor — bu, YÖN BAŞINA tek sabit nokta (kampüs) olduğunu, dolayısıyla TEK dokunuşla (diğer nokta zaten bilinen kampüs) mesafe hesaplanabileceğini ima ediyor. Bu paket bunu VARSAYMADI, açık soru olarak bırakıyor.
  **ÇÖZÜLDÜ (resolved in konuşma kararı, 2026-10-30, WP-3.5 onayı sırasında):** A1 **tek dokunuşlu**dur; her yön kendi bağımsız tek-tap'i. Gidiş ve dönüş aynı trip değildir. A3/A4'ün açık-RiderTrip/boarding-alighting state machine'i A1'e uygulanmayacak — WP-3.6 teslimi bu kararla kodlanacak.
- **KR-A2**: A3'te "biniş maksimum ücreti" neyin maksimumu — route'un tüm O-D'leri arasındaki en yüksek `RouteStopFare`'i mi, yoksa `RouteFare`'in ayrı bir "tavan" alanı mı? A4'te ise matrisin o biniş-noktasından ERİŞİLEBİLİR tüm hücrelerinin maksimumu mu? Model/alan eklenmesi gerekebilir.
  **ÇÖZÜLDÜ (resolved in konuşma kararı + WP-3.5 kodu, commit `2a6a051`/`8060901`):** P, önceden doğrulanmış, kayıt-zamanında sabitlenmiş bir tavandır — asla anlık hesaplanmaz. A4 için `MaxGrossFareCard/Cash >= route'un TÜM aktif matris hücrelerinin maksimumu` kayıt/aktivasyon anında zorunlu (`route_fares.Service.validate`, `RouteStopFareLister.ListActiveForRoute`). A3 için eşdeğer bir maksimum-mesafe validator'ı HENÜZ YOK — bu yüzden A3 `ErrA3MaximumDistancePolicyNotConfigured` ile aktivasyonda koşulsuz reddediliyor (WP-3.5 kapanış maddesi 6); A3 enum/modelde var ama aktif kullanılabilir değil.
- **KR-A3**: `autoCloseOne`'ın A3/A4 seferlerini nasıl kapatacağı (G-9) — sefer tamamlanmadan sefer sonlanırsa biniş-tahsilatının ne olacağı mali bir karar, kod kararı değil.
  **ÇÖZÜLDÜ (resolved in konuşma kararı, 2026-10-30):** Auto-close, P'yi nihai tutar olarak tutar, ayrı bir durumla (`RiderTripStatusAutoClosedMissingAlight`) işaretler — otomatik iade/ayarlama YOK. Gecikmeli-iniş (late-alight) uzlaşması kendi penceresi/sonucuyla AYRI bir akıştır ama bu round'da İMPLEMENTE EDİLMEDİ — yalnız model/durum hazırlığı (`RiderTripStatusAdjustedAfterLateAlight`) var. Gerçek davranış **WP-3.8**'e aittir, hâlâ bekliyor.

---

## 11. Bağımlılık olarak kaydedilenler (bu paketin kapsamı dışı)

Firma devri (WP-11), Global Market (WP-12), POS entegrasyonu (WP-9), offline cihaz kuyruğu (WP-7), raporlama geliştirmesi — bu paketteki yeni alanlar (`OccurredAt`, `FareAlgorithm` vb.) bu işlerin HER BİRİNİ etkileyecek (offline kuyruk `OccurredAt`'i dolduracak istemci; raporlama yeni `CollectedAmount`/`RefundAmount` alanlarını gösterecek) — sadece bağımlılık olarak not edildi, bu paketin parçası değil.

---

## 12. Açık yayın engelleri (2026-10-30 eklendi — WP-3.5 kapanışı bunları KAPATMADI, kapatılmış gibi raporlanmayacak)

- **R1 — Unique index preflight:** `internal/platform/indexes.EnsureAll`, DuplicateKey çakışmasında kritik partial unique index'i (G-3, `fleet_tally_rider_trips__tenant_site_person_open_unique`) **sessizce atlayabiliyor** (hata döndürmüyor, sadece log uyarısı yazıyor). `rider_trips.PreflightDuplicateOpenTrips`, canlı yayın/migration prosedürünün **zorunlu ön kontrolü** olarak işaretli kalıyor — bu turda `EnsureAll`'ın genel davranışı değiştirilmedi.
- **R2 — A3:** A3, enum/model ve ortak mali altyapıda (Debit/Credit/snapshot/idempotency) bulunabilir. Ancak `ErrA3MaximumDistancePolicyNotConfigured` guard'ı (WP-3.5 kapanış maddesi 6, `route_fares.Service.validate`) kaldırılmayacak — A3 hiçbir route'ta aktif algoritma olarak kullanılamaz, kendi maksimum-mesafe validatorü yazılana kadar.
- **R3 — Auto-close / late-alight:** `RiderTripStatusAutoClosedMissingAlight` ve `RiderTripStatusAdjustedAfterLateAlight` durumlarının MODELDE bulunması, bu özelliklerin UYGULANMIŞ olduğu anlamına gelmiyor. Gerçek davranış implementasyonu **WP-3.8**'e ait, bu turda da yapılmadı.
- **R4 — Production verification:** Gerçek üretim verisinde duplicate-open-trip preflight (`PreflightDuplicateOpenTrips`) **çalıştırılmadı** — bu sandbox'ta production erişimi yok. Yalnız disposable/sentetik yerel MongoDB'de test edildi. Canlıya alınmadan önce gerçek veride ayrıca çalıştırılması R1'in bir parçasıdır.
