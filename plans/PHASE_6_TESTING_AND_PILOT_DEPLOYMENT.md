---
title: Phase 6 - Testing, Hardening & Go-Live
version: 1.0.0
status: TESTED - PILOT EXECUTION PENDING
last_updated: 2026-03-09
classification: Internal - Implementation Plan
related_documents:
  - ./AYSTEK_CORE_API_GATEWAY_THIRD_PARTY_ACCESS_CONTROL.md
  - ./PHASE_6_INTERNAL_PILOT_RUNBOOK.md
  - ./PHASE_6_PARTNER_ONBOARDING_CHECKLIST.md
  - ../user-guide/05-operations/API_GATEWAY_SECURITY_AUDIT_CHECKLIST.md
  - ../user-guide/04-api/API_CLIENT_ONBOARDING_BEST_PRACTICES.md
---

# Phase 6: Testing, Hardening & Go-Live (1-2 Weeks)

## Executive Summary

Phase 6 focuses on **comprehensive test coverage**, **infrastructure validation**, and **controlled pilot deployment** with internal team and first partner. All test suites validate:

1. **Functional correctness**: Multi-tenant isolation, concurrent access, permission enforcement
2. **Performance**: Rate-limit accuracy, quota fairness, latency under load
3. **Security**: Secret rotation, revocation, IP bypass detection, token expiry
4. **Operational readiness**: Monitoring validation, alert testing, rollback procedures

## Progress Update (2026-03-09)

- [x] Phase 6 test fixture stack was rebuilt (`go/aystek-core/tests/fixtures.go`) and the suite now compiles/runs end-to-end.
- [x] Validation run completed successfully: `go test ./tests/...`.
- [x] Internal pilot and partner onboarding runbooks are available:
  - `plans/PHASE_6_INTERNAL_PILOT_RUNBOOK.md`
  - `plans/PHASE_6_PARTNER_ONBOARDING_CHECKLIST.md`
- [ ] Production pilot execution and first partner live traffic window remain pending scheduling.

---

## Testing Strategy

### 1. Integration Tests (Location: `go/aystek-core/tests/api_gateway_integration_test.go`)

**Objective:** Validate OAuth2 flow, multi-tenant isolation, scope enforcement, and authorization logic in integrated environment.

#### Test Categories

**1.1 OAuth2 Token Flow**
- [ ] Valid `client_credentials` grant returns `access_token` + `refresh_token`
- [ ] `refresh_token` rotates on reuse (sliding window validation)
- [ ] Token claims contain correct `client_id`, `scope`, `tenant_id`, `site_id`
- [ ] Expired token rejected with `401 Unauthorized`
- [ ] Invalid `client_secret` rejected with `401 Unauthorized`
- [ ] Revoked API client cannot obtain new token (immediate revocation)
- [ ] Token scope reduced if API client scope downgraded mid-session (reclaim on next issue)

**1.2 Multi-Tenant Isolation**
- [ ] API client scoped to tenant A cannot access tenant B routes (403 Forbidden)
- [ ] Audit logs show correct `tenant_id` and `site_id` in context key verification
- [ ] Rate limit quotas isolated per tenant (tenant A's 1000 req/hour ≠ tenant B's)
- [ ] Two concurrent API clients from same tenant share quota fairly
- [ ] Concurrent requests from different tenants execute without cross-tenant data leak

**1.3 Route Authorization**
- [ ] API client with `products:read` scope can `GET /products`, cannot `POST /products`
- [ ] API client with `orders:create` scope can `POST /orders`, cannot `DELETE /orders`
- [ ] Route whitelist enforced: `GET /admin` denied unless explicitly whitelisted
- [ ] Scope binding to method: `read:orders` → GET only, not POST/DELETE
- [ ] Scope inflation attempt rejected (no privilege escalation)

**1.4 Concurrent Requests**
- [ ] 100 simultaneous requests from same API client succeed without race condition
- [ ] Quota tracking accurate with 100 concurrent requests (not duplicated)
- [ ] Request ordering preserved (no request interleaving)
- [ ] Audit logs record all 100 requests, no missing entries
- [ ] Response latency consistent (no queuing timeout)

#### Test Implementation Pattern

