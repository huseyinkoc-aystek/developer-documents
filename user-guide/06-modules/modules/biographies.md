# Biyografiler — biographies

**Sürüm:** 1.0.0 | **Son Güncelleme:** 2026-03-31

## Genel Bakış

`biographies` modülü, site içinde kişi profillerini (gazeteci, yazar, konuk yazar, uzman vb.) yönetir. Biyografiler slug tabanlı genel URL'lerle yayımlanır ve diğer içerik modülleriyle ilişkilendirilebilir.

## İzinler

| İzin | Açıklama |
|------|----------|
| `biographies:create` | Yeni biyografi oluşturma |
| `biographies:read` | Biyografi listeleme ve detay görüntüleme |
| `biographies:update` | Biyografi düzenleme |
| `biographies:delete` | Biyografi silme |

## API Endpoint'leri

### Yönetim Endpoint'leri

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/biographies` | Biyografi listesi (sayfalama, filtreleme) | `biographies:read` |
| `GET` | `/biographies/:id` | Tek biyografi detayı (ID ile) | `biographies:read` |
| `POST` | `/biographies` | Yeni biyografi oluştur | `biographies:create` |
| `PUT` | `/biographies/:id` | Biyografi güncelle | `biographies:update` |
| `DELETE` | `/biographies/:id` | Biyografi sil | `biographies:delete` |

### Genel (Public) Endpoint'ler

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/public/biographies` | Yayımlanan biyografilerin listesi | Yok |
| `GET` | `/public/biographies/:slug` | Slug ile tek biyografi | Yok |

## Veri Modeli

```json
{
  "id": "uuid",
  "site_id": "uuid",
  "slug": "ali-yilmaz",
  "full_name": "Ali Yılmaz",
  "title": "Kıdemli Muhabir",
  "bio": "...",
  "avatar_url": "https://...",
  "social_links": {
    "twitter": "https://twitter.com/...",
    "linkedin": "https://linkedin.com/..."
  },
  "is_published": true,
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-01-01T00:00:00Z"
}
```

## Çok-Tenant Davranışı

Her biyografi tek bir `site_id` ile ilişkilendirilir. Listeleme endpoint'leri her zaman aktif tenant'ın site bağlamında filtrelenir. Farklı tenant'ların biyografileri birbirinden izole edilmiştir.

## Örnek İstek

```http
POST /biographies
Content-Type: application/json
Authorization: Bearer <token>

{
  "slug": "ahmet-celik",
  "full_name": "Ahmet Çelik",
  "title": "Editör",
  "bio": "20 yıllık gazetecilik deneyimiyle...",
  "is_published": true
}
```

## Hata Kodları

| Kod | Açıklama |
|-----|----------|
| `400` | Geçersiz giriş verisi (slug çakışması, zorunlu alan eksik) |
| `403` | İzin yetersiz |
| `404` | Biyografi bulunamadı |
| `409` | Slug zaten kullanımda |
