# Kiko UI 重构 - AI 执行指导文档

## 文档概述

本文档专门为 AI 代理执行 Kiko UI 重构任务而设计，旨在确保 AI 严格按照既定方案执行，避免自由发挥和幻觉，提高重构成功率。

**重要提醒**: 本文档与 `kiko-ui-refactor-plan.md` 配套使用，AI 必须同时遵循两个文档的所有约束。

---

## 1. 核心约束规则

### 1.1 严格禁止的行为

**🚫 绝对禁止以下行为，任何违反将被视为重构失败**：

#### 功能相关禁止
- ❌ **禁止增加新功能** - 只能进行功能平移，不得添加原有库中没有的方法和功能
- ❌ **禁止使用纯 JS 动画** - 所有动画必须基于 Tailwind CSS 4.1 实现
- ❌ **禁止保留 jQuery 依赖** - 必须完全移除所有 jQuery 相关代码
- ❌ **禁止修改现有 API 行为逻辑** - 保持所有公共方法的输入输出行为完全一致
- ❌ **禁止过度设计和功能膨胀** - 避免引入复杂的架构模式和抽象层
- ❌ **禁止创建复杂的配置管理系统** - 保持简单的配置结构

#### 技术实现禁止
- ❌ **禁止使用过时的 JavaScript 语法** - 必须使用 ES2023 语法
- ❌ **禁止忽略 ESLint 错误和警告** - 代码必须通过所有 ESLint 检查
- ❌ **禁止手动 DOM 操作替代自动绑定** - 必须使用 MutationObserver 实现自动绑定
- ❌ **禁止创建与 Tailwind CSS 冲突的自定义样式** - 严格基于 Tailwind CSS 4.1

### 1.2 必须遵循的原则

**✅ 以下原则必须严格遵循**：

#### 核心原则
- ✅ **功能平移** - 所有原有功能必须在新版本中完全保持
- ✅ **ES2023 语法** - 使用现代 JavaScript 语法和特性
- ✅ **ESLint 零错误零警告** - 代码必须通过所有静态检查
- ✅ **自动初始化** - kiko.js 加载后自动初始化，无需手动调用
- ✅ **自动感知与绑定** - 组件自动感知 HTML 属性并绑定功能
- ✅ **MutationObserver 监听 DOM 变化** - 动态内容自动获得事件绑定
- ✅ **基于 Tailwind CSS 4.1 实现动画** - 所有动画效果使用 Tailwind CSS 类

---

## 2. 技术实施约束

### 2.1 代码质量要求

#### 函数复杂度限制
```javascript
// ✅ 正确：函数长度控制
const processData = (data) => {
    // 函数体不超过 50 行
    // 参数不超过 4 个
    // 嵌套层级不超过 3 层
    // 圈复杂度不超过 8
};

// ❌ 错误：超出复杂度限制
const complexFunction = (a, b, c, d, e) => { // 参数过多
    // 超过 50 行的函数体
    // 深度嵌套的逻辑
};
```

#### ESLint 通过率要求
- **错误数量**: 0 个
- **警告数量**: 0 个
- **代码覆盖率**: ≥ 90%
- **测试通过率**: 100%

### 2.2 模块化规范

#### 导出风格
```javascript
// ✅ 正确：命名导出
export const KikoCore = {
    select: (selector) => { /* ... */ },
    addClass: (element, className) => { /* ... */ }
};

export const KikoUI = {
    notify: { /* ... */ },
    box: { /* ... */ }
};

// ❌ 错误：默认导出或混合导出
export default KikoCore;
export { KikoUI as default };
```

#### 入口文件结构
```javascript
// kiko-next/src/index.js
import { KikoCore } from './core/index.js';
import { KikoUI } from './ui/index.js';
import { KikoUtils } from './utils/index.js';
import { KikoI18n } from './i18n/index.js';

// 自动初始化
document.addEventListener('DOMContentLoaded', () => {
    KikoCore.init();
});

// 全局暴露
window.kiko = {
    ...KikoCore,
    ui: KikoUI,
    utils: KikoUtils,
    i18n: KikoI18n
};
```

