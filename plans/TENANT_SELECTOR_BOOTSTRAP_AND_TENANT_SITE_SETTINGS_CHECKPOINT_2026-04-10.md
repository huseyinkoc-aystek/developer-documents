# Tenant Selector Bootstrap And Tenant-Site Settings Checkpoint

Date: 2026-04-10

## Scope Closed

- `/tenant-selector` artık sadece seçim ekranı değildir.
- Admin/platform admin için:
  - dil değiştirme
  - tenant oluşturma
  - site oluşturma
  akışı tek yüzeyde açılmıştır.
- `tenants` koleksiyonu tenant selector için birincil kaynak haline getirilmiştir; tenant görünürlüğü artık yalnızca `stores` varlığına bağlı değildir.
- `Settings` altında ayrı tenant-site ayar yüzeyi açılmıştır:
  - `/admin/settings/tenant-site`
- yeni izin seed’i eklenmiştir:
  - `settings:tenant_site_settings`
- bu izin admin role force-sync zinciriyle admin’e atanır.

## Backend Changes

- `/me/context` artık tenant/site bootstrap capability bilgisi döner.
- yeni bootstrap endpoint’leri:
  - `GET /me/bootstrap-sites`
  - `POST /me/bootstrap-tenant`
  - `POST /me/bootstrap-site`
- tenant listesi `tenants` koleksiyonu + `stores` fallback ile çözülür.
- store response artık `tenant_id` döner.

## Frontend Changes

- `tenant-selector` ekranında:
  - locale switch
  - tenant create form
  - site create form
  - `setup` placeholder yönlendirmesi yerine gerçek create flow
- yeni settings yüzeyi:
  - selected tenant / selected site snapshot
  - site directory
  - site create
  - current site defaults update
  - tenant-site public branding update
- settings dashboard içine tenant-site surface link’i eklendi.

## I18n Coverage Added

- `tenantSelector.*` anahtarları seed’e eklendi.
- `settings.tenantSite.*` anahtarları seed’e eklendi.
- `menu.settings.tenant_site` anahtarı seed’e eklendi.

## Verification

- `cd /workspace/go/aystek-core && go test ./internal/modules/system/... ./internal/modules/stores/... ./internal/modules/settings/... ./internal/seeds ./cmd/api`
- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/tenant-selector/page.tsx' --file 'app/(dashboard)/admin/settings/tenant-site/page.tsx' --file 'app/(dashboard)/admin/settings/dashboard/page.tsx' --file 'lib/api/me.ts'`

## Next Correct Step

- `settings.*` namespace içinde kalan sessiz raw-key riskini sürdürmek:
  - settings tenant-site yeni yüzeyinin key kullanımını matrise dahil etmek
  - settings dashboard/general surfaces ile tenant-site slice’ını aynı verification setine bağlamak
