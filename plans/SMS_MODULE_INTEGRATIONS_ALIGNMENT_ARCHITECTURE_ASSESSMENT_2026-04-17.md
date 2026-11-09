# SMS Module Integrations Alignment Architecture Assessment

Date: 2026-04-17
Owner: Codex
Scope: `aystek-core`, `aystek-core-next`

## Summary

SMS capability is currently split across two domains:

1. `internal/modules/sms`
   - shared provider, credential, policy, outbox, DLR handling
   - integration-platform style capability
2. `internal/modules/stock/sms`
   - stock-device / SIM command templates and executions
   - tightly coupled to stock asset, tenant-site scope, and telematics-adjacent operations

This means a blind package move into `integrations` would mix two different responsibilities:

- integration infrastructure
- stock-device command orchestration

## Decision

### Frontend

Move the operator-facing stock SMS command surface under `Integrations` in admin navigation and route design.

Reason:
- operators mentally treat SIM command delivery as external gateway orchestration
- the page should sit next to connectors / hooks / runtime-style operations
- this is a UX and operational alignment improvement

### Backend

Do **not** physically move `internal/modules/stock/sms` into `internal/modules/integrations` yet.

Reason:
- it is stock-domain scoped, not generic
- it depends on tenant/site scoped stock context
- it persists stock-specific command template and execution data
- it acts on stock-linked SIM identifiers rather than abstract external integrations

## Proper Future Extraction Boundary

If we later introduce a generic orchestration layer, extraction should happen into a new abstraction such as:

- `internal/modules/integrations/sim_commands`
- or `internal/modules/integrations/device_command_orchestration`

Only then should device-vendor command adapters live in integrations-facing infrastructure while stock remains the source of truth for target resolution.

## Current Data Reality

Available domain anchors already exist:

- stock inventory items carry `msisdn`, `imei`, `operator`
- stock vehicle assets resolve mounted SIM/GPS relationships
- telematics lifecycle/provisioning already models SIM readiness

This is enough to plan an enterprise SIM command flow without doing a premature package move.

## Recommended Shape

### Keep

- `internal/modules/sms` as shared SMS infrastructure
- `internal/modules/stock/sms` as stock-scoped command domain

### Change

- move stock SMS command navigation under `Integrations`
- redesign the stock SMS command page to match integrations visual language
- keep legacy stock route as compatibility redirect

## Device Vendor / SIM Command Implication

Device command delivery should be modeled as:

1. target resolution
   - vehicle
   - mounted HAT/SIM asset
   - authoritative `msisdn`
2. device vendor adapter selection
   - `trio`
   - `mobiliz`
   - `arvento`
   - `n2_mobil`
   - future vendors
3. command intent resolution
   - based on lifecycle / provisioning / stock status
4. SMS template rendering and send
5. callback / reconciliation / audit

That is an orchestration problem, not a simple page relocation problem.

## Conclusion

The correct enterprise move now is:

- move **UI placement** to `Integrations`
- keep **stock command backend** in stock domain
- document and then implement device-vendor command orchestration as a distinct future slice
