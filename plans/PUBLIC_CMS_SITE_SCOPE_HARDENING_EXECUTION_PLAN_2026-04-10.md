# Public CMS Site Scope Hardening Execution Plan

Date: 2026-04-10

## Goal

Bring public CMS modules to a consistent tenant+site model so one tenant can safely operate multiple public sites without content bleed.

## Scope

- pages
- menus
- sliders
- media

## Slice 1

- add `site_id` to domain models and persistence
- create indexes for `tenant_id + site_id`
- backfill migration for existing records

## Slice 2

- update admin handlers and repositories
- update public handlers and repositories
- ensure `tenantfilter.WithTenantAndSite(...)` is used where required

## Slice 3

- reseed `aystek` public CMS content against resolved site context
- verify `aysteknolojileri.com` and `www.aysteknolojileri.com`
- verify admin preview/edit flows for the same site

## Closure Criteria

- no tenant-level bleed between multiple sites of the same tenant
- public page/menu/slider/media fetches are site-scoped
- `aystek` public site still renders correctly after migration
