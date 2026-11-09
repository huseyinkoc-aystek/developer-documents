---
Title: Content Ingestion and Content Ops Modules
Version: 2.1.0
Status: Enterprise Certified
Last Updated: 2026-02-18
Related Documents:
  - ../01-architecture/CONTENT_OPERATIONS_ARCHITECTURE.md
  - ../04-api/API_REFERENCE.md
  - ../05-operations/OPERATIONS_GUIDE.md
---

# Table of Contents
1. [Module: Content Ingestion](#1-module-content-ingestion)
2. [Module: Moderation Engine](#2-module-moderation-engine)
3. [Module: Content Ops Settings](#3-module-content-ops-settings)

## 1. Module: Content Ingestion
### 1.1 Purpose
Ingest external content streams through controlled connectors and convert source payloads into moderation-ready items.

### 1.2 Key Features
- Multi-source connector support (RSS/API/Agency/Scraper/Webhook patterns)
- Schedule-based and manual run execution
- Retry/backoff and rate-limit-aware runtime controls
- Dedupe/fingerprint model for duplicate suppression
- Job and item observability surfaces
### 1.3 Architecture
- Connector registry
- Scheduler and worker runtime
- Ingestion jobs and items
- Deduplication/fingerprints
- Outbox events for downstream propagation

### 1.4 Configurable Parameters
- Connector schedule (cron/interval)
- Retry and backoff
- Rate limits
- Mapping and normalization rules
- Dedupe strategy toggles

### 1.5 Scaling Notes
- Horizontal worker scaling by tenant/site partitions
- Queue-based execution and idempotency keys
- Read model rollups for high-volume dashboards

### 1.6 Limitations
- External source quality varies and requires moderation controls.
- Auto-publish behavior must remain policy-gated, not default.

## 2. Module: Moderation Engine
### 2.1 Workflow
`pending_review -> approved/rejected -> scheduled/published`

### 2.2 Key Features
- Moderation-first publishing guardrail
- Approve/reject/edit/schedule action model
- Reviewer responsibility traceability via audit logs
- Queue filtering and bulk-operation compatibility
### 2.3 Operational SLA Considerations
- Queue aging thresholds
- Reviewer assignment policy
- Priority lanes for critical sources/categories

### 2.4 Multi-Step Approval Option
Optional staged approval can be enabled for high-risk categories:
- Reviewer approval
- Senior editor approval
- Publish authorization

## 3. Module: Content Ops Settings
### 3.1 Policy vs Feature Flags
- Policy: operational rules that define expected editorial behavior.
- Feature Flags: runtime capability toggles and controlled rollout levers.

### 3.2 Key Features
- Site-scoped policy control plane
- Runtime feature-flag governance
- Quota and limit enforcement surface
- Ad behavior control (refresh/popup/device)
- Secure and versioned injection code registry
### 3.3 Limits and Runtime Impact
Limits directly influence:
- Ingestion throughput
- Moderation backlog ceilings
- Injection code payload acceptance
- Connector count boundaries

### 3.4 Ad Policy Integration
Ad policy controls are site-scoped and runtime-bound, covering refresh/popup/device behavior.

### 3.5 Injection Code Governance
- Versioned, auditable lifecycle
- Sanitization and size checks
- Confirm-before-save UX requirement
- Soft delete and rollback-safe operation
