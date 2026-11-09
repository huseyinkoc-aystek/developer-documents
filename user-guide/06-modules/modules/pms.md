# Mülk Yönetim Sistemi — pms

**Sürüm:** 1.0.0 | **Son Güncelleme:** 2026-03-31

## Genel Bakış

`pms` (Property Management System — Mülk Yönetim Sistemi), site platformuna entegre edilmiş kapsamlı bir mülk, kiracı ve bakım yönetimi modülüdür. Apartman siteleri, konut kompleksleri ve ticari mülk yöneticileri için tasarlanmıştır.

### Temel Yetenekler

| Alan | Açıklama |
|------|----------|
| **Mülk Yönetimi** | Bağımsız birimler, kat mülkiyeti, mülk portföyü |
| **Kiracı / Sakin Takibi** | Sözleşme, kira/aidat ödemeleri, sorumluluk |
| **Aidat Planları** | Dönemsel aidat planı oluşturma ve aktivasyon |
| **İş Emirleri** | Bakım talebi, atama, takip, harcama |
| **Tedarikçiler** | Bakım firmaları, hizmet sözleşmeleri |
| **Ölçüm / Sayaç** | Su, elektrik, doğalgaz sayaç okuma |
| **Mutabakat** | Ödeme mutabakat işlemleri |
| **Sakin Portalı** | Sakinlerin kendi bilgilerini ve işlemlerini görebildiği self-servis portal |
| **Webhook** | Dış sistemlere olay bildirimi |
| **Dışa Aktarım** | Excel / PDF rapor çıktıları |

---

## Mimari

```

                     PMS Modülü                          │

  Properties  │  Occupancies │   Billing    │  Metering  │
  (Mülkler)   │  (Kiracı)    │  (Aidat)     │  (Sayaç)   │
go/aystek-core cd
  Work Orders │  Vendors     │ Reconcil.    │  Reports   │
  (İş Emri)   │  (Tedarikçi) │  (Mutabakat) │  (Rapor)   │

  Dashboard | Portfolio | Legal | Audit | Exports        │

              Portal (Sakin Self-Servis)                  │

```

---

## İzinler

| İzin | Açıklama |
|------|----------|
| `pms_properties:read` | Mülk görüntüleme |
| `pms_properties:create` | Mülk oluşturma |
| `pms_properties:update` | Mülk düzenleme |
| `pms_properties:delete` | Mülk silme |
| `pms_occupancies:read` | Kiracı / sakin görüntüleme |
| `pms_occupancies:create` | Kiracı oluşturma |
| `pms_occupancies:update` | Kiracı güncelleme |
| `pms_occupancies:delete` | Kiracı silme |
| `pms_billing_plans:read` | Aidat planı görüntüleme |
| `pms_billing_plans:create` | Aidat planı oluşturma |
| `pms_billing_plans:update` | Aidat planı güncelleme |
| `pms_billing_plans:activate` | Aidat planı aktivasyonu |
| `pms_billing_plans:archive` | Aidat planı arşivleme |
| `pms_dashboard:read` | PMS kontrol paneli görüntüleme |

---

## API Endpoint'leri

### Mülkler (`/pms/properties`)

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/pms/properties` | Mülk listesi | `pms_properties:read` |
| `GET` | `/pms/properties/:id` | Mülk detayı | `pms_properties:read` |
| `POST` | `/pms/properties` | Yeni mülk oluştur | `pms_properties:create` |
| `PUT` | `/pms/properties/:id` | Mülk güncelle | `pms_properties:update` |
| `DELETE` | `/pms/properties/:id` | Mülk sil | `pms_properties:delete` |

### Kiracı / Sakin (`/pms/occupancies`)

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/pms/occupancies` | Sakin listesi | `pms_occupancies:read` |
| `GET` | `/pms/occupancies/:id` | Sakin detayı | `pms_occupancies:read` |
| `POST` | `/pms/occupancies` | Yeni sakin kaydı | `pms_occupancies:create` |
| `PUT` | `/pms/occupancies/:id` | Sakin güncelle | `pms_occupancies:update` |
| `DELETE` | `/pms/occupancies/:id` | Sakin kaydı sil | `pms_occupancies:delete` |
| `GET` | `/pms/occupancies/responsibility` | Sorumluluk çözümleme | `pms_occupancies:read` |

