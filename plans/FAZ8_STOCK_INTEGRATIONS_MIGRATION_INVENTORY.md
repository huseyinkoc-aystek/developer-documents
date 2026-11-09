# Faz 8 Stock Integrations Migration Inventory

Bu envanter, `stock/integrations` kalıntılarının hangi sınırda kaldığını ve hangilerinin yeni `integrations` veya `telematics` operasyon yüzeyine taşındığını netleştirir.

## Sonuç

- Aktif stock UI içinde legacy `stock/integrations` ekranı yok.
- Kalan bağımlılık `API + permission + backend module` seviyesinde.
- `stock/locations` ownership sızıntısı kapatıldı; telematics artık stock modeline değil, kendi projection modeline yazar ve stock tarafına adapter ile düşer.
- `stock/integrations` şu an tam kaldırılabilir durumda değil; ürün içinde `compatibility only` olarak ele alınmalı.

## Envanter

| Alan | Bugünkü Yüzey | Karar | Hedef |
|---|---|---|---|
| Stock current/history locations | `/admin/stock/locations/current`, `/admin/stock/locations/history` | Stock içinde kalır | Telematics projection adapter ile beslenir |
| Vehicle assets / subscriptions / reports | `/admin/stock/vehicle-assets`, `/admin/stock/subscriptions`, `/admin/stock/reports/*` | Stock içinde kalır | Stock operasyon modülü |
| Legacy stock integrations connectors/jobs | `/admin/stock/integrations/*` API, `stock_integrations` permission/module | Uyumluluk için korunur | Yeni yönetişim `/admin/integrations` ve `/admin/telematics` |
| Command / callback / partner governance | `/admin/telematics/commands*` | Taşındı | Telematics command center ve partner governance |

## Kod Seviyesi Kanıt

- Frontend çağrı noktaları:
  - `lib/api/stock-ecosystem.ts` içinde legacy `listIntegrationConnectors`, `createIntegrationJob` vb. duruyor.
  - Aktif stock app route ağacında bu API metodlarını kullanan ekran bulunmuyor.
- Backend kalıntıları:
  - `internal/modules/stock/integrations/*`
  - `stock_integrations` permission ve role seed kayıtları
- Ownership cleanup:
  - `internal/modules/telematics/telemetry/model/model.go`
  - `internal/modules/telematics/telemetry/projection/stock_adapter.go`
  - `internal/modules/telematics/telemetry/service/service.go`

## Faz 8 İçin Sonraki Teknik Adımlar

1. `stock_integrations` permission ve menu görünürlüğünü runtime usage ölçerek daraltmak.
2. Legacy connector/job API’leri için deprecation response header veya audit event eklemek.
3. Partner/command ve genel integrations tarafında eşdeğer coverage tamamlandıktan sonra `stock/integrations` backend module’ünü read-only compatibility katmanına çekmek.
4. Son aşamada `stock-ecosystem.ts` içindeki deprecated integration client’larını ayrı bir legacy namespace’e taşımak.

## Finalizasyon

Faz 8 kapanış kararı ve retirement kriterleri ayrı dokümanda finalize edilmiştir:

- [FAZ8_CLOSURE_AND_RETIREMENT_DECISION.md](/workspace/aystek-core-dev-docs/plans/FAZ8_CLOSURE_AND_RETIREMENT_DECISION.md)
