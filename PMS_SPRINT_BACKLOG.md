# PMS SPRINT BACKLOG (Enterprise Roadmap)

**Document Version:** 1.0.0  
**Generated:** 2026-02-21  
**Status:** Q1 Planning (Sprint-15+)

---

## Overview

Bu backlog **PMS Suite** for enterprise go-live readiness'i define eder. TASK_037-TASK_040 tamamlandıktan sonra kalan işler.

---

## P0: Release Blockers (Security/Stability/Correctness)

### P0-001: Access Suite State Machine Hardening & Edge Cases

**Priority:** CRITICAL  
**Sprint:** Sprint-14.B (Must Complete)  
**Owner:** PMS/Access Team

**Why:**
- Visitors/vehicles state machines incomplete edge case handling
- Double checkout/exit → should 409 but may silently accept
- Tenant/site scoping validation missing in access endpoints
- PTS ingest signature replay vulnerability (no idempotency key)

**Modules Affected:**
- `internal/modules/pms/access/visitors`
- `internal/modules/pms/access/vehicles`

**Acceptance Criteria:**
```
✅ Double checkout attempts return 409 pms.error.already_checked_out
✅ Double exit attempts return 409 pms.error.already_exited
✅ Wrong tenant/site access attempts return 404 or 403
✅ PTS ingest: repeated payload with same action returns 409 (or idempotency guard)
✅ All state transitions validated before execution
✅ Unit tests: 15+ edge case scenarios
✅ go test ./internal/modules/pms/access/... ✅
```

**Implementation Tasks:**
1. Add idempotency guards to PTS ingest endpoint (external_event_id or timestamp-based)
2. Implement state transition matrix validation (open → out only valid from "in" state)
3. Add tenant/site scoping to access repository filters
4. Write 20+ unit tests for state machine edge cases
5. Audit logging for rejected state transitions

**Test Scenarios:**
- Single visitor: in → out → out (attempt) → 409 ✅
- Vehicle double-exit: in → out → out (attempt) → 409 ✅
- Cross-tenant access attempt: tenant_a_user tries access tenant_b_visitor → 404 ✅
- PTS signature mismatch: reject 401 ✅
- PTS replay: same payload twice → 409 (idempotency) ✅

**Estimated Effort:** 3-4 days  
**Risk:** High (access control + state machine = security critical)

---

### P0-002: Payment Webhook Auto-Finalize Verification & Idempotency

**Priority:** CRITICAL  
**Sprint:** Sprint-14 (Post-TASK_040)  
**Owner:** Billing/Payments Team

**Why:**
- TASK_040 implemented webhook → finalize orchestration
- Webhook retries could trigger duplicate finalize attempts
- Idempotency guard may not be watertight; needs verification
- Payment allocation must be atomic (no double-charge)

**Modules Affected:**
- `internal/modules/payments` (callback orchestration)
- `internal/modules/pms/portal` (finalize endpoint)
- `internal/modules/pms/billing` (allocation check)

**Acceptance Criteria:**
```
✅ Webhook callback triggers finalize via service method
✅ Finalize idempotency key prevents duplicate execution
✅ Duplicate finalize attempt returns 409 or success (consistent)
✅ Allocation only executed once per intent
✅ Outbox event pms.payment.received emitted exactly once
✅ End-to-end test: webhook → payment creation → finalize → charge allocation
✅ go test ./internal/modules/pms/portal/... -run Finalize ✅
```

**Implementation Tasks:**
1. Write integration test: mock webhook callback → verify finalize called
2. Verify idempotency key is stored and checked before finalize
3. Check allocation idempotency via billing_idempotency collection
4. Add test for webhook retry (simulate 3+ duplicate callbacks)
5. Verify outbox event de-duplication (pms.payment.received logged only once)

**Test Scenarios:**
- Webhook delivers successfully → finalize executed → charge allocated ✅
- Webhook delivers, duplicate callback received → finalize idempotent ✅
- Webhook delivers, duplicate callback + eventual finalize call → no double-charge ✅
- Payment allocation already done (stale payment) → 409 or success ✅

**Estimated Effort:** 2-3 days  
**Risk:** High (payment + allocation = data integrity + revenue risk)

---

### P0-003: Metering Mixed Allocation Engine Correctness

