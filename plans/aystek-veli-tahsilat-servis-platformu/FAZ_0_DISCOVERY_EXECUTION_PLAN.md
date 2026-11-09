# Faz 0 — Keşif ve Mimari Kilitleme: Uygulama Planı

**Tarih:** 2026-05-11  
**Durum:** Uygulamaya hazır  
**Bağımlı Dokümanlar:** `MASTER_WORKPLAN.md`, `ADR_001_DOMAIN_BOUNDARIES_AND_NAMESPACE.md`  
**Amaç:** Sprint 01 (Faz 1) başlamadan önce kritik belirsizlikleri kapatmak, mimari riskleri sınıflandırmak ve kalıcı referans materyaller üretmek.

---

## 1. Neden Faz 0 Gerekli?

Sprint 01, AYSTEK Core içinde `transport_platform` temelini atar. Ancak bazı temel kararlar henüz netleşmemişse sprint ortasında mimari geri dönüşler yaşanır. Faz 0 şu soruları cevaplar:

| Soru | Yanıtlanmadan sprint yapmak riski |
|---|---|
| SAP sürümü ve bağlantı tipi nedir? | `transport_core/sapdispatch` adapter contract yazılamaz |
| Hangi banka POS sağlayıcıları aktif? | `PaymentQuote` BIN routing modeli bilinemez |
| Legacy DB'ye erişim var mı? | `GuardianLink` person lookup adapter geliştirilemez |
| Kan grubu verisi için hukuki dayanak var mı? | `DynamicFieldDefinition` KVKK sınıfı ayarlanamaz |
| Mobil uygulama mimarisi ne? | Public endpoint tasarımı bilinemez |
| TOSLA entegrasyon modeli nedir? | `TopUpProviderBinding` contract yazılamaz |
| Kardeş indirimi için belge kanıtı şart mı? | `DiscountRule.family_group` modeli bilinemez |

---

## 2. Faz 0 Kapsam ve Çıktılar

### 2.1 Üretilecek Dokümanlar

| Doküman | Konum | Sorumlu konu |
|---|---|---|
| `DISCOVERY_FINDINGS.md` | Bu klasör | Genel bulgular ve sınıflandırılmış belirsizlikler |
| `LEGACY_MAPPING_MATRIX.md` | Bu klasör | PHP/Yii2 modelleri → yeni domain modeli eşleme tablosu |
| `BANK_AND_POS_DECISION_RECORD.md` | Bu klasör | Aktif banka POS sağlayıcıları, MID/TID yönetimi, BIN fallback kararı |
| `SAP_ADAPTER_DECISION_RECORD.md` | Bu klasör | SAP sürümü, bağlantı tipi, dispatch payload şeması |
| `KVKK_SPECIAL_CATEGORY_DATA_POLICY.md` | Bu klasör | Kan grubu ve diğer özel nitelikli verilerin hukuki dayanak ve saklama politikası |

### 2.2 Keşif Soruları ve Belirsizlik Sınıflandırması

Her belirsizlik şu sınıflardan birine alınır:
- **BLOCKER:** Cevap alınmadan sprint başlatılamaz.
- **ASSUMPTION:** Sprint ilerleyebilir; varsayım makul, sonradan değiştirilebilir.
- **ACCEPTED RISK:** Bilinçli olarak ertelendi; riski azaltmak için önlem alındı.

---

## 3. Keşif Soruları

### 3.1 SAP Entegrasyonu (BLOCKER)

- [ ] Gürsel Turizm'in kullandığı SAP sürümü nedir? (ECC, S/4HANA, Business One)
- [ ] Entegrasyon tipi nedir? OData, IDoc, BAPI/RFC, SAP PI/PO, SAP Integration Suite?
- [ ] Ödeme/tahsilat kaydı için hangi SAP transaction tipi kullanılıyor? (örn. FI-AR incoming payment)
- [ ] SAP'ye gönderilecek zorunlu alanlar neler? (örn. müşteri kodu, referans no, tutar, para birimi)
- [ ] SAP test ortamı var mı? Erişim bilgileri nasıl paylaşılacak?
- [ ] Mevcut GurselVB22 uygulaması SAP'ye nasıl bağlanıyor? Hangi API/RFC kullanılıyor?