#### 解构导入规范
```javascript
// ✅ 正确：解构导入
import { select, addClass, removeClass } from './core/dom.js';
import { notify, box } from './ui/index.js';

// ❌ 错误：整体导入
import * as DOM from './core/dom.js';
import UI from './ui/index.js';
```

#### 目录结构要求
```
kiko-next/
├── src/
│   ├── index.js          # 主入口文件
│   ├── core/             # 核心模块
│   │   ├── index.js      # 核心模块入口
│   │   ├── dom.js        # DOM 操作
│   │   ├── events.js     # 事件处理
│   │   ├── ajax.js       # AJAX 请求
│   │   └── auto-binder.js # 自动绑定器
│   ├── ui/               # UI 组件
│   │   ├── index.js      # UI 组件入口
│   │   ├── notify.js     # 通知组件
│   │   ├── box.js        # 对话框组件
│   │   ├── tooltip.js    # 提示组件
│   │   ├── tab.js        # 标签页组件
│   │   ├── grid.js       # 表格组件
│   │   ├── pager.js      # 页面行为管理
│   │   └── valid.js      # 验证组件
│   ├── utils/            # 工具模块
│   │   └── index.js      # 工具模块入口
│   └── i18n/             # 国际化模块
│       └── index.js      # 国际化入口
├── tests/                # 测试文件
├── dist/                 # 构建产物
└── build.js              # 构建脚本
```

### 2.3 API 替换映射

#### jQuery 到 KikoCore 的映射
```javascript
// jQuery API → KikoCore API 映射表

// DOM 选择和操作
$('selector')              → kiko.select('selector')
$(element).addClass()      → kiko.addClass(element, className)
$(element).removeClass()   → kiko.removeClass(element, className)
$(element).toggleClass()   → kiko.toggleClass(element, className)
$(element).hasClass()      → kiko.hasClass(element, className)
$(element).attr()          → kiko.attr(element, name, value)
$(element).removeAttr()    → kiko.removeAttr(element, name)
$(element).data()          → kiko.data(element, key, value)
$(element).html()          → kiko.html(element, content)
$(element).text()          → kiko.text(element, content)
$(element).val()           → kiko.val(element, value)
$(element).show()          → kiko.show(element)
$(element).hide()          → kiko.hide(element)
$(element).remove()        → kiko.remove(element)
$(element).append()        → kiko.append(element, content)
$(element).prepend()       → kiko.prepend(element, content)

// 事件处理
$(element).on()            → kiko.on(element, event, handler)
$(element).off()           → kiko.off(element, event, handler)
$(element).trigger()       → kiko.trigger(element, event)
$(document).ready()        → kiko.ready(callback)

// AJAX 请求
$.ajax()                   → kiko.ajax(options)
$.get()                    → kiko.get(url, data)
$.post()                   → kiko.post(url, data)

// 动画效果（使用 Tailwind CSS）
$(element).fadeIn()        → kiko.fadeIn(element)  // 添加 Tailwind 淡入类
$(element).fadeOut()       → kiko.fadeOut(element) // 添加 Tailwind 淡出类
$(element).slideUp()       → kiko.slideUp(element) // 添加 Tailwind 滑动类
$(element).slideDown()     → kiko.slideDown(element) // 添加 Tailwind 滑动类

// 工具方法
$.extend()                 → kiko.extend(target, source)
$.each()                   → kiko.each(array, callback)
$.map()                    → kiko.map(array, callback)
$.grep()                   → kiko.filter(array, callback)
$.inArray()                → kiko.indexOf(array, item)
$.isArray()                → Array.isArray()
$.isFunction()             → typeof fn === 'function'
$.isPlainObject()          → kiko.isPlainObject(obj)
$.trim()                   → str.trim()
```

---

## 3. 阶段执行约束

### 3.1 阶段完成标准

