### Uno UI 语义化样式与插件扩展系统重构

**Role:** 你是一位精通 Alpine.js、Tailwind CSS 和架构设计的资深前端专家。
**Task:** 请对 Uno UI 项目进行深度重构，目标是建立一套“插件级”的样式配置方案，消除源码中的硬编码类名。

---

#### 1. 核心层重构 (src/index.js & src/core/uno.js)

在 `Uno` 核心类中构建样式管理中心：

- **初始化配置**: 在 `this.options` 中引入 `ui` 属性，结构如下：

```javascript
ui: {
  base: { rounded: 'rounded-lg', primary: 'bg-blue-600', shadow: 'shadow-md' },
  box: {}, notify: {}, tooltip: {}, tab: {}, valid: {}
}
```

- **样式注入方法**: 实现核心方法 `getStyle(comp, element, defaultCls)`：
  1. **优先查找配置**: 检查 `this.options.ui[comp][element]`。如果存在且为字符串，直接返回（覆盖模式）。
  2. **全局因子回退 (由调用方处理)**: `getStyle` 本身不自动匹配 `base`，而是要求调用方（组件）在构建 `defaultCls` 时主动引用 `ui.base` 中的变量。
  3. **返回结果**: 返回最终的 Tailwind 类名字符串。

  _示例实现_:

  ```javascript
  getStyle(comp, elem, defaultCls) {
      const config = this.options.ui?.[comp]?.[elem];
      return config !== undefined ? config : defaultCls;
  }
  ```

#### 2. 插件层重构 (逻辑替换)

将以下插件中的 HTML 模板字符串及动态类名绑定进行重构。**关键原则**: 组件内部必须先读取 `ui.base` 构建默认样式，再调用 `getStyle` 允许用户覆盖。

- **Box (`src/components/box/box.js`)**:
  - **Overlay**: `getStyle('box', 'overlay', 'fixed inset-0 bg-black ...')`
  - **Container**: 构建默认值时使用 `ui.base.rounded` 和 `ui.base.shadow`。
  - **Buttons**: `okBtn` 和 `cancelBtn` 的默认样式应引用 `ui.base.primary` (对于 OK 按钮)。
  - **默认值保持原样**，但允许通过 `ui.box.okBtn` 彻底更换外观。

- **Notify (`src/components/notify/notify.js`)**:
  - **类型映射**: 将 `success/error/info/warning` 的配色映射表抽离至 `ui.notify.types`。
  - **容器**: 引用 `ui.base.rounded`。

- **Tooltip (`src/components/tooltip/tooltip.js`)**:
  - **样式抽离**: 将黑色背景、文字大小、圆角抽离至 `ui.tooltip.content`。
  - **默认引用**: 默认引用 `ui.base.rounded`。

- **Tab (`src/components/tab/tab.js`)**:
  - **彻底移除** `.current` 类名的硬编码。
  - **配置项**: 引入 `ui.tab.triggerActive` (默认: `'border-b-2 border-blue-500'`) 和 `ui.tab.triggerInactive`。
  - **逻辑修改**: 修改 `trigger()` 方法，根据 `isActive` 状态返回对应的配置类名。

- **Valid (`src/components/valid/valid.js`)**:
  - **样式抽离**: 将错误提示文字的类名（如 `text-red-500`）和输入框报错边框（如 `is-invalid` 或 `border-red-500`）抽离至 `ui.valid`。

#### 3. E2E 测试加固 (Playwright)

更新 `uno/e2e/` 目录下的相关 `.spec.ts` 文件：

- **新增样式验证测试**: 编写一个用例，在 `Uno.init` 时设置 `ui.base.rounded = 'rounded-none'`，验证 Box、Notify 和 Tooltip 的 DOM 节点是否均成功移除了圆角（检查 class 是否包含 `rounded-none` 或不包含 `rounded-lg`）。
- **Tab 激活态测试**: 验证当自定义 `triggerActive` 为特定的 Tailwind 类（如 `text-red-600`）时，点击切换能正确应用该类名。
- **稳定性**: 将测试中所有的 `waitForTimeout` 替换为 Web-first Assertions（如 `expect().toBeVisible()`）。

#### 4. 文档说明

- 更新 `README.md`，添加「Customization (主题定制)」章节。
- 展示如何通过 `window.UnoConfig.ui` 实现全局样式的“一键直角化”或“品牌色定义”。

---

**验证标准 (Definition of Done):**

1. 源码中不再出现硬编码的颜色（如 `bg-blue-600`）和圆角（如 `rounded-lg`）字符串（应引用 `ui.base` 或通过 `getStyle` 获取）。
2. 所有组件外观必须能通过 `UnoConfig` 完成深度定制。
3. 现有的全量 E2E 测试必须 100% 通过。
4. 在不提供任何自定义配置时，UI 表现必须与重构前完全一致。

#### 5. 手动验证场景 (Manual Testing Scenarios)

为了验证 `Uno.box` 全局 API 的可用性以及内容源加载功能，请执行以下手动测试：

**场景：Box 组件 API 调用与内容加载验证**

1.  **环境准备**:
    - 在浏览器中打开 `uno/examples/components/box.html` 文件（建议使用 Live Server 或本地静态服务器）。
    - 打开浏览器的开发者工具 (F12)，切换到 Console 面板。

2.  **测试步骤 1: 验证 AJAX 内容加载**
    - **操作**: 滚动到页面中部的 "Content Sources" 区域。
    - **操作**: 点击 "Load JSON (AJAX)" 按钮。
    - **预期结果**:
      - 屏幕中央弹出一个 Dialog 对话框。
      - 对话框标题为 "AJAX Content"。
      - 对话框内容显示从 `/examples/data/sample.json` 加载的 JSON 数据（格式化后的代码块）。
      - 控制台无报错信息 (特别是无 `Uno.box.dialog is undefined` 错误)。
    - **操作**: 点击对话框右上角的关闭按钮 (X) 或点击遮罩层关闭对话框。

3.  **测试步骤 2: 验证 DOM 内容加载**
    - **操作**: 在 "Content Sources" 区域，点击 "Load from DOM" 按钮。
    - **预期结果**:
      - 屏幕中央弹出一个 Dialog 对话框。
      - 对话框标题为 "DOM Content"。
      - 对话框内容显示 "This content is hidden in the DOM..." 等文本，且样式正确（包含粉色背景的提示框）。
      - 控制台无报错信息。

4.  **测试步骤 3: 验证 API 交互 (Confirm/Alert)**
    - **操作**: 在浏览器控制台中直接输入以下命令并回车：
      ```javascript
      Uno.box
        .confirm('Are you sure you want to proceed?', {
          title: 'Manual Test',
          okButtonText: 'Yes, do it!',
          cancelButtonText: 'No, wait'
        })
        .then((result) => console.log('Confirm result:', result))
      ```
    - **预期结果**:
      - 弹出一个 Confirm 对话框。
      - 点击 "Yes, do it!" 按钮，控制台输出 `Confirm result: true`。
      - 再次运行命令，点击 "No, wait" 按钮，控制台输出 `Confirm result: false`。

5.  **验收通过标准**:
    - 所有对话框均能正常弹出和关闭。
    - AJAX 和 DOM 内容均能正确加载和渲染。
    - `Uno.box` 对象在全局作用域下可用，且所有方法 (`dialog`, `confirm`, `alert` 等) 均能正常工作。
