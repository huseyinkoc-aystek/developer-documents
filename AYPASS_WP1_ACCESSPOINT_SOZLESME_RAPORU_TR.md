# AYPASS WP-1 — AccessPoint Düzenlemesinde Plaka ve POS Kaybının Önlenmesi

**Tarih:** 2026-10-29 · **Hazırlayan:** Claude (Sonnet 5.5) · **Durum:** uygulandı, yerel dalda, **birleştirilmedi / push edilmedi**; Hüseyin Koç incelemesi bekliyor.

> **Bu teslim şunu göstermez:** keşfin tamamlandığını, mali altyapının doğrulandığını ya da uygulamanın yayına hazır olduğunu. P00 hâlâ "kapsam eksikleri kayıtlı keşif raporu"dur (`AYPASS_P00_KESIF_RAPORU_TR.md` §8.1). Cihaz–AccessPoint bağ denetimi (WP-2) açık güvenlik bulgusudur; ele alınmadan **üretim yayınına hazır sonucu çıkarılmaz**.

---

## 1. Referanslar

| Öğe | Değer |
|---|---|
| Çalışılan depo | `/workspace/node/aystek-core-next` (npm adı `aystek-core-next`, remote `admin-panel-next`) |
| Worktree | `/workspace/node/aystek-core-next-wp1` (ayrı; `main` ve `approval-policy-ui` çalışma ağaçlarına müdahale yok) |
| Dal | `fix/aypass-accesspoint-contract` (yeni; adı çakışmadı) |
| Başlangıç SHA'sı | `c5e0994` — `c5e09943c85d989cd72b57954278f54c7217e2f4`, depoda doğrulandı, `origin/main` ile aynı |
| WP-1 commit'leri | `f477dc8` (davranışı değiştirmeyen çıkarma) → `46b3a29` (düzeltme + testler) |
| Backend sözleşme referansı | `aystek-core` `3a3bb86` (dal `aypass-approval-ui`); DTO ve `Service.Update` `git show 3a3bb86:…` ile okundu. Backend dalı değiştirilmedi. |
| Not | `c5e0994` ≠ "canlıdaki sürüm" ve ≠ kesin birleştirme hedefi. Birleştirme hedefi paket sonunda ayrıca değerlendirilecek. |

---

## 2. Kök neden

Backend `Service.Update` (`3a3bb86`) bu alanları **koşulsuz** atar:

```go
current.VehiclePlate = strings.ToUpper(strings.TrimSpace(in.VehiclePlate))   // service.go:688
current.PosEnabled   = in.PosEnabled                                            // service.go:718
```

İstekte alan yoksa Go bunu `""` / `false` olarak okur ve mevcut değerin üzerine yazar. Panel zincirinde iki kopukluk vardı:

1. **`parseAccessPoint`** `vehicle_plate` alanını okumuyordu → form `res.vehicle_plate ?? ""` ile boş yükleniyordu → kayıtta `|| undefined` ile alan istekten düşüyordu.
2. **`pos_enabled`** panelde hiçbir yerde (tip, parser, form, payload) yoktu → istekte hiç bulunmuyordu → her kayıt POS'u `false` yapıyordu.

Yalnızca ad düzenlense bile plaka silinir ve POS kapanırdı. Plakasız `fare_transport` erişim noktasında olay `ErrTransportVehicleRequired` ile reddedilir.

### 2.1 "Alanı göndermemenin etkisi" — çalıştırılarak doğrulandı

