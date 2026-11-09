---
Title: Documentation Governance Rules
Version: 2.0.0
Status: Approved
Last Updated: 2026-02-12
Related Documents:
  - ../README.md
  - ../CHANGELOG.md
  - ../01-architecture/Enterprise_Isolation_Whitepaper.md

---
# Table of Contents
1. [Governance Rules](#governance-rules)
2. [CI Expectations](#ci-expectations)
3. [Terminology Standard](#terminology-standard)

## Governance Rules
1. Any isolation-impacting code change must include corresponding document updates.
2. No merge is approved with broken references in `user-guide/`.
3. Historical reports are archived, not deleted.
4. Bilingual documents must clearly separate Turkish and English sections.

## CI Expectations
1. CI must verify reference integrity for local markdown links in `user-guide/`.
2. CI must fail when required enterprise references are missing from index documents.
3. CI should run `security-gates` in code repository for contract alignment.

## Terminology Standard
Use these terms consistently:
- `tenant/site`
- `isolation`
- `dominance`
- `tombstone/tombstoned`
- `deterministic outcome`
