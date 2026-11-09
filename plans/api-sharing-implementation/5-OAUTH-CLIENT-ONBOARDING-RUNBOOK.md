# OAuth Client Onboarding Runbook

**Status:** implemented in backend as of 2026-03-11, still `sandbox/internal` operational scope only

## Scope

Bu runbook, `oauth2_clients` onboarding lifecycle'ını admin control-plane üzerinden nasıl yöneteceğimizi tanımlar.

- Secret lifecycle ve rotation
- Approval / rejection state machine
- Sandbox activation / revocation
- Admin API surface
- Seed/menu/permission matrix

## Current Backend State

Backend tarafında aşağıdaki yüzeyler mevcut:

- `POST /admin/api-sharing/oauth-clients`
- `GET /admin/api-sharing/oauth-clients`
- `GET /admin/api-sharing/oauth-clients/:id`
- `PUT /admin/api-sharing/oauth-clients/:id`
- `POST /admin/api-sharing/oauth-clients/:id/submit-approval`
- `POST /admin/api-sharing/oauth-clients/:id/approve`
- `POST /admin/api-sharing/oauth-clients/:id/reject`
- `POST /admin/api-sharing/oauth-clients/:id/rotate-secret`
- `POST /admin/api-sharing/oauth-clients/:id/sandbox/activate`
- `POST /admin/api-sharing/oauth-clients/:id/sandbox/revoke`

## State Machine

### Onboarding Status

- `draft`
- `pending_approval`
- `approved`
- `rejected`

Allowed transitions:

1. `draft -> pending_approval`
2. `rejected -> pending_approval`
3. `pending_approval -> approved`
4. `pending_approval -> rejected`

### Runtime Status

- `inactive`
- `active`
- `suspended`
- `revoked`

### Exposure Level

- `internal_only`
- `sandbox`
- `production`

Current enforcement rule:

- Token issuance için client şu koşulları sağlamalı:
  - `status=active`
  - `onboarding_status=approved`
  - `exposure_level in (sandbox, production)`

## Secret Lifecycle

### Confidential Clients

- Secret admin endpoint ile üretilir.
- Plain secret sadece rotation response'unda bir kez gösterilir.
- DB'de yalnızca bcrypt hash tutulur.
- Yeni secret üretildiğinde önceki aktif secret `grace` durumuna alınır.
- Grace window şu an 24 saattir.
- Opsiyonel `secret_ttl_days` verilebilir.

### Public Clients

- Secret rotation yasaktır.
- PKCE zorunludur.

## Operational Procedure

### 1. Create

Client oluştur:

- `status=inactive`
- `onboarding_status=draft`
- `exposure_level=internal_only`

### 2. Submit For Approval

Security/Platform review öncesi:

- redirect URI whitelist
- allowed scopes
- site scope
- partner/legal metadata

Sonra `submit-approval` çağrılır.

### 3. Approve

Approval sonrası client:

- `onboarding_status=approved`
- hala `inactive` ve `internal_only` kalabilir

Bu aşama production/sandbox erişimi açmaz.

### 4. Rotate Secret

Sadece confidential client için:

- `rotate-secret`
- secret güvenli kanal üzerinden partnere teslim edilir
- secret ticket/secret manager dışında tekrar paylaşılmaz

### 5. Sandbox Activate

Aktivasyon ön koşulları:

- `onboarding_status=approved`
- confidential ise aktif secret mevcut
- grant ve scope seti hazır

Aktivasyon sonucu:

- `status=active`
- `exposure_level=sandbox`

### 6. Sandbox Revoke

Revocation durumları:

- abuse / compromise
- legal request
- expired pilot
- partner offboarding

Sonuç:

- `status=inactive`
- `exposure_level=internal_only`

Gerekirse ardından secret rotation ve grant revoke yapılır.

## Permission Matrix

Bu yüzey mevcut `integrations` modül permission setini kullanır.

- `integrations:read`
  - list/get oauth clients
  - list grants
  - read audit views
- `integrations:create`
  - create oauth client
  - create grants
- `integrations:update`
  - update oauth client
  - rotate secret
- `integrations:approve`
  - submit for approval
  - approve
  - reject
- `integrations:activate`
  - sandbox activate
- `integrations:deactivate`
  - sandbox revoke

Not:

- `internal/seeds/permission_seed.go` içinde `integrations` modülü için `approve`, `activate`, `deactivate` action'ları zaten seed ediliyor.

## Seed / Menu Surface

Eklenen seed tarafları:

- `internal/seeds/api_sharing_seed.go`
  - örnek public sandbox client
  - örnek pending confidential client
- `internal/seeds/menu_seed.go`
  - API Sharing parent menu
  - OAuth Clients
  - Tenant Grants
  - Access Audit
- `internal/seeds/translation_seed.go`
  - onboarding action notices
  - menu title keys

## Admin UI Note

Bu workspace içinde `aystek-core-next` admin UI repo yok. Bu yüzden bu turda:

- backend API contract
- menu/translation seeds
- onboarding state model

tamamlandı.

Gerçek Next.js ekranları ayrı frontend repo/workspace içinde açılmalı.

## Recommended Next Step

1. Frontend admin screens for onboarding timeline, secret rotation modal and audit feed
2. Grant + OAuth client linking UX
3. Secret delivery integration with Vault/KMS-backed channel
4. Offboarding automation: revoke grant + revoke sandbox + rotate secret + audit export
