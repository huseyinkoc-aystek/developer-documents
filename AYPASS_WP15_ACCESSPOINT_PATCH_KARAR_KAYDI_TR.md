# WP-15 / KR-13 — AccessPoint Güncelleme Sözleşmesi Karar Kaydı

**Tarih:** 2026-10-29 · **Durum:** karar kaydı + sözleşme taslağı + kabul test planı. **Bu belgede üretim kodu yazılmadı, hiçbir dosya değiştirilmedi.**

**Referans:** backend `aystek-core` `3a3bb86` (dal `aypass-approval-ui`); panel `aystek-core-next` `c5e0994` + WP-1 (`fix/aypass-accesspoint-contract`, `a9752d6`). Kapsam: `PUT /admin/access-points/:id` → `handler.AdminHandler.Update` → `toInput` → `service.Input` → `service.Service.Update` (`internal/modules/access_points/handler/admin.go:127-149,434-487`; `internal/modules/access_points/service/service.go` `Update`).

---

## 1. Karar (KR-13)

Seçilen yaklaşım, kullanıcının verdiği beş kuralla **alan varlığını izleyen kısmi güncelleme**:

1. **Gönderilmeyen alan:** mevcut değer korunur.
2. **Açık `false`:** yetki ve iş kuralları uygunsa uygulanır.
3. **Açık boş string:** ilgili alanın temizlenmesine izin veriliyorsa temizlenir; aksi hâlde doğrulama hatası döner.
4. **`null`:** anlamı alan bazında ayrıca tanımlanmamışsa sessiz sıfırlama yapılmaz (aşağıdaki tablonun "null anlamı" sütunu doldurulmamış her alan için `null` → **doğrulama hatası**, varsayım yapılmaz).
5. **`device_bindings` için açık boş liste (`[]`):** varsayılan değer değildir; cihaz bağlarını kaldırmaya yönelik **açık ve yetkili** bir işlemdir.

Bu, mevcut tam-PUT'un (her alan koşulsuz atanır) yerini alır. **Tam-PUT korunmuyor** — kullanıcı net tercih bildirdi ("Yeni kanonik güncelleme yaklaşımı olarak..."); bu kararla "tam-PUT + eksik alana 400" seçeneği (KR-13'ün diğer şıkkı) elenmiştir.

---

## 2. İstemci envanteri

| İstemci | Durum |
|---|---|
| `admin-panel-next` (`lib/api/access-points.ts:159`, `accessPointsApi.update`) | **Tek bilinen çağıran.** WP-1 sonrası (bu oturum) her alanı **açıkça** gönderiyor (`buildAccessPointUpdatePayload`), yani yeni PATCH semantiğiyle de **otomatik uyumlu**: gönderdiği her alan "bu benim bildiğim/istediğim değer" anlamına geliyor, gönderilmeyen hiçbir alan yok (device_bindings dâhil, her zaman bir dizi olarak gönderiliyor — boş dizi de dâhil, 2.3'teki yorumla tutarlı). |
| Diğer backend modülleri / `cmd/*` araçları | `grep` ile tarandı: bu PUT ucuna giden başka bir Go çağrısı **bulunamadı**. |
| Flutter uygulamaları | `/admin/access-points/:id` PUT'u çağıran bir Flutter kodu **bulunamadı** (cihaz uygulamaları bu uca değil `/device/access-points/*`'e gider, farklı bir yüzey). |
| Swagger/OpenAPI | **Kontrol edilmedi.** Devreye almadan önce regenerasyon/karşılaştırma gerekir — ayrı kabul kriteri (bkz. §5). |
| Harici/3. parti API client | **Doğrulanmadı.** `/admin/*` insan-admin yüzeyi; api_client gateway'de bu uç yok (§ aşağıdaki not). |

