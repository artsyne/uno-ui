# Uno UI 插件优化执行文档

## 审查范围与现状

**范围：** `/uno/src`（ESM 源码）、`/uno/dist`（构建产物）、`/uno/examples`（示例）、`/uno/e2e`（Playwright）。

**入口与初始化：**

- 入口文件：`/uno/src/index.js`
- 自动注册：在 `deferLoadingAlpine` 与 `alpine:init` 阶段注册 store / directive
- 已注册插件：`tooltip` / `notify` / `box` / `valid` / `tab`，以及 `i18n` store

**兼容性与依赖一致性：**

- `package.json` 依赖：`alpinejs@^3.15.2`、`htmx.org@^2.0.8`
- `examples` 中存在版本差异：部分页面引用 `htmx@2.0.4`、`alpine@3.14.8` 或 `3.x.x`（`/uno/examples/components/*.html`）
- 建议：示例页与构建依赖版本对齐，避免“示例可用但真实项目不可用”的误导

**关键风险点（来自现状代码扫描）：**

- `innerHTML` / `x-html` 注入点：`tooltip`、`tab`、`box`（详见 P0-1）
- `window` 全局污染：`window.Alpine` / `window.htmx` / `window.deferLoadingAlpine`（详见 P0-4）
- 事件监听清理不足：`tab` 的 `hashchange`（详见 P0-2）

## 审查报告（摘要）

**总体评价：** 组件功能覆盖较完整，E2E 覆盖度高（Playwright），但当前实现更偏“业务内嵌式插件”，要达到开源组件库的通用性标准，需要补齐安全收口、生命周期/兼容性、a11y、以及性能基线与回归机制。

**最优先的整改方向：**

- 安全：统一 HTML 注入与默认净化，明确 trusted 路径（P0-1）
- 兼容：避免覆盖宿主全局对象，重复初始化需有稳定判定（P0-4）
- 稳定：全局事件监听器必须可清理，避免热更新/重复挂载泄漏（P0-2）
- 性能：减少 `Alpine.initTree` 的高频调用，建立可比的 perf 基线（P1-4 + 性能基线与目标）
- 可访问性：modal focus trap/ESC、tab 的 aria 契约补齐（P1-9）

---

## 测试用例补充计划（按插件）

**原则：**

- 每个优化点必须配套新增测试：`e2e` 自动化 + `examples` 手动验证
- 新增能力优先用测试“锁死契约”（避免重构时回归）
- 每阶段完成后同步更新 `/uno/doc/*-api.md`：对齐 API、参数约束、兼容性说明、以及测试覆盖说明

### Box

**建议新增 e2e（`/uno/e2e/box.spec.ts`）：**

- XSS：`content`/`targetId`/`ajaxUrl`/`hxUrl` 返回内容的净化验证（script/on\* 属性不生效）
- 兼容：`width/height` 非 number 输入时的行为（拒绝/回退到 class），避免 `pxpx`
- 生命周期：关闭后 body scroll 还原（含连续打开多个 modal 的计数正确性）
- 竞态：关闭后异步内容回写被忽略（ajax/htmx in-flight）
- a11y：ESC 关闭（可配置）、Focus Trap、关闭后焦点回到触发源

**建议补充 example 手动用例（`/uno/examples/components/box.html`）：**

- 新增“安全测试区”：展示恶意 HTML 的渲染结果与预期（被净化）
- 新增“并发/取消测试区”：快速打开/关闭、快速切换不同内容源（ajax/htmx）
- 新增“a11y 测试区”：键盘 Tab 循环、ESC 行为、焦点恢复

### Tooltip

**建议新增 e2e（`/uno/e2e/tooltip.spec.ts`）：**

- XSS：`content` 与 `target` DOM 读取内容的净化验证
- 性能：相同 content 多次 show 不重复 `initTree`（如果引入缓存/幂等）
- 交互：滚动/resize 后位置更新策略（如果引入）
- 兼容：页面存在多个 tooltip 指令实例时，clickOutside 监听器不累积

**建议补充 example 手动用例（`/uno/examples/components/tooltip.html`）：**

- 新增“恶意 content 示例”（script、onerror、onclick）
- 新增“长列表/大 HTML”示例，用于观察性能与布局抖动

### Notify

**建议新增 e2e（`/uno/e2e/notify.spec.ts`）：**

