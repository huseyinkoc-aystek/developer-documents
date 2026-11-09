# 🚀 Integration Scenario: Vehicle Assets → GPS Tracking

**Status:** Enterprise-Grade, Production-Ready  
**Date:** March 3, 2026  
**Architecture:** Event-Driven + RabbitMQ Async + Multi-Language UX  

---

## 1. Business Scenario

### Admin Workflow (Current State)
```
[Admin] → Opens "Stok" (Vehicle Assets) Screen
         → CREATE/UPDATE/DELETE Vehicle Mount
         → Operation saves to MongoDB
         → (END - no integration)
```

### Desired Workflow (After Implementation)
```
[Admin] → Opens "Stok" (Vehicle Assets) Screen
         → [NEW] Sees "Integration Configuration" section
             ├─ Available fields from event payload (28 languages)
             ├─ Field mapping suggestions
             └─ Matching hooks preview
         → CREATE/UPDATE/DELETE Vehicle Mount
         → Backend:
             ├─ Saves vehicle asset to MongoDB
             ├─ Publishes EVENT: stock.vehicle_asset.created/updated/deleted/unmounted
             └─ (END of sync operation)
         ↓ [ASYNC via RabbitMQ]
         → Integration Worker picks up event
         → Matches hook by event_key + tenant/site
         → Builds request using field_mapping ("license_plate" → "plaka", etc)
         → Calls external API (SOAP XML)
         → Records execution result
         ├─ Success: LOGGED + status recorded
         ├─ Failure: RETRY policy applied (exponential backoff)
         └─ Exhausted: MOVED to DLQ for manual replay
         ↓ [FRONTEND OBSERVABILITY]
         → Admin sees real-time status via SSE
         → Execution history in "Integrations" tab
```

---

## 2. Data Flow Diagram

```
DOMAIN MODULE (Stock)
    ↓
    Service.Create(VehicleAsset)
    ├─ validate
    ├─ insert to MongoDB ✅
    └─ publishVehicleAssetEvent("stock.vehicle_asset.created")
         ↓
PLATFORM EVENTS
    ├─ Type: "stock.vehicle_asset.created"
    ├─ TenantID: auto-scoped
    ├─ SiteID: auto-scoped
    ├─ Payload: {license_plate, serial_no, msisdn, integration_key, custom_gps, ...}
    └─ timestamp
         ↓
INTEGRATIONS EVENTBRIDGE LISTENER
    ├─ Hear event on channel
    ├─ Find matching hooks:
    │   WHERE event_key = "stock.vehicle_asset.created"
    │     AND connector.status = "active"
    │     AND hook.enabled = true
    │     AND tenant_id = event.tenant_id
    │     AND site_id = event.site_id
    └─ For each hook:
         ├─ Extract action from hook.action_id
         ├─ Build RabbitMQ job:
         │   {
         │     "hook_id": "...",
         │     "action_id": "...",
         │     "connector_id": "...",
         │     "event_payload": {...},
         │     "field_mapping": {"plaka": "license_plate", ...},
         │     "request_mapping": {...},
         │     "attempt": 1,
         │     "tenant_id": "...",
         │     "site_id": "..."
         │   }
         ├─ Publish to RabbitMQ queue: "integrations.hooks.async"
         └─ Record: Hook triggered for this event
             ↓
RABBITMQ QUEUE
    ├─ Durable queue: integrations.hooks.async
    ├─ TTL: 7 days (configurable)
    └─ Consumers: N worker instances (configurable concurrency)
             ↓
INTEGRATION WORKER (cmd/workers/integration_worker.go)
    ├─ Consume message from queue
    ├─ Extract job metadata + event_payload
    ├─ Load action template + request_mapping
    ├─ Build request body (request_builder.go):
    │   ├─ Apply field_mapping:
    │   │   event_payload["license_plate"] 
    │   │   → template field["plaka"]
    │   ├─ Render template:
    │   │   <plaka>{{.plaka}}</plaka>  ← "34ABC123"
    │   └─ Execute HTTP/SOAP call
    ├─ Handle response:
    │   ├─ 200-299: Record execution (SUCCESS)
    │   ├─ 4xx/5xx: Determine retry logic
    │   │   ├─ Retryable (500,502,503): Retry with backoff
    │   │   └─ Not retryable (400,401,403): DLQ
    │   └─ Retry exhausted: SEND TO DLQ
    ├─ Ack message to RabbitMQ
    └─ Update execution_status in MongoDB
             ↓
EXECUTION LOG (MongoDB collection: integration_executions)
    ├─ hook_id
    ├─ action_id
    ├─ event_key
    ├─ event_payload
    ├─ request_body
    ├─ response_body
    ├─ http_status
    ├─ error_message
    ├─ retry_count
    ├─ created_at
    └─ status: success|failed|retrying|dlq
             ↓
FRONTEND (aystek-core-next)
    ├─ Vehicle Asset detail page
    ├─ New tab: "Integration History"
    │   ├─ Real-time status (SSE)
    │   ├─ Last execution result
    │   ├─ Retry count
    │   └─ Error details (if failed)
    └─ During form edit:
        ├─ Show matching hooks
        ├─ Show available fields preview
        └─ Show field mapping suggestions
```

