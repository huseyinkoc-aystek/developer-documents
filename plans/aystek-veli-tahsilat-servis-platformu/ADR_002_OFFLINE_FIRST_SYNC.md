# ADR 002 — Offline-First Geçiş Senkronizasyon Protokolü

**Tarih:** 2026-05-11  
**Durum:** Kabul edildi  
**Karar türü:** Kritik altyapı / platform protokolü  
**Kapsam:** `transport_access` domain — cihaz tabanlı geçiş işleme, offline kayıt ve sunucu senkronizasyonu  
**İlgili ADR:** `ADR_001_DOMAIN_BOUNDARIES_AND_NAMESPACE.md`

---

## 1. Bağlam

Kullanıcı anlatımında "tüm bu uygulamalar için en önemlisi; sunucunun durması veya internet sorunlarının yaşandığı acil durumlarda, uygulamaların çevrimdışı olarak aksamadan çalışabilmesidir" şeklinde tanımlanmış; bu gereksinim Non-Negotiable Standartlar arasındadır.

Servis araçlarındaki okuyucu cihazlar şu ortamlarda çalışır:

- Tünel içi, kampüs içi, kötü hücre kapsaması olan bölgeler.
- Sabah 07:00 servis yoğunluğu gibi anlık trafik baskısı.
- Sunucu planlı veya plansız bakım pencereleri.

Bu koşullarda okuyucunun "internet yok → geçiş reddedildi" davranışı kabul edilemez. Aynı zamanda offline kaldığı süre boyunca sınırsız ve kontrolsüz geçişe izin vermek finansal risk üretir.

---

## 2. Karar

`transport_access/offline_sync` paketi, cihaz bazlı append-only transaction log ve senkronizasyon protokolü üzerine inşa edilecektir.

### Temel İlkeler

1. **Cihaz yetkili kayıt olarak hareket eder.** Internet kesintisinde cihaz, son policy snapshot'ına göre kararını kendisi verir.
2. **Append-only log.** Offline transaction'lar silinmez, sadece eklenir. Cihaz imzası taşır.
3. **Idempotent sync.** Aynı batch iki kez gönderilse sunucu duplicate oluşturmaz.
4. **Sınırlı offline özgürlük.** Risk limitleri (max işlem, max tutar, max süre) yapılandırılabilir; aşılırsa cihaz okumayı durdurur.
5. **Conflict resolution deterministik.** Çakışma çözümü kuralları önceden tanımlıdır; operatör inisiyatifine bırakılmaz.
6. **Fraud audit.** Offline batch sunucuya ulaştığında otomatik fraud kontrol yapılır.

---

## 3. Cihaz Policy Snapshot

Cihaz online olduğunda sunucudan şu snapshot'ı alır ve yerel şifreli depoda saklar:

```json
{
  "snapshot_id": "snap_abc123",
  "device_id": "dev_001",
  "tenant_id": "gursel-turizm",
  "site_id": "aractakip",
  "valid_until": "2026-05-11T06:00:00Z",
  "max_offline_minutes": 240,
  "max_offline_transactions": 200,
  "max_offline_total_amount": 5000.00,
  "currency": "TRY",
  "fare_rules": [...],
  "group_discounts": [...],
  "credit_policies": [...],
  "refund_window_minutes": 15,
  "route_context": {...},
  "hmac_signature": "sha256:..."
}
```

**Kurallar:**
- `valid_until` geçtikten sonra cihaz yalnızca `credit_policies` içindeki kullanıcılar için offline geçişe izin verir.
- `max_offline_minutes` dolduğunda cihaz uyarı üretir; operatör kararı olmadan devam edemez.
- Snapshot HMAC ile imzalanır; cihaz tahrif tespit eder.

---

## 4. Offline Transaction Log Formatı

Her offline transaction şu formatta cihazın yerel append-only log dosyasına yazılır:

```json
{
  "device_seq": 1042,
  "device_id": "dev_001",
  "snapshot_id": "snap_abc123",
  "occurred_at": "2026-05-11T08:22:15.123Z",
  "instrument_token": "tok_abc...xyz",
  "instrument_type": "mifare_uid_hash | qr_token | mobile_token",
  "transaction_type": "boarding | alighting | refund | guest_ticket",
  "route_id": "route_001",
  "direction_id": "dir_sabanci_to_kadikoy",
  "resolved_fare_amount": 12.50,
  "resolved_discount_pct": 0.0,
  "resolved_credit_used": false,
  "wallet_balance_before": 87.50,
  "wallet_balance_after": 75.00,
  "decision_source": "snapshot | cached_policy | emergency_credit",
  "offline_mode": true,
  "device_signature": "hmac:sha256:..."
}
```

**Kurallar:**
- `device_seq` cihaz bazında monoton artar; sunucu sıra boşluğu tespit eder.
- `device_signature` snapshot imzası ile türetilmiş cihaz anahtarıyla imzalanır; sunucu doğrular.
- `instrument_token` ham UID değil, cihazda türetilmiş token'dır (PAN/CVV analogu ham UID sunucuya gönderilmez).

---

## 5. Offline Sync Protokolü

### 5.1 Batch Yapısı

```json
{
  "batch_id": "batch_dev001_1042",
  "device_id": "dev_001",
  "snapshot_id": "snap_abc123",
  "batch_created_at": "2026-05-11T08:45:00Z",
  "transactions": [...],
  "device_signature": "hmac:sha256:..."
}
```

