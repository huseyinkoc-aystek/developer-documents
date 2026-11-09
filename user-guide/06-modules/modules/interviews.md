# Röportajlar — interviews

**Sürüm:** 1.0.0 | **Son Güncelleme:** 2026-03-31

## Genel Bakış

`interviews` modülü, yönetim panelinde röportaj içeriklerinin oluşturulmasını ve yayımlanmasını sağlar. Röportajlar soru-cevap formatında ya da uzun form metin olarak girilebilir; biyografi modülüyle ilişkilendirilebilir.

## İzinler

| İzin | Açıklama |
|------|----------|
| `interviews:create` | Röportaj oluşturma |
| `interviews:read` | Listeleme ve detay |
| `interviews:update` | Düzenleme |
| `interviews:delete` | Silme |

## API Endpoint'leri

### Yönetim Endpoint'leri

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/interviews` | Röportaj listesi | `interviews:read` |
| `GET` | `/interviews/:id` | Tek röportaj detayı | `interviews:read` |
| `POST` | `/interviews` | Yeni röportaj oluştur | `interviews:create` |
| `PUT` | `/interviews/:id` | Röportaj güncelle | `interviews:update` |
| `DELETE` | `/interviews/:id` | Röportaj sil | `interviews:delete` |

### Genel (Public) Endpoint'ler

| Metod | Yol | Açıklama |
|-------|-----|----------|
| `GET` | `/public/interviews` | Yayımlanan röportajlar |
| `GET` | `/public/interviews/:slug` | Slug ile röportaj detayı |

## Veri Modeli

```json
{
  "id": "uuid",
  "site_id": "uuid",
  "slug": "ahmet-celik-roportaj-2026",
  "title": "Teknoloji ve Gelecek",
  "interviewee": {
    "name": "Ahmet Çelik",
    "title": "CEO, ABC Teknoloji",
    "biography_id": "uuid"
  },
  "cover_image_url": "https://...",
  "content": "...",
  "qa_format": [
    { "question": "Soru?", "answer": "Cevap..." }
  ],
  "is_published": true,
  "published_at": "2026-01-15T09:00:00Z",
  "tags": ["teknoloji", "startup"],
  "created_at": "2026-01-14T00:00:00Z"
}
```

## Çok-Tenant Davranışı

Röportajlar `site_id` ile izole edilir. `biography_id` alanı aynı tenant'a ait bir biyografiye referans vermelidir; çapraz tenant erişimi reddedilir.

## Hata Kodları

| Kod | Açıklama |
|-----|----------|
| `400` | Zorunlu alan eksik |
| `404` | Röportaj bulunamadı |
| `409` | Slug çakışması |
