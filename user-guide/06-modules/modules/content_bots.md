# İçerik Botları — content_bots

**Sürüm:** 1.0.0 | **Son Güncelleme:** 2026-03-31

## Genel Bakış

`content_bots` modülü, dış kaynaklardan otomatik içerik çekme iş akışlarını yönetir. Bot bağlayıcıları (connector) tanımlanır, test edilir ve düzenli aralıklarla içerik çekme görevi tetiklenir. Çekilen içerikler moderasyon kuyruğuna düşer.

## Mimari

```
[Dış Kaynak] → [Bot Connector] → [Ingestion Worker] → [Moderasyon Kuyruğu] → [Yayım]
```

## İzinler

| İzin | Açıklama |
|------|----------|
| `content_bots:create` | Bağlayıcı oluşturma |
| `content_bots:read` | Listeleme ve detay |
| `content_bots:update` | Bağlayıcı düzenleme |
| `content_bots:delete` | Bağlayıcı silme |
| `content_bots:run` | Manuel tetikleme |

## API Endpoint'leri

### Bot Bağlayıcıları

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/bot-connectors` | Bağlayıcı listesi | `content_bots:read` |
| `GET` | `/bot-connectors/:id` | Bağlayıcı detayı | `content_bots:read` |
| `POST` | `/bot-connectors` | Yeni bağlayıcı oluştur | `content_bots:create` |
| `PUT` | `/bot-connectors/:id` | Bağlayıcı güncelle | `content_bots:update` |
| `DELETE` | `/bot-connectors/:id` | Bağlayıcı sil | `content_bots:delete` |
| `POST` | `/bot-connectors/:id/disable` | Bağlayıcıyı devre dışı bırak | `content_bots:update` |
| `POST` | `/bot-connectors/:id/test-connection` | Bağlantıyı test et | `content_bots:read` |
| `POST` | `/bot-connectors/:id/run` | İçerik çekmeyi manuel tetikle | `content_bots:run` |

### Ingestion

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/ingestion/connectors` | Tüm ingestion bağlayıcıları | `content_bots:read` |

## Veri Modeli

```json
{
  "id": "uuid",
  "site_id": "uuid",
  "name": "RSS Bot - Haber Ajansı",
  "source_type": "rss",
  "source_url": "https://example.com/feed.xml",
  "schedule_cron": "*/30 * * * *",
  "is_active": true,
  "last_run_at": "2026-01-15T10:00:00Z",
  "last_run_status": "success",
  "fetched_count": 150,
  "created_at": "2026-01-01T00:00:00Z"
}
```

## Desteklenen Kaynak Türleri

| Tür | Açıklama |
|-----|----------|
| `rss` | RSS / Atom besleme URL'si |
| `api` | Özel REST API uç noktası |
| `scraper` | Web kazıma kuralları |
| `webhook` | Gelen webhook dinleyici |

## Çok-Tenant Davranışı

Her bağlayıcı `site_id` ile izole edilir. Çalışma zamanı logları tenant bazında filtrelenir.

## Hata Kodları

| Kod | Açıklama |
|-----|----------|
| `400` | Geçersiz cron ifadesi veya kaynak URL |
| `404` | Bağlayıcı bulunamadı |
| `422` | Bağlantı testi başarısız |
| `503` | Kaynak ulaşılamaz |