- HX-Trigger：模拟后端触发，统一事件入口能驱动 notify（如果引入统一 request）
- 定位：更多 position 组合（center-left/right、top-center/bottom-center）
- 资源：大量通知（例如 100 条）下的自动清理与 DOM 数量稳定

**建议补充 example 手动用例（`/uno/examples/components/notify.html`）：**

- 新增“HX-Trigger 演示区”（展示与 htmx/后端配合模式）
- 新增“高频通知压力区”（用于观察堆叠、关闭、性能）

### Tab

**建议新增 e2e（`/uno/e2e/tab.spec.ts` + `/uno/e2e/tab-htmx.spec.ts`）：**

- 内存：hashchange 监听器不随重复挂载增长（需要暴露计数或用 page.evaluate 检测）
- XSS：ajax/hxUrl 内容净化（若 Tab 支持渲染 HTML）
- 竞态：快速切换 tab 时取消/忽略旧请求回写
- a11y：`aria-controls`/`aria-labelledby` 与 manual 模式焦点行为

**建议补充 example 手动用例（`/uno/examples/components/tab.html`）：**

- 新增“manual 模式演示区”
- 新增“快速切换 + 异步加载”演示区（便于观察竞态）

### Valid

**建议新增 e2e（`/uno/e2e/valid.spec.ts`）：**

- ReDoS：超长/异常 regexp 输入的兜底验证（不应卡死，且给出错误）
- Callback：callback 抛错/超时的用户提示与状态恢复
- HTMX：valid-submit 与 hx-post 链路完整（成功/失败）

**建议补充 example 手动用例（`/uno/examples/components/valid.html`）：**

- 新增“regexp 输入示例”（正常/异常/超长）
- 新增“异步校验失败/超时”示例（可视化 validating/错误）

### i18n（核心）

**建议新增 e2e（`/uno/e2e/i18n.spec.ts`）：**

- 未提供语言包时的 fallback 行为（en/zh-CN 回退）
- 动态 extend 后的 UI 文案更新（如将新 key 注入）

**建议补充 example 手动用例（`/uno/examples/i18n.html`）：**

- 新增“动态 extend 示例”（运行时注入新字典）

## 🔴 P0 - 安全漏洞（必须立即修复）

### P0-1: XSS 防护 - HTML 注入统一收口（难度：中）

**问题：** 多个插件直接使用 `innerHTML`/`x-html`，缺少统一防护，且部分内容来自远端请求或 DOM 读取，存在 XSS 风险。

**现状定位（必须优先整改）：**

- `tooltip`：`body.innerHTML = this.content`（`/uno/src/components/tooltip/tooltip.js:42`）
- `tooltip`：从 `target.innerHTML` 读取并直接注入（`/uno/src/components/tooltip/tooltip.js:150-156`）
- `tab`：加载占位与远端内容使用 `innerHTML`（`/uno/src/components/tab/tab.js:77-92`、`113-131`）
- `box`：模板使用 `x-html="item.content"`（`/uno/src/components/box/box.js:288`）
- `box`：`target.innerHTML` 读取（`/uno/src/components/box/box.js:88`）

**方案：** 创建统一的 HTML 渲染入口 + 默认净化（不引入新依赖的前提下先落地最小安全集）。

```javascript
// src/utils/dom.js
export function sanitizeHTML(html, options = {}) {
  const {
    allowedTags = [
      'b',
      'i',
      'strong',
      'em',
      'u',
      'br',
      'p',
      'ul',
      'ol',
      'li',
      'code',
      'pre',
      'a',
      'span',
      'div'
    ],
    allowedAttrs = { a: ['href', 'target', 'rel'], '*': ['class'] }
  } = options

  const template = document.createElement('template')
  template.innerHTML = String(html ?? '')

  const walker = document.createTreeWalker(template.content, NodeFilter.SHOW_ELEMENT)
  const toRemove = []

  while (walker.nextNode()) {
    const el = walker.currentNode
    const tag = el.tagName.toLowerCase()

    if (!allowedTags.includes(tag)) {
      toRemove.push(el)
      continue
    }

    for (const attr of [...el.attributes]) {
      const name = attr.name.toLowerCase()
      const allow = (allowedAttrs[tag] || []).concat(allowedAttrs['*'] || [])
      if (!allow.includes(name)) el.removeAttribute(attr.name)
    }

    if (tag === 'a') {
      const href = (el.getAttribute('href') || '').trim()
      const isSafe =
        href.startsWith('/') ||
        href.startsWith('https://') ||
        href.startsWith('http://') ||
        href.startsWith('#')
      if (!isSafe) el.removeAttribute('href')
      if (!el.getAttribute('rel')) el.setAttribute('rel', 'noopener noreferrer')
    }
  }

  toRemove.forEach((n) => n.replaceWith(document.createTextNode(n.textContent || '')))
  return template.innerHTML
}

export function setHTML(el, html, { trusted = false, alpine } = {}) {
  if (!el) return
  el.innerHTML = trusted ? String(html ?? '') : sanitizeHTML(html)
  alpine?.initTree?.(el)
}
```

