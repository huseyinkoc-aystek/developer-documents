# TASK_029_PMS_WEBHOOKS_DELIVERY_WORKER

## Scope
Sprint-11.B kapsamında webhook teslimat hattı claim-lock + retry worker modeline harden edildi.

## Completed
- Delivery durumları genişletildi:
  - `pending|processing|sent|failed|dead_letter`
- Atomic claim eklendi (`findOneAndUpdate`):
  - `ClaimNextDelivery`
  - `locked_at`, `locked_by`, `lock_expires_at`
  - stuck delivery recovery
- Deterministic retry/backoff:
  - `1m, 5m, 30m, 2h, 12h`
- Max retry sonrası dead-letter akışı.
- List endpoint delivery summary döndürür hale getirildi.
- Admin webhooks sayfasında sent/failed/dead/pending özet kartları eklendi.

## Files
- `go/aystek-core/internal/modules/pms/webhooks/model.go`
- `go/aystek-core/internal/modules/pms/webhooks/repository.go`
- `go/aystek-core/internal/modules/pms/webhooks/service.go`
- `go/aystek-core/internal/modules/pms/webhooks/handler.go`
- `go/aystek-core/internal/modules/pms/webhooks/service_test.go`
- `go/aystek-core/internal/modules/pms/webhooks/indexes.go`
- `node/aystek-core-next/app/(dashboard)/admin/pms/webhooks/page.tsx`
- `node/aystek-core-next/lib/api/pms.ts`
- `node/aystek-core-next/lib/types/pms.ts`
- `go/aystek-core/internal/seeds/translation_seed.go`

## Validation
- `go test ./internal/modules/pms/... ./cmd/api`
- `npx eslint app/(dashboard)/admin/pms/webhooks/page.tsx lib/api/pms.ts lib/types/pms.ts`
