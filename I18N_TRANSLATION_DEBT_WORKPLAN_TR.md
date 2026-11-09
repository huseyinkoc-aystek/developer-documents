# I18N Çeviri Borcu — İş Planı

**Sahibi:** Hüseyin Koç · **Oluşturma:** 2026-09-21 · **Kaynak:** `done/TASK_059_SITE_LANGUAGES.md` (üçüncü tur)
**Durum:** WP-1 tamamlandı (2026-09-22) · WP-2 / WP-3 planlandı, başlanmadı

README standardı: her kullanıcıya görünen metin bir çeviri **anahtarı** olmalı ve
`translation_seed.go` içinde **28 dilde** tohumlanmalı. Eski (legacy) alanlarda bu sağlanmıyor;
aşağıdaki iş paketleri borcu ölçülebilir ve geri dönüşsüz şekilde kapatır.

## 1. Ölçüm (2026-09-21)

| Ölçüt | Değer | Nasıl ölçülür |
|---|---|---|
| `languages.*` anahtarı, 19 dilde Türkçe metin | **62 anahtar / 1200 boşluk** | `TRANSLATION_REPORT_PREFIX=languages. go test ./internal/seeds -run TestTranslationCoverageReport -v` |
| Frontend'de kullanılan ama hiçbir seed'de olmayan anahtar | **2259 anahtar** | `cd aystek-core-next && npm run check-i18n` (baseline: `scripts/i18n-baseline.json`) |
| Herhangi bir dilde eksik / Türkçe kopya / (Latin dışı dilde) İngilizce kopya | **8307 anahtar / 106.230 boşluk** | `TRANSLATION_REPORT_PREFIX='*' go test ./internal/seeds -run TestTranslationCoverageReport -v` |
| Ağır: ≥10 dilde Türkçe yedek | **1592 anahtar / 37.230 boşluk** | seed edilmiş DB'de `values[lang] === values.tr` sayımı |

Aktif diller (seed'de `is_active`): **tr, en, de, fr, ar, fa, zh, ja, pt, es, ru**. Diğer 17 dil kayıtlı ama pasif —
**öncelik aktif dillerdedir**. WP-1'de aktif olup eksik olanlar: **zh, ja, fa, ar, ru**.

Bu dokümandaki sayılar, kabul kapıları (Bölüm 4) çalıştırılınca yeniden doğrulanır; yanlışsa doküman güncellenir.

## 2. İş paketleri

### WP-1 — Dil Yönetimi sayfası çevirileri (`languages.*`) · Öncelik: YÜKSEK · Boyut: ~1200 metin

Dil yönetim sayfasının kendisi çoklu dilli yönetimin giriş noktasıdır; 19 dilde Türkçe görünmesi kabul edilemez.

- **Kapsam:** 62 anahtar × (zh, ja, fa, ar, kk, uz, ru, uk, pl, hr, cs, bs, sr, sq, hy, az-cyrl, ky, tk, mn) = 1178, ayrıca `az` (16), `languages.delete.message` için de/fr/it/es/pt/pt-br (6).
- **Kaynak metinler:** `I18N_WP1_LANGUAGES_SOURCE_TEXTS_TR.md`
- **Desen:** yeni overlay dosyası `internal/seeds/languages_admin_translations_seed.go` (`site_language_translations_seed.go` ile aynı yaklaşım); eski `translation_seed.go` satırlarına dokunulmaz.
- **Öncelik sırası:** önce aktif diller (zh, ja, fa, ar, ru), sonra kalan 14 dil.
- **Bitti tanımı:** `fullyTranslatedPrefixes` listesine `"languages."` eklenmiş ve `go test ./internal/seeds` yeşil.

### WP-2 — Seed'siz anahtarlar (`check-i18n` baseline) · Öncelik: ORTA · Boyut: 2259 anahtar × 28 dil

Bu anahtarlar kodda `t("anahtar", "varsayılan metin")` ile kullanılıyor ama hiçbir seed'de yok; Türkçe dışındaki dillerde
yalnızca kod içi varsayılan görünür. Dağılım (baseline'ın ilk öneki): telematics 495, fleet_tally 350, content_ingestion 84,
ops 62, common 60, chatbot 58, wallets 58, currencies 53, activityLogs 46, access_points 45, persons 45, institutions 42,
access_point_discounts 40, feed 38 …

- **Yöntem:** modül önekini modül önekine kapat: anahtarları koddaki varsayılan metinden (tr) çıkar, yeni bir overlay/seed dosyasında 28 dile çevir, `npm run check-i18n -- --update-baseline` ile baseline'ı küçült.
- **Bitti tanımı (paket başına):** o önek baseline'dan tamamen çıkmış, `fullyTranslatedPrefixes` listesine eklenmiş.
- **Kural:** yeni kod baseline'a **yeni** anahtar ekleyemez (`npm run check-i18n` hata verir) — borç büyümez.

**`common` alt-paketi tamamlandı (2026-09-22, Claude):** 60 eksik `common.*` anahtarından 3'ü
(`confirm_delete`, `reset_filters`, `no_records`) zaten var olan, tam çevrili camelCase eşdeğerleriyle
(`confirmDelete`, `resetFilters`, `noRecords`) aynıydı — yeni seed eklemek yerine 23 dosyadaki
kullanım noktaları bu mevcut anahtarlara yönlendirildi (aystek-core-next). Kalan 57 anahtar
translation_seed.go'ya 28 dilde tam çevrili literal map girdisi olarak eklendi (overlay değil —
bu anahtarlar hiç yoktu, "mevcut satıra dokunma" kısıtı geçerli değildi). Ayrıca yönlendirilen
`common.noRecords`'un kendi 19 dilde Türkçe-kopya hatası (WP-3 tipi borç) fırsat bu fırsattayken
düzeltildi. Baseline 2259 → 2199 (-60).
**ÖNEMLİ BULGU:** `"common."` öneki `fullyTranslatedPrefixes`'e EKLENMEDİ — önekte halihazırda
77 var olan anahtarda 971 boşluk var (çoğu `allLangs()` ile İngilizce doldurulmuş, klasik WP-3
borcu). Bu, WP-2'nin "eksik anahtar" kapsamının çok dışında; ayrı bir WP-3 artımı olarak ele
alınmalı. `common.*` prefix'i WP-3'ün öncelik listesine eklenmeli (bkz. Bölüm 2, WP-3).

**`currencies` alt-paketi tamamlandı (2026-09-22, Claude):** 53 eksik `currencies.*` anahtarı
(`app/(dashboard)/admin/currencies/page.tsx`) 28 dilde tam çevrili literal map girdisi olarak
eklendi — tek dosyada kullanılıyordu, yinelenen (near-duplicate) anahtar bulunmadı. Baseline
2199 → 2146 (-53). `"currencies."` öneki `fullyTranslatedPrefixes`'e EKLENMEDİ — pre-existing
`currencies.fx.*` alt-öneğinde 468 boşluk / 39 anahtar var (farklı özellik: kur listesi,
`allLangs`/literal-EN ile seed edilmiş); bu WP-3'ün öncelik listesine eklendi (bkz. Bölüm 2, WP-3).

**`persons` alt-paketi tamamlandı (2026-09-22, Claude):** 45 eksik `persons.*` anahtarı
(`app/(dashboard)/admin/persons/page.tsx` ve `.../persons/[id]/page.tsx`) 28 dilde tam çevrili
literal map girdisi olarak eklendi, yinelenen anahtar bulunmadı. Doğrulamada 17 "same as Turkish"
yanlış-pozitifi çıktı — hepsi ya Azerice-Türkçe Oğuz Türkçesi ortak kökü (`Bağla`, `Ad Soyad`,
`Var`, `Silindi`) ya da uluslararası ödünç kelime (`Telefon`: az/bs/cs/de/hr/pl/tk/uz; `Rol`:
az/es/tk/uz) — `identicalToTurkishAllowed`'a gerekçeli olarak eklendi. Baseline 2146 → 2101 (-45).
`"persons."` öneki `fullyTranslatedPrefixes`'e EKLENMEDİ — önekte halihazırda 13 var olan
anahtarda 156 boşluk var (`persons.affiliation.*`, `persons.role.*`, `persons.status.*`,
`persons.field.route`; `allLangs`/literal-EN ile seed edilmiş); bu WP-3'ün öncelik listesine
eklendi (bkz. Bölüm 2, WP-3).

**`institutions` alt-paketi tamamlandı (2026-09-22, Claude):** 42 eksik `institutions.*` anahtarı
(`app/(dashboard)/admin/institutions/*.tsx`) 28 dilde tam çevrili literal map girdisi olarak
eklendi, yinelenen anahtar bulunmadı. 6 anahtar (`confirm_delete_message`, `field.order`,
`field.status`, `notice.created`, `notice.deleted`, `notice.updated`) ile `filters.type`/`modal.new`
mevcut `persons.*`/`common.*`/kendi `field.type` çevirisiyle aynı metni yeniden kullandı (aynı
kavram, farklı ad alanı — yinelenen anahtar sorunu değil). Doğrulamada 8 yanlış-pozitif çıktı:
5'i Azerice-Türkçe Oğuz Türkçesi ortak kökü (`Ad`, `Sıra`, `Silindi`) ve uluslararası ödünç kelime
`Tip` (az/bs/hr) — `identicalToTurkishAllowed`'a gerekçeli olarak eklendi; ayrıca `field.slug`'ın
kk/az-cyrl değerleri yanlışlıkla bare-Latin "Slug" bırakılmıştı (diğer Kiril dilleriyle tutarsız),
"Слаг" olarak düzeltildi. Baseline 2101 → 2059 (-42). `"institutions."` öneki
`fullyTranslatedPrefixes`'e EKLENMEDİ — önekte halihazırda 5 var olan anahtarda (`status.active`,
`status.inactive`, `type.customer_company`, `type.institution`, `type.supplier`) 60 boşluk var
(`allLangs`/literal-EN ile seed edilmiş); bu WP-3'ün öncelik listesine eklendi (bkz. Bölüm 2, WP-3).

**`wallets` alt-paketi tamamlandı (2026-09-22, Claude):** 58 eksik `wallets.*` anahtarı
(`app/(dashboard)/admin/wallets/*.tsx`) 28 dilde tam çevrili literal map girdisi olarak eklendi,
yinelenen anahtar bulunmadı. 4 anahtar (`field.source`, `col.status`, `filter.status`,
`modal.new`) mevcut `persons.*`/kendi `btn.new` çevirisiyle aynı metni yeniden kullandı (aynı
kavram, farklı ad alanı). Doğrulamada 7 yanlış-pozitif çıktı: hepsi Azerice/Türkmence-Türkçe
Oğuz Türkçesi ortak kökü veya uluslararası ödünç kelime (`Cüzdan`, `Kart`/`QR`, `Blokla`) —
`identicalToTurkishAllowed`'a gerekçeli olarak eklendi; ayrıca `reconcile.ledger`'ın az-cyrl
değeri yanlışlıkla bare-Latin "Ledger" bırakılmıştı (diğer Kiril dilleriyle tutarsız), "Дәфтәр"
olarak düzeltildi. Baseline 2059 → 2001 (-58). `"wallets."` öneki `fullyTranslatedPrefixes`'e
EKLENMEDİ — önekte halihazırda 2 var olan anahtarda (`status.active`, `status.frozen`) 24 boşluk
var (`allLangs`/literal-EN ile seed edilmiş); bu WP-3'ün öncelik listesine eklendi (bkz. Bölüm 2,
WP-3).

**`telematics` (495 anahtar) alt-modül-bazlı işleniyor** — tek seferde ele alınamayacak kadar
büyük olduğu için `telematics.*` alt-önekleri kendi WP-2 modülü gibi ayrı ayrı kapatılıyor:
`settings` 29, `can_decoders` 35, `alerts` 39, `driver_sessions` 42, `commands` 47,
`vehicle360` 61, `alert_rules` 81, `dashboard` 161 (toplam 495). Sıra en küçükten en büyüğe.

**`telematics.settings` alt-paketi tamamlandı (2026-09-22, Claude):** 29 eksik
`telematics.settings.*` anahtarı (`app/(dashboard)/admin/telematics/settings/page.tsx`) 28 dilde
tam çevrili literal map girdisi olarak eklendi, yinelenen anahtar bulunmadı. **Not:** bu modülün
frontend varsayılan metni Türkçe değil İngilizce yazılmış (`t(key, "İngilizce metin")`) — bu
WP-2'de ilk kez görülen bir durum; Türkçe değerler frontend varsayılanından kopyalanmadı, gerçek
çeviri olarak yazıldı. Doğrulamada 2 gap çıktı: `teltonika_tcp.tenant_id`'nin az-cyrl/hy
değerleri yanlışlıkla bare-Latin "Tenant ID" bırakılmıştı (diğer dillerle tutarsız), "Тенант ID"
/ "Վարձակալի ID" olarak düzeltildi. Baseline 2001 → 1972 (-29). `"telematics.settings."` öneği
`fullyTranslatedPrefixes`'e eklendi (önekte pre-existing boşluk yoktu, tüm anahtarlar yeni).

**`telematics.can_decoders` alt-paketi tamamlandı (2026-09-22, Claude):** 35 eksik
`telematics.can_decoders.*` anahtarı (`app/(dashboard)/admin/telematics/can-decoders/page.tsx`)
28 dilde tam çevrili literal map girdisi olarak eklendi, yinelenen anahtar bulunmadı. Bu modülün
de frontend varsayılanı İngilizce yazılmıştı (settings ile aynı desen), Türkçe değerler gerçek
çeviri olarak yazıldı. Doğrulamada 9 yanlış-pozitif çıktı: hepsi Azerice-Türkçe Oğuz Türkçesi
ortak kökü (`Son Test`, `Ad`, `Yük JSON`) veya uluslararası ödünç kelime `Protokol`
(az/bs/cs/hr/tk/uz) — `identicalToTurkishAllowed`'a gerekçeli olarak eklendi. Baseline
1972 → 1937 (-35). `"telematics.can_decoders."` öneki `fullyTranslatedPrefixes`'e eklendi
(pre-existing boşluk yoktu, tüm anahtarlar yeni).

**`telematics.alerts` alt-paketi tamamlandı (2026-09-22, Claude):** 39 eksik
`telematics.alerts.*` anahtarı (`app/(dashboard)/admin/telematics/alerts/page.tsx`) 28 dilde
tam çevrili literal map girdisi olarak eklendi, yinelenen anahtar bulunmadı. Bu modülün de
frontend varsayılanı İngilizce yazılmıştı, Türkçe değerler gerçek çeviri olarak yazıldı.
Doğrulamada 3 yanlış-pozitif çıktı (Azerice-Türkçe ortak kökü/ödünç kelime: `field.device`,
`field.message`, `state.overdue`) — gerekçeli olarak eklendi. Baseline 1937 → 1898 (-39).
**ÖNEMLİ BULGU:** `"telematics.alerts."` öneki `fullyTranslatedPrefixes`'e EKLENMEDİ — önekte
halihazırda var olan tek anahtar `telematics.alerts.notice.acknowledged`'da 11 boşluk var
(`allLangs` ile seed edilmiş, klasik WP-3 borcu). Bu bulgu WP-3'ün öncelik listesine eklendi.

**`telematics.driver_sessions` alt-paketi tamamlandı (2026-09-22, Claude):** 42 eksik
`telematics.driver_sessions.*` anahtarı (`app/(dashboard)/admin/telematics/driver-sessions/page.tsx`
ve `dashboard`/`vehicles/[vehicleId]` sayfalarındaki kullanımlar) 28 dilde tam çevrili literal
map girdisi olarak eklendi, yinelenen anahtar bulunmadı, önek tamamen yeniydi (pre-existing
boşluk yoktu). Doğrulamada 6 yanlış-pozitif çıktı: 4'ü Azerice-Türkçe Oğuz Türkçesi ortak kökü
(`completed`, `device_id`, `driver_name`, `started`) — gerekçeli eklendi; 2'si (`identifier`,
`identifier_type`) ise gerçek kök benzerliği değil, tembel kopyalamaydı — Azerice değerleri
"İdentifikator"/"İdentifikator Tipi" (gerçek Azerice loanword) olarak düzeltildi. Baseline
1898 → 1856 (-42). `"telematics.driver_sessions."` öneki `fullyTranslatedPrefixes`'e eklendi.

