# TASK-015: PMS Sprint-15 Production Hardening Checklist

**Date:** February 21, 2026  
**Sprint:** Sprint-15  
**Branch:** `feature/pms-s15-hardening`  
**Owner:** Production Hardening Engineer  
**Status:** 🔄 IN PROGRESS

---

## 📋 EXECUTIVE SUMMARY

Sprint-15 hardening focuses on **production readiness audit** for AYSTEK PMS. Building on the metering hardening foundation (PMS_PRODUCTION_HARDENING_REPORT.md), this sprint extends idempotency guarantees, state machine standardization, and concurrency safety across all PMS modules.

**Scope:** `internal/modules/pms/**` + Next.js PMS screens + seeds + docs  
**Out of Scope:** Non-PMS modules (auth, tenant, ads, CMS, etc.)

---

## 🎯 P0: RELEASE GATE REQUIREMENTS

### P0-1: Data Integrity & Idempotency Audit ✅

**Modules Requiring Idempotency:**

| Module | Critical Path | Idempotency Mechanism | Status |
|--------|--------------|----------------------|--------|
| **Billing** | Charge Publish | `pms_billing_idempotency` (TTL 30d) | ✅ EXISTS |
| **Billing** | Payment Allocate | Request hash dedup | ✅ EXISTS (P0-002) |
| **Metering** | Allocation Publish | Unique index + `idempotency_key` | ✅ ENHANCED (S14) |
| **Reconciliation** | Confirm Match | `external_txn_id` unique per tenant | ⚠️ VERIFY |
| **Portal Payments** | Finalize → Allocate | Multi-layer nonce + intent key | ✅ EXISTS (P0-002) |
| **Exports** | Job Finalize | Atomic claim via `lock_id` | ⚠️ VERIFY |
| **Webhooks** | Delivery Retry | Delivery attempt dedup | ⚠️ VERIFY |

**Test Requirements:**
- [ ] **Replay Protection Test:** Same idempotency key + same payload → cached response (no new DB write)
- [ ] **Conflict Detection Test:** Same idempotency key + different payload → 409 Conflict
- [ ] **Concurrency Test:** 50 concurrent requests with same key → only 1 DB write

**Files to Audit:**
- `internal/modules/pms/billing/service/service.go` → PublishCharges, AllocatePayment
- `internal/modules/pms/reconciliation/service/service.go` → ConfirmMatch
- `internal/modules/pms/exports/service/service.go` → FinalizeJob
- `internal/modules/pms/portal/service/service.go` → FinalizePayment

---

### P0-2: State Machine & 409 Semantics Standardization ⚠️

**Current State:**
- **Access Module:** ✅ Implements proper 409 handling (visitor/vehicle already exited)
- **Metering Module:** ⚠️ Has `CanTransitionTo()` but no 409 mapping in handler
- **Work Orders:** ⚠️ State machine exists, 409 handling unclear
- **Reconciliation:** ⚠️ State transitions exist, 409 handling unclear

**Target Standard (Based on Access Module):**

```go
// PATTERN: Repository layer returns domain error
var (
    ErrInvalidStateTransition = errors.New("invalid state transition")
    ErrAlreadyPublished = errors.New("allocation already published")
    ErrAlreadyConfirmed = errors.New("reconciliation already confirmed")
    // ... per module
)

// PATTERN: Handler maps to 409
func handleErr(c *gin.Context, err error) {
    if errors.Is(err, service.ErrInvalidStateTransition) {
        response.Error(c, 409, "pms.error.invalid_state_transition", nil)
        return
    }
    if errors.Is(err, service.ErrAlreadyPublished) {
        response.Error(c, 409, "pms.metering.error.already_published", nil)
        return
    }
    // ...
}
```

**Modules to Standardize:**

| Module | State Machine | Current 409 Handling | Action Needed |
|--------|--------------|---------------------|---------------|
| **Access** | in → out | ✅ DONE (ErrVisitorAlreadyCheckedOut) | None |
| **Metering** | draft → previewed → published → archived | ⚠️ PARTIAL (has CanTransitionTo) | Add handler 409 mapping |
| **Work Orders** | open → in_progress → done/cancelled | ❌ MISSING | Add domain errors + 409 mapping |
| **Reconciliation** | new → matched → confirmed/rejected | ❌ MISSING | Add domain errors + 409 mapping |
| **Billing** | open → allocated_full → paid | ⚠️ PARTIAL | Verify 409 handling |
| **Webhooks** | pending → processing → sent/failed | ⚠️ PARTIAL | Verify 409 handling |

