# JavaScript Architecture Guide

This document walks through the browser-side code under `static/js/`, explaining how each module contributes to the login → profile experience. It is written to be friendly to readers who are new to JavaScript, so the focus is on what happens, why it happens, and how the different pieces fit together.

## Runtime Flow Overview
- **HTML bootstrapping**: `templates/index.html` loads `static/js/main.js` as an ES module, which means the browser resolves the `import` statements and executes the module after the DOM is ready.
- **App startup**: `main.js` wires up the router with the two view renderers (login + profile), probes the server to learn the current authentication state, and then asks the router to render the correct screen for the current URL.
- **Routing**: `router.js` watches navigation events (initial load, `history.pushState`, back/forward buttons). It decides whether the login or profile view should be on screen based on the current path and the cached auth status.
- **Views**: The login and profile screens live in `static/js/views/`. They are responsible for attaching DOM elements, binding event listeners, and calling service functions.
- **Data + services**: `static/js/services/auth.js` talks to the Go backend over `fetch` and pushes results into a simple in-memory store defined in `static/js/state.js`. `static/js/utils/profile.js` reshapes raw API payloads into the exact structures expected by the UI components.
- **Visual components**: `static/js/charts/` contains DOM/SVG builders that turn the processed data into charts. `static/js/ui/` houses small helpers for shared DOM elements and CSS class management.

The sections below dive into each module in detail.

## Entry Point: `static/js/main.js`
- **Purpose**: Provide a single startup routine that wires together all other modules.
- **Key functions**:
  - `start()`: orchestrates startup. It registers the login/profile renderers with the router, determines which path should be rendered first (current location vs `/login` fallback), and calls `bootstrapAuthState()` to ask the server whether the user has an active session.
  - `run()`: calls `start()` and logs any unexpected startup failures. It is invoked once the DOM is ready (either immediately or via `DOMContentLoaded`).
- **Design considerations**:
  - Waiting for `bootstrapAuthState()` ensures the router knows whether the user is authenticated before rendering, so deep links to `/profile` can redirect to `/login` without a flash of the wrong screen.
  - Optional `renderOptions` pass status messages (for example, “session expired”) to the first screen render after startup.

## Shared Constants: `static/js/constants.js`
- **Purpose**: Centralize the canonical route strings (`/`, `/login`, `/profile`).
- **Why it matters**: Using a shared constant makes it impossible for two modules to drift apart (for example, one using `/login/` with a trailing slash), which would break route comparisons.

## Global State Container: `static/js/state.js`
- **Purpose**: Hold the minimal authentication context inside the browser.
- **Structure**: A simple object containing `status` (`"unknown"`, `"anonymous"`, `"authenticated"`) and `profileData` (last known profile payload).
- **Key functions**:
  - `getAuthState()`: returns the singleton state object so other modules can read it.
  - `setAuthState(status, profileData)`: mutates the state to reflect new login state and optionally stores fresh profile data.
  - `clearProfileData()`: helper to remove cached profile data without changing the status.
- **Design note**: There is no Redux/Vuex-like machinery; the plain object works because the app has only two screens and no complicated cross-component synchronization.

## Router: `static/js/router.js`
- **Purpose**: Provide client-side navigation so the app can switch between login/profile without full page reloads.
- **Key exports**:
  - `configureRouter({ loginRenderer, profileRenderer })`: stores the two render functions supplied by `main.js`.
  - `initRouter({ initialPath, renderOptions })`: registers a `popstate` listener (handles back/forward buttons), then calls `handleRoute` for the initial path. The `replace` flag prevents duplicate history entries for the bootstrap render.
  - `navigate(path, options)`: used by views to change screens programmatically (for example, after logging in).
- **Internal helpers**:
  - `normalizePath(path)`: converts arbitrary strings into clean paths (`/profile`, `/login`) and falls back to the home route.
  - `applyHistory(path, replace)`: wraps `history.pushState` / `history.replaceState` to keep the browser address bar in sync with the rendered view.
  - `handleRoute(path, options)`: core decision maker. It:
    1. Normalizes the requested path.
    2. Looks at `getAuthState().status` and rewrites the target route when needed (anonymous users are forced to login, authenticated users are steered toward the profile).
    3. Updates the history stack if the route actually changes.
    4. Avoids redundant renders unless the `renderOptions` object changed (useful for surfacing status messages on the same screen).
    5. Calls either `renderProfile()` or `renderLogin()` with the optional `renderOptions`.
