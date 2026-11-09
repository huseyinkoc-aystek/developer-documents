# ADR 003 — Çoklu Banka POS Routing ve BIN Tabanlı Acquirer Seçimi

**Tarih:** 2026-05-11  
**Durum:** Kabul edildi (banka profil detayları Faz 0 keşfinde tamamlanacak)  
**Karar türü:** Ödeme altyapısı / acquirer routing  
**Kapsam:** `transport_portal/payments` ve `payments/gateways` modülü — BIN yönlendirme ve taksit sunumu  
**İlgili ADR:** `ADR_001_DOMAIN_BOUNDARIES_AND_NAMESPACE.md`  
**Bağımlı keşif:** `BANK_AND_POS_DECISION_RECORD.md` (Faz 0'da oluşturulacak)

---

## 1. Bağlam

Legacy `GurselVB22/controllers/SiteController.php` içinde her banka için ayrı action metodları vardır:

```
garantitaksit, garantiodeme3d, garantiodeme3donay
yapikreditaksit, yapikrediodeme3d, yapikrediodeme3donay
isbanktaksit, isbankodeme3d, isbankodeme3donay
finanstaksit, finansyenitaksit, finansodeme3d, finansyeniodeme3d, finansodeme3donay, finansyeniodeme3donay
akbanktaksit, akbankodeme3d, akbankodeme3donay
```

Bu hardcoded yaklaşım:
- Her yeni banka için yeni kod yazılmasını zorunlu kılıyor.
- Taksit tablolarının bakımını zorlaştırıyor.
- BIN eşleşmesi olmayan kartlar için kullanıcı manuel banka seçmek zorunda.
- Çoklu MID/TID yönetimini imkânsız kılıyor.

Yeni sistemde kullanıcı anlatımı: "Kartların numaralarına göre hangi bankaya aitse ya da bizde sanal posu olmayan kart ise ona göre sanalpos alt yapısı seçilecek. Taksitler buna göre iletilecek."

---

## 2. Karar

BIN/IIN tabanlı otomatik acquirer routing yapılacaktır. Kart numarasının ilk 6-8 hanesi (BIN/IIN) tanınırsa eşleşen acquirer profiline yönlendirilir. Tanınamazsa yapılandırılmış fallback acquirer kullanılır.

---

## 3. Mimari Bileşenler

### 3.1 AcquirerProfile (PaymentGateway kayıtları)

Mevcut `payments/gateways` modülü extend edilir. Her acquirer profili:

```json
{
  "gateway_id": "garanti_pos_main",
  "provider": "garanti",
  "tenant_id": "gursel-turizm",
  "merchant_id": "...",
  "terminal_id": "...",
  "priority": 10,
  "supported_card_types": ["visa", "mastercard", "troy"],
  "max_installments": 12,
  "currency": "TRY",
  "is_fallback": false,
  "status": "active"
}
```

Birden fazla sağlayıcı aynı tenant altında aktif olabilir. `priority` alanı sıralamayı belirler.

### 3.2 BINRecord

BIN veritabanı `payments/gateways` veya ayrı `transport_portal/payments` içinde yönetilebilir. Yapı:

```json
{
  "bin_prefix": "454965",
  "card_brand": "visa",
  "issuing_bank_code": "garanti",
  "issuing_bank_name": "T. Garanti Bankası A.Ş.",
  "card_type": "credit",
  "country_code": "TR",
  "preferred_gateway_id": "garanti_pos_main"
}
```

**BIN veritabanı kaynağı (Faz 0'da kararlaştırılacak seçenekler):**

| Seçenek | Avantaj | Dezavantaj |
|---|---|---|
| Yerel admin-managed tablo | Kontrol, no API bağımlılığı | Manuel bakım gerektir |
| Harici BIN lookup API (Mastercard/Visa, Interpos) | Otomatik güncel | API maliyeti, latency |
| Hibrit: yerel cache + harici fallback | Hız + güncellik | Kompleks senkronizasyon |

**Varsayılan karar (Faz 0'a kadar):** Hibrit yaklaşım. Yönetim panelinden düzenlenebilen yerel BIN kaydı + tanınmayan BIN için harici lookup fallback. Harici API yoksa yerel tablo yeterli.

### 3.3 AcquirerRoutingPolicy

Routing kuralları öncelik sırasıyla:

1. BIN → `preferred_gateway_id` eşleşmesi (en yüksek öncelik).
2. BIN card brand + tenant-level `acquirer_priority` sıralaması.
3. Fallback: `is_fallback: true` olan aktif gateway.
4. Hiçbiri yoksa: `PaymentQuote` oluşturulamaz; kullanıcıya hata mesajı.

---

## 4. İki Aşamalı Quote Modeli

### 4.1 Provisional Quote (Kart Bilgisi Girmeden Önce)

Kullanıcı kart girmeden önce:

```
POST /public/transport-platform/payment-quotes
{
  "debt_account_id": "...",
  "debt_item_ids": ["..."],
  "payment_method_hint": "card"
}
```

Sunucu yanıtı: Tüm aktif acquirer'lerin taksit seçeneklerinin birleşimi (worst-case estimation). Kullanıcı genel olarak "12 taksit var" bilgisini görür.

```json
{
  "quote_id": "q_001",
  "status": "provisional",
  "total_amount": 1200.00,
  "currency": "TRY",
  "available_installments": [1, 3, 6, 9, 12],
  "installment_details": [
    {"count": 1, "per_installment": 1200.00, "total": 1200.00},
    {"count": 3, "per_installment": 407.40, "total": 1222.20}
  ],
  "note": "provisional_all_acquirers"
}
```

### 4.2 Final Quote (BIN Girildikten Sonra)

Kullanıcı kart numarasının ilk 6-8 hanesini girdiğinde:

```
POST /public/transport-platform/payment-quotes/{id}/finalize
{
  "bin_prefix": "454965"
}
```

Sunucu:
1. BIN lookup → issuing_bank → `preferred_gateway_id`.
2. O gateway'e özel taksit tablosunu çeker.
3. Quote güncellenir: `status: "final"`, `acquirer_id` set edilir.
4. Artık taksitler o bankanın gerçek taksit koşullarına göre listelenmiştir.

```json
{
  "quote_id": "q_001",
  "status": "final",
  "acquirer_id": "garanti_pos_main",
  "issuing_bank": "Garanti",
  "card_brand": "visa",
  "available_installments": [1, 3, 6, 9, 12],
  "installment_details": [...]
}
```

---

## 5. Taksit Tablosu Yönetimi

**Seçenek A — Admin panelinden manuel giriş:** Her gateway için max taksit, ek maliyet oranları admin'den girilir. Basit; gerçek zamanlı değildir.

**Seçenek B — Banka API'sinden gerçek zamanlı çekme:** Banka API'si destekliyorsa (örn. Garanti/İşbank bazı sağlayıcılar sunuyor). Kompleks; her banka için ayrı integration.

**Varsayılan karar:** Seçenek A ile başlanır. Admin panelinde `PaymentGatewayInstallmentTable` yönetimi yapılır. Faz 5'te gerçek zamanlı API entegrasyonu eklenebilir.

---

## 6. Kart Bankasının Bilinmediği Durum (No BIN Match)

BIN veritabanında eşleşme yoksa:

1. Fallback acquirer (`is_fallback: true`) seçilir.
2. Kullanıcıya "Kartınızın bankası tespit edilemedi, standart taksitler uygulanacak" mesajı gösterilir (i18n key: `transport_portal.payment.bin_unknown_fallback`).
3. Fallback acquirer'ın taksit tablosu gösterilir.
4. Eğer fallback acquirer yoksa kullanıcı "Lütfen banka transferi ile ödeme yapın" yönlendirmesine alınır.

---

## 7. TOSLA / Üçüncü Taraf Top-up Sağlayıcısı

TOSLA, kart POS'undan farklıdır: cüzdan bakiyesi yükler, doğrudan borç ödemesi yapmaz.

Ayrım:
- `transport_wallet/topups` — TOSLA gibi sağlayıcılarla cüzdan yükleme.
- `transport_portal/payments` — Borç ödeme için banka POS 3DS.

TOSLA yükleme fatura sorumluluğu: `TopUpProviderBinding.invoice_to: "provider" | "user"` alanı. Eğer `provider` seçilmişse yükleme SAP'ye TOSLA adına dispatch edilir; kullanıcı faturası oluşmaz.

---

## 8. Güvenlik Gereksinimleri

- PAN/CVV AYSTEK Core'a girmez. Kart verisi doğrudan acquirer'ın hosted component veya JS widget'ına girer.
- BIN prefix (6-8 hane) müşteri tarafından geliyor ve düşük riskli veridir; ancak rate limiting ile kötüye kullanım önlenir.
- `PaymentQuote` token tabanlıdır; quote ID ile başka kullanıcının ödeme akışına geçilemez.
- 3DS callback token tek kullanımlık ve TTL'lidir.
- Tüm payment attempt'ler `transport_payment_attempts` audit log'una düşer.

---

## 9. Faz 0 Kararlarına Bağımlılık

Bu ADR, `BANK_AND_POS_DECISION_RECORD.md` tamamlandığında şu alanlar doldurulacaktır:

| Alan | Faz 0'da netleşecek |
|---|---|
| Aktif banka listesi | Hangi gateway profilleri seed edilecek |
| BIN veritabanı kaynağı | Yerel tablo mı, harici API mi, hibrit mi |
| Taksit tablosu güncelleme frekansı | Manuel admin mi, API mi |
| TOSLA merchant agreement | TOSLA top-up endpoint ve invoice model |
| Fallback acquirer kararı | Hangi banka fallback olacak |

---

## 10. Sonuç

Çoklu banka POS routing, BIN/IIN tabanlı otomatik acquirer seçimi ile yönetilir. İki aşamalı quote modeli (provisional → final) kullanıcıya önce genel taksit bilgisi verir, kart BIN'i girilince kesin acquirer ve taksit koşulları gösterilir. Tüm acquirer profilleri admin panelinden yönetilir; kod değişikliği gerekmez.
