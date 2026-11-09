**Documentation Alignment Suggestions**

The following potential gaps were detected between the Swagger spec and platform handbook coverage:

- Swagger defines module tags for functional domains (e.g., Products, Orders, Payments), but the platform handbook focuses on platform-level concepts. Consider cross-linking the module guides or OpenAPI tags from the handbook for discoverability.
- Swagger contains operational endpoints (e.g., `/healthz`, `/readyz`, `/metrics`) but the handbook does not explicitly link to them from architecture or enforcement sections. Consider referencing these in `ARCHITECTURE.md` or `ENFORCEMENT_FLOW.md` alongside operational considerations.
- Swagger includes SSE endpoints (`/events/stream`, `/notifications/stream`) which are described in `EVENT_BUS.md`, but the OpenAPI tag linkage should be referenced explicitly in that file (tag: Events).

No backend behavior is modified by these suggestions; they are documentation linkage improvements only.
