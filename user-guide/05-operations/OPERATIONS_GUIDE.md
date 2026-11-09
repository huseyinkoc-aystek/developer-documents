---
Title: Operations Guide
Version: 2.5.0
Status: Enterprise Certified
Last Updated: 2026-10-28
Related Documents:
  - ./DEPLOYMENT_SCENARIOS.md
  - ./API_GATEWAY_SECRET_ROTATION_RUNBOOK.md
  - ./API_GATEWAY_SECURITY_AUDIT_CHECKLIST.md
  - ./API_GATEWAY_SLO_SLA_ALERTING.md
  - ./INTEGRATIONS_EVENTBRIDGE_SMOKE_RUNBOOK.md
  - ./INTEGRATIONS_CHAIN_RUNBOOK_QA_MATRIX.md
  - ./INTEGRATIONS_SCHEDULED_PULL_RUNBOOK_QA_MATRIX.md
  - ./INTEGRATIONS_SHARED_SERVICE_RUNBOOK_QA_MATRIX.md
  - ./INTEGRATIONS_SHARED_SERVICE_ALERT_THRESHOLDS.md
  - ./GITLAB_PROTECTED_ENV_SETUP.md
  - ../04-api/API_REFERENCE.md
  - ../01-architecture/CONTENT_OPERATIONS_ARCHITECTURE.md
---

