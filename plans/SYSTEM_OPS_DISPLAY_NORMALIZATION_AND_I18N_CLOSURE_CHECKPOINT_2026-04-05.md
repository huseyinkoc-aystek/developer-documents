# System Ops Display Normalization and I18N Closure Checkpoint

Date: 2026-04-05

## Scope
- Next.js `System Ops` kullanıcı yüzeylerinde backend kaynaklı ham display değerleri
- `status`, `backup type`, `trigger source` ve locale-aware zaman gösterimi
- `translation_seed.go` içindeki karşılık anahtarlar

## Completed
- Ortak display helper eklendi:
  - `lib/system-ops/display.ts`
- Normalize edilen alanlar:
  - status label
  - backup type label
  - trigger source label
- Güncellenen yüzeyler:
  - `components/system-ops/SystemOpsStatusBadge.tsx`
  - `components/system-ops/AuditEntryRow.tsx`
  - `app/(dashboard)/admin/system/page.tsx`
  - `app/(dashboard)/admin/system/overview/page.tsx`
  - `app/(dashboard)/admin/system/health/page.tsx`
  - `app/(dashboard)/admin/system/backups/page.tsx`
  - `app/(dashboard)/admin/system/license/page.tsx`
  - `app/(dashboard)/admin/system/updates/[id]/page.tsx`
- Seed’e eklenen anahtarlar:
  - `system_ops.license.status.unlicensed`
  - `system_ops.trigger_source.manual`
  - `system_ops.trigger_source.api`
  - `system_ops.trigger_source.system`
  - `system_ops.trigger_source.scheduler`
  - `system_ops.trigger_source.cli`
  - `system_ops.trigger_source.rollback`
  - `system_ops.trigger_source.guardrail`

## Verification
- System Ops Next anahtar taraması: `NO_MISSING_KEYS`
- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/system/page.tsx' --file 'app/(dashboard)/admin/system/overview/page.tsx' --file 'app/(dashboard)/admin/system/health/page.tsx' --file 'app/(dashboard)/admin/system/backups/page.tsx' --file 'app/(dashboard)/admin/system/license/page.tsx' --file 'app/(dashboard)/admin/system/updates/[id]/page.tsx' --file 'components/system-ops/SystemOpsStatusBadge.tsx' --file 'components/system-ops/AuditEntryRow.tsx' --file 'lib/system-ops/display.ts'`
- `cd /workspace/go/aystek-core && go test ./internal/seeds ./cmd/api`

## Current State
- System Ops modülünde yalnızca statik metinler değil, backend’den gelen operasyonel display değerleri de normalize edilerek çoklu dil yapısına bağlandı.
- Teknik ve runtime içerikler gerektiği yerde ham kalabilir; kullanıcıya dönük operasyon etiketleri artık normalize edilmiş durumda.
