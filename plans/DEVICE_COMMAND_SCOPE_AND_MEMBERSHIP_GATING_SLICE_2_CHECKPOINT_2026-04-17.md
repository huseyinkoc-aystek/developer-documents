# Device Command Scope And Membership Gating Slice 2 Checkpoint

Date: 2026-04-17

## Completed

- Login auth error rendering now falls back safely even if backend returns raw `auth.*` / `login.*` style codes.
- Login and tenant selector both use the shared modal locale picker with a cleaner pill-style trigger.
- SMS device command center now treats resolved SIM target + resolved vendor route as the primary path.
- Privileged route override is explicit and limited to higher-scope users instead of being the default workflow.
- Device vendor command registry read-model was expanded with:
  - vendor summary cards
  - vendor / intent / product type registry table
- Dashboard admin/module affordance was tightened further:
  - user-management surfaces are no longer shown just because a user can manage sites
  - tenant/user management visibility is more closely aligned to membership and tenant-management scope
- Users membership batch panel gained a faster scoped-admin action:
  - apply the current viewer scope directly to selected users

## Verification

- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(auth)/login/page.tsx' --file 'app/tenant-selector/page.tsx' --file 'components/i18n/LocaleModalPicker.tsx' --file 'app/(dashboard)/page.tsx' --file 'app/(dashboard)/admin/users/page.tsx' --file 'app/(dashboard)/admin/integrations/device-commands/SmsCommandCenterPage.tsx'`
- `cd /workspace/go/aystek-core && go test ./internal/modules/stock/integrations/device_commands/... ./internal/seeds ./cmd/api`

## Next

- Push device command center from “route strongly required” to “resolved route required except explicit privileged override policy”.
- Add dedicated admin CRUD/read-model surfaces for the device vendor command registry instead of relying only on template-derived summaries.
- Bind launchpad and more dashboard module entry points to membership scope with the same strictness used in context switching.
