import D from "alpinejs";
let B = !1, F = null;
const f = {
  /**
   * Set debug mode
   * @param {boolean} debug
   */
  setDebug(i) {
    B = i;
  },
  /**
   * Set custom error reporter
   * @param {Function} reporter - (error, ...args) => void
   */
  setReporter(i) {
    F = i;
  },
  /**
   * Log info message
   * @param {...any} args - Arguments to log
   */
  info(...i) {
    B && typeof console < "u" && typeof console.log == "function" && console.log("[Uno Info]", ...i);
  },
  /**
   * Log debug message
   * @param {...any} args - Arguments to log
   */
  debug(...i) {
    B && typeof console < "u" && typeof console.log == "function" && console.log("[Uno Debug]", ...i);
  },
  /**
   * Log warning message
   * @param {...any} args - Arguments to log
   */
  warn(...i) {
    B && typeof console < "u" && typeof console.warn == "function" && console.warn("[Uno Warn]", ...i);
  },
  /**
   * Log error message
   * @param {...any} args - Arguments to log
   */
  error(...i) {
    if (F)
      try {
        F(...i);
      } catch (t) {
        B && typeof console < "u" && typeof console.error == "function" && console.error("[Uno Error] Reporter failed:", t);
      }
    B && typeof console < "u" && typeof console.error == "function" && console.error("[Uno Error]", ...i);
  }
}, Z = /* @__PURE__ */ new WeakMap();
let tt = 0;
function j(i, t) {
  if (!i) return !1;
  const n = String(t), e = `${tt}:${n}`;
  let o = Z.get(i);
  return o || (o = /* @__PURE__ */ new Set(), Z.set(i, o)), o.has(e) ? !1 : (o.add(e), !0);
}
function P(i) {
  return i ? i.replace(/<script\b[^>]*>([\s\S]*?)<\/script>/gim, "").replace(/<img\b[^>]*>/gim, "").replace(/\s+on[a-z]+="[^"]*"/gim, "").replace(/\s+on[a-z]+='[^']*'/gim, "").replace(/\s+href="javascript:[^"]*"/gim, "").replace(/\s+href='javascript:[^']*'/gim, "").replace(/javascript:/gim, "") : "";
}
function et(i, t, { trusted: n = !1, alpine: e } = {}) {
  i && (i.innerHTML = n ? String(t ?? "") : P(String(t ?? "")), e?.initTree?.(i));
}
const q = {
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
      return s.id = i, window.Uno?.options?.csp?.nonce && !n.nonce && (n.nonce = window.Uno.options.csp.nonce), Object.entries(n).forEach(([r, a]) => {
        s.setAttribute(r, a);
      }), t && (s.innerHTML = t), e ? e.appendChild(s) : document.documentElement.appendChild(s), s;
    } catch (s) {
      return f.error(`DomManager: Failed to inject container #${i}`, s), null;
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
        f.error(`DomManager: Failed to inject styles #${i}`, n);
      }
  },
  // Exports
  setHTML: et,
  // Secure HTML setting
  sanitize: P
  // HTML sanitization
}, nt = (i, t) => {
  try {
    return i.getResponseHeader(t);
  } catch {
    return null;
  }
}, A = (i) => {
  const t = i?.detail, n = t?.xhr;
  if (n)
    return {
      status: n.status,
      text: n.responseText ?? "",
      header: (o) => nt(n, o)
    };
  const e = t?.ctx?.response;
  return e ? {
    status: e.status,
    text: t.ctx.text ?? "",
    header: (o) => {
      try {
        return e.headers?.get(o) ?? null;
      } catch {
        return null;
      }
    }
  } : null;
}, ot = (i) => i?.detail?.xhr;
function U(i) {
  return Object.prototype.toString.call(i) === "[object Object]";
}
function k(i, t) {
  if (!U(i)) return U(t) ? { ...t } : t;
  if (!U(t)) return t === void 0 ? { ...i } : t;
  const n = { ...i };
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
}, $ = (i) => {
  if (!i) return L.BASE;
  const t = i.toUpperCase();
  return L[t] !== void 0 ? L[t] : L.BASE;
};
class it {
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
      ajax: (o, s, r = {}) => typeof window > "u" || typeof window.htmx?.ajax != "function" ? Promise.reject(new Error("htmx.ajax is not available")) : window.htmx.ajax(o, s, r),
      /**
       * Fetch an HTML fragment without inserting it into the DOM.
       * Useful for getting content to display in a modal or tooltip.
       *
       * @param {string} method - HTTP method
       * @param {string} url - The URL to request
       * @param {Object} [ajaxOptions={}] - HTMX ajax options
       * @returns {Promise<string>} The HTML fragment string
       */
      fragment: async (o, s, r = {}) => {
        if (typeof document > "u") throw new Error("document is not available");
        const a = document.createElement("div");
        let l = "";
        const c = (d) => {
          d.preventDefault(), l = d.detail.serverResponse, d.detail.shouldSwap = !1;
        };
        a.addEventListener("htmx:beforeSwap", c);
        try {
          const d = r.swap || "innerHTML";
          return await this.htmx.ajax(o, s, { ...r, target: a, swap: d }), l;
        } catch (d) {
          throw f.error("HTMX fragment fetch failed:", { method: o, url: s, error: d }), d;
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
        if (!o || typeof o != "object" || (this._htmxRequestHandlers.push(o), typeof window > "u" || typeof document > "u") || !window.htmx || !j(document, "uno:htmx:request-hooks")) return;
        const s = ot, r = (c, d) => {
          for (const p of this._htmxRequestHandlers) {
            const g = p?.[c];
            if (typeof g == "function")
              try {
                g(d);
              } catch (w) {
                f.error("Uno.htmx.onRequest handler failed", w);
              }
          }
        }, a = (c) => {
          const d = s(c);
          if (d && this._htmxActiveXhrs) {
            if (this._htmxActiveXhrs.has(d)) return;
            this._htmxActiveXhrs.add(d);
          }
          r("start", c);
        }, l = (c) => {
          const d = s(c);
          if (d && this._htmxActiveXhrs) {
            if (!this._htmxActiveXhrs.has(d)) return;
            this._htmxActiveXhrs.delete(d);
          }
          r("end", c);
        };
        document.addEventListener("htmx:beforeRequest", a), document.addEventListener("htmx:afterRequest", l), document.addEventListener("htmx:responseError", l), document.addEventListener("htmx:sendError", l), document.addEventListener("htmx:timeout", l), document.addEventListener("htmx:xhr:progress", (c) => r("progress", c));
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
          showSuccess: s = !0,
          showError: r = !0,
          showProgress: a = !0,
          successMessage: l = "Operation completed",
          errorMessage: c = "Operation failed"
        } = o, d = this;
        if (typeof document > "u") {
          f.warn("bindNotifyToHtmx: document is not available");
          return;
        }
        if (!d.htmx.available()) {
          f.warn("bindNotifyToHtmx: HTMX is not available");
          return;
        }
        const p = () => typeof window < "u" && window.Alpine ? window.Alpine.store("uno")?.notify : null;
        let g = 0;
        const w = (h) => {
          if (g++, a) {
            const x = p();
            x && typeof x.startProgress == "function" && x.startProgress();
          }
        }, E = (h) => {
          if (g = Math.max(0, g - 1), a && g === 0) {
            const x = p();
            x && typeof x.stopProgress == "function" && x.stopProgress();
          }
          if (s) {
            const x = A(h);
            if (!(x?.status >= 400)) {
              const m = p();
              if (m && typeof m.success == "function") {
                let v = l;
                const C = x?.header("X-Notify-Message");
                C && (v = C), x?.header("X-Notify-Suppress") || m.success(v);
              }
            }
          }
        }, b = (h) => {
          if (g = Math.max(0, g - 1), a && g === 0) {
            const x = p();
            x && typeof x.stopProgress == "function" && x.stopProgress();
          }
          if (r) {
            const x = p();
            if (x && typeof x.error == "function") {
              const u = A(h);
              let m = c;
              const v = u?.header("X-Notify-Error");
              v && (m = v), x.error(m);
            }
          }
        };
        j(document, "uno:htmx:bind-notify") && (document.addEventListener("htmx:beforeRequest", w), document.addEventListener("htmx:afterRequest", E), document.addEventListener("htmx:responseError", b), document.addEventListener("htmx:sendError", b), document.addEventListener("htmx:timeout", b), f.debug("bindNotifyToHtmx: HTMX events bound to notify"));
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
          closeBtn: "inline-flex rounded-md text-gray-500 hover:text-gray-700 hover:bg-black/5 dark:text-slate-400 dark:hover:text-slate-200 dark:hover:bg-white/10 focus:outline-none focus:ring-2 focus:ring-indigo-500",
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
    let s = e;
    for (const r of o) {
      if (s == null) return;
      s = s[r];
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
    return !t || typeof t != "object" ? (f.warn("setConfig called with invalid argument, expected object"), this) : (this.initialized ? (f.info("Uno configuration updated after initialization"), this.options = k(this.options, t)) : this._preInitConfig = k(this._preInitConfig, t), this);
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
      return f.warn("Uno has already been initialized."), D || (typeof window < "u" ? window.Alpine : null);
    let e = D;
    return !e && typeof window < "u" && window.Alpine && (e = window.Alpine), t && (t.version || typeof t.start == "function" || typeof t.store == "function") && (e = t), this.options = k(this.options, this._preInitConfig), f.setDebug(this.options.debug !== !1), this.options.onError && f.setReporter(this.options.onError), e && typeof e.store == "function" ? (this._doInit(e), e) : typeof document > "u" ? (f.error("Alpine not found and document is unavailable. Cannot initialize Uno."), null) : (f.debug("Alpine not ready, waiting for alpine:init event..."), new Promise((s) => {
      document.addEventListener(
        "alpine:init",
        () => {
          if (e = window.Alpine || D, !e) {
            f.error("Alpine still not found after alpine:init event"), s(null);
            return;
          }
          f.debug("Alpine ready, continuing Uno initialization"), this._doInit(e), s(e);
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
    f.info(`Uno v${this.version} initializing...`), typeof window < "u" && typeof document < "u" && window.htmx && typeof t?.initTree == "function" && j(document, "uno:htmx-afterSwap:initTree") && document.addEventListener("htmx:afterSwap", (n) => {
      const e = n?.detail?.elt;
      e && t.initTree(e);
    }), this.plugins.forEach(({ plugin: n, options: e }) => {
      try {
        const o = k(this.options, e);
        n(t, o);
      } catch (o) {
        f.error("Failed to load plugin", o);
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
    if (!(typeof window > "u") && j(window, "uno:global-error-handlers")) {
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
const y = new it(), st = {
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
    pwd: "Password must be 8-72 chars, contain letters and numbers, and can include _!@#$%^&*()+-=",
    strongpwd: "Password must be 8-80 chars, contain letters, numbers, and at least 1 special char (_@#$%^&*()+-=)"
  },
  charCount: "Remaining: {count}"
}, rt = {
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
    pwd: "密码需 8-72 位，必须包含字母和数字，可含 _!@#$%^&*()+-= 特殊字符",
    strongpwd: "密码需 8-80 位，必须包含字母、数字和至少 1 个特殊字符 (_@#$%^&*()+-=)"
  },
  charCount: "剩余: {count}"
}, at = {
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
    pwd: "密碼需 8-72 位，必須包含字母和數字，可含 _!@#$%^&*()+-= 特殊字符",
    strongpwd: "密碼需 8-80 位，必須包含字母、數字和至少 1 個特殊字符 (_@#$%^&*()+-=)"
  },
  charCount: "剩餘: {count}"
};
const V = {
  en: st,
  "zh-CN": rt,
  "zh-TW": at
}, K = Object.keys(V);
function lt() {
  if (typeof localStorage < "u") {
    const i = localStorage.getItem("uno.lang");
    if (i && K.includes(i)) return i;
  }
  if (typeof navigator < "u") {
    const i = navigator.language, t = i.startsWith("zh") ? i.includes("TW") || i.includes("Hant") ? "zh-TW" : "zh-CN" : "en";
    return K.includes(t) ? t : "en";
  }
  return "en";
}
function ct(i, t) {
  const n = t.split(".");
  let e = i;
  for (const o of n) {
    if (!e || typeof e != "object" || !(o in e)) return;
    e = e[o];
  }
  return e;
}
function W(i, t) {
  for (const n in t)
    t[n] && typeof t[n] == "object" ? (i[n] || Object.assign(i, { [n]: {} }), W(i[n], t[n])) : Object.assign(i, { [n]: t[n] });
}
function dt(i) {
  const t = i.store("i18n"), n = {
    locale: lt(),
    messages: V,
    loading: !1,
    init() {
      typeof document < "u" && (document.documentElement.lang = this.locale);
    },
    t(e, o = {}) {
      const s = this.locale, r = this.messages[s] || this.messages.en, a = ct(r, e);
      return (typeof a == "string" ? a : e).replace(/\{(\w+)\}/g, (c, d) => d in o ? String(o[d]) : `{${d}}`);
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
      W(this.messages, e);
    },
    async loadLocale(e, o) {
      if (!this.messages[e]) {
        this.loading = !0;
        try {
          const s = await fetch(o);
          if (!s.ok) throw new Error(`Failed to load ${e}`);
          const r = await s.json();
          this.messages[e] = r, this.setLocale(e);
        } catch (s) {
          f.error(`[Uno] Failed to load locale ${e}:`, s);
        } finally {
          this.loading = !1;
        }
      }
    }
  };
  t ? (f.info("[Uno] Merging into existing i18n store"), t.messages || (t.messages = {}), W(t.messages, V), t.t || (t.t = n.t), t.setLocale || (t.setLocale = n.setLocale), t.setLang || (t.setLang = n.setLang), t.loadLocale || (t.loadLocale = n.loadLocale)) : i.store("i18n", n), i.directive("t", (e, { expression: o, modifiers: s }, { effect: r }) => {
    const a = () => i.store("i18n");
    r(() => {
      const l = a();
      l.locale;
      const c = o, d = l.t(c);
      s.length > 0 ? s.forEach((p) => {
        e.setAttribute(p, d);
      }) : e.textContent = d;
    });
  });
}
function R(i, t, n) {
  if (!i || typeof i.store != "function") {
    f.error("[Uno] Alpine.store is not available");
    return;
  }
  const e = i.store("uno");
  e ? i.store("uno", {
    ...e,
    [t]: n
  }) : i.store("uno", { [t]: n }), i.store(t) || i.store(t, i.store("uno")?.[t]);
}
const ut = (i) => {
  const t = i || "center", n = {
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
}, ft = (i) => {
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
}, gt = (i) => {
  const t = i || "top-right";
  let n = "flex-col", e = "items-end", o = "space-y-2";
  return t.includes("bottom") && (n = "flex-col-reverse", o = "space-y-reverse space-y-2"), t.includes("left") ? e = "items-start" : t.includes("center") ? t === "center" ? e = "items-center justify-center" : e = "items-center" : e = "items-end", { flexDir: n, alignItems: e, spaceY: o };
};
function ht(i, t, n, e = {}) {
  const o = typeof e?.offset == "number" && Number.isFinite(e.offset) && e.offset >= 0 ? e.offset : 8, s = i.getBoundingClientRect(), r = t.getBoundingClientRect(), a = r.width, l = r.height, c = s.width, d = s.height, p = { top: s.top, left: s.left }, g = window.innerWidth, w = window.innerHeight, E = 0, b = 0, h = w, x = g, u = (S, z) => {
    if (z) {
      if (S + a > x)
        return x - a - 10;
      if (S < b)
        return b + 10;
    } else {
      if (S + l > h)
        return h - l - 10;
      if (S < E)
        return E + 10;
    }
    return S;
  };
  let m, v, C = ["top", "bottom", "left", "right"].includes(n) ? n : "bottom";
  switch (n) {
    case "left":
      v = p.left - a - o, v < b && (v = p.left + c + o, C = "right"), m = u(p.top + d / 2 - l / 2, !1);
      break;
    case "top":
      m = p.top - l - o, m < E && (m = p.top + d + o, C = "bottom"), v = p.left + (c - a) / 2, v = u(v, !0);
      break;
    case "right":
      v = p.left + c + o, v + a > x && (v = p.left - a - o, C = "left"), m = u(p.top + d / 2 - l / 2, !1);
      break;
    default:
      m = p.top + d + o, m + l > h && (m = p.top - l - o, C = "top"), v = p.left + (c - a) / 2, v = u(v, !0);
      break;
  }
  const I = 8, T = I / 2, _ = (S, z, Q) => Math.min(Math.max(S, z), Math.max(z, Q)), N = p.left + c / 2, Y = p.top + d / 2;
  let H;
  switch (C) {
    case "top":
      H = {
        top: l - T,
        left: _(N - v - T, 0, a - I)
      };
      break;
    case "left":
      H = {
        top: _(Y - m - T, 0, l - I),
        left: a - T
      };
      break;
    case "right":
      H = {
        top: _(Y - m - T, 0, l - I),
        left: -T
      };
      break;
    default:
      H = {
        top: -T,
        left: _(N - v - T, 0, a - I)
      };
      break;
  }
  return { top: m, left: v, position: C, arrow: H };
}
const pt = () => (y.options?.ui?.layers?.tooltip || $("TOOLTIP"), `
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
const mt = (i) => {
  const t = () => {
    M || typeof requestAnimationFrame == "function" && (M = requestAnimationFrame(() => {
      M = 0;
      const e = i.store("uno")?.tooltip;
      e && e.updatePosition();
    }));
  }, n = () => {
    if (O || typeof window > "u" || typeof document > "u") return;
    const e = () => t(), o = () => t(), s = (a) => {
      if (a.key !== "Escape") return;
      const l = i.store("uno")?.tooltip;
      l?.visible && l.hide();
    }, r = (a) => {
      const l = i.store("uno")?.tooltip;
      if (!l?.visible || l.triggerMode !== "click" || l.clickOutside !== !0) return;
      const c = a.target, d = document.getElementById("uno-tooltip-content"), p = l.triggerEl && l.triggerEl.contains(c), g = d && d.contains(c);
      !p && !g && l.hide();
    };
    window.addEventListener("scroll", e, !0), window.addEventListener("resize", o), window.addEventListener("keydown", s), window.addEventListener("click", r), O = () => {
      window.removeEventListener("scroll", e, !0), window.removeEventListener("resize", o), window.removeEventListener("keydown", s), window.removeEventListener("click", r), O = null;
    };
  };
  return {
    visible: !1,
    content: "",
    placement: "top",
    top: -9999,
    left: -9999,
    zIndex: $("TOOLTIP"),
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
      const o = y.getConfig("tooltip", e);
      if (f.info("Tooltip resolved config", o), !o.target && typeof o.content == "string" && o.content.startsWith("#") && (o.content.includes(" ") || document.querySelector(o.content) && (o.target = o.content)), this.triggerEl = o.element || this.triggerEl, o.target) {
        const s = document.querySelector(o.target);
        s ? this.content = P(s.innerHTML) : (f.warn(`Tooltip target element not found: ${o.target}`), this.content = P(o.content || ""));
      } else
        this.content = P(o.content || "");
      this.placement = o.placement, this.triggerMode = o.trigger, this.clickOutside = o.clickOutside, this.offset = o.offset, this.arrow = o.arrow !== !1, this.zIndex = y.options?.ui?.layers?.tooltip || $("TOOLTIP"), this.hxUrl = o.hxUrl, this.hxHeaders = o.hxHeaders, this.hxTrigger = o.hxTrigger || "load", this.visible = !0, this.ensureContainer(), n(), i.nextTick(() => {
        this.updatePosition();
      });
    },
    hide() {
      this.visible = !1, this.triggerEl = null, O && (O(), O = null), M && (cancelAnimationFrame(M), M = 0);
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
      const o = ht(this.triggerEl, e, this.placement, {
        offset: this.offset
      });
      this.top = o.top, this.left = o.left, this.placement = o.position || this.placement, o.arrow && (this.arrowPosition = o.arrow);
    },
    ensureContainer() {
      const e = "uno-tooltip-manager";
      if (document.getElementById(e)) return;
      const o = y.options?.ui?.layers?.tooltip || $("TOOLTIP"), s = pt();
      q.ensureContainer(e, s, {
        "x-data": "{}",
        style: `position: absolute; top: 0; left: 0; width: 0; height: 0; z-index: ${o};`
      });
      const r = document.getElementById(e);
      r && i && !r._x_dataStack && i.initTree(r);
    }
  };
}, xt = (i) => {
  const t = (n, { expression: e, modifiers: o }, { evaluateLater: s, effect: r, cleanup: a }) => {
    const l = e ? s(e) : null, c = (u) => (o.length > 0 && f.warn(
      "[Uno] x-uno-tooltip modifiers are deprecated in v2.0. Use configuration object instead."
    ), { ...typeof u == "object" && u !== null ? u : u ? { content: u } : {}, element: n }), d = (u) => {
      const m = i.store("uno")?.tooltip;
      m && m.show(c(u));
    }, p = (u) => {
      const m = i.store("uno")?.tooltip;
      m && (m.visible && m.triggerEl === n ? m.hide() : d(u));
    };
    let g = null;
    l && r(() => {
      l((u) => {
        g = u;
      });
    });
    const w = () => {
      if (e && g === null) return;
      const u = c(g);
      y.getConfig("tooltip", u).trigger === "hover" && d(g);
    }, E = () => {
      const u = i.store("uno")?.tooltip;
      u && u.triggerMode === "hover" && u.triggerEl === n && u.scheduleHide();
    }, b = (u) => {
      const m = c(g);
      y.getConfig("tooltip", m).trigger === "click" && (u.preventDefault(), p(g));
    }, h = () => {
      const u = c(g);
      y.getConfig("tooltip", u).trigger === "focus" && d(g);
    }, x = () => {
      const u = i.store("uno")?.tooltip;
      u && u.triggerMode === "focus" && u.triggerEl === n && u.hide();
    };
    n.addEventListener("mouseenter", w), n.addEventListener("mouseleave", E), n.addEventListener("click", b), n.addEventListener("focus", h), n.addEventListener("blur", x), a(() => {
      n.removeEventListener("mouseenter", w), n.removeEventListener("mouseleave", E), n.removeEventListener("click", b), n.removeEventListener("focus", h), n.removeEventListener("blur", x);
      const u = i.store("uno")?.tooltip;
      u && u.triggerEl === n && u.hide();
    });
  };
  i.directive("uno-tooltip", t);
};
const yt = (i, t = {}) => {
  const n = {
    placement: "top",
    trigger: "hover",
    delay: 300,
    offset: 8,
    arrow: !0,
    clickOutside: !0
  };
  y.registerPluginDefaults("tooltip", {
    ...n,
    ...t.plugins?.tooltip,
    ...t.tooltip
  });
  const e = mt(i);
  R(i, "tooltip", e), xt(i);
}, bt = () => `
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
`, vt = () => `
    <div :class="Uno.getStyle('notify', 'progressWrap')"
         style="z-index: ${$("MAX")}"
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
  `, wt = (i) => ({
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
      const o = this.groups[n].find((s) => s.id === t);
      if (o) return o;
    }
    return null;
  },
  onHtmxError(t, n) {
    this._onHtmxEnd(n), f.error("HTMX Error:", t), this.error("Content load failed"), n.onError && n.onError(t);
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
      zIndex: y.options?.ui?.layers?.notify || $("NOTIFY")
    };
  },
  getContainerClass(t) {
    const n = y.getStyle("notify", "container"), e = ft(t), o = gt(t);
    return [n, e, o.flexDir, o.alignItems].join(" ");
  },
  ensureContainer() {
    const t = "uno-notify-container", n = bt(), e = q.ensureContainer(t, n, {
      "x-data": "{}"
      // Independent Alpine component
    });
    e && i.initTree && i.initTree(e);
    const o = "uno-global-progress", s = vt(), r = q.ensureContainer(o, s, {
      "x-data": "{}"
    });
    r && i.initTree && i.initTree(r);
  },
  show(t, n = {}) {
    this.ensureContainer();
    const e = y.getConfig("notify", n), o = this.nextId++;
    let s = t;
    i.store("i18n") && typeof i.store("i18n").t == "function" && (s = i.store("i18n").t(t));
    const r = e.duration !== void 0 ? e.duration : 3e3;
    f.info("[Notify] Show:", { id: o, message: s, duration: r, config: e });
    const a = {
      id: o,
      message: s,
      isVisible: !0,
      // For x-show
      closable: !0,
      duration: r,
      content: e.hxUrl ? "" : void 0,
      // Initialize content for HTMX-loaded content
      ...e
    }, l = e.position || "top-right";
    return this.groups[l] ? this.groups[l].push(a) : (this.groups["top-right"] || (this.groups["top-right"] = []), this.groups["top-right"].push(a)), r > 0 && setTimeout(() => {
      this.remove(o);
    }, r), o;
  },
  remove(t) {
    Object.keys(this.groups).forEach((n) => {
      const e = this.groups[n].findIndex((o) => o.id === t);
      e !== -1 && (this.groups[n][e].isVisible = !1, setTimeout(() => {
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
}), Et = (i) => {
  i.directive("uno-notify", (t, { expression: n }, { evaluateLater: e }) => {
    if (!n) {
      f.warn("[Uno] x-uno-notify directive requires a configuration object or message string.");
      return;
    }
    const o = e(n), s = () => {
      const r = i.store("uno")?.notify;
      r && o((a) => {
        let l = "", c = {};
        if (typeof a == "string")
          l = a;
        else if (typeof a == "object" && a !== null)
          l = a.message || "", c = a;
        else {
          f.warn("[Uno] x-uno-notify directive value must be a string or object", a);
          return;
        }
        r.show(l, c);
      });
    };
    t.addEventListener("click", s);
  });
};
const Ct = (i, t = {}) => {
  const n = {
    duration: 3e3,
    position: "top-right"
  };
  y.registerPluginDefaults("notify", {
    ...n,
    ...t.plugins?.notify,
    ...t.notify
  });
  const e = wt(i);
  R(i, "notify", e), Et(i);
}, kt = () => `
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
           :style="'z-index:' + (${y.options?.ui?.layers?.box || $("BOX")} + index * 10)"
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
                  aria-label="Close">
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
`, Tt = (i) => ({
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
      const { onOk: o, onCancel: s, onClose: r, ...a } = n;
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
          typeof s == "function" && s(), e(!1);
        },
        onClose: () => {
          typeof r == "function" && r(), e(!1);
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
      const { onOk: s, onCancel: r, ...a } = e, l = "prompt-" + Math.random().toString(36).slice(2, 11), c = y.getStyle("box", "promptWrapper"), d = y.getStyle("box", "promptText"), p = y.getStyle("box", "promptInput"), g = y.dom.sanitize(t), w = y.dom.sanitize(n), E = `
        <div class="${c}">
          <p class="${d}">${g}</p>
          <input id="${l}" type="text" class="${p}" value="${w}">
        </div>
      `;
      this.open("confirm", {
        title: a.title || "Prompt",
        ...a,
        content: E,
        _isTrustedPrompt: !0,
        // Internal flag for prompt content
        onOk: () => {
          const b = typeof s == "function" ? s() : void 0;
          if (b === !1) return !1;
          if (b instanceof Promise)
            return b.then((x) => {
              if (x === !1) return !1;
              const u = document.getElementById(l);
              o(u ? u.value : null);
            });
          const h = document.getElementById(l);
          o(h ? h.value : null);
        },
        onCancel: () => {
          typeof r == "function" && r(), o(null);
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
    const s = y.getConfig("box", o), r = s.boxId || Math.random().toString(36).slice(2, 10);
    let a = s.showOkButton, l = s.showCancelButton;
    e === "confirm" ? (a = a !== !1, l = l !== !1) : e === "dialog" && (a = a || !1, l = l || !1);
    const c = {
      id: r,
      type: e,
      title: s.title || "",
      content: s._isTrustedPrompt ? s.content || "" : y.dom.sanitize(s.content || ""),
      // Security: Always sanitize content (trustedHTML removed)
      // Exception: Internal prompt content is pre-escaped
      _isTrustedPrompt: s._isTrustedPrompt || !1,
      width: s.width !== void 0 ? s.width : s.widthClass ? void 0 : s.defaultWidth || 400,
      height: s.height !== void 0 ? s.height : (s.heightClass, void 0),
      widthClass: s.widthClass,
      heightClass: s.heightClass,
      containerClass: s.containerClass,
      autoFocus: s.autoFocus !== void 0 ? s.autoFocus : !0,
      position: s.position || s.defaultPosition || "center",
      modal: s.modal !== void 0 ? s.modal : !0,
      closeOnEsc: s.closeOnEsc !== void 0 ? s.closeOnEsc : !1,
      closeOnOverlayClick: s.closeOnOverlayClick !== void 0 ? s.closeOnOverlayClick : !1,
      showCloseButton: s.showCloseButton !== void 0 ? s.showCloseButton : !0,
      showOkButton: a,
      showCancelButton: l,
      okText: s.okText || "OK",
      cancelText: s.cancelText || "Cancel",
      okButtonClass: s.okButtonClass,
      cancelButtonClass: s.cancelButtonClass,
      extraButtons: Array.isArray(s.extraButtons) ? s.extraButtons : [],
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
        const g = typeof document < "u" ? document.getElementById("uno-box-" + this.id) : null, w = () => {
          const h = i.store("uno")?.box;
          h && h.remove(this.id), this.onClose && this.onClose();
        };
        if (!g) {
          w();
          return;
        }
        let E = !1;
        const b = (h) => {
          h && h.target !== g || E || (E = !0, g.removeEventListener("transitionend", b), w());
        };
        g.addEventListener("transitionend", b), setTimeout(b, 350);
      }
    };
    this.items.push(c), this.ensureContainer();
    const d = this.items[this.items.length - 1], p = () => {
      d.isOpen = !0, this.scheduleAutoFocus(d, i);
    };
    return i && typeof i.nextTick == "function" && i.nextTick(p), d;
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
    return ut(t);
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
      const r = document.getElementById(e);
      r && (t.autoFocus && this.focusFirstInput(r) || this.focusNeutralTarget(r, t));
    }, s = (r) => {
      typeof requestAnimationFrame == "function" ? requestAnimationFrame(() => requestAnimationFrame(r)) : setTimeout(r, 0);
    };
    n && typeof n.nextTick == "function" ? n.nextTick(() => s(o)) : s(o);
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
    ].join(","), e = t.querySelectorAll(n), o = (s) => {
      for (let r = s; r && r.nodeType === 1; r = r.parentElement) {
        const a = window.getComputedStyle(r);
        if (a.display === "none" || a.visibility === "hidden") return !1;
      }
      return !0;
    };
    for (const s of e)
      if (o(s))
        return s.focus(), !0;
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
      return i && typeof i.initTree == "function" && !n._x_dataStack && i.initTree(n), n;
    const o = {
      "x-data": "{}",
      class: "relative",
      style: `z-index: ${y.options?.ui?.layers?.box || $("BOX")};`,
      "x-show": "$store.uno.box.items.length > 0",
      "x-cloak": ""
      // Hide until Alpine is ready
    }, s = kt(), r = y.dom.ensureContainer(t, s, o);
    return r && i && typeof i.initTree == "function" && i.initTree(r), r;
  }
}), $t = [
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
function G(i, t = {}) {
  const { initialFocus: n = !0 } = t, e = (c) => {
    if (!(c instanceof HTMLElement)) return !1;
    const d = window.getComputedStyle(c);
    return d.display === "none" || d.visibility === "hidden" ? !1 : c.getClientRects().length > 0;
  }, o = () => Array.from(i.querySelectorAll($t)).filter((c) => e(c)), s = () => o()[0], r = () => {
    const c = o();
    return c[c.length - 1];
  }, a = (c) => {
    if (c.key === "Tab") {
      if (c.shiftKey) {
        if (document.activeElement === s()) {
          c.preventDefault();
          const d = r();
          d && d.focus();
        }
      } else if (document.activeElement === r()) {
        c.preventDefault();
        const d = s();
        d && d.focus();
      }
    }
  }, l = () => {
    if (i.contains(document.activeElement)) return;
    const c = s();
    c && c.focus();
  };
  return n && (typeof requestAnimationFrame == "function" ? requestAnimationFrame(() => requestAnimationFrame(l)) : setTimeout(l, 0)), i.addEventListener("keydown", a), () => {
    i.removeEventListener("keydown", a);
  };
}
const St = (i) => {
  i.directive("uno-trap", (t, { expression: n }, { cleanup: e, evaluateLater: o }) => {
    if (n)
      o(n)((r) => {
        const a = G(t, r || {});
        e(a);
      });
    else {
      const s = G(t);
      e(s);
    }
  }), i.directive("uno-box", (t, { expression: n }, { evaluateLater: e, cleanup: o }) => {
    if (!n) {
      f.warn("[Uno] x-uno-box directive requires a configuration object.");
      return;
    }
    const s = e(n), r = (a) => {
      a.preventDefault();
      const l = i.store("uno")?.box;
      if (!l) {
        f.warn("[Uno] Box store not found. Is Uno initialized?");
        return;
      }
      s((c) => {
        c && typeof c == "object" ? l.open(c) : f.warn("[Uno] x-uno-box directive value must be an object", c);
      });
    };
    t.addEventListener("click", r), o(() => t.removeEventListener("click", r));
  });
};
const Lt = (i, t = {}) => {
  const n = {
    defaultWidth: 400,
    defaultPosition: "center",
    okText: "OK",
    cancelText: "Cancel",
    closeOnEsc: !1,
    closeOnOverlayClick: !1
  };
  y.registerPluginDefaults("box", { ...n, ...t.plugins?.box, ...t.box });
  const e = Tt(i);
  R(i, "box", e), St(i);
}, J = {
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
}, It = (i) => ({
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
    const o = this.normalizeRules(n), s = y.getConfig("valid", e);
    this.rules[t] = o.map((r) => ({
      emptyQuit: s.emptyQuit !== !1,
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
    const o = this.rules[t];
    if (!o || o.length === 0) return !0;
    delete this.errors[t], this.timers[t] && (clearTimeout(this.timers[t]), delete this.timers[t]);
    const s = n == null ? "" : String(n), r = s.trim() === "";
    for (const a of o) {
      if (e === "input" && (a.trigger === "blur" || a.realtime === !1))
        continue;
      const l = a.type === "required" || a.required === !0;
      if (!(r && !l && a.emptyQuit !== !1)) {
        if (a.debounce > 0 && e === "input")
          return this.validating[t] = !0, new Promise((c) => {
            this.timers[t] = setTimeout(async () => {
              const d = await this.executeRule(t, n, a, s);
              delete this.validating[t], c(!!d);
            }, a.debounce);
          });
        this.validating[t] = !0;
        try {
          if (!await this.executeRule(t, n, a, s))
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
    let s = !0, r = "";
    const a = e.type;
    if (a === "required" || e.required)
      o.trim() === "" && (s = !1, r = e.msg || this.getI18nMsg("valid.required") || "This field is required");
    else if (J[a])
      J[a].test(o) || (s = !1, r = e.msg || this.getI18nMsg(`valid.${a}`) || `Invalid ${a} format`);
    else if (a === "length")
      e.min !== void 0 && o.length < e.min ? (s = !1, r = e.msg || this.getI18nMsg("valid.min", { n: e.min }) || `Min length is ${e.min}`) : e.max !== void 0 && o.length > e.max && (s = !1, r = e.msg || this.getI18nMsg("valid.max", { n: e.max }) || `Max length is ${e.max}`);
    else if (a === "range") {
      const l = parseFloat(o);
      isNaN(l) ? (s = !1, r = e.msg || this.getI18nMsg("valid.number") || "Must be a number") : e.min !== void 0 && l < e.min ? (s = !1, r = e.msg || this.getI18nMsg("valid.rangeMin", { n: e.min }) || `Min value is ${e.min}`) : e.max !== void 0 && l > e.max && (s = !1, r = e.msg || this.getI18nMsg("valid.rangeMax", { n: e.max }) || `Max value is ${e.max}`);
    } else if (a === "pattern" || e.pattern)
      (e.pattern instanceof RegExp ? e.pattern : new RegExp(e.pattern)).test(o) || (s = !1, r = e.msg || this.getI18nMsg("valid.pattern") || "Invalid format");
    else if (a === "match" || e.match) {
      const l = e.match || e.selector || e.target;
      if (!l)
        return !0;
      const c = document.querySelector(l.startsWith("#") ? l : `[name="${l}"]`) || document.getElementById(l);
      c && o !== c.value && (s = !1, r = e.msg || this.getI18nMsg("valid.match") || "Fields do not match");
    } else if (a === "custom" || a === "validator" || e.fn || this.customValidators[a]) {
      let l = e.fn || e.name;
      if (typeof l == "string" ? l = this.customValidators[l] : !l && this.customValidators[a] && (l = this.customValidators[a]), typeof l == "function")
        try {
          const c = await l(n, t, e);
          c === !1 ? (s = !1, r = e.msg || "Validation failed") : typeof c == "string" ? (s = !1, r = c) : typeof c == "object" && c.valid === !1 && (s = !1, r = c.message || e.msg || "Validation failed");
        } catch (c) {
          f.error(`[Valid] Error in custom validator for ${t}`, c), s = !1, r = "Validation error";
        }
      else
        f.warn(`[Valid] Custom validator not found: ${a || e.fn}`);
    }
    return s || (this.errors[t] = r), s;
  },
  /**
   * Validates all fields in a form.
   *
   * @param {HTMLFormElement} form - The form element to validate
   * @returns {Promise<boolean>} True if the form is valid
   */
  async validateForm(t) {
    const n = t.querySelectorAll("input, select, textarea");
    let e = !0, o = null;
    for (const s of n) {
      if (s.disabled) continue;
      const r = s.name || s.id;
      this.rules[r] && (await this.validate(r, s.value, "submit") || (e = !1, o || (o = s)));
    }
    return !e && o && (y.getConfig("valid") || {}).scrollOnError !== !1 && (o.scrollIntoView({ behavior: "smooth", block: "center" }), o.focus({ preventScroll: !0 })), e;
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
}), X = "data-uno-counter-wrapper", Bt = (i) => i === "flex" || i === "inline-flex" || i === "grid" || i === "inline-grid", Ut = (i) => {
  const t = i.parentElement;
  if (!t) return null;
  if (t.hasAttribute(X)) return t;
  const n = window.getComputedStyle(t);
  if (Bt(n.display))
    return n.position === "static" && t.tagName !== "FORM" && t.tagName !== "BODY" && (t.style.position = "relative"), t;
  const e = i.ownerDocument.createElement("div");
  e.setAttribute(X, ""), e.style.position = "relative";
  const o = window.getComputedStyle(i).display;
  return (o === "inline" || o === "inline-block") && (e.style.display = "inline-block"), t.insertBefore(e, i), e.appendChild(i), e;
}, Mt = (i, t) => {
  if (!t || !t.hasAttribute(X)) return;
  const n = t.parentNode;
  n && (n.insertBefore(i, t), n.removeChild(t));
}, Ot = (i) => {
  i.directive("uno-valid", (t, { expression: n }, { evaluateLater: e, cleanup: o }) => {
    if (t.tagName === "FORM") return;
    const s = t.name || t.id;
    if (!s) {
      f.warn("[Uno Valid] Element has no name or id:", t);
      return;
    }
    const r = Uno.getConfig("valid"), a = Uno.options.ui?.valid || {};
    if (r.autoError !== !1) {
      const g = `$store.uno.valid.errors['${s}']`, w = `[x-text="${g.replace(/"/g, '\\"')}"]`, E = t.form || t.closest("form");
      if (!(E ? E.querySelector(w) : document.querySelector(w))) {
        let h;
        if (typeof r.createErrorElement == "function")
          h = r.createErrorElement(s, g);
        else {
          h = document.createElement("div"), h.setAttribute("x-text", g), h.setAttribute("x-show", g);
          const x = r.errorClass || a.helperText || "text-xs text-red-500 mt-1";
          h.className = x;
        }
        h && (t.after(h), i.initTree(h));
      }
    }
    const l = n ? e(n) : () => {
    }, c = (g) => {
      i.store("uno").valid.addRules(s, g);
    };
    n ? l((g) => {
      c(g);
      const w = i.store("uno").valid;
      if (w && w.normalizeRules) {
        const b = w.normalizeRules(g).find(
          (h) => h.type === "length" || h.maxLength !== void 0 || h.max !== void 0
        );
        if (b) {
          const h = b.max !== void 0 ? b.max : b.maxLength;
          if (b.stopOnMax === !0 && h !== void 0 && t.setAttribute("maxlength", h), b.showCounter === !0 && h !== void 0) {
            const x = `${s}-counter`;
            let u = document.getElementById(x);
            const m = Ut(t);
            u || (u = document.createElement("span"), u.id = x, u.className = "absolute bottom-1 right-2 text-xs text-gray-400 pointer-events-none transition-colors duration-200"), u.parentElement !== m && t.after(u);
            const v = () => {
              const C = t.value.length;
              u.textContent = `${C}/${h}`, C > h ? (u.classList.add("text-red-500"), u.classList.remove("text-gray-400")) : (u.classList.remove("text-red-500"), u.classList.add("text-gray-400"));
            };
            v(), t.addEventListener("input", v), o(() => {
              t.removeEventListener("input", v), u && u.parentNode && u.parentNode.removeChild(u), Mt(t, m);
            });
          }
        }
      }
    }) : c([]);
    const d = () => {
      i.store("uno").valid.validate(s, t.value, "input");
    }, p = () => {
      i.store("uno").valid.validate(s, t.value, "blur");
    };
    t.addEventListener("input", d), t.addEventListener("blur", p), o(() => {
      t.removeEventListener("input", d), t.removeEventListener("blur", p);
    });
  }), typeof document < "u" && (document.body.addEventListener("htmx:confirm", async (t) => {
    const n = t.detail.elt, e = n.tagName === "FORM" ? n : n.closest("form");
    if (e && (e.hasAttribute("x-uno-valid") || e.querySelector("[x-uno-valid]"))) {
      t.preventDefault();
      try {
        await i.store("uno").valid.validateForm(e) && t.detail.issueRequest(!0);
      } catch (o) {
        f.error("[Valid] Form validation failed unexpectedly, request blocked", o);
      }
    }
  }), document.body.addEventListener("htmx:afterRequest", (t) => {
    const n = A(t);
    if (n && n.status === 422)
      try {
        const e = JSON.parse(n.text);
        e.errors && (i.store("uno").valid.showErrors(e.errors), f.info("[Valid] Server validation errors applied", e.errors));
      } catch {
      }
  }));
};
const Ht = (i, t = {}) => {
  const n = {
    realtime: !0,
    showError: !0,
    scrollOnError: !0
  };
  y.registerPluginDefaults("valid", {
    ...n,
    ...t.plugins?.valid,
    ...t.valid
  });
  const e = It(i);
  R(i, "valid", e), Ot(i);
}, Pt = (i) => ({
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
        const { group: e, name: o } = this.registry[n];
        f.info(`[Tab] checkHash match: ${o}`), this.activate({ group: e, name: o, updateHash: !1 });
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
    const n = y.getConfig("tab", {
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
}), Rt = (i) => {
  i.directive("uno-tab", (t, { expression: n }, { evaluateLater: e, cleanup: o }) => {
    if (!n) {
      f.warn("[Uno] x-uno-tab directive requires a configuration object.");
      return;
    }
    const s = e(n);
    s((a) => {
      let l = {};
      if (typeof a == "object" && a !== null)
        l = a;
      else if (typeof a == "string")
        l = { name: a };
      else {
        f.warn("[Uno] x-uno-tab directive value must be a string or object", a);
        return;
      }
      l.group || (l.group = "default"), l.name && i.store("uno").tab.register(l.name, l.group);
    });
    const r = (a) => {
      a.preventDefault(), s((l) => {
        let c = {};
        typeof l == "object" && l !== null ? c = l : typeof l == "string" && (c = { name: l }), c.name && i.store("uno").tab.activate(c);
      });
    };
    t.addEventListener("click", r), o(() => t.removeEventListener("click", r));
  });
};
const _t = (i, t = {}) => {
  const n = {
    defaultTab: "",
    hash: !0,
    activeClass: "is-active",
    inactiveClass: ""
  };
  y.registerPluginDefaults("tab", { ...n, ...t.plugins?.tab, ...t.tab });
  const e = Pt();
  R(i, "tab", e), Rt(i);
};
y.dom = q;
typeof window < "u" && (window.Uno = y);
y.use(dt).use(yt).use(Ct).use(Lt).use(Ht).use(_t);
if (typeof window < "u") {
  const i = () => {
    window.Alpine && !y.initialized && y.init(window.Alpine);
  };
  window.Alpine ? i() : document.addEventListener(
    "alpine:init",
    () => {
      i();
    },
    { once: !0 }
  );
}
export {
  q as DomManager,
  y as Uno,
  Lt as box,
  y as default,
  Ct as notify,
  _t as tab,
  yt as tooltip,
  Ht as valid
};
