#!/usr/bin/env bash
set -euo pipefail

# Post-deploy smoke check for integrations scheduled pull pipeline.
# Requires: curl, jq

: "${API_BASE_URL:?API_BASE_URL is required (example: http://localhost:8080)}"
: "${AUTH_TOKEN:?AUTH_TOKEN is required (Bearer token with integrations read/create/execute permissions)}"

SCHEDULE_ID="${SCHEDULE_ID:-}"
SCHEDULE_KEY="${SCHEDULE_KEY:-}"
TIMEOUT_SEC="${TIMEOUT_SEC:-90}"
POLL_INTERVAL_SEC="${POLL_INTERVAL_SEC:-3}"
DRY_RUN="${DRY_RUN:-true}"
EXPECTED_FINAL_STATUS="${EXPECTED_FINAL_STATUS:-success}"
TEST_PAYLOAD_JSON="${TEST_PAYLOAD_JSON:-{\"smoke\":true,\"source\":\"pipeline\"}}"

if ! command -v jq >/dev/null 2>&1; then
  echo "[pull-smoke] jq not found" >&2
  exit 1
fi

if ! echo "$TEST_PAYLOAD_JSON" | jq -e . >/dev/null 2>&1; then
  echo "[pull-smoke] TEST_PAYLOAD_JSON is not valid JSON" >&2
  exit 1
fi

hdr_auth=( -H "Authorization: Bearer $AUTH_TOKEN" )
hdr_json=( -H "Content-Type: application/json" )

curl_json() {
  curl -sS "${hdr_auth[@]}" "${hdr_json[@]}" "$@"
}

echo "[pull-smoke] Health checks"
for endpoint in /healthz /readyz; do
  code="$(curl -sS -o /dev/null -w '%{http_code}' "$API_BASE_URL$endpoint")"
  if [[ "$code" != "200" ]]; then
    echo "[pull-smoke] $endpoint failed with HTTP $code" >&2
    exit 1
  fi
  echo "[pull-smoke] $endpoint -> 200"
done

if [[ -z "$SCHEDULE_ID" ]]; then
  if [[ -n "$SCHEDULE_KEY" ]]; then
    echo "[pull-smoke] Resolving active schedule by key=$SCHEDULE_KEY"
    resp="$(curl -sS "${hdr_auth[@]}" "$API_BASE_URL/admin/integrations/pull-schedules?q=$SCHEDULE_KEY&status=active&page=1&limit=50")"
    SCHEDULE_ID="$(echo "$resp" | jq -r --arg key "$SCHEDULE_KEY" '.data.items[] | select(.key == $key) | .id' | head -n1)"
  else
    echo "[pull-smoke] Resolving first active schedule"
    resp="$(curl -sS "${hdr_auth[@]}" "$API_BASE_URL/admin/integrations/pull-schedules?status=active&page=1&limit=1")"
    SCHEDULE_ID="$(echo "$resp" | jq -r '.data.items[0].id // empty')"
  fi
fi

if [[ -z "$SCHEDULE_ID" ]]; then
  echo "[pull-smoke] No active pull schedule found" >&2
  exit 1
fi
echo "[pull-smoke] Using schedule_id=$SCHEDULE_ID"

test_req="$(jq -n --argjson payload "$TEST_PAYLOAD_JSON" --argjson dry "$DRY_RUN" '{payload:$payload,dry_run:$dry}')"
test_resp="$(curl_json -X POST "$API_BASE_URL/admin/integrations/pull-schedules/$SCHEDULE_ID/test-run" -d "$test_req")"
run_id="$(echo "$test_resp" | jq -r '.data.run_id // empty')"

if [[ -z "$run_id" ]]; then
  echo "[pull-smoke] test-run did not return run_id" >&2
  echo "$test_resp" | jq . >&2 || true
  exit 1
fi
echo "[pull-smoke] Run created: $run_id"

deadline=$(( $(date +%s) + TIMEOUT_SEC ))
while true; do
  now="$(date +%s)"
  if (( now > deadline )); then
    echo "[pull-smoke] Timeout waiting for pull run completion" >&2
    exit 1
  fi

  detail="$(curl -sS "${hdr_auth[@]}" "$API_BASE_URL/admin/integrations/pull-runs/$run_id")"
  status="$(echo "$detail" | jq -r '.data.status // empty')"

  if [[ "$status" == "success" || "$status" == "failed" || "$status" == "cancelled" ]]; then
    echo "[pull-smoke] Final status=$status"
    echo "$detail" | jq '{id:.data.id,status:.data.status,run_key:.data.run_key,error_code:.data.error_code,error_message:.data.error_message,records_fetched:.data.records_fetched,records_upserted:.data.records_upserted,duration_ms:.data.duration_ms}'
    if [[ "$status" != "$EXPECTED_FINAL_STATUS" ]]; then
      echo "[pull-smoke] Expected status=$EXPECTED_FINAL_STATUS but got $status" >&2
      exit 1
    fi
    echo "[pull-smoke] PASS"
    exit 0
  fi

  echo "[pull-smoke] Waiting... status=$status"
  sleep "$POLL_INTERVAL_SEC"
done
