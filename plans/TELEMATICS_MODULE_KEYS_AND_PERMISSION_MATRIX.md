# Telematics Module Keys and Permission Matrix

## Module Keys

| Module Key | Purpose | Initial Menu Visibility |
|---|---|---|
| `telematics` | Root dashboard and module gate anchor | hidden |
| `telematics_devices` | Logical device registry and health metadata | hidden |
| `telematics_telemetry` | Current state ingestion and telemetry state reads | hidden |

## Initial Route Map

| Route | Module Gate | Permission Module | Actions |
|---|---|---|---|
| `/admin/telematics/dashboard` | `telematics` | `telematics` | `read` |
| `/admin/telematics/devices` | `telematics` | `telematics_devices` | `read`, `create` |
| `/admin/telematics/devices/:id` | `telematics` | `telematics_devices` | `read`, `update`, `delete` |
| `/admin/telematics/telemetry/current-state` | `telematics` | `telematics_telemetry` | `read` |
| `/admin/telematics/telemetry/current-state/ingest` | `telematics` | `telematics_telemetry` | `execute` |

## Initial Permission Actions

### `telematics`

- `read`

### `telematics_devices`

- `create`
- `read`
- `update`
- `delete`

### `telematics_telemetry`

- `read`
- `execute`

## Notes

- Menu seeds are intentionally hidden in this phase because frontend surfaces are not yet implemented.
- The root module key `telematics` is still required so tenant/site/module gating can be enabled before full UI rollout.
- Additional modules such as `telematics_trips`, `telematics_geofences`, `telematics_alerts` and `telematics_commands` will be added in later phases instead of overloading the foundation skeleton.