```go
func TestMultiTenantIsolation(t *testing.T) {
    // Setup: Create two API clients in different tenants
    tenantA := setupTenant(t, "tenant-a")
    tenantB := setupTenant(t, "tenant-b")
    
    clientA := createAPIClient(t, tenantA)
    clientB := createAPIClient(t, tenantB)
    
    // Act: Client A tries to access tenant B product endpoint
    tokenA := obtainToken(t, clientA)
    resp, err := callAPIWithToken(t, 
        "GET /products?tenant_id="+tenantB.ID, tokenA)
    
    // Assert: Should be denied
    if err != ErrForbidden {
        t.Fatalf("expected 403, got %v", err)
    }
}
```

---

### 2. Load Tests (Location: `go/aystek-core/tests/api_gateway_load_test.go`)

**Objective:** Validate rate-limit accuracy, quota fairness, and latency consistency under sustained high throughput.

#### Load Test Scenarios

**2.1 Rate-Limit Accuracy**
- [ ] Apply rate limit: 100 req/min per API client
- [ ] Send 120 requests in 60 seconds
- [ ] First 100 succeed with 200 OK
- [ ] Request 101-120 rejected with `429 Too Many Requests`
- [ ] First request after 60s allowed (quota resets)
- [ ] Validate atomic quota tracking (no double-counting)

**2.2 Quota Fairness (Multiple API Clients)**
- [ ] Create 3 API clients, each with 1000 req/hour quota
- [ ] Client A sends 500 req/hour, Client B sends 400 req/hour, Client C sends 300 req/hour
- [ ] All succeed (total 1200 < 3000 available)
- [ ] Client A sends 600 req/hour (exceeds its 1000 limit)
- [ ] Client A's next request returns 429; Clients B/C unaffected
- [ ] Validate per-client quota tracking (not shared pool)

**2.3 Latency Consistency**
- [ ] Sustained 1000 req/min for 10 minutes
- [ ] p50 latency < 100ms
- [ ] p95 latency < 500ms
- [ ] p99 latency < 1000ms
- [ ] No latency degradation after quota resets (no cleanup stalls)

**2.4 Concurrent Tenant Load**
- [ ] 10 tenants × 5 API clients each = 50 concurrent clients
- [ ] Each sends 100 req/min (total 5000 req/min cluster-wide)
- [ ] All requests succeed within fairness window
- [ ] Tenant A's spike (1000 req/min) only affects tenant A's quota, not B-J
- [ ] Audit logs complete and accurate across all 50 clients

#### Load Test Implementation Pattern

```go
func BenchmarkRateLimitAccuracy(b *testing.B) {
    client := createAPIClient(b, "test-tenant")
    token := obtainToken(b, client)
    
    limit := 100 // requests per minute
    successCount := 0
    rateLimitedCount := 0
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        resp, _ := callAPIWithToken(b, "GET /products", token)
        if resp.StatusCode == 200 {
            successCount++
        } else if resp.StatusCode == 429 {
            rateLimitedCount++
        }
    }
    
    if successCount != limit {
        b.Fatalf("expected %d successful, got %d", limit, successCount)
    }
}
```

---

### 3. Security Tests (Location: `go/aystek-core/tests/api_gateway_security_test.go`)

**Objective:** Validate secret rotation, revocation, IP whitelist, and token expiry enforcement.

#### Security Test Scenarios

**3.1 Secret Rotation**
- [ ] Create API client with `secret_v1`
- [ ] Rotate to `secret_v2`: Old token still valid during grace period (5 min)
- [ ] After grace period, `secret_v1` rejected (only `secret_v2` works)
- [ ] Rotate back to `secret_v1`: New tokens issued with `secret_v1`
- [ ] Audit log records all rotation events with timestamp and operator

**3.2 Revocation**
- [ ] Create API client, obtain active token
- [ ] Revoke API client
- [ ] Existing token rejected with `401 Unauthorized` (immediate)
- [ ] New token request fails (client no longer exists)
- [ ] Audit log shows revocation event with operator and reason

