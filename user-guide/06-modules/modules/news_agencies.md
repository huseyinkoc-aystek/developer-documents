# Haber Ajansları — news_agencies

**Sürüm:** 1.0.0 | **Son Güncelleme:** 2026-03-31

## Genel Bakış

`news_agencies` modülü; AA, Reuters gibi dış haber ajanslarıyla entegrasyonu yönetir. Sağlayıcı (provider) API kimlik bilgileri ve bağlayıcı (connector) konfigürasyonları bu modül üzerinden tanımlanır.

## Mimari

```
[Ajans API] → [Provider Kimlik Bilgileri] → [Connector] → [İçerik Kuyruğu]
```

## İzinler

| İzin | Açıklama |
|------|----------|
| `news_agencies:create` | Sağlayıcı / bağlayıcı oluşturma |
| `news_agencies:read` | Listeleme ve detay |
| `news_agencies:update` | Düzenleme |
| `news_agencies:delete` | Silme |
| `news_agencies:run` | Bağlayıcıyı manuel tetikleme |

## API Endpoint'leri

### Sağlayıcılar (Providers)

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/news-agencies/providers` | Sağlayıcı listesi | `news_agencies:read` |
| `GET` | `/news-agencies/providers/:id` | Sağlayıcı detayı | `news_agencies:read` |
| `POST` | `/news-agencies/providers` | Yeni sağlayıcı oluştur | `news_agencies:create` |
| `PUT` | `/news-agencies/providers/:id` | Sağlayıcı güncelle | `news_agencies:update` |
| `DELETE` | `/news-agencies/providers/:id` | Sağlayıcı sil | `news_agencies:delete` |
| `POST` | `/news-agencies/providers/:id/disable` | Sağlayıcıyı devre dışı bırak | `news_agencies:update` |

### Bağlayıcılar (Connectors)

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/news-agencies/connectors` | Bağlayıcı listesi | `news_agencies:read` |
| `GET` | `/news-agencies/connectors/:id` | Bağlayıcı detayı | `news_agencies:read` |
| `POST` | `/news-agencies/connectors` | Yeni bağlayıcı oluştur | `news_agencies:create` |
| `PUT` | `/news-agencies/connectors/:id` | Bağlayıcı güncelle | `news_agencies:update` |
| `DELETE` | `/news-agencies/connectors/:id` | Bağlayıcı sil | `news_agencies:delete` |
| `POST` | `/news-agencies/connectors/:id/test` | Bağlantıyı test et | `news_agencies:read` |
| `POST` | `/news-agencies/connectors/:id/run` | Manuel tetikle | `news_agencies:run` |
| `GET` | `/news-agencies/connectors/:id/history` | Çalışma geçmişi | `news_agencies:read` |

## Veri Modelleri

### Provider

```json
{
  "id": "uuid",
  "site_id": "uuid",
  "name": "Anadolu Ajansı",
  "provider_type": "aa",
  "api_key": "***",
  "api_endpoint": "https://api.aa.com.tr/v1",
  "is_active": true,
  "created_at": "2026-01-01T00:00:00Z"
}
```

### Connector

```json
{
  "id": "uuid",
  "provider_id": "uuid",
  "name": "AA Ekonomi Haberleri",
  "category_filter": ["ekonomi", "finans"],
  "schedule_cron": "*/15 * * * *",
  "last_run_at": "2026-01-15T10:00:00Z",
  "last_run_status": "success",
  "is_active": true
}
```

## Çok-Tenant Davranışı

Sağlayıcı kimlik bilgileri tenant bazında şifrelenerek saklanır. Farklı tenant'lar aynı ajans için farklı API anahtarları kullanabilir.

## Hata Kodları

| Kod | Açıklama |
|-----|----------|
| `422` | API testi başarısız (geçersiz kimlik bilgileri) |
| `503` | Ajans API'si ulaşılamaz |
