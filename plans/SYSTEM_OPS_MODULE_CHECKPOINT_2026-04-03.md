# System Ops Checkpoint

Date: 2026-04-03
Purpose: VSCode / Copilot reset durumunda hızlı devam noktası bırakmak.

## Verified Gaps

- Maintenance UI sadece `tr/en` destekliyordu.
- `system-ops` API client ve types katmanı drift halindeydi.
- Update start cevabı `record_id` dönmediği için detail yönlendirmesi kırılabiliyordu.
- Update log SSE akışı sadece `job_id` ile çözülebiliyordu.
- System Ops sayfaları arasında ortak operasyon rail’i yoktu.

## Implemented Fixes

- Dynamic multi-language maintenance editor
- Shared `SystemOpsRail`
- API normalization layer
- `record_id` support in update/rollback start responses
- log stream resolution by `job_id` or `record_id`

## Primary Files

- `/workspace/node/aystek-core-next/lib/api/system-ops.ts`
- `/workspace/node/aystek-core-next/lib/types/system-ops.ts`
- `/workspace/node/aystek-core-next/components/system-ops/SystemOpsRail.tsx`
- `/workspace/node/aystek-core-next/app/(dashboard)/admin/system/maintenance/page.tsx`
- `/workspace/go/aystek-core/internal/modules/system_ops/handler/admin.go`
- `/workspace/go/aystek-core/internal/modules/system_ops/service/service.go`
- `/workspace/go/aystek-core/internal/modules/system_ops/dto/system_ops.go`
- `/workspace/go/aystek-core/internal/seeds/translation_seed.go`

## Remaining Nice-to-Haves

- rollbacks/migrations pages için rail bağını tamamlamak
- locale-aware formatting helper çıkarımı
- shared page shell
