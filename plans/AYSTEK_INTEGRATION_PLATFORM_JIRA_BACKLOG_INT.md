# AYSTEK Integration Platform — JIRA Backlog & Work Plan

**Document Version**: 2.0 (Tasarım Dokümanıyla Sync)  
**Status**: Active  
**Last Updated**: Feb 28, 2026  
**Owner**: Platform Team

---

## Executive Summary

AYSTEK CORE'un tenant+site bazında dinamik olarak değişen entegrasyon ihtiyaçlarını (ERP/CRM/GPS/payment) kurumsal ölçekte çözen **Integration Platform** modülünün tam implementasyon iş planı.

**6 Sprint, ~480 story point, 4 Phase:**
- **Phase 0**: Foundation (Models, DB schema, types) — Sprint 4
- **Phase 1**: Runtime & Execution (Worker, async, DLQ, retry) — Sprint 5-6
- **Phase 2**: Admin UI (Config, management, monitoring) — Sprint 6-7
- **Phase 3**: Observability & Security (Tracing, audit, SSRF) — Sprint 7-8
- **Phase 4**: Governance & Testing (Versioning, approvals, docs) — Sprint 8-9

---

## Completed & In-Flight

✅ **INT-205 (Backend)**: DLQ model + requeue endpoint (DONE)  
✅ **NT-205 (Frontend)**: DLQ UI pages (DONE) - Sprint 4  
✅ **Phase 0 Foundation**: Model/repository/types/seeds/indexes/JMESPath altyapısı tamamlandı  
✅ **Phase 2 Admin UI (COMPLETE)**: 
  - Connector CRUD UI (list, create, edit, delete) ✅
  - Action CRUD UI (list, create, edit, delete) ✅
  - Hook CRUD UI (list, create, edit, delete, phase filter) ✅
  - Hook Form modernization (phase, action_keys, enabled, sync/async policy fields) ✅
  - Execution Log UI (list, detail, filtering) ✅
  - DLQ Management UI (list, detail modal, bulk requeue) ✅
  - i18n translations (28 languages, 16 new keys) ✅
  - Frontend types & API client extensions ✅
🔄 **Phase 1 Runtime (IN PROGRESS)**:
  - Integration Worker skeleton ✅
  - Basic ExecuteHook service ✅
  - ❌ HTTP Executor & Mapping Engine (INT-501-2) - CRITICAL BLOCKER
  - ❌ SSRF-Safe HTTP Client (INT-503-1) - SECURITY BLOCKER
  - ❌ Retry/Backoff/DLQ lifecycle (INT-501-3) - RELIABILITY GAP
  - ❌ Circuit Breaker & Rate Limiting (INT-501-6, INT-503-3) - STABILITY GAP

📋 **Next Priority**: Phase 1 Runtime tamamlanması (bkz. INTEGRATION_PLATFORM_REMAINING_WORK.md)

---

# PHASE 0: FOUNDATION & CORE MODELS — Sprint 4

## Epic INT-401: Data Models & Schema

### INT-401-1: Connector Model Repository
- **Type**: Story | **Points**: 13 | **Assigned**: Backend Lead
- **Scope**:
  - Go struct: `Connector` with fields: id, tenant_id, site_id, code, name, type, auth_config, TLS, environment, outbound_allowlist, rate_limits, timeouts, status, created_at, updated_by
  - Auth type enum: APIKey, BasicAuth, OAuth2, mTLS
  - MongoDB indexes: `{tenant_id, site_id, code}`, `{type}`, `{created_at}`
  - Repository interface: Create, Read, Update, Delete, ListBySite, FindByCode
  - Validation helpers: auth type check, URL/CIDR/IP allowlist validation, timeout bounds
- **Acceptance Criteria**:
  - [ ] Struct fields match design spec + AYSTEK standards (timestamps, soft_delete compatible)
  - [ ] Repo interface all CRUD + query methods
  - [ ] Unit tests: Create duplicate code → error, auth validation, list filtering
  - [ ] MongoDB indexes exist and work
