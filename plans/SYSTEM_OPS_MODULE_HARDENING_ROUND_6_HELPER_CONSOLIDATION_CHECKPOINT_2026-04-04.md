# System Ops Module Hardening Round 6 Checkpoint

Date: 2026-04-04

## Scope Completed

- `components/system-ops/` altında yeni ortak `SystemOpsBanner` eklendi.
- Eski helper katmanı yeni shared helper setine konsolide edildi:
  - `StatusBadge` -> `SystemOpsStatusBadge`
  - `LiveLogTerminal` -> `SystemOpsLogTerminal`
  - `MaintenanceBanner` -> `SystemOpsBanner`
  - `SystemHealthCard` -> yeni status badge çizgisi
  - `UpdateProgressCard` -> `SystemOpsStatusBadge`
  - `AuditEntryRow` -> `SystemOpsStatusBadge`
- Barrel export güncellendi.

## Pages Moved Further To Shared Stack

- `/admin/system/license`
- `/admin/system/node`
- `/admin/system/health`
- `/admin/system/backups`
- `/admin/system/updates`
- `/admin/system/rollbacks`
- `/admin/system/migrations`

## Shared Patterns Applied

- ortak action button class kullanımı
- ortak banner kullanımı
- ortak empty-state kullanımı
- ortak table shell kullanımı
- node disk görünümünde shared disk bar kullanımı

## Verification

- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/system/license/page.tsx' --file 'app/(dashboard)/admin/system/node/page.tsx' --file 'app/(dashboard)/admin/system/health/page.tsx' --file 'app/(dashboard)/admin/system/backups/page.tsx' --file 'app/(dashboard)/admin/system/updates/page.tsx' --file 'app/(dashboard)/admin/system/rollbacks/page.tsx' --file 'app/(dashboard)/admin/system/migrations/page.tsx' --file 'app/(dashboard)/admin/system/page.tsx' --file 'app/(dashboard)/admin/system/updates/[id]/page.tsx' --file 'components/system-ops/SystemOpsBanner.tsx' --file 'components/system-ops/StatusBadge.tsx' --file 'components/system-ops/LiveLogTerminal.tsx' --file 'components/system-ops/MaintenanceBanner.tsx' --file 'components/system-ops/SystemHealthCard.tsx' --file 'components/system-ops/UpdateProgressCard.tsx' --file 'components/system-ops/AuditEntryRow.tsx' --file 'components/system-ops/index.ts'`
- `cd /workspace/go/aystek-core && go test ./internal/modules/system_ops/... ./cmd/api`

## Recommended Next Step

- `system/maintenance`, `system/version`, `system/overview`, `system/updates/start`, `system/rollbacks/start` yüzeylerinde kalan lokal info blocks ve form feedback alanlarını da `SystemOpsBanner` ile hizalamak.
- ardından `components/system-ops/` için shared column/header helpers çıkarıp tablo başlık tekrarlarını azaltmak.
