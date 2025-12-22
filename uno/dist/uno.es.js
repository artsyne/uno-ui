import X from "alpinejs";
import gt from "htmx.org";
let P = !1, K = null;
const L = {
  /**
   * Set debug mode
   * @param {boolean} debug
   */
  setDebug(n) {
    P = n;
  },
  /**
   * Set custom error reporter
   * @param {Function} reporter - (error, ...args) => void
   */
  setReporter(n) {
    K = n;
  },
  /**
   * Log info message
   * @param {...any} args - Arguments to log
   */
  info(...n) {
    P && typeof console < "u" && typeof console.log == "function" && console.log("[Uno Info]", ...n);
  },
  /**
   * Log warning message
   * @param {...any} args - Arguments to log
   */
  warn(...n) {
    P && typeof console < "u" && typeof console.warn == "function" && console.warn("[Uno Warn]", ...n);
  },
  /**
   * Log error message
   * @param {...any} args - Arguments to log
   */
  error(...n) {
    if (K)
      try {
        K(...n);
      } catch (i) {
        P && typeof console < "u" && typeof console.error == "function" && console.error("[Uno Error] Reporter failed:", i);
      }
    P && typeof console < "u" && typeof console.error == "function" && console.error("[Uno Error]", ...n);
  }
};
function U(n) {
  return Object.prototype.toString.call(n) === "[object Object]";
}
function V(n, i) {
  if (!U(n)) return U(i) ? { ...i } : i;
  if (!U(i)) return i === void 0 ? { ...n } : i;
  const o = { ...n };
  for (const [l, e] of Object.entries(i)) {
    if (U(e) && U(o[l])) {
      o[l] = V(o[l], e);
      continue;
    }
    if (U(e)) {
      o[l] = V({}, e);
      continue;
    }
    o[l] = e;
  }
  return o;
}
class mt {
  constructor() {
    this.version = "0.3.0", this.plugins = [], this.initialized = !1, this.options = {
      debug: !1,
      csp: {
        nonce: null
        // For Content Security Policy
      },
      onError: null,
      // Custom error handler
      plugins: {},
      ui: {
        base: { rounded: "rounded-lg", primary: "bg-blue-600", shadow: "shadow-md" },
        layers: { boxBase: 1e4, tooltipBase: 2e4, notifyBase: 3e4 },
        box: {},
        notify: {},
        tooltip: {},
        tab: {},
        valid: {}
      }
    };
  }
  _getUiValue(i, o) {
    const l = this.options.ui?.[i];
    if (!l) return;
    if (typeof o != "string" || o.indexOf(".") === -1) return l?.[o];
    const e = o.split(".").filter(Boolean);
    let r = l;
    for (const u of e) {
      if (r == null) return;
      r = r[u];
    }
    return r;
  }
  /**
   * Get style class from configuration or fallback
   * Support "+" prefix to append class to default
   * @param {string} comp - Component name
   * @param {string} elem - Element name
   * @param {string} defaultCls - Default class string
   * @returns {string} Final class string
   */
  getStyle(i, o, l) {
    const e = this._getUiValue(i, o);
    return e === void 0 ? l : typeof e == "string" && e.trim().indexOf("+") === 0 ? `${l} ${e.trim().slice(1)}`.trim() : e;
  }
  /**
   * Set global configuration (before or after init)
   * @param {Object} config - Global configuration object
   */
  setConfig(i) {
    return this.options = V(this.options, i), this.initialized && L.info("Uno configuration updated", this.options), this;
  }
  /**
   * Register a plugin
   * @param {Function} plugin - Plugin function (Alpine, options) => void
   * @param {Object} options - Plugin options
   * @returns {UnoInstance}
   */
  use(i, o = {}) {
    return this.initialized && L.warn("Uno.use() called after initialization. Plugin might not work correctly."), this.plugins.push({ plugin: i, options: o }), this;
  }
  /**
   * Initialize Uno
   * @param {Object|Function} alpineOrOptions - Alpine instance or Configuration options
   * @param {Object} options - Configuration options (if first arg is Alpine)
   * @returns {Object} Alpine instance
   */
  init(i = {}, o = {}) {
    if (this.initialized)
      return L.warn("Uno has already been initialized."), X;
    let l = X, e = o;
    return !l && typeof window < "u" && window.Alpine && (l = window.Alpine), i && (i.version || typeof i.start == "function" || typeof i.store == "function") ? l = i : e = i || o, this.options = V(this.options, e), L.setDebug(this.options.debug !== !1), this.options.onError && L.setReporter(this.options.onError), L.info(`Uno v${this.version} initializing...`), this.plugins.forEach(({ plugin: r, options: u }) => {
      try {
        const f = V(this.options, u);
        r(l, f);
      } catch (f) {
        L.error("Failed to load plugin", f);
      }
    }), this.initialized = !0, L.info("Uno initialized"), l;
  }
}
const g = new mt(), st = {
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
}, ht = {
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
}, yt = {
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
const rt = {
  en: st,
  "zh-CN": ht,
  "zh-TW": yt
}, et = Object.keys(rt);
function bt() {
  if (typeof localStorage < "u") {
    const n = localStorage.getItem("uno.lang");
    if (n && et.includes(n)) return n;
  }
  if (typeof navigator < "u") {
    const n = navigator.language, i = n.startsWith("zh") ? n.includes("TW") || n.includes("Hant") ? "zh-TW" : "zh-CN" : "en";
    return et.includes(i) ? i : "en";
  }
  return "en";
}
function vt(n, i) {
  const o = i.split(".");
  let l = n;
  for (const e of o) {
    if (!l || typeof l != "object" || !(e in l)) return;
    l = l[e];
  }
  return l;
}
function at(n, i) {
  for (const o in i)
    i[o] && typeof i[o] == "object" ? (n[o] || Object.assign(n, { [o]: {} }), at(n[o], i[o])) : Object.assign(n, { [o]: i[o] });
}
function wt(n) {
  n.store("i18n", {
    locale: bt(),
    messages: rt,
    init() {
      typeof document < "u" && (document.documentElement.lang = this.locale);
    },
    t(i, o = {}) {
      const l = this.locale, e = this.messages[l] || this.messages.en, r = vt(e, i);
      return (typeof r == "string" ? r : i).replace(/\{(\w+)\}/g, (f, d) => d in o ? String(o[d]) : `{${d}}`);
    },
    setLang(i) {
      if (!this.messages[i]) {
        L.warn(`[Uno] Language ${i} not found`);
        return;
      }
      this.locale = i, typeof localStorage < "u" && localStorage.setItem("uno.lang", i), typeof document < "u" && (document.documentElement.lang = i);
    },
    extend(i) {
      for (const o in i)
        this.messages[o] ? at(this.messages[o], i[o]) : this.messages[o] = i[o];
    }
  });
}
function q(n) {
  return String(n ?? "").replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;").replace(/'/g, "&#39;");
}
const nt = /* @__PURE__ */ new WeakMap();
function lt(n, i) {
  if (!n) return !1;
  const o = String(i);
  let l = nt.get(n);
  return l || (l = /* @__PURE__ */ new Set(), nt.set(n, l)), l.has(o) ? !1 : (l.add(o), !0);
}
function A(n, i = {}) {
  const {
    allowedTags: o = [
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
    allowedAttrs: l = { a: ["href", "target", "rel"], "*": ["class"] }
  } = i, r = String(n ?? "").replace(/<(script|style|iframe|object|embed)\b[\s\S]*?<\/\1\s*>/gi, "").replace(/<(video|audio|svg|math)\b[\s\S]*?<\/\1\s*>/gi, "").replace(/<(img|link|meta|base|source|track)\b[^>]*\/?>/gi, ""), u = document.createElement("template");
  u.innerHTML = r;
  const f = document.createTreeWalker(u.content, NodeFilter.SHOW_ELEMENT), d = [], p = [];
  for (; f.nextNode(); ) {
    const a = f.currentNode, t = a.tagName.toLowerCase();
    if (["script", "style", "iframe", "object", "embed"].includes(t)) {
      p.push(a);
      continue;
    }
    if (!o.includes(t)) {
      d.push(a);
      continue;
    }
    for (const s of [...a.attributes]) {
      const c = s.name.toLowerCase();
      (l[t] || []).concat(l["*"] || []).includes(c) || a.removeAttribute(s.name);
    }
    if (t === "a") {
      const s = (a.getAttribute("href") || "").trim();
      s.startsWith("/") || s.startsWith("https://") || s.startsWith("http://") || s.startsWith("#") || a.removeAttribute("href"), a.getAttribute("rel") || a.setAttribute("rel", "noopener noreferrer");
    }
  }
  for (const a of p)
    a.remove();
  for (const a of d) {
    const t = a.textContent || "";
    t ? a.replaceWith(document.createTextNode(t)) : a.remove();
  }
  return u.innerHTML;
}
function Y(n, i, { trusted: o = !1, alpine: l } = {}) {
  n && (n.innerHTML = o ? String(i ?? "") : A(i), l?.initTree?.(n));
}
const J = {
  /**
   * Ensure a container element exists in the DOM.
   * @param {string} id - Element ID
   * @param {string} html - Inner HTML content
   * @param {Object} attributes - Attributes to set
   * @param {HTMLElement} parent - Parent element (default: document.body)
   */
  ensureContainer(n, i = "", o = {}, l = document.body) {
    if (!document.getElementById(n))
      try {
        const e = document.createElement("div");
        return e.id = n, window.Uno?.options?.csp?.nonce && !o.nonce && (o.nonce = window.Uno.options.csp.nonce), Object.entries(o).forEach(([r, u]) => {
          e.setAttribute(r, u);
        }), i && (e.innerHTML = i), l ? l.appendChild(e) : document.documentElement.appendChild(e), e;
      } catch (e) {
        L.error(`DomManager: Failed to inject container #${n}`, e);
      }
  },
  /**
   * Inject CSS styles with CSP nonce support.
   * @param {string} id - Style element ID
   * @param {string} css - CSS content
   */
  injectStyles(n, i) {
    if (!document.getElementById(n))
      try {
        const o = document.createElement("style");
        o.id = n, window.Uno?.options?.csp?.nonce && o.setAttribute("nonce", window.Uno.options.csp.nonce), o.textContent = i, document.head.appendChild(o);
      } catch (o) {
        L.error(`DomManager: Failed to inject styles #${n}`, o);
      }
  },
  // Exports
  setHTML: Y,
  // Secure HTML setting
  sanitize: A,
  // HTML sanitization
  escape: q
  // HTML escaping
}, xt = [
  "top-left",
  "top-center",
  "top-right",
  "center-left",
  "center",
  "center-right",
  "bottom-left",
  "bottom-center",
  "bottom-right"
], Ct = (n) => {
  const i = n || "center", o = {
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
  return o[i] || o.center;
}, kt = (n) => {
  const i = n || "top-right", o = {
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
  return o[i] || o["top-right"];
}, $t = (n) => {
  const i = n || "top-right";
  let o = "flex-col", l = "items-end", e = "space-y-2";
  return i.includes("bottom") && (o = "flex-col-reverse", e = "space-y-reverse space-y-2"), i.includes("left") ? l = "items-start" : i.includes("center") ? i === "center" ? l = "items-center justify-center" : l = "items-center" : l = "items-end", { flexDir: o, alignItems: l, spaceY: e };
};
function Et(n, i, o) {
  const l = n.getBoundingClientRect(), e = i.getBoundingClientRect(), r = e.width, u = e.height, f = l.width, d = l.height, p = { top: l.top, left: l.left }, a = window.innerWidth, t = window.innerHeight, s = 0, c = 0, v = t, m = a, h = (C, $) => {
    if ($) {
      if (C + r > m)
        return m - r - 10;
      if (C < c)
        return c + 10;
    } else {
      if (C + u > v)
        return v - u - 10;
      if (C < s)
        return s + 10;
    }
    return C;
  };
  let y, b, w = o;
  switch (o) {
    case "left":
      b = p.left - r - 5, b < c && (b = p.left + f + 5, w = "right"), y = h(p.top + d / 2 - u / 2, !1);
      break;
    case "top":
      y = p.top - u - 8, y < s && (y = p.top + d + 8, w = "bottom"), b = p.left + (f - r) / 2, b = h(b, !0);
      break;
    case "right":
      b = p.left + f + 5, b + r > m && (b = p.left - r - 5, w = "left"), y = h(p.top + d / 2 - u / 2, !1);
      break;
    case "bottom":
    default:
      y = p.top + d + 8, y + u > v && (y = p.top - u - 8, w = "top"), b = p.left + (f - r) / 2, b = h(b, !0);
      break;
  }
  return { top: y, left: b, position: w };
}
const St = (n, i = {}) => {
  const o = i.plugins?.tooltip || i.tooltip || {}, l = () => {
    const m = typeof window < "u" ? window.Uno?.options : g.options;
    return m?.plugins?.tooltip || m?.tooltip || o;
  };
  if (!n) {
    L.error("Tooltip: Alpine is undefined. Plugin initialization failed.");
    return;
  }
  let e = null, r = null, u = null, f = 0, d = null;
  const p = () => {
    let m = Number(g.options?.ui?.layers?.tooltipBase ?? 2e4);
    return Number.isFinite(m) || (m = 2e4), m;
  }, a = () => e && r && u && typeof document < "u" && document.contains(e) ? !0 : (e = typeof document < "u" ? document.getElementById("uno-tooltip-container") : null, r = typeof document < "u" ? document.getElementById("uno-tooltip-content") : null, u = typeof document < "u" ? document.getElementById("uno-tooltip-body") : null, !!(e && r && u)), t = () => {
    f || typeof requestAnimationFrame == "function" && (f = requestAnimationFrame(() => {
      f = 0, n.store("tooltip")?.updatePosition?.();
    }));
  }, s = () => {
    if (d || typeof window > "u" || typeof document > "u") return;
    const m = () => t(), h = () => t(), y = (w) => {
      if (w.key !== "Escape") return;
      const C = n.store("tooltip");
      C?.visible && C.hide();
    }, b = (w) => {
      const C = n.store("tooltip");
      if (!C?.visible || C.triggerMode !== "click" || C.clickOutside !== !0 || !a()) return;
      const $ = w.target, k = C.eventEl && typeof C.eventEl.contains == "function" && C.eventEl.contains($), B = r && r.contains($);
      !k && !B && C.hide();
    };
    window.addEventListener("scroll", m, !0), window.addEventListener("resize", h), window.addEventListener("keydown", y), window.addEventListener("click", b), d = () => {
      window.removeEventListener("scroll", m, !0), window.removeEventListener("resize", h), window.removeEventListener("keydown", y), window.removeEventListener("click", b), d = null;
    };
  };
  n.store("tooltip", {
    visible: !1,
    content: "",
    placement: "top",
    top: -9999,
    // Start off-screen
    left: -9999,
    zIndex: p(),
    triggerEl: null,
    eventEl: null,
    triggerMode: "hover",
    clickOutside: !0,
    onShow: null,
    async show(m, h, y = {}) {
      this.triggerEl = m, this.eventEl = y.eventEl || m, this.content = h, this.placement = y.placement || l().placement || "top", this.onShow = y.onShow, this.triggerMode = y.trigger || "hover", this.clickOutside = y.clickOutside !== void 0 ? y.clickOutside : this.triggerMode === "click", this.zIndex = p(), this.visible = !0, this.ensureContainer(), s(), a() && Y(u, this.content, { trusted: y.trustedHTML === !0, alpine: n }), await n.nextTick(), this.updatePosition(), this.onShow && this.onShow(u);
    },
    hide() {
      this.visible = !1, this.triggerEl = null, this.eventEl = null, d && d(), f && (cancelAnimationFrame(f), f = 0);
    },
    updatePosition() {
      if (!this.triggerEl || !a()) return;
      const m = Et(this.triggerEl, r, this.placement);
      this.top = m.top, this.left = m.left;
    },
    ensureContainer() {
      c(), a();
    },
    destroy() {
      this.hide();
    }
  });
  const c = () => {
    if (typeof document > "u") return;
    const h = `fixed px-3 py-2 text-sm font-medium ${g.getStyle("base", "rounded", "rounded-lg")} shadow-sm whitespace-nowrap pointer-events-none border bg-white text-gray-900 border-gray-200 dark:bg-slate-800 dark:text-white dark:border-slate-700`, y = g.getStyle("tooltip", "content", h), b = g.getStyle("tooltip", "body", "pointer-events-auto"), w = `
          <div id="uno-tooltip-content"
               x-show="$store.tooltip.visible"
               role="tooltip"
               x-transition:enter="transition ease-out duration-200"
               x-transition:enter-start="opacity-0 scale-95"
               x-transition:enter-end="opacity-100 scale-100"
               x-transition:leave="transition ease-in duration-150"
               x-transition:leave-start="opacity-100 scale-100"
               x-transition:leave-end="opacity-0 scale-95"
               class="${y}"
              :style="\`top: \${$store.tooltip.top}px; left: \${$store.tooltip.left}px; z-index: \${$store.tooltip.zIndex};\`"
              style="display: none;">
              <div id="uno-tooltip-body" class="${b}"></div>
          </div>
        `;
    J.ensureContainer("uno-tooltip-container", w, {
      "x-data": "{}"
    });
  };
  n.directive(
    "tooltip",
    (m, { expression: h, modifiers: y }, { evaluateLater: b, effect: w, cleanup: C }) => {
      const $ = b(h);
      let k = {
        content: "",
        placement: null,
        trigger: null,
        clickOutside: void 0
      };
      const B = () => {
        let M = l().placement || "top";
        return k.placement && (M = k.placement), y.includes("top") && (M = "top"), y.includes("bottom") && (M = "bottom"), y.includes("left") && (M = "left"), y.includes("right") && (M = "right"), M;
      }, z = () => {
        const E = l();
        return y.includes("click") ? "click" : y.includes("focus") ? "focus" : y.includes("hover") ? "hover" : k.trigger ? k.trigger : E.trigger || "hover";
      }, H = () => (typeof m.disabled == "boolean" && m.disabled === !0 || m.getAttribute("aria-disabled") === "true" || m.getAttribute("disabled") !== null) && m.parentElement || m, N = (E) => {
        if (typeof E == "string" ? k.content = E : typeof E == "object" && (k = { ...k, ...E }), k.content && k.content.startsWith("#") && !k.content.includes(" ")) {
          const M = document.querySelector(k.content);
          M && (k.content = M.innerHTML);
        } else if (k.target) {
          const M = document.querySelector(k.target);
          M && (k.content = M.innerHTML);
        }
      };
      w(() => {
        $((E) => {
          N(E);
        });
      });
      const j = () => {
        const E = n.store("tooltip"), M = z();
        m.setAttribute("aria-describedby", "uno-tooltip-content"), E.show(m, k.content, {
          ...k,
          placement: B(),
          trigger: M,
          eventEl: H()
        });
      }, D = () => {
        m.removeAttribute("aria-describedby"), n.store("tooltip").hide();
      }, F = () => {
        const E = n.store("tooltip");
        E.visible && E.triggerEl === m ? E.hide() : j();
      }, _ = () => {
        z() === "click" && F();
      }, R = () => {
        z() === "hover" && j();
      }, O = () => {
        z() === "hover" && D();
      }, x = () => {
        z() === "focus" && j();
      }, T = () => {
        z() === "focus" && D();
      }, I = H();
      I.addEventListener("click", _), I.addEventListener("mouseenter", R), I.addEventListener("mouseleave", O), I.addEventListener("focus", x), I.addEventListener("blur", T), C(() => {
        I.removeEventListener("click", _), I.removeEventListener("mouseenter", R), I.removeEventListener("mouseleave", O), I.removeEventListener("focus", x), I.removeEventListener("blur", T), D();
      });
    }
  );
  const v = {
    show: (m, h, y) => n.store("tooltip").show(m, h, y),
    hide: () => n.store("tooltip").hide(),
    updatePosition: () => n.store("tooltip").updatePosition(),
    destroy: () => n.store("tooltip").destroy()
  };
  g.tooltip = v, typeof window < "u" && window.Uno && (window.Uno.tooltip = v);
};
const Tt = (n, i = {}) => {
  const o = i.plugins?.notify || i.notify || {}, l = () => {
    const t = typeof window < "u" ? window.Uno?.options : g.options;
    return t?.plugins?.notify || t?.notify || o;
  }, e = () => g.options?.csp?.nonce || typeof window < "u" && window.Uno?.options?.csp?.nonce, r = {
    types: {
      success: "bg-green-50 dark:bg-green-950/40",
      error: "bg-red-50 dark:bg-red-950/40",
      warning: "bg-yellow-50 dark:bg-yellow-950/40",
      info: "bg-blue-50 dark:bg-blue-950/40"
    },
    icons: {
      success: "text-green-500 dark:text-green-300",
      error: "text-red-500 dark:text-red-300",
      warning: "text-yellow-500 dark:text-yellow-300",
      info: "text-blue-500 dark:text-blue-300"
    }
  }, u = (t) => {
    const s = r.types[t] || "bg-white dark:bg-slate-900";
    return g.getStyle("notify", `types.${t}`, s);
  }, f = (t) => {
    const s = r.icons[t] || "text-gray-400 dark:text-gray-500";
    return g.getStyle("notify", `icons.${t}`, s);
  }, d = () => {
    if (document.getElementById("uno-notify-styles")) return;
    const t = document.createElement("style");
    t.id = "uno-notify-styles";
    const s = e();
    s && t.setAttribute("nonce", s), t.textContent = `
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
    `, document.head.appendChild(t);
  };
  let p = null;
  n.store("notify", {
    items: [],
    nextId: 1,
    activeRequests: 0,
    activeProgress: 0,
    // Notification Methods
    show(t, s = {}) {
      const c = l(), {
        type: v = "info",
        timeout: m = c.duration !== void 0 ? c.duration : 3e3,
        position: h = s.position || c.position || "top-right",
        closeIndicator: y = !1,
        showCloseButton: b = !0
      } = s, w = this.nextId++, C = {
        id: w,
        text: t,
        type: v,
        position: h,
        timeout: m,
        closeIndicator: y,
        showCloseButton: b,
        visible: !0
        // For transition
      };
      this.items.unshift(C), m > 0 && setTimeout(() => {
        this.close(w);
      }, m), this.ensureContainer();
    },
    close(t) {
      this.items.find((c) => c.id === t) && (this.items = this.items.filter((c) => c.id !== t));
    },
    success(t, s = {}) {
      const c = typeof s == "number" ? { timeout: s } : s;
      this.show(t, { type: "success", ...c });
    },
    warning(t, s = {}) {
      const c = typeof s == "number" ? { timeout: s } : s;
      this.show(t, { type: "warning", ...c });
    },
    error(t, s = {}) {
      const c = typeof s == "number" ? { timeout: s } : s;
      this.show(t, { type: "error", ...c });
    },
    info(t, s = {}) {
      const c = typeof s == "number" ? { timeout: s } : s;
      this.show(t, { type: "info", ...c });
    },
    loading(t) {
      t ? this.activeRequests++ : this.activeRequests = Math.max(0, this.activeRequests - 1), this.ensureContainer();
    },
    progress(t) {
      t ? this.activeProgress++ : this.activeProgress = Math.max(0, this.activeProgress - 1), this.ensureContainer();
    },
    // Style Helpers
    getTypeClass(t) {
      return u(t);
    },
    getIconClass(t) {
      return f(t);
    },
    ensureContainer() {
      if (p && document.contains(p) || (p = document.getElementById("uno-notify-container"), p)) return;
      d();
      let t = Number(g.options?.ui?.layers?.notifyBase ?? 3e4);
      Number.isFinite(t) || (t = 3e4);
      const s = document.createElement("div");
      s.id = "uno-notify-container";
      const c = e();
      c && s.setAttribute("nonce", c), s.setAttribute("x-data", "{}");
      const v = g.getStyle("base", "rounded", "rounded-lg"), m = g.getStyle(
        "notify",
        "item",
        `pointer-events-auto w-full max-w-sm overflow-hidden ${v} shadow-lg ring-1 ring-black ring-opacity-5`
      ), h = g.getStyle("notify", "content", "p-4"), y = g.getStyle("notify", "row", "flex items-start"), b = g.getStyle("notify", "iconWrap", "flex-shrink-0"), w = g.getStyle("notify", "textWrap", "ml-3 w-0 flex-1 pt-0.5"), C = g.getStyle(
        "notify",
        "text",
        "text-sm font-medium text-gray-900 dark:text-slate-100"
      ), $ = g.getStyle("notify", "closeWrap", "ml-4 flex flex-shrink-0"), k = g.getStyle(
        "notify",
        "closeBtn",
        "inline-flex rounded-md bg-white dark:bg-slate-900 text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
      ), B = g.getStyle("notify", "closeIcon", "h-5 w-5"), z = g.getStyle(
        "notify",
        "indicatorTrack",
        "h-1 bg-gray-200 dark:bg-slate-700 overflow-hidden"
      ), H = g.getStyle(
        "notify",
        "indicatorBar",
        "h-full bg-blue-500 origin-left"
      ), N = (j) => {
        const D = kt(j), F = $t(j);
        return `
          <div class="${g.getStyle(
          "notify",
          "container",
          `fixed flex ${F.flexDir} ${F.alignItems} gap-2 p-4 pointer-events-none ${D}`
        )}" style="z-index: ${t};" x-show="$store.notify.itemsFor('${j}').length">
            <template x-for="item in $store.notify.itemsFor('${j}')" :key="item.id">
              <div
                data-testid="notify-item"
                x-show="item.visible"
                x-transition:enter="transition ease-out duration-300"
                x-transition:enter-start="opacity-0 transform translate-y-2 scale-95"
                x-transition:enter-end="opacity-100 transform translate-y-0 scale-100"
                x-transition:leave="transition ease-in duration-200"
                x-transition:leave-start="opacity-100 transform translate-y-0 scale-100"
                x-transition:leave-end="opacity-0 transform translate-y-2 scale-95"
                class="${m}"
                :class="$store.notify.getTypeClass(item.type)"
              >
                <div class="${h}">
                  <div class="${y}">
                    <div class="${b}">
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
                    <div class="${w}">
                      <p class="${C}" x-text="item.text"></p>
                    </div>
                    <div class="${$}" x-show="item.showCloseButton">
                      <button @click="$store.notify.close(item.id)" class="${k}">
                        <span class="sr-only">Close</span>
                        <svg class="${B}" viewBox="0 0 20 20" fill="currentColor">
                          <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                        </svg>
                      </button>
                    </div>
                  </div>
                </div>
                <div x-show="item.closeIndicator" class="${z}">
                   <div class="${H}"
                        :style="\`animation: notify-progress \${item.timeout}ms linear forwards\`"></div>
                </div>
              </div>
            </template>
          </div>
        `;
      };
      s.innerHTML = `
        <div
          id="uno-progress-bar"
          x-show="$store.notify.activeProgress > 0"
          class="fixed top-0 left-0 right-0 h-1 bg-transparent"
          style="z-index: ${t + 1};"
        >
          <div class="uno-progress-anim h-full bg-blue-500"></div>
        </div>
        ${xt.map(N).join("")}
      `, document.body.appendChild(s), p = s, typeof n?.initTree == "function" && n.initTree(s);
    },
    itemsFor(t) {
      const s = t || "top-right";
      return this.items.filter((c) => (c.position || "top-right") === s);
    }
  });
  const a = {
    show: (t, s) => n.store("notify").show(t, s),
    success: (t, s) => n.store("notify").success(t, s),
    error: (t, s) => n.store("notify").error(t, s),
    warning: (t, s) => n.store("notify").warning(t, s),
    info: (t, s) => n.store("notify").info(t, s),
    loading: (t) => n.store("notify").loading(t),
    progress: (t) => n.store("notify").progress(t),
    close: (t) => n.store("notify").close(t),
    closeAll: () => {
      n.store("notify").items = [];
    }
  };
  g.notify = a, typeof window < "u" && window.Uno && (window.Uno.notify = a);
};
class Lt {
  constructor(i = 50) {
    this.maxSize = i, this.cache = /* @__PURE__ */ new Map();
  }
  get(i) {
    if (!this.cache.has(i)) return;
    const o = this.cache.get(i);
    return this.cache.delete(i), this.cache.set(i, o), o;
  }
  set(i, o) {
    if (this.cache.has(i) && this.cache.delete(i), this.cache.set(i, o), this.cache.size > this.maxSize) {
      const l = this.cache.keys().next().value;
      this.cache.delete(l);
    }
  }
  has(i) {
    return this.cache.has(i);
  }
  clear() {
    this.cache.clear();
  }
}
const W = new Lt(50);
function Mt(n) {
  const i = String(n ?? "").trim();
  if (!i) return [];
  try {
    const o = JSON.parse(i);
    if (typeof o == "string") return [[o, {}]];
    if (o && typeof o == "object") return Object.entries(o);
  } catch {
  }
  return i.split(",").map((o) => o.trim()).filter(Boolean).map((o) => [o, {}]);
}
function Bt(n, i, o) {
  const l = n || document.body || document.documentElement;
  if (l) {
    if (typeof window < "u" && window.htmx?.trigger) {
      window.htmx.trigger(l, i, o);
      return;
    }
    l.dispatchEvent(new CustomEvent(i, { bubbles: !0, detail: o }));
  }
}
function It(n, i) {
  const o = n.headers.get("HX-Trigger");
  if (o)
    for (const [l, e] of Mt(o))
      Bt(i, l, e);
}
async function zt(n, i = {}) {
  const {
    method: o = "GET",
    htmxTarget: l,
    swap: e = "innerHTML",
    htmxConfig: r = {},
    cache: u = !1,
    cacheKey: f
  } = i, d = f || n;
  if (u && o === "GET" && W.has(d))
    return W.get(d);
  if (typeof window < "u" && window.htmx?.ajax) {
    const p = l || document.createElement("div");
    l || (p.style.display = "none"), l || (document.body ? document.body.appendChild(p) : document.documentElement.appendChild(p));
    try {
      await window.htmx.ajax(o, n, { target: p, swap: e, ...r });
      const a = p.innerHTML;
      return u && o === "GET" && W.set(d, a), a;
    } finally {
      l || p.remove();
    }
  }
  return jt(n, i);
}
async function Q(n, i = {}) {
  const {
    method: o = "GET",
    headers: l = {},
    body: e,
    credentials: r = "same-origin",
    triggerTarget: u = typeof document < "u" ? document.body : null,
    cache: f = !1,
    cacheKey: d,
    signal: p
    // AbortSignal
  } = i, a = d || n;
  if (f && o === "GET" && W.has(a))
    return W.get(a);
  const t = {
    "HX-Request": "true",
    "X-Requested-With": "XMLHttpRequest",
    ...l
  }, s = await fetch(n, { method: o, headers: t, body: e, credentials: r, signal: p });
  if (It(s, u), !s.ok) {
    const m = new Error(`Request failed: ${s.status}`);
    throw m.status = s.status, m;
  }
  const c = s.headers.get("content-type") || "";
  let v;
  return c.includes("application/json") ? v = await s.json() : v = await s.text(), f && o === "GET" && W.set(a, v), v;
}
async function jt(n, i = {}) {
  const o = await Q(n, i);
  return typeof o == "string" ? o : JSON.stringify(o, null, 2);
}
const Dt = [
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
function _t(n, i = {}) {
  const o = (d) => {
    if (!(d instanceof HTMLElement)) return !1;
    const p = window.getComputedStyle(d);
    return p.display === "none" || p.visibility === "hidden" ? !1 : d.getClientRects().length > 0;
  }, l = () => Array.from(n.querySelectorAll(Dt)).filter((d) => o(d)), e = () => l()[0], r = () => {
    const d = l();
    return d[d.length - 1];
  }, u = (d) => {
    if (d.key === "Tab") {
      if (d.shiftKey) {
        if (document.activeElement === e()) {
          d.preventDefault();
          const p = r();
          p && p.focus();
        }
      } else if (document.activeElement === r()) {
        d.preventDefault();
        const p = e();
        p && p.focus();
      }
    }
  }, f = () => {
    if (n.contains(document.activeElement)) return;
    const d = e();
    d && d.focus();
  };
  return typeof requestAnimationFrame == "function" ? requestAnimationFrame(() => requestAnimationFrame(f)) : setTimeout(f, 0), n.addEventListener("keydown", u), () => {
    n.removeEventListener("keydown", u);
  };
}
const Ht = (n, i = {}) => {
  const o = i.plugins?.box || i.box || {}, l = () => {
    const a = typeof window < "u" ? window.Uno?.options : g.options;
    return a?.plugins?.box || a?.box || o;
  }, e = (a, t) => g.getStyle("box", a, t), r = (a) => `<div class="${e("error", "p-4 text-red-500")}">${q(a)}</div>`;
  let u = null;
  const f = /* @__PURE__ */ new Map(), d = (a, t) => {
    const s = `${t}:${a}`, c = f.get(s);
    if (c && typeof document < "u" && document.contains(c))
      return c;
    const v = typeof document < "u" ? document.getElementById(t + a) || u?.querySelector?.(`#${t + a}`) : null;
    return v && f.set(s, v), v;
  };
  n.directive("uno-trap", (a, { expression: t }, { cleanup: s }) => {
    const c = _t(a);
    s(c);
  }), n.store("box", {
    items: [],
    // Public API
    alert(a, t = {}) {
      const s = typeof t == "string" ? { title: t } : t;
      return this.open("dialog", {
        content: a,
        showOkButton: !0,
        ...s
      });
    },
    confirm(a, t = {}) {
      return new Promise((s) => {
        this.open("confirm", {
          content: a,
          onOk: () => s(!0),
          onCancel: () => s(!1),
          onClose: () => s(!1),
          // Closing without choice = cancel
          ...t
        });
      });
    },
    prompt(a, t = "", s = {}) {
      return new Promise((c) => {
        const v = "prompt-" + Math.random().toString(36).slice(2, 11), m = e("promptWrapper", "space-y-3"), h = e("promptText", "text-slate-700 dark:text-slate-200"), y = e(
          "promptInput",
          "w-full px-3 py-2 border rounded focus:ring-2 focus:ring-blue-500 outline-none"
        ), b = `
          <div class="${m}">
            <p class="${h}">${q(a)}</p>
            <input id="${v}" type="text" class="${y}" value="${q(t)}">
          </div>
        `;
        this.open("confirm", {
          title: s.title || "Prompt",
          content: b,
          trustedHTML: !0,
          onOk: () => {
            const w = document.getElementById(v);
            c(w ? w.value : null);
          },
          onCancel: () => c(null),
          ...s
        });
      });
    },
    async open(a, t = {}) {
      let s = "dialog", c = t;
      const v = l();
      typeof a == "object" ? (c = a, s = c.type || "dialog") : s = a;
      const m = c.boxId || Math.random().toString(36).slice(2, 10), h = {
        id: m,
        type: s,
        title: c.title || "",
        content: c.content || "",
        trustedHTML: c.trustedHTML === !0,
        width: c.width !== void 0 ? c.width : c.widthClass ? void 0 : v.defaultWidth || 400,
        height: c.height,
        widthClass: c.widthClass,
        heightClass: c.heightClass,
        containerClass: c.containerClass,
        modal: c.modal !== !1,
        // default true
        overlayOpacity: c.overlayOpacity ?? 0.5,
        position: c.position || "center",
        showCloseButton: c.showCloseButton !== void 0 ? c.showCloseButton : !0,
        // Buttons
        okText: c.okButtonText,
        cancelText: c.cancelButtonText,
        onOk: c.onOk,
        onCancel: c.onCancel,
        onClose: c.onClose,
        onBeforeClose: c.onBeforeClose,
        onContentLoad: c.onContentLoad,
        showOk: s === "confirm" || s === "message" || s === "dialog" && c.showOkButton,
        showCancel: s === "confirm" || s === "dialog" && c.showCancelButton,
        // Cache Option
        cache: c.cache === !0,
        // State
        loading: !1,
        visible: !0,
        abortController: null
      };
      h.trustedHTML || (h.content = A(h.content), h.trustedHTML = !0);
      const y = n.store("i18n");
      if (h.okText || (h.okText = y ? y.t("box.ok") : "OK"), h.cancelText || (h.cancelText = y ? y.t("box.cancel") : "Cancel"), c.targetId) {
        const b = document.getElementById(c.targetId);
        b ? (h.content = b.innerHTML, c.trustedHTML || (h.content = A(h.content)), h.trustedHTML = !0) : (L.warn("Box: Target element not found:", c.targetId), h.content = r("Error: Content target not found"), h.trustedHTML = !0);
      } else if (c.hxUrl) {
        h.loading = !0, this.add(h);
        try {
          const b = await zt(c.hxUrl, { method: "GET", cache: h.cache }), w = this.items.find((C) => C.id === m);
          w && (w.content = c.trustedHTML ? b : A(b), w.trustedHTML = !0, w.loading = !1, w.onContentLoad && n.nextTick(() => {
            w.onContentLoad(d(m, "uno-box-content"));
          }));
        } catch (b) {
          L.error("Box: HTMX Load Failed", b);
          const w = this.items.find((C) => C.id === m);
          w && (w.loading = !1, w.trustedHTML = !0, w.content = r("Error loading content via HTMX"));
        }
        return m;
      } else if (c.ajaxUrl) {
        h.loading = !0, h.abortController = new AbortController(), this.add(h);
        try {
          const b = await Q(c.ajaxUrl, {
            cache: h.cache,
            signal: h.abortController.signal
          }), w = e(
            "ajaxPre",
            "bg-gray-100 dark:bg-slate-800 p-3 rounded overflow-x-auto whitespace-pre-wrap"
          );
          let C = typeof b == "object" ? `<pre class="${w}">${q(JSON.stringify(b, null, 2))}</pre>` : String(b);
          c.trustedHTML || (C = A(C));
          const $ = this.items.find((k) => k.id === m);
          $ && ($.content = C, $.trustedHTML = !0, $.loading = !1, n.nextTick(() => {
            const k = d(m, "uno-box-content");
            k && $.onContentLoad && $.onContentLoad(k);
          }));
        } catch (b) {
          if (b.name === "AbortError")
            return L.info("Box: Request cancelled"), m;
          L.error("Box: Ajax failed", b);
          const w = this.items.find((C) => C.id === m);
          w && (w.content = r("Error loading content"), w.trustedHTML = !0, w.loading = !1);
        }
        return m;
      } else if (c.iframeUrl) {
        const b = e("iframe", "w-full h-full");
        h.content = `<iframe frameborder="0" class="${b}" src="${q(c.iframeUrl)}"></iframe>`, h.trustedHTML = !0;
      }
      return this.add(h), m;
    },
    add(a) {
      this.items.push(a), this.ensureContainer();
    },
    close(a) {
      const t = this.items.find((s) => s.id === a);
      t && (t.abortController && t.abortController.abort(), t.onBeforeClose && t.onBeforeClose(), this.items = this.items.filter((s) => s.id !== a), f.delete(`uno-box-content:${a}`), t.onClose && t.onClose());
    },
    closeAll() {
      this.items.forEach((a) => {
        a.abortController && a.abortController.abort();
      }), this.items = [], f.clear();
    },
    getPositionClass(a) {
      const t = e(
        "positionBase",
        "fixed inset-0 flex transition-opacity duration-300 p-4 pointer-events-none"
      ), s = Ct(a);
      return `${t} ${s}`;
    },
    ensureContainer() {
      p();
    }
  });
  const p = () => {
    if (typeof document < "u" && document.getElementById("uno-box-container")) return;
    const a = g.getStyle("base", "rounded", "rounded-lg"), t = g.getStyle("base", "shadow", "shadow-md"), s = g.getStyle("base", "primary", "bg-blue-600");
    let c = Number(g.options.ui?.layers?.boxBase ?? 1e4);
    Number.isFinite(c) || (c = 1e4);
    let v = Number(g.options.ui?.layers?.boxStep ?? 1);
    (!Number.isFinite(v) || v <= 0) && (v = 1);
    const h = g.getStyle("box", "overlay", "absolute inset-0 bg-black transition-opacity duration-300 pointer-events-auto"), y = `relative bg-white dark:bg-gray-900 ${a} ${t} ring-1 ring-black/5 overflow-hidden flex flex-col max-h-[90vh] transition-all duration-300 transform pointer-events-auto outline-none`, b = g.getStyle("box", "container", y), C = g.getStyle("box", "cancelBtn", "px-4 py-2 bg-slate-100 dark:bg-slate-700 border border-slate-200 dark:border-slate-600 rounded-md text-sm font-medium text-slate-800 dark:text-slate-200 hover:bg-slate-200 dark:hover:bg-slate-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors"), $ = `px-4 py-2 ${s} border border-transparent rounded-md text-sm font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500`, k = g.getStyle("box", "okBtn", $), z = e("header", "flex justify-between items-center p-4 border-b border-slate-200 dark:border-slate-700"), N = e("title", "text-lg font-semibold text-slate-800 dark:text-slate-200"), D = e("closeBtn", "text-slate-400 hover:text-slate-500 focus:outline-none"), _ = e("closeBtnFloating", "absolute top-2 right-2 z-10 text-slate-400 hover:text-slate-500 focus:outline-none"), O = e("closeIcon", "h-6 w-6"), T = e("closeIconSmall", "h-5 w-5"), E = e("content", "p-4 overflow-auto flex-1 text-slate-600 dark:text-slate-300"), dt = e("loadingWrap", "flex justify-center items-center h-full"), ut = e("spinner", "animate-spin h-8 w-8 text-blue-500"), ft = e("footer", "p-4 border-t border-slate-200 dark:border-slate-700 flex justify-end space-x-3 bg-slate-50 dark:bg-slate-800"), pt = `
      <template x-for="(item, index) in $store.box.items" :key="item.id">
          <div :id="'uno-box-outer-container' + item.id"
               :class="$store.box.getPositionClass(item.position)"
               :style="\`z-index: \${${c} + (index * ${v})}\`"
               role="dialog"
               aria-modal="true"
               @keydown.escape.window="if(item.modal) $store.box.close(item.id)">
               
               <!-- Overlay -->
               <div :id="'uno-box-overlay' + item.id"
                    class="${h}"
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
                    class="${b}"
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
                    <div :id="'uno-box-header' + item.id" x-show="item.title" class="${z}">
                        <h3 :id="'uno-box-title' + item.id" class="${N}" x-text="item.title"></h3>
                        <button :id="'uno-box-close' + item.id"
                                x-show="item.showCloseButton"
                                @click="$store.box.close(item.id)"
                                class="${D}">
                            <span class="sr-only">Close</span>
                            <svg class="${O}" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button>
                    </div>
                    
                    <!-- Close button if no title -->
                    <button x-show="!item.title && item.showCloseButton" 
                            :id="'uno-box-close' + item.id"
                            @click="$store.box.close(item.id)" 
                            class="${_}">
                        <span class="sr-only">Close</span>
                        <svg class="${T}" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>

                    <!-- Content -->
                    <div :id="'uno-box-wrapper' + item.id" class="${E}">
                        <div x-show="item.loading" class="${dt}">
                            <svg class="${ut}" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                            </svg>
                        </div>
                        <div :id="'uno-box-content' + item.id" x-show="!item.loading" x-html="item.content"></div>
                    </div>

                    <!-- Footer -->
                    <div :id="'uno-box-footer' + item.id" x-show="item.showOk || item.showCancel" class="${ft}">
                        <button :id="'uno-box-cancel' + item.id"
                                x-show="item.showCancel"
                                @click="item.onCancel ? (item.onCancel(), $store.box.close(item.id)) : $store.box.close(item.id)"
                                x-text="item.cancelText"
                                class="${C}">
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
                                class="${k}">
                        </button>
                    </div>
               </div>
          </div>
      </template>
    `;
    u = J.ensureContainer("uno-box-container", pt, {
      "x-data": "{}"
    });
  };
  typeof window < "u" && lt(document, "uno:box:htmx-afterSwap") && document.addEventListener("htmx:afterSwap", (a) => {
    a.detail.elt && n && n.initTree(a.detail.elt);
  }), g.box = {
    dialog: (a) => n.store("box").open("dialog", a),
    confirm: (a, t) => n.store("box").confirm(a, t),
    alert: (a, t) => n.store("box").alert(a, t),
    prompt: (a, t, s) => n.store("box").prompt(a, t, s),
    close: (a) => n.store("box").close(a),
    closeAll: () => n.store("box").closeAll()
  };
}, ot = {
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
function qt(n) {
  const i = n.match(/^\/(.*?)\/([gimsuy]*)$/);
  return i ? { pattern: i[1], flags: i[2] } : { pattern: n, flags: "" };
}
function Nt(n, i = {}) {
  const { defaultFlags: o = "i", maxPatternLength: l = 200 } = i, { pattern: e, flags: r } = qt(String(n ?? ""));
  if (!e || e.length > l) return null;
  const u = String(r || o || "");
  if (!/^[gimsuy]*$/.test(u)) return null;
  try {
    return new RegExp(e, u);
  } catch {
    return null;
  }
}
function Ft(n) {
  try {
    return n ? JSON.parse(n) : [];
  } catch {
    return [];
  }
}
const it = /* @__PURE__ */ new WeakMap(), G = Symbol("uno.valid.debounce.cancelled"), Ot = (n) => String(n || "").split(/\s+/).filter(Boolean), S = (n, i, o) => {
  if (n)
    for (const l of Ot(i))
      o ? n.classList.add(l) : n.classList.remove(l);
}, Rt = (n, i = {}) => String(n ?? "").replace(/\{(\w+)\}/g, (o, l) => l in i ? String(i[l]) : `{${l}}`), Ut = (n, i, o, l) => {
  if (!n) return l();
  n._debounceMap || (n._debounceMap = /* @__PURE__ */ new Map());
  const e = n._debounceMap, r = e.get(i);
  r && (clearTimeout(r.timer), r.resolve(G), e.delete(i));
  const u = Math.max(0, Number(o) || 0);
  return u === 0 ? l() : new Promise((f, d) => {
    const p = setTimeout(async () => {
      e.delete(i);
      try {
        const a = await l();
        f(a);
      } catch (a) {
        d(a);
      }
    }, u);
    e.set(i, { timer: p, resolve: f });
  });
};
const At = (n) => {
  const i = {
    invalid: "<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 12 12' fill='none' stroke='currentColor' class='w-4 h-4'><circle cx='6' cy='6' r='4.5'/><path stroke-linecap='round' d='M5.8 3.6h.4L6 6.5zM6 8.2a.6.6 0 110-1.2.6.6 0 010 1.2z'/></svg>",
    valid: "<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 8 8' class='w-4 h-4'><path fill='currentColor' d='M2.3 6.73L.6 4.53c-.4-1.04.46-1.4 1.1-.8l1.1 1.4 3.4-3.8c.6-.63 1.6-.27 1.2.7l-4 4.6c-.43.5-.8.4-1.1.1z'/></svg>"
  }, o = (e, r) => {
    const u = n.store("i18n"), f = `valid.${e}`;
    if (u && typeof u.t == "function") {
      const p = u.t(f, r);
      if (p !== f) return p;
    }
    const d = st?.valid?.[e];
    return typeof d == "string" ? Rt(d, r) : f;
  };
  n.store("valid", {
    callbacks: {},
    registerCallback(e, r) {
      typeof r == "function" && (this.callbacks[e] = r);
    },
    async validateForm(e) {
      if (!e) return !1;
      if (e._x_dataStack) {
        const r = e._x_dataStack.find((u) => u.validateAll);
        if (r)
          return await r.validateAll();
      }
      return !1;
    },
    clearForm(e) {
      if (e && e._x_dataStack) {
        const r = e._x_dataStack.find((u) => u.reset);
        r && r.reset();
      }
    }
  });
  const l = {
    // Validate all forms or specific form?
    // Mostly helper to set state manually
    set: (e, { valid: r, message: u }) => {
      e && e.dispatchEvent(
        new CustomEvent("valid-set", {
          detail: { valid: r, message: u, name: e.name },
          bubbles: !0
        })
      );
    },
    validate: async (e) => {
      const r = e?.tagName === "FORM" ? e : e?.closest?.("form");
      return n.store("valid").validateForm(r);
    },
    clear: (e) => {
      const r = e?.tagName === "FORM" ? e : e?.closest?.("form");
      n.store("valid").clearForm(r);
    }
  };
  g.valid = l, typeof window < "u" && window.Uno && (window.Uno.valid = l), n.data("validForm", () => ({
    fields: {},
    // Reactive state: { [name]: { value, error, valid, ... } }
    init() {
      this.$el.noValidate = !0, this.$el.addEventListener("valid-set", (e) => {
        const { name: r, valid: u, message: f } = e.detail;
        if (r && this.fields[r])
          this.fields[r].valid = u, this.fields[r].error = f || null, u || (this.fields[r].touched = !0);
        else if (e.target && e.target.name && this.fields[e.target.name]) {
          const d = e.target.name;
          this.fields[d].valid = u, this.fields[d].error = f || null, u || (this.fields[d].touched = !0);
        }
      });
    },
    // Reset validation state (API: Uno.valid.clear)
    reset() {
      const e = g.getStyle("valid", "stateInvalid", "is-invalid"), r = g.getStyle("valid", "stateValid", "is-valid"), u = g.getStyle("valid", "stateValidating", "validating"), f = g.getStyle("valid", "inputError", "!border-red-500"), d = g.getStyle("valid", "inputValid", "!border-green-500"), p = g.getStyle("valid", "inputPaddingIcon", "pr-10");
      Object.keys(this.fields).forEach((a) => {
        const t = this.fields[a];
        if (t._validationSeq += 1, t.error = null, t.valid = !0, t.touched = !1, t.dirty = !1, t.validating = !1, t._debounceMap) {
          for (const c of t._debounceMap.values())
            clearTimeout(c.timer), c.resolve(G);
          t._debounceMap.clear();
        }
        S(t.el, e, !1), S(t.el, r, !1), S(t.el, u, !1), S(t.el, f, !1), S(t.el, d, !1), S(t.el, p, !1), t.icon && (t.icon.innerHTML = "", t.icon.style.opacity = "0"), t.spinner && (t.spinner.style.display = "");
        const s = t.wrapper || t.el?.closest?.(".uno-input-wrapper");
        s && s.querySelectorAll("[data-valid-error-for]").forEach((c) => {
          c.getAttribute("data-valid-error-for") === a && c.remove();
        });
      });
    },
    // Register field (called by x-validate directive)
    register(e, r, u = [], f = null) {
      if (this.fields[e]) return;
      const d = [...u];
      if ((r.hasAttribute("required") || r.dataset.validRequired === "true") && d.push({ type: "required" }), (r.type === "email" || r.dataset.validType === "email") && d.push("email"), (r.type === "url" || r.dataset.validType === "url") && d.push("url"), r.dataset.validType) {
        const p = r.dataset.validType;
        p !== "email" && p !== "url" && d.push(p);
      }
      r.dataset.validPattern && d.push({ type: "regexp", expression: r.dataset.validPattern });
      for (const p of d) {
        const a = typeof p == "object" && p ? p : null;
        if (a?.type === "compareto" && typeof a.destid == "string" && !a._destEl) {
          const t = document.getElementById(a.destid);
          a._destEl = t || null, a._destResolved = !0;
        }
      }
      this.fields[e] = {
        el: r,
        wrapper: f?.wrapper || null,
        icon: f?.icon || null,
        spinner: f?.spinner || null,
        rules: d,
        value: r.value,
        error: null,
        valid: !0,
        dirty: !1,
        touched: !1,
        validating: !1,
        _validationSeq: 0
      };
    },
    update(e, r) {
      this.fields[e] && (this.fields[e].value = r, this.fields[e].dirty = !0, this.validate(e));
    },
    touch(e) {
      this.fields[e] && (this.fields[e].touched = !0, this.validate(e));
    },
    async validate(e) {
      const r = this.fields[e];
      if (!r) return !0;
      const u = r._validationSeq += 1;
      r.validating = !0, r.error = null;
      const f = (r.value || "").trim();
      let d = null;
      for (const p of r.rules) {
        const a = p.type || p, t = typeof p == "object" ? p : {};
        if (!f) {
          if (a === "required") {
            d = t.msg || o("required");
            break;
          }
          if (a === "emptyquit") break;
          continue;
        }
        if (ot[a])
          ot[a].test(f) || (d = t.msg || o(a) || o("pattern"));
        else if (a === "regexp" && t.expression) {
          let s = it.get(t);
          if (!s) {
            const c = String(t.expression);
            s = Nt(c, { defaultFlags: "i" }), s && it.set(t, s);
          }
          s && s.test(f) || (d = t.msg || o("pattern"));
        } else if (a === "limit")
          t.min && f.length < t.min && (d = t.msg || o("min", { n: t.min })), t.max && f.length > t.max && (d = t.msg || o("max", { n: t.max }));
        else if (a === "range") {
          const s = parseFloat(f);
          t.min !== void 0 && t.max !== void 0 && (s < t.min || s > t.max) ? d = t.msg || o("range", { min: t.min, max: t.max }) : t.min !== void 0 && s < t.min ? d = t.msg || o("rangeMin", { n: t.min }) : t.max !== void 0 && s > t.max && (d = t.msg || o("rangeMax", { n: t.max }));
        } else if (a === "compareto") {
          !t._destResolved && typeof t.destid == "string" && (t._destEl = document.getElementById(t.destid) || null, t._destResolved = !0);
          const s = t._destEl;
          s && f !== s.value && (d = t.msg || o("compareto"));
        } else if (a === "tag") {
          const s = t.mark || ",", c = f.split(s).filter((v) => v.trim());
          t.max && c.length > t.max && (d = t.msg || o("tagMax", { n: t.max }));
        }
        if (d) break;
        if (a === "callback" && t.fn) {
          const s = n.store("valid").callbacks[t.fn];
          if (s)
            try {
              const c = Number.isFinite(t.debounce) ? t.debounce : Number.isFinite(g.options?.plugins?.valid?.callbackDebounce) ? g.options.plugins.valid.callbackDebounce : Number.isFinite(g.options?.valid?.callbackDebounce) ? g.options.valid.callbackDebounce : 300, v = await Ut(
                r,
                `callback:${t.fn}`,
                c,
                () => s(r.el.id, f)
              );
              if (v === G || r._validationSeq !== u) return r.valid;
              v === !1 || typeof v == "object" && v.valid === !1 ? d = (typeof v == "object" ? v.message : "") || t.msg || o("callback") : typeof v == "string" && (d = v);
            } catch (c) {
              L.warn("[Uno.valid] Async validation failed", c), d = t.msg || o("callback");
            }
        }
        if (d) break;
      }
      return r._validationSeq !== u || (r.error = d, r.valid = !d, r.validating = !1), r.valid;
    },
    async validateAll() {
      const e = Object.keys(this.fields);
      return (await Promise.all(e.map((u) => this.validate(u)))).every((u) => u);
    },
    async submit(e) {
      if (e && (e.preventDefault(), e.stopPropagation()), Object.keys(this.fields).forEach((u) => {
        this.fields[u].touched = !0;
      }), await this.validateAll())
        this.$el.dispatchEvent(new CustomEvent("valid-submit", { bubbles: !0 })), !this.$el.hasAttribute("hx-post") && !this.$el.hasAttribute("hx-get") && this.$el.submit();
      else {
        const u = Object.keys(this.fields).find((f) => !this.fields[f].valid);
        u && this.fields[u].el && this.fields[u].el.focus();
      }
    }
  })), n.directive("validate", (e, r, { effect: u, cleanup: f }) => {
    const d = e.getAttribute("name");
    if (!d) return;
    const p = g.getStyle("valid", "wrapper", "relative w-full"), a = g.getStyle(
      "valid",
      "iconBase",
      "absolute right-0 top-0 bottom-0 flex items-center justify-center w-10 pointer-events-none transition-opacity duration-200"
    ), t = g.getStyle("valid", "inputError", "!border-red-500"), s = g.getStyle("valid", "inputValid", "!border-green-500"), c = g.getStyle("valid", "iconError", "text-red-500 dark:text-red-400"), v = g.getStyle("valid", "iconValid", "text-green-500 dark:text-green-400"), m = g.getStyle("valid", "stateInvalid", "is-invalid"), h = g.getStyle("valid", "stateValid", "is-valid"), y = g.getStyle("valid", "stateValidating", "validating"), b = g.getStyle("valid", "inputPaddingIcon", "pr-10"), w = () => {
      let x = e.parentNode;
      x.classList.contains("uno-input-wrapper") ? S(x, p, !0) : (x = document.createElement("div"), x.className = "uno-input-wrapper", S(x, p, !0), e.parentNode.insertBefore(x, e), x.appendChild(e));
      let T = x.querySelector(".uno-valid-icon");
      return T ? T.className = `uno-valid-icon ${a}` : (T = document.createElement("div"), T.className = `uno-valid-icon ${a}`, T.style.opacity = "0", x.appendChild(T)), { wrapper: x, icon: T };
    }, { wrapper: C, icon: $ } = w();
    let k = null, B = C.nextElementSibling, z = 0;
    for (; B && z < 3 && !(B.nodeType === 1 && (B.classList.contains("animate-spin") ? k = B : k = B.querySelector(".animate-spin"), k)); )
      B = B.nextElementSibling, z++;
    let H = [];
    const N = e.getAttribute("valid");
    if (N) {
      const x = Ft(N);
      H = Array.isArray(x) ? x : [];
    }
    const j = H.find(
      (x) => (x.type === "limit" || x === "limit") && typeof x == "object"
    ), D = j?.max, F = j?.show === !0 && Number.isFinite(D);
    let _ = null;
    u(() => {
      const x = n.$data(e);
      x && typeof x.register == "function" && (_ = x, x.register(d, e, H, { wrapper: C, icon: $, spinner: k }));
    });
    const R = () => {
      const x = _ || n.$data(e);
      F && D && e.value.length > D && (e.value = e.value.substring(0, D)), x && x.update && x.update(d, e.value);
    }, O = () => {
      const x = _ || n.$data(e);
      x && x.touch && x.touch(d);
    };
    e.addEventListener("input", R), e.addEventListener("blur", O), u(() => {
      const x = _ || n.$data(e);
      if (x && x.fields && x.fields[d]) {
        const T = x.fields[d], I = !!T.error, E = T.valid && T.dirty && !T.error;
        S(e, m, I), S(e, h, E), S(e, y, T.validating), $.innerHTML = "", $.style.opacity = "0", I ? (S(e, t, !0), S(e, s, !1), S(e, b, !0), $.innerHTML = i.invalid, $.className = `uno-valid-icon ${a} ${c}`, $.style.opacity = "1") : E ? (S(e, s, !0), S(e, t, !1), S(e, b, !0), $.innerHTML = i.valid, $.className = `uno-valid-icon ${a} ${v}`, $.style.opacity = "1") : (S(e, t, !1), S(e, s, !1), S(e, b, !1)), k && (k.style.display = T.validating ? "block" : "");
      }
    }), f(() => {
      e.removeEventListener("input", R), e.removeEventListener("blur", O);
    });
  }), n.directive("valid-error", (e, { expression: r }, { effect: u }) => {
    const f = r;
    e.style.display = "none", u(() => {
      const d = n.$data(e);
      if (d && d.fields && d.fields[f]) {
        const a = d.fields[f].error;
        if (a) {
          e.textContent = a, e.style.display = "block";
          const t = g.getStyle(
            "valid",
            "helperText",
            "text-red-600 dark:text-red-400 text-xs"
          );
          e.className = t;
        } else
          e.textContent = "", e.style.display = "none";
      }
    });
  }), n.directive("valid-count", (e, { expression: r }, { effect: u }) => {
    const f = r;
    u(() => {
      const d = n.$data(e);
      if (d && d.fields && d.fields[f]) {
        const p = d.fields[f], a = p.rules.find((t) => t.type === "limit" || t === "limit");
        if (a && a.max) {
          const t = (p.value || "").length, s = a.max - t, c = n.store("i18n");
          e.textContent = c && c.t("charCount", { count: s }) || `Remaining: ${s}`, e.style.display = "block";
        } else
          e.style.display = "none";
      }
    });
  });
};
let ct;
const tt = (n = null, i = {}) => ({
  activeTab: n,
  options: i,
  _cleanup: [],
  _pending: /* @__PURE__ */ new Map(),
  _seq: 0,
  init() {
    if (!lt(this.$el, "tabs_init")) return;
    this.$nextTick(() => {
      const l = window.location.hash.slice(1);
      let e = !1;
      if (l) {
        const r = this.$el.querySelector(`[data-tab-hash="${l}"]`);
        r && (this.activeTab = r.getAttribute("data-tab-target"), e = !0);
      }
      if (!e && !this.activeTab) {
        const r = this.$el.querySelector('[role="tab"]');
        r && (this.activeTab = r.getAttribute("data-tab-target"));
      }
    });
    const o = () => {
      const l = window.location.hash.slice(1);
      if (l) {
        const e = this.$el.querySelector(`[data-tab-hash="${l}"]`);
        e && (this.activeTab = e.getAttribute("data-tab-target"));
      }
    };
    window.addEventListener("hashchange", o), this._cleanup.push(() => window.removeEventListener("hashchange", o));
  },
  destroy() {
    for (const o of this._pending.values())
      o?.abortController?.abort?.();
    this._pending.clear(), this._cleanup.forEach((o) => o()), this._cleanup = [];
  },
  isActive(o) {
    return this.activeTab === o;
  },
  async open(o, l = null, e = null, r = null) {
    this.activeTab = o, e && window.location.hash.slice(1) !== e && window.history.replaceState(null, null, `#${e}`);
    const u = l || r;
    if (!u) return;
    const f = document.getElementById(o);
    if (!f || f.getAttribute("data-tab-loaded") === "true") return;
    const d = this._pending.get(o);
    d?.abortController && d.abortController.abort();
    const p = new AbortController(), a = this._seq += 1;
    this._pending.set(o, { abortController: p, seq: a });
    const t = ct || window.Alpine, s = g.getStyle("base", "rounded", "rounded-lg"), c = g.getStyle(
      "tab",
      "loading",
      `animate-pulse h-32 w-full ${s} bg-gray-200 dark:bg-slate-700`
    ), v = g.getStyle("tab", "error", "p-4 text-red-600 dark:text-red-400");
    Y(f, `<div class="${c}"></div>`, { trusted: !1, alpine: t });
    try {
      const m = await Q(u, { cache: !0, signal: p.signal }), h = this._pending.get(o);
      if (!h || h.seq !== a) return;
      const y = typeof m == "string" ? m : `<pre class="whitespace-pre-wrap text-sm">${q(
        JSON.stringify(m, null, 2)
      )}</pre>`;
      Y(f, y, { trusted: !1, alpine: t }), f.setAttribute("data-tab-loaded", "true");
    } catch (m) {
      const h = this._pending.get(o);
      if (!h || h.seq !== a || m?.name === "AbortError") return;
      const y = m?.message || "Unknown error";
      Y(f, `<div class="${v}">${q(r ? "Error loading content via HTMX: " : "Error loading content: ")}${q(y)}</div>`, {
        trusted: !1,
        alpine: t
      });
    } finally {
      const m = this._pending.get(o);
      m && m.seq === a && this._pending.delete(o);
    }
  },
  trigger(o, l = null, e = null, r = null) {
    e && typeof window < "u" && window.location.hash.slice(1) === e && (this.activeTab || (this.activeTab = o));
    const u = g.getStyle("tab", "active", "border-b-2 border-blue-500"), f = g.getStyle("tab", "inactive", "text-slate-500 hover:text-slate-700"), d = g.getStyle("tab", "triggerActive", u), p = g.getStyle("tab", "triggerInactive", f);
    return {
      "@click.prevent": `open('${o}', ${l ? `'${l}'` : "null"}, ${e ? `'${e}'` : "null"}, ${r ? `'${r}'` : "null"})`,
      ":class": `isActive('${o}') ? '${d}' : '${p}'`,
      role: "tab",
      ":aria-selected": `isActive('${o}') ? 'true' : 'false'`,
      "data-tab-target": o,
      "data-tab-hash": e || null,
      "data-tab-url": l || null,
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
  panel(o) {
    return {
      "x-show": `isActive('${o}')`,
      role: "tabpanel",
      id: o,
      tabindex: "0"
    };
  }
});
typeof window < "u" && (window.unoTabs = tt);
const Wt = (n) => {
  ct = n, n.data("unoTabs", tt), g.tab = {
    open: (i) => {
      const o = document.querySelector(`[data-tab-target="${i}"]`);
      o ? o.click() : window.location.hash = i;
    }
  };
};
g.dom = J;
typeof window < "u" && (window.Alpine = window.Alpine || X, window.htmx = window.htmx || gt, window.unoTabs = window.unoTabs || tt, window.Uno = g);
g.use(wt).use(St).use(Tt).use(Ht).use(At).use(Wt);
function Z() {
  if (g.initialized) return;
  const n = window.UnoConfig || {};
  g.init(window.Alpine, n);
}
if (typeof window < "u") {
  const n = window.deferLoadingAlpine;
  window.deferLoadingAlpine = (i) => {
    Z(), typeof n == "function" ? n(i) : i();
  }, document.addEventListener("alpine:init", () => {
    g.initialized || Z();
  }), document.addEventListener("DOMContentLoaded", () => {
    !g.initialized && window.Alpine && Z();
  }), !g.initialized && window.Alpine && Z();
}
export {
  g as default
};
