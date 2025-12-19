# Uno UI 插件安全与性能优化（最终版）

**Role:** 你是一位资深的前端架构师，擅长 Alpine.js、HTMX 和安全加固。

**Task:** 请根据以下「修改清单」对 Uno UI 项目进行重构优化。

---

## 🔴 P0 - 安全加固（必须立即执行）

### 1. Tab 组件内容净化

**位置：** `src/components/tabs.js` 的 `open()` 方法

**当前问题：** AJAX 加载的内容使用 `trusted: true`，存在 XSS 风险

**修改要求：**
```javascript
// 在 Tab 的 open() 方法中，处理 AJAX URL 时
if (ajaxUrl) {
  const content = await fetchContent(ajaxUrl);
  DomManager.setHTML(element, content, {
    trusted: false,  // ✅ 改为 false，默认净化
    alpine: Alpine
  });
}

// 在处理 HTMX URL 时同样应用
if (hxUrl) {
  const content = await requestViaHTMX(hxUrl);
  DomManager.setHTML(element, content, {
    trusted: false,  // ✅ 改为 false
    alpine: Alpine
  });
}
```

**验证：** 加载包含 `<script>alert(1)</script>` 的内容时，脚本不应执行。

---

### 2. Box 组件内容净化优化

**位置：** `src/components/box/box.js` 的 `open()` 方法

**当前问题：** 用户直接传入的 `content` 应该默认净化

**修改要求：**
```javascript
// 在 Box 组件的 open() 方法中
async open(type, options = {}) {
  // ...
  const item = {
    content: options.content || "",
    trustedHTML: options.trustedHTML === true,
    // ...
  };
  
  // ✅ 添加：如果未标记为 trusted，则净化 content
  if (!item.trustedHTML && item.content) {
    item.content = DomManager.sanitize(item.content);
    item.trustedHTML = true; // 净化后标记为已处理
  }
  
  // 其余逻辑保持不变
}
```

**说明：** 
- `prompt()` 方法中因为使用了 `escapeHTML()` 包裹用户输入，可以保持 `trustedHTML: true`
- AJAX/HTMX 加载的内容应在加载后净化

---

## 🟡 P1 - 架构完整性与内存管理

### 3. 完善 DomManager 导出

**位置：** `src/utils/dom.js`

**修改要求：**
```javascript
// 在 DomManager 对象中添加导出
export const DomManager = {
  ensureContainer(id, html, attrs, parent) { /* 现有代码 */ },
  injectStyles(id, css) { /* 现有代码 */ },
  
  // ✅ 新增导出
  setHTML: setHTML,        // 安全设置 HTML
  sanitize: sanitizeHTML,  // HTML 净化
  escape: escapeHTML       // HTML 转义
};

// 在主文件中确保挂载到全局
window.Uno = {
  dom: DomManager,
  // ...
};
```

**验证：** `window.Uno.dom.setHTML` 和 `window.Uno.dom.escape` 可正常调用。

---

### 4. 请求缓存引入 LRU 策略

**位置：** `src/utils/request.js`

**修改要求：**
```javascript
// ✅ 实现简单的 LRU 缓存类
class LRUCache {
  constructor(maxSize = 50) {
    this.maxSize = maxSize;
    this.cache = new Map();
  }
  
  get(key) {
    if (!this.cache.has(key)) return undefined;
    const value = this.cache.get(key);
    // 更新访问顺序：删除后重新插入
    this.cache.delete(key);
    this.cache.set(key, value);
    return value;
  }
  
  set(key, value) {
    if (this.cache.has(key)) {
      this.cache.delete(key);
    }
    this.cache.set(key, value);
    
    // 超出容量时删除最早的
    if (this.cache.size > this.maxSize) {
      const firstKey = this.cache.keys().next().value;
      this.cache.delete(firstKey);
    }
  }
  
  has(key) {
    return this.cache.has(key);
  }
}

// 替换原有的 Map
const requestCache = new LRUCache(50);

// 在请求方法中使用
if (cache && method === 'GET') {
  if (requestCache.has(url)) {
    return requestCache.get(url);
  }
}

// 缓存结果
if (cache && method === 'GET') {
  requestCache.set(url, result);
}
```

**验证：** 连续请求 100 个不同 URL 后，缓存大小应维持在 50 个。

---

### 5. 统一 AJAX 请求头

**位置：** `src/utils/request.js` 的 `requestViaFetch` 方法

**修改要求：**
```javascript
async function requestViaFetch(url, options = {}) {
  const {
    method = 'GET',
    headers = {},
    body,
    credentials = 'same-origin',
    signal
  } = options;
  
  // ✅ 统一添加 HTMX 风格的请求头
  const requestHeaders = {
    'HX-Request': 'true',
    'X-Requested-With': 'XMLHttpRequest',
    ...headers
  };
  
  const response = await fetch(url, {
    method,
    headers: requestHeaders,
    body,
    credentials,
    signal
  });
  
  // 解析 HX-Trigger
  parseAndEmitHXTriggers(response, options.triggerTarget);
  
  // 其余代码保持不变
}
```

**目的：** 确保后端能识别所有 AJAX 请求并返回正确的片段格式。

---

## 🟢 P2 - 性能优化与代码质量

### 6. 正则缓存 WeakMap 化

**位置：** `src/components/valid.js` 的 `validate()` 方法