**Implementation Tasks:**
- [ ] Define domain errors in `service/errors.go` for each module
- [ ] Add state transition validation before DB writes (use existing `CanTransitionTo()` methods)
- [ ] Map domain errors to 409 in handlers
- [ ] Write unit tests: `TestInvalidTransitionReturns409`
- [ ] Update i18n error keys: `pms.<module>.error.invalid_transition`

**Test Pattern:**
```go
func TestPublishAllocation_InvalidTransition_Returns409(t *testing.T) {
    // Setup: allocation already published
    alloc := &model.UtilityAllocation{Status: model.AllocationStatusPublished}
    repo.allocations = []model.UtilityAllocation{*alloc}
    
    // Attempt: publish again
    _, err := svc.PublishAllocation(ctx, input)
    
    // Assert: error is ErrAlreadyPublished
    assert.ErrorIs(t, err, service.ErrAlreadyPublished)
}
```

---

### P0-3: Transaction & Partial Failure Tolerance 🔍

**Critical Orchestration Chains:**

#### Chain 1: Metering Publish → Billing Charge Creation
```
1. repo.InsertAllocation(allocation)        ← Write 1
2. for each unit:
     billingSvc.CreateManualCharge(...)      ← Write 2..N
3. repo.SetAllocationPublished(...)         ← Write N+1
4. repo.InsertOutboxEvent(...)              ← Write N+2
```

**Risk:** Step 2 fails partway → allocation exists but charges incomplete  
**Current Mitigation:** Idempotency allows safe retry (republish creates missing charges)  
**Enhancement Needed:**
- [ ] Document compensation pattern in code comments
- [ ] Add test: `TestPublishAllocation_PartialChargeFailure_RetrySucceeds`
- [ ] Verify: Republish with same `idempotency_key` completes missing charges without duplicates

---

#### Chain 2: Reconciliation Confirm → Payment Create → Allocate
```
1. paymentSvc.CreatePayment(...)            ← Write 1 (core finance)
2. billingSvc.AllocatePayment(...)          ← Write 2 (PMS bridge)
3. repo.UpdateTransactionStatus(confirmed)  ← Write 3
4. repo.InsertOutboxEvent(...)              ← Write 4
```

**Risk:** Step 2 fails → payment created but not allocated  
**Current Mitigation:** ⚠️ VERIFY - is ConfirmMatch idempotent?  
**Enhancement Needed:**
- [ ] Verify: ConfirmMatch checks if already confirmed (409 guard)
- [ ] Add test: `TestConfirmMatch_AlreadyConfirmed_Returns409`
- [ ] Add test: `TestConfirmMatch_AllocationFailure_PaymentStillCreated_RetryAllocates`

---

#### Chain 3: Portal Payment Finalize → Allocate
```
1. Check nonce (5min window)                 ← Idempotency Layer 1
2. paymentIntentSvc.Finalize(idem_key)       ← Idempotency Layer 2
3. Check allocation hash                     ← Idempotency Layer 3
4. billingSvc.AllocatePayment(...)           ← Write
5. repo.InsertOutboxEvent(...)               ← Write
```

**Status:** ✅ VERIFIED in P0-002 (webhook auto-finalize)  
**Action:** None (already hardened)

---

#### Chain 4: Export Job Finalize
```
1. Claim job (atomic update with lock_id)    ← Claim
2. Generate file                             ← Compute
3. Upload to storage                         ← External I/O
4. repo.SetJobCompleted(...)                 ← Write
5. repo.InsertOutboxEvent(...)               ← Write
```

**Risk:** Step 3 fails → job stuck in "running", lock expires  
**Current Mitigation:** Lock expiration allows retry (stale lock recovery)  
**Enhancement Needed:**
- [ ] Verify: Lock expiration logic exists
- [ ] Add test: `TestExportJob_StaleLock_RecoverySucceeds`
- [ ] Add test: `TestExportJob_50Concurrent_OnlyOneSucceeds`

---

**Deliverables:**
- [ ] Add code comments documenting compensation patterns
- [ ] Add 8 partial failure tolerance tests (2 per chain)
- [ ] Update TASK_015_SUMMARY.md with transaction boundary analysis

---

## 🎯 P1: OPERATIONAL READINESS

### P1-1: Observability Enhancement 📊

**Current State (Metering Module):**
✅ 7 Prometheus metrics (preview_duration, publish_duration, failures, republish, etc.)  
✅ PII-safe labels (no tenant_id)  
✅ Low cardinality

**Target: Extend to All Hot Paths**

