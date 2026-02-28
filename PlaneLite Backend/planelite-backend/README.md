# Plane-Lite Backend

Production-ready Go backend for a multi-workspace project management system (Plane/Linear/Jira-like).

## Tech stack

- **Language:** Go 1.22+
- **Database:** MongoDB Atlas (driver: `go.mongodb.org/mongo-driver`)
- **Auth:** JWT (access token), bcrypt for passwords
- **Config:** env vars + `.env` (godotenv)
- **HTTP:** `net/http` (no framework)

## Setup

1. **Clone and enter the project**
   ```bash
   cd planelite-backend
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Environment**
   Create a `.env` in the project root:
   ```env
   PORT=8080
   MONGO_URI=mongodb+srv://user:pass@cluster.mongodb.net/
   DB_NAME=planelite
   JWT_SECRET=your-secret-min-32-chars
   JWT_EXPIRY_HOURS=24
   ```
   `MONGO_URI` and `JWT_SECRET` are required; the server will exit on startup if they are missing.

4. **Run**
   ```bash
   go run ./cmd/server
   ```
   You should see: `Server running on :8080` (or whatever `PORT` is in `.env`).

## Quick test (is it working?)

1. **Health check** (no auth; pings MongoDB):
   ```bash
   curl http://localhost:8080/health
   ```
   Expected: `OK` and status 200. If MongoDB is unreachable, returns 503 so load balancers can failover.

2. **Signup** (creates a user in MongoDB):
   ```bash
   curl -X POST http://localhost:8080/auth/signup \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"secret123","role":"USER"}'
   ```
   Expected: JSON with `access_token`, `user_id`, `email`, `role`. The user is stored in the **`users`** collection in MongoDB (in the database named `DB_NAME`, e.g. `planelite`).

3. **Login** (same user, get a new token):
   ```bash
   curl -X POST http://localhost:8080/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"secret123"}'
   ```
   Expected: JSON with `access_token`, `user_id`, `email`, `role`.

4. **Authenticated route** (use the token from signup/login):
   ```bash
   curl http://localhost:8080/me -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
   ```
   Replace `YOUR_ACCESS_TOKEN` with the `access_token` from the signup or login response. Expected: your user info (id, email, role, no password).

**Yes, your user is added to MongoDB.** Signup calls the user repository’s `Create`, which does `InsertOne` into the **`users`** collection. You can confirm in MongoDB Atlas: open the `planelite` (or your `DB_NAME`) database and check the **`users`** collection after a successful signup.

## API (overview)

- **Auth:** `POST /auth/signup`, `POST /auth/login` (email + password, returns JWT).
- **Me:** `GET /me` (Bearer token).
- **Workspaces:** `POST /workspaces` (admin only), `GET /workspaces`, `GET /workspaces/{id}`, `POST /workspaces/{id}/members`, `GET /workspaces/{id}/members`, `POST /workspaces/{id}/members/{mid}/approve`.
- **Projects:** `POST /workspaces/{id}/projects`, `GET /workspaces/{id}/projects`, `GET /workspaces/{id}/projects/{pid}`.
- **Tasks:** `POST /workspaces/{id}/projects/{pid}/tasks`, `GET /workspaces/{id}/projects/{pid}/tasks`, `GET/ PATCH /workspaces/{id}/projects/{pid}/tasks/{tid}`.

Roles: `ADMIN`, `PROJECT_MANAGER`, `USER`. Only ADMIN can create workspaces; only PROJECT_MANAGER (or ADMIN) can create tasks; users can update task status/priority.

## Architecture

- **Handler → Service → Repository** per domain (auth, user, workspace, project, task).
- Business rules in services; repositories only talk to MongoDB; handlers only parse request/response.
- Auth middleware validates JWT and sets user in context; role and workspace-access middleware enforce permissions.
- Indexes: `users.email` (unique), `memberships (user_id, workspace_id)` (unique).

## Production-oriented behaviour

- **Graceful shutdown:** On SIGINT/SIGTERM the server stops accepting new requests, waits up to 15s for in-flight requests, then disconnects MongoDB. Set `ENV=development` for dev; omit or set to `production` for prod.
- **Deep health:** `GET /health` pings MongoDB; returns 503 if DB is down so load balancers can mark the instance unhealthy.
- **Structured logging:** Uses zap (dev: pretty; prod: JSON). Log level and format controlled by `ENV`.
- **Request ID:** Every response includes `X-Request-ID` (from header or generated) for tracing and log correlation.
- **Server timeouts:** `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, `IdleTimeout` set to avoid slow clients holding connections.
