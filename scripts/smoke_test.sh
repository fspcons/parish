#!/usr/bin/env bash
#
# Happy-path smoke tests for the Parish API.
#
# Prerequisites:
#   1. Infrastructure running        (make infra-start)
#   2. API server running            (make run-local)
#      — run-local seeds an admin user via ADMIN_EMAIL / ADMIN_PASSWORD env vars.
#
# Usage:
#   make smoke-test
#   # or directly:
#   ./scripts/smoke_test.sh [BASE_URL] [EMULATOR_HOST]

set -euo pipefail

BASE="${1:-http://localhost:8080}"
EMULATOR_HOST="${2:-localhost:8081}"
ADMIN_EMAIL="${ADMIN_EMAIL:-admin@parish.local}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-Admin@Str0ng!Pass}"
COOKIE_JAR=$(mktemp)
PASS=0
FAIL=0

# shellcheck disable=SC2329
cleanup() {
  printf "\n\033[1;34m=== Cleanup ===\033[0m\n"
  printf "  Resetting Datastore emulator at %s ...\n" "$EMULATOR_HOST"
  if curl -s -X POST "http://${EMULATOR_HOST}/reset" >/dev/null 2>&1; then
    printf "  \033[32m✓ Emulator data reset\033[0m\n"
  else
    printf "  \033[33m⚠ Could not reset emulator (is it running?)\033[0m\n"
  fi
  rm -f "$COOKIE_JAR"
}
trap cleanup EXIT

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------
section() { printf "\n\033[1;34m=== %s ===\033[0m\n" "$1"; }

# call METHOD PATH [BODY]
# Prints status code + response, sets $STATUS and $BODY.
call() {
  local method="$1" path="$2" body="${3:-}"
  local curl_args=(-s -w "\n%{http_code}" -X "$method" -b "$COOKIE_JAR" -c "$COOKIE_JAR")
  curl_args+=(-H "Content-Type: application/json")

  if [[ -n "$body" ]]; then
    curl_args+=(-d "$body")
  fi

  local raw
  raw=$(curl "${curl_args[@]}" "${BASE}${path}")
  STATUS=$(echo "$raw" | tail -n1)
  BODY=$(echo "$raw" | sed '$d')
}

assert_status() {
  local expected="$1" label="$2"
  if [[ "$STATUS" == "$expected" ]]; then
    printf "  \033[32m✓ %s  (HTTP %s)\033[0m\n" "$label" "$STATUS"
    PASS=$((PASS + 1))
  else
    printf "  \033[31m✗ %s  (expected %s, got %s)\033[0m\n" "$label" "$expected" "$STATUS"
    printf "    %s\n" "$BODY"
    FAIL=$((FAIL + 1))
  fi
}

# Extract a JSON field value (simple jq-free fallback).
json_field() {
  echo "$BODY" | grep -o "\"$1\":\"[^\"]*\"" | head -1 | cut -d'"' -f4
}

# ---------------------------------------------------------------------------
# 1. Health check
# ---------------------------------------------------------------------------
section "Health"
call GET /health
assert_status 200 "GET /health"

# ---------------------------------------------------------------------------
# 2. Auth: Login as seeded admin
# ---------------------------------------------------------------------------
section "Auth (seeded admin)"

call POST /api/auth/login "{\"email\":\"${ADMIN_EMAIL}\",\"password\":\"${ADMIN_PASSWORD}\"}"
assert_status 200 "POST /api/auth/login (seeded admin)"

call POST /api/auth/login "{\"email\":\"${ADMIN_EMAIL}\",\"password\":\"wrongpassword\"}"
assert_status 401 "POST /api/auth/login (bad password)"

# Re-login to restore valid cookie after the failed attempt.
call POST /api/auth/login "{\"email\":\"${ADMIN_EMAIL}\",\"password\":\"${ADMIN_PASSWORD}\"}"
assert_status 200 "POST /api/auth/login (restore session)"

# ---------------------------------------------------------------------------
# 3. Roles CRUD
# ---------------------------------------------------------------------------
section "Roles"

call POST /api/roles '{"name":"editor","description":"Content editor","permissions":[{"resource":"events","read":true,"write":true},{"resource":"materials","read":true,"write":true}]}'
assert_status 201 "POST /api/roles (create)"
ROLE_ID=$(json_field "id")
printf "    created role id=%s\n" "$ROLE_ID"

call GET /api/roles
assert_status 200 "GET /api/roles (list)"

call GET "/api/roles/${ROLE_ID}"
assert_status 200 "GET /api/roles/:id"

call PUT "/api/roles/${ROLE_ID}" '{"name":"editor-updated","description":"Updated editor","permissions":[{"resource":"events","read":true,"write":true},{"resource":"materials","read":true,"write":false}]}'
assert_status 200 "PUT /api/roles/:id"

call DELETE "/api/roles/${ROLE_ID}"
assert_status 200 "DELETE /api/roles/:id"

# ---------------------------------------------------------------------------
# 4. Schedule (GET is public, PUT is protected)
# ---------------------------------------------------------------------------
section "Schedule"

