# Device Command Pending Override Queue Checkpoint

Date: 2026-04-21

## Scope
- Pending override isteklerini execution stream içinden ayırıp ayrı governance kuyruğuna taşımak
- Approver ve approval timestamp bilgisini görünür read-model olarak eklemek
- Policy audit/history görünürlüğünü actor bilgisiyle güçlendirmek

## Completed
- `CommandExecution` modeline `approved_at` ve `approved_by` alanları eklendi.
- Execution dispatch update zinciri approval metadata yazacak şekilde genişletildi.
- Override approval service akışı artık dispatch sırasında approver ve approval timestamp bilgisi kaydediyor.
- Registry sayfasına ayrı `Pending Override Queue` bölümü eklendi.
- Queue satırlarında şu alanlar görünür hale geldi:
  - command intent
  - vendor key
  - product type
  - phone
  - override reason
  - override ticket ref
  - requested at
  - approved by
  - approved at
- Queue içinden `Approve Override` aksiyonu çalışır hale getirildi.
- Policy audit/history tablosuna `actor_id` kolonu eklendi.
- Next type ve API parser katmanı approval metadata ve actor bilgisiyle hizalandı.
- Yeni i18n anahtarları translation seed’e eklendi.

## Files
- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/model/model.go`
- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/repository/interface.go`
- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/repository/mongo.go`
- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/service/service.go`
- `/workspace/go/aystek-core/internal/modules/stock/integrations/device_commands/service/service_test.go`
- `/workspace/go/aystek-core/internal/seeds/translation_seed.go`
- `/workspace/node/aystek-core-next/lib/types/stock-ecosystem.ts`
- `/workspace/node/aystek-core-next/lib/api/stock-ecosystem.ts`
- `/workspace/node/aystek-core-next/app/(dashboard)/admin/integrations/device-commands/registry/page.tsx`

## Verification
- `cd /workspace/go/aystek-core && go test ./internal/modules/stock/integrations/device_commands/... ./internal/seeds ./cmd/api`
- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/integrations/device-commands/SmsCommandCenterPage.tsx' --file 'app/(dashboard)/admin/integrations/device-commands/registry/page.tsx' --file 'components/layout/LaunchpadMenu.tsx' --file 'lib/api/stock-ecosystem.ts' --file 'lib/types/stock-ecosystem.ts'`

## Outcome
- Pending override akışları artık ayrı governance kuyruğunda görünür.
- Approval izi execution read-model üzerinde kaybolmuyor.
- Override approval ve policy history daha denetlenebilir hale geldi.
