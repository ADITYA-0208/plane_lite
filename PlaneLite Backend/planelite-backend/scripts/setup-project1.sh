#!/usr/bin/env bash
# Create admin adityakhanna, org "project1", user aarohi, and add aarohi to project1.
# Requires: server running on localhost:8080, jq installed.

set -e
BASE="${BASE:-http://localhost:8080}"

echo "1. Creating admin adityakhanna..."
ADMIN_RESP=$(curl -s -X POST "$BASE/auth/signup" \
  -H "Content-Type: application/json" \
  -d '{"email":"adityakhanna@org.com","password":"admin123","role":"ADMIN"}')
if echo "$ADMIN_RESP" | jq -e '.error' >/dev/null 2>&1; then
  echo "   Admin may already exist; logging in..."
  ADMIN_RESP=$(curl -s -X POST "$BASE/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"adityakhanna@org.com","password":"admin123"}')
fi
export ADMIN_TOKEN=$(echo "$ADMIN_RESP" | jq -r '.data.access_token')
if [ -z "$ADMIN_TOKEN" ] || [ "$ADMIN_TOKEN" = "null" ]; then
  echo "   Failed to get admin token. Response: $ADMIN_RESP"
  exit 1
fi
echo "   Admin token set."

echo "2. Creating org project1 (or using existing)..."
WS_RESP=$(curl -s -X POST "$BASE/workspaces" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"name":"project1"}')
if echo "$WS_RESP" | jq -e '.error' >/dev/null 2>&1; then
  if echo "$WS_RESP" | jq -e '.error == "conflict (e.g. duplicate)"' >/dev/null 2>&1; then
    echo "   Admin already has a workspace; fetching existing..."
    WS_LIST=$(curl -s -X GET "$BASE/workspaces" -H "Authorization: Bearer $ADMIN_TOKEN")
    export WORKSPACE_ID=$(echo "$WS_LIST" | jq -r '.data[0]._id // .data[0].ID')
  else
    echo "   Failed to create workspace. Response: $WS_RESP"
    exit 1
  fi
else
  export WORKSPACE_ID=$(echo "$WS_RESP" | jq -r '.data._id // .data.ID')
fi
if [ -z "$WORKSPACE_ID" ] || [ "$WORKSPACE_ID" = "null" ] || [ "$WORKSPACE_ID" = "000000000000000000000000" ]; then
  echo "   Create returned zero ID (server may be running old code). Trying GET /workspaces..."
  WS_LIST=$(curl -s -X GET "$BASE/workspaces" -H "Authorization: Bearer $ADMIN_TOKEN")
  WORKSPACE_ID=$(echo "$WS_LIST" | jq -r '.data[0]._id // .data[0].ID')
fi
if [ -z "$WORKSPACE_ID" ] || [ "$WORKSPACE_ID" = "null" ] || [ "$WORKSPACE_ID" = "000000000000000000000000" ]; then
  echo "   Failed to get a valid workspace id."
  echo "   Do this:"
  echo "   1. Restart the backend server (so workspace create returns a real ID)."
  echo "   2. In MongoDB, delete the bad workspace: db.workspaces.deleteOne({ _id: ObjectId(\"000000000000000000000000\") })"
  echo "   3. Re-run this script."
  echo "   Response was: $WS_RESP"
  exit 1
fi
export WORKSPACE_ID
echo "   WORKSPACE_ID=$WORKSPACE_ID"

echo "3. Creating user aarohi..."
USER_RESP=$(curl -s -X POST "$BASE/auth/signup" \
  -H "Content-Type: application/json" \
  -d '{"email":"aarohi@org.com","password":"user123","role":"USER"}')
if echo "$USER_RESP" | jq -e '.error' >/dev/null 2>&1; then
  echo "   User may already exist; getting user_id via login..."
  USER_RESP=$(curl -s -X POST "$BASE/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"aarohi@org.com","password":"user123"}')
fi
export USER_ID=$(echo "$USER_RESP" | jq -r '.data.user_id')
if [ -z "$USER_ID" ] || [ "$USER_ID" = "null" ]; then
  echo "   Failed to get user id. Response: $USER_RESP"
  exit 1
fi
echo "   USER_ID=$USER_ID"

echo "4. Adding aarohi to project1..."
ADD_RESP=$(curl -s -X POST "$BASE/workspaces/$WORKSPACE_ID/members" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d "$(jq -n --arg uid "$USER_ID" '{user_id: $uid}')")
if echo "$ADD_RESP" | jq -e '.error' >/dev/null 2>&1; then
  echo "   Failed to add member. Response: $ADD_RESP"
  exit 1
fi
echo "   Done. aarohi added to project1 (PENDING)."

echo "5. Listing members (requests)..."
MEMBERS_RESP=$(curl -s "$BASE/workspaces/$WORKSPACE_ID/members" -H "Authorization: Bearer $ADMIN_TOKEN")
echo "$MEMBERS_RESP" | jq .
echo ""

# Approve all PENDING memberships
PENDING_IDS=$(echo "$MEMBERS_RESP" | jq -r '.data[]? | select(.status == "PENDING") | ._id // .ID' 2>/dev/null)
if [ -z "$PENDING_IDS" ]; then
  echo "6. No PENDING requests to approve."
else
  echo "6. Approving PENDING request(s)..."
  while IFS= read -r mid; do
    [ -z "$mid" ] && continue
    APPROVE_RESP=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/workspaces/$WORKSPACE_ID/members/$mid/approve" \
      -H "Authorization: Bearer $ADMIN_TOKEN")
    if [ "$APPROVE_RESP" = "204" ]; then
      echo "   Approved membership $mid"
    else
      echo "   Failed to approve $mid (HTTP $APPROVE_RESP)"
    fi
  done <<< "$PENDING_IDS"
fi

echo ""
echo "Summary:"
echo "  Admin:     adityakhanna@org.com (ADMIN_TOKEN in env)"
echo "  Workspace: project1 (WORKSPACE_ID=$WORKSPACE_ID)"
echo "  User:      aarohi@org.com (USER_ID=$USER_ID)"