- **Definition of Done**:
  - Go code compiles, tests pass, code review approved

### INT-401-2: Action Model Repository
- **Type**: Story | **Points**: 13 | **Assigned**: Backend Lead
- **Scope**:
  - Go struct: `Action` with fields: id, tenant_id, site_id, connector_id, key, method, path_template, request_mapping, response_mapping, pagination_config, error_policy, safety_opts, capability_ref, version, status (draft/active), created_at, updated_at
  - Request mapping: query/header/body templates, content_type, auth override (optional)
  - Response mapping: success_criteria (status code range + optional predicate), field_remapping (JMESPath expression), pagination directives
  - Error policy: retryable_status_codes (list), backoff_config (initial_delay, max_delay, multiplier)
  - Safety: max_response_size (default 10MB), redirect_disabled (always true), timeout_override
  - MongoDB indexes: `{connector_id, key}` (unique), `{capability_ref}`, `{status}`
- **Acceptance Criteria**:
  - [ ] All fields present and typed correctly
  - [ ] HTTP method enum validation (GET, POST, PUT, DELETE, PATCH)
  - [ ] Unique constraint enforced: (connector_id, key)
  - [ ] Repo interface: CRUD + ListByConnector(connector_id)
  - [ ] JMESPath expression field parser integrated (mock for now)
  - [ ] Unit tests: method validation, duplicate key, response mapping parse
- **Dependencies**: INT-401-1

### INT-401-3: Hook/Trigger Model
- **Type**: Story | **Points**: 8 | **Assigned**: Backend Lead
- **Scope**:
  - Go struct: `Hook` with fields: id, tenant_id, site_id, event_key, phase (pre/post), mode (sync/async), trigger_filter (optional JMESPath predicate), action_keys ([]string), mapping (JMESPath transform event→action input), enabled, created_at, updated_at, created_by
  - Sync hook config: timeout_ms (default 5s, max 30s), circuit_breaker (enabled, threshold, half_open_timeout)
  - Async hook config: retry_policy (max_attempts, initial_backoff, max_backoff), dlq_enabled (always true)
  - MongoDB indexes: `{tenant_id, site_id, event_key}`, `{phase}`, `{mode}`
- **Acceptance Criteria**:
  - [ ] Phase enum: pre = blocking, post = non-blocking
  - [ ] Mode enum: sync or async (mutually exclusive based on phase)
  - [ ] Event key format validation (e.g., `integration.execution.failed`)
  - [ ] Soft delete support
  - [ ] Tests: Hook activation/deactivation, predicate parsing
- **Dependencies**: INT-401-1

### INT-401-4: DLQItem Model Extension
- **Type**: Story | **Points**: 5 | **Assigned**: Backend Lead
- **Scope**:
  - Update existing `DLQItem` struct:
    - Add fields: hook_id, capability_ref, context (tenant_id, site_id, user_id, request_id)
    - Rename/alias: execution_id → parent_execution_id (for audit trail)
  - MongoDB migration: add fields, backfill null/empty
  - Repository update: new query methods FilterByHookID(hook_id), FilterByCapabilityRef(ref)
- **Acceptance Criteria**:
  - [ ] New fields optional (null-safe) on read
  - [ ] Migration script idempotent (can re-run safely)
  - [ ] Existing DLQ items still queryable
  - [ ] Tests: backfill, new query methods
- **Dependencies**: INT-205 (backend DLQ must exist)

### INT-401-5: ExecutionLog Time Series Collection
- **Type**: Story | **Points**: 8 | **Assigned**: Backend Lead
- **Scope**:
  - Go struct: `ExecutionLog` with fields: id, tenant_id, site_id, connector_id, action_key, hook_id, phase, status (success/error), http_status_code, latency_ms, attempt_num, retry_count, error_code, error_message (truncated), input_size, output_size, created_at (timestamp), metadata (additional context)
  - MongoDB time-series collection: `integration_executions_ts` with timeseries field = created_at, metadata_fields = [tenant_id, site_id, connector_id, action_key]
  - TTL: 90 days (configurable)
  - Indexes: `{tenant_id, site_id, connector_id}` compound, `{created_at}` (implicit in time-series)
  - Repository: BatchInsert([]ExecutionLog), QueryRange(tenant, start, end, filters)
