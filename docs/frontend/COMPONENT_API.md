# Hebe 前端组件 API 文档

本文档详细描述了 `/website/static/src/js/modules` 下定义的 Alpine.js 组件及其 API。

## 1. 布局组件 (Layout Components)

### 1.1 Navbar

导航栏组件，负责管理顶部导航、移动端菜单、语言切换和主题切换。

- **源文件**: `static/src/js/modules/layout.js`
- **使用方式**: `x-data="navbar"`

#### 状态 (State)

_无_

#### 方法 (Methods)

| 方法名   | 参数 | 返回值 | 描述                     |
| :------- | :--- | :----- | :----------------------- |
| `init()` | 无   | `void` | 组件初始化钩子，打印日志 |

#### 依赖 (Dependencies)

- **Store**: `$store.theme` (用于主题切换)
- **Store**: `$store.i18n` (用于语言切换)

### 1.2 Footer

页脚组件，负责显示版权信息。

- **源文件**: `static/src/js/modules/layout.js`
- **使用方式**: `x-data="footer"` (可选，如果仅包含静态内容)

#### 状态 (State)

_无_

#### 方法 (Methods)

| 方法名   | 参数 | 返回值 | 描述                     |
| :------- | :--- | :----- | :----------------------- |
| `init()` | 无   | `void` | 组件初始化钩子，打印日志 |

## 2. 全局状态 (Global Stores)

### 2.1 Theme Store

管理应用的主题（亮色/暗色模式）。

- **源文件**: `static/src/js/stores/theme.js`
- **访问方式**: `$store.theme`

#### 状态 (State)

| 属性名 | 类型      | 描述                                      |
| :----- | :-------- | :---------------------------------------- |
| `on`   | `Boolean` | `true` 表示暗色模式，`false` 表示亮色模式 |

#### 方法 (Methods)

| 方法名     | 描述                                  |
| :--------- | :------------------------------------ |
| `toggle()` | 切换主题模式，并持久化到 LocalStorage |

### 2.2 I18n Store

管理应用的语言设置。

- **源文件**: `static/src/js/stores/i18n.js`
- **访问方式**: `$store.i18n`

#### 状态 (State)

| 属性名    | 类型     | 描述                               |
| :-------- | :------- | :--------------------------------- |
| `current` | `String` | 当前语言代码 (e.g., 'zh-cn', 'en') |

#### 方法 (Methods)

| 方法名      | 参数           | 描述                             |
| :---------- | :------------- | :------------------------------- |
| `set(lang)` | `lang: String` | 切换语言，设置 Cookie 并刷新页面 |

---
