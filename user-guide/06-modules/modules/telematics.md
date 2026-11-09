# Araç Takip ve Telematik — telematics

**Sürüm:** 1.0.0 | **Son Güncelleme:** 2026-03-31

## Genel Bakış

`telematics` modülü, araç filolarının gerçek zamanlı takibini, saha teknisyenlerinin saha biletlerini ve seyahat raporlarını yönetir. Hem yönetim paneli hem de mobil uygulama (saha personeli) API'leri içerir.

## Kullanım Senaryoları

| Senaryo | Açıklama |
|---------|----------|
| **Filo Yöneticisi** | Araç konumları, anlık durum, seyahat geçmişi |
| **Saha Teknisyeni** | Kendi araç verilerini ve saha biletlerini mobil üzerinden görme |
| **Saha Denetçisi** | Saha biletleri oluşturma, atama ve kapama |
| **Denetim** | Tüm telematik olayların denetim kaydı |

## İzinler

| İzin | Açıklama |
|------|----------|
| `telematics:read` | Araç ve telematik verilerini görüntüleme |
| `telematics:update` | Güncelleme |
| `telematics:field_tickets` | Saha biletlerini yönetme |

## API Endpoint'leri

### Yönetim Endpoint'leri

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/telematics/dashboard` | Filo özet tablosu | `telematics:read` |
| `GET` | `/telematics/vehicles/:id/overview` | Araç genel görünümü | `telematics:read` |

### Mobil Endpoint'ler (Saha Personeli)

| Metod | Yol | Açıklama | Kimlik Doğrulama |
|-------|-----|----------|------------------|
| `GET` | `/mobile/telematics/me/trips` | Kendi seyahatlerim | JWT (saha personeli) |
| `GET` | `/mobile/telematics/me/trips/:id` | Seyahat detayı | JWT (saha personeli) |
| `GET` | `/mobile/telematics/me/field-tickets` | Kendi saha biletlerim | JWT (saha personeli) |
| `GET` | `/mobile/telematics/me/field-tickets/:id` | Saha bileti detayı | JWT (saha personeli) |
| `POST` | `/mobile/telematics/me/field-tickets` | Yeni saha bileti oluştur | JWT (saha personeli) |
| `PUT` | `/mobile/telematics/me/field-tickets/:id` | Bilet güncelle / kapat | JWT (saha personeli) |

### Denetim

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/telematics/audit` | Denetim log listesi | `telematics:read` |
| `GET` | `/telematics/audit/:id` | Denetim log detayı | `telematics:read` |

## Veri Modelleri

### Araç Genel Görünümü

```json
{
  "vehicle_id": "uuid",
  "plate": "34 ABC 123",
  "status": "moving",
  "current_location": {
    "lat": 41.015,
    "lng": 28.979,
    "address": "Kadıköy, İstanbul",
    "timestamp": "2026-01-15T10:30:00Z"
  },
  "speed_kmh": 45,
  "driver": {
    "id": "uuid",
    "name": "Mehmet Yılmaz"
  },
  "last_event": "ignition_on",
  "odometer_km": 125400
}
```

### Seyahat (Trip)

```json
{
  "id": "uuid",
  "vehicle_id": "uuid",
  "driver_id": "uuid",
  "start_location": { "lat": 41.0, "lng": 28.9 },
  "end_location": { "lat": 40.9, "lng": 29.1 },
  "started_at": "2026-01-15T08:00:00Z",
  "ended_at": "2026-01-15T09:30:00Z",
  "distance_km": 42.5,
  "max_speed_kmh": 92,
  "avg_speed_kmh": 65,
  "status": "completed"
}
```

### Saha Bileti (Field Ticket)

```json
{
  "id": "uuid",
  "vehicle_id": "uuid",
  "assigned_to": "uuid",
  "title": "Arıza Bildirimi",
  "description": "Ön sol lastik basıncı düşük",
  "priority": "high",
  "status": "open",
  "location": { "lat": 40.9, "lng": 29.0 },
  "created_at": "2026-01-15T10:00:00Z",
  "resolved_at": null
}
```

## Kontrol Paneli Özet Metrikleri

```json
{
  "total_vehicles": 45,
  "active_vehicles": 28,
  "idle_vehicles": 12,
  "offline_vehicles": 5,
  "open_field_tickets": 7,
  "today_total_distance_km": 1240,
  "alerts": [
    { "type": "speeding", "count": 3 },
    { "type": "geofence_breach", "count": 1 }
  ]
}
```

## Çok-Tenant Davranışı

Tüm araç ve saha bileti verileri `site_id` ile izole edilir. Saha personeli JWT'sinde `vehicle_id` ve `driver_id` bağlamı bulunur; kullanıcı yalnızca kendi atanmış verilere erişebilir.

## Hata Kodları

| Kod | Açıklama |
|-----|----------|
| `400` | Geçersiz koordinat veya tarih aralığı |
| `403` | Başka sürücünün verisine erişim girişimi |
| `404` | Araç / seyahat / bilet bulunamadı |
| `422` | GPS verisi beklenen formatta değil |
