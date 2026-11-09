# SMS Device Command Intent And Auth Locale Refinement Checkpoint

Date: 2026-04-17

## Closed In This Slice

- Login ekranındaki auth hata kodu artık ham key olarak görünmemeli.
- Login ve tenant-selector dil seçimi tekil, modal tabanlı, daha tutarlı bir picker üzerinden çalışıyor.
- Kullanılmayan `LocaleSelectSwitcher` kaldırıldı.
- Stock device command execution artık sadece serbest mesaj gönderimi değil:
  - `command_intent`
  - `vendor_key`
  - `product_type`
  alanlarını da taşıyor.
- Backend artık `template_id` verilmemişse `command_intent + vendor_key + product_type` üzerinden en uygun command template’i çözebiliyor.
- Execution kayıtları artık çözülmüş intent/vendor/product metadata’sını saklıyor.
- Integrations SMS command center quick-send akışı resolved target metadata’sını backend’e geçiriyor.
- Dashboard ve users yönetim yüzeyinde membership scope affordance’ları bir adım daha görünür hale geldi.

## Validation

- `cd /workspace/go/aystek-core && go test ./internal/modules/stock/integrations/device_commands/... ./internal/seeds ./cmd/api`
- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(auth)/login/page.tsx' --file 'app/tenant-selector/page.tsx' --file 'components/i18n/LocaleModalPicker.tsx' --file 'app/(dashboard)/page.tsx' --file 'app/(dashboard)/admin/users/page.tsx' --file 'app/(dashboard)/admin/integrations/device-commands/SmsCommandCenterPage.tsx' --file 'lib/api/stock-ecosystem.ts' --file 'lib/types/stock-ecosystem.ts'`

## Remaining Best Next Step

- Device vendor command registry için read-model / admin management yüzeyini büyütmek
- Serbest provider seçiminden daha da çıkıp resolved route’u baskın hale getirmek
- Dashboard modül görünürlüğünü membership scope ile daha sert bağlamak
