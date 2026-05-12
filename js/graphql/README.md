# Zone01 Profile Dashboard

The Zone01 Profile Dashboard is a Go + vanilla JavaScript web application that lets Zone01 learners authenticate with their platform credentials and explore a personalized snapshot of their progress, collaborators, and XP history. A production deployment is available at https://graphql-y217.onrender.com.

## Table of Contents
- [Key Features](#key-features)
- [Application Architecture](#application-architecture)
- [Directory Layout](#directory-layout)
- [Backend Overview](#backend-overview)
- [Frontend Overview](#frontend-overview)
- [REST Endpoints](#rest-endpoints)
- [Local Development](#local-development)
- [Configuration](#configuration)
- [Deployment Notes](#deployment-notes)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)

## Key Features
- **Zone01 authentication** – Sign in with your Zone01 identifier/password; successful logins persist a secure session cookie.
- **GraphQL-powered insights** – Aggregates multiple GraphQL queries (profile, XP, project progress, collaborators) against the Zone01 Hasura API.
- **Interactive profile dashboard** – Rich UI that visualizes total XP, per-day XP progress, completed projects, and collaborator frequency.
- **Vanilla JS SPA routing** – Lightweight client-side router keeps navigation fast and dependency-free.
- **Graceful session handling** – Detects expired sessions, surfaces actionable status messages, and guides users back to the login view.

## Application Architecture
```
main.go (Go 1.21 HTTP server)
 ├─ router/               : request routing
 │   └─ router.go         : wires static assets + JSON API handlers
 ├─ handlers/             : HTTP handlers for HTML, auth, profile
 │   └─ handlers.go
 ├─ graphql/              : outbound Zone01 GraphQL client + queries
 │   ├─ graphql_Client.go
 │   └─ queries.go
 ├─ session/              : JWT parsing helpers for cookie sessions
 ├─ response/             : Consistent JSON response wrapper
 └─ static/               : ES modules, charts, UI helpers, CSS
      ├─ js/
      └─ css/
```

High-level flow:
1. `main.go` starts an HTTP server on port `8080` and delegates routing to `router.New()`.
2. The router serves static assets under `/static` and exposes three JSON endpoints (`/api/auth/login`, `/api/auth/logout`, `/api/profile`).
3. Authentication requests call `graphql.SignIn`, storing the returned JWT in an HTTP-only cookie.
4. Profile requests reuse the JWT to execute multiple GraphQL queries and consolidate the results into one JSON payload.
5. The frontend initializes via `static/js/main.js`, handles SPA navigation, and renders login/profile views based on authentication state.

## Directory Layout
```
go.mod
main.go
router/
handlers/
graphql/
session/
response/
static/
  ├─ css/
  └─ js/
templates/
```
Key frontend folders:
- `static/js/views/` – Login and profile view renderers.
- `static/js/services/` – Fetch wrappers for auth + profile data.
- `static/js/utils/` – Data shaping helpers (formatting, aggregations).
- `static/js/charts/` – SVG chart builders for XP progress and collaborators.
- `static/js/ui/` – DOM utilities and layout mode toggles.

## Backend Overview
- Written in Go 1.21, depends only on the standard library.
- Uses `text/template` to serve a single page from `templates/index.html`.
- Persists JWT tokens in an `HttpOnly`, `SameSite=Lax` cookie named `zone01_session`.
- Executes outbound requests with `net/http` against `https://platform.zone01.gr/api`.
- GraphQL helpers automatically detect API errors and surface them with appropriate HTTP status codes.

## Frontend Overview
- Vanilla ES modules loaded directly from `static/js/main.js`.
- Lightweight router guards `/login` and `/profile` routes while respecting browser history.
- Centralized auth state lives in `static/js/state.js`; views react to state changes and status messages.
- Charts are rendered using raw SVG to avoid bundlers or charting dependencies.
- Styling lives in `static/css/styles.css`, tuned for the two layout modes (`mode-login`, `mode-profile`).

## REST Endpoints
| Method | Path              | Description |
|--------|-------------------|-------------|
| POST   | `/api/auth/login` | Authenticates a Zone01 identifier/password pair, stores the JWT cookie, and returns `user_id`.
| POST   | `/api/auth/logout`| Clears the session cookie and returns a success message.
| POST   | `/api/profile`    | Requires a valid session cookie; fetches user profile, XP totals, project history, and collaborators via GraphQL.

Static assets are available under `/static/*`, and the SPA shell is served at `/` (any unknown route falls back to the frontend router).

## Local Development
1. **Install prerequisites**
   - Go 1.21+
   - A modern browser (for ES module support)
2. **Clone the repository**
   ```bash
   git clone <your-fork-url>
   cd graphql
   ```
3. **Run the development server**
   ```bash
   go run .
   ```
   The server starts on `http://localhost:8080`. Static assets are served directly—no bundler or Node.js toolchain is required.
4. **Sign in with Zone01 credentials**
   - Use your Zone01 identifier and password on the login screen.
   - Successful login redirects to `/profile` and loads dashboard data from the live Zone01 API.

## Configuration
The `graphql` package currently targets the production Zone01 platform (`https://platform.zone01.gr/api`). If you need to point at a different environment:
- Update the `defaultBaseAPI` constant in `graphql/graphql_Client.go`.
- Ensure CORS and TLS requirements for your alternate endpoint mirror the platform defaults.

All session handling is cookie-based; there are no additional environment variables to configure.

## Deployment Notes
- The production deployment at https://graphql-y217.onrender.com runs the same Go binary served here.
- Render (and similar platforms) can build and run the project with:
  ```bash
  go build -o server ./...
  ./server
  ```
  Make sure the deploy environment exposes port `8080` (or set `$PORT` and adapt `ListenAddr`).
- Because the app reaches out to the live Zone01 API, outbound internet access must remain enabled in your hosting provider.

## Troubleshooting
- **Invalid credentials** – The Zone01 API returns `401 Unauthorized`; ensure the identifier/password is correct and that the account has platform access.
- **Session expired** – Re-login; the frontend surfaces clear status messages when cookies expire or are missing.
- **GraphQL errors** – API validation issues bubble up as `400 Bad Request` with a generic message; inspect server logs for detailed responses when debugging.
- **Network failures** – Frontend displays friendly fallback messages when `/api/profile` cannot be reached; check server logs and outbound network connectivity.

## Contributing
1. Fork the repository and create a feature branch.
2. Make your changes, keeping Go code formatted (`gofmt`) and JS lint-friendly.
3. Validate endpoints manually (`/api/auth/login`, `/api/profile`), then open a pull request describing changes and test coverage.

Feel free to open issues for feature requests, bug reports, or clarifications about the Zone01 API integration.