**3.3 IP Whitelist Bypass Prevention**
- [ ] Create API client with IP whitelist: `10.0.0.0/8`, `192.168.1.0/24`
- [ ] Request from 10.1.1.1 succeeds (in CIDR range)
- [ ] Request from 192.168.1.100 succeeds (in CIDR range)
- [ ] Request from 172.16.0.1 rejected with `403 Forbidden` (outside CIDR range)
- [ ] Request with spoofed `X-Forwarded-For: 10.0.0.1` rejected (only `X-Real-IP` or socket address trusted)
- [ ] Request with crafted URL encoding (e.g., `10.0.0.%31`) rejected

**3.4 Token Expiry**
- [ ] Access token issued with 15 min expiry
- [ ] After 15 min, token rejected with `401 Unauthorized`
- [ ] Refresh token issued with 7 day expiry
- [ ] After 7 days, refresh token rejected, full re-auth required

**3.5 Scope Isolation**
- [ ] API client with `products:read` scope cannot access `orders:write` endpoint
- [ ] Scope claim in JWT verified against route requirements
- [ ] No scope inflation via token manipulation (HMAC signature verified)

#### Security Test Implementation Pattern

```go
func TestSecretRotation(t *testing.T) {
    client := createAPIClient(t, "test-tenant")
    secretV1 := client.Secret
    
    // Obtain token with secret_v1
    tokenV1 := obtainToken(t, client)
    
    // Rotate secret
    secretV2, _ := rotateSecret(t, client)
    
    // Old token still valid during grace period
    resp, _ := callAPIWithToken(t, "GET /products", tokenV1)
    if resp.StatusCode != 200 {
        t.Fatalf("expected old token valid during grace period")
    }
    
    // After grace period, old token rejected
    time.Sleep(6 * time.Minute) // grace period + 1min
    resp, _ = callAPIWithToken(t, "GET /products", tokenV1)
    if resp.StatusCode != 401 {
        t.Fatalf("expected old token rejected after grace period")
    }
    
    // New token works with secret_v2
    tokenV2 := obtainTokenWithSecret(t, client, secretV2)
    resp, _ = callAPIWithToken(t, "GET /products", tokenV2)
    if resp.StatusCode != 200 {
        t.Fatalf("expected new token valid")
    }
}
```

---

## Pilot Deployment Plan

### Phase 6a: Internal Team Pilot (Week 1)

**Objective:** Validate API Gateway in controlled production-like environment with AYSTEK engineering team.

#### Participants
- API Gateway team (2-3 engineers)
- Backend integration team (2 engineers)
- Frontend team (1 engineer)
- DevOps/SRE (1 engineer)

#### Schedule
- **Day 1–2**: Internal API client provisioning, secret setup, scope configuration
- **Day 3–5**: Live traffic validation, monitoring review, incident response drill
- **Day 6–7**: Feedback collection, bug fixes, documentation updates

#### Internal Pilot Checklist

- [ ] **Setup**
  - [ ] Provision 5 API clients (one per team)
  - [ ] Configure scopes: `products:read`, `orders:read`, `orders:create`, `users:read`
  - [ ] Configure IP whitelist: `10.0.0.0/8` (internal networks)
  - [ ] Set rate limit: 10,000 req/hour per client
  - [ ] Set quota: 100,000 req/day per client

- [ ] **Integration**
  - [ ] Each team integrates API client into their service
  - [ ] Token refresh logic validated (15 min token TTL, 7 day refresh)
  - [ ] Error handling verified (401/403/429 responses)
  - [ ] Audit log export tested (`GET /admin/api-clients/:id/audit-logs/export?format=csv`)

- [ ] **Monitoring**
  - [ ] Grafana dashboard reviewed (request rate, error rate, latency)
  - [ ] Prometheus alerts tested (spike in 5xx, latency spike, auth failures)
  - [ ] Alert routing verified (Slack, PagerDuty, email)

- [ ] **Hardening**
  - [ ] Secret rotation tested (90-day policy, immediate post-incident)
  - [ ] Revocation tested (immediate token invalidation)
  - [ ] IP whitelist bypass prevention validated
  - [ ] CSV export formula injection prevention verified

- [ ] **Feedback**
  - [ ] Team collects issues/feature requests
  - [ ] Documentation review (API Onboarding Guide, Best Practices)
  - [ ] Operational runbook validation (Secret Rotation, SLO/SLA)