#### 每个阶段必须满足的条件
```markdown
## 阶段完成检查清单

### 功能对等性检查
- [ ] 所有原有功能在新实现中完全保持
- [ ] 公共 API 行为与原版本完全一致
- [ ] 所有配置选项和默认值保持不变
- [ ] 错误处理逻辑与原版本一致

### ESLint 通过检查
- [ ] 运行 `npm run lint` 无错误
- [ ] 运行 `npm run lint` 无警告
- [ ] 代码格式符合 Prettier 规范
- [ ] 所有文件通过 ESLint 检查

### 测试页面验证
- [ ] 创建对应的测试页面
- [ ] 所有功能在测试页面中正常工作
- [ ] 测试页面覆盖所有主要使用场景
- [ ] 浏览器兼容性测试通过

### 单元测试覆盖
- [ ] 核心模块测试覆盖率 ≥ 90%
- [ ] UI 组件测试覆盖率 ≥ 90%
- [ ] 所有测试用例通过
- [ ] 边界条件测试完整

### API 文档更新
- [ ] 更新 API 文档
- [ ] 添加使用示例
- [ ] 更新迁移指南
- [ ] 文档与实现保持同步
```

### 3.2 测试要求

#### 测试对象
- **编译后的 `kiko.js`** - 对最终构建产物进行测试
- **各个模块** - 对每个独立模块进行单元测试
- **集成功能** - 对模块间交互进行集成测试

#### 测试类型
```javascript
// 单元测试示例
describe('KikoCore.select', () => {
    test('should select elements by CSS selector', () => {
        document.body.innerHTML = '<div class="test">Test</div>';
        const elements = kiko.select('.test');
        expect(elements.length).toBe(1);
        expect(elements[0].textContent).toBe('Test');
    });

    test('should return empty array for non-existent selector', () => {
        const elements = kiko.select('.non-existent');
        expect(elements.length).toBe(0);
    });
});

// 集成测试示例
describe('KikoUI.notify integration', () => {
    test('should show and hide notification', async () => {
        kiko.ui.notify.success('Test message');
        const notification = kiko.select('.kiko-notify')[0];
        expect(notification).toBeTruthy();
        expect(notification.textContent).toContain('Test message');
        
        // 测试自动隐藏
        await new Promise(resolve => setTimeout(resolve, 3000));
        expect(kiko.select('.kiko-notify').length).toBe(0);
    });
});
```

#### 覆盖率要求
- **核心模块**: 100% 覆盖率
- **UI 组件**: ≥ 90% 覆盖率
- **工具模块**: ≥ 90% 覆盖率
- **整体项目**: ≥ 90% 覆盖率

#### 浏览器兼容性
- **Chrome 110+**: 完全支持
- **Edge 110+**: 完全支持
- **Firefox 128 ESR+**: 完全支持
- **Safari 16+**: 基本支持

---

## 4. 防止幻觉的具体提示

### 4.1 功能边界提示

#### 严格参考原有组件功能列表
```javascript
// ✅ 正确：严格按照原有功能实现
// 原有 notify.js 功能：success, error, warning, info, loading, progress
const KikoNotify = {
    success: (message, options) => { /* 实现成功通知 */ },
    error: (message, options) => { /* 实现错误通知 */ },
    warning: (message, options) => { /* 实现警告通知 */ },
    info: (message, options) => { /* 实现信息通知 */ },
    loading: (message, options) => { /* 实现加载通知 */ },
    progress: (percent, options) => { /* 实现进度通知 */ }
};

// ❌ 错误：添加原有库中没有的功能
const KikoNotify = {
    // ... 原有功能
    toast: () => { /* 新增的 toast 功能 - 禁止 */ },
    snackbar: () => { /* 新增的 snackbar 功能 - 禁止 */ },
    modal: () => { /* 新增的 modal 功能 - 禁止 */ }
};
```

#### 功能对照检查表
```markdown
## 组件功能对照表

### notify.js 组件
- [ ] success() - 成功消息通知
- [ ] error() - 错误消息通知
- [ ] warning() - 警告消息通知
- [ ] info() - 信息消息通知
- [ ] loading() - 加载状态显示
- [ ] progress() - 进度条显示
- [ ] hide() - 隐藏通知
- [ ] hideAll() - 隐藏所有通知

### box.js 组件
- [ ] alert() - 警告对话框
- [ ] confirm() - 确认对话框
- [ ] prompt() - 输入对话框
- [ ] dialog() - 自定义对话框
- [ ] iframe() - iframe 对话框
- [ ] ajax() - AJAX 内容对话框
- [ ] close() - 关闭对话框
- [ ] closeAll() - 关闭所有对话框

### 其他组件...
```

