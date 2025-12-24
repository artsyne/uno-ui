import Y from "alpinejs";
let j = !1, P = null;
const $ = {
  /**
   * Set debug mode
   * @param {boolean} debug
   */
  setDebug(n) {
    j = n;
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
    j && typeof console < "u" && typeof console.log == "function" && console.log("[Uno Info]", ...n);
  },
  /**
   * Log warning message
   * @param {...any} args - Arguments to log
   */
  warn(...n) {
    j && typeof console < "u" && typeof console.warn == "function" && console.warn("[Uno Warn]", ...n);
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
        j && typeof console < "u" && typeof console.error == "function" && console.error("[Uno Error] Reporter failed:", s);
      }
    j && typeof console < "u" && typeof console.error == "function" && console.error("[Uno Error]", ...n);
  }
};
function et(n) {
  return String(n ?? "").replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;").replace(/'/g, "&#39;");
}
const X = /* @__PURE__ */ new WeakMap();
function N(n, s) {
  if (!n) return !1;
  const l = String(s);
  let i = X.get(n);
  return i || (i = /* @__PURE__ */ new Set(), X.set(n, i)), i.has(l) ? !1 : (i.add(l), !0);
}
function _(n) {
  return n ? n.replace(/<script\b[^>]*>([\s\S]*?)<\/script>/gim, "").replace(/<img\b[^>]*>/gim, "").replace(/\s+on[a-z]+="[^"]*"/gim, "").replace(/\s+on[a-z]+='[^']*'/gim, "").replace(/javascript:/gim, "") : "";
}
function J(n, s, { trusted: l = !1, alpine: i } = {}) {
  n && (n.innerHTML = l ? String(s ?? "") : _(String(s ?? "")), i?.initTree?.(n));
}
const G = {
  /**
   * Ensure a container element exists in the DOM.
   * @param {string} id - Element ID
   * @param {string} html - Inner HTML content
   * @param {Object} attributes - Attributes to set
   * @param {HTMLElement} parent - Parent element (default: document.body)
   */
  ensureContainer(n, s = "", l = {}, i = document.body) {
    if (!document.getElementById(n))
      try {
        const d = document.createElement("div");
        return d.id = n, window.Uno?.options?.csp?.nonce && !l.nonce && (l.nonce = window.Uno.options.csp.nonce), Object.entries(l).forEach(([t, r]) => {
          d.setAttribute(t, r);
        }), s && (d.innerHTML = s), i ? i.appendChild(d) : document.documentElement.appendChild(d), d;
      } catch (d) {
        $.error(`DomManager: Failed to inject container #${n}`, d);
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
        const l = document.createElement("style");
        l.id = n, window.Uno?.options?.csp?.nonce && l.setAttribute("nonce", window.Uno.options.csp.nonce), l.textContent = s, document.head.appendChild(l);
      } catch (l) {
        $.error(`DomManager: Failed to inject styles #${n}`, l);
      }
  },
  // Exports
  setHTML: J,
  // Secure HTML setting
  escape: et,
  // HTML escaping
  sanitize: _
  // HTML sanitization
};
function A(n) {
  return Object.prototype.toString.call(n) === "[object Object]";
}
function W(n, s) {
  if (!A(n)) return A(s) ? { ...s } : s;
  if (!A(s)) return s === void 0 ? { ...n } : s;
  const l = { ...n };
  for (const [i, d] of Object.entries(s)) {
    if (A(d) && A(l[i])) {
      l[i] = W(l[i], d);
      continue;
    }
    if (A(d)) {
      l[i] = W({}, d);
      continue;
    }
    if (typeof d == "string" && d.trim().indexOf("+") === 0 && typeof l[i] == "string") {
      l[i] = `${l[i]} ${d.trim().slice(1)}`.trim();
      continue;
    }
    l[i] = d;
  }
  return l;
}
class nt {
  constructor() {
    this.version = "0.5.0", this.plugins = [], this.initialized = !1, this._htmxRequestHandlers = [], this._htmxActiveXhrs = typeof WeakSet == "function" ? /* @__PURE__ */ new WeakSet() : null, this.htmx = {
      available: () => typeof window < "u" && typeof window.htmx?.ajax == "function",
      ajax: (l, i, d = {}) => typeof window > "u" || typeof window.htmx?.ajax != "function" ? Promise.reject(new Error("htmx.ajax is not available")) : window.htmx.ajax(l, i, d),
      fragment: async (l, i, d = {}) => {
        if (typeof document > "u") throw new Error("document is not available");
        const t = document.createElement("div");
        let r = "";
        const a = (o) => {
          o.preventDefault(), r = o.detail.serverResponse, o.detail.shouldSwap = !1;
        };
        t.addEventListener("htmx:beforeSwap", a);
        try {
          const o = d.swap || "innerHTML";
          return await this.htmx.ajax(l, i, { ...d, target: t, swap: o }), r;
        } finally {
          t.removeEventListener("htmx:beforeSwap", a), t.remove();
        }
      },
      onRequest: (l = {}) => {
        if (!l || typeof l != "object" || (this._htmxRequestHandlers.push(l), typeof window > "u" || typeof document > "u") || !window.htmx || !N(document, "uno:htmx:request-hooks")) return;
        const i = (a) => a?.detail?.xhr, d = (a, o) => {
          for (const m of this._htmxRequestHandlers) {
            const v = m?.[a];
            if (typeof v == "function")
              try {
                v(o);
              } catch (e) {
                $.error("Uno.htmx.onRequest handler failed", e);
              }
          }
        }, t = (a) => {
          const o = i(a);
          if (o && this._htmxActiveXhrs) {
            if (this._htmxActiveXhrs.has(o)) return;
            this._htmxActiveXhrs.add(o);
          }
          d("start", a);
        }, r = (a) => {
          const o = i(a);
          if (o && this._htmxActiveXhrs) {
            if (!this._htmxActiveXhrs.has(o)) return;
            this._htmxActiveXhrs.delete(o);
          }
          d("end", a);
        };
        document.addEventListener("htmx:beforeRequest", t), document.addEventListener("htmx:afterRequest", r), document.addEventListener("htmx:responseError", r), document.addEventListener("htmx:sendError", r), document.addEventListener("htmx:timeout", r), document.addEventListener("htmx:xhr:progress", (a) => d("progress", a));
      }
    };
    const s = "rounded-lg shadow-md border border-gray-200 dark:border-gray-700";
    this.options = {
      debug: !1,
      csp: {
        nonce: null
        // For Content Security Policy
      },
      onError: null,
      // Custom error handler
      plugins: {},
      ui: {
        base: s,
        layers: { boxBase: 1e4, tooltipBase: 2e4, notifyBase: 3e4 },
        box: {
          overlay: "fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity",
          dialogWrapper: "fixed inset-0 z-10 w-screen overflow-y-auto",
          itemWrapper: "flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0",
          container: "relative transform overflow-hidden {base} bg-white dark:bg-gray-900 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg flex flex-col max-h-[90vh] pointer-events-auto outline-none",
          cancelBtn: "mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:mt-0 sm:w-auto",
          okBtn: "inline-flex w-full justify-center rounded-md bg-blue-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 sm:ml-3 sm:w-auto",
          header: "flex justify-between items-center p-4 border-b border-slate-200 dark:border-slate-700",
          title: "text-base font-semibold leading-6 text-gray-900 dark:text-slate-200",
          closeBtn: "text-slate-400 hover:text-slate-500 focus:outline-none",
          closeBtnFloating: "absolute top-2 right-2 z-10 text-slate-400 hover:text-slate-500 focus:outline-none",
          closeIcon: "h-6 w-6",
          closeIconSmall: "h-5 w-5",
          content: "px-4 pb-4 pt-5 sm:p-6 sm:pb-4 overflow-auto flex-1 text-slate-600 dark:text-slate-300",
          loadingWrap: "flex justify-center items-center h-full",
          spinner: "animate-spin h-8 w-8 text-blue-500",
          footer: "bg-gray-50 px-4 py-3 sm:flex sm:flex-row-reverse sm:px-6 border-t border-slate-200 dark:border-slate-700 dark:bg-slate-800",
          error: "p-4 text-red-500",
          promptWrapper: "space-y-3",
          promptText: "text-slate-700 dark:text-slate-200",
          promptInput: "w-full px-3 py-2 border rounded focus:ring-2 focus:ring-blue-500 outline-none",
          positionBase: "fixed inset-0 flex transition-opacity duration-300 p-4 pointer-events-none"
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
          indicatorTrack: "h-1 bg-gray-200 dark:bg-slate-700 overflow-hidden",
          indicatorBar: "h-full bg-blue-500 origin-left",
          progressWrap: "fixed top-0 left-0 right-0 h-1 bg-transparent",
          progressBar: "uno-progress-anim h-full bg-blue-500 animate-pulse",
          loadingWrap: "fixed inset-0 flex items-center justify-center bg-white/50 dark:bg-black/50 backdrop-blur-sm transition-opacity duration-300 pointer-events-auto",
          spinner: "animate-spin h-10 w-10 text-blue-600 dark:text-blue-400",
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
          content: "fixed px-3 py-2 text-sm font-medium {base} whitespace-nowrap pointer-events-none bg-white text-gray-900 dark:bg-slate-800 dark:text-white",
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
  }
  _getUiValue(s, l) {
    const i = this.options.ui?.[s];
    if (!i) return;
    if (typeof l != "string" || l.indexOf(".") === -1) return i?.[l];
    const d = l.split(".").filter(Boolean);
    let t = i;
    for (const r of d) {
      if (t == null) return;
      t = t[r];
    }
    return t;
  }
  _resolveUiTokens(s) {
    if (typeof s != "string") return s;
    const l = this.options.ui?.base;
    return typeof l == "string" && s.indexOf("{base}") !== -1 ? s.replaceAll("{base}", l) : s;
  }
  /**
   * Get style class from configuration or fallback
   * Support "+" prefix to append class to default
   * @param {string} comp - Component name
   * @param {string} elem - Element name
   * @param {string} defaultCls - Default class string
   * @returns {string} Final class string
   */
  getStyle(s, l, i) {
    const d = this._getUiValue(s, l);
    return d === void 0 ? this._resolveUiTokens(i) : typeof d == "string" && d.trim().indexOf("+") === 0 ? this._resolveUiTokens(`${i} ${d.trim().slice(1)}`.trim()) : this._resolveUiTokens(d);
  }
  /**
   * Set global configuration (before or after init)
   * @param {Object} config - Global configuration object
   */
  setConfig(s) {
    return this.options = W(this.options, s), this.initialized && $.info("Uno configuration updated", this.options), this;
  }
  /**
   * Register a plugin
   * @param {Function} plugin - Plugin function (Alpine, options) => void
   * @param {Object} options - Plugin options
   * @returns {UnoInstance}
   */
  use(s, l = {}) {
    return this.initialized && $.warn("Uno.use() called after initialization. Plugin might not work correctly."), this.plugins.push({ plugin: s, options: l }), this;
  }
  /**
   * Initialize Uno
   * @param {Object|Function} alpineOrOptions - Alpine instance or Configuration options
   * @param {Object} options - Configuration options (if first arg is Alpine)
   * @returns {Object} Alpine instance
   */
  init(s = {}, l = {}) {
    if (this.initialized)
      return $.warn("Uno has already been initialized."), Y;
    let i = Y, d = l;
    return !i && typeof window < "u" && window.Alpine && (i = window.Alpine), s && (s.version || typeof s.start == "function" || typeof s.store == "function") ? i = s : d = s || l, this.options = W(this.options, d), $.setDebug(this.options.debug !== !1), this.options.onError && $.setReporter(this.options.onError), $.info(`Uno v${this.version} initializing...`), typeof window < "u" && typeof document < "u" && window.htmx && typeof i?.initTree == "function" && N(document, "uno:htmx-afterSwap:initTree") && document.addEventListener("htmx:afterSwap", (t) => {
      const r = t?.detail?.elt;
      r && i.initTree(r);
    }), this.plugins.forEach(({ plugin: t, options: r }) => {
      try {
        const a = W(this.options, r);
        t(i, a);
      } catch (a) {
        $.error("Failed to load plugin", a);
      }
    }), this.initialized = !0, $.info("Uno initialized"), i;
  }
}
const g = new nt(), ot = {
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
}, it = {
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
}, rt = {
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
const Q = {
  en: ot,
  "zh-CN": it,
  "zh-TW": rt
}, Z = Object.keys(Q);
function st() {
  if (typeof localStorage < "u") {
    const n = localStorage.getItem("uno.lang");
    if (n && Z.includes(n)) return n;
  }
  if (typeof navigator < "u") {
    const n = navigator.language, s = n.startsWith("zh") ? n.includes("TW") || n.includes("Hant") ? "zh-TW" : "zh-CN" : "en";
    return Z.includes(s) ? s : "en";
  }
  return "en";
}
function at(n, s) {
  const l = s.split(".");
  let i = n;
  for (const d of l) {
    if (!i || typeof i != "object" || !(d in i)) return;
    i = i[d];
  }
  return i;
}
function tt(n, s) {
  for (const l in s)
    s[l] && typeof s[l] == "object" ? (n[l] || Object.assign(n, { [l]: {} }), tt(n[l], s[l])) : Object.assign(n, { [l]: s[l] });
}
function lt(n) {
  n.store("i18n", {
    locale: st(),
    messages: Q,
    init() {
      typeof document < "u" && (document.documentElement.lang = this.locale);
    },
    t(s, l = {}) {
      const i = this.locale, d = this.messages[i] || this.messages.en, t = at(d, s);
      return (typeof t == "string" ? t : s).replace(/\{(\w+)\}/g, (a, o) => o in l ? String(l[o]) : `{${o}}`);
    },
    setLang(s) {
      if (!this.messages[s]) {
        $.warn(`[Uno] Language ${s} not found`);
        return;
      }
      this.locale = s, typeof localStorage < "u" && localStorage.setItem("uno.lang", s), typeof document < "u" && (document.documentElement.lang = s);
    },
    extend(s) {
      for (const l in s)
        this.messages[l] ? tt(this.messages[l], s[l]) : this.messages[l] = s[l];
    }
  });
}
const dt = [
  "top-left",
  "top-center",
  "top-right",
  "center-left",
  "center",
  "center-right",
  "bottom-left",
  "bottom-center",
  "bottom-right"
], ct = (n) => {
  const s = n || "top-right", l = {
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
  return l[s] || l["top-right"];
}, ut = (n) => {
  const s = n || "top-right";
  let l = "flex-col", i = "items-end", d = "space-y-2";
  return s.includes("bottom") && (l = "flex-col-reverse", d = "space-y-reverse space-y-2"), s.includes("left") ? i = "items-start" : s.includes("center") ? s === "center" ? i = "items-center justify-center" : i = "items-center" : i = "items-end", { flexDir: l, alignItems: i, spaceY: d };
};
function ft(n, s, l) {
  const i = n.getBoundingClientRect(), d = s.getBoundingClientRect(), t = d.width, r = d.height, a = i.width, o = i.height, m = { top: i.top, left: i.left }, v = window.innerWidth, e = window.innerHeight, c = 0, f = 0, h = e, w = v, u = (C, b) => {
    if (b) {
      if (C + t > w)
        return w - t - 10;
      if (C < f)
        return f + 10;
    } else {
      if (C + r > h)
        return h - r - 10;
      if (C < c)
        return c + 10;
    }
    return C;
  };
  let y, p, E = l;
  switch (l) {
    case "left":
      p = m.left - t - 5, p < f && (p = m.left + a + 5, E = "right"), y = u(m.top + o / 2 - r / 2, !1);
      break;
    case "top":
      y = m.top - r - 8, y < c && (y = m.top + o + 8, E = "bottom"), p = m.left + (a - t) / 2, p = u(p, !0);
      break;
    case "right":
      p = m.left + a + 5, p + t > w && (p = m.left - t - 5, E = "left"), y = u(m.top + o / 2 - r / 2, !1);
      break;
    case "bottom":
    default:
      y = m.top + o + 8, y + r > h && (y = m.top - r - 8, E = "top"), p = m.left + (a - t) / 2, p = u(p, !0);
      break;
  }
  return { top: y, left: p, position: E };
}
const mt = (n, s = {}) => {
  const l = s.plugins?.tooltip || s.tooltip || {}, i = () => {
    const u = typeof window < "u" ? window.Uno?.options : g.options;
    return u?.plugins?.tooltip || u?.tooltip || l;
  };
  if (!n) {
    $.error("Tooltip: Alpine is undefined. Plugin initialization failed.");
    return;
  }
  let d = null, t = null, r = null, a = 0, o = null;
  const m = () => {
    let u = Number(g.options?.ui?.layers?.tooltipBase ?? 2e4);
    return Number.isFinite(u) || (u = 2e4), u;
  }, v = () => d && t && r && typeof document < "u" && document.contains(d) ? !0 : (d = typeof document < "u" ? document.getElementById("uno-tooltip-container") : null, t = typeof document < "u" ? document.getElementById("uno-tooltip-content") : null, r = typeof document < "u" ? document.getElementById("uno-tooltip-body") : null, !!(d && t && r)), e = () => {
    a || typeof requestAnimationFrame == "function" && (a = requestAnimationFrame(() => {
      a = 0, n.store("tooltip")?.updatePosition?.();
    }));
  }, c = () => {
    if (o || typeof window > "u" || typeof document > "u") return;
    const u = () => e(), y = () => e(), p = (C) => {
      if (C.key !== "Escape") return;
      const b = n.store("tooltip");
      b?.visible && b.hide();
    }, E = (C) => {
      const b = n.store("tooltip");
      if (!b?.visible || b.triggerMode !== "click" || b.clickOutside !== !0 || !v()) return;
      const M = C.target, x = b.eventEl && typeof b.eventEl.contains == "function" && b.eventEl.contains(M), U = t && t.contains(M);
      !x && !U && b.hide();
    };
    window.addEventListener("scroll", u, !0), window.addEventListener("resize", y), window.addEventListener("keydown", p), window.addEventListener("click", E), o = () => {
      window.removeEventListener("scroll", u, !0), window.removeEventListener("resize", y), window.removeEventListener("keydown", p), window.removeEventListener("click", E), o = null;
    };
  };
  n.store("tooltip", {
    visible: !1,
    content: "",
    placement: "top",
    top: -9999,
    // Start off-screen
    left: -9999,
    zIndex: m(),
    triggerEl: null,
    eventEl: null,
    triggerMode: "hover",
    clickOutside: !0,
    onShow: null,
    async show(u, y, p = {}) {
      this.triggerEl = u, this.eventEl = p.eventEl || u, this.content = y, this.placement = p.placement || i().placement || "top", this.onShow = p.onShow, this.triggerMode = p.trigger || "hover", this.clickOutside = p.clickOutside !== void 0 ? p.clickOutside : this.triggerMode === "click", this.zIndex = m(), this.visible = !0, this.ensureContainer(), c(), v() && J(r, this.content, { trusted: p.trustedHTML === !0, alpine: n }), await n.nextTick(), this.updatePosition(), this.onShow && this.onShow(r);
    },
    hide() {
      this.visible = !1, this.triggerEl = null, this.eventEl = null, o && o(), a && (cancelAnimationFrame(a), a = 0);
    },
    updatePosition() {
      if (!this.triggerEl || !v()) return;
      const u = ft(this.triggerEl, t, this.placement);
      this.top = u.top, this.left = u.left;
    },
    ensureContainer() {
      f(), v();
    },
    destroy() {
      this.hide();
    }
  });
  const f = () => {
    if (typeof document > "u") return;
    const u = g.getStyle("tooltip", "content"), y = g.getStyle("tooltip", "body"), p = `
          <div id="uno-tooltip-content"
               x-show="$store.tooltip.visible"
               role="tooltip"
               x-transition:enter="transition ease-out duration-200"
               x-transition:enter-start="opacity-0 scale-95"
               x-transition:enter-end="opacity-100 scale-100"
               x-transition:leave="transition ease-in duration-150"
               x-transition:leave-start="opacity-100 scale-100"
               x-transition:leave-end="opacity-0 scale-95"
               class="${u}"
              :style="\`top: \${$store.tooltip.top}px; left: \${$store.tooltip.left}px; z-index: \${$store.tooltip.zIndex};\`"
              style="display: none;">
              <div id="uno-tooltip-body" class="${y}"></div>
          </div>
        `;
    G.ensureContainer("uno-tooltip-container", p, {
      "x-data": "{}"
    });
  }, h = (u, { expression: y, modifiers: p }, { evaluateLater: E, effect: C, cleanup: b }) => {
    const M = E(y);
    let x = {
      content: "",
      placement: null,
      trigger: null,
      clickOutside: void 0
    };
    const U = () => {
      const k = u.getAttribute("data-uno-placement");
      return k || (x.placement ? x.placement : p.includes("top") ? "top" : p.includes("bottom") ? "bottom" : p.includes("left") ? "left" : p.includes("right") ? "right" : i().placement || "top");
    }, B = () => {
      const k = u.getAttribute("data-uno-trigger");
      return k || (x.trigger ? x.trigger : p.includes("click") ? "click" : p.includes("focus") ? "focus" : p.includes("hover") ? "hover" : i().trigger || "hover");
    }, R = () => (typeof u.disabled == "boolean" && u.disabled === !0 || u.getAttribute("aria-disabled") === "true" || u.getAttribute("disabled") !== null) && u.parentElement || u, D = (k) => {
      if (typeof k == "string" ? x.content = k : typeof k == "object" && k !== null && (x = { ...x, ...k }), x.content && x.content.startsWith("#") && !x.content.includes(" ")) {
        const L = document.querySelector(x.content);
        L && (x.content = L.innerHTML);
      } else if (x.target) {
        const L = document.querySelector(x.target);
        L && (x.content = L.innerHTML);
      }
    };
    C(() => {
      M((k) => {
        D(k);
      });
    });
    const I = () => {
      const k = n.store("tooltip"), L = B();
      u.setAttribute("aria-describedby", "uno-tooltip-content"), k.show(u, x.content, {
        ...x,
        placement: U(),
        trigger: L,
        eventEl: R()
      });
    }, z = () => {
      u.removeAttribute("aria-describedby"), n.store("tooltip").hide();
    }, T = () => {
      const k = n.store("tooltip");
      k.visible && k.triggerEl === u ? k.hide() : I();
    }, H = () => {
      B() === "click" && T();
    }, q = () => {
      B() === "hover" && I();
    }, V = () => {
      B() === "hover" && z();
    }, F = () => {
      B() === "focus" && I();
    }, O = () => {
      B() === "focus" && z();
    }, S = R();
    S.addEventListener("click", H), S.addEventListener("mouseenter", q), S.addEventListener("mouseleave", V), S.addEventListener("focus", F), S.addEventListener("blur", O), b(() => {
      S.removeEventListener("click", H), S.removeEventListener("mouseenter", q), S.removeEventListener("mouseleave", V), S.removeEventListener("focus", F), S.removeEventListener("blur", O), z();
    });
  };
  n.directive("uno-tooltip", h), n.directive("tooltip", h);
  const w = {
    show: (u, y, p) => n.store("tooltip").show(u, y, p),
    hide: () => n.store("tooltip").hide(),
    updatePosition: () => n.store("tooltip").updatePosition(),
    destroy: () => n.store("tooltip").destroy()
  };
  g.tooltip = w, typeof window < "u" && window.Uno && (window.Uno.tooltip = w);
};
const pt = (n, s = {}) => {
  const l = s.plugins?.notify || s.notify || {}, i = () => {
    const e = typeof window < "u" ? window.Uno?.options : g.options;
    return e?.plugins?.notify || e?.notify || l;
  }, d = () => g.options?.csp?.nonce || typeof window < "u" && window.Uno?.options?.csp?.nonce, t = (e) => g.getStyle("notify", `types.${e}`, ""), r = (e) => g.getStyle("notify", `icons.${e}`, "");
  let a = null;
  const o = () => {
    if (typeof document > "u") return null;
    if (a && document.contains(a)) return a;
    const e = document.getElementById("uno-notify-container") || document.querySelector("[data-uno-notify]");
    return e && (a = e), e;
  }, m = () => {
    const e = o();
    if (!e) return {};
    const c = e.getAttribute("data-uno-position") || void 0, f = e.getAttribute("data-uno-timeout"), h = f !== null && f !== "" ? Number(f) : void 0, w = e.getAttribute("data-uno-type") || void 0;
    return {
      position: c,
      timeout: Number.isFinite(h) ? h : void 0,
      type: w
    };
  };
  n.store("notify", {
    items: [],
    nextId: 1,
    activeRequests: 0,
    activeProgress: 0,
    // Notification Methods
    show(e, c = {}) {
      const f = i(), h = m(), {
        type: w = c.type || f.type || h.type || "info",
        timeout: u = c.timeout ?? (f.duration !== void 0 ? f.duration : h.timeout ?? 3e3),
        position: y = c.position || f.position || h.position || "top-right",
        closeIndicator: p = !1,
        showCloseButton: E = !0
      } = c, C = this.nextId++, b = {
        id: C,
        text: e,
        type: w,
        position: y,
        timeout: u,
        closeIndicator: p,
        showCloseButton: E,
        visible: !0,
        indicatorWidth: 100
      };
      this.items.unshift(b), u > 0 && setTimeout(() => {
        this.close(C);
      }, u), this.ensureContainer();
    },
    close(e) {
      this.items.find((f) => f.id === e) && (this.items = this.items.filter((f) => f.id !== e));
    },
    success(e, c = {}) {
      const f = typeof c == "number" ? { timeout: c } : c;
      this.show(e, { type: "success", ...f });
    },
    warning(e, c = {}) {
      const f = typeof c == "number" ? { timeout: c } : c;
      this.show(e, { type: "warning", ...f });
    },
    error(e, c = {}) {
      const f = typeof c == "number" ? { timeout: c } : c;
      this.show(e, { type: "error", ...f });
    },
    info(e, c = {}) {
      const f = typeof c == "number" ? { timeout: c } : c;
      this.show(e, { type: "info", ...f });
    },
    loading(e) {
      e ? this.activeRequests++ : this.activeRequests = Math.max(0, this.activeRequests - 1), this.ensureContainer();
    },
    progress(e) {
      e ? this.activeProgress++ : this.activeProgress = Math.max(0, this.activeProgress - 1), this.ensureContainer();
    },
    // Style Helpers
    getTypeClass(e) {
      return t(e);
    },
    getIconClass(e) {
      return r(e);
    },
    ensureContainer() {
      if (a && document.contains(a) || (a = document.getElementById("uno-notify-container"), a) || (a = document.querySelector("[data-uno-notify]"), a)) return;
      let e = Number(g.options?.ui?.layers?.notifyBase ?? 3e4);
      Number.isFinite(e) || (e = 3e4);
      const c = document.createElement("div");
      c.id = "uno-notify-container", c.setAttribute("data-uno-notify", "");
      const f = d();
      f && c.setAttribute("nonce", f), c.setAttribute("x-data", "{}");
      const h = g.getStyle("notify", "item"), w = g.getStyle("notify", "content"), u = g.getStyle("notify", "row"), y = g.getStyle("notify", "iconWrap"), p = g.getStyle("notify", "textWrap"), E = g.getStyle("notify", "text"), C = g.getStyle("notify", "closeWrap"), b = g.getStyle("notify", "closeBtn"), M = g.getStyle("notify", "closeIcon"), x = g.getStyle("notify", "indicatorTrack"), U = g.getStyle("notify", "indicatorBar"), B = g.getStyle("notify", "progressWrap"), R = g.getStyle("notify", "progressBar"), D = g.getStyle("notify", "loadingWrap"), I = g.getStyle("notify", "spinner"), z = (T) => {
        const H = ct(T), q = ut(T), V = T === "top-center" || T === "bottom-center" || T === "center" ? "w-full" : "";
        return `
          <div class="${`${g.getStyle("notify", "container")} ${q.flexDir} ${q.alignItems} ${V} ${H}`}" style="z-index: ${e};" x-show="$store.notify.itemsFor('${T}').length">
            <template x-for="item in $store.notify.itemsFor('${T}')" :key="item.id">
              <div
                data-testid="notify-item"
                x-show="item.visible"
                x-transition:enter="transition ease-out duration-300"
                x-transition:enter-start="opacity-0 transform translate-y-2 scale-95"
                x-transition:enter-end="opacity-100 transform translate-y-0 scale-100"
                x-transition:leave="transition ease-in duration-200"
                x-transition:leave-start="opacity-100 transform translate-y-0 scale-100"
                x-transition:leave-end="opacity-0 transform translate-y-2 scale-95"
                class="${h}"
                :class="$store.notify.getTypeClass(item.type)"
              >
                <div class="${w}">
                  <div class="${u}">
                    <div class="${y}">
                      <template x-if="item.type === 'success'">
                        <svg class="h-6 w-6" :class="$store.notify.getIconClass('success')" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                      </template>
                      <template x-if="item.type === 'error'">
                        <svg class="h-6 w-6" :class="$store.notify.getIconClass('error')" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                      </template>
                      <template x-if="item.type === 'warning'">
                        <svg class="h-6 w-6" :class="$store.notify.getIconClass('warning')" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                        </svg>
                      </template>
                      <template x-if="item.type === 'info'">
                        <svg class="h-6 w-6" :class="$store.notify.getIconClass('info')" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                      </template>
                    </div>
                    <div class="${p}">
                      <p class="${E}" x-text="item.text"></p>
                    </div>
                    <div class="${C}" x-show="item.showCloseButton">
                      <button @click="$store.notify.close(item.id)" class="${b}">
                        <span class="sr-only">Close</span>
                        <svg class="${M}" viewBox="0 0 20 20" fill="currentColor">
                          <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                        </svg>
                      </button>
                    </div>
                  </div>
                </div>
                <div x-show="item.closeIndicator" class="${x}">
                   <div
                        class="${U}"
                        x-init="if(item.timeout > 0) { item.indicatorWidth = 100; $nextTick(() => (item.indicatorWidth = 0)) }"
                        :style="\`width: \${item.indicatorWidth}%; \${item.timeout > 0 ? 'transition: width ' + item.timeout + 'ms linear;' : ''}\`"
                    ></div>
                </div>
              </div>
            </template>
          </div>
        `;
      };
      c.innerHTML = `
        <div
          id="uno-progress-bar"
          x-show="$store.notify.activeProgress > 0"
          class="${B}"
          style="z-index: ${e + 1};"
        >
          <div class="${R}"></div>
        </div>
        <div
          id="uno-loading-indicator"
          x-show="$store.notify.activeRequests > 0"
          class="${D}"
          style="z-index: ${e + 2};"
        >
            <svg class="${I}" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
        </div>
        ${dt.map(z).join("")}
      `, document.body.appendChild(c), a = c, typeof n?.initTree == "function" && n.initTree(c);
    },
    itemsFor(e) {
      const c = e || "top-right";
      return this.items.filter((f) => (f.position || "top-right") === c);
    }
  }), typeof window < "u" && typeof document < "u" && typeof g.htmx?.onRequest == "function" && g.htmx.available() && g.htmx.onRequest({
    start: () => {
      const e = n.store("notify");
      e.loading(!0), e.progress(!0);
    },
    end: () => {
      const e = n.store("notify");
      e.loading(!1), e.progress(!1);
    }
  }), typeof window < "u" && typeof document < "u" && N(document, "uno:notify:event") && document.addEventListener("uno:notify", (e) => {
    const c = e?.detail || {}, f = c.text || c.message || c.title || "";
    f && n.store("notify").show(f, c);
  });
  const v = {
    show: (e, c) => n.store("notify").show(e, c),
    success: (e, c) => n.store("notify").success(e, c),
    error: (e, c) => n.store("notify").error(e, c),
    warning: (e, c) => n.store("notify").warning(e, c),
    info: (e, c) => n.store("notify").info(e, c),
    loading: (e) => n.store("notify").loading(e),
    progress: (e) => n.store("notify").progress(e),
    close: (e) => n.store("notify").close(e),
    closeAll: () => {
      n.store("notify").items = [];
    }
  };
  g.notify = v, typeof window < "u" && window.Uno && (window.Uno.notify = v);
}, gt = [
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
function ht(n, s = {}) {
  const l = (o) => {
    if (!(o instanceof HTMLElement)) return !1;
    const m = window.getComputedStyle(o);
    return m.display === "none" || m.visibility === "hidden" ? !1 : o.getClientRects().length > 0;
  }, i = () => Array.from(n.querySelectorAll(gt)).filter((o) => l(o)), d = () => i()[0], t = () => {
    const o = i();
    return o[o.length - 1];
  }, r = (o) => {
    if (o.key === "Tab") {
      if (o.shiftKey) {
        if (document.activeElement === d()) {
          o.preventDefault();
          const m = t();
          m && m.focus();
        }
      } else if (document.activeElement === t()) {
        o.preventDefault();
        const m = d();
        m && m.focus();
      }
    }
  }, a = () => {
    if (n.contains(document.activeElement)) return;
    const o = d();
    o && o.focus();
  };
  return typeof requestAnimationFrame == "function" ? requestAnimationFrame(() => requestAnimationFrame(a)) : setTimeout(a, 0), n.addEventListener("keydown", r), () => {
    n.removeEventListener("keydown", r);
  };
}
const yt = (n, s = {}) => {
  const l = s.plugins?.box || s.box || {}, i = () => {
    const t = typeof window < "u" ? window.Uno?.options : g.options;
    return t?.plugins?.box || t?.box || l;
  }, d = (t) => g.getStyle("box", t);
  n.directive("uno-trap", (t, { expression: r }, { cleanup: a }) => {
    const o = ht(t);
    a(o);
  }), n.directive("uno-box", (t, { expression: r }, { evaluateLater: a, cleanup: o }) => {
    const m = r ? a(r) : null, v = (h) => {
      if (h == null) return;
      const w = String(h).trim().toLowerCase();
      if (w === "true" || w === "1") return !0;
      if (w === "false" || w === "0") return !1;
    }, e = (h) => {
      if (h == null || h === "") return;
      const w = Number(h);
      return Number.isFinite(w) ? w : void 0;
    }, c = () => {
      const h = t.getAttribute("data-uno-box") || t.getAttribute("data-uno-box-type") || "dialog", w = t.getAttribute("data-uno-title") || "", u = t.getAttribute("data-uno-content") || "", y = t.getAttribute("data-uno-hx-url") || void 0, p = t.getAttribute("data-uno-position") || void 0, E = v(t.getAttribute("data-uno-modal")), C = v(t.getAttribute("data-uno-trusted-html")), b = e(t.getAttribute("data-uno-width")), M = e(t.getAttribute("data-uno-height")), x = t.getAttribute("data-uno-width-class") || void 0, U = t.getAttribute("data-uno-height-class") || void 0;
      return {
        type: h,
        title: w,
        content: u,
        hxUrl: y,
        position: p,
        modal: E,
        trustedHTML: C,
        width: b,
        height: M,
        widthClass: x,
        heightClass: U
      };
    }, f = (h) => {
      h.preventDefault();
      const w = n.store("box");
      if (w) {
        if (m) {
          let u;
          if (m((y) => {
            u = y;
          }), u && typeof u == "object") {
            w.open(u);
            return;
          }
        }
        w.open(c());
      }
    };
    t.addEventListener("click", f), o(() => t.removeEventListener("click", f));
  }), n.store("box", {
    items: [],
    // Public API
    alert(t, r = {}) {
      const a = typeof r == "string" ? { title: r } : r;
      return this.open("dialog", {
        content: t,
        showOkButton: !0,
        ...a
      });
    },
    confirm(t, r = {}) {
      return new Promise((a) => {
        this.open("confirm", {
          content: t,
          onOk: () => a(!0),
          onCancel: () => a(!1),
          onClose: () => a(!1),
          // Closing without choice = cancel
          ...r
        });
      });
    },
    prompt(t, r = "", a = {}) {
      return new Promise((o) => {
        const m = "prompt-" + Math.random().toString(36).slice(2, 11), v = d("promptWrapper"), e = d("promptText"), c = d("promptInput"), f = `
          <div class="${v}">
            <p class="${e}">${t}</p>
            <input id="${m}" type="text" class="${c}" value="${r}">
          </div>
        `;
        this.open("confirm", {
          title: a.title || "Prompt",
          content: f,
          onOk: () => {
            const h = document.getElementById(m);
            o(h ? h.value : null);
          },
          onCancel: () => o(null),
          ...a
        });
      });
    },
    async open(t, r = {}) {
      let a = "dialog", o = r;
      const m = i();
      typeof t == "object" ? (o = t, a = o.type || "dialog") : a = t;
      const e = {
        id: o.boxId || Math.random().toString(36).slice(2, 10),
        type: a,
        title: o.title || "",
        content: o.content || "",
        trustedHTML: o.trustedHTML === !0,
        width: o.width !== void 0 ? o.width : o.widthClass ? void 0 : m.defaultWidth || 400,
        height: o.height,
        widthClass: o.widthClass,
        heightClass: o.heightClass,
        position: o.position || m.defaultPosition || "center",
        modal: o.modal !== void 0 ? o.modal : !0,
        showCloseButton: o.showCloseButton !== void 0 ? o.showCloseButton : !0,
        showOkButton: o.showOkButton || !1,
        showCancelButton: o.showCancelButton || !1,
        okText: o.okText || m.okText || "OK",
        cancelText: o.cancelText || m.cancelText || "Cancel",
        onOk: o.onOk,
        onCancel: o.onCancel,
        onClose: o.onClose,
        hxUrl: o.hxUrl,
        hxHeaders: o.hxHeaders,
        // Internal state
        isOpen: !1,
        // Methods
        close() {
          this.isOpen = !1, setTimeout(() => {
            n.store("box").remove(this.id), this.onClose && this.onClose();
          }, 300);
        }
      };
      return this.items.push(e), this.ensureContainer(), setTimeout(() => {
        const c = this.items.find((f) => f.id === e.id);
        c && (c.isOpen = !0);
      }, 50), e;
    },
    remove(t) {
      this.items = this.items.filter((r) => r.id !== t);
    },
    closeAll() {
      this.items = [];
    },
    ensureContainer() {
      const t = "uno-box-manager";
      if (typeof document < "u" && document.getElementById(t)) return;
      const a = {
        "x-data": "",
        class: "relative",
        style: `z-index: ${g.options?.ui?.layers?.boxBase || 1e4}; display: none;`,
        "x-show": "$store.box.items.length > 0",
        role: "dialog",
        "aria-modal": "true"
      }, m = g.dom.ensureContainer(t, `
             <!-- Backdrop -->
             <div :class="Uno.getStyle('box', 'overlay')"
                  x-show="$store.box.items.length > 0"
                  x-transition:enter="ease-out duration-300"
                  x-transition:enter-start="opacity-0"
                  x-transition:enter-end="opacity-100"
                  x-transition:leave="ease-in duration-200"
                  x-transition:leave-start="opacity-100"
                  x-transition:leave-end="opacity-0"></div>

             <div :class="Uno.getStyle('box', 'dialogWrapper')">
                <template x-for="item in $store.box.items" :key="item.id">
                   <div :id="'uno-box-outer-container' + item.id"
                        :class="Uno.getStyle('box', 'itemWrapper')"
                        x-show="item.isOpen">

                        <div :id="'uno-box-container' + item.id"
                             :class="[Uno.getStyle('box', 'container'), item.widthClass]"
                             :style="item.width ? 'width:' + item.width + 'px; max-width:' + item.width + 'px;' : ''"
                             x-show="item.isOpen">

                             <!-- Content -->
                             <div :class="Uno.getStyle('box', 'content')">
                                 <h3 :class="Uno.getStyle('box', 'title')"
                                     x-text="item.title"
                                     x-show="item.title"></h3>
                                 <div class="mt-2" :id="'uno-box-content' + item.id">
                                    <template x-if="!item.hxUrl">
                                      <div x-html="item.trustedHTML ? item.content : Uno.dom.sanitize(item.content)"></div>
                                    </template>
                                    <template x-if="item.hxUrl">
                                      <div x-init="htmx.process($el)"
                                           :hx-get="item.hxUrl"
                                           hx-trigger="load"
                                           hx-swap="innerHTML"
                                           :hx-headers="item.hxHeaders ? JSON.stringify(item.hxHeaders) : null">
                                        <!-- Loading state could go here -->
                                        <div class="animate-pulse">Loading...</div>
                                      </div>
                                    </template>
                                 </div>
                             </div>

                             <!-- Footer -->
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
      `, a);
      m && n && typeof n.initTree == "function" && n.initTree(m);
    }
  });
}, K = {
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
}, wt = (n) => {
  n.store("unoValid", {
    errors: {},
    setError(t, r) {
      this.errors[t] = r;
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
    const r = t.detail.elt, a = r.tagName === "FORM" ? r : r.closest("form");
    a && a.hasAttribute("data-uno-valid") && (a.checkValidity() || (t.preventDefault(), t.stopPropagation(), l(a), a.dispatchEvent(new CustomEvent("uno:validation-failed", { bubbles: !0 }))));
  }), document.body.addEventListener("input", (t) => {
    const r = t.target;
    r.form && r.form.hasAttribute("data-uno-valid") && (r.getAttribute("aria-invalid") === "true" && i(r), r.name && r.form.querySelectorAll(`[data-uno-match="${r.name}"]`).forEach((o) => {
      (o.getAttribute("aria-invalid") === "true" || o.value) && i(o);
    }));
  }), document.body.addEventListener(
    "blur",
    (t) => {
      const r = t.target;
      (r.tagName === "INPUT" || r.tagName === "TEXTAREA" || r.tagName === "SELECT") && r.form && r.form.hasAttribute("data-uno-valid") && i(r);
    },
    !0
  ), document.body.addEventListener("input", (t) => {
    const r = t.target;
    r.matches && r.matches("input, textarea") && s(r);
  }));
  function s(t) {
    let r = null;
    if (t.hasAttribute("data-uno-count")) {
      const a = t.getAttribute("data-uno-count");
      a && (r = document.querySelector(a));
    }
    if (!r && t.form) {
      const a = t.closest("div");
      a && (r = a.querySelector(".valid-count span") || a.querySelector(".valid-count"));
    }
    if (r) {
      const a = t.value.length;
      r.textContent = a;
    }
  }
  function l(t) {
    const r = t.querySelectorAll("input, select, textarea");
    let a = null;
    return r.forEach((o) => {
      i(o) || a || (a = o);
    }), a && a.focus(), !a;
  }
  function i(t) {
    if (t.disabled || t.type === "hidden") return !0;
    t.setCustomValidity("");
    let r = [];
    if (t.hasAttribute("valid"))
      try {
        const e = JSON.parse(t.getAttribute("valid"));
        Array.isArray(e) && (r = e);
      } catch (e) {
        console.warn("Uno Valid: Invalid JSON rules", e);
      }
    const a = t.getAttribute("data-valid-type");
    a && r.push({ type: a }), t.hasAttribute("required") && !r.some((e) => e.type === "required") && r.unshift({ type: "required" });
    for (const e of r) {
      if (e.type === "emptyquit") {
        if (!t.value) {
          t.setCustomValidity("");
          break;
        }
        continue;
      }
      if (e.type === "required") {
        if (t.validity.valueMissing) {
          t.setCustomValidity(
            e.msg || d("valid.required") || "This field is required"
          );
          break;
        }
        continue;
      }
      if (K[e.type]) {
        if (t.value && !K[e.type].test(t.value)) {
          t.setCustomValidity(e.msg || d(`valid.${e.type}`) || "Invalid format");
          break;
        }
        continue;
      }
      if (typeof window < "u" && window.UnoValidators && window.UnoValidators[e.type]) {
        const c = window.UnoValidators[e.type](t.value, t);
        if (c !== !0 && c !== void 0 && c !== "") {
          const f = typeof c == "string" ? c : d(`valid.${e.type}`) || "Validation failed";
          t.setCustomValidity(e.msg || f);
          break;
        }
        continue;
      }
      if (typeof window < "u" && window[e.type] && typeof window[e.type] == "function") {
        const c = window[e.type](t.value, t);
        if (c !== !0 && c !== void 0 && c !== "") {
          const f = typeof c == "string" ? c : "Validation failed";
          t.setCustomValidity(e.msg || f);
          break;
        }
        continue;
      }
    }
    if (t.hasAttribute("data-uno-match")) {
      const e = t.getAttribute("data-uno-match"), c = t.form.querySelector(`[name="${e}"]`);
      c && t.value !== c.value && t.setCustomValidity(t.getAttribute("data-uno-error-match") || "Fields do not match");
    }
    const o = t.checkValidity(), m = g.getStyle("valid", "errorClass", "is-invalid");
    o ? t.classList.remove(m) : t.classList.add(m);
    const v = g.getStyle("valid", "parentErrorClass");
    if (v) {
      const c = (typeof window < "u" && window.Uno?.options?.valid || {}).parentSelector || ".field, .form-group, .input-group, div", f = t.closest(c);
      f && (o ? f.classList.remove(v) : f.classList.add(v));
    }
    if (t.setAttribute("aria-invalid", !o), o)
      t.name && n.store("unoValid").clearError(t.name);
    else {
      let e = t.validationMessage;
      if (t.validity.valueMissing)
        e = t.getAttribute("data-uno-error-required") || d("valid.required") || e;
      else if (t.validity.typeMismatch) {
        let c = "valid.pattern";
        t.type === "email" && (c = "valid.email"), t.type === "url" && (c = "valid.url"), e = t.getAttribute("data-uno-error-type") || d(c) || e;
      } else t.validity.patternMismatch ? e = t.getAttribute("data-uno-error-pattern") || d("valid.pattern") || e : t.validity.tooShort ? e = t.getAttribute("data-uno-error-minlength") || d("valid.min") || e : t.validity.tooLong ? e = t.getAttribute("data-uno-error-maxlength") || d("valid.max") || e : t.validity.rangeUnderflow ? e = t.getAttribute("data-uno-error-min") || d("valid.rangeMin") || e : t.validity.rangeOverflow && (e = t.getAttribute("data-uno-error-max") || d("valid.rangeMax") || e);
      t.name && n.store("unoValid").setError(t.name, e);
    }
    return o;
  }
  function d(t) {
    return typeof window < "u" && window.Alpine && window.Alpine.store("i18n") ? window.Alpine.store("i18n").t(t) : null;
  }
}, xt = (n) => {
  typeof document < "u" && document.body.addEventListener("htmx:beforeSwap", (s) => {
    const l = s.detail.requestConfig.elt;
    l && l.getAttribute("role") === "tab" && (s.detail.serverResponse = _(s.detail.serverResponse));
  }), n.data("unoTabs", (s = "", l = {}) => ({
    activeTab: s,
    init() {
      this.$nextTick(() => {
        if (l.hash !== !1) {
          const i = () => {
            const d = window.location.hash.slice(1);
            d ? this.open(d, { updateHash: !1 }) : this.activeTab && this.load(this.activeTab);
          };
          i(), window.addEventListener("hashchange", i);
        } else
          this.activeTab && this.load(this.activeTab);
      });
    },
    isActive(i) {
      return this.activeTab === i;
    },
    /**
     * Load tab content if it's an AJAX tab.
     * Uses a retry mechanism to ensure DOM attributes are applied by Alpine before querying.
     */
    load(i) {
      const d = (t = 0) => {
        const r = this.$root.querySelector(`[role="tab"][data-tab="${i}"]`);
        if (!r) {
          if (t < 5) {
            setTimeout(() => d(t + 1), 20);
            return;
          }
          return;
        }
        const a = r.getAttribute("data-hx-url");
        if (!a) return;
        const o = r.getAttribute("data-hx-target") || `#panel-${i}`, m = document.querySelector(o);
        typeof htmx < "u" && m ? htmx.ajax("GET", a, {
          target: o,
          swap: "innerHTML",
          source: r
          // Allows hx-indicator to work if placed on the tab button
        }).then(() => {
        }) : (typeof htmx > "u" && console.error("[Uno] HTMX not found"), m || console.error(`[Uno] Target element not found: ${o}`));
      };
      d();
    },
    open(i, d = {}) {
      const { updateHash: t = !0 } = d;
      this.activeTab = i, this.load(i), l.hash !== !1 && t && i && window.location.hash.slice(1) !== i && window.history.replaceState(null, null, `#${i}`);
    },
    trigger(i, d = {}) {
      const t = typeof window < "u" && window.Uno?.options?.ui?.tab || {}, r = t.triggerActive || "is-active", a = t.triggerInactive || "", o = {
        role: "tab",
        "data-tab": i,
        // Critical for locating the element later
        ":aria-selected": `isActive('${i}').toString()`,
        ":class": `isActive('${i}') ? '${r}' : '${a}'`,
        "@click.prevent": `open('${i}')`
      };
      return d.url && (o["data-hx-url"] = d.url, o["data-hx-target"] = d.target || `#panel-${i}`), o;
    },
    tablist: {
      role: "tablist"
    },
    panel(i) {
      return {
        "x-show": `isActive('${i}')`,
        role: "tabpanel",
        "aria-labelledby": `tab-${i}`
      };
    }
  }));
};
g.dom = G;
typeof window < "u" && (window.Uno = g);
g.use(lt).use(mt).use(pt).use(yt).use(wt).use(xt);
if (typeof window < "u" && typeof document < "u") {
  const n = () => {
    if (g.initialized) return;
    const s = window.UnoConfig || {};
    g.init(window.Alpine, s);
  };
  window.Alpine ? n() : document.addEventListener("alpine:init", n);
}
export {
  g as default
};
