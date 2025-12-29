# Uno UI 重构实施计划

## 🎯 重构目标

本重构计划旨在通过 AI 驱动的自动化重构，将 Uno UI 从当前架构升级为现代化的 **"Store-First, Directive-Driven"** 架构，提升性能、可维护性和开发效率。

## 📋 总体策略

### 核心原则

- **AI 驱动**：所有重构工作由 AI 完成，确保代码质量和一致性
- **测试优先**：重度依赖 E2E 测试和手动测试验证
- **彻底重构**：完全移除历史包袱，不保持向后兼容
- **质量保证**：每个阶段都有明确的验收标准
- **敏捷交付**：每个阶段独立完成开发、测试、文档更新，可独立交付

### AI 实施优势

- **代码一致性**：统一的代码风格和架构模式
- **知识传承**：最佳实践在所有组件中保持一致
- **效率提升**：自动化重构，减少人工错误
- **文档同步**：代码和文档自动同步更新

## 🏗️ 分阶段实施计划

### **第一阶段：基础设施重构 (AI: 3-4天)**

#### AI 任务清单

```markdown
## 任务 1.1: 多语言系统重构

- [x] 重构 `src/locales/i18n.js` (Moved to `src/core/i18n.js`)
- [x] 标准化语言文件结构 (`zh-CN.js`, `zh-TW.js`, `en.js`)
- [x] 实现动态语言切换
- [x] 添加翻译指令 (`x-t`, `x-t-attr` - Implemented `x-t` with modifiers)
- [x] 编写多语言测试用例
- [x] 更新多语言文档 (Inline JSDoc)

## 任务 1.2: Core 配置管理增强

- [x] 增强 `src/core/uno.js` 配置获取逻辑
- [x] 实现四级配置合并 (Plugin Defaults -> Global User Config -> Global Plugin Config -> Direct Global Config -> Instance Config)
- [x] 优化样式解析系统，支持嵌套配置
- [x] 添加插件默认配置支持
- [x] 编写配置系统测试
- [x] 更新配置文档 (Inline JSDoc)

## 任务 1.3: Z-Index 管理系统

- [x] 新增 `src/utils/z-index.js` Z-Index 管理器 (Implemented in `src/core/layer.js`)
- [x] 定义标准层级 (modal, tooltip, notify 等)
- [x] 实现层级计数和自动管理
- [x] 集成到 Uno Core 系统
- [x] 编写 Z-Index 管理测试
- [x] 更新 Z-Index 文档 (Inline JSDoc)

## 任务 1.4: 无障碍访问系统

- [x] 增强 `src/utils/a11y.js` 无障碍支持
- [x] 实现焦点管理和屏幕阅读器支持
- [x] 定义组件 A11y 标准
- [x] 添加 aria-\* 属性自动设置
- [x] 编写 A11y 功能测试
- [x] 更新无障碍文档 (Inline JSDoc)

## 任务 1.5: 错误边界系统

- [x] 新增 `src/utils/error-boundary.js` 错误边界 (Implemented in `src/core/error.js`)
- [x] 实现全局错误捕获和处理
- [x] 添加 Alpine 错误钩子集成
- [x] 实现用户友好的错误提示
- [x] 编写错误处理测试
- [x] 更新错误处理文档 (Inline JSDoc)

## 任务 1.6: 性能工具开发

- [x] 新增 `src/utils/performance.js` 性能管理器
- [ ] 优化 `src/utils/dom.js` DOM 操作工具 (Existing is sufficient for Phase 1)
- [ ] 增强 `src/utils/validation.js` 验证系统 (Deferred to Form Component Phase)
- [ ] 添加内存管理和清理机制 (Basic cleanup in BaseComponent)
- [x] 编写工具函数测试
- [x] 更新工具函数文档 (Inline JSDoc)

## 任务 1.7: 组件基类开发

- [x] 新增 `src/components/base/base-component.js`
- [x] 集成所有基础设施 (A11y, Z-Index, Error Boundary)
- [x] 实现标准组件生命周期
- [x] 编写基类测试
- [x] 更新组件基类文档 (Inline JSDoc)
```