### 4.2 动画实现约束

#### 所有动画必须使用 Tailwind CSS 4.1 类
```javascript
// ✅ 正确：使用 Tailwind CSS 动画类
const fadeIn = (element) => {
    element.classList.remove('opacity-0');
    element.classList.add('opacity-100', 'transition-opacity', 'duration-300');
};

const fadeOut = (element) => {
    element.classList.remove('opacity-100');
    element.classList.add('opacity-0', 'transition-opacity', 'duration-300');
};

const slideDown = (element) => {
    element.classList.remove('max-h-0');
    element.classList.add('max-h-screen', 'transition-all', 'duration-300');
};

// ❌ 错误：使用纯 JavaScript 动画
const fadeInWrong = (element) => {
    let opacity = 0;
    const timer = setInterval(() => {
        opacity += 0.1;
        element.style.opacity = opacity;
        if (opacity >= 1) clearInterval(timer);
    }, 30);
};

// ❌ 错误：使用 CSS 关键帧动画
const slideDownWrong = (element) => {
    element.style.animation = 'slideDown 0.3s ease-out';
};
```

#### Tailwind CSS 动画类映射
```css
/* 淡入淡出 */
.fade-in { @apply opacity-100 transition-opacity duration-300; }
.fade-out { @apply opacity-0 transition-opacity duration-300; }

/* 滑动效果 */
.slide-down { @apply max-h-screen transition-all duration-300; }
.slide-up { @apply max-h-0 transition-all duration-300; }

/* 缩放效果 */
.scale-in { @apply scale-100 transition-transform duration-300; }
.scale-out { @apply scale-0 transition-transform duration-300; }

/* 位移效果 */
.translate-in { @apply translate-x-0 transition-transform duration-300; }
.translate-out { @apply translate-x-full transition-transform duration-300; }
```

### 4.3 自动绑定实现约束

#### MutationObserver 实现模式
```javascript
// ✅ 正确：使用 MutationObserver 实现自动绑定
class AutoBinder {
    constructor() {
        this.observer = null;
        this.bindings = new Map();
    }

    init() {
        // 初始绑定现有元素
        this.bindExistingElements();
        
        // 监听 DOM 变化
        this.observer = new MutationObserver((mutations) => {
            mutations.forEach((mutation) => {
                if (mutation.type === 'childList') {
                    mutation.addedNodes.forEach((node) => {
                        if (node.nodeType === Node.ELEMENT_NODE) {
                            this.bindElement(node);
                            this.bindDescendants(node);
                        }
                    });
                }
            });
        });

        this.observer.observe(document.body, {
            childList: true,
            subtree: true
        });
    }

    bindElement(element) {
        // 检查 data-kiko-* 属性
        const kikoAttrs = Array.from(element.attributes)
            .filter(attr => attr.name.startsWith('data-kiko-'));
        
        kikoAttrs.forEach(attr => {
            const component = attr.name.replace('data-kiko-', '');
            const config = this.parseConfig(attr.value);
            this.bindComponent(element, component, config);
        });
    }

    bindComponent(element, component, config) {
        // 防重复绑定
        const bindingKey = `${component}-${element.dataset.kikoId || Math.random()}`;
        if (this.bindings.has(bindingKey)) return;
        
        // 绑定组件
        switch (component) {
            case 'tooltip':
                kiko.ui.tooltip.bind(element, config);
                break;
            case 'notify':
                kiko.ui.notify.bind(element, config);
                break;
            // ... 其他组件
        }
        
        this.bindings.set(bindingKey, { element, component, config });
    }
}

// ❌ 错误：手动绑定方式
const bindElementsManually = () => {
    document.querySelectorAll('[data-kiko-tooltip]').forEach(element => {
        kiko.ui.tooltip.bind(element);
    });
};
```

