# P1 Master Implementation Roadmap

**Amaç:** Tüm P1 özellikleri (G1-G4) sistematik, risksiz şekilde canlıya al  
**Toplam Efor:** 95 saat  
**İdeal Zaman:** 12 gün (8h/gün) veya 2.5 hafta (30h/hafta)  
**Hedef Tarih:** Şubat 28, 2026

---

## 1. Execution Strategy Overview

### ⚠️ Critical Constraints

1. **Zero Breaking Changes** — Existing P0 orders must work without changes
2. **Feature Flags** — Each P1 feature independently toggleable
3. **Backward Compatible** — Old data models must work with new code
4. **Database Migrations** — Must support rollback
5. **Production Testing** — Load test before go-live

### ✅ Advantages

1. **All code paths documented** — No guessing on implementation
2. **Dependencies mapped** — Know what to build when
3. **Tests pre-written** — Unit + integration scenarios ready
4. **Go backend mature** — P0 proof: 75+ tests passing, version control working

---

## 2. Feature Dependency Graph

```
P0 Foundation (100% complete)
├─ Order state machine ✅
├─ Version control ✅
├─ Webhook replay ✅
├─ Credit memo system ✅
├─ Refund engine ✅
└─ Shipping carriers ✅

P1-G1: Tax Rules (23h) ⬅️ CRITICAL PATH - MUST BE FIRST
├─ Depends: P0 order system, credit memo
├─ Blocks: P1-G2 (B2B), P1-G4 (returns tax)
└─ Can parallel: P1-G3 (stock)

P1-G3: Stock Allocation (28h) ⬅️ INDEPENDENT
├─ Depends: P0 order system only
├─ Blocks: P1-G2 (B2B uses allocation)
└─ Can parallel: P1-G1 (tax), but finish before G2

P1-G2: B2B Workflows (24h) ⬅️ DEPENDENT ON G1 + G3
├─ Depends: P1-G1 (tax), P1-G3 (stock allocation)
├─ PO model (independent)
├─ Approval workflow (independent)
└─ Payment terms (tax + refund dependent)

P1-G4: Returns/RMA (20h) ⬅️ PARTIALLY DEPENDENT
├─ Depends: P0 invoices, P0 refunds, P1-G1 (tax on credit)
├─ Restock logic: P1-G3 (stock)
└─ Can start: Before G2, after G1
```

---

## 3. Recommended Execution Sequence

### Sprint 1: Tax Rules Foundation (Days 1-3, 23 hours)

**Why First?**
- Tax affects ALL orders (P1-G2, P1-G4)
- Simplest to test in isolation
- Must be production-ready before other features

**Phase 1.1: Checkout Concurrency (Day 1, 6h)**

| Time | Task | Files | Owner |
|------|------|-------|-------|
| 0-1h | Add Version field to CheckoutSession | `checkout/model/checkout_session.go` | Dev 1 |
| 1-2h | UpdateWithVersion() in repository | `checkout/repository/mongo.go` | Dev 1 |
| 2-3h | Modify CalculateTotals() for concurrency | `checkout/service/service.go` | Dev 1 |
| 3-4h | Fix discount application atomicity | `checkout/service/service.go` | Dev 2 |
| 4-5h | Unit tests (concurrent scenarios) | `checkout/service/*_test.go` | Dev 2 |
| 5-6h | Integration test + load test (100 concurrent) | `checkout/service/*_test.go` | QA |

**Deliverable:** ✅ Checkout race conditions fixed, 100 concurrent orders < 5s

**Phase 1.2: EU Reverse Charge (Day 2, 8h)**

