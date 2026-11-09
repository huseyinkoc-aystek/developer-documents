# SMS Servisi — sms

**Sürüm:** 1.0.0 | **Son Güncelleme:** 2026-03-31

## Genel Bakış

`sms` modülü, platform içinden SMS gönderimi için sağlayıcı (provider) entegrasyonlarını, mesaj şablonlarını ve gönderim politikalarını yönetir. Çoklu sağlayıcı desteği sunarak yük devretme (failover) ve maliyet optimizasyonu sağlar.

## İzinler

| İzin | Açıklama |
|------|----------|
| `sms:create` | Sağlayıcı / şablon oluşturma |
| `sms:read` | Listeleme ve detay |
| `sms:update` | Düzenleme |
| `sms:delete` | Silme |

## API Endpoint'leri

### Sağlayıcılar (Providers)

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/sms/providers` | Sağlayıcı listesi | `sms:read` |
| `GET` | `/sms/providers/:id` | Sağlayıcı detayı | `sms:read` |
| `POST` | `/sms/providers` | Yeni sağlayıcı ekle | `sms:create` |
| `PUT` | `/sms/providers/:id` | Sağlayıcı güncelle | `sms:update` |
| `DELETE` | `/sms/providers/:id` | Sağlayıcı sil | `sms:delete` |
| `POST` | `/sms/providers/:id/credentials/rotate` | API kimlik bilgilerini yenile | `sms:update` |
| `POST` | `/sms/providers/:id/credentials/activate` | Kimlik bilgilerini aktive et | `sms:update` |

### Şablonlar (Templates)

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/sms/templates` | Şablon listesi | `sms:read` |
| `GET` | `/sms/templates/:id` | Şablon detayı | `sms:read` |
| `POST` | `/sms/templates` | Yeni şablon oluştur | `sms:create` |
| `PUT` | `/sms/templates/:id` | Şablon güncelle | `sms:update` |
| `DELETE` | `/sms/templates/:id` | Şablon sil | `sms:delete` |

### Politikalar (Policies)

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/sms/policies` | SMS gönderim politikası | `sms:read` |
| `PUT` | `/sms/policies` | Politikayı güncelle | `sms:update` |

## Veri Modelleri

### Provider

```json
{
  "id": "uuid",
  "site_id": "uuid",
  "name": "Netgsm",
  "provider_type": "netgsm",
  "sender_name": "AYSTEK",
  "is_active": true,
  "is_default": true,
  "created_at": "2026-01-01T00:00:00Z"
}
```

### Şablon

```json
{
  "id": "uuid",
  "site_id": "uuid",
  "name": "otp_verification",
  "body": "Doğrulama kodunuz: {{code}}. Bu kod 5 dakika geçerlidir.",
  "variables": ["code"],
  "max_length": 160,
  "is_active": true
}
```

### Politika

```json
{
  "daily_send_limit": 10000,
  "per_user_hourly_limit": 5,
  "blackout_hours": { "start": "23:00", "end": "08:00" },
  "opt_out_keyword": "IPTAL",
  "default_provider_id": "uuid"
}
```

## Kimlik Bilgisi Rotasyonu

Sağlayıcı API anahtarları dönemsel olarak değiştirilebilir. Rotasyon akışı:

1. `POST /sms/providers/:id/credentials/rotate` — Yeni anahtar yüklenir (henüz aktif değil)
2. Test gönderimi yapılır
3. `POST /sms/providers/:id/credentials/activate` — Yeni anahtar aktive edilir, eski iptal olur

## Çok-Tenant Davranışı

Her SMS sağlayıcısı `site_id` ile izole edilir. API anahtarları sunucuda şifreli saklanır; API yanıtlarında maskelenerek gösterilir (`***`).

## Hata Kodları

| Kod | Açıklama |
|-----|----------|
| `400` | Geçersiz şablon değişkeni veya mesaj uzunluğu |
| `404` | Sağlayıcı / şablon bulunamadı |
| `422` | Sağlayıcı API doğrulama hatası |
| `429` | Günlük / saatlik limit aşıldı |
