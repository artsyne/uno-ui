package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/contracts"
	"github.com/biuren/hebe/pkg"
	"github.com/gin-gonic/gin"
)

/**
 * 认证助手模块
 *
 * 提供统一的用户认证、授权相关功能，包括：
 * - 用户登录状态检查
 * - 用户数据获取
 * - 认证中间件
 * - 重定向控制
 *
 * 注意：此模块依赖于 middlewares.SharedData() 中间件注入的认证信息
 * 使用前请确保在路由中正确配置了 SharedData 中间件
 *
 * @author hebe
 * @version 1.0.0
 */



// AuthHelper 认证助手结构体
// 封装了所有与用户认证相关的操作，提供类型安全和错误处理
type AuthHelper struct {
	config *config.Config
}

// NewAuthHelper 创建新的认证助手实例
// 虽然提供了全局实例，但也可以通过此方法创建独立实例
//
// Returns:
//   - *AuthHelper: 认证助手实例
func NewAuthHelper(config *config.Config) *AuthHelper {
	return &AuthHelper{
		config: config,
	}
}

// IsAuthenticated 检查用户是否已通过认证
// 这是核心的认证检查方法，会验证用户的登录状态和数据完整性
//
// 依赖条件：
//   - 必须先使用 middlewares.SharedData() 中间件注入认证信息
//   - 中间件会将 IsLogin 和 AuthUser 注入到上下文中
//
// Parameters:
//   - c: Gin上下文对象，包含请求和会话信息
//
// Returns:
//   - bool: 用户是否已认证
//   - *Login: 用户登录信息，如果未认证则返回nil
//
// Example:
//
//	// 路由配置（必需）
//	router.Use(middlewares.SharedData())
//
//	// 在处理函数中使用
//	isAuth, user := helpers.Auth.IsAuthenticated(c)
//	if isAuth {
//	    fmt.Printf("用户 %s 已登录", user.Username)
//	}
func (a *AuthHelper) IsAuthenticated(c *gin.Context) (bool, *types.AuthUser) {
	// 检查登录状态标志
	isLogin, exists := c.Get("IsLogin")
	if !exists || isLogin != true {
		return false, nil
	}

	// 获取用户会话数据
	userSession, exists := c.Get("AuthUser")
	if !exists {
		return false, nil
	}

	// 类型断言，确保数据类型正确
	userData, ok := userSession.(*types.AuthUser)
	if !ok {
		return false, nil
	}

	return true, userData
}

