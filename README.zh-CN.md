# Uno UI

<p align="left">
  <a href="README.md"><img alt="English" src="https://img.shields.io/badge/Language-English-blue?style=flat-square"></a>
  <a href="README.zh-CN.md"><img alt="简体中文" src="https://img.shields.io/badge/Language-%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87-red?style=flat-square"></a>
</p>

[![CI](https://github.com/artsyne/uno-ui/actions/workflows/ci.yml/badge.svg)](https://github.com/artsyne/uno-ui/actions/workflows/ci.yml)
[![npm](https://img.shields.io/npm/v/@artsyne/uno-ui.svg?style=flat-square)](https://www.npmjs.com/package/@artsyne/uno-ui)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://opensource.org/licenses/MIT)
[![Version](https://img.shields.io/badge/version-0.10.0-blue.svg?style=flat-square)](CHANGELOG.md)

Uno 是一个基于 **Alpine.js** 与 **Tailwind CSS 4** 的轻量 UI 组件库，采用
**Store-First、Directive-Driven** 架构，可选集成 **HTMX 2**。

## 当前状态

- 版本：`0.10.0`
- npm：已发布为 `@artsyne/uno-ui@0.10.0`。
- 项目网站：https://artsyne.github.io/uno-ui/
- 开源发布说明：[RELEASE.md](RELEASE.md)。
- 仓库内的 `dist/uno.umd.js` 与 `dist/uno.es.js` 是当前可用的构建产物；CI 会校验它们与
  `src/` 一致。
- Uno 不发布独立 CSS 文件，组件样式由 Tailwind CSS 根据模板中的类名生成。

## 特性

- **5 个组件**：Box（对话框）、Notify（通知）、Tooltip（提示）、Tab（标签页）、Valid（表单校验）。
- **Store-First**：通过 `Alpine.store('uno')` 编程式控制，指令负责声明式绑定。
- **轻量运行时**：唯一运行时依赖是 Alpine.js；HTMX 是可选 peer dependency。
- **懒初始化**：对话框、通知、Tooltip 的 DOM 容器在首次使用时才注入。
- **11 个内置语言包**：`en`、`zh-CN`、`zh-TW`、`es`、`ja`、`de`、`fr`、`ko`、`pt`、`ru`、`ar`。
- **TypeScript 类型**：随包提供 `types/index.d.ts`。
- **HTMX 2**：仅使用 XHR 响应，支持内容加载、进度条、表单拦截与 422 错误映射。

## 环境要求

| 依赖 | 版本 | 说明 |
| :--- | :--- | :--- |
| Alpine.js | 3.x（测试版本 3.17.3） | 必需，唯一运行时依赖。 |
| Tailwind CSS | 4.x | 必需；Uno 不发布 CSS。 |
| HTMX | 2.x | 可选。 |
| Node.js | >= 20.19（CI 使用 `.nvmrc` 的 24） | 仅开发、构建、测试需要。 |

Uno 不维护正式的浏览器测试矩阵。实际可用范围受 Alpine.js 3 与 Tailwind CSS 4 的浏览器要求
约束；请以二者的文档和你自己的目标浏览器测试为准。

## 快速开始

### 运行仓库示例

```bash
git clone https://github.com/artsyne/uno-ui.git
cd uno-ui
npm ci
npm run build
npm run dev
```

`npm run dev` 会打开 `http://localhost:5176/examples/`（固定端口，`strictPort`）。
组件页面位于 `examples/components/`。

### 在页面中使用本地 UMD 构建

```html
<!doctype html>
<html lang="zh-CN">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Uno 示例</title>

    <!-- Tailwind 4 浏览器构建：适合原型，不适合生产 -->
    <script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>

    <!-- 先加载 Uno，再加载 Alpine -->
    <script src="/vendor/uno-ui/dist/uno.umd.js" defer></script>
    <script src="/js/uno-config.js" defer></script>
    <script src="https://cdn.jsdelivr.net/npm/alpinejs@3.17.3/dist/cdn.min.js" defer></script>
  </head>
  <body>
    <button x-data @click="$store.uno.box.alert('你好，Uno')">打开对话框</button>
  </body>
</html>
```

```javascript
// /js/uno-config.js
window.Uno.setConfig({
  debug: false,
  ui: {
    base: 'rounded-lg shadow-md border border-gray-200 dark:border-gray-700'
  }
})
```

**加载顺序很重要**：Uno 需要在 Alpine 启动前完成注册。上面的顺序（Uno 在前、Alpine 在后）
是 UMD 的推荐路径。ESM 及另一种注册路径见
[Getting Started](docs/getting-started.md#alpine-loading-paths)。

### Tailwind 4 扫描 Uno 的类名

生产构建需要把包内 `src` 加入 Tailwind 的 `@source`，否则组件类名不会被生成：

```css
@import 'tailwindcss';
@source '../node_modules/@artsyne/uno-ui/src';
```

完整说明（暗色模式、浏览器构建、自定义主题）见 [Theming](docs/theming.md)。

### npm 与 CDN

```bash
npm install @artsyne/uno-ui alpinejs
```

```html
<script src="https://cdn.jsdelivr.net/npm/@artsyne/uno-ui@0.10.0/dist/uno.umd.js" defer></script>
```

jsDelivr 与 unpkg 均可访问已发布的 tarball。仓库内的 `dist/` 仍可用于本地开发和未发布的改动。

## 组件

| 组件 | 说明 | 文档 |
| :--- | :--- | :--- |
| **Box** | 模态框、Alert、Confirm、Prompt；支持堆叠、定位、HTMX 内容与 fill 布局。 | [Box API](docs/box-api.md) |
| **Notify** | 9 个位置的 Toast 通知、全局进度条、HTMX 内容。 | [Notify API](docs/notify-api.md) |
| **Tooltip** | Hover / Click / Focus 提示，自动翻转、视口夹取、HTMX 内容。 | [Tooltip API](docs/tooltip-api.md) |
| **Tab** | 标签组、hash 深链、HTMX 面板、APG 键盘模式。 | [Tab API](docs/tab-api.md) |
| **Valid** | 数组规则、异步校验、字符计数器、错误摘要、422 映射。 | [Valid API](docs/valid-api.md) |

## 文档索引

### 入门与概念

- [Getting Started](docs/getting-started.md) — 安装、Uno 与 Alpine 的加载顺序、Tailwind、HTMX、CSP、i18n。
- [Core Concepts](docs/concepts.md) — Store、Directive、Template、Layer、Action、生命周期。
- [Configuration](docs/configuration.md) — 4 级配置合并、默认值、`ui` token、`{base}` 与 `+`。
- [Theming and Tailwind](docs/theming.md) — `@source`、暗色模式、主题覆盖。
- [i18n](docs/i18n.md) — 11 个语言包、检测顺序、`extend`、`loadLocale`、`x-t`。

### 集成

- [HTMX Integration](docs/htmx.md) — XHR-only、事件、`Uno.htmx.*`、`afterSwap`、422。
- [CSP](docs/csp.md) — nonce 的真实范围与 Alpine 的 `'unsafe-eval'` 约束。
- [Accessibility](docs/accessibility.md) — 逐组件语义/键盘、axe 范围、手工验收。
- [TypeScript Usage](docs/typescript-usage.md) — 类型、tsconfig、已知缺口。

### 排障与升级

- [Troubleshooting](docs/troubleshooting.md) — 按症状排查。
- [Error Handling](docs/error-handling.md) — 全局错误、`onError`、HTMX 失败。
- [Migration Guide](docs/migration.md) — 0.7 / 0.8 / 0.9 → 0.10。
- [Test Reports](docs/test-reports.md) — 测试命令、axe 范围、CI 门禁。

## 配置与主题示例

```javascript
window.Uno.setConfig({
  debug: false,
  plugins: {
    box: { closeOnEsc: true },
    notify: { duration: 4000 }
  },
  ui: {
    base: 'rounded-lg shadow-md border border-gray-200 dark:border-gray-700',
    box: {
      okBtn: '+min-w-24'
    },
    layers: { box: 1400, notify: 1600, tooltip: 1700 }
  }
})
```

`{base}` 会在解析样式时替换为 `ui.base`；`+` 前缀表示追加而不是覆盖。
完整 token 表见 [Configuration](docs/configuration.md#ui-token-reference)。

## 国际化

```javascript
document.addEventListener('alpine:init', () => {
  Alpine.store('i18n').setLang('zh-CN')
})
```

语言选择会写入 `localStorage` 的 `uno.lang`。11 个语言包随包加载，无需额外请求。详见
[i18n](docs/i18n.md)。

## 安全说明

- Box 的 `content` 与 Tooltip 的 `content` 会经过 `Uno.dom.sanitize()`，但该实现是正则清洗，
  不是完整的 HTML sanitizer，可能被绕过。
- Notify/Tooltip/Box 的 HTMX 响应按可信服务端 HTML 处理，Uno 不会再次清洗。
- CSP nonce 只作用于 Uno 创建的容器与 `injectStyles()` 创建的 `<style>`；它不会让 Alpine 的
  表达式求值自动满足严格 CSP。
- 漏洞报告请使用 [SECURITY.md](SECURITY.md) 中的私密渠道。

## 开发

```bash
npm ci
npm run lint
npm run test:run
npm run build
```

- `npm run dev`：Vite 开发服务器。
- `npm run build`：构建到 `dist/`。
- `npm run test:run`：CI 使用的单次测试。
- `npm run test`：Vitest watch 模式。

`dist/` 是仓库跟踪的发布产物：修改 `src/` 后必须重新构建并提交。CI 会执行
`git diff --exit-code dist`。

## 项目结构

```text
uno-ui/
├── dist/                 # 构建产物（ESM + UMD）
├── docs/                 # 文档
├── examples/             # 可运行的组件示例
├── src/
│   ├── components/       # box / notify / tooltip / tab / valid
│   ├── core/             # uno / i18n / layer / plugin-registry
│   ├── locales/          # 11 个语言包
│   ├── utils/            # dom / position / validation / a11y / log / object
│   └── index.js          # 入口
├── tests/                # 单元测试、axe 扫描、类型测试
├── types/                # TypeScript 声明
├── CONTRIBUTING.md
├── SECURITY.md
├── CHANGELOG.md
└── package.json
```

## 贡献与发布

- 贡献流程、代码风格、测试与 `dist` 政策：[CONTRIBUTING.md](CONTRIBUTING.md)
- 变更记录：[CHANGELOG.md](CHANGELOG.md)
- 安全政策：[SECURITY.md](SECURITY.md)

## 许可证

MIT License，见 [LICENSE](LICENSE)。
