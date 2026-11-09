# Login Tenant Selector Locale Refinement And Scope UI Checkpoint

Date: 2026-04-12

## Completed
- Login dil seçimi üstten kaldırılıp form altına, ortalı ve modern dropdown/listbox olarak taşındı.
- Tenant selector dil listesi `getOperationalLanguages(...)` ile aktif/tekil dillerden üretilir hale getirildi.
- Tenant selector ve login yüzeylerinde kalan sert-coded İngilizce fallback’ler locale-aware fallback ile kapatıldı.
- Locale dropdown katmanı için z-index ve renk uyumu güçlendirildi.
- `useMenu` cache anahtarı `tenant + site + lang` bazında çalışmaya devam edecek şekilde lint-clean hale getirildi.
- Site switch sonrası stale görünüm riskini azaltmak için tam navigation/reload davranışı korundu.
- Translation fetch limit client/server tarafında `20000` seviyesine çıkarıldı.
- Translation seed performansı bulk prefetch + batch upsert yaklaşımıyla hızlandırıldı.
- Users ekranında platform admin vs tenant admin ayrımı daha görünür hale getirildi.
- Users batch access paneline daha yoğun scope summary kartları ve matrix kolonları eklendi.

## Verification
- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(auth)/login/page.tsx' --file 'app/tenant-selector/page.tsx' --file 'components/i18n/LocaleSelectSwitcher.tsx' --file 'components/context/SiteSwitcher.tsx' --file 'hooks/useMenu.ts' --file 'lib/api/translations.ts' --file 'lib/api/translations.server.ts' --file 'lib/i18n/localeText.ts' --file 'app/(dashboard)/admin/users/page.tsx'`
- `cd /workspace/go/aystek-core && go test ./internal/seeds ./cmd/api`

## Notes
- Translation seed yavaşlığının ana nedeni MongoDB’den çok önceki seed algoritmasının per-key read + per-key upsert modeli olmasıydı.
- Canlı görünümde değişiklikleri net görmek için backend restart ve tarayıcı hard refresh önerilir.
- Yetki modelinin sonraki slice’ında dashboard/modül görünürlüğünün membership scope’a daha doğrudan bağlanması devam ettirilebilir.
