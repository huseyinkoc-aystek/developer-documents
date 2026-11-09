# Device Command Read-Model Labels and Integrations Scope Cleanup

Date: 2026-04-21

## Scope Closed

- Enriched device command registry/read-model displays with shared vendor, command intent, and stock product type labels.
- Removed remaining legacy `/admin/sms/*` duplication pressure by filtering stale launchpad entries and consolidating navigation around `Integrations > Device Commands`.
- Replaced manual `product_type` entry on device command template and policy forms with stock-backed active product type selections.

## Key Decisions

- `product_type` should be chosen from stock inventory settings, not free-text.
- Vendor, command intent, and product type labels should be rendered through one shared read-model helper instead of repeated inline formatting.
- Until legacy menu rows are fully reseeded everywhere, launchpad should hide `/admin/sms/*` items on the client side to avoid duplicate operational surfaces.

## Files Updated

### Next.js

- `lib/device-commands/read-model.ts`
  - new shared label helpers for:
    - vendor
    - command intent
    - product type
- `app/(dashboard)/admin/integrations/device-commands/SmsCommandCenterPage.tsx`
  - registry and execution surface now use shared label helpers
  - command intent select now shows readable labels
  - product type display resolves against stock product type labels
- `app/(dashboard)/admin/integrations/device-commands/registry/page.tsx`
  - pending queue, policy audit, policy listing, registry rows, and template detail sections now use shared labels
- `app/(dashboard)/admin/integrations/device-commands/templates/TemplateForm.tsx`
  - active stock product types loaded from inventory settings
  - `product_type` is now a select
- `app/(dashboard)/admin/integrations/device-commands/policies/PolicyForm.tsx`
  - active stock product types loaded from inventory settings
  - `product_type` is now a select
- `components/layout/LaunchpadMenu.tsx`
  - stale legacy `/admin/sms/*` items are filtered out to avoid duplicate navigation

## Validation

### Next

```bash
cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/integrations/device-commands/SmsCommandCenterPage.tsx' --file 'app/(dashboard)/admin/integrations/device-commands/registry/page.tsx' --file 'app/(dashboard)/admin/integrations/device-commands/templates/TemplateForm.tsx' --file 'app/(dashboard)/admin/integrations/device-commands/policies/PolicyForm.tsx' --file 'components/layout/LaunchpadMenu.tsx' --file 'lib/device-commands/read-model.ts'
```

Result: passed

### Go

```bash
cd /workspace/go/aystek-core && go test ./internal/modules/stock/integrations/device_commands/... ./internal/seeds ./cmd/api
```

Result: passed

## Outcome

- Device command governance surfaces now read more like a coherent domain view instead of raw code/value dumps.
- Product type handling is aligned to stock truth.
- Integrations navigation is cleaner and less likely to show duplicated legacy SMS paths.
