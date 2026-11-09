---
Başlık: System Operations (Sistem İşlemleri) Modülü
Sürüm: 1.0.0
Durum: Genel Erişilebilirlik (GA)
Son Güncelleme: 2026-03-31
Modül Yolu: /internal/modules/system_ops
Öncelik: P0 - Kritik
Bağımlılıklar: MongoDB, Redis, RabbitMQ
İlgili Belgeler:
  - system.md
  - maintenance.md
  - deployment.md
  - ../05-operations/SUNUCU_KURULUM_REHBERI.md
---

# System Operations Modülü

## Özet

**System Operations (system_ops)** modülü; AYSTEK Core platformunun **güncelleme, geri alma (rollback), veritabanı migrasyonu, yedekleme, bakım modu ve denetim logu** işlemlerini yöneten P0-Kritik modüldür.

Sunucu üzerinde çalışan platform sürümünü yönetim paneli veya installer/CLI koordinasyonuyla güncelleme; herhangi bir sorun durumunda önceki release'e kontrollü geri dönme; backup, migration, maintenance, health, license ve guardrail işlemlerini izlenebilir audit logu ile kayıt altına alma temel işlevleridir.

**Temel Yetenekler:**
- Panel veya CLI üzerinden güvenli güncelleme başlatma
- Guardrail policy ile korunan güncelleme / rollback akışları
- Update ve rollback kayıtları, canlı log stream'i ve detay ekranları
- Kalıcı maintenance state + çoklu dil destekli maintenance mesajı
- Veritabanı migrasyon yönetimi
- Manuel ve otomatik yedekleme
- Güncelleme ilerlemesini Server-Sent Events (SSE) ile canlı izleme
- Tüm operasyonlar için denetim logu (audit log)
- License activation ve operational policy yönetimi

**Erişim modeli:** `system:*` ve `system_ops:*` izinleri **yalnızca platform kullanıcılarına** aittir. Tenant `admin` rolüne verilmez (`internal/seeds/role_seed.go` → `platformOnlyPermissionModules`); AYSTEK platform kullanıcıları (`type=platform`, `@aystek.com` veya `platform` rolü) modül RBAC'ını bypass eder. Müşteri tenant admini `/admin/system/*` uçlarında `403` alır (salt-okunur `/admin/system/maintenance/status` okunabilir kalır).

**Ana Endpoint Grubu:** `/admin/system/`  
**Veritabanı Koleksiyonları:** `system_updates`, `system_migrations`, `system_backups`, `system_audit_log`, `system_maintenance_state`, `system_guardrail_policies`, `system_license_state`

> **Uygulama Gerçekliği Notu (2026-04-05):** Bu modül artık tarihsel dokümanlarda geçen ayrı bir `aystek-updater` daemon'una dayanmaz. Güncelleme ve rollback kayıtları `system_ops` modülünde tutulur; sunucu tarafı release/backup/rollback işlemleri `aystek-core-installer` ile koordine edilir. Ayrıca MongoDB migration'ları rollback sırasında geri alınmaz.

---

## İçindekiler

