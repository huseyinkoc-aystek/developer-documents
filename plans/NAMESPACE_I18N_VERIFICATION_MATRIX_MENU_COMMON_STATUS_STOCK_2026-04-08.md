# Namespace I18n Verification Matrix: `menu.*`, `common.*`, `status.*`, `stock.*`

Date: 2026-04-08

## Scope

This pass compares UI-used translation keys against the seed source for these namespaces:

- `menu.*`
- `common.*`
- `status.*`
- `stock.*`

Scanned UI surface:

- `/workspace/node/aystek-core-next`

Seed source:

- `/workspace/go/aystek-core/internal/seeds/translation_seed.go`

## Summary

- Total UI-used keys across scanned namespaces: `698`
- Total missing in seed: `154`

### By Namespace

- `menu.*`
  - UI-used: `49`
  - Missing: `3`

- `common.*`
  - UI-used: `126`
  - Missing: `51`

- `status.*`
  - UI-used: `4`
  - Missing: `0`

- `stock.*`
  - UI-used: `519`
  - Missing: `100`

## Missing Keys

### `menu.*`

- `menu.chatbot_canned`
- `menu.operations.messaging`
- `menu.operations.messaging.sms`

### `common.*`

- `common.action`
- `common.activate`
- `common.active`
- `common.add`
- `common.all_rules`
- `common.archive`
- `common.archived`
- `common.confirm_delete`
- `common.created`
- `common.created_at`
- `common.csv`
- `common.date`
- `common.deleteFailed`
- `common.deleted`
- `common.deleting`
- `common.details`
- `common.draft`
- `common.feature_coming_soon`
- `common.field`
- `common.id`
- `common.inactive`
- `common.lastUpdated`
- `common.last_run`
- `common.loadFailed`
- `common.manage`
- `common.message`
- `common.noResults`
- `common.no_data`
- `common.no_permission`
- `common.ok`
- `common.open`
- `common.order`
- `common.other`
- `common.pause`
- `common.payment`
- `common.play`
- `common.publish`
- `common.published`
- `common.quickLinks`
- `common.saveFailed`
- `common.saved`
- `common.selectOption`
- `common.send`
- `common.test`
- `common.tryAgain`
- `common.updateFailed`
- `common.updated_at`
- `common.upload`
- `common.version`
- `common.view_all`
- `common.xlsx`

### `status.*`

- No missing keys.

### `stock.*`

