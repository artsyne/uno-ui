package render

import "github.com/gin-gonic/gin"

// ContextBuilder 上下文构建器
type ContextBuilder struct {
	localizationManager *LocalizationManager
}

// NewContextBuilder 创建上下文构建器
func NewContextBuilder(localizationManager *LocalizationManager) *ContextBuilder {
	return &ContextBuilder{
		localizationManager: localizationManager,
	}
}

// BuildContext 构建渲染上下文
// 上下文获取认证信息是user 模块的 middleware 中的 SharedData() 中间件注入
// 逻辑是解密cookie获取客户端存储的AuthUser对象
// todo：客户端cookie 处理可能有安全风险
func (cb *ContextBuilder) BuildContext(c *gin.Context) *RenderContext {
	ctx := &RenderContext{}

	// 获取登录状态
	if isLogin, exists := c.Get("IsLogin"); exists {
		if login, ok := isLogin.(bool); ok {
			ctx.IsLogin = login
		}
	}

	// 获取认证用户并转换为接口类型
	if authUserRaw, exists := c.Get("AuthUser"); exists && authUserRaw != nil {
		if renderUser, ok := authUserRaw.(RenderUser); ok {
			ctx.AuthUser = renderUser
		}
	}

	// 智能检测当前语言
	ctx.Lang = cb.localizationManager.GetLocaleFromContext(c)

	return ctx
}
