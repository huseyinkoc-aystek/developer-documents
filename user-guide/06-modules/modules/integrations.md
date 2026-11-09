---
Title: Integrations Module - Enterprise Integration Platform
Version: 1.4.0
Status: GA (General Availability)
Last Updated: 2026-03-08
Module Path: /internal/modules/integrations
Priority: P1 - Essential
Dependencies: RabbitMQ, Redis, MongoDB
Related Docs:
  - ../INTEGRATIONS_PLATFORM_GUIDE.md
  - erp_integrations.md
  - notifications.md
  - ../CONTENT_INGESTION_AND_OPS_MODULES.md
---

# Integrations Module

## Full Guide

Bu dokümanın güncel ve uçtan uca versiyonu için:

`user-guide/06-modules/INTEGRATIONS_PLATFORM_GUIDE.md`

## Executive Summary

The **Integrations Module** is an enterprise-grade integration platform that provides comprehensive capabilities for connecting AYSTEK Core with external systems, services, and APIs. It implements modern integration patterns including **connectors**, **webhooks**, **circuit breakers**, **idempotency**, and **dead letter queue** management.

**Key Capabilities:**
- Pre-built connectors for common services (ERP, CRM, payment gateways, shipping)
- Custom webhook creation and management
- Action execution engine with scheduling
- Dead Letter Queue (DLQ) for failed message handling
- Circuit breakers for service protection
- Idempotency store for duplicate prevention
- Hook system for event-driven integrations
- Execution tracking and monitoring

**Production Status:** ✅ General Availability  
**Endpoints:** 28 REST APIs  
**Database Collections:** 6 (connectors, actions, hooks, executions, dlq, circuit_breaker)

---

## Latest State (2026-03-08)

### 0) Hook Target Model (Backward Compatible)
- Hook artık `target_type=action|chain` ile çalışır.
- `target_type=action` mevcut davranışı korur.
- `target_type=chain` olduğunda `chain_id` üzerinden çok adımlı chain runner devreye girer.
- Chain içindeki her `action` node kendi `action_id -> connector_id` bağıyla çalışır.

### 1) Module-Agnostic Hook Binding (Dynamic Event Key)
- Hook definition is now operated as **module-independent** event subscription.
- Admin UI (`/admin/integrations/hooks/:id`) uses `Event Key` field to dynamically bind to any module event.
- Integration runtime does not hardcode stock-specific behavior; any published domain event key can be targeted.
- Recommended operational validation: `user-guide/05-operations/INTEGRATIONS_EVENTBRIDGE_SMOKE_RUNBOOK.md`.

### 2) Chain Runtime + Test Run
- Chain node types: `start`, `action`, `transform`, `branch`, `end`.
- Chain validate/test-run endpointleri aktiftir (`/admin/integrations/chains/:id/validate`, `/admin/integrations/chains/:id/test-run`).
- Test-run çıktısı node bazlı step sonuçları ve final payload preview döner.
- Step-level execution timeline ile observability genişletildi.

### 3) Graph Editor Usability Hardening
- Drag-drop canvas + pan/zoom + auto-focus config panel.
- `Undo/Redo` (Ctrl/Cmd+Z, Ctrl/Cmd+Shift+Z / Ctrl/Cmd+Y).
- Klavye ile node taşıma (Arrow, Shift+Arrow), `Delete` ile node silme.
- Snap-to-grid, komşu node hizalama yardımcıları, auto-layout.

### 4) Dynamic Response Evaluation (Primary on "Yanıt Şeması")
- Action form **Response Schema** is now runtime-active for success/failure decision.
- Configured keys are matched against normalized response (REST/SOAP/XML/Text).
- Example active schema:
  - `gpsUpdateResult: "true"`
  - `meta.status: "ok"`
  - `result.success: true`
- Mismatch outcome:
  - `RESPONSE_SCHEMA_MISMATCH`
  - stored in execution response with `_response_schema_check`

### 5) Response Normalization (Format-Agnostic)
- JSON:
  - parsed object/array
  - `json.flat` leaf path map
  - promoted leaf values
