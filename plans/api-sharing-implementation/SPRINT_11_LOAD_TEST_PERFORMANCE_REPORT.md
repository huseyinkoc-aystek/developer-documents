# Sprint 11 - Load Testing and Performance Report

**Sprint Window:** Planned for 2026-08-05 to 2026-08-18  
**Prepared On:** 2026-03-11  
**Module:** `internal/modules/integrations/api_sharing`

> Status note: this is an early evidence draft prepared before the planned sprint window. It must not be treated as final Sprint 11 closure.

## Scope

This report captures Sprint 11 performance work implemented and measured in this cycle:

- Added in-memory lookup cache for operation and scope reads in service layer.
- Added repeatable benchmark suite for lookup and scope validation paths.
- Added k6 load-test scenario for `POST /internal/api-sharing/policy-check` with 1000 RPS target profile.

## Implemented Changes

### 1. Service Optimization

- File: `/workspace/go/aystek-core/internal/modules/integrations/api_sharing/service/service.go`
- Added TTL-based lookup caches:
  - Operation cache (`operation_id` keyed)
  - Scope cache (`scope_id` keyed)
- Added cache population/invalidation on create/update/delete paths.
- Added bounded cache eviction (expired-first, then single-entry eviction when max size is reached).

### 2. Performance Test Assets

- Benchmarks file: `/workspace/go/aystek-core/internal/modules/integrations/api_sharing/service/performance_benchmark_test.go`
- k6 scenario: `/workspace/go/aystek-core/internal/modules/integrations/api_sharing/tests/load/policy_check_k6.js`
- k6 usage notes: `/workspace/go/aystek-core/internal/modules/integrations/api_sharing/tests/load/README.md`
- Runner hardening:
  - `run_staging_k6.sh` now exports `p(99)` in summary output via `--summary-trend-stats`
  - `run_first_staging_once.sh` now preserves summary generation even when a k6 threshold fails

### 3. Internal Route Path Correction

- File: `/workspace/go/aystek-core/cmd/api/main.go`
- `POST /internal/api-sharing/policy-check` was previously mounted on the full protected middleware chain.
- Route registration was moved from `protectedGroup` to a lightweight internal group so gateway-facing policy checks no longer pay auth/store/CSRF/permission overhead.
- This change turned the main bottleneck from router/middleware composition into tail-latency tuning.

## Benchmark Evidence

### Command

```bash
go test ./internal/modules/integrations/api_sharing/service \
  -run '^$' \
  -bench 'BenchmarkGetOperation_(With|Without)Cache(_SlowRepository)?$|BenchmarkValidateScopes_WithCache$' \
  -benchmem \
  -count=1
```

### Results (Service-Level)

| Benchmark | ns/op | allocs/op | Notes |
|---|---:|---:|---|
| `BenchmarkGetOperation_WithCache` | 96.23 | 1 | In-memory mock repository (fast path baseline) |
| `BenchmarkGetOperation_WithoutCache` | 90.53 | 1 | In-memory mock repository (baseline) |
| `BenchmarkGetOperation_WithCache_SlowRepository` | 91.48 | 1 | Simulated repository latency (`750us`) hidden by cache |
| `BenchmarkGetOperation_WithoutCache_SlowRepository` | 2,000,184 | 1 | Simulated repository latency on every request |
| `BenchmarkValidateScopes_WithCache` | 264.9 | 3 | Scope validation with warm cache |

## Interpretation

- In purely in-memory repository paths, cache overhead is negligible and expected to be close to baseline.
- Under simulated repository latency, cache provides a major gain:
  - Cached path: ~91 ns/op
  - Uncached path: ~2 ms/op
- This indicates the optimization target is correct for real deployments where DB/network latency dominates.

## 1000 RPS Validation Status

- [x] Service-level performance benchmark artifacts produced.
- [x] 1000 RPS k6 scenario defined with p95/p99/error thresholds.
- [x] Representative end-to-end local run executed against the corrected internal route.
- [ ] Full staging execution report for k6 scenario attached.
- [ ] End-to-end gateway + Redis + Mongo profile captured in staging.

## Executed Evidence (2026-03-11)

No active staging URL/token pair was available during this cycle, so Sprint 11 evidence was produced on an isolated local instance of the same repository running in `APP_ENV=staging`.

### A. Pre-fix bottleneck run (`api-sharing-s11-20260311T130040Z`)

This run exercised `POST /internal/api-sharing/policy-check` while the route was still mounted on the full protected middleware chain.

| Metric | Value | Target | Status |
|---|---:|---:|---|
| Request rate | 61.25 req/s | >=1000 req/s | fail |
| p95 latency | 21.97 s | <500 ms | fail |
| p99 latency | 22.44 s | <1000 ms | fail |
| Error rate | 0.00% | <0.1% | pass |
| Throttle rate | 0.00% | <1% | pass |
| Dropped iterations | 280799 (`895.83/s`) | minimize | severe saturation |

Artifacts:
- Run log: `/workspace/go/aystek-core/internal/modules/integrations/api_sharing/tests/load/artifacts/api-sharing-s11-20260311T130040Z_run.log`
- Summary JSON: `/workspace/go/aystek-core/internal/modules/integrations/api_sharing/tests/load/artifacts/api-sharing-s11-20260311T130040Z_summary.json`
- Summary markdown: `/workspace/go/aystek-core/internal/modules/integrations/api_sharing/tests/load/artifacts/api-sharing-s11-20260311T130040Z_summary_report.md`

