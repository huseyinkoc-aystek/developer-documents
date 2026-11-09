# System Ops Locale-Aware Fallback Fix Checkpoint

Date: 2026-04-05

## Problem
- `System Ops` modülünde bazı yüzeyler Türkçe, bazı yüzeyler İngilizce görünüyordu.
- Kök neden: bazı ekranlar yeni `system_ops.*` anahtarlarını kullanıyordu; bu anahtarlar DB'de henüz reconcile edilmemişse `t(...)` İngilizce fallback metne düşüyordu.
- Aynı modülde rail tarafının Türkçe görünmesinin nedeni, orada daha eski ve DB'de mevcut `menu.system_ops.*` anahtarlarının kullanılmasıydı.

## Fix
- Locale-aware fallback helper eklendi:
  - `lib/system-ops/text.ts`
- Bu helper aşağıdaki yüzeylere uygulandı:
  - `components/system-ops/SystemOpsPageShell.tsx`
  - `components/system-ops/SystemOpsRail.tsx`
  - `app/(dashboard)/admin/system/page.tsx`
  - `app/(dashboard)/admin/system/overview/page.tsx`

## Result
- DB tarafında yeni translation key henüz yoksa bile TR kullanıcı artık İngilizce fallback görmüyor.
- System Ops ekranları seed timing’ine daha az bağımlı hale geldi.

## Verification
- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/system/page.tsx' --file 'app/(dashboard)/admin/system/overview/page.tsx' --file 'components/system-ops/SystemOpsRail.tsx' --file 'components/system-ops/SystemOpsPageShell.tsx' --file 'lib/system-ops/text.ts'`