#### Success Criteria (Internal Pilot)
- ✅ 100 consecutive hours without incident
- ✅ All team members able to obtain token and call scoped endpoints
- ✅ Monitoring dashboard reflects all traffic accurately
- ✅ No unexpected 401/403/429 responses
- ✅ Audit logs complete and queryable
- ✅ Team confidence: 4/5 stars or higher

---

### Phase 6b: Partner Onboarding (Week 2)

**Objective:** Successfully onboard first third-party API client and validate end-to-end integration, audit, and support workflows.

#### Partner Selection Criteria
- Existing AYSTEK customer (reduces trust barrier)
- Limited scope initially (1–3 endpoints)
- Dedicated technical contact (for support/feedback loop)
- Use case: mobile app integration, ERP synchronization, or dashboard

#### Partner Onboarding Workflow

1. **Planning Phase (Day 1–2)**
   - [ ] Partner identifies required endpoints (e.g., `GET /products`, `POST /orders`)
   - [ ] Partner provides IP whitelist (CIDR blocks, data center IPs, VPN gateways)
   - [ ] Partner commits to secret rotation cadence (90 days or per incident)
   - [ ] Draft KVKK data processing agreement if applicable

2. **API Client Setup (Day 3)**
   - [ ] AYSTEK admin creates API client in dashboard
   - [ ] Scope set: `products:read`, `orders:create` (minimum)
   - [ ] Route whitelist: Only partner-approved endpoints
   - [ ] Rate limit: 1,000 req/hour (initial, tuned per feedback)
   - [ ] Quota: 50,000 req/day
   - [ ] IP whitelist: Partner-provided CIDR blocks
   - [ ] Generate secret (`client_id` + `client_secret`)

3. **Partner Integration (Day 4–5)**
   - [ ] Provide API client credentials via secure channel (encrypted, one-time delivery)
   - [ ] Partner implements OAuth2 Client Credentials flow:
     ```bash
     curl -X POST https://api.aystek.com/auth/oauth2/token \
       -H "Content-Type: application/x-www-form-urlencoded" \
       -d "grant_type=client_credentials&
            client_id=p_acme_001&
            client_secret=<secret>&
            scope=products:read,orders:create"
     ```
   - [ ] Partner tests token caching (e.g., refresh at 14 min before expiry)
   - [ ] Partner validates error handling (401/403/429 responses)
   - [ ] Partner runs smoke tests on all scoped endpoints

4. **UAT & Audit (Day 6–7)**
   - [ ] Partner executes 100 production-like requests per endpoint
   - [ ] AYSTEK monitors dashboard for anomalies
   - [ ] Partner exports audit logs for compliance validation
   - [ ] Verify audit log records: `client_id`, `method`, `path`, `status`, `timestamp`, `tenant_id`, `site_id`
   - [ ] Partner confirms secret rotation capability (test rotate, reissue token, validate)

5. **Go-Live & Support (Week 2, Day 5+)**
   - [ ] Partner switches traffic from legacy integration to OAuth2 gateway
   - [ ] On-call support available (engineering Slack channel)
   - [ ] Partner monitors error rates, latency, quota usage
   - [ ] AYSTEK validates partner quota fairness (no cross-tenant impact)
   - [ ] Bid-weekly check-in for 1 month (feedback, performance tuning)

#### Partner Onboarding Checklist

- [ ] **Pre-Onboarding**
  - [ ] Partner use case documented
  - [ ] Data sensitivity classified (PII/non-PII)
  - [ ] IP whitelist provided
  - [ ] KVKK agreement (if applicable) signed
  - [ ] Support contact assigned

- [ ] **API Client Setup**
  - [ ] Client created in admin UI
  - [ ] Scope/route whitelist configured and partner-approved
  - [ ] Rate limit/quota thresholds set
  - [ ] IP whitelist validated against partner infrastructure
  - [ ] Secret securely delivered (encrypted, one-time link with 24h expiry)

- [ ] **Integration**
  - [ ] Partner token request succeeds
  - [ ] All scoped endpoints return expected data
  - [ ] Out-of-scope endpoint returns 403 Forbidden
  - [ ] IP whitelist blocks unauthorized IPs (403 Forbidden)
  - [ ] Token expiry handled gracefully (refresh token flows)

