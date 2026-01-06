import j from "alpinejs";
let I = !1, B = null;
const g = {
  /**
   * Set debug mode
   * @param {boolean} debug
   */
  setDebug(i) {
    I = i;
  },
  /**
   * Set custom error reporter
   * @param {Function} reporter - (error, ...args) => void
   */
  setReporter(i) {
    B = i;
  },
  /**
   * Log info message
   * @param {...any} args - Arguments to log
   */
  info(...i) {
    I && typeof console < "u" && typeof console.log == "function" && console.log("[Uno Info]", ...i);
  },
  /**
   * Log warning message
   * @param {...any} args - Arguments to log
   */
  warn(...i) {
    I && typeof console < "u" && typeof console.warn == "function" && console.warn("[Uno Warn]", ...i);
  },
  /**
   * Log error message
   * @param {...any} args - Arguments to log
   */
  error(...i) {
    if (B)
      try {
        B(...i);
      } catch (t) {
        I && typeof console < "u" && typeof console.error == "function" && console.error("[Uno Error] Reporter failed:", t);
      }
    I && typeof console < "u" && typeof console.error == "function" && console.error("[Uno Error]", ...i);
  }
}, z = /* @__PURE__ */ new WeakMap();
function M(i, t) {
  if (!i) return !1;
  const n = String(t);
  let e = z.get(i);
  return e || (e = /* @__PURE__ */ new Set(), z.set(i, e)), e.has(n) ? !1 : (e.add(n), !0);
}
function P(i) {
  return i ? i.replace(/<script\b[^>]*>([\s\S]*?)<\/script>/gim, "").replace(/<img\b[^>]*>/gim, "").replace(/\s+on[a-z]+="[^"]*"/gim, "").replace(/\s+on[a-z]+='[^']*'/gim, "").replace(/\s+href="javascript:[^"]*"/gim, "").replace(/\s+href='javascript:[^']*'/gim, "").replace(/javascript:/gim, "") : "";
}
function V(i, t, { trusted: n = !1, alpine: e } = {}) {
  i && (i.innerHTML = n ? String(t ?? "") : P(String(t ?? "")), e?.initTree?.(i));
}
const O = {
  /**
   * Ensure a container element exists in the DOM, creating it if necessary.
   *
   * This method implements idempotent container management:
   * - Returns existing container if already present (by ID)
   * - Creates and injects new container if not found
   * - Automatically handles CSP nonce injection
   * - Safely handles errors with logging
   *
   * Container Creation Process:
   * 1. Check for existing element with matching ID
   * 2. If found, return it immediately (no duplicate creation)
   * 3. If not found:
   *    a. Create new div element
   *    b. Set the specified ID
   *    c. Auto-inject CSP nonce from global Uno config (if available)
   *    d. Apply all provided attributes
   *    e. Set innerHTML if content provided
   *    f. Append to specified parent (default: document.body)
   *
   * CSP Nonce Injection:
   * - Automatically reads nonce from window.Uno.options.csp.nonce
   * - Only injects if not already specified in attributes
   * - Enables inline styles/scripts in CSP-protected environments
   *
   * Error Handling:
   * - Wraps creation in try-catch
   * - Logs errors via Uno's logging system
   * - Returns null on failure (graceful degradation)
   *
   * Edge Cases:
   * - Parent null/undefined: Falls back to document.documentElement
   * - Duplicate ID: Returns existing element (prevents duplicates)
   * - Empty HTML: Creates container without content
   * - No attributes: Creates minimal container with just ID
   *
   * @param {string} id - Unique element ID (used for lookup and creation)
   * @param {string} [html=''] - Inner HTML content for the container
   * @param {Object} [attributes={}] - HTML attributes to set (key-value pairs)
   * @param {HTMLElement} [parent=document.body] - Parent element for appending
   * @returns {HTMLElement|null} The container element, or null if creation failed
   *
   * @example
   * // Simple container
   * const container = DomManager.ensureContainer('my-container')
   *
   * @example
   * // Container with content and attributes
   * const container = DomManager.ensureContainer(
   *   'tooltip-container',
   *   '<div class="tooltip">Content</div>',
   *   { 'x-data': '{}', 'class': 'fixed top-0' }
   * )
   *
   * @example
   * // Idempotent: Safe to call multiple times
   * const c1 = DomManager.ensureContainer('my-id')
   * const c2 = DomManager.ensureContainer('my-id')
   * console.log(c1 === c2) // true (same element)
   */
  ensureContainer(i, t = "", n = {}, e = document.body) {
    const o = document.getElementById(i);
    if (o) return o;
    try {
      const s = document.createElement("div");
      return s.id = i, window.Uno?.options?.csp?.nonce && !n.nonce && (n.nonce = window.Uno.options.csp.nonce), Object.entries(n).forEach(([a, c]) => {
        s.setAttribute(a, c);
      }), t && (s.innerHTML = t), e ? e.appendChild(s) : document.documentElement.appendChild(s), s;
    } catch (s) {
      return g.error(`DomManager: Failed to inject container #${i}`, s), null;
    }
  },
  /**
   * Inject CSS styles into the document head with CSP nonce support.
   *
   * This method provides safe, idempotent style injection with automatic
   * Content Security Policy (CSP) nonce handling.
   *
   * Features:
   * - **Idempotent**: Safe to call multiple times (won't create duplicates)
   * - **CSP Support**: Automatically injects nonce attribute
   * - **Error Handling**: Graceful failure with logging
   * - **Head Injection**: Styles go in <head> for proper CSS precedence
   *
   * Injection Process:
   * 1. Check if style element with this ID already exists
   * 2. If exists, return immediately (no duplicate creation)
   * 3. If not exists:
   *    a. Create new <style> element
   *    b. Set the specified ID
   *    c. Auto-inject CSP nonce from window.Uno.options.csp.nonce
   *    d. Set CSS content as textContent (safe from XSS)
   *    e. Append to document.head
   *
   * CSP Compliance:
   * - Reads nonce from global Uno configuration
   * - Enables inline styles in strict CSP environments
   * - Example CSP header: `style-src 'nonce-abc123'`
   *
   * Performance:
   * - ID-based duplicate prevention (O(1) lookup)
   * - No style re-processing on subsequent calls
   * - Suitable for dynamic style injection
   *
   * Use Cases:
   * - Component-specific styles
   * - Theme customization
   * - Dynamic CSS generation
   * - Plugin style injection
   *
   * @param {string} id - Unique style element ID (prevents duplicates)
   * @param {string} css - CSS content to inject (raw CSS text)
   *
   * @example
   * // Inject component styles
   * DomManager.injectStyles('tooltip-styles', `
   *   .tooltip { position: fixed; z-index: 9999; }
   *   .tooltip-arrow { width: 0; height: 0; }
   * `)
   *
   * @example
   * // Idempotent: Safe to call multiple times
   * DomManager.injectStyles('my-styles', '.foo { color: red; }')
   * DomManager.injectStyles('my-styles', '.foo { color: red; }') // No-op
   *
   * @see {@link https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP|MDN: Content Security Policy}
   */
  injectStyles(i, t) {
    if (!document.getElementById(i))
      try {
        const n = document.createElement("style");
        n.id = i, window.Uno?.options?.csp?.nonce && n.setAttribute("nonce", window.Uno.options.csp.nonce), n.textContent = t, document.head.appendChild(n);
      } catch (n) {
        g.error(`DomManager: Failed to inject styles #${i}`, n);
      }
  },
  // Exports
  setHTML: V,
  // Secure HTML setting
  sanitize: P
  // HTML sanitization
};
function L(i) {
  return Object.prototype.toString.call(i) === "[object Object]";
}
function E(i, t) {
  if (!L(i)) return L(t) ? { ...t } : t;
  if (!L(t)) return t === void 0 ? { ...i } : t;
  const n = { ...i };
  for (const [e, o] of Object.entries(t)) {
    if (L(o) && L(n[e])) {
      n[e] = E(n[e], o);
      continue;
    }
    if (L(o)) {
      n[e] = E({}, o);
      continue;
    }
    if (typeof o == "string" && o.trim().indexOf("+") === 0 && typeof n[e] == "string") {
      n[e] = `${n[e]} ${o.trim().slice(1)}`.trim();
      continue;
    }
    n[e] = o;
  }
  return n;
}
const C = {
  BASE: 0,
  DROPDOWN: 1e3,
  STICKY: 1100,
  FIXED: 1200,
  BOX_BACKDROP: 1300,
  BOX: 1400,
  POPOVER: 1500,
  NOTIFY: 1600,
  TOOLTIP: 1700,
  MAX: 9999
}, k = (i) => {
  if (!i) return C.BASE;
  const t = i.toUpperCase();
  return C[t] !== void 0 ? C[t] : C.BASE;
};
function U(i, t, n) {
  if (!i || typeof i.store != "function") {
    console.error("[Uno] Alpine.store is not available");
    return;
  }
  const e = i.store("uno");
  e ? i.store("uno", {
    ...e,
    [t]: n
  }) : i.store("uno", { [t]: n }), i.store(t) || i.store(t, i.store("uno")?.[t]);
}
const A = (i) => {
  const t = i || "center", n = {
    center: "items-center justify-center",
    "top-center": "items-start justify-center",
    "top-left": "items-start justify-start",
    "top-right": "items-start justify-end",
    "bottom-center": "items-end justify-center",
    "bottom-left": "items-end justify-start",
    "bottom-right": "items-end justify-end",
    "center-left": "items-center justify-start",
    "center-right": "items-center justify-end"
  };
  return n[t] || n.center;
}, F = (i) => {
  const t = i || "top-right", n = {
    "top-right": "top-0 right-0",
    "top-left": "top-0 left-0",
    "top-center": "top-0 left-0 right-0",
    // Full width for center alignment
    center: "inset-0",
    // Full screen for center
    "center-right": "top-1/2 right-0 -translate-y-1/2",
    "center-left": "top-1/2 left-0 -translate-y-1/2",
    "bottom-right": "bottom-0 right-0",
    "bottom-left": "bottom-0 left-0",
    "bottom-center": "bottom-0 left-0 right-0"
  };
  return n[t] || n["top-right"];
}, W = (i) => {
  const t = i || "top-right";
  let n = "flex-col", e = "items-end", o = "space-y-2";
  return t.includes("bottom") && (n = "flex-col-reverse", o = "space-y-reverse space-y-2"), t.includes("left") ? e = "items-start" : t.includes("center") ? t === "center" ? e = "items-center justify-center" : e = "items-center" : e = "items-end", { flexDir: n, alignItems: e, spaceY: o };
};
function X(i, t, n) {
  const e = i.getBoundingClientRect(), o = t.getBoundingClientRect(), s = o.width, a = o.height, c = e.width, r = e.height, l = { top: e.top, left: e.left }, u = window.innerWidth, p = window.innerHeight, f = 0, y = 0, b = p, v = u, x = (T, q) => {
    if (q) {
      if (T + s > v)
        return v - s - 10;
      if (T < y)
        return y + 10;
    } else {
      if (T + a > b)
        return b - a - 10;
      if (T < f)
        return f + 10;
    }
    return T;
  };
  let w, d, m = n;
  switch (n) {
    case "left":
      d = l.left - s - 5, d < y && (d = l.left + c + 5, m = "right"), w = x(l.top + r / 2 - a / 2, !1);
      break;
    case "top":
      w = l.top - a - 8, w < f && (w = l.top + r + 8, m = "bottom"), d = l.left + (c - s) / 2, d = x(d, !0);
      break;
    case "right":
      d = l.left + c + 5, d + s > v && (d = l.left - s - 5, m = "left"), w = x(l.top + r / 2 - a / 2, !1);
      break;
    case "bottom":
    default:
      w = l.top + r + 8, w + a > b && (w = l.top - a - 8, m = "top"), d = l.left + (c - s) / 2, d = x(d, !0);
      break;
  }
  return { top: w, left: d, position: m };
}
const Y = () => `
  <!-- Backdrop / Overlay -->
  <div id="uno-box-overlay"
       :class="Uno.getStyle('box', 'overlay')"
       x-show="$store.uno.box.items.length > 0"
       x-transition:enter="ease-out duration-300"
       x-transition:enter-start="opacity-0"
       x-transition:enter-end="opacity-100"
       x-transition:leave="ease-in duration-200"
       x-transition:leave-start="opacity-100"
       x-transition:leave-end="opacity-0"
       @click="$store.uno.box.items.length && $store.uno.box.items[$store.uno.box.items.length - 1].closeOnOverlayClick && $store.uno.box.items[$store.uno.box.items.length - 1].close()"></div>

  <div :class="Uno.getStyle('box', 'dialogWrapper')"
       @click.self="$store.uno.box.items.length && $store.uno.box.items[$store.uno.box.items.length - 1].closeOnOverlayClick && $store.uno.box.items[$store.uno.box.items.length - 1].close()">
    <template x-for="item in $store.uno.box.items" :key="item.id">
      <div :id="'uno-box-outer-container' + item.id"
           :class="[Uno.getStyle('box', 'itemWrapper'), $store.uno.box.getPositionClasses(item.position)]"
           :style="'z-index:' + (Uno.options?.ui?.layers?.box || getZIndex('BOX'))"
           x-show="item.isOpen"
           x-transition:enter="ease-out duration-300"
           x-transition:enter-start="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
           x-transition:enter-end="opacity-100 translate-y-0 sm:scale-100"
           x-transition:leave="ease-in duration-200"
           x-transition:leave-start="opacity-100 translate-y-0 sm:scale-100"
           x-transition:leave-end="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
           role="dialog"
           :aria-modal="item.modal !== false ? 'true' : 'false'"
           :aria-labelledby="item.title ? 'uno-box-title-' + item.id : undefined"
           :aria-describedby="'uno-box-content-' + item.id"
           @keydown.escape.window="item.isOpen && item.closeOnEsc && item.close()">

        <div :id="'uno-box-container' + item.id"
             :class="[Uno.getStyle('box', 'container'), item.widthClass, item.heightClass, item.containerClass]"
             :style="(item.width ? 'width:' + item.width + 'px; max-width:' + item.width + 'px;' : '') + (item.height ? ' height:' + item.height + 'px;' : '')"
             x-uno-trap>

          <!-- Close Button -->
          <button type="button"
                  x-show="item.showCloseButton"
                  :class="Uno.getStyle('box', 'closeBtnFloating')"
                  @click="item.close()"
                  aria-label="Close">
            <svg :class="Uno.getStyle('box', 'closeIconSmall')" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>

          <!-- Content Area -->
          <div :class="Uno.getStyle('box', 'content')">
            <h3 :id="'uno-box-title-' + item.id"
                :class="Uno.getStyle('box', 'title')"
                x-text="item.title"
                x-show="item.title"></h3>
            <div class="mt-2" :id="'uno-box-content-' + item.id">
              <template x-if="!item.hxUrl">
                <div x-html="item.content"></div>
              </template>
              <template x-if="item.hxUrl">
                <div x-init="htmx.process($el); $el.addEventListener('htmx:afterSwap', () => $store.uno.box.onHtmxContentLoaded($el, item))"
                     :hx-get="item.hxMethod === 'GET' ? item.hxUrl : undefined"
                     :hx-post="item.hxMethod === 'POST' ? item.hxUrl : undefined"
                     :hx-trigger="item.hxTrigger"
                     hx-swap="innerHTML"
                     :hx-headers="item.hxHeaders ? JSON.stringify(item.hxHeaders) : null">
                  <div :class="Uno.getStyle('box', 'loadingWrap')">
                    <svg :class="Uno.getStyle('box', 'spinner')" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                  </div>
                </div>
              </template>
            </div>
          </div>

          <!-- Footer with Buttons -->
          <div :class="Uno.getStyle('box', 'footer')"
               x-show="item.showOkButton || item.showCancelButton">
            <button type="button"
                    :id="'uno-box-ok' + item.id"
                    :class="Uno.getStyle('box', 'okBtn')"
                    x-show="item.showOkButton"
                    x-text="item.okText"
                    @click="item.onOk ? (item.onOk() !== false && item.close()) : item.close()"></button>
            <button type="button"
                    :id="'uno-box-cancel' + item.id"
                    :class="Uno.getStyle('box', 'cancelBtn')"
                    x-show="item.showCancelButton"
                    x-text="item.cancelText"
                    @click="item.onCancel ? (item.onCancel() !== false && item.close()) : item.close()"></button>
          </div>
        </div>
      </div>
    </template>
  </div>
`, Z = (i) => ({
  items: [],
  /**
   * Initialize the box component
   * Ensures container exists when store is initialized
   */
  init() {
    this.ensureContainer();
  },
  // ============================================================
  // Public API - Convenience Methods
  // ============================================================
  /**
   * Show an alert dialog
   * @param {string} message - Alert message
   * @param {object} options - Additional options
   * @returns {object} Box item instance
   */
  alert(t, n = {}) {
    const e = typeof n == "string" ? { title: n } : n;
    return this.open("dialog", {
      content: t,
      showOkButton: !0,
      ...e
    });
  },
  /**
   * Show a confirm dialog
   * @param {string} message - Confirmation message
   * @param {object} options - Additional options
   * @returns {Promise<boolean>} Resolves true if OK, false otherwise
   */
  confirm(t, n = {}) {
    return new Promise((e) => {
      this.open("confirm", {
        content: t,
        onOk: () => e(!0),
        onCancel: () => e(!1),
        onClose: () => e(!1),
        // Closing without choice = cancel
        ...n
      });
    });
  },
  /**
   * Show a prompt dialog
   * @param {string} message - Prompt message
   * @param {string} defaultValue - Default input value
   * @param {object} options - Additional options
   * @returns {Promise<string|null>} Resolves with input value or null
   */
  prompt(t, n = "", e = {}) {
    return new Promise((o) => {
      const s = "prompt-" + Math.random().toString(36).slice(2, 11), a = h.getStyle("box", "promptWrapper"), c = h.getStyle("box", "promptText"), r = h.getStyle("box", "promptInput"), l = h.dom.sanitize(t), u = h.dom.sanitize(n), p = `
        <div class="${a}">
          <p class="${c}">${l}</p>
          <input id="${s}" type="text" class="${r}" value="${u}">
        </div>
      `;
      this.open("confirm", {
        title: e.title || "Prompt",
        content: p,
        _isTrustedPrompt: !0,
        // Internal flag for prompt content
        onOk: () => {
          const f = document.getElementById(s);
          o(f ? f.value : null);
        },
        onCancel: () => o(null),
        ...e
      });
    });
  },
  /**
   * Show a dialog (alias for open)
   * @param {object} options - Dialog options
   * @returns {object} Box item instance
   */
  dialog(t = {}) {
    return this.open("dialog", t);
  },
  // ============================================================
  // Core Methods
  // ============================================================
  /**
   * Open a box dialog
   * Uses 4-level configuration merge: Plugin Defaults -> Global -> Plugin Config -> Instance
   * @param {string|object} typeOrOpts - Type string or options object
   * @param {object} opts - Options if first arg is type
   * @returns {object} Box item instance
   */
  open(t, n = {}) {
    let e = "dialog", o = n;
    typeof t == "object" ? (o = t, e = o.type || "dialog") : e = t;
    const s = h.getConfig("box", o), a = s.boxId || Math.random().toString(36).slice(2, 10);
    let c = s.showOkButton, r = s.showCancelButton;
    e === "confirm" ? (c = c !== !1, r = r !== !1) : e === "dialog" && (c = c || !1, r = r || !1);
    const l = {
      id: a,
      type: e,
      title: s.title || "",
      content: s._isTrustedPrompt ? s.content || "" : h.dom.sanitize(s.content || ""),
      // Security: Always sanitize content (trustedHTML removed)
      // Exception: Internal prompt content is pre-escaped
      _isTrustedPrompt: s._isTrustedPrompt || !1,
      width: s.width !== void 0 ? s.width : s.widthClass ? void 0 : s.defaultWidth || 400,
      height: s.height,
      widthClass: s.widthClass,
      heightClass: s.heightClass,
      containerClass: s.containerClass,
      position: s.position || s.defaultPosition || "center",
      modal: s.modal !== void 0 ? s.modal : !0,
      closeOnEsc: s.closeOnEsc !== void 0 ? s.closeOnEsc : !1,
      closeOnOverlayClick: s.closeOnOverlayClick !== void 0 ? s.closeOnOverlayClick : !1,
      showCloseButton: s.showCloseButton !== void 0 ? s.showCloseButton : !0,
      showOkButton: c,
      showCancelButton: r,
      okText: s.okText || "OK",
      cancelText: s.cancelText || "Cancel",
      onOk: s.onOk,
      onCancel: s.onCancel,
      onBeforeClose: s.onBeforeClose,
      onClose: s.onClose,
      onContentLoad: s.onContentLoad,
      hxUrl: s.hxUrl,
      hxMethod: s.hxMethod || "GET",
      hxHeaders: s.hxHeaders,
      hxTrigger: s.hxTrigger || "load",
      // Internal state - start closed to trigger animation
      isOpen: !1,
      /**
       * Close this box dialog
       */
      close() {
        this.onBeforeClose && this.onBeforeClose(), this.isOpen = !1;
        const f = typeof document < "u" ? document.getElementById("uno-box-" + this.id) : null, y = () => {
          const x = i.store("uno")?.box;
          x && x.remove(this.id), this.onClose && this.onClose();
        };
        if (!f) {
          y();
          return;
        }
        let b = !1;
        const v = (x) => {
          x && x.target !== f || b || (b = !0, f.removeEventListener("transitionend", v), y());
        };
        f.addEventListener("transitionend", v), setTimeout(v, 350);
      }
    };
    this.items.push(l), this.ensureContainer();
    const u = this.items[this.items.length - 1], p = () => {
      u.isOpen = !0;
    };
    return i && typeof i.nextTick == "function" ? i.nextTick(p) : setTimeout(p, 0), u;
  },
  /**
   * Close a specific box by ID
   * @param {string} id - Box ID
   */
  close(t) {
    const n = this.items.find((e) => e.id === t);
    n && n.close();
  },
  /**
   * Remove a box from the items array
   * @param {string} id - Box ID
   */
  remove(t) {
    const n = this.items.findIndex((e) => e.id === t);
    n !== -1 && this.items.splice(n, 1);
  },
  /**
   * Close all open boxes
   */
  closeAll() {
    this.items.forEach((t) => {
      t.onBeforeClose && t.onBeforeClose(), t.onClose && t.onClose();
    }), this.items = [];
  },
  /**
   * Get position CSS classes for 9-point grid
   * @param {string} position - Position name
   * @returns {string} CSS classes
   */
  getPositionClasses(t) {
    return A(t);
  },
  /**
   * Handle HTMX content load completion
   * @param {HTMLElement} el - Content element
   * @param {object} item - Box item
   */
  onHtmxContentLoaded(t, n) {
    n.onContentLoad && typeof n.onContentLoad == "function" && n.onContentLoad(t), window.Alpine && window.Alpine.initTree && window.Alpine.initTree(t);
  },
  /**
   * Ensure the box container exists in the DOM
   */
  ensureContainer() {
    const t = "uno-box-manager", n = typeof document < "u" && document.getElementById(t);
    if (n)
      return i && typeof i.initTree == "function" && !n._x_dataStack && i.initTree(n), n;
    const o = {
      "x-data": "{}",
      class: "relative",
      style: `z-index: ${h.options?.ui?.layers?.box || k("BOX")};`,
      "x-show": "$store.uno.box.items.length > 0",
      "x-cloak": ""
      // Hide until Alpine is ready
    }, s = Y(), a = h.dom.ensureContainer(t, s, o);
    return a && i && typeof i.initTree == "function" && i.initTree(a), a;
  }
}), N = [
  "a[href]",
  "area[href]",
  "input:not([disabled])",
  "select:not([disabled])",
  "textarea:not([disabled])",
  "button:not([disabled])",
  "iframe",
  "object",
  "embed",
  "[contenteditable]",
  '[tabindex]:not([tabindex^="-"])'
].join(",");
function K(i, t = {}) {
  const n = (r) => {
    if (!(r instanceof HTMLElement)) return !1;
    const l = window.getComputedStyle(r);
    return l.display === "none" || l.visibility === "hidden" ? !1 : r.getClientRects().length > 0;
  }, e = () => Array.from(i.querySelectorAll(N)).filter((r) => n(r)), o = () => e()[0], s = () => {
    const r = e();
    return r[r.length - 1];
  }, a = (r) => {
    if (r.key === "Tab") {
      if (r.shiftKey) {
        if (document.activeElement === o()) {
          r.preventDefault();
          const l = s();
          l && l.focus();
        }
      } else if (document.activeElement === s()) {
        r.preventDefault();
        const l = o();
        l && l.focus();
      }
    }
  }, c = () => {
    if (i.contains(document.activeElement)) return;
    const r = o();
    r && r.focus();
  };
  return typeof requestAnimationFrame == "function" ? requestAnimationFrame(() => requestAnimationFrame(c)) : setTimeout(c, 0), i.addEventListener("keydown", a), () => {
    i.removeEventListener("keydown", a);
  };
}
const G = (i) => {
  i.directive("uno-trap", (t, { expression: n }, { cleanup: e }) => {
    const o = K(t);
    e(o);
  }), i.directive("uno-box", (t, { expression: n }, { evaluateLater: e, cleanup: o }) => {
    if (!n) {
      console.warn("[Uno] x-uno-box directive requires a configuration object.");
      return;
    }
    const s = e(n), a = (c) => {
      c.preventDefault();
      const r = i.store("uno")?.box;
      if (!r) {
        console.warn("[Uno] Box store not found. Is Uno initialized?");
        return;
      }
      s((l) => {
        l && typeof l == "object" ? r.open(l) : console.warn("[Uno] x-uno-box directive value must be an object", l);
      });
    };
    t.addEventListener("click", a), o(() => t.removeEventListener("click", a));
  });
};
const J = (i, t = {}) => {
  const n = {
    defaultWidth: 400,
    defaultPosition: "center",
    okText: "OK",
    cancelText: "Cancel",
    closeOnEsc: !1,
    closeOnOverlayClick: !1
  };
  h.registerPluginDefaults("box", { ...n, ...t.plugins?.box, ...t.box });
  const e = Z(i);
  U(i, "box", e), G(i);
};
class Q {
  /**
   * Create an UnoInstance.
   *
   * @param {Object} [options={}] - Global configuration options
   */
  constructor(t = {}) {
    this.version = "0.6.0", this.options = t, this.pluginDefaults = {}, this.plugins = [], this.initialized = !1, this._htmxRequestHandlers = [], this._htmxActiveXhrs = typeof WeakSet == "function" ? /* @__PURE__ */ new WeakSet() : null, this.htmx = {
      /**
       * Check if HTMX is available.
       * @returns {boolean} True if HTMX is loaded and available
       */
      available: () => typeof window < "u" && typeof window.htmx?.ajax == "function",
      /**
       * Wrapper for htmx.ajax.
       *
       * @param {string} method - HTTP method (GET, POST, etc.)
       * @param {string} url - The URL to request
       * @param {Object} [ajaxOptions={}] - HTMX ajax options
       * @returns {Promise} Promise resolving to the result
       */
      ajax: (o, s, a = {}) => typeof window > "u" || typeof window.htmx?.ajax != "function" ? Promise.reject(new Error("htmx.ajax is not available")) : window.htmx.ajax(o, s, a),
      /**
       * Fetch an HTML fragment without inserting it into the DOM.
       * Useful for getting content to display in a modal or tooltip.
       *
       * @param {string} method - HTTP method
       * @param {string} url - The URL to request
       * @param {Object} [ajaxOptions={}] - HTMX ajax options
       * @returns {Promise<string>} The HTML fragment string
       */
      fragment: async (o, s, a = {}) => {
        if (typeof document > "u") throw new Error("document is not available");
        const c = document.createElement("div");
        let r = "";
        const l = (u) => {
          u.preventDefault(), r = u.detail.serverResponse, u.detail.shouldSwap = !1;
        };
        c.addEventListener("htmx:beforeSwap", l);
        try {
          const u = a.swap || "innerHTML";
          return await this.htmx.ajax(o, s, { ...a, target: c, swap: u }), r;
        } catch (u) {
          throw g.error("HTMX fragment fetch failed:", { method: o, url: s, error: u }), u;
        } finally {
          c.removeEventListener("htmx:beforeSwap", l), c.remove();
        }
      },
      /**
       * Register handlers for HTMX request lifecycle events.
       *
       * @param {Object} handlers - The handlers object
       * @param {Function} [handlers.start] - Called when a request starts
       * @param {Function} [handlers.end] - Called when a request ends (success, error, or timeout)
       */
      onRequest: (o = {}) => {
        if (!o || typeof o != "object" || (this._htmxRequestHandlers.push(o), typeof window > "u" || typeof document > "u") || !window.htmx || !M(document, "uno:htmx:request-hooks")) return;
        const s = (l) => l?.detail?.xhr, a = (l, u) => {
          for (const p of this._htmxRequestHandlers) {
            const f = p?.[l];
            if (typeof f == "function")
              try {
                f(u);
              } catch (y) {
                g.error("Uno.htmx.onRequest handler failed", y);
              }
          }
        }, c = (l) => {
          const u = s(l);
          if (u && this._htmxActiveXhrs) {
            if (this._htmxActiveXhrs.has(u)) return;
            this._htmxActiveXhrs.add(u);
          }
          a("start", l);
        }, r = (l) => {
          const u = s(l);
          if (u && this._htmxActiveXhrs) {
            if (!this._htmxActiveXhrs.has(u)) return;
            this._htmxActiveXhrs.delete(u);
          }
          a("end", l);
        };
        document.addEventListener("htmx:beforeRequest", c), document.addEventListener("htmx:afterRequest", r), document.addEventListener("htmx:responseError", r), document.addEventListener("htmx:sendError", r), document.addEventListener("htmx:timeout", r), document.addEventListener("htmx:xhr:progress", (l) => a("progress", l));
      },
      // Global HTMX event binding for notify functionality - progress bar only
      bindNotifyToHtmx: () => {
        typeof window > "u" || typeof document > "u" || !window.htmx || !window.Alpine || M(document, "uno:htmx:bind-notify") && (document.addEventListener("htmx:beforeRequest", () => {
          window.Alpine.store("uno").notify.startProgress();
        }), document.addEventListener("htmx:afterRequest", () => {
          window.Alpine.store("uno").notify.stopProgress();
        }), document.addEventListener("htmx:responseError", () => {
          window.Alpine.store("uno").notify.stopProgress(), window.Alpine.store("uno").notify.error("Request failed");
        }), document.addEventListener("htmx:sendError", () => {
          window.Alpine.store("uno").notify.stopProgress(), window.Alpine.store("uno").notify.error("Request failed");
        }), document.addEventListener("htmx:timeout", () => {
          window.Alpine.store("uno").notify.stopProgress(), window.Alpine.store("uno").notify.warning("Request timed out");
        }));
      }
    };
    const e = {
      debug: !1,
      csp: {
        nonce: null
        // For Content Security Policy
      },
      onError: null,
      // Custom error handler
      plugins: {},
      ui: {
        base: "rounded-lg shadow-md border border-gray-200 dark:border-gray-700",
        layers: { box: C.BOX, tooltip: C.TOOLTIP, notify: C.NOTIFY },
        box: {
          overlay: "fixed inset-0 bg-gray-500/75 transition-opacity",
          dialogWrapper: "fixed inset-0 z-10 w-screen overflow-y-auto",
          itemWrapper: "flex min-h-full p-4 text-center sm:p-0 pointer-events-none",
          container: "relative transform overflow-hidden {base} bg-white dark:bg-gray-900 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg flex flex-col max-h-[90vh] pointer-events-auto outline-none",
          cancelBtn: "mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:mt-0 sm:w-auto",
          okBtn: "inline-flex w-full justify-center rounded-md bg-blue-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 sm:ml-3 sm:w-auto",
          title: "text-base font-semibold leading-6 text-gray-900 dark:text-slate-200",
          closeBtnFloating: "absolute top-2 right-2 z-10 text-slate-400 hover:text-slate-500 focus:outline-none",
          closeIconSmall: "h-5 w-5",
          content: "px-4 pb-4 pt-5 sm:p-6 sm:pb-4 overflow-auto flex-1 text-slate-600 dark:text-slate-300",
          loadingWrap: "flex justify-center items-center h-full",
          spinner: "animate-spin h-8 w-8 text-blue-500",
          footer: "bg-gray-50 px-4 py-3 sm:flex sm:flex-row-reverse sm:px-6 border-t border-slate-200 dark:border-slate-700 dark:bg-slate-800",
          promptWrapper: "space-y-3",
          promptText: "text-slate-700 dark:text-slate-200",
          promptInput: "w-full px-3 py-2 border rounded focus:ring-2 focus:ring-blue-500 outline-none"
        },
        notify: {
          container: "fixed flex gap-2 p-4 pointer-events-none",
          item: "pointer-events-auto max-w-[calc(100vw-2rem)] overflow-hidden {base}",
          content: "p-4",
          row: "flex items-start",
          iconWrap: "flex-shrink-0",
          textWrap: "ml-3 min-w-0 flex-1 pt-0.5",
          text: "text-sm font-medium text-gray-900 dark:text-slate-100 break-words",
          closeWrap: "ml-4 flex flex-shrink-0",
          closeBtn: "inline-flex rounded-md bg-white dark:bg-slate-900 text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2",
          closeIcon: "h-5 w-5",
          progressWrap: "fixed top-0 left-0 right-0 h-1 bg-gray-200 dark:bg-gray-700 pointer-events-none",
          progressBar: "h-full bg-gradient-to-r from-blue-500 to-indigo-600 shadow-sm shadow-blue-500/50",
          types: {
            success: "bg-green-100 dark:bg-green-950",
            error: "bg-red-100 dark:bg-red-950",
            warning: "bg-yellow-100 dark:bg-yellow-950",
            info: "bg-blue-100 dark:bg-blue-950"
          },
          icons: {
            success: "text-green-500 dark:text-green-300",
            error: "text-red-500 dark:text-red-300",
            warning: "text-yellow-500 dark:text-yellow-300",
            info: "text-blue-500 dark:text-blue-300"
          }
        },
        tooltip: {
          content: "fixed px-3 py-2 text-sm font-medium {base} whitespace-nowrap bg-white text-gray-900 dark:bg-slate-800 dark:text-white",
          body: "pointer-events-auto"
        },
        tab: {
          active: "border-b-2 border-blue-500",
          inactive: "text-slate-500 hover:text-slate-700",
          loading: "animate-pulse h-32 w-full rounded-lg bg-gray-200 dark:bg-slate-700",
          error: "p-4 text-red-600 dark:text-red-400"
        },
        valid: {
          wrapper: "relative w-full",
          iconBase: "absolute right-0 top-0 bottom-0 flex items-center justify-center w-10 pointer-events-none transition-opacity duration-200",
          inputError: "!border-red-500",
          inputValid: "!border-green-500",
          iconError: "text-red-500 dark:text-red-400",
          iconValid: "text-green-500 dark:text-green-400",
          stateInvalid: "is-invalid",
          stateValid: "is-valid",
          stateValidating: "validating",
          inputPaddingIcon: "pr-10",
          helperText: "text-red-600 dark:text-red-400 text-xs"
        }
      }
    };
    this.options = E(e, t);
  }
  /**
   * ✅ Unified configuration retrieval method
   * Supports 4-level merge: Plugin Defaults -> Global User Config -> Global Plugin Config -> Direct Global Config -> Instance Config
   * Note: We use mergeDeep to ensure deep merging of options.
   *
   * Level 1: Plugin Defaults (this.pluginDefaults[componentName])
   * Level 2: Global User Config (this.options.plugins[componentName])
   * Level 3: Direct Global Config (this.options[componentName]) - Legacy/Shorthand
   * Level 4: Instance Config (instanceOpts)
   *
   * @param {string} componentName - Name of the component
   * @param {Object} [instanceOpts={}] - Instance-specific options
   * @returns {Object} The merged configuration object
   */
  getConfig(t, n = {}) {
    const e = this.options || {};
    let o = this.pluginDefaults[t] ? { ...this.pluginDefaults[t] } : {};
    return e.plugins && e.plugins[t] && (o = E(o, e.plugins[t])), e[t] && (o = E(o, e[t])), n && (o = E(o, n)), o;
  }
  /**
   * ✅ Register plugin defaults
   *
   * @param {string} name - Component name
   * @param {Object} defaults - Default configuration object
   */
  registerPluginDefaults(t, n) {
    this.pluginDefaults[t] = n;
  }
  _getUiValue(t, n) {
    const e = this.options.ui?.[t];
    if (!e) return;
    if (typeof n != "string" || n.indexOf(".") === -1) return e?.[n];
    const o = n.split(".").filter(Boolean);
    let s = e;
    for (const a of o) {
      if (s == null) return;
      s = s[a];
    }
    return s;
  }
  _resolveUiTokens(t) {
    if (typeof t != "string") return t;
    const n = this.options.ui?.base;
    return typeof n == "string" && t.indexOf("{base}") !== -1 ? t.replaceAll("{base}", n) : t;
  }
  /**
   * Get style class from configuration or fallback
   * Support "+" prefix to append class to default
   * @param {string} comp - Component name
   * @param {string} elem - Element name
   * @param {string} defaultCls - Default class string
   * @returns {string} Final class string
   */
  getStyle(t, n, e) {
    const o = this._getUiValue(t, n);
    return o === void 0 ? this._resolveUiTokens(e) : typeof o == "string" && o.trim().indexOf("+") === 0 ? this._resolveUiTokens(`${e} ${o.trim().slice(1)}`.trim()) : this._resolveUiTokens(o);
  }
  /**
   * Set global configuration (before or after init)
   * @param {Object} config - Global configuration object
   */
  setConfig(t) {
    return this.options = E(this.options, t), this.initialized && g.info("Uno configuration updated", this.options), this;
  }
  /**
   * Register a plugin
   * @param {Function} plugin - Plugin function (Alpine, options) => void
   * @param {Object} options - Plugin options
   * @returns {UnoInstance}
   */
  use(t, n = {}) {
    return this.initialized && g.warn("Uno.use() called after initialization. Plugin might not work correctly."), this.plugins.push({ plugin: t, options: n }), this;
  }
  /**
   * Initialize Uno
   * @param {Object|Function} alpineOrOptions - Alpine instance or Configuration options
   * @param {Object} options - Configuration options (if first arg is Alpine)
   * @returns {Object} Alpine instance
   */
  init(t = {}, n = {}) {
    if (this.initialized)
      return g.warn("Uno has already been initialized."), j;
    let e = j, o = n;
    return !e && typeof window < "u" && window.Alpine && (e = window.Alpine), t && (t.version || typeof t.start == "function" || typeof t.store == "function") ? e = t : o = t || n, this.options = E(this.options, o), g.setDebug(this.options.debug !== !1), this.options.onError && g.setReporter(this.options.onError), g.info(`Uno v${this.version} initializing...`), typeof window < "u" && typeof document < "u" && window.htmx && typeof e?.initTree == "function" && M(document, "uno:htmx-afterSwap:initTree") && document.addEventListener("htmx:afterSwap", (s) => {
      const a = s?.detail?.elt;
      a && e.initTree(a);
    }), this.plugins.forEach(({ plugin: s, options: a }) => {
      try {
        const c = E(this.options, a);
        s(e, c);
      } catch (c) {
        g.error("Failed to load plugin", c);
      }
    }), this._setupGlobalErrorHandlers(), this.initialized = !0, g.info("Uno initialized"), e;
  }
  /**
   * Setup global error boundaries for unhandled errors and promise rejections.
   * This ensures that all errors are properly logged and reported through the
   * configured error handler.
   *
   * Error Handling Strategy:
   * 1. Catches unhandled Promise rejections
   * 2. Catches uncaught JavaScript errors
   * 3. Reports errors through log.error() which triggers onError callback
   * 4. Prevents duplicate error handling using WeakSet
   *
   * @private
   */
  _setupGlobalErrorHandlers() {
    if (!(typeof window > "u") && M(window, "uno:global-error-handlers")) {
      const t = typeof WeakSet == "function" ? /* @__PURE__ */ new WeakSet() : null;
      window.addEventListener("unhandledrejection", (n) => {
        const e = n.reason;
        if (t && e instanceof Error) {
          if (t.has(e)) return;
          t.add(e);
        }
        g.error("Unhandled Promise Rejection:", e);
      }), window.addEventListener("error", (n) => {
        const e = n.error;
        if (t && e instanceof Error) {
          if (t.has(e)) return;
          t.add(e);
        }
        g.error("Uncaught Error:", {
          message: n.message,
          filename: n.filename,
          lineno: n.lineno,
          colno: n.colno,
          error: e
        });
      }), g.debug("Global error handlers installed");
    }
  }
}
const h = new Q(), tt = {
  box: {
    ok: "OK",
    cancel: "Cancel",
    loading: "Loading",
    close: "Close"
  },
  valid: {
    required: "Required",
    pattern: "Invalid format",
    min: "At least {n} characters",
    max: "Max {n} characters",
    range: "Must be between {min} and {max}",
    rangeMin: "Min {n}",
    rangeMax: "Max {n}",
    match: "Fields do not match",
    compareto: "Not match",
    tagMax: "Max {n} items",
    callback: "Invalid",
    validname: "Letters, numbers, - and _ only",
    en: "English letters only",
    domain: "Invalid domain format",
    mobile: "Invalid mobile number",
    int: "Integer required",
    email: "Invalid email address",
    date: "Invalid date format (YYYY-MM-DD)",
    url: "Invalid URL",
    pwd: "At least 8 chars, 1 letter, 1 number",
    strongpwd: "At least 8 chars, 1 letter, 1 number, 1 special char"
  },
  charCount: "Remaining: {count}"
}, et = {
  box: {
    ok: "确认",
    cancel: "取消",
    loading: "加载中",
    close: "关闭"
  },
  valid: {
    required: "必填",
    pattern: "格式无效",
    min: "至少 {n} 个字符",
    max: "最多 {n} 个字符",
    range: "必须在 {min} 到 {max} 之间",
    rangeMin: "最小 {n}",
    rangeMax: "最大 {n}",
    match: "输入不一致",
    compareto: "不一致",
    tagMax: "最多 {n} 项",
    callback: "校验失败",
    validname: "仅限字母、数字、- 和 _",
    en: "仅限英文字母",
    domain: "无效的域名格式",
    mobile: "无效的手机号",
    int: "需要整数",
    email: "无效的邮箱地址",
    date: "无效的日期格式 (YYYY-MM-DD)",
    url: "无效的 URL",
    pwd: "至少 8 位，包含 1 个字母和 1 个数字",
    strongpwd: "至少 8 位，包含字母、数字和特殊字符"
  },
  charCount: "剩余: {count}"
}, nt = {
  box: {
    ok: "確認",
    cancel: "取消",
    loading: "加載中",
    close: "關閉"
  },
  valid: {
    required: "必填",
    pattern: "格式無效",
    min: "至少 {n} 個字符",
    max: "最多 {n} 個字符",
    range: "必須在 {min} 到 {max} 之間",
    rangeMin: "最小 {n}",
    rangeMax: "最大 {n}",
    match: "輸入不一致",
    compareto: "不一致",
    tagMax: "最多 {n} 項",
    callback: "校驗失敗",
    validname: "僅限字母、數字、- 和 _",
    en: "僅限英文字母",
    domain: "無效的域名格式",
    mobile: "無效的手機號",
    int: "需要整數",
    email: "無效的電子郵件地址",
    date: "無效的日期格式 (YYYY-MM-DD)",
    url: "無效的 URL",
    pwd: "至少 8 位，包含 1 個字母和 1 個數字",
    strongpwd: "至少 8 位，包含字母、數字和特殊字符"
  },
  charCount: "剩餘: {count}"
};
const R = {
  en: tt,
  "zh-CN": et,
  "zh-TW": nt
}, _ = Object.keys(R);
function ot() {
  if (typeof localStorage < "u") {
    const i = localStorage.getItem("uno.lang");
    if (i && _.includes(i)) return i;
  }
  if (typeof navigator < "u") {
    const i = navigator.language, t = i.startsWith("zh") ? i.includes("TW") || i.includes("Hant") ? "zh-TW" : "zh-CN" : "en";
    return _.includes(t) ? t : "en";
  }
  return "en";
}
function it(i, t) {
  const n = t.split(".");
  let e = i;
  for (const o of n) {
    if (!e || typeof e != "object" || !(o in e)) return;
    e = e[o];
  }
  return e;
}
function D(i, t) {
  for (const n in t)
    t[n] && typeof t[n] == "object" ? (i[n] || Object.assign(i, { [n]: {} }), D(i[n], t[n])) : Object.assign(i, { [n]: t[n] });
}
function st(i) {
  i.store("i18n", {
    locale: ot(),
    messages: R,
    loading: !1,
    init() {
      typeof document < "u" && (document.documentElement.lang = this.locale);
    },
    t(t, n = {}) {
      const e = this.locale, o = this.messages[e] || this.messages.en, s = it(o, t);
      return (typeof s == "string" ? s : t).replace(/\{(\w+)\}/g, (c, r) => r in n ? String(n[r]) : `{${r}}`);
    },
    setLocale(t) {
      if (!this.messages[t]) {
        g.warn(`[Uno] Language ${t} not found. Please load it first using loadLocale().`);
        return;
      }
      this.locale = t, typeof localStorage < "u" && localStorage.setItem("uno.lang", t), typeof document < "u" && (document.documentElement.lang = t);
    },
    // Alias for backward compatibility
    setLang(t) {
      this.setLocale(t);
    },
    async loadLocale(t, n) {
      if (!this.messages[t]) {
        this.loading = !0;
        try {
          const e = await fetch(n);
          if (!e.ok) throw new Error(`Failed to load locale: ${e.statusText}`);
          const o = await e.json();
          return this.extend({ [t]: o }), !0;
        } catch (e) {
          throw g.error(`[Uno] Failed to load locale ${t}`, e), e;
        } finally {
          this.loading = !1;
        }
      }
    },
    extend(t) {
      for (const n in t)
        this.messages[n] ? D(this.messages[n], t[n]) : this.messages[n] = t[n];
    }
  }), i.directive("t", (t, { expression: n, modifiers: e }, { effect: o }) => {
    const s = () => i.store("i18n");
    o(() => {
      const a = s();
      a.locale;
      const c = n, r = a.t(c);
      e.length > 0 ? e.forEach((l) => {
        t.setAttribute(l, r);
      }) : t.textContent = r;
    });
  });
}
const rt = () => (h.options?.ui?.layers?.tooltip || k("TOOLTIP"), `
    <div id="uno-tooltip-content"
         x-show="$store.uno.tooltip.visible"
         role="tooltip"
         x-on:mouseenter="$store.uno.tooltip.cancelHide()"
         x-on:mouseleave="$store.uno.tooltip.scheduleHide()"
         x-transition:enter="transition ease-out duration-200"
         x-transition:enter-start="opacity-0 scale-95"
         x-transition:enter-end="opacity-100 scale-100"
         x-transition:leave="transition ease-in duration-150"
         x-transition:leave-start="opacity-100 scale-100"
         x-transition:leave-end="opacity-0 scale-95"
         :class="Uno.getStyle('tooltip', 'content')"
         :style="{ top: $store.uno.tooltip.top + 'px', left: $store.uno.tooltip.left + 'px', zIndex: $store.uno.tooltip.zIndex }">
         
         <!-- Arrow -->
         <div x-show="$store.uno.tooltip.arrow"
              :style="{ top: $store.uno.tooltip.arrowPosition.top + 'px', left: $store.uno.tooltip.arrowPosition.left + 'px' }">
         </div>

         <!-- Content -->
         <div id="uno-tooltip-body" :class="Uno.getStyle('tooltip', 'body')">
            <template x-if="!$store.uno.tooltip.hxUrl">
              <div x-html="$store.uno.tooltip.content"></div>
            </template>
            <template x-if="$store.uno.tooltip.hxUrl">
               <div x-init="htmx.process($el)"
                    :hx-get="$store.uno.tooltip.hxUrl"
                    :hx-trigger="$store.uno.tooltip.hxTrigger"
                    hx-swap="innerHTML"
                    :hx-headers="$store.uno.tooltip.hxHeaders ? JSON.stringify($store.uno.tooltip.hxHeaders) : null">
                    <!-- Loading State -->
                    <div class="uno-tooltip-loading">Loading...</div>
               </div>
            </template>
         </div>
    </div>
  `);
