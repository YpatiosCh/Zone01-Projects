# AuraJS

A lightweight JavaScript framework built from scratch with **zero dependencies**. AuraJS provides everything needed to build interactive web applications — from a simple TodoMVC to a real-time multiplayer game — all rendered through the DOM.

Rather than relying on external libraries, every piece of this framework was written by hand to understand and control every layer of the stack: how elements get created, how state flows through the app, how routes are matched, and how components communicate. The result is a small, transparent system where nothing is hidden behind abstractions you can't inspect.

---

## Table of Contents

1. [Quick Start](#quick-start)
2. [Project Structure](#project-structure)
3. [Virtual DOM](#virtual-dom)
   - [Why a Virtual DOM?](#why-a-virtual-dom)
   - [Creating Elements with `h()`](#creating-elements-with-h)
   - [How `h()` Works Internally](#how-h-works-internally)
   - [Nesting Elements](#nesting-elements)
   - [Dynamic Children (Lists)](#dynamic-children-lists)
   - [Conditional Rendering](#conditional-rendering)
   - [Rendering to the Real DOM](#rendering-to-the-real-dom)
   - [Mounting a Tree](#mounting-a-tree)
4. [Attributes & Events](#attributes--events)
   - [HTML Attributes](#html-attributes)
   - [Class Names](#class-names)
   - [Style Objects](#style-objects)
   - [Boolean Attributes](#boolean-attributes)
   - [Event Handlers](#event-handlers)
   - [Why Stable Event Wrappers?](#why-stable-event-wrappers)
5. [Diffing & Patching](#diffing--patching)
   - [Why Diffing?](#why-diffing)
   - [The Four Cases](#the-four-cases)
   - [Attribute Diffing](#attribute-diffing)
   - [Children Diffing (Index-Based)](#children-diffing-index-based)
   - [Keyed Reconciliation](#keyed-reconciliation)
   - [When to Use Keys](#when-to-use-keys)
6. [State Management](#state-management)
   - [Why a Centralized Store?](#why-a-centralized-store)
   - [The Unidirectional Data Flow](#the-unidirectional-data-flow)
   - [Creating a Store](#creating-a-store)
   - [Dispatching Actions](#dispatching-actions)
   - [Subscribing to Changes](#subscribing-to-changes)
   - [Batched Dispatches](#batched-dispatches)
   - [Combining Reducers](#combining-reducers)
   - [Why Deep Clone the Initial State?](#why-deep-clone-the-initial-state)
7. [Routing](#routing)
   - [Why Hash-Based Routing?](#why-hash-based-routing)
   - [How the Router Works](#how-the-router-works)
   - [Setting Up Routes](#setting-up-routes)
   - [Writing a Route Handler](#writing-a-route-handler)
   - [The `ctx` Object](#the-ctx-object)
   - [Navigating Between Routes](#navigating-between-routes)
   - [Dynamic Route Parameters](#dynamic-route-parameters)
   - [404 Handling](#404-handling)
   - [Standalone Router Usage](#standalone-router-usage)
8. [Event Bus](#event-bus)
   - [Why an Event Bus?](#why-an-event-bus)
   - [Creating and Using the Event Bus](#creating-and-using-the-event-bus)
   - [How Emit Works (Snapshot Safety)](#how-emit-works-snapshot-safety)
   - [Event Delegation](#event-delegation)
9. [The App — `createApp()` Connects Everything](#the-app--createapp-connects-everything)
   - [Why `createApp()`?](#why-createapp)
   - [What Happens When You Call `createApp()`](#what-happens-when-you-call-createapp)
   - [The Render Cycle](#the-render-cycle)
   - [Config Options](#config-options)
   - [The Returned App Object](#the-returned-app-object)
10. [Plugin System](#plugin-system)
    - [Why Plugins?](#why-plugins)
    - [How Plugins Work](#how-plugins-work)
    - [Writing a Plugin](#writing-a-plugin)
    - [Plugin Lifecycle](#plugin-lifecycle)
    - [Real-World Example: Mock Server Plugin](#real-world-example-mock-server-plugin)
11. [Lifecycle & Cleanup](#lifecycle--cleanup)
12. [Full Walkthrough: How Everything Connects](#full-walkthrough-how-everything-connects)
13. [Running the TodoMVC Example](#running-the-todomvc-example)
14. [API Reference](#api-reference)

---

## Quick Start

### Minimal example

```html
<div id="app"></div>
<script type="module">
  import { createApp, h } from "./AuraJS/index.js";

  createApp({
    root: "#app",
    state: { count: 0 },
    reducer(state, action) {
      switch (action.type) {
        case "INC": return { ...state, count: state.count + 1 };
        default:    return state;
      }
    },
    view(state, dispatch) {
      return h("div", {},
        h("p", {}, `Count: ${state.count}`),
        h("button", { onClick: () => dispatch({ type: "INC" }) }, "+1")
      );
    },
  });
</script>
```

What happens here:

1. `createApp` creates a store, router, and event bus behind the scenes
2. The `view` function is called with the current state — it returns a virtual DOM tree
3. The framework renders that tree into real DOM inside `#app`
4. When the button is clicked, `dispatch({ type: "INC" })` sends an action to the reducer
5. The reducer produces a new state with `count + 1`
6. The store notifies subscribers, which triggers a re-render
7. The framework diffs the old virtual tree against the new one and patches only what changed — in this case, the text inside `<p>`

The framework handles rendering, diffing, and re-rendering automatically whenever state changes. You never touch the DOM directly.

---

## Project Structure

```
mini-framework/
├── AuraJS/
│   ├── index.js              ← single import entry point
│   └── core/
│       ├── vdom.js           ← virtual DOM: h(), renderNode(), mount(), diff()
│       ├── state.js          ← store: createStore(), combineReducers()
│       ├── router.js         ← hash-based routing: createRouter()
│       ├── events.js         ← event bus + delegation: createEventBus(), delegate()
│       └── app.js            ← createApp() glue + plugin system
└── todo-app/
    ├── index.html
    ├── app.js                ← app entry point using createApp()
    ├── state/
    │   ├── state.js          ← initial state definition
    │   └── reducer.js        ← all action handlers
    ├── views/
    │   └── todoList.js       ← main view (route handler)
    ├── components/
    │   ├── todoItem.js       ← single todo item component
    │   └── filterLink.js     ← filter navigation link component
    └── styles/
        └── todomvc.css
```

Everything is imported from a single entry point: `AuraJS/index.js`. This re-exports all public APIs from the core modules so you only ever need one import line:

```js
import { createApp, h, createStore, createRouter, createEventBus } from "./AuraJS/index.js";
```

---

## Virtual DOM

### Why a Virtual DOM?

Directly manipulating the real DOM is expensive. Every call to `document.createElement`, `element.setAttribute`, or `parent.removeChild` triggers the browser to recalculate layout, styles, and repaint the screen. In a complex app where many things change at once, this becomes a performance bottleneck.

The virtual DOM solves this by introducing a lightweight intermediary. Instead of modifying the real DOM directly, you describe **what the UI should look like** as a plain JavaScript object tree. The framework then compares the new tree to the previous one (a process called **diffing**) and computes the **minimal set of real DOM operations** needed to bring the screen up to date. This means:

- If only one attribute changed on one element, only that attribute is updated
- If a list was reordered, existing DOM nodes are moved rather than destroyed and recreated
- If nothing changed, nothing happens

You never call `document.createElement` yourself. You describe your UI declaratively and let the framework handle the imperative DOM work.

### Creating Elements with `h()`

The `h()` function (short for "hyperscript") is how you create virtual DOM nodes. It is the equivalent of writing HTML, but in JavaScript:

```js
import { h } from "./AuraJS/index.js";

// This creates a virtual <div>Hello</div>
const vNode = h("div", {}, "Hello");
```

The function signature is:

```js
h(tag, attrs, ...children)
```

| Argument      | Type   | Description                                   |
| ------------- | ------ | --------------------------------------------- |
| `tag`         | string | HTML tag name (`"div"`, `"input"`, `"span"`)  |
| `attrs`       | object | Attributes, properties, and event handlers     |
| `...children` | any    | Child vNodes, strings, numbers, or arrays      |

The `attrs` parameter is always the second argument and should be an object (or `{}` if you have no attributes). Everything after `attrs` is treated as children.

```js
// A paragraph with text
h("p", {}, "Hello, world!")

// An input with attributes
h("input", { type: "text", placeholder: "Type here..." })

// A div with a class and two children
h("div", { class: "card" },
  h("h2", {}, "Title"),
  h("p", {}, "Description")
)
```

### How `h()` Works Internally

When you call `h("div", { class: "card", key: "item-1" }, "Hello")`, the function does the following:

1. **Extracts the `key`** from `attrs` if it exists. Keys are used by the diffing algorithm for list reconciliation — they are not rendered to the DOM. The `key` is stored on the vNode object separately and removed from `attrs`.

2. **Flattens children** using `.flat(Infinity)`. This means you can pass arrays of children, nested arrays, or spread arrays — they all get flattened into a single list. This is what makes `...items.map(...)` work seamlessly.

3. **Filters out falsy values**. Children that are `null`, `undefined`, `false`, or `true` are silently removed. This is what enables conditional rendering with `&&` and ternary operators — `false && h("div", {}, "hidden")` simply disappears.

4. **Coerces primitives to strings**. Numbers and strings become string children. Objects (other vNodes) pass through unchanged.

The result is a plain object like this:

```js
{
  tag: "div",
  key: "item-1",
  attrs: { class: "card" },    // key has been removed
  children: ["Hello"]           // flattened, filtered, coerced
}
```

This object is the virtual DOM node. It is cheap to create (no DOM allocation) and cheap to compare.

### Nesting Elements

You build complex UIs by nesting `h()` calls. Each call returns a vNode, which can be passed as a child to another `h()` call:

```js
// A navigation bar with links
h("nav", { class: "navbar" },
  h("a", { href: "#/" }, "Home"),
  h("a", { href: "#/about" }, "About"),
  h("a", { href: "#/contact" }, "Contact")
)
```

```js
// A form with labeled inputs
h("form", { onSubmit: handleSubmit },
  h("div", { class: "form-group" },
    h("label", { for: "email" }, "Email"),
    h("input", { id: "email", type: "email", placeholder: "you@example.com" })
  ),
  h("div", { class: "form-group" },
    h("label", { for: "password" }, "Password"),
    h("input", { id: "password", type: "password" })
  ),
  h("button", { type: "submit" }, "Log In")
)
```

There is no limit to nesting depth. The framework recursively renders and diffs the entire tree.

### Dynamic Children (Lists)

To render a list of items, use `.map()` to transform an array of data into an array of vNodes, and spread it into the parent's children:

```js
const fruits = ["Apple", "Banana", "Cherry"];

h("ul", {},
  ...fruits.map(fruit => h("li", {}, fruit))
)
```

This works because `h()` flattens all children with `.flat(Infinity)`. You could also pass the array directly without spreading:

```js
h("ul", {},
  fruits.map(fruit => h("li", {}, fruit))
)
```

Both produce the same result — three `<li>` children inside a `<ul>`.

For lists that can reorder, add keys (see [Keyed Reconciliation](#keyed-reconciliation)):

```js
const todos = [
  { id: 1, text: "Buy milk" },
  { id: 2, text: "Walk the dog" },
];

h("ul", {},
  ...todos.map(todo =>
    h("li", { key: todo.id }, todo.text)
  )
)
```

### Conditional Rendering

Since falsy values are filtered from children, you can use standard JavaScript for conditional rendering:

```js
// Ternary: show one thing or another
h("div", {},
  isLoggedIn
    ? h("p", {}, "Welcome back!")
    : h("p", {}, "Please log in")
)

// Logical AND: show or hide
h("div", {},
  h("h1", {}, "Dashboard"),
  hasErrors && h("div", { class: "alert" }, "Something went wrong!"),
  isLoading && h("span", {}, "Loading...")
)
```

When `hasErrors` is `false`, the expression `false && h(...)` evaluates to `false`, which `h()` filters out. The element simply doesn't appear in the tree — no special `v-if` directive needed.

### Rendering to the Real DOM

The `renderNode()` function takes a vNode and creates the corresponding real DOM node. It is called internally by the framework, but you can use it directly if needed:

```js
import { renderNode } from "./AuraJS/index.js";

const vNode = h("div", { class: "card" },
  h("p", {}, "Hello")
);

const domNode = renderNode(vNode);
// domNode is now a real <div class="card"><p>Hello</p></div>
```

How it works internally:

1. If the vNode is a string or number, it creates a `TextNode` via `document.createTextNode()`
2. Otherwise, it creates an element via `document.createElement(vNode.tag)`
3. It iterates over `vNode.attrs` and calls `setAttr()` for each attribute (handling events, styles, booleans, etc.)
4. It recursively calls `renderNode()` for each child and appends it

### Mounting a Tree

The `mount()` function renders a vNode tree and places it into a DOM container:

```js
import { mount, h } from "./AuraJS/index.js";

const tree = h("div", {}, h("p", {}, "Mounted!"));
const container = document.querySelector("#app");

mount(tree, container);
```

`mount()` clears the container's contents (`innerHTML = ""`), renders the vNode tree into real DOM, and appends it. It is used for the **first render** of the application. After the first render, all subsequent updates use `diff()` instead.

---

## Attributes & Events

### HTML Attributes

Pass any standard HTML attribute as a key in the `attrs` object:

```js
h("input", {
  type: "text",
  placeholder: "Enter your name",
  class: "input-field",
  id: "name-input",
  disabled: true,
})
```

The framework uses `el.setAttribute(key, value)` for regular attributes. Special cases (events, styles, booleans, value) are handled differently as described below.

### Class Names

Both `class` and `className` are supported. `className` is automatically mapped to the `class` attribute:

```js
h("div", { class: "container active" })
h("div", { className: "container active" })  // equivalent
```

### Style Objects

Pass a plain JavaScript object for inline styles. Property names use camelCase (just like `element.style` in JavaScript):

```js
h("div", {
  style: {
    backgroundColor: "red",
    padding: "10px",
    fontSize: "16px",
    border: "1px solid black",
  }
})
```

When the framework applies styles, it sets each property individually on `el.style` rather than replacing the entire style string. During diffing, only **changed style properties** are updated and **removed properties** are cleared — unchanged ones are left alone. This avoids unnecessary style recalculations.

### Boolean Attributes

Boolean attributes (`checked`, `disabled`, `autofocus`, `selected`, etc.) are handled specially:

```js
h("input", { type: "checkbox", checked: true })
h("button", { disabled: false })
```

- When `true`: the attribute is set on the element (`el.setAttribute(key, "")`) AND the DOM property is set (`el[key] = true`)
- When `false`: the attribute is removed (`el.removeAttribute(key)`) AND the DOM property is set (`el[key] = false`)

This dual approach is necessary because some HTML elements (like checkboxes) rely on the DOM property for their visual state, while CSS selectors and accessibility tools rely on the attribute.

### Event Handlers

Attach event handlers by prefixing any DOM event name with `on` and capitalizing the first letter of the event name:

```js
// Click handler
h("button", {
  onClick: (e) => console.log("clicked!", e),
}, "Click me")

// Input handler
h("input", {
  onInput: (e) => dispatch({ type: "SET_TEXT", text: e.target.value }),
})

// Keyboard handler
h("input", {
  onKeydown: (e) => {
    if (e.key === "Enter") dispatch({ type: "SUBMIT" });
  },
})

// Mouse handlers
h("div", {
  onMouseenter: () => dispatch({ type: "HOVER_START" }),
  onMouseleave: () => dispatch({ type: "HOVER_END" }),
})

// Double-click
h("label", {
  onDblclick: () => dispatch({ type: "EDIT_TODO", id: todo.id }),
}, todo.text)

// Form submission
h("form", {
  onSubmit: (e) => {
    e.preventDefault();
    dispatch({ type: "SUBMIT_FORM" });
  },
}, /* children */)

// Focus and blur
h("input", {
  onFocus: () => console.log("focused"),
  onBlur: () => console.log("blurred"),
})
```

The naming convention is: `on` + `EventName`. The framework strips the `on` prefix and lowercases the rest to get the actual DOM event name (`onClick` → `click`, `onKeydown` → `keydown`, `onDblclick` → `dblclick`).

### Why Stable Event Wrappers?

In a naive implementation, every re-render would detach the old event listener and attach a new one — because the handler function reference changes on every render (arrow functions in your view create a new closure each time). This is wasteful.

AuraJS avoids this by using **stable wrappers**. The first time an event handler is set on an element, the framework:

1. Stores the actual handler in `el.__listeners[eventName]`
2. Creates a permanent wrapper function: `(e) => el.__listeners[eventName](e)`
3. Attaches the wrapper to the DOM with `addEventListener` — this only happens once

On subsequent renders, the framework simply updates `el.__listeners[eventName]` to point to the new handler. The wrapper is still the same function reference, so no `removeEventListener`/`addEventListener` cycle is needed. The DOM listener stays in place and always calls whatever the current handler is.

This means zero listener churn regardless of how often you re-render.

---

## Diffing & Patching

### Why Diffing?

On the first render, the framework creates the entire DOM tree from scratch using `mount()`. But after state changes, rebuilding the whole DOM would be wasteful — most of the UI is unchanged.

Instead, the framework:

1. Calls your `view()` function to get a **new** virtual tree
2. Compares the new tree against the **previous** virtual tree
3. Applies **only the differences** to the real DOM

This is the `diff()` function. It walks both trees in parallel and decides for each node what to do.

### The Four Cases

`diff(parent, oldTree, newTree, index)` handles four cases:

**Case 1 — New node added:** The old tree has nothing at this position, but the new tree has a node. The framework calls `renderNode()` and appends it.

```
old: null     →  new: <p>Hello</p>    →  appendChild(renderNode(newNode))
```

**Case 2 — Old node removed:** The old tree has a node, but the new tree has nothing at this position. The framework removes the corresponding DOM node.

```
old: <p>Hello</p>    →  new: null    →  removeChild(existingNode)
```

**Case 3 — Node replaced:** The node has fundamentally changed — different tag, different type (text vs element), different key, or different text content. The framework creates a brand new DOM node and replaces the old one.

```
old: <p>Hello</p>    →  new: <div>World</div>    →  replaceChild(renderNode(newNode), oldNode)
```

The `changed()` function determines this by checking:
- Type mismatch (string vs object, i.e. text node vs element)
- Text content changed (for text nodes)
- Different tag names (`p` vs `div`)
- Different keys

**Case 4 — Same node, diff in-place:** The node has the same tag and key. The framework diffs the attributes and then recursively diffs the children. No DOM node is created or destroyed — the existing one is updated in place.

```
old: <div class="a">    →  new: <div class="b">    →  setAttribute("class", "b")
```

### Attribute Diffing

When two elements have the same tag, their attributes are compared using two `for...in` loops (this is more efficient than allocating a `Set` of all keys):

1. **Iterate new attrs:** For each attribute in the new tree, if it differs from the old value, update it. Style objects get special treatment — they are diffed property-by-property (see below).

2. **Iterate old attrs:** For each attribute in the old tree, if it no longer exists in the new tree, remove it.

**Style diffing:** Rather than replacing the entire style string, the framework compares individual CSS properties:

```js
// Old style: { color: "red", fontSize: "16px" }
// New style: { color: "blue", padding: "10px" }

// Result:
// - color: updated from "red" to "blue"
// - fontSize: removed (set to "")
// - padding: added "10px"
```

This is important because replacing the entire style attribute would cause the browser to recalculate all style properties, even unchanged ones.

### Children Diffing (Index-Based)

When children don't have keys, the framework uses **index-based** comparison:

1. **Diff overlapping children:** Walk both lists up to the shorter length, diffing each pair at the same index
2. **Remove extras:** If the old list was longer, remove trailing DOM nodes (iterated in reverse to avoid index shifting)
3. **Append new:** If the new list is longer, create and append the new nodes

This is simple and efficient for static lists that don't reorder. However, for dynamic lists (where items move, get inserted in the middle, or get removed from arbitrary positions), this can lead to unnecessary DOM work — which is why keyed reconciliation exists.

### Keyed Reconciliation

When any child in a list has a `key` attribute, the framework switches to a **keyed diffing** algorithm. This is critical for performance in dynamic lists.

```js
// Add a key to each item so the differ can track them
h("ul", {},
  ...items.map(item =>
    h("li", { key: item.id }, item.name)
  )
)
```

**How the keyed algorithm works step by step:**

1. **Build a map from old children:** `Map<key, { vNode, index }>`. This lets us look up any old child by its key in O(1) time.

2. **Walk new children:** For each new child:
   - If its key exists in the old map → **reuse** the existing DOM node. Diff its attributes and children in-place (no destroy/recreate). Add the DOM node to the new ordering.
   - If its key is not in the old map → **create** a new DOM node via `renderNode()`.

3. **Remove unmatched old nodes:** Any old key that wasn't matched by a new child gets its DOM node removed.

4. **Reorder DOM nodes:** Walk the new ordering and use `insertBefore()` to move DOM nodes into the correct position. If a node is already in the right spot, nothing happens.

5. **Trim trailing extras:** Remove any leftover DOM nodes beyond the new list's length.

**Why this matters — an example:**

Suppose you have a list `[A, B, C, D]` and reorder it to `[D, A, B, C]`:

- Without keys: the differ sees index 0 changed from A to D, index 1 changed from B to A, etc. It **replaces the content of every node** — 4 updates.
- With keys: the differ sees D moved to the front. It **moves one DOM node** — 1 operation. A, B, C stay in place with no changes.

This difference is even more significant when list items have complex subtrees, focus state, animations, or form inputs that would lose their state if destroyed and recreated.

### When to Use Keys

Use keys whenever your list is **dynamic** — items can be added, removed, or reordered:

```js
// Todo items — can be added, deleted, reordered
...todos.map(todo => h("li", { key: todo.id }, todo.text))

// Players in a game — can join, leave, move around
...players.map(player =>
  h("div", { key: player.id, style: { left: `${player.x}px`, top: `${player.y}px` } },
    player.name
  )
)
```

Keys must be **unique among siblings** and **stable** (the same item should always have the same key). Use IDs from your data, not array indices.

You **don't need keys** for static lists that never change order, like a navigation menu.

---

## State Management

### Why a Centralized Store?

In any interactive application, state (the data that determines what the UI shows) needs to live somewhere. The simplest approach — scattering state across DOM elements or global variables — quickly becomes unmanageable: you lose track of what changed, when, and why.

AuraJS uses a **centralized store** inspired by the Flux/Redux pattern. This means:

- **All application state** lives in one place — a single JavaScript object
- **State can only change** through dispatched actions (plain objects describing what happened)
- **A reducer function** takes the current state and an action, and returns the new state
- **Subscribers** are notified whenever state changes, triggering re-renders

This architecture makes the app **predictable** — given the same state, the view always renders the same output. Given the same sequence of actions, you always end up with the same state. There are no hidden side effects or mystery mutations.

### The Unidirectional Data Flow

```
User clicks button
       ↓
dispatch({ type: "ADD_TODO", text: "Buy milk" })
       ↓
reducer(currentState, action) → newState
       ↓
store notifies all subscribers
       ↓
view(newState) → new virtual DOM tree
       ↓
diff(oldTree, newTree) → minimal DOM patches
       ↓
User sees updated UI
```

Data flows in **one direction**: actions go up (from user interaction to the store), state flows down (from the store to the view). This is intentional — it eliminates circular dependencies and makes debugging straightforward. If the UI is wrong, you check the state. If the state is wrong, you check the reducer. If the reducer is wrong, you check the action.

### Creating a Store

```js
import { createStore } from "./AuraJS/index.js";

// The reducer is a pure function: (state, action) => newState
function reducer(state, action) {
  switch (action.type) {
    case "ADD_ITEM":
      return { ...state, items: [...state.items, action.item] };
    case "REMOVE_ITEM":
      return { ...state, items: state.items.filter(i => i.id !== action.id) };
    default:
      return state;  // IMPORTANT: always return state for unknown actions
  }
}

// Create the store with a reducer and initial state
const store = createStore(reducer, { items: [] });
```

The `createStore` function returns an object with four methods:

| Method        | Description                                                    |
| ------------- | -------------------------------------------------------------- |
| `getState()`  | Returns the current state object                               |
| `dispatch(action)` | Runs the reducer with the action, notifies subscribers   |
| `subscribe(fn)`    | Registers a listener; returns an unsubscribe function    |
| `batch(fn)`        | Groups multiple dispatches into a single notification    |

### Dispatching Actions

An action is a plain object with a `type` property and optional payload data:

```js
// Simple action
store.dispatch({ type: "INCREMENT" });

// Action with payload
store.dispatch({ type: "ADD_TODO", text: "Buy milk" });

// Action with multiple payload fields
store.dispatch({ type: "MOVE_PLAYER", playerId: 1, x: 100, y: 200 });
```

When `dispatch` is called:

1. The reducer is called with the current state and the action
2. The reducer returns a new state object (it must **never** mutate the old state)
3. If the new state is a different reference than the old state (`!==`), subscribers are notified
4. If the new state is the same reference (e.g. the `default` case returned `state`), nothing happens — no wasted re-renders

**This reference equality check is why reducers must return new objects** when something changes. Returning `{ ...state, count: state.count + 1 }` creates a new object, which triggers a re-render. Mutating `state.count++` and returning `state` would not trigger anything because it is the same reference.

### Subscribing to Changes

```js
const unsubscribe = store.subscribe((newState, oldState) => {
  console.log("State changed!");
  console.log("Old:", oldState);
  console.log("New:", newState);
});

// Stop listening when you no longer need it:
unsubscribe();
```

Subscribers receive both the **new state** and the **previous state**. This makes it easy to detect what changed — you can compare specific properties between old and new.

The subscriber list is **snapshotted** before notification — if a subscriber unsubscribes during notification, it won't affect the current notification cycle. This prevents subtle bugs where removing a listener shifts the array and causes other listeners to be skipped.

### Batched Dispatches

Sometimes you need to dispatch multiple actions at once — for example, setting a player's X position, Y position, and health simultaneously. Without batching, each dispatch triggers a separate re-render:

```js
// BAD: 3 dispatches = 3 re-renders (2 are wasted)
store.dispatch({ type: "SET_X", x: 100 });     // render 1
store.dispatch({ type: "SET_Y", y: 200 });     // render 2
store.dispatch({ type: "SET_HP", health: 75 }); // render 3
```

With `batch()`, all dispatches are grouped and subscribers fire only once at the end:

```js
// GOOD: 3 dispatches = 1 re-render
store.batch(() => {
  store.dispatch({ type: "SET_X", x: 100 });
  store.dispatch({ type: "SET_Y", y: 200 });
  store.dispatch({ type: "SET_HP", health: 75 });
});
// Subscribers fire ONCE here, with the final combined state
```

**How batching works internally:**

- A `batchDepth` counter tracks nesting level
- When `batch()` is called, `batchDepth` increments
- While inside a batch, dispatches still run the reducer (state updates immediately) but subscriber notifications are **deferred**
- The first dispatch that actually changes state captures `batchPrev` — the pre-batch state
- When the outermost `batch()` completes (`batchDepth` returns to 0), subscribers fire once with the final state and the `batchPrev` as the old state
- A `try/finally` block ensures `batchDepth` always decrements, even if an error is thrown

**Nested batches** are supported. Only the outermost batch triggers notification:

```js
store.batch(() => {
  store.dispatch({ type: "A" });
  store.batch(() => {
    store.dispatch({ type: "B" });
    store.dispatch({ type: "C" });
  });
  // Inner batch ended, but outer batch is still open — no notification yet
  store.dispatch({ type: "D" });
});
// All 4 actions applied. Single notification here.
```

### Combining Reducers

As your app grows, a single reducer becomes unwieldy. `combineReducers()` lets you split state into independent slices, each managed by its own reducer:

```js
import { combineReducers, createStore } from "./AuraJS/index.js";

function todosReducer(state = [], action) {
  switch (action.type) {
    case "ADD_TODO": return [...state, action.todo];
    case "REMOVE_TODO": return state.filter(t => t.id !== action.id);
    default: return state;
  }
}

function userReducer(state = null, action) {
  switch (action.type) {
    case "LOGIN": return action.user;
    case "LOGOUT": return null;
    default: return state;
  }
}

const rootReducer = combineReducers({
  todos: todosReducer,
  user: userReducer,
});

const store = createStore(rootReducer, {
  todos: [],
  user: null,
});
```

`combineReducers` works by:

1. Calling each sub-reducer with its slice of state and the action
2. Checking if any slice changed (by reference `!==`)
3. If anything changed, returning a new combined object
4. If nothing changed, returning the **same** state reference (which means no re-render)

Each sub-reducer only sees its own slice. `todosReducer` receives `state.todos` (the array), not the entire state object.

### Why Deep Clone the Initial State?

When creating a store, the initial state is **deep-cloned**. This prevents a subtle bug: if you define your initial state as a module-level object and create multiple store instances (e.g. in tests), they would share the same object reference. Mutating one would corrupt the other. Deep cloning ensures each store gets its own independent copy.

The deep clone handles: plain objects, arrays, `Date`, `RegExp`, `Map`, `Set`, and nested combinations of all of these.

---

## Routing

### Why Hash-Based Routing?

Single-page applications (SPAs) need a way to show different views (pages) without full page reloads. There are two common approaches:

1. **History API routing** (`/lobby`, `/game/42`) — cleaner URLs but requires server-side configuration to redirect all paths to `index.html`
2. **Hash-based routing** (`/#/lobby`, `/#/game/42`) — works out of the box with any static file server because the hash fragment is never sent to the server

AuraJS uses hash-based routing because it requires **zero server configuration**. You can open the HTML file directly, serve it with any static server, or deploy it anywhere — routing just works. The browser fires a `hashchange` event whenever the hash changes, which the router listens to.

### How the Router Works

```
User clicks link or calls navigate("/lobby")
       ↓
window.location.hash changes to "#/lobby"
       ↓
Browser fires "hashchange" event
       ↓
Router's hashchange listener fires the callback
       ↓
Callback triggers a re-render
       ↓
During render, app calls router.resolve()
       ↓
Router matches "/lobby" against registered route patterns
       ↓
Returns { handler: LobbyView, params: {} }
       ↓
The handler is passed to your view via ctx.handler
       ↓
ctx.handler(state, dispatch, ctx) → vNode tree for the lobby page
```

The router itself does **not** render anything. It is purely a matching engine. It takes a URL path, matches it against registered patterns, and returns the corresponding handler function. The rendering is done by the view function, which calls the handler to get the virtual DOM tree.

Internally, route patterns like `"/game/:id"` are compiled into regular expressions (`/^\/game\/([^\/]+)$/`) when registered. This compilation happens once at registration time, so matching at render time is a simple regex test.

### Setting Up Routes

Routes are defined as an array of `{ path, handler }` objects in your `createApp` config:

```js
import { createApp, h } from "./AuraJS/index.js";
import Home from "./views/home.js";
import About from "./views/about.js";
import UserProfile from "./views/userProfile.js";

const app = createApp({
  root: "#app",
  reducer: myReducer,
  state: initialState,

  routes: [
    { path: "/", handler: Home },
    { path: "/about", handler: About },
    { path: "/user/:id", handler: UserProfile },
  ],

  view(state, dispatch, ctx) {
    return h("div", {},
      // Navigation bar (always visible)
      h("nav", {},
        h("a", { href: "#/" }, "Home"),
        h("a", { href: "#/about" }, "About"),
      ),
      // Render the matched route handler
      ctx.handler(state, dispatch, ctx)
    );
  },
});
```

Routes are matched **in order** — the first route whose pattern matches the current URL wins. Place more specific routes before less specific ones.

### Writing a Route Handler

A route handler is a plain function that receives `(state, dispatch, ctx)` and returns a virtual DOM tree. It is exactly like the main `view` function, but responsible for one specific page:

```js
// views/home.js
import { h } from "../AuraJS/index.js";

export default function Home(state, dispatch, ctx) {
  return h("div", { class: "home-page" },
    h("h1", {}, "Welcome Home"),
    h("p", {}, `You have ${state.items.length} items.`),
    h("button", {
      onClick: () => ctx.navigate("/about"),
    }, "Go to About page")
  );
}
```

### The `ctx` Object

Every view function (including route handlers) receives a `ctx` (context) object as its third argument. This object is the bridge between the routing system, the event bus, and your view:

| Property   | Type     | Description                                                   |
| ---------- | -------- | ------------------------------------------------------------- |
| `handler`  | Function | The matched route handler function (the current page's view)  |
| `params`   | Object   | Dynamic route parameters extracted from the URL               |
| `path`     | string   | The current URL path (e.g. `"/game/42"`)                      |
| `navigate` | Function | Programmatically change the URL hash to navigate to a route   |
| `events`   | Object   | The event bus instance for cross-component communication      |

The `ctx` object is constructed fresh on every render by `createApp`. It reflects the current state of the router and provides access to shared framework services without needing global imports.

### Navigating Between Routes

There are two ways to navigate:

**1. Declarative navigation with anchor tags:**

```js
h("a", { href: "#/" }, "Home")
h("a", { href: "#/about" }, "About")
h("a", { href: "#/game/42" }, "Game 42")
```

Since these are standard `<a>` tags with `href`, clicking them changes `window.location.hash`, which fires the `hashchange` event and triggers a re-render.

**2. Programmatic navigation with `ctx.navigate()`:**

```js
// Inside a view function
h("button", {
  onClick: () => ctx.navigate("/lobby"),
}, "Go to Lobby")

// Navigate after an async operation
async function handleSubmit() {
  await saveData();
  ctx.navigate("/success");
}
```

`navigate(path)` simply sets `window.location.hash = path`, which triggers the same hashchange flow. It is a convenience method so you don't have to write `window.location.hash = "#/lobby"` yourself.

### Dynamic Route Parameters

Routes support `:param` segments that match any value and extract it as a named parameter:

```js
routes: [
  { path: "/user/:id", handler: UserProfile },
  { path: "/posts/:postId/comments/:commentId", handler: Comment },
]
```

When the URL is `/#/user/42`, the router matches `/user/:id` and extracts `{ id: "42" }`. When the URL is `/#/posts/5/comments/12`, it extracts `{ postId: "5", commentId: "12" }`.

Access parameters via `ctx.params`:

```js
export default function UserProfile(state, dispatch, ctx) {
  const userId = ctx.params.id;  // "42"

  return h("div", {},
    h("h1", {}, `User Profile: ${userId}`),
    h("button", {
      onClick: () => ctx.navigate("/"),
    }, "Back to Home")
  );
}
```

**How parameter extraction works internally:** The route path `"/user/:id"` is compiled to a regex `/^\/user\/([^\/]+)$/`. The `:id` becomes a capture group. When the URL matches, the captured values are paired with the parameter names extracted from the original path pattern.

### 404 Handling

If no registered route matches the current URL, the framework provides a **built-in default 404 page** — a centered "404 / Page not found" message. This works out of the box with zero configuration.

To customize the 404 page, use `setNotFound()` on the router:

```js
const app = createApp({ /* ... */ });

app.router.setNotFound((state, dispatch, ctx) => {
  return h("div", { class: "not-found" },
    h("h1", {}, "404"),
    h("p", {}, "The page you're looking for doesn't exist."),
    h("button", {
      onClick: () => ctx.navigate("/"),
    }, "Go Home")
  );
});
```

The custom 404 handler is a regular view function — it receives `(state, dispatch, ctx)` and returns a vNode tree, just like any route handler.

### Standalone Router Usage

The router can be used independently of `createApp()`. This is useful if you want to integrate routing into an existing app or use it for a different purpose:

```js
import { createRouter } from "./AuraJS/index.js";

const router = createRouter();

// Register routes
router.addRoute("/", homeHandler);
router.addRoute("/about", aboutHandler);
router.addRoute("/user/:id", userHandler);

// Start listening for hash changes
router.start(() => {
  const match = router.resolve();
  if (match) {
    // match.handler is the view function
    // match.params is the extracted parameters
    renderSomehow(match.handler, match.params);
  }
});

// Navigate programmatically
router.navigate("/about");

// Get current path
router.getCurrentPath(); // "/about"

// Cleanup when done
router.destroy();
```

---

## Event Bus

### Why an Event Bus?

In a component-based architecture, components sometimes need to communicate without being in a direct parent-child relationship. For example:

- A WebSocket connection receives a message and needs to notify multiple views
- A plugin needs to signal that something happened
- Two unrelated parts of the UI need to coordinate

Passing callbacks through props works for parent-child communication, but for **decoupled, cross-cutting concerns**, a pub/sub event bus is more appropriate. Components publish events without knowing who's listening, and listeners subscribe without knowing who's emitting.

### Creating and Using the Event Bus

```js
import { createEventBus } from "./AuraJS/index.js";

const events = createEventBus();
```

**Subscribe to an event:**

```js
const unsubscribe = events.on("user:login", (userData) => {
  console.log(`${userData.name} logged in`);
});
```

`on()` returns an unsubscribe function. Call it to stop listening.

**Emit an event:**

```js
events.emit("user:login", { name: "Alice", id: 1 });
```

All handlers registered for `"user:login"` will be called with the data object.

**One-time listener:**

```js
events.once("app:ready", () => {
  console.log("App has booted — this only fires once");
});
```

The handler is automatically unsubscribed after it fires once.

**Remove a specific handler:**

```js
function myHandler(data) { /* ... */ }
events.on("some:event", myHandler);
events.off("some:event", myHandler);
```

**Clear handlers:**

```js
events.clear("user:login");  // Clear handlers for one specific event
events.clear();               // Clear ALL handlers for ALL events
```

### How Emit Works (Snapshot Safety)

When `emit()` is called, the handler array is **snapshotted** (copied via `.slice()`) before iteration. This means:

- A handler can safely unsubscribe itself during emission without affecting other handlers
- A handler can safely subscribe new handlers during emission without causing infinite loops

This is a deliberate design choice. Without snapshotting, removing a handler during iteration would shift array indices and cause handlers to be skipped or called twice.

### Event Delegation

For situations where you need traditional DOM event delegation (one listener on a parent that handles events for dynamically created children), AuraJS provides the `delegate()` utility:

```js
import { delegate } from "./AuraJS/index.js";

const cleanup = delegate(
  document.querySelector(".todo-list"),  // parent element
  "click",                                // event type
  ".delete-btn",                          // CSS selector for targets
  (event, matchedElement) => {            // handler
    const id = matchedElement.dataset.id;
    store.dispatch({ type: "DELETE_TODO", id });
  }
);

// Later, remove the delegation:
cleanup();
```

**How it works:** A single event listener is attached to the `root` element. When an event fires, it uses `event.target.closest(selector)` to find the nearest ancestor matching the CSS selector. If found and contained within the root, the handler is called with both the original event and the matched element.

This is useful for cases where the virtual DOM is not managing the elements (e.g. integrating with third-party libraries) or when you want explicit DOM-level delegation outside the framework's render cycle.

---

## The App — `createApp()` Connects Everything

### Why `createApp()`?

The virtual DOM, store, router, and event bus are independent modules that can each be used on their own. But building an app requires wiring them together:

- The store needs to trigger re-renders when state changes
- The router needs to trigger re-renders when the URL changes
- The view function needs access to state, dispatch, route info, and the event bus
- Plugins need access to all of the above
- Cleanup needs to be coordinated across all systems

`createApp()` is the **glue function** that connects all these pieces. It creates instances of each subsystem, wires the data flow, manages the render loop, and provides a clean public API.

### What Happens When You Call `createApp()`

Here is the exact sequence of operations:

```js
const app = createApp({
  root: "#app",
  reducer: myReducer,
  state: { count: 0 },
  view: myView,
  routes: [{ path: "/", handler: Home }],
  plugins: [myPlugin],
});
```

1. **Validates config** — throws if `view` or `reducer` is missing
2. **Resolves root element** — if `root` is a CSS selector string, queries the DOM; if it's already an element, uses it directly. Throws if not found.
3. **Creates the store** — `createStore(reducer, state)` with deep-cloned initial state
4. **Creates the router** — `createRouter()` with no routes yet
5. **Creates the event bus** — `createEventBus()` with no handlers
6. **Subscribes to store changes** — `store.subscribe(() => render())` so every state change triggers a re-render
7. **Registers all routes** — iterates the `routes` array and calls `router.addRoute(path, handler)` for each
8. **Initializes plugins** — calls each plugin function with `{ store, router, events, render }`. If a plugin returns a function, it is saved as a cleanup function for later teardown
9. **Starts the router** — `router.start(() => render())` begins listening for `hashchange` events. This also fires the callback **immediately**, which triggers the **first render**

After this, the app is live. State changes re-render. Route changes re-render. Plugins are active. Everything is connected.

### The Render Cycle

The render function is the beating heart of the app. Here is what happens on every render:

1. **Read current state** from the store via `store.getState()`
2. **Resolve the current route** via `router.resolve()` — returns `{ handler, params }` or falls back to the default 404 handler
3. **Build the context object** (`ctx`) with: `path`, `navigate`, `events`, `handler`, `params`
4. **Call the view function** with `view(state, dispatch, ctx)` — returns a new vNode tree
5. **First render?** Call `mount(newTree, rootEl)` to create the entire DOM from scratch
6. **Subsequent render?** Call `diff(rootEl, prevTree, newTree)` to patch only what changed
7. **Store the new tree** as `prevTree` for the next diff

**Render batching with `requestAnimationFrame`:**

Renders are scheduled via `requestAnimationFrame` (rAF) to ensure:

- Multiple rapid state changes (even outside of `batch()`) don't cause multiple synchronous DOM updates
- Rendering happens at the browser's optimal time (aligned with the display refresh rate)
- The UI never blocks the main thread with excessive DOM work

The render function uses two flags to manage this:

- `isRendering` — prevents re-entrant renders. If a render is already scheduled, subsequent calls set `renderQueued = true` instead of scheduling another rAF.
- `renderQueued` — after a render completes, if this flag is true, another render is immediately scheduled. This ensures no state change is ever missed.

### Config Options

| Option    | Type             | Required | Description                                                |
| --------- | ---------------- | -------- | ---------------------------------------------------------- |
| `root`    | string / Element | yes      | CSS selector or DOM element to mount the app into          |
| `reducer` | Function         | yes      | `(state, action) => newState` — pure function              |
| `state`   | Object           | no       | Initial state (defaults to `{}`)                           |
| `view`    | Function         | yes      | `(state, dispatch, ctx) => vNode` — returns the UI         |
| `routes`  | Array            | no       | `[{ path, handler }]` route definitions                    |
| `plugins` | Array            | no       | Plugin functions to initialize                             |

### The Returned App Object

`createApp()` returns an object that gives you access to the internals:

```js
const app = createApp({ /* ... */ });

app.store     // Store instance: { getState, dispatch, subscribe, batch }
app.router    // Router instance: { addRoute, setNotFound, navigate, resolve, ... }
app.events    // Event bus instance: { on, off, emit, once, clear }
app.render()  // Manually trigger a re-render
app.destroy() // Full teardown (see Lifecycle & Cleanup)
```

You will primarily use `app.destroy()` for cleanup. The other properties are useful for debugging, testing, or advanced use cases like registering additional routes or event handlers after initialization.

---

## Plugin System

### Why Plugins?

The core framework handles state, rendering, routing, and events. But real applications often need additional functionality: WebSocket connections, logging, analytics, persistence, mock servers for development, etc.

Rather than bloating the core with every possible feature, AuraJS provides a **plugin system** that lets you extend the framework cleanly. Plugins have full access to the store, router, event bus, and render function — they can do anything, from dispatching actions in response to external events, to intercepting route changes, to bridging the app with a WebSocket server.

### How Plugins Work

A plugin is a function that receives an object containing the framework's core systems:

```js
function myPlugin({ store, router, events, render }) {
  // Plugin has access to everything:
  // - store.getState(), store.dispatch(), store.subscribe(), store.batch()
  // - router.navigate(), router.addRoute(), etc.
  // - events.on(), events.emit(), etc.
  // - render() to trigger re-renders manually

  // Optionally return a cleanup function
  return () => {
    // Teardown logic here
  };
}
```

| Property | Type     | Description                                              |
| -------- | -------- | -------------------------------------------------------- |
| `store`  | Object   | The store instance: `getState`, `dispatch`, `subscribe`, `batch` |
| `router` | Object   | The router instance: `navigate`, `addRoute`, `resolve`, etc.     |
| `events` | Object   | The event bus instance: `on`, `off`, `emit`, `once`, `clear`     |
| `render` | Function | Manually trigger an app re-render                               |

### Writing a Plugin

**Simple logging plugin:**

```js
function loggerPlugin({ store }) {
  const unsub = store.subscribe((newState, oldState) => {
    console.log("State changed:", { oldState, newState });
  });

  return () => unsub();
}
```

**Persistence plugin (save state to localStorage):**

```js
function persistPlugin({ store }) {
  // Restore state on load
  const saved = localStorage.getItem("app-state");
  if (saved) {
    store.dispatch({ type: "RESTORE_STATE", state: JSON.parse(saved) });
  }

  // Save state on every change
  const unsub = store.subscribe((state) => {
    localStorage.setItem("app-state", JSON.stringify(state));
  });

  return () => unsub();
}
```

**Plugin factory pattern (configurable plugins):**

When a plugin needs configuration, wrap it in a factory function:

```js
function analyticsPlugin(config) {
  return function plugin({ store, events }) {
    events.on("page:view", (page) => {
      sendToAnalytics(config.apiKey, page);
    });

    const unsub = store.subscribe((state) => {
      trackStateChange(config.apiKey, state);
    });

    return () => unsub();
  };
}

// Usage:
createApp({
  plugins: [analyticsPlugin({ apiKey: "abc123" })],
  // ...
});
```

### Plugin Lifecycle

1. **Initialization:** Plugins are called **after** the store, router, and event bus are created, but **before** the first render. This means plugins can set up listeners, modify the router, or dispatch initial actions before the UI appears.

2. **Active:** While the app is running, plugins participate through the subscriptions and event handlers they set up during initialization. They react to state changes, route changes, and custom events.

3. **Cleanup:** When `app.destroy()` is called, the framework calls every plugin's cleanup function (if one was returned). This is where plugins should unsubscribe from events, close connections, and release resources.

### Real-World Example: Mock Server Plugin

The included `mockServerPlugin` demonstrates how plugins bridge external systems with the framework using the event bus as a communication channel:

```js
export function mockServerPlugin() {
  return function plugin({ store, router, events }) {
    // Listen for messages the client "sends to the server"
    const unsub = events.on("ws:send", (msg) => {
      console.log("SERVER received:", msg);

      // Process the message and dispatch state changes back
      if (msg.type === "MOVE_PLAYER_1") {
        store.dispatch({ type: "MOVE_PLAYER_1", move: 2 });
        router.navigate("/lobby");
      }

      if (msg.type === "MOVE_PLAYER_2") {
        store.dispatch({ type: "MOVE_PLAYER_2", move: 2 });
      }
    });

    // Cleanup: remove the event listener
    return () => unsub();
  };
}
```

**How this works:**

1. Views emit `events.emit("ws:send", { type: "MOVE_PLAYER_1" })` — this represents a client sending a message
2. The plugin receives the message via `events.on("ws:send", ...)` — it acts as the "server"
3. The plugin processes the message and dispatches actions back to the store — the "server response"
4. The store update triggers a re-render — the UI reflects the server's response

This pattern decouples the views from the communication layer. The views don't know (or care) whether messages go to a real WebSocket server or a mock — they just emit events. Swapping the mock plugin for a real WebSocket plugin requires no view changes.

---

## Lifecycle & Cleanup

### `app.destroy()`

Calling `destroy()` on the app instance performs a full, ordered teardown:

1. **Cancels pending renders** — cancels any queued `requestAnimationFrame`
2. **Runs plugin cleanups** — calls every cleanup function returned by plugins
3. **Destroys the router** — removes the `hashchange` event listener from `window`
4. **Clears the event bus** — removes all event handlers
5. **Resets the vDOM** — sets the internal previous tree to `null`

```js
const app = createApp({ /* ... */ });

// Later, when shutting down:
app.destroy();
```

### Individual Subsystem Cleanup

Each subsystem also supports standalone cleanup, useful when you're using them independently of `createApp`:

```js
// Router: stop listening for hash changes
router.destroy();

// Event bus: clear specific or all handlers
events.clear("user:login");  // clear one event
events.clear();               // clear all events

// Store: unsubscribe a specific listener
const unsub = store.subscribe(myListener);
unsub();
```

---

## Full Walkthrough: How Everything Connects

To solidify understanding, here is a complete trace of what happens from app creation to a user clicking a button:

### Step 1: App Creation

```js
const app = createApp({
  root: "#app",
  state: { count: 0 },
  reducer(state, action) {
    if (action.type === "INC") return { ...state, count: state.count + 1 };
    return state;
  },
  routes: [{ path: "/", handler: CounterPage }],
  view(state, dispatch, ctx) {
    return ctx.handler(state, dispatch, ctx);
  },
});
```

- Store created with `{ count: 0 }` (deep cloned)
- Router created, `/` route registered
- Event bus created
- Store subscribed to trigger `render()` on changes
- Router started → fires callback immediately → `render()` called

### Step 2: First Render

- `render()` schedules a `requestAnimationFrame`
- On the next frame:
  - `store.getState()` returns `{ count: 0 }`
  - `router.resolve()` matches `/` → `{ handler: CounterPage, params: {} }`
  - `ctx` is built: `{ path: "/", handler: CounterPage, params: {}, navigate, events }`
  - `view(state, dispatch, ctx)` calls `ctx.handler(state, dispatch, ctx)` → calls `CounterPage`
  - `CounterPage` returns a vNode tree: `h("div", {}, h("p", {}, "Count: 0"), h("button", { onClick: ... }, "+1"))`
  - Since `prevTree` is `null`, `mount(newTree, rootEl)` is called
  - `renderNode()` recursively creates real DOM: `<div><p>Count: 0</p><button>+1</button></div>`
  - Container is cleared and the DOM tree is appended
  - `prevTree` is set to the vNode tree

### Step 3: User Clicks "+1"

- The button's `onClick` handler fires: `dispatch({ type: "INC" })`
- Inside `dispatch`:
  - Reducer called: `reducer({ count: 0 }, { type: "INC" })` → `{ count: 1 }`
  - New state `{ count: 1 }` is a different reference than old `{ count: 0 }` → subscribers are notified
  - The store's subscriber calls `render()`

### Step 4: Re-Render

- `render()` schedules a `requestAnimationFrame`
- On the next frame:
  - `store.getState()` returns `{ count: 1 }`
  - `view()` is called → `CounterPage` returns a new vNode tree with `"Count: 1"` instead of `"Count: 0"`
  - `diff(rootEl, prevTree, newTree)` is called
  - The differ walks the tree:
    - `<div>` → same tag → diff children
    - `<p>` → same tag → diff children
    - `"Count: 0"` vs `"Count: 1"` → text changed → `replaceChild(newTextNode, oldTextNode)`
    - `<button>` → same tag → diff attrs → `onClick` handler reference updated (via stable wrapper, no DOM work) → diff children → `"+1"` unchanged
  - `prevTree` updated to the new tree
- The user sees "Count: 1" on screen. Only the text node inside `<p>` was touched.

---

## Running the TodoMVC Example

The `todo-app/` directory contains a full TodoMVC implementation built with AuraJS. It demonstrates CRUD operations, routing, filtering, inline editing, keyboard handling, and more.

### With a local server

```bash
# Using npx (Node.js required)
npx serve .

# Then open in your browser:
# http://localhost:3000/todo-app/
```

```bash
# Or with Python
python3 -m http.server 8000

# Then open:
# http://localhost:8000/todo-app/
```

### Features demonstrated

| Feature              | How it works                                                       |
| -------------------- | ------------------------------------------------------------------ |
| **Add todos**        | Type in the input, press Enter → dispatches `ADD_TODO`             |
| **Toggle completion**| Click the checkbox → dispatches `TOGGLE_TODO`                      |
| **Edit todos**       | Double-click a label → dispatches `EDIT_TODO`, shows input         |
| **Save edits**       | Press Enter or blur the edit input → dispatches `SAVE_TODO`        |
| **Cancel edits**     | Press Escape → dispatches `CANCEL_EDIT`                            |
| **Delete todos**     | Click the X button → dispatches `DELETE_TODO`                      |
| **Filter**           | Click All / Active / Completed → hash-based routing via `ctx.path` |
| **Toggle all**       | Click the down-arrow → dispatches `TOGGLE_ALL`                     |
| **Clear completed**  | Click "Clear completed" button → dispatches `CLEAR_COMPLETED`      |

The TodoMVC example is a good reference for how to structure an app with AuraJS: separate files for state, reducer, views, and reusable components.

---

## API Reference

### Exports from `AuraJS/index.js`

| Export              | Module           | Description                               |
| ------------------- | ---------------- | ----------------------------------------- |
| `createApp`         | `core/app.js`    | Create and wire a full application        |
| `h`                 | `core/vdom.js`   | Create a virtual DOM node                 |
| `renderNode`        | `core/vdom.js`   | Render a vNode to a real DOM node         |
| `mount`             | `core/vdom.js`   | Mount a vNode tree into a DOM container   |
| `diff`              | `core/vdom.js`   | Diff and patch two vNode trees            |
| `createStore`       | `core/state.js`  | Create a centralized state store          |
| `combineReducers`   | `core/state.js`  | Combine multiple reducers into one        |
| `createRouter`      | `core/router.js` | Create a hash-based router                |
| `createEventBus`    | `core/events.js` | Create a pub/sub event bus                |
| `delegate`          | `core/events.js` | Attach a delegated DOM event listener     |

### `h(tag, attrs, ...children)` → vNode

Creates a virtual DOM node. Children are flattened, falsy values filtered, primitives coerced to strings. If `attrs.key` is provided, it is extracted for keyed reconciliation and not rendered to the DOM.

### `renderNode(vNode)` → DOM Node

Converts a vNode (or string) into a real DOM node. Recursively renders children. Handles text nodes, elements, attributes, events, and styles.

### `mount(vNode, container)` → DOM Element

Renders a vNode tree and appends it to the container (after clearing the container). Used for the initial render.

### `diff(parent, oldTree, newTree, index?)` → void

Compares two vNode trees and applies the minimal set of DOM operations to update the real DOM. Handles additions, removals, replacements, attribute changes, and child reconciliation (both index-based and keyed).

### `createStore(reducer, initialState?)` → Store

Creates a centralized state store.

- `getState()` → returns the current state object
- `dispatch(action)` → runs the reducer with the action, notifies subscribers if state changed (by reference)
- `subscribe(fn)` → registers a listener that receives `(newState, prevState)`; returns an unsubscribe function
- `batch(fn)` → executes `fn`, deferring all subscriber notifications until `fn` completes. Supports nesting.

### `combineReducers(reducerMap)` → Function

Takes an object mapping state keys to reducer functions. Returns a single reducer that delegates to each sub-reducer for its slice of state. Returns the same state reference if nothing changed.

### `createRouter()` → Router

Creates a hash-based router.

- `addRoute(path, handler)` → register a route pattern (supports `:param` dynamic segments) and its view handler function
- `setNotFound(handler)` → set a custom 404 handler
- `resolve()` → match the current URL hash against routes; returns `{ handler, params }` or falls back to notFoundHandler, or `null`
- `start(callback)` → listen for `hashchange` events; fires callback on every change and once immediately
- `destroy()` → remove the `hashchange` listener
- `navigate(path)` → set `window.location.hash` to trigger navigation
- `getCurrentPath()` → returns the current hash path (e.g. `"/lobby"`)

### `createEventBus()` → EventBus

Creates a pub/sub event bus.

- `on(event, handler)` → subscribe to an event; returns an unsubscribe function
- `off(event, handler)` → remove a specific handler from an event
- `emit(event, data)` → notify all subscribers for the event (handler list is snapshotted before iteration for safety)
- `once(event, handler)` → subscribe to an event; handler fires once then auto-unsubscribes
- `clear(event?)` → if event name provided, clear that event's handlers; otherwise clear all handlers for all events

### `delegate(root, eventType, selector, handler)` → Function

Attaches a single event listener to `root` that fires `handler(event, matchedElement)` when an event of type `eventType` occurs on an element matching the CSS `selector` within `root`. Uses `event.target.closest(selector)` for matching. Returns a cleanup function that removes the listener.
