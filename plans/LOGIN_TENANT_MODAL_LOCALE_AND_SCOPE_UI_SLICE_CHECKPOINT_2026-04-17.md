# Login / Tenant Modal Locale And Scope UI Slice Checkpoint

Date: 2026-04-17

## Scope
- login ekranında raw i18n key sızıntısını kesmek
- login ve tenant-selector dil seçimini modal tabanlı daha uyumlu yüzeye taşımak
- tenant-selector içinde locale-safe fallback ile mixed language/raw key riskini azaltmak
- users batch access panelinde daha hızlı bulk scope uygulama aksiyonu açmak

## Completed
- `app/(auth)/login/page.tsx`
  - locale-safe `localize(...)` helper eklendi
  - auth hata mesajlarında raw key yerine anlamlı fallback zorlandı
  - login görünür metinleri locale-safe fallback ile hizalandı
- `app/tenant-selector/page.tsx`
  - tenant/site seçim yüzeyindeki kritik başlık, aksiyon, placeholder ve empty-state metinleri locale-safe fallback ile hizalandı
- `components/i18n/LocaleModalPicker.tsx`
  - modal tetikleyici daha modern card/pill hibrit görünümüne taşındı
- `app/(dashboard)/admin/users/page.tsx`
  - platform admin için `Tüm görünür tenantları ekle` bulk aksiyonu açıldı
- `go/aystek-core/internal/seeds/translation_seed.go`
  - `users.context.batch.applyAllVisible` translation key eklendi

## Verification
- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(auth)/login/page.tsx' --file 'app/tenant-selector/page.tsx' --file 'components/i18n/LocaleModalPicker.tsx' --file 'app/(dashboard)/admin/users/page.tsx'`
- `cd /workspace/go/aystek-core && go test ./internal/seeds ./cmd/api`

## Notes
- Login tarafındaki `401` hatası i18n problemi değil; auth response raw-key sızıntısı ayrıca kapatıldı.
- Tenant/site selector’da hâlâ karışık dil görülürse en olası neden stale bundle/cache olur. Backend restart + hard refresh ile runtime doğrulama yapılmalı.
