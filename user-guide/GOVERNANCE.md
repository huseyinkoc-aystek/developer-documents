---
Title: Documentation Governance
Version: 2.1.0
Status: Enterprise Certified
Last Updated: 2026-02-18
Related Documents:
  - ./README.md
  - ./CHANGELOG.md
  - ./01-architecture/Enterprise_Isolation_Whitepaper.md
  - ./01-architecture/CONTENT_OPERATIONS_ARCHITECTURE.md
---

# Table of Contents
1. [Governance Scope](#1-governance-scope)
2. [Mandatory Update Rule](#2-mandatory-update-rule)
3. [Versioning Rule](#3-versioning-rule)
4. [CI Enforcement](#4-ci-enforcement)
5. [Authority Rule](#5-authority-rule)
6. [Content Ops & Ingestion Documentation Rule](#6-content-ops--ingestion-documentation-rule)

## 1. Governance Scope
This policy governs all documentation under `user-guide/`.

## 2. Mandatory Update Rule
Any change impacting tenant/site isolation, context resolution, error model, security controls, moderation workflow, ingestion runtime, or site operational policies must include documentation updates in the same change set.

## 3. Versioning Rule
All enterprise documents follow `MAJOR.MINOR.PATCH`.
- Use `MAJOR` for structural and governance changes.
- Use `MINOR` for new approved content and large expansions.
- Use `PATCH` for corrections and consistency fixes.

## 4. CI Enforcement
CI must enforce:
1. Internal link validation across `user-guide/`.
2. No stale references to deprecated documentation paths.
3. Security gate execution (`make security-gates`) as part of release validation.
4. API reference consistency checks against runtime routes for critical admin surfaces.

## 5. Authority Rule
`01-architecture/Enterprise_Isolation_Whitepaper.md` is the authoritative narrative document for isolation model statements.
Supporting guides must reference, not duplicate, whitepaper-level explanations.

## 6. Content Ops & Ingestion Documentation Rule
For `content_ingestion` and `content_ops_settings` changes, updates must include:
- Architecture changes (bounded context and event flow)
- Security impact (RBAC, secret handling, injection controls)
- API contract changes (endpoint/request/response)
- Operational impact (monitoring and runbooks)
- Compliance impact (audit evidence and traceability)