**收益：**

- 安全问题有唯一收口
- 未来替换净化库只改一处
- 自动触发 Alpine 初始化

**落地要点：**

- `box` 的 `x-html` 建议改为：仅允许渲染已净化 HTML（`item.contentSafe`），或默认 `x-text` 并提供 `trustedHTML` 显式开关
- `tab`/`tooltip` 的 DOM 读取（`target.innerHTML`）必须默认净化，且应提供 `trusted: true` 的白名单路径（仅用于内部静态模板）

---

### P0-2: 内存泄漏 - 统一事件清理机制（难度：低）

**问题：** Tab 等插件的全局事件监听器无清理。

**方案：** 引入幂等标记 + 清理机制

```javascript
// 幂等检查
Uno.utils.once = (el, key) => {
  const flag = `__uno_${key}`;
  if (el[flag]) return false;
  el[flag] = true;
  return true;
};

// Tab 组件
init() {
  if (!Uno.utils.once(this.$el, 'tabs_init')) return;

  this._cleanup = [];

  if (this.options.hash) {
    const handler = () => { /* ... */ };
    window.addEventListener('hashchange', handler);
    this._cleanup.push(() => window.removeEventListener('hashchange', handler));
  }
},

destroy() {
  this._cleanup?.forEach(fn => fn());
  this._cleanup = [];
}
```

**关键点：**

- 所有全局事件都记录清理函数
- Alpine directive 的 cleanup 钩子调用 destroy
- 使用 once 防止重复初始化

**现状定位（必须优先整改）：**

- `tab`：`window.addEventListener('hashchange', () => { ... })` 未移除（`/uno/src/components/tab/tab.js:35-47`）

---

### P0-3: HTMX 统一集成（难度：中）

**问题：** Box 用 fetch，无法利用 HTMX 的 HX-Trigger 机制。

**方案：** 创建统一请求接口

```javascript
// utils/request.js
Uno.request = async (url, options = {}) => {
  const { preferHTMX = true, target, onTrigger } = options

  // 优先 HTMX
  if (preferHTMX && window.htmx) {
    return new Promise((resolve, reject) => {
      const temp = target || document.createElement('div')

      // 监听 HX-Trigger
      if (onTrigger) {
        document.addEventListener('htmx:trigger', onTrigger, { once: true })
      }

      const cleanup = () => {
        temp.removeEventListener('htmx:afterSwap', handleSuccess)
        if (!target) temp.remove()
      }

      const handleSuccess = (e) => {
        resolve(temp.innerHTML)
        cleanup()
      }

      temp.addEventListener('htmx:afterSwap', handleSuccess)
      htmx.ajax('GET', url, { target: temp })
    })
  }

  // 回退 fetch
  return await fetch(url).then((r) => r.text())
}
```

**使用：**

```javascript
// Box 插件中
await Uno.request(ajaxUrl, {
  preferHTMX: true,
  onTrigger: (e) => {
    // 自动处理后端的 HX-Trigger
    if (e.detail.showNotification) {
      Alpine.store('notify').success(e.detail.showNotification)
    }
  }
})
```

---

### P0-4: 全局对象污染与重复初始化收口（难度：中）

**问题：** 入口文件直接写入 `window.Alpine` / `window.htmx` / `window.deferLoadingAlpine`，可能覆盖宿主已有实例或破坏其他库的初始化链路。

**现状定位：**

