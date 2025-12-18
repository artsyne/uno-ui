# Golang 服务端优化方案 (HTMX+Alpine 集成)

**状态**: 草稿 (Draft)
**目标目录**: `/Users/yasong.cy/go/src/idiary.ai/hebe/internal/app`
**背景**: 针对现代 HTMX+Alpine 架构进行重构，移除遗留技术债务，并通过 OpenAPI 实现标准化。

## 1. 执行摘要 (Executive Summary)

本计划概述了 Golang 服务端架构的优化方案，旨在支持新的 HTMX+Alpine 前端技术栈，同时确保严格遵守 RESTful 最佳实践和 OpenAPI 标准。主要目标是消除技术债务（遗留兼容层）、统一响应处理、实现泛型验证，并确保 100% 的测试覆盖率。

## 2. 第一阶段：代码审查与清理（移除遗留代码）

**目标**: 移除所有仅为兼容旧版技术栈（如旧的 jQuery/AJAX 模式、已弃用的 JSON 结构）而保留的代码路径。

### 2.1 审计与识别

- **行动**: 扫描 `internal/app/controller` 和 `internal/app/view` 以查找：
  - 已弃用的结构体字段（标记为 `deprecated` 或未使用）。
  - 处理旧版客户端的条件逻辑（例如 `if req.Version < 2 ...`）。
  - 新版 Alpine/HTMX 前端不再调用的端点。
- **交付物**: `docs/LEGACY_REMOVAL_CHECKLIST.md`

### 2.2 清理执行

- **行动**:
  - 移除识别出的遗留代码。
  - 移除 Go 结构体/处理器中残留的 `data-action` 处理逻辑（替换为 RESTful 路由）。
  - 确保 `internal/app/view/response.go` 仅支持标准化的响应格式。

## 3. 第二阶段：架构优化

### 3.1 统一响应器实现 (`internal/app/view`)

**目标**: 创建一个集中的响应处理器，强制执行符合 OpenAPI 规范的 JSON 结构和标准化的错误处理。

- **标准响应结构**:
  ```go
  type APIResponse[T any] struct {
      Success   bool   `json:"success" example:"true"`
      Message   string `json:"message,omitempty" example:"操作成功"`
      Data      T      `json:"data,omitempty"`
      RequestID string `json:"request_id" example:"req_123456"`
  }
  ```
- **标准错误结构**:
  ```go
  type APIErrorResponse struct {
      Success   bool              `json:"success" example:"false"`
      Code      int               `json:"code" example:"4001001"`
      Message   string            `json:"message" example:"验证失败"`
      Details   map[string]string `json:"details,omitempty"`
      RequestID string            `json:"request_id"`
  }
  ```
- **实现细节**:
  - 重构 `ViewService` 以使用泛型 (`T any`)。
  - 集成 `slog` 实现自动请求/响应日志记录（状态码、耗时、载荷大小）。
  - 自动支持 HTMX 交互的 `HX-Trigger` 响应头注入。

### 3.2 泛型验证中间件 (`internal/app/middleware`)

**目标**: 移除控制器中手动的 `c.ShouldBind` 检查，替换为基于 OpenAPI 3.0 标签的声明式、中间件驱动的验证系统。

- **方法**:
  - 使用 `go-playground/validator` (v10) 并配合符合 OpenAPI 的自定义标签。
  - 创建 `Validate[T any]() gin.HandlerFunc` 中间件。
- **工作流**:
  1.  中间件将请求体/查询参数绑定到结构体 `T`。
  2.  根据结构体标签验证 `T`（例如 `binding:"required,email"`）。
  3.  如果无效 -> 立即返回标准 `APIErrorResponse` (400)。
  4.  如果有效 -> 将 `T` 存储在 `gin.Context` 中（例如 `c.Set("validatedRequest", req)`）。
- **控制器简化**:
  - 控制器不再执行绑定/验证，只需获取已验证的结构体：
    ```go
    req := c.MustGet("validatedRequest").(*CreateDiaryRequest)
    ```

## 4. 第三阶段：文档规范 (OpenSpec)

**目标**: 采用 Spec-Driven Development (Fission-AI/OpenSpec) 流程管理变更。

- **规范**: 遵循 [Fission-AI/OpenSpec](https://github.com/Fission-AI/OpenSpec) 标准。
- **工作流**:
  1.  **建立基准**: 在 `openspec/specs/server-api.md` 中描述当前服务端架构和 API 规范作为 Source of Truth。
  2.  **变更管理**:
      - 所有重构任务（如"统一响应器"、"泛型验证"）必须先在 `openspec/changes/` 下创建 Proposal。
      - Proposal 需包含具体的接口定义变更、结构体修改计划。
  3.  **实施与归档**:
      - 代码实现需严格对应 Proposal 中的 Tasks。
      - 任务完成后，更新 `openspec/specs/server-api.md` 并归档 Proposal。

- **API 文档输出**:
  - 虽然开发流程采用 OpenSpec，但最终对外接口文档仍建议生成标准 OpenAPI 3.0 (YAML/JSON) 以供前端集成工具使用。
  - 在 `openspec/specs/` 中维护 API 的语义描述。

## 5. 实施计划与要求

### 5.1 分步执行

1.  **准备**: 安装更新的依赖项（validator v10, swaggo）。
2.  **核心重构**:
    - 更新 `internal/utils/errors` 以严格映射 HTTP 状态码。
    - 实现 `internal/app/view/responder.go`（新的泛型响应器）。
    - 实现 `internal/app/middleware/validator.go`。
3.  **模块重构（迭代进行）**:
    - 重构 `User` 模块（认证/设置）。
    - 重构 `Diary` 模块。
    - 重构 `Post` 模块。
    - 重构 `Comment` 模块。
4.  **测试更新**:
    - 更新 `tests/unit` 和 `tests/integration` 以反映新的响应结构。

### 5.2 实施要求

- **RESTful**: 使用正确的动词（GET, POST, PUT, DELETE）代替 `action` 参数。
- **测试**: 保持 >90% 的覆盖率。所有新中间件必须有单元测试。
- **兼容性**: 确保 `internal/app` 重构不破坏 `internal/domain` 逻辑。

## 6. 交付物

1.  **重构后的代码库**: `internal/app` 完全迁移到泛型响应器/验证器。
2.  **OpenAPI 规范**: `docs/openapi.yaml`（生成的或手动的）。
3.  **测试套件**: `/Users/yasong.cy/go/src/idiary.ai/hebe/tests` 中的所有测试通过 (Green)。
4.  **变更日志**: `docs/CHANGELOG_SERVER_OPTIMIZATION.md`。

---

**下一步**:

1.  审查本计划。
2.  批准后，我将从 **第一阶段：代码审查与清理** 开始执行。