Backend dalına dokunmadan, `3a3bb86`'nın **ayrık bir scratch worktree**'sinde atılabilir bir Go testi çalıştırıldı (commit'lenmedi, worktree sonra silindi; backend dalı/ağacı değişmedi):

```
omitted    -> plate="" pos=false name="Servis 2"          # alan istekte yok → SİLİNİR
sent back  -> plate="34 ABC 123" pos=true name="Servis 2" # alan geri gönderilir → KORUNUR
```

Sonuç: sadece alanları payload'dan çıkarmak çözüm değildir, **tam tersine hatayı üretir**. Doğru davranış: mevcut değeri **her güncellemede geri göndermek**.

### 2.2 Sözleşme gerçeği — talimatla ayrışan nokta (önemli)

Talimat `pos_enabled=false` ile alanın hiç gelmemesini ayırmamı istedi. Backend yanıt DTO'su bunu **telde ayırt edilemez** kılıyor:

```go
VehiclePlate string `json:"vehicle_plate,omitempty"`   // dto.go (AccessPointResponse)
PosEnabled   bool   `json:"pos_enabled,omitempty"`     // false → anahtar hiç yazılmaz
RequiresApproval bool `json:"requires_approval"`        // her zaman yazılır
```

`false` ve `""` Go tarafından **anahtarsız** üretilir; aynı tanım `main` `f48748f`'de de vardır. Bu yüzden "alan gelmedi" ile "false" panelde ayırt edilemez ve her POS'u kapalı kayıt "eksik" sayılırsa **hiçbir erişim noktası kaydedilemez**.

**Uygulanan çözüm** (sözleşmeye sadık, yıkıcı varsayılan uydurmaz):

| Yanıt durumu | Sonuç |
|---|---|
| Alan var ve doğru tipte | değer aynen kullanılır |
| Alan yok **ve** yanıt tam (`requires_approval` boolean olarak var) | `""` / `false` = backend'in "sıfır değer" kodlaması; bilinen değer |
| Alan `null` / yanlış tip | **okunamaz** (`undefined`) |
| Yanıt tam değil (`requires_approval` yok), nesne değil, `null`, `[]`, `""` | alan yoksa **okunamaz** |
| Kayıt henüz yüklenmedi | **okunamaz** |

Okunamayan değer varsa **güncelleme gönderilmez**; kullanıcıya `common.saveFailed` ("Kaydedilemedi") + alan adları gösterilir. Böylece boş plaka/`false` uydurulup yıkıcı PUT gitmez.

Kalan sınır: yanıt **tam görünüp** sunucu yine de plakayı eksik yazarsa (ör. omitempty olmayan farklı bir backend sürümü ya da araya giren bir ara katman) panel bunu bilemez. Canlı sürüm doğrulanmadığı için (P00 R-07) bu risk açık.

---

## 3. Yapılan sınırlı düzeltme

| Dosya | Değişiklik |
|---|---|
| `lib/types/access-points.ts` | `pos_enabled?: boolean` eklendi; `vehicle_plate`/`pos_enabled` için "`undefined` = okunamadı" notu. |
| `lib/api/access-points.ts` | `parseAccessPoint` `vehicle_plate` ve `pos_enabled` okur (yukarıdaki tablo); test için `__accessPointsPrivate__` (projedeki `__smsPrivate__` kalıbı). |
| `lib/forms/accessPointPayload.ts` (yeni) | `accessPointToEditForm`, `buildAccessPointUpdatePayload`, `prepareAccessPointUpdate`. Mevcut `lib/forms/authorPayload.ts` kalıbı; yeni bağımlılık/katman yok. `vehicle_plate` ve `pos_enabled` her zaman gönderilir (`""` ve `false` dahil). |
| `app/(dashboard)/admin/access-points/[id]/page.tsx` | Form `pos_enabled` taşır; yükleme/kayıt yardımcılara bağlandı; okunamayan değerde kayıt engellenir. |
| `tests/access-point-contract.test.ts` (yeni) | 13 regresyon/sözleşme testi. |

Diff: `git -C /workspace/node/aystek-core-next-wp1 diff c5e0994..fix/aypass-accesspoint-contract` — ayrıca `AYPASS_WP1_accesspoint-contract.diff` (526 satır; 5 dosya, +401/−22).

**Bilinçli olarak yapılmayanlar:** sayfalama (wallets/persons), `birth_date` biçimi, `fee=0`, create formuna plaka, olay parser'ı, cihaz yetkilendirmesi, mali çekirdek, tarife/offline/devir/Global Market. Create sayfası **hiç değişmedi** (diff'te yok).

---

## 4. Testler

### 4.1 Önce hata (regresyon testi, düzeltmeden önce)

Davranışı değiştirmeyen çıkarma (`f477dc8`) üzerinde, ilk koşu (12 test): **10 başarısız, 2 geçti**.
- Parser testleri gerçek hatayı gösterdi (`expected: '34 ABC 123'`, plaka `undefined`).
- Diğer 8 test, henüz var olmayan `prepareAccessPointUpdate` nedeniyle `TypeError` verdi (yeni API).
- Geçen 2 test, backend sözleşme testleriydi (aşağıda).

Gerçek kod üzerinde GET→PUT zinciri (plaka+POS dolu yanıt, yalnız ad düzenlendi):