let $ = 0, S = null;
const at = (i) => {
  const t = () => {
    $ || typeof requestAnimationFrame == "function" && ($ = requestAnimationFrame(() => {
      $ = 0;
      const e = i.store("uno")?.tooltip;
      e && e.updatePosition();
    }));
  }, n = () => {
    if (S || typeof window > "u" || typeof document > "u") return;
    const e = () => t(), o = () => t(), s = (c) => {
      if (c.key !== "Escape") return;
      const r = i.store("uno")?.tooltip;
      r?.visible && r.hide();
    }, a = (c) => {
      const r = i.store("uno")?.tooltip;
      if (!r?.visible || r.triggerMode !== "click" || r.clickOutside !== !0) return;
      const l = c.target, u = document.getElementById("uno-tooltip-content"), p = r.triggerEl && r.triggerEl.contains(l), f = u && u.contains(l);
      !p && !f && r.hide();
    };
    window.addEventListener("scroll", e, !0), window.addEventListener("resize", o), window.addEventListener("keydown", s), window.addEventListener("click", a), S = () => {
      window.removeEventListener("scroll", e, !0), window.removeEventListener("resize", o), window.removeEventListener("keydown", s), window.removeEventListener("click", a), S = null;
    };
  };
  return {
    visible: !1,
    content: "",
    placement: "top",
    top: -9999,
    left: -9999,
    zIndex: k("TOOLTIP"),
    triggerEl: null,
    triggerMode: "hover",
    clickOutside: !0,
    arrow: !0,
    arrowPosition: { top: 0, left: 0 },
    // HTMX support
    hxUrl: null,
    hxHeaders: null,
    hxTrigger: "load",
    init() {
      this.ensureContainer();
    },
    show(e = {}) {
      g.info("Tooltip show called", e);
      const o = h.getConfig("tooltip", e);
      if (g.info("Tooltip resolved config", o), this.triggerEl = o.element || this.triggerEl, o.target) {
        const s = document.querySelector(o.target);
        s ? this.content = P(s.innerHTML) : (g.warn(`Tooltip target element not found: ${o.target}`), this.content = o.content || "");
      } else
        this.content = o.content || "";
      this.placement = o.placement, this.triggerMode = o.trigger, this.clickOutside = o.clickOutside, this.offset = o.offset, this.arrow = o.arrow !== !1, this.zIndex = h.options?.ui?.layers?.tooltip || k("TOOLTIP"), this.hxUrl = o.hxUrl, this.hxHeaders = o.hxHeaders, this.hxTrigger = o.hxTrigger || "load", this.visible = !0, this.ensureContainer(), n(), i.nextTick(() => {
        this.updatePosition();
      });
    },
    hide() {
      this.visible = !1, this.triggerEl = null, S && (S(), S = null), $ && (cancelAnimationFrame($), $ = 0);
    },
    cancelHide() {
      this.hideTimeout && (clearTimeout(this.hideTimeout), this.hideTimeout = null);
    },
    scheduleHide(e = 300) {
      this.cancelHide(), this.hideTimeout = setTimeout(() => {
        this.hide();
      }, e);
    },
    updatePosition() {
      if (!this.visible || !this.triggerEl) return;
      const e = document.getElementById("uno-tooltip-content");
      if (!e) return;
      const o = X(this.triggerEl, e, this.placement, {
        offset: this.offset,
        arrow: this.arrow
      });
      this.top = o.top, this.left = o.left, this.placement = o.placement, o.arrow && (this.arrowPosition = o.arrow);
    },
    ensureContainer() {
      const e = "uno-tooltip-manager";
      if (document.getElementById(e)) return;
      const o = h.options?.ui?.layers?.tooltip || k("TOOLTIP"), s = rt();
      O.ensureContainer(e, s, {
        "x-data": "{}",
        style: `position: absolute; top: 0; left: 0; width: 0; height: 0; z-index: ${o};`
      });
      const a = document.getElementById(e);
      a && i && !a._x_dataStack && i.initTree(a);
    }
  };
}, lt = (i) => {
  const t = (n, { expression: e, modifiers: o }, { evaluateLater: s, effect: a, cleanup: c }) => {
    const r = e ? s(e) : null, l = (d) => (o.length > 0 && console.warn(
      "[Uno] x-uno-tooltip modifiers are deprecated in v2.0. Use configuration object instead."
    ), { ...typeof d == "object" && d !== null ? d : d ? { content: d } : {}, element: n }), u = (d) => {
      const m = i.store("uno")?.tooltip;
      m && m.show(l(d));
    }, p = (d) => {
      const m = i.store("uno")?.tooltip;
      m && (m.visible && m.triggerEl === n ? m.hide() : u(d));
    };
    let f = null;
    r && a(() => {
      r((d) => {
        f = d;
      });
    });
    const y = () => {
      if (e && f === null) return;
      const d = l(f);
      h.getConfig("tooltip", d).trigger === "hover" && u(f);
    }, b = () => {
      const d = i.store("uno")?.tooltip;
      d && d.triggerMode === "hover" && d.triggerEl === n && d.scheduleHide();
    }, v = (d) => {
      const m = l(f);
      h.getConfig("tooltip", m).trigger === "click" && (d.preventDefault(), p(f));
    }, x = () => {
      const d = l(f);
      h.getConfig("tooltip", d).trigger === "focus" && u(f);
    }, w = () => {
      const d = i.store("uno")?.tooltip;
      d && d.triggerMode === "focus" && d.triggerEl === n && d.hide();
    };
    n.addEventListener("mouseenter", y), n.addEventListener("mouseleave", b), n.addEventListener("click", v), n.addEventListener("focus", x), n.addEventListener("blur", w), c(() => {
      n.removeEventListener("mouseenter", y), n.removeEventListener("mouseleave", b), n.removeEventListener("click", v), n.removeEventListener("focus", x), n.removeEventListener("blur", w);
      const d = i.store("uno")?.tooltip;
      d && d.triggerEl === n && d.hide();
    });
  };
  i.directive("uno-tooltip", t);
};
const ct = (i, t = {}) => {
  const n = {
    placement: "top",
    trigger: "hover",
    delay: 300,
    offset: 8,
    arrow: !0,
    clickOutside: !0
  };
  h.registerPluginDefaults("tooltip", {
    ...n,
    ...t.plugins?.tooltip,
    ...t.tooltip
  });
  const e = at(i);
  U(i, "tooltip", e), lt(i);
}, dt = () => `
  <template x-for="(items, position) in $store.uno.notify.groups" :key="position">
    <div :class="$store.uno.notify.getContainerClass(position)"
         :style="$store.uno.notify.getContainerStyle(position)">
      
      <template x-for="item in items" :key="item.id">
        <div :class="[
            Uno.getStyle('notify', 'item'),
            Uno.getStyle('notify', 'types.' + item.type)
          ]"
          role="alert"
          aria-live="assertive"
          x-show="item.isVisible"
          x-transition:enter="transition ease-out duration-300"
          x-transition:enter-start="opacity-0 transform scale-90"
          x-transition:enter-end="opacity-100 transform scale-100"
          x-transition:leave="transition ease-in duration-200"
          x-transition:leave-start="opacity-100 transform scale-100"
          x-transition:leave-end="opacity-0 transform scale-90"
          :data-id="item.id"
          data-testid="notify-item"
         >
          
          <div :class="Uno.getStyle('notify', 'content')">
            <!-- Standard Content -->
            <template x-if="!item.hxUrl">
              <div :class="Uno.getStyle('notify', 'row')">
                 <!-- Icon -->
                 <div x-show="item.icon" :class="Uno.getStyle('notify', 'iconWrap')">
                    <i :class="[item.icon, Uno.getStyle('notify', 'icons.' + item.type)]"></i>
                 </div>
                 
                 <!-- Content -->
                 <div :class="Uno.getStyle('notify', 'textWrap')">
                    <p :class="Uno.getStyle('notify', 'text')" x-text="item.message"></p>
                 </div>
                 
                 <!-- Close Button -->
                 <div x-show="item.closable" :class="Uno.getStyle('notify', 'closeWrap')">
                    <button type="button" 
                            @click="$store.uno.notify.remove(item.id)"
                            :class="Uno.getStyle('notify', 'closeBtn')">
                      <span class="sr-only">Close</span>
                      <svg :class="Uno.getStyle('notify', 'closeIcon')" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                      </svg>
                    </button>
                 </div>
              </div>
            </template>

            <!-- HTMX Content -->
            <template x-if="item.hxUrl">
               <div class="w-full"
                    x-init="
                      $el.setAttribute('hx-get', item.hxUrl);
                      $el.setAttribute('hx-trigger', item.hxTrigger || 'load');
                      $el.setAttribute('hx-swap', item.hxSwap || 'innerHTML');
                      if (item.hxHeaders) $el.setAttribute('hx-headers', JSON.stringify(item.hxHeaders));
                      htmx.process($el);
                    "
                    @htmx:beforeRequest.stop="$store.uno.notify._onHtmxBegin(item)"
                    @htmx:afterRequest.stop="$store.uno.notify._onHtmxEnd(item)"
                    @htmx:afterSwap.stop="$store.uno.notify.onContentLoaded($el, item)"
                    @htmx:responseError.stop="$store.uno.notify.onHtmxError($event.detail, item)"
                    @htmx:sendError.stop="$store.uno.notify.onHtmxError($event.detail, item)"
                    @htmx:timeout.stop="$store.uno.notify.onHtmxError($event.detail, item)"
               >
                  <div x-html="item.content || ''"></div>
               </div>
            </template>
          </div>

        </div>
      </template>
    </div>
  </template>
`, ut = () => `
    <div :class="Uno.getStyle('notify', 'progressWrap')"
         style="z-index: ${k("MAX")}"
         x-show="$store.uno.notify.progressVisible"
         x-transition:enter="transition-opacity duration-200"
         x-transition:enter-start="opacity-0"
         x-transition:enter-end="opacity-100"
         x-transition:leave="transition-opacity duration-200"
         x-transition:leave-start="opacity-100"
         x-transition:leave-end="opacity-0">
        <div :class="Uno.getStyle('notify', 'progressBar')"
             :style="{ width: $store.uno.notify.progressPercent + '%' }">
        </div>
    </div>
  `, ft = (i) => ({
  groups: {
    "top-left": [],
    "top-center": [],
    "top-right": [],
    "center-left": [],
    center: [],
    "center-right": [],
    "bottom-left": [],
    "bottom-center": [],
    "bottom-right": []
  },
  nextId: 1,
  progressRequests: 0,
  progressVisible: !1,
  progressPercent: 0,
  _progressInterval: null,
  init() {
    this.ensureContainer();
  },
  // HTMX Event Handlers
  onContentLoaded(t, n) {
    this._onHtmxEnd(n);
    const e = t.innerHTML || t.textContent || "", o = this._findNotificationById(n.id);
    o && (o.content = e), n.onContentLoad && n.onContentLoad(t), window.Alpine && window.Alpine.initTree && window.Alpine.initTree(t);
  },
  // Helper method to find notification by ID
  _findNotificationById(t) {
    for (const n of Object.keys(this.groups)) {
      const o = this.groups[n].find((s) => s.id === t);
      if (o) return o;
    }
    return null;
  },
  onHtmxError(t, n) {
    this._onHtmxEnd(n), console.error("HTMX Error:", t), this.error("Content load failed"), n.onError && n.onError(t);
  },
  _onHtmxBegin(t) {
    !t || t.__unoHtmxActive || (t.__unoHtmxActive = !0, this.progressRequests++, this.progressRequests === 1 && this.startProgress());
  },
  _onHtmxEnd(t) {
    !t || !t.__unoHtmxActive || (t.__unoHtmxActive = !1, this.progressRequests > 0 && this.progressRequests--, this.progressRequests === 0 && this.stopProgress());
  },
  startProgress() {
    this._progressInterval && clearInterval(this._progressInterval), this.progressVisible = !0, this.progressPercent = 0, this._progressInterval = setInterval(() => {
      this.progressPercent >= 100 ? this.progressPercent = 0 : this.progressPercent += 10;
    }, 100);
  },
  stopProgress() {
    this._progressInterval && clearInterval(this._progressInterval), this.progressPercent = 100, setTimeout(() => {
      this.progressVisible = !1, setTimeout(() => {
        this.progressPercent = 0;
      }, 200);
    }, 300);
  },
  getContainerStyle() {
    return {
      zIndex: h.options?.ui?.layers?.notify || k("NOTIFY")
    };
  },
  getContainerClass(t) {
    const n = h.getStyle("notify", "container"), e = F(t), o = W(t);
    return [n, e, o.flexDir, o.alignItems].join(" ");
  },
  ensureContainer() {
    const t = "uno-notify-container", n = dt(), e = O.ensureContainer(t, n, {
      "x-data": "{}"
      // Independent Alpine component
    });
    e && i.initTree && i.initTree(e);
    const o = "uno-global-progress", s = ut(), a = O.ensureContainer(o, s, {
      "x-data": "{}"
    });
    a && i.initTree && i.initTree(a);
  },
  show(t, n = {}) {
    this.ensureContainer();
    const e = h.getConfig("notify", n), o = this.nextId++;
    let s = t;
    i.store("i18n") && typeof i.store("i18n").t == "function" && (s = i.store("i18n").t(t));
    const a = e.duration !== void 0 ? e.duration : 3e3;
    g.info("[Notify] Show:", { id: o, message: s, duration: a, config: e });
    const c = {
      id: o,
      message: s,
      isVisible: !0,
      // For x-show
      closable: !0,
      duration: a,
      content: e.hxUrl ? "" : void 0,
      // Initialize content for HTMX-loaded content
      ...e
    }, r = e.position || "top-right";
    return this.groups[r] ? this.groups[r].push(c) : (this.groups["top-right"] || (this.groups["top-right"] = []), this.groups["top-right"].push(c)), a > 0 && setTimeout(() => {
      this.remove(o);
    }, a), o;
  },
  remove(t) {
    Object.keys(this.groups).forEach((n) => {
      const e = this.groups[n].findIndex((o) => o.id === t);
      e !== -1 && (this.groups[n][e].visible = !1, setTimeout(() => {
        const o = this.groups[n].findIndex((s) => s.id === t);
        o !== -1 && this.groups[n].splice(o, 1);
      }, 300));
    });
  },
  // 3. Proxy Methods
  success(t, n = {}) {
    return this.show(t, { type: "success", ...n });
  },
  error(t, n = {}) {
    return this.show(t, { type: "error", ...n });
  },
  warning(t, n = {}) {
    return this.show(t, { type: "warning", ...n });
  },
  info(t, n = {}) {
    return this.show(t, { type: "info", ...n });
  }
}), gt = (i) => {
  i.directive("uno-notify", (t, { expression: n }, { evaluateLater: e }) => {
    if (!n) {
      console.warn(
        "[Uno] x-uno-notify directive requires a configuration object or message string."
      );
      return;
    }
    const o = e(n), s = () => {
      const a = i.store("uno")?.notify;
      a && o((c) => {
        let r = "", l = {};
        if (typeof c == "string")
          r = c;
        else if (typeof c == "object" && c !== null)
          r = c.message || "", l = c;
        else {
          console.warn("[Uno] x-uno-notify directive value must be a string or object", c);
          return;
        }
        a.show(r, l);
      });
    };
    t.addEventListener("click", s);
  });
};
const ht = (i, t = {}) => {
  const n = {
    duration: 3e3,
    position: "top-right"
  };
  h.registerPluginDefaults("notify", {
    ...n,
    ...t.plugins?.notify,
    ...t.notify
  });
  const e = ft(i);
  U(i, "notify", e), e.init(), gt(i);
}, H = {
  // --- Standard Types ---
  email: /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/,
  url: /^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([/\w .-]*)*\/?$/i,
  date: /^(\d{4}-\d{2}-\d{2}|\d{2}\/\d{2}\/\d{4}|\d{2}\.\d{2}\.\d{4})$/,
  // --- Numeric ---
  int: /^-?\d+$/,
  number: /^-?\d*(\.\d+)?$/,
  // Added floating point support for general 'number'
  // --- String Content ---
  validname: /^[a-zA-Z0-9_-]+$/,
  en: /^[a-zA-Z]+$/,
  domain: new RegExp("^((?!-)[A-Za-z0-9-]{1,63}(?<!-)\\.)+[A-Za-z]{2,6}$"),
  mobile: /^\+?[\d\s-]{7,}$/,
  // --- Security ---
  pwd: /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$/,
  // Letters + Numbers, 8+
  strongpwd: /^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$/
  // + Special Char
}, mt = (i) => ({
  /**
   * Validation error messages keyed by field name.
   * @type {Object.<string, string>}
   */
  errors: {},
  /**
   * Validation rules keyed by field name.
   * @type {Object.<string, Object>}
   */
  rules: {},
  /**
   * Custom validators registry.
   * @type {Object.<string, Function>}
   */
  customValidators: {},
  /**
   * Adds or updates validation rules for a field.
   *
   * @param {string} fieldName - The unique identifier for the field (name or id)
   * @param {Object} [instanceOpts={}] - Instance-specific configuration
   */
  addRules(t, n = {}) {
    this.rules[t] = h.getConfig("valid", n);
  },
  /**
   * Registers a custom validator function.
   *
   * @param {string} name - The name of the validator
   * @param {Function} validatorFn - The validator function.
   */
  registerValidator(t, n) {
    typeof n == "function" ? this.customValidators[t] = n : g.error(`[Valid] Invalid validator function for ${t}`);
  },
  /**
   * Validates a single field against its rules.
   *
   * @param {string} fieldName - The field identifier
   * @param {string|number} value - The value to validate
   * @returns {boolean} True if valid, false otherwise
   */
  validate(t, n) {
    const e = this.rules[t];
    if (!e) return !0;
    let o = !0, s = "";
    const a = n == null ? "" : String(n), c = a.trim() === "";
    if (e.required && c && (o = !1, s = e.msgRequired || this.getI18nMsg("valid.required") || "This field is required"), o && c && !e.required)
      return delete this.errors[t], !0;
    if (o) {
      const r = e.type || (e.email ? "email" : null);
      r && H[r] && (H[r].test(a) || (o = !1, s = e[`msg${r.charAt(0).toUpperCase() + r.slice(1)}`] || e.msgType || this.getI18nMsg(`valid.${r}`) || `Invalid ${r} format`));
    }
    if (o && e.minLength !== null && e.minLength !== void 0 && a.length < e.minLength && (o = !1, s = e.msgMinLength || this.getI18nMsg("valid.min", { n: e.minLength }) || `Minimum length is ${e.minLength}`), o && e.maxLength !== null && e.maxLength !== void 0 && a.length > e.maxLength && (o = !1, s = e.msgMaxLength || this.getI18nMsg("valid.max", { n: e.maxLength }) || `Maximum length is ${e.maxLength}`), o && (e.min !== null && e.min !== void 0 || e.max !== null && e.max !== void 0))
      if (H.number.test(a)) {
        const r = parseFloat(a);
        e.min !== null && e.min !== void 0 && r < e.min && (o = !1, s = e.msgMin || this.getI18nMsg("valid.rangeMin", { n: e.min }) || `Value must be at least ${e.min}`), o && e.max !== null && e.max !== void 0 && r > e.max && (o = !1, s = e.msgMax || this.getI18nMsg("valid.rangeMax", { n: e.max }) || `Value must be at most ${e.max}`);
      } else
        o = !1, s = e.msgNumber || this.getI18nMsg("valid.number") || "Must be a number";
    if (o && e.pattern)
      try {
        (e.pattern instanceof RegExp ? e.pattern : new RegExp(e.pattern)).test(a) || (o = !1, s = e.msgPattern || this.getI18nMsg("valid.pattern") || "Invalid format");
      } catch (r) {
        g.error(`[Valid] Invalid regex pattern for ${t}`, r);
      }
    if (o && e.match) {
      const r = document.querySelector(
        e.match.startsWith("#") ? e.match : `[name="${e.match}"]`
      ) || document.getElementById(e.match);
      r && a !== r.value && (o = !1, s = e.msgMatch || this.getI18nMsg("valid.match") || "Fields do not match");
    }
    if (o && e.validator) {
      let r = e.validator;
      if (typeof e.validator == "string" && (r = this.customValidators[e.validator], r || g.warn(`[Valid] Validator "${e.validator}" not found`)), typeof r == "function") {
        const l = r(n, t);
        l === !1 ? (o = !1, s = e.msgValidator || "Validation failed") : typeof l == "string" ? (o = !1, s = l) : typeof l == "object" && l.valid === !1 && (o = !1, s = l.message || e.msgValidator || "Validation failed");
      }
    }
    return o ? delete this.errors[t] : this.errors[t] = s, o;
  },
  /**
   * Validates all fields in a form.
   *
   * @param {HTMLFormElement} form - The form element to validate
   * @returns {boolean} True if the form is valid
   */
  validateForm(t) {
    const n = t.querySelectorAll("input, select, textarea");
    let e = !0, o = null;
    return n.forEach((s) => {
      if (s.disabled) return;
      const a = s.name || s.id;
      this.rules[a] && (this.validate(a, s.value) || (e = !1, o || (o = s)));
    }), !e && o && (h.getConfig("valid") || {}).scrollOnError !== !1 && (o.scrollIntoView({ behavior: "smooth", block: "center" }), o.focus({ preventScroll: !0 })), e;
  },
  showError(t, n) {
    this.errors[t] = n;
  },
  hideError(t) {
    delete this.errors[t];
  },
  clearErrors() {
    this.errors = {};
  },
  showErrors(t) {
    Object.assign(this.errors, t);
  },
  getI18nMsg(t, n) {
    const e = i.store("i18n");
    return e && typeof e.t == "function" ? e.t(t, n) : null;
  }
}), pt = (i) => {
  i.directive("uno-valid", (t, { expression: n }, { evaluateLater: e, cleanup: o }) => {
    if (t.tagName === "FORM") return;
    const s = t.name || t.id;
    if (!s) {
      console.warn("[Uno Valid] Element has no name or id:", t);
      return;
    }
    const a = n ? e(n) : () => {
    }, r = {
      // Standard HTML5
      required: t.hasAttribute("required"),
      minLength: t.getAttribute("minlength") || null,
      maxLength: t.getAttribute("maxlength") || null,
      min: t.getAttribute("min") || null,
      max: t.getAttribute("max") || null,
      pattern: t.getAttribute("pattern") || null,
      type: t.getAttribute("type") || null
    }, l = (f) => {
      i.store("uno").valid.addRules(s, { ...r, ...f });
    };
    n ? a((f) => {
      l(typeof f == "object" ? f : {});
    }) : l({});
    const u = () => {
      const f = i.store("uno").valid.rules[s];
      f && f.realtime !== !1 && i.store("uno").valid.validate(s, t.value);
    }, p = () => {
      i.store("uno").valid.validate(s, t.value);
    };
    t.addEventListener("input", u), t.addEventListener("blur", p), o(() => {
      t.removeEventListener("input", u), t.removeEventListener("blur", p);
    });
  }), typeof document < "u" && (document.body.addEventListener("htmx:confirm", (t) => {
    const n = t.detail.elt, e = n.tagName === "FORM" ? n : n.closest("form");
    e && e.hasAttribute("x-uno-valid") && (i.store("uno").valid.validateForm(e) || (t.preventDefault(), t.stopPropagation(), g.info("[Valid] Form validation failed. Request blocked."), e.dispatchEvent(
      new CustomEvent("uno:validation-failed", {
        bubbles: !0,
        detail: { errors: i.store("uno").valid.errors }
      })
    )));
  }), document.body.addEventListener("htmx:afterRequest", (t) => {
    const n = t.detail.xhr;
    if (n.status === 422)
      try {
        const e = JSON.parse(n.responseText);
        e.errors && (i.store("uno").valid.showErrors(e.errors), g.info("[Valid] Server validation errors applied", e.errors));
      } catch {
      }
  }));
};
const xt = (i, t = {}) => {
  const n = {
    realtime: !0,
    showError: !0,
    scrollOnError: !0
  };
  h.registerPluginDefaults("valid", {
    ...n,
    ...t.plugins?.valid,
    ...t.valid
  });
  const e = mt(i);
  U(i, "valid", e), pt(i);
}, yt = (i) => ({
  // State: { [group]: activeTabName }
  groups: {},
  // Registry: { [tabName]: { group, name, ... } }
  registry: {},
  init() {
    const t = () => {
      let n = window.location.hash.slice(1);
      if (!n && window.location.href.includes("#") && (n = window.location.href.split("#")[1]), g.info(`[Tab] checkHash: ${n}`), n && this.registry[n]) {
        const { group: e, name: o } = this.registry[n];
        g.info(`[Tab] checkHash match: ${o}`), this.activate({ group: e, name: o, updateHash: !1 });
      }
    };
    t(), window.addEventListener("hashchange", t), window.addEventListener("load", t);
  },
  /**
   * Activate a tab
   * @param {object} instanceOpts - Options { name, group, ... }
   */
  activate(t = {}) {
    const n = h.getConfig("tab", {
      group: "default",
      hash: !0,
      ...t
    });
    n.name && (g.info(
      `[Tab] Activating ${n.name} in ${n.group}. Hash update: ${n.updateHash !== !1}`
    ), this.registry[n.name] || (this.registry[n.name] = { group: n.group, name: n.name }), this.groups[n.group] = n.name, n.hash && n.updateHash !== !1 && window.location.hash.slice(1) !== n.name && window.history.replaceState(null, null, `#${n.name}`), n.updateHash === !1 && setTimeout(() => {
      const e = document.getElementById(`tab-${n.name}`);
      e && e.hasAttribute("hx-get") && (g.info(`[Tab] Triggering HTMX for ${n.name}`), window.htmx && window.htmx.trigger(e, "click"));
    }, 0));
  },
  /**
   * Check if a tab is active
   * @param {string} name - Tab name
   * @param {string} group - Group name
   * @returns {boolean}
   */
  isActive(t, n = "default") {
    return this.groups[n] === t;
  },
  /**
   * Register a tab (helper for directives or manual setup)
   * @param {string} name
   * @param {string} group
   */
  register(t, n = "default") {
    g.info(`[Tab] Registering ${t} in ${n}. Hash: ${window.location.hash}`), this.registry[t] = { name: t, group: n }, this.groups[n] || (this.groups[n] = t);
    let e = window.location.hash.slice(1);
    !e && window.location.href.includes("#") && (e = window.location.href.split("#")[1]), e === t && (g.info(`[Tab] Hash match for ${t}. Activating.`), this.activate({ group: n, name: t, updateHash: !1 }));
  }
}), wt = (i) => {
  i.directive("uno-tab", (t, { expression: n }, { evaluateLater: e, cleanup: o }) => {
    if (!n) {
      console.warn("[Uno] x-uno-tab directive requires a configuration object.");
      return;
    }
    const s = e(n);
    s((c) => {
      let r = {};
      if (typeof c == "object" && c !== null)
        r = c;
      else if (typeof c == "string")
        r = { name: c };
      else {
        console.warn("[Uno] x-uno-tab directive value must be a string or object", c);
        return;
      }
      r.group || (r.group = "default"), r.name && i.store("uno").tab.register(r.name, r.group);
    });
    const a = (c) => {
      c.preventDefault(), s((r) => {
        let l = {};
        typeof r == "object" && r !== null ? l = r : typeof r == "string" && (l = { name: r }), l.name && i.store("uno").tab.activate(l);
      });
    };
    t.addEventListener("click", a), o(() => t.removeEventListener("click", a));
  });
};
const vt = (i, t = {}) => {
  const n = {
    defaultTab: "",
    hash: !0,
    activeClass: "is-active",
    inactiveClass: ""
  };
  h.registerPluginDefaults("tab", { ...n, ...t.plugins?.tab, ...t.tab });
  const e = yt();
  U(i, "tab", e), i.store("uno").tab.init(), wt(i);
};
h.dom = O;
typeof window < "u" && (window.Uno = h);
h.use(st).use(ct).use(ht).use(J).use(xt).use(vt);
if (typeof window < "u" && typeof document < "u") {
  const i = () => {
    if (h.initialized) return;
    const t = window.UnoConfig || {};
    h.init(window.Alpine, t);
  };
  window.Alpine ? i() : document.addEventListener("alpine:init", i);
}
export {
  O as DomManager,
  h as Uno,
  J as box,
  h as default,
  ht as notify,
  vt as tab,
  ct as tooltip,
  xt as valid
};