- `window.Alpine = Alpine`、`window.htmx = htmx`（`/uno/src/index.js:18-19`）
- 覆盖 `window.deferLoadingAlpine`（`/uno/src/index.js:47-57`）
- 以 `alpine.store('tooltip')` 判断“是否已安装”，存在与宿主 store 名冲突的误判风险（`/uno/src/index.js:28-33`）

**方案：**

- 仅在 `window.Alpine`/`window.htmx` 不存在时才赋值，避免覆盖
- 安装标记使用专用 flag（如 `window.__UNO_INSTALLED__ = true`）或 `alpine.store('__uno')` 私有 store
- 对 `deferLoadingAlpine` 做链式合并（保留原先的 hook）

**收益：**

- 降低与宿主应用/第三方脚本的兼容性风险
- 避免重复注册导致的事件监听器叠加与内存增长

---

### P0-5: 正则表达式输入的 DoS 风险（ReDoS）（难度：低）

**问题：** `valid` 支持通过 `data-valid-pattern` 或规则 `{ type: 'regexp', expression: '...' }` 动态构造 `RegExp`，如果表达式来自不可信输入或可被用户配置，可能触发灾难性回溯导致 UI 卡死。

**现状定位：**

- `new RegExp(pattern, flags || 'i')`（`/uno/src/components/valid/valid.js`，regexp 分支）

**方案：**

- 限制：pattern 长度上限、flags 白名单（`gimsuy`）、禁止嵌套量词等高风险结构（最小可行是长度 + try/catch）
- 兜底：`new RegExp` 失败时给出明确错误并视为校验失败，而不是静默吞掉

**收益：**

- 避免被错误/恶意配置拖垮前端主线程

## 🟡 P1 - 架构优化（提升可维护性）

### P1-1: 插件命名空间 - 物理解耦（难度：中）

**问题：** 插件逻辑独立但耦合在同一文件。

**方案：** 形式化插件结构（即便在同一文件）

```javascript
// 插件注册协议
Uno.plugins = {}

Uno.use = (name, plugin) => {
  Uno.plugins[name] = plugin

  // 注册 Alpine store
  if (plugin.store) {
    Alpine.store(name, plugin.store)
  }

  // 注册指令
  if (plugin.directive) {
    Alpine.directive(name, plugin.directive)
  }

  // 执行初始化
  plugin.init?.()
}

// 使用
Uno.use('box', {
  store: { items: [], open() {} },
  init() {
    /* 注入容器 */
  }
})
```

**收益：**

- 插件边界清晰
- 易于独立测试
- 未来可拆分文件

---

### P1-2: 初始化顺序保证（难度：中）

**问题：** 插件、Alpine、HTMX 初始化顺序隐式。

**方案：** 显式初始化队列

```javascript
Uno._initQueue = []
Uno._initialized = false

Uno.use = (name, plugin) => {
  Uno.plugins[name] = plugin
  Uno._initQueue.push({ name, plugin })
}

Uno.init = () => {
  if (Uno._initialized) return

  // 1. 先注册所有 store
  Uno._initQueue.forEach(({ name, plugin }) => {
    if (plugin.store) Alpine.store(name, plugin.store)
  })

  // 2. 再执行插件 init
  Uno._initQueue.forEach(({ plugin }) => {
    plugin.init?.()
  })

  // 3. 最后绑定 HTMX
  document.addEventListener('htmx:afterSwap', (e) => {
    Alpine.initTree(e.detail.elt)
  })

  Uno._initialized = true
}

// Alpine 启动后执行
document.addEventListener('alpine:init', () => Uno.init())
```

---

### P1-3: DOM 注入统一管理（难度：低）

**问题：** 每个插件各自注入容器。

**方案：** 统一容器管理

```javascript
Uno.dom.ensureContainer = (id, html) => {
  let container = document.getElementById(id);

  if (!container) {
    container = document.createElement('div');
    container.id = id;
    container.setAttribute('x-data', '{}');
    container.innerHTML = html;
    document.body.appendChild(container);
  }

  return container;
};

// 插件使用
ensureContainer() {
  Uno.dom.ensureContainer('uno-box-container', this.template);
}
```

---

### P1-4: 性能优化 - DOM 查询缓存（难度：低）

**原则：** 频繁访问的元素缓存，销毁时清理。

