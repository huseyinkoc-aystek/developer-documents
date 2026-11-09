# AYSTEK Integration Platform — Kalan İşler ve Öncelik Sıralaması

**Oluşturulma Tarihi**: 28 Şubat 2026  
**Son Durum Özeti**: Phase 0 ve Phase 2 Admin UI tamamlandı, Phase 1 Runtime kısmi tamamlandı  
**Hedef**: Production-ready kurumsal entegrasyon platformu

---

## 📊 Genel Durum Özeti

### ✅ Tamamlanmış (Phase 0 + Phase 2)
- **Backend Models & Repositories**: Connector, Action, Hook, Execution, DLQ modelleri tam
- **Seeds & Permissions**: 28 dil i18n, RBAC permissions, menu entries
- **Admin UI - CRUD**: Connector, Action, Hook, Execution Log, DLQ tam yönetim arayüzleri
- **Hook Modernization**: Phase 0/1 field alignment (phase, action_keys, enabled, sync/async policy)
- **Integration Worker**: Temel worker altyapısı kuruldu
- **Service Layer**: ExecuteHook, RetryExecution, DLQ management metodları

### 🔶 Kısmi Tamamlanmış (Phase 1)
- **Worker Infrastructure**: Worker skeleton var, ama HTTP executor eksik
- **Execution Engine**: Hook execution temel mantık var, ama mapping/transformation eksik
- **Retry/DLQ**: DLQ model ve UI tam, ama retry mekanizması eksik

### ❌ Eksik (Phase 1 + Phase 3 + Phase 4)
- **HTTP Client Hardening**: SSRF-safe client eksik
- **Mapping DSL**: JMESPath evaluation engine eksik
- **Circuit Breaker & Rate Limiting**: Henüz yok
- **Observability**: Metrics/tracing kısmi
- **OAuth2 Refresh**: Token yenileme yok
- **Versioning & Governance**: Henüz yok
- **Test-Run UI**: Action ve Hook test-run modalları eksik

---

## 🎯 Phase 1: Runtime & Execution Engine (Kritik - Production Blocker)

### **INT-501-2: HTTP Executor & Mapping Engine** [Priority: P0] [Points: 21]
**Durum**: 🔴 Eksik  
**Kapsam**:
- [ ] JMESPath expression evaluator integration
  - Go library: github.com/jmespath/go-jmespath
  - Expression validation (syntax check)
  - Evaluation with timeout (max 50ms)
  - Helper functions: `now()`, `uuid()`, `hash()`, `base64()`
- [ ] HTTP Request Builder
  - Template rendering (Mustache-like)
  - Query parameter injection
  - Header injection (auth headers, custom headers)
  - Body transformation (JSON → HTTP request body)
- [ ] HTTP Response Parser
  - Status code validation
  - Response body extraction
  - JMESPath mapping (remote response → canonical format)
  - Error extraction and classification (retryable vs permanent)
- [ ] ExecutePostHook Service Implementation
  - Load connector config (base_url, auth, allowlist)
  - Build request via mapping
  - Execute HTTP call
  - Parse response
  - Log to ExecutionLog collection
  - Return success/error signal
- [ ] Tests: mapping evaluation, HTTP request building, response parsing, timeout scenarios

**Bağımlılıklar**: INT-503-1 (SSRF-safe HTTP client)  
**Etki**: Bu olmadan hooks hiçbir şey yapmıyor (production blocker)  
**Tahmini Süre**: 2-3 gün

---

### **INT-503-1: SSRF-Safe HTTP Client** [Priority: P0] [Points: 13]
**Durum**: 🔴 Eksik  
**Kapsam**:
- [ ] `internal/platform/http/security_client` package
- [ ] Outbound allowlist validation
  - Domain wildcard support (*.example.com)
  - IP/CIDR allowlist
  - Port filtering
  - Scheme validation (https preferred)
- [ ] SSRF prevention (OWASP guidelines)
  - DNS resolution + private IP check (RFC1918, loopback, link-local)
  - Metadata endpoint blocking (169.254.169.254)
  - Redirect disabled (no 30x following)
  - DNS rebinding detection