- `stock.dashboard.legacy_reports.alert_policy`
- `stock.dashboard.legacy_reports.alert_policy_cooldown`
- `stock.dashboard.legacy_reports.alert_policy_failed`
- `stock.dashboard.legacy_reports.alert_policy_meta`
- `stock.dashboard.legacy_reports.auto_report`
- `stock.dashboard.legacy_reports.auto_report_meta`
- `stock.dashboard.legacy_reports.done`
- `stock.dashboard.legacy_reports.eyebrow`
- `stock.dashboard.legacy_reports.failed`
- `stock.dashboard.legacy_reports.last_alert`
- `stock.dashboard.legacy_reports.last_auto_run`
- `stock.dashboard.legacy_reports.last_export`
- `stock.dashboard.legacy_reports.manage_policy`
- `stock.dashboard.legacy_reports.open_audit`
- `stock.dashboard.legacy_reports.open_jobs`
- `stock.dashboard.legacy_reports.open_stream`
- `stock.dashboard.legacy_reports.queued`
- `stock.dashboard.legacy_reports.quick_actions`
- `stock.dashboard.legacy_reports.title`
- `stock.dashboard.legacy_reports.total`
- `stock.dashboard.ops_policy_reports.auto_report_hint`
- `stock.dashboard.ops_policy_reports.auto_report_meta`
- `stock.dashboard.ops_policy_reports.auto_report_needed`
- `stock.dashboard.ops_policy_reports.auto_report_not_needed`
- `stock.dashboard.ops_policy_reports.auto_report_policy`
- `stock.dashboard.ops_policy_reports.auto_report_recommended`
- `stock.dashboard.ops_policy_reports.done`
- `stock.dashboard.ops_policy_reports.failed`
- `stock.dashboard.ops_policy_reports.last_auto_run`
- `stock.dashboard.ops_policy_reports.last_export`
- `stock.dashboard.ops_policy_reports.last_notice`
- `stock.dashboard.ops_policy_reports.manage_policy`
- `stock.dashboard.ops_policy_reports.notice_policy`
- `stock.dashboard.ops_policy_reports.notice_policy_cooldown`
- `stock.dashboard.ops_policy_reports.notice_policy_failed`
- `stock.dashboard.ops_policy_reports.notice_policy_meta`
- `stock.dashboard.ops_policy_reports.open_audit`
- `stock.dashboard.ops_policy_reports.open_jobs`
- `stock.dashboard.ops_policy_reports.open_stream`
- `stock.dashboard.ops_policy_reports.queued`
- `stock.dashboard.ops_policy_reports.title`
- `stock.dashboard.ops_policy_reports.total`
- `stock.portal.field_reports.assets.create`
- `stock.portal.field_reports.assets.recent`
- `stock.portal.field_reports.assets.subtitle`
- `stock.portal.field_reports.assets.title`
- `stock.portal.locations.current.subtitle`
- `stock.portal.locations.current.title`
- `stock.portal.notice.saved`
- `stock.portal.sales.create`
- `stock.portal.sales.subtitle`
- `stock.portal.sales.title`
- `stock.portal.subtitle`
- `stock.portal.title`
- `stock.reports.type.legacy_retirement`
- `stock.reports.type.lifecycle_regression_governance`
- `stock.reports.type.ops_policy_changes`
- `stock.reports.type.protocol_adapter_governance`
- `stock.sales.fields.checkout_session_help`
- `stock.sales.validation.product_item_required`
- `stock.settings.legacy_retirement.access_24h`
- `stock.settings.legacy_retirement.alert_cooldown`
- `stock.settings.legacy_retirement.alert_enabled`
- `stock.settings.legacy_retirement.alert_failed_threshold`
- `stock.settings.legacy_retirement.alert_policy`
- `stock.settings.legacy_retirement.alert_queue_threshold`
- `stock.settings.legacy_retirement.auto_disabled_at`
- `stock.settings.legacy_retirement.auto_report`
- `stock.settings.legacy_retirement.auto_report_enabled`
- `stock.settings.legacy_retirement.auto_report_format`
- `stock.settings.legacy_retirement.auto_report_interval`
- `stock.settings.legacy_retirement.blocked_24h`
- `stock.settings.legacy_retirement.description`
- `stock.settings.legacy_retirement.enabled`
- `stock.settings.legacy_retirement.gate_state`
- `stock.settings.legacy_retirement.interval`
- `stock.settings.legacy_retirement.last_access`
- `stock.settings.legacy_retirement.last_alert_at`
- `stock.settings.legacy_retirement.last_alert_state`
- `stock.settings.legacy_retirement.last_auto_report_run`
- `stock.settings.legacy_retirement.last_report_status`
- `stock.settings.legacy_retirement.live_snapshot`
- `stock.settings.legacy_retirement.retirement_ready`
- `stock.settings.legacy_retirement.title`
- `stock.settings.legacy_retirement.zero_access_days`
- `stock.settings.ops_policy.alert_cooldown`
- `stock.settings.ops_policy.alert_enabled`
- `stock.settings.ops_policy.alert_failed_threshold`
- `stock.settings.ops_policy.alert_policy`
- `stock.settings.ops_policy.alert_queue_threshold`
- `stock.settings.ops_policy.auto_report`
- `stock.settings.ops_policy.auto_report_enabled`
- `stock.settings.ops_policy.auto_report_format`
- `stock.settings.ops_policy.auto_report_interval`
- `stock.settings.ops_policy.last_alert_at`
- `stock.settings.ops_policy.last_alert_state`
- `stock.settings.ops_policy.last_auto_report_run`
- `stock.settings.ops_policy.last_report_status`
- `stock.settings.sales_auto_redirect.description.with_explicit_types`
- `stock.subscriptions.validation.pricing_required_for_date`

## Interpretation

- `status.*` is currently clean.
- `menu.*` has a small, fast-to-fix gap.
- `common.*` has a medium shared-foundation gap and should be normalized centrally.
- `stock.*` has the largest residual risk and should be handled slice-by-slice rather than patched ad hoc.

## Artifacts

- `/workspace/namespace_i18n_scan.json`
- `/workspace/ns_ui_keys.txt`
- `/workspace/ns_seed_keys.txt`
