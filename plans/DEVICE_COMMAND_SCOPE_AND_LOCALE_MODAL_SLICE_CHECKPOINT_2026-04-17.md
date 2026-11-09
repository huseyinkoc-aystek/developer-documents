# Device Command Scope And Locale Modal Slice Checkpoint

Date: 2026-04-17

## Completed

- Login error rendering no longer leaks raw auth keys when backend returns `auth.*` / `login.*` style codes.
- Login language selector and tenant-selector language selector now use the shared modal picker with a cleaner pill trigger.
- SMS device command center now pushes operators toward:
  - resolved SIM target
  - resolved or explicitly chosen vendor route
  - approved template / command intent
- Manual free-form partner routing was reduced; resolved vendor route is now surfaced as the primary source of truth.
- Device command registry gained a stronger read-model summary block grouped by vendor / intent / product type.
- Dashboard global visibility tightened:
  - global metrics and analytics stay more visible for platform/global admins
  - scoped admins see a narrower management surface
- Users bulk membership panel now exposes clearer platform-vs-scoped affordance and a faster “apply visible scope” action.

## Verification

- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(auth)/login/page.tsx' --file 'app/tenant-selector/page.tsx' --file 'components/i18n/LocaleModalPicker.tsx' --file 'app/(dashboard)/page.tsx' --file 'app/(dashboard)/admin/users/page.tsx' --file 'app/(dashboard)/admin/integrations/device-commands/SmsCommandCenterPage.tsx'`
- `cd /workspace/go/aystek-core && go test ./internal/modules/stock/integrations/device_commands/... ./internal/seeds ./cmd/api`

## Next

- Bind dashboard module cards and launch surfaces to membership scope even more strictly.
- Add a denser device vendor command registry admin/read-model surface.
- Move from “route strongly encouraged” toward “resolved route required except explicit privileged override” in the device command center.
