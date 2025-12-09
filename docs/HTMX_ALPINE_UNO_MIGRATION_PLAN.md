# Hebe 前端架构重构计划：Kiko 迁移至 Uno + HTMX + Alpine + Tailwind v4

**文档版本**: v1.5
**创建日期**: 2025-12-09
**状态**: 提案 (Proposed)

---

## 1. 现有架构技术评估与重构原则

### 1.1 核心原则 (Key Principles)

本重构计划严格遵循以下原则，以降低风险并确保系统稳定性：

1.  **不向后兼容 (No Backward Compatibility)**: 重构完成后，将彻底移除 `kiko.js`, `jquery.js` 及其相关初始化代码，不保留兼容层。系统将完全运行在现代技术栈上。
2.  **保留 Tailwind CSS v4 (Keep Tailwind v4)**: 继续使用现有的 Tailwind CSS v4 样式体系。不引入 UnoCSS，确保现有 UI 视觉效果 100% 不变。
3.  **服务端接口最小侵入 (Minimal Server Impact)**:
    - **禁止修改** Service 层和 Repository 层逻辑。
    - **禁止修改** 现有的 JSON API 路由和参数结构。
    - **允许微调** Controller 层以支持 **Content Negotiation**（内容协商），即根据请求头（如 `HX-Request`）决定返回 HTML 片段还是完整页面。
4.  **功能与样式零回归 (Zero Regression)**: 重构后的页面在视觉上和功能流程上必须与旧版完全一致。
5.  **CSS 代码保留 (Keep CSS As-Is)**: `/website/static/src/css` 目录下的代码不需要重构。当前基于 Tailwind CSS v4 的样式定义已经是最佳实践，应予以保留并直接复用。
6.  **JS 代码规范 (Structured JS)**: **禁止在 HTML 模板中编写内联 JavaScript**。所有 JS 逻辑必须位于 `/website/static/src/js` 目录中，通过模块化方式组织。

### 1.2 重构核心目标

当前重构的核心是通过引入 **Uno UI** (作为核心 UI 框架), **HTMX** (处理服务端交互), **Alpine.js** (处理前端交互) 和 **Tailwind CSS v4** (保持现有样式)，来替换老旧的 `Kiko UI`、`jQuery` 技术栈，实现前端架构的现代化，同时保持视觉和功能的完全一致。

### 1.3 现有架构痛点

- **高耦合**: 业务逻辑通过 `data-action` 属性与 `kiko.js` 中的 jQuery 选择器强绑定。
- **状态混乱**: 状态散落在 DOM 属性中，缺乏统一管理。
- **维护困难**: 每次修改 UI 结构都需要同步修改 JS 里的 DOM 操作代码。

---

## 2. Uno + HTMX + Alpine + Tailwind v4 技术方案

### 2.1 技术选型

- **Uno UI**: 核心 UI 组件库（Dialogs, Notifications, Tabs, Validation），替换 Kiko UI 组件。
  - **部署方式**: 使用已编译文件，路径为 `/website/static/dist/js/uno.umd.js`。
- **HTMX**: 负责与服务端交互（AJAX 替代），处理 HTML 片段的替换与更新。
- **Alpine.js**: 负责纯前端的交互逻辑（如 Dropdown、Modal、Tabs、Toast），替代 jQuery 插件。
- **Tailwind CSS v4**: 保持现状，利用 v4 的即时编译特性。

### 2.2 服务端适配策略

为了在不修改 Service 接口的前提下支持 HTMX，我们将采用 **ViewService 增强模式**：

1.  **请求识别**: 后端通过 `HX-Request` Header 识别 HTMX 请求。
2.  **片段渲染**:
    - **常规请求**: `Render(c, "diary/index", data)` -> 渲染完整 HTML（含 Header/Footer）。
    - **HTMX 请求**: `Render(c, "diary/index", data)` -> 自动识别，仅渲染 `content` 区块的 HTML 片段。
