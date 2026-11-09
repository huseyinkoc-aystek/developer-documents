# System Ops Module Slice 1 Checkpoint

Date: 2026-04-04
Scope: `SYSTEM_OPS_MODULE_FEATURE_COMPLETION_AND_CLOSURE_EXECUTION_PLAN_2026-04-04.md`
Slice: `Slice 1 - Update and Rollback Realization`

## Completed

- `system_ops` backend update start contract now persists:
  - `tarball_url`
  - `confirm_token` validation
  - `options.take_backup`
  - `options.maintenance_mode`
  - `options.auto_rollback_on_health_fail`
  - `options.dry_run`
- rollback start flow now supports:
  - blank target version
  - previous-version resolution from latest non-rollback record
  - filesystem symlink fallback
  - confirm token validation
- rollback listing no longer filters after pagination; service-level rollback listing is now explicit
- update check endpoint moved from fixed stub to env-driven response:
  - `SYSTEM_OPS_CURRENT_VERSION`
  - `SYSTEM_OPS_LATEST_VERSION`
  - `SYSTEM_OPS_LATEST_CHANNEL`
  - `SYSTEM_OPS_RELEASE_NOTES`
  - `SYSTEM_OPS_RELEASE_DATE`
- Next admin `updates/start` screen now exposes tarball URL, confirm token, and update options
- Next admin `rollbacks/start` screen now supports:
  - optional target version
  - previous-version shortcut
  - confirm token field
- translation seed now includes Slice 1 update/rollback UI keys
- backend service tests added for:
  - previous-version resolution
  - symlink fallback
  - confirm token validation

## Files

### Go

- `/workspace/go/aystek-core/internal/modules/system_ops/model/update.go`
- `/workspace/go/aystek-core/internal/modules/system_ops/dto/system_ops.go`
- `/workspace/go/aystek-core/internal/modules/system_ops/repository/interface.go`
- `/workspace/go/aystek-core/internal/modules/system_ops/repository/mongo.go`
- `/workspace/go/aystek-core/internal/modules/system_ops/service/service.go`
- `/workspace/go/aystek-core/internal/modules/system_ops/service/service_test.go`
- `/workspace/go/aystek-core/internal/modules/system_ops/handler/admin.go`
- `/workspace/go/aystek-core/internal/seeds/translation_seed.go`

### Next

- `/workspace/node/aystek-core-next/lib/types/system-ops.ts`
- `/workspace/node/aystek-core-next/lib/api/system-ops.ts`
- `/workspace/node/aystek-core-next/app/(dashboard)/admin/system/updates/start/page.tsx`
- `/workspace/node/aystek-core-next/app/(dashboard)/admin/system/rollbacks/start/page.tsx`

## Verification

- `cd /workspace/go/aystek-core && go test ./internal/modules/system_ops/... ./internal/seeds ./cmd/api`
- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/system/updates/start/page.tsx' --file 'app/(dashboard)/admin/system/rollbacks/start/page.tsx' --file 'lib/api/system-ops.ts' --file 'lib/types/system-ops.ts' --file 'components/system-ops/SystemOpsFieldHintAction.tsx'`

## Remaining Inside Slice 1

- decide whether update availability should stay env-driven or move behind a provider/service abstraction
- optionally expose backend-side effective update plan in update detail output so selected options are visible post-start

## Next Recommended Step

- continue with `Slice 2: Maintenance and Backup Realization`
- highest-value targets:
  - persistent maintenance state
  - real backup lifecycle semantics
  - backup note usage
  - admin/UI contract alignment for manual backup flow