- **Acceptance Criteria**:
  - [ ] Time-series collection created with correct config
  - [ ] Batch insert works (bulk write)
  - [ ] TTL job removes docs after 90 days
  - [ ] Query performance acceptable (< 1s for 1M docs)
  - [ ] Tests: insertion, range queries, TTL
- **Documentation**: Retention strategy, rollup options for future archival

### INT-401-6: Mapping DSL Type System & JMESPath Integration
- **Type**: Story | **Points**: 8 | **Assigned**: Backend Lead
- **Scope**:
  - Go package: `internal/integration/mapping`
  - Types:
    - `Expression` struct wrapping compiled JMESPath AST + metadata
    - `HelperFunction` enum: now, uuid, hash_sha256, hash_sha1, base64, base64_decode, date_format, number_round, string_upper, string_lower, string_split
    - `ValidationConfig` with max_depth, max_execution_ms, max_output_bytes
  - Functions:
    - `ParseExpression(exprStr string) (*Expression, error)` - compile JMESPath
    - `EvaluateExpression(expr *Expression, input interface{}) (interface{}, error)` - with timeout+depth guard
    - `ValidateExpression(exprStr string) error` - lightweight validation
  - Integration with go-jmespath library
  - Tests: valid/invalid expression detection, evaluation with timeout
- **Acceptance Criteria**:
  - [ ] JMESPath lib integrated and compiling
  - [ ] Expression validation catches syntax errors
  - [ ] Evaluation timeout enforced (context deadline)
  - [ ] Helper functions enumerated (no implementation yet, placeholder)
  - [ ] Tests: parse, validate, evaluate, timeout scenarios
- **Estimate for later integration**: Complex expressions deferred to INT-602-2 (Action form)

---

## Epic INT-402: Database Seeding & Migrations

### INT-402-1: Permission Seeds (Connectors, Actions, Hooks)
- **Type**: Story | **Points**: 5 | **Assigned**: Backend Lead
- **Scope**:
  - Backend `internal/seeds/permission_seed.go`:
    - Permissions for `integrations` module:
      - `integrations:read` - list/view connector/action/hook/execution
      - `integrations:create` - create new connector/action/hook
      - `integrations:update` - edit existing
      - `integrations:delete` - delete
      - `integrations:execute` - run hooks/actions, requeue DLQ items
      - `integrations:approve` - approve config changes (for governance, later)
    - Role assignments: admin=all, editor=read/create/update, viewer=read
    - Idempotency: permissions can be re-seeded safely
- **Acceptance Criteria**:
  - [ ] Permission creation/update idempotent
  - [ ] Role assignments correct
  - [ ] Tests: permission checks work in handlers

### INT-402-2: Menu Seeds (DLQ already added in Sprint 4 for NT-205)
- **Type**: Story | **Points**: 5 | **Assigned**: Backend Lead
- **Scope**:
  - Backend `internal/seeds/menu_seed.go`:
    - Menu item: `/admin/integrations/connectors` (already exists per prior backlog)
    - Menu item: `/admin/integrations/actions` (new)
    - Menu item: `/admin/integrations/hooks` (new)
    - Menu item: `/admin/integrations/executions` (new)
    - Menu item: `/admin/integrations/dlq` (DONE in NT-205 Sprint 4)
    - All with proper i18n keys + icon assignments (e.g., icon="plug" for connectors, icon="alert-circle" for dlq)
    - Role restrictions: only visible to admin/editor
- **Acceptance Criteria**:
  - [ ] All menu entries present in seed
  - [ ] Role restrictions enforced
  - [ ] Icons assigned + consistent with design