---

## 3. Frontend UX Components

### 3.1 Vehicle Asset Form - Configuration Section

**When editing/creating a vehicle asset with linked GPS product:**

```tsx
<Card className="border-blue-200 bg-blue-50">
  <CardHeader>
    <h3>🔗 Integration Configuration</h3>
    <p className="text-xs text-gray-600">This operation will trigger integrations</p>
  </CardHeader>
  
  <MatchingHooks />     {/* Shows hooks that match this event */}
  <AvailableFields/>    {/* Shows dropdown of event payload fields (28 languages) */}
  <FieldMappingGuide/>  {/* Read-only: current mapping for reference */}
</Card>
```

### 3.2 AvailableFields Component

Shows user exactly what fields will be available when event fires:

```tsx
interface AvailableFieldsProps {
  action: IntegrationAction;
  eventKey: string;      // "stock.vehicle_asset.created"
  locale: string;         // Current language code
}

// Displays:
// [TR] Available Event Fields (Mevcut Event Alanları)
// 
// license_plate     → Araç plakası (REQUIRED for template)
// serial_no         → GPS cihaz numarası (IMEI)
// msisdn            → Hat numarası (MSISDN) 
// operation         → İşlem durumu  
// integration_key   → Brand GPS bilgisi (OPTIONAL)
// custom_gps        → Özel GPS coordinates (OPTIONAL)
// vehicle_name      → Araç adı
// vehicle_status    → Araç durumu
// product_type      → Ürün türü
//
// Each field shows:
// - Field name (English identifier)
// - Description (translated to locale)
// - Icon: ⭐ REQUIRED | ℹ️ OPTIONAL | ⚠️ RARELY SET
```

### 3.3 MatchingHooks Component

Shows which hooks will execute on this operation:

```tsx
// MATCHING HOOKS FOR: stock.vehicle_asset.created
//
// ✅ GPS Bildirim Servisi (GPS_UPDATE)
//    Connector: GPS Bildirim Servisi (SOAP)
//    Mode: ASYNC (will run in background)
//    Retry: 3 attempts, exponential backoff
//    Last run: 2 hours ago (Status: SUCCESS)
//
// ✅ Telematics Integration
//    Connector: Telematics API (HTTP/JSON)
//    Mode: SYNC (must complete before save)
//    Status: ACTIVE
```

### 3.4 Integration History Tab