**Hedef çıktı:** `SAP_ADAPTER_DECISION_RECORD.md` içinde dispatch payload şeması ve adapter interface contract taslağı.

### 3.2 Banka POS Sağlayıcıları (BLOCKER)

- [ ] Aktif sanal POS sağlayıcılar hangileri? (Garanti, İşbank, YapıKredi, Finans, Akbank — hepsi aktif mi?)
- [ ] Her banka için MID (Merchant ID) ve TID (Terminal ID) yönetimi nasıl yapılıyor?
- [ ] BIN veritabanı: harici servis mi (örn. Interpos, Mastercard/Visa BIN API), kendi tablosu mu?
- [ ] Kart bankası bilinmiyorsa (BIN match yok) fallback POS hangisi?
- [ ] Taksit tabloları: banka API'sinden gerçek zamanlı mı alınıyor, yoksa admin panelinden mi giriliyor?
- [ ] 3DS callback URL'i sabit mi, yoksa dynamic mi üretiliyor?
- [ ] TOSLA entegrasyonu: TOSLA'nın kendi API'si mi var? Merchant agreement var mı?
- [ ] TOSLA üzerinden yapılan yüklemelerin faturası kime kesiliyor — TOSLA mı, kullanıcı mı?

**Hedef çıktı:** `BANK_AND_POS_DECISION_RECORD.md` içinde BIN routing stratejisi ve her sağlayıcı için adapter interface contract taslağı.

### 3.3 Legacy Veritabanı ve Kişi Lookup (BLOCKER)

- [ ] GurselVB22 hangi veritabanına bağlanıyor? (MySQL/MariaDB tablo yapısı?)
- [ ] Kişi araması için hangi alanlar kullanılıyor? (TC + okul kodu yeterli mi, ek validasyon var mı?)
- [ ] Yeni sistemde legacy DB'ye okuma bağlantısı verilecek mi? (read-only bridge)
- [ ] Legacy'den yeni sisteme veri migration planlanıyor mu? Hangi veriler taşınacak?
- [ ] Borç verileri SAP'ten mi geliyor, yoksa legacy DB'de mi?
- [ ] Çocuk/personel kaydında TC ile birlikte başka kimlik alanı var mı? (örn. okul numarası, sicil no)

**Hedef çıktı:** `LEGACY_MAPPING_MATRIX.md` içinde legacy alan → yeni domain model eşleme tablosu.

### 3.4 KVKK ve Özel Nitelikli Veri (BLOCKER)

- [ ] Kan grubu toplanması için hukuki dayanak nedir? (açık rıza mı, kanuni yükümlülük mü, hayati menfaat mi?)
- [ ] Kan grubu verisi için ayrı bir aydınlatma metni hazırlanacak mı?
- [ ] Verilerin tutulma süresi, silinme politikası ve erişim günlüğü için kurumun politikası var mı?
- [ ] Kan grubu dışında başka özel nitelikli veri toplanacak mı? (sağlık verisi, engellilik durumu vb.)
- [ ] KVKK VERBİS kaydı var mı? Yeni veri işleme aktivitesi eklenmesi gerekiyor mu?

**Hedef çıktı:** `KVKK_SPECIAL_CATEGORY_DATA_POLICY.md` içinde veri sınıfı, hukuki dayanak, saklama ve erişim politikası.

### 3.5 Kardeş İndirimi ve İndirim Politikaları (ASSUMPTION geçici)

- [ ] Kardeş indirimi için belge kanıtı (dosya yükleme) zorunlu mu, yoksa sistem otomatik tespit edebilir mi?
- [ ] Kardeş indirimi: yüzde mi, sabit tutar mı, taksit adedi mi?
- [ ] Birden fazla kardeş varsa ek indirim var mı?
- [ ] Kardeş indirimi kurum bazlı mı, genel mi?
- [ ] Diğer indirim politikaları: peşin ödeme indirimi, erken kayıt indirimi yüzdeleri neler?

