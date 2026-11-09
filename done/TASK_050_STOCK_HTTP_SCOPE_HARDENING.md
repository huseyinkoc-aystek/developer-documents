# TASK_050 — STOCK HTTP Scope Hardening

## Kapsam
- Stock admin handler seviyesinde tenant/site scope doğrulaması sertleştirildi.
- HTTP-level negatif testlerle erişim sınırları doğrulandı.

## Uygulama
- `internal/modules/stock/handler/admin.go`
  - `tenant_id` fallback (`default`) kaldırıldı.
  - `requireTenantID` ve `requireSiteID` guard’ları eklendi.
  - Etkilenen endpointler:
    - `ListVariantStock`
    - `BulkUpdateVariantStockAsync`
    - `QueueExportStockCSV`
    - `ListJobs`

- Portal sales object-scope testleri:
  - `internal/modules/stock/sales/handler/portal_test.go`
  - Kullanıcı bağlamı olmayan çağrıda `401`
  - `PortalList` çağrısında `assignee_user_id` zorunlu user scope
  - `PortalCreate` çağrısında `assignee_user_id/created_by` context user’dan set edilir

- Admin scope testleri:
  - `internal/modules/stock/handler/admin_scope_test.go`
  - `tenant_id` olmadan export job queue çağrısı `401`
  - `site_id` olmadan export job queue çağrısı `400`
  - `tenant_id` olmadan jobs list çağrısı `401`

## Doğrulama
- `go test ./internal/modules/stock/handler/... ./internal/modules/stock/sales/handler/... ./internal/modules/stock/... ./cmd/api` ✅

