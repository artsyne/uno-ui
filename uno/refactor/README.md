# Uno UI 重构文档

本目录包含 Uno UI 框架的完整重构方案和实施计划。

## 📁 文档结构

### 🏗️ 基础设施

- **`infrastructure.md`** - 核心基础设施重构规范
  - Core 配置管理系统
  - 多语言系统重构
  - Z-Index 层级管理
  - 无障碍访问 (A11y) 系统
  - 错误边界系统
  - 工具函数标准化
  - 性能管理系统

### 🔧 组件重构方案

- **`box.md`** - Box 组件重构规范
- **`tab.md`** - Tab 组件重构规范
- **`tooltip.md`** - Tooltip 组件重构规范
- **`notify.md`** - Notify 组件重构规范
- **`valid.md`** - Valid 组件重构规范

### 📋 实施计划

- **`implementation-plan.md`** - 详细的重构实施计划
  - 五个阶段的详细任务清单
  - AI 执行指南
  - 测试策略和验收标准
  - 每个阶段的 Review 任务

### 📊 质量保证

- **`review-logs.md`** - 重构阶段 Review 日志
  - 每个阶段的详细 Review 结果

  - 问题发现和解决方案记录

  - 质量指标达成情况

  - 最终重构总结

### 🚫 不向后兼容原则

- **`no-backward-compatibility.md`** - 不向后兼容原则
  - 完全移除历史包袱的决策

  - 不兼容变更清单

  - 迁移策略和时间表

  - 重构原则和最佳实践

### 📝 JSDoc 注释标准

- **`jsdoc-standards.md`** - JSDoc 注释标准
  - 函数、类、配置对象的注释规范

  - 最佳实践和避免事项

  - 工具推荐和验证清单

  - 开源项目注释标准

## 🎯 重构目标

将 Uno UI 重构为 **面向 Alpine + HTMX + Tailwind CSS 最佳实践的现代化组件库**：

1. **彻底现代化** - 完全拥抱 Alpine + HTMX + Tailwind CSS 生态
2. **统一架构模式** - 所有组件遵循相同的 Store-First, Directive-Driven 架构
3. **性能优化** - 显著提升运行时性能和内存使用
4. **多语言支持** - 完整的国际化解决方案
5. **HTMX 原生集成** - 完全依赖 HTMX 生态，移除所有 AJAX 封装
6. **无障碍访问** - 完整的 A11y 支持，符合现代标准
7. **错误处理** - 健壮的错误边界和异常处理
8. **层级管理** - 统一的 Z-Index 管理，避免遮挡问题
9. **开发效率** - 提升代码可维护性和开发体验

## 🚫 不向后兼容原则

**Uno UI v2.0 将完全重构，不保持向后兼容：**

1. **API 变更**：所有 API 将重新设计，不兼容旧版本
2. **配置变更**：配置系统完全重构，旧配置无效
3. **方法重命名**：Store 方法按新规范命名，旧方法名移除
4. **HTML 属性变更**：只支持新的 `x-uno-*` 指令，移除 `data-uno-*` 属性
5. **依赖更新**：明确要求 Alpine + HTMX + Tailwind CSS
6. **移除历史包袱**：删除所有遗留代码和兼容性适配层

**迁移策略：** 提供详细的迁移指南，但不是自动兼容

## 🚀 使用方式

这些文档专为 AI 驱动的重构工作设计，包含详细的执行步骤和质量要求。

### 📋 完整重构执行

```bash
# AI 执行完整重构
"Step 1: 全局上下文加载
请仔细阅读以下核心规范文档：
- [实施计划] /refactor/implementation-plan.md
- [不兼容准则] /refactor/no-backward-compatibility.md (严禁保留历史包袱)
- [注释标准] /refactor/jsdoc-standards.md (所有代码必须符合 JSDoc 标准)

Step 2: 执行策略确认
确认你已理解以下核心原则：
- 严格遵循 5 个阶段的顺序执行，不可跳跃。
- 采用 'Store-First, Directive-Driven' 架构模式。
- 所有新代码必须面向 Alpine + HTMX + Tailwind CSS 最佳实践。
- 遇到不兼容变更时，参考 no-backward-compatibility.md，不提供适配层。

Step 3: 阶段性执行循环
对于每个阶段，请执行：
1. 阅读该阶段特定的组件/基础设施文档。
2. 按照 implementation-plan.md 拆解任务。
3. 执行 TDD 循环 (Red-Green-Refactor)。
4. 进行自我 Review 并记录到 /refactor/review-logs.md。

质量门控：
- 代码覆盖率 > 90%
- E2E 测试通过率 100%
- 无任何 Linter 警告
- 所有公共 API 具备标准 JSDoc 注释

错误处理：
- 遇到阻碍性问题，暂停并记录到日志，不要强行 Hack。
- 确保每个阶段提交前都可独立回滚。"
```

