#!/usr/bin/env bash
set -euo pipefail

# Post-deploy smoke check for integrations shared service observability.
# Requires: curl, jq

: "${API_BASE_URL:?API_BASE_URL is required (example: http://localhost:8080)}"
: "${AUTH_TOKEN:?AUTH_TOKEN is required (Bearer token with integrations read permissions)}"

SHARED_SERVICE_ID="${SHARED_SERVICE_ID:-}"
SHARED_SERVICE_GRANT_ID="${SHARED_SERVICE_GRANT_ID:-}"
WINDOW_DAYS="${WINDOW_DAYS:-7}"

if ! command -v jq >/dev/null 2>&1; then
  echo "[shared-smoke] jq not found" >&2
  exit 1
fi

if [[ -z "$SHARED_SERVICE_ID" || -z "$SHARED_SERVICE_GRANT_ID" ]]; then
  echo "[shared-smoke] SHARED_SERVICE_ID and SHARED_SERVICE_GRANT_ID are required" >&2
  exit 1
fi

hdr_auth=(-H "Authorization: Bearer $AUTH_TOKEN")

fetch_json() {
  local url="$1"
  curl -sS "${hdr_auth[@]}" "$url"
}

echo "[shared-smoke] Health checks"
for endpoint in /healthz /readyz; do
  code="$(curl -sS -o /dev/null -w '%{http_code}' "$API_BASE_URL$endpoint")"
  if [[ "$code" != "200" ]]; then
    echo "[shared-smoke] $endpoint failed with HTTP $code" >&2
    exit 1
  fi
done

svc_url="$API_BASE_URL/admin/integrations/shared-services/$SHARED_SERVICE_ID/usage-summary?days=$WINDOW_DAYS"
gr_url="$API_BASE_URL/admin/integrations/shared-service-grants/$SHARED_SERVICE_GRANT_ID/usage-summary?days=$WINDOW_DAYS"

echo "[shared-smoke] Checking service usage summary"
svc="$(fetch_json "$svc_url")"
if ! echo "$svc" | jq -e '.data.totals.requests_total != null and .data.outcomes != null and .data.days != null' >/dev/null; then
  echo "[shared-smoke] Service summary schema check failed" >&2
  echo "$svc" | jq . >&2 || true
  exit 1
fi
echo "$svc" | jq '{window_days:.data.window_days,service:.data.service,totals:.data.totals,warnings:(.data.warnings // [])}'

echo "[shared-smoke] Checking grant usage summary"
gr="$(fetch_json "$gr_url")"
if ! echo "$gr" | jq -e '.data.grant.id != null and .data.totals.requests_total != null and .data.days != null' >/dev/null; then
  echo "[shared-smoke] Grant summary schema check failed" >&2
  echo "$gr" | jq . >&2 || true
  exit 1
fi
echo "$gr" | jq '{window_days:.data.window_days,grant:.data.grant,effective_limits:.data.effective_limits,totals:.data.totals,warnings:(.data.warnings // [])}'

echo "[shared-smoke] PASS"