| Time | Task | Files | Owner |
|------|------|-------|-------|
| 0-1h | Add EU config fields to tax model | `taxes/model/tax_rule.go` | Dev 1 |
| 1-2h | Company EU registration fields | `b2b/model/company.go` | Dev 1 |
| 2-3h | Invoice EU reverse charge fields | `billing/invoices/model/invoice.go` | Dev 2 |
| 3-4h | Tax service: EU reverse charge logic | `taxes/service/service.go` | Dev 2 |
| 4-5h | E-invoice XML generation | `billing/invoices/service/e_invoice.go` | Dev 1 |
| 5-6h | Unit tests (EU scenarios) | `taxes/service/*_test.go` | Dev 2 |
| 6-7h | Integration test (Cross-border B2B) | `taxes/service/*_test.go` | QA |
| 7-8h | Compliance reporting API | `taxes/handler/reporting_handler.go` | Dev 1 |

**Deliverable:** ✅ EU reverse charge working, compliance reports ready

**Phase 1.3: Tax Brackets (Day 3, 4h)**

| Time | Task | Files | Owner |
|------|------|-------|-------|
| 0-1h | Tax bracket model | `taxes/model/tax_bracket.go` | Dev 1 |
| 1-2h | Bracket selection logic | `taxes/service/bracket_selection.go` | Dev 2 |
| 2-3h | Unit tests | `taxes/service/*_test.go` | Dev 2 |
| 3-4h | Integration test + load test | `taxes/service/*_test.go` | QA |

**Deliverable:** ✅ Tax brackets configured, tested, 1000 calc/s performance

**End of Sprint 1 Checklist:**
- [ ] All P1-G1 code committed + reviewed
- [ ] 20+ unit tests passing
- [ ] 5 integration scenarios passing
- [ ] Load test: 100 concurrent orders < 5s
- [ ] Documentation updated
- [ ] Feature flag: `TAX_RULES_V2` (default: false for safety)

---

### Sprint 2A: Stock Allocation (Days 4-6, 28 hours)

**Can Run Parallel to Sprint 1** — Fully independent

| Day | Phase | Hours | Deliverable |
|-----|-------|-------|-------------|
| 4 | 1.1-1.4: Backorder system | 6h | BO model + CRUD + repository |
| 5 | 2.1-2.4: Allocation strategies | 8h | 3 strategies implemented, API working |
| 6 | 3.1-3.4: Safety stock + reorder | 6h | Reorder scheduler working |
| 7 | 4.1-4.4: Testing + deployment | 8h | All tests passing, ready for prod |

**End of Sprint 2A Checklist:**
- [ ] Backorder system tested (FIFO fulfillment)
- [ ] 3 allocation strategies working
- [ ] Multi-warehouse allocation < 100ms per order
- [ ] Safety stock check working (6h interval)
- [ ] Feature flag: `STOCK_ALLOCATION_V2` (default: false)

---

### Sprint 2B: Returns/RMA (Days 4-5, 20 hours)

**Can Run Parallel to Sprint 1 + 2A**

| Day | Phase | Hours | Deliverable |
|-----|-------|-------|-------------|
| 4 | 1.1-1.4: RMA model + CRUD | 5h | RMA basic workflow |
| 5 | 2.1-2.3: Inspection system | 5h | Reception + condition assessment |
| 6 | 3.1-3.4: Refund + restock | 5h | Credit memo + refund integration |
| 7 | 4.1-4.4: Testing + deployment | 5h | All tests passing |

**End of Sprint 2B Checklist:**
- [ ] RMA workflow: Request → Authorize → Receive → Inspect → Process
- [ ] Refund calculations by reason (100%, 90%, 50%)
- [ ] Restock logic working
- [ ] 5 integration scenarios passing
- [ ] Feature flag: `RETURNS_RMA_V1` (default: false)

---

### Sprint 3: B2B Workflows (Days 8-10, 24 hours)

**Must Wait For:** P1-G1 (tax) ✅ + P1-G3 (stock) ✅

| Day | Phase | Hours | Deliverable |
|-----|-------|-------|-------------|
| 8 | 1.1-1.4: PO model + CRUD | 8h | Purchase order system |
| 9 | 2.1-2.4: Approval workflow | 8h | Multi-level approvals working |
| 10 | 3.1-3.3: Payment terms + invoice | 8h | Payment terms integrated |

**Note:** G2 starts DAY 8 (after G1 + G3 complete on Day 7)

