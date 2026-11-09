# API Sharing Implementation - Quick Reference Card

**Prepared for:** Engineering Teams  
**Timeline:** 6-9 months (12 sprints)  
**Next Milestone:** Sprint 11 Staging Load Validation

> Status note: this is a planning-oriented quick reference, not a statement that all listed sprint scopes are already implemented or production-approved. Use it together with `0-ARCHITECTURE-REVIEW-AND-REMEDIATION-PLAN.md`, `4-IMPLEMENTATION-TIMELINE.md`, and the Sprint 11/12 evidence documents.

---

## 📍 Where to Start

### If you're starting Sprint 1 (Operations & Scopes)
1. **Read these first** (30 min)
   - [INDEX.md](INDEX.md) — Full roadmap
   - [4-IMPLEMENTATION-TIMELINE.md](4-IMPLEMENTATION-TIMELINE.md) — Sprint 1 section
   - [3-DB-SCHEMA-DESIGN.md](3-DB-SCHEMA-DESIGN.md) — Collections overview

2. **Review code** (60 min)
   - `/workspace/go/aystek-core/internal/modules/integrations/api_sharing/model/api_sharing.go`
   - `/workspace/go/aystek-core/internal/modules/integrations/api_sharing/repository/interfaces.go`
   - `/workspace/go/aystek-core/internal/modules/integrations/api_sharing/service/interfaces.go`

3. **Start coding**
   - Implement `repositories.ApiOperationRepository` (MongoDB)
   - Implement `service.OperationService` (business logic)
   - Extend the current admin handler surface in `internal/modules/integrations/api_sharing/handler/admin.go`
   - Wire or verify routes in `internal/modules/integrations/api_sharing/routes.go`

---

### If you're starting Sprint 4 (OAuth2/OIDC)
1. **Read these first** (45 min)
   - [1-RFC-COMPLIANCE-CHECKLIST.md](1-RFC-COMPLIANCE-CHECKLIST.md) — OAuth2 section (RFC 6749, 7519, 7662, 7009, 9700)
   - [4-IMPLEMENTATION-TIMELINE.md](4-IMPLEMENTATION-TIMELINE.md) — Sprint 4 section
   - RFC 6749 (OAuth2 Framework) — https://tools.ietf.org/html/rfc6749

2. **Key endpoints** (need implementation)
   ```
   POST   /oauth/token                    # client_credentials + auth code
   GET    /oauth/authorize                # auth code with PKCE
   GET    /oauth/jwks                     # public key set (RFC 7517)
   POST   /oauth/introspect               # token info (RFC 7662)
   POST   /oauth/revoke                   # token revocation (RFC 7009)
   ```

3. **Required packages**
   ```bash
   go get github.com/golang-jwt/jwt/v5
   go get golang.org/x/crypto/pkce        # for PKCE (RFC 7636)
   ```

---

### If you're starting Sprint 6 (Audit & KVKK)
1. **Read these first** (60 min)
   - [2-KVKK-COMPLIANCE-MAPPING.md](2-KVKK-COMPLIANCE-MAPPING.md) — Turkish DPL requirements
   - [4-IMPLEMENTATION-TIMELINE.md](4-IMPLEMENTATION-TIMELINE.md) — Sprint 6 section
   - Article focus: 11 (SAR), 17 (Deletion)

2. **Key responsibilities**
   - Audit logging with PII masking (AuditService)
   - Subject Access Request endpoint (SAR)
   - Deletion/anonymization job (Kubernetes CronJob)
   - TTL enforcement (MongoDB)

---

### If you're starting Sprint 9 (Testing & Security)
1. **Read these first** (90 min)
   - [1-RFC-COMPLIANCE-CHECKLIST.md](1-RFC-COMPLIANCE-CHECKLIST.md) — Full spec
   - [4-IMPLEMENTATION-TIMELINE.md](4-IMPLEMENTATION-TIMELINE.md) — Sprint 9 section
   - OWASP API Top 10 — https://owasp.org/www-project-api-security/

2. **Test matrix** (required)
   - Contract tests: OAuth2 flows (client_credentials, auth code + PKCE)
   - Integration tests: Grant → Token → API call → audit log
   - Security tests: BOLA, BFLA, data exposure, rate limit bypass
   - KVKK tests: Consent withdrawal → audit deletion, SAR export, TTL enforcement
   - Chaos tests: Gateway/auth failure, network partition, Redis outage

3. **Tools needed**
   ```bash
   go get github.com/pact-foundation/pact-go              # Contract testing
   go get gotest.tools                                    # Test utilities
   go get github.com/stretchr/testify                     # Assertions
   ```