1. [Mimari Genel Bakış](#1-mimari-genel-bakış)
2. [Güncelleme Yönetimi](#2-güncelleme-yönetimi)
3. [Rollback Yönetimi](#3-rollback-yönetimi)
4. [Migrasyon Yönetimi](#4-migrasyon-yönetimi)
5. [Yedekleme Yönetimi](#5-yedekleme-yönetimi)
6. [Bakım Modu](#6-bakım-modu)
7. [Sağlık ve Durum](#7-sağlık-ve-durum)
8. [Denetim Logu](#8-denetim-logu)
9. [API Referansı](#9-api-referansı)
10. [Veri Modelleri](#10-veri-modelleri)
11. [Güvenlik ve İzinler](#11-güvenlik-ve-i̇zinler)
12. [Sorun Giderme](#12-sorun-giderme)

---

## 1. Mimari Genel Bakış

```
┌─────────────────────────────────────────────────────────────────────┐
│                     system_ops Modülü                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │  Handler (HTTP)                                               │  │
│  │  ├── admin.go   → güncelleme, rollback, bakım, yedek         │  │
│  │  └── health.go  → /health, /health/live, /health/ready       │  │
│  └──────────────────────────────────────────────────────────────┘  │
│           │                                                          │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │  Service (İş Mantığı)                                         │  │
│  │  ├── StartUpdate()       → guardrail kontrollü update kaydı  │  │
│  │  ├── StartRollback()     → onaylı rollback kaydı             │  │
│  │  ├── TriggerMigrations() → migration record + apply/fail     │  │
│  │  ├── TriggerBackup()     → backup lifecycle                  │  │
│  │  ├── EnableMaintenance() → kalıcı maintenance state          │  │
│  │  └── GuardrailPolicy()   → operasyon güvence katmanı         │  │
│  └──────────────────────────────────────────────────────────────┘  │
│           │                                                          │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │  Worker                                                        │  │
│  │  └── manager.go → async job / SSE fan-out yöneticisi         │  │
│  └──────────────────────────────────────────────────────────────┘  │
│           │                                                          │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │  Repository (MongoDB)                                          │  │
│  │  ├── system_updates            → update + rollback kayıtları  │  │
│  │  ├── system_migrations         → migrasyon durumları          │  │
│  │  ├── system_backups            → yedek kayıtları              │  │
│  │  ├── system_audit_log          → tüm operasyon logları        │  │
│  │  ├── system_maintenance_state  → kalıcı bakım modu            │  │
│  │  ├── system_guardrail_policies → operasyon policy             │  │
│  │  └── system_license_state      → lisans durumu                │  │
│  └──────────────────────────────────────────────────────────────┘  │
│                                                                      │
│  Harici Bağımlılıklar:                                               │
│  ├── Redis  → cache ve mevcut runtime bağımlılıkları                │
│  ├── RabbitMQ → bağımlılık sağlık görünürlüğü                       │
│  └── MongoDB → tüm kayıtlar                                         │
└─────────────────────────────────────────────────────────────────────┘
```

### Dosya Yapısı

```
internal/modules/system_ops/
├── model/
│   ├── update.go      # UpdateRecord, UpdateJob, UpdateStep modelleri
│   ├── migration.go   # MigrationRecord modeli
│   ├── backup.go      # BackupRecord modeli
│   ├── health.go      # HealthStatus, ServiceStatus modelleri
│   └── audit.go       # AuditEntry modeli
├── repository/
│   ├── interface.go   # Repository arayüzü
│   └── mongo.go       # MongoDB implementasyonu
├── service/
│   └── service.go     # Tüm iş mantığı
├── handler/
│   ├── admin.go       # Admin HTTP handler'lar
│   └── health.go      # Sağlık endpoint'leri (RabbitMQ dahil)
├── worker/
│   └── manager.go       # Async update log fan-out yöneticisi
└── routes.go          # Route tanımları
```

---

## 2. Güncelleme Yönetimi

### 2.1 Güncelleme Akışı

Bir güncelleme başlatıldığında bugünkü üretim-gerçekliği aşağıdaki operasyon çerçevesiyle açıklanmalıdır:

| Adım | İşlem | Başarısız Olursa |
|------|-------|-----------------|
| 1 | Guardrail policy doğrulanır | Durdur |
| 2 | Update kaydı ve canlı log kanalı açılır | Durdur |
| 3 | Gerekli ise backup alınır | Uyarı / politika |
| 4 | Release ve config hazırlığı yapılır | Durdur |
| 5 | İdempotent migration/seed uygulanır | Rollback |
| 6 | `current` symlink yeni release'e döndürülür | Rollback |
| 7 | Servisler yeniden başlatılır | Rollback |
| 8 | Sağlık kontrolü yapılır | Rollback |
| 9 | Audit ve log kapanışı yapılır | — |

> ⚡ Release rotasyonunda `current` ve `previous` symlink'leri kullanılır. MongoDB migration'ları ise rollback sırasında geri alınmaz.

### 2.2 Güncelleme Başlatma (API)

```http
POST /admin/system/updates/start
Authorization: Bearer {admin_token}
Content-Type: application/json

{
  "target_version": "1.3.0",
  "channel": "stable"
}
```

Yanıt:
```json
{
  "success": true,
  "data": {
    "id": "upd_abc123",
    "record_id": "665f...",
    "status": "running",
    "stream_url": "/admin/system/updates/upd_abc123/logs"
  }
}
```

### 2.3 Güncelleme İlerleme Logu (SSE)

```http
GET /admin/system/updates/{jobId}/logs
Authorization: Bearer {admin_token}
Accept: text/event-stream
```

Stream formatı:
```
data: {"time":"2026-03-31T10:01:00Z","level":"info","step":"preflight","message":"Disk alanı: 2.3GB mevcut ✅"}
data: {"time":"2026-03-31T10:01:05Z","level":"info","step":"download","message":"Paket indiriliyor: 45%"}
data: {"time":"2026-03-31T10:03:00Z","level":"success","step":"complete","message":"Güncelleme tamamlandı: v1.3.0"}
```

Yönetim panelinde bu stream `LiveLogTerminal` bileşeni üzerinden canlı gösterilir.

### 2.4 Güncelleme Kontrolü

```http
GET /admin/system/updates/check
Authorization: Bearer {admin_token}
```

Yanıt:
```json
{
  "data": {
    "update_available": true,
    "latest_version": "1.3.0",
    "current_version": "1.2.0",
    "channel": "stable",
    "release_notes_url": "https://updates.aystek.com/releases/1.3.0/notes"
  }
}
```

---

## 3. Rollback Yönetimi

### 3.1 Rollback Akışı

| Adım | İşlem |
|------|-------|
| 1 | Rollback onayı ve guardrail doğrulanır |
| 2 | Gerekli ise pre-rollback backup alınır |
| 3 | Servisler durdurulur |
| 4 | `current` symlink önceki sürüme döndürülür |
| 5 | Servisler başlatılır |
| 6 | Sağlık kontrolü yapılır |
| 7 | Audit log kaydedilir |

> ⚠️ MongoDB migration'ları geri alınmaz. Rollback, release ve servis seviyesinde çalışır.

### 3.2 Rollback Başlatma (API)

```http
POST /admin/system/rollbacks/start
Authorization: Bearer {admin_token}
Content-Type: application/json

{
  "target_version": "1.1.0",
  "reason": "Kritik hata: ödeme modülünde gerileme"
}
```

> ⚠️ **Dikkat:** Rollback işlemi 2 adımlı onay gerektirir. Panel üzerinde onay kelimesi yazılmadan işlem başlamaz.

---

## 4. Migrasyon Yönetimi

### 4.1 Migrasyon Durumu Listesi

```http
GET /admin/system/migrations?page=1&limit=50
Authorization: Bearer {admin_token}
```

Yanıt:
```json
{
  "data": [
    {
      "_id": "665f...",
      "version": "1.2.0",
      "name": "0042_add_telematics_index",
      "status": "success",
      "applied_at": "2026-03-15T09:00:00Z"
    },
    {
      "_id": "665f...",
      "name": "0043_add_system_ops_collections",
      "status": "pending"
    }
  ],
  "total": 43
}
```

Migrasyon durumları: `pending` | `success` | `failed`

### 4.2 Manuel Migrasyon Tetikleme

```http
POST /admin/system/migrations/run
Authorization: Bearer {admin_token}
```

> Güncelleme sırasında migrasyonlar otomatik çalışır. Bu endpoint yalnızca geliştirici/operatör tarafından manuel tetiklemek için kullanılır.

---

## 5. Yedekleme Yönetimi

### 5.1 Yedek Tetikleme

```http
POST /admin/system/backups/trigger
Authorization: Bearer {admin_token}
```

Yedek türleri:
- `pre_update` — güncelleme öncesi otomatik alınır
- `manual` — operatör tarafından talep edilir
- `scheduled` — yapılandırılmış zamanlamaya göre

### 5.2 Yedek Listesi

```http
GET /admin/system/backups?page=1&limit=20
Authorization: Bearer {admin_token}
```

Yanıt alanları: `_id`, `type`, `status`, `path`, `size_bytes`, `triggered_by`, `started_at`, `completed_at`

---

## 6. Bakım Modu

### 6.1 Nasıl Çalışır

Bakım modu açıldığında iki katmanda uygulanır:

1. **Redis katmanı (birincil):** `SET system:maintenance_mode true` Redis'e yazılır. `MaintenanceMiddleware`, her istekte bu key'i kontrol eder (200ms timeout).
2. **Settings katmanı (yedek):** Redis erişilemeyen durumlarda `maintenance.enabled` ayarına bakılır.

```
İstek gelir
    ↓
/admin veya /health* → Her zaman geçer (operatör erişimini kapatma!)
    ↓
Redis'te system:maintenance_mode == "true"? → 503 + Retry-After: 120
    ↓
settings.maintenance.enabled == true? → 503 + Retry-After: 120
    ↓
Normal akış
```

### 6.2 Bakım Modunu Açma

```http
POST /admin/system/maintenance/enable
Authorization: Bearer {admin_token}
Content-Type: application/json

{
  "message": {
    "tr": "Sistem güncellemesi yapılmaktadır. Lütfen 15 dakika içinde tekrar deneyin.",
    "en": "Scheduled maintenance in progress. Please try again in 15 minutes."
  }
}
```

### 6.3 Bakım Modunu Kapatma

```http
POST /admin/system/maintenance/disable
Authorization: Bearer {admin_token}
```

### 6.4 Bakım Modu Durumu

```http
GET /admin/system/maintenance
Authorization: Bearer {admin_token}
```

Yanıt:
```json
{
  "data": {
    "enabled": true,
    "message": { "tr": "Güncelleme devam ediyor..." },
    "enabled_by": "admin@aystek.com",
    "enabled_at": "2026-03-31T10:00:00Z"
  }
}
```

### 6.5 Operasyon Erişimi

```bash
# Bakım modu, guardrail ve lisans operasyonları yönetim paneli veya
# ilgili admin API endpoint'leri üzerinden yürütülür.
# Host-level kurulum/update/rollback işlemleri için kanonik araç:
sudo aystek-core-installer status
sudo aystek-core-installer update
sudo aystek-core-installer rollback
```

---

## 7. Sağlık ve Durum

### 7.1 Detaylı Sağlık Kontrolü

```http
GET /admin/system/health
```

Bu endpoint **herkesin erişebildiği** (auth gerekmez) sağlık endpoint'idir. Yük dengeleyiciler ve izleme sistemleri tarafından kullanılır.

Yanıt:
```json
{
  "data": {
    "status": "healthy",
    "version": {
      "backend": "1.2.0",
      "git_commit": "abc123",
      "build_time": "2026-03-31T08:00:00Z",
      "go_version": "go1.22.2"
    },
    "uptime_seconds": 86400,
    "services": {
      "mongodb":  { "name": "mongodb",  "status": "ok", "latency_ms": 2 },
      "redis":    { "name": "redis",    "status": "ok", "latency_ms": 1 },
      "rabbitmq": { "name": "rabbitmq", "status": "ok", "latency_ms": 3 }
    },
    "node": {
      "hostname": "prod-server-01",
      "os": "linux",
      "arch": "amd64",
      "num_cpu": 4,
      "go_routines": 42,
      "memory": { "alloc_mb": 128, "sys_mb": 256 },
      "disk": { "total_gb": 100, "free_gb": 60, "used_pct": 40 }
    }
  }
}
```

Genel durum (status) değerleri:
- `healthy` — tüm servisler normal
- `degraded` — bir veya birden fazla servis yavaş/hatalı ama platform çalışıyor
- `unhealthy` — kritik servis çalışmıyor

### 7.2 Canlılık ve Hazırlık Probları

```http
GET /health/live    # Kubernetes liveness probe
GET /health/ready   # Kubernetes readiness probe
```

Bu endpoint'ler bakım modu sırasında da erişilebilir. `MaintenanceMiddleware` bunları her zaman geçirir.

---

## 8. Denetim Logu

Tüm sistem operasyonları otomatik olarak `system_audit_log` koleksiyonuna kaydedilir.

### 8.1 Audit Log Listesi

```http
GET /admin/system/audit?page=1&limit=50&action=update_start
Authorization: Bearer {admin_token}
```

Desteklenen filtreler: `action`, `actor`, `page`, `limit`

Yanıt örneği:
```json
{
  "data": [
    {
      "_id": "665f...",
      "action": "update_start",
      "actor": "admin@aystek.com",
      "actor_ip": "10.0.0.5",
      "resource": "update",
      "resource_id": "upd_abc123",
      "before": { "version": "1.2.0" },
      "after":  { "version": "1.3.0" },
      "result": "success",
      "created_at": "2026-03-31T10:00:00Z"
    }
  ]
}
```

Kayıt edilen başlıca `action` değerleri:

| Action | Açıklama |
|--------|---------|
| `update_start` | Güncelleme başlatıldı |
| `update_complete` | Güncelleme tamamlandı |
| `update_failed` | Güncelleme başarısız |
| `rollback_start` | Rollback başlatıldı |
| `rollback_complete` | Rollback tamamlandı |
| `maintenance_enable` | Bakım modu açıldı |
| `maintenance_disable` | Bakım modu kapatıldı |
| `migration_run` | Migrasyon çalıştırıldı |
| `backup_trigger` | Yedekleme başlatıldı |
| `license_activate` | Lisans aktive edildi |

---

## 9. API Referansı

### Tüm Endpoint'ler

| Metod | Yol | Açıklama | Auth |
|-------|-----|---------|------|
| GET | `/health` | Detaylı sağlık | ❌ |
| GET | `/health/live` | Canlılık probu | ❌ |
| GET | `/health/ready` | Hazırlık probu | ❌ |
| GET | `/admin/system/version` | Platform sürümü | ✅ |
| GET | `/admin/system/node` | Sunucu metrikleri | ✅ |
| GET | `/admin/system/services` | Servis durumları | ✅ |
| GET | `/admin/system/updates/check` | Güncelleme var mı? | ✅ |
| GET | `/admin/system/updates` | Güncelleme geçmişi | ✅ |
| POST | `/admin/system/updates/start` | Güncelleme başlat | ✅ admin |
| GET | `/admin/system/updates/:id` | Güncelleme detayı | ✅ |
| GET | `/admin/system/updates/:id/logs` | Canlı log (SSE) | ✅ |
| GET | `/admin/system/rollbacks` | Rollback geçmişi | ✅ |
| POST | `/admin/system/rollbacks/start` | Rollback başlat | ✅ admin |
| GET | `/admin/system/migrations` | Migrasyon listesi | ✅ |
| POST | `/admin/system/migrations/run` | Migrasyon çalıştır | ✅ admin |
| GET | `/admin/system/backups` | Yedek listesi | ✅ |
| POST | `/admin/system/backups/trigger` | Yedek al | ✅ operator |
| GET | `/admin/system/maintenance` | Bakım modu durumu | ✅ |
| POST | `/admin/system/maintenance/enable` | Bakım modunu aç | ✅ operator |
| POST | `/admin/system/maintenance/disable` | Bakım modunu kapat | ✅ operator |
| GET | `/admin/system/audit` | Audit log listesi | ✅ |
| GET | `/admin/system/license` | Lisans durumu | ✅ admin |
| POST | `/admin/system/license/activate` | Lisans aktive et | ✅ super_admin |

---

## 10. Veri Modelleri

### UpdateRecord

```go
type UpdateRecord struct {
    ID              primitive.ObjectID  `bson:"_id"`
    FromVersion     string              `bson:"from_version"`
    ToVersion       string              `bson:"to_version"`
    Channel         string              `bson:"channel"`       // stable|rc|beta
    Status          string              `bson:"status"`        // pending|running|success|failed|rolled_back
    TriggerSource   string              `bson:"trigger_source"` // panel|cli|api|schedule
    TriggeredBy     string              `bson:"triggered_by"`
    StartedAt       time.Time           `bson:"started_at"`
    CompletedAt     *time.Time          `bson:"completed_at,omitempty"`
    DurationSeconds *float64            `bson:"duration_seconds,omitempty"`
    ErrorMessage    string              `bson:"error_message,omitempty"`
}
```

### MigrationRecord

```go
type MigrationRecord struct {
    ID           primitive.ObjectID `bson:"_id"`
    Version      string             `bson:"version"`
    Name         string             `bson:"name"`
    Status       string             `bson:"status"` // pending|success|failed
    AppliedAt    *time.Time         `bson:"applied_at,omitempty"`
    ErrorMessage string             `bson:"error_message,omitempty"`
}
```

### AuditEntry

```go
type AuditEntry struct {
    ID         primitive.ObjectID     `bson:"_id"`
    Action     string                 `bson:"action"`
    Actor      string                 `bson:"actor"`
    ActorIP    string                 `bson:"actor_ip,omitempty"`
    Resource   string                 `bson:"resource"`
    ResourceID string                 `bson:"resource_id,omitempty"`
    Before     map[string]interface{} `bson:"before,omitempty"`
    After      map[string]interface{} `bson:"after,omitempty"`
    Result     string                 `bson:"result"` // success|failed
    CreatedAt  time.Time              `bson:"created_at"`
}
```

---

## 11. Güvenlik ve İzinler

### RBAC Tablosu

| İzin | Açıklama | Super Admin | Admin | Operator | Görüntüleyici |
|------|---------|:-----------:|:-----:|:--------:|:-------------:|
| `system_ops:read` | Tüm ekranları görüntüle | ✅ | ✅ | ✅ | ✅ |
| `system_ops:update_start` | Güncelleme başlat | ✅ | ✅ | ❌ | ❌ |
| `system_ops:rollback` | Rollback başlat | ✅ | ✅ | ❌ | ❌ |
| `system_ops:maintenance` | Bakım modu yönet | ✅ | ✅ | ✅ | ❌ |
| `system_ops:backup` | Yedek tetikle | ✅ | ✅ | ✅ | ❌ |
| `system_ops:config` | Sistem ayarı güncelle | ✅ | ❌ | ❌ | ❌ |
| `system_ops:license` | Lisans yönet | ✅ | ❌ | ❌ | ❌ |

### Kritik Operasyon Güvenceleri

- **Güncelleme/Rollback:** 2 adımlı onay gerektirir (panel confirm_token)
- **Eşzamanlı güncelleme koruması:** Redis dağıtık kilit — aynı anda ikinci güncelleme başlatılamaz
- **Download doğrulama:** SHA256 checksum + isteğe bağlı GPG imza
- **Rate limiting:** Update start: 2 istek/30sn
- **Her operasyon:** IP adresiyle birlikte `system_audit_log`'a yazılır

---

## 12. Sorun Giderme

### Güncelleme Takıldı

```bash
# Çalışan iş var mı
curl -s http://127.0.0.1:9090/admin/system/updates \
  -H "Authorization: Bearer TOKEN" | python3 -m json.tool

# Bakım modu takıldı mı
redis-cli GET system:maintenance_mode

# Bakım modunu manuel kaldır
redis-cli DEL system:maintenance_mode

# Symlink durumu
ls -la /opt/aystek-core/current
ls -la /opt/aystek-core/previous
ls /opt/aystek-core/releases/
```

### Rollback Başarısız

```bash
# Acil manuel rollback
systemctl stop aystek-core aystek-panel
PREV="1.1.0"
ln -sfn /opt/aystek-core/releases/v${PREV} /opt/aystek-core/current
systemctl start aystek-core aystek-panel
curl -s http://127.0.0.1:9090/health
```

### Migrasyon Başarısız

```bash
# Migrasyon durumu
curl -s "http://127.0.0.1:9090/admin/system/migrations?limit=10" \
  -H "Authorization: Bearer TOKEN" | python3 -m json.tool

# Migrasyon detay logu
tail -100 /var/log/aystek-core/migrations/migration.log
```