| Module | Hot Path | Metric Name | Labels | Status |
|--------|----------|------------|--------|--------|
| **Billing** | Publish Charges | `pms_billing_publish_duration_seconds` | method | ⚠️ ADD |
| **Billing** | Allocate Payment | `pms_billing_allocate_duration_seconds` | - | ⚠️ ADD |
| **Reconciliation** | Confirm Match | `pms_reconciliation_confirm_duration_seconds` | - | ⚠️ ADD |
| **Exports** | Job Processing | `pms_exports_job_duration_seconds` | format | ⚠️ ADD |
| **Webhooks** | Delivery Attempt | `pms_webhooks_delivery_duration_seconds` | status | ⚠️ ADD |
| **Work Orders** | Status Change | `pms_work_orders_status_change_total` | from, to | ⚠️ ADD |

**Request-ID Propagation:**
- [ ] Verify: Middleware injects `X-Request-ID` header
- [ ] Verify: Context carries request-id through service calls
- [ ] Verify: Outbox events include request-id for tracing

**Log Standards:**
- [ ] All service methods log entry/exit with request-id
- [ ] Errors logged with full context (tenant, site, entity_id, error)
- [ ] Sensitive data (passwords, tokens) never logged

---

### P1-2: Worker Concurrency & Lock Hardening 🔒

**Webhooks Delivery Worker:**

Current Pattern:
```go
// Claim deliveries WHERE status=pending AND next_retry_at <= NOW
deliveries := repo.ListPendingDeliveries(ctx, limit: 10)
for _, d := range deliveries {
    repo.UpdateStatus(d.ID, "processing")  ← Race condition?
    deliver(d)
    repo.UpdateStatus(d.ID, "sent")
}
```

**Risk:** Multiple workers claim same delivery  
**Fix:** Atomic claim pattern (like exports)

```go
// PATTERN: Atomic claim via UpdateOne with lock_id
lockID := uuid.New()
result := repo.ClaimDelivery(ctx, deliveryID, lockID, lockExpiry)
if result.ModifiedCount == 0 {
    return nil // Another worker claimed it
}
// Proceed with delivery...
```

**Tasks:**
- [ ] Implement atomic claim for webhook deliveries
- [ ] Add lock expiration recovery (stuck worker detection)
- [ ] Add test: `TestWebhookDelivery_100Concurrent_NoDoubleDelivery`
- [ ] Add metric: `pms_webhooks_claim_conflicts_total`

---

**Exports Runner:**

Status: ⚠️ VERIFY - does atomic claim exist?

**Tasks:**
- [ ] Review `exports/service/service.go` → ProcessJob method
- [ ] Verify: Lock-based claim with expiration
- [ ] Add test: `TestExportJob_StaleLockRecovery`
- [ ] Add metric: `pms_exports_stale_locks_recovered_total`

---

### P1-3: Retention & TTL Configuration 📅

**Current TTL Indexes:**

| Collection | Index | TTL | Status |
|-----------|-------|-----|--------|
| `pms_outbox` | `created_at` | 30 days | ✅ EXISTS (S14) |
| `pms_billing_idempotency` | `created_at` | 30 days | ✅ EXISTS |
| `pms_webhook_deliveries` | ❓ | ❓ | ⚠️ VERIFY |
| `pms_export_jobs` | ❓ | ❓ | ⚠️ VERIFY |
| `pms_work_order_events` | ❓ | ❓ | ⚠️ VERIFY |

**Environment Variables (Check `.env.example`):**
```bash
PMS_OUTBOX_TTL_DAYS=30
PMS_IDEMPOTENCY_TTL_DAYS=30
PMS_WEBHOOK_DELIVERY_TTL_DAYS=90    # To be added?
PMS_EXPORT_JOB_TTL_DAYS=90          # To be added?
```

**Tasks:**
- [ ] Verify: All TTL env vars documented in `.env.example`
- [ ] Add TTL indexes where missing (webhook_deliveries, export_jobs)
- [ ] Update seed files to include TTL configurations
- [ ] Document TTL policy in TASK_015_SUMMARY.md

---

## 📊 TESTING STRATEGY

### Idempotency Proof Tests (Per Critical Path)

