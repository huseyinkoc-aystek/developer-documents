# İlanlar — classifieds

**Sürüm:** 1.0.0 | **Son Güncelleme:** 2026-03-31

## Genel Bakış

`classifieds` modülü, kullanıcıların ilan oluşturabildiği ve yönetebildiği bir ilan bölümü sunar. İlanlar; taslak, inceleme, yayım, satıldı ve arşiv durumlarından geçer. Editörler inceleme sürecini yönetir.

## Durum Akışı

```
Taslak → İncelemede → Yayımlandı → Satıldı / Arşivlendi
               ↓
           Reddedildi
```

## İzinler

| İzin | Açıklama |
|------|----------|
| `classifieds:create` | Yeni ilan oluşturma |
| `classifieds:read` | İlan listeleme ve detay görüntüleme |
| `classifieds:update` | İlan düzenleme |
| `classifieds:delete` | İlan silme |
| `classifieds:publish` | İlan yayımlama |
| `classifieds:review` | İnceleme sürecini yönetme |

## API Endpoint'leri

### Yönetim Endpoint'leri

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/classifieds` | İlan listesi | `classifieds:read` |
| `GET` | `/classifieds/:id` | Tek ilan detayı | `classifieds:read` |
| `POST` | `/classifieds` | Yeni ilan oluştur | `classifieds:create` |
| `PUT` | `/classifieds/:id` | İlan güncelle | `classifieds:update` |
| `DELETE` | `/classifieds/:id` | İlan sil | `classifieds:delete` |
| `POST` | `/classifieds/:id/submit-review` | İlanı incelemeye gönder | `classifieds:update` |
| `POST` | `/classifieds/:id/publish` | İlanı yayımla | `classifieds:publish` |
| `POST` | `/classifieds/:id/reject` | İlanı reddet | `classifieds:review` |
| `POST` | `/classifieds/:id/sold` | İlanı satıldı olarak işaretle | `classifieds:update` |
| `POST` | `/classifieds/:id/archive` | İlanı arşivle | `classifieds:update` |

### Genel (Public) Endpoint'ler

| Metod | Yol | Açıklama |
|-------|-----|----------|
| `GET` | `/public/classifieds` | Yayımlanan ilanların listesi |
| `GET` | `/public/classifieds/:slug` | Slug ile tek ilan |

## Veri Modeli

```json
{
  "id": "uuid",
  "site_id": "uuid",
  "slug": "satilik-daire-2024",
  "title": "Satılık 3+1 Daire",
  "description": "...",
  "category": "emlak",
  "price": 2500000,
  "currency": "TRY",
  "contact_info": { "phone": "...", "email": "..." },
  "images": ["https://..."],
  "status": "published",
  "created_by": "uuid",
  "expires_at": "2026-06-01T00:00:00Z",
  "created_at": "2026-01-01T00:00:00Z"
}
```

## Çok-Tenant Davranışı

Her ilan `site_id` ile izole edilir. İnceleme ve yayım izinleri tenant rollerine göre kontrol edilir.

## Hata Kodları

| Kod | Açıklama |
|-----|----------|
| `400` | Geçersiz durum geçişi |
| `403` | İzin yetersiz |
| `404` | İlan bulunamadı |
| `409` | Slug çakışması |