```tsx
// INTEGRATION EXECUTION HISTORY
//
// Last 10 executions for this vehicle asset:
//
// 2026-03-03 10:45:23 | ✅ SUCCESS  | GPS Bildirim   | 245ms
// 2026-03-03 09:30:12 | ✅ SUCCESS  | GPS Bildirim   | 312ms
// 2026-03-03 08:15:45 | ⚠️ RETRY-1  | GPS Bildirim   | failed connection
// 2026-03-03 08:16:52 | ✅ SUCCESS  | GPS Bildirim   | 289ms (retry succeeded)
// 2026-03-02 22:00:00 | 🔴 FAILED  | GPS Bildirim   | Max retries exhausted
//                                                      → Click to view details / retry manually
```

---

## 4. Backend Implementation Details

### 4.1 Event Definition

Already exists - examples:

```typescript
// Event published from service.go
events.Publish(ctx, events.Event{
  Type:       "stock.vehicle_asset.created",  // or .updated/.deleted/.unmounted
  Scope:      "stock",
  TenantID:   tenantID,
  SiteID:     siteID,
  EntityType: "vehicle_asset",
  EntityID:   asset.ID.Hex(),
  Payload: map[string]interface{}{
    "license_plate":   "34ABC123",      // from Vehicle
    "serial_no":       "ABC-12345",     // from ProductItem.SerialNo or IMEI
    "msisdn":          "5551234567",    // from ProductItem.MSISDN
    "operation":       "created",       // event type
    "integration_key": "gps_service_key", // from Brand.IntegrationKey
    "custom_gps":      "...",           // from ProductItem.CustomData
    "vehicle_name":    "...",           // from Vehicle.Name
    "vehicle_status":  "active",        // from Vehicle.Status
    "product_type":    "gps",           // from ProductItem.ProductType
  },
})
```

### 4.2 Hook Seed Entry

GPS action hook for vehicle asset events:

```go
// File: internal/seeds/integrations_seed.go
// Already exists - no changes needed

// Hook example:
{
  "event_key":    "stock.vehicle_asset.created",
  "connector_id": ObjectID(...),
  "action_id":    ObjectID(...),  // GPS_UPDATE action
  "action_keys":  ["GPS_UPDATE"],
  "mode":         "async",        // RabbitMQ queue job
  "phase":        "post",         // After vehicle asset is created
  "enabled":      true,
  "retry_limit":  3,
  "timeout_ms":   30000,
}
```

### 4.3 RabbitMQ Job Structure

Queue message publishing in eventbridge listener:

```go
// internal/modules/integrations/eventbridge/queue.go

type HookExecutionJob struct {
  HookID          primitive.ObjectID              `json:"hook_id"`
  ActionID        primitive.ObjectID              `json:"action_id"`
  ConnectorID     primitive.ObjectID              `json:"connector_id"`
  EventKey        string                          `json:"event_key"`
  EventPayload    map[string]interface{}          `json:"event_payload"`
  RequestMapping  *model.RequestMapping           `json:"request_mapping"`
  FieldMapping    map[string]string               `json:"field_mapping"` // Convenience copy
  Attempt         int                             `json:"attempt"`
  MaxRetries      int                             `json:"max_retries"`
  CreatedAt       time.Time                       `json:"created_at"`
  TenantID        string                          `json:"tenant_id"`
  SiteID          string                          `json:"site_id"`
}

// Queue: integrations.hooks.async
// Durable: true
// TTL: 7 days
// Consumers: configurable (env: INTEGRATION_WORKER_CONCURRENCY)
```

### 4.4 Worker Processing

