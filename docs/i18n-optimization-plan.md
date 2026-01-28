# iDiary.ai 国际化深度重构与品牌优化方案

本文档定义了 idiary.ai 前端 UI 国际化系统的深度重构计划。本方案旨在彻底解决术语混淆问题，建立可扩展的、语义清晰的翻译键值体系，并统一品牌形象。

**执行原则**: 采用面向终态的设计方案，不考虑向后兼容。所有修改必须符合当前最佳实践标准。

## 1. 核心术语与品牌标准化 (Terminology & Branding)

### 1.1 品牌识别

- **Domain**: `idiary.ai`
- **Brand Voice**: 现代、智能、简洁、温暖。
- **Privacy First**: 在涉及个人页面的命名和显示上，优先考虑隐私保护，避免暴露敏感信息。

### 1.2 核心术语映射 (Core Term Mapping)

彻底废弃 "Notebook" 概念在 UI 层的展示，采用更符合用户直觉的 "Diary" (日记本) 和 "Entry" (日记篇章) 二元体系。同时扩展其他核心业务对象的术语定义。

| 概念实体 (Entity)   | 旧术语 (Old Terms)     | **新标准术语 (New Standard)**   | 说明 (Description)                 |
| :------------------ | :--------------------- | :------------------------------ | :--------------------------------- |
| **日记本 (容器)**   | Notebook / Diary       | **Diary**                       | 用户创建的用于存放内容的集合容器。 |
| **日记篇章 (内容)** | Diary / Post           | **Entry**                       | 日记本中的具体文章/记录。          |
| **新建日记本**      | Add Notebook           | **New Diary**                   | 创建新容器的动作。                 |
| **写日记**          | Write Diary / Post     | **Write Entry** / **New Entry** | 创建新内容的动作。                 |
| **我的日记列表**    | My Diaries             | **Diaries List**                | 用户的日记本列表页。               |
| **日记详情**        | Diary Detail           | **Entry Detail**                | 查看单篇日记的页面。               |
| **评论/留言**       | Comments               | **Comments**                    | 用户互动内容。                     |
| **点赞**            | Like                   | **Like**                        | 社交互动动作。                     |
| **个人资料**        | Profile                | **Profile**                     | 用户个人信息页。                   |
| **用户设置**        | User Settings          | **Settings**                    | 账户级设置。                       |
| **作者**            | Author                 | **Author**                      | 内容创作者。                       |
| **保存操作**        | Save                   | **Save**                        | 通用保存动作。                     |
| **编辑操作**        | Edit                   | **Edit**                        | 通用编辑动作。                     |
| **网络异常**        | Net Error              | **Network Error**               | 网络请求失败的通用提示。           |
| **操作反馈**        | Operation Success/Fail | **Success / Failed**            | 通用操作结果反馈。                 |

## 2. 页面与路由命名优化 (Page & Route Optimization)

### 2.1 路由重构 (Route Refactoring)

为了实现语义统一，我们将对前端路由 (Browser URL) 和 API 路由进行标准化重构。

| 类型 (Type)      | 旧路由 (Old Route) | **新路由 (New Route)** | 说明 (Note)                               |
| :--------------- | :----------------- | :--------------------- | :---------------------------------------- |
| **View: List**   | `/d/{id}`          | **/diary/{id}**        | 日记本列表页，使用全称增强可读性。        |
| **View: Detail** | `/p/{id}`          | **/entry/{id}**        | 单篇日记详情页，`p` (post) 改为 `entry`。 |
| **API: Create**  | `POST /post`       | **POST /entry**        | 创建日记篇章。                            |
| **API: Update**  | `PUT /post/{id}`   | **PUT /entry/{id}**    | 更新日记篇章。                            |
| **API: Form**    | `/form/post`       | **/form/entry**        | 获取日记编辑表单片段。                    |
| **API: Form**    | `/form/diary`      | **/form/diary**        | 获取日记本设置表单片段。                  |
| **View: User**   | `/u/{id}`          | **/u/{id}**            | (保持不变) 用户主页简写符合惯例。         |
| **List**         | `/my-diaries`      | **/diaries**           | 我的日记列表页。                          |