**End of Sprint 3 Checklist:**
- [ ] PO creation → approval chain → order conversion working
- [ ] Multi-level approvals (manager > credit > executive)
- [ ] Payment terms (Net 30, 2/10 Net 30, deposits)
- [ ] Credit reservation working
- [ ] 3 integration scenarios passing
- [ ] Feature flag: `B2B_WORKFLOWS_V1` (default: false)

---

## 4. Parallel Execution Timeline

```
Week 1:
├─ Days 1-3: P1-G1 Tax Rules (SEQUENTIAL, must complete)
├─ Days 1-7: P1-G3 Stock (parallel to G1, can overlap)
└─ Days 1-5: P1-G4 Returns (parallel to G1+G3, can overlap)

Week 2:
├─ Day 8-10: P1-G2 B2B (AFTER G1+G3 complete)
└─ Days 8-12: Testing + deployment prep
```

**Reality Timeline:**
- If 2 teams (3 devs + 1 QA) → 12 days possible
- If 1 team → 18-20 days (sequential)
- Recommended: Sequential within features, parallel between independent features

---

## 5. Code Review & Quality Gate

### Before Each Day's Work

```bash
# 1. Pull latest P0 tests to verify nothing broke
go test ./internal/modules/orders/... -v

# 2. Run P0 payment workflow
go test ./internal/modules/payments/... -v

# 3. Load test baseline
go test -bench=BenchmarkOrderCreation -benchmem
```

### Before Each PR Merge

**Checklist:**
- [ ] All tests passing (`go test ./...`)
- [ ] Coverage > 80% for new code
- [ ] No breaking changes to P0 APIs
- [ ] Database migrations reversible
- [ ] Feature flag implemented (OFF by default)
- [ ] Rollback procedure documented
- [ ] 2+ reviewers approval
- [ ] Load test passed (if performance-critical)

### Code Review Workflow

```
Developer → Create feature branch (e.g., feature/p1-g1-tax-phase1)
           → Implement Phase 1 tasks
           → Run local tests
           → Create PR (reference plan doc)
           → Automated tests run
           → 2+ manual reviews
           → Merge to develop
           → Deploy to staging
           → QA validation
           → Tag release (e.g., v1.2.0-p1g1-phase1)
           → Production deployment
```

---

## 6. Database Migrations Strategy

### Deployment Pattern

```go
// 1. Add new fields as OPTIONAL (nullable in MongoDB)
type Invoice struct {
    // OLD
    Amount float64
    
    // NEW (optional, backward compatible)
    PaymentTermsID *primitive.ObjectID `bson:"payment_terms_id,omitempty"`
}

// 2. Code handles missing fields
if invoice.PaymentTermsID == nil {
    // Use default behavior (P0)
} else {
    // Use new behavior (P1)
}

// 3. After 1 week production stability → backfill data
db.invoices.updateMany({paymentTermsID: null}, {$set: {paymentTermsID: ObjectId(...)}})

// 4. Make field mandatory in next major version
```

### Rollback Procedure

```bash
# If critical issue found in prod:
1. Disable feature flag: TAX_RULES_V2=false
2. Redeploy previous binary
3. All new orders fall back to P0 tax logic
4. No data loss, no migration rollback needed
```

---

## 7. Feature Flags Configuration

**Environment Variables:**

```bash
# Tax Rules V2 (default: disabled for safety)
P1_G1_CHECKOUT_CONCURRENCY=false      # Day 1: Enable after testing
P1_G1_EU_REVERSE_CHARGE=false         # Day 2: Enable after EU compliance check
P1_G1_TAX_BRACKETS=false              # Day 3: Enable after bracket config

# Stock Allocation V2
P1_G3_BACKORDER_SYSTEM=false          # Day 4: Enable for returns + B2B
P1_G3_ALLOCATION_STRATEGIES=false     # Day 5: Enable multi-warehouse
P1_G3_SAFETY_STOCK=false              # Day 6: Enable reorder automation

# Returns/RMA V1
P1_G4_RMA_WORKFLOW=false              # Day 4: Enable RMA system
P1_G4_INSPECTION_SYSTEM=false         # Day 5: Enable warehouse inspection

# B2B Workflows V1
P1_G2_PURCHASE_ORDERS=false           # Day 8: Enable after G1+G3
P1_G2_APPROVAL_WORKFLOW=false         # Day 9: Enable approvals
P1_G2_PAYMENT_TERMS=false             # Day 10: Enable payment terms
```