- **Design intent**: Encapsulate all route-guard logic so views only focus on DOM work. The router is deliberately un-opinionated about the DOM—it just delegates rendering to whichever function matches the selected route.

## Auth Services: `static/js/services/auth.js`
- **Purpose**: Handle every network interaction related to authentication and profile data.
- **Core helper**:
  - `requestProfileData()`: POSTs to `/api/profile`, parses JSON responses when available, and standardizes success/error replies so callers never deal with raw `Response` objects.
- **Public API**:
  - `bootstrapAuthState()`: Used at startup. After calling `requestProfileData()`, it sets the global auth state to `"authenticated"` (with fetched data) or `"anonymous"`. It returns a lightweight status object so the caller can decide which screen to render.
  - `fetchProfileData()`: Re-fetches the profile for the profile view. It updates the global state based on whether the response was successful or returned 401/403 (session expired).
  - `loginUser(identifier, password)`: Sends credentials to `/api/auth/login`. On success, the global state flips to `"authenticated"`. On failure it surfaces a human-friendly error message.
  - `logoutUser()`: POSTs to `/api/auth/logout`. Success returns a message that can be shown to the user; failures are normalized just like in the other helpers.
- **Design considerations**:
  - `credentials: "include"` ensures cookies/session tokens are sent with requests.
  - All functions catch network errors and return `{ ok: false, message }` objects instead of throwing, so UI code can show feedback without `try/catch` everywhere.

## UI Utilities: `static/js/ui/dom.js`
- **Purpose**: Manage shared DOM elements used by multiple views.
- **Key responsibilities**:
  - Cache and return the app root (`<main id="app">`).
  - Lazily create (and later reuse) the global status message paragraph for inline notifications. This element can be appended to different containers depending on the current layout.
  - Provide `updateStatus(message, type)` to set text and apply `"status-error"` / `"status-success"` CSS classes.
  - Offer a small helper `removeElementById(id)` to cleanly dispose of old DOM nodes when switching views.

## Layout Utilities: `static/js/ui/layout.js`
- **Purpose**: Toggle top-level CSS classes on `<body>` so stylesheets can adapt the layout to the current screen (`mode-login` vs `mode-profile`).
- **How it works**: Removes any known mode classes first, then applies `mode-${mode}` if it is one of the allowed values. This keeps class lists tidy and prevents conflicting layout rules.

## Login Screen: `static/js/views/login.js`
- **Responsibilities**:
  - Teardown any lingering profile-specific elements (`profile-data`, logout button) before rendering the login panel. This ensures a clean slate whenever an anonymous user returns to the login screen.
  - Build the login form programmatically (`buildLoginForm()`), wiring it for client-side submission handling so the UI stays on the same page.
  - Use `ensureStatusMessage()` to mount the shared status element inside the login panel and then feed it messages through `updateStatus`.
- **Submission flow**:
  1. The `handleLoginSubmit()` listener prevents the default form POST.
  2. It validates that both identifier and password fields are filled. Missing data results in an inline error message.
  3. While waiting for `loginUser()`, the submit button is disabled and its label switches to “Signing in…” for feedback.
  4. If the login service answers with `{ ok: false }`, the error message is displayed and the user can try again.
  5. On success, `navigate(ROUTES.PROFILE, { replace: true })` lets the router take over (which will render the profile screen and update the URL to `/profile`).
- **Design intent**: Keep the view self-contained—no hidden dependencies on DOM elements outside the panel. All mutations go through small helpers so the login screen can be re-rendered repeatedly without leaking event listeners or duplicate nodes.

