---
title: API Gateway — Internal Pilot Rollout Runbook
version: 1.0.0
date: 2026-03-09
audience: AYSTEK Engineering Team, DevOps
classification: Internal - Operational
---

# API Gateway — Internal Pilot Rollout Runbook

## Quick Reference

| Phase | Timeline | Participants | Objective |
|-------|----------|--------------|-----------|
| **Setup** | Day 1-2 | DevOps, Product | Provision 5 API clients for internal teams |
| **Integration** | Day 3-5 | Backend, Frontend, Integrations | Teams integrate API client into services |
| **Validation** | Day 6-7 | QA, Monitoring | 100-hour incident-free window |
| **Feedback & Fixes** | Day 7-8 | All | Collect feedback, patch critical issues |

---

## Pre-Pilot Checklist (Day 0)

### Infrastructure Readiness
- [ ] API Gateway deployed to staging environment
- [ ] MongoDB audit log collection running (`api_audit_logs`)
- [ ] Redis rate-limit counters operational
- [ ] Prometheus metrics endpoint accessible
- [ ] Grafana dashboard provisioned and accessible
- [ ] Alert rules configured and tested (no false positives)

### Team Preparation
- [ ] On-call schedule established (24/7 coverage Week 1)
- [ ] Incident response procedure documented
- [ ] Support Slack channel created (#aystek-api-gateway-pilot)
- [ ] Escalation contacts identified and verified
- [ ] Team training completed (OAuth2 flow, error handling, API client UI)

### Documentation
- [ ] API Onboarding Guide reviewed by team
- [ ] Integration Guide distributed
- [ ] Runbooks accessible (secret rotation, SLO/SLA, security audit)
- [ ] Sample code (Go, JavaScript, Python) prepared

---

## Day 1-2: Setup Phase

### Activity 1.1: Create API Clients (30 min)

**Responsible:** Product Manager + DevOps

**Steps:**

1. Access Admin Dashboard: `https://api.aystek.dev/admin/api-clients`

2. Create 5 API clients (one per team):

   ```
   Team 1: Backend Integration (Services)
   - Name: backend-integration-internal
   - Scope: products:read, orders:read, orders:create, users:read
   - Rate Limit: 10,000 req/hour
   - Quota: 100,000 req/day
   - IP Whitelist: 10.0.0.0/8
   
   Team 2: Frontend (Dashboard, Portal)
   - Name: frontend-dashboard-internal
   - Scope: products:read, orders:read, users:read
   - Rate Limit: 5,000 req/hour
   - Quota: 50,000 req/day
   - IP Whitelist: 10.0.0.0/8
   
   Team 3: Integrations (3rd-party sync)
   - Name: integrations-sync-internal
   - Scope: products:read, products:write, orders:read, orders:write
   - Rate Limit: 3,000 req/hour
   - Quota: 30,000 req/day
   - IP Whitelist: 10.0.0.0/8
   
   Team 4: QA (Testing, Validation)
   - Name: qa-testing-internal
   - Scope: products:read, orders:read, orders:create
   - Rate Limit: 2,000 req/hour
   - Quota: 20,000 req/day
   - IP Whitelist: 10.0.0.0/8
   
   Team 5: Mobile (iOS/Android)
   - Name: mobile-app-internal
   - Scope: products:read, orders:read
   - Rate Limit: 1,000 req/hour (conservative to prevent dev testing impact)
   - Quota: 10,000 req/day
   - IP Whitelist: 10.0.0.0/8
   ```

3. Record credentials securely:
   ```
   Store in 1Password vault: api-gateway-pilot-credentials
   Share client_id + client_secret via encrypted 1Password link (24h expiry)
   ```

4. Verify each client in list view

**Acceptance:** All 5 clients created and visible in admin UI

---

### Activity 1.2: Configure Monitoring (1 hour)

**Responsible:** DevOps + Monitoring

**Steps:**

1. Verify Prometheus scrape target:
   ```bash
   curl -s http://prometheus:9090/api/v1/targets | jq '.data.activeTargets[] | select(.labels.job=="api-gateway")'
   ```

2. Verify Grafana dashboard loads:
   ```bash
   Visit: http://grafana:3000/d/api-gateway-third-party
   Check panels: request rate, error rate, latency, auth failures, rate-limits
   ```

3. Test alert routing (simulate spike):
   ```bash
   # Trigger test alert
   curl -X POST http://alertmanager:9093/api/v2/alerts -H "Content-Type: application/json" \
     -d '[{"labels":{"alertname":"APIGatewayHighErrorRate","severity":"critical"},"annotations":{"summary":"Test alert"}}]'
   
   # Verify Slack notification received in #aystek-alerts
   ```

4. Tune alert thresholds based on baseline:
   - If <5% 5xx rate during normal load, set alert threshold to >10% for 5 min
   - If median latency <100ms, set high latency alert to >500ms

**Acceptance:** All monitoring dashboards and alerts operational

---

## Day 3-5: Integration Phase

### Activity 3.1: Team Integration Kickoff (30 min)

**Responsible:** Tech Lead

**Meeting Schedule:**
```
Day 3, 10am: All team leads
- Review API client credentials (share encrypted link)
- Walk through integration flow (OAuth2 token request → API call)
- Review error handling (401/403/429 responses)
- Q&A on scope, rate limits, quotas
```

**Share with Each Team:**
1. Encrypted 1Password link with `client_id` + `client_secret`
2. Integration guide (`INTEGRATION_GUIDE.md`)
3. Sample code (Go HTTP client, JavaScript fetch, Python requests)
4. Postman collection for manual testing

---

### Activity 3.2: Backend Integration (Day 3-4)

**Responsible:** Backend Team

**Checklist:**

- [ ] **Obtain Token**
  ```go
  // Backend Service integration
  func init() {
      client := &http.Client{}
      req, _ := http.NewRequest("POST", "https://api.aystek.dev/auth/oauth2/token", 
          strings.NewReader(fmt.Sprintf(
              "grant_type=client_credentials&client_id=%s&client_secret=%s&scope=%s",
              os.Getenv("API_CLIENT_ID"),
              os.Getenv("API_CLIENT_SECRET"),
              "products:read,orders:read,orders:create,users:read",
          )))
      req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
      resp, _ := client.Do(req)
      // Cache token + refresh before expiry (15 min)
  }
  ```

- [ ] **Call Protected Endpoint**
  ```go
  // Use cached token in Authorization header
  req, _ := http.NewRequest("GET", 
      "https://api.aystek.dev/products?limit=100",
      nil)
  req.Header.Set("Authorization", "Bearer "+cachedToken)
  resp, _ := client.Do(req)
  ```

- [ ] **Handle Errors**
  ```go
  switch resp.StatusCode {
  case 401: // Token expired or invalid
      // Refresh token or obtain new one
  case 403: // Scope or route unauthorized
      // Check scopes against API client config
  case 429: // Rate limited
      // Implement exponential backoff
      // Check X-RateLimit-Reset header
  }
  ```

- [ ] **Test with 100 Requests**
  ```bash
  # Simple load test
  for i in {1..100}; do
    curl -H "Authorization: Bearer $TOKEN" \
         https://api.aystek.dev/products > /dev/null
  done
  ```

**Success Criteria:**
- All 100 requests succeed (no 401/403 errors)
- Response time < 500ms each
- No rate-limit exceeded (429) unless exceeding limits

---

### Activity 3.3: Frontend Integration (Day 3-4)

**Responsible:** Frontend Team

**Checklist:**

- [ ] **Implement Token Caching** (JavaScript example)
  ```javascript
  class APIClient {
    constructor(clientId, clientSecret) {
      this.clientId = clientId;
      this.clientSecret = clientSecret;
      this.token = null;
      this.tokenExpiresAt = null;
    }
    
    async getToken() {
      if (this.token && this.tokenExpiresAt > Date.now() + 60000) {
        return this.token; // Reuse if valid for >1 min
      }
      
      const response = await fetch('https://api.aystek.dev/auth/oauth2/token', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          grant_type: 'client_credentials',
          client_id: this.clientId,
          client_secret: this.clientSecret,
          scope: 'products:read,orders:read,users:read',
        }),
      });
      
      const data = await response.json();
      this.token = data.access_token;
      this.tokenExpiresAt = Date.now() + (data.expires_in * 1000);
      return this.token;
    }
    
    async call(method, path, body) {
      const token = await this.getToken();
      const response = await fetch(`https://api.aystek.dev${path}`, {
        method,
        headers: { 'Authorization': `Bearer ${token}` },
        body: body ? JSON.stringify(body) : undefined,
      });
      
      if (response.status === 401) {
        this.token = null; // Invalidate token
        return this.call(method, path, body); // Retry
      }
      
      return response;
    }
  }
  ```

- [ ] **Test Dashboard Pages**
  - Products list page (GET /products)
  - Orders list page (GET /orders)
  - Order creation (POST /orders)
  - Error page (test 403 for unauthorized scope)

- [ ] **Monitor Network Tab**
  - OAuth2 token request completes in <1s
  - API endpoint requests < 500ms
  - No 401 due to expired token (token cached properly)

---

### Activity 3.4: QA Testing (Day 4-5)

**Responsible:** QA Team

**Test Scenarios:**

- [ ] Happy Path
  ```
  1. Obtain token with QA client credentials
  2. GET /products → 200 OK
  3. GET /orders → 200 OK
  4. POST /orders with valid payload → 201 Created
  5. Verify response contains expected fields
  ```

- [ ] Scope Validation
  ```
  1. QA client has scope: products:read, orders:read, orders:create
  2. GET /products → 200 OK (allowed)
  3. DELETE /products/:id → 403 Forbidden (not in scope)
  4. POST /orders → 201 Created (allowed)
  5. DELETE /orders/:id → 403 Forbidden (not in scope)
  ```

- [ ] Rate Limiting
  ```
  1. QA client rate limit: 2,000 req/hour
  2. Send 50 requests → all succeed
  3. Verify X-RateLimit-Limit, X-RateLimit-Remaining headers
  4. Verify each request decrements remaining count
  ```

- [ ] Error Handling
  ```
  1. Invalid token → 401 Unauthorized
  2. Expired token → 401 Unauthorized
  3. No Authorization header → 401 Unauthorized
  4. Unauthorized scope → 403 Forbidden
  5. Rate limit exceeded → 429 Too Many Requests
  ```

**Document Results:**
- [ ] Test summary in AYSTEK Jira (test_pilot_integration)
- [ ] Screenshots of API responses
- [ ] Performance metrics (p50/p95 latency)

---

## Day 6-7: Validation Phase

### Activity 6.1: 100-Hour Incident-Free Window

**Responsible:** On-call Engineer + Monitoring

**Monitoring Checklist:**

- [ ] **Hourly Status Check (every 4 hours):**
  ```
  Time    | 5xx Rate | Latency P95 | Auth Fails | Rate Limits | Notes
  --------|----------|-------------|-----------|-------------|-------
  Hour 0  | 0.1%     | <50ms       | 0         | 0           | Pilot started
  Hour 4  | 0.2%     | <60ms       | 0         | 0           | All teams active
  ...
  Hour 96 | 0.1%     | <70ms       | 0         | 0           | Still stable
  Hour 100| 0.1%     | <50ms       | 0         | 0           | ✅ Milestone reached
  ```

- [ ] **Alert Response SLA:**
  - [ ] Any alert → respond within 5 minutes
  - [ ] Critical alert (5xx > 5%) → page on-call immediately
  - [ ] Investigate root cause and implement fix

- [ ] **Team Feedback Collection (Daily):**
  - [ ] #aystek-api-gateway-pilot Slack channel for issues
  - [ ] Daily standup: "Any issues with API client integration?"
  - [ ] Log all feedback in issue tracker

---

### Activity 6.2: Audit Log Export Validation (Day 6)

**Responsible:** Compliance Officer

**Steps:**

1. Export audit logs for all 5 API clients:
   ```bash
   curl -H "Authorization: Bearer admin-token" \
     "https://api.aystek.dev/admin/api-clients/{client_id}/audit-logs/export?format=csv" \
     > audit-logs-{client_id}.csv
   ```

2. Verify CSV contains all required fields:
   ```
   client_id, tenant_id, method, path, status_code, timestamp, request_id, response_time_ms
   ```

3. Check record count matches expected API calls
   - Backend: ~5000 requests
   - Frontend: ~3000 requests
   - QA: ~1000 requests

4. Verify no PII exposed in CSV (data classification correct)

5. Test formula injection protection:
   ```bash
   # Verify dangerous characters escaped
   grep "^'" audit-logs-*.csv | head -5
   ```

---

### Activity 6.3: Secret Rotation Test (Day 6-7)

**Responsible:** DevOps / Security Team

**Steps:**

1. Select one API client for rotation test (e.g., backend-integration-internal):
   ```bash
   curl -X POST \
     -H "Authorization: Bearer admin-token" \
     "https://api.aystek.dev/admin/api-clients/{client_id}/rotate-secret"
   ```

2. Record new secret securely

3. Verify old secret works during 5-minute grace period:
   ```bash
   # Old secret should still work
   curl -X POST "https://api.aystek.dev/auth/oauth2/token" \
     -d "grant_type=client_credentials&client_id=...&client_secret={OLD_SECRET}"
   # Should return 200 OK
   ```

4. Verify new secret works:
   ```bash
   # New secret should work
   curl -X POST "https://api.aystek.dev/auth/oauth2/token" \
     -d "grant_type=client_credentials&client_id=...&client_secret={NEW_SECRET}"
   # Should return 200 OK
   ```

5. Wait 6 minutes, verify old secret rejected:
   ```bash
   sleep 360
   curl -X POST "https://api.aystek.dev/auth/oauth2/token" \
     -d "grant_type=client_credentials&client_id=...&client_secret={OLD_SECRET}"
   # Should return 401 Unauthorized
   ```

6. Verify rotation event logged:
   ```bash
   curl -H "Authorization: Bearer admin-token" \
     "https://api.aystek.dev/admin/api-clients/{client_id}/secret-rotations" \
     | jq '.rotations[0]'
   # Should show: client_id, rotated_at, rotated_by, reason
   ```

---

## Day 7-8: Feedback & Fixes

### Activity 7.1: Team Retrospective

**Responsible:** Tech Lead

**Meeting:** Day 7, 2pm (all teams)

**Agenda:**

1. **What Went Well?** (10 min)
   - Smooth OAuth2 integration?
   - Good error messages?
   - Docs clear and helpful?

2. **What Challenges Did You Face?** (15 min)
   - Token caching issues?
   - Confusing rate-limit headers?
   - Rate limits too aggressive?
   - Scope requirements unclear?

3. **Action Items** (10 min)
   - Adjust rate limits for teams needing more capacity
   - Update docs if integration confusing
   - File bugs for any errors or crashes

**Document Results:** Retrospective notes in Confluence

---

### Activity 7.2: Bug Triage & Fixes (Day 7-8)

**Responsible:** API Gateway Team

**Process:**

1. Review all feedback from Slack:
   ```bash
   # Export all messages from #aystek-api-gateway-pilot
   docker exec slack_exporter slack-export-bulk "#aystek-api-gateway-pilot" > pilot-feedback.txt
   ```

2. Categorize issues:
   - **Critical** (P0): 5xx errors, token failing, audit logs missing
   - **High** (P1): Rate-limit too aggressive, confusing error messages
   - **Medium** (P2): Documentation improvements, nice-to-have features

3. Fix critical/high-priority issues same day

4. Verify fixes with QA team on Day 8

5. Update runbooks/docs if needed

---

## Post-Pilot Success Criteria

### Operational Metrics
- [ ] **Availability:** ≥99.5% (max 3.6 minutes downtime in 100 hours)
- [ ] **Error Rate:** <0.1% (5xx errors)
- [ ] **Latency P95:** <500ms
- [ ] **Auth Failures:** <0.1%
- [ ] **Rate-Limit Fairness:** All teams stay within limits fairly

### Team Feedback
- [ ] **Developer Satisfaction:** ≥4/5 stars (teams polled on ease of integration)
- [ ] **Documentation:** Clear and actionable
- [ ] **Support:** Issues resolved within 30 minutes
- [ ] **Critical Bugs:** 0 (no production incidents requiring emergency fix)

### Audit & Compliance
- [ ] **Audit Logs:** Complete and queryable for all requests
- [ ] **Secret Rotation:** Verified working with 5-minute grace period
- [ ] **IP Whitelist:** Enforcement verified (non-whitelisted IPs blocked)
- [ ] **KVKK/GDPR:** Data classification correct, no PII exposure in logs

### Go-Live Readiness
- [ ] **Internal Pilot:** Incident-free 100 hours ✅
- [ ] **Documentation Updated:** Based on team feedback
- [ ] **Monitoring Tuned:** Alert thresholds validated with real traffic
- [ ] **Team Training:** Completed (engineers confident with OAuth2 flow)
- [ ] **Support Plan:** On-call schedule confirmed, escalation clear

---

## Escalation & Support

### On-Call Support (Week 1)

**Hours:** 24/7 (including weekends)

**Slack Channel:** #aystek-api-gateway-pilot
- Response time: <30 minutes

**Escalation Chain:**
1. API Gateway On-Call → API Gateway Team Lead
2. API Gateway Team Lead → DevOps Manager (if infrastructure issue)
3. DevOps Manager → VP Engineering (if critical incident)

**Common Issues & Fixes:**

| Issue | Symptom | Fix |
|-------|---------|-----|
| Token Expired | 401 Unauthorized | Refresh token via `/auth/oauth2/token` with `refresh_token` grant |
| Rate Limited | 429 Too Many Requests | Wait for `X-RateLimit-Reset` header before retrying |
| Scope Unauthorized | 403 Forbidden | Verify scope in admin UI includes required endpoints |
| IP Blocked | 403 Forbidden (with IP in logs) | Whitelist source IP in admin UI |
| Quota Exceeded | 429 Too Many Requests | Contact ops to increase quota |

---

## Sign-Off & Phase 6 Completion

**Phase 6 Complete When:**

✅ All pilot success criteria met  
✅ 100-hour incident-free window achieved  
✅ Team feedback incorporated  
✅ Documentation updated  
✅ Monitoring validated  

**Approval Signature:**
- [ ] API Gateway Team Lead: ___________________ Date: ___
- [ ] DevOps Lead: ___________________ Date: ___
- [ ] Security Lead: ___________________ Date: ___

**Next Step:** Proceed to Phase 7 (Production Operations & First Partner Onboarding)

---

**Last Updated:** March 9, 2026  
**Owner:** API Gateway Team Lead