**Enabling Strategy:**

```
Day 1 (Prod Deploy):  All flags OFF (safe, proves P0 still works)
Day 2-7 (Staging):    Gradually enable each flag, run integration tests
Week 2 (Early Adopters): Enable for specific B2B customers only
Week 3 (Full Rollout):  Enable globally with monitoring
```

---

## 8. Monitoring & Observability

### Prometheus Metrics to Add

```go
// P1-G1: Tax calculations
tax_calculation_duration_seconds
tax_calculation_errors_total
tax_reverse_charge_detected_total
tax_brackets_used_by_tier

// P1-G3: Stock allocation
allocation_strategy_choice     // Which strategy used?
backorder_created_total
allocation_latency_percentiles (p50, p95, p99)
warehouse_fragmentation_ratio

// P1-G2: B2B
po_approval_chain_length
po_authorization_time_seconds
payment_terms_adoption_rate

// P1-G4: Returns
rma_authorization_rate
rma_processing_time_days
refund_approval_rate
return_condition_distribution
```

### Alerting Thresholds

```yaml
alerts:
  - name: TaxCalculationLatency
    threshold: p99 > 100ms
    severity: warning
    
  - name: AllocationFailureRate
    threshold: > 5%
    severity: critical
    
  - name: RMABacklog
    threshold: > 100 open for > 7 days
    severity: warning
    
  - name: B2BApprovalStalled
    threshold: > 5 awaiting approval > 24h
    severity: warning
```

---

## 9. Risk Assessment & Mitigation

### P1-G1 Tax (HIGH RISK - Financial)

**Risk:** Incorrect tax calculation → compliance issue + revenue impact

**Mitigation:**
1. ✅ Phase 1: Checkout concurrency (zero financial impact)
2. ✅ Phase 2: EU reverse charge (specific market only, togglable)
3. ✅ Phase 3: Brackets (verified against tax authority rates)
4. ✅ Testing: 5 scenarios covering 10 countries/states
5. ✅ Load test: 1000 tax calcs/sec verified
6. ✅ Audit: All calculations logged for compliance review
7. ✅ Rollback: Feature flag allows instant revert to P0

**Decision Gate:** Finance approval before enabling EU reverse charge

---

### P1-G3 Stock (MEDIUM RISK - Operations)

**Risk:** Backorder mismanagement → customer satisfaction impact

**Mitigation:**
1. ✅ Phase 1: Backorder system manual (humans decide)
2. ✅ Phase 2: Allocation strategies (A/B test each)
3. ✅ Phase 3: Reorder automation (configurable thresholds)
4. ✅ Testing: Concurrent order + stock race scenarios
5. ✅ Monitoring: Allocation failure rate dashboard
6. ✅ Rollback: Fall back to "first warehouse" if issues

**Decision Gate:** Warehouse manager sign-off before automation

---

### P1-G2 B2B (MEDIUM RISK - Process)

**Risk:** PO approval delays → B2B customer frustration

**Mitigation:**
1. ✅ PO model + CRUD (low risk)
2. ✅ Approval workflow (test with subset of B2B customers first)
3. ✅ Payment terms (validate against existing invoices)
4. ✅ Testing: Full approval chain scenarios
5. ✅ Monitoring: PO approval SLA dashboard
6. ✅ Rollback: Disable B2B workflows, orders go back to B2C

**Decision Gate:** B2B customer pilot feedback before full rollout

---

### P1-G4 Returns (LOW RISK - New Feature)

**Risk:** RMA abuse → higher refund rates

