# Public Site CMS Parity And Default Tenant Bootstrap Checkpoint

**Date:** 2026-04-10  
**Status:** In Progress / Verified Checkpoint  
**Scope:** `/workspace/node/aysteknoloji-com` + `/workspace/go/aystek-core`

---

## Scope Closed In This Checkpoint

- `SystemNotInitialized` kök nedeni kapatıldı.
- `default` tenant için otomatik `admin-panel` site seed’i devreye alındı.
- soft-delete olmuş store kayıtlarının seed ile restore edilmesi sağlandı.
- `aystek` tenant public CMS seed’lerinde `site_id` bağlamı düzeltildi.
- `hakkimizda`, `cozumler`, `aystek-core`, `iletisim` sayfaları CMS-first SSR wrapper yapısına taşındı.
- blog / SEO / OG metadata hattındaki statik kalıntılar temizlendi.

## Implemented Changes

### Backend
- `internal/modules/system/service/bootstrap_service.go`
  - `TenantCount()` erken `0` dönüşü kaldırıldı; stores fallback yeniden aktif.
- `internal/seeds/store_tax_seed.go`
  - `default` tenant için `admin-panel` site eklendi.
  - soft-delete restore için `$unset: { deleted_at: "" }` akışı eklendi.
- `internal/seeds/util.go`
  - `resolveSeedSiteID()` ve tenant-site seed context yardımcıları kullanıma alındı.
- `internal/seeds/aystek_cms_seed.go`
  - slider/page seed akışlarına `site_id` enjekte edildi.

### Frontend / Public Site
- `lib/cms/server-page.ts`
  - cached CMS page resolver + metadata helper sertleştirildi.
- `lib/cms/server-post.ts`
  - OG alanları ayrı geçirilir hale getirildi.
- `lib/seo/site.ts`
  - `ogTitle` / `ogDescription` ayrımı eklendi.
- `app/hakkimizda/page.tsx`
- `app/cozumler/page.tsx`
- `app/aystek-core/page.tsx`
- `app/iletisim/page.tsx`
  - tamamı server wrapper + `generateMetadata()` pattern’ine geçti.

## Verified Runtime Evidence

### 1) Health Check
```bash
curl -sf http://localhost:9090/healthz
```
Response:
```json
{"success":true,"data":{"status":"ok"}}
```

### 2) Auth Login
```bash
curl -sf -X POST http://localhost:9090/auth/login-by-username \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: default" \
  -d '{"identifier":"huseyinkoc","password":"L5834!"}'
```
Result: access token + permissions payload returned successfully.

### 3) Store Verification
```bash
curl -sf -H "X-Tenant-ID: default" -H "Authorization: Bearer <token>" \
  http://localhost:9090/admin/stores
```
Verified result:
- `admin-panel`
- `default`

Both stores are visible for the `default` tenant.

## Remaining Exact Next Steps

1. **Homepage CMS-first preload closure**
   - server-side home hero / slider / featured content preload
   - homepage metadata from CMS SEO fields

2. **Language resolution hardening**
   - `cookie + Accept-Language` resolver
   - `tr / en / es` first-render verification

3. **Media + canonical hardening**
   - `featured_image`, `cover_image`, `og_image` normalization
   - canonical / sitemap / robots verification

4. **Public API smoke matrix**
   - `/public/posts`
   - `/public/pages`
   - `/public/sliders`
   - `/public/navmenus`

## Operational Note

Bu servis için aktif route prefix’leri şunlardır:
- `/auth/...`
- `/admin/...`
- `/public/...`

Frontend tarafında hâlâ `/api/v1/...` varsayımı olan yerler varsa ayrıca hizalanmalıdır.
