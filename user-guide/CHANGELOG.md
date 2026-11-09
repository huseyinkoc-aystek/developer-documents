---
Title: User Guide Changelog
Version: 3.0.8
Status: Enterprise Certified - Production Ready
Last Updated: 2026-06-05
Classification: Internal - Enterprise Documentation
Related Documents:
  - ./README.md
  - ./GOVERNANCE.md
---

# Table of Contents
1. [Versioning Convention](#1-versioning-convention)
2. [Releases](#2-releases)

## 1. Versioning Convention
Documentation versioning follows `MAJOR.MINOR.PATCH`.
- `MAJOR`: structural suite changes or governance baseline shifts.
- `MINOR`: new approved documents or major section expansion.
- `PATCH`: corrections, reference fixes, wording, and consistency updates.

## 2. Releases

## [3.0.8] – Fleet Tally Türkçe Kullanım Kılavuzu

### Added
- `06-modules/FLEET_TALLY_USER_GUIDE_TR.md`
  - Fleet Tally / GPS Çetele Kontrol için kullanıcı ve operasyon odaklı Türkçe rehber eklendi.
  - Ekranlar, kurulum sırası, Gürsel seed kayıtları, sağlayıcılar, zamanlayıcılar, run/job takibi, sonuç kodları ve troubleshooting akışı açıklandı.

### Updated
- `README.md`
  - Modules bölümüne Fleet Tally kullanıcı kılavuzu referansı eklendi.

### Notes
- Bu sürüm dokümantasyon odaklıdır; kullanıcıların `/admin/fleet-tally/*` ekranlarında hangi adımı nasıl takip edeceğini standartlaştırır.

## [3.0.7] – One-Page Guide Deep Dive Expansion

### Updated
- `00-overview/AYSTEK_CORE_ONE_PAGE_PLATFORM_GUIDE.md`
  - Rehber detay seviyesi artırıldı:
    - backend/frontend katman ayrımı
    - modül aileleri ve kritik kapsam
    - integrations/chain/scheduled pull derin özet
    - güvenlik/uyum/SLO-operasyon checklist
    - sık sorunlar ve onboarding akışı

### Notes
- Bu sürüm, tek sayfa rehberin “yönetici özeti” seviyesinden “operasyonel başvuru” seviyesine genişletilmesidir.

## [3.0.6] – One-Page Platform Master Guide Added

### Added
- `00-overview/AYSTEK_CORE_ONE_PAGE_PLATFORM_GUIDE.md`
  - Tek sayfada platform özeti:
    - Go microservice katmanları ve modül yetenekleri
    - Next admin panel ekran kapsamı
    - Integrations (event-driven, chain, scheduled pull) operasyon modeli
    - Uçtan uca kullanım senaryoları ve başlangıç okuma sırası

### Updated
- `README.md`
  - “Recommended First Read” bölümüne one-page master guide referansı eklendi.

### Notes
- Bu sürüm, kullanıcıların sistemi baştan sona hızlı kavraması için yüksek seviyeli rehber ekleme güncellemesidir.

## [3.0.5] – Integrations User Usage Guide Expansion

### Updated
- `06-modules/INTEGRATIONS_PLATFORM_GUIDE.md`
  - Kullanıcı odaklı adım adım işletim rehberi eklendi:
    - event-driven kurulum
    - action test
    - chain test-run
    - scheduled pull oluşturma/test/izleme/retry
    - ekran bazlı troubleshooting yönlendirmesi
- `06-modules/modules/integrations.md`
  - Quick Start bölümü eklendi.
- `04-api/INTEGRATION_GUIDE.md`
  - Scheduled pull endpoint seti ve API kullanım akışı eklendi.

### Notes
- Bu sürüm, teknik özellikleri kullanıcı işletim adımlarıyla eşleştiren dokümantasyon genişletmesidir.

## [3.0.4] – Scheduled Pull UI + Runbook Documentation Sync

### Updated
- `README.md`
  - Operations bölümüne scheduled pull runbook referansı eklendi.
- `05-operations/OPERATIONS_GUIDE.md`
  - Scheduled pull smoke/QA adımları (`3.8`) eklendi.
- `06-modules/INTEGRATIONS_PLATFORM_GUIDE.md`
  - Scheduled pull API yüzeyi ve UI kapsamı güncellendi.
- `06-modules/modules/integrations.md`
  - Son durum bölümüne scheduled pull phase-4 tamamlanma notları eklendi.

### Added
- `05-operations/INTEGRATIONS_SCHEDULED_PULL_RUNBOOK_QA_MATRIX.md`
  - Scheduler/worker/pull-run/retry-DLQ/cursor/ingestion observability doğrulama matrisi.

### Notes
- Bu sürüm belge + UI kapsamı senkronizasyon güncellemesidir.

## [3.0.3] – Integrations Chain + Operations Documentation Sync

### Updated
- `README.md`
  - Operations bölümüne chain runbook/QA matrix ve GitLab protected setup referansları eklendi.
- `01-architecture/MODULE_INDEX.md`
  - `integrations` modülü için eski `documentation needed` notu kaldırıldı; mevcut modül dokümanına bağlandı.
- `05-operations/OPERATIONS_GUIDE.md`
  - Chain smoke/doğrulama ve GitLab protected environment güvenlik kontrol adımları eklendi.
- `06-modules/modules/integrations.md`
  - Son durum bölümü chain target modeli, chain test-run, chain graph editor usability hardening ve operasyon script referanslarıyla güncellendi.

### Added
- GitLab API tabanlı güvenlik doğrulama script referansı dokümantasyona işlendi:
  - `go/aystek-core/tools/ops/gitlab_verify_integrations_smoke_security.sh`

### Notes
- Bu sürüm belge güncellemesidir; kod tarafında chain runtime ve smoke pipeline özellikleriyle uyumludur.

## [3.0.2] – Integrations End-to-End Platform Guide Added

### Added
- `06-modules/INTEGRATIONS_PLATFORM_GUIDE.md` created as full-scope integrations handbook:
  - connectors/actions/hooks/executions/DLQ lifecycle
  - payload builder flow (mapping + enrichment + defaults)
  - response schema runtime behavior (JSON/XML/SOAP/Text)
  - dynamic event catalog behavior
  - admin UI flow and API surface summary
  - RabbitMQ operational validation references

### Updated
- `06-modules/modules/integrations.md` linked to full end-to-end guide.
- `README.md` modules index includes new integrations platform guide entry.

### Notes
- This update is documentation-only and reflects the latest integrations runtime/admin capabilities as of 2026-03-06.

## [3.0.1] – Integrations Runtime Documentation Update

### Updated
- `04-api/INTEGRATION_GUIDE.md` aligned with current integrations runtime behavior:
  - action test flow
  - dynamic response normalization
  - response schema based success evaluation
- `06-modules/modules/integrations.md` updated with latest production state:
  - response schema active decision model
  - format-agnostic response capture (REST/SOAP/XML/Text)
  - connector SSRF allowlist policy model
  - SOAPAction operational guidance

### Notes
- Documentation now reflects 2026-03-05 integration pipeline behavior and admin UI usage.

## [3.0.0] – Complete System Inventory & Enterprise Modernization

### Added
- **Complete Module Index Update**: Comprehensive 80-module system inventory with categorization
- **Module Classification System**: P0-P3 priority levels and maturity model (5 levels)
- **Module Dependency Matrix**: Critical path dependencies and integration map
- **Integration Platform Documentation**: Connectors, webhooks, DLQ, circuit breakers, idempotency
- **Stock Ecosystem Expansion**: Vehicle fleet, GPS tracking, subscription rentals, SMS integration
- **28-Language Translation Standard**: Enterprise-level multi-language support documentation
- **API Endpoint Summary**: 704 endpoints across 8 categories with authentication patterns
- **Module Health Dashboard**: Code coverage, uptime, security audit metrics

### Integration Platform
- **Integrations Module**: Complete integration platform with connectors, actions, hooks, executions
- **Dead Letter Queue (DLQ)**: Modern UI with gradient design, statistics, bulk operations
- **Circuit Breakers**: Service protection and resilience patterns
- **Idempotency Store**: Redis-based duplicate prevention
- **Webhook Management**: Event-driven integration patterns

### UI/UX Modernization
- **Modern Admin UI Patterns**: Gradient headers (purple-fuchsia), backdrop blur, animations
- **Component Standardization**: IntegrationsAdminHeader, AdminDataTable, TableToolbar
- **Statistics Cards**: Color-coded with gradients (emerald, yellow, blue badges)
- **Status Badges**: Inline color coding with Lucide React icons
- **Modal Animations**: Fade-in, zoom-in-95, duration-200/300

### Architecture
- **Module Maturity Model**: 5-level maturity assessment with criteria
- **Documentation Coverage**: 64/80 modules documented (16 require documentation)
- **Dependency Analysis**: Core and integration dependency tables
- **Performance SLA**: <200ms p95 (current: 178ms), 99.92% uptime

### Security & Compliance
- **Multi-layer Authentication**: JWT + Permissions + API Keys
- **Security Audit**: Pass status maintained
- **Code Coverage**: 84% (target: >80%)

### Operations
- **Menu System**: Added 6 integrations menu entries with proper hierarchy
- **Translation Keys**: Added integrations translations with explicit 28-language definitions
- **Build Verification**: Go backend and Next.js frontend compilation verified

### Documentation Scope Updated
- `01-architecture/MODULE_INDEX.md`: Complete rewrite with 80 modules
- `user-guide/README.md`: Version bump to 3.0.0, updated metadata
- `user-guide/CHANGELOG.md`: Added 3.0.0 release notes

## [2.1.0] – Enterprise Content Ops & Ingestion Platform

### Added
- Connector registry documentation for ingestion surface
- Moderation workflow documentation and state transition model
- Site policy engine documentation
- Feature flag engine documentation
- Injection code registry documentation
- Ad policy engine documentation
- Multi-tenant and multi-site enforcement details in API/security/architecture sections
- New architecture document: `01-architecture/CONTENT_OPERATIONS_ARCHITECTURE.md`
- New module document: `06-modules/CONTENT_INGESTION_AND_OPS_MODULES.md`

### Security
- Sanitized injection code storage model documented
- Encrypted connector secrets and rotation model documented

### Architecture
- New bounded contexts documented:
  - `content_ingestion`
  - `content_ops_settings`

### Documentation Scope Updated
- `00-overview`
- `01-architecture`
- `02-security`
- `03-compliance`
- `04-api`
- `05-operations`
- `06-modules`
- `07-executive`
- `README.md`
- `GOVERNANCE.md`

## [2.0.0] – User Guide Enterprise Restructure
- User guide layout restructured into numbered enterprise sections (`00` to `07`).
- Enterprise documents integrated and renamed into canonical publication names.
- Existing guides relocated to architecture, security, API, operations, and modules sections.
- Historical/internal materials preserved under `archive/` without deletion.
- Cross-document references normalized to new structure.
- Root governance document introduced: `GOVERNANCE.md`.
