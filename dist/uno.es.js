import Y from "alpinejs";
let B = !1, N = null;
const h = {
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
    N = i;
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
    if (N)
      try {
        N(...i);
      } catch (e) {
        B && typeof console < "u" && typeof console.error == "function" && console.error("[Uno Error] Reporter failed:", e);
      }
    B && typeof console < "u" && typeof console.error == "function" && console.error("[Uno Error]", ...i);
  }
}, ie = /* @__PURE__ */ new WeakMap();
let ce = 0;
function j(i, e) {
  if (!i) return !1;
  const t = String(e), n = `${ce}:${t}`;
  let r = ie.get(i);
  return r || (r = /* @__PURE__ */ new Set(), ie.set(i, r)), r.has(n) ? !1 : (r.add(n), !0);
}
function R(i) {
  return i ? i.replace(/<script\b[^>]*>([\s\S]*?)<\/script>/gim, "").replace(/<img\b[^>]*>/gim, "").replace(/\s+on[a-z]+="[^"]*"/gim, "").replace(/\s+on[a-z]+='[^']*'/gim, "").replace(/\s+href="javascript:[^"]*"/gim, "").replace(/\s+href='javascript:[^']*'/gim, "").replace(/javascript:/gim, "") : "";
}
function de(i, e, { trusted: t = !1, alpine: n } = {}) {
  i && (i.innerHTML = t ? String(e ?? "") : R(String(e ?? "")), n?.initTree?.(i));
}
const V = {
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
  ensureContainer(i, e = "", t = {}, n = document.body) {
    const r = document.getElementById(i);
    if (r) return r;
    try {
      const o = document.createElement("div");
      return o.id = i, window.Uno?.options?.csp?.nonce && !t.nonce && (t.nonce = window.Uno.options.csp.nonce), Object.entries(t).forEach(([s, a]) => {
        o.setAttribute(s, a);
      }), e && (o.innerHTML = e), n ? n.appendChild(o) : document.documentElement.appendChild(o), o;
    } catch (o) {
      return h.error(`DomManager: Failed to inject container #${i}`, o), null;
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
  injectStyles(i, e) {
    if (!document.getElementById(i))
      try {
        const t = document.createElement("style");
        t.id = i, window.Uno?.options?.csp?.nonce && t.setAttribute("nonce", window.Uno.options.csp.nonce), t.textContent = e, document.head.appendChild(t);
      } catch (t) {
        h.error(`DomManager: Failed to inject styles #${i}`, t);
      }
  },
  // Exports
  setHTML: de,
  // Secure HTML setting
  sanitize: R
  // HTML sanitization
}, ue = (i, e) => {
  try {
    return i.getResponseHeader(e);
  } catch {
    return null;
  }
}, K = (i) => {
  const t = i?.detail?.xhr;
  return t ? {
    status: t.status,
    text: t.responseText ?? "",
    header: (n) => ue(t, n)
  } : null;
}, fe = (i) => i?.detail?.xhr;
function _(i) {
  return Object.prototype.toString.call(i) === "[object Object]";
}
function M(i, e) {
  if (!_(i)) return _(e) ? { ...e } : e;
  if (!_(e)) return e === void 0 ? { ...i } : e;
  const t = { ...i };
  for (const [n, r] of Object.entries(e)) {
    if (_(r) && _(t[n])) {
      t[n] = M(t[n], r);
      continue;
    }
    if (_(r)) {
      t[n] = M({}, r);
      continue;
    }
    if (typeof r == "string" && r.trim().indexOf("+") === 0 && typeof t[n] == "string") {
      t[n] = `${t[n]} ${r.trim().slice(1)}`.trim();
      continue;
    }
    t[n] = r;
  }
  return t;
}
const I = {
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
}, S = (i) => {
  if (!i) return I.BASE;
  const e = i.toUpperCase();
  return I[e] !== void 0 ? I[e] : I.BASE;
};
class me {
  /**
   * Create an UnoInstance.
   *
   * @param {Object} [options={}] - Global configuration options
   */
  constructor(e = {}) {
    this.version = "0.10.1", this.options = e, this.pluginDefaults = {}, this.plugins = [], this.initialized = !1, this._preInitConfig = {}, this._htmxRequestHandlers = [], this._htmxActiveXhrs = typeof WeakSet == "function" ? /* @__PURE__ */ new WeakSet() : null, this.htmx = {
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
      ajax: (r, o, s = {}) => typeof window > "u" || typeof window.htmx?.ajax != "function" ? Promise.reject(new Error("htmx.ajax is not available")) : window.htmx.ajax(r, o, s),
      /**
       * Fetch an HTML fragment without inserting it into the DOM.
       * Useful for getting content to display in a modal or tooltip.
       *
       * @param {string} method - HTTP method
       * @param {string} url - The URL to request
       * @param {Object} [ajaxOptions={}] - HTMX ajax options
       * @returns {Promise<string>} The HTML fragment string
       */
      fragment: async (r, o, s = {}) => {
        if (typeof document > "u") throw new Error("document is not available");
        const a = document.createElement("div");
        let c = "";
        const l = (d) => {
          d.preventDefault(), c = d.detail.serverResponse, d.detail.shouldSwap = !1;
        };
        a.addEventListener("htmx:beforeSwap", l);
        try {
          const d = s.swap || "innerHTML";
          return await this.htmx.ajax(r, o, { ...s, target: a, swap: d }), c;
        } catch (d) {
          throw h.error("HTMX fragment fetch failed:", { method: r, url: o, error: d }), d;
        } finally {
          a.removeEventListener("htmx:beforeSwap", l), a.remove();
        }
      },
      /**
       * Register handlers for HTMX request lifecycle events.
       *
       * @param {Object} handlers - The handlers object
       * @param {Function} [handlers.start] - Called when a request starts
       * @param {Function} [handlers.end] - Called when a request ends (success, error, or timeout)
       */
      onRequest: (r = {}) => {
        if (!r || typeof r != "object" || (this._htmxRequestHandlers.push(r), typeof window > "u" || typeof document > "u") || !window.htmx || !j(document, "uno:htmx:request-hooks")) return;
        const o = fe, s = (l, d) => {
          for (const m of this._htmxRequestHandlers) {
            const u = m?.[l];
            if (typeof u == "function")
              try {
                u(d);
              } catch (g) {
                h.error("Uno.htmx.onRequest handler failed", g);
              }
          }
        }, a = (l) => {
          const d = o(l);
          if (d && this._htmxActiveXhrs) {
            if (this._htmxActiveXhrs.has(d)) return;
            this._htmxActiveXhrs.add(d);
          }
          s("start", l);
        }, c = (l) => {
          const d = o(l);
          if (d && this._htmxActiveXhrs) {
            if (!this._htmxActiveXhrs.has(d)) return;
            this._htmxActiveXhrs.delete(d);
          }
          s("end", l);
        };
        document.addEventListener("htmx:beforeRequest", a), document.addEventListener("htmx:afterRequest", c), document.addEventListener("htmx:responseError", c), document.addEventListener("htmx:sendError", c), document.addEventListener("htmx:timeout", c), document.addEventListener("htmx:xhr:progress", (l) => s("progress", l));
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
      bindNotifyToHtmx: (r = {}) => {
        const {
          showSuccess: o = !0,
          showError: s = !0,
          showProgress: a = !0,
          successMessage: c = "Operation completed",
          errorMessage: l = "Operation failed"
        } = r, d = this;
        if (typeof document > "u") {
          h.warn("bindNotifyToHtmx: document is not available");
          return;
        }
        if (!d.htmx.available()) {
          h.warn("bindNotifyToHtmx: HTMX is not available");
          return;
        }
        const m = () => typeof window < "u" && window.Alpine ? window.Alpine.store("uno")?.notify : null;
        let u = 0;
        const g = (p) => {
          if (u++, a) {
            const x = m();
            x && typeof x.startProgress == "function" && x.startProgress();
          }
        }, w = (p) => {
          if (u = Math.max(0, u - 1), a && u === 0) {
            const x = m();
            x && typeof x.stopProgress == "function" && x.stopProgress();
          }
          if (o) {
            const x = K(p);
            if (!(x?.status >= 400)) {
              const v = m();
              if (v && typeof v.success == "function") {
                let C = c;
                const E = x?.header("X-Notify-Message");
                E && (C = E), x?.header("X-Notify-Suppress") || v.success(C);
              }
            }
          }
        }, y = (p) => {
          if (u = Math.max(0, u - 1), a && u === 0) {
            const x = m();
            x && typeof x.stopProgress == "function" && x.stopProgress();
          }
          if (s) {
            const x = m();
            if (x && typeof x.error == "function") {
              const f = K(p);
              let v = l;
              const C = f?.header("X-Notify-Error");
              C && (v = C), x.error(v);
            }
          }
        };
        j(document, "uno:htmx:bind-notify") && (document.addEventListener("htmx:beforeRequest", g), document.addEventListener("htmx:afterRequest", w), document.addEventListener("htmx:responseError", y), document.addEventListener("htmx:sendError", y), document.addEventListener("htmx:timeout", y), h.debug("bindNotifyToHtmx: HTMX events bound to notify"));
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
        layers: { box: I.BOX, tooltip: I.TOOLTIP, notify: I.NOTIFY },
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
    let r = this.pluginDefaults[e] ? { ...this.pluginDefaults[e] } : {};
    return n.plugins && n.plugins[e] && (r = M(r, n.plugins[e])), n[e] && (r = M(r, n[e])), t && (r = M(r, t)), r;
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
    const r = t.split(".").filter(Boolean);
    let o = n;
    for (const s of r) {
      if (o == null) return;
      o = o[s];
    }
    return o;
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
    const r = this._getUiValue(e, t);
    return r === void 0 ? this._resolveUiTokens(n) : typeof r == "string" && r.trim().indexOf("+") === 0 ? this._resolveUiTokens(`${n} ${r.trim().slice(1)}`.trim()) : this._resolveUiTokens(r);
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
      return h.warn("Uno has already been initialized."), Y || (typeof window < "u" ? window.Alpine : null);
    let n = Y;
    return !n && typeof window < "u" && window.Alpine && (n = window.Alpine), e && (e.version || typeof e.start == "function" || typeof e.store == "function") && (n = e), this.options = M(this.options, this._preInitConfig), h.setDebug(this.options.debug !== !1), this.options.onError && h.setReporter(this.options.onError), n && typeof n.store == "function" ? (this._doInit(n), n) : typeof document > "u" ? (h.error("Alpine not found and document is unavailable. Cannot initialize Uno."), null) : (h.debug("Alpine not ready, waiting for alpine:init event..."), new Promise((o) => {
      document.addEventListener(
        "alpine:init",
        () => {
          if (n = window.Alpine || Y, !n) {
            h.error("Alpine still not found after alpine:init event"), o(null);
            return;
          }
          h.debug("Alpine ready, continuing Uno initialization"), this._doInit(n), o(n);
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
    h.info(`Uno v${this.version} initializing...`), typeof window < "u" && typeof document < "u" && window.htmx && typeof e?.initTree == "function" && j(document, "uno:htmx-afterSwap:initTree") && document.addEventListener("htmx:afterSwap", (t) => {
      const n = t?.detail?.elt;
      n && e.initTree(n);
    }), this.plugins.forEach(({ plugin: t, options: n }) => {
      try {
        const r = M(this.options, n);
        t(e, r);
      } catch (r) {
        h.error("Failed to load plugin", r);
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
    if (!(typeof window > "u") && j(window, "uno:global-error-handlers")) {
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
const b = new me(), ge = {
  box: {
    ok: "OK",
    cancel: "Cancel",
    loading: "Loading",
    close: "Close",
    dialog: "Dialog",
    prompt: "Input"
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
    strongpwd: "Password must be 8-80 chars, contain letters, numbers, and at least 1 special char (@$!%*#?&)",
    summary_title: "Please fix the following problems"
  },
  charCount: "Remaining: {count}"
}, he = {
  box: {
    ok: "确认",
    cancel: "取消",
    loading: "加载中",
    close: "关闭",
    dialog: "对话框",
    prompt: "输入"
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
    strongpwd: "密码需 8-80 位，必须包含字母、数字和至少 1 个特殊字符 (@$!%*#?&)",
    summary_title: "请修正以下问题"
  },
  charCount: "剩余: {count}"
}, pe = {
  box: {
    ok: "確認",
    cancel: "取消",
    loading: "載入中",
    close: "關閉",
    dialog: "對話框",
    prompt: "輸入"
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
    strongpwd: "密碼需 8-80 個字元，必須包含字母、數字和至少 1 個特殊字元 (@$!%*#?&)",
    summary_title: "請修正以下問題"
  },
  charCount: "剩餘: {count}"
}, xe = {
  box: {
    ok: "Aceptar",
    cancel: "Cancelar",
    loading: "Cargando",
    close: "Cerrar",
    dialog: "Diálogo",
    prompt: "Entrada"
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
    strongpwd: "La contraseña debe tener entre 8 y 80 caracteres, incluir letras, números y al menos 1 carácter especial (@$!%*#?&)",
    summary_title: "Corrige los siguientes problemas"
  },
  charCount: "Restantes: {count}"
}, ye = {
  box: {
    ok: "確認",
    cancel: "キャンセル",
    loading: "読み込み中",
    close: "閉じる",
    dialog: "ダイアログ",
    prompt: "入力"
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
    strongpwd: "パスワードは 8〜80 文字で、英字・数字を含み、記号（@$!%*#?&）を 1 文字以上含めてください",
    summary_title: "次の問題を修正してください"
  },
  charCount: "残り：{count}"
}, be = {
  box: {
    ok: "OK",
    cancel: "Abbrechen",
    loading: "Wird geladen",
    close: "Schließen",
    dialog: "Dialog",
    prompt: "Eingabe"
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
    strongpwd: "Das Passwort muss 8-80 Zeichen lang sein, Buchstaben, Zahlen und mindestens ein Sonderzeichen (@$!%*#?&) enthalten",
    summary_title: "Bitte behebe die folgenden Probleme"
  },
  charCount: "Verbleibend: {count}"
}, ve = {
  box: {
    ok: "Valider",
    cancel: "Annuler",
    loading: "Chargement",
    close: "Fermer",
    dialog: "Boîte de dialogue",
    prompt: "Saisie"
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
    strongpwd: "Le mot de passe doit comporter entre 8 et 80 caractères, contenir des lettres, des chiffres et au moins un caractère spécial (@$!%*#?&)",
    summary_title: "Veuillez corriger les problèmes suivants"
  },
  charCount: "Caractères restants : {count}"
}, we = {
  box: {
    ok: "확인",
    cancel: "취소",
    loading: "불러오는 중",
    close: "닫기",
    dialog: "대화 상자",
    prompt: "입력"
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
    strongpwd: "비밀번호는 8~80자이며 영문, 숫자, 특수문자 1자 이상(@$!%*#?&)을 포함해야 해요",
    summary_title: "다음 항목을 수정해 주세요"
  },
  charCount: "남은 글자: {count}"
}, Ee = {
  box: {
    ok: "Confirmar",
    cancel: "Cancelar",
    loading: "Carregando",
    close: "Fechar",
    dialog: "Caixa de diálogo",
    prompt: "Entrada"
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
    strongpwd: "A senha deve ter de 8 a 80 caracteres e conter letras, números e pelo menos 1 caractere especial (@$!%*#?&)",
    summary_title: "Corrija os seguintes problemas"
  },
  charCount: "Restantes: {count}"
}, ke = {
  box: {
    ok: "ОК",
    cancel: "Отмена",
    loading: "Загрузка",
    close: "Закрыть",
    dialog: "Диалоговое окно",
    prompt: "Ввод"
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
    strongpwd: "Пароль должен содержать от 8 до 80 символов, включать буквы, цифры и хотя бы один специальный символ (@$!%*#?&)",
    summary_title: "Пожалуйста, исправьте следующие ошибки"
  },
  charCount: "Осталось: {count}"
}, Ce = {
  box: {
    ok: "موافق",
    cancel: "إلغاء",
    loading: "جارٍ التحميل",
    close: "إغلاق",
    dialog: "مربع حوار",
    prompt: "إدخال"
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
    strongpwd: "يجب أن تتكون كلمة المرور من 8 إلى 80 حرفًا، وأن تحتوي على أحرف وأرقام، ورمز خاص واحد على الأقل (@$!%*#?&)",
    summary_title: "يُرجى تصحيح المشكلات التالية"
  },
  charCount: "المتبقي: {count}"
};
const Z = {
  en: ge,
  "zh-CN": he,
  "zh-TW": pe,
  es: xe,
  ja: ye,
  de: be,
  fr: ve,
  ko: we,
  pt: Ee,
  ru: ke,
  ar: Ce
}, G = Object.keys(Z);
function z(i) {
  if (!i) return "";
  const e = String(i).trim().replace(/_/g, "-"), t = e.toLowerCase();
  if (t === "zh" || t.startsWith("zh-"))
    return t.includes("hant") ? "zh-TW" : t.includes("hans") ? "zh-CN" : /(?:^|-)(?:tw|hk|mo)(?:-|$)/.test(t) ? "zh-TW" : "zh-CN";
  const n = t.split("-")[0];
  return G.find((o) => o.toLowerCase() === t || o.toLowerCase() === n) || e;
}
function Me() {
  if (typeof localStorage < "u") {
    const i = z(localStorage.getItem("uno.lang"));
    if (i && G.includes(i)) return i;
  }
  if (typeof navigator < "u") {
    const i = [navigator.language, ...navigator.languages || []];
    for (const e of i) {
      const t = z(e);
      if (G.includes(t)) return t;
    }
    return "en";
  }
  return "en";
}
function $e(i, e) {
  const t = e.split(".");
  let n = i;
  for (const r of t) {
    if (!n || typeof n != "object" || !(r in n)) return;
    n = n[r];
  }
  return n;
}
function J(i, e) {
  for (const t in e)
    e[t] && typeof e[t] == "object" ? (i[t] || Object.assign(i, { [t]: {} }), J(i[t], e[t])) : Object.assign(i, { [t]: e[t] });
}
function Se(i) {
  const e = i.store("i18n"), t = {
    locale: Me(),
    messages: Z,
    loading: !1,
    // <html lang> 只由服务端模板输出；Uno 只维护组件词条语言，不回写文档属性。
    init() {
    },
    t(n, r = {}) {
      const o = z(this.locale), s = this.messages[o] || this.messages.en, a = $e(s, n);
      return (typeof a == "string" ? a : n).replace(/\{(\w+)\}/g, (l, d) => d in r ? String(r[d]) : `{${d}}`);
    },
    setLocale(n) {
      const r = z(n);
      if (!this.messages[r]) {
        h.warn(`[Uno] Language ${n} not found. Please load it first using loadLocale().`);
        return;
      }
      this.locale = r, typeof localStorage < "u" && localStorage.setItem("uno.lang", r);
    },
    // Alias for backward compatibility
    setLang(n) {
      this.setLocale(n);
    },
    extend(n) {
      J(this.messages, n);
    },
    async loadLocale(n, r) {
      const o = z(n);
      if (!this.messages[o]) {
        this.loading = !0;
        try {
          const s = await fetch(r);
          if (!s.ok) throw new Error(`Failed to load ${n}`);
          const a = await s.json();
          this.messages[o] = a, this.setLocale(o);
        } catch (s) {
          h.error(`[Uno] Failed to load locale ${n}:`, s);
        } finally {
          this.loading = !1;
        }
      }
    }
  };
  e ? (h.info("[Uno] Merging into existing i18n store"), e.messages || (e.messages = {}), J(e.messages, Z), e.t || (e.t = t.t), e.setLocale || (e.setLocale = t.setLocale), e.setLang || (e.setLang = t.setLang), e.loadLocale || (e.loadLocale = t.loadLocale)) : i.store("i18n", t), i.directive("t", (n, { expression: r, modifiers: o }, { effect: s }) => {
    const a = () => i.store("i18n");
    s(() => {
      const c = a();
      c.locale;
      const l = r, d = c.t(l);
      o.length > 0 ? o.forEach((m) => {
        n.setAttribute(m, d);
      }) : n.textContent = d;
    });
  });
}
function H(i, e, t) {
  if (!i || typeof i.store != "function") {
    h.error("[Uno] Alpine.store is not available");
    return;
  }
  const n = i.store("uno");
  n ? i.store("uno", {
    ...n,
    [e]: t
  }) : i.store("uno", { [e]: t }), i.store(e) || i.store(e, i.store("uno")?.[e]);
}
const Te = (i) => {
  const e = i || "center", t = {
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
}, Ie = (i) => {
  const e = i || "top-right", t = {
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
}, Le = (i) => {
  const e = i || "top-right";
  let t = "flex-col", n = "items-end", r = "space-y-2";
  return e.includes("bottom") && (t = "flex-col-reverse", r = "space-y-reverse space-y-2"), e.includes("left") ? n = "items-start" : e.includes("center") ? e === "center" ? n = "items-center justify-center" : n = "items-center" : n = "items-end", { flexDir: t, alignItems: n, spaceY: r };
};
function Be(i, e, t, n = {}) {
  const r = typeof n?.offset == "number" && Number.isFinite(n.offset) && n.offset >= 0 ? n.offset : 8, o = i.getBoundingClientRect(), s = e.getBoundingClientRect(), c = (typeof getComputedStyle == "function" && getComputedStyle(i).direction === "rtl" ? "rtl" : "ltr") === "rtl" ? t === "left" ? "right" : t === "right" ? "left" : t : t, l = s.width, d = s.height, m = o.width, u = o.height, g = { top: o.top, left: o.left }, w = window.innerWidth, y = window.innerHeight, p = 0, x = 0, f = y, v = w, C = (T, q) => {
    if (q) {
      if (T + l > v)
        return Math.max(x + 10, v - l - 10);
      if (T < x)
        return x + 10;
    } else {
      if (T + d > f)
        return f - d - 10;
      if (T < p)
        return p + 10;
    }
    return T;
  };
  let E, k, L = ["top", "bottom", "left", "right"].includes(c) ? c : "bottom";
  switch (c) {
    case "left":
      k = g.left - l - r, k < x && (k = g.left + m + r, L = "right"), E = C(g.top + u / 2 - d / 2, !1);
      break;
    case "top":
      E = g.top - d - r, E < p && (E = g.top + u + r, L = "bottom"), k = g.left + (m - l) / 2, k = C(k, !0);
      break;
    case "right":
      k = g.left + m + r, k + l > v && (k = g.left - l - r, L = "left"), E = C(g.top + u / 2 - d / 2, !1);
      break;
    default:
      E = g.top + u + r, E + d > f && (E = g.top - d - r, L = "top"), k = g.left + (m - l) / 2, k = C(k, !0);
      break;
  }
  const O = 8, $ = O / 2, F = (T, q, le) => Math.min(Math.max(T, q), Math.max(q, le)), ne = g.left + m / 2, oe = g.top + u / 2;
  let D;
  switch (L) {
    case "top":
      D = {
        top: d - $,
        left: F(ne - k - $, 0, l - O)
      };
      break;
    case "left":
      D = {
        top: F(oe - E - $, 0, d - O),
        left: l - $
      };
      break;
    case "right":
      D = {
        top: F(oe - E - $, 0, d - O),
        left: -$
      };
      break;
    default:
      D = {
        top: -$,
        left: F(ne - k - $, 0, l - O)
      };
      break;
  }
  return { top: E, left: k, position: L, arrow: D };
}
const Q = [
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
function re(i, e = {}) {
  const { initialFocus: t = !0, enabled: n = !0 } = e;
  if (!n) return () => {
  };
  const r = (l) => {
    if (!(l instanceof HTMLElement)) return !1;
    const d = window.getComputedStyle(l);
    return d.display === "none" || d.visibility === "hidden" ? !1 : l.getClientRects().length > 0;
  }, o = () => Array.from(i.querySelectorAll(Q)).filter((l) => r(l)), s = () => o()[0], a = (l) => {
    if (l.key !== "Tab") return;
    const d = o();
    if (!d.length) return;
    const m = d[0], u = d[d.length - 1], g = document.activeElement, w = d.some(
      (p) => g.compareDocumentPosition(p) & Node.DOCUMENT_POSITION_FOLLOWING
    ), y = d.some(
      (p) => g.compareDocumentPosition(p) & Node.DOCUMENT_POSITION_PRECEDING
    );
    l.shiftKey ? (g === m || !y) && (l.preventDefault(), u.focus()) : (g === u || !w) && (l.preventDefault(), m.focus());
  }, c = () => {
    if (i.contains(document.activeElement)) return;
    const l = s();
    l && l.focus();
  };
  return t && (typeof requestAnimationFrame == "function" ? requestAnimationFrame(() => requestAnimationFrame(c)) : setTimeout(c, 0)), i.addEventListener("keydown", a), () => {
    i.removeEventListener("keydown", a);
  };
}
const _e = () => (b.options?.ui?.layers?.tooltip || S("TOOLTIP"), `
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
let A = 0, U = null, P = null;
const ee = "uno-tooltip-content", Ae = (i) => {
  const e = (i.getAttribute("aria-describedby") || "").split(/\s+/).filter(Boolean);
  e.includes(ee) || e.push(ee), i.setAttribute("aria-describedby", e.join(" "));
}, Ue = (i) => {
  const e = (i.getAttribute("aria-describedby") || "").split(/\s+/).filter((t) => t && t !== ee);
  e.length ? i.setAttribute("aria-describedby", e.join(" ")) : i.removeAttribute("aria-describedby");
}, Oe = (i) => {
  const e = () => {
    A || typeof requestAnimationFrame == "function" && (A = requestAnimationFrame(() => {
      A = 0;
      const o = i.store("uno")?.tooltip;
      o && o.updatePosition();
    }));
  }, t = () => {
    P && (P.disconnect(), P = null);
  }, n = () => {
    if (t(), typeof ResizeObserver != "function") return;
    const o = document.getElementById("uno-tooltip-content");
    o && (P = new ResizeObserver(() => e()), P.observe(o));
  }, r = () => {
    if (U || typeof window > "u" || typeof document > "u") return;
    const o = () => e(), s = () => e(), a = (l) => {
      if (l.key !== "Escape") return;
      const d = i.store("uno")?.tooltip;
      d?.visible && d.hide();
    }, c = (l) => {
      const d = i.store("uno")?.tooltip;
      if (!d?.visible || d.triggerMode !== "click" || d.clickOutside !== !0) return;
      const m = l.target, u = document.getElementById("uno-tooltip-content"), g = d.triggerEl && d.triggerEl.contains(m), w = u && u.contains(m);
      !g && !w && d.hide();
    };
    window.addEventListener("scroll", o, !0), window.addEventListener("resize", s), window.addEventListener("keydown", a), window.addEventListener("click", c), U = () => {
      window.removeEventListener("scroll", o, !0), window.removeEventListener("resize", s), window.removeEventListener("keydown", a), window.removeEventListener("click", c), U = null;
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
    // A3 keyboard bindings (click trigger popup): keydown handlers
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
    show(o = {}) {
      h.info("Tooltip show called", o), this.cancelHide();
      const s = b.getConfig("tooltip", o);
      if (h.info("Tooltip resolved config", s), !s.target && typeof s.content == "string" && s.content.startsWith("#") && (s.content.includes(" ") || document.querySelector(s.content) && (s.target = s.content)), this.triggerEl = s.element || this.triggerEl, s.target) {
        const a = document.querySelector(s.target);
        a ? this.content = R(a.innerHTML) : (h.warn(`Tooltip target element not found: ${s.target}`), this.content = R(s.content || ""));
      } else
        this.content = R(s.content || "");
      this.placement = s.placement, this.triggerMode = s.trigger, this.clickOutside = s.clickOutside, this.offset = s.offset, this.arrow = s.arrow !== !1, this.zIndex = b.options?.ui?.layers?.tooltip || S("TOOLTIP"), this.hxUrl = s.hxUrl, this.hxHeaders = s.hxHeaders, this.hxTrigger = s.hxTrigger || "load", this.visible = !0, this.ensureContainer(), r(), n(), this.triggerEl && (Ae(this.triggerEl), this.triggerMode === "click" && (this.triggerEl.setAttribute("aria-haspopup", "true"), this.triggerEl.setAttribute("aria-expanded", "true"), this.bindPopupKeyboard())), i.nextTick(() => {
        this.updatePosition(), this.focusPopupIfNeeded();
      });
    },
    hide() {
      this.cancelHide();
      const o = this.triggerEl, s = this.triggerMode === "click", a = document.getElementById("uno-tooltip-content"), c = document.activeElement, l = !!a && !!c && a.contains(c);
      this.visible = !1, this.triggerEl = null, a && this.onPopupKeydown && (a.removeEventListener("keydown", this.onPopupKeydown), this.onPopupKeydown = null), o && this.onTriggerKeydown && (o.removeEventListener("keydown", this.onTriggerKeydown), this.onTriggerKeydown = null), o && o.setAttribute && (Ue(o), s && o.setAttribute("aria-expanded", "false")), U && (U(), U = null), t(), A && (cancelAnimationFrame(A), A = 0), s && o && typeof o.focus == "function" && (l || c === document.body) && o.focus();
    },
    /**
     * A3 click 弹层键盘绑定（APG menu-button）：Tab 关闭弹层并离开（焦点先还原触发器，
     * 浏览器默认 Tab 动作随后落到触发器之后的元素，不 trap）；role=menu 内容方向键导航；
     * 锚点菜单项 Space 合成 click（Enter 由原生 click 覆盖，按钮 Space 亦原生）；
     * 锚点触发器 Space 合成关闭（仅有关闭态需合成，打开由 click 覆盖）。
     */
    bindPopupKeyboard() {
      const o = document.getElementById("uno-tooltip-content");
      o && (this.onPopupKeydown = (s) => {
        if (s.key === "Tab") {
          this.hide();
          return;
        }
        if (s.key === " " || s.key === "Spacebar") {
          const d = document.activeElement;
          d && d.tagName === "A" && o.contains(d) && (s.preventDefault(), d.click());
          return;
        }
        if (!["ArrowDown", "ArrowUp", "Home", "End"].includes(s.key) || !o.querySelector('[role="menu"]')) return;
        const a = Array.from(o.querySelectorAll(Q));
        if (a.length === 0) return;
        s.preventDefault();
        const c = a.indexOf(document.activeElement);
        let l = 0;
        s.key === "ArrowDown" ? l = (c + 1) % a.length : s.key === "ArrowUp" ? l = (c - 1 + a.length) % a.length : s.key === "End" && (l = a.length - 1), a[l] && a[l].focus();
      }, o.addEventListener("keydown", this.onPopupKeydown), this.triggerEl && this.triggerEl.tagName === "A" && (this.onTriggerKeydown = (s) => {
        s.key !== " " && s.key !== "Spacebar" || (s.preventDefault(), this.visible && this.triggerEl === s.currentTarget && this.hide());
      }, this.triggerEl.addEventListener("keydown", this.onTriggerKeydown)));
    },
    /**
     * A3 焦点移入弹层首个可聚焦元素，键盘用户落在刚打开的披露内容内（2.4.3）；
     * 弹层无可聚焦内容时为 no-op（焦点留在触发器）。
     */
    focusPopupIfNeeded() {
      if (this.triggerMode !== "click") return;
      const o = document.getElementById("uno-tooltip-content");
      if (!o || o.contains(document.activeElement)) return;
      const s = o.querySelector(Q);
      s && s.focus();
    },
    cancelHide() {
      this.hideTimeout && (clearTimeout(this.hideTimeout), this.hideTimeout = null);
    },
    scheduleHide(o = 300) {
      this.cancelHide(), this.hideTimeout = setTimeout(() => {
        this.hide();
      }, o);
    },
    /**
     * Resolve the offset used by the next positioning pass.
     *
     * A number is used as-is. A function is invoked on every positioning pass
     * (open, scroll, resize, size change) with the live geometry, so callers can
     * derive the gap instead of freezing it at trigger time — breakpoint switches
     * and layout changes then stay correct while the popup is open. A throwing or
     * non-numeric result falls back to the default gap rather than breaking the
     * popup.
     *
     * @param {HTMLElement} tooltipEl - The popup element being positioned
     * @returns {number|undefined} Pixels between trigger and popup
     */
    resolveOffset(o) {
      if (typeof this.offset != "function") return this.offset;
      try {
        const s = this.offset({
          triggerEl: this.triggerEl,
          tooltipEl: o,
          placement: this.placement
        });
        return typeof s == "number" && Number.isFinite(s) ? s : void 0;
      } catch (s) {
        h.warn("[Uno] tooltip offset() threw, falling back to the default gap", s);
        return;
      }
    },
    updatePosition() {
      if (!this.visible || !this.triggerEl) return;
      const o = document.getElementById("uno-tooltip-content");
      if (!o) return;
      const s = Be(this.triggerEl, o, this.placement, {
        offset: this.resolveOffset(o)
      });
      this.top = s.top, this.left = s.left, this.placement = s.position || this.placement, s.arrow && (this.arrowPosition = s.arrow);
    },
    ensureContainer() {
      const o = "uno-tooltip-manager";
      if (document.getElementById(o)) return;
      const s = b.options?.ui?.layers?.tooltip || S("TOOLTIP"), a = _e();
      V.ensureContainer(o, a, {
        "x-data": "{}",
        style: `position: absolute; top: 0; left: 0; width: 0; height: 0; z-index: ${s};`
      });
      const c = document.getElementById(o);
      c && i && !c._x_dataStack && i.initTree(c);
    }
  };
}, De = (i) => {
  const e = (t, { expression: n, modifiers: r }, { evaluateLater: o, effect: s, cleanup: a }) => {
    const c = n ? o(n) : null, l = (f) => (r.length > 0 && h.warn(
      "[Uno] x-uno-tooltip modifiers are deprecated in v2.0. Use configuration object instead."
    ), { ...typeof f == "object" && f !== null ? f : f ? { content: f } : {}, element: t }), d = (f) => {
      const v = i.store("uno")?.tooltip;
      v && v.show(l(f));
    }, m = (f) => {
      const v = i.store("uno")?.tooltip;
      v && (v.visible && v.triggerEl === t ? v.hide() : d(f));
    };
    let u = null;
    c && s(() => {
      c((f) => {
        u = f;
      });
    });
    const g = () => {
      if (n && u === null) return;
      const f = l(u);
      b.getConfig("tooltip", f).trigger === "hover" && d(u);
    }, w = () => {
      const f = i.store("uno")?.tooltip;
      f && f.triggerMode === "hover" && f.triggerEl === t && f.scheduleHide();
    }, y = (f) => {
      const v = l(u);
      b.getConfig("tooltip", v).trigger === "click" && (f.preventDefault(), m(u));
    }, p = () => {
      if (n && u === null) return;
      const f = l(u), v = b.getConfig("tooltip", f);
      (v.trigger === "focus" || v.trigger === "hover") && d(u);
    }, x = () => {
      const f = i.store("uno")?.tooltip;
      !f || f.triggerEl !== t || (f.triggerMode === "focus" ? f.hide() : f.triggerMode === "hover" && !t.matches(":hover") && f.scheduleHide());
    };
    t.addEventListener("mouseenter", g), t.addEventListener("mouseleave", w), t.addEventListener("click", y), t.addEventListener("focus", p), t.addEventListener("blur", x), a(() => {
      t.removeEventListener("mouseenter", g), t.removeEventListener("mouseleave", w), t.removeEventListener("click", y), t.removeEventListener("focus", p), t.removeEventListener("blur", x);
      const f = i.store("uno")?.tooltip;
      f && f.triggerEl === t && f.hide();
    });
  };
  i.directive("uno-tooltip", e);
};
const Pe = (i, e = {}) => {
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
  const n = Oe(i);
  H(i, "tooltip", n), De(i);
}, Re = () => `
  <template x-for="(items, position) in $store.uno.notify.groups" :key="position">
    <div :class="$store.uno.notify.getContainerClass(position)"
         :style="$store.uno.notify.getContainerStyle(position)">
      
      <template x-for="item in items" :key="item.id">
        <div :class="[
            Uno.getStyle('notify', 'item'),
            Uno.getStyle('notify', 'types.' + item.type)
          ]"
          x-show="item.isVisible"
          x-transition:enter="transition ease-out duration-300"
          x-transition:enter-start="opacity-0 transform scale-90"
          x-transition:enter-end="opacity-100 transform scale-100"
          x-transition:leave="transition ease-in duration-200"
          x-transition:leave-start="opacity-100 transform scale-100"
          x-transition:leave-end="opacity-0 transform scale-90"
          :data-id="item.id"
          data-testid="notify-item"
          @mouseenter="$store.uno.notify.pause(item.id, 'hover')"
          @mouseleave="$store.uno.notify.resume(item.id, 'hover')"
          @focusin="$store.uno.notify.pause(item.id, 'focus')"
          @focusout="!$el.contains($event.relatedTarget) && $store.uno.notify.resume(item.id, 'focus')"
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
                    <p :class="Uno.getStyle('notify', 'text')" x-text="item.message"
                       :role="item.type === 'error' || item.type === 'warning' ? 'alert' : 'status'"
                       :aria-live="item.type === 'error' || item.type === 'warning' ? 'assertive' : 'polite'"></p>
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
                  <div x-html="item.content || ''"
                       :role="item.type === 'error' || item.type === 'warning' ? 'alert' : 'status'"
                       :aria-live="item.type === 'error' || item.type === 'warning' ? 'assertive' : 'polite'"></div>
               </div>
            </template>
          </div>

        </div>
      </template>
    </div>
  </template>
`, ze = () => `
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
  `, He = (i) => ({
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
    const n = e.innerHTML || e.textContent || "", r = this._findNotificationById(t.id);
    r && (r.content = n), t.onContentLoad && t.onContentLoad(e), window.Alpine && window.Alpine.initTree && window.Alpine.initTree(e);
  },
  // Helper method to find notification by ID
  _findNotificationById(e) {
    for (const t of Object.keys(this.groups)) {
      const r = this.groups[t].find((o) => o.id === e);
      if (r) return r;
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
    const t = b.getStyle("notify", "container"), n = Ie(e), r = Le(e);
    return [t, n, r.flexDir, r.alignItems].join(" ");
  },
  ensureContainer() {
    const e = "uno-notify-container", t = Re(), n = V.ensureContainer(e, t, {
      "x-data": "{}"
      // Independent Alpine component
    });
    n && i.initTree && i.initTree(n);
    const r = "uno-global-progress", o = ze(), s = V.ensureContainer(r, o, {
      "x-data": "{}"
    });
    s && i.initTree && i.initTree(s);
  },
  show(e, t = {}) {
    this.ensureContainer();
    const n = b.getConfig("notify", t), r = this.nextId++;
    let o = e;
    i.store("i18n") && typeof i.store("i18n").t == "function" && (o = i.store("i18n").t(e));
    const s = n.duration !== void 0 ? n.duration : 3e3;
    h.info("[Notify] Show:", { id: r, message: o, duration: s, config: n });
    const a = {
      id: r,
      message: o,
      isVisible: !0,
      // For x-show
      closable: !0,
      duration: s,
      content: n.hxUrl ? "" : void 0,
      // Initialize content for HTMX-loaded content
      ...n,
      // 可暂停的自动关闭计时——保留剩余时长，hover/focus 期间不走表
      _remaining: s,
      _timer: null,
      _startedAt: 0,
      _paused: { hover: !1, focus: !1 }
    }, c = n.position || "top-right";
    return this.groups[c] ? this.groups[c].push(a) : (this.groups["top-right"] || (this.groups["top-right"] = []), this.groups["top-right"].push(a)), this._startDismiss(this._findNotificationById(r)), r;
  },
  /**
   * 启动自动关闭计时。仅对 duration>0、可见、未暂停且未在计时的通知生效；
   * 以 _remaining 为延时，暂停后恢复时续走剩余时长而非重置。
   */
  _startDismiss(e) {
    !e || !(e.duration > 0) || !e.isVisible || e._timer || e._paused.hover || e._paused.focus || (e._startedAt = Date.now(), e._timer = setTimeout(() => {
      e._timer = null, this.remove(e.id);
    }, e._remaining));
  },
  /**
   * 暂停自动关闭（鼠标悬停或键盘聚焦）。结算已走时长并保留剩余，
   * 使读者/键盘用户有充足时间阅读与操作，不被超时打断。
   * @param {number} id
   * @param {'hover'|'focus'} reason
   */
  pause(e, t) {
    const n = this._findNotificationById(e);
    n && (t && (n._paused[t] = !0), n._timer && (clearTimeout(n._timer), n._timer = null, n._remaining = Math.max(0, n._remaining - (Date.now() - n._startedAt))));
  },
  /**
   * 恢复自动关闭。仅当 hover 与 focus 两种暂停都解除后续走剩余时长。
   * @param {number} id
   * @param {'hover'|'focus'} reason
   */
  resume(e, t) {
    const n = this._findNotificationById(e);
    n && (t && (n._paused[t] = !1), !n._paused.hover && !n._paused.focus && this._startDismiss(n));
  },
  remove(e) {
    const t = this._findNotificationById(e);
    if (t && t._timer && (clearTimeout(t._timer), t._timer = null), typeof document < "u") {
      const n = document.querySelector(`[data-id="${e}"]`);
      if (n && n.contains(document.activeElement)) {
        const r = Array.from(document.querySelectorAll('[data-testid="notify-item"]')).filter((o) => o !== n).map((o) => o.querySelector("button")).find(Boolean);
        r && r.focus();
      }
    }
    Object.keys(this.groups).forEach((n) => {
      const r = this.groups[n].findIndex((o) => o.id === e);
      r !== -1 && (this.groups[n][r].isVisible = !1, setTimeout(() => {
        const o = this.groups[n].findIndex((s) => s.id === e);
        o !== -1 && this.groups[n].splice(o, 1);
      }, 300));
    });
  },
  // 3. Proxy Methods
  success(e, t = {}) {
    return this.show(e, { type: "success", ...t });
  },
  error(e, t = {}) {
    return this.show(e, { type: "error", duration: 0, ...t });
  },
  warning(e, t = {}) {
    return this.show(e, { type: "warning", ...t });
  },
  info(e, t = {}) {
    return this.show(e, { type: "info", ...t });
  }
}), Fe = (i) => {
  i.directive("uno-notify", (e, { expression: t }, { evaluateLater: n }) => {
    if (!t) {
      h.warn("[Uno] x-uno-notify directive requires a configuration object or message string.");
      return;
    }
    const r = n(t), o = () => {
      const s = i.store("uno")?.notify;
      s && r((a) => {
        let c = "", l = {};
        if (typeof a == "string")
          c = a;
        else if (typeof a == "object" && a !== null)
          c = a.message || "", l = a;
        else {
          h.warn("[Uno] x-uno-notify directive value must be a string or object", a);
          return;
        }
        s.show(c, l);
      });
    };
    e.addEventListener("click", o);
  });
};
const qe = (i, e = {}) => {
  const t = {
    duration: 3e3,
    position: "top-right"
  };
  b.registerPluginDefaults("notify", {
    ...t,
    ...e.plugins?.notify,
    ...e.notify
  });
  const n = He(i);
  H(i, "notify", n), Fe(i);
}, je = () => `
  <!-- Backdrop / Overlay（仅模态弹窗存在时显示，非模态不遮挡背景） -->
  <div id="uno-box-overlay"
       :class="Uno.getStyle('box', 'overlay')"
       x-show="$store.uno.box.hasModalOpen()"
       x-transition:enter="ease-out duration-300"
       x-transition:enter-start="opacity-0"
       x-transition:enter-end="opacity-100"
       x-transition:leave="ease-in duration-200"
       x-transition:leave-start="opacity-100"
       x-transition:leave-end="opacity-0"
       @click="$store.uno.box.items.length && $store.uno.box.items[$store.uno.box.items.length - 1].closeOnOverlayClick && $store.uno.box.items[$store.uno.box.items.length - 1].close()"></div>

  <div :class="[Uno.getStyle('box', 'dialogWrapper'), 'grid grid-cols-1', $store.uno.box.hasModalOpen() ? '' : 'pointer-events-none']"
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
           :aria-labelledby="item.labelledBy || (item.title ? 'uno-box-title-' + item.id : null)"
           :aria-label="item.labelledBy || item.title ? null : $store.i18n.t('box.dialog')"
           :aria-describedby="'uno-box-content-' + item.id"
           @keydown.escape.window="$store.uno.box.isTopmost(item) && item.closeOnEsc && !$store.uno.tooltip?.visible && item.close()">

        <div :id="'uno-box-container' + item.id"
             :class="[Uno.getStyle('box', 'container'), item.widthClass, item.heightClass || (item.fill && !item.height ? Uno.getStyle('box', 'containerFill') : ''), item.containerClass, 'focus:outline-none']"
             :style="(item.width ? 'width:' + item.width + 'px; max-width:' + item.width + 'px;' : (item.widthClass ? 'max-width:none;' : '')) + (item.height ? ' height:' + item.height + 'px;' : '')"
             tabindex="-1"
             x-uno-trap="{ initialFocus: false, enabled: item.modal !== false }">

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
`, Ve = (i) => ({
  items: [],
  /**
   * Dialog stack management
   * @type {Array}
   */
  stack: [],
  // 记录本组件 inert 过的背景节点，关闭时只恢复这些节点，保留应用原有 inert 状态。
  _backgroundInert: /* @__PURE__ */ new Set(),
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
      const { onOk: r, onCancel: o, onClose: s, ...a } = t;
      this.open("confirm", {
        ...a,
        content: e,
        onOk: () => {
          const c = typeof r == "function" ? r() : void 0;
          if (c === !1) return !1;
          if (c instanceof Promise)
            return c.then((l) => {
              if (l === !1) return !1;
              n(!0);
            });
          n(!0);
        },
        onCancel: () => {
          typeof o == "function" && o(), n(!1);
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
    return new Promise((r) => {
      const { onOk: o, onCancel: s, ...a } = n, c = "prompt-" + Math.random().toString(36).slice(2, 11), l = b.getStyle("box", "promptWrapper"), d = b.getStyle("box", "promptText"), m = b.getStyle("box", "promptInput"), u = b.dom.sanitize(e), g = b.dom.sanitize(t), w = `
        <div class="${l}">
          <label for="${c}" class="block ${d}">${u}</label>
          <input id="${c}" type="text" class="${m}" value="${g}">
        </div>
      `;
      this.open("confirm", {
        ...a,
        title: a.title || i.store("i18n")?.t?.("box.prompt") || "Prompt",
        content: w,
        _isTrustedPrompt: !0,
        // Internal flag for prompt content
        onOk: () => {
          const y = typeof o == "function" ? o() : void 0;
          if (y === !1) return !1;
          if (y instanceof Promise)
            return y.then((x) => {
              if (x === !1) return !1;
              const f = document.getElementById(c);
              r(f ? f.value : null);
            });
          const p = document.getElementById(c);
          r(p ? p.value : null);
        },
        onCancel: () => {
          typeof s == "function" && s(), r(null);
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
    let n = "dialog", r = t;
    typeof e == "object" ? (r = e, n = r.type || "dialog") : n = e;
    const o = b.getConfig("box", r), s = o.boxId || Math.random().toString(36).slice(2, 10);
    let a = o.showOkButton, c = o.showCancelButton;
    n === "confirm" ? (a = a !== !1, c = c !== !1) : n === "dialog" && (a = a || !1, c = c || !1);
    const l = {
      id: s,
      type: n,
      title: o.title || "",
      // 外部/内容内标题元素 id（aria-labelledby 优先于 title 与通用兜底名），
      // 供阅读弹窗等「标题在异步内容里」的场景关联真实标题，避免重复展示
      labelledBy: o.labelledBy || "",
      content: o._isTrustedPrompt ? o.content || "" : b.dom.sanitize(o.content || ""),
      // Security: Always sanitize content (trustedHTML removed)
      // Exception: Internal prompt content is pre-escaped
      _isTrustedPrompt: o._isTrustedPrompt || !1,
      width: o.width !== void 0 ? o.width : o.widthClass ? void 0 : o.defaultWidth || 400,
      height: o.height !== void 0 ? o.height : (o.heightClass, void 0),
      widthClass: o.widthClass,
      heightClass: o.heightClass,
      containerClass: o.containerClass,
      fill: o.fill === !0,
      autoFocus: o.autoFocus !== void 0 ? o.autoFocus : !0,
      position: o.position || o.defaultPosition || "center",
      modal: o.modal !== void 0 ? o.modal : !0,
      closeOnEsc: o.closeOnEsc !== void 0 ? o.closeOnEsc : !1,
      closeOnOverlayClick: o.closeOnOverlayClick !== void 0 ? o.closeOnOverlayClick : !1,
      showCloseButton: o.showCloseButton !== void 0 ? o.showCloseButton : !0,
      showOkButton: a,
      showCancelButton: c,
      okText: o.okText || "OK",
      cancelText: o.cancelText || "Cancel",
      okButtonClass: o.okButtonClass,
      cancelButtonClass: o.cancelButtonClass,
      extraButtons: Array.isArray(o.extraButtons) ? o.extraButtons : [],
      onOk: o.onOk,
      onCancel: o.onCancel,
      onBeforeClose: o.onBeforeClose,
      onClose: o.onClose,
      onContentLoad: o.onContentLoad,
      hxUrl: o.hxUrl,
      hxMethod: o.hxMethod || "GET",
      hxHeaders: o.hxHeaders,
      hxTrigger: o.hxTrigger || "load",
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
        const u = this.opener;
        if (!u || u === document.body || !u.isConnected) return;
        const g = document.getElementById("uno-box-outer-container" + this.id), w = document.activeElement;
        (g && w && g.contains(w) || w === document.body) && u.focus();
      },
      /**
       * Close this box dialog
       */
      close() {
        this.onBeforeClose && this.onBeforeClose(), this.isOpen = !1;
        const u = i.store("uno")?.box;
        u && u.syncUnderlayInert(), this.restoreFocus();
        const g = typeof document < "u" ? document.getElementById("uno-box-outer-container" + this.id) : null, w = () => {
          const x = i.store("uno")?.box;
          x && x.remove(this.id), this.onClose && this.onClose();
        };
        if (!g) {
          w();
          return;
        }
        let y = !1;
        const p = (x) => {
          x && x.target !== g || y || (y = !0, g.removeEventListener("transitionend", p), w());
        };
        g.addEventListener("transitionend", p), setTimeout(p, 350);
      }
    };
    this.items.push(l), this.ensureContainer();
    const d = this.items[this.items.length - 1], m = () => {
      d.isOpen = !0, this.scheduleAutoFocus(d, i), this.syncUnderlayInert();
    };
    return i && typeof i.nextTick == "function" && i.nextTick(m), d;
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
   * Esc 只由最顶层打开的实例消费，下层弹窗不随之关闭；
   * 顶层判定与 syncUnderlayInert 同口径（items 顺序 + isOpen）。
   */
  isTopmost(e) {
    const t = this.items.filter((n) => n.isOpen);
    return t.length > 0 && t[t.length - 1] === e;
  },
  /**
   * 存在打开的模态弹窗时遮罩与全屏 wrapper 才拦截背景；
   * 非模态期间背景保持可点击（模板据此释放 wrapper 的 pointer-events）。
   */
  hasModalOpen() {
    return this.items.some((e) => e.isOpen && e.modal !== !1);
  },
  /**
   * B5 多层弹窗：仅顶层可交互/可感知，下层 inert + aria-hidden，
   * 屏幕阅读器与键盘不穿透底层（不改现有 scheduleAutoFocus 契约）。
   * 顶层为非模态时背景与下层均保持可交互，不做 inert 隔离。
   * 模态顶层时另以原生 inert 隔离页面背景，杜绝任何引擎的焦点逃逸。
   */
  syncUnderlayInert() {
    if (typeof document > "u") return;
    const e = this.items.filter((n) => n.isOpen), t = e.length > 0 && e[e.length - 1].modal !== !1;
    e.forEach((n, r) => {
      const o = document.getElementById("uno-box-outer-container" + n.id);
      if (!o) return;
      const s = t && r < e.length - 1;
      o.toggleAttribute("inert", s), s ? o.setAttribute("aria-hidden", "true") : o.removeAttribute("aria-hidden");
    }), this.syncBackgroundInert(t);
  },
  /**
   * 模态顶层时用原生 inert 隔离弹窗管理器（#uno-box-manager）之外的 body
   * 子节点。inert 同时移出焦点序列与无障碍树，背景无法被 Tab 或程序聚焦命中，
   * 不再依赖手写的 Tab 候选过滤来防逃逸。关闭（顶层非模态或无弹窗）时仅恢复
   * 本组件 inert 过的节点，应用原本就 inert 的节点保持原状。
   * @param {boolean} topModal - 当前顶层是否为模态弹窗
   */
  syncBackgroundInert(e) {
    if (!document.body) return;
    const t = this._backgroundInert;
    if (!e) {
      t.forEach((r) => r.removeAttribute("inert")), t.clear();
      return;
    }
    const n = document.getElementById("uno-box-manager");
    Array.from(document.body.children).forEach((r) => {
      r === n || r.hasAttribute("inert") || (r.setAttribute("inert", ""), t.add(r));
    });
  },
  /**
   * Close all open boxes
   */
  closeAll() {
    this.items.forEach((e) => {
      e.onBeforeClose && e.onBeforeClose(), e.onClose && e.onClose();
    }), this.items = [], this.syncUnderlayInert();
  },
  /**
   * Get position CSS classes for 9-point grid
   * @param {string} position - Position name
   * @returns {string} CSS classes
   */
  getPositionClasses(e) {
    return Te(e);
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
    const n = "uno-box-container" + e.id, r = () => {
      if (!e.isOpen) return;
      const s = document.getElementById(n);
      s && (e.autoFocus && this.focusFirstInput(s) || this.focusNeutralTarget(s, e));
    }, o = (s) => {
      typeof requestAnimationFrame == "function" ? requestAnimationFrame(() => requestAnimationFrame(s)) : setTimeout(s, 0);
    };
    t && typeof t.nextTick == "function" ? t.nextTick(() => o(r)) : o(r);
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
    ].join(","), n = e.querySelectorAll(t), r = (o) => {
      for (let s = o; s && s.nodeType === 1; s = s.parentElement) {
        const a = window.getComputedStyle(s);
        if (a.display === "none" || a.visibility === "hidden") return !1;
      }
      return !0;
    };
    for (const o of n)
      if (r(o))
        return o.focus(), !0;
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
    const n = t && t.title ? document.getElementById("uno-box-title-" + t.id) : null, r = n && e.contains(n) ? n : e;
    return r.hasAttribute("tabindex") || r.setAttribute("tabindex", "-1"), r.focus({ preventScroll: !0 }), !0;
  },
  /**
   * Ensure the box container exists in the DOM
   */
  ensureContainer() {
    const e = "uno-box-manager", t = typeof document < "u" && document.getElementById(e);
    if (t)
      return i && typeof i.initTree == "function" && !t._x_dataStack && i.initTree(t), t;
    const r = {
      "x-data": "{}",
      class: "relative",
      style: `z-index: ${b.options?.ui?.layers?.box || S("BOX")};`,
      "x-show": "$store.uno.box.items.length > 0",
      "x-cloak": ""
      // Hide until Alpine is ready
    }, o = je(), s = b.dom.ensureContainer(e, o, r);
    return s && i && typeof i.initTree == "function" && i.initTree(s), s;
  }
}), Ye = (i) => {
  i.directive("uno-trap", (e, { expression: t }, { cleanup: n, evaluateLater: r }) => {
    if (t)
      r(t)((s) => {
        const a = re(e, s || {});
        n(a);
      });
    else {
      const o = re(e);
      n(o);
    }
  }), i.directive("uno-box", (e, { expression: t }, { evaluateLater: n, cleanup: r }) => {
    if (!t) {
      h.warn("[Uno] x-uno-box directive requires a configuration object.");
      return;
    }
    const o = n(t), s = (a) => {
      a.preventDefault();
      const c = i.store("uno")?.box;
      if (!c) {
        h.warn("[Uno] Box store not found. Is Uno initialized?");
        return;
      }
      o((l) => {
        l && typeof l == "object" ? c.open(l) : h.warn("[Uno] x-uno-box directive value must be an object", l);
      });
    };
    e.addEventListener("click", s), r(() => e.removeEventListener("click", s));
  });
};
const Ne = (i, e = {}) => {
  const t = {
    defaultWidth: 400,
    defaultPosition: "center",
    okText: "OK",
    cancelText: "Cancel",
    closeOnEsc: !1,
    closeOnOverlayClick: !1
  };
  b.registerPluginDefaults("box", { ...t, ...e.plugins?.box, ...e.box });
  const n = Ve(i);
  H(i, "box", n), Ye(i);
}, se = {
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
  // Password must be 8-80 chars, contain letters, numbers, and at least 1 special char (@$!%*#?&)
};
let We = 0, Xe = 0;
const W = () => typeof window < "u" && typeof window.matchMedia == "function" && window.matchMedia("(prefers-reduced-motion: reduce)").matches ? "auto" : "smooth", Ke = (i, e) => {
  const t = i.labels && i.labels[0] || i.closest && i.closest("label");
  if (t && t.textContent.trim()) return t.textContent.trim();
  const n = i.getAttribute && i.getAttribute("aria-label");
  if (n && n.trim()) return n.trim();
  const r = i.getAttribute && i.getAttribute("aria-labelledby");
  if (r) {
    const o = document.getElementById(r.trim().split(/\s+/)[0]);
    if (o && o.textContent.trim()) return o.textContent.trim();
  }
  return e;
}, Ze = (i) => ({
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
    const r = this.normalizeRules(t), o = b.getConfig("valid", n);
    this.rules[e] = r.map((s) => ({
      emptyQuit: o.emptyQuit !== !1,
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
    const r = this.rules[e];
    if (!r || r.length === 0) return !0;
    delete this.errors[e], this.syncAriaInvalid(e), this.timers[e] && (clearTimeout(this.timers[e]), delete this.timers[e]);
    const o = t == null ? "" : String(t), s = o.trim() === "";
    for (const a of r) {
      if (n === "input" && (a.trigger === "blur" || a.realtime === !1))
        continue;
      const c = a.type === "required" || a.required === !0;
      if (!(s && !c && a.emptyQuit !== !1)) {
        if (a.debounce > 0 && n === "input")
          return this.validating[e] = !0, new Promise((l) => {
            this.timers[e] = setTimeout(async () => {
              const d = await this.executeRule(e, t, a, o);
              delete this.validating[e], l(!!d);
            }, a.debounce);
          });
        this.validating[e] = !0;
        try {
          if (!await this.executeRule(e, t, a, o))
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
  async executeRule(e, t, n, r) {
    let o = !0, s = "";
    const a = n.type;
    if (a === "required" || n.required)
      r.trim() === "" && (o = !1, s = n.msg || this.getI18nMsg("valid.required") || "This field is required");
    else if (se[a])
      se[a].test(r) || (o = !1, s = n.msg || this.getI18nMsg(`valid.${a}`) || `Invalid ${a} format`);
    else if (a === "length")
      n.min !== void 0 && r.length < n.min ? (o = !1, s = n.msg || this.getI18nMsg("valid.min", { n: n.min }) || `Min length is ${n.min}`) : n.max !== void 0 && r.length > n.max && (o = !1, s = n.msg || this.getI18nMsg("valid.max", { n: n.max }) || `Max length is ${n.max}`);
    else if (a === "range") {
      const c = parseFloat(r);
      isNaN(c) ? (o = !1, s = n.msg || this.getI18nMsg("valid.number") || "Must be a number") : n.min !== void 0 && c < n.min ? (o = !1, s = n.msg || this.getI18nMsg("valid.rangeMin", { n: n.min }) || `Min value is ${n.min}`) : n.max !== void 0 && c > n.max && (o = !1, s = n.msg || this.getI18nMsg("valid.rangeMax", { n: n.max }) || `Max value is ${n.max}`);
    } else if (a === "pattern" || n.pattern)
      (n.pattern instanceof RegExp ? n.pattern : new RegExp(n.pattern)).test(r) || (o = !1, s = n.msg || this.getI18nMsg("valid.pattern") || "Invalid format");
    else if (a === "match" || n.match) {
      const c = n.match || n.selector || n.target;
      if (!c)
        return !0;
      const l = document.querySelector(c.startsWith("#") ? c : `[name="${c}"]`) || document.getElementById(c);
      l && r !== l.value && (o = !1, s = n.msg || this.getI18nMsg("valid.match") || "Fields do not match");
    } else if (a === "custom" || a === "validator" || n.fn || this.customValidators[a]) {
      let c = n.fn || n.name;
      if (typeof c == "string" ? c = this.customValidators[c] : !c && this.customValidators[a] && (c = this.customValidators[a]), typeof c == "function")
        try {
          const l = await c(t, e, n);
          l === !1 ? (o = !1, s = n.msg || "Validation failed") : typeof l == "string" ? (o = !1, s = l) : typeof l == "object" && l.valid === !1 && (o = !1, s = l.message || n.msg || "Validation failed");
        } catch (l) {
          h.error(`[Valid] Error in custom validator for ${e}`, l), o = !1, s = "Validation error";
        }
      else
        h.warn(`[Valid] Custom validator not found: ${a || n.fn}`);
    }
    return o || (this.errors[e] = s, this.syncAriaInvalid(e)), o;
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
    let n = !0, r = null;
    const o = [];
    this.removeErrorSummary(e);
    for (const s of t) {
      if (s.disabled) continue;
      const a = s.name || s.id;
      this.rules[a] && (await this.validate(a, s.value, "submit") || (n = !1, r || (r = s), o.push({ input: s, fieldName: a })));
    }
    if (!n && r) {
      const s = b.getConfig("valid") || {};
      if (o.length >= 2) {
        const a = this.renderErrorSummary(e, o);
        a && s.scrollOnError !== !1 && (a.focus({ preventScroll: !0 }), a.scrollIntoView({ behavior: W(), block: "center" }));
      } else s.scrollOnError !== !1 && (r.focus({ preventScroll: !0 }), r.scrollIntoView({ behavior: W(), block: "center" }));
    }
    return n;
  },
  /**
   * Renders the top-of-form error summary for multi-field failures
   * (E4 / WCAG 3.3.1 / 3.3.3, GOV.UK pattern): a focusable `role="alert"`
   * container labelled by its title, with one link per failed field. Each link
   * carries the field name plus its error message and anchors the real control
   * (`input.id`, minting a stable id when the control has none) so clicking
   * focuses the control itself rather than the non-focusable inline error div.
   * Focus and scroll are separate steps and scrolling honors
   * prefers-reduced-motion. Visual classes (`ui-error-summary*`) are styled by
   * the app stylesheet (website tailwind.css), mirroring the `is-invalid` pattern.
   *
   * @param {HTMLFormElement} form - The form receiving the summary.
   * @param {Array<{input: HTMLElement, fieldName: string}>} failedFields
   * @returns {HTMLElement|null} The summary element, or null when unsupported.
   */
  renderErrorSummary(e, t) {
    if (typeof document > "u" || !e) return null;
    const n = document.createElement("div");
    n.className = "ui-error-summary", n.dataset.validSummary = "true", n.setAttribute("role", "alert"), n.setAttribute("tabindex", "-1");
    const r = `valid-error-summary-title-${++We}`;
    n.setAttribute("aria-labelledby", r);
    const o = document.createElement("p");
    o.className = "ui-error-summary-title", o.id = r, o.textContent = this.getI18nMsg("valid.summary_title") || "Please fix the following problems", n.appendChild(o);
    const s = document.createElement("ul");
    return s.className = "ui-error-summary-list", t.forEach(({ input: a, fieldName: c }) => {
      a.id || (a.id = `valid-field-${++Xe}`);
      const l = a.id, d = this.errors[c] || "", m = Ke(a, c), u = document.createElement("a");
      u.className = "ui-error-summary-link", u.setAttribute("href", `#${l}`), u.textContent = m && m !== d ? `${m}: ${d}` : d, u.addEventListener("click", (w) => {
        w.preventDefault();
        const y = document.getElementById(l) || a;
        y && (y.focus({ preventScroll: !0 }), y.scrollIntoView({ behavior: W(), block: "center" }));
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
    const n = i.store("i18n");
    return n && typeof n.t == "function" ? n.t(e, t) : null;
  }
}), te = "data-uno-counter-wrapper", Ge = (i) => i === "flex" || i === "inline-flex" || i === "grid" || i === "inline-grid", Je = (i) => {
  const e = i.parentElement;
  if (!e) return null;
  if (e.hasAttribute(te)) return e;
  const t = window.getComputedStyle(e);
  if (Ge(t.display))
    return t.position === "static" && e.tagName !== "FORM" && e.tagName !== "BODY" && (e.style.position = "relative"), e;
  const n = i.ownerDocument.createElement("div");
  n.setAttribute(te, ""), n.style.position = "relative";
  const r = window.getComputedStyle(i).display;
  return (r === "inline" || r === "inline-block") && (n.style.display = "inline-block"), e.insertBefore(n, i), n.appendChild(i), n;
}, Qe = (i, e) => {
  if (!e || !e.hasAttribute(te)) return;
  const t = e.parentNode;
  t && (t.insertBefore(i, e), t.removeChild(e));
}, et = (i) => `${i}-error`, ae = (i, e) => {
  i.id || (i.id = et(e)), i.hasAttribute("role") || i.setAttribute("role", "alert");
}, X = (i, e) => {
  if (!e.id) return;
  const t = (i.getAttribute("aria-describedby") || "").split(/\s+/).filter(Boolean);
  t.includes(e.id) || t.push(e.id), i.setAttribute("aria-describedby", t.join(" "));
}, tt = (i) => {
  i.directive("uno-valid", (e, { expression: t }, { evaluateLater: n, cleanup: r }) => {
    if (e.tagName === "FORM") return;
    const o = e.name || e.id;
    if (!o) {
      h.warn("[Uno Valid] Element has no name or id:", e);
      return;
    }
    const s = Uno.getConfig("valid"), a = Uno.options.ui?.valid || {};
    if (s.autoError !== !1) {
      const u = `$store.uno.valid.errors['${o}']`, g = `[x-text="${u.replace(/"/g, '\\"')}"]`, w = e.form || e.closest("form"), y = w ? w.querySelector(g) : document.querySelector(g);
      if (y)
        ae(y, o), X(e, y);
      else {
        let p;
        if (typeof s.createErrorElement == "function")
          p = s.createErrorElement(o, u);
        else {
          p = document.createElement("div"), p.setAttribute("x-text", u), p.setAttribute("x-show", u);
          const x = s.errorClass || a.helperText || "text-xs text-red-500 mt-1";
          p.className = x;
        }
        p && (ae(p, o), e.after(p), i.initTree(p), X(e, p));
      }
    }
    const c = t ? n(t) : () => {
    }, l = (u) => {
      i.store("uno").valid.addRules(o, u);
    };
    t ? c((u) => {
      l(u);
      const g = i.store("uno").valid;
      if (g && g.normalizeRules) {
        const y = g.normalizeRules(u).find(
          (p) => p.type === "length" || p.maxLength !== void 0 || p.max !== void 0
        );
        if (y) {
          const p = y.max !== void 0 ? y.max : y.maxLength;
          if (y.stopOnMax === !0 && p !== void 0 && e.setAttribute("maxlength", p), y.showCounter === !0 && p !== void 0) {
            const x = `${o}-counter`;
            let f = document.getElementById(x);
            const v = Je(e);
            f || (f = document.createElement("span"), f.id = x, f.className = "absolute bottom-1 end-2 text-xs text-gray-600 dark:text-gray-400 pointer-events-none transition-colors duration-200"), f.parentElement !== v && e.after(f), X(e, f);
            const C = () => {
              const E = e.value.length;
              f.textContent = `${E}/${p}`;
              const k = g.getI18nMsg ? g.getI18nMsg("valid.counter", { current: E, max: p }) : null;
              f.setAttribute("aria-label", k || `${E}/${p}`), E > p ? (f.classList.add("text-red-600", "dark:text-red-400"), f.classList.remove("text-gray-600", "dark:text-gray-400")) : (f.classList.remove("text-red-600", "dark:text-red-400"), f.classList.add("text-gray-600", "dark:text-gray-400"));
            };
            C(), e.addEventListener("input", C), r(() => {
              e.removeEventListener("input", C), f && f.parentNode && f.parentNode.removeChild(f), Qe(e, v);
            });
          }
        }
      }
    }) : l([]);
    const d = () => {
      i.store("uno").valid.validate(o, e.value, "input");
    }, m = () => {
      i.store("uno").valid.validate(o, e.value, "blur");
    };
    e.addEventListener("input", d), e.addEventListener("blur", m), r(() => {
      e.removeEventListener("input", d), e.removeEventListener("blur", m);
    });
  }), typeof document < "u" && (document.body.addEventListener("htmx:confirm", async (e) => {
    const t = e.detail.elt, n = t.tagName === "FORM" ? t : t.closest("form");
    if (n && (n.hasAttribute("x-uno-valid") || n.querySelector("[x-uno-valid]"))) {
      e.preventDefault();
      try {
        await i.store("uno").valid.validateForm(n) && e.detail.issueRequest(!0);
      } catch (r) {
        h.error("[Valid] Form validation failed unexpectedly, request blocked", r);
      }
    }
  }), document.body.addEventListener("htmx:afterRequest", (e) => {
    const t = K(e);
    if (t && t.status === 422)
      try {
        const n = JSON.parse(t.text);
        n.errors && (i.store("uno").valid.showErrors(n.errors), h.info("[Valid] Server validation errors applied", n.errors));
      } catch {
      }
  }));
};
const nt = (i, e = {}) => {
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
  const n = Ze(i);
  H(i, "valid", n), tt(i);
}, ot = (i) => ({
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
        const { group: n, name: r } = this.registry[t];
        h.info(`[Tab] checkHash match: ${r}`), this.activate({ group: n, name: r, updateHash: !1 });
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
    const r = this.order[e];
    if (!r || r.length === 0) return null;
    const o = r.indexOf(t);
    return n === "Home" ? r[0] : n === "End" ? r[r.length - 1] : o === -1 ? r[0] : n === "ArrowRight" ? r[(o + 1) % r.length] : n === "ArrowLeft" ? r[(o - 1 + r.length) % r.length] : null;
  }
}), it = (i) => {
  i.directive("uno-tab", (e, { expression: t }, { evaluateLater: n, effect: r, cleanup: o }) => {
    if (!t) {
      h.warn("[Uno] x-uno-tab directive requires a configuration object.");
      return;
    }
    const s = n(t);
    let a = null;
    const c = (m) => typeof m == "object" && m !== null ? { ...m } : typeof m == "string" ? { name: m } : null;
    s((m) => {
      const u = c(m);
      if (!u) {
        h.warn("[Uno] x-uno-tab directive value must be a string or object", m);
        return;
      }
      u.group || (u.group = "default"), a = u, u.name && (i.store("uno").tab.register(u.name, u.group), e.setAttribute("role", "tab"), e.id || (e.id = "tab-" + u.name), e.setAttribute("aria-controls", u.panel || "panel-" + u.name));
    }), r(() => {
      if (!a || !a.name) return;
      const m = i.store("uno").tab.isActive(a.name, a.group);
      e.setAttribute("aria-selected", m ? "true" : "false"), e.setAttribute("tabindex", m ? "0" : "-1");
    });
    const l = (m) => {
      m.preventDefault(), s((u) => {
        const g = c(u);
        g && (g.group || (g.group = "default"), a = g, g.name && i.store("uno").tab.activate(g));
      });
    }, d = (m) => {
      if (m.key === " " || m.key === "Spacebar") {
        m.preventDefault(), m.repeat || e.click();
        return;
      }
      if (!["ArrowRight", "ArrowLeft", "Home", "End"].includes(m.key) || !a || !a.name) return;
      m.preventDefault();
      const u = i.store("uno").tab, g = e.closest("[dir]") || document.documentElement, y = typeof getComputedStyle == "function" && getComputedStyle(g).direction === "rtl" ? m.key === "ArrowRight" ? "ArrowLeft" : m.key === "ArrowLeft" ? "ArrowRight" : m.key : m.key, p = u.rotate(a.group, a.name, y);
      if (!p) return;
      const x = document.getElementById("tab-" + p);
      x && (x.focus(), x.click());
    };
    e.addEventListener("click", l), e.addEventListener("keydown", d), o(() => {
      e.removeEventListener("click", l), e.removeEventListener("keydown", d);
    });
  });
};
const rt = (i, e = {}) => {
  const t = {
    defaultTab: "",
    hash: !0,
    activeClass: "is-active",
    inactiveClass: ""
  };
  b.registerPluginDefaults("tab", { ...t, ...e.plugins?.tab, ...e.tab });
  const n = ot();
  H(i, "tab", n), it(i);
};
b.dom = V;
typeof window < "u" && (window.Uno = b);
b.use(Se).use(Pe).use(qe).use(Ne).use(nt).use(rt);
if (typeof window < "u") {
  const i = () => {
    window.Alpine && !b.initialized && b.init(window.Alpine);
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
  V as DomManager,
  b as Uno,
  Ne as box,
  b as default,
  qe as notify,
  rt as tab,
  Pe as tooltip,
  nt as valid
};
