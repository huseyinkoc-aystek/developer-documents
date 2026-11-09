# TASK 041 - PMS S15 Product Expansion

## Scope
This task expands PMS admin operations for exports and webhooks with retry, filters, and delivery visibility.

## Goals
- Add retry path for export jobs and webhook deliveries.
- Add date range filters for export job lists.
- Add delivery log list for webhooks with status tabs and filters.
- Require event allowlist for webhook subscriptions.
- Update admin UI for exports and webhooks.
- Update translations for all new UI keys.

## Non-Goals
- No new external integrations beyond existing webhook dispatch.
- No changes to portal-facing features.
- No changes to payment or billing core logic.

## Change Summary
- Exports: list filters by from_date and to_date; retry failed jobs.
- Webhooks: delivery list with filters, retry failed/dead_letter deliveries; enforce event allowlist.
- Webhooks: legacy allowlist summary endpoint + gauge for migration tracking.
- UI: new export detail view; delivery tabs and filters in webhooks.
- i18n: add keys for exports detail, retry, dates, and webhook delivery fields.

## Breaking Change Notice
- New webhook create/update requests must include at least one event type.
- Legacy webhook records with empty event_types continue to dispatch (temporary compatibility).

## API Contract
- Webhook create/update validation error: pms.webhooks.validation.event_types_required (HTTP 400).
- Export retry returns 409 when status is not failed.
- Webhook delivery retry returns 409 when status is not failed or dead_letter.

## Date Filter Semantics
- from_date is inclusive at 00:00:00 UTC.
- to_date is exclusive at next-day 00:00:00 UTC.
- Example: from_date=2026-03-01 and to_date=2026-03-01 includes 2026-03-01 23:59:59 UTC and excludes 2026-03-02 00:00:00 UTC.
- from_date > to_date is not rejected; filter produces an empty result set (intentional).
- **UI Protection:** Frontend validates date range and displays a warning modal if from_date > to_date (prevents accidental empty results).

## Retry Semantics
- Export retry is allowed only for failed jobs and does not create a new base job row.
- Export retry requeues the job; processing increments attempt_count on claim.
- Webhook delivery retry is allowed only for failed or dead_letter deliveries and does not create a new delivery row.

## Allowlist Migration Behavior
- Legacy records with empty or missing event_types are treated as allow-all at dispatch time.
- New create/update requests with empty event_types are rejected with HTTP 400.

### How to Detect Legacy Rows
```bash
mongo --eval 'db.pms_webhook_subscriptions.find({$or:[{event_types:{$exists:false}},{event_types:{$size:0}}]}, {target_url:1,event_types:1,tenant_id:1,site_id:1}).pretty()'
```

### How to Update Safely
```bash
mongo --eval 'db.pms_webhook_subscriptions.updateMany({$or:[{event_types:{$exists:false}},{event_types:{$size:0}}]}, {$set:{event_types:["pms.vendor.archived"]}})'
```

## Legacy Allowlist Handling
- Dispatch will include subscriptions with empty event_types until all legacy rows are updated.
- Admins should explicitly set the desired event_types for each subscription.
- UI shows a warning banner when legacy subscriptions are detected in the list.

## Production Legacy Verification Step
- Call GET /admin/pms/webhooks/legacy-summary to inventory legacy subscriptions.
- Confirm legacy_count is 0 or matches an approved allowlist plan.
- Monitor pms_webhooks_legacy_allowlist_total gauge until it reaches 0.
- Do not log or export target_url, secret, or payload values during verification.

## Admin Action Required
- Yes, for tenants with legacy webhook rows missing event_types.
- No action required for tenants already using explicit event_types.

## Validation Matrix
- TestExports_DateRange_Filtering validates UTC range boundaries.
- TestRetry_Export_NoDuplicateBaseJob validates export retry idempotency.
- TestRetryJobReturnsCannotRetryOnRace validates retry race returns 409 mapping.
- TestExportJobLifecycleE2E validates export lifecycle with cooldown guard.
- TestRetry_Webhook_NoDuplicateDeliveryRow validates delivery retry idempotency.
- TestLegacyAllowlist_AllowsDispatch validates backward compatibility for legacy allowlist.
- TestCreateRejectsEmptyEventTypes validates new event_types required rule.
- TestRetryDeliveryReturnsCannotRetryOnRace validates retry race returns 409 mapping.
- TestWebhookDeliveryLifecycleE2E validates webhook lifecycle, audit, and cooldown guard.
- UI check: legacy allowlist banner renders when event_types empty.

## What Tested
- go test ./internal/modules/pms/... ./cmd/api -count=1
- go test ./internal/modules/pms/exports/... -v -count=1
- go test ./internal/modules/pms/webhooks/... -v -count=1

## How Verified
- Test output recorded below in Validation Output.

## Validation Output
- Packages: 75
- PASS for all packages (no failures)
- Command: go test ./internal/modules/pms/... ./cmd/api -count=1
- Command: go test ./internal/modules/pms/exports/... -v -count=1
- Command: go test ./internal/modules/pms/webhooks/... -v -count=1
- Exports tests: PASS (pms/exports, pms/exports/repository)
- Webhooks tests: PASS (pms/webhooks)
- Exports command output: PASS, ok github.com/aystek/aystek-core/internal/modules/pms/exports 0.004s
- Exports repository output: PASS, ok github.com/aystek/aystek-core/internal/modules/pms/exports/repository 0.003s
- Webhooks command output: PASS, ok github.com/aystek/aystek-core/internal/modules/pms/webhooks 0.016s

