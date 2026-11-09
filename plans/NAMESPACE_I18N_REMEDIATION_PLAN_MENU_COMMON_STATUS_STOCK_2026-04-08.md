# Namespace I18n Remediation Plan: `menu.*`, `common.*`, `status.*`, `stock.*`

Date: 2026-04-08

## Goal

Close the remaining raw-key risk outside the `settings.*` namespace by remediating the namespaces that still have seed gaps.

## Current State

- `status.*`: complete for scanned UI usage
- `menu.*`: small missing set
- `common.*`: broad shared-key deficit
- `stock.*`: large domain-specific deficit concentrated in dashboard/report/policy and portal flows

## Execution Slices

### Slice 1: `menu.*` Quick Closure

Target:

- `menu.chatbot_canned`
- `menu.operations.messaging`
- `menu.operations.messaging.sms`

Outcome:

- eliminate the remaining menu raw-key risk for scanned surfaces

### Slice 2: `common.*` Foundation Closure

Target:

- add the missing shared action/status/table/file/runtime keys

Priority groups:

- action verbs
- lifecycle verbs
- table/date/file keys
- empty/error/retry helpers

Outcome:

- reduce repeated fallback literals across modules
- close shared raw-key risk at the base namespace

### Slice 3: `stock.*` Dashboard/Reports Closure

Target:

- `stock.dashboard.legacy_reports.*`
- `stock.dashboard.ops_policy_reports.*`
- `stock.reports.type.*`

Outcome:

- close the highest-visibility stock governance/reporting raw-key risk

### Slice 4: `stock.*` Settings Closure

Target:

- `stock.settings.legacy_retirement.*`
- `stock.settings.ops_policy.*`
- `stock.settings.sales_auto_redirect.description.with_explicit_types`

Outcome:

- remove raw-key exposure from stock settings and policy surfaces

### Slice 5: `stock.*` Portal/Sales/Subscription Closure

Target:

- `stock.portal.*`
- `stock.sales.*` residuals
- `stock.subscriptions.validation.pricing_required_for_date`

Outcome:

- close stock operational and portal-facing residual gaps

## Verification

For each slice:

- rerun namespace scan
- confirm missing count reduced for the targeted namespace group
- lint touched Next files
- run seed tests

## Closure Condition

This remediation track is complete when:

- scanned UI-used `menu.*`, `common.*`, `status.*`, `stock.*` keys all exist in seed
- no raw-key regressions remain in the scanned surfaces
