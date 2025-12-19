
# 🚀 AI 执行指令：Uno UI 语义化架构与高性能标准重构 (v0.4.0)

**Role:** 资深前端架构师 (Alpine.js / Tailwind / 高性能 JS 专家)
**Task:** 深度重构 `uno.es.js` 源码，将 Uno UI 从“个人插件”提升至“工业级开源标准”。

---

### 1. 核心层：构建样式总线 (Theme Engine)

* **重构 `X` 类 (Core Engine, 对应源码中的核心类)**：
* 在 `this.options` 中初始化 `ui` 对象，建立三级级联体系：`UnoConfig.ui[组件][元素]` > `UnoConfig.ui.base[元素]` > `defaultCls`。
* **实现 `getStyle(comp, element, defaultCls)` 方法**：
* **替换逻辑**：若配置存在且不以 `+` 开头，完全替换默认类名。
* **追加逻辑**：若配置以 `+` 开头，则将配置值（去重后）叠加到默认类名上。
* **示例**：`getStyle('box', 'okBtn', 'bg-blue-500')` 在配置为 `{okBtn: '+rounded-full'}` 时应返回 `bg-blue-500 rounded-full`。





---

### 2. 性能与安全层 (Optimization & Security)

* **请求缓存 (LRU Strategy)**：
* 在核心请求方法中引入容量为 50 的 LRU 缓存，缓存 GET 请求产生的 HTML/JSON 结果。


* **正则解析缓存 (Performance)**：
* 在 `valid` 插件中使用全局 `WeakMap` 缓存已编译的正则对象，键为元素的校验规则字符串。


* **安全加固 (Security Default)**：
* 修改 `Tab` 和 `Box`：默认 `trusted: false`。
* 在此模式下，AJAX 加载的内容必须强制流经 `Uno.dom.sanitizeHTML`，且**物理移除**所有 `<script>` 标签以杜绝 XSS。
* 仅当显式传参 `trusted: true` 时允许执行脚本。



---

### 3. 组件层：语义化槽位化 (Style Slots)

重构以下插件，移除硬编码字符串，统一调用 `getStyle`：

* **Box**: 抽离 `overlay`, `container`, `okBtn`, `cancelBtn`。
* **Notify**: 抽离 `container`, `item`；状态色（success 等）映射至 `ui.notify.types`。
* **Tooltip**: 抽离 `content`, `arrow`；圆角必须联动 `ui.base.rounded`。
* **Tab**: **彻底废弃** `.current` 类名，改为从 `ui.tab.active` 读取（默认值 "current"）。
* **Valid**: 抽离 `inputInvalid` (默认 `is-invalid`) 和 `errorMsg` (默认 `text-red-500`)。

---

### 4. API 层：补全全量公共接口 (Public API) 🌟

确保所有组件都能通过 `window.Uno` 对象进行 JS 命令式控制：

* **Uno.notify**: 补全 `success/error/info/warn(msg, opts)` 快捷方法。
* **Uno.box**: 新增 `Uno.box.closeAll()`；标准化 `Uno.box.open(opts)` 接口。
* **Uno.tooltip (新增控制权)**:
* `Uno.tooltip.show(target, content, opts)`: 手动为指定 DOM 弹出提示。
* `Uno.tooltip.hide(target)`: 手动关闭。


* **Uno.valid (深度增强)**:
* `Uno.valid.validate(form)`: 触发全表单校验，返回 `Promise<Boolean>`。
* `Uno.valid.clear(form)`: **新增**，重置表单验证状态（移除所有 `is-invalid` 类名和错误消息）。
* `Uno.valid.set(el, {valid, message})`: 确保该手动设置方法同步触发样式总线的类名更新。



---

### 5. 交付要求与自动化测试

* **代码质量**：保持 `Uno.use()` 架构，确保混淆后体积增量控制在 10% 以内。
* **零破坏性变更**：默认配置下，视觉表现必须与 v0.3.0 完全一致。
* **Playwright E2E 测试用例**：
1. **回归测试**：设置 `ui.base.rounded = 'rounded-none'`，断言 Box 和 Tooltip 不含圆角类名。
2. **API 测试**：执行 `Uno.valid.set(input, {valid: false})`，断言该 input 包含 `is-invalid` 类名。
3. **安全测试**：在 `trusted: false` 模式下加载含 `<script>` 的 AJAX 内容，验证脚本未执行。



---

**AI 助手请开始执行：**

1. 请先输出重构后的核心 `X` 类（Core）及 `getStyle` 实现。
2. 随后按 **Box -> Notify -> Tooltip -> Tab -> Valid** 的顺序提供修改后的插件代码块。
3. 最后提供配套的 Playwright 测试脚本。