### B. Post-fix corrected-route run (`api-sharing-s11-20260311T131424Z`)

This run exercised the corrected internal route registration after moving policy-check off the protected middleware chain.

| Metric | Value | Target | Status |
|---|---:|---:|---|
| Request rate | 987.12 req/s | >=1000 req/s | near miss |
| p95 latency | 302.39 ms | <500 ms | pass |
| p99 latency | 1.34 s | <1000 ms | fail |
| Error rate | 0.00% | <0.1% | pass |
| Throttle rate | 0.00% | <1% | pass |
| Dropped iterations | 3849 (`12.83/s`) | minimize | major improvement |

Artifacts:
- Run log: `/workspace/go/aystek-core/internal/modules/integrations/api_sharing/tests/load/artifacts/api-sharing-s11-20260311T131424Z_run.log`
- Summary JSON: `/workspace/go/aystek-core/internal/modules/integrations/api_sharing/tests/load/artifacts/api-sharing-s11-20260311T131424Z_summary.json`

## Interpretation

- The dominant bottleneck was not Mongo lookup latency after caching; it was incorrect route placement behind the full protected middleware stack.
- Moving `policy-check` to a lightweight internal route improved throughput from `61.25 req/s` to `987.12 req/s` and reduced p95 from `21.97s` to `302.39ms`.
- Current baseline clears `p95`, `error rate`, and `throttle rate` targets, narrowly misses the sustained request-rate target, and still fails the `p99 < 1000ms` gate.
- Sprint 11 is therefore materially advanced, but not production-ready until staging validation and p99 tail reduction are complete.

## Staging Execution Command Set (Prometheus/Grafana)

Quick one-shot flow:

```bash
cd /workspace/go/aystek-core
export BASE_URL="https://<staging-api-host>"
export AUTH_TOKEN="<staging-token>"
export K6_PROMETHEUS_RW_SERVER_URL="https://<prometheus-host>/api/v1/write"
internal/modules/integrations/api_sharing/tests/load/run_first_staging_once.sh
```

Notes:
- The one-shot flow now still generates summary artifacts even if a threshold fails.
- `run_staging_k6.sh` now exports `p(99)` into the summary JSON for future runs.

### 1. Execute k6 run with Prometheus remote write export

```bash
cd /workspace/go/aystek-core

export BASE_URL="https://<staging-api-host>"
export AUTH_TOKEN="<staging-token>"
export TENANT_ID="tenant-1"
export CLIENT_ID="partner-acme"
export ENVIRONMENT="staging"
export OUTPUT_MODE="prometheus-rw"
export TEST_ID="api-sharing-s11-$(date -u +%Y%m%dT%H%M%SZ)"

export K6_PROMETHEUS_RW_SERVER_URL="https://<prometheus-host>/api/v1/write"
export K6_PROMETHEUS_RW_TREND_STATS="p(90),p(95),p(99),avg,max"
export K6_PROMETHEUS_RW_PUSH_INTERVAL="5s"

internal/modules/integrations/api_sharing/tests/load/run_staging_k6.sh
```

### 2. Generate markdown summary artifact from k6 summary JSON

```bash
cd /workspace/go/aystek-core

SUMMARY_FILE="internal/modules/integrations/api_sharing/tests/load/artifacts/${TEST_ID}_summary.json"
internal/modules/integrations/api_sharing/tests/load/summarize_k6_results.sh "$SUMMARY_FILE"
```

### 3. Export Grafana dashboard + panel artifacts

```bash
cd /workspace/go/aystek-core

export GRAFANA_URL="https://<grafana-host>"
export GRAFANA_TOKEN="<grafana-api-token>"
export DASHBOARD_UID="<dashboard-uid>"
export PANEL_IDS="2,4,6,8"
export FROM_TS="now-30m"
export TO_TS="now"

internal/modules/integrations/api_sharing/tests/load/export_grafana_artifacts.sh
```

### 4. Prometheus query snippets

```promql
sum(rate(k6_http_reqs_total{testid="$TEST_ID"}[1m]))
```

```promql
avg_over_time(k6_http_req_duration_p95{testid="$TEST_ID"}[5m])
```

```promql
avg_over_time(k6_http_req_duration_p99{testid="$TEST_ID"}[5m])
```

```promql
avg_over_time(k6_http_req_failed_rate{testid="$TEST_ID"}[5m])
```

```promql
sum(rate(api_sharing_throttled_total{environment="staging"}[5m]))
/
sum(rate(api_sharing_requests_total{environment="staging"}[5m]))
```

Note: metric suffixes may vary by k6 output mode/version; adapt query names where needed.

## Bottleneck Review

- **Database query latency:** materially reduced for repeated operation/scope reads via service cache.
- **Router/middleware composition:** confirmed as the primary bottleneck and fixed via internal route remounting.
- **Redis rate-limit latency:** still pending dedicated staging profiling under mixed load.
- **JWT signature validation and mixed-token traffic:** still pending targeted staging profiling.
- **Tail latency:** remaining open performance issue is `p99`, not average or p95 latency.

## Next Actions

1. Execute the same scenario in real staging and capture Prometheus/Grafana evidence.
2. Reduce `p99` from `1.34s` to `<1.0s` using focused profiling of long-tail requests.
3. Validate Redis, Mongo, and HPA behavior under mixed policy-check + token traffic before Sprint 12 approval.
