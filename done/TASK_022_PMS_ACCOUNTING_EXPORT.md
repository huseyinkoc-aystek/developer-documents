# TASK_022 — PMS Accounting Export

## Teslim
- Yeni endpoint:
  - `GET /admin/pms/exports/accounting?period=YYYY-MM&format=csv`

## İçerik
- CSV satır tipleri:
  - `charge`
  - `allocation`
  - `vendor_expense`

## Kurallar
- Tenant + site scope zorunlu.
- Export çağrısı audit log’a yazılır.
- Dashboard üzerinden CSV indirme aksiyonu eklendi.
