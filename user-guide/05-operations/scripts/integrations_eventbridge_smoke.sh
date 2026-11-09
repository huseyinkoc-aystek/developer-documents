#!/usr/bin/env bash
set -euo pipefail

# Post-deploy smoke check for integrations hook pipeline (module-agnostic).
# Requires: curl, jq

: "${API_BASE_URL:?API_BASE_URL is required (example: http://localhost:8080)}"
: "${AUTH_TOKEN:?AUTH_TOKEN is required (Bearer token with integrations read/execute permissions)}"

EVENT_KEY="${EVENT_KEY:-stock.vehicle_asset.created}"
HOOK_ID="${HOOK_ID:-}"
TIMEOUT_SEC="${TIMEOUT_SEC:-90}"
POLL_INTERVAL_SEC="${POLL_INTERVAL_SEC:-3}"
TRIGGER_TYPE="${TRIGGER_TYPE:-pipeline_smoke}"
EXPECTED_FINAL_STATUS="${EXPECTED_FINAL_STATUS:-success}"
PAYLOAD_JSON="${PAYLOAD_JSON:-{\"smoke\":true,\"event_key\":\"$EVENT_KEY\",\"source\":\"pipeline\"}}"

if ! command -v jq >/dev/null 2>&1; then
  echo "[smoke] jq not found" >&2
  exit 1
fi

if ! echo "$PAYLOAD_JSON" | jq -e . >/dev/null 2>&1; then
  echo "[smoke] PAYLOAD_JSON is not valid JSON" >&2
  exit 1
fi

hdr_auth=( -H "Authorization: Bearer $AUTH_TOKEN" )
hdr_json=( -H "Content-Type: application/json" )

curl_json() {
  curl -sS "${hdr_auth[@]}" "${hdr_json[@]}" "$@"
}

echo "[smoke] Health checks"
for endpoint in /healthz /readyz; do
  code="$(curl -sS -o /dev/null -w '%{http_code}' "$API_BASE_URL$endpoint")"
  if [[ "$code" != "200" ]]; then
    echo "[smoke] $endpoint failed with HTTP $code" >&2
    exit 1
  fi
  echo "[smoke] $endpoint -> 200"
done

if [[ -z "$HOOK_ID" ]]; then
  echo "[smoke] Resolving hook for event_key=$EVENT_KEY"
  resp="$(curl -sS "${hdr_auth[@]}" "$API_BASE_URL/admin/integrations/hooks?event_key=$EVENT_KEY&status=active&page=1&limit=1")"
  HOOK_ID="$(echo "$resp" | jq -r '.data.items[0].id // empty')"
  if [[ -z "$HOOK_ID" ]]; then
    echo "[smoke] No active hook found for event_key=$EVENT_KEY" >&2
    echo "$resp" | jq . >&2 || true
    exit 1
  fi
fi

echo "[smoke] Using hook_id=$HOOK_ID"

idem="smoke-$(date +%s)-$RANDOM"
exec_req="$(jq -n --arg tt "$TRIGGER_TYPE" --arg ik "$idem" --argjson payload "$PAYLOAD_JSON" '{trigger_type:$tt,idempotency_key:$ik,payload:$payload}')"
exec_resp="$(curl_json -X POST "$API_BASE_URL/admin/integrations/hooks/$HOOK_ID/execute" -d "$exec_req")"
execution_id="$(echo "$exec_resp" | jq -r '.data.id // empty')"

if [[ -z "$execution_id" ]]; then
  echo "[smoke] Failed to enqueue execution" >&2
  echo "$exec_resp" | jq . >&2 || true
  exit 1
fi

echo "[smoke] Execution created: $execution_id"

deadline=$(( $(date +%s) + TIMEOUT_SEC ))
while true; do
  now="$(date +%s)"
  if (( now > deadline )); then
    echo "[smoke] Timeout waiting for execution completion" >&2
    exit 1
  fi

  detail="$(curl -sS "${hdr_auth[@]}" "$API_BASE_URL/admin/integrations/executions/$execution_id")"
  status="$(echo "$detail" | jq -r '.data.status // empty')"

  if [[ "$status" == "success" || "$status" == "failed" ]]; then
    echo "[smoke] Final status=$status"
    echo "$detail" | jq '{id:.data.id,status:.data.status,trigger_type:.data.trigger_type,error_code:.data.error_code,error_message:.data.error_message,duration_ms:.data.duration_ms}'
    if [[ "$status" != "$EXPECTED_FINAL_STATUS" ]]; then
      echo "[smoke] Expected status=$EXPECTED_FINAL_STATUS but got $status" >&2
      exit 1
    fi
    echo "[smoke] PASS"
    exit 0
  fi

  echo "[smoke] Waiting... status=$status"
  sleep "$POLL_INTERVAL_SEC"
done