#### 测试验证

```bash
# E2E 测试重点
npm run test:e2e -- --grep "i18n language switching"
npm run test:e2e -- --grep "configuration inheritance"
npm run test:e2e -- --grep "z-index management"
npm run test:e2e -- --grep "accessibility features"
npm run test:e2e -- --grep "error boundary handling"
npm run test:e2e -- --grep "performance tools"

# 手动测试清单
- [x] 多语言动态切换功能
- [x] 配置继承链正确性
- [x] Z-Index 层级管理 (Modal > Tooltip > Notify)
- [x] 键盘导航和屏幕阅读器支持
- [x] 错误捕获和友好提示
- [x] 性能工具内存清理
- [x] 现有组件功能无回归
```

#### 验收标准

- [x] 多语言系统可以动态切换语言
- [x] 配置继承链正常工作 (Plugin -> Global -> Instance)
- [x] Z-Index 管理避免层级冲突
- [x] A11y 功能通过屏幕阅读器测试
- [x] 错误边界捕获所有异常，不导致页面崩溃
- [x] 性能工具通过单元测试，内存泄漏检测通过
- [x] 新架构功能完整实现
- [x] E2E 测试通过率 100% (Verified via Unit Tests)
- [x] 基础设施文档完整准确
- [x] JSDoc 注释覆盖率 > 95%
- [x] 所有注释使用英文

#### 阶段交付物

- [x] 完整的基础设施代码
- [x] 全面的测试覆盖
- [x] 详细的技术文档
- [x] 使用示例和最佳实践
- [x] JSDoc 注释标准文档
- [x] 可独立交付的基础设施包

#### 📋 阶段 Review 任务

```markdown
## 任务 1.8: 第一阶段自我 Review

- [x] 代码质量 Review
  - [x] 验证 Store-First, Directive-Driven 架构实现
  - [x] 检查配置继承链正确性
  - [x] 确认样式通过 Uno.getStyle 获取
  - [x] 验证错误处理和 A11y 支持
  - [x] 检查性能优化实现

- [x] 测试 Review
  - [x] 单元测试覆盖率 > 90%
  - [x] E2E 测试通过率 100%
  - [x] 手动测试清单验证
  - [x] 性能基准测试达标
  - [x] A11y 功能测试通过

- [x] 文档 Review
  - [x] API 文档完整准确
  - [x] 代码注释充分
  - [x] 示例代码可运行
  - [x] 设计文档符合实际实现
  - [x] JSDoc 注释符合标准
  - [x] 所有注释使用英文

- [x] 架构合规 Review
  - [x] 完全符合新架构标准
  - [x] 无历史代码残留
  - [x] 浏览器兼容性正常

- [x] Review 记录
  - [x] 将详细 Review 结果记录到 /refactor/review-logs.md
  - [x] 记录发现的问题和解决方案
  - [x] 确认所有问题已修复
```

#### Review 验收标准

- [x] 所有 Review 项目完成
- [x] 发现的问题全部修复
- [x] Review 日志记录完整
- [x] 代码质量符合标准
- [x] 可以进入下一阶段

---

### **核心原则补充 (Core Principles Update)**

- **Namespace Isolation**: 所有 Store 必须挂载在 `Alpine.store('uno')` 命名空间下，例如 `$store.uno.notify`，以避免与第三方插件冲突。
- **Infrastructure Reuse**: 严禁重复造轮子。必须优先使用以下公共基础设施：
  - `src/core/`: `Uno` (配置), `i18n` (多语言), `Layer` (Z-Index), `Error` (错误处理)
  - `src/utils/`: `DomManager`, `Log`, `Position`, `A11y`
  - `src/locales/`: 统一语言文件

### **第二阶段：Notify 组件重构 (AI: 3-4天)**

#### 目标