// RequireGuest 中间件：要求用户处于游客状态（未登录）
// 适用于登录页、注册页等不需要用户登录的页面
// 如果用户已登录，会自动重定向到指定URL
//
// Parameters:
//   - redirectURL: 可选的重定向URL，默认为 "/"
//
// Returns:
//   - gin.HandlerFunc: Gin中间件函数
//
// Example:
//
//	// 使用默认重定向URL "/"
//	router.Use(helpers.Auth.RequireGuest())
//
//	// 使用自定义重定向URL
//	router.Use(helpers.Auth.RequireGuest("/dashboard"))
func (a *AuthHelper) RequireGuest(redirectURL ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isAuth, _ := a.IsAuthenticated(c); isAuth {
			// 确定重定向URL
			redirect := "/"
			if len(redirectURL) > 0 && redirectURL[0] != "" {
				redirect = redirectURL[0]
			}

			c.Redirect(http.StatusFound, redirect)
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireAuth 中间件：要求用户已通过认证
// 适用于需要登录才能访问的页面，如个人中心、管理后台等
// 如果用户未登录，会重定向到登录页，并保存原始请求URL以便登录后跳回
//
// Parameters:
//   - loginURL: 可选的登录页URL，默认为 "/login"
//
// Returns:
//   - gin.HandlerFunc: Gin中间件函数
//
// Example:
//
//	// 使用默认登录页 "/login"
//	authGroup.Use(helpers.Auth.RequireAuth())
//
//	// 使用自定义登录页
//	adminGroup.Use(helpers.Auth.RequireAuth("/admin/login"))
func (a *AuthHelper) RequireAuth(loginURL ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isAuth, _ := a.IsAuthenticated(c); !isAuth {
			// 确定登录页URL
			redirect := "/login"
			if len(loginURL) > 0 && loginURL[0] != "" {
				redirect = loginURL[0]
			}

			// 保存原始请求URL，用于登录后重定向
			originalURL := c.Request.URL.String()
			redirectWithReturn := fmt.Sprintf("%s?redirect=%s", redirect, url.QueryEscape(originalURL))

			c.Redirect(http.StatusFound, redirectWithReturn)
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetCurrentUser 安全地获取当前登录用户信息
// 推荐在大多数场景下使用此方法，因为它提供了明确的状态返回
//
// Parameters:
//   - c: Gin上下文对象
//
// Returns:
//   - *Login: 用户登录信息，如果未认证则返回nil
//   - bool: 是否成功获取用户信息
//
// Example:
//
//	user, ok := helpers.Auth.GetCurrentUser(c)
//	if !ok {
//	    c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
//	    return
//	}
//	// 使用用户信息
//	fmt.Printf("当前用户：%s", user.Username)
func (a *AuthHelper) GetCurrentUser(c *gin.Context) (*types.AuthUser, bool) {
	isAuth, userData := a.IsAuthenticated(c)
	return userData, isAuth
}

// MustGetCurrentUser 强制获取当前登录用户信息
// 适用于已经通过认证中间件保护的路由
// 如果用户未认证会引发panic，应谨慎使用
//
// Parameters:
//   - c: Gin上下文对象
//
// Returns:
//   - *Login: 用户登录信息
//
// Panics:
//   - 当用户未认证时会panic
//
// Example:
//
//	// 在已认证的路由中使用
//	user := helpers.Auth.MustGetCurrentUser(c)
//	fmt.Printf("当前用户：%s", user.Username)
func (a *AuthHelper) MustGetCurrentUser(c *gin.Context) *types.AuthUser {
	userData, ok := a.GetCurrentUser(c)
	if !ok {
		panic("用户未认证：在未认证的上下文中调用了 MustGetCurrentUser")
	}
	return userData
}

// IsLogin 检查用户是否已登录（简化版本）
// 这是 IsAuthenticated 的简化版本，只返回布尔值
// 适用于只需要判断登录状态而不需要用户数据的场景
//
// Parameters:
//   - c: Gin上下文对象
//
// Returns:
//   - bool: 用户是否已登录
//
// Example:
//
//	if helpers.Auth.IsLogin(c) {
//	    // 用户已登录的逻辑
//	}
func (a *AuthHelper) IsLogin(c *gin.Context) bool {
	isAuth, _ := a.IsAuthenticated(c)
	return isAuth
}

// IfLoginGotoIndex 如果已登录则跳转到首页
// 适用于登录页、注册页等，避免已登录用户重复访问
//
// Parameters:
//   - c: Gin上下文对象
//
// Example:
//
//	// 在登录页处理函数开头使用
//	helpers.Auth.IfLoginGotoIndex(c)
func (a *AuthHelper) IfLoginGotoIndex(c *gin.Context) {
	if a.IsLogin(c) {
		c.Redirect(http.StatusFound, "/")
		c.Abort()
	}
}

// GetLogin 获取登录用户信息（兼容性方法）
// 这是 GetCurrentUser 的兼容性版本，只返回用户数据
// 如果用户未登录返回 nil
//
// Parameters:
//   - c: Gin上下文对象
//
// Returns:
//   - *Login: 用户登录信息，如果未认证则返回nil
//
// Example:
//
//	user := helpers.Auth.GetLogin(c)
//	if user != nil {
//	    fmt.Printf("当前用户：%s", user.Username)
//	}
func (a *AuthHelper) GetLogin(c *gin.Context) *types.AuthUser {
	userData, _ := a.GetCurrentUser(c)
	return userData
}

// SetLoginCookie 设置登录Cookie
func (a *AuthHelper) SetLoginCookie(c *gin.Context, authUser *types.AuthUser) {
	// 序列化用户数据
	userData, _ := json.Marshal(authUser)

	// 加密用户数据
	encryptedData, err := pkg.AESEncrypt(string(userData), a.config.Cookie.Secret)
	if err != nil {
		return
	}

	// 设置登录cookie
	pkg.WriteCookie(c.Writer, "ud", encryptedData)
}

// ParseUserFromCookie 
func (a *AuthHelper) ParseUserFromCookie(c *gin.Context) *types.AuthUser {
	cookieEnUserData := pkg.ReadCookie(c.Request, "ud")
	if len(cookieEnUserData) < 1 {
		return nil
	}

	byteUserData, err := pkg.AESDecrypt(cookieEnUserData, a.config.Cookie.Secret)
	if err != nil || byteUserData == nil {
		return nil
	}

	var userData types.AuthUser
	err = json.Unmarshal(byteUserData, &userData)
	if err != nil {
		return nil
	}

	return &userData
}