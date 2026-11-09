# TASK_024 — PMS Observability & Security Hardening

## Observability
- PMS admin route zincirine request log middleware eklendi.
- Log alanları:
  - `tenant_id`
  - `site_id`
  - `route`
  - `method`
  - `status`
  - `duration_ms`

- Basit metric counter eklendi (`internal/modules/pms/metrics`):
  - dashboard summary
  - billing preview/publish/allocate
  - reconciliation propose/confirm/reject

## Security Hardening
- Portal rate limit env-configurable:
  - `PMS_PORTAL_RATE_LIMIT_PER_MIN`
  - `PMS_PORTAL_RATE_LIMIT_WINDOW_SEC`
- Billing idempotency retention için TTL index eklendi:
  - `pms_billing_idempotency__created_ttl_30d`

## Test
- PMS modül testleri ve `cmd/api` derleme doğrulandı.
