import j from "alpinejs";
import z from "htmx.org";
const _ = {
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
}, q = {
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
}, R = {
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
const I = {
  en: _,
  "zh-CN": q,
  "zh-TW": R
}, M = Object.keys(I);
function N() {
  if (typeof localStorage < "u") {
    const s = localStorage.getItem("uno.lang");
    if (s && M.includes(s)) return s;
  }
  if (typeof navigator < "u") {
    const s = navigator.language, l = s.startsWith("zh") ? s.includes("TW") || s.includes("Hant") ? "zh-TW" : "zh-CN" : "en";
    return M.includes(l) ? l : "en";
  }
  return "en";
}
function P(s, l) {
  const e = l.split(".");
  let t = s;
  for (const i of e) {
    if (!t || typeof t != "object" || !(i in t)) return;
    t = t[i];
  }
  return t;
}
function B(s, l) {
  for (const e in l)
    l[e] && typeof l[e] == "object" ? (s[e] || Object.assign(s, { [e]: {} }), B(s[e], l[e])) : Object.assign(s, { [e]: l[e] });
}
function U(s) {
  s.store("i18n", {
    locale: N(),
    messages: I,
    init() {
      typeof document < "u" && (document.documentElement.lang = this.locale);
    },
    t(l, e = {}) {
      const t = this.locale, i = this.messages[t] || this.messages.en, n = P(i, l);
      return (typeof n == "string" ? n : l).replace(/\{(\w+)\}/g, (a, c) => c in e ? String(e[c]) : `{${c}}`);
    },
    setLang(l) {
      if (!this.messages[l]) {
        console.warn(`[Uno] Language ${l} not found`);
        return;
      }
      this.locale = l, typeof localStorage < "u" && localStorage.setItem("uno.lang", l), typeof document < "u" && (document.documentElement.lang = l);
    },
    extend(l) {
      for (const e in l)
        this.messages[e] ? B(this.messages[e], l[e]) : this.messages[e] = l[e];
    }
  });
}
const x = {
  /**
   * Log info message
   * @param {...any} args - Arguments to log
   */
  info(...s) {
    console.log("[Uno Info]", ...s);
  },
  /**
   * Log warning message
   * @param {...any} args - Arguments to log
   */
  warn(...s) {
    console.warn("[Uno Warn]", ...s);
  },
  /**
   * Log error message
   * @param {...any} args - Arguments to log
   */
  error(...s) {
    console.error("[Uno Error]", ...s);
  }
}, W = (s) => {
  const l = s || "center", e = {
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
  return e[l] || e.center;
}, D = (s) => {
  const l = s || "top-right", e = {
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
  return e[l] || e["top-right"];
}, F = (s) => {
  const l = s || "top-right";
  let e = "flex-col", t = "items-end", i = "space-y-2";
  return l.includes("bottom") && (e = "flex-col-reverse", i = "space-y-reverse space-y-2"), l.includes("left") ? t = "items-start" : l.includes("center") ? l === "center" ? t = "items-center justify-center" : t = "items-center" : t = "items-end", { flexDir: e, alignItems: t, spaceY: i };
};
function Y(s, l, e) {
  const t = s.getBoundingClientRect(), i = l.getBoundingClientRect(), n = i.width, r = i.height, a = t.width, c = t.height, o = { top: t.top, left: t.left }, d = window.innerWidth, u = window.innerHeight, p = 0, w = 0, f = u, m = d, g = (y, C) => {
    if (C) {
      if (y + n > m)
        return m - n - 10;
      if (y < w)
        return w + 10;
    } else {
      if (y + r > f)
        return f - r - 10;
      if (y < p)
        return p + 10;
    }
    return y;
  };
  let h, b, v = e;
  switch (e) {
    case "left":
      b = o.left - n - 5, b < w && (b = o.left + a + 5, v = "right"), h = g(o.top + c / 2 - r / 2, !1);
      break;
    case "top":
      h = o.top - r - 8, h < p && (h = o.top + c + 8, v = "bottom"), b = o.left + (a - n) / 2, b = g(b, !0);
      break;
    case "right":
      b = o.left + a + 5, b + n > m && (b = o.left - n - 5, v = "left"), h = g(o.top + c / 2 - r / 2, !1);
      break;
    case "bottom":
    default:
      h = o.top + c + 8, h + r > f && (h = o.top - r - 8, v = "top"), b = o.left + (a - n) / 2, b = g(b, !0);
      break;
  }
  return { top: h, left: b, position: v };
}
function T(s) {
  return String(s ?? "").replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;").replace(/'/g, "&#39;");
}
const H = /* @__PURE__ */ new WeakMap();
function L(s, l) {
  if (!s) return !1;
  const e = String(l);
  let t = H.get(s);
  return t || (t = /* @__PURE__ */ new Set(), H.set(s, t)), t.has(e) ? !1 : (t.add(e), !0);
}
function Z(s, l = {}) {
  console.log("[Sanitize] Input:", s);
  const {
    allowedTags: e = [
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
  } = l;
  console.log("[Sanitize] Allowed Tags:", e);
  const i = document.createElement("template");
  i.innerHTML = String(s ?? "");
  const n = document.createTreeWalker(i.content, NodeFilter.SHOW_ELEMENT), r = [];
  for (; n.nextNode(); ) {
    const a = n.currentNode, c = a.tagName.toLowerCase();
    if (!["script", "style", "iframe", "object", "embed"].includes(c)) {
      if (!e.includes(c)) {
        r.push(a);
        continue;
      }
      for (const o of [...a.attributes]) {
        const d = o.name.toLowerCase();
        (t[c] || []).concat(t["*"] || []).includes(d) || a.removeAttribute(o.name);
      }
      if (c === "a") {
        const o = (a.getAttribute("href") || "").trim(), d = o.startsWith("/") || o.startsWith("https://") || o.startsWith("http://") || o.startsWith("#");
        console.log("[Sanitize] a tag:", o, "Safe:", d), d || a.removeAttribute("href"), a.getAttribute("rel") || a.setAttribute("rel", "noopener noreferrer");
      }
    }
  }
  for (const a of r)
    a.replaceWith(document.createTextNode(a.textContent || ""));
  return i.innerHTML;
}
function k(s, l, { trusted: e = !1, alpine: t } = {}) {
  s && (s.innerHTML = e ? String(l ?? "") : Z(l), t?.initTree?.(s));
}
const V = (s) => {
  if (!s) {
    x.error("Tooltip: Alpine is undefined. Plugin initialization failed.");
    return;
  }
  s.store("tooltip", {
    visible: !1,
    content: "",
    placement: "top",
    top: -9999,
    // Start off-screen
    left: -9999,
    triggerEl: null,
    onShow: null,
    async show(e, t, i = {}) {
      this.triggerEl = e, this.content = t, this.placement = i.placement || "top", this.onShow = i.onShow, this.visible = !0, this.ensureContainer();
      const n = document.getElementById("uno-tooltip-body");
      n && k(n, this.content, { trusted: i.trustedHTML === !0, alpine: s }), await s.nextTick(), this.updatePosition(), this.onShow && this.onShow(document.getElementById("uno-tooltip-body"));
    },
    hide() {
      this.visible = !1, this.triggerEl = null;
    },
    updatePosition() {
      if (!this.triggerEl) return;
      const e = document.getElementById("uno-tooltip-content");
      if (!e) return;
      const t = Y(this.triggerEl, e, this.placement);
      this.top = t.top, this.left = t.left;
    },
    ensureContainer() {
      document.getElementById("uno-tooltip-container") || l();
    }
  });
  const l = () => {
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
               class="fixed px-3 py-2 text-sm font-medium rounded-lg shadow-sm whitespace-nowrap pointer-events-none border bg-white text-gray-900 border-gray-200 dark:bg-slate-800 dark:text-white dark:border-slate-700 z-[20000]"
              :style="\`top: \${$store.tooltip.top}px; left: \${$store.tooltip.left}px;\`"
              style="display: none;">
              <div id="uno-tooltip-body" class="pointer-events-auto"></div>
          </div>
        `, document.body ? document.body.appendChild(e) : document.documentElement.appendChild(e);
      } catch (e) {
        x.error("Tooltip: Failed to inject container", e);
      }
  };
  document.body ? l() : document.addEventListener("DOMContentLoaded", l), s.directive(
    "tooltip",
    (e, { expression: t, modifiers: i }, { evaluateLater: n, effect: r, cleanup: a }) => {
      const c = n(t);
      let o = {
        content: "",
        placement: i[0] || "top",
        trigger: i.includes("click") ? "click" : "hover"
      };
      i.includes("top") && (o.placement = "top"), i.includes("bottom") && (o.placement = "bottom"), i.includes("left") && (o.placement = "left"), i.includes("right") && (o.placement = "right");
      const d = (f) => {
        if (typeof f == "string" ? o.content = f : typeof f == "object" && (o = { ...o, ...f }), o.content && o.content.startsWith("#") && !o.content.includes(" ")) {
          const m = document.querySelector(o.content);
          m && (o.content = m.innerHTML);
        } else if (o.target) {
          const m = document.querySelector(o.target);
          m && (o.content = m.innerHTML);
        }
      };
      r(() => {
        c((f) => {
          d(f);
        });
      });
      const u = () => s.store("tooltip").show(e, o.content, o), p = () => s.store("tooltip").hide(), w = () => {
        const f = s.store("tooltip");
        f.visible && f.triggerEl === e ? f.hide() : u();
      };
      if (o.trigger === "click") {
        e.addEventListener("click", w);
        const f = (m) => {
          const g = s.store("tooltip"), h = document.getElementById("uno-tooltip-content");
          g.visible && g.triggerEl === e && !e.contains(m.target) && (!h || !h.contains(m.target)) && g.hide();
        };
        window.addEventListener("click", f), a(() => window.removeEventListener("click", f));
      } else
        e.addEventListener("mouseenter", u), e.addEventListener("mouseleave", p), e.addEventListener("focus", u), e.addEventListener("blur", p);
      a(() => {
        o.trigger === "click" ? e.removeEventListener("click", w) : (e.removeEventListener("mouseenter", u), e.removeEventListener("mouseleave", p), e.removeEventListener("focus", u), e.removeEventListener("blur", p));
      });
    }
  );
};
const X = (s) => {
  s.store("notify", {
    items: [],
    nextId: 1,
    activeRequests: 0,
    activeProgress: 0,
    // Notification Methods
    show(t, i = {}) {
      const {
        type: n = "info",
        timeout: r = 3e3,
        position: a = "top-right",
        closeIndicator: c = !1,
        showCloseButton: o = !0
      } = i, d = this.nextId++, u = {
        id: d,
        text: t,
        type: n,
        position: a,
        closeIndicator: c,
        showCloseButton: o,
        visible: !0
        // For transition
      };
      this.items.unshift(u), r > 0 && setTimeout(() => {
        this.close(d);
      }, r), this.ensureContainer();
    },
    close(t) {
      this.items.find((n) => n.id === t) && (this.items = this.items.filter((n) => n.id !== t));
    },
    success(t, i = {}) {
      const n = typeof i == "number" ? { timeout: i } : i;
      this.show(t, { type: "success", ...n });
    },
    warning(t, i = {}) {
      const n = typeof i == "number" ? { timeout: i } : i;
      this.show(t, { type: "warning", ...n });
    },
    error(t, i = {}) {
      const n = typeof i == "number" ? { timeout: i } : i;
      this.show(t, { type: "error", ...n });
    },
    info(t, i = {}) {
      const n = typeof i == "number" ? { timeout: i } : i;
      this.show(t, { type: "info", ...n });
    },
    // Loading Indicator Methods
    loading(t) {
      t ? this.activeRequests++ : this.activeRequests = Math.max(0, this.activeRequests - 1), this.ensureContainer();
    },
    // Progress Bar Methods
    progress(t) {
      t ? this.activeProgress++ : this.activeProgress = Math.max(0, this.activeProgress - 1), this.ensureContainer();
    },
    // Container Injection
    ensureContainer() {
      document.getElementById("uno-notify-container") || l();
    }
  });
  const l = () => {
    try {
      const t = document.createElement("div");
      t.id = "uno-notify-container", t.setAttribute("x-data", "{}");
      const r = `
        <style>
          @keyframes uno-progress-indeterminate {
            0% { left: -35%; width: 35%; }
            100% { left: 100%; width: 35%; }
          }
          .uno-progress-anim {
            animation: uno-progress-indeterminate 1.5s infinite linear;
          }
        </style>

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
        const c = D(a), { flexDir: o, alignItems: d, spaceY: u } = F(a);
        let p = `i.position === '${a}'`;
        return a === "top-center" && (p += " || i.position === 'top'"), a === "bottom-center" && (p += " || i.position === 'bottom'"), `
        <!-- ${a} -->
        <div class="fixed ${c} z-[30000] p-4 flex ${o} ${u} ${d} pointer-events-none">
            <template x-for="item in $store.notify.items.filter(i => ${p})" :key="item.id">
                ${e()}
            </template>
        </div>`;
      }).join(`
`)}
      `;
      t.innerHTML = r, document.body.appendChild(t);
    } catch (t) {
      x.error("Notify: Failed to inject container", t);
    }
  }, e = () => `
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
  document.body ? l() : document.addEventListener("DOMContentLoaded", l);
};
function G(s) {
  const l = String(s ?? "").trim();
  if (!l) return [];
  try {
    const e = JSON.parse(l);
    if (typeof e == "string") return [[e, {}]];
    if (e && typeof e == "object") return Object.entries(e);
  } catch {
  }
  return l.split(",").map((e) => e.trim()).filter(Boolean).map((e) => [e, {}]);
}
function J(s, l, e) {
  const t = s || document.body || document.documentElement;
  if (t) {
    if (typeof window < "u" && window.htmx?.trigger) {
      window.htmx.trigger(t, l, e);
      return;
    }
    t.dispatchEvent(new CustomEvent(l, { bubbles: !0, detail: e }));
  }
}
function K(s, l) {
  const e = s.headers.get("HX-Trigger");
  if (e)
    for (const [t, i] of G(e))
      J(l, t, i);
}
async function Q(s, l = {}) {
  const {
    method: e = "GET",
    htmxTarget: t,
    swap: i = "innerHTML",
    htmxConfig: n = {}
  } = l;
  if (typeof window < "u" && window.htmx?.ajax) {
    const r = t || document.createElement("div");
    t || (r.style.display = "none"), t || (document.body ? document.body.appendChild(r) : document.documentElement.appendChild(r));
    try {
      return await window.htmx.ajax(e, s, { target: r, swap: i, ...n }), r.innerHTML;
    } finally {
      t || r.remove();
    }
  }
  return O(s, l);
}
async function A(s, l = {}) {
  const {
    method: e = "GET",
    headers: t,
    body: i,
    credentials: n = "same-origin",
    triggerTarget: r = typeof document < "u" ? document.body : null
  } = l, a = await fetch(s, { method: e, headers: t, body: i, credentials: n });
  if (K(a, r), !a.ok) {
    const o = new Error(`Request failed: ${a.status}`);
    throw o.status = a.status, o;
  }
  return (a.headers.get("content-type") || "").includes("application/json") ? a.json() : a.text();
}
async function O(s, l = {}) {
  const e = await A(s, l);
  return typeof e == "string" ? e : JSON.stringify(e, null, 2);
}
const tt = (s) => {
  L(s, "uno:directive:safe-html") && s.directive("safe-html", (e, { expression: t }, { evaluateLater: i, effect: n }) => {
    const r = i(t);
    let a = null, c = null;
    n(() => {
      r((o) => {
        let d = "", u = !1;
        typeof o == "string" ? d = o : typeof o == "object" && o !== null && (d = o.html || "", u = o.trusted === !0), !(d === a && u === c) && (a = d, c = u, k(e, d, { trusted: u, alpine: s }));
      });
    });
  }), s.store("box", {
    items: [],
    // Public API
    alert(e, t = {}) {
      const i = typeof t == "string" ? { title: t } : t;
      return this.open("dialog", {
        content: e,
        showOkButton: !0,
        ...i
      });
    },
    confirm(e, t = {}) {
      return new Promise((i) => {
        this.open("confirm", {
          content: e,
          onOk: () => i(!0),
          onCancel: () => i(!1),
          onClose: () => i(!1),
          // Closing without choice = cancel
          ...t
        });
      });
    },
    prompt(e, t = "", i = {}) {
      return new Promise((n) => {
        const r = "prompt-" + Math.random().toString(36).substr(2, 9), a = `
          <div class="space-y-3">
            <p>${T(e)}</p>
            <input id="${r}" type="text" class="w-full px-3 py-2 border rounded focus:ring-2 focus:ring-blue-500 outline-none" value="${T(
          t
        )}">
          </div>
        `;
        this.open("confirm", {
          title: i.title || "Prompt",
          content: a,
          trustedHTML: !0,
          onOk: () => {
            const c = document.getElementById(r);
            n(c ? c.value : null);
          },
          onCancel: () => n(null),
          ...i
        });
      });
    },
    async open(e, t = {}) {
      let i = "dialog", n = t;
      typeof e == "object" ? (n = e, i = n.type || "dialog") : i = e;
      const r = n.boxId || Math.random().toString(36).slice(2, 10), a = {
        id: r,
        type: i,
        title: n.title || "",
        content: n.content || "",
        trustedHTML: n.trustedHTML === !0,
        width: n.width !== void 0 ? n.width : n.widthClass ? null : 400,
        height: n.height,
        widthClass: n.widthClass,
        heightClass: n.heightClass,
        containerClass: n.containerClass,
        modal: n.modal !== !1,
        // default true
        overlayOpacity: n.overlayOpacity ?? 0.5,
        position: n.position || "center",
        showCloseButton: !0,
        // default true?
        // Buttons
        okText: n.okButtonText,
        cancelText: n.cancelButtonText,
        onOk: n.onOk,
        onCancel: n.onCancel,
        onClose: n.onClose,
        onBeforeClose: n.onBeforeClose,
        onContentLoad: n.onContentLoad,
        // Callback after content is loaded (ajax/htmx)
        showOk: i === "confirm" || i === "message" || i === "dialog" && n.showOkButton,
        showCancel: i === "confirm" || i === "dialog" && n.showCancelButton,
        // State
        loading: !1,
        visible: !0
      }, c = s.store("i18n");
      if (a.okText || (a.okText = c ? c.t("box.ok") : "OK"), a.cancelText || (a.cancelText = c ? c.t("box.cancel") : "Cancel"), n.targetId) {
        const o = document.getElementById(n.targetId);
        o ? (a.content = o.innerHTML, a.trustedHTML = !0) : (x.warn("Box: Target element not found:", n.targetId), a.content = '<div class="p-4 text-red-500">Error: Content target not found</div>', a.trustedHTML = !0);
      } else if (n.hxUrl) {
        a.loading = !0, this.add(a);
        try {
          const o = await Q(n.hxUrl, { method: "GET" }), d = this.items.find((u) => u.id === r);
          d && (d.content = o, d.trustedHTML = !0, d.loading = !1, d.onContentLoad && s.nextTick(() => {
            d.onContentLoad(document.getElementById("uno-box-content" + r));
          }));
        } catch (o) {
          x.error("Box: HTMX Load Failed", o);
          const d = this.items.find((u) => u.id === r);
          d && (d.loading = !1, d.trustedHTML = !0, d.content = '<div class="p-4 text-red-500">Error loading content via HTMX</div>');
        }
        return r;
      } else if (n.ajaxUrl) {
        a.loading = !0, this.add(a);
        try {
          const o = await A(n.ajaxUrl), d = typeof o == "object" ? `<pre class="bg-gray-100 p-3 rounded overflow-x-auto whitespace-pre-wrap">${JSON.stringify(o, null, 2)}</pre>` : String(o), u = this.items.find((p) => p.id === r);
          u && (u.content = d, u.trustedHTML = !0, u.loading = !1, s.nextTick(() => {
            const p = document.getElementById("uno-box-content" + r);
            p && u.onContentLoad && u.onContentLoad(p);
          }));
        } catch (o) {
          x.error("Box: Ajax failed", o);
          const d = this.items.find((u) => u.id === r);
          d && (d.content = '<div class="text-red-500">Error loading content</div>', d.trustedHTML = !0, d.loading = !1);
        }
        return r;
      } else n.iframeUrl && (a.content = `<iframe frameborder="0" class="w-full h-full" src="${T(n.iframeUrl)}"></iframe>`, a.trustedHTML = !0);
      return this.add(a), r;
    },
    add(e) {
      this.items.push(e), this.ensureContainer();
    },
    close(e) {
      const t = this.items.find((i) => i.id === e);
      t && (t.onBeforeClose && t.onBeforeClose(), this.items = this.items.filter((i) => i.id !== e), t.onClose && t.onClose());
    },
    closeAll() {
      this.items = [];
    },
    getPositionClass(e) {
      const t = "fixed inset-0 flex transition-opacity duration-300 p-4 pointer-events-none", i = W(e);
      return `${t} ${i}`;
    },
    ensureContainer() {
      document.getElementById("uno-box-container") || l();
    }
  });
  const l = () => {
    try {
      const e = document.createElement("div");
      e.id = "uno-box-container", e.setAttribute("x-data", "{}");
      const t = `
        <template x-for="(item, index) in $store.box.items" :key="item.id">
            <div :id="'uno-box-outer-container' + item.id"
                 :class="$store.box.getPositionClass(item.position)"
                 :style="\`z-index: \${10000 + index}\`"
                 role="dialog"
                 aria-modal="true">
                 
                 <!-- Overlay: Covers viewport, blocks interaction, does NOT close modal on click -->
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
                      class="relative bg-white dark:bg-gray-900 rounded-lg shadow-2xl ring-1 ring-black/5 overflow-hidden flex flex-col max-h-[90vh] transition-all duration-300 transform pointer-events-auto"
                      :class="[item.widthClass, item.heightClass, item.containerClass]"
                      :style="\`\${item.width ? 'width: ' + item.width + 'px;' : ''} \${item.height ? 'height:' + item.height + 'px' : ''}\`"
                      x-show="item.visible"
                      x-transition:enter="ease-out duration-300"
                      x-transition:enter-start="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                      x-transition:enter-end="opacity-100 translate-y-0 sm:scale-100"
                      x-transition:leave="ease-in duration-200"
                      x-transition:leave-start="opacity-100 translate-y-0 sm:scale-100"
                      x-transition:leave-end="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95">
                      
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
                          <div :id="'uno-box-content' + item.id" x-show="!item.loading" x-safe-html="{ html: item.content, trusted: item.trustedHTML }"></div>
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
      `;
      e.innerHTML = t, document.body.appendChild(e);
    } catch (e) {
      x.error("Box: Failed to inject container", e);
    }
  };
  document.body ? l() : L(document, "uno:box:domcontentloaded") && document.addEventListener("DOMContentLoaded", l), L(document, "uno:box:delegator") && document.addEventListener("click", (e) => {
    const t = e.target.closest("[data-box-type]");
    if (!t) return;
    const i = t.getAttribute("data-box-type");
    if (!["dialog", "confirm", "message"].includes(i)) return;
    e.preventDefault();
    const n = t.dataset, r = {
      title: n.boxTitle,
      content: n.boxContent,
      width: n.boxWidth ? parseInt(n.boxWidth) : void 0,
      height: n.boxHeight ? parseInt(n.boxHeight) : void 0,
      ajaxUrl: n.boxAjaxUrl,
      hxUrl: n.boxHxUrl,
      iframeUrl: n.boxIframeUrl,
      targetId: n.boxTargetId,
      position: n.boxPosition,
      okButtonText: n.boxOkText,
      cancelButtonText: n.boxCancelText,
      overlayOpacity: n.boxOverlayOpacity ? parseFloat(n.boxOverlayOpacity) : void 0,
      showCloseButton: n.boxShowClose !== "false",
      boxId: n.boxId,
      cache: n.boxCache === "true",
      trustedHTML: n.boxTrusted === "true"
    };
    n.boxModal === "false" && (r.modal = !1), s.store("box").open(i, r);
  }), typeof window < "u" && L(document, "uno:box:htmx-afterSwap") && document.addEventListener("htmx:afterSwap", (e) => {
    e.detail.elt && s && s.initTree(e.detail.elt);
  });
}, S = {
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
function et(s) {
  const l = s.match(/^\/(.*?)\/([gimsuy]*)$/);
  return l ? { pattern: l[1], flags: l[2] } : { pattern: s, flags: "" };
}
function nt(s, l = {}) {
  const { defaultFlags: e = "i", maxPatternLength: t = 200 } = l, { pattern: i, flags: n } = et(String(s ?? ""));
  if (!i || i.length > t) return null;
  const r = String(n || e || "");
  if (!/^[gimsuy]*$/.test(r)) return null;
  try {
    return new RegExp(i, r);
  } catch {
    return null;
  }
}
function it(s) {
  try {
    return s ? JSON.parse(s) : [];
  } catch {
    return [];
  }
}
const ot = (s) => {
  const l = {
    invalid: "<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 12 12' fill='none' stroke='currentColor' class='w-4 h-4'><circle cx='6' cy='6' r='4.5'/><path stroke-linecap='round' d='M5.8 3.6h.4L6 6.5zM6 8.2a.6.6 0 110-1.2.6.6 0 010 1.2z'/></svg>",
    valid: "<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 8 8' class='w-4 h-4'><path fill='currentColor' d='M2.3 6.73L.6 4.53c-.4-1.04.46-1.4 1.1-.8l1.1 1.4 3.4-3.8c.6-.63 1.6-.27 1.2.7l-4 4.6c-.43.5-.8.4-1.1.1z'/></svg>"
  }, e = (t, i) => {
    const n = s.store("i18n");
    if (!n) return t;
    const r = {
      required: "Required",
      pattern: "Invalid format",
      min: `At least ${i?.n} characters`,
      max: `Max ${i?.n} characters`,
      int: "Integer required",
      email: "Invalid email address",
      url: "Invalid URL"
    };
    return n.t(`valid.${t}`, i) || r[t] || t;
  };
  s.store("valid", {
    callbacks: {},
    registerCallback(t, i) {
      typeof i == "function" && (this.callbacks[t] = i);
    },
    async validateForm(t) {
      if (!t) return !1;
      if (t._x_dataStack) {
        const i = t._x_dataStack.find((n) => n.validateAll);
        if (i)
          return await i.validateAll();
      }
      return !0;
    }
  }), s.data("validForm", () => ({
    fields: {},
    // Reactive state: { [name]: { value, error, valid, ... } }
    init() {
      this.$el.noValidate = !0;
    },
    // Register field (called by x-validate directive)
    register(t, i, n = []) {
      if (this.fields[t]) return;
      const r = [...n];
      if ((i.hasAttribute("required") || i.dataset.validRequired === "true") && r.push({ type: "required" }), (i.type === "email" || i.dataset.validType === "email") && r.push("email"), (i.type === "url" || i.dataset.validType === "url") && r.push("url"), i.dataset.validType) {
        const a = i.dataset.validType;
        a !== "email" && a !== "url" && r.push(a);
      }
      i.dataset.validPattern && r.push({ type: "regexp", expression: i.dataset.validPattern }), this.fields[t] = {
        el: i,
        rules: r,
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
      const n = (i.value || "").trim();
      let r = null;
      for (const a of i.rules) {
        const c = a.type || a, o = typeof a == "object" ? a : {};
        if (!n) {
          if (c === "required") {
            r = o.msg || e("required");
            break;
          }
          if (c === "emptyquit") break;
          continue;
        }
        if (S[c])
          S[c].test(n) || (r = o.msg || e(c) || e("pattern"));
        else if (c === "regexp" && o.expression) {
          const d = String(o.expression);
          o.__reSrc !== d && (o.__re = nt(d, { defaultFlags: "i" }), o.__reSrc = d);
          const u = o.__re;
          u && u.test(n) || (r = o.msg || e("pattern"));
        } else if (c === "limit")
          o.min && n.length < o.min && (r = o.msg || e("min", { n: o.min })), o.max && n.length > o.max && (r = o.msg || e("max", { n: o.max }));
        else if (c === "range") {
          const d = parseFloat(n);
          o.min !== void 0 && d < o.min && (r = o.msg || `Min ${o.min}`), o.max !== void 0 && d > o.max && (r = o.msg || `Max ${o.max}`);
        } else if (c === "compareto") {
          const d = document.getElementById(o.destid);
          d && n !== d.value && (r = o.msg || "Not match");
        } else if (c === "tag") {
          const d = o.mark || ",", u = n.split(d).filter((p) => p.trim());
          o.max && u.length > o.max && (r = o.msg || `Max ${o.max} items`);
        }
        if (r) break;
        if (c === "callback" && o.fn) {
          const d = s.store("valid").callbacks[o.fn];
          if (d)
            try {
              const u = await d(i.el.id, n);
              u === !1 || typeof u == "object" && u.valid === !1 ? r = (typeof u == "object" ? u.message : "") || "Invalid" : typeof u == "string" && (r = u);
            } catch {
            }
        }
        if (r) break;
      }
      return i.error = r, i.valid = !r, i.validating = !1, i.valid;
    },
    async validateAll() {
      const t = Object.keys(this.fields);
      return (await Promise.all(t.map((n) => this.validate(n)))).every((n) => n);
    },
    async submit(t) {
      if (t && (t.preventDefault(), t.stopPropagation()), Object.keys(this.fields).forEach((n) => {
        this.fields[n].touched = !0;
      }), await this.validateAll())
        this.$el.dispatchEvent(new CustomEvent("valid-submit", { bubbles: !0 })), !this.$el.hasAttribute("hx-post") && !this.$el.hasAttribute("hx-get") && this.$el.submit();
      else {
        const n = Object.keys(this.fields).find((r) => !this.fields[r].valid);
        n && this.fields[n].el && this.fields[n].el.focus();
      }
    }
  })), s.directive("validate", (t, i, { effect: n, cleanup: r }) => {
    const a = t.getAttribute("name");
    if (!a) return;
    const c = () => {
      let f = t.parentNode;
      f.classList.contains("uno-input-wrapper") || (f = document.createElement("div"), f.className = "relative w-full uno-input-wrapper", t.parentNode.insertBefore(f, t), f.appendChild(t));
      let m = f.querySelector(".uno-valid-icon");
      return m || (m = document.createElement("div"), m.className = "uno-valid-icon absolute right-0 top-0 bottom-0 flex items-center justify-center w-10 pointer-events-none transition-opacity duration-200", m.style.opacity = "0", f.appendChild(m)), { wrapper: f, icon: m };
    }, { icon: o } = c();
    let d = [];
    const u = t.getAttribute("valid");
    u && (d = it(u)), n(() => {
      const f = s.$data(t);
      f && typeof f.register == "function" && f.register(a, t, d);
    });
    const p = () => {
      const f = s.$data(t);
      if (f && f.update && f.update(a, t.value), d.find((m) => (m.type === "limit" || m === "limit") && m.show === !0)) {
        const m = d.find((g) => g.type === "limit" || g === "limit");
        m && m.max && t.value.length > m.max && (t.value = t.value.substring(0, m.max), f && f.update && f.update(a, t.value));
      }
    }, w = () => {
      const f = s.$data(t);
      f && f.touch && f.touch(a);
    };
    t.addEventListener("input", p), t.addEventListener("blur", w), n(() => {
      const f = s.$data(t);
      if (f && f.fields && f.fields[a]) {
        const m = f.fields[a], g = !!m.error, h = m.valid && m.dirty && !m.error;
        t.classList.toggle("is-invalid", g), t.classList.toggle("is-valid", h), t.classList.toggle("validating", m.validating), o.innerHTML = "", o.style.opacity = "0", g ? (t.classList.add("!border-red-500", "pr-10"), t.classList.remove("!border-green-500"), o.innerHTML = l.invalid, o.className = "uno-valid-icon absolute right-0 top-0 bottom-0 flex items-center justify-center w-10 pointer-events-none transition-opacity duration-200 text-red-500 dark:text-red-400", o.style.opacity = "1") : h ? (t.classList.add("!border-green-500", "pr-10"), t.classList.remove("!border-red-500"), o.innerHTML = l.valid, o.className = "uno-valid-icon absolute right-0 top-0 bottom-0 flex items-center justify-center w-10 pointer-events-none transition-opacity duration-200 text-green-500 dark:text-green-400", o.style.opacity = "1") : t.classList.remove("!border-red-500", "!border-green-500", "pr-10");
        let v = t.parentNode.nextElementSibling, y = null, C = 0;
        for (; v && C < 3 && !(v.nodeType === 1 && (v.classList.contains("animate-spin") ? y = v : y = v.querySelector(".animate-spin"), y)); )
          v = v.nextElementSibling, C++;
        y && (y.style.display = m.validating ? "block" : "");
      }
    }), r(() => {
      t.removeEventListener("input", p), t.removeEventListener("blur", w);
    });
  }), s.directive("valid-error", (t, { expression: i }, { effect: n }) => {
    const r = i;
    t.style.display = "none", n(() => {
      const a = s.$data(t);
      if (a && a.fields && a.fields[r]) {
        const o = a.fields[r].error;
        o ? (t.textContent = o, t.style.display = "block", t.classList.add("text-red-600", "dark:text-red-400", "text-xs"), t.classList.remove("text-gray-400")) : t.style.display = "none";
      }
    });
  }), s.directive("valid-count", (t, { expression: i }, { effect: n }) => {
    const r = i;
    n(() => {
      const a = s.$data(t);
      if (a && a.fields && a.fields[r]) {
        const c = a.fields[r], o = c.rules.find((d) => d.type === "limit" || d === "limit");
        if (o && o.max) {
          const d = (c.value || "").length, u = o.max - d, p = s.store("i18n");
          t.textContent = p && p.t("charCount", { count: u }) || `Remaining: ${u}`, t.style.display = "block";
        } else
          t.style.display = "none";
      }
    });
  });
}, $ = (s = "", l = {}) => ({
  active: s,
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
      const e = () => {
        const t = window.location.hash.substring(1);
        if (t) {
          const i = this.$el.querySelector(`[data-tab-hash="${t}"]`);
          if (i) {
            const n = i.getAttribute("data-tab-target");
            n && (this.active = n);
          }
        }
      };
      this.$nextTick(() => e()), this._hashHandler = () => e(), window.addEventListener("hashchange", this._hashHandler);
    }
  },
  destroy() {
    this._hashHandler && window.removeEventListener("hashchange", this._hashHandler);
  },
  isActive(e) {
    return this.active === e;
  },
  async open(e, t = null, i = null, n = null) {
    if (this.active = e, this.options.hash && i && window.location.hash.substring(1) !== i && history.pushState(null, null, "#" + i), n) {
      if (typeof window < "u" && window.htmx) {
        const r = document.getElementById(e);
        if (r && !this.loaded[e]) {
          k(
            r,
            `<div class="animate-pulse p-4 space-y-2">
               <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-3/4"></div>
               <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/2"></div>
             </div>`,
            { trusted: !0, alpine: window.Alpine || window.top?.Alpine }
          );
          try {
            await window.htmx.ajax("GET", n, {
              target: `#${e}`,
              swap: "innerHTML"
            }), r.setAttribute("data-tab-loaded", "true"), this.loaded[e] = !0, r.dispatchEvent(new CustomEvent("tab-loaded", { bubbles: !0 }));
          } catch {
            k(r, '<div class="text-red-500 p-4">Error loading content via HTMX</div>', {
              trusted: !0
            });
          }
        }
      } else {
        console.warn("Uno Tab: HTMX requested but window.htmx is not defined.");
        const r = document.getElementById(e);
        r && k(
          r,
          `<div class="text-amber-600 bg-amber-50 p-4 rounded border border-amber-200">
            <strong>Warning:</strong> HTMX library not found. Please include htmx.org script.
          </div>`,
          { trusted: !0 }
        );
      }
      return;
    }
    if (t) {
      const r = document.getElementById(e);
      if (r && !this.loaded[e]) {
        k(
          r,
          `<div class="animate-pulse p-4 space-y-2">
             <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-3/4"></div>
             <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/2"></div>
           </div>`,
          { trusted: !0, alpine: window.Alpine || window.top?.Alpine }
        );
        try {
          const a = await O(t, { method: "GET" });
          k(r, a, { trusted: !0, alpine: window.Alpine || window.top?.Alpine }), r.setAttribute("data-tab-loaded", "true"), this.loaded[e] = !0, r.dispatchEvent(new CustomEvent("tab-loaded", { bubbles: !0 }));
        } catch (a) {
          k(
            r,
            `<div class="text-red-500 p-4">Error loading content: ${a?.message || "Unknown error"}</div>`,
            { trusted: !0 }
          );
        }
      }
    }
  },
  trigger(e, t = null, i = null, n = null) {
    return {
      "@click.prevent": () => this.open(e, t, i, n),
      ":class": () => this.isActive(e) ? "current" : "",
      role: "tab",
      "aria-selected": () => this.isActive(e).toString(),
      "data-tab-target": e,
      "data-tab-hash": i || null,
      "data-tab-url": t || null,
      "data-tab-hx-url": n || null
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
      "x-show": () => this.isActive(e),
      role: "tabpanel",
      id: e
    };
  }
});
typeof window < "u" && (window.unoTabs = $);
const st = (s) => {
  s.data("unoTabs", $);
};
typeof window < "u" && (window.Alpine = window.Alpine || j, window.htmx = window.htmx || z, window.unoTabs = window.unoTabs || $);
const E = (s) => {
  if (!s) {
    x.error("Uno: Alpine is undefined, cannot register plugins");
    return;
  }
  if (s.store && s.store("tooltip")) {
    x.info("Uno: Plugins already registered");
    return;
  }
  x.info("Uno: Registering plugins"), U(s), V(s), X(s), tt(s), ot(s), st(s);
}, rt = () => {
  const s = j || window.Alpine;
  if (s && (x.info("Uno: Alpine found, registering"), E(s)), typeof window < "u") {
    const l = window.deferLoadingAlpine;
    window.deferLoadingAlpine = (e) => {
      let t = !1;
      const i = () => {
        t || (t = !0, e());
      };
      try {
        E(window.Alpine);
      } catch {
      }
      typeof l == "function" && l !== window.deferLoadingAlpine ? l(i) : i();
    };
  }
  document.addEventListener("alpine:init", () => {
    E(window.Alpine);
  });
};
typeof window < "u" && window.__UNO_INSTALLED__ ? x.info("Uno: Already installed, skipping init") : (typeof window < "u" && (window.__UNO_INSTALLED__ = !0), rt());