**`telematics.commands` alt-paketi tamamlandı (2026-09-22, Claude):** 47 eksik
`telematics.commands.*` anahtarı (`app/(dashboard)/admin/telematics/commands/*.tsx` ve
dashboard/settings/activity-logs sayfalarındaki kullanımlar) 28 dilde tam çevrili literal map
girdisi olarak eklendi, yinelenen anahtar bulunmadı. `device_id`/`vehicle_id` mevcut
`telematics.driver_sessions.*` çevirisini aynen yeniden kullandı (aynı kavram, farklı ad alanı).
Doğrulamada 5 yanlış-pozitif çıktı: Azerice-Türkçe ortak kökü/uluslararası ödünç kelime
(`device_id`, `payload`, `field.policy` × bs/hr/sq) — gerekçeli eklendi. Ayrıca `field.sla`'nın
11 Latin-dışı dildeki değeri yanlışlıkla bare-Latin "SLA" bırakılmıştı (İngilizce ile birebir
aynı) — her biri için gerçek çeviri/native ifade yazıldı (ör. ru "Соглашение об уровне
обслуживания", zh "服务级别协议", ja "サービスレベル合意"). Baseline 1856 → 1809 (-47).
**ÖNEMLİ BULGU:** `"telematics.commands."` öneki `fullyTranslatedPrefixes`'e EKLENMEDİ — önekte
halihazırda var olan 7 `notice.*` anahtarında (approved, created, dispatched,
execution_sync_processed, execution_synced, rejected, result_updated) 77 boşluk var
(`allLangs` ile seed edilmiş, klasik WP-3 borcu). Bu bulgu WP-3'ün öncelik listesine eklendi.

**`telematics.vehicle360` alt-paketi tamamlandı (2026-09-22, Claude):** 61 eksik
`telematics.vehicle360.*` anahtarı (`app/(dashboard)/admin/telematics/vehicles/[vehicleId]/page.tsx`
ve driver-sessions sayfasındaki kullanımlar) 28 dilde tam çevrili literal map girdisi olarak
eklendi, yinelenen anahtar bulunmadı, önek tamamen yeniydi (pre-existing boşluk yoktu).
Doğrulamada 8 yanlış-pozitif çıktı: 3'ü Azerice-Türkçe ortak kökü (`can_metric.last_test`,
`replay.at`, `slot.driver_started`) — gerekçeli eklendi; 3'ü (`slot.can_obd`, `slot.driver`,
`slot.driver_identity`) belirsiz kök benzerliğiydi — Azerice değerleri gerçek/daha net
ifadelerle ("Slotu", "Şəxsiyyət") düzeltildi; ayrıca `can_metric.rpm`'in az-cyrl/ja değerleri
yanlışlıkla bare-Latin "RPM" bırakılmıştı (İngilizce ile birebir aynı) — gerçek çeviri/native
ifadeyle düzeltildi (ja "回転数", az-cyrl "Фырланма Сүрәти"). Baseline 1809 → 1748 (-61).
`"telematics.vehicle360."` öneki `fullyTranslatedPrefixes`'e eklendi.

**`telematics.alert_rules` alt-paketi tamamlandı (2026-09-22, Claude):** 81 eksik
`telematics.alert_rules.*` anahtarı (`app/(dashboard)/admin/telematics/alert-rules/*.tsx` ve
dashboard sayfasındaki kullanımlar) 28 dilde tam çevrili literal map girdisi olarak eklendi,
yinelenen anahtar bulunmadı. Doğrulamada 7 yanlış-pozitif çıktı: `name|az` (Azerice-Türkçe ortak
kökü) ve `code` (az/bs/hr/pl/tk/uz — uluslararası ödünç kelime, `languages.table.code` ile aynı) —
gerekçeli olarak eklendi. Baseline 1748 → 1667 (-81). **ÖNEMLİ BULGU:**
`"telematics.alert_rules."` öneki `fullyTranslatedPrefixes`'e EKLENMEDİ — önekte halihazırda var
olan 2 `notice.*` anahtarında (`deleted`, `saved`) 22 boşluk var (`allLangs` ile seed edilmiş,
klasik WP-3 borcu, yeni eklenen 81 anahtardan tamamen bağımsız). Bu bulgu WP-3'ün öncelik
listesine eklendi.

**`telematics.dashboard` alt-paketi tamamlandı (2026-09-22, Claude):** 161 eksik
`telematics.dashboard.*` anahtarı (`app/(dashboard)/admin/telematics/dashboard/page.tsx`) 28
dilde tam çevrili literal map girdisi olarak eklendi, yinelenen anahtar bulunmadı — bu, WP-2'de
tek seferde eklenen en büyük anahtar kümesiydi (4 alt-parti halinde üretildi: `kpi.*` 41,
`metric.*` 21, `section.*` 40, kalan 59). Doğrulamada 1 yanlış-pozitif çıktı:
`partner_governance.sla_overdue_short`'un 11 Latin-dışı dildeki değeri yanlışlıkla bare-Latin
"SLA" bırakılmıştı — `telematics.commands.field.sla` ile aynı gerçek çeviri/native ifadeler
kullanılarak düzeltildi. Baseline 1667 → 1506 (-161). **ÖNEMLİ BULGU:**
`"telematics.dashboard."` öneki `fullyTranslatedPrefixes`'e EKLENMEDİ — önekte halihazırda var
olan tek anahtar `telematics.dashboard.ready`'de 11 boşluk var (`allLangs` ile seed edilmiş,
klasik WP-3 borcu, yeni eklenen 161 anahtardan tamamen bağımsız). Bu bulgu WP-3'ün öncelik
listesine eklendi.

**🎯 KİLOMETRE TAŞI: `telematics.*` (495 anahtar) TAMAMLANDI (2026-09-22, Claude).** Sekiz
alt-paket (`settings` 29, `can_decoders` 35, `alerts` 39, `driver_sessions` 42, `commands` 47,
`vehicle360` 61, `alert_rules` 81, `dashboard` 161 = 495) tek tek kapatıldı. Sıradaki ve WP-2'nin
son büyük modülü: `fleet_tally` (350 anahtar, 15 alt-önek: `route_definitions` 39,
`scheduled_routes` 38, `supplier_invoices` 33, `trip_assignments` 33, `connectors` 28,
`gps_providers` 28, `schedules` 26, `vehicle_alarms` 26, `control_types` 24,
`vehicle_mappings` 21, `drivers` 20, `verifications` 14, `movement_reports` 10,
`route_fares` 5, `runs` 5 — en küçükten başlanıyor).

**`fleet_tally.route_fares` + `fleet_tally.runs` alt-paketleri tamamlandı (2026-09-22,
Claude):** 10 eksik anahtar (5+5) (`app/(dashboard)/admin/fleet-tally/scheduled-routes/[id]/page.tsx`
ve `app/(dashboard)/admin/fleet-tally/runs/*.tsx`) 28 dilde tam çevrili literal map girdisi
olarak eklendi, yinelenen anahtar bulunmadı. Doğrulamada 1 yanlış-pozitif çıktı
(`runs.col.timing|az` — Azerice-Türkçe ortak kökü, "Zaman") — gerekçeli olarak eklendi.
Baseline 1506 → 1496 (-10). Her iki önek de `fullyTranslatedPrefixes`'e EKLENMEDİ — her ikisinde
de çok sayıda pre-existing `allLangs()` anahtarı var (`route_fares`: badge/btn/col/field grubu;
`runs`: action/col/detail/job/status grubu, çok büyük — ayrı bir WP-3 taraması gerektirir).

**`fleet_tally.movement_reports` alt-paketi tamamlandı (2026-09-22, Claude):** 10 eksik
`fleet_tally.movement_reports.*` anahtarı (`app/(dashboard)/admin/fleet-tally/movement-reports/page.tsx`)
28 dilde tam çevrili literal map girdisi olarak eklendi, yinelenen anahtar bulunmadı, önek
tamamen yeniydi (pre-existing boşluk yoktu). Baseline 1496 → 1486 (-10).
`"fleet_tally.movement_reports."` öneki `fullyTranslatedPrefixes`'e eklendi.

**`fleet_tally.verifications` alt-paketi tamamlandı (2026-09-22, Claude):** 14 eksik
`fleet_tally.verifications.*` anahtarı (`app/(dashboard)/admin/fleet-tally/verifications/page.tsx`)
28 dilde tam çevrili literal map girdisi olarak eklendi, yinelenen anahtar bulunmadı.
`btn.retry`/`search_placeholder` mevcut `retry.title`/`fleet_tally.movement_reports.search_placeholder`
çevirisini aynen yeniden kullandı. Doğrulamada iki gap kendi eklediğim anahtarlarda çıktı ve
düzeltildi: `col.push`'un kk/ky/mn/ru/sr/uk değerleri ile `col.gps_provider`'ın 12 Latin-dışı
dildeki değerleri yanlışlıkla bare-Latin "Push"/"GPS" bırakılmıştı (İngilizce ile birebir aynı) —
gerçek çeviri/native ifadeyle düzeltildi (ör. `col.push` ru "Пуш", `col.gps_provider` ar
"نظام GPS", zh "GPS定位"). Baseline 1486 → 1472 (-14). `"fleet_tally.verifications."` öneki
`fullyTranslatedPrefixes`'e EKLENMEDİ — önekte halihazırda 7 var olan anahtarda (`col.checked_at`,
`col.job_date`, `col.location_count`, `col.plate`, `col.schedule`, `col.status`, `title`) boşluk
var (`allLangs` ile seed edilmiş, klasik WP-3 borcu, yeni eklenen 14 anahtardan tamamen
bağımsız).

**`fleet_tally.drivers` alt-paketi tamamlandı (2026-09-22, Claude):** 20 eksik
`fleet_tally.drivers.*` anahtarı (`app/(dashboard)/admin/fleet-tally/drivers/*.tsx`) 28 dilde
tam çevrili literal map girdisi olarak eklendi, yinelenen anahtar bulunmadı. 6 anahtar
(`col.status`, `field.status`, `notice.created`, `notice.deleted`, `notice.updated`,
`section.basics`) mevcut `persons.*`/`institutions.*` çevirisini aynen yeniden kullandı.
Doğrulamada 1 yanlış-pozitif çıktı (`notice.deleted|az` — Azerice-Türkçe ortak kökü, aynı
`persons.notice.deleted|az` gerekçesi) — gerekçeli olarak eklendi. Baseline 1472 → 1452 (-20).
`"fleet_tally.drivers."` öneki `fullyTranslatedPrefixes`'e eklendi (pre-existing boşluk yoktu).

**`fleet_tally.vehicle_mappings` alt-paketi tamamlandı (2026-09-22, Claude):** 21 eksik
`fleet_tally.vehicle_mappings.*` anahtarı (`app/(dashboard)/admin/fleet-tally/vehicle-mappings/page.tsx`)
28 dilde tam çevrili literal map girdisi olarak eklendi, yinelenen anahtar bulunmadı (mevcut
`col.plate`/`col.external_id`/`col.note`/`title` gibi anahtarlar aynı sayfada tablo başlığı
olarak farklı bir amaçla kullanılıyor, gerçek çakışma değil). `modal.new` mevcut `btn.new`,
`notice.created`/`deleted`/`updated` mevcut `persons.*`/`fleet_tally.drivers.notice.*`
çevirisini aynen yeniden kullandı. Doğrulamada 3 gap çıktı: `notice.deleted|az` (Azerice-Türkçe
ortak kökü) ve `source.manual|fr` ("Manuel" — Fransızca ve Türkçede aynı yazım) gerekçeli olarak
eklendi; `source.stock`'un 12 Latin-dışı dildeki değeri yanlışlıkla bare-Latin "Stock" (modül
adı) bırakılmıştı — gerçek transliterasyonla düzeltildi (ör. ja "ストック", ru "Сток", ar
"ستوك"). Baseline 1452 → 1431 (-21). `"fleet_tally.vehicle_mappings."`
öneki `fullyTranslatedPrefixes`'e EKLENMEDİ — önekte halihazırda 11 var olan anahtarda (`add`,
`col.external_id`, `col.note`, `col.plate`, `edit`, `external_id_hint`, `import_csv`,
`import_success`, `provider`, `title`, `validation.required_fields`) boşluk var (`allLangs` ile
seed edilmiş, klasik WP-3 borcu, yeni eklenen 21 anahtardan tamamen bağımsız).

**`fleet_tally.control_types` alt-paketi tamamlandı (2026-09-22, Claude):** 24 eksik
`fleet_tally.control_types.*` anahtarı (`app/(dashboard)/admin/fleet-tally/control-types/*.tsx`)
28 dilde tam çevrili literal map girdisi olarak eklendi, yinelenen anahtar bulunmadı. `modal.new`
mevcut `btn.new`, `col.status`/`field.status`/`notice.created`/`notice.deleted`/`notice.updated`/
`section.basics` mevcut `persons.*`/`institutions.*` çevirisini aynen yeniden kullandı.
Doğrulamada 2 yanlış-pozitif çıktı (`field.name|az` — "Kod" uluslararası ödünç kelime + "Adı"
ortak kök; `notice.deleted|az` — aynı `persons.notice.deleted|az` gerekçesi) — gerekçeli olarak
eklendi. Baseline 1431 → 1407 (-24). `"fleet_tally.control_types."` öneki
`fullyTranslatedPrefixes`'e EKLENMEDİ — önekte halihazırda 6 var olan anahtarda (`add`,
`col.active`, `col.algorithm`, `col.name`, `edit`, `title`) boşluk var (`allLangs` ile seed
edilmiş, klasik WP-3 borcu, yeni eklenen 24 anahtardan tamamen bağımsız).

**`fleet_tally.schedules` alt-paketi tamamlandı (2026-09-22, Claude):** 26 eksik
`fleet_tally.schedules.*` anahtarı (`app/(dashboard)/admin/fleet-tally/schedules/*.tsx`) 28
dilde tam çevrili literal map girdisi olarak eklendi, yinelenen anahtar bulunmadı. `modal.new`
mevcut `btn.new`, `col.status`/`field.is_active`/`notice.*`/`search_placeholder` mevcut
`persons.*`/`fleet_tally.control_types.*` çevirisini aynen yeniden kullandı. Doğrulamada 6 gap
çıktı: 2'si gerçek hata idi ve düzeltildi (`col.config`'in az değeri "Ayarlar" tembelce Türkçeden
kopyalanmıştı, gerçek Azerice "Tənzimləmələr" ile değiştirildi; `col.lookback`'in az-cyrl/kk/ky/mn
değerleri yanlışlıkla bare-Latin "Lookback" bırakılmıştı, fonetik transliterasyonla düzeltildi:
kk/az-cyrl/ky "Лукбек", mn "Лукбэк"); 4'ü Azerice-Türkçe ortak kökü idi (`col.name`,
`field.connector`, `field.lookback_days`, `notice.deleted`) — gerekçeli olarak eklendi. Baseline
1407 → 1381 (-26). `"fleet_tally.schedules."` öneki `fullyTranslatedPrefixes`'e EKLENMEDİ —
önekte halihazırda 10 var olan anahtarda (`add`, `col.active`, `col.connector`, `col.cron_expr`,
`col.display_name`, `col.last_run`, `col.lookback_days`, `edit`, `run_now`, `title`) boşluk var
(`allLangs` ile seed edilmiş, klasik WP-3 borcu, yeni eklenen 26 anahtardan tamamen bağımsız).

