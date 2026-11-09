# Device Command Resolved Target Enforcement

Date: 2026-04-18

## Completed

- Device command send flow now requires a resolved stock target identity in normal operation.
- `resolved_target_id` is now part of the request contract and execution record.
- Backend validates that the resolved target exists inside current tenant/site scope.
- Non-override requests now require all three:
  - `resolved_target_id`
  - `vendor_key`
  - `product_type`
- Override remains privileged-only and bypasses resolved target requirement only for platform-level flows.

## Updated Files

- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/model/model.go`
- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/dto/dto.go`
- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/repository/interface.go`
- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/repository/mongo.go`
- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/service/service.go`
- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/handler/admin.go`
- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/service/service_test.go`
- `/workspace/node/aystek-core-next/app/(dashboard)/admin/integrations/device-commands/SmsCommandCenterPage.tsx`
- `/workspace/node/aystek-core-next/lib/api/stock-ecosystem.ts`
- `/workspace/node/aystek-core-next/lib/types/stock-ecosystem.ts`

## Validation

- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/integrations/device-commands/SmsCommandCenterPage.tsx' --file 'app/(dashboard)/admin/integrations/device-commands/registry/page.tsx' --file 'lib/api/stock-ecosystem.ts' --file 'lib/types/stock-ecosystem.ts'`
- `cd /workspace/go/aystek-core && go test ./internal/modules/stock/integrations/device_commands/... ./internal/seeds ./cmd/api`