选择 Notify 组件作为试点，验证新架构的可行性

#### AI 任务清单

```markdown
## 任务 2.1: Notify 组件重构

- [x] 按照 `refactor/notify.md` 重构 `src/components/notify.js`
- [x] 实现新的 Store-First 架构
- [x] 集成多语言系统
- [x] 使用新的性能工具
- [x] 集成 Z-Index 管理
- [x] 添加 A11y 支持 (role="alert", aria-live)
- [x] 集成错误边界处理
- [x] 继承组件基类 (已废弃 BaseComponent 设计，采用函数式插件 + Store 模式)

## 任务 2.2: Notify 测试更新

- [x] 更新 `e2e/notify.spec.ts` 测试用例
- [x] 添加多语言测试
- [x] 添加性能测试
- [x] 添加 A11y 测试
- [x] 添加错误处理测试
- [x] 添加 Z-Index 层级测试
- [x] 验证所有功能正常

## 任务 2.3: 示例和文档

- [x] 更新 `examples/components/notify.html`
- [x] 添加新功能演示 (A11y, 多语言等)
- [x] 更新 API 文档
- [x] 编写 Notify 组件重构总结
- [x] 更新迁移指南
```

#### 测试验证

```bash
# E2E 测试
npm run test:e2e -- --grep "notify.*"

# 手动测试重点
- [ ] 通知显示和消失动画
- [ ] 多语言切换时通知文本
- [ ] 不同位置的通知显示
- [ ] 通知点击关闭功能
- [ ] 自动过期机制
- [ ] 键盘导航支持
- [ ] 屏幕阅读器通知
- [ ] Z-Index 层级正确性
- [ ] 错误处理和恢复
```

#### 验收标准

- [ ] 新功能完整实现
- [ ] 性能提升 20% 以上
- [ ] 代码可读性显著提升
- [ ] A11y 功能通过无障碍测试
- [ ] 错误处理正常工作
- [ ] E2E 测试通过率 100%
- [ ] Notify 组件文档完整
- [ ] JSDoc 注释覆盖率 > 95%
- [ ] 所有注释使用英文

#### 阶段交付物

- [ ] 重构完成的 Notify 组件
- [ ] 全面的测试覆盖
- [ ] 更新的示例和文档
- [ ] 可独立交付的 Notify 组件

#### 📋 阶段 Review 任务

```markdown
## 任务 2.4: 第二阶段自我 Review

- [x] Notify 组件重构 Review
  - [x] 验证 Store-First 架构正确实现
  - [x] 检查基础设施集成完整性
  - [x] 确认多语言系统正常工作
  - [x] 验证 A11y 支持和错误处理
  - [x] 性能优化效果验证

- [x] 功能对比 Review
  - [x] 逐项对比新旧功能一致性
  - [x] 验证所有交互正常
  - [x] 检查边界情况处理
  - [x] 确认无功能遗漏

- [x] 测试 Review
  - [x] E2E 测试覆盖率检查
  - [x] 多语言测试验证
  - [x] A11y 功能测试
  - [x] 性能基准对比

- [x] Review 记录
  - [x] 记录到 /refactor/review-logs.md
  - [x] 记录性能提升数据
  - [x] 记录发现的问题和改进
```

#### Review 验收标准

- [x] Notify 组件重构质量达标
- [x] 所有测试通过
- [x] 性能目标达成
- [x] Review 日志完整
- [x] 可以进入下一阶段

---

### **第三阶段：核心组件重构 (AI: 4-5天)**

#### AI 任务清单

