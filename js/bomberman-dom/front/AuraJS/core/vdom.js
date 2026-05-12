// ============================================================
// Virtual DOM — The heart of AuraJS
// ============================================================
// Instead of manipulating the real DOM directly, we build a
// lightweight JS object tree (virtual DOM), diff it against the
// previous tree, and apply only the minimal set of changes.
// ============================================================

/**
 * Creates a virtual DOM node (vNode).
 *
 * Builds a lightweight JS object that describes a DOM element. These objects
 * are cheap to create and compare, and the framework uses them to compute
 * the minimal set of real DOM operations needed on each render.
 *
 * Children are automatically flattened (nested arrays are OK), and falsy
 * values (`null`, `undefined`, `false`, `true`) are filtered out — this
 * enables conditional rendering with `&&` and ternaries.
 *
 * If `attrs.key` is provided, it is extracted onto the vNode for keyed
 * reconciliation and will **not** be rendered as a DOM attribute.
 *
 * @param {string} tag        — HTML tag name (e.g. `"div"`, `"input"`, `"span"`)
 * @param {Object} [attrs={}] — Attributes, properties, and event handlers.
 *   Supports: HTML attrs (`class`, `id`, `type`), boolean attrs (`checked`,
 *   `disabled`), style objects (`{ color: "red" }`), `className` shorthand,
 *   and event handlers prefixed with `on` (`onClick`, `onInput`, `onKeydown`).
 * @param {...(Object|string|number|Array|null|boolean)} children — Child
 *   vNodes, strings, numbers, or nested arrays. Falsy values are ignored.
 * @returns {{ tag: string, key: string|null, attrs: Object, children: Array }}
 *   A plain vNode object used by `renderNode()`, `mount()`, and `diff()`.
 *
 * @example
 * // Simple element with text
 * h("p", {}, "Hello, world!")
 *
 * @example
 * // Nested elements with attributes
 * h("div", { class: "card" },
 *   h("h2", {}, "Title"),
 *   h("p", {}, "Description")
 * )
 *
 * @example
 * // Dynamic list with keys
 * h("ul", {},
 *   ...items.map(item => h("li", { key: item.id }, item.name))
 * )
 *
 * @example
 * // Conditional rendering
 * h("div", {},
 *   isLoggedIn && h("span", {}, "Welcome back!"),
 *   hasError ? h("p", { class: "error" }, errorMsg) : null
 * )
 */
export function h(tag, attrs = {}, ...children) {
  const a = attrs || {};
  const key = a.key != null ? a.key : null;

  // Remove key from attrs so it's not rendered to DOM
  if (key !== null) {
    const { key: _k, ...rest } = a;
    return {
      tag,
      key,
      attrs: rest,
      children: children
        .flat(Infinity)
        .filter((c) => c !== null && c !== undefined && c !== false && c !== true)
        .map((c) => (typeof c === "object" ? c : String(c))),
    };
  }

  return {
    tag,
    key: null,
    attrs: a,
    children: children
      .flat(Infinity)
      .filter((c) => c !== null && c !== undefined && c !== false && c !== true)
      .map((c) => (typeof c === "object" ? c : String(c))),
  };
}

/**
 * Converts a vNode (or string/number) into a real DOM node.
 *
 * Recursively walks the vNode tree: creates elements via
 * `document.createElement`, applies attributes and event handlers via
 * `setAttr()`, and appends children. Text nodes (strings/numbers) are
 * created with `document.createTextNode`.
 *
 * This is called internally by `mount()` and `diff()` — you typically
 * don't need to call it directly unless building custom rendering logic.
 *
 * @param {Object|string|number} vNode — A vNode object from `h()`, or a
 *   string/number for text nodes.
 * @returns {Node} A real DOM node (Element or Text) ready to be inserted
 *   into the document.
 *
 * @example
 * const vNode = h("div", { class: "card" }, h("p", {}, "Hello"));
 * const domNode = renderNode(vNode);
 * // domNode is a real <div class="card"><p>Hello</p></div>
 */
export function renderNode(vNode) {
  // Text nodes
  if (typeof vNode === "string" || typeof vNode === "number") {
    return document.createTextNode(String(vNode));
  }

  const el = document.createElement(vNode.tag);

  // Apply attributes
  for (const key in vNode.attrs) {
    setAttr(el, key, vNode.attrs[key]);
  }

  // Recursively render children
  for (const child of vNode.children || []) {
    el.appendChild(renderNode(child));
  }

  return el;
}

