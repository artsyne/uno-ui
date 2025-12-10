# Hebe 布局使用指南

本文档介绍了 Hebe 项目中基于 Alpine.js + HTMX 的响应式布局系统的使用方法。

## 1. 布局结构概览

Hebe 的公共布局由以下核心部分组成：

1.  **Navbar (`views/navbar.html`)**: 顶部导航栏，包含 Logo、菜单、用户头像、主题/语言切换。
2.  **Content Area**: 页面主体内容区域，通常由具体的业务模板渲染。
3.  **Footer (`views/footer.html`)**: 底部版权信息。

### 1.1 技术栈

- **Alpine.js**: 用于处理前端交互状态（如下拉菜单、移动端侧边栏）。
- **Tailwind CSS v4**: 用于响应式布局样式。
- **HTMX**: (可选) 用于处理局部内容更新。

## 2. 如何使用布局组件

### 2.1 引入 Navbar

在 HTML 模板中，Navbar 作为一个独立的 Go 模板块 (`{{define "navbar"}}`) 存在。

```html
<!-- 在主布局文件中引入 -->
{{template "navbar" .}}
```

**关键属性**:
- `x-data="navbar"`: 绑定 Alpine.js 的 `navbar` 组件数据。
- `x-show="mobileMenuOpen"`: 控制移动端菜单的显示。

### 2.2 引入 Footer

Footer 同样是一个独立的 Go 模板块 (`{{define "footer"}}`)。

```html
<!-- 在主布局文件中引入 -->
{{template "footer" .}}
```

**关键属性**:
- `x-data="footer"`: 绑定 Alpine.js 的 `footer` 组件数据。
- `x-text="currentYear"`: 动态显示当前年份。

## 3. 开发新页面

在开发新页面时，请遵循以下布局规范：

### 3.1 页面容器

所有页面内容应包裹在标准的容器中，以确保在不同设备上的显示效果一致。

```html
<main class="container mx-auto px-4 py-8 min-h-screen">
    <!-- 页面内容 -->
</main>
```

### 3.2 响应式设计

使用 Tailwind CSS 的响应式前缀来适配不同屏幕：

- **Mobile (Default)**: 默认样式适用于移动端。
- **Tablet (`md:`)**: 平板设备样式。
- **Desktop (`lg:`)**: 桌面设备样式。

**示例**:

```html
<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
    <!-- 在移动端单列，平板双列，桌面三列 -->
</div>
```

## 4. 状态管理

布局相关的状态（如菜单展开、主题）由 Alpine.js 的 Global Stores 和 Components 管理。

- **获取主题状态**: `Alpine.store('theme').on`
- **获取当前语言**: `Alpine.store('i18n').current`
- **切换移动端菜单**: 调用 `navbar` 组件的 `toggleMobileMenu()` 方法。

## 5. 常见问题

### Q: 如何修改导航菜单项？
A: 直接修改 `views/navbar.html` 文件中的 HTML 结构。如果是动态菜单，请修改 `static/src/js/modules/layout.js` 中的数据逻辑。

### Q: 为什么修改了 JS 文件但没生效？
A: 请确保运行了前端构建命令 `npm run build:frontend` (或 `npm run build:js:dev`)，并强制刷新浏览器以清除缓存。

---
