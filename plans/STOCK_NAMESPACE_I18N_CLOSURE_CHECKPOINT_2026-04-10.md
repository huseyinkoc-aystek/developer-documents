# Stock Namespace I18n Closure Checkpoint

Date: 2026-04-10

## Scope

This checkpoint closes the `stock.*` namespace seed gap identified in:

- `NAMESPACE_I18N_VERIFICATION_MATRIX_MENU_COMMON_STATUS_STOCK_2026-04-08.md`
- `NAMESPACE_I18N_REMEDIATION_PLAN_MENU_COMMON_STATUS_STOCK_2026-04-08.md`

## Baseline

- `stock.*` UI-used keys: `519`
- missing seed keys before remediation: `100`

## Remediation

Added the missing `stock.*` translation entries to:

- `/workspace/go/aystek-core/internal/seeds/translation_seed.go`

Closed groups:

- `stock.dashboard.legacy_reports.*`
- `stock.dashboard.ops_policy_reports.*`
- `stock.portal.*`
- `stock.reports.type.*`
- `stock.sales.fields.checkout_session_help`
- `stock.sales.validation.product_item_required`
- `stock.settings.legacy_retirement.*`
- `stock.settings.ops_policy.*`
- `stock.settings.sales_auto_redirect.description.with_explicit_types`
- `stock.subscriptions.validation.pricing_required_for_date`

## Verification

- residual missing `stock.*` keys after remediation: `0`
- command:
  - `cd /workspace/go/aystek-core && go test ./internal/seeds ./cmd/api`

## Outcome

- `stock.*` namespace raw-key risk is closed at the seed layer for the scanned UI usage.
- Any remaining raw-key issue on a live screen after this slice should be treated as:
  - seed reconcile not yet executed on the active database, or
  - a runtime translation loading/hydration issue rather than a missing seed entry.
