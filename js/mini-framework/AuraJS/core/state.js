// ============================================================
// State Management — Centralized Reactive Store
// ============================================================
// Inspired by the Flux/Redux pattern: a single store holds all
// app state. State is updated ONLY through actions dispatched
// to a reducer. Subscribers are notified on every change.
// ============================================================

/**
 * Creates a centralized state store.
 *
 * All application state lives in one place. State can only change through
 * dispatched actions — plain objects with a `type` property — which are
 * processed by a pure reducer function `(state, action) => newState`.
 *
 * The initial state is **deep-cloned** on creation to prevent accidental
 * mutation of the original object (important when reusing the same initial
 * state across tests or multiple instances).
 *
 * Change detection uses **reference equality** (`!==`). Reducers must return
 * a new object when state changes (e.g. `{ ...state, count: state.count + 1 }`)
 * — mutating and returning the same object will not trigger subscribers.
 *
 * @param {Function} reducer — Pure function: `(state, action) => newState`.
 *   Must handle unknown action types by returning the current state unchanged.
 * @param {Object} [initialState={}] — The starting state. Deep-cloned internally.
 * @returns {{ getState: Function, dispatch: Function, subscribe: Function, batch: Function }}
 *
 * @example
 * const store = createStore(
 *   (state, action) => {
 *     if (action.type === "INC") return { ...state, count: state.count + 1 };
 *     return state;
 *   },
 *   { count: 0 }
 * );
 *
 * store.subscribe((newState, oldState) => console.log(newState.count));
 * store.dispatch({ type: "INC" }); // logs: 1
 */
export function createStore(reducer, initialState = {}) {
  let state = deepClone(initialState);
  let listeners = [];
  let batchDepth = 0;
  let batchDirty = false;
  let batchPrev = null;

  /**
   * Returns the current state object.
   *
   * @returns {Object} The current application state.
   */
  function getState() {
    return state;
  }

  /**
   * Dispatches an action to the reducer to produce a new state.
   *
   * Calls `reducer(currentState, action)` and, if the returned state is a
   * different reference (`!==`), notifies all subscribers. If inside a
   * `batch()`, notification is deferred until the outermost batch completes.
   *
   * @param {{ type: string, [key: string]: * }} action — A plain object with
   *   a `type` property and optional payload fields.
   *
   * @example
   * store.dispatch({ type: "ADD_TODO", text: "Buy milk" });
   */
  function dispatch(action) {
    const prev = state;
    state = reducer(state, action);

    if (state !== prev) {
      if (batchDepth > 0) {
        // Inside a batch: defer notification
        if (!batchDirty) {
          batchPrev = prev;
        }
        batchDirty = true;
      } else {
        notify(state, prev);
      }
    }
  }

  /**
   * Notifies all subscribers with the current and previous state.
   * Snapshots the listener array before iterating to prevent issues
   * if a subscriber unsubscribes during notification.
   *
   * @param {Object} current — The new state.
   * @param {Object} prev    — The previous state.
   */
  function notify(current, prev) {
    const snapshot = listeners.slice();
    snapshot.forEach((fn) => fn(current, prev));
  }

  /**
   * Registers a listener that is called whenever state changes.
   *
   * The listener receives `(newState, prevState)`, making it easy to
   * detect what changed by comparing specific properties.
   *
   * @param {Function} fn — Callback: `(newState, prevState) => void`.
   * @returns {Function} An unsubscribe function — call it to stop listening.
   *
   * @example
   * const unsub = store.subscribe((newState, oldState) => {
   *   console.log("Count changed:", oldState.count, "→", newState.count);
   * });
   * unsub(); // stop listening
   */
  function subscribe(fn) {
    listeners.push(fn);
    return () => {
      listeners = listeners.filter((l) => l !== fn);
    };
  }

  /**
   * Groups multiple dispatches into a single subscriber notification.
   *
   * Without batching, each `dispatch()` triggers a re-render. `batch()`
   * defers all notifications until the callback completes, then fires
   * subscribers **once** with the final state.
   *
   * Supports nesting — only the outermost `batch()` triggers notification.
   * Uses `try/finally` to ensure the batch depth always decrements, even
   * if the callback throws.
   *
   * @param {Function} fn — A callback that calls `dispatch()` one or more times.
   *
   * @example
   * store.batch(() => {
   *   store.dispatch({ type: "SET_X", x: 100 });
   *   store.dispatch({ type: "SET_Y", y: 200 });
   * });
   * // Subscribers fire ONCE with the final state
   */
  function batch(fn) {
    batchDepth++;
    try {
      fn();
    } finally {
      batchDepth--;
      if (batchDepth === 0 && batchDirty) {
        batchDirty = false;
        const prev = batchPrev;
        batchPrev = null;
        notify(state, prev);
      }
    }
  }

  return { getState, dispatch, subscribe, batch };
}

/**
 * Combines multiple reducers into a single root reducer.
 *
 * Each key in `reducerMap` corresponds to a slice of the state tree.
 * Each sub-reducer only receives its own slice (not the entire state).
 *
 * Uses **reference equality** to detect changes — if no sub-reducer
 * returns a new reference, the original state object is returned
 * unchanged (which means no re-render).
 *
 * @param {Object<string, Function>} reducerMap — An object mapping state
 *   keys to their reducer functions. Each reducer: `(sliceState, action) => newSliceState`.
 * @returns {Function} A combined reducer: `(state, action) => newState`.
 *
 * @example
 * const rootReducer = combineReducers({
 *   todos: todosReducer,   // manages state.todos
 *   user: userReducer,     // manages state.user
 * });
 *
 * const store = createStore(rootReducer, { todos: [], user: null });
 */
export function combineReducers(reducerMap) {
  return (state = {}, action) => {
    const next = {};
    let changed = false;
    for (const [key, reducer] of Object.entries(reducerMap)) {
      const prev = state[key];
      next[key] = reducer(prev, action);
      if (next[key] !== prev) changed = true;
    }
    return changed ? next : state;
  };
}

// ============================================================
// Helpers
// ============================================================

/**
 * Deep-clones a value, handling plain objects, arrays, Date, RegExp, Map, and Set.
 *
 * Used internally to clone the initial state passed to `createStore()`,
 * ensuring each store instance gets an independent copy.
 *
 * @param {*} obj — The value to clone.
 * @returns {*} A deep copy of the value.
 */
function deepClone(obj) {
  if (obj === null || typeof obj !== "object") return obj;
  if (obj instanceof Date) return new Date(obj.getTime());
  if (obj instanceof RegExp) return new RegExp(obj.source, obj.flags);
  if (obj instanceof Map)
    return new Map([...obj].map(([k, v]) => [deepClone(k), deepClone(v)]));
  if (obj instanceof Set) return new Set([...obj].map(deepClone));
  if (Array.isArray(obj)) return obj.map(deepClone);
  const clone = {};
  for (const key of Object.keys(obj)) {
    clone[key] = deepClone(obj[key]);
  }
  return clone;
}