#### data-kiko-* 属性规范
```html
<!-- ✅ 正确：使用 data-kiko-* 属性 -->
<button data-kiko-tooltip="{\"text\":\"Click me\",\"position\":\"top\"}">Button</button>
<div data-kiko-notify="{\"type\":\"success\",\"message\":\"Success!\"}">Trigger</div>
<form data-kiko-valid="{\"rules\":{\"email\":\"required|email\"}}">Form</form>

<!-- ❌ 错误：使用其他属性名 -->
<button tooltip="Click me">Button</button>
<div notify-type="success">Trigger</div>
<form validate="true">Form</form>
```

#### 事件委托实现
```javascript
// ✅ 正确：使用事件委托
class EventManager {
    constructor() {
        this.delegatedEvents = new Map();
    }

    delegate(selector, event, handler) {
        const delegateHandler = (e) => {
            const target = e.target.closest(selector);
            if (target) {
                handler.call(target, e);
            }
        };
        
        document.addEventListener(event, delegateHandler);
        this.delegatedEvents.set(`${selector}-${event}`, delegateHandler);
    }

    undelegate(selector, event) {
        const key = `${selector}-${event}`;
        const handler = this.delegatedEvents.get(key);
        if (handler) {
            document.removeEventListener(event, handler);
            this.delegatedEvents.delete(key);
        }
    }
}

// 使用示例
const eventManager = new EventManager();
eventManager.delegate('[data-kiko-tooltip]', 'mouseenter', function(e) {
    kiko.ui.tooltip.show(this);
});

// ❌ 错误：直接绑定事件
document.querySelectorAll('[data-kiko-tooltip]').forEach(element => {
    element.addEventListener('mouseenter', (e) => {
        kiko.ui.tooltip.show(element);
    });
});
```

#### 防重复绑定机制
```javascript
// ✅ 正确：防重复绑定
const bindTooltip = (element, config) => {
    // 检查是否已绑定
    if (element.dataset.kikoTooltipBound === 'true') {
        return;
    }
    
    // 执行绑定逻辑
    // ...
    
    // 标记已绑定
    element.dataset.kikoTooltipBound = 'true';
};

// ❌ 错误：重复绑定
const bindTooltipWrong = (element, config) => {
    // 没有检查重复绑定，可能导致多次绑定
    element.addEventListener('mouseenter', showTooltip);
    element.addEventListener('mouseleave', hideTooltip);
};
```

---

## 5. 代码审查检查清单

### 5.1 提交前检查

```markdown
## 代码提交前检查清单

### 基础检查
- [ ] 代码编译无错误
- [ ] ESLint 检查通过（0 错误 0 警告）
- [ ] Prettier 格式化完成
- [ ] 所有测试用例通过
- [ ] 测试覆盖率达标（≥90%）

### 功能检查
- [ ] 功能与原版本完全对等
- [ ] 所有公共 API 行为一致
- [ ] 配置选项和默认值保持不变
- [ ] 错误处理逻辑正确
- [ ] 国际化文本正确

### 技术检查
- [ ] 使用 ES2023 语法
- [ ] 没有 jQuery 依赖
- [ ] 动画基于 Tailwind CSS 4.1
- [ ] 使用 MutationObserver 自动绑定
- [ ] 事件委托正确实现
- [ ] 防重复绑定机制有效

### 性能检查
- [ ] 没有内存泄漏
- [ ] 事件监听器正确清理
- [ ] DOM 操作优化
- [ ] 构建产物大小合理

### 文档检查
- [ ] API 文档更新
- [ ] 使用示例完整
- [ ] 迁移指南准确
- [ ] 注释清晰易懂
```

### 5.2 代码质量标准

#### 命名规范
```javascript
// ✅ 正确：清晰的命名
const showNotification = (message, type) => { /* ... */ };
const hideAllNotifications = () => { /* ... */ };
const validateEmailAddress = (email) => { /* ... */ };

// ❌ 错误：模糊的命名
const show = (msg, t) => { /* ... */ };
const hideAll = () => { /* ... */ };
const validate = (e) => { /* ... */ };
```

