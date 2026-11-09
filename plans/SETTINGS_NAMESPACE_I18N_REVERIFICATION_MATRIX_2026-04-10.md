# Settings Namespace I18n Reverification Matrix

Date: 2026-04-10

## Scope

This reverification pass includes:

- existing `settings.*` admin surfaces
- newly added `/admin/settings/tenant-site`
- settings dashboard and general settings surfaces

Validation rule:

- only real translation calls were counted
- source of truth pattern: `t("settings.*")`
- plain data-field access like `settings.logo_url` was excluded from the missing-key calculation

## Result

- UI-used `settings.*` i18n keys: `186`
- Seeded `settings.*` translation keys present for those uses: `186`
- Missing seeded key count: `0`

## Conclusion

`settings.*` namespace is now closed at the translation-key level for actual Next admin i18n usage.

This includes the new tenant-site settings surface.

## Residual Risk

If raw keys still appear on a running environment, the remaining likely causes are operational, not seed-definition gaps:

1. seed reconcile has not yet run against the active database
2. frontend translation hydration loaded stale cache before updated seed data became available
3. runtime is serving an older backend/frontend build pair

## Evidence

Artifacts written during scan:

- `/workspace/settings_namespace_scan_files.json`
- `/workspace/settings_namespace_i18n_only_scan.json`

Scan result:

- `ui_count = 186`
- `missing_count = 0`

## Next Correct Step

Move from `settings.*` closure to runtime verification:

- seed reconcile on active environment
- hard refresh / cache clear
- spot-check settings dashboard, general settings, tenant-site settings, and shipping settings for raw-key regressions
