# TASK_016 — PMS Portal Work Orders

## Kapsam
- Portal work-order akisi urun seviyesine genisletildi.
- Yeni endpointler:
  - `GET /portal/pms/me/work-orders?unit_id=&status=&page=`
  - `GET /portal/pms/me/work-orders/:id`
  - `POST /portal/pms/me/work-orders/:id/comment`

## Teknik Notlar
- `me` kullanicisi sadece aktif occupancy ile bagli oldugu unitler icin liste/detay/comment islemi yapabilir.
- `unit_id` filtreli listede erisim yoksa servis `errPortalNoAccess` dondurur.
- Portal rotalarinda hafif in-memory rate limit korunur.
- Comment aksiyonu audit logda `portal_action=true` etiketi ile loglanir.

## Test
- `internal/modules/pms/portal/service_test.go`:
  - occupancy yoksa erisim hatasi
  - occupancy varsa liste/detay
  - comment basarili akisi