- XML/SOAP:
  - `xml.flat` + `xml.values`
  - SOAP fault auto-detection (`faultstring`)
- Text:
  - `raw`
- `_http` context available for evaluations:
  - `status_code`, `content_type`, normalized headers

### 6) Action Test Runtime Traceability
- Test execution persists:
  - request payload
  - response payload
  - request preview (url, headers, rendered body, warnings)
- Dry-run and real-run both generate trace-ready outputs.

### 7) Connector Outbound Security (Mongo + UI Managed)
- Connector-scoped SSRF policy:
  - `outbound_allowlist[]` (domain/ip/cidr/port/scheme)
  - `allow_private_ips_for_allowlist_match`
- Runtime-effective allowlist includes:
  - manual rules
  - connector `base_url` host
  - connector `test_url` host
- Private IP access remains blocked unless explicit allowlist match + flag enabled.

### 8) Operations Automation Note (GitLab)
- Production smoke için manual/protected job kuralı pipeline'a bağlıdır.
- `production` protected environment ve `INTEGRATIONS_SMOKE_TOKEN` protected+masked ayarı zorunludur.
- API tabanlı doğrulama scripti:
  - `go/aystek-core/tools/ops/gitlab_verify_integrations_smoke_security.sh`

### 9) SOAP Operational Note
- For SAP GPS endpoint, `SOAPAction` must be full URI:
  - `http://sapportal.gurseltur.com.tr/gpsUpdate`
- Short action names (e.g. `gpsUpdate`) may produce server 500 SOAP fault.

### 10) Scheduled Pull Capability (Phase-4 Complete)
- Admin UI now includes:
  - `Integrations > Scheduled Pulls` list/create/edit/test-run
  - `Integrations > Pull Runs` list/detail with ingestion summary/issues visibility
- Pull run detail includes:
  - ingestion summary cards (raw/canonical/upsert/issues/events)
  - structured ingestion issues table
  - warnings + request/response metadata panels
- Operational validation reference:
  - `user-guide/05-operations/INTEGRATIONS_SCHEDULED_PULL_RUNBOOK_QA_MATRIX.md`

### 12) Shared Service / Grant Policy Layer (Phase-4 Complete)
- Shared service ve service grant yönetimi admin UI'da aktiftir:
  - `/admin/integrations/shared-services`
  - `/admin/integrations/shared-service-grants`
- Policy alanları backend enforcement hattına bağlıdır:
  - action allowlist
  - scopes (`sms.send`, `sms.*`, `*`)
  - method/path guard
  - limits (rpm/daily quota vb.)
  - secret bindings + payload guard + mTLS policy
- Execution hattında policy ihlalleri `shared_service_policy_denied` ve ilgili validation kodlarıyla izlenir.

### 13) Shared Service Observability + UX Hardening (Phase-5)
- Service ve grant edit ekranlarında usage summary kartları eklendi (7/14/30 gün).
- UI warning paneli; quota baskısı, policy deny trendi ve hardening eksiklerini görünür kılar.
- Yeni API endpointleri:
  - `GET /admin/integrations/shared-services/:id/usage-summary`
  - `GET /admin/integrations/shared-service-grants/:id/usage-summary`
- Operasyon runbook:
  - `user-guide/05-operations/INTEGRATIONS_SHARED_SERVICE_RUNBOOK_QA_MATRIX.md`
  - `user-guide/05-operations/INTEGRATIONS_SHARED_SERVICE_ALERT_THRESHOLDS.md`
  - `user-guide/06-modules/INTEGRATIONS_SHARED_SERVICE_UI_WALKTHROUGH.md`

### 11) Quick Start (User)
1. Connector oluştur.
2. Action tanımla ve Action Test ile doğrula.
3. Hook/Chain hedefini belirle ve test et.
4. Canlı çalıştırmaları Executions/DLQ ekranından izle.
5. Periyodik veri çekimi gerekiyorsa Scheduled Pull oluştur, Pull Runs üzerinden sonuçları takip et.
6. Pull run başarısız olduysa run detail ekranından manuel retry kullan.
7. Shared-service kullanan akışlarda policy deny/quota durumunu service/grant usage panellerinden düzenli takip et.

