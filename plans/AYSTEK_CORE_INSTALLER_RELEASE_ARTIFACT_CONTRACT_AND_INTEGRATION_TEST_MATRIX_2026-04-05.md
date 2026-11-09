# AYSTEK Core Installer Release Artifact Contract and Integration Test Matrix

**Tarih:** 2026-04-05  
**Durum:** Drafted for Execution  
**Bağlı Plan:** [AYSTEK_CORE_INSTALLER_AND_SERVER_SETUP_ALIGNMENT_EXECUTION_PLAN_2026-04-05.md](/workspace/aystek-core-dev-docs/plans/AYSTEK_CORE_INSTALLER_AND_SERVER_SETUP_ALIGNMENT_EXECUTION_PLAN_2026-04-05.md)

## Amaç

`aystek-core-installer` için release artifact beklentilerini tek sözleşmede netleştirmek ve bugünkü test boşluğunu uygulanabilir bir integration test matrisine dönüştürmek.

Bugünkü ana problem kodun tamamen belgesiz olması değil; artifact beklentileri ile gerçek test kapsamının dağınık ve örtük kalmasıdır.

## 1. Mevcut Gerçeklik

Installer bugün iki ana deployment akışı taşıyor:

1. **Native**
   - `CoreBinaryPath`
   - `PanelBuildPath`
   - release dizini oluşturma
   - `current` / `previous` symlink rotasyonu

2. **Docker**
   - `CoreImageTag`
   - `PanelImageTag`
   - `docker-compose.yml` + `.env` üretimi

Bootstrap tarafı ise release artifact olarak doğrudan installer binary archive'ını indiriyor:

- `aystek-core-installer_<version>_<os>_<arch>.tar.gz`

Burada açık kalan nokta şu: uygulama release artifact sözleşmesi installer repo içinde tam ve merkezi olarak yazılı değil. Özellikle native deployment için `aystek-core` binary + panel build nasıl paketlenmeli sorusu operatör veya release engineering bilgisine bırakılmış durumda.

## 2. Canonical Artifact Contract

### 2.1 Installer Release Artifact

Bootstrap ve dağıtım sistemlerinin beklediği canonical artifact:

- `aystek-core-installer_<version>_<os>_<arch>.tar.gz`

Arşiv içinde beklenen binary:

- `aystek-core-installer`

Örnek:

- `aystek-core-installer_2.0.0_linux_amd64.tar.gz`
- `aystek-core-installer_2.0.0_linux_arm64.tar.gz`

### 2.2 Native Application Release Contract

Native mod için installer şu girdileri net bekler:

- `release_version`
- `core_binary_path`
- `panel_build_path`

Bu akışta canonical runtime yerleşimi:

- `/opt/aystek-core/releases/v{release_version}/bin/aystek-core`
- `/opt/aystek-core/releases/v{release_version}/panel/`
- `/opt/aystek-core/current`
- `/opt/aystek-core/previous`

Bu yüzden release engineering için canonical native paketleme çıktısı şu bilgileri güvenilir biçimde sağlamalıdır:

1. `aystek-core` Linux binary
2. panel production build directory
3. sürüm numarası
4. checksum / provenance bilgisi

### 2.3 Docker Application Release Contract

Docker mod için canonical input:

- `core_image_tag`
- `panel_image_tag`

Örnek:

- `ghcr.io/aystek/aystek-core:2.0.0`
- `ghcr.io/aystek/aystek-panel:2.0.0`

### 2.4 Scope Contract

Installer şu scope/mod kombinasyonlarını destekler:

1. `scope=admin`, `mode=native`
2. `scope=admin`, `mode=docker`
3. `scope=admin`, `mode=mongodb-only`
4. `scope=tenant-site`, `mode=native`
5. `scope=tenant-site`, `mode=docker`

`tenant-site` akışında canonical ek girdiler:

- `tenant_id`
- `tenant_name`
- `site_id`
- `site_name`
- `tenant_admin_api_url`

## 3. Bugünkü Boşluklar

### 3.1 Artifact Sözleşmesi Boşluğu

Bugün installer repo içinde:

- installer archive adı net
- docker image input net
- native deploy runtime yerleşimi net

Ama şu iki başlık tam kapalı değil:

1. native uygulama release artefact'ının dış formatı
2. release engineering pipeline'ının hangi checksum/provenance alanlarını zorunlu kıldığı

### 3.2 Test Boşluğu

`go test ./...` geçiyor, fakat repo içinde gerçek integration test katmanı yok. Bu nedenle şu riskler hâlâ belge/varsayım seviyesinde:

- bootstrap archive gerçekten doğru isimle açılıyor mu
- generated template + systemd + nginx + docker env dosyaları host üzerinde beklendiği gibi davranıyor mu
- update/rollback/backup akışları gerçek süreçte beklendiği gibi çalışıyor mu
- tenant-site isolation aynı host üstünde port/path/state çakışması üretmiyor mu

