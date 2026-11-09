# TASK_052 Summary

## Tamamlananlar
- Portal stock için eksik Next.js sayfaları eklendi:
  - `/portal`
  - `/portal/stock`
  - `/portal/stock/field-reports/assets`
  - `/portal/stock/locations/current`
- `/portal/stock/sales` sayfası hardcode fallback metinlerinden temizlendi.
- Portal stock i18n anahtarları `translation_seed.go` içine eklendi.
- Temel smoke test eklendi (`tests/stock-portal-pages.test.ts`).

## Etki
- Backend tarafında daha önce açılan portal route’lar artık UI’dan erişilebilir.
- Portal flow’da i18n key kaçakları azaltıldı.
- Sprint-7 parity hardening için portal UI boşlukları kapatıldı.

## Doğrulama komutları
```bash
cd /workspace/node/aystek-core-next
npx eslint app/portal/page.tsx app/portal/stock/page.tsx app/portal/stock/sales/page.tsx app/portal/stock/field-reports/assets/page.tsx app/portal/stock/locations/current/page.tsx tests/stock-portal-pages.test.ts
npx tsx --test tests/stock-portal-pages.test.ts

cd /workspace/go/aystek-core
go test ./internal/modules/stock/... ./cmd/api
```
