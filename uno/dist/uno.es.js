import I from "alpinejs";
import V from "htmx.org";
let D = !0, S = null;
const g = {
  /**
   * Set debug mode
   * @param {boolean} debug
   */
  setDebug(n) {
    D = n;
  },
  /**
   * Set custom error reporter
   * @param {Function} reporter - (error, ...args) => void
   */
  setReporter(n) {
    S = n;
  },
  /**
   * Log info message
   * @param {...any} args - Arguments to log
   */
  info(...n) {
    D && console.log("[Uno Info]", ...n);
  },
  /**
   * Log warning message
   * @param {...any} args - Arguments to log
   */
  warn(...n) {
    console.warn("[Uno Warn]", ...n);
  },
  /**
   * Log error message
   * @param {...any} args - Arguments to log
   */
  error(...n) {
    if (console.error("[Uno Error]", ...n), S)
      try {
        S(...n);
      } catch (l) {
        console.error("[Uno Error] Reporter failed:", l);
      }
  }
};
class X {
  constructor() {
    this.version = "0.3.0", this.plugins = [], this.initialized = !1, this.options = {
      debug: !0,
      csp: {
        nonce: null
        // For Content Security Policy
      },
      onError: null
      // Custom error handler
    };
  }
  /**
   * Set global configuration (before or after init)
   * @param {Object} config - Global configuration object
   */
  setConfig(l) {
    const i = { ...this.options };
    return Object.keys(l).forEach((t) => {
      typeof l[t] == "object" && l[t] !== null && !Array.isArray(l[t]) && i[t] ? i[t] = { ...i[t], ...l[t] } : i[t] = l[t];
    }), this.options = i, this.initialized && g.info("Uno configuration updated", this.options), this;
  }
  /**
   * Register a plugin
   * @param {Function} plugin - Plugin function (Alpine, options) => void
   * @param {Object} options - Plugin options
   * @returns {UnoInstance}
   */
  use(l, i = {}) {
    return this.initialized && g.warn("Uno.use() called after initialization. Plugin might not work correctly."), this.plugins.push({ plugin: l, options: i }), this;
  }
  /**
   * Initialize Uno
   * @param {Object|Function} alpineOrOptions - Alpine instance or Configuration options
   * @param {Object} options - Configuration options (if first arg is Alpine)
   * @returns {Object} Alpine instance
   */
  init(l = {}, i = {}) {
    if (this.initialized)
      return g.warn("Uno has already been initialized."), I;
    let t = I, d = i;
    return l && (l.version || typeof l.start == "function" || typeof l.store == "function") ? t = l : d = l || i, this.options = { ...this.options, ...d }, g.setDebug(this.options.debug !== !1), this.options.onError && g.setReporter(this.options.onError), g.info(`Uno v${this.version} initializing...`), this.plugins.forEach(({ plugin: r, options: e }) => {
      try {
        const o = { ...this.options, ...e };
        r(t, o);
      } catch (o) {
        g.error("Failed to load plugin", o);
      }
    }), this.initialized = !0, g.info("Uno initialized"), t;
  }
}
const E = new X(), J = {
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
}, Q = {
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
const F = {
  en: J,
  "zh-CN": K,
  "zh-TW": Q
}, P = Object.keys(F);
function tt() {
  if (typeof localStorage < "u") {
    const n = localStorage.getItem("uno.lang");
    if (n && P.includes(n)) return n;
  }
  if (typeof navigator < "u") {
    const n = navigator.language, l = n.startsWith("zh") ? n.includes("TW") || n.includes("Hant") ? "zh-TW" : "zh-CN" : "en";
    return P.includes(l) ? l : "en";
  }
  return "en";
}
function et(n, l) {
  const i = l.split(".");
  let t = n;
  for (const d of i) {
    if (!t || typeof t != "object" || !(d in t)) return;
    t = t[d];
  }
  return t;
}
function Y(n, l) {
  for (const i in l)
    l[i] && typeof l[i] == "object" ? (n[i] || Object.assign(n, { [i]: {} }), Y(n[i], l[i])) : Object.assign(n, { [i]: l[i] });
}
function nt(n) {
  n.store("i18n", {
    locale: tt(),
    messages: F,
    init() {
      typeof document < "u" && (document.documentElement.lang = this.locale);
    },
    t(l, i = {}) {
      const t = this.locale, d = this.messages[t] || this.messages.en, r = et(d, l);
      return (typeof r == "string" ? r : l).replace(/\{(\w+)\}/g, (o, s) => s in i ? String(i[s]) : `{${s}}`);
    },
    setLang(l) {
      if (!this.messages[l]) {
        console.warn(`[Uno] Language ${l} not found`);
        return;
      }
      this.locale = l, typeof localStorage < "u" && localStorage.setItem("uno.lang", l), typeof document < "u" && (document.documentElement.lang = l);
    },
    extend(l) {
      for (const i in l)
        this.messages[i] ? Y(this.messages[i], l[i]) : this.messages[i] = l[i];
    }
  });
}
const it = (n) => {
  const l = n || "center", i = {
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
  return i[l] || i.center;
}, ot = (n) => {
  const l = n || "top-right", i = {
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
  return i[l] || i["top-right"];
}, st = (n) => {
  const l = n || "top-right";
  let i = "flex-col", t = "items-end", d = "space-y-2";
  return l.includes("bottom") && (i = "flex-col-reverse", d = "space-y-reverse space-y-2"), l.includes("left") ? t = "items-start" : l.includes("center") ? l === "center" ? t = "items-center justify-center" : t = "items-center" : t = "items-end", { flexDir: i, alignItems: t, spaceY: d };
};
function rt(n, l, i) {
  const t = n.getBoundingClientRect(), d = l.getBoundingClientRect(), r = d.width, e = d.height, o = t.width, s = t.height, a = { top: t.top, left: t.left }, c = window.innerWidth, m = window.innerHeight, f = 0, h = 0, u = m, p = c, C = (k, T) => {
    if (T) {
      if (k + r > p)
        return p - r - 10;
      if (k < h)
        return h + 10;
    } else {
      if (k + e > u)
        return u - e - 10;
      if (k < f)
        return f + 10;
    }
    return k;
  };
  let x, v, w = i;
  switch (i) {
    case "left":
      v = a.left - r - 5, v < h && (v = a.left + o + 5, w = "right"), x = C(a.top + s / 2 - e / 2, !1);
      break;
    case "top":
      x = a.top - e - 8, x < f && (x = a.top + s + 8, w = "bottom"), v = a.left + (o - r) / 2, v = C(v, !0);
      break;
    case "right":
      v = a.left + o + 5, v + r > p && (v = a.left - r - 5, w = "left"), x = C(a.top + s / 2 - e / 2, !1);
      break;
    case "bottom":
    default:
      x = a.top + s + 8, x + e > u && (x = a.top - e - 8, w = "top"), v = a.left + (o - r) / 2, v = C(v, !0);
      break;
  }
  return { top: x, left: v, position: w };
}
function B(n) {
  return String(n ?? "").replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;").replace(/'/g, "&#39;");
}
const _ = /* @__PURE__ */ new WeakMap();
function W(n, l) {
  if (!n) return !1;
  const i = String(l);
  let t = _.get(n);
  return t || (t = /* @__PURE__ */ new Set(), _.set(n, t)), t.has(i) ? !1 : (t.add(i), !0);
}
function M(n, l = {}) {
  const {
    allowedTags: i = [
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
    allowedAttrs: t = { a: ["href", "target", "rel"], "*": ["class"] }
  } = l, d = document.createElement("template");
  d.innerHTML = String(n ?? "");
  const r = document.createTreeWalker(d.content, NodeFilter.SHOW_ELEMENT), e = [];
  for (; r.nextNode(); ) {
    const o = r.currentNode, s = o.tagName.toLowerCase();
    if (!["script", "style", "iframe", "object", "embed"].includes(s)) {
      if (!i.includes(s)) {
        e.push(o);
        continue;
      }
      for (const a of [...o.attributes]) {
        const c = a.name.toLowerCase();
        (t[s] || []).concat(t["*"] || []).includes(c) || o.removeAttribute(a.name);
      }
      if (s === "a") {
        const a = (o.getAttribute("href") || "").trim();
        a.startsWith("/") || a.startsWith("https://") || a.startsWith("http://") || a.startsWith("#") || o.removeAttribute("href"), o.getAttribute("rel") || o.setAttribute("rel", "noopener noreferrer");
      }
    }
  }
  for (const o of e)
    o.replaceWith(document.createTextNode(o.textContent || ""));
  return d.innerHTML;
}
function L(n, l, { trusted: i = !1, alpine: t } = {}) {
  n && (n.innerHTML = i ? String(l ?? "") : M(l), t?.initTree?.(n));
}
const j = {
  /**
   * Ensure a container element exists in the DOM.
   * @param {string} id - Element ID
   * @param {string} html - Inner HTML content
   * @param {Object} attributes - Attributes to set
   * @param {HTMLElement} parent - Parent element (default: document.body)
   */
  ensureContainer(n, l = "", i = {}, t = document.body) {
    if (!document.getElementById(n))
      try {
        const d = document.createElement("div");
        return d.id = n, Object.entries(i).forEach(([r, e]) => {
          d.setAttribute(r, e);
        }), l && (d.innerHTML = l), t ? t.appendChild(d) : document.documentElement.appendChild(d), d;
      } catch (d) {
        g.error(`DomManager: Failed to inject container #${n}`, d);
      }
  },
  /**
   * Inject CSS styles with CSP nonce support.
   * @param {string} id - Style element ID
   * @param {string} css - CSS content
   */
  injectStyles(n, l) {
    if (!document.getElementById(n))
      try {
        const i = document.createElement("style");
        i.id = n, E.options.csp?.nonce && i.setAttribute("nonce", E.options.csp.nonce), i.textContent = l, document.head.appendChild(i);
      } catch (i) {
        g.error(`DomManager: Failed to inject styles #${n}`, i);
      }
  }
};
const at = (n, l = {}) => {
  const i = l.tooltip || {}, t = () => typeof window < "u" && window.Uno?.options?.tooltip || i;
  if (!n) {
    g.error("Tooltip: Alpine is undefined. Plugin initialization failed.");
    return;
  }
  n.store("tooltip", {
    visible: !1,
    content: "",
    placement: "top",
    top: -9999,
    // Start off-screen
    left: -9999,
    triggerEl: null,
    onShow: null,
    async show(r, e, o = {}) {
      this.triggerEl = r, this.content = e, this.placement = o.placement || "top", this.onShow = o.onShow, this.visible = !0, this.ensureContainer();
      const s = document.getElementById("uno-tooltip-body");
      s && L(s, this.content, { trusted: o.trustedHTML === !0, alpine: n }), await n.nextTick(), this.updatePosition(), this.onShow && this.onShow(document.getElementById("uno-tooltip-body"));
    },
    hide() {
      this.visible = !1, this.triggerEl = null;
    },
    updatePosition() {
      if (!this.triggerEl) return;
      const r = document.getElementById("uno-tooltip-content");
      if (!r) return;
      const e = rt(this.triggerEl, r, this.placement);
      this.top = e.top, this.left = e.left;
    },
    ensureContainer() {
      d();
    }
  });
  const d = () => {
    j.ensureContainer("uno-tooltip-container", `
          <div id="uno-tooltip-content"
               x-show="$store.tooltip.visible"
               role="tooltip"
               x-transition:enter="transition ease-out duration-200"
               x-transition:enter-start="opacity-0 scale-95"
               x-transition:enter-end="opacity-100 scale-100"
               x-transition:leave="transition ease-in duration-150"
               x-transition:leave-start="opacity-100 scale-100"
               x-transition:leave-end="opacity-0 scale-95"
               class="fixed px-3 py-2 text-sm font-medium rounded-lg shadow-sm whitespace-nowrap pointer-events-none border bg-white text-gray-900 border-gray-200 dark:bg-slate-800 dark:text-white dark:border-slate-700 z-[20000]"
              :style="\`top: \${$store.tooltip.top}px; left: \${$store.tooltip.left}px;\`"
              style="display: none;">
              <div id="uno-tooltip-body" class="pointer-events-auto"></div>
          </div>
        `, {
      "x-data": "{}"
    });
  };
  n.directive(
    "tooltip",
    (r, { expression: e, modifiers: o }, { evaluateLater: s, effect: a, cleanup: c }) => {
      const m = s(e);
      let f = {
        content: "",
        placement: "top",
        trigger: "hover"
      };
      const h = () => {
        let y = t().placement || "top";
        return o.includes("top") && (y = "top"), o.includes("bottom") && (y = "bottom"), o.includes("left") && (y = "left"), o.includes("right") && (y = "right"), y;
      }, u = () => {
        const b = t();
        return o.includes("click") ? "click" : b.trigger || "hover";
      }, p = (b) => {
        if (typeof b == "string" ? f.content = b : typeof b == "object" && (f = { ...f, ...b }), f.content && f.content.startsWith("#") && !f.content.includes(" ")) {
          const y = document.querySelector(f.content);
          y && (f.content = y.innerHTML);
        } else if (f.target) {
          const y = document.querySelector(f.target);
          y && (f.content = y.innerHTML);
        }
      };
      a(() => {
        m((b) => {
          p(b);
        });
      });
      const C = () => {
        const b = n.store("tooltip");
        r.setAttribute("aria-describedby", "uno-tooltip-content"), b.show(r, f.content, { ...f, placement: h() });
      }, x = () => {
        r.removeAttribute("aria-describedby"), n.store("tooltip").hide();
      }, v = () => {
        const b = n.store("tooltip");
        b.visible && b.triggerEl === r ? b.hide() : C();
      }, w = () => {
        u() === "click" && v();
      }, k = () => {
        u() === "hover" && C();
      }, T = () => {
        u() === "hover" && x();
      }, R = () => {
        u() === "hover" && C();
      }, U = () => {
        u() === "hover" && x();
      }, q = (b) => {
        if (u() !== "click") return;
        const y = n.store("tooltip"), O = document.getElementById("uno-tooltip-content");
        y.visible && y.triggerEl === r && !r.contains(b.target) && (!O || !O.contains(b.target)) && y.hide();
      };
      r.addEventListener("click", w), r.addEventListener("mouseenter", k), r.addEventListener("mouseleave", T), r.addEventListener("focus", R), r.addEventListener("blur", U), window.addEventListener("click", q), c(() => {
        r.removeEventListener("click", w), r.removeEventListener("mouseenter", k), r.removeEventListener("mouseleave", T), r.removeEventListener("focus", R), r.removeEventListener("blur", U), window.removeEventListener("click", q);
      });
    }
  );
};
const lt = (n, l = {}) => {
  const i = l.notify || {}, t = () => typeof window < "u" && window.Uno?.options?.notify || i;
  n.store("notify", {
    items: [],
    nextId: 1,
    activeRequests: 0,
    activeProgress: 0,
    // Notification Methods
    show(e, o = {}) {
      const s = t(), {
        type: a = "info",
        timeout: c = s.duration !== void 0 ? s.duration : 3e3,
        position: m = s.position || "top-right",
        closeIndicator: f = !1,
        showCloseButton: h = !0
      } = o, u = this.nextId++, p = {
        id: u,
        text: e,
        type: a,
        position: m,
        closeIndicator: f,
        showCloseButton: h,
        visible: !0
        // For transition
      };
      this.items.unshift(p), c > 0 && setTimeout(() => {
        this.close(u);
      }, c), this.ensureContainer();
    },
    close(e) {
      this.items.find((s) => s.id === e) && (this.items = this.items.filter((s) => s.id !== e));
    },
    success(e, o = {}) {
      const s = typeof o == "number" ? { timeout: o } : o;
      this.show(e, { type: "success", ...s });
    },
    warning(e, o = {}) {
      const s = typeof o == "number" ? { timeout: o } : o;
      this.show(e, { type: "warning", ...s });
    },
    error(e, o = {}) {
      const s = typeof o == "number" ? { timeout: o } : o;
      this.show(e, { type: "error", ...s });
    },
    info(e, o = {}) {
      const s = typeof o == "number" ? { timeout: o } : o;
      this.show(e, { type: "info", ...s });
    },
    // Loading Indicator Methods
    loading(e) {
      e ? this.activeRequests++ : this.activeRequests = Math.max(0, this.activeRequests - 1), this.ensureContainer();
    },
    // Progress Bar Methods
    progress(e) {
      e ? this.activeProgress++ : this.activeProgress = Math.max(0, this.activeProgress - 1), this.ensureContainer();
    },
    // Container Injection
    ensureContainer() {
      d();
    }
  });
  const d = () => {
    j.injectStyles(
      "uno-notify-styles",
      `
      @keyframes uno-progress-indeterminate {
        0% { left: -35%; width: 35%; }
        100% { left: 100%; width: 35%; }
      }
      .uno-progress-anim {
        animation: uno-progress-indeterminate 1.5s infinite linear;
      }
    `
    );
    const s = `
      <!-- Progress Bar (Top Fixed) -->
      <div id="uno-progress-bar"
           x-show="$store.notify.activeProgress > 0"
           class="fixed top-0 left-0 w-full h-1 bg-blue-100 z-[30200]"
           style="display: none;">
           <div class="absolute top-0 h-full bg-blue-600 uno-progress-anim"></div>
      </div>

      <!-- Loading Indicator (Top Right Fixed) -->
      <div id="uno-loading-indicator" 
           x-show="$store.notify.activeRequests > 0"
           class="fixed top-4 right-4 z-[30150] bg-blue-100 border-blue-500 text-blue-700 p-2 rounded-full shadow-md flex items-center dark:bg-blue-900 dark:text-blue-100"
           style="display: none;">
           <svg class="animate-spin -ml-1 mr-2 h-5 w-5 text-blue-500 dark:text-blue-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
           <span>Loading</span>
      </div>

      ${[
      "top-right",
      "top-left",
      "top-center",
      "center",
      "bottom-right",
      "bottom-left",
      "bottom-center"
    ].map((a) => {
      const c = ot(a), { flexDir: m, alignItems: f, spaceY: h } = st(a);
      let u = `i.position === '${a}'`;
      return a === "top-center" && (u += " || i.position === 'top'"), a === "bottom-center" && (u += " || i.position === 'bottom'"), `
      <!-- ${a} -->
      <div class="fixed ${c} z-[30000] p-4 flex ${m} ${h} ${f} pointer-events-none">
          <template x-for="item in $store.notify.items.filter(i => ${u})" :key="item.id">
              ${r()}
          </template>
      </div>`;
    }).join(`
`)}
    `;
    j.ensureContainer("uno-notify-container", s, {
      "x-data": "{}"
    });
  }, r = () => `
    <div x-show="true"
         data-testid="notify-item"
         x-transition:enter="transition ease-out duration-300"
         x-transition:enter-start="opacity-0 transform scale-90 translate-y-2"
         x-transition:enter-end="opacity-100 transform scale-100 translate-y-0"
         x-transition:leave="transition ease-in duration-300"
         x-transition:leave-start="opacity-100 transform scale-100 translate-y-0"
         x-transition:leave-end="opacity-0 transform scale-90 translate-y-2"
         class="pointer-events-auto relative w-80 max-w-full border-l-4 pl-3 pr-3 py-4 rounded shadow-sm dark:shadow-md overflow-hidden"
         :class="{
            'bg-green-100 dark:bg-green-800 border-green-500 text-green-700 dark:text-green-100': item.type === 'success',
            'bg-yellow-100 dark:bg-yellow-800 border-yellow-500 text-yellow-700 dark:text-yellow-100': item.type === 'warning',
            'bg-red-100 dark:bg-red-800 border-red-500 text-red-700 dark:text-red-100': item.type === 'error',
            'bg-blue-100 dark:bg-blue-800 border-blue-500 text-blue-700 dark:text-blue-100': item.type === 'info'
         }"
         role="alert">
         
         <div class="flex items-start justify-between">
            <div x-text="item.text" class="text-sm font-medium"></div>
            <button x-show="item.showCloseButton" 
                    @click="$store.notify.close(item.id)" 
                    class="ml-4 opacity-50 hover:opacity-100 transition-opacity focus:outline-none">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
            </button>
         </div>
    </div>
    `;
}, $ = /* @__PURE__ */ new Map();
function dt(n) {
  const l = String(n ?? "").trim();
  if (!l) return [];
  try {
    const i = JSON.parse(l);
    if (typeof i == "string") return [[i, {}]];
    if (i && typeof i == "object") return Object.entries(i);
  } catch {
  }
  return l.split(",").map((i) => i.trim()).filter(Boolean).map((i) => [i, {}]);
}
function ct(n, l, i) {
  const t = n || document.body || document.documentElement;
  if (t) {
    if (typeof window < "u" && window.htmx?.trigger) {
      window.htmx.trigger(t, l, i);
      return;
    }
    t.dispatchEvent(new CustomEvent(l, { bubbles: !0, detail: i }));
  }
}
function ut(n, l) {
  const i = n.headers.get("HX-Trigger");
  if (i)
    for (const [t, d] of dt(i))
      ct(l, t, d);
}
async function ft(n, l = {}) {
  const {
    method: i = "GET",
    htmxTarget: t,
    swap: d = "innerHTML",
    htmxConfig: r = {},
    cache: e = !1
  } = l;
  if (e && i === "GET" && $.has(n))
    return $.get(n);
  if (typeof window < "u" && window.htmx?.ajax) {
    const o = t || document.createElement("div");
    t || (o.style.display = "none"), t || (document.body ? document.body.appendChild(o) : document.documentElement.appendChild(o));
    try {
      await window.htmx.ajax(i, n, { target: o, swap: d, ...r });
      const s = o.innerHTML;
      return e && i === "GET" && $.set(n, s), s;
    } finally {
      t || o.remove();
    }
  }
  return G(n, l);
}
async function Z(n, l = {}) {
  const {
    method: i = "GET",
    headers: t,
    body: d,
    credentials: r = "same-origin",
    triggerTarget: e = typeof document < "u" ? document.body : null,
    cache: o = !1,
    signal: s
    // AbortSignal
  } = l;
  if (o && i === "GET" && $.has(n))
    return $.get(n);
  const a = await fetch(n, { method: i, headers: t, body: d, credentials: r, signal: s });
  if (ut(a, e), !a.ok) {
    const f = new Error(`Request failed: ${a.status}`);
    throw f.status = a.status, f;
  }
  const c = a.headers.get("content-type") || "";
  let m;
  return c.includes("application/json") ? m = await a.json() : m = await a.text(), o && i === "GET" && $.set(n, m), m;
}
async function G(n, l = {}) {
  const i = await Z(n, l);
  return typeof i == "string" ? i : JSON.stringify(i, null, 2);
}
const mt = [
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
function pt(n, l = {}) {
  const i = () => Array.from(n.querySelectorAll(mt)).filter((e) => e.offsetParent !== null), t = () => i()[0], d = () => {
    const e = i();
    return e[e.length - 1];
  }, r = (e) => {
    if (e.key === "Tab") {
      if (e.shiftKey) {
        if (document.activeElement === t()) {
          e.preventDefault();
          const o = d();
          o && o.focus();
        }
      } else if (document.activeElement === d()) {
        e.preventDefault();
        const o = t();
        o && o.focus();
      }
    }
  };
  return setTimeout(() => {
    const e = t();
    e && e.focus();
  }, 50), n.addEventListener("keydown", r), () => {
    n.removeEventListener("keydown", r);
  };
}
const ht = (n, l = {}) => {
  const i = l.box || {}, t = () => typeof window < "u" && window.Uno?.options?.box || i;
  n.directive("uno-trap", (r, { expression: e }, { cleanup: o }) => {
    const s = pt(r);
    o(s);
  }), n.store("box", {
    items: [],
    // Public API
    alert(r, e = {}) {
      const o = typeof e == "string" ? { title: e } : e;
      return this.open("dialog", {
        content: r,
        showOkButton: !0,
        ...o
      });
    },
    confirm(r, e = {}) {
      return new Promise((o) => {
        this.open("confirm", {
          content: r,
          onOk: () => o(!0),
          onCancel: () => o(!1),
          onClose: () => o(!1),
          // Closing without choice = cancel
          ...e
        });
      });
    },
    prompt(r, e = "", o = {}) {
      return new Promise((s) => {
        const a = "prompt-" + Math.random().toString(36).substr(2, 9), c = `
          <div class="space-y-3">
            <p>${B(r)}</p>
            <input id="${a}" type="text" class="w-full px-3 py-2 border rounded focus:ring-2 focus:ring-blue-500 outline-none" value="${B(
          e
        )}">
          </div>
        `;
        this.open("confirm", {
          title: o.title || "Prompt",
          content: c,
          trustedHTML: !0,
          onOk: () => {
            const m = document.getElementById(a);
            s(m ? m.value : null);
          },
          onCancel: () => s(null),
          ...o
        });
      });
    },
    async open(r, e = {}) {
      let o = "dialog", s = e;
      typeof r == "object" ? (s = r, o = s.type || "dialog") : o = r;
      const a = s.boxId || Math.random().toString(36).slice(2, 10), c = {
        id: a,
        type: o,
        title: s.title || "",
        content: s.content || "",
        trustedHTML: s.trustedHTML === !0,
        width: s.width !== void 0 ? s.width : s.widthClass ? void 0 : t().defaultWidth || 400,
        height: s.height,
        widthClass: s.widthClass,
        heightClass: s.heightClass,
        containerClass: s.containerClass,
        modal: s.modal !== !1,
        // default true
        overlayOpacity: s.overlayOpacity ?? 0.5,
        position: s.position || "center",
        showCloseButton: s.showCloseButton !== void 0 ? s.showCloseButton : !0,
        // Buttons
        okText: s.okButtonText,
        cancelText: s.cancelButtonText,
        onOk: s.onOk,
        onCancel: s.onCancel,
        onClose: s.onClose,
        onBeforeClose: s.onBeforeClose,
        onContentLoad: s.onContentLoad,
        showOk: o === "confirm" || o === "message" || o === "dialog" && s.showOkButton,
        showCancel: o === "confirm" || o === "dialog" && s.showCancelButton,
        // Cache Option (P1-8)
        cache: s.cache === !0,
        // State
        loading: !1,
        visible: !0,
        abortController: null
      };
      c.trustedHTML || (c.content = M(c.content), c.trustedHTML = !0);
      const m = n.store("i18n");
      if (c.okText || (c.okText = m ? m.t("box.ok") : "OK"), c.cancelText || (c.cancelText = m ? m.t("box.cancel") : "Cancel"), s.targetId) {
        const f = document.getElementById(s.targetId);
        f ? (c.content = f.innerHTML, s.trustedHTML || (c.content = M(c.content)), c.trustedHTML = !0) : (g.warn("Box: Target element not found:", s.targetId), c.content = '<div class="p-4 text-red-500">Error: Content target not found</div>', c.trustedHTML = !0);
      } else if (s.hxUrl) {
        c.loading = !0, this.add(c);
        try {
          const f = await ft(s.hxUrl, { method: "GET", cache: c.cache }), h = this.items.find((u) => u.id === a);
          h && (h.content = s.trustedHTML ? f : M(f), h.trustedHTML = !0, h.loading = !1, h.onContentLoad && n.nextTick(() => {
            h.onContentLoad(document.getElementById("uno-box-content" + a));
          }));
        } catch (f) {
          g.error("Box: HTMX Load Failed", f);
          const h = this.items.find((u) => u.id === a);
          h && (h.loading = !1, h.trustedHTML = !0, h.content = '<div class="p-4 text-red-500">Error loading content via HTMX</div>');
        }
        return a;
      } else if (s.ajaxUrl) {
        c.loading = !0, c.abortController = new AbortController(), this.add(c);
        try {
          const f = await Z(s.ajaxUrl, {
            cache: c.cache,
            signal: c.abortController.signal
          });
          let h = typeof f == "object" ? `<pre class="bg-gray-100 p-3 rounded overflow-x-auto whitespace-pre-wrap">${JSON.stringify(f, null, 2)}</pre>` : String(f);
          s.trustedHTML || (h = M(h));
          const u = this.items.find((p) => p.id === a);
          u && (u.content = h, u.trustedHTML = !0, u.loading = !1, n.nextTick(() => {
            const p = document.getElementById("uno-box-content" + a);
            p && u.onContentLoad && u.onContentLoad(p);
          }));
        } catch (f) {
          if (f.name === "AbortError") {
            g.info("Box: Request cancelled");
            return;
          }
          g.error("Box: Ajax failed", f);
          const h = this.items.find((u) => u.id === a);
          h && (h.content = '<div class="text-red-500">Error loading content</div>', h.trustedHTML = !0, h.loading = !1);
        }
        return a;
      } else s.iframeUrl && (c.content = `<iframe frameborder="0" class="w-full h-full" src="${B(s.iframeUrl)}"></iframe>`, c.trustedHTML = !0);
      return this.add(c), a;
    },
    add(r) {
      this.items.push(r), this.ensureContainer();
    },
    close(r) {
      const e = this.items.find((o) => o.id === r);
      e && (e.abortController && e.abortController.abort(), e.onBeforeClose && e.onBeforeClose(), this.items = this.items.filter((o) => o.id !== r), e.onClose && e.onClose());
    },
    closeAll() {
      this.items.forEach((r) => {
        r.abortController && r.abortController.abort();
      }), this.items = [];
    },
    getPositionClass(r) {
      const e = "fixed inset-0 flex transition-opacity duration-300 p-4 pointer-events-none", o = it(r);
      return `${e} ${o}`;
    },
    ensureContainer() {
      d();
    }
  });
  const d = () => {
    j.ensureContainer("uno-box-container", `
      <template x-for="(item, index) in $store.box.items" :key="item.id">
          <div :id="'uno-box-outer-container' + item.id"
               :class="$store.box.getPositionClass(item.position)"
               :style="\`z-index: \${10000 + index}\`"
               role="dialog"
               aria-modal="true"
               @keydown.escape.window="if(item.modal) $store.box.close(item.id)">
               
               <!-- Overlay -->
               <div :id="'uno-box-overlay' + item.id"
                    class="absolute inset-0 bg-black transition-opacity duration-300 pointer-events-auto"
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
                    class="relative bg-white dark:bg-gray-900 rounded-lg shadow-2xl ring-1 ring-black/5 overflow-hidden flex flex-col max-h-[90vh] transition-all duration-300 transform pointer-events-auto outline-none"
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
                                class="px-4 py-2 bg-slate-100 dark:bg-slate-700 border border-slate-200 dark:border-slate-600 rounded-md text-sm font-medium text-slate-800 dark:text-slate-200 hover:bg-slate-200 dark:hover:bg-slate-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors">
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
                                class="px-4 py-2 bg-blue-600 border border-transparent rounded-md text-sm font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
                        </button>
                    </div>
               </div>
          </div>
      </template>
    `, {
      "x-data": "{}"
    });
  };
  W(document, "uno:box:delegator") && document.addEventListener("click", (r) => {
    const e = r.target.closest("[data-box-type]");
    if (!e) return;
    const o = e.getAttribute("data-box-type");
    if (!["dialog", "confirm", "message"].includes(o)) return;
    r.preventDefault();
    const s = e.dataset, a = {
      title: s.boxTitle,
      content: s.boxContent,
      width: s.boxWidth ? parseInt(s.boxWidth) : void 0,
      height: s.boxHeight ? parseInt(s.boxHeight) : void 0,
      ajaxUrl: s.boxAjaxUrl,
      hxUrl: s.boxHxUrl,
      iframeUrl: s.boxIframeUrl,
      targetId: s.boxTargetId,
      position: s.boxPosition,
      okButtonText: s.boxOkText,
      cancelButtonText: s.boxCancelText,
      overlayOpacity: s.boxOverlayOpacity ? parseFloat(s.boxOverlayOpacity) : void 0,
      showCloseButton: s.boxShowClose !== "false",
      boxId: s.boxId,
      cache: s.boxCache === "true",
      trustedHTML: s.boxTrusted === "true"
    };
    s.boxModal === "false" && (a.modal = !1), n.store("box").open(o, a);
  }), typeof window < "u" && W(document, "uno:box:htmx-afterSwap") && document.addEventListener("htmx:afterSwap", (r) => {
    r.detail.elt && n && n.initTree(r.detail.elt);
  });
}, N = {
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
function gt(n) {
  const l = n.match(/^\/(.*?)\/([gimsuy]*)$/);
  return l ? { pattern: l[1], flags: l[2] } : { pattern: n, flags: "" };
}
function bt(n, l = {}) {
  const { defaultFlags: i = "i", maxPatternLength: t = 200 } = l, { pattern: d, flags: r } = gt(String(n ?? ""));
  if (!d || d.length > t) return null;
  const e = String(r || i || "");
  if (!/^[gimsuy]*$/.test(e)) return null;
  try {
    return new RegExp(d, e);
  } catch {
    return null;
  }
}
function vt(n) {
  try {
    return n ? JSON.parse(n) : [];
  } catch {
    return [];
  }
}
const wt = (n) => {
  const l = {
    invalid: "<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 12 12' fill='none' stroke='currentColor' class='w-4 h-4'><circle cx='6' cy='6' r='4.5'/><path stroke-linecap='round' d='M5.8 3.6h.4L6 6.5zM6 8.2a.6.6 0 110-1.2.6.6 0 010 1.2z'/></svg>",
    valid: "<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 8 8' class='w-4 h-4'><path fill='currentColor' d='M2.3 6.73L.6 4.53c-.4-1.04.46-1.4 1.1-.8l1.1 1.4 3.4-3.8c.6-.63 1.6-.27 1.2.7l-4 4.6c-.43.5-.8.4-1.1.1z'/></svg>"
  }, i = (t, d) => {
    const r = n.store("i18n");
    if (!r) return t;
    const e = {
      required: "Required",
      pattern: "Invalid format",
      min: `At least ${d?.n} characters`,
      max: `Max ${d?.n} characters`,
      int: "Integer required",
      email: "Invalid email address",
      url: "Invalid URL"
    };
    return r.t(`valid.${t}`, d) || e[t] || t;
  };
  n.store("valid", {
    callbacks: {},
    registerCallback(t, d) {
      typeof d == "function" && (this.callbacks[t] = d);
    },
    async validateForm(t) {
      if (!t) return !1;
      if (t._x_dataStack) {
        const d = t._x_dataStack.find((r) => r.validateAll);
        if (d)
          return await d.validateAll();
      }
      return !0;
    }
  }), n.data("validForm", () => ({
    fields: {},
    // Reactive state: { [name]: { value, error, valid, ... } }
    init() {
      this.$el.noValidate = !0;
    },
    // Register field (called by x-validate directive)
    register(t, d, r = []) {
      if (this.fields[t]) return;
      const e = [...r];
      if ((d.hasAttribute("required") || d.dataset.validRequired === "true") && e.push({ type: "required" }), (d.type === "email" || d.dataset.validType === "email") && e.push("email"), (d.type === "url" || d.dataset.validType === "url") && e.push("url"), d.dataset.validType) {
        const o = d.dataset.validType;
        o !== "email" && o !== "url" && e.push(o);
      }
      d.dataset.validPattern && e.push({ type: "regexp", expression: d.dataset.validPattern }), this.fields[t] = {
        el: d,
        rules: e,
        value: d.value,
        error: null,
        valid: !0,
        dirty: !1,
        touched: !1,
        validating: !1
      };
    },
    update(t, d) {
      this.fields[t] && (this.fields[t].value = d, this.fields[t].dirty = !0, this.validate(t));
    },
    touch(t) {
      this.fields[t] && (this.fields[t].touched = !0, this.validate(t));
    },
    async validate(t) {
      const d = this.fields[t];
      if (!d) return !0;
      d.validating = !0, d.error = null;
      const r = (d.value || "").trim();
      let e = null;
      for (const o of d.rules) {
        const s = o.type || o, a = typeof o == "object" ? o : {};
        if (!r) {
          if (s === "required") {
            e = a.msg || i("required");
            break;
          }
          if (s === "emptyquit") break;
          continue;
        }
        if (N[s])
          N[s].test(r) || (e = a.msg || i(s) || i("pattern"));
        else if (s === "regexp" && a.expression) {
          const c = String(a.expression);
          a.__reSrc !== c && (a.__re = bt(c, { defaultFlags: "i" }), a.__reSrc = c);
          const m = a.__re;
          m && m.test(r) || (e = a.msg || i("pattern"));
        } else if (s === "limit")
          a.min && r.length < a.min && (e = a.msg || i("min", { n: a.min })), a.max && r.length > a.max && (e = a.msg || i("max", { n: a.max }));
        else if (s === "range") {
          const c = parseFloat(r);
          a.min !== void 0 && c < a.min && (e = a.msg || `Min ${a.min}`), a.max !== void 0 && c > a.max && (e = a.msg || `Max ${a.max}`);
        } else if (s === "compareto") {
          const c = document.getElementById(a.destid);
          c && r !== c.value && (e = a.msg || "Not match");
        } else if (s === "tag") {
          const c = a.mark || ",", m = r.split(c).filter((f) => f.trim());
          a.max && m.length > a.max && (e = a.msg || `Max ${a.max} items`);
        }
        if (e) break;
        if (s === "callback" && a.fn) {
          const c = n.store("valid").callbacks[a.fn];
          if (c)
            try {
              const m = await c(d.el.id, r);
              m === !1 || typeof m == "object" && m.valid === !1 ? e = (typeof m == "object" ? m.message : "") || "Invalid" : typeof m == "string" && (e = m);
            } catch {
            }
        }
        if (e) break;
      }
      return d.error = e, d.valid = !e, d.validating = !1, d.valid;
    },
    async validateAll() {
      const t = Object.keys(this.fields);
      return (await Promise.all(t.map((r) => this.validate(r)))).every((r) => r);
    },
    async submit(t) {
      if (t && (t.preventDefault(), t.stopPropagation()), Object.keys(this.fields).forEach((r) => {
        this.fields[r].touched = !0;
      }), await this.validateAll())
        this.$el.dispatchEvent(new CustomEvent("valid-submit", { bubbles: !0 })), !this.$el.hasAttribute("hx-post") && !this.$el.hasAttribute("hx-get") && this.$el.submit();
      else {
        const r = Object.keys(this.fields).find((e) => !this.fields[e].valid);
        r && this.fields[r].el && this.fields[r].el.focus();
      }
    }
  })), n.directive("validate", (t, d, { effect: r, cleanup: e }) => {
    const o = t.getAttribute("name");
    if (!o) return;
    const s = () => {
      let u = t.parentNode;
      u.classList.contains("uno-input-wrapper") || (u = document.createElement("div"), u.className = "relative w-full uno-input-wrapper", t.parentNode.insertBefore(u, t), u.appendChild(t));
      let p = u.querySelector(".uno-valid-icon");
      return p || (p = document.createElement("div"), p.className = "uno-valid-icon absolute right-0 top-0 bottom-0 flex items-center justify-center w-10 pointer-events-none transition-opacity duration-200", p.style.opacity = "0", u.appendChild(p)), { wrapper: u, icon: p };
    }, { icon: a } = s();
    let c = [];
    const m = t.getAttribute("valid");
    m && (c = vt(m)), r(() => {
      const u = n.$data(t);
      u && typeof u.register == "function" && u.register(o, t, c);
    });
    const f = () => {
      const u = n.$data(t);
      if (u && u.update && u.update(o, t.value), c.find((p) => (p.type === "limit" || p === "limit") && p.show === !0)) {
        const p = c.find((C) => C.type === "limit" || C === "limit");
        p && p.max && t.value.length > p.max && (t.value = t.value.substring(0, p.max), u && u.update && u.update(o, t.value));
      }
    }, h = () => {
      const u = n.$data(t);
      u && u.touch && u.touch(o);
    };
    t.addEventListener("input", f), t.addEventListener("blur", h), r(() => {
      const u = n.$data(t);
      if (u && u.fields && u.fields[o]) {
        const p = u.fields[o], C = !!p.error, x = p.valid && p.dirty && !p.error;
        t.classList.toggle("is-invalid", C), t.classList.toggle("is-valid", x), t.classList.toggle("validating", p.validating), a.innerHTML = "", a.style.opacity = "0", C ? (t.classList.add("!border-red-500", "pr-10"), t.classList.remove("!border-green-500"), a.innerHTML = l.invalid, a.className = "uno-valid-icon absolute right-0 top-0 bottom-0 flex items-center justify-center w-10 pointer-events-none transition-opacity duration-200 text-red-500 dark:text-red-400", a.style.opacity = "1") : x ? (t.classList.add("!border-green-500", "pr-10"), t.classList.remove("!border-red-500"), a.innerHTML = l.valid, a.className = "uno-valid-icon absolute right-0 top-0 bottom-0 flex items-center justify-center w-10 pointer-events-none transition-opacity duration-200 text-green-500 dark:text-green-400", a.style.opacity = "1") : t.classList.remove("!border-red-500", "!border-green-500", "pr-10");
        let w = t.parentNode.nextElementSibling, k = null, T = 0;
        for (; w && T < 3 && !(w.nodeType === 1 && (w.classList.contains("animate-spin") ? k = w : k = w.querySelector(".animate-spin"), k)); )
          w = w.nextElementSibling, T++;
        k && (k.style.display = p.validating ? "block" : "");
      }
    }), e(() => {
      t.removeEventListener("input", f), t.removeEventListener("blur", h);
    });
  }), n.directive("valid-error", (t, { expression: d }, { effect: r }) => {
    const e = d;
    t.style.display = "none", r(() => {
      const o = n.$data(t);
      if (o && o.fields && o.fields[e]) {
        const a = o.fields[e].error;
        a ? (t.textContent = a, t.style.display = "block", t.classList.add("text-red-600", "dark:text-red-400", "text-xs"), t.classList.remove("text-gray-400")) : t.style.display = "none";
      }
    });
  }), n.directive("valid-count", (t, { expression: d }, { effect: r }) => {
    const e = d;
    r(() => {
      const o = n.$data(t);
      if (o && o.fields && o.fields[e]) {
        const s = o.fields[e], a = s.rules.find((c) => c.type === "limit" || c === "limit");
        if (a && a.max) {
          const c = (s.value || "").length, m = a.max - c, f = n.store("i18n");
          t.textContent = f && f.t("charCount", { count: m }) || `Remaining: ${m}`, t.style.display = "block";
        } else
          t.style.display = "none";
      }
    });
  });
}, A = (n = "", l = {}) => ({
  active: n,
  options: {
    hash: !0,
    manual: !1,
    // if true, arrow keys only move focus, don't select
    ...l
  },
  loaded: {},
  // Track loaded AJAX content
  init() {
    if (this.options && this.options.hash) {
      const i = () => {
        const t = window.location.hash.substring(1);
        if (t) {
          const d = this.$el.querySelector(`[data-tab-hash="${t}"]`);
          if (d) {
            const r = d.getAttribute("data-tab-target");
            r && (this.active = r);
          }
        }
      };
      this.$nextTick(() => i()), this._hashHandler = () => i(), window.addEventListener("hashchange", this._hashHandler);
    }
  },
  destroy() {
    this._hashHandler && window.removeEventListener("hashchange", this._hashHandler);
  },
  isActive(i) {
    return this.active === i;
  },
  async open(i, t = null, d = null, r = null) {
    if (this.active = i, this.options.hash && d && window.location.hash.substring(1) !== d && history.pushState(null, null, "#" + d), r) {
      if (typeof window < "u" && window.htmx) {
        const e = document.getElementById(i);
        if (e && !this.loaded[i]) {
          L(
            e,
            `<div class="animate-pulse p-4 space-y-2">
               <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-3/4"></div>
               <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/2"></div>
             </div>`,
            { trusted: !0, alpine: window.Alpine || window.top?.Alpine }
          );
          try {
            await window.htmx.ajax("GET", r, {
              target: `#${i}`,
              swap: "innerHTML"
            }), e.setAttribute("data-tab-loaded", "true"), this.loaded[i] = !0, e.dispatchEvent(new CustomEvent("tab-loaded", { bubbles: !0 }));
          } catch {
            L(e, '<div class="text-red-500 p-4">Error loading content via HTMX</div>', {
              trusted: !0
            });
          }
        }
      } else {
        console.warn("Uno Tab: HTMX requested but window.htmx is not defined.");
        const e = document.getElementById(i);
        e && L(
          e,
          `<div class="text-amber-600 bg-amber-50 p-4 rounded border border-amber-200">
            <strong>Warning:</strong> HTMX library not found. Please include htmx.org script.
          </div>`,
          { trusted: !0 }
        );
      }
      return;
    }
    if (t) {
      const e = document.getElementById(i);
      if (e && !this.loaded[i]) {
        L(
          e,
          `<div class="animate-pulse p-4 space-y-2">
             <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-3/4"></div>
             <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/2"></div>
           </div>`,
          { trusted: !0, alpine: window.Alpine || window.top?.Alpine }
        );
        try {
          const o = await G(t, { method: "GET" });
          L(e, o, { trusted: !0, alpine: window.Alpine || window.top?.Alpine }), e.setAttribute("data-tab-loaded", "true"), this.loaded[i] = !0, e.dispatchEvent(new CustomEvent("tab-loaded", { bubbles: !0 }));
        } catch (o) {
          L(
            e,
            `<div class="text-red-500 p-4">Error loading content: ${o?.message || "Unknown error"}</div>`,
            { trusted: !0 }
          );
        }
      }
    }
  },
  trigger(i, t = null, d = null, r = null) {
    return {
      "@click.prevent": () => this.open(i, t, d, r),
      ":class": () => this.isActive(i) ? "current" : "",
      role: "tab",
      "aria-selected": () => this.isActive(i).toString(),
      "data-tab-target": i,
      "data-tab-hash": d || null,
      "data-tab-url": t || null,
      "data-tab-hx-url": r || null
    };
  },
  tablist: {
    role: "tablist",
    "@keydown.right.prevent": "next()",
    "@keydown.left.prevent": "prev()",
    "@keydown.home.prevent": "first()",
    "@keydown.end.prevent": "last()"
  },
  panel(i) {
    return {
      "x-show": () => this.isActive(i),
      role: "tabpanel",
      id: i
    };
  }
});
typeof window < "u" && (window.unoTabs = A);
const xt = (n) => {
  n.data("unoTabs", A);
};
typeof window < "u" && (window.Alpine = window.Alpine || I, window.htmx = window.htmx || V, window.unoTabs = window.unoTabs || A, window.Uno = E);
E.use(nt).use(at).use(lt).use(ht).use(wt).use(xt);
const z = () => {
  if (typeof window > "u") return {};
  const n = window.UnoConfig;
  return n && typeof n == "object" ? n : {};
}, H = () => {
  if (E.initialized) return;
  const n = z();
  g.setDebug(n.debug !== !1), E.init(typeof window < "u" ? window.Alpine : void 0, n);
};
if (typeof window < "u") {
  const n = window.deferLoadingAlpine;
  if (window.deferLoadingAlpine = (l) => {
    H(), typeof n == "function" ? n(l) : l();
  }, document.addEventListener("alpine:init", () => {
    if (!E.initialized) {
      const l = z();
      g.setDebug(l.debug !== !1), l.debug !== !1 && g.info("Uno: alpine:init detected, initializing..."), H();
    }
  }), document.addEventListener("DOMContentLoaded", () => {
    !E.initialized && window.Alpine && H();
  }), !E.initialized && window.Alpine) {
    const l = z();
    g.setDebug(l.debug !== !1), g.info("Uno: Auto-initializing synchronously"), H();
  }
}
export {
  E as default
};