---

### If you're starting Sprint 11 (Load & Performance)
1. **Read these first** (45 min)
   - [SPRINT_11_LOAD_TEST_PERFORMANCE_REPORT.md](SPRINT_11_LOAD_TEST_PERFORMANCE_REPORT.md)
   - [SPRINT_11_CAPACITY_PLANNING.md](SPRINT_11_CAPACITY_PLANNING.md)
   - [4-IMPLEMENTATION-TIMELINE.md](4-IMPLEMENTATION-TIMELINE.md) — Sprint 11 section

2. **Run benchmarks**
   ```bash
   go test ./internal/modules/integrations/api_sharing/service \
     -run '^$' \
     -bench 'BenchmarkGetOperation_(With|Without)Cache(_SlowRepository)?$|BenchmarkValidateScopes_WithCache$' \
     -benchmem
   ```

3. **Run staging load profile with Prometheus export**
   ```bash
      cd /workspace/go/aystek-core

    export BASE_URL="https://<staging-api-host>"
    export AUTH_TOKEN="<staging-token>"
    export K6_PROMETHEUS_RW_SERVER_URL="https://<prometheus-host>/api/v1/write"

      internal/modules/integrations/api_sharing/tests/load/run_first_staging_once.sh
   ```

---

### If you're starting Sprint 12 (Go-Live Readiness)
1. **Read these first** (30 min)
   - [SPRINT_12_PRODUCTION_READINESS_CHECKLIST.md](SPRINT_12_PRODUCTION_READINESS_CHECKLIST.md)
   - [SPRINT_12_DEPLOYMENT_RUNBOOK.md](SPRINT_12_DEPLOYMENT_RUNBOOK.md)

2. **Validate readiness gates**
   - RFC + KVKK verification evidence linked
   - Staging load test evidence attached
   - Sign-off matrix completed (Backend, SRE, Security, Compliance, Product)

3. **Perform rollout rehearsal**
   - Blue/green test cutover
   - Rollback test
   - Incident drill based on runbook

---

## 🏗️ Code Structure

```
/workspace/go/aystek-core/internal/modules/integrations/api_sharing/
│
├── model/                          # Domain models
│   └── api_sharing.go
│
├── dto/                            # Request/response DTOs
│   └── api_sharing.go
│
├── repository/                     # Data layer
│   ├── interfaces.go               # Repository contracts
│   ├── mongo.go                    # Mongo implementation (grants/consents/audit/oauth/etc.)
│   ├── operation_mongo.go          # Operation repository implementation
│   └── scope_mongo.go              # Scope repository implementation
│
├── service/                        # Business layer
│   ├── interfaces.go
│   ├── service.go                  # Main service implementation
│   ├── jwt_service.go              # JWT/JWKS utilities
│   └── performance_benchmark_test.go
│
├── handler/                        # HTTP handlers
│   ├── admin.go
│   ├── internal.go
│   ├── oauth.go
│   ├── operation_handler.go
│   ├── registry.go
│   └── scope_handler.go
│
├── middleware/                     # Policy enforcement + audit middleware
│   └── policy_enforcement.go
│
├── routes.go                       # Route registration
└── tests/load/
   ├── policy_check_k6.js          # Sprint 11 load profile
   ├── run_staging_k6.sh           # k6 + Prometheus remote write runner
   ├── summarize_k6_results.sh     # summary json -> markdown
   ├── run_first_staging_once.sh   # one-shot staging runner
   └── export_grafana_artifacts.sh # dashboard/panel export helper
```

---

## 🔄 Dependency Flow

```
HTTP Layer (handler)
    ↓
Service Layer (business logic)
    ↓
Repository Layer (data access)
    ↓
MongoDB (data storage)
    +→ Redis (for rate limiting cache)
    +→ JWKS endpoint (for token validation)
```

**Example:** Sprint 1 HTTP endpoint flow
```
POST /admin/api-sharing/operations                    # operation_handler.CreateOperation()
  → OperationService.RegisterOperation()             # business logic
    → ApiOperationRepository.Create()                # save to MongoDB
      → db.api_operations.insertOne()                # MongoDB insert
```

---

## 📋 Checklist for First Implementation

### Setup (Day 1)
- [ ] Clone/checkout code at `/workspace/go/aystek-core/`
- [ ] Review models at `internal/modules/integrations/api_sharing/model/`
- [ ] Confirm existing repository/service layout before creating any new directories
- [ ] Read [INDEX.md](INDEX.md) (15 min)