```
ÖNCE : parsed.vehicle_plate = undefined | parsed.pos_enabled = undefined
       PUT gövdesinde vehicle_plate: false | pos_enabled: false   (alanlar yok)
SONRA: parsed.vehicle_plate = 34 ABC 123 | parsed.pos_enabled = true
       wire: {"vehicle_plate":"34 ABC 123","pos_enabled":true,"name":"Servis 34 (yeni ad)"}
```

### 4.2 Sonra

`tests/access-point-contract.test.ts`: **13/13 geçti**.

| Zorunlu durum | Test |
|---|---|
| Dolu plaka + `pos_enabled=true`, yalnız ad düzenleme → ikisi korunur (tüm diğer alanlar `deepEqual`) | 3 |
| `pos_enabled=false` açık `false` olarak gider (anahtar var); `true→false` ve `false→true` doğru | 4, 5 |
| Bilinçli plaka değişikliği gider; bilinçli silme `""` olarak gider (anahtar var) | 6 |
| Eksik/`null`/yanlış tipli yanıtta yıkıcı varsayılanla güncelleme yapılmaz | 8, 9, 11 |
| Ulaşım dışı (yemekhane) AccessPoint bozulmaz (satış ürünleri, cihaz bağları, plaka `""`) | 7 |
| Omitempty sıfır değerler "bilinen boş" okunur; mevcut değer tam olmayan gövdede de güvenilir | 2, 10 |
| Parser plaka/POS'u korur | 1 |
| Backend DTO etiketleri (`omitempty`) ve `Update`'in koşulsuz atamaları — `aystek-core` ağacından okunur (`AYSTEK_CORE_DIR`, varsayılan `../../go/aystek-core`; yoksa atlanır) | 12, 13 |

Fixture'lar elle seçilmedi: Go'nun `omitempty` kodlamasını taklit eden `goMarshal` fonksiyonundan geçer. Yine de bunlar **panel tarafı testleridir**; backend ile gerçek uçtan uca entegrasyon (HTTP) çalıştırılmadı. Backend davranışını doğrulayan tek çalışan kanıt §2.1'deki atılabilir Go testidir.