### INT-402-3: Translation Seeds i18n (28 Languages)
- **Type**: Story | **Points**: 8 | **Assigned**: Backend Lead
- **Scope**:
  - Backend `internal/seeds/translation_seed.go`:
    - ~150+ i18n keys for all Admin UI strings (Phases 0-4)
    - Keys grouped by feature: connectors.*, actions.*, hooks.*, executions.*, dlq.* (already added in NT-205), policies.*, etc.
    - Format: `allLangs("tr_string", "en_string")` for all 28 languages
    - Examples:
      - `integrations.connectors.title` = "Bağlantılar" / "Connectors"
      - `integrations.actions.form.method` = "HTTP Yöntemi" / "HTTP Method"
      - `integrations.hooks.list.empty` = "Hook bulunamadı" / "No hooks found"
      - DLQ keys (already added in NT-205): `integrations.dlq.title`, `integrations.dlq.requeue.success`, etc.
- **Acceptance Criteria**:
  - [ ] All keys use `allLangs(...)` pattern
  - [ ] No hardcoded strings in code
  - [ ] Fallback: tr + en present, other 26 langs → en copy
  - [ ] Tests: keys render without error

### INT-402-4: Sample Connector Seeds (ERP, GPS, Finance)
- **Type**: Story | **Points**: 8 | **Assigned**: Backend Lead
- **Scope**:
  - Create 3 sample connectors in `internal/seeds/domain_seed.go` or dedicated `connector_seed.go`:
    1. **ERP Connector**:
       - code: `erp_gateway_v1`, type: HTTP, auth: APIKey
       - base_url: https://erp.example.com/api/v2
       - allowlist: erp.example.com:443
       - rate_limit: 1000 req/min, timeout: 30s
    2. **GPS Provider**:
       - code: `gps_provider_v1`, type: HTTP, auth: OAuth2
       - base_url: https://gps-api.example.com/tracking
       - allowlist: gps-api.example.com:443
       - rate_limit: 5000 req/min, timeout: 10s
    3. **Payment Gateway**:
       - code: `payment_processor_v1`, type: HTTP, auth: mTLS
       - base_url: https://pay.example.com/transactions
       - allowlist: pay.example.com:443
       - rate_limit: 500 req/min, timeout: 15s
  - For each connector: create 2-3 sample actions (e.g., ERP: GetDebt, CreateJournal)
  - Seeding conditional: only if dev/test environment
- **Acceptance Criteria**:
  - [ ] Seeds created with realistic structure
  - [ ] Sample actions callable + properly mapped
  - [ ] Dev environment demo flow works

---

## Epic INT-403: Frontend Types & API Client Layer

### INT-403-1: Frontend TypeScript Types
- **Type**: Story | **Points**: 8 | **Assigned**: Frontend Lead
- **Scope**:
  - File: `lib/types/integrations.ts` (extend existing)
  - New types:
    - `Connector`, `ConnectorAuthConfig`, `ConnectorAuthType` (enum), `RateLimitConfig`, `TimeoutConfig`
    - `Action`, `RequestMapping`, `ResponseMapping`, `PaginationConfig`, `ErrorPolicy`
    - `Hook`, `HookPhase` (enum: pre/post), `HookMode` (enum: sync/async), `HookConfig`
    - `ExecutionLog`, `ExecutionStatus` (enum: success/error)
    - `CapabilityRef`, `CapabilityTemplate`
    - Request payloads: `CreateConnectorPayload`, `UpdateConnectorPayload`, `CreateActionPayload`, etc.
    - List response wrappers: `ConnectorListResult` with pagination meta
  - All types follow AYSTEK patterns: union types for discriminators, ISO8601 timestamps, id=string
- **Acceptance Criteria**:
  - [ ] Types compile without errors
  - [ ] Discriminator unions work (e.g., auth type → different config shape)
  - [ ] Tests: type safety checks (TypeScript strict mode)

