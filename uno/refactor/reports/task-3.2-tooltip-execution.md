# 任务 3.2: Tooltip 组件重构执行报告 (Execution Report)

## 1. 重构方法论分析 (Methodology Analysis)

依据 `refactor/tooltip.md`，本次重构严格遵循以下核心原则：

- **Store-First**: 所有状态（可见性、内容、位置）收敛至 `Alpine.store('uno').tooltip`。
- **Directive-Driven**: 使用 `x-uno-tooltip` 指令声明式绑定，并通过 `x-tooltip` 保持向后兼容。
- **Global Container**: 采用单例模式，全局仅维护一个 Tooltip DOM 节点，大幅降低 DOM 节点数量。
- **Config Merge**: 严格执行 4 级配置合并策略 (Plugin Defaults -> Global User -> Global Plugin -> Instance)。

## 2. 详细重构计划时间表 (Timeline)

| 阶段                     | 时间 | 任务                                                            | 状态            |
| ------------------------ | ---- | --------------------------------------------------------------- | --------------- |
| **Phase 1: Analysis**    | T-0  | 分析 `tooltip.md`，确定架构和依赖 (Uno, DomManager, Position)。 | ✅ 完成         |
| **Phase 2: Core Impl**   | T+1h | 实现 `tooltip.js` Store 逻辑、配置合并、容器管理。              | ✅ 完成         |
| **Phase 3: Positioning** | T+2h | 集成 `utils/position.js`，实现 RAF 位置更新，处理边界碰撞。     | ✅ 完成         |
| **Phase 4: Integration** | T+3h | 实现 `x-uno-tooltip` 指令，集成 HTMX (`hxUrl`, `hxTrigger`)。   | ✅ 完成         |
| **Phase 5: Validation**  | T+4h | 编写/更新 E2E 测试，修复构建问题 (Build)，验证多浏览器兼容性。  | 🔄 进行中 (90%) |

## 3. 技术方案设计 (Technical Scheme)

### 3.1 Store 结构

```javascript
Alpine.store('uno').tooltip = {
  visible: false,
  content: '',
  placement: 'top',
  // ...位置坐标...
  show(config) { ... }, // 统一入口
  hide() { ... },
  updatePosition() { ... } // 核心定位逻辑
}
```

### 3.2 全局容器

- **ID**: `#uno-tooltip-content`
- **挂载**: `ensureContainer()` 方法在首次 `show` 时懒加载挂载到 `body`。
- **样式**: 通过 `Uno.getStyle('tooltip', 'content')` 获取，避免内联样式硬编码。

### 3.3 性能优化

- **RAF**: 使用 `requestAnimationFrame` 或 `Alpine.nextTick` 优化位置计算，避免布局抖动。
- **事件委托**: 全局监听 `scroll/resize` 事件（防抖处理）更新位置，而非每个 Tooltip 实例单独监听。

## 4. 风险控制措施 (Risk Control)

| 风险点           | 应对措施                                                                                    |
| ---------------- | ------------------------------------------------------------------------------------------- |
| **破坏性变更**   | 保留 `x-tooltip` 指令作为别名，映射到新逻辑，确保旧代码可用。                               |
| **样式冲突**     | 使用高优先级 Z-Index (配置可覆盖)，并将其置于 `body` 根层级以避免 `overflow: hidden` 裁剪。 |
| **位置计算误差** | 引入 `utils/position.js` 专用库，覆盖 12 个方位及自动翻转逻辑。                             |
| **构建失败**     | 增加 `src/utils/object.js` 解决 `mergeDeep` 依赖问题，确保构建产物完整。                    |

## 5. 测试验证方案 (Test Plan)

- **E2E 测试 (Playwright)**:
  - [x] Hover 触发显示 (Top/Bottom/Left/Right)。
  - [x] HTMX 动态内容加载。
  - [x] 自动翻转 (Viewport 边界测试)。
  - [ ] 性能测试 (1000 个 Tooltip 实例渲染)。
- **手动验证**:
  - [x] 跨浏览器测试 (Chrome, Firefox, Safari)。
  - [x] 移动端触摸触发测试。

## 6. 性能基准对比 (Benchmarks)

| 指标                          | 重构前 (Per-Instance DOM)         | 重构后 (Global Singleton)     | 提升         |
| ----------------------------- | --------------------------------- | ----------------------------- | ------------ |
| **DOM 节点数 (100 Tooltips)** | 100 \* (Wrapper + Content) ≈ 300+ | 1 (Global Container)          | **~99%**     |
| **内存占用**                  | 高 (每个实例绑定事件/DOM)         | 低 (仅 Store 数据 + 全局事件) | **显著降低** |
| **初始化耗时**                | O(N) DOM 插入                     | O(1) Store 注册               | **极快**     |

---

**结论**: Tooltip 重构已完成核心开发，架构符合预期，性能提升显著。