/**
 * Renders a vNode tree and mounts it into a DOM container.
 *
 * Clears the container's existing content, converts the vNode tree into
 * real DOM via `renderNode()`, and appends it. Used for the **first render**
 * of the application. Subsequent updates should use `diff()` instead.
 *
 * @param {Object} vNode      — Root vNode tree created with `h()`.
 * @param {Element} container — DOM element to mount into (e.g. `document.querySelector("#app")`).
 * @returns {Element} The rendered root DOM element.
 *
 * @example
 * const tree = h("div", {}, h("p", {}, "Mounted!"));
 * mount(tree, document.querySelector("#app"));
 */
export function mount(vNode, container) {
  const dom = renderNode(vNode);
  container.innerHTML = "";
  container.appendChild(dom);
  return dom;
}

// ============================================================
// Diffing & Patching
// ============================================================

/**
 * Diffs two vNode trees and patches the real DOM in-place.
 *
 * Compares `oldTree` and `newTree` at the given child `index` within
 * `parent`, and applies the minimal DOM mutations needed:
 *
 * 1. **Add** — `oldTree` is null, `newTree` exists → `appendChild`
 * 2. **Remove** — `oldTree` exists, `newTree` is null → `removeChild`
 * 3. **Replace** — tag, type, key, or text changed → `replaceChild`
 * 4. **Update** — same tag/key → diff attributes, then recursively diff children
 *
 * Children are diffed using index-based comparison by default. If any child
 * has a `key` attribute, the algorithm switches to **keyed reconciliation**,
 * which reuses and reorders existing DOM nodes for optimal performance.
 *
 * @param {Element} parent  — The parent DOM node containing the children to diff.
 * @param {Object|string|null} oldTree — The previous vNode (from the last render).
 * @param {Object|string|null} newTree — The new vNode (from the current render).
 * @param {number} [index=0] — The child index within `parent` to diff at.
 * @returns {void}
 *
 * @example
 * // Typically called by the framework's render loop:
 * diff(rootEl, previousVTree, newVTree);
 */
export function diff(parent, oldTree, newTree, index = 0) {
  const existingNode = parent.childNodes[index];

  // Null safety: both null or newTree null with no existing node
  if (newTree == null && oldTree == null) return;
  if (newTree == null && !existingNode) return;

  // 1. New node added
  if (oldTree == null) {
    parent.appendChild(renderNode(newTree));
    return;
  }

  // 2. Old node removed
  if (newTree == null) {
    if (existingNode) parent.removeChild(existingNode);
    return;
  }

  // 3. Node replaced (different tag or text changed)
  if (changed(oldTree, newTree)) {
    parent.replaceChild(renderNode(newTree), existingNode);
    return;
  }

  // 4. Same tag — diff attrs + children
  if (newTree.tag) {
    diffAttrs(existingNode, oldTree.attrs, newTree.attrs);
    diffChildren(existingNode, oldTree.children, newTree.children);
  }
}

/**
 * Determines whether two vNodes are fundamentally different and require
 * a full DOM replacement rather than an in-place update.
 *
 * Returns `true` if any of the following are true:
 * - Different JS types (e.g. string vs object — text vs element)
 * - Both are text nodes with different content
 * - Different HTML tag names
 * - Different `key` values
 *
 * @param {Object|string} a — Old vNode.
 * @param {Object|string} b — New vNode.
 * @returns {boolean} `true` if the node must be replaced.
 */
function changed(a, b) {
  if (typeof a !== typeof b) return true;
  if (typeof a === "string" && a !== b) return true;
  if (a.tag !== b.tag) return true;
  if (a.key !== b.key) return true;
  return false;
}

/**
 * Diffs the attributes of two vNodes and patches the real DOM element.
 *
 * Uses two `for...in` loops (no Set allocation):
 * 1. Iterate `newAttrs` — set anything changed or new.
 * 2. Iterate `oldAttrs` — remove anything not in `newAttrs`.
 *
 * Style objects are diffed property-by-property via `diffStyle()` rather
 * than replaced wholesale, avoiding unnecessary browser style recalculations.
 *
 * @param {Element} el          — The real DOM element to update.
 * @param {Object}  [oldAttrs]  — Previous attributes.
 * @param {Object}  [newAttrs]  — New attributes.
 */
function diffAttrs(el, oldAttrs = {}, newAttrs = {}) {
  // Iterate newAttrs: update anything changed or new
  for (const key in newAttrs) {
    if (key === "key") continue;
    const newVal = newAttrs[key];
    const oldVal = oldAttrs[key];
    if (key === "style" && typeof newVal === "object" && typeof oldVal === "object") {
      diffStyle(el, oldVal, newVal);
    } else if (oldVal !== newVal) {
      setAttr(el, key, newVal);
    }
  }

  // Iterate oldAttrs: remove anything not in newAttrs
  for (const key in oldAttrs) {
    if (key === "key") continue;
    if (!(key in newAttrs)) {
      removeAttr(el, key, oldAttrs[key]);
    }
  }
}

