# P1 Remaining Work Tasks (5% Outstanding)

**Priority Level:** All tasks are NON-BLOCKING for go-live. Recommend Phase 2 (post-launch).

---

## 🟡 Optional Tasks (2-3 Days Work)

### **P1-E9: Admin CSV Bulk Operations Endpoints** (2-3 hours)
**Epic:** Inventory & Warehouse Operations  
**Current State:** CSV import/export service exists (`/internal/modules/stock/service/csv_service.go`)  
**Outstanding:** Admin HTTP endpoints to trigger bulk operations

#### **Tasks:**
```
[ ] Create /admin/stock/bulk-import endpoint
    - Accept CSV file upload
    - Validate schema (SKU, quantity, warehouse)
    - Trigger AsyncImportJob in background
    - Return job ID + status polling URL

[ ] Create /admin/stock/bulk-export endpoint
    - Generate current stock snapshot CSV
    - Download as file
    - Log export in activity_logs

[ ] Create /admin/jobs/{jobId}/status endpoint
    - Poll import job progress
    - Show processed/failed rows
    - List errors (if any)

[ ] Create /admin/stock/import-history endpoint
    - List past import jobs
    - Show success rate + row counts
    - Link to error reports
```

**Implementation Guide:**
- File: Create `/internal/modules/stock/handler/bulk_operations.go`
- Dependencies: csv_service.go (already exists), jobs module (queue management)
- Routes: Register in `/internal/modules/stock/routes.go`
- Tests: Add 3-4 integration tests for upload + validation + job creation

**Effort:** 2-3 hours (straightforward CRUD + file handling)

---

### **P1-X1: OpenAPI Auto-Generation** (3-4 hours)
**Cross-Cutting Epic:** API Standardization  
**Current State:** Manual OpenAPI spec at `/docs/openapi.yaml`  
**Outstanding:** Auto-generation from Go code annotations

#### **Tasks:**
```
[ ] Add go:generate comment annotations to handlers
    - e.g., // @Router /admin/products [GET]
    - Include request/response models
    - Document error responses

[ ] Install swag CLI tool
    $ go install github.com/swaggo/swag/cmd/swag@latest

[ ] Configure swag generation
    - Add Makefile target: make swagger-gen
    - Output to docs/swagger.json + docs/openapi.yaml
    - CI/CD integration (auto-generate on each commit)

[ ] Validate generated spec
    - Test against swagger-ui
    - Compare with manual spec (find gaps)
    - Document any manual additions needed
```

