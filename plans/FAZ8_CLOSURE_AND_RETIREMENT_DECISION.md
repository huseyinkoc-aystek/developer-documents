# Faz 8 Closure And Retirement Decision

Bu doküman, Faz 8 sonunda oluşan merkezi audit, export ve governance zincirini özetler ve `stock_integrations` legacy yüzeyi için son migration-retirement kararını netleştirir.

## Kapanış Özeti

- Telematics ve integrations operasyon zinciri artık audit, export, scheduled governance, notice, acknowledgement ve remediation döngüsüyle kapanmıştır.
- `stock/locations` ownership sızıntısı kapatılmıştır; telematics yazma modeli stock domainine doğrudan bağımlı değildir.
- `stock_integrations` aktif ürün yüzeyi olmaktan çıkarılmış, `compatibility-read-only` katmanına çekilmiş ve ölçülebilir retirement lifecycle ile yönetilir hale getirilmiştir.

## Final Karar

`stock_integrations` modülü artık yeni geliştirme yüzeyi değildir.

Karar:
- Yeni operasyon ve ürün geliştirmeleri `integrations` ve `telematics` merkezlerinde devam eder.
- `stock_integrations` yalnız geriye dönük uyumluluk için korunur.
- Modül read-only compatibility katmanı olarak çalışır.
- Runtime usage sıfıra indiğinde ve belirlenen gözlem süresi boyunca rollback ihtiyacı oluşmadığında modül tamamen kaldırılabilir.

## Bugünkü Runtime Durumu

- Legacy API yüzeyi deprecation header’ları ile yayınlanır.
- Write operasyonları bloklanır ve audit event bırakır.
- Read erişimleri ölçülür.
- Kullanım belirli süre sıfır kalırsa compatibility gate otomatik disable edilebilir.
- Gerekirse admin yüzeyinden tek tık re-enable yapılabilir.
- Retirement lifecycle activity logs, stock dashboard ve integrations dashboard içinde görünürdür.
- `legacy_retirement` ve `ops_policy_changes` export/reporting zinciri tamamlanmıştır.

## Kalan Aktif Sınır

Aktif ve kalıcı:
- `stock/locations`
- `stock/reports`
- `telematics/*`
- `integrations/*`

Uyumluluk için geçici:
- `internal/modules/stock/integrations/*`
- ilgili legacy permission ve API contract yüzeyi
- frontend tarafındaki deprecated `stock-ecosystem` legacy integration client alanları

## Tam Retirement İçin Exit Criteria

Bir sonraki kaldırma kararı şu koşullar birlikte sağlandığında verilmelidir:

1. `legacy_integrations_access_24h == 0` durumu tanımlı gözlem penceresi boyunca korunmalı.
2. Auto-disable sonrası `compatibility_reenabled` olayı gözlenmemeli.
3. `deprecated_write_blocked` olayları operasyonel gürültü seviyesinin altına inmeli.
4. Stock tarafında aktif UI veya job akışı bu modüle ihtiyaç duymamalı.
5. Ops policy ve legacy retirement export’larında kaldırma riskini gösteren açık notice kalmamalı.

## Uygulanacak Son Retirement Kararı

Exit criteria sağlandığında:
- `stock_integrations` backend module kaldırılır.
- legacy permission kayıtları temizlenir.
- deprecated frontend client namespace kaldırılır.
- deprecation header/middleware kaldırılır.
- retirement karar dokümanı final removal notu ile güncellenir.

## Rollback Politikası

Tam kaldırma öncesi aşamada rollback yalnız compatibility gate re-enable ile yapılır.

Tam kaldırma sonrası rollback:
- yalnız planlı hotfix veya restore kararı ile yapılmalıdır
- otomatik geri açma politikası artık kullanılmamalıdır

## Referanslar

- [FAZ8_STOCK_INTEGRATIONS_MIGRATION_INVENTORY.md](/workspace/aystek-core-dev-docs/plans/FAZ8_STOCK_INTEGRATIONS_MIGRATION_INVENTORY.md)
- [TELEMATICS_FOUNDATION_ADR_002_LOCATION_OWNERSHIP_AND_COMPATIBILITY.md](/workspace/aystek-core-dev-docs/plans/TELEMATICS_FOUNDATION_ADR_002_LOCATION_OWNERSHIP_AND_COMPATIBILITY.md)
