# TASK_031_PMS_RETENTION_POLICIES

## Scope
Sprint-11.D kapsamında PMS export/webhook/reconciliation verileri için retention politikası teknik olarak tanımlandı.

## Completed
- Mongo TTL indexleri eklendi:
  - `pms_export_files.created_at`
  - `pms_webhook_deliveries.created_at`
  - `pms_bank_transactions.created_at` (opsiyonel)
- Retention env konfigürasyonları:
  - `PMS_EXPORT_FILES_TTL_DAYS` (default: `60`, `<=0` disable)
  - `PMS_WEBHOOK_DELIVERIES_TTL_DAYS` (default: `30`, `<=0` disable)
  - `PMS_BANK_TXN_TTL_DAYS` (default yok, sadece set edilirse aktif)
- Admin UI tarafında retention bilgilendirme notu eklendi:
  - exports ekranı
  - webhooks ekranı

## Files
- `go/aystek-core/internal/modules/pms/exports/indexes.go`
- `go/aystek-core/internal/modules/pms/webhooks/indexes.go`
- `go/aystek-core/internal/modules/pms/reconciliation/indexes.go`
- `node/aystek-core-next/app/(dashboard)/admin/pms/exports/page.tsx`
- `node/aystek-core-next/app/(dashboard)/admin/pms/webhooks/page.tsx`
- `go/aystek-core/internal/seeds/translation_seed.go`

## Validation
- `go test ./internal/modules/pms/... ./cmd/api`
- `npx eslint app/(dashboard)/admin/pms/exports/page.tsx app/(dashboard)/admin/pms/webhooks/page.tsx`
