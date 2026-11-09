# Device Command Product Type Catalog and Legacy SMS Redirect Cleanup

Date: 2026-04-21

## Scope Closed

- Replaced free-text `product_type` entry points with stock inventory product type selections on device command template and policy forms.
- Simplified legacy `/admin/sms/*` navigation by redirecting it into the `Integrations > Device Commands` flow.
- Removed remaining legacy SMS links from the device command center so the admin surface no longer points users back into duplicated SMS screens.
- Repointed seeded menu entries so SMS navigation resolves into the integrations/device-commands route family.

## Key Decisions

- `product_type` should not be typed manually for command registry/policy forms.
- The authoritative source is stock inventory settings:
  - active product types from `getInventorySettings().product_types`
- Legacy `/admin/sms/*` screens were creating duplication and visual inconsistency.
- For now, the correct cleanup is redirect-based consolidation instead of keeping two parallel admin surfaces alive.

## Files Updated

### Next.js

- `app/(dashboard)/admin/integrations/device-commands/templates/TemplateForm.tsx`
  - loads active stock product types
  - uses a select instead of free-text input
- `app/(dashboard)/admin/integrations/device-commands/policies/PolicyForm.tsx`
  - loads active stock product types
  - uses a select instead of free-text input
- `app/(dashboard)/admin/integrations/device-commands/SmsCommandCenterPage.tsx`
  - removed remaining legacy `/admin/sms/*` shortcuts
  - replaced them with integrations-native navigation
- legacy routes now redirect:
  - `app/(dashboard)/admin/sms/providers/page.tsx`
  - `app/(dashboard)/admin/sms/credentials/page.tsx`
  - `app/(dashboard)/admin/sms/messages/page.tsx`
  - `app/(dashboard)/admin/sms/outbox/page.tsx`
  - `app/(dashboard)/admin/sms/policies/page.tsx`
  - `app/(dashboard)/admin/sms/templates/page.tsx`

### Go

- `internal/seeds/menu_seed.go`
  - SMS-related menu paths now point to the `integrations/device-commands` route family

## Validation

### Next

```bash
cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/integrations/device-commands/SmsCommandCenterPage.tsx' --file 'app/(dashboard)/admin/integrations/device-commands/templates/TemplateForm.tsx' --file 'app/(dashboard)/admin/integrations/device-commands/policies/PolicyForm.tsx' --file 'app/(dashboard)/admin/sms/providers/page.tsx' --file 'app/(dashboard)/admin/sms/credentials/page.tsx' --file 'app/(dashboard)/admin/sms/messages/page.tsx' --file 'app/(dashboard)/admin/sms/outbox/page.tsx' --file 'app/(dashboard)/admin/sms/policies/page.tsx' --file 'app/(dashboard)/admin/sms/templates/page.tsx'
```

Result: passed

### Go

```bash
cd /workspace/go/aystek-core && go test ./internal/modules/stock/integrations/device_commands/... ./internal/seeds ./cmd/api
```

Result: passed

## Outcome

- Device command governance surfaces now rely on stock-defined product types instead of manual text entry.
- The duplicated SMS admin route family is no longer the primary operational surface.
- SMS navigation is now more consistent with the integration-led architecture direction.
