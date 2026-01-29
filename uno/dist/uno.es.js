import M from "alpinejs";
let $ = !1, z = null;
const f = {
  /**
   * Set debug mode
   * @param {boolean} debug
   */
  setDebug(s) {
    $ = s;
  },
  /**
   * Set custom error reporter
   * @param {Function} reporter - (error, ...args) => void
   */
  setReporter(s) {
    z = s;
  },
  /**
   * Log info message
   * @param {...any} args - Arguments to log
   */
  info(...s) {
    $ && typeof console < "u" && typeof console.log == "function" && console.log("[Uno Info]", ...s);
  },
  /**
   * Log debug message
   * @param {...any} args - Arguments to log
   */
  debug(...s) {
    $ && typeof console < "u" && typeof console.log == "function" && console.log("[Uno Debug]", ...s);
  },
  /**
   * Log warning message
   * @param {...any} args - Arguments to log
   */
  warn(...s) {
    $ && typeof console < "u" && typeof console.warn == "function" && console.warn("[Uno Warn]", ...s);
  },
  /**
   * Log error message
   * @param {...any} args - Arguments to log
   */
  error(...s) {
    if (z)
      try {
        z(...s);
      } catch (t) {
        $ && typeof console < "u" && typeof console.error == "function" && console.error("[Uno Error] Reporter failed:", t);
      }
    $ && typeof console < "u" && typeof console.error == "function" && console.error("[Uno Error]", ...s);
  }
}, R = /* @__PURE__ */ new WeakMap();
function H(s, t) {
  if (!s) return !1;
  const n = String(t);
  let e = R.get(s);
  return e || (e = /* @__PURE__ */ new Set(), R.set(s, e)), e.has(n) ? !1 : (e.add(n), !0);
}
function P(s) {
  return s ? s.replace(/<script\b[^>]*>([\s\S]*?)<\/script>/gim, "").replace(/<img\b[^>]*>/gim, "").replace(/\s+on[a-z]+="[^"]*"/gim, "").replace(/\s+on[a-z]+='[^']*'/gim, "").replace(/\s+href="javascript:[^"]*"/gim, "").replace(/\s+href='javascript:[^']*'/gim, "").replace(/javascript:/gim, "") : "";
}
function A(s, t, { trusted: n = !1, alpine: e } = {}) {
  s && (s.innerHTML = n ? String(t ?? "") : P(String(t ?? "")), e?.initTree?.(s));
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
  ensureContainer(s, t = "", n = {}, e = document.body) {
    const i = document.getElementById(s);
    if (i) return i;
    try {
      const o = document.createElement("div");
      return o.id = s, window.Uno?.options?.csp?.nonce && !n.nonce && (n.nonce = window.Uno.options.csp.nonce), Object.entries(n).forEach(([r, l]) => {
        o.setAttribute(r, l);
      }), t && (o.innerHTML = t), e ? e.appendChild(o) : document.documentElement.appendChild(o), o;
    } catch (o) {
      return f.error(`DomManager: Failed to inject container #${s}`, o), null;
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
  injectStyles(s, t) {
    if (!document.getElementById(s))
      try {
        const n = document.createElement("style");
        n.id = s, window.Uno?.options?.csp?.nonce && n.setAttribute("nonce", window.Uno.options.csp.nonce), n.textContent = t, document.head.appendChild(n);
      } catch (n) {
        f.error(`DomManager: Failed to inject styles #${s}`, n);
      }
  },
  // Exports
  setHTML: A,
  // Secure HTML setting
  sanitize: P
  // HTML sanitization
};
function L(s) {
  return Object.prototype.toString.call(s) === "[object Object]";
}
function C(s, t) {
  if (!L(s)) return L(t) ? { ...t } : t;
  if (!L(t)) return t === void 0 ? { ...s } : t;
  const n = { ...s };
  for (const [e, i] of Object.entries(t)) {
    if (L(i) && L(n[e])) {
      n[e] = C(n[e], i);
      continue;
    }
    if (L(i)) {
      n[e] = C({}, i);
      continue;
    }
    if (typeof i == "string" && i.trim().indexOf("+") === 0 && typeof n[e] == "string") {
      n[e] = `${n[e]} ${i.trim().slice(1)}`.trim();
      continue;
    }
    n[e] = i;
  }
  return n;
}
const T = {
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
}, k = (s) => {
  if (!s) return T.BASE;
  const t = s.toUpperCase();
  return T[t] !== void 0 ? T[t] : T.BASE;
};
class V {
  /**
   * Create an UnoInstance.
   *
   * @param {Object} [options={}] - Global configuration options
   */
  constructor(t = {}) {
    this.version = "0.8.0", this.options = t, this.pluginDefaults = {}, this.plugins = [], this.initialized = !1, this._preInitConfig = {}, this._htmxRequestHandlers = [], this._htmxActiveXhrs = typeof WeakSet == "function" ? /* @__PURE__ */ new WeakSet() : null, this.htmx = {
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
      ajax: (i, o, r = {}) => typeof window > "u" || typeof window.htmx?.ajax != "function" ? Promise.reject(new Error("htmx.ajax is not available")) : window.htmx.ajax(i, o, r),
      /**
       * Fetch an HTML fragment without inserting it into the DOM.
       * Useful for getting content to display in a modal or tooltip.
       *
       * @param {string} method - HTTP method
       * @param {string} url - The URL to request
       * @param {Object} [ajaxOptions={}] - HTMX ajax options
       * @returns {Promise<string>} The HTML fragment string
       */
      fragment: async (i, o, r = {}) => {
        if (typeof document > "u") throw new Error("document is not available");
        const l = document.createElement("div");
        let c = "";
        const a = (u) => {
          u.preventDefault(), c = u.detail.serverResponse, u.detail.shouldSwap = !1;
        };
        l.addEventListener("htmx:beforeSwap", a);
        try {
          const u = r.swap || "innerHTML";
          return await this.htmx.ajax(i, o, { ...r, target: l, swap: u }), c;
        } catch (u) {
          throw f.error("HTMX fragment fetch failed:", { method: i, url: o, error: u }), u;
        } finally {
          l.removeEventListener("htmx:beforeSwap", a), l.remove();
        }
      },
      /**
       * Register handlers for HTMX request lifecycle events.
       *
       * @param {Object} handlers - The handlers object
       * @param {Function} [handlers.start] - Called when a request starts
       * @param {Function} [handlers.end] - Called when a request ends (success, error, or timeout)
       */
      onRequest: (i = {}) => {
        if (!i || typeof i != "object" || (this._htmxRequestHandlers.push(i), typeof window > "u" || typeof document > "u") || !window.htmx || !H(document, "uno:htmx:request-hooks")) return;
        const o = (a) => a?.detail?.xhr, r = (a, u) => {
          for (const w of this._htmxRequestHandlers) {
            const g = w?.[a];
            if (typeof g == "function")
              try {
                g(u);
              } catch (x) {
                f.error("Uno.htmx.onRequest handler failed", x);
              }
          }
        }, l = (a) => {
          const u = o(a);
          if (u && this._htmxActiveXhrs) {
            if (this._htmxActiveXhrs.has(u)) return;
            this._htmxActiveXhrs.add(u);
          }
          r("start", a);
        }, c = (a) => {
          const u = o(a);
          if (u && this._htmxActiveXhrs) {
            if (!this._htmxActiveXhrs.has(u)) return;
            this._htmxActiveXhrs.delete(u);
          }
          r("end", a);
        };
        document.addEventListener("htmx:beforeRequest", l), document.addEventListener("htmx:afterRequest", c), document.addEventListener("htmx:responseError", c), document.addEventListener("htmx:sendError", c), document.addEventListener("htmx:timeout", c), document.addEventListener("htmx:xhr:progress", (a) => r("progress", a));
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
        layers: { box: T.BOX, tooltip: T.TOOLTIP, notify: T.NOTIFY },
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
          helperText: "text-red-600 dark:text-red-400 text-xs mt-1"
        }
      }
    };
    this.options = C(e, t);
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
    let i = this.pluginDefaults[t] ? { ...this.pluginDefaults[t] } : {};
    return e.plugins && e.plugins[t] && (i = C(i, e.plugins[t])), e[t] && (i = C(i, e[t])), n && (i = C(i, n)), i;
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
    const i = n.split(".").filter(Boolean);
    let o = e;
    for (const r of i) {
      if (o == null) return;
      o = o[r];
    }
    return o;
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
    const i = this._getUiValue(t, n);
    return i === void 0 ? this._resolveUiTokens(e) : typeof i == "string" && i.trim().indexOf("+") === 0 ? this._resolveUiTokens(`${e} ${i.trim().slice(1)}`.trim()) : this._resolveUiTokens(i);
  }
  /**
   * Set configuration before or after initialization
   *
   * This is the primary way to configure Uno. Can be called:
   * 1. Before Uno auto-initializes (recommended)
   * 2. After initialization (for runtime config updates)
   *
   * @param {Object} config - Configuration object
   * @returns {UnoInstance} this (for chaining)
   *
   * @example
   * // Before initialization
   * Uno.setConfig({
   *   debug: true,
   *   ui: { base: 'rounded-lg shadow-md' }
   * })
   *
   * // Chaining
   * Uno.setConfig({ debug: true }).setConfig({ ui: { base: 'shadow' } })
   */
  setConfig(t) {
    return !t || typeof t != "object" ? (f.warn("setConfig called with invalid argument, expected object"), this) : (this.initialized ? (f.info("Uno configuration updated after initialization"), this.options = C(this.options, t)) : this._preInitConfig = C(this._preInitConfig, t), this);
  }
  /**
   * Get current global configuration
   * @returns {Object} Current global configuration object
   */
  getGlobalConfig() {
    return this.initialized ? this.options : this._preInitConfig;
  }
  /**
   * Register a plugin
   * @param {Function} plugin - Plugin function (Alpine, options) => void
   * @param {Object} options - Plugin options
   * @returns {UnoInstance}
   */
  use(t, n = {}) {
    return this.initialized && f.warn("Uno.use() called after initialization. Plugin might not work correctly."), this.plugins.push({ plugin: t, options: n }), this;
  }
  /**
   * Initialize Uno - Automatically waits for Alpine if needed
   *
   * This method intelligently handles Alpine initialization:
   * 1. If Alpine is already initialized → Initialize Uno immediately
   * 2. If Alpine is not ready → Wait for alpine:init event
   *
   * Usage:
   *   Uno.init()                    // Uses window.Alpine automatically
   *   Uno.init(config)              // Pass configuration
   *   Uno.init(Alpine, config)      // Explicit Alpine instance (legacy)
   *
   * @param {Object|Function} alpineOrOptions - Alpine instance or Configuration options
   * @param {Object} options - Configuration options (if first arg is Alpine)
   * @returns {Promise<Object>} Promise resolving to Alpine instance
   */
  async init(t = {}, n = {}) {
    if (this.initialized)
      return f.warn("Uno has already been initialized."), M || (typeof window < "u" ? window.Alpine : null);
    let e = M;
    return !e && typeof window < "u" && window.Alpine && (e = window.Alpine), t && (t.version || typeof t.start == "function" || typeof t.store == "function") && (e = t), this.options = C(this.options, this._preInitConfig), f.setDebug(this.options.debug !== !1), this.options.onError && f.setReporter(this.options.onError), e && typeof e.store == "function" ? (this._doInit(e), e) : typeof document > "u" ? (f.error("Alpine not found and document is unavailable. Cannot initialize Uno."), null) : (f.debug("Alpine not ready, waiting for alpine:init event..."), new Promise((o) => {
      document.addEventListener(
        "alpine:init",
        () => {
          if (e = window.Alpine || M, !e) {
            f.error("Alpine still not found after alpine:init event"), o(null);
            return;
          }
          f.debug("Alpine ready, continuing Uno initialization"), this._doInit(e), o(e);
        },
        { once: !0 }
      );
    }));
  }
  /**
   * Internal initialization logic
   * @private
   * @param {Object} alpine - Alpine instance
   */
  _doInit(t) {
    f.info(`Uno v${this.version} initializing...`), typeof window < "u" && typeof document < "u" && window.htmx && typeof t?.initTree == "function" && H(document, "uno:htmx-afterSwap:initTree") && document.addEventListener("htmx:afterSwap", (n) => {
      const e = n?.detail?.elt;
      e && t.initTree(e);
    }), this.plugins.forEach(({ plugin: n, options: e }) => {
      try {
        const i = C(this.options, e);
        n(t, i);
      } catch (i) {
        f.error("Failed to load plugin", i);
      }
    }), this._setupGlobalErrorHandlers(), this.initialized = !0, f.info("Uno initialized");
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
    if (!(typeof window > "u") && H(window, "uno:global-error-handlers")) {
      const t = typeof WeakSet == "function" ? /* @__PURE__ */ new WeakSet() : null;
      window.addEventListener("unhandledrejection", (n) => {
        const e = n.reason;
        if (t && e instanceof Error) {
          if (t.has(e)) return;
          t.add(e);
        }
        f.error("Unhandled Promise Rejection:", e);
      }), window.addEventListener("error", (n) => {
        const e = n.error;
        if (t && e instanceof Error) {
          if (t.has(e)) return;
          t.add(e);
        }
        f.error("Uncaught Error:", {
          message: n.message,
          filename: n.filename,
          lineno: n.lineno,
          colno: n.colno,
          error: e
        });
      }), f.debug("Global error handlers installed");
    }
  }
}
const h = new V(), W = {
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
    validname: "Letters, numbers, - and _ only (3-16 chars)",
    en: "English letters only",
    domain: "Invalid domain format",
    mobile: "Invalid mobile number",
    int: "Integer required",
    email: "Invalid email address",
    date: "Invalid date format (YYYY-MM-DD)",
    url: "Invalid URL",
    pwd: "Password must be 8-80 chars, contain letters and numbers, and can include _@#$%^&*()+-=",
    strongpwd: "Password must be 8-80 chars, contain letters, numbers, and at least 1 special char (_@#$%^&*()+-=)"
  },
  charCount: "Remaining: {count}"
}, X = {
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
    validname: "仅限字母、数字、- 和 _ (3-16 个字符)",
    en: "仅限英文字母",
    domain: "无效的域名格式",
    mobile: "无效的手机号",
    int: "需要整数",
    email: "无效的邮箱地址",
    date: "无效的日期格式 (YYYY-MM-DD)",
    url: "无效的 URL",
    pwd: "密码需 8-80 位，必须包含字母和数字，可含 _@#$%^&*()+-= 特殊字符",
    strongpwd: "密码需 8-80 位，必须包含字母、数字和至少 1 个特殊字符 (_@#$%^&*()+-=)"
  },
  charCount: "剩余: {count}"
}, Y = {
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
    pwd: "密碼需 8-80 位，必須包含字母和數字，可含 _@#$%^&*()+-= 特殊字符",
    strongpwd: "密碼需 8-80 位，必須包含字母、數字和至少 1 個特殊字符 (_@#$%^&*()+-=)"
  },
  charCount: "剩餘: {count}"
};
const _ = {
  en: W,
  "zh-CN": X,
  "zh-TW": Y
}, D = Object.keys(_);
function Z() {
  if (typeof localStorage < "u") {
    const s = localStorage.getItem("uno.lang");
    if (s && D.includes(s)) return s;
  }
  if (typeof navigator < "u") {
    const s = navigator.language, t = s.startsWith("zh") ? s.includes("TW") || s.includes("Hant") ? "zh-TW" : "zh-CN" : "en";
    return D.includes(t) ? t : "en";
  }
  return "en";
}
function N(s, t) {
  const n = t.split(".");
  let e = s;
  for (const i of n) {
    if (!e || typeof e != "object" || !(i in e)) return;
    e = e[i];
  }
  return e;
}
function j(s, t) {
  for (const n in t)
    t[n] && typeof t[n] == "object" ? (s[n] || Object.assign(s, { [n]: {} }), j(s[n], t[n])) : Object.assign(s, { [n]: t[n] });
}
function K(s) {
  const t = s.store("i18n"), n = {
    locale: Z(),
    messages: _,
    loading: !1,
    init() {
      typeof document < "u" && (document.documentElement.lang = this.locale);
    },
    t(e, i = {}) {
      const o = this.locale, r = this.messages[o] || this.messages.en, l = N(r, e);
      return (typeof l == "string" ? l : e).replace(/\{(\w+)\}/g, (a, u) => u in i ? String(i[u]) : `{${u}}`);
    },
    setLocale(e) {
      if (!this.messages[e]) {
        f.warn(`[Uno] Language ${e} not found. Please load it first using loadLocale().`);
        return;
      }
      this.locale = e, typeof localStorage < "u" && localStorage.setItem("uno.lang", e), typeof document < "u" && (document.documentElement.lang = e);
    },
    // Alias for backward compatibility
    setLang(e) {
      this.setLocale(e);
    },
    extend(e) {
      j(this.messages, e);
    },
    async loadLocale(e, i) {
      if (!this.messages[e]) {
        this.loading = !0;
        try {
          const o = await fetch(i);
          if (!o.ok) throw new Error(`Failed to load ${e}`);
          const r = await o.json();
          this.messages[e] = r, this.setLocale(e);
        } catch (o) {
          f.error(`[Uno] Failed to load locale ${e}:`, o);
        } finally {
          this.loading = !1;
        }
      }
    }
  };
  t ? (f.info("[Uno] Merging into existing i18n store"), t.messages || (t.messages = {}), j(t.messages, _), t.t || (t.t = n.t), t.setLocale || (t.setLocale = n.setLocale), t.setLang || (t.setLang = n.setLang), t.loadLocale || (t.loadLocale = n.loadLocale)) : s.store("i18n", n), s.directive("t", (e, { expression: i, modifiers: o }, { effect: r }) => {
    const l = () => s.store("i18n");
    r(() => {
      const c = l();
      c.locale;
      const a = i, u = c.t(a);
      o.length > 0 ? o.forEach((w) => {
        e.setAttribute(w, u);
      }) : e.textContent = u;
    });
  });
}
function U(s, t, n) {
  if (!s || typeof s.store != "function") {
    f.error("[Uno] Alpine.store is not available");
    return;
  }
  const e = s.store("uno");
  e ? s.store("uno", {
    ...e,
    [t]: n
  }) : s.store("uno", { [t]: n }), s.store(t) || s.store(t, s.store("uno")?.[t]);
}
const G = (s) => {
  const t = s || "center", n = {
    center: "items-center justify-center",
    "top-center": "items-start justify-center",
    "top-left": "items-start justify-start sm:px-8",
    "top-right": "items-start justify-end sm:px-8",
    "bottom-center": "items-end justify-center",
    "bottom-left": "items-end justify-start sm:px-8",
    "bottom-right": "items-end justify-end sm:px-8",
    "center-left": "items-center justify-start sm:px-8",
    "center-right": "items-center justify-end sm:px-8"
  };
  return n[t] || n.center;
}, J = (s) => {
  const t = s || "top-right", n = {
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
}, Q = (s) => {
  const t = s || "top-right";
  let n = "flex-col", e = "items-end", i = "space-y-2";
  return t.includes("bottom") && (n = "flex-col-reverse", i = "space-y-reverse space-y-2"), t.includes("left") ? e = "items-start" : t.includes("center") ? t === "center" ? e = "items-center justify-center" : e = "items-center" : e = "items-end", { flexDir: n, alignItems: e, spaceY: i };
};
function tt(s, t, n) {
  const e = s.getBoundingClientRect(), i = t.getBoundingClientRect(), o = i.width, r = i.height, l = e.width, c = e.height, a = { top: e.top, left: e.left }, u = window.innerWidth, w = window.innerHeight, g = 0, x = 0, b = w, y = u, m = (E, B) => {
    if (B) {
      if (E + o > y)
        return y - o - 10;
      if (E < x)
        return x + 10;
    } else {
      if (E + r > b)
        return b - r - 10;
      if (E < g)
        return g + 10;
    }
    return E;
  };
  let v, d, p = n;
  switch (n) {
    case "left":
      d = a.left - o - 5, d < x && (d = a.left + l + 5, p = "right"), v = m(a.top + c / 2 - r / 2, !1);
      break;
    case "top":
      v = a.top - r - 8, v < g && (v = a.top + c + 8, p = "bottom"), d = a.left + (l - o) / 2, d = m(d, !0);
      break;
    case "right":
      d = a.left + l + 5, d + o > y && (d = a.left - o - 5, p = "left"), v = m(a.top + c / 2 - r / 2, !1);
      break;
    default:
      v = a.top + c + 8, v + r > b && (v = a.top - r - 8, p = "top"), d = a.left + (l - o) / 2, d = m(d, !0);
      break;
  }
  return { top: v, left: d, position: p };
}
const et = () => (h.options?.ui?.layers?.tooltip || k("TOOLTIP"), `
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
let S = 0, I = null;
const nt = (s) => {
  const t = () => {
    S || typeof requestAnimationFrame == "function" && (S = requestAnimationFrame(() => {
      S = 0;
      const e = s.store("uno")?.tooltip;
      e && e.updatePosition();
    }));
  }, n = () => {
    if (I || typeof window > "u" || typeof document > "u") return;
    const e = () => t(), i = () => t(), o = (l) => {
      if (l.key !== "Escape") return;
      const c = s.store("uno")?.tooltip;
      c?.visible && c.hide();
    }, r = (l) => {
      const c = s.store("uno")?.tooltip;
      if (!c?.visible || c.triggerMode !== "click" || c.clickOutside !== !0) return;
      const a = l.target, u = document.getElementById("uno-tooltip-content"), w = c.triggerEl && c.triggerEl.contains(a), g = u && u.contains(a);
      !w && !g && c.hide();
    };
    window.addEventListener("scroll", e, !0), window.addEventListener("resize", i), window.addEventListener("keydown", o), window.addEventListener("click", r), I = () => {
      window.removeEventListener("scroll", e, !0), window.removeEventListener("resize", i), window.removeEventListener("keydown", o), window.removeEventListener("click", r), I = null;
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
    /**
     * Active tooltip ID
     * @type {string|null}
     */
    active: null,
    show(e = {}) {
      f.info("Tooltip show called", e);
      const i = h.getConfig("tooltip", e);
      if (f.info("Tooltip resolved config", i), !i.target && typeof i.content == "string" && i.content.startsWith("#") && (i.content.includes(" ") || document.querySelector(i.content) && (i.target = i.content)), this.triggerEl = i.element || this.triggerEl, i.target) {
        const o = document.querySelector(i.target);
        o ? this.content = P(o.innerHTML) : (f.warn(`Tooltip target element not found: ${i.target}`), this.content = i.content || "");
      } else
        this.content = i.content || "";
      this.placement = i.placement, this.triggerMode = i.trigger, this.clickOutside = i.clickOutside, this.offset = i.offset, this.arrow = i.arrow !== !1, this.zIndex = h.options?.ui?.layers?.tooltip || k("TOOLTIP"), this.hxUrl = i.hxUrl, this.hxHeaders = i.hxHeaders, this.hxTrigger = i.hxTrigger || "load", this.visible = !0, this.ensureContainer(), n(), s.nextTick(() => {
        this.updatePosition();
      });
    },
    hide() {
      this.visible = !1, this.triggerEl = null, I && (I(), I = null), S && (cancelAnimationFrame(S), S = 0);
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
      const i = tt(this.triggerEl, e, this.placement, {
        offset: this.offset,
        arrow: this.arrow
      });
      this.top = i.top, this.left = i.left, this.placement = i.placement, i.arrow && (this.arrowPosition = i.arrow);
    },
    ensureContainer() {
      const e = "uno-tooltip-manager";
      if (document.getElementById(e)) return;
      const i = h.options?.ui?.layers?.tooltip || k("TOOLTIP"), o = et();
      O.ensureContainer(e, o, {
        "x-data": "{}",
        style: `position: absolute; top: 0; left: 0; width: 0; height: 0; z-index: ${i};`
      });
      const r = document.getElementById(e);
      r && s && !r._x_dataStack && s.initTree(r);
    }
  };
}, ot = (s) => {
  const t = (n, { expression: e, modifiers: i }, { evaluateLater: o, effect: r, cleanup: l }) => {
    const c = e ? o(e) : null, a = (d) => (i.length > 0 && f.warn(
      "[Uno] x-uno-tooltip modifiers are deprecated in v2.0. Use configuration object instead."
    ), { ...typeof d == "object" && d !== null ? d : d ? { content: d } : {}, element: n }), u = (d) => {
      const p = s.store("uno")?.tooltip;
      p && p.show(a(d));
    }, w = (d) => {
      const p = s.store("uno")?.tooltip;
      p && (p.visible && p.triggerEl === n ? p.hide() : u(d));
    };
    let g = null;
    c && r(() => {
      c((d) => {
        g = d;
      });
    });
    const x = () => {
      if (e && g === null) return;
      const d = a(g);
      h.getConfig("tooltip", d).trigger === "hover" && u(g);
    }, b = () => {
      const d = s.store("uno")?.tooltip;
      d && d.triggerMode === "hover" && d.triggerEl === n && d.scheduleHide();
    }, y = (d) => {
      const p = a(g);
      h.getConfig("tooltip", p).trigger === "click" && (d.preventDefault(), w(g));
    }, m = () => {
      const d = a(g);
      h.getConfig("tooltip", d).trigger === "focus" && u(g);
    }, v = () => {
      const d = s.store("uno")?.tooltip;
      d && d.triggerMode === "focus" && d.triggerEl === n && d.hide();
    };
    n.addEventListener("mouseenter", x), n.addEventListener("mouseleave", b), n.addEventListener("click", y), n.addEventListener("focus", m), n.addEventListener("blur", v), l(() => {
      n.removeEventListener("mouseenter", x), n.removeEventListener("mouseleave", b), n.removeEventListener("click", y), n.removeEventListener("focus", m), n.removeEventListener("blur", v);
      const d = s.store("uno")?.tooltip;
      d && d.triggerEl === n && d.hide();
    });
  };
  s.directive("uno-tooltip", t);
};
const it = (s, t = {}) => {
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
  const e = nt(s);
  U(s, "tooltip", e), ot(s);
}, st = () => `
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
`, rt = () => `
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
  `, at = (s) => ({
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
  /**
   * Notification queue
   * @type {Array}
   */
  notifications: [],
  // HTMX Event Handlers
  onContentLoaded(t, n) {
    this._onHtmxEnd(n);
    const e = t.innerHTML || t.textContent || "", i = this._findNotificationById(n.id);
    i && (i.content = e), n.onContentLoad && n.onContentLoad(t), window.Alpine && window.Alpine.initTree && window.Alpine.initTree(t);
  },
  // Helper method to find notification by ID
  _findNotificationById(t) {
    for (const n of Object.keys(this.groups)) {
      const i = this.groups[n].find((o) => o.id === t);
      if (i) return i;
    }
    return null;
  },
  onHtmxError(t, n) {
    this._onHtmxEnd(n), f.error("HTMX Error:", t), this.error("Content load failed"), n.onError && n.onError(t);
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
    const n = h.getStyle("notify", "container"), e = J(t), i = Q(t);
    return [n, e, i.flexDir, i.alignItems].join(" ");
  },
  ensureContainer() {
    const t = "uno-notify-container", n = st(), e = O.ensureContainer(t, n, {
      "x-data": "{}"
      // Independent Alpine component
    });
    e && s.initTree && s.initTree(e);
    const i = "uno-global-progress", o = rt(), r = O.ensureContainer(i, o, {
      "x-data": "{}"
    });
    r && s.initTree && s.initTree(r);
  },
  show(t, n = {}) {
    this.ensureContainer();
    const e = h.getConfig("notify", n), i = this.nextId++;
    let o = t;
    s.store("i18n") && typeof s.store("i18n").t == "function" && (o = s.store("i18n").t(t));
    const r = e.duration !== void 0 ? e.duration : 3e3;
    f.info("[Notify] Show:", { id: i, message: o, duration: r, config: e });
    const l = {
      id: i,
      message: o,
      isVisible: !0,
      // For x-show
      closable: !0,
      duration: r,
      content: e.hxUrl ? "" : void 0,
      // Initialize content for HTMX-loaded content
      ...e
    }, c = e.position || "top-right";
    return this.groups[c] ? this.groups[c].push(l) : (this.groups["top-right"] || (this.groups["top-right"] = []), this.groups["top-right"].push(l)), r > 0 && setTimeout(() => {
      this.remove(i);
    }, r), i;
  },
  remove(t) {
    Object.keys(this.groups).forEach((n) => {
      const e = this.groups[n].findIndex((i) => i.id === t);
      e !== -1 && (this.groups[n][e].visible = !1, setTimeout(() => {
        const i = this.groups[n].findIndex((o) => o.id === t);
        i !== -1 && this.groups[n].splice(i, 1);
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
}), lt = (s) => {
  s.directive("uno-notify", (t, { expression: n }, { evaluateLater: e }) => {
    if (!n) {
      f.warn("[Uno] x-uno-notify directive requires a configuration object or message string.");
      return;
    }
    const i = e(n), o = () => {
      const r = s.store("uno")?.notify;
      r && i((l) => {
        let c = "", a = {};
        if (typeof l == "string")
          c = l;
        else if (typeof l == "object" && l !== null)
          c = l.message || "", a = l;
        else {
          f.warn("[Uno] x-uno-notify directive value must be a string or object", l);
          return;
        }
        r.show(c, a);
      });
    };
    t.addEventListener("click", o);
  });
};
const ct = (s, t = {}) => {
  const n = {
    duration: 3e3,
    position: "top-right"
  };
  h.registerPluginDefaults("notify", {
    ...n,
    ...t.plugins?.notify,
    ...t.notify
  });
  const e = at(s);
  U(s, "notify", e), lt(s);
}, dt = () => `
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

  <div :class="[Uno.getStyle('box', 'dialogWrapper'), 'grid grid-cols-1']"
       @click.self="$store.uno.box.items.length && $store.uno.box.items[$store.uno.box.items.length - 1].closeOnOverlayClick && $store.uno.box.items[$store.uno.box.items.length - 1].close()">
    <template x-for="(item, index) in $store.uno.box.items" :key="item.id">
      <div :id="'uno-box-outer-container' + item.id"
           :class="[Uno.getStyle('box', 'itemWrapper'), $store.uno.box.getPositionClasses(item.position), 'col-start-1 row-start-1']"
           :style="'z-index:' + (${h.options?.ui?.layers?.box || k("BOX")} + index * 10)"
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
             :style="(item.width ? 'width:' + item.width + 'px; max-width:' + item.width + 'px;' : (item.widthClass ? 'max-width:none;' : '')) + (item.height ? ' height:' + item.height + 'px;' : '')"
             x-uno-trap="{ initialFocus: !item.autoFocus }">

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
                    :class="[Uno.getStyle('box', 'okBtn'), item.okButtonClass]"
                    x-show="item.showOkButton"
                    x-text="item.okText"
                    @click="$store.uno.box.handleOk(item)"></button>
            <button type="button"
                    :id="'uno-box-cancel' + item.id"
                    :class="[Uno.getStyle('box', 'cancelBtn'), item.cancelButtonClass]"
                    x-show="item.showCancelButton"
                    x-text="item.cancelText"
                    @click="item.onCancel ? (item.onCancel() !== false && item.close()) : item.close()"></button>
          </div>
        </div>
      </div>
    </template>
  </div>
`, ut = (s) => ({
  items: [],
  /**
   * Dialog stack management
   * @type {Array}
   */
  stack: [],
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
    return new Promise((i) => {
      const o = "prompt-" + Math.random().toString(36).slice(2, 11), r = h.getStyle("box", "promptWrapper"), l = h.getStyle("box", "promptText"), c = h.getStyle("box", "promptInput"), a = h.dom.sanitize(t), u = h.dom.sanitize(n), w = `
        <div class="${r}">
          <p class="${l}">${a}</p>
          <input id="${o}" type="text" class="${c}" value="${u}">
        </div>
      `;
      this.open("confirm", {
        title: e.title || "Prompt",
        content: w,
        _isTrustedPrompt: !0,
        // Internal flag for prompt content
        onOk: () => {
          const g = document.getElementById(o);
          i(g ? g.value : null);
        },
        onCancel: () => i(null),
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
    let e = "dialog", i = n;
    typeof t == "object" ? (i = t, e = i.type || "dialog") : e = t;
    const o = h.getConfig("box", i), r = o.boxId || Math.random().toString(36).slice(2, 10);
    let l = o.showOkButton, c = o.showCancelButton;
    e === "confirm" ? (l = l !== !1, c = c !== !1) : e === "dialog" && (l = l || !1, c = c || !1);
    const a = {
      id: r,
      type: e,
      title: o.title || "",
      content: o._isTrustedPrompt ? o.content || "" : h.dom.sanitize(o.content || ""),
      // Security: Always sanitize content (trustedHTML removed)
      // Exception: Internal prompt content is pre-escaped
      _isTrustedPrompt: o._isTrustedPrompt || !1,
      width: o.width !== void 0 ? o.width : o.widthClass ? void 0 : o.defaultWidth || 400,
      height: o.height !== void 0 ? o.height : (o.heightClass, void 0),
      widthClass: o.widthClass,
      heightClass: o.heightClass,
      containerClass: o.containerClass,
      autoFocus: o.autoFocus !== void 0 ? o.autoFocus : !0,
      position: o.position || o.defaultPosition || "center",
      modal: o.modal !== void 0 ? o.modal : !0,
      closeOnEsc: o.closeOnEsc !== void 0 ? o.closeOnEsc : !1,
      closeOnOverlayClick: o.closeOnOverlayClick !== void 0 ? o.closeOnOverlayClick : !1,
      showCloseButton: o.showCloseButton !== void 0 ? o.showCloseButton : !0,
      showOkButton: l,
      showCancelButton: c,
      okText: o.okText || "OK",
      cancelText: o.cancelText || "Cancel",
      okButtonClass: o.okButtonClass,
      cancelButtonClass: o.cancelButtonClass,
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
        const g = typeof document < "u" ? document.getElementById("uno-box-" + this.id) : null, x = () => {
          const m = s.store("uno")?.box;
          m && m.remove(this.id), this.onClose && this.onClose();
        };
        if (!g) {
          x();
          return;
        }
        let b = !1;
        const y = (m) => {
          m && m.target !== g || b || (b = !0, g.removeEventListener("transitionend", y), x());
        };
        g.addEventListener("transitionend", y), setTimeout(y, 350);
      }
    };
    this.items.push(a), this.ensureContainer();
    const u = this.items[this.items.length - 1], w = () => {
      u.isOpen = !0, this.scheduleAutoFocus(u, s);
    };
    return s && typeof s.nextTick == "function" && s.nextTick(w), u;
  },
  /**
   * Handle OK button click
   * Supports synchronous and asynchronous onOk callbacks
   * @param {object} item - Box item
   */
  async handleOk(t) {
    if (t)
      if (t.onOk) {
        const n = t.onOk();
        n instanceof Promise ? await n !== !1 && t.close() : n !== !1 && t.close();
      } else
        t.close();
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
    return G(t);
  },
  /**
   * Handle HTMX content load completion
   * @param {HTMLElement} el - Content element
   * @param {object} item - Box item
   */
  onHtmxContentLoaded(t, n) {
    n.onContentLoad && typeof n.onContentLoad == "function" && n.onContentLoad(t), window.Alpine && window.Alpine.initTree && window.Alpine.initTree(t), this.scheduleAutoFocus(n, window.Alpine);
  },
  scheduleAutoFocus(t, n) {
    if (!t.autoFocus) return;
    const e = "uno-box-container" + t.id, i = () => {
      const r = document.getElementById(e);
      r && this.focusFirstInput(r);
    }, o = (r) => {
      typeof requestAnimationFrame == "function" ? requestAnimationFrame(() => requestAnimationFrame(r)) : setTimeout(r, 0);
    };
    n && typeof n.nextTick == "function" ? n.nextTick(() => o(i)) : o(i);
  },
  /**
   * Focus the first input element in the container
   * @param {HTMLElement} container - Container element
   * @returns {boolean} True if an element was focused
   */
  focusFirstInput(t) {
    if (!t) return !1;
    const n = [
      'input:not([type="hidden"]):not([type="button"]):not([type="submit"]):not([type="reset"]):not([type="image"]):not([disabled]):not([readonly])',
      "textarea:not([disabled]):not([readonly])",
      "select:not([disabled]):not([readonly])"
    ].join(","), e = t.querySelectorAll(n);
    for (const i of e) {
      const o = window.getComputedStyle(i);
      if (i.offsetParent !== null && o.visibility !== "hidden" && o.display !== "none")
        return i.focus(), !0;
    }
    return !1;
  },
  /**
   * Ensure the box container exists in the DOM
   */
  ensureContainer() {
    const t = "uno-box-manager", n = typeof document < "u" && document.getElementById(t);
    if (n)
      return s && typeof s.initTree == "function" && !n._x_dataStack && s.initTree(n), n;
    const i = {
      "x-data": "{}",
      class: "relative",
      style: `z-index: ${h.options?.ui?.layers?.box || k("BOX")};`,
      "x-show": "$store.uno.box.items.length > 0",
      "x-cloak": ""
      // Hide until Alpine is ready
    }, o = dt(), r = h.dom.ensureContainer(t, o, i);
    return r && s && typeof s.initTree == "function" && s.initTree(r), r;
  }
}), ft = [
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
function F(s, t = {}) {
  const { initialFocus: n = !0 } = t, e = (a) => {
    if (!(a instanceof HTMLElement)) return !1;
    const u = window.getComputedStyle(a);
    return u.display === "none" || u.visibility === "hidden" ? !1 : a.getClientRects().length > 0;
  }, i = () => Array.from(s.querySelectorAll(ft)).filter((a) => e(a)), o = () => i()[0], r = () => {
    const a = i();
    return a[a.length - 1];
  }, l = (a) => {
    if (a.key === "Tab") {
      if (a.shiftKey) {
        if (document.activeElement === o()) {
          a.preventDefault();
          const u = r();
          u && u.focus();
        }
      } else if (document.activeElement === r()) {
        a.preventDefault();
        const u = o();
        u && u.focus();
      }
    }
  }, c = () => {
    if (s.contains(document.activeElement)) return;
    const a = o();
    a && a.focus();
  };
  return n && (typeof requestAnimationFrame == "function" ? requestAnimationFrame(() => requestAnimationFrame(c)) : setTimeout(c, 0)), s.addEventListener("keydown", l), () => {
    s.removeEventListener("keydown", l);
  };
}
const gt = (s) => {
  s.directive("uno-trap", (t, { expression: n }, { cleanup: e, evaluateLater: i }) => {
    if (n)
      i(n)((r) => {
        const l = F(t, r || {});
        e(l);
      });
    else {
      const o = F(t);
      e(o);
    }
  }), s.directive("uno-box", (t, { expression: n }, { evaluateLater: e, cleanup: i }) => {
    if (!n) {
      f.warn("[Uno] x-uno-box directive requires a configuration object.");
      return;
    }
    const o = e(n), r = (l) => {
      l.preventDefault();
      const c = s.store("uno")?.box;
      if (!c) {
        f.warn("[Uno] Box store not found. Is Uno initialized?");
        return;
      }
      o((a) => {
        a && typeof a == "object" ? c.open(a) : f.warn("[Uno] x-uno-box directive value must be an object", a);
      });
    };
    t.addEventListener("click", r), i(() => t.removeEventListener("click", r));
  });
};
const ht = (s, t = {}) => {
  const n = {
    defaultWidth: 400,
    defaultPosition: "center",
    okText: "OK",
    cancelText: "Cancel",
    closeOnEsc: !1,
    closeOnOverlayClick: !1
  };
  h.registerPluginDefaults("box", { ...n, ...t.plugins?.box, ...t.box });
  const e = ut(s);
  U(s, "box", e), gt(s);
}, q = {
  // --- Standard Types ---
  email: /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)+$/,
  url: /^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([/\w .-]*)*\/?$/i,
  date: /^(\d{4}-\d{2}-\d{2}|\d{2}\/\d{2}\/\d{4}|\d{2}\.\d{2}\.\d{4})$/,
  // --- Numeric ---
  int: /^-?\d+$/,
  number: /^-?\d*(\.\d+)?$/,
  // Added floating point support for general 'number'
  // --- String Content ---
  validname: /^[a-zA-Z0-9_-]{3,16}$/,
  en: /^[a-zA-Z]+$/,
  domain: new RegExp("^((?!-)[A-Za-z0-9-]{1,63}(?<!-)\\.)+[A-Za-z]{2,6}$"),
  mobile: /^\+?[\d\s-]{7,}$/,
  // --- Network & Colors (New) ---
  ipv4: /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/,
  hex: /^#?([a-f0-9]{6}|[a-f0-9]{3})$/i,
  // --- Security ---
  pwd: /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d_@#$%^&*()+-=]{8,80}$/,
  // Password must be 8-80 chars, contain letters and numbers, and can include _@#$%^&*()+-= special characters
  strongpwd: /^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,80}$/
  // Password must be 8-80 chars, contain letters, numbers, and at least 1 special char (_@#$%^&*()+-=)
}, mt = (s) => ({
  /**
   * Validation error messages keyed by field name.
   * @type {Object.<string, string>}
   */
  errors: {},
  /**
   * Validation rules keyed by field name.
   * Rules are stored as normalized arrays of rule objects.
   * @type {Object.<string, Array<Object>>}
   */
  rules: {},
  /**
   * Fields currently undergoing async validation.
   * @type {Object.<string, boolean>}
   */
  validating: {},
  /**
   * Custom validators registry.
   * @type {Object.<string, Function>}
   */
  customValidators: {},
  /**
   * Debounce timers for async validation.
   * @type {Object.<string, number>}
   */
  timers: {},
  /**
   * Adds or updates validation rules for a field.
   *
   * @param {string} fieldName - The unique identifier for the field (name or id)
   * @param {Array|Object} config - The validation configuration
   * @param {Object} [instanceOpts={}] - Instance-specific global options
   */
  addRules(t, n, e = {}) {
    const i = this.normalizeRules(n), o = h.getConfig("valid", e);
    this.rules[t] = i.map((r) => ({
      emptyQuit: o.emptyQuit !== !1,
      // Default to true
      ...r
    }));
  },
  /**
   * Normalizes configuration into a standard array of rule objects.
   *
   * @param {Array|Object|string} config
   * @returns {Array<Object>}
   */
  normalizeRules(t) {
    return t ? typeof t == "string" ? [{ type: t }] : typeof t == "object" && !Array.isArray(t) ? [t] : Array.isArray(t) ? t.map((n) => typeof n == "string" ? { type: n } : n) : [] : [];
  },
  /**
   * Registers a custom validator function.
   *
   * @param {string} name - The name of the validator
   * @param {Function} validatorFn - The validator function.
   */
  registerValidator(t, n) {
    typeof n == "function" ? this.customValidators[t] = n : f.error(`[Valid] Invalid validator function for ${t}`);
  },
  /**
   * Validates a single field against its rules.
   *
   * @param {string} fieldName - The field identifier
   * @param {string|number} value - The value to validate
   * @param {string} [eventType='input'] - The event triggering validation ('input', 'blur', 'submit')
   * @returns {Promise<boolean>|boolean} True if valid (can be async)
   */
  async validate(t, n, e = "input") {
    const i = this.rules[t];
    if (!i || i.length === 0) return !0;
    delete this.errors[t], this.timers[t] && (clearTimeout(this.timers[t]), delete this.timers[t]);
    const o = n == null ? "" : String(n), r = o.trim() === "";
    for (const l of i) {
      if (e === "input" && (l.trigger === "blur" || l.realtime === !1))
        continue;
      const c = l.type === "required" || l.required === !0;
      if (!(r && !c && l.emptyQuit !== !1)) {
        if (l.debounce > 0 && e === "input")
          return this.validating[t] = !0, new Promise((a) => {
            this.timers[t] = setTimeout(async () => {
              const u = await this.executeRule(t, n, l, o);
              delete this.validating[t], a(!!u);
            }, l.debounce);
          });
        this.validating[t] = !0;
        try {
          if (!await this.executeRule(t, n, l, o))
            return !1;
        } finally {
          delete this.validating[t];
        }
      }
    }
    return !0;
  },
  /**
   * Executes a single rule.
   *
   * @param {string} fieldName
   * @param {any} value
   * @param {Object} rule
   * @param {string} strValue
   * @returns {Promise<boolean>}
   */
  async executeRule(t, n, e, i) {
    let o = !0, r = "";
    const l = e.type;
    if (l === "required" || e.required)
      i.trim() === "" && (o = !1, r = e.msg || this.getI18nMsg("valid.required") || "This field is required");
    else if (q[l])
      q[l].test(i) || (o = !1, r = e.msg || this.getI18nMsg(`valid.${l}`) || `Invalid ${l} format`);
    else if (l === "length")
      e.min !== void 0 && i.length < e.min ? (o = !1, r = e.msg || this.getI18nMsg("valid.min", { n: e.min }) || `Min length is ${e.min}`) : e.max !== void 0 && i.length > e.max && (o = !1, r = e.msg || this.getI18nMsg("valid.max", { n: e.max }) || `Max length is ${e.max}`);
    else if (l === "range") {
      const c = parseFloat(i);
      isNaN(c) ? (o = !1, r = e.msg || this.getI18nMsg("valid.number") || "Must be a number") : e.min !== void 0 && c < e.min ? (o = !1, r = e.msg || this.getI18nMsg("valid.rangeMin", { n: e.min }) || `Min value is ${e.min}`) : e.max !== void 0 && c > e.max && (o = !1, r = e.msg || this.getI18nMsg("valid.rangeMax", { n: e.max }) || `Max value is ${e.max}`);
    } else if (l === "pattern" || e.pattern)
      (e.pattern instanceof RegExp ? e.pattern : new RegExp(e.pattern)).test(i) || (o = !1, r = e.msg || this.getI18nMsg("valid.pattern") || "Invalid format");
    else if (l === "match" || e.match) {
      const c = e.match || e.selector || e.target;
      if (!c)
        return !0;
      const a = document.querySelector(c.startsWith("#") ? c : `[name="${c}"]`) || document.getElementById(c);
      a && i !== a.value && (o = !1, r = e.msg || this.getI18nMsg("valid.match") || "Fields do not match");
    } else if (l === "custom" || l === "validator" || e.fn || this.customValidators[l]) {
      let c = e.fn || e.name;
      if (typeof c == "string" ? c = this.customValidators[c] : !c && this.customValidators[l] && (c = this.customValidators[l]), typeof c == "function")
        try {
          const a = await c(n, t, e);
          a === !1 ? (o = !1, r = e.msg || "Validation failed") : typeof a == "string" ? (o = !1, r = a) : typeof a == "object" && a.valid === !1 && (o = !1, r = a.message || e.msg || "Validation failed");
        } catch (a) {
          f.error(`[Valid] Error in custom validator for ${t}`, a), o = !1, r = "Validation error";
        }
      else
        f.warn(`[Valid] Custom validator not found: ${l || e.fn}`);
    }
    return o || (this.errors[t] = r), o;
  },
  /**
   * Validates all fields in a form.
   *
   * @param {HTMLFormElement} form - The form element to validate
   * @returns {Promise<boolean>} True if the form is valid
   */
  async validateForm(t) {
    const n = t.querySelectorAll("input, select, textarea");
    let e = !0, i = null;
    for (const o of n) {
      if (o.disabled) continue;
      const r = o.name || o.id;
      this.rules[r] && (await this.validate(r, o.value, "submit") || (e = !1, i || (i = o)));
    }
    return !e && i && (h.getConfig("valid") || {}).scrollOnError !== !1 && (i.scrollIntoView({ behavior: "smooth", block: "center" }), i.focus({ preventScroll: !0 })), e;
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
    const e = s.store("i18n");
    return e && typeof e.t == "function" ? e.t(t, n) : null;
  }
}), pt = (s) => {
  s.directive("uno-valid", (t, { expression: n }, { evaluateLater: e, cleanup: i }) => {
    if (t.tagName === "FORM") return;
    const o = t.name || t.id;
    if (!o) {
      f.warn("[Uno Valid] Element has no name or id:", t);
      return;
    }
    const r = Uno.getConfig("valid"), l = Uno.options.ui?.valid || {};
    if (r.autoError !== !1) {
      const g = `$store.uno.valid.errors['${o}']`, x = `[x-text="${g.replace(/"/g, '\\"')}"]`, b = t.form || t.closest("form");
      if (!(b ? b.querySelector(x) : document.querySelector(x))) {
        let m;
        if (typeof r.createErrorElement == "function")
          m = r.createErrorElement(o, g);
        else {
          m = document.createElement("div"), m.setAttribute("x-text", g), m.setAttribute("x-show", g);
          const v = r.errorClass || l.helperText || "text-xs text-red-500 mt-1";
          m.className = v;
        }
        m && (t.after(m), s.initTree(m));
      }
    }
    const c = n ? e(n) : () => {
    }, a = (g) => {
      s.store("uno").valid.addRules(o, g);
    };
    n ? c((g) => {
      a(g);
      const x = s.store("uno").valid;
      if (x && x.normalizeRules) {
        const y = x.normalizeRules(g).find((m) => m.type === "length" || m.maxLength !== void 0 || m.max !== void 0);
        if (y) {
          const m = y.max !== void 0 ? y.max : y.maxLength;
          if (y.stopOnMax === !0 && m !== void 0 && t.setAttribute("maxlength", m), y.showCounter === !0 && m !== void 0) {
            const v = `${o}-counter`;
            let d = document.getElementById(v);
            const p = t.parentElement;
            p && window.getComputedStyle(p).position === "static" && p.tagName !== "FORM" && p.tagName !== "BODY" && (p.style.position = "relative"), d || (d = document.createElement("span"), d.id = v, d.className = "absolute bottom-1 right-2 text-xs text-gray-400 pointer-events-none transition-colors duration-200", t.after(d));
            const E = () => {
              const B = t.value.length;
              d.textContent = `${B}/${m}`, B > m ? (d.classList.add("text-red-500"), d.classList.remove("text-gray-400")) : (d.classList.remove("text-red-500"), d.classList.add("text-gray-400"));
            };
            E(), t.addEventListener("input", E), i(() => {
              t.removeEventListener("input", E), d && d.parentNode && d.parentNode.removeChild(d);
            });
          }
        }
      }
    }) : a([]);
    const u = () => {
      s.store("uno").valid.validate(o, t.value, "input");
    }, w = () => {
      s.store("uno").valid.validate(o, t.value, "blur");
    };
    t.addEventListener("input", u), t.addEventListener("blur", w), i(() => {
      t.removeEventListener("input", u), t.removeEventListener("blur", w);
    });
  }), typeof document < "u" && (document.body.addEventListener("htmx:confirm", async (t) => {
    const n = t.detail.elt, e = n.tagName === "FORM" ? n : n.closest("form");
    e && e.hasAttribute("x-uno-valid") && (t.preventDefault(), await s.store("uno").valid.validateForm(e) && t.detail.issueRequest(!0));
  }), document.body.addEventListener("htmx:afterRequest", (t) => {
    const n = t.detail.xhr;
    if (n.status === 422)
      try {
        const e = JSON.parse(n.responseText);
        e.errors && (s.store("uno").valid.showErrors(e.errors), f.info("[Valid] Server validation errors applied", e.errors));
      } catch {
      }
  }));
};
const xt = (s, t = {}) => {
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
  const e = mt(s);
  U(s, "valid", e), pt(s);
}, yt = (s) => ({
  // State: { [group]: activeTabName }
  groups: {},
  // Registry: { [tabName]: { group, name, ... } }
  registry: {},
  // Initialization state
  initialized: !1,
  /**
   * Ensure the store is initialized (lazy loading)
   * Sets up hash change listeners only when needed
   */
  ensureInit() {
    if (this.initialized) return;
    this.initialized = !0;
    const t = () => {
      let n = window.location.hash.slice(1);
      if (!n && window.location.href.includes("#") && (n = window.location.href.split("#")[1]), f.info(`[Tab] checkHash: ${n}`), n && this.registry[n]) {
        const { group: e, name: i } = this.registry[n];
        f.info(`[Tab] checkHash match: ${i}`), this.activate({ group: e, name: i, updateHash: !1 });
      }
    };
    t(), window.addEventListener("hashchange", t), window.addEventListener("load", t);
  },
  /**
   * Active tab ID
   * @type {string|null}
   */
  active: null,
  /**
   * Activate a tab
   * @param {object} instanceOpts - Options { name, group, ... }
   */
  activate(t = {}) {
    this.ensureInit();
    const n = h.getConfig("tab", {
      group: "default",
      hash: !0,
      ...t
    });
    n.name && (f.info(
      `[Tab] Activating ${n.name} in ${n.group}. Hash update: ${n.updateHash !== !1}`
    ), this.registry[n.name] || (this.registry[n.name] = { group: n.group, name: n.name }), this.groups[n.group] = n.name, n.hash && n.updateHash !== !1 && window.location.hash.slice(1) !== n.name && window.history.replaceState(null, null, `#${n.name}`), n.updateHash === !1 && setTimeout(() => {
      const e = document.getElementById(`tab-${n.name}`);
      e && e.hasAttribute("hx-get") && (f.info(`[Tab] Triggering HTMX for ${n.name}`), window.htmx && window.htmx.trigger(e, "click"));
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
    this.ensureInit(), f.info(`[Tab] Registering ${t} in ${n}. Hash: ${window.location.hash}`), this.registry[t] = { name: t, group: n }, this.groups[n] || (this.groups[n] = t);
    let e = window.location.hash.slice(1);
    !e && window.location.href.includes("#") && (e = window.location.href.split("#")[1]), e === t && (f.info(`[Tab] Hash match for ${t}. Activating.`), this.activate({ group: n, name: t, updateHash: !1 }));
  }
}), vt = (s) => {
  s.directive("uno-tab", (t, { expression: n }, { evaluateLater: e, cleanup: i }) => {
    if (!n) {
      f.warn("[Uno] x-uno-tab directive requires a configuration object.");
      return;
    }
    const o = e(n);
    o((l) => {
      let c = {};
      if (typeof l == "object" && l !== null)
        c = l;
      else if (typeof l == "string")
        c = { name: l };
      else {
        f.warn("[Uno] x-uno-tab directive value must be a string or object", l);
        return;
      }
      c.group || (c.group = "default"), c.name && s.store("uno").tab.register(c.name, c.group);
    });
    const r = (l) => {
      l.preventDefault(), o((c) => {
        let a = {};
        typeof c == "object" && c !== null ? a = c : typeof c == "string" && (a = { name: c }), a.name && s.store("uno").tab.activate(a);
      });
    };
    t.addEventListener("click", r), i(() => t.removeEventListener("click", r));
  });
};
const wt = (s, t = {}) => {
  const n = {
    defaultTab: "",
    hash: !0,
    activeClass: "is-active",
    inactiveClass: ""
  };
  h.registerPluginDefaults("tab", { ...n, ...t.plugins?.tab, ...t.tab });
  const e = yt();
  U(s, "tab", e), vt(s);
};
h.dom = O;
typeof window < "u" && (window.Uno = h);
h.use(K).use(it).use(ct).use(ht).use(xt).use(wt);
if (typeof window < "u") {
  const s = () => {
    window.Alpine && !h.initialized && h.init(window.Alpine);
  };
  window.Alpine ? s() : document.addEventListener(
    "alpine:init",
    () => {
      s();
    },
    { once: !0 }
  );
}
export {
  O as DomManager,
  h as Uno,
  ht as box,
  h as default,
  ct as notify,
  wt as tab,
  it as tooltip,
  xt as valid
};