**Priority:** CRITICAL  
**Sprint:** Sprint-14.C (Must Complete)  
**Owner:** Metering Team

**Why:**
- TASK_039 implemented metering but allocation model incomplete
- "Mixed" method requires BOTH fixed_base_cost (equal split) + consumption_cost (by_consumption)
- If formula wrong: charges incorrect, tenant disputes inevitable
- Publish must be idempotent (no duplicate charges on retry)

**Modules Affected:**
- `internal/modules/pms/metering`
- `internal/modules/pms/billing` (bridge CreateManualCharge)

**Acceptance Criteria:**
```
✅ Meter allocation preview shows correct split:
    - fixed_base_cost: allocated equally across N units
    - consumption_cost: allocated by consumption ratio (unit_reading / total_reading)
✅ Publish creates charges via bridge with charge_type=utility_share
✅ Republish same allocation: no new charges (idempotent via charge_type dedup)
✅ All 3 methods work: equal_split, by_consumption, mixed
✅ Backward compat: enum aliases (equal→equal_split, consumption→by_consumption)
✅ Unit tests: 20+ allocation scenarios
✅ go test ./internal/modules/pms/metering/... ✅
```

**Implementation Tasks:**
1. Implement mixed allocation formula:
   ```
   per_unit_fixed = fixed_base_cost / num_units
   consumption_ratio = unit_reading / sum_all_readings
   per_unit_consumption = consumption_cost * consumption_ratio
   unit_total = per_unit_fixed + per_unit_consumption
   ```
2. Verify bridge charge_type=utility_share uniqueness (all 3 method charges must use this type)
3. Add idempotency check: don't re-create charges if already published for same allocation
4. Write 25+ unit test scenarios (fixed, consumption, mixed variants)
5. Test backward compatibility (legacy enum values)

**Test Scenarios:**
- 3 units, mixed 50€ fixed + 100€ consumption (readings: 100, 200, 200):
  - Unit 1: 50/3 + 100*(100/500) = €36.67 ✅
  - Unit 2: 50/3 + 100*(200/500) = €53.33 ✅
  - Unit 3: 50/3 + 100*(200/500) = €53.33 ✅
  - Total: €143.33 (50 + 100) = ✅
- Republish same allocation: no new charges ✅
- Legacy enum "equal" works same as "equal_split" ✅

**Estimated Effort:** 3-4 days  
**Risk:** High (billing correctness + significant number of calculations)

---

## P1: Enterprise Features (Go-Live Ready)

### P1-001: Portfolio Cross-Site KPI Dashboard

**Priority:** HIGH  
**Sprint:** Sprint-15 (Early)  
**Owner:** Dashboard Team

**Why:**
- TASK_034 created portfolio basic structure
- Enterprise requirement: tenant sees aggregated KPIs across multi-site portfolio
- Portfolio view is critical for management company scenarios
- Missing: top-5 overdue units ranking, site drill-down correctness

**Modules Affected:**
- `internal/modules/pms/portfolio`
- `internal/modules/pms/dashboard` (site-level ref)

**Acceptance Criteria:**
```
✅ GET /admin/pms/portfolio/dashboard?period=YYYY-MM returns:
  {
    tenant_kpis: {
      total_units, total_due, total_paid, total_remaining,
      collection_rate, avg_occupancy_rate
    },
    sites: [
      {site_id, site_name, kpis: {...}},
      ...
    ],
    top_overdue_units: [
      {unit_id, site_id, amount_due, days_overdue},
      ...
    ] (top 5)
  }
✅ Optional site_id filter returns single site data
✅ Tenant isolation: can't view other tenant's portfolio
✅ Performance: response < 1s for 100+ sites
✅ UI page displays KPI cards + site grid + overdue table
✅ npx eslint app/(dashboard)/admin/pms/portfolio/page.tsx ✅
```

**Implementation Tasks:**
1. Implement portfolio dashboard service:
   - Aggregate KPIs from all sites (or filtered site)
   - Calculate collection_rate = paid / due
   - Calculate avg_occupancy_rate
2. Add top-5 overdue query (order by days_overdue desc, limit 5)
3. Add site drill-down data structure
4. Create Next.js portfolio page (`app/(dashboard)/admin/pms/portfolio/dashboard/page.tsx`)
5. Add i18n keys for portfolio UI (`pms.portfolio.*`)