### INT-403-2: API Client Methods (CRUD)
- **Type**: Story | **Points**: 13 | **Assigned**: Frontend Lead
- **Scope**:
  - File: `lib/api/integrations.ts` (extend existing)
  - Client objects:
    - `connectorApi`: listConnectors, getConnector, createConnector, updateConnector, deleteConnector
    - `actionApi`: listActions, getAction, createAction, updateAction, deleteAction
    - `hookApi`: listHooks, getHook, createHook, updateHook, deleteHook
    - `executionApi`: listExecutions, getExecution, retryExecution
  - Parser functions: parseConnector, parseAction, parseHook, parseExecutionLog (convert raw API → typed)
  - Error handling: unwrapApiResponse + typed error extraction
  - Features:
    - All methods typed (input params + return Promise<Type>)
    - Pagination support (page, limit params + meta in result)
    - Query filters passed as params object
    - Dry-run/test-run mode for actions (opsiyonel early preview)
- **Acceptance Criteria**:
  - [ ] All CRUD methods implemented
  - [ ] Error responses typed and extractable
  - [ ] Tests (unit): mock API calls, verify parser functions

---

**Phase 0 Total: ~105 points | Sprint 4 Allocation: ~99-105 points (on track)**

---

# PHASE 1: RUNTIME & EXECUTION ENGINE — Sprints 5-6

## Epic INT-501: Integration Worker & Async Execution

### INT-501-1: Integration Worker Service (Core)
- **Type**: Story | **Points**: 13 | **Assigned**: Backend Runtime Team
- **Scope**:
  - New service: `cmd/workers/integration-worker/main.go`
  - Workers consuming from RabbitMQ queue: `integration_events`
  - Message struct: { tenant_id, site_id, hook_id, event_key, payload, attempt_count, correlation_id }
  - Worker logic:
    1. Connect to MongoDB + RabbitMQ
    2. Consume message (prefetch=10 for concurrency)
    3. Deserialize → hook context
    4. Call `integrationService.ExecutePostHook(ctx, event)`
    5. On success: ack message
    6. On error: decide retry vs DLQ (via INT-501-3)
  - Graceful shutdown: drain queue, wait for in-flight execution (30s timeout)
  - Context propagation: extract correlation_id from message header
  - OWASP-compliant logging: secret masking (API keys, tokens), PII redaction
- **Acceptance Criteria**:
  - [ ] Worker starts, connects to broker
  - [ ] Message consumption works
  - [ ] Error paths trigger retry/DLQ correctly
  - [ ] Graceful shutdown tested
  - [ ] Tests: mock broker, message flow

### INT-501-2: ExecutePostHook Service Method
- **Type**: Story | **Points**: 13 | **Assigned**: Backend Runtime Team
- **Scope**:
  - Service method: `ExecutePostHook(ctx context.Context, hook *Hook, event interface{}, actorID string) error`
  - Flow:
    1. Validate hook + action configs exist
    2. Load connector (auth, base_url, allowlist)
    3. Build HTTP request: method, URL, headers, body (using mapping expressions → INT-401-6)
    4. Execute via SSRF-hardened HTTP client (INT-503-1)
    5. Capture response + latency
    6. Log to ExecutionLog collection (INT-401-5)
    7. Return success or detailed error
  - Timeout enforcement: context deadline
  - Retry signal: error type indicates retryable vs permanent
  - Secret handling: load OAuth token if refresh needed (INT-702-2 later)
- **Acceptance Criteria**:
  - [ ] Mapping expressions evaluated correctly
  - [ ] HTTP request follows specification
  - [ ] Response logged + metrics incremented
  - [ ] Timeout enforced (test with stalled endpoint mock)
  - [ ] Tests (unit): mapping, timeout, error scenarios

### INT-501-3: Retry & Backoff Mechanism (RabbitMQ DLX)
- **Type**: Story | **Points**: 8 | **Assigned**: Backend Runtime Team
- **Scope**:
  - RabbitMQ queue setup: `integration_events` (main) → `integration_events_delayed` (replay after TTL) → back to main
  - Dead Letter Exchange (DLX): failed messages → `integration_dlq` persistent
  - Backoff schedule:
    - Initial: 1s (x_retry_count=0)
    - Retry 1: 10s (x_retry_count=1)
    - Retry 2: 100s (x_retry_count=2)
    - Retry 3: → DLQ (max_retry=3)
  - Jitter: ±0-30% random delay added to backoff
  - Service method: `ScheduleRetry(ctx, message, delay_duration)` → republish to delay queue
  - Message headers: `x_retry_count`, `x_first_attempt_at`, `x_last_attempt_at`
