---
Başlık: Deployment Senaryoları
Sürüm: 3.0.0
Durum: Üretime Hazır
Son Güncelleme: 2026-10-28
İlgili Belgeler:
  - ./SUNUCU_KURULUM_REHBERI.md
  - ./OPERATIONS_GUIDE.md
  - ../01-architecture/Distributed_Zero_Trust_Tenant_Architecture_Spec.md
  - ../06-modules/modules/system_ops.md
---

# Deployment Senaryoları

## İçindekiler

1. [Tek Tenant — Tek Site](#1-tek-tenant--tek-site)
2. [Tek Tenant — Çok Site](#2-tek-tenant--çok-site)
3. [Çok Tenant SaaS](#3-çok-tenant-saas)
4. [Güncelleme / Rollback Senaryoları](#4-güncelleme--rollback-senaryoları)
5. [Felaket Kurtarma (DR)](#5-felaket-kurtarma-dr)
6. [Yatay Ölçekleme (Gelecek Faz)](#6-yatay-ölçekleme-gelecek-faz)
7. [AYSTEK Haber Public-Site Senaryosu](#7-aystek-haber-public-site-senaryosu)
8. [2026-10-28 Onay politikaları ve tohum (seed) değişiklikleri](#8-2026-10-28-onay-politikaları-ve-tohum-seed-değişiklikleri)

---

## 1. Tek Tenant — Tek Site

**Kullanım:** Küçük/orta ölçekli işletmeler, tek marka, tek domain.

```
[Tek CentOS 9 Sunucusu]
  ├── Nginx (:443/:80)
  ├── aystek-core   (Go API :9090)
  ├── aystek-panel  (Next.js :3000)
  ├── MongoDB       (:27017 — dahili)
  ├── Redis         (:6379 — dahili)
  └── RabbitMQ      (:5672 — dahili)
```

**Özellikler:**
- Tüm bileşenler tek sunucuda çalışır
- systemd ile yönetilir
- `aystek-core-installer` ve `system_ops` birlikte lifecycle yönetimini üstlenir
- Yedek `/var/backups/aystek-core/` dizinine alınır

**Ölçek:** 10.000'e kadar günlük aktif kullanıcı, 100k/gün içerik işleme

**Kurulum:** [Sunucu Kurulum Rehberi](./SUNUCU_KURULUM_REHBERI.md)

---

## 2. Tek Tenant — Çok Site

**Kullanım:** Tek kurum, birden fazla web sitesi/marka (örn. medya grubu).

```
[Tek Sunucu]
  ├── aystek-core (Go API)
  │     ├── Tenant: "anamedya"
  │     │     ├── Site: "gazetex.com"
  │     │     ├── Site: "sporx.com"
  │     │     └── Site: "ekonomix.com"
  │     └── (site başına izole ayarlar)
  └── Nginx (sanal host başına domain routing)
```

**Bağlam Modeli:**
- Her site, HTTP başlıklarından çözümlenir (`X-Site-Code` veya domain eşleme)
- Site kapsamlı içerik operasyonları, limitler ve bayraklar bağımsız yapılandırılır
- Kullanıcı izinleri site düzeyinde atanabilir

---

## 3. Çok Tenant SaaS

**Kullanım:** Birden fazla bağımsız müşteri aynı platform üzerinde barındırılır.

```
[Yük Dengeleyici]
      │
  ┌───┴───┐
  │  App  │  ← aystek-core + aystek-panel
  └───┬───┘
      ├── Tenant: "sirketA" → Site: "a.com"
      ├── Tenant: "sirketB" → Site: "b.com"
      └── Tenant: "sirketC" → Site: "c.com"

[Paylaşılan Altyapı]
  ├── MongoDB (tenant başına DB veya koleksiyon izolasyonu)
  ├── Redis   (tenant-prefixli key izolasyonu)
  └── RabbitMQ (tenant-prefixli kuyruk izolasyonu)
```

**Güvenlik:**
- Her istek tenant + site bağlamı doğrulandıktan sonra işlenir
- Çapraz tenant veri sızıntısı middleware katmanında engellenir
- RBAC izinleri tenant kapsamlıdır
- Onay politikaları tenant ve site bazında tutulur; kayıtlı politika yoksa refund onayı kapalı kalır, B2B satın alma gönderimi başarısız olur (bkz. Bölüm 8)

---

## 4. Güncelleme / Rollback Senaryoları

### 4.1 Başarılı Güncelleme (Normal Akış)

```
Tetikleyici: Panel / `aystek-core-installer`
      ↓
Ön kontroller (disk, bağlantı, guardrail doğrulamaları)
      ↓
Gerekli ise sürüm/artefact bilgisi çözülür
      ↓
Yedek al (/var/backups/aystek-core/pre-update-vX.Y.Z/)
      ↓
Bakım modu AÇ
      ↓
systemctl stop aystek-core aystek-panel
      ↓
İdempotent DB migration çalıştır ve tam `aystek-core` seed zincirini `SEED_FULL_RESET=false` ile çalıştır (silme yok; bazı kayıtlar kod tanımına göre yeniden yazılır, bkz. Bölüm 8)
      ↓
ln -sfn releases/v1.3.0 current   ← ATOMİK
      ↓
systemctl start aystek-core aystek-panel
      ↓
/health kontrolü (60sn timeout) → 200 OK?
      ↓
Bakım modu KAPAT
      ↓
Audit log + Bildirim
```

**Tahmini Kesinti Süresi:** PATCH: ~30sn | MINOR: ~60sn | MAJOR: 2-5dk

### 4.2 Başarısız Güncelleme → Otomatik Rollback

Sağlık kontrolü başarısız olursa operatör `aystek-core-installer rollback` veya System Ops yüzeyi ile rollback başlatır:

```
Sağlık kontrolü BAŞARISIZ (60sn timeout)
      ↓
Otomatik rollback tetiklenir
      ↓
systemctl stop aystek-core aystek-panel
      ↓
ln -sfn releases/v1.2.0 current   ← Eski sürüme dön
      ↓
systemctl start aystek-core aystek-panel
      ↓
/health kontrolü
      ↓
Bakım modu KAPAT
      ↓
Uyarı audit log + Bildirim
```

### 4.3 Manuel Rollback

Panel: **Sistem > Rollback** → Hedef sürüm seç → onay token akışı

> **Not:** Geri alma uygulama sürümünü geri döndürür. Güncellemenin tam seed zincirinde yeniden yazdığı kayıtlar (ör. `admin` rolünün izin seti) otomatik olarak önceki değerlerine dönmez; gerekirse güncelleme öncesi yedekten geri yükleyin (Bölüm 8). Gürsel GPS yapılandırmaları güncellemede değişmediği için bunlar için geri yükleme gerekmez.

CLI:
```bash
sudo aystek-core-installer rollback
```

---

## 5. Felaket Kurtarma (DR)

### Yedekleme Stratejisi

| Tür | Sıklık | Saklama | Konum |
|-----|--------|---------|-------|
| Güncelleme öncesi otomatik | Her güncelleme | 30 gün | `/var/backups/aystek-core/pre-update-*/` |
| Günlük zamanlanmış | 02:00 TSI | 30 gün | `/var/backups/aystek-core/scheduled/` |
| Manuel | İsteğe bağlı | — | `/var/backups/aystek-core/manual/` |

### Kurtarma Adımları

```bash
# 1. Mevcut release kopyala (güvenlik için)
cp -r /opt/aystek-core/releases/v1.2.0 /opt/aystek-core/releases/v1.2.0.bak

# 2. Yeni sunucuya yedekten geri yükle
rsync -avz /var/backups/aystek-core/db/ yeni-sunucu:/tmp/db-restore/

# 3. MongoDB geri yükleme
mongorestore --uri "mongodb://aystek:PAROLA@127.0.0.1:27017/aystek_core"   /tmp/db-restore/

# 4. Aystek Core'u kur (kurulum tam seed zincirini çalıştırır: admin rolünün izin seti kod tanımına döner; Gürsel GPS yapılandırmaları yedekteki haliyle kalır (yedekte olmayanlar kod değerleriyle eklenir); editor, reviewer ve fleet_viewer rolleri yedekteki haliyle kalır; bkz. Bölüm 8)
sudo aystek-core-installer install --config /etc/aystek-core/installer.json --non-interactive

# 5. Env dosyalarını kopyala ve servisleri başlat
```

### RTO / RPO Hedefleri

| Metrik | Hedef |
|--------|-------|
| RTO (Kurtarma Süresi) | < 2 saat |
| RPO (Veri Kaybı) | < 24 saat (zamanlanmış yedek) |
| Güncelleme rollback | < 5 dakika |

---

## 6. Yatay Ölçekleme (Gelecek Faz)

> Bu senaryo şu an için yol haritasında olup henüz üretim onaylı değildir.

```
[HAProxy / AWS ALB]
      ├── [App 1] aystek-core + aystek-panel
      ├── [App 2] aystek-core + aystek-panel
      └── [App N] ...

[Paylaşılan Altyapı Katmanı]
  ├── MongoDB Replica Set (3 node)
  ├── Redis Sentinel / Cluster
  └── RabbitMQ Cluster
```

**Gereksinimler:**
- Tüm app node'larında `/opt/aystek-core/` dizin yapısı aynı
- `uploads/` ve `shared/` dizinleri NFS veya S3 üzerinden paylaşılır
- Güncelleme orchestratörü tüm node'ları koordineli olarak günceller
- Oturum verisi Redis'te merkezi saklandığı için durumsuz (stateless) ölçekleme mümkündür

---

## 7. AYSTEK Haber Public-Site Senaryosu

AYSTEK Haber, yeni bir Core instance'ı değil; mevcut Core'a store ve tenant
header'larıyla bağlanan bağımsız Next.js container'ıdır. Önerilen topoloji:

```text
Internet → Nginx/TLS → 127.0.0.1:3001 aystekhaber-com
                              ↓ aystek-net
                         aystek-core:9090
                              ↓
                  site-scoped MongoDB içerikleri
```

Container yalnız loopback üzerinde yayınlanır. Nginx primary domain ve
opsiyonel `www` alias'ını karşılar. Public-site ortamı `API_BASE_URL`,
`AYSTEK_STORE_DOMAIN=aystekhaber.com` ve `AYSTEK_TENANT_ID=aystek` sözleşmesine
uymalıdır. Ayrıntılı kurulum ve rollback adımları için
`05-operations/SUNUCU_KURULUM_REHBERI.md` Bölüm 15'e bakın.

### Küresel Hat ve ortak haber kaynağı

`kureselhat.com`, aynı `aystek` tenant'ında farklı bir `site_id` kullanır.
Yerel geliştirmede iki domain tek `aystekhaber-com` source checkout'ını paylaşır;
sunucuda ise her domain aynı repository/sürümden gelen ayrı source dizini,
install/config/log dizini, runtime env, port, container ve Nginx/TLS vhost ile
çalışır. AYSTEK Haber `3001`, Küresel Hat `3002` loopback portunu kullanır.

Tek image iki deployment'ta kullanılabilir; site kimliği image'a gömülmez.
`SITE_URL`, `SITE_DOMAIN_ALIASES`, `AYSTEK_STORE_DOMAIN` ve `AYSTEK_TENANT_ID`
runtime sözleşmesidir. Host yalnız kabul kontrolüdür ve tenant/site seçemez.
`go run ./cmd/seed --news-site=kureselhat` yalnız Küresel Hat başlangıç kayıtlarını
ekler; normal seed zinciri de yeni kaydı içerir ve örnek haber/yazar/galeri
aktarılmaz. Kod tarafındaki marka/SEO/feed ve domain izolasyonu tamamlanmıştır;
gerçek içerik, hukuki sayfalar, logo, DNS ve TLS yayın onayı için hâlâ gereklidir.
Yayın kontrol listesi için sunucu kurulum rehberi Bölüm 15.1'e bakın.

---

## 8. 2026-10-28 Onay politikaları ve tohum (seed) değişiklikleri

Onay politikaları (iade onayı ve B2B satın alma zinciri) tenant ve site bazında tutulur. Politikası kayıtlı olmayan bir kapıda kod varsayılanı kullanılır. Seed yalnızca varsayılan tenant `aystek` ve onun yönetim paneli store'u için politika ekler. Dağıtım senaryosuna göre etkiler:

| Senaryo | Etki | Operatör işlemi |
|---|---|---|
| Tek tenant — tek site (Bölüm 1) | Refund onayı kapalı kalır. B2B zinciri yalnızca `aystek` tenant'ında varsayılan olarak etkindir; zincirin rolleri (`store_manager`, `credit_dept`, `executive`) oluşturulana kadar adımlar onaylanamaz. Kendi tenant'ınız için politika kaydedilmediyse B2B gönderimi `approvals.error.chain_not_configured` ile başarısız olur. | Rolleri oluşturun veya zinciri mevcut rollerle güncelleyin; politikayı kendi tenant'ınız için kaydedin. |
| Tek tenant — çok site (Bölüm 2) | Politika site bazındadır. B2B siparişi alınacak her site için ayrıca politika gerekir. | Her site için politikayı kaydedin. |
| Çok tenant SaaS (Bölüm 3) | Yeni tenant ve site için politika kaydedilmeden satın alma siparişi gönderimi başarısız olur. | Her tenant ve site için politikayı kaydedin. |
| Güncelleme (Bölüm 4.1) | Installer ile güncelleme tam seed zincirini çalıştırır (yeni izin modülleri, `admin` rolünün izin seti, varsayılan politikalar). `editor`, `reviewer` ve `fleet_viewer` rollerinin panel değişiklikleri korunur; bu rollere yeni varsayılan izin eklenmez. Installer dışı dağıtımda hedefli `--roles` gerekir. | Installer dışı dağıtımda `aystek-seed --roles` çalıştırın; sürüm bu rollere yeni varsayılan izin ekliyorsa installer ile güncellemeden sonra da çalıştırın; her durumda sonucu kontrol edin. |
| Felaket kurtarma (Bölüm 5) | Geri yüklemeden sonra installer install tam seed zincirini çalıştırır; `admin` rolünün izin seti kod tanımına döner. Gürsel GPS yapılandırmaları yedekteki haliyle kalır (yedekte olmayanlar kod değerleriyle eklenir). `editor`, `reviewer` ve `fleet_viewer` rolleri yedekteki haliyle kalır. | Geri yüklemeden sonra bu kayıtları kontrol edin. |

Ayrıntılı komutlar ve kontrol listesi: `SUNUCU_KURULUM_REHBERI.md` Bölüm 10.4 ve `OPERATIONS_GUIDE.md` Bölüm 3.15.