### 2.2 前端 JS 方法与变量重命名 (JS Refactoring)

在 `website/static/src/js/modules/diary.js` 及相关模块中，需彻底清理 `Post` 术语，并优化 `Diary` 相关操作的命名以区分"容器"与"内容"。

- **方法重命名 (Method Renaming)**:
  - `openPostModal` $\rightarrow$ `openEntryModal`
  - `openCreatePost` $\rightarrow$ `openCreateEntry`
  - `openEditPost` $\rightarrow$ `openEditEntry`
  - `_renderPostCard` $\rightarrow$ `_renderEntryCard`
  - `_handlePostSubmit` $\rightarrow$ `_handleEntrySubmit`
  - `_updatePostStats` $\rightarrow$ `_updateEntryStats`
  - `openEditDiary` $\rightarrow$ `openDiarySettings` (明确语义：这是修改日记本设置，而非内容)
  - `_handleDiarySubmit` $\rightarrow$ `_handleDiarySettingsSubmit`

- **变量与参数重命名 (Variable Renaming)**:
  - **数据字段**: API 响应及数据对象中的 `post_id` $\rightarrow$ `entry_id` (需后端配合，或在前端适配层转换)。
  - **DOM ID**: `#post-card-{id}` $\rightarrow$ `#entry-card-{id}`。
  - **DOM ID**: `#post-form-container` $\rightarrow$ `#entry-form-container`。
  - **URL Params**: `?post_id=...` $\rightarrow$ `?entry_id=...`。

### 2.3 涉及范围

- **Page Titles**: 页面 `<title>` 及头部导航。
- **UI Components**: 侧边栏、个人主页卡片标题。
- **Routes**: 检查路由定义。
- **API Interfaces**: 确保后端返回的字段描述与前端一致。

## 3. 多语言 Key 全面重构体系 (Key Refactoring Architecture)

建立统一的、层级分明的 Key 命名规范，**彻底清理**所有包含 `notebook` 等具体业务实现的词汇。

### 3.1 命名空间结构 (Namespace Structure)

新的 Key 结构应遵循 `[Context].[Object].[Action/Attribute]` 的模式。

```json
{
  "common": { ... },      // 通用词汇 (Save, Edit, Cancel)
  "page": {               // 页面级标题
    "title": {
      "diaries_list": "Diaries List",
      "entry_detail": "Entry Detail",
      "user_settings": "Settings"
    }
  },
  "diary": {              // 实体：日记本 (The Container)
    "label": "Diary",
    "create": "New Diary",
    "settings": "Diary Settings",
    "list": { "title": "Diaries List" }
  },
  "entry": {              // 实体：篇章 (The Content)
    "label": "Entry",
    "create": "New Entry",
    "edit": "Edit Entry",
    "delete": "Delete Entry",
    "form": { "title": "Post Entry", ... }
  },
  "comment": { ... },     // 实体：评论
  "auth": { ... },        // 认证相关
  "member": { ... }       // 会员展示相关
}
```

### 3.2 迁移策略与规则 (Migration Strategy & Rules)

不再列举具体的 Key 映射表，而是执行以下**自动化迁移策略**：

1.  **对象分离原则 (Object Separation)**:
    - 识别 Key 的上下文。如果 Key 描述的是容器（如标题、描述、隐私设置、封面），迁移至 `diary.*`。
    - 如果 Key 描述的是内容（如正文、天气、心情、发表），迁移至 `entry.*`。

2.  **关键词映射规则 (Keyword Mapping)**:
    - 包含 `notebook` 的 Key $\rightarrow$ 强制迁移至 `diary` 命名空间，并移除 `notebook` 字样。
    - 包含 `post` 的 Key $\rightarrow$ 强制迁移至 `entry` 命名空间。
    - 原 `diary.*` 下涉及 "写/编辑/删除" 具体篇章的 Key $\rightarrow$ 迁移至 `entry.*`。

3.  **动作标准化 (Action Standardization)**:
    - `add` / `new` $\rightarrow$ 统一为 `create` (e.g., `diary.create`).
    - `write` $\rightarrow$ 统一为 `create` 或 `write` (保持一致性).
    - `update` / `modify` $\rightarrow$ 统一为 `edit`.
    - `remove` $\rightarrow$ 统一为 `delete`.

