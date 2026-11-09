# Device Command Registry And Membership Gating Slice 3

Date: 2026-04-18

## Completed

- Backend override policy for device command send is now enforced.
  - `allow_override` is accepted by API.
  - only platform users can activate override.
  - non-override requests now require resolved route metadata (`vendor_key` and `product_type`).

- Device vendor command registry read-model endpoint was added.
  - `GET /admin/stock/integrations/device-commands/registry`
  - grouped by `vendor_key + command_intent + product_type`
  - returns template, active, inactive counts and last update timestamp.

- Separate admin/read-model surface was added in Next.
  - `/admin/integrations/device-commands/registry`
  - includes summary cards, registry table and template library links.

- Command center was tightened.
  - override flag is sent explicitly to backend.
  - registry summary now comes from backend read-model instead of only local client aggregation.

- Launchpad membership filtering was hardened.
  - scoped users no longer see selected out-of-scope entries such as:
    - `/admin/users`
    - `/admin/settings/tenant-site`
    - `/admin/system`

- Dashboard scope gating was tightened.
  - scoped admins no longer see the same global user/content stat block as platform-wide users.
  - scoped dashboard now emphasizes visible tenant/site scope and site-management capability.

## Validation

- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/integrations/device-commands/SmsCommandCenterPage.tsx' --file 'app/(dashboard)/admin/integrations/device-commands/registry/page.tsx' --file 'components/layout/LaunchpadMenu.tsx' --file 'app/(dashboard)/page.tsx' --file 'lib/api/stock-ecosystem.ts' --file 'lib/types/stock-ecosystem.ts'`
- `cd /workspace/go/aystek-core && go test ./internal/modules/stock/integrations/device_commands/... ./internal/seeds ./cmd/api`

## Next

- make resolved target identity stricter than vendor/product metadata alone
- add registry CRUD affordances beyond summary/read-model
- turn override policy into a dedicated backend policy surface instead of implicit platform-user check
