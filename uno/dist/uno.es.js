import z from "alpinejs";
let A = !1, P = null;
const v = {
  /**
   * Set debug mode
   * @param {boolean} debug
   */
  setDebug(n) {
    A = n;
  },
  /**
   * Set custom error reporter
   * @param {Function} reporter - (error, ...args) => void
   */
  setReporter(n) {
    P = n;
  },
  /**
   * Log info message
   * @param {...any} args - Arguments to log
   */
  info(...n) {
    A && typeof console < "u" && typeof console.log == "function" && console.log("[Uno Info]", ...n);
  },
  /**
   * Log warning message
   * @param {...any} args - Arguments to log
   */
  warn(...n) {
    A && typeof console < "u" && typeof console.warn == "function" && console.warn("[Uno Warn]", ...n);
  },
  /**
   * Log error message
   * @param {...any} args - Arguments to log
   */
  error(...n) {
    if (P)
      try {
        P(...n);
      } catch (s) {
        A && typeof console < "u" && typeof console.error == "function" && console.error("[Uno Error] Reporter failed:", s);
      }
    A && typeof console < "u" && typeof console.error == "function" && console.error("[Uno Error]", ...n);
  }
};
function N(n) {
  return String(n ?? "").replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;").replace(/'/g, "&#39;");
}
const q = /* @__PURE__ */ new WeakMap();
function _(n, s) {
  if (!n) return !1;
  const a = String(s);
  let l = q.get(n);
  return l || (l = /* @__PURE__ */ new Set(), q.set(n, l)), l.has(a) ? !1 : (l.add(a), !0);
}
function j(n) {
  return n ? n.replace(/<script\b[^>]*>([\s\S]*?)<\/script>/gim, "").replace(/<img\b[^>]*>/gim, "").replace(/\s+on[a-z]+="[^"]*"/gim, "").replace(/\s+on[a-z]+='[^']*'/gim, "").replace(/\s+href="javascript:[^"]*"/gim, "").replace(/\s+href='javascript:[^']*'/gim, "").replace(/javascript:/gim, "") : "";
}
function X(n, s, { trusted: a = !1, alpine: l } = {}) {
  n && (n.innerHTML = a ? String(s ?? "") : j(String(s ?? "")), l?.initTree?.(n));
}
const I = {
  /**
   * Ensure a container element exists in the DOM.
   * @param {string} id - Element ID
   * @param {string} html - Inner HTML content
   * @param {Object} attributes - Attributes to set
   * @param {HTMLElement} parent - Parent element (default: document.body)
   */
  ensureContainer(n, s = "", a = {}, l = document.body) {
    const o = document.getElementById(n);
    if (o) return o;
    try {
      const t = document.createElement("div");
      return t.id = n, window.Uno?.options?.csp?.nonce && !a.nonce && (a.nonce = window.Uno.options.csp.nonce), Object.entries(a).forEach(([i, r]) => {
        t.setAttribute(i, r);
      }), s && (t.innerHTML = s), l ? l.appendChild(t) : document.documentElement.appendChild(t), t;
    } catch (t) {
      return v.error(`DomManager: Failed to inject container #${n}`, t), null;
    }
  },
  /**
   * Inject CSS styles with CSP nonce support.
   * @param {string} id - Style element ID
   * @param {string} css - CSS content
   */
  injectStyles(n, s) {
    if (!document.getElementById(n))
      try {
        const a = document.createElement("style");
        a.id = n, window.Uno?.options?.csp?.nonce && a.setAttribute("nonce", window.Uno.options.csp.nonce), a.textContent = s, document.head.appendChild(a);
      } catch (a) {
        v.error(`DomManager: Failed to inject styles #${n}`, a);
      }
  },
  // Exports
  setHTML: X,
  // Secure HTML setting
  escape: N,
  // HTML escaping
  sanitize: j
  // HTML sanitization
};
function $(n) {
  return Object.prototype.toString.call(n) === "[object Object]";
}
function S(n, s) {
  if (!$(n)) return $(s) ? { ...s } : s;
  if (!$(s)) return s === void 0 ? { ...n } : s;
  const a = { ...n };
  for (const [l, o] of Object.entries(s)) {
    if ($(o) && $(a[l])) {
      a[l] = S(a[l], o);
      continue;
    }
    if ($(o)) {
      a[l] = S({}, o);
      continue;
    }
    if (typeof o == "string" && o.trim().indexOf("+") === 0 && typeof a[l] == "string") {
      a[l] = `${a[l]} ${o.trim().slice(1)}`.trim();
      continue;
    }
    a[l] = o;
  }
  return a;
}
class Y {
  constructor(s = {}) {
    this.version = "0.6.0", this.options = s, this.pluginDefaults = {}, this.plugins = [], this.initialized = !1, this._htmxRequestHandlers = [], this._htmxActiveXhrs = typeof WeakSet == "function" ? /* @__PURE__ */ new WeakSet() : null, this.htmx = {
      available: () => typeof window < "u" && typeof window.htmx?.ajax == "function",
      ajax: (o, t, i = {}) => typeof window > "u" || typeof window.htmx?.ajax != "function" ? Promise.reject(new Error("htmx.ajax is not available")) : window.htmx.ajax(o, t, i),
      fragment: async (o, t, i = {}) => {
        if (typeof document > "u") throw new Error("document is not available");
        const r = document.createElement("div");
        let d = "";
        const e = (u) => {
          u.preventDefault(), d = u.detail.serverResponse, u.detail.shouldSwap = !1;
        };
        r.addEventListener("htmx:beforeSwap", e);
        try {
          const u = i.swap || "innerHTML";
          return await this.htmx.ajax(o, t, { ...i, target: r, swap: u }), d;
        } finally {
          r.removeEventListener("htmx:beforeSwap", e), r.remove();
        }
      },
      onRequest: (o = {}) => {
        if (!o || typeof o != "object" || (this._htmxRequestHandlers.push(o), typeof window > "u" || typeof document > "u") || !window.htmx || !_(document, "uno:htmx:request-hooks")) return;
        const t = (e) => e?.detail?.xhr, i = (e, u) => {
          for (const f of this._htmxRequestHandlers) {
            const c = f?.[e];
            if (typeof c == "function")
              try {
                c(u);
              } catch (g) {
                v.error("Uno.htmx.onRequest handler failed", g);
              }
          }
        }, r = (e) => {
          const u = t(e);
          if (u && this._htmxActiveXhrs) {
            if (this._htmxActiveXhrs.has(u)) return;
            this._htmxActiveXhrs.add(u);
          }
          i("start", e);
        }, d = (e) => {
          const u = t(e);
          if (u && this._htmxActiveXhrs) {
            if (!this._htmxActiveXhrs.has(u)) return;
            this._htmxActiveXhrs.delete(u);
          }
          i("end", e);
        };
        document.addEventListener("htmx:beforeRequest", r), document.addEventListener("htmx:afterRequest", d), document.addEventListener("htmx:responseError", d), document.addEventListener("htmx:sendError", d), document.addEventListener("htmx:timeout", d), document.addEventListener("htmx:xhr:progress", (e) => i("progress", e));
      },
      // Global HTMX event binding for notify functionality - progress bar only
      bindNotifyToHtmx: () => {
        typeof window > "u" || typeof document > "u" || !window.htmx || !window.Alpine || _(document, "uno:htmx:bind-notify") && (document.addEventListener("htmx:beforeRequest", () => {
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
        layers: { boxBase: 1e4, tooltipBase: 2e4, notifyBase: 3e4 },
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
    this.options = S(l, s);
  }
  /**
   * ✅ Unified configuration retrieval method
   * Supports 4-level merge: Plugin Defaults -> Global User Config -> Global Plugin Config -> Direct Global Config -> Instance Config
   * Note: We use mergeDeep to ensure deep merging of options.
   *
   * Level 1: Plugin Defaults (this.pluginDefaults[componentName])
   * Level 2: Global Plugin Config (this.options.plugins[componentName])
   * Level 3: Direct Global Config (this.options[componentName]) - Legacy/Shorthand
   * Level 4: Instance Config (instanceOpts)
   */
  getConfig(s, a = {}) {
    const l = this.options || {};
    let o = this.pluginDefaults[s] ? { ...this.pluginDefaults[s] } : {};
    return l.plugins && l.plugins[s] && (o = S(o, l.plugins[s])), l[s] && (o = S(o, l[s])), a && (o = S(o, a)), o;
  }
  /**
   * ✅ Register plugin defaults
   */
  registerPluginDefaults(s, a) {
    this.pluginDefaults[s] = a;
  }
  _getUiValue(s, a) {
    const l = this.options.ui?.[s];
    if (!l) return;
    if (typeof a != "string" || a.indexOf(".") === -1) return l?.[a];
    const o = a.split(".").filter(Boolean);
    let t = l;
    for (const i of o) {
      if (t == null) return;
      t = t[i];
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
    const o = this._getUiValue(s, a);
    return o === void 0 ? this._resolveUiTokens(l) : typeof o == "string" && o.trim().indexOf("+") === 0 ? this._resolveUiTokens(`${l} ${o.trim().slice(1)}`.trim()) : this._resolveUiTokens(o);
  }
  /**
   * Set global configuration (before or after init)
   * @param {Object} config - Global configuration object
   */
  setConfig(s) {
    return this.options = S(this.options, s), this.initialized && v.info("Uno configuration updated", this.options), this;
  }
  /**
   * Register a plugin
   * @param {Function} plugin - Plugin function (Alpine, options) => void
   * @param {Object} options - Plugin options
   * @returns {UnoInstance}
   */
  use(s, a = {}) {
    return this.initialized && v.warn("Uno.use() called after initialization. Plugin might not work correctly."), this.plugins.push({ plugin: s, options: a }), this;
  }
  /**
   * Initialize Uno
   * @param {Object|Function} alpineOrOptions - Alpine instance or Configuration options
   * @param {Object} options - Configuration options (if first arg is Alpine)
   * @returns {Object} Alpine instance
   */
  init(s = {}, a = {}) {
    if (this.initialized)
      return v.warn("Uno has already been initialized."), z;
    let l = z, o = a;
    return !l && typeof window < "u" && window.Alpine && (l = window.Alpine), s && (s.version || typeof s.start == "function" || typeof s.store == "function") ? l = s : o = s || a, this.options = S(this.options, o), v.setDebug(this.options.debug !== !1), this.options.onError && v.setReporter(this.options.onError), v.info(`Uno v${this.version} initializing...`), typeof window < "u" && typeof document < "u" && window.htmx && typeof l?.initTree == "function" && _(document, "uno:htmx-afterSwap:initTree") && document.addEventListener("htmx:afterSwap", (t) => {
      const i = t?.detail?.elt;
      i && l.initTree(i);
    }), this.plugins.forEach(({ plugin: t, options: i }) => {
      try {
        const r = S(this.options, i);
        t(l, r);
      } catch (r) {
        v.error("Failed to load plugin", r);
      }
    }), this.initialized = !0, v.info("Uno initialized"), l;
  }
}
const b = new Y(), Z = {
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
}, K = {
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
}, G = {
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
const W = {
  en: Z,
  "zh-CN": K,
  "zh-TW": G
}, R = Object.keys(W);
function J() {
  if (typeof localStorage < "u") {
    const n = localStorage.getItem("uno.lang");
    if (n && R.includes(n)) return n;
  }
  if (typeof navigator < "u") {
    const n = navigator.language, s = n.startsWith("zh") ? n.includes("TW") || n.includes("Hant") ? "zh-TW" : "zh-CN" : "en";
    return R.includes(s) ? s : "en";
  }
  return "en";
}
function Q(n, s) {
  const a = s.split(".");
  let l = n;
  for (const o of a) {
    if (!l || typeof l != "object" || !(o in l)) return;
    l = l[o];
  }
  return l;
}
function F(n, s) {
  for (const a in s)
    s[a] && typeof s[a] == "object" ? (n[a] || Object.assign(n, { [a]: {} }), F(n[a], s[a])) : Object.assign(n, { [a]: s[a] });
}
function tt(n) {
  n.store("i18n", {
    locale: J(),
    messages: W,
    loading: !1,
    init() {
      typeof document < "u" && (document.documentElement.lang = this.locale);
    },
    t(s, a = {}) {
      const l = this.locale, o = this.messages[l] || this.messages.en, t = Q(o, s);
      return (typeof t == "string" ? t : s).replace(/\{(\w+)\}/g, (r, d) => d in a ? String(a[d]) : `{${d}}`);
    },
    setLocale(s) {
      if (!this.messages[s]) {
        v.warn(`[Uno] Language ${s} not found. Please load it first using loadLocale().`);
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
          const o = await l.json();
          return this.extend({ [s]: o }), !0;
        } catch (l) {
          throw v.error(`[Uno] Failed to load locale ${s}`, l), l;
        } finally {
          this.loading = !1;
        }
      }
    },
    extend(s) {
      for (const a in s)
        this.messages[a] ? F(this.messages[a], s[a]) : this.messages[a] = s[a];
    }
  }), n.directive("t", (s, { expression: a, modifiers: l }, { effect: o, evaluateLater: t }) => {
    const i = () => n.store("i18n");
    o(() => {
      const r = i();
      r.locale;
      const d = a, e = r.t(d);
      l.length > 0 ? l.forEach((u) => {
        s.setAttribute(u, e);
      }) : s.textContent = e;
    });
  });
}
const et = (n) => {
  const s = n || "center", a = {
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
}, ot = (n) => {
  const s = n || "top-right", a = {
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
}, nt = (n) => {
  const s = n || "top-right";
  let a = "flex-col", l = "items-end", o = "space-y-2";
  return s.includes("bottom") && (a = "flex-col-reverse", o = "space-y-reverse space-y-2"), s.includes("left") ? l = "items-start" : s.includes("center") ? s === "center" ? l = "items-center justify-center" : l = "items-center" : l = "items-end", { flexDir: a, alignItems: l, spaceY: o };
};
function it(n, s, a) {
  const l = n.getBoundingClientRect(), o = s.getBoundingClientRect(), t = o.width, i = o.height, r = l.width, d = l.height, e = { top: l.top, left: l.left }, u = window.innerWidth, f = window.innerHeight, c = 0, g = 0, h = f, x = u, m = (k, L) => {
    if (L) {
      if (k + t > x)
        return x - t - 10;
      if (k < g)
        return g + 10;
    } else {
      if (k + i > h)
        return h - i - 10;
      if (k < c)
        return c + 10;
    }
    return k;
  };
  let w, y, E = a;
  switch (a) {
    case "left":
      y = e.left - t - 5, y < g && (y = e.left + r + 5, E = "right"), w = m(e.top + d / 2 - i / 2, !1);
      break;
    case "top":
      w = e.top - i - 8, w < c && (w = e.top + d + 8, E = "bottom"), y = e.left + (r - t) / 2, y = m(y, !0);
      break;
    case "right":
      y = e.left + r + 5, y + t > x && (y = e.left - t - 5, E = "left"), w = m(e.top + d / 2 - i / 2, !1);
      break;
    case "bottom":
    default:
      w = e.top + d + 8, w + i > h && (w = e.top - i - 8, E = "top"), y = e.left + (r - t) / 2, y = m(y, !0);
      break;
  }
  return { top: w, left: y, position: E };
}
function H(n, s, a) {
  if (!n || typeof n.store != "function") {
    console.error("[Uno] Alpine.store is not available");
    return;
  }
  const l = n.store("uno");
  l ? n.store("uno", {
    ...l,
    [s]: a
  }) : n.store("uno", { [s]: a }), n.store(s) || n.store(s, n.store("uno")?.[s]);
}
const st = (n, s = {}) => {
  const a = {
    placement: "top",
    trigger: "hover",
    delay: 300,
    offset: 8,
    arrow: !0,
    clickOutside: !0
  };
  b.registerPluginDefaults("tooltip", {
    ...a,
    ...s.plugins?.tooltip,
    ...s.tooltip
  });
  let l = 0, o = null;
  const t = () => {
    l || typeof requestAnimationFrame == "function" && (l = requestAnimationFrame(() => {
      l = 0;
      const e = n.store("uno")?.tooltip;
      e && e.updatePosition();
    }));
  }, i = () => {
    if (o || typeof window > "u" || typeof document > "u") return;
    const e = () => t(), u = () => t(), f = (g) => {
      if (g.key !== "Escape") return;
      const h = n.store("uno")?.tooltip;
      h?.visible && h.hide();
    }, c = (g) => {
      const h = n.store("uno")?.tooltip;
      if (!h?.visible || h.triggerMode !== "click" || h.clickOutside !== !0) return;
      const x = g.target, m = document.getElementById("uno-tooltip-content"), w = h.triggerEl && h.triggerEl.contains(x), y = m && m.contains(x);
      !w && !y && h.hide();
    };
    window.addEventListener("scroll", e, !0), window.addEventListener("resize", u), window.addEventListener("keydown", f), window.addEventListener("click", c), o = () => {
      window.removeEventListener("scroll", e, !0), window.removeEventListener("resize", u), window.removeEventListener("keydown", f), window.removeEventListener("click", c), o = null;
    };
  };
  H(n, "tooltip", {
    visible: !1,
    content: "",
    placement: "top",
    top: -9999,
    left: -9999,
    zIndex: 2e4,
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
      v.info("Tooltip show called", e);
      const u = b.getConfig("tooltip", e);
      if (v.info("Tooltip resolved config", u), this.triggerEl = u.element || this.triggerEl, u.target) {
        const f = document.querySelector(u.target);
        f ? this.content = j(f.innerHTML) : (v.warn(`Tooltip target element not found: ${u.target}`), this.content = u.content || "");
      } else
        this.content = u.content || "";
      this.placement = u.placement, this.triggerMode = u.trigger, this.clickOutside = u.clickOutside, this.offset = u.offset, this.arrow = u.arrow !== !1, this.zIndex = b.options?.ui?.layers?.tooltip || 2e4, this.hxUrl = u.hxUrl, this.hxHeaders = u.hxHeaders, this.hxTrigger = u.hxTrigger || "load", this.visible = !0, this.ensureContainer(), i(), n.nextTick(() => {
        this.updatePosition();
      });
    },
    hide() {
      this.visible = !1, this.triggerEl = null, o && (o(), o = null), l && (cancelAnimationFrame(l), l = 0);
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
      const u = it(this.triggerEl, e, this.placement, {
        offset: this.offset,
        arrow: this.arrow
      });
      this.top = u.top, this.left = u.left, this.placement = u.placement, u.arrow && (this.arrowPosition = u.arrow);
    },
    ensureContainer() {
      const e = "uno-tooltip-manager";
      if (document.getElementById(e)) return;
      const u = b.options?.ui?.layers?.tooltip || 2e4;
      I.ensureContainer(e, `
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
        style: `position: absolute; top: 0; left: 0; width: 0; height: 0; z-index: ${u};`
      });
      const c = document.getElementById(e);
      c && n && !c._x_dataStack && n.initTree(c);
    }
  });
  const d = (e, { expression: u, modifiers: f }, { evaluateLater: c, effect: g, cleanup: h }) => {
    const x = u ? c(u) : null, m = (p) => {
      const T = typeof p == "object" && p !== null ? p : p ? { content: p } : {}, C = {};
      return e.dataset.unoContent && (C.content = e.dataset.unoContent), e.dataset.unoTarget && (C.target = e.dataset.unoTarget), e.dataset.unoTheme && (C.theme = e.dataset.unoTheme), e.dataset.unoPlacement && (C.placement = e.dataset.unoPlacement), e.dataset.unoTrigger && (C.trigger = e.dataset.unoTrigger), e.dataset.unoOffset && (C.offset = parseInt(e.dataset.unoOffset)), e.dataset.unoArrow && (C.arrow = e.dataset.unoArrow === "true"), e.dataset.unoClickOutside && (C.clickOutside = e.dataset.unoClickOutside === "true"), e.hasAttribute("hx-get") && (C.hxUrl = e.getAttribute("hx-get")), e.hasAttribute("hx-headers") && (C.hxHeaders = e.getAttribute("hx-headers")), e.hasAttribute("hx-trigger") && (C.hxTrigger = e.getAttribute("hx-trigger")), f.includes("top") && (C.placement = "top"), f.includes("bottom") && (C.placement = "bottom"), f.includes("left") && (C.placement = "left"), f.includes("right") && (C.placement = "right"), f.includes("click") && (C.trigger = "click"), f.includes("hover") && (C.trigger = "hover"), f.includes("focus") && (C.trigger = "focus"), { ...S(T, C), element: e };
    }, w = (p) => {
      const T = n.store("uno")?.tooltip;
      T && T.show(m(p));
    }, y = (p) => {
      const T = n.store("uno")?.tooltip;
      T && (T.visible && T.triggerEl === e ? T.hide() : w(p));
    };
    let E = null;
    x && g(() => {
      x((p) => {
        E = p;
      });
    });
    const k = () => {
      const p = m(E);
      b.getConfig("tooltip", p).trigger === "hover" && w(E);
    }, L = () => {
      const p = n.store("uno")?.tooltip;
      p && p.triggerMode === "hover" && p.triggerEl === e && p.scheduleHide();
    }, U = (p) => {
      const T = m(E);
      b.getConfig("tooltip", T).trigger === "click" && (p.preventDefault(), y(E));
    }, M = () => {
      const p = m(E);
      b.getConfig("tooltip", p).trigger === "focus" && w(E);
    }, O = () => {
      const p = n.store("uno")?.tooltip;
      p && p.triggerMode === "focus" && p.triggerEl === e && p.hide();
    };
    e.addEventListener("mouseenter", k), e.addEventListener("mouseleave", L), e.addEventListener("click", U), e.addEventListener("focus", M), e.addEventListener("blur", O), h(() => {
      e.removeEventListener("mouseenter", k), e.removeEventListener("mouseleave", L), e.removeEventListener("click", U), e.removeEventListener("focus", M), e.removeEventListener("blur", O);
      const p = n.store("uno")?.tooltip;
      p && p.triggerEl === e && p.hide();
    });
  };
  n.directive("uno-tooltip", d), n.directive("tooltip", d);
}, B = {
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
}, D = (n) => {
  if (!n) return B.BASE;
  const s = n.toUpperCase();
  return B[s] !== void 0 ? B[s] : B.BASE;
}, rt = (n, s = {}) => {
  const a = (o = {}) => b.getConfig("notify", o), l = {
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
    onContentLoaded(o, t) {
      this._onHtmxEnd(t);
      const i = o.innerHTML || o.textContent || "", r = this._findNotificationById(t.id);
      r && (r.content = i), t.onContentLoad && t.onContentLoad(o), window.Alpine && window.Alpine.initTree && window.Alpine.initTree(o);
    },
    // Helper method to find notification by ID
    _findNotificationById(o) {
      for (const t of Object.keys(this.groups)) {
        const r = this.groups[t].find((d) => d.id === o);
        if (r) return r;
      }
      return null;
    },
    onHtmxError(o, t) {
      this._onHtmxEnd(t), console.error("HTMX Error:", o), this.error("Content load failed"), t.onError && t.onError(o);
    },
    _onHtmxBegin(o) {
      !o || o.__unoHtmxActive || (o.__unoHtmxActive = !0, this.progressRequests++, this.progressRequests === 1 && this.startProgress());
    },
    _onHtmxEnd(o) {
      !o || !o.__unoHtmxActive || (o.__unoHtmxActive = !1, this.progressRequests > 0 && this.progressRequests--, this.progressRequests === 0 && this.stopProgress());
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
    getContainerStyle(o) {
      return {
        zIndex: b.options?.ui?.layers?.notifyBase || D("TOAST")
      };
    },
    getContainerClass(o) {
      const t = b.getStyle("notify", "container"), i = ot(o), r = nt(o);
      return [t, i, r.flexDir, r.alignItems].join(" ");
    },
    ensureContainer() {
      const i = I.ensureContainer("uno-notify-container", `
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
      i && n.initTree && n.initTree(i);
      const r = "uno-global-progress", e = `
          <div :class="Uno.getStyle('notify', 'progressWrap')"
               style="z-index: ${D("MODAL") + 2100}"
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
      `, u = I.ensureContainer(r, e, {
        "x-data": "{}"
      });
      u && n.initTree && n.initTree(u);
    },
    show(o, t = {}) {
      this.ensureContainer();
      const i = a(t), r = this.nextId++;
      let d = o;
      n.store("i18n") && typeof n.store("i18n").t == "function" && (d = n.store("i18n").t(o));
      const e = i.duration !== void 0 ? i.duration : 3e3;
      v.info("[Notify] Show:", { id: r, message: d, duration: e, config: i });
      const u = {
        id: r,
        message: d,
        isVisible: !0,
        // For x-show
        closable: !0,
        duration: e,
        content: i.hxUrl ? "" : void 0,
        // Initialize content for HTMX-loaded content
        ...i
      }, f = i.position || "top-right";
      return this.groups[f] ? this.groups[f].push(u) : (this.groups["top-right"] || (this.groups["top-right"] = []), this.groups["top-right"].push(u)), e > 0 && setTimeout(() => {
        this.remove(r);
      }, e), r;
    },
    remove(o) {
      Object.keys(this.groups).forEach((t) => {
        const i = this.groups[t].findIndex((r) => r.id === o);
        i !== -1 && (this.groups[t][i].visible = !1, setTimeout(() => {
          const r = this.groups[t].findIndex((d) => d.id === o);
          r !== -1 && this.groups[t].splice(r, 1);
        }, 300));
      });
    },
    // 3. Proxy Methods
    success(o, t = {}) {
      return this.show(o, { type: "success", ...t });
    },
    error(o, t = {}) {
      return this.show(o, { type: "error", ...t });
    },
    warning(o, t = {}) {
      return this.show(o, { type: "warning", ...t });
    },
    info(o, t = {}) {
      return this.show(o, { type: "info", ...t });
    }
    // 4. Loading Helpers - Removed to simplify the component
    // Loading indicator functionality has been removed to reduce complexity
  };
  H(n, "notify", l), l.init(), n.directive("uno-notify", (o, { expression: t }, { evaluate: i }) => {
    const r = () => {
      const d = n.store("uno")?.notify;
      if (!d) return;
      const e = {
        type: o.getAttribute("data-uno-type") || "info",
        position: o.getAttribute("data-uno-position"),
        duration: o.getAttribute("data-uno-duration") ? Number(o.getAttribute("data-uno-duration")) : void 0,
        closable: o.getAttribute("data-uno-closable") !== "false"
      };
      Object.keys(e).forEach(
        (c) => e[c] === void 0 && delete e[c]
      );
      let u = "", f = {};
      if (t) {
        const c = i(t);
        Array.isArray(c) ? (u = c[0], f = c[1] || {}) : typeof c == "object" && c !== null ? f = c : u = String(c);
      } else
        u = o.getAttribute("data-uno-message") || "";
      d.show(u, { ...e, ...f });
    };
    o.addEventListener("click", r);
  });
}, at = [
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
function lt(n, s = {}) {
  const a = (d) => {
    if (!(d instanceof HTMLElement)) return !1;
    const e = window.getComputedStyle(d);
    return e.display === "none" || e.visibility === "hidden" ? !1 : d.getClientRects().length > 0;
  }, l = () => Array.from(n.querySelectorAll(at)).filter((d) => a(d)), o = () => l()[0], t = () => {
    const d = l();
    return d[d.length - 1];
  }, i = (d) => {
    if (d.key === "Tab") {
      if (d.shiftKey) {
        if (document.activeElement === o()) {
          d.preventDefault();
          const e = t();
          e && e.focus();
        }
      } else if (document.activeElement === t()) {
        d.preventDefault();
        const e = o();
        e && e.focus();
      }
    }
  }, r = () => {
    if (n.contains(document.activeElement)) return;
    const d = o();
    d && d.focus();
  };
  return typeof requestAnimationFrame == "function" ? requestAnimationFrame(() => requestAnimationFrame(r)) : setTimeout(r, 0), n.addEventListener("keydown", i), () => {
    n.removeEventListener("keydown", i);
  };
}
const dt = (n, s = {}) => {
  const a = {
    defaultWidth: 400,
    defaultPosition: "center",
    okText: "OK",
    cancelText: "Cancel",
    closeOnEsc: !1,
    closeOnOverlayClick: !1
  };
  b.registerPluginDefaults("box", { ...a, ...s.plugins?.box, ...s.box });
  const l = (t) => b.getStyle("box", t);
  n.directive("uno-trap", (t, { expression: i }, { cleanup: r }) => {
    const d = lt(t);
    r(d);
  }), n.directive("uno-box", (t, { expression: i }, { evaluateLater: r, cleanup: d }) => {
    const u = i && i.trim().length > 0 ? r(i) : null, f = (x) => {
      if (x == null) return;
      const m = String(x).trim().toLowerCase();
      if (m === "true" || m === "1") return !0;
      if (m === "false" || m === "0") return !1;
    }, c = (x) => {
      if (x == null || x === "") return;
      const m = Number(x);
      return Number.isFinite(m) ? m : void 0;
    }, g = () => {
      const x = t.getAttribute("data-uno-box") || t.getAttribute("data-uno-box-type") || "dialog", m = t.getAttribute("data-uno-title") || "", w = t.getAttribute("data-uno-content") || "", y = t.getAttribute("data-uno-hx-url") || void 0, E = t.getAttribute("data-uno-position") || void 0, k = f(t.getAttribute("data-uno-modal")), L = f(t.getAttribute("data-uno-close-on-esc")), U = f(t.getAttribute("data-uno-close-on-overlay-click")), M = c(t.getAttribute("data-uno-width")), O = c(t.getAttribute("data-uno-height")), p = t.getAttribute("data-uno-width-class") || void 0, T = t.getAttribute("data-uno-height-class") || void 0;
      return {
        type: x,
        title: m,
        content: w,
        hxUrl: y,
        position: E,
        modal: k,
        closeOnEsc: L,
        closeOnOverlayClick: U,
        width: M,
        height: O,
        widthClass: p,
        heightClass: T
      };
    }, h = (x) => {
      x.preventDefault();
      const m = n.store("uno")?.box;
      if (!m) {
        console.warn("[Uno] Box store not found. Is Uno initialized?");
        return;
      }
      if (u) {
        u((w) => {
          if (w && typeof w == "object")
            m.open(w);
          else if (!w) {
            const y = g();
            (y.title || y.content || y.hxUrl) && m.open(y);
          }
        });
        return;
      }
      m.open(g());
    };
    t.addEventListener("click", h), d(() => t.removeEventListener("click", h));
  }), H(n, "box", {
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
    alert(t, i = {}) {
      const r = typeof i == "string" ? { title: i } : i;
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
    confirm(t, i = {}) {
      return new Promise((r) => {
        this.open("confirm", {
          content: t,
          onOk: () => r(!0),
          onCancel: () => r(!1),
          onClose: () => r(!1),
          // Closing without choice = cancel
          ...i
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
    prompt(t, i = "", r = {}) {
      return new Promise((d) => {
        const e = "prompt-" + Math.random().toString(36).slice(2, 11), u = l("promptWrapper"), f = l("promptText"), c = l("promptInput"), g = b.dom.escape(t), h = b.dom.escape(i), x = `
          <div class="${u}">
            <p class="${f}">${g}</p>
            <input id="${e}" type="text" class="${c}" value="${h}">
          </div>
        `;
        this.open("confirm", {
          title: r.title || "Prompt",
          content: x,
          _isTrustedPrompt: !0,
          // Internal flag for prompt content
          onOk: () => {
            const m = document.getElementById(e);
            d(m ? m.value : null);
          },
          onCancel: () => d(null),
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
    open(t, i = {}) {
      let r = "dialog", d = i;
      typeof t == "object" ? (d = t, r = d.type || "dialog") : r = t;
      const e = b.getConfig("box", d), u = e.boxId || Math.random().toString(36).slice(2, 10);
      let f = e.showOkButton, c = e.showCancelButton;
      r === "confirm" ? (f = f !== !1, c = c !== !1) : r === "dialog" && (f = f || !1, c = c || !1);
      const g = {
        id: u,
        type: r,
        title: e.title || "",
        content: e._isTrustedPrompt ? e.content || "" : b.dom.sanitize(e.content || ""),
        // Security: Always sanitize content (trustedHTML removed)
        // Exception: Internal prompt content is pre-escaped
        _isTrustedPrompt: e._isTrustedPrompt || !1,
        width: e.width !== void 0 ? e.width : e.widthClass ? void 0 : e.defaultWidth || 400,
        height: e.height,
        widthClass: e.widthClass,
        heightClass: e.heightClass,
        containerClass: e.containerClass,
        position: e.position || e.defaultPosition || "center",
        modal: e.modal !== void 0 ? e.modal : !0,
        closeOnEsc: e.closeOnEsc !== void 0 ? e.closeOnEsc : !1,
        closeOnOverlayClick: e.closeOnOverlayClick !== void 0 ? e.closeOnOverlayClick : !1,
        showCloseButton: e.showCloseButton !== void 0 ? e.showCloseButton : !0,
        showOkButton: f,
        showCancelButton: c,
        okText: e.okText || "OK",
        cancelText: e.cancelText || "Cancel",
        onOk: e.onOk,
        onCancel: e.onCancel,
        onBeforeClose: e.onBeforeClose,
        onClose: e.onClose,
        onContentLoad: e.onContentLoad,
        hxUrl: e.hxUrl,
        hxMethod: e.hxMethod || "GET",
        hxHeaders: e.hxHeaders,
        hxTrigger: e.hxTrigger || "load",
        // Internal state - start closed to trigger animation
        isOpen: !1,
        /**
         * Close this box dialog
         */
        close() {
          this.onBeforeClose && this.onBeforeClose(), this.isOpen = !1;
          const m = typeof document < "u" ? document.getElementById("uno-box-" + this.id) : null, w = () => {
            const k = n.store("uno")?.box;
            k && k.remove(this.id), this.onClose && this.onClose();
          };
          if (!m) {
            w();
            return;
          }
          let y = !1;
          const E = (k) => {
            k && k.target !== m || y || (y = !0, m.removeEventListener("transitionend", E), w());
          };
          m.addEventListener("transitionend", E), setTimeout(E, 350);
        }
      };
      this.items.push(g), this.ensureContainer();
      const h = this.items[this.items.length - 1], x = () => {
        h.isOpen = !0;
      };
      return n && typeof n.nextTick == "function" ? n.nextTick(x) : setTimeout(x, 0), h;
    },
    /**
     * Close a specific box by ID
     * @param {string} id - Box ID
     */
    close(t) {
      const i = this.items.find((r) => r.id === t);
      i && i.close();
    },
    /**
     * Remove a box from the items array
     * @param {string} id - Box ID
     */
    remove(t) {
      const i = this.items.findIndex((r) => r.id === t);
      i !== -1 && this.items.splice(i, 1);
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
      return et(t);
    },
    /**
     * Handle HTMX content load completion
     * @param {HTMLElement} el - Content element
     * @param {object} item - Box item
     */
    onHtmxContentLoaded(t, i) {
      i.onContentLoad && typeof i.onContentLoad == "function" && i.onContentLoad(t), window.Alpine && window.Alpine.initTree && window.Alpine.initTree(t);
    },
    /**
     * Ensure the box container exists in the DOM
     */
    ensureContainer() {
      const t = "uno-box-manager", i = typeof document < "u" && document.getElementById(t);
      if (i)
        return n && typeof n.initTree == "function" && !i._x_dataStack && n.initTree(i), i;
      const d = {
        "x-data": "{}",
        class: "relative",
        style: `z-index: ${b.options?.ui?.layers?.boxBase || 1e4};`,
        "x-show": "$store.uno.box.items.length > 0",
        "x-cloak": ""
        // Hide until Alpine is ready
      }, u = b.dom.ensureContainer(t, `
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
                 :style="'z-index:' + (Uno.options?.ui?.layers?.boxBase || 10000)"
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
      `, d);
      return u && n && typeof n.initTree == "function" && n.initTree(u), u;
    }
  });
}, V = {
  // Common Regex Patterns
  email: /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/,
  url: /^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([/\w .-]*)*\/?$/,
  username: /^[a-zA-Z0-9_-]{3,16}$/,
  number: /^\d+$/,
  integer: /^-?\d+$/,
  date: /^\d{4}-\d{2}-\d{2}$/,
  time: /^\d{2}:\d{2}$/,
  phone: /^\+?[\d\s-]{10,}$/,
  alphanumeric: /^[a-zA-Z0-9]+$/
}, ct = (n) => {
  n.store("unoValid", {
    errors: {},
    setError(t, i) {
      this.errors[t] = i;
    },
    clearError(t) {
      delete this.errors[t];
    },
    getMessage(t) {
      return this.errors[t] || "";
    }
  }), typeof window < "u" && (window.$uno || (window.$uno = {}), window.$uno.valid = {
    getMessage: (t) => n.store("unoValid").getMessage(t)
  }), typeof document < "u" && (document.body.addEventListener("htmx:confirm", (t) => {
    const i = t.detail.elt, r = i.tagName === "FORM" ? i : i.closest("form");
    r && r.hasAttribute("data-uno-valid") && (r.checkValidity() || (t.preventDefault(), t.stopPropagation(), a(r), r.dispatchEvent(new CustomEvent("uno:validation-failed", { bubbles: !0 }))));
  }), document.body.addEventListener("input", (t) => {
    const i = t.target;
    i.form && i.form.hasAttribute("data-uno-valid") && (i.getAttribute("aria-invalid") === "true" && l(i), i.name && i.form.querySelectorAll(`[data-uno-match="${i.name}"]`).forEach((d) => {
      (d.getAttribute("aria-invalid") === "true" || d.value) && l(d);
    }));
  }), document.body.addEventListener(
    "blur",
    (t) => {
      const i = t.target;
      (i.tagName === "INPUT" || i.tagName === "TEXTAREA" || i.tagName === "SELECT") && i.form && i.form.hasAttribute("data-uno-valid") && l(i);
    },
    !0
  ), document.body.addEventListener("input", (t) => {
    const i = t.target;
    i.matches && i.matches("input, textarea") && s(i);
  }));
  function s(t) {
    let i = null;
    if (t.hasAttribute("data-uno-count")) {
      const r = t.getAttribute("data-uno-count");
      r && (i = document.querySelector(r));
    }
    if (!i && t.form) {
      const r = t.closest("div");
      r && (i = r.querySelector(".valid-count span") || r.querySelector(".valid-count"));
    }
    if (i) {
      const r = t.value.length;
      i.textContent = r;
    }
  }
  function a(t) {
    const i = t.querySelectorAll("input, select, textarea");
    let r = null;
    return i.forEach((d) => {
      l(d) || r || (r = d);
    }), r && r.focus(), !r;
  }
  function l(t) {
    if (t.disabled || t.type === "hidden") return !0;
    t.setCustomValidity("");
    let i = [];
    const r = t.getAttribute("data-uno-rules") || t.getAttribute("valid");
    if (r)
      try {
        const c = JSON.parse(r);
        Array.isArray(c) && (i = c);
      } catch (c) {
        console.warn("Uno Valid: Invalid JSON rules", c);
      }
    const d = t.getAttribute("data-uno-validator") || t.getAttribute("data-valid-type");
    d && i.push({ type: d }), t.hasAttribute("required") && !i.some((c) => c.type === "required") && i.unshift({ type: "required" });
    for (const c of i) {
      if (c.type === "emptyquit") {
        if (!t.value) {
          t.setCustomValidity("");
          break;
        }
        continue;
      }
      if (c.type === "required") {
        if (t.validity.valueMissing) {
          t.setCustomValidity(
            c.msg || o("valid.required") || "This field is required"
          );
          break;
        }
        continue;
      }
      if (V[c.type]) {
        if (t.value && !V[c.type].test(t.value)) {
          t.setCustomValidity(c.msg || o(`valid.${c.type}`) || "Invalid format");
          break;
        }
        continue;
      }
      if (typeof window < "u" && window.UnoValidators && window.UnoValidators[c.type]) {
        const g = window.UnoValidators[c.type](t.value, t);
        if (g !== !0 && g !== void 0 && g !== "") {
          const h = typeof g == "string" ? g : o(`valid.${c.type}`) || "Validation failed";
          t.setCustomValidity(c.msg || h);
          break;
        }
        continue;
      }
      if (typeof window < "u" && window[c.type] && typeof window[c.type] == "function") {
        const g = window[c.type](t.value, t);
        if (g !== !0 && g !== void 0 && g !== "") {
          const h = typeof g == "string" ? g : "Validation failed";
          t.setCustomValidity(c.msg || h);
          break;
        }
        continue;
      }
    }
    if (t.hasAttribute("data-uno-match")) {
      const c = t.getAttribute("data-uno-match"), g = t.form.querySelector(`[name="${c}"]`);
      g && t.value !== g.value && t.setCustomValidity(t.getAttribute("data-uno-error-match") || "Fields do not match");
    }
    const e = t.checkValidity(), u = b.getStyle("valid", "errorClass", "is-invalid");
    e ? t.classList.remove(u) : t.classList.add(u);
    const f = b.getStyle("valid", "parentErrorClass");
    if (f) {
      const g = (typeof window < "u" && window.Uno?.options?.valid || {}).parentSelector || ".field, .form-group, .input-group, div", h = t.closest(g);
      h && (e ? h.classList.remove(f) : h.classList.add(f));
    }
    if (t.setAttribute("aria-invalid", !e), e)
      t.name && n.store("unoValid").clearError(t.name);
    else {
      let c = t.validationMessage;
      if (t.validity.valueMissing)
        c = t.getAttribute("data-uno-error-required") || o("valid.required") || c;
      else if (t.validity.typeMismatch) {
        let g = "valid.pattern";
        t.type === "email" && (g = "valid.email"), t.type === "url" && (g = "valid.url"), c = t.getAttribute("data-uno-error-type") || o(g) || c;
      } else t.validity.patternMismatch ? c = t.getAttribute("data-uno-error-pattern") || o("valid.pattern") || c : t.validity.tooShort ? c = t.getAttribute("data-uno-error-minlength") || o("valid.min") || c : t.validity.tooLong ? c = t.getAttribute("data-uno-error-maxlength") || o("valid.max") || c : t.validity.rangeUnderflow ? c = t.getAttribute("data-uno-error-min") || o("valid.rangeMin") || c : t.validity.rangeOverflow && (c = t.getAttribute("data-uno-error-max") || o("valid.rangeMax") || c);
      t.name && n.store("unoValid").setError(t.name, c);
    }
    return e;
  }
  function o(t) {
    return typeof window < "u" && window.Alpine && window.Alpine.store("i18n") ? window.Alpine.store("i18n").t(t) : null;
  }
};
const ut = (n, s = {}) => {
  const a = {
    defaultTab: "",
    hash: !0,
    activeClass: "is-active",
    inactiveClass: ""
  };
  b.registerPluginDefaults("tab", { ...a, ...s.plugins?.tab, ...s.tab }), H(n, "tab", {
    // State: { [group]: activeTabName }
    groups: {},
    // Registry: { [tabName]: { group, name, ... } }
    registry: {},
    init() {
      const o = () => {
        let t = window.location.hash.slice(1);
        if (!t && window.location.href.includes("#") && (t = window.location.href.split("#")[1]), v.info(`[Tab] checkHash: ${t}`), t && this.registry[t]) {
          const { group: i, name: r } = this.registry[t];
          v.info(`[Tab] checkHash match: ${r}`), this.activate({ group: i, name: r, updateHash: !1 });
        }
      };
      o(), window.addEventListener("hashchange", o), window.addEventListener("load", o);
    },
    /**
     * Activate a tab
     * @param {object} instanceOpts - Options { name, group, ... }
     */
    activate(o = {}) {
      const t = (typeof window < "u" ? window.Uno?.options : b.options) || {}, i = a, d = {
        group: "default",
        hash: !0,
        ...t.plugins?.tab || t.tab || i,
        ...o
      };
      d.name && (v.info(
        `[Tab] Activating ${d.name} in ${d.group}. Hash update: ${d.updateHash !== !1}`
      ), this.registry[d.name] || (this.registry[d.name] = { group: d.group, name: d.name }), this.groups[d.group] = d.name, d.hash && d.updateHash !== !1 && window.location.hash.slice(1) !== d.name && window.history.replaceState(null, null, `#${d.name}`), d.updateHash === !1 && setTimeout(() => {
        const e = document.querySelector(`[data-uno-name="${d.name}"]`) || document.getElementById(`tab-${d.name}`);
        e && e.hasAttribute("hx-get") && (v.info(`[Tab] Triggering HTMX for ${d.name}`), window.htmx && window.htmx.trigger(e, "click"));
      }, 0));
    },
    /**
     * Check if a tab is active
     * @param {string} name - Tab name
     * @param {string} group - Group name
     * @returns {boolean}
     */
    isActive(o, t = "default") {
      return this.groups[t] === o;
    },
    /**
     * Register a tab (helper for directives or manual setup)
     * @param {string} name
     * @param {string} group
     */
    register(o, t = "default") {
      v.info(`[Tab] Registering ${o} in ${t}. Hash: ${window.location.hash}`), this.registry[o] = { name: o, group: t }, this.groups[t] || (this.groups[t] = o);
      let i = window.location.hash.slice(1);
      !i && window.location.href.includes("#") && (i = window.location.href.split("#")[1]), i === o && (v.info(`[Tab] Hash match for ${o}. Activating.`), this.activate({ group: t, name: o, updateHash: !1 }));
    }
  }), n.store("uno").tab.init(), n.directive("uno-tab", (o, { expression: t }, { evaluateLater: i, cleanup: r }) => {
    const d = t ? i(t) : () => {
    }, e = () => ({
      name: o.getAttribute("data-uno-name") || o.getAttribute("data-tab") || o.id?.replace("tab-", ""),
      group: o.getAttribute("data-uno-group") || "default",
      url: o.getAttribute("data-uno-url"),
      // For HTMX reference if needed
      target: o.getAttribute("data-uno-target")
    });
    let f = { ...e() };
    t ? d((g) => {
      typeof g == "object" ? f = { ...f, ...g } : typeof g == "string" && (f.name = g), n.store("uno").tab.register(f.name, f.group);
    }) : n.store("uno").tab.register(f.name, f.group);
    const c = (g) => {
      g.preventDefault(), t ? d((h) => {
        const x = typeof h == "object" ? h : { name: h };
        n.store("uno").tab.activate({ ...e(), ...x });
      }) : n.store("uno").tab.activate(e());
    };
    o.addEventListener("click", c), r(() => o.removeEventListener("click", c));
  });
};
b.dom = I;
typeof window < "u" && (window.Uno = b);
b.use(tt).use(st).use(rt).use(dt).use(ct).use(ut);
if (typeof window < "u" && typeof document < "u") {
  const n = () => {
    if (b.initialized) return;
    const s = window.UnoConfig || {};
    b.init(window.Alpine, s);
  };
  window.Alpine ? n() : document.addEventListener("alpine:init", n);
}
export {
  b as default
};
