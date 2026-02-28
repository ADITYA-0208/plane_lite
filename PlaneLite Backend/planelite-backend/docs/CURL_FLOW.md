# Curl flow: Admin → Org → Add people → Admin sees requests → Admin approves

Base URL: `http://localhost:8080` (change if your server runs elsewhere.)

---

## All admin curl requests (quick reference)

Every request that uses the admin token (`Authorization: Bearer $ADMIN_TOKEN`). Set `ADMIN_TOKEN` first (signup or login with an ADMIN user).

| # | Method | Path | Purpose |
|---|--------|------|--------|
| 1 | POST | `/auth/signup` | Create admin (or login if exists) |
| 2 | POST | `/auth/login` | Get admin token when admin already exists |
| 3 | POST | `/workspaces` | Create org (admin only; one workspace per admin) |
| 4 | GET | `/workspaces` | List my workspaces (admin only) |
| 5 | GET | `/workspaces/{id}` | Get workspace by id (need workspace access) |
| 6 | POST | `/workspaces/{id}/members` | Add user to org (by `user_id`) |
| 7 | GET | `/workspaces/{id}/members` | List members (pending + approved) |
| 8 | POST | `/workspaces/{id}/members/{mid}/approve` | Approve a pending membership |

**Copy-paste (admin token required for 3–8):**

```bash
BASE=http://localhost:8080

# 1. Get admin token (signup or login)
export ADMIN_TOKEN=$(curl -s -X POST $BASE/auth/login -H "Content-Type: application/json" \
  -d '{"email":"admin@org.com","password":"admin123"}' | jq -r '.data.access_token')

# 2. Create org
curl -s -X POST $BASE/workspaces -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" -d '{"name":"My Org"}'

# 3. List my workspaces   1
curl -s $BASE/workspaces -H "Authorization: Bearer $ADMIN_TOKEN"

# 4. Get workspace (set WORKSPACE_ID first)  2 
curl -s $BASE/workspaces/$WORKSPACE_ID -H "Authorization: Bearer $ADMIN_TOKEN"

# 5. Add member (set WORKSPACE_ID and USER_ID)
curl -s -X POST "$BASE/workspaces/$WORKSPACE_ID/members" -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" -d "$(jq -n --arg uid "$USER_ID" '{user_id: $uid}')"

# 6. List members
curl -s "$BASE/workspaces/$WORKSPACE_ID/members" -H "Authorization: Bearer $ADMIN_TOKEN"

# 7. Approve membership (set MEMBERSHIP_ID)
curl -s -X POST "$BASE/workspaces/$WORKSPACE_ID/members/$MEMBERSHIP_ID/approve" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

---

## 1. Create Admin (signup)

```bash
curl -s -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@org.com","password":"admin123","role":"ADMIN"}'
```

**Save from response:** `access_token` → use as `ADMIN_TOKEN`, and `user_id` (admin’s id).

Example (save the token):
```bash
# If you have jq:
export ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@org.com","password":"admin123","role":"ADMIN"}' | jq -r '.data.access_token')
```

---

## 2. Admin creates org (workspace)

```bash
curl -s -X POST http://localhost:8080/workspaces \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"name":"My Org"}'
```

**Save from response:** `data._id` → use as `WORKSPACE_ID`.

**If you get `{"error":"unauthorized"}`:** The server only accepts a valid JWT. Ensure the token is set: `echo "$ADMIN_TOKEN"` (should print a long string, not empty). If you exported using signup and the email already existed, the response was `{"error":"conflict (e.g. duplicate)"}` so `.data.access_token` is missing and `ADMIN_TOKEN` is empty. Use **login** to get a token: `export ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/auth/login -H "Content-Type: application/json" -d '{"email":"admin@org.com","password":"admin123"}' | jq -r '.data.access_token')` (use the same email/password you signed up with).

Example:
```bash
export WORKSPACE_ID=$(curl -s -X POST http://localhost:8080/workspaces \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"name":"My Org"}' | jq -r '.data._id // .data.ID')
```

If `echo "$WORKSPACE_ID"` shows `null`, inspect the raw response: run the same `curl` (without the `| jq ...` pipe) and check for `{"error":"..."}`. Fix the token or role; the success body has the workspace id under `data._id` or `data.ID`.

---

## 3. Create Project Manager and User (“people”)

**Project Manager:**
```bash
curl -s -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"pm@org.com","password":"pm123","role":"PROJECT_MANAGER"}'
```

**Save** `user_id` from response → use as `PM_USER_ID`.

**Regular User:**
```bash
curl -s -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"dev@org.com","password":"dev123","role":"USER"}'
```

**Save** `user_id` from response → use as `USER_ID`.

Example (with jq):
```bash
export PM_USER_ID=$(curl -s -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"pm@org.com","password":"pm123","role":"PROJECT_MANAGER"}' | jq -r '.data.user_id')

