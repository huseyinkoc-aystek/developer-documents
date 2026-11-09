# Bulk SMS Persisted Segments And Approval History Checkpoint

Date: 2026-04-22

## Scope

This slice moved bulk SMS operations beyond UI-only presets and added persisted reusable audience segments plus a dedicated approval history read-model.

## Completed

- Added persisted bulk SMS segment presets.
- Added CRUD-style endpoints for segment presets.
- Added approval history endpoint backed by approval-required campaigns.
- Expanded the bulk SMS page with:
  - reusable saved segments
  - current-filter save flow
  - approval history section
  - richer campaign/governance visibility
- Added `sms:bulk_manage` permission seed.
- Added missing translation keys for segment management and approval history UI.

## Key Files

- `/workspace/go/aystek-core/internal/modules/sms/model/model.go`
- `/workspace/go/aystek-core/internal/modules/sms/dto/dto.go`
- `/workspace/go/aystek-core/internal/modules/sms/repository/interface.go`
- `/workspace/go/aystek-core/internal/modules/sms/repository/mongo.go`
- `/workspace/go/aystek-core/internal/modules/sms/service/service.go`
- `/workspace/go/aystek-core/internal/modules/sms/handler/admin.go`
- `/workspace/go/aystek-core/internal/modules/sms/routes.go`
- `/workspace/go/aystek-core/internal/modules/sms/service/service_test.go`
- `/workspace/go/aystek-core/internal/seeds/permission_seed.go`
- `/workspace/go/aystek-core/internal/seeds/translation_seed.go`
- `/workspace/node/aystek-core-next/app/(dashboard)/admin/integrations/bulk-sms/page.tsx`
- `/workspace/node/aystek-core-next/lib/api/sms.ts`
- `/workspace/node/aystek-core-next/lib/types/sms.ts`

## Verification

- `cd /workspace/go/aystek-core && go test ./internal/modules/sms/... ./internal/modules/users/... ./internal/seeds ./cmd/api`
- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/integrations/bulk-sms/page.tsx' --file 'lib/api/sms.ts' --file 'lib/types/sms.ts'`

## Next Logical Step

- Add editable system-level segment presets instead of code-defined defaults only.
- Add richer approval history filtering and approver read-model enrichment.
- Add message-template-driven segment defaults and governance trend reporting.
