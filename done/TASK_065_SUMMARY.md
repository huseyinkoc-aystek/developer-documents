# TASK 065 Özeti

AYPASS geçiş sisteminde Rev.06'nın onay (maker-checker) katmanı kodlandı. Yeni bir onay motoru, kiracı ve site bazlı dinamik politikalar, tek kullanımlık onay talepleri ve bunlara bağlı ödeme iadesi, B2B satın alma siparişi zinciri, banka havalesi ve sefer onayı kapıları eklendi. Eski kullanılmayan onay modülü silindi. Panelde onay politikası ve onay talebi ekranları yapıldı. Gürsel GPS yapılandırmasının yönetici değişikliklerini ezmesi sorunu da kapatıldı (`c27c3b5`).

Bu iş Rev.06 kabul kriterlerini tamamlamaz. Model 1 ve Model 4 uygulanmamıştır. Canlı ortam doğrulanmamıştır. Teknik ayrıntılar `TASK_065_AYPASS_APPROVALS.md` dosyasındadır.

**Tarih:** 2026-10-28 · **Durum:** Backend, seed ve panel işi tamamlandı ve origin'e gönderilmiş dallara aktarıldı (`74a5453`, `7b03fe7`). Yerel derleme ve test sonuçları alındı. Canlı ortam doğrulanmadı. Rev.06 kabul kriterleri henüz karşılanmıyor.

**Backend:** `/workspace/go/aystek-core`, dal `aypass-approval-ui`, taban `f48748f`, birleşik ağaç `f286fd0`, origin'e gönderilmiş dal başı `74a5453` (birleşme commit'i; P0 dalı `7eb7bfa`, `93b4232`, `25abcc1` ile birleşti).
**Panel:** `/workspace/node/aystek-core-next`, dal `approval-policy-ui`, `main` (`56907eb`) üzerine beş commit; origin'e gönderilmiş dal başı `7b03fe7` (plaka gerekçesi commit'i dahil).
**Installer belgeleri:** `/workspace/go/aystek-core-installer`, commit edilmemiş.

---

## Ne yapıldı

### 1. Onay motoru

Yeni `internal/modules/approvals` modülü kuruldu. Kapı kataloğu kodda tanımlıdır ve iki kapı içerir: `payments.refund` ve `b2b.purchase_order`. Eşik, etkinlik, para birimi, denetçi izni ve zincir adımları kiracı ve site bazlı politika verisidir. Yönetici uçları politika okuma ve güncellemeyi, onay talebi listelemeyi ve karar vermeyi sağlar.

Motor şu kuralları uygular:

- Bekleyen talepler kapı, konu ve tutar bazında tekildir (kısmi benzersiz dizin).
- Onay talebi tek kullanımlıktır. Tüketim koşullu bir güncellemedir.
- Talebi açan kişi kendi talebini onaylayamaz.
- Karar veren, talebe kaydedilmiş denetçi iznine sahip olmalıdır.
- Kararlar ve politika değişiklikleri etkinlik günlüğüne yazılır.
- Listeler sayfalıdır (varsayılan 50, en fazla 200).
- Hatalar `approvals.error.*` çeviri anahtarlarıdır (15 anahtar).

### 2. Ödeme iadesi kapısı

`payments/intents` içindeki iade akışı `RefundWithApproval` ile onay kapısından geçer. Politika kapalıysa ya da tutar eşik altındaysa iade doğrudan yapılır. Eşik üstü bir iade 409 `approvals.error.required` döndürür ve `approval_request_id` verir. Denetçi onayladıktan sonra aynı anahtar ve aynı kimlikle yapılan istek talebi tüketerek iadeyi gerçekleştirir.

Sağlayıcı iadesi başarısız olursa rezervasyon geri alınır ve onay `approved` durumuna döner. `returns` modülü onay bekleyen bu durumu başarısızlık saymaz. İade anahtarı zorunludur; anahtarsız istek 400 `payment_intents.idempotency_key_required` döner (`a00f69d`).

### 3. B2B satın alma siparişleri

Sipariş onay zinciri politikadan oluşturulur. Adımlar sıralıdır: yalnızca mevcut adım karara bağlanabilir. Reddetme de mevcut adımın rolünü ister. Onay veya red için rol istek gövdesinden okunmaz; oturumdan alınır. `submitted_by` kaydedilir ve talep sahibi kendi siparişini onaylayamaz. Politika yoksa veya zincir yapılandırılmamışsa gönderim kapalı başarısız olur (`chain_not_configured`). B2B uçları `/admin/b2b` altına taşındı ve `b2b` takma adı `b2b_dealers` modülüne çözülüyor (`de3c4fa`).

