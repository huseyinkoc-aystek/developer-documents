# Gitignore Release And Installer Audit Checkpoint

Date: 2026-04-29

## Scope

- `/workspace/go/aystek-core`
- `/workspace/node/aystek-core-next`
- `/workspace/node/aysteknoloji-com`
- `/workspace/go/aystek-core-installer`

## Findings

- `aystek-core` `.gitignore` was already comparatively mature and aligned with Go/runtime/release artifacts.
- `aystek-core-next` `.gitignore` was too weak for a production Next.js repo and risked drift around env examples, build outputs, logs, and editor files.
- `aysteknoloji-com` `.gitignore` was too weak for a production Next.js/public-site repo and did not clearly preserve example env files.
- `aystek-core-installer` `.gitignore` was too weak for a Go release/installer repo and did not cover host binary, dist outputs, coverage artifacts, and release archives.
- `aystek-core-installer` also had the built host binary tracked in Git, which is a release hygiene risk and can create false confidence around installer state.

## Applied Remediation

- Hardened `.gitignore` in:
  - `/workspace/node/aystek-core-next/.gitignore`
  - `/workspace/node/aysteknoloji-com/.gitignore`
  - `/workspace/go/aystek-core-installer/.gitignore`
- Preserved example env files explicitly where needed:
  - `.env.example`
  - `.env.local.example`
- Added ignores for:
  - Next.js build/cache outputs
  - coverage artifacts
  - release archives and checksums
  - logs/tmp scratch
  - OS/editor noise
- Removed tracked installer host binary from Git index:
  - `/workspace/go/aystek-core-installer/aystek-core-installer`

## Operational Meaning

- GHCR/release packaging is less likely to be polluted by local build outputs.
- Installer repo is less likely to ship misleading tracked binaries.
- Example env files remain versioned while secret-bearing local env files stay ignored.
- Next/public site repos now better match production packaging and CI realities.

## Validation

- `git check-ignore` used to verify:
  - examples are preserved
  - local env/build outputs are ignored
- `git status --ignored=matching` used to confirm expected ignore behavior