### First Sprint Task (Days 1-5)
- [ ] Compare the current repo state with `4-IMPLEMENTATION-TIMELINE.md` before starting any Sprint 1 task
- [ ] Reuse existing `api_sharing` repository/service/handler structure where present
- [ ] Fill only the gaps that remain unimplemented or unverified
- [ ] Add or refresh tests for the specific gap you touched

### Code Review Checklist
- [ ] All models use BSON/JSON tags (for MongoDB serialization)
- [ ] Error handling follows existing patterns (e.g., `auth/handler/`)
- [ ] Tests cover happy path + error cases
- [ ] Comments document complex logic
- [ ] No PII logged without masking
- [ ] Soft-delete pattern for compliance records

---

## 🔐 Security Checklist (Per Sprint)

### For all sprints:
- [ ] No hardcoded secrets (use environment variables)
- [ ] Input validation on all request parameters
- [ ] SQL injection prevention (not applicable to MongoDB, but prevent NoSQL injection)
- [ ] CORS configured properly
- [ ] PII masking in logs + audit trail

### For auth-related sprints (3-4):
- [ ] Token signing validated (RS256)
- [ ] Scope validation enforced
- [ ] Rate limits + quota honored
- [ ] IP allowlist checked
- [ ] Revoked tokens rejected

### For KVKK/audit sprints (6):
- [ ] Consent withdrawal → audit logs deleted
- [ ] SAR export includes all user data
- [ ] TTL deletion job runs + verifiable
- [ ] PII anonymization applied

---

## 🚀 Quick Links

**Master Documents:**
- [INDEX.md](INDEX.md) — Full roadmap (start here)
- [4-IMPLEMENTATION-TIMELINE.md](4-IMPLEMENTATION-TIMELINE.md) — Sprint-by-sprint plan

**Compliance & Design:**
- [1-RFC-COMPLIANCE-CHECKLIST.md](1-RFC-COMPLIANCE-CHECKLIST.md) — OAuth2/JWT/Rate limiting requirements
- [2-KVKK-COMPLIANCE-MAPPING.md](2-KVKK-COMPLIANCE-MAPPING.md) — Turkish DPL mapping
- [3-DB-SCHEMA-DESIGN.md](3-DB-SCHEMA-DESIGN.md) — Database design

**Code References:**
- Models: `/workspace/go/aystek-core/internal/modules/integrations/api_sharing/model/api_sharing.go`
- Repository interfaces: `/workspace/go/aystek-core/internal/modules/integrations/api_sharing/repository/interfaces.go`
- Service interfaces: `/workspace/go/aystek-core/internal/modules/integrations/api_sharing/service/interfaces.go`
- DTOs: `/workspace/go/aystek-core/internal/modules/integrations/api_sharing/dto/api_sharing.go`

**Existing Patterns to Follow:**
- Repository pattern: `/workspace/go/aystek-core/internal/modules/auth/repository/`
- Service pattern: `/workspace/go/aystek-core/internal/modules/auth/service/`
- Handler pattern: `/workspace/go/aystek-core/internal/modules/auth/handler/`

**External References:**
- OAuth2 (RFC 6749): https://tools.ietf.org/html/rfc6749
- JWT (RFC 7519): https://tools.ietf.org/html/rfc7519
- Rate Limiting (RFC 6585): https://tools.ietf.org/html/rfc6585
- OAuth2 BCP (RFC 9700): https://tools.ietf.org/html/rfc9700
- KVKK 6698: https://www.mevzuat.gov.tr/MevzuatMetin/1.5.6698.pdf

---

## 📞 Getting Help

**For Q&A during implementation:**

1. **Code structure questions:** Check [INDEX.md](INDEX.md) section "Code Organization"
2. **API endpoint contract:** Check [4-IMPLEMENTATION-TIMELINE.md](4-IMPLEMENTATION-TIMELINE.md) sprint details
3. **RFC compliance requirements:** Check [1-RFC-COMPLIANCE-CHECKLIST.md](1-RFC-COMPLIANCE-CHECKLIST.md)
4. **KVKK legal requirements:** Check [2-KVKK-COMPLIANCE-MAPPING.md](2-KVKK-COMPLIANCE-MAPPING.md)
5. **Database schema:** Check [3-DB-SCHEMA-DESIGN.md](3-DB-SCHEMA-DESIGN.md)
6. **Existing patterns:** Review `/workspace/go/aystek-core/internal/modules/auth/` (similar structure)

**Escalation path:**
- Technical lead → Architecture team → Project manager

---

**Version:** 1.1  
**Last Updated:** 2026-03-11  
**Next Update:** Sprint 11 load-test evidence review
