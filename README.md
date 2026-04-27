# ElevatEd Backend

A production-quality REST API backend for an education platform, built in **Go** with a clean layered architecture, custom-built authentication system, PostgreSQL, Docker, and Tailscale for secure remote access.

> **Stack:** Go · PostgreSQL · Fiber · Docker · Tailscale · Argon2id · JWT

---

## Authentication System

Auth was built from scratch rather than delegating to a third-party service. The system implements a **dual-token architecture** with short-lived JWTs and long-lived refresh tokens backed by a persistent session store.

### How It Works

```
Login Request
    │
    ├─ Argon2id hash comparison against stored password hash
    │
    ├─ JWT issued (3hr expiry) — custom claims with embedded user details
    │
    └─ Refresh token issued (24hr expiry)
           │
           └─ Session record written to USER_SESSIONS
                  (token hash, user ID, IP address, user agent, timestamps)

Subsequent Requests
    │
    ├─ Bearer token extracted + validated via regex
    ├─ JWT signature verified (HMAC-SHA256)
    ├─ Expiry checked
    └─ Claims decoded → user context available to handlers

Token Refresh
    │
    ├─ Refresh token looked up in USER_SESSIONS
    ├─ Hash compared, expiry checked, IP recorded
    └─ New JWT issued

Logout
    └─ Session expires_at set to now → token immediately invalidated
```

### Security Decisions

**Argon2id for password hashing** — chosen over bcrypt/scrypt as the winner of the Password Hashing Competition, with configurable memory (64MB), parallelism (4 threads), and iteration parameters tuned for the deployment environment.

**Embedded user details in JWT claims** — avoids a database lookup on every authenticated request. A custom `CustomClaimStruct` wraps `jwt.RegisteredClaims` with a full `UserDetails` payload, so handlers have immediate access to user context without hitting PostgreSQL.

**Session table for refresh tokens** — refresh tokens are not self-contained. They're validated against `USER_SESSIONS`, which means logout is real: setting `expires_at = now()` immediately invalidates the session server-side with no waiting for token expiry.

**IP address tracking** — each session records the originating IP, enabling anomaly detection and audit trails.

---

## Architecture

The codebase follows a clean layered separation:

```
ElevatEd_backend/
├── main.go              # Entry point — wires dependencies, starts server
├── router/              # Route definitions + middleware (auth, CORS, logging)
├── handler/             # HTTP handlers — parse requests, call functions, return responses
├── functions/           # Business logic — auth, queries, transformations
├── structs/             # Shared data models — request/response/DB entities
├── database/            # PostgreSQL connection + schema migrations
├── docs/                # API endpoint documentation
├── .air.toml            # Hot reload config (Air)
├── Dockerfile           # Container config
└── apitest.http         # In-editor API test file
```

**Request flow:**
```
HTTP Request → Router (middleware) → Handler → Functions → Database → Response
```

Each layer has a single responsibility. Handlers don't touch the database directly — all queries go through `functions/`, keeping business logic testable and handlers thin.

---

## API Modules

Full endpoint documentation lives in `/docs`:

| Module | Docs | Description |
|--------|------|-------------|
| Auth | [authEndPoints.md](docs/authEndPoints.md) | Register, login, logout, token refresh |
| Appointments | [appointmentEndPoints.md](docs/appointmentEndPoints.md) | Schedule, view, update, cancel |
| Events | [eventEndPoints.md](docs/eventEndPoints.md) | Create, manage, filter educational events |
| Grades | [gradeEndPoint.md](docs/gradeEndPoint.md) | Submit, retrieve, track student grades |

---

## Infrastructure

### Tailscale
Used for two purposes: secure remote access to the development environment, and encrypted access to the PostgreSQL database without exposing it to the public internet. The database is only reachable over the Tailscale network — no open ports.

### Docker
The backend is fully containerized with live reload support via Air:

```bash
# Build
docker build -t elevated_backend .

# Run with live reload
docker run -p 3000:3000 --rm \
  -v $(pwd):/app \
  -v /app/tmp \
  --name docker-air \
  elevated_backend
```

---

## Getting Started

### Prerequisites
- Go 1.24.5+
- PostgreSQL
- Tailscale (for database access)

### Setup

```bash
git clone https://github.com/ElevatEd-Devs/ElevatEd_backend.git
cd ElevatEd_backend
go mod download
```

Create a `.env` file:
```
DATABASE_URL=postgresql://user:password@localhost:5432/elevated_db
JWT_SECRET=your_jwt_secret
HASHING_SALT=your_argon2_salt
PORT=3000
```

### Run

```bash
# With live reload (recommended)
air

# Standard
go run main.go
```

### Ubuntu: Air PATH Setup
```bash
export PATH=$PATH:$HOME/go/bin
source ~/.bashrc
air
```

---

## API Testing

The repo includes `apitest.http` for in-editor testing with the VS Code REST Client extension. Alternatively use Postman or Bruno.

```bash
curl -X GET http://localhost:3000/api/v1/health
```

---

## Related

- [ElevatEd Frontend](https://github.com/ElevatEd-Devs/ElevatEd_frontend)