### 5.2 Sunucu Tarafı İşlem Sırası

1. Batch alındı → `transport_offline_batches` koleksiyonuna `received` statüsünde yazıldı.
2. Signature doğrulaması: başarısız → `signature_invalid` statüsü; batch reddedildi; operasyon alarmı üretildi.
3. Snapshot geçerlilik kontrolü: `valid_until` vs. `occurred_at`.
4. `device_seq` sıra kontrolü: boşluk varsa → `gap_detected` uyarısı; işlem durdurulmaz ama audit kaydı yapılır.
5. İdempotency check: `batch_id` daha önce işlenmiş mi?
6. Transaction-by-transaction uygulama:
   - Wallet bakiyesi güncelleme (gerçek uygulama).
   - Fare policy karşılaştırma: offline karar vs. sunucu kararı. Fark varsa → `fare_discrepancy` kaydı.
   - İade penceresi doldurulmuş mu kontrolü.
7. `MovementReport` güncelleme: batch transaction'ları ilgili `Trip`'e eklendi.
8. Fraud analizi: `offline_fraud_review` event'i yayınlandı.
9. Batch statüsü: `applied` veya `applied_with_discrepancies`.

### 5.3 Conflict Resolution Kuralları

| Çakışma türü | Çözüm kuralı |
|---|---|
| Aynı instrument 15 dk içinde iki kez boarding | İkincisi `duplicate_prevented` olarak işaretlenir; ücret alınmaz |
| Offline bakiye yetersiz ama credit policy vardı | `credit_applied_offline` olarak kaydedilir; gerçek ledger'a düşer |
| Offline fare !== sunucu fare | `fare_discrepancy` kaydı; fark ledger'a reconciliation queue'ya gider |
| Snapshot süresi dolmuş sırada geçiş | `stale_snapshot` flag ile kaydedilir; audit review kuyruğuna alınır |
| Offline boarding ama online alighting | Online alighting geçerlidir; boarding `reconstructed` olarak işaretlenir |
| Sefer sonrası gelen batch (trip kapanmış) | Batch uygulanır; trip `retroactively_updated` olarak işaretlenir |

---

## 6. Risk Limitleri

Cihaz yapılandırmasına göre (snapshot'ta gelir):

| Limit | Varsayılan değer | Aşıldığında |
|---|---|---|
| `max_offline_minutes` | 240 dk (4 saat) | Cihaz sesli/görsel uyarı verir; operatör override gerektir |
| `max_offline_transactions` | 200 işlem | Cihaz yeni geçişleri kabul etmez |
| `max_offline_total_amount` | 5.000 TRY | Cihaz yeni geçişleri kabul etmez |
| `max_stale_snapshot_minutes` | 480 dk (8 saat) | Snapshot kabul edilmez; sadece emergency_credit kullanıcıları geçer |

Bu değerler `transport_access_devices` koleksiyonundaki cihaz kaydından veya `transport_access_rules` tenant-level policy'den gelir.

---

## 7. Offline'da Yapılmayacaklar

| Yasak eylem | Gerekçe |
|---|---|
| Yeni abonelik aktivasyonu | Finansal risk; backend onayı gerekir |
| Yüksek riskli misafir bilet satışı (online ödeme gerektirir) | Ödeme tamamlanmadan geçiş imkânı |
| Limitsiz eksi bakiye açma | Sonsuz kredi riski |
| PAN/CVV veya ham Mifare UID'nin sunucuya gönderilmesi | Güvenlik; token kullanılır |
| Eski policy snapshot ile sınırsız çalışma | Fiyat manipülasyonu riski |
| Havale veya banka POS işlemi | İnternet bağlantısı zorunludur |

---

## 8. Cihaz Güvenlik Modeli

- Her cihaz unique key pair ile sağlanır; provision sırasında sunucu sertifika authority'sine kaydedilir.
- Snapshot ve transaction imzaları bu key ile türetilir.
- Cihaz revoke edilebilir: `transport_access_devices.status = revoked`; revoke edilen cihazın batch'leri reddedilir.
- Cihaz kayıp/çalıntı durumunda remote wipe ve key invalidation süreci tanımlanmalıdır (Faz 6 scope'u).

---

## 9. Frontend ve Operasyon Görünümü

**Admin panel:**
- `transport_access/offline_sync` batch listesi: `device_id`, `batch_id`, `received_at`, `status`, `discrepancy_count`.
- Discrepancy detay ekranı: her transaction için offline vs. server karar karşılaştırması.
- Stale snapshot uyarı dashboard'u.

**Event'ler:**

| Event key | Tüketici |
|---|---|
| `transport.offline_batch.received` | ops dashboard, fraud review |
| `transport.offline_batch.applied` | trip ledger, movement report |
| `transport.offline_batch.signature_invalid` | security alert |
| `transport.offline_batch.gap_detected` | ops review |
| `transport.offline_fare.discrepancy` | reconciliation queue |

---

## 10. Sonuç

Offline-first geçiş sistemi, cihazda policy snapshot + append-only transaction log + idempotent batch sync üçgeni üzerine inşa edilir. Risk limitleri yapılandırılabilir; conflict resolution deterministik kurallara göre otomatiktir. Sunucu hiçbir zaman "primary" yerine "authoritative reconciler" konumundadır: cihaz karar verir, sunucu doğrular ve muhasebeleştirir.

Bu tasarım, internet kesintisini "graceful degradation" ile karşılar ve finansal tutarlılığı korur.