**Implementation Guide:**
- Tool: Use [swag](https://github.com/swaggo/swag) (Go doc annotations → OpenAPI)
- Alternative: Use [ogen](https://github.com/ogen-go/ogen) (code generation from OpenAPI)
- Routes file: Update each handler with swagger tags

**Effort:** 3-4 hours (tooling + integration + validation)

---

### **P1-X2: PII Masking Rules** (1-2 hours)
**Cross-Cutting Epic:** RBAC + Audit + Compliance  
**Current State:** Activity logs exist, no PII masking  
**Outstanding:** Mask sensitive fields (email, phone, payment details)

#### **Tasks:**
```
[ ] Create /internal/middleware/pii_masker.go
    - Middleware to mask PII in logs before storage
    - Mask patterns:
      * Email: user@example.com → u****@example.com
      * Phone: +90555123456 → +90555****56
      * Card: 4532123456789012 → ****9012
      * SSN: 12345678901 → ****01
    - Applied at activity_logs write time

[ ] Add PII config rules
    - Define maskable fields per module
    - Configurable mask patterns (regex)
    - Test coverage: unit tests for each pattern

[ ] Update activity_logs service
    - Before SaveActivityLog() call masking middleware
    - Ensure masked data is final (no unmasking)
    - Log audit trail of what was masked
```

**Implementation Guide:**
- Pattern: Middleware chain pattern (similar to correlation ID injection)
- Storage: Don't store original PII in logs (PII goes to dedicated encrypted store if needed)
- Audit: Track what was masked (for compliance)

**Effort:** 1-2 hours (straightforward middleware)

---

### **P1-X3: Dead Letter Queue (DLQ) + Event Reprocess UI** (4-6 hours)
**Cross-Cutting Epic:** Event-Driven Reliability  
**Current State:** Event infrastructure exists, no DLQ for failed events  
**Outstanding:** Capture + retry mechanism for failed event processing

#### **Tasks:**
```
[ ] Create DLQ Model
    - FailedEvent struct:
      * EventID, EventType, Payload, ErrorMessage
      * FailedAt, LastRetryAt, RetryCount
      * Status (pending_retry, dead, archived)

[ ] DLQ Service (save + query)
    - SaveFailedEvent(event, error) → DLQ store
    - ListFailedEvents(filter) → paginated results
    - RetryFailedEvent(eventID) → re-process

[ ] Event Handler Wrapper
    - Wrap all event handlers with try-catch
    - On error: save to DLQ + log
    - Max 3 retries before marking as dead
    - Exponential backoff (1s, 5s, 30s)

[ ] Reprocess UI Endpoints
    - GET /admin/events/dlq → list failed events
    - POST /admin/events/dlq/{eventID}/retry → manual retry
    - DELETE /admin/events/dlq/{eventID} → archive event
    - POST /admin/events/dlq/bulk-retry → retry failed batch

[ ] Background Retry Job
    - Cron job: Every 5 minutes
    - Auto-retry pending_retry events (max 3 times)
    - Alert on events stuck > 1 hour
```

**Implementation Guide:**
- Storage: MongoDB collection `failed_events`
- Pattern: Standard error recovery (retry + backoff + DLQ)
- Admin UI: CRUD endpoints (will have UI component built in Phase 2)

**Effort:** 4-6 hours (service + handler + background job)

---

### **P1-X4: Localization Coverage Audit** (1-2 hours)
**Cross-Cutting Epic:** Seed & Localization Governance  
**Current State:** 28 languages supported, some modules may have incomplete translations  
**Outstanding:** Identify + fill translation gaps

#### **Tasks:**
```
[ ] Run Translation Coverage Audit
    - Script: List all translation keys in code
    - Compare with translation_seed.go
    - Report: Which modules have <90% coverage
    - Output: Coverage report (CSV)

[ ] Fill Critical Translation Gaps
    - Identify missing keys in EN, TR, DE
    - Add translations for: payment, order, shipping, returns
    - Test: Ensure UI renders correctly in each language

[ ] Update Translation Seed
    - Add missing keys to translation_seed.go
    - Ensure pluralization rules handled
    - Document language fallback (e.g., PT-BR → PT if not found)

[ ] Add Translation Validator
    - Middleware/startup check
    - Warn if key missing in language
    - Test in CI/CD
```

**Implementation Guide:**
- Tool: [gotext](https://github.com/go-echarts/gotext) or simple grep-based script
- Coverage Target: ≥95% for top 5 languages (EN, TR, DE, FR, IT)
- Other languages: ≥80% acceptable for Phase 1

**Effort:** 1-2 hours (audit + gap filling)

---

## 📅 Recommended Schedule

### **Go-Live (Late February 2026)**
- Ship with P1-E6 through P1-E11 ✅
- Ship with P1-X1 through P1-X4 (core functionality) ✅
- Skip: Optional admin UI, auto-generation, advanced DLQ

### **Phase 2 (Week 1 Post-Launch)**
```
Day 1-2: P1-E9 admin CSV bulk endpoints (2-3 hours)
Day 2-3: P1-X1 OpenAPI auto-generation (3-4 hours)
Day 3-4: P1-X2 PII masking (1-2 hours)
```

### **Phase 2 (Week 2 Post-Launch)**
```
Day 1-3: P1-X3 DLQ + reprocess UI (4-6 hours)
Day 3-4: P1-X4 translation coverage audit (1-2 hours)
Day 4-5: P2-E13 marketplace MVP (start)
```

---

## 🎯 Success Criteria

### **For Go-Live (P1 at 98%)**
- ✅ P0 critical paths passing
- ✅ P1 epics E6-E11 feature-complete
- ✅ RBAC + audit infrastructure working
- ✅ Seed data loaded (permissions, localization)
- ✅ Event architecture operational
- ✅ OpenAPI documentation available (manual is fine)

### **For Phase 2 (P1 at 100%)**
- ✅ All optional admin endpoints implemented
- ✅ OpenAPI auto-generation integrated
- ✅ PII masking enforced
- ✅ DLQ + reprocess UI operational
- ✅ Translation coverage ≥95% (top 5 languages)

---

## 📝 Tracking

**Update:** 2026-02-02  
**Owner:** Engineering Team  
**Status:** Ready for backlog prioritization

---

### Quick Reference: File Paths

| Task | Key Files | Priority |
|------|-----------|----------|
| CSV Admin Endpoints | `/internal/modules/stock/handler/bulk_operations.go` | 🟡 Low |
| OpenAPI Auto-Gen | `/docs/openapi.yaml`, `Makefile` | 🟡 Low |
| PII Masking | `/internal/middleware/pii_masker.go` | 🟡 Medium |
| DLQ + Reprocess | `/internal/modules/events/dlq/`, `/internal/platform/jobs/retry_job.go` | 🟡 Medium |
| Localization Audit | `/internal/seeds/translation_seed.go` | 🟡 Low |