- [ ] Response validation
  - Size limit (default 10MB)
  - Content-type allowlist (optional)
- [ ] Timeout enforcement
  - Connect timeout (5s)
  - Read timeout (15s)
  - Overall request budget
- [ ] Context propagation (correlation IDs, tracing)
- [ ] Logging (masked secrets)
- [ ] Tests: allowlist pass/fail, private IP rejection, redirect block, size limit

**Bağımlılıklar**: Yok  
**Etki**: Security blocker - bu olmadan production'a çıkamaz  
**Tahmini Süre**: 2 gün

---

### **INT-501-3: Retry & Backoff Mechanism** [Priority: P0] [Points: 8]
**Durum**: 🔴 Eksik  
**Kapsam**:
- [ ] RabbitMQ DLX (Dead Letter Exchange) setup
  - `integration_events` → `integration_events_delayed` → back to main
- [ ] Exponential backoff schedule
  - Retry 0: 1s
  - Retry 1: 10s
  - Retry 2: 100s
  - Retry 3: → DLQ (max_retry=3)
- [ ] Jitter (±30% random delay)
- [ ] Message headers: `x_retry_count`, `x_first_attempt_at`, `x_last_attempt_at`
- [ ] Service method: `ScheduleRetry(ctx, message, delay_duration)`
- [ ] Tests: backoff calculation, jitter variance, max retry threshold

**Bağımlılıklar**: INT-501-1 (Worker varsa bağlanır)  
**Etki**: Güvenilirlik - bu olmadan failed execution'lar DLQ'ya düşer ama retry yok  
**Tahmini Süre**: 1.5 gün

---

### **INT-501-5: Idempotency Key & Dedup** [Priority: P1] [Points: 8]
**Durum**: 🔴 Eksik  
**Kapsam**:
- [ ] Redis-backed idempotency store
  - Key format: `idempotency:{tenant}:{hook_id}:{event_id}:{attempt}`
  - TTL: 24 hours
- [ ] Worker logic: check Redis before execute
  - If exists + success → skip execution (return cached result)
  - If not exists → execute + set key on success
- [ ] Service method: `CheckIdempotent(ctx, key) → (exists bool, error)`
- [ ] Dedup logging: "Skipped duplicate message for event_id=X"
- [ ] Tests: dedup detection, TTL expiry, concurrent access

**Bağımlılıklar**: INT-501-2 (HTTP executor)  
**Etki**: Data correctness - duplicate execution'ları engeller  
**Tahmini Süre**: 1.5 gün

---

### **INT-501-6: Circuit Breaker** [Priority: P1] [Points: 8]
**Durum**: 🔴 Eksik  
**Kapsam**:
- [ ] Circuit breaker library integration (Sony/gobreaker veya go-kit)
- [ ] Config per connector:
  - Failure threshold: 5 consecutive failures
  - Open timeout: 30s (half-open test period)
  - Success threshold for recovery: 2 successes
- [ ] States: Closed → Open → Half-Open → Closed/Open
- [ ] Service method: `ExecuteWithCircuitBreaker(ctx, connector_id, executeFunc)`
- [ ] Metrics: circuit state transitions, open duration
- [ ] Logging: state changes ("Circuit breaker opened for connector=X")
- [ ] Tests: state transitions, fail-fast in open state, recovery

**Bağımlılıklar**: INT-503-1 (HTTP client)  
**Etki**: Stability - cascading failure prevention  
**Tahmini Süre**: 1.5 gün

---

### **INT-503-3: Rate Limit Enforcement** [Priority: P1] [Points: 8]
**Durum**: 🔴 Eksik  
**Kapsam**:
- [ ] Redis-backed sliding window rate limiter
  - Key: `ratelimit:{tenant_id}:{connector_id}:requests`
  - Window: 60s rolling
  - Metric: count requests in last 60s
- [ ] Service method: `AcquireRateLimit(ctx, tenantID, connectorID, qps_limit) → (allowed, remaining, reset_after)`
- [ ] Connector config: `rate_limit_per_minute` (default 10000)
- [ ] Handler integration: pre-request check in hook execution
- [ ] Tests: limit crossing, sliding window, multiple concurrent requests

