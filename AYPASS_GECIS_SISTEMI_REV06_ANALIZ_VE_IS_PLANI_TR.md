# AYPASS Geçiş Sistemi Rev.06 — Gereksinim Analizi, Kod Karşılaştırması ve İş Planı

**Sürüm:** 0.9 (taslak — bulgu bölümleri denetim sonuçlarıyla tamamlanıyor)
**Tarih:** 2026-10-28
**Kaynak doküman:** `AYSTEK_Gecis_Sistemi_Rev06.pdf` (Rev.06, 27.10.2026, 44 sayfa). Kullanıcı (Hüseyin Koç) sohbete ekledi; **PDF depoda yoktur**, bu dosyadaki envanter PDF metninden çıkarılmıştır.
**İlişkili dokümanlar:** `README_FOR_AI_ASSISTANTS.md` (bağlayıcı standartlar), `GURSEL_YENI_NESIL_DONUSUM_MASTER_TR.md` (Increment 1-169, son kayıt 2026-08-28; Karar #1-#31), `AYSHUTTLEGURSEL_CURRENT_STATE.md` (legacy envanteri).

---

## 0. Durum & Devam Notu (Sohbet Geçmişi Olmasa Bile Önce Bunu Okuyun)

- **Tarih:** 2026-10-28 · **Aşama:** Kod denetimi tamamlandı (5/5 denetçi raporu ve ana doğrulamalar). Excel işlendi; kimlik ve test eşlemesi `AYPASS_REV06_EXCEL_ESLEME_TR.md` dosyasında.
- **Girdi durumu:** Excel alındı (`AYSTEK_Gecis_Sistemi_Rev06-xlsx.pdf`, 41 sayfa). `.xlsx` yüklemesi desteklenmediği için PDF dökümü kullanıldı; sayımlar Excel'in kendi notlarıyla tutarlı. §9 artık Excel kimlikleriyle eşlenmiş durumda.
- **Kod durumu:** Hiçbir kaynak dosya değiştirilmedi; derleme veya test çalıştırılmadı; commit yapılmadı. `aystek-core` çalışma ağacında `.env` değişik; commit'e dahil edilmeyecek.
- **Doğrulama etiketleri:** **BAĞIMSIZ** = satır okunarak veya grep ile bu oturumda doğrulandı. **DENETÇİ** = denetçi raporunda kanıt satırıyla verildi, tekrar edilmedi (bazıları örneklemeyle kontrol edildi). **ÇIKARIM** = koddan akıl yürütüldü, çalıştırılmadı.
- **Sonraki adım:** (1) AD-01…AD-15 kararlarının onayı, (2) §8 Faz 0 için onay, (3) Excel karar sahiplerinin yanıtları (P, D, DK).

---

## 1. Yönetici Özeti

Rev.06, dört işletim modeli ve ulaşım tarifesi (A1-A4) için tasarım ve kabul gereksinimleridir; uygulanmış olduğunu iddia etmez. Mevcut AYPASS çekirdeği (kişi, kurum, cüzdan, kart/QR, erişim noktası, sefer kaydı, ledger) bu gereksinimlerin önemli bir kısmının temelini taşıyor. Ancak dört yapısal boşluk var:

1. **Çevrimdışı ve sıralı olay yok.** Cihaz kalıcı kuyruk tutmuyor, olay zamanı ve cihaz sırası yok, toplu eşitleme yok (Z05, Z06, Z08). Ağ kesintisinde geçiş "reddedildi" görünüyor ve olay kayboluyor.
2. **Para kesin hassasiyette değil, eşzamanlılık güvencesi zayıf.** Tüm tutarlar `float64`. Cüzdan bakiyesi okuma-hesaplama-yazma ile güncelleniyor (Z09, Z10).
3. **Sağlayıcı ataması ve devir modeli yok.** Geçerlilik tarihli "kurum × hizmet → sağlayıcı" kaydı bulunamadı. Devir iş akışı, bakiye devri ve atama çakışma kuralı kodda yok (M-02…M-04, DV-01…DV-08).
4. **Model 4 çekirdeği yok.** Cep (pocket), merchant/şube hiyerarşisi, hakediş (settlement), fonlama otoritesi ve outbox yok (MK-04…MK-10).

Buna ek olarak, Rev.06'dan bağımsız olarak **mevcut kodda üretim riski taşıyan hatalar** bulundu (P0):

- **P0-1 (BAĞIMSIZ):** `insufficient_balance` olaylarında `ChargedAmount` ücret tutarına set ediliyor; `RefundEvent` ise `Result`'a bakmadan iade ediyor. Tahsil edilmemiş ücret cüzdana geri yazılabilir.
- **P0-2 (BAĞIMSIZ):** GİB SOAP istemcisi taşıma hatasında sahte UUID döndürüp hata vermiyor; fatura "issued" olabilir. Gerçek GİB ortamının prod'da kurulu olup olmadığı doğrulanmadı.
- **P0-3 (BAĞIMSIZ):** `IssueEInvoice` adaptör hata döndüğünde `nil` sonucu okuyup panic yapabilir.
- **P0-4 (DENETÇİ, repository tarafı BAĞIMSIZ):** Cüzdan bakiyesi koşulsuz `$set` ile yazılıyor; eşzamanlı borçlandırmalarda taban aşılabilir.
- **P0-5 ve P0-6 (DENETÇİ):** Callback replay'inde kredi kaybı riski; ödeme iadesinde toplam iade ≤ tahsilat kontrolü yok.

Denetçi raporları tamamlandı. Ek **P1** bulgular: (a) `TripAssignment` güncellemesi onay, devir ve sürücü alanlarını yazmıyor; manuel istisna akışı kalıcı değil. (b) `RecordEvent` pasif AccessPoint'i kabul ediyor. (c) Ücretsiz ulaşım takibi plaka ve aktif sefer zincirine bağlı; Mercedes/LCW ücretsiz takibi bu zincire takılıyor. (d) Plaka değişimi faturalamayı sessizce atlatıyor. (e) Stok plaka tarihçesinde aktör istek gövdesinden alınıyor.
**Güvenlik (P0-S):** `cmd/migrate_gursel` içinde sabit bağlantı dizgisi var, derlenmiş binary git'te izleniyor, seed'de düz metin parola ve token alanları bulunuyor. Değerler bu dokümanda yazılmadı; rotasyon ve temizlik gerekli.

**Öneri:** P0 maddeleri Rev.06 işinden önce ayrı bir artım olarak kapatılsın (§8 Faz 0). Bunlar para ve yasal fatura etkisi taşıdığı için kod değişikliği onayınıza bağlı; onay gelmeden kodda değişiklik yapılmadı.

---

## 2. Kaynak, Yöntem ve Sınırlar

- **Kaynak:** Rev.06 PDF'i (44 sf.): s.1-6 iş modelleri, s.7-12 ulaşım tarifeleri (A1-A4), s.13-20 zorunlu kurallar Z01-Z14, s.21-36 Model 4 (Global Market), s.37-44 Model 1 sağlayıcı devri.
- **Doküman statüsü (PDF'in kendi ifadesi):** "Tasarım ve kabul gereksinimi; uygulanmış/canlı iddiası değildir. Testler tanımlandı; yazılımda çalıştırılmadı." Bu analiz de aynı dürüstlük kuralını izler: **bir gereksinimin kodda karşılandığı yalnızca kod okunarak doğrulandıysa yazılır.**
- **Excel sayımları (doğrulandı):** Z: 14 gereksinim (Z01–Z14), 24 kabul testi (K01–K24), 11 açık karar (P01–P11). Market: 24 tasarım kuralı (M01–M24), 30 kabul senaryosu (MK01–MK30), 14 karar (D01–D14). Devir: 20 gereksinim (DV01–DV20), 10 iş adımı, 12 karar (DK01–DK12), 26 test (DT01–DT26), 16 veri alanı grubu. Ulaşım: 6 algoritma satırı, 25 güzergâh ayarı (01–25), 15 kontrol senaryosu (T01–T15), 9 formül örneği (FX1–FX9).
- **Kod tabanı referansı:** `/workspace/go/aystek-core` HEAD `f48748f` (2026-10-27); çalışma ağacında yalnızca `.env` değişik (git-tracked sızıntı — **asla stage/commit edilmez**, bkz. hafıza notu). Master plan artım kayıtları 2026-08-28'de (Increment 169) bitiyor; sonrasındaki AYPASS değişiklikleri master plana işlenmemiş olabilir, bu yüzden bulgular **koda** dayanır, plan metnine değil.
- **Kod denetimi:** 2026-10-28'de beş alanda salt-okunur denetim yapıldı (para/defter/cüzdan; cihaz/offline/Flutter; fatura/ödeme; ulaşım tarife-sefer-araç; tenant/topoloji). Doğrulama etiketleri §0'da tanımlıdır.
- **Kapsam dışı:** Hukuki/mevzuat değerlendirmesi (PDF'in kendisi canlı öncesi hukuk/finans onayı şart koşar; §10'da risk olarak işaretlenir), PDF'te yer almayan yeni iş kuralları.

---

## 3. Gereksinim Envanteri (ID'li)

> **Not:** Z01–Z14 Excel'deki resmi kimliklerdir. Aşağıdaki T, MK ve DV kimlikleri iç gruplamadır (PDF bölümlerinden türetildi). Resmi kimlikler: M01–M24, F01–F12, MK01–MK30, DV01–DV20, DT01–DT26, DK01–DK12, FX1–FX9. Eşleme: `AYPASS_REV06_EXCEL_ESLEME_TR.md`.


ID şeması: **M** = iş modeli, **T** = ulaşım tarifesi, **Z** = zorunlu kural (PDF'te Z01-Z14), **MK** = Model 4 (Global Market), **DV** = Model 1 sağlayıcı devri. "Kabul kontrolü" sütunu PDF'teki ifadelerin özetidir. Durum sütunları §5'te doldurulur.

### 3.1 İş modelleri (s.1-6)

| ID | Gereksinim (öz) | Sf. |
|---|---|---|
| M-01 | Dört model tek platformda; ortak kişi hesabı ve izlenebilir işlemler; hizmet, finans ve kurum yetkileri ayrı. Model 4, M1/M2/M3'ü iptal etmez, isteğe bağlı ortak kabul ağıdır; fabrika QR/kart nokta takibi ayrı yetenek olarak korunur. | 1 |
| M-02 | **Model 1:** kurum ve hizmet (ulaşım/yemek-kafe/kütüphane/fotokopi) başına ayrı hizmet firması seçilebilir; yüklenen tutarları ilgili hizmet firması faturalar. | 2 |
| M-03 | **Model 2:** merkezi operatör (Gürsel) hizmeti kendisi verir veya taşerona devreder; kullanıcılar kendi kurum grubunda doğrulanır; hizmeti veren fatura keser; taşeron değişince önceki taşeronun kayıtları kalır. | 3 |
| M-04 | **Model 3:** tek kurum (Sabancı), sektör bazında ayrı firmalar, ihale devri kritik; kurum içi kullanıcılar tek grup. | 4 |
| M-05 | Kurum × hizmet kapsam matrisi (VAR/YOK/TAKİP); Mercedes/LCW ulaşım = ücretsiz TAKİP. Ulaşım sabit/indirimli + 4 algoritma; yemek ürün başına; kütüphane ücret modeli **açık**; fotokopi kullanım bedeli. | 5 |
| M-06 | Fabrika/tesis QR-kart geçiş takibi: tanımlı noktalar; olay alanları kurum-tesis-nokta-zaman (+ sonuç/kimlik cihaz verisine göre); yalnız veri aktarabilen noktalar; her kurum yalnız kendi kayıtlarını görür. | 6 |

### 3.2 Ulaşım tarifesi (s.7-12)

> **Not:** Excel'de resmi kimlikler algoritma satırları (F0, Ücretsiz, A1–A4), güzergâh ayarları (01–25) ve kontrol senaryolarıdır (T01–T15). Bu tablodaki T-xx iç gruplamadır.


| ID | Gereksinim (öz) | Sf. |
|---|---|---|
| T-F0 | Mevcut sabit/indirimli ücret korunur; ücretsiz takip (tahsilatsız geçiş/personel kaydı) korunur. | 7 |
| T-A1 | **Sabit uç + mesafe:** kampüs/fabrika bilinir; sabah **binişte**, akşam **inişte** okutulur; brüt = mesafe tarifesi(d) (km başı veya dilim); net = varsa indirim sonrası; doğrudan tahsilat. Sabah/akşam mesafeleri ayrı belirlenir. | 8 |
| T-A2 | **Durak bazlı tarife:** anahtar = güzergâh + yön + durak; mesafeden ücret üretilmez; varsa indirim. A4 ile birbirinin yerine kullanılmaz. | 8 |
| T-A3 | **Maksimum + mesafe + iade:** binişte maksimum P tahsil; inişte 2. okutma ile mesafeden net N; iade R = P − N; indirim yalnız nihai ücrete **bir kez**; P olası tüm net ücretleri karşılayan tavan olmalı. | 9 |
| T-A4 | **Maksimum + matris + iade:** binişte P tahsil; iniş okutmasında matris[biniş, iniş] tutarı; yön asimetrik olabilir; **tanımsız hücre sıfır ücret değildir**; indirim sonrası fark iade. | 10 |
| T-CFG | Güzergâh bazlı ayar: kapsam+algoritma (F0/A1-A4), sabit uç/duraklar, mesafe kaynağı+tarife, durak fiyat tablosu, matris, tahsilat zamanı, maksimum+iade, indirim/ücretsiz, eksik-hatalı işlem politikası, sürüm/yuvarlama/yetki/eşleştirme kayıtları. Süre, ceza, bakiye limiti, iade kanalı varsayılmaz. | 11 |
| T-EXC | Tahsilat − iade = nihai ücret. İstisnalar: tekrarlı okutma (tek etki), ikinci okutma yok (tavan kalır/inceleme/düzeltme — güzergâhta seçilir), yetersiz bakiye (reddet/borçlandır; başarısız ödeme gizlenmez), iade başarısız (bekleyen durum; ikinci para gönderilmez), tarife/konum eksik (varsayımsal fiyat yok; inceleme), net > maksimum (işaretle; negatif iade/izinsiz ek tahsilat yok). Zorunlu kayıt: kurum, kişi hesabı, olay/yolculuk kimliği, kart/QR, sefer, araç, tarife sürümü, mali işlem ilişkileri. | 12 |

### 3.3 Zorunlu kurallar Z01-Z14 (s.13-20)

| ID | Gereksinim (öz) | Kabul kontrolü (PDF) | Sf. |
|---|---|---|---|
| Z01 | Ücret bağlamı (kurum, güzergâh, yön, sefer, tarife sürümü) **araç/plaka atamasından ayrılır**; doğrulanmış sefer oturumu, yetkili seçim veya önceden tanımlı profil ile tarife bulunur; araç atama eksikliği tek başına geçişi engellemez; tarife belirsizse tutar uydurulmaz, **olay kaydedilir**, onaylı istisna akışı uygulanır. | Araç ataması olmadan doğrulanmış seferle doğru tarife bulunur; belirsizde varsayımsal tahsilat yok. | 13 |
| Z02 | Araç sabit `vehicle_id`; plaka değişebilir gösterim; sefer-araç ilişkisi başlangıç/bitiş zamanlı; düzeltme eski ilişkiyi silmeden, yapan kişi+gerekçe+zaman ile; yanlış plaka düzeltmesi mali hareket üretmez, bakiye devretmez, geçmiş ücreti değiştirmez. | Araç değişiminden önce/sonra hareketler ayrı izlenir; tahsilat sayısı ve bakiye değişmez. | 13 |
| Z03 | Bakiye **kişi hesabına** aittir (kart/QR yalnız erişim kimliği); yeni kart aynı hesaba bağlanınca bakiye/borç/geçmiş korunur, iade/devir hareketi oluşmaz; eski kart iptal; kimlik-hesap eşleşmesinin **geçerlilik tarihi** saklanır; gecikmiş çevrimdışı olay **okutma anındaki** hesaba bağlanır. | Kart yenileme hesabı değiştirmez; iptal öncesi geçerli gecikmiş olay yeni kart sahibine/başka hesaba yazılmaz. | 14 |
| Z04 | Karta bakiye/borç/hareket **yazılmaz**; kart/QR hesabı belirler; kimlik+cihaz doğrulama kuralları; QR geçerlilik süresi/yeniden kullanım/iptal; çevrimdışı önbellek yaşı ve cihaz yetkileri; yerel kayıtlar asgari veri, şifreli, değişikliğe karşı korumalı; aynı QR/kartın farklı çevrimdışı cihazlarda kullanım riski yönetilir. | Kart belleğine yazma olmadan geçiş işler; iptal ve tekrar kullanım çevrimiçi/çevrimdışı ayrı test edilir. | 14 |
| Z05 | Onaylı çevrimdışı profille geçiş olayı **başarı bildirilmeden önce** cihazın kalıcı kuyruğuna yazılır (olay kimliği, cihaz, kurum, hesap eşleşmesi, güzergâh-sefer, araç, durak, okutma zamanı, tarife/indirim sürümü, biliniyorsa tutar); uygulama kapanması/yeniden başlama olay kaybettirmez; **geçiş kabulü ile mali tahsilat onayı ayrı durumlar**. | Kesinti sonrası kuyruk geri gelir; merkezi tahsilat onayı olmadan tahsil edildi sayılmaz. | 15 |
| Z06 | Her olay için değişmeyen benzersiz anahtar, tek mali etki; sunucu kalıcı kabul etmeden cihaz kuyruğu temizlenmez; olay kabulü ile mali sonuç ayrı teyit; **olay zamanı ve sunucu kabul zamanı ayrı**; geç/sırasız olaylar ilişkilendirilir; sonradan değişen tarife eski yolculuğa uygulanmaz; eşitleme hataları izlenebilir. | Aynı olay çok kez gelse tek borç/alacak; kayıp yanıt sonrası yeniden deneme sonucu bozmaz. | 15 |
| Z07 | Çevrimiçi alt bakiye sınırı kurum/hesap politikası, **işlem sonrası bakiye** üzerinden kontrol (−15 TL evrensel değil); çevrimdışı kabul (azami süre/tutar/adet + sonradan borçlandırma) **ayrı politika**; güzergâh seçimi kurum sınırlarını aşmaz; verilmiş hizmetin alacağı kaybolmaz (onaylı borç kaydı/yetkili inceleme), tahsil edilmiş nakit gibi gösterilmez. | 0 TL hesapta 120 TL çevrimiçi ücret −15 TL sınırıyla onaylanmaz; çevrimdışı alacak ayrı politika/durumla izlenir. | 16 |
| Z08 | Tüm algoritmalar için olay+sefer+tarife sürümü çevrimdışı korunur; A3/A4 biniş-iniş aynı yolculuğa bağlanır, iniş önce gelirse eşleşme beklenir; P merkezi mali kayıtta doğrulanmadan P−N iadesi gerçekleşmiş yazılamaz; eşitlemede biniş tahsilatı/izinli borç ve iniş düzeltmesi tekilleştirilir; sorunlu olay saklanır, **mali istisna** durumuna alınır. Çok cihazda sınırın kesin korunacağı varsayılamaz (onaylı risk + mutabakat). | Çevrimdışı biniş-iniş ters sırada gelse tek yolculuk, tek mali etki, doğru nihai ücret. | 16 |
| Z09 | Yükleme/tahsilat/iade/izinli borç/düzeltme hareketleri aynı `account_id` ve işlem ilişkileriyle **ortak denetlenebilir hareket modeli**; kesinleşmiş hareket silinmez/değişmez, bağlı ters/düzeltme kaydı; bakiye güncellemesi ile mali kayıt **atomik** veya eşdeğer güvenilirlikte; tek mantıksal kaynak ≠ tek fiziksel koleksiyon. | Bakiye açılış + onaylı hareketlerden yeniden hesaplanır; rapor kopyası ayrı mali gerçek üretmez. | 17 |
| Z10 | Online yükleme kesin ödeme onayı + tekil dış işlem referansı; tekrar bildirim ikinci yükleme üretmez; eşzamanlı yükleme/tahsilat/iade bakiye kaybettirmez; negatif bakiyeye yükleme bakiyeyi artırır, aynı borç ikinci kez kesilmez; ayrı alacak izleniyorsa mahsup kaydı açık; **tutarlar para birimiyle kesin hassasiyette**; ödeme/hesap/rapor mutabakatı. | Aynı ödeme bildirimi iki kez gelse tek yükleme; −120 TL hesaba 200 TL yükleme 80 TL verir. | 17 |
| Z11 | Faturalanacak her yükleme özgün mali hareket kimliğiyle; mali kayıt + fatura işi tutarlı oluşur (iş kaybolmaz); ayrı kuyruk olabilir ama ikinci mali kaynak olamaz; veri çekmek ≠ fatura oluşması; tekil istek anahtarı, deneme geçmişi, dış referans, fatura kimliği; belirsiz yanıt önce **sorgulanır**, kör tekrar yok. Durumlar: bekliyor/işleniyor (kontrollü sahiplenme, devralma) → gönderildi/sonuç belirsiz → faturalandı \| hata/yeniden deneme. Yükleme faturası ≠ yolculuk tüketim kaydı. | Fatura dış sistemde oluşup yanıt kaybolsa sorgulamayla eşleşir; aynı yükleme yeniden faturalanmaz. | 18 |
| Z12 | Raporlar doğrulanmış hareket kaynağından ve ortak tanımlarla; kasa (gerçek yükleme/tahsilat, iade, borçlandırma, bekleyen) ayrı; kullanım ortalamasının pay/payda/dönem/iptal sefer davranışı belgelenir; araç raporunda personel türü dağılımı/toplam geçiş/tekil kişi ayrı; özetten özgün olaya inilir; rapor güncellik zamanı ve gecikmiş çevrimdışı olaylar görünür. | Toplamlar kaynak hareketlerle uyuşur; sorgu performansı gerçekçi hacimle ölçülür. | 19 |
| Z13 | Kurum, kişi, hesap, kart/QR, araç, güzergâh, sefer, durak, hizmet sağlayıcı ilişkileri **sabit ID**; ad/plaka ilişki anahtarı olamaz; olay anındaki personel türü/ad-plaka görüntüsü/tarife/işletici tarihsel kopya veya sürümlü ilişki; değişiklik izi (kim/ne zaman/ne) ve kurum erişim sınırları. | Ad/plaka değişimi ilişkileri bozmaz; eski tarihli rapor olay anındaki sınıflandırmayı korur. | 19 |
| Z14 | Göç/canlıya çıkış: kart-kişi-hesap eşleşmesi, mükerrer kişi, sahipsiz kart doğrulanır; bakiye/borç ile hareket toplamları, açık yolculuklar, bekleyen iadeler, faturalar, cihaz kuyrukları mutabık; kesim zamanı, geri dönüş planı, sorumlular; eski hareketlerin yeniden tahsilat/fatura üretmesi engellenir; hacim/indeks/arşiv kararları ölçüme dayanır. | Pilot ve mutabakat onayı olmadan canlıya geçiş tamamlandı sayılmaz; test tanımı ≠ test geçti. | 20 |

Canlıya geçiş öncesi dört doğrulama grubu (s.20): (1) operasyon + kart değişimi, (2) kesinti + tekrar + sıra, (3) bakiye + yükleme + fatura, (4) rapor + mutabakat.

### 3.4 Model 4 — Global Market (s.21-36)

> **Not:** MK-xx iç gruplamadır. Resmi kimlikler M01–M24 (tasarım kuralları), F01–F12 (akışlar), MK01–MK30 (kabul), D01–D14 (kararlar).


| ID | Gereksinim (öz) | Sf. |
|---|---|---|
| MK-01 | Yapı: kurum/sponsor (hak+bütçe, fon finansal ortağa) → kişi hesabı + cepler → üye market/şube/kasa; **yetkili banka/e-para ortağı** fon, ödeme onayı, tahsilat/iade, işyeri ödemesi ve mutabakatı yürütür. Ağ işletmecisi ticari rolü fon tutma/ödeme lisansı vermez. İlk kapsam Türkiye/TL. | 21 |
| MK-02 | Ar-Ge dersleri: yalnız okuyucu değil **market operasyon paneli**; satış-hakediş-banka teyidi ayrı; amaç/sponsor bazlı ayrı cepler. Hiçbir marka entegrasyonu iddia edilmez. | 22 |
| MK-03 | Finansal model: A) finansal ortak (**önerilen başlangıç**), B) sınırlı program (ayrı hukuki inceleme), C) kendi yapısı (ilk faz dışı). Canlı öncesi zorunlu onay: ödeme emrini kimin aldığı/onayladığı, ortağın güncel izinleri, kredi yasağı, sözleşme sorumlulukları hukuk/finansça doğrulanır. | 23 |
| MK-04 | Kişi > doğrulanmış kurum üyeliği > finansal hesap > ihraççı/para birimi/program/sponsor. **Cepler** (ulaşım, yemek/gıda, market, promosyon); kart/QR aynı kişiyi bulur, harcanacak cep o program ve market için uygun olmalı; **bakiyeler otomatik birleşmez**, kısıtlı hak serbest cebe/nakde dönüşmez; kurum ayrılışında hak/kişisel fon politikaları ayrı. Kart belleğine yazma yok. | 24 |
| MK-05 | Üye market hiyerarşisi: market tüzel kişiliği > şube > kasa/terminal > yetkili kullanıcı; alanlar: kabul ağı, uygunluk (alt seviye üst kuralı genişletemez), ticari profil (komisyon/vade/iade bedeli/rezerv/sözleşme sürümü — sayısal değerler onayla), finansal kimlik (IBAN değişimi çift kontrol), teknik profil (pasif terminal işlem yapamaz). Öncelik: mevzuat/ihraççı > program > kurum hakkı > işyeri sözleşmesi > şube/kasa; işlem anı sürümü saklanır. **Market satışında güzergâh/araç/sefer ataması aranmaz.** | 25 |
| MK-06 | Kasa akışı: sepet → kimlik (QR/kart) → onay (market+tutar) → kesin tahsilat → fiş/hakediş. **Kasada başarı = finansal ortağın doğrulanmış kesin tahsilat sonucu.** Kart seri numarası yetmez (ortak onaylı kart/cihaz veya ek müşteri onayı); QR işlemle birlikte doğrulanır (süreli, tek kullanımlık; tutar+işyeri+para birimine bağlı); banka PIN'i genel ekranda toplanmaz. Ulaşımdaki maksimum-tahsilat/iade algoritması kullanılmaz. | 26 |
| MK-07 | Ödeme durumları: CREATED → AUTHORIZED → CAPTURED → REFUNDED; PENDING_EXTERNAL (aynı intent+sağlayıcı referansıyla sorgula; zaman aşımı ≠ ret); CANCELLED/RELEASED; IDEMPOTENCY (aynı anahtar+farklı tutar/sepet = çatışma, yeni tahsilat yok). Tek aşamalı satış destekleniyorsa doğrudan CAPTURED; desteklenmeyen rezervasyon varsayılmaz. | 27 |
| MK-08 | Tek mantıksal doğruluk kaynağı: finansal ortak (harcama yetkisi), AYSTEK mali iz (çift taraflı kayıt + orijinal kaynak dağılımı), rapor modeli (yeniden üretilir, as-of + mutabakat farkı). Katmanlar: kimlik/kapsam; ödeme-iade (intent/capture/allocation/refund/provider_ref); defter-olay (journal/entry/**outbox/inbox**; yerel mali kayıt + outbox atomik); hakediş-belge. Dış API DB transaction'ında değil: sağlayıcı sorgusu + tekilleştirilmiş yerel kayıt + dönemsel mutabakat. | 28 |
| MK-09 | İade: capture kimliği zorunlu, toplam iade ≤ ilk tahsilat, iade kullanılan cep/fon kaynağına döner; komisyon iadesi ayrı sözleşme kuralı. Çok kaynak: ilk faz aynı ihraççıda onaylı cep bölüşümü; parça başarısızsa işlem tamam sayılmaz. Ödeme sonrası iade: sonraki dönem mahsup/rezerv/alacak; eski ödenmiş hakediş dosyası değişmez. Kısıtlı kaynak iadeyle serbest cebe/nakde dönmez. | 29 |
| MK-10 | Hakediş: kesin satış − iade − komisyon − diğer kesinti − rezerv + fonlanmış destek − önceki ödemeler; hesaplandı → onaylandı → gönderildi → banka teyidi → mutabık; "gönderildi ≠ ödendi"; yanıt kaybında aynı banka referansı sorgulanır; rapor kırılımı kurum/ihraççı/para birimi/işyeri/şube/dönem > işlem > banka referansı; farklı fonlar izinsiz mahsup edilmez. | 30 |
| MK-11 | Belge profili: dört ayrı ekonomik olay (fonlama, market satışı, komisyon/hizmet, iade/iptal); Rev.04 yükleme faturalama kuralı tüm market programlarına otomatik taşınmaz; alındı/gönderildi/oluştu ayrı; belge/vergi düzeni mali müşavir ve hukukça yazılı onaylanır. | 31 |
| MK-12 | Market kesintisi: ulaşım politikası kopyalanmaz; **çevrimiçi onay + yeterli uygun bakiye** şart; kesintide sepet saklanır/alternatif kanal; okutma başarılı tahsilat gibi gösterilmez; **otomatik eksi bakiye kapalı**; istisnai offline kabul ayrı faz (risk üstlenicisi, hukuki dayanak, limitler, mutabakat). | 32 |
| MK-13 | Güvenlik/gizlilik/rapor: yetkili kasa + cihaz kaydı + iptal/askı kontrolü; tutar/market bağlı onay, doğrulanmış callback, tek kullanımlık QR; IBAN/sözleşme değişiminde çift kontrol; market kendi satışını, kurum yetkili programı görür, işverenin kişisel sepete varsayılan erişimi yok; finansal rapor sözlüğü (ciro AYSTEK geliri değil); indeks+benzersiz işlem anahtarları; ortak-defter-banka mutabakatı. | 33 |
| MK-14 | Uygulama sırası: Faz 0 rol/sözleşme → Faz 1 online pilot (TR/TL) → kabul/mutabakat → Faz 2 kontrollü büyüme. Kontrol grupları: kimlik/kapsam, tahsilat/tekrar, iade/hakediş, belge/kesinti/gizlilik. | 34 |
| MK-15 | İlk faz dışı: sınırsız offline, ulaşım borcunu markete taşıma, her markanın kartı, onaysız fonlar arası transfer, yurt dışı/çok para birimi. "Global" = kurumlar arası sözleşmeli kabul ağı. | 35 |

### 3.5 Model 1 sağlayıcı devri (s.37-44)

> **Not:** DV-xx iç gruplamadır. Resmi kimlikler DV01–DV20 (gereksinim), DK01–DK12 (karar), DT01–DT26 (test).


| ID | Gereksinim (öz) | Sf. |
|---|---|---|
| DV-01 | Korunan yapı: platform + kurum + kişi hesapları + kart/QR ilişkileri. Devir kapsamı: kurum + hizmet + sözleşme (gerekirse hat/sefer grubu); diğer hizmetler/kurumlar etkilenmez. T öncesi başlayan seferler ve eski mali işlemler kendi bağlamında kapanır; T ve sonrası yeni atamaya gider. Beşinci model değil, Model 1'e eklenen yetenek. | 37 |
| DV-02 | İş akışı: 6 kapı (kapsam/envanter; mali plan/protokol; yeni firma hazırlığı; prova ve kesim onayı; T anında hizmet geçişi; mutabakat/kapanış). Durumlar: Taslak → Hazır/onaylı → Planlandı → Operasyon devredildi → Mali kapanış bekliyor → Kapandı. İhtilafta otomatik kapanış yok; kurum+hukuk+finans istisna kararı kaydedilir. Hizmete başlama ile mali kapanış ayrı izlenir. | 38 |
| DV-03 | Sürümlü atama veri modeli: `institution_id+service_id+scope_id`; `provider_id+service_assignment_id` (yeni firma yeni kimlik, eski silinmez); `valid_from/valid_to` **[başlangıç, bitiş)** (eski bitiş = T = yeni başlangıç); `trip_id/journey_id` + sürümler (başlamış sefer sabit); `seller_id + liable_entity_id` (işletici/hukuki satıcı/fon sorumlusu aynı varsayılmaz); `payment_profile_id + version`; `event_id + devir_id + zamanlar`. Aynı kapsamda aktif atamalar çakışmaz. | 39 |
| DV-04 | Bakiye devri ≠ hizmet devri: A) aynı cüzdan devam eder, B) fon eski firmada kalır (iade/kapanış), C) onaylı mali devir (tek transfer ID, eşli azaltım/kabul, çift harcanabilir bakiye/mükerrer fatura yok). Ayrı kalemler: nakit, eski borç/alacak, promosyon, kurum katkısı, abonman, bloke/bekleyen iade. Borç otomatik taşınmaz. | 40 |
| DV-05 | Kesim anı senaryoları: eski sefer/yeni dönem inişi (eski sözleşme; A3/A4 aynı P ve tarife, ikinci maksimum tahsilat yok); geç offline aktarım (**olay-anı** atamasıyla eski dönem, geliş zamanı firmayı seçtirmez); süresi dolmuş profil (olay saklanır, istisna/yenileme); geç ödeme onayı (özgün fon/hesap/ödeme profiline tek kayıt); saat/bağlam güvenilmez (tahminle tahsilat yok, inceleme). | 41 |
| DV-06 | POS/iade/fatura/yetki: yeni satış → yeni onaylı POS (fiziksel cihazın aynı kalması onay yerine geçmez; eski şifre/anahtar kopyalanmaz); eski iade → özgün işlem zinciri; belge → işlemin asıl tarafı; erişim → asgari kapsam ve rol. AYSTEK teknik geçişi yönetir, fon toplama/banka adına onay yetkisi yaratmaz. | 42 |
| DV-07 | Mutabakat/kapanış/geri dönüş: kesim envanteri (hesap+para birimi+fon sorumlusu bazında), ilk tutanak korunur, geç olaylar ayrı düzeltme sürümü; çözülmeyen kalemin sahibi/durumu/takip tarihi; geri dönüş: yeni mali işlem yoksa plan iptali, varsa yeni atama sürümü + onaylı düzeltme/ters kayıt + mutabakat. | 43 |
| DV-08 | Kabul grupları: kapsam+tarihçe; hesap+mali devir; sefer+offline; POS+iade+fatura; yetki+kapanış. Gerçek T, yeni firma, bakiye planı ve banka onayı atanmadı. | 44 |

---

## 4. Tasarım İlkeleri ve Hedef Mimari

### 4.1 Gereksinimlerin ortak köküne dair tespit

Z01-Z14, MK ve DV'nin büyük kısmı **beş mühendislik ilkesine** indirgenir. Bunlar ayrı ayrı "özellik" olarak değil, platformun ortak sözleşmesi olarak ele alınmalıdır; aksi halde her modül aynı sorunu kendi kopyasıyla çözer (README: mükerrer yapı yasağı).

| # | İlke | Karşıladığı gereksinimler |
|---|---|---|
| P1 | **Para = değer nesnesi.** Tutar + para birimi tek tipte, kesin hassasiyette (tam sayı küçük birim), yuvarlama kuralı sürümlü tarifede. `float64` ile tutar taşınmaz. | Z10, T-CFG (yuvarlama), MK-08, DV-07 |
| P2 | **Çift-zamanlı olay.** Her olay `occurred_at` (geçerlilik zamanı) ve `received_at` (kabul zamanı) taşır. Atama, kart-hesap eşleşmesi, tarife sürümü, indirim kuralı, ödeme profili **`occurred_at`** ile çözülür; düzeltme yeni sürümdür. | Z02, Z03, Z06, Z08, Z13, DV-03, DV-05 |
| P3 | **Sürümlü referans veri.** `[valid_from, valid_to)` yarı-açık aralık; aynı kapsamda çakışma yok (depo seviyesinde doğrulanır); silme yok, kapatma var. | T-CFG, Z03, DV-03, MK-05 |
| P4 | **Değişmez defter + izdüşüm.** Kesinleşmiş hareket değişmez; ters kayıt üretilir; bakiye izdüşümü journal ile **aynı transaction'da** güncellenir; yan etkiler (event, fatura işi, bildirim) outbox ile. | Z09, Z10, Z11, MK-08 |
| P5 | **Komut idempotensisi.** Cihaz/dış sistem olay kimliği benzersiz anahtardır; tekrar = aynı sonucun yeniden sunumu; aynı anahtar + farklı içerik = çatışma. | Z05, Z06, Z10, Z11, MK-07 |

Tamamlayıcı ilkeler: **P6** strateji nesneleri (tarife algoritması, fonlama otoritesi, ödeme sağlayıcı portu, politika nesneleri); **P7** açık oturum soyutlaması (iki okutmalı akışlar: A3/A4 ulaşım, olası kütüphane ödünç/iade); **P8** durum makineleri model katmanında `CanTransition` ile (README modeli); **P9** olay anı anlık görüntüsü (sağlayıcı/satıcı/fon sorumlusu/ödeme profili/tarife sürümü olaya yazılır — Z13); **P10** mevcut modül genişletilir, yeni modül yalnız kavram yoksa açılır.

### 4.2 Hedef alan modeli ve mevcut modüllerle eşleme

Her satır için README'nin "önce tara, genişlet, mükerrer yazma" kuralı uygulandı: mevcut bir modül varsa genişletilir; yoksa yeni modül önerilir.

| Gereksinim kümesi | Mevcut karşılık | Hedef eylem | Not |
|---|---|---|---|
| Para değer nesnesi (P1) | `currencies.Decimals` ve `RoundAmount` (yalnız `orders`); tutarlar `float64` (ledger, wallets, access_points, payments, route_fares) | **Değiştir:** minor-unit `int64` + para birimi; mevcut `currencies` bilgisini kullan; migration | Tüm para modüllerini etkiler (AD-01) |
| Çift zamanlı olay (P2) | `access_points.AccessEvent.created_at` (sunucu zamanı) | **Genişlet:** `occurred_at` (cihaz), `received_at` (sunucu), `device_seq`, cihaz kimliği | Z05, Z06 |
| Sürümlü referans (P3) | `route_fares`, `trip_contracts`, `products/pricing`, `taxes` (`valid_from/to` kalıbı mevcut) | **Kullan:** kalıp hazır; sağlayıcı ataması için yeni agrega | README: önce tara |
| Sağlayıcı ataması (M-02…M-04, DV-03) | Yok (BULUNAMADI). En yakın: `trip_contracts` (rota/araç kapsamı yok), `Vehicle.SupplierID` tek işaretçi | **Yeni modül** `service_assignments` (AD-03): `institution × service × scope → provider`, `[valid_from, valid_to)`, çakışma yasağı | Model 1/2/3'ün ortak çekirdeği |
| Devir iş akışı (DV-02) | `trip_assignments` durumları; `approvals` (bağlı değil) | **Yeni** durum makinesi (`service_transfers`); onay kapıları `approvals` ile | 6 kapı, ihtilaf kuralı |
| Değişmez defter (P4) | `ledger` journal/entries; silme yok; ters kayıt yok | **Genişlet:** `reversal_of` alanı, void/reversal servisi; cüzdan cache'ine doğrudan yazım kapatılır | Z09 |
| Outbox (P4) | Yok. `platform/events` süreç-içi ve kalıcı değil; `platform/queue` ve `platform/rabbitmq` ayrı. Modül-yerel outbox koleksiyonları 3 yerde kopya (sms, pms/webhooks, digital_editions) | **Yeni ortak paket** `platform/outbox`; mevcut üç kopya bu pakete taşınır | README: mükerrer yasağı |
| Komut idempotensisi (P5) | `platform/idempotency` (Redis, 24 saat, get-then-set; para modülleri kullanmıyor); ledger partial unique index | **Genişlet:** DB-backed, payload parmak izi, aynı anahtar + farklı içerik = çatışma hatası | Z06, Z10, MK-07 |
| Strateji portları (P6) | `payments/intents.GatewayAdapter` (port mevcut); tarife algoritmaları access_points içinde gömülü | **Soyutla:** `FareAlgorithm` (F0, A1–A4, metered) ve `FundingAuthority` (iç defter / ortak destekli) | AD-02, AK-03 |
| Kart ve QR (Z03, Z04) | `wallet_instruments` (`IssueCard`, `GenerateQr`, `ValidateCard`) | **Genişlet:** QR tek kullanım (atomik tüketim), revoke, nonce; çevrimdışı için imzalı token seçeneği | AD-06 |
| Olay zamanı kuralları (P2, P9) | Aktif sefer ve SUL zaman kontrolü sunucu saatine göre | **Değiştir:** çözüm `occurred_at` ile yapılsın; geç olay kabul kuralı tanımlansın | C-02 |
| Model 4 — cep (pocket) | `wallets` (person + currency tekil) | **Genişlet:** `wallet_pockets` (program/sponsor/ihraççı bazlı) veya tekil anahtarın genişletilmesi; migration | MK-04 |
| Merchant > şube > kasa (MK-05) | `institutions` (Type + hiyerarşi), `access_points` (kasa/terminal) | **Genişlet:** `merchant` tipi, yetkili kullanıcı, ticari profil sürümü | — |
| Hakediş ve settlement (MK-10) | Yok | **Yeni modül** `settlements` | — |
| Maker-checker | `approvals` silindi (AD-04, 2026-10-28); yeni tasarım gerekli | **Yeniden tasarla:** ObjectID, tenant+site kapsamı, istek sahibi ≠ onaylayan (servis katmanında, gerçek Mongo testi), rol kontrolü, kalıcı karar kaydı | — |
| Fatura işi (Z11) | `billing/invoices` (senkron), `invoiceable_records` (read-model) | **Değiştir:** `invoice_jobs` durum makinesi (kuyruk, lease, deneme, dış referans); `journal_id ↔ invoice_id` bağı | P0-2, P0-3 |
| Raporlama (Z12) | `access_points` özet uçları (`KindSummary`, `RouteSummary`, `PersonSummary`) | **Genişlet:** kasa kalemleri ayrımı, reddedilen ve iade filtresi, as-of zamanı, drill-down | C-14 |

---

## 5. Kod Karşılaştırması — Bulgular (Kanıtlı)

### 5.1 Para, defter ve cüzdan

| Gereksinim | Kod durumu | Kanıt | Doğrulama | İlişki |
|---|---|---|---|---|
| Z10 — kesin hassasiyet | Tüm tutarlar `float64`: ledger `Entry.Amount`, `Wallet.CachedBalance`, `AccessEvent.ChargedAmount`, `PaymentIntent`, `RouteFare`. Yuvarlama yalnız ürün satışında (`money()`, sabit 2 hane). Ulaşım ücreti (`baseFee*(1-pct/100)`, `km*rate`) yuvarlanmadan borçlanıyor. `currencies.RoundAmount` yalnız `orders` tarafından kullanılıyor. | `ledger/model/journal.go:51`; `wallets/model/wallet.go:68`; `access_points/service/service.go:801-803, 1224-1227` | DENETÇİ | YOK |
| Z09 — atomiklik | Journal ve entry Mongo transaction içinde yazılıyor (`InsertJournal`). Cüzdan cache'i ayrı bir adım. | `ledger/repository/mongo.go:155-226` | DENETÇİ | KISMİ |
| Z09 — silinmez / ters kayıt | Journal ve entry için update/delete yok. Ancak ters kayıt (void/reversal) yolu da yok: `JournalVoid` sabiti hiç kullanılmıyor. `POST /ledger/journals` cüzdan hesabına cache'i güncellemeden yazabiliyor. | `ledger/model/journal.go:13`; `ledger/routes.go:16` | BAĞIMSIZ (sabit, route); DENETÇİ (kullanım yokluğu) | YOK |
| Z10 — eşzamanlılık | `applyBalanceChange`: okunan bakiye + delta hesaplanıp koşulsuz `$set` ile yazılıyor. Taban kontrolü de okunan değerle yapılıyor. | `wallets/repository/mongo.go:111-127` (BAĞIMSIZ); `wallets/service/service.go:732-775` (DENETÇİ) | BAĞIMSIZ + DENETÇİ | YOK (**P0-4**) |
| Z10 — tekrar bildirim | Ledger'da `(tenant_id, idempotency_key)` üzerinde partial unique index var. Aynı anahtar farklı tutarla gelirse ilk journal sessizce döner (payload parmak izi yok). | `ledger/repository/indexes.go:18-23` | BAĞIMSIZ | KISMİ |
| Z09 — yeniden hesaplanabilirlik | `Reconcile` toplar, düzeltmez, alarm üretmez. Pending havale entry'leri de toplamda görünür; sahte drift oluşur. Zamanlayıcı yok (yalnız `POST /wallets/reconcile-all`). | `wallets/service/service.go:913-934`; `cmd/api/main.go` zamanlayıcılarında yok | DENETÇİ | KISMİ |
| Z07 — online taban | Taban sırası: `FloorExempt` → `BalanceFloorOverride` → kurum `MinWalletBalance` → 0. Kontrol işlem sonrası bakiyeyle yapılıyor. Çok kurumlu kişide "ilk" affiliation (`created_at DESC [0]`) kullanılıyor; erişim noktasının kurumu kullanılmıyor. | `wallets/service/service.go:97-124`; `persons/repository/mongo.go:312-331` | DENETÇİ | KISMİ |
| Z07 — alacak kaydı | `insufficient_balance` yalnız `AccessEvent` olarak kalıyor; borç veya alacak kaydı oluşmuyor. Otomatik kapanışta `Debit` hatası yutulup `ChargedAmount` yazılıyor. | `access_points/service/service.go:1207-1208`; `fleet_tally/rider_trips/service/service.go:215, 223` | DENETÇİ | YOK |
| Z09/Z10 — para birimi tutarlılığı | Journal içinde tek para birimi kontrolü var. Cüzdan↔journal, intent↔cüzdan ve erişim noktası↔cüzdan karşılaştırması yok. | `ledger/service/service.go:158-171` | DENETÇİ | KISMİ |
| MK-04 — cepler | Cüzdan, kişi + para birimi başına tekil (`unique(tenant, site, person, currency)`). Program, sponsor veya ihraççı alanı yok. | `wallets/repository/indexes.go:15-24` | BAĞIMSIZ | YOK |
| Model 4 — hold / kullanılabilir bakiye | `available_balance`, `held_amount`, `reserved`, `pending_balance` yok. En yakın mekanizmalar: `JournalPendingReconciliation` (yalnız havale) ve intents `requires_capture` (kart provizyonu). | GREP (DENETÇİ) | DENETÇİ | YOK |

### 5.2 Cihaz geçidi, çevrimdışı dayanıklılık, kart ve QR

| Gereksinim | Kod durumu | Kanıt | Doğrulama | İlişki |
|---|---|---|---|---|
| Z05 — kalıcı kuyruk | Flutter'da yerel kuyruk yok. `pubspec.yaml` sqflite/drift/hive/isar/shared_preferences/connectivity içermiyor. Tek kalıcı depo `flutter_secure_storage` (token ve secret). Başarı bildirimi sunucu yanıtından sonra gösteriliyor. | `aypass_erisim_app/pubspec.yaml:44-66` | DENETÇİ | YOK |
| Z05 — olay zamanı | Olay isteğinde `occurred_at`, `device_time` veya `scanned_at` yok. `AccessEvent` tek zaman alanı taşıyor: `created_at = time.Now()` (sunucu). | `access_points/dto/dto.go:50-60`; `access_points/model/access_event.go:156` | DENETÇİ | YOK |
| Z06 — tekrar gönderim | `idempotency_key` isteğe bağlı. Aynı anahtar aynı olayı döndürüyor (201). Anahtar kapsamı tenant + site; içerik karşılaştırması yok. Cash settlement ve walk-on'da erken arama yok (tekrar 400, ÇIKARIM). | `access_points/service/service.go:933-938, 1619-1649`; `access_points/repository/indexes.go:59-69` | DENETÇİ | KISMİ |
| Z06 — sıra ve toplu eşitleme | `device_seq` yok, batch ucu yok. | GREP (DENETÇİ) | DENETÇİ | YOK |
| Z05 — geçiş kabulü ≠ mali tahsilat | Tek modal ("Geçiş Onaylandı" + tahsilat). Ayrı durum yok. | `aypass_erisim_app/lib/screens/validation_result_screen.dart:30-33` | DENETÇİ | YOK |
| Z05/Z06 — ağ kesintisi | Hata sınıfı ayrımı yok; genel catch → "Geçiş Reddedildi". Sunucuda işlenip yanıtı kaybolan istek için yeniden okutma ikinci borçlandırma riski doğurur. | `aypass_erisim_app/lib/services/validation_service.dart:84-96` | DENETÇİ + ÇIKARIM | YOK (P1) |
| Reddedilen okutma kaydı | `AccessEventResultDenied` hiç üretilmiyor. Reddedilen okutma `AccessEvent` bırakmıyor, yalnız audit log. | `access_points/model/access_event.go:59` | DENETÇİ | YOK |
| Z08 — geç olay | Taşıma olayı aktif sefer şartı taşıyor. Sefer bittikten sonra gelen olay reddediliyor; test bu davranışı sabitliyor. | `access_points/service/service_test.go:3143` | DENETÇİ | ÇELİŞKİ |
| Z04 — QR | 32 bayt opak token; yalnız SHA-256 saklanıyor; TTL 90 sn. `ConsumedAt` yalnız denetim amaçlı, tek kullanım zorunlu değil. Revoke yok. İmza yok: her doğrulama sunucu DB'si gerektiriyor. Tutar veya işyeri bağı yok. | `wallet_instruments/model/qr.go:14-29`; `wallet_instruments/service/service.go:247-297` | DENETÇİ | KISMİ |
| Z03/Z04 — kart | `ValidateCard` yalnız güncel durumu okuyor (lost/blocked/replaced → 400). Olay anındaki geçerlilik değerlendirmesi yok. Kart↔cüzdan bağı değişmiyor (uyumlu), ancak geçerlilik tarihi yok. | `wallet_instruments/service/service.go:119-153`; `access_blocks/service/service.go:95-101` | DENETÇİ | KISMİ |
| MK-06 — UID ile borçlandırma | Kart UID'si tek başına cüzdan borçlandırabiliyor. Ulaşım için kabul edilebilir mi kararı gerekiyor; Model 4 için MK-06 bunu açıkça yasaklıyor. | `wallet_instruments/service/service.go` | DENETÇİ | ÇELİŞKİ (Model 4) |
| Z04 — cihaz yetkisi | OAuth2 client_credentials (RS256 JWT, 15 dk), her istekte DB'den durum kontrolü, rate limit (PIN'li cihaz 10 rps, 5000/gün), rota yetkilendirmesi. Kara liste, nonce, cihaz başına limit ve çoklu cihaz riski yönetimi yok. | `auth/middleware/api_client_gateway.go:183-349` | DENETÇİ | KISMİ |
| Z04 — çoklu çevrimdışı cihaz | Çevrimdışı doğrulama olmadığı için kısıt yalnız çevrimiçi anlık bakiye (taban). Anti-passback yok. | — | ÇIKARIM | YOK |
| Z13 — kimlik kaynağı | `AccessEvent` kullanılan kart veya QR'ı (UID ya da token hash) saklamıyor. `person_id` istemciden geliyor ve cüzdanla tutarlılığı doğrulanmıyor (denetçi okuması). **Çelişki:** master plan Increment 148, cüzdan sahibinin backend'de çözüldüğünü yazıyor. Denetçi bunu kodda teyit edemedi; doğrulanacak. | `access_points/service/service.go:940-956`; `access_points/dto/dto.go:50-60` | DENETÇİ | ÇELİŞKİ (doğrulanacak) |
| MK-12 — başarı gösterimi | Başarı yalnız sunucu yanıtından sonra gösteriliyor; sepet ve ürün listesi yalnız bellekte. | `aypass_erisim_app/lib/services/canteen_cart_service.dart:9-45` | DENETÇİ | KISMİ |
| Model 4 — ödeme onayı | Onay isteği tipi (tutar + işyeri + süre, onayla/reddet) yok. FCM yalnız başlık ve gövde gösteriyor. PIN alanı yok (MK-06 "PIN toplanmaz" uyumlu). | `aypass_kullanici_app/lib/services/fcm_service.dart:160-175` | DENETÇİ | YOK |
| Cihaz izin listeleri | `invitation.go` (PIN ile yeni istemci) `…/me/continuing-riders` ve `…/cancel` uçlarını içermiyor; `admin.go` AP-bağlı istemcide rota listesini yeniden yazıyor ve `/device/qr|card/validate` içermiyor. 403 olası. | `auth/handler/invitation.go:192-206`; `auth/handler/admin.go:33-57` | DENETÇİ + ÇIKARIM | YOK (P1) |
| Cihaz PIN'i | 6 hex (24 bit), 30 dk geçerli, herkese açık redeem. `ip_whitelist` saklanıyor; uygulandığı yer bulunamadı. Kaba kuvvet ve IP kısıtı açısından risk. | `auth/handler/invitation.go:54-60, 101` | DENETÇİ | RİSK (P1) |

### 5.3 Fatura, ödeme, iade ve hakediş

| Gereksinim | Kod durumu | Kanıt | Doğrulama | İlişki |
|---|---|---|---|---|
| Z11 — iş kuyruğu / durum makinesi | Yok. Fatura gönderimi senkron, admin isteğiyle tetikleniyor. `Invoice.Status` ∈ {draft, issued, cancelled, failed}; `EInvoice.Status` serbest string. Claim/lease, deneme geçmişi, sonraki deneme zamanı, tekil istek anahtarı yok. | `billing/invoices/model/invoice.go:10-17, 38-46` | DENETÇİ | YOK |
| Z11 — sahte başarı | GİB SOAP istemcisi `callSOAP` hatasında `generateMockUUID` döndürüyor ve `nil` hata veriyor. `submitToGIB` → `SendInvoice` çağrısı bu sonucu başarı yoluna sokuyor. | `billing/invoices/adapters/gib_soap.go:76-86` (BAĞIMSIZ); `gib_einvoice.go:399` (BAĞIMSIZ) | BAĞIMSIZ | ÇELİŞKİ (**P0-2**) |
| Z11 — nil panic | `IssueEInvoice`: `result, err := adapter.Issue(...)` sonrası `err` kontrolünden önce `result.Status` okunuyor. GİB `Issue`, config veya fatura doğrulama hatasında `nil, err` dönüyor. | `billing/invoices/service/service.go:265-271` (BAĞIMSIZ); `gib_einvoice.go:131-137` (BAĞIMSIZ) | BAĞIMSIZ | YOK (**P0-3**) |
| Z11 — çift fatura | `IssueEInvoice` mevcut `ExternalID` veya durumu kontrol etmeden yeniden gönderiyor. GİB UUID'si `IssuedAt` boşken `time.Now()`'a bağlı; tekrarda değişir. | `billing/invoices/service/service.go:256-294`; `gib_einvoice.go:153` | DENETÇİ | YOK |
| Z11 — belirsiz yanıt | Her adaptör hatası `failed` yapılıyor; "gönderildi / sonuç belirsiz" ayrımı yok. `QueryEInvoiceStatus` elle çağrılabiliyor, ancak `QueryStatus` her zaman `ACCEPTED` dönüyor ve UBL XML "placeholder". | `billing/invoices/service/service.go:256-294`; `gib_soap.go:92` | DENETÇİ | KISMİ |
| Z11 — faturalandı ↔ yükleme bağı | Journal ile fatura arasında bağ yok. `InvoiceMark` faturaya referans vermiyor; read-model yalnız elle işaretleme yapıyor. | `billing/invoiceable_records/model/model.go:15-51` | DENETÇİ | YOK |
| invoiceable_records — kapsam | `ListInvoiceable` her kaynak için yalnız en yeni 1000 journal'ı çekip işaretlileri eliyor; eski işaretsiz kayıtlar görünmeden kalabilir. Journal'da `site_id` yok; işaret ve ayar site-scoped (ÇIKARIM: çok siteli tenant'ta tutarsız). | `billing/invoiceable_records/service/service.go:37-71`; `ledger/model/journal.go:21-36` | DENETÇİ | KISMİ |
| MK-07 — durum eşlemesi | `intents` durumları: requires_payment_method/confirmation/action/capture, succeeded, cancelled, failed. CREATED≈requires_confirmation, AUTHORIZED≈requires_capture (yalnız Stripe'ta `manual_capture`; nestpay/qnb/isbank/garanti'de Capture/Cancel/Refund `ErrOperationNotSupported`), CAPTURED≈succeeded. REFUNDED, RELEASED ve PENDING_EXTERNAL yok. Sağlayıcı sorgulama metodu yok. | `payments/intents/model/intent.go:12-18` (BAĞIMSIZ); `payments/intents/service/adapters.go:44-51` (DENETÇİ) | KISMİ | — |
| MK-07 — idempotency | `capture/refund/cancel/callback` için `(tenant, intent, op, key)` unique. `Create`'te anahtar yok. Aynı anahtar farklı tutarla gelirse eski sonuç sessizce dönüyor. Başarısız op kaydı, sonraki başarılı op'u duplicate-key ile yutabilir. | `payments/intents/repository/indexes.go:36-39`; `payments/intents/service/service.go:477-531` | DENETÇİ | ÇELİŞKİ (**P0-6**) |
| MK-09 — iade üst sınırı | Yalnız `status == succeeded` ve `amount > 0` kontrolü var. Toplam iade ≤ capture kontrolü ve `refunded_amount` alanı yok. | `payments/intents/service/service.go:477-531` | BAĞIMSIZ (kontrol yokluğu) + DENETÇİ | YOK (**P0-6**) |
| MK-09 — iade kimliği | Refund kimliği yalnız `ProviderPayload["refund"]` içinde tutuluyor; son iade öncekinin üstüne yazıyor. Capture kimliği yok. | `payments/intents/service/service.go:510` | BAĞIMSIZ | YOK |
| Callback — replay ve kurtarma | `recordOperation(callback)`, `succeededCb`'den önce yazılıyor. `succeededCb` hata verirse kayıt kalıyor; yeniden gelen bildirim replay sayılıp atlanıyor → tahsil edilmiş, kredisi işlenmemiş intent. Webhook retry servisi ve job main'de bağlı değil. | `payments/intents/service/service.go:745-788` (DENETÇİ); `cmd/api/main.go` (GREP: bağlı değil) | DENETÇİ | ÇELİŞKİ (**P0-5**) |
| Callback — tutar doğrulaması | `CallbackResult.Amount` ve `Currency` alanları var, ancak `HandleCallback` bunları intent ile karşılaştırmıyor. | GREP (DENETÇİ) | DENETÇİ | YOK |
| Callback — imza | Güçlü imza: Stripe HMAC + 5 dk, PayPal verify API, iyzico ve PayTR hash. QNB imzasız; iyzico'da başlık yoksa doğrulama atlanıyor. Replay nonce'u imzadan önce yazılıyor. | `payments/intents/service/adapters.go` | DENETÇİ | KISMİ |
| Mock adaptör | `APP_ENV = production|prod` ise mock kaydı yapılmıyor (BAĞIMSIZ). Diğer ortamlarda imzasız callback kabul eden mock aktif. | `payments/intents/service/adapters.go:102-105` | BAĞIMSIZ | RİSK (P1) |
| Yükleme (wallet_topup) | Tutar, sunucudaki `intent.Amount` ile krediliyor; anahtar `payment-intent:<id>`; çift kredi ledger anahtarıyla engelli. Intent ↔ cüzdan para birimi kontrolü yok. | `wallets/service/payment_intent_bridge.go:26-37` | DENETÇİ | KISMİ |
| Cihaz yüklemesi | `/device/wallet-topups` anahtarı `topup:<code>:<extTxID>` (zorunlu). Aynı txid farklı cüzdan veya tutarla gelirse ilk journal'a düşer; yanıt yeni tutarı yankılar. | `wallets/topup_sources/service/service.go:97-132` | DENETÇİ | KISMİ |
| Havale | Beyan → pending (entry'ler hemen ledger'da) → admin onayı. Reddetme veya iptal yok. Referans boşsa rastgele üretiliyor; tekrar beyan yeni kayıt açıyor. Anahtar cüzdana bağlı; aynı banka referansı iki cüzdana beyan edilebilir. | `wallets/service/service.go:393-400` | DENETÇİ | KISMİ |
| Negatif bakiye ve mahsup | Özel kod yok: `CachedBalance + delta` → −120 + 200 = 80 (doğru sonuç, ama testsiz). Wallets ↔ receivables mahsup mekanizması yok. | GREP (DENETÇİ) | DENETÇİ | KARŞILANIYOR (testsiz) / YOK (mahsup) |
| Hakediş, settlement, komisyon, IBAN çift kontrol | Bulunamadı. Kodda "settlement" nakit tahsilat anlamında kullanılıyor. IBAN yalnız storefront `BankAccount` ve `wallets.Refund` serbest metin hesap no. | GREP (DENETÇİ) | DENETÇİ | YOK |
| Onay (maker-checker) | `approvals` yazılmış ama `cmd/api/main.go`'ya bağlı değil (grep: 0 sonuç). `Approve`, onaylayanın rolünü ve istek sahibi ≠ onaylayan kuralını denetlemiyor. | `cmd/api/main.go` (BAĞIMSIZ); `approvals/service/approvals.go:71-111` (DENETÇİ) | BAĞIMSIZ | YOK |
| Yazılı ama bağlı olmayan kod | payments webhook retry, `internal/queue/topology.go` (fatura kuyrukları). `approvals` ve `integrations/governance/approval.go` silindi (AD-04). README kuralı: kullanılmayan kod ya bağlanır ya silinir; kalan iki kalem için karar gerekiyor (AD-05). | — | DENETÇİ | — |
| ERP çakışması | `erp_integrations` (Logo/Mikro/Nebim InvoiceDTO push) ile `billing/invoices` ERP adaptörleri (Logo/Netsis) mükerrer. | — | DENETÇİ | YOK (mükerrer) |

### 5.4 Ulaşım tarifesi, sefer, araç kimliği ve yolculuk durumu (ulaşım denetçisi — DENETÇİ)

- **A3/A4 primitifleri:** Biniş `Debit(P)` olarak yapılabilir. İniş `CreditFromAccount(P−N)` ikinci çağrısı HTTP'de açık değil; `RefundEvent` tek çağrı, yalnız tam `ChargedAmount` iadesi yapıyor. Dolayısıyla kısmi iade bugün ifade edilemiyor (DENETÇİ).
- **Metered giriş:** `chargeMetered` girişte hold veya yetkilendirme yapmıyor; çıkışta N'yi borçlandırıyor (kod yorumu BAĞIMSIZ okundu: `access_points/service/service.go:1449-1457`).
- **Tarife sürümü:** Olay, tarife sürüm kimliğini değil yalnız `BaseFee`, `DiscountRuleID`, `DiscountPercentage` alanlarını yazıyor (DENETÇİ). Sonradan değişen tarifenin eski yolculuğa uygulanmadığı garanti değil.
- **Aktif sefer plakayla çözülüyor:** Master plan Karar #30 ve Increment 149, aktif seferi plaka ile `in_progress` kaydı üzerinden buluyor (doküman okundu). Z13, "ad/plaka ilişki anahtarı olamaz" diyor → **ÇELİŞKİ C-08**; kod tarafı ulaşım denetçisiyle doğrulanacak.

### 5.5 Tenant, sözleşme ve tedarikçi modelleri

Sözleşme ve tedarikçi modülleri (DENETÇİ, fatura/ödeme denetçisi raporu):

- **Sağlayıcı ataması bulunamadı:** Geçerlilik tarihli "sağlayıcı-kapsam ataması" yok (`ServiceAssignment`, `provider_assignment`, `ownership_history` aranıp bulunamadı). En yakın kavram `trip_contracts` (tedarikçi × müşteri firma, sefer oranı, `EffectiveFrom/To`, çakışma kontrolü, `FindActiveForPair(at)`); ancak rota ve araç kapsamı içermiyor ve T-anı snapshot'ı yok.
- **Fatura kesen taraf modelde yok:** Hiçbir model "hizmeti veren" alanı taşımıyor; GİB VKN ortam değişkeninden tek değer olarak geliyor.
- **Model 2 kendi aracı:** `Vehicle.SupplierID` boşsa `trip_billing` satır üretmiyor. Yani Gürsel'in kendi araçları için "hizmeti veren fatura keser" temsil edilmiyor.
- **`trip_billing`:** Sefer bitişinde best-effort çalışıyor; eksik yapılandırmada sessiz atlıyor, yeniden deneme yok. Fatura kesmiyor. Durum akışı: pending → submitted → approved | disputed.
- **`supplier_invoices`:** Tedarikçinin dış faturası; `SupplierAmount` toplamıyla %1 toleranslı mutabakat. `invoice_number` unique değil.
- **`receivables`:** Kişi/kurum borcu (open/paid/cancelled), `PaymentIntentID`, `SourceRef` (SAP).
- **Tek işaretçi:** `Vehicle.SupplierID` tek sağlayıcı işaretçisi; geçmişi tutulmuyor (`PlateHistory` yalnız plaka).
- **İzolasyon katmanları (DENETÇİ):** `tenantfilter` tenant ve site zorlar. `institutionfilter` opt-in üçüncü katman; tek seçili kurum, hiyerarşi yok; kısıtsız kullanıcı herhangi bir `X-Institution-ID` seçebiliyor. `site_bridge` salt-okuma site grant'ı; fiilen yalnız `stock_vehicles` kullanıyor. `institution_directory` dış rehberden aynı siteye Person ve Affiliation yazar. `api_sharing` tek tenant'a bağlı dış OAuth2 istemcisi.
- **Veri yaşam alanı:** Person, Wallet, Institution, AccessPoint, AccessEvent tenant ve site ile filtreleniyor; ledger, payments, subscriptions, contracts, products/pricing, billing ve approvals yalnız tenant ile filtreleniyor. Ortak kişi hesabı yalnız aynı site içinde; site ve tenant arası yok. Gürsel tenant'ında beş site var (aractakip, admin-panel, aypass, veli, comfortpay): kurum, kişi, cüzdan, kart ve AP aypass'ta; veli ayrı host; SAP borç senkronu comfortpay'de. **ÇIKARIM:** aynı kişi üç sitede ayrı Person olabilir.
- **Topoloji (AD-13):** A (tek tenant ve site, `institutions.Type=supplier`) mevcut kodla en tutarlı; eksikleri sağlayıcı ataması, rakip izolasyonu ve fatura kesen taraf. B (sağlayıcı başına tenant) ve C (kurum başına tenant) ortak kişi hesabıyla **çelişir**: paylaşılan kurum kaydı tenant_id ile kopyalanır.
- **Affiliation ve taban:** Person ve kurum başına tek kayıt; durum {active, withdrawn, pending_approval}; geçmiş ve sürüm yok; birleştirme (merge) yok. `wallets.resolveInstitution` `affiliations[0]` kullanıyor (en yeni, durum fark etmez); çok kurumlu kişide taban politikası belirsiz.
- **Sözleşme ve sağlayıcı:** `Vehicle.SupplierID` tek işaretçi, geçmişsiz. `trip_contracts` (supplier × customer, [from, to), çakışma kontrolü; rota ve araç kapsamı yok). `scheduled_routes` ve indirim kurallarında geçişte iki kayıt aynı anda aktif olabilir. Hizmet–kurum–sağlayıcı için tek agrega yok; zincir RouteDefinition → ScheduledRoute → TripAssignment → Vehicle.SupplierID.
- **Kurum ACL (Karar #25):** ACL'li yönetici başka kurumun AP ve olay verisini göremiyor (access_points list/get/events, access_blocks, discounts, receivables, supplier_invoices, vehicles). Ancak persons, wallets, wallet_instruments, institutions, ledger, trip_contracts ve RecordEvent/RefundEvent yazma uçlarında kontrol yok. "Kurum yetkili programı görür" ve alan maskeleme (SaleItems) ifade edilemiyor.
- **Onay ve denetim izi:** `approvals` `main.go`'ya bağlı değil; chain doldurulmuyor, audit yazılmıyor. `telematics/commands` iki aşamalı onay içeriyor (domain'e bağlı). `activity_logs` domain modüllerinde çağrılmıyor; fiili iz `CreatedBy` ve `UpdatedBy` (son yazar).
- **Göç (Z14):** `cmd/migrate_gursel` yalnız stok ve araç takip verisini taşıyor (ürünler, satılanlar, abonelikler, iadeler → stock_*). Kart, kişi, bakiye, hareket, borç, açık yolculuk ve bekleyen iade göçü yok; mutabakat ve kesim aracı yok.
- **Güvenlik (AD-15, P0):** `cmd/migrate_gursel` içinde bağlantı dizgisi (parola dahil) sabit; derlenmiş binary git'te izleniyor; `gursel_turizm_seed.go` içinde parola ve token alanları düz metin. Değerler bu dokümanda yazılmadı; rotasyon gerekli.
- **Veri tutarlılığı:** `AccessEvent.LedgerEntryID` aslında journal ID tutuyor (ad uyuşmazlığı). `LedgerAccount.Code` için unique index yok. `invoiceable_records` tenant geneli journal listeler (limit 1000, sayfalama yok); `InvoiceMark` site bazlı; siteler arası mükerrer fatura riski (ÇIKARIM).
- **Modül kaydı (mükerrerlik):** Yüksek risk: trip_contracts, contracts, stock/vehicles (SupplierID ve PlateHistory). Orta: scheduled_routes, access_blocks, subscriptions, institutions (taraf), billing DistributionRule.

| Gereksinim | Kod durumu | Kanıt | Doğrulama | İlişki |
|---|---|---|---|---|
| Z01 — araçsız yol | Transport olayı plaka + aktif TripAssignment + SUL ister; eksikse `return nil, err`: olay, ledger, RiderTrip yazılmaz; HTTP 400 | `access_points/service/service.go:932-1012` | DENETÇİ | ÇELİŞKİ (K01, K02) |
| Z01/Z05 — istisna | `AccessEventResultDenied` hiç üretilmiyor; pending ve review yok | `access_points/model/access_event.go:59` | DENETÇİ | ÇELİŞKİ (K02) |
| Ücretsiz takip (ulaşım) | `free` modu da plaka + aktif sefer + SUL ister; Mercedes/LCW personel takibi bu zincire bağlı | `service.go:430-432, 991` | DENETÇİ | ÇELİŞKİ (T15, FX9, C-18) |
| F0 tarifesi | CardFee/CashFee; net = brüt × (1 − %); indirim tek kural, yalnız cüzdan kanalında | `service.go:1080-1105, 1151, 1163` | DENETÇİ | KARŞILANIYOR (FX1) |
| A1 — sabit uç + mesafe | Metered: giriş-çıkış GPS haversine × km oranı; sabit uç, mesafe dilimi, sabah/akşam yön yok; ücret çıkışta | `service.go:1519-1535`; `rider_trips/service/service.go:176-227` | DENETÇİ + ÇIKARIM | ÇELİŞKİ (T01, T02) |
| A2 — durak | RouteStop fiyat taşımaz; O-D matrisi yalnız çıkışta kullanılır | `fleet_tally/route_stop_fares/model/route_stop_fare.go:21-46` | DENETÇİ | KISMİ (T03) |
| A3 — maksimum + iade | P, binişte tahsil ve R = P − N yok; hold yok | `service.go:1449-1457` | BAĞIMSIZ (yorum) + DENETÇİ | KISMİ (N var; P ve R yok) |
| A4 — matris | Yönlü O-D; tanımsız hücre km oranına düşer (sıfır değil ✔); durak toleransı dışında sessiz km düşüşü | `service.go:1524-1536, 1344-1353` | DENETÇİ | KISMİ (T06, T13) |
| Algoritma düzeyi | Excel: güzergâh düzeyinde (kurum > güzergâh > yön); kod: `AccessPoint.PricingMode` | `service.go:990-1011`; `access_point.go:114` | DENETÇİ | ÇELİŞKİ (tasarım; C-20, AD-12) |
| Tavan, yuvarlama, para birimi | Tavan yok; yuvarlama yalnız yemekhane `money()`; RouteFare para birimi yok | `service.go:801-803` | DENETÇİ | YOK (T14, AK-05) |
| Tarife sürümü | `route_fares` dönemleri ve süreklilik var; kayıt PUT ile yerinde düzenlenir; fiyat `time.Now()` ile; olayda tarife sürüm kimliği yok | `route_fares/service/service.go:76-150`; `service.go:1087, 1117, 1525` | DENETÇİ | ÇELİŞKİ (K11, C-19) |
| Olay bağlamı | Kurum snapshot'ı, kart/QR kimliği, tarife sürümü, araç ID yok | `access_event.go:89-152` | DENETÇİ | KISMİ (Z13) |
| Plaka anahtarı | AP.VehiclePlate ↔ TripAssignment; alarms, verifications, movement_reports, trip_billing; `approach_notifier` `FindByPlate` yalnız `current_plate` | `service.go:433-440`; `trip_assignment.go:71` | DENETÇİ | ÇELİŞKİ (Z13, K23, C-21) |
| Stok araç kimliği | `_id` sabit; `PlateHistory` var; gerekçe yok; aktör istek gövdesinden | `stock/vehicles/dto/dto.go:41`; `handler/admin.go:265` | DENETÇİ | KISMİ (Z02); RİSK (C-22) |
| Durum makinesi | open → closed / auto_closed / cancelled; OpenTrip ve CloseTrip hataları yutulur; admin kapama ucu yok | `service.go:1488, 1549` | DENETÇİ | KISMİ |
| Sırasız olay | Olay zamanı yok (`CreatedAt = time.Now()`); "sonraki dokunuş = iniş" | `service.go:1830` | DENETÇİ | ÇELİŞKİ (K15, C-19) |
| İkinci okutma yok | Sefer bitince planlı mesafe × oran doğrudan Debit; indirim yok; anahtar yok; hata yutulur, ChargedAmount yazılır | `rider_trips/service/service.go:212-217` | DENETÇİ | ÇELİŞKİ (T10, T11; C-23, C-24) |
| İade | Tam ChargedAmount; pencere `RefundWindowMinutes`; Result kontrolü yok → `insufficient_balance` olayı iade edilebilir | `service.go:328-391` | BAĞIMSIZ | ÇELİŞKİ (P0-1; C-26) |
| Tekrar okutma | Anahtar isteğe bağlı; AYPASS uygulamaları göndermiyor; QR 90 sn çok kullanımlı; yarışta ham E11000 | `service.go:933-938`; `repository/indexes.go:59-69`; `wallet_instruments/model/qr.go:24-26` | DENETÇİ | ÇELİŞKİ (T09) |
| Pasif terminal | RecordEvent `AccessPoint.Status` denetlemez | `service.go` (grep) | BAĞIMSIZ | ÇELİŞKİ (MK05; C-25) |
### 5.6 Raporlama (Z12) — DENETÇİ

- `reports` genel rapor tasarlayıcı. Alan raporları `access_points` içinde: `KindSummary`, `RouteSummary` (tekil kişi `$addToSet`), `PersonSummary`.
- Eksikler: result ve refunded süzülmüyor; reddedilen ve iade edilen olaylar ciroya giriyor (`repository/mongo.go:466-625`). Yükleme, iade, borçlandırma ve bekleyen işlemler tek kasa raporunda ayrışmıyor. Güzergâh ortalaması ve pay/payda tanımı yok. `TripStatus` yalnız planned, in_progress, completed; iptal sefer kavramı yok. Araç (plaka) bazlı özet ve personel türü (`AffiliationRole`) dağılımı yok; olayda yalnız `RiderCategory` var. Özetten olaya drill-down bağı yok. `as_of` ve `generated_at` yok.
- Durum: **ÇELİŞKİ** (Z12; C-14).

---

## 6. Mevcut Karar ve Uygulamalarla Çatışmalar

| # | Gereksinim | Mevcut davranış | Tür | Bölüm |
|---|---|---|---|---|
| C-01 | Z07: verilmiş hizmetin alacağı kaybolmaz; onaylı borç kaydı veya inceleme oluşur | `insufficient_balance` yalnız `AccessEvent`; borç veya alacak kaydı yok | ÇELİŞKİ | §5.1 |
| C-02 | Z05/Z08: geç ve sırasız olay kaydedilir, eşleşme beklenir | Aktif sefer yoksa olay reddediliyor; test bunu sabitliyor | ÇELİŞKİ | §5.2 |
| C-03 | Z10: para kesin hassasiyette tutulur | `float64`; ulaşım ücreti yuvarlanmadan borçlanıyor | ÇELİŞKİ | §5.1 |
| C-04 | Z11: belirsiz yanıt sorgulanır, kör tekrar yapılmaz | GİB sahte UUID ile başarı; her hata `failed`; `IssueEInvoice` çift gönderim | ÇELİŞKİ | §5.3 |
| C-05 | MK-07: aynı anahtar + farklı içerik = çatışma; yeni tahsilat yok | Eski sonuç sessizce dönüyor | ÇELİŞKİ | §5.3 |
| C-06 | MK-09: toplam iade ≤ ilk tahsilat; iade kullanılan kaynağa döner | Üst sınır yok; kaynak/cep kavramı yok | YOK | §5.3 |
| C-07 | Z05: geçiş kabulü ile mali tahsilat ayrı durumlar | Tek modal, tek sonuç | YOK | §5.2 |
| C-08 | Z13: ad/plaka ilişki anahtarı olamaz | Aktif sefer plaka ile çözülüyor (Karar #30, Increment 149) | ÇELİŞKİ (kod doğrulanacak) | §5.4 |
| C-09 | Z01: araç ataması olmadan doğrulanmış sefer oturumuyla tarife bulunur; belirsizde olay kaydedilir | Aktif `TripAssignment` (plakaya bağlı) şart; olay üretilmiyor | ÇELİŞKİ (ulaşım denetçisi doğrulayacak) | §5.4 |
| C-10 | Z03: kart yenileme hesabı değiştirmez; geçmiş korunur | Kart↔cüzdan bağı değişmiyor | KARŞILANIYOR | §5.2 |
| C-11 | Z04: karta veri yazılmaz | NFC yalnız UID okuyor, yazmıyor | KARŞILANIYOR | §5.2 |
| C-12 | MK-06: kart UID tek başına para çekme yetkisi değildir | UID ile borçlandırma mümkün | ÇELİŞKİ (Model 4) | §5.2 |
| C-13 | Z13/Z03: kişi kimliği istemciden değil kanıttan çözülür | `person_id` istemciden; master plan Increment 148 aksini söylüyor | DOĞRULANACAK | §5.2 |
| C-14 | Z12: reddedilen ve iade edilen olaylar kasa toplamından ayrı | Raporlar result ve refunded süzmez; charged_amount toplar | ÇELİŞKİ | §5.6 |
| C-15 | Z09: kesinleşmiş hareket değişmez; bakiye hareketlerden yeniden hesaplanabilir | `Reconcile` düzeltmez; pending girdiler sahte drift üretir; cüzdan cache'i doğrudan yazılabiliyor | ÇELİŞKİ | §5.1 |
| C-16 | DV-04: bakiye devri ≠ hizmet devri | Devir modeli yok | YOK | §5.5 |
| C-17 | Z01: tarife belirsizse olay kaydedilir; onaylı istisna akışı | Hata döner; olay, ledger ve RiderTrip yok; istisna durumu yok | ÇELİŞKİ | §5.4 |
| C-18 | Excel ücretsiz takip (ulaşım): tahsilatsız geçiş ve personel kaydı; Mercedes/LCW korunur | Ulaşım ücretsiz modu plaka + aktif sefer + SUL ister | ÇELİŞKİ | §5.4 |
| C-19 | Z06/Z08: sırasız olay eşleşmeyi bekler; tarife olay anına göre | Sıra bağımlı; fiyat `time.Now()` ile; açık yolculuk çıkış tarifesini alır | ÇELİŞKİ | §5.4 |
| C-20 | Excel s.11: algoritma güzergâh (kurum > güzergâh > yön) düzeyinde | Algoritma `AccessPoint.PricingMode`'da | ÇELİŞKİ (tasarım) | §5.4 |
| C-21 | Z13: ad veya plaka ilişki anahtarı olamaz | Plaka anahtar: AP ↔ TripAssignment, alarms, verifications, billing | ÇELİŞKİ | §5.4 |
| C-22 | Z02: düzeltmede kişi, gerekçe, zaman kaydedilir | Stok tarihçesinde aktör istek gövdesinden, gerekçe yok; transfer kim ve neden saklamıyor | ÇELİŞKİ | §5.4 |
| C-23 | T10 / Excel: süre ve ceza varsayılmaz; politika güzergâhta seçilir | Sefer bitince planlı mesafe × oran otomatik borçlandırılır | ÇELİŞKİ | §5.4 |
| C-24 | T11: tahsil edilmiş gibi gösterilmez | Otomatik kapanışta Debit hatası yutulup `ChargedAmount` yazılır | ÇELİŞKİ | §5.4 |
| C-25 | MK05 / Excel s.25: pasif terminal işlem yapamaz | RecordEvent `AccessPoint.Status` denetlemez | ÇELİŞKİ | §5.4 |
| C-26 | Z09: kesinleşmiş hareket değişmez; ters kayıt ile düzeltilir | Ters kayıt yok; RefundEvent Result kontrolsüz; gelir hesabı ters çevrilmiyor (ÇIKARIM) | ÇELİŞKİ | §5.4 |
| C-27 | MK18 / M16: kümülatif iade ≤ ilk tahsilat | Üst sınır yok | ÇELİŞKİ | §5.3 |
| C-28 | MK14 / M12: aynı anahtar + farklı içerik = çatışma | Eski sonuç sessizce döner (intents ve ledger) | ÇELİŞKİ | §5.3 |
| C-29 | DV03 / DV01: sürümlü sağlayıcı ataması, aktif çakışma yok | Yok; `trip_contracts` emsal (rota ve araç kapsamı yok) | YOK | §5.5 |
| C-30 | Model 1: her kurum kendi sağlayıcısını seçer | Sağlayıcı kavramı yok (Institution Type etiket) | YOK | §5.5 |
| C-31 | Z14: göç (kart, bakiye, hareket, açık yolculuk) | Yalnız stok göçü | YOK | §5.5 |

---

## 7. Doküman İçi Belirsizlikler ve Açık Kararlar

PDF'in kendisi birçok değeri bilinçli olarak "yetkililerce onaylanmalı" bırakır (s.20, s.35). Aşağıdakiler, yazılımın **seçmek zorunda kaldığı** noktalardır; her biri için önerilen varsayılan ve karar sahibi verilmiştir. "Öneri" kullanıcı onayına kadar **taslaktır**.

| ID | Belirsizlik | Neden önemli | Öneri (varsayılan) | Karar sahibi |
|---|---|---|---|---|
| AK-01 | **Topoloji ve izolasyon birimi.** M1'de "her kurum kendi sağlayıcısını seçer" ve "ortak kişi hesabı" birlikte isteniyor; sağlayıcı bir tenant mı, kurum içinde bir kayıt mı? | Mevcut mimari tenant/site veri izolasyonunu bilinçli zorluyor; seçim ortak kişi hesabını, fatura tarafını ve rakip sağlayıcı verisinin görünürlüğünü belirler. | §4.2 ve §5.5'te kod bulgularıyla birlikte. | İş sahibi + mimari |
| AK-02 | **Cep (pocket) yalnız Model 4'e mi özgü?** DV-04 "fon sorumlusu" bazlı bakiye ve DV-03 `payment_profile_id` ister; M2'de "hizmeti veren fatura keser"; MK-04 cepleri. Ortak hesapta farklı sağlayıcıların fonları yükleme anında hangi satıcıya ait? | Yükleme faturası (Z11) satıcıyı yükleme anında bilmeyi gerektirir; tek bakiyeli cüzdanda M1'in çok-sağlayıcılı durumu fon/fatura sahipliğini belirsizleştirir. | **Pocket Model 4'ten önce, M1 devir işinin temeli olarak gelir**: bakiye = (hesap, fon sorumlusu/ihraççı, program, para birimi) cebi. | Mimari |
| AK-03 | **İhraççı / fon sahibi ve hukuki profil.** Ulaşım/yemek yüklemeleri bugün kimin yükümlülüğü? Market için ortak-destekli, ulaşım için iç defter mi? | PDF s.23/s.42 yazılımın fon tutma/ödeme onayı yetkisi yaratmadığını vurgular; bu karar hukuk/finans dışı çözülemez. | Yazılım `FundingAuthority` arayüzüyle **her iki modeli** taşır (iç defter / ortak-destekli); hangisinin hangi programda geçerli olduğu yapılandırma verisi olur. | Hukuk/finans + iş sahibi |
| AK-04 | **Market ağı kapsamı:** yalnız Gürsel tenant'ındaki kurumlar mı, tenant'lar arası mı? | Tenant'lar arası ağ platform seviyesinde yeni bir kimlik/dizin katmanı demektir. | Faz 1: tenant içi; platform dizini ileri faz. | İş sahibi |
| AK-05 | **Para hassasiyeti ve yuvarlama.** Tam sayı küçük birim mi, Decimal128 mi? Yuvarlama (yukarı/aşağı/yarıya) nerede tanımlı? | Z10 "kesin hassasiyet" ister; mevcut tutarlar `float64`. | Tam sayı küçük birim (kuruş) + `Money` değer nesnesi; yuvarlama kuralı tarife sürümünde (T-CFG "para birimi/yuvarlama"). | Mimari |
| AK-06 | **A3/A4: P tahsilatı için yetersiz bakiye.** Binişte P (tavan) tahsil edilir; gerçek net N çok daha küçükse bile P'yi karşılayamayan yolcu reddedilir mi? | UX/operasyon: bakiyesi N'ye yetip P'ye yetmeyen yolcu araca binemez. PDF "reddet/borçlandır politikası belirlenmeli" der (s.12). | Güzergâh başına `InsufficientBalancePolicy`: `reject` (varsayılan), `debit_to_floor`, `allow_and_flag`. Kayıt anında **P ≥ olası tüm N** doğrulaması (ihlal = kayıt reddi). | Operasyon |
| AK-07 | **İkinci okutma yok** politikası (tavan kalır / süre sonunda otomatik kapanış / inceleme). | PDF "güzergâhta seçilmeli" der; varsayım üretilmemeli. | `NoExitPolicy` enum + zaman aşımı; varsayılan: P kalır + inceleme kuyruğu; N tahmin edilmez. | Operasyon |
| AK-08 | **A1 mesafe kaynağı ve okutma anı.** Mesafe: durak tablosu mu, iki nokta GPS'i mi, mesafe sağlayıcısı mı? Sabah binişte/akşam inişte okutma cihaz akışında nasıl ayrışır? | Algoritma ve cihaz UX'i buna bağlı. | `DistanceSource` enum (`stop_table`, `gps_haversine`, `distance_provider` — mevcut `distance_providers` modülü yeniden kullanılır); okutma anı güzergâh-yön ayarının parçası. | Operasyon |
| AK-09 | **Çevrimdışı risk politikası ve QR profili sayısal değerleri** (azami süre/tutar/adet, önbellek yaşı, QR TTL/yeniden kullanım). | PDF s.20: "Limitler, QR profili … yetkililerce onaylanmalıdır." | Çevrimdışı kabul **varsayılan kapalı**; değerler kurum/hesap profili verisi, onayla atanır. | İş sahibi + güvenlik |
| AK-10 | **Taban (floor) politikası çok-kurumlu kişide hangi kuruma bağlı?** | Bir kişinin birden çok `InstitutionAffiliation`'ı olabilir; Z07 "kurum/hesap politikası" der. | Hizmeti veren **atamanın** kapsamındaki kurumun politikası; belirsizse en kısıtlayıcı. | Mimari |
| AK-11 | **Kütüphane ücret modeli** (PDF: araştırılacak — süre mi sabit mi). | M-05 açık karar. | Açık-oturum soyutlaması (P7) ile ödünç/iade eşlemesi; (a) sabit işlem ücreti, (b) süre bazlı, (c) yalnız gecikme cezası seçenekleri yapılandırılabilir; karar iş sahibinde. | İş sahibi |
| AK-12 | **Rapor tanımları** (kullanım ortalaması pay/payda/dönem/iptal sefer; kasa kalemleri). | PDF "belgelenmeli" der ama tanımı vermez; s.20 yetkili onayı ister. | Yazılı tanım + test; iş sahibi onayı. | İş sahibi |
| AK-13 | **Devir:** T için saat kaynağı/zaman dilimi, "saat güvenilmez" eşiği, iki katmanlı atama (kurum↔ana operatör, ana operatör↔taşeron). | DV-03/DV-05; M-03'te fatura kesen taraf taşeron olabilir. | T sunucu saatine (UTC) göre, kurum yerel saat dilimiyle gösterilir; saat sapması eşiği yapılandırma verisi; atamada `parent_assignment_id` ile alt atama. | Mimari + operasyon |
| AK-14 | ~~Excel ekleri elde yok~~ → **çözüldü** (PDF dökümü alındı) | Test kataloğu ve karar listesi tamamlandı | Eşleme: `AYPASS_REV06_EXCEL_ESLEME_TR.md` | — |

**Excel karar listeleri ile eşleme.** P01–P11 (Z), D01–D14 (market) ve DK01–DK12 (devir) kararları Excel'de sahipleriyle tanımlı ve bekliyor. AK-01…AK-13 bu kararların analiz tarafındaki karşılıklarıdır. Eşleme tabloları `AYPASS_REV06_EXCEL_ESLEME_TR.md` §5 bölümündedir. Öneri kararları AD-01…AD-15 olarak §11'de numaralanmıştır.

---

## 8. İş Planı (Artımlar)

Sıralama: 0 → 1 → 2, ardından 3 ve 4 paralel, sonra 5, 6-7, 8, 9. Büyüklükler göreli ölçüdür (S/M/L/XL); takvim tahmini içermez. Kod değişikliği gerektiren her faz için README kontrol listesi uygulanır (28 dil, menü ve görünüm profilleri, seed, gerçek MongoDB testi, gerçek `go run`).

| Faz | Kapsam | Gereksinim | Bağımlılık | Büyüklük | Onay |
|---|---|---|---|---|---|
| **0** | **P0 düzeltmeleri:** (0.1) `RefundEvent` yalnız gerçekten tahsil edilmiş olayları iade etsin; `insufficient_balance` için `ChargedAmount` = 0 veya ayrı alan. (0.2) GİB SOAP taşıma hatası hata olarak dönsün; mock UUID yalnız test ortamında; `IssueEInvoice` nil guard ve mevcut `ExternalID` varsa yeniden gönderme yok. (0.3) Ödeme: toplam iade ≤ capture (`refunded_amount`); aynı anahtar + farklı payload → çatışma hatası; refund anahtarı zorunlu. (0.4) Cüzdan: koşullu güncelleme (`balance == expected` filtresi + yeniden deneme) veya `$inc` + taban filtresi; gerçek MongoDB'de goroutine'li eşzamanlılık testi. (0.5) Callback: `succeededCb` hatası kaydı kalıcı yapmasın; yeniden işlenebilir durum (`credit_pending`). (0.6) Altyapı: canlı Mongo replica set doğrulaması; `journal_entries(journal_id)` index; `ledger_accounts` için `(tenant, site, code)` unique (önce veri denetimi). (0.7) `approvals`, webhook retry ve queue topology için bağlama/silme kararı. (0.8) `TripAssignment` güncellemesine onay, devir ve sürücü alanlarını ekle. (0.9) `RecordEvent` içinde pasif AccessPoint'i reddet. (0.10) Sabit kimlik bilgilerini koddan ve git'ten kaldır; rotasyon (kullanıcı ve ops; değerler yazılmadı). (0.11) Stok plaka değişiminde aktörü oturumdan al; gerekçe alanını zorunlu yap. | Z09, Z10, Z11, MK-07, MK-09 | yok | S–M | **Onay gerekli** (AD-10) |
| **1** | Para temeli: minor-unit `int64`; `occurred_at` ve `received_at`; eşzamanlı güncelleme; yuvarlama sürümlü tarifede | Z09, Z10, P1, P2 | Faz 0 | L (migration) | AD-01 |
| **2** | İdempotensi ve outbox: DB-backed idempotency, payload parmak izi, `platform/outbox`, ledger reversal | Z06, Z09, Z11, MK-08 | Faz 1 | M | AD-09 |
| **3** | Ulaşım algoritmaları: A3/A4 (tahsil + iade, iki bağlı journal), A4 matris (`route_stop_fares` ile), tarife sürüm kimliği olaya yazılsın, geç ve sırasız olay | T-A3, T-A4, Z01, Z02, Z08 | Faz 1–2; C-08/C-09 kararı | L | AD-02, C-08, C-09 |
| **4** | Sağlayıcı ataması ve devir: `service_assignments`, devir durum makinesi, bakiye devri A/B/C | M-02…M-04, DV-01…DV-08, Z14 | Faz 1 | L | AD-03 |
| **5** | Çevrimdışı dayanıklılık: Flutter kalıcı kuyruk, batch sync (sunucu ack'inden sonra silme), QR tek kullanım ve revoke, cihaz saat sapması, geçiş kabulü ile tahsilat ayrımı, reddedilen olay kaydı | Z04, Z05, Z06, Z08, MK-12 | Faz 2 | XL (Flutter + backend) | AD-06, AD-08 |
| **6** | Fatura işi: `invoice_jobs` durum makinesi, GİB gerçek hata sınıfları, dış referans sorgusu | Z11 | Faz 0–2 | M | AD-05 |
| **7** | Raporlama: kasa kalemleri, reddedilen/iade filtresi, güzergâh ortalaması tanımı, araç ve personel türü, as-of | Z12 | Faz 1 | M | AK-12 |
| **8** | Model 4 çekirdeği: cep, `FundingAuthority`, merchant, hakediş, maker-checker, belge profili | MK-01…MK-15 | Faz 1, 2, 4 + hukuk/finans onayı | XL | AK-03 (hukuk/finans), partner seçimi |
| **9** | Canlıya geçiş: Z14 göçü, mutabakat, pilot | Z14, kabul grupları | Hepsi | L | Pilot ve mutabakat onayı |

**Uygulama durumu (Faz 0, 2026-10-28):**

- Commit'ler (`aypass-rev06-faz0` dalı, push edilmedi): `b2398bd` cüzdan/iade/callback/ledger; `884b0a7` GİB; `c79cc9b` trip onay alanları, stok aktörü, migrate kaynağı. `eb581c7` kullanılmayan `approvals` ve `governance/approval.go` kodunu sildi (AD-04). İlk commit kendi başına `go build` ve `go vet` ile doğrulandı; HEAD derleniyor.

- ✅ **0.1 RefundEvent:** Yalnız başarılı ve cüzdandan tahsil edilmiş olaylar iade edilir (`Result == granted`, `LedgerEntryID` dolu). Yetersiz bakiye olayı iade edilemez. Test: `TestRefundEvent_RejectsInsufficientBalanceEvent`. Borç tutarı (`ChargedAmount`) nakit yerleşimi için olduğu gibi korundu.
- ✅ **0.8 TripAssignment:** Güncelleme onay, devir ve sürücü alanlarını yazar. Test: `TestUpdate_PersistsApprovalAndDriverFields` (gerçek MongoDB; `AYSTEK_TEST_MONGO_URI`).
- ✅ **0.9 Pasif erişim noktası:** Olay kaydı, nakit yerleşim ve walk-on reddedilir (`ErrAccessPointInactive`). Test: `TestRecordEvent_RejectsInactiveAccessPoint`. İptal işlemi düzeltici adım olduğu için kapsam dışı bırakıldı.
- ✅ **0.11:** Stok plaka ve araç güncellemelerinde aktör oturumdan alınır; gövdedeki `updated_by` kaldırıldı. Plaka gerçekten değişirse gerekçe zorunludur (kırpılmış, 1–500 karakter; boş ve uzun gerekçe `stock.vehicles.validation.plate_reason_required` / `plate_reason_too_long` ile reddedilir); gerekçe `stock_vehicle_plate_history` kaydına `reason` olarak yazılır. Panel gerekçe alanını ekledi. Commit'ler: `93b4232`, `25abcc1` (backend), `7b03fe7` (panel). Açık: plaka geçmişi yazımındaki hata önceden olduğu gibi yutuluyor (`UpdateVehicle` ve `AddPlate`); geçmiş kaydı olmadan plaka değişebilir.
- ◐ **0.10 (kısmi):** `cmd/migrate_gursel` kaynak bağlantısı artık `MIGRATE_SOURCE_URI` veya `-src-uri` ile alınır; koddaki kimlik bilgisi (kod ve yorum) kaldırıldı. Kullanıcı kararı: `gursel_turizm_seed.go` içindeki düz metin test değerleri koda kalır (test ortamı için); canlıda değiştirilecek. Gürsel GPS sağlayıcı yapılandırmaları (kimlik bilgileri dahil) artık yalnızca ilk oluşturulurken yazılır (`c27c3b5`); tam zincir mevcut kayıtları ezmez. Üretim compose dosyasında `SEED_ON_START` `false` (`ccb13e2`). Açık kalanlar: (b) git'te izlenen `cmd/migrate_gursel/migrate_gursel` binary'si, (c) geçmiş commit'lerdeki değerlerin rotasyonu (değerler git geçmişinde mevcut).
- ✅ **0.2 GİB:** Taşıma, HTTP ve parse hataları artık hata döner; sahte UUID üretim yolundan kaldırıldı. Zaten gönderilmiş fatura (ExternalID dolu) yeniden gönderilmez (`ErrAlreadyIssued`, HTTP 409). Başarı olmadan `issued` yazılmaz; `nil` sonuçta panik yok. `QueryStatus` durum alanı yoksa `ACCEPTED` dönmez, hata döner. Testler: `adapters/gib_soap_test.go`, `service/issue_einvoice_test.go`, `handler/issue_einvoice_test.go`. Açık: eşzamanlı iki gönderim için atomik koşullu güncelleme gerekir (kapsam dışı).
- ✅ **0.4 Cüzdan eşzamanlılığı:** Mutlak `$set` (`UpdateCachedBalance`) kaldırıldı; yerine `ApplyBalanceChange` geldi: tek `FindOneAndUpdate` ile `$inc`, filtrede taban (`$expr`), durum ve journal tekrarı koşulu. Borçlandırma ve iade önce rezervasyonla, sonra journal ile yazılır; journal yazılamazsa rezervasyon telafi edilir. Son 50 journal penceresi eşzamanlı çift havale onayını engeller. Testler gerçek MongoDB'de, `-race` ile: 30 eşzamanlı borçlandırma tabanı aşmadı ve bakiye ledger ile eşleşti; 10 eşzamanlı aynı anahtar tek kez uygulandı; 5 eşzamanlı havale onayı tek kredi üretti; sahte ledger hatasında bakiye değişmedi. Açık: çökme anında cache ledger'ın altında kalabilir (Reconcile'da negatif drift olarak görünür, harcanabilir bakiye oluşmaz).
- ◐ **0.6 Ledger index:** `journal_entries(journal_id)` eklendi (unique değil); test: `TestVerifyModule_LedgerJournalIDIndexIsRegistered`. Yarış güvenliği (`7eb7bfa`): sistem hesabı yeni kayıtlarda deterministik `_id` ile insert-only açılır (`model.SystemAccountID`), mevcut hesaplar değişmez; `ConfirmPendingJournal` koşullu güncelleme yapar, ikinci onay `ErrJournalNotPending` döner. Testler: 20 eşzamanlı ilk `ensureSystemAccount` çağrısı tek hesap, 10 eşzamanlı onay tek başarı; `-race` ile geçti. Açık: `ledger_accounts (tenant, site, code)` unique indeksi veri denetiminden sonra; soft-delete sonrası yeniden açma ve eski ile yeni binary'nin aynı anda çalıştığı deploy penceresinde mükerrer hesap riski tamamen kapanmadı; canlı replica set doğrulaması yapılmadı; `ConfirmBankTransfer` yeniden denemesi, journal posted olup bakiye uygulanamadıysa bakiyeyi uygulamıyor (0.4 ile aynı tipte açık not).
- ◐ **0.3 Ödeme iadesi ve idempotency:** İade rezervasyonu tek koşullu `UpdateOne` + `$inc` ile yapılır (`refunded_amount + talep ≤ tahsil`); sağlayıcı hatasında `ReleaseRefund` telafi eder. Anahtar zorunlu; aynı anahtar + aynı tutar replay, farklı tutar `ErrIdempotencyConflict` (422); işlenmekte olan anahtar 409; başarısız deneme aynı anahtarla yeniden denenebilir. Ters yönlü `Update` artık `refunded_amount` yazmaz. Returns iade anahtarı `return:<returnID>` oldu (saniye çözünürlüklü `Number` çakışması giderildi). Testler gerçek MongoDB'de: sınırdaki iade, 10 eşzamanlı iade (90 TL, tek kabul), telafi, çakışma, replay. Açık: (a) Panel iade isteğinde `idempotency_key` gönderir (`usePaymentIntentRefund`). Anahtarsız iade isteği artık reddedilir (400, `payment_intents.idempotency_key_required`); backend istek başına anahtar üretmez (`a00f69d`). (b) `refunded_amount` alanı önceden iade almış kayıtlarda yok, 0 sayılıyor; backfill gerekir (kapsam dışı). (c) Sağlayıcı tarafında iade anahtarı gönderilmiyor. (d) Süreç çökmesinde kayıt `pending` kalır ve o anahtar `in_progress` döner; elle müdahale gerekir. (e) Yeni hata anahtarları (`refund_exceeds_capture`, `idempotency_key_required`, `idempotency_in_progress`, `callback_credit_pending`, GİB `already_issued`, `access_point.error.inactive` ve diğerleri) 28 dilde eklendi (`module_error_translations_seed.go`, commit `931e70f`).
- ✅ **0.5 Callback replay:** Operasyon kaydı, kanca (`succeededCb`) başarılı olduktan sonra yazılır; kanca hatası `ErrCallbackCreditPending` olarak 503 döner, sağlayıcı yeniden teslim eder. Testler gerçek MongoDB'de: başarılı yeniden teslimde tek kredi, kanca hatasından sonra yeniden teslimde kredi, bekleyen kredi yeniden denenebilir. Açık: `syncOrderPayment` hata verirse yeniden teslimde senkron tekrar çalışmaz (önceden var olan sorun, kapsam dışı).
- ◐ **0.7 `approvals` (AD-04, karar: sil):** `internal/modules/approvals` ve import edilmeyen `integrations/governance/approval.go` silindi. Silme öncesi bulgular: `Approve` öz-onayı engellemiyordu, model site alanı taşımıyordu, kimlikler uuid string'ti. Bu konuyu kapsayan tek testler yer tutucuydu (`assert.NotNil(t, t)`, "self-approval prevention" dahil) ve kaldırıldı. Yeniden tasarım uygulandı: ObjectID, tenant+site kapsamı, istek sahibi ≠ onaylayan (servis katmanında, gerçek Mongo testiyle), kalıcı karar kaydı (approvals modülü; politikalar, istekler, B2B zinciri ve iade kapısı). Kapılar: manuel sefer istisnası (`trip_assignments.approval_status`, onay kapısı uygulandı: `approve` izni ve öz-onay yasağı), IBAN değişikliği çift kontrolü (MK-05, MK-13), devir kapıları (DV-02), ödeme emri onayı (MK-03). Yetki modeli (kullanıcı kararı, 2026-10-28): roller şimdilik `admin`; onay yetkisi rol adına değil izin anahtarına bağlı (`fleet_tally_trip_assignments:approve`). Roller ve izinler admin API'sinden dinamik yönetilir; kod rol adı kontrolü yapmaz. Uygulandı: manuel sefer onay endpoint'i `execute` yerine `approve` iznini ister; seferi kaydeden kullanıcı onaylayamaz (`ErrSelfApproval`, 403). Testler: servis (öz-onay reddi, bağımsız onay, kaynağı olmayan eski kayıt) ve middleware (tam aksiyon eşleşmesi). Açık: (a) prod'da `admin` rolüne izin senkronu için `cmd/seed --roles` çalıştırılmalı; seed prod'da otomatik çalışmaz. (b) Özel roller için yalnızca `approve` yetmez; coarse modül kapısı için `fleet_tally:*` önekli bir izin de gerekir. (c) Kaynağı olmayan eski manuel seferler öz-onay kontrolünden muaf. (d) Yeni hata anahtarları 28 dilde eklendi. Kapı listesi (IBAN, devir, ödeme emri) kararı hâlâ açık.

---

## 9. Kabul Test Kataloğu

Excel'deki kabul setleri birebir alındı: **K01–K24** (Z kabul), **MK01–MK30** (market), **DT01–DT26** (devir), **T01–T15** (ulaşım kontrol senaryoları) ve **FX1–FX9** (formül örnekleri). Excel'de hepsi "Tanımlandı; çalıştırılmadı". Aşağıdaki durumlar kod okumasına dayanır; **hiçbir test çalıştırılmadı**. Test-bazlı tam tablo: `AYPASS_REV06_EXCEL_ESLEME_TR.md`.

| Set | Toplam | KARŞILANIYOR | KISMİ | YOK | ÇELİŞKİ |
|---|---|---|---|---|---|
| K01–K24 (Z kabul) | 24 | 1 | 7 | 8 | 8 |
| MK01–MK30 (market) | 30 | 0 | 12 | 12 | 6 |
| DT01–DT26 (devir) | 26 | 0 | 3 | 22 | 1 |
| T01–T15 (ulaşım) | 15 | 1 | 6 | 2 | 6 |
| FX1–FX9 (formül) | 9 | 1 | 6 | 1 | 1 |
| **Toplam** | **104** | **3** | **34** | **45** | **22** |

Gereksinim düzeyi: Z01–Z14 → 0 karşılanıyor, 5 kısmi, 2 yok, 7 çelişki. M01–M24 → 0 karşılanıyor, 10 kısmi, 10 yok, 4 çelişki. DV01–DV20 → 0 karşılanıyor, 4 kısmi, 16 yok.

**Yorum:** Hiçbir kabul grubu tam karşılanmıyor. En yakın alanlar: kart ve QR'ın karta yazmaması (K06), yükleme idempotensi (K18), yönlü O-D matrisi (T07). Ortak eksikler: sürümlü sağlayıcı ataması, çevrimdışı katman (kuyruk, sıra, olay zamanı), outbox, cep (pocket), hakediş, fatura işi.

---

## 10. Riskler

| ID | Risk | Etki | Olasılık (öneri) | Azaltma |
|---|---|---|---|---|
| R-01 | Tahsil edilmemiş ücretin iadesi (P0-1) | Doğrudan para üretimi | Orta (admin tetikli) | Faz 0.1 |
| R-02 | Sahte e-fatura "issued" durumu (P0-2) | Yasal uyum, yanlış muhasebe | Bilinmiyor (prod GİB ortamı doğrulanmalı) | Faz 0.2, Faz 6 |
| R-03 | `float64` yuvarlama sapması | Kuruş kayması, mutabakat kopması | Yüksek (ulaşım ücretinde yuvarlama yok) | Faz 1 |
| R-04 | Eşzamanlı borçlandırmada taban aşımı (P0-4) | Sınırsız negatif bakiye | Orta | Faz 0.4 |
| R-05 | Callback replay'inde kredi kaybı (P0-5) | Tahsil edilmiş ama kredisiz intent | Düşük–orta | Faz 0.5 |
| R-06 | Çevrimdışı olay kaybı veya ağ kesintisinde çift borçlandırma | Müşteri şikayeti, mutabakat farkı | Yüksek (uygulama davranışı) | Faz 5; geçici olarak UI'da "belirsiz" durumu |
| R-07 | Replica set yokluğunda transaction hatası veya tutarsızlık | Ledger yazımı başarısız olabilir | Bilinmiyor (canlı doğrulanmadı) | Faz 0.6 |
| R-08 | Model 4 hukuki yapısı (ödeme kuruluşu izni, finansal ortak seçimi) | Proje durdurma | Yüksek | MK-03 hukuk onayı; kod öncesi karar |
| R-09 | Tenant izolasyonu ile kurumlar arası market ağı çelişkisi | Mimari yeniden kurgu | Yüksek | AK-04, AD-03 |
| R-10 | Eski verinin göçü (Z14) | Mükerrer tahsilat veya fatura | Orta | Z14 göç planı, Faz 9 |
| R-11 | Test kapsamı: sahte repository, eşzamanlılık ve gerçek Mongo testi yok | Regresyonlar gizli kalır | Yüksek | Faz 0'da gerçek MongoDB testleri |
| R-12 | Bağlı olmayan kod (webhook retry, queue topology) | Yanlış güven, bakım yükü | Düşük–orta (`approvals` silindi) | AD-05 |
| R-13 | `invoiceable_records` 1000 journal sınırı | Faturalanacak kayıtlar eksik listelenir | Orta | Faz 6 |
| R-14 | Mock adaptörün production dışında imzasız kabulü | Test ortamında sahte ödeme | Düşük (production korumalı) | Faz 0 kontrolü |
| R-15 | Sabit kimlik bilgileri koddan ve git'ten erişilebilir: `cmd/migrate_gursel` içinde bağlantı dizgisi, git'te izlenen derlenmiş binary, seed'de düz metin anahtar değerleri | Kimlik bilgisi sızıntısı; kaynak veritabanına erişim | Yüksek (depo erişimi olan herkes) | Faz 0.10: kaldır, rotasyon, binary'yi izlemeden çıkar (değerler yazılmadı) |
| R-16 | Pasif AccessPoint (`Status`) denetlenmiyor | Pasif terminalden geçiş ve tahsilat | Orta | Faz 0.9 |
| R-17 | Onay alanları kalıcı yazılmıyor (`TripAssignment` güncellemesi) | Manuel istisna ve onay akışı çalışmıyor | Yüksek | Faz 0.8 |
| R-18 | Plaka değişimi sonrası faturalama sessizce atlanıyor; rapor hata veriyor | Eksik fatura ve hakediş | Orta–Yüksek | AD-14; Faz 1 |
| R-19 | Stok plaka tarihçesinde aktör istek gövdesinden | Denetim izi güvenilmez | Orta | Faz 0.11 |
| R-20 | Topoloji B/C (sağlayıcı veya kurum başına tenant) ortak kişi hesabını bozar | Mimari yeniden kurgu | Yüksek | AD-13 |

---

## 11. Karar Günlüğü Önerileri (Master Doküman §12'ye İşlenecek)

Tüm kararlar henüz onaylanmadı. Her biri için öneri verilmiştir; onay gelmeden kod değişikliği yapılmaz.

**Not:** Excel'deki market kararları (D01–D14) ile bu öneriler ayrı numaralanır: öneriler AD-xx, Excel kararları D, P ve DK. Eşleme: `AYPASS_REV06_EXCEL_ESLEME_TR.md` §5.


| No | Konu | Öneri | Gerekçe | Durum |
|---|---|---|---|---|
| AD-01 | Para temsili | `int64` minor-unit + para birimi; yuvarlama tarife sürümünde | Z10 kesin hassasiyet | Onay bekliyor |
| AD-02 | A3/A4 modeli | Tahsil + iade (iki bağlı journal); girişte hold yok | Doküman s.9: fiili tahsilat P, iade P−N | Onay bekliyor |
| AD-03 | Sağlayıcı ataması | Yeni modül `service_assignments`; `trip_contracts` genişletilmez | Kavram yok; çakışma kuralı ve `valid_to` gerekli | Onay bekliyor |
| AD-04 | `approvals` | Sil ve Rev.06 kapılarına göre yeniden tasarla | Kullanıcı kararı: sil (2026-10-28). Yeniden tasarımın kapı listesi ve rol matrisi kararı bekliyor | Silindi; yeniden tasarım uygulandı (2026-10-28); kapı listesi açık |
| AD-05 | Webhook retry ve queue topology | Bağla (callback kayıpları için) | P0-5 | Onay bekliyor |
| AD-06 | Çevrimdışı kapsamı | Yalnız onaylı profil ve kart; varsayılan kapalı; limitler kurum profilinde | AK-09; Z05, Z08 | Onay bekliyor |
| AD-07 | Negatif bakiye varsayılanı | Taban kurum politikasından gelsin; −15 evrensel varsayılan değil | Z07 | Onay bekliyor (AK-06) |
| AD-08 | Reddedilen okutmalar | Kaydedilsin (`denied` olay) | Z06 izlenebilirlik | Onay bekliyor |
| AD-09 | Mongo dağıtımı | Replica set zorunlu; canlıda doğrulanacak | Transaction gerekli | Onay bekliyor (altyapı) |
| AD-10 | P0 düzeltmelerinin zamanlaması | Rev.06 işinden önce, ayrı artım | Para ve yasal risk | Uygulandı: Faz 0, kullanıcının "İş planlarını yerine getirmeye başla" talimatıyla yürütüldü |
| AD-11 | Ücretsiz ulaşım takibi | Ulaşım ücretsiz modu plaka + aktif sefer zincirinden bağımsız olsun (Excel: "korunur") | Mercedes/LCW personel takibi bu zincire bağlı kalırsa kırılır | Onay bekliyor |
| AD-12 | Algoritma seçimi düzeyi | Güzergâh (kurum > güzergâh > yön) düzeyinde ayar; `AccessPoint.PricingMode` yalnız cihaz ve yetki davranışı | Excel s.11; kod AP düzeyinde | Onay bekliyor |
| AD-13 | Topoloji | A (tek tenant ve site, `institutions` supplier); B ve C ortak kişi hesabını bozar | Mevcut kod ve izolasyon kuralları | Onay bekliyor (AK-01) |
| AD-14 | Plaka anahtarı | Ana ilişki anahtarı stok araç `_id`; plaka değişebilir gösterim | Z02, Z13 | Onay bekliyor |
| AD-15 | Güvenlik | Sabit kimlik bilgilerini koddan ve git'ten kaldır; binary'yi izlemeden çıkar; rotasyon | Düz metin değerler ve git'te binary. Kullanıcı kararı: seed test değerleri koda kalır, canlıda değiştirilir | Kısmi karar; GPS kimlik bilgileri yalnızca ilk oluşturmada yazılır (`c27c3b5`); binary ve rotasyon açık; değerler yazılmadı |

---

## 12. Rol ve Yetki Matrisi (Kapı Kararı İçin Öneri, 2026-10-28)

Bu bölüm §11'deki açık kapı kararını (IBAN çift kontrolü, devir kapıları, ödeme emri onayı) karşılar. Kullanıcı rolleri şimdilik `admin` tutmak, yetkinin sistem devreye alındığında dinamik devredilebilmesini istedi. Her eşleme izin anahtarına (`modül:aksiyon`) dayanır; rol adı koda yazılmaz. Roller ve izinler admin API'sinde oluşturulur ve rollere atanır.

### 12.1 İlkeler

- **Yetki izin anahtarıdır.** Roller veridir (README: yönetici tarafından yönetilen davranış modellenir, koda gömülmez).
- **Admin tüm izinleri alır** (README RULE #2.1, seed ile). Bu nedenle zincirin güvenliği rol ayrımına değil, kod içindeki "yapan ≠ onaylayan" kuralına dayanır. Admin'in bir işi hem hazırlayıp hem onaylaması yönetim politikasıyla yasaklanmalı; kodda zorunlu olan yalnızca öz-onay kontrolüdür.
- **Her rol, kullandığı modülde okuma izni almalıdır** (ör. `wallets:read`). Genel `PermissionMiddleware` yalnızca modül önekini (`wallets:`) arar; `wallets_topup_sources:read` gibi bir izin bu kapıyı geçmez.
- **Yeni modül izinleri** (`merchant_bank_accounts`, `payment_orders`, `settlements`, `service_transfers`) seed modül listesine eklenmeden katalogda yoktur.

### 12.2 Öneri Roller

| Rol anahtarı | Görev | Yapabilir (maker) | Onaylayabilir (checker) | Kısıt |
|---|---|---|---|---|
| `kurum_yoneticisi` | Kurum yöneticisi | Sefer kaydı (`fleet_tally_trip_assignments:execute`) | Manuel sefer onayı (`fleet_tally_trip_assignments:approve`); devir kapısı kurum onayı (`service_transfers:approve`, ileride) | Yalnızca kendi kurumu; kendi kaydettiği seferi onaylayamaz |
| `kasa_gorevlisi` | Kasa / market kasası | Kasa tahsilat kaydı (`access_points:execute`) | — | Onay ve iade yetkisi yok; yalnızca kendi terminali |
| `bilet_gorevlisi` | Biletçi (walk-on, biniş) | Biniş ve bilet kaydı (`access_points:execute`) | — | Onay ve iade yetkisi yok; yalnızca kendi noktası |
| `mali_hazirlayici` | Mali operatör | IBAN değişiklik talebi (`merchant_bank_accounts:update`, ileride); ödeme emri hazırlama (`payment_orders:create`, ileride) | — | Hazırladığını onaylayamaz |
| `mali_onayci` | Mali onaycı | Havale onayı (`wallets:approve`); IBAN onayı (`merchant_bank_accounts:approve`, ileride); ödeme emri onayı (`payment_orders:approve`, ileride); hakediş onayı (`settlements:approve`, ileride) | Yukarıdakilerin tümü | Hazırladığı işi onaylayamaz |
| `devir_yoneticisi` | Devir yöneticisi | Devir adımlarını hazırlar (`service_transfers:create`, `service_transfers:update`, ileride) | — | Onay yapmaz |
| `devir_onayci` | Devir onaycısı (kapı bazlı) | — | Devir kapıları (`service_transfers:approve`, ileride); kapı–rol eşlemesi veri olarak tutulur | Her kapıda hazırlayan ≠ onaylayan |
| `odeme_gonderici` | Bankaya gönderen (opsiyonel, ayrı rol önerilir) | Ödeme emrini gönderir (`payment_orders:execute`, ileride) | — | Gönderildi ≠ ödendi; onaylayan göndermez |
| `admin` | Tenant yöneticisi | Seed ile tüm izinler | Tümü (politika ile sınırlı) | Yukarıdaki yönetim politikası geçerli |

Platform kullanıcıları (`user_type=platform`) mevcut politikayla modül RBAC'ini atlar; tenant rolü değildir.

### 12.3 Kapı → İzin Eşlemesi

| Kapı | Kaynak | Maker izni | Checker izni | Kural | Kod durumu |
|---|---|---|---|---|---|
| Manuel sefer istisnası | Z01, §5 P1-a | `fleet_tally_trip_assignments:execute` | `fleet_tally_trip_assignments:approve` | Kaydeden ≠ onaylayan | ✅ Uygulandı (`45a2e15`); test var (BAĞIMSIZ) |
| Banka havalesi onayı | Z10 | Müşteri beyanı (`/me/wallet/bank-transfer-declarations`) | `wallets:approve` | Beyan eden ≠ onaylayan | ✅ Uygulandı (bu tur); rota izni + servis kontrolü + test (BAĞIMSIZ) |
| Ödeme niyeti: confirm, capture, cancel, iade | MK-09, Z10 | `payments:execute` | Eşik üstü iade için ayrı onay (karar bekliyor, §12.4) | Toplam iade ≤ tahsil | ◐ Granüler izin kapısı uygulandı (bu tur); eşik onayı yok |
| B2B bayi sipariş onayı | Rev.06 dışı (mevcut kod) | — | `b2b_dealers:approve` ve adım rolü (oturumdan) | Sipariş veren ≠ onaylayan | ✅ Zincir politikadan üretilir (bu tur); rota izin kapısı, sunucu tarafı rol ve öz-onay uygulandı |
| IBAN / finansal kimlik değişimi | MK-05, MK-13 "çift kontrol" | `merchant_bank_accounts:update` | `merchant_bank_accounts:approve` | Hazırlayan ≠ onaylayan; onay olmadan yeni IBAN'a ödeme yok | ⏸ Kodda yok (üye/market modülü yok) |
| Ödeme emri (fon / market ödemesi) | MK-03 "ödeme emrini kimin aldığı ve onayladığı" | `payment_orders:create` | `payment_orders:approve`; gönderim `payment_orders:execute` | Üç ayrı rol önerilir; gönderildi ≠ ödendi | ⏸ Kodda yok (hakediş/ödeme emri modülü yok) |
| Hakediş | MK-10 | `settlements:create` | `settlements:approve` | Hesaplayan ≠ onaylayan; hesaplandı → onaylandı → gönderildi → banka teyidi | ⏸ Kodda yok |
| Devir kapıları (6 kapı) | DV-02 | `service_transfers:create` / `:update` | `service_transfers:approve` (kapı–rol eşlemesi veri) | Her kapıda hazırlayan ≠ onaylayan; ihtilafta otomatik kapanış yok | ⏸ Kodda yok (devir modülü yok) |

### 12.4 Açık Karar Noktaları (Kullanıcı)

1. **İade eşiği (MK-09):** mekanizma uygulandı (§12.7): eşik, aktiflik ve checker izni `approval_policies` verisinde, yönetici API'siyle değişir. Varsayılan kapalı. Canlıda hangi tutarın üstünün onay gerektireceği değeri kullanıcı kararıdır.
2. **Ödeme emri (MK-03):** üç rol (hazırlayan, onaylayan, gönderen) zorunlu mu, iki rol yeterli mi?
3. **Devir kapıları (DV-02):** altı kapının her biri için onaylayan rol (kurum, mali, hukuk) eşlemesi.
4. **Admin politikası:** admin, onay rolü ile hazırlama rolünü aynı kişide tutabilir mi? (Öneri: hayır.)
5. **B2B sipariş eşikleri:** 50.000 ve 100.000 eşikleri ile onaylayıcı rol adları (`store_manager`, `credit_dept`, `executive`) kodda sabit. Veri olarak yönetilmesi README'nin dinamik yönetim kuralına uygun olur; onaylıyor musunuz?

### 12.5 Bu Turda Kodla Doğrulanan Açıklar

- **B2B PO onay (kapatıldı):** `b2b_dealers:approve` izni gerekiyor; onaylayıcı rolü artık oturumdan alınıyor. İstek gövdesindeki `role` alanı kaldırıldı; Go JSON bağlama bilinmeyen alanları yok saydığı için panel isteği bozulmaz. Gönderen siparişi onaylayamaz veya reddedemez. Gönderen bilgisi yeni `submitted_by` alanında tutulur; bu alan boş olan eski siparişlerde öz-onay kontrolü uygulanmaz.
- **Havale onayı:** önceden granüler izin ve beyan eden ≠ onaylayan kontrolü yoktu. Düzeltildi (bu tur).
- **Ödeme niyeti rotaları:** capture, cancel, refund ve confirm yalnızca genel `payments:*` kapısındaydı; `payments:read` sahibi iade yapabiliyordu. Düzeltildi (bu tur, `payments:execute`).
- **Sabit rol adları ve eşikler:** B2B onay zinciri kodda sabit; README'nin dinamik yönetim kuralına aykırı.

### 12.6 Doğrulama Durumu

- Rota, middleware ve servis bulguları dosya yollarıyla kodda doğrulandı (BAĞIMSIZ).
- Havale self-confirm, sefer öz-onay ve izin tam eşleşmesi testleri geçti. Rota kapıları HTTP düzeyinde test edilmedi.
- Rol tasarımı bir öneridir (ÇIKARIM); kullanıcı onayı gerekir.
- Yeni hata metinleri düz İngilizce (modül stiliyle aynı). README çeviri kuralı gereği 28 dilli anahtarlara taşınmalı; bu borç açık.

### 12.7 Uygulanan Yapı: Dinamik Onay Politikaları (2026-10-28)

Kullanıcı kararı: eşik, aktiflik, checker izni ve rol adları dinamik olmalı; rol adları sabit değildir, İngilizce makine anahtarları kullanılır. Dış araştırma (maker-checker ve ayrım-görevi kaynakları) şu ilkeleri verdi: dört göz kontrolü yüksek tutarda uygulanır ama tek yüksek eşik küçük tekrarlı işlemleri denetimsiz bırakabilir; başlatan kişi onaylayamaz; serbest bırakma onay durumuna bağlı olmalı; onay ayarlarını değiştirme yetkisi kısıtlanmalı; rol anahtarı değişmez, görünen ad dile göre ayrı tutulmalı.

- **Kapı kataloğu (kod):** teknik kapı anahtarları. Şu an tek kapı: `payments.refund` (modül `payments`). Yeni kapı eklemek kod değişikliğidir; kapının ayarları veridir.
- **Politika (veri, tenant ve site):** `enabled`, `threshold` ve `currency`, `required_permission`. İzin, katalogda var olduğu doğrulanarak seçilir. Politika kaydı yoksa varsayılan: kapalı, izin `<modül>:approve`.
- **Onay talebi (veri):** talep açıldığında gerekli izin snapshot'lanır; politika sonradan değişse de bekleyen talep aynı izni ister. Aynı konu ve tutar için tek bekleyen talep vardır (kısmi benzersiz indeks). Onay tek kullanımlıktır: koşullu güncelleme ile gerçek MongoDB'de 10 eşzamanlı talepten yalnızca biri kazanır.
- **Kurallar:** öz-onay yasak. Checker, talebin snapshot iznine sahip olmalı; `execute` onay yerine geçmez; `modül:*` ve `*:*` kabul edilir. Para birimi politikayla uyuşmazsa onay gerekir (kapalı yönde hata).
- **İade akışı:** eşik üstü ilk istek 409 ile `approval_request_id` döner. Checker onaylar. Aynı idempotency anahtarı ve `approval_request_id` ile yeniden istek, onayı tüketip iadeyi yapar. Onaylı talep ikinci kez kullanılamaz.
- **Uç noktalar:** `GET /admin/approval-gates`, `GET` ve `PUT /admin/approval-policies/:gate`, `GET /admin/approval-requests`, `POST /admin/approval-requests/:id/approve` ve `/reject`.
- **Roller:** kod rol adı bilmez. İngilizce örnek anahtarlar (`institution_manager`, `cashier`, `ticket_agent`, `finance_maker`, `finance_checker`, `payment_releaser`) yalnızca dokümantasyondur; roller admin tarafından oluşturulur.
- **Çeviri:** 11 hata anahtarı ve 1 kapı etiketi; 28 dilde seed (`allLangsMulti`: 9 dil açık metin, 19 dil İngilizce yedek). Kapsam testleri geçti.
- **Seed ve izin kataloğu:** `approval_policies` ve `approval_requests` modülleri izin kataloğunda. Prod için `cmd/seed --roles` gerekir (`SEED_ON_START` kullanılmaz). Politika için seed dosyası yoktur; politikalar admin API ile oluşturulur. Bu, README RULE #3'ten bilinçli bir sapmadır: varsayılanlar kod katalogundadır, veri yöneticiye aittir.

**Açık kalanlar:**
1. **Returns iadesi:** kapıdan geçiyor (`approval_request_id` alanı). Onay bekleyen iade "başarısız" işaretlenmiyor; test ile korunuyor.
2. **Havale onayı ve manuel sefer onayı:** politika motoruna taşınmayacak (karar: kullanıcı onayladı). Her ikisi de zaten iki kişi kuralıyla korunuyor (izin + öz-onay); eşik gerektirmiyorlar. Eşik gerekirse ayrı bir karar olarak yeniden açılır.
3. **B2B sipariş zinciri:** dinamik. Politika (`b2b.purchase_order`, zincir kapısı) sıralı adımları tutar: rol adı, eşik ve eşiğin kaynağı (sabit tutar veya bayi kredi limiti). Sipariş gönderilirken zincir politikadan üretilir. Onaylayıcı rolü oturumdan alınır; gönderen onaylayamaz. Varsayılan zincir seed ile gelir (mağaza yöneticisi > 50.000, bayi kredi limiti, yönetici > 100.000); politika yoksa sipariş gönderilemez (kapalı başarısızlık). Varsayılan rol adları (`store_manager`, `credit_dept`, `executive`) örnektir; rol kurulmadan politika kaydedilemez.
4. **Panel ekranları** (`aystek-core-next`) bu sandbox'ta yok; yapılmadı.
5. **Rol görünen adları:** uygulandı. Rol anahtarı (`name`) İngilizce küçük harf ve değişmez; görünen adlar `labels` haritasında dile göre tutulur (28 dil kodu, en yedeği). Oluşturma sırasında anahtar doğrulanır, sonradan değiştirilemez. Panel rol adlarını bu haritadan göstermeli; panel bu sandbox'ta olmadığı için yapılmadı.
6. **Kapı kataloğu** IBAN, ödeme emri, hakediş ve devir modülleri gelince genişleyecek. Bu modüller henüz yok.
7. **Başarısız iade ve onay:** sağlayıcı iadeyi kabul etmeden önce hata verirse (rezervasyon, rezerve edilemeyen tutar, sağlayıcı hatası) tüketilmiş onay geri açılır ve aynı talep yeniden denenebilir. Sağlayıcı kabulünden sonraki hatalarda onay tüketilmiş kalır; para çıkmıştır ve yeniden kullanılamaz. Test ile korunuyor.

## Ek A — Kanıt İndeksi

Yollar `/workspace/go/aystek-core/` köküne göredir. BAĞIMSIZ satırlar bu oturumda okundu veya grep ile doğrulandı.

| Bulgu | Kanıt | Doğrulama |
|---|---|---|
| RefundEvent `Result` kontrolü yok | `internal/modules/access_points/service/service.go:328-343` | BAĞIMSIZ |
| Yetersiz bakiyede `ChargedAmount` = fee | `access_points/service/service.go:1195-1208`, `1845-1862` | BAĞIMSIZ |
| Cüzdan cache'i koşulsuz `$set` | `wallets/repository/mongo.go:111-127` | BAĞIMSIZ |
| Journal idempotency partial unique | `ledger/repository/indexes.go:18-23` | BAĞIMSIZ |
| `ledger_accounts.code` unique index yok | `ledger/repository/indexes.go` (grep: `code` eşleşmesi yok) | BAĞIMSIZ |
| `JournalVoid` kullanılmıyor | `ledger/model/journal.go:13` | BAĞIMSIZ |
| Cüzdan (person, currency) tekil | `wallets/repository/indexes.go:15-24` | BAĞIMSIZ |
| GİB SOAP hatada mock UUID, nil hata | `billing/invoices/adapters/gib_soap.go:76-86`; `gib_einvoice.go:399` | BAĞIMSIZ |
| `IssueEInvoice` nil sonuç okuması | `billing/invoices/service/service.go:265-271`; `gib_einvoice.go:131-137` | BAĞIMSIZ |
| Refund yalnız status ve tutar kontrolü | `payments/intents/service/service.go:477-531`; `:510` | BAĞIMSIZ |
| Intent durum sabitleri | `payments/intents/model/intent.go:12-18` | BAĞIMSIZ |
| Mock adaptör production koruması | `payments/intents/service/adapters.go:102-105` | BAĞIMSIZ |
| `approvals` main'e bağlı değil | `cmd/api/main.go` (grep: 0 sonuç) | BAĞIMSIZ |
| Production Mongo `--replSet` yok | `deployments/production/docker-compose.yml:11` | BAĞIMSIZ |
| Metered girişte hold yok (yorum) | `access_points/service/service.go:1449-1457` | BAĞIMSIZ |
| Aktif sefer şartlı transport olayı; geç olay reddi | `access_points/service/service_test.go:3143` | DENETÇİ |
| Olay zamanı alanı yok | `access_points/dto/dto.go:50-60`; `access_points/model/access_event.go:156` | DENETÇİ |
| Reddedilen olay kaydedilmiyor | `access_points/model/access_event.go:59` | DENETÇİ |
| Flutter'da yerel kuyruk yok | `aypass_erisim_app/pubspec.yaml:44-66` | DENETÇİ |
| QR opak token, TTL 90 sn, tek kullanım yok | `wallet_instruments/model/qr.go:14-29`; `wallet_instruments/service/service.go:247-297` | DENETÇİ |
| Cihaz yetkilendirme zinciri | `auth/middleware/api_client_gateway.go:183-349` | DENETÇİ |
| Cihaz izin listeleri tutarsız | `auth/handler/invitation.go:192-206`; `auth/handler/admin.go:33-57` | DENETÇİ |
| Cihaz PIN'i 24 bit, public redeem | `auth/handler/invitation.go:54-60, 101` | DENETÇİ |
| Replay: `recordOperation` callback'ten önce | `payments/intents/service/service.go:745-788` | DENETÇİ |
| Invoice/Status modeli | `billing/invoices/model/invoice.go:10-17, 38-46` | DENETÇİ |
| Tarihli sağlayıcı ataması yok | GREP (`ServiceAssignment`, `provider_assignment`, `ownership_history`) | DENETÇİ |
| Hakediş/IBAN/komisyon yok | GREP (`payout`, `settlement`, `hakedi`, `iban`, `commission`) | DENETÇİ |
| Kart UID ile borçlandırma | `wallet_instruments/service/service.go` | DENETÇİ |
| Master plan Increment 148 vs denetçi okuması | `GURSEL_YENI_NESIL_DONUSUM_MASTER_TR.md` Karar #29 / Increment 148; `access_points/service/service.go:940-956` | Çelişki, doğrulanacak |
| TripAssignment `Update` $set onay, devir ve sürücü alanlarını içermiyor | `fleet_tally/trip_assignments/repository/mongo.go:57-69` | BAĞIMSIZ |
| RecordEvent `AccessPoint.Status` denetlemiyor | `access_points/service/service.go` (grep) | BAĞIMSIZ |
| `insufficient_balance` olayında `ChargedAmount` = fee | `access_points/service/service.go:1207-1208, 1845-1862` | BAĞIMSIZ |
| `cmd/migrate_gursel`: sabit bağlantı dizgisi ve git'te izlenen binary (değerler yazılmadı) | `cmd/migrate_gursel/main.go:8, 43`; `cmd/migrate_gursel/migrate_gursel` | BAĞIMSIZ |
| Seed'de düz metin parola ve token alanları (değerler yazılmadı) | `internal/seeds/gursel_turizm_seed.go:391, 406, 422` | BAĞIMSIZ |
| Otomatik kapanış borçlandırması hatayı yutar | `fleet_tally/rider_trips/service/service.go:212-217` | DENETÇİ |
| Rapor toplamları result ve refunded süzmez | `access_points/repository/mongo.go:466-625` | DENETÇİ |
| Stok plaka tarihçesinde aktör istek gövdesinden | `stock/vehicles/dto/dto.go:41`; `stock/vehicles/handler/admin.go:265` | DENETÇİ |
