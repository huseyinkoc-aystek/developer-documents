# Bulk SMS System Presets And Approver Read Model Checkpoint

Date: 2026-04-22

## Scope

This slice moved bulk SMS presets from UI-only defaults toward backend-managed system presets and improved approval history readability with approver-aware read-model data.

## Completed

- Added backend-managed system bulk segment presets.
- Added persisted preset CRUD endpoints for tenant/site scope.
- Added approval history endpoint with approver name enrichment.
- Updated bulk SMS UI to use backend-driven presets instead of only hardcoded UI options.
- Updated approval history UI to prefer readable approver names over raw actor IDs.
- Added `sms:bulk_manage` permission seed usage for preset management.

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
- `/workspace/node/aystek-core-next/app/(dashboard)/admin/integrations/bulk-sms/page.tsx`
- `/workspace/node/aystek-core-next/lib/api/sms.ts`
- `/workspace/node/aystek-core-next/lib/types/sms.ts`

## Verification

- `cd /workspace/go/aystek-core && go test ./internal/modules/sms/... ./internal/modules/users/... ./internal/seeds ./cmd/api`
- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/integrations/bulk-sms/page.tsx' --file 'lib/api/sms.ts' --file 'lib/types/sms.ts'`

## Next Logical Step

- Add editable lifecycle rules for system presets instead of only auto-seeding missing defaults.
- Enrich approval history with approver role and template/provider labels.
- Add trend reporting for quota and rate consumption over time.