**Mitigation:**
1. ✅ RMA model + authorization
2. ✅ Return window enforcement (30 days default)
3. ✅ Condition assessment + refund adjustment
4. ✅ Testing: Return reason distribution analysis
5. ✅ Monitoring: RMA aging, refund rate by reason
6. ✅ Rollback: Feature flag disables RMA creation

**Decision Gate:** 1 week monitoring after launch

---

## 10. Success Criteria

### Phase Completion

**P1-G1 Complete When:**
- [ ] All 3 phases code-complete + merged
- [ ] 20+ unit tests passing
- [ ] 5 integration scenarios passing
- [ ] Load test: 100 concurrent orders < 5s
- [ ] 2 manual reviews + approved
- [ ] Production staging: Feature flag OFF (verified P0 works)

**P1-G3 Complete When:**
- [ ] All 4 phases code-complete + merged
- [ ] Backorder FIFO tested
- [ ] 3 allocation strategies tested
- [ ] Reorder check working (6h interval)
- [ ] Multi-warehouse allocation < 100ms
- [ ] Staging: Basic backorder flow tested

**P1-G2 Complete When:**
- [ ] All 3 phases code-complete + merged
- [ ] PO workflow: Create → Approve → Order
- [ ] Multi-level approvals tested
- [ ] Payment terms integration tested
- [ ] Staging: End-to-end B2B flow tested
- [ ] Feature flag OFF (P0 B2C still works)

**P1-G4 Complete When:**
- [ ] All 4 phases code-complete + merged
- [ ] RMA workflow: Request → Process
- [ ] Refund calculations by reason verified
- [ ] Restock logic tested
- [ ] Staging: Happy path return tested
- [ ] Feature flag OFF (P0 no returns yet)

### Go-Live Readiness

**Day 12 Checkpoint:**

```
✅ P0 Tests:       75+ passing (zero regressions)
✅ P1-G1 Tests:    20+ passing
✅ P1-G2 Tests:    15+ passing
✅ P1-G3 Tests:    18+ passing
✅ P1-G4 Tests:    12+ passing
---
✅ Total:          140+ tests passing

✅ Code Review:    2+ approvals per PR
✅ Load Test:      100 concurrent orders < 5s
✅ Staging:        All feature flows verified
✅ Feature Flags:  All OFF (P0 baseline working)
✅ Monitoring:     Dashboards + alerts configured
✅ Rollback Plan:  Documented + tested

🟢 PRODUCTION READY
```

---

## 11. Deployment Sequence

### Day 12: Staging Deployment (1h)

```bash
# 1. Build with all P1 code
git checkout develop
go build -o aystek-core ./cmd/api

# 2. Run P0 tests (must pass)
go test ./internal/modules/orders/... -v
go test ./internal/modules/payments/... -v
go test ./internal/modules/billing/... -v
go test ./internal/modules/shipping/... -v

# 3. Deploy to staging with feature flags OFF
docker push aystek:v1.2.0-p1-ready
helm upgrade --install aystek ./charts/aystek \
  -f values-staging.yaml \
  --set env.P1_G1_CHECKOUT_CONCURRENCY=false \
  --set env.P1_G3_ALLOCATION_STRATEGIES=false \
  --set env.P1_G2_PURCHASE_ORDERS=false \
  --set env.P1_G4_RMA_WORKFLOW=false

# 4. QA: Verify P0 orders work (smoke test)
curl -X POST http://staging-api/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"items": [{"product_id": "...", "quantity": 1}]}'
# Expected: Order created successfully

# 5. Monitoring: Check error rates, latency
# (All metrics should show P0 baseline)
```

### Days 13-19: Feature Rollout (staged)