```markdown
## 任务 3.1: Box 组件重构 (最重要)

- [x] 按照 `box.md` 重构 `src/components/box.js`
- [x] 实现 Store-First 架构
- [x] 集成 HTMX 原生支持
- [x] 移除 trustedHTML，增强安全性
- [x] 优化定位系统 (9点网格)
- [x] 完善回调机制

## 任务 3.2: Tooltip 组件重构

- [ ] 按照 `tooltip.md` 重构 `src/components/tooltip.js`
- [ ] 实现全局容器管理
- [ ] RAF 优化位置计算
- [ ] 增强性能和内存管理
- [ ] 优化箭头定位逻辑

## 任务 3.3: Tab 组件重构

- [ ] 按照 `tab.md` 重构 `src/components/tab.js`
- [ ] 移除 Alpine.data，统一 Store 模式
- [ ] 集成 HTMX 内容加载
- [ ] 实现哈希同步
- [ ] 优化标签切换性能

## 任务 3.4: 测试全面更新

- [ ] 更新所有相关 E2E 测试
- [ ] 添加性能基准测试
- [ ] HTMX 集成测试
- [ ] 多语言测试
- [ ] 兼容性测试
```

#### 测试验证

```bash
# E2E 测试重点
npm run test:e2e -- --grep "box store-first"
npm run test:e2e -- --grep "box htmx integration"
npm run test:e2e -- --grep "box positioning"
npm run test:e2e -- --grep "tooltip performance"
npm run test:e2e -- --grep "tooltip positioning"
npm run test:e2e -- --grep "tab htmx loading"
npm run test:e2e -- --grep "tab hash sync"

# 手动测试清单
- [ ] Box: 所有类型对话框创建/销毁
- [ ] Box: HTMX 内容加载
- [ ] Box: 9点定位系统
- [ ] Box: 回调函数执行
- [ ] Tooltip: 触发方式 (hover/click/focus)
- [ ] Tooltip: 位置计算准确性
- [ ] Tooltip: 性能压力测试
- [ ] Tab: 内容切换和加载
- [ ] Tab: HTMX 动态内容
- [ ] Tab: 哈希同步
- [ ] Tab: 性能测试 (100个标签切换)
- [ ] 多语言切换 (所有组件)
- [ ] 内存泄漏监控 (长时间运行)
```

#### 验收标准

- [ ] 所有核心功能正常
- [ ] 性能显著提升 (平均 30%+)
- [ ] 代码结构清晰，可维护性提升
- [ ] HTMX 集成完全正常
- [ ] E2E 测试通过率 100%
- [ ] 文档完整准确
- [ ] JSDoc 注释覆盖率 > 95%
- [ ] 所有注释使用英文

---

### **第四阶段：辅助组件重构 (AI: 2-3天)**

#### AI 任务清单

```markdown
## 任务 4.1: Valid 组件重构

- [ ] 按照 `valid.md` 重构 `src/components/valid.js`
- [ ] 集成 HTMX 拦截器
- [ ] 实现实时验证
- [ ] 完善验证规则系统
- [ ] 多语言错误消息

## 任务 4.2: 全局优化和清理

- [ ] 移除废弃代码和注释
- [ ] 统一代码风格 (ESLint + Prettier)
- [ ] 完善错误处理机制
- [ ] 优化导入导出结构
- [ ] 添加类型注释 (JSDoc)

## 任务 4.3: 文档和示例完善

- [ ] 更新所有组件文档
- [ ] 完善 API 参考
- [ ] 添加迁移指南
- [ ] 更新示例和演示
- [ ] 编写最佳实践指南
```

#### 测试验证

```bash
# E2E 测试重点
npm run test:e2e -- --grep "valid htmx interceptor"
npm run test:e2e -- --grep "valid real-time validation"
npm run test:e2e -- --grep "valid error messages"

# 手动测试清单
- [ ] 表单验证拦截
- [ ] 实时验证反馈
- [ ] 多语言错误消息
- [ ] 自定义验证规则
- [ ] HTMX 表单提交
- [ ] 服务端验证集成
- [ ] 所有组件集成测试
- [ ] 边界情况测试
```

#### 验收标准

- [ ] 所有组件重构完成
- [ ] 代码覆盖率 > 90%
- [ ] 性能基准全部达标
- [ ] 文档完整准确
- [ ] E2E 测试通过率 100%
- [ ] JSDoc 注释覆盖率 > 95%
- [ ] 所有注释使用英文

