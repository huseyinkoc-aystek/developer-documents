# TASK 025 Summary

PMS export flow was upgraded from direct endpoint usage to a job-based Exports Center.

## Delivered
- Job lifecycle API (`create/list/get/download`) under `/admin/pms/exports/jobs`.
- Async in-process runner (MVP) with persisted files.
- Outbox and audit integration for export jobs.
- New admin page `/admin/pms/exports` with filters and download action.
- Menu, permission, and i18n keys updated.

## Notes
- Existing direct export endpoint remains for backward compatibility.
- Exports Center is now the canonical admin UX for PMS exports.