**Test Scenarios:**
- Tenant with 3 sites, different KPIs: response aggregates correctly ✅
- Top-5 overdue: returns 5 units with max aging ✅
- Site filter: returns only filtered site KPIs ✅
- Cross-tenant access attempt: returns other tenant 404 ✅
- Performance: 100 sites, <1s response ✅

**Estimated Effort:** 4-5 days  
**Risk:** Medium (aggregation query performance on large tenant)

---

### P1-002: Reconciliation Advanced Scoring & Matching Heuristics

**Priority:** HIGH  
**Sprint:** Sprint-15 (Mid)  
**Owner:** Reconciliation Team

**Why:**
- TASK_018 created basic scoring
- Enterprise users need intelligent matching (not just amount closeness)
- Heuristic must be deterministic (replay-consistent)
- False positives lead to manual review overhead

**Modules Affected:**
- `internal/modules/pms/reconciliation`

**Acceptance Criteria:**
```
✅ Match scoring considers:
  - Amount closeness (±2% tolerance or configurable)
  - Party name similarity (exact or fuzzy match)
  - Transaction date proximity (within 3 days window)
  - Description keyword match
✅ Scoring deterministic: same payload → same score
✅ Tiebreaker rule: creation_date ascending (stable sort)
✅ Minimum score threshold (e.g., >70 points) to auto-propose
✅ UI shows match confidence percentage
✅ Admin can override/reject all proposed matches
✅ Unit tests: 30+ scoring scenarios
```

**Implementation Tasks:**
1. Define scoring rubric (weights for each heuristic)
2. Implement amount closeness calc (within tolerance)
3. Implement party name matching (exact first, fuzzy if no exact match)
4. Implement date proximity scoring
5. Add deterministic tiebreaker (sort by created_at)
6. Threshold-based proposal (>70 points triggers auto-propose)
7. Write 40+ unit tests covering scoring edge cases

**Test Scenarios:**
- Amount exact match: 100 points ✅
- Amount within ±1%: 90 points ✅
- Amount outside ±5%: 0 points ✅
- Party name exact: +20 points ✅
- Party name fuzzy (1 char diff): +10 points ✅
- Date within 1 day: +15 points ✅
- Same amount + party + date: 100+ points (high confidence) ✅
- Tiebreaker deterministic: two identical txns → always same sort order ✅

**Estimated Effort:** 4-5 days  
**Risk:** Medium (scoring heuristic non-determinism hard to debug)

---

### P1-003: Export Job Concurrency & Lock Hardening

**Priority:** MEDIUM-HIGH  
**Sprint:** Sprint-15 (Mid)  
**Owner:** Exports Team

**Why:**
- TASK_030 hardened single-worker claim-lock
- Enterprise scale: need multi-worker support (100+ concurrent exports)
- Current single-lock may bottleneck large sites
- Distributed lock required for multi-deployment scenarios

**Modules Affected:**
- `internal/modules/pms/exports`
- `internal/queue/exports_worker.go`

**Acceptance Criteria:**
```
✅ Multi-worker concurrency support:
  - Each worker claims lock atomically
  - Lock TTL: 5 min; auto-unlock on completion
  - Stuck job recovery: lock_expires_at > now → reclaim
✅ Distributed lock via MongoDB atomic update:
  updateOne({_id, status: queued}, {$set: {status: running, lock_id, lock_expires_at}})
  ← only one worker succeeds
✅ File integrity: no concurrent writes to same file
✅ Retry logic: max 3 attempts, dead-letter after
✅ Stress test: 100 concurrent export jobs → all succeed
```

**Implementation Tasks:**
1. Implement distributed atomic claim (via MongoDB updateOne)
2. Add claim failure handling (another worker already claimed)
3. Implement lock expiration check (stuck recovery)
4. Write stress test: spawn 100 export jobs, verify atomicity
5. Add metrics: exports_jobs_processing gauge

**Test Scenarios:**
- Worker 1 claims job → succeeds, lock_id set ✅
- Worker 2 attempts same job → fails (different lock_id) ✅
- Lock expires (5 min) → Worker 3 reclaims → succeeds ✅
- 100 concurrent jobs on 4 workers → all process without collision ✅

