import q from "alpinejs";
let U = !1, B = null;
const m = {
  /**
   * Set debug mode
   * @param {boolean} debug
   */
  setDebug(i) {
    U = i;
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
    U && typeof console < "u" && typeof console.log == "function" && console.log("[Uno Info]", ...i);
  },
  /**
   * Log warning message
   * @param {...any} args - Arguments to log
   */
  warn(...i) {
    U && typeof console < "u" && typeof console.warn == "function" && console.warn("[Uno Warn]", ...i);
  },
  /**
   * Log error message
   * @param {...any} args - Arguments to log
   */
  error(...i) {
    if (B)
      try {
        B(...i);
      } catch (s) {
        U && typeof console < "u" && typeof console.error == "function" && console.error("[Uno Error] Reporter failed:", s);
      }
    U && typeof console < "u" && typeof console.error == "function" && console.error("[Uno Error]", ...i);
  }
};
function Z(i) {
  return String(i ?? "").replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;").replace(/'/g, "&#39;");
}
const A = /* @__PURE__ */ new WeakMap();
function P(i, s) {
  if (!i) return !1;
  const a = String(s);
  let l = A.get(i);
  return l || (l = /* @__PURE__ */ new Set(), A.set(i, l)), l.has(a) ? !1 : (l.add(a), !0);
}
function z(i) {
  return i ? i.replace(/<script\b[^>]*>([\s\S]*?)<\/script>/gim, "").replace(/<img\b[^>]*>/gim, "").replace(/\s+on[a-z]+="[^"]*"/gim, "").replace(/\s+on[a-z]+='[^']*'/gim, "").replace(/\s+href="javascript:[^"]*"/gim, "").replace(/\s+href='javascript:[^']*'/gim, "").replace(/javascript:/gim, "") : "";
}
function X(i, s, { trusted: a = !1, alpine: l } = {}) {
  i && (i.innerHTML = a ? String(s ?? "") : z(String(s ?? "")), l?.initTree?.(i));
}
const H = {
  /**
   * Ensure a container element exists in the DOM.
   * @param {string} id - Element ID
   * @param {string} html - Inner HTML content
   * @param {Object} attributes - Attributes to set
   * @param {HTMLElement} parent - Parent element (default: document.body)
   */
  ensureContainer(i, s = "", a = {}, l = document.body) {
    const n = document.getElementById(i);
    if (n) return n;
    try {
      const t = document.createElement("div");
      return t.id = i, window.Uno?.options?.csp?.nonce && !a.nonce && (a.nonce = window.Uno.options.csp.nonce), Object.entries(a).forEach(([e, r]) => {
        t.setAttribute(e, r);
      }), s && (t.innerHTML = s), l ? l.appendChild(t) : document.documentElement.appendChild(t), t;
    } catch (t) {
      return m.error(`DomManager: Failed to inject container #${i}`, t), null;
    }
  },
  /**
   * Inject CSS styles with CSP nonce support.
   * @param {string} id - Style element ID
   * @param {string} css - CSS content
   */
  injectStyles(i, s) {
    if (!document.getElementById(i))
      try {
        const a = document.createElement("style");
        a.id = i, window.Uno?.options?.csp?.nonce && a.setAttribute("nonce", window.Uno.options.csp.nonce), a.textContent = s, document.head.appendChild(a);
      } catch (a) {
        m.error(`DomManager: Failed to inject styles #${i}`, a);
      }
  },
  // Exports
  setHTML: X,
  // Secure HTML setting
  escape: Z,
  // HTML escaping
  sanitize: z
  // HTML sanitization
};
function S(i) {
  return Object.prototype.toString.call(i) === "[object Object]";
}
function T(i, s) {
  if (!S(i)) return S(s) ? { ...s } : s;
  if (!S(s)) return s === void 0 ? { ...i } : s;
  const a = { ...i };
  for (const [l, n] of Object.entries(s)) {
    if (S(n) && S(a[l])) {
      a[l] = T(a[l], n);
      continue;
    }
    if (S(n)) {
      a[l] = T({}, n);
      continue;
    }
    if (typeof n == "string" && n.trim().indexOf("+") === 0 && typeof a[l] == "string") {
      a[l] = `${a[l]} ${n.trim().slice(1)}`.trim();
      continue;
    }
    a[l] = n;
  }
  return a;
}
const $ = {
  BASE: 0,
  DROPDOWN: 1e3,
  STICKY: 1100,
  FIXED: 1200,
  MODAL_BACKDROP: 1300,
  MODAL: 1400,
  POPOVER: 1500,
  TOAST: 1600,
  TOOLTIP: 1700,
  MAX: 9999
}, M = (i) => {
  if (!i) return $.BASE;
  const s = i.toUpperCase();
  return $[s] !== void 0 ? $[s] : $.BASE;
}, Y = [
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
function K(i, s = {}) {
  const a = (c) => {
    if (!(c instanceof HTMLElement)) return !1;
    const o = window.getComputedStyle(c);
    return o.display === "none" || o.visibility === "hidden" ? !1 : c.getClientRects().length > 0;
  }, l = () => Array.from(i.querySelectorAll(Y)).filter((c) => a(c)), n = () => l()[0], t = () => {
    const c = l();
    return c[c.length - 1];
  }, e = (c) => {
    if (c.key === "Tab") {
      if (c.shiftKey) {
        if (document.activeElement === n()) {
          c.preventDefault();
          const o = t();
          o && o.focus();
        }
      } else if (document.activeElement === t()) {
        c.preventDefault();
        const o = n();
        o && o.focus();
      }
    }
  }, r = () => {
    if (i.contains(document.activeElement)) return;
    const c = n();
    c && c.focus();
  };
  return typeof requestAnimationFrame == "function" ? requestAnimationFrame(() => requestAnimationFrame(r)) : setTimeout(r, 0), i.addEventListener("keydown", e), () => {
    i.removeEventListener("keydown", e);
  };
}
function I(i, s, a) {
  if (!i || typeof i.store != "function") {
    console.error("[Uno] Alpine.store is not available");
    return;
  }
  const l = i.store("uno");
  l ? i.store("uno", {
    ...l,
    [s]: a
  }) : i.store("uno", { [s]: a }), i.store(s) || i.store(s, i.store("uno")?.[s]);
}
const N = (i) => {
  const s = i || "center", a = {
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
  return a[s] || a.center;
}, J = (i) => {
  const s = i || "top-right", a = {
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
  return a[s] || a["top-right"];
}, G = (i) => {
  const s = i || "top-right";
  let a = "flex-col", l = "items-end", n = "space-y-2";
  return s.includes("bottom") && (a = "flex-col-reverse", n = "space-y-reverse space-y-2"), s.includes("left") ? l = "items-start" : s.includes("center") ? s === "center" ? l = "items-center justify-center" : l = "items-center" : l = "items-end", { flexDir: a, alignItems: l, spaceY: n };
};
function Q(i, s, a) {
  const l = i.getBoundingClientRect(), n = s.getBoundingClientRect(), t = n.width, e = n.height, r = l.width, c = l.height, o = { top: l.top, left: l.left }, d = window.innerWidth, u = window.innerHeight, f = 0, p = 0, x = u, y = d, v = (C, O) => {
    if (O) {
      if (C + t > y)
        return y - t - 10;
      if (C < p)
        return p + 10;
    } else {
      if (C + e > x)
        return x - e - 10;
      if (C < f)
        return f + 10;
    }
    return C;
  };
  let b, w, E = a;
  switch (a) {
    case "left":
      w = o.left - t - 5, w < p && (w = o.left + r + 5, E = "right"), b = v(o.top + c / 2 - e / 2, !1);
      break;
    case "top":
      b = o.top - e - 8, b < f && (b = o.top + c + 8, E = "bottom"), w = o.left + (r - t) / 2, w = v(w, !0);
      break;
    case "right":
      w = o.left + r + 5, w + t > y && (w = o.left - t - 5, E = "left"), b = v(o.top + c / 2 - e / 2, !1);
      break;
    case "bottom":
    default:
      b = o.top + c + 8, b + e > x && (b = o.top - e - 8, E = "top"), w = o.left + (r - t) / 2, w = v(w, !0);
      break;
  }
  return { top: b, left: w, position: E };
}
const tt = (i, s = {}) => {
  const a = {
    defaultWidth: 400,
    defaultPosition: "center",
    okText: "OK",
    cancelText: "Cancel",
    closeOnEsc: !1,
    closeOnOverlayClick: !1
  };
  h.registerPluginDefaults("box", { ...a, ...s.plugins?.box, ...s.box });
  const l = (t) => h.getStyle("box", t);
  i.directive("uno-trap", (t, { expression: e }, { cleanup: r }) => {
    const c = K(t);
    r(c);
  }), i.directive("uno-box", (t, { expression: e }, { evaluateLater: r, cleanup: c }) => {
    if (!e) {
      console.warn("[Uno] x-uno-box directive requires a configuration object.");
      return;
    }
    const o = r(e), d = (u) => {
      u.preventDefault();
      const f = i.store("uno")?.box;
      if (!f) {
        console.warn("[Uno] Box store not found. Is Uno initialized?");
        return;
      }
      o((p) => {
        p && typeof p == "object" ? f.open(p) : console.warn("[Uno] x-uno-box directive value must be an object", p);
      });
    };
    t.addEventListener("click", d), c(() => t.removeEventListener("click", d));
  }), I(i, "box", {
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
    alert(t, e = {}) {
      const r = typeof e == "string" ? { title: e } : e;
      return this.open("dialog", {
        content: t,
        showOkButton: !0,
        ...r
      });
    },
    /**
     * Show a confirm dialog
     * @param {string} message - Confirmation message
     * @param {object} options - Additional options
     * @returns {Promise<boolean>} Resolves true if OK, false otherwise
     */
    confirm(t, e = {}) {
      return new Promise((r) => {
        this.open("confirm", {
          content: t,
          onOk: () => r(!0),
          onCancel: () => r(!1),
          onClose: () => r(!1),
          // Closing without choice = cancel
          ...e
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
    prompt(t, e = "", r = {}) {
      return new Promise((c) => {
        const o = "prompt-" + Math.random().toString(36).slice(2, 11), d = l("promptWrapper"), u = l("promptText"), f = l("promptInput"), p = h.dom.escape(t), x = h.dom.escape(e), y = `
          <div class="${d}">
            <p class="${u}">${p}</p>
            <input id="${o}" type="text" class="${f}" value="${x}">
          </div>
        `;
        this.open("confirm", {
          title: r.title || "Prompt",
          content: y,
          _isTrustedPrompt: !0,
          // Internal flag for prompt content
          onOk: () => {
            const v = document.getElementById(o);
            c(v ? v.value : null);
          },
          onCancel: () => c(null),
          ...r
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
    open(t, e = {}) {
      let r = "dialog", c = e;
      typeof t == "object" ? (c = t, r = c.type || "dialog") : r = t;
      const o = h.getConfig("box", c), d = o.boxId || Math.random().toString(36).slice(2, 10);
      let u = o.showOkButton, f = o.showCancelButton;
      r === "confirm" ? (u = u !== !1, f = f !== !1) : r === "dialog" && (u = u || !1, f = f || !1);
      const p = {
        id: d,
        type: r,
        title: o.title || "",
        content: o._isTrustedPrompt ? o.content || "" : h.dom.sanitize(o.content || ""),
        // Security: Always sanitize content (trustedHTML removed)
        // Exception: Internal prompt content is pre-escaped
        _isTrustedPrompt: o._isTrustedPrompt || !1,
        width: o.width !== void 0 ? o.width : o.widthClass ? void 0 : o.defaultWidth || 400,
        height: o.height,
        widthClass: o.widthClass,
        heightClass: o.heightClass,
        containerClass: o.containerClass,
        position: o.position || o.defaultPosition || "center",
        modal: o.modal !== void 0 ? o.modal : !0,
        closeOnEsc: o.closeOnEsc !== void 0 ? o.closeOnEsc : !1,
        closeOnOverlayClick: o.closeOnOverlayClick !== void 0 ? o.closeOnOverlayClick : !1,
        showCloseButton: o.showCloseButton !== void 0 ? o.showCloseButton : !0,
        showOkButton: u,
        showCancelButton: f,
        okText: o.okText || "OK",
        cancelText: o.cancelText || "Cancel",
        onOk: o.onOk,
        onCancel: o.onCancel,
        onBeforeClose: o.onBeforeClose,
        onClose: o.onClose,
        onContentLoad: o.onContentLoad,
        hxUrl: o.hxUrl,
        hxMethod: o.hxMethod || "GET",
        hxHeaders: o.hxHeaders,
        hxTrigger: o.hxTrigger || "load",
        // Internal state - start closed to trigger animation
        isOpen: !1,
        /**
         * Close this box dialog
         */
        close() {
          this.onBeforeClose && this.onBeforeClose(), this.isOpen = !1;
          const v = typeof document < "u" ? document.getElementById("uno-box-" + this.id) : null, b = () => {
            const C = i.store("uno")?.box;
            C && C.remove(this.id), this.onClose && this.onClose();
          };
          if (!v) {
            b();
            return;
          }
          let w = !1;
          const E = (C) => {
            C && C.target !== v || w || (w = !0, v.removeEventListener("transitionend", E), b());
          };
          v.addEventListener("transitionend", E), setTimeout(E, 350);
        }
      };
      this.items.push(p), this.ensureContainer();
      const x = this.items[this.items.length - 1], y = () => {
        x.isOpen = !0;
      };
      return i && typeof i.nextTick == "function" ? i.nextTick(y) : setTimeout(y, 0), x;
    },
    /**
     * Close a specific box by ID
     * @param {string} id - Box ID
     */
    close(t) {
      const e = this.items.find((r) => r.id === t);
      e && e.close();
    },
    /**
     * Remove a box from the items array
     * @param {string} id - Box ID
     */
    remove(t) {
      const e = this.items.findIndex((r) => r.id === t);
      e !== -1 && this.items.splice(e, 1);
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
      return N(t);
    },
    /**
     * Handle HTMX content load completion
     * @param {HTMLElement} el - Content element
     * @param {object} item - Box item
     */
    onHtmxContentLoaded(t, e) {
      e.onContentLoad && typeof e.onContentLoad == "function" && e.onContentLoad(t), window.Alpine && window.Alpine.initTree && window.Alpine.initTree(t);
    },
    /**
     * Ensure the box container exists in the DOM
     */
    ensureContainer() {
      const t = "uno-box-manager", e = typeof document < "u" && document.getElementById(t);
      if (e)
        return i && typeof i.initTree == "function" && !e._x_dataStack && i.initTree(e), e;
      const c = {
        "x-data": "{}",
        class: "relative",
        style: `z-index: ${h.options?.ui?.layers?.box || M("MODAL")};`,
        "x-show": "$store.uno.box.items.length > 0",
        "x-cloak": ""
        // Hide until Alpine is ready
      }, d = h.dom.ensureContainer(t, `
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
                 :style="'z-index:' + (Uno.options?.ui?.layers?.box || 1400)"
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
      `, c);
      return d && i && typeof i.initTree == "function" && i.initTree(d), d;
    }
  });
};
class et {
  /**
   * Create an UnoInstance.
   *
   * @param {Object} [options={}] - Global configuration options
   */
  constructor(s = {}) {
    this.version = "0.6.0", this.options = s, this.pluginDefaults = {}, this.plugins = [], this.initialized = !1, this._htmxRequestHandlers = [], this._htmxActiveXhrs = typeof WeakSet == "function" ? /* @__PURE__ */ new WeakSet() : null, this.htmx = {
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
      ajax: (n, t, e = {}) => typeof window > "u" || typeof window.htmx?.ajax != "function" ? Promise.reject(new Error("htmx.ajax is not available")) : window.htmx.ajax(n, t, e),
      /**
       * Fetch an HTML fragment without inserting it into the DOM.
       * Useful for getting content to display in a modal or tooltip.
       *
       * @param {string} method - HTTP method
       * @param {string} url - The URL to request
       * @param {Object} [ajaxOptions={}] - HTMX ajax options
       * @returns {Promise<string>} The HTML fragment string
       */
      fragment: async (n, t, e = {}) => {
        if (typeof document > "u") throw new Error("document is not available");
        const r = document.createElement("div");
        let c = "";
        const o = (d) => {
          d.preventDefault(), c = d.detail.serverResponse, d.detail.shouldSwap = !1;
        };
        r.addEventListener("htmx:beforeSwap", o);
        try {
          const d = e.swap || "innerHTML";
          return await this.htmx.ajax(n, t, { ...e, target: r, swap: d }), c;
        } finally {
          r.removeEventListener("htmx:beforeSwap", o), r.remove();
        }
      },
      /**
       * Register handlers for HTMX request lifecycle events.
       *
       * @param {Object} handlers - The handlers object
       * @param {Function} [handlers.start] - Called when a request starts
       * @param {Function} [handlers.end] - Called when a request ends (success, error, or timeout)
       */
      onRequest: (n = {}) => {
        if (!n || typeof n != "object" || (this._htmxRequestHandlers.push(n), typeof window > "u" || typeof document > "u") || !window.htmx || !P(document, "uno:htmx:request-hooks")) return;
        const t = (o) => o?.detail?.xhr, e = (o, d) => {
          for (const u of this._htmxRequestHandlers) {
            const f = u?.[o];
            if (typeof f == "function")
              try {
                f(d);
              } catch (p) {
                m.error("Uno.htmx.onRequest handler failed", p);
              }
          }
        }, r = (o) => {
          const d = t(o);
          if (d && this._htmxActiveXhrs) {
            if (this._htmxActiveXhrs.has(d)) return;
            this._htmxActiveXhrs.add(d);
          }
          e("start", o);
        }, c = (o) => {
          const d = t(o);
          if (d && this._htmxActiveXhrs) {
            if (!this._htmxActiveXhrs.has(d)) return;
            this._htmxActiveXhrs.delete(d);
          }
          e("end", o);
        };
        document.addEventListener("htmx:beforeRequest", r), document.addEventListener("htmx:afterRequest", c), document.addEventListener("htmx:responseError", c), document.addEventListener("htmx:sendError", c), document.addEventListener("htmx:timeout", c), document.addEventListener("htmx:xhr:progress", (o) => e("progress", o));
      },
      // Global HTMX event binding for notify functionality - progress bar only
      bindNotifyToHtmx: () => {
        typeof window > "u" || typeof document > "u" || !window.htmx || !window.Alpine || P(document, "uno:htmx:bind-notify") && (document.addEventListener("htmx:beforeRequest", () => {
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
    const l = {
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
        layers: { box: $.MODAL, tooltip: $.TOOLTIP, notify: $.TOAST },
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
    this.options = T(l, s);
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
  getConfig(s, a = {}) {
    const l = this.options || {};
    let n = this.pluginDefaults[s] ? { ...this.pluginDefaults[s] } : {};
    return l.plugins && l.plugins[s] && (n = T(n, l.plugins[s])), l[s] && (n = T(n, l[s])), a && (n = T(n, a)), n;
  }
  /**
   * ✅ Register plugin defaults
   *
   * @param {string} name - Component name
   * @param {Object} defaults - Default configuration object
   */
  registerPluginDefaults(s, a) {
    this.pluginDefaults[s] = a;
  }
  _getUiValue(s, a) {
    const l = this.options.ui?.[s];
    if (!l) return;
    if (typeof a != "string" || a.indexOf(".") === -1) return l?.[a];
    const n = a.split(".").filter(Boolean);
    let t = l;
    for (const e of n) {
      if (t == null) return;
      t = t[e];
    }
    return t;
  }
  _resolveUiTokens(s) {
    if (typeof s != "string") return s;
    const a = this.options.ui?.base;
    return typeof a == "string" && s.indexOf("{base}") !== -1 ? s.replaceAll("{base}", a) : s;
  }
  /**
   * Get style class from configuration or fallback
   * Support "+" prefix to append class to default
   * @param {string} comp - Component name
   * @param {string} elem - Element name
   * @param {string} defaultCls - Default class string
   * @returns {string} Final class string
   */
  getStyle(s, a, l) {
    const n = this._getUiValue(s, a);
    return n === void 0 ? this._resolveUiTokens(l) : typeof n == "string" && n.trim().indexOf("+") === 0 ? this._resolveUiTokens(`${l} ${n.trim().slice(1)}`.trim()) : this._resolveUiTokens(n);
  }
  /**
   * Set global configuration (before or after init)
   * @param {Object} config - Global configuration object
   */
  setConfig(s) {
    return this.options = T(this.options, s), this.initialized && m.info("Uno configuration updated", this.options), this;
  }
  /**
   * Register a plugin
   * @param {Function} plugin - Plugin function (Alpine, options) => void
   * @param {Object} options - Plugin options
   * @returns {UnoInstance}
   */
  use(s, a = {}) {
    return this.initialized && m.warn("Uno.use() called after initialization. Plugin might not work correctly."), this.plugins.push({ plugin: s, options: a }), this;
  }
  /**
   * Initialize Uno
   * @param {Object|Function} alpineOrOptions - Alpine instance or Configuration options
   * @param {Object} options - Configuration options (if first arg is Alpine)
   * @returns {Object} Alpine instance
   */
  init(s = {}, a = {}) {
    if (this.initialized)
      return m.warn("Uno has already been initialized."), q;
    let l = q, n = a;
    return !l && typeof window < "u" && window.Alpine && (l = window.Alpine), s && (s.version || typeof s.start == "function" || typeof s.store == "function") ? l = s : n = s || a, this.options = T(this.options, n), m.setDebug(this.options.debug !== !1), this.options.onError && m.setReporter(this.options.onError), m.info(`Uno v${this.version} initializing...`), typeof window < "u" && typeof document < "u" && window.htmx && typeof l?.initTree == "function" && P(document, "uno:htmx-afterSwap:initTree") && document.addEventListener("htmx:afterSwap", (t) => {
      const e = t?.detail?.elt;
      e && l.initTree(e);
    }), this.plugins.forEach(({ plugin: t, options: e }) => {
      try {
        const r = T(this.options, e);
        t(l, r);
      } catch (r) {
        m.error("Failed to load plugin", r);
      }
    }), this.initialized = !0, m.info("Uno initialized"), l;
  }
}
const h = new et(), nt = {
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
}, ot = {
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
}, it = {
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
const F = {
  en: nt,
  "zh-CN": ot,
  "zh-TW": it
}, V = Object.keys(F);
function st() {
  if (typeof localStorage < "u") {
    const i = localStorage.getItem("uno.lang");
    if (i && V.includes(i)) return i;
  }
  if (typeof navigator < "u") {
    const i = navigator.language, s = i.startsWith("zh") ? i.includes("TW") || i.includes("Hant") ? "zh-TW" : "zh-CN" : "en";
    return V.includes(s) ? s : "en";
  }
  return "en";
}
function rt(i, s) {
  const a = s.split(".");
  let l = i;
  for (const n of a) {
    if (!l || typeof l != "object" || !(n in l)) return;
    l = l[n];
  }
  return l;
}
function W(i, s) {
  for (const a in s)
    s[a] && typeof s[a] == "object" ? (i[a] || Object.assign(i, { [a]: {} }), W(i[a], s[a])) : Object.assign(i, { [a]: s[a] });
}
function at(i) {
  i.store("i18n", {
    locale: st(),
    messages: F,
    loading: !1,
    init() {
      typeof document < "u" && (document.documentElement.lang = this.locale);
    },
    t(s, a = {}) {
      const l = this.locale, n = this.messages[l] || this.messages.en, t = rt(n, s);
      return (typeof t == "string" ? t : s).replace(/\{(\w+)\}/g, (r, c) => c in a ? String(a[c]) : `{${c}}`);
    },
    setLocale(s) {
      if (!this.messages[s]) {
        m.warn(`[Uno] Language ${s} not found. Please load it first using loadLocale().`);
        return;
      }
      this.locale = s, typeof localStorage < "u" && localStorage.setItem("uno.lang", s), typeof document < "u" && (document.documentElement.lang = s);
    },
    // Alias for backward compatibility
    setLang(s) {
      this.setLocale(s);
    },
    async loadLocale(s, a) {
      if (!this.messages[s]) {
        this.loading = !0;
        try {
          const l = await fetch(a);
          if (!l.ok) throw new Error(`Failed to load locale: ${l.statusText}`);
          const n = await l.json();
          return this.extend({ [s]: n }), !0;
        } catch (l) {
          throw m.error(`[Uno] Failed to load locale ${s}`, l), l;
        } finally {
          this.loading = !1;
        }
      }
    },
    extend(s) {
      for (const a in s)
        this.messages[a] ? W(this.messages[a], s[a]) : this.messages[a] = s[a];
    }
  }), i.directive("t", (s, { expression: a, modifiers: l }, { effect: n }) => {
    const t = () => i.store("i18n");
    n(() => {
      const e = t();
      e.locale;
      const r = a, c = e.t(r);
      l.length > 0 ? l.forEach((o) => {
        s.setAttribute(o, c);
      }) : s.textContent = c;
    });
  });
}
const lt = (i, s = {}) => {
  const a = {
    placement: "top",
    trigger: "hover",
    delay: 300,
    offset: 8,
    arrow: !0,
    clickOutside: !0
  };
  h.registerPluginDefaults("tooltip", {
    ...a,
    ...s.plugins?.tooltip,
    ...s.tooltip
  });
  let l = 0, n = null;
  const t = () => {
    l || typeof requestAnimationFrame == "function" && (l = requestAnimationFrame(() => {
      l = 0;
      const o = i.store("uno")?.tooltip;
      o && o.updatePosition();
    }));
  }, e = () => {
    if (n || typeof window > "u" || typeof document > "u") return;
    const o = () => t(), d = () => t(), u = (p) => {
      if (p.key !== "Escape") return;
      const x = i.store("uno")?.tooltip;
      x?.visible && x.hide();
    }, f = (p) => {
      const x = i.store("uno")?.tooltip;
      if (!x?.visible || x.triggerMode !== "click" || x.clickOutside !== !0) return;
      const y = p.target, v = document.getElementById("uno-tooltip-content"), b = x.triggerEl && x.triggerEl.contains(y), w = v && v.contains(y);
      !b && !w && x.hide();
    };
    window.addEventListener("scroll", o, !0), window.addEventListener("resize", d), window.addEventListener("keydown", u), window.addEventListener("click", f), n = () => {
      window.removeEventListener("scroll", o, !0), window.removeEventListener("resize", d), window.removeEventListener("keydown", u), window.removeEventListener("click", f), n = null;
    };
  }, r = {
    visible: !1,
    content: "",
    placement: "top",
    top: -9999,
    left: -9999,
    zIndex: M("TOOLTIP"),
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
    show(o = {}) {
      m.info("Tooltip show called", o);
      const d = h.getConfig("tooltip", o);
      if (m.info("Tooltip resolved config", d), this.triggerEl = d.element || this.triggerEl, d.target) {
        const u = document.querySelector(d.target);
        u ? this.content = z(u.innerHTML) : (m.warn(`Tooltip target element not found: ${d.target}`), this.content = d.content || "");
      } else
        this.content = d.content || "";
      this.placement = d.placement, this.triggerMode = d.trigger, this.clickOutside = d.clickOutside, this.offset = d.offset, this.arrow = d.arrow !== !1, this.zIndex = h.options?.ui?.layers?.tooltip || M("TOOLTIP"), this.hxUrl = d.hxUrl, this.hxHeaders = d.hxHeaders, this.hxTrigger = d.hxTrigger || "load", this.visible = !0, this.ensureContainer(), e(), i.nextTick(() => {
        this.updatePosition();
      });
    },
    hide() {
      this.visible = !1, this.triggerEl = null, n && (n(), n = null), l && (cancelAnimationFrame(l), l = 0);
    },
    cancelHide() {
      this.hideTimeout && (clearTimeout(this.hideTimeout), this.hideTimeout = null);
    },
    scheduleHide(o = 300) {
      this.cancelHide(), this.hideTimeout = setTimeout(() => {
        this.hide();
      }, o);
    },
    updatePosition() {
      if (!this.visible || !this.triggerEl) return;
      const o = document.getElementById("uno-tooltip-content");
      if (!o) return;
      const d = Q(this.triggerEl, o, this.placement, {
        offset: this.offset,
        arrow: this.arrow
      });
      this.top = d.top, this.left = d.left, this.placement = d.placement, d.arrow && (this.arrowPosition = d.arrow);
    },
    ensureContainer() {
      const o = "uno-tooltip-manager";
      if (document.getElementById(o)) return;
      const d = h.options?.ui?.layers?.tooltip || M("TOOLTIP");
      H.ensureContainer(o, `
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
      `, {
        "x-data": "{}",
        style: `position: absolute; top: 0; left: 0; width: 0; height: 0; z-index: ${d};`
      });
      const f = document.getElementById(o);
      f && i && !f._x_dataStack && i.initTree(f);
    }
  };
  I(i, "tooltip", r);
  const c = (o, { expression: d, modifiers: u }, { evaluateLater: f, effect: p, cleanup: x }) => {
    const y = d ? f(d) : null, v = (g) => {
      const k = typeof g == "object" && g !== null ? g : g ? { content: g } : {}, L = {};
      return u.includes("top") && (L.placement = "top"), u.includes("bottom") && (L.placement = "bottom"), u.includes("left") && (L.placement = "left"), u.includes("right") && (L.placement = "right"), u.includes("click") && (L.trigger = "click"), u.includes("hover") && (L.trigger = "hover"), u.includes("focus") && (L.trigger = "focus"), { ...T(k, L), element: o };
    }, b = (g) => {
      const k = i.store("uno")?.tooltip;
      k && k.show(v(g));
    }, w = (g) => {
      const k = i.store("uno")?.tooltip;
      k && (k.visible && k.triggerEl === o ? k.hide() : b(g));
    };
    let E = null;
    y && p(() => {
      y((g) => {
        E = g;
      });
    });
    const C = () => {
      const g = v(E);
      h.getConfig("tooltip", g).trigger === "hover" && b(E);
    }, O = () => {
      const g = i.store("uno")?.tooltip;
      g && g.triggerMode === "hover" && g.triggerEl === o && g.scheduleHide();
    }, _ = (g) => {
      const k = v(E);
      h.getConfig("tooltip", k).trigger === "click" && (g.preventDefault(), w(E));
    }, R = () => {
      const g = v(E);
      h.getConfig("tooltip", g).trigger === "focus" && b(E);
    }, D = () => {
      const g = i.store("uno")?.tooltip;
      g && g.triggerMode === "focus" && g.triggerEl === o && g.hide();
    };
    o.addEventListener("mouseenter", C), o.addEventListener("mouseleave", O), o.addEventListener("click", _), o.addEventListener("focus", R), o.addEventListener("blur", D), x(() => {
      o.removeEventListener("mouseenter", C), o.removeEventListener("mouseleave", O), o.removeEventListener("click", _), o.removeEventListener("focus", R), o.removeEventListener("blur", D);
      const g = i.store("uno")?.tooltip;
      g && g.triggerEl === o && g.hide();
    });
  };
  i.directive("uno-tooltip", c);
}, ct = (i, s = {}) => {
  const a = {
    duration: 3e3,
    position: "top-right"
  };
  h.registerPluginDefaults("notify", {
    ...a,
    ...s.plugins?.notify,
    ...s.notify
  });
  const l = {
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
    onContentLoaded(n, t) {
      this._onHtmxEnd(t);
      const e = n.innerHTML || n.textContent || "", r = this._findNotificationById(t.id);
      r && (r.content = e), t.onContentLoad && t.onContentLoad(n), window.Alpine && window.Alpine.initTree && window.Alpine.initTree(n);
    },
    // Helper method to find notification by ID
    _findNotificationById(n) {
      for (const t of Object.keys(this.groups)) {
        const r = this.groups[t].find((c) => c.id === n);
        if (r) return r;
      }
      return null;
    },
    onHtmxError(n, t) {
      this._onHtmxEnd(t), console.error("HTMX Error:", n), this.error("Content load failed"), t.onError && t.onError(n);
    },
    _onHtmxBegin(n) {
      !n || n.__unoHtmxActive || (n.__unoHtmxActive = !0, this.progressRequests++, this.progressRequests === 1 && this.startProgress());
    },
    _onHtmxEnd(n) {
      !n || !n.__unoHtmxActive || (n.__unoHtmxActive = !1, this.progressRequests > 0 && this.progressRequests--, this.progressRequests === 0 && this.stopProgress());
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
        zIndex: h.options?.ui?.layers?.notify || M("TOAST")
      };
    },
    getContainerClass(n) {
      const t = h.getStyle("notify", "container"), e = J(n), r = G(n);
      return [t, e, r.flexDir, r.alignItems].join(" ");
    },
    ensureContainer() {
      const e = H.ensureContainer("uno-notify-container", `
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
        `, {
        "x-data": "{}"
        // Independent Alpine component
      });
      e && i.initTree && i.initTree(e);
      const r = "uno-global-progress", o = `
          <div :class="Uno.getStyle('notify', 'progressWrap')"
               style="z-index: ${M("MODAL") + 2100}"
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
      `, d = H.ensureContainer(r, o, {
        "x-data": "{}"
      });
      d && i.initTree && i.initTree(d);
    },
    show(n, t = {}) {
      this.ensureContainer();
      const e = h.getConfig("notify", t), r = this.nextId++;
      let c = n;
      i.store("i18n") && typeof i.store("i18n").t == "function" && (c = i.store("i18n").t(n));
      const o = e.duration !== void 0 ? e.duration : 3e3;
      m.info("[Notify] Show:", { id: r, message: c, duration: o, config: e });
      const d = {
        id: r,
        message: c,
        isVisible: !0,
        // For x-show
        closable: !0,
        duration: o,
        content: e.hxUrl ? "" : void 0,
        // Initialize content for HTMX-loaded content
        ...e
      }, u = e.position || "top-right";
      return this.groups[u] ? this.groups[u].push(d) : (this.groups["top-right"] || (this.groups["top-right"] = []), this.groups["top-right"].push(d)), o > 0 && setTimeout(() => {
        this.remove(r);
      }, o), r;
    },
    remove(n) {
      Object.keys(this.groups).forEach((t) => {
        const e = this.groups[t].findIndex((r) => r.id === n);
        e !== -1 && (this.groups[t][e].visible = !1, setTimeout(() => {
          const r = this.groups[t].findIndex((c) => c.id === n);
          r !== -1 && this.groups[t].splice(r, 1);
        }, 300));
      });
    },
    // 3. Proxy Methods
    success(n, t = {}) {
      return this.show(n, { type: "success", ...t });
    },
    error(n, t = {}) {
      return this.show(n, { type: "error", ...t });
    },
    warning(n, t = {}) {
      return this.show(n, { type: "warning", ...t });
    },
    info(n, t = {}) {
      return this.show(n, { type: "info", ...t });
    }
    // 4. Loading Helpers - Removed to simplify the component
    // Loading indicator functionality has been removed to reduce complexity
  };
  I(i, "notify", l), l.init(), i.directive("uno-notify", (n, { expression: t }, { evaluateLater: e }) => {
    if (!t) {
      console.warn("[Uno] x-uno-notify directive requires a configuration object or message string.");
      return;
    }
    const r = e(t), c = () => {
      const o = i.store("uno")?.notify;
      o && r((d) => {
        let u = "", f = {};
        if (typeof d == "string")
          u = d;
        else if (typeof d == "object" && d !== null)
          u = d.message || "", f = d;
        else {
          console.warn("[Uno] x-uno-notify directive value must be a string or object", d);
          return;
        }
        o.show(u, f);
      });
    };
    n.addEventListener("click", c);
  });
}, j = {
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
}, dt = (i, s = {}) => {
  const a = {
    realtime: !0,
    showError: !0,
    scrollOnError: !0
  };
  h.registerPluginDefaults("valid", {
    ...a,
    ...s.plugins?.valid,
    ...s.valid
  }), I(i, "valid", {
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
     * Merges configuration from 4 levels:
     * 1. Plugin Defaults
     * 2. Global User Config
     * 3. Global Plugin Config
     * 4. Instance Config (passed argument)
     *
     * @param {string} fieldName - The unique identifier for the field (name or id)
     * @param {Object} [instanceOpts={}] - Instance-specific configuration
     */
    addRules(n, t = {}) {
      this.rules[n] = h.getConfig("valid", t);
    },
    /**
     * Registers a custom validator function.
     *
     * @param {string} name - The name of the validator
     * @param {Function} validatorFn - The validator function. Receives (value, fieldId). Returns true, false, string (error msg), or Promise.
     */
    registerValidator(n, t) {
      typeof t == "function" ? this.customValidators[n] = t : m.error(`[Valid] Invalid validator function for ${n}`);
    },
    /**
     * Validates a single field against its rules.
     *
     * @param {string} fieldName - The field identifier
     * @param {string|number} value - The value to validate
     * @returns {boolean} True if valid, false otherwise
     */
    validate(n, t) {
      const e = this.rules[n];
      if (!e) return !0;
      let r = !0, c = "";
      const o = t == null ? "" : String(t), d = o.trim() === "";
      if (e.required && d && (r = !1, c = e.msgRequired || this.getI18nMsg("valid.required") || "This field is required"), r && d && !e.required)
        return delete this.errors[n], !0;
      if (r) {
        const u = e.type || (e.email ? "email" : null);
        u && j[u] && (j[u].test(o) || (r = !1, c = e[`msg${u.charAt(0).toUpperCase() + u.slice(1)}`] || e.msgType || this.getI18nMsg(`valid.${u}`) || `Invalid ${u} format`));
      }
      if (r && e.minLength !== null && e.minLength !== void 0 && o.length < e.minLength && (r = !1, c = e.msgMinLength || this.getI18nMsg("valid.min", { n: e.minLength }) || `Minimum length is ${e.minLength}`), r && e.maxLength !== null && e.maxLength !== void 0 && o.length > e.maxLength && (r = !1, c = e.msgMaxLength || this.getI18nMsg("valid.max", { n: e.maxLength }) || `Maximum length is ${e.maxLength}`), r && (e.min !== null && e.min !== void 0 || e.max !== null && e.max !== void 0))
        if (j.number.test(o)) {
          const u = parseFloat(o);
          e.min !== null && e.min !== void 0 && u < e.min && (r = !1, c = e.msgMin || this.getI18nMsg("valid.rangeMin", { n: e.min }) || `Value must be at least ${e.min}`), r && e.max !== null && e.max !== void 0 && u > e.max && (r = !1, c = e.msgMax || this.getI18nMsg("valid.rangeMax", { n: e.max }) || `Value must be at most ${e.max}`);
        } else
          r = !1, c = e.msgNumber || this.getI18nMsg("valid.number") || "Must be a number";
      if (r && e.pattern)
        try {
          (e.pattern instanceof RegExp ? e.pattern : new RegExp(e.pattern)).test(o) || (r = !1, c = e.msgPattern || this.getI18nMsg("valid.pattern") || "Invalid format");
        } catch (u) {
          m.error(`[Valid] Invalid regex pattern for ${n}`, u);
        }
      if (r && e.match) {
        const u = document.querySelector(
          e.match.startsWith("#") ? e.match : `[name="${e.match}"]`
        ) || document.getElementById(e.match);
        u && o !== u.value && (r = !1, c = e.msgMatch || this.getI18nMsg("valid.match") || "Fields do not match");
      }
      if (r && e.validator) {
        let u = e.validator;
        if (typeof e.validator == "string" && (u = this.customValidators[e.validator], u || m.warn(`[Valid] Validator "${e.validator}" not found`)), typeof u == "function") {
          const f = u(t, n);
          f === !1 ? (r = !1, c = e.msgValidator || "Validation failed") : typeof f == "string" ? (r = !1, c = f) : typeof f == "object" && f.valid === !1 && (r = !1, c = f.message || e.msgValidator || "Validation failed");
        }
      }
      return r ? delete this.errors[n] : this.errors[n] = c, r;
    },
    /**
     * Validates all fields in a form.
     *
     * @param {HTMLFormElement} form - The form element to validate
     * @returns {boolean} True if the form is valid
     */
    validateForm(n) {
      const t = n.querySelectorAll("input, select, textarea");
      let e = !0, r = null;
      return t.forEach((c) => {
        if (c.disabled) return;
        const o = c.name || c.id;
        this.rules[o] && (this.validate(o, c.value) || (e = !1, r || (r = c)));
      }), !e && r && (h.getConfig("valid") || {}).scrollOnError !== !1 && (r.scrollIntoView({ behavior: "smooth", block: "center" }), r.focus({ preventScroll: !0 })), e;
    },
    /**
     * Manually sets an error for a field.
     *
     * @param {string} fieldName - The field identifier
     * @param {string} message - The error message
     */
    showError(n, t) {
      this.errors[n] = t;
    },
    /**
     * Clears error for a specific field.
     *
     * @param {string} fieldName
     */
    hideError(n) {
      delete this.errors[n];
    },
    /**
     * Clears all validation errors.
     */
    clearErrors() {
      this.errors = {};
    },
    /**
     * Bulk sets errors (e.g., from server response).
     *
     * @param {Object.<string, string>} errorsObj - Object mapping field names to error messages
     */
    showErrors(n) {
      Object.assign(this.errors, n);
    },
    /**
     * Helper to get I18n message safely.
     *
     * @private
     * @param {string} key - Translation key
     * @param {Object} [params] - Translation parameters
     * @returns {string|null} The translated string or null
     */
    getI18nMsg(n, t) {
      const e = i.store("i18n");
      return e && typeof e.t == "function" ? e.t(n, t) : null;
    }
  }), i.directive("uno-valid", (n, { expression: t }, { evaluateLater: e, cleanup: r }) => {
    if (n.tagName === "FORM") return;
    const c = n.name || n.id;
    if (!c) {
      console.warn("[Uno Valid] Element has no name or id:", n);
      return;
    }
    const o = t ? e(t) : () => {
    }, u = {
      // Standard HTML5
      required: n.hasAttribute("required"),
      minLength: n.getAttribute("minlength") || null,
      maxLength: n.getAttribute("maxlength") || null,
      min: n.getAttribute("min") || null,
      max: n.getAttribute("max") || null,
      pattern: n.getAttribute("pattern") || null,
      type: n.getAttribute("type") || null
    }, f = (y) => {
      i.store("uno").valid.addRules(c, { ...u, ...y });
    };
    t ? o((y) => {
      f(typeof y == "object" ? y : {});
    }) : f({});
    const p = () => {
      const y = i.store("uno").valid.rules[c];
      y && y.realtime !== !1 && i.store("uno").valid.validate(c, n.value);
    }, x = () => {
      i.store("uno").valid.validate(c, n.value);
    };
    n.addEventListener("input", p), n.addEventListener("blur", x), r(() => {
      n.removeEventListener("input", p), n.removeEventListener("blur", x);
    });
  }), typeof document < "u" && (document.body.addEventListener("htmx:confirm", (n) => {
    const t = n.detail.elt, e = t.tagName === "FORM" ? t : t.closest("form");
    e && e.hasAttribute("x-uno-valid") && (i.store("uno").valid.validateForm(e) || (n.preventDefault(), n.stopPropagation(), m.info("[Valid] Form validation failed. Request blocked."), e.dispatchEvent(
      new CustomEvent("uno:validation-failed", {
        bubbles: !0,
        detail: { errors: i.store("uno").valid.errors }
      })
    )));
  }), document.body.addEventListener("htmx:afterRequest", (n) => {
    const t = n.detail.xhr;
    if (t.status === 422)
      try {
        const e = JSON.parse(t.responseText);
        e.errors && (i.store("uno").valid.showErrors(e.errors), m.info("[Valid] Server validation errors applied", e.errors));
      } catch {
      }
  }));
};
const ut = (i, s = {}) => {
  const a = {
    defaultTab: "",
    hash: !0,
    activeClass: "is-active",
    inactiveClass: ""
  };
  h.registerPluginDefaults("tab", { ...a, ...s.plugins?.tab, ...s.tab }), I(i, "tab", {
    // State: { [group]: activeTabName }
    groups: {},
    // Registry: { [tabName]: { group, name, ... } }
    registry: {},
    init() {
      const n = () => {
        let t = window.location.hash.slice(1);
        if (!t && window.location.href.includes("#") && (t = window.location.href.split("#")[1]), m.info(`[Tab] checkHash: ${t}`), t && this.registry[t]) {
          const { group: e, name: r } = this.registry[t];
          m.info(`[Tab] checkHash match: ${r}`), this.activate({ group: e, name: r, updateHash: !1 });
        }
      };
      n(), window.addEventListener("hashchange", n), window.addEventListener("load", n);
    },
    /**
     * Activate a tab
     * @param {object} instanceOpts - Options { name, group, ... }
     */
    activate(n = {}) {
      const t = h.getConfig("tab", {
        group: "default",
        hash: !0,
        ...n
      });
      t.name && (m.info(
        `[Tab] Activating ${t.name} in ${t.group}. Hash update: ${t.updateHash !== !1}`
      ), this.registry[t.name] || (this.registry[t.name] = { group: t.group, name: t.name }), this.groups[t.group] = t.name, t.hash && t.updateHash !== !1 && window.location.hash.slice(1) !== t.name && window.history.replaceState(null, null, `#${t.name}`), t.updateHash === !1 && setTimeout(() => {
        const e = document.getElementById(`tab-${t.name}`);
        e && e.hasAttribute("hx-get") && (m.info(`[Tab] Triggering HTMX for ${t.name}`), window.htmx && window.htmx.trigger(e, "click"));
      }, 0));
    },
    /**
     * Check if a tab is active
     * @param {string} name - Tab name
     * @param {string} group - Group name
     * @returns {boolean}
     */
    isActive(n, t = "default") {
      return this.groups[t] === n;
    },
    /**
     * Register a tab (helper for directives or manual setup)
     * @param {string} name
     * @param {string} group
     */
    register(n, t = "default") {
      m.info(`[Tab] Registering ${n} in ${t}. Hash: ${window.location.hash}`), this.registry[n] = { name: n, group: t }, this.groups[t] || (this.groups[t] = n);
      let e = window.location.hash.slice(1);
      !e && window.location.href.includes("#") && (e = window.location.href.split("#")[1]), e === n && (m.info(`[Tab] Hash match for ${n}. Activating.`), this.activate({ group: t, name: n, updateHash: !1 }));
    }
  }), i.store("uno").tab.init(), i.directive("uno-tab", (n, { expression: t }, { evaluateLater: e, cleanup: r }) => {
    if (!t) {
      console.warn("[Uno] x-uno-tab directive requires a configuration object.");
      return;
    }
    const c = e(t);
    c((d) => {
      let u = {};
      if (typeof d == "object" && d !== null)
        u = d;
      else if (typeof d == "string")
        u = { name: d };
      else {
        console.warn("[Uno] x-uno-tab directive value must be a string or object", d);
        return;
      }
      u.group || (u.group = "default"), u.name && i.store("uno").tab.register(u.name, u.group);
    });
    const o = (d) => {
      d.preventDefault(), c((u) => {
        let f = {};
        typeof u == "object" && u !== null ? f = u : typeof u == "string" && (f = { name: u }), f.name && i.store("uno").tab.activate(f);
      });
    };
    n.addEventListener("click", o), r(() => n.removeEventListener("click", o));
  });
};
h.dom = H;
typeof window < "u" && (window.Uno = h);
h.use(at).use(lt).use(ct).use(tt).use(dt).use(ut);
if (typeof window < "u" && typeof document < "u") {
  const i = () => {
    if (h.initialized) return;
    const s = window.UnoConfig || {};
    h.init(window.Alpine, s);
  };
  window.Alpine ? i() : document.addEventListener("alpine:init", i);
}
export {
  H as DomManager,
  h as Uno,
  tt as box,
  h as default,
  ct as notify,
  ut as tab,
  lt as tooltip,
  dt as valid
};
