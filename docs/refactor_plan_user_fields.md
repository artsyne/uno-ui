# 用户字段标准化重构方案 (已完成)

本以此文档详细阐述将项目中的用户相关字段进行标准化和优化的完整方案。核心目标是将 `user_name` 重命名为 `handle`，`full_name` 重命名为 `display_name`，并消除数据冗余，统一通过关联查询获取用户显示信息。

**状态**: ✅ 所有重构任务已于 2026-01-28 完成。

## 1. 核心目标

- **字段重命名**:
  - `user_name` -> `handle` (用户唯一标识符/句柄)
  - `full_name` -> `display_name` (用户显示名称)
- **去冗余**: 移除所有业务表中冗余存储的 `user_name` 字段。
- **关联查询**: 强制要求显示用户信息的场景必须关联 `br_user` 表。

## 2. 数据库层重构 (Database) [已完成]

### 2.1 Schema 变更

在 `sql/schema.sql` 中已执行以下变更：

1.  **表 `br_user`**:
    - `CHANGE COLUMN user_name handle varchar(100) ...`
    - `CHANGE COLUMN full_name display_name varchar(128) ...`
    - 更新索引名称: `user_name_UNIQUE` -> `handle_UNIQUE`。

2.  **表 `br_comment`**:
    - `DROP COLUMN add_user_name`

3.  **表 `br_diary`**:
    - `DROP COLUMN add_user_name`
    - `DROP COLUMN last_action_user_name`

4.  **表 `br_diary_member`**:
    - `DROP COLUMN user_name`

5.  **表 `br_entry`**:
    - `DROP COLUMN add_user_name`

### 2.2 数据迁移策略

- **开发环境**: 更新 `schema.sql` 并重置数据库（或编写 ALTER 脚本）。
- **生产环境**: 编写 SQL 迁移脚本，需注意：
  - 先修改 `br_user` 列名。
  - (可选) 在删除冗余列之前，确保代码已更新为关联查询，避免服务中断期间无法显示用户名。
  - 或者采用停机维护方式一次性执行。

## 3. 应用代码层重构 (Backend - Go) [已完成]

### 3.1 实体与模型 (Structs)

修改以下文件中的结构体定义及 JSON/Form/DB 标签：

- **`internal/domain/user/types.go`**:
  - `User`: `UserName` -> `Handle`, `FullName` -> `DisplayName`
  - `RegisterRequest`, `CheckUserNameRequest` 等 DTO 同步修改。
- **`internal/contracts/structs.go`**:
  - `AuthUser`: `UserName` -> `Handle`, `FullName` -> `DisplayName`
- **`internal/domain/entry/types.go`**:
  - `Entry`: 删除 `AddUserName`
- **`internal/domain/comment/types.go`**:
  - `Comment`: 删除 `AddUserName`

### 3.2 数据访问层 (Repository)

- **`internal/infra/repository/user.go`**:
  - 更新 CRUD SQL 语句，使用 `handle` 和 `display_name`。
  - `GetUserByUserName` -> `GetUserByHandle`。
  - `IsExistUserName` -> `IsExistHandle`。

- **关联查询改造**:
  - 修改 `EntryRepository`, `CommentRepository`, `DiaryRepository` 的查询方法（如 `GetEntry`, `ListEntries`）。
  - 使用 `JOIN br_user ON ...` 获取作者的 `handle` and `display_name`。
  - 或者在 Service 层进行数据组装（视性能需求而定，推荐 JOIN）。
  - 定义新的扩展结构体（如 `EntryWithAuthor`）或在原结构体中添加 `Author *User` 字段（非 DB 映射字段）来承载用户信息。

### 3.3 业务逻辑与控制器 (Service & Controller)

- **`internal/domain/user/service.go`**: 适配 Repository 变更。
- **`internal/app/controller/user`**:
  - `Register`: 接收 `handle` 参数。
  - `CheckUserName` -> `CheckHandle`。
- **模板函数 (`internal/app/templates/auth.go`)**:
  - `getUserName` -> `getHandle` (或保留方法名但返回 Handle)。
  - `getFullName` -> `getDisplayName`。

## 4. 前端代码层重构 (Frontend - JS/HTML) [已完成]

### 4.1 模板文件 (Templates)

搜索并替换 `website/views` 下所有 HTML 文件：

- `.User.UserName` -> `.User.Handle`
- `.User.FullName` -> `.User.DisplayName`
- `user_name` (form name) -> `handle`
- `full_name` (form name) -> `display_name`
- 更新 `hx-post` 参数和验证规则。

### 4.2 JavaScript 逻辑

- **`website/static/src/js/modules/auth.js`**:
  - 更新验证器：`/check-username` -> `/check-handle`。
  - 参数名 `user_name` -> `handle`。
- **API 调用**: 检查所有 `fetch` 或 `htmx` 调用，确保载荷字段名匹配。

### 4.3 国际化 (i18n)

- `website/locales/*.json`:
  - 虽然 Key 可以不改（如 `user.settings.username_label`），但建议检查 Value 中的占位符。
  - 如果 Key 包含字段名（如 `auth.label.username`），建议也同步更新为 `auth.label.handle` 以保持语义一致（可选）。

## 5. 接口与协议 (API) [已完成]

- **API 响应格式变更**:
  - 所有返回用户对象的接口，JSON 字段将变为 `handle` 和 `display_name`。
- **API 请求参数变更**:
  - 注册、更新资料等接口参数变更。

## 6. 测试与验证 (Testing) [已完成]

### 6.1 单元测试

- 修改 `internal/domain/user/service_test.go` 等。
- 修改 `tests/unit` 下所有涉及 User 构造的测试。

### 6.2 集成测试

- 修改 `tests/integration` 和 `tests/e2e`。
- 重点验证：
  - 注册流程（字段映射正确）。
  - 日记/评论列表显示（关联查询正确获取作者名）。
  - 个人设置页更新（更新 Handle/DisplayName）。

## 7. 执行步骤 (完成记录)

1.  **数据库**: 修改 `schema.sql`。 (已完成)
2.  **后端**:
    - 修改 Structs (Types)。 (已完成)
    - 修改 Repository (SQL)。 (已完成)
    - 修改 Service/Controller。 (已完成)
    - 修改 Templates Helper。 (已完成)
3.  **前端**:
    - 修改 HTML 模板。 (已完成)
    - 修改 JS 验证逻辑。 (已完成)
      - 修改 i18n 资源文件 (UI 文案: Username->Handle (Reverted to Username), FullName->DisplayName)。 (已完成)
        - **注意**: 最终决定在 UI 层保留 "Username/用户名" 的称呼，因为更符合用户习惯；但在代码和数据库层面使用 `Handle`。
4.  **清理**:


    - 删除冗余列的 Go 代码引用。 (已完成)

5.  **测试**:
    - 修复编译错误。 (已完成)
    - 运行测试套件。 (已完成)
    - 手动冒烟测试。 (已完成)

## 8. 注意事项

- **非兼容性变更**: 此变更将破坏现有 API 契约，必须通知所有客户端（当前主要为内部 Web 前端）。
- **性能监控**: 引入 JOIN 查询后，需关注列表页查询性能，确保 `br_user.user_id` 有正确索引（主键已有）。
