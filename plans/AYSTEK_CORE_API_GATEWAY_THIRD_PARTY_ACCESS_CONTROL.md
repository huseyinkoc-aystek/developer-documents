---
title: AYSTEK CORE — API Gateway & Third-Party Access Control
version: 1.0.0
status: Draft - Architecture Design
last_updated: 2026-03-09
classification: Internal - Enterprise Architecture
related_documents:
  - ../user-guide/04-api/API_REFERENCE.md
  - ../user-guide/02-security/AUTHENTICATION_AND_AUTHORIZATION.md
  - ../README_FOR_AI_ASSISTANTS.md
---

# AYSTEK CORE API Gateway & Third-Party Access Control

## Executive Summary

AYSTEK CORE, internal API endpoint'lerini (GET /products, POST /orders vb.) üçüncü taraf yazılımlara **OAuth2 Client Credentials** tabanlı, **tenant-site scoped**, **permission-limited**, **rate-limited** ve **audit-trail'li** şekilde açacaktır.

Model: AWS API Gateway, Kong, Stripe, GitHub'ın "third-party API consumption" yaklaşımına uygun.

**Temel Hedefler:**
- ✅ Third-party yazılımlar AYSTEK API'sını çağırabilir (REST/JSON)
- ✅ IP, user credential, scope ve rate-limit ile kontrol
- ✅ Tenant-site isolation (bir API user sadece kendi tenant'ına erişebilir)
- ✅ Enterprise-grade UI (API client create, scope select, key rotate, usage monitor)
- ✅ Audit trail (kim ne zaman hangi endpoint'i çağırdı)
- ✅ Security best practice'ler (OAuth2, secret rotation, TLS, idempotency)

**MVP Timeline:** 4–6 ay (backend + frontend + operasyon)

---

## Architecture Overview

### System Components

```
┌─────────────────────────────────────────────────────────────┐
│ Third-Party Applications (Partner, Vendor, Agency)         │
│ • Mobile apps                                              │
│ • Web dashboards                                           │
│ • ERP systems                                              │
│ • Marketing platforms                                      │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       │ OAuth2 Client Credentials Flow
                       │ (POST /auth/token + access_token)
                       │
        ┌──────────────▼───────────────┐
        │  API Gateway Auth Layer      │
        │  • Token validation          │
        │  • Scope enforcement         │
        │  • IP check (optional)       │
        │  • Rate limit / quota        │
        └──────────────┬───────────────┘
                       │
        ┌──────────────▼───────────────┐
        │  Authorization Layer         │
        │  • Route whitelist           │
        │  • Method check (GET/POST)   │
        │  • API User permissions      │
        └──────────────┬───────────────┘
                       │
        ┌──────────────▼───────────────┐
        │  Internal AYSTEK Routes      │
        │  • GET /products             │
        │  • POST /orders              │
        │  • GET /inventory            │
        │  • ... (permitted set only)  │
        └──────────────┬───────────────┘
                       │
        ┌──────────────▼───────────────┐
        │  Audit & Observability       │
        │  • Execution logs            │
        │  • Usage metrics             │
        │  • Security alerts           │
        └──────────────────────────────┘
```

### OAuth2 Client Credentials Flow (RFC 6749)

```
Third-Party App
      │
      │ 1. POST /auth/oauth2/token
      │    client_id=...
      │    client_secret=...
      │    grant_type=client_credentials
      │    scope=products:read,orders:read,orders:create
      │
      ▼
  API Auth Endpoint
      │
      │ 2. Validate client + secret
      │    Check scope against API user's allowed scopes
      │    Generate JWT access_token (15 min TTL) + refresh_token (7 days)
      │
      ▼ (success)
  Third-Party App receives:
      {
        "access_token": "eyJhbGc...",
        "token_type": "Bearer",
        "expires_in": 900,
        "refresh_token": "ref_..."
      }
      │
      │ 3. Call protected endpoint
      │    GET /products
      │    Authorization: Bearer eyJhbGc...
      │
      ▼
  API Gateway
      │
      │ 4. Extract JWT, verify signature + scope
      │    Check rate limit bucket for this API user
      │    Check IP whitelist (if configured)
      │
      ▼ (authorized)
  Internal Service
      │
      │ 5. Process request
      │    Filter by tenant_id (from JWT claims)
      │    Filter by site_id (if provided)
      │
      ▼
  Response (200 OK + data)
      │
      │ 6. Log execution (audit trail)
      │    • api_user_id
      │    • method + path
      │    • status code
      │    • latency
      │    • error (if any)
      │
      ▼
  Audit Store
```

---

## Data Model

### 1. API Client (Expansion of Users Table)

**Location:** `internal/modules/users/model/user.go` (add struct)

```go
// User types
const (
    UserTypePlatform   = "platform"     // AYSTEK team
    UserTypeTenant     = "tenant"       // Tenant user (UI login)
    UserTypeAPIClient  = "api_client"   // Third-party API access
)

type User struct {
    ID              primitive.ObjectID `bson:"_id,omitempty"`
    TenantID        primitive.ObjectID `bson:"tenant_id"`
    SiteID          *primitive.ObjectID `bson:"site_id"` // null = all sites
    
    Type            string `bson:"type"` // platform|tenant|api_client
    
    // For api_client type:
    Name            string `bson:"name"` // "Stripe Integration", "Mobile App"
    Description     string `bson:"description"`
    ClientID        string `bson:"client_id"` // unique, indexed
    ClientSecretRef string `bson:"client_secret_ref"` // KMS pointer, never stored plaintext
    
    AllowedScopes   []string `bson:"allowed_scopes"` // ["products:read", "orders:create"]
    AllowedRoutes   []RoutePermission `bson:"allowed_routes"` // explicit route + method list
    
    IPWhitelist     []string `bson:"ip_whitelist"` // optional; if empty = allow all
    RateLimitRPS    int `bson:"rate_limit_rps"` // requests per second
    QuotaDaily      int `bson:"quota_daily"` // max requests per day
    
    Status          string `bson:"status"` // active|inactive|revoked
    LastKeyRotation time.Time `bson:"last_key_rotation"`
    
    CreatedAt       time.Time `bson:"created_at"`
    UpdatedAt       time.Time `bson:"updated_at"`
    Metadata        map[string]interface{} `bson:"metadata"` // tags, contact info, etc.
}

type RoutePermission struct {
    Method  string `bson:"method"` // GET, POST, PUT, DELETE, PATCH
    Path    string `bson:"path"`   // /products, /orders, etc.
    Scopes  []string `bson:"scopes"` // scopes required for this route
}
```

### 2. API Usage Meter (Quota Tracking)

**Location:** `internal/modules/integrations/model/api_usage_meter.go`

```go
type APIUsageMeter struct {
    ID              primitive.ObjectID `bson:"_id,omitempty"`
    APIUserID       primitive.ObjectID `bson:"api_user_id"`
    TenantID        primitive.ObjectID `bson:"tenant_id"`
    
    // Daily window (reset at UTC midnight)
    Date            string `bson:"date"` // "2026-03-08"
    RequestCount    int `bson:"request_count"`
    SuccessCount    int `bson:"success_count"`
    ErrorCount      int `bson:"error_count"`
    
    // Rate limit bucket (Redis + periodic Mongo sync)
    TokensRemaining int `bson:"tokens_remaining"` // derived from Redis
    
    CreatedAt       time.Time `bson:"created_at"`
    UpdatedAt       time.Time `bson:"updated_at"`
}
```

### 3. API Call Audit Log

**Location:** `internal/modules/integrations/model/api_audit_log.go`

```go
type APIAuditLog struct {
    ID              primitive.ObjectID `bson:"_id,omitempty"`
    APIUserID       primitive.ObjectID `bson:"api_user_id"`
    TenantID        primitive.ObjectID `bson:"tenant_id"`
    
    Method          string `bson:"method"`
    Path            string `bson:"path"`
    QueryParams     string `bson:"query_params"` // log (no sensitive data)
    ResponseStatus  int `bson:"response_status"`
    
    RequestSize     int `bson:"request_size"` // bytes
    ResponseSize    int `bson:"response_size"`
    LatencyMS       int64 `bson:"latency_ms"`
    
    ErrorMessage    string `bson:"error_message"` // if error
    ClientIP        string `bson:"client_ip"`
    UserAgent       string `bson:"user_agent"`
    
    CreatedAt       time.Time `bson:"created_at"`
}
```

---

## OAuth2 Implementation Details

### 1. Token Generation Endpoint

**Endpoint:** `POST /auth/oauth2/token`

**Request:**
```json
{
  "grant_type": "client_credentials",
  "client_id": "api_client_xyz",
  "client_secret": "secret_abc123xyz",
  "scope": "products:read orders:read orders:create"
}
```

**Response (200 OK):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 900,
  "refresh_token": "ref_xyz123...",
  "scope": "products:read orders:read orders:create"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "invalid_client",
  "error_description": "Client authentication failed"
}
```

**Errors:**
- `invalid_request`: Missing required parameters
- `invalid_client`: client_id/secret mismatch or client revoked
- `invalid_scope`: Requested scope exceeds allowed scope
- `unauthorized_client`: Client not allowed for this grant type
- `unsupported_grant_type`: Only client_credentials supported for API clients

### 2. JWT Structure (access_token)

```
Header:
{
  "alg": "HS256",
  "typ": "JWT",
  "kid": "current_key_id"
}

Payload (Claims):
{
  "iss": "https://api.aystek.local",
  "sub": "api_client_xyz",
  "aud": ["api.aystek.local"],
  "iat": 1709959200,
  "exp": 1709960100,
  "scope": "products:read orders:read orders:create",
  "tenant_id": "tenant_abc123",
  "site_id": "site_def456", // null if all sites
  "api_user_id": "user_xyz",
  "type": "api_client"
}
```

### 3. Token Validation in Handler Middleware

```go
// Middleware: extract + validate JWT
func APIGatewayAuthMiddleware(c *gin.Context) {
    auth := c.GetHeader("Authorization")
    if auth == "" {
        c.JSON(401, gin.H{"error": "missing_authorization"})
        c.Abort()
        return
    }
    
    // Parse "Bearer <token>"
    parts := strings.Split(auth, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        c.JSON(400, gin.H{"error": "invalid_authorization_header"})
        c.Abort()
        return
    }
    
    token := parts[1]
    claims, err := jwt.ParseAndValidate(token, publicKey)
    if err != nil {
        c.JSON(401, gin.H{"error": "invalid_token"})
        c.Abort()
        return
    }
    
    // Store in context
    c.Set("api_user_id", claims.APIUserID)
    c.Set("tenant_id", claims.TenantID)
    c.Set("site_id", claims.SiteID)
    c.Set("scopes", claims.Scope)
    c.Next()
}
```

---

## Route Authorization & Permission Mapping

### Route Whitelist Model

Each API user has an `allowed_routes` list:

```go
type RoutePermission struct {
    Method  string   // GET, POST, PUT, DELETE, PATCH
    Path    string   // /products, /orders, /inventory, etc.
    Scopes  []string // required scope(s) for this route
}

// Example:
user.AllowedRoutes = []RoutePermission{
    {Method: "GET", Path: "/products", Scopes: []string{"products:read"}},
    {Method: "GET", Path: "/products/:id", Scopes: []string{"products:read"}},
    {Method: "POST", Path: "/orders", Scopes: []string{"orders:create"}},
    {Method: "GET", Path: "/orders", Scopes: []string{"orders:read"}},
    {Method: "GET", Path: "/orders/:id", Scopes: []string{"orders:read"}},
    // Explicitly excluded:
    // NOT /admin/*, /users/*, /settings/*
}
```

### Route Authorization Middleware

```go
func RouteAuthorizationMiddleware(c *gin.Context) {
    apiUserID := c.GetString("api_user_id")
    scopes := strings.Split(c.GetString("scopes"), " ")
    method := c.Request.Method
    path := c.Request.URL.Path
    
    // Find API user from DB/cache
    apiUser, _ := getUserByID(apiUserID)
    
    // Check if route + method is whitelisted
    allowed := false
    requiredScopes := []string{}
    for _, route := range apiUser.AllowedRoutes {
        if route.Method == method && matchPath(route.Path, path) {
            allowed = true
            requiredScopes = route.Scopes
            break
        }
    }
    
    if !allowed {
        c.JSON(403, gin.H{"error": "endpoint_not_permitted"})
        c.Abort()
        return
    }
    
    // Check if user has required scopes
    for _, required := range requiredScopes {
        if !contains(scopes, required) {
            c.JSON(403, gin.H{"error": "insufficient_scope"})
            c.Abort()
            return
        }
    }
    
    c.Next()
}
```

---

## Rate Limit & Quota Strategy

### Token Bucket Algorithm (Redis)

```
Per API user: TOKEN_BUCKET = "rate_limit:{api_user_id}"

Refill strategy:
  - RPS (requests per second): fixed interval refill
  - Burst capacity: 10 × RPS (to allow small bursts)
  - Redis TTL: match daily quota window

Example (RPS=10):
  - Burst capacity: 100 tokens
  - Refill rate: 10 tokens per second
  - When empty: reject with 429 + Retry-After header
```

### Quota Window (Daily)

```
Daily quota: max requests per UTC calendar day

Tracking (Redis + Mongo rollup):
  - Redis counter: "quota:{api_user_id}:{date}"
  - TTL: 24h + 1h buffer
  - Midnight UTC: close day, store to Mongo, reset counter

Response (when limit hit):
  HTTP/1.1 429 Too Many Requests
  X-RateLimit-Limit: 1000
  X-RateLimit-Remaining: 0
  X-RateLimit-Reset: 1709961600
  Retry-After: 3600
  
  {
    "error": "quota_exceeded",
    "message": "Daily quota (1000 requests) exhausted",
    "reset_at": "2026-03-09T00:00:00Z"
  }
```

### Rate Limit Middleware

```go
func RateLimitMiddleware(c *gin.Context) {
    apiUserID := c.GetString("api_user_id")
    
    // Check RPS bucket
    tokens, err := rateLimit.Allow(apiUserID, "rps")
    if !tokens {
        retryAfter := calculateRetryAfter(apiUserID)
        c.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
        c.JSON(429, gin.H{
            "error": "rate_limit_exceeded",
            "retry_after": retryAfter,
        })
        c.Abort()
        return
    }
    
    // Check daily quota
    quota, err := meter.GetDailyQuota(apiUserID)
    if quota.RequestCount >= quota.DailyLimit {
        c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", tomorrowUTC))
        c.JSON(429, gin.H{
            "error": "quota_exceeded",
            "reset_at": tomorrow.Format(time.RFC3339),
        })
        c.Abort()
        return
    }
    
    c.Next()
}
```

---

## Security Measures

### 1. Client Secret Management

**NEVER store plaintext secrets in Mongo.**

```go
type User struct {
    // Store only a reference to KMS/Vault
    ClientSecretRef string `bson:"client_secret_ref"`
    // Example: "kms:aws/us-east-1/key-id/secret-abc123"
}

// On token generation:
func generateToken(clientID, clientSecretFromRequest string) (tokenOrError) {
    // 1. Get secret KMS reference from User
    user, _ := findByClientID(clientID)
    
    // 2. Call KMS/Vault with reference
    actualSecret, err := kms.Decrypt(user.ClientSecretRef)
    
    // 3. Compare
    if actualSecret != clientSecretFromRequest {
        return 401, "invalid_client"
    }
    
    return // OK, generate token
}
```

### 2. Secret Rotation

```go
// Endpoint: POST /admin/api-clients/{id}/rotate-secret
func RotateAPIClientSecret(c *gin.Context) {
    apiClientID := c.Param("id")
    
    // Generate new secret
    newSecret := generateSecureRandom(32)
    
    // Encrypt and store KMS reference
    kmsRef := kms.Encrypt(newSecret, metadata)
    
    // Update User
    user.LastKeyRotation = time.Now()
    user.ClientSecretRef = kmsRef
    updateUser(user)
    
    // Return new secret (only once!) to UI
    // User must save it; we don't store plaintext
    c.JSON(200, gin.H{
        "client_secret": newSecret,
        "expires_in": "2026-03-15T08:00:00Z",
        "warning": "Copy and save this secret securely. It will not be shown again.",
    })
}
```

### 3. IP Whitelist (Optional)

```go
type User struct {
    IPWhitelist []string // ["203.0.113.0/24", "203.0.113.50"]
}

// In middleware:
func CheckIPWhitelist(c *gin.Context) {
    apiUser := getAPIUser(c)
    
    if len(apiUser.IPWhitelist) == 0 {
        c.Next() // No restriction
        return
    }
    
    clientIP := c.ClientIP()
    if !isIPInWhitelist(clientIP, apiUser.IPWhitelist) {
        c.JSON(403, gin.H{"error": "client_ip_not_whitelisted"})
        c.Abort()
        return
    }
    
    c.Next()
}
```

### 4. Audit Trail

All API calls logged with tenant/user/endpoint/status/latency/error context.

---

## Frontend: Admin UI for API Clients

### Screens to Implement

#### Screen 1: API Clients List
- Table: Name, Status, Created At, Last Activity, Usage (% quota)
- Actions: Create, Edit, Rotate Key, Revoke, Usage Logs
- Filters: Status (active/inactive), Tenant, Created Date
- Bulk: Revoke multiple

#### Screen 2: Create/Edit API Client
- Form fields:
  - Name (required, unique per tenant)
  - Description (optional)
  - Allowed Routes (multi-select table)
    - ✅ GET /products
    - ✅ POST /orders
    - ✅ GET /orders
    - ❌ GET /admin/* (disabled)
  - Scopes (auto-generated from allowed routes or manual select)
  - IP Whitelist (optional, comma-separated CIDR)
  - Rate Limit RPS (default 10)
  - Daily Quota (default 10000)
  - Status (active/inactive)
  - Metadata (tags, contact person, etc.)
- Actions: Save, Test Connection, Reset Secret
- Validation: Ensure routes + scopes are feasible

#### Screen 3: API Client Details & Secret Management
- Display:
  - Client ID (copyable)
  - Client Secret (hidden; show on rotate only)
  - Status, Created At, Last Activity
  - Allowed Routes (list view)
  - Current Rate Limit & Quota settings
- Actions:
  - Rotate Secret (modal: "Are you sure?" → newSecret shown once)
  - Revoke (immediate; no grace period)
  - Usage Logs (link to audit log screen)
  - Edit (go to edit form)

#### Screen 4: API Client Usage Logs
- Table: Timestamp, Method, Path, Status, Latency, Response Size, Client IP, Error
- Filters: Date range, Method, Path, Status Code (200/400/401/403/429/500)
- Charts:
  - Daily request trend
  - Status code distribution (200 vs 4xx vs 5xx)
  - Top endpoints called
  - Error rate over time
- Export: CSV (for audit export)

#### Screen 5: Quota & Rate Limit Monitor
- Current quota usage (pie chart: used/remaining)
- RPS meter (gauge: current RPS vs limit)
- Daily quota reset countdown
- Alert configuration:
  - "Alert when quota reaches X%"
  - "Alert when RPS exceeds Y"
- Historical quota usage (last 30 days)

### UI/UX Design Principles

**Enterprise Grade:**
- Consistent with AYSTEK's modern, minimal design (Tailwind gradient headers)
- Dark mode friendly
- Responsive (mobile support, but primary use is desktop)
- Accessibility (WCAG 2.1 AA)

**Information Architecture:**
- Breadcrumb: Dashboard → Security → API Clients
- Clear distinction between "Create" and "Configure"
- Show warnings for sensitive actions (Revoke, Rotate Secret)
- Copy-to-clipboard for client_id
- One-time secret display with manual save requirement

**Validation & Error Handling:**
- Real-time route conflict detection
- Scope mismatch warnings
- IP format validation (CIDR notation guide)
- Clear error messages on failed API operations

---

## Implementation Plan (MVP 4–6 months)

### Phase 1: Data Model & Backend Auth (Weeks 1–6)

- [x] Extend User model (api_client type, scopes, routes, limits)
- [x] Create APIUsageMeter model
- [x] Create APIAuditLog model
- [x] Implement OAuth2 client_credentials endpoint (`POST /auth/oauth2/token`)
- [x] Implement JWT generation and validation
- [x] KMS/Vault integration for secret storage
- [ ] Unit tests (auth flow, token validation, scope enforcement)

### Phase 2: Rate Limiting (Redis) + Route Authorization (Weeks 7–10) ✅ COMPLETE

- [x] Redis token bucket implementation
- [ ] Daily quota tracking (Redis counter + Mongo rollup) ⏳ **Enhancement: Scheduled for deployment alongside Quota Monitor UI**
- [x] Rate limit middleware
- [x] Quota enforcement middleware
- [x] Route authorization middleware (method + path + scope check)
- [x] Tests (429 responses, Retry-After header, quota reset)

**Completed Components:**
- ✅ `internal/modules/auth/ratelimit/redis_store.go` - Token bucket with RPS + daily quota enforcement
- ✅ `APIClientRateLimit()` middleware - Enforces 10 RPS + 100,000 daily quota per API client
- ✅ `APIClientRouteAuthorization()` middleware - Validates route permissions + OAuth2 scope requirements
- ✅ HTTP 429 responses with Retry-After header for RPS exhaustion
- ✅ HTTP 429 responses with quota_resets_at timestamp for daily quota exhaustion

**Test Coverage:**
- ✅ RPS limit exceeded returns 429 with retry_after=1
- ✅ Daily quota exceeded returns 429 with quota reset timestamp
- ✅ Route not in whitelist returns 403 route_not_permitted
- ✅ Insufficient OAuth2 scope returns 403 insufficient_scope

**Remaining Enhancement (Non-blocking for Phase 3):**
- ⏳ Mongo daily quota rollup for historical analytics (will be added in Phase 4 with usage dashboard)

### Phase 3: Audit Trail + Observability (Weeks 11–14) ✅ COMPLETE

- [x] Audit logging middleware
- [x] API audit log model + Mongo repository/service
- [x] Prometheus metrics (requests, latency, auth failures, rate-limit rejections, route denials)
- [x] Internal API route catalog (document all endpoints)
- [x] Tests (unauthorized access, scope mismatch, audit trail)

**Completed Components:**
- ✅ `internal/modules/auth/audit/` - Full audit module with model/repository/service layers
- ✅ `internal/modules/auth/metrics/api_gateway_metrics.go` - 5 Prometheus metrics types
- ✅ `APIClientAuditTrail()` middleware - Records all requests with tenant/user/method/path/status/latency
- ✅ Middleware test suite - 11 comprehensive test cases covering auth, rate-limit, route-authorization, audit scenarios
- ✅ `API_GATEWAY_ROUTE_CATALOG.md` - Complete route documentation with schemas, examples, error codes

**Test Results:**
```
=== ALL 11 TESTS PASS ===
TestAPIClientAuthRequired_MissingToken
TestAPIClientAuthRequired_InvalidTokenFormat
TestAPIClientAuthRequired_InvalidClientType
TestAPIClientAuthRequired_Success
TestAPIClientRateLimit_RPSExceeded
TestAPIClientRateLimit_DailyQuotaExceeded
TestAPIClientRouteAuthorization_RouteNotPermitted
TestAPIClientRouteAuthorization_InsufficientScope
TestAPIClientRouteAuthorization_Success
TestAPIClientAuditTrail_RecordsSuccessfulRequest
TestAPIClientAuditTrail_RecordsFailedRequest
```

**Integration Status:**
- ✅ Middleware chain wired in `cmd/api/main.go` lines 1535-1545
- ✅ API gateway routes exposed: GET/POST /products, GET/POST /orders
- ✅ Audit logs stored in `api_audit_logs` collection with indexes on tenant_id+created_at
- ✅ Metrics exported on `/metrics` endpoint for Prometheus scraping

### Phase 4: Frontend Admin Screens (Weeks 15–20) ✅ COMPLETE

- [x] API Clients List screen
- [x] Create/Edit API Client screen
- [x] API Client Details & Secret Management
- [x] Usage Logs screen
- [x] Quota & Rate Limit Monitor
- [x] Tests (form validation, API integration)

**Implemented UI Components (Next.js):**
- ✅ `node/aystek-core-next/app/(dashboard)/admin/api-clients/page.tsx` - API client list, filtering, status actions (suspend/reactivate/revoke)
- ✅ `node/aystek-core-next/app/(dashboard)/admin/api-clients/create/page.tsx` - Create form with scopes, routes, limits, IP whitelist
- ✅ `node/aystek-core-next/app/(dashboard)/admin/api-clients/[id]/page.tsx` - Detail/edit, secret rotation, revoke flow
- ✅ `node/aystek-core-next/app/(dashboard)/admin/api-clients/[id]/audit-logs/page.tsx` - Audit log explorer with filters
- ✅ `node/aystek-core-next/app/(dashboard)/admin/api-clients/quota-monitor/page.tsx` - Quota and rate-limit monitor with threshold alerts/reset countdown
- ✅ `node/aystek-core-next/lib/api/api-clients.ts` - Frontend API client service layer and typings
- ✅ `node/aystek-core-next/lib/api/api-client-validation.ts` - Shared form/IP validation helpers

**Backend APIs Completed for Phase 4:**
- ✅ `go/aystek-core/internal/modules/auth/handler/api_client_admin.go` - Admin CRUD, secret rotation, scope/route catalog, audit logs
- ✅ `go/aystek-core/internal/modules/auth/routes.go` - `/admin/api-clients/*` route registration with module permissions
- ✅ `go/aystek-core/cmd/api/main.go` - Handler initialization and admin route wiring

**Validation & Tests:**
- ✅ ESLint clean for all modified API client frontend files
- ✅ Go compile/tests pass for auth + API main package
- ✅ `tests/api-client-validation.test.ts` - Form/IP validation tests
- ✅ `tests/api-clients-api.test.ts` - API response mapping/integration tests

### Phase 5: Operations & Security (Weeks 21–24) ✅ COMPLETE

- [x] Secret manager provider selection from environment (`inline`, `aws_kms`)
- [x] AWS KMS-backed OAuth2 secret manager (`awskms:` reference format)
- [x] Production hardening guardrails
  - explicit `OAUTH2_SECRET_PROVIDER` requirement in production
  - inline provider disabled by default in production (explicit opt-in only)
- [x] Provider/router tests for production safeguards and legacy secret compatibility
- [x] Secret rotation runbook (`user-guide/05-operations/API_GATEWAY_SECRET_ROTATION_RUNBOOK.md`)
- [x] Audit log retention and export
  - `api_audit_logs` TTL index with `API_GATEWAY_AUDIT_RETENTION_DAYS` (default 90)
  - `GET /admin/api-clients/:id/audit-logs/export?format=csv|json`
- [x] KVKK compliance checks (data classification, third-party agreement)
  - `user-guide/03-compliance/API_GATEWAY_KVKK_GDPR_COMPLIANCE_MATRIX.md`
- [x] Security audit (SSRF, injection, replay attack mitigation)
  - `user-guide/05-operations/API_GATEWAY_SECURITY_AUDIT_CHECKLIST.md`
- [x] SLA/SLO definition
  - `user-guide/05-operations/API_GATEWAY_SLO_SLA_ALERTING.md`
- [x] Monitoring and alerting (Prometheus metrics, Grafana dashboard)
  - `go/aystek-core/deployments/prometheus/alerts.yml`
  - `go/aystek-core/deployments/grafana/dashboards/api-gateway-third-party.json`
- [x] Documentation (API client onboarding guide, best practices)
  - `user-guide/04-api/API_CLIENT_ONBOARDING_BEST_PRACTICES.md`

### Phase 6: Testing, Hardening & Go-Live (Weeks 25–26) 🚧 IN PROGRESS

- [x] Integration tests (multi-tenant isolation, concurrent requests)
  - `go/aystek-core/tests/api_gateway_integration_test.go`
- [x] Load tests (rate limit accuracy, quota fairness)
  - `go/aystek-core/tests/api_gateway_load_test.go`
- [x] Security tests (secret rotation, revocation, IP whitelist bypass attempts)
  - `go/aystek-core/tests/api_gateway_security_test.go`
- [x] Pilot rollout with internal team (runbook/checklist prepared)
  - `plans/PHASE_6_INTERNAL_PILOT_RUNBOOK.md`
- [x] Partner onboarding (first third-party API client) runbook/checklist prepared
  - `plans/PHASE_6_PARTNER_ONBOARDING_CHECKLIST.md`

Operational execution window (internal pilot traffic + first partner production calls) is pending business scheduling.

---

## Security Checklist

- [ ] All secrets encrypted at rest (KMS/Vault, never plaintext Mongo)
- [ ] TLS 1.3+ for all API communication
- [ ] Secret rotation enforced (min 90 days; max 365 days)
- [ ] IP whitelist validation (CIDR notation, no bypass)
- [ ] Audit logging (every API call with tenant/user/endpoint/status)
- [ ] Rate limit evasion prevention (consistent per-user tracking)
- [ ] Token expiry enforcement (15 min access_token, 7 day refresh_token)
- [ ] Scope binding to routes (no scope inflation)
- [ ] Tenant/site isolation in response filtering
- [ ] KVKK compliance: data classification, third-party SLA integration
- [ ] GDPR Art. 32 compliance: encryption in transit, secure key management, audit trail

---

## Compliance & Governance

### KVKK (Turkish Personal Data Protection Law)

- **Data Classification:** Identify which API routes expose PII
  - Example: `/customers/{id}` → contains name, email, phone (PII)
  - Mark as "pii=true" in route catalog
- **Third-Party Agreement:** If API client is external vendor, KVKK "data processor" agreement required
- **Data Residency:** Document if data stays in Turkey (KVKK Kurumu requirement)

### GDPR (EU regulation)

- **Art. 32 (Security):** Implement encryption, audit trail, access control ✅
- **Art. 28 (Processor):** If third-party API client processes EU citizen data, formal DPA (Data Processing Agreement) required

### Tracking

- [ ] API Client veri sınıflandırması tablosu (PII/non-PII routes)
- [ ] Third-party SLA template (KVKK data processor clause)
- [ ] Data residency policy
- [ ] Encryption + audit checklist (GDPR Art. 32)

---

## Monitoring & SLO/SLA

### Metrics (OpenTelemetry)

```
integrations_api_requests_total{
  api_user_id,
  method,
  path,
  status_code,
  tenant_id
}

integrations_api_rate_limit_exceeded_total{
  api_user_id,
  reason (rps|daily_quota)
}

integrations_api_latency_ms_p95{
  api_user_id,
  path
}

integrations_api_secret_rotation_age_days{
  api_user_id
}
```

### SLA/SLO Definitions

- **Availability:** API gateway + auth layer < 0.1% error rate monthly
- **Latency:** p95 API response < 500ms (internal services only; external depends on vendor)
- **Quota Accuracy:** Token bucket test: ±0.5% fairness across users
- **Security:** Zero secret leakage incidents; all rotations < 7 days
- **Audit Trail Retention:** 90 days (compliance storage), 7 years (archive)

---

## Going Forward: Phases 2+

Once MVP is live:
- **GraphQL API wrapper** (optional; advanced consumers)
- **Webhook outbound** (API client subscribes to events)
- **API key (not OAuth2)** deprecation path (legacy support 6 months)
- **Cost tracking** per API client (if charging model introduced)
- **mTLS** support (for high-security enterprise partners)

---

## References

- RFC 6749: OAuth 2.0 Authorization Framework https://tools.ietf.org/html/rfc6749
- RFC 7519: JSON Web Token (JWT) https://tools.ietf.org/html/rfc7519
- RFC 6585: HTTP 429 Too Many Requests https://tools.ietf.org/html/rfc6585
- AWS API Gateway: https://docs.aws.amazon.com/apigateway/
- Kong: https://docs.konghq.com/gateway/
- Stripe API Security: https://stripe.com/docs/security
- GitHub API: https://docs.github.com/en/rest
- KVKK Kurumu Kılavuzları: https://www.kvkk.gov.tr
- GDPR Art. 32: https://gdpr-info.eu/art-32-gdpr/

---

**Document Status:** ✅ **STANDARDS COMPLIANCE COMPLETE** (Phase 6 Tests Complete + Menu/Translation/i18n Fixed - March 9, 2026)  
**Next Steps:** 1) Schedule production pilot window and execute runbook 2) Complete first partner go-live and capture sign-off evidence 3) **Standards Compliance:** Menu entry added to security group, 28-language translations added, all UI text converted to i18n keys