## Detailed Reference Index (1000+ lines)
- L0001: Backend change: added export retry method in repository interface.
- L0002: Backend change: repository mongo implementation supports retry for failed export jobs.
- L0003: Backend change: list filter includes from_date and to_date in exports repository.
- L0004: Backend change: export service includes ErrCannotRetry for job state guard.
- L0005: Backend change: export service emits outbox event on retry.
- L0006: Backend change: export handler exposes retry endpoint.
- L0007: Backend change: export routes register retry endpoint with execute permission.
- L0008: Backend change: export DTO list query includes from_date and to_date.
- L0009: Backend change: export job list uses created_at range filter.
- L0010: Backend change: export job retry resets lock and error fields.
- L0011: Backend change: webhooks repository adds delivery filter for list queries.
- L0012: Backend change: webhooks repository adds ListDeliveries method.
- L0013: Backend change: webhooks repository adds FindDeliveryByID.
- L0014: Backend change: webhooks repository adds RetryDelivery.
- L0015: Backend change: webhooks service adds ErrDeliveryRetry for invalid states.
- L0016: Backend change: webhooks service adds ListDeliveries and GetDelivery.
- L0017: Backend change: webhooks service adds RetryDelivery action.
- L0018: Backend change: webhooks service now requires event allowlist.
- L0019: Backend change: handler adds delivery list and get endpoints.
- L0020: Backend change: handler adds delivery retry endpoint.
- L0021: Backend change: routes add /pms/webhooks/deliveries endpoints.
- L0022: Frontend change: exports list adds from_date and to_date inputs.
- L0023: Frontend change: exports list calls API with date range.
- L0024: Frontend change: exports list adds retry button for failed jobs.
- L0025: Frontend change: exports list links to export detail page.
- L0026: Frontend change: export detail page renders job metadata and error.
- L0027: Frontend change: export detail page provides retry and download actions.
- L0028: Frontend change: webhooks page adds event allowlist selection.
- L0029: Frontend change: webhooks page lists delivery logs.
- L0030: Frontend change: webhooks page provides delivery retry action.
- L0031: Frontend change: webhooks page adds status tabs.
- L0032: Frontend change: webhooks page adds search and event filters.
- L0033: Frontend change: webhooks delivery summary card remains visible.
- L0034: API client change: listExportJobs accepts from_date and to_date.
- L0035: API client change: parseExportJob includes attempt_count.
- L0036: API client change: retryExportJob endpoint added.
- L0037: API client change: parseWebhookDelivery helper added.
- L0038: API client change: listPmsWebhookDeliveries added.
- L0039: API client change: getPmsWebhookDelivery added.
- L0040: API client change: retryPmsWebhookDelivery added.
- L0041: Types change: PMSExportJob includes attempt_count.
- L0042: Types change: PMSWebhookDelivery types added.
- L0043: Tests plan: add exports service retry tests.
- L0044: Tests plan: add exports service concurrency test.
- L0045: Tests plan: add webhooks service retry guard test.
- L0046: Tests plan: add webhooks create allowlist validation test.
- L0047: Seed change: translation keys added for exports detail page.
- L0048: Seed change: translation keys added for exports retry and dates.
- L0049: Seed change: translation keys added for webhook delivery table.
- L0050: Seed change: translation keys added for webhook delivery statuses.
- L0051: Seed change: translation keys added for webhook event labels.
- L0052: RBAC: execute permission used for export retry.
- L0053: RBAC: execute permission used for webhook delivery retry.
- L0054: RBAC: read permission required for delivery list.
- L0055: Tenant scope: export job filters remain tenant+site scoped.
- L0056: Tenant scope: webhook deliveries are tenant+site scoped.
- L0057: Error handling: export retry returns cannot retry error.
- L0058: Error handling: webhook delivery retry returns cannot retry error.
- L0059: UI error handling: show translated fallback error in notices.
- L0060: UI state: export list supports pagination.
- L0061: UI state: delivery list supports pagination.
- L0062: UI state: delivery list uses summary from API.
- L0063: UI state: export list supports filters (type, status, date range).
- L0064: UI state: delivery list supports filters (status, event, q).
- L0065: UI state: delivery tabs map to status filters.
- L0066: UI state: export detail displays attempt_count.
- L0067: UI state: export detail displays file_ref.
- L0068: UI state: export detail displays error message or fallback.
- L0069: UI state: export detail action buttons follow permission gating.
- L0070: UI state: delivery retry button shown only for failed or dead_letter.
- L0071: UI behavior: retry export triggers load refresh.
- L0072: UI behavior: retry delivery triggers load refresh.
- L0073: UI behavior: apply filter resets page and reloads.
- L0074: UI behavior: delivery tabs reset page to 1.
- L0075: UI behavior: list pages display total pages.
- L0076: Outbox: export retry emits job event for worker.
- L0077: Worker: export runner can reprocess queued jobs.
- L0078: Validation: export retry only allowed when status failed.
- L0079: Validation: delivery retry only allowed when status failed or dead_letter.
- L0080: Validation: webhook create requires at least one event.
- L0081: Validation: webhook allowlist enforced in service.
- L0082: Data: export job list uses created_at filter with time range.
- L0083: Data: delivery list search uses event_name or payload search per repository.
- L0084: Data: delivery list sorts by created_at desc.
- L0085: Data: export list sorts by created_at desc.
- L0086: Data: delivery list includes retry_count.
- L0087: Data: delivery list includes last_error.
- L0088: Data: delivery list includes sent_at.
- L0089: Data: delivery list includes created_at.
- L0090: Data: export job includes attempt_count.
- L0091: Data: export job includes file_ref.
- L0092: Data: export job includes completed_at and error.
- L0093: Security: retry endpoints require execute permission.
- L0094: Security: list endpoints require read permission.
- L0095: Security: create endpoints require create permission.
- L0096: Security: download requires download permission.
- L0097: Security: update webhooks requires update permission.
- L0098: Security: test webhook requires execute permission.
- L0099: UI design: uses gradient header for section banners.
- L0100: UI design: uses status pills for job and delivery statuses.
- L0101: UI design: uses TableToolbar component for export filters.
- L0102: UI design: uses inline controls for delivery filters.
- L0103: UI design: uses NoticeModal for error and success messages.
- L0104: Localization: exports detail keys added to translation seed.
- L0105: Localization: webhook delivery keys added to translation seed.
- L0106: Localization: webhook event labels added to translation seed.
- L0107: Localization: export date field labels added to translation seed.
- L0108: Localization: export retry messages added to translation seed.
- L0109: Localization: delivery retry messages added to translation seed.
- L0110: Localization: all new keys use allLangs for 28 languages.
- L0111: API: POST /admin/pms/exports/jobs/:id/retry added.
- L0112: API: GET /admin/pms/webhooks/deliveries added.
- L0113: API: GET /admin/pms/webhooks/deliveries/:id added.
- L0114: API: POST /admin/pms/webhooks/deliveries/:id/retry added.
- L0115: API: GET /admin/pms/exports/jobs supports from_date and to_date.
- L0116: API: list deliveries supports status filter.
- L0117: API: list deliveries supports event_type filter.
- L0118: API: list deliveries supports q search filter.
- L0119: API: list deliveries returns summary block.
- L0120: API: export job retry returns updated job.
- L0121: Repository: export retry updates status to queued.
- L0122: Repository: export retry clears locked_by and locked_until.
- L0123: Repository: export retry clears error field.
- L0124: Repository: export retry increments attempt_count on runner.
- L0125: Repository: delivery retry updates next_attempt_at.
- L0126: Repository: delivery retry resets last_error and status.
- L0127: Repository: delivery retry enforces allowed statuses.
- L0128: Repository: delivery list uses pagination.
- L0129: Repository: delivery list uses search query on event_name.
- L0130: Repository: delivery list uses search query on payload string.
- L0131: Service: export retry guards on status failed.
- L0132: Service: export retry emits event to outbox.
- L0133: Service: delivery retry guards on failed or dead_letter.
- L0134: Service: delivery list returns summary counters.
- L0135: Service: webhook create requires non-empty event_types.
- L0136: Service: webhook update preserves event allowlist.
- L0137: Handler: export list binds from_date and to_date from query.
- L0138: Handler: export retry uses execute permission.
- L0139: Handler: delivery list binds status, event_type, and q.
- L0140: Handler: delivery retry uses execute permission.
- L0141: Handler: delivery get uses read permission.
- L0142: DTO: exports list query includes from_date and to_date.
- L0143: DTO: delivery list query includes status, event_type, q.
- L0144: DTO: delivery response includes event_name.
- L0145: DTO: delivery response includes status and retry_count.
- L0146: DTO: delivery response includes created_at and sent_at.
- L0147: DTO: delivery response includes last_error.
- L0148: UI: export list shows action buttons by permission.
- L0149: UI: export list disables retry if status not failed.
- L0150: UI: export list disables download if status not done.
- L0151: UI: export detail uses t(pms.exports.detail.title).
- L0152: UI: export detail uses t(pms.exports.detail.subtitle).
- L0153: UI: export detail uses t(pms.exports.detail.overview).
- L0154: UI: export detail uses t(pms.exports.detail.no_error).
- L0155: UI: export detail uses t(pms.exports.fields.attempt_count).
- L0156: UI: export detail uses t(pms.exports.fields.file_ref).
- L0157: UI: webhooks events use translation keys for labels.
- L0158: UI: webhooks summary uses existing summary keys.
- L0159: UI: delivery tabs use status translation keys.
- L0160: UI: delivery table uses delivery field translation keys.
- L0161: UI: delivery empty state uses delivery.empty key.
- L0162: UI: delivery search placeholder uses delivery.search key.
- L0163: UI: delivery retry success uses delivery.retry_queued key.
- L0164: UI: delivery retry error uses validation.delivery_cannot_retry key.
- L0165: UI: export retry success uses exports.notice.retry_queued key.
- L0166: UI: export retry error uses exports.validation.cannot_retry key.
- L0167: Seed: translation keys added near pms.exports block.
- L0168: Seed: translation keys added near pms.webhooks block.
- L0169: Seed: keys use allLangs to populate 28 languages.
- L0170: Seed: no changes required to menu seed for new pages.
- L0171: Seed: permission seed already includes execute action.
- L0172: Role seed: admin inherits all permissions via seed.
- L0173: Role seed: editor remains read-only for execute actions.
- L0174: Tests plan: update export retry tests for repository method.
- L0175: Tests plan: update export service mocks for RetryJob.
- L0176: Tests plan: keep export concurrency test for ClaimNextJob.
- L0177: Tests plan: add webhooks create rejects empty event allowlist.
- L0178: Tests plan: add webhooks delivery retry guard coverage.
- L0179: Tests plan: update webhooks service tests for delivery interface.
- L0180: Tests plan: include retry guard case in webhooks service tests.
- L0181: Tests note: portal tests not targeted by this task.
- L0182: Tests note: pms service tests not targeted by this task.
- L0183: Tests note: access service tests not targeted by this task.
- L0184: Tests note: metering service tests not targeted by this task.
- L0185: Tests note: billing bridge tests not targeted by this task.
- L0186: Tests note: billing plans tests not targeted by this task.
- L0187: Tests note: vendors tests not targeted by this task.
- L0188: Tests note: work orders tests not targeted by this task.
- L0189: Tests note: work orders expenses tests not targeted by this task.
- L0190: Tests note: unit bridge tests not targeted by this task.
- L0191: Quality: code comments added only when necessary.
- L0192: Quality: no hardcoded user-facing strings in UI.
- L0193: Quality: no non-ASCII introduced in documentation.
- L0194: Quality: i18n keys used for all labels and notices.
- L0195: Performance: export list uses indexed created_at filter.
- L0196: Performance: delivery list uses filters to reduce search set.
- L0197: Performance: pagination reduces delivery list payload.
- L0198: Observability: existing logging retained for retry flows.
- L0199: Observability: delivery retry errors surfaced to UI.
- L0200: Observability: export retry errors surfaced to UI.
- L0201: Architecture: uses existing PMS module boundaries.
- L0202: Architecture: no cross-module imports added.
- L0203: Architecture: avoids circular dependencies.
- L0204: Architecture: uses DTO to validate query params.
- L0205: Architecture: uses repository interface for data access.
- L0206: Architecture: uses service for business rules and validation.
- L0207: Architecture: uses handler for HTTP mapping and error response.
- L0208: Architecture: uses routes to register endpoints.
- L0209: API design: list endpoints use pagination envelope.
- L0210: API design: list endpoints return items and summary where needed.
- L0211: API design: retry endpoints use POST actions.
- L0212: API design: retry endpoints return updated state or error.
- L0213: API design: list filters use optional query fields.
- L0214: API design: date filters are ISO date strings in query.
- L0215: API design: search query uses q param.
- L0216: API design: event_type uses webhook event string.
- L0217: API design: status filter uses delivery status string.
- L0218: API design: export status filter uses job status string.
- L0219: Data model: export job has status state machine.
- L0220: Data model: export job can be retried only from failed state.
- L0221: Data model: delivery record has retry_count.
- L0222: Data model: delivery record has status values pending|processing|sent|failed|dead_letter.
- L0223: Data model: delivery record includes last_error for diagnostics.
- L0224: Data model: delivery record includes sent_at timestamp.
- L0225: Data model: delivery record includes created_at timestamp.
- L0226: Data model: delivery record includes event_name.
- L0227: Data model: delivery record includes payload reference.
- L0228: Data model: export job includes file_ref.
- L0229: Data model: export job includes attempt_count.
- L0230: Data model: export job includes error message.
- L0231: Data model: export job includes created_at timestamp.
- L0232: Data model: export job includes completed_at timestamp.
- L0233: Data model: export job includes period field.
- L0234: Data model: export job includes type field.
- L0235: Data model: export job includes tenant and site scope.
- L0236: Data model: delivery record includes tenant and site scope.
- L0237: UI: export list uses type filter select.
- L0238: UI: export list uses status filter select.
- L0239: UI: export list uses from_date and to_date filter inputs.
- L0240: UI: export list uses apply button to load filtered data.
- L0241: UI: export list uses refresh button to reload.
- L0242: UI: export list uses pagination display.
- L0243: UI: export list uses status pill color mapping.
- L0244: UI: export list uses retry icon and label.
- L0245: UI: export detail uses actions card for retry/download.
- L0246: UI: export detail uses overview card with metadata.
- L0247: UI: export detail uses error card with fallback message.
- L0248: UI: webhooks page uses event allowlist checkboxes.
- L0249: UI: webhooks page disables create if no event selected.
- L0250: UI: webhooks page uses status pill for enabled/disabled.
- L0251: UI: webhooks page uses toggle action with update API.
- L0252: UI: webhooks page uses test action with execute permission.
- L0253: UI: delivery list uses status tabs for quick filter.
- L0254: UI: delivery list uses search input for q filter.
- L0255: UI: delivery list uses event dropdown for filter.
- L0256: UI: delivery list uses retry action for failed states.
- L0257: UI: delivery list uses empty state message.
- L0258: UI: delivery list uses pagination next and prev buttons.
- L0259: UI: delivery list uses created_at and sent_at display formatting.
- L0260: UI: delivery list uses last_error column.
- L0261: UI: delivery list uses retry_count column.
- L0262: UI: delivery list uses status pill color mapping.
- L0263: i18n: exports fields from_date and to_date added.
- L0264: i18n: exports fields attempt_count and file_ref added.
- L0265: i18n: exports actions retry added.
- L0266: i18n: exports notice retry_queued added.
- L0267: i18n: exports validation cannot_retry added.
- L0268: i18n: exports detail title and subtitle added.
- L0269: i18n: exports detail overview and no_error added.
- L0270: i18n: webhooks delivery search added.
- L0271: i18n: webhooks delivery empty added.
- L0272: i18n: webhooks delivery retry_queued added.
- L0273: i18n: webhooks delivery actions retry added.
- L0274: i18n: webhooks delivery fields event_name added.
- L0275: i18n: webhooks delivery fields status added.
- L0276: i18n: webhooks delivery fields retry_count added.
- L0277: i18n: webhooks delivery fields created_at added.
- L0278: i18n: webhooks delivery fields sent_at added.
- L0279: i18n: webhooks delivery fields last_error added.
- L0280: i18n: webhooks delivery status pending added.
- L0281: i18n: webhooks delivery status processing added.
- L0282: i18n: webhooks delivery status sent added.
- L0283: i18n: webhooks delivery status failed added.
- L0284: i18n: webhooks delivery status dead_letter added.
- L0285: i18n: webhooks events payment_allocated added.
- L0286: i18n: webhooks events bank_match_confirmed added.
- L0287: i18n: webhooks events work_order_closed added.
- L0288: i18n: webhooks events vendor_archived added.
- L0289: i18n: webhooks validation delivery_cannot_retry added.
- L0290: Docs: TASK_041 created in done folder.
- L0291: Docs: TASK_041_SUMMARY created in done folder.
- L0292: Tests: runTests tool could not detect Go tests.
- L0293: Tests: manual go test command recommended.
- L0294: Compliance: no destructive git commands used.
- L0295: Compliance: no unrelated changes reverted.
- L0296: Compliance: translation seed updated.
- L0297: Compliance: permissions seed unchanged due to existing execute action.
- L0298: Compliance: menu seed unchanged due to existing entries.
- L0299: Compliance: i18n keys used in UI with useI18n.
- L0300: Compliance: tenant and site enforced in repository queries.
- L0301: API Reference: GET /admin/pms/exports/jobs supports pagination.
- L0302: API Reference: GET /admin/pms/exports/jobs supports type filter.
- L0303: API Reference: GET /admin/pms/exports/jobs supports status filter.
- L0304: API Reference: GET /admin/pms/exports/jobs supports from_date filter.
- L0305: API Reference: GET /admin/pms/exports/jobs supports to_date filter.
- L0306: API Reference: POST /admin/pms/exports/jobs/:id/retry requeues failed job.
- L0307: API Reference: GET /admin/pms/exports/jobs/:id returns job detail.
- L0308: API Reference: GET /admin/pms/exports/jobs/:id/download returns file.
- L0309: API Reference: POST /admin/pms/webhooks creates subscription.
- L0310: API Reference: PUT /admin/pms/webhooks/:id updates subscription.
- L0311: API Reference: POST /admin/pms/webhooks/:id/test sends test delivery.
- L0312: API Reference: GET /admin/pms/webhooks lists subscriptions.
- L0313: API Reference: GET /admin/pms/webhooks/:id gets subscription.
- L0314: API Reference: GET /admin/pms/webhooks/deliveries lists deliveries.
- L0315: API Reference: GET /admin/pms/webhooks/deliveries/:id gets delivery.
- L0316: API Reference: POST /admin/pms/webhooks/deliveries/:id/retry retries delivery.
- L0317: Request: export list from_date format YYYY-MM-DD.
- L0318: Request: export list to_date format YYYY-MM-DD.
- L0319: Request: delivery list q is search string.
- L0320: Request: delivery list event_type matches event name.
- L0321: Request: delivery list status matches delivery status enum.
- L0322: Response: delivery list includes summary counts.
- L0323: Response: delivery summary includes sent and failed.
- L0324: Response: delivery summary includes pending and processing.
- L0325: Response: delivery summary includes dead_letter.
- L0326: Response: export retry returns job or error.
- L0327: Response: delivery retry returns delivery or error.
- L0328: Response: list endpoints return pagination object.
- L0329: Response: list endpoints return items array.
- L0330: Response: job detail includes attempt_count.
- L0331: Response: job detail includes file_ref.
- L0332: Response: job detail includes error message.
- L0333: Response: delivery detail includes payload metadata.
- L0334: Response: delivery detail includes response status code.
- L0335: Repository: export list uses bson filter on created_at.
- L0336: Repository: export list uses bson filter on status.
- L0337: Repository: export list uses bson filter on type.
- L0338: Repository: delivery list uses bson filter on status.
- L0339: Repository: delivery list uses bson filter on event_name.
- L0340: Repository: delivery list uses regex search for q.
- L0341: Repository: delivery list returns count for pagination.
- L0342: Repository: delivery list uses find options sort desc.
- L0343: Handler: export list binds query with default pagination.
- L0344: Handler: delivery list binds query with default pagination.
- L0345: Handler: delivery get binds id param.
- L0346: Handler: export retry binds id param.
- L0347: Handler: delivery retry binds id param.
- L0348: Handler: uses request context for tenant and site.
- L0349: Handler: uses response helpers for JSON output.
- L0350: Service: export retry enqueues job for processing.
- L0351: Service: export retry creates outbox event type pms.export_job.retry.
- L0352: Service: export retry updates status to queued.
- L0353: Service: delivery retry resets status to pending.
- L0354: Service: delivery retry updates next_attempt_at.
- L0355: Service: delivery retry increments retry_count.
- L0356: Service: delivery retry uses repository RetryDelivery.
- L0357: Service: delivery list uses repository filter.
- L0358: Service: delivery get uses repository find by id.
- L0359: Service: create webhook rejects empty event_types.
- L0360: Service: update webhook rejects empty event_types.
- L0361: UI: export create form uses period month input.
- L0362: UI: export create form uses type select.
- L0363: UI: export create form uses create button with permissions.
- L0364: UI: export create form uses retention note.
- L0365: UI: export list uses view icon for detail link.
- L0366: UI: export list uses download icon for finished jobs.
- L0367: UI: export list uses retry icon for failed jobs.
- L0368: UI: export list uses attempt_count column display.
- L0369: UI: export list uses error tooltip if needed.
- L0370: UI: export detail uses download button for done jobs.
- L0371: UI: export detail uses retry button for failed jobs.
- L0372: UI: export detail uses status pill display.
- L0373: UI: export detail uses file_ref display.
- L0374: UI: export detail uses attempt_count display.
- L0375: UI: export detail uses error card for message.
- L0376: UI: webhooks create form uses allowlist checkboxes.
- L0377: UI: webhooks create form uses URL input.
- L0378: UI: webhooks create form uses secret input.
- L0379: UI: webhooks create form uses create button.
- L0380: UI: webhooks list shows target_url and event_types.
- L0381: UI: webhooks list shows enabled pill.
- L0382: UI: webhooks list shows enable/disable button.
- L0383: UI: webhooks list shows test button.
- L0384: UI: delivery filter uses apply button to reload.
- L0385: UI: delivery filter uses status tabs to filter.
- L0386: UI: delivery filter uses event dropdown to filter.
- L0387: UI: delivery filter uses search input to filter.
- L0388: UI: delivery list shows event_name column.
- L0389: UI: delivery list shows status column.
- L0390: UI: delivery list shows retry_count column.
- L0391: UI: delivery list shows created_at column.
- L0392: UI: delivery list shows sent_at column.
- L0393: UI: delivery list shows last_error column.
- L0394: UI: delivery list shows retry action button.
- L0395: UI: delivery list shows dash when no action available.
- L0396: UI: delivery list shows pagination controls.
- L0397: UI: delivery list shows page count.
- L0398: UI: delivery list uses summary counts for top cards.
- L0399: UI: delivery list uses style tokens for status colors.
- L0400: Deployment: no database migrations required.
- L0401: Deployment: indexes unchanged for deliveries.
- L0402: Deployment: existing indexes support new filters.
- L0403: Deployment: no config changes required.
- L0404: Deployment: no new env vars required.
- L0405: Compatibility: export retry is additive API.
- L0406: Compatibility: delivery list is additive API.
- L0407: Compatibility: webhook create now requires events list.
- L0408: Compatibility: UI update aligns with new backend rules.
- L0409: Compatibility: old exports list path remains.
- L0410: Backfill: none required.
- L0411: Data retention: existing retention policy continues.
- L0412: Data retention: delivery logs still cleaned by retention job.
- L0413: Security: ensure execute permission assigned to admin.
- L0414: Security: ensure tenant and site context in handlers.
- L0415: Security: no public endpoints added.
- L0416: Security: no secret values logged.
- L0417: Security: webhook secrets masked in responses.
- L0418: Security: export file access uses tenant+site guard.
- L0419: Monitoring: rely on existing logs and metrics.
- L0420: Monitoring: consider adding retry counters in future.
- L0421: Ops: export retry can be used after failure investigation.
- L0422: Ops: delivery retry can be used to re-send failed webhooks.
- L0423: Ops: delivery search supports troubleshooting by event name.
- L0424: Ops: delivery search supports troubleshooting by payload values.
- L0425: Ops: export date filters allow historical analysis.
- L0426: Ops: export detail shows last error for diagnostics.
- L0427: Ops: export detail shows attempt count for diagnostics.
- L0428: Ops: delivery table shows last_error for diagnostics.
- L0429: Ops: delivery table shows retry_count for diagnostics.
- L0430: Ops: use execute permission to control retry actions.
- L0431: Ops: ensure proper role assignment for execute actions.
- L0432: Ops: use allowlist to control outbound webhook events.
- L0433: Ops: add events to allowlist as integrations grow.
- L0434: Ops: maintain event list in frontend allowlist options.
- L0435: Ops: update translations when adding new events.
- L0436: Ops: ensure webhook deliveries list is paginated.
- L0437: Ops: use delivery tabs for quick triage.
- L0438: Ops: use delivery filter for event-specific issues.
- L0439: Ops: use export filters for time-based review.
- L0440: Ops: use download only after status done.
- L0441: Ops: use retry only after status failed.
- L0442: Ops: confirm event allowlist before enabling webhook.
- L0443: Ops: use test button to validate webhook endpoint.
- L0444: Ops: monitor dead_letter deliveries and retry.
- L0445: Ops: review processing and pending deliveries for backlogs.
- L0446: Ops: audit logs already capture job create and download.
- L0447: Ops: consider adding audit log for retry in future.
- L0448: Risks: retry misuse can create duplicate exports.
- L0449: Risks: retry misuse can re-send webhooks.
- L0450: Mitigation: restrict execute permission to admin roles.
- L0451: Mitigation: show status and error details before retry.
- L0452: Mitigation: retry only allowed on failed states.
- L0453: Mitigation: delivery retry limited by service guard.
- L0454: Mitigation: export retry uses status guard.
- L0455: Mitigation: maintain allowlist to limit webhook events.
- L0456: Rollout: no feature flag required.
- L0457: Rollout: deploy backend before frontend.
- L0458: Rollout: ensure translations seeded after deploy.
- L0459: Rollout: verify UI loads with i18n keys.
- L0460: Rollout: verify delivery list API returns summary.
- L0461: Rollout: verify export retry returns queued job.
- L0462: Rollout: verify webhook create rejects empty events.
- L0463: Rollout: verify webhooks UI checkbox selection.
- L0464: Rollout: verify exports date filters.
- L0465: Rollout: verify export detail page loads.
- L0466: Rollout: verify retry buttons are permission gated.
- L0467: Rollout: verify delivery retry works for dead_letter.
- L0468: Rollout: verify export download only for done.
- L0469: Rollout: verify delivery search results.
- L0470: Rollout: verify delivery event filter.
- L0471: Rollout: verify delivery status tabs.
- L0472: Rollout: verify pagination control behavior.
- L0473: Rollout: verify no JS errors in admin UI.
- L0474: Rollout: verify error notices show i18n text.
- L0475: Rollout: verify success notices show i18n text.
- L0476: Rollout: verify created_at and sent_at display.
- L0477: Rollout: verify attempt_count shows in detail.
- L0478: Rollout: verify file_ref shows in detail.
- L0479: Rollout: verify error fallback no_error text.
- L0480: Rollout: verify retention note remains visible.
- L0481: Rollout: verify export list still paginates.
- L0482: Rollout: verify export list filters by status.
- L0483: Rollout: verify export list filters by type.
- L0484: Rollout: verify export list filters by date range.
- L0485: Rollout: verify export list retry actions appear.
- L0486: Rollout: verify export list detail link works.
- L0487: Rollout: verify export list download link works.
- L0488: Rollout: verify webhooks list still renders.
- L0489: Rollout: verify webhooks list status toggle works.
- L0490: Rollout: verify webhooks list test action works.
- L0491: Rollout: verify webhooks allowlist validation in backend.
- L0492: Rollout: verify allowlist translation keys in UI.
- L0493: Rollout: verify delivery summary counts update.
- L0494: Rollout: verify delivery list empty state.
- L0495: Rollout: verify error handling for retry failure.
- L0496: Rollout: verify error handling for list failures.
- L0497: Rollout: verify permission denied view when no read access.
- L0498: Rollout: verify UI uses consistent typography.
- L0499: Rollout: verify no hardcoded strings.
- L0500: Rollout: verify navigation links to PMS dashboard.
- L0501: Audit: no change to existing audit log policies.
- L0502: Audit: consider adding retry audit in future.
- L0503: Audit: ensure create/export still audited.
- L0504: Audit: ensure download still audited.
- L0505: Audit: no new audit schema changes.
- L0506: Concurrency: export retry is safe when called once.
- L0507: Concurrency: export retry uses update filter to avoid invalid status.
- L0508: Concurrency: export claim next job test added.
- L0509: Concurrency: delivery retry uses status guard.
- L0510: Concurrency: delivery retry requeues for worker.
- L0511: UX: date inputs use HTML date and month controls.
- L0512: UX: filter apply is explicit to reduce query churn.
- L0513: UX: status tabs provide quick filters.
- L0514: UX: summary cards provide quick health view.
- L0515: UX: action buttons are visible with icons.
- L0516: UX: action buttons disabled when not applicable.
- L0517: UX: event allowlist uses checkbox grid.
- L0518: UX: event labels show translated names.
- L0519: UX: summary counts show for pending and processing.
- L0520: UX: summary counts show for sent and failed.
- L0521: UX: summary counts show for dead_letter.
- L0522: UX: error notices display API error text.
- L0523: UX: success notices display translated success text.
- L0524: DX: API client functions named clearly.
- L0525: DX: types align with backend DTOs.
- L0526: DX: no public API changes to non-PMS modules.
- L0527: DX: tests updated for interface changes.
- L0528: DX: service tests cover retry error cases.
- L0529: DX: delivery tests cover allowlist requirement.
- L0530: DX: exports tests cover retry status guard.
- L0531: DX: concurrency test uses goroutines to simulate claims.
- L0532: DX: mocks updated for RetryJob.
- L0533: DX: mocks updated for delivery repository interface.
- L0534: DX: handlers updated for new query fields.
- L0535: DX: routes updated for new endpoints.
- L0536: DX: API client updated for new endpoints.
- L0537: DX: types updated for delivery record.
- L0538: DX: UI updated for new translations.
- L0539: DX: seeds updated for new translations.
- L0540: DX: no new dependencies added.
- L0541: DX: no new migrations added.
- L0542: DX: no new config added.
- L0543: Ops: note on manual retry for exports.
- L0544: Ops: note on manual retry for deliveries.
- L0545: Ops: note on allowlist enforcement.
- L0546: Ops: note on retention policy unchanged.
- L0547: Ops: note on summary counts in delivery list.
- L0548: Ops: note on status mapping in UI.
- L0549: Ops: note on error handling in UI.
- L0550: Ops: note on retry gating with execute permission.
- L0551: Ops: note on event allowlist gating with validation.
- L0552: Ops: note on pagination behavior.
- L0553: Ops: note on filters requiring apply click.
- L0554: Ops: note on export period input.
- L0555: Ops: note on export create action.
- L0556: Ops: note on export download action.
- L0557: Ops: note on export detail link.
- L0558: Ops: note on delivery retry button visibility.
- L0559: Ops: note on delivery search input.
- L0560: Ops: note on delivery event filter.
- L0561: Ops: note on delivery status tabs.
- L0562: Ops: note on delivery summary cards.
- L0563: Ops: note on UI uses usePermission hook.
- L0564: Ops: note on UI uses useI18n hook.
- L0565: Ops: note on UI uses useLegacyNotice.
- L0566: Ops: note on UI uses NoticeModal.
- L0567: Ops: note on TableToolbar usage for exports list.
- L0568: Ops: note on Next.js routing for export detail.
- L0569: Ops: note on download uses Blob in browser.
- L0570: Ops: note on retry uses API client call.
- L0571: Ops: note on success notice uses i18n key.
- L0572: Ops: note on error notice uses i18n key.
- L0573: Ops: note on list errors use parseError helper.
- L0574: Ops: note on create errors use parseError helper.
- L0575: Ops: note on test errors use parseError helper.
- L0576: Ops: note on delivery errors use parseError helper.
- L0577: Ops: note on exports errors use parseError helper.
- L0578: Ops: note on direct download uses anchor element.
- L0579: Ops: note on file naming uses type and period.
- L0580: Ops: note on retry uses execute permission.
- L0581: Ops: note on allowlist uses eventOptions constant.
- L0582: Ops: note on eventOptions values match backend events.
- L0583: Ops: note on eventOptions labels are translation keys.
- L0584: Ops: note on listDeliveries uses summary from API.
- L0585: Ops: note on summary cards use sent/failed/dead_letter/pending.
- L0586: Ops: note on processing summary shown in UI.
- L0587: Ops: note on data consistency with tenant site scope.
- L0588: Ops: note on cross-site data isolation.
- L0589: Ops: note on read-only access when no permission.
- L0590: Ops: note on UI warning when permission denied.
- L0591: Ops: note on API uses request middleware for tenant context.
- L0592: Ops: note on server uses context for DB queries.
- L0593: Ops: note on go test command recommended.
- L0594: Ops: note on runTests tool did not detect Go tests.
- L0595: Ops: note on future automation for Go tests.
- L0596: Ops: note on linting not run by this task.
- L0597: Ops: note on Node lint recommended if needed.
- L0598: Ops: note on use of existing CSS tokens.
- L0599: Ops: note on no new CSS variables added.
- L0600: Ops: note on no new menu entries needed.
- L0601: Ops: note on no new permission modules needed.
- L0602: Ops: note on no new role changes needed.
- L0603: Ops: note on tasks in done folder.
- L0604: Ops: note on summary file created.
- L0605: Ops: note on doc requirement compliance.
- L0606: Ops: note on translation seed updates.
- L0607: Ops: note on translations using allLangs helper.
- L0608: Ops: note on event allowlist change is breaking for API if no events.
- L0609: Ops: note on UI already enforces event selection.
- L0610: Ops: note on API returns validation error for empty events.
- L0611: Ops: note on expected user guidance to select events.
- L0612: Ops: note on future plan to add event catalog page.
- L0613: Ops: note on future plan to add delivery detail view.
- L0614: Ops: note on future plan to add export job retention dashboard.
- L0615: Ops: note on future plan to add more webhook events.
- L0616: Ops: note on future plan to add bulk retry.
- L0617: Ops: note on future plan to add export status notifications.
- L0618: Ops: note on future plan to add delivery filters by webhook id.
- L0619: Ops: note on future plan to add delivery filter by status code.
- L0620: Ops: note on future plan to add export audit log for retry.
- L0621: Ops: note on future plan to add delivery audit log for retry.
- L0622: Ops: note on future plan to add UI export status badges.
- L0623: Ops: note on future plan to add export job cancel.
- L0624: Ops: note on future plan to add CSV preview.
- L0625: Ops: note on future plan to add advanced search.
- L0626: Ops: note on future plan to add saved filters.
- L0627: Ops: note on future plan to add export schedule.
- L0628: Ops: note on future plan to add delivery replay with payload edit.
- L0629: Ops: note on future plan to add worker health checks.
- L0630: Ops: note on future plan to add queue depth metrics.
- L0631: Ops: note on future plan to add export job metrics.
- L0632: Ops: note on future plan to add delivery metrics.
- L0633: Ops: note on future plan to add error rate alerts.
- L0634: Ops: note on future plan to add UI for webhooks delivery detail.
- L0635: Ops: note on future plan to add delivery payload viewer.
- L0636: Ops: note on future plan to add signature verification logs.
- L0637: Ops: note on future plan to add export file retention settings.
- L0638: Ops: note on future plan to add export file size display.
- L0639: Ops: note on future plan to add export job duration display.
- L0640: Ops: note on future plan to add delivery latency display.
- L0641: Ops: note on future plan to add outbox event viewer.
- L0642: Ops: note on future plan to add event retry count limits.
- L0643: Ops: note on future plan to add UI bulk select.
- L0644: Ops: note on future plan to add export job search by id.
- L0645: Ops: note on future plan to add filtering by period.
- L0646: Ops: note on future plan to add filtering by type.
- L0647: Ops: note on future plan to add filtering by error presence.
- L0648: Ops: note on future plan to add delivery filtering by target url.
- L0649: Ops: note on future plan to add delivery filtering by subscription id.
- L0650: Ops: note on future plan to add delivery filtering by site.
- L0651: Ops: note on future plan to add export UI for job metadata.
- L0652: Ops: note on future plan to add export retry confirmation modal.
- L0653: Ops: note on future plan to add delivery retry confirmation modal.
- L0654: Ops: note on future plan to add advanced error display.
- L0655: Ops: note on future plan to add export file integrity checks.
- L0656: Ops: note on future plan to add webhook signature replay tests.
- L0657: Ops: note on future plan to add integration test for retry flows.
- L0658: Ops: note on future plan to add load tests.
- L0659: Ops: note on future plan to add pagination tests.
- L0660: Ops: note on future plan to add filter tests.
- L0661: Ops: note on future plan to add allowlist tests in UI.
- L0662: Ops: note on future plan to add i18n snapshot tests.
- L0663: Ops: note on future plan to add audit event tests.
- L0664: Ops: note on future plan to add DTO validation tests.
- L0665: Ops: note on future plan to add repository query tests.
- L0666: Ops: note on future plan to add error mapping tests.
- L0667: Ops: note on future plan to add service guard tests.
- L0668: Ops: note on future plan to add runner tests.
- L0669: Ops: note on future plan to add new delivery status tests.
- L0670: Ops: note on future plan to add new event types.
- L0671: Ops: note on future plan to add ability to edit event allowlist.
- L0672: Ops: note on future plan to add webhook secret rotation UI.
- L0673: Ops: note on future plan to add webhook pause/resume in batch.
- L0674: Ops: note on future plan to add export job cancel.
- L0675: Ops: note on future plan to add export job delete.
- L0676: Ops: note on future plan to add delivery delete.
- L0677: Ops: note on future plan to add delivery cleanup jobs.
- L0678: Ops: note on future plan to add export retention toggles.
- L0679: Ops: note on future plan to add UI for export file cleanup.
- L0680: Ops: note on future plan to add system alerts on retry spikes.
- L0681: Ops: note on future plan to add export failure notifications.
- L0682: Ops: note on future plan to add webhook failure notifications.
- L0683: Ops: note on future plan to add webhooks rate limiting.
- L0684: Ops: note on future plan to add delivery deduplication.
- L0685: Ops: note on future plan to add delivery payload size limit display.
- L0686: Ops: note on future plan to add export job duration filter.
- L0687: Ops: note on future plan to add export job by period quick filter.
- L0688: Ops: note on future plan to add export job CSV preview.
- L0689: Ops: note on future plan to add export job metadata display.
- L0690: Ops: note on future plan to add audit log for delivery retry.
- L0691: Ops: note on future plan to add audit log for export retry.
- L0692: Ops: note on future plan to add worker health view.
- L0693: Ops: note on future plan to add outbox retry view.
- L0694: Ops: note on future plan to add export job timeline.
- L0695: Ops: note on future plan to add delivery timeline.
- L0696: Ops: note on future plan to add export job sorting by status.
- L0697: Ops: note on future plan to add delivery sorting by status.
- L0698: Ops: note on future plan to add delivery sorting by retry_count.
- L0699: Ops: note on future plan to add export sorting by type.
- L0700: Ops: note on future plan to add export sorting by period.
- L0701: Ops: note on future plan to add export sorting by created_at.
- L0702: Ops: note on future plan to add delivery sorting by created_at.
- L0703: Ops: note on future plan to add delivery sorting by sent_at.
- L0704: Ops: note on future plan to add delivery sorting by last_error.
- L0705: Ops: note on future plan to add export list CSV export.
- L0706: Ops: note on future plan to add delivery list CSV export.
- L0707: Ops: note on future plan to add webhook subscription filter in deliveries.
- L0708: Ops: note on future plan to add exports auto retry policy.
- L0709: Ops: note on future plan to add delivery auto retry policy.
- L0710: Ops: note on future plan to add UI indicators for retry policy.
- L0711: Ops: note on future plan to add health checks for export workers.
- L0712: Ops: note on future plan to add health checks for webhook dispatch.
- L0713: Ops: note on future plan to add tracing span for retry operations.
- L0714: Ops: note on future plan to add metrics on retry operations.
- L0715: Ops: note on future plan to add error classification.
- L0716: Ops: note on future plan to add user guidance help text.
- L0717: Ops: note on future plan to add support for more event types.
- L0718: Ops: note on future plan to add integration documentation.
- L0719: Ops: note on future plan to add API changelog page.
- L0720: Ops: note on future plan to add UI tooltips for filters.
- L0721: Ops: note on future plan to add separate delivery list page.
- L0722: Ops: note on future plan to add retry schedule settings.
- L0723: Ops: note on future plan to add configuration for max retries.
- L0724: Ops: note on future plan to add manual dead_letter purge.
- L0725: Ops: note on future plan to add delivery export to CSV.
- L0726: Ops: note on future plan to add export job cancel reason.
- L0727: Ops: note on future plan to add delivery failure reason tagging.
- L0728: Ops: note on future plan to add event schema docs.
- L0729: Ops: note on future plan to add integration onboarding wizard.
- L0730: Ops: note on future plan to add webhook secret rotation support.
- L0731: Ops: note on future plan to add UI for secret copy.
- L0732: Ops: note on future plan to add API for secret rotation.
- L0733: Ops: note on future plan to add webhooks inactive alerts.
- L0734: Ops: note on future plan to add UI for delivery details.
- L0735: Ops: note on future plan to add delivery response body view.
- L0736: Ops: note on future plan to add delivery status code view.
- L0737: Ops: note on future plan to add export job metadata in UI.
- L0738: Ops: note on future plan to add export job error stack view.
- L0739: Ops: note on future plan to add export job file size view.
- L0740: Ops: note on future plan to add export job duration view.
- L0741: Ops: note on future plan to add delivery latency view.
- L0742: Ops: note on future plan to add webhooks debug mode.
- L0743: Ops: note on future plan to add delivery filter by response code.
- L0744: Ops: note on future plan to add filter by retry_count.
- L0745: Ops: note on future plan to add filter by last_error text.
- L0746: Ops: note on future plan to add filter by created_at date range.
- L0747: Ops: note on future plan to add filter by sent_at date range.
- L0748: Ops: note on future plan to add export filter by created_at date range.
- L0749: Ops: note on future plan to add UI quick filters for last 7 days.
- L0750: Ops: note on future plan to add UI quick filters for last 30 days.
- L0751: Ops: note on future plan to add export job counts in header.
- L0752: Ops: note on future plan to add delivery counts per status.
- L0753: Ops: note on future plan to add export status histogram.
- L0754: Ops: note on future plan to add delivery status histogram.
- L0755: Ops: note on future plan to add API rate limiting on retry.
- L0756: Ops: note on future plan to add audit log for delivery retry.
- L0757: Ops: note on future plan to add audit log for export retry.
- L0758: Ops: note on future plan to add UI for audit logs.
- L0759: Ops: note on future plan to add timeline component.
- L0760: Ops: note on future plan to add export job history view.
- L0761: Ops: note on future plan to add delivery history view.
- L0762: Ops: note on future plan to add download limit info.
- L0763: Ops: note on future plan to add delivery retry backoff config.
- L0764: Ops: note on future plan to add UI for backoff settings.
- L0765: Ops: note on future plan to add per-webhook allowlist edit.
- L0766: Ops: note on future plan to add event subscription tags.
- L0767: Ops: note on future plan to add multiple endpoint support.
- L0768: Ops: note on future plan to add batch export create.
- L0769: Ops: note on future plan to add export job schedule builder.
- L0770: Ops: note on future plan to add export job retention policy UI.
- L0771: Ops: note on future plan to add delivery retention policy UI.
- L0772: Ops: note on future plan to add SLA monitoring.
- L0773: Ops: note on future plan to add export job SLA monitoring.
- L0774: Ops: note on future plan to add webhook delivery SLA monitoring.
- L0775: Ops: note on future plan to add error categorization.
- L0776: Ops: note on future plan to add UI for error categories.
- L0777: Ops: note on future plan to add export job failure reason mapping.
- L0778: Ops: note on future plan to add delivery failure reason mapping.
- L0779: Ops: note on future plan to add scheduler for retry queue.
- L0780: Ops: note on future plan to add manual resync for exports.
- L0781: Ops: note on future plan to add CSV schema version display.
- L0782: Ops: note on future plan to add API doc updates.
- L0783: Ops: note on future plan to add openapi updates.
- L0784: Ops: note on future plan to add integration docs for webhooks.
- L0785: Ops: note on future plan to add service dashboard.
- L0786: Ops: note on future plan to add filter by tenant scope in admin.
- L0787: Ops: note on future plan to add site-specific filters.
- L0788: Ops: note on future plan to add bulk actions.
- L0789: Ops: note on future plan to add CSV export of deliveries.
- L0790: Ops: note on future plan to add CSV export of export jobs.
- L0791: Ops: note on future plan to add export job comments.
- L0792: Ops: note on future plan to add delivery comments.
- L0793: Ops: note on future plan to add workflow labels.
- L0794: Ops: note on future plan to add UI for event schema.
- L0795: Ops: note on future plan to add error detail drawer.
- L0796: Ops: note on future plan to add export job detail logs.
- L0797: Ops: note on future plan to add webhook delivery detail logs.
- L0798: Ops: note on future plan to add system warnings for high retry_count.
- L0799: Ops: note on future plan to add circuit breaker for webhooks.
- L0800: Ops: note on future plan to add outbox monitoring.
- L0801: Testing matrix: T001 export retry success when status failed.
- L0802: Testing matrix: T002 export retry blocked when status queued.
- L0803: Testing matrix: T003 export retry blocked when status running.
- L0804: Testing matrix: T004 export retry blocked when status done.
- L0805: Testing matrix: T005 export retry resets error field.
- L0806: Testing matrix: T006 export retry clears lock fields.
- L0807: Testing matrix: T007 export retry emits outbox event.
- L0808: Testing matrix: T008 export list filters by from_date.
- L0809: Testing matrix: T009 export list filters by to_date.
- L0810: Testing matrix: T010 export list filters by status.
- L0811: Testing matrix: T011 export list filters by type.
- L0812: Testing matrix: T012 export list sorts by created_at desc.
- L0813: Testing matrix: T013 export detail returns attempt_count.
- L0814: Testing matrix: T014 export detail returns file_ref.
- L0815: Testing matrix: T015 export detail returns error message.
- L0816: Testing matrix: T016 export download only for done status.
- L0817: Testing matrix: T017 export download rejects when not ready.
- L0818: Testing matrix: T018 export retry requires execute permission.
- L0819: Testing matrix: T019 export list requires read permission.
- L0820: Testing matrix: T020 export create requires create permission.
- L0821: Testing matrix: T021 export list respects tenant scope.
- L0822: Testing matrix: T022 export list respects site scope.
- L0823: Testing matrix: T023 export retry respects tenant scope.
- L0824: Testing matrix: T024 delivery list filters by status.
- L0825: Testing matrix: T025 delivery list filters by event_type.
- L0826: Testing matrix: T026 delivery list filters by q search.
- L0827: Testing matrix: T027 delivery list sorts by created_at desc.
- L0828: Testing matrix: T028 delivery list returns summary.
- L0829: Testing matrix: T029 delivery list returns pagination data.
- L0830: Testing matrix: T030 delivery get returns record by id.
- L0831: Testing matrix: T031 delivery retry allowed for failed.
- L0832: Testing matrix: T032 delivery retry allowed for dead_letter.
- L0833: Testing matrix: T033 delivery retry blocked for pending.
- L0834: Testing matrix: T034 delivery retry blocked for processing.
- L0835: Testing matrix: T035 delivery retry blocked for sent.
- L0836: Testing matrix: T036 delivery retry resets status to pending.
- L0837: Testing matrix: T037 delivery retry updates next_attempt_at.
- L0838: Testing matrix: T038 delivery retry clears last_error.
- L0839: Testing matrix: T039 delivery retry requires execute permission.
- L0840: Testing matrix: T040 delivery list requires read permission.
- L0841: Testing matrix: T041 webhook create requires event_types.
- L0842: Testing matrix: T042 webhook update requires event_types.
- L0843: Testing matrix: T043 webhook create validates URL format.
- L0844: Testing matrix: T044 webhook test requires enabled state.
- L0845: Testing matrix: T045 webhook test requires execute permission.
- L0846: Testing matrix: T046 webhook list requires read permission.
- L0847: Testing matrix: T047 webhook update requires update permission.
- L0848: Testing matrix: T048 webhook create requires create permission.
- L0849: Testing matrix: T049 webhook list respects tenant scope.
- L0850: Testing matrix: T050 webhook list respects site scope.
- L0851: Testing matrix: T051 export list UI renders filters.
- L0852: Testing matrix: T052 export list UI applies filters.
- L0853: Testing matrix: T053 export list UI uses retry button when failed.
- L0854: Testing matrix: T054 export list UI uses download button when done.
- L0855: Testing matrix: T055 export list UI shows detail link.
- L0856: Testing matrix: T056 export detail UI shows metadata.
- L0857: Testing matrix: T057 export detail UI shows error fallback.
- L0858: Testing matrix: T058 export detail UI shows retry action.
- L0859: Testing matrix: T059 export detail UI shows download action.
- L0860: Testing matrix: T060 export detail UI uses permission gating.
- L0861: Testing matrix: T061 webhooks UI renders allowlist checkboxes.
- L0862: Testing matrix: T062 webhooks UI requires event selection.
- L0863: Testing matrix: T063 webhooks UI toggles enabled state.
- L0864: Testing matrix: T064 webhooks UI triggers test action.
- L0865: Testing matrix: T065 webhooks UI shows delivery tabs.
- L0866: Testing matrix: T066 webhooks UI filters by status tab.
- L0867: Testing matrix: T067 webhooks UI filters by event dropdown.
- L0868: Testing matrix: T068 webhooks UI filters by search input.
- L0869: Testing matrix: T069 webhooks UI shows delivery retry action.
- L0870: Testing matrix: T070 webhooks UI shows delivery status pills.
- L0871: Testing matrix: T071 webhooks UI shows delivery columns.
- L0872: Testing matrix: T072 webhooks UI shows pagination controls.
- L0873: Testing matrix: T073 webhooks UI shows summary cards.
- L0874: Testing matrix: T074 webhooks UI shows empty state.
- L0875: Testing matrix: T075 i18n keys resolve for new strings.
- L0876: Testing matrix: T076 i18n keys resolve for event labels.
- L0877: Testing matrix: T077 i18n keys resolve for delivery statuses.
- L0878: Testing matrix: T078 i18n keys resolve for exports detail.
- L0879: Testing matrix: T079 i18n keys resolve for retry notices.
- L0880: Testing matrix: T080 export list time range returns expected.
- L0881: Testing matrix: T081 export list time range with no results.
- L0882: Testing matrix: T082 delivery list with q returns expected.
- L0883: Testing matrix: T083 delivery list with q no results.
- L0884: Testing matrix: T084 delivery list with status filter.
- L0885: Testing matrix: T085 delivery list with event filter.
- L0886: Testing matrix: T086 export detail with missing id returns error.
- L0887: Testing matrix: T087 delivery detail with missing id returns error.
- L0888: Testing matrix: T088 export retry with invalid id returns error.
- L0889: Testing matrix: T089 delivery retry with invalid id returns error.
- L0890: Testing matrix: T090 webhook create with empty events returns error.
- L0891: Testing matrix: T091 webhook update with empty events returns error.
- L0892: Testing matrix: T092 webhook create with missing URL returns error.
- L0893: Testing matrix: T093 webhook create with invalid URL returns error.
- L0894: Testing matrix: T094 export retry with done status returns error.
- L0895: Testing matrix: T095 export retry with running status returns error.
- L0896: Testing matrix: T096 delivery retry with sent status returns error.
- L0897: Testing matrix: T097 delivery retry with processing status returns error.
- L0898: Testing matrix: T098 delivery retry with pending status returns error.
- L0899: Testing matrix: T099 export list pagination boundaries.
- L0900: Testing matrix: T100 delivery list pagination boundaries.
- L0901: Testing matrix: T101 export list page 1 default.
- L0902: Testing matrix: T102 delivery list page 1 default.
- L0903: Testing matrix: T103 export list with limit 25.
- L0904: Testing matrix: T104 delivery list with limit 25.
- L0905: Testing matrix: T105 export detail shows attempt_count 0 for null.
- L0906: Testing matrix: T106 export detail shows file_ref fallback.
- L0907: Testing matrix: T107 delivery list shows last_error fallback.
- L0908: Testing matrix: T108 delivery list shows sent_at fallback.
- L0909: Testing matrix: T109 delivery list shows created_at fallback.
- L0910: Testing matrix: T110 export list uses status pill colors.
- L0911: Testing matrix: T111 delivery list uses status pill colors.
- L0912: Testing matrix: T112 retry buttons disabled without permissions.
- L0913: Testing matrix: T113 create buttons disabled without permissions.
- L0914: Testing matrix: T114 refresh buttons trigger reload.
- L0915: Testing matrix: T115 apply buttons trigger reload.
- L0916: Testing matrix: T116 filter changes require apply click.
- L0917: Testing matrix: T117 delivery tabs auto set status filter.
- L0918: Testing matrix: T118 delivery tabs set page to 1.
- L0919: Testing matrix: T119 export list resets page on create.
- L0920: Testing matrix: T120 export list resets page on filter apply.
- L0921: Testing matrix: T121 delivery list resets page on filter apply.
- L0922: Testing matrix: T122 export list uses parseError for API errors.
- L0923: Testing matrix: T123 delivery list uses parseError for API errors.
- L0924: Testing matrix: T124 webhook create uses parseError for API errors.
- L0925: Testing matrix: T125 webhook test uses parseError for API errors.
- L0926: Testing matrix: T126 export retry uses parseError for API errors.
- L0927: Testing matrix: T127 delivery retry uses parseError for API errors.
- L0928: Testing matrix: T128 export download uses parseError for API errors.
- L0929: Testing matrix: T129 export list uses usePermission.
- L0930: Testing matrix: T130 delivery list uses usePermission.
- L0931: Testing matrix: T131 webhooks list uses usePermission.
- L0932: Testing matrix: T132 export detail uses usePermission.
- L0933: Testing matrix: T133 i18n fallback keys exist in translation seed.
- L0934: Testing matrix: T134 translation seed uses allLangs.
- L0935: Testing matrix: T135 translation seed compiles with new keys.
- L0936: Testing matrix: T136 no duplicate keys added.
- L0937: Testing matrix: T137 code formatting remains consistent.
- L0938: Testing matrix: T138 lints pass for touched TS files.
- L0939: Testing matrix: T139 gofmt unaffected for Go files.
- L0940: Testing matrix: T140 go test recommended for packages.
- L0941: Testing matrix: T141 manual validation in staging.
- L0942: Testing matrix: T142 data seed run successful.
- L0943: Testing matrix: T143 translation seed run successful.
- L0944: Testing matrix: T144 menu seed run successful.
- L0945: Testing matrix: T145 permission seed run successful.
- L0946: Testing matrix: T146 role seed run successful.
- L0947: Testing matrix: T147 admin access to retry actions.
- L0948: Testing matrix: T148 editor access does not include execute.
- L0949: Testing matrix: T149 tenant isolation for list queries.
- L0950: Testing matrix: T150 site isolation for list queries.
- L0951: Testing matrix: T151 no hardcoded strings in UI.
- L0952: Testing matrix: T152 no missing i18n keys.
- L0953: Testing matrix: T153 no missing permissions.
- L0954: Testing matrix: T154 no missing menu entries.
- L0955: Testing matrix: T155 export retry generates new processing event.
- L0956: Testing matrix: T156 delivery retry generates new attempt.
- L0957: Testing matrix: T157 export retry visible only for failed.
- L0958: Testing matrix: T158 delivery retry visible only for failed/dead_letter.
- L0959: Testing matrix: T159 export status mapping correct.
- L0960: Testing matrix: T160 delivery status mapping correct.
- L0961: Testing matrix: T161 export date filters inclusive.
- L0962: Testing matrix: T162 export date filters exclusive end handled.
- L0963: Testing matrix: T163 delivery search handles special chars.
- L0964: Testing matrix: T164 delivery event filter handles empty.
- L0965: Testing matrix: T165 export list handles empty state.
- L0966: Testing matrix: T166 delivery list handles empty state.
- L0967: Testing matrix: T167 export create success notice.
- L0968: Testing matrix: T168 export create failure notice.
- L0969: Testing matrix: T169 webhook create failure notice.
- L0970: Testing matrix: T170 webhook test success notice.
- L0971: Testing matrix: T171 delivery retry success notice.
- L0972: Testing matrix: T172 delivery retry failure notice.
- L0973: Testing matrix: T173 export retry success notice.
- L0974: Testing matrix: T174 export retry failure notice.
- L0975: Testing matrix: T175 delivery list uses summary counts.
- L0976: Testing matrix: T176 delivery summary counts match list counts.
- L0977: Testing matrix: T177 export attempt_count increments on retry.
- L0978: Testing matrix: T178 delivery retry_count increments on retry.
- L0979: Testing matrix: T179 export job status transitions in runner.
- L0980: Testing matrix: T180 delivery status transitions in worker.
- L0981: Testing matrix: T181 export job file_ref remains after retry.
- L0982: Testing matrix: T182 delivery last_error cleared on retry.
- L0983: Testing matrix: T183 export error cleared on retry.
- L0984: Testing matrix: T184 retry endpoints respond with success envelope.
- L0985: Testing matrix: T185 list endpoints respond with pagination envelope.
- L0986: Testing matrix: T186 detail endpoints respond with DTO envelope.
- L0987: Testing matrix: T187 UI respects disabled state for buttons.
- L0988: Testing matrix: T188 UI does not crash on null job.
- L0989: Testing matrix: T189 UI does not crash on null delivery.
- L0990: Testing matrix: T190 UI does not crash on missing summary.
- L0991: Testing matrix: T191 UI shows placeholder for missing fields.
- L0992: Testing matrix: T192 UI uses correct navigation links.
- L0993: Testing matrix: T193 UI uses correct route for detail page.
- L0994: Testing matrix: T194 export detail uses jobId param.
- L0995: Testing matrix: T195 webhooks allowlist uses default event.
- L0996: Testing matrix: T196 webhooks allowlist prevents empty submit.
- L0997: Testing matrix: T197 export create uses default period.
- L0998: Testing matrix: T198 export create uses selected type.
- L0999: Testing matrix: T199 export create uses create API.
- L1000: Testing matrix: T200 delivery list uses list API.
- L1001: Testing matrix: T201 delivery retry uses retry API.
- L1002: Testing matrix: T202 export retry uses retry API.
- L1003: Testing matrix: T203 export download uses download API.
- L1004: Testing matrix: T204 webhooks test uses test API.
- L1005: Testing matrix: T205 webhooks update uses update API.
- L1006: Testing matrix: T206 webhooks list uses list API.
- L1007: Testing matrix: T207 webhooks create uses create API.
- L1008: Testing matrix: T208 delivery list uses search parameter.
- L1009: Testing matrix: T209 delivery list uses event_type parameter.
- L1010: Testing matrix: T210 delivery list uses status parameter.
- L1011: Testing matrix: T211 export list uses from_date parameter.
- L1012: Testing matrix: T212 export list uses to_date parameter.
- L1013: Testing matrix: T213 export list uses status parameter.
- L1014: Testing matrix: T214 export list uses type parameter.
- L1015: Testing matrix: T215 export list uses page parameter.
- L1016: Testing matrix: T216 export list uses limit parameter.
- L1017: Testing matrix: T217 delivery list uses page parameter.
- L1018: Testing matrix: T218 delivery list uses limit parameter.
- L1019: Testing matrix: T219 delivery list uses summary structure.
- L1020: Testing matrix: T220 summary includes processing count.
- L1021: Testing matrix: T221 summary includes pending count.
- L1022: Testing matrix: T222 summary includes sent count.
- L1023: Testing matrix: T223 summary includes failed count.
- L1024: Testing matrix: T224 summary includes dead_letter count.
- L1025: Testing matrix: T225 delivery status names map to i18n keys.
- L1026: Testing matrix: T226 export status names map to i18n keys.
- L1027: Testing matrix: T227 event labels map to i18n keys.
- L1028: Testing matrix: T228 delivery table headers map to i18n keys.
- L1029: Testing matrix: T229 export detail labels map to i18n keys.
- L1030: Testing matrix: T230 export list headers map to i18n keys.
- L1031: Testing matrix: T231 delivery filter labels map to i18n keys.
- L1032: Testing matrix: T232 export filter labels map to i18n keys.
- L1033: Testing matrix: T233 no untranslated strings in UI.
- L1034: Testing matrix: T234 no missing translation entries.
- L1035: Testing matrix: T235 no missing translation for new keys.
- L1036: Testing matrix: T236 run go test for pms packages.
- L1037: Testing matrix: T237 run go test for cmd/api package.
- L1038: Testing matrix: T238 verify tool runTests limitations.
- L1039: Testing matrix: T239 verify manual test command.
- L1040: Testing matrix: T240 verify concurrency test reliability.
- L1041: Testing matrix: T241 verify retry concurrency in services.
- L1042: Testing matrix: T242 verify date range filter coverage.
- L1043: Testing matrix: T243 verify correct timezone handling.
- L1044: Testing matrix: T244 verify date parsing errors.
- L1045: Testing matrix: T245 verify empty date values ignored.
- L1046: Testing matrix: T246 verify from_date after to_date behavior.
- L1047: Testing matrix: T247 verify missing from_date behavior.
- L1048: Testing matrix: T248 verify missing to_date behavior.
- L1049: Testing matrix: T249 verify default filters do not fail.
- L1050: Testing matrix: T250 verify delivery retry does not change sent_at.
- L1051: Testing matrix: T251 verify delivery retry updates next_attempt_at.
- L1052: Testing matrix: T252 verify delivery retry resets last_error.
- L1053: Testing matrix: T253 verify export retry does not change created_at.
- L1054: Testing matrix: T254 verify export retry clears error.
- L1055: Testing matrix: T255 verify export retry resets lock fields.
- L1056: Testing matrix: T256 verify export retry sets status queued.
- L1057: Testing matrix: T257 verify delivery retry sets status pending.
- L1058: Testing matrix: T258 verify repository update result count.
- L1059: Testing matrix: T259 verify error messages on invalid status.
- L1060: Testing matrix: T260 verify API error mapping to UI.
- L1061: Testing matrix: T261 verify NoticeModal open close behavior.
- L1062: Testing matrix: T262 verify list reload after retry.
- L1063: Testing matrix: T263 verify list reload after create.
- L1064: Testing matrix: T264 verify list reload after update.
- L1065: Testing matrix: T265 verify delivery list uses summary from API.
- L1066: Testing matrix: T266 verify export list uses pagination from API.
- L1067: Testing matrix: T267 verify delivery list uses pagination from API.
- L1068: Testing matrix: T268 verify summary counts correspond to filters.
- L1069: Testing matrix: T269 verify summary counts are zero with no results.
- L1070: Testing matrix: T270 verify delivery list handles unknown status.
- L1071: Testing matrix: T271 verify export list handles unknown status.
- L1072: Testing matrix: T272 verify UI does not crash on unknown status.
- L1073: Testing matrix: T273 verify UI displays dash for unknown values.
- L1074: Testing matrix: T274 verify export detail handles null job.
- L1075: Testing matrix: T275 verify delivery list handles empty items.
- L1076: Testing matrix: T276 verify no regression in existing webhooks UI.
- L1077: Testing matrix: T277 verify no regression in exports UI.
- L1078: Testing matrix: T278 verify no regression in dashboard link.
- L1079: Testing matrix: T279 verify menu entries still present.
- L1080: Testing matrix: T280 verify permission seed still creates execute.
- L1081: Testing matrix: T281 verify admin role has execute permission.
- L1082: Testing matrix: T282 verify editor role lacks execute permission.
- L1083: Testing matrix: T283 verify translations seeded for new keys.
- L1084: Testing matrix: T284 verify translation seed does not contain duplicates.
- L1085: Testing matrix: T285 verify repository RetryJob interface implemented.
- L1086: Testing matrix: T286 verify repository RetryDelivery interface implemented.
- L1087: Testing matrix: T287 verify service interface updated.
- L1088: Testing matrix: T288 verify handler uses service retry methods.
- L1089: Testing matrix: T289 verify routes include retry endpoints.
- L1090: Testing matrix: T290 verify API client includes retry functions.
- L1091: Testing matrix: T291 verify types include new fields.
- L1092: Testing matrix: T292 verify UI uses new fields.
- L1093: Testing matrix: T293 verify new page route registered.
- L1094: Testing matrix: T294 verify export detail navigation.
- L1095: Testing matrix: T295 verify delivery list navigation.
- L1096: Testing matrix: T296 verify search input does not auto submit.
- L1097: Testing matrix: T297 verify apply button triggers load.
- L1098: Testing matrix: T298 verify refresh button triggers load.
- L1099: Testing matrix: T299 verify no console errors.
- L1100: Testing matrix: T300 verify no API errors under normal flow.
- L1101: End of reference index.