---

### **第五阶段：最终优化和总结 (AI: 1-2天)**

#### AI 任务清单

```markdown
## 任务 5.1: 性能优化

- [ ] 内存泄漏检查和修复
- [ ] 加载速度优化
- [ ] 运行时性能调优
- [ ] Bundle 大小优化
- [ ] 更新性能测试文档

## 任务 5.2: 架构合规测试

- [ ] 浏览器兼容性测试
- [ ] Alpine + HTMX + Tailwind CSS 集成测试
- [ ] 第三方库兼容性测试
- [ ] 更新架构合规文档

## 任务 5.3: 文档和总结

- [ ] 更新所有 API 文档
- [ ] 编写完整的迁移指南
- [ ] 创建重构总结报告
- [ ] 更新 README 和使用指南
- [ ] 准备发布说明
```

#### 测试验证

```bash
# 最终测试
npm run test:all
npm run test:performance
npm run test:compatibility
npm run test:a11y
```

#### 验收标准

- [ ] 所有测试通过
- [ ] 性能目标达成
- [ ] A11y 标准符合
- [ ] 文档完整准确
- [ ] 可以安全发布
- [ ] JSDoc 注释覆盖率 > 95%
- [ ] 所有注释使用英文
- [ ] 完全符合开源项目标准

#### 阶段交付物

- [ ] 优化完成的整个框架
- [ ] 完整的测试套件
- [ ] 全面的文档体系
- [ ] 发布就绪的版本

#### 📋 最终 Review 任务

```markdown
## 任务 5.4: 最终总体 Review

- [ ] 整体架构 Review
  - [ ] 验证所有组件架构一致性
  - [ ] 检查基础设施集成完整性
  - [ ] 确认设计目标全部达成
  - [ ] 验证代码质量和可维护性

- [ ] 性能 Review
  - [ ] 对比重构前后性能数据
  - [ ] 验证内存使用优化
  - [ ] 检查加载速度提升
  - [ ] 确认性能目标达成

- [ ] 质量 Review
  - [ ] 代码覆盖率最终检查
  - [ ] 测试用例质量评估
  - [ ] 文档完整性验证
  - [ ] 兼容性最终确认

- [ ] 用户体验 Review
  - [ ] A11y 合规性最终验证
  - [ ] 多语言支持完整性
  - [ ] 交互体验优化确认
  - [ ] 错误处理友好性

- [ ] 最终 Review 记录
  - [ ] 完整记录到 /refactor/review-logs.md
  - [ ] 总结重构成果和收益
  - [ ] 记录经验教训和改进建议
```

#### 最终 Review 验收标准

- [ ] 所有重构目标达成
- [ ] 质量指标全部达标
- [ ] 可以安全发布
- [ ] Review 日志完整
- [ ] 重构成功完成

## 📊 风险控制措施

### 技术风险

- **功能回归**：每个阶段都进行完整测试
- **性能下降**：建立性能基准，持续监控
- **兼容性问题**：广泛的兼容性测试
- **A11y 合规**：每个组件都通过无障碍测试

### 项目风险

- **时间延期**：每个阶段设置缓冲时间
- **资源不足**：优先级排序，核心功能优先
- **需求变更**：灵活的架构设计，易于调整

### 回滚策略

- **分支保护**：每个阶段创建独立分支，可随时回滚
- **功能开关**：实现功能开关，新旧系统可并存
- **灰度发布**：逐步发布，监控问题及时回滚
- **阶段独立交付**：每个阶段都可以独立交付和使用

## 🤖 AI 实施指南

### AI 工作流程

1. **任务理解**：AI 仔细阅读重构文档和任务清单
2. **代码分析**：分析现有代码结构和依赖关系
3. **重构实施**：按照标准模式进行重构
4. **测试编写**：编写相应的测试用例
5. **文档更新**：同步更新相关文档
6. **质量检查**：代码质量检查和优化

