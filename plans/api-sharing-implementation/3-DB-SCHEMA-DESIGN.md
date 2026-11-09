# API Sharing - MongoDB Schema & Indexes

## Collections Overview

| Collection | Purpose | TTL | Key Indexes |
|---|---|---|---|
| `api_operations` | API endpoint registry | None | operation_id, method+path, service_name |
| `api_scopes` | Permission scopes | None | scope_id, module |
| `api_grants` | Tenant-scoped grants | None | tenant_id, client_id, status |
| `api_consents` | User consents (KVKK) | 7 years | tenant_id, user_id, client_id, status |
| `api_audit_logs` | Access logs (KVKK) | 90 days | tenant_id, client_id, timestamp, deleted_at |
| `api_usage` | Usage metrics | None | tenant_id, client_id, period |
| `partner_dpa_records` | DPA tracking (KVKK) | None | tenant_id, client_id |

---

## 1. api_operations (Collection)

### Document Structure
```json
{
  "_id": ObjectId,
  "operation_id": "orders_list_v1",  // unique
  "method": "GET",
  "path": "/api/v1/orders",  // unique per method
  "service_name": "orders-service",
  "summary": "List all orders",
  "description": "Retrieve paginated list of orders with filters",
  "tags": ["orders", "commerce"],
  "required_scopes": ["orders.read"],
  "deprecated": false,
  "sunset_date": null,
  "rate_limit_override": {
    "requests_per_second": 50,
    "requests_per_minute": 1000,
    "burst_size": 100
  },
  "created_at": ISODate("2026-03-09T10:00:00Z"),
  "updated_at": ISODate("2026-03-09T10:00:00Z"),
  "deleted_at": null
}
```

### Indexes
```javascript
db.api_operations.createIndex({ operation_id: 1 }, { unique: true })
db.api_operations.createIndex({ method: 1, path: 1 }, { unique: true })
db.api_operations.createIndex({ service_name: 1 })
db.api_operations.createIndex({ deprecated: 1, deleted_at: 1 })
db.api_operations.createIndex({ tags: 1 })
db.api_operations.createIndex({ "required_scopes": 1 })
```

---

## 2. api_scopes (Collection)

### Document Structure
```json
{
  "_id": ObjectId,
  "scope_id": "orders.read",  // unique
  "description": "Read access to orders",
  "module": "orders",
  "action": "read",
  "is_global": true,  // or tenant-specific
  "created_at": ISODate("2026-03-09T10:00:00Z"),
  "updated_at": ISODate("2026-03-09T10:00:00Z")
}
```

### Indexes
```javascript
db.api_scopes.createIndex({ scope_id: 1 }, { unique: true })
db.api_scopes.createIndex({ module: 1, action: 1 })
db.api_scopes.createIndex({ is_global: 1 })
```

---

## 3. api_grants (Collection)

### Document Structure
```json
{
  "_id": ObjectId,
  "tenant_id": "tenant_123",
  "site_id": null,  // null = all sites
  "client_id": "partner_acme",
  "scopes": ["orders.read", "orders.write"],
  "allowed_operations": [],  // if empty, use scopes; if set, restrict to these ops
  "ip_allowlist": ["203.0.113.0/24", "192.168.0.0/16"],
  "rate_limit_policy": {
    "requests_per_second": 10,
    "requests_per_minute": 300,
    "burst_size": 20
  },
  "quota_policy": {
    "daily_quota": 10000,
    "monthly_quota": 200000,
    "max_payload_kb": 256
  },
  "processing_basis": "consent",  // or "contract", "legal_obligation"
  "consent_id": "con_xxxxx",  // if processing_basis="consent"
  "status": "active",  // "active", "suspended", "revoked"
  "expires_at": null,
  "created_at": ISODate("2026-03-09T10:00:00Z"),
  "updated_at": ISODate("2026-03-09T10:00:00Z"),
  "created_by": "user_admin_123",
  "revoked_at": null,
  "revoked_by": null,
  "revoked_reason": null,
  "deleted_at": null
}
```

### Indexes
```javascript
db.api_grants.createIndex({ tenant_id: 1, client_id: 1 })
db.api_grants.createIndex({ tenant_id: 1, status: 1 })
db.api_grants.createIndex({ client_id: 1 })
db.api_grants.createIndex({ "scopes": 1 })
db.api_grants.createIndex({ status: 1, expires_at: 1 })
db.api_grants.createIndex({ consent_id: 1 })
```

