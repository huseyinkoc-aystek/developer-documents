# Device Command Policy And Active Binding Enforcement Checkpoint

Date: 2026-04-18

## Scope

This slice advanced the stock device command flow from vendor/product validation to true resolved-target enforcement and permission-based policy control.

## Completed

1. Active binding and SIM-state validation was added to resolved target enforcement.
   - `resolved_target_id` is no longer checked only for existence.
   - The backend now resolves active vehicle-asset binding, integration route, product type, MSISDN, and inventory state.
   - Normal execution now fails when the target is:
     - unbound
     - missing MSISDN
     - not in an active inventory state

2. Registry surface moved from read-only toward real administration.
   - Registry summary/read-model remains available.
   - Template delete support was added.
   - Registry management routes now require dedicated registry permission.

3. Privileged override is no longer tied only to platform-user identity.
   - New permissions were introduced:
     - `stock_sms:registry_manage`
     - `stock_sms:route_override`
   - Override execution now depends on permission policy instead of a simple platform-user check.

4. Membership-aware UI visibility was tightened.
   - Launchpad now hides registry and user-management entries when permissions are missing.
   - Dashboard scope cards use membership-management capability instead of broad visibility assumptions.

## Key Files

- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/model/model.go`
- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/repository/interface.go`
- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/repository/mongo.go`
- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/service/service.go`
- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/handler/admin.go`
- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/routes.go`
- `/workspace/go/aystek-core/internal/modules/system/handler/bootstrap.go`
- `/workspace/go/aystek-core/internal/seeds/permission_seed.go`
- `/workspace/go/aystek-core/internal/seeds/translation_seed.go`
- `/workspace/node/aystek-core-next/app/(dashboard)/admin/integrations/device-commands/SmsCommandCenterPage.tsx`
- `/workspace/node/aystek-core-next/app/(dashboard)/admin/integrations/device-commands/registry/page.tsx`
- `/workspace/node/aystek-core-next/app/(dashboard)/page.tsx`
- `/workspace/node/aystek-core-next/components/layout/LaunchpadMenu.tsx`
- `/workspace/node/aystek-core-next/lib/api/me.ts`
- `/workspace/node/aystek-core-next/lib/api/stock-ecosystem.ts`
- `/workspace/node/aystek-core-next/lib/types/stock-ecosystem.ts`

## Verification

```bash
cd /workspace/go/aystek-core && go test ./internal/modules/stock/integrations/device_commands/... ./internal/modules/system/... ./internal/seeds ./cmd/api
cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/integrations/device-commands/SmsCommandCenterPage.tsx' --file 'app/(dashboard)/admin/integrations/device-commands/registry/page.tsx' --file 'app/(dashboard)/page.tsx' --file 'components/layout/LaunchpadMenu.tsx' --file 'lib/api/me.ts' --file 'lib/api/stock-ecosystem.ts' --file 'lib/types/stock-ecosystem.ts'
```

## Next

The next meaningful slice is:

1. Tie `resolved_target_id` to stronger active SIM and asset-state constraints where available.
2. Open full registry CRUD, not only delete/read-model.
3. Harden override policy further so privileged override becomes an auditable governance path instead of a convenience path.
