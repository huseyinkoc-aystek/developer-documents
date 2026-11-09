# Dijital Yayınlar — digital_editions

**Sürüm:** 1.0.0 | **Son Güncelleme:** 2026-03-31

## Genel Bakış

`digital_editions` modülü, gazete / dergi gibi periyodik içeriklerin dijital baskılarını (e-gazete, e-dergi, PDF arşiv) yönetir. Yayımlanmış baskılar genel API üzerinden slug bazlı erişime açılır.

## İzinler

| İzin | Açıklama |
|------|----------|
| `digital_editions:create` | Yeni baskı oluşturma |
| `digital_editions:read` | Listeleme ve detay |
| `digital_editions:update` | Düzenleme |
| `digital_editions:delete` | Silme |

## API Endpoint'leri

### Yönetim Endpoint'leri

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/digital-editions` | Baskı listesi | `digital_editions:read` |
| `GET` | `/digital-editions/:id` | Tek baskı detayı | `digital_editions:read` |
| `POST` | `/digital-editions` | Yeni baskı oluştur | `digital_editions:create` |
| `PUT` | `/digital-editions/:id` | Baskı güncelle | `digital_editions:update` |
| `DELETE` | `/digital-editions/:id` | Baskı sil | `digital_editions:delete` |

### Genel (Public) Endpoint'ler

| Metod | Yol | Açıklama |
|-------|-----|----------|
| `GET` | `/public/digital-editions` | Baskı listesi |
| `GET` | `/public/digital-editions/:slug` | Slug ile baskı detayı |

## Veri Modeli

```json
{
  "id": "uuid",
  "site_id": "uuid",
  "slug": "gazete-15-ocak-2026",
  "title": "15 Ocak 2026 Sayısı",
  "edition_date": "2026-01-15",
  "cover_image_url": "https://...",
  "file_url": "https://.../edition.pdf",
  "file_size_mb": 12.5,
  "page_count": 24,
  "is_published": true,
  "published_at": "2026-01-15T06:00:00Z",
  "created_at": "2026-01-14T20:00:00Z"
}
```

## Çok-Tenant Davranışı

Baskılar `site_id` ile izole edilir. Genel API yalnızca `is_published: true` kayıtları döndürür.

## Hata Kodları

| Kod | Açıklama |
|-----|----------|
| `400` | Geçersiz tarih veya dosya formatı |
| `404` | Baskı bulunamadı |
| `409` | Slug çakışması |
