import j from "alpinejs";
import O from "htmx.org";
const D = {
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
}, W = {
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
}, Y = {
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
const U = {
  en: D,
  "zh-CN": W,
  "zh-TW": Y
}, _ = Object.keys(U);
function Z() {
  if (typeof localStorage < "u") {
    const r = localStorage.getItem("uno.lang");
    if (r && _.includes(r)) return r;
  }
  if (typeof navigator < "u") {
    const r = navigator.language, s = r.startsWith("zh") ? r.includes("TW") || r.includes("Hant") ? "zh-TW" : "zh-CN" : "en";
    return _.includes(s) ? s : "en";
  }
  return "en";
}
function V(r, s) {
  const e = s.split(".");
  let t = r;
  for (const i of e) {
    if (!t || typeof t != "object" || !(i in t)) return;
    t = t[i];
  }
  return t;
}
function N(r, s) {
  for (const e in s)
    s[e] && typeof s[e] == "object" ? (r[e] || Object.assign(r, { [e]: {} }), N(r[e], s[e])) : Object.assign(r, { [e]: s[e] });
}
function F(r) {
  r.store("i18n", {
    locale: Z(),
    messages: U,
    init() {
      typeof document < "u" && (document.documentElement.lang = this.locale);
    },
    t(s, e = {}) {
      const t = this.locale, i = this.messages[t] || this.messages.en, f = V(i, s);
      return (typeof f == "string" ? f : s).replace(/\{(\w+)\}/g, (d, l) => l in e ? String(e[l]) : `{${l}}`);
    },
    setLang(s) {
      if (!this.messages[s]) {
        console.warn(`[Uno] Language ${s} not found`);
        return;
      }
      this.locale = s, typeof localStorage < "u" && localStorage.setItem("uno.lang", s), typeof document < "u" && (document.documentElement.lang = s);
    },
    extend(s) {
      for (const e in s)
        this.messages[e] ? N(this.messages[e], s[e]) : this.messages[e] = s[e];
    }
  });
}
const E = {
  /**
   * Log info message
   * @param {...any} args - Arguments to log
   */
  info(...r) {
    console.log("[Uno Info]", ...r);
  },
  /**
   * Log warning message
   * @param {...any} args - Arguments to log
   */
  warn(...r) {
    console.warn("[Uno Warn]", ...r);
  },
  /**
   * Log error message
   * @param {...any} args - Arguments to log
   */
  error(...r) {
    console.error("[Uno Error]", ...r);
  }
};
function J(r, s, e) {
  const t = r.getBoundingClientRect(), i = s.getBoundingClientRect(), f = i.width, g = i.height, d = t.width, l = t.height, n = { top: t.top, left: t.left }, p = window.innerWidth, m = window.innerHeight, h = 0, y = 0, b = m, x = p, k = (o, u) => {
    if (u) {
      if (o + f > x)
        return x - f - 10;
      if (o < y)
        return y + 10;
    } else {
      if (o + g > b)
        return b - g - 10;
      if (o < h)
        return h + 10;
    }
    return o;
  };
  let L, c, a = e;
  switch (e) {
    case "left":
      c = n.left - f - 5, c < y && (c = n.left + d + 5, a = "right"), L = k(n.top + l / 2 - g / 2, !1);
      break;
    case "top":
      L = n.top - g - 8, L < h && (L = n.top + l + 8, a = "bottom"), c = n.left + (d - f) / 2, c = k(c, !0);
      break;
    case "right":
      c = n.left + d + 5, c + f > x && (c = n.left - f - 5, a = "left"), L = k(n.top + l / 2 - g / 2, !1);
      break;
    case "bottom":
    default:
      L = n.top + l + 8, L + g > b && (L = n.top - g - 8, a = "top"), c = n.left + (d - f) / 2, c = k(c, !0);
      break;
  }
  return { top: L, left: c, position: a };
}
const K = (r) => {
  if (!r) {
    E.error("Tooltip: Alpine is undefined. Plugin initialization failed.");
    return;
  }
  r.store("tooltip", {
    visible: !1,
    content: "",
    placement: "top",
    top: -9999,
    // Start off-screen
    left: -9999,
    theme: "dark",
    // 'dark' | 'light'
    triggerEl: null,
    async show(e, t, i = {}) {
      this.triggerEl = e, this.content = t, this.placement = i.placement || "top", this.theme = i.theme || "dark", this.visible = !0, this.ensureContainer(), await r.nextTick(), this.updatePosition();
    },
    hide() {
      this.visible = !1, this.triggerEl = null;
    },
    updatePosition() {
      if (!this.triggerEl) return;
      const e = document.getElementById("uno-tooltip-content");
      if (!e) return;
      const t = J(this.triggerEl, e, this.placement);
      this.top = t.top, this.left = t.left;
    },
    ensureContainer() {
      document.getElementById("uno-tooltip-container") || s();
    }
  });
  const s = () => {
    if (!document.getElementById("uno-tooltip-container"))
      try {
        const e = document.createElement("div");
        e.id = "uno-tooltip-container", e.setAttribute("x-data", "{}"), e.innerHTML = `
          <div id="uno-tooltip-content"
               x-show="$store.tooltip.visible"
               x-transition:enter="transition ease-out duration-200"
               x-transition:enter-start="opacity-0 scale-95"
               x-transition:enter-end="opacity-100 scale-100"
               x-transition:leave="transition ease-in duration-150"
               x-transition:leave-start="opacity-100 scale-100"
               x-transition:leave-end="opacity-0 scale-95"
               class="fixed z-[9999] px-3 py-2 text-sm font-medium rounded-lg shadow-sm whitespace-nowrap pointer-events-none border"
               :class="$store.tooltip.theme === 'light' ? 'bg-white text-gray-900 border-gray-200' : 'bg-slate-800 text-white border-slate-700'"
               :style="\`top: \${$store.tooltip.top}px; left: \${$store.tooltip.left}px\`"
               style="display: none;">
               <span x-html="$store.tooltip.content"></span>
          </div>
        `, document.body ? document.body.appendChild(e) : document.documentElement.appendChild(e);
      } catch (e) {
        E.error("Tooltip: Failed to inject container", e);
      }
  };
  document.body ? s() : document.addEventListener("DOMContentLoaded", s), r.directive(
    "tooltip",
    (e, { expression: t, modifiers: i }, { evaluateLater: f, effect: g, cleanup: d }) => {
      const l = f(t);
      let n = {
        content: "",
        placement: i[0] || "top",
        trigger: i.includes("click") ? "click" : "hover",
        theme: "dark"
      };
      i.includes("top") && (n.placement = "top"), i.includes("bottom") && (n.placement = "bottom"), i.includes("left") && (n.placement = "left"), i.includes("right") && (n.placement = "right"), i.includes("light") && (n.theme = "light");
      const p = (b) => {
        if (typeof b == "string" ? n.content = b : typeof b == "object" && (n = { ...n, ...b }), n.content && n.content.startsWith("#") && !n.content.includes(" ")) {
          const x = document.querySelector(n.content);
          x && (n.content = x.innerHTML);
        } else if (n.target) {
          const x = document.querySelector(n.target);
          x && (n.content = x.innerHTML);
        }
      };
      g(() => {
        l((b) => {
          p(b);
        });
      });
      const m = () => r.store("tooltip").show(e, n.content, n), h = () => r.store("tooltip").hide(), y = () => {
        const b = r.store("tooltip");
        b.visible && b.triggerEl === e ? b.hide() : m();
      };
      if (n.trigger === "click") {
        e.addEventListener("click", y);
        const b = (x) => {
          const k = r.store("tooltip");
          k.visible && k.triggerEl === e && !e.contains(x.target) && k.hide();
        };
        window.addEventListener("click", b), d(() => window.removeEventListener("click", b));
      } else
        e.addEventListener("mouseenter", m), e.addEventListener("mouseleave", h), e.addEventListener("focus", m), e.addEventListener("blur", h);
      d(() => {
        n.trigger === "click" ? e.removeEventListener("click", y) : (e.removeEventListener("mouseenter", m), e.removeEventListener("mouseleave", h), e.removeEventListener("focus", m), e.removeEventListener("blur", h));
      });
    }
  );
};
function P(r) {
  const s = { position: "fixed" };
  switch (r) {
    case "top":
      s.top = "0px", s.left = "0px", s.right = "0px";
      break;
    case "bottom":
      s.bottom = "0px", s.left = "0px", s.right = "0px";
      break;
    case "top-right":
      s.top = "20px", s.right = "20px";
      break;
    case "bottom-right":
      s.bottom = "12px", s.right = "20px";
      break;
    case "top-left":
      s.top = "20px", s.left = "20px";
      break;
    case "bottom-left":
      s.bottom = "12px", s.left = "20px";
      break;
    case "center":
      s.top = "50%", s.left = "50%", s.transform = "translate(-50%, -50%)";
      break;
    case "top-center":
      s.top = "0px", s.left = "50%", s.transform = "translateX(-50%)";
      break;
    case "bottom-center":
      s.bottom = "0px", s.left = "50%", s.transform = "translateX(-50%)";
      break;
    default:
      s.top = "20px", s.right = "20px";
  }
  return s;
}
let A = null, M = !1;
function X(r, s = 1e3) {
  if (!r) return;
  M = !0;
  const e = () => {
    M && (r.style.transition = "none", r.style.width = "0%", r.offsetHeight, setTimeout(() => {
      M && (r.style.transition = `width ${s}ms ease-in-out`, r.style.width = "100%", A = setTimeout(e, s + 200));
    }, 50));
  };
  e();
}
function G(r) {
  M = !1, A && (clearTimeout(A), A = null), r && (r.style.width = "0%");
}
const Q = (r) => {
  r.data("notify", ({ position: s = "top-right" } = {}) => ({
    items: [],
    position: s,
    nextId: 1,
    showCloseButton: !1,
    closeIndicatorEnabled: !1,
    activeRequests: 0,
    progressBarRunning: !1,
    get containerClasses() {
      switch (this.position) {
        case "top-left":
          return "flex flex-col space-y-2 items-start";
        case "bottom-right":
          return "flex flex-col-reverse space-y-2 space-y-reverse items-end";
        case "bottom-left":
          return "flex flex-col-reverse space-y-2 space-y-reverse items-start";
        case "top-right":
        default:
          return "flex flex-col space-y-2 items-end";
      }
    },
    get containerStyle() {
      return P(this.position);
    },
    itemClasses(e) {
      const t = "relative border-l-4 pl-3 pr-3 py-4 mb-2 rounded shadow-sm dark:shadow-md transition-colors duration-300 overflow-hidden animate-fadeIn", i = {
        success: "bg-green-100 dark:bg-green-800 border-green-500 dark:border-green-400 text-green-700 dark:text-green-100",
        warning: "bg-yellow-100 dark:bg-yellow-800 border-yellow-500 dark:border-yellow-400 text-yellow-700 dark:text-yellow-100",
        error: "bg-red-100 dark:bg-red-800 border-red-500 dark:border-red-400 text-red-700 dark:text-red-100",
        info: "bg-blue-100 dark:bg-blue-800 border-blue-500 dark:border-blue-400 text-blue-700 dark:text-blue-100"
      }, f = i[e] || i.error;
      return `${t} ${f}`;
    },
    /**
     * Show a notification
     * @param {string} text - Message text
     * @param {Object} options - Options { type, timeout, position, closeIndicator, showCloseButton }
     */
    show(e, { type: t = "error", timeout: i, position: f, closeIndicator: g, showCloseButton: d } = {}) {
      const l = this.nextId++, n = f || this.position;
      E.info("notify.show", { text: e, type: t, timeout: i, position: n });
      const p = {
        id: l,
        text: e,
        type: t,
        targetPosition: n,
        showCloseButton: !!d,
        closeIndicator: !!g
      };
      this.items = [p, ...this.items];
      let m = i;
      if (m === void 0)
        switch (t) {
          case "success":
            m = 1e3;
            break;
          case "warning":
            m = 2e3;
            break;
          case "error":
          default:
            m = 3e3;
            break;
        }
      if (m && m > 0) {
        const h = l;
        setTimeout(() => this.remove(h), m);
      }
    },
    /**
     * Remove a notification by ID
     * @param {number} id
     */
    remove(e) {
      this.items = this.items.filter((t) => t.id !== e);
    },
    success(e, t) {
      const i = e || this.$store?.i18n?.t("notify.operation.success") || "Success";
      typeof t == "number" ? this.show(i, { type: "success", timeout: t }) : this.show(i, { type: "success", ...t || {} });
    },
    warning(e, t) {
      typeof t == "number" ? this.show(e, { type: "warning", timeout: t }) : this.show(e, { type: "warning", ...t || {} });
    },
    error(e, t) {
      typeof t == "number" ? this.show(e, { type: "error", timeout: t }) : this.show(e, { type: "error", ...t || {} });
    },
    /**
     * Show or hide loading indicator
     * @param {boolean} show
     * @param {string|Object} styleOrOptions - 'indicator' | 'progressbar' | { style, position }
     */
    loading(e, t) {
      let i = "indicator", f = null;
      typeof t == "string" ? i = t : typeof t == "object" && t && (i = t.style || "indicator", f = t.position), e ? (this.activeRequests++, i === "progressbar" ? (this._ensureProgressBar(), this._showProgressBar()) : (this._ensureLoadingIndicator(), this._showLoadingIndicator(f))) : (this.activeRequests = Math.max(0, this.activeRequests - 1), this.activeRequests === 0 && (this._hideLoadingIndicator(), this._hideProgressBar()));
    },
    _ensureLoadingIndicator() {
      if (!document.getElementById("uno-loading-indicator")) {
        const e = document.createElement("div");
        e.id = "uno-loading-indicator", e.className = "fixed top-4 right-4 z-[150] bg-blue-100 border-blue-500 text-blue-700 p-2 rounded-full shadow-md flex items-center hidden dark:bg-blue-900 dark:text-blue-100", e.innerHTML = '<svg class="animate-spin -ml-1 mr-2 h-5 w-5 text-blue-500 dark:text-blue-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg><span>Loading</span>', document.body.appendChild(e);
      }
    },
    _showLoadingIndicator() {
      const e = document.getElementById("uno-loading-indicator");
      e && (e.style.display = "flex");
    },
    _hideLoadingIndicator() {
      const e = document.getElementById("uno-loading-indicator");
      e && (e.style.display = "none");
    },
    _progressFill() {
      const e = document.getElementById("uno-progress-bar");
      return e ? e.querySelector(".uno-progress-fill") : null;
    },
    _ensureProgressBar() {
      if (!document.getElementById("uno-progress-bar")) {
        const e = document.createElement("div");
        e.id = "uno-progress-bar", e.className = "fixed top-0 left-0 right-0 z-[9999] bg-blue-100 dark:bg-blue-900 hidden";
        const t = document.createElement("div");
        t.className = "uno-progress-fill bg-blue-500 dark:bg-blue-400 transition-all duration-1000 ease-in-out", t.style.height = "4px", t.style.width = "0%", e.appendChild(t), document.body.appendChild(e);
      }
    },
    _showProgressBar() {
      const e = document.getElementById("uno-progress-bar"), t = this._progressFill();
      !e || !t || (e.style.display = "block", this.progressBarRunning = !0, X(t, 1e3));
    },
    _hideProgressBar() {
      const e = document.getElementById("uno-progress-bar"), t = this._progressFill();
      !e || !t || (e.style.display = "none", this.progressBarRunning = !1, G(t));
    }
  }));
};
let S = 0, B = "";
function R() {
  S === 0 && (B = document.body.style.overflow || "", document.body.style.overflow = "hidden"), S++;
}
function tt() {
  S = Math.max(0, S - 1), S === 0 && (document.body.style.overflow = B || "", B = "");
}
async function et(r) {
  const s = await fetch(r);
  return (s.headers.get("content-type") || "").includes("application/json") ? s.json() : s.text();
}
const nt = (r) => {
  const s = {};
  let e = "uno-box-";
  function t() {
    return Math.random().toString(36).slice(2, 10);
  }
  function i(c, a) {
    return `#${e}${a}${c}`;
  }
  function f(c) {
    const a = c.dataset, o = {
      boxTitle: "title",
      boxContent: "content",
      boxWidth: "width",
      boxHeight: "height",
      boxModal: "modal",
      boxPosition: "position",
      boxAjaxUrl: "ajaxUrl",
      boxIframeUrl: "iframeUrl",
      boxOkText: "okButtonText",
      boxCancelText: "cancelButtonText",
      boxCache: "cache",
      boxId: "boxId",
      boxWidthClass: "widthClass",
      boxHeightClass: "heightClass",
      boxContainerClass: "containerClass",
      boxContentOverflow: "contentOverflow",
      boxButtonAlign: "buttonAlign"
    }, u = {};
    for (const $ in o) {
      const w = a[$];
      if (w !== void 0) {
        const C = o[$];
        C === "width" || C === "height" ? u[C] = parseInt(w, 10) : C === "modal" || C === "cache" ? u[C] = w === "true" || w === !0 : u[C] = w;
      }
    }
    return u;
  }
  function g(c, a) {
    try {
      const o = r?.store("i18n");
      return o && o.t(`box.${c}`) || a;
    } catch {
      return a;
    }
  }
  function d(c) {
    const a = c.uniqueId, o = !!c.title, u = (c.containerClass || "").trim(), $ = o ? `
      <div id="${e}header${a}" class="flex justify-between items-center p-4 border-b border-gray-300 dark:border-gray-600">
        <div id="${e}title${a}" class="text-lg font-semibold text-gray-900 dark:text-gray-100"></div>
        <button id="${e}close${a}" class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200" title="${g("close", "Close")}">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
        </button>
      </div>
    ` : "", w = o ? "" : `
      <button id="${e}close${a}" class="absolute top-2 right-2 z-[101] text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 bg-white dark:bg-gray-900 p-1" title="Close">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
      </button>
    `, C = o ? "p-4 text-gray-900 dark:text-gray-100" : "pt-12 px-4 pb-4 text-gray-900 dark:text-gray-100", I = document.createElement("div");
    I.id = `${e}outer-container${a}`, I.className = "hidden fixed inset-0 z-[103] p-4 flex items-center justify-center";
    const q = document.createElement("div");
    q.id = `${e}container${a}`, q.className = `bg-white dark:bg-gray-900 border-8 border-gray-200 dark:border-gray-600 shadow-lg overflow-hidden transition-colors duration-300 ${u}`, q.style.position = "relative", q.innerHTML = `${$}${w}<div id="${e}wrapper${a}" class="${C}"><div id="${e}content${a}" class="overflow-auto"></div></div><div id="${e}footer${a}" class="p-4 border-t border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-100 hidden"></div>`, I.appendChild(q);
    const v = document.createElement("div");
    return v.id = `${e}overlay${a}`, v.className = "hidden fixed inset-0 bg-black z-[102]", v.style.opacity = String(c.overlayOpacity ?? 0.5), document.body.appendChild(I), document.body.appendChild(v), { outer: I, container: q, overlay: v };
  }
  function l(c, a) {
    const o = document.querySelector(i(c, "container")), u = document.querySelector(i(c, "content")), $ = document.querySelector(i(c, "wrapper"));
    if (!o || !u) return;
    a.widthClass || a.heightClass ? (o.style.width = "", o.style.height = "", o.style.maxWidth = "", o.style.maxHeight = "", a.widthClass && o.classList.add(...a.widthClass.split(" ")), a.heightClass && (o.classList.add(...a.heightClass.split(" ")), $ && ($.classList.add("flex", "flex-col", "h-full"), u.classList.add("flex-1", "min-h-0")))) : window.innerWidth <= 768 ? (o.style.width = "calc(100% - 2rem)", o.style.height = "calc(100% - 2rem)", $ ? ($.classList.add("flex", "flex-col", "h-full"), u.classList.add("flex-1", "min-h-0")) : u.style.height = "100%") : (a.width && (o.style.width = `${a.width}px`), a.height && (u.style.height = `${a.height}px`)), u.classList.remove("overflow-auto", "overflow-hidden", "overflow-visible");
    const w = a.contentOverflow || "auto";
    w === "hidden" ? u.classList.add("overflow-hidden") : w === "visible" ? u.classList.add("overflow-visible") : u.classList.add("overflow-auto");
  }
  function n(c, a) {
    const o = document.querySelector(i(c, "outer-container"));
    if (!o) return;
    const u = P(a || "center");
    Object.assign(o.style, u), o.classList.remove("hidden");
  }
  function p(c, a) {
    const o = document.querySelector(i(c, "close"));
    o && o.addEventListener("click", () => b({ uniqueId: c, modal: a }));
  }
  async function m(c, a) {
    const o = document.querySelector(i(c, "content"));
    if (o)
      if (typeof a.content == "string")
        o.innerHTML = a.content;
      else if (typeof a.targetId == "string") {
        const u = document.getElementById(a.targetId);
        u && (o.innerHTML = u.innerHTML);
      } else if (typeof a.ajaxUrl == "string") {
        o.innerHTML = `<div class="flex justify-center items-center h-full"><div class="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 dark:border-gray-100"></div><span class="ml-2 text-gray-900 dark:text-gray-100">${g("loading", "Loading")}</span></div>`;
        try {
          const u = await et(a.ajaxUrl);
          typeof u == "object" && u !== null ? o.innerHTML = `<pre class="bg-gray-100 p-3 rounded overflow-x-auto whitespace-pre-wrap">${JSON.stringify(u, null, 2)}</pre>` : o.innerHTML = String(u);
        } catch (u) {
          E.error("box ajax error", u);
        }
      } else typeof a.iframeUrl == "string" && (o.innerHTML = `<iframe frameborder="0" class="w-full h-full" src="${a.iframeUrl}"></iframe>`);
  }
  function h(c, a, o) {
    const u = document.querySelector(i(c, "footer"));
    if (!u) return;
    const $ = { left: "justify-start", center: "justify-center", right: "justify-end" }[o.buttonAlign] || "justify-end";
    u.classList.remove("hidden"), u.className = `p-4 border-t border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-100 flex ${$} space-x-4`;
    const w = o.okButtonText || g("ok", "OK"), C = o.cancelButtonText || g("cancel", "Cancel"), I = a === "confirm" || a === "message" || a === "dialog" && o.showOkButton === !0, q = a === "confirm" || a === "dialog" && o.showCancelButton === !0;
    if (u.innerHTML = "", I) {
      const v = document.createElement("button");
      v.id = `${e}ok${c}`, v.className = "px-4 py-2 bg-blue-500 text-white hover:bg-blue-600 transition-colors duration-200", v.textContent = w, v.addEventListener("click", () => {
        let z = !0;
        typeof o.onOk == "function" && (z = o.onOk() !== !1), z && b({ uniqueId: c, modal: o.modal !== !1 });
      }), u.appendChild(v);
    }
    if (q) {
      const v = document.createElement("button");
      v.id = `${e}cancel${c}`, v.className = "px-4 py-2 bg-gray-300 dark:bg-gray-600 text-gray-700 dark:text-gray-200 hover:bg-gray-400 dark:hover:bg-gray-500 transition-colors duration-200", v.textContent = C, v.addEventListener("click", () => {
        b({ uniqueId: c, modal: o.modal !== !1 }), typeof o.onCancel == "function" && o.onCancel();
      }), u.appendChild(v);
    }
    !I && !q && u.classList.add("hidden");
  }
  async function y(c, a) {
    const o = {
      modal: !0,
      position: "center",
      width: 400,
      height: void 0,
      overlayOpacity: 0.5,
      responsive: !0,
      contentOverflow: "auto",
      ...a
    };
    typeof o.onBeforeOpen == "function" && o.onBeforeOpen();
    let u = o.boxId || t();
    if (o.uniqueId = u, o.cache && o.boxId && s[o.boxId])
      return u = s[o.boxId].uniqueId, document.querySelector(i(u, "outer-container"))?.classList.remove("hidden"), document.querySelector(i(u, "overlay"))?.classList.toggle("hidden", o.modal === !1), o.modal !== !1 && R(), h(u, c, o), p(u, o.modal !== !1), typeof o.onOpen == "function" && o.onOpen(), u;
    const $ = d(o);
    if (o.cache && o.boxId && (s[o.boxId] = { uniqueId: u, modal: o.modal !== !1 }), o.title) {
      const w = document.querySelector(i(u, "title"));
      w && (w.textContent = o.title);
    }
    return o.modal !== !1 && ($.overlay.classList.remove("hidden"), R()), n(u, o.position), l(u, o), await m(u, o), h(u, c, o), p(u, o.modal !== !1), typeof o.onOpen == "function" && o.onOpen(), u;
  }
  function b({ uniqueId: c, modal: a }) {
    document.querySelector(i(c, "outer-container"))?.classList.add("hidden"), document.querySelector(i(c, "overlay"))?.classList.add("hidden"), a && tt();
  }
  function x() {
    Object.keys(s).forEach((c) => {
      const a = s[c].uniqueId;
      document.querySelector(i(a, "outer-container"))?.remove(), document.querySelector(i(a, "overlay"))?.remove(), delete s[c];
    });
  }
  function k(c) {
    const a = s[c];
    if (!a) return;
    const o = a.uniqueId;
    document.querySelector(i(o, "outer-container"))?.remove(), document.querySelector(i(o, "overlay"))?.remove(), delete s[c];
  }
  function L() {
    document.addEventListener("click", (c) => {
      const a = c.target.closest("[data-box-type]");
      if (!a) return;
      c.preventDefault();
      const o = a.getAttribute("data-box-type");
      if (!["dialog", "confirm", "message"].includes(o)) return;
      const u = f(a);
      y(o, u);
    });
  }
  L(), r.store("box", {
    open: y,
    close: b,
    clearCache: x,
    clearCacheById: k
  });
}, H = {
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
function it(r) {
  const s = r.match(/^\/(.*?)\/([gimsuy]*)$/);
  return s ? { pattern: s[1], flags: s[2] } : { pattern: r, flags: "" };
}
function ot(r) {
  try {
    return r ? JSON.parse(r) : [];
  } catch {
    return [];
  }
}
const rt = (r) => {
  const s = {
    invalid: `url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 12 12' width='12' height='12' fill='none' stroke='%23ef4444'%3e%3ccircle cx='6' cy='6' r='4.5'/%3e%3cpath stroke-linecap='round' d='M5.8 3.6h.4L6 6.5zM6 8.2af.6.6 0 110-1.2.6.6 0 010 1.2z'/%3e%3c/svg%3e")`,
    valid: `url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 8 8'%3e%3cpath fill='%2322c55e' d='M2.3 6.73L.6 4.53c-.4-1.04.46-1.4 1.1-.8l1.1 1.4 3.4-3.8c.6-.63 1.6-.27 1.2.7l-4 4.6c-.43.5-.8.4-1.1.1z'/%3e%3c/svg%3e")`
  }, e = (t, i) => {
    const f = r.store("i18n");
    if (!f) return t;
    const g = {
      required: "Required",
      pattern: "Invalid format",
      min: `At least ${i?.n} characters`,
      max: `Max ${i?.n} characters`,
      int: "Integer required",
      email: "Invalid email address",
      url: "Invalid URL"
    };
    return f.t(`valid.${t}`, i) || g[t] || t;
  };
  r.store("valid", {
    callbacks: {},
    registerCallback(t, i) {
      typeof i == "function" && (this.callbacks[t] = i);
    }
  }), r.data("validForm", () => ({
    fields: {},
    // Reactive state: { [name]: { value, error, valid, ... } }
    init() {
      this.$el.noValidate = !0;
    },
    // Register field (called by x-validate directive)
    register(t, i, f = []) {
      if (this.fields[t]) return;
      const g = [...f];
      if ((i.hasAttribute("required") || i.dataset.validRequired === "true") && g.push({ type: "required" }), (i.type === "email" || i.dataset.validType === "email") && g.push("email"), (i.type === "url" || i.dataset.validType === "url") && g.push("url"), i.dataset.validType) {
        const d = i.dataset.validType;
        d !== "email" && d !== "url" && g.push(d);
      }
      i.dataset.validPattern && g.push({ type: "regexp", expression: i.dataset.validPattern }), this.fields[t] = {
        el: i,
        rules: g,
        value: i.value,
        error: null,
        valid: !0,
        dirty: !1,
        touched: !1,
        validating: !1
      };
    },
    update(t, i) {
      this.fields[t] && (this.fields[t].value = i, this.fields[t].dirty = !0, this.validate(t));
    },
    touch(t) {
      this.fields[t] && (this.fields[t].touched = !0, this.validate(t));
    },
    async validate(t) {
      const i = this.fields[t];
      if (!i) return !0;
      i.validating = !0, i.error = null;
      const f = (i.value || "").trim();
      let g = null;
      for (const d of i.rules) {
        const l = d.type || d, n = typeof d == "object" ? d : {};
        if (!f) {
          if (l === "required") {
            g = n.msg || e("required");
            break;
          }
          if (l === "emptyquit") break;
          continue;
        }
        if (H[l])
          H[l].test(f) || (g = n.msg || e(l) || e("pattern"));
        else if (l === "regexp" && n.expression) {
          const { pattern: p, flags: m } = it(n.expression);
          new RegExp(p, m || "i").test(f) || (g = n.msg || e("pattern"));
        } else if (l === "limit")
          n.min && f.length < n.min && (g = n.msg || e("min", { n: n.min })), n.max && f.length > n.max && (g = n.msg || e("max", { n: n.max }));
        else if (l === "range") {
          const p = parseFloat(f);
          n.min !== void 0 && p < n.min && (g = n.msg || `Min ${n.min}`), n.max !== void 0 && p > n.max && (g = n.msg || `Max ${n.max}`);
        } else if (l === "compareto") {
          const p = document.getElementById(n.destid);
          p && f !== p.value && (g = n.msg || "Not match");
        } else if (l === "tag") {
          const p = n.mark || ",", m = f.split(p).filter((h) => h.trim());
          n.max && m.length > n.max && (g = n.msg || `Max ${n.max} items`);
        }
        if (g) break;
        if (l === "callback" && n.fn) {
          const p = r.store("valid").callbacks[n.fn];
          if (p)
            try {
              const m = await p(i.el.id, f);
              m === !1 || typeof m == "object" && m.valid === !1 ? g = (typeof m == "object" ? m.message : "") || "Invalid" : typeof m == "string" && (g = m);
            } catch {
            }
        }
        if (g) break;
      }
      return i.error = g, i.valid = !g, i.validating = !1, i.valid;
    },
    async validateAll() {
      const t = Object.keys(this.fields);
      return (await Promise.all(t.map((f) => this.validate(f)))).every((f) => f);
    },
    async submit(t) {
      if (t && (t.preventDefault(), t.stopPropagation()), Object.keys(this.fields).forEach((f) => {
        this.fields[f].touched = !0;
      }), await this.validateAll())
        this.$el.dispatchEvent(new CustomEvent("valid-submit", { bubbles: !0 })), !this.$el.hasAttribute("hx-post") && !this.$el.hasAttribute("hx-get") && this.$el.submit();
      else {
        const f = Object.keys(this.fields).find((g) => !this.fields[g].valid);
        f && this.fields[f].el && this.fields[f].el.focus();
      }
    }
  })), r.directive("validate", (t, i, { effect: f, cleanup: g }) => {
    const d = t.getAttribute("name");
    if (!d) return;
    let l = [];
    const n = t.getAttribute("valid");
    n && (l = ot(n)), f(() => {
      const h = r.$data(t);
      h && typeof h.register == "function" && h.register(d, t, l);
    });
    const p = () => {
      const h = r.$data(t);
      if (h && h.update && h.update(d, t.value), l.find((y) => (y.type === "limit" || y === "limit") && y.show === !0)) {
        const y = l.find((b) => b.type === "limit" || b === "limit");
        y && y.max && t.value.length > y.max && (t.value = t.value.substring(0, y.max), h && h.update && h.update(d, t.value));
      }
    }, m = () => {
      const h = r.$data(t);
      h && h.touch && h.touch(d);
    };
    t.addEventListener("input", p), t.addEventListener("blur", m), f(() => {
      const h = r.$data(t);
      if (h && h.fields && h.fields[d]) {
        const y = h.fields[d], b = !!y.error, x = y.valid && y.dirty && !y.error;
        t.classList.toggle("is-invalid", b), t.classList.toggle("is-valid", x), t.classList.toggle("validating", y.validating), b ? (t.classList.add(
          "!border-red-500",
          "pr-10",
          "bg-no-repeat",
          "bg-[right_0.75rem_center]",
          "bg-[length:1rem_1rem]"
        ), t.classList.remove("!border-green-500"), t.style.backgroundImage = s.invalid) : x ? (t.classList.add(
          "!border-green-500",
          "pr-10",
          "bg-no-repeat",
          "bg-[right_0.75rem_center]",
          "bg-[length:1rem_1rem]"
        ), t.classList.remove("!border-red-500"), t.style.backgroundImage = s.valid) : (t.classList.remove(
          "!border-red-500",
          "!border-green-500",
          "pr-10",
          "bg-no-repeat",
          "bg-[right_0.75rem_center]",
          "bg-[length:1rem_1rem]"
        ), t.style.backgroundImage = "");
        const k = t.nextElementSibling;
        if (k && k.tagName === "DIV") {
          const L = k.querySelector(".animate-spin") || t.parentNode.querySelector(".animate-spin");
          L && (L.style.display = y.validating ? "block" : "");
        }
      }
    }), g(() => {
      t.removeEventListener("input", p), t.removeEventListener("blur", m);
    });
  }), r.directive("valid-error", (t, { expression: i }, { effect: f }) => {
    const g = i;
    t.style.display = "none", f(() => {
      const d = r.$data(t);
      if (d && d.fields && d.fields[g]) {
        const n = d.fields[g].error;
        n ? (t.textContent = n, t.style.display = "block", t.classList.add("text-red-600", "text-xs"), t.classList.remove("text-gray-400")) : t.style.display = "none";
      }
    });
  }), r.directive("valid-count", (t, { expression: i }, { effect: f }) => {
    const g = i;
    f(() => {
      const d = r.$data(t);
      if (d && d.fields && d.fields[g]) {
        const l = d.fields[g], n = l.rules.find((p) => p.type === "limit" || p === "limit");
        if (n && n.max) {
          const p = (l.value || "").length, m = n.max - p, h = r.store("i18n");
          t.textContent = h && h.t("charCount", { count: m }) || `Remaining: ${m}`, t.style.display = "block";
        } else
          t.style.display = "none";
      }
    });
  });
};
const st = (r) => {
  const s = {
    activeClass: "text-blue-600 dark:text-blue-400 border-b-2 border-blue-500 dark:border-blue-400",
    inactiveClass: "text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200",
    baseTabClass: "px-4 py-2 transition-colors duration-150 cursor-pointer"
  }, e = (d, l) => {
    const n = s.activeClass.split(" "), p = s.inactiveClass.split(" ");
    l ? (d.classList.add("current"), d.classList.remove(...p), d.classList.add(...n)) : (d.classList.remove("current"), d.classList.remove(...n), d.classList.add(...p));
  }, t = (d, l) => {
    const n = document.querySelectorAll(`[data-tab-group="${d}"], [tabgroup="${d}"]`);
    n.forEach((m) => {
      const y = (m.getAttribute("data-tab-target") || m.getAttribute("tabtarget")) === l;
      e(m, y);
    });
    const p = /* @__PURE__ */ new Set();
    n.forEach((m) => {
      const h = m.getAttribute("data-tab-target") || m.getAttribute("tabtarget");
      h && p.add(h);
    }), p.forEach((m) => {
      const h = document.getElementById(m);
      h && (m === l ? (h.classList.remove("hidden"), h.style.display = "") : h.classList.add("hidden"));
    });
  }, i = () => {
    const d = window.location.hash.substring(1);
    if (!d) return;
    const l = document.querySelector(`[data-tab-hash="${d}"], [rel="${d}"]`);
    if (l) {
      const n = l.getAttribute("data-tab-group") || l.getAttribute("tabgroup"), p = l.getAttribute("data-tab-target") || l.getAttribute("tabtarget");
      n && p && t(n, p);
    }
  }, f = () => {
    document.querySelectorAll("[data-tab-group], [tabgroup]").forEach((l) => {
      if (!l.className.includes("transition-colors")) {
        const p = s.baseTabClass.split(" ");
        l.classList.add(...p);
      }
      const n = l.classList.contains("current") || l.getAttribute("data-tab-active") === "true";
      e(l, n);
    });
  };
  document.addEventListener("click", (d) => {
    const l = d.target.closest("[data-tab-target], [tabtarget]");
    if (!l) return;
    const n = l.getAttribute("data-tab-group") || l.getAttribute("tabgroup"), p = l.getAttribute("data-tab-target") || l.getAttribute("tabtarget"), m = l.getAttribute("data-tab-hash") || l.getAttribute("rel"), h = l.getAttribute("data-tab-url");
    if (n && p) {
      if (d.preventDefault(), h) {
        const y = document.getElementById(p);
        y && !y.getAttribute("data-tab-loaded") && (y.innerHTML = `<div class="animate-pulse p-4 space-y-2">
             <div class="h-4 bg-gray-200 rounded w-3/4"></div>
             <div class="h-4 bg-gray-200 rounded w-1/2"></div>
           </div>`, fetch(h).then((b) => {
          if (!b.ok) throw new Error("Network response was not ok");
          return b.text();
        }).then((b) => {
          y.innerHTML = b, y.setAttribute("data-tab-loaded", "true");
        }).catch((b) => {
          y.innerHTML = `<div class="text-red-500 p-4">Error loading content: ${b.message}</div>`;
        }));
      }
      t(n, p), m && history.replaceState(null, null, "#" + m);
    }
  });
  const g = () => {
    i(), f(), window.addEventListener("hashchange", i);
  };
  document.readyState === "loading" ? document.addEventListener("DOMContentLoaded", g) : g(), r.store("tabs", {
    switch(d, l) {
      t(d, l);
    }
  });
};
window.Alpine = j;
window.htmx = O;
E.info("Uno: Script loaded. Import Alpine:", !!j, "Window Alpine:", !!window.Alpine);
const T = (r) => {
  if (!r) {
    E.error("Uno: Alpine is undefined, cannot register plugins");
    return;
  }
  if (r.store && r.store("tooltip")) {
    E.info("Uno: Plugins already registered");
    return;
  }
  E.info("Uno: Registering plugins"), F(r), K(r), Q(r), nt(r), rt(r), st(r);
}, at = () => {
  let r = j || window.Alpine;
  r ? (E.info("Uno: Alpine found, registering"), T(r)) : (E.info("Uno: Alpine not found, waiting for alpine:init"), document.addEventListener("alpine:init", () => {
    E.info("Uno: alpine:init fired"), T(window.Alpine);
  }), document.addEventListener("DOMContentLoaded", () => {
    window.Alpine && (E.info("Uno: DOMContentLoaded fallback"), T(window.Alpine));
  }));
};
at();