**Varsayım (geçici):** Kardeş indirimi, aynı `guardian_link`'e bağlı birden fazla aktif `GuardianLink`'in aynı kurumda olması koşuluyla `rule_type: family_group` olarak modellenir; belge yükleme ikinci fazda eklenir.

### 3.6 Mobil Uygulama ve Rol-bazlı Kanallar (ASSUMPTION geçici)

- [ ] Şoför ekranı, taşınan kişi uygulaması, araç sahibi ekranı — bunlar yeni uygulama mı, yoksa mevcut uygulamaların AYSTEK Core API'ye bağlanması mı?
- [ ] Mevcut `/workspace/android/GurselBakiyeCekme*` projeleri bu kanallar için mi? Kullanılacak mı?
- [ ] `aysteknoloji-com` public sitesi bu kanallara ev sahipliği yapacak mı?
- [ ] Mifare kart okuyucusu: hangi donanım? AYLAID, ACR, diğer? SDK/protokol var mı?
- [ ] Offline mod için cihaz platform hedefi nedir? Android? Raspberry Pi? Dedicated terminal?

**Varsayım (geçici):** Sprint 01-06 boyunca API endpoint tasarımı role-based public endpoints içerir; client uygulaması ayrı scope olarak planlanır.

### 3.7 Servis Yaklaşma Bildirimi (ACCEPTED RISK)

- [ ] Yaklaşma eşiği ne olacak? Mesafe mi (örn. 500 metre), tahmini varış süresi mi?
- [ ] Bildirim kanalı: push notification, SMS, ya da her ikisi?
- [ ] Opt-in zorunlu mu? Veri saklama politikası?
- [ ] GPS veri gecikmesi ne kadar? (N2Mobil/WeTrack polling frekansı?)

**Kabul edilen risk:** Bu özellik Faz 6-7'ye ertelendi. Telematics ve notifications altyapısı önceki fazlarda kurulacak; yaklaşma bildirimi bu omurga üstüne Faz 6'da eklenecek.

---

## 4. Legacy Sistem Analizi

### 4.1 GurselVB22 (Veli Borç Takip)

Analiz edilecek dosyalar:

```
/workspace/php/GurselVB22/
├── controllers/
│   ├── SiteController.php     → Tüm banka POS akışları (Garanti, İşbank, YapıKredi, Finans, Akbank)
│   └── HavaleController.php   → Havale bildirimi akışı
├── models/
│   ├── OgrenciBilgileri.php   → Öğrenci/kişi veri yapısı
│   ├── OgrenciSMSOnay.php     → SMS OTP akışı
│   ├── ServisOnKayitForm.php  → Ön kayıt formu alanları
│   ├── ServistenAyrilmaTalebiForm.php → Ayrılma talebi
│   ├── VSozlesmesi.php        → Sözleşme şablon değişkenleri
│   ├── KardesIndirimi.php     → Kardeş indirimi kuralı
│   └── HavaleForm.php         → Havale form alanları
└── views/                     → Form alanlarından yeni UI gereksinimleri çıkarılabilir
```

**Görev:** Her model için alan listesi çıkar → `LEGACY_MAPPING_MATRIX.md`'e işle.

### 4.2 AYShuttleGursel (Araç Takip / KT)

Analiz edilecek dosyalar:

```
/workspace/php/AYShuttleGursel/
├── controllers/
│   └── SiteController.php     → Müşteri giriş akışları
├── models/
│   ├── Kullanicilar.php       → KT kullanıcı modeli
│   ├── KullaniciDurumlari.php → Kullanıcı durum geçişleri
│   └── SifreSureleri.php      → Şifre/erişim süre politikası
```

**Görev:** KT'de kullanılan kullanıcı durumları ve müşteri giriş alanları → `transport_access` ve `transport_core` person modeline mapping.

---

## 5. Belirsizlik Sınıflandırma Tablosu

