# API Sharing Architecture Review and Remediation Plan

**Review Date:** 2026-03-11  
**Scope:** `/workspace/aystek-core-dev-docs/plans/AYSTEK CORE Go Mikroservislerinde Güvenli, İzlenebilir ve Tenant‑Scoped API Paylaşımı.md`, `/workspace/aystek-core-dev-docs/plans/api-sharing-implementation/*`, `/workspace/go/aystek-core/internal/modules/integrations/api_sharing`

## Executive Verdict

Yeni API-sharing yönü stratejik olarak doğru: dışa açılacak AYSTEK Core route'larını ayrı bir control-plane + policy enforcement katmanıyla tenant-scoped, audit edilebilir ve güvenli hale getirme hedefi AYSTEK Core'a yakışan yöndür.

Ancak mevcut implementation paketi bugün için **enterprise-ready değildir**. Kod tabanı "çöp" değildir; aksine repo/service/handler/test iskeleti ve bazı yararlı building block'lar üretilmiştir. Fakat kritik güvenlik ve yetkilendirme boşlukları nedeniyle sistem en fazla **internal prototype / controlled pilot** seviyesinde kabul edilmelidir.

## What Is Correct and Worth Keeping

- External gateway / edge ayrıştırması doğru yönde kurgulanmış.
- Operation, scope, grant, consent, audit, OAuth/JWKS, OpenAPI publishing için domain iskeleti değerli.
- `api_sharing` modülünün ayrı module olarak açılması AYSTEK Core modüler yapısına uygundur.
- Benchmark ve k6 varlıkları, production iddiası için yetersiz olsa da tekrar kullanılabilir.
- KVKK, audit trail, retention, DPA, OpenAPI ve RFC checklist dokümanları tam çöpe atılacak içerikler değildir; fakat önemli ölçüde gerçek duruma göre revize edilmelidir.

## Findings

### 1. Critical: token issuance is not constrained by the grant, and runtime authorization is not operation-based

`IssueClientCredentialsToken` çağrısı istemcinin istediği scope listesini doğrudan JWT içine yazıyor; grant ile karşılaştırmıyor. Runtime tarafında `PolicyEnforcementMiddleware` route/operation gereksinimini çözmeden yalnızca token'dan okuduğu scope'ların grant içinde olup olmadığını kontrol ediyor. Sonuç: `allowed_operations` fiilen zorlanmıyor ve operasyon bazlı yetkilendirme yok.