**Template:**
```go
func TestXXX_Idempotency_SameKey_SamePayload_Cached(t *testing.T) {
    // First call: creates entity
    result1, err := svc.XXX(ctx, input)
    assert.NoError(t, err)
    
    // Second call: same idempotency_key + same payload
    result2, err := svc.XXX(ctx, input)
    assert.NoError(t, err)
    
    // Verify: result2 is cached (no new DB write)
    assert.Equal(t, result1.ID, result2.ID)
    assert.Equal(t, 1, len(repo.entities)) // Only 1 entity created
}

func TestXXX_Idempotency_SameKey_DifferentPayload_409(t *testing.T) {
    // First call: creates entity with payload A
    input1 := Input{IdempotencyKey: "key-123", Amount: 100}
    _, err := svc.XXX(ctx, input1)
    assert.NoError(t, err)
    
    // Second call: same key but different payload B
    input2 := Input{IdempotencyKey: "key-123", Amount: 200}
    _, err = svc.XXX(ctx, input2)
    
    // Verify: 409 Conflict
    assert.ErrorIs(t, err, service.ErrIdempotencyConflict)
}
```

**Tests to Add:**
- [ ] `TestPublishCharges_Idempotency` (billing)
- [ ] `TestAllocatePayment_Idempotency` (billing)
- [ ] `TestConfirmMatch_Idempotency` (reconciliation)
- [ ] `TestFinalizeJob_Idempotency` (exports)

---

### State Machine 409 Tests (Per Module)

**Template:**
```go
func TestXXX_InvalidTransition_Returns409(t *testing.T) {
    // Setup: entity in terminal state
    entity := &model.XXX{Status: "done"}
    repo.entities = []model.XXX{*entity}
    
    // Attempt: invalid transition
    err := svc.ChangeStatus(ctx, entity.ID, "in_progress")
    
    // Verify: domain error
    assert.ErrorIs(t, err, service.ErrInvalidStateTransition)
    
    // Verify: HTTP mapping in handler test
    rec := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(rec)
    handler.ChangeStatus(c)
    assert.Equal(t, 409, rec.Code)
}
```

**Tests to Add:**
- [ ] `TestPublishAllocation_AlreadyPublished_409` (metering)
- [ ] `TestChangeWorkOrderStatus_InvalidTransition_409` (work_orders)
- [ ] `TestConfirmMatch_AlreadyConfirmed_409` (reconciliation)

---

### Concurrency Stress Tests (50-100 goroutines)

**Template:**
```go
func TestXXX_50Concurrent_Idempotency(t *testing.T) {
    var wg sync.WaitGroup
    var successCount int32
    
    for i := 0; i < 50; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            _, err := svc.XXX(ctx, input) // Same idempotency key
            if err == nil {
                atomic.AddInt32(&successCount, 1)
            }
        }()
    }
    
    wg.Wait()
    
    assert.Equal(t, int32(50), successCount) // All succeed (cached)
    assert.Equal(t, 1, len(repo.entities))   // Only 1 entity created
}
```

**Tests to Add:**
- [ ] `TestPublishCharges_50Concurrent` (billing)
- [ ] `TestAllocatePayment_50Concurrent` (billing)
- [ ] `TestWebhookDelivery_100Concurrent_NoDuplicates` (webhooks)
- [ ] `TestExportJob_50Concurrent_OnlyOneSucceeds` (exports)

---

## 📝 ERROR CONTRACT STABILIZATION

**Current Error Keys (from PMS_MASTER.md):**

```
pms.error.invalid_status_transition
pms.error.already_checked_out
pms.error.occupancy_overlap
pms.error.unit_not_found
pms.error.access_denied
pms.portal.error.unit_access_denied
```

**Standardization Tasks:**
- [ ] Audit all `response.Error(c, ...)` calls in PMS handlers
- [ ] Create comprehensive error key inventory
- [ ] Ensure i18n keys exist for all error codes
- [ ] Update Next.js toast messages to use error keys
- [ ] Document error contract in TASK_015_SUMMARY.md

**Error Key Pattern:**
```
pms.<module>.error.<error_type>

Examples:
pms.metering.error.already_published
pms.billing.error.insufficient_balance
pms.reconciliation.error.already_confirmed
pms.work_orders.error.invalid_transition
pms.exports.error.job_already_completed
```

---

## 🚀 IMPLEMENTATION PLAN

### Phase 1: Audit & Documentation (Day 1)
- [x] Analyze PMS module structure
- [x] Create hardening checklist (this document)
- [ ] Audit idempotency patterns across modules
- [ ] Audit state machine implementations
- [ ] Document current vs target state

### Phase 2: P0-1 Idempotency (Day 2-3)
- [ ] Add idempotency tests for billing (publish, allocate)
- [ ] Add idempotency tests for reconciliation (confirm)
- [ ] Add idempotency tests for exports (finalize)
- [ ] Verify webhook delivery deduplication
- [ ] Run: `go test ./internal/modules/pms/... -count=1`

