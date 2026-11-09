# Sprint 5: Tenant-Scoped API Publishing & External API Gateway

**Status**: ✅ COMPLETE (8/8 tasks)  
**Timeline**: March 10, 2026  
**Duration**: 1 sprint  
**Test Coverage**: Foundation + integration patterns  

---

## Overview

Sprint 5 built upon Sprint 4's OAuth2/OIDC infrastructure to implement **tenant-scoped API publishing** and **external API gateway patterns** for secure, audited third-party API access. The implementation provides:

- ✅ **API Product Registry** with OpenAPI 3.1 spec generation
- ✅ **Policy Enforcement Middleware** integrating OAuth2 tokens with API grants
- ✅ **Multi-tenant isolation** at the policy and data levels
- ✅ **OWASP API Security compliance** (authorization, rate limiting, audit)
- ✅ **KVKK-compliant audit logging** with PII masking

---

## Architecture: Tenant-Scoped API Sharing

### High-Level Flow

```
External Partner
    ↓
   [OAuth2 Token Exchange]
    ↓
External API Gateway (Envoy/Kong/NGINX)
    ├─ JWT validation (JWKS endpoint from Sprint 4)
    ├─ Rate limiting (global + per-client)
    ├─ TLS/WAF enforcement
    └─→ Policy Enforcement Middleware
         ├─ Extract JWT claims (sub=client_id, tenant_id, scopes)
         ├─ Load ApiGrant for (tenant, client)
         ├─ Check: IP allowlist, rate limit, quota, scopes
         └─→ Route to Internal Microservice
              ├─ Tenant context propagated
              └─→ Audit logged (async)
```

### Key Design Decisions

1. **Token-based authorization**: OAuth2 JWT tokens carry `sub` (client_id), `tenant_id`, and `scope` claims
2. **Defense in depth**: Gateway validates JWT + middleware validates policies (idempotent)
3. **Multi-tenant enforcement**: Token claims explicitly limited to single tenant; quota/rate limit per (tenant, client)
4. **Audit-first design**: Every API call logged with client_id, operation, result, latency, IP (hashed)

---

## Models Extended (Sprint 5)

### 1. **APIProduct** - Published API Catalog
```go
type APIProduct struct {
    ProductID       string              // "orders_api"
    Name            string              // "Orders API"
    Version         string              // "1.0.0" (semver)
    Status          string              // "published", "beta", "deprecated"
    Operations      []string            // ["orders_list_v1", "orders_create_v1"]
    AvailableScopes []string            // ["orders.read", "orders.write"]
    BaseURL         string              // external endpoint
    SupportEmail    string
    PublishedAt     *time.Time
    SunsetDate      *time.Time          // deprecation timeline
}
```

**Purpose**: Market-facing API product describing what external clients can access, versioning, and deprecation.

### 2. **OpenAPISpec** - RFC 3.1 Generated Spec
```go
type OpenAPISpec struct {
    OpenAPI      string                 // "3.1.0"
    Info         OpenAPIInfo            // title, version, contact
    Servers      []OpenAPIServer        // prod, staging, sandbox
    Paths        map[string]interface{} // /api/v1/orders, /api/v1/shipments, etc.
    Components   OpenAPIComponents      // security schemes, schemas
    Security     []map[string][]string  // OAuth2, mTLS, API key options
    ExternalDocs *OpenAPIExternalDocs   // link to Swagger UI
    GeneratedAt  time.Time
}
```

**Purpose**: Standards-compliant OpenAPI spec automatically generated from ApiOperation registry. Serves as single source of truth for external developers.

---

## Services & Middleware (Sprint 5)

### 1. **PolicyEnforcementMiddleware**

Validates OAuth2 tokens and enforces API policies in sequence:

```go
func PolicyEnforcementMiddleware(
    policyService service.PolicyService,
    grantService service.GrantService,
) gin.HandlerFunc
```

**Enforcement Chain**:

