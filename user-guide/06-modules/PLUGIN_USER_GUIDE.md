---
Title: Plugin User Guide
Version: 2.0.0
Status: Enterprise Certified
Last Updated: 2026-02-12
Related Documents:
  - ./MARKETPLACE_GUIDE.md
  - ../04-api/API_REFERENCE.md
  - ../01-architecture/Enterprise_Isolation_Whitepaper.md
---

# Table of Contents
1. [Plugin Runtime Principle](#1-plugin-runtime-principle)
2. [Lifecycle Overview](#2-lifecycle-overview)
3. [Authority Boundaries](#3-authority-boundaries)

## 1. Plugin Runtime Principle
Plugin execution authority is backend-side. UI consumes metadata and must not compute lifecycle or enforcement state.

## 2. Lifecycle Overview
```text
Installed -> Enabled -> Effective Enabled
```
Effective state is gated by capability, permissions, feature flags, and license constraints.

## 3. Authority Boundaries
- No plugin runtime authority in frontend.
- No seed/menu/translation ownership in plugin packages.
- No bypass of backend-enforced capability and context checks.
