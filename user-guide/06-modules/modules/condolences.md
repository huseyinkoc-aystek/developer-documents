# Başsağlığı İlanları — condolences

**Sürüm:** 1.0.0 | **Son Güncelleme:** 2026-03-31

## Genel Bakış

`condolences` modülü, site üzerinde yayımlanan başsağlığı ilanlarını yönetir. İlanlar yönetim panelinden oluşturulup onaylanır; genel API üzerinden anonim erişime açılabilir.

## İzinler

| İzin | Açıklama |
|------|----------|
| `condolences:create` | Başsağlığı ilanı oluşturma |
| `condolences:read` | Listeleme ve detay görüntüleme |
| `condolences:update` | Düzenleme |
| `condolences:delete` | Silme |

## API Endpoint'leri

### Yönetim Endpoint'leri

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/condolences` | İlan listesi | `condolences:read` |
| `GET` | `/condolences/:id` | Tek ilan detayı | `condolences:read` |
| `POST` | `/condolences` | Yeni ilan oluştur | `condolences:create` |
| `PUT` | `/condolences/:id` | İlan güncelle | `condolences:update` |
| `DELETE` | `/condolences/:id` | İlan sil | `condolences:delete` |

### Genel (Public) Endpoint'ler

| Metod | Yol | Açıklama |
|-------|-----|----------|
| `GET` | `/public/condolences` | Yayımlanan ilanların listesi |
| `GET` | `/public/condolences/:slug` | Slug ile ilan detayı |

## Veri Modeli

```json
{
  "id": "uuid",
  "site_id": "uuid",
  "slug": "mehmet-yilmaz-2026",
  "name": "Mehmet Yılmaz",
  "date_of_birth": "1945-03-15",
  "date_of_death": "2026-01-10",
  "message": "Merhum için...",
  "contact_name": "Yılmaz Ailesi",
  "contact_phone": "+90 ...",
  "is_published": true,
  "published_at": "2026-01-11T00:00:00Z",
  "created_at": "2026-01-10T00:00:00Z"
}
```

## Çok-Tenant Davranışı

Her ilan `site_id` ile izole edilir. Genel API'de yalnızca `is_published: true` olan ilanlar döner.

## Hata Kodları

| Kod | Açıklama |
|-----|----------|
| `400` | Geçersiz tarih formatı veya eksik alan |
| `403` | İzin yetersiz |
| `404` | İlan bulunamadı |
