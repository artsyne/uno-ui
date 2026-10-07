import D from "alpinejs";
let B = !1, q = null;
const m = {
  /**
   * Set debug mode
   * @param {boolean} debug
   */
  setDebug(r) {
    B = r;
  },
  /**
   * Set custom error reporter
   * @param {Function} reporter - (error, ...args) => void
   */
  setReporter(r) {
    q = r;
  },
  /**
   * Log info message
   * @param {...any} args - Arguments to log
   */
  info(...r) {
    B && typeof console < "u" && typeof console.log == "function" && console.log("[Uno Info]", ...r);
  },
  /**
   * Log debug message
   * @param {...any} args - Arguments to log
   */
  debug(...r) {
    B && typeof console < "u" && typeof console.log == "function" && console.log("[Uno Debug]", ...r);
  },
  /**
   * Log warning message
   * @param {...any} args - Arguments to log
   */
  warn(...r) {
    B && typeof console < "u" && typeof console.warn == "function" && console.warn("[Uno Warn]", ...r);
  },
  /**
   * Log error message
   * @param {...any} args - Arguments to log
   */
  error(...r) {
    if (q)
      try {
        q(...r);
      } catch (t) {
        B && typeof console < "u" && typeof console.error == "function" && console.error("[Uno Error] Reporter failed:", t);
      }
    B && typeof console < "u" && typeof console.error == "function" && console.error("[Uno Error]", ...r);
  }
}, J = /* @__PURE__ */ new WeakMap();
let ot = 0;
function z(r, t) {
  if (!r) return !1;
  const n = String(t), e = `${ot}:${n}`;
  let o = J.get(r);
  return o || (o = /* @__PURE__ */ new Set(), J.set(r, o)), o.has(e) ? !1 : (o.add(e), !0);
}
function A(r) {
  return r ? r.replace(/<script\b[^>]*>([\s\S]*?)<\/script>/gim, "").replace(/<img\b[^>]*>/gim, "").replace(/\s+on[a-z]+="[^"]*"/gim, "").replace(/\s+on[a-z]+='[^']*'/gim, "").replace(/\s+href="javascript:[^"]*"/gim, "").replace(/\s+href='javascript:[^']*'/gim, "").replace(/javascript:/gim, "") : "";
}
function it(r, t, { trusted: n = !1, alpine: e } = {}) {
  r && (r.innerHTML = n ? String(t ?? "") : A(String(t ?? "")), e?.initTree?.(r));
}
const j = {
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
  ensureContainer(r, t = "", n = {}, e = document.body) {
    const o = document.getElementById(r);
    if (o) return o;
    try {
      const i = document.createElement("div");
      return i.id = r, window.Uno?.options?.csp?.nonce && !n.nonce && (n.nonce = window.Uno.options.csp.nonce), Object.entries(n).forEach(([s, a]) => {
        i.setAttribute(s, a);
      }), t && (i.innerHTML = t), e ? e.appendChild(i) : document.documentElement.appendChild(i), i;
    } catch (i) {
      return m.error(`DomManager: Failed to inject container #${r}`, i), null;
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
  injectStyles(r, t) {
    if (!document.getElementById(r))
      try {
        const n = document.createElement("style");
        n.id = r, window.Uno?.options?.csp?.nonce && n.setAttribute("nonce", window.Uno.options.csp.nonce), n.textContent = t, document.head.appendChild(n);
      } catch (n) {
        m.error(`DomManager: Failed to inject styles #${r}`, n);
      }
  },
  // Exports
  setHTML: it,
  // Secure HTML setting
  sanitize: A
  // HTML sanitization
}, rt = (r, t) => {
  try {
    return r.getResponseHeader(t);
  } catch {
    return null;
  }
}, V = (r) => {
  const n = r?.detail?.xhr;
  return n ? {
    status: n.status,
    text: n.responseText ?? "",
    header: (e) => rt(n, e)
  } : null;
}, st = (r) => r?.detail?.xhr;
function U(r) {
  return Object.prototype.toString.call(r) === "[object Object]";
}
function k(r, t) {
  if (!U(r)) return U(t) ? { ...t } : t;
  if (!U(t)) return t === void 0 ? { ...r } : t;
  const n = { ...r };
  for (const [e, o] of Object.entries(t)) {
    if (U(o) && U(n[e])) {
      n[e] = k(n[e], o);
      continue;
    }
    if (U(o)) {
      n[e] = k({}, o);
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
const L = {
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
}, S = (r) => {
  if (!r) return L.BASE;
  const t = r.toUpperCase();
  return L[t] !== void 0 ? L[t] : L.BASE;
};
class at {
  /**
   * Create an UnoInstance.
   *
   * @param {Object} [options={}] - Global configuration options
   */
  constructor(t = {}) {
    this.version = "0.9.0", this.options = t, this.pluginDefaults = {}, this.plugins = [], this.initialized = !1, this._preInitConfig = {}, this._htmxRequestHandlers = [], this._htmxActiveXhrs = typeof WeakSet == "function" ? /* @__PURE__ */ new WeakSet() : null, this.htmx = {
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
      ajax: (o, i, s = {}) => typeof window > "u" || typeof window.htmx?.ajax != "function" ? Promise.reject(new Error("htmx.ajax is not available")) : window.htmx.ajax(o, i, s),
      /**
       * Fetch an HTML fragment without inserting it into the DOM.
       * Useful for getting content to display in a modal or tooltip.
       *
       * @param {string} method - HTTP method
       * @param {string} url - The URL to request
       * @param {Object} [ajaxOptions={}] - HTMX ajax options
       * @returns {Promise<string>} The HTML fragment string
       */
      fragment: async (o, i, s = {}) => {
        if (typeof document > "u") throw new Error("document is not available");
        const a = document.createElement("div");
        let l = "";
        const c = (u) => {
          u.preventDefault(), l = u.detail.serverResponse, u.detail.shouldSwap = !1;
        };
        a.addEventListener("htmx:beforeSwap", c);
        try {
          const u = s.swap || "innerHTML";
          return await this.htmx.ajax(o, i, { ...s, target: a, swap: u }), l;
        } catch (u) {
          throw m.error("HTMX fragment fetch failed:", { method: o, url: i, error: u }), u;
        } finally {
          a.removeEventListener("htmx:beforeSwap", c), a.remove();
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
        if (!o || typeof o != "object" || (this._htmxRequestHandlers.push(o), typeof window > "u" || typeof document > "u") || !window.htmx || !z(document, "uno:htmx:request-hooks")) return;
        const i = st, s = (c, u) => {
          for (const g of this._htmxRequestHandlers) {
            const d = g?.[c];
            if (typeof d == "function")
              try {
                d(u);
              } catch (h) {
                m.error("Uno.htmx.onRequest handler failed", h);
              }
          }
        }, a = (c) => {
          const u = i(c);
          if (u && this._htmxActiveXhrs) {
            if (this._htmxActiveXhrs.has(u)) return;
            this._htmxActiveXhrs.add(u);
          }
          s("start", c);
        }, l = (c) => {
          const u = i(c);
          if (u && this._htmxActiveXhrs) {
            if (!this._htmxActiveXhrs.has(u)) return;
            this._htmxActiveXhrs.delete(u);
          }
          s("end", c);
        };
        document.addEventListener("htmx:beforeRequest", a), document.addEventListener("htmx:afterRequest", l), document.addEventListener("htmx:responseError", l), document.addEventListener("htmx:sendError", l), document.addEventListener("htmx:timeout", l), document.addEventListener("htmx:xhr:progress", (c) => s("progress", c));
      },
      /**
       * Bind notify component to HTMX events.
       * Automatically shows notifications for HTMX requests:
       * - Shows progress bar during active requests
       * - Shows success notification after successful requests
       * - Shows error notification on failed requests
       *
       * @param {Object} [options={}] - Configuration options
       * @param {boolean} [options.showSuccess=true] - Show success notification after successful requests
       * @param {boolean} [options.showError=true] - Show error notification on failed requests
       * @param {boolean} [options.showProgress=true] - Show progress bar during requests
       * @param {string} [options.successMessage='Operation completed'] - Success message template
       * @param {string} [options.errorMessage='Operation failed'] - Error message template
       */
      bindNotifyToHtmx: (o = {}) => {
        const {
          showSuccess: i = !0,
          showError: s = !0,
          showProgress: a = !0,
          successMessage: l = "Operation completed",
          errorMessage: c = "Operation failed"
        } = o, u = this;
        if (typeof document > "u") {
          m.warn("bindNotifyToHtmx: document is not available");
          return;
        }
        if (!u.htmx.available()) {
          m.warn("bindNotifyToHtmx: HTMX is not available");
          return;
        }
        const g = () => typeof window < "u" && window.Alpine ? window.Alpine.store("uno")?.notify : null;
        let d = 0;
        const h = (p) => {
          if (d++, a) {
            const x = g();
            x && typeof x.startProgress == "function" && x.startProgress();
          }
        }, w = (p) => {
          if (d = Math.max(0, d - 1), a && d === 0) {
            const x = g();
            x && typeof x.stopProgress == "function" && x.stopProgress();
          }
          if (i) {
            const x = V(p);
            if (!(x?.status >= 400)) {
              const y = g();
              if (y && typeof y.success == "function") {
                let E = l;
                const C = x?.header("X-Notify-Message");
                C && (E = C), x?.header("X-Notify-Suppress") || y.success(E);
              }
            }
          }
        }, b = (p) => {
          if (d = Math.max(0, d - 1), a && d === 0) {
            const x = g();
            x && typeof x.stopProgress == "function" && x.stopProgress();
          }
          if (s) {
            const x = g();
            if (x && typeof x.error == "function") {
              const f = V(p);
              let y = c;
              const E = f?.header("X-Notify-Error");
              E && (y = E), x.error(y);
            }
          }
        };
        z(document, "uno:htmx:bind-notify") && (document.addEventListener("htmx:beforeRequest", h), document.addEventListener("htmx:afterRequest", w), document.addEventListener("htmx:responseError", b), document.addEventListener("htmx:sendError", b), document.addEventListener("htmx:timeout", b), m.debug("bindNotifyToHtmx: HTMX events bound to notify"));
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
        layers: { box: L.BOX, tooltip: L.TOOLTIP, notify: L.NOTIFY },
        box: {
          overlay: "fixed inset-0 bg-gray-500/75 transition-opacity",
          dialogWrapper: "fixed inset-0 z-10 w-screen overflow-y-auto",
          itemWrapper: "flex min-h-full p-4 text-center sm:p-0 pointer-events-none",
          container: "relative transform overflow-hidden {base} bg-white dark:bg-gray-900 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg flex flex-col max-h-[90vh] pointer-events-auto outline-none",
          cancelBtn: "mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:mt-0 sm:w-auto",
          extraBtn: "mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-normal text-gray-700 hover:bg-gray-50 sm:mt-0 sm:w-auto",
          okBtn: "inline-flex w-full justify-center rounded-md bg-blue-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 sm:ml-3 sm:w-auto",
          title: "text-base font-semibold leading-6 text-gray-900 dark:text-slate-200",
          // 目标尺寸（D5，WCAG 2.5.8）：p-1 使图标命中区达 28×28；-m-1 抵消绝对定位偏移
          // 非文本对比度（E7，WCAG 1.4.11）：亮色 slate-500（hover slate-600）；
          // 暗底配对 slate-400/300
          closeBtnFloating: "absolute top-2 right-2 z-10 p-1 -m-1 text-slate-500 hover:text-slate-600 dark:text-slate-400 dark:hover:text-slate-300 focus:outline-none",
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
          // 目标尺寸（D5，WCAG 2.5.8）：p-1 使图标命中区达 28×28；-m-1 抵消行内占位
          closeBtn: "inline-flex rounded-md p-1 -m-1 text-gray-500 hover:text-gray-700 hover:bg-black/5 dark:text-slate-400 dark:hover:text-slate-200 dark:hover:bg-white/10 focus:outline-none focus:ring-2 focus:ring-indigo-500",
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
          body: "pointer-events-auto",
          // Arrow is an 8px square rotated 45deg, half of it poking out of the
          // tooltip edge. bg-inherit picks up the tooltip background so it reads
          // as a pointer instead of a white patch.
          arrow: "absolute w-2 h-2 rotate-45 bg-inherit border-gray-200 dark:border-gray-700",
          // Only the two sides facing away from the tooltip carry a border, so
          // the outline stays continuous around the tip. Keys are the resolved
          // placement; values assume a clockwise 45deg rotation.
          arrowSides: {
            top: "border-b border-r",
            // tooltip above trigger -> tip points down
            bottom: "border-t border-l",
            // tip points up
            left: "border-t border-r",
            // tip points right
            right: "border-b border-l"
            // tip points left
          }
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
    this.options = k(e, t);
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
    return e.plugins && e.plugins[t] && (o = k(o, e.plugins[t])), e[t] && (o = k(o, e[t])), n && (o = k(o, n)), o;
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
    let i = e;
    for (const s of o) {
      if (i == null) return;
      i = i[s];
    }
    return i;
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
    return !t || typeof t != "object" ? (m.warn("setConfig called with invalid argument, expected object"), this) : (this.initialized ? (m.info("Uno configuration updated after initialization"), this.options = k(this.options, t)) : this._preInitConfig = k(this._preInitConfig, t), this);
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
    return this.initialized && m.warn("Uno.use() called after initialization. Plugin might not work correctly."), this.plugins.push({ plugin: t, options: n }), this;
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
      return m.warn("Uno has already been initialized."), D || (typeof window < "u" ? window.Alpine : null);
    let e = D;
    return !e && typeof window < "u" && window.Alpine && (e = window.Alpine), t && (t.version || typeof t.start == "function" || typeof t.store == "function") && (e = t), this.options = k(this.options, this._preInitConfig), m.setDebug(this.options.debug !== !1), this.options.onError && m.setReporter(this.options.onError), e && typeof e.store == "function" ? (this._doInit(e), e) : typeof document > "u" ? (m.error("Alpine not found and document is unavailable. Cannot initialize Uno."), null) : (m.debug("Alpine not ready, waiting for alpine:init event..."), new Promise((i) => {
      document.addEventListener(
        "alpine:init",
        () => {
          if (e = window.Alpine || D, !e) {
            m.error("Alpine still not found after alpine:init event"), i(null);
            return;
          }
          m.debug("Alpine ready, continuing Uno initialization"), this._doInit(e), i(e);
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
    m.info(`Uno v${this.version} initializing...`), typeof window < "u" && typeof document < "u" && window.htmx && typeof t?.initTree == "function" && z(document, "uno:htmx-afterSwap:initTree") && document.addEventListener("htmx:afterSwap", (n) => {
      const e = n?.detail?.elt;
      e && t.initTree(e);
    }), this.plugins.forEach(({ plugin: n, options: e }) => {
      try {
        const o = k(this.options, e);
        n(t, o);
      } catch (o) {
        m.error("Failed to load plugin", o);
      }
    }), this._setupGlobalErrorHandlers(), this.initialized = !0, m.info("Uno initialized");
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
    if (!(typeof window > "u") && z(window, "uno:global-error-handlers")) {
      const t = typeof WeakSet == "function" ? /* @__PURE__ */ new WeakSet() : null;
      window.addEventListener("unhandledrejection", (n) => {
        const e = n.reason;
        if (t && e instanceof Error) {
          if (t.has(e)) return;
          t.add(e);
        }
        m.error("Unhandled Promise Rejection:", e);
      }), window.addEventListener("error", (n) => {
        const e = n.error;
        if (t && e instanceof Error) {
          if (t.has(e)) return;
          t.add(e);
        }
        m.error("Uncaught Error:", {
          message: n.message,
          filename: n.filename,
          lineno: n.lineno,
          colno: n.colno,
          error: e
        });
      }), m.debug("Global error handlers installed");
    }
  }
}
const v = new at(), lt = {
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
    counter: "{current} of {max} characters entered",
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
    pwd: "Password must be 8-72 chars, contain letters and numbers, and can include _!@#$%^&*()+-=",
    strongpwd: "Password must be 8-80 chars, contain letters, numbers, and at least 1 special char (_@#$%^&*()+-=)",
    summary_title: "Please fix the following problems"
  },
  charCount: "Remaining: {count}"
}, ct = {
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
    counter: "已输入 {current} 字，共 {max} 字",
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
    pwd: "密码需 8-72 位，必须包含字母和数字，可含 _!@#$%^&*()+-= 特殊字符",
    strongpwd: "密码需 8-80 位，必须包含字母、数字和至少 1 个特殊字符 (_@#$%^&*()+-=)",
    summary_title: "请修正以下问题"
  },
  charCount: "剩余: {count}"
}, dt = {
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
    counter: "已輸入 {current} 字，共 {max} 字",
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
    pwd: "密碼需 8-72 位，必須包含字母和數字，可含 _!@#$%^&*()+-= 特殊字符",
    strongpwd: "密碼需 8-80 位，必須包含字母、數字和至少 1 個特殊字符 (_@#$%^&*()+-=)",
    summary_title: "請修正以下問題"
  },
  charCount: "剩餘: {count}"
};
const W = {
  en: lt,
  "zh-CN": ct,
  "zh-TW": dt
}, Q = Object.keys(W);
function ut() {
  if (typeof localStorage < "u") {
    const r = localStorage.getItem("uno.lang");
    if (r && Q.includes(r)) return r;
  }
  if (typeof navigator < "u") {
    const r = navigator.language, t = r.startsWith("zh") ? r.includes("TW") || r.includes("Hant") ? "zh-TW" : "zh-CN" : "en";
    return Q.includes(t) ? t : "en";
  }
  return "en";
}
function ft(r, t) {
  const n = t.split(".");
  let e = r;
  for (const o of n) {
    if (!e || typeof e != "object" || !(o in e)) return;
    e = e[o];
  }
  return e;
}
function N(r, t) {
  for (const n in t)
    t[n] && typeof t[n] == "object" ? (r[n] || Object.assign(r, { [n]: {} }), N(r[n], t[n])) : Object.assign(r, { [n]: t[n] });
}
function gt(r) {
  const t = r.store("i18n"), n = {
    locale: ut(),
    messages: W,
    loading: !1,
    init() {
      typeof document < "u" && (document.documentElement.lang = this.locale);
    },
    t(e, o = {}) {
      const i = this.locale, s = this.messages[i] || this.messages.en, a = ft(s, e);
      return (typeof a == "string" ? a : e).replace(/\{(\w+)\}/g, (c, u) => u in o ? String(o[u]) : `{${u}}`);
    },
    setLocale(e) {
      if (!this.messages[e]) {
        m.warn(`[Uno] Language ${e} not found. Please load it first using loadLocale().`);
        return;
      }
      this.locale = e, typeof localStorage < "u" && localStorage.setItem("uno.lang", e), typeof document < "u" && (document.documentElement.lang = e);
    },
    // Alias for backward compatibility
    setLang(e) {
      this.setLocale(e);
    },
    extend(e) {
      N(this.messages, e);
    },
    async loadLocale(e, o) {
      if (!this.messages[e]) {
        this.loading = !0;
        try {
          const i = await fetch(o);
          if (!i.ok) throw new Error(`Failed to load ${e}`);
          const s = await i.json();
          this.messages[e] = s, this.setLocale(e);
        } catch (i) {
          m.error(`[Uno] Failed to load locale ${e}:`, i);
        } finally {
          this.loading = !1;
        }
      }
    }
  };
  t ? (m.info("[Uno] Merging into existing i18n store"), t.messages || (t.messages = {}), N(t.messages, W), t.t || (t.t = n.t), t.setLocale || (t.setLocale = n.setLocale), t.setLang || (t.setLang = n.setLang), t.loadLocale || (t.loadLocale = n.loadLocale)) : r.store("i18n", n), r.directive("t", (e, { expression: o, modifiers: i }, { effect: s }) => {
    const a = () => r.store("i18n");
    s(() => {
      const l = a();
      l.locale;
      const c = o, u = l.t(c);
      i.length > 0 ? i.forEach((g) => {
        e.setAttribute(g, u);
      }) : e.textContent = u;
    });
  });
}
function H(r, t, n) {
  if (!r || typeof r.store != "function") {
    m.error("[Uno] Alpine.store is not available");
    return;
  }
  const e = r.store("uno");
  e ? r.store("uno", {
    ...e,
    [t]: n
  }) : r.store("uno", { [t]: n }), r.store(t) || r.store(t, r.store("uno")?.[t]);
}
const mt = (r) => {
  const t = r || "center", n = {
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
}, ht = (r) => {
  const t = r || "top-right", n = {
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
}, pt = (r) => {
  const t = r || "top-right";
  let n = "flex-col", e = "items-end", o = "space-y-2";
  return t.includes("bottom") && (n = "flex-col-reverse", o = "space-y-reverse space-y-2"), t.includes("left") ? e = "items-start" : t.includes("center") ? t === "center" ? e = "items-center justify-center" : e = "items-center" : e = "items-end", { flexDir: n, alignItems: e, spaceY: o };
};
function xt(r, t, n, e = {}) {
  const o = typeof e?.offset == "number" && Number.isFinite(e.offset) && e.offset >= 0 ? e.offset : 8, i = r.getBoundingClientRect(), s = t.getBoundingClientRect(), a = s.width, l = s.height, c = i.width, u = i.height, g = { top: i.top, left: i.left }, d = window.innerWidth, h = window.innerHeight, w = 0, b = 0, p = h, x = d, f = ($, R) => {
    if (R) {
      if ($ + a > x)
        return x - a - 10;
      if ($ < b)
        return b + 10;
    } else {
      if ($ + l > p)
        return p - l - 10;
      if ($ < w)
        return w + 10;
    }
    return $;
  };
  let y, E, C = ["top", "bottom", "left", "right"].includes(n) ? n : "bottom";
  switch (n) {
    case "left":
      E = g.left - a - o, E < b && (E = g.left + c + o, C = "right"), y = f(g.top + u / 2 - l / 2, !1);
      break;
    case "top":
      y = g.top - l - o, y < w && (y = g.top + u + o, C = "bottom"), E = g.left + (c - a) / 2, E = f(E, !0);
      break;
    case "right":
      E = g.left + c + o, E + a > x && (E = g.left - a - o, C = "left"), y = f(g.top + u / 2 - l / 2, !1);
      break;
    default:
      y = g.top + u + o, y + l > p && (y = g.top - l - o, C = "top"), E = g.left + (c - a) / 2, E = f(E, !0);
      break;
  }
  const T = 8, I = T / 2, _ = ($, R, nt) => Math.min(Math.max($, R), Math.max(R, nt)), Z = g.left + c / 2, G = g.top + u / 2;
  let P;
  switch (C) {
    case "top":
      P = {
        top: l - I,
        left: _(Z - E - I, 0, a - T)
      };
      break;
    case "left":
      P = {
        top: _(G - y - I, 0, l - T),
        left: a - I
      };
      break;
    case "right":
      P = {
        top: _(G - y - I, 0, l - T),
        left: -I
      };
      break;
    default:
      P = {
        top: -I,
        left: _(Z - E - I, 0, a - T)
      };
      break;
  }
  return { top: y, left: E, position: C, arrow: P };
}
const X = [
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
function K(r, t = {}) {
  const { initialFocus: n = !0 } = t, e = (l) => {
    if (!(l instanceof HTMLElement)) return !1;
    const c = window.getComputedStyle(l);
    return c.display === "none" || c.visibility === "hidden" ? !1 : l.getClientRects().length > 0;
  }, o = () => Array.from(r.querySelectorAll(X)).filter((l) => e(l)), i = () => o()[0], s = (l) => {
    if (l.key !== "Tab") return;
    const c = o();
    if (!c.length) return;
    const u = c[0], g = c[c.length - 1], d = document.activeElement, h = c.some(
      (b) => d.compareDocumentPosition(b) & Node.DOCUMENT_POSITION_FOLLOWING
    ), w = c.some(
      (b) => d.compareDocumentPosition(b) & Node.DOCUMENT_POSITION_PRECEDING
    );
    l.shiftKey ? (d === u || !w) && (l.preventDefault(), g.focus()) : (d === g || !h) && (l.preventDefault(), u.focus());
  }, a = () => {
    if (r.contains(document.activeElement)) return;
    const l = i();
    l && l.focus();
  };
  return n && (typeof requestAnimationFrame == "function" ? requestAnimationFrame(() => requestAnimationFrame(a)) : setTimeout(a, 0)), r.addEventListener("keydown", s), () => {
    r.removeEventListener("keydown", s);
  };
}
const yt = () => (v.options?.ui?.layers?.tooltip || S("TOOLTIP"), `
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
         
         <!-- Arrow: 8px rotated square, anchored by calculatePosition().
              Border sides follow the resolved placement so the tip continues
              the tooltip outline instead of punching a hole in it. -->
         <div x-show="$store.uno.tooltip.arrow"
              :class="[
                Uno.getStyle('tooltip', 'arrow'),
                Uno.getStyle('tooltip', 'arrowSides.' + $store.uno.tooltip.placement, '')
              ]"
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
let M = 0, O = null;
const bt = (r) => {
  const t = () => {
    M || typeof requestAnimationFrame == "function" && (M = requestAnimationFrame(() => {
      M = 0;
      const e = r.store("uno")?.tooltip;
      e && e.updatePosition();
    }));
  }, n = () => {
    if (O || typeof window > "u" || typeof document > "u") return;
    const e = () => t(), o = () => t(), i = (a) => {
      if (a.key !== "Escape") return;
      const l = r.store("uno")?.tooltip;
      l?.visible && l.hide();
    }, s = (a) => {
      const l = r.store("uno")?.tooltip;
      if (!l?.visible || l.triggerMode !== "click" || l.clickOutside !== !0) return;
      const c = a.target, u = document.getElementById("uno-tooltip-content"), g = l.triggerEl && l.triggerEl.contains(c), d = u && u.contains(c);
      !g && !d && l.hide();
    };
    window.addEventListener("scroll", e, !0), window.addEventListener("resize", o), window.addEventListener("keydown", i), window.addEventListener("click", s), O = () => {
      window.removeEventListener("scroll", e, !0), window.removeEventListener("resize", o), window.removeEventListener("keydown", i), window.removeEventListener("click", s), O = null;
    };
  };
  return {
    visible: !1,
    content: "",
    placement: "top",
    top: -9999,
    left: -9999,
    zIndex: S("TOOLTIP"),
    triggerEl: null,
    triggerMode: "hover",
    clickOutside: !0,
    arrow: !0,
    arrowPosition: { top: 0, left: 0 },
    // A3 keyboard bindings (click trigger popup): trap cleanup + keydown handlers
    cleanupTrap: null,
    onPopupKeydown: null,
    onTriggerKeydown: null,
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
      m.info("Tooltip show called", e);
      const o = v.getConfig("tooltip", e);
      if (m.info("Tooltip resolved config", o), !o.target && typeof o.content == "string" && o.content.startsWith("#") && (o.content.includes(" ") || document.querySelector(o.content) && (o.target = o.content)), this.triggerEl = o.element || this.triggerEl, o.target) {
        const i = document.querySelector(o.target);
        i ? this.content = A(i.innerHTML) : (m.warn(`Tooltip target element not found: ${o.target}`), this.content = A(o.content || ""));
      } else
        this.content = A(o.content || "");
      this.placement = o.placement, this.triggerMode = o.trigger, this.clickOutside = o.clickOutside, this.offset = o.offset, this.arrow = o.arrow !== !1, this.zIndex = v.options?.ui?.layers?.tooltip || S("TOOLTIP"), this.hxUrl = o.hxUrl, this.hxHeaders = o.hxHeaders, this.hxTrigger = o.hxTrigger || "load", this.visible = !0, this.ensureContainer(), n(), this.triggerMode === "click" && this.triggerEl && (this.triggerEl.setAttribute("aria-haspopup", "true"), this.triggerEl.setAttribute("aria-expanded", "true"), this.triggerEl.setAttribute("aria-describedby", "uno-tooltip-content"), this.bindPopupKeyboard()), r.nextTick(() => {
        this.updatePosition(), this.focusPopupIfNeeded();
      });
    },
    hide() {
      const e = this.triggerEl, o = this.triggerMode === "click", i = document.getElementById("uno-tooltip-content"), s = document.activeElement, a = !!i && !!s && i.contains(s);
      this.visible = !1, this.triggerEl = null, this.cleanupTrap && (this.cleanupTrap(), this.cleanupTrap = null), i && this.onPopupKeydown && (i.removeEventListener("keydown", this.onPopupKeydown), this.onPopupKeydown = null), e && this.onTriggerKeydown && (e.removeEventListener("keydown", this.onTriggerKeydown), this.onTriggerKeydown = null), o && e && e.setAttribute && (e.setAttribute("aria-expanded", "false"), e.removeAttribute("aria-describedby")), O && (O(), O = null), M && (cancelAnimationFrame(M), M = 0), o && e && typeof e.focus == "function" && (a || s === document.body) && e.focus();
    },
    /**
     * A3 click 弹层键盘绑定：Tab 循环（trapFocus）、role=menu 内容方向键导航、
     * 锚点触发器 Space 合成（Enter 在锚点/按钮上由原生 click 覆盖，按钮 Space 亦原生，
     * 仅锚点需合成以免双重 toggle）。
     */
    bindPopupKeyboard() {
      const e = document.getElementById("uno-tooltip-content");
      e && (this.cleanupTrap = K(e, { initialFocus: !1 }), this.onPopupKeydown = (o) => {
        if (!["ArrowDown", "ArrowUp", "Home", "End"].includes(o.key) || !e.querySelector('[role="menu"]')) return;
        const i = Array.from(e.querySelectorAll(X));
        if (i.length === 0) return;
        o.preventDefault();
        const s = i.indexOf(document.activeElement);
        let a = 0;
        o.key === "ArrowDown" ? a = (s + 1) % i.length : o.key === "ArrowUp" ? a = (s - 1 + i.length) % i.length : o.key === "End" && (a = i.length - 1), i[a] && i[a].focus();
      }, e.addEventListener("keydown", this.onPopupKeydown), this.triggerEl && this.triggerEl.tagName === "A" && (this.onTriggerKeydown = (o) => {
        o.key !== " " && o.key !== "Spacebar" || (o.preventDefault(), this.visible && this.triggerEl === o.currentTarget && this.hide());
      }, this.triggerEl.addEventListener("keydown", this.onTriggerKeydown)));
    },
    /**
     * A3 焦点移入弹层首个可聚焦元素，键盘用户落在刚打开的披露内容内（2.4.3）；
     * 弹层无可聚焦内容时为 no-op（焦点留在触发器）。
     */
    focusPopupIfNeeded() {
      if (this.triggerMode !== "click") return;
      const e = document.getElementById("uno-tooltip-content");
      if (!e || e.contains(document.activeElement)) return;
      const o = e.querySelector(X);
      o && o.focus();
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
      const o = xt(this.triggerEl, e, this.placement, {
        offset: this.offset
      });
      this.top = o.top, this.left = o.left, this.placement = o.position || this.placement, o.arrow && (this.arrowPosition = o.arrow);
    },
    ensureContainer() {
      const e = "uno-tooltip-manager";
      if (document.getElementById(e)) return;
      const o = v.options?.ui?.layers?.tooltip || S("TOOLTIP"), i = yt();
      j.ensureContainer(e, i, {
        "x-data": "{}",
        style: `position: absolute; top: 0; left: 0; width: 0; height: 0; z-index: ${o};`
      });
      const s = document.getElementById(e);
      s && r && !s._x_dataStack && r.initTree(s);
    }
  };
}, vt = (r) => {
  const t = (n, { expression: e, modifiers: o }, { evaluateLater: i, effect: s, cleanup: a }) => {
    const l = e ? i(e) : null, c = (f) => (o.length > 0 && m.warn(
      "[Uno] x-uno-tooltip modifiers are deprecated in v2.0. Use configuration object instead."
    ), { ...typeof f == "object" && f !== null ? f : f ? { content: f } : {}, element: n }), u = (f) => {
      const y = r.store("uno")?.tooltip;
      y && y.show(c(f));
    }, g = (f) => {
      const y = r.store("uno")?.tooltip;
      y && (y.visible && y.triggerEl === n ? y.hide() : u(f));
    };
    let d = null;
    l && s(() => {
      l((f) => {
        d = f;
      });
    });
    const h = () => {
      if (e && d === null) return;
      const f = c(d);
      v.getConfig("tooltip", f).trigger === "hover" && u(d);
    }, w = () => {
      const f = r.store("uno")?.tooltip;
      f && f.triggerMode === "hover" && f.triggerEl === n && f.scheduleHide();
    }, b = (f) => {
      const y = c(d);
      v.getConfig("tooltip", y).trigger === "click" && (f.preventDefault(), g(d));
    }, p = () => {
      const f = c(d);
      v.getConfig("tooltip", f).trigger === "focus" && u(d);
    }, x = () => {
      const f = r.store("uno")?.tooltip;
      f && f.triggerMode === "focus" && f.triggerEl === n && f.hide();
    };
    n.addEventListener("mouseenter", h), n.addEventListener("mouseleave", w), n.addEventListener("click", b), n.addEventListener("focus", p), n.addEventListener("blur", x), a(() => {
      n.removeEventListener("mouseenter", h), n.removeEventListener("mouseleave", w), n.removeEventListener("click", b), n.removeEventListener("focus", p), n.removeEventListener("blur", x);
      const f = r.store("uno")?.tooltip;
      f && f.triggerEl === n && f.hide();
    });
  };
  r.directive("uno-tooltip", t);
};
const wt = (r, t = {}) => {
  const n = {
    placement: "top",
    trigger: "hover",
    delay: 300,
    offset: 8,
    arrow: !0,
    clickOutside: !0
  };
  v.registerPluginDefaults("tooltip", {
    ...n,
    ...t.plugins?.tooltip,
    ...t.tooltip
  });
  const e = bt(r);
  H(r, "tooltip", e), vt(r);
}, Et = () => `
  <template x-for="(items, position) in $store.uno.notify.groups" :key="position">
    <div :class="$store.uno.notify.getContainerClass(position)"
         :style="$store.uno.notify.getContainerStyle(position)">
      
      <template x-for="item in items" :key="item.id">
        <div :class="[
            Uno.getStyle('notify', 'item'),
            Uno.getStyle('notify', 'types.' + item.type)
          ]"
          :role="item.type === 'error' || item.type === 'warning' ? 'alert' : 'status'"
          :aria-live="item.type === 'error' || item.type === 'warning' ? 'assertive' : 'polite'"
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
                      <span class="sr-only" x-t="box.close"></span>
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
`, Ct = () => `
    <div :class="Uno.getStyle('notify', 'progressWrap')"
         style="z-index: ${S("MAX")}"
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
  `, kt = (r) => ({
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
    const e = t.innerHTML || t.textContent || "", o = this._findNotificationById(n.id);
    o && (o.content = e), n.onContentLoad && n.onContentLoad(t), window.Alpine && window.Alpine.initTree && window.Alpine.initTree(t);
  },
  // Helper method to find notification by ID
  _findNotificationById(t) {
    for (const n of Object.keys(this.groups)) {
      const o = this.groups[n].find((i) => i.id === t);
      if (o) return o;
    }
    return null;
  },
  onHtmxError(t, n) {
    this._onHtmxEnd(n), m.error("HTMX Error:", t), this.error("Content load failed"), n.onError && n.onError(t);
  },
  _onHtmxBegin(t) {
    !t || t.__unoHtmxActive || (t.__unoHtmxActive = !0, this.beginProgress());
  },
  _onHtmxEnd(t) {
    !t || !t.__unoHtmxActive || (t.__unoHtmxActive = !1, this.endProgress());
  },
  /**
   * Begin tracking an active request. Shows the progress bar on the first
   * concurrent request. Public API — always pair with endProgress().
   */
  beginProgress() {
    this.progressRequests++, this.progressRequests === 1 && this.startProgress();
  },
  /**
   * End tracking of an active request. Hides the progress bar when the last
   * concurrent request finishes. Safe to call more often than beginProgress().
   */
  endProgress() {
    this.progressRequests > 0 && this.progressRequests--, this.progressRequests === 0 && this.stopProgress();
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
      zIndex: v.options?.ui?.layers?.notify || S("NOTIFY")
    };
  },
  getContainerClass(t) {
    const n = v.getStyle("notify", "container"), e = ht(t), o = pt(t);
    return [n, e, o.flexDir, o.alignItems].join(" ");
  },
  ensureContainer() {
    const t = "uno-notify-container", n = Et(), e = j.ensureContainer(t, n, {
      "x-data": "{}"
      // Independent Alpine component
    });
    e && r.initTree && r.initTree(e);
    const o = "uno-global-progress", i = Ct(), s = j.ensureContainer(o, i, {
      "x-data": "{}"
    });
    s && r.initTree && r.initTree(s);
  },
  show(t, n = {}) {
    this.ensureContainer();
    const e = v.getConfig("notify", n), o = this.nextId++;
    let i = t;
    r.store("i18n") && typeof r.store("i18n").t == "function" && (i = r.store("i18n").t(t));
    const s = e.duration !== void 0 ? e.duration : 3e3;
    m.info("[Notify] Show:", { id: o, message: i, duration: s, config: e });
    const a = {
      id: o,
      message: i,
      isVisible: !0,
      // For x-show
      closable: !0,
      duration: s,
      content: e.hxUrl ? "" : void 0,
      // Initialize content for HTMX-loaded content
      ...e
    }, l = e.position || "top-right";
    return this.groups[l] ? this.groups[l].push(a) : (this.groups["top-right"] || (this.groups["top-right"] = []), this.groups["top-right"].push(a)), s > 0 && setTimeout(() => {
      this.remove(o);
    }, s), o;
  },
  remove(t) {
    Object.keys(this.groups).forEach((n) => {
      const e = this.groups[n].findIndex((o) => o.id === t);
      e !== -1 && (this.groups[n][e].isVisible = !1, setTimeout(() => {
        const o = this.groups[n].findIndex((i) => i.id === t);
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
}), Tt = (r) => {
  r.directive("uno-notify", (t, { expression: n }, { evaluateLater: e }) => {
    if (!n) {
      m.warn("[Uno] x-uno-notify directive requires a configuration object or message string.");
      return;
    }
    const o = e(n), i = () => {
      const s = r.store("uno")?.notify;
      s && o((a) => {
        let l = "", c = {};
        if (typeof a == "string")
          l = a;
        else if (typeof a == "object" && a !== null)
          l = a.message || "", c = a;
        else {
          m.warn("[Uno] x-uno-notify directive value must be a string or object", a);
          return;
        }
        s.show(l, c);
      });
    };
    t.addEventListener("click", i);
  });
};
const It = (r, t = {}) => {
  const n = {
    duration: 3e3,
    position: "top-right"
  };
  v.registerPluginDefaults("notify", {
    ...n,
    ...t.plugins?.notify,
    ...t.notify
  });
  const e = kt(r);
  H(r, "notify", e), Tt(r);
}, St = () => `
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
           :style="'z-index:' + (${v.options?.ui?.layers?.box || S("BOX")} + index * 10)"
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
             :class="[Uno.getStyle('box', 'container'), item.widthClass, item.heightClass, item.containerClass, 'focus:outline-none']"
             :style="(item.width ? 'width:' + item.width + 'px; max-width:' + item.width + 'px;' : (item.widthClass ? 'max-width:none;' : '')) + (item.height ? ' height:' + item.height + 'px;' : '')"
             tabindex="-1"
             x-uno-trap="{ initialFocus: false }">

          <!-- Close Button -->
          <button type="button"
                  x-show="item.showCloseButton"
                  :class="Uno.getStyle('box', 'closeBtnFloating')"
                  @click="item.close()"
                  x-t.aria-label="box.close">
            <svg :class="Uno.getStyle('box', 'closeIconSmall')" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>

          <!-- Content Area -->
          <div :class="Uno.getStyle('box', 'content')">
            <h3 :id="'uno-box-title-' + item.id"
                :class="[Uno.getStyle('box', 'title'), 'focus:outline-none']"
                tabindex="-1"
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
               x-show="item.showOkButton || item.showCancelButton || (item.extraButtons && item.extraButtons.length > 0)">
            <!-- Extra buttons: position='before' (visually LEFT of OK; sm:order-last counters flex-row-reverse) -->
            <template x-for="(btn, xbi) in (item.extraButtons || []).filter(b => (b.position || 'before') === 'before')" :key="'xb-'+xbi">
              <button type="button"
                      :class="[btn.buttonClass || Uno.getStyle('box', 'extraBtn'), 'sm:order-last']"
                      x-text="btn.text"
                      @click="$store.uno.box.handleExtraButton(item, btn)"></button>
            </template>
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
            <!-- Extra buttons: position='after' (rendered after Cancel button) -->
            <template x-for="(btn, xai) in (item.extraButtons || []).filter(b => b.position === 'after')" :key="'xa-'+xai">
              <button type="button"
                      :class="[btn.buttonClass || Uno.getStyle('box', 'extraBtn')]"
                      x-text="btn.text"
                      @click="$store.uno.box.handleExtraButton(item, btn)"></button>
            </template>
          </div>
        </div>
      </div>
    </template>
  </div>
`, $t = (r) => ({
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
   * @param {object} options - Additional options. User-supplied onOk/onCancel/onClose
   *   callbacks are composed with (not replacing) the Promise resolver: the promise
   *   resolves true on OK / false otherwise. Like handleOk, an onOk returning false
   *   (or resolving to false) keeps the dialog open and the promise pending.
   * @returns {Promise<boolean>} Resolves true if OK, false otherwise
   */
  confirm(t, n = {}) {
    return new Promise((e) => {
      const { onOk: o, onCancel: i, onClose: s, ...a } = n;
      this.open("confirm", {
        ...a,
        content: t,
        onOk: () => {
          const l = typeof o == "function" ? o() : void 0;
          if (l === !1) return !1;
          if (l instanceof Promise)
            return l.then((c) => {
              if (c === !1) return !1;
              e(!0);
            });
          e(!0);
        },
        onCancel: () => {
          typeof i == "function" && i(), e(!1);
        },
        onClose: () => {
          typeof s == "function" && s(), e(!1);
        }
      });
    });
  },
  /**
   * Show a prompt dialog
   * @param {string} message - Prompt message
   * @param {string} defaultValue - Default input value
   * @param {object} options - Additional options. User-supplied onOk/onCancel are
   *   composed with the Promise resolver (see confirm()); onOk returning false
   *   keeps the dialog open and the promise pending.
   * @returns {Promise<string|null>} Resolves with input value or null
   */
  prompt(t, n = "", e = {}) {
    return new Promise((o) => {
      const { onOk: i, onCancel: s, ...a } = e, l = "prompt-" + Math.random().toString(36).slice(2, 11), c = v.getStyle("box", "promptWrapper"), u = v.getStyle("box", "promptText"), g = v.getStyle("box", "promptInput"), d = v.dom.sanitize(t), h = v.dom.sanitize(n), w = `
        <div class="${c}">
          <p class="${u}">${d}</p>
          <input id="${l}" type="text" class="${g}" value="${h}">
        </div>
      `;
      this.open("confirm", {
        title: a.title || "Prompt",
        ...a,
        content: w,
        _isTrustedPrompt: !0,
        // Internal flag for prompt content
        onOk: () => {
          const b = typeof i == "function" ? i() : void 0;
          if (b === !1) return !1;
          if (b instanceof Promise)
            return b.then((x) => {
              if (x === !1) return !1;
              const f = document.getElementById(l);
              o(f ? f.value : null);
            });
          const p = document.getElementById(l);
          o(p ? p.value : null);
        },
        onCancel: () => {
          typeof s == "function" && s(), o(null);
        }
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
    const i = v.getConfig("box", o), s = i.boxId || Math.random().toString(36).slice(2, 10);
    let a = i.showOkButton, l = i.showCancelButton;
    e === "confirm" ? (a = a !== !1, l = l !== !1) : e === "dialog" && (a = a || !1, l = l || !1);
    const c = {
      id: s,
      type: e,
      title: i.title || "",
      content: i._isTrustedPrompt ? i.content || "" : v.dom.sanitize(i.content || ""),
      // Security: Always sanitize content (trustedHTML removed)
      // Exception: Internal prompt content is pre-escaped
      _isTrustedPrompt: i._isTrustedPrompt || !1,
      width: i.width !== void 0 ? i.width : i.widthClass ? void 0 : i.defaultWidth || 400,
      height: i.height !== void 0 ? i.height : (i.heightClass, void 0),
      widthClass: i.widthClass,
      heightClass: i.heightClass,
      containerClass: i.containerClass,
      autoFocus: i.autoFocus !== void 0 ? i.autoFocus : !0,
      position: i.position || i.defaultPosition || "center",
      modal: i.modal !== void 0 ? i.modal : !0,
      closeOnEsc: i.closeOnEsc !== void 0 ? i.closeOnEsc : !1,
      closeOnOverlayClick: i.closeOnOverlayClick !== void 0 ? i.closeOnOverlayClick : !1,
      showCloseButton: i.showCloseButton !== void 0 ? i.showCloseButton : !0,
      showOkButton: a,
      showCancelButton: l,
      okText: i.okText || "OK",
      cancelText: i.cancelText || "Cancel",
      okButtonClass: i.okButtonClass,
      cancelButtonClass: i.cancelButtonClass,
      extraButtons: Array.isArray(i.extraButtons) ? i.extraButtons : [],
      onOk: i.onOk,
      onCancel: i.onCancel,
      onBeforeClose: i.onBeforeClose,
      onClose: i.onClose,
      onContentLoad: i.onContentLoad,
      hxUrl: i.hxUrl,
      hxMethod: i.hxMethod || "GET",
      hxHeaders: i.hxHeaders,
      hxTrigger: i.hxTrigger || "load",
      // B5：记录 opener，关闭时焦点还原（2.4.3）
      opener: typeof document < "u" ? document.activeElement : null,
      // Internal state - start closed to trigger animation
      isOpen: !1,
      /**
       * B5 焦点还原：焦点仍在本 box 内（OK/Cancel/关闭按钮）或无处停留（body）时
       * 还原至 opener；焦点已被外部接管（如新弹窗/页面导航）时不抢占。
       * 容器锚点必须与 box.template.js 的 :id 绑定一致（uno-box-outer-container{id}，
       * 即 role="dialog" 元素）——与模板漂移会导致 contains 判定恒为 false。
       */
      restoreFocus() {
        const d = this.opener;
        if (!d || d === document.body || !d.isConnected) return;
        const h = document.getElementById("uno-box-outer-container" + this.id), w = document.activeElement;
        (h && w && h.contains(w) || w === document.body) && d.focus();
      },
      /**
       * Close this box dialog
       */
      close() {
        this.onBeforeClose && this.onBeforeClose(), this.isOpen = !1;
        const d = r.store("uno")?.box;
        d && d.syncUnderlayInert(), this.restoreFocus();
        const h = typeof document < "u" ? document.getElementById("uno-box-outer-container" + this.id) : null, w = () => {
          const x = r.store("uno")?.box;
          x && x.remove(this.id), this.onClose && this.onClose();
        };
        if (!h) {
          w();
          return;
        }
        let b = !1;
        const p = (x) => {
          x && x.target !== h || b || (b = !0, h.removeEventListener("transitionend", p), w());
        };
        h.addEventListener("transitionend", p), setTimeout(p, 350);
      }
    };
    this.items.push(c), this.ensureContainer();
    const u = this.items[this.items.length - 1], g = () => {
      u.isOpen = !0, this.scheduleAutoFocus(u, r), this.syncUnderlayInert();
    };
    return r && typeof r.nextTick == "function" && r.nextTick(g), u;
  },
  /**
   * Handle extra button click
   * Supports synchronous and asynchronous onClick callbacks, consistent with handleOk
   * @param {object} item - Box item
   * @param {object} btn - Extra button config
   */
  async handleExtraButton(t, n) {
    if (!t || !n?.onClick) return;
    const e = n.onClick(t);
    e instanceof Promise ? await e !== !1 && t.close() : e !== !1 && t.close();
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
    n !== -1 && (this.items.splice(n, 1), this.syncUnderlayInert());
  },
  /**
   * B5 多层弹窗：仅顶层可交互/可感知，下层 inert + aria-hidden，
   * 屏幕阅读器与键盘不穿透底层（不改现有 scheduleAutoFocus 契约）。
   */
  syncUnderlayInert() {
    if (typeof document > "u") return;
    const t = this.items.filter((n) => n.isOpen);
    t.forEach((n, e) => {
      const o = document.getElementById("uno-box-outer-container" + n.id);
      if (!o) return;
      const i = e < t.length - 1;
      o.toggleAttribute("inert", i), i ? o.setAttribute("aria-hidden", "true") : o.removeAttribute("aria-hidden");
    });
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
    return mt(t);
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
    const e = "uno-box-container" + t.id, o = () => {
      const s = document.getElementById(e);
      s && (t.autoFocus && this.focusFirstInput(s) || this.focusNeutralTarget(s, t));
    }, i = (s) => {
      typeof requestAnimationFrame == "function" ? requestAnimationFrame(() => requestAnimationFrame(s)) : setTimeout(s, 0);
    };
    n && typeof n.nextTick == "function" ? n.nextTick(() => i(o)) : i(o);
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
    ].join(","), e = t.querySelectorAll(n), o = (i) => {
      for (let s = i; s && s.nodeType === 1; s = s.parentElement) {
        const a = window.getComputedStyle(s);
        if (a.display === "none" || a.visibility === "hidden") return !1;
      }
      return !0;
    };
    for (const i of e)
      if (o(i))
        return i.focus(), !0;
    return !1;
  },
  /**
   * Focus the dialog's neutral target (WAI-ARIA APG dialog pattern):
   * the title heading when present, otherwise the container itself.
   * Keeps initial focus off chrome controls such as the close button.
   * @param {HTMLElement} container - Container element
   * @param {object} item - Box item
   * @returns {boolean} True if a neutral target was focused
   */
  focusNeutralTarget(t, n) {
    if (!t || t.contains(document.activeElement)) return !1;
    const e = n && n.title ? document.getElementById("uno-box-title-" + n.id) : null, o = e && t.contains(e) ? e : t;
    return o.hasAttribute("tabindex") || o.setAttribute("tabindex", "-1"), o.focus({ preventScroll: !0 }), !0;
  },
  /**
   * Ensure the box container exists in the DOM
   */
  ensureContainer() {
    const t = "uno-box-manager", n = typeof document < "u" && document.getElementById(t);
    if (n)
      return r && typeof r.initTree == "function" && !n._x_dataStack && r.initTree(n), n;
    const o = {
      "x-data": "{}",
      class: "relative",
      style: `z-index: ${v.options?.ui?.layers?.box || S("BOX")};`,
      "x-show": "$store.uno.box.items.length > 0",
      "x-cloak": ""
      // Hide until Alpine is ready
    }, i = St(), s = v.dom.ensureContainer(t, i, o);
    return s && r && typeof r.initTree == "function" && r.initTree(s), s;
  }
}), Lt = (r) => {
  r.directive("uno-trap", (t, { expression: n }, { cleanup: e, evaluateLater: o }) => {
    if (n)
      o(n)((s) => {
        const a = K(t, s || {});
        e(a);
      });
    else {
      const i = K(t);
      e(i);
    }
  }), r.directive("uno-box", (t, { expression: n }, { evaluateLater: e, cleanup: o }) => {
    if (!n) {
      m.warn("[Uno] x-uno-box directive requires a configuration object.");
      return;
    }
    const i = e(n), s = (a) => {
      a.preventDefault();
      const l = r.store("uno")?.box;
      if (!l) {
        m.warn("[Uno] Box store not found. Is Uno initialized?");
        return;
      }
      i((c) => {
        c && typeof c == "object" ? l.open(c) : m.warn("[Uno] x-uno-box directive value must be an object", c);
      });
    };
    t.addEventListener("click", s), o(() => t.removeEventListener("click", s));
  });
};
const Bt = (r, t = {}) => {
  const n = {
    defaultWidth: 400,
    defaultPosition: "center",
    okText: "OK",
    cancelText: "Cancel",
    closeOnEsc: !1,
    closeOnOverlayClick: !1
  };
  v.registerPluginDefaults("box", { ...n, ...t.plugins?.box, ...t.box });
  const e = $t(r);
  H(r, "box", e), Lt(r);
}, tt = {
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
  pwd: /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d!_@#$%^&*()+\-=]{8,72}$/,
  // Password must be 8-72 chars, contain letters and numbers, and can include _!@#$%^&*()+-= special characters
  strongpwd: /^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,80}$/
  // Password must be 8-80 chars, contain letters, numbers, and at least 1 special char (_@#$%^&*()+-=)
};
let Ut = 0;
const Mt = (r) => ({
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
    const o = this.normalizeRules(n), i = v.getConfig("valid", e);
    this.rules[t] = o.map((s) => ({
      emptyQuit: i.emptyQuit !== !1,
      // Default to true
      ...s
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
    typeof n == "function" ? this.customValidators[t] = n : m.error(`[Valid] Invalid validator function for ${t}`);
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
    const o = this.rules[t];
    if (!o || o.length === 0) return !0;
    delete this.errors[t], this.syncAriaInvalid(t), this.timers[t] && (clearTimeout(this.timers[t]), delete this.timers[t]);
    const i = n == null ? "" : String(n), s = i.trim() === "";
    for (const a of o) {
      if (e === "input" && (a.trigger === "blur" || a.realtime === !1))
        continue;
      const l = a.type === "required" || a.required === !0;
      if (!(s && !l && a.emptyQuit !== !1)) {
        if (a.debounce > 0 && e === "input")
          return this.validating[t] = !0, new Promise((c) => {
            this.timers[t] = setTimeout(async () => {
              const u = await this.executeRule(t, n, a, i);
              delete this.validating[t], c(!!u);
            }, a.debounce);
          });
        this.validating[t] = !0;
        try {
          if (!await this.executeRule(t, n, a, i))
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
  async executeRule(t, n, e, o) {
    let i = !0, s = "";
    const a = e.type;
    if (a === "required" || e.required)
      o.trim() === "" && (i = !1, s = e.msg || this.getI18nMsg("valid.required") || "This field is required");
    else if (tt[a])
      tt[a].test(o) || (i = !1, s = e.msg || this.getI18nMsg(`valid.${a}`) || `Invalid ${a} format`);
    else if (a === "length")
      e.min !== void 0 && o.length < e.min ? (i = !1, s = e.msg || this.getI18nMsg("valid.min", { n: e.min }) || `Min length is ${e.min}`) : e.max !== void 0 && o.length > e.max && (i = !1, s = e.msg || this.getI18nMsg("valid.max", { n: e.max }) || `Max length is ${e.max}`);
    else if (a === "range") {
      const l = parseFloat(o);
      isNaN(l) ? (i = !1, s = e.msg || this.getI18nMsg("valid.number") || "Must be a number") : e.min !== void 0 && l < e.min ? (i = !1, s = e.msg || this.getI18nMsg("valid.rangeMin", { n: e.min }) || `Min value is ${e.min}`) : e.max !== void 0 && l > e.max && (i = !1, s = e.msg || this.getI18nMsg("valid.rangeMax", { n: e.max }) || `Max value is ${e.max}`);
    } else if (a === "pattern" || e.pattern)
      (e.pattern instanceof RegExp ? e.pattern : new RegExp(e.pattern)).test(o) || (i = !1, s = e.msg || this.getI18nMsg("valid.pattern") || "Invalid format");
    else if (a === "match" || e.match) {
      const l = e.match || e.selector || e.target;
      if (!l)
        return !0;
      const c = document.querySelector(l.startsWith("#") ? l : `[name="${l}"]`) || document.getElementById(l);
      c && o !== c.value && (i = !1, s = e.msg || this.getI18nMsg("valid.match") || "Fields do not match");
    } else if (a === "custom" || a === "validator" || e.fn || this.customValidators[a]) {
      let l = e.fn || e.name;
      if (typeof l == "string" ? l = this.customValidators[l] : !l && this.customValidators[a] && (l = this.customValidators[a]), typeof l == "function")
        try {
          const c = await l(n, t, e);
          c === !1 ? (i = !1, s = e.msg || "Validation failed") : typeof c == "string" ? (i = !1, s = c) : typeof c == "object" && c.valid === !1 && (i = !1, s = c.message || e.msg || "Validation failed");
        } catch (c) {
          m.error(`[Valid] Error in custom validator for ${t}`, c), i = !1, s = "Validation error";
        }
      else
        m.warn(`[Valid] Custom validator not found: ${a || e.fn}`);
    }
    return i || (this.errors[t] = s, this.syncAriaInvalid(t)), i;
  },
  /**
   * Validates all fields in a form.
   *
   * When two or more fields fail, a focusable error summary is rendered at
   * the top of the form (E4 / WCAG 3.3.1 / 3.3.3, GOV.UK pattern) and takes
   * focus; with a single failure the first offending field receives focus as
   * before. Inline messages are always kept in place (never replaced).
   *
   * @param {HTMLFormElement} form - The form element to validate
   * @returns {Promise<boolean>} True if the form is valid
   */
  async validateForm(t) {
    const n = t.querySelectorAll("input, select, textarea");
    let e = !0, o = null;
    const i = [];
    this.removeErrorSummary(t);
    for (const s of n) {
      if (s.disabled) continue;
      const a = s.name || s.id;
      this.rules[a] && (await this.validate(a, s.value, "submit") || (e = !1, o || (o = s), i.push({ input: s, fieldName: a })));
    }
    if (!e && o) {
      const s = v.getConfig("valid") || {};
      if (i.length >= 2) {
        const a = this.renderErrorSummary(t, i);
        a && s.scrollOnError !== !1 && (a.scrollIntoView({ behavior: "smooth", block: "center" }), a.focus({ preventScroll: !0 }));
      } else s.scrollOnError !== !1 && (o.scrollIntoView({ behavior: "smooth", block: "center" }), o.focus({ preventScroll: !0 }));
    }
    return e;
  },
  /**
   * Renders the top-of-form error summary for multi-field failures
   * (E4 / WCAG 3.3.1 / 3.3.3, GOV.UK pattern): a focusable `role="alert"`
   * container labelled by its title, with one link per failed field pointing
   * at the field (`input.id`, falling back to the `{field}-error` element the
   * directive generates). Clicking a link focuses the field instead of
   * navigating. Visual classes (`ui-error-summary*`) are styled by the app
   * stylesheet (website tailwind.css), mirroring the `is-invalid` pattern.
   *
   * @param {HTMLFormElement} form - The form receiving the summary.
   * @param {Array<{input: HTMLElement, fieldName: string}>} failedFields
   * @returns {HTMLElement|null} The summary element, or null when unsupported.
   */
  renderErrorSummary(t, n) {
    if (typeof document > "u" || !t) return null;
    const e = document.createElement("div");
    e.className = "ui-error-summary", e.dataset.validSummary = "true", e.setAttribute("role", "alert"), e.setAttribute("tabindex", "-1");
    const o = `valid-error-summary-title-${++Ut}`;
    e.setAttribute("aria-labelledby", o);
    const i = document.createElement("p");
    i.className = "ui-error-summary-title", i.id = o, i.textContent = this.getI18nMsg("valid.summary_title") || "Please fix the following problems", e.appendChild(i);
    const s = document.createElement("ul");
    return s.className = "ui-error-summary-list", n.forEach(({ input: a, fieldName: l }) => {
      const c = a.id || `${l}-error`, u = document.createElement("a");
      u.className = "ui-error-summary-link", u.setAttribute("href", `#${c}`), u.textContent = this.errors[l] || "", u.addEventListener("click", (d) => {
        d.preventDefault();
        const h = document.getElementById(c);
        h && (h.focus({ preventScroll: !0 }), h.scrollIntoView({ behavior: "smooth", block: "center" }));
      });
      const g = document.createElement("li");
      g.appendChild(u), s.appendChild(g);
    }), e.appendChild(s), t.prepend(e), e;
  },
  /**
   * Removes a previously rendered error summary (matched by the
   * `data-valid-summary` flag written in renderErrorSummary).
   *
   * @param {HTMLFormElement} form - The form to clean.
   */
  removeErrorSummary(t) {
    if (!t || typeof t.querySelector != "function") return;
    const n = t.querySelector("[data-valid-summary]");
    n && n.remove();
  },
  showError(t, n) {
    this.errors[t] = n, this.syncAriaInvalid(t);
  },
  hideError(t) {
    delete this.errors[t], this.syncAriaInvalid(t);
  },
  clearErrors() {
    const t = Object.keys(this.errors);
    this.errors = {}, t.forEach((n) => this.syncAriaInvalid(n));
  },
  showErrors(t) {
    Object.assign(this.errors, t), Object.keys(t || {}).forEach((n) => this.syncAriaInvalid(n));
  },
  /**
   * Resolves the DOM control backing a field name.
   *
   * Field names come from `el.name || el.id`, so look the control up by name
   * first and fall back to id.
   *
   * @param {string} fieldName - The field identifier.
   * @returns {HTMLElement|null} The matching control, or null when detached.
   */
  getFieldElement(t) {
    return typeof document > "u" || !t ? null : document.querySelector(`[name="${t}"]`) || document.getElementById(t);
  },
  /**
   * Mirrors the field's error state onto its `aria-invalid` attribute so screen
   * readers announce the control as invalid (WCAG 3.3.1). Called from every
   * `errors` mutation site so client-side validation and 422 server errors stay
   * in sync with the exposed state.
   *
   * @param {string} fieldName - The field identifier.
   */
  syncAriaInvalid(t) {
    const n = this.getFieldElement(t);
    n && (this.errors[t] ? n.setAttribute("aria-invalid", "true") : n.removeAttribute("aria-invalid"));
  },
  getI18nMsg(t, n) {
    const e = r.store("i18n");
    return e && typeof e.t == "function" ? e.t(t, n) : null;
  }
}), Y = "data-uno-counter-wrapper", Ot = (r) => r === "flex" || r === "inline-flex" || r === "grid" || r === "inline-grid", Pt = (r) => {
  const t = r.parentElement;
  if (!t) return null;
  if (t.hasAttribute(Y)) return t;
  const n = window.getComputedStyle(t);
  if (Ot(n.display))
    return n.position === "static" && t.tagName !== "FORM" && t.tagName !== "BODY" && (t.style.position = "relative"), t;
  const e = r.ownerDocument.createElement("div");
  e.setAttribute(Y, ""), e.style.position = "relative";
  const o = window.getComputedStyle(r).display;
  return (o === "inline" || o === "inline-block") && (e.style.display = "inline-block"), t.insertBefore(e, r), e.appendChild(r), e;
}, At = (r, t) => {
  if (!t || !t.hasAttribute(Y)) return;
  const n = t.parentNode;
  n && (n.insertBefore(r, t), n.removeChild(t));
}, Ht = (r) => `${r}-error`, et = (r, t) => {
  r.id || (r.id = Ht(t)), r.hasAttribute("role") || r.setAttribute("role", "alert");
}, F = (r, t) => {
  if (!t.id) return;
  const n = (r.getAttribute("aria-describedby") || "").split(/\s+/).filter(Boolean);
  n.includes(t.id) || n.push(t.id), r.setAttribute("aria-describedby", n.join(" "));
}, _t = (r) => {
  r.directive("uno-valid", (t, { expression: n }, { evaluateLater: e, cleanup: o }) => {
    if (t.tagName === "FORM") return;
    const i = t.name || t.id;
    if (!i) {
      m.warn("[Uno Valid] Element has no name or id:", t);
      return;
    }
    const s = Uno.getConfig("valid"), a = Uno.options.ui?.valid || {};
    if (s.autoError !== !1) {
      const d = `$store.uno.valid.errors['${i}']`, h = `[x-text="${d.replace(/"/g, '\\"')}"]`, w = t.form || t.closest("form"), b = w ? w.querySelector(h) : document.querySelector(h);
      if (b)
        et(b, i), F(t, b);
      else {
        let p;
        if (typeof s.createErrorElement == "function")
          p = s.createErrorElement(i, d);
        else {
          p = document.createElement("div"), p.setAttribute("x-text", d), p.setAttribute("x-show", d);
          const x = s.errorClass || a.helperText || "text-xs text-red-500 mt-1";
          p.className = x;
        }
        p && (et(p, i), t.after(p), r.initTree(p), F(t, p));
      }
    }
    const l = n ? e(n) : () => {
    }, c = (d) => {
      r.store("uno").valid.addRules(i, d);
    };
    n ? l((d) => {
      c(d);
      const h = r.store("uno").valid;
      if (h && h.normalizeRules) {
        const b = h.normalizeRules(d).find(
          (p) => p.type === "length" || p.maxLength !== void 0 || p.max !== void 0
        );
        if (b) {
          const p = b.max !== void 0 ? b.max : b.maxLength;
          if (b.stopOnMax === !0 && p !== void 0 && t.setAttribute("maxlength", p), b.showCounter === !0 && p !== void 0) {
            const x = `${i}-counter`;
            let f = document.getElementById(x);
            const y = Pt(t);
            f || (f = document.createElement("span"), f.id = x, f.className = "absolute bottom-1 right-2 text-xs text-gray-600 dark:text-gray-400 pointer-events-none transition-colors duration-200"), f.parentElement !== y && t.after(f), F(t, f);
            const E = () => {
              const C = t.value.length;
              f.textContent = `${C}/${p}`;
              const T = h.getI18nMsg ? h.getI18nMsg("valid.counter", { current: C, max: p }) : null;
              f.setAttribute("aria-label", T || `${C}/${p}`), C > p ? (f.classList.add("text-red-600", "dark:text-red-400"), f.classList.remove("text-gray-600", "dark:text-gray-400")) : (f.classList.remove("text-red-600", "dark:text-red-400"), f.classList.add("text-gray-600", "dark:text-gray-400"));
            };
            E(), t.addEventListener("input", E), o(() => {
              t.removeEventListener("input", E), f && f.parentNode && f.parentNode.removeChild(f), At(t, y);
            });
          }
        }
      }
    }) : c([]);
    const u = () => {
      r.store("uno").valid.validate(i, t.value, "input");
    }, g = () => {
      r.store("uno").valid.validate(i, t.value, "blur");
    };
    t.addEventListener("input", u), t.addEventListener("blur", g), o(() => {
      t.removeEventListener("input", u), t.removeEventListener("blur", g);
    });
  }), typeof document < "u" && (document.body.addEventListener("htmx:confirm", async (t) => {
    const n = t.detail.elt, e = n.tagName === "FORM" ? n : n.closest("form");
    if (e && (e.hasAttribute("x-uno-valid") || e.querySelector("[x-uno-valid]"))) {
      t.preventDefault();
      try {
        await r.store("uno").valid.validateForm(e) && t.detail.issueRequest(!0);
      } catch (o) {
        m.error("[Valid] Form validation failed unexpectedly, request blocked", o);
      }
    }
  }), document.body.addEventListener("htmx:afterRequest", (t) => {
    const n = V(t);
    if (n && n.status === 422)
      try {
        const e = JSON.parse(n.text);
        e.errors && (r.store("uno").valid.showErrors(e.errors), m.info("[Valid] Server validation errors applied", e.errors));
      } catch {
      }
  }));
};
const Rt = (r, t = {}) => {
  const n = {
    realtime: !0,
    showError: !0,
    scrollOnError: !0
  };
  v.registerPluginDefaults("valid", {
    ...n,
    ...t.plugins?.valid,
    ...t.valid
  });
  const e = Mt(r);
  H(r, "valid", e), _t(r);
}, zt = (r) => ({
  // State: { [group]: activeTabName }
  groups: {},
  // Registry: { [tabName]: { group, name, ... } }
  registry: {},
  // B4: { [group]: [tabName...] } 组内注册序（DOM 序代理），供方向键旋转
  order: {},
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
      if (!n && window.location.href.includes("#") && (n = window.location.href.split("#")[1]), m.info(`[Tab] checkHash: ${n}`), n && this.registry[n]) {
        const { group: e, name: o } = this.registry[n];
        m.info(`[Tab] checkHash match: ${o}`), this.activate({ group: e, name: o, updateHash: !1 });
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
    const n = v.getConfig("tab", {
      group: "default",
      hash: !0,
      ...t
    });
    n.name && (m.info(
      `[Tab] Activating ${n.name} in ${n.group}. Hash update: ${n.updateHash !== !1}`
    ), this.registry[n.name] || (this.registry[n.name] = { group: n.group, name: n.name }), this.groups[n.group] = n.name, typeof document < "u" && document.querySelectorAll('[role="tabpanel"][data-tab-group="' + n.group + '"]').forEach((e) => e.setAttribute("aria-labelledby", "tab-" + n.name)), n.hash && n.updateHash !== !1 && window.location.hash.slice(1) !== n.name && window.history.replaceState(null, null, `#${n.name}`), n.updateHash === !1 && setTimeout(() => {
      const e = document.getElementById(`tab-${n.name}`);
      e && e.hasAttribute("hx-get") && (m.info(`[Tab] Triggering HTMX for ${n.name}`), window.htmx && window.htmx.trigger(e, "click"));
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
    this.ensureInit(), m.info(`[Tab] Registering ${t} in ${n}. Hash: ${window.location.hash}`), this.registry[t] = { name: t, group: n }, this.order[n] || (this.order[n] = []), this.order[n].includes(t) || this.order[n].push(t), this.groups[n] || (this.groups[n] = t);
    let e = window.location.hash.slice(1);
    !e && window.location.href.includes("#") && (e = window.location.href.split("#")[1]), e === t && (m.info(`[Tab] Hash match for ${t}. Activating.`), this.activate({ group: n, name: t, updateHash: !1 }));
  },
  /**
   * B4 APG 键盘旋转：计算目标 tab 名（不激活），调用方负责 activate + focus
   * （roving tabindex）。←→ 循环旋转，Home/End 跳首尾。
   * @param {string} group
   * @param {string} current
   * @param {string} key - ArrowRight/ArrowLeft/Home/End
   * @returns {string|null}
   */
  rotate(t, n, e) {
    const o = this.order[t];
    if (!o || o.length === 0) return null;
    const i = o.indexOf(n);
    return e === "Home" ? o[0] : e === "End" ? o[o.length - 1] : i === -1 ? o[0] : e === "ArrowRight" ? o[(i + 1) % o.length] : e === "ArrowLeft" ? o[(i - 1 + o.length) % o.length] : null;
  }
}), jt = (r) => {
  r.directive("uno-tab", (t, { expression: n }, { evaluateLater: e, effect: o, cleanup: i }) => {
    if (!n) {
      m.warn("[Uno] x-uno-tab directive requires a configuration object.");
      return;
    }
    const s = e(n);
    let a = null;
    const l = (g) => typeof g == "object" && g !== null ? { ...g } : typeof g == "string" ? { name: g } : null;
    s((g) => {
      const d = l(g);
      if (!d) {
        m.warn("[Uno] x-uno-tab directive value must be a string or object", g);
        return;
      }
      d.group || (d.group = "default"), a = d, d.name && (r.store("uno").tab.register(d.name, d.group), t.setAttribute("role", "tab"), t.id || (t.id = "tab-" + d.name), t.setAttribute("aria-controls", d.panel || "panel-" + d.name));
    }), o(() => {
      if (!a || !a.name) return;
      const g = r.store("uno").tab.isActive(a.name, a.group);
      t.setAttribute("aria-selected", g ? "true" : "false"), t.setAttribute("tabindex", g ? "0" : "-1");
    });
    const c = (g) => {
      g.preventDefault(), s((d) => {
        const h = l(d);
        h && (h.group || (h.group = "default"), a = h, h.name && r.store("uno").tab.activate(h));
      });
    }, u = (g) => {
      if (!["ArrowRight", "ArrowLeft", "Home", "End"].includes(g.key) || !a || !a.name) return;
      g.preventDefault();
      const h = r.store("uno").tab.rotate(a.group, a.name, g.key);
      if (!h) return;
      const w = document.getElementById("tab-" + h);
      w && (w.focus(), w.click());
    };
    t.addEventListener("click", c), t.addEventListener("keydown", u), i(() => {
      t.removeEventListener("click", c), t.removeEventListener("keydown", u);
    });
  });
};
const Dt = (r, t = {}) => {
  const n = {
    defaultTab: "",
    hash: !0,
    activeClass: "is-active",
    inactiveClass: ""
  };
  v.registerPluginDefaults("tab", { ...n, ...t.plugins?.tab, ...t.tab });
  const e = zt();
  H(r, "tab", e), jt(r);
};
v.dom = j;
typeof window < "u" && (window.Uno = v);
v.use(gt).use(wt).use(It).use(Bt).use(Rt).use(Dt);
if (typeof window < "u") {
  const r = () => {
    window.Alpine && !v.initialized && v.init(window.Alpine);
  };
  window.Alpine ? r() : document.addEventListener(
    "alpine:init",
    () => {
      r();
    },
    { once: !0 }
  );
}
export {
  j as DomManager,
  v as Uno,
  Bt as box,
  v as default,
  It as notify,
  Dt as tab,
  wt as tooltip,
  Rt as valid
};
