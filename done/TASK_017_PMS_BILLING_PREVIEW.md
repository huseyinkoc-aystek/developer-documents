# TASK_017 — PMS Billing Preview + Plan Selection

## Kapsam
- Dry-run tahakkuk endpointi eklendi:
  - `POST /admin/pms/billing/dues/preview`
- Publish ile ayni hesap motoru kullanildi.
- Plan secimi `effective_from <= periodStart` kuralina gore yapildi.

## Teknik Notlar
- `computePreview` fonksiyonu hem preview hem publish tarafinda ortak.
- `plan_id` verilmezse aktif ve doneme uygun plan secilir.
- `area_based` planlarda unit alan bilgisi okunur; alan yoksa satir skip edilir (default_amount yoksa).
- Preview sonucu:
  - `created_count`, `skipped_count`, `total_due`
  - unit bazli `will_create/skip_reason`.

## UI
- `/admin/pms/billing` ekranina `Tahakkuk Onizle` aksiyonu eklendi.
- Onizleme bolumu unit bazli liste ve toplam degerleri gosterir.

## Test
- `internal/modules/pms/billing/bridge/service/service_test.go`:
  - preview/publish uyumu
  - effective_from plan secimi
  - idempotency davranisi korunumu
