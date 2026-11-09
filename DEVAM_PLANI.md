# PMS DEVAM PLANI (Continuation Plan for Sprint-15+)

**Author:** Principal Engineer (AI Analysis)  
**Date:** 2026-02-21  
**Status:** Ready for Handoff

---

## 📋 Executive Summary

Sprint-14 (TASK_037 → TASK_040) **Access Suite** ve **Payment Webhook** ağırlıklı olmak üzere tamamlandı. PMS Suite enterprise-ready midir?

**Verdict:** 90% feature-complete; **3 critical P0 items** kaldı:

1. ✅ **Access Suite Edge Cases** (State machine validation)
2. ✅ **Payment Webhook Idempotency** (Double-finalize guard)
3. ✅ **Metering Mixed Engine** (Fixed + consumption formula)

Bu 3 item'i **Sprint-14.B & 14.C** içinde kapatmalısınız.

Sonra **Sprint-15** ile P1 features'a geçiş yapabilirsiniz (Portfolio, Reconciliation scoring, Export concurrency).

---

## 🎯 İlk Yapılacaklar (Next 48 Hours)

### Öncelik 1: Kod Auditı

**Yapılacak:**
1. Read [PMS_MASTER.md](PMS_MASTER.md) Section B (Bounded Contexts)
   - Tüm 16 modules'ı 30 dakikada skim
   - Critical state machines tanımla

2. Read [PMS_SPRINT_BACKLOG.md](PMS_SPRINT_BACKLOG.md) Section P0-001 thru P0-003
   - Edge case senaryolarını not et
   - Test coverage gaps'ı tespit et

3. Run existing tests:
   ```bash
   cd /workspace/go/aystek-core
   go test ./internal/modules/pms/access/... -v
   go test ./internal/modules/pms/billing/... -v
   go test ./internal/modules/pms/metering/... -v
   ```
   - Failures'ı not et
   - Test counts'ı say (target: 100+ tests)

### Öncelik 2: Branching Strategy

**Branches:**
- `feature/pms-p0-access-hardening` ← P0-001 implementation
- `feature/pms-p0-payment-webhook-verify` ← P0-002 verification
- `feature/pms-p0-metering-mixed` ← P0-003 implementation

**PR Scope:** Dar tut (~150-200 lines max per PR)

**Approval:** Code review + QA test run required

### Öncelik 3: Test Suite Expansion

**Her P0 için 20+ test ekle:**

```bash
# P0-001: Access hardening
go test ./internal/modules/pms/access/visitors -run TestDoubleCheckout
go test ./internal/modules/pms/access/vehicles -run TestDoubleExit
go test ./internal/modules/pms/access -run TestTenantIsolation

# P0-002: Payment webhook
go test ./internal/modules/pms/portal -run TestFinalizeIdempotency
go test ./cmd/api -run TestWebhookPaymentAutoFinalize

# P0-003: Metering mixed
go test ./internal/modules/pms/metering -run TestMixedAllocation
go test ./internal/modules/pms/billing -run TestUtilityShareIdempotency
```

---

## 📅 Week-by-Week Roadmap (Sprint-14.B → Sprint-15)

### Week 1: P0-001 Access Hardening (Feb 21-27)

| Day | Task | Owner | Acceptance |
|-----|------|-------|-----------|
| Fri-22 | Read PMS_MASTER.md + audit current access code | Dev | Understand state machines |
| Mon-25 | Edge case test suite (20+ tests) | QA | Tests written, some fail |
| Tue-26 | State machine validation code | Dev | All tests pass |
| Wed-27 | Tenant/site scoping + idempotency guards | Dev | Code review ready |
| Thu-28 | PR review + integration test | QA | go test ./internal/modules/pms/access ✅ |

**Deliverable:** `feature/pms-p0-access-hardening` merged to dev

---

### Week 2: P0-002 Payment Webhook (Mar 3-9)

| Day | Task | Owner | Acceptance |
|-----|------|-------|-----------|
| Mon-3 | Understand TASK_040 + webhook callback flow | Dev | Diagram sketch |
| Tue-4 | End-to-end integration test (webhook → finalize) | QA | Test written, fails |
| Wed-5 | Idempotency guard verification | Dev | Code review |
| Thu-6 | Duplicate finalize scenarios | QA | 5 scenarios tested ✅ |
| Fri-7 | Outbox event deduplication | Dev | go test ./cmd/api ✅ |