# Table of Contents
1. [Runtime Health](#1-runtime-health)
2. [Content Ingestion Monitoring](#2-content-ingestion-monitoring)
3. [Operational Playbooks](#3-operational-playbooks)
4. [Context Safety Checks](#4-context-safety-checks)
5. [Operational Baseline](#5-operational-baseline)

## 1. Runtime Health
- `GET /healthz`: liveness-level health check
- `GET /readyz`: readiness-level dependency check
- `GET /metrics`: observability/monitoring surface

## 2. Content Ingestion Monitoring
Track at minimum:
- Connector health status (`healthy/warning/critical`)
- Job success/failure rates
- Error thresholds and retry patterns
- Duplicate rate trend
- Moderation backlog size and aging

Recommended operational SLOs:
- Connector run success ratio per 24h
- Median moderation lead time
- Duplicate suppression effectiveness

## 3. Operational Playbooks
### 3.1 Connector Failure Handling
1. Detect connector in `warning/critical` state.
2. Review recent run errors and trace IDs.
3. Trigger connection test and dry-run if supported.
4. Rotate secret if credential issue suspected.
5. Disable connector via endpoint if recovery fails.

### 3.2 Mass Rejection Workflow
1. Filter moderation queue by connector/time window.
2. Bulk reject with mandatory reason.
3. Create incident note and audit reference.
4. Re-enable only after source validation.

### 3.3 Injection Code Rollback
1. Identify problematic code key and version.
2. Soft delete or deactivate current version.
3. Re-apply previous approved version.
4. Validate rendering and CSP/security behavior.

### 3.4 Emergency Feature Flag Disable
1. Toggle critical feature flag off in site-scoped settings.
2. Verify runtime behavior and queue stabilization.
3. Record governance incident and resolution notes.

### 3.5 Integrations EventBridge Smoke (All Modules)
1. Use `05-operations/INTEGRATIONS_EVENTBRIDGE_SMOKE_RUNBOOK.md`.
2. Validate event-key based dynamic module binding from Integrations Hook screen.
3. Confirm queue consume + execution success path.
4. Confirm retry + DLQ behavior with controlled failure.

### 3.6 Integrations Chain Smoke & QA Matrix
1. Use `05-operations/INTEGRATIONS_CHAIN_RUNBOOK_QA_MATRIX.md`.
2. Validate chain graph (`valid=true`) before runtime tests.
3. Run chain test-run in both `dry_run=true` and `dry_run=false`.
4. Verify step-level observability (node status/duration/error) in execution details.
5. Confirm branch/retry/DLQ/idempotency scenarios from QA matrix.

### 3.7 GitLab Protected Environment Security Check
1. Apply UI settings from `05-operations/GITLAB_PROTECTED_ENV_SETUP.md`:
   - Protected Environment: `production`
   - Variable: `INTEGRATIONS_SMOKE_TOKEN` as `protected + masked`
2. Verify by API script:
   - `go/aystek-core/tools/ops/gitlab_verify_integrations_smoke_security.sh`

### 3.8 Integrations Scheduled Pull Smoke & QA Matrix
1. Use `05-operations/INTEGRATIONS_SCHEDULED_PULL_RUNBOOK_QA_MATRIX.md`.
2. Validate schedule create/update + test-run from admin UI.
3. Validate scheduler -> pull-run creation -> worker execution path.
4. Verify pull-run detail ingestion summary/issues visibility.
5. Confirm retry/backoff/DLQ and success-only cursor commit scenarios.

### 3.9 Integrations Shared Service Smoke & QA Matrix
1. Use `05-operations/INTEGRATIONS_SHARED_SERVICE_RUNBOOK_QA_MATRIX.md`.
2. Validate shared service/grant CRUD ve policy hardening alanları.
3. Validate execution hattında scope/method/path/limit/secret/mTLS enforcement.
4. Validate usage-summary endpointleri ve UI warning paneli.
5. Pipeline smoke için script çalıştırın:
   - `05-operations/scripts/integrations_shared_service_smoke.sh`

### 3.10 Integrations Shared Service Alert Thresholds
1. Use `05-operations/INTEGRATIONS_SHARED_SERVICE_ALERT_THRESHOLDS.md`.
2. Prometheus alert rulesını environment bazlı eşiğe göre uygulayın.
3. Grafana dashboard panel setini shared service/grant outcome etiketleriyle yayınlayın.

### 3.11 API Gateway Secret Rotation & Audit Export
1. Use `05-operations/API_GATEWAY_SECRET_ROTATION_RUNBOOK.md`.
2. `POST /admin/api-clients/:id/rotate-secret` ile rotation tetikleyin.
3. Rotation sonrasi token yenileme dogrulamasini yapin (`/auth/oauth2/token`).
4. `GET /admin/api-clients/:id/audit-logs/export?format=csv` ile ilk trafik penceresini disariya alin.
5. `API_GATEWAY_AUDIT_RETENTION_DAYS` degerini regulator kurallarina gore ortam bazli ayarlayin.

### 3.12 API Gateway Security Audit
1. Use `05-operations/API_GATEWAY_SECURITY_AUDIT_CHECKLIST.md`.
2. SSRF deny, injection guard ve replay kontrollerini audit penceresinde dogrulayin.
3. Audit export + metrik trendi ile denetim kanitlarini saklayin.

### 3.13 API Gateway SLO/SLA & Alerting
1. Use `05-operations/API_GATEWAY_SLO_SLA_ALERTING.md`.
2. Prometheus ruleset: `go/aystek-core/deployments/prometheus/alerts.yml`.
3. Grafana dashboard: `go/aystek-core/deployments/grafana/dashboards/api-gateway-third-party.json`.
4. Haftalik SLO review + aylik SLA raporu cikarin.

### 3.14 AYSTEK Haber Public-Site Operasyonu

1. İlk kurulum, güncelleme, doğrulama ve kontrollü geri dönüş için
   `05-operations/SUNUCU_KURULUM_REHBERI.md` Bölüm 15'i kullanın.
2. Her release sonrasında Core posts, authors ve sports public uçlarını
   `X-Tenant-ID: aystek` ve `X-Store-Domain: aystekhaber.com` ile doğrulayın.
3. Site tarafında ana sayfa, haber/yazar detayları, e-bülten, `robots.txt`,
   `sitemap.xml` ve `rss.xml` smoke kontrollerini çalıştırın.
4. Let's Encrypt sertifikasının hem primary domaini hem yapılandırılmış
   `public_site_domain_aliases` kayıtlarını kapsadığını kontrol edin.
5. İçerik güncellemesinde `SEED_FULL_RESET=true` kullanmayın. `SEED_FULL_RESET=false` veri silmez, ancak tam seed zinciri bazı kayıtları kod tanımına göre yeniden yazar; canlıda hedefli seed tercih edin (bkz. 3.15).

### 3.15 2026-10-28 Onay politikaları ve tohum (seed) değişiklikleri

Ayrıntılı açıklama ve komutlar `SUNUCU_KURULUM_REHBERI.md` Bölüm 10.4'tedir. Dağıtım sonrası kontrol listesi:

1. Hedefli seed'i çalıştırın: `aystek-seed --roles` (installer `update` tam zinciri `editor`, `reviewer` ve `fleet_viewer` rollerine eksik varsayılan izinleri eklemez; sürüm bu rollere yeni varsayılan izin ekliyorsa komutu çalıştırın; sonucu yine de kontrol edin). Bu komut veri silmez ve `editor`, `reviewer` ve `fleet_viewer` rollerinden izin kaldırmaz; eksik varsayılan izinleri ekler ve `admin` rolünü kod tanımıyla eşitler. Panelden kaldırılmış bir varsayılan izin bu komutla yeniden eklenir.
2. Admin panelinde Onay Politikaları ekranında `b2b.purchase_order` ve `payments.refund` politikalarını doğrulayın. Varsayılan B2B zincirindeki `store_manager`, `credit_dept` ve `executive` rollerini oluşturun veya zinciri mevcut rollerle güncelleyin; aksi halde zincir adımları onaylanamaz.
3. İzleme: `approvals.error.chain_not_configured` (HTTP 400), ilgili tenant/site için B2B politikasının eksik veya devre dışı olduğunu gösterir. `approvals.error.required` (HTTP 409, iade) beklenen bir durumdur; onay talebi açıldığını belirtir. Bekleyen talepler için `GET /admin/approval-requests?status=pending` kullanın.
4. Politika değiştirme: `PUT /admin/approval-policies/<gate>` politikanın tamamını yazar. Eşik, para birimi, izin ve adım alanlarını korumak için mevcut değerleri de gönderin (B2B için `steps` zorunludur). `payments.refund` için `enabled: false` iade onayını kapatır (eski davranış). `b2b.purchase_order` politikasını devre dışı bırakmak PO gönderimini kapatır (kapalı başarısızlık); bu kararı bilinçli verin.
5. Migrasyon aracı (`cmd/migrate_gursel`): çalıştırmadan önce `MIGRATE_SOURCE_URI` ortam değişkenini veya `-src-uri` bayrağını kaynak bağlantı dizesiyle ayarlayın. Kodda varsayılan kaynak yoktur; `SOURCE_MONGO_URI` tanımlıysa değerini geçersiz kılar. Bağlantı dizesini belgelere, ticket'lara veya komut geçmişine yazmayın.
6. Canlı ortamda `SEED_ON_START` değerini kurulumdan sonra `false` tutun. Gürsel GPS sağlayıcı yapılandırmaları (`seedGurselTurizimGpsProviders`) yalnızca eksikse eklenir; mevcut yapılandırmaların kimlik bilgileri, temel URL, aktif bayrağı, görünen ad ve saat dilimi farkı tam zincirde değişmez. Değişiklikleri admin panelinden yapın.

## 4. Context Safety Checks
Operational verification should include:
1. Tenant-required routes reject missing tenant context.
2. Site-required routes reject missing site context.
3. Security events include traceable identifiers.

## 5. Operational Baseline

### Küresel Hat shared news-site deployment checkpoint (2026-09-27)

- Use the existing Core source with `go run ./cmd/seed --news-site=kureselhat`
  and the target Mongo environment; do not run a full seed/reset just to add it.
- The active `aystek` tenant must exist. Canonical or alias domain conflicts
  stop provisioning instead of moving another tenant's domain.
- Defaults are insert-only. Subsequent runs preserve admin changes and deleted
  defaults; they are not migrations that reapply edited template text.
- Verify apex and `www` resolve to the same site and no demo posts are exposed.
- Local development shares `/workspace/node/aystekhaber-com`; production/server
  deployments do not share a writable checkout or runtime state. Use separate
  source/install/config/log directories, ports, containers and vhosts per domain.
- The same immutable image is valid for both sites. Identity must be supplied by
  `SITE_URL`, `SITE_DOMAIN_ALIASES`, `AYSTEK_STORE_DOMAIN` and `AYSTEK_TENANT_ID`.
- Validate homepage canonical/OpenGraph identity, `robots.txt`, `sitemap.xml`,
  `rss.xml`, `manifest.webmanifest`, alias admission and cross-domain rejection.
- Size `PUBLIC_RATE_REQ` for SSR traffic. All Core calls from a renderer may share
  a network identity; a 429 is an operational capacity error, not missing content.
- The development DB was provisioned; no production DNS/TLS or proxy was changed.
- The shared renderer contract is complete. `public_site_ready=false` remains
  descriptive metadata, not an access-control mechanism; set it only after the
  editorial/legal/branding and infrastructure launch checklist is complete.
- The latest local verification reported MongoDB, Redis and RabbitMQ healthy.
  Production health must still be rechecked immediately before rollout.

### General baseline

- Enforce predictable environment settings across dev/stage/prod.
- Keep integration and security-gate checks in release pipeline.
- Track deprecation policy and documentation changelog together.
- Keep worker idempotency checks enabled in all non-local environments.
