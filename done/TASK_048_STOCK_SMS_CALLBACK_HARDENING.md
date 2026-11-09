# TASK_048 — STOCK SMS Callback Hardening

## Kapsam
- `stock_sms` callback hattına HMAC imza doğrulama ve replay koruması eklendi.
- Güvenlik doğrulaması service katmanına taşındı.

## Uygulanan Değişiklikler
- Yeni callback security akışı:
  - `AuthorizeCallback(ctx, provider, nonce, signature, timestamp, rawBody)`
  - Secret env: `STOCK_SMS_CALLBACK_SECRET`
  - İmza formatı: HMAC-SHA256 (`hex` veya `base64` kabul)
  - Zaman damgası doğrulaması: ±300 saniye pencere
  - Replay koruması: callback receipt insert + unique nonce

- Repository/model genişletmeleri:
  - `stock_sms_callback_receipts` koleksiyonu
  - `CallbackReceipt` modeli
  - `CreateCallbackReceipt` repository metodu

- Indexler:
  - Unique: `tenant_id + site_id + provider + nonce`
  - TTL: `received_at` (24 saat)

## API Uyum
- Mevcut endpoint korunur: `POST /admin/stock/sms/callback`
- Header/body fallback destekleri:
  - Provider: `X-Stock-Sms-Provider` veya `provider`
  - Nonce: `X-Stock-Sms-Nonce` veya `nonce`
  - Signature: `X-Aystek-Signature` / `X-Stock-Sms-Signature` veya `signature`
  - Timestamp: `X-Stock-Sms-Timestamp` veya `timestamp`

## i18n Anahtarları
- `stock.sms.validation.callback_nonce_required`
- `stock.sms.validation.callback_signature_invalid`
- `stock.sms.validation.callback_replay`
- `stock.sms.validation.callback_timestamp_invalid`

## Testler
- `service_test.go`:
  - valid signature
  - replay detection
  - invalid signature
  - timestamp window validation

## Doğrulama
- `go test ./internal/modules/stock/sms/... ./internal/modules/stock/... ./cmd/api` ✅

