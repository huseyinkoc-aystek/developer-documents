# API Sharing Implementation Timeline & Sprint Plan

**Timeline:** 6-9 months (MVP)  
**Start Date:** 2026-03-17  
**Target Go-Live:** 2026-10-10  

## Status Guardrail

This timeline mixes planned work and partially implemented prototype items.

- `[x]` means "some repository/code/test work exists in the repo".
- `[ ]` means "not implemented or not verified".
- No sprint in this document should be interpreted as production-approved unless separately evidenced.

Read [0-ARCHITECTURE-REVIEW-AND-REMEDIATION-PLAN.md](0-ARCHITECTURE-REVIEW-AND-REMEDIATION-PLAN.md) first.

## Phase 0: Architecture Correction and Exposure Freeze

**Dates:** 2026-03-11 to 2026-03-24  
**Goal:** stop false go-live assumptions and close P0 design gaps before continuing planned sprints

**Tasks:**
- [ ] Reclassify module status as prototype / internal pilot only
- [ ] Prevent external partner rollout before security remediation
- [x] Fix token issuance to be grant-bounded
- [x] Add app-layer JWT verification or introspection enforcement
- [x] Design client registry and redirect URI whitelist model
- [x] Fix usage aggregation schema and migration plan
- [x] Expose `oauth2_clients` via admin CRUD for controlled onboarding
- [x] Add onboarding workflow: approval state machine, secret rotation, sandbox activation/revocation
- [x] Rewrite sprint status docs to separate "planned", "implemented", "verified"

**Deliverables:**
- Updated architecture review
- P0 remediation backlog
- Exposure/go-live freeze note

---

## Sprint Breakdown (12 sprints = 6 months base + 4-6 weeks for hardening)

### **PHASE 1: Design & Modeling (Sprint 1-2, Weeks 1-4)**

#### Sprint 1: API Operations & Scope Model
**Dates:** 2026-03-17 to 2026-03-31 (2 weeks)  
**Goal:** Foundation for operation registry & scope system

**Tasks:**
- [x] Create models (ApiOperation, Scope, RateLimitPolicy, QuotaPolicy)
- [x] Create repository interfaces
- [x] Create service interfaces
- [x] Design OpenAPI structure for operations
- [x] Implement MongoDB repository for ApiOperation
- [x] Implement MongoDB repository for Scope
- [x] Write unit tests (repo layer — operation + scope repository, 14 tests)
- [x] Write unit tests (service layer — OperationService + ScopeService, 9 tests)
- [ ] API endpoint: `POST /admin/api-sharing/operations` (register operation)
- [ ] API endpoint: `GET /admin/api-sharing/operations` (list)
- [ ] API endpoint: `GET /admin/api-sharing/scopes` (list available scopes)

**Deliverables:**
- Repository implementations + tests
- 2 CRUD endpoints
- RFC 9700 security checklist validation

**Owner:** Backend Team (2-3 engineers)

---

#### Sprint 2: Multi-Tenant Grants & Consents Model
**Dates:** 2026-04-01 to 2026-04-14 (2 weeks)  
**Goal:** Tenant-scoped permissions + KVKK consent tracking

**Tasks:**
- [x] Create models (ApiGrant, ApiConsent, ConsentAuditLog, PartnerDPA)
- [x] Design KVKK compliance mapping
- [x] Create repository interfaces for Grant/Consent
- [x] Implement MongoDB repository for ApiGrant
- [x] Implement MongoDB repository for ApiConsent
- [ ] Implement MongoDB repository for PartnerDPA
- [x] API endpoint: `POST /admin/api-sharing/grants` (create grant)
- [x] API endpoint: `GET /admin/api-sharing/grants` (list by tenant)
- [x] API endpoint: `POST /admin/api-sharing/grants/:id/revoke` (revoke grant)
- [x] API endpoint: `POST /users/me/consents` (create consent - user endpoint)
- [x] API endpoint: `GET /users/me/consents` (list user consents)
- [x] API endpoint: `DELETE /users/me/consents/{id}` (withdraw consent)
- [x] API endpoint: `POST /internal/api-sharing/policy-check` (gateway policy check)
- [x] Write handler-level tests (Sprint 2 flows — 2 test functions covering grants, consents, policy check)
- [ ] Consent form UI design (frontend)
- [ ] Write integration tests

