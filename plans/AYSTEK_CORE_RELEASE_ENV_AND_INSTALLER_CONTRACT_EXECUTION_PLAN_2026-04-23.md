# AYSTEK Core Release, Env, And Installer Contract Execution Plan

Date: 2026-04-23
Status: In Progress
Scope: `aystek-core`, `aystek-core-next`, `aystek-core-installer`

## Objective

Align private `ghcr.io` release packaging, runtime environment contracts, and installer-generated configuration so that:
- `aystek-core` and `aystek-core-next` can be released as private GHCR images without ad-hoc env patching
- `aystek-core-installer` produces correct `core.env`, `panel.env`, and Docker Compose inputs
- custom host ports for MongoDB, Redis, and RabbitMQ remain valid and predictable in both native and Docker modes
- deployment docs, examples, and installer runtime behavior follow one canonical contract

## Key Findings

### 1. Frontend env contract is fragmented

Observed runtime/doc mismatch:
- Next README documents `NEXT_PUBLIC_API_BASE_URL`, `NEXT_PUBLIC_SSE_BASE_URL`, `NEXT_PUBLIC_APP_ENV`
- installer `panel.env.tmpl` currently emits only `NEXT_PUBLIC_API_URL`
- production compose injects only legacy `NEXT_PUBLIC_API_URL`
- Next Dockerfile only exposes build args for `NEXT_PUBLIC_API_URL` and `NEXT_PUBLIC_SITE_URL`

Impact:
- installer and compose can produce a working but incomplete panel runtime
- SSE/public API base resolution depends on fallback behavior instead of canonical values
- release images are harder to reason about because documented env keys and generated env keys diverge

### 2. Core/panel image naming is not fully unified

Observed:
- `aystek-core` production compose defaults to `ghcr.io/aystek/aystek-core` and `ghcr.io/aystek/aystek-panel`
- installer defaults use `ghcr.io/huseyinkoc-aystek/...`
- docs mention `ghcr.io/aystek/...`

Impact:
- release contract is ambiguous
- installer and manual compose users may pull from different registries/namespaces

### 3. Installer env templates are generous but not canonical

Observed:
- core template emits many Mongo aliases (`MONGO_URI`, `MONGO_DATABASE`, `MONGO_DB`, `MONGO_DBNAME`, `MONGO_DB_NAME`) which is acceptable for compatibility
- panel template does not emit the canonical Next env set documented by the frontend
- Docker env/compose templates use slightly different naming than production compose examples

Impact:
- generated config works by compatibility rather than by explicit contract
- future runtime behavior changes risk breaking installer assumptions

### 4. Port handling is mostly correct, but contract needs tightening

Observed:
- installer and production compose bind Mongo/Redis/RabbitMQ to `127.0.0.1`
- custom ports are supported in schema and templates
- naming differs between production compose (`RABBITMQ_PORT`, `RABBITMQ_MANAGEMENT_PORT`) and installer docker env (`AMQP_PORT`, `RABBIT_MGMT_PORT`)

Impact:
- behavior is correct, but variable naming drift increases operator confusion and doc debt

### 5. Local source deployment path is under-specified

Observed:
- native deploy already understands `core_binary_path` and `panel_build_path`
- docker deploy only understands image tags + registry auth
- installer wizard does not expose a first-class artifact source selection for fresh installs
- operator needs to be able to point at local `aystek-core` / `aystek-core-next` directories on the server instead of private GHCR

Impact:
- installer contract is incomplete for private/offline/self-hosted release workflows
- operators are pushed toward GHCR even when local source deployment is the intended operating model

## Execution Slices

### Slice 1: Canonical Env Contract
- define canonical panel env set:
  - `NEXT_PUBLIC_API_BASE_URL`
  - `NEXT_PUBLIC_API_URL` (compatibility)
  - `NEXT_PUBLIC_SSE_BASE_URL`
  - `NEXT_PUBLIC_APP_ENV`
  - `INTERNAL_API_URL`
- define canonical core env set for deployment/install docs around:
  - `PORT`
  - `MONGODB_URI` plus compatibility aliases
  - `REDIS_ADDR`
  - `REDIS_PASSWORD`
  - `RABBITMQ_URL`
  - `JWT_SECRET`
  - `ALLOWED_ORIGINS`
- update installer env templates and compose examples to emit the canonical set

### Slice 2: GHCR Release Contract Alignment
- unify image names and env variable names used by:
  - production compose
  - installer defaults
  - installer docs
- confirm installer registry login flow matches private GHCR release expectations
- ensure no old namespace remains as default

### Slice 3: Port And Service Binding Validation
- validate local/external Mongo/Redis/RabbitMQ contract in installer templates
- ensure custom host ports map correctly in Docker mode
- ensure native env templates point to the operator-selected host/port values
- document expected behavior clearly for local vs external service modes

### Slice 4: Documentation And Verification
- update:
  - `/workspace/aystek-core-dev-docs/README_FOR_AI_ASSISTANTS.md` implications already satisfied by synchronized docs requirement
  - installer usage guide
  - server setup guide
  - production deployment README/examples
  - frontend README if canonical env set changes
- run build/test verification across relevant repos

### Slice 5: Local Source Artifact Contract
- add first-class installer `source_mode`
- support:
  - `docker + registry`
  - `docker + local-source`
  - `native + local-artifacts`
  - `native + local-source`
- make deployers actually consume local source directories
- document toolchain/runtime expectations for local-source installs

### Slice 6: Private GHCR Release And Source-Mode Verification
- align release workflow to canonical GHCR namespace
- add explicit release checklist for private GHCR vs local-source operator choice
- add CI contract validation for:
  - `docker + registry`
  - `docker + local-source`
  - `native + local-artifacts`
  - `native + local-source`
- keep installer P0 disposable matrix separate from source-mode fixture validation

## First Concrete Fix Set
- write canonical env values into installer `panel.env.tmpl`
- align production compose with canonical panel env names
- align Next Dockerfile build args/env support with canonical names
- normalize default GHCR namespace/image names across installer + production compose

## Success Criteria
- installer-generated `panel.env` contains the same canonical public/runtime env keys documented by Next
- production compose and installer use the same default GHCR image coordinates
- custom Mongo/Redis/RabbitMQ ports are preserved correctly in generated config
- docs and examples no longer disagree on primary env keys
- verification passes with Go tests and targeted Next lint/build checks
