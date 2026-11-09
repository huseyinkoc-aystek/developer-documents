# Tenant Site Authorization Slice 1 Foundation

Date: 2026-04-11

## Goal

Start the enterprise-grade authorization model for:

- user -> tenant memberships
- user -> site memberships
- controlled context switching
- admin-managed access boundaries for tenant/site selection flows

## Why This Slice Exists

Current runtime now supports better tenant/site selection and bootstrap flows, but it still does **not** enforce a complete per-user tenant/site membership model. That means selector UX can be correct while authorization remains too coarse.

## Slice 1 Deliverables

1. Permission foundation
- `users:tenant_membership_read`
- `users:tenant_membership_manage`
- `users:site_membership_read`
- `users:site_membership_manage`
- `users:context_switch_manage`

2. Domain planning boundary
- membership records should be explicit, not inferred from raw tenant IDs on user objects
- context switching should read membership scope, not just currently selected tenant/site

3. Next slice target
- add membership persistence model
- add admin CRUD/read endpoints
- bind tenant selector and header context switcher to effective permissions

## Status

- Slice 1 foundation started
- permission seed foundation added
- full membership model still pending
