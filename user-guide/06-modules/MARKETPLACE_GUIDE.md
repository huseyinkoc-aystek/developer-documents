---
Title: Marketplace Guide
Version: 2.1.0
Status: Enterprise Certified
Last Updated: 2026-02-18
Related Documents:
  - ./PLUGIN_USER_GUIDE.md
  - ./CONTENT_INGESTION_AND_OPS_MODULES.md
  - ../04-api/API_REFERENCE.md
---

# Table of Contents
1. [Scope](#1-scope)
2. [Catalog Statuses](#2-catalog-statuses)
3. [Contract Boundaries](#3-contract-boundaries)
4. [Provider/Connector Marketplace Alignment](#4-providerconnector-marketplace-alignment)

## 1. Scope
Marketplace is a read-only plugin catalog for metadata, verification state, compatibility, and licensing context.

## 2. Catalog Statuses
- `verified`
- `rejected`
- `deprecated`
- `archived`

## 3. Contract Boundaries
- Marketplace metadata is not runtime authority.
- Plugin runtime state remains backend-governed via plugin infrastructure.

## 4. Provider/Connector Marketplace Alignment
For content ingestion ecosystems, marketplace principles apply to provider/connectors as follows:
- Connector templates define capability metadata, not tenant secrets.
- Tenant/site-specific operational config remains in scoped runtime modules.
- Commercial terms and provider status are governed by RBAC + audit + policy controls.

For implementation details, refer to `./CONTENT_INGESTION_AND_OPS_MODULES.md`.