export USER_ID=$(curl -s -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"dev@org.com","password":"dev123","role":"USER"}' | jq -r '.data.user_id')
```

If signup returns `{"error":"conflict (e.g. duplicate)"}`, there is no `.data.user_id`. Use **login** to get the existing user’s ID:  
`export USER_ID=$(curl -s -X POST http://localhost:8080/auth/login -H "Content-Type: application/json" -d '{"email":"dev@org.com","password":"dev123"}' | jq -r '.data.user_id')` (and similarly for PM with `pm@org.com` / `pm123`).

---

## 4. Admin adds people to the org (they become PENDING)

**Add Project Manager to org:**
```bash
curl -s -X POST "http://localhost:8080/workspaces/$WORKSPACE_ID/members" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d "{\"user_id\":\"$PM_USER_ID\"}"
```

**Add User to org:**
```bash
curl -s -X POST "http://localhost:8080/workspaces/$WORKSPACE_ID/members" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d "{\"user_id\":\"$USER_ID\"}"
```

If you get `{"error":"bad request"}`: the API expects `user_id` and the workspace path `id` to be valid 24‑character hex IDs. Check `echo "$WORKSPACE_ID"` and `echo "$USER_ID"` (both must be non‑empty). If the user already existed, set `USER_ID` from login: `export USER_ID=$(curl -s -X POST http://localhost:8080/auth/login -H "Content-Type: application/json" -d '{"email":"dev@org.com","password":"dev123"}' | jq -r '.data.user_id')`. You can also build the body with jq: `-d "$(jq -n --arg uid "$USER_ID" '{user_id: $uid}')"`.

Each response has `data._id` → that is the **membership id** (`MEMBERSHIP_ID`) for the next step.

---

## 5. Admin sees all requests (pending + approved members)

```bash
curl -s "http://localhost:8080/workspaces/$WORKSPACE_ID/members" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

You’ll see each membership with `status`: `"PENDING"` or `"APPROVED"`.  
**Save** the `_id` of any PENDING membership to approve it.

Example (get first PENDING membership id with jq):
```bash
export MEMBERSHIP_ID=$(curl -s "http://localhost:8080/workspaces/$WORKSPACE_ID/members" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq -r '.data[0]._id')
```

---

## 6. Admin approves a request (membership → APPROVED)

```bash
curl -s -X POST "http://localhost:8080/workspaces/$WORKSPACE_ID/members/$MEMBERSHIP_ID/approve" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

No body. Success = 204 No Content.

Repeat with another `MEMBERSHIP_ID` to approve more people.

---

## 7. (Optional) Admin sees members again (all APPROVED after approval)

```bash
curl -s "http://localhost:8080/workspaces/$WORKSPACE_ID/members" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

---

## One-shot copy-paste (bash, with jq)

Run these in order. Requires `jq`.

```bash
BASE=http://localhost:8080

# 1. Admin
ADMIN_TOKEN=$(curl -s -X POST $BASE/auth/signup -H "Content-Type: application/json" \
  -d '{"email":"admin@org.com","password":"admin123","role":"ADMIN"}' | jq -r '.data.access_token')

# 2. Org
WORKSPACE_ID=$(curl -s -X POST $BASE/workspaces -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" -d '{"name":"My Org"}' | jq -r '.data._id // .data.ID')

# 3. People (PM + User)
PM_USER_ID=$(curl -s -X POST $BASE/auth/signup -H "Content-Type: application/json" \
  -d '{"email":"pm@org.com","password":"pm123","role":"PROJECT_MANAGER"}' | jq -r '.data.user_id')
USER_ID=$(curl -s -X POST $BASE/auth/signup -H "Content-Type: application/json" \
  -d '{"email":"dev@org.com","password":"dev123","role":"USER"}' | jq -r '.data.user_id')

# 4. Admin adds people to org (PENDING)
curl -s -X POST "$BASE/workspaces/$WORKSPACE_ID/members" -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" -d "{\"user_id\":\"$PM_USER_ID\"}"
curl -s -X POST "$BASE/workspaces/$WORKSPACE_ID/members" -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" -d "{\"user_id\":\"$USER_ID\"}"

# 5. Admin sees requests (list members; status = PENDING or APPROVED)
curl -s "$BASE/workspaces/$WORKSPACE_ID/members" -H "Authorization: Bearer $ADMIN_TOKEN" | jq .

# 6. Approve first pending (replace MEMBERSHIP_ID with _id from step 5)
MEMBERSHIP_ID=$(curl -s "$BASE/workspaces/$WORKSPACE_ID/members" -H "Authorization: Bearer $ADMIN_TOKEN" | jq -r '.data[0]._id')
curl -s -X POST "$BASE/workspaces/$WORKSPACE_ID/members/$MEMBERSHIP_ID/approve" \
  -H "Authorization: Bearer $ADMIN_TOKEN" -w "\nHTTP %{http_code}\n"

# 7. Admin sees members again (approved)
curl -s "$BASE/workspaces/$WORKSPACE_ID/members" -H "Authorization: Bearer $ADMIN_TOKEN" | jq .
```

---

## Summary

| Step | Who        | Action | Curl |
|------|------------|--------|------|
| 1    | -          | Add admin | `POST /auth/signup` with `"role":"ADMIN"` |
| 2    | Admin      | Create org | `POST /workspaces` with `Authorization: Bearer $ADMIN_TOKEN` |
| 3    | -          | Add PM / User | `POST /auth/signup` with `role` PROJECT_MANAGER or USER |
| 4    | Admin      | Add people to org | `POST /workspaces/{id}/members` with `user_id` |
| 5    | Admin      | See requests | `GET /workspaces/{id}/members` (check `status`: PENDING) |
| 6    | Admin      | Approve request | `POST /workspaces/{id}/members/{mid}/approve` |

After approval, that user has access to the org (workspace) and can call workspace-scoped APIs with their own token.
