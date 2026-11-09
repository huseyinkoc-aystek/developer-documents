---
Title: Content Operations Architecture
Version: 2.1.0
Status: Enterprise Certified
Last Updated: 2026-02-18
Related Documents:
  - ./TENANT_AND_SITE_GUIDE.md
  - ../04-api/API_REFERENCE.md
  - ../06-modules/CONTENT_INGESTION_AND_OPS_MODULES.md
---

# Table of Contents
1. [Bounded Contexts](#1-bounded-contexts)
2. [Architecture Features](#2-architecture-features)
3. [Content Ingestion Context](#3-content-ingestion-context)
4. [Content Ops Settings Context](#4-content-ops-settings-context)
5. [Collections and Persistence](#5-collections-and-persistence)
6. [Event Flow and Outbox](#6-event-flow-and-outbox)
7. [Tenant and Site Scoping](#7-tenant-and-site-scoping)
8. [Moderation State Machine](#8-moderation-state-machine)

## 1. Bounded Contexts
### A) `content_ingestion`
- Connector Registry
- Job Scheduler
- Worker runtime
- Ingestion items
- Moderation queue
- Fingerprints/deduplication
- Idempotency controls

### B) `content_ops_settings`
- Policy engine
- Feature flags engine
- Limits/quota model
- Ad policy model
- Injection code registry

## 2. Architecture Features
- DDD-aligned bounded contexts separating ingestion runtime from operational governance.
- Connector model decoupled from moderation and publishing decisions.
- CQRS-ready read patterns via list endpoints, pagination, and rollup-compatible structures.
- Idempotency and dedupe primitives included in ingestion design.
- Outbox event contracts for post-commit consistency and async integration.

## 3. Content Ingestion Context
### 2.1 Core Components
- Connector Registry: configurable connectors (RSS/API/Agency/Scraper/etc.)
- Job Scheduler: schedule-based execution
- Worker Model: fetch/normalize/dedupe/persist
- Moderation Queue: approve/reject/edit/schedule operations

### 2.2 Reliability Model
- Idempotency key per job/item workflow
- Fingerprint dedupe via canonical URL hash + external ID + dedupe key
- Retry/backoff and health scoring per connector

## 4. Content Ops Settings Context
### 3.1 Policy Engine
Defines moderation and operational behavior per site.

### 3.2 Feature Flags Engine
Controls optional runtime capabilities and emergency toggles.

### 3.3 Limits & Quotas
Defines ingestion, moderation backlog, connector count, and injection size caps.

### 3.4 Ad Policy
Defines ad refresh/popup/device-level policy controls.

### 3.5 Injection Code Registry
Stores controlled, versioned, auditable code snippets by scope/environment/device.

## 5. Collections and Persistence
### 4.1 `content_ingestion`
- `bot_connectors`
- `ingestion_jobs`
- `ingested_items`
- `editorial_drafts`
- `content_bots_outbox`

### 4.2 `content_ops_settings`
- `content_ops_site_policies`
- `content_ops_feature_flags`
- `content_ops_content_limits`
- `content_ops_ad_policies`
- `content_ops_injection_codes`
- `content_ops_settings_outbox`

## 6. Event Flow and Outbox
### 5.1 Ingestion Events
- `bot_connector.created|updated|disabled|deleted`
- `ingestion_job.started|finished|failed`
- `ingestion_item.created|normalized|queued`
- `moderation_item.approved|rejected|edited|published`

### 5.2 Settings Events
- `site_policy.updated`
- `feature_flags.updated`
- `content_limits.updated`
- `ad_policy.updated`
- `injection_code.created|updated|deleted`

### 5.3 Outbox Guarantees
Each write-side mutation produces deterministic outbox records scoped by tenant and site for replay-safe downstream consumption.

## 7. Tenant and Site Scoping
- Every read/write path is tenant-scoped.
- Site-scoped endpoints require `siteId` route/context.
- Tenant/site are derived from auth/context middleware; never trusted from request payload.

## 8. Moderation State Machine
| State | Next States |
| --- | --- |
| pending_review | approved, rejected |
| approved | scheduled, published |
| scheduled | published |
| rejected | — |

Operational rule: bot-ingested content is moderation-first by default; direct publish requires explicit policy permission.
