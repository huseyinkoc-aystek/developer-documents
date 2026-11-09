# TASK 027 — PMS Observability v2

## Scope
Expanded PMS observability with request duration buckets, request id propagation, and action counters for new critical flows.

## Implemented
- Updated `internal/modules/pms/request_middleware.go`
  - adds `request_id` to PMS request log line (`X-Request-Id` / context fallback)
  - records request duration bucket metrics
- Extended `internal/modules/pms/metrics/metrics.go`
  - new counters:
    - exports create/download
    - webhook dispatch
  - duration buckets:
    - `<=100ms`
    - `<=500ms`
    - `>500ms`
- Hooked counters into:
  - export job create/download
  - webhook dispatch path

## Notes
- Existing PMS request structured log fields retained:
  - tenant_id
  - site_id
  - route
  - method
  - status
  - duration_ms
