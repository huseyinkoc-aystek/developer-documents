# System Ops Module Feature Roadmap

Date: 2026-04-05
Status: Planned
Owner: Codex

## Purpose

System Ops modülü closure aşamasını tamamladı.

Bu dokümanın amacı, bundan sonra açılacak işleri:
- `closure`
- `hardening`
- `stub completion`

olarak değil, doğrudan `feature roadmap` olarak çerçevelemektir.

## Roadmap Principles

Yeni System Ops işleri varsayılan olarak şu sınıflardan birine girmelidir:

1. Operator productivity features
2. Operational automation features
3. Governance and policy intelligence
4. External integration and provider-backed capabilities
5. Regression safety expansion

## Feature Tracks

### Track 1: Release Operations

Focus:
- release channel policy analytics
- update preflight summary and safety preview
- rollback recommendation engine
- update approval workflow
- release wave visibility

Candidate slices:
- update provider abstraction
- release artifact registry integration
- health gate evaluation summary
- rollback readiness score

### Track 2: Backup and Recovery Operations

Focus:
- restore workflow
- retention governance
- backup scope selection
- restore validation checklist
- recovery confidence score

Candidate slices:
- backup restore simulation
- restore request and approval flow
- retention policy scheduler
- backup storage target visibility

### Track 3: Maintenance Scheduling

Focus:
- scheduled maintenance windows
- expiry and auto-disable
- public impact preview
- operator note and audience targeting
- maintenance calendar visibility

Candidate slices:
- scheduled start/end persistence
- countdown and expiry UI
- public-facing maintenance payload preview
- maintenance recurrence templates

### Track 4: Migration Governance

Focus:
- migration batch detail
- failure classification
- rollbackability summary
- migration output log detail
- migration approval gate

Candidate slices:
- migration step log stream
- batch detail page
- reversible/non-reversible risk summary
- failure reason taxonomy

### Track 5: License and Entitlement Management

Focus:
- external license authority integration
- entitlement source visibility
- activation actor and history
- renewal workflow
- tenant/seat overage warnings

Candidate slices:
- signed entitlement validation
- remote activation adapter
- license event history
- entitlement delta reporting

### Track 6: Guardrails and Governance Intelligence

Focus:
- dual-control approvals
- action risk scoring
- operator-specific policy overlays
- policy change history
- policy drift alerts

Candidate slices:
- approval workflow layer
- guardrail violation summary
- policy audit timeline
- escalation notice chain

### Track 7: Verification Expansion

Focus:
- handler tests
- repository tests
- worker recovery tests
- API contract tests
- admin page smoke tests

Candidate slices:
- handler-level HTTP tests
- repository filter tests
- SSE/log stream tests
- frontend API normalization tests

## Prioritization Guidance

If a new System Ops initiative is proposed, prioritize in this order:

1. operator-visible operational value
2. safety/risk reduction
3. repeatability/automation
4. external integration depth
5. additional polish

## Recommended First Feature Roadmap Slice

The highest-value first post-closure feature slice is:

`Release Operations and Restore Readiness`

This combines:
- update provider abstraction
- rollback recommendation summary
- backup restore intent visibility
- operational approval hooks

Reason:
- it builds directly on the closure work already completed
- it improves real operator decisions rather than only extending internals
- it creates a clean bridge from System Ops closure into platform operations maturity

## Boundary

If future work is primarily about:
- fixing a regression
- restoring a broken contract
- re-aligning drift between Go and Next

that work should be treated as maintenance.

If future work introduces new operator capability, policy depth, automation, or integration, it belongs in this roadmap.