**Deliverable:** `feature/pms-p0-payment-webhook-verify` merged to dev

---

### Week 3: P0-003 Metering Mixed (Mar 10-16)

| Day | Task | Owner | Acceptance |
|-----|------|-------|-----------|
| Mon-10 | Design mixed allocation formula | Dev | Whiteboard approved |
| Tue-11 | Implement mixed method in allocation service | Dev | Code review |
| Wed-12 | Bridge charge_type=utility_share idempotency | Dev | Integration test |
| Thu-13 | Backward compat enum aliases (equal→equal_split) | Dev | 30+ test scenarios pass |
| Fri-14 | Performance test (1000 units allocation) | QA | <5s response ✅ |

**Deliverable:** `feature/pms-p0-metering-mixed` merged to dev

---

### Week 4: Integration & Go-Live Prep (Mar 17-23)

| Day | Task | Owner | Acceptance |
|-----|------|-------|-----------|
| Mon-17 | All P0 items merged to dev | Dev | No merge conflicts |
| Tue-18 | Full PMS test suite run | QA | go test ./internal/modules/pms/... ✅ |
| Wed-19 | Seed verification (menu/permission/translation/domain) | DevOps | All seeds load ✅ |
| Thu-20 | Frontend lint + portal pages | Frontend | npm run lint 0 errors ✅ |
| Fri-21 | UAT environment deployment | DevOps | PMS live on staging ✅ |

**Deliverable:** Sprint-14 complete, ready for Sprint-15

---

### Week 5-6: Sprint-15 P1 Features (Mar 24-Apr 6)

**Parallel streams:**

| Stream | Lead | Items | Est. |
|--------|------|-------|------|
| **Analytics** | Dashboard Dev | P1-001 Portfolio | 5 days |
| **Billing** | Billing Dev | P1-002 Reconciliation Scoring | 5 days |
| **Infrastructure** | DevOps | P1-003 Export Concurrency + P1-004 Webhook Monitoring | 7 days |
| **Portal** | Portal Dev | P1-005 Portal Rate Limiting | 3 days |

---

## 🔍 Kod Tarama Checklist

**Before Sprint-15 starts, ensure:**