Bu tablo Faz 0 sonunda doldurulmalı ve `DISCOVERY_FINDINGS.md` içinde publish edilmeli:

| Belirsizlik | Sınıf | Yanıt | Etkilenen sprint |
|---|---|---|---|
| SAP sürümü ve bağlantı tipi | BLOCKER | — | Faz 6 |
| Aktif banka POS sağlayıcı listesi | BLOCKER | — | Faz 5 |
| BIN veritabanı kaynağı | BLOCKER | — | Faz 5 |
| Legacy DB okuma erişimi | BLOCKER | — | Faz 2 |
| KVKK kan grubu hukuki dayanak | BLOCKER | — | Faz 3 |
| TOSLA entegrasyon modeli | BLOCKER | — | Faz 5 |
| Kardeş indirimi belge zorunluluğu | ASSUMPTION | family_group rule, belge Faz 2+ | Faz 2 |
| Mobil uygulama platform kararı | ASSUMPTION | API-first, client ayrı scope | Faz 1+ |
| Servis yaklaşma bildirim eşiği | ACCEPTED RISK | Faz 6'ya ertelendi | Faz 6 |
| Offline cihaz donanım platformu | ASSUMPTION | ADR_002'ye göre Android öncelikli | Faz 4+ |

---

## 6. Teknik Hazırlık Kontrol Listesi

Faz 0 bitiminde bu kontrol listesi tamamlanmış olmalıdır:

### 6.1 Dokümanlar

- [ ] `DISCOVERY_FINDINGS.md` oluşturuldu; tüm belirsizlikler BLOCKER/ASSUMPTION/ACCEPTED RISK olarak sınıflandırıldı
- [ ] `LEGACY_MAPPING_MATRIX.md` oluşturuldu; GurselVB22 + AYShuttleGursel model → domain model tabloları dolduruldu
- [ ] `BANK_AND_POS_DECISION_RECORD.md` oluşturuldu; aktif POS sağlayıcıları ve BIN routing kararı netleştirildi
- [ ] `SAP_ADAPTER_DECISION_RECORD.md` oluşturuldu; SAP sürümü ve adapter interface taslağı hazır
- [ ] `KVKK_SPECIAL_CATEGORY_DATA_POLICY.md` oluşturuldu; kan grubu ve özel veri politikası onaylı

### 6.2 Mimari Kilitleme

- [ ] `ADR_001` SAP dispatch namespace kararı güncellendi (`transport_core/sapdispatch`) ✅
- [ ] `ADR_002` Offline-first sync protokolü onaylandı ✅
- [ ] `ADR_003` Multi-banka POS routing kararı onaylandı ✅
- [ ] Sprint 01 (Faz 1) başlatmak için BLOCKER listesi boş

### 6.3 Erişimler

- [ ] AYSTEK Core dev ortamı erişimi aktif
- [ ] Legacy DB'ye (GurselVB22) salt okunur erişim sağlandı veya model dump alındı
- [ ] SAP test ortamı erişimi veya sandbox credentials planlandı

---

## 7. Faz 0 Definition of Done

Faz 0 tamamlandı sayılmak için:

- Tüm 5 keşif dokümanı oluşturuldu ve içerik onaylandı.
- BLOCKER listesindeki tüm maddeler yanıtlanmış ya da "ACCEPTED RISK" olarak kararlaştırılmış.
- `ADR_002` ve `ADR_003` onaylanmış.
- Sprint 01 başlatma kararı verilmiş.
- Bu klasörde `FAZ_0_COMPLETION_CHECKPOINT.md` oluşturulmuş.

---

## 8. Süre Tahmini

Faz 0 için yeterli bilgiyle (legacy DB erişimi + SAP bilgisi + banka POS bilgisi) 3-5 iş günü beklenmektedir. BLOCKER'lar yanıtlanmadan sprint başlatılmamalıdır.

Kritik yol: **SAP bilgisi** ve **legacy DB erişimi** en uzun bekleme sürelerine neden olabilir. Bu iki konuyu paralel yürütmek Faz 0'ı kısaltır.
