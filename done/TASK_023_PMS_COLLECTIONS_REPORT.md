# TASK_023 — PMS Collections Report

## Teslim
- Yeni endpoint:
  - `GET /admin/pms/reports/collections?period=YYYY-MM`

## Hesaplanan metrikler
- `total_due`
- `total_paid`
- `total_remaining`
- `collection_rate`
- aging buckets:
  - `0_30`
  - `31_60`
  - `61_90`
  - `90_plus`

## UI
- `/admin/pms/dashboard` üzerinde collections özet kartı/bucket görünümü eklendi.