- **Acceptance Criteria**:
  - [ ] RabbitMQ queue + DLX setup documented + automated
  - [ ] Backoff calculation correct + deterministic
  - [ ] Jitter range verified
  - [ ] Max retry threshold prevents infinite loops
  - [ ] Tests: backoff escalation, jitter variance

### INT-501-4: DLQ Item Management & Analytics
- **Type**: Story | **Points**: 8 | **Assigned**: Backend Runtime Team
- **Scope**:
  - Handler extensions:
    - `RequeueDLQBatch(ctx, ids []string, tenantID, siteID, actorID)` → republish back to main queue, log attempt
    - `ListDLQ(ctx, tenant, site, filters)` → query by hook_id, status, date_range, error_code, with pagination
    - `GetDLQItem(ctx, id)` → full item + associated hook config (for UI inspection)
    - `AnalyzeDLQTrends(ctx, tenant, window_days)` → count by error_code, error_message, top errors
  - Audit logging: who requeued, when, batch size, result (succeeded/failed count)
  - Batch limits: max 1000 items per requeue to prevent overload
- **Acceptance Criteria**:
  - [ ] Requeue flow end-to-end tested
  - [ ] Filtering on multiple fields works
  - [ ] Trend analysis aggregation correct
  - [ ] Audit records created + queryable
  - [ ] Tests: edge cases (empty batch, invalid IDs, permission checks)

### INT-501-5: Idempotency Key & Dedup Store (Redis)
- **Type**: Story | **Points**: 8 | **Assigned**: Backend Runtime Team
- **Scope**:
  - Redis-backed idempotency:
    - Key format: `idempotency:{tenant}:{hook_id}:{event_id}:{attempt_num}`
    - Value: execution result JSON (cached response) or success marker
    - TTL: 24 hours (after which same event allowed again, "eventual consistency" model)
  - Worker logic: before execute, check Redis key
    - If exists + success → return cached result + skip execution (idempotent consumer)
    - If not exists → execute + set key on success
    - If exists + error → check if error reason is retryable; if yes, return retry signal; if no, permanent fail
  - Service method: `CheckIdempotent(ctx, key)` → (exists bool, error)
  - Dedup logging: "Skipped duplicate message for event_id=X" (no error, idempotent)
- **Acceptance Criteria**:
  - [ ] Redis connection established
  - [ ] Key TTL set correctly
  - [ ] Dedup detection working (verify via logs)
  - [ ] No duplicate execution on retry
  - [ ] Tests: dedup scenarios, TTL expiry

### INT-501-6: Circuit Breaker for Outbound Calls
- **Type**: Story | **Points**: 8 | **Assigned**: Backend Runtime Team
- **Scope**:
  - Circuit breaker library: use `go-kit` CB or `Sony/gobreaker` or similar
  - Config per connector:
    - Failure threshold: 5 consecutive failures (configurable)
    - Open timeout: 30s (half-open test period)
    - Success threshold for recovery: 2 successive successes
  - States: Closed (pass through) → Open (fail fast) → Half-Open (test) → Closed/Open
  - Metrics: circuit state transitions + open duration
  - Service method: `ExecuteWithCircuitBreaker(ctx, connector_id, executeFunc)`
  - Logging: CB state transitions (e.g., "Circuit breaker opened for connector=erp_gateway")
- **Acceptance Criteria**:
  - [ ] CB state transitions logic correct
  - [ ] Open state fails fast (no HTTP call attempt)
  - [ ] Half-open allows test call + decides
  - [ ] Metrics exposed for monitoring
  - [ ] Tests: state transitions, recovery