```javascript
// Tooltip 示例
updatePosition() {
  if (!this._contentEl) {
    this._contentEl = document.getElementById("uno-tooltip-content");
  }
  // 使用 this._contentEl
},

hide() {
  this._contentEl = null; // 清理缓存
}
```

---

### P1-5: API 契约一致性与兼容性（难度：中）

**问题：** 当前 API 对“入参类型/默认值/向后兼容”约束不够明确，存在隐性 bug 与宿主集成风险。

**现状定位：**

- `box`：`width`/`height` 在渲染时按 number 拼接 `px`，但示例中存在传 string 的用法，可能生成 `pxpx`（`/uno/src/components/box/box.js:53-57`、`250-251` 与 `/uno/examples/components/box.html:250-258`）
- `box`：读取 `dataset` 的 `cache` 字段但未实现缓存行为（`/uno/src/components/box/box.js:357-359`）
- `tab`：`open(tab, ajaxUrl, hash, hxUrl)` 参数顺序耦合且易误用；建议参数对象化（`/uno/src/components/tab/tab.js:54-133`）

**方案：**

- 对外 API 采用“参数对象”而非位置参数（尤其 `tab.open`），并明确必填/可选项
- 对 `box.width/height` 做输入规范：只接受 `number`（像素），字符串改走 `widthClass/heightClass`
- 将向后兼容层（`data-*` 触发、旧参数）集中到一个适配器，核心 store API 保持干净

**收益：**

- 降低集成踩坑成本（更可预测的 API）
- 降低线上不可预期 UI bug 的概率

---

### P1-6: 错误处理与日志策略统一（难度：低）

**问题：**

- `log` 永久开启 `DEBUG = true`，生产环境噪音与性能开销不可控（`/uno/src/utils/log.js:9-35`）
- `ajax` 工具缺少 `res.ok` 检查与超时/取消机制（`/uno/src/utils/ajax.js:13-18`）
- 组件内存在 `console.warn` 与 `log` 混用（例如 i18n、tab），不利于统一收口

**方案：**

- `DEBUG` 以 `globalThis.UnoConfig?.debug` 或 `import.meta.env.MODE` 驱动，并提供 `log.level`
- `ajax` 补齐：`res.ok`、超时、错误标准化（返回 `{ ok, status, data }` 或抛出统一错误对象）
- 约定：仅 `log` 输出；`console.*` 只在 `debug` 模式存在

**收益：**

- 生产环境更干净，排障更可控
- 错误与告警更结构化（便于接入监控/埋点）

---

### P1-7: CSP/样式注入兼容（难度：中）

**问题：** 部分组件通过 `innerHTML` 注入 `<style>` 或大量内联样式，严格 CSP（禁用 `unsafe-inline`）下可能失效。

**现状定位：**

- `notify` 注入 `<style>`（`/uno/src/components/notify/notify.js:162-191`）
- `box`/`tooltip` 依赖大量内联 `style` 与模板字符串拼接（`/uno/src/components/box/box.js:250-251`、`/uno/src/components/tooltip/tooltip.js:100-101`）

**方案：**

- 将关键动画/样式迁移为 class（由 Tailwind utilities 或预置 CSS 文件承载）
- `z-index`、尺寸等必须动态的样式保持最小化，其他改为 class

**收益：**

- 能在更严格的安全策略（CSP）下工作
- 降低运行期字符串拼接与 DOM 注入开销

---

### P1-8: 请求取消/并发控制与缓存（难度：中）

**问题：**

- `tab` 的 `fetch` 与 `box` 的 `ajaxUrl` 加载无取消能力，用户快速切换/关闭时会产生无效工作与竞态写入
- 远端加载失败时的反馈不一致，且缺少重试/退避策略

**方案：**

- 对 `fetch` 路径引入 `AbortController`：同一容器重复请求时取消上一请求
- 对 `htmx.ajax` 路径在 close/destroy 时触发取消（使用 htmx 的请求管理能力或外层标记忽略回写）
- 实现 `cache`：以 `url -> html/text` 缓存并提供 `cacheTTL`/`forceReload` 选项

**收益：**

- 降低 CPU/网络开销，提升 UI 响应性
- 降低竞态导致的“内容闪回”类问题

---

### P1-9: 可访问性与交互一致性（难度：中）

