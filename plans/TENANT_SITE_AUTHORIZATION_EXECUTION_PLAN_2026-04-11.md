# Tenant & Site Authorization Execution Plan — 2026-04-11

## Goal
Move from global admin/platform access assumptions to explicit tenant/site-scoped authorization.

## Why
Current system now supports stronger tenant/site selection and bootstrap flows, but it still does not model which user can access which tenant and which sites under that tenant.

## Slice 1
- Add tenant membership model
- Add site membership / site grant model
- Seed admin/platform defaults safely
- Enforce tenant access in `/me/context` and bootstrap flows

## Slice 2
- Enforce authorized site lists in selector, header switchers, and context middleware
- Prevent cross-site access even within same tenant unless granted
- Add permission-aware UX for create/read/update actions

## Slice 3
- Admin UI for membership assignment
- User -> tenant assignments
- User -> site assignments
- Audit log for grant/revoke operations

## Slice 4
- Seed / migration / backfill plan for existing users and stores
- Regression verification matrix for admin, tenant admin, limited operator