```go
// cmd/workers/integration_worker.go - Main consumer

func (w *Worker) processJob(ctx context.Context, job HookExecutionJob) error {
  // 1. Load connector + action + credentials
  connector, _ := w.connectorsRepo.FindByID(ctx, job.ConnectorID)
  action, _ := w.actionsRepo.FindByID(ctx, job.ActionID)
  
  // 2. Build request using field_mapping
  requestBody, _ := w.executor.buildRequestBody(
    job.EventPayload,
    action.RequestMapping,
  )
  
  // 3. Execute HTTP call
  resp, err := http.Do(ctx, connector, action, requestBody)
  
  // 4. Record execution result
  execution := model.Execution{
    HookID:        job.HookID,
    ActionID:      job.ActionID,
    EventKey:      job.EventKey,
    EventPayload:  job.EventPayload,
    RequestBody:   requestBody,
    ResponseBody:  resp.Body,
    HTTPStatus:    resp.StatusCode,
    RetryCount:    job.Attempt - 1,
    Status:        deriveStatus(resp, err),
    CreatedAt:     time.Now(),
  }
  w.executionsRepo.Create(ctx, execution)
  
  // 5. Retry logic
  if resp.StatusCode >= 500 && job.Attempt < job.MaxRetries {
    delay := calculateBackoff(job.Attempt)
    w.queue.PublishWithDelay(ctx, job, delay)
    return nil  // Don't nack - will be processed later
  }
  
  // 6. DLQ if exhausted
  if job.Attempt >= job.MaxRetries && err != nil {
    w.dlqRepo.Create(ctx, DLQRecord{
      OriginalJob: job,
      ErrorReason: err.Error(),
      CreatedAt:   time.Now(),
    })
    return nil  // Don't nack - moved to DLQ
  }
  
  return nil  // Ack message - success or permanent failure
}
```

---

## 5. Frontend Components (Next.js)

### 5.1 AvailableFields Component

```tsx
// components/integrations/AvailableFields.tsx

interface AvailableFieldsProps {
  action?: IntegrationAction;
  availableFields?: Record<string, string>;
  locale: string;
}

export function AvailableFields({
  action,
  availableFields,
  locale,
}: AvailableFieldsProps) {
  const { t } = useI18n();
  
  const fields = availableFields || action?.request_mapping?.available_event_fields || {};
  
  if (Object.keys(fields).length === 0) {
    return null;
  }

  return (
    <div className="rounded-lg border border-cyan-200 bg-cyan-50 p-4">
      <h4 className="mb-3 font-semibold text-cyan-900">
        📋 {t("integrations.available_event_fields")}
      </h4>
      
      <div className="grid grid-cols-1 gap-2 sm:grid-cols-2 md:grid-cols-3">
        {Object.entries(fields).map(([fieldName, description]) => (
          <div key={fieldName} className="rounded bg-white p-2 text-xs">
            <code className="block rounded bg-gray-100 px-1.5 py-1 font-mono text-cyan-700">
              {fieldName}
            </code>
            <p className="mt-1 text-cyan-800">{description}</p>
          </div>
        ))}
      </div>
      
      <p className="mt-3 border-t border-cyan-200 pt-2 text-xs text-cyan-700">
        💡 {t("integrations.available_fields_hint")}
      </p>
    </div>
  );
}
```

### 5.2 MatchingHooks Component