---

## 4. api_consents (Collection) - KVKK Compliance

### Document Structure
```json
{
  "_id": ObjectId,
  "tenant_id": "tenant_123",
  "client_id": "partner_acme",
  "user_id": "user_456",
  "scopes": ["orders.read", "shipments.write"],
  "processing_purpose": "Integration with ACME ERP System",
  "processing_basis": "consent",
  "consent_version": "1.0",
  "language": "tr_TR",
  "status": "active",  // "active", "withdrawn", "expired"
  "issued_at": ISODate("2026-03-09T10:00:00Z"),
  "expires_at": ISODate("2027-03-09T10:00:00Z"),
  "withdrawn_at": null,
  "ip_address_hash": "sha256_hash_of_203.0.113.50",
  "user_agent_hash": "sha256_hash_of_user_agent",
  "consent_url": "/consents/con_xxxxx/details",
  "audit_trail": [
    {
      "action": "ISSUED",
      "timestamp": ISODate("2026-03-09T10:00:00Z"),
      "ip_hash": "sha256_..."
    },
    {
      "action": "WITHDRAWN",
      "timestamp": ISODate("2026-03-15T14:30:00Z"),
      "reason": "User requested"
    }
  ],
  "deleted_at": null
}
```

### Indexes
```javascript
db.api_consents.createIndex({ tenant_id: 1, user_id: 1 })
db.api_consents.createIndex({ tenant_id: 1, client_id: 1, user_id: 1 }, { unique: true, sparse: true })
db.api_consents.createIndex({ user_id: 1, status: 1 })
db.api_consents.createIndex({ client_id: 1, status: 1 })
db.api_consents.createIndex({ status: 1, expires_at: 1 })
db.api_consents.createIndex({ deleted_at: 1 })
// TTL Index: 7 years (for KVKK retention)
db.api_consents.createIndex({ created_at: 1 }, { expireAfterSeconds: 220752000 })
```

---

## 5. api_audit_logs (Collection) - KVKK Compliance

### Document Structure
```json
{
  "_id": ObjectId,
  "tenant_id": "tenant_123",
  "client_id": "partner_acme",
  "user_id": null,  // if user-bound
  "operation_id": "orders_list_v1",
  "method": "GET",
  "path": "/api/v1/orders",
  "query_params_hash": "sha256_hash_of_query",
  "request_payload_hash": "sha256_hash",  // sensitive data not stored
  "response_status": 200,
  "response_size_bytes": 4096,
  "latency_ms": 245,
  "result": "ALLOW",  // "ALLOW", "DENY", "THROTTLE"
  "denial_reason": null,  // "SCOPE_DENIED", "IP_BLOCKED", "QUOTA_EXCEEDED"
  "scopes_checked": ["orders.read"],
  "ip_address_hash": "sha256_hash_203.0.113.x",
  "timestamp": ISODate("2026-03-09T10:00:00Z"),
  "trace_id": "W3C_trace_id",
  "deleted_at": null  // soft delete flag
}
```

### Indexes
```javascript
db.api_audit_logs.createIndex({ tenant_id: 1, timestamp: -1 })
db.api_audit_logs.createIndex({ client_id: 1, timestamp: -1 })
db.api_audit_logs.createIndex({ operation_id: 1, timestamp: -1 })
db.api_audit_logs.createIndex({ user_id: 1, timestamp: -1 }, { sparse: true })
db.api_audit_logs.createIndex({ result: 1, timestamp: -1 })
db.api_audit_logs.createIndex({ deleted_at: 1 })
// TTL Index: 90 days (for KVKK retention)
db.api_audit_logs.createIndex({ created_at: 1 }, { expireAfterSeconds: 7776000 })
```

---

## 6. api_usage (Collection)

### Document Structure
```json
{
  "_id": ObjectId,
  "tenant_id": "tenant_123",
  "client_id": "partner_acme",
  "operation_id": "orders_list_v1",
  "period": "2026-03",  // YYYY-MM for monthly aggregation
  "request_count": 4523,
  "error_count": 12,
  "throttled_count": 5,
  "total_latency_ms": 1104837,  // for avg calculation
  "created_at": ISODate("2026-03-01T00:00:00Z"),
  "updated_at": ISODate("2026-03-09T14:30:00Z")
}
```