### Aidat Planları (`/pms/billing`)

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/pms/billing/dues-plans` | Aidat planı listesi | `pms_billing_plans:read` |
| `GET` | `/pms/billing/dues-plans/:id` | Plan detayı | `pms_billing_plans:read` |
| `POST` | `/pms/billing/dues-plans` | Yeni plan oluştur | `pms_billing_plans:create` |
| `PUT` | `/pms/billing/dues-plans/:id` | Plan güncelle | `pms_billing_plans:update` |
| `POST` | `/pms/billing/dues-plans/:id/activate` | Planı aktive et | `pms_billing_plans:activate` |
| `POST` | `/pms/billing/dues-plans/:id/archive` | Planı arşivle | `pms_billing_plans:archive` |

### Kontrol Paneli

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/pms/dashboard` | PMS özet metrikleri | `pms_dashboard:read` |

### Sakin Portalı (`/portal/pms`)

| Metod | Yol | Açıklama | Kimlik Doğrulama |
|-------|-----|----------|------------------|
| `GET` | `/portal/pms/me/units` | Sakinim birimlerim | JWT (sakin) |
| `GET` | `/portal/pms/me/dues` | Aidat borcum | JWT (sakin) |
| `GET` | `/portal/pms/me/payments` | Ödeme geçmişim | JWT (sakin) |
| `POST` | `/portal/pms/me/payments/initiate_card_payment` | Kart ödemesi başlat | JWT (sakin) |
| `POST` | `/portal/pms/me/payments/intents/:intent_id/finalize` | Ödeme tamamla | JWT (sakin) |
| `GET` | `/portal/pms/me/work-orders` | İş emirlerim | JWT (sakin) |
| `GET` | `/portal/pms/me/work-orders/:id` | İş emri detayı | JWT (sakin) |
| `POST` | `/portal/pms/me/work-orders` | Yeni bakım talebi | JWT (sakin) |
| `POST` | `/portal/pms/me/work-orders/:id/comment` | Talebe yorum ekle | JWT (sakin) |

---

## Temel Veri Modelleri

### Mülk (Property)

```json
{
  "id": "uuid",
  "site_id": "uuid",
  "name": "Güneş Sitesi A Blok",
  "type": "apartment_block",
  "address": "Ataşehir, İstanbul",
  "unit_count": 24,
  "created_at": "2026-01-01T00:00:00Z"
}
```

### Sakin (Occupancy)

```json
{
  "id": "uuid",
  "property_id": "uuid",
  "unit_id": "uuid",
  "tenant_name": "Ayşe Kaya",
  "tenant_type": "owner",
  "move_in_date": "2025-01-01",
  "move_out_date": null,
  "monthly_dues": 850.00,
  "currency": "TRY",
  "status": "active"
}
```

### Aidat Planı (Dues Plan)

```json
{
  "id": "uuid",
  "site_id": "uuid",
  "name": "2026 Ocak Aidatı",
  "amount": 850.00,
  "currency": "TRY",
  "billing_period": "monthly",
  "due_day": 5,
  "status": "active",
  "activated_at": "2026-01-01T00:00:00Z"
}
```

---

## Portal Rate Limiting

Sakin portalı endpoint'lerinde kötüye kullanımı önlemek için rate limiter uygulanır:

| Endpoint Grubu | Limit |
|----------------|-------|
| İş emri listeleme | 30 istek/dakika |
| İş emri detayı | 60 istek/dakika |
| İş emri oluşturma | 10 istek/dakika |
| İş emri yorum | 20 istek/dakika |

---

## Çok-Tenant Davranışı

PMS tüm işlemlerde `site_id` ile tam izolasyon uygular. Sakin portal erişimi JWT tabanlı olup kullanıcı yalnızca kendi birimine (unit) ait verileri görebilir.

---

## Ödeme Akışı

```
Sakin → POST /initiate_card_payment → [Ödeme Sağlayıcısı] → intent oluşur
                                                              ↓
Sakin → POST /intents/:id/finalize → Ödeme onaylanır / reddedilir
```

---

## Hata Kodları

| Kod | Açıklama |
|-----|----------|
| `400` | Geçersiz veri |
| `403` | İzin yetersiz / farklı tenant erişimi |
| `404` | Kaynak bulunamadı |
| `409` | Çakışan kayıt (aktif plan, aktif sakin vb.) |
| `429` | Rate limit aşıldı (portal endpoint'leri) |
| `422` | Ödeme başlatma başarısız |
