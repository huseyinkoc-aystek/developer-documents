# ADR-002 Location Ownership and Compatibility

**Status:** Accepted  
**Date:** March 13, 2026

## Context

`stock/locations` is active in backend routes, seed menus and frontend pages, but location history and live tracking are not fundamentally stock concerns. A direct cut-over would break existing admin and portal flows.

## Decision

`telematics/telemetry` becomes the canonical write-owner for live state and future historical telemetry.

During migration:

1. `stock/locations` remains available.
2. `stock/locations` is treated as a compatibility projection or facade.
3. Existing stock location pages continue to operate while new telematics surfaces are introduced.

## Consequences

- No new strategic write logic will be added to `stock/locations`.
- Existing stock location pages can be migrated incrementally instead of with a hard cut.
- Telemetry, replay, trip and alert pipelines can be designed around a dedicated operational model.