### Indexes
```javascript
db.api_usage.createIndex({ tenant_id: 1, period: -1 })
db.api_usage.createIndex({ client_id: 1, period: -1 })
db.api_usage.createIndex({ tenant_id: 1, client_id: 1, operation_id: 1, period: -1 }, { unique: true })
db.api_usage.createIndex({ operation_id: 1 })
```

---

## 7. partner_dpa_records (Collection) - KVKK Compliance

### Document Structure
```json
{
  "_id": ObjectId,
  "tenant_id": "tenant_123",
  "client_id": "partner_acme",
  "dpa_version": "1.0",
  "partner_country": "DE",  // ISO-3166-1 (Germany)
  "scc_applied": true,  // Standard Contractual Clauses
  "signed_at": ISODate("2026-03-05T09:00:00Z"),
  "signed_by_email": "legal@acme.com",
  "expires_at": ISODate("2027-03-05T09:00:00Z"),
  "created_at": ISODate("2026-03-05T09:00:00Z"),
  "updated_at": ISODate("2026-03-05T09:00:00Z"),
  "deleted_at": null
}
```

### Indexes
```javascript
db.partner_dpa_records.createIndex({ tenant_id: 1, client_id: 1 })
db.partner_dpa_records.createIndex({ tenant_id: 1 })
db.partner_dpa_records.createIndex({ client_id: 1 })
db.partner_dpa_records.createIndex({ expires_at: 1 })
```

---

## Maintenance Jobs (KVKK Compliance)

### 1. TTL Auto-Deletion (MongoDB built-in)
- `api_audit_logs`: 90 days (via TTL index)
- `api_consents`: 7 years (via TTL index)

**Note:** MongoDB TTL works on `created_at` field. For soft-deleted records, use a separate cleanup job.

### 2. Soft Delete Cleanup Job (run weekly)
```javascript
// Hard delete soft-deleted audit logs older than 30 days after soft-delete
db.api_audit_logs.deleteMany({
  deleted_at: { $lt: new Date(Date.now() - 30*24*60*60*1000) }
})

// Hard delete soft-deleted consents (after natural TTL + grace period)
db.api_consents.deleteMany({
  deleted_at: { $lt: new Date(Date.now() - 90*24*60*60*1000) }
})
```

### 3. Consent Expiration Job (run daily)
```javascript
db.api_consents.updateMany(
  { expires_at: { $lt: new Date() }, status: "active" },
  { $set: { status: "expired" } }
)
```

### 4. Grant Expiration Job (run daily)
```javascript
db.api_grants.updateMany(
  { expires_at: { $lt: new Date() }, status: "active" },
  { $set: { status: "suspended", updated_at: new Date() } }
)
```

---

## Data Export/Anonymization (KVKK SAR/GDPR)

### Subject Access Request (SAR) - User data export
```javascript
// Fetch all records for user_id
db.api_consents.find({ user_id: "user_456" })
db.api_audit_logs.find({ user_id: "user_456" })
// Export to JSON for user
```

### Anonymization (KVKK deletion)
```javascript
// Anonymize: replace PII with hashes/nulls
db.api_audit_logs.updateMany(
  { user_id: "user_456" },
  {
    $set: {
      user_id: null,
      ip_address_hash: "ANONYMIZED",
      query_params_hash: "ANONYMIZED"
    }
  }
)
```

---

## Backup & Disaster Recovery Strategy

- **Full backup:** Daily (incremental)
- **Point-in-time recovery:** 30 days
- **Geo-replication:** Multi-region (optional, for HA)
- **Test DR quarterly**

---

## Performance Tuning

| Operation | Expected QPS | Latency Target |
|---|---|---|
| Create grant | 10 | 50ms |
| Check access (policy eval) | 1000 | 10ms |
| Log audit | 1000 | 5ms |
| Query usage report | 10 | 200ms |
| Export user data (SAR) | 1 | 1s |

**Optimization:**
- Use connection pooling (MongoDB Atlas, +100 connections)
- Compound indexes for common query patterns
- Caching for scopes/operations (Redis, 1hr TTL)

---

**Document Version:** 1.0  
**Last Updated:** 2026-03-09  
**Status:** Ready for Implementation