### 🏗️ 第一阶段：基础设施重构

```bash
# AI 执行第一阶段 (Infrastructure)
"Step 1: 核心上下文加载
请深入阅读以下文档：
- [核心规范] /refactor/infrastructure.md
- [不兼容准则] /refactor/no-backward-compatibility.md (关注配置系统变更)
- [注释标准] /refactor/jsdoc-standards.md

Step 2: 架构预演 (思维链)
在编写代码前，请先在思维链中推演：
1. **配置合并逻辑**：如何实现 4-Level Merge (PluginDefaults -> GlobalPlugin -> GlobalComp -> Instance)？
2. **i18n 加载**：如何设计异步加载机制以支持按需获取语言包？
3. **Z-Index**：如何构建统一的层级管理系统避免魔法数字？

Step 3: 增量式 TDD 执行
请按照 implementation-plan.md 的第一阶段任务，逐一执行 Red-Green-Refactor：
1. [Config] 实现 Core 配置系统，支持深层合并和默认值覆盖。
2. [i18n] 实现多语言 Store，支持 `setLocale` 和 `t` 函数。
3. [Layer] 实现 Z-Index Manager，定义基础层级常量。
4. [A11y] 实现无障碍焦点管理和键盘导航工具函数。
5. [Error] 实现全局错误边界 (Error Boundary) 和日志收集。

Step 4: 反模式最终检查
在提交代码前，自我检查：
- ❌ 是否保留了旧版的 `window.UnoConfig` (应完全移除)？
- ❌ 是否使用了复杂的原型继承 (应使用组合或简单的对象合并)？
- ❌ 是否遗漏了关键函数的 JSDoc 注释？
- ❌ 是否使用了 `var` 或非标准 ES6+ 语法？

验收标准：
- 运行 `npm run test:unit` 确保基础设施单元测试 100% 通过。
- 验证配置合并逻辑在所有 4 个层级都能正确覆盖。
- 验证 i18n 切换语言时无页面闪烁。
- 代码符合 jsdoc-standards.md 要求。"
```

### 🧪 第二阶段：试点组件重构 (Notify)

```bash
# AI 执行第二阶段 (Notify 组件)
"Step 1: 核心上下文加载
请深入阅读以下文档：
- [组件规范] /refactor/notify.md
- [不兼容准则] /refactor/no-backward-compatibility.md (注意 API 和 HTML 属性变更)
- [注释标准] /refactor/jsdoc-standards.md

Step 2: 架构预演 (思维链)
在编写代码前，请先在思维链中构建 NotifyStore 的接口定义：
1. **Store-First**：如何通过 `Alpine.store('uno').notify` 管理 items 数组？
2. **触发机制**：如何统一 Directive, Global API, Event 三种触发方式？
3. **动画管理**：如何利用 Alpine 的 `x-transition` 配合 CSS 实现平滑进出？

Step 3: 增量式 TDD 执行
任务拆解如下，请逐一执行 Red-Green-Refactor 循环：
1. [Core] 实现 NotifyStore：仅包含 items 数组和 add/remove 基础方法。
2. [UI] 实现 `x-uno-notify` 指令：负责渲染 UI 列表，绑定 A11y 属性 (role='alert')。
3. [API] 实现 `Alpine.store('uno').notify.success/error` 等代理方法。
4. [Feat] 集成多语言 (i18n) 和 Z-Index 管理。
5. [Anim] 实现进入/退出动画。

Step 4: 反模式最终检查
在提交代码前，自我检查：
- ❌ 在指令中直接操作 innerHTML (应通过 Alpine 响应式数据驱动)？
- ❌ 使用 setTimeout 硬编码动画时间 (应监听 transitionend)？
- ❌ 硬编码样式 (应使用 Uno.getStyle)？
- ❌ 是否保留了 `Uno.notify.show` 旧 API (应完全移除)？

验收标准：
- 运行 `npm run test:e2e:notify` 确保通过。
- 屏幕阅读器能正确朗读新通知。
- 快速连续触发 10 条通知不卡顿。
- 所有公共方法包含 JSDoc 注释。"
```

### 🎯 第三阶段：核心组件重构