**问题：** 作为通用 UI 插件，`box`/`tab`/`tooltip` 的无障碍与键盘交互需要更接近 WAI-ARIA 模式，否则在企业/开源场景会被视为质量缺陷。

**现状定位：**

- `tab` 已包含 `role="tablist"`/`role="tab"`/`role="tabpanel"` 与部分方向键（`/uno/src/components/tab/tab.js:148-162`），但缺少 `aria-controls`、panel 的 `aria-labelledby`，以及 `manual` 模式的完整焦点管理
- `box` 作为 modal：缺少 ESC 关闭、Focus Trap、关闭后焦点回到触发源（store 里没有 `triggerEl`/`returnFocus` 逻辑）

**方案：**

- Tab：补齐 `aria-controls`（指向 panel id）、panel 追加 `aria-labelledby`（指向 tab id），并在 `manual: true` 时仅移动 focus 不切换 active
- Box：实现 Focus Trap（Tab/Shift+Tab 循环）、ESC 关闭（可配置）、打开时聚焦第一个可聚焦元素，关闭后恢复到 trigger

**预期效果：**

- 键盘可完整操作、屏幕阅读器信息更准确
- 与主站/其他框架的 a11y 策略更兼容

---

### P1-10: 按需注入与页面负载控制（难度：低）

**问题：** `box`/`notify`/`tooltip` 默认 eager 注入全局容器，即使页面只使用其中一个组件，也会产生额外 DOM、样式与 Alpine 观察开销。

**现状定位：**

- `box` eager 注入（`/uno/src/components/box/box.js:324-329`）
- `notify` eager 注入（`/uno/src/components/notify/notify.js:232-237`）
- `tooltip` eager 注入（`/uno/src/components/tooltip/tooltip.js:117-122`）

**方案：**

- 默认 lazy：首次调用 `show/open` 或首次命中指令时注入容器
- 提供 `UnoConfig.eager` 允许需要全局占位的场景主动开启

**收益：**

- 首屏 DOM 更小，降低 Alpine init 与样式计算开销

## 🟢 P2 - 增强功能（可选）

### P2-1: 插件配置化（难度：低）

把各插件样式，配置等提取到配置，提供更丰富的定制选项，可在外部配置文件中覆盖默认值。下面是建议的例子：

```javascript
window.UnoConfig = {
  security: {
    enableSanitization: true,
    allowedTags: ['b', 'i', 'strong', 'a']
  },
  classes: {
    notify: {
      success: 'bg-green-100 border-green-500',
      error: 'bg-red-100 border-red-500'
    }
  },
  validation: {
    timeout: 300, // 防抖延迟
    showIcon: true
  }
}
```

### P2-2: 调试模式（难度：低）

```javascript
Uno.debug = process.env.NODE_ENV === 'development'

Uno.log = {
  info(...args) {
    if (Uno.debug) console.log('[Uno]', ...args)
  },
  perf(label, fn) {
    if (!Uno.debug) return fn()
    const start = performance.now()
    const result = fn()
    console.log(`[Uno] ${label}: ${performance.now() - start}ms`)
    return result
  }
}
```

---

## 实施路线图

### 阶段交付物要求（强制）

每个阶段的每一项优化点完成后，都必须同步交付以下“新增部分”，否则视为未完成：

- **e2e 自动测试**：为新增/变更的行为增加或更新 Playwright 用例（放在对应 `e2e/*.spec.ts`）
- **examples 手动测试**：为新增/变更的行为补充可操作的示例入口（放在对应 `examples/*.html`）
- **API 文档更新**：同步更新 `/uno/doc/*-api.md`（接口签名、参数约束、兼容性、示例代码、测试覆盖说明）

### 第 1 阶段

- [x] P0-1 XSS 防护 - 创建 `src/utils/dom.js` 并替换所有 HTML 注入点
- [x] P0-2 内存泄漏 - 添加 once 标记和清理机制
- [x] P0-3 HTMX 统一 - 创建统一请求入口并支持 HX-Trigger
- [x] P0-4 全局污染收口 - 安装标记与 hook 链式合并
- [x] P0-5 regexp ReDoS - 限制表达式与兜底策略

**第 1 阶段新增部分（必须补充）：**

