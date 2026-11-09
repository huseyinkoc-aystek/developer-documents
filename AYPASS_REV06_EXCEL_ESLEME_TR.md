# AYPASS Rev.06 — Excel Kimlik Eşlemesi ve Kod Durumu

- **Tarih:** 2026-10-28
- **Kaynak:** `AYSTEK_Gecis_Sistemi_Rev06-xlsx.pdf` (Excel'in PDF dökümü, 41 sayfa; kullanıcı iletti). Asıl `.xlsx` dosyası yüklenemedi.
- **Kod referansı:** `/workspace/go/aystek-core`, HEAD `f48748f`. Çalışma ağacında yalnız `.env` değişik; commit'e dahil edilmez.
- **Ana analiz:** `AYPASS_GECIS_SISTEMI_REV06_ANALIZ_VE_IS_PLANI_TR.md`
- **Yöntem ve sınır:** Kod okuma, grep ve beş denetçi raporu. Kritik iddialar bu oturumda bağımsız doğrulandı (BAĞIMSIZ). **Hiçbir test çalıştırılmadı, hiçbir kod değiştirilmedi, commit yapılmadı.** Canlı veritabanı ve üretim GİB yapılandırması doğrulanmadı.

## Durum sözlüğü

| Durum | Anlamı |
|---|---|
| KARŞILANIYOR | Kod kabul kriterini karşılıyor |
| KISMİ | Parça var; kriterin bir kısmı eksik |
| YOK | Özellik yok |
| ÇELİŞKİ | Kod var ama kabul kriterine aykırı ya da sessiz hata üretiyor |
| KARAR | İş tarafı kararı bekleniyor (Excel: "Karar / değer bekliyor") |

## 1. Excel sayımları

| Grup | Sayı | Kimlikler |
|---|---|---|
| Z (Dikkat Edilmesi Zorunludur) | 14 gereksinim · 24 kabul testi · 11 açık karar | Z01–Z14 · K01–K24 · P01–P11 |
| Market | 24 tasarım kuralı · 30 kabul senaryosu · 14 karar | M01–M24 · MK01–MK30 · D01–D14 |
| Market akış ve veri | 12 işlem akışı · 13 veri varlığı · 5 temsili hakediş örneği | F01–F12 |
| Devir (Model 1) | 20 gereksinim · 10 iş adımı · 12 karar · 26 test · 16 veri alanı grubu | DV01–DV20 · DK01–DK12 · DT01–DT26 |
| Ulaşım | 6 algoritma satırı · 25 güzergâh ayarı · 15 kontrol senaryosu · 9 formül örneği | F0, Ücretsiz, A1–A4 · 01–25 · T01–T15 · FX1–FX9 |
| Kaynak | 13 | S01–S11, R06-S1, R06-S2 |

Sayılar Excel'in kendi notlarıyla tutarlı. Excel içinde kimlik çakışması bulunmadı.

## 2. Gereksinimler — kod durumu

### 2.1 Z01–Z14 (Dikkat Edilmesi Zorunludur)

| ID | Gereksinim (özet) | Durum | Ana neden |
|---|---|---|---|
| Z01 | Ücret bağlamı araçtan ayrı; tarife belirsizse olay kaydı | ÇELİŞKİ | Araçsız yol yok; hata olur, olay yazılmaz |
| Z02 | Sabit araç ID; plaka değişikliği tarihçeli ve gerekçeli | KISMİ | Stok plaka tarihçesi var; gerekçe ve aktör yok; AP ve alarm plakaları güncellenmiyor |
| Z03 | Bakiye kişi hesabında; kart değişimi bakiyeyi korur | KISMİ | Kart↔cüzdan bağı sabit; geçerlilik tarihi yok; test bakiye doğrulamıyor |
| Z04 | Karta yazma yok; QR/kart kuralları | KISMİ | Karta yazma yok ✔; QR tek kullanım ve revoke yok; çevrimdışı kontrol yok |
| Z05 | Çevrimdışı kalıcı kuyruk; kabul ≠ tahsilat | YOK | Kuyruk yok; hata "Geçiş Reddedildi" olarak görünür |
| Z06 | Benzersiz olay anahtarı; olay zamanı ve sıra | KISMİ | Idempotency anahtarı var ama isteğe bağlı; istemci göndermiyor; olay zamanı ve sıra yok |
| Z07 | Online alt sınır işlem sonrası; çevrimdışı alacak izlenir | ÇELİŞKİ | Alacak kaydı yok; eşzamanlı aşım riski |
| Z08 | Algoritmalar çevrimdışıda; sırasız olay eşleşir | ÇELİŞKİ | Sıra bağımlı; A3/A4 yok; otomatik borçlandırma |
| Z09 | Tek mali kaynak; silinmez; ters kayıt | KISMİ | Journal transaction ✔; ters kayıt yok; cüzdan cache'i doğrudan yazılabiliyor |
| Z10 | Tekil referans; kesin hassasiyet; eşzamanlılık | ÇELİŞKİ | Tekrar bildirim idempotent ✔; float64; eşzamanlı güncelleme kaybı |
| Z11 | Fatura işi kalıcı; belirsiz yanıt sorgulanır | ÇELİŞKİ | İş kuyruğu yok; GİB taşıma hatasında sahte başarı; çift fatura |
| Z12 | Kasa, güzergâh ve araç raporları ayrışır | ÇELİŞKİ | Reddedilen ve iade edilen olaylar toplama giriyor |
| Z13 | Sabit ID; ad/plaka ilişki anahtarı değil | ÇELİŞKİ | Plaka anahtar olarak kullanılıyor; olay anı kopyaları eksik |
| Z14 | Göç ve canlıya çıkış kontrolü | YOK | Göç aracı yalnız stok; mutabakat ve kesim kaydı yok |

Özet: 0 karşılanıyor · 5 kısmi · 2 yok · 7 çelişki.

### 2.2 M01–M24 (Market tasarım kuralları)

| ID | Konu | Durum | Ana neden |
|---|---|---|---|
| M01 | Ağ ve program (network_id, program_id) | YOK | Ağ/program kavramı yok |
| M02 | Finansal rol ayrımı (issuer, provider, legal_entity) | YOK | Bu roller modelde yok |
| M03 | Kişi hesabı ve kimlik | KISMİ | Kart↔hesap bağı ✔; kurumlar arası birleştirme yok |
| M04 | Ayrı cüzdan/cep | YOK | Kişi + para birimi başına tek cüzdan |
| M05 | Üye işyeri kabulü (merchant > şube > terminal) | KISMİ | institutions hiyerarşisi ve AP var; tüzel kişilik, sözleşme, IBAN yok |
| M06 | Kural önceliği (policy_version) | YOK | Politika sürümü yok |
| M07 | Uygun sepet denetimi | YOK | Uygunluk denetimi yok |
| M08 | Kart ve QR ile ödeme onayı | ÇELİŞKİ | UID ile borçlandırma mümkün (ÇIKARIM: istemci akışı) |
| M09 | Dinamik ve bağlı QR | KISMİ | TTL 90 sn ✔; tek kullanım ve tutar/işyeri bağı yok |
| M10 | Bakiye ve eşzamanlılık (hold, atomik kontrol) | ÇELİŞKİ | Hold yok; okuma-hesap-yazma yarışı |
| M11 | Ödeme yaşam döngüsü | KISMİ | Intents durumları var; REFUNDED ve PENDING_EXTERNAL yok |
| M12 | Tekilleştirme ve kesinti (sorgulama) | KISMİ | Op bazlı idempotency ✔; sağlayıcı sorgusu yok; çatışma kontrolü yok |
| M13 | Tek mali kaynak, ters kayıt | KISMİ | Journal var; ters kayıt yok |
| M14 | Dış ortak ve yerel defter mutabakatı | YOK | Dış durum ve mutabakat yok |
| M15 | Parçalı ödeme | YOK | Bölüşümlü ödeme yok |
| M16 | Tam/kısmi iade (kümülatif ≤ tahsilat) | ÇELİŞKİ | Üst sınır yok (P0-6) |
| M17 | Market hakedişi | YOK | Hakediş modeli yok |
| M18 | Banka ödemesi ve mutabakat | YOK | Banka ödeme akışı yok |
| M19 | Belge ve fatura profili | KISMİ | Genel fatura var; profil yok; GİB sorunları |
| M20 | Market çevrimdışı profili (varsayılan kapalı) | KISMİ | Çevrimiçi varsayılan ✔; sepet korunmuyor |
| M21 | Eksi bakiye ve risk sınırı | KISMİ | Taban kurum politikasından (varsayılan 0 ✔); market kuralı yok |
| M22 | Veri ve erişim sınırı | KISMİ | Kurum kapsamı ✔; merchant kapsamı ve alan maskeleme yok |
| M23 | Rapor ve performans (as_of, business_date) | ÇELİŞKİ | as_of yok; toplam tanımı yanlış |
| M24 | Canlıya çıkış kapıları | YOK | Kapı mekanizması yok |

Özet: 0 karşılanıyor · 10 kısmi · 10 yok · 4 çelişki.

Market işlem akışları (F01–F12) ve veri modeli (13 varlık) için: F01 YOK, F02 KISMİ, F03 KISMİ, F04 ÇELİŞKİ, F05 YOK, F06 KISMİ, F07 YOK, F08 KISMİ, F09 ÇELİŞKİ, F10 YOK, F11 KISMİ, F12 KISMİ. Veri modelinde 13 varlıktan 7'si kısmi, 1'i çelişki (payment_allocation/refund), 5'i yok (network/program, policy, outbox/inbox, settlement, fiscal/invoice_job). Cep (pocket) account satırında "kısmi" olarak işaretli; cep modeli yok.

### 2.3 DV01–DV20 (Model 1 devir gereksinimleri)

| ID | Gereksinim (özet) | Durum | Ana neden |
|---|---|---|---|
| DV01 | Devir kapsamı kurum + hizmet + sözleşme | YOK | Devir/atama modeli yok |
| DV02 | Firma kimlikleri korunur, yeni provider_id | YOK | Provider kavramı yok |
| DV03 | Sürümlü hizmet ataması, [başlangıç, bitiş) | YOK | Emsal yalnız trip_contracts (rota/araç kapsamı yok) |
| DV04 | Devir dosyası ve onay | YOK | approvals bağlı değil |
| DV05 | Kesim anı ve sefer (T, UTC, saat dilimi) | YOK | Sefer T'ye bağlı değil |
| DV06 | Kişi hesabı ve kart devir sırasında korunur | KISMİ | Kişi-merkezli cüzdan ✔; firma bağı yok |
| DV07 | Pozitif bakiye ve fon sorumlusu | YOK | Fon sorumlusu yok |
| DV08 | Borç, promosyon, abonman ayrı envanter | YOK | Tek nakit bakiye |
| DV09 | Güvenli mali devir (transfer_batch) | YOK | Transfer yok |
| DV10 | Açık yolculuk T öncesi eski bağlamla kapanır | YOK | A3/A4 yok |
| DV11 | Gecikmiş offline olay eski bağlamla | YOK | Çevrimdışı yok |
| DV12 | Cihaz ve kesintisizlik (profil, kuyruk) | YOK | Profil ve kuyruk yok |
| DV13 | Kesim sınırında ödeme (T öncesi/sonrası) | YOK | T çapası yok |
| DV14 | POS ve tahsilat profili | YOK | Merchant/terminal profil devri yok |
| DV15 | İade ve harcama itirazı (original_transaction_id) | KISMİ | Olay ve intent referansı var; kanal ve kısmi iade yok |
| DV16 | Fatura ve belgeler özgün tarafa bağlı | YOK | Belge sahibi modellenmiyor |
| DV17 | Yetki ve veri ayrımı | KISMİ | Kurum ACL kısmi; devir sonrası yetki kapanışı yok |
| DV18 | Mutabakat ve denetim | KISMİ | Reconcile var; düzeltme sürümü yok |
| DV19 | Operasyon ve mali kapanış ayrımı | YOK | Durum makinesi yok |
| DV20 | Geri dönüş ve tekrar | YOK | Plan ve geri dönüş yok |

Özet: 0 karşılanıyor · 4 kısmi · 16 yok. Devir iş akışı (10 adım, durum makinesi, 6 onay kapısı): tamamı YOK.

## 3. Kabul testleri — test bazlı eşleme

Durumlar kod okumasına dayanır. Hiçbir test çalıştırılmadı.

### 3.1 K01–K24 (Operasyon kabul testleri, Z)

| ID | Gereksinim | Senaryo (kısa) | Durum | Gerekçe |
|---|---|---|---|---|
| K01 | Z01 | Araç ataması yok; doğru tarife bulunur; geçiş bloke olmaz | ÇELİŞKİ | Araçsız yol yok; plaka ve aktif sefer zorunlu, eksikse 400 |
| K02 | Z01 | Tarife belirsiz; olay kaybolmaz; istisna durumu | ÇELİŞKİ | Hata döner; AccessEvent, ledger, RiderTrip yok; istisna durumu yok (sıfır/rastgele ücret yok, iyi) |
| K03 | Z02 | Yanlış plaka düzeltmesi; tarihçe ve gerekçe | KISMİ | Stok plaka tarihçesi var; gerekçe yok; aktör istek gövdesinden |
| K04 | Z03 | Kart kaybı/yenileme; hesap aynı kalır | KISMİ | Bağ sabit; test bakiyeyi doğrulamıyor |
| K05 | Z03 | Gecikmiş eski kart olayı → okutma anındaki hesap | YOK | Çevrimdışı olay yok; eski kart reddedilir |
| K06 | Z04 | Karta veri yazılmaz | KARŞILANIYOR | NFC yalnız UID okur (yalnız çevrimiçi) |
| K07 | Z04 | QR tekrarı, iptal, önbellek | KISMİ | 90 sn TTL ✔; tek kullanım ve revoke yok |
| K08 | Z05 | Sunucu kesintisinde yerel kayıt | YOK | Yerel kuyruk yok |
| K09 | Z05 | Uygulama yeniden başlar; kuyruk korunur | YOK | Kuyruk yok |
| K10 | Z06 | Yanıt kaybı; aynı sonuç, tek mali etki | KISMİ | idempotency_key var (isteğe bağlı); AYPASS uygulamaları göndermiyor |
| K11 | Z06 | Tarife sonradan değişti; eski tarife geçerli | ÇELİŞKİ | Açık yolculuk çıkış anı tarifesini alır (time.Now); olayda tarife sürüm kimliği yok |
| K12 | Z07 | −15 sınırı; 16 TL onaylanmaz; eşzamanlı aşım yok | ÇELİŞKİ | Tek istek doğru; okuma-hesap-$set nedeniyle eşzamanlı aşım mümkün |
| K13 | Z07 | Çevrimdışı 120 TL alacak ayrı izlenir | YOK | Çevrimdışı yok; insufficient yalnız AccessEvent, alacak kaydı yok |
| K14 | Z07 | Birden fazla offline cihaz; risk görünür | YOK | Çevrimdışı yok |
| K15 | Z08 | A3/A4 iniş önce gelir; tek yolculuk | ÇELİŞKİ | A3/A4 yok; metered "sonraki dokunuş = iniş" erken inişi yeni biniş sayar |
| K16 | Z09 | Açılış + hareketler = bakiye; silinmez | KISMİ | Posted journal silinmez ✔; açılış modeli yok; Reconcile pending'i sayar |
| K17 | Z09 | Rapor yeniden kurulur; aynı sonuç | KISMİ | Rapor on-demand ✔; toplamlar refunded ve denied ayırmıyor |
| K18 | Z10 | Ödeme bildirimi iki kez → tek yükleme ve tek fatura işi | KISMİ | Yükleme idempotent ✔ (payment-intent:<id>, topup:<code>:<txid>); fatura işi yok |
| K19 | Z10 | −120 + 200 = 80; eşzamanlı işlem kaybolmaz | ÇELİŞKİ | Aritmetik doğru; eşzamanlı güncelleme kaybı (P0-4) |
| K20 | Z11 | Fatura işi oluşurken kesinti; yükleme kaybolmaz | YOK | Fatura işi yok; fatura senkron ve admin tetikli |
| K21 | Z11 | Fatura oluştu, yanıt kayboldu; sorgula, tekrar etme | ÇELİŞKİ | GİB taşıma hatasında sahte UUID ile başarı; ExternalID kontrolü yok |
| K22 | Z11 | İşçi yarıda kaldı; iş devralınır | YOK | Worker ve lease yok |
| K23 | Z12/Z13 | Plaka veya personel türü değişti; geç olay; doğru sayım | ÇELİŞKİ | FindByPlate current_plate; plaka değişince eski sefer bulunamaz; billing sessiz atlar |
| K24 | Z14 | Geçiş provası; mutabakat; kesim zamanı | YOK | Göç aracı yalnız stok; mutabakat yok |

Özet: 1 karşılanıyor · 7 kısmi · 8 yok · 8 çelişki.

### 3.2 MK01–MK30 (Market kabul senaryoları)

| ID | Gereksinim | Senaryo (kısa) | Durum | Gerekçe |
|---|---|---|---|---|
| MK01 | M01/M06 | Güzergâhsız market satışı | KISMİ | Yemekhane satış yolu (sale_items) var; market ve cep yok |
| MK02 | M01/M22 | Kurumlar arası ortak market | YOK | Ağ ve program kavramı yok |
| MK03 | M03 | Kart değişimi; bakiye aynı | KISMİ | Kod uyumlu; test bakiye doğrulamıyor |
| MK04 | M04 | Yanlış cep | YOK | Cep yok |
| MK05 | M05/M06 | Pasif işyeri veya terminal işlem yapamaz | ÇELİŞKİ | RecordEvent AccessPoint.Status denetlemez (BAĞIMSIZ, grep) |
| MK06 | M07 | Karma sepet; uygun kısım ayrılır | YOK | Uygunluk denetimi yok |
| MK07 | M08 | Kart seri no ile tahsilat yapılmaz | ÇELİŞKİ | UID doğrulanıp wallet_id ile borçlandırılıyor (ÇIKARIM) |
| MK08 | M09 | Tekrarlanan veya süresi geçmiş QR | KISMİ | TTL ✔; 90 sn içinde çok kullanım |
| MK09 | M09 | Tutar veya işyeri değişince eski onay geçersiz | YOK | QR tutar ve işyerine bağlı değil |
| MK10 | M10 | İki kasada eşzamanlı ödeme; e-para eksiye düşmez | ÇELİŞKİ | Okuma-hesap-yazma; hold yok |
| MK11 | M10/M21 | Yetersiz bakiye varsayılan ret | KISMİ | Yetersiz olay kaydı ve ret ✔; kurum tabanı negatif olabilir |
| MK12 | M12 | Yanıt kaybı → sorgula; ikinci capture yok | KISMİ | Op bazlı idempotency; sağlayıcı sorgusu yok |
| MK13 | M12 | Tekrarlı webhook tek mali kayıt | KISMİ | Callback idempotent ✔ (test); hakediş kalemi yok |
| MK14 | M12 | Aynı anahtar, farklı tutar → çatışma | ÇELİŞKİ | Eski sonuç sessizce döner (intents ve ledger) |
| MK15 | M13/M14 | Yerel kayıt kesintisi → dış referansla kurulum | YOK | Outbox yok |
| MK16 | M11 | Capture öncesi iptal | KISMİ | Cancel op var; yalnız Stripe destekli |
| MK17 | M16 | Tam iade, orijinal kaynağa ters kayıt | KISMİ | Refund var; ters kayıt yok |
| MK18 | M16 | Kümülatif iade tahsilatı aşamaz (60 + 50 > 100) | ÇELİŞKİ | Üst sınır yok (P0-6) |
| MK19 | M15/M16 | Bölüşümlü cep iadesi | YOK | Cep yok |
| MK20 | M15 | Ödeme parçası başarısız | YOK | Bölüşümlü ödeme yok |
| MK21 | M17/M18 | Dönem sonu mutabakat ve hakediş | YOK | Hakediş yok |
| MK22 | M18 | Banka gönderildi, yanıt yok | YOK | Banka ödeme akışı yok |
| MK23 | M18 | Ödenmiş hakediş sonrası iade | YOK | Hakediş yok |
| MK24 | M19 | Yükleme ve satış belgesi ayrımı | KISMİ | invoiceable_records kaynağa göre ayırır; belge profili yok |
| MK25 | M19 | Fatura oluştu, yanıt yok | ÇELİŞKİ | GİB sahte başarı; sorgu yok |
| MK26 | M20 | Çevrimdışı market: başarı yok, sepet korunur | KISMİ | Başarı verilmiyor ✔; sepet korunmuyor |
| MK27 | M20/M21 | Çevrimdışı kayıttan sonra borçlandırma yok | YOK | Market çevrimdışı yok (ulaşımdaki otomatik borçlandırma C-23 ayrı risk) |
| MK28 | M22 | Kişisel alışveriş gizliliği | KISMİ | Kurum kapsamı var; SaleItems maskesiz |
| MK29 | M23 | Rapor yeniden üretilir, toplamlar aynı | KISMİ | Deterministik; toplam tanımı yanlış (refunded dahil) |
| MK30 | M24 | Pilot ve geri dönüş | YOK | Mekanizma yok |

Özet: 0 karşılanıyor · 12 kısmi · 12 yok · 6 çelişki.

### 3.3 DT01–DT26 (Devir kabul testleri)

| ID | Gereksinim | Senaryo (kısa) | Durum | Gerekçe |
|---|---|---|---|---|
| DT01 | DV01 | Kapsam izolasyonu | YOK | Devir modeli yok |
| DT02 | DV02 | Yeni firma provider_id; Gürsel kaydı silinmez | YOK | Provider kavramı yok |
| DT03 | DV03 | T-1 eski, T yeni atama; çakışma yok | YOK | Sürümlü atama yok |
| DT04 | DV04 | Onay eksik; etkinleştirme bekler | YOK | approvals bağlı değil |
| DT05 | DV06 | Hesap ve kart korunur | KISMİ | Kişi-merkezli cüzdan ✔; firma bağı yok |
| DT06 | DV07 | Bağımsız cüzdan; yeni kabul ilişkisi | KISMİ | Cüzdan bağımsız; kabul ilişkisi yok |
| DT07 | DV07 | Eski firma bakiyesi harcanabilir sayılmaz | YOK | Fon sorumlusu yok |
| DT08 | DV09 | Parti tekrarı; tek kredi | YOK | Transfer batch yok |
| DT09 | DV09 | Yarım devir; çift kullanım yok | YOK | Transfer yok |
| DT10 | DV08 | Eski borç onaysız devredilmez | YOK | Borç ayrı envanter yok |
| DT11 | DV08 | Promosyon ve abonman ayrımı | YOK | Tek nakit bakiye |
| DT12 | DV10 | A3/A4 kesimi aşar; ikinci tahsilat yok | YOK | A3/A4 yok |
| DT13 | DV05 | A1/A2 akşam; sefer T-1 başlar | YOK | Sefer T'ye bağlı değil |
| DT14 | DV11 | Geç offline olay eski bağlamla | YOK | Çevrimdışı yok |
| DT15 | DV11/12 | Eski profil kullanımı | YOK | Profil süresi yok |
| DT16 | DV11 | Güvenilmez saat → inceleme | YOK | Cihaz saati kullanılmıyor; inceleme yok |
| DT17 | DV12 | Cihaz yeniden atama; kuyruk korunur | YOK | Kuyruk yok |
| DT18 | DV13 | Geç yükleme onayı; tek kredi | KISMİ | Top-up idempotent ✔; T çapası yok |
| DT19 | DV13/14 | POS sonucu belirsiz → aynı referansla sorgu | YOK | Sorgu yok |
| DT20 | DV14 | Yeni POS yönü; merchant ve terminal onaylı | YOK | Profil devri yok |
| DT21 | DV15 | Eski kısmi iade; tek iade | YOK | Kısmi iade kanalı yok |
| DT22 | DV16 | Fatura yanıt kaybı; devirde ikinci fatura yok | ÇELİŞKİ | GİB sahte başarı |
| DT23 | DV17 | Yetki sınırı; kapsam dışı erişim yok | YOK | Firma yetkisi yok |
| DT24 | DV18 | Geç olay ve rapor; düzeltme sürümü | YOK | Düzeltme sürümü yok |
| DT25 | DV19 | Operasyon devredildi; mali kapanış açık | YOK | Durum makinesi yok |
| DT26 | DV20 | Mali işlem sonrası geri dönüş | YOK | Plan yok |

Özet: 0 karşılanıyor · 3 kısmi · 22 yok · 1 çelişki.

### 3.4 T01–T15 (Ulaşım kontrol senaryoları) ve FX1–FX9 (formül örnekleri)

| ID | Senaryo (kısa) | Durum | Gerekçe |
|---|---|---|---|
| T01 | A1 sabah: biniş tek okutma; mesafe tarifesi ve indirim; biniş anında tahsilat | ÇELİŞKİ | Metered ücret çıkışta alınır; biniş tahsilatı yok; sabit uç yok |
| T02 | A1 akşam: kampüsten çıkış; iniş okutması; akşam mesafesi ayrı | ÇELİŞKİ | Akşam tek okutma yönsüz; biniş sayılır (ÇIKARIM) |
| T03 | A2 durak: sabit fiyat; mesafeden üretilmez | KISMİ | O-D matrisi var; biniş anında kullanılmıyor |
| T04 | A3 indirimsiz: P tahsili; N = brüt; R = P − N | KISMİ | Net hesap var (çıkış); P ve R yok |
| T05 | A3 indirimli: P=100, brüt=60, %20 → N=48, R=52 | KISMİ | N hesabı var; P ve R yok |
| T06 | A4 matris: A>B=50, tavan 80, %20 → N=40, R=40 | KISMİ | N=50×0,8 ✔; P ve R yok |
| T07 | A4 ters yön: A>B ve B>A farklı; simetri yok | KARŞILANIYOR | Yönlü O-D; TestReverseDirectionDoesNotCollide |
| T08 | A3/A4 R=0: iade yok, ikinci tahsilat yok | KISMİ | Tek ücret çıkışta; R kavramı yok |
| T09 | Tekrar okutma tek etki | ÇELİŞKİ | İstemci anahtar göndermiyor; iki olay, iki tahsilat (ÇIKARIM) |
| T10 | İkinci okutma yok: politika güzergâhta seçilir, süre ve ceza varsayılmaz | ÇELİŞKİ | Sefer bitince planlı mesafe × oran otomatik borçlandırılır |
| T11 | Yetersiz bakiye: tahsil edilmiş gibi gösterilmez | ÇELİŞKİ | Otomatik kapanışta Debit hatası yutulup ChargedAmount yazılır |
| T12 | İade hatası: bekleyen durum; güvenli yeniden deneme | YOK | Bekleyen durum yok; toplu iade hataları yutulur |
| T13 | Eksik tarife: sıfır veya tahmini fiyat yok | KISMİ | ErrActiveRouteFareRequired ✔; durak eşleşmezse km yoluna sessiz düşüş |
| T14 | N > P: negatif iade yok; incelemeye alınır | YOK | Tavan kavramı yok |
| T15 | Mercedes/LCW ücretsiz takip ve fabrika kayıtları korunur | ÇELİŞKİ | Ulaşım ücretsiz modu plaka + aktif sefer + SUL ister |

Özet (T): 1 karşılanıyor · 6 kısmi · 2 yok · 6 çelişki.

| ID | Formül örneği | Durum | Gerekçe |
|---|---|---|---|
| FX1 | F0: 50, %10 → net 45 | KARŞILANIYOR | net = brüt × (1 − %) |
| FX2 | A1 sabah: 12 km × 4 = 48, %20 → 38,40 | KISMİ | Aritmetik ✔; biniş tahsilatı yok |
| FX3 | A1 akşam: 15 × 4 = 60, %20 → 48 | KISMİ | Aritmetik ✔; yön ayrımı yok |
| FX4 | A2 durak: 60, %20 → 48 | KISMİ | Net ✔; biniş anında yok |
| FX5 | A3: P=100, N=48, R=52 | KISMİ | N ✔; P ve R yok |
| FX6 | A4 matris: 50 → N=40; P=80 → R=40 | KISMİ | N ✔; P ve R yok |
| FX7 | A3 iade sıfır: P=N=100 → R=0 | KISMİ | N ✔; R=0 kapanış kavramı yok |
| FX8 | A3 tavan hatası: P=40 < N=48 → TAVAN YETERSİZ | YOK | Tavan kontrolü ve işaret yok |
| FX9 | Ücretsiz takip U0 = 0 | ÇELİŞKİ | Plaka + aktif sefer zinciri şart |

Özet (FX): 1 karşılanıyor · 6 kısmi · 1 yok · 1 çelişki.

### 3.5 Güzergâh bazlı ayarlar (01–25)

| No | Ayar | Durum | Gerekçe |
|---|---|---|---|
| 01 | Kurum / güzergâh kapsamı | YOK | Ayar hiyerarşisi yok; algoritma AccessPoint'te |
| 02 | Sefer yönü / sabah-akşam | YOK | Yön modellenmiyor |
| 03 | Tarife algoritması seçimi | KISMİ | AccessPoint.PricingMode (AD-12 tasarım farkı) |
| 04 | Ücretsiz geçiş / takip | ÇELİŞKİ | Plaka + aktif sefer zinciri |
| 05 | Sabit nokta (kampüs / fabrika) | YOK | Sabit uç kavramı yok |
| 06 | Mesafe kaynağı | KISMİ | Yalnız GPS haversine; durak tablosu seçimi yok |
| 07 | Mesafe tarifesi (km / dilim) | KISMİ | km oranı var; dilim yok |
| 08 | Durak fiyat tablosu | KISMİ | O-D matrisi; yalnız çıkışta |
| 09 | Noktalar arası matris | KISMİ | O-D ve yön ✔ |
| 10 | Maksimum ilk tahsilat | YOK | P yok |
| 11 | İndirim kuralı | KISMİ | Yalnız cüzdan kanalı; "bir kez" ✔ |
| 12 | Tahsilat anı (A1/A2 sabah-akşam, A3/A4 maksimum) | ÇELİŞKİ | Sabah/akşam ayrımı yok |
| 13 | İade hesabı ve kanalı | YOK | İade kanalı tanımı yok |
| 14 | İkinci okutma yoksa politika | ÇELİŞKİ | Politikasız otomatik borçlandırma |
| 15 | Yetersiz bakiye politikası | KISMİ | Ret ✔; borç politikası seçimi yok |
| 16 | Yanlış sefer / eşleşme | KISMİ | Plaka tabanlı; ilişki anahtarı sorunlu |
| 17 | Tekrar okutma / iade hatası | ÇELİŞKİ | İstemci anahtar göndermiyor; hatalar yutuluyor |
| 18 | Tarife sürümü / geçerlilik | ÇELİŞKİ | Fiyat time.Now() ile; yerinde PUT |
| 19 | Para birimi / yuvarlama | YOK | Yuvarlama yok; RouteFare para birimi yok |
| 20 | Denetim ve işlem kayıtları | KISMİ | Olay snapshot'ı eksik (kurum, kart, tarife sürümü, araç ID) |
| 21 | Sefer profili / araçtan ayrı | ÇELİŞKİ | Plaka zorunlu |
| 22 | Online borç politikası (kurum + hesap + güzergâh) | KISMİ | Taban kurum politikası ✔; güzergâh seçimi yok |
| 23 | Offline kabul / risk profili | YOK | Çevrimdışı yok |
| 24 | Offline A3/A4 kapanış | YOK | Çevrimdışı ve A3/A4 yok |
| 25 | Kişi hesabı / kimlikler | KISMİ | Kişi hesabı ✔; kart↔hesap geçerlilik tarihi yok |

Özet: 0 karşılanıyor · 11 kısmi · 8 yok · 6 çelişki.

## 4. Market akışları ve veri modeli

Akış durumları F01–F12 bölüm 2.2 altında özetlendi. Veri modeli varlıkları (13):

| Varlık | Durum |
|---|---|
| person / membership | KISMİ |
| account / wallet / pocket | KISMİ (pocket yok) |
| credential_binding | KISMİ |
| network / program / enrollment | YOK |
| merchant / branch / terminal | KISMİ |
| policy / agreement | YOK |
| payment_intent / attempt | KISMİ |
| payment_allocation / refund | ÇELİŞKİ (kümülatif sınır yok) |
| journal / journal_entry | KISMİ |
| outbox / inbox | YOK |
| settlement / payout | YOK |
| fiscal_event / invoice_job | YOK |
| report_projection / audit | KISMİ |

## 5. Karar listeleri → analiz kararları

Excel kararları sahipleriyle tanımlı ve tümü "Karar / değer bekliyor" veya benzeri durumda. Bu tablo, her Excel kararının analiz tarafındaki karşılığını gösterir. Öneri kararları AD-01…AD-15 olarak ana analizin §11 bölümünde numaralanmıştır.

### 5.1 Z kararları (P01–P11)

| ID | Konu | Analiz karşılığı |
|---|---|---|
| P01 | Sefer ve tarife doğrulama | AD-11, AD-12; K01, K02 |
| P02 | Kurum ve hesap kapsamı | AD-13; AK-01, AK-10 |
| P03 | Çevrimiçi borç limiti | AD-07; AK-06 |
| P04 | Çevrimdışı risk profili | AD-06; AK-09 |
| P05 | QR ve kart kimlik kuralları | Ana analiz §5.2 (Z03/Z04) |
| P06 | A3/A4 offline kapanışı | AD-02; AK-07 |
| P07 | Borç ve yükleme mutabakatı | AD-01, AD-07; AK-06, AK-10 |
| P08 | Fatura entegrasyonu | Faz 6 (invoice_jobs) |
| P09 | Rapor sözlüğü | AK-12; §5.6 |
| P10 | Performans ve saklama | Faz 1 ve Faz 7 (ölçüm) |
| P11 | Geçiş ve geri dönüş | Faz 9 (Z14 göçü) |

### 5.2 Market kararları (D01–D14)

| ID | Konu | Analiz karşılığı |
|---|---|---|
| D01 | Hukuki ve finansal model | AK-03 |
| D02 | İş ortağı ve yetkisi | AK-03 |
| D03 | Ağ işletmecisi | AK-04 |
| D04 | İlk kapsam (Türkiye / TL, kayıtlı kurum) | AK-04 |
| D05 | Bireysel fonlama | AK-02 |
| D06 | Cep ve kurum ayrımı | AK-02 |
| D07 | Ürün uygunluğu | M07 (Faz 8) |
| D08 | Kart ve QR doğrulaması | P05 |
| D09 | Komisyon ve hakediş | Faz 8 |
| D10 | Belge ve fatura profili | M19 |
| D11 | Offline ve borç | AD-06, AD-07 |
| D12 | Rapor ve saklama | AK-12 |
| D13 | Pilot başarı eşikleri | Faz 9 |
| D14 | Genişleme | Kapsam dışı (ilk faz) |

### 5.3 Devir kararları (DK01–DK12)

| ID | Konu | Analiz karşılığı |
|---|---|---|
| DK01 | Kapsam ve T | AK-13 |
| DK02 | Hukuki roller | AK-03 |
| DK03 | Pozitif bakiye | AK-02; DV07 |
| DK04 | Borç ve mahsup | AK-10; AD-07 |
| DK05 | Promosyon ve abonman | DV08 |
| DK06 | Açık yolculuk | AK-07; DV10 |
| DK07 | Offline cihazlar | AD-06 |
| DK08 | POS ve ödeme | DV13, DV14 (Faz 4) |
| DK09 | Belge düzeni | DV16 (Faz 4 ve 6) |
| DK10 | Veri ve yetki | DV17 (Faz 4) |
| DK11 | Mutabakat ve ihtilaf | DV18, DV19 |
| DK12 | Geri dönüş | DV04, DV20 |

## 6. Excel ile kod arasındaki tutarsızlıklar

1. **Ücretsiz takip (Excel s.1, 7, 9, 15):** "Mercedes/LCW ücretsiz personel takibi korunur; tahsilat yapmadan geçiş kaydı." Kodda ulaşım ücretsiz modu plaka ve aktif sefer ister → C-18.
2. **Algoritma seçimi düzeyi (Excel s.11):** "Tek ayar merkezi: güzergâh (kurum > güzergâh > yön)." Kodda algoritma AccessPoint.PricingMode'da → C-20, AD-12.
3. **"Dört yeni algoritma ayrı hizmet türü değil" (Excel s.9):** Kod her algoritmayı ayrı pricing mode olarak modelliyor (aynı kök sorun).
4. **Sabah biniş / akşam iniş (Excel T01, T02):** Kodda metered tek okutmayı yön bilmeden biniş sayıyor → C-19.
5. **"Süre, ceza, bakiye limiti ve iade kanalı varsayılmaz" (Excel s.11):** Kodda sefer bitince otomatik borçlandırma var → C-23.
6. **N > P tavan (Excel FX8):** Kodda P kavramı yok; uyarı üretilmiyor.
7. **Durak eşleşmezse km'ye sessiz düşüş:** Excel "tanımsız hücre sıfır değildir" diyor; sessiz algoritma değişimi ayrıca risk → T13 KISMİ.
8. **Excel içi:** Kimlik ve sayım tutarsızlığı bulunmadı.

## 7. Kapsam notu

- Bu dosya kod okumasına dayanır. Testler çalıştırılmadı; canlı veri ve üretim yapılandırması doğrulanmadı.
- Kod değişikliği için kullanıcı onayı beklenmektedir (ana analiz §8 Faz 0).
- Güvenlik notu: Sabit kimlik bilgisi değerleri bu dokümanda yazılmadı. Konum: `cmd/migrate_gursel` ve `internal/seeds/gursel_turizm_seed.go`.