4.  **层级扁平化 (Flattening)**:
    - 控制 Key 的深度不超过 3 层 (e.g., `Context.Object.Attribute`)。

## 4. 文件结构与代码调整 (File Structure & Code Changes)

### 4.1 前端语言文件拆分 (Frontend Locale Splitting)

为了与后端模板的多语言策略保持一致，并支持未来扩展，前端 JS 语言文件将进行拆分。

- **当前结构**:
  - `website/static/src/js/locales/data.js`: 一个文件包含所有语言对象 (`export default { en: {...}, 'zh-cn': {...} }`)。

- **目标结构**:
  采用 ES Modules 分离语言文件，通过 `index.js` 统一聚合导出。
  1.  `website/static/src/js/locales/en.js`:
      ```javascript
      export default {
        'common.save': 'Save',
        // ...
      };
      ```
  2.  `website/static/src/js/locales/zh-cn.js`:
      ```javascript
      export default {
        'common.save': '保存',
        // ...
      };
      ```
  3.  `website/static/src/js/locales/index.js`:

      ```javascript
      import en from './en.js';
      import zhCn from './zh-cn.js';

      export default {
        en,
        'zh-cn': zhCn,
      };
      ```

### 4.2 涉及文件范围

1.  **Locale Files (Definition)**:
    - `website/locales/*.json` (Backend): 对应修改后端翻译 Key。
    - `website/static/src/js/locales/*.js` (Frontend): 执行上述拆分和 Key 重构。
2.  **View Templates (Usage)**:
    - `website/views/**/*.html`: 全局搜索并替换旧 Key。
3.  **Frontend Logic (Usage)**:
    - `website/static/src/js/**/*.js`: 更新 `$t()` 或 `getMessage()` 调用。

### 4.3 文件重命名与模块拆分 (File Renaming & Splitting)

基于 "Object Separation" 原则，对 Views 和 JS 模块进行物理文件层面的重构。

#### 4.3.1 Views 重命名

| 原文件/目录 (Original)           | **新文件/目录 (New)**                 | 说明 (Note)                    |
| :------------------------------- | :------------------------------------ | :----------------------------- |
| `website/views/post/`            | **`website/views/entry/`**            | 目录重命名，存放篇章相关视图。 |
| `website/views/post/detail.html` | **`website/views/entry/detail.html`** | 篇章详情页视图。               |
| `website/views/form/post.html`   | **`website/views/form/entry.html`**   | 篇章编辑表单片段。             |

#### 4.3.2 JS 模块重构

| 原文件 (Original)     | **新文件 (New)**                 | 核心组件重命名 (Component Renaming)                                                                          |
| :-------------------- | :------------------------------- | :----------------------------------------------------------------------------------------------------------- |
| `js/modules/diary.js` | **`js/modules/diary-detail.js`** | `DiaryList` $\rightarrow$ **`DiaryDetail`** (负责日记本详情页逻辑)                                           |
| (New File)            | **`js/modules/entry.js`**        | (新建) 提取原 `diary.js` 中 `_renderPostCard` 等与 Entry 渲染高度相关的逻辑，或作为 Entry 详情页的独立模块。 |

_注意：JS 文件重命名后，需同步更新 `website/static/src/js/core/init.js` 中的引用和 `Alpine.data` 注册。_

## 5. 实施标准与测试 (Standards & Testing)

### 5.1 代码标准

- **面向终态**: 不保留 `notebook` 相关的 alias，代码中只允许出现 `diary` (容器) 和 `entry` (内容) 的语义。
- **结构清晰**: JSON/JS 文件需按字母顺序或逻辑模块排序，保持整洁。

### 5.2 测试计划

1.  **静态检查**: 使用 Grep 扫描代码库，确保没有残留的 `notebook` 关键字用于 UI 显示。
2.  **功能测试**:
    - **主要流程**: 创建日记本 -> 写日记 -> 查看列表 -> 编辑/删除。
3.  **回归测试**: 确保未涉及的模块（如 Auth, User Settings）未受影响。

---

_文档生成日期: 2026-01-19_
