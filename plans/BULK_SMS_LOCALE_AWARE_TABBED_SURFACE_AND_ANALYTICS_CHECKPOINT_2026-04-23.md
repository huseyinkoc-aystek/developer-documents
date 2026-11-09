# Bulk SMS Locale-Aware Tabbed Surface And Analytics Checkpoint

Date: 2026-04-23

## Scope

This slice hardened the `Integrations > Bulk SMS` surface to align with project UI and i18n standards.

## Completed

- Removed the cluttered all-in-one governance layout from the main bulk SMS page.
- Kept the main surface focused on:
  - segments
  - audience
  - compose
  - campaigns
- Moved governance/analytics concerns into a dedicated page:
  - `/admin/integrations/bulk-sms/analytics`
- Added a visible analytics entry point from the main bulk SMS surface.
- Introduced locale-aware system preset display behavior:
  - system presets now carry `system_key`
  - active language locale-specific preset is preferred
  - generic approved preset is hidden when a locale-specific approved preset exists for the active language
- Replaced hardwired visible preset labels with localized display labels driven by system keys.
- Added analytics snapshot cards on the main page to preserve situational awareness without clutter.
- Added dedicated analytics page sections for:
  - governance time series
  - provider breakdown
  - template breakdown
  - approver-role breakdown
  - approval history drill-down
- Added local approval history filters on analytics page:
  - status
  - provider
  - template
  - approver role
  - search

## Backend changes

- `BulkSegmentPreset` now supports `system_key`.
- `BulkSegmentPresetSummary` exposes `system_key`.
- System preset generation now sets deterministic keys:
  - `approved_all`
  - `approved_locale`
  - `pending_approval`
  - `all_audience`
- Service layer infers system key for older records that do not yet persist it.

## Frontend changes

- Main page:
  - `/workspace/node/aystek-core-next/app/(dashboard)/admin/integrations/bulk-sms/page.tsx`
- Analytics page:
  - `/workspace/node/aystek-core-next/app/(dashboard)/admin/integrations/bulk-sms/analytics/page.tsx`
- API/types:
  - `/workspace/node/aystek-core-next/lib/api/sms.ts`
  - `/workspace/node/aystek-core-next/lib/types/sms.ts`

## i18n additions

New `sms.bulk.*` keys were added for:

- analytics entry and snapshot
- analytics page title/subtitle
- analytics preview labels
- history drill-down search
- governance snapshot hint
- locale-aware system preset display labels

## Verification

- Next lint passed:
  - bulk SMS main page
  - bulk SMS analytics page
  - sms API/types
- Go test passed:
  - `./internal/modules/sms/...`
  - `./internal/modules/users/...`
  - `./internal/seeds`
  - `./cmd/api`

## Result

Bulk SMS is now substantially closer to the project’s enterprise UI expectations:

- less crowded main surface
- cleaner operational flow
- analytics/governance separated
- active-language-aware system preset visibility
- stronger i18n compliance
