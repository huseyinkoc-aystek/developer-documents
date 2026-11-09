# Firmalar — companies

**Sürüm:** 1.0.0 | **Son Güncelleme:** 2026-03-31

## Genel Bakış

`companies` modülü, site bünyesindeki firmaların ve sektörlerin yönetimini sağlar. Firma profilleri; sektör kategorisi, iletişim bilgileri ve sosyal medya bağlantıları içerir.

## İzinler

| İzin | Açıklama |
|------|----------|
| `companies:create` | Yeni firma / sektör oluşturma |
| `companies:read` | Listeleme ve detay görüntüleme |
| `companies:update` | Firma / sektör düzenleme |
| `companies:delete` | Firma / sektör silme |

## API Endpoint'leri

### Sektörler

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/company-sectors` | Sektör listesi | `companies:read` |
| `GET` | `/company-sectors/:id` | Tek sektör detayı | `companies:read` |
| `POST` | `/company-sectors` | Yeni sektör oluştur | `companies:create` |
| `PUT` | `/company-sectors/:id` | Sektör güncelle | `companies:update` |
| `DELETE` | `/company-sectors/:id` | Sektör sil | `companies:delete` |

### Firmalar

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/companies` | Firma listesi | `companies:read` |
| `GET` | `/companies/:id` | Tek firma detayı | `companies:read` |
| `POST` | `/companies` | Yeni firma oluştur | `companies:create` |
| `PUT` | `/companies/:id` | Firma güncelle | `companies:update` |
| `DELETE` | `/companies/:id` | Firma sil | `companies:delete` |

### Genel (Public) Endpoint'ler

| Metod | Yol | Açıklama |
|-------|-----|----------|
| `GET` | `/public/companies` | Firma listesi |
| `GET` | `/public/companies/:slug` | Firma detayı |
| `GET` | `/public/company-sectors` | Sektör listesi |

## Veri Modeli

```json
{
  "id": "uuid",
  "site_id": "uuid",
  "slug": "abc-teknoloji",
  "name": "ABC Teknoloji A.Ş.",
  "sector_id": "uuid",
  "description": "...",
  "logo_url": "https://...",
  "website": "https://abctech.com",
  "address": "İstanbul, Türkiye",
  "phone": "+90 212 000 0000",
  "email": "info@abctech.com",
  "social_links": { "linkedin": "..." },
  "is_published": true,
  "created_at": "2026-01-01T00:00:00Z"
}
```

## Çok-Tenant Davranışı

Her firma kaydı `site_id` ile tam izolasyon altındadır. Sektörler de aynı şekilde tenant'a özgüdür.

## Hata Kodları

| Kod | Açıklama |
|-----|----------|
| `400` | Geçersiz veri |
| `404` | Firma / sektör bulunamadı |
| `409` | Slug çakışması |