/**
 * Diffs individual CSS style properties between old and new style objects.
 *
 * Only changed properties are updated; removed properties are cleared to `""`.
 * This avoids replacing the entire `style` attribute, which would force the
 * browser to recalculate all style properties.
 *
 * @param {Element} el       — The real DOM element.
 * @param {Object}  oldStyle — Previous style properties.
 * @param {Object}  newStyle — New style properties.
 */
function diffStyle(el, oldStyle, newStyle) {
  for (const prop in newStyle) {
    if (oldStyle[prop] !== newStyle[prop]) {
      el.style[prop] = newStyle[prop];
    }
  }
  for (const prop in oldStyle) {
    if (!(prop in newStyle)) {
      el.style[prop] = "";
    }
  }
}

/**
 * Diffs the children of two vNodes and patches the parent DOM element.
 *
 * If any child has a `key` attribute, delegates to `diffKeyedChildren()`
 * for efficient reorder/reuse. Otherwise uses index-based comparison:
 * 1. Diff overlapping children pair-by-pair.
 * 2. Remove extra old children from the end (reverse order avoids index shifting).
 * 3. Append new children beyond the old list's length.
 *
 * @param {Element} parent         — The parent DOM element.
 * @param {Array}   [oldChildren]  — Previous child vNodes.
 * @param {Array}   [newChildren]  — New child vNodes.
 */
function diffChildren(parent, oldChildren = [], newChildren = []) {
  if (hasKeys(newChildren) || hasKeys(oldChildren)) {
    diffKeyedChildren(parent, oldChildren, newChildren);
    return;
  }

  const min = Math.min(oldChildren.length, newChildren.length);

  for (let i = 0; i < min; i++) {
    diff(parent, oldChildren[i], newChildren[i], i);
  }

  for (let i = oldChildren.length - 1; i >= min; i--) {
    const node = parent.childNodes[i];
    if (node) parent.removeChild(node);
  }

  for (let i = min; i < newChildren.length; i++) {
    parent.appendChild(renderNode(newChildren[i]));
  }
}

/**
 * Checks whether any child in a list has a `key` attribute.
 *
 * Used to decide between index-based diffing and keyed reconciliation.
 * Returns `true` as soon as one keyed child is found (early exit).
 *
 * @param {Array} children — Array of vNodes (or strings).
 * @returns {boolean} `true` if at least one child has a non-null key.
 */
function hasKeys(children) {
  for (let i = 0; i < children.length; i++) {
    if (typeof children[i] === "object" && children[i] !== null && children[i].key != null) {
      return true;
    }
  }
  return false;
}

/**
 * Keyed reconciliation algorithm for efficient list diffing.
 *
 * Instead of comparing children by index (which causes unnecessary
 * destroy/recreate when items reorder), this algorithm matches children
 * by their `key` attribute:
 *
 * 1. Build a `Map<key, { vNode, index }>` from old children (O(1) lookup).
 * 2. Walk new children — reuse matched DOM nodes (diff attrs/children in-place),
 *    create new nodes for unmatched keys.
 * 3. Remove old DOM nodes whose keys don't appear in the new list.
 * 4. Reorder DOM nodes to match the new order via `insertBefore()`.
 * 5. Trim any trailing extra DOM nodes.
 *
 * @param {Element} parent       — The parent DOM element.
 * @param {Array}   oldChildren  — Previous child vNodes (with keys).
 * @param {Array}   newChildren  — New child vNodes (with keys).
 */
function diffKeyedChildren(parent, oldChildren, newChildren) {
  const oldKeyMap = new Map();
  for (let i = 0; i < oldChildren.length; i++) {
    const child = oldChildren[i];
    if (typeof child === "object" && child !== null && child.key != null) {
      oldKeyMap.set(child.key, { vNode: child, index: i });
    }
  }

  const newDomOrder = [];
  const usedKeys = new Set();

  for (let i = 0; i < newChildren.length; i++) {
    const newChild = newChildren[i];
    const key = typeof newChild === "object" && newChild !== null ? newChild.key : null;

    if (key != null && oldKeyMap.has(key)) {
      const old = oldKeyMap.get(key);
      usedKeys.add(key);
      const domNode = parent.childNodes[old.index];

      if (domNode) {
        diffAttrs(domNode, old.vNode.attrs, newChild.attrs);
        diffChildren(domNode, old.vNode.children, newChild.children);
        newDomOrder.push(domNode);
      } else {
        newDomOrder.push(renderNode(newChild));
      }
    } else {
      newDomOrder.push(renderNode(newChild));
    }
  }

  for (const [key, old] of oldKeyMap) {
    if (!usedKeys.has(key)) {
      const domNode = parent.childNodes[old.index];
      if (domNode) parent.removeChild(domNode);
    }
  }

  for (let i = 0; i < newDomOrder.length; i++) {
    const node = newDomOrder[i];
    const current = parent.childNodes[i];

    if (!current) {
      parent.appendChild(node);
    } else if (current !== node) {
      parent.insertBefore(node, current);
    }
  }

  while (parent.childNodes.length > newDomOrder.length) {
    parent.removeChild(parent.lastChild);
  }
}

