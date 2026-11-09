# Site İçerik Operasyon Ayarları — content_ops_settings

**Sürüm:** 1.0.0 | **Son Güncelleme:** 2026-03-31

## Genel Bakış

`content_ops_settings` modülü, her site için içerik operasyon politikalarını, özellik bayraklarını, içerik limitlerini, reklam politikalarını ve kod enjeksiyon ayarlarını merkezi olarak yönetir.

## İzinler

| İzin | Açıklama |
|------|----------|
| `content_ops_settings:read` | Ayarları görüntüleme |
| `content_ops_settings:update` | Ayarları güncelleme |

## API Endpoint'leri

Tüm endpoint'ler `/sites/:siteId` prefix'i altında yer alır.

| Metod | Yol | Açıklama | Yetki |
|-------|-----|----------|-------|
| `GET` | `/sites/:siteId/policy` | İçerik politikasını getir | `content_ops_settings:read` |
| `PUT` | `/sites/:siteId/policy` | İçerik politikasını güncelle | `content_ops_settings:update` |
| `GET` | `/sites/:siteId/feature-flags` | Özellik bayraklarını getir | `content_ops_settings:read` |
| `PUT` | `/sites/:siteId/feature-flags` | Özellik bayraklarını güncelle | `content_ops_settings:update` |
| `GET` | `/sites/:siteId/content-limits` | İçerik limitlerini getir | `content_ops_settings:read` |
| `PUT` | `/sites/:siteId/content-limits` | İçerik limitlerini güncelle | `content_ops_settings:update` |
| `GET` | `/sites/:siteId/ad-policy` | Reklam politikasını getir | `content_ops_settings:read` |
| `PUT` | `/sites/:siteId/ad-policy` | Reklam politikasını güncelle | `content_ops_settings:update` |
| `GET` | `/sites/:siteId/injection-codes` | Enjeksiyon kodlarını getir | `content_ops_settings:read` |
| `PUT` | `/sites/:siteId/injection-codes` | Enjeksiyon kodlarını güncelle | `content_ops_settings:update` |

## Ayar Kategorileri

### İçerik Politikası (`policy`)

```json
{
  "moderation_required": true,
  "auto_publish_on_approve": true,
  "allow_external_authors": false,
  "default_content_language": "tr"
}
```

### Özellik Bayrakları (`feature-flags`)

```json
{
  "comments_enabled": true,
  "ratings_enabled": false,
  "newsletters_enabled": true,
  "classifieds_enabled": true,
  "events_enabled": true
}
```

### İçerik Limitleri (`content-limits`)

```json
{
  "max_image_size_mb": 10,
  "max_images_per_post": 50,
  "max_video_duration_sec": 600,
  "max_tags_per_post": 20
}
```

### Reklam Politikası (`ad-policy`)

```json
{
  "ads_enabled": true,
  "ad_provider": "google_adsense",
  "ad_slots": ["header", "sidebar", "in-content"],
  "blocked_categories": ["adult", "gambling"]
}
```

### Enjeksiyon Kodları (`injection-codes`)

```json
{
  "head_scripts": "<script>...</script>",
  "body_start_scripts": "",
  "body_end_scripts": "<script>...</script>"
}
```

> **Güvenlik Uyarısı:** Enjeksiyon kodları yalnızca güvenilir yöneticiler tarafından güncellenmelidir. XSS riskine karşı kod içeriği sunucu tarafında validate edilir.

## Çok-Tenant Davranışı

Her site kendi izole ayar kümesine sahiptir. `:siteId` parametresi JWT'deki tenant bağlamıyla doğrulanır; farklı bir sitenin ayarlarına erişim `403` hatasıyla reddedilir.

## Hata Kodları

| Kod | Açıklama |
|-----|----------|
| `400` | Geçersiz ayar değeri |
| `403` | Başka sitenin ayarlarına erişim girişimi |
| `404` | Site bulunamadı |
