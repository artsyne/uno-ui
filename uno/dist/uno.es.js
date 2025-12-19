import O from "alpinejs";
import J from "htmx.org";
let W = !0, R = null;
const C = {
  /**
   * Set debug mode
   * @param {boolean} debug
   */
  setDebug(e) {
    W = e;
  },
  /**
   * Set custom error reporter
   * @param {Function} reporter - (error, ...args) => void
   */
  setReporter(e) {
    R = e;
  },
  /**
   * Log info message
   * @param {...any} args - Arguments to log
   */
  info(...e) {
    W && console.log("[Uno Info]", ...e);
  },
  /**
   * Log warning message
   * @param {...any} args - Arguments to log
   */
  warn(...e) {
    console.warn("[Uno Warn]", ...e);
  },
  /**
   * Log error message
   * @param {...any} args - Arguments to log
   */
  error(...e) {
    if (console.error("[Uno Error]", ...e), R)
      try {
        R(...e);
      } catch (a) {
        console.error("[Uno Error] Reporter failed:", a);
      }
  }
};
class Q {
  constructor() {
    this.version = "0.3.0", this.plugins = [], this.initialized = !1, this.options = {
      debug: !0,
      csp: {
        nonce: null
        // For Content Security Policy
      },
      onError: null,
      // Custom error handler
      ui: {
        base: { rounded: "rounded-lg", primary: "bg-blue-600", shadow: "shadow-md" },
        box: {},
        notify: {},
        tooltip: {},
        tab: {},
        valid: {}
      }
    };
  }
  /**
   * Get style class from configuration or fallback
   * Support "+" prefix to append class to default
   * @param {string} comp - Component name
   * @param {string} elem - Element name
   * @param {string} defaultCls - Default class string
   * @returns {string} Final class string
   */
  getStyle(a, r, c) {
    const i = this.options.ui?.[a]?.[r];
    return i === void 0 ? c : typeof i == "string" && i.trim().indexOf("+") === 0 ? `${c} ${i.trim().slice(1)}`.trim() : i;
  }
  /**
   * Set global configuration (before or after init)
   * @param {Object} config - Global configuration object
   */
  setConfig(a) {
    const r = { ...this.options };
    return Object.keys(a).forEach((c) => {
      typeof a[c] == "object" && a[c] !== null && !Array.isArray(a[c]) && r[c] ? r[c] = { ...r[c], ...a[c] } : r[c] = a[c];
    }), this.options = r, this.initialized && C.info("Uno configuration updated", this.options), this;
  }
  /**
   * Register a plugin
   * @param {Function} plugin - Plugin function (Alpine, options) => void
   * @param {Object} options - Plugin options
   * @returns {UnoInstance}
   */
  use(a, r = {}) {
    return this.initialized && C.warn("Uno.use() called after initialization. Plugin might not work correctly."), this.plugins.push({ plugin: a, options: r }), this;
  }
  /**
   * Initialize Uno
   * @param {Object|Function} alpineOrOptions - Alpine instance or Configuration options
   * @param {Object} options - Configuration options (if first arg is Alpine)
   * @returns {Object} Alpine instance
   */
  init(a = {}, r = {}) {
    if (this.initialized)
      return C.warn("Uno has already been initialized."), O;
    let c = O, i = r;
    return !c && typeof window < "u" && window.Alpine && (c = window.Alpine), a && (a.version || typeof a.start == "function" || typeof a.store == "function") ? c = a : i = a || r, this.options = { ...this.options, ...i }, C.setDebug(this.options.debug !== !1), this.options.onError && C.setReporter(this.options.onError), C.info(`Uno v${this.version} initializing...`), this.plugins.forEach(({ plugin: o, options: s }) => {
      try {
        const n = { ...this.options, ...s };
        o(c, n);
      } catch (n) {
        C.error("Failed to load plugin", n);
      }
    }), this.initialized = !0, C.info("Uno initialized"), c;
  }
}
const p = new Q(), tt = {
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
const Z = {
  en: tt,
  "zh-CN": et,
  "zh-TW": nt
}, _ = Object.keys(Z);
function it() {
  if (typeof localStorage < "u") {
    const e = localStorage.getItem("uno.lang");
    if (e && _.includes(e)) return e;
  }
  if (typeof navigator < "u") {
    const e = navigator.language, a = e.startsWith("zh") ? e.includes("TW") || e.includes("Hant") ? "zh-TW" : "zh-CN" : "en";
    return _.includes(a) ? a : "en";
  }
  return "en";
}
function ot(e, a) {
  const r = a.split(".");
  let c = e;
  for (const i of r) {
    if (!c || typeof c != "object" || !(i in c)) return;
    c = c[i];
  }
  return c;
}
function G(e, a) {
  for (const r in a)
    a[r] && typeof a[r] == "object" ? (e[r] || Object.assign(e, { [r]: {} }), G(e[r], a[r])) : Object.assign(e, { [r]: a[r] });
}
function st(e) {
  e.store("i18n", {
    locale: it(),
    messages: Z,
    init() {
      typeof document < "u" && (document.documentElement.lang = this.locale);
    },
    t(a, r = {}) {
      const c = this.locale, i = this.messages[c] || this.messages.en, o = ot(i, a);
      return (typeof o == "string" ? o : a).replace(/\{(\w+)\}/g, (n, t) => t in r ? String(r[t]) : `{${t}}`);
    },
    setLang(a) {
      if (!this.messages[a]) {
        console.warn(`[Uno] Language ${a} not found`);
        return;
      }
      this.locale = a, typeof localStorage < "u" && localStorage.setItem("uno.lang", a), typeof document < "u" && (document.documentElement.lang = a);
    },
    extend(a) {
      for (const r in a)
        this.messages[r] ? G(this.messages[r], a[r]) : this.messages[r] = a[r];
    }
  });
}
function B(e) {
  return String(e ?? "").replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;").replace(/'/g, "&#39;");
}
const V = /* @__PURE__ */ new WeakMap();
function A(e, a) {
  if (!e) return !1;
  const r = String(a);
  let c = V.get(e);
  return c || (c = /* @__PURE__ */ new Set(), V.set(e, c)), c.has(r) ? !1 : (c.add(r), !0);
}
function S(e, a = {}) {
  const {
    allowedTags: r = [
      "b",
      "i",
      "strong",
      "em",
      "u",
      "br",
      "p",
      "ul",
      "ol",
      "li",
      "code",
      "pre",
      "a",
      "span",
      "div"
    ],
    allowedAttrs: c = { a: ["href", "target", "rel"], "*": ["class"] }
  } = a, o = String(e ?? "").replace(/<(script|style|iframe|object|embed)\b[\s\S]*?<\/\1\s*>/gi, "").replace(/<(video|audio|svg|math)\b[\s\S]*?<\/\1\s*>/gi, "").replace(/<(img|link|meta|base|source|track)\b[^>]*\/?>/gi, ""), s = document.createElement("template");
  s.innerHTML = o;
  const n = document.createTreeWalker(s.content, NodeFilter.SHOW_ELEMENT), t = [], d = [];
  for (; n.nextNode(); ) {
    const l = n.currentNode, f = l.tagName.toLowerCase();
    if (["script", "style", "iframe", "object", "embed"].includes(f)) {
      d.push(l);
      continue;
    }
    if (!r.includes(f)) {
      t.push(l);
      continue;
    }
    for (const m of [...l.attributes]) {
      const u = m.name.toLowerCase();
      (c[f] || []).concat(c["*"] || []).includes(u) || l.removeAttribute(m.name);
    }
    if (f === "a") {
      const m = (l.getAttribute("href") || "").trim();
      m.startsWith("/") || m.startsWith("https://") || m.startsWith("http://") || m.startsWith("#") || l.removeAttribute("href"), l.getAttribute("rel") || l.setAttribute("rel", "noopener noreferrer");
    }
  }
  for (const l of d)
    l.remove();
  for (const l of t) {
    const f = l.textContent || "";
    f ? l.replaceWith(document.createTextNode(f)) : l.remove();
  }
  return s.innerHTML;
}
function z(e, a, { trusted: r = !1, alpine: c } = {}) {
  e && (e.innerHTML = r ? String(a ?? "") : S(a), c?.initTree?.(e));
}
const N = {
  /**
   * Ensure a container element exists in the DOM.
   * @param {string} id - Element ID
   * @param {string} html - Inner HTML content
   * @param {Object} attributes - Attributes to set
   * @param {HTMLElement} parent - Parent element (default: document.body)
   */
  ensureContainer(e, a = "", r = {}, c = document.body) {
    if (!document.getElementById(e))
      try {
        const i = document.createElement("div");
        return i.id = e, window.Uno?.options?.csp?.nonce && !r.nonce && (r.nonce = window.Uno.options.csp.nonce), Object.entries(r).forEach(([o, s]) => {
          i.setAttribute(o, s);
        }), a && (i.innerHTML = a), c ? c.appendChild(i) : document.documentElement.appendChild(i), i;
      } catch (i) {
        C.error(`DomManager: Failed to inject container #${e}`, i);
      }
  },
  /**
   * Inject CSS styles with CSP nonce support.
   * @param {string} id - Style element ID
   * @param {string} css - CSS content
   */
  injectStyles(e, a) {
    if (!document.getElementById(e))
      try {
        const r = document.createElement("style");
        r.id = e, window.Uno?.options?.csp?.nonce && r.setAttribute("nonce", window.Uno.options.csp.nonce), r.textContent = a, document.head.appendChild(r);
      } catch (r) {
        C.error(`DomManager: Failed to inject styles #${e}`, r);
      }
  },
  // Exports
  setHTML: z,
  // Secure HTML setting
  sanitize: S,
  // HTML sanitization
  escape: B
  // HTML escaping
}, rt = [
  "top-left",
  "top-center",
  "top-right",
  "center-left",
  "center",
  "center-right",
  "bottom-left",
  "bottom-center",
  "bottom-right"
], at = (e) => {
  const a = e || "center", r = {
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
  return r[a] || r.center;
}, lt = (e) => {
  const a = e || "top-right", r = {
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
  return r[a] || r["top-right"];
}, ct = (e) => {
  const a = e || "top-right";
  let r = "flex-col", c = "items-end", i = "space-y-2";
  return a.includes("bottom") && (r = "flex-col-reverse", i = "space-y-reverse space-y-2"), a.includes("left") ? c = "items-start" : a.includes("center") ? a === "center" ? c = "items-center justify-center" : c = "items-center" : c = "items-end", { flexDir: r, alignItems: c, spaceY: i };
};
function dt(e, a, r) {
  const c = e.getBoundingClientRect(), i = a.getBoundingClientRect(), o = i.width, s = i.height, n = c.width, t = c.height, d = { top: c.top, left: c.left }, l = window.innerWidth, f = window.innerHeight, m = 0, u = 0, b = f, h = l, g = (k, $) => {
    if ($) {
      if (k + o > h)
        return h - o - 10;
      if (k < u)
        return u + 10;
    } else {
      if (k + s > b)
        return b - s - 10;
      if (k < m)
        return m + 10;
    }
    return k;
  };
  let y, w, T = r;
  switch (r) {
    case "left":
      w = d.left - o - 5, w < u && (w = d.left + n + 5, T = "right"), y = g(d.top + t / 2 - s / 2, !1);
      break;
    case "top":
      y = d.top - s - 8, y < m && (y = d.top + t + 8, T = "bottom"), w = d.left + (n - o) / 2, w = g(w, !0);
      break;
    case "right":
      w = d.left + n + 5, w + o > h && (w = d.left - o - 5, T = "left"), y = g(d.top + t / 2 - s / 2, !1);
      break;
    case "bottom":
    default:
      y = d.top + t + 8, y + s > b && (y = d.top - s - 8, T = "top"), w = d.left + (n - o) / 2, w = g(w, !0);
      break;
  }
  return { top: y, left: w, position: T };
}
const ut = (e, a = {}) => {
  const r = a.tooltip || {}, c = () => typeof window < "u" && window.Uno?.options?.tooltip || r;
  if (!e) {
    C.error("Tooltip: Alpine is undefined. Plugin initialization failed.");
    return;
  }
  e.store("tooltip", {
    visible: !1,
    content: "",
    placement: "top",
    top: -9999,
    // Start off-screen
    left: -9999,
    triggerEl: null,
    onShow: null,
    async show(s, n, t = {}) {
      this.triggerEl = s, this.content = n, this.placement = t.placement || c().placement || "top", this.onShow = t.onShow, this.visible = !0, this.ensureContainer();
      const d = document.getElementById("uno-tooltip-body");
      d && z(d, this.content, { trusted: t.trustedHTML === !0, alpine: e }), await e.nextTick(), this.updatePosition(), this.onShow && this.onShow(document.getElementById("uno-tooltip-body"));
    },
    hide() {
      this.visible = !1, this.triggerEl = null;
    },
    updatePosition() {
      if (!this.triggerEl) return;
      const s = document.getElementById("uno-tooltip-content");
      if (!s) return;
      const n = dt(this.triggerEl, s, this.placement);
      this.top = n.top, this.left = n.left;
    },
    ensureContainer() {
      i();
    }
  });
  const i = () => {
    const n = `fixed px-3 py-2 text-sm font-medium ${p.getStyle("base", "rounded", "rounded-lg")} shadow-sm whitespace-nowrap pointer-events-none border bg-white text-gray-900 border-gray-200 dark:bg-slate-800 dark:text-white dark:border-slate-700 z-[20000]`, d = `
          <div id="uno-tooltip-content"
               x-show="$store.tooltip.visible"
               role="tooltip"
               x-transition:enter="transition ease-out duration-200"
               x-transition:enter-start="opacity-0 scale-95"
               x-transition:enter-end="opacity-100 scale-100"
               x-transition:leave="transition ease-in duration-150"
               x-transition:leave-start="opacity-100 scale-100"
               x-transition:leave-end="opacity-0 scale-95"
               class="${p.getStyle("tooltip", "content", n)}"
              :style="\`top: \${$store.tooltip.top}px; left: \${$store.tooltip.left}px;\`"
              style="display: none;">
              <div id="uno-tooltip-body" class="pointer-events-auto"></div>
          </div>
        `;
    N.ensureContainer("uno-tooltip-container", d, {
      "x-data": "{}"
    });
  };
  e.directive(
    "tooltip",
    (s, { expression: n, modifiers: t }, { evaluateLater: d, effect: l, cleanup: f }) => {
      const m = d(n);
      let u = {
        content: "",
        placement: null,
        trigger: null
      };
      const b = () => {
        let x = c().placement || "top";
        return u.placement && (x = u.placement), t.includes("top") && (x = "top"), t.includes("bottom") && (x = "bottom"), t.includes("left") && (x = "left"), t.includes("right") && (x = "right"), x;
      }, h = () => {
        const v = c();
        return t.includes("click") ? "click" : u.trigger ? u.trigger : v.trigger || "hover";
      }, g = (v) => {
        if (typeof v == "string" ? u.content = v : typeof v == "object" && (u = { ...u, ...v }), u.content && u.content.startsWith("#") && !u.content.includes(" ")) {
          const x = document.querySelector(u.content);
          x && (u.content = x.innerHTML);
        } else if (u.target) {
          const x = document.querySelector(u.target);
          x && (u.content = x.innerHTML);
        }
      };
      l(() => {
        m((v) => {
          g(v);
        });
      });
      const y = () => {
        const v = e.store("tooltip");
        s.setAttribute("aria-describedby", "uno-tooltip-content"), v.show(s, u.content, { ...u, placement: b() });
      }, w = () => {
        s.removeAttribute("aria-describedby"), e.store("tooltip").hide();
      }, T = () => {
        const v = e.store("tooltip");
        v.visible && v.triggerEl === s ? v.hide() : y();
      }, k = () => {
        h() === "click" && T();
      }, $ = () => {
        h() === "hover" && y();
      }, L = () => {
        h() === "hover" && w();
      }, H = (v) => {
        if (h() !== "click") return;
        const x = e.store("tooltip"), E = document.getElementById("uno-tooltip-content");
        x.visible && x.triggerEl === s && !s.contains(v.target) && (!E || !E.contains(v.target)) && x.hide();
      };
      s.addEventListener("click", k), s.addEventListener("mouseenter", $), s.addEventListener("mouseleave", L), s.addEventListener("focus", $), s.addEventListener("blur", L), window.addEventListener("click", H), f(() => {
        s.removeEventListener("click", k), s.removeEventListener("mouseenter", $), s.removeEventListener("mouseleave", L), s.removeEventListener("focus", $), s.removeEventListener("blur", L), window.removeEventListener("click", H), w();
      });
    }
  );
  const o = {
    show: (s, n, t) => e.store("tooltip").show(s, n, t),
    hide: () => e.store("tooltip").hide(),
    updatePosition: () => e.store("tooltip").updatePosition()
  };
  p.tooltip = o, typeof window < "u" && window.Uno && (window.Uno.tooltip = o);
};
const ft = (e, a = {}) => {
  const r = a.notify || {}, c = () => typeof window < "u" && window.Uno?.options?.notify || r, i = () => p.options?.csp?.nonce || typeof window < "u" && window.Uno?.options?.csp?.nonce, o = () => {
    if (document.getElementById("uno-notify-styles")) return;
    const n = document.createElement("style");
    n.id = "uno-notify-styles";
    const t = i();
    t && n.setAttribute("nonce", t), n.textContent = `
      @keyframes notify-progress {
        from { width: 100%; }
        to { width: 0%; }
      }

      @keyframes uno-progress-indeterminate {
        0% { transform: translateX(-100%); }
        100% { transform: translateX(250%); }
      }

      #uno-progress-bar .uno-progress-anim {
        animation: uno-progress-indeterminate 1.2s ease-in-out infinite;
        width: 40%;
      }
    `, document.head.appendChild(n);
  };
  e.store("notify", {
    items: [],
    nextId: 1,
    activeRequests: 0,
    activeProgress: 0,
    // Notification Methods
    show(n, t = {}) {
      const d = c(), {
        type: l = "info",
        timeout: f = d.duration !== void 0 ? d.duration : 3e3,
        position: m = t.position || d.position || "top-right",
        closeIndicator: u = !1,
        showCloseButton: b = !0
      } = t, h = this.nextId++, g = {
        id: h,
        text: n,
        type: l,
        position: m,
        timeout: f,
        closeIndicator: u,
        showCloseButton: b,
        visible: !0
        // For transition
      };
      this.items.unshift(g), f > 0 && setTimeout(() => {
        this.close(h);
      }, f), this.ensureContainer();
    },
    close(n) {
      this.items.find((d) => d.id === n) && (this.items = this.items.filter((d) => d.id !== n));
    },
    success(n, t = {}) {
      const d = typeof t == "number" ? { timeout: t } : t;
      this.show(n, { type: "success", ...d });
    },
    warning(n, t = {}) {
      const d = typeof t == "number" ? { timeout: t } : t;
      this.show(n, { type: "warning", ...d });
    },
    error(n, t = {}) {
      const d = typeof t == "number" ? { timeout: t } : t;
      this.show(n, { type: "error", ...d });
    },
    info(n, t = {}) {
      const d = typeof t == "number" ? { timeout: t } : t;
      this.show(n, { type: "info", ...d });
    },
    loading(n) {
      n ? this.activeRequests++ : this.activeRequests = Math.max(0, this.activeRequests - 1), this.ensureContainer();
    },
    progress(n) {
      n ? this.activeProgress++ : this.activeProgress = Math.max(0, this.activeProgress - 1), this.ensureContainer();
    },
    // Style Helpers
    getTypeClass(n) {
      const t = {
        success: "bg-green-50",
        error: "bg-red-50",
        warning: "bg-yellow-50",
        info: "bg-blue-50"
      };
      return p.getStyle("notify", n, t[n] || "bg-white");
    },
    getIconClass(n) {
      const t = {
        success: "text-green-400",
        error: "text-red-400",
        warning: "text-yellow-400",
        info: "text-blue-400"
      };
      return p.getStyle(
        "notify",
        `icon${n.charAt(0).toUpperCase() + n.slice(1)}`,
        t[n] || "text-gray-400"
      );
    },
    ensureContainer() {
      if (document.getElementById("uno-notify-container")) return;
      o();
      const n = document.createElement("div");
      n.id = "uno-notify-container";
      const t = i();
      t && n.setAttribute("nonce", t), n.setAttribute("x-data", "{}");
      const d = p.getStyle("base", "rounded", "rounded-lg"), l = p.getStyle(
        "notify",
        "item",
        `pointer-events-auto w-full max-w-sm overflow-hidden ${d} shadow-lg ring-1 ring-black ring-opacity-5`
      ), f = (m) => {
        const u = lt(m), b = ct(m);
        return `
          <div class="${p.getStyle(
          "notify",
          "container",
          `fixed z-[30000] flex ${b.flexDir} ${b.alignItems} gap-2 p-4 pointer-events-none ${u}`
        )}" x-show="$store.notify.itemsFor('${m}').length">
            <template x-for="item in $store.notify.itemsFor('${m}')" :key="item.id">
              <div
                data-testid="notify-item"
                x-show="item.visible"
                x-transition:enter="transition ease-out duration-300"
                x-transition:enter-start="opacity-0 transform translate-y-2 scale-95"
                x-transition:enter-end="opacity-100 transform translate-y-0 scale-100"
                x-transition:leave="transition ease-in duration-200"
                x-transition:leave-start="opacity-100 transform translate-y-0 scale-100"
                x-transition:leave-end="opacity-0 transform translate-y-2 scale-95"
                class="${l}"
                :class="$store.notify.getTypeClass(item.type)"
              >
                <div class="p-4">
                  <div class="flex items-start">
                    <div class="flex-shrink-0">
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
                    <div class="ml-3 w-0 flex-1 pt-0.5">
                      <p class="text-sm font-medium text-gray-900" x-text="item.text"></p>
                    </div>
                    <div class="ml-4 flex flex-shrink-0" x-show="item.showCloseButton">
                      <button @click="$store.notify.close(item.id)" class="inline-flex rounded-md bg-white text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
                        <span class="sr-only">Close</span>
                        <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                          <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                        </svg>
                      </button>
                    </div>
                  </div>
                </div>
                <div x-show="item.closeIndicator" class="h-1 bg-gray-200 overflow-hidden">
                   <div class="h-full bg-blue-500 origin-left"
                        :style="\`animation: notify-progress \${item.timeout}ms linear forwards\`"></div>
                </div>
              </div>
            </template>
          </div>
        `;
      };
      n.innerHTML = `
        <div
          id="uno-progress-bar"
          x-show="$store.notify.activeProgress > 0"
          class="fixed top-0 left-0 right-0 h-1 bg-transparent z-[30001]"
        >
          <div class="uno-progress-anim h-full bg-blue-500"></div>
        </div>
        ${rt.map(f).join("")}
      `, document.body.appendChild(n), typeof e?.initTree == "function" && e.initTree(n);
    },
    itemsFor(n) {
      const t = n || c().position || "top-right";
      return this.items.filter((d) => (d.position || "top-right") === t);
    }
  });
  const s = {
    show: (n, t) => e.store("notify").show(n, t),
    success: (n, t) => e.store("notify").success(n, t),
    error: (n, t) => e.store("notify").error(n, t),
    warning: (n, t) => e.store("notify").warning(n, t),
    info: (n, t) => e.store("notify").info(n, t),
    loading: (n) => e.store("notify").loading(n),
    progress: (n) => e.store("notify").progress(n),
    close: (n) => e.store("notify").close(n),
    closeAll: () => {
      e.store("notify").items = [];
    }
  };
  p.notify = s, typeof window < "u" && window.Uno && (window.Uno.notify = s);
};
class mt {
  constructor(a = 50) {
    this.maxSize = a, this.cache = /* @__PURE__ */ new Map();
  }
  get(a) {
    if (!this.cache.has(a)) return;
    const r = this.cache.get(a);
    return this.cache.delete(a), this.cache.set(a, r), r;
  }
  set(a, r) {
    if (this.cache.has(a) && this.cache.delete(a), this.cache.set(a, r), this.cache.size > this.maxSize) {
      const c = this.cache.keys().next().value;
      this.cache.delete(c);
    }
  }
  has(a) {
    return this.cache.has(a);
  }
  clear() {
    this.cache.clear();
  }
}
const M = new mt(50);
function ht(e) {
  const a = String(e ?? "").trim();
  if (!a) return [];
  try {
    const r = JSON.parse(a);
    if (typeof r == "string") return [[r, {}]];
    if (r && typeof r == "object") return Object.entries(r);
  } catch {
  }
  return a.split(",").map((r) => r.trim()).filter(Boolean).map((r) => [r, {}]);
}
function pt(e, a, r) {
  const c = e || document.body || document.documentElement;
  if (c) {
    if (typeof window < "u" && window.htmx?.trigger) {
      window.htmx.trigger(c, a, r);
      return;
    }
    c.dispatchEvent(new CustomEvent(a, { bubbles: !0, detail: r }));
  }
}
function gt(e, a) {
  const r = e.headers.get("HX-Trigger");
  if (r)
    for (const [c, i] of ht(r))
      pt(a, c, i);
}
async function bt(e, a = {}) {
  const {
    method: r = "GET",
    htmxTarget: c,
    swap: i = "innerHTML",
    htmxConfig: o = {},
    cache: s = !1
  } = a;
  if (s && r === "GET" && M.has(e))
    return M.get(e);
  if (typeof window < "u" && window.htmx?.ajax) {
    const n = c || document.createElement("div");
    c || (n.style.display = "none"), c || (document.body ? document.body.appendChild(n) : document.documentElement.appendChild(n));
    try {
      await window.htmx.ajax(r, e, { target: n, swap: i, ...o });
      const t = n.innerHTML;
      return s && r === "GET" && M.set(e, t), t;
    } finally {
      c || n.remove();
    }
  }
  return vt(e, a);
}
async function K(e, a = {}) {
  const {
    method: r = "GET",
    headers: c = {},
    body: i,
    credentials: o = "same-origin",
    triggerTarget: s = typeof document < "u" ? document.body : null,
    cache: n = !1,
    signal: t
    // AbortSignal
  } = a;
  if (n && r === "GET" && M.has(e))
    return M.get(e);
  const d = {
    "HX-Request": "true",
    "X-Requested-With": "XMLHttpRequest",
    ...c
  }, l = await fetch(e, { method: r, headers: d, body: i, credentials: o, signal: t });
  if (gt(l, s), !l.ok) {
    const u = new Error(`Request failed: ${l.status}`);
    throw u.status = l.status, u;
  }
  const f = l.headers.get("content-type") || "";
  let m;
  return f.includes("application/json") ? m = await l.json() : m = await l.text(), n && r === "GET" && M.set(e, m), m;
}
async function vt(e, a = {}) {
  const r = await K(e, a);
  return typeof r == "string" ? r : JSON.stringify(r, null, 2);
}
const wt = [
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
function yt(e, a = {}) {
  const r = (t) => {
    if (!(t instanceof HTMLElement)) return !1;
    const d = window.getComputedStyle(t);
    return d.display === "none" || d.visibility === "hidden" ? !1 : t.getClientRects().length > 0;
  }, c = () => Array.from(e.querySelectorAll(wt)).filter((t) => r(t)), i = () => c()[0], o = () => {
    const t = c();
    return t[t.length - 1];
  }, s = (t) => {
    if (t.key === "Tab") {
      if (t.shiftKey) {
        if (document.activeElement === i()) {
          t.preventDefault();
          const d = o();
          d && d.focus();
        }
      } else if (document.activeElement === o()) {
        t.preventDefault();
        const d = i();
        d && d.focus();
      }
    }
  }, n = () => {
    if (e.contains(document.activeElement)) return;
    const t = i();
    t && t.focus();
  };
  return typeof requestAnimationFrame == "function" ? requestAnimationFrame(() => requestAnimationFrame(n)) : setTimeout(n, 0), e.addEventListener("keydown", s), () => {
    e.removeEventListener("keydown", s);
  };
}
const xt = (e, a = {}) => {
  const r = a.box || {}, c = () => typeof window < "u" && window.Uno?.options?.box || r;
  e.directive("uno-trap", (o, { expression: s }, { cleanup: n }) => {
    const t = yt(o);
    n(t);
  }), e.store("box", {
    items: [],
    // Public API
    alert(o, s = {}) {
      const n = typeof s == "string" ? { title: s } : s;
      return this.open("dialog", {
        content: o,
        showOkButton: !0,
        ...n
      });
    },
    confirm(o, s = {}) {
      return new Promise((n) => {
        this.open("confirm", {
          content: o,
          onOk: () => n(!0),
          onCancel: () => n(!1),
          onClose: () => n(!1),
          // Closing without choice = cancel
          ...s
        });
      });
    },
    prompt(o, s = "", n = {}) {
      return new Promise((t) => {
        const d = "prompt-" + Math.random().toString(36).substr(2, 9), l = `
          <div class="space-y-3">
            <p>${B(o)}</p>
            <input id="${d}" type="text" class="w-full px-3 py-2 border rounded focus:ring-2 focus:ring-blue-500 outline-none" value="${B(
          s
        )}">
          </div>
        `;
        this.open("confirm", {
          title: n.title || "Prompt",
          content: l,
          trustedHTML: !0,
          onOk: () => {
            const f = document.getElementById(d);
            t(f ? f.value : null);
          },
          onCancel: () => t(null),
          ...n
        });
      });
    },
    async open(o, s = {}) {
      let n = "dialog", t = s;
      typeof o == "object" ? (t = o, n = t.type || "dialog") : n = o;
      const d = t.boxId || Math.random().toString(36).slice(2, 10), l = {
        id: d,
        type: n,
        title: t.title || "",
        content: t.content || "",
        trustedHTML: t.trustedHTML === !0,
        width: t.width !== void 0 ? t.width : t.widthClass ? void 0 : c().defaultWidth || 400,
        height: t.height,
        widthClass: t.widthClass,
        heightClass: t.heightClass,
        containerClass: t.containerClass,
        modal: t.modal !== !1,
        // default true
        overlayOpacity: t.overlayOpacity ?? 0.5,
        position: t.position || "center",
        showCloseButton: t.showCloseButton !== void 0 ? t.showCloseButton : !0,
        // Buttons
        okText: t.okButtonText,
        cancelText: t.cancelButtonText,
        onOk: t.onOk,
        onCancel: t.onCancel,
        onClose: t.onClose,
        onBeforeClose: t.onBeforeClose,
        onContentLoad: t.onContentLoad,
        showOk: n === "confirm" || n === "message" || n === "dialog" && t.showOkButton,
        showCancel: n === "confirm" || n === "dialog" && t.showCancelButton,
        // Cache Option
        cache: t.cache === !0,
        // State
        loading: !1,
        visible: !0,
        abortController: null
      };
      l.trustedHTML || (l.content = S(l.content), l.trustedHTML = !0);
      const f = e.store("i18n");
      if (l.okText || (l.okText = f ? f.t("box.ok") : "OK"), l.cancelText || (l.cancelText = f ? f.t("box.cancel") : "Cancel"), t.targetId) {
        const m = document.getElementById(t.targetId);
        m ? (l.content = m.innerHTML, t.trustedHTML || (l.content = S(l.content)), l.trustedHTML = !0) : (C.warn("Box: Target element not found:", t.targetId), l.content = '<div class="p-4 text-red-500">Error: Content target not found</div>', l.trustedHTML = !0);
      } else if (t.hxUrl) {
        l.loading = !0, this.add(l);
        try {
          const m = await bt(t.hxUrl, { method: "GET", cache: l.cache }), u = this.items.find((b) => b.id === d);
          u && (u.content = t.trustedHTML ? m : S(m), u.trustedHTML = !0, u.loading = !1, u.onContentLoad && e.nextTick(() => {
            u.onContentLoad(document.getElementById("uno-box-content" + d));
          }));
        } catch (m) {
          C.error("Box: HTMX Load Failed", m);
          const u = this.items.find((b) => b.id === d);
          u && (u.loading = !1, u.trustedHTML = !0, u.content = '<div class="p-4 text-red-500">Error loading content via HTMX</div>');
        }
        return d;
      } else if (t.ajaxUrl) {
        l.loading = !0, l.abortController = new AbortController(), this.add(l);
        try {
          const m = await K(t.ajaxUrl, {
            cache: l.cache,
            signal: l.abortController.signal
          });
          let u = typeof m == "object" ? `<pre class="bg-gray-100 p-3 rounded overflow-x-auto whitespace-pre-wrap">${JSON.stringify(m, null, 2)}</pre>` : String(m);
          t.trustedHTML || (u = S(u));
          const b = this.items.find((h) => h.id === d);
          b && (b.content = u, b.trustedHTML = !0, b.loading = !1, e.nextTick(() => {
            const h = document.getElementById("uno-box-content" + d);
            h && b.onContentLoad && b.onContentLoad(h);
          }));
        } catch (m) {
          if (m.name === "AbortError") {
            C.info("Box: Request cancelled");
            return;
          }
          C.error("Box: Ajax failed", m);
          const u = this.items.find((b) => b.id === d);
          u && (u.content = '<div class="text-red-500">Error loading content</div>', u.trustedHTML = !0, u.loading = !1);
        }
        return d;
      } else t.iframeUrl && (l.content = `<iframe frameborder="0" class="w-full h-full" src="${B(t.iframeUrl)}"></iframe>`, l.trustedHTML = !0);
      return this.add(l), d;
    },
    add(o) {
      this.items.push(o), this.ensureContainer();
    },
    close(o) {
      const s = this.items.find((n) => n.id === o);
      s && (s.abortController && s.abortController.abort(), s.onBeforeClose && s.onBeforeClose(), this.items = this.items.filter((n) => n.id !== o), s.onClose && s.onClose());
    },
    closeAll() {
      this.items.forEach((o) => {
        o.abortController && o.abortController.abort();
      }), this.items = [];
    },
    getPositionClass(o) {
      const s = "fixed inset-0 flex transition-opacity duration-300 p-4 pointer-events-none", n = at(o);
      return `${s} ${n}`;
    },
    ensureContainer() {
      i();
    }
  });
  const i = () => {
    const { rounded: o, shadow: s, primary: n } = p.options.ui.base, d = p.getStyle("box", "overlay", "absolute inset-0 bg-black transition-opacity duration-300 pointer-events-auto"), l = `relative bg-white dark:bg-gray-900 ${o} ${s} ring-1 ring-black/5 overflow-hidden flex flex-col max-h-[90vh] transition-all duration-300 transform pointer-events-auto outline-none`, f = p.getStyle("box", "container", l), u = p.getStyle("box", "cancelBtn", "px-4 py-2 bg-slate-100 dark:bg-slate-700 border border-slate-200 dark:border-slate-600 rounded-md text-sm font-medium text-slate-800 dark:text-slate-200 hover:bg-slate-200 dark:hover:bg-slate-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors"), b = `px-4 py-2 ${n} border border-transparent rounded-md text-sm font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500`, h = p.getStyle("box", "okBtn", b), g = `
      <template x-for="(item, index) in $store.box.items" :key="item.id">
          <div :id="'uno-box-outer-container' + item.id"
               :class="$store.box.getPositionClass(item.position)"
               :style="\`z-index: \${10000 + index}\`"
               role="dialog"
               aria-modal="true"
               @keydown.escape.window="if(item.modal) $store.box.close(item.id)">
               
               <!-- Overlay -->
               <div :id="'uno-box-overlay' + item.id"
                    class="${d}"
                    :style="\`opacity: \${item.overlayOpacity}\`"
                    x-show="item.modal"
                    x-transition:enter="ease-out duration-300"
                    x-transition:enter-start="opacity-0"
                    x-transition:enter-end="opacity-100"
                    x-transition:leave="ease-in duration-200"
                    x-transition:leave-start="opacity-100"
                    x-transition:leave-end="opacity-0"></div>

               <!-- Box Container -->
               <div :id="'uno-box-container' + item.id"
                    class="${f}"
                    :class="[item.widthClass, item.heightClass, item.containerClass]"
                    :style="\`\${item.width ? 'width: ' + item.width + 'px;' : ''} \${item.height ? 'height:' + item.height + 'px' : ''}\`"
                    x-show="item.visible"
                    x-transition:enter="ease-out duration-300"
                    x-transition:enter-start="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                    x-transition:enter-end="opacity-100 translate-y-0 sm:scale-100"
                    x-transition:leave="ease-in duration-200"
                    x-transition:leave-start="opacity-100 translate-y-0 sm:scale-100"
                    x-transition:leave-end="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                    x-uno-trap="item.visible"
                    tabindex="-1">
                    
                    <!-- Header -->
                    <div :id="'uno-box-header' + item.id" x-show="item.title" class="flex justify-between items-center p-4 border-b border-slate-200 dark:border-slate-700">
                        <h3 :id="'uno-box-title' + item.id" class="text-lg font-semibold text-slate-800 dark:text-slate-200" x-text="item.title"></h3>
                        <button :id="'uno-box-close' + item.id" @click="$store.box.close(item.id)" class="text-slate-400 hover:text-slate-500 focus:outline-none">
                            <span class="sr-only">Close</span>
                            <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button>
                    </div>
                    
                    <!-- Close button if no title -->
                    <button x-show="!item.title && item.showCloseButton" 
                            :id="'uno-box-close' + item.id"
                            @click="$store.box.close(item.id)" 
                            class="absolute top-2 right-2 z-10 text-slate-400 hover:text-slate-500 focus:outline-none">
                        <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>

                    <!-- Content -->
                    <div :id="'uno-box-wrapper' + item.id" class="p-4 overflow-auto flex-1 text-slate-600 dark:text-slate-300">
                        <div x-show="item.loading" class="flex justify-center items-center h-full">
                            <svg class="animate-spin h-8 w-8 text-blue-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                            </svg>
                        </div>
                        <div :id="'uno-box-content' + item.id" x-show="!item.loading" x-html="item.content"></div>
                    </div>

                    <!-- Footer -->
                    <div :id="'uno-box-footer' + item.id" x-show="item.showOk || item.showCancel" class="p-4 border-t border-slate-200 dark:border-slate-700 flex justify-end space-x-3 bg-slate-50 dark:bg-slate-800">
                        <button :id="'uno-box-cancel' + item.id"
                                x-show="item.showCancel"
                                @click="item.onCancel ? (item.onCancel(), $store.box.close(item.id)) : $store.box.close(item.id)"
                                x-text="item.cancelText"
                                class="${u}">
                        </button>
                        <button :id="'uno-box-ok' + item.id"
                                x-show="item.showOk"
                                @click="async () => {
                                    let shouldClose = true;
                                    if (item.onOk) {
                                        const result = item.onOk();
                                        shouldClose = result instanceof Promise ? await result : result;
                                    }
                                    if (shouldClose !== false) $store.box.close(item.id);
                                }"
                                x-text="item.okText"
                                class="${h}">
                        </button>
                    </div>
               </div>
          </div>
      </template>
    `;
    N.ensureContainer("uno-box-container", g, {
      "x-data": "{}"
    });
  };
  A(document, "uno:box:delegator") && document.addEventListener("click", (o) => {
    const s = o.target.closest("[data-box-type]");
    if (!s) return;
    const n = s.getAttribute("data-box-type");
    if (!["dialog", "confirm", "message"].includes(n)) return;
    o.preventDefault();
    const t = s.dataset, d = {
      title: t.boxTitle,
      content: t.boxContent,
      width: t.boxWidth ? parseInt(t.boxWidth) : void 0,
      height: t.boxHeight ? parseInt(t.boxHeight) : void 0,
      ajaxUrl: t.boxAjaxUrl,
      hxUrl: t.boxHxUrl,
      iframeUrl: t.boxIframeUrl,
      targetId: t.boxTargetId,
      position: t.boxPosition,
      okButtonText: t.boxOkText,
      cancelButtonText: t.boxCancelText,
      overlayOpacity: t.boxOverlayOpacity ? parseFloat(t.boxOverlayOpacity) : void 0,
      showCloseButton: t.boxShowClose !== "false",
      boxId: t.boxId,
      cache: t.boxCache === "true",
      trustedHTML: t.boxTrusted === "true"
    };
    t.boxModal === "false" && (d.modal = !1), e.store("box").open(n, d);
  }), typeof window < "u" && A(document, "uno:box:htmx-afterSwap") && document.addEventListener("htmx:afterSwap", (o) => {
    o.detail.elt && e && e.initTree(o.detail.elt);
  }), p.box = {
    dialog: (o) => e.store("box").open("dialog", o),
    confirm: (o, s) => e.store("box").confirm(o, s),
    alert: (o, s) => e.store("box").alert(o, s),
    prompt: (o, s, n) => e.store("box").prompt(o, s, n),
    close: (o) => e.store("box").close(o),
    closeAll: () => e.store("box").closeAll()
  };
}, Y = {
  validname: /^[a-zA-Z0-9_-]+$/,
  en: /^[a-zA-Z]+$/,
  domain: new RegExp("^((?!-)[A-Za-z0-9-]{1,63}(?<!-)\\.)+[A-Za-z]{2,6}$"),
  mobile: /^\+?[\d\s-]{7,}$/,
  int: /^-?\d+$/,
  email: /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/,
  date: /^(\d{4}-\d{2}-\d{2}|\d{2}\/\d{2}\/\d{4}|\d{2}\.\d{2}\.\d{4})$/,
  url: /^(https?:\/\/)?([\da-z.-]+)\.([a-z.-]{2,6})([/\w .-]*)*\/?$/i,
  pwd: /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$/,
  strongpwd: /^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$/
};
function Ct(e) {
  const a = e.match(/^\/(.*?)\/([gimsuy]*)$/);
  return a ? { pattern: a[1], flags: a[2] } : { pattern: e, flags: "" };
}
function kt(e, a = {}) {
  const { defaultFlags: r = "i", maxPatternLength: c = 200 } = a, { pattern: i, flags: o } = Ct(String(e ?? ""));
  if (!i || i.length > c) return null;
  const s = String(o || r || "");
  if (!/^[gimsuy]*$/.test(s)) return null;
  try {
    return new RegExp(i, s);
  } catch {
    return null;
  }
}
function Tt(e) {
  try {
    return e ? JSON.parse(e) : [];
  } catch {
    return [];
  }
}
const X = /* @__PURE__ */ new WeakMap();
const Et = (e) => {
  const a = {
    invalid: "<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 12 12' fill='none' stroke='currentColor' class='w-4 h-4'><circle cx='6' cy='6' r='4.5'/><path stroke-linecap='round' d='M5.8 3.6h.4L6 6.5zM6 8.2a.6.6 0 110-1.2.6.6 0 010 1.2z'/></svg>",
    valid: "<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 8 8' class='w-4 h-4'><path fill='currentColor' d='M2.3 6.73L.6 4.53c-.4-1.04.46-1.4 1.1-.8l1.1 1.4 3.4-3.8c.6-.63 1.6-.27 1.2.7l-4 4.6c-.43.5-.8.4-1.1.1z'/></svg>"
  }, r = (i, o) => {
    const s = e.store("i18n");
    if (!s) return i;
    const n = {
      required: "Required",
      pattern: "Invalid format",
      min: `At least ${o?.n} characters`,
      max: `Max ${o?.n} characters`,
      int: "Integer required",
      email: "Invalid email address",
      url: "Invalid URL"
    };
    return s.t(`valid.${i}`, o) || n[i] || i;
  };
  e.store("valid", {
    callbacks: {},
    registerCallback(i, o) {
      typeof o == "function" && (this.callbacks[i] = o);
    },
    async validateForm(i) {
      if (!i) return !1;
      if (i._x_dataStack) {
        const o = i._x_dataStack.find((s) => s.validateAll);
        if (o)
          return await o.validateAll();
      }
    },
    clearForm(i) {
      if (i && i._x_dataStack) {
        const o = i._x_dataStack.find((s) => s.reset);
        o && o.reset();
      }
    }
  });
  const c = {
    // Validate all forms or specific form?
    // Mostly helper to set state manually
    set: (i, { valid: o, message: s }) => {
      i && i.dispatchEvent(
        new CustomEvent("valid-set", {
          detail: { valid: o, message: s, name: i.name },
          bubbles: !0
        })
      );
    },
    validate: async (i) => {
      const o = i?.tagName === "FORM" ? i : i?.closest?.("form");
      return e.store("valid").validateForm(o);
    },
    clear: (i) => {
      const o = i?.tagName === "FORM" ? i : i?.closest?.("form");
      e.store("valid").clearForm(o);
    }
  };
  p.valid = c, typeof window < "u" && window.Uno && (window.Uno.valid = c), e.data("validForm", () => ({
    fields: {},
    // Reactive state: { [name]: { value, error, valid, ... } }
    init() {
      this.$el.noValidate = !0, this.$el.addEventListener("valid-set", (i) => {
        const { name: o, valid: s, message: n } = i.detail;
        if (o && this.fields[o])
          this.fields[o].valid = s, this.fields[o].error = n || null, s || (this.fields[o].touched = !0);
        else if (i.target && i.target.name && this.fields[i.target.name]) {
          const t = i.target.name;
          this.fields[t].valid = s, this.fields[t].error = n || null, s || (this.fields[t].touched = !0);
        }
      });
    },
    // Reset validation state (API: Uno.valid.clear)
    reset() {
      Object.keys(this.fields).forEach((i) => {
        const o = this.fields[i];
        o.error = null, o.valid = !0, o.touched = !1, o.dirty = !1, o.validating = !1;
      });
    },
    // Register field (called by x-validate directive)
    register(i, o, s = []) {
      if (this.fields[i]) return;
      const n = [...s];
      if ((o.hasAttribute("required") || o.dataset.validRequired === "true") && n.push({ type: "required" }), (o.type === "email" || o.dataset.validType === "email") && n.push("email"), (o.type === "url" || o.dataset.validType === "url") && n.push("url"), o.dataset.validType) {
        const t = o.dataset.validType;
        t !== "email" && t !== "url" && n.push(t);
      }
      o.dataset.validPattern && n.push({ type: "regexp", expression: o.dataset.validPattern }), this.fields[i] = {
        el: o,
        rules: n,
        value: o.value,
        error: null,
        valid: !0,
        dirty: !1,
        touched: !1,
        validating: !1
      };
    },
    update(i, o) {
      this.fields[i] && (this.fields[i].value = o, this.fields[i].dirty = !0, this.validate(i));
    },
    touch(i) {
      this.fields[i] && (this.fields[i].touched = !0, this.validate(i));
    },
    async validate(i) {
      const o = this.fields[i];
      if (!o) return !0;
      o.validating = !0, o.error = null;
      const s = (o.value || "").trim();
      let n = null;
      for (const t of o.rules) {
        const d = t.type || t, l = typeof t == "object" ? t : {};
        if (!s) {
          if (d === "required") {
            n = l.msg || r("required");
            break;
          }
          if (d === "emptyquit") break;
          continue;
        }
        if (Y[d])
          Y[d].test(s) || (n = l.msg || r(d) || r("pattern"));
        else if (d === "regexp" && l.expression) {
          let f = X.get(l);
          if (!f) {
            const m = String(l.expression);
            f = kt(m, { defaultFlags: "i" }), f && X.set(l, f);
          }
          f && f.test(s) || (n = l.msg || r("pattern"));
        } else if (d === "limit")
          l.min && s.length < l.min && (n = l.msg || r("min", { n: l.min })), l.max && s.length > l.max && (n = l.msg || r("max", { n: l.max }));
        else if (d === "range") {
          const f = parseFloat(s);
          l.min !== void 0 && f < l.min && (n = l.msg || `Min ${l.min}`), l.max !== void 0 && f > l.max && (n = l.msg || `Max ${l.max}`);
        } else if (d === "compareto") {
          const f = document.getElementById(l.destid);
          f && s !== f.value && (n = l.msg || "Not match");
        } else if (d === "tag") {
          const f = l.mark || ",", m = s.split(f).filter((u) => u.trim());
          l.max && m.length > l.max && (n = l.msg || `Max ${l.max} items`);
        }
        if (n) break;
        if (d === "callback" && l.fn) {
          const f = e.store("valid").callbacks[l.fn];
          if (f)
            try {
              const m = await f(o.el.id, s);
              m === !1 || typeof m == "object" && m.valid === !1 ? n = (typeof m == "object" ? m.message : "") || "Invalid" : typeof m == "string" && (n = m);
            } catch {
            }
        }
        if (n) break;
      }
      return o.error = n, o.valid = !n, o.validating = !1, o.valid;
    },
    async validateAll() {
      const i = Object.keys(this.fields);
      return (await Promise.all(i.map((s) => this.validate(s)))).every((s) => s);
    },
    async submit(i) {
      if (i && (i.preventDefault(), i.stopPropagation()), Object.keys(this.fields).forEach((s) => {
        this.fields[s].touched = !0;
      }), await this.validateAll())
        this.$el.dispatchEvent(new CustomEvent("valid-submit", { bubbles: !0 })), !this.$el.hasAttribute("hx-post") && !this.$el.hasAttribute("hx-get") && this.$el.submit();
      else {
        const s = Object.keys(this.fields).find((n) => !this.fields[n].valid);
        s && this.fields[s].el && this.fields[s].el.focus();
      }
    }
  })), e.directive("validate", (i, o, { effect: s, cleanup: n }) => {
    const t = i.getAttribute("name");
    if (!t) return;
    const d = () => {
      let h = i.parentNode;
      h.classList.contains("uno-input-wrapper") || (h = document.createElement("div"), h.className = "relative w-full uno-input-wrapper", i.parentNode.insertBefore(h, i), h.appendChild(i));
      let g = h.querySelector(".uno-valid-icon");
      return g || (g = document.createElement("div"), g.className = "uno-valid-icon absolute right-0 top-0 bottom-0 flex items-center justify-center w-10 pointer-events-none transition-opacity duration-200", g.style.opacity = "0", h.appendChild(g)), { wrapper: h, icon: g };
    }, { icon: l } = d();
    let f = [];
    const m = i.getAttribute("valid");
    m && (f = Tt(m)), s(() => {
      const h = e.$data(i);
      h && typeof h.register == "function" && h.register(t, i, f);
    });
    const u = () => {
      const h = e.$data(i);
      if (h && h.update && h.update(t, i.value), f.find((g) => (g.type === "limit" || g === "limit") && g.show === !0)) {
        const g = f.find((y) => y.type === "limit" || y === "limit");
        g && g.max && i.value.length > g.max && (i.value = i.value.substring(0, g.max), h && h.update && h.update(t, i.value));
      }
    }, b = () => {
      const h = e.$data(i);
      h && h.touch && h.touch(t);
    };
    i.addEventListener("input", u), i.addEventListener("blur", b), s(() => {
      const h = e.$data(i);
      if (h && h.fields && h.fields[t]) {
        const g = h.fields[t], y = !!g.error, w = g.valid && g.dirty && !g.error;
        i.classList.toggle("is-invalid", y), i.classList.toggle("is-valid", w), i.classList.toggle("validating", g.validating), l.innerHTML = "", l.style.opacity = "0";
        const T = p.getStyle("valid", "inputError", "!border-red-500"), k = p.getStyle("valid", "inputValid", "!border-green-500"), $ = p.getStyle("valid", "iconError", "text-red-500 dark:text-red-400"), L = p.getStyle("valid", "iconValid", "text-green-500 dark:text-green-400"), H = (U) => U.split(" ").filter(Boolean).forEach((q) => i.classList.add(q)), v = (U) => U.split(" ").filter(Boolean).forEach((q) => i.classList.remove(q));
        y ? (H(T), i.classList.add("pr-10"), v(k), l.innerHTML = a.invalid, l.className = `uno-valid-icon absolute right-0 top-0 bottom-0 flex items-center justify-center w-10 pointer-events-none transition-opacity duration-200 ${$}`, l.style.opacity = "1") : w ? (H(k), i.classList.add("pr-10"), v(T), l.innerHTML = a.valid, l.className = `uno-valid-icon absolute right-0 top-0 bottom-0 flex items-center justify-center w-10 pointer-events-none transition-opacity duration-200 ${L}`, l.style.opacity = "1") : (v(T), v(k), i.classList.remove("pr-10"));
        let E = i.parentNode.nextElementSibling, I = null, P = 0;
        for (; E && P < 3 && !(E.nodeType === 1 && (E.classList.contains("animate-spin") ? I = E : I = E.querySelector(".animate-spin"), I)); )
          E = E.nextElementSibling, P++;
        I && (I.style.display = g.validating ? "block" : "");
      }
    }), n(() => {
      i.removeEventListener("input", u), i.removeEventListener("blur", b);
    });
  }), e.directive("valid-error", (i, { expression: o }, { effect: s }) => {
    const n = o;
    i.style.display = "none", s(() => {
      const t = e.$data(i);
      if (t && t.fields && t.fields[n]) {
        const l = t.fields[n].error;
        if (l) {
          i.textContent = l, i.style.display = "block";
          const f = p.getStyle(
            "valid",
            "helperText",
            "text-red-600 dark:text-red-400 text-xs"
          );
          i.className = f;
        } else
          i.style.display = "none";
      }
    });
  }), e.directive("valid-count", (i, { expression: o }, { effect: s }) => {
    const n = o;
    s(() => {
      const t = e.$data(i);
      if (t && t.fields && t.fields[n]) {
        const d = t.fields[n], l = d.rules.find((f) => f.type === "limit" || f === "limit");
        if (l && l.max) {
          const f = (d.value || "").length, m = l.max - f, u = e.store("i18n");
          i.textContent = u && u.t("charCount", { count: m }) || `Remaining: ${m}`, i.style.display = "block";
        } else
          i.style.display = "none";
      }
    });
  });
};
let D;
const F = () => ({
  activeTab: null,
  _cleanup: [],
  init() {
    if (!A(this.$el, "tabs_init")) return;
    this.$nextTick(() => {
      const a = window.location.hash.slice(1);
      let r = !1;
      if (a) {
        const c = this.$el.querySelector(`[data-tab-hash="${a}"]`);
        c && (this.activeTab = c.getAttribute("data-tab-target"), r = !0);
      }
      if (!r && !this.activeTab) {
        const c = this.$el.querySelector('[role="tab"]');
        c && (this.activeTab = c.getAttribute("data-tab-target"));
      }
    });
    const e = () => {
      const a = window.location.hash.slice(1);
      if (a) {
        const r = this.$el.querySelector(`[data-tab-hash="${a}"]`);
        r && (this.activeTab = r.getAttribute("data-tab-target"));
      }
    };
    window.addEventListener("hashchange", e), this._cleanup.push(() => window.removeEventListener("hashchange", e));
  },
  destroy() {
    this._cleanup.forEach((e) => e()), this._cleanup = [];
  },
  isActive(e) {
    return this.activeTab === e;
  },
  async open(e, a = null, r = null, c = null) {
    if (this.activeTab = e, r && window.location.hash.slice(1) !== r && window.history.replaceState(null, null, `#${r}`), a) {
      const i = document.getElementById(e);
      if (i && i.getAttribute("data-tab-loaded") !== "true")
        try {
          i.innerHTML = p.options.ui?.box?.loading || '<div class="animate-pulse bg-gray-200 h-32 w-full rounded"></div>';
          const o = await fetch(a);
          if (!o.ok) throw new Error("Network response was not ok");
          const s = await o.text();
          z(i, s, { trusted: !1, alpine: D || window.Alpine }), i.setAttribute("data-tab-loaded", "true");
        } catch (o) {
          i.innerHTML = `<div class="text-red-500 p-4">Error loading content: ${o?.message || "Unknown error"}</div>`;
        }
    }
    if (c) {
      const i = document.getElementById(e);
      if (i && i.getAttribute("data-tab-loaded") !== "true")
        try {
          i.innerHTML = p.options.ui?.box?.loading || '<div class="animate-pulse bg-gray-200 h-32 w-full rounded"></div>';
          const o = await fetch(c, { headers: { "HX-Request": "true" } });
          if (!o.ok) throw new Error("Network response was not ok");
          const s = await o.text();
          z(i, s, { trusted: !1, alpine: D || window.Alpine }), i.setAttribute("data-tab-loaded", "true");
        } catch {
          i.innerHTML = '<div class="text-red-500 p-4">Error loading content via HTMX</div>';
        }
    }
  },
  trigger(e, a = null, r = null, c = null) {
    r && typeof window < "u" && window.location.hash.slice(1) === r && (this.activeTab || (this.activeTab = e));
    const i = p.getStyle("tab", "active", "border-b-2 border-blue-500"), o = p.getStyle("tab", "inactive", "text-slate-500 hover:text-slate-700"), s = p.getStyle("tab", "triggerActive", i), n = p.getStyle("tab", "triggerInactive", o);
    return {
      "@click.prevent": `open('${e}', ${a ? `'${a}'` : "null"}, ${r ? `'${r}'` : "null"}, ${c ? `'${c}'` : "null"})`,
      ":class": `isActive('${e}') ? '${s}' : '${n}'`,
      role: "tab",
      ":aria-selected": `isActive('${e}') ? 'true' : 'false'`,
      "data-tab-target": e,
      "data-tab-hash": r || null,
      "data-tab-url": a || null,
      "data-tab-hx-url": c || null
    };
  },
  tablist: {
    role: "tablist",
    "@keydown.right.prevent": "next()",
    "@keydown.left.prevent": "prev()",
    "@keydown.home.prevent": "first()",
    "@keydown.end.prevent": "last()"
  },
  panel(e) {
    return {
      "x-show": `isActive('${e}')`,
      role: "tabpanel",
      id: e,
      tabindex: "0"
    };
  }
});
typeof window < "u" && (window.unoTabs = F);
const $t = (e) => {
  D = e, e.data("unoTabs", F), p.tab = {
    open: (a) => {
      const r = document.querySelector(`[data-tab-target="${a}"]`);
      r ? r.click() : window.location.hash = a;
    }
  };
};
p.dom = N;
typeof window < "u" && (window.Alpine = window.Alpine || O, window.htmx = window.htmx || J, window.unoTabs = window.unoTabs || F, window.Uno = p);
p.use(st).use(ut).use(ft).use(xt).use(Et).use($t);
function j() {
  if (p.initialized) return;
  const e = window.UnoConfig || {};
  C.setDebug(e.debug !== !1), e.debug !== !1 && C.info("Uno: Initializing..."), p.init(window.Alpine, e);
}
if (typeof window < "u") {
  const e = window.deferLoadingAlpine;
  window.deferLoadingAlpine = (a) => {
    j(), typeof e == "function" ? e(a) : a();
  }, document.addEventListener("alpine:init", () => {
    p.initialized || j();
  }), document.addEventListener("DOMContentLoaded", () => {
    !p.initialized && window.Alpine && j();
  }), !p.initialized && window.Alpine && j();
}
export {
  p as default
};
