package view

import (
	"log/slog"

	"github.com/biuren/hebe/internal/shared/i18n"
	"github.com/biuren/hebe/internal/shared/types"
	"github.com/gin-gonic/gin"
)

// ContextBuilder 上下文构建器
type ContextBuilder struct {
	localizationManager i18n.I18nManager
}

// NewContextBuilder 创建上下文构建器
func NewContextBuilder(localizationManager i18n.I18nManager) *ContextBuilder {
	return &ContextBuilder{
		localizationManager: localizationManager,
	}
}

// BuildContext 构建渲染上下文
// 上下文获取认证信息是user 模块的 middleware 中的 SharedData() 中间件注入
// 逻辑是解密cookie获取客户端存储的AuthUser对象
// todo：客户端cookie 处理可能有安全风险
func (cb *ContextBuilder) BuildContext(c *gin.Context) *types.RenderContext {
	ctx := &types.RenderContext{
		IsLogin:  cb.getLoginStatus(c),
		AuthUser: cb.getAuthUser(c),
		Lang:     cb.localizationManager.GetLocaleFromContext(c),
	}

	return ctx
}

// getLoginStatus 提取登录状态
func (cb *ContextBuilder) getLoginStatus(c *gin.Context) bool {
	if isLogin, exists := c.Get("IsLogin"); exists {
		if login, ok := isLogin.(bool); ok {
			return login
		}
		slog.Warn("IsLogin 类型断言失败", "value", isLogin)
	}
	return false
}

// getAuthUser 提取认证用户信息
func (cb *ContextBuilder) getAuthUser(c *gin.Context) *types.AuthUser {
	if authUserRaw, exists := c.Get("AuthUser"); exists && authUserRaw != nil {
		if authUser, ok := authUserRaw.(*types.AuthUser); ok {
			return authUser
		}
		slog.Warn("AuthUser 类型断言失败", "value", authUserRaw)
	}
	return nil
}
