# TASK_048 — Özet

- `stock_sms` callback endpointi için HMAC signature + replay protection eklendi.
- Callback güvenliği tenant/site scoped receipt koleksiyonu ile kalıcılaştırıldı.
- Timestamp window ve signature doğrulaması aktif.
- Hata anahtarları translation seed’e eklendi ve stock testleri yeşil.