> Not: `/admin/access-points/:id` yalnızca insan-admin yüzeyinde (`routes.go`); device gateway'de (`/device/access-points/:id/...`) AccessPoint'in kendisini güncelleyen bir uç **yok** — WP-2'nin kapsadığı dört uç (`events`, `cash-settlement`, `cancel`, `walk-on-cash-entry`) işlem kaydeder, AccessPoint alanlarını değiştirmez. Bu karar kaydı WP-2'yi etkilemez.

---

## 3. Sözleşme taslağı

### 3.1 Mekanizma

Go'da "anahtar hiç yok" ile "anahtar var, değeri sıfır" ayrımı salt `*T` işaretçilerle **güvenilir yapılamaz** (ikisi de `nil` işaretçiye çözülür — `encoding/json`, eksik anahtar ile açık `null`'u aynı şekilde `nil` yapar). Bu nedenle öneri:

- Handler, gövdeyi önce `map[string]json.RawMessage` olarak okur (anahtar varlığını doğrudan sorgulamak için).
- Her alan için: anahtar **yoksa** → dokunma. Anahtar **varsa** → `json.RawMessage`'ı alanın tipine çöz; `null` ise ayrıca işaretle (aşağıdaki tablo).
- Bu, mevcut DTO/`toInput` deseniyle (el yazması alan-alan dönüşüm) aynı tarzda, yeni bir jenerik `Optional[T]` tipi icat etmeden uygulanabilir — kod tabanının başka yerlerinde böyle bir genel örüntü yok, bu yüzden tek bir yeni soyutlama eklemek yerine mevcut üslupla devam etmek öneriliyor. (Alternatif: her alan için `*T` + ayrı bir `"alan" in rawMap` kontrolü — aynı sonucu verir, biraz daha fazla kod.)
- `service.Input`'a bu alanları taşımak için: her "değişti mi" bilgisini taşıyan ince bir yapı (`PatchInput` veya mevcut `Input`'a `Set<Alan>` bool bayrakları eklemek) servis katmanına geçmelidir — **servis katmanının kendisi "yetkili kaynak"** olduğundan (canonical architecture), bu ayrımı handler'da çözüp servise "zaten temiz" bir `Input` göndermek (yalnızca değişecek alanları içeren) DEĞİL; servisin kendisi hangi alanın "dokunulmadı" olduğunu bilmeli — aksi hâlde servis, HTTP katmanına özgü "var/yok" bilgisini kaybeder ve başka bir çağıran (varsa) aynı yarı-otomatik korumayı alamaz.

### 3.2 Alan bazlı tablo