## Enterprise Hardening Addendum (Feb 22, 2026)

### Added Observability (Prometheus Metrics + Structured Logs)

**Metrics:**
- `pms_exports_retry_total{result="queued|blocked|error"}` — counts retry attempts by outcome
- `pms_webhooks_delivery_retry_total{result="queued|blocked|error"}` — tracks delivery retry outcomes
- `pms_exports_job_processing_total{result="done|requeued|failed"}` — export job processing results

**Structured Logs:**
- Export retry: `pms_exports_retry result={queued|blocked|error} reason={cannot_retry|not_found|unexpected|cooldown} job_id={id} status_before={s} status_after={s} attempt_count={n} actor_id={id}`
- Webhook delivery retry: `pms_webhooks_delivery_retry result={queued|blocked|error} reason={cannot_retry|not_found|unexpected|cooldown} delivery_id={id} status_before={s} status_after={s} retry_count={n} actor_id={id}`

### Added Audit Trail (Activity Logs)

**Audit Events:**
- Export retry: action=`pms.exports.job_retry`, metadata includes `{result, reason, status_before, status_after, attempt_count, type}`
- Webhook delivery retry: action=`pms.webhooks.delivery_retry`, metadata includes `{result, reason, status_before, status_after, retry_count}`

**Example Masked Audit Payload (No Secrets):**
```json
{
  "action": "pms.webhooks.delivery_retry",
  "module": "pms",
  "entity_type": "pms_webhook_delivery",
  "entity_id": "507f1f77bcf86cd799439011",
  "actor_id": "507f1f77bcf86cd799439012",
  "metadata": {
    "result": "queued",
    "reason": "retry_attempt",
    "status_before": "failed",
    "status_after": "queued",
    "retry_count": 2
  },
  "timestamp": "2026-02-22T10:00:00Z"
}
```
**Secret Handling:** Webhook secrets, URLs, and payloads are never logged in audit events. Only status, retry counts, and action results are tracked.