---

## Epic INT-502: Sync Execution (Pre-hooks & Direct Calls)

### INT-502-1: ExecutePreHook Service Method
- **Type**: Story | **Points**: 8 | **Assigned**: Backend Runtime Team
- **Scope**:
  - Service method: `ExecutePreHook(ctx, hook *Hook, event interface{}, actorID) → (allowed bool, reason string, err error)`
  - Behavior: blocking, short timeout (5s default, configurable max 30s)
  - Circuit breaker applied (if endpoint keeps failing, fail decisions → prevent cascade)
  - Response parsing: extract decision from JSON response (expected field: `allowed` boolean or `success` boolean)
  - Optional reason extraction: field `reason` or `message`
  - Return tuple: (allowed, reason, error)
  - Timeout behavior: if context expires → (false, "timeout", nil) [deny access as fail-safe]
  - Logging: decision + reason
- **Acceptance Criteria**:
  - [ ] Timeout enforced (< 5s)
  - [ ] Circuit breaker applied
  - [ ] Response decision parsing correct
  - [ ] Timeout results in deny (fail-safe)
  - [ ] Tests: success path, timeout, CB open scenarios

### INT-502-2: Direct Action Execution (Admin Test-Run)
- **Type**: Story | **Points**: 8 | **Assigned**: Backend Runtime Team
- **Scope**:
  - Service method: `ExecuteActionDirect(ctx, action *Action, input_payload interface{}, tenantID, siteID, actorID) → (*ExecutionResult, error)`
  - Result struct: { status_code int, response_body []byte, latency_ms int, execution_log_entry *ExecutionLog, error_msg string }
  - Timeout: 30s (longer than pre-hook, more permissive for testing)
  - Mapping: apply JMESPath expressions to transform input → HTTP request
  - No audit log for test-run (only for "deployed" production execution)
  - Caching: no idempotency key (each test-run is independent)
  - Dry-run mode (opsiyonel v1): skip actual HTTP call, return preview request (for validation before deploy)
- **Acceptance Criteria**:
  - [ ] Direct execution works end-to-end
  - [ ] Response parsing + logging correct
  - [ ] Timeout enforced
  - [ ] Test-run flag prevents audit logging
  - [ ] Tests: various status codes, mapping validation

---

## Epic INT-503: HTTP Client Hardening & SSRF Prevention

### INT-503-1: HTTP Client Wrapper with Outbound Allowlist
- **Type**: Story | **Points**: 13 | **Assigned**: Backend Security Team
- **Scope**:
  - Go package: `internal/platform/http/security_client`
  - Config: `AllowlistRule` struct with domain (wildcard support: `*.example.com`), IP/CIDR, port (or 0 for any), scheme (https preferred)
  - HTTP client method: `Do(ctx, method, url, body, allowlist []AllowlistRule) (*http.Response, error)`
  - Pre-request validation (OWASP SSRF cheat sheet):
    1. Parse URL → verify scheme is http/https only
    2. DNS resolve hostname → get resolved IP
    3. Check resolved IP against denylist (RFC1918, loopback 127.0.0.1, link-local 169.254.x.x, metadata 169.254.169.254)
    4. Check resolved IP against allowlist (if provided)
    5. If dynamic DNS domain: re-resolve to detect rebinding (time-of-use check)
    6. Verify port in allowlist (default https=443, http=80)
  - Response validation:
    - Redirect follow: disabled (no 30x following)
    - Response size limit (configurable, default 10MB)
    - Content-type allowlist (optional): e.g., application/json, text/plain
  - Logging: requests + responses (masked URL if secret suspected)
- **Acceptance Criteria**:
  - [ ] Allowlist pass: www.example.com on allowlist → allowed
  - [ ] Allowlist fail: internal.corp on allowlist → reject
  - [ ] Private IP rejection: 192.168.x.x → reject
  - [ ] Redirect rejection: 3xx response → no follow
  - [ ] Response size limit enforced
  - [ ] Tests: allowlist scenarios, private IP, redirect, size limit