**修改要求：**
```javascript
// ✅ 在文件顶部添加
const regexCache = new WeakMap();

// 在 validate() 方法中，处理 regexp 类型时
if (ruleType === "regexp" && rule.expression) {
  let regex = regexCache.get(rule);
  
  if (!regex) {
    regex = createRegex(rule.expression, { defaultFlags: "i" });
    if (regex) {
      regexCache.set(rule, regex);
    }
  }
  
  if (regex && !regex.test(value)) {
    error = rule.msg || getMessage("pattern");
  }
}
```

**说明：** 避免在 Alpine 响应式对象上存储正则，防止内存泄漏。

---

### 7. CSP Nonce 容器支持

**位置：** `src/utils/dom.js` 的 `ensureContainer` 方法

**修改要求：**
```javascript
ensureContainer(id, html = "", attrs = {}, parent = document.body) {
  if (document.getElementById(id)) return;
  
  try {
    const element = document.createElement("div");
    element.id = id;
    
    // ✅ 自动注入 CSP nonce
    if (window.Uno?.options?.csp?.nonce && !attrs.nonce) {
      attrs.nonce = window.Uno.options.csp.nonce;
    }
    
    // 设置属性
    Object.entries(attrs).forEach(([key, value]) => {
      element.setAttribute(key, value);
    });
    
    if (html) {
      element.innerHTML = html;
    }
    
    // 其余代码保持不变
  } catch (error) {
    logger.error(`Failed to inject container #${id}`, error);
  }
}
```

---

### 8. 优化初始化逻辑

**位置：** `src/index.js` 底部的初始化代码

**修改要求：**
```javascript
// ✅ 统一的初始化函数
function initializeUno() {
  if (Uno.initialized) return;
  
  const config = window.UnoConfig || {};
  logger.setDebug(config.debug !== false);
  
  if (config.debug !== false) {
    logger.info("Uno: Initializing...");
  }
  
  Uno.init(window.Alpine, config);
}

// 监听 Alpine 初始化
document.addEventListener("alpine:init", () => {
  if (!Uno.initialized) {
    initializeUno();
  }
});

// DOMContentLoaded 备用
document.addEventListener("DOMContentLoaded", () => {
  if (!Uno.initialized && window.Alpine) {
    initializeUno();
  }
});

// 同步初始化（Alpine 已加载）
if (!Uno.initialized && window.Alpine) {
  initializeUno();
}
```

**说明：** 移除重复的配置读取和 debug 设置。

---

## 📝 验证清单（Definition of Done）

修改完成后，请确保通过以下验证：

### 安全验证
```javascript
// ✓ 1. Tab XSS 防护
const tab = unoTabs('test');
await tab.open('tab1', '/api/malicious-content'); // 返回 <script>alert(1)</script>
// 期望：脚本不执行

// ✓ 2. Box XSS 防护
Alpine.store('box').open('dialog', {
  content: '<img src=x onerror=alert(1)>'
});
// 期望：onerror 被移除
```

### 性能验证
```javascript
// ✓ 3. 缓存 LRU 验证
for (let i = 0; i < 100; i++) {
  await request(`/api/test${i}`, { cache: true });
}
// 期望：缓存大小 ≤ 50

// ✓ 4. 正则缓存验证
for (let i = 0; i < 1000; i++) {
  await validForm.validate('email');
}
// 期望：性能稳定，无内存泄漏
```

### API 验证
```javascript
// ✓ 5. DomManager 导出验证
const element = document.createElement('div');
window.Uno.dom.setHTML(element, '<b>test</b>');
const escaped = window.Uno.dom.escape('<script>');
// 期望：方法正常调用

// ✓ 6. HTMX 请求头验证
// 在后端检查请求头
// 期望：包含 HX-Request: true
```

### 兼容性验证
```javascript
// ✓ 7. CSP nonce 验证
window.UnoConfig = { csp: { nonce: 'abc123' } };
Alpine.store('box').open('dialog', { content: 'Test' });
// 期望：Box 容器有 nonce="abc123"

// ✓ 8. 向后兼容验证
// 期望：所有现有 API 调用不受影响
```

---

## 🎯 执行优先级

1. **立即执行（今天）**：P0-1, P0-2（安全修复）
2. **本周完成**：P1-3, P1-4, P1-5（架构完善）
3. **有时间优化**：P2-6, P2-7, P2-8（性能提升）

---

## 📋 修改文件清单

| 文件 | 修改内容 | 优先级 | 预计时间 |
|------|---------|--------|----------|
| `src/components/tabs.js` | Tab 内容净化 | P0 | 10 分钟 |
| `src/components/box/box.js` | Box 内容净化 | P0 | 10 分钟 |
| `src/utils/dom.js` | 导出 API + CSP | P1/P2 | 15 分钟 |
| `src/utils/request.js` | LRU 缓存 + 请求头 | P1 | 20 分钟 |
| `src/components/valid.js` | 正则缓存 | P2 | 10 分钟 |
| `src/index.js` | 初始化优化 | P2 | 5 分钟 |

**总计预估时间：** 约 70 分钟

---

## ⚠️ 注意事项

1. **向后兼容**：所有修改不应破坏现有 API
2. **渐进增强**：优先完成 P0，逐步推进 P1/P2
3. **测试覆盖**：每完成一项修改，立即进行对应验证，根据需要在 uno/e2e 和 uno/examples中增加相关测试
4. **文档更新**：完成后更新 README 中的安全最佳实践部分，还有 uno/doc 中的相关文档

---

**开始重构！**