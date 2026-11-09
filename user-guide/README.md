---
Title: AYSTEK Core User Guide Index
Version: 3.0.8
Status: Enterprise Certified - Production Ready
Last Updated: 2026-06-05
Classification: Internal - Enterprise Documentation
Related Documents:
  - ./CHANGELOG.md
  - ./GOVERNANCE.md
---

# Table of Contents
1. [Purpose](#1-purpose)
2. [Architecture](#2-architecture)
3. [Security](#3-security)
4. [Compliance](#4-compliance)
5. [API](#5-api)
6. [Operations](#6-operations)
7. [Modules](#7-modules)
8. [Executive](#8-executive)
9. [Archive](#9-archive)

## 1. Purpose
This index defines the enterprise documentation layout for AYSTEK Core under `user-guide/`.

### Recommended First Read
- `00-overview/AYSTEK_CORE_ONE_PAGE_PLATFORM_GUIDE.md`: Go microservice + Next admin panel platform capabilities in a single page (business + technical overview).

### Key Features Covered in This Guide
- Multi-tenant and multi-site isolated architecture
- Event-driven content ingestion and moderation workflow
- Site-scoped operational governance (policy/flags/limits/ad/injection)
- Enterprise security and compliance evidence model
- Operational playbooks for high-volume digital publishing

## 2. Architecture
- `00-overview/OVERVIEW.md`: Platform scope, boundaries, and canonical ownership map.
- `01-architecture/TENANT_AND_SITE_GUIDE.md`: Tenant/site contract and enforcement matrix.
- `01-architecture/MODULE_INDEX.md`: Module documentation map.
- `01-architecture/CONTENT_OPERATIONS_ARCHITECTURE.md`: Content ingestion and site operations bounded contexts.
- `01-architecture/Enterprise_Isolation_Whitepaper.md`: Authoritative isolation whitepaper.
- `01-architecture/Distributed_Zero_Trust_Tenant_Architecture_Spec.md`: Distributed zero-trust architecture specification.

## 3. Security
- `02-security/AUTHENTICATION_AND_AUTHORIZATION.md`: AuthN/AuthZ + ingestion/settings security controls.
- `02-security/Enterprise_Security_Architecture_Booklet.md`: Security-layer model for enterprise review.
- `02-security/TENANT_SITE_SECURITY_MATRIX.md`: Deterministic tenant/site requirement matrix.

## 4. Compliance
- `03-compliance/Regulatory_Compliance_Technical_Annex.md`: SOC2/ISO/GDPR/KVKK mapping and audit considerations.

## 5. API
- `04-api/API_REFERENCE.md`: Runtime API contract summary for enterprise ingestion and content ops settings.
- `04-api/INTEGRATION_GUIDE.md`: Integrator and partner usage guidance.
- `04-api/ERROR_CODE_REFERENCE.md`: Canonical error contract and codes.

## 6. Operations
- `05-operations/DEPLOYMENT_SCENARIOS.md`: Deployment topology patterns.
- `05-operations/OPERATIONS_GUIDE.md`: Operational lifecycle, monitoring, and playbooks.
- `05-operations/INTEGRATIONS_EVENTBRIDGE_SMOKE_RUNBOOK.md`: Post-deploy smoke verification for module-agnostic integration event pipeline.
- `05-operations/INTEGRATIONS_CHAIN_RUNBOOK_QA_MATRIX.md`: Chain runtime smoke steps and QA matrix (branch/retry/DLQ/idempotency).
- `05-operations/INTEGRATIONS_SCHEDULED_PULL_RUNBOOK_QA_MATRIX.md`: Scheduled pull smoke/test matrix (scheduler/worker/retry/DLQ/cursor/ingestion observability).
- `05-operations/GITLAB_PROTECTED_ENV_SETUP.md`: Protected environment + protected/masked variable setup and API verification script usage.

## 7. Modules
- `06-modules/modules/aystek_haber.md`: AYSTEK Haber public-site veri sahipliği, API, seed, yönetim ve operasyon sözleşmesi.
- `06-modules/CONTENT_INGESTION_AND_OPS_MODULES.md`: Content ingestion, moderation, and content ops module guide.
- `06-modules/MARKETPLACE_GUIDE.md`: Marketplace metadata and consumption.
- `06-modules/INTEGRATIONS_PLATFORM_GUIDE.md`: End-to-end integrations platform guide (connectors/actions/hooks/executions/DLQ, mapping, response schema, security, operations).
- `06-modules/FLEET_TALLY_USER_GUIDE_TR.md`: GPS çetele kontrol / Fleet Tally kullanıcı ve operasyon kılavuzu.
- `06-modules/PLUGIN_USER_GUIDE.md`: Plugin consumption and governance constraints.
- `06-modules/modules/`: Module-by-module reference files.

## 8. Executive
- `07-executive/Board_Level_Risk_Isolation_Report.md`: Board-level risk narrative and strategic positioning.

## 9. Archive
- `archive/`: Historical and superseded documents preserved for traceability.
- `archive/internal/`: Internal reference snapshots retained without deletion.
