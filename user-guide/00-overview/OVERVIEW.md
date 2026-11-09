---
Başlık: Platform Genel Bakış
Sürüm: 2.2.0
Durum: Üretime Hazır (Enterprise)
Son Güncelleme: 2026-03-31
İlgili Belgeler:
  - ../01-architecture/Enterprise_Isolation_Whitepaper.md
  - ../01-architecture/TENANT_AND_SITE_GUIDE.md
  - ../05-operations/SUNUCU_KURULUM_REHBERI.md
  - ../04-api/API_REFERENCE.md
---

# AYSTEK Core — Platform Genel Bakış

## İçindekiler

1. [Platform Tanımı](#1-platform-tanımı)
2. [Sistem Sınırları](#2-sistem-sınırları)
3. [Otorite Modeli](#3-otorite-modeli)
4. [Temel Özellikler](#4-temel-özellikler)
5. [Modül Listesi](#5-modül-listesi)
6. [Deployment Mimarisi](#6-deployment-mimarisi)
7. [Referans Haritası](#7-referans-haritası)

---

## 1. Platform Tanımı

**AYSTEK Core**, çok-tenant, çok-site destekli bir enterprise platform altyapısıdır. Kimlik doğrulama, yetkilendirme, modül kapısı, lisanslama, özellik bayrakları ve eklenti çalıştırma kontrollerini tek sunucu tarafı otorite modeliyle yönetir.

Platformun temel bileşenleri:

| Bileşen | Teknoloji | Rol |
|---------|-----------|-----|
| `aystek-core` | Go (REST API) | İş mantığı, veri, yetkilendirme |
| `aystek-panel` | Next.js (React) | Yönetim arayüzü |
| `aystek-core-installer` | Go CLI | Kurulum, update, rollback, backup, validate yaşam döngüsü yöneticisi |
| `aystek-cli` | Go binary | Sysadmin komut satırı aracı |
| MongoDB | NoSQL | Ana veri deposu |
| Redis | Önbellek/mesajlaşma | Oturum, kilit, bakım modu |
| RabbitMQ | Mesaj kuyruğu | Async iş işleme |

---

## 2. Sistem Sınırları

```
Admin UI (tüketici)  ---> AYSTEK Core API (otorite) ---> Veri/Olay Altyapısı
Genel İstemciler     ---> AYSTEK Core API (otorite) ---> Veri/Olay Altyapısı
aystek-core-installer ---> Sunucu / release / servis katmanı
```

---

## 3. Otorite Modeli

1. UI katmanları tüketicidir, politika motorları değildir.
2. Tenant/site izolasyonu sunucu tarafında zorlanır.
3. Hata yanıtı sözleşmesi deterministiktir ve öngörülebilirdir.
4. Genel ve admin yüzeyleri açık bağlam gereksinimlerine sahiptir.

---

## 4. Temel Özellikler

### İçerik Operasyonları
- Bağlayıcı tabanlı içerik akışı (RSS/API/Ajans/Scraper/Webhook) merkezi yönetimle
- Moderasyon öncelikli editoryal iş akışı (onayla/reddet/düzenle/zamanla)
- Site kapsamlı içerik operasyon ayarları (politika, bayraklar, kotalar)
- Outbox destekli olay akışı (güvenli entegrasyon ve tekrar oynatma)

### E-Ticaret ve Sipariş Yönetimi
- Ürün katalogu, stok, fiyatlandırma, iskonto, sipariş, iade
- B2B iş akışları, toplu fiyatlandırma, sözleşme yönetimi
- Çoklu para birimi, döviz kuru yönetimi
- Ödeme entegrasyonları, taksit desteği

### Çok-Tenant/Çok-Site
- Her tenant için bağımsız izolasyon
- Site başına özelleştirilebilir ayarlar
- Tenant/site bağlamı her okuma/yazma yolunda zorlanır

### Sistem Güncelleme ve Bakım
- Panel veya installer/CLI üzerinden kontrollü güncelleme
- Release symlink rotasyonu + health validation odaklı güncelleme akışı
- Tek tıkla rollback ve operatör onaylı geri alma
- Kalıcı maintenance state + guardrail policy desteği

### Güvenlik ve Uyumluluk
- JWT tabanlı kimlik doğrulama
- RBAC rol tabanlı erişim kontrolü
- KVKK/GDPR uyum araçları
- Kapsamlı denetim logu

---

## 5. Modül Listesi

Platformun Go backend'inde yer alan modüller:

| Kategori | Modüller |
|----------|---------|
| Kullanıcı ve Güvenlik | `auth`, `users`, `permissions`, `security` |
| İçerik | `posts`, `pages`, `media`, `taxonomy`, `tags`, `translations`, `localized_content` |
| İçerik Operasyonları | `content_ops_settings`, `content_bots`, `integrations`, `news_sources`, `news_agencies` |
| E-Ticaret | `catalog`, `products`, `carts`, `orders`, `checkout`, `payments`, `shipping`, `returns`, `stock` |
| Fiyatlandırma | `pricing`, `discounts`, `currencies`, `fx`, `taxes`, `installments` |
| B2B | `b2b`, `contracts`, `companies` |
| CRM | `crm`, `contacts`, `billing`, `ledger`, `subscriptions`, `subscription_contracts` |
| Haber/Medya | `biographies`, `interviews`, `digital_editions`, `classifieds`, `condolences`, `events`, `galleries` |
| Platform | `settings`, `plugins`, `licenses`, `notifications`, `languages`, `menus` |
| Sistem | `system`, `system_ops`, `maintenance`, `deployment` |
| Diğer | `ai`, `chatbot`, `telematics`, `pms`, `sms`, `surveys`, `forum`, `ads` |

Tam modül dokümantasyonu: [Modüller Dizini](./../06-modules/)

---

## 6. Deployment Mimarisi

```
[CentOS 9 Sunucu]
  ├── Nginx (443 → :9090 API, /panel → :3000 Panel)
  ├── aystek-core.service    (Go API — port 9090)
  ├── aystek-panel.service   (Next.js — port 3000)
  ├── mongod.service         (MongoDB)
  ├── redis.service          (Redis)
  └── rabbitmq-server.service (RabbitMQ)

/opt/aystek-core/
  ├── releases/
  │   ├── v1.1.0/   (önceki)
  │   └── v1.2.0/   (aktif)
  ├── current -> releases/v1.2.0
  └── previous -> releases/v1.1.0
```

Kurulum detayları: [Sunucu Kurulum Rehberi](../05-operations/SUNUCU_KURULUM_REHBERI.md)

---

## 7. Referans Haritası

| Konu | Belge |
|------|-------|
| İlk kurulum | [SUNUCU_KURULUM_REHBERI.md](../05-operations/SUNUCU_KURULUM_REHBERI.md) |
| Güncelleme/Rollback | [system_ops.md](../06-modules/modules/system_ops.md) |
| Bakım modu | [maintenance.md](../06-modules/modules/maintenance.md) |
| Mimari | [Enterprise_Isolation_Whitepaper.md](../01-architecture/Enterprise_Isolation_Whitepaper.md) |
| Güvenlik | [Enterprise_Security_Architecture_Booklet.md](../02-security/Enterprise_Security_Architecture_Booklet.md) |
| API Referansı | [API_REFERENCE.md](../04-api/API_REFERENCE.md) |
| Operasyonlar | [OPERATIONS_GUIDE.md](../05-operations/OPERATIONS_GUIDE.md) |
