# TASK_052 — Stock Portal UI Parity

## Kapsam
Sprint-7 parity hardening sonrası portal stock endpointleri için Next.js portal ekranları ve temel smoke doğrulaması tamamlandı.

## Uygulanan değişiklikler

### 1) Portal giriş yüzeyi
- `node/aystek-core-next/app/portal/page.tsx`
- `node/aystek-core-next/app/portal/stock/page.tsx`

Portal kullanıcıları için tek noktadan erişim sağlandı:
- `/portal/stock/sales`
- `/portal/stock/field-reports/assets`
- `/portal/stock/locations/current`

### 2) Portal stock ekranları
- `node/aystek-core-next/app/portal/stock/sales/page.tsx`
  - hardcoded fallback metinleri kaldırıldı
  - tüm metinler i18n anahtarından okunuyor
- `node/aystek-core-next/app/portal/stock/field-reports/assets/page.tsx`
  - varlık talebi oluşturma (portal endpointi)
  - son talepler listesi
- `node/aystek-core-next/app/portal/stock/locations/current/page.tsx`
  - anlık konum listesi + araç filtresi
  - pagination/limit/search toolbar uyumu

### 3) i18n seed tamamlama
- `go/aystek-core/internal/seeds/translation_seed.go`

Eklenen anahtarlar:
- `portal.title`, `portal.subtitle`
- `stock.portal.title`, `stock.portal.subtitle`
- `stock.portal.notice.saved`
- `stock.portal.sales.*`
- `stock.portal.field_reports.assets.*`
- `stock.portal.locations.current.*`

### 4) Smoke test
- `node/aystek-core-next/tests/stock-portal-pages.test.ts`

Portal stock sayfa dosyalarının varlığını doğrulayan temel smoke test eklendi.

## Doğrulama

### Frontend
```bash
cd /workspace/node/aystek-core-next
npx eslint app/portal/page.tsx \
  app/portal/stock/page.tsx \
  app/portal/stock/sales/page.tsx \
  app/portal/stock/field-reports/assets/page.tsx \
  app/portal/stock/locations/current/page.tsx \
  tests/stock-portal-pages.test.ts
npx tsx --test tests/stock-portal-pages.test.ts
```

### Backend (stock kapsamı)
```bash
cd /workspace/go/aystek-core
go test ./internal/modules/stock/... ./cmd/api
```

## Sonuç
- Portal stock endpointlerinin UI karşılığı tamamlandı.
- i18n seed eksikleri kapatıldı.
- Kapsam PMS dışı modüllere taşmadan, mevcut plan akışına sadık şekilde kapatıldı.