## Profile Screen: `static/js/views/profile.js`
- **Responsibilities**:
  - Switch the layout to profile mode, clear login-only DOM nodes, and ensure the floating logout button is mounted on the document body.
  - Read cached profile data from `getAuthState()`. If data already exists (for example, from startup bootstrap), render it immediately; otherwise, show a “Loading” skeleton while fresh data is fetched.
  - Delegate to `loadProfileData()` to asynchronously pull the latest profile, handle session expiry (redirect back to login with an error message), and update the global store.
- **Rendering pipeline**:
  1. `renderProfileSections()` receives the combined API payload and uses the utility functions in `utils/profile.js` to derive human-friendly slices of data.
  2. It constructs the DOM tree in logical chunks—summary header, charts section—clearing the container each time to avoid stale fragments.
  3. The logout button is wired to `handleLogout()`, which calls `logoutUser()`, handles failure messages, and otherwise redirects back to the login view with a success status.
- **Design intent**: Separate concerns between data shaping (utilities) and DOM building (view). This allows the same raw API payload to power multiple components while keeping DOM creation readable.
- **Note about status messaging**: The profile view intentionally leaves the global status element detached unless a status needs to be surfaced (a deliberate product choice). When `updateStatus` is called—for example, after logout failure—it will operate on the shared status node created by `ensureStatusMessage()`.

## Profile Data Utilities: `static/js/utils/profile.js`
- **Purpose**: Convert verbose GraphQL-style responses into flat arrays/objects the UI can use directly.
- **Key helpers**:
  - `deriveUserProfile(userResponse)`: Extracts the first user record, builds a friendly object with name, login, contact info, and join date.
  - `extractTotalXP(xpResponse)`: Pulls the total XP sum and ensures the result is a number.
  - `formatXPAmount(amount)`: Formats numbers into B/KB/MB strings for display and re-used by charts.
  - `deriveCollaboratorCounts(collaboratorsResponse)`: Tallies how many projects each collaborator appears in and returns a sorted array for the bar chart.
  - `deriveProgressSeries(progressResponse)`: Groups transaction records by day and outputs a date/value series suited for cumulative plotting.
  - `deriveProjectList(projectsResponse)`: Creates a sorted list of unique project names (last segment of the path, prettified).
  - `formatISODate(isoString)`: Converts ISO timestamps into localized human-readable dates.
- **Design intent**: Keep views focused on presentation. If the backend response format changes, only these utilities need updating; the rest of the UI can remain untouched.

## Chart Builders: `static/js/charts/`
- **`collaborators.js`**:
  - Accepts an array of `{ login, count }` objects.
  - Builds a `<section>` containing a title, subtitle, and an `<svg>` bar chart sized to fit the data.
  - Calculates bar widths proportionally to the maximum collaborator count and adds textual labels (handles both inside/outside placement based on available space).
  - Adds an axis line to visually ground the bars.
- **`progress.js`**:
  - Accepts an array of `{ date: Date, value: number }` points representing XP earned per day.
  - Converts the raw series into a cumulative (running total) sequence so the chart shows overall progress.
  - Builds an `<svg>` line chart with axes, tick marks, grid lines, and dots for each data point.
  - Uses helper scales to map dates and values to pixel positions and leverages `formatXPAmount` for y-axis labels.
- **Design intent**: Keep chart creation deterministic and DOM-based. Because these functions only depend on parameters, they can be tested in isolation or replaced with a different visualization library later.

## Bringing It All Together
1. The browser loads `main.js`, which registers view renderers and asks the server for the current profile.
2. The router decides which screen to show. If you’re not signed in, it forces the login view; if you are, it takes you straight to the profile (even when you land on `/` or `/login`).
3. The login view builds a form, handles submission, and on success tells the router to switch to `/profile`.
4. The profile view ensures layout/UI scaffolding is in place, renders cached data instantly when available, and then refreshes the data from the server. While waiting, it keeps the user informed via status messages and a loading skeleton.
5. Utility modules massage API responses into neatly structured objects, and chart modules turn those objects into visuals without worrying about navigation or state management.

The result is a small but cohesive single-page application: routing lives in one place, state is centralized, views are predictable, and the code is organized so each module has a clear, single responsibility.
