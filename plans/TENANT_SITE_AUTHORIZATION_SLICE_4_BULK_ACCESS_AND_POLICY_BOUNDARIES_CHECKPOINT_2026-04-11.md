# Tenant Site Authorization Slice 4 Checkpoint

Date: 2026-04-11

## Scope Closed

This slice completed the real management layer for tenant and site access governance.

- Bulk tenant/site access update endpoint added for admin users
- Policy boundaries tightened between platform admins and scoped tenant admins
- Users list UI extended with bulk tenant/site access assignment
- Translation seed coverage completed for new bulk-access messaging

## Backend Changes

- Added `PATCH /admin/users/context-memberships/bulk`
- Added bulk DTO and service contract for context membership updates
- Enforced policy boundaries for:
  - read access to target users
  - update access to target users
  - membership assignments outside actor scope
  - platform-managed users vs scoped tenant admins

Files:

- `/workspace/go/aystek-core/internal/modules/users/dto/dto.go`
- `/workspace/go/aystek-core/internal/modules/users/service/interface.go`
- `/workspace/go/aystek-core/internal/modules/users/service/service.go`
- `/workspace/go/aystek-core/internal/modules/users/handler/admin.go`
- `/workspace/go/aystek-core/internal/modules/users/routes.go`
- `/workspace/go/aystek-core/internal/seeds/translation_seed.go`

## Frontend Changes

- Users list page now supports selecting multiple users and applying the same tenant/site access set
- Tenant/site access summary is visible on each user card
- Bulk panel loads available tenants and scoped sites from the authenticated context

Files:

- `/workspace/node/aystek-core-next/app/(dashboard)/admin/users/page.tsx`
- `/workspace/node/aystek-core-next/lib/api/users.ts`

## Verification

Passed:

```bash
cd /workspace/go/aystek-core && go test ./internal/modules/users/... ./internal/modules/system/... ./internal/seeds ./cmd/api
cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/users/page.tsx' --file 'lib/api/users.ts' --file 'components/context/TenantSwitcher.tsx' --file 'components/context/SiteSwitcher.tsx' --file 'app/(auth)/login/page.tsx' --file 'app/tenant-selector/page.tsx' --file 'components/i18n/LocalePillSwitcher.tsx' --file 'lib/stores/language.ts' --file 'lib/tenant.ts' --file 'middleware.ts'
```

## Current Status

Closed:

- bulk tenant/site access management
- scoped policy boundary enforcement
- users list management surface

Next:

- tenant/site memberships for non-platform admins should be surfaced in a more operational matrix view
- header context switch can be hardened further against out-of-scope selections by reading these memberships directly from a normalized authorization payload
