// ============================================================
// Event Handling — Custom Event Bus & Delegation
// ============================================================
// A pub/sub event system for decoupled cross-component
// communication, plus a DOM event delegation utility.
// ============================================================

/**
 * Creates a pub/sub event bus for decoupled cross-component communication.
 *
 * The event bus lets unrelated parts of the application communicate without
 * direct references to each other. Components publish named events with
 * `emit()` and subscribe to them with `on()` — neither side needs to know
 * who the other is.
 *
 * This is especially useful for:
 * - Plugin-to-view communication (e.g. a WebSocket plugin dispatching events)
 * - Cross-cutting concerns (analytics, logging)
 * - Coordinating unrelated UI sections
 *
 * The event bus is also available inside views via `ctx.events` when using
 * `createApp()`.
 *
 * @returns {{
 *   on: Function,
 *   off: Function,
 *   emit: Function,
 *   once: Function,
 *   clear: Function
 * }}
 *
 * @example
 * const events = createEventBus();
 *
 * const unsub = events.on("user:login", (user) => {
 *   console.log(`${user.name} logged in`);
 * });
 *
 * events.emit("user:login", { name: "Alice" });
 * unsub(); // stop listening
 */
export function createEventBus() {
  let handlers = {};

  /**
   * Subscribes a handler to a named event.
   *
   * The handler will be called every time `emit(event, data)` is called
   * with the matching event name. Returns an unsubscribe function for
   * convenient cleanup.
   *
   * @param {string}   event   — Event name (e.g. `"user:login"`, `"ws:message"`).
   * @param {Function} handler — Callback that receives the emitted data.
   * @returns {Function} An unsubscribe function — call it to remove this handler.
   *
   * @example
   * const unsub = events.on("todo:added", (todo) => {
   *   console.log("New todo:", todo.text);
   * });
   * unsub(); // stop listening
   */
  function on(event, handler) {
    if (!handlers[event]) handlers[event] = [];
    handlers[event].push(handler);

    return () => off(event, handler);
  }

  /**
   * Removes a specific handler from a named event.
   *
   * If you need to unsubscribe, prefer using the function returned by `on()`
   * instead — it's more convenient and less error-prone.
   *
   * @param {string}   event   — Event name.
   * @param {Function} handler — The exact function reference to remove.
   */
  function off(event, handler) {
    if (!handlers[event]) return;
    handlers[event] = handlers[event].filter((h) => h !== handler);
  }

  /**
   * Emits a named event, calling all subscribed handlers with the given data.
   *
   * The handler array is **snapshotted** (`.slice()`) before iteration. This
   * means it is safe for a handler to unsubscribe itself or subscribe new
   * handlers during emission — neither action will affect the current
   * notification cycle.
   *
   * @param {string} event — Event name to emit.
   * @param {*}      data  — Payload passed to each handler.
   *
   * @example
   * events.emit("player:moved", { id: 1, x: 100, y: 200 });
   */
  function emit(event, data) {
    const list = handlers[event];
    if (!list) return;
    const snapshot = list.slice();
    snapshot.forEach((h) => h(data));
  }

  /**
   * Subscribes a handler that fires **only once**, then auto-unsubscribes.
   *
   * Useful for one-time setup events like `"app:ready"` or waiting for
   * a single response.
   *
   * @param {string}   event   — Event name.
   * @param {Function} handler — Callback that receives the emitted data.
   * @returns {Function} An unsubscribe function (in case you need to cancel before it fires).
   *
   * @example
   * events.once("app:ready", () => console.log("App booted!"));
   */
  function once(event, handler) {
    const wrapper = (data) => {
      handler(data);
      off(event, wrapper);
    };
    return on(event, wrapper);
  }

  /**
   * Clears event handlers.
   *
   * - With an event name: removes all handlers for that specific event.
   * - Without arguments: removes **all** handlers for **all** events.
   *
   * Called automatically by `app.destroy()` to clean up the entire bus.
   *
   * @param {string} [event] — If provided, clear only this event's handlers.
   *   If omitted, clear everything.
   *
   * @example
   * events.clear("user:login");  // clear one event
   * events.clear();               // clear all events
   */
  function clear(event) {
    if (event) {
      delete handlers[event];
    } else {
      handlers = {};
    }
  }

  return { on, off, emit, once, clear };
}

// ============================================================
// Delegate — Declarative DOM Event Delegation
// ============================================================

/**
 * Attaches a single delegated event listener to a parent element.
 *
 * Instead of attaching individual listeners to every child element (which
 * is expensive for dynamic lists), one listener is placed on the `root`
 * container. When an event fires, `event.target.closest(selector)` is used
 * to find the nearest ancestor matching the CSS selector. If found (and
 * contained within `root`), the handler is called.
 *
 * This pattern survives re-renders and works with dynamically added elements
 * without needing to re-attach listeners.
 *
 * @param {Element}  root      — Container element to attach the listener to.
 * @param {string}   eventType — DOM event type (e.g. `"click"`, `"input"`, `"keydown"`).
 * @param {string}   selector  — CSS selector to match target elements within `root`.
 * @param {Function} handler   — Callback: `(event, matchedElement) => void`.
 *   Receives the original DOM event and the element that matched the selector.
 * @returns {Function} A cleanup function — call it to remove the listener.
 *
 * @example
 * const cleanup = delegate(
 *   document.querySelector(".todo-list"),
 *   "click",
 *   ".delete-btn",
 *   (event, el) => {
 *     store.dispatch({ type: "DELETE_TODO", id: el.dataset.id });
 *   }
 * );
 *
 * cleanup(); // remove the delegation
 */
export function delegate(root, eventType, selector, handler) {
  const listener = (e) => {
    const target = e.target.closest(selector);
    if (target && root.contains(target)) {
      handler(e, target);
    }
  };

  root.addEventListener(eventType, listener);

  return () => root.removeEventListener(eventType, listener);
}
