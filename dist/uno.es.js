import V from "alpinejs";
let U = !1, Y = null;
const h = {
  /**
   * Set debug mode
   * @param {boolean} debug
   */
  setDebug(r) {
    U = r;
  },
  /**
   * Set custom error reporter
   * @param {Function} reporter - (error, ...args) => void
   */
  setReporter(r) {
    Y = r;
  },
  /**
   * Log info message
   * @param {...any} args - Arguments to log
   */
  info(...r) {
    U && typeof console < "u" && typeof console.log == "function" && console.log("[Uno Info]", ...r);
  },
  /**
   * Log debug message
   * @param {...any} args - Arguments to log
   */
  debug(...r) {
    U && typeof console < "u" && typeof console.log == "function" && console.log("[Uno Debug]", ...r);
  },
  /**
   * Log warning message
   * @param {...any} args - Arguments to log
   */
  warn(...r) {
    U && typeof console < "u" && typeof console.warn == "function" && console.warn("[Uno Warn]", ...r);
  },
  /**
   * Log error message
   * @param {...any} args - Arguments to log
   */
  error(...r) {
    if (Y)
      try {
        Y(...r);
      } catch (e) {
        U && typeof console < "u" && typeof console.error == "function" && console.error("[Uno Error] Reporter failed:", e);
      }
    U && typeof console < "u" && typeof console.error == "function" && console.error("[Uno Error]", ...r);
  }
}, ne = /* @__PURE__ */ new WeakMap();
let se = 0;
function q(r, e) {
  if (!r) return !1;
  const t = String(e), n = `${se}:${t}`;
  let o = ne.get(r);
  return o || (o = /* @__PURE__ */ new Set(), ne.set(r, o)), o.has(n) ? !1 : (o.add(n), !0);
}
function D(r) {
  return r ? r.replace(/<script\b[^>]*>([\s\S]*?)<\/script>/gim, "").replace(/<img\b[^>]*>/gim, "").replace(/\s+on[a-z]+="[^"]*"/gim, "").replace(/\s+on[a-z]+='[^']*'/gim, "").replace(/\s+href="javascript:[^"]*"/gim, "").replace(/\s+href='javascript:[^']*'/gim, "").replace(/javascript:/gim, "") : "";
}
function ae(r, e, { trusted: t = !1, alpine: n } = {}) {
  r && (r.innerHTML = t ? String(e ?? "") : D(String(e ?? "")), n?.initTree?.(r));
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
  ensureContainer(r, e = "", t = {}, n = document.body) {
    const o = document.getElementById(r);
    if (o) return o;
    try {
      const i = document.createElement("div");
      return i.id = r, window.Uno?.options?.csp?.nonce && !t.nonce && (t.nonce = window.Uno.options.csp.nonce), Object.entries(t).forEach(([s, a]) => {
        i.setAttribute(s, a);
      }), e && (i.innerHTML = e), n ? n.appendChild(i) : document.documentElement.appendChild(i), i;
    } catch (i) {
      return h.error(`DomManager: Failed to inject container #${r}`, i), null;
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
  injectStyles(r, e) {
    if (!document.getElementById(r))
      try {
        const t = document.createElement("style");
        t.id = r, window.Uno?.options?.csp?.nonce && t.setAttribute("nonce", window.Uno.options.csp.nonce), t.textContent = e, document.head.appendChild(t);
      } catch (t) {
        h.error(`DomManager: Failed to inject styles #${r}`, t);
      }
  },
  // Exports
  setHTML: ae,
  // Secure HTML setting
  sanitize: D
  // HTML sanitization
}, le = (r, e) => {
  try {
    return r.getResponseHeader(e);
  } catch {
    return null;
  }
}, N = (r) => {
  const t = r?.detail?.xhr;
  return t ? {
    status: t.status,
    text: t.responseText ?? "",
    header: (n) => le(t, n)
  } : null;
}, ce = (r) => r?.detail?.xhr;
function B(r) {
  return Object.prototype.toString.call(r) === "[object Object]";
}
function M(r, e) {
  if (!B(r)) return B(e) ? { ...e } : e;
  if (!B(e)) return e === void 0 ? { ...r } : e;
  const t = { ...r };
  for (const [n, o] of Object.entries(e)) {
    if (B(o) && B(t[n])) {
      t[n] = M(t[n], o);
      continue;
    }
    if (B(o)) {
      t[n] = M({}, o);
      continue;
    }
    if (typeof o == "string" && o.trim().indexOf("+") === 0 && typeof t[n] == "string") {
      t[n] = `${t[n]} ${o.trim().slice(1)}`.trim();
      continue;
    }
    t[n] = o;
  }
  return t;
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
  const e = r.toUpperCase();
  return L[e] !== void 0 ? L[e] : L.BASE;
};
class de {
  /**
   * Create an UnoInstance.
   *
   * @param {Object} [options={}] - Global configuration options
   */
  constructor(e = {}) {
    this.version = "0.9.0", this.options = e, this.pluginDefaults = {}, this.plugins = [], this.initialized = !1, this._preInitConfig = {}, this._htmxRequestHandlers = [], this._htmxActiveXhrs = typeof WeakSet == "function" ? /* @__PURE__ */ new WeakSet() : null, this.htmx = {
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
          throw h.error("HTMX fragment fetch failed:", { method: o, url: i, error: u }), u;
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
        if (!o || typeof o != "object" || (this._htmxRequestHandlers.push(o), typeof window > "u" || typeof document > "u") || !window.htmx || !q(document, "uno:htmx:request-hooks")) return;
        const i = ce, s = (c, u) => {
          for (const g of this._htmxRequestHandlers) {
            const d = g?.[c];
            if (typeof d == "function")
              try {
                d(u);
              } catch (m) {
                h.error("Uno.htmx.onRequest handler failed", m);
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
          h.warn("bindNotifyToHtmx: document is not available");
          return;
        }
        if (!u.htmx.available()) {
          h.warn("bindNotifyToHtmx: HTMX is not available");
          return;
        }
        const g = () => typeof window < "u" && window.Alpine ? window.Alpine.store("uno")?.notify : null;
        let d = 0;
        const m = (p) => {
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
            const x = N(p);
            if (!(x?.status >= 400)) {
              const v = g();
              if (v && typeof v.success == "function") {
                let k = l;
                const E = x?.header("X-Notify-Message");
                E && (k = E), x?.header("X-Notify-Suppress") || v.success(k);
              }
            }
          }
        }, y = (p) => {
          if (d = Math.max(0, d - 1), a && d === 0) {
            const x = g();
            x && typeof x.stopProgress == "function" && x.stopProgress();
          }
          if (s) {
            const x = g();
            if (x && typeof x.error == "function") {
              const f = N(p);
              let v = c;
              const k = f?.header("X-Notify-Error");
              k && (v = k), x.error(v);
            }
          }
        };
        q(document, "uno:htmx:bind-notify") && (document.addEventListener("htmx:beforeRequest", m), document.addEventListener("htmx:afterRequest", w), document.addEventListener("htmx:responseError", y), document.addEventListener("htmx:sendError", y), document.addEventListener("htmx:timeout", y), h.debug("bindNotifyToHtmx: HTMX events bound to notify"));
      }
    };
    const n = {
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
          container: "relative transform overflow-hidden {base} bg-white dark:bg-gray-900 text-start shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg flex flex-col max-h-[calc(100dvh-2rem)] sm:max-h-[calc(100dvh-4rem)] pointer-events-auto outline-none",
          // fill 模式：未指定 height/heightClass 时的默认高度；dvh 排除移动端地址栏，
          // 手机扣除 itemWrapper p-4，sm 起扣除 container my-8
          containerFill: "h-[calc(100dvh-2rem)] sm:h-[min(85dvh,56rem)]",
          cancelBtn: "mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:mt-0 sm:w-auto",
          extraBtn: "mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-normal text-gray-700 hover:bg-gray-50 sm:mt-0 sm:w-auto",
          okBtn: "inline-flex w-full justify-center rounded-md bg-blue-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 sm:ms-3 sm:w-auto",
          title: "text-base font-semibold leading-6 text-gray-900 dark:text-slate-200",
          // 目标尺寸（D5，WCAG 2.5.8）：p-1 使图标命中区达 28×28；-m-1 抵消绝对定位偏移
          // 非文本对比度（E7，WCAG 1.4.11）：亮色 slate-500（hover slate-600）；
          // 暗底配对 slate-400/300
          closeBtnFloating: "absolute top-2 end-2 z-10 p-1 -m-1 text-slate-500 hover:text-slate-600 dark:text-slate-400 dark:hover:text-slate-300 focus:outline-none",
          closeIconSmall: "h-5 w-5",
          content: "px-4 pb-4 pt-5 sm:p-6 sm:pb-4 overflow-auto flex-1 text-slate-600 dark:text-slate-300",
          // fill 模式：内容区与正文容器逐层传递剩余高度，供加载内容以 flex-1 撑满
          contentFill: "flex flex-col min-h-0",
          bodyFill: "flex-1 min-h-0 flex flex-col",
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
          textWrap: "ms-3 min-w-0 flex-1 pt-0.5",
          text: "text-sm font-medium text-gray-900 dark:text-slate-100 break-words",
          closeWrap: "ms-4 flex flex-shrink-0",
          // 目标尺寸（D5，WCAG 2.5.8）：p-1 使图标命中区达 28×28；-m-1 抵消行内占位
          closeBtn: "inline-flex rounded-md p-1 -m-1 text-gray-500 hover:text-gray-700 hover:bg-black/5 dark:text-slate-400 dark:hover:text-slate-200 dark:hover:bg-white/10 focus:outline-none focus:ring-2 focus:ring-indigo-500",
          closeIcon: "h-5 w-5",
          progressWrap: "fixed top-0 start-0 end-0 h-1 bg-gray-200 dark:bg-gray-700 pointer-events-none",
          progressBar: "h-full bg-gradient-to-r from-blue-500 to-indigo-600 shadow-sm shadow-blue-500/50 rtl:bg-gradient-to-l",
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
            top: "border-b border-e",
            // tooltip above trigger -> tip points down
            bottom: "border-t border-s",
            // tip points up
            left: "border-t border-e",
            // tip points right
            right: "border-b border-s"
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
          iconBase: "absolute end-0 top-0 bottom-0 flex items-center justify-center w-10 pointer-events-none transition-opacity duration-200",
          inputError: "!border-red-500",
          inputValid: "!border-green-500",
          iconError: "text-red-500 dark:text-red-400",
          iconValid: "text-green-500 dark:text-green-400",
          stateInvalid: "is-invalid",
          stateValid: "is-valid",
          stateValidating: "validating",
          inputPaddingIcon: "pe-10",
          helperText: "text-red-600 dark:text-red-400 text-xs mt-1"
        }
      }
    };
    this.options = M(n, e);
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
  getConfig(e, t = {}) {
    const n = this.options || {};
    let o = this.pluginDefaults[e] ? { ...this.pluginDefaults[e] } : {};
    return n.plugins && n.plugins[e] && (o = M(o, n.plugins[e])), n[e] && (o = M(o, n[e])), t && (o = M(o, t)), o;
  }
  /**
   * ✅ Register plugin defaults
   *
   * @param {string} name - Component name
   * @param {Object} defaults - Default configuration object
   */
  registerPluginDefaults(e, t) {
    this.pluginDefaults[e] = t;
  }
  _getUiValue(e, t) {
    const n = this.options.ui?.[e];
    if (!n) return;
    if (typeof t != "string" || t.indexOf(".") === -1) return n?.[t];
    const o = t.split(".").filter(Boolean);
    let i = n;
    for (const s of o) {
      if (i == null) return;
      i = i[s];
    }
    return i;
  }
  _resolveUiTokens(e) {
    if (typeof e != "string") return e;
    const t = this.options.ui?.base;
    return typeof t == "string" && e.indexOf("{base}") !== -1 ? e.replaceAll("{base}", t) : e;
  }
  /**
   * Get style class from configuration or fallback
   * Support "+" prefix to append class to default
   * @param {string} comp - Component name
   * @param {string} elem - Element name
   * @param {string} defaultCls - Default class string
   * @returns {string} Final class string
   */
  getStyle(e, t, n) {
    const o = this._getUiValue(e, t);
    return o === void 0 ? this._resolveUiTokens(n) : typeof o == "string" && o.trim().indexOf("+") === 0 ? this._resolveUiTokens(`${n} ${o.trim().slice(1)}`.trim()) : this._resolveUiTokens(o);
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
  setConfig(e) {
    return !e || typeof e != "object" ? (h.warn("setConfig called with invalid argument, expected object"), this) : (this.initialized ? (h.info("Uno configuration updated after initialization"), this.options = M(this.options, e)) : this._preInitConfig = M(this._preInitConfig, e), this);
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
  use(e, t = {}) {
    return this.initialized && h.warn("Uno.use() called after initialization. Plugin might not work correctly."), this.plugins.push({ plugin: e, options: t }), this;
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
  async init(e = {}, t = {}) {
    if (this.initialized)
      return h.warn("Uno has already been initialized."), V || (typeof window < "u" ? window.Alpine : null);
    let n = V;
    return !n && typeof window < "u" && window.Alpine && (n = window.Alpine), e && (e.version || typeof e.start == "function" || typeof e.store == "function") && (n = e), this.options = M(this.options, this._preInitConfig), h.setDebug(this.options.debug !== !1), this.options.onError && h.setReporter(this.options.onError), n && typeof n.store == "function" ? (this._doInit(n), n) : typeof document > "u" ? (h.error("Alpine not found and document is unavailable. Cannot initialize Uno."), null) : (h.debug("Alpine not ready, waiting for alpine:init event..."), new Promise((i) => {
      document.addEventListener(
        "alpine:init",
        () => {
          if (n = window.Alpine || V, !n) {
            h.error("Alpine still not found after alpine:init event"), i(null);
            return;
          }
          h.debug("Alpine ready, continuing Uno initialization"), this._doInit(n), i(n);
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
  _doInit(e) {
    h.info(`Uno v${this.version} initializing...`), typeof window < "u" && typeof document < "u" && window.htmx && typeof e?.initTree == "function" && q(document, "uno:htmx-afterSwap:initTree") && document.addEventListener("htmx:afterSwap", (t) => {
      const n = t?.detail?.elt;
      n && e.initTree(n);
    }), this.plugins.forEach(({ plugin: t, options: n }) => {
      try {
        const o = M(this.options, n);
        t(e, o);
      } catch (o) {
        h.error("Failed to load plugin", o);
      }
    }), this._setupGlobalErrorHandlers(), this.initialized = !0, h.info("Uno initialized");
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
    if (!(typeof window > "u") && q(window, "uno:global-error-handlers")) {
      const e = typeof WeakSet == "function" ? /* @__PURE__ */ new WeakSet() : null;
      window.addEventListener("unhandledrejection", (t) => {
        const n = t.reason;
        if (e && n instanceof Error) {
          if (e.has(n)) return;
          e.add(n);
        }
        h.error("Unhandled Promise Rejection:", n);
      }), window.addEventListener("error", (t) => {
        const n = t.error;
        if (e && n instanceof Error) {
          if (e.has(n)) return;
          e.add(n);
        }
        h.error("Uncaught Error:", {
          message: t.message,
          filename: t.filename,
          lineno: t.lineno,
          colno: t.colno,
          error: n
        });
      }), h.debug("Global error handlers installed");
    }
  }
}
const b = new de(), ue = {
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
    compareto: "Does not match",
    tagMax: "Max {n} items",
    callback: "Validation failed",
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
}, fe = {
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
}, me = {
  box: {
    ok: "確認",
    cancel: "取消",
    loading: "載入中",
    close: "關閉"
  },
  valid: {
    required: "必填",
    pattern: "格式無效",
    min: "至少 {n} 個字元",
    max: "最多 {n} 個字元",
    counter: "已輸入 {current} 字，共 {max} 字",
    range: "必須在 {min} 到 {max} 之間",
    rangeMin: "最小 {n}",
    rangeMax: "最大 {n}",
    match: "輸入不一致",
    compareto: "不一致",
    tagMax: "最多 {n} 項",
    callback: "驗證失敗",
    validname: "僅限字母、數字、- 和 _ (3-16 個字元)",
    en: "僅限英文字母",
    domain: "無效的域名格式",
    mobile: "無效的手機號碼",
    int: "需要整數",
    email: "無效的電子郵件地址",
    date: "無效的日期格式 (YYYY-MM-DD)",
    url: "無效的 URL",
    pwd: "密碼需 8-72 個字元，必須包含字母和數字，可含 _!@#$%^&*()+-= 特殊字元",
    strongpwd: "密碼需 8-80 個字元，必須包含字母、數字和至少 1 個特殊字元 (_@#$%^&*()+-=)",
    summary_title: "請修正以下問題"
  },
  charCount: "剩餘: {count}"
}, ge = {
  box: {
    ok: "Aceptar",
    cancel: "Cancelar",
    loading: "Cargando",
    close: "Cerrar"
  },
  valid: {
    required: "Obligatorio",
    pattern: "Formato no válido",
    min: "Al menos {n} caracteres",
    max: "Máximo {n} caracteres",
    counter: "{current} de {max} caracteres ingresados",
    range: "Debe estar entre {min} y {max}",
    rangeMin: "Mín. {n}",
    rangeMax: "Máx. {n}",
    match: "Los campos no coinciden",
    compareto: "No coincide",
    tagMax: "Máximo {n} elementos",
    callback: "Error de validación",
    validname: "Solo letras, números, guion (-) y guion bajo (_), de 3 a 16 caracteres",
    en: "Solo letras del alfabeto inglés",
    domain: "Formato de dominio no válido",
    mobile: "Número de celular no válido",
    int: "Se requiere un número entero",
    email: "Correo electrónico no válido",
    date: "Formato de fecha no válido (YYYY-MM-DD)",
    url: "URL no válida",
    pwd: "La contraseña debe tener entre 8 y 72 caracteres, incluir letras y números y puede contener _!@#$%^&*()+-=",
    strongpwd: "La contraseña debe tener entre 8 y 80 caracteres, incluir letras, números y al menos 1 carácter especial (_@#$%^&*()+-=)",
    summary_title: "Corrige los siguientes problemas"
  },
  charCount: "Restantes: {count}"
}, he = {
  box: {
    ok: "確認",
    cancel: "キャンセル",
    loading: "読み込み中",
    close: "閉じる"
  },
  valid: {
    required: "必須です",
    pattern: "形式が正しくありません",
    min: "{n} 文字以上で入力してください",
    max: "{n} 文字以内で入力してください",
    counter: "{current}/{max} 文字入力済み",
    range: "{min}〜{max} の範囲で入力してください",
    rangeMin: "{n} 以上",
    rangeMax: "{n} 以下",
    match: "値が一致しません",
    compareto: "一致しません",
    tagMax: "{n} 件まで追加できます",
    callback: "入力内容の検証に失敗しました",
    validname: "半角英数字とハイフン（-）、アンダースコア（_）のみ使用できます（3〜16 文字）",
    en: "英字のみ入力できます",
    domain: "ドメインの形式が正しくありません",
    mobile: "携帯電話番号の形式が正しくありません",
    int: "整数で入力してください",
    email: "メールアドレスの形式が正しくありません",
    date: "日付の形式が正しくありません（YYYY-MM-DD）",
    url: "URL の形式が正しくありません",
    pwd: "パスワードは 8〜72 文字で、英字と数字を含み、_!@#$%^&*()+-= も使用できます",
    strongpwd: "パスワードは 8〜80 文字で、英字・数字を含み、記号（_@#$%^&*()+-=）を 1 文字以上含めてください",
    summary_title: "次の問題を修正してください"
  },
  charCount: "残り：{count}"
}, pe = {
  box: {
    ok: "OK",
    cancel: "Abbrechen",
    loading: "Wird geladen",
    close: "Schließen"
  },
  valid: {
    required: "Erforderlich",
    pattern: "Ungültiges Format",
    min: "Mindestens {n} Zeichen",
    max: "Max. {n} Zeichen",
    counter: "{current} von {max} Zeichen eingegeben",
    range: "Muss zwischen {min} und {max} liegen",
    rangeMin: "Min. {n}",
    rangeMax: "Max. {n}",
    match: "Felder stimmen nicht überein",
    compareto: "Stimmt nicht überein",
    tagMax: "Max. {n} Elemente",
    callback: "Validierung fehlgeschlagen",
    validname: "3-16 Zeichen; nur Buchstaben, Zahlen, Bindestrich (-) und Unterstrich (_)",
    en: "Nur englische Buchstaben",
    domain: "Ungültiges Domainformat",
    mobile: "Ungültige Mobilnummer",
    int: "Ganzzahl erforderlich",
    email: "Ungültige E-Mail-Adresse",
    date: "Ungültiges Datumsformat (YYYY-MM-DD)",
    url: "Ungültige URL",
    pwd: "Das Passwort muss 8-72 Zeichen lang sein, Buchstaben und Zahlen enthalten; außerdem sind _!@#$%^&*()+-= erlaubt",
    strongpwd: "Das Passwort muss 8-80 Zeichen lang sein, Buchstaben, Zahlen und mindestens ein Sonderzeichen (_@#$%^&*()+-=) enthalten",
    summary_title: "Bitte behebe die folgenden Probleme"
  },
  charCount: "Verbleibend: {count}"
}, xe = {
  box: {
    ok: "Valider",
    cancel: "Annuler",
    loading: "Chargement",
    close: "Fermer"
  },
  valid: {
    required: "Obligatoire",
    pattern: "Format invalide",
    min: "Au moins {n} caractères",
    max: "Maximum {n} caractères",
    counter: "{current} sur {max} caractères saisis",
    range: "Doit être compris entre {min} et {max}",
    rangeMin: "Min. {n}",
    rangeMax: "Max. {n}",
    match: "Les champs ne correspondent pas",
    compareto: "Ne correspond pas",
    tagMax: "Maximum {n} éléments",
    callback: "Échec de la validation",
    validname: "Lettres, chiffres, tiret (-) et tiret bas (_) uniquement (3 à 16 caractères)",
    en: "Lettres anglaises uniquement",
    domain: "Format de domaine invalide",
    mobile: "Numéro de mobile invalide",
    int: "Nombre entier requis",
    email: "Adresse e-mail invalide",
    date: "Format de date invalide (AAAA-MM-JJ)",
    url: "URL invalide",
    pwd: "Le mot de passe doit comporter entre 8 et 72 caractères, contenir des lettres et des chiffres, et peut inclure _!@#$%^&*()+-=",
    strongpwd: "Le mot de passe doit comporter entre 8 et 80 caractères, contenir des lettres, des chiffres et au moins un caractère spécial (_@#$%^&*()+-=)",
    summary_title: "Veuillez corriger les problèmes suivants"
  },
  charCount: "Caractères restants : {count}"
}, ye = {
  box: {
    ok: "확인",
    cancel: "취소",
    loading: "불러오는 중",
    close: "닫기"
  },
  valid: {
    required: "필수 입력 항목이에요",
    pattern: "형식이 올바르지 않아요",
    min: "{n}자 이상 입력해 주세요",
    max: "{n}자 이하로 입력해 주세요",
    counter: "{max}자 중 {current}자를 입력했어요",
    range: "{min}에서 {max} 사이여야 해요",
    rangeMin: "최소 {n}",
    rangeMax: "최대 {n}",
    match: "입력한 값이 일치하지 않아요",
    compareto: "일치하지 않아요",
    tagMax: "최대 {n}개",
    callback: "검증에 실패했어요",
    validname: "영문, 숫자, 하이픈(-), 밑줄(_)만 사용할 수 있어요 (3~16자)",
    en: "영문만 입력할 수 있어요",
    domain: "도메인 형식이 올바르지 않아요",
    mobile: "휴대폰 번호가 올바르지 않아요",
    int: "정수만 입력할 수 있어요",
    email: "이메일 주소가 올바르지 않아요",
    date: "날짜 형식이 올바르지 않아요 (YYYY-MM-DD)",
    url: "URL 형식이 올바르지 않아요",
    pwd: "비밀번호는 8~72자이며 영문과 숫자를 포함해야 하고, _!@#$%^&*()+-=를 사용할 수 있어요",
    strongpwd: "비밀번호는 8~80자이며 영문, 숫자, 특수문자 1자 이상(_@#$%^&*()+-=)을 포함해야 해요",
    summary_title: "다음 항목을 수정해 주세요"
  },
  charCount: "남은 글자: {count}"
}, be = {
  box: {
    ok: "Confirmar",
    cancel: "Cancelar",
    loading: "Carregando",
    close: "Fechar"
  },
  valid: {
    required: "Obrigatório",
    pattern: "Formato inválido",
    min: "Mínimo de {n} caracteres",
    max: "Máximo de {n} caracteres",
    counter: "{current} de {max} caracteres inseridos",
    range: "Deve estar entre {min} e {max}",
    rangeMin: "Mín. {n}",
    rangeMax: "Máx. {n}",
    match: "Os campos não coincidem",
    compareto: "Não coincide",
    tagMax: "Máximo de {n} itens",
    callback: "Falha na validação",
    validname: "Apenas letras, números, - e _ (3 a 16 caracteres)",
    en: "Apenas letras do alfabeto inglês",
    domain: "Formato de domínio inválido",
    mobile: "Número de celular inválido",
    int: "Informe um número inteiro",
    email: "Endereço de e-mail inválido",
    date: "Formato de data inválido (YYYY-MM-DD)",
    url: "URL inválida",
    pwd: "A senha deve ter de 8 a 72 caracteres e conter letras e números; pode incluir _!@#$%^&*()+-=",
    strongpwd: "A senha deve ter de 8 a 80 caracteres e conter letras, números e pelo menos 1 caractere especial (_@#$%^&*()+-=)",
    summary_title: "Corrija os seguintes problemas"
  },
  charCount: "Restantes: {count}"
}, ve = {
  box: {
    ok: "ОК",
    cancel: "Отмена",
    loading: "Загрузка",
    close: "Закрыть"
  },
  valid: {
    required: "Обязательное поле",
    pattern: "Неверный формат",
    min: "Не менее {n} символов",
    max: "Не более {n} символов",
    counter: "Введено {current} из {max} символов",
    range: "Должно быть от {min} до {max}",
    rangeMin: "Мин. {n}",
    rangeMax: "Макс. {n}",
    match: "Поля не совпадают",
    compareto: "Не совпадает",
    tagMax: "Максимум {n} элементов",
    callback: "Не удалось выполнить проверку",
    validname: "Только буквы, цифры, дефис (-) и подчёркивание (_); длина 3–16 символов",
    en: "Только латинские буквы",
    domain: "Неверный формат домена",
    mobile: "Неверный номер мобильного телефона",
    int: "Требуется целое число",
    email: "Неверный адрес электронной почты",
    date: "Неверный формат даты (YYYY-MM-DD)",
    url: "Неверная ссылка",
    pwd: "Пароль должен содержать от 8 до 72 символов, включать буквы и цифры; допускаются символы _!@#$%^&*()+-=",
    strongpwd: "Пароль должен содержать от 8 до 80 символов, включать буквы, цифры и хотя бы один специальный символ (_@#$%^&*()+-=)",
    summary_title: "Пожалуйста, исправьте следующие ошибки"
  },
  charCount: "Осталось: {count}"
}, we = {
  box: {
    ok: "موافق",
    cancel: "إلغاء",
    loading: "جارٍ التحميل",
    close: "إغلاق"
  },
  valid: {
    required: "مطلوب",
    pattern: "صيغة غير صالحة",
    min: "يلزم {n} أحرف على الأقل",
    max: "الحد الأقصى {n} حرفًا",
    counter: "تم إدخال {current} من أصل {max} حرفًا",
    range: "يجب أن تكون القيمة بين {min} و{max}",
    rangeMin: "الحد الأدنى {n}",
    rangeMax: "الحد الأقصى {n}",
    match: "لا يتطابق الحقلان",
    compareto: "غير متطابق",
    tagMax: "الحد الأقصى {n} عناصر",
    callback: "فشل التحقق",
    validname: "الأحرف والأرقام والشرطة (-) والشرطة السفلية (_) فقط (من 3 إلى 16 حرفًا)",
    en: "أحرف إنجليزية فقط",
    domain: "صيغة النطاق غير صالحة",
    mobile: "رقم الجوال غير صالح",
    int: "يلزم عدد صحيح",
    email: "عنوان البريد الإلكتروني غير صالح",
    date: "صيغة التاريخ غير صالحة (YYYY-MM-DD)",
    url: "رابط غير صالح",
    pwd: "يجب أن تتكون كلمة المرور من 8 إلى 72 حرفًا، وأن تحتوي على أحرف وأرقام، ويمكن أن تتضمن _!@#$%^&*()+-=",
    strongpwd: "يجب أن تتكون كلمة المرور من 8 إلى 80 حرفًا، وأن تحتوي على أحرف وأرقام، ورمز خاص واحد على الأقل (_@#$%^&*()+-=)",
    summary_title: "يُرجى تصحيح المشكلات التالية"
  },
  charCount: "المتبقي: {count}"
};
const X = {
  en: ue,
  "zh-CN": fe,
  "zh-TW": me,
  es: ge,
  ja: he,
  de: pe,
  fr: xe,
  ko: ye,
  pt: be,
  ru: ve,
  ar: we
}, K = Object.keys(X);
function R(r) {
  if (!r) return "";
  const e = String(r).trim().replace(/_/g, "-"), t = e.toLowerCase();
  if (t === "zh" || t.startsWith("zh-"))
    return t.includes("hant") ? "zh-TW" : t.includes("hans") ? "zh-CN" : /(?:^|-)(?:tw|hk|mo)(?:-|$)/.test(t) ? "zh-TW" : "zh-CN";
  const n = t.split("-")[0];
  return K.find((i) => i.toLowerCase() === t || i.toLowerCase() === n) || e;
}
function Ee() {
  if (typeof localStorage < "u") {
    const r = R(localStorage.getItem("uno.lang"));
    if (r && K.includes(r)) return r;
  }
  if (typeof navigator < "u") {
    const r = [navigator.language, ...navigator.languages || []];
    for (const e of r) {
      const t = R(e);
      if (K.includes(t)) return t;
    }
    return "en";
  }
  return "en";
}
function Ce(r, e) {
  const t = e.split(".");
  let n = r;
  for (const o of t) {
    if (!n || typeof n != "object" || !(o in n)) return;
    n = n[o];
  }
  return n;
}
function Z(r, e) {
  for (const t in e)
    e[t] && typeof e[t] == "object" ? (r[t] || Object.assign(r, { [t]: {} }), Z(r[t], e[t])) : Object.assign(r, { [t]: e[t] });
}
function ke(r) {
  const e = r.store("i18n"), t = {
    locale: Ee(),
    messages: X,
    loading: !1,
    // <html lang> 只由服务端模板输出；Uno 只维护组件词条语言，不回写文档属性。
    init() {
    },
    t(n, o = {}) {
      const i = R(this.locale), s = this.messages[i] || this.messages.en, a = Ce(s, n);
      return (typeof a == "string" ? a : n).replace(/\{(\w+)\}/g, (c, u) => u in o ? String(o[u]) : `{${u}}`);
    },
    setLocale(n) {
      const o = R(n);
      if (!this.messages[o]) {
        h.warn(`[Uno] Language ${n} not found. Please load it first using loadLocale().`);
        return;
      }
      this.locale = o, typeof localStorage < "u" && localStorage.setItem("uno.lang", o);
    },
    // Alias for backward compatibility
    setLang(n) {
      this.setLocale(n);
    },
    extend(n) {
      Z(this.messages, n);
    },
    async loadLocale(n, o) {
      const i = R(n);
      if (!this.messages[i]) {
        this.loading = !0;
        try {
          const s = await fetch(o);
          if (!s.ok) throw new Error(`Failed to load ${n}`);
          const a = await s.json();
          this.messages[i] = a, this.setLocale(i);
        } catch (s) {
          h.error(`[Uno] Failed to load locale ${n}:`, s);
        } finally {
          this.loading = !1;
        }
      }
    }
  };
  e ? (h.info("[Uno] Merging into existing i18n store"), e.messages || (e.messages = {}), Z(e.messages, X), e.t || (e.t = t.t), e.setLocale || (e.setLocale = t.setLocale), e.setLang || (e.setLang = t.setLang), e.loadLocale || (e.loadLocale = t.loadLocale)) : r.store("i18n", t), r.directive("t", (n, { expression: o, modifiers: i }, { effect: s }) => {
    const a = () => r.store("i18n");
    s(() => {
      const l = a();
      l.locale;
      const c = o, u = l.t(c);
      i.length > 0 ? i.forEach((g) => {
        n.setAttribute(g, u);
      }) : n.textContent = u;
    });
  });
}
function F(r, e, t) {
  if (!r || typeof r.store != "function") {
    h.error("[Uno] Alpine.store is not available");
    return;
  }
  const n = r.store("uno");
  n ? r.store("uno", {
    ...n,
    [e]: t
  }) : r.store("uno", { [e]: t }), r.store(e) || r.store(e, r.store("uno")?.[e]);
}
const Me = (r) => {
  const e = r || "center", t = {
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
  return t[e] || t.center;
}, Te = (r) => {
  const e = r || "top-right", t = {
    "top-right": "top-0 end-0",
    "top-left": "top-0 start-0",
    "top-center": "top-0 start-0 end-0",
    // Full width for center alignment
    center: "inset-0",
    // Full screen for center
    "center-right": "top-1/2 end-0 -translate-y-1/2",
    "center-left": "top-1/2 start-0 -translate-y-1/2",
    "bottom-right": "bottom-0 end-0",
    "bottom-left": "bottom-0 start-0",
    "bottom-center": "bottom-0 start-0 end-0"
  };
  return t[e] || t["top-right"];
}, Se = (r) => {
  const e = r || "top-right";
  let t = "flex-col", n = "items-end", o = "space-y-2";
  return e.includes("bottom") && (t = "flex-col-reverse", o = "space-y-reverse space-y-2"), e.includes("left") ? n = "items-start" : e.includes("center") ? e === "center" ? n = "items-center justify-center" : n = "items-center" : n = "items-end", { flexDir: t, alignItems: n, spaceY: o };
};
function $e(r, e, t, n = {}) {
  const o = typeof n?.offset == "number" && Number.isFinite(n.offset) && n.offset >= 0 ? n.offset : 8, i = r.getBoundingClientRect(), s = e.getBoundingClientRect(), l = (typeof getComputedStyle == "function" && getComputedStyle(r).direction === "rtl" ? "rtl" : "ltr") === "rtl" ? t === "left" ? "right" : t === "right" ? "left" : t : t, c = s.width, u = s.height, g = i.width, d = i.height, m = { top: i.top, left: i.left }, w = window.innerWidth, y = window.innerHeight, p = 0, x = 0, f = y, v = w, k = ($, H) => {
    if (H) {
      if ($ + c > v)
        return v - c - 10;
      if ($ < x)
        return x + 10;
    } else {
      if ($ + u > f)
        return f - u - 10;
      if ($ < p)
        return p + 10;
    }
    return $;
  };
  let E, C, I = ["top", "bottom", "left", "right"].includes(l) ? l : "bottom";
  switch (l) {
    case "left":
      C = m.left - c - o, C < x && (C = m.left + g + o, I = "right"), E = k(m.top + d / 2 - u / 2, !1);
      break;
    case "top":
      E = m.top - u - o, E < p && (E = m.top + d + o, I = "bottom"), C = m.left + (g - c) / 2, C = k(C, !0);
      break;
    case "right":
      C = m.left + g + o, C + c > v && (C = m.left - c - o, I = "left"), E = k(m.top + d / 2 - u / 2, !1);
      break;
    default:
      E = m.top + d + o, E + u > f && (E = m.top - u - o, I = "top"), C = m.left + (g - c) / 2, C = k(C, !0);
      break;
  }
  const O = 8, T = O / 2, z = ($, H, re) => Math.min(Math.max($, H), Math.max(H, re)), ee = m.left + g / 2, te = m.top + d / 2;
  let P;
  switch (I) {
    case "top":
      P = {
        top: u - T,
        left: z(ee - C - T, 0, c - O)
      };
      break;
    case "left":
      P = {
        top: z(te - E - T, 0, u - O),
        left: c - T
      };
      break;
    case "right":
      P = {
        top: z(te - E - T, 0, u - O),
        left: -T
      };
      break;
    default:
      P = {
        top: -T,
        left: z(ee - C - T, 0, c - O)
      };
      break;
  }
  return { top: E, left: C, position: I, arrow: P };
}
const G = [
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
function J(r, e = {}) {
  const { initialFocus: t = !0 } = e, n = (l) => {
    if (!(l instanceof HTMLElement)) return !1;
    const c = window.getComputedStyle(l);
    return c.display === "none" || c.visibility === "hidden" ? !1 : l.getClientRects().length > 0;
  }, o = () => Array.from(r.querySelectorAll(G)).filter((l) => n(l)), i = () => o()[0], s = (l) => {
    if (l.key !== "Tab") return;
    const c = o();
    if (!c.length) return;
    const u = c[0], g = c[c.length - 1], d = document.activeElement, m = c.some(
      (y) => d.compareDocumentPosition(y) & Node.DOCUMENT_POSITION_FOLLOWING
    ), w = c.some(
      (y) => d.compareDocumentPosition(y) & Node.DOCUMENT_POSITION_PRECEDING
    );
    l.shiftKey ? (d === u || !w) && (l.preventDefault(), g.focus()) : (d === g || !m) && (l.preventDefault(), u.focus());
  }, a = () => {
    if (r.contains(document.activeElement)) return;
    const l = i();
    l && l.focus();
  };
  return t && (typeof requestAnimationFrame == "function" ? requestAnimationFrame(() => requestAnimationFrame(a)) : setTimeout(a, 0)), r.addEventListener("keydown", s), () => {
    r.removeEventListener("keydown", s);
  };
}
const Le = () => (b.options?.ui?.layers?.tooltip || S("TOOLTIP"), `
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
let _ = 0, A = null;
const Ie = (r) => {
  const e = () => {
    _ || typeof requestAnimationFrame == "function" && (_ = requestAnimationFrame(() => {
      _ = 0;
      const n = r.store("uno")?.tooltip;
      n && n.updatePosition();
    }));
  }, t = () => {
    if (A || typeof window > "u" || typeof document > "u") return;
    const n = () => e(), o = () => e(), i = (a) => {
      if (a.key !== "Escape") return;
      const l = r.store("uno")?.tooltip;
      l?.visible && l.hide();
    }, s = (a) => {
      const l = r.store("uno")?.tooltip;
      if (!l?.visible || l.triggerMode !== "click" || l.clickOutside !== !0) return;
      const c = a.target, u = document.getElementById("uno-tooltip-content"), g = l.triggerEl && l.triggerEl.contains(c), d = u && u.contains(c);
      !g && !d && l.hide();
    };
    window.addEventListener("scroll", n, !0), window.addEventListener("resize", o), window.addEventListener("keydown", i), window.addEventListener("click", s), A = () => {
      window.removeEventListener("scroll", n, !0), window.removeEventListener("resize", o), window.removeEventListener("keydown", i), window.removeEventListener("click", s), A = null;
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
    show(n = {}) {
      h.info("Tooltip show called", n);
      const o = b.getConfig("tooltip", n);
      if (h.info("Tooltip resolved config", o), !o.target && typeof o.content == "string" && o.content.startsWith("#") && (o.content.includes(" ") || document.querySelector(o.content) && (o.target = o.content)), this.triggerEl = o.element || this.triggerEl, o.target) {
        const i = document.querySelector(o.target);
        i ? this.content = D(i.innerHTML) : (h.warn(`Tooltip target element not found: ${o.target}`), this.content = D(o.content || ""));
      } else
        this.content = D(o.content || "");
      this.placement = o.placement, this.triggerMode = o.trigger, this.clickOutside = o.clickOutside, this.offset = o.offset, this.arrow = o.arrow !== !1, this.zIndex = b.options?.ui?.layers?.tooltip || S("TOOLTIP"), this.hxUrl = o.hxUrl, this.hxHeaders = o.hxHeaders, this.hxTrigger = o.hxTrigger || "load", this.visible = !0, this.ensureContainer(), t(), this.triggerMode === "click" && this.triggerEl && (this.triggerEl.setAttribute("aria-haspopup", "true"), this.triggerEl.setAttribute("aria-expanded", "true"), this.triggerEl.setAttribute("aria-describedby", "uno-tooltip-content"), this.bindPopupKeyboard()), r.nextTick(() => {
        this.updatePosition(), this.focusPopupIfNeeded();
      });
    },
    hide() {
      const n = this.triggerEl, o = this.triggerMode === "click", i = document.getElementById("uno-tooltip-content"), s = document.activeElement, a = !!i && !!s && i.contains(s);
      this.visible = !1, this.triggerEl = null, this.cleanupTrap && (this.cleanupTrap(), this.cleanupTrap = null), i && this.onPopupKeydown && (i.removeEventListener("keydown", this.onPopupKeydown), this.onPopupKeydown = null), n && this.onTriggerKeydown && (n.removeEventListener("keydown", this.onTriggerKeydown), this.onTriggerKeydown = null), o && n && n.setAttribute && (n.setAttribute("aria-expanded", "false"), n.removeAttribute("aria-describedby")), A && (A(), A = null), _ && (cancelAnimationFrame(_), _ = 0), o && n && typeof n.focus == "function" && (a || s === document.body) && n.focus();
    },
    /**
     * A3 click 弹层键盘绑定：Tab 循环（trapFocus）、role=menu 内容方向键导航、
     * 锚点触发器 Space 合成（Enter 在锚点/按钮上由原生 click 覆盖，按钮 Space 亦原生，
     * 仅锚点需合成以免双重 toggle）。
     */
    bindPopupKeyboard() {
      const n = document.getElementById("uno-tooltip-content");
      n && (this.cleanupTrap = J(n, { initialFocus: !1 }), this.onPopupKeydown = (o) => {
        if (!["ArrowDown", "ArrowUp", "Home", "End"].includes(o.key) || !n.querySelector('[role="menu"]')) return;
        const i = Array.from(n.querySelectorAll(G));
        if (i.length === 0) return;
        o.preventDefault();
        const s = i.indexOf(document.activeElement);
        let a = 0;
        o.key === "ArrowDown" ? a = (s + 1) % i.length : o.key === "ArrowUp" ? a = (s - 1 + i.length) % i.length : o.key === "End" && (a = i.length - 1), i[a] && i[a].focus();
      }, n.addEventListener("keydown", this.onPopupKeydown), this.triggerEl && this.triggerEl.tagName === "A" && (this.onTriggerKeydown = (o) => {
        o.key !== " " && o.key !== "Spacebar" || (o.preventDefault(), this.visible && this.triggerEl === o.currentTarget && this.hide());
      }, this.triggerEl.addEventListener("keydown", this.onTriggerKeydown)));
    },
    /**
     * A3 焦点移入弹层首个可聚焦元素，键盘用户落在刚打开的披露内容内（2.4.3）；
     * 弹层无可聚焦内容时为 no-op（焦点留在触发器）。
     */
    focusPopupIfNeeded() {
      if (this.triggerMode !== "click") return;
      const n = document.getElementById("uno-tooltip-content");
      if (!n || n.contains(document.activeElement)) return;
      const o = n.querySelector(G);
      o && o.focus();
    },
    cancelHide() {
      this.hideTimeout && (clearTimeout(this.hideTimeout), this.hideTimeout = null);
    },
    scheduleHide(n = 300) {
      this.cancelHide(), this.hideTimeout = setTimeout(() => {
        this.hide();
      }, n);
    },
    updatePosition() {
      if (!this.visible || !this.triggerEl) return;
      const n = document.getElementById("uno-tooltip-content");
      if (!n) return;
      const o = $e(this.triggerEl, n, this.placement, {
        offset: this.offset
      });
      this.top = o.top, this.left = o.left, this.placement = o.position || this.placement, o.arrow && (this.arrowPosition = o.arrow);
    },
    ensureContainer() {
      const n = "uno-tooltip-manager";
      if (document.getElementById(n)) return;
      const o = b.options?.ui?.layers?.tooltip || S("TOOLTIP"), i = Le();
      j.ensureContainer(n, i, {
        "x-data": "{}",
        style: `position: absolute; top: 0; left: 0; width: 0; height: 0; z-index: ${o};`
      });
      const s = document.getElementById(n);
      s && r && !s._x_dataStack && r.initTree(s);
    }
  };
}, Ue = (r) => {
  const e = (t, { expression: n, modifiers: o }, { evaluateLater: i, effect: s, cleanup: a }) => {
    const l = n ? i(n) : null, c = (f) => (o.length > 0 && h.warn(
      "[Uno] x-uno-tooltip modifiers are deprecated in v2.0. Use configuration object instead."
    ), { ...typeof f == "object" && f !== null ? f : f ? { content: f } : {}, element: t }), u = (f) => {
      const v = r.store("uno")?.tooltip;
      v && v.show(c(f));
    }, g = (f) => {
      const v = r.store("uno")?.tooltip;
      v && (v.visible && v.triggerEl === t ? v.hide() : u(f));
    };
    let d = null;
    l && s(() => {
      l((f) => {
        d = f;
      });
    });
    const m = () => {
      if (n && d === null) return;
      const f = c(d);
      b.getConfig("tooltip", f).trigger === "hover" && u(d);
    }, w = () => {
      const f = r.store("uno")?.tooltip;
      f && f.triggerMode === "hover" && f.triggerEl === t && f.scheduleHide();
    }, y = (f) => {
      const v = c(d);
      b.getConfig("tooltip", v).trigger === "click" && (f.preventDefault(), g(d));
    }, p = () => {
      const f = c(d);
      b.getConfig("tooltip", f).trigger === "focus" && u(d);
    }, x = () => {
      const f = r.store("uno")?.tooltip;
      f && f.triggerMode === "focus" && f.triggerEl === t && f.hide();
    };
    t.addEventListener("mouseenter", m), t.addEventListener("mouseleave", w), t.addEventListener("click", y), t.addEventListener("focus", p), t.addEventListener("blur", x), a(() => {
      t.removeEventListener("mouseenter", m), t.removeEventListener("mouseleave", w), t.removeEventListener("click", y), t.removeEventListener("focus", p), t.removeEventListener("blur", x);
      const f = r.store("uno")?.tooltip;
      f && f.triggerEl === t && f.hide();
    });
  };
  r.directive("uno-tooltip", e);
};
const Be = (r, e = {}) => {
  const t = {
    placement: "top",
    trigger: "hover",
    delay: 300,
    offset: 8,
    arrow: !0,
    clickOutside: !0
  };
  b.registerPluginDefaults("tooltip", {
    ...t,
    ...e.plugins?.tooltip,
    ...e.tooltip
  });
  const n = Ie(r);
  F(r, "tooltip", n), Ue(r);
}, _e = () => `
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
`, Ae = () => `
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
  `, Oe = (r) => ({
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
  onContentLoaded(e, t) {
    this._onHtmxEnd(t);
    const n = e.innerHTML || e.textContent || "", o = this._findNotificationById(t.id);
    o && (o.content = n), t.onContentLoad && t.onContentLoad(e), window.Alpine && window.Alpine.initTree && window.Alpine.initTree(e);
  },
  // Helper method to find notification by ID
  _findNotificationById(e) {
    for (const t of Object.keys(this.groups)) {
      const o = this.groups[t].find((i) => i.id === e);
      if (o) return o;
    }
    return null;
  },
  onHtmxError(e, t) {
    this._onHtmxEnd(t), h.error("HTMX Error:", e), this.error("Content load failed"), t.onError && t.onError(e);
  },
  _onHtmxBegin(e) {
    !e || e.__unoHtmxActive || (e.__unoHtmxActive = !0, this.beginProgress());
  },
  _onHtmxEnd(e) {
    !e || !e.__unoHtmxActive || (e.__unoHtmxActive = !1, this.endProgress());
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
      zIndex: b.options?.ui?.layers?.notify || S("NOTIFY")
    };
  },
  getContainerClass(e) {
    const t = b.getStyle("notify", "container"), n = Te(e), o = Se(e);
    return [t, n, o.flexDir, o.alignItems].join(" ");
  },
  ensureContainer() {
    const e = "uno-notify-container", t = _e(), n = j.ensureContainer(e, t, {
      "x-data": "{}"
      // Independent Alpine component
    });
    n && r.initTree && r.initTree(n);
    const o = "uno-global-progress", i = Ae(), s = j.ensureContainer(o, i, {
      "x-data": "{}"
    });
    s && r.initTree && r.initTree(s);
  },
  show(e, t = {}) {
    this.ensureContainer();
    const n = b.getConfig("notify", t), o = this.nextId++;
    let i = e;
    r.store("i18n") && typeof r.store("i18n").t == "function" && (i = r.store("i18n").t(e));
    const s = n.duration !== void 0 ? n.duration : 3e3;
    h.info("[Notify] Show:", { id: o, message: i, duration: s, config: n });
    const a = {
      id: o,
      message: i,
      isVisible: !0,
      // For x-show
      closable: !0,
      duration: s,
      content: n.hxUrl ? "" : void 0,
      // Initialize content for HTMX-loaded content
      ...n
    }, l = n.position || "top-right";
    return this.groups[l] ? this.groups[l].push(a) : (this.groups["top-right"] || (this.groups["top-right"] = []), this.groups["top-right"].push(a)), s > 0 && setTimeout(() => {
      this.remove(o);
    }, s), o;
  },
  remove(e) {
    Object.keys(this.groups).forEach((t) => {
      const n = this.groups[t].findIndex((o) => o.id === e);
      n !== -1 && (this.groups[t][n].isVisible = !1, setTimeout(() => {
        const o = this.groups[t].findIndex((i) => i.id === e);
        o !== -1 && this.groups[t].splice(o, 1);
      }, 300));
    });
  },
  // 3. Proxy Methods
  success(e, t = {}) {
    return this.show(e, { type: "success", ...t });
  },
  error(e, t = {}) {
    return this.show(e, { type: "error", ...t });
  },
  warning(e, t = {}) {
    return this.show(e, { type: "warning", ...t });
  },
  info(e, t = {}) {
    return this.show(e, { type: "info", ...t });
  }
}), Pe = (r) => {
  r.directive("uno-notify", (e, { expression: t }, { evaluateLater: n }) => {
    if (!t) {
      h.warn("[Uno] x-uno-notify directive requires a configuration object or message string.");
      return;
    }
    const o = n(t), i = () => {
      const s = r.store("uno")?.notify;
      s && o((a) => {
        let l = "", c = {};
        if (typeof a == "string")
          l = a;
        else if (typeof a == "object" && a !== null)
          l = a.message || "", c = a;
        else {
          h.warn("[Uno] x-uno-notify directive value must be a string or object", a);
          return;
        }
        s.show(l, c);
      });
    };
    e.addEventListener("click", i);
  });
};
const De = (r, e = {}) => {
  const t = {
    duration: 3e3,
    position: "top-right"
  };
  b.registerPluginDefaults("notify", {
    ...t,
    ...e.plugins?.notify,
    ...e.notify
  });
  const n = Oe(r);
  F(r, "notify", n), Pe(r);
}, Re = () => `
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
           :style="'z-index:' + (${b.options?.ui?.layers?.box || S("BOX")} + index * 10)"
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
             :class="[Uno.getStyle('box', 'container'), item.widthClass, item.heightClass || (item.fill && !item.height ? Uno.getStyle('box', 'containerFill') : ''), item.containerClass, 'focus:outline-none']"
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

          <!-- Content Area（fill 模式下至加载容器逐层 flex 传递剩余高度） -->
          <div :class="[Uno.getStyle('box', 'content'), item.fill ? Uno.getStyle('box', 'contentFill') : '']">
            <h3 :id="'uno-box-title-' + item.id"
                :class="[Uno.getStyle('box', 'title'), 'focus:outline-none']"
                tabindex="-1"
                x-text="item.title"
                x-show="item.title"></h3>
            <div class="mt-2" :class="item.fill ? Uno.getStyle('box', 'bodyFill') : ''" :id="'uno-box-content-' + item.id">
              <template x-if="!item.hxUrl">
                <div :class="item.fill ? Uno.getStyle('box', 'bodyFill') : ''" x-html="item.content"></div>
              </template>
              <template x-if="item.hxUrl">
                <div x-init="htmx.process($el); $el.addEventListener('htmx:afterSwap', () => $store.uno.box.onHtmxContentLoaded($el, item))"
                     :class="item.fill ? Uno.getStyle('box', 'bodyFill') : ''"
                     :hx-get="item.hxMethod === 'GET' ? item.hxUrl : undefined"
                     :hx-post="item.hxMethod === 'POST' ? item.hxUrl : undefined"
                     :hx-trigger="item.hxTrigger"
                     hx-swap="innerHTML"
                     :hx-headers="item.hxHeaders ? JSON.stringify(item.hxHeaders) : null">
                  <div :class="[Uno.getStyle('box', 'loadingWrap'), item.fill ? 'flex-1' : '']">
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
`, Fe = (r) => ({
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
  alert(e, t = {}) {
    const n = typeof t == "string" ? { title: t } : t;
    return this.open("dialog", {
      content: e,
      showOkButton: !0,
      ...n
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
  confirm(e, t = {}) {
    return new Promise((n) => {
      const { onOk: o, onCancel: i, onClose: s, ...a } = t;
      this.open("confirm", {
        ...a,
        content: e,
        onOk: () => {
          const l = typeof o == "function" ? o() : void 0;
          if (l === !1) return !1;
          if (l instanceof Promise)
            return l.then((c) => {
              if (c === !1) return !1;
              n(!0);
            });
          n(!0);
        },
        onCancel: () => {
          typeof i == "function" && i(), n(!1);
        },
        onClose: () => {
          typeof s == "function" && s(), n(!1);
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
  prompt(e, t = "", n = {}) {
    return new Promise((o) => {
      const { onOk: i, onCancel: s, ...a } = n, l = "prompt-" + Math.random().toString(36).slice(2, 11), c = b.getStyle("box", "promptWrapper"), u = b.getStyle("box", "promptText"), g = b.getStyle("box", "promptInput"), d = b.dom.sanitize(e), m = b.dom.sanitize(t), w = `
        <div class="${c}">
          <p class="${u}">${d}</p>
          <input id="${l}" type="text" class="${g}" value="${m}">
        </div>
      `;
      this.open("confirm", {
        title: a.title || "Prompt",
        ...a,
        content: w,
        _isTrustedPrompt: !0,
        // Internal flag for prompt content
        onOk: () => {
          const y = typeof i == "function" ? i() : void 0;
          if (y === !1) return !1;
          if (y instanceof Promise)
            return y.then((x) => {
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
  dialog(e = {}) {
    return this.open("dialog", e);
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
  open(e, t = {}) {
    let n = "dialog", o = t;
    typeof e == "object" ? (o = e, n = o.type || "dialog") : n = e;
    const i = b.getConfig("box", o), s = i.boxId || Math.random().toString(36).slice(2, 10);
    let a = i.showOkButton, l = i.showCancelButton;
    n === "confirm" ? (a = a !== !1, l = l !== !1) : n === "dialog" && (a = a || !1, l = l || !1);
    const c = {
      id: s,
      type: n,
      title: i.title || "",
      content: i._isTrustedPrompt ? i.content || "" : b.dom.sanitize(i.content || ""),
      // Security: Always sanitize content (trustedHTML removed)
      // Exception: Internal prompt content is pre-escaped
      _isTrustedPrompt: i._isTrustedPrompt || !1,
      width: i.width !== void 0 ? i.width : i.widthClass ? void 0 : i.defaultWidth || 400,
      height: i.height !== void 0 ? i.height : (i.heightClass, void 0),
      widthClass: i.widthClass,
      heightClass: i.heightClass,
      containerClass: i.containerClass,
      fill: i.fill === !0,
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
        const m = document.getElementById("uno-box-outer-container" + this.id), w = document.activeElement;
        (m && w && m.contains(w) || w === document.body) && d.focus();
      },
      /**
       * Close this box dialog
       */
      close() {
        this.onBeforeClose && this.onBeforeClose(), this.isOpen = !1;
        const d = r.store("uno")?.box;
        d && d.syncUnderlayInert(), this.restoreFocus();
        const m = typeof document < "u" ? document.getElementById("uno-box-outer-container" + this.id) : null, w = () => {
          const x = r.store("uno")?.box;
          x && x.remove(this.id), this.onClose && this.onClose();
        };
        if (!m) {
          w();
          return;
        }
        let y = !1;
        const p = (x) => {
          x && x.target !== m || y || (y = !0, m.removeEventListener("transitionend", p), w());
        };
        m.addEventListener("transitionend", p), setTimeout(p, 350);
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
  async handleExtraButton(e, t) {
    if (!e || !t?.onClick) return;
    const n = t.onClick(e);
    n instanceof Promise ? await n !== !1 && e.close() : n !== !1 && e.close();
  },
  /**
   * Handle OK button click
   * Supports synchronous and asynchronous onOk callbacks
   * @param {object} item - Box item
   */
  async handleOk(e) {
    if (e)
      if (e.onOk) {
        const t = e.onOk();
        t instanceof Promise ? await t !== !1 && e.close() : t !== !1 && e.close();
      } else
        e.close();
  },
  /**
   * Close a specific box by ID
   * @param {string} id - Box ID
   */
  close(e) {
    const t = this.items.find((n) => n.id === e);
    t && t.close();
  },
  /**
   * Remove a box from the items array
   * @param {string} id - Box ID
   */
  remove(e) {
    const t = this.items.findIndex((n) => n.id === e);
    t !== -1 && (this.items.splice(t, 1), this.syncUnderlayInert());
  },
  /**
   * B5 多层弹窗：仅顶层可交互/可感知，下层 inert + aria-hidden，
   * 屏幕阅读器与键盘不穿透底层（不改现有 scheduleAutoFocus 契约）。
   */
  syncUnderlayInert() {
    if (typeof document > "u") return;
    const e = this.items.filter((t) => t.isOpen);
    e.forEach((t, n) => {
      const o = document.getElementById("uno-box-outer-container" + t.id);
      if (!o) return;
      const i = n < e.length - 1;
      o.toggleAttribute("inert", i), i ? o.setAttribute("aria-hidden", "true") : o.removeAttribute("aria-hidden");
    });
  },
  /**
   * Close all open boxes
   */
  closeAll() {
    this.items.forEach((e) => {
      e.onBeforeClose && e.onBeforeClose(), e.onClose && e.onClose();
    }), this.items = [];
  },
  /**
   * Get position CSS classes for 9-point grid
   * @param {string} position - Position name
   * @returns {string} CSS classes
   */
  getPositionClasses(e) {
    return Me(e);
  },
  /**
   * Handle HTMX content load completion
   * @param {HTMLElement} el - Content element
   * @param {object} item - Box item
   */
  onHtmxContentLoaded(e, t) {
    t.onContentLoad && typeof t.onContentLoad == "function" && t.onContentLoad(e), window.Alpine && window.Alpine.initTree && window.Alpine.initTree(e), this.scheduleAutoFocus(t, window.Alpine);
  },
  scheduleAutoFocus(e, t) {
    const n = "uno-box-container" + e.id, o = () => {
      const s = document.getElementById(n);
      s && (e.autoFocus && this.focusFirstInput(s) || this.focusNeutralTarget(s, e));
    }, i = (s) => {
      typeof requestAnimationFrame == "function" ? requestAnimationFrame(() => requestAnimationFrame(s)) : setTimeout(s, 0);
    };
    t && typeof t.nextTick == "function" ? t.nextTick(() => i(o)) : i(o);
  },
  /**
   * Focus the first input element in the container
   * @param {HTMLElement} container - Container element
   * @returns {boolean} True if an element was focused
   */
  focusFirstInput(e) {
    if (!e) return !1;
    const t = [
      'input:not([type="hidden"]):not([type="button"]):not([type="submit"]):not([type="reset"]):not([type="image"]):not([disabled]):not([readonly])',
      "textarea:not([disabled]):not([readonly])",
      "select:not([disabled]):not([readonly])"
    ].join(","), n = e.querySelectorAll(t), o = (i) => {
      for (let s = i; s && s.nodeType === 1; s = s.parentElement) {
        const a = window.getComputedStyle(s);
        if (a.display === "none" || a.visibility === "hidden") return !1;
      }
      return !0;
    };
    for (const i of n)
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
  focusNeutralTarget(e, t) {
    if (!e || e.contains(document.activeElement)) return !1;
    const n = t && t.title ? document.getElementById("uno-box-title-" + t.id) : null, o = n && e.contains(n) ? n : e;
    return o.hasAttribute("tabindex") || o.setAttribute("tabindex", "-1"), o.focus({ preventScroll: !0 }), !0;
  },
  /**
   * Ensure the box container exists in the DOM
   */
  ensureContainer() {
    const e = "uno-box-manager", t = typeof document < "u" && document.getElementById(e);
    if (t)
      return r && typeof r.initTree == "function" && !t._x_dataStack && r.initTree(t), t;
    const o = {
      "x-data": "{}",
      class: "relative",
      style: `z-index: ${b.options?.ui?.layers?.box || S("BOX")};`,
      "x-show": "$store.uno.box.items.length > 0",
      "x-cloak": ""
      // Hide until Alpine is ready
    }, i = Re(), s = b.dom.ensureContainer(e, i, o);
    return s && r && typeof r.initTree == "function" && r.initTree(s), s;
  }
}), ze = (r) => {
  r.directive("uno-trap", (e, { expression: t }, { cleanup: n, evaluateLater: o }) => {
    if (t)
      o(t)((s) => {
        const a = J(e, s || {});
        n(a);
      });
    else {
      const i = J(e);
      n(i);
    }
  }), r.directive("uno-box", (e, { expression: t }, { evaluateLater: n, cleanup: o }) => {
    if (!t) {
      h.warn("[Uno] x-uno-box directive requires a configuration object.");
      return;
    }
    const i = n(t), s = (a) => {
      a.preventDefault();
      const l = r.store("uno")?.box;
      if (!l) {
        h.warn("[Uno] Box store not found. Is Uno initialized?");
        return;
      }
      i((c) => {
        c && typeof c == "object" ? l.open(c) : h.warn("[Uno] x-uno-box directive value must be an object", c);
      });
    };
    e.addEventListener("click", s), o(() => e.removeEventListener("click", s));
  });
};
const He = (r, e = {}) => {
  const t = {
    defaultWidth: 400,
    defaultPosition: "center",
    okText: "OK",
    cancelText: "Cancel",
    closeOnEsc: !1,
    closeOnOverlayClick: !1
  };
  b.registerPluginDefaults("box", { ...t, ...e.plugins?.box, ...e.box });
  const n = Fe(r);
  F(r, "box", n), ze(r);
}, oe = {
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
let qe = 0;
const je = (r) => ({
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
  addRules(e, t, n = {}) {
    const o = this.normalizeRules(t), i = b.getConfig("valid", n);
    this.rules[e] = o.map((s) => ({
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
  normalizeRules(e) {
    return e ? typeof e == "string" ? [{ type: e }] : typeof e == "object" && !Array.isArray(e) ? [e] : Array.isArray(e) ? e.map((t) => typeof t == "string" ? { type: t } : t) : [] : [];
  },
  /**
   * Registers a custom validator function.
   *
   * @param {string} name - The name of the validator
   * @param {Function} validatorFn - The validator function.
   */
  registerValidator(e, t) {
    typeof t == "function" ? this.customValidators[e] = t : h.error(`[Valid] Invalid validator function for ${e}`);
  },
  /**
   * Validates a single field against its rules.
   *
   * @param {string} fieldName - The field identifier
   * @param {string|number} value - The value to validate
   * @param {string} [eventType='input'] - The event triggering validation ('input', 'blur', 'submit')
   * @returns {Promise<boolean>|boolean} True if valid (can be async)
   */
  async validate(e, t, n = "input") {
    const o = this.rules[e];
    if (!o || o.length === 0) return !0;
    delete this.errors[e], this.syncAriaInvalid(e), this.timers[e] && (clearTimeout(this.timers[e]), delete this.timers[e]);
    const i = t == null ? "" : String(t), s = i.trim() === "";
    for (const a of o) {
      if (n === "input" && (a.trigger === "blur" || a.realtime === !1))
        continue;
      const l = a.type === "required" || a.required === !0;
      if (!(s && !l && a.emptyQuit !== !1)) {
        if (a.debounce > 0 && n === "input")
          return this.validating[e] = !0, new Promise((c) => {
            this.timers[e] = setTimeout(async () => {
              const u = await this.executeRule(e, t, a, i);
              delete this.validating[e], c(!!u);
            }, a.debounce);
          });
        this.validating[e] = !0;
        try {
          if (!await this.executeRule(e, t, a, i))
            return !1;
        } finally {
          delete this.validating[e];
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
  async executeRule(e, t, n, o) {
    let i = !0, s = "";
    const a = n.type;
    if (a === "required" || n.required)
      o.trim() === "" && (i = !1, s = n.msg || this.getI18nMsg("valid.required") || "This field is required");
    else if (oe[a])
      oe[a].test(o) || (i = !1, s = n.msg || this.getI18nMsg(`valid.${a}`) || `Invalid ${a} format`);
    else if (a === "length")
      n.min !== void 0 && o.length < n.min ? (i = !1, s = n.msg || this.getI18nMsg("valid.min", { n: n.min }) || `Min length is ${n.min}`) : n.max !== void 0 && o.length > n.max && (i = !1, s = n.msg || this.getI18nMsg("valid.max", { n: n.max }) || `Max length is ${n.max}`);
    else if (a === "range") {
      const l = parseFloat(o);
      isNaN(l) ? (i = !1, s = n.msg || this.getI18nMsg("valid.number") || "Must be a number") : n.min !== void 0 && l < n.min ? (i = !1, s = n.msg || this.getI18nMsg("valid.rangeMin", { n: n.min }) || `Min value is ${n.min}`) : n.max !== void 0 && l > n.max && (i = !1, s = n.msg || this.getI18nMsg("valid.rangeMax", { n: n.max }) || `Max value is ${n.max}`);
    } else if (a === "pattern" || n.pattern)
      (n.pattern instanceof RegExp ? n.pattern : new RegExp(n.pattern)).test(o) || (i = !1, s = n.msg || this.getI18nMsg("valid.pattern") || "Invalid format");
    else if (a === "match" || n.match) {
      const l = n.match || n.selector || n.target;
      if (!l)
        return !0;
      const c = document.querySelector(l.startsWith("#") ? l : `[name="${l}"]`) || document.getElementById(l);
      c && o !== c.value && (i = !1, s = n.msg || this.getI18nMsg("valid.match") || "Fields do not match");
    } else if (a === "custom" || a === "validator" || n.fn || this.customValidators[a]) {
      let l = n.fn || n.name;
      if (typeof l == "string" ? l = this.customValidators[l] : !l && this.customValidators[a] && (l = this.customValidators[a]), typeof l == "function")
        try {
          const c = await l(t, e, n);
          c === !1 ? (i = !1, s = n.msg || "Validation failed") : typeof c == "string" ? (i = !1, s = c) : typeof c == "object" && c.valid === !1 && (i = !1, s = c.message || n.msg || "Validation failed");
        } catch (c) {
          h.error(`[Valid] Error in custom validator for ${e}`, c), i = !1, s = "Validation error";
        }
      else
        h.warn(`[Valid] Custom validator not found: ${a || n.fn}`);
    }
    return i || (this.errors[e] = s, this.syncAriaInvalid(e)), i;
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
  async validateForm(e) {
    const t = e.querySelectorAll("input, select, textarea");
    let n = !0, o = null;
    const i = [];
    this.removeErrorSummary(e);
    for (const s of t) {
      if (s.disabled) continue;
      const a = s.name || s.id;
      this.rules[a] && (await this.validate(a, s.value, "submit") || (n = !1, o || (o = s), i.push({ input: s, fieldName: a })));
    }
    if (!n && o) {
      const s = b.getConfig("valid") || {};
      if (i.length >= 2) {
        const a = this.renderErrorSummary(e, i);
        a && s.scrollOnError !== !1 && (a.scrollIntoView({ behavior: "smooth", block: "center" }), a.focus({ preventScroll: !0 }));
      } else s.scrollOnError !== !1 && (o.scrollIntoView({ behavior: "smooth", block: "center" }), o.focus({ preventScroll: !0 }));
    }
    return n;
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
  renderErrorSummary(e, t) {
    if (typeof document > "u" || !e) return null;
    const n = document.createElement("div");
    n.className = "ui-error-summary", n.dataset.validSummary = "true", n.setAttribute("role", "alert"), n.setAttribute("tabindex", "-1");
    const o = `valid-error-summary-title-${++qe}`;
    n.setAttribute("aria-labelledby", o);
    const i = document.createElement("p");
    i.className = "ui-error-summary-title", i.id = o, i.textContent = this.getI18nMsg("valid.summary_title") || "Please fix the following problems", n.appendChild(i);
    const s = document.createElement("ul");
    return s.className = "ui-error-summary-list", t.forEach(({ input: a, fieldName: l }) => {
      const c = a.id || `${l}-error`, u = document.createElement("a");
      u.className = "ui-error-summary-link", u.setAttribute("href", `#${c}`), u.textContent = this.errors[l] || "", u.addEventListener("click", (d) => {
        d.preventDefault();
        const m = document.getElementById(c);
        m && (m.focus({ preventScroll: !0 }), m.scrollIntoView({ behavior: "smooth", block: "center" }));
      });
      const g = document.createElement("li");
      g.appendChild(u), s.appendChild(g);
    }), n.appendChild(s), e.prepend(n), n;
  },
  /**
   * Removes a previously rendered error summary (matched by the
   * `data-valid-summary` flag written in renderErrorSummary).
   *
   * @param {HTMLFormElement} form - The form to clean.
   */
  removeErrorSummary(e) {
    if (!e || typeof e.querySelector != "function") return;
    const t = e.querySelector("[data-valid-summary]");
    t && t.remove();
  },
  showError(e, t) {
    this.errors[e] = t, this.syncAriaInvalid(e);
  },
  hideError(e) {
    delete this.errors[e], this.syncAriaInvalid(e);
  },
  clearErrors() {
    const e = Object.keys(this.errors);
    this.errors = {}, e.forEach((t) => this.syncAriaInvalid(t));
  },
  showErrors(e) {
    Object.assign(this.errors, e), Object.keys(e || {}).forEach((t) => this.syncAriaInvalid(t));
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
  getFieldElement(e) {
    return typeof document > "u" || !e ? null : document.querySelector(`[name="${e}"]`) || document.getElementById(e);
  },
  /**
   * Mirrors the field's error state onto its `aria-invalid` attribute so screen
   * readers announce the control as invalid (WCAG 3.3.1). Called from every
   * `errors` mutation site so client-side validation and 422 server errors stay
   * in sync with the exposed state.
   *
   * @param {string} fieldName - The field identifier.
   */
  syncAriaInvalid(e) {
    const t = this.getFieldElement(e);
    t && (this.errors[e] ? t.setAttribute("aria-invalid", "true") : t.removeAttribute("aria-invalid"));
  },
  getI18nMsg(e, t) {
    const n = r.store("i18n");
    return n && typeof n.t == "function" ? n.t(e, t) : null;
  }
}), Q = "data-uno-counter-wrapper", Ve = (r) => r === "flex" || r === "inline-flex" || r === "grid" || r === "inline-grid", Ye = (r) => {
  const e = r.parentElement;
  if (!e) return null;
  if (e.hasAttribute(Q)) return e;
  const t = window.getComputedStyle(e);
  if (Ve(t.display))
    return t.position === "static" && e.tagName !== "FORM" && e.tagName !== "BODY" && (e.style.position = "relative"), e;
  const n = r.ownerDocument.createElement("div");
  n.setAttribute(Q, ""), n.style.position = "relative";
  const o = window.getComputedStyle(r).display;
  return (o === "inline" || o === "inline-block") && (n.style.display = "inline-block"), e.insertBefore(n, r), n.appendChild(r), n;
}, We = (r, e) => {
  if (!e || !e.hasAttribute(Q)) return;
  const t = e.parentNode;
  t && (t.insertBefore(r, e), t.removeChild(e));
}, Ne = (r) => `${r}-error`, ie = (r, e) => {
  r.id || (r.id = Ne(e)), r.hasAttribute("role") || r.setAttribute("role", "alert");
}, W = (r, e) => {
  if (!e.id) return;
  const t = (r.getAttribute("aria-describedby") || "").split(/\s+/).filter(Boolean);
  t.includes(e.id) || t.push(e.id), r.setAttribute("aria-describedby", t.join(" "));
}, Xe = (r) => {
  r.directive("uno-valid", (e, { expression: t }, { evaluateLater: n, cleanup: o }) => {
    if (e.tagName === "FORM") return;
    const i = e.name || e.id;
    if (!i) {
      h.warn("[Uno Valid] Element has no name or id:", e);
      return;
    }
    const s = Uno.getConfig("valid"), a = Uno.options.ui?.valid || {};
    if (s.autoError !== !1) {
      const d = `$store.uno.valid.errors['${i}']`, m = `[x-text="${d.replace(/"/g, '\\"')}"]`, w = e.form || e.closest("form"), y = w ? w.querySelector(m) : document.querySelector(m);
      if (y)
        ie(y, i), W(e, y);
      else {
        let p;
        if (typeof s.createErrorElement == "function")
          p = s.createErrorElement(i, d);
        else {
          p = document.createElement("div"), p.setAttribute("x-text", d), p.setAttribute("x-show", d);
          const x = s.errorClass || a.helperText || "text-xs text-red-500 mt-1";
          p.className = x;
        }
        p && (ie(p, i), e.after(p), r.initTree(p), W(e, p));
      }
    }
    const l = t ? n(t) : () => {
    }, c = (d) => {
      r.store("uno").valid.addRules(i, d);
    };
    t ? l((d) => {
      c(d);
      const m = r.store("uno").valid;
      if (m && m.normalizeRules) {
        const y = m.normalizeRules(d).find(
          (p) => p.type === "length" || p.maxLength !== void 0 || p.max !== void 0
        );
        if (y) {
          const p = y.max !== void 0 ? y.max : y.maxLength;
          if (y.stopOnMax === !0 && p !== void 0 && e.setAttribute("maxlength", p), y.showCounter === !0 && p !== void 0) {
            const x = `${i}-counter`;
            let f = document.getElementById(x);
            const v = Ye(e);
            f || (f = document.createElement("span"), f.id = x, f.className = "absolute bottom-1 end-2 text-xs text-gray-600 dark:text-gray-400 pointer-events-none transition-colors duration-200"), f.parentElement !== v && e.after(f), W(e, f);
            const k = () => {
              const E = e.value.length;
              f.textContent = `${E}/${p}`;
              const C = m.getI18nMsg ? m.getI18nMsg("valid.counter", { current: E, max: p }) : null;
              f.setAttribute("aria-label", C || `${E}/${p}`), E > p ? (f.classList.add("text-red-600", "dark:text-red-400"), f.classList.remove("text-gray-600", "dark:text-gray-400")) : (f.classList.remove("text-red-600", "dark:text-red-400"), f.classList.add("text-gray-600", "dark:text-gray-400"));
            };
            k(), e.addEventListener("input", k), o(() => {
              e.removeEventListener("input", k), f && f.parentNode && f.parentNode.removeChild(f), We(e, v);
            });
          }
        }
      }
    }) : c([]);
    const u = () => {
      r.store("uno").valid.validate(i, e.value, "input");
    }, g = () => {
      r.store("uno").valid.validate(i, e.value, "blur");
    };
    e.addEventListener("input", u), e.addEventListener("blur", g), o(() => {
      e.removeEventListener("input", u), e.removeEventListener("blur", g);
    });
  }), typeof document < "u" && (document.body.addEventListener("htmx:confirm", async (e) => {
    const t = e.detail.elt, n = t.tagName === "FORM" ? t : t.closest("form");
    if (n && (n.hasAttribute("x-uno-valid") || n.querySelector("[x-uno-valid]"))) {
      e.preventDefault();
      try {
        await r.store("uno").valid.validateForm(n) && e.detail.issueRequest(!0);
      } catch (o) {
        h.error("[Valid] Form validation failed unexpectedly, request blocked", o);
      }
    }
  }), document.body.addEventListener("htmx:afterRequest", (e) => {
    const t = N(e);
    if (t && t.status === 422)
      try {
        const n = JSON.parse(t.text);
        n.errors && (r.store("uno").valid.showErrors(n.errors), h.info("[Valid] Server validation errors applied", n.errors));
      } catch {
      }
  }));
};
const Ke = (r, e = {}) => {
  const t = {
    realtime: !0,
    showError: !0,
    scrollOnError: !0
  };
  b.registerPluginDefaults("valid", {
    ...t,
    ...e.plugins?.valid,
    ...e.valid
  });
  const n = je(r);
  F(r, "valid", n), Xe(r);
}, Ze = (r) => ({
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
    const e = () => {
      let t = window.location.hash.slice(1);
      if (!t && window.location.href.includes("#") && (t = window.location.href.split("#")[1]), h.info(`[Tab] checkHash: ${t}`), t && this.registry[t]) {
        const { group: n, name: o } = this.registry[t];
        h.info(`[Tab] checkHash match: ${o}`), this.activate({ group: n, name: o, updateHash: !1 });
      }
    };
    e(), window.addEventListener("hashchange", e), window.addEventListener("load", e);
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
  activate(e = {}) {
    this.ensureInit();
    const t = b.getConfig("tab", {
      group: "default",
      hash: !0,
      ...e
    });
    t.name && (h.info(
      `[Tab] Activating ${t.name} in ${t.group}. Hash update: ${t.updateHash !== !1}`
    ), this.registry[t.name] || (this.registry[t.name] = { group: t.group, name: t.name }), this.groups[t.group] = t.name, typeof document < "u" && document.querySelectorAll('[role="tabpanel"][data-tab-group="' + t.group + '"]').forEach((n) => n.setAttribute("aria-labelledby", "tab-" + t.name)), t.hash && t.updateHash !== !1 && window.location.hash.slice(1) !== t.name && window.history.replaceState(null, null, `#${t.name}`), t.updateHash === !1 && setTimeout(() => {
      const n = document.getElementById(`tab-${t.name}`);
      n && n.hasAttribute("hx-get") && (h.info(`[Tab] Triggering HTMX for ${t.name}`), window.htmx && window.htmx.trigger(n, "click"));
    }, 0));
  },
  /**
   * Check if a tab is active
   * @param {string} name - Tab name
   * @param {string} group - Group name
   * @returns {boolean}
   */
  isActive(e, t = "default") {
    return this.groups[t] === e;
  },
  /**
   * Register a tab (helper for directives or manual setup)
   * @param {string} name
   * @param {string} group
   */
  register(e, t = "default") {
    this.ensureInit(), h.info(`[Tab] Registering ${e} in ${t}. Hash: ${window.location.hash}`), this.registry[e] = { name: e, group: t }, this.order[t] || (this.order[t] = []), this.order[t].includes(e) || this.order[t].push(e), this.groups[t] || (this.groups[t] = e);
    let n = window.location.hash.slice(1);
    !n && window.location.href.includes("#") && (n = window.location.href.split("#")[1]), n === e && (h.info(`[Tab] Hash match for ${e}. Activating.`), this.activate({ group: t, name: e, updateHash: !1 }));
  },
  /**
   * B4 APG 键盘旋转：计算目标 tab 名（不激活），调用方负责 activate + focus
   * （roving tabindex）。←→ 循环旋转，Home/End 跳首尾。
   * @param {string} group
   * @param {string} current
   * @param {string} key - ArrowRight/ArrowLeft/Home/End
   * @returns {string|null}
   */
  rotate(e, t, n) {
    const o = this.order[e];
    if (!o || o.length === 0) return null;
    const i = o.indexOf(t);
    return n === "Home" ? o[0] : n === "End" ? o[o.length - 1] : i === -1 ? o[0] : n === "ArrowRight" ? o[(i + 1) % o.length] : n === "ArrowLeft" ? o[(i - 1 + o.length) % o.length] : null;
  }
}), Ge = (r) => {
  r.directive("uno-tab", (e, { expression: t }, { evaluateLater: n, effect: o, cleanup: i }) => {
    if (!t) {
      h.warn("[Uno] x-uno-tab directive requires a configuration object.");
      return;
    }
    const s = n(t);
    let a = null;
    const l = (g) => typeof g == "object" && g !== null ? { ...g } : typeof g == "string" ? { name: g } : null;
    s((g) => {
      const d = l(g);
      if (!d) {
        h.warn("[Uno] x-uno-tab directive value must be a string or object", g);
        return;
      }
      d.group || (d.group = "default"), a = d, d.name && (r.store("uno").tab.register(d.name, d.group), e.setAttribute("role", "tab"), e.id || (e.id = "tab-" + d.name), e.setAttribute("aria-controls", d.panel || "panel-" + d.name));
    }), o(() => {
      if (!a || !a.name) return;
      const g = r.store("uno").tab.isActive(a.name, a.group);
      e.setAttribute("aria-selected", g ? "true" : "false"), e.setAttribute("tabindex", g ? "0" : "-1");
    });
    const c = (g) => {
      g.preventDefault(), s((d) => {
        const m = l(d);
        m && (m.group || (m.group = "default"), a = m, m.name && r.store("uno").tab.activate(m));
      });
    }, u = (g) => {
      if (!["ArrowRight", "ArrowLeft", "Home", "End"].includes(g.key) || !a || !a.name) return;
      g.preventDefault();
      const d = r.store("uno").tab, m = e.closest("[dir]") || document.documentElement, y = typeof getComputedStyle == "function" && getComputedStyle(m).direction === "rtl" ? g.key === "ArrowRight" ? "ArrowLeft" : g.key === "ArrowLeft" ? "ArrowRight" : g.key : g.key, p = d.rotate(a.group, a.name, y);
      if (!p) return;
      const x = document.getElementById("tab-" + p);
      x && (x.focus(), x.click());
    };
    e.addEventListener("click", c), e.addEventListener("keydown", u), i(() => {
      e.removeEventListener("click", c), e.removeEventListener("keydown", u);
    });
  });
};
const Je = (r, e = {}) => {
  const t = {
    defaultTab: "",
    hash: !0,
    activeClass: "is-active",
    inactiveClass: ""
  };
  b.registerPluginDefaults("tab", { ...t, ...e.plugins?.tab, ...e.tab });
  const n = Ze();
  F(r, "tab", n), Ge(r);
};
b.dom = j;
typeof window < "u" && (window.Uno = b);
b.use(ke).use(Be).use(De).use(He).use(Ke).use(Je);
if (typeof window < "u") {
  const r = () => {
    window.Alpine && !b.initialized && b.init(window.Alpine);
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
  b as Uno,
  He as box,
  b as default,
  De as notify,
  Je as tab,
  Be as tooltip,
  Ke as valid
};