- [ ] **Audit & Compliance**
  - [ ] Audit logs complete (all requests recorded)
  - [ ] CSV export validated by partner
  - [ ] KVKK data classification correct
  - [ ] Compliance evidence collected (audit log snapshots)

- [ ] **Documentation**
  - [ ] Partner onboarding guide reviewed (`API_CLIENT_ONBOARDING_BEST_PRACTICES.md`)
  - [ ] Integration guide provided (secret rotation, error handling, best practices)
  - [ ] Support escalation procedure communicated
  - [ ] Incident notification procedure documented

- [ ] **Go-Live**
  - [ ] Partner production traffic validated for 24 hours
  - [ ] No quota fairness issues observed
  - [ ] Alert tuning complete (no noise)
  - [ ] Partner satisfied with monitoring/audit capabilities
  - [ ] Bid-weekly check-in scheduled for 1 month

#### Partner Communication Template

```markdown
## AYSTEK API Gateway — Partner Integration Guide

### API Client Credentials
- **Client ID**: p_acme_001
- **Secret**: [Encrypted, one-time delivery link]
- **Base URL**: https://api.aystek.com
- **Scope**: products:read, orders:create
- **Rate Limit**: 1,000 req/hour
- **Quota**: 50,000 req/day

### Quick Start
1. Obtain access token: POST /auth/oauth2/token
2. Call API endpoint: GET /products (with Authorization: Bearer <token>)
3. Handle token expiry: Refresh at 14 min before 15 min expiry
4. Monitor quota: Check rate-limit headers (X-RateLimit-Remaining, X-RateLimit-Reset)

### Support
- Slack: #aystek-api-support (24/7)
- Email: api-support@aystek.com
- Incident escalation: pagerduty link

### Secret Rotation
- Schedule: Every 90 days (next rotation: June 9, 2026)
- Process: POST /admin/api-clients/p_acme_001/rotate-secret
- Grace period: Old secret valid for 5 min after rotation
- Testing: Test new secret before decommissioning old one

### Audit & Compliance
- Export logs: GET /admin/api-clients/p_acme_001/audit-logs/export?format=csv
- Log retention: 90 days
- Compliance: KVKK data processor agreement signed and on file
```

---

## Testing Infrastructure Setup

### Test Fixtures & Mocks

**Location:** `go/aystek-core/tests/fixtures.go`

```go
// CreateTestTenant sets up a tenant with site, user, and OAuth2 config
func CreateTestTenant(t *testing.T, name string) *models.Tenant { ... }

// CreateTestAPIClient provisions an API client with default scopes/limits
func CreateTestAPIClient(t *testing.T, tenant *models.Tenant) *models.APIClient { ... }

// ObtainTestToken issues an OAuth2 access token for testing
func ObtainTestToken(t *testing.T, client *models.APIClient) string { ... }

// MockPrometheusRegistry provides in-memory metrics for testing
func MockPrometheusRegistry(t *testing.T) prometheus.Registerer { ... }

// MockMongoCollection mocks MongoDB for audit log tests
func MockMongoCollection(t *testing.T) *mongo.Collection { ... }
```

### CI/CD Integration

**GitHub Actions Workflow:** `.github/workflows/phase6-tests.yml`

```yaml
name: Phase 6 - Integration & Security Tests

on: [push, pull_request]

jobs:
  integration-tests:
    runs-on: ubuntu-latest
    services:
      mongodb:
        image: mongo:6.0
        options: --health-cmd="mongosh --eval 'db.adminCommand(\"ping\")'" --health-interval 10s --health-timeout 5s --health-retries 5
      redis:
        image: redis:7.0
        options: --health-cmd "redis-cli ping" --health-interval 10s --health-timeout 5s --health-retries 5
    
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      
      - name: Run Integration Tests
        run: go test -v -race ./tests/api_gateway_integration_test.go
      
      - name: Run Security Tests
        run: go test -v ./tests/api_gateway_security_test.go
      
      - name: Run Load Tests (Quick)
        run: go test -benchtime=10s -bench=. ./tests/api_gateway_load_test.go
```

---

## Phase 6 Completion Checklist

### Testing ✅/❌
- [ ] Integration tests pass (all 15 test cases)
- [ ] Load tests pass (rate limit, quota fairness, latency SLA)
- [ ] Security tests pass (rotation, revocation, IP bypass, expiry)
- [ ] Test coverage: >80% on auth + API client modules
- [ ] CI/CD pipeline green (all checks pass)

