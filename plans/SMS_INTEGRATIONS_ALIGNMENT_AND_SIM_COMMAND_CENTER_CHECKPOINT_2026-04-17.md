# SMS Integrations Alignment And SIM Command Center Checkpoint

Date: 2026-04-17

## Completed

- stock SMS command menu entry moved under `Integrations`
- messaging group nested under `Integrations`
- new integrations-native route added:
  - `/admin/integrations/sms-commands`
- legacy stock list route redirected:
  - `/admin/stock/sms/templates`
- integrations-native SMS command center UI added
- architecture decision recorded:
  - frontend alignment yes
  - backend stock SMS physical relocation no
- device vendor / SIM orchestration analysis recorded

## Verified

- Next lint passed for new route and legacy redirect
- Go seed/api test passed

## Follow-up

- open next slice for device vendor command orchestration
- define `command_intent` contract
- add device vendor command registry and state matrix
