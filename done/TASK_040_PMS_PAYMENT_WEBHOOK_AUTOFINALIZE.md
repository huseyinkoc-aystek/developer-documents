# TASK 040 - PMS Payment Webhook Auto-Finalize

## Scope
Sprint-13.C kapsamında provider callback sonrası PMS portal finalize server-side otomatik çalışır hale getirildi.

## Implementation
- Payments intents service:
  - succeeded callback hook eklendi (`WithSucceededCallback`)
  - callback akışında `succeeded` durumda hook tetiklenir
- API wiring (`cmd/api/main.go`):
  - `entity_type == pms_portal_payment` için PMS portal finalize çağrılır
- PMS portal service:
  - `FinalizeIntentPayment(...)` public API eklendi
  - idempotency key ile finalize guard korundu
  - `pms.payment.received` outbox duplicate guard eklendi

## Guarantees
- Webhook tekrar gelse de finalize idempotent davranır.
- Allocation tekrar çalıştırılmaz.
- Outbox duplicate event emit engellenir.

## Notes
- Core payment intent akışı korunmuştur, PMS tarafı yalnız bridge/orchestration ile entegre edilmiştir.