- [x] e2e：新增 XSS/净化/事件清理/HTMX 触发/regexp 兜底用例（按插件归档）
- [x] examples：新增“安全测试区”和“错误/异常兜底区”（便于人工复核）
- [x] doc：更新 `/uno/doc/box-api.md`、`tooltip-api.md`、`tab-api.md`、`valid-api.md`（增加安全与输入约束说明）

### 第 2 阶段

6. P1-1 插件命名空间 - 形式化 `Uno.use`
7. P1-2 初始化顺序 - 显式 `Uno.init()`
8. P1-3 DOM 统一管理
9. P1-4 性能优化
10. P1-5 API 契约一致性与兼容性
11. P1-6 错误处理与日志策略统一
12. P1-7 CSP/样式注入兼容
13. P1-8 请求取消/并发控制与缓存
14. P1-9 可访问性与交互一致性
15. P1-10 按需注入与页面负载控制

**第 2 阶段新增部分（必须补充）：**

- e2e：新增 a11y（focus trap/ESC/tab aria）、取消/并发/缓存、初始化幂等与全局污染兼容用例
- examples：新增“a11y 演示区”“取消/并发/缓存演示区”“懒注入演示区”
- doc：更新全部 `/uno/doc/*-api.md`（增加初始化/生命周期、错误处理、CSP、缓存与取消、a11y 说明）

### 第 3 阶段

16. P2 增强功能按需实施

**第 3 阶段新增部分（必须补充）：**

- e2e：针对新增配置项/调试能力/主题定制新增用例
- examples：新增“配置覆盖演示页”或在现有页追加配置区块
- doc：更新 `/uno/doc/*-api.md` 中的配置项与默认值表

---

## 关键原则

1. **安全优先**：所有 HTML 注入必须经过 `Uno.dom.setHTML`
2. **生命周期清晰**：init → use → destroy 三段式
3. **HTMX 优先**：能用 htmx.ajax 就不用 fetch
4. **幂等性保证**：所有初始化都通过 `once()` 检查
5. **统一出口**：DOM/请求/事件都有统一管理接口

---

## 性能基线与目标

**基线说明：** 数据来自本仓库 Playwright 本地运行 `npm test -- -g "perf:"`（3 浏览器并行）。不同机器会有波动，建议用同一 CI/硬件做对比。

| 场景                       | Chromium | Firefox | WebKit | 目标（优化后）        |
| -------------------------- | -------: | ------: | -----: | --------------------- |
| Box：打开/关闭 50 次       |    493ms |   600ms |  117ms | \< 300ms（Chromium）  |
| Tooltip：show/hide 100 次  |   3321ms |  3342ms | 6657ms | \< 1500ms（Chromium） |
| Tab：切换 200 次（同一组） |   18.7ms |    34ms |   16ms | \< 15ms（Chromium）   |

**重点解释：**

- Tooltip 基线偏高的主要原因通常是每次 `show` 都会 `innerHTML` + `Alpine.initTree`（`/uno/src/components/tooltip/tooltip.js:41-44`），建议引入“内容不变不重建”与 DOM 缓存策略（见 P1-4）。

---

## 验证清单

需要增加 e2e 自动化测试覆盖，确保所有功能在真实环境中正常运行。也同步丰富 example 中的手动测试用例，确保测试用例覆盖所有场景。

```javascript
// 安全测试
✓ XSS 注入测试（<script>alert(1)</script>）
✓ HTML 属性注入（onclick="alert(1)"）
✓ Template 内容净化
✓ regexp 表达式异常/超长兜底（避免页面卡死）

// 性能测试
✓ 创建/销毁 1000 个 tooltip 无内存泄漏
✓ 连续打开/关闭 100 个 box 性能稳定
✓ Tab 切换 1000 次 hashchange 监听器数量不变
✓ 基线采集：`npm test -- -g "perf:"` 输出对比表
✓ 对比口径：同一浏览器、同一机器、同一数据集

// 集成测试
✓ HTMX afterSwap 后 Alpine 正常工作
✓ 后端 HX-Trigger 能触发前端 notify
✓ 多个插件并存无冲突
```

---

## 总结

这份文档整合了：

- ✅ 原有的 3 个核心优化（XSS、内存泄漏、HTMX 统一）
- ✅ review 中的架构建议（插件命名空间、初始化顺序、DOM 统一）
- ✅ 实用的性能优化（缓存、幂等、配置化）
- ✅ 明确的实施路线和验证方法
