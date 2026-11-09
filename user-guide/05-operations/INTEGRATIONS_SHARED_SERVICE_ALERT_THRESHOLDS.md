---
Title: Integrations Shared Service Alert Thresholds (Prometheus/Grafana)
Version: 1.0.0
Status: Active
Last Updated: 2026-03-08
Related Documents:
  - ./INTEGRATIONS_SHARED_SERVICE_RUNBOOK_QA_MATRIX.md
  - ./OPERATIONS_GUIDE.md
  - ../06-modules/INTEGRATIONS_PLATFORM_GUIDE.md
---

# Integrations Shared Service Alert Thresholds

Bu doküman Shared Service policy hattı için önerilen Prometheus/Grafana alarm eşiklerini tanımlar.

## 1. Metric Surface
Kod bazlı metrik:
1. `aystek_integrations_shared_service_usage_total{tenant_id,site_id,shared_service_id,grant_id,outcome}`

İlgili outcome örnekleri:
1. `execution_success`
2. `execution_failed`
3. `daily_quota_denied`
4. `rpm_denied`
5. `scope_denied`
6. `method_denied`
7. `path_denied`
8. `payload_guard_denied`
9. `mtls_denied`
10. `secret_binding_denied`
11. `grant_inactive`
12. `service_inactive`

## 2. Severity Model
1. `P1`: müşteri etkisi yüksek, anlık müdahale gerekir.
2. `P2`: işlev bozuluyor, çalışma saatinde müdahale gerekir.
3. `P3`: risk artıyor, backlog iyileştirme gerekir.

## 3. Alert Rules (PromQL)
### 3.1 Shared Service Error Ratio High (P1)
Amaç: Son 15 dakikada başarısızlık oranı %10+ ise alarm üret.

```promql
(
  sum(rate(aystek_integrations_shared_service_usage_total{outcome="execution_failed"}[15m]))
/
  clamp_min(sum(rate(aystek_integrations_shared_service_usage_total[15m])), 0.001)
) > 0.10
```

### 3.2 Policy Deny Burst (P2)
Amaç: Son 10 dakikada policy deny toplamı yüksekse alarm üret.

```promql
sum(increase(aystek_integrations_shared_service_usage_total{outcome=~".*denied|grant_inactive|service_inactive|connector_mismatch"}[10m])) > 50
```

### 3.3 Quota Pressure High (P2)
Amaç: Kota kaynaklı deny artışını yakala.

```promql
sum(increase(aystek_integrations_shared_service_usage_total{outcome=~"daily_quota_denied|rpm_denied"}[10m])) > 20
```

### 3.4 Secret/MTLS Misconfiguration (P1)
Amaç: Güvenlik policy hatalarını hızlı tespit et.

```promql
sum(increase(aystek_integrations_shared_service_usage_total{outcome=~"secret_binding_denied|mtls_denied"}[5m])) > 0
```

### 3.5 No Success With Traffic (P1)
Amaç: Trafik var ama success yoksa alarm üret.

```promql
sum(increase(aystek_integrations_shared_service_usage_total[15m])) > 30
and
sum(increase(aystek_integrations_shared_service_usage_total{outcome="execution_success"}[15m])) == 0
```

## 4. Suggested PrometheusRule
```yaml
groups:
  - name: integrations-shared-service-alerts
    rules:
      - alert: IntegrationsSharedServiceErrorRatioHigh
        expr: (sum(rate(aystek_integrations_shared_service_usage_total{outcome="execution_failed"}[15m])) / clamp_min(sum(rate(aystek_integrations_shared_service_usage_total[15m])), 0.001)) > 0.10
        for: 10m
        labels:
          severity: critical
          team: integrations
        annotations:
          summary: "Shared service execution failure ratio is above 10%"
          description: "Check executions, grant policies, and downstream availability."

      - alert: IntegrationsSharedServicePolicyDenyBurst
        expr: sum(increase(aystek_integrations_shared_service_usage_total{outcome=~".*denied|grant_inactive|service_inactive|connector_mismatch"}[10m])) > 50
        for: 5m
        labels:
          severity: warning
          team: integrations
        annotations:
          summary: "Shared service policy deny burst detected"
          description: "Review service/grant method-path-scope-limits and latest deployments."

      - alert: IntegrationsSharedServiceSecurityPolicyBlocked
        expr: sum(increase(aystek_integrations_shared_service_usage_total{outcome=~"secret_binding_denied|mtls_denied"}[5m])) > 0
        for: 2m
        labels:
          severity: critical
          team: integrations
        annotations:
          summary: "Shared service secret or mTLS policy is blocking traffic"
          description: "Verify secret bindings/values and mTLS cert-key configuration."
```

## 5. Grafana Dashboard Panels
1. `Shared Traffic`: total request rate (5m).
2. `Success vs Failed`: stacked outcome graph (`execution_success`, `execution_failed`).
3. `Policy Denied Breakdown`: deny reason table (`scope_denied`, `method_denied`, `path_denied`, `payload_guard_denied`, `mtls_denied`).
4. `Quota Denied`: `daily_quota_denied`, `rpm_denied` trend.
5. `Top Impacted Grants`: `grant_id` bazlı deny/success ratio.

Provisioning (dashboard as-code):
1. Dashboard JSON: `/workspace/go/aystek-core/deployments/grafana/dashboards/integrations-shared-service.json`
2. Provider YAML: `/workspace/go/aystek-core/deployments/grafana/provisioning/dashboards/integrations-shared-service.yaml`

## 6. Triage Mapping
1. `execution_failed` yükselirse: downstream endpoint, timeout, auth kontrolü.
2. `scope/method/path denied` yükselirse: grant/service policy drift kontrolü.
3. `secret_binding_denied` veya `mtls_denied`: config rollback + secret/cert doğrulaması.
4. `daily_quota_denied/rpm_denied`: limit planı ve traffic shaping kontrolü.