References:
- [service.go](/workspace/go/aystek-core/internal/modules/integrations/api_sharing/service/service.go#L1746)
- [policy_enforcement.go](/workspace/go/aystek-core/internal/modules/integrations/api_sharing/middleware/policy_enforcement.go#L50)
- [policy_enforcement.go](/workspace/go/aystek-core/internal/modules/integrations/api_sharing/middleware/policy_enforcement.go#L156)

Impact:
- Tasarlanan "tenant-scoped API publishing" hedefi tam olarak sağlanmıyor.
- Scope kontrolü operation-centric değil, caller-centric kalıyor.
- `allowed_operations` alanı tasarımda var ama enforcement zincirinde yok.

Required remediation:
- Token issue sırasında requested scope seti grant tarafından daraltılmalı.
- Gateway request'i `operation_id` veya canonical route id'ye çözümlemeli.
- Authorization kararı `grant + operation.required_scopes + allowed_operations + tenant/site scope` ile verilmelidir.

### 2. Critical: application middleware trusts unverified JWT claims

External middleware JWT'yi `ParseUnverified` ile açıyor ve claim'leri güvenilir kabul ediyor. Doküman "validation done via gateway" diyor; fakat uygulama içinde de `/external/v1` route'ları mount edilmiş durumda. Enterprise sistemde gateway bypass, yanlış ingress kuralı veya internal exposure varsayımına güvenilmez.

References:
- [policy_enforcement.go](/workspace/go/aystek-core/internal/modules/integrations/api_sharing/middleware/policy_enforcement.go#L50)
- [routes.go](/workspace/go/aystek-core/internal/modules/integrations/api_sharing/routes.go#L149)

Impact:
- Defense-in-depth yok.
- Yanlış ağ konfigürasyonunda forged JWT claim'leri policy kararını etkileyebilir.

Required remediation:
- App katmanında JWKS veya introspection ile signature/exp/iss/aud doğrulaması zorunlu olmalı.
- `/external/v1` yalnızca gateway arkasından erişilebilir hale getirilmeli; network policy ve trusted proxy sınırları açık tanımlanmalı.

### 3. High: authorization-code flow does not validate redirect URI against a registered client

`InitiateAuthCode` gelen `redirectURI` değerini olduğu gibi kaydediyor. Exchange aşamasında yalnızca ilk girilen değerle eşitlik kontrol ediliyor; registered redirect URI whitelist yok. Bu, gerçek OAuth client registry olmadan authorization-code flow'un güvenli olmadığı anlamına gelir.

References:
- [service.go](/workspace/go/aystek-core/internal/modules/integrations/api_sharing/service/service.go#L1791)
- [service.go](/workspace/go/aystek-core/internal/modules/integrations/api_sharing/service/service.go#L1851)

Impact:
- OAuth authorization code flow production için uygun değil.
- Redirect hijack / yanlış client onboarding riski var.

Required remediation:
- Ayrı `api_clients` kaydı ve registered redirect URI listesi oluşturulmalı.
- `authorize` ve `token` akışları registered client metadata ile doğrulanmalı.
- Public/confidential client ayrımı zorlanmalı.

### 4. High: client authentication silently succeeds when no secret is stored

`AuthenticateClient` içinde `HashedClientSecret` boşsa çağrı başarılı kabul ediliyor. Bu backward-compatibility kolaylığı olabilir; fakat external API platformunda bu kabul edilemez varsayılandır.

Reference:
- [service.go](/workspace/go/aystek-core/internal/modules/integrations/api_sharing/service/service.go#L2249)

Impact:
- `client_id + tenant_id` bilen taraf token alabilir hale gelir.
- Confidential client modeli fiilen delinmiş olur.

Required remediation:
- Secret'siz grant ile OAuth token issuance kapatılmalı.
- Public client senaryosu gerekiyorsa ayrı client type ve PKCE-only akışına taşınmalı.

### 5. High: usage metering schema cannot safely store per-operation aggregates

`api_usage` unique index'i `{tenant_id, client_id, period}` üzerinde; fakat repository upsert filtreleri `{tenant_id, client_id, operation_id, period}` kullanıyor. Bir tenant/client altında ikinci operation geldiğinde unique index çakışması veya aggregation bozulması oluşur. Bu, quota/billing/observability verisini güvenilmez yapar.

References:
- [indexes.go](/workspace/go/aystek-core/internal/modules/integrations/api_sharing/indexes.go#L130)
- [mongo.go](/workspace/go/aystek-core/internal/modules/integrations/api_sharing/repository/mongo.go#L1170)

Impact:
- Per-operation usage, billing ve quota raporu doğru tutulamaz.
- Enterprise raporlama ve abuse detection zayıflar.

Required remediation:
- Unique index `{tenant_id, client_id, operation_id, period}` olacak şekilde düzeltilmeli.
- Migration planı yazılmalı.

### 6. Medium: implementation status documents overstate readiness

Implementation docs future sprint'leri olmuş gibi yazıyor, yanlış repo path'leri (`aystek_coreX`) içeriyor ve staging kanıtı yokken completion tonu taşıyor. Bu, delivery ve governance açısından risklidir.

References:
- [SPRINT_11_LOAD_TEST_PERFORMANCE_REPORT.md](/workspace/aystek-core-dev-docs/plans/api-sharing-implementation/SPRINT_11_LOAD_TEST_PERFORMANCE_REPORT.md#L3)
- [SPRINT_11_LOAD_TEST_PERFORMANCE_REPORT.md](/workspace/aystek-core-dev-docs/plans/api-sharing-implementation/SPRINT_11_LOAD_TEST_PERFORMANCE_REPORT.md#L19)
- [INDEX.md](/workspace/aystek-core-dev-docs/plans/api-sharing-implementation/INDEX.md#L1)

Impact:
- Yönetim katmanında yanlış ilerleme algısı oluşur.
- "Production readiness" kararı hatalı veriyle alınabilir.

Required remediation:
- Tüm sprint belgeleri "planned / prototype / verified" şeklinde yeniden etiketlenmeli.
- Gerçek doğrulama kanıtı olmayan tamamlandı iddiaları kaldırılmalı.

## Useful Work vs Rework / Discard

### Keep and continue

- Model, DTO, repository ve handler iskeletleri
- JWT/JWKS/refresh-token temel yapısı
- OpenAPI generation başlangıcı
- Load-test script ve artifact üretim zinciri
- KVKK/DPA/audit kavram seti

### Rework before external exposure

- Policy enforcement zinciri
- Scope/operation authorization modeli
- OAuth client registry ve redirect-uri yönetimi
- Usage metering ve quota veri modeli
- Sprint durum belgeleri ve readiness iddiaları

### Do not treat as completed

- Production-ready OAuth2/OIDC server
- Enterprise-grade external gateway authorization
- KVKK-complete retention and deletion posture
- 1000 RPS verified go-live readiness

## Remediation Work Plan

### Phase 0: Stop incorrect go-live assumptions

- API-sharing status resmi olarak `prototype` olarak işaretlenecek.
- Dış partner onboarding durdurulacak; yalnızca internal sandbox kullanımı açık kalacak.
- Sprint 11/12 belgeleri "evidence pending" ve "not approved for production" biçiminde revize edilecek.

### Phase 1: Security Closure

- Token issuance grant-bounded hale getirilecek.
- App-side JWT validation eklenecek.
- Secret'siz client authentication kapatılacak.
- Redirect URI whitelist ve client registry tasarlanacak.
- Threat model ve abuse cases yeniden yazılacak.

### Phase 2: Enterprise Authorization Model

- Route -> operation_id çözümleme standardı eklenecek.
- `allowed_operations`, `required_scopes`, `tenant_id`, `site_id` birlikte enforce edilecek.
- External edge contract netleştirilecek: gateway responsibility vs app responsibility.
- Usage/quota index ve aggregation modeli düzeltilecek.

### Phase 3: Compliance, Operations and Real Readiness

- KVKK retention, export, deletion ve audit kapsamı veri envanteri ile tekrar bağlanacak.
- Staging evidence üretilecek: gateway + Redis + Mongo + OTel dahil.
- Production readiness checklist gerçek kanıtla doldurulacak.

## Delivery Backlog

1. [x] P0 security fixes for token issuance, JWT validation and secret enforcement
2. [x] Client registry + redirect URI whitelist + client type model
3. [x] Operation resolution and allowed-operations enforcement
4. [x] Usage index migration and quota correctness tests
5. [x] Admin CRUD surface for `oauth2_clients` and controlled onboarding metadata
6. [x] Onboarding workflow: approval state machine, secret rotation, sandbox activation/revocation runbook
7. [x] Sprint docs normalization and real readiness evidence collection

## Acceptance Bar for "Enterprise / AYSTEK Core Grade"

- No unsigned/unverified token trust anywhere in the request path
- No client receives scopes beyond its active grant
- Every external request is authorized against a concrete `operation_id`
- Tenant/site scope is enforced in both edge and downstream service layer
- Audit, quota and usage data are internally consistent and queryable
- Staging and pre-production evidence exists for performance and security gates

## Current Status Statement

Bu çalışma yön olarak doğru, fakat bugünkü haliyle "enterprise ve AYSTEK Core'a yakışır bitmiş çözüm" değildir. En doğru tanım:

**`Directionally correct, partially useful, but still a prototype that requires targeted remediation before any serious external exposure.`**
