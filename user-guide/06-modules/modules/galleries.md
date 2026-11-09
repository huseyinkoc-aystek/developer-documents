# Galeriler — galleries

**Sürüm:** 1.0.0 | **Son Güncelleme:** 2026-03-31

## Genel Bakış

`galleries` modülü, fotoğraf ve medya galerilerinin oluşturulmasını, düzenlenmesini ve yayımlanmasını sağlar. Galeriler birden fazla medya ögesini içerebilir; slug tabanlı genel URL'lerle erişilebilir.

## İzinler

| İzin | Açıklama |
|------|----------|
| `galleries:create` | Yeni galeri oluşturma |
| `galleries:read` | Listeleme ve detay |
| `galleries:update` | Düzenleme |
| `galleries:delete` | Silme |

## API Endpoint'leri

### Yönetim Endpoint'leri

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/galleries` | Galeri listesi | `galleries:read` |
| `GET` | `/galleries/:id` | Tek galeri detayı | `galleries:read` |
| `POST` | `/galleries` | Yeni galeri oluştur | `galleries:create` |
| `PUT` | `/galleries/:id` | Galeri güncelle | `galleries:update` |
| `DELETE` | `/galleries/:id` | Galeri sil | `galleries:delete` |

### Genel (Public) Endpoint'ler

| Metod | Yol | Açıklama |
|-------|-----|----------|
| `GET` | `/public/galleries` | Yayımlanan galeriler |
| `GET` | `/public/galleries/:slug` | Slug ile galeri detayı |

## Veri Modeli

```json
{
  "id": "uuid",
  "site_id": "uuid",
  "slug": "istanbul-itfaiye-2026",
  "title": "İstanbul İtfaiye Tatbikatı",
  "description": "...",
  "cover_image_url": "https://...",
  "items": [
    {
      "order": 1,
      "url": "https://.../img1.jpg",
      "caption": "Tatbikat başlangıcı",
      "type": "image"
    }
  ],
  "is_published": true,
  "published_at": "2026-01-15T09:00:00Z",
  "tags": ["itfaiye", "tatbikat"],
  "created_at": "2026-01-14T00:00:00Z"
}
```

## Çok-Tenant Davranışı

Galeriler `site_id` ile izole edilir. Genel API yalnızca `is_published: true` kayıtları gösterir.

## Hata Kodları

| Kod | Açıklama |
|-----|----------|
| `400` | Geçersiz medya URL'si veya eksik başlık |
| `404` | Galeri bulunamadı |
| `409` | Slug çakışması |