3.  **JSON 接口复用**: 对于必须调用现有 JSON API 的场景（极少数），使用 Alpine.js 发起 `fetch` 请求并处理响应，或使用 HTMX 的 `hx-ext="json-enc"` 扩展。但优先推荐后端 Controller 适配返回 HTML。

### 2.3 前端资源架构 (Frontend Architecture)

我们将重新设计 `/website/static/src/js` 目录结构，以模块化方式组织代码，替代原有的 `app.js` 和 `diary.js`。

**目录结构规划**:

```text
/website/static/src/js/
├── core/               # 核心框架初始化
│   └── init.js         # 初始化 Alpine, HTMX, Uno 配置
├── stores/             # Alpine 全局状态 (Global Stores)
│   ├── theme.js        # 主题管理 (Dark Mode)
│   ├── i18n.js         # 多语言状态管理
│   └── user.js         # 用户登录状态
├── modules/            # 业务模块逻辑 (Alpine Data Components)
│   ├── auth.js         # 登录/注册表单逻辑
│   ├── diary.js        # 日记相关交互 (点赞, 删除)
│   ├── layout.js       # 布局交互 (导航栏, 页脚)
│   └── settings.js     # 用户设置逻辑
└── utils/              # 工具函数
    └── http.js         # 简单的 Fetch 封装 (可选)
```

**资源引用路径**:

- **Uno UI**: `/public/dist/js/uno.umd.js` (对应物理路径 `/website/static/dist/js/uno.umd.js`)
- **App CSS**: `/public/src/css/tailwind.css` (对应物理路径 `/website/static/src/css/tailwind.css`)

### 2.4 前端构建方案 (Frontend Build Plan)

本重构引入了新的构建流程，需确保构建脚本能够自动化处理以下任务：

1.  **CSS 构建 (CSS Build)**:
    - **策略**: **保持现状 (As-Is)**。
    - **说明**: 继续沿用项目现有的 CSS 构建流程（基于 Tailwind CSS v4）。`/website/static/src/css` 目录下的所有文件保持不变，构建输出路径仍为 `/website/static/dist/css`。

2.  **JS 构建 (JS Build)**:
    - **策略**: **模块化打包 (Modular Bundling)**。
    - **说明**: 需要更新构建脚本，支持将 `/website/static/src/js` 下的模块化文件（core, stores, modules, utils）打包处理。
      - **开发环境**: 可直接通过 `<script type="module">` 引用各独立文件以便调试，或使用简单的合并策略。
      - **生产环境**: 需将所有模块打包、压缩为 `app.min.js` (或其他命名)，并输出到 `/website/static/dist/js` 目录。
    - **规范**: 构建过程需支持 ES Module 语法，并确保不破坏代码的执行顺序（如 `core/init.js` 必须最先执行）。

3.  **Uno UI 资源同步 (Uno UI Sync)**:
    - **策略**: **自动复制 (Auto-Copy)**。
    - **逻辑**: 每次执行前端构建时，必须包含一个自动复制步骤。
    - **源文件**: `/Users/yasong.cy/go/src/idiary.ai/hebe/uno/dist/uno.umd.js`
    - **目标路径**: `/Users/yasong.cy/go/src/idiary.ai/hebe/website/static/dist/js/uno.umd.js`
    - **要求**: 不需要重新编译 Uno 项目，仅执行文件复制操作，确保 Website 始终使用最新的 Uno UI 产物。

---

## 3. 分阶段实施路线图 (Phased Roadmap)

采用按**文件目录结构**划分的渐进式重构策略。

### **Phase 1: 基础设施与静态资源 (预计 3 天)**

