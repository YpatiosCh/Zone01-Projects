// ============================================================
// Router — Hash-based Client-Side Routing
// ============================================================
// Uses the URL hash fragment (e.g. /#/lobby) to determine which
// view to render. The router itself does NOT render anything —
// it simply notifies the app when the route changes, and the
// app calls resolve() to get the matched handler.
//
// Hash-based routing was chosen because it requires zero server
// configuration — no catch-all redirects, no rewrite rules. The
// hash fragment is never sent to the server, so any static file
// server (or even opening the HTML file directly) works.
//
// Flow:
//   1. Routes are registered via addRoute(path, handler)
//   2. start(callback) listens for hashchange events
//   3. On hash change, callback fires → triggers app re-render
//   4. During render, app calls resolve() to match the current
//      URL to a route and get its handler + params
//   5. The handler (your view function) is passed to the main
//      view via ctx.handler, which calls it to get the vNode tree
// ============================================================

/**
 * Creates a hash-based client-side router.
 *
 * Routes are registered with `addRoute(path, handler)` where `path` can
 * include dynamic segments (e.g. `"/game/:id"`). The router listens for
 * `hashchange` events and notifies the app to re-render. During render,
 * `resolve()` matches the current URL hash against registered patterns
 * and returns the corresponding handler and extracted parameters.
 *
 * The router does **not** render anything itself — it is purely a matching
 * engine. Rendering is handled by the view function that calls the matched
 * handler.
 *
 * @returns {{
 *   addRoute: Function,
 *   setNotFound: Function,
 *   navigate: Function,
 *   start: Function,
 *   destroy: Function,
 *   getCurrentPath: Function,
 *   resolve: Function
 * }}
 *
 * @example
 * const router = createRouter();
 * router.addRoute("/", HomeView);
 * router.addRoute("/user/:id", UserView);
 * router.start(() => {
 *   const match = router.resolve();
 *   // match.handler === UserView, match.params === { id: "42" }
 * });
 */
export function createRouter() {
  const routes = [];
  let notFoundHandler = null;
  let hashHandler = null;

  /**
   * Registers a route pattern and its view handler.
   *
   * Supports static paths (`"/lobby"`) and dynamic segments (`"/game/:id"`).
   * Dynamic segments match any non-slash characters and are extracted as
   * named parameters accessible via `ctx.params`.
   *
   * The path is compiled to a regular expression at registration time, so
   * matching at render time is a fast regex test.
   *
   * Routes are matched **in registration order** — place more specific
   * routes before less specific ones.
   *
   * @param {string}   path    — Route pattern (e.g. `"/"`, `"/game/:id"`).
   * @param {Function} handler — View function: `(state, dispatch, ctx) => vNode`.
   *
   * @example
   * router.addRoute("/", HomeView);
   * router.addRoute("/user/:id", UserProfileView);
   * router.addRoute("/posts/:postId/comments/:commentId", CommentView);
   */
  function addRoute(path, handler) {
    routes.push({ path, regex: pathToRegex(path), handler });
  }

  /**
   * Sets a custom 404 handler that is used when no route matches.
   *
   * If not set, `createApp()` provides a built-in default 404 page.
   * The custom handler is a regular view function — it receives
   * `(state, dispatch, ctx)` and returns a vNode tree.
   *
   * @param {Function} handler — View function for 404: `(state, dispatch, ctx) => vNode`.
   *
   * @example
   * router.setNotFound((state, dispatch, ctx) => {
   *   return h("div", {},
   *     h("h1", {}, "Page not found"),
   *     h("button", { onClick: () => ctx.navigate("/") }, "Go Home")
   *   );
   * });
   */
  function setNotFound(handler) {
    notFoundHandler = handler;
  }

  /**
   * Programmatically navigates to a new route.
   *
   * Sets `window.location.hash`, which fires a `hashchange` event and
   * triggers a re-render through the router's `start()` callback.
   *
   * @param {string} path — Target route path (e.g. `"/lobby"`, `"/game/42"`).
   *
   * @example
   * router.navigate("/lobby");
   * // or from a view: ctx.navigate("/game/42");
   */
  function navigate(path) {
    window.location.hash = path;
  }

  /**
   * Returns the current route path from the URL hash.
   *
   * Strips the leading `#` from `window.location.hash`. Defaults to `"/"`
   * if the hash is empty.
   *
   * @returns {string} The current path (e.g. `"/lobby"`, `"/game/42"`, `"/"`).
   */
  function getCurrentPath() {
    return window.location.hash.slice(1) || "/";
  }

  /**
   * Matches the current URL hash against all registered routes.
   *
   * Returns the **first** matching route's handler and extracted parameters.
   * If no route matches, falls back to the `notFoundHandler` (if set via
   * `setNotFound()`), or returns `null` (in which case `createApp()` uses
   * its built-in default 404 page).
   *
   * @returns {{ handler: Function, params: Object } | null}
   *   The matched route's handler and params, or `null` if no match and
   *   no notFoundHandler is set.
   *
   * @example
   * // URL is /#/game/42, route registered as "/game/:id"
   * const match = router.resolve();
   * // match === { handler: GameView, params: { id: "42" } }
   */
  function resolve() {
    const path = getCurrentPath();

    for (const route of routes) {
      const match = path.match(route.regex);
      if (match) {
        return { handler: route.handler, params: extractParams(route.path, match) };
      }
    }

    return notFoundHandler ? { handler: notFoundHandler, params: {} } : null;
  }

  /**
   * Starts listening for URL hash changes.
   *
   * Attaches a `hashchange` event listener to `window` and fires the
   * callback on every hash change. Also fires the callback **once
   * immediately** for the initial route (so the app renders on load).
   *
   * The callback receives no arguments — the app calls `resolve()` during
   * its render cycle to get the current route match.
   *
   * @param {Function} callback — Invoked on every route change (typically triggers a re-render).
   */
  function start(callback) {
    hashHandler = () => callback();
    window.addEventListener("hashchange", hashHandler);
    hashHandler();
  }

  /**
   * Removes the `hashchange` event listener and stops the router.
   *
   * Called automatically by `app.destroy()`. Can also be called manually
   * when using the router standalone.
   */
  function destroy() {
    if (hashHandler) {
      window.removeEventListener("hashchange", hashHandler);
      hashHandler = null;
    }
  }

  return { addRoute, setNotFound, navigate, start, destroy, getCurrentPath, resolve };
}

// ---- Helpers ----

/**
 * Compiles a route path pattern into a regular expression.
 *
 * Escapes forward slashes and replaces `:param` segments with capture groups
 * that match any non-slash characters.
 *
 * @param {string} path — Route pattern (e.g. `"/game/:id"`).
 * @returns {RegExp} A regex that matches the path (e.g. `/^\/game\/([^\/]+)$/`).
 */
function pathToRegex(path) {
  const pattern = path
    .replace(/\//g, "\\/")
    .replace(/:(\w+)/g, "([^\\/]+)");
  return new RegExp(`^${pattern}$`);
}

/**
 * Extracts named parameters from a regex match result.
 *
 * Pairs the `:param` names from the original path pattern with the
 * captured values from the regex match.
 *
 * @param {string}          path  — Original route pattern (e.g. `"/game/:id"`).
 * @param {RegExpMatchArray} match — Result of `currentPath.match(routeRegex)`.
 * @returns {Object} Parameter key-value pairs (e.g. `{ id: "42" }`).
 */
function extractParams(path, match) {
  const keys = (path.match(/:(\w+)/g) || []).map((k) => k.slice(1));
  const params = {};
  keys.forEach((key, i) => {
    params[key] = match[i + 1];
  });
  return params;
}
