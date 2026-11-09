# TASK_049 — STOCK Security Negative Tests Expansion

## Kapsam
- `stock/security` negatif test paketi genişletildi.
- Tenant/site context eksikliği için sales/returns/sms callback security akışları kapsandı.

## Eklenen Test Senaryoları
- `stock_sales` create çağrısında tenant/site context yoksa reject
- `stock_returns` create çağrısında tenant/site context yoksa reject
- `stock_sms` callback authorize çağrısında tenant/site context yoksa reject

## Dosya
- `go/aystek-core/internal/modules/stock/security/security_negative_test.go`

## Doğrulama
- `go test ./internal/modules/stock/security/... ./internal/modules/stock/... ./cmd/api` ✅