**Bağımlılıklar**: INT-501-2 (HTTP executor)  
**Etki**: API quota compliance, denial of service prevention  
**Tahmini Süre**: 1.5 gün

---

### **INT-502-1: ExecutePreHook (Sync Blocking)** [Priority: P2] [Points: 8]
**Durum**: 🔴 Eksik (Post-hook var, pre-hook eksik)  
**Kapsam**:
- [ ] Service method: `ExecutePreHook(ctx, hook, event, actorID) → (allowed bool, reason string, err error)`
- [ ] Blocking behavior: short timeout (5s default, max 30s)
- [ ] Circuit breaker applied
- [ ] Response parsing: extract `allowed` or `success` boolean
- [ ] Timeout behavior: deny access as fail-safe
- [ ] Logging: decision + reason
- [ ] Tests: success path, timeout, circuit breaker open scenarios

**Bağımlılıklar**: INT-503-1 (HTTP client), INT-501-6 (Circuit breaker)  
**Etki**: Pre-hook workflows (ör. fraud check before payment)  
**Tahmini Süre**: 1 gün

---

### **INT-502-2: Direct Action Test-Run** [Priority: P2] [Points: 8]
**Durum**: 🔴 Eksik  
**Kapsam**:
- [ ] Service method: `ExecuteActionDirect(ctx, action, input_payload, tenantID, siteID, actorID) → (*ExecutionResult, error)`
- [ ] Result struct: { status_code, response_body, latency_ms, execution_log_entry, error_msg }
- [ ] Timeout: 30s (longer than pre-hook)
- [ ] Mapping: JMESPath transform input → HTTP request
- [ ] No audit log (test-run only)
- [ ] No idempotency key (each test independent)
- [ ] Dry-run mode (optional): skip HTTP call, return preview request
- [ ] Tests: various status codes, mapping validation

**Bağımlılıklar**: INT-501-2 (HTTP executor)  
**Etki**: Admin UI test-run feature (nice-to-have, not blocker)  
**Tahmini Süre**: 1 gün

---

## 🔍 Phase 3: Observability & Security (Post-MVP)

### **INT-701: Tracing & Metrics** [Priority: P2] [Points: 13]
**Durum**: 🟡 Kısmi (Temel loglar var, structured metrics eksik)  
**Kapsam**:
- [ ] OpenTelemetry context propagation
  - W3C Trace Context headers (traceparent/tracestate)
  - Correlation ID injection
  - Span creation per execution
- [ ] Prometheus metrics
  - `integration_execution_total{connector,action,status}`
  - `integration_execution_duration_seconds{connector,action}`
  - `integration_retry_total{reason}`
  - `integration_dlq_total`
  - `integration_circuit_breaker_state{connector}`
- [ ] Structured logging (OWASP guidelines)
  - Secret masking (API keys, OAuth tokens)
  - PII redaction
  - Security relevant events marked

**Bağımlılıklar**: INT-501-2 (execution engine)  
**Etki**: Operability, debugging  
**Tahmini Süre**: 2 gün

---

### **INT-702: OAuth2 Token Refresh** [Priority: P2] [Points: 13]
**Durum**: 🔴 Eksik  
**Kapsam**:
- [ ] OAuth2 credential storage (encrypted in MongoDB)
- [ ] Token expiry detection
- [ ] Refresh token flow (RFC 6749)
  - Endpoint: connector.auth_config.token_url
  - Refresh request with refresh_token
  - New access_token + refresh_token storage
- [ ] Token revocation support (RFC 7009)
- [ ] Rotation strategy (HashiCorp Vault integration optional)
- [ ] Service method: `RefreshOAuthToken(ctx, connector_id) → error`
- [ ] Tests: token expiry, refresh success/failure, rotation

**Bağımlılıklar**: INT-503-1 (HTTP client for token endpoint)  
**Etki**: Long-lived integration sustainability  
**Tahmini Süre**: 2 gün

---