**Deliverables:**
- Grant/Consent repositories + tests
- 5 CRUD endpoints (admin + user)
- KVKK aydınlatma metni template
- DPA template (SCC)

**Owner:** Backend Team (2-3) + Frontend Team (1)

---

### **PHASE 2: Control Plane (Sprint 3-5, Weeks 5-12)**

#### Sprint 3: Policy Engine & Access Control
**Dates:** 2026-04-15 to 2026-04-28 (2 weeks)  
**Goal:** IP allowlist, rate limit, quota enforcement

**Tasks:**
- [x] Design policy evaluation logic
- [x] Create service interfaces (PolicyService)
- [x] Implement IP allowlist check (CIDR parsing)
- [x] Implement scope validator
- [x] Write policy service unit tests (EvaluatePolicy — 3 scenarios)
- [ ] Implement rate limiter (token bucket, Redis-backed)
- [ ] Implement quota checker (daily/monthly tracking)
- [ ] Redis integration for rate limit state
- [ ] API endpoint: `POST /internal/api-sharing/policy-check` (gateway call)
- [ ] Write load tests (429 + Retry-After)
- [ ] Write security tests (IP spoofing, quota overflow)

**Deliverables:**
- PolicyService implementation + Redis integration
- 429 + Retry-After response handling
- Load test results (1000 RPS target)

**Owner:** Backend + Infra Team (2-3)

---

#### Sprint 4: OAuth2/OIDC Authorization Server
**Dates:** 2026-04-29 to 2026-05-12 (2 weeks)  
**Goal:** Token issuance, JWT, JWKS, introspection, revocation (RFC 7519, 7662, 7009, 9700)

**Tasks:**
- [ ] Design token structure (JWT claims, scopes, tenant_id, site_ids)
- [ ] Implement JWT signing (RS256, key rotation strategy)
- [ ] API endpoint: `/oauth/token` (client_credentials flow)
  - [ ] Parse grant_type, client_id, client_secret
  - [ ] Validate client (from api_clients)
  - [ ] Generate JWT token
  - [ ] Handle errors (invalid_client, invalid_grant)
- [ ] API endpoint: `/oauth/authorize` (authorization code flow)
  - [ ] Parse response_type, client_id, redirect_uri, scope, state
  - [ ] Validate client redirect URI
  - [ ] Display consent screen (scopes)
  - [ ] Generate authorization code
- [ ] API endpoint: `/oauth/jwks` (public JWKS endpoint)
- [ ] API endpoint: `/oauth/introspect` (RFC 7662)
  - [ ] Validate token
  - [ ] Return scope, client_id, exp, sub, username
- [ ] API endpoint: `/oauth/revoke` (RFC 7009)
  - [ ] Mark token revoked in cache/DB
  - [ ] Optional: revocation list distribution
- [ ] Implement PKCE support (RFC 7636) for auth code flow
- [ ] Write contract tests (Pact) for OAuth2 flows
- [ ] Write security tests (token expiry, revocation, signature validation)

**Deliverables:**
- 5 OAuth endpoints (token, authorize, jwks, introspect, revoke)
- JWT/JWKS sign+verify
- Contract tests (OAuth2 flows)
- RFC 9700 security validation checklist

**Owner:** Backend Team (3-4 engineers)

---

#### Sprint 5: API Registry & OpenAPI Publishing
**Dates:** 2026-05-13 to 2026-05-26 (2 weeks)  
**Goal:** Dynamic OpenAPI spec generation + Swagger/ReDoc UI

**Tasks:**
- [ ] Design OpenAPI 3.1 structure (paths, operations, security schemes)
- [ ] Implement OpenAPI spec generator (from ApiOperation models)
- [ ] API endpoint: `GET /api-registry/openapi.json` (OpenAPI spec v3.1)
- [ ] API endpoint: `GET /api-registry/swagger-ui` (HTML UI)
- [ ] API endpoint: `GET /api-registry/redoc-ui` (HTML UI)
- [ ] Implement spec versioning (semver)
- [ ] Add security schemes to spec (OAuth2 client_credentials + auth code)
- [ ] Map scopes to operations in spec
- [ ] Write tests (spec validation, security scheme generation)

**Deliverables:**
- OpenAPI spec generation + versioning
- Swagger UI + ReDoc UI (public endpoints)
- OpenAPI 3.1 specification document

