# System Ops Next I18N Completeness Checkpoint

Date: 2026-04-05

## Scope
- Next.js `System Ops` sayfaları
- `components/system-ops/` shared UI katmanı
- `translation_seed.go` içindeki `system_ops.*`, `menu.system_ops.*`, `common.enabled`, `common.disabled`

## Completed
- System Ops sayfalarında kullanılan i18n anahtarları seed ile karşılaştırıldı.
- Eksik kalan anahtarlar eklendi:
  - `common.enabled`
  - `common.disabled`
  - `system_ops.license.expired`
  - `system_ops.audit.filter_help`
  - `system_ops.audit.table_help`
  - `system_ops.maintenance.placeholder_generic`
  - `system_ops.updates.detail.back_help`
- `overview` sayfasındaki maintenance locale-map gösterimi düzeltildi.
- `updates/[id]` log saat gösterimi locale-aware hale getirildi.
- `audit` ve `updates/[id]` semantik fallback metinleri kendi i18n anahtarlarına taşındı.
- `maintenance` locale placeholder davranışı tek dilli fallback yerine genel çok dilli yordamla hizalandı.

## Verification
- System Ops Next i18n anahtar taraması: `NO_MISSING_KEYS`
- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/system/page.tsx' --file 'app/(dashboard)/admin/system/overview/page.tsx' --file 'app/(dashboard)/admin/system/maintenance/page.tsx' --file 'app/(dashboard)/admin/system/audit/page.tsx' --file 'app/(dashboard)/admin/system/updates/[id]/page.tsx' --file 'components/system-ops/AuditEntryRow.tsx' --file 'components/system-ops/SystemOpsKeyValueRow.tsx' --file 'components/system-ops/SystemOpsMetricCard.tsx' --file 'components/system-ops/SystemOpsNavTile.tsx' --file 'components/system-ops/SystemOpsPageShell.tsx' --file 'components/system-ops/SystemOpsStatList.tsx' --file 'components/system-ops/SystemOpsTableShell.tsx' --file 'components/system-ops/DiskUsageBar.tsx' --file 'components/system-ops/SystemOpsStatusBadge.tsx' --file 'lib/system-ops/formatters.ts'`
- `cd /workspace/go/aystek-core && go test ./internal/seeds ./cmd/api`

## Current State
- System Ops Next tarafında kullanıcıya dönük metinler artık i18n anahtarları üzerinden çözümleniyor.
- Teknik terimler ve backend kaynaklı runtime içerikler hariç modül çoklu dil sözleşmesiyle hizalı.
