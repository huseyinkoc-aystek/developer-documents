# AYSTEK Core — Uzak Sunucu Kurulum Rehberi

**Sürüm:** 2.0.0  
**Tarih:** 2026-04-03  
**Son Güncelleme:** 2026-08-04 — Docker + host Nginx üç-domain TLS yenileme notu  
**Hedef Ortam:** CentOS 9 (birincil), Ubuntu 22.04+ (ikincil)  
**Durum:** Operasyonel Olarak Geçerli — `aystek-core-installer` 2.x ile hizalı

---

## İçindekiler

1. [Gereksinimler](#1-gereksinimler)
2. [Ön Hazırlık — Sunucu Yöneticisi](#2-ön-hazırlık--sunucu-yöneticisi)
3. [Bağımlılık Kurulumu](#3-bağımlılık-kurulumu)
4. [AYSTEK Core Paketi İndirme ve Doğrulama](#4-aystek-core-paketi-indirme-ve-doğrulama)
5. [İlk Kurulum](#5-ilk-kurulum)
6. [Ortam Dosyalarını Yapılandırma ve Port Değiştirme](#6-ortam-dosyalarını-yapılandırma)
7. [systemd Servis Aktivasyonu](#7-systemd-servis-aktivasyonu)
8. [Nginx Reverse Proxy Yapılandırması](#8-nginx-reverse-proxy-yapılandırması)
9. [İlk Başlatma ve Lisans Aktivasyonu](#9-i̇lk-başlatma-ve-lisans-aktivasyonu)
10. [Güncelleme Prosedürü](#10-güncelleme-prosedürü)
11. [Rollback Prosedürü](#11-rollback-prosedürü)
12. [Sağlık Kontrolü ve Sorun Giderme](#12-sağlık-kontrolü-ve-sorun-giderme)
13. [Kaldırma (Uninstall)](#13-kaldırma-uninstall)
14. [Dizin Yapısı Referansı](#14-dizin-yapısı-referansı)
15. [AYSTEK Haber Kurulum Yol Haritası](#15-aystek-haber-kurulum-yol-haritası)

---

## 1. Gereksinimler

> **Mevcut Gerçeklik Notu (2026-04-05):** Bu rehberin kanonik kurulum yolu `aystek-core-installer` 2.x kullanımıdır. Manual `install.sh`, eski `aystek-updater` daemon anlatısı veya tarihsel `14 adımlı updater akışı` referansları artık kanonik işletim modeli değildir. Installer, native/docker/mongodb-only/tenant-site akışlarını yöneten ana yaşam döngüsü aracıdır. Installer davranışı değiştiğinde bu rehber aynı işte güncellenmelidir.
>
> **Artifact kaynağı notu (2026-04-24):** Installer yalnızca private GHCR image akışına bağlı değildir. Kurulumda artifact kaynağı olarak şu seçenekler desteklenmelidir:
> - `docker + registry images`
> - `docker + local source directories`
> - `native + local artifacts`
> - `native + local source directories`
> Sunucuda `aystek-core` / `aystek-core-next` repo dizinleri varsa installer bu yerel kaynakları da kurulum algoritmasında birinci sınıf seçenek olarak ele almalıdır.
> Docker kurulumunda operatör yanlışlıkla image tag alanına `/home/.../aystek-core` gibi yerel bir dizin yolu verirse installer bunu registry image olarak zorlamaz; local-source olarak yorumlar ve registry kimlik doğrulama adımını atlar.
>
> **Toolchain doğrulama notu:** `native + local-source` modunda hedef sunucuda `go` ve `npm`; `docker + local-source` modunda `docker` erişimi ve ilgili Dockerfile’lar bulunmalıdır. Installer bu prerequisite’leri deploy başlamadan önce doğrular.
>
> **Secret doğrulama notu:** MongoDB, Redis, RabbitMQ, SMTP ve registry secret alanları boş bırakılmıyorsa installer bunları gizli olarak iki kez ister. Böylece TTY üzerinde görünmeyen parola girişlerinde yanlış değerle devam etme riski düşer.
>
> **Apache / Nginx birlikte çalışma notu:** Aynı IP ve aynı `80/443` portlarında Apache/httpd ile Nginx’i birlikte bind etmek doğru değildir. Installer bu durumda artık özel karar ekranı açar: `Apache → Nginx migration`, `Apache'yi koru / Nginx'i kapat`, `kurulumu durdur`. Migration yalnızca desteklenen Apache virtual host, reverse proxy ve TLS desenleri için otomatik yapılır; yapılandırmalar önce yedeklenir, destek dışı kurallar raporlanır ve installer güvenli şekilde durur.
>
> **Aktif Apache site sözleşmesi:** Installer Apache tarafında aktifliği klasör adına göre varsaymaz. Ana config dosyasından (`/etc/apache2/apache2.conf` veya `/etc/httpd/conf/httpd.conf`) başlar, `ServerRoot` değerini çözer ve `Include` / `IncludeOptional` zincirini takip ederek Apache'nin gerçekten yüklediği `.conf` dosyalarını belirler. `sites-enabled`, `sites-available`, `conf.d` veya özel bir dizin ancak bu include graph içinde yer alıyorsa aktif sayılır. Ana config okunamazsa son çare olarak konvansiyonel dizinler taranır ve migration raporunda bu belirsizlik görünür tutulur. Symlink hedefleri tekilleştirilir; aynı domainin HTTP/HTTPS vhost'ları tek Nginx çıktısına merge edilir.
>
> **AI destek sınırı:** Gerekirse migration raporu bir yapay zeka modeliyle yorumlatılabilir; ancak otomatik kurulum güvenliği deterministic parser + `nginx -t` + operatör onayı üzerine kuruludur. LLM önerisi tek başına Apache'yi durdurma, config üretme veya sertifika yolu seçme kararı vermez.
>
> **Uninstall davranışı:** Installer uninstall sırasında Nginx paketini silmez; sadece AYSTEK tarafından yazılan Nginx konfigürasyonlarını kaldırır. Eğer aynı kurulum Apache'den Nginx'e migration yapmışsa ve Nginx önceden sistemde kurulu değilse, uninstall Apache servisini yeniden etkinleştirmeyi dener. Eğer Nginx zaten sistemde önceden kuruluysa Apache otomatik geri alınmaz. Bu geri dönüş davranışı uninstall onay ekranında açıkça gösterilir ve migration raporuna da yazılır.
>
> **Seed garantisi:** Installer fresh install ve update sırasında gerçek `aystek-core` seed zincirini de çalıştırır (`SEED_FULL_RESET=false`; veri silinmez). Bu nedenle `translation_seed`, `menu_seed` ve benzeri çekirdek seed içerikleri ortam tipinden bağımsız olarak veritabanına reconcile edilir. Zincir bazı mevcut kayıtları kod tanımına göre yeniden yazar (`admin` rolünün izin seti); Gürsel GPS sağlayıcı yapılandırmaları yalnızca eksikse eklenir, mevcut olanların kimlik bilgileri, temel URL, aktif bayrağı, görünen ad ve saat dilimi farkı değişmez; `editor`, `reviewer` ve `fleet_viewer` rollerinin panelden yapılan izin değişikliklerine dokunmaz; ayrıntı için Bölüm 10.4.
>
> **Docker port sözleşmesi:** Docker modunda operatörün seçtiği MongoDB/Redis/RabbitMQ portları host-side bind portlarıdır. Container içi servis keşfi sabit Compose DNS adlarıyla yapılır: `mongodb:27017`, `redis:6379`, `rabbitmq:5672`. Bu ayrım bozulursa `aystek-core` container'ı sağlıklı başlayamaz.
>
> **Docker healthcheck notu:** İlk açılışta seed/reconcile uzun sürebileceği için `aystek-core` healthcheck penceresi geniş tutulmalıdır. Güncel installer core için `start_period: 180s`, `retries: 20`; panel için `start_period: 90s`, `retries: 10` üretir.
>
> **Progress ve diagnostik notu:** Installer `install`, `update`, `repair` sırasında adım numarasıyla ilerleme gösterir. Docker stack başlatılırken `aystek-core` unhealthy kalırsa hata çıktısına otomatik olarak `docker compose ps`, `aystek-core` logları, backing service logları ve container health inspect özeti eklenmelidir.

### 1.1 Donanım Gereksinimleri

| Bileşen | Minimum | Önerilen (Üretim) |
|---------|---------|-------------------|
| CPU | 2 çekirdek | 4+ çekirdek |
| RAM | 4 GB | 8+ GB |
| Disk (/)  | 20 GB | 40+ GB |
| Disk (/var) | 20 GB | 50+ GB (log + backup) |
| Ağ | 100 Mbit | 1 Gbit |

### 1.2 Yazılım Gereksinimleri

| Bağımlılık | Sürüm | Zorunlu | Not |
|-----------|-------|---------|-----|
| CentOS Stream 9 / Ubuntu 22.04 | — | ✅ | Native veya Docker host olarak |
| MongoDB | 6.0+ | ✅ | Native kurulumda gerekli |
| Redis | 7.0+ | ✅ | Native kurulumda gerekli |
| RabbitMQ | 3.12+ | ✅ | Native kurulumda gerekli |
| Node.js | 20 LTS | ✅ | Native panel çalıştırma için |
| Nginx | 1.24+ | Önerilir | Reverse proxy ve TLS sonlandırma |
| Docker Engine | 24.x+ | Docker kurulumunda ✅ | `docker --version` ile doğrulayın |
| Docker Compose Plugin | v2.x | Docker kurulumunda ✅ | `docker compose version` ile doğrulayın |
| systemd | 250+ | ✅ | Native servis yönetimi için |
| curl, tar, sha256sum | — | ✅ | Paket doğrulama ve sağlık kontrolleri |
| **aystek-core-installer** | **v2.0.0+** | **⭐ Önerilir** | **Kurulum / update / rollback / backup / validate işlemlerini otomatik yönetir** |

### 1.3 Ağ Portları

| Varsayılan Port | Servis | Erişim | Değiştirilebilir mi? | Güvenlik Notu |
|-----------------|--------|--------|----------------------|---------------|
| 443 | HTTPS (Nginx) | Genel | ✅ | İsterseniz `8443` gibi özel port kullanabilirsiniz |
| 80 | HTTP → HTTPS yönlendirme | Genel | ✅ | İsterseniz `8080` kullanabilirsiniz |
| 9090 | AYSTEK Core API | Sadece dahili (Nginx proxy) | ✅ (`APP_PORT`) | Güvenlik için `19090` gibi yüksek port önerilir |
| 3000 | AYSTEK Panel (Next.js) | Sadece dahili (Nginx proxy) | ✅ (`PORT`) | Güvenlik için `13000` gibi yüksek port önerilir |
| 27017 | MongoDB | Sadece dahili | ✅ | Public bind yapmayın; `127.0.0.1` veya Docker internal network kullanın |
| 6379 | Redis | Sadece dahili | ✅ | `requirepass` ve `bind 127.0.0.1` önerilir |
| 5672 | RabbitMQ AMQP | Sadece dahili | ✅ | Sadece uygulama ağına açın |
| 15672 | RabbitMQ Yönetim | Sadece VPN/dahili | ✅ | Public erişime açmayın; gerekirse farklı port kullanın |

> ✅ **Önerilen güvenli yaklaşım:** dış dünyaya yalnızca `80/443` açık olsun. API, panel ve altyapı servisleri ya `127.0.0.1` üstünden ya da Docker private network içinde çalışsın.

---

## 2. Ön Hazırlık — Sunucu Yöneticisi

### 2.1 SELinux ve Güvenlik Duvarı

```bash
# CentOS 9 — SELinux enforcing modunda bırakın (üretim için)
getenforce         # Enforcing olmalı

# Nginx için gerekli SELinux izni
setsebool -P httpd_can_network_connect 1

# Güvenlik için örnek özel portlar
APP_PORT=19090
PANEL_PORT=13000

# Güvenlik duvarı kuralları
firewall-cmd --permanent --add-service=http
firewall-cmd --permanent --add-service=https

# Not: API ve panel portlarını mümkünse public açmayın.
# Gerçekten gerekiyorsa sadece dahili ağ/VPN kaynaklarına izin verin.
firewall-cmd --permanent --add-rich-rule="rule family='ipv4' source address='10.0.0.0/8' port protocol='tcp' port='${APP_PORT}' accept"
firewall-cmd --permanent --add-rich-rule="rule family='ipv4' source address='10.0.0.0/8' port protocol='tcp' port='${PANEL_PORT}' accept"
firewall-cmd --reload
```

### 2.2 Zaman Senkronizasyonu

```bash
# Sunucu saatinin doğru olması kritik (JWT, audit log)
timedatectl set-timezone Europe/Istanbul
systemctl enable --now chronyd
chronyc tracking
```

### 2.3 Sistem Güncellemeleri

```bash
# CentOS 9
dnf update -y

# Ubuntu 22.04
apt update && apt upgrade -y
```

---

## 3. Bağımlılık Kurulumu

> **Not:** Bu bölümdeki yerel bağımlılık kurulumları, yalnızca installer akışında ilgili servisler `local` seçildiğinde otomatik yapılır. `external`, `docker`, `mongodb-only` ve `tenant-site` senaryolarında seçime göre farklı adımlar uygulanır. Bu bölümü yalnızca bağımlılıkları ayrı ayrı yönetmeniz gerekiyorsa uygulayın.

### 3.1 MongoDB 6.0 (CentOS 9)

```bash
cat > /etc/yum.repos.d/mongodb-org-6.0.repo << 'EOF'
[mongodb-org-6.0]
name=MongoDB Repository
baseurl=https://repo.mongodb.org/yum/redhat/9/mongodb-org/6.0/x86_64/
gpgcheck=1
enabled=1
gpgkey=https://www.mongodb.org/static/pgp/server-6.0.asc
EOF

dnf install -y mongodb-org
systemctl enable --now mongod

# Üretim yapılandırması (yetkilendirme açın)
mongosh --eval "
  use admin
  db.createUser({
    user: 'aystek',
    pwd: 'GÜÇLÜ_PAROLA_BURAYA',
    roles: [{ role: 'dbOwner', db: 'aystek_core' }]
  })
"
```

MongoDB `/etc/mongod.conf` içinde güvenliği etkinleştirin:
```yaml
security:
  authorization: enabled
net:
  bindIp: 127.0.0.1
```

### 3.2 Redis 7 (CentOS 9)

```bash
dnf install -y redis
systemctl enable --now redis

# Üretim için redis.conf
# /etc/redis/redis.conf içinde:
#   bind 127.0.0.1
#   requirepass GÜÇLÜ_REDIS_PAROLA
#   maxmemory 512mb
#   maxmemory-policy allkeys-lru
```

### 3.3 RabbitMQ (CentOS 9)

```bash
# Erlang önce kurulur
dnf install -y erlang

# RabbitMQ repository
curl -s https://packagecloud.io/install/repositories/rabbitmq/rabbitmq-server/script.rpm.sh | bash
dnf install -y rabbitmq-server

systemctl enable --now rabbitmq-server

# Yönetim arayüzü (opsiyonel, VPN arkasında kullanın)
rabbitmq-plugins enable rabbitmq_management

# Kullanıcı oluştur
rabbitmqctl add_user aystek GÜÇLÜ_PAROLA
rabbitmqctl set_user_tags aystek administrator
rabbitmqctl set_permissions -p / aystek ".*" ".*" ".*"
```

### 3.4 Node.js 20 LTS

```bash
# CentOS 9
curl -fsSL https://rpm.nodesource.com/setup_20.x | bash -
dnf install -y nodejs

# Ubuntu 22.04
curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
apt install -y nodejs

node --version   # v20.x.x olmalı
```

### 3.5 Nginx

```bash
# CentOS 9
dnf install -y nginx
systemctl enable nginx

# Ubuntu 22.04
apt install -y nginx
```

### 3.6 Docker Engine + Docker Compose Plugin (Alternatif Kurulum Yolu)

```bash
# CentOS 9
dnf install -y dnf-plugins-core
# Docker resmi repository
# shellcheck disable=SC2086
dnf config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
dnf install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
systemctl enable --now docker

# Ubuntu 22.04
apt update
apt install -y ca-certificates curl gnupg
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
chmod a+r /etc/apt/keyrings/docker.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo \"$VERSION_CODENAME\") stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
apt update
apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
systemctl enable --now docker

# Doğrulama
docker --version
docker compose version
```

### 3.7 Docker İçin Güvenlik Varsayılanları

- Container portlarını mümkün olduğunca `127.0.0.1:<host_port>:<container_port>` biçiminde bind edin.
- MongoDB, Redis ve RabbitMQ portlarını public internete açmayın.
- Varsayılan portlar yerine örneğin `19090`, `13000`, `12701`, `16379`, `15673`, `25672` gibi yüksek portlar kullanın.
- Dış erişim için yalnızca Nginx üzerinden `80/443` yayınlayın.
- `docker compose` ile çalışan servislerde kalıcı veriler için named volume veya host volume kullanın.

---

## 4. AYSTEK Core Paketi İndirme ve Doğrulama

> **Not:** `aystek-core-installer` binary'sini kullanıyorsanız paket indirme işlemi de otomatik gerçekleştirilir. Bu bölüm yalnızca manuel kurulum senaryolarında veya özel indirme günlüğü gerektiğinde uygulanır.

### 4.1 Paketi İndirme

```bash
# AYSTEK güncelleme sunucusundan son sürümü indirin
AYSTEK_VERSION="1.2.0"
PACKAGE_URL="https://updates.aysteknolojileri.com/releases/aystek-core-v${AYSTEK_VERSION}-linux-amd64.tar.gz"

mkdir -p /tmp/aystek-install
cd /tmp/aystek-install

curl -fsSL "${PACKAGE_URL}" -o "aystek-core-v${AYSTEK_VERSION}.tar.gz"
curl -fsSL "${PACKAGE_URL}.sha256" -o "aystek-core-v${AYSTEK_VERSION}.tar.gz.sha256"
```

### 4.2 Bütünlük Doğrulama

```bash
# SHA256 kontrolü — KRİTİK, atlamamayın
sha256sum -c "aystek-core-v${AYSTEK_VERSION}.tar.gz.sha256"
# Çıktı: aystek-core-vX.Y.Z.tar.gz: OK  olmalı

# GPG imza doğrulama (isteğe bağlı ama önerilir)
curl -fsSL https://updates.aysteknolojileri.com/keys/aystek-release.gpg | gpg --import
gpg --verify "aystek-core-v${AYSTEK_VERSION}.tar.gz.sig" \
              "aystek-core-v${AYSTEK_VERSION}.tar.gz"
```

> ⚠️ **Uyarı:** SHA256 eşleşmezse kurulumu durdurun. Paketi yeniden indirin.

---

## 5. İlk Kurulum

### 5.0 aystek-core-installer CLI ile Kurulum (Önerilen Yol)

`aystek-core-installer`, seçilen moda göre bağımlılık kurulumu, yapılandırma, systemd veya docker entegrasyonu, backup servisi ve doğrulamayı tek akışta yönetir.

```bash
# Binary'yi sunucuya kopyalayın veya AYSTEK dağıtım paketinden yükleyin
install -m 755 aystek-core-installer /usr/local/bin/aystek-core-installer

# Etkileşimli kurulum sihirbazını başlatın
sudo aystek-core-installer install

# Yapılandırma dosyasıyla etkileşimsiz kurulum (CI/CD için)
sudo aystek-core-installer install --config /etc/aystek-core/installer.json --non-interactive
```

Sihirbaz sizden alacağı yanıtlara göre tüm adımları otomatik yürütür:

| Adım | Otomatik Yapılan |
|------|-----------------|
| Bağımlılık kontrolü | Seçilen moda göre gerekli yerel bağımlılıklar |
| Ortam dosyaları | `/etc/aystek-core/core.env`, `panel.env` oluşturulur |
| systemd servisleri | Native modda unit dosyaları yazılır, `enable` edilir |
| Güvenlik duvarı | Private bind ve seçilen kurallara göre yapılandırılır |
| Yedekleme timer'ı | Backup etkinse systemd timer kaydedilir |
| Sağlık kontrolü | Kurulum sonrası otomatik doğrulama |

> **İlk admin bootstrap kuralı:** `admin` ve `tenant-site` scope kurulumlarında, veritabanı sıfırdan kuruluyorsa ve dump geri yüklenmiyorsa installer ilk admin kullanıcı bilgilerini ister. Bu bilgi AYSTEK Core `users` ve `roles` koleksiyonlarına seed edilerek ilk giriş için hazır admin hesabı oluşturulur. Dump ile geri yükleme yapılan kurulumlarda bu prompt beklenmez.

> `mongodb-only` modunu seçerseniz yalnızca MongoDB kurulur; uygulama sunucusu bu DB'ye dışarıdan bağlanacak şekilde yapılandırılır.

Kurulum tamamlandıktan sonra doğrulama:

```bash
sudo aystek-core-installer status
curl -s http://127.0.0.1:9090/health | python3 -m json.tool
```

---

### 5.1 Manuel Kurulum — Tarihsel / Özel Paketleme Senaryosu

```bash
# Paketi geçici dizine aç
tar -xzf "aystek-core-v${AYSTEK_VERSION}.tar.gz" -C /tmp/aystek-install/
cd "/tmp/aystek-install/aystek-core-v${AYSTEK_VERSION}"

# Bu yol yalnızca release engineering tarafından hazırlanmış özel paketlerde geçerlidir.
# Güncel kanonik yol: sudo aystek-core-installer install
```

> ⭐ **Güncel önerilen yol `sudo aystek-core-installer install` komutudur.** Manual script tabanlı yol, yalnızca özel release paketlerinde veya kontrollü iç operasyon senaryolarında kullanılmalıdır.

Betik aşağıdakileri otomatik olarak gerçekleştirir:

| Adım | İşlem |
|------|-------|
| 1 | OS uyumluluk ve disk alanı kontrolü |
| 2 | `aystek` sistem kullanıcısı oluşturma |
| 3 | Dizin yapısını oluşturma (`/opt/aystek-core`, `/etc/aystek-core`, vs.) |
| 4 | Binary dosyaları kopyalama |
| 5 | Symlink oluşturma: `current → releases/v2.0.0` |
| 6 | systemd unit dosyaları kopyalama |
| 7 | Örnek ortam dosyaları oluşturma |
| 8 | Dosya izinleri ayarlama |

### 5.2 Kurulum Sonrası Dizin Yapısı (Native Referans)

```
/opt/aystek-core/
├── releases/
│   └── v2.0.0/                ← aktif release
│       ├── bin/
│       │   ├── aystek-core    ← Go API binary
│       │   └── aystek-cli     ← CLI aracı
│       ├── panel/             ← Next.js build
│       ├── scripts/           ← Kurulum scriptleri
│       └── config/            ← Yapılandırma şablonları
├── current -> releases/v2.0.0 ← aktif symlink
└── shared/uploads/            ← tenant yüklemeleri (sürümler arası paylaşılır)
```

### 5.3 Docker Compose ile Kurulum (Alternatif)

> AYSTEK dağıtım paketiyle gelen `docker-compose.aystek-core.example.yml` ve `docker.env.example` dosyalarını başlangıç şablonu olarak kullanın.

#### 5.3.1 Çalışma Dizini ve `.env` Dosyası

```bash
mkdir -p /opt/aystek-docker/{data/mongodb,data/redis,data/rabbitmq,logs}
cd /opt/aystek-docker

# Paketle gelen örnek dosyaları kopyalayın
cp docker-compose.aystek-core.example.yml ./docker-compose.yml
cp docker.env.example ./.env

# Gizli bilgileri düzenleyin
nano .env
```

#### 5.3.2 Servisleri Başlatma

```bash
docker compose pull
docker compose up -d
docker compose ps
```

#### 5.3.3 Sağlık Kontrolü

```bash
curl -s http://127.0.0.1:${APP_PORT:-19090}/health | python3 -m json.tool
```

> Güvenlik için API/panel portlarını doğrudan internete açmayın; Nginx reverse proxy üzerinden yayınlayın.

---

## 6. Ortam Dosyalarını Yapılandırma

> 🔒 **Güvenlik:** `.env` dosyaları `root:aystek 640` izniyle korunur. Asla sürüm kontrolüne eklemeyin.

### 6.1 Ana API Ortam Dosyası

```bash
cp /opt/aystek-core/current/config/core.env.example /etc/aystek-core/core.env
chmod 640 /etc/aystek-core/core.env
chown root:aystek /etc/aystek-core/core.env
nano /etc/aystek-core/core.env
```

Düzenlenmesi gereken kritik değerler:

```bash
# /etc/aystek-core/core.env

APP_ENV=production
APP_PORT=9090   # güvenlik için örn. 19090 kullanılabilir

# MongoDB — parola olarak yeni oluşturduğunuzu kullanın
MONGODB_URI=mongodb://aystek:GÜÇLÜ_PAROLA@127.0.0.1:27017/aystek_core?authSource=admin

# Redis
REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=GÜÇLÜ_REDIS_PAROLA

# RabbitMQ
RABBITMQ_URL=amqp://aystek:GÜÇLÜ_PAROLA@127.0.0.1:5672/

# JWT — 64+ karakterlik rastgele bir değer üretin
# openssl rand -hex 32 ile üretin
JWT_SECRET=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

# Üçüncü taraf sırlarının (ödeme ağ geçidi anahtarları, SMTP hesap şifreleri) veritabanında
# şifreli saklanması için. Installer ENCRYPTION_KEY ile aynı değeri yazar; sonradan değiştirmeyin
# (kayıtlı sırlar çözülemez olur ve panelden yeniden girilmelidir).
PAYMENT_GATEWAY_ENCRYPTION_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

# CORS — örnek production domainleri
ALLOWED_ORIGINS=https://core.aysteknolojileri.com,https://panel.core.aysteknolojileri.com

# Log
LOG_LEVEL=info
LOG_FORMAT=json
LOG_DIR=/var/log/aystek-core/core
```

### 6.2 Panel Ortam Dosyası

```bash
cp /opt/aystek-core/current/config/panel.env.example /etc/aystek-core/panel.env
chmod 640 /etc/aystek-core/panel.env
chown root:aystek /etc/aystek-core/panel.env
nano /etc/aystek-core/panel.env
```

```bash
# /etc/aystek-core/panel.env
PORT=3000   # güvenlik için örn. 13000 kullanılabilir
NEXT_PUBLIC_API_BASE_URL=https://core.aysteknolojileri.com
NEXT_PUBLIC_API_URL=https://core.aysteknolojileri.com
NEXT_PUBLIC_SSE_BASE_URL=https://core.aysteknolojileri.com
NEXT_PUBLIC_APP_ENV=production
NEXT_PUBLIC_SITE_URL=https://panel.core.aysteknolojileri.com
NEXT_API_URL=http://app:9090
API_URL=http://app:9090
INTERNAL_API_URL=http://app:9090
NODE_ENV=production
```

Panel tarafında canonical env seti birlikte tutulmalıdır:
- `NEXT_PUBLIC_API_BASE_URL`
- `NEXT_PUBLIC_API_URL` (compatibility)
- `NEXT_PUBLIC_SSE_BASE_URL`
- `NEXT_PUBLIC_APP_ENV`
- `NEXT_PUBLIC_SITE_URL`
- `NEXT_API_URL`
- `API_URL`
- `INTERNAL_API_URL`

### 6.3 Installer ve Runtime Yapılandırması

```bash
cp /etc/aystek-core/installer.json /etc/aystek-core/installer.json.bak
nano /etc/aystek-core/installer.json
```

```bash
# installer.json içinde örnek kritik alanlar
"release_version": "2.0.0",
"api_port": 9090,
"panel_port": 3000,
"mongo_mode": "local",
"backup_enabled": true,
"backup_schedule": "*-*-* 02:00:00"
```

> Runtime environment için `core.env` ve `panel.env`; installer davranışı için `installer.json` kanonik dosyalardır. Tarihsel `update.conf` anlatısı artık ana referans değildir.

### 6.4 Port Değiştirme ve Güvenlik Sertleştirme

Aşağıdaki örnek değerler doğrudan kullanılabilir:

| Bileşen | Varsayılan | Önerilen Güvenli Örnek | Nereden Değişir? |
|--------|------------|------------------------|------------------|
| API | `9090` | `19090` | `core.env` → `APP_PORT` |
| Panel | `3000` | `13000` | `panel.env` → `PORT` |
| MongoDB | `27017` | `12701` | `mongod.conf` veya Docker port mapping |
| Redis | `6379` | `16379` | `redis.conf` veya Docker port mapping |
| RabbitMQ AMQP | `5672` | `15673` | `rabbitmq.conf` veya Docker port mapping |
| RabbitMQ UI | `15672` | `25672` | `rabbitmq.conf` veya Docker port mapping |

Port değiştirirken şu 5 yeri birlikte güncelleyin:

1. `core.env` içindeki `APP_PORT`
2. `panel.env` içindeki `PORT`
3. `installer.json` ve gerekirse `core.env` içindeki ilgili port/URL alanları
4. Nginx upstream tanımları (`127.0.0.1:9090` / `127.0.0.1:3000`)
5. Sağlık kontrolü ve firewall komutları

---

## 7. systemd Servis Aktivasyonu

> **Not:** `aystek-core-installer install` (`native` modda) systemd unit dosyalarını otomatik kopyalar, `daemon-reload` çalıştırır ve servisleri `enable` eder. Bu bölüm; manuel kurulum tamamlandıktan sonra veya servis yapılandırmasını doğrudan yönetmek istediğinizde referans alınabilir.

### 7.1 Native / systemd Servis Yönetimi

```bash
# Unit dosyaları aystek-core-installer veya kontrollü release süreci tarafından /etc/systemd/system/ içine kopyalandı
systemctl daemon-reload

# Servisleri etkinleştir (otomatik başlatma)
systemctl enable aystek-core aystek-panel

# Sırayla başlat
systemctl start aystek-core
sleep 5
systemctl start aystek-panel

# Durumu kontrol et
systemctl status aystek-core aystek-panel
```

```bash
# Logları izlemek için
journalctl -u aystek-core -f
journalctl -u aystek-panel -f
tail -f /var/log/aystek-core/audit.log
```

### 7.2 Docker Compose ile Servis Yönetimi

```bash
cd /opt/aystek-docker

# Servisleri başlat
docker compose up -d

# Durumu kontrol et
docker compose ps

# Belirli servisi yeniden başlat
docker compose restart aystek-core
docker compose restart aystek-panel

# Servis durdur / tekrar başlat
docker compose stop
docker compose start

# Tam kapat (veri kalır)
docker compose down

# Logları canlı izle
docker compose logs -f aystek-core
docker compose logs -f aystek-panel
docker compose logs -f mongodb redis rabbitmq
```

#### 7.2.1 Sunucu açılışında otomatik başlatma

```bash
# Docker servisi açılışta aktif olsun
systemctl enable --now docker

# Compose servisleri restart policy sayesinde otomatik kalkar.
# Örnek compose içinde: restart: unless-stopped
```

#### 7.2.2 Servis içine girme / bakım komutları

```bash
# API container shell
docker compose exec aystek-core sh

# Panel container shell
docker compose exec aystek-panel sh

# Mongo shell
docker compose exec mongodb mongosh -u "$MONGO_ROOT_USERNAME" -p "$MONGO_ROOT_PASSWORD"
```

---

## 8. Nginx Reverse Proxy Yapılandırması

```bash
cat > /etc/nginx/conf.d/aystek-core.conf << 'NGINX'
# AYSTEK Core — Nginx Yapılandırması
upstream aystek_api {
    server 127.0.0.1:9090;   # özel port kullanıyorsanız örn. 19090
    keepalive 32;
}

upstream aystek_panel {
    server 127.0.0.1:3000;   # özel port kullanıyorsanız örn. 13000
    keepalive 32;
}

# HTTP → HTTPS yönlendirme
server {
    listen 80;
    server_name core.aysteknolojileri.com panel.core.aysteknolojileri.com;
    return 301 https://$host$request_uri;
}

# API
server {
    listen 443 ssl http2;
    server_name core.aysteknolojileri.com;

    ssl_certificate     /etc/letsencrypt/live/core.aysteknolojileri.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/core.aysteknolojileri.com/privkey.pem;
    ssl_protocols       TLSv1.2 TLSv1.3;
    ssl_ciphers         HIGH:!aNULL:!MD5;

    # SSE için (güncelleme log stream)
    proxy_buffering off;
    proxy_cache off;

    location / {
        proxy_pass         http://aystek_api;
        proxy_set_header   Host $host;
        proxy_set_header   X-Real-IP $remote_addr;
        proxy_set_header   X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto $scheme;
        proxy_read_timeout 300s;
    }
}

# Yönetim Paneli
server {
    listen 443 ssl http2;
    server_name panel.core.aysteknolojileri.com;

    ssl_certificate     /etc/letsencrypt/live/core.aysteknolojileri.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/core.aysteknolojileri.com/privkey.pem;
    ssl_protocols       TLSv1.2 TLSv1.3;

    location / {
        proxy_pass       http://aystek_panel;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
NGINX

nginx -t && systemctl reload nginx
```

### 8.1 SSL Sertifikası (Let's Encrypt)

```bash
# CentOS 9
dnf install -y epel-release certbot certbot-nginx
certbot --nginx -d core.aysteknolojileri.com -d panel.core.aysteknolojileri.com

# Otomatik yenileme (RHEL ailesi)
systemctl enable --now certbot-renew.timer

# Debian/Ubuntu paketlerinde timer adı genellikle certbot.timer'dır
systemctl enable --now certbot.timer
```

Installer iki alan adı için tek bir SAN sertifikası üretir. Canonical lineage
`/etc/letsencrypt/live/core.aysteknolojileri.com/` altındadır ve sertifika hem
`core.aysteknolojileri.com` hem `panel.core.aysteknolojileri.com` adını kapsar.
Nginx'teki iki HTTPS vhost aynı `fullchain.pem` ve `privkey.pem` dosyalarını
kullanabilir.

#### 8.1.1 Süresi Dolmuş Sertifikayı Acil Yenileme

Docker kurulumunda API ve panel container'ları özel host portlarına publish
edilir; 80/443 portları host Nginx'e aittir. Certbot host üzerinde çalıştırılır
ve bu işlem için uygulama container'ları durdurulmaz. Renewal yapılandırmasında
`authenticator = standalone` varsa host Nginx port 80'i kullandığı için yenileme
başarısız olur. Nginx plugin'i kurularak lineage `nginx` authenticator'a
geçirilmelidir.

```bash
# Ön kontrol
dig +short core.aysteknolojileri.com
dig +short panel.core.aysteknolojileri.com
nginx -t
certbot certificates

# RHEL/CentOS/Rocky/AlmaLinux: Nginx plugin
dnf install -y epel-release
dnf install -y certbot-nginx || dnf install -y python3-certbot-nginx
certbot plugins

# Mevcut SAN lineage'ını nginx authenticator'a geçir, yenile ve Nginx'e yükle
certbot certonly --nginx \
  --cert-name core.aysteknolojileri.com \
  -d core.aysteknolojileri.com \
  -d panel.core.aysteknolojileri.com \
  --force-renewal
nginx -t && systemctl reload nginx

# Sunulan sertifikayı iki SNI adıyla doğrula
for domain in core.aysteknolojileri.com panel.core.aysteknolojileri.com; do
  echo | openssl s_client -connect "${domain}:443" -servername "$domain" 2>/dev/null \
    | openssl x509 -noout -subject -issuer -dates -ext subjectAltName
done
```

`renew` başarısız olur veya lineage iki alan adını kapsamıyorsa aynı lineage'ı
iki SAN ile yeniden üretin:

```bash
certbot certonly --nginx \
  --cert-name core.aysteknolojileri.com \
  -d core.aysteknolojileri.com \
  -d panel.core.aysteknolojileri.com \
  --force-renewal
nginx -t && systemctl reload nginx
```

Kalıcı doğrulama:

```bash
systemctl list-unit-files | grep -E '^certbot.*timer'
systemctl status certbot-renew.timer --no-pager || systemctl status certbot.timer --no-pager
certbot renew --cert-name core.aysteknolojileri.com --dry-run
```

Installer başarılı yenilemeden sonra Nginx'i güvenli biçimde yeniden yükleyen
`/etc/letsencrypt/renewal-hooks/deploy/reload-nginx.sh` hook'unu kurar. Uygun
systemd timer yoksa `/etc/cron.d/certbot-renew` fallback'i oluşturulur. Tamamlanmış
bir kurulumda `aystek-core-installer repair --verbose` TLS bakım adımını da yeniden
çalıştırır; sertifika dosyasının yalnızca varlığını değil tarih ve SAN kapsamını
kontrol eder.

#### 8.1.2 Operasyon Kaydı — Üç Domain Otomatik Yenileme (2026-08-04)

Üretimde Core API ve panel Docker container'larında, 80/443 trafiği host Nginx
üzerinde çalışır. Certbot host üzerinde çalıştırılır; TLS yenilemesi için Docker
container'ları durdurulmaz.

| Lineage | Kapsadığı domainler | Yönetim modeli |
|---|---|---|
| `core.aysteknolojileri.com` | `core.aysteknolojileri.com`, `panel.core.aysteknolojileri.com` | Installer tarafından yönetilen ortak SAN sertifikası |
| `mobinav.aysteknolojileri.com` | `mobinav.aysteknolojileri.com` | Apache migration'dan kalan legacy lineage; bir kez Nginx authenticator'a geçirilir |

Geçmişte Core lineage'ının `standalone`, Mobinav lineage'ının `apache`
authenticator kullanması nedeniyle yenilemeler sırasıyla port 80 çakışması ve
başarısız `httpd` restart hatası vermiştir. Kalıcı çözüm iki lineage'ı da
`nginx` authenticator'a geçirmek, Certbot timer'ı etkinleştirmek ve yenileme
sonrasında Nginx reload hook'unu çalıştırmaktır.

```bash
# Authenticator durumları
grep '^authenticator' \
  /etc/letsencrypt/renewal/core.aysteknolojileri.com.conf \
  /etc/letsencrypt/renewal/mobinav.aysteknolojileri.com.conf

# Gerekirse Core + Panel SAN lineage geçişi
certbot certonly --nginx \
  --cert-name core.aysteknolojileri.com \
  -d core.aysteknolojileri.com \
  -d panel.core.aysteknolojileri.com \
  --force-renewal

# Gerekirse Mobinav legacy lineage geçişi
certbot certonly --nginx \
  --cert-name mobinav.aysteknolojileri.com \
  -d mobinav.aysteknolojileri.com \
  --force-renewal

# Timer ve doğrulama
systemctl enable --now certbot-renew.timer || systemctl enable --now certbot.timer
certbot renew --cert-name core.aysteknolojileri.com --dry-run
certbot renew --cert-name mobinav.aysteknolojileri.com --dry-run
certbot renew --dry-run
```

Beklenen son durumda iki renewal dosyasında da `authenticator = nginx` bulunur
ve global dry-run `all simulated renewals succeeded` sonucu verir. Mobinav,
installer'ın `api_domain` / `panel_domain` sertifikasının parçası değildir;
lineage geçişi tek seferlik operatör adımıdır, sonraki otomatik yenilemeler ortak
Certbot timer/cron tarafından yapılır.

---

## 9. İlk Başlatma ve Lisans Aktivasyonu

### 9.1 Sağlık Kontrolü

```bash
# API yanıt veriyor mu?
curl -s http://127.0.0.1:9090/health | python3 -m json.tool
```

Beklenen yanıt:
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "version": { "backend": "1.2.0", ... },
    "services": {
      "mongodb":  { "status": "ok", "latency_ms": 2 },
      "redis":    { "status": "ok", "latency_ms": 1 },
      "rabbitmq": { "status": "ok", "latency_ms": 3 }
    }
  }
}
```

### 9.2 Lisans Aktivasyonu

```bash
# Lisans aktivasyonunu System Ops API veya yönetim paneli üzerinden yapın
curl -X POST http://127.0.0.1:9090/admin/system/license/activate \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"license_key": "LICENSE-XXXXX-XXXXX-XXXXX"}'
```

#### Docker için lisans aktivasyonu

```bash
cd /opt/aystek-docker

# Container içinden API erişimi ya da yönetim paneliyle aynı activation akışı kullanın.
```

### 9.3 İlk Yönetici Oluşturma

```bash
# İlk yönetici kullanıcı oluşturma, seed veya yönetim API akışıyla yapılmalıdır.
# Projenin kullanıcı/seed standartlarını izleyin.
```

#### Docker için ilk yönetici oluşturma

```bash
cd /opt/aystek-docker
# Docker dağıtımında da aynı seed veya yönetim API akışını kullanın.
```

---

## 10. Güncelleme Prosedürü

### 10.0 aystek-core-installer ile Güncelleme (Önerilen CLI Yolu)

```bash
# Güncellemeyi başlat — tüm adımlar otomatik yürütülür
sudo aystek-core-installer update

# Etkileşimsiz (CI/CD veya script)
sudo aystek-core-installer update --non-interactive
```

**`aystek-core-installer update` şunları otomatik gerçekleştirir:**

1. **Otomatik yedek** — güncelleme başlamadan önce `pre-update` etiketiyle mongodump al
2. Installer'a ait idempotent DB migration'ı ve tam `aystek-core` seed zincirini (`SEED_FULL_RESET=false`) çalıştırır; etkisi için Bölüm 10.4'e bakın
3. Config ve release klasörlerini hazırlar
4. Symlink'i yeni sürüme döndürür
5. Servisleri yeniden başlatır
6. Sağlık kontrolü yapar ve audit log yazar

> Yedek başarısız olursa güncelleme yine de devam eder (non-fatal); hata audit log'a yazılır.

---

### 10.1 Panel Üzerinden Güncelleme (Önerilir)

1. Yönetim paneline giriş yapın
2. **Sistem > Güncellemeler** menüsüne gidin
3. **Güncelleme Kontrolü** tuşuna tıklayın
4. Mevcut güncelleme varsa **Güncellemeyi Başlat** seçin
5. Güncelleme adımlarını canlı log akışından takip edin
6. Güncelleme tamamlandığında sağlık kontrolü otomatik yapılır

### 10.2 CLI ile Güncelleme

```bash
# Güncellemeyi başlat (aystek-core-installer — önerilen)
sudo aystek-core-installer update

# İlerlemeyi izle
tail -f /var/log/aystek-core/audit.log
```

### 10.2.1 Docker Compose ile Güncelleme

```bash
cd /opt/aystek-docker

# 1) Önce mevcut volume ve env dosyalarını yedekleyin
tar -czf /var/backups/aystek-docker-$(date +%F-%H%M).tar.gz .env docker-compose.yml

# 2) Yeni image tag'i .env içinde güncelleyin
sed -i 's|AYSTEK_CORE_IMAGE=.*|AYSTEK_CORE_IMAGE=ghcr.io/aystek/aystek-core:1.3.0|' .env
sed -i 's|AYSTEK_PANEL_IMAGE=.*|AYSTEK_PANEL_IMAGE=ghcr.io/aystek/aystek-panel:1.3.0|' .env

# 3) Yeni image'ları çekin ve servisleri yeniden oluşturun
docker compose pull
docker compose up -d --force-recreate

# 4) Sağlık ve log kontrolü
docker compose ps
docker compose logs --tail=200 aystek-core
curl -s http://127.0.0.1:${APP_PORT:-19090}/health | python3 -m json.tool
```

> **Not:** Bu manuel akış installer'ı atlar. Seed zinciri yalnızca konteyner başlatılırken `SEED_ON_START=true` ise çalışır; üretim compose dosyasında (`deployments/production/docker-compose.yml`) bu değer `false`'tur, installer kendi seed konteynerinde bayrağı açıkça `true` yapar; örnek compose dosyasında ise `true`'dur. Onay politikası ve yeni izin modülleri için Bölüm 10.4'teki hedefli `--roles` adımını uygulayın.

### 10.3 Güncelleme Çerçevesi (Güncel)

```
1. ÖN KONTROLLER   → State ve guardrail kontrolleri
2. YEDEK ALMA      → pre-update mongodump
3. MİGRASYON       → idempotent schema/seed akışı
4. CONFIG WRITE    → env ve runtime config güncellemesi
5. DEPLOY          → release dizini hazırlanır
6. SYMLINK         → current → yeni sürüm (atomik)
7. SERVİS BAŞLAT   → aystek-core + aystek-panel
8. SAĞLIK KONTROLÜ → validate / health
9. AUDIT LOG       → installer audit kaydı
```

---

### 10.4 2026-10-28 Onay politikaları ve tohum (seed) değişiklikleri

Bu sürüm onay (approval) modülünü, buna bağlı izin, rol ve seed değişikliklerini getirir. Mevcut bir veritabanına dağıtımdan sonra aşağıdaki adımlar uygulanmalıdır.

**Ne değişti**

- **Onay kapıları (gate):** Politikalar tenant ve site bazında tutulur. Politikası kayıtlı olmayan bir kapı kod varsayılanıyla çalışır.
  - `payments.refund`: tek eşikli onay. Varsayılan **kapalıdır**; onaylayıcının `payments:approve` izni olmalıdır.
  - `b2b.purchase_order`: onay adımlarının listesi. Varsayılan adımlar: `store_manager` (tutar > 50.000), `credit_dept` (tutar bayinin kredi limitini aşarsa; yalnızca limit > 0 ise) ve `executive` (tutar > 100.000). Adımlar sırayla zorunlu değildir (bkz. *Onay akışı*).
- **Yeni izin modülleri:** `approval_policies`, `approval_requests`.
- **Yeni admin uçları** (`/admin` altında): `GET /approval-gates`, `GET` ve `PUT /approval-policies/:gate`, `GET /approval-requests`, `POST /approval-requests/:id/approve`, `POST /approval-requests/:id/reject`. Politika okuma `approval_policies:read`, yazma `approval_policies:update` izni ister. Talep listeleme ve karar uçları `approval_requests:read` ister; karar ayrıca talebin kaydedilmiş iznini (ör. `payments:approve`) gerektirir.
- **Rol anahtarları:** Rolün `name` alanı küçük İngilizce harf, rakam ve `_` içerir, harfle başlar ve oluşturulduktan sonra değiştirilemez. Görünen adlar `labels` alanında dile göre tutulur.
- **Migrasyon aracı** (`cmd/migrate_gursel`): Kaynak bağlantı dizesi kodda yoktur. `MIGRATE_SOURCE_URI` ortam değişkeni veya `-src-uri` bayrağı zorunludur. `SOURCE_MONGO_URI` tanımlıysa değerini geçersiz kılar; ancak `MIGRATE_SOURCE_URI` veya `-src-uri` yine de boş olmamalıdır.

**Seed davranışı**

- `RunApprovalsSeed` (`internal/seeds/approvals_seed.go`) varsayılan politikaları yalnızca **eksikse** ekler; mevcut politikanın üzerine yazmaz. Hedef, varsayılan tenant `aystek` ve yönetim paneli store'udur (kod `admin-panel`).
- `go run ./cmd/seed --roles` (container içinde `aystek-seed --roles`) veri silmeden şunları çalıştırır: eksik izinleri ekler (`SeedPermissions`); `editor`, `reviewer` ve `fleet_viewer` rollerine eksik varsayılan izinleri ekler ve `admin` rolünü kod tanımıyla eşitler (`GrantRoleDefaults`); varsayılan onay politikalarını ekler (`RunApprovalsSeed`).
  - **Dikkat:** `editor`, `reviewer` ve `fleet_viewer` rollerinden hiçbir izin kaldırılmaz. Ancak panelden kaldırılmış bir varsayılan izin bu komutla yeniden eklenir; komutu bunu istediğinizde çalıştırın. `admin` rolü, `system` ve `system_ops` modülleri dışındaki tüm izinleri alır ve kod tanımına göre eşitlenir. Public-site yönetici rollerinin (ör. `balci_manager`) izin setleri ise `SeedRoles` çalıştığında ve `--public-site ... --create-manager` ile kod tanımına göre yeniden yazılır.
- Tam seed zinciri (`seeds.Run`, `SEED_ON_START=true`) `SeedPermissions`, `SeedRoles` ve `RunApprovalsSeed` dahil diğer tüm seed'leri içerir. `SeedRoles` mevcut `editor`, `reviewer` ve `fleet_viewer` rollerinin izinlerine dokunmaz, yalnızca eksik rolleri oluşturur; bu roller için yeni varsayılan izin eklemek `--roles` gerektirir. Installer kurulum ve update sırasında bunu `SEED_FULL_RESET=false` ile çalıştırır.
- **Gürsel GPS sağlayıcıları:** `seedGurselTurizimGpsProviders` (`internal/seeds/gursel_turizm_seed.go`), `fleet_tally_gps_provider_configs` koleksiyonunda trio, n2, mobiliz ve arvento yapılandırmalarını yalnızca eksikse ekler; ekleme sırasında kimlik bilgileri, temel URL, aktif bayrağı, görünen ad ve saat dilimi farkını kod değerleriyle yazar. Site referansı onarımı hariç, mevcut kayıtların bu alanlarına hiçbir tam zincir, güncelleme veya konteyner başlatması dokunmaz; panelden yapılan değişiklikler korunur. Yapılandırmayı değiştirmek için admin panelini kullanın. Kimlik bilgileri bu dokümanda yer almaz.

**Canlı ortam kuralları**

- Kurulum tamamlandıktan sonra canlı veritabanında tam seed zinciri çalıştırmayın; hedefli komutları (`--roles`, `--news-site`, `--public-site`) kullanın.
- `SEED_ON_START` değeri `true` olan bir konteyner her başlatmada tam zinciri çalıştırır. `deployments/production/docker-compose.yml` bu değeri `false` ile gelir; installer kendi tek seferlik seed konteynerinde bayrağı açıkça `true` yapar. Örnek compose dosyasında (`docker-compose.aystek-core.example.yml`) bu değer `true`'dur; canlı ortamda kurulumdan sonra `false` yapın.
- Tek tek seed hataları API başlangıcını durdurmaz ve `seed warning (non-fatal)` olarak loglanır; installer bunları başarısızlık saymaz, log'u kontrol edin. Zincirde bir seed hata verirse o adımdan sonraki seed'ler çalıştırılmaz (ör. `RunApprovalsSeed` hata verirse çeviriler ve menüler yazılmaz); roller bu adımdan önce yazıldığı için etkilenmez.
- Komutu kaynak dizininden çalıştırıyorsanız `.env` dosyası yalnızca ortamda tanımlı olmayan değişkenleri doldurur. Üretim `MONGODB_URI` ve `MONGO_DB` değerlerini komuttan önce dışa aktarın.

**Dağıtım sonrası adımlar (mevcut veritabanı)**

Installer ile `update` yaptıysanız tam zincir `editor`, `reviewer` ve `fleet_viewer` rollerine yeni varsayılan izin eklemez; sürüm bu rollere yeni varsayılan izin ekliyorsa 1. adımı da uygulayın. Yine de 2–4. adımları kontrol edin. Installer dışında (ör. Docker Compose ile imaj güncellemesi) dağıttıysanız 1–4. adımların tamamını uygulayın.

1. Hedefli seed'i üretim veritabanı ortamıyla çalıştırın:

```bash
docker run --rm --network aystek-net --env-file /etc/aystek-core/core.env \
  --entrypoint /usr/local/bin/aystek-seed \
  "$(docker inspect --format '{{.Config.Image}}' aystek-core)" --roles
```

   Kaynak kodu olan bir ortamda: üretim değişkenleri dışa aktarılmış halde `go run ./cmd/seed --roles`.

2. Yönetim panelinde **Onay Politikaları** ekranını açın; `b2b.purchase_order` ve `payments.refund` politikalarını kontrol edin.
3. Varsayılan B2B zinciri `store_manager`, `credit_dept` ve `executive` rollerini kullanır. Bu roller hiçbir seed tarafından oluşturulmaz ve kod gönderimde rol adlarının varlığını denetlemez. Rolleri oluşturup kullanıcılara atamadan bu adımlar onaylanamaz ve sipariş bekler. Rolleri oluşturun ya da zinciri mevcut rollerle güncelleyin (`PUT /admin/approval-policies/b2b.purchase_order`; `steps` alanı gönderilir, sunucu rol adlarını doğrular).
4. Etkin B2B zinciri politikası kayıtlı olmayan (veya devre dışı bırakılmış) tenant ve site için satın alma siparişi gönderimi HTTP 400 ve `approvals.error.chain_not_configured` ile başarısız olur (kapalı başarısızlık). Diğer tenant ve site'lar için politikayı Onay Politikaları ekranından kaydedin.

**Onay akışı**

- *İade:* Politika etkinse ve tutar eşiği aşıyorsa `POST /admin/payment-intents/:id/refunds` HTTP 409 döner; gövdede `{"error":"approvals.error.required","approval_request_id":"<id>"}` bulunur. Talebi açandan farklı, `payments:approve` izni olan bir kullanıcı `POST /admin/approval-requests/<id>/approve` ile onaylar. Aynı iade isteği aynı tutar ve `approval_request_id` ile yeniden gönderilir; önceki deneme başarısız kapandığı için aynı idempotency anahtarı da kullanılabilir. Onay tüketilir ve iade tamamlanır. Tutar değişirse, talep onaylı değilse veya onay zaten kullanılmışsa HTTP 409 `approvals.error.not_granted` döner. Sağlayıcı iadeyi gerçekleştirmeden hata verirse onay yeniden `approved` olur.
- *B2B sipariş:* Gönderimde zincir, politikadaki tutar eşiklerine göre üretilir; hiçbir eşik aşılmazsa sipariş gönderimde otomatik onaylanır. Adımlar sıralıdır: yalnızca ilk bekleyen adım karara bağlanır ve onu yalnızca o adımın rolünü taşıyan kullanıcı onaylayabilir veya reddedebilir; rotalar ayrıca `b2b_dealers:approve` izni ister. Siparişi gönderen kişi onaylayamaz ve reddedemez.

**Menü davranışı**

Onay Politikaları ve Onay Talepleri menü öğeleri `approval_policies:read` veya `approval_requests:read` izni olan kullanıcılara gösterilir; bu öğeler rol adına bağlı değildir. Bir grup, görünür bir alt öğesi olduğu sürece görünür kalır. Bu nedenle `editor` rolü, grup düzeyinde yalnızca `admin` rolüne tanımlı olan ama alt öğeleri editöre açık grupları da görür: sipariş (yorumlar), ödemeler (storefront ödemeleri) ve operasyon (telematik, PMS). Bu girişler modül kapısı ve görünüm profili ile ayrıca daraltılabilir.

---

## 11. Rollback Prosedürü

### 11.0 aystek-core-installer ile Rollback (Önerilen CLI Yolu)

```bash
# Önceki sürüme geri dön — otomatik yedek alınır
sudo aystek-core-installer rollback
```

**`aystek-core-installer rollback` şunları otomatik gerçekleştirir:**

1. **Otomatik yedek** — geri alma başlamadan önce `pre-rollback` etiketiyle mevcut DB yedeği alınır
2. Servisleri durdurur
3. Symlink'i `previous` sürüme döndürür
4. Servisleri önceki sürümle yeniden başlatır
5. Sağlık kontrolü yapar

---

> **Not:** Rollback uygulama sürümünü geri döndürür. Güncellemenin tam seed zincirinde yeniden yazdığı kayıtlar (ör. `admin` rolünün izin seti) önceki değerlerine dönmez. Gürsel GPS yapılandırmaları güncellemede değişmediği için bunlar için geri yükleme gerekmez. Gerekirse güncelleme öncesi yedekten geri yükleyin (Bölüm 10.0).

### 11.1 Panel Üzerinden Rollback

1. **Sistem > Rollback** menüsüne gidin
2. Hedef sürümü seçin
3. Onay penceresinde **"Rollback'i Onayla"** yazıp onaylayın
4. İlerlemeyi takip edin

### 11.2 CLI ile Rollback

```bash
# Mevcut ve önceki sürümü göster
ls -la /opt/aystek-core/releases/
readlink /opt/aystek-core/current      # aktif sürüm
readlink /opt/aystek-core/previous     # rollback hedefi

# Rollback çalıştır (aystek-core-installer — önerilen, otomatik yedek alır)
sudo aystek-core-installer rollback
```

### 11.2.1 Docker Compose ile Rollback

```bash
cd /opt/aystek-docker

# Önceki stabil image tag'ine dönün
sed -i 's|AYSTEK_CORE_IMAGE=.*|AYSTEK_CORE_IMAGE=ghcr.io/aystek/aystek-core:1.2.0|' .env
sed -i 's|AYSTEK_PANEL_IMAGE=.*|AYSTEK_PANEL_IMAGE=ghcr.io/aystek/aystek-panel:1.2.0|' .env

# Eski sürümü geri al
docker compose pull
docker compose up -d --force-recreate

# Doğrula
docker compose ps
curl -s http://127.0.0.1:${APP_PORT:-19090}/health | python3 -m json.tool
```

### 11.3 Acil Manuel Rollback

Servisler tamamen yanıt vermiyorsa:

```bash
# Servisleri durdur
systemctl stop aystek-core aystek-panel

# Symlink'i önceki sürüme döndür
PREV_VERSION="1.1.0"
ln -sfn /opt/aystek-core/releases/v${PREV_VERSION} /opt/aystek-core/current

# Servisleri başlat
systemctl start aystek-core aystek-panel

# Sağlık kontrolü
curl -s http://127.0.0.1:9090/health
```

---

## 12. Sağlık Kontrolü ve Sorun Giderme

### 12.1 Servis Durumları

```bash
# Tüm AYSTEK servislerinin durumu
systemctl status aystek-core aystek-panel

# Bağımlılık servisleri
systemctl status mongod redis rabbitmq-server nginx
```

### 12.2 Sağlık Endpoint'leri

| Endpoint | Amaç |
|----------|------|
| `GET /health` | Detaylı sağlık (versiyon, servisler, bellek) |
| `GET /health/live` | Canlılık kontrolü (sadece OK/ERROR) |
| `GET /health/ready` | Hazırlık kontrolü (yük dengeleyici için) |
| `GET /admin/system/node` | Sunucu metrikleri (CPU, RAM, disk) |

### 12.3 Yaygın Sorunlar

| Sorun | Kontrol | Çözüm |
|-------|---------|-------|
| MongoDB bağlantısı yok | `mongosh --eval "db.ping()"` | URI ve parolayı kontrol et |
| Redis bağlantısı yok | `redis-cli -a PAROLA ping` | `requirepass` ayarını kontrol et |
| RabbitMQ bağlantısı yok | `rabbitmq-diagnostics ping` | Kullanıcı izinlerini kontrol et |
| API portu yanıt vermiyor | `ss -tlnp | grep 9090` (veya özel `APP_PORT`) | `aystek-core` servisini veya container'ını kontrol et |
| Panel portu yanıt vermiyor | `ss -tlnp | grep 3000` (veya özel `PORT`) | `aystek-panel` servisini veya container'ını kontrol et |
| Bakim modu takılı kaldı | Redis'te key var mı? | `redis-cli DEL system:maintenance_mode` |

### 12.3.1 Docker Modda `aystek-core` Unhealthy

`docker compose up -d` sırasında `dependency failed to start: container aystek-core is unhealthy` hatası alınırsa önce uygulama container logu okunmalıdır. Yeni installer bu bilgiyi hata çıktısına otomatik ekler; manuel kontrol için:

```bash
cd /opt/aystek-docker
docker compose ps
docker compose logs --tail=200 aystek-core
docker compose logs --tail=100 mongodb redis rabbitmq
docker inspect --format '{{json .State.Health}}' aystek-core
```

Kontrol sırası:

1. `core.env` içinde Docker mod için MongoDB adresi `mongodb:27017`, Redis adresi `redis:6379`, RabbitMQ adresi `rabbitmq:5672` olmalıdır. Host portları yalnızca dışarıdan localhost erişimi içindir.
2. Kullanılan core/panel image tag'leri aynı release çizgisine ait olmalıdır; eski panel/core tag karışımı startup veya health contract kırabilir.
3. Local-source kurulumda build edilen image gerçekten güncel kaynak koddan üretilmiş olmalıdır.
4. MongoDB, Redis, RabbitMQ secret değerleri compose container env'i ile `core.env` arasında uyuşmalıdır.
5. İlk açılış seed/reconcile süresi uzun olabilir. Güncel compose template core healthcheck penceresini genişletir; eski binary/template ile kurulum yapılıyorsa stack erken unhealthy olabilir.
5. `/health/live` endpoint'i açılmadan önce seed/startup hatası varsa `aystek-core` logu asıl sebebi verir.

### 12.4 Log Konumları

```bash
# Uygulama logları
tail -f /var/log/aystek-core/core/app.log
tail -f /var/log/aystek-core/panel/app.log
tail -f /var/log/aystek-core/audit.log
tail -f /var/log/aystek-core/migrations/migration.log

# systemd journal (anlık)
journalctl -u aystek-core -n 100 --no-pager
journalctl -u aystek-panel -n 100 --no-pager
```

---

## 13. Kaldırma (Uninstall)

> ⚠️ **Uyarı:** Bu işlem geri alınamaz. Önce veritabanı yedeği alın.

### 13.0 aystek-core-installer ile Kaldırma (Önerilir)

```bash
# Standart kaldırma — uygulama dosyaları kaldırılır, veri saklanır
sudo aystek-core-installer uninstall

# Tam silme — tüm veri ve sırlar da kalıcı olarak silinir (geri alınamaz)
sudo aystek-core-installer uninstall --purge
```

`--purge` bayrağı kullanıldığında installer kritik onay ekranı gösterir; onaylamak için MongoDB veritabanı adını yazmanız gerekir.

---

### 13.1 Manuel Kaldırma (Native)

```bash
# Servisleri durdur ve devre dışı bırak
systemctl stop aystek-core aystek-panel
systemctl disable aystek-core aystek-panel

# Unit dosyalarını kaldır
rm /etc/systemd/system/aystek-*.service
systemctl daemon-reload

# İsteğe bağlı yardımcı symlink varsa kaldır
rm -f /usr/local/bin/aystek-cli

# Uygulama dosyaları — VERİ YEDEĞİ ALINDIKTAN SONRA
rm -rf /opt/aystek-core

# Yapılandırma ve loglar (dikkatli olun)
rm -rf /etc/aystek-core
rm -rf /var/log/aystek-core
rm -rf /var/lib/aystek-core

# Kullanıcıyı kaldır
userdel aystek
```

### 13.2 Manuel Kaldırma (Docker)

```bash
cd /opt/aystek-docker

# Container'ları durdur ve sil
docker compose down

# Volume'leri de silmek isterseniz (geri alınamaz)
docker compose down -v

# Image temizliği (opsiyonel)
docker image rm ghcr.io/aystek/aystek-core:latest || true
docker image rm ghcr.io/aystek/aystek-panel:latest || true

# Çalışma dizinini silmeden önce yedek alın
cd /
rm -rf /opt/aystek-docker
```

---

## 14. Dizin Yapısı Referansı

```
/opt/aystek-core/
├── releases/
│   ├── v1.0.0/                 # Arşiv (rollback için 30 gün saklanır)
│   ├── v1.1.0/                 # Önceki release (previous symlink)
│   └── v2.0.0/                 # Aktif release (current symlink)
├── current -> releases/v2.0.0  # Aktif symlink
├── previous -> releases/v1.1.0 # Rollback symlink
└── shared/
    ├── uploads/                # Tenant yüklemeleri (sürümler arası paylaşılır)
    └── cache/

/etc/aystek-core/
├── core.env                    # API ortam değişkenleri (root:aystek 640)
├── panel.env                   # Panel ortam değişkenleri
├── installer.json              # Installer yapılandırması
└── license.key                 # Lisans dosyası

/var/log/aystek-core/
├── core/app.log                # API logları
├── panel/app.log               # Panel logları
├── audit.log                   # Installer audit logları
└── migrations/migration.log    # DB migrasyon logları

/var/lib/aystek-core/
├── data/                       # Kalıcı uygulama verisi
├── backups/                    # Güncelleme öncesi otomatik yedekler
└── migrations/state.json       # Migrasyon durum takipçisi

/var/backups/aystek-core/
├── db/                         # Veritabanı yedekleri
└── config/                     # Yapılandırma yedekleri
```

---

## 15. AYSTEK Haber Kurulum Yol Haritası

### 15.1 Mimari ve Sıralama

> **Küresel Hat ortak şablon kurulumu (27 Eylül 2026):** Yeni haber/magazin
> domainleri için localde `aystekhaber-com` projesi kopyalanmaz. İlk yeni site,
> AYSTEK Mühendislik (`aystek`)
> altında `kureselhat` / `kureselhat.com` kimliğiyle tanımlandı.
> Frontend runtime marka/domain/SEO/feed sözleşmesi tamamlandı. DNS, TLS ve
> üretim deployment'ı ayrıca doğrulanmalıdır.
>
> Core kaynak dizininde hedef Mongo ortamı ayarlandıktan sonra:
>
> ```bash
> go run ./cmd/seed --news-site=kureselhat
> ```
>
> Aktif `aystek` tenant'ı bulunmalıdır. Komut tüm seed'leri veya reset'i
> çalıştırmaz; tekrarında mevcut ID'leri ve yönetici düzenlemelerini korur.
> Kayıtlar: 1 store, TR/EN içeren 1 ana sayfa, 9 public menü,
> 2 site ayarı ve 1 editoryal profil. Var olan tenant kategorileri kullanılır.
> Yeni haber/yazar/galeri/maç demosu eklenmez. `www.kureselhat.com` aynı siteye
> alias olarak bağlanır. Dil kataloğu tenant'tan miras alınır; siteye özel dil
> kısıtlaması yönetim panelinden ayrıca yapılabilir.
>
> Doğrulama (API adresini kendi ortamınıza göre değiştirin):
>
> ```bash
> curl -fsS -H 'X-Store-Domain: kureselhat.com' \
>   http://localhost:9090/public/pages/tr/anasayfa
> curl -fsS -H 'X-Store-Domain: www.kureselhat.com' \
>   http://localhost:9090/public/pages/tr/anasayfa
> ```
>
> Her iki yanıtta aynı `site_id` ve `Küresel Hat` başlığı beklenir.
> Yeni sitenin haber listesi editör içerik girene kadar boş olmalıdır.
> `public_site_ready=false` yayın hazırlık notudur, bakım modu değildir.
> Ortak frontend'in domain çözümlemesi, marka/SEO/RSS ve tanınmayan domain reddi
> yerel production build ile doğrulandı. Canlıya çıkmadan önce gerçek içerik,
> kurumsal/hukuki sayfalar, logo/favicon, sosyal/analytics ayarları, DNS ve TLS
> testleri tamamlanmalıdır. Aşağıdaki demo adımları tek başına Küresel Hat yayın
> onayı değildir.

Örnek şablon içeriği ayrıca istendiğinde, çalışan Core image'ındaki
`aystek-seed` ile `--news-site=kureselhat --demo-content` çalıştırılır (tam
komut için installer kullanım kılavuzuna bakın). 121 örnek haber, 5 yazar,
6 galeri ve `[ÖRNEK]` etiketli 2025-2026 Süper Lig snapshot'ı yalnız bu siteye
insert-only eklenir. Spor verisi canlı değildir; eski örnek maçlar “canlı”
gösterilmez. Mevcut/pasifleştirilmiş editör kaydı korunur. Manşet altındaki
spor paneli boşsa site header'larıyla
`GET /public/sports/leagues/super-lig` çağrısının `200` ve `active=true`
döndüğünü doğrulayın; `404` ilgili siteye spor snapshot'ı eklenmediğini veya
pasif olduğunu gösterir.

Local geliştirmede iki site aynı `/workspace/node/aystekhaber-com` checkout'ından
ayrı `NEXT_DIST_DIR` ve portlarla çalışır. Sunucuda kod ve runtime state ayrıdır:

| Domain | Source | Install | Config | Log | Port |
|---|---|---|---|---|---:|
| `aystekhaber.com` | `/opt/aystek-public/aystekhaber-com/source` | `/opt/aystek-public/aystekhaber-com` | `/etc/aystek-public/aystekhaber-com` | `/var/log/aystek-public/aystekhaber-com` | 3001 |
| `kureselhat.com` | `/opt/aystek-public/kureselhat-com/source` | `/opt/aystek-public/kureselhat-com` | `/etc/aystek-public/kureselhat-com` | `/var/log/aystek-public/kureselhat-com` | 3002 |
| `balci.aysteknolojileri.com` (e-ticaret, bkz. 15.7) | `/opt/aystek-public/balci-aysteknolojileri-com/source` | `/opt/aystek-public/balci-aysteknolojileri-com` | `/etc/aystek-public/balci-aysteknolojileri-com` | `/var/log/aystek-public/balci-aysteknolojileri-com` | 3003 |

Source dizinleri aynı repository ve immutable sürümden hazırlanır; ortak writable
dizin, env, compose project/container veya port kullanılmaz. Farklı tasarım
localde de ayrı projedir.

`aystekhaber.com`, ayrı veritabanı veya ayrı Core instance'ı kullanmaz. AYSTEK
tenant'ı altındaki `aystekhaber` site/store kaydına bağlı bağımsız Next.js
public-site container'ıdır. Kurulum sırası değiştirilmemelidir:

1. AYSTEK Core'un AYSTEK Haber store ve seed kodunu içeren sürümünü dağıtın.
2. Core seed zincirini `SEED_FULL_RESET=false` ile çalıştırın.
3. Core public posts API'sini tenant/store header'larıyla doğrulayın.
4. `aystekhaber-com` image'ını yayınlayın veya yerel kaynaktan build edin.
5. Installer `public-site` scope ile siteyi kurun.
6. Lokal container, Nginx, TLS ve dış domain kontrollerini tamamlayın.

Core seed'i çalışmadan frontend kurulursa uygulama açılır ancak içerik alanları
boş görünür. Frontend'in başka site içeriğine fallback yapması beklenmez ve
tenant izolasyonu gereği yapılmamalıdır.

### 15.2 DNS ve Sertifika Ön Kontrolü

```bash
getent ahosts aystekhaber.com
curl -fsS https://core.aysteknolojileri.com/health/live
docker network inspect aystek-net >/dev/null
```

`aystekhaber.com` sunucu IP'sine yönelmeli; `80/443` portları Let's Encrypt ve
HTTPS trafiği için açık olmalıdır. `www.aystekhaber.com` da kullanılacaksa DNS,
Nginx `server_name` ve sertifika SAN listesine birlikte eklenmelidir.

### 15.3 Core Seed ve API Doğrulaması

> **Dikkat:** Bu komut tam seed zincirini çalıştırır. Canlı veritabanında yalnızca bu Core sürümü Haber içeriği eklediğinde ve yedek alındıktan sonra çalıştırın. Zincir `admin` rolünün izin setini kod tanımına göre yeniden yazabilir; Gürsel GPS yapılandırmalarının ayar alanlarını değiştirmez, yalnızca eksik yapılandırmaları kod değerleriyle ekler; `editor`, `reviewer` ve `fleet_viewer` rollerinin panelden yapılan izin değişiklikleri korunur (Bölüm 10.4).

```bash
docker run --rm \
  --network aystek-net \
  --env-file /etc/aystek-core/core.env \
  -e SEED_ON_START=true \
  -e SEED_FULL_RESET=false \
  -e SEED_EXIT_AFTER_RUN=true \
  ghcr.io/huseyinkoc-aystek/aystek-core:<core-sürümü>

docker exec aystek-public-aystekhaber-com node -e '
  fetch(process.env.API_BASE_URL + "/public/posts?lang=tr&page=1&limit=100", {
    headers: {"X-Tenant-ID": "aystek", "X-Store-Domain": "aystekhaber.com"}
  }).then(async response => {
    if (!response.ok) throw new Error(`posts: HTTP ${response.status}`);
    console.log((await response.json()).data?.items?.length ?? 0);
  }).catch(error => { console.error(error); process.exit(1); });
'
```

Dış Core Nginx vhost'u `X-Store-Domain` başlığını Core domainine sabitler;
haber sitesine özgü API kontrolü ilgili public-site container'ından yapılır.

Seed şu başlangıç verilerini idempotent olarak uzlaştırır:

- `aystek` tenant'ı altında `aystekhaber` store ve `aystekhaber.com` domaini
- siteye bağlı ana sayfa ve branding ayarları
- public header menüsü ve post placement kayıtları
- 8 haber kategorisi
- Figma şablonundaki 70 test haberi
- Core author API'sinden sunulan 5 yönetilebilir köşe yazarı
- Süper Lig puan durumu ve son maç snapshot'ı
- AYSTEK Haber editoryal görünüm profili
- Türkçe ve İngilizce post/page/category lokalizasyonları

### 15.4 Installer ile Kurulum

Registry kurulumu için installer reposundaki
`docs/examples/aystekhaber-registry.json`, local-source kurulumu için
`docs/examples/aystekhaber-local-source.json` kullanılır. Küresel Hat karşılıkları
`docs/examples/kureselhat-registry.json` ve
`docs/examples/kureselhat-local-source.json` dosyalarıdır. GitHub'dan otomatik
ilk klonlama ve AI incelemeli güncelleme için ayrıca
`docs/examples/kureselhat-github-ai.json` vardır. Sunucuya yüklenen temiz,
commit edilmiş yerel Git deposu için
`docs/examples/kureselhat-uploaded-git-ai.json` alternatiftir.

```bash
# Registry image ile; üretimde latest yerine immutable semver kullanın.
sudo aystek-core-installer install --non-interactive \
  --config /root/aystekhaber-installer.json

# Geliştirme/aynı host kaynağı ile alternatif:
sudo aystek-core-installer install --non-interactive \
  --config docs/examples/aystekhaber-local-source.json

# Aynı renderer, ayrı server deployment
sudo aystek-core-installer install --non-interactive \
  --config docs/examples/kureselhat-local-source.json

# GitHub profili: önce örneği /root altına kopyalayıp TLS e-postası, Core URL,
# DNS alias ve SSH salt-okunur deploy key erişimini doğrulayın.
sudo aystek-core-installer install --non-interactive \
  --config /root/kureselhat-installer.json
```

Temel sözleşme:

| Alan | Değer |
|---|---|
| `scope` | `public-site` |
| `public_site_domain` | `aystekhaber.com` |
| `public_site_domain_aliases` | `["www.aystekhaber.com"]` (DNS hazırsa) |
| `store_domain` | `aystekhaber.com` |
| `tenant_id` | `aystek` |
| `public_site_port` | `3001` |
| Registry image | `ghcr.io/huseyinkoc-aystek/aystekhaber-com:<semver>` |
| Local image | `aystek/public-aystekhaber-com:local` |
| Install path | `/opt/aystek-public/aystekhaber-com` |
| `tls_email` | Let's Encrypt kayıt e-postası; zorunlu |

Küresel Hat profilinde domain/store `kureselhat.com`, alias
`www.kureselhat.com`, port `3002`, install path
`/opt/aystek-public/kureselhat-com` ve source path
`/opt/aystek-public/kureselhat-com/source` olur. Local ve AI güncelleme
image'ları domain slug'ına özel etiketlenir; env ve filesystem/runtime kimliği
de paylaşılmaz.

Docker sağlık kontrolleri panel için konteyner içindeki
`http://127.0.0.1:<panel_port>/api/health`, public site için
`http://127.0.0.1:3000/api/health` adresini kullanır. Bu rotalar oturum ve
Core API çağrısı olmadan `{"status":"ok"}` döner. Bu süreç kontrolüdür; domain/TLS ve Core API
erişimi ayrıca test edilir. Güncelleme öncesinde veritabanı, yapılandırma ve
`/var/lib/aystek-core/uploads` alanının yedeğini alıp doğrulayın; upload
bind mount'unu yeniden oluşturmayın veya temizlemeyin.
Core, internetten proxy ile gelen public API istekleri için normal IP sınırını
korurken Docker içi doğrudan renderer isteklerine ayrı
`PUBLIC_INTERNAL_RATE_REQ` (varsayılan 600/dakika) kotası uygular. Bu değer
Core env dosyasından ayarlanabilir; proxy başlığı taşıyan istekler iç kotaya
geçmez.
Manuel TLS kullanılıyorsa Nginx `ssl_certificate` yoluna yaprak + ara
sertifikaları içeren fullchain dosyası verin. Ayrı
`ssl_trusted_certificate` dosyası istemciye zinciri tek başına göndermez;
`nginx -t` ve sertifika doğrulamalı `curl` ile test edin.
`install --dry-run` durum dosyası veya audit logu yazmaz ve henüz kurulmamış
servisin portunu yoklamaz; gerçek kurulum
sonrası `validate` ile canlı kontrolleri gerçekleştirin.
Public-site `update` MongoDB yedeği veya seed çalıştırmaz; başarılı imaj
etiketi ve kaynak yolu site state dosyasına yazılır. Core/admin güncellemesinin
veritabanı yedeği ayrı işlem olarak sürer.

GitHub veya yüklenen Git profili ilk kurulumda mevcut klasörü silmeden yalnız
boş checkout'a klonlar; sonraki manuel `update` temiz checkout üzerinde
fast-forward yapar. Yükleme yolu kullanılıyorsa
`/srv/aystek-source/kureselhat-com` tam bir Git reposu olmalı;
`git -C /srv/aystek-source/kureselhat-com status --porcelain` boş dönmeli ve
son değişiklikler `main` dalına commit edilmiş olmalıdır. Düz veya yarım
yüklenmiş klasör otomatik dağıtılmaz. İki remote aynı checkout'ta otomatik
birleştirilmez.
Başka remote/dal, yerel değişiklik veya ayrışmış geçmişte durur. Core için
hedefli `go run ./cmd/seed --news-site=kureselhat` adımı yine gereklidir;
frontend kurulumu bu seed'i otomatik çalıştırmaz.

### 15.4.1 Yerel AI ile Periyodik Public-Site Güncellemesi

Bu özellik GitHub veya commit edilmiş sunucu-içi Git kaynaklı public-site
profilinde kullanılabilir.
Sunucudaki 16 vCPU/31 GiB RAM ve gerçek GPU olmaması nedeniyle küçük bir yerel
açık-ağırlıklı model (profil örneğinde `qwen3:4b`) ile başlayıp CPU gecikmesini
ölçün. Core'daki `AystekAIEngine` şu anda kural/tabanlı bilgi motorudur;
installer'daki Ollama release-advisor aynı eğitilmiş model değildir. Bu
katman ücretli hizmetlerin performansına eşdeğerlik vaat etmez. Ollama'yı
resmî kaynaktan kurun, yalnız localhost'ta tutun, modeli `ollama pull qwen3:4b`
ile indirin. Sonra:

Installer reposundaki `docs/systemd/ollama-aystek.service` örneği modeli yalnız
`127.0.0.1:11434` üzerinde çalıştırır; CPU kotası 4 çekirdek eşdeğeri ve bellek
sınırı 8 GiB'dır. Önce `ollama` servis kullanıcısını ve yazılabilir
`/var/lib/ollama/models` dizinini hazırlayın. `docs/examples/kureselhat-auto-update-ssh.conf`
dosyasındaki GitHub deploy anahtarı yalnız Küresel Hat servis instance'ına
drop-in olarak bağlanmalıdır; anahtarın kendisi repoya konmaz.

```bash
sudo aystek-core-installer auto-update --non-interactive \
  --config /etc/aystek-public/kureselhat-com/installer.json

# İlk kabul testinden sonra, açık operatör kararıyla:
sudo install -m 0644 docs/systemd/aystek-public-auto-update@.service /etc/systemd/system/
sudo install -m 0644 docs/systemd/aystek-public-auto-update@.timer /etc/systemd/system/
sudo install -d -m 0755 /etc/systemd/system/aystek-public-auto-update@kureselhat-com.service.d
sudo install -m 0644 docs/examples/kureselhat-auto-update-ssh.conf \
  /etc/systemd/system/aystek-public-auto-update@kureselhat-com.service.d/ssh-key.conf
sudo systemctl daemon-reload
sudo systemctl enable --now aystek-public-auto-update@kureselhat-com.timer
```

Timer varsayılan kurulumda etkinleşmez. Güncelleme servisi değişen dosyaları
ve sınırlı kod farkını yerel modele değerlendirir; model yoksa veya altyapı/bağımlılık dosyaları
değişmişse dağıtımı tutar. Onaylanan SHA ayrı image'a build edilir, sağlıksız
olursa eski Compose/image geri alınır. CI, korumalı branch, editoryal onay ve
üretim yük testi yine gereklidir. Sitenin yayın hazır olma durumu bu
otomasyonla değişmez; içerik/hukuk/marka/DNS/TLS işleri ayrıca tamamlanmalıdır.

Installer, sertifika henüz yoksa önce yalnız HTTP vhost üretir. Certbot primary
domain ile tüm alias'ları aynı SAN sertifikasına aldıktan sonra public-site
vhost'unu HTTPS olarak yeniden yazar, `nginx -t` çalıştırır ve Nginx'i reload
eder. Bu nedenle Core API/panel domainleri public-site sertifikasına karışmaz.
Staging kurulumunda Nginx/TLS kapalıysa DNS hazır olduğunda bayraklar açılıp
`update` çalıştırılır; installer yalnız ilgili public-site vhost/sertifika
adımlarını ekler, Core veritabanını veya seed'ini çalıştırmaz. Let’s Encrypt
öncesi hem kök domain hem `www` alias'ı yetkili DNS ve genel çözücülerde doğru
sunucu IP'sine çözülmeli, `tls_email` gerçek bir iletişim adresi olmalıdır.

### 15.5 Sağlık ve İçerik Kontrolü

```bash
sudo aystek-core-installer validate \
  --config /root/aystekhaber-installer.json
docker compose -f /opt/aystek-public/aystekhaber-com/docker-compose.yml ps
curl -fsS http://127.0.0.1:3001/ >/dev/null
curl -fsS https://aystekhaber.com/ >/dev/null
curl -fsS https://aystekhaber.com/robots.txt >/dev/null
curl -fsS https://aystekhaber.com/sitemap.xml >/dev/null
curl -fsS https://aystekhaber.com/rss.xml >/dev/null
```

Tarayıcı kontrolünde logo, header menüsü, son dakika bandı, manşet, sağ haberler,
çok okunanlar, reklam/sponsor, video, şehir sekmeli yerel haberler, son maçlar,
puan tablosu, Core author API'sinden gelen yazar kartları/profilleri, mobil
uygulama alanı, e-bülten, footer, mobil menü ve haber detay rotası
doğrulanmalıdır. İçerik gelmiyorsa container içindeki
`API_BASE_URL`, `AYSTEK_STORE_DOMAIN` ve `AYSTEK_TENANT_ID` değerlerini kontrol
edin.

### 15.6 Güncelleme ve Geri Dönüş

Public-site güncellemesi yeni image'ı pull eder veya local source'u yeniden
build eder, container'ı yeniden oluşturur ve ardından validate çalıştırır:

```bash
sudo aystek-core-installer update --non-interactive \
  --config /root/aystekhaber-installer.json
```

Core tarafında yeni seed içeriği varsa update sonrasında 15.3'teki seed komutu
`SEED_FULL_RESET=false` ile tekrar çalıştırılır. Bu tam zincirdir; Bölüm 15.3'teki uyarıya uyun.

Public-site için otomatik önceki image seçimi henüz `rollback` sözleşmesine
dahil değildir. Güvenli geri dönüş için image'ları semver ile yayınlayın,
config'teki `public_site_image_tag` değerini bilinen sağlam sürüme çekin ve
installer `update` + `validate` çalıştırın. Bu işlem site verisini silmez.
`auto-update` kendi güncellemesinde sağlık kontrolü başarısız olursa eski
image'ı geri yükler; genel `rollback` komutunu genişletmez.

Eksiksiz operatör komutları ve sorun giderme ayrıntıları:
`/workspace/go/aystek-core-installer/docs/KULLANIM_KILAVUZU.md`, Bölüm 21.

---

### 15.7 Balcı (Rucher du VELAY) E-Ticaret Sitesi

`balci.aysteknolojileri.com`, haber sitelerinden farklı olarak **ayrı bir Next.js projesidir**
(`/workspace/node/balci-aysteknolojileri-com`; Core CMS + e-ticaret altyapısına bağlı mağaza).
Kurulum sırası ve sözleşme haber siteleriyle aynıdır (kendi dizin/port/image/state'i: 3003),
ek olarak:

1. **Seed (Core görüntüsü ile, hedefli, reset yok):**
   `aystek-seed --public-site=balci` (store, ana sayfa, diller, görünüm profili, Fransa KDV
   kuralları, pasif ödeme seçenekleri). Yalnız test/UAT için örnek katalog: `--demo-catalog`
   (`[DEMO]` ürünler, kargo seçenekleri, demo SEPA hesabı — **üretimde çalıştırmayın**).
   Mağaza yöneticisi hesabı: `aystek-seed --public-site=balci --create-manager` (rol `balci_manager`,
   yalnız balci sitesine kapsamlı; menü rolleri ve `balci_shop_manager` görünüm profili eşitlenir; parola
   `--manager-password`/`SEED_MANAGER_PASSWORD` ile verilir, yoksa rastgele üretilip terminale bir kez yazılır).
2. **Core ortamı:** `APP_ENV=production` (mock ödeme adaptörü kapanır; UAT için `uat`/`staging`),
   **SMTP açık ve doğrulanmış** (hesap doğrulama ve şifre sıfırlama e-postaları Core'dan gider; SMTP yoksa
   müşteriler giriş yapamaz), `PUBLIC_RATE_REQ` ziyaretçi başına değerlendirilir (BFF ziyaretçi IP'sini iletir).
3. **Nginx:** installer'ın public-site vhost'u `X-Forwarded-For`'u `$remote_addr` ile **üzerine yazar**
   (Core oran sınırı ilk girdiyi anahtar alır). Önüne CDN/ters vekil konulursa gerçek istemci IP'sini
   güvenilir tek kaynaktan (CDN'in kendi başlığı) `$remote_addr`'a taşıyan `real_ip` ayarı şarttır.
4. **Yönetim paneli:** ürünleri siteye yayınla (*Commerce → Storefront*), ödeme seçeneklerini aç
   (*Storefront Payments*: IBAN/BIC; Stripe/PayPal ağ geçidi anahtarları + `webhook_secret` / `webhook_id`),
   gerçek iletişim bilgisi, logo ve yasal sayfalar (`mentions-legales`, `cgv`, `confidentialite`).
5. Profiller: `docs/examples/balci-registry.json`, `docs/examples/balci-local-source.json`
   (TLS e-postası `info@aysteknolojileri.com`).

```bash
docker run --rm --network aystek-net --env-file /etc/aystek-core/core.env \
  --entrypoint /usr/local/bin/aystek-seed \
  "$(docker inspect --format '{{.Config.Image}}' aystek-core)" --public-site=balci
sudo aystek-core-installer install --non-interactive \
  --config docs/examples/balci-local-source.json
curl -fsS -H 'X-Store-Domain: balci.aysteknolojileri.com' -H 'X-Tenant-ID: aystek' \
  http://127.0.0.1:9090/public/storefront/payment-methods
```

Sağlık: `https://balci.aysteknolojileri.com/api/health`. Kayıtlı müşteri akışı (`/compte`) ve
tüm ödeme akışları, gerçek ödeme yapılmadan önce sandbox anahtarlarıyla uçtan uca denenmelidir.

## Hızlı Başvuru — Sık Kullanılan Komutlar

```bash
# ── aystek-core-installer komutları ──────────────────────────────────────────────
sudo aystek-core-installer status                  # tüm servislerin durumu
sudo aystek-core-installer update                  # güncelle (otomatik yedek alır)
sudo aystek-core-installer rollback                # geri al (otomatik yedek alır)
sudo aystek-core-installer backup                  # manuel yedek al
sudo aystek-core-installer backup --list           # yedekleri listele
sudo aystek-core-installer backup --restore /var/backups/aystek-core/db/2024-01-19_14-23-11_pre-update
sudo aystek-core-installer repair                  # servisler ve yapılandırmayı onar
sudo aystek-core-installer validate                # kurulum bütünlüğünü doğrula
sudo aystek-core-installer version                 # installer sürümü

# ── Servis durumu (native) ───────────────────────────────────────────────────
systemctl status aystek-core aystek-panel

# ── Docker durumu (compose) ──────────────────────────────────────────────────
cd /opt/aystek-docker
docker compose ps
docker compose restart aystek-core
docker compose logs -f aystek-core

# ── Core hedefli seed (veri silmez; container içinde, bkz. Bölüm 10.4) ───────
docker run --rm --network aystek-net --env-file /etc/aystek-core/core.env \
  --entrypoint /usr/local/bin/aystek-seed \
  "$(docker inspect --format '{{.Config.Image}}' aystek-core)" --roles

# ── Sağlık kontrolü ──────────────────────────────────────────────────────────
curl -s http://127.0.0.1:9090/health | python3 -m json.tool   # veya özel APP_PORT

# ── Canlı log izleme ─────────────────────────────────────────────────────────
journalctl -u aystek-core -f
```
