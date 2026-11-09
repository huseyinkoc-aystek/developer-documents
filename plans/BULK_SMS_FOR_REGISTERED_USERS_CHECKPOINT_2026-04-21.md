# Bulk SMS For Registered Users Checkpoint

Date: 2026-04-21

## Scope

Tenant kapsamındaki kayıtlı kullanıcılar arasından telefon numarası olan abonelere toplu SMS gönderimi için backend + Next admin yüzeyi eklendi.

## Backend

- SMS modülüne audience listeleme endpoint’i eklendi:
  - `GET /admin/sms/broadcasts/audience`
- SMS modülüne bulk send endpoint’i eklendi:
  - `POST /admin/sms/broadcasts/send`
- Audience kaynağı `users` domaini oldu; `phone_number` bulunan tenant-scope kullanıcılar listeleniyor.
- Bulk send akışı:
  - kullanıcı ID listesi alır
  - tenant kapsamı içinde kullanıcıları çözer
  - mevcut SMS provider/policy zincirini kullanarak tekil gönderim yapar
  - sonuçları toplu response ile döner
  - activity log özeti yazar

## Frontend

- Yeni integrations yüzeyi eklendi:
  - `/admin/integrations/bulk-sms`
- Sayfa düzeni:
  - integrations header
  - audience KPI kartları
  - filtrelenebilir audience listesi
  - ayrı composer paneli
  - ayrı sonuç tablosu
- `product_type` benzeri serbest alan yok; kullanıcı seçimi audience listesi üzerinden yapılıyor.

## Seeds

- yeni menu anahtarı:
  - `menu.integrations.bulk_sms`
- yeni SMS bulk çeviri anahtarları:
  - `sms.bulk.*`
- permission seed tarafında SMS için:
  - `read`
  - `bulk_read`
  - `send_bulk`

## Validation

- `cd /workspace/go/aystek-core && go test ./internal/modules/sms/... ./internal/modules/users/... ./internal/seeds ./cmd/api`
- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/integrations/bulk-sms/page.tsx' --file 'lib/api/sms.ts' --file 'lib/types/sms.ts'`

## Notes

- Audience kaynağı newsletter değil, kayıtlı kullanıcılardır.
- Bu slice dağıtım, onay ve campaign segmentation katmanına girmeden, mevcut SMS provider yapısıyla uyumlu minimal ama gerçek bir bulk SMS operasyon yüzeyi sağlar.