| Alan | Tip | Gönderilmezse | `false`/`""`/`[]` açık değeri | `null` anlamı | Yetki (öneri) |
|---|---|---|---|---|---|
| `name` | string | korunur | `""` → **doğrulama hatası** (ad boş olamaz, mevcut `Create` kuralıyla tutarlı) | tanımsız → hata | `access_points.update` |
| `vehicle_plate` | string | korunur | `""` → **temizlenebilir** (yalnız `fare_transport` dışı kind'lerde; `fare_transport`'ta boş plaka `ErrTransportVehicleRequired` ile reddedilir — mevcut iş kuralı) | tanımsız → hata | `access_points.update` |
| `kind` | enum | korunur | — (boş string geçersiz enum, hata) | tanımsız → hata | `access_points.update` |
| `institution_id` | ObjectID\|null | korunur | — | **tanımlı:** kurumu kaldırır (bugünkü `toInput`'un `nil` davranışına en yakın karşılık) | `access_points.update` |
| `pricing_mode` | enum | korunur | — | tanımsız → hata | `access_points.update` |
| `fee` | float | korunur | `0` → **açık değer olarak uygulanır** (bugünkü `if in.Fee > 0` davranışı KALDIRILIR — D-03'ün kök nedeni budur; ücret artık 0'a indirilebilir) | tanımsız → hata (negatif de hata) | `access_points.update` |
| `currency` | string | korunur | `""` → hata (para birimi boş olamaz) | tanımsız → hata | `access_points.update` |
| `requires_approval` | bool | korunur | `false` → uygulanır | tanımsız → hata | `access_points.update` |
| `cash_fallback_enabled` | bool | korunur | `false` → uygulanır; `true` iş kuralına tabi (kurum nakit ödemeyi kapatmışsa `ErrCashPaymentsDisabledForInstitution`, mevcut davranış) | tanımsız → hata | `access_points.update` |
| `pos_enabled` | bool | korunur | `false`/`true` → uygulanır (D-02'nin kalıcı çözümü) | tanımsız → hata | `access_points.update` |
| `device_bindings` | []string | korunur | `[]` → **açık, yetkili kaldırma** (§1 madde 5) | tanımsız → hata (liste alanında null yerine `[]` kullanılır) | **ayrı yetki önerisi:** `access_points.manage_devices` (aşağıda §3.3) |
| `subscription_service_key` | string | korunur | `""` → temizlenebilir (bugünkü "boşsa yok say" davranışı KALDIRILIR) | tanımsız → hata | `access_points.update` |
| `scheduled_route_ids` | []ObjectID | korunur | `[]` → açık, yetkili kaldırma (bugünkü "nil=dokunma, []=değiştir" ayrımı zaten buna yakın; farkı: artık `[]` JSON'da **açıkça** `[]` olarak gelmeli, Go nil-slice/boş-slice karışıklığına değil, rawMap anahtar kontrolüne dayanır) | tanımsız → hata | `access_points.update` |
| `sale_product_ids` | []ObjectID | korunur | `[]` → açık kaldırma | tanımsız → hata | `access_points.update` (+ kind/pricing_mode uyumu, mevcut kural) |
| `status` | enum | korunur | — | tanımsız → hata | `access_points.update` (durum değişimi için ayrı yetki şu an yok, önerilmiyor) |

### 3.3 Alan bazlı yetki önerisi

`device_bindings`, hangi **fiziksel cihazın** bu erişim noktasında işlem yapabileceğini belirler — WP-2'nin yeni koruduğu tam da bu alan. Yalnızca `access_points.update` yetkisiyle rastgele bir admin'in bunu değiştirebilmesi, WP-2'nin sağladığı korumayı "paneli kullanabilen herkes" seviyesine indirir. **Öneri:** `device_bindings` alanının değişmesi (ekleme VEYA çıkarma) için ayrı bir yetki (`access_points.manage_devices`) şart koşulsun; bu yetkiye sahip olmayan bir çağıranın gövdesinde `device_bindings` anahtarı varsa (değer ne olursa olsun, mevcut değere eşit olsa dahi) istek **403** ile reddedilsin — "sessizce yok sayma" değil, çünkü sessiz yok sayma da bir tür belirsizlik yaratır. Diğer tüm alanlar mevcut `access_points.update` ile kalır. Bu öneri **uygulanmadı**; ayrı onay gerektirir (yeni permission seed kaydı + rol matrisi güncellemesi + panel tarafı yetki kontrolü).

### 3.4 Yanıtta açık temsil

Kullanıcının talebi: *"Yanıtta false ve boş string gibi geçerli değerler açıkça temsil edilsin; istemci başka bir alan üzerinden tahmin etmesin."* Bugün `AccessPointResponse` birçok alanı `omitempty` ile işaretliyor (WP-1 raporunda belgelendi: `vehicle_plate`, `pos_enabled`, `cash_fallback_enabled`, `fee`, `currency`, `device_bindings`, `subscription_service_key`, `scheduled_route_ids`, `sale_product_ids`). PATCH semantiği devreye alınırken **yanıt DTO'sunun `omitempty` etiketleri kaldırılmalı** (GET/PUT yanıtı her alanı her zaman içersin — `false`, `""`, `[]` dâhil) — aksi hâlde istemciler (panelin WP-1'deki `looksLikeAccessPointResponse` gibi) hâlâ "yok mu, boş mu" tahminine muhtaç kalır. Bu, backend'in **ayrı, görünür bir değişikliğidir** (yanıt şekli değişir) ve panel tarafında WP-1'in omitempty-okuma mantığının **kaldırılmasını/basitleştirilmesini** gerektirir — geçiş planına madde olarak eklendi (§4).

---

## 4. Geçiş planı

1. **Geriye dönük uyumluluk (istek tarafı):** WP-1 sonrası panel zaten her alanı açıkça gönderiyor → yeni semantikte davranışı **değişmez** (gönderilen her alan uygulanır, tıpkı şimdi olduğu gibi). Hiçbir panel kodu bu adımda kırılmaz.
2. **Davranış değişikliği (varsayılan istemciler için):** Eğer doğrulanmamış bir başka istemci bugün KISMİ gövde gönderiyorsa (envanterde bulunamadı ama kesin değil), o istemcinin "gönderilmeyen alan sıfırlanır" bekleyen hiçbir mantığı YOKTUR (bugünkü davranış istenmeyen bir yan etkiydi, kasıtlı kullanılan bir sözleşme değildi) — yeni semantik onu **daha güvenli** hâle getirir, kırmaz. Tek gerçek risk: birisi BİLEREK "alanı silmek için göndermeme" triki kullanıyorsa (olası değil, kanıt yok) davranışı değişir.
3. **Yanıt şekli değişikliği (§3.4):** `omitempty` kaldırma, panelin WP-1'deki `looksLikeAccessPointResponse`/`readOmitEmpty*` mantığını **basitleştirir** (artık "eksik=bilinen sıfır" çıkarımına gerek kalmaz, her alan her zaman gelir) — ayrı bir panel PR'ı, WP-15'in parçası değil.
4. **`device_bindings` ayrı yetki (§3.3) kabul edilirse:** panelin POS benzeri bir "cihaz yönetimi" ekranı/izni olmayan rollerde bu alanı **göndermemesi** gerekir (bugün form her zaman `device_bindings: form.device_bindings` gönderiyor — yetkisiz bir admin bunu değişmemiş haliyle geri gönderse bile reddedilir, §3.3'teki kural). Panel tarafında ayrı bir iş kaydı.
5. **Yöntem adı:** URL ve HTTP metodu (`PUT /admin/access-points/:id`) **değişmiyor** — yalnızca gövde yorumlama semantiği değişiyor. Bir `PATCH` takma adı eklemek (aynı handler'a) istenirse ayrı, küçük bir karar.
6. **Swagger/OpenAPI:** regen edilmeli; bu belge kapsamında **yapılmadı**.

### 4.1 Çözülmeyen risk: "eski veriyi geri yazarak ezme" (lost update)

Kullanıcının riski: *"Panelin eski GET verisini geri göndererek arada yapılmış bir cihaz bağı/POS/onay ayarı değişikliğini ezmesi riski ele alınsın."*

**Bu karar kaydı bunu TAM çözmüyor, dürüstçe işaretliyorum:** Yeni sözleşim "gönderilmeyen alan korunur" der, ama **panel (WP-1 sonrası) her alanı HER ZAMAN gönderiyor** — kullanıcı formda sadece adı değiştirse de, `device_bindings`/`pos_enabled` gibi alanları da (kendi yüklediği, muhtemelen BAYATLAMIŞ GET anlık görüntüsünden) geri gönderiyor. Yani:

- Backend sözleşmesi (bu belge) "gönderilmeyeni koru" diyor → **doğru ve gerekli bir temel.**
- Ama panel hâlâ HER ŞEYİ gönderdiği için, iki admin aynı AccessPoint'i aynı anda açıp biri `device_bindings`'i değiştirirken diğeri `name`'i değiştirip kaydederse, ikinci kayıt birincinin `device_bindings` değişikliğini **kendi bayat kopyasıyla ezer** — PATCH semantiği bunu **önlemez**, çünkü panel teknik olarak "göndermiyor" değil, "(bayat) değerini gönderiyor".
- **Gerçek çözüm iki parçalı:** (a) bu belgedeki backend sözleşmesi (alan varlığı izleme) — **gerekli ön koşul**; (b) panelin "yalnızca kullanıcının GERÇEKTEN dokunduğu alanları gönder" davranışına geçmesi (optimistic/dirty-field diffing) — **ayrı, yapılmamış bir panel işi**. (b) olmadan (a) tek başına bu riski kapatmaz; sadece "backend'in varsayılan olarak her şeyi sıfırlaması" riskini kapatır (D-01/D-02'nin kök nedeni), "panelin bayat veriyle geri yazması" riski AYRI kalır.
- Alternatif/ek çözüm: `updated_at`/sürüm alanına dayalı optimistik kilit (If-Match benzeri) — bu belgede **önerilmedi, kararlaştırılmadı**, ayrı bir karar kaydı gerektirir.

---

## 5. Kabul test planı (WP-15 devreye alındığında)

Tüm senaryolar backend `Service.Update` ve `handler.Update` seviyesinde, mevcut `service_test.go` desenine uygun sahte repo'larla:

| # | Senaryo | Beklenen |
|---|---|---|
| 1 | Yalnız `name` gönder | diğer tüm alanlar (incl. `vehicle_plate`, `pos_enabled`, `device_bindings`, `requires_approval`, `cash_fallback_enabled`) **değişmez** |
| 2 | `pos_enabled: false` açıkça gönder (önceden `true`) | `false` **uygulanır** |
| 3 | `vehicle_plate: ""` gönder, kind≠`fare_transport` | **temizlenir** |
| 4 | `vehicle_plate: ""` gönder, kind=`fare_transport` | **doğrulama hatası** (mevcut `ErrTransportVehicleRequired` kuralıyla tutarlı) |
| 5 | `fee: 0` açıkça gönder | **uygulanır** (D-03 kapanır) |
| 6 | `device_bindings: []` açıkça gönder | **tüm bağlar kaldırılır** (açık işlem olarak) |
| 7 | `device_bindings` anahtarını hiç gönderme | bağlar **değişmez** |
| 8 | `institution_id: null` açıkça gönder | kurum **kaldırılır** |
| 9 | Tanımsız bir alanda (örn. `name: null`) `null` gönder | **doğrulama hatası**, sessiz sıfırlama yok |
| 10 | (§3.3 kabul edilirse) `access_points.manage_devices` yetkisi olmayan çağıran, `device_bindings` alanını (değişmemiş değeriyle bile) gönderir | **403**, hiçbir alan uygulanmaz (yalnız `device_bindings` değil, bütün istek) |
| 11 | Yanıt (`GET`/`PUT` dönüşü) — `pos_enabled: false`, `vehicle_plate: ""`, `device_bindings: []` olan bir kayıt | yanıt gövdesinde bu alanlar **anahtar olarak mevcut** (omitempty kaldırıldıktan sonra) |
| 12 | Eşzamanlı iki güncelleme: A yalnız `device_bindings`'i değiştirir, B (A'dan SONRA okuyup) yalnız `name`'i değiştirip A'nın değişikliğinden SONRA kaydeder, B'nin gövdesinde `device_bindings` **anahtarı yok** | A'nın `device_bindings` değişikliği **korunur** (bu, panel "yalnız değişeni gönder"e geçtiğinde geçerli olur — §4.1; panel geçmeden bu test yalnızca BACKEND sözleşmesini kanıtlar, uçtan uca riski kapatmaz) |
| 13 | Mevcut 91+ `access_points` testi | regresyon yok |
| 14 | Panel WP-1 sonrası payload'u (her alan açık) aynı backend'e karşı | **davranış aynı kalır** (geçiş planı madde 1'in kanıtı) |

---

## 6. Bu belgede yapılmayanlar

Üretim kodu yazılmadı; `Input`/`toInput`/`Update`/DTO değiştirilmedi; permission seed'i değiştirilmedi; swagger regen edilmedi; panelin "yalnız değişeni gönder" davranışına geçişi yapılmadı (§4.1 açık risk olarak kaldı). WP-1 ve WP-2'nin üretim koduna bu belgeyle **hiçbir dokunuş yapılmadı**.
