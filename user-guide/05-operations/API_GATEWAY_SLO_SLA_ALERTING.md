---
Title: API Gateway SLO/SLA & Alerting Baseline
Version: 1.0.0
Status: Active
Last Updated: 2026-03-09
Related Documents:
  - ./OPERATIONS_GUIDE.md
  - ./API_GATEWAY_SECURITY_AUDIT_CHECKLIST.md
  - ../../plans/AYSTEK_CORE_API_GATEWAY_THIRD_PARTY_ACCESS_CONTROL.md
---

# API Gateway SLO/SLA & Alerting Baseline

Bu dokuman, API Gateway third-party access katmani icin SLO/SLA hedeflerini ve Prometheus/Grafana izleme baseline'ini tanimlar.

## 1. SLO Targets
1. Availability: aylik %99.9 basarili cevap orani (5xx haric).
2. Latency: p95 `request_duration_seconds` < 500ms.
3. Auth Failure Guardrail: auth failure ratio (15m) < %5.
4. Rate-limit Fairness: tenant/client bazli deny burst anomalileri 10 dk icinde tespit.

## 2. SLA Statement (Partner Facing)
1. Planned maintenance disinda API availability hedefi %99.9.
2. Critical incidentlerde ilk yanit suresi: 30 dakika.
3. Production incidents icin RCA teslim suresi: 3 is gunu.
4. Secret compromise suphelerinde acil rotation proseduru aninda uygulanir.

## 3. Metric Surface
1. `aystek_api_gateway_requests_total{tenant_id,api_user_id,method,route,status}`
2. `aystek_api_gateway_request_duration_seconds_bucket{tenant_id,api_user_id,method,route,le}`
3. `aystek_api_gateway_auth_failures_total{reason}`
4. `aystek_api_gateway_rate_limit_exceeded_total{tenant_id,api_user_id,reason}`
5. `aystek_api_gateway_route_denied_total{tenant_id,api_user_id,reason,method,route}`

## 4. Prometheus Ruleset
1. Rules file: `go/aystek-core/deployments/prometheus/alerts.yml`
2. Group name: `aystek-core_api_gateway`
3. Baseline alerts:
   - `APIGatewayHighErrorRate`
   - `APIGatewayHighLatencyP95`
   - `APIGatewayAuthFailureSpike`
   - `APIGatewayRouteDeniedBurst`
   - `APIGatewayRateLimitBurst`

## 5. Grafana Dashboard
1. Dashboard JSON:
   - `go/aystek-core/deployments/grafana/dashboards/api-gateway-third-party.json`
2. Provisioning YAML:
   - `go/aystek-core/deployments/grafana/provisioning/dashboards/api-gateway-third-party.yaml`
3. Panel baseline:
   - Request rate by route/status
   - Error ratio
   - p95 latency
   - Auth failure trend
   - Rate-limit/route-deny trend

## 6. Operational Review
1. Haftalik: SLO trend review + noisy alert tuning.
2. Aylik: SLA raporu + partner etkisi degerlendirmesi.
3. Incident sonrasi: alert threshold kalibrasyonu ve runbook update.