Oluşturma sayfası ve ulaşım dışı akışlar için ayrı bir bileşen testi **yazılmadı**: create sayfası değişmedi (diff'te yok); ulaşım dışı edit akışı test 7 ile yardımcı düzeyinde kapsanır.

### 4.3 Kontrol özeti

| Kontrol | Baseline (`c5e0994`) | WP-1 sonrası | Not |
|---|---|---|---|
| `npm test` (`tsx --test`) | 13 test: 11 geçti, **2 başarısız** | 26 test: 24 geçti, **2 başarısız** | Aynı 2 başarısız test: `getTenantContext rejects invalid tenant id`, `buildTenantHeaders skips invalid tenant id` — **önceden var, bu değişiklikle ilgisiz.** |
| `tsc --noEmit` | 328 hata | 328 hata | Hata kümesi **birebir aynı**; AccessPoint dosyalarında hata yok. (Bir ara `!update.ok` daraltması `strict: false` nedeniyle 1 yeni hata üretti; `=== false` ile giderildi.) |
| `next lint` (değişen 5 dosya) | `lib/api/access-points.ts`: 2 `no-explicit-any` | aynı 2 hata | Önceden var (`(res as any).pagination`, satırlar +21 kaydı). Bu değişiklik 1 `no-unused-vars` getirmişti; giderildi. |
| `npm run check-i18n` | 1205 eksik (1155 baseline), 59 yeni eksik, exit 1 | **çıktı birebir aynı** | Yeni çeviri anahtarı eklenmedi: `common.saveFailed` (seed'lerde var) yeniden kullanıldı. |

**Çalıştırılmayanlar:** tarayıcıda/Next dev sunucusunda elle doğrulama (sayfa `handleSubmit` koruması yalnızca tsc ile doğrulandı, bileşen/e2e testi yok); `next build`; backend ile HTTP entegrasyonu; canlıya karşı herhangi bir şey.

---

## 5. Korunan kullanıcı değişiklikleri

| Öğe | Durum |
|---|---|
| Panel `main` ağacı: `tsconfig.tsbuildinfo` (M), `incognito_check.mjs` (??) | Dokunulmadı; aynı durumda. |
| Panel `main` `c5e0994`, yerel `approval-policy-ui` `434617d` (origin'in 1 commit önünde), `origin/approval-policy-ui` `7b03fe7` | Dokunulmadı. |
| Backend `3a3bb86` ve `.env` (M, izlenen dosya) | Dokunulmadı; içerik okunmadı/yazdırılmadı. Backend'e ayrı scratch worktree ile bakıldı, **silindi**; `git worktree list` yalnız ana ağacı gösteriyor. |
| `aypass_erisim_app` `android/gradle.properties` (M) | Dokunulmadı. |
| Scratchpad'deki eski `wt-panel` (prunable) | Dokunulmadı. |
| WP-1 worktree'sinde `node_modules` | Ana ağacın `node_modules`'üne **symlink**; `git`te izlenmez (untracked görünür), commit'e girmedi. Gerekirse silinebilir. |

Secret, gerçek kişi/kart bilgisi veya ham QR token yazdırılmadı. Push, PR, merge, deploy, migration, canlı erişim yapılmadı.

---

## 6. `approval-policy-ui` uyumluluk notu

Salt okunur karşılaştırma (`c5e0994` ↔ `origin/approval-policy-ui` `7b03fe7` ve yerel `434617d`):

- AccessPoint dosyalarında (`lib/api/access-points.ts`, `lib/types/access-points.ts`, `app/(dashboard)/admin/access-points/**`) **hiçbir fark yok**; o dalda bu yollara dokunan commit yok.
- Daldaki 5 commit: onay politikası/talep ekranları, rol ekranları, `usePaymentIntentRefund`, `lib/api/approvals.ts`, ve **stok `VehicleForm.tsx`** (stok *araç* plakası için neden sorusu). Bu stok aracıdır, AccessPoint'in `vehicle_plate`'i değildir; ayrı varlık.
- Dalın `tests/` altında yalnız onay testleri var (`approvals-*.test.ts`); WP-1'in test dosyasıyla ad çakışması yok.
- Otomatik merge/rebase/cherry-pick yapılmadı.

Sonuç: WP-1 değişikliği ile `approval-policy-ui` arasında dosya çakışması beklenmiyor. Bu bir birleştirme önerisi değildir (KR-8).

---

## 7. Kalan riskler ve takip kayıtları

1. **Backend güncelleme sözleşmesi (WP-15, yeni).** Düzeltme yalnız bu panel istemcisini korur. Alanı göndermeyen başka her istemci (başka panel sürümü, betik, mobil) plakayı/POS'u/`device_bindings`'i/`requires_approval`'ı/`cash_fallback_enabled`'ı hâlâ sıfırlar. Karar: işaretçi/`PATCH` semantiği mi, belgelenmiş tam-PUT + eksik alana 400 mü (KR-13).
2. **POS için arayüz kontrolü yok.** Kullanıcı POS'u panelden değiştiremez; alan yalnızca korunur ve yardımcı düzeyinde değiştirilebilir. Kontrol eklemek 28 dilli çeviri gerektirir; ayrı karar.
3. **Create formunda plaka alanı yok** (kapsam dışı): yeni `fare_transport` erişim noktası önce plakasız oluşur, plaka düzenleme sayfasından girilir (artık kaybolmaz).
4. **Daha önce silinmiş plakalar** geri yazılmadı, tahmin edilmedi; geçmiş veri onarımı yok. Canlıda zaten silinmiş plaka olup olmadığı **doğrulanmadı** (R-06).
5. **Canlı backend sürümü** doğrulanmadı (R-07); `omitempty`/`requires_approval` varsayımı `3a3bb86` ve `f48748f` için doğrulandı, canlı için değil.
6. **Kapsam dışı panel bulguları açık:** D-03 (`fee=0`), D-04 (wallets/persons sayfalama), D-05 (`birth_date`), D-06 (olay parser alanları).
7. **WP-2 (cihaz–erişim noktası bağ denetimi) sıradaki öncelik.** `/device/access-points/:id/events` ve ilişkili uçlar çağıran cihazın o erişim noktasına bağlı olup olmadığına bakmıyor (D-15). Ayrı paket olarak ele alınacak.
8. **Mali doğruluk bu pakette doğrulanmadı.** Mongo olmadığı için cüzdan/ledger eşzamanlılık testlerinin 12'si atlandı (137 geçti · 12 atlandı · 0 başarısız); bu sonuç mali doğruluk onayı olarak kullanılamaz.
