# TASK 026 Summary

PMS webhook integration surface is available with tenant/site-scoped subscriptions, signed deliveries, and event dispatch hooks for key PMS finance/ops events.

## Delivered
- Subscription CRUD + test endpoint.
- HMAC signature (`X-Aystek-Signature`).
- Delivery logging with retry metadata.
- Dispatch bridge connected to PMS event emitters.