- **目标**: 建立新老共存的环境，搭建新的 JS 架构。
- **涉及文件**: `views/header.html`, `static/src/js/core/init.js`, `static/src/js/stores/*.js`
- **任务**:
  1.  **JS 架构搭建**:
      - 创建 `static/src/js/core/init.js`: 引入 Alpine.js 和 HTMX 的初始化配置。
      - 创建 `static/src/js/stores/theme.js`: 迁移原 `theme.js` 逻辑到 Alpine Store。
      - 创建 `static/src/js/stores/i18n.js`: 迁移原 `i18n.js` 逻辑到 Alpine Store。
  2.  **资源引用更新**:
      - 修改 `views/header.html`，移除旧的 `kiko` 初始化脚本。
      - 引入 `/public/dist/js/uno.umd.js`。
      - 引入新的模块化 JS 文件 (开发阶段可分别引入，生产阶段合并)。
  3.  **后端扩展**: 创建 `internal/app/view/htmx.go` 扩展，实现 `RenderPartial` 逻辑。

- **验证与测试 (Verification & Testing)**:
  - [ ] **浏览器控制台检查**: 确保 Alpine.js, HTMX, Uno UI 初始化成功，无 JS 错误。
  - [ ] **资源加载验证**: 检查 Network 面板，确认 `uno.umd.js` 和新模块化 JS 文件均返回 200 OK。
  - [ ] **Store 状态验证**: 在控制台输入 `Alpine.store('theme')`，确认能获取到主题状态对象。

### **Phase 2: 公共布局与导航 (预计 5 天)**

- **目标**: 迁移全局导航、页脚和用户菜单，验证 Alpine.js 交互能力。
- **涉及文件**: `views/navbar.html`, `views/footer.html`, `static/src/js/modules/layout.js`
- **任务**:
  1.  **创建 Layout 模块**: 新建 `static/src/js/modules/layout.js`，定义 `navbar` 和 `footer` 的 Alpine 组件。
  2.  **用户下拉菜单**: 在 `navbar.html` 中使用 `x-data="layout.navbar()"`，移除 `data-tooltip-target`。
  3.  **语言切换**: 使用 Uno 的 i18n 支持，在 `layout.js` 中调用 `Alpine.store('i18n')`。
  4.  **暗色模式**: 绑定 `x-data` 到 `Alpine.store('theme')`。

- **验证与测试 (Verification & Testing)**:
  - [ ] **视觉回归**: 对比重构前后 Navbar 和 Footer 的截图，确保布局、间距、字体完全一致。
  - [ ] **交互功能**:
    - 点击用户头像，下拉菜单应平滑展开/收起。
    - 点击语言切换，Cookie 中的 `idiary-lang` 应更新，且页面自动刷新。
    - 点击暗色模式切换，HTML 标签应正确添加/移除 `dark` 类，页面配色即时变化。

### **Phase 3: 核心业务 - 浏览与交互 (预计 7 天)**

- **目标**: 迁移高频只读页面和轻量级交互。
- **涉及文件**: `views/diary/index.html`, `static/src/js/modules/diary.js`
- **任务**:
  1.  **创建 Diary 模块**: 新建 `static/src/js/modules/diary.js`，处理点赞、删除等逻辑。
  2.  **日记列表翻页**: 使用 HTMX (`hx-get`, `hx-target`) 实现无刷新翻页。
  3.  **点赞功能**: 将 `data-action="like"` 替换为 `hx-post`，或在 `diary.js` 中封装 Alpine Action。
  4.  **删除功能**: 使用 Uno 的 `Box` 组件。在 `diary.js` 中定义 `confirmDelete` 方法，触发 Uno Box 弹窗。

- **验证与测试 (Verification & Testing)**:
  - [ ] **列表渲染**: 首页日记列表加载正常，样式无崩坏。
  - [ ] **HTMX 交互**:
    - 点击“下一页”，列表应局部刷新，页面不闪烁，Network 面板可见 `HX-Request: true`。
    - 点击“点赞”，计数器应即时 +1。
  - [ ] **Uno 组件**:
    - 点击“删除”，应弹出 Uno 风格的确认框。
    - 确认删除后，该条目应从列表中平滑消失。

### **Phase 4: 核心业务 - 表单与提交 (预计 10 天)**