```bash
# AI 执行第三阶段 (Box, Tooltip, Tab)
"Step 1: 核心上下文加载
请深入阅读以下文档：
- [Box规范] /refactor/box.md
- [Tooltip规范] /refactor/tooltip.md
- [Tab规范] /refactor/tab.md
- [不兼容准则] /refactor/no-backward-compatibility.md
- [注释标准] /refactor/jsdoc-standards.md

Step 2: 架构预演 (思维链)
请针对每个组件进行预演：
1. **Box**：如何实现 HTMX 原生集成 (hx-target) 和 9点定位系统？
2. **Tooltip**：如何设计全局容器以避免 `overflow: hidden` 截断，并使用 RAF 优化位置计算？
3. **Tab**：如何实现 URL Hash 双向同步和动态内容加载？

Step 3: 增量式 TDD 执行
请按组件顺序执行 TDD 循环：

**A. Box 组件**
1. [Core] 实现 BoxStore，管理实例和 Z-Index。
2. [UI] 实现 `x-uno-box` 指令，处理 `hx-get` 等属性。
3. [Feat] 实现模态遮罩 (Overlay) 和键盘 ESC 关闭。

**B. Tooltip 组件**
1. [Core] 实现 TooltipStore，管理激活状态。
2. [UI] 实现 `x-uno-tooltip` 指令，绑定 mouseenter/leave 事件。
3. [Perf] 实现 RAF 轮询定位计算 (仅在显示时)。

**C. Tab 组件**
1. [Core] 实现 TabStore，管理 activeTab 状态。
2. [UI] 实现 `x-uno-tab` 指令和 `x-uno-tab-panel`。
3. [Feat] 实现 Hash 监听和同步。

Step 4: 反模式最终检查
在提交代码前，自我检查：
- ❌ Box 是否使用了自定义 AJAX (必须使用 HTMX)？
- ❌ Tooltip 是否在不可见时仍在计算位置 (性能杀手)？
- ❌ Tab 是否在切换时强制刷新页面 (应为无刷新)？
- ❌ 是否兼容了旧版 `data-uno-*` 属性 (应仅支持 `x-uno-*`)？

验收标准：
- 三个组件的 E2E 测试全部通过。
- Box 支持 HTMX 加载内容。
- Tooltip 在滚动时位置跟随流畅。
- Tab 切换时 URL Hash 正确更新。
- 代码符合 JSDoc 标准。"
```

### 🔧 第四阶段：辅助组件重构

```bash
# AI 执行第四阶段 (Valid & Global)
"Step 1: 核心上下文加载
请深入阅读以下文档：
- [Valid规范] /refactor/valid.md
- [不兼容准则] /refactor/no-backward-compatibility.md
- [注释标准] /refactor/jsdoc-standards.md

Step 2: 架构预演 (思维链)
请预演 HTMX 验证集成：
1. 如何拦截 `htmx:validation:validate` 事件？
2. 如何处理 `hx-target-error` 扩展？
3. 如何实现实时验证与提交验证的协同？

Step 3: 增量式 TDD 执行
1. [Core] 实现 ValidStore，管理错误消息和验证规则。
2. [UI] 实现 `x-uno-valid` 指令，支持 `realtime` 配置。
3. [HTMX] 实现 HTMX 事件拦截器，阻止无效表单提交。
4. [Cleanup] 扫描全项目，移除所有已废弃的旧代码 (如 `data-uno-*` 选择器)。

Step 4: 反模式最终检查
- ❌ 是否使用了复杂的自定义验证逻辑 (应优先使用 HTML5 验证 API)？
- ❌ 是否在验证失败时依然触发了 HTMX 请求？
- ❌ 全局 CSS 中是否残留了旧版本样式？

验收标准：
- 表单验证能正确拦截 HTMX 请求。
- 服务端返回的错误能正确显示。
- 项目中不再包含任何 `legacy` 代码。"
```

### 📝 第五阶段：最终优化和总结

```bash
# AI 执行第五阶段 (Optimization)
"Step 1: 核心上下文加载
请阅读 /refactor/implementation-plan.md 和所有组件文档，建立全局视野。

Step 2: 性能与质量审计
执行以下审计命令：
1. `npm run test:e2e` (全量回归)
2. `npm run analyze` (Bundle 分析)
3. Chrome DevTools Memory Profile (手动)

Step 3: 优化执行
1. [Perf] 优化 Bundle Size (Tree-shaking, 移除未使用依赖)。
2. [Perf] 解决发现的内存泄漏问题。
3. [Doc] 生成最终 API 文档和迁移指南。
4. [Log] 汇总所有 Review 日志，生成最终报告。

Step 4: 最终检查
- ❌ 是否存在任何 `console.log` 或调试代码？
- ❌ 是否有未完成的 TODO？
- ❌ 关键路径是否有性能瓶颈？

验收标准：
- Bundle Size < 50KB (gzipped)。
- 所有测试 100% 通过。
- 文档完整且无拼写错误。"
```

### 📊 单独任务执行

