**Plugin Runtime**
This document describes plugin evaluation and runtime interactions.

**Runtime Model (Aligned With Canonical Evaluation Order)**
```text
Auth -> Tenant -> Site -> ModuleGate -> RBAC -> LicenseGate -> Flags -> Hook Evaluator
```

**Metadata Source Of Truth**
Admin interfaces consume backend metadata from:
- `/admin/plugins/metadata` (see `swagger.yaml#/tag/Plugins`)

**Capabilities**
Capabilities are declared by plugins and filtered by server-side evaluation. Plugins never bypass permission checks.

**Lifecycle Hooks**
Hooks are executed by the runtime engine during:
- install
- enable
- disable
- upgrade
- remove

**Failure Containment**
Hook failures are contained; they do not bypass isolation or compromise tenant boundaries.

**License And Flags**
- License always overrides flags and capability evaluation.
- Flags refine behavior; they never enable a disabled plugin.

**Example Metadata Response**
```json
{
  "success": true,
  "data": {
    "plugins": [
      {
        "key": "webhook-dispatcher",
        "is_enabled": true,
        "effective_enabled": true,
        "allowed_capabilities": ["integrations:send"],
        "license": {"plugin": true},
        "flags": {"effective": {"webhooks.enabled": true}}
      }
    ]
  }
}
```
