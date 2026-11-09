# Sprint 11 - Capacity Planning

**Prepared On:** 2026-03-11  
**Module:** API Sharing (`policy-check`, OAuth, audit, usage)

## Objective

Estimate baseline resource and scaling requirements to sustain `1000 RPS` policy-check traffic while preserving latency and error SLOs.

## Planning Assumptions

- Traffic profile:
  - Baseline: 1000 RPS sustained
  - Peak burst: up to 1500 RPS for short windows
- Traffic mix:
  - `policy-check`: 70%
  - token/introspection/revoke: 20%
  - registry/audit/admin endpoints: 10%
- SLO targets:
  - `p95 < 500ms`
  - `p99 < 1000ms`
  - error rate `< 0.1%`

## Service Throughput Notes

From Sprint 11 benchmark and load evidence:

- Cache keeps repeated operation/scope lookups effectively CPU-bound inside the service.
- The dominant bottleneck discovered in end-to-end testing was not the cacheable Mongo path; it was incorrect route placement behind the full protected middleware chain.
- After remounting `POST /internal/api-sharing/policy-check` onto a lightweight internal route, measured local throughput and latency improved materially.

Measured evidence:

| Run | Request rate | p95 | p99 | Error rate | Throttle rate |
|---|---:|---:|---:|---:|---:|
| Pre-fix (`api-sharing-s11-20260311T130040Z`) | 61.25 req/s | 21.97 s | 22.44 s | 0.00% | 0.00% |
| Post-fix (`api-sharing-s11-20260311T131424Z`) | 987.12 req/s | 302.39 ms | 1.34 s | 0.00% | 0.00% |

Decision:

- Keep operation/scope lookup cache enabled in all non-test environments.
- Keep internal policy-check off the protected middleware chain.
- Use staging load tests to finalize Redis/JWT/Mongo tuning under mixed traffic, but use the corrected local baseline for immediate pod sizing and cost-envelope planning.

## Final Capacity Model

### API Sharing Service (Go)

- **Replicas:** start at 4 pods
- **Per-pod requests:** `500m CPU`, `512Mi memory`
- **Per-pod limits:** `1000m CPU`, `1Gi memory`
- **HPA target:** CPU 60%, memory 70%
- **HPA range:** min 4, max 12
- **Operational note:** one corrected local instance already reached `987.12 req/s`; keeping 4 pods as baseline provides failover and mixed-traffic headroom while the `p99` tail is still being tuned.
- **Pre-scale trigger for staging/production rehearsals:** add 2 replicas if CPU holds above `70%` for 5 minutes or if observed `p95` rises above `350ms` during a sustained window.

### Redis (Rate Limit + Quota Counters)

- **Mode:** managed HA (primary + replica) or equivalent.
- **Memory baseline:** 2 GiB.
- **Current finding:** corrected policy-check baseline is not Redis-bound; Redis should remain out of the hot path until Redis-backed rate-limit/quota enforcement is explicitly enabled.
- **Connection guidance:** size for roughly `100-150` client connections per app pod when rate/quota checks are enabled.
- **Alert thresholds:** memory usage `>65%`, p95 Redis latency `>5ms`, evictions `>0`.

### MongoDB

- Ensure indexes for hottest paths remain healthy and covered:
  - `api_operations.operation_id`
  - `api_operations.method + path`
  - `api_scopes.scope_id`
  - grant and audit indexes per schema design doc
- Operation/scope cache should keep most policy-check reads off Mongo in steady state.
- Reserve write headroom for grant/admin writes and selected audit events, not for a full success-log write on every allowed policy-check.

## Storage Growth Estimate (Audit Logs)

### Rejected model: full success audit for every policy-check

- Approx log size: `1.5 KB` per request.
- Daily requests at 1000 RPS: `86,400,000`.
- Raw/day: ~`123.6 GB`.
- 90-day soft-retention raw envelope: ~`11.1 TB`.

This model is not cost-efficient for the production baseline and should not be used on the primary Mongo cluster.

### Recommended model: detailed deny/error/admin logs + sampled success logs

- Keep full-detail documents for deny, throttle, auth-failure, admin, and token-management events.
- Keep aggregated usage counters for successful allowed policy-checks.
- If success sampling is required, cap it at `1%` before hot storage.
- At `1%` success sampling, hot success-log volume falls to roughly `1.24 GB/day` before compression, or about `111 GB` over 90 days.

## Final Cloud Cost Projection

Illustrative managed-cloud monthly envelope for the corrected baseline:

| Component | Baseline assumption | Monthly envelope |
|---|---|---:|
| App compute | 4 Go pods, `500m/512Mi` requests, modest average utilization | `$120-$220` |
| Redis HA | 2 GiB managed tier | `$70-$140` |
| MongoDB | small managed replica set + ~`100-150 GB` hot data + backups | `$180-$350` |
| Observability | metrics, traces, sampled success logs, alerting | `$120-$300` |
| Buffer | network, backup, egress, variance | `$50-$150` |
| **Total steady-state envelope** | corrected-route baseline | **`$540-$1,160 / month`** |

Cost warning:

- If full success-audit documents are retained in hot storage for all policy-check requests, storage and log-ingestion cost can rise into multi-thousand-dollar monthly territory and invalidate this baseline.

## Cost Planning Inputs

- Compute: based on `4-12` replica HPA range and mixed policy-check/token load.
- Redis: HA baseline sized for token revocation now and rate/quota counters later.
- Mongo: storage + IOPS sized for sampled/sparse audit writes, not full unsampled success logging.
- Observability: keep high-cardinality request detail sampled; use counters for bulk success traffic.

## Sprint 11 Capacity Checklist

- [x] Initial capacity assumptions documented.
- [x] Service cache optimization reflected in model.
- [x] Corrected-route load-test data folded into sizing assumptions.
- [x] Final pod sizing + HPA profile proposed.
- [x] Final cloud cost projection prepared.
- [ ] Final cloud cost projection approved by Infra/SRE.

## Handoff to Sprint 12

Use this document as input for production readiness and deployment runbook sign-off.
