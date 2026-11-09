# Haber Kaynakları — news_sources

**Sürüm:** 1.0.0 | **Son Güncelleme:** 2026-03-31

## Genel Bakış

`news_sources` modülü, RSS / web kazıma yoluyla içerik toplanan haber kaynaklarını yönetir. Her kaynak için tarama geçmişi, başarı metrikleri ve hata logları takip edilir.

## İzinler

| İzin | Açıklama |
|------|----------|
| `news_sources:create` | Kaynak oluşturma |
| `news_sources:read` | Listeleme ve detay |
| `news_sources:update` | Düzenleme |
| `news_sources:delete` | Silme |

## API Endpoint'leri

### Haber Kaynakları

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/news-sources` | Kaynak listesi | `news_sources:read` |
| `GET` | `/news-sources/:id` | Tek kaynak detayı | `news_sources:read` |
| `POST` | `/news-sources` | Yeni kaynak oluştur | `news_sources:create` |
| `PUT` | `/news-sources/:id` | Kaynak güncelle | `news_sources:update` |
| `DELETE` | `/news-sources/:id` | Kaynak sil | `news_sources:delete` |
| `POST` | `/news-sources/:id/disable` | Kaynağı devre dışı bırak | `news_sources:update` |

### İzleme ve Geçmiş

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/news-sources/:id/crawl-history` | Tarama geçmişi | `news_sources:read` |
| `GET` | `/news-sources/:id/metrics` | Başarı / hata metrikleri | `news_sources:read` |
| `POST` | `/internal/news-sources/:id/crawl-result` | Tarama sonucu kaydet (dahili) | Sistem |

## Veri Modeli

```json
{
  "id": "uuid",
  "site_id": "uuid",
  "name": "Hürriyet RSS",
  "source_url": "https://www.hurriyet.com.tr/rss/anasayfa",
  "source_type": "rss",
  "schedule_cron": "*/20 * * * *",
  "is_active": true,
  "last_crawl_at": "2026-01-15T10:00:00Z",
  "last_crawl_status": "success",
  "total_fetched": 4200,
  "total_errors": 12,
  "created_at": "2026-01-01T00:00:00Z"
}
```

## Tarama Geçmişi Yanıtı

```json
{
  "data": [
    {
      "run_id": "uuid",
      "started_at": "2026-01-15T10:00:00Z",
      "completed_at": "2026-01-15T10:01:30Z",
      "fetched_count": 15,
      "new_items": 3,
      "status": "success",
      "error_message": null
    }
  ],
  "pagination": { "page": 1, "total": 250, "per_page": 20 }
}
```

## Çok-Tenant Davranışı

Her haber kaynağı `site_id` ile izole edilir. Dahili `crawl-result` endpoint'i yalnızca sistem servisleri tarafından çağrılır.

## Hata Kodları

| Kod | Açıklama |
|-----|----------|
| `400` | Geçersiz URL veya cron ifadesi |
| `404` | Kaynak bulunamadı |
| `422` | Kaynak erişilemez veya geçersiz format |