### **INT-703: Webhook Signature Verification** [Priority: P3] [Points: 8]
**Durum**: 🔴 Eksik  
**Kapsam**:
- [ ] Inbound webhook handler (`POST /webhooks/:hook_id`)
- [ ] Signature verification methods:
  - HMAC-SHA256 (header: X-Signature)
  - Timestamp validation (replay attack prevention)
  - Secret key per hook
- [ ] Request validation
  - Body hash verification
  - Timestamp freshness (max 5 minutes)
- [ ] Tests: valid signature, invalid signature, replay attack, expired timestamp

**Bağımlılıklar**: Yok (independent feature)  
**Etki**: Inbound webhook security  
**Tahmini Süre**: 1.5 gün

---

## 🏛️ Phase 4: Governance & Testing (Future Enhancement)

### **INT-801: Versioning & Approvals** [Priority: P3] [Points: 13]
**Durum**: 🔴 Eksik  
**Kapsam**:
- [ ] Config version tracking (`version` field in Connector/Action/Hook)
- [ ] Draft → Review → Active workflow
- [ ] Approval gate for production deployment
- [ ] Rollback capability (restore previous version)
- [ ] Audit trail (who approved, when, what changed)

**Bağımlılıklar**: Yok  
**Etki**: Governance, compliance  
**Tahmini Süre**: 2 gün

---

### **INT-802: Integration Testing Suite** [Priority: P2] [Points: 13]
**Durum**: 🔴 Eksik  
**Kapsam**:
- [ ] Worker → Execution → DLQ end-to-end flows
- [ ] Mock HTTP server for external API simulation
- [ ] Retry behavior testing (exponential backoff validation)
- [ ] Circuit breaker state transitions
- [ ] Idempotency dedup scenarios
- [ ] Load testing (1000 concurrent executions)
- [ ] Target: 95%+ code coverage for core services

**Bağımlılıklar**: INT-501-* (all runtime features)  
**Etki**: Quality assurance, regression prevention  
**Tahmini Süre**: 3 gün

---

### **INT-803: Documentation & Runbooks** [Priority: P2] [Points: 8]
**Durum**: 🟡 Kısmi (Tasarım dokümanları var, operational runbooks eksik)  
**Kapsam**:
- [ ] API specification (OpenAPI/Swagger)
- [ ] Capability catalog (sample connectors + actions)
- [ ] Operational runbooks:
  - Worker deployment & scaling
  - DLQ triage & requeue procedures
  - Circuit breaker manual override
  - Incident response (high latency, rate limit exceeded)
