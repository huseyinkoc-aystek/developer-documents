**Event Bus**
AYSTEK Core uses an internal event bus and SSE streams for operational and user-facing events.

**SSE Streams**
The platform exposes SSE endpoints for notifications:
- `/events/stream`
- `/notifications/stream`

Refer to `swagger.yaml#/tag/Events` for endpoints.

**Flow**
```text
Domain Event
  -> Event Bus
  -> SSE Stream
  -> Client (with access token or cookie)
```

**Authentication**
SSE endpoints require authentication. The access token may be provided via:
- `Authorization: Bearer <token>`
- `access_token` cookie (for EventSource)

**Example Connection**
```http
GET /events/stream?token=<access_token> HTTP/1.1
Host: api.example.com
```
