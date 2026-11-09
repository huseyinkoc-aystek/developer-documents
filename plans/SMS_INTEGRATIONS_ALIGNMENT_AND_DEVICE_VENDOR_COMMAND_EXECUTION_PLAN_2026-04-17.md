# SMS Integrations Alignment And Device Vendor Command Execution Plan

Date: 2026-04-17
Owner: Codex
Status: In Progress

## Objective

Deliver a coherent SMS command operating model for stock-bound SIM assets without unnecessary backend churn.

## Slice 1: Integrations Navigation Alignment

Goal:
- move stock SMS command entry under `Integrations`
- preserve backward compatibility for existing stock route consumers

Deliverables:
- integrations child menu entry for SMS commands
- legacy stock route redirect
- no backend package relocation

## Slice 2: Integrations-Native Admin Surface

Goal:
- align the stock SMS command page with the integrations design system

Deliverables:
- integrations header
- orchestration-aware summary cards
- SIM / identifier lookup helper
- command composer
- template and execution visibility

## Slice 3: SIM Target Resolution Contract

Goal:
- define how command targets are resolved from stock truth

Required rules:
- mounted HAT/SIM asset is preferred target
- `msisdn` is authoritative when available
- `imei` and vehicle identifier support lookup, not final delivery target by themselves
- unresolved targets must block send

## Slice 4: Device Vendor Command Planning

Goal:
- define enterprise contract for brand/device specific SMS-driven command delivery

Required outputs:
- supported command intents
- vendor resolution contract
- variable rendering rules
- vendor-specific command entries
- command-by-state matrix
- callback / delivery reconciliation
- audit, retry, dead-letter, and escalation expectations

## Slice 5: Governance And Operational Controls

Goal:
- avoid freeform unsafe command sending

Controls:
- role-based access
- tenant/site scoped visibility
- template approval states
- provider policy linkage
- command execution audit trail

## Not In This Slice

- premature physical relocation of stock SMS backend package
- speculative vendor implementation without confirmed vendor command contract
- duplicate orchestration modules

## Success Criteria

- operators find SMS commands under `Integrations`
- page visually matches integrations admin surfaces
- stock/SIM command flow is documented with real domain anchors
- device vendor command work is decomposed into implementable next slices, not hand-wavy backlog
