# TASK_014_PMS_DUES_PLANS_UI

## Scope
Sprint-7.A.1 kapsamında PMS Dues Plan yönetimi için admin UI akışı tamamlandı.

## Completed
- `GET/POST/PUT/ACTIVATE/ARCHIVE` akışlarıyla uyumlu Dues Plan list/create/edit ekranları tamamlandı.
- Sayfalarda RBAC gating eklendi:
  - `pms_billing_plans.read`
  - `pms_billing_plans.create`
  - `pms_billing_plans.update`
  - `pms_billing_plans.activate`
  - `pms_billing_plans.archive`
- Liste ekranında filtreler (`property`, `status`) ve satır aksiyonları (edit/activate/archive) permission-aware hale getirildi.
- Billing overview ekranında aktif plan göstergesi eklendi:
  - aktif plan varsa `code` badge metni
  - yoksa “no active plan” bilgisi
- Billing overview üst aksiyonunda plans yönetim bağlantısı eklendi (`Manage Plans`).
- UI metinleri i18n key üzerinden taşındı; yeni anahtarlar seed’e eklendi.

## Files
- `node/aystek-core-next/app/(dashboard)/admin/pms/billing/dues-plans/page.tsx`
- `node/aystek-core-next/app/(dashboard)/admin/pms/billing/dues-plans/new/page.tsx`
- `node/aystek-core-next/app/(dashboard)/admin/pms/billing/dues-plans/[id]/page.tsx`
- `node/aystek-core-next/app/(dashboard)/admin/pms/billing/page.tsx`
- `go/aystek-core/internal/seeds/translation_seed.go`

## Validation
- Target ESLint PMS billing pages + `lib/api/pms.ts` + `lib/types/pms.ts`
- `go test ./internal/modules/pms/... ./cmd/api`

## Notes
- PMS dışı `go test ./...` kırıkları kapsam dışı bırakıldı.
