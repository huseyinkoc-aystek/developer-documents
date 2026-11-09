# Device Vendor Command Registry Slice 1 Checkpoint

Date: 2026-04-17
Scope: SMS integrations alignment / vendor command registry foundation

## Completed

- Stock SMS template model now carries vendor-registry fields:
  - `command_intent`
  - `vendor_key`
  - `product_type`
  - `required_variables`
- Backend create/update flow normalizes these fields before persistence.
- Stock SMS repository persists registry fields.
- Stock SMS service test coverage now includes normalization behavior.
- SMS template admin form now exposes registry fields with i18n-backed labels and placeholders.
- Integrations SMS command center now surfaces registry metadata in:
  - template library table
  - quick-send template metadata card
- Integrations SMS command center no longer frames vendor routing around only `Trio`.
  - current selectable vendor routes:
    - `trio`
    - `mobiliz`
    - `arvento`
    - `n2_mobil`
- Translation seed updated for registry-related SMS keys.
- Stock SMS template indexes now include vendor/intent/status discovery index.

## Validation

- `cd /workspace/go/aystek-core && go test ./internal/modules/stock/sms/... ./internal/seeds ./cmd/api`
- `cd /workspace/node/aystek-core-next && npm run lint -- --file 'app/(dashboard)/admin/stock/sms/templates/TemplateForm.tsx' --file 'app/(dashboard)/admin/integrations/sms-commands/SmsCommandCenterPage.tsx' --file 'lib/api/stock-ecosystem.ts' --file 'lib/types/stock-ecosystem.ts'`

## Architectural Position

- Frontend/menu placement under `Integrations` remains correct.
- Backend stock SMS domain remains under stock for now.
- Registry foundation is vendor-agnostic and does not force a `Trio`-specific module shape.

## Next Recommended Slice

- Define `command_intent` read-model and resolution rules:
  - derive vendor from device/catalog/integration anchors
  - derive product type from stock asset class
  - validate required variables before execution
- Move quick-send from `freeform provider + phone` toward:
  - resolved SIM target
  - resolved vendor route
  - approved command intent
