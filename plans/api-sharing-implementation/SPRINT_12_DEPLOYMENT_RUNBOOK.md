# Sprint 12 - Deployment Runbook

**Prepared On:** 2026-03-11  
**System:** API Sharing  
**Strategy:** Blue/Green with fast rollback

> Status note: this runbook is a prepared rollout template. It does not imply deployment approval. Use it only together with an approved `SPRINT_12_PRODUCTION_READINESS_CHECKLIST.md` and attached staging evidence.

## 1. Pre-Deployment Checks

- Confirm production readiness checklist status (`SPRINT_12_PRODUCTION_READINESS_CHECKLIST.md`).
- Confirm latest container image digest and SBOM availability.
- Confirm DB migrations (if any) are backward compatible.
- Confirm Redis and Mongo health are green.
- Confirm alert channels and on-call engineer availability.

## 2. Release Inputs

- Version tag: `TBD`
- Deployment window: `TBD`
- Change set:
  - Operation/scope lookup cache optimization
  - Performance benchmark assets
  - Load test scenario assets

## 3. Blue/Green Procedure

1. Deploy new version to green environment.
2. Run smoke tests on green:
   - `/oauth/jwks`
   - `/oauth/introspect`
   - `/internal/api-sharing/policy-check`
   - `/api-registry/openapi.json`
3. Verify telemetry on green for 10-15 minutes.
4. Shift 10% traffic to green and monitor:
   - error rate
   - p95/p99 latency
   - throttled rate
5. If healthy, ramp to 50%, then 100%.
6. Keep blue environment warm until rollback window closes.

## 4. Health Criteria During Rollout

- Error rate remains `< 0.1%`
- p95 latency remains `< 500ms`
- p99 latency remains `< 1000ms`
- No critical alert remains open for >5 minutes

## 5. Rollback Plan

Trigger rollback immediately when any condition is met:

- sustained 5xx spike above threshold
- auth/token validation failures above baseline
- severe Redis/Mongo saturation with request degradation

Rollback steps:

1. Route 100% traffic back to blue.
2. Scale down green deployment.
3. Verify recovery on blue within 5 minutes.
4. Open incident ticket with timeline and metrics snapshot.
5. Start post-incident review before next rollout attempt.

## 6. Post-Deployment Validation

- Confirm request throughput and latency return to expected baselines.
- Confirm audit logs and usage counters continue to ingest correctly.
- Confirm no increase in scope-denied false positives.
- Confirm no unexpected token revocation or introspection failures.

## 7. Incident Quick Actions

### Token/Authentication Failures

- Validate JWT signing key availability.
- Validate JWKS endpoint response.
- Check token revocation store health.

### Rate Limiter or Quota Anomalies

- Check Redis latency, error rate, and key eviction pressure.
- Validate rate-limit key churn and TTL behavior.
- If required, temporarily apply safer per-client limits.

### Policy-Check Latency Regression

- Inspect operation/scope lookup cache hit rate.
- Inspect Mongo query latency and lock contention.
- Evaluate temporary traffic shaping while tuning.

## 8. Runbook Ownership

- Primary: SRE On-Call
- Secondary: Backend API Sharing owner
- Escalation: Security + Infra leads
