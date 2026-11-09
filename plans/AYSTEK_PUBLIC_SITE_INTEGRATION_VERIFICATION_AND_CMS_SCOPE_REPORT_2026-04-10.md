# AYSTEK Public Site Integration Verification And CMS Scope Report

Date: 2026-04-10

## Verified In This Slice

- `default` tenant seed already contains the admin panel site:
  - `code=admin-panel`
  - `domain=core.aysteknolojileri.com`
- `aystek` tenant seed already contains the corporate public site:
  - `code=aystek`
  - `domain=aysteknolojileri.com`
- `aysteknoloji-com` public Next app is wired to the intended backend context:
  - `X-Tenant-ID=aystek`
  - `X-Store-Domain=aysteknolojileri.com`
- AYSTEK public CMS seed is present and active:
  - media
  - taxonomy
  - sliders
  - pages
  - posts
  - forum

## Hardening Added In Code

- Localhost admin resolution now prefers `admin-panel` before `default`.
- Bootstrap tenant/site routes are now context-bootstrap safe.
- Tenant/site bootstrap mutations no longer depend on `/auth/csrf` prefetch.
- Store domain resolution now supports `metadata.domain_aliases` in addition to the primary `domain` field.

Files:

- `/workspace/go/aystek-core/internal/platform/middleware/store_context.go`
- `/workspace/go/aystek-core/internal/platform/middleware/bootstrap_guard.go`
- `/workspace/go/aystek-core/internal/platform/middleware/site_isolation.go`
- `/workspace/go/aystek-core/internal/modules/stores/repository/mongo.go`
- `/workspace/go/aystek-core/internal/modules/stores/repository/lookup.go`
- `/workspace/node/aystek-core-next/lib/api/requestContext.ts`
- `/workspace/node/aystek-core-next/lib/api/me.ts`
- `/workspace/node/aystek-core-next/app/tenant-selector/page.tsx`
- `/workspace/node/aystek-core-next/app/(auth)/login/page.tsx`

## CMS Scope Reality Check

The current public CMS domain model is not equally site-scoped across all modules.

### Fully tenant+site scoped

- `posts`
- post geo / relations / placements / media mappings

### Tenant-scoped, not fully site-scoped

- `pages`
- `menus`
- `sliders`
- `media`

## Practical Meaning

For the current `aystek` corporate setup, this is acceptable enough to serve one primary public site correctly because that tenant has one main public site contract.

However, for true enterprise-grade **multi-site public CMS under one tenant**, the following residual risk remains:

- page sets are shared at tenant level
- nav menus are shared at tenant level
- sliders are shared at tenant level
- media library is shared at tenant level

This means a second public site under the same tenant could not be cleanly isolated yet.

## Recommended Next Slice

`Public CMS Site Scope Hardening`

Target:

- add `site_id` to pages
- add `site_id` to menus
- add `site_id` to sliders
- add `site_id` to media
- migrate existing `aystek` content to its resolved site id
- update public/admin repositories and handlers to enforce tenant+site filtering consistently

## Validation Performed

- `cd /workspace/go/aystek-core && go test ./internal/platform/middleware ./internal/modules/stores/repository ./internal/seeds ./cmd/api`

## Validation Limitation

- `/workspace/node/aysteknoloji-com` lint could not be executed non-interactively because the repo currently opens the initial Next ESLint setup prompt instead of running a committed lint configuration.