**Owner:** Backend Team (2) + Frontend Team (1)

---

### **PHASE 3: Advanced Features (Sprint 6-7, Weeks 13-20)**

#### Sprint 6: Comprehensive Audit & KVKK Compliance
**Dates:** 2026-05-27 to 2026-06-09 (2 weeks)  
**Goal:** Audit logging, retention enforcement, SAR/GDPR

**Tasks:**
- [x] Design audit log schema (PII masking)
- [x] Create audit repository
- [ ] Implement AuditService
- [ ] Log all API accesses (in gateway layer)
  - [ ] tenant_id, client_id, user_id, operation_id, ip, result, latency
  - [ ] Mask PII: hash emails, phone numbers, IP addresses
- [ ] API endpoint: `GET /admin/api-sharing/audit-logs` (query logs)
- [ ] API endpoint: `GET /users/me/data-export` (Subject Access Request - SAR)
  - [ ] Export user consents, granted operations, audit logs as JSON
  - [ ] 7-day generation window
- [ ] Implement retention policies (MongoDB TTL)
  - [ ] Audit logs: 90 days soft-delete → 30 days hard-delete
  - [ ] Consents: 7 years retention → hard-delete
- [ ] Implement deletion/anonymization job
  - [ ] On consent withdrawal: mark audit logs as deleted
  - [ ] Scheduled job: hard-delete old soft-deleted records
- [ ] Implement anonymization logic (replace PII with hashes/nulls)
- [ ] Write KVKK compliance tests (retention verification)

**Deliverables:**
- AuditService + MongoDB TTL cleanup
- SAR endpoint + JSON export
- Deletion/anonymization job (Kubernetes CronJob)
- KVKK compliance report generator

**Owner:** Backend Team (2-3) + Legal Team

---

#### Sprint 7: External Gateway & Policy Enforcement
**Dates:** 2026-06-10 to 2026-06-23 (2 weeks)  
**Goal:** Gateway integration bridge, policy enforcement

**Tasks:**
- [ ] Design gateway architecture (Envoy/Kong/NGINX)
- [ ] Implement gateway policy evaluation middleware
  - [ ] Fetch JWT token
  - [ ] Validate signature via JWKS endpoint
  - [ ] Extract tenant_id, client_id, scopes from token
  - [ ] Call PolicyService (/internal/api-sharing/policy-check)
  - [ ] Check IP allowlist
  - [ ] Check rate limit → if exceeded, return 429 + Retry-After
  - [ ] Check quota
  - [ ] Check scopes vs operation requirements
- [ ] Implement request context enrichment
  - [ ] Add X-Tenant-ID, X-Client-ID, X-User-ID headers
  - [ ] Add W3C Trace-ID for OpenTelemetry
- [ ] Implement audit logging in gateway
  - [ ] Log all API requests with policy decision
  - [ ] Mask sensitive data (PII)
- [ ] Gateway configuration (TLS, WAF, CORS)
- [ ] Write integration tests (gateway + policy)

**Deliverables:**
- Gateway policy middleware (RFC 9700 compliant)
- Request context enrichment
- Audit logging in gateway

**Owner:** Infra + Backend Team (2-3)

---

### **PHASE 4: Observability & Testing (Sprint 8-9, Weeks 21-28)**

#### Sprint 8: Observability & Monitoring
**Dates:** 2026-06-24 to 2026-07-07 (2 weeks)  
**Goal:** OpenTelemetry instrumentation + OTel backend integration

**Tasks:**
- [ ] Add OTel instrumentation to all services
  - [ ] Traces: request path, policy check, database queries
  - [ ] Metrics: requests_total, throttled_total, auth_fail_total, scope_denied_total, latency_p95, p99
  - [ ] Logs: structured JSON logs, PII masking
- [ ] Export to OTel backend (Jaeger, Prometheus, Loki)
- [ ] Create Grafana dashboards
  - [ ] Overview: RPS, error rate, latency, throttled count
  - [ ] Per-tenant: usage, scope violations, IP blocks
  - [ ] Per-operation: latency patterns, error distribution
- [ ] Create alert rules
  - [ ] Quota exceeded (per tenant)
  - [ ] High error rate (>1%)
  - [ ] Auth failures (>10/min)
  - [ ] Scope denials (>100/day)