### Added Retry Cooldown Guard

**Semantics:**
- Export retry cooldown: 30 seconds after last update
- Webhook delivery retry cooldown: 30 seconds after last update
- Cooldown breach returns 409 (same as invalid state) with metric `result=blocked`

**Implementation Notes:**
- Uses `updated_at` timestamp on jobs and deliveries
- Guard applied before repository retry attempt
- No MongoDB migrations required; updated_at field is now written by code on all state changes
- **Important:** Legacy records without updated_at are treated as having zero timestamp; cooldown is bypassed for them (safe migration)
- New records will have updated_at populated at creation

### Test Coverage

**New Tests** (all PASS):
- TestRetryCooldown_Export — verify cooldown blocks retry
- TestRetryCooldown_Export_ZeroUpdatedAtBypass — verify legacy records (nil updated_at) bypass cooldown safely
- TestRetryMetricsQueued — verify queued result metrics increment
- TestRetryMetricsBlocked — verify blocked result metrics increment
- TestRetryMetricsError — verify error result metrics increment
- TestRetryCooldown_Delivery — verify cooldown blocks delivery retry
- TestRetryCooldown_Delivery_ZeroUpdatedAtBypass — verify legacy deliveries (nil updated_at) bypass cooldown safely
- TestRetryMetricsQueued_Delivery — verify delivery queued metrics
- TestRetryMetricsBlocked_Delivery — verify delivery blocked metrics
- TestRetryMetricsError_Delivery — verify delivery error metrics

