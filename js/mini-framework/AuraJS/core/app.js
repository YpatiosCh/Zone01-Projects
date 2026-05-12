// ============================================================
// AuraJS App — The Glue
// ============================================================
// Wires together the virtual DOM, store, router, and event bus
// into a single cohesive API. This is the main entry point for
// building applications with AuraJS.
// ============================================================

import { mount, diff, h } from "./vdom.js";
import { createStore } from "./state.js";
import { createRouter } from "./router.js";
import { createEventBus } from "./events.js";

/**
 * Creates an AuraJS application.
 *
 * This is the main entry point that wires all framework subsystems together:
 * - Creates a **store** (centralized state management)
 * - Creates a **router** (hash-based URL routing)
 * - Creates an **event bus** (pub/sub cross-component communication)
 * - Subscribes to state changes and route changes to trigger re-renders
 * - Initializes **plugins** with access to all subsystems
 * - Manages the **render loop** with `requestAnimationFrame` batching
 *
 * The initialization sequence is:
 * 1. Validate config and resolve the root DOM element
 * 2. Create store, router, and event bus
 * 3. Subscribe to store → re-render on state changes
 * 4. Register all routes
 * 5. Initialize plugins (they can set up listeners, dispatch actions, etc.)
 * 6. Start the router (fires the initial render)
 *
 * @param {Object}           config            — Application configuration.
 * @param {string|Element}   config.root       — CSS selector (e.g. `"#app"`) or DOM element to mount into.
 * @param {Function}         config.reducer    — Pure reducer function: `(state, action) => newState`.
 * @param {Function}         config.view       — View function: `(state, dispatch, ctx) => vNode`.
 *   Called on every render. Returns a virtual DOM tree that describes the current UI.
 * @param {Object}           [config.state={}] — Initial application state (deep-cloned internally).
 * @param {Array<{path: string, handler: Function}>} [config.routes=[]] — Route definitions.
 *   Each route has a `path` (e.g. `"/"`, `"/game/:id"`) and a `handler` view function.
 * @param {Array<Function>}  [config.plugins=[]] — Plugin functions. Each receives
 *   `{ store, router, events, render }` and may return a cleanup function.
 * @returns {{
 *   store: Object,
 *   router: Object,
 *   events: Object,
 *   render: Function,
 *   destroy: Function
 * }} The app instance with access to all subsystems and lifecycle methods.
 *
 * @example
 * import { createApp, h } from "./AuraJS/index.js";
 *
 * const app = createApp({
 *   root: "#app",
 *   state: { count: 0 },
 *   reducer(state, action) {
 *     if (action.type === "INC") return { ...state, count: state.count + 1 };
 *     return state;
 *   },
 *   routes: [
 *     { path: "/", handler: HomePage },
 *     { path: "/about", handler: AboutPage },
 *   ],
 *   view(state, dispatch, ctx) {
 *     return ctx.handler(state, dispatch, ctx);
 *   },
 * });
 */
export function createApp(config) {
  const {
    root: rootSelector,
    reducer,
    state: initialState = {},
    view,
    routes = [],
    plugins = [],
  } = config;

  if (!view) {
    throw new Error('AuraJS: "view" function is required.');
  }
  if (!reducer) {
    throw new Error('AuraJS: "reducer" function is required.');
  }

  // Resolve root element
  const rootEl =
    typeof rootSelector === "string"
      ? document.querySelector(rootSelector)
      : rootSelector;

  if (!rootEl) {
    throw new Error(`AuraJS: Root element "${rootSelector}" not found.`);
  }

  // Core systems
  const store = createStore(reducer, initialState);
  const router = createRouter();
  const events = createEventBus();

  // Track previous vDOM for diffing
  let prevTree = null;
  let isRendering = false;
  let renderQueued = false;
  let rafId = null;

  /**
   * Re-renders the entire application.
   *
   * Schedules a render via `requestAnimationFrame` to batch rapid state
   * changes into a single DOM update. If a render is already scheduled,
   * sets a flag to re-render again after the current one completes (so
   * no state change is ever missed).
   *
   * On each render:
   * 1. Reads the current state from the store
   * 2. Resolves the current route (handler + params)
   * 3. Builds the `ctx` object (path, navigate, events, handler, params)
   * 4. Calls the `view()` function to get a new vNode tree
   * 5. First render → `mount()` | Subsequent → `diff()` (minimal patches)
   * 6. Stores the new tree as `prevTree` for the next diff
   */
  function render() {
    if (isRendering) {
      renderQueued = true;
      return;
    }

    isRendering = true;

    rafId = requestAnimationFrame(() => {
      try {
        const state = store.getState();
        const currentRoute = router.getCurrentPath();
        const match = router.resolve();
        const newTree = view(state, store.dispatch, {
          path: currentRoute,
          navigate: router.navigate,
          events,
          handler: match ? match.handler : defaultNotFound,
          params: match ? match.params : {},
        });

        if (prevTree === null) {
          mount(newTree, rootEl);
        } else {
          diff(rootEl, prevTree, newTree);
        }

        prevTree = newTree;
      } catch (err) {
        console.error("AuraJS render error:", err);
      } finally {
        isRendering = false;

        if (renderQueued) {
          renderQueued = false;
          render();
        }
      }
    });
  }

  // Re-render on state changes
  store.subscribe(() => render());

  // Register routes
  for (const route of routes) {
    router.addRoute(route.path, route.handler);
  }

  // Initialize plugins before first render
  const pluginCleanups = [];
  for (const plugin of plugins) {
    const cleanup = plugin({ store, router, events, render });
    if (typeof cleanup === "function") {
      pluginCleanups.push(cleanup);
    }
  }

  // Re-render on route changes (also triggers initial render)
  router.start(() => render());

  /**
   * Tears down the entire application.
   *
   * Performs a full ordered cleanup:
   * 1. Cancels any pending `requestAnimationFrame`
   * 2. Calls all plugin cleanup functions
   * 3. Removes the router's `hashchange` listener
   * 4. Clears all event bus handlers
   * 5. Resets the internal vNode tree reference
   */
  function destroy() {
    if (rafId != null) {
      cancelAnimationFrame(rafId);
      rafId = null;
    }

    for (const cleanup of pluginCleanups) {
      cleanup();
    }

    router.destroy();
    events.clear();
    prevTree = null;
  }

  // Public API
  return {
    store,
    router,
    events,
    render,
    destroy,
  };
}

/**
 * Built-in default 404 page rendered when no route matches and no
 * custom `setNotFound()` handler has been registered.
 *
 * @returns {Object} A vNode tree displaying "404 / Page not found".
 */
function defaultNotFound() {
  return h("div", { style: "display:flex;flex-direction:column;align-items:center;justify-content:center;height:100vh" },
    h("h1", {}, "404"),
    h("p", {}, "Page not found"),
  );
}

// Re-export everything for convenience
export { h, renderNode, mount, diff } from "./vdom.js";
export { createStore, combineReducers } from "./state.js";
export { createRouter } from "./router.js";
export { createEventBus, delegate } from "./events.js";
