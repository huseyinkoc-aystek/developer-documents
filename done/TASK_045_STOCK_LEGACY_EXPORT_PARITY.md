# TASK_045 — STOCK Legacy Export Parity Hardening

## Kapsam
- Legacy `stock/export` akışı Sprint-7 parity hardening seviyesine getirildi.
- CSV güvenliği ve async export runner, legacy endpoint yüzeyine uyarlandı.

## Backend
- Yeni async endpointler:
  - `POST /admin/stock/export/jobs`
  - `GET /admin/stock/export/jobs/:id/download`
- Mevcut endpoint:
  - `GET /admin/stock/export` için CSV formula sanitize eklendi.
- Job sistemi:
  - Yeni job tipi: `stock.export_csv`
  - Queue factory hydration:
    - `stock.bulk_import`
    - `stock.export_csv`
  - `job_status` kaydı genişletildi (`site_id`, `file_name`, `file_data`).

## Next.js
- `stock reports exports` ekranına legacy export job akışı eklendi:
  - warehouse_id girerek legacy export job oluşturma
  - legacy job listesi (export job tipine filtreli)
  - completed job için download aksiyonu
  - pending/processing polling

## i18n
- `stock.reports.*` altında legacy export için ek anahtarlar eklendi:
  - `stock.reports.fields.job_id`
  - `stock.reports.fields.warehouse_id`
  - `stock.reports.legacy_status.*`
  - `stock.reports.exports.legacy_title`
  - `stock.reports.exports.legacy_create`
  - `stock.reports.validation.warehouse_required`

## Doğrulama
- `go test ./internal/modules/stock/... ./cmd/api` ✅
- `npx eslint app/(dashboard)/admin/stock/reports/exports/jobs/page.tsx lib/api/stock-ecosystem.ts lib/types/stock-ecosystem.ts` ✅
