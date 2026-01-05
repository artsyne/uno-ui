# Uno UI v2.0 不向后兼容原则

## 🚫 核心原则

**Uno UI v2.0 将彻底重构，完全移除历史包袱，面向 Alpine + HTMX + Tailwind CSS 最佳实践。**

---

## 📋 不兼容变更清单

### 1. API 变更

#### ❌ 移除的 API

```javascript
// ❌ 旧版本 API（完全移除）
Uno.box.open('dialog', { title: '标题' })
Uno.notify.show('消息', { type: 'success' })
Uno.tab.show('tab1', { hash: true })
Uno.tooltip.show('提示内容')
Uno.valid.validate(form)

// ❌ 旧版本配置（完全移除）
window.UnoConfig = {
  box: { defaultWidth: 400 },
  notify: { position: 'top-right' }
}
```

#### ✅ 新版本 API

```javascript
// ✅ 新版本 API
Alpine.store('uno').box.open({ type: 'dialog', title: '标题' })
Alpine.store('uno').notify.success('消息')
Alpine.store('uno').tab.switch('tab1', { hash: true })
Alpine.store('uno').tooltip.show('提示内容')
Alpine.store('uno').valid.validate(form)

// ✅ 新版本配置
window.UnoConfig = {
  plugins: {
    box: { defaultWidth: 400 },
    notify: { position: 'top-right' }
  }
}
```

### 2. HTML 属性变更

#### ❌ 移除的属性

```html
<!-- ❌ 旧版本属性（完全移除） -->
<button data-uno-box='{"type":"dialog"}'>打开</button>
<div data-uno-notify='{"type":"success"}'>通知</div>
<button data-uno-tab='{"defaultTab":"tab1"}'>标签</button>
<span data-uno-tooltip='{"content":"提示"}'>悬停</span>
<form data-uno-valid='{"realtime":true}'>表单</form>
```

#### ✅ 新版本指令

```html
<!-- ✅ 新版本指令 -->
<button x-uno-box='{"type":"dialog"}'>打开</button>
<div x-uno-notify='{"type":"success"}'>通知</div>
<button x-uno-tab='{"defaultTab":"tab1"}'>标签</button>
<span x-uno-tooltip='{"content":"提示"}'>悬停</span>
<form x-uno-valid='{"realtime":true}'>表单</form>
```

### 3. 配置系统变更

#### ❌ 移除的配置选项

```javascript
// ❌ 旧版本配置（完全移除）
{
  ajaxUrl: '/api/data',           // 改为 hxUrl
  trustedHTML: true,             // 完全移除
  sanitizeHTML: false,           // 完全移除
  showOkButton: true,            // 移除，通过类型控制
  showCancelButton: true         // 移除，通过类型控制
}
```

#### ✅ 新版本配置

```javascript
// ✅ 新版本配置
{
  hxUrl: '/api/data',            // HTMX 标准
  type: 'dialog',               // 类型驱动
  onOk: () => {},                // 回调函数
  onCancel: () => {}             // 回调函数
}
```

### 4. 依赖要求变更

#### ❌ 移除的依赖

- 自定义 AJAX 封装
- 内置 HTML 清理函数
- 复杂的配置继承逻辑

#### ✅ 强制依赖

- Alpine.js 3.x
- HTMX 2.x
- Tailwind CSS 3.x

---

## 🎯 重构原则

### 1. 完全移除历史包袱

- 不保留任何旧版本兼容代码
- 不提供兼容性适配层
- 不支持渐进式迁移

### 2. 面向最佳实践

- 完全拥抱 Alpine + HTMX + Tailwind CSS 生态
- 遵循各框架的官方最佳实践
- 移除所有反模式实现

### 3. 架构优先

- 架构一致性优先于向后兼容
- 代码质量优先于功能完整性
- 长期维护性优先于短期便利

---

## 📚 迁移策略

### 迁移指南

提供详细的迁移指南，但不是自动兼容：

1. **API 迁移**：提供新旧 API 对应表
2. **配置迁移**：提供配置转换工具
3. **HTML 迁移**：提供属性替换指南
4. **示例迁移**：提供完整示例对比

### 迁移工具

```bash
# 提供迁移辅助工具
npm install @uno/migration-tools

# 自动转换配置
uno-migrate-config old-config.json new-config.json

# 自动转换 HTML
uno-migrate-html ./old-templates ./new-templates
```

---

## ⚠️ 重要提醒

### 对用户的影响

1. **必须更新代码**：现有代码无法直接运行在新版本上
2. **必须重新学习**：API 和使用方式有重大变更
3. **必须迁移配置**：配置系统完全重构

### 对开发者的好处

1. **更清晰的架构**：统一的开发模式
2. **更好的性能**：移除兼容层，性能更优
3. **更易维护**：代码更简洁，逻辑更清晰
4. **更好的生态**：完全拥抱主流框架生态

---

## 🚀 迁移时间表

### v1.x 停止支持

- **2025年Q2**：停止新功能开发
- **2025年Q3**：停止安全更新
- **2025年Q4**：完全停止支持

### v2.0 发布

- **2025年Q1**：Alpha 版本发布
- **2025年Q2**：Beta 版本发布
- **2025年Q3**：正式版本发布

---

_不向后兼容是为了更好的未来。彻底的重构将带来更好的开发体验和更长的生命周期。_
