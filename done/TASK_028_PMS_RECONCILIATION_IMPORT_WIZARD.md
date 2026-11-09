# TASK_028_PMS_RECONCILIATION_IMPORT_WIZARD

## Scope
Sprint-11.A kapsamında PMS reconciliation için operasyon odaklı CSV import wizard eklendi.

## Completed
- Admin UI: `/admin/pms/reconciliation/import`
  - CSV paste
  - delimiter seçimi (`,`, `;`, `TAB`)
  - kolon mapping
  - ilk 20 satır preview
  - satır bazlı validasyon hataları
  - 200’lük chunk batch ingest
- Backend ingest hattı batch payload (`items[]`) desteğiyle genişletildi.
- Batch sonuç modeli eklendi:
  - `created_count`, `duplicate_count`, `failed_count`, `errors[]`
- Audit aktiviteleri eklendi:
  - `pms.reconciliation.import.started`
  - `pms.reconciliation.import.completed`

## Files
- `go/aystek-core/internal/modules/pms/reconciliation/dto/dto.go`
- `go/aystek-core/internal/modules/pms/reconciliation/handler/admin.go`
- `go/aystek-core/internal/modules/pms/reconciliation/service/interface.go`
- `go/aystek-core/internal/modules/pms/reconciliation/service/service.go`
- `go/aystek-core/internal/modules/pms/reconciliation/service/service_test.go`
- `node/aystek-core-next/app/(dashboard)/admin/pms/reconciliation/import/page.tsx`
- `node/aystek-core-next/app/(dashboard)/admin/pms/reconciliation/page.tsx`
- `node/aystek-core-next/lib/api/pms.ts`
- `node/aystek-core-next/lib/types/pms.ts`
- `go/aystek-core/internal/seeds/translation_seed.go`

## Validation
- `go test ./internal/modules/pms/... ./cmd/api`
- `npx eslint app/(dashboard)/admin/pms/reconciliation/import/page.tsx lib/api/pms.ts lib/types/pms.ts`