### Internal Pilot ✅/❌
- [ ] 5 API clients provisioned for internal teams
- [ ] 100 consecutive hours without incident
- [ ] Team feedback positive (4/5 stars or higher)
- [ ] Documentation reviewed and updated
- [ ] Monitoring alerts tuned and validated

### Partner Onboarding ✅/❌
- [ ] First partner API client successfully provisioned
- [ ] Partner integration tested and validated
- [ ] Audit logs complete and partner-reviewed
- [ ] KVKK/compliance agreement signed
- [ ] Partner data classification documented
- [ ] Go-live successful (24 hour production validation)

### Operational Readiness ✅/❌
- [ ] Runbooks validated (secret rotation, SLO/SLA, security audit)
- [ ] Monitoring dashboard reflects all metrics accurately
- [ ] Alert rules tested and configured (no noise)
- [ ] Incident response procedure documented and drilled
- [ ] Rollback procedure documented and tested
- [ ] Support escalation clear (on-call, Slack, PagerDuty)

### Documentation ✅/❌
- [ ] Phase 6 completion report written
- [ ] Go-live readiness assessed as "Ready for Production" ✅
- [ ] Next steps for Phase 7 (ongoing operations) documented
- [ ] Lessons learned captured
- [ ] Knowledge transfer completed to operations team

---

## Risks & Mitigation

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|-----------|
| Secret rotation grace period too short (5 min) | Medium | High | Test with partner, adjust to 15 min if needed |
| Rate limit quota tracking inaccuracy under load | Low | High | Use Redis atomic counters + testing with 1000s concurrent requests |
| IP whitelist bypass via X-Forwarded-For spoofing | Medium | High | Only trust socket address (requires reverse proxy validation) |
| Partner data exposure via audit export CSV | Low | Critical | Implement data masking for PII fields in export |
| Monitoring alert fatigue during internal pilot | High | Medium | Baseline thresholds from load tests, tune weekly |
| Partner goes live without proper error handling | High | Medium | Provide sample error handling code in onboarding guide |
| Rollback procedure too slow (>10 min) | Medium | High | Test rollback end-to-end, document DNS migration procedure |

---

## Success Metrics (Phase 6 Exit Criteria)

| Metric | Target | Status |
|--------|--------|--------|
| Test coverage (auth + API client modules) | >80% | ⏳ Pending |
| Integration test pass rate | 100% | ⏳ Pending |
| Load test RPS accuracy | ±2% | ⏳ Pending |
| Internal pilot incidents | <1 | ⏳ Pending |
| Internal pilot team satisfaction | ≥4/5 stars | ⏳ Pending |
| Partner integration time | <7 days | ⏳ Pending |
| Partner go-live success | First time right | ⏳ Pending |
| Production SLA achievement | ≥99.9% | ⏳ Pending |

---

## Go-Live Decision Gate

**Phase 6 → Phase 7 (Production)** transition requires:

1. ✅ All integration/security/load tests pass
2. ✅ Internal pilot 100-hour incident-free window
3. ✅ First partner successfully integrated and validated
4. ✅ Monitoring/alerting tuned and operational
5. ✅ Runbooks executed and team trained
6. ✅ Rollback procedure tested and <5 min RTO

**Approval:** API Gateway team lead + DevOps lead + Security lead

**Post-Go-Live:** 24/7 on-call, daily SLA review, weekly operational review

---

## Timeline (1-2 Weeks Estimated)

```
Week 1 (Days 1-5):
  Day 1-2: Test implementation + infrastructure setup
  Day 3-4: Internal pilot deployment + initial validation
  Day 5: Internal pilot feedback, bug fixes

Week 2 (Days 6-10):
  Day 6-7: Partner selection + onboarding prep
  Day 8-9: Partner integration + UAT
  Day 10: Go-live decision + Phase 7 handoff
```

**Fast-track option (if all tests pass quickly):**
- Compress to 7-10 days by running internal pilot + partner onboarding in parallel
- Requires dedicated staff and 24/7 support availability

---

**Last Updated:** March 9, 2026  
**Owner:** API Gateway Team Lead  
**Approval:** [Pending sign-off]
