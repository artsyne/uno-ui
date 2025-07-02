# iDiary 错误处理和响应服务使用示例

本文档提供了iDiary项目中错误处理和响应服务的完整使用示例，展示了如何在实际业务场景中使用新的错误处理体系。

## 目录

1. [ResponseService 使用示例](#responseservice-使用示例)
2. [增强版错误处理示例](#增强版错误处理示例)
3. [实际业务场景](#实际业务场景)
4. [最佳实践指南](#最佳实践指南)

## ResponseService 使用示例

### 基本Handler结构

```go
package response

import (
	"github.com/biuren/hebe/internal/app/view"
	errors "github.com/biuren/hebe/internal/core/error"
	"github.com/gin-gonic/gin"
)

// ExampleHandler 示例Handler，展示如何使用ResponseService
type ExampleHandler struct {
	responseService ResponseServiceInterface
	viewService     view.ViewServiceInterface
}

// NewExampleHandler 创建示例Handler
func NewExampleHandler(responseService ResponseServiceInterface, viewService view.ViewServiceInterface) *ExampleHandler {
	return &ExampleHandler{
		responseService: responseService,
		viewService:     viewService,
	}
}
```

### 1. API成功响应示例

```go
// HandleAPISuccess 示例：API成功响应
func (h *ExampleHandler) HandleAPISuccess(c *gin.Context) {
	data := map[string]interface{}{
		"message": "操作成功",
		"id":      123,
		"items":   []string{"item1", "item2"},
	}

	// 使用统一的成功响应格式
	// 输出：{"code": 0, "message": "success", "data": {...}}
	h.responseService.JSONSuccess(c, data)
}
```

### 2. API错误响应示例

```go
// HandleAPIError 示例：API错误响应
func (h *ExampleHandler) HandleAPIError(c *gin.Context) {
	// 模拟业务错误
	err := errors.UserNotExists

	// 使用统一的错误响应格式
	// 自动记录错误日志，返回标准格式
	h.responseService.JSONError(c, err)
}
```

### 3. HTML页面渲染示例

```go
// HandlePageRender 示例：HTML页面渲染
func (h *ExampleHandler) HandlePageRender(c *gin.Context) {
	// 创建页面数据
	pageData := h.viewService.PageData(c, "diary.title")

	// 添加页面特定数据
	pageData.Data["diaries"] = []map[string]interface{}{
		{"id": 1, "title": "今天的日记", "content": "今天天气不错..."},
		{"id": 2, "title": "昨天的回忆", "content": "昨天去了公园..."},
	}

	// 渲染HTML页面
	h.responseService.HTML(c, "diary/index.html", pageData)
}
```

### 4. 重定向示例

```go
// HandleRedirect 示例：重定向
func (h *ExampleHandler) HandleRedirect(c *gin.Context) {
	// 重定向到登录页面
	h.responseService.Redirect(c, "/passport/login")
}
```

## 增强版错误处理示例

### 基本Handler结构

```go
package errors

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// ExampleHandler 示例Handler，展示如何使用增强版错误处理
type ExampleHandler struct {
	// 在实际应用中，这里会注入相关服务
}
```

### 1. 用户注册错误处理示例

```go
// ExampleUserRegistration 示例：用户注册中的错误处理
func (h *ExampleHandler) ExampleUserRegistration(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")

	// 1. 参数验证错误示例
	if username == "" {
		err := NewValidationError("username_required").
			WithMeta(map[string]interface{}{
				"field": "username",
				"value": username,
			}).
			WithContext("user", "register").
			WithRequestID(c.GetString("request_id")).
			WithDetails("用户名不能为空")

		// 通过c.Error()传递给中间件处理
		_ = c.Error(err)
		return
	}

	// 2. 业务错误示例：用户已存在
	if h.userExists(username) {
		err := ErrUserExists.
			WithMeta(map[string]interface{}{
				"username": username,
				"attempt":  "registration",
			}).
			WithRequestID(c.GetString("request_id"))

		_ = c.Error(err)
		return
	}

	// 3. 数据库错误示例
	if err := h.saveUser(username, email); err != nil {
		dbErr := NewDatabaseError("user_creation_failed").
			WithCause(err).
			WithMeta(map[string]interface{}{
				"username": username,
				"email":    email,
				"table":    "users",
			}).
			WithContext("user", "register").
			WithRequestID(c.GetString("request_id")).
			WithDetails("创建用户记录时发生数据库错误")

		_ = c.Error(dbErr)
		return
	}

	// 成功响应
	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data":    gin.H{"message": "用户注册成功"},
	})
}
```

### 2. 日记访问权限检查示例

```go
// ExampleDiaryAccess 示例：日记访问权限检查
func (h *ExampleHandler) ExampleDiaryAccess(c *gin.Context) {
	diaryID := c.Param("id")
	userID := c.GetString("user_id")

	// 1. 资源不存在错误
	diary := h.getDiary(diaryID)
	if diary == nil {
		err := ErrDiaryNotFound.
			WithMeta(map[string]interface{}{
				"diary_id": diaryID,
				"user_id":  userID,
			}).
			WithRequestID(c.GetString("request_id"))

		_ = c.Error(err)
		return
	}

	// 2. 权限错误示例
	if !h.canAccessDiary(userID, diaryID) {
		err := ErrDiaryAccessDenied.
			WithMeta(map[string]interface{}{
				"diary_id":     diaryID,
				"user_id":      userID,
				"diary_owner":  diary["owner_id"],
				"access_level": "read",
			}).
			WithRequestID(c.GetString("request_id")).
			WithDetails("用户无权访问此日记")

		_ = c.Error(err)
		return
	}

	// 成功返回日记内容
	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data":    diary,
	})
}
```

### 3. 外部网络请求错误处理示例

```go
// ExampleNetworkRequest 示例：外部网络请求错误处理
func (h *ExampleHandler) ExampleNetworkRequest(c *gin.Context) {
	// 模拟调用外部API
	data, err := h.callExternalAPI()
	if err != nil {
		// 包装网络错误
		networkErr := NewNetworkError("external_api_failed").
			WithCause(err).
			WithMeta(map[string]interface{}{
				"api_endpoint": "https://api.example.com/data",
				"timeout":      "30s",
				"retry_count":  3,
			}).
			WithContext("external", "api_call").
			WithRequestID(c.GetString("request_id")).
			WithDetails("调用外部API失败")

		_ = c.Error(networkErr)
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}
```

### 4. 旧错误转新错误示例

```go
// ExampleErrorConversion 示例：旧错误转新错误
func (h *ExampleHandler) ExampleErrorConversion(c *gin.Context) {
	// 假设从旧代码中获得AppError
	legacyErr := UserNotExists // 这是一个*AppError

	// 转换为增强版错误
	enhancedErr := FromLegacyError(legacyErr).
		WithRequestID(c.GetString("request_id")).
		WithUserID(c.GetString("user_id")).
		WithContext("user", "profile_access")

	_ = c.Error(enhancedErr)
}
```

### 5. 错误链式调用完整示例

```go
// ExampleErrorChaining 示例：错误链式调用
func (h *ExampleHandler) ExampleErrorChaining() {
	// 创建一个完整的错误，展示所有可用的链式方法
	err := NewUserError(CodePasswordWrong, "password_incorrect").
		WithMeta(map[string]interface{}{
			"username":      "john_doe",
			"login_attempt": 3,
			"ip_address":    "192.168.1.100",
			"user_agent":    "Mozilla/5.0...",
		}).
		WithDetails("用户输入了错误的密码，这是第3次尝试").
		WithRequestID("req-12345").
		WithUserID(12345).
		WithContext("auth", "login")

	// 可以继续添加更多信息
	if err.GetSeverity() == SeverityMedium {
		// 如果是中等严重性错误，记录额外信息
		err.WithMeta(map[string]interface{}{
			"should_alert": true,
			"alert_level":  "warning",
		})
	}

	// 使用错误
	fmt.Printf("Error: %s\n", err.Error())
	fmt.Printf("Code: %d\n", err.Code)
	fmt.Printf("Severity: %v\n", err.GetSeverity())
	fmt.Printf("Type: %v\n", err.GetType())
	fmt.Printf("Recoverable: %v\n", err.IsRecoverable())
}
```

## 实际业务场景

### 场景1：用户登录流程

```go
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErr := NewValidationError("invalid_login_request").
			WithMeta(map[string]interface{}{
				"validation_error": err.Error(),
			}).
			WithRequestID(c.GetString("request_id"))

		h.responseService.JSONError(c, validationErr)
		return
	}

	user, err := h.userService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		// Service层返回的错误已经是增强版错误
		h.responseService.JSONError(c, err)
		return
	}

	// 生成JWT token
	token, err := h.generateToken(user.ID)
	if err != nil {
		tokenErr := NewInternalError("token_generation_failed").
			WithCause(err).
			WithUserID(user.ID).
			WithRequestID(c.GetString("request_id"))

		h.responseService.JSONError(c, tokenErr)
		return
	}

	// 成功响应
	h.responseService.JSONSuccess(c, map[string]interface{}{
		"user":  user,
		"token": token,
	})
}
```

### 场景2：日记创建流程

```go
func (h *DiaryHandler) CreateDiary(c *gin.Context) {
	userID := c.GetString("user_id")

	var req CreateDiaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErr := NewValidationError("invalid_diary_data").
			WithMeta(map[string]interface{}{
				"user_id": userID,
				"error":   err.Error(),
			}).
			WithRequestID(c.GetString("request_id"))

		h.responseService.JSONError(c, validationErr)
		return
	}

	// 检查用户是否达到日记本上限
	if count, err := h.diaryService.GetUserDiaryCount(userID); err != nil {
		h.responseService.JSONError(c, err)
		return
	} else if count >= MaxDiaryLimit {
		limitErr := ErrDiaryLimitExceeded.
			WithMeta(map[string]interface{}{
				"user_id":      userID,
				"current_count": count,
				"limit":        MaxDiaryLimit,
			}).
			WithRequestID(c.GetString("request_id"))

		h.responseService.JSONError(c, limitErr)
		return
	}

	// 创建日记
	diary, err := h.diaryService.CreateDiary(userID, &req)
	if err != nil {
		h.responseService.JSONError(c, err)
		return
	}

	h.responseService.JSONSuccess(c, diary)
}
```

### 场景3：文件上传处理

```go
func (h *FileHandler) UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		uploadErr := NewValidationError("no_file_uploaded").
			WithMeta(map[string]interface{}{
				"form_field": "image",
			}).
			WithRequestID(c.GetString("request_id"))

		h.responseService.JSONError(c, uploadErr)
		return
	}
	defer file.Close()

	// 检查文件大小
	if header.Size > MaxFileSize {
		sizeErr := NewValidationError("file_too_large").
			WithMeta(map[string]interface{}{
				"file_size": header.Size,
				"max_size":  MaxFileSize,
				"filename":  header.Filename,
			}).
			WithRequestID(c.GetString("request_id"))

		h.responseService.JSONError(c, sizeErr)
		return
	}

	// 保存文件
	savedPath, err := h.fileService.SaveUploadedFile(file, header)
	if err != nil {
		// fileService会返回增强版错误
		h.responseService.JSONError(c, err)
		return
	}

	h.responseService.JSONSuccess(c, map[string]interface{}{
		"file_path": savedPath,
		"file_size": header.Size,
	})
}
```

## 最佳实践指南

### 1. ResponseService 使用规范

```go
// ✅ 正确：使用统一响应服务
func (h *Handler) GoodExample(c *gin.Context) {
	h.responseService.JSONSuccess(c, data)
}

// ❌ 错误：直接使用gin.Context
func (h *Handler) BadExample(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}
```

### 2. 错误创建规范

```go
// ✅ 正确：使用预定义错误和链式调用
err := ErrUserNotFound.
	WithRequestID(c.GetString("request_id")).
	WithUserID(targetUserID)

// ✅ 正确：使用模块特定错误创建函数
err := NewUserError(CodePasswordWrong, "password_incorrect").
	WithMeta(map[string]interface{}{
		"attempt_count": 3,
	})

// ❌ 错误：直接使用旧版错误
err := errors.New("user not found")
```

### 3. 错误传递规范

```go
// ✅ 正确：在Handler中使用responseService
func (h *Handler) HandlerExample(c *gin.Context) {
	if err := h.service.DoSomething(); err != nil {
		h.responseService.JSONError(c, err)
		return
	}
}

// ✅ 正确：在Service中直接返回错误
func (s *Service) ServiceExample() error {
	if something {
		return ErrUserNotFound.WithMeta(map[string]interface{}{
			"user_id": userID,
		})
	}
	return nil
}

// ✅ 正确：在中间件中使用c.Error
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := validateSomething(); err != nil {
			_ = c.Error(err) // 传递给错误处理中间件
			return
		}
		c.Next()
	}
}
```

### 4. 日志记录规范

```go
// ✅ 系统会自动记录结构化日志，包含：
// - 错误类型和严重程度
// - 请求ID和用户ID
// - 模块和操作信息
// - 业务元数据
// - 原始错误和堆栈信息

// ❌ 不需要手动记录错误日志
log.Error("user login failed") // 不推荐
```

### 5. 向后兼容规范

```go
// ✅ 正确：旧代码可以继续使用
func LegacyFunction() error {
	return errors.UserNotExists // 旧版AppError
}

// ✅ 正确：需要时转换为增强版
func NewFunction(c *gin.Context) {
	if err := LegacyFunction(); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			enhancedErr := FromLegacyError(appErr).
				WithRequestID(c.GetString("request_id"))

			h.responseService.JSONError(c, enhancedErr)
			return
		}
	}
}
```

### 6. HTML页面渲染规范

```go
// ✅ 正确：使用ViewService创建结构化页面数据
func (h *Handler) RenderPage(c *gin.Context) {
	pageData := h.viewService.PageData(c, "page.title")
	pageData.Data["content"] = content

	h.responseService.HTML(c, "template.html", pageData)
}

// ❌ 错误：直接使用gin.Context渲染
func (h *Handler) BadRender(c *gin.Context) {
	c.HTML(http.StatusOK, "template.html", data)
}
```

## 使用收益

通过使用新的错误处理和响应服务体系，您将获得：

1. **统一的API响应格式**：`{"code": 0, "message": "success", "data": {...}}`
2. **自动的错误日志记录**：结构化日志包含丰富的上下文信息
3. **完整的错误链路追踪**：通过请求ID追踪错误发生的完整流程
4. **智能的错误严重程度分级**：便于运维监控和告警
5. **良好的向后兼容性**：现有代码无需大规模修改
6. **类型安全的页面数据**：结构化的PageData提供更好的开发体验

这套体系显著提升了错误定位效率和开发体验，为后续的系统维护和扩展奠定了坚实基础。
