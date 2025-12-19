# Uno UI 修复与优化实施计划 (基于 v0.4.0 标准)

## 1. 现状评估与问题分析

经过对现有代码库（`src/core`, `src/components`, `src/utils`）与优化文档（特别是 `uno 优化4.md`）的对比审查，发现当前项目处于 **v0.4.0 架构的半成品状态**。主要存在以下问题：

1.  **核心样式引擎不完整**: `Uno.getStyle` 仅实现了配置读取，**缺失了 `+` 追加模式**（即无法通过配置 `+class` 来追加样式而不覆盖默认值）。
2.  **组件样式硬编码**:
    *   **Notify**: 颜色样式（`bg-blue-50`, `text-blue-400` 等）仍硬编码在模板中，未映射至 `ui.notify.types`。
    *   **Tab**: 仍存在硬编码默认值，未完全走配置驱动。
3.  **API 缺失**:
    *   **Valid**: 缺失 `Uno.valid.clear()` 接口。
4.  **测试失败**: `tab.spec.ts` 中的 Deep linking 和样式定制测试失败，说明重构破坏了原有逻辑。
5.  **代码一致性**: 部分优化文档中的要求（如 `uno 优化2.md` 的安全加固）已落实，但 `uno 优化4.md` 的语义化重构尚未完成。

## 2. 修复目标

将项目完全对齐 `uno 优化4.md` 标准，修复现有 BUG，并确保所有 E2E 测试通过。

## 3. 详细实施步骤

### 阶段一：核心引擎修复 (Core Engine)

*   **目标**: 完善样式总线。
*   **文件**: `src/core/uno.js`
*   **任务**:
    *   修改 `getStyle` 方法，支持以 `+` 开头的配置值进行类名追加。
    *   确保 `ui` 配置结构完整初始化。

### 阶段二：组件语义化重构 (Components Refactor)

按顺序对插件进行“去硬编码”和“API补全”：

#### 1. Notify 组件 (`src/components/notify.js`)
*   **样式抽离**: 将 `info`, `success`, `warning`, `error` 的背景色、边框色、图标色抽离至 `Uno.options.ui.notify.types`。
*   **逻辑修改**: 模板中不再直接写 Tailwind 颜色类名，而是通过 helper 函数读取配置。

#### 2. Tab 组件 (`src/components/tab.js`)
*   **样式配置**: 确保 `active` 和 `inactive` 类名完全从 `ui.tab` 读取。
*   **BUG 修复**: 修复 `deep linking` 测试失败的问题（检查 `init` 时 hash 读取时机与 `aria-selected` 绑定）。
*   **安全加固**: 确认 `open` 方法中的内容加载默认使用 `trusted: false`。

#### 3. Valid 组件 (`src/components/valid.js`)
*   **API 新增**: 实现 `Uno.valid.clear(form)` 方法，用于重置验证状态。
*   **样式抽离**: 确保错误样式的类名（`is-invalid` 等）通过 `getStyle` 获取。

#### 4. Tooltip 组件 (`src/components/tooltip.js`)
*   **API 增强**: 确保 `show/hide` 公共 API 逻辑正确。
*   **样式联动**: 确认圆角样式引用 `ui.base.rounded`。

### 阶段三：测试加固 (Testing)

*   **文件**: `e2e/*.spec.ts`
*   **任务**:
    *   **Style Test**: 新增测试用例，验证 `+` 追加模式是否生效（如 `ui.box.okBtn = '+font-bold'`）。
    *   **Valid Test**: 新增 `Uno.valid.clear()` 的测试用例。
    *   **Fix Existing**: 修复 `tab.spec.ts` 中的失败用例。

## 4. 执行计划

我将按照以下顺序执行代码修改：

1.  **Core**: 修正 `src/core/uno.js` 中的 `getStyle`。
2.  **Tab**: 优先修复 `src/components/tab.js` 的测试错误并完成样式重构（因为这是当前报错的痛点）。
3.  **Notify**: 重构 `src/components/notify.js` 样式系统。
4.  **Valid**: 补充 `src/components/valid.js` 的 `clear` API。
5.  **Verify**: 运行 `npm run test` 验证修复结果。

---
**确认**: 请确认以上计划，我将立即开始执行。