call GET /api/schedule
assert_status 200 "GET /api/schedule"

call PUT /api/schedule '{"monday":"7:00 AM","tuesday":"7:00 AM","wednesday":"7:00 AM","thursday":"7:00 AM","friday":"7:00 AM","saturday":"9:00 AM","sunday":"8:00 AM, 10:00 AM"}'
assert_status 200 "PUT /api/schedule"

# ---------------------------------------------------------------------------
# 5. Events CRUD
# ---------------------------------------------------------------------------
section "Events"

call GET /api/events
assert_status 200 "GET /api/events (list)"

call POST /api/events '{"title":"Sunday Mass","description":"Weekly celebration","date":"2026-03-01","location":"Main Church","origin":"parish"}'
assert_status 201 "POST /api/events"
EVENT_ID=$(json_field "id")
printf "    created event id=%s\n" "$EVENT_ID"

call GET "/api/events/${EVENT_ID}"
assert_status 200 "GET /api/events/:id"

call PUT "/api/events/${EVENT_ID}" '{"title":"Sunday Mass (Updated)","description":"Updated","date":"2026-03-01","location":"Main Church","origin":"parish"}'
assert_status 200 "PUT /api/events/:id"

call DELETE "/api/events/${EVENT_ID}"
assert_status 200 "DELETE /api/events/:id"

# ---------------------------------------------------------------------------
# 6. Parish Groups CRUD
# ---------------------------------------------------------------------------
section "Parish Groups"

call GET /api/parish-groups
assert_status 200 "GET /api/parish-groups (list)"

call POST /api/parish-groups '{"title":"Youth Group","description":"Young adults ministry","manager":"John Doe","active":true}'
assert_status 201 "POST /api/parish-groups"
GROUP_ID=$(json_field "id")
printf "    created group id=%s\n" "$GROUP_ID"

call GET "/api/parish-groups/${GROUP_ID}"
assert_status 200 "GET /api/parish-groups/:id"

call PUT "/api/parish-groups/${GROUP_ID}" '{"title":"Youth Group (Updated)","description":"Updated","manager":"Jane Doe","active":true}'
assert_status 200 "PUT /api/parish-groups/:id"

call DELETE "/api/parish-groups/${GROUP_ID}"
assert_status 200 "DELETE /api/parish-groups/:id"

# ---------------------------------------------------------------------------
# 7. Materials CRUD
# ---------------------------------------------------------------------------
section "Materials"

call GET /api/materials
assert_status 200 "GET /api/materials (list)"

call GET "/api/materials?type=videos"
assert_status 200 "GET /api/materials?type=videos"

call GET "/api/materials?label=catechism"
assert_status 200 "GET /api/materials?label=catechism"

call POST /api/materials '{"title":"Intro to Faith","type":"videos","description":"Catechism series","url":"https://example.com/video1","label":"catechism"}'
assert_status 201 "POST /api/materials"
MATERIAL_ID=$(json_field "id")
printf "    created material id=%s\n" "$MATERIAL_ID"

call GET "/api/materials/${MATERIAL_ID}"
assert_status 200 "GET /api/materials/:id"

call PUT "/api/materials/${MATERIAL_ID}" '{"title":"Intro to Faith (Updated)","type":"documents","description":"Updated","url":"https://example.com/doc1","label":"catechism"}'
assert_status 200 "PUT /api/materials/:id"

call DELETE "/api/materials/${MATERIAL_ID}"
assert_status 200 "DELETE /api/materials/:id"

# ---------------------------------------------------------------------------
# 8. Register + User role assignment
# ---------------------------------------------------------------------------
section "Register & Role Assignment"

# Create a role to assign.
call POST /api/roles '{"name":"viewer","description":"Read-only viewer","permissions":[{"resource":"events","read":true,"write":false}]}'
assert_status 201 "POST /api/roles (create viewer)"
VIEWER_ROLE_ID=$(json_field "id")

# Register a new (non-admin) user via the public register endpoint.
call POST /api/auth/register '{"email":"viewer@parish.local","name":"Viewer","password":"viewer123"}'
assert_status 201 "POST /api/auth/register"
VIEWER_USER_ID=$(json_field "id")
printf "    created viewer user id=%s\n" "$VIEWER_USER_ID"

# Assign the viewer role to the new user (admin is still logged in).
call PUT "/api/users/${VIEWER_USER_ID}/roles" "{\"roleIds\":[\"${VIEWER_ROLE_ID}\"]}"
assert_status 200 "PUT /api/users/:id/roles"

# ---------------------------------------------------------------------------
# 9. Logout
# ---------------------------------------------------------------------------
section "Logout"

call POST /api/auth/logout
assert_status 200 "POST /api/auth/logout"

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
section "Results"
TOTAL=$((PASS + FAIL))
printf "  %d / %d passed" "$PASS" "$TOTAL"
if [[ "$FAIL" -gt 0 ]]; then
  printf " (\033[31m%d failed\033[0m)" "$FAIL"
fi
printf "\n\n"

exit "$FAIL"