### 4. Banka havalesi maker-checker

Havale onayı `wallets:approve` izni gerektirir ve onaylayan, beyan edenden farklı olmalıdır (aynı ise 403 `wallets.error.bank_transfer_self_confirm`). Bakiye değişikliği tek bir atomik `FindOneAndUpdate` ile yapılır. Uygulanan defter kayıtlarının penceresi 50'dir. Onay izni `wallets:approve` artık tohumlanıyor (`18e27f8`).

### 5. Manuel sefer onayı

Sefer incelemesi `fleet_tally_trip_assignments:approve` iznini ister (taban commit'te `execute` idi). Onaylayan, oluşturandan farklı olmalıdır. Onay alanları (`approval_status`, `driver_profile_id`, `transferred_to_trip_id`) güncellemede kalıcı olarak yazılır.

### 6. Erişim noktaları

Pasif bir erişim noktası olay kaydını, nakit kapanışını ve yaya binişi kaydını reddeder (`access_points.error.inactive`).

### 7. Roller ve izinler

Rol anahtarı `^[a-z][a-z0-9_]{1,62}$` desenine uymalıdır ve sonradan değiştirilemez. Etiketler dil koduna göre saklanır; geri dönüş sırası dil, İngilizce, anahtardır. Dil kodu ve etiket denetimleri yapılır. Sistem rolleri değiştirilemez ve silinemez.

### 8. Faturalama

GİB taşıma hatası artık boş değil hata döndürür (`gib.submission_failed`). Sahte UUID üretimi kaldırıldı. Zaten kesilmiş fatura yeniden gönderilmez (HTTP 409). Boş sağlayıcı sonucu hata olarak döner (`invoices.provider_empty_result`). Mock e-invoice ve ERP adaptörleri yalnızca production dışında kaydedilir (`33467c8`).

### 9. Menüler

Onay menüleri rol adı yerine izinle kapılanır (`approval_policies:read`, `approval_requests:read`). Bir grup, görünür bir alt öğesi olduğu sürece korunur. Menü yazımında izinler yalnızca alan `nil` değilse yazılır.

### 10. Seed'ler

Yeni `approvals_seed` politikaları yalnızca eksikse ekler (`$setOnInsert`); mevcut politikalar üzerine yazılmaz. Modül hata çevirileri 21 anahtar ve 28 dil olarak eklendi (iki plaka gerekçesi anahtarı dahil). Onay arayüzü çevirileri ve rol çevirileri de eklendi. Seed sırasında roller onay politikasından önce gelir. `SeedRoles` editor, reviewer ve fleet_viewer rollerini yalnızca eksikse oluşturur ve izinlerine bir daha dokunmaz. `aystek-seed --roles` eksik varsayılan izinleri ekler, hiçbir şey kaldırmaz. `admin` rolü her çalıştırmada platform dışı tüm izinlere eşitlenir.

### 11. Testler

`Test{Method}_{Scenario}` biçimi ve testify ile yeniden adlandırma yapıldı (`f31c571`). Rol tohumlama için iki yeni test ve sıralı B2B zinciri için iki yeni test yazıldı. Rol tohumlama testi, eski üzerine yazma davranışı geri konduğunda başarısız olur (sahibin bildirimi).

### 12. Panel

- `/admin/approval-policies`: kapı listesi, politika düzenleme ve zincir adımı düzenleyicisi.
- `/admin/approval-requests`: bekleyen talepler, onay ve red (not ile), `?request=<id>` ile doğrudan erişim.
- Sipariş sayfasında iade bildirimi: bekleme durumu, yeniden deneme ve iptal.
- Rol sayfasında çok dilli etiket düzenleyicisi. Ad eksikse hata sayfa içinde gösterilir; `alert` kullanılmaz.
- API istemcisi ve tipler backend JSON adlarıyla eşleşir (yollar kod okumasıyla eşleştirildi).

### 13. Operasyon belgeleri

Kurulum rehberi, işletim kılavuzu ve dağıtım senaryoları ile installer README ve kullanım kılavuzu onay politikası, seed davranışı ve `--roles` akışını anlatacak şekilde güncellendi. Installer belgeleri commit edilmemiştir.

### 14. Gürsel GPS yapılandırması (`c27c3b5`)

Gürsel GPS sağlayıcı yapılandırmaları artık yalnızca ilk oluşturmada yazılır. Tam seed zinciri, güncelleme ve konteyner başlatma mevcut kayıtların `credentials`, `base_url`, `is_active`, `display_name` ve `timezone_offset` alanlarını değiştirmez. Panelden yapılan değişiklik korunur. Regresyon testi: `TestSeedGurselTurizimGpsProviders_KeepsAdminChangesOnRerun`. Bu konu artık açık risk değildir.

### 15. Araç plakası değişikliği gerekçesi (`93b4232`, panel `7b03fe7`)

Plaka değişince gerekçe zorunlu hale geldi. Gerekçe kırpılır, en fazla 500 karakter olur ve yalnızca plaka gerçekten değiştiğinde istenir. Gerekçe `stock_vehicle_plate_history.reason` alanına yazılır. Yeni çeviri anahtarları: `stock.vehicles.validation.plate_reason_required`, `stock.vehicles.validation.plate_reason_too_long` ve `stock.vehicles.fields.plate_reason` (28 dil). Panel araç formunda gerekçe alanını gösterir. Geçmiş kaydındaki yazım hatalarının yutulması açık risktir ("Bilinen eksikler ve riskler", madde 25).

---

## Temel kararlar

### Sahip kararları

- AD-04 onay kayıtları: kullanılmayan modül silinir ve yeniden tasarlanır. Onaylandı. Silme `eb581c7` ile yapıldı; yeniden tasarım bu işte kuruldu.
- Ödeme iadesi eşiği ve B2B zinciri kod değil, politika verisidir. Rol anahtarları İngilizcedir ve değişmezdir. Etiketler dile göredir. Sahibin yönlendirmesi.
- Banka havalesi ve manuel sefer onayları politika değil, izin artı kendi kendine onay kontrolü olarak kalır. Sahip bunu kabul etti.
- Seed test kimlik bilgileri kodda kalır ve üretimde değiştirilecektir. Sahip kararı.
- Editör rolünün menü görünürlüğü değişikliği onaylandı. Editör artık görmeye izinli olduğu alt öğeleri olan üç grubu görür: ticari siparişler, ticari ödemeler ve operasyon.
- Ana dil incelemesi, bir ülke projesi kapsama girene ya da iş tamamlanana kadar ertelenir. Sahip kararı.
- Sistem rolleri `editor`, `reviewer` ve `fleet_viewer` oluşturulduktan sonra panelden yönetilir; seed'ler bunların izinlerinin üzerine yazmaz. Karar 2026-10-28 tarihinde, sahibin "doğru olanı yap" yönergesi üzerine asistan tarafından verildi ve bu işte uygulandı.
- Gürsel GPS yapılandırması yalnızca ilk oluşturmada yazılır. Karar 2026-10-28 tarihinde onaylandı ve `c27c3b5` ile uygulandı.

### Tasarım kararları

- Kapı anahtarları kodda sabittir. Eşik, etkinlik ve denetçi izni veridir ve yönetici API'si ile değiştirilir.
- Politikalar kiracı ve site bazındadır. Bir site başka bir sitenin politikasını görmez.
- Denetçi izni onay talebi oluşturulurken talebe kopyalanır (snapshot). Sonradan yapılan politika değişikliği bekleyen bir kararın denetçi kümesini değiştirmez.
- Zincir gerektiren kapılar kapalı başarısız çalışır. Politika yoksa işlem durdurulur.
- Tek kullanımlık onay tüketimi koşullu güncellemedir. Eşzamanlı iki istekte yalnızca biri kazanır.
- Onay talebinin süre dolumu yoktur. Bu konu "Bilinen eksikler ve riskler" bölümünde listelenmiştir.
- Yönetici API'si menülerin izinlerini yazmaz. İzinler yalnızca seed'lerle tanımlanır.
- Politika güncellemesi (PUT) tüm politikayı değiştirir. Gönderilmeyen alanlar varsayılana döner.
- İade anahtarı zorunludur. Anahtarsız istek 400 `payment_intents.idempotency_key_required` ile reddedilir; onay ya da para hareketi başlamaz (`a00f69d`).

### Varsayılan zincir konusu

Brief'teki "varsayılan zincir yok" ifadesi aystek kiracısı için doğru değildir. Kod varsayılanında zincir yoktur ve politikası olmayan tenant ve site kapalı başarısız çalışır. `approvals_seed` varsayılan kiracının (`aystek`) varsayılan sitesi için etkin üç adımlı B2B zinciri ekler: `store_manager` (50.000 tutar eşiği), `credit_dept` (bayi kredi limiti) ve `executive` (100.000 tutar eşiği). Politikası kayıtlı olmayan veya devre dışı tenant ve site kapalı başarısız olur.

---

## Tamamlama kontrol listesi

Bu liste `README_FOR_AI_ASSISTANTS.md` içindeki tamamlama kontrol listesine göre hazırlanmıştır.

### Son doğrulama

- [x] Testler: sahibin listelediği paketler geçer. Önceden var olan iki başarısızlık (TestNewsSite*, panel kiracı kimliği testleri) belgelenmiştir.
- [x] Proje derlenir: `go build ./...` geçer.
- [x] Lint: Değişen Go dosyaları `gofmt` ile temizdir. Panelde değişen dosyalardaki üç ESLint hatası önceden var olan satırlardadır.
- [x] Ana akışa bağlandı: `cmd/api/main.go` içinde onay servisi, rotalar ve iade onay kapısı bağlıdır.
- [x] Yinelenen veya eski uygulama kalmadı: eski onay modülü silindi.

### Seed verisi

- [x] `translation_seed.go` 28 dil için güncellendi (ayrı katmanlar üzerinden).
- [x] `menu_seed.go` güncellendi: iki onay menüsü `security` grubu altında.
- [x] Alan seed dosyası oluşturuldu: `approvals_seed.go`.
- [x] Yeni UI metinleri çeviri anahtarı kullanır. Panelde sabit İngilizce metin yoktur.
- [x] Yönetici tarafından yönetilen varsayılanlar seed edildi: politikalar ve roller.
- [ ] Yeni menü anahtarlarının görünüm profillerine (RULE #2.2) kaydı bu belgede doğrulanmadı.

### Temizlik

- [x] Eski onay modülü ve governance onay dosyası silindi.
- [x] Panelde native `alert` yerine sayfa içi mesaj kullanıldı (rol sayfası).
- [ ] Kullanılmayan çeviri, menü ve izin kayıtlarının tam taraması yapılmadı.
- [x] Geçici uyumluluk kodu: yoktur.

### Belgeler

- [x] Kod yorumları: yeni ve değişen fonksiyonlarda Go ve TypeScript yorumları.
- [x] Görev belgesi: `TASK_065_AYPASS_APPROVALS.md` (1899 satır).
- [x] Özet belgesi: bu dosya.
- [x] Dosyalar `done/` klasöründe.
- [x] Operasyon belgeleri güncellendi: kurulum rehberi, işletim kılavuzu, dağıtım senaryoları, installer README ve kullanım kılavuzu. Installer belgeleri commit edilmemiştir.
- [ ] Operasyon belgelerindeki eski GPS uyarıları (`OPERATIONS_GUIDE.md` 3.15, `KULLANIM_KILAVUZU.md` 20.8) `c27c3b5` ile güncellenmedi.

### Bildirim

- [ ] Tamamlanma bildirimi yapılacak. Bu özet ve görev belgesi bildirim için kullanılabilir.
- [x] Faz 0 0.6 ve 0.11 kapatıldı (koordinatör bildirimi); 0.6 için kalan açıklar Bölüm 15.2'de.

---

## Test sonuçları

Bu sonuçlar sahibin bildirdiği şekilde, `f286fd0` ağacı üzerinde alınmıştır. `c27c3b5`, `a00f69d`, `de3c4fa`, `18e27f8` ve `33467c8` için sahibin listesinde sonuç yoktur. Bu özet yazılırken testler yeniden çalıştırılmamıştır.

### Backend

- `go build ./...` geçer.
- `gofmt`: bu işte değişen 103 mevcut Go dosyası temiz.
- `go vet ./...`: iki tanı. İkisi de bu işin dokunmadığı paketlerde: `personal_assistant/service` test dosyasında derleme hatası ve `stock/subscription_pricings/service` içinde kendine atama uyarısı.
- `go test` geçen paketler: approvals (repository, service); b2b/dealers/service; wallets (repository, service, topup_sources/service); fleet_tally/trip_assignments (repository, service); permissions (handler, service); returns/service; menus (handler, repository); access_points/service; platform/middleware; payments/intents (handler, service); billing/invoices (adapters, handler, service); orders/handler; service_reservations/service; seeds (rol, provisioning, çeviri kapsamı, menü ve profil testleri).
- Mongo bağlantılı testler yerel bir replica set'e karşı çalıştı.
- Yedi GİB sandbox testi atlanır. Taban commit'te de atlanıyordu. Düzeltme: atlama mesajlarından yalnızca biri `GIB_ENVIRONMENT` adını içerir. Diğer altısı GİB yapılandırmasının veya kimlik bilgilerinin eksikliğini belirtir.
- Önceden var olan başarısızlık: `internal/seeds` içindeki `TestNewsSite*` testleri `7c6ae1e` commit'inde de başarısızdır.

### Rota kontrolü

`cmd/api` son ağaçtan derlendi. `APP_ENV=test` ile seed'ler kapalı ve geçici bir veritabanına karşı 40 saniye çalıştırıldı. Panik ve rota çakışması yoktu. Onay rotaları kayıtlıydı.

### Panel

- `npm test`: 41 testin 39'u geçer. İki başarısızlık önceden vardır ve `main` üzerinde de vardır.
- `tsc --noEmit`: 328 hata. Bu işten önceki sayıyla aynı.
- `npm run check-i18n`: 1 ile çıkar. Eksik anahtarlar `main` ile aynıdır.
- ESLint (değişen dosyalar): üç `no-explicit-any`, önceden var olan satırlarda.

### Düzeltilen sayılar

- Sahibin "91 test fonksiyonu f31c571 ile yeniden adlandırıldı" ifadesi, taban commit'e göre 20 dosyadaki yeni test adlarının sayısıdır (95'ten 185'e). Commit'in diff'i 70 yeniden adlandırma gösterir.
- Yedi GİB atlamasının altısı `GIB_ENVIRONMENT` ile ilgili değildir (yukarıda).

### Son doğrulama (74a5453, koordinatör bildirimi)

Aşağıdakiler koordinatör tarafından bildirilmiştir. Bu özetin yazımında yeniden çalıştırılmamıştır.

- `go build ./...` geçer. `go vet ./...` yalnızca iki önceden var olan bulgu verir.
- Hedefli testler geçer: ledger, wallets, stock vehicles, access points, payments handler, b2b, middleware, invoices, seeds (çeviri, rol, GPS, izin). Ledger ve wallets eşzamanlılık testleri `-race` ile geçer.
- Route probe: `/admin/stock/vehicles` ve `/admin/b2b/purchase-orders` 401 döndü. Eski `/api/admin/b2b` 200 döndü; rota kaydı yok, yanıtın kaynağı doğrulanmadı.
- Panel (`7b03fe7`): `npm test` 39/41 (2 önceden var olan başarısızlık), değişen dosyalarda ESLint 0 sorun, `tsc` 328 (taban), yeni i18n eksiği yok.

### Bu özette iddia edilmeyenler

- Yeni GPS regresyon testinin sonucu. Test kodda vardır, ama sonucu sahibin listesinde yoktur.
- Bu özetin yazımında yeniden çalıştırılmayan testler (koordinatörün bildirdiği sonuçlar hariç).
- Canlı istekler, panel tarayıcı testleri ve Flutter derlemesi.

---

## Bilinen eksikler ve riskler

Riskler önem sırasına göre listelenmiştir. Her madde kod okumasına dayanır ya da bu özette doğrulanmamıştır. Ayrıntılar görev belgesinin Bölüm 14'ündedir.

1. B2B rota yolu (düzeltildi, doğrulanmadı): B2B uçları `de3c4fa` ile `/admin/b2b` altına taşındı; `b2b` takma adı `b2b_dealers` modülüne çözülüyor. Önceki `/api/admin/b2b` yapısında modül `api` çözülüyordu ve platform dışı kullanıcılar 403 alıyordu (kod okumasına göre). Koordinatör probe (`74a5453`): `/admin/b2b/purchase-orders` 401 (rota var); eski `/api/admin/b2b` 200 döndü, rota kaydı yok ve kaynağı doğrulanmadı. Yetkili istekle doğrulama yapılmadı; eski yolu kullanan istemciler güncellenmelidir.
2. `wallets:approve` izni artık tohumlanıyor (`18e27f8`). Uçtan uca doğrulama yapılmadı.
3. Sahte sağlayıcılar (açık risk, düzeltilmedi): mock e-invoice ve ERP adaptörleri `33467c8` ile yalnızca production dışında kaydediliyor. Mock PDF renderer, mock kargo ve mock SMS her ortamda kayıtlı kalır.
4. İade serbest bırakma kararı bekliyor: `adapter.Refund` hata döndürdüğünde rezervasyon ve onay serbest bırakılır. Zaman aşımında sağlayıcı iadeyi yapmış olabilir. Rev.06 MK-07 (PDF sayfa 27) PENDING_EXTERNAL sorgulamasını ister; zaman aşımı ret sayılmamalıdır. Karar açık (Açık kararlar, madde 4).
5. Sağlayıcı kabul ettikten sonra `repo.Update` başarısız olursa iade 400 döner. Para hareket etmiştir ama onay tüketilmiş kalır. Serbest bırakma API'si yoktur.
6. Talep tüketildikten sonra sağlayıcı çağrısından önce süreç çökerse talep tüketilmiş kalır. Kurtarma için yeni talep gerekir.
7. Bekleyen ve onaylanmış talepler süresiz kalır. Süre dolumu yoktur.
8. Onay talebi iade kaydına bağlanmaz. `approval_request_id` ödeme niyetinde, işlemde veya iade ödemesinde saklanmaz.
9. Varsayılan B2B zincirinin rolleri (`store_manager`, `credit_dept`, `executive`) hiçbir seed tarafından oluşturulmaz. Bu roller yaratılmadan eşiği aşan siparişler bekleyen kalır.
10. Varsayılan zincir: aystek kiracısının varsayılan sitesi için seed etkin üç adımlı B2B zinciri ekler; "varsayılan zincir yok" ifadesi bu kiracı için geçerli değildir (Bölüm "Varsayılan zincir konusu").
11. Boş zincirde otomatik onay: hiçbir adım uygulanmıyorsa sipariş gönderildiği anda onaylanır. Davranış belgelenmiştir.
12. Satın alma siparişi gönderiminde bayi okuma hatası yok sayılır. Hata durumunda kredi limiti 0 kabul edilir ve kredi adımı düşer (kod okumasına göre).
13. Banka havalesi onayında kısmi hata: defter kaydı işlendikten sonra bakiye uygulanmazsa yeniden onay bakiyeyi uygulamadan başarılı döner. Tutarsızlık yeniden mutabakatta görünür (kod okumasına göre).
14. Rol adları benzersiz değildir. `CreateRole` tekillik denetimi yapmaz. Roles koleksiyonunda dizin bulunmadı.
15. Eşiksiz etkin ödeme iadesi kapısı her iadeyi onaya bağlar.
16. Panelde talep listesi en fazla 100 satır ve yalnızca ilk sayfadır.
17. Sipariş sayfasında İngilizce yer tutucular (`0.00`, iade nedeni) önceden vardır.
18. Analiz belgesindeki `ensureSystemAccount` yarış durumu bu işte doğrulanmadı.
19. Test yardımcı fonksiyonları hâlâ `t.Fatalf` kullanır.
20. Flutter uygulamaları (`aypass_erisim_app`, `aypass_kullanici_app`) bu işte derlenmedi ve doğrulanmadı.
21. Üretim sunucusu (`aypass.gurseltur.com.tr`) bu ortamdan erişilemedi. Canlı durum doğrulanmamıştır.
22. Installer belgeleri commit edilmemiştir.
23. Depo `.env` dosyasını izler. Proje notlarına göre SMTP ve Firebase gizli bilgileri içerir ve bunlar gönderilmiştir. Üretimden önce döndürülmeli ve git geçmişi temizlenmelidir. Bu özette hiçbir değer yazılmamıştır.
24. Anahtar zorunluluğu (`a00f69d`): iade için `idempotency_key` göndermeyen istemciler artık 400 alır. Panel gönderir; diğer istemciler denetlenmedi.
25. Plaka geçmişi yazım hataları yutuluyor (Z02 açığı): `Create` (service.go 75) ve `UpdateVehicle` (service.go 128-129) geçmiş yazma hatalarını yok sayıyor. `AddPlate` hataları döndürür; ancak geçmişi kapatma, geçmiş ekleme ve araç güncelleme işlem sınırı olmadan ayrı yapılır.
26. Stok araçlarının `stock.vehicles.validation.plate_required`, `duplicate_plate`, `status_invalid`, `invalid_payload` ve `stock.vehicles.error.not_found` anahtarları hiçbir seed'de yok; panel ham anahtar gösterebilir (koordinatör bildirimi).
27. `internal/platform/middleware/module_gate.go` taban commit `f48748f`'de de biçimlendirilmemiş; bu dal dosyaya yalnızca alias satırını ekledi (`de3c4fa`).
### Rev.06 durumu

- Z01-Z14 arasında tamamen karşılanan gereksinim yoktur. Gereksinimler: 0 karşılanıyor, 5 kısmi, 2 yok, 7 çelişki.
- 104 kabul testi (f48748f üzerinde hesaplandı): 3 karşılanıyor, 34 kısmi, 45 yok, 22 çelişki. Bu sayılar bu işten sonra yeniden hesaplanmamıştır.
- Model 1 (DV01-DV20, DT01-DT26) sağlayıcı transferi uygulanmamıştır.
- Model 4 (M01-M24, MK01-MK30, D01-D14, F01-F12) Küresel Pazar uygulanmamıştır.
- Rev.06 analiz belgesinin §12.7 ve bazı maddeleri bu işten sonra eskimiştir (ayrıntı görev belgesinin Bölüm 14.5'inde).

---

## Açık kararlar

Dört karar açıktır:

1. `cmd/migrate_gursel` ikili dosyasının git'ten kaldırılması. Sahibin kararı bekleniyor.
2. Kimlik bilgisi rotasyonu ve git geçmişinin temizlenmesi. Üretimden önce yapılmalıdır. Kapsam: `.env` dosyası ve taban commit'teki `cmd/migrate_gursel` kaynak bağlantı bilgisi. Değerler bu özette yazılmadı.
3. Ana dil incelemesi: `ky`, `kk`, `mn`, `hy`, `az-cyrl`, `uz`, `tk`. Ertelendi.
4. İade serbest bırakma kararı (karar bekliyor): PENDING_EXTERNAL sorgulamasının tasarlanıp tasarlanmayacağı ve serbest bırakma kuralının buna göre daraltılıp daraltılmayacağı (Rev.06 MK-07, PDF sayfa 27).

---

## Sonraki adımlar

### Öncelik sırasına göre teknik işler

1. B2B rotalarının (`/admin/b2b`, `de3c4fa`) çalışan bir istekle doğrulanması. Tenant yöneticileri B2B'yi kullanmadan önce yapılmalıdır.
2. Banka havalesi onayının (`wallets:approve`, `18e27f8`) uçtan uca doğrulanması.
3. Mock PDF renderer, mock kargo ve mock SMS sağlayıcılarının production için kapatılması veya kayıttan çıkarılması.
4. Adaptör sözleşmesinin netleştirilmesi: `adapter.Refund` hatasının "para hareket etmedi" anlamına geldiğinin doğrulanması.
5. Sağlayıcı kabulünden sonraki `repo.Update` hatası için telafi mekanizması.
6. Onay talebi süre dolumu ve süresi dolmuş talepler için zamanlayıcı.
7. Refund kayıtlarında `approval_request_id` saklanması.
8. Rev.06 kabul matrisinin yeniden hesaplanması.
9. Model 1 (DV, DT) ve Model 4 (M, MK, D, F) uygulamaları.
10. Operasyon belgelerindeki eski GPS uyarılarının düzeltilmesi (`OPERATIONS_GUIDE.md` 3.15, `KULLANIM_KILAVUZU.md` 20.8).
11. Rev.06 analiz belgesi §12.7 ve ilgili maddelerinin güncellenmesi.
12. Gürsel GPS regresyon testinin bir MongoDB veritabanına karşı çalıştırılması.
13. Flutter uygulamalarının derlenmesi ve doğrulanması.
14. Üretim sunucusuna erişildiğinde canlı durum kontrolü.
15. Faz 0 0.6 kalan açıkları: `ledger_accounts` unique indeksi (veri denetiminden sonra), soft-delete ve eski+yeni binary riski, canlı replica set doğrulaması, `ConfirmBankTransfer` yeniden denemesi.
16. İade serbest bırakma kararının verilmesi (Açık kararlar, madde 4).
17. Plaka geçmişi yazım hatalarının işlenmesi (risk 25).
18. Stok araçları seed anahtarlarının eklenmesi (risk 26).
19. `module_gate.go` biçimlendirmesi (risk 27).
20. Eski `/api/admin/b2b` 200 yanıtının kaynağının bulunması (risk 1).

### Operasyon adımları (kurulum sonrası)

1. `store_manager`, `credit_dept` ve `executive` rollerini oluşturun.
2. Onay politikaları ekranında `b2b.purchase_order` zincirini ve `payments.refund` kapısını gözden geçirin.
3. İade kapısını açmadan önce bir eşik belirleyin.
4. Mevcut kurulumlarda `aystek-seed --roles` komutunu çalıştırın.
5. `chain_not_configured` (B2B gönderiminde 400) ve `approvals.error.required` (iade için 409) yanıtlarını izleyin. İkincisi beklenen davranıştır.
6. `SEED_ON_START` değerini kurulumdan sonra `false` tutun.
7. İade isteği gönderen her istemci `idempotency_key` göndermelidir.
8. B2B istemcileri `/admin/b2b` yoluna geçmelidir (`de3c4fa`).
9. Araç plakası değiştirilirken gerekçe zorunludur; eksik gerekçe `stock.vehicles.validation.plate_reason_required` ile döner.

---

## Sürümler

### Backend (`aypass-approval-ui`, origin'e gönderilmiş, HEAD `74a5453`)

Taban `f48748f` ile `f286fd0` arasında 25 commit ve 111 değişen dosya (31 eklendi, 73 değişti, 7 silindi; 9614 satır eklendi, 2431 satır silindi).

Başlıca commit'ler:

- `eb581c7` Remove unwired approvals module and unused governance approval workflow (AD-04)
- `45a2e15` Gate manual trip review on approve permission and block self-approval (AD-04)
- `54cf0bf` Gate bank transfer confirmation, PO approval and payment operations by permission
- `b96c261` Add dynamic approval policies and gate refunds above threshold (Rev.06)
- `60cd404` Read B2B approver role from the session and block submitter approval; gate returns refunds
- `269ae67` Move the B2B purchase order approval chain into admin-managed policy data
- `7c6ae1e` Give a refund's approval back when the provider never took the money
- `a712aca` Keep role keys English and immutable; store per-language role labels
- `3666773` Add approval admin menus and 28-language approval translations
- `129ed3c` Grant admin menu entries by permission as well as role
- `229c562` Require ordered approval and step-role rejection for purchase orders
- `21a8d61` Seed roles before the approvals policy
- `ccb13e2` Keep panel-managed role permissions across seed runs
- `f31c571` Name new tests Test{Method}_{Scenario} and assert with testify
- `f286fd0` Merge branch 'aypass-role-defaults' into aypass-approval-compliance
- `c27c3b5` Write Gürsel GPS provider configurations only when they are created
- `a00f69d` Require an idempotency key for admin refunds (dal başı)

### Bu turda belgeye eklenen commit'ler

- `74a5453` Merge branch 'aypass-faz0-p0-fixes' into aypass-approval-ui
- `25abcc1` Update the vehicle incidents test fake to the plate reason signatures
- `93b4232` Require a reason for every stock vehicle plate change
- `7eb7bfa` Make system account creation and journal confirmation race-safe
- `33467c8` Register mock e-invoice and ERP adapters only outside production
- `18e27f8` Grant wallets:approve in the permission seed
- `de3c4fa` Mount B2B routes under /admin and resolve them to the dealers module
- `a00f69d` Require an idempotency key for admin refunds
- `c27c3b5` Write Gürsel GPS provider configurations only when they are created
- `ccb13e2` Keep panel-managed role permissions across seed runs
- `f31c571` Name new tests Test{Method}_{Scenario} and assert with testify

### Panel (`approval-policy-ui`, origin'e gönderilmiş, HEAD `7b03fe7`)

- `e5217b2` Add approval policy and request admin screens
- `ab6fd9b` Read chain gates from the backend chain flag instead of a key list
- `e89592d` Map the backend chain flag in the approval gate client
- `dcf9d38` Replace the role page alert with an in-page message and match step cards to the card standard
- `7b03fe7` Ask for a reason when a stock vehicle plate changes

---

## Belgeler ve kaynaklar

| Belge | Konum |
|---|---|
| Görev belgesi (teknik referans) | `/workspace/aystek-core-dev-docs/done/TASK_065_AYPASS_APPROVALS.md` |
| Özet (bu belge) | `/workspace/aystek-core-dev-docs/done/TASK_065_SUMMARY.md` |
| Rev.06 analiz ve iş planı | `/workspace/aystek-core-dev-docs/AYPASS_GECIS_SISTEMI_REV06_ANALIZ_VE_IS_PLANI_TR.md` |
| Rev.06 Excel eşlemesi | `/workspace/aystek-core-dev-docs/AYPASS_REV06_EXCEL_ESLEME_TR.md` |
| Installer README | `/workspace/go/aystek-core-installer/README.md` (commit edilmemiş) |
| Installer kullanım kılavuzu | `/workspace/go/aystek-core-installer/docs/KULLANIM_KILAVUZU.md` (commit edilmemiş) |

Detay: `TASK_065_AYPASS_APPROVALS.md`.