**Estimated Effort:** 3-4 days  
**Risk:** Medium (distributed lock bug causes data corruption)

---

### P1-004: Webhook Delivery Monitoring & Retry Instrumentation

**Priority:** MEDIUM-HIGH  
**Sprint:** Sprint-15 (Late)  
**Owner:** Webhooks Team

**Why:**
- TASK_026/TASK_029 implemented webhook delivery + worker
- Enterprise needs visibility: which delivery failed? why? retry status?
- Silent failures lead to integration gaps (partner data sync breaks)
- Admin needs dashboard for webhook health

**Modules Affected:**
- `internal/modules/pms/webhooks`
- `internal/queue/webhooks_worker.go`

**Acceptance Criteria:**
```
✅ Delivery model extended:
  - retry_count: current attempt number
  - last_error_message: human-readable error
  - last_error_code: HTTP status or error enum
  - next_retry_at: RFC3339 timestamp
✅ Admin endpoint: GET /admin/pms/webhooks/deliveries?status=failed&limit=100
✅ Dashboard: failed delivery count, retry distribution chart
✅ Metrics: webhook_delivery_retries_total{status, subscription}
✅ Audit log: webhook_retry_attempted, webhook_delivery_succeeded/failed
```

**Implementation Tasks:**
1. Extend delivery schema (retry_count, last_error_message, last_error_code, next_retry_at)
2. Update worker to populate error fields on failure
3. Create delivery list endpoint with filters (status, subscription, date range)
4. Create admin UI page: `/admin/pms/webhooks` (deliveries tab)
5. Add Prometheus metrics for retry distribution

**Test Scenarios:**
- Delivery fails with HTTP 500 → retry_count=1, error_code=500 ✅
- After 5 retries with all failures → status=dead_letter, audit log recorded ✅
- Admin filters by status=failed → sees 15 failed deliveries ✅
- Metrics show retry distribution (1 retry: 50%, 2 retries: 30%, etc) ✅

**Estimated Effort:** 3-4 days  
**Risk:** Low (observability only; no business logic change)

---

### P1-005: Portal Rate Limiting Distributed

**Priority:** MEDIUM  
**Sprint:** Sprint-16 (Late)  
**Owner:** Portal Team

**Why:**
- TASK_035 implemented in-memory rate limit
- Single-instance MVP; multi-deployment needed for enterprise
- DDoS protection insufficient (in-memory limits per-instance)
- Redis-backed rate limit required for distributed scenarios

**Modules Affected:**
- `internal/modules/pms/portal`

**Acceptance Criteria:**
```
✅ Rate limit moved to Redis:
  - Key: portal:ratelimit:{tenant_id}:{user_id}
  - Window: configurable (default PMS_PORTAL_RATE_LIMIT_WINDOW_SEC=60)
  - Limit: configurable (default PMS_PORTAL_RATE_LIMIT_PER_MIN=100)
✅ Consistent across all deployment instances
✅ Graceful degradation: if Redis down, allow requests (loose rate limit)
✅ Metrics: portal_ratelimit_exceeded_total
```