// ============================================================
// Attribute Helpers
// ============================================================

/**
 * Sets a single attribute, property, or event handler on a DOM element.
 *
 * Handles the following cases:
 * - **`null`/`undefined` values** → delegates to `removeAttr()`.
 * - **Event handlers** (`onClick`, `onInput`, etc.) → uses stable wrappers
 *   so the listener is attached once and the handler reference is swapped
 *   on updates, avoiding detach/reattach overhead per render.
 * - **Style objects** → sets each CSS property individually on `el.style`.
 * - **Boolean attributes** (`checked`, `disabled`) → sets both the DOM
 *   property and the HTML attribute for full browser compatibility.
 * - **`className`** → mapped to `el.setAttribute("class", ...)`.
 * - **`value`** → set as a DOM property (required for input elements).
 * - **All other attributes** → `el.setAttribute(key, value)`.
 *
 * @param {Element} el    — Target DOM element.
 * @param {string}  key   — Attribute name (e.g. `"class"`, `"onClick"`, `"style"`).
 * @param {*}       value — Attribute value.
 */
function setAttr(el, key, value) {
  if (value === null || value === undefined) {
    removeAttr(el, key, value);
    return;
  }

  // Event listeners: onClick -> click
  // Uses a stable wrapper so listeners aren't detached/reattached every render.
  if (key.startsWith("on") && typeof value === "function") {
    const event = key.slice(2).toLowerCase();
    if (!el.__listeners) {
      el.__listeners = {};
      el.__wrappers = {};
    }
    el.__listeners[event] = value;
    if (!el.__wrappers[event]) {
      el.__wrappers[event] = (e) => el.__listeners[event](e);
      el.addEventListener(event, el.__wrappers[event]);
    }
    return;
  }

  // Style object support
  if (key === "style" && typeof value === "object") {
    for (const [prop, val] of Object.entries(value)) {
      el.style[prop] = val;
    }
    return;
  }

  // Boolean attributes (checked, disabled, autofocus, etc.)
  if (typeof value === "boolean") {
    if (value) {
      el.setAttribute(key, "");
      el[key] = true;
    } else {
      el.removeAttribute(key);
      el[key] = false;
    }
    return;
  }

  // className shorthand
  if (key === "className") {
    el.setAttribute("class", value);
    return;
  }

  // Special: value (for inputs)
  if (key === "value") {
    el.value = value;
    return;
  }

  el.setAttribute(key, value);
}

/**
 * Removes a single attribute, property, or event handler from a DOM element.
 *
 * Handles the following cases:
 * - **Event handlers** → removes the stable wrapper via `removeEventListener`
 *   and cleans up `el.__listeners` and `el.__wrappers`.
 * - **`className`** → removes the `class` attribute.
 * - **`value`** → resets to `""`.
 * - **Boolean properties** (`checked`, `disabled`, `selected`) → sets to `false`.
 * - **Style objects** → removes the entire `style` attribute.
 * - **All other attributes** → `el.removeAttribute(key)`.
 *
 * @param {Element} el    — Target DOM element.
 * @param {string}  key   — Attribute name to remove.
 * @param {*}       value — The current value (used for type-checking style objects).
 */
function removeAttr(el, key, value) {
  if (key.startsWith("on")) {
    const event = key.slice(2).toLowerCase();
    if (el.__wrappers && el.__wrappers[event]) {
      el.removeEventListener(event, el.__wrappers[event]);
      delete el.__wrappers[event];
      if (el.__listeners) delete el.__listeners[event];
      return;
    }
  }

  if (key === "className") {
    el.removeAttribute("class");
    return;
  }

  if (key === "value") {
    el.value = "";
  } else if (key === "checked" || key === "disabled" || key === "selected") {
    el[key] = false;
  } else if (key === "style" && typeof value === "object") {
    el.removeAttribute("style");
    return;
  }

  el.removeAttribute(key);
}