- [ ] Read [PMS_MASTER.md](PMS_MASTER.md) completely
- [ ] Read [PMS_SPRINT_BACKLOG.md](PMS_SPRINT_BACKLOG.md) P0 section
- [ ] `go test ./internal/modules/pms/... ./cmd/api` passes (100%+ coverage)
- [ ] `npm run lint 'app/(dashboard)/admin/pms/**'` 0 errors
- [ ] Database indexes verified (TTL on export_files, webhook_deliveries, billing_idempotency)
- [ ] Environment variables documented (PMS_*) in `.env.example`
- [ ] All seeds load without error (menu, permission, role, translation, domain)
- [ ] Admin portal pages responsive ([/admin/pms/dashboard](http://localhost:3000/admin/pms/dashboard) accessible)
- [ ] Portal endpoints guarded (occupancy check enforced)
- [ ] Outbox events (15+ event types) flowing to subscribers

---

## 📊 Metrics & Monitoring

**Setup before production:**

```
# Backend metrics
go install github.com/prometheus/client_golang/cmd/promtool@latest
curl http://localhost:8080/metrics | grep pms_

# Examples:
pms_operations_total{operation="publish", status="success", module="billing"} 
pms_outbox_events_emitted_total{event_type="pms.billing.dues.published"}
pms_webhook_deliveries_total{status="sent"}
pms_export_jobs_total{type="accounting", status="done"}

# Frontend type checking
npm run type-check app/(dashboard)/admin/pms

# Frontend lint
npx eslint 'app/(dashboard)/admin/pms/**/*.{ts,tsx}' --fix
```

**SLO Baselines (MVP target):**
- Dues publish: P95 < 2s
- Reconciliation match: 100% deterministic
- Export job: P95 < 10s
- Webhook delivery (w/ retry): >99% success
- Portal uptime: >99.5%
- Payment allocation: 100% correctness (audit-verified)

---

## 🚨 Critical Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|-----------|
| **Payment webhook replay** | Double-charge | idempotency key guard + test scenario P0-002 |
| **Metering formula error** | Billing disputes | Manual QA verification + 40+ unit tests |
| **Access state machine bypass** | Security breach | Comprehensive edge-case tests P0-001 |
| **Portfolio aggregation slow** | Performance regression | Query optimization + <1s SLO verified |
| **Distributed lock collision** | Export data corruption | Stress test 100 concurrent jobs |
| **Webhook silent failure** | Missed integrations | Monitoring dashboard + retry insights |

**Daily standup:** Review P0-001/P0-002/P0-003 status; escalate blocks same-day.

---

## 📚 Reference Documents

**Core PMS Docs:**
1. [PMS_MASTER.md](PMS_MASTER.md) ← architecture + all 16 modules + API surface
2. [PMS_SPRINT_BACKLOG.md](PMS_SPRINT_BACKLOG.md) ← P0/P1/P2 items + timeline
3. [README_FOR_AI_ASSISTANTS.md](README_FOR_AI_ASSISTANTS.md) ← AYSTEK standards (critical!)

**Implementation References:**
- `go/aystek-core/internal/modules/pms/` ← all PMS code
- `node/aystek-core-next/app/(dashboard)/admin/pms/` ← admin UI
- `node/aystek-core-next/lib/api/pms.ts` ← API client

**Historical Context:**
- `aystek-core-dev-docs/done/TASK_037_PMS_ACCESS_VISITORS.md` ← most recent completed work
- `aystek-core-dev-docs/done/TASK_040_PMS_PAYMENT_WEBHOOK_AUTOFINALIZE.md` ← webhook orchestration

---

## 🤝 Handoff Checklist

**Before handing off to Copilot/QA:**

- [x] PMS_MASTER.md written (enterprise architecture)
- [x] PMS_SPRINT_BACKLOG.md written (P0/P1/P2 roadmap)
- [x] This Devam Planı written (week-by-week tasks)
- [ ] All P0 items understood by dev team
- [ ] Sprint-14.B & 14.C tasks assigned
- [ ] Test suite reviewed (coverage gaps plugged)
- [ ] Go-live readiness gate defined

**Kick-off for Sprint-14.B:** Monday Feb 25, 9:00 AM

---

## 💡 Key Insights

### What Went Well ✅

- **Modular design:** All PMS business logic isolated in `/internal/modules/pms/`
- **Outbox pattern:** Event-driven integration solid (15+ event types)
- **RBAC integration:** Role-based access well-established
- **Multi-tenancy:** Tenant + site scoping rigorously enforced
- **Test-first approach:** Most modules have service-level tests

### What Needs Attention 🔴

1. **Access Suite edge cases:** Double checkout/exit scenarios incomplete
2. **Payment webhook idempotency:** Needs verification testing (no e2e test yet)
3. **Metering mixed formula:** Formula correctness critical; needs 40+ test cases
4. **Portfolio performance:** Aggregation queries may slow on 100+ sites
5. **Webhook monitoring:** No admin UI for delivery status/retry history

### No-Go Items ❌

- **Do NOT:** Bypass tenant/site scoping (security-critical)
- **Do NOT:** Skip idempotency tests for payment flows (data integrity)
- **Do NOT:** Release without comprehensive metering test suite
- **Do NOT:** Deploy without Portfolio performance verification
- **Do NOT:** Go live without webhook monitoring dashboard

---

## 📞 Escalation Path

**If blocked on:**

| Blocker | Owner | Contact |
|---------|-------|---------|
| Webhook callback not triggering finalize | Payments team | Escalate ASAP (P0-002 critical) |
| Metering formula disagreement | Billing team | Whiteboard session required |
| Portal access guard failing | Security team | Review occupancy logic |
| Export concurrency deadlock | DevOps team | DB lock analysis required |
| i18n translation missing | i18n team | Generate 28-language keys |

---

## 🎯 Success = Sprint-15 Ready

By end of **Week 4 (Mar 23):**

✅ All P0 items merged & tested  
✅ PMS suite deployable to staging  
✅ Admin portal fully functional  
✅ Portal endpoints occupancy-gated  
✅ Outbox events flowing  
✅ All seeds loading  
✅ go test ./internal/modules/pms/... **100% passing**  
✅ npm run lint **0 errors**  
✅ Documentation complete (PMS_MASTER.md + backlog + this plan)  

**Then: Sprint-15 P1 features proceed with confidence.**

---

**END OF DEVAM PLANI**

*Principal Engineer Handoff | Ready for Implementation*
