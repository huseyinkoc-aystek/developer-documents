# Device Command Registry CRUD And Override Governance Checkpoint

Date: 2026-04-19

## Scope

This slice moved the device command registry and privileged override flow one step closer to enterprise-grade governance.

## Completed

1. Privileged override now requires an explicit audit reason.
   - `override_reason` was added to the send-command contract.
   - Backend rejects override requests without a reason.
   - Execution records now store:
     - `override_used`
     - `override_reason`

2. Override was connected to the audit trail.
   - A dedicated activity log event is now emitted when a command is sent with privileged override.
   - This event includes resolved target, route, product type, and the provided reason.

3. Registry UI became more operational.
   - Registry page now exposes clearer create/edit/delete actions.
   - Command center now shows override reason controls in the actual send flow.
   - Execution stream now exposes override usage and reason data.

4. Launchpad gating was tightened for the device commands area.
   - Command center requires `stock_sms:read` or `stock_sms:execute`.
   - Registry and template management require `stock_sms:registry_manage`.

## Key Files

- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/model/model.go`
- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/dto/dto.go`
- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/service/service.go`
- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/handler/admin.go`
- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/service/service_test.go`
- `/workspace/go/aystek-core/internal/seeds/translation_seed.go`
- `/workspace/node/aystek-core-next/app/(dashboard)/admin/integrations/device-commands/SmsCommandCenterPage.tsx`
- `/workspace/node/aystek-core-next/app/(dashboard)/admin/integrations/device-commands/registry/page.tsx`
- `/workspace/node/aystek-core-next/components/layout/LaunchpadMenu.tsx`
- `/workspace/node/aystek-core-next/lib/api/stock-ecosystem.ts`
- `/workspace/node/aystek-core-next/lib/types/stock-ecosystem.ts`

## Verification

```bash
cd /workspace/go/aystek-core && go test ./internal/modules/stock/integrations/device_commands/... ./internal/seeds ./cmd/api
cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/integrations/device-commands/SmsCommandCenterPage.tsx' --file 'app/(dashboard)/admin/integrations/device-commands/registry/page.tsx' --file 'components/layout/LaunchpadMenu.tsx' --file 'lib/api/stock-ecosystem.ts' --file 'lib/types/stock-ecosystem.ts'
```

## Next

The next meaningful slice is:

1. Open full registry CRUD beyond create/edit/delete into policy-managed vendor-route records.
2. Add stronger route-resolution governance so override becomes a rare exception flow.
3. Bind dashboard and admin affordance visibility even more tightly to membership and command permissions.
