# Tenant/Site Authorization Slice 2 Checkpoint

Date: 2026-04-11

## Scope
- login raw-key fallback hardening
- tenant/site selector membership-aware filtering
- user tenant/site membership persistence
- admin user detail membership management UI

## Backend
- Added `context_memberships` to user model.
- Added user endpoints:
  - `GET /admin/users/:id/context-memberships`
  - `PATCH /admin/users/:id/context-memberships`
- Bound `/me/context`, `/me/select-tenant`, `/me/bootstrap-sites` to effective user memberships for non-platform users.

## Frontend
- `login` and `tenant-selector` now avoid rendering raw keys when a translation entry equals the key or preload cache is partial.
- Added merge-based i18n preload behavior so auth/bootstrap translations are not dropped by an already-hydrated partial cache.
- Added tenant/site access management section to admin user detail page.

## Verification
- `cd /workspace/go/aystek-core && go test ./internal/modules/users/... ./internal/modules/system/... ./internal/seeds ./cmd/api`
- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(auth)/login/page.tsx' --file 'app/tenant-selector/page.tsx' --file 'app/(dashboard)/admin/users/[id]/page.tsx' --file 'lib/stores/i18n.ts' --file 'hooks/useI18n.ts' --file 'lib/api/users.ts'`

## Next Recommended Slice
- persist site selection as a first-class user context action if header switching is required globally
- add dedicated admin user membership page/table for bulk management
- optionally separate platform-wide tenant membership admin from tenant-local site access admin