### Phase 3: P0-2 State Machine Standardization (Day 3-4)
- [ ] Define domain errors per module (service/errors.go)
- [ ] Add 409 mapping in handlers (metering, work_orders, reconciliation)
- [ ] Write state machine 409 tests (3 modules)
- [ ] Update i18n error keys
- [ ] Run: `go test ./internal/modules/pms/... -count=1`

### Phase 4: P0-3 Partial Failure Tests (Day 4-5)
- [ ] Add partial failure tests for metering publish
- [ ] Add partial failure tests for reconciliation confirm
- [ ] Add partial failure tests for export finalize
- [ ] Document compensation patterns
- [ ] Run: `go test ./internal/modules/pms/... -count=1`

### Phase 5: P1 Operational Readiness (Day 5-6)
- [ ] Add Prometheus metrics to hot paths (billing, reconciliation, exports)
- [ ] Implement atomic claim for webhooks
- [ ] Add concurrency stress tests (50-100 goroutines)
- [ ] Verify/add TTL indexes
- [ ] Run: `go test ./internal/modules/pms/... -count=1`

### Phase 6: Release Notes & Summary (Day 6)
- [ ] Create TASK_015_SUMMARY.md with risk analysis
- [ ] Document migration steps (index creation)
- [ ] Update PMS_PRODUCTION_HARDENING_REPORT.md with S15 additions
- [ ] Create PR with small, logical commits
- [ ] Final validation: `go test ./internal/modules/pms/... ./cmd/api -count=1`

---

## ✅ ACCEPTANCE CRITERIA

### Code Quality
- [ ] All new code follows existing PMS patterns
- [ ] No regressions: All existing tests pass
- [ ] All new tests pass: `go test ./internal/modules/pms/... -count=1 ✅`
- [ ] Lint clean: No new linter warnings

### Test Coverage
- [ ] Idempotency tests: 4+ modules covered
- [ ] State machine 409 tests: 3+ modules covered
- [ ] Partial failure tests: 3+ orchestration chains covered
- [ ] Concurrency stress tests: 4+ hot paths covered
- [ ] Total new tests: 20+ added

### Documentation
- [ ] TASK_015_SUMMARY.md complete with risk analysis
- [ ] All code changes documented with comments
- [ ] Error contract inventory complete
- [ ] Migration guide (index creation) documented

### Operational
- [ ] Prometheus metrics added to 5+ hot paths
- [ ] TTL indexes verified/added
- [ ] Request-ID propagation verified
- [ ] Worker concurrency hardened (webhooks, exports)

---

## 📊 RISK ASSESSMENT

| Risk | Severity | Likelihood | Mitigation |
|------|----------|------------|------------|
| Breaking existing functionality | 🔴 HIGH | 🟢 LOW | Comprehensive test suite + zero regression requirement |
| Performance degradation | 🟡 MEDIUM | 🟢 LOW | Load tests with 1000 units + 50-100 concurrent goroutines |
| Incomplete idempotency coverage | 🟡 MEDIUM | 🟡 MEDIUM | Systematic audit + proof tests per critical path |
| State machine inconsistency | 🟡 MEDIUM | 🟡 MEDIUM | Standardize on Access module pattern + tests |
| Worker race conditions | 🔴 HIGH | 🟡 MEDIUM | Atomic claim pattern + concurrency stress tests |

---

## 📚 REFERENCES

- **Foundation:** [PMS_PRODUCTION_HARDENING_REPORT.md](/workspace/PMS_PRODUCTION_HARDENING_REPORT.md) - Metering hardening (S14)
- **Architecture:** [PMS_MASTER.md](/workspace/aystek-core-dev-docs/PMS_MASTER.md) - Module specifications
- **Previous Work:** [P0_002_COMPLETION_SUMMARY.md](/workspace/P0_002_COMPLETION_SUMMARY.md) - Payment auto-finalize
- **Previous Work:** [P0_003_FINAL_VALIDATION_REPORT.md](/workspace/P0_003_FINAL_VALIDATION_REPORT.md) - Metering mixed allocation

---

## 📞 STAKEHOLDERS

**Engineering:** Backend team (Go), Frontend team (Next.js), DevOps (metrics/alerts)  
**QA:** Test coverage validation, regression testing  
**Product:** Risk analysis review, go/no-go decision

---

**Status:** 🔄 IN PROGRESS  
**Last Updated:** February 21, 2026  
**Next Review:** Phase 2 completion (idempotency tests)