#### 注释规范
```javascript
// ✅ 正确：有意义的注释
/**
 * 显示通知消息
 * @param {string} message - 通知内容
 * @param {string} type - 通知类型 (success|error|warning|info)
 * @param {Object} options - 配置选项
 * @param {number} options.duration - 显示时长（毫秒）
 * @param {boolean} options.closable - 是否可关闭
 * @returns {HTMLElement} 通知元素
 */
const showNotification = (message, type, options = {}) => {
    // 合并默认配置
    const config = { duration: 3000, closable: true, ...options };
    
    // 创建通知元素
    const element = createNotificationElement(message, type, config);
    
    // 添加到 DOM
    document.body.appendChild(element);
    
    return element;
};

// ❌ 错误：无用的注释
// 显示通知
const showNotification = (message, type, options) => {
    // 设置配置
    const config = { ...options };
    // 创建元素
    const element = createElement();
    // 返回元素
    return element;
};
```

#### 错误处理规范
```javascript
// ✅ 正确：完整的错误处理
const loadContent = async (url) => {
    try {
        const response = await fetch(url);
        
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Failed to load content:', error);
        kiko.ui.notify.error(`加载失败: ${error.message}`);
        throw error; // 重新抛出以便上层处理
    }
};

// ❌ 错误：忽略错误处理
const loadContent = async (url) => {
    const response = await fetch(url);
    const data = await response.json();
    return data;
};
```

---

## 6. 错误处理指导

### 6.1 常见问题预防

#### 内存泄漏预防
```javascript
// ✅ 正确：清理事件监听器
class Component {
    constructor(element) {
        this.element = element;
        this.handlers = new Map();
    }
    
    addEventListener(event, handler) {
        this.element.addEventListener(event, handler);
        this.handlers.set(event, handler);
    }
    
    destroy() {
        // 清理所有事件监听器
        this.handlers.forEach((handler, event) => {
            this.element.removeEventListener(event, handler);
        });
        this.handlers.clear();
        
        // 清理 DOM 引用
        this.element = null;
    }
}

// ❌ 错误：不清理事件监听器
class Component {
    constructor(element) {
        this.element = element;
        this.element.addEventListener('click', this.handleClick);
    }
    
    // 没有 destroy 方法，造成内存泄漏
}
```

#### DOM 操作安全
```javascript
// ✅ 正确：安全的 DOM 操作
const safeQuerySelector = (selector, context = document) => {
    try {
        return context.querySelector(selector);
    } catch (error) {
        console.warn(`Invalid selector: ${selector}`, error);
        return null;
    }
};

const safeSetHTML = (element, html) => {
    if (!element || typeof html !== 'string') {
        console.warn('Invalid element or HTML content');
        return;
    }
    
    // 简单的 XSS 防护
    const sanitizedHTML = html.replace(/<script[^>]*>.*?<\/script>/gi, '');
    element.innerHTML = sanitizedHTML;
};

// ❌ 错误：不安全的 DOM 操作
const unsafeSetHTML = (element, html) => {
    element.innerHTML = html; // 可能导致 XSS
};
```

### 6.2 调试约束

#### 调试信息输出
```javascript
// ✅ 正确：条件调试输出
const DEBUG = process.env.NODE_ENV === 'development';

const debugLog = (message, data) => {
    if (DEBUG) {
        console.log(`[Kiko Debug] ${message}`, data);
    }
};

const debugWarn = (message, data) => {
    if (DEBUG) {
        console.warn(`[Kiko Warning] ${message}`, data);
    }
};

// 使用示例
debugLog('Component initialized', { element, config });
debugWarn('Deprecated API used', { method: 'oldMethod' });

// ❌ 错误：生产环境输出调试信息
console.log('Debug info'); // 会在生产环境输出
console.debug('Debug data'); // 会在生产环境输出
```

