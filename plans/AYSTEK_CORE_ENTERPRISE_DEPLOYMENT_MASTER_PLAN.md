# AYSTEK CORE — Enterprise Deployment, Packaging, Versioning & Update Management  
## Master Teknik Tasarım Dokümanı

**Sürüm:** 1.0.0  
**Tarih:** 2026-03-31  
**Hazırlayan:** Enterprise Solution Architect  
**Durum:** MASTER PLAN — Implementasyon öncesi onay belgesi  
**Hedef Ortam:** CentOS 9 (birincil), Ubuntu 22.04+ (ikincil), Containerized (opsiyonel)  

---

## İçindekiler

1. [Executive Summary](#1-executive-summary)
2. [Önerilen Mimari](#2-önerilen-mimari)
3. [Neden Bu Mimari](#3-neden-bu-mimari)
4. [Alternatiflerin Karşılaştırması](#4-alternatiflerin-karşılaştırması)
5. [Deployment Topology](#5-deployment-topology)
6. [Paketleme Stratejisi](#6-paketleme-stratejisi)
7. [Sürümleme ve Release Yönetimi](#7-sürümleme-ve-release-yönetimi)
8. [Update / Rollback Akış Diyagramı](#8-update--rollback-akış-diyagramı)
9. [Yönetim Paneli Modül Tasarımı](#9-yönetim-paneli-modül-tasarımı)
10. [Go Backend API Tasarımı](#10-go-backend-api-tasarımı)
11. [Next.js UI Modül Tasarımı](#11-nextjs-ui-modül-tasarımı)
12. [Veri Modeli / Koleksiyon Önerileri](#12-veri-modeli--koleksiyon-önerileri)
13. [Dizin Yapısı ve systemd Servis Planı](#13-dizin-yapısı-ve-systemd-servis-planı)
14. [Güvenlik ve Audit Planı](#14-güvenlik-ve-audit-planı)
15. [CI/CD ve Release Pipeline Planı](#15-cicd-ve-release-pipeline-planı)
16. [Fazlara Bölünmüş Implementasyon Roadmap](#16-fazlara-bölünmüş-implementasyon-roadmap)
17. [Her Faz İçin Görev Listesi](#17-her-faz-için-görev-listesi)
18. [Riskler ve Mitigasyonlar](#18-riskler-ve-mitigasyonlar)
19. [Acceptance Criteria](#19-acceptance-criteria)
20. [Dosya Dosya Implementasyon Çalışma Planı](#20-dosya-dosya-implementasyon-çalışma-planı)

---

## 1. Executive Summary

AYSTEK CORE, çok-tenant/çok-site destekli, Go backend + Next.js admin frontend'den oluşan bir enterprise platform. Şu an için deployment süreci manuel ve ölçeklenemez. Bu plan, platformu aşağıdaki yeteneklerle donatmayı hedefler:

| Yetenek | Hedef |
|---|---|
| Dağıtım | Tek komutla kurulum, sıfır bağımlılık sürtüşmesi |
| Güncelleme | Panel üzerinden veya CLI ile, near-zero-downtime |
| Rollback | Tek tık veya tek komut, otomatik tetikleme |
| Bütünlük | Checksum + imzalı paketler, tedarik zinciri güvenliği |
| Gözlemlenebilirlik | Health, metrics, logs, migration durumu panelden izlenebilir |
| Güvenlik | RBAC + audit log + onay akışı + sistemd hardenig + SELinux |
| CI/CD | Git tag → build → test → sign → publish → notify |

**Temel karar:** Tarball + installer script + systemd + symlink-tabanlı release yapısı **hibrit yaklaşım** seçilmiştir. Detay Bölüm 4'te.

---

## 2. Önerilen Mimari

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        AYSTEK CORE DEPLOYMENT STACK                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                               │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                      UPDATE SERVER (AYSTEK HQ)                        │   │
│  │  - Release manifest JSON API                                          │   │
│  │  - Signed tarball storage (S3 / CDN / private)                       │   │
│  │  - License validation API                                             │   │
│  │  - Version compatibility API                                          │   │
│  └──────────────────────────┬─────────────────────────────────────────┘   │
│                               │ HTTPS (mutual TLS veya signed JWT)          │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                    CUSTOMER SERVER (CentOS 9)                          │   │
│  │                                                                        │   │
│  │  ┌──────────────────┐   ┌──────────────────┐   ┌─────────────────┐  │   │
│  │  │   Nginx (443)    │   │  aystek-updater   │   │  aystek-agent   │  │   │
│  │  │  Reverse Proxy   │   │  (Go binary)      │   │  (monitoring)   │  │   │
│  │  │  SSL Termination │   │  Update handler   │   │  Health checks  │  │   │
│  │  └────────┬─────────┘   └────────┬─────────┘   └────────┬────────┘  │   │
│  │           │                       │                        │           │   │
│  │  ┌────────┴─────────┐   ┌────────┴─────────┐             │           │   │
│  │  │  aystek-core     │   │  /opt/aystek-core  │            │           │   │
│  │  │  (Go API :9090)  │   │  releases/         │            │           │   │
│  │  │  systemd unit    │   │  current -> v1.2.0 │            │           │   │
│  │  └────────┬─────────┘   │  v1.1.0/           │            │           │   │
│  │           │              │  v1.2.0/           │            │           │   │
│  │  ┌────────┴─────────┐   └────────────────────┘            │           │   │
│  │  │  aystek-panel    │                                       │           │   │
│  │  │  (Next.js :3000) │   ┌────────────────────┐            │           │   │
│  │  │  systemd unit    │   │  MongoDB / Redis    │←───────────┘           │   │
│  │  └──────────────────┘   │  /var/lib/aystek    │                        │   │
│  │                          └────────────────────┘                        │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Bileşenler

| Bileşen | Rol | Dil/Teknoloji |
|---|---|---|
| `aystek-core` | Go REST API | Go binary, systemd |
| `aystek-panel` | Next.js admin UI | Node.js, systemd veya static |
| `aystek-updater` | Update orchestrator | Go binary, systemd (opsiyonel daemon) |
| `aystek-agent` | Health + metrics + log forwarder | Go binary, systemd |
| `aystek-cli` | Sysadmin CLI aracı | Go binary, /usr/local/bin |
| `nginx` | Reverse proxy + SSL | rpm paketi |
| Update Server | Manifest + paketi sunan API | Harici (AYSTEK HQ yönetiminde) |

---

## 3. Neden Bu Mimari

### Kararın Özü

**"Symlink-based release + tarball + installer + systemd + aystek-updater daemon"** kombinasyonu seçilmiştir. Bu kararın arkasındaki 5 temel gerekçe:

1. **Bağımlılık yok:** Docker gerektirmiyor. CentOS 9 minimal kurulumda bile çalışır.
2. **Operatör özgürlüğü:** Sysadmin tanıdık araçlarla (systemctl, ls, ln) her şeyi anlayabilir.
3. **Rollback deterministik:** `/opt/aystek-core/current` symlink'i bir önceki sürüme döndürmek 1 komut.
4. **Audit trail:** Her dosya değişikliği, her symlink değişimi kayıt altına alınır.
5. **CI/CD entegrasyon kolaylığı:** Artifact = tarball + manifest JSON. Herhangi bir CI sistemiyle üretilebilir.

### Reddetme Gerekçeleri

- **RPM:** Build pipeline'ı complex, `%pre/%post` scriptler test edilmesi güç, rollback native RPM'de zayıf.
- **Docker-only:** Hedef sunucular container runtime garantisi vermiyor; Kubernetes overhead fazla.
- **Ansible-only:** Configuration management aracı, deployment lifecycle yönetimi için yetersiz.

---

## 4. Alternatiflerin Karşılaştırması

| Kriter | RPM Paketi | Tarball + Installer | Container (Docker/K8s) | **Hibrit (Seçilen)** |
|---|---|---|---|---|
| Kurulum kolaylığı | Orta | Yüksek | Düşük (runtime gerekir) | **Çok Yüksek** |
| CentOS 9 uyumu | Yüksek | Yüksek | Orta | **Yüksek** |
| Rollback kalitesi | Zayıf | İyi | Çok İyi | **Çok İyi** |
| Zero-downtime | Zor | Mümkün | Kolay | **Mümkün** |
| Bağımlılık yönetimi | RPM deps | Manuel | Image içinde | **Entegre** |
| Operatör öğrenme eğrisi | Düşük | Düşük | Yüksek | **Düşük** |
| Panel entegrasyonu | Zor | Kolay | Orta | **Kolay** |
| Paket imzalama | GPG | SHA256+GPG | Docker Content Trust | **SHA256+GPG** |
| Uzun vadeli sürdürülebilirlik | Orta | Yüksek | Çok Yüksek | **Yüksek** |
| Öneri | ❌ | ✅ Tek başına yeterli | ❌ (opsiyonel katman) | ✅✅ **SEÇİLDİ** |

**Varsayım:** Müşteri sunucularında Docker yoktur ve zorunlu kılınamaz. Container desteği gelecek versiyonda opsiyonel katman olarak eklenecektir.

---

## 5. Deployment Topology

### Tek Sunucu (Standart Kurulum)

```
[CentOS 9 Server]
  ├── Nginx (443 → :9090 proxy, /panel → :3000 proxy)
  ├── aystek-core.service    (Go API, port 9090)
  ├── aystek-panel.service   (Next.js, port 3000)
  ├── aystek-agent.service   (Background health reporter)
  ├── mongod.service         (MongoDB 6.x)
  └── redis.service          (Redis 7.x)
```

### Çok Sunucu (Horizontal Scale, Gelecek Faz)

```
[Load Balancer / HAProxy]
    ├── [App Server 1] aystek-core + aystek-panel
    ├── [App Server 2] aystek-core + aystek-panel
    └── [DB Server]    MongoDB replicaset + Redis sentinel
```

### Update Flow Topology

```
AYSTEK Update Server
    └──[HTTPS]──► aystek-updater (müşteri sunucusu)
                      ├── manifest kontrol
                      ├── paket indir
                      ├── verify (SHA256 + GPG)
                      ├── backup al
                      ├── maintenance mode aç
                      ├── migration çalıştır
                      ├── symlink switch
                      ├── servis restart
                      ├── health check
                      └── maintenance mode kapat / rollback
```

---

## 6. Paketleme Stratejisi

### 6.1 Paket Formatı

**Format:** `.tar.gz` + signed SHA256 manifest

Her release için üretilen dosyalar:

```
aystek-core-v1.2.0-linux-amd64.tar.gz
aystek-core-v1.2.0-linux-amd64.tar.gz.sha256
aystek-core-v1.2.0-linux-amd64.tar.gz.sig   (GPG imza)
aystek-core-v1.2.0.manifest.json
```

### 6.2 Tarball İçeriği

```
aystek-core-v1.2.0/
├── manifest.json               # Release metadata
├── bin/
│   ├── aystek-core             # Go API binary
│   ├── aystek-updater          # Updater binary
│   └── aystek-agent            # Agent binary
├── panel/
│   ├── .next/                  # Next.js production build
│   ├── public/
│   └── package.json
├── migrations/
│   ├── 0001_initial.js
│   └── 0042_add_telematics.js
├── scripts/
│   ├── install.sh
│   ├── upgrade.sh
│   ├── rollback.sh
│   ├── pre-install.sh
│   ├── post-install.sh
│   ├── pre-upgrade.sh
│   └── post-upgrade.sh
├── config/
│   ├── aystek-core.env.example
│   ├── nginx.conf.template
│   └── systemd/
│       ├── aystek-core.service
│       ├── aystek-panel.service
│       └── aystek-agent.service
└── public/
    └── static/                 # Static web assets
```

### 6.3 Enterprise Dizin Yapısı (Sunucu)

```
/opt/aystek-core/
├── releases/
│   ├── v1.0.0/                 # Arşivlenmiş eski release
│   ├── v1.1.0/                 # Bir önceki release (rollback için saklanır)
│   └── v1.2.0/                 # Aktif release
├── current -> releases/v1.2.0  # Symlink (atomic switch ile değişir)
├── previous -> releases/v1.1.0 # Rollback hedefi
├── staging/                    # Yeni release hazırlık alanı
└── shared/
    ├── uploads/                # Tenant upload dosyaları (sürümler arası paylaşılır)
    └── cache/

/etc/aystek-core/
├── core.env                    # Ana environment dosyası (root:aystek 640)
├── panel.env                   # Panel environment
├── update.conf                 # Updater konfigürasyonu
└── license.key                 # Lisans dosyası (şifreli)

/var/log/aystek-core/
├── core/
│   ├── app.log
│   └── error.log
├── panel/
│   └── app.log
├── updater/
│   ├── update.log
│   └── audit.log
├── agent/
│   └── health.log
└── migrations/
    └── migration.log

/var/lib/aystek-core/
├── data/                       # Kalıcı uygulama verisi
├── backups/
│   ├── pre-update-v1.2.0/      # Update öncesi backup
│   └── scheduled/
└── migrations/
    └── state.json              # Migration state tracker

/var/backups/aystek-core/       # Sistem seviyesi backup
├── db/
├── config/
└── releases/

/tmp/aystek-core/               # Geçici dosyalar (update staging)
```

### 6.4 systemd Servis Planı

**aystek-core.service** (Go API):
- `After=network-online.target mongod.service redis.service`
- `Restart=always`, `RestartSec=5`
- `User=aystek`, `Group=aystek`
- `WorkingDirectory=/opt/aystek-core/current`
- `ExecStart=/opt/aystek-core/current/bin/aystek-core`
- `EnvironmentFile=/etc/aystek-core/core.env`
- Hardening: `PrivateTmp=true`, `NoNewPrivileges=true`, `ProtectSystem=strict`, `ReadWritePaths=/var/log/aystek-core /var/lib/aystek-core`

**aystek-panel.service** (Next.js):
- `After=aystek-core.service`
- `ExecStart=/usr/local/bin/node /opt/aystek-core/current/panel/server.js`
- `EnvironmentFile=/etc/aystek-core/panel.env`

**aystek-agent.service**:
- Timer tabanlı veya daemon: 30 saniyede bir health check → `/var/lib/aystek-core/health.json`

### 6.5 Installer Scripts

| Script | Tetiklenme | Görev |
|---|---|---|
| `pre-install.sh` | İlk kurulum öncesi | OS uyumluluk, bağımlılık, disk alanı, kullanıcı oluşturma |
| `post-install.sh` | İlk kurulum sonrası | systemd enable, env oluşturma, lisans doğrulama, ilk seed |
| `pre-upgrade.sh` | Güncelleme öncesi | Backup al, maintenance mode aç, queue drain, preflight |
| `post-upgrade.sh` | Güncelleme sonrası | Migration çalıştır, symlink güncelle, servis restart, health check |
| `rollback.sh` | Rollback tetiklendiğinde | Önceki symlink'e dön, servis restart, DB rollback (opsiyonel) |
| `uninstall.sh` | Kaldırma | Servis durdur, veri yedekle, dizin temizle |

---

## 7. Sürümleme ve Release Yönetimi

### 7.1 Semantic Versioning

Format: `MAJOR.MINOR.PATCH[-CHANNEL][+BUILD]`

```
1.2.0              → stable release
1.2.1-rc.1         → release candidate
1.3.0-beta.2       → beta
2.0.0-internal.1   → internal/nightly
1.2.0+20260331     → build metadata
```

| Segment | Anlam |
|---|---|
| MAJOR | Breaking change (API veya schema kırıcı değişiklik) |
| MINOR | Geriye dönük uyumlu yeni özellik |
| PATCH | Bug fix, güvenlik yaması |
| CHANNEL | stable / rc / beta / internal |

### 7.2 Release Channels

| Kanal | Kitle | Auto-update? | Approval gerekir? |
|---|---|---|---|
| `stable` | Production müşterileri | Hayır (manuel) | Evet |
| `rc` | Erken benimseyenler | Hayır | Evet |
| `beta` | Test ortamları | Opsiyonel | Evet |
| `internal` | AYSTEK dev servers | Evet | Hayır |

### 7.3 Release Manifest Örneği

```json
{
  "version": "1.2.0",
  "channel": "stable",
  "released_at": "2026-03-31T10:00:00Z",
  "build_id": "20260331-abc1234",
  "components": {
    "backend": {
      "version": "1.2.0",
      "binary": "aystek-core",
      "checksum_sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4..."
    },
    "frontend": {
      "version": "1.2.0",
      "build_hash": "abc1234def567",
      "checksum_sha256": "a87ff679a2f3e71d9181a67b7542122c..."
    },
    "updater": {
      "version": "1.1.0",
      "min_version": "1.0.0"
    }
  },
  "compatibility": {
    "min_backend_version": "1.1.0",
    "min_frontend_version": "1.1.0",
    "min_updater_version": "1.0.0",
    "schema_version": 42,
    "migration_version": 42,
    "breaking_changes": false
  },
  "migrations": [
    {
      "id": "0042_add_telematics_commands",
      "description": "Adds telematics commands collection indexes",
      "reversible": true,
      "estimated_duration_seconds": 5,
      "affects_collections": ["telematics_commands"],
      "breaking": false
    }
  ],
  "release_notes": {
    "tr": "## v1.2.0\n\n### Yenilikler\n...",
    "en": "## v1.2.0\n\n### What's New\n..."
  },
  "download_url": "https://releases.aystek.com/stable/v1.2.0/aystek-core-v1.2.0-linux-amd64.tar.gz",
  "signature": "...",
  "rollback_supported": true,
  "requires_maintenance": false,
  "estimated_downtime_seconds": 0,
  "preflight_checks": [
    "disk_space_min_500mb",
    "mongodb_connection",
    "redis_connection"
  ]
}
```

### 7.4 Compatibility Matrix

| Backend | Frontend | Schema | Migration | Min Updater | Kanal |
|---|---|---|---|---|---|
| 1.0.x | 1.0.x | 1–10 | 1–10 | 1.0.0 | stable |
| 1.1.x | 1.1.x | 11–30 | 11–30 | 1.0.0 | stable |
| 1.2.x | 1.2.x | 31–42 | 31–42 | 1.1.0 | stable |
| 2.0.x | 2.0.x | 43+ | 43+ | 2.0.0 | stable (breaking) |

### 7.5 Politikalar

**Upgrade Policy:**  
- PATCH: Doğrudan uygulanabilir, bakım modu gerekmez  
- MINOR: Opsiyonel bakım modu, migration otomatik  
- MAJOR: Bakım modu zorunlu, DBA onayı, rollback planı hazır olmalı  

**Rollback Policy:**  
- Her başarılı update → önceki release 30 gün saklanır  
- Rollback penceresi: PATCH için 7 gün, MINOR için 30 gün, MAJOR için 90 gün  
- DB rollback: Reversible migration varsa desteklenir; yoksa "forward-only" uyarısı verilir  

**Deprecation Policy:**  
- API endpoint deprecation: 2 MINOR sürüm önceden `Sunset` header ile bildirim  
- Config key deprecation: Changelog'da belgelenir, bir MAJOR sürüm boyunca eski key desteklenir  

---

## 8. Update / Rollback Akış Diyagramı

### 8.1 Normal Update Akışı

```
[Tetikleyici: Panel / CLI / Schedule / API]
        │
        ▼
[1. PREFLIGHT CHECKS]
    ├── Disk alanı ≥ 500MB?
    ├── Mevcut versiyon uyumlu mu?
    ├── Lisans geçerli mi?
    ├── MongoDB bağlantısı var mı?
    └── Açık işlemler / queue var mı?
        │
        ▼ (tümü OK)
[2. MANIFEST FETCH]
    ├── Update Server'dan manifest.json çek
    ├── GPG imzasını doğrula
    └── Version compatibility kontrolü
        │
        ▼ (uyumlu)
[3. PACKAGE DOWNLOAD]
    ├── Tarball'ı /tmp/aystek-core/ 'a indir
    ├── SHA256 checksum doğrula
    └── GPG imzasını doğrula
        │
        ▼ (doğrulandı)
[4. BACKUP]
    ├── Mevcut release'i /var/backups/aystek-core/pre-update-vX.Y.Z/ 'e kopyala
    └── MongoDB snapshot al (opsiyonel, config'e bağlı)
        │
        ▼
[5. STAGING]
    ├── Tarball'ı /opt/aystek-core/staging/ 'e aç
    ├── Migration dry-run (--check flag)
    └── Binary validate (version flag çalışıyor mu?)
        │
        ▼ (staging OK)
[6. MAINTENANCE MODE]
    ├── API: 503 + retry-after header dön
    └── Panel: "Güncelleme devam ediyor" ekranı göster
        │
        ▼
[7. QUEUE DRAIN]
    └── RabbitMQ worker'larını graceful shutdown (max 30sn)
        │
        ▼
[8. SERVICE STOP]
    └── systemctl stop aystek-core aystek-panel
        │
        ▼
[9. MIGRATION]
    ├── Pending migration'ları sırayla çalıştır
    ├── Her migration'ı logla
    └── Hata olursa → [ROLLBACK]
        │
        ▼ (migration OK)
[10. SYMLINK SWITCH] ← ATOMIC
    ├── previous -> current (eski release)
    ├── current -> releases/vNew
    └── chmod/chown düzelt
        │
        ▼
[11. SERVICE START]
    └── systemctl start aystek-core aystek-panel
        │
        ▼
[12. HEALTH CHECK] (max 60sn timeout)
    ├── /health endpoint → 200?
    ├── /health/ready → 200?
    └── Hata olursa → [ROLLBACK]
        │
        ▼ (healthy)
[13. MAINTENANCE MODE OFF]
        │
        ▼
[14. AUDIT LOG + NOTIFY]
    ├── Update log'a yaz
    ├── Panel'i bilgilendir (SSE/WebSocket)
    └── Alert sistemine bildir (opsiyonel)
        │
        ▼
[✅ UPDATE COMPLETE]
```

### 8.2 Rollback Akışı

```
[Tetikleyici: Otomatik (health fail) / Manuel (panel/CLI)]
        │
        ▼
[1. MAINTENANCE MODE] (zaten aktifse geç)
        │
        ▼
[2. SERVICE STOP]
        │
        ▼
[3. MIGRATION ROLLBACK] (reversible ise)
    ├── Migration state'den rollback aksiyonlarını belirle
    └── Geri al veya "forward-only" uyarısı ver
        │
        ▼
[4. SYMLINK REVERT] ← ATOMIC
    └── current -> previous
        │
        ▼
[5. SERVICE START]
        │
        ▼
[6. HEALTH CHECK]
    └── Hata olursa → Acil durum alertı + manual müdahale gerekiyor
        │
        ▼ (healthy)
[7. MAINTENANCE MODE OFF]
        │
        ▼
[8. AUDIT LOG + NOTIFY]
        │
        ▼
[✅ ROLLBACK COMPLETE]
```

---

## 9. Yönetim Paneli Modül Tasarımı

### 9.1 Modül Adı ve Konumu

**Go modül:** `internal/modules/system_ops/`  
**Next.js sayfaları:** `app/(dashboard)/admin/system/`  

### 9.2 Özellik Listesi

| # | Özellik | Öncelik |
|---|---|---|
| 1 | Mevcut sürüm ve bileşen versiyonları | P0 |
| 2 | Sunucu / node bilgileri (OS, RAM, Disk, Uptime) | P0 |
| 3 | Servis durumları (core, panel, agent, nginx) | P0 |
| 4 | Health check durumu (son kontrol zamanı) | P0 |
| 5 | Update kontrolü (son sürüm var mı?) | P0 |
| 6 | Update başlatma (tek tık) | P0 |
| 7 | Update ilerleme durumu (canlı log stream) | P0 |
| 8 | Rollback başlatma | P0 |
| 9 | Rollback geçmişi | P0 |
| 10 | Release notes görüntüleme | P1 |
| 11 | Bakım modu açma/kapama | P0 |
| 12 | Migration durumu ve geçmişi | P1 |
| 13 | Backup tetikleme ve liste | P1 |
| 14 | Servis log stream | P1 |
| 15 | Lisans/abonelik durumu | P1 |
| 16 | Audit log (tüm operasyon geçmişi) | P1 |
| 17 | Zamanlanmış update konfigürasyonu | P2 |
| 18 | Update bildirimleri ve politika | P2 |
| 19 | Multi-node görünüm (gelecek faz) | P3 |

### 9.3 Ekran Listesi

```
/admin/system/
├── /overview          → Dashboard: versiyon, sağlık, servis durumları
├── /updates           → Update kontrolü, başlatma, ilerleme
├── /updates/history   → Geçmiş update listesi
├── /rollback          → Rollback ekranı ve geçmiş
├── /migrations        → Migration durumu ve geçmişi
├── /maintenance       → Bakım modu yönetimi
├── /backups           → Backup listesi ve tetikleme
├── /logs              → Canlı log stream
├── /health            → Detaylı health check sonuçları
├── /node              → Sunucu/node bilgileri
├── /license           → Lisans ve abonelik durumu
└── /audit             → Operasyon audit log
```

### 9.4 Permission / RBAC Planı

| Permission | Açıklama | Super Admin | Admin | Operator | Viewer |
|---|---|---|---|---|---|
| `system_ops:read` | Tüm ekranları görüntüle | ✅ | ✅ | ✅ | ✅ |
| `system_ops:update_start` | Update başlat | ✅ | ✅ | ❌ | ❌ |
| `system_ops:rollback` | Rollback başlat | ✅ | ✅ | ❌ | ❌ |
| `system_ops:maintenance` | Bakım modu yönet | ✅ | ✅ | ✅ | ❌ |
| `system_ops:backup` | Backup tetikle | ✅ | ✅ | ✅ | ❌ |
| `system_ops:config` | Sistem config güncelle | ✅ | ❌ | ❌ | ❌ |
| `system_ops:license` | Lisans yönet | ✅ | ❌ | ❌ | ❌ |

**Kritik operasyonlar (update/rollback) için 2-step confirmation zorunludur.**

### 9.5 UI Component Listesi

| Component | Amaç |
|---|---|
| `SystemOverviewCard` | Versiyon, sağlık, uptime özeti |
| `ServiceStatusWidget` | Her servisi badge ile göster |
| `VersionBadge` | Mevcut vs son sürüm karşılaştırması |
| `UpdateWizard` | Adım adım update başlatma modal |
| `UpdateProgressStream` | Canlı log stream (SSE) |
| `RollbackConfirmModal` | 2-step rollback onayı |
| `MaintenanceModeToggle` | Bakım modu açma/kapatma + mesaj |
| `MigrationStatusTable` | Migration listesi ve durumları |
| `BackupListTable` | Backup listesi + restore butonu |
| `LogStreamViewer` | Canlı log (SSE + filtre) |
| `NodeInfoPanel` | OS, CPU, RAM, disk bilgileri |
| `LicenseStatusCard` | Lisans geçerlilik, abonelik |
| `AuditLogTable` | Operasyon geçmişi |
| `ReleaseNoteViewer` | Markdown render |
| `HealthCheckDetail` | Bileşen bazlı health sonuçları |

---

## 10. Go Backend API Tasarımı

### 10.1 Route Grubu

```
/admin/system/
├── GET    /health                  → Detaylı health check
├── GET    /health/ready            → Readiness (servis hazır mı?)
├── GET    /health/live             → Liveness
├── GET    /version                 → Mevcut versiyon bilgisi
├── GET    /node                    → Sunucu node bilgileri
├── GET    /services                → Tüm servis durumları
│
├── GET    /updates/check           → Update var mı? manifest kontrol
├── GET    /updates                 → Update geçmişi
├── POST   /updates/start           → Update başlat (async job tetikler)
├── GET    /updates/:id             → Belirli update'in durumu
├── GET    /updates/:id/logs        → Update log stream (SSE)
│
├── GET    /rollbacks               → Rollback geçmişi + mevcut rollback hedefleri
├── POST   /rollbacks/start         → Rollback başlat
├── GET    /rollbacks/:id           → Rollback durumu
│
├── GET    /migrations              → Migration listesi ve durumları
├── POST   /migrations/run          → Manuel migration tetikle
│
├── GET    /maintenance             → Bakım modu durumu
├── POST   /maintenance/enable      → Bakım modunu aç
├── POST   /maintenance/disable     → Bakım modunu kapat
│
├── GET    /backups                 → Backup listesi
├── POST   /backups/trigger         → Manuel backup başlat
├── POST   /backups/:id/restore     → Restore (kritik: 2-step)
│
├── GET    /logs                    → Log stream (SSE, filtreli)
│
├── GET    /license                 → Lisans durumu
├── POST   /license/activate        → Lisans aktivasyonu
│
└── GET    /audit                   → Operasyon audit log listesi
```

### 10.2 Kritik Endpoint Detayları

**GET /admin/system/version**
```json
{
  "backend": {
    "version": "1.2.0",
    "build_id": "20260331-abc1234",
    "build_time": "2026-03-31T10:00:00Z",
    "commit": "abc1234",
    "go_version": "1.24.1"
  },
  "frontend": {
    "version": "1.2.0",
    "build_hash": "def567"
  },
  "updater": {
    "version": "1.1.0"
  },
  "schema_version": 42,
  "migration_version": 42
}
```

**POST /admin/system/updates/start** (işlem asenkron, job ID döner)
```json
{
  "target_version": "1.3.0",
  "channel": "stable",
  "options": {
    "take_backup": true,
    "maintenance_mode": true,
    "auto_rollback_on_health_fail": true,
    "dry_run": false
  }
}
```

Response:
```json
{
  "job_id": "upd_abc123",
  "status": "queued",
  "stream_url": "/admin/system/updates/upd_abc123/logs"
}
```

**GET /admin/system/updates/:id/logs** (SSE stream)
```
data: {"time":"2026-03-31T10:01:00Z","level":"info","step":"preflight","message":"Disk alanı kontrolü: 2.3GB mevcut ✅"}
data: {"time":"2026-03-31T10:01:01Z","level":"info","step":"download","message":"Paket indiriliyor: 45%"}
...
data: {"time":"2026-03-31T10:03:00Z","level":"success","step":"complete","message":"Update tamamlandı: v1.3.0"}
```

### 10.3 Versiyon Bilgisini Binary'ye Gömme (Go)

```go
// cmd/api/version.go
var (
    Version    = "dev"
    BuildTime  = "unknown"
    CommitHash = "unknown"
    GoVersion  = runtime.Version()
)
// Build flags: -ldflags "-X main.Version=1.2.0 -X main.BuildTime=..."
```

---

## 11. Next.js UI Modül Tasarımı

### 11.1 Sayfa Yapısı

```
app/(dashboard)/admin/system/
├── page.tsx                        → /overview (redirect)
├── overview/
│   └── page.tsx                    → SystemOverviewPage
├── updates/
│   ├── page.tsx                    → UpdatesPage (history + check)
│   └── [jobId]/
│       └── page.tsx                → UpdateJobDetailPage (canlı log)
├── rollback/
│   └── page.tsx                    → RollbackPage
├── migrations/
│   └── page.tsx                    → MigrationsPage
├── maintenance/
│   └── page.tsx                    → MaintenancePage
├── backups/
│   └── page.tsx                    → BackupsPage
├── logs/
│   └── page.tsx                    → LogStreamPage
├── node/
│   └── page.tsx                    → NodeInfoPage
├── license/
│   └── page.tsx                    → LicensePage
└── audit/
    └── page.tsx                    → AuditLogPage
```

### 11.2 API Client (`lib/api/system-ops.ts`)

```typescript
export const systemOpsApi = {
  getVersion():         Promise<VersionInfo>
  getHealth():          Promise<HealthStatus>
  getNodeInfo():        Promise<NodeInfo>
  getServices():        Promise<ServiceStatus[]>
  checkUpdates():       Promise<UpdateCheckResult>
  getUpdateHistory():   Promise<UpdateRecord[]>
  startUpdate(opts):    Promise<UpdateJob>
  getRollbacks():       Promise<RollbackOption[]>
  startRollback(opts):  Promise<RollbackJob>
  getMigrations():      Promise<MigrationRecord[]>
  getMaintenanceMode(): Promise<MaintenanceModeStatus>
  setMaintenanceMode(enabled, message): Promise<void>
  getBackups():         Promise<BackupRecord[]>
  triggerBackup():      Promise<BackupJob>
  getLicense():         Promise<LicenseStatus>
  getAuditLog(params):  Promise<PaginatedResult<AuditEntry>>
}
```

### 11.3 Canlı Log SSE Entegrasyonu

```typescript
// hooks/useUpdateLogStream.ts
export function useUpdateLogStream(jobId: string) {
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [status, setStatus] = useState<'running' | 'complete' | 'failed'>('running');

  useEffect(() => {
    const sse = new EventSource(`/admin/system/updates/${jobId}/logs`);
    sse.onmessage = (e) => {
      const entry = JSON.parse(e.data);
      setLogs(prev => [...prev, entry]);
      if (entry.step === 'complete') setStatus('complete');
      if (entry.level === 'error') setStatus('failed');
    };
    return () => sse.close();
  }, [jobId]);

  return { logs, status };
}
```

---

## 12. Veri Modeli / Koleksiyon Önerileri

### 12.1 MongoDB Koleksiyonları

**`system_updates`** — Update kayıtları
```javascript
{
  _id: ObjectId,
  job_id: "upd_abc123",
  tenant_id: ObjectId,        // Multi-tenant: hangi tenant için
  type: "update" | "rollback",
  status: "queued" | "running" | "success" | "failed" | "rolled_back",
  from_version: "1.1.0",
  to_version: "1.2.0",
  triggered_by: ObjectId,     // User ID
  trigger_type: "manual" | "scheduled" | "api" | "auto",
  options: { take_backup, maintenance_mode, auto_rollback },
  steps: [
    { name: "preflight", status, started_at, completed_at, output }
  ],
  started_at: ISODate,
  completed_at: ISODate,
  duration_ms: 12300,
  error: "...",
  rollback_available: true,
  created_at: ISODate
}
```

**`system_migrations`** — Migration kayıtları
```javascript
{
  _id: ObjectId,
  migration_id: "0042_add_telematics",
  version: "1.2.0",
  status: "pending" | "running" | "success" | "failed" | "rolled_back",
  direction: "up" | "down",
  started_at: ISODate,
  completed_at: ISODate,
  duration_ms: 500,
  reversible: true,
  affects_collections: ["telematics_commands"],
  error: ""
}
```

**`system_backups`** — Backup kayıtları
```javascript
{
  _id: ObjectId,
  backup_id: "bak_xyz789",
  type: "pre_update" | "scheduled" | "manual",
  trigger_version: "1.2.0",
  status: "running" | "success" | "failed",
  path: "/var/backups/aystek-core/pre-update-v1.2.0",
  size_bytes: 52428800,
  components: ["binary", "config", "db_snapshot"],
  created_at: ISODate,
  expires_at: ISODate
}
```

**`system_audit_log`** — Audit kayıtları
```javascript
{
  _id: ObjectId,
  action: "update.started" | "rollback.started" | "maintenance.enabled" | ...,
  actor_id: ObjectId,
  actor_email: "admin@tenant.com",
  ip_address: "1.2.3.4",
  resource: "update:upd_abc123",
  details: { from_version, to_version, options },
  outcome: "success" | "failed" | "pending",
  created_at: ISODate
}
```

**`system_health_snapshots`** — TTL collection (son 24 saat)
```javascript
{
  _id: ObjectId,
  snapshot_at: ISODate,
  services: { core: "up", panel: "up", mongodb: "up", redis: "up" },
  metrics: { cpu_pct: 12, mem_pct: 45, disk_pct: 33 },
  response_time_ms: 45,
  ttl_expire: ISODate  // TTL index, 24 saat
}
```

---

## 13. Dizin Yapısı ve systemd Servis Planı

### 13.1 Modül Dizin Yapısı (Go)

```
internal/modules/system_ops/
├── model/
│   ├── update.go              # UpdateRecord, UpdateJob, UpdateStep
│   ├── migration.go           # MigrationRecord
│   ├── backup.go              # BackupRecord
│   ├── health.go              # HealthStatus, ServiceStatus
│   └── audit.go               # AuditEntry
├── dto/
│   └── system_ops.go          # Request/Response DTOs
├── repository/
│   ├── interface.go
│   └── mongo.go
├── service/
│   ├── service.go             # Ana orchestration
│   ├── health_service.go      # Health check logic
│   ├── update_service.go      # Update flow orchestration
│   ├── migration_service.go   # Migration runner
│   ├── backup_service.go      # Backup logic
│   └── service_test.go
├── handler/
│   ├── admin.go               # Admin HTTP handlers
│   └── health.go              # Public health endpoint
├── worker/
│   └── update_worker.go       # Async update job runner
└── routes.go
```

### 13.2 Updater Binary Dizin Yapısı (Go)

```
cmd/updater/
├── main.go                    # Updater entry point
├── preflight.go               # Preflight check runner
├── downloader.go              # Manifest + package fetch
├── verifier.go                # Checksum + GPG verify
├── migrator.go                # Migration runner
├── switcher.go                # Symlink switch atomic
└── rollback.go                # Rollback logic
```

### 13.3 CLI Binary Dizin Yapısı

```
cmd/cli/
├── main.go                    # cobra root
├── cmd/
│   ├── update.go              # aystek-cli update [--channel stable]
│   ├── rollback.go            # aystek-cli rollback [--version 1.1.0]
│   ├── status.go              # aystek-cli status
│   ├── maintenance.go         # aystek-cli maintenance [on|off]
│   ├── migrate.go             # aystek-cli migrate [run|status|rollback]
│   ├── backup.go              # aystek-cli backup [create|list|restore]
│   └── version.go             # aystek-cli version
```

### 13.4 Scripts Dizini

```
scripts/deploy/
├── install.sh
├── upgrade.sh
├── rollback.sh
├── uninstall.sh
├── backup.sh
├── health-check.sh
├── maintenance-on.sh
├── maintenance-off.sh
└── lib/
    ├── common.sh              # Shared functions
    ├── logger.sh              # Colored logging
    └── systemd.sh             # systemd helpers
```

---

## 14. Güvenlik ve Audit Planı

### 14.1 Paket Güvenliği

| Katman | Mekanizma |
|---|---|
| Paket bütünlüğü | SHA256 checksum |
| Paket imzalama | GPG detached signature (AYSTEK özel key) |
| Transport | HTTPS + sertifika pinleme (update client'ta) |
| Supply chain | SBOM (CycloneDX format), Go module proxy |
| Vulnerability scan | `govulncheck` + Trivy (CI pipeline'da) |

### 14.2 Sistem Güvenliği

| Katman | Uygulama |
|---|---|
| Service user | `aystek` system user, loginShell=/sbin/nologin |
| File permissions | Binary: 750, Config: 640 owner=root:aystek, Logs: 640 |
| systemd hardening | `NoNewPrivileges=true`, `ProtectSystem=strict`, `PrivateTmp=true`, `CapabilityBoundingSet=` |
| SELinux | Custom policy module: `aystek_core.te` |
| Firewall | nftables/firewalld: sadece 443 ve 9090 (internal) |
| Secret management | `.env` at-rest şifreleme veya HashiCorp Vault entegrasyonu |

### 14.3 API Güvenliği

| Kural | Uygulama |
|---|---|
| Role-based auth | Super Admin + Admin = update/rollback; Operator = maintenance |
| 2-step confirmation | Update/rollback için `confirm_token` mekanizması |
| Rate limiting | Update start: 2 req/30sn; CLI token: 5 req/sn |
| Audit log | Tüm sistem operasyonu audit_log'a yazılır |
| Update token | CLI/API için ayrı kısa ömürlü token |

### 14.4 Operasyon Güvenliği

| Kural | Açıklama |
|---|---|
| Approval workflow | MAJOR update için ikinci admin onayı |
| Rollback lock | Bir rollback devam ederken ikincisi engellensin |
| Concurrent update lock | Eş zamanlı update engellensin (distributed lock: Redis) |
| Anti-tamper | Binary hash boot sırasında doğrulanır (opsiyonel) |
| Download retry limit | Başarısız indirme maks 3 deneme |
| Timeout | Her update adımı için timeout (preflight: 30sn, download: 300sn, migration: custom) |

---

## 15. CI/CD ve Release Pipeline Planı

### 15.1 Branch Stratejisi

```
main            → production-ready kod (tagged releases)
develop         → aktif geliştirme
release/vX.Y.Z  → release hazırlık branch
feature/*       → özellik branchleri
fix/*           → bug fix branchleri
hotfix/*        → kritik fix (main'e direkt merge)
```

### 15.2 Release Pipeline

```
[Git Tag: v1.2.0]
        │
        ▼
[1. BUILD GATE]
    ├── go test ./...
    ├── npm run test (panel)
    ├── go vet + staticcheck
    └── govulncheck
        │
        ▼ (tümü OK)
[2. SECURITY SCAN]
    ├── Trivy (Go binary + panel)
    ├── SBOM üret (CycloneDX)
    └── License check
        │
        ▼
[3. ARTIFACT BUILD]
    ├── go build -ldflags "-X Version=1.2.0 ..."
    ├── npm run build (Next.js production)
    ├── Tarball oluştur
    ├── Manifest JSON üret
    └── SHA256 checksum üret
        │
        ▼
[4. PACKAGE SIGN]
    └── GPG sign (AYSTEK CI signing key, HSM veya secrets vault)
        │
        ▼
[5. STAGE PUBLISH]
    ├── Staging update server'a yükle
    └── Internal test trigger
        │
        ▼ (internal test OK)
[6. MANUAL APPROVAL GATE]
    └── Release Manager onayı
        │
        ▼
[7. PRODUCTION PUBLISH]
    ├── Update server'a yükle (stable channel)
    ├── Manifest API güncelle
    └── Changelog publish
        │
        ▼
[8. NOTIFY]
    ├── Müşteri portalına bildirim
    └── Slack/e-posta release duyurusu
```

### 15.3 Version Bump Stratejisi

```bash
# Otomatik PATCH bump (fix/ merger)
git tag -a v$(semver bump patch) -m "Release v$(semver bump patch)"

# Manuel MINOR/MAJOR (release/ merger)
# release/v1.2.0 → main merge → tag v1.2.0
```

### 15.4 Rollback Artifact Retention

- Son 3 stable release → sonsuza kadar saklanır  
- Beta/RC → 90 gün saklanır  
- Internal → 30 gün saklanır  

---

## 16. Fazlara Bölünmüş Implementasyon Roadmap

### Faz Özeti

| Faz | Ad | Süre | Öncelik |
|---|---|---|---|
| **F1** | Altyapı Temeli | 1 hafta | P0 |
| **F2** | Paketleme + Installer | 1 hafta | P0 |
| **F3** | Updater Binary + CLI | 1 hafta | P0 |
| **F4** | system_ops Go Modülü | 1 hafta | P0 |
| **F5** | Panel UI Modülü | 1 hafta | P0 |
| **F6** | Maintenance Mode | 3 gün | P0 |
| **F7** | Migration Altyapısı | 1 hafta | P1 |
| **F8** | Backup Sistemi | 3 gün | P1 |
| **F9** | Update Server API | 1 hafta | P1 |
| **F10** | CI/CD Pipeline | 1 hafta | P1 |
| **F11** | Güvenlik Sertleştirme | 3 gün | P2 |
| **F12** | Gözlemlenebilirlik | 3 gün | P2 |

---

## 17. Her Faz İçin Görev Listesi

### F1: Altyapı Temeli

- [ ] `/opt/aystek-core/` dizin yapısı şeması finalize
- [ ] `aystek` sistem kullanıcısı oluşturma scripti
- [ ] systemd unit dosyaları (core, panel, agent) temel versiyonları
- [ ] `/etc/aystek-core/core.env.example` oluştur
- [ ] Dizin permission matrisini uygula
- [ ] Graceful shutdown (mevcut `signal.NotifyContext` kodunu doğrula ✅ — zaten yapıldı)
- [ ] Version embed (ldflags): `cmd/api/version.go`
- [ ] `/health` ve `/health/ready` endpointleri (mevcut varsa doğrula)

### F2: Paketleme + Installer

- [ ] `scripts/build/package.sh` — tarball üretim scripti
- [ ] `scripts/build/sign.sh` — GPG imzalama
- [ ] `scripts/build/manifest.sh` — manifest.json üretimi
- [ ] `scripts/deploy/install.sh` — ilk kurulum
- [ ] `scripts/deploy/upgrade.sh` — güncelleme scripti
- [ ] `scripts/deploy/rollback.sh` — rollback scripti
- [ ] `scripts/deploy/lib/common.sh` — paylaşılan fonksiyonlar
- [ ] Symlink yönetim kritik: `ln -sfn` atomic switch
- [ ] Next.js standalone build konfigürasyonu (`next.config.js`)

### F3: Updater Binary + CLI

- [ ] `cmd/updater/main.go` — cobra tabanlı updater
- [ ] `cmd/updater/preflight.go` — disk, bağlantı, uyumluluk kontrolleri
- [ ] `cmd/updater/downloader.go` — manifest fetch + tarball download
- [ ] `cmd/updater/verifier.go` — SHA256 + GPG verify
- [ ] `cmd/updater/switcher.go` — atomic symlink switch
- [ ] `cmd/cli/main.go` — sysadmin CLI
- [ ] aystek-agent temel health reporter
- [ ] systemd unit için aystek-updater.service

### F4: system_ops Go Modülü

- [ ] `internal/modules/system_ops/model/` — tüm modeller
- [ ] `internal/modules/system_ops/repository/interface.go`
- [ ] `internal/modules/system_ops/repository/mongo.go`
- [ ] `internal/modules/system_ops/service/service.go`
- [ ] `internal/modules/system_ops/service/health_service.go`
- [ ] `internal/modules/system_ops/service/update_service.go`
- [ ] `internal/modules/system_ops/handler/admin.go`
- [ ] `internal/modules/system_ops/handler/health.go`
- [ ] `internal/modules/system_ops/routes.go`
- [ ] `cmd/api/main.go`'ya wiring
- [ ] SSE log stream endpoint
- [ ] Update job async worker
- [ ] translation_seed.go: system_ops çevirileri (28 dil)
- [ ] menu_seed.go: System / Operations menü grubu
- [ ] permission_seed.go: system_ops:* permissions
- [ ] role_seed.go: editor ve admin atamaları
- [ ] `internal/seeds/system_ops_seed.go`

### F5: Panel UI Modülü

- [ ] `lib/api/system-ops.ts` API client
- [ ] `lib/types/system-ops.ts` TypeScript tipleri
- [ ] `app/(dashboard)/admin/system/overview/page.tsx`
- [ ] `app/(dashboard)/admin/system/updates/page.tsx`
- [ ] `app/(dashboard)/admin/system/updates/[jobId]/page.tsx`
- [ ] `app/(dashboard)/admin/system/rollback/page.tsx`
- [ ] `app/(dashboard)/admin/system/maintenance/page.tsx`
- [ ] `app/(dashboard)/admin/system/health/page.tsx`
- [ ] `app/(dashboard)/admin/system/node/page.tsx`
- [ ] `app/(dashboard)/admin/system/audit/page.tsx`
- [ ] `hooks/useUpdateLogStream.ts`
- [ ] UI components (UpdateWizard, RollbackConfirmModal, vb.)

### F6: Maintenance Mode

- [ ] `internal/platform/middleware/maintenance.go` — mevcut varsa geliştir
- [ ] Maintenance mode Redis key: `system:maintenance_mode`
- [ ] Panel UI: MaintenanceModeToggle component
- [ ] Update sırasında otomatik aktifleşme hook'u
- [ ] Maintenance mesajı özelleştirme (tenant bazlı)

### F7: Migration Altyapısı

- [ ] `internal/platform/migrations/runner.go` — migration runner
- [ ] `internal/platform/migrations/registry.go` — migration kaydı
- [ ] Migration state collection (`system_migrations`)
- [ ] Dry-run modu
- [ ] Reversible migration interface
- [ ] Migration CLI komutu: `aystek-cli migrate`

### F8: Backup Sistemi

- [ ] `scripts/deploy/backup.sh` — binary + config + db backup
- [ ] MongoDB dump entegrasyonu (mongodump)
- [ ] Backup kayıt koleksiyonu (`system_backups`)
- [ ] Otomatik pre-update backup hook
- [ ] Panel UI: backup liste + tetikle

### F9: Update Server API (Harici / AYSTEK HQ)

- [ ] Manifest storage + serving API spesifikasyonu
- [ ] License validation API spesifikasyonu
- [ ] Signed artifact hosting stratejisi (S3 + CloudFront veya eşdeğeri)
- [ ] Version compatibility API

### F10: CI/CD Pipeline

- [ ] `.github/workflows/release.yml` (veya GitLab CI eşdeğeri)
- [ ] Build + test gate
- [ ] Security scan (govulncheck + Trivy)
- [ ] SBOM üretimi (CycloneDX)
- [ ] Package build + sign
- [ ] Staging publish
- [ ] Manual approval gate
- [ ] Production publish
- [ ] Changelog otomasyonu (git-cliff veya conventional-changelog)

### F11: Güvenlik Sertleştirme

- [ ] systemd hardening direktifleri finalize
- [ ] SELinux policy modülü (`aystek_core.te`)
- [ ] File permission audit scripti
- [ ] Update server sertifika pinleme
- [ ] 2-step confirmation implementation
- [ ] Concurrent update/rollback lock (Redis)
- [ ] Anti-tamper binary hash check (opsiyonel)

### F12: Gözlemlenebilirlik

- [ ] Prometheus metrics: update duration, rollback count, migration duration
- [ ] health snapshot TTL collection
- [ ] Log aggregation konfigürasyonu (opsiyonel: Loki veya file-based)
- [ ] Alert tanımları: update failed, rollback triggered, maintenance mode stuck

---

## 18. Riskler ve Mitigasyonlar

| Risk | Olasılık | Etki | Mitigasyon |
|---|---|---|---|
| Migration sırasında veri kaybı | Düşük | Kritik | Pre-update backup zorunlu; reversible migration tercihi |
| Symlink switch sonrası boot fail | Düşük | Yüksek | Auto-rollback: health check fail → otomatik geri dön |
| Update Server erişim kaybı | Orta | Orta | Offline update modu: yerel paket kabul et |
| GPG key compromise | Çok Düşük | Kritik | Key rotation prosedürü; eski imzaları geçersiz kıl |
| Disk dolması (update sırasında) | Orta | Yüksek | Preflight: min 500MB disk kontrolü |
| Concurrent update fail | Orta | Orta | Redis distributed lock |
| MongoDB migration timeout | Orta | Yüksek | Long-running migration için background job stratejisi |
| Next.js build incompatibility | Düşük | Orta | Compatibility matrix; frontend version pinleme |
| Bakım modu açık kalırsa | Düşük | Yüksek | Max maintenance duration timer (varsayılan: 30 dakika) |
| Tenant verisi etkilenmesi | Düşük | Kritik | Migration tenant-aware test suite |

---

## 19. Acceptance Criteria

### F1-F3 (Altyapı + Paketleme + Updater)
- [ ] `install.sh` ile CentOS 9 minimal'a kurulum ≤ 10 dakika
- [ ] Tüm systemd servisleri `active (running)` durumuna geliyor
- [ ] `aystek-cli version` mevcut versiyonu basıyor
- [ ] Symlink switch atomic: yarım kalmış switch durumunda rollback çalışıyor

### F4-F5 (Go Modül + Panel UI)
- [ ] `/admin/system/health` 200 dönüyor ve tüm bileşenler gösteriliyor
- [ ] Panel üzerinden mevcut versiyon görünüyor
- [ ] Update check API sonucu panelde görünüyor
- [ ] Canlı log stream SSE çalışıyor

### F6-F8 (Maintenance + Migration + Backup)
- [ ] Maintenance mode açıkken `/api/*` 503 dönüyor
- [ ] Migration runner pending migration'ları çalıştırıyor
- [ ] Pre-update backup `/var/backups/aystek-core/`'a yazılıyor

### F9-F10 (Update Server + CI/CD)
- [ ] Tag push → artifact build → signed tarball otomatik üretiliyor
- [ ] `aystek-cli update` komutu latest stable versiyona güncelliyor

### F11-F12 (Güvenlik + Observability)
- [ ] `systemd-analyze security aystek-core.service` skoru ≥ 8/10
- [ ] Update/rollback işlemleri audit_log'a yazılıyor
- [ ] Prometheus endpoint update metrikleri sunuyor

### Genel
- [ ] Near-zero-downtime update: aystek-core servis kesintisi ≤ 5 saniye
- [ ] Rollback ≤ 60 saniyede tamamlanıyor
- [ ] Tüm kritik operasyonlar audit log'a yansıyor
- [ ] Health check fail → otomatik rollback tetikleniyor

---

## 20. Dosya Dosya Implementasyon Çalışma Planı

Implementasyona geçildiğinde her faz için aşağıdaki sıra izlenecektir. Her "devam et" komutunda bir sonraki grup tamamlanacak.

### Grup 1: Version Embed + Health (F1 başlangıcı)
```
cmd/api/version.go                    → Versiyon bilgisi struct + ldflags
internal/modules/system_ops/model/health.go  → HealthStatus, ServiceStatus model
internal/modules/system_ops/handler/health.go → /health, /health/ready handler
```

### Grup 2: systemd + Dizin Yapısı (F1 devam)
```
config/systemd/aystek-core.service    → Hardened systemd unit
config/systemd/aystek-panel.service
config/systemd/aystek-agent.service
scripts/deploy/lib/common.sh          → Paylaşılan bash fonksiyonları
scripts/deploy/lib/logger.sh          → Renkli logging
```

### Grup 3: Installer Scriptleri (F2)
```
scripts/deploy/install.sh
scripts/deploy/upgrade.sh
scripts/deploy/rollback.sh
scripts/build/package.sh
scripts/build/manifest.sh
```

### Grup 4: Updater Binary (F3)
```
cmd/updater/main.go
cmd/updater/preflight.go
cmd/updater/downloader.go
cmd/updater/verifier.go
cmd/updater/switcher.go
cmd/cli/main.go
cmd/cli/cmd/update.go
cmd/cli/cmd/status.go
cmd/cli/cmd/rollback.go
```

### Grup 5: system_ops Go Modülü (F4)
```
internal/modules/system_ops/model/update.go
internal/modules/system_ops/model/migration.go
internal/modules/system_ops/model/backup.go
internal/modules/system_ops/model/audit.go
internal/modules/system_ops/dto/system_ops.go
internal/modules/system_ops/repository/interface.go
internal/modules/system_ops/repository/mongo.go
internal/modules/system_ops/service/service.go
internal/modules/system_ops/service/health_service.go
internal/modules/system_ops/service/update_service.go
internal/modules/system_ops/handler/admin.go
internal/modules/system_ops/routes.go
internal/seeds/system_ops_seed.go
internal/seeds/translation_seed.go  (system_ops keys ekle)
internal/seeds/menu_seed.go         (System menü grubu ekle)
internal/seeds/permission_seed.go   (system_ops:* ekle)
cmd/api/main.go                     (system_ops wiring)
```

### Grup 6: Update Worker + SSE Stream (F4 devam)
```
internal/modules/system_ops/worker/update_worker.go
internal/modules/system_ops/service/migration_service.go
internal/modules/system_ops/service/backup_service.go
```

### Grup 7: Maintenance Mode Middleware (F6)
```
internal/platform/middleware/maintenance.go  (geliştir/oluştur)
```

### Grup 8: Panel UI Altyapısı (F5)
```
lib/types/system-ops.ts
lib/api/system-ops.ts
hooks/useUpdateLogStream.ts
hooks/useSystemHealth.ts
app/(dashboard)/admin/system/overview/page.tsx
app/(dashboard)/admin/system/updates/page.tsx
app/(dashboard)/admin/system/updates/[jobId]/page.tsx
app/(dashboard)/admin/system/rollback/page.tsx
app/(dashboard)/admin/system/maintenance/page.tsx
app/(dashboard)/admin/system/health/page.tsx
app/(dashboard)/admin/system/node/page.tsx
app/(dashboard)/admin/system/audit/page.tsx
```

### Grup 9: UI Components (F5 devam)
```
app/components/system/SystemOverviewCard.tsx
app/components/system/ServiceStatusWidget.tsx
app/components/system/UpdateWizard.tsx
app/components/system/UpdateProgressStream.tsx
app/components/system/RollbackConfirmModal.tsx
app/components/system/MaintenanceModeToggle.tsx
app/components/system/MigrationStatusTable.tsx
app/components/system/BackupListTable.tsx
app/components/system/LogStreamViewer.tsx
app/components/system/NodeInfoPanel.tsx
app/components/system/LicenseStatusCard.tsx
app/components/system/AuditLogTable.tsx
```

### Grup 10: CI/CD Pipeline (F10)
```
.github/workflows/release.yml
.github/workflows/security-scan.yml
scripts/build/sign.sh
scripts/build/sbom.sh
```

---

## Notlar ve Varsayımlar

**Varsayım 1:** AYSTEK platformu tek müşteri = tek kurulum (tenant/site mantığı uygulama içinde, OS/server seviyesinde değil). Multi-node kurulum gelecek faz.

**Varsayım 2:** Update Server (manifest API + artifact hosting) AYSTEK HQ tarafından ayrı bir altyapıda yönetilir. Bu planda specification yapılır, implementasyon ayrı servistir.

**Varsayım 3:** Müşteri sunucusunda internet erişimi var. Offline kurulum/güncelleme opsiyonel özellik olarak ileriki fazda ele alınacak.

**Varsayım 4:** MongoDB ve Redis kurulumu müşteri sorumluluğunda veya installer script içinde paketlenmiş. Bu plan uygulama katmanına odaklanır.

**Varsayım 5:** Next.js `output: "standalone"` ile build edilecek. Node.js runtime sunucuda mevcut olmalı. Alternatif: Nginx + static export (API route yoksa) — API routes varsa standalone zorunlu.

---

*Bu doküman implementasyon için onay beklemektedir. "Devam et" komutu ile Grup 1'den başlanacaktır.*