- [ ] Admin training guide (screenshot'lu Admin UI walkthrough)
- [ ] Developer guide (custom connector ekleme)

**Bağımlılıklar**: Yok  
**Etki**: Adoption, maintainability  
**Tahmini Süre**: 2 gün

---

## 🎨 Phase 2 Admin UI Enhancements (Nice-to-Have)

### **INT-602-2: Action Test-Run Modal** [Priority: P3] [Points: 8]
**Durum**: 🔴 Eksik (Action CRUD tam, test-run UI eksik)  
**Kapsam**:
- [ ] Modal component: `ActionTestRunModal.tsx`
- [ ] Input payload editor (JSON textarea with validation)
- [ ] "Test Run" button in Action list page
- [ ] Result viewer:
  - HTTP status code
  - Response body (syntax highlighted JSON)
  - Latency (ms)
  - Request preview (method, URL, headers, body)
- [ ] Error display (user-friendly error messages)
- [ ] Tests: modal open/close, payload validation, API integration

**Bağımlılıklar**: INT-502-2 (Backend ExecuteActionDirect method)  
**Etki**: Admin UX improvement (not blocker)  
**Tahmini Süre**: 0.5 gün

---

### **INT-603-2: Hook Manual Test Modal** [Priority: P3] [Points: 8]
**Durum**: 🔴 Eksik (Hook CRUD tam, test execution UI eksik)  
**Kapsam**:
- [ ] Modal component: `HookTestExecutionModal.tsx`
- [ ] Test payload input (JSON editor)
- [ ] "Execute Test" button in Hook list page
- [ ] Execution result viewer:
  - Success/failure status
  - Response body
  - Latency
  - Action execution chain (which actions ran)
  - Error messages
- [ ] Request/response inspector (collapsible sections)
- [ ] Tests: modal integration, payload validation

**Bağımlılıklar**: INT-501-2 (Backend ExecuteHook method tam çalışır durumda)  
**Etki**: Admin UX improvement (not blocker)  
**Tahmini Süre**: 0.5 gün

---

### **INT-604-3: Execution Retry UI Enhancement** [Priority: P4] [Points: 5]
**Durum**: 🟡 Kısmi (Retry endpoint var, batch retry yok)  
**Kapsam**:
- [ ] Batch selection in Execution list (checkboxes)
- [ ] "Retry Selected" button
- [ ] Confirmation modal with retry count limit
- [ ] Progress indicator (X/Y retried)
- [ ] Tests: batch retry flow

**Bağımlılıklar**: Backend RetryExecution method (zaten var)  
**Etki**: Operational convenience  
**Tahmini Süre**: 0.5 gün

---

## 📊 Story Point ve Süre Tahminleri

### Phase 1 (Production Blocker)
| Epic | Story | Points | Tahmini Süre | Öncelik |
|------|-------|--------|--------------|---------|
| INT-501-2 | HTTP Executor & Mapping | 21 | 2-3 gün | P0 🔴 |
| INT-503-1 | SSRF-Safe HTTP Client | 13 | 2 gün | P0 🔴 |
| INT-501-3 | Retry & Backoff | 8 | 1.5 gün | P0 🔴 |
| INT-501-5 | Idempotency & Dedup | 8 | 1.5 gün | P1 🟡 |
| INT-501-6 | Circuit Breaker | 8 | 1.5 gün | P1 🟡 |
| INT-503-3 | Rate Limiting | 8 | 1.5 gün | P1 🟡 |
| INT-502-1 | ExecutePreHook | 8 | 1 gün | P2 🟢 |
| INT-502-2 | Direct Action Test-Run | 8 | 1 gün | P2 🟢 |
| **TOPLAM** | | **82 points** | **~12-14 gün** | |

### Phase 3 (Post-MVP)
| Epic | Story | Points | Tahmini Süre | Öncelik |
|------|-------|--------|--------------|---------|
| INT-701 | Tracing & Metrics | 13 | 2 gün | P2 🟢 |
| INT-702 | OAuth2 Refresh | 13 | 2 gün | P2 🟢 |
| INT-703 | Webhook Verification | 8 | 1.5 gün | P3 🔵 |
| **TOPLAM** | | **34 points** | **~5-6 gün** | |

### Phase 4 (Future)
| Epic | Story | Points | Tahmini Süre | Öncelik |
|------|-------|--------|--------------|---------|
| INT-801 | Versioning & Approvals | 13 | 2 gün | P3 🔵 |
| INT-802 | Integration Tests | 13 | 3 gün | P2 🟢 |
| INT-803 | Documentation | 8 | 2 gün | P2 🟢 |
| **TOPLAM** | | **34 points** | **~7 gün** | |

### Admin UI Enhancements
| Epic | Story | Points | Tahmini Süre | Öncelik |
|------|-------|--------|--------------|---------|
| INT-602-2 | Action Test-Run Modal | 8 | 0.5 gün | P3 🔵 |
| INT-603-2 | Hook Test Modal | 8 | 0.5 gün | P3 🔵 |
| INT-604-3 | Execution Batch Retry | 5 | 0.5 gün | P4 ⚪ |
| **TOPLAM** | | **21 points** | **~1.5 gün** | |

---

## 🎯 Önerilen Çalışma Sırası (Minimum Viable Product'a Giden Yol)

### Sprint A: Core Runtime (MVP Kritik) — ~6-7 gün
**Hedef**: Hook execution end-to-end çalışır hale getirmek

1. **INT-503-1: SSRF-Safe HTTP Client** [2 gün] ← Önce güvenlik
2. **INT-501-2: HTTP Executor & Mapping Engine** [2-3 gün] ← Sonra execution
3. **INT-501-3: Retry & Backoff Mechanism** [1.5 gün] ← Güvenilirlik

**Sonuç**: Hooks execute edilebilir, retry ile DLQ'ya düşer, güvenli HTTP client var → **Minimal Production Ready**

---

### Sprint B: Reliability & Performance — ~4-5 gün
**Hedef**: Production SLA'yı karşılamak

4. **INT-501-5: Idempotency & Dedup** [1.5 gün]
5. **INT-501-6: Circuit Breaker** [1.5 gün]
6. **INT-503-3: Rate Limiting** [1.5 gün]
7. **INT-502-1: ExecutePreHook** [1 gün]

**Sonuç**: Duplicate execution yok, cascading failure engellenir, rate limit uygulanır → **Production Hardened**

---

### Sprint C: Observability & UX — ~3-4 gün
**Hedef**: Operations ve admin teams'i desteklemek

8. **INT-701: Tracing & Metrics** [2 gün]
9. **INT-502-2: Direct Action Test-Run** [1 gün] ← Backend
10. **INT-602-2 + INT-603-2: Test-Run Modals** [1 gün] ← Frontend

**Sonuç**: Metrics/tracing tam, admin UI'da test-run var → **Operational Excellence**

---

### Sprint D: Integrations & Compliance — ~4-5 gün
**Hedef**: Long-lived integrations ve security compliance

11. **INT-702: OAuth2 Token Refresh** [2 gün]
12. **INT-703: Webhook Signature Verification** [1.5 gün]
13. **INT-802: Integration Testing** [2-3 gün] ← Critical for regression prevention

**Sonuç**: OAuth refresh var, webhooks secure, test coverage yüksek → **Enterprise Ready**

---

### Sprint E: Governance & Future-Proof — ~4 gün
**Hedef**: Scaling ve governance hazırlığı

14. **INT-801: Versioning & Approvals** [2 gün]
15. **INT-803: Documentation & Runbooks** [2 gün]
16. **INT-604-3: Execution Batch Retry UI** [0.5 gün] ← Quick win

**Sonuç**: Config versioning var, dokümantasyon tam, runbook'lar hazır → **Scale Ready**

---

## 🚀 Minimum Viable Product (MVP) Definition

**MVP = Sprint A + Sprint B tamamlandığında:**
- ✅ Hooks execute edilebilir (sync + async)
- ✅ HTTP client SSRF-safe
- ✅ Retry/backoff/DLQ tam çalışır
- ✅ Idempotency/circuit breaker/rate limit aktif
- ✅ Admin UI tam (CRUD + list/filter)
- ✅ Basic logging + error handling

**Production Cutover Kriterleri:**
- [ ] End-to-end test: Hook oluştur → event tetikle → execution log'da görün → başarı
- [ ] Error scenario test: 5 consecutive fail → circuit breaker açılsın
- [ ] Rate limit test: Connector limit'i aş → 429 error alsın
- [ ] DLQ test: Max retry aşılınca DLQ'ya düşsün
- [ ] Security test: Private IP'ye istek gönder → reddedilsin
- [ ] Metrics test: Prometheus endpoint'te `integration_execution_total` görünsün

**Tahmini MVP Teslim Süresi**: 10-12 iş günü (2.5 hafta, tek kişi çalışırsa)

---

## 📌 Sonraki Adımlar

### Hemen Yapılacaklar (Bu Hafta)
1. **INT-503-1 başla**: SSRF-safe HTTP client (en kritik güvenlik blocker)
2. **INT-501-2 hazırlık**: JMESPath library research + mapping spec finalize
3. **Backlog refinement**: Sprint A için detaylı task breakdown

### Orta Vadede (2-3 Hafta)
- Sprint A + B tamamla → MVP cutover
- Sprint C başla → Observability + test-run UX

### Uzun Vadede (1-2 Ay)
- Sprint D + E → Enterprise-ready platform
- Production rollout + monitor + iterate
- Community feedback → backlog refine

---

**Doküman Sahibi**: Platform Team  
**Güncelleme Sıklığı**: Her sprint sonunda  
**İletişim**: Platform Slack kanalı #integration-platform