```bash
# 执行特定子任务
"Step 1: 上下文加载
请阅读 /refactor/infrastructure.md 和 /refactor/jsdoc-standards.md。

Step 2: TDD 执行
请实现多语言系统重构：
1. 编写测试。
2. 实现 i18n Store。
3. 确保符合 JSDoc 标准。
"

"Step 1: 上下文加载
请阅读 /refactor/box.md 和 /refactor/no-backward-compatibility.md。

Step 2: TDD 执行
请重构 Box 组件：
1. 编写 E2E 测试。
2. 实现 Store 和 Directive。
3. 确保集成 HTMX。
"
```

## ⚠️ 重要提醒

### 执行原则

1. **严格按文档执行**：不得偏离设计文档的架构标准
2. **测试驱动**：每个功能必须有对应的测试
3. **质量优先**：代码质量和可维护性优先于速度
4. **不向后兼容**：严格遵循 `no-backward-compatibility.md`，不保留历史包袱
5. **完整 Review**：每个阶段必须进行全面的自我 Review

### 质量标准

- 代码覆盖率 > 90%
- E2E 测试通过率 100%
- 性能提升达到预期目标
- A11y 合规率 100%
- 无任何 Legacy 代码残留

### Review 要求

- 每个阶段完成后必须进行自我 Review
- Review 结果必须记录到 `/refactor/review-logs.md`
- 发现问题必须立即修复
- 不符合标准的需要二次确认和修改

## 🎯 重构亮点

### 新增基础设施

- **Z-Index 管理**：统一的层级管理，避免组件遮挡问题
- **无障碍访问 (A11y)**：完整的屏幕阅读器支持和键盘导航
- **错误边界**：全局错误捕获，防止组件错误导致页面崩溃
- **性能工具**：RAF 防抖、内存管理、性能优化

### 🚫 不向后兼容原则

- **彻底重构**：完全移除历史包袱，不保留任何兼容代码
- **面向最佳实践**：完全拥抱 Alpine + HTMX + Tailwind CSS 生态
- **架构优先**：架构一致性优先于向后兼容
- **长期维护**：为更好的未来而重构

### 敏捷交付方式

- **阶段独立**：每个阶段都可以独立交付和使用
- **测试驱动**：每个阶段都有完整的 E2E 测试和手动测试
- **文档同步**：每个阶段都更新相关文档
- **风险可控**：每个阶段都可以独立回滚

### 架构标准

- **Store-First, Directive-Driven**：统一的核心架构模式

- **四级配置合并**：Plugin Defaults -> Global User Config -> Global Plugin Config -> Direct Global Config -> Instance Config

- **样式与模板分离**：通过 Uno.getStyle 获取样式

- **HTMX 原生集成**：内容加载全权交给 HTMX

- **统一方法命名**：按组件类型规范命名 Store 方法

- **JSDoc 注释标准**：统一的英文注释和文档标准

## 📊 预期收益

### 技术指标

- **代码体积减少** 20%+
- **运行时性能提升** 30%+
- **内存使用减少** 25%+
- **测试覆盖率** > 90%
- **Z-Index 层级冲突** 0%
- **错误捕获率** 100%

### 用户体验指标

- **页面加载速度提升** 20%+
- **交互响应时间** < 100ms
- **错误率** < 0.1%
- **A11y 合规率** 100%
- **多语言支持完整度** 100%

### 开发效率指标

- **新组件开发时间减少** 50%+
- **Bug 修复时间减少** 40%+
- **代码可维护性评分** > 8/10
- **组件重构标准化程度** > 95%

## 🔄 实施流程 (敏捷交付)

### 阶段独立交付

每个阶段都可以独立完成开发、测试、文档更新并交付使用：

1. **第一阶段** - 基础设施重构 (3-4天)
   - 独立交付完整的基础设施包
   - 包含 Z-Index、A11y、错误边界等系统

2. **第二阶段** - 试点组件重构 (2-3天)
   - 独立交付 Notify 组件
   - 验证新架构可行性

3. **第三阶段** - 核心组件重构 (4-5天)
   - 独立交付 Box、Tooltip、Tab 组件
   - 建立组件重构标准

4. **第四阶段** - 辅助组件重构 (2-3天)
   - 独立交付 Valid 组件和整体优化
   - 完善所有功能

5. **第五阶段** - 最终优化和总结 (1-2天)
   - 独立交付完整框架
   - 发布就绪版本

### 质量保证

- **每个阶段**：独立完成 E2E 测试、手动测试、文档更新
- **风险控制**：每个阶段都可独立回滚
- **持续交付**：每个阶段都有可用的交付物

---

_重构完成后，Uno UI 将成为一个现代化、高性能、易维护的前端组件库。_
