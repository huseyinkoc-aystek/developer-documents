# Device Command Vendor Route Policy Slice Checkpoint

Date: 2026-04-19

## Scope
- Template-level registry summary was extended into real vendor-route policy records.
- Policy CRUD surface was opened in Next admin.
- Override path was further narrowed through policy requirements.
- Launchpad and dashboard affordances were tightened with membership + command permission scope.

## Backend
- Added vendor-route policy read/write model under:
  - `internal/modules/stock/integrations/device_commands/model`
  - `internal/modules/stock/integrations/device_commands/repository`
  - `internal/modules/stock/integrations/device_commands/service`
  - `internal/modules/stock/integrations/device_commands/handler`
- Added policy CRUD routes:
  - `GET /policies`
  - `POST /policies`
  - `GET /policies/:id`
  - `PUT /policies/:id`
  - `DELETE /policies/:id`
- Added `stock_sms:policy_manage` permission.
- Strengthened `SendCommand` policy behavior:
  - normal flow respects `ExecutionEnabled`
  - override flow now requires a matching policy with `AllowOverride`
  - override remains gated by existing permission flow

## Frontend
- Added policy CRUD pages:
  - `app/(dashboard)/admin/integrations/device-commands/policies/new/page.tsx`
  - `app/(dashboard)/admin/integrations/device-commands/policies/[id]/page.tsx`
- Extended registry page to show vendor-route policies and allow create/edit/delete.
- Added direct policy affordances in command center.
- Tightened launchpad gating for policy routes.
- Added dashboard scope cards for device command center and vendor-route policies.

## i18n
- Added policy-related translation keys for:
  - titles
  - subtitles
  - empty state
  - save notice
  - execution/override labels
  - dashboard scope cards

## Verification
- `cd /workspace/go/aystek-core && go test ./internal/modules/stock/integrations/device_commands/... ./internal/seeds ./cmd/api`
- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/integrations/device-commands/SmsCommandCenterPage.tsx' --file 'app/(dashboard)/admin/integrations/device-commands/registry/page.tsx' --file 'app/(dashboard)/admin/integrations/device-commands/policies/PolicyForm.tsx' --file 'app/(dashboard)/admin/integrations/device-commands/policies/new/page.tsx' --file 'app/(dashboard)/admin/integrations/device-commands/policies/[id]/page.tsx' --file 'components/layout/LaunchpadMenu.tsx' --file 'app/(dashboard)/page.tsx' --file 'lib/api/stock-ecosystem.ts' --file 'lib/types/stock-ecosystem.ts'`

## Next Slice
- Move from policy records into richer vendor-route governance.
- Add audit-focused CRUD actions around policy changes.
- Tighten dashboard/admin visibility with combined membership + device command permissions.