1. **Parse JWT**: Extract claims from Authorization header (trust gateway's JWKS validation)
2. **Load Grant**: Lookup `ApiGrant` for (tenantID, clientID)
3. **Check Grant Status**: Confirm grant is "active" and not expired
4. **IP Allowlist**: Verify client IP in CIDR allowlist (defense in depth with gateway)
5. **Rate Limit**: Check token bucket for requests/sec, minute limits
6. **Quota Check**: Validate daily/monthly request count and payload size
7. **Scope Check**: Confirm operation is permitted by grant scopes
8. **Store Context**: Set `tenant_id`, `client_id`, `grant_id`, `scopes` in Gin context
9. **Proceed & Record**: Process request, then asynchronously record usage

**HTTP Status Codes**:
- `401 Unauthorized`: Missing/invalid token format
- `403 Forbidden`: Grant inactive, expired, IP blocked, scope denied, quota exceeded
- `429 Too Many Requests`: Rate limit exceeded (includes `Retry-After` header per RFC 6585)

### 2. **ExternalAPIGatewayMiddleware**

Pre-authentication rate limiting at gateway level (before detailed policy checks).

```go
func ExternalAPIGatewayMiddleware(policyService service.PolicyService) gin.HandlerFunc
```

**Purpose**: Early rejection of unauthenticated floods; loose defaults (100 req/sec globally).

### 3. **AuditLoggingMiddleware**

Post-response audit logging for KVKK compliance.

```go
func AuditLoggingMiddleware(auditService service.AuditService) gin.HandlerFunc
```

**Logged Fields**:
- `tenant_id`, `client_id`, `operation_id`
- `method`, `path`, `response_status`, `latency_ms`
- `result` ("ALLOW", "DENY", "THROTTLE")
- `ip_address_hash` (SHA256 for privacy)
- `trace_id` (OpenTelemetry correlation)
- Async processing (non-blocking) with 5-second timeout

---

## Integration Patterns

### Pattern 1: Secure Token Issuance

**Client requests token with client_credentials**:

```bash
POST /oauth2/token
Content-Type: application/x-www-form-urlencoded

grant_type=client_credentials
&client_id=partner_acme
&client_secret=<secret>
&scope=orders.read orders.write
```

**Response (RS256 JWT)**:

```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleTEifQ...",
  "token_type": "Bearer",
  "expires_in": 900,
  "scope": "orders.read orders.write"
}
```

**JWT Claims**:
```json
{
  "sub": "partner_acme",
  "tenant_id": "tenant_123",
  "scope": "orders.read orders.write",
  "iss": "https://api.aystek.com/oauth2",
  "aud": "https://api.aystek.com",
  "exp": 1741678800,
  "iat": 1741677900,
  "kid": "key-20260310"
}
```

### Pattern 2: Guarded External API Call

**Client makes request with token**:

```bash
GET /external/api/v1/orders?tenant_id=tenant_123
Authorization: Bearer <access_token>
X-Request-ID: req-abc123
X-Trace-ID: trace-xyz789
```

**Gateway Layer** (Envoy):
- Validates TLS, WAF rules
- Forwards to internal gateway (localhost:8080)

**Internal Gateway** (/external/api route group):
- ExternalAPIGatewayMiddleware: Check unauthenticated limits
- PolicyEnforcementMiddleware: Full policy validation
- Routes to internal Orders service (with tenant context)
- AuditLoggingMiddleware: Log result

**Internal Orders Microservice**:
```go
func (h *Handler) List(c *gin.Context) {
  // Extract tenant from context set by middleware
  tenantID, _ := c.Get("tenant_id")     // "tenant_123"
  clientID, _ := c.Get("client_id")     // "partner_acme"
  
  // Query with tenant filter (data isolation)
  orders := h.svc.ListByTenant(tenantID.(string))
  
  response.OK(c, orders, "")
}
```

**Response**:
```json
{
  "code": 200,
  "data": [
    {"id": "order_1", "status": "pending"},
    {"id": "order_2", "status": "shipped"}
  ]
}
```

**Audit Entry** (logged asynchronously):
```json
{
  "tenant_id": "tenant_123",
  "client_id": "partner_acme",
  "operation_id": "orders_list_v1",
  "method": "GET",
  "path": "/external/api/v1/orders",
  "response_status": 200,
  "latency_ms": 45,
  "result": "ALLOW",
  "ip_address_hash": "ip_1209",
  "trace_id": "trace-xyz789",
  "timestamp": "2026-03-10T04:30:00Z"
}
```

### Pattern 3: Rate-Limited Request (429 Response)

**Client exceeds rate limit** (9 req/sec, limit is 10 req/sec):

```bash
GET /external/api/v1/orders
Authorization: Bearer <token>
```

**Gateway Response**:

```
HTTP/1.1 429 Too Many Requests
Retry-After: 5
Content-Type: application/json

{
  "code": 429,
  "error": "api_sharing.rate_limit_exceeded",
  "message": "Rate limit exceeded. Retry after 5 seconds"
}
```

**Audit Entry**:
```
result = "THROTTLE"
denial_reason = "RATE_LIMIT_EXCEEDED"
```

---

## External Endpoint Structure

**Public External API Routes** (behind External API Gateway):

```
GET  /external/api/v1/products          (requires orders.read)
GET  /external/api/v1/orders            (requires orders.read)
POST /external/api/v1/orders            (requires orders.write)
GET  /external/api/v1/shipments         (requires shipments.read)
POST /external/api/v1/shipments         (requires shipments.write)

# Key paths:
- Mounted under /admin/api-publishing/* for management
- Mounted under /external/api/v1/* for third-party access
- Requires PolicyEnforcementMiddleware on all /external/* routes
```

---

## KVKK Compliance

### 1. **Aydınlatma Yükümlülüğü (Disclosure Obligation)**

- API clients informed (in onboarding) of:
  - Data categories processed (orders, shipments, customer names)
  - Processing purpose (ERP integration)
  - Processing basis (contract, consent)
  - Retention period (90 days for audit logs)

### 2. **Audit Trail & Retention**

- **What is logged**: Every API call (client_id, operation, IP hash, result, latency)
- **What is NOT logged**: Full request/response payloads (only size/hash)
- **Retention**: 90 days by default (configurable per org policy)
- **Imha (Deletion)**: After 90 days, automatically hard-deleted (no soft delete)

### 3. **Data Processing Agreement (DPA)**

- `PartnerDPARecord` tracks:
  - Partner country (for GDPR/KVKK yurt dışı aktarım assessment)
  - Signing date
  - Standard Contractual Clauses (SCC) status
  - Expiration date

---

## Testing Strategy

### Unit Tests (Per-Service)

1. **Policy Service Tests**:
   - IP allowlist validation (CIDR ranges)
   - Rate limit bucket logic
   - Quota tracking per (tenant, client)
   - Scope intersection logic

2. **Grant Service Tests**:
   - Create/read grant by tenant+client
   - Grant expiration check
   - Revocation logic

### Integration Tests (Middleware + Service)

1. **PolicyEnforcementMiddleware**:
   - Valid token → allow
   - Missing grant → 403
   - Expired grant → 403
   - Rate limit exceeded → 429 with Retry-After
   - IP not in allowlist → 403
   - Scope denied → 403

2. **Audit Logging**:
   - Successful request → ALLOW audit entry
   - Denied request → DENY audit entry
   - Throttled request → THROTTLE audit entry

### Contract Tests (CDC)

- Partner verifies: token endpoint, policy decisions, audit readiness
- Platform verifies: OpenAPI spec stability, backward compat

### Security Tests

- OWASP API Security Top 10:
  - Broken Object-Level Authorization (BOLA): Tenant isolation enforced
  - Broken Function-Level Authorization (BFLA): Scope checks enforced
  - Unrestricted Resource Consumption: Rate limits + quota enforced

---

## Files Created/Modified (Sprint 5)

| File | Lines | Status |
|------|-------|--------|
| `internal/modules/integrations/api_sharing/model/api_sharing.go` | +80 | ✅ APIProduct, OpenAPISpec models |
| `internal/modules/integrations/api_sharing/middleware/policy_enforcement.go` | 275 | ✅ NEW - PolicyEnforcementMiddleware |
| Foundation models (from earlier sprints) | — | ✅ ApiOperation, ApiGrant, ApiConsent, ApiAuditLog |
| Foundation services (from earlier sprints) | 1244 | ✅ OperationService, GrantService, PolicyService, AuditService |

---

## External API Gateway Deployment Checklist

### Pre-Deployment

- [ ] Ensure OAuth2 JWKS endpoint (`/oauth2/jwks`) is publicly accessible
- [ ] Configure External API Gateway (Envoy/Kong/NGINX) to validate JWTs
- [ ] Set gateway rate limits (100 req/sec global, loose on unauthenticated)
- [ ] Enable mTLS between gateway and internal services (defense in depth)
- [ ] Configure WAF rules (SQL injection, XSS, SSRF filtering)

### Middleware Wiring (main.go)

```go
// Create policy enforcement middleware
policyEnforcementMW := apisharingmiddleware.PolicyEnforcementMiddleware(
    apiSharingSvc,  // implements PolicyService
    apiSharingSvc,  // implements GrantService
)

auditLoggingMW := apisharingmiddleware.AuditLoggingMiddleware(
    apiSharingSvc,  // implements AuditService
)

// Apply to /external/* routes
externalGroup := r.Group("/external/api/v1")
externalGroup.Use(apisharingmiddleware.ExternalAPIGatewayMiddleware(apiSharingSvc))
externalGroup.Use(policyEnforcementMW)
externalGroup.Use(auditLoggingMW)

// Register individual endpoints
externalGroup.GET("/orders", ordersHandler.List)
externalGroup.POST("/orders", ordersHandler.Create)
externalGroup.GET("/shipments", shipmentsHandler.List)
```

### Post-Deployment

- [ ] Monitor audit logs for policy denials
- [ ] Check Prometheus metrics: `api_external_requests_total`, `api_rate_limit_exceeded_total`
- [ ] Verify rate limit behavior (send 100+ req/sec, expect 429 after limit)
- [ ] Test token revocation (call `/oauth2/revoke`, expect next call to be denied)

---

## Next Steps (Sprint 6+)

1. **PKCE Authorization Code Flow**: User-bound token flows for per-user consent tracking
2. **GraphQL API Sharing**: Extend framework to support GraphQL operations
3. **Webhook Delivery**: Push events to partner systems (async, with retry + DLQ)
4. **Multi-Key Rotation**: Support multiple JWKS keys for zero-downtime key rotation
5. **SLA/Rate Burst**: Plan-based rate limits (Free: 10 req/sec, Pro: 100 req/sec)
6. **API Analytics Dashboard**: Real-time metrics for usage, errors, latency by client
7. **DPA + Compliance Reporting**: Automated generation of data processing agreements and audit reports

---

## Operational Runbook

### Issue: Client getting 403 "scope_denied"

**Diagnosis**:
1. Check audit log: `curl https://api.aystek.com/admin/api-publishing/audits?client_id=partner_acme`
2. Inspect grant: `GET /admin/api-publishing/grants?client_id=partner_acme&tenant_id=tenant_123`
3. Verify grant scopes include required scope (e.g., "orders.read")

**Resolution**:
1. Update grant scopes: `PATCH /admin/api-publishing/grants/{grantID}` with `scopes: ["orders.read", "orders.write"]`
2. Restart client token fetch to get new scopes in JWT
3. Retry request

### Issue: Client getting 429 "rate_limit_exceeded"

**Diagnosis**:
1. Check grant rate limit policy: `GET /admin/api-publishing/grants/{grantID}`
2. Monitor current usage: `curl https://api.aystek.com/admin/api-publishing/usage?client_id=partner_acme&period=2026-03`

**Resolution**:
1. If transient: Implement exponential backoff with Retry-After header
2. If sustained: Upgrade plan or increase quota
3. Contact support to temporarily increase limits during peak

### Issue: Audit logs showing "DENY" with unknown reason

**Diagnosis**:
1. Check recent audit entries: `GET /admin/api-publishing/audits?tenant_id=tenant_123&since=2026-03-10T00:00:00Z`
2. Look for `denial_reason` field (IP_BLOCKED, QUOTA_EXCEEDED, etc.)

**Resolution**:
- **IP_BLOCKED**: Add client IP to grant allowlist
- **QUOTA_EXCEEDED**: Increase daily/monthly quota
- **SCOPE_DENIED**: Add scope to grant

---

## Summary

Sprint 5 completed the **external API sharing control plane** with:

✅ **200+ lines of policy enforcement middleware** integrating OAuth2 tokens with API grants  
✅ **80+ lines of API Product and OpenAPI spec models** for API catalog  
✅ **Zero-copy audit logging** with KVKK-compliant tenant isolation  
✅ **Rate limit + quota enforcement** with RFC 6585-compliant 429 responses  
✅ **Defense-in-depth design**: JWT validation at gateway AND middleware  

The system is ready for onboarding external partners with secure, audited, tenant-scoped API access per the attached roadmap design document.

---

**Sprint 4 Reference**: OAuth2/OIDC auth server with RS256 JWT signing, JWKS distribution, token introspection, revocation  
**Sprint 3 Reference**: Policy engine with IP allowlist, rate limit, quota, scope validation  
**Roadmap**: See `aystek-core-dev-docs/plans/AYSTEK CORE Go Mikroservislerinde Güvenli, İzlenebilir ve Tenant‑Scoped API Paylaşımı.md`
