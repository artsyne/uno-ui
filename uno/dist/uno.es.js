import M from "alpinejs";
let $ = !1, O = null;
const d = {
  /**
   * Set debug mode
   * @param {boolean} debug
   */
  setDebug(o) {
    $ = o;
  },
  /**
   * Set custom error reporter
   * @param {Function} reporter - (error, ...args) => void
   */
  setReporter(o) {
    O = o;
  },
  /**
   * Log info message
   * @param {...any} args - Arguments to log
   */
  info(...o) {
    $ && typeof console < "u" && typeof console.log == "function" && console.log("[Uno Info]", ...o);
  },
  /**
   * Log debug message
   * @param {...any} args - Arguments to log
   */
  debug(...o) {
    $ && typeof console < "u" && typeof console.log == "function" && console.log("[Uno Debug]", ...o);
  },
  /**
   * Log warning message
   * @param {...any} args - Arguments to log
   */
  warn(...o) {
    $ && typeof console < "u" && typeof console.warn == "function" && console.warn("[Uno Warn]", ...o);
  },
  /**
   * Log error message
   * @param {...any} args - Arguments to log
   */
  error(...o) {
    if (O)
      try {
        O(...o);
      } catch (t) {
        $ && typeof console < "u" && typeof console.error == "function" && console.error("[Uno Error] Reporter failed:", t);
      }
    $ && typeof console < "u" && typeof console.error == "function" && console.error("[Uno Error]", ...o);
  }
}, P = /* @__PURE__ */ new WeakMap();
function H(o, t) {
  if (!o) return !1;
  const n = String(t);
  let e = P.get(o);
  return e || (e = /* @__PURE__ */ new Set(), P.set(o, e)), e.has(n) ? !1 : (e.add(n), !0);
}
function _(o) {
  return o ? o.replace(/<script\b[^>]*>([\s\S]*?)<\/script>/gim, "").replace(/<img\b[^>]*>/gim, "").replace(/\s+on[a-z]+="[^"]*"/gim, "").replace(/\s+on[a-z]+='[^']*'/gim, "").replace(/\s+href="javascript:[^"]*"/gim, "").replace(/\s+href='javascript:[^']*'/gim, "").replace(/javascript:/gim, "") : "";
}
function V(o, t, { trusted: n = !1, alpine: e } = {}) {
  o && (o.innerHTML = n ? String(t ?? "") : _(String(t ?? "")), e?.initTree?.(o));
}
const B = {
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
  ensureContainer(o, t = "", n = {}, e = document.body) {
    const i = document.getElementById(o);
    if (i) return i;
    try {
      const s = document.createElement("div");
      return s.id = o, window.Uno?.options?.csp?.nonce && !n.nonce && (n.nonce = window.Uno.options.csp.nonce), Object.entries(n).forEach(([a, l]) => {
        s.setAttribute(a, l);
      }), t && (s.innerHTML = t), e ? e.appendChild(s) : document.documentElement.appendChild(s), s;
    } catch (s) {
      return d.error(`DomManager: Failed to inject container #${o}`, s), null;
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
  injectStyles(o, t) {
    if (!document.getElementById(o))
      try {
        const n = document.createElement("style");
        n.id = o, window.Uno?.options?.csp?.nonce && n.setAttribute("nonce", window.Uno.options.csp.nonce), n.textContent = t, document.head.appendChild(n);
      } catch (n) {
        d.error(`DomManager: Failed to inject styles #${o}`, n);
      }
  },
  // Exports
  setHTML: V,
  // Secure HTML setting
  sanitize: _
  // HTML sanitization
};
function L(o) {
  return Object.prototype.toString.call(o) === "[object Object]";
}
function b(o, t) {
  if (!L(o)) return L(t) ? { ...t } : t;
  if (!L(t)) return t === void 0 ? { ...o } : t;
  const n = { ...o };
  for (const [e, i] of Object.entries(t)) {
    if (L(i) && L(n[e])) {
      n[e] = b(n[e], i);
      continue;
    }
    if (L(i)) {
      n[e] = b({}, i);
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
const k = {
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
}, E = (o) => {
  if (!o) return k.BASE;
  const t = o.toUpperCase();
  return k[t] !== void 0 ? k[t] : k.BASE;
};
class W {
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
      ajax: (i, s, a = {}) => typeof window > "u" || typeof window.htmx?.ajax != "function" ? Promise.reject(new Error("htmx.ajax is not available")) : window.htmx.ajax(i, s, a),
      /**
       * Fetch an HTML fragment without inserting it into the DOM.
       * Useful for getting content to display in a modal or tooltip.
       *
       * @param {string} method - HTTP method
       * @param {string} url - The URL to request
       * @param {Object} [ajaxOptions={}] - HTMX ajax options
       * @returns {Promise<string>} The HTML fragment string
       */
      fragment: async (i, s, a = {}) => {
        if (typeof document > "u") throw new Error("document is not available");
        const l = document.createElement("div");
        let r = "";
        const c = (u) => {
          u.preventDefault(), r = u.detail.serverResponse, u.detail.shouldSwap = !1;
        };
        l.addEventListener("htmx:beforeSwap", c);
        try {
          const u = a.swap || "innerHTML";
          return await this.htmx.ajax(i, s, { ...a, target: l, swap: u }), r;
        } catch (u) {
          throw d.error("HTMX fragment fetch failed:", { method: i, url: s, error: u }), u;
        } finally {
          l.removeEventListener("htmx:beforeSwap", c), l.remove();
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
        const s = (c) => c?.detail?.xhr, a = (c, u) => {
          for (const m of this._htmxRequestHandlers) {
            const h = m?.[c];
            if (typeof h == "function")
              try {
                h(u);
              } catch (v) {
                d.error("Uno.htmx.onRequest handler failed", v);
              }
          }
        }, l = (c) => {
          const u = s(c);
          if (u && this._htmxActiveXhrs) {
            if (this._htmxActiveXhrs.has(u)) return;
            this._htmxActiveXhrs.add(u);
          }
          a("start", c);
        }, r = (c) => {
          const u = s(c);
          if (u && this._htmxActiveXhrs) {
            if (!this._htmxActiveXhrs.has(u)) return;
            this._htmxActiveXhrs.delete(u);
          }
          a("end", c);
        };
        document.addEventListener("htmx:beforeRequest", l), document.addEventListener("htmx:afterRequest", r), document.addEventListener("htmx:responseError", r), document.addEventListener("htmx:sendError", r), document.addEventListener("htmx:timeout", r), document.addEventListener("htmx:xhr:progress", (c) => a("progress", c));
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
        layers: { box: k.BOX, tooltip: k.TOOLTIP, notify: k.NOTIFY },
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
    this.options = b(e, t);
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
    return e.plugins && e.plugins[t] && (i = b(i, e.plugins[t])), e[t] && (i = b(i, e[t])), n && (i = b(i, n)), i;
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
    let s = e;
    for (const a of i) {
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
    return !t || typeof t != "object" ? (d.warn("setConfig called with invalid argument, expected object"), this) : (this.initialized ? (d.info("Uno configuration updated after initialization"), this.options = b(this.options, t)) : this._preInitConfig = b(this._preInitConfig, t), this);
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
    return this.initialized && d.warn("Uno.use() called after initialization. Plugin might not work correctly."), this.plugins.push({ plugin: t, options: n }), this;
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
      return d.warn("Uno has already been initialized."), M || (typeof window < "u" ? window.Alpine : null);
    let e = M;
    return !e && typeof window < "u" && window.Alpine && (e = window.Alpine), t && (t.version || typeof t.start == "function" || typeof t.store == "function") && (e = t), this.options = b(this.options, this._preInitConfig), d.setDebug(this.options.debug !== !1), this.options.onError && d.setReporter(this.options.onError), e && typeof e.store == "function" ? (this._doInit(e), e) : typeof document > "u" ? (d.error("Alpine not found and document is unavailable. Cannot initialize Uno."), null) : (d.debug("Alpine not ready, waiting for alpine:init event..."), new Promise((s) => {
      document.addEventListener(
        "alpine:init",
        () => {
          if (e = window.Alpine || M, !e) {
            d.error("Alpine still not found after alpine:init event"), s(null);
            return;
          }
          d.debug("Alpine ready, continuing Uno initialization"), this._doInit(e), s(e);
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
    d.info(`Uno v${this.version} initializing...`), typeof window < "u" && typeof document < "u" && window.htmx && typeof t?.initTree == "function" && H(document, "uno:htmx-afterSwap:initTree") && document.addEventListener("htmx:afterSwap", (n) => {
      const e = n?.detail?.elt;
      e && t.initTree(e);
    }), this.plugins.forEach(({ plugin: n, options: e }) => {
      try {
        const i = b(this.options, e);
        n(t, i);
      } catch (i) {
        d.error("Failed to load plugin", i);
      }
    }), this._setupGlobalErrorHandlers(), this.initialized = !0, d.info("Uno initialized");
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
        d.error("Unhandled Promise Rejection:", e);
      }), window.addEventListener("error", (n) => {
        const e = n.error;
        if (t && e instanceof Error) {
          if (t.has(e)) return;
          t.add(e);
        }
        d.error("Uncaught Error:", {
          message: n.message,
          filename: n.filename,
          lineno: n.lineno,
          colno: n.colno,
          error: e
        });
      }), d.debug("Global error handlers installed");
    }
  }
}
const g = new W(), F = {
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
}, A = {
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
}, X = {
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
const z = {
  en: F,
  "zh-CN": A,
  "zh-TW": X
}, R = Object.keys(z);
function Y() {
  if (typeof localStorage < "u") {
    const o = localStorage.getItem("uno.lang");
    if (o && R.includes(o)) return o;
  }
  if (typeof navigator < "u") {
    const o = navigator.language, t = o.startsWith("zh") ? o.includes("TW") || o.includes("Hant") ? "zh-TW" : "zh-CN" : "en";
    return R.includes(t) ? t : "en";
  }
  return "en";
}
function Z(o, t) {
  const n = t.split(".");
  let e = o;
  for (const i of n) {
    if (!e || typeof e != "object" || !(i in e)) return;
    e = e[i];
  }
  return e;
}
function j(o, t) {
  for (const n in t)
    t[n] && typeof t[n] == "object" ? (o[n] || Object.assign(o, { [n]: {} }), j(o[n], t[n])) : Object.assign(o, { [n]: t[n] });
}
function K(o) {
  const t = o.store("i18n"), n = {
    locale: Y(),
    messages: z,
    loading: !1,
    init() {
      typeof document < "u" && (document.documentElement.lang = this.locale);
    },
    t(e, i = {}) {
      const s = this.locale, a = this.messages[s] || this.messages.en, l = Z(a, e);
      return (typeof l == "string" ? l : e).replace(/\{(\w+)\}/g, (c, u) => u in i ? String(i[u]) : `{${u}}`);
    },
    setLocale(e) {
      if (!this.messages[e]) {
        d.warn(`[Uno] Language ${e} not found. Please load it first using loadLocale().`);
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
          const s = await fetch(i);
          if (!s.ok) throw new Error(`Failed to load ${e}`);
          const a = await s.json();
          this.messages[e] = a, this.setLocale(e);
        } catch (s) {
          d.error(`[Uno] Failed to load locale ${e}:`, s);
        } finally {
          this.loading = !1;
        }
      }
    }
  };
  t ? (d.info("[Uno] Merging into existing i18n store"), t.messages || (t.messages = {}), j(t.messages, z), t.t || (t.t = n.t), t.setLocale || (t.setLocale = n.setLocale), t.setLang || (t.setLang = n.setLang), t.loadLocale || (t.loadLocale = n.loadLocale)) : o.store("i18n", n), o.directive("t", (e, { expression: i, modifiers: s }, { effect: a }) => {
    const l = () => o.store("i18n");
    a(() => {
      const r = l();
      r.locale;
      const c = i, u = r.t(c);
      s.length > 0 ? s.forEach((m) => {
        e.setAttribute(m, u);
      }) : e.textContent = u;
    });
  });
}
function U(o, t, n) {
  if (!o || typeof o.store != "function") {
    d.error("[Uno] Alpine.store is not available");
    return;
  }
  const e = o.store("uno");
  e ? o.store("uno", {
    ...e,
    [t]: n
  }) : o.store("uno", { [t]: n }), o.store(t) || o.store(t, o.store("uno")?.[t]);
}
const G = (o) => {
  const t = o || "center", n = {
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
}, N = (o) => {
  const t = o || "top-right", n = {
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
}, J = (o) => {
  const t = o || "top-right";
  let n = "flex-col", e = "items-end", i = "space-y-2";
  return t.includes("bottom") && (n = "flex-col-reverse", i = "space-y-reverse space-y-2"), t.includes("left") ? e = "items-start" : t.includes("center") ? t === "center" ? e = "items-center justify-center" : e = "items-center" : e = "items-end", { flexDir: n, alignItems: e, spaceY: i };
};
function Q(o, t, n) {
  const e = o.getBoundingClientRect(), i = t.getBoundingClientRect(), s = i.width, a = i.height, l = e.width, r = e.height, c = { top: e.top, left: e.left }, u = window.innerWidth, m = window.innerHeight, h = 0, v = 0, C = m, w = u, x = (T, q) => {
    if (q) {
      if (T + s > w)
        return w - s - 10;
      if (T < v)
        return v + 10;
    } else {
      if (T + a > C)
        return C - a - 10;
      if (T < h)
        return h + 10;
    }
    return T;
  };
  let y, f, p = n;
  switch (n) {
    case "left":
      f = c.left - s - 5, f < v && (f = c.left + l + 5, p = "right"), y = x(c.top + r / 2 - a / 2, !1);
      break;
    case "top":
      y = c.top - a - 8, y < h && (y = c.top + r + 8, p = "bottom"), f = c.left + (l - s) / 2, f = x(f, !0);
      break;
    case "right":
      f = c.left + l + 5, f + s > w && (f = c.left - s - 5, p = "left"), y = x(c.top + r / 2 - a / 2, !1);
      break;
    case "bottom":
    default:
      y = c.top + r + 8, y + a > C && (y = c.top - a - 8, p = "top"), f = c.left + (l - s) / 2, f = x(f, !0);
      break;
  }
  return { top: y, left: f, position: p };
}
const tt = () => (g.options?.ui?.layers?.tooltip || E("TOOLTIP"), `
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
const et = (o) => {
  const t = () => {
    S || typeof requestAnimationFrame == "function" && (S = requestAnimationFrame(() => {
      S = 0;
      const e = o.store("uno")?.tooltip;
      e && e.updatePosition();
    }));
  }, n = () => {
    if (I || typeof window > "u" || typeof document > "u") return;
    const e = () => t(), i = () => t(), s = (l) => {
      if (l.key !== "Escape") return;
      const r = o.store("uno")?.tooltip;
      r?.visible && r.hide();
    }, a = (l) => {
      const r = o.store("uno")?.tooltip;
      if (!r?.visible || r.triggerMode !== "click" || r.clickOutside !== !0) return;
      const c = l.target, u = document.getElementById("uno-tooltip-content"), m = r.triggerEl && r.triggerEl.contains(c), h = u && u.contains(c);
      !m && !h && r.hide();
    };
    window.addEventListener("scroll", e, !0), window.addEventListener("resize", i), window.addEventListener("keydown", s), window.addEventListener("click", a), I = () => {
      window.removeEventListener("scroll", e, !0), window.removeEventListener("resize", i), window.removeEventListener("keydown", s), window.removeEventListener("click", a), I = null;
    };
  };
  return {
    visible: !1,
    content: "",
    placement: "top",
    top: -9999,
    left: -9999,
    zIndex: E("TOOLTIP"),
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
      d.info("Tooltip show called", e);
      const i = g.getConfig("tooltip", e);
      if (d.info("Tooltip resolved config", i), !i.target && typeof i.content == "string" && i.content.startsWith("#") && (i.content.includes(" ") || document.querySelector(i.content) && (i.target = i.content)), this.triggerEl = i.element || this.triggerEl, i.target) {
        const s = document.querySelector(i.target);
        s ? this.content = _(s.innerHTML) : (d.warn(`Tooltip target element not found: ${i.target}`), this.content = i.content || "");
      } else
        this.content = i.content || "";
      this.placement = i.placement, this.triggerMode = i.trigger, this.clickOutside = i.clickOutside, this.offset = i.offset, this.arrow = i.arrow !== !1, this.zIndex = g.options?.ui?.layers?.tooltip || E("TOOLTIP"), this.hxUrl = i.hxUrl, this.hxHeaders = i.hxHeaders, this.hxTrigger = i.hxTrigger || "load", this.visible = !0, this.ensureContainer(), n(), o.nextTick(() => {
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
      const i = Q(this.triggerEl, e, this.placement, {
        offset: this.offset,
        arrow: this.arrow
      });
      this.top = i.top, this.left = i.left, this.placement = i.placement, i.arrow && (this.arrowPosition = i.arrow);
    },
    ensureContainer() {
      const e = "uno-tooltip-manager";
      if (document.getElementById(e)) return;
      const i = g.options?.ui?.layers?.tooltip || E("TOOLTIP"), s = tt();
      B.ensureContainer(e, s, {
        "x-data": "{}",
        style: `position: absolute; top: 0; left: 0; width: 0; height: 0; z-index: ${i};`
      });
      const a = document.getElementById(e);
      a && o && !a._x_dataStack && o.initTree(a);
    }
  };
}, nt = (o) => {
  const t = (n, { expression: e, modifiers: i }, { evaluateLater: s, effect: a, cleanup: l }) => {
    const r = e ? s(e) : null, c = (f) => (i.length > 0 && d.warn(
      "[Uno] x-uno-tooltip modifiers are deprecated in v2.0. Use configuration object instead."
    ), { ...typeof f == "object" && f !== null ? f : f ? { content: f } : {}, element: n }), u = (f) => {
      const p = o.store("uno")?.tooltip;
      p && p.show(c(f));
    }, m = (f) => {
      const p = o.store("uno")?.tooltip;
      p && (p.visible && p.triggerEl === n ? p.hide() : u(f));
    };
    let h = null;
    r && a(() => {
      r((f) => {
        h = f;
      });
    });
    const v = () => {
      if (e && h === null) return;
      const f = c(h);
      g.getConfig("tooltip", f).trigger === "hover" && u(h);
    }, C = () => {
      const f = o.store("uno")?.tooltip;
      f && f.triggerMode === "hover" && f.triggerEl === n && f.scheduleHide();
    }, w = (f) => {
      const p = c(h);
      g.getConfig("tooltip", p).trigger === "click" && (f.preventDefault(), m(h));
    }, x = () => {
      const f = c(h);
      g.getConfig("tooltip", f).trigger === "focus" && u(h);
    }, y = () => {
      const f = o.store("uno")?.tooltip;
      f && f.triggerMode === "focus" && f.triggerEl === n && f.hide();
    };
    n.addEventListener("mouseenter", v), n.addEventListener("mouseleave", C), n.addEventListener("click", w), n.addEventListener("focus", x), n.addEventListener("blur", y), l(() => {
      n.removeEventListener("mouseenter", v), n.removeEventListener("mouseleave", C), n.removeEventListener("click", w), n.removeEventListener("focus", x), n.removeEventListener("blur", y);
      const f = o.store("uno")?.tooltip;
      f && f.triggerEl === n && f.hide();
    });
  };
  o.directive("uno-tooltip", t);
};
const ot = (o, t = {}) => {
  const n = {
    placement: "top",
    trigger: "hover",
    delay: 300,
    offset: 8,
    arrow: !0,
    clickOutside: !0
  };
  g.registerPluginDefaults("tooltip", {
    ...n,
    ...t.plugins?.tooltip,
    ...t.tooltip
  });
  const e = et(o);
  U(o, "tooltip", e), nt(o);
}, it = () => `
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
`, st = () => `
    <div :class="Uno.getStyle('notify', 'progressWrap')"
         style="z-index: ${E("MAX")}"
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
  `, rt = (o) => ({
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
      const i = this.groups[n].find((s) => s.id === t);
      if (i) return i;
    }
    return null;
  },
  onHtmxError(t, n) {
    this._onHtmxEnd(n), d.error("HTMX Error:", t), this.error("Content load failed"), n.onError && n.onError(t);
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
      zIndex: g.options?.ui?.layers?.notify || E("NOTIFY")
    };
  },
  getContainerClass(t) {
    const n = g.getStyle("notify", "container"), e = N(t), i = J(t);
    return [n, e, i.flexDir, i.alignItems].join(" ");
  },
  ensureContainer() {
    const t = "uno-notify-container", n = it(), e = B.ensureContainer(t, n, {
      "x-data": "{}"
      // Independent Alpine component
    });
    e && o.initTree && o.initTree(e);
    const i = "uno-global-progress", s = st(), a = B.ensureContainer(i, s, {
      "x-data": "{}"
    });
    a && o.initTree && o.initTree(a);
  },
  show(t, n = {}) {
    this.ensureContainer();
    const e = g.getConfig("notify", n), i = this.nextId++;
    let s = t;
    o.store("i18n") && typeof o.store("i18n").t == "function" && (s = o.store("i18n").t(t));
    const a = e.duration !== void 0 ? e.duration : 3e3;
    d.info("[Notify] Show:", { id: i, message: s, duration: a, config: e });
    const l = {
      id: i,
      message: s,
      isVisible: !0,
      // For x-show
      closable: !0,
      duration: a,
      content: e.hxUrl ? "" : void 0,
      // Initialize content for HTMX-loaded content
      ...e
    }, r = e.position || "top-right";
    return this.groups[r] ? this.groups[r].push(l) : (this.groups["top-right"] || (this.groups["top-right"] = []), this.groups["top-right"].push(l)), a > 0 && setTimeout(() => {
      this.remove(i);
    }, a), i;
  },
  remove(t) {
    Object.keys(this.groups).forEach((n) => {
      const e = this.groups[n].findIndex((i) => i.id === t);
      e !== -1 && (this.groups[n][e].visible = !1, setTimeout(() => {
        const i = this.groups[n].findIndex((s) => s.id === t);
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
}), at = (o) => {
  o.directive("uno-notify", (t, { expression: n }, { evaluateLater: e }) => {
    if (!n) {
      d.warn("[Uno] x-uno-notify directive requires a configuration object or message string.");
      return;
    }
    const i = e(n), s = () => {
      const a = o.store("uno")?.notify;
      a && i((l) => {
        let r = "", c = {};
        if (typeof l == "string")
          r = l;
        else if (typeof l == "object" && l !== null)
          r = l.message || "", c = l;
        else {
          d.warn("[Uno] x-uno-notify directive value must be a string or object", l);
          return;
        }
        a.show(r, c);
      });
    };
    t.addEventListener("click", s);
  });
};
const lt = (o, t = {}) => {
  const n = {
    duration: 3e3,
    position: "top-right"
  };
  g.registerPluginDefaults("notify", {
    ...n,
    ...t.plugins?.notify,
    ...t.notify
  });
  const e = rt(o);
  U(o, "notify", e), at(o);
}, ct = () => `
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
           :style="'z-index:' + (${g.options?.ui?.layers?.box || E("BOX")} + index * 10)"
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
`, dt = (o) => ({
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
      const s = "prompt-" + Math.random().toString(36).slice(2, 11), a = g.getStyle("box", "promptWrapper"), l = g.getStyle("box", "promptText"), r = g.getStyle("box", "promptInput"), c = g.dom.sanitize(t), u = g.dom.sanitize(n), m = `
        <div class="${a}">
          <p class="${l}">${c}</p>
          <input id="${s}" type="text" class="${r}" value="${u}">
        </div>
      `;
      this.open("confirm", {
        title: e.title || "Prompt",
        content: m,
        _isTrustedPrompt: !0,
        // Internal flag for prompt content
        onOk: () => {
          const h = document.getElementById(s);
          i(h ? h.value : null);
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
    const s = g.getConfig("box", i), a = s.boxId || Math.random().toString(36).slice(2, 10);
    let l = s.showOkButton, r = s.showCancelButton;
    e === "confirm" ? (l = l !== !1, r = r !== !1) : e === "dialog" && (l = l || !1, r = r || !1);
    const c = {
      id: a,
      type: e,
      title: s.title || "",
      content: s._isTrustedPrompt ? s.content || "" : g.dom.sanitize(s.content || ""),
      // Security: Always sanitize content (trustedHTML removed)
      // Exception: Internal prompt content is pre-escaped
      _isTrustedPrompt: s._isTrustedPrompt || !1,
      width: s.width !== void 0 ? s.width : s.widthClass ? void 0 : s.defaultWidth || 400,
      height: s.height !== void 0 ? s.height : (s.heightClass, void 0),
      widthClass: s.widthClass,
      heightClass: s.heightClass,
      containerClass: s.containerClass,
      position: s.position || s.defaultPosition || "center",
      modal: s.modal !== void 0 ? s.modal : !0,
      closeOnEsc: s.closeOnEsc !== void 0 ? s.closeOnEsc : !1,
      closeOnOverlayClick: s.closeOnOverlayClick !== void 0 ? s.closeOnOverlayClick : !1,
      showCloseButton: s.showCloseButton !== void 0 ? s.showCloseButton : !0,
      showOkButton: l,
      showCancelButton: r,
      okText: s.okText || "OK",
      cancelText: s.cancelText || "Cancel",
      okButtonClass: s.okButtonClass,
      cancelButtonClass: s.cancelButtonClass,
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
        const h = typeof document < "u" ? document.getElementById("uno-box-" + this.id) : null, v = () => {
          const x = o.store("uno")?.box;
          x && x.remove(this.id), this.onClose && this.onClose();
        };
        if (!h) {
          v();
          return;
        }
        let C = !1;
        const w = (x) => {
          x && x.target !== h || C || (C = !0, h.removeEventListener("transitionend", w), v());
        };
        h.addEventListener("transitionend", w), setTimeout(w, 350);
      }
    };
    this.items.push(c), this.ensureContainer();
    const u = this.items[this.items.length - 1], m = () => {
      u.isOpen = !0;
    };
    return o && typeof o.nextTick == "function" ? o.nextTick(m) : setTimeout(m, 0), u;
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
    n.onContentLoad && typeof n.onContentLoad == "function" && n.onContentLoad(t), window.Alpine && window.Alpine.initTree && window.Alpine.initTree(t);
  },
  /**
   * Ensure the box container exists in the DOM
   */
  ensureContainer() {
    const t = "uno-box-manager", n = typeof document < "u" && document.getElementById(t);
    if (n)
      return o && typeof o.initTree == "function" && !n._x_dataStack && o.initTree(n), n;
    const i = {
      "x-data": "{}",
      class: "relative",
      style: `z-index: ${g.options?.ui?.layers?.box || E("BOX")};`,
      "x-show": "$store.uno.box.items.length > 0",
      "x-cloak": ""
      // Hide until Alpine is ready
    }, s = ct(), a = g.dom.ensureContainer(t, s, i);
    return a && o && typeof o.initTree == "function" && o.initTree(a), a;
  }
}), ut = [
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
function ft(o, t = {}) {
  const n = (r) => {
    if (!(r instanceof HTMLElement)) return !1;
    const c = window.getComputedStyle(r);
    return c.display === "none" || c.visibility === "hidden" ? !1 : r.getClientRects().length > 0;
  }, e = () => Array.from(o.querySelectorAll(ut)).filter((r) => n(r)), i = () => e()[0], s = () => {
    const r = e();
    return r[r.length - 1];
  }, a = (r) => {
    if (r.key === "Tab") {
      if (r.shiftKey) {
        if (document.activeElement === i()) {
          r.preventDefault();
          const c = s();
          c && c.focus();
        }
      } else if (document.activeElement === s()) {
        r.preventDefault();
        const c = i();
        c && c.focus();
      }
    }
  }, l = () => {
    if (o.contains(document.activeElement)) return;
    const r = i();
    r && r.focus();
  };
  return typeof requestAnimationFrame == "function" ? requestAnimationFrame(() => requestAnimationFrame(l)) : setTimeout(l, 0), o.addEventListener("keydown", a), () => {
    o.removeEventListener("keydown", a);
  };
}
const gt = (o) => {
  o.directive("uno-trap", (t, { expression: n }, { cleanup: e }) => {
    const i = ft(t);
    e(i);
  }), o.directive("uno-box", (t, { expression: n }, { evaluateLater: e, cleanup: i }) => {
    if (!n) {
      d.warn("[Uno] x-uno-box directive requires a configuration object.");
      return;
    }
    const s = e(n), a = (l) => {
      l.preventDefault();
      const r = o.store("uno")?.box;
      if (!r) {
        d.warn("[Uno] Box store not found. Is Uno initialized?");
        return;
      }
      s((c) => {
        c && typeof c == "object" ? r.open(c) : d.warn("[Uno] x-uno-box directive value must be an object", c);
      });
    };
    t.addEventListener("click", a), i(() => t.removeEventListener("click", a));
  });
};
const ht = (o, t = {}) => {
  const n = {
    defaultWidth: 400,
    defaultPosition: "center",
    okText: "OK",
    cancelText: "Cancel",
    closeOnEsc: !1,
    closeOnOverlayClick: !1
  };
  g.registerPluginDefaults("box", { ...n, ...t.plugins?.box, ...t.box });
  const e = dt(o);
  U(o, "box", e), gt(o);
}, D = {
  // --- Standard Types ---
  email: /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)+$/,
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
  // --- Network & Colors (New) ---
  ipv4: /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/,
  hex: /^#?([a-f0-9]{6}|[a-f0-9]{3})$/i,
  // --- Security ---
  pwd: /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$/,
  // Letters + Numbers, 8+
  strongpwd: /^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$/
  // + Special Char
}, pt = (o) => ({
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
    const i = this.normalizeRules(n), s = g.getConfig("valid", e);
    this.rules[t] = i.map((a) => ({
      emptyQuit: s.emptyQuit !== !1,
      // Default to true
      ...a
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
    typeof n == "function" ? this.customValidators[t] = n : d.error(`[Valid] Invalid validator function for ${t}`);
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
    const s = n == null ? "" : String(n), a = s.trim() === "";
    for (const l of i) {
      if (e === "input" && (l.trigger === "blur" || l.realtime === !1))
        continue;
      const r = l.type === "required" || l.required === !0;
      if (!(a && !r && l.emptyQuit !== !1)) {
        if (l.debounce > 0 && e === "input")
          return this.validating[t] = !0, new Promise((c) => {
            this.timers[t] = setTimeout(async () => {
              const u = await this.executeRule(t, n, l, s);
              delete this.validating[t], c(!!u);
            }, l.debounce);
          });
        this.validating[t] = !0;
        try {
          if (!await this.executeRule(t, n, l, s))
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
    let s = !0, a = "";
    const l = e.type;
    if (l === "required" || e.required)
      i.trim() === "" && (s = !1, a = e.msg || this.getI18nMsg("valid.required") || "This field is required");
    else if (D[l])
      D[l].test(i) || (s = !1, a = e.msg || this.getI18nMsg(`valid.${l}`) || `Invalid ${l} format`);
    else if (l === "length")
      e.min !== void 0 && i.length < e.min ? (s = !1, a = e.msg || this.getI18nMsg("valid.min", { n: e.min }) || `Min length is ${e.min}`) : e.max !== void 0 && i.length > e.max && (s = !1, a = e.msg || this.getI18nMsg("valid.max", { n: e.max }) || `Max length is ${e.max}`);
    else if (l === "range") {
      const r = parseFloat(i);
      isNaN(r) ? (s = !1, a = e.msg || this.getI18nMsg("valid.number") || "Must be a number") : e.min !== void 0 && r < e.min ? (s = !1, a = e.msg || this.getI18nMsg("valid.rangeMin", { n: e.min }) || `Min value is ${e.min}`) : e.max !== void 0 && r > e.max && (s = !1, a = e.msg || this.getI18nMsg("valid.rangeMax", { n: e.max }) || `Max value is ${e.max}`);
    } else if (l === "pattern" || e.pattern)
      (e.pattern instanceof RegExp ? e.pattern : new RegExp(e.pattern)).test(i) || (s = !1, a = e.msg || this.getI18nMsg("valid.pattern") || "Invalid format");
    else if (l === "match" || e.match) {
      const r = e.match || e.selector || e.target;
      if (!r)
        return !0;
      const c = document.querySelector(r.startsWith("#") ? r : `[name="${r}"]`) || document.getElementById(r);
      c && i !== c.value && (s = !1, a = e.msg || this.getI18nMsg("valid.match") || "Fields do not match");
    } else if (l === "custom" || l === "validator" || e.fn || this.customValidators[l]) {
      let r = e.fn || e.name;
      if (typeof r == "string" ? r = this.customValidators[r] : !r && this.customValidators[l] && (r = this.customValidators[l]), typeof r == "function")
        try {
          const c = await r(n, t, e);
          c === !1 ? (s = !1, a = e.msg || "Validation failed") : typeof c == "string" ? (s = !1, a = c) : typeof c == "object" && c.valid === !1 && (s = !1, a = c.message || e.msg || "Validation failed");
        } catch (c) {
          d.error(`[Valid] Error in custom validator for ${t}`, c), s = !1, a = "Validation error";
        }
      else
        d.warn(`[Valid] Custom validator not found: ${l || e.fn}`);
    }
    return s || (this.errors[t] = a), s;
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
    for (const s of n) {
      if (s.disabled) continue;
      const a = s.name || s.id;
      this.rules[a] && (await this.validate(a, s.value, "submit") || (e = !1, i || (i = s)));
    }
    return !e && i && (g.getConfig("valid") || {}).scrollOnError !== !1 && (i.scrollIntoView({ behavior: "smooth", block: "center" }), i.focus({ preventScroll: !0 })), e;
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
    const e = o.store("i18n");
    return e && typeof e.t == "function" ? e.t(t, n) : null;
  }
}), mt = (o) => {
  o.directive("uno-valid", (t, { expression: n }, { evaluateLater: e, cleanup: i }) => {
    if (t.tagName === "FORM") return;
    const s = t.name || t.id;
    if (!s) {
      d.warn("[Uno Valid] Element has no name or id:", t);
      return;
    }
    const a = n ? e(n) : () => {
    }, l = (u) => {
      o.store("uno").valid.addRules(s, u);
    };
    n ? a((u) => {
      l(u);
    }) : l([]);
    const r = () => {
      o.store("uno").valid.validate(s, t.value, "input");
    }, c = () => {
      o.store("uno").valid.validate(s, t.value, "blur");
    };
    t.addEventListener("input", r), t.addEventListener("blur", c), i(() => {
      t.removeEventListener("input", r), t.removeEventListener("blur", c);
    });
  }), typeof document < "u" && (document.body.addEventListener("htmx:confirm", async (t) => {
    const n = t.detail.elt, e = n.tagName === "FORM" ? n : n.closest("form");
    e && e.hasAttribute("x-uno-valid") && (t.preventDefault(), await o.store("uno").valid.validateForm(e) && t.detail.issueRequest(!0));
  }), document.body.addEventListener("htmx:afterRequest", (t) => {
    const n = t.detail.xhr;
    if (n.status === 422)
      try {
        const e = JSON.parse(n.responseText);
        e.errors && (o.store("uno").valid.showErrors(e.errors), d.info("[Valid] Server validation errors applied", e.errors));
      } catch {
      }
  }));
};
const xt = (o, t = {}) => {
  const n = {
    realtime: !0,
    showError: !0,
    scrollOnError: !0
  };
  g.registerPluginDefaults("valid", {
    ...n,
    ...t.plugins?.valid,
    ...t.valid
  });
  const e = pt(o);
  U(o, "valid", e), mt(o);
}, yt = (o) => ({
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
      if (!n && window.location.href.includes("#") && (n = window.location.href.split("#")[1]), d.info(`[Tab] checkHash: ${n}`), n && this.registry[n]) {
        const { group: e, name: i } = this.registry[n];
        d.info(`[Tab] checkHash match: ${i}`), this.activate({ group: e, name: i, updateHash: !1 });
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
    const n = g.getConfig("tab", {
      group: "default",
      hash: !0,
      ...t
    });
    n.name && (d.info(
      `[Tab] Activating ${n.name} in ${n.group}. Hash update: ${n.updateHash !== !1}`
    ), this.registry[n.name] || (this.registry[n.name] = { group: n.group, name: n.name }), this.groups[n.group] = n.name, n.hash && n.updateHash !== !1 && window.location.hash.slice(1) !== n.name && window.history.replaceState(null, null, `#${n.name}`), n.updateHash === !1 && setTimeout(() => {
      const e = document.getElementById(`tab-${n.name}`);
      e && e.hasAttribute("hx-get") && (d.info(`[Tab] Triggering HTMX for ${n.name}`), window.htmx && window.htmx.trigger(e, "click"));
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
    this.ensureInit(), d.info(`[Tab] Registering ${t} in ${n}. Hash: ${window.location.hash}`), this.registry[t] = { name: t, group: n }, this.groups[n] || (this.groups[n] = t);
    let e = window.location.hash.slice(1);
    !e && window.location.href.includes("#") && (e = window.location.href.split("#")[1]), e === t && (d.info(`[Tab] Hash match for ${t}. Activating.`), this.activate({ group: n, name: t, updateHash: !1 }));
  }
}), vt = (o) => {
  o.directive("uno-tab", (t, { expression: n }, { evaluateLater: e, cleanup: i }) => {
    if (!n) {
      d.warn("[Uno] x-uno-tab directive requires a configuration object.");
      return;
    }
    const s = e(n);
    s((l) => {
      let r = {};
      if (typeof l == "object" && l !== null)
        r = l;
      else if (typeof l == "string")
        r = { name: l };
      else {
        d.warn("[Uno] x-uno-tab directive value must be a string or object", l);
        return;
      }
      r.group || (r.group = "default"), r.name && o.store("uno").tab.register(r.name, r.group);
    });
    const a = (l) => {
      l.preventDefault(), s((r) => {
        let c = {};
        typeof r == "object" && r !== null ? c = r : typeof r == "string" && (c = { name: r }), c.name && o.store("uno").tab.activate(c);
      });
    };
    t.addEventListener("click", a), i(() => t.removeEventListener("click", a));
  });
};
const wt = (o, t = {}) => {
  const n = {
    defaultTab: "",
    hash: !0,
    activeClass: "is-active",
    inactiveClass: ""
  };
  g.registerPluginDefaults("tab", { ...n, ...t.plugins?.tab, ...t.tab });
  const e = yt();
  U(o, "tab", e), vt(o);
};
g.dom = B;
typeof window < "u" && (window.Uno = g);
g.use(K).use(ot).use(lt).use(ht).use(xt).use(wt);
if (typeof window < "u") {
  const o = () => {
    window.Alpine && !g.initialized && g.init(window.Alpine);
  };
  window.Alpine ? o() : document.addEventListener(
    "alpine:init",
    () => {
      o();
    },
    { once: !0 }
  );
}
export {
  B as DomManager,
  g as Uno,
  ht as box,
  g as default,
  lt as notify,
  wt as tab,
  ot as tooltip,
  xt as valid
};