#### 错误边界处理
```javascript
// ✅ 正确：错误边界处理
const withErrorBoundary = (fn, fallback) => {
    return (...args) => {
        try {
            return fn(...args);
        } catch (error) {
            console.error('Error in function:', error);
            return typeof fallback === 'function' ? fallback(error) : fallback;
        }
    };
};

// 使用示例
const safeShowNotification = withErrorBoundary(
    (message, type) => kiko.ui.notify.show(message, type),
    (error) => {
        console.error('Failed to show notification:', error);
        // 降级处理：使用原生 alert
        alert(message);
    }
);

// ❌ 错误：不处理异常
const showNotification = (message, type) => {
    // 如果出错会导致整个应用崩溃
    kiko.ui.notify.show(message, type);
};
```

---

## 7. 成功标准验证

### 7.1 最终验收标准

```markdown
## 重构成功验收标准

### 技术标准
- [ ] 完全移除 jQuery 依赖
- [ ] 所有原有功能正常工作（功能平移）
- [ ] 基于 Tailwind CSS 4.1 实现所有动画效果
- [ ] 自动初始化机制正常工作
- [ ] 自动感知与绑定功能正常
- [ ] 动态内容支持正常
- [ ] 性能不低于原版本
- [ ] 支持目标浏览器（Chrome 110+, Edge 110+, Firefox 128 ESR+）
- [ ] 构建产物大小减少 ≥ 30%
- [ ] 版本号升级至 2.0.0

### 质量标准
- [ ] 代码覆盖率 ≥ 90%
- [ ] 所有测试页面功能正常
- [ ] 整个项目无 ESLint 错误和警告
- [ ] 所有动画基于 Tailwind CSS 4.1
- [ ] 文档完整性 ≥ 95%

### 部署标准
- [ ] 在 kiko-next 目录成功构建
- [ ] 生产环境稳定运行
- [ ] 用户体验保持一致
- [ ] 开发效率不降低
```

### 7.2 验收测试流程

```javascript
// 自动化验收测试脚本
const runAcceptanceTests = async () => {
    console.log('开始验收测试...');
    
    // 1. 构建测试
    console.log('1. 执行构建测试');
    const buildResult = await runBuild();
    if (!buildResult.success) {
        throw new Error('构建失败');
    }
    
    // 2. ESLint 检查
    console.log('2. 执行 ESLint 检查');
    const lintResult = await runESLint();
    if (lintResult.errors > 0 || lintResult.warnings > 0) {
        throw new Error(`ESLint 检查失败: ${lintResult.errors} 错误, ${lintResult.warnings} 警告`);
    }
    
    // 3. 单元测试
    console.log('3. 执行单元测试');
    const testResult = await runTests();
    if (testResult.coverage < 90) {
        throw new Error(`测试覆盖率不足: ${testResult.coverage}%`);
    }
    
    // 4. 功能测试
    console.log('4. 执行功能测试');
    const functionalResult = await runFunctionalTests();
    if (!functionalResult.success) {
        throw new Error('功能测试失败');
    }
    
    // 5. 性能测试
    console.log('5. 执行性能测试');
    const performanceResult = await runPerformanceTests();
    if (performanceResult.score < 80) {
        throw new Error(`性能测试不达标: ${performanceResult.score}`);
    }
    
    console.log('✅ 所有验收测试通过');
};
```

---

## 8. 总结

本文档为 AI 代理执行 Kiko UI 重构任务提供了全面的指导和约束。通过严格遵循这些规则和标准，可以确保重构任务的成功完成，避免常见的错误和幻觉。

### 关键要点回顾

1. **严格功能平移** - 不增加新功能，保持原有行为
2. **ES2023 + Tailwind CSS 4.1** - 使用现代技术栈
3. **自动绑定机制** - 基于 MutationObserver 实现
4. **零 ESLint 错误** - 确保代码质量
5. **完整测试覆盖** - 保证功能正确性

### 执行建议

- 每个阶段完成后进行完整的验收测试
- 遇到问题时参考本文档的约束和指导
- 保持与 `kiko-ui-refactor-plan.md` 的一致性
- 定期更新文档以反映最新的约束和标准

---

**文档版本**: v1.0  
**创建日期**: 2024年  
**适用项目**: Kiko UI 2.0.0 重构  
**目标目录**: `/Users/yasong.cy/go/src/idiary.ai/hebe/kiko-next`  
**配套文档**: `kiko-ui-refactor-plan.md`