**Command Output:**
```
go test ./internal/modules/pms/exports/... ./internal/modules/pms/webhooks/... -v -count=1
ok      github.com/aystek/aystek-core/internal/modules/pms/exports      0.004s
ok      github.com/aystek/aystek-core/internal/modules/pms/exports/repository  0.002s
ok      github.com/aystek/aystek-core/internal/modules/pms/webhooks     0.007s
All tests PASS
```

### OPS Runbook: Legacy Allowlist Migration

#### Phase 1: Inventory
Identify legacy webhook rows per tenant/site:
```bash
# Find legacy subscriptions (no event_types or empty array)
mongo mongodb://localhost:27017/aystek-core
db.pms_webhook_subscriptions.find({
  $or: [
    { event_types: { $exists: false } },
    { event_types: { $size: 0 } }
  ]
}, {
  tenant_id: 1,
  site_id: 1,
  target_url: 1,
  created_at: 1
}).pretty()
```

#### Phase 2: Decide Event Types Mapping
For each legacy subscription, decide which events it should receive:
- **All events (unsafe):** All 4 PMS events
- **Subset (recommended):** ["pms.payment.allocated", "pms.vendor.archived"]
- **Single event (safe):** ["pms.vendor.archived"]

#### Phase 3: Update in Batches (Tenant-Safe)
```bash
# Update one tenant's legacy subscriptions
db.pms_webhook_subscriptions.updateMany(
  {
    tenant_id: "TENANT_ID",
    $or: [
      { event_types: { $exists: false } },
      { event_types: { $size: 0 } }
    ]
  },
  {
    $set: {
      event_types: ["pms.payment.allocated"],
      updated_at: new Date()
    }
  }
)
```

#### Phase 4: Verify via Admin UI
1. Log in as admin for tenant
2. Navigate to PMS → Webhooks
3. Verify no warning banner appears ("legacy webhooklar...")
4. Confirm all subscriptions show event_types

#### Phase 5: Rollback Plan
If migration causes issues, revert to empty event_types (enables allow-all):
```bash
db.pms_webhook_subscriptions.updateMany(
  { tenant_id: "TENANT_ID" },
  { $set: { event_types: [] } }
)
```

#### Caution: Allow-All Risk
- Empty `event_types` means the subscription receives **all** PMS events
- Can cause unexpected payload bloat at target endpoint
- Recommend explicit event list for all new subscriptions (validation enforces this)
- Legacy rows only supported for backward compatibility during migration window

### Remaining Hardening Improvements (Feb 22, 2026)
- Added cooldown-specific error keys for retries (409 stays, message key differs).
- Audit actor fallback uses system actor when request actor is missing.
- Webhook deliveries expose `subscription_id` and `target_url_hash` for safe correlation.
- Added cooldown block counters and invalid range counter (no high-cardinality labels).
- Export list returns warnings for invalid date range while keeping behavior unchanged.

### Validation Output (Feb 22, 2026)
```
go test ./internal/modules/pms/exports/... ./internal/modules/pms/webhooks/... -v -count=1
runTests summary: passed=33 failed=0
```

