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

## 3. 日记业务组件 (Diary Components)

### 3.1 DiaryList

日记列表组件，负责管理日记列表页的交互，包括统计数据加载、删除日记和点赞。

- **源文件**: `static/src/js/modules/diary.js`
- **使用方式**: `x-data="diaryList"`

#### 状态 (State)

| 属性名  | 类型     | 描述                                                                                           |
| :------ | :------- | :--------------------------------------------------------------------------------------------- |
| `stats` | `Object` | 存储每篇日记的统计信息，Key 为 `postId`，Value 为包含 `view_num`, `like_num`, `liked` 等的对象 |

#### 方法 (Methods)

| 方法名                           | 参数                | 描述                                                         |
| :------------------------------- | :------------------ | :----------------------------------------------------------- |
| `init()`                         | 无                  | 初始化组件，触发 `loadStats()`，并监听 `htmx:afterSwap` 事件 |
| `loadStats()`                    | 无                  | 扫描页面上的 `data-post-id` 元素，批量获取统计数据和点赞状态 |
| `confirmDelete(diaryId, postId)` | `diaryId`, `postId` | 弹出确认框，确认后调用 `deletePost`                          |
| `deletePost(diaryId, postId)`    | `diaryId`, `postId` | 发送删除请求，成功后从 DOM 中移除对应卡片                    |
| `toggleLike(postId)`             | `postId`            | 切换点赞状态，更新 `stats` 中的数据                          |
| `openCreatePost(diaryId)`        | `diaryId`           | 打开写日记对话框 (Uno Box) 并加载表单                        |
| `openEditPost(diaryId, postId)`  | `diaryId`, `postId` | 打开编辑日记对话框 (Uno Box) 并加载表单                      |
| `openEditDiary(diaryId)`         | `diaryId`           | 打开日记本设置对话框 (Uno Box) 并加载表单                    |
| `openPostModal(url)`             | `url`               | 在模态框中加载指定 URL 的内容 (快速查看)                     |

### 3.2 LikeButton

(已集成到 `DiaryList` 的 `toggleLike` 逻辑中，作为独立组件定义保留以备后用)

- **源文件**: `static/src/js/modules/diary.js`
- **使用方式**: `x-data="likeButton(count, liked, postId)"` (目前主要通过 `DiaryList` 直接管理)

---

## 4. 表单验证 (Form Validation)

Hebe 项目使用 Uno UI 的 `valid` 插件进行表单验证。

### 4.1 验证规则配置

从 v0.8.0 开始，验证规则必须通过 `x-uno-valid` 属性以数组形式配置，不再支持 HTML5 原生属性（如 `required`, `pattern` 等）。

```html
<!-- ✅ 正确示例 -->
<input x-uno-valid="['required', { type: 'length', min: 5, max: 20 }]" />

<!-- ❌ 错误示例 (v0.8.0+ 已废弃) -->
<input x-uno-valid required minlength="5" />
```

### 4.2 自定义验证器

自定义验证器在 `static/src/js/modules/auth.js` 中注册。

- **checkUsername**: 检查用户名是否已被注册（异步）。
- **checkEmail**: 检查邮箱是否已被注册（异步）。
- **matchPwd**: 检查确认密码是否与新密码一致。

详细 API 文档请参考 `uno/doc/valid-api.md`。