**`fleet_tally.vehicle_alarms` alt-paketi tamamlandı (2026-09-22, Claude):** 26 eksik
`fleet_tally.vehicle_alarms.*` anahtarı (`app/(dashboard)/admin/fleet-tally/vehicle-alarms/*.tsx`)
28 dilde tam çevrili literal map girdisi olarak eklendi, yinelenen anahtar bulunmadı, önek
tamamen yeniydi (pre-existing boşluk yoktu). `field.alarm_time` mevcut `col.alarm_time`,
`col.status`/`field.is_active`/`notice.*`/`search_placeholder`/`section.basics` mevcut
`persons.*`/`fleet_tally.schedules.*`/`fleet_tally.movement_reports.*` çevirisini aynen
yeniden kullandı. Doğrulamada 1 yanlış-pozitif çıktı (`notice.deleted|az` — aynı
`persons.notice.deleted|az` gerekçesi) — gerekçeli olarak eklendi. Baseline 1381 → 1355 (-26).
`"fleet_tally.vehicle_alarms."` öneki `fullyTranslatedPrefixes`'e eklendi.

**`fleet_tally.connectors` alt-paketi tamamlandı (2026-09-23, Claude):** 28 eksik
`fleet_tally.connectors.*` anahtarı (`app/(dashboard)/admin/fleet-tally/connectors/*.tsx`) 28
dilde tam çevrili literal map girdisi olarak eklendi, yinelenen anahtar bulunmadı. 6 anahtar
(`field.is_active`, `notice.created`/`notice.deleted`/`notice.updated`, `search_placeholder`,
`modal.new`) mevcut `persons.*`/kendi `btn.new` çevirisini aynen yeniden kullandı. Doğrulamada
3 yanlış-pozitif çıktı — hepsi gerekçeli olarak eklendi: `notice.deleted|az` (aynı
`persons.notice.deleted|az` gerekçesi), `btn.new|az`/`modal.new|az` ("Yeni Bağlayıcı" — "Yeni"
ve "Bağlayıcı" ikisi de standart Oğuz Türkçesi ortak kökenli, aynı `fleet_tally.schedules.field.connector|az`
gerekçesi). Baseline 1355 → 1327 (-28). `"fleet_tally.connectors."` öneki `fullyTranslatedPrefixes`'e
EKLENMEDİ — önekte halihazırda var olan 11 anahtardan 4'ünde (`title`, `edit`, `col.display_name`,
`validation.integration_connector_*` — 4 anahtar) boşluk var (`allLangs` ile seed edilmiş, klasik
WP-3 borcu, yeni eklenen 28 anahtardan tamamen bağımsız).

**`fleet_tally.gps_providers` alt-paketi tamamlandı (2026-09-23, Claude):** 28 eksik
`fleet_tally.gps_providers.*` anahtarı (`app/(dashboard)/admin/fleet-tally/gps-providers/*.tsx`) 28
dilde tam çevrili literal map girdisi olarak eklendi, yinelenen anahtar bulunmadı. 16 anahtar
(`col.credentials`, `col.status`, `col.type`, `credentials.configured`/`edit_help_set`/`edit_help_unset`/`missing`,
`field.display_name`, `field.is_active`, `notice.created`/`deleted`/`updated`, `search_placeholder`,
`section.advanced`/`basics`/`credentials`) mevcut `fleet_tally.connectors.*`/`persons.*` çevirisini aynen
yeniden kullandı (`modal.new` da kendi `btn.new`'ini). Doğrulamada 1 yanlış-pozitif çıktı
(`notice.deleted|az` — aynı `persons.notice.deleted|az` gerekçesi) — gerekçeli olarak eklendi.
Ayrıca bu artımda geçmiş `fleet_tally.connectors.*` ekleme betiğinin son satırında (son anahtar
`section.credentials` ile sonraki mevcut `fleet_tally.runs.detail.col.job_date` anahtarı arasında)
eksik satır sonu karakterinden kaynaklanan bir biçim hatası fark edildi ve düzeltildi (Go sözdizimi
geçerliydi, derleme/testler zaten geçiyordu — salt biçim tutarlılığı düzeltmesi). Baseline: 1327 → 1299.
`"fleet_tally.gps_providers."` öneki `fullyTranslatedPrefixes`'e EKLENMEDİ — önekte halihazırda var
olan 12 anahtardan 9'unda (`add`, `col.active`, `col.base_url`, `col.display_name`, `col.provider_type`,
`credentials.hint`, `edit`, `notice.deleted`, `title`) boşluk var (`allLangs` ile seed edilmiş, klasik
WP-3 borcu, yeni eklenen 28 anahtardan tamamen bağımsız).

**`fleet_tally.supplier_invoices` alt-paketi tamamlandı (2026-09-23, Claude):** 33 eksik
`fleet_tally.supplier_invoices.*` anahtarı (`app/(dashboard)/admin/fleet-tally/supplier-invoices/*.tsx`)
28 dilde tam çevrili literal map girdisi olarak eklendi, yinelenen anahtar bulunmadı, önek tamamen
yeniydi (pre-existing boşluk yoktu). 4 anahtar (`col.amount`, `col.invoice_number`, `col.status`,
`col.supplier`) kendi eşdeğer `field.*`/`filters.status` anahtarlarını dahili olarak yeniden kullandı.
Doğrulamada 0 yanlış-pozitif/boşluk çıktı — tamamen temiz. Baseline: 1299 → 1266.
`"fleet_tally.supplier_invoices."` öneki `fullyTranslatedPrefixes`'e eklendi (önek tamamen yeniydi,
pre-existing anahtar yoktu).

**`fleet_tally.trip_assignments` alt-paketi tamamlandı (2026-09-23, Claude):** 33 eksik
`fleet_tally.trip_assignments.*` anahtarı (`app/(dashboard)/admin/fleet-tally/trip-assignments/page.tsx`)
28 dilde tam çevrili literal map girdisi olarak eklendi, yinelenen anahtar yoktu (1 anahtar —
`col.plate` — kendi `field.plate` çevirisini dahili yeniden kullandı). Doğrulamada 4 yanlış-pozitif
çıktı — hepsi gerekçeli olarak eklendi: `btn.end|az` ("Bitir" — "bitirmək" ortak Oğuz Türkçesi kökü),
`col.timing|az` ("Zaman" — aynı `fleet_tally.runs.col.timing|az` gerekçesi), `source.manual|fr`
("Manuel" — aynı `fleet_tally.vehicle_mappings.source.manual|fr` gerekçesi), `status.completed|az`
("Tamamlandı" — "tamamlamaq" ortak Arapça-ödünç+Türkçe-ek kökü). Baseline: 1266 → 1233.
`"fleet_tally.trip_assignments."` öneki EKLENMEDİ — önekte halihazırda var olan 8 anahtarda
(`btn.plan`, `col.planned_date`, `field.planned_date`, `field.scheduled_route`, `modal.plan`,
`notice.planned`, `source.planned`, `status.planned`) `allLangs` kaynaklı WP-3 borcu var, yeni
eklenen 33 anahtardan tamamen bağımsız.

**`fleet_tally.scheduled_routes` alt-paketi tamamlandı (2026-09-23, Claude):** 38 eksik
`fleet_tally.scheduled_routes.*` anahtarı (`app/(dashboard)/admin/fleet-tally/scheduled-routes/*.tsx`)
28 dilde tam çevrili literal map girdisi olarak eklendi, yinelenen anahtar yoktu. `modal.new`=`btn.new`
ve `section.schedule`=`col.schedule` dahili olarak yeniden kullandı; `notice.created`/`deleted`/`updated`
mevcut `persons.notice.*` çevirisini yeniden kullandı. Doğrulamada 6 sorun çıktı: 2 gerçek hata
düzeltildi (`col.schedule`/`section.schedule` ve `field.schedule_entry`'nin az değeri tembelce
Türkçe "Zamanlama" kopyalanmıştı — Azerbaycan dilinde bu bağlamda doğru kelime "Cədvəl" olduğundan
gerçek çeviriyle düzeltildi), 3 yanlış-pozitif gerekçeli olarak eklendi (`notice.deleted|az` — aynı
`persons.notice.deleted|az` gerekçesi; `schedule_modal.entry|az`/`|tk` — "Giriş" ortak Oğuz Türkçesi
kelimesi). Baseline: 1233 → 1195. `"fleet_tally.scheduled_routes."` öneki EKLENMEDİ — önekte
halihazırda var olan 18 anahtarda (`btn.clone`, `col.validity`, `col.visible`,
`customer_company_hint`, `description_create`, `edit_hint`, `fares_hint`, `field.customer_company`,
`field.effective_from`, `field.effective_to`, `field.visible_to_riders`, `modal.clone`, `no_end_date`,
`section.fares`, `section.validity`, `validity_hint`, `visible_no`, `visible_yes`) `allLangs` kaynaklı
WP-3 borcu var, yeni eklenen 38 anahtardan tamamen bağımsız.

**`fleet_tally.route_definitions` alt-paketi tamamlandı (2026-09-23, Claude):** 39 eksik
`fleet_tally.route_definitions.*` anahtarı (`app/(dashboard)/admin/fleet-tally/route-definitions/*.tsx`)
28 dilde tam çevrili literal map girdisi olarak eklendi, yinelenen anahtar yoktu, önek tamamen
yeniydi (pre-existing anahtar yoktu). `modal.new`=`btn.new` dahili, `col.name`=`field.name` ve
`notice.created`/`deleted`/`updated` mevcut `persons.notice.*` çevirisini yeniden kullandı.
Doğrulamada 2 sorun çıktı: `subtitle`'ın az-cyrl değeri (modül adı "Fleet Tally", marka adı olarak
İngilizce/Latin bırakıldı) İngilizce ile birebir eşleştiğinden Latin-dışı yazı sistemi kuralına
takıldı — diğer Latin-dışı dillerle tutarlı şekilde "Fleet Tally модулу" olarak düzeltildi;
`notice.deleted|az` yanlış-pozitifi aynı `persons.notice.deleted|az` gerekçesiyle eklendi.
Baseline: 1195 → 1156. `"fleet_tally.route_definitions."` öneki `fullyTranslatedPrefixes`'e eklendi
(önek tamamen yeniydi, pre-existing borç yoktu).

**🎯 KİLOMETRE TAŞI: `fleet_tally.*` (350 anahtar, 15 alt-paket) TAMAMLANDI (2026-09-23, Claude).**
Tüm alt-paketler tek tek kapatıldı: `route_fares` 5, `runs` 5, `movement_reports` 10,
`verifications` 14, `drivers` 20, `vehicle_mappings` 21, `control_types` 24, `schedules` 26,
`vehicle_alarms` 26, `connectors` 28, `gps_providers` 28, `supplier_invoices` 33,
`trip_assignments` 33, `scheduled_routes` 38, `route_definitions` 39 = 350. Bununla birlikte
WP-2'nin şu anda listelenen tüm seed'siz anahtar backlog'u (telematics.* 495 + fleet_tally.* 350
= 845 anahtar, 28 dilde tam çeviri) tamamlanmış oldu.

### WP-3 — Türkçe / İngilizce yedeğe düşmüş mevcut anahtarlar · Öncelik: DÜŞÜK-ORTA · Boyut: 1592 ağır anahtar (37.230 metin)