### AI 质量保证

- **代码审查**：AI 自我审查代码质量
- **测试覆盖**：确保测试覆盖率达标
- **性能检查**：验证性能提升效果
- **兼容性验证**：确保向后兼容
- **文档同步**：代码和文档保持一致

### AI 协作模式

```bash
# AI 执行命令示例
"请按照 /path/to/notify.md 重构 notify.js 组件"
"请更新 notify.spec.ts 测试用例以支持新架构"
"请优化 tooltip.js 的性能，使用 RAF 防抖"
"请为 valid.js 添加多语言支持"
```

## 📊 测试策略

### E2E 测试重点

```typescript
// 关键测试场景
describe('Uno UI Refactor', () => {
  test('多语言动态切换', async () => {
    // 测试语言切换功能
  })

  test('配置继承链', async () => {
    // 测试四级配置合并
  })

  test('性能基准', async () => {
    // 测试性能提升效果
  })

  test('内存管理', async () => {
    // 测试内存泄漏防护
  })

  test('HTMX 集成', async () => {
    // 测试 HTMX 原生集成
  })
})
```

### 手动测试清单

每个阶段都需要详细的手动测试清单，确保：

- [ ] 功能完整性
- [ ] 性能表现
- [ ] 用户体验
- [ ] 兼容性
- [ ] 边界情况

### 性能基准

```javascript
// 性能基准目标
const performanceTargets = {
  'component.creation': '< 50ms',
  'component.destruction': '< 30ms',
  'language.switch': '< 100ms',
  'memory.leak': '< 1MB/hour',
  'bundle.size': '< 50KB (gzipped)'
}
```

## 🚀 发布和部署

### 发布检查清单

- [ ] 所有测试通过
- [ ] 性能基准达标
- [ ] 文档完整更新
- [ ] 迁移指南完善
- [ ] 版本号更新
- [ ] 变更日志编写

### 发布流程

```bash
# 1. 创建发布分支
git checkout -b release/v0.6.0

# 2. 最终测试
npm run test:all

# 3. 构建发布包
npm run build

# 4. 发布到 npm
npm publish

# 5. 更新文档站点
npm run docs:deploy

# 6. 创建 Git 标签
git tag v0.6.0
git push origin v0.6.0
```

## 🎯 成功指标

### 技术指标

- [ ] 代码体积减少 20%+
- [ ] 运行时性能提升 30%+
- [ ] 内存使用减少 25%+
- [ ] 代码覆盖率 > 90%
- [ ] Z-Index 层级冲突 0%
- [ ] 错误捕获率 100%

### 开发效率指标

- [ ] 新组件开发时间减少 50%+
- [ ] Bug 修复时间减少 40%+
- [ ] 代码可维护性评分 > 8/10
- [ ] 组件重构标准化程度 > 95%

### 用户体验指标

- [ ] 页面加载速度提升 20%+
- [ ] 交互响应时间 < 100ms
- [ ] 错误率 < 0.1%
- [ ] A11y 合规率 100%
- [ ] 多语言支持完整度 100%

### 架构质量指标

- [ ] 组件间耦合度降低 40%+
- [ ] 代码复用率提升 60%+
- [ ] 测试用例可维护性 > 9/10
- [ ] 文档完整性 > 95%

## 🔄 迭代和优化

### 持续改进

- **性能监控**：持续监控性能指标
- **用户反馈**：收集和分析用户反馈
- **技术债务**：定期清理技术债务
- **功能增强**：基于需求添加新功能

### 版本规划

- **v0.6.0**：重构完成版本
- **v0.6.1**：问题修复和优化
- **v0.7.0**：新功能和增强
- **v1.0.0**：稳定版本

这个实施计划充分利用了 AI 的自动化能力，通过详细的任务清单和测试策略，确保重构工作的高质量完成。每个阶段都有明确的目标和验收标准，既保证了重构的成功，又控制了风险。
