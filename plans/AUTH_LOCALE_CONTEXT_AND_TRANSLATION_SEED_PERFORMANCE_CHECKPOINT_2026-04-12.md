# Auth Locale, Context, and Translation Seed Performance Checkpoint

Date: 2026-04-12

## Scope
- Login dil seçiminin modern dropdown/listbox yapısına taşınması
- Tenant selector dil seçiminin modern dropdown/listbox yapısına taşınması
- Login ve tenant selector çoklu dil davranışının locale-aware fallback ile sıkılaştırılması
- Site değişiminde stale UI/context riskinin azaltılması
- Menu cache’in tenant/site/lang bağlamına göre ayrıştırılması
- Translation fetch limit kaynaklı eksik sözlük yükleme riskinin kapatılması
- Translation seed performansının per-key read/write modelinden bulk-upsert modeline taşınması

## Completed
- `login` ekranında dil seçimi modern dropdown/listbox olarak konumlandırıldı.
- `tenant-selector` ekranında dil seçimi modern dropdown/listbox olarak hizalandı.
- `login` dil değişiminde `router.refresh()` eklendi.
- `tenant-selector` dil değişiminde `ensureLang(... publicOnly: false)` + `router.refresh()` korunarak locale değişiminin etkisi güçlendirildi.
- `useMenu` cache anahtarı `tenant + site + lang` bazına taşındı.
- `SiteSwitcher` site değişiminde sadece client refresh yerine tam navigation/reload yapacak şekilde sıkılaştırıldı.
- Client translation fetch limit `6000 -> 20000` yükseltildi.
- Server translation preload limit `6000 -> 20000` yükseltildi.
- Translation seed path’i tek tek existence read + upsert yerine prefetch + `BulkWrite` modeline taşındı.

## Verification
- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(auth)/login/page.tsx' --file 'app/tenant-selector/page.tsx' --file 'components/i18n/LocaleSelectSwitcher.tsx' --file 'components/context/SiteSwitcher.tsx' --file 'hooks/useMenu.ts' --file 'lib/api/translations.ts' --file 'lib/api/translations.server.ts' --file 'lib/i18n/localeText.ts' --file 'app/(dashboard)/admin/users/page.tsx'`
- `cd /workspace/go/aystek-core && go test ./internal/seeds ./cmd/api`

## Key Root Causes
- Translation payload fetch zinciri ilk `6000` kayıtla sınırlıydı; tenant başına `8534+` kayıt olan ortamlarda bazı key’ler client’a hiç ulaşmıyordu.
- Menu cache yalnızca `lang` ile anahtarlanıyordu; tenant/site değişiminde stale menu sonucu oluşabiliyordu.
- Translation seed her key için existence read + upsert yaptığı için Mongo tarafında tenant başına yaklaşık iki katı işlem üretiyordu.

## Runtime Note
- Bu değişikliklerin canlı etkisini görmek için backend restart ve tarayıcı hard refresh önerilir.