**✅ `common.*` TAMAMLANDI (2026-09-23, Claude).** 76 anahtardaki 952 boşluk (76'sı `TRANSLATION_REPORT_PREFIX=common.` ile bulundu; 3'ü hy'de Türkçe yedek olarak kalan `confirmDelete`/`deleteConfirmation`/`deleteConfirmationMultiple`, `copyright` 19 dilde Türkçe yedek, kalan 72 anahtarın çoğu literal-EN veya literal-TR ile "sahte" seed edilmişti — her dile aynı İngilizce/Türkçe kelime kopyalanmıştı) tek tek düzeltildi. 72 anahtar tam 28-dil literal olarak yeniden yazıldı (mevcut doğru çeviriler korunarak sadece hatalı diller hedeflendi 4 özel anahtar için: confirmDelete/deleteConfirmation/deleteConfirmationMultiple/copyright). Doğrulamada ek 40 boşluk çıktı (ilk düzeltme turunda bazı yeni çeviriler de eksik/yanlış-pozitifti): 1 gerçek hata düzeltildi (`bulk_delete|az` tembel Türkçe kopyaydı, gerçek Azerice "Kütləvi sil" ile düzeltildi), `notAvailable`/`seo`/`sku`/`slug` anahtarlarının bazı Latin-dışı dillerdeki bare-acronym tutarsızlıkları düzeltildi (ja/zh/ru/uk/kk/ky/mn/sr/az-cyrl için doğal ifadeler veya parantezli açıklamalar eklendi), 11 yanlış-pozitif gerekçelendirildi (Azerice/İspanyolca/Türkmence/İtalyanca ortak Oğuz Türkçesi kökleri ve uluslararası ödünç kelimeler). **Doğrulama sonrası `common.` öneki TAMAMEN temiz çıktı** (sadece bu 76 anahtar değil, önek altındaki TÜM anahtarlar) — `"common."` `fullyTranslatedPrefixes`'e eklendi. **⚠️ Önemli teknik ders:** ilk düzeltme denemesinde Python regex tabanlı "brace-matching" replace scripti `allLangs(...)` fonksiyon-çağrısı formatındaki anahtarlarda (`{` yerine `(` ile başlayan) yanlışlıkla dosyanın çok ilerisindeki bir sonraki `{` karakterine kadar arama yaptı ve 1011 satırlık istenmeyen silme yaptı — `git checkout` ile geri alındı ve script hem `{...}` hem `allLangs(...)` biçimlerini ayrı ayrı algılayacak şekilde düzeltildi (bkz. `feedback` niteliğinde not: WP-3'te var olan anahtarları toplu değiştirirken her zaman `git diff --stat` ile beklenmeyen büyük silme/ekleme oranı kontrol edilmeli).

**✅ `currencies.*` TAMAMLANDI (2026-09-23, Claude).** Başlangıçta `currencies.fx.*` (468/39) olarak
işaretlenmişti, ancak tam `TRANSLATION_REPORT_PREFIX=currencies.` taraması 45 anahtarda 491 boşluk
ortaya çıkardı (fx. alt-öneği dışında `currencies.form.code/name`, `currencies.fxSync.*`,
`currencies.fx_rate_*` de dahil). Tüm 45 anahtar 28 dilde tam çevrildi (çoğu literal-EN ile "sahte"
seed edilmişti). Doğrulamada ek 40 boşluk çıktı: 10 yanlış-pozitif gerekçelendirildi (Azerice/Fransızca
ortak Oğuz Türkçesi kökleri/uluslararası ödünç kelimeler: "Ad", "Sabit", "Manuel", "Kod", "Statik"),
geri kalanı `source.api`/`fxSync.provider.http`/`fxSync.form.headers` anahtarlarının Latin-dışı
dillerdeki bare-acronym/İngilizce-loanword tutarsızlıkları — doğal ifadelerle düzeltildi (örn. ja
"API経由", ru "Через HTTP API", kk "Тақырып JSON"). Doğrulama sonrası `currencies.` öneki TAMAMEN
temiz — `"currencies."` `fullyTranslatedPrefixes`'e eklendi.

**✅ `persons.*` TAMAMLANDI (2026-09-23, Claude).** 13 anahtarda (`affiliation.role.*` 4,
`affiliation.status.*` 3, `field.route`, `role.alumni`/`role.corporate`, `status.active`/`merged`/`suspended`)
tespit edilen 156 boşluk düzeltildi — tümü `allLangs(tr,en)` ile seed edilmiş, 28 dilde tam çevrildi.
Doğrulamada 1 yanlış-pozitif çıktı (`affiliation.role.staff|pl` — "Personel" Fransızca kökenli
uluslararası ödünç kelime, Lehçede de aynen kullanılıyor) — gerekçeli olarak eklendi. Doğrulama sonrası
`persons.` öneki TAMAMEN temiz (daha önceki WP-2 artımında eklenen 45 anahtarla birlikte) —
`"persons."` `fullyTranslatedPrefixes`'e eklendi.

**✅ `institutions.*` TAMAMLANDI (2026-09-23, Claude).** 5 anahtarda (`status.active`/`inactive`,
`type.customer_company`/`institution`/`supplier`) tespit edilen 60 boşluk düzeltildi — tümü
`allLangs(tr,en)` ile seed edilmiş, 28 dilde tam çevrildi, yanlış-pozitif/allowlist gerekmedi
(ilk doğrulama turunda 0 boşluk). Doğrulama sonrası `institutions.` öneki TAMAMEN temiz (daha önceki
WP-2 artımında eklenen 42 anahtarla birlikte) — `"institutions."` `fullyTranslatedPrefixes`'e eklendi.

**✅ `wallets.*` TAMAMLANDI (2026-09-23, Claude).** 2 anahtarda (`status.active`, `status.frozen`)
tespit edilen 24 boşluk düzeltildi — `allLangs(tr,en)` ile seed edilmiş, 28 dilde tam çevrildi,
0 yanlış-pozitif. Doğrulama sonrası `wallets.` öneki TAMAMEN temiz (daha önceki WP-2 artımında
eklenen 58 anahtarla birlikte) — `"wallets."` `fullyTranslatedPrefixes`'e eklendi.

**🎯 KİLOMETRE TAŞI:** WP-2 sırasında keşfedilen tüm büyük-önek boşlukları (common, currencies,
persons, institutions, wallets — toplam 1703 boşluk/141 anahtar) tamamlandı.

**✅ Ertelenmiş telematics/fleet_tally tekil-anahtar boşlukları TAMAMLANDI (2026-09-23, Claude).**
8 önek (`telematics.alerts.`, `telematics.commands.`, `telematics.alert_rules.`, `telematics.dashboard.`,
`fleet_tally.connectors.`, `fleet_tally.gps_providers.`, `fleet_tally.trip_assignments.`,
`fleet_tally.scheduled_routes.`) tarandı — belgede tahmin edilen ~30 anahtar yerine gerçekte 56 anahtarda
672 boşluk bulundu (fleet_tally.connectors./scheduled_routes. önceden düşünülenden daha büyük çıktı).
Tüm 56 anahtar 28 dilde tam çevrildi (çoğu `allLangs(tr,en)` ile seed edilmişti). Doğrulamada 5
yanlış-pozitif gerekçelendirildi (Azerice/Türkmence ortak Oğuz Türkçesi kökleri ve uluslararası ödünç
kelimeler: "Bağlayıcı Tipi", "Klonla", "Görünürlük", "Gizli"). Doğrulama sonrası tüm 8 önek TAMAMEN
temiz — hepsi `fullyTranslatedPrefixes`'e eklendi.

**✅ `plugins.*` TAMAMLANDI (2026-09-23, Claude).** Doküman 121 anahtar tahmin etmişti; gerçek tarama
158 anahtarda 2486 boşluk buldu. 156 anahtar 28 dilde tam çevrildi (4 alt-parti halinde üretildi — büyük
modüller için kurulan yöntemle). Doğrulamada ek 62 boşluk çıktı: 4 yanlış-pozitif gerekçelendirildi
(Azerice ortak Oğuz Türkçesi kökleri ve "Hook" gibi paylaşılan teknik terimler), 2 gerçek hata düzeltildi
(`form.settings`/`form.settings_invalid`'in az değeri tembelce Türkçe "Ayarlar" kopyalanmıştı — Azerbaycan
dilinde doğru kelime "Tənzimləmələr" olduğundan düzeltildi, aynı `fleet_tally.scheduled_routes.col.config`
hatası), `form.sha` (SHA256) ve `form.websitePlaceholder` (https://example.com) 9-10 Latin-dışı dilde
bare-acronym/literal-URL tutarsızlığı doğal ifade/yerelleştirilmiş örnek alan adlarıyla düzeltildi.
**Bilinçli istisna:** `form.config_schemaPlaceholder`, `form.settingsPlaceholder` (JSON örnek metni
`{"fields":[]}`/`{"enabled":true}`) ve `form.entry_filePlaceholder` (dosya adı örneği `plugin.so`) —
3 anahtar, 35 boşluk — kasıtlı olarak çevrilmedi: bunlar kod/dosya adı örnekleridir, doğal dil değil,
her dilde aynı kalmaları doğrudur. Bu nedenle `"plugins."` `fullyTranslatedPrefixes`'e EKLENMEDİ
(test'in "Latin dışı dilde İngilizce metin" sezgiseli kod örneklerini de yanlışlıkla işaretliyor,
allowlist mekanizması bu durum için tasarlanmamış — test mantığını değiştirmek bu düzeltmenin kapsamı
dışında bırakıldı).

**⚠️ `integrations.*` BEKLENENDEN ÇOK DAHA BÜYÜK ÇIKTI (2026-09-23 keşfi).** Doküman 84 anahtar tahmin
etmişti; gerçek tarama **1232 anahtar, 14.784 boşluk** buldu (14 kat büyük) — 36 alt modüle yayılıyor
(en büyükleri: `api_sharing.*` 2316 boşluk, `telematics.*` entegrasyon alt-sistemi 2280, `dashboard.*`
1524, `shared_services.*`/`chains.*` 876'şar, `fields.*` 780, `pull_schedules.*` 768,
`runtime_alert_policy.*` 732, `test.*` 720, `validation.*` 600, kalan 26 alt modül ~3780). Bu, bugüne
kadar tamamlanan en büyük paketin (`plugins.*`, 156 anahtar) ~8 katı — ayrı, kendi başına büyük bir
görev olarak ele alınmalı. Kullanıcı onayıyla ŞİMDİLİK ATLANDI, sıradaki daha küçük öneklerle devam
ediliyor; `integrations.*` listenin sonuna, kendi büyüklüğüyle orantılı ayrı bir çalışma bloğu olarak
not düşüldü.

**✅ `contracts.*` TAMAMLANDI (2026-09-23, Claude).** Doküman 82 anahtar tahmin etmişti; gerçek tarama
83 anahtarda 2051 boşluk buldu — tr/en/de doğruydu, kalan 25 dilin tamamı tembelce Türkçe kopyalanmıştı.
83 anahtarın tamamı 28 dilde tam çevrildi (3 alt-parti halinde üretildi). Doğrulamada 3 boşluk çıktı:
2 yanlış-pozitif gerekçelendirildi (`contents.form.language|az`/`|tk` — "Dil" ortak Oğuz Türkçesi kelimesi),
1 gerçek hata düzeltildi (`table.end|az`'ın "Bitiş" değeri standart Azerbaycanca değildi — "Son" ile
düzeltildi). Doğrulama sonrası `contracts.` öneki TAMAMEN temiz — `"contracts."` `fullyTranslatedPrefixes`'e
eklendi.

**✅ `menus.*` TAMAMLANDI (2026-09-23, Claude).** Doküman 81 anahtar tahmin etmişti; gerçek tarama
81 anahtarda 1523 boşluk buldu (tahmin doğruydu). 81 anahtarın tamamı 28 dilde tam çevrildi (3
alt-parti halinde üretildi). Doğrulamada ek boşluklar çıktı: `form.url`/`table.url` (URL) ve
`form.rel` (Rel) ve `form.seoSection` (SEO) anahtarlarının Latin-dışı dillerdeki bare-acronym
tutarsızlıkları doğal ifadelerle düzeltildi (ör. ja "URLアドレス", ru "Атрибут Rel", kk "SEO
баптаулары"), 14 yanlış-pozitif gerekçelendirildi (Azerice/Türkmence ortak Oğuz Türkçesi kökleri
ve uluslararası ödünç kelimeler: "Sil", "ARIA Etiketi", "başlığı", "Sıra", "-- Kök --", "SEO
Başlığı"). **Bilinçli istisna:** `form.typePlaceholder` (literal kod-değeri örneği
`backend, frontend, footer...`) — 1 anahtar, 10 boşluk — kasıtlı olarak çevrilmedi (aynı
`plugins.form.entry_filePlaceholder` mantığı: kod örneği, doğal dil değil). Bu nedenle `"menus."`
`fullyTranslatedPrefixes`'e EKLENMEDİ.

**✅ `posts.*` TAMAMLANDI (2026-09-23, Claude).** Doküman 72 anahtar tahmin etmişti; gerçek tarama
140 anahtarda 1591 boşluk buldu (~2 kat büyük). 140 anahtarın tamamı 28 dilde tam çevrildi (5
alt-parti halinde üretildi). Doğrulamada ek 24 boşluk çıktı: `editor.html` (HTML), `form.slug`/`slug`
(Slug) anahtarlarının Latin-dışı dillerdeki bare-acronym tutarsızlıkları doğal ifadelerle düzeltildi
(ör. ja "HTMLコード", kk "Мекенжай атауы (Slug)"), `form.seoTabTitle`'ın az değerindeki "Ayarlar"
tembel Türkçe kopyası "Tənzimləmələri" ile düzeltildi (aynı önceki `plugins.form.settings` hatası),
5 yanlış-pozitif gerekçelendirildi (Azerice/Türkmence ortak Oghuz Türkçesi kökleri: "Sil", "Diller",
"SEO Başlığı"). Doğrulama sonrası `posts.` öneki TAMAMEN temiz — `"posts."` `fullyTranslatedPrefixes`'e
eklendi.

**✅ `marketplaces.*` TAMAMLANDI (2026-09-23, Claude).** Doküman 68 anahtar tahmin etmişti; gerçek
tarama da 68 anahtarda 1768 boşluk buldu (tahmin doğruydu) — tr/en dışında hiçbir dil seed edilmemişti
(literal kopya değil, tamamen eksik). 68 anahtarın tamamı 28 dilde tam çevrildi (3 alt-parti halinde
üretildi; `provider.*` alt-önekindeki marka adları — Amazon, eBay, Etsy, Shopify, Trendyol,
Hepsiburada, N11 — Latin-dışı 12 dil için fonetik transliterasyonla işlendi, GPS-sağlayıcı/Stock
emsaliyle tutarlı). Doğrulamada ek 22 boşluk çıktı: `provider.n11`'in bazı Latin-dışı dillerde bare
kalması ve `provider.*`/`form.webhookSecret`'ın az-cyrl'de (marka adları PROTECTED listesinde
korunduğu için) İngilizce ile birebir kalması doğal ifade/fonetik transliterasyonla düzeltildi,
6 yanlış-pozitif gerekçelendirildi (Azerice ortak Oghuz Türkçesi kökleri ve Farsça kökenli ödünç
kelime: "Ad", "Mağaza", "Alındı"). Doğrulama sonrası `marketplaces.` öneki TAMAMEN temiz —
`"marketplaces."` `fullyTranslatedPrefixes`'e eklendi.

**✅ `orders.*` TAMAMLANDI (2026-09-23, Claude).** Doküman 60 anahtar tahmin etmişti; gerçek tarama
63 anahtarda 1566 boşluk buldu (yakın tahmin). 63 anahtarın tamamı 28 dilde tam çevrildi (2 alt-parti
halinde üretildi). Doğrulamada 8 boşluk çıktı, tümü yanlış-pozitif olarak gerekçelendirildi (Azerice/
Boşnakça/Hırvatça ortak Oghuz Türkçesi kökleri ve uluslararası ödünç kelimeler: "Aç", "Tip",
"Tamamlandı", "Vergi" — çoğu daha önceki modüllerde aynı gerekçeyle zaten kabul edilmiş kalıplar).
Doğrulama sonrası `orders.` öneki TAMAMEN temiz — `"orders."` `fullyTranslatedPrefixes`'e eklendi.

**✅ `dashboard.*` TAMAMLANDI (2026-09-23, Claude).** Doküman 56 anahtar tahmin etmişti; gerçek tarama
64 anahtarda 1384 boşluk buldu (yakın tahmin — `scope.*` alt-öneği rol-temelli görünürlük UI'ı için
uzun açıklama metinleri içeriyordu). 64 anahtarın tamamı 28 dilde tam çevrildi (4 alt-parti halinde
üretildi). Doğrulamada ek 15 boşluk çıktı: `performance.p95` (API P95) 12 Latin-dışı dilde bare kaldı
— "gecikme/latency" anlamı eklenerek düzeltildi (ör. kk "API P95 кідірісі"), `title`'ın kk/ky değeri
bare "Dashboard" kalmıştı — fonetik "Дашборд" ile düzeltildi, 1 yanlış-pozitif gerekçelendirildi
(`performance.title|az` — "Performans" uluslararası ödünç kelime). Doğrulama sonrası `dashboard.`
öneki TAMAMEN temiz — `"dashboard."` `fullyTranslatedPrefixes`'e eklendi.

**✅ `payment_links.*` TAMAMLANDI (2026-09-23, Claude).** Doküman 54 anahtar tahmin etmişti; gerçek
tarama da 54 anahtarda 1366 boşluk buldu (tahmin tam doğruydu). 54 anahtarın tamamı 28 dilde tam
çevrildi (2 alt-parti halinde üretildi). Doğrulamada ek 13 boşluk çıktı: `form.paymentIntentId`/
`public.intent` (Payment Intent — Stripe'a özgü teknik terim) 5 Latin-dışı dilde bare kaldı, parantezli
açıklamayla düzeltildi (ör. kk "Payment Intent (төлем ниеті)"), `public.token`/`table.token`'ın zh
değeri bare "Token" kalmıştı — "令牌（Token）" ile düzeltildi, 1 yanlış-pozitif gerekçelendirildi
(`public.open3ds|az` — "Doğrulama" ortak Oghuz Türkçesi kökü). Doğrulama sonrası `payment_links.`
öneki TAMAMEN temiz — `"payment_links."` `fullyTranslatedPrefixes`'e eklendi.

**✅ `users.*` TAMAMLANDI (2026-09-23, Claude).** Doküman 49 anahtar tahmin etmişti; gerçek tarama
87 anahtarda 1668 boşluk buldu (~1.8x tahmin — `context.*` alt-öneği toplu tenant/site erişim
matrisi UI'ı için 34 anahtarlık ayrı bir alt-modül içeriyordu). 83 anahtar 28 dilde tam çevrildi
(4 alt-parti halinde üretildi); kalan 4 anahtar (`form.name`, `form.phone`, `form.surname`,
`stats.languages`) zaten tam doluydu, sadece belirli dillerde meşru ortak-kök gerekçesi gerekiyordu.
Doğrulamada 2 gerçek hata bulundu ve düzeltildi: `context.scope.tenants` ve `form.companyName`
Azerice değerleri Türkçe metnin birebir kopyası kalmıştı — genuine Azerice ifadelerle
("Görünən tenant sayı: {{count}}", "Şirkət adı") değiştirildi. 23 yanlış-pozitif gerekçelendirildi
(Azerice/Boşnakça/Çekçe/Almanca/Hırvatça/Lehçe/Türkmence/Özbekçe ortak Oğuz Türkçesi kökleri ve
uluslararası ödünç kelimeler: "Etiket", "Ad Soyad", "Sil", "Ad", "Soyad", "Diller", "Adres",
"Telefon"). Doğrulama sonrası `users.` öneki TAMAMEN temiz — `"users."` `fullyTranslatedPrefixes`'e
eklendi.

**✅ `permissions.*` TAMAMLANDI (2026-09-23, Claude).** Doküman 48 anahtar tahmin etmişti; gerçek
tarama tam 48 anahtarda 914 boşluk buldu (tahmin tam doğruydu). Bu modülde farklı bir desen vardı:
9 dil (`tr/en/de/fr/it/es/pt/pt-br/az`) zaten doğru çevrilmişti, ama kalan 19 dilin (zh, ja, fa, ar,
kk, uz, ru, uk, pl, hr, cs, bs, sr, sq, hy, az-cyrl, ky, tk, mn) tamamı Türkçe metnin birebir
kopyasıydı — 48 anahtarın tamamında bu 19 dil genuine çevrildi (3 alt-parti halinde). Doğrulamada
3 boşluk çıktı: `notice.update.message` (`{{module}} - {{action}}` — sadece değişkenlerden oluşan,
çevrilecek doğal dil metni içermeyen bir kalıp) 12 Latin-dışı dilde İngilizce ile birebir kalıyordu —
tipografik olarak yerelleştirilmiş em-dash ("–") ile düzeltildi (kısa çizgi yerine), 2 yanlış-pozitif
gerekçelendirildi (`table.delete`/`table.label` Azerice — "Sil"/"Etiket", users.* emsalleriyle aynı
gerekçe). Doğrulama sonrası `permissions.` öneki TAMAMEN temiz — `"permissions."`
`fullyTranslatedPrefixes`'e eklendi.

**✅ `discounts.*` TAMAMLANDI (2026-09-23, Claude).** Doküman 44 anahtar tahmin etmişti; gerçek
tarama tam 44 anahtarda 1144 boşluk buldu (tahmin tam doğruydu). Tüm 44 anahtar sadece `{tr, en}`
ile seed edilmişti — 26 dilin tamamı genuine çevrildi (2 alt-parti halinde). Doğrulamada 15
yanlış-pozitif çıktı, tümü gerekçelendirildi (Azerice ortak Oğuz Türkçesi kökleri: "Ad", "Sabit";
Azerice/Boşnakça/Hırvatça/Lehçe/Türkmence/Özbekçe ortak uluslararası ödünç kelime: "Kod" —
languages.table.code emsaliyle aynı gerekçe). Doğrulama sonrası `discounts.` öneki TAMAMEN temiz —
`"discounts."` `fullyTranslatedPrefixes`'e eklendi.

**✅ `editor.*` TAMAMLANDI (2026-09-23, Claude).** Doküman 42 anahtar tahmin etmişti; gerçek tarama
tam 42 anahtarda 767 boşluk buldu (tahmin tam doğruydu). `permissions.*` ile aynı karma desen:
9 dil (`tr/en/de/fr/it/es/pt/pt-br/az`) çoğunlukla doğruydu, kalan 19 dil çoğu anahtarda Türkçe
metnin birebir kopyasıydı (2 alt-parti halinde genuine çevrildi); ayrıca bare "Embed"/"HTML"
teknik terimleri (`embedModal.badge`, `embedModal.urlLabel`, `toolbar.embed`, `toolbar.source`)
12 Latin-dışı dilde İngilizce ile birebir kalıyordu — bu 4 anahtar için sadece o 12 dile genuine
yerel karşılık/nitelik eklendi (ör. zh "嵌入", kk "HTML коды"). Doğrulamada 7 yanlış-pozitif
gerekçelendirildi (Azerice ortak Oğuz Türkçesi kökleri/ödünç kelimeler: "Alıntı", "İtalik",
"Geri al"; "Kod" — languages.table.code emsaliyle aynı gerekçe, az/pl/tk/uz). **`embedModal.
urlPlaceholder`** (`https://www.youtube.com/watch?v=...`) kasıtlı olarak TÜM dillerde aynı
bırakıldı — bu bir örnek URL'dir, doğal dil içermez (plugins./menus. kod-örneği istisnasıyla
aynı kategori). Bu nedenle `"editor."` `fullyTranslatedPrefixes`'e EKLENMEDİ (tek kalan boşluk bu
kasıtlı istisna, 12 dil × 1 anahtar = 12 boşluk).

Bu listedeki isimlendirilmiş öncelikler TAMAMEN tükendi (users/permissions/discounts/editor/
categories dahil). **Sıradaki öncelik artık Bölüm 9'daki tam kapsam taramasından geliyor** —
bu liste yalnızca tarihsel referans için korunuyor, güncel değil.

En çok etkilenen önekler (≥10 dilde Türkçe yedek): plugins 121, integrations 84, contracts 82, menus 81, posts 72,
marketplaces 68, languages 62 (→ WP-1), orders 60, dashboard 56, payment_links 54, users 49, permissions 48, discounts 44,
editor 42, categories 40,
**currencies 468 boşluk / 39 anahtar (2026-09-22'de WP-2 sırasında keşfedildi — `currencies.fx.*` alt-öneği, farklı bir özellik: kur listesi; `allLangs`/literal-EN ile seed edilmiş)**,
**persons 156 boşluk / 13 anahtar (2026-09-22'de WP-2 sırasında keşfedildi — `persons.affiliation.*`, `persons.role.*`, `persons.status.*`, `persons.field.route` alt-önekleri; `allLangs`/literal-EN ile seed edilmiş, yeni eklenen 45 `persons.*` anahtarından tamamen bağımsız)**,
**institutions 60 boşluk / 5 anahtar (2026-09-22'de WP-2 sırasında keşfedildi — `institutions.status.*`, `institutions.type.*` alt-önekleri; `allLangs`/literal-EN ile seed edilmiş, yeni eklenen 42 `institutions.*` anahtarından tamamen bağımsız)**,
**wallets 24 boşluk / 2 anahtar (2026-09-22'de WP-2 sırasında keşfedildi — `wallets.status.*` alt-öneği; `allLangs`/literal-EN ile seed edilmiş, yeni eklenen 58 `wallets.*` anahtarından tamamen bağımsız)**,
**telematics.alerts.notice.acknowledged 11 boşluk / 1 anahtar (2026-09-22'de WP-2 sırasında keşfedildi; `allLangs` ile seed edilmiş, yeni eklenen 39 `telematics.alerts.*` anahtarından tamamen bağımsız)**,
**telematics.commands.notice.* 77 boşluk / 7 anahtar (2026-09-22'de WP-2 sırasında keşfedildi; `allLangs` ile seed edilmiş, yeni eklenen 47 `telematics.commands.*` anahtarından tamamen bağımsız)**,
**telematics.alert_rules.notice.* 22 boşluk / 2 anahtar (2026-09-22'de WP-2 sırasında keşfedildi; `allLangs` ile seed edilmiş, yeni eklenen 81 `telematics.alert_rules.*` anahtarından tamamen bağımsız)**,
**telematics.dashboard.ready 11 boşluk / 1 anahtar (2026-09-22'de WP-2 sırasında keşfedildi; `allLangs` ile seed edilmiş, yeni eklenen 161 `telematics.dashboard.*` anahtarından tamamen bağımsız)**,
**fleet_tally.connectors.{title,edit,col.display_name,validation.integration_connector_*} 4 anahtar (2026-09-23'te WP-2 sırasında keşfedildi; `allLangs`/literal-EN ile seed edilmiş, yeni eklenen 28 `fleet_tally.connectors.*` anahtarından tamamen bağımsız)**,
**fleet_tally.gps_providers.{add,col.active,col.base_url,col.display_name,col.provider_type,credentials.hint,edit,title} 8 anahtar (2026-09-23'te WP-2 sırasında keşfedildi; `allLangs`/literal-EN ile seed edilmiş, yeni eklenen 28 `fleet_tally.gps_providers.*` anahtarından tamamen bağımsız)**,
**fleet_tally.trip_assignments.{btn.plan,col.planned_date,field.planned_date,field.scheduled_route,modal.plan,notice.planned,source.planned,status.planned} 8 anahtar (2026-09-23'te WP-2 sırasında keşfedildi; `allLangs` ile seed edilmiş, yeni eklenen 33 `fleet_tally.trip_assignments.*` anahtarından tamamen bağımsız)**,
**fleet_tally.scheduled_routes.{btn.clone,col.validity,col.visible,customer_company_hint,description_create,edit_hint,fares_hint,field.customer_company,field.effective_from,field.effective_to,field.visible_to_riders,modal.clone,no_end_date,section.fares,section.validity,validity_hint,visible_no,visible_yes} 18 anahtar (2026-09-23'te WP-2 sırasında keşfedildi; `allLangs` ile seed edilmiş, yeni eklenen 38 `fleet_tally.scheduled_routes.*` anahtarından tamamen bağımsız)**.

- **Yöntem:** önek önek; her önek tamamlandığında `fullyTranslatedPrefixes` listesine eklenir (ratchet — geri gidilemez).
- **Not:** yeni seed'lerde yaygın `allLangs(tr, en)` yardımcısı tüm dilleri İngilizce ile doldurur; Latin dışı 12 dilde bu boşluk sayılır. Kalıcı çözüm için yeni anahtarlar `allLangsMulti` + overlay ya da tam 28 dille yazılmalı.

## 3. Sıra ve bağımlılıklar

1. **WP-1** (bağımsız, hemen başlanabilir; prompt Bölüm 5).
2. **WP-2** öncelik: kullanıcıların en çok gördüğü modüller (`common`, `currencies`, `persons`, `institutions`, `wallets`), sonra telematics / fleet_tally.
3. **WP-3** öncelik: `menus`, `dashboard`, `users`, `permissions` (yönetim kabuğu), sonra iş modülleri.
4. Her paket: teknik referans `done/TASK_06x_*.md` + özet; README "Task Documentation" kuralı.

## 4. Kabul kapıları (her paket için)

```bash
cd /workspace/go/aystek-core
gofmt -l internal/seeds            # çıktı boş olmalı
go vet ./internal/seeds
go test ./internal/seeds           # TestFullyTranslatedPrefixesHaveNoGaps dahil
TRANSLATION_REPORT_PREFIX=<önek> go test ./internal/seeds -run TestTranslationCoverageReport -v   # 0 anahtar / 0 boşluk

cd /workspace/node/aystek-core-next
npm run check-i18n                 # WP-2 için: baseline küçülmeli, yeni eksik olmamalı
```

Seed sonrası doğrulama (isteğe bağlı): `go run ./cmd/seed` + admin panelde dili değiştirip ilgili sayfayı açmak.

## 5. Çeviri kuralları (tüm paketler)

- Yer tutucular (`{{name}}`, `{{code}}`, `{{count}}`) **birebir** korunur; metin içinde konumları dil dilbilgisine göre değişebilir.
- Yazı sistemleri: `sr` **Kiril** (mevcut kayıtlar öyle), `az` Latin, `az-cyrl` = Azerice metnin Kiril transliterasyonu (1958–1991 alfabesi), `uz` Latin (kesme işareti `‘`), `bs`/`hr` Latin, `kk`/`ky`/`mn`/`ru`/`uk` Kiril, `tk` Latin, `zh` Basitleştirilmiş Çince, `fa` (yarım boşluk/ZWNJ doğru kullanılır), `ar` Modern Standart Arapça.
- Terim tutarlılığı: ortak terimler (Varsayılan, Kaydet, Sil, Düzenle, Etkin/Pasif, Tenant, Site) için `internal/seeds/site_language_translations_seed.go` dosyasındaki karşılıklar kullanılır.
- Türkçe ya da İngilizce metin **kopyalanmaz**; yalnızca marka/teknik ad (AYSTEK, API, URL …) aynı kalabilir.
- Arayüz metni kısa ve eylem odaklı olmalı (buton, sütun başlığı, bildirim).
- Bir dilde metin Türkçe ile **dilbilgisel olarak aynıysa** (ör. Azericede "Kod"), test bunu "same as Turkish" sayar; gerekçe yorumuyla `translation_coverage_test.go` içindeki `identicalToTurkishAllowed["anahtar|dil"] = true` listesine eklenir.

## 6. İlerleme

| Paket | Durum | Tarih | Not |
|---|---|---|---|
| WP-1 `languages.*` | Tamamlandı | 2026-09-22 | Claude (bu sesyon) yaptı, planı tersine çevirdi; `internal/seeds/languages_admin_translations_seed.go` overlay + `fullyTranslatedPrefixes`'e "languages." eklendi; 1200/1200 boşluk kapatıldı, 23'ü (Azerice-Türkçe ortak kök + "Kod" ödünç kelimesi) `identicalToTurkishAllowed` ile gerekçelendirildi |
| WP-2 seed'siz anahtarlar | Devam ediyor | 2026-09-22 | `common` alt-paketi tamamlandı (Claude): 57 yeni anahtar 28 dilde eklendi, 3 anahtar mevcut camelCase eşdeğerlerine yönlendirildi, `common.noRecords`'un kendi 19-dil hatası da düzeltildi. Baseline: 2259 → 2199. `currencies` alt-paketi tamamlandı (Claude): 53 yeni anahtar (`app/(dashboard)/admin/currencies/page.tsx`) 28 dilde eklendi — tek dosyada, yinelenen anahtar yoktu. Baseline: 2199 → 2146. `persons` alt-paketi tamamlandı (Claude): 45 yeni anahtar (`app/(dashboard)/admin/persons/*.tsx`) 28 dilde eklendi, yinelenen anahtar yoktu, 17 Azerice/uluslararası-ödünç-kelime yanlış-pozitifi `identicalToTurkishAllowed`'a gerekçelendirildi. Baseline: 2146 → 2101. `institutions` alt-paketi tamamlandı (Claude): 42 yeni anahtar (`app/(dashboard)/admin/institutions/*.tsx`) 28 dilde eklendi, yinelenen anahtar yoktu (8 anahtar mevcut persons.*/common.* çevirisini yeniden kullandı), 8 yanlış-pozitif gerekçelendirildi + `field.slug`'ın kk/az-cyrl tutarsızlığı düzeltildi. Baseline: 2101 → 2059. `wallets` alt-paketi tamamlandı (Claude): 58 yeni anahtar (`app/(dashboard)/admin/wallets/*.tsx`) 28 dilde eklendi, yinelenen anahtar yoktu (4 anahtar mevcut persons.*/kendi btn.new çevirisini yeniden kullandı), 7 yanlış-pozitif gerekçelendirildi + `reconcile.ledger`'ın az-cyrl tutarsızlığı düzeltildi. Baseline: 2059 → 2001. `telematics.settings` alt-paketi tamamlandı (Claude): 29 yeni anahtar (`app/(dashboard)/admin/telematics/settings/page.tsx`) 28 dilde eklendi — bu modülün frontend varsayılanı İngilizce yazılmıştı, Türkçe değerler gerçek çeviri olarak yazıldı. 2 gap (`teltonika_tcp.tenant_id` az-cyrl/hy) düzeltildi. Baseline: 2001 → 1972. `"telematics.settings."` öneki `fullyTranslatedPrefixes`'e eklendi (pre-existing boşluk yoktu). `telematics.*` 495 anahtar tek modül olamayacak kadar büyük olduğundan alt-önek bazlı işleniyor (bkz. Bölüm 2 WP-2 notu). `telematics.can_decoders` alt-paketi tamamlandı (Claude): 35 yeni anahtar 28 dilde eklendi, yinelenen anahtar yoktu, 9 yanlış-pozitif gerekçelendirildi. Baseline: 1972 → 1937. `"telematics.can_decoders."` öneki de `fullyTranslatedPrefixes`'e eklendi. `telematics.alerts` alt-paketi tamamlandı (Claude): 39 yeni anahtar 28 dilde eklendi, yinelenen anahtar yoktu, 3 yanlış-pozitif gerekçelendirildi. Baseline: 1937 → 1898. `"telematics.alerts."` öneki EKLENMEDİ (pre-existing `notice.acknowledged` borcu var, bkz. Bölüm 2 WP-3). `telematics.driver_sessions` alt-paketi tamamlandı (Claude): 42 yeni anahtar 28 dilde eklendi, yinelenen anahtar yoktu, 6 yanlış-pozitif çıktı (4 gerekçelendirildi, 2'si tembel kopyalama olduğundan gerçek çeviriyle düzeltildi). Baseline: 1898 → 1856. `"telematics.driver_sessions."` öneki de `fullyTranslatedPrefixes`'e eklendi. `telematics.commands` alt-paketi tamamlandı (Claude): 47 yeni anahtar 28 dilde eklendi, yinelenen anahtar yoktu (2 anahtar mevcut driver_sessions.* çevirisini yeniden kullandı), 5 yanlış-pozitif gerekçelendirildi + `field.sla`'nın 11 dildeki bare-Latin tutarsızlığı düzeltildi. Baseline: 1856 → 1809. `"telematics.commands."` öneki EKLENMEDİ (pre-existing `notice.*` borcu var, bkz. Bölüm 2 WP-3). `telematics.vehicle360` alt-paketi tamamlandı (Claude): 61 yeni anahtar 28 dilde eklendi, yinelenen anahtar yoktu, 8 yanlış-pozitif çıktı (3 gerekçelendirildi, 3'ü daha net Azerice ifadeyle düzeltildi, `can_metric.rpm`'in az-cyrl/ja bare-Latin tutarsızlığı düzeltildi). Baseline: 1809 → 1748. `"telematics.vehicle360."` öneki de `fullyTranslatedPrefixes`'e eklendi. `telematics.alert_rules` alt-paketi tamamlandı (Claude): 81 yeni anahtar 28 dilde eklendi, yinelenen anahtar yoktu, 7 yanlış-pozitif gerekçelendirildi. Baseline: 1748 → 1667. `"telematics.alert_rules."` öneki EKLENMEDİ (pre-existing `notice.*` borcu var, bkz. Bölüm 2 WP-3). `telematics.dashboard` alt-paketi tamamlandı (Claude): 161 yeni anahtar 28 dilde eklendi (4 alt-parti halinde üretildi), yinelenen anahtar yoktu, 1 yanlış-pozitif düzeltildi. Baseline: 1667 → 1506. `"telematics.dashboard."` öneki EKLENMEDİ (pre-existing `ready` borcu var, bkz. Bölüm 2 WP-3). **`telematics.*` (495 anahtar, 8 alt-paket) TAMAMLANDI.** `fleet_tally.route_fares`+`fleet_tally.runs` alt-paketleri tamamlandı (Claude): 10 yeni anahtar (5+5) 28 dilde eklendi, yinelenen anahtar yoktu, 1 yanlış-pozitif gerekçelendirildi. Baseline: 1506 → 1496. Her iki önek de EKLENMEDİ (pre-existing borç var). `fleet_tally.movement_reports` alt-paketi tamamlandı (Claude): 10 yeni anahtar 28 dilde eklendi, yinelenen anahtar yoktu. Baseline: 1496 → 1486. `"fleet_tally.movement_reports."` öneki `fullyTranslatedPrefixes`'e eklendi. `fleet_tally.verifications` alt-paketi tamamlandı (Claude): 14 yeni anahtar 28 dilde eklendi, yinelenen anahtar yoktu (2 anahtar mevcut çeviri yeniden kullandı), 2 gap (`col.push`, `col.gps_provider` bare-Latin tutarsızlığı) düzeltildi. Baseline: 1486 → 1472. `"fleet_tally.verifications."` öneki EKLENMEDİ (pre-existing borç var, bkz. Bölüm 2 WP-3). `fleet_tally.drivers` alt-paketi tamamlandı (Claude): 20 yeni anahtar 28 dilde eklendi, yinelenen anahtar yoktu (6 anahtar mevcut persons.*/institutions.* çevirisini yeniden kullandı), 1 yanlış-pozitif gerekçelendirildi. Baseline: 1472 → 1452. `"fleet_tally.drivers."` öneki `fullyTranslatedPrefixes`'e eklendi. `fleet_tally.vehicle_mappings` alt-paketi tamamlandı (Claude): 21 yeni anahtar 28 dilde eklendi, yinelenen anahtar yoktu (4 anahtar mevcut çeviri yeniden kullandı), 3 gap düzeltildi/gerekçelendirildi. Baseline: 1452 → 1431. `"fleet_tally.vehicle_mappings."` öneki EKLENMEDİ (pre-existing borç var, bkz. Bölüm 2 WP-3). `fleet_tally.control_types` alt-paketi tamamlandı (Claude): 24 yeni anahtar 28 dilde eklendi, yinelenen anahtar yoktu (6 anahtar mevcut çeviri yeniden kullandı), 2 yanlış-pozitif gerekçelendirildi. Baseline: 1431 → 1407. `"fleet_tally.control_types."` öneki EKLENMEDİ (pre-existing borç var, bkz. Bölüm 2 WP-3). `fleet_tally.schedules` alt-paketi tamamlandı (Claude): 26 yeni anahtar 28 dilde eklendi, yinelenen anahtar yoktu (5 anahtar mevcut çeviri yeniden kullandı), 2 gerçek hata düzeltildi + 4 yanlış-pozitif gerekçelendirildi. Baseline: 1407 → 1381. `"fleet_tally.schedules."` öneki EKLENMEDİ (pre-existing borç var, bkz. Bölüm 2 WP-3). `fleet_tally.vehicle_alarms` alt-paketi tamamlandı (Claude): 26 yeni anahtar 28 dilde eklendi, yinelenen anahtar yoktu (6 anahtar mevcut çeviri yeniden kullandı), 1 yanlış-pozitif gerekçelendirildi. Baseline: 1381 → 1355. `"fleet_tally.vehicle_alarms."` öneki `fullyTranslatedPrefixes`'e eklendi. `fleet_tally.connectors` alt-paketi tamamlandı (Claude): 28 yeni anahtar 28 dilde eklendi, yinelenen anahtar yoktu (6 anahtar mevcut çeviri yeniden kullandı), 3 yanlış-pozitif gerekçelendirildi. Baseline: 1355 → 1327. `"fleet_tally.connectors."` öneki EKLENMEDİ (pre-existing borç var, bkz. Bölüm 2 WP-3). `fleet_tally.gps_providers` alt-paketi tamamlandı (Claude): 28 yeni anahtar 28 dilde eklendi, yinelenen anahtar yoktu (16 anahtar mevcut çeviri yeniden kullandı), 1 yanlış-pozitif gerekçelendirildi. Baseline: 1327 → 1299. `"fleet_tally.gps_providers."` öneki EKLENMEDİ (pre-existing borç var, bkz. Bölüm 2 WP-3). `fleet_tally.supplier_invoices` alt-paketi tamamlandı (Claude): 33 yeni anahtar 28 dilde eklendi, yinelenen anahtar yoktu (4 anahtar dahili yeniden kullanım), 0 yanlış-pozitif/boşluk — tamamen temiz. Baseline: 1299 → 1266. `"fleet_tally.supplier_invoices."` öneki `fullyTranslatedPrefixes`'e eklendi (önek tamamen yeniydi). `fleet_tally.trip_assignments` alt-paketi tamamlandı (Claude): 33 yeni anahtar 28 dilde eklendi, yinelenen anahtar yoktu (1 anahtar dahili yeniden kullanım), 4 yanlış-pozitif gerekçelendirildi. Baseline: 1266 → 1233. `"fleet_tally.trip_assignments."` öneki EKLENMEDİ (pre-existing borç var, bkz. Bölüm 2 WP-3). `fleet_tally.scheduled_routes` alt-paketi tamamlandı (Claude): 38 yeni anahtar 28 dilde eklendi, yinelenen anahtar yoktu (3 anahtar dahili/persons.* yeniden kullanım), 2 gerçek hata düzeltildi + 3 yanlış-pozitif gerekçelendirildi. Baseline: 1233 → 1195. `"fleet_tally.scheduled_routes."` öneki EKLENMEDİ (pre-existing borç var, bkz. Bölüm 2 WP-3). `fleet_tally.route_definitions` alt-paketi tamamlandı (Claude): 39 yeni anahtar 28 dilde eklendi, yinelenen anahtar yoktu (3 anahtar dahili/persons.* yeniden kullanım), 2 sorun düzeltildi/gerekçelendirildi — tamamen yeni önek, temiz. Baseline: 1195 → 1156. `"fleet_tally.route_definitions."` öneki `fullyTranslatedPrefixes`'e eklendi. 🎯 **`fleet_tally.*` (350 anahtar, 15 alt-paket) TAMAMLANDI** — bununla WP-2'nin şu anda listelenen tüm seed'siz anahtar backlog'u (telematics.* 495 + fleet_tally.* 350 = 845 anahtar) tamamlandı |
| WP-3 yedeğe düşmüş anahtarlar | Devam ediyor | 2026-09-23 | Tamamlanan önekler (`fullyTranslatedPrefixes`'e eklendi): `contracts.` (82 anahtar), `posts.` (72), `marketplaces.` (68), `orders.` (63), `dashboard.` (64), `payment_links.` (54), `users.` (87), `permissions.` (48), `discounts.` (44), `categories.` (69). `plugins.`, `menus.` ve `editor.` tamamlandı ama kasıtlı kod/örnek-URL istisnaları nedeniyle listeye EKLENMEDİ. `integrations.*` (1232 anahtar, 14784 boşluk) ayrı büyük görev olarak ertelendi. **Bölüm 2'nin isimlendirilmiş önceliği TÜKENDİ** — tam kapsam taraması 7057 anahtar/83690 boşluk kalan öneki ortaya çıkardı (bkz. Bölüm 9, güncel öncelik kaynağı). `stock.*` (922 anahtar) ve `pms.*` (533 anahtar) alt-önek bazlı TAMAMLANDI, detaylı ilerleme Bölüm 10/11'de. Sırada `forum.` (4164 boşluk). |

## 7. Diğer yapay zeka araçları için prompt (WP-1)

Aşağıdaki bloğu olduğu gibi kopyalayın.

````text
ROL
AYSTEK Core projesinde çalışan kıdemli bir yerelleştirme (i18n) mühendisisin. Görevin WP-1: Dil Yönetimi sayfasının
(`languages.*`) çeviri borcunu kapatmak. Bu görevde kod mantığı yazmıyorsun; çeviri verisi ve bir küçük test listesi
güncelliyorsun. Kullanıcıyla (Hüseyin Koç) sohbet dili TÜRKÇE'dir; kod, yorum ve commit gövdesi mevcut proje
geleneğini izler (kod İngilizce, commit mesajı Türkçe).

ZORUNLU OKUMA (başlamadan, bu sırayla)
1. /workspace/aystek-core-dev-docs/README_FOR_AI_ASSISTANTS.md   (anayasa; özellikle RULE #1 çeviri seed'i ve "Remove obsolete work")
2. /workspace/aystek-core-dev-docs/I18N_TRANSLATION_DEBT_WORKPLAN_TR.md   (WP-1 ve Bölüm 5 çeviri kuralları)
3. /workspace/aystek-core-dev-docs/I18N_WP1_LANGUAGES_SOURCE_TEXTS_TR.md   (62 anahtarın tr/en/de kaynak metinleri)
4. /workspace/go/aystek-core/internal/seeds/site_language_translations_seed.go   (uygulanacak desen + terim tutarlılığı)
5. /workspace/go/aystek-core/internal/seeds/translation_coverage_test.go   (kabul testi)

BAŞLANGIÇ ÖLÇÜMÜ
cd /workspace/go/aystek-core
TRANSLATION_REPORT_PREFIX=languages. go test ./internal/seeds -run TestTranslationCoverageReport -v
Beklenen başlangıç: "62 keys with gaps, 1200 gaps". Bu çıktı, dosyadaki tablodan önceliklidir.

YAPILACAKLAR
1. Yeni dosya: internal/seeds/languages_admin_translations_seed.go
   - `var languagesAdminTranslations = map[string]map[string]string{ "<dil>": { "<anahtar>": "<çeviri>", ... }, ... }`
   - `func applyLanguagesAdminTranslations(entries map[string]map[string]string)` — her dil/anahtar için mevcut
     yardımcı `setEntryLanguage(entries, key, lang, text)` kullanılır (site_language_translations_seed.go içinde tanımlı).
2. `buildTranslationEntries()` içinde `applySiteLanguageTranslations(entries)` satırından hemen sonra
   `applyLanguagesAdminTranslations(entries)` çağrısını ekle (translation_seed.go).
3. Kapsanacak diller: zh, ja, fa, ar, ru (AKTİF diller — önce bunlar), sonra kk, uz, uk, pl, hr, cs, bs, sr, sq, hy,
   az-cyrl, ky, tk, mn; ayrıca `az` (16 anahtar) ve `languages.delete.message` için de, fr, it, es, pt, pt-br.
   Hangi dil/anahtar eksik: rapor çıktısındaki her "[dil: neden]" satırı.
4. Kabul: `translation_coverage_test.go` içindeki `fullyTranslatedPrefixes` listesine "languages." ekle.

KURALLAR
- Mevcut translation_seed.go satırlarını DÜZENLEME (overlay deseni); başka anahtarlara dokunma; hiçbir anahtarı silme.
- Yer tutucular `{{name}}`, `{{code}}`, `{{count}}` birebir korunur.
- Türkçe/İngilizce metni kopyalama; makine çevirisi yer tutucusu ("TODO", "[tr]" vb.) bırakma.
- Yazı sistemleri: sr KİRİL; az Latin; az-cyrl = Azerice (Latin) metnin Kiril transliterasyonu (a→а b→б c→ҹ ç→ч d→д e→е ə→ә
  f→ф g→ҝ ğ→ғ h→һ x→х ı→ы i→и j→ж k→к q→г l→л m→м n→н o→о ö→ө p→п r→р s→с ş→ш t→т u→у ü→ү v→в y→ј z→з);
  uz Latin (kesme işareti ‘); bs/hr/tk/pl/cs/sq Latin; kk/ky/mn/ru/uk Kiril; zh Basitleştirilmiş Çince; fa'da ZWNJ doğru;
  ar Modern Standart Arapça; hy Ermenice.
- Terimler site_language_translations_seed.go ile tutarlı olmalı (Varsayılan, Kaydet, Sil, Düzenle, Tenant, Site).
- Arayüz metni kısa ve eylem odaklı (buton, sütun başlığı, bildirim).
- Dilbilgisel olarak Türkçe ile aynı olan tekil metinler (ör. az "Kod") için testteki `identicalToTurkishAllowed["anahtar|dil"] = true` listesine GEREKÇE YORUMUYLA ekle; bunu çeviri yapmaktan kaçınmak için kullanma.
- Çalışma yöntemi: 3-4 parti (1: zh ja fa ar ru · 2: kk uz uk pl hr · 3: cs bs sr sq hy · 4: az-cyrl ky tk mn + az + tekil boşluklar).
  Her partiden sonra raporu yeniden çalıştır; boşluk sayısı düşmeli.
- Git: yalnızca kendi dosyalarını stage et (`git add` ile tek tek; `.env` ASLA), commit mesajı Türkçe, sonuna
  "Co-Authored-By: <araç adın> <noreply@...>" ekle. PUSH ETME; push'u proje sahibi yapar.

DOĞRULAMA (hepsi geçmeli)
cd /workspace/go/aystek-core
gofmt -l internal/seeds                      # boş olmalı
go vet ./internal/seeds
go test ./internal/seeds                     # TestFullyTranslatedPrefixesHaveNoGaps dahil
TRANSLATION_REPORT_PREFIX=languages. go test ./internal/seeds -run TestTranslationCoverageReport -v   # "0 keys with gaps"

TESLİM RAPORU (Türkçe, kısa)
- Dil başına tamamlanan anahtar sayısı ve toplam çeviri sayısı (hedef 1200)
- Yukarıdaki dört komutun çıktısı (son satırlar)
- Emin olmadığın / ana dil konuşuru gözden geçirmesi gereken çeviriler (anahtar + dil + neden)
- Dokunduğun dosyaların listesi
- Ardından /workspace/aystek-core-dev-docs/I18N_TRANSLATION_DEBT_WORKPLAN_TR.md Bölüm 6 tablosunda WP-1 satırını "Tamamlandı" yap.
````

## 8. Sonraki prompt'lar

WP-2 ve WP-3 için prompt, WP-1 tamamlanıp desen doğrulandıktan sonra bu dosyanın Bölüm 7 yapısıyla (rol → zorunlu okuma →
ölçüm → yapılacaklar → kurallar → doğrulama → rapor) hazırlanır; önek listesi ve kabul komutları Bölüm 2 ve 4'te tanımlıdır.

## 9. Tam Kapsam Taraması ve Güncel İş Planı (2026-09-23)

Bölüm 2'deki "en çok etkilenen önekler" listesi doküman oluşturulurken yapılan kısmi bir taramaya
dayanıyordu ve WP-3'ün isimlendirilmiş önceliklerinin (permissions, discounts, editor, categories)
tamamlanmasının ardından **tükendi**. Bu bölüm, kod tabanının TAMAMINI kapsayan güncel bir tarama
ile üretildi ve artık WP-3'ün (integrations.* hariç) resmi öncelik kaynağıdır — Bölüm 2'nin yerini
alır, onu geçersiz kılar.

**Ölçüm komutu:**
```bash
cd /workspace/go/aystek-core
TRANSLATION_REPORT_PREFIX='*' go test ./internal/seeds -run TestTranslationCoverageReport -v > /tmp/full_scan_report.txt 2>&1
# ardından üst-düzey önek (ilk nokta öncesi) bazında grupla
```

Toplam: **7057 anahtar / 83690 boşluk**, 107 farklı üst-düzey önekte (2026-09-23 tam tarama, `integrations.*` dahil).

**Güncel kalan (2026-09-25): 5254 anahtar / 61998 boşluk** (`integrations.*` dahil).
Bu sayaç `forum.*` ve `system_ops.action.*` tamamlandıktan sonra aynı tam tarama yöntemiyle
hesaplandı; aşağıdaki tablo başlangıç envanterini korur, tamamlanan öneklerin ayrıntısı sonraki
ilerleme bölümlerindedir.

| Önek | Anahtar | Boşluk | Not |
|---|---:|---:|---|
| `integrations` | 1232 | 14784 |  |
| `stock` | 922 | 11120 |  |
| `pms` | 533 | 6396 |  |
| `forum` | 347 | 4164 |  |
| `system_ops` | 296 | 3552 |  |
| `menu` | 290 | 3550 |  |
| `public_site` | 320 | 2816 |  |
| `sms` | 226 | 2712 |  |
| `fleet_tally` | 188 | 2234 | WP-2 enumerated backlog (350 anahtar) TAMAMLANDI; bu sayı WP-2 kapsamı DIŞINDAKİ diğer fleet_tally alt-öneklerindeki pre-existing borç (bkz. Bölüm 6 — route_fares/runs/verifications/vehicle_mappings/control_types/schedules kısmen çevrildi ama ratchet'e eklenemedi) |
| `shipping` | 138 | 1656 |  |
| `personal_assistant` | 127 | 1524 |  |
| `events` | 117 | 1404 |  |
| `settings` | 121 | 1343 |  |
| `classifieds` | 99 | 1188 |  |
| `news_sources` | 90 | 1080 |  |
| `ai` | 88 | 1056 |  |
| `telematics` | 86 | 1032 | WP-2 enumerated backlog (495 anahtar) TAMAMLANDI; bu sayı WP-2 kapsamı DIŞINDAKİ diğer telematics alt-önekleri (ör. telematics.events, telematics.geofences vb. — henüz taranmadı) |
| `digital_editions` | 86 | 1032 |  |
| `content_ops_settings` | 83 | 996 |  |
| `coupons` | 38 | 988 |  |
| `returns` | 37 | 962 |  |
| `invoices` | 79 | 948 |  |
| `payment_intents` | 40 | 928 |  |
| `content_bots` | 77 | 924 |  |
| `news_agencies` | 75 | 900 |  |
| `payments` | 36 | 900 |  |
| `condolences` | 73 | 876 |  |
| `checkout` | 31 | 806 |  |
| `units` | 31 | 775 |  |
| `taxes` | 61 | 732 |  |
| `biographies` | 61 | 732 |  |
| `galleries` | 58 | 696 |  |
| `translations` | 36 | 651 |  |
| `contacts` | 33 | 632 |  |
| `news` | 49 | 588 |  |
| `interviews` | 48 | 576 |  |
| `tenantSelector` | 43 | 516 |  |
| `stores` | 40 | 469 |  |
| `payment_gateways` | 18 | 468 |  |
| `ads` | 36 | 399 |  |
| `sports` | 33 | 396 |  |
| `variants` | 15 | 390 |  |
| `licenses` | 30 | 360 |  |
| `payment_risk` | 28 | 336 |  |
| `message` | 36 | 327 |  |
| `access_points` | 25 | 300 |  |
| `library` | 23 | 276 |  |
| `discount` | 9 | 234 |  |
| `access_blocks` | 19 | 228 |  |
| `cart` | 6 | 156 |  |
| `media` | 46 | 135 |  |
| `page` | 6 | 114 |  |
| `service_reservations` | 8 | 96 |  |
| `module` | 8 | 96 |  |
| `header` | 9 | 84 |  |
| `block` | 14 | 78 |  |
| `profile` | 22 | 65 |  |
| `graphql` | 39 | 63 |  |
| `receivables` | 5 | 60 |  |
| `wallet_instruments` | 4 | 48 |  |
| `erp` | 10 | 47 |  |
| `authors` | 8 | 44 |  |
| `products` | 12 | 42 |  |
| `contact` | 7 | 38 |  |
| `status` | 2 | 38 |  |
| `errors` | 3 | 36 |  |
| `sliders` | 13 | 35 |  |
| `plugins` | 3 | 35 | kasıtlı istisna: kod/JSON örnekleri (WP-3'te işlendi, bkz. Bölüm 6) |
| `localizedContent` | 22 | 30 |  |
| `register` | 21 | 27 |  |
| `view_profiles` | 14 | 27 |  |
| `comments` | 7 | 25 |  |
| `mail_templates` | 5 | 24 |  |
| `actions` | 2 | 24 |  |
| `catalog` | 7 | 20 |  |
| `category` | 1 | 19 |  |
| `survey` | 12 | 18 |  |
| `pages` | 8 | 18 |  |
| `reset` | 16 | 16 |  |
| `order` | 14 | 16 |  |
| `seo` | 4 | 16 |  |
| `references` | 11 | 13 |  |
| `forgot` | 9 | 13 |  |
| `editor` | 1 | 12 | kasıtlı istisna: embedModal.urlPlaceholder (örnek URL, plugins./menus. ile aynı kategori) |
| `subscriptions` | 11 | 11 |  |
| `menus` | 1 | 11 | kasıtlı istisna: kod/dosya-adı örnekleri (WP-3'te işlendi, bkz. Bölüm 6) |
| `verify` | 10 | 10 |  |
| `productMedia` | 3 | 10 |  |
| `login` | 1 | 9 |  |
| `tags` | 7 | 8 |  |
| `companies` | 1 | 8 |  |
| `launchpad` | 3 | 7 |  |
| `product` | 6 | 6 |  |
| `notifications` | 5 | 6 |  |
| `roles` | 3 | 4 |  |
| `error` | 3 | 3 |  |
| `table` | 2 | 2 |  |
| `sidebar` | 2 | 2 |  |
| `noaccess` | 2 | 2 |  |
| `installments` | 2 | 2 |  |
| `deployment` | 2 | 2 |  |
| `activityLogs` | 2 | 2 |  |
| `uat` | 1 | 1 |  |
| `subscription_contracts` | 1 | 1 |  |
| `iconPicker` | 1 | 1 |  |
| `golive` | 1 | 1 |  |
| `activity` | 1 | 1 |  |
**Notlar:**
- Bu tablo `integrations.*` dahil TÜM kalan boşlukları kapsar (integrations.* ayrı büyük görev
  olarak ertelenmeye devam ediyor, kullanıcı onayı gerekiyor — bkz. Bölüm 2).
- `fleet_tally`/`telematics` satırları YANILTICI görünebilir: WP-2'nin enumerated backlog'u
  (845 anahtar, sırasıyla 350/495) TAMAMEN bitti. Kalan sayılar o backlog'un DIŞINDaki,
  daha önce keşfedilmemiş alt-öneklerden geliyor (bkz. Bölüm 6'daki "EKLENMEDİ" notları ve
  henüz hiç taranmamış alt-modüller).
- `editor`/`menus`/`plugins` satırlarındaki küçük kalıntılar kasıtlı istisnalardır (kod/örnek-URL
  içerir, doğal dil değildir) — WP-3'te zaten işlendi, tekrar dokunulmamalı.
- `languages` bu listede YOK (WP-1 ile tamamen temiz, `fullyTranslatedPrefixes`'te).

**Öncelik sırası (kullanıcı kararı, 2026-09-23):** en büyükten başlayarak sırayla devam.
`integrations.*` (14784 boşluk) zaten ayrı ertelenmiş görev olduğundan, bir sonraki hedef
**`stock`** (922 anahtar / 11120 boşluk) — muhtemelen kendi içinde alt-önek bazlı parçalanacak kadar
büyük (fleet_tally/telematics emsaliyle aynı yöntem). Ardından `pms` (6396), `forum` (4164),
`system_ops` (3552), `menu` (3550), `public_site` (2816), `sms` (2712), ... tabloda yukarıdan
aşağıya sırayla.

**Yöntem:** Bölüm 4'teki standart WP-3 pipeline'ı aynen uygulanır (extract → çevir → güvenli
replace → doğrula → gerekçelendir/düzelt → ratchet kararı → tam test → commit+push → devam).
Her üst-düzey önek kendi içinde büyükse (>150 anahtar), fleet_tally/telematics emsaliyle aynı
şekilde alt-önek bazlı (`stock.products.`, `stock.warehouses.` gibi) parçalanır ve her alt-önek
kendi commit'iyle ilerler.


## 10. `stock.*` İlerleme Detayı

`stock.*` toplam 922 anahtar / 11120 boşluk (2026-09-23 tam tarama). Bölüm 9'daki önceliğe göre
en büyükten başlanarak alt-önek alt-önek işleniyor. Tüm alt-önekler aynı desende çıktı: her anahtar
`allLangs(tr, en)` ile seed edilmiş, 16 Latin yazılı dil zaten `en` (İngilizce) ile dolu kalıyor
(WP-2 allLangs davranışı, ratchet kriterine göre engelleyici değil), sadece 12 Latin-dışı dilde
boşluk sayılıyor. Yöntem: benzersiz `tr` metinlerini çıkar (çoğu anahtar aynı metni paylaşıyor) →
12 dile genuine çevir → tüm anahtarlara eşle → güvenli replace → doğrula → ratchet.

| Alt-önek | Anahtar | Boşluk | Durum | Not |
|---|---:|---:|---|---|
| `stock.settings.` | 213 | 2556 | ✅ Tamamlandı (2026-09-23) | 150 benzersiz ifade |
| `stock.dashboard.` | 154 | 1848 | ✅ Tamamlandı (2026-09-23) | 141 benzersiz ifade |
| `stock.sms.` | 151 | 1812 | ✅ Tamamlandı (2026-09-23) | 147 benzersiz ifade; 4 gerçek hata düzeltildi (fields.override/registry.vendorCount bazı dillerde bare "Override"/"Vendor" bırakılmış, vendor.arvento/mobiliz/n2_mobil/trio marka adları zh/ja'da fonetik çevrilmemiş) |
| `stock.reports.` | 74 | 888 | ✅ Tamamlandı (2026-09-23) | 72 benzersiz ifade; "GPS + HAT"(EN "GPS + SIM") ilk taslakta 12 dilde İngilizce ile çakışıyordu, "SIM" yerel/fonetik karşılıklarla düzeltildi |
| `stock.inventory.` | 66 | 792 | ✅ Tamamlandı (2026-09-23) | 66 benzersiz ifade; IMEI/MSISDN evrensel tanımlayıcıları yerel "numara" niteliğiyle düzeltildi |
| `stock.subscriptions.` | 38 | 456 | ✅ Tamamlandı (2026-09-23) | 38 benzersiz ifade, sorunsuz |
| `stock.sales.` | 37 | 444 | ✅ Tamamlandı (2026-09-23) | 37 benzersiz ifade, sorunsuz |
| `stock.vehicle_assets.` | 34 | 408 | ✅ 33/34 tamamlandı (2026-09-23) | `suggest.plate_placeholder` ("34ABC123") kasıtlı istisna — örnek plaka formatı, doğal dil değil (editor.urlPlaceholder ile aynı kategori); bu tek anahtar nedeniyle `fullyTranslatedPrefixes`'e EKLENMEDİ |
| `stock.catalog.` | 33 | 396 | ✅ Tamamlandı (2026-09-23) | 33 benzersiz ifade; GPS bare akronimi 12 Latin-dışı dilde İngilizce ile çakışıyordu, yerel niteliklerle düzeltildi (ör. ru "GPS-система"); HAT anahtarının (EN "SIM") ja değeri de bare "SIM" kalmıştı, "SIMカード" ile düzeltildi |
| `stock.vehicles.` | 19 | 228 | ✅ Tamamlandı (2026-09-23) | 19 benzersiz ifade, sorunsuz |
| `stock.subscription_pricing.` | 19 | 228 | ✅ Tamamlandı (2026-09-23) | 19 benzersiz ifade, sorunsuz |
| `stock.locations.` | 19 | 228 | ✅ Tamamlandı (2026-09-23) | 19 benzersiz ifade; Heading/Latitude/Longitude/Recorded At/Source/Speed/Vehicle gibi bare İngilizce GPS/telemetri alan adları (tr=en) 12 Latin-dışı dilde genuine coğrafya terminolojisiyle çevrildi |
| `stock.transfer.` | 15 | 180 | ✅ Tamamlandı (2026-09-24) | 15 benzersiz ifade, sorunsuz |
| `stock.returns.` | 15 | 180 | ✅ Tamamlandı (2026-09-24) | 15 benzersiz ifade, sorunsuz |
| `stock.assignments.` | 13 | 156 | ✅ Tamamlandı (2026-09-24) | 13 benzersiz ifade, sorunsuz |
| `stock.portal.` | 12 | 144 | ✅ Tamamlandı (2026-09-24) | 12 benzersiz ifade, sorunsuz |
| diğer küçük stock alt-önekleri (outOfStock/list/lastUpdated/activity/stats/reserved/onHand/integrations/confirm/available) | 10 | 118 | ✅ Tamamlandı (2026-09-24) | `stock.stats.total`'da gerçek hata bulundu (Türkçe değeri bile "Stats Total" bırakılmıştı, hiç yerelleştirilmemiş) — "Toplam"/"Total" ile düzeltildi. Diğer 9 anahtar sorunsuz. |

🎯 **`stock.*` (922 anahtar, 16 alt-paket + 10 tekil anahtar) TAMAMLANDI** — tek kalan boşluk `stock.vehicle_assets.suggest.plate_placeholder` (12 boşluk), kasıtlı istisna (örnek plaka formatı). Sırada Bölüm 9'daki bir sonraki öncelik: `pms.` (6396 boşluk).

**Not:** anahtar/boşluk tahminleri Bölüm 9'daki ilk taramadan; her alt-önek başlanırken
`TRANSLATION_REPORT_PREFIX=stock.X.` ile kesin sayı doğrulanmalı (tahminler ~%5 sapabilir,
görüldüğü üzere subscriptions/sales tahminleri birebir tutmuştu). `stock.*` bitince Bölüm 9'daki
sıraya göre `pms.` (6396 boşluk) başlar.

## 11. `pms.*` İlerleme Detayı

`pms.*` toplam 533 anahtar / 6396 boşluk (2026-09-23 tam tarama). Aynı yöntem: `allLangs(tr,en)`
ile seed edilmiş her anahtar, 16 Latin yazılı dil zaten `en` ile dolu (engelleyici değil), 12
Latin-dışı dilde boşluk sayılıyor. Benzersiz `tr` metinlerini çıkar → 12 dile genuine çevir → tüm
anahtarlara eşle → güvenli replace → doğrula → ratchet.

| Alt-önek | Anahtar | Boşluk | Durum | Not |
|---|---:|---:|---|---|
| `pms.work_orders.` | 84 | 1008 | ✅ Tamamlandı (2026-09-24) | 80 benzersiz ifade, sorunsuz |
| `pms.billing.` | 79 | 948 | ✅ Tamamlandı (2026-09-24) | 72 benzersiz ifade, sorunsuz |
| `pms.reconciliation.` | 66 | 792 | ✅ Tamamlandı (2026-09-24) | 64 benzersiz ifade, sorunsuz |
| `pms.webhooks.` | 43 | 516 | ✅ Tamamlandı (2026-09-24) | 38 benzersiz ifade; 4 anahtarda Türkçe kaynak metninde eksik aksan (ı/ş/ü/ğ) bug'ı da düzeltildi |
| `pms.exports.` | 43 | 516 | ✅ Tamamlandı (2026-09-24) | 43 benzersiz ifade; 1 anahtarda Türkçe kaynak metninde eksik aksan bug'ı da düzeltildi |
| `pms.metering.` | 34 | 408 | ✅ Tamamlandı (2026-09-24) | 34 benzersiz ifade, sorunsuz |
| `pms.access.` | 31 | 372 | ✅ Tamamlandı (2026-09-24) | 31 benzersiz ifade; "PTS" akronimi için Latin-dışı dillerde yerel niteleyici eklendi |
| `pms.occupancies.` | 30 | 360 | ✅ Tamamlandı (2026-09-24) | 30 benzersiz ifade, sorunsuz |
| `pms.properties.` | 29 | 348 | ✅ Tamamlandı (2026-09-24) | 29 benzersiz ifade, sorunsuz |
| `pms.vendors.` | 20 | 240 | ✅ Tamamlandı (2026-09-24) | 20 benzersiz ifade, sorunsuz |
| `pms.units.` | 19 | 228 | ✅ Tamamlandı (2026-09-24) | 18 benzersiz ifade; "PMS" akronimi için Latin-dışı dillerde yerel niteleyici eklendi |
| `pms.portal.` | 17 | 204 | ✅ Tamamlandı (2026-09-24) | 17 benzersiz ifade, sorunsuz |
| `pms.dashboard.` | 12 | 144 | ✅ Tamamlandı (2026-09-24) | 12 benzersiz ifade; `pms.dashboard.title` anahtarında Türkçe kaynak hiç yerelleştirilmemiş bug'ı bulundu ve düzeltildi (tr="PMS Dashboard"→"PMS Panosu", tüm 28 dil için genuine çeviri) |
| `pms.portfolio.` | 8 | 96 | ✅ Tamamlandı (2026-09-24) | 8 benzersiz ifade, sorunsuz |
| `pms.reports.` | 5 | 60 | ✅ Tamamlandı (2026-09-24) | 5 benzersiz ifade, sorunsuz |
| `pms.legal.` | 5 | 60 | ✅ Tamamlandı (2026-09-24) | 5 benzersiz ifade, sorunsuz |
| `pms.notice.` | 3 | 36 | ✅ Tamamlandı (2026-09-24) | 3 benzersiz ifade, sorunsuz |
| diğer tekil anahtarlar (title/subtitle/responsibility/confirm/coming_soon) | 5 | 60 | ✅ Tamamlandı (2026-09-24) | 5 benzersiz ifade, sorunsuz |

🎯 **`pms.*` (533 anahtar, 16 alt-paket + 5 tekil anahtar) TAMAMLANDI** — tam `pms.` prefix taraması 0 boşluk doğrulandı. Sırada Bölüm 9'daki bir sonraki öncelik: `forum.` (4164 boşluk).

## 12. `forum.*` İlerleme Detayı

`forum.*` toplam 347 anahtar / 4164 boşluk (2026-09-24 tam tarama). Aynı yöntem: `allLangs(tr,en)`
ile seed edilmiş her anahtar, 16 Latin yazılı dil zaten `en` ile dolu (engelleyici değil), 12
Latin-dışı dilde boşluk sayılıyor. Benzersiz `tr` metinlerini çıkar → 12 dile genuine çevir → tüm
anahtarlara eşle → güvenli replace → doğrula → ratchet.

| Alt-önek | Anahtar | Boşluk | Durum | Not |
|---|---:|---:|---|---|
| `forum.field.` | 47 | 564 | ✅ Tamamlandı (2026-09-24) | allLangs değil, açık map literal formatı; 40 benzersiz ifade + 7 özel anahtar (Forum/ID/IP/Plan/Slug kognatları + 2 format-örneği placeholder), her biri için Latin-dışı dillere genuine çeviri/niteleyici eklendi |
| `forum.ai.` | 35 | 420 | ✅ Tamamlandı (2026-09-24) | Açık map literal formatı; 32 benzersiz ifade + 3 özel anahtar (Model/Global/Tenant kognatları) |
| `forum.badges.` | 27 | 324 | ✅ Tamamlandı (2026-09-24) | Açık map literal; 27 benzersiz ifade; **12/27 anahtarda Türkçe kaynak metni bozuk bulundu ve düzeltildi** (İngilizce kelimeler doğrudan karışmıştı: "Oluştur Badge"→"Rozet Oluştur", "Kural deleted"→"Kural silindi.", "user badges yok"→"Kullanıcı rozeti yok" vb.) |
| `forum.sponsored.` | 22 | 264 | ✅ Tamamlandı (2026-09-24) | Açık map literal; 21 benzersiz ifade + 1 özel anahtar (Sponsor URL); 9/22 anahtarda Türkçe kaynak metni bozuk bulundu ve düzeltildi (aynı desen: "Sil Ad Placement"→"Reklam Alanını Sil" vb.) |
| `forum.moderation.` | 22 | 264 | ✅ Tamamlandı (2026-09-24) | Açık map literal; 22 benzersiz ifade; **10/22 anahtarda Türkçe kaynak metni bozuk bulundu ve düzeltildi** ("Ata Rapor"→"Raporu Ata", "Başarısız: create ban"→"Yasak oluşturulamadı." vb.); Azerbaycanca Kiril metinleri 1958–1991 alfabe kuralına göre doğrulandı |
| `forum.threads.` | 21 | 252 | ✅ Tamamlandı (2026-09-24) | Açık map literal; 21 benzersiz ifade; **16/21 anahtarın Türkçe kaynak metni düzeltildi veya netleştirildi** ("İşlem failed"→"İşlem başarısız oldu.", "Oluştur Konu"→"Konu Oluştur", teknik ID örnekleri `konu_id_*` / `gönderi_id_*` biçimine getirildi); Latin-dışı dillerde örnek alanların formatı korunarak yerelleştirildi |
| `forum.public.` | 21 | 252 | ✅ Tamamlandı (2026-09-24) | Açık map literal; 21 benzersiz ifade; tümü 12 Latin-dışı dilde tamamlandı; `empty` Türkçe metni çoğul kaynak anlamıyla uyumlu olacak şekilde "Kullanılabilir forum bulunmuyor." olarak düzeltildi |
| `forum.access.` | 19 | 228 | ✅ Tamamlandı (2026-09-24) | Açık map literal; 19 benzersiz ifade; tümü 12 Latin-dışı dilde tamamlandı; 4 Türkçe kaynak/örnek alanı netleştirildi (`Access Kurallar`→`Erişim Kuralları`, `Doğrulama gerektir`→`Doğrulanmış hesap gerektir`, erişim yükleme metni ve rol örneği) |
| `forum.subscriptions.` | 18 | 216 | ✅ Tamamlandı (2026-09-24) | Açık map literal; 18 benzersiz ifade; tümü 12 Latin-dışı dilde tamamlandı; **11/18 Türkçe kaynak metni düzeltildi veya netleştirildi** (`Oluştur Subscription`→`Abonelik Oluştur`, `access plans yok`→`Erişim planı bulunmuyor.`, `Deneme days`→`Deneme süresi (gün)`, `Ara by user ID`→`Kullanıcı ID'sine göre ara` vb.) |
| `forum.reputation.` | 17 | 204 | ✅ Tamamlandı (2026-09-24) | Açık map literal; 17 benzersiz ifade; tümü 12 Latin-dışı dilde tamamlandı; **10/17 Türkçe kaynak metni düzeltildi** (`ledger entries yok`→`Defter kaydı bulunmuyor.`, `Oluştur Kural`→`Kural Oluştur`, `Kural saved`→`Kural kaydedildi.`, kullanıcı arama metni vb.) |
| `forum.posts.` | 15 | 180 | ✅ Tamamlandı (2026-09-24) | Açık map literal; 15 benzersiz ifade; tümü 12 Latin-dışı dilde tamamlandı; **10/15 Türkçe kaynak metni düzeltildi** (`Oluştur Gönderi`→`Gönderi Oluştur`, `Sil failed`→`Gönderi silinemedi.`, `reactions yüklenemedi`→`Tepkiler yüklenemedi.`, kayıt bildirimleri vb.) |
| `forum.status.` | 12 | 144 | ✅ Tamamlandı (2026-09-24) | Açık map literal; 12 benzersiz durum etiketi; Türkçe kaynaklar temizdi, tümü 12 Latin-dışı dilde tamamlandı |
| `forum.forums.` | 10 | 120 | ✅ Tamamlandı (2026-09-24) | Açık map literal; 10 forum yönetimi metni 12 Latin-dışı dilde tamamlandı; 4 bozuk Türkçe kaynak metin düzeltildi |
| `forum.categories.` | 10 | 120 | ✅ Tamamlandı (2026-09-24) | Açık map literal; 10 kategori yönetimi metni 12 Latin-dışı dilde tamamlandı; 7 bozuk Türkçe kaynak metin düzeltildi |
| `forum.stats.` | 7 | 84 | ✅ Tamamlandı (2026-09-24) | Açık map literal; 7 forum istatistik etiketi 12 Latin-dışı dilde tamamlandı; Türkçe kaynaklar temizdi |
| `forum.pm.` | 7 | 84 | ✅ Tamamlandı (2026-09-25) | Açık map literal; 7 özel mesajlaşma metni 12 Latin-dışı dilde tamamlandı; 2 eksik anlamlı Türkçe kaynak metin düzeltildi |
| `forum.indexes.` | 7 | 84 | ✅ Tamamlandı (2026-09-25) | Açık map literal; 7 MongoDB indeks doğrulama metni 12 Latin-dışı dilde tamamlandı; Türkçe kaynaklar temizdi |
| `forum.settings.` | 5 | 60 | ✅ Tamamlandı (2026-09-25) | Açık map literal; 5 forum içerik ayarı 12 Latin-dışı dilde tamamlandı; 1 bozuk Türkçe kaynak metin düzeltildi |
| `forum.report.` | 5 | 60 | ✅ Tamamlandı (2026-09-25) | Açık map literal; 5 rapor durum etiketi 12 Latin-dışı dilde tamamlandı; 2 Türkçe durum etiketi anlam açısından düzeltildi |
| `forum.overview.` | 4 | 48 | ✅ Tamamlandı (2026-09-25) | Açık map literal; 4 forum genel bakış metni 12 Latin-dışı dilde tamamlandı; Türkçe kaynaklar temizdi |
| `forum.analytics.` | 4 | 48 | ✅ Tamamlandı (2026-09-25) | Açık map literal; 4 forum analitik metni 12 Latin-dışı dilde tamamlandı; Türkçe kaynaklar temizdi |
| `forum.type.` | 3 | 36 | ✅ Tamamlandı (2026-09-25) | Açık map literal; 3 forum konu türü etiketi 12 Latin-dışı dilde tamamlandı; Türkçe kaynaklar temizdi |
| `forum.config.` | 3 | 36 | ✅ Tamamlandı (2026-09-25) | Açık map literal; 3 forum yapılandırma metni 12 Latin-dışı dilde tamamlandı; Türkçe kaynaklar temizdi |
| `forum.ban.` | 3 | 36 | ✅ Tamamlandı (2026-09-25) | Açık map literal; 3 yasak hedef türü 12 Latin-dışı dilde tamamlandı; yalın IP etiketi doğal IP adresi karşılıklarıyla yerelleştirildi |
| `forum.reasonModal.` | 3 | 36 | ✅ Tamamlandı (2026-09-25) | Açık map literal; 3 işlem gerekçesi metni 12 Latin-dışı dilde tamamlandı; 2 Türkçe metin kaynak anlamıyla eşlendi |

`forum.*` bitince Bölüm 9'daki sıraya göre `system_ops.` (3552 boşluk) başlar.

**✅ `forum.*` TAMAMLANDI (2026-09-25).** 347 anahtar / 4164 boşluğun tamamı kapatıldı;
`TRANSLATION_REPORT_PREFIX=forum.` doğrulaması `0 anahtar / 0 boşluk` verdi. Alt-önek
mandalları mükerrer taramayı önlemek için tek `"forum."` regresyon mandalında birleştirildi.
Sıradaki paket Bölüm 9'a göre `system_ops.` (296 anahtar / 3552 boşluk).

## 13. `system_ops.*` İlerleme Detayı

`system_ops.*` başlangıç ölçümü: **296 anahtar / 3552 boşluk** (2026-09-25). Alt-gruplar,
tam rapordaki benzersiz anahtarların ikinci segmentine göre ayrıldı; 17 adet doğrudan
`system_ops.<anahtar>` ortak metni ayrıca `shared` grubunda tutuldu. Küçükten büyüğe ilerlenir;
her satır kapsam raporu, regresyon mandalı, `gofmt`, `go vet`, tam seed testi ve izole commit ile kapanır.

| Alt-önek / grup | Anahtar | Boşluk | Durum |
|---|---:|---:|---|
| `system_ops.action.` | 2 | 24 | ✅ Tamamlandı (2026-09-25) |
| `system_ops.form.` | 1 | 12 | Bekliyor |
| `system_ops.ops_rail.` | 2 | 24 | Bekliyor |
| `system_ops.channel.` | 3 | 36 | Bekliyor |
| `system_ops.dashboard.` | 3 | 36 | Bekliyor |
| `system_ops.log_stream.` | 3 | 36 | Bekliyor |
| `system_ops.logs.` | 4 | 48 | Bekliyor |
| `system_ops.status.` | 4 | 48 | Bekliyor |
| `system_ops.trigger_source.` | 7 | 84 | Bekliyor |
| `system_ops.backups.` | 12 | 144 | Bekliyor |
| `system_ops.migrations.` | 12 | 144 | Bekliyor |
| `system_ops.version.` | 12 | 144 | Bekliyor |
| `system_ops.policies.` | 13 | 156 | Bekliyor |
| `system_ops.health.` | 14 | 168 | Bekliyor |
| `system_ops.audit.` | 17 | 204 | Bekliyor |
| `system_ops.rollbacks.` | 17 | 204 | Bekliyor |
| `system_ops.maintenance.` | 18 | 216 | Bekliyor |
| `system_ops.overview.` | 29 | 348 | Bekliyor |
| `system_ops.node.` | 31 | 372 | Bekliyor |
| `system_ops.license.` | 32 | 384 | Bekliyor |
| `system_ops.updates.` | 43 | 516 | Bekliyor |
| `system_ops` ortak (`system_ops.<anahtar>`) | 17 | 204 | Bekliyor |

Güncel System Ops kalanı: **294 anahtar / 3528 boşluk**.
