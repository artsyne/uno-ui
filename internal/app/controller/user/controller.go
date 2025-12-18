package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/biuren/hebe/internal/app/view"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/gin-gonic/gin"
)

// UserController 用户HTTP处理层实现
// 负责处理所有与用户相关的HTTP请求，包括用户认证、注册、个人信息管理等
//
// 字段:
//   - config: 应用配置信息
//   - service: 用户业务逻辑服务
//   - auth: 用户认证辅助工具
//   - renderService: 视图渲染服务
type UserController struct {
	config        *config.Config
	service       *user.Service
	auth          *user.AuthHelper
	renderService *view.ViewService
}

// NewUserController 创建用户HTTP处理层实例
// 通过依赖注入的方式初始化用户控制器，确保所有必要的服务都已正确配置
//
// 参数:
//   - config: 应用配置信息
//   - service: 用户业务逻辑服务
//   - renderService: 视图渲染服务
//   - auth: 用户认证辅助工具
//
// 返回:
//   - *UserController: 初始化完成的用户控制器实例
func NewUserController(config *config.Config, service *user.Service, renderService *view.ViewService, auth *user.AuthHelper) *UserController {
	return &UserController{
		config:        config,
		service:       service,
		renderService: renderService,
		auth:          auth,
	}
}

// isHtmxRequest 检查是否为HTMX请求
func (ctrl *UserController) isHtmxRequest(c *gin.Context) bool {
	return c.GetHeader("HX-Request") == "true"
}

// htmxError 返回HTMX错误响应（内联样式）
func (ctrl *UserController) htmxError(c *gin.Context, code int, message string) {
	html := fmt.Sprintf(`<span class="text-red-500 text-sm">%s</span>`, message)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(code, html)
}

// htmxAlert 返回HTMX错误警告框（块级样式）
func (ctrl *UserController) htmxAlert(c *gin.Context, code int, message string) {
	html := fmt.Sprintf(`<div class="p-4 mb-4 text-sm text-red-800 rounded-lg bg-red-50 dark:bg-gray-800 dark:text-red-400" role="alert"><span class="font-medium">Error!</span> %s</div>`, message)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(code, html)
}

// htmxSuccess 返回HTMX成功响应
func (ctrl *UserController) htmxSuccess(c *gin.Context, message string) {
	html := fmt.Sprintf(`<span class="text-green-500 text-sm">%s</span>`, message)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// htmxToast 返回HTMX Toast触发响应
func (ctrl *UserController) htmxToast(c *gin.Context, msgType, message string) {
	payload := map[string]map[string]string{
		"showToast": {
			"type":    msgType,
			"message": message,
		},
	}
	jsonBytes, _ := json.Marshal(payload)
	c.Header("HX-Trigger", string(jsonBytes))
	c.Status(http.StatusOK)
}
