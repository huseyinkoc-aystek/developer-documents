# Env Example Contract Alignment Checkpoint

Date: 2026-04-29

## Scope

- `/workspace/go/aystek-core/.env.example`
- `/workspace/go/aystek-core/config/env/core.env.example`
- `/workspace/go/aystek-core/config/env/panel.env.example`
- `/workspace/node/aystek-core-next/.env.example`
- `/workspace/node/aysteknoloji-com/.env.example`
- `/workspace/node/aysteknoloji-com/.env.local.example`

## Findings

- `aystek-core` root `.env.example` was stale and mixed obsolete keys with keys not used by the current runtime contract.
- `aystek-core` runtime examples were missing some compatibility aliases that the code still resolves (`PORT`, `MONGO_URI`, `MONGO_DB*`, `INTERNAL_API_TOKEN`).
- SMTP variable names in examples were inconsistent with the real mail config reader (`SMTP_USER`, `SMTP_PASS`, `SMTP_FROM_EMAIL`, `SMTP_FROM_NAME`, `SMTP_TLS`).
- `aystek-core-next` example was missing `NEXT_API_URL` and did not document optional pinned tenant/site bootstrap variables.
- `aysteknoloji-com` examples were incomplete for SSR/public fetch configuration and did not document optional site code explicitly.

## Applied Remediation

- Rewrote `/workspace/go/aystek-core/.env.example` to match the current canonical local/developer contract.
- Expanded `/workspace/go/aystek-core/config/env/core.env.example` with compatibility aliases and SMTP/internal token fields.
- Expanded `/workspace/go/aystek-core/config/env/panel.env.example` with `NEXT_API_URL` and optional fixed context variables.
- Expanded `/workspace/node/aystek-core-next/.env.example` with `NEXT_API_URL` and optional fixed context variables.
- Expanded `/workspace/node/aysteknoloji-com/.env.example` and `.env.local.example` to include the full SSR/public-site runtime contract.

## Validation

- Searched real env usage from code (`process.env.*`, `os.Getenv`, config resolvers).
- Re-read all example files after patching.
- Ran:
  - `cd /workspace/go/aystek-core && go test ./cmd/api ./internal/platform/config ./internal/seeds`