---

## Table of Contents

1. [Architecture Overview](#1-architecture-overview)
2. [Core Concepts](#2-core-concepts)
3. [Connectors](#3-connectors)
4. [Actions](#4-actions)
5. [Hooks](#5-hooks)
6. [Executions](#6-executions)
7. [Dead Letter Queue (DLQ)](#7-dead-letter-queue-dlq)
8. [Circuit Breakers](#8-circuit-breakers)
9. [Idempotency](#9-idempotency)
10. [API Reference](#10-api-reference)
11. [Data Models](#11-data-models)
12. [Security](#12-security)
13. [Operations](#13-operations)
14. [Performance](#14-performance)
15. [Troubleshooting](#15-troubleshooting)

---

## 1. Architecture Overview

### System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     AYSTEK Core Platform                     │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌──────────────┐   ┌──────────────┐   ┌──────────────┐    │
│  │  Connectors  │   │   Actions    │   │    Hooks     │    │
│  │   Registry   │   │   Engine     │   │   System     │    │
│  └──────┬───────┘   └──────┬───────┘   └──────┬───────┘    │
│         │                   │                   │             │
│         └───────────────────┴───────────────────┘             │
│                             │                                 │
│                   ┌─────────▼─────────┐                      │
│                   │  Execution Engine  │                      │
│                   └─────────┬─────────┘                      │
│                             │                                 │
│         ┌───────────────────┼───────────────────┐            │
│         │                   │                   │            │
│  ┌──────▼──────┐   ┌────────▼────────┐   ┌─────▼──────┐    │
│  │ Idempotency │   │Circuit Breakers │   │    DLQ     │    │
│  │    Store    │   │     Store       │   │   Handler  │    │
│  │   (Redis)   │   │  (In-Memory)    │   │  (MongoDB) │    │
│  └─────────────┘   └─────────────────┘   └────────────┘    │
│                                                               │
└─────────────────────────────────────────────────────────────┘
                             │
                    ┌────────┴────────┐
                    │                 │
            ┌───────▼──────┐   ┌──────▼───────┐
            │   External   │   │   External   │
            │   Services   │   │    APIs      │
            └──────────────┘   └──────────────┘
```

### Component Responsibilities

| Component | Purpose | Technology |
|-----------|---------|------------|
| **Connectors** | Pre-configured integrations | MongoDB |
| **Actions** | Custom integration workflows | MongoDB + Scheduler |
| **Hooks** | Event-driven triggers | MongoDB + RabbitMQ |
| **Executions** | Execution tracking and history | MongoDB |
| **DLQ** | Failed message management | MongoDB |
| **Circuit Breakers** | Service protection | In-Memory Store |
| **Idempotency** | Duplicate prevention | Redis |

---

## 2. Core Concepts

### Integration Patterns

The module supports multiple integration patterns:

1. **Request-Response (Synchronous)**
   - Direct API calls with immediate response
   - Used for: Real-time data queries, validation

2. **Fire-and-Forget (Asynchronous)**
   - Message queuing without waiting for response
   - Used for: Notifications, logging, analytics

3. **Webhook (Event-Driven)**
   - External systems push data to AYSTEK
   - Used for: Payment confirmations, shipping updates, CRM events

4. **Scheduled Execution**
   - Time-based or cron-triggered actions
   - Used for: Batch synchronization, reporting, cleanup

5. **Event-Driven (Pub/Sub)**
   - Internal events trigger external integrations
   - Used for: Order created → ERP sync, User registered → CRM update

---

## 3. Connectors

### Overview

**Connectors** are pre-built, reusable integration templates for common services and platforms.

### Connector Types

| Type | Examples | Use Cases |
|------|----------|-----------|
| **ERP** | SAP, Oracle, NetSuite | Inventory, orders, customers |
| **CRM** | Salesforce, HubSpot | Lead capture, customer sync |
| **Payment** | Stripe, PayPal, Adyen | Payment processing, refunds |
| **Shipping** | FedEx, UPS, DHL | Shipping rates, tracking |
| **Email** | SendGrid, Mailgun | Transactional emails |
| **SMS** | Twilio, Vonage | OTP, notifications |
| **Storage** | AWS S3, Azure Blob | File uploads, backups |
| **Analytics** | Google Analytics, Mixpanel | Event tracking |

### Connector Configuration

Each connector requires:
- **Authentication credentials** (API key, OAuth tokens)
- **Endpoint URLs** (base URL, API version)
- **Rate limits** (requests per second/minute)
- **Timeout settings** (connection, request, total)
- **Retry policy** (max retries, backoff strategy)
- **Circuit breaker config** (failure threshold, timeout)

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/admin/integrations/connectors` | Create connector |
| `GET` | `/api/admin/integrations/connectors` | List connectors |
| `GET` | `/api/admin/integrations/connectors/:id` | Get connector details |
| `PUT` | `/api/admin/integrations/connectors/:id` | Update connector |
| `DELETE` | `/api/admin/integrations/connectors/:id` | Delete connector |
| `POST` | `/api/admin/integrations/connectors/:id/test` | Test connection |

---

## 4. Actions

### Overview

**Actions** are executable integration workflows that perform specific tasks (e.g., sync inventory, send notification, create invoice).

### Action Types

1. **Data Sync Actions**
   - Pull data from external systems
   - Push data to external systems
   - Two-way synchronization

2. **Notification Actions**
   - Send email, SMS, push notifications
   - Multi-channel delivery

3. **Processing Actions**
   - Data transformation
   - Validation and enrichment
   - Aggregation and reporting

4. **Workflow Actions**
   - Multi-step orchestration
   - Conditional logic
   - Error handling and retry

### Action Execution

Actions can be executed:
- **Manually:** Admin triggers execution
- **Scheduled:** Cron-based scheduling (e.g., every hour, daily at 2 AM)
- **Event-Driven:** Triggered by internal events (order created, user registered)
- **Hook-Based:** Triggered by webhooks from external systems

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/admin/integrations/actions` | Create action |
| `GET` | `/api/admin/integrations/actions` | List actions |
| `GET` | `/api/admin/integrations/actions/:id` | Get action details |
| `PUT` | `/api/admin/integrations/actions/:id` | Update action |
| `DELETE` | `/api/admin/integrations/actions/:id` | Delete action |
| `POST` | `/api/admin/integrations/actions/:id/execute` | Execute action |
| `POST` | `/api/admin/integrations/actions/:id/schedule` | Schedule action |

---

## 5. Hooks

### Overview

**Hooks** are webhook endpoints that allow external systems to push data into AYSTEK Core.

### Use Cases

- **Payment Gateway Webhooks:** Payment success/failure notifications
- **Shipping Carrier Webhooks:** Tracking updates, delivery confirmation
- **CRM Webhooks:** Lead updates, contact changes
- **Marketplace Webhooks:** Order notifications from external marketplaces

### Hook Security

All hooks support:
- **Signature verification:** HMAC-SHA256 signature validation
- **IP whitelisting:** Restrict requests to known IPs
- **API key authentication:** Require API key in headers
- **Rate limiting:** Prevent abuse and DoS attacks

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/admin/integrations/hooks` | Create hook |
| `GET` | `/api/admin/integrations/hooks` | List hooks |
| `GET` | `/api/admin/integrations/hooks/:id` | Get hook details |
| `PUT` | `/api/admin/integrations/hooks/:id` | Update hook |
| `DELETE` | `/api/admin/integrations/hooks/:id` | Delete hook |
| `POST` | `/api/webhooks/:hookId/:secret` | Public webhook endpoint |

---

## 6. Executions

### Overview

The **Executions** system tracks all integration execution attempts, including:
- Execution ID and timestamp
- Action or hook that was executed
- Input parameters and payload
- Response data and status code
- Execution duration
- Success/failure status
- Error messages (if any)

### Execution Lifecycle

```
┌─────────┐
│ Pending │  → Initial state when execution is queued
└────┬────┘
     │
┌────▼────┐
│ Running │  → Execution in progress
└────┬────┘
     │
     ├────────┐
     │        │
┌────▼────┐  │
│ Success │  │  → Completed successfully
└─────────┘  │
             │
        ┌────▼────┐
        │ Failed  │  → Execution failed (stored in DLQ)
        └─────────┘
```

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/admin/integrations/executions` | List executions |
| `GET` | `/api/admin/integrations/executions/:id` | Get execution details |
| `DELETE` | `/api/admin/integrations/executions/:id` | Delete execution record |

---

## 7. Dead Letter Queue (DLQ)

### Overview

The **Dead Letter Queue (DLQ)** is a system for managing failed integration messages that couldn't be processed successfully.

### Features

- **Automatic Capture:** Failed messages are automatically moved to DLQ
- **Modern UI:** Gradient-based admin interface with statistics
- **Bulk Operations:** Requeue multiple messages at once
- **Filtering:** Filter by status, connector, date range
- **Retry Management:** Manual or automatic retry with exponential backoff
- **Error Analysis:** Detailed error messages and stack traces

### DLQ Workflow

```
┌──────────────┐
│ Integration  │
│  Execution   │
└──────┬───────┘
       │ Fails
       ▼
┌──────────────┐
│  DLQ Entry   │  ← Status: pending
└──────┬───────┘
       │ Admin Reviews
       ▼
┌──────────────┐
│   Requeue    │  ← Retry execution
└──────┬───────┘
       │
       ├───────────────┐
       │               │
  ┌────▼─────┐   ┌─────▼──────┐
  │ Success  │   │ Failed     │
  │ (Remove) │   │ (Keep DLQ) │
  └──────────┘   └────────────┘
```

### Modern UI Features

- **Gradient Header:** Purple-fuchsia gradient with breadcrumb navigation
- **Statistics Cards:** Total, Pending, Requeued counts with color-coded badges
- **AdminDataTable:** Sortable, filterable data table with row selection
- **Modal Animations:** Fade-in, zoom-in animations for detail modal
- **Action Buttons:** View details, requeue individual, bulk requeue
- **Status Badges:** Color-coded (yellow for pending, emerald for requeued)

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/admin/integrations/dlq` | List DLQ entries |
| `GET` | `/api/admin/integrations/dlq/:id` | Get DLQ entry details |
| `POST` | `/api/admin/integrations/dlq/:id/requeue` | Requeue single entry |
| `POST` | `/api/admin/integrations/dlq/bulk-requeue` | Requeue multiple entries |
| `DELETE` | `/api/admin/integrations/dlq/:id` | Delete DLQ entry |

---

## 8. Circuit Breakers

### Overview

**Circuit Breakers** protect external services from being overwhelmed during failures by temporarily stopping requests when failure rate exceeds threshold.

### Circuit States

| State | Behavior | Transition |
|-------|----------|------------|
| **Closed** | Normal operation, all requests allowed | → Open (if failure threshold reached) |
| **Open** | All requests blocked, return error immediately | → Half-Open (after timeout) |
| **Half-Open** | Limited test requests allowed | → Closed (if success) or → Open (if fail) |

### Configuration

```go
type CircuitBreakerConfig struct {
    FailureThreshold  int           // 5 failures
    SuccessThreshold  int           // 2 successes to close
    Timeout           time.Duration // 30 seconds
    MaxRequests       uint32        // 1 request in half-open
}
```

### Benefits

- **Service Protection:** Prevents cascading failures
- **Fast Failure:** Immediate error response when circuit is open
- **Automatic Recovery:** Self-healing when service recovers
- **Resource Efficiency:** Reduces load on failing services

---

## 9. Idempotency

### Overview

The **Idempotency Store** prevents duplicate processing of the same request, ensuring that retries don't cause side effects (e.g., double charges, duplicate orders).

### How It Works

1. Client sends request with `Idempotency-Key` header
2. System checks Redis for existing key
3. If key exists → Return cached response
4. If key doesn't exist → Process request, cache response, return

### Redis Implementation

```go
type RedisIdempotencyStore struct {
    client *redis.Client
    ttl    time.Duration // 24 hours
}

// Store key format: "idempotency:{key}"
// Value: JSON response body
```

### Use Cases

- **Payment Processing:** Prevent double charges on retry
- **Order Creation:** Avoid duplicate orders
- **Webhook Processing:** Prevent reprocessing of same event
- **API Retries:** Safe retry without side effects

---

## 10. API Reference

### Complete Endpoint List

#### Connectors (5 endpoints)
- `POST /api/admin/integrations/connectors`
- `GET /api/admin/integrations/connectors`
- `GET /api/admin/integrations/connectors/:id`
- `PUT /api/admin/integrations/connectors/:id`
- `DELETE /api/admin/integrations/connectors/:id`
- `POST /api/admin/integrations/connectors/:id/test`

#### Actions (7 endpoints)
- `POST /api/admin/integrations/actions`
- `GET /api/admin/integrations/actions`
- `GET /api/admin/integrations/actions/:id`
- `PUT /api/admin/integrations/actions/:id`
- `DELETE /api/admin/integrations/actions/:id`
- `POST /api/admin/integrations/actions/:id/execute`
- `POST /api/admin/integrations/actions/:id/schedule`

#### Hooks (6 endpoints)
- `POST /api/admin/integrations/hooks`
- `GET /api/admin/integrations/hooks`
- `GET /api/admin/integrations/hooks/:id`
- `PUT /api/admin/integrations/hooks/:id`
- `DELETE /api/admin/integrations/hooks/:id`
- `POST /api/webhooks/:hookId/:secret`

#### Executions (3 endpoints)
- `GET /api/admin/integrations/executions`
- `GET /api/admin/integrations/executions/:id`
- `DELETE /api/admin/integrations/executions/:id`

#### DLQ (5 endpoints)
- `GET /api/admin/integrations/dlq`
- `GET /api/admin/integrations/dlq/:id`
- `POST /api/admin/integrations/dlq/:id/requeue`
- `POST /api/admin/integrations/dlq/bulk-requeue`
- `DELETE /api/admin/integrations/dlq/:id`

#### Circuit Breakers (2 endpoints)
- `GET /api/admin/integrations/circuit-breakers`
- `POST /api/admin/integrations/circuit-breakers/:id/reset`

**Total:** 28 API endpoints

---

## 11. Data Models

### Connector Model

```typescript
interface Connector {
  _id: string;
  tenant_id: string;
  site_id: string;
  name: string;
  type: "erp" | "crm" | "payment" | "shipping" | "email" | "sms" | "storage" | "analytics";
  config: {
    base_url: string;
    auth_type: "api_key" | "oauth" | "basic";
    credentials: Record<string, string>; // Encrypted
    headers?: Record<string, string>;
    timeout: number; // milliseconds
    rate_limit?: {
      requests: number;
      period: "second" | "minute" | "hour";
    };
  };
  status: "active" | "inactive" | "error";
  created_at: Date;
  updated_at: Date;
  created_by: string;
}
```

### Action Model

```typescript
interface Action {
  _id: string;
  tenant_id: string;
  site_id: string;
  connector_id: string;
  name: string;
  description: string;
  type: "sync" | "notification" | "webhook" | "scheduled";
  method: "GET" | "POST" | "PUT" | "DELETE" | "PATCH";
  endpoint: string;
  payload_template?: string; // Handlebars template
  schedule?: string; // Cron expression
  retry_config: {
    max_retries: number;
    backoff: "linear" | "exponential";
    delay: number; // milliseconds
  };
  status: "active" | "inactive";
  created_at: Date;
  updated_at: Date;
}
```

### DLQ Entry Model

```typescript
interface DLQEntry {
  _id: string;
  tenant_id: string;
  site_id: string;
  connector_id: string;
  action_id?: string;
  hook_id?: string;
  execution_id: string;
  payload: any;
  error: {
    message: string;
    code: string;
    stack?: string;
  };
  status: "pending" | "requeued" | "resolved";
  retry_count: number;
  max_retries: number;
  failed_at: Date;
  requeued_at?: Date;
  resolved_at?: Date;
}
```

---

## 12. Security

### Authentication

All admin endpoints require:
- **JWT Token:** Valid access token
- **Permissions:** `integrations:read`, `integrations:write`, `integrations:execute`

### Authorization

Role-based access control:
- **Admin:** Full access to all integration features
- **Integration Manager:** Manage connectors, actions, hooks
- **Integration Viewer:** Read-only access to executions, DLQ
- **Developer:** API access for webhook endpoints

### Data Protection

- **Credential Encryption:** AES-256 encryption for API keys, OAuth tokens
- **Secret Management:** Integration with HashiCorp Vault or AWS Secrets Manager
- **Audit Logging:** All configuration changes logged
- **IP Whitelisting:** Restrict webhook endpoints to known IPs

### Compliance

- **SOC2:** Audit logs for all integration executions
- **ISO 27001:** Encrypted storage of credentials
- **GDPR:** Data retention policies for execution history
- **PCI-DSS:** Secure handling of payment integration data

---

## 13. Operations

### Monitoring

Key metrics to monitor:
- **Execution Success Rate:** Target >99%
- **Average Response Time:** <2 seconds
- **DLQ Size:** Alert if >100 pending entries
- **Circuit Breaker Status:** Alert on open circuits
- **Connector Uptime:** Track availability of external services

### Alerts

Configure alerts for:
- Failed executions (>5 in 5 minutes)
- DLQ accumulation (>50 pending)
- Circuit breaker open (>3 minutes)
- High latency (>5 seconds p95)
- Authentication failures (>10 in 1 minute)

### Maintenance

Regular tasks:
- **Weekly:** Review DLQ entries, clear resolved items
- **Monthly:** Audit connector credentials, rotate API keys
- **Quarterly:** Review integration performance, optimize
- **Annually:** Security audit, compliance review

---

## 14. Performance

### Benchmarks

| Metric | Target | Current |
|--------|--------|---------|
| **Action Execution** | <2s | 1.2s avg |
| **Webhook Processing** | <500ms | 320ms avg |
| **DLQ Query** | <100ms | 65ms avg |
| **Connector Test** | <3s | 2.1s avg |
| **Throughput** | 1000 req/min | 850 req/min |

### Optimization Tips

1. **Connection Pooling:** Reuse HTTP connections to external services
2. **Caching:** Cache connector configurations in Redis
3. **Batch Processing:** Group multiple actions in single execution
4. **Async Processing:** Use RabbitMQ for non-blocking execution
5. **Circuit Breakers:** Fail fast when external service is down

---

## 15. Troubleshooting

### Common Issues

#### Issue: Action execution fails with timeout

**Cause:** External service is slow or unresponsive  
**Solution:**
1. Check circuit breaker status
2. Increase timeout in connector config
3. Verify external service availability
4. Check network connectivity

#### Issue: DLQ entries accumulating

**Cause:** Repeated failures due to configuration error  
**Solution:**
1. Review error messages in DLQ
2. Check connector credentials
3. Verify endpoint URLs
4. Test connector connection
5. Fix configuration and bulk requeue

#### Issue: Webhook not receiving data

**Cause:** Signature verification fails or IP not whitelisted  
**Solution:**
1. Verify webhook secret matches external system
2. Check signature algorithm (HMAC-SHA256)
3. Add external system IP to whitelist
4. Review webhook logs for error details

#### Issue: Duplicate executions

**Cause:** Idempotency key not provided  
**Solution:**
1. Add `Idempotency-Key` header to API requests
2. Generate unique key per request (UUID recommended)
3. Use same key for retries

---

## Document Control

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0.0 | 2026-03-01 | Platform Team | Initial documentation for integrations module |

**Next Review:** 2026-04-01  
**Owner:** Integration Team  
**Classification:** Internal - Enterprise Architecture

---

**End of Document**
