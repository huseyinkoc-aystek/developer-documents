# Bulk SMS Audience Information Architecture And Locale Preset Normalization Checkpoint

Date: 2026-04-23

## Scope
- Bulk SMS `Hedef Kitle` sekmesindeki dağınık ve sıkışık filtre yapısını düzeltmek
- System preset kartlarında aktif locale dışında görünen EN/TR varyantlarını normalize etmek

## Completed
- `Hedef Kitle` sekmesi, gerçek bir operasyon yüzeyine dönüştürüldü:
  - sol panel: filtreler, seçim ve sayfalama
  - sağ panel: filtrelenmiş audience listesi
- Sayfa üstündeki global analytics/governance kartları audience sekmesinden çıkarıldı
- Yeni audience başlık/içerik metinleri için i18n anahtarları eklendi
- SMS service tarafında system preset kayıtları canonical `system_key` ve locale-neutral adlarla normalize edilmeye başlandı
  - `approved_locale` preset adı artık locale’a göre tekrar eden `Onaylı EN/TR Aboneler` yerine `Onaylı Aboneler`
  - mevcut eski sistem preset kayıtları da upsert sırasında normalize ediliyor

## Files
- `/workspace/node/aystek-core-next/app/(dashboard)/admin/integrations/bulk-sms/page.tsx`
- `/workspace/go/aystek-core/internal/modules/sms/service/service.go`
- `/workspace/go/aystek-core/internal/seeds/translation_seed.go`

## Verification
- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/integrations/bulk-sms/page.tsx'`
- `cd /workspace/go/aystek-core && go test ./internal/modules/sms/... ./internal/seeds ./cmd/api`

## Expected runtime effect
- Audience sekmesi daha sade ve nefesli görünür
- Aktif locale dışında system preset kartları görünmez
- Eski EN/TR isimli system preset kayıtları zamanla canonical görünüme oturur