- [ ] Document oncall runbooks

**Deliverables:**
- OTel instrumentation + OTLP export
- Grafana dashboards (5+)
- Alert rules (PagerDuty/AlertManager)
- Runbook documentation

**Owner:** Infra + Backend Team (2)

---

#### Sprint 9: Testing & Security Validation
**Dates:** 2026-07-08 to 2026-07-21 (2 weeks)  
**Goal:** Contract, integration, security, & chaos tests

**Tasks:**
- [ ] Write contract tests (Pact)
  - [ ] OAuth2 token flow (client_credentials, auth code + PKCE)
  - [ ] Policy check API
  - [ ] Gateway policy enforcement
- [ ] Write integration tests
  - [ ] End-to-end: create grant → generate token → call API → audit log
  - [ ] Rate limit: 429 + Retry-After behavior
  - [ ] IP allowlist: block unauthorized IPs
  - [ ] Scope validation: deny missing scopes
  - [ ] TTL: token expiry + refresh
- [ ] Write security tests (OWASP API Top 10)
  - [ ] BOLA (Broken Object Level Authorization): cross-tenant access
  - [ ] BFLA (Broken Function Level Authorization): operation-level checks
  - [ ] BAPI1: excessive data exposure (PII in responses)
  - [ ] BAPI3: unrestricted resource consumption (quota bypass)
  - [ ] BAPI4: lack of rate limiting (429 generation)
- [ ] Chaos tests
  - [ ] Kill gateway pod → failover
  - [ ] Auth server latency → request timeout
  - [ ] Network partition → timeout + retry
  - [ ] Redis outage → fallback to in-memory rate limit
- [ ] KVKK compliance test suite
  - [ ] Consent withdrawal → audit logs deleted
  - [ ] TTL enforcement → old logs deleted
  - [ ] SAR functionality → all user data exported
  - [ ] Anonymization → PII masked/hashed

**Deliverables:**
- Contract tests (Pact, >50 scenarios)
- Integration test suite (>100 tests)
- Security test suite (OWASP API Top 10)
- Chaos test reports
- KVKK compliance test results

**Owner:** QA + Security Team (3-4)

---

### **PHASE 5: Hardening & Go-Live (Sprint 10-12, Weeks 29-36)**

#### Sprint 10: Pilot & Partner Onboarding
**Dates:** 2026-07-22 to 2026-08-04 (2 weeks)  
**Goal:** Sandbox environment + test partner integration

**Tasks:**
- [ ] Set up sandbox environment (separate databases, Redis)
- [ ] Onboard 1-2 test partners
  - [ ] Create API clients
  - [ ] Generate credentials (client_id, client_secret)
  - [ ] Grant scopes (orders.read, shipments.write)
  - [ ] Provide OpenAPI docs + code examples
- [ ] Monitor pilot metrics
  - [ ] API success rate
  - [ ] Latency patterns
  - [ ] Rate limit hit rate
- [ ] Gather feedback + bug fixes
- [ ] SRE runbooks for common issues
  - [ ] Token revocation
  - [ ] Quota reset
  - [ ] IP allowlist updates
  - [ ] Incident response (auth failure, rate limiter gone wild)

**Deliverables:**
- Sandbox environment fully operational
- 1-2 partners successfully integrated
- SRE runbooks + on-call playbooks
- Known issues log with mitigation

**Owner:** Backend + Infra + SRE Team (3-4)

---

#### Sprint 11: Load Testing & Performance Tuning
**Dates:** 2026-08-05 to 2026-08-18 (2 weeks)  
**Goal:** Performance validation + capacity planning

**Sprint 11 Progress (% Complete):** **72%**

**Progress Model (weighted):**

| Workstream | Weight | Current Status | Completion |
|---|---:|---|---:|
| Load execution + SLO validation | 45% | representative local run completed, 3/5 gates pass, throughput near target, staging still pending | 65% |
| Bottleneck analysis | 20% | protected middleware-chain bottleneck identified and fixed; Redis/JWT staging profile still pending | 85% |
| Performance tuning rollout | 20% | scope/operation cache, internal route correction, and load helper hardening completed; infra tuning pending | 70% |
| Capacity and cost planning | 15% | measured baseline folded into sizing and cost model; Infra/SRE approval pending | 80% |

