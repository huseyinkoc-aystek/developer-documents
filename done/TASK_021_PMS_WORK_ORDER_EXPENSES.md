# TASK_021 — PMS Work Order Expenses

## Teslim
- Yeni alt modül: `internal/modules/pms/work_orders/expenses/*`
- Endpointler:
  - `POST /admin/pms/work-orders/:id/expenses`
  - `GET /admin/pms/work-orders/:id/expenses`

## Kurallar
- Expense eklenince work-order timeline’a `expense_added` event yazılır.
- Audit log yazılır.
- Outbox event: `pms.work_order.expense_added`.
- `bill_to=unit` ise billing bridge üzerinden manual charge orchestration tetiklenir.

## UI
- `/admin/pms/work-orders` ekranına expense paneli eklendi:
  - vendor seçimi
  - amount/currency/date/note
  - bill_to seçimi
  - mevcut expense listesi