- **目标**: 处理最复杂的表单交互，验证服务端校验反馈机制。
- **涉及文件**: `views/form/*.html`, `views/passport/*.html`, `static/src/js/modules/auth.js`
- **任务**:
  1.  **创建 Auth 模块**: 新建 `static/src/js/modules/auth.js`。
  2.  **登录/注册**: 改造为 `hx-post`，服务端返回带错误信息的 HTML 片段更新表单区域。使用 Uno `Valid` 组件进行表单验证。
  3.  **日记发布**: 确保富文本编辑器与 Alpine 组件的生命周期 (`init`, `destroy`) 集成。

- **验证与测试 (Verification & Testing)**:
  - [ ] **表单验证**: 提交空表单，应显示 Uno Valid 组件定义的错误提示样式。
  - [ ] **提交反馈**:
    - 登录成功，应跳转至首页。
    - 登录失败，应在表单区域内显示错误信息，不刷新页面。
  - [ ] **富文本编辑器**: 确保编辑器能正常加载、输入内容，并随表单一起提交。

### **Phase 5: 清理与验收 (预计 5 天)**

- **目标**: 移除旧代码，确保系统整洁。
- **涉及文件**: `views/base.html`, `static/src/js/`
- **任务**:
  1.  **移除 Kiko/jQuery 依赖**:
      - 检查并删除 `header.html` 中对 `kiko.js`, `jquery.js`, `idiary.js` 的引用。
      - 删除 `static/src/js` 目录下的旧文件：`app.js`, `diary.js`, `kiko*.js`, `jquery*.js`。
  2.  **迁移逻辑清理**:
      - 扫描所有 HTML 模板，确保无残留的 `data-action` 属性（除非已被新逻辑复用）。
      - 确认所有 `onclick` 内联事件已被 `x-on:click` 或 `hx-trigger` 替代。
  3.  **全站回归测试**:
      - 执行完整的视觉回归测试套件。
      - 模拟用户完整操作路径（注册 -> 登录 -> 写日记 -> 点赞 -> 删除 -> 登出）。

- **验证与测试 (Verification & Testing)**:
  - [ ] **控制台零错误**: 全站浏览，控制台不应出现 404 (资源丢失) 或 JS Error。
  - [ ] **网络请求检查**: 确认不再加载任何 jQuery 或 Kiko 相关资源。
  - [ ] **最终验收**: 确认所有功能点均已迁移，且用户体验无退化。

---

## 4. 测试策略总纲 (Testing Strategy Overview)

> **注意**: 具体的测试用例已分散至上述各个阶段的“验证与测试”环节中。本节仅概述使用的工具和通用方法。

### 4.1 工具链

- **Playwright/Puppeteer**: 用于自动化视觉回归测试（截图对比）和端到端（E2E）功能测试。
- **Chrome DevTools**: 用于开发过程中的手动验证（Console, Network, Elements）。

### 4.2 核心测试类型

- **视觉回归 (Visual Regression)**: 像素级对比重构前后的页面截图，确保样式零差异。
- **功能验证 (Functional Verification)**: 验证交互逻辑（点击、提交、弹窗）是否符合预期。
- **接口验证 (API Verification)**: 确保 HTMX 请求正确发送 `HX-Request` 头，且服务端正确返回 HTML 片段。

---

## 5. 风险登记与应对

| 风险点           | 可能性 | 影响 | 应对措施                                                                                                                                        |
| :--------------- | :----- | :--- | :---------------------------------------------------------------------------------------------------------------------------------------------- |
| **JS 依赖丢失**  | 高     | 中   | 某些功能（如 Markdown 渲染、时间格式化）可能依赖了 jQuery 插件。需提前梳理 `static/src/js/utils.js`，将必要逻辑移植为原生 JS 或 Alpine Helper。 |
| **样式覆盖**     | 中     | 低   | 移除 `idiary.min.css` 可能导致部分自定义样式丢失。需对比 CSS 文件，确保所有自定义样式都在 `tailwind.css` 或内联类中。                           |
| **复杂交互卡顿** | 低     | 中   | HTMX 频繁请求可能导致服务器压力略增。对高频操作（如搜索）使用 `hx-trigger="keyup changed delay:500ms"` 防抖。                                   |

---
