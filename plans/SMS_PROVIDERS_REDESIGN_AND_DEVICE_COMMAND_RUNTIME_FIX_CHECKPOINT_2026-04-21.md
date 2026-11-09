# SMS Providers Redesign and Device Command Runtime Fix Checkpoint

Date: 2026-04-21

## Scope Closed

- Fixed the `resolvedVendorKey` runtime crash on the device commands center.
- Redesigned the legacy `/admin/sms/providers` page to better match the AYSTEK Core admin design language.
- Split the legacy SMS providers surface into clearer sections instead of mixing config form and operational lists in one indistinct block.
- Added and aligned `sms.*` translation coverage so raw keys are less likely to leak on the legacy SMS surface.
- Continued governance work already in progress with pending override queue visibility, actor metadata, and approval read-model support.

## Key Technical Changes

### Next.js

- `app/(dashboard)/admin/integrations/device-commands/SmsCommandCenterPage.tsx`
  - Fixed `resolvedVendorKey` initialization order so the page no longer crashes at runtime.
- `app/(dashboard)/admin/sms/providers/page.tsx`
  - Reworked the page into:
    - hero/header section
    - sub-navigation pills
    - provider summary cards
    - provider configuration card
    - test send card
    - provider registry/list section
  - Added locale-safe `t(key, fallback)` usage across visible labels and actions.

### Go

- `internal/seeds/translation_seed.go`
  - Added/expanded `sms.*` translation keys needed by the redesigned SMS providers UI.

## Validation

### Next

```bash
cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/integrations/device-commands/SmsCommandCenterPage.tsx' --file 'app/(dashboard)/admin/integrations/device-commands/registry/page.tsx' --file 'app/(dashboard)/admin/sms/providers/page.tsx' --file 'components/layout/LaunchpadMenu.tsx' --file 'lib/api/stock-ecosystem.ts' --file 'lib/types/stock-ecosystem.ts'
```

Result: passed

### Go

```bash
cd /workspace/go/aystek-core && go test ./internal/modules/stock/integrations/device_commands/... ./internal/seeds ./cmd/api
```

Result: passed

## Outcome

- The device commands center no longer fails with the reported `ReferenceError`.
- The legacy SMS providers page is closer to the established admin visual system and reads as a structured management surface instead of a mixed one-page utility form.
- The current governance queue and approval metadata work remains intact and validated after the SMS page changes.