**Implementation Tasks:**
1. Move rate limit logic to Redis (lua script for atomic increment)
2. Add Redis connection pool
3. Implement graceful degradation (if Redis unavailable, don't block)
4. Add metrics counter for exceeded limits
5. Write integration tests (mock Redis)

**Test Scenarios:**
- 100 requests from same user within 60s → requests 1-100 succeed, 101+ get 429 ✅
- Requests from different users: each has own limit bucket ✅
- Redis down: requests still allowed (degraded) ✅
- Metrics show exceeded_total counter incrementing ✅

**Estimated Effort:** 2-3 days  
**Risk:** Low (distributed caching, well-established pattern)

---

## P2: Nice-to-Have (Post-MVP)

### P2-001: Bulk Operations

**Priority:** LOW  
**Sprint:** Sprint-16+  
**Owner:** Admin Team

**Why:**
- Enterprise users manage 100+ work orders; single-edit tedious
- Bulk status toggle, vendor archive, dues repub needed
- Dry-run mode for safety (preview before commit)

**Scope:**
- Bulk work order status change
- Bulk vendor archive
- Bulk dues republish with preview

**Estimated Effort:** 5-6 days  
**Risk:** Low (non-critical feature)

---

### P2-002: Advanced Filters & Search

**Priority:** LOW  
**Sprint:** Sprint-16+  
**Owner:** UX Team

**Why:**
- Lists support basic filters only
- Date range filters for reconciliation, exports
- Priority/category filters for work orders
- Full-text search on descriptions

**Scope:**
- Date range picker on reconciliation queue
- Work order filters: priority, category, assigned_to
- Full-text search: vendors, units

**Estimated Effort:** 4-5 days  
**Risk:** Low (UX polish)

---

### P2-003: PDF Receipt Generation & Email Dispatch

**Priority:** LOW  
**Sprint:** Sprint-16+  
**Owner:** Legal/Portal Team

**Why:**
- Currently: JSON-only receipt
- Enterprise: PDF receipts for tenant portal
- Email dispatch: receipt to tenant email

**Scope:**
- PDF rendering: payment receipt with logo, terms, signature
- Email template + dispatch
- Portal download link

**Estimated Effort:** 5-7 days (wkhtmltopdf setup complexity)  
**Risk:** Medium (external library + email delivery)

---

### P2-004: Portal Multi-Language Support

**Priority:** LOW  
**Sprint:** Sprint-17  
**Owner:** i18n Team

**Why:**
- Portal UI currently hardcoded English
- Tenant preference for locale in portal
- i18n key coverage incomplete (pms.portal.*)

**Scope:**
- Add language selector to portal
- Ensure all portal UI keys translated (28 languages)
- Store user language preference

**Estimated Effort:** 3-4 days  
**Risk:** Low (i18n pattern established)

---

### P2-005: Webhook Event Filtering & Delivery Logs

**Priority:** LOW  
**Sprint:** Sprint-17  
**Owner:** Webhooks Team

**Why:**
- All events dispatched to all subscriptions (no filtering)
- UX improvement: filter events per subscription
- Delivery logs UI for debugging

**Scope:**
- Webhook subscription: select events to receive
- Delivery log viewer (timestamps, payloads, responses)

**Estimated Effort:** 3-4 days  
**Risk:** Low (UI + query enhancement)

---

## Summary: Backlog Timeline Estimate

| Phase | Sprints | Items | Est. Effort |
|-------|---------|-------|------------|
| **P0 Release Blockers** | Sprint-14.B, 14.C | 3 items | 10 days |
| **P1 Enterprise Features** | Sprint-15, 16 | 5 items | 19-21 days |
| **P2 Nice-to-Have** | Sprint-16+ | 5 items | 20-25 days |
| **TOTAL** | - | **13 items** | **49-56 days** |

**Recommended Staffing:** 3 engineers full-time = 2-3 full sprints (3 weeks per sprint)

---

## Dependency Tree

```
P0-001 (Access Hardening)
  └─ Can proceed independently

P0-002 (Payment Webhook)
  └─ Depends on: P0-001 (for integration stability)

P0-003 (Metering Mixed)
  └─ Can proceed independently

↓

P1-001 (Portfolio Dashboard)
  └─ Depends on: P0-1, P0-2, P0-3 (all P0 complete)

P1-002 (Reconciliation Scoring)
  └─ Can proceed independently (P0-002 optional)

P1-003 (Export Concurrency)
  └─ Can proceed independently

P1-004 (Webhook Monitoring)
  └─ Should wait: P1-002 or P1-003 for precedent on error handling

P1-005 (Portal Rate Limiting)
  └─ Can proceed independently

↓

P2-* (All nice-to-have)
  └─ Can proceed after P1 complete
```

**Critical Path:** P0-1 → P0-2 → P0-3 → P1-1 = **~18-20 days minimum**

---

## Success Metrics

| Metric | Target | Owner |
|--------|--------|-------|
| All P0 tests pass | 100% ✅ | QA |
| go test ./internal/modules/pms/... | 100% passing | QA |
| npm run lint on /admin/pms/** | 0 errors | Frontend |
| Portfolio dashboard response time | <1s (100 sites) | Backend |
| Export concurrency stress test | 100% atomicity | Backend |
| Webhook retry success rate | >99% (with retries) | QA |
| Portal uptime | >99.5% | DevOps |

---

**END OF PMS_SPRINT_BACKLOG.md**

*Enterprise Roadmap | Q1 2026*