## 4. Integration Test Matrix

### 4.1 Test Katmanları

#### A. Unit Tests

Hızlı, deterministik, root gerektirmeyen testler.

Kapsam:

- config default/validation
- state manager
- tenant registry
- rollback safety checks
- template render helpers
- backup label/path helper'ları

#### B. Local Smoke Tests

CI veya geliştirici makinesinde çalışabilecek, sınırlı yan etki üreten testler.

Kapsam:

- `go build ./cmd/aystek-core-installer`
- `bash -n scripts/bootstrap.sh`
- release archive name consistency
- template rendering dry-run smoke

#### C. Host Integration Tests

VM/container tabanlı, root/systemd/network gerektiren testler.

Kapsam:

- native install
- native update
- native rollback
- scheduled backup timer install
- nginx vhost generation
- firewall integration

#### D. Scenario Tests

Daha geniş operasyon akışları.

Kapsam:

- mongodb-only
- tenant-site first install
- tenant-site second install on same host
- docker mode install/update
- uninstall vs uninstall --purge

## 5. Test Matrisi

| ID | Senaryo | Katman | Beklenen Sonuç |
|----|---------|--------|----------------|
| IT-01 | Installer binary build | Local smoke | `aystek-core-installer` binary üretilir |
| IT-02 | Bootstrap shell syntax | Local smoke | script parse hatası vermez |
| IT-03 | Bootstrap archive naming | Local smoke | indirilen artifact adı canonical formata uyar |
| IT-04 | Native config render | Unit/local smoke | `core.env`, `panel.env`, systemd ve nginx dosyaları render edilir |
| IT-05 | Native first install | Host integration | dizinler, config, symlink, service unit'leri oluşur |
| IT-06 | Native update | Host integration | pre-update backup + symlink rotasyonu + validate çalışır |
| IT-07 | Native rollback | Host integration | pre-rollback backup + previous symlink dönüşü çalışır |
| IT-08 | Rollback migration limitation | Unit/host | DB migration rollback yapılmadığı açıkça doğrulanır |
| IT-09 | Docker install | Host integration | compose/env üretilir ve stack kalkar |
| IT-10 | MongoDB-only install | Host integration | yalnız MongoDB ve backup path kurulumu yapılır |
| IT-11 | Tenant-site native install | Host integration | scoped path/db/state/registry entry üretilir |
| IT-12 | Second tenant-site install | Host integration | ilk tenant ile çakışmadan ikinci instance açılır |
| IT-13 | Backup command manual/list/restore | Host integration | yedek alınır, listelenir, restore path çözülür |
| IT-14 | Uninstall safe | Host integration | app kaldırılır, data/config korunur |
| IT-15 | Uninstall purge | Host integration | kritik onaydan sonra veri ve config kaldırılır |

## 6. Önceliklendirme

### P0

- IT-03 Bootstrap archive naming
- IT-05 Native first install
- IT-06 Native update
- IT-07 Native rollback
- IT-11 Tenant-site native install

### P1

- IT-09 Docker install
- IT-10 MongoDB-only install
- IT-13 Backup command lifecycle
- IT-14 Safe uninstall

### P2

- IT-12 Multi-tenant-site collision test
- IT-15 Purge destructive path
- provenance/checksum expansion

## 7. Uygulama Önerisi

### İlk uygulanacak test harness

1. ephemeral VM veya privileged container tabanlı smoke harness
2. release archive name assertion script
3. native install/update/rollback golden scenario
4. tenant-site isolation scenario

### Repo içi önerilen yerleşim

- `/workspace/go/aystek-core-installer/tests/integration/`
- `/workspace/go/aystek-core-installer/tests/fixtures/`
- `/workspace/go/aystek-core-installer/tests/scripts/`

## 8. Çıkış Kriteri

Bu slice tamamlandı sayılabilmesi için:

1. canonical artifact contract yazılı olmalı
2. test matrix plan seviyesinde onaylı olmalı
3. en az P0 senaryolar için harness tasarımı açılmış olmalı
4. kurulum rehberleri bu contract ile çelişmemeli

## 9. Durum Güncellemesi (2026-04-05)

P0 harness iskeleti açıldı:

- `/workspace/go/aystek-core-installer/tests/integration/`
- `/workspace/go/aystek-core-installer/tests/fixtures/`
- `/workspace/go/aystek-core-installer/tests/scripts/`

Scaffold edilen akışlar:

- `IT-03`
- `IT-05`
- `IT-06`
- `IT-07`
- `IT-11`

Ek olarak fresh-DB bootstrap admin akışı installer wizard ve Mongo seed katmanına eklendi; böylece `admin` ve `tenant-site` scope kurulumlarda boş veritabanı için ilk admin kullanıcısı kurulum sırasında tanımlanabiliyor.
