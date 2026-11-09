# ADR-001 Telematics Domain Boundaries

**Status:** Accepted  
**Date:** March 13, 2026

## Context

AYSTEK Core already contains active stock, integration, location, PMS and notification capabilities. The vehicle tracking initiative must not duplicate these domains or create conflicting ownership.

## Decision

The fleet tracking solution will be split into three primary bounded contexts:

1. `integrations`
   Owns vendor connector governance, webhooks, scheduled pulls, partner APIs, retry, DLQ and event-hook orchestration.
2. `stock`
   Owns vehicle master-data, physical inventory, mounted assets, sales, returns and stock-driven operational lifecycle.
3. `telematics`
   Owns live telemetry, current state, trips, geofences, alerts, driver behavior, CAN/OBD and command workflows.

`stock/locations` is retained only as a compatibility projection until the telematics read surface fully replaces it.

## Consequences

- New telematics vendor integrations will be modeled on top of `integrations`, not `stock/integrations`.
- `stock/vehicles`, `stock/inventory` and `stock/vehicle_assets` remain the authoritative master-data chain in early phases.
- New tracking workflows can evolve independently without polluting the stock domain.
- A later vehicle-domain re-segmentation remains possible, but it is not part of the foundation phase.