```tsx
// components/integrations/MatchingHooks.tsx

interface MatchingHooksProps {
  vehicleAsset: VehicleAsset;
  eventKey: "stock.vehicle_asset.created" | "stock.vehicle_asset.updated" | ...;
}

export function MatchingHooks({ vehicleAsset, eventKey }: MatchingHooksProps) {
  const [hooks, setHooks] = useState<IntegrationHook[]>([]);
  const [loading, setLoading] = useState(true);
  const { t } = useI18n();

  useEffect(() => {
    integrationsApi
      .listHooks({
        event_key: eventKey,
        status: "active",
      })
      .then(setHooks)
      .finally(() => setLoading(false));
  }, [eventKey]);

  if (loading) return <Spinner />;
  if (hooks.length === 0) {
    return (
      <div className="rounded bg-gray-50 p-3 text-sm text-gray-600">
        {t("integrations.no_matching_hooks")}
      </div>
    );
  }

  return (
    <div className="space-y-2">
      <h4 className="text-sm font-semibold">
        🔗 {t("integrations.matching_hooks_for_event", { event: eventKey })}
      </h4>
      
      {hooks.map((hook) => (
        <HookStatusCard key={hook.id} hook={hook} vehicleAssetId={vehicleAsset.id} />
      ))}
    </div>
  );
}

function HookStatusCard({ hook, vehicleAssetId }: { hook: IntegrationHook; vehicleAssetId: string }) {
  const [lastExecution, setLastExecution] = useState<Execution | null>(null);
  const { t } = useI18n();

  useEffect(() => {
    // Fetch last execution for this hook
    integrationsApi
      .getLastExecution({
        hook_id: hook.id,
        entity_id: vehicleAssetId,
      })
      .then(setLastExecution)
      .catch(() => {});
  }, [hook.id, vehicleAssetId]);

  return (
    <div className="flex items-center rounded border border-blue-200 bg-blue-50 p-3">
      <div className="flex-1">
        <h5 className="font-medium text-blue-900">{hook.name}</h5>
        <p className="text-xs text-blue-700">
          {t("integrations.mode")}: {hook.mode === "async" ? "⏱️ ASYNC" : "⚡ SYNC"} |{" "}
          {t("integrations.retry_limit")}: {hook.retry_limit}
        </p>
        
        {lastExecution && (
          <p className="mt-1 text-xs text-blue-600">
            <StatusIcon status={lastExecution.status} /> Last run: {formatTime(lastExecution.created_at)}
          </p>
        )}
      </div>
      
      <div className="ml-3">
        {lastExecution?.status === "success" && <CheckCircle className="text-green-500" />}
        {lastExecution?.status === "failed" && <AlertCircle className="text-red-500" />}
        {!lastExecution && <Info className="text-gray-400" />}
      </div>
    </div>
  );
}
```

### 5.3 Integration History Tab

```tsx
// components/integrations/IntegrationHistoryTab.tsx

interface ExecutionListProps {
  vehicleAssetId: string;
}

export function IntegrationHistoryTab({ vehicleAssetId }: ExecutionListProps) {
  const [executions, setExecutions] = useState<Execution[]>([]);
  const [loading, setLoading] = useState(true);
  const { t } = useI18n();

  useEffect(() => {
    const unsubscribe = integrationsApi.subscribeToExecutions(
      { entity_id: vehicleAssetId },
      setExecutions,
    );
    setLoading(false);
    return unsubscribe;
  }, [vehicleAssetId]);

  if (loading) return <Spinner />;
  if (executions.length === 0) {
    return <div className="text-sm text-gray-500">{t("integrations.no_executions")}</div>;
  }

  return (
    <table className="w-full text-sm">
      <thead>
        <tr>
          <th>{t("common.time")}</th>
          <th>{t("integrations.hook")}</th>
          <th>{t("integrations.status")}</th>
          <th>{t("integrations.duration")}</th>
          <th>{t("common.actions")}</th>
        </tr>
      </thead>
      <tbody>
        {executions.map((exec) => (
          <tr key={exec.id} className="border-t">
            <td className="py-2">{formatDateTime(exec.created_at)}</td>
            <td>{exec.hook_name}</td>
            <td>
              <StatusBadge status={exec.status} />
              {exec.status === "retrying" && ` (${exec.retry_count}/${exec.max_retries})`}
            </td>
            <td>{exec.duration_ms}ms</td>
            <td>
              {exec.status === "failed" && (
                <button
                  onClick={() => handleRetry(exec.id)}
                  className="text-blue-600 hover:underline"
                >
                  {t("common.retry")}
                </button>
              )}
              <button
                onClick={() => setDetailOpen(exec.id)}
                className="ml-2 text-gray-600 hover:underline"
              >
                {t("common.details")}
              </button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
```

---

## 6. Translation Keys (28 Languages)

Add to `translation_seed.go`:

```go
"integrations.available_event_fields": {
  "tr": "Mevcut Event Alanları",
  "en": "Available Event Fields",
  "de": "Verfügbare Ereignisfelder",
  "fr": "Champs d'événement disponibles",
  // ... 24 more languages ...
},

"integrations.available_fields_hint": {
  "tr": "Bu alanlar event yayınlandığında payload'da mevcut olacaktır",
  "en": "These fields will be available in the payload when the event is published",
  // ... 26 more ...
},

"integrations.matching_hooks_for_event": {
  "tr": "Bu event için eşleşen entegrasyonlar ({{event}})",
  "en": "Matching integrations for this event ({{event}})",
  // ... 26 more ...
},

"integrations.no_matching_hooks": {
  "tr": "Bu event için eşleşen entegrasyon yok",
  "en": "No matching integrations for this event",
  // ... 26 more ...
},
```

---

## 7. Sequence Diagram

```
ADMIN              FRONTEND           BACKEND            RABBITMQ           WORKER
  │                   │                   │                   │                 │
  ├── Create Vehicle ─┤                   │                   │                 │
  │   Asset (mount)   │                   │                   │                 │
  │                   ├── Show ────────────────────────────────────────────────>│
  │                   │ AvailableFields   │                   │                 │
  │                   │ MatchingHooks     │                   │                 │
  │                   │                   │                   │                 │
  │                   ├── POST ──────────>│                   │                 │
  │                   │ /save             │                   │                 │
  │                   │                   ├─ Validate ────────┤                 │
  │                   │                   ├─ Insert to DB ✅  │                 │
  │                   │                   ├─ Publish Event    │                 │
  │                   │                   │  "stock.vehicle   │                 │
  │                   │                   │   _asset.created" │                 │
  │                   │                   │────────────────>│                   │
  │                   │                   │                 │─ Match Hook     │
  │                   │                   │                 │─ Build Job      │
  │                   │                   │                 │─ ENQUEUE ────────>│
  │                   │                   │                 │                 ├─ Process
  │                   │                   │                 │                 ├─ Build Req
  │                   │<────── OK ────────┤                 │                 ├─ HTTP Call
  │                   │ (sync done)       │                 │                 ├─ Record
  │                   │                   │                 │                 │
  │                   │ [SSE updates]     │                 │                 │
  │                   │<────────────────────────────────────────────>         │
  │                   │ status: processing                  │                 │
  │                   │                   │                 │                 │<─ Ack
  │                   │<────────────────────────────────────────────┤          │
  │                   │ status: success/failed              │                 │
  │                   │                   │                 │                 │
```

---

## 8. Implementation Checklist

### Backend (Go)

- [x] Event publishing (already done in service.go)
- [ ] Eventbridge listener for stock events
- [ ] RabbitMQ job publishing
- [ ] Worker implementation
- [ ] Execution logging
- [ ] DLQ handling
- [ ] Retry policy

### Frontend (Next.js)

- [ ] AvailableFields component
- [ ] MatchingHooks component
- [ ] IntegrationHistoryTab component
- [ ] SSE subscription for real-time status
- [ ] Translation keys (28 languages)
- [ ] Integration history API
- [ ] Manual retry button

### Infrastructure

- [x] RabbitMQ (already running)
- [x] MongoDB execution logging
- [ ] Worker startup in docker-compose
- [ ] Environment variables (concurrency, retry limits)

---

## 9. Next Steps

1. **Implement Eventbridge Listener** - Listen to `stock.*` events, find matching hooks
2. **Implement RabbitMQ Job Publishing** - Enqueue HookExecutionJob
3. **Enhance Worker** - Implement field mapping, request building, retry/DLQ
4. **Frontend Components** - Build UX for available fields + history
5. **Testing** - E2E test: create vehicle asset → execute integration → verify result
6. **Documentation** - Add to admin panel help

---

**END OF SCENARIO DOCUMENT**