**Computed Overall:** `0.45*65 + 0.20*85 + 0.20*70 + 0.15*80 = 72.25%`

**Current Metric Gates (corrected local baseline):**

| Metric Gate | Target | Current | Status |
|---|---:|---:|---:|
| Sustained request rate | >=1000 req/s | 987.12 req/s | near miss |
| p95 latency | <500 ms | 302.39 ms | pass |
| p99 latency | <1000 ms | 1.34 s | fail |
| Error rate | <0.1% | 0.00% | pass |
| Throttle rate | <1% | 0.00% | pass |

Sprint 11 now has real execution evidence and a verified root-cause fix. It remains below hardening-ready status because staging evidence is still missing and the `p99 < 1000ms` gate is still open. Passing the same profile in staging with `p99` brought below `1.0s` moves Sprint 11 to `85%+`.

**Deliverables:**
- Updated load/performance report with pre-fix and post-fix evidence: `SPRINT_11_LOAD_TEST_PERFORMANCE_REPORT.md`
- Performance optimization checklist with route/middleware correction (partial)
- Updated capacity planning document: `SPRINT_11_CAPACITY_PLANNING.md`
- HPA + resource request updates (recommended, not yet validated in staging)

**Owner:** Infra Team (2-3)

---

#### Sprint 12: Hardening & Go-Live Readiness
**Dates:** 2026-08-19 to 2026-09-01 (2 weeks)  
**Goal:** Final validation + production deployment

**Tasks:**
- [x] Production readiness checklist (draft)
  - [ ] All RFC requirements covered
  - [ ] KVKK compliance verified
  - [ ] Security test coverage (>90%)
  - [ ] Documentation complete
  - [ ] Runbooks validated
- [ ] Staging environment full system test
  - [ ] Production data volume
  - [ ] Real partner scenarios
  - [ ] Disaster recovery drill (backup/restore)
- [ ] SRE sign-off
  - [ ] Monitoring + alerting operational
  - [ ] Incident response procedures tested
  - [ ] On-call rotation established
- [ ] Final security audit (internal or external)
- [ ] Legal/Compliance sign-off (KVKK)
- [ ] Product owner approval
- [x] Deployment runbook (draft)
  - [x] Blue-green deployment strategy
  - [x] Rollback plan
  - [x] Health checks

**Deliverables:**
- Production readiness checklist draft: `SPRINT_12_PRODUCTION_READINESS_CHECKLIST.md`
- Staging test results (all pass) (pending)
- Security audit report + sign-off (pending)
- KVKK compliance approval (pending)
- Deployment runbook draft: `SPRINT_12_DEPLOYMENT_RUNBOOK.md`

**Owner:** All Teams + Leadership

---

## Contingency Buffer (2-4 weeks)

If issues discovered in load testing or security audit:
- Weeks 37-40: Extended hardening/bug fixing
- Target Go-Live: 2026-10-31 (worst case)

---

## Parallel Workstreams (can overlap phases)

### Documentation (Throughout)
- OpenAPI specification (publish in Sprint 5)
- API client integration guide (publish in Sprint 6)
- KVKK compliance guide (publish in Sprint 6)
- Runbooks & on-call guides (publish in Sprint 10)

### Training (Weeks 32-34)
- Partner onboarding flow training
- SRE incident response drills
- Admin panel training (for operations team)

### Release Management (Week 35+)
- Release notes preparation
- Beta program setup (if applicable)
- Communications (blog, whitepaper)

---

## Success Criteria

| Milestone | Target | Current |
|---|---|---|
| RFC 9700 compliance (OAuth2) | 100% | 0% |
| KVKK compliance (audit, retention) | 100% | 0% |
| API test coverage | >90% | 0% |
| Load test: 1000 RPS | p95 <500ms | Local corrected run: 987.12 req/s, p95 302.39ms, p99 1.34s, staging pending |
| Security tests (OWASP) | >80% | 0% |
| Partner pilot | 2 partners | 0% |
| Documentation | 100% complete | 25% (Sprint 11/12 docs added) |
| Go-Live readiness | All sign-offs | In progress |

---

**Document Version:** 1.1  
**Next Review:** 2026-04-01 (end of Phase 1)  
**Prepared by:** Architecture Team  
**Last Updated:** 2026-03-11