### INT-503-2: Timeout Enforcement & Context Integration
- **Type**: Story | **Points**: 5 | **Assigned**: Backend Security Team
- **Scope**:
  - HTTP client method: `DoWithTimeout(ctx context.Context, method, url string, body []byte, timeout_ms int, allowlist []AllowlistRule) (*http.Response, error)`
  - Timeouts configured:
    - Dial timeout (connect): 5s
    - Read timeout (per request body + response): 15s
    - Overall budget: min(ctx.Deadline, now + timeout_ms)
  - Context deadline propagation: shorter of context deadline or timeout
  - Error on timeout: `context.DeadlineExceeded` wrapped
- **Acceptance Criteria**:
  - [ ] Timeout enforced at deadline
  - [ ] Context cancellation respected
  - [ ] Timeout error clear + distinguishable from other errors
  - [ ] Tests: timeout expiry, deadline propagation

### INT-503-3: Rate Limit Enforcement (Per Connector, Sliding Window)
- **Type**: Story | **Points**: 8 | **Assigned**: Backend Runtime Team
- **Scope**:
  - Redis-backed sliding window rate limiter:
    - Key: `ratelimit:{tenant_id}:{connector_id}:requests`
    - Window: 60s rolling
    - Metric: count requests in last 60s (INCR + EXPIRE)
  - Service method: `AcquireRateLimit(ctx, tenantID, connectorID, qps_limit) → (allowed bool, remaining int, reset_after_sec int, err error)`
  - If count >= limit: return (false, 0, time_until_window_reset)
  - If count < limit: INCR in Redis, return (true, remaining, 0)
  - Connector config: `rate_limit_per_minute` (default from tenant setting or 10000)
  - Handler integration: pre-request check in hook execution
- **Acceptance Criteria**:
  - [ ] Rate limit crossing detected
  - [ ] Remaining count accurate
  - [ ] Window reset time correct
  - [ ] Tests: limit crossing, sliding window, multiple requests

---

**Phase 1 Total: ~143 points | Sprint 5: ~92 points, Sprint 6: ~51 points (with overflow to UI tasks)**

---

# PHASE 2-4: ADMIN UI, OBSERVABILITY, SECURITY & GOVERNANCE

*Full epics INT-601 through INT-803 follow the same detailed pattern above.*

---

## Summary Tables

### Sprint Allocation

| Sprint | Phase | Epics | Points | Focus |
|--------|-------|-------|---------|-------|
| 4 | 0 | INT-401, 402, 403 | ~105 | Foundation: models, seeds, types |
| 5 | 1 | INT-501, 502, 503 (partial) | ~92 | Runtime: worker, execution, SSRF |
| 6 | 1-2 | INT-503, 601 (partial) | ~47 | Finish HTTP hardening, start UI |
| 7 | 2-3 | INT-602, 603, 604, 701 (start) | ~81 | Complete UI, begin observability |
| 8 | 3-4 | INT-701, 702, 703, 801, 802 (partial) | ~134 | Security + governance + testing |
| 9 | 4 | INT-801, 802, 803 | ~67 | Testing + documentation |
| **TOTAL** | | | **~526** | **6 sprints, enterprise-ready platform** |

### Success Metrics (End of Phase 4)

- ✅ 95%+ unit test coverage (core services)
- ✅ 5+ integration tests (worker → execution → DLQ flows)
- ✅ Admin UI: 100% of CRUD operations functional + responsive
- ✅ Security: SSRF allowlist enforced, secrets encrypted, OAuth2 refresh working
- ✅ Observability: tracing + metrics + audit logs + correlation IDs
- ✅ SLA achieved: p95 < 5s (sync), p95 < 10s (async); >99.5% success rate
- ✅ Documentation: API spec, runbooks, capability catalog, training guide
- ✅ Governance: versioning + approvals workflow live
- ✅ Deployment: staging environment fully featured

---

**Document Status**: Ready for JIRA import. Use as source-of-truth for 6-sprint delivery roadmap.