```
Day 13: Enable P1-G1 (tax) for 1% of orders
        → Monitor for 24h
        → If OK: 25%
        
Day 14: 50% P1-G1 + Enable P1-G3 (stock) for 5%
        → Monitor for 24h
        → If OK: 25% stock
        
Day 15: 75% P1-G1 + 25% P1-G3 + Enable P1-G4 (returns) for early adopters
        → Monitor for 48h
        
Day 16: 100% P1-G1 + 50% P1-G3 + Enable P1-G2 (B2B) for pilot customers
        → Monitor for 72h
        → Customer feedback collection
        
Day 17-19: Full rollout if metrics green
           All flags → 100% enabled
           Monitor for issues
```

### Production Monitoring Dashboard

```
Real-time Metrics:
├─ Orders/min (should match P0 baseline ±5%)
├─ Tax calc latency (p99 < 50ms)
├─ Allocation success rate (>99%)
├─ RMA creation rate (expect < 2% of orders)
├─ B2B PO approval time (SLA: <4h)
├─ Error rate (should be < 0.1%)
└─ Database query times (p99 < 100ms)

Alerts (page on-call immediately):
├─ Order creation failing > 1%
├─ Tax calculation errors > 0.5%
├─ Allocation failures > 1%
├─ Unusual refund spike (>5x baseline)
└─ Database latency > 500ms
```

---

## 12. Communication Plan

### Stakeholder Updates

**Finance Team:**
- [ ] Day 7: Tax rules ready for compliance review
- [ ] Day 8: EU reverse charge implementation details
- [ ] Day 12: Staging deployment complete, ready for Day 13 prod rollout

**Warehouse Team:**
- [ ] Day 5: Stock allocation strategy overview
- [ ] Day 7: Reorder automation testing schedule
- [ ] Day 10: Returns inspection workflow walkthrough

**B2B Sales Team:**
- [ ] Day 8: PO workflow ready for customer pilot
- [ ] Day 10: Payment terms configuration options
- [ ] Day 16: B2B feature live for select customers

**Customer Support:**
- [ ] Day 1: RMA feature launching (training provided)
- [ ] Day 12: Staging deployment (internal testing)
- [ ] Day 13: P1 features rolling out (monitor support tickets)
- [ ] Day 19: Full rollout complete (monitoring continues)

---

## 13. Document Requirements

### By End of Sprint 1 (Day 7)

- [ ] Tax Rules Implementation Guide (admin config)
- [ ] EU Reverse Charge Compliance Report
- [ ] Tax Bracket Configuration Manual
- [ ] Rollback Procedure (1-page checklist)

### By End of Sprint 2 (Day 11)

- [ ] Stock Allocation Strategy Guide
- [ ] Backorder Management Manual
- [ ] RMA Workflow Documentation (customer-facing)
- [ ] Returns Inspection Procedure (warehouse)
- [ ] B2B Purchase Order Guide (for customers)
- [ ] Approval Workflow Configuration (for admins)

### Final (Day 12)

- [ ] Production Deployment Runbook
- [ ] Monitoring Dashboard Setup Guide
- [ ] Incident Response Playbook (if issues)
- [ ] Performance Baseline Report (P0 vs P1)

---

## Summary

**P1 Implementation Master Timeline:**

| Week | Focus | Hours | Status |
|------|-------|-------|--------|
| **1** | P1-G1 Tax (primary) + G3 Stock (parallel) + G4 Returns (parallel) | 51h | Sequential then parallel |
| **2** | P1-G2 B2B (after G1+G3) + Testing/Deployment | 44h | Dependent start |
| **Total** | All P1 features production-ready | **95h** | ✅ Ready |

**Next Action:**
1. Approve this roadmap
2. Assign developers to sprints
3. Start Day 1: P1-G1 Phase 1 (Checkout concurrency fixes)
4. Run P0 regression test at start of each day
5. Daily standup: 15min status + blockers

**Risks Mitigated:**
- ✅ Zero breaking changes (feature flags)
- ✅ Incremental delivery (3 days → 1 week → 2 weeks)
- ✅ Rollback ready (feature flags + DB migration strategy)
- ✅ Quality gates (tests + code review + load testing)
- ✅ Monitoring dashboards (track every metric)

**Success Metric:** Day 19 → All P1 features 100% enabled, 0% error rate increase vs P0
