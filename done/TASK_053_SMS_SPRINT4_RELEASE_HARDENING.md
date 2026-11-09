# TASK_053 — SMS Sprint-4 Release Hardening

## Kapsam
Sprint-4 kapanışı için SMS shared modülünde release hardening tamamlandı:
- Admin UI parity (Providers, Credentials, Policies, Messages, Outbox, Templates)
- Outbox yönetimi (enqueue/list/process/retry)
- Template CRUD (backend + Next.js)
- DLR callback + replay/idempotency zinciri
- Seed/i18n/menu güncellemeleri
- Smoke test seti ve doğrulama komutları

## Tamamlanan Teknik İşler

### 1) Backend Hardening
- SMS outbox endpointleri:
  - `POST /admin/sms/outbox/jobs`
  - `GET /admin/sms/outbox/jobs`
  - `POST /admin/sms/outbox/jobs/process`
  - `POST /admin/sms/outbox/jobs/:job_id/retry`
- Public DLR endpoint:
  - `POST /public/sms/webhooks/:provider_id/dlr`
- Policy failover execution:
  - Primary provider send hatasında failover zinciri ile sonraki provider denenir.
- DLR replay/idempotency:
  - `event_id` yoksa `x-webhook-event-id` header fallback.
  - Duplicate event insert idempotent kabul edilir.

### 2) Templates Modülü (SMS içinde)
- Endpointler:
  - `GET /admin/sms/templates`
  - `POST /admin/sms/templates`
  - `PUT /admin/sms/templates/:id`
  - `DELETE /admin/sms/templates/:id`
- Koleksiyon: `sms_templates`
- Index:
  - unique partial: `tenant_id + site_id + code` (`deleted_at` yokken)
  - list index: `tenant_id + site_id + status + updated_at`

### 3) Next.js Admin UI
- Yeni sayfalar:
  - `/admin/sms/outbox`
  - `/admin/sms/templates`
- Typed client genişletmeleri:
  - `listOutbox`, `enqueueOutbox`, `processOutbox`, `retryOutbox`
  - `listTemplates`, `upsertTemplate`, `deleteTemplate`
- SMS ekranlarında ortak sekme navigasyonu parity:
  - Providers / Credentials / Policies / Messages / Outbox / Templates

### 4) Seeds ve i18n
- `menu_seed.go`:
  - `menu.operations.messaging.sms_outbox`
  - `menu.operations.messaging.sms_templates`
- `translation_seed.go`:
  - `sms.outbox.*`
  - `sms.templates.*`

### 5) Test ve Smoke
- Go service testleri:
  - failover chain execution
  - DLR event id zorunluluğu + header fallback
  - duplicate DLR idempotent davranış
  - outbox process success
  - template upsert/list
- Next.js smoke/unit:
  - `tests/sms-api.test.ts` parser doğrulamaları

## Doğrulama Komutları

### Backend
```bash
cd /workspace/go/aystek-core
go test -count=1 ./internal/modules/sms/... ./cmd/api
```

### Frontend
```bash
cd /workspace/node/aystek-core-next
npx eslint 'app/(dashboard)/admin/sms/**/*.tsx' lib/api/sms.ts lib/types/sms.ts
npm run test -- tests/sms-api.test.ts
```

## Manual QA Checklist
- [ ] `/admin/sms/providers` açılır, provider config kaydedilir.
- [ ] `/admin/sms/credentials` rotate/activate akışı çalışır.
- [ ] `/admin/sms/policies` save + version conflict davranışı doğrulanır.
- [ ] `/admin/sms/messages` filtre ve listeleme çalışır.
- [ ] `/admin/sms/outbox` enqueue + process + retry akışı çalışır.
- [ ] `/admin/sms/templates` create/update/delete akışı çalışır.
- [ ] Public DLR callback duplicate event’te idempotent davranır.
- [ ] Tenant/site context olmadan endpointler hata verir.

## Notlar
- Sprint-4 release hardening SMS boundary içinde tamamlandı.
- Stock/PMS dışı kırıklara müdahale edilmedi.

