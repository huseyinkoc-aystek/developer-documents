# Etkinlikler — events

**Sürüm:** 1.0.0 | **Son Güncelleme:** 2026-03-31

## Genel Bakış

`events` modülü, site üzerinde etkinliklerin (konser, sergi, seminer, fuar vb.) oluşturulmasını, yönetilmesini ve yayımlanmasını sağlar. Etkinlikler taslak, yayımlı, iptal ve arşiv durumlarında bulunabilir.

## Durum Akışı

```
Taslak → Yayımlı → Arşivlendi
           ↓
         İptal
```

## İzinler

| İzin | Açıklama |
|------|----------|
| `events:create` | Yeni etkinlik oluşturma |
| `events:read` | Listeleme ve detay |
| `events:update` | Etkinlik düzenleme |
| `events:delete` | Etkinlik silme |
| `events:publish` | Etkinlik yayımlama |

## API Endpoint'leri

### Yönetim Endpoint'leri

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/events` | Etkinlik listesi | `events:read` |
| `GET` | `/events/:id` | Tek etkinlik detayı | `events:read` |
| `POST` | `/events` | Yeni etkinlik oluştur | `events:create` |
| `PUT` | `/events/:id` | Etkinlik güncelle | `events:update` |
| `DELETE` | `/events/:id` | Etkinlik sil | `events:delete` |
| `POST` | `/events/:id/publish` | Etkinliği yayımla | `events:publish` |
| `POST` | `/events/:id/cancel` | Etkinliği iptal et | `events:update` |
| `POST` | `/events/:id/archive` | Etkinliği arşivle | `events:update` |

### Genel (Public) Endpoint'ler

| Metod | Yol | Açıklama |
|-------|-----|----------|
| `GET` | `/public/events` | Yayımlanan etkinlikler |
| `GET` | `/public/events/:slug` | Slug ile etkinlik detayı |

## Veri Modeli

```json
{
  "id": "uuid",
  "site_id": "uuid",
  "slug": "bahar-konseri-2026",
  "title": "Bahar Konseri",
  "description": "...",
  "category": "muzik",
  "cover_image_url": "https://...",
  "start_date": "2026-04-15T19:00:00Z",
  "end_date": "2026-04-15T22:00:00Z",
  "venue": {
    "name": "Atatürk Kültür Merkezi",
    "address": "İstanbul",
    "coordinates": { "lat": 41.0, "lng": 28.9 }
  },
  "ticket_url": "https://biletix.com/...",
  "is_free": false,
  "price_range": "100-500 TRY",
  "status": "published",
  "created_at": "2026-01-01T00:00:00Z"
}
```

## Çok-Tenant Davranışı

Her etkinlik `site_id` ile izole edilir. Genel API yalnızca `status: published` olan etkinlikleri döndürür.

## Hata Kodları

| Kod | Açıklama |
|-----|----------|
| `400` | Geçersiz tarih aralığı veya durum geçişi |
| `403` | İzin yetersiz |
| `404` | Etkinlik bulunamadı |
| `409` | Slug çakışması |
