package page

import (
	"fmt"
	"strings"

	"github.com/biuren/hebe/internal/app/render/types"
	"github.com/gin-gonic/gin"
)

// BreadcrumbStrategy 面包屑策略接口
type BreadcrumbStrategy interface {
	CanHandle(path string) bool
	Generate(c *gin.Context, path string) []types.BreadcrumbItem
}

// BaseBreadcrumbStrategy 基础面包屑策略
type BaseBreadcrumbStrategy struct{}

// GetHomeBreadcrumb 获取首页面包屑
func (b *BaseBreadcrumbStrategy) GetHomeBreadcrumb() []types.BreadcrumbItem {
	return []types.BreadcrumbItem{
		{Name: "首页", URL: "/"},
	}
}

// ExtractUserInfo 提取用户信息
func (b *BaseBreadcrumbStrategy) ExtractUserInfo(c *gin.Context) (string, string, bool) {
	if authUser, exists := c.Get("AuthUser"); exists {
		if user, ok := authUser.(map[string]interface{}); ok {
			if userName, exists := user["user_name"]; exists {
				if userID, exists := user["user_id"]; exists {
					return fmt.Sprintf("%v", userName), fmt.Sprintf("%v", userID), true
				}
			}
		}
	}
	return "", "", false
}

// ===================== 用户页面策略 =====================

// UserBreadcrumbStrategy 用户页面面包屑策略
type UserBreadcrumbStrategy struct {
	BaseBreadcrumbStrategy
}

// NewUserBreadcrumbStrategy 创建用户页面面包屑策略
func NewUserBreadcrumbStrategy() BreadcrumbStrategy {
	return &UserBreadcrumbStrategy{}
}

// CanHandle 检查是否可以处理该路径
func (u *UserBreadcrumbStrategy) CanHandle(path string) bool {
	return strings.HasPrefix(path, "/u/")
}

// Generate 生成用户页面面包屑
func (u *UserBreadcrumbStrategy) Generate(c *gin.Context, path string) []types.BreadcrumbItem {
	breadcrumb := u.GetHomeBreadcrumb()

	// 添加用户信息
	if userName, userID, exists := u.ExtractUserInfo(c); exists {
		breadcrumb = append(breadcrumb, types.BreadcrumbItem{
			Name: userName,
			URL:  fmt.Sprintf("/u/%s", userID),
		})
	}

	// 添加具体页面
	if strings.Contains(path, "/settings") {
		breadcrumb = append(breadcrumb, types.BreadcrumbItem{
			Name: "设置",
		})
	}

	return breadcrumb
}

// ===================== 日记页面策略 =====================

// DiaryBreadcrumbStrategy 日记页面面包屑策略
type DiaryBreadcrumbStrategy struct {
	BaseBreadcrumbStrategy
}

// NewDiaryBreadcrumbStrategy 创建日记页面面包屑策略
func NewDiaryBreadcrumbStrategy() BreadcrumbStrategy {
	return &DiaryBreadcrumbStrategy{}
}

// CanHandle 检查是否可以处理该路径
func (d *DiaryBreadcrumbStrategy) CanHandle(path string) bool {
	return strings.HasPrefix(path, "/d/")
}

// Generate 生成日记页面面包屑
func (d *DiaryBreadcrumbStrategy) Generate(c *gin.Context, path string) []types.BreadcrumbItem {
	breadcrumb := d.GetHomeBreadcrumb()

	// 添加用户信息
	if userName, userID, exists := d.ExtractUserInfo(c); exists {
		breadcrumb = append(breadcrumb, types.BreadcrumbItem{
			Name: userName,
			URL:  fmt.Sprintf("/u/%s", userID),
		})
	}

	// TODO: 从路径或数据中提取日记本信息
	// 例如: /d/123 -> 日记本名称
	// 这里需要从数据库或缓存中获取日记本信息

	return breadcrumb
}

// ===================== 文章页面策略 =====================

// PostBreadcrumbStrategy 文章页面面包屑策略
type PostBreadcrumbStrategy struct {
	BaseBreadcrumbStrategy
}

// NewPostBreadcrumbStrategy 创建文章页面面包屑策略
func NewPostBreadcrumbStrategy() BreadcrumbStrategy {
	return &PostBreadcrumbStrategy{}
}

// CanHandle 检查是否可以处理该路径
func (p *PostBreadcrumbStrategy) CanHandle(path string) bool {
	return strings.HasPrefix(path, "/p/")
}

// Generate 生成文章页面面包屑
func (p *PostBreadcrumbStrategy) Generate(c *gin.Context, path string) []types.BreadcrumbItem {
	breadcrumb := p.GetHomeBreadcrumb()

	// 添加用户信息
	if userName, userID, exists := p.ExtractUserInfo(c); exists {
		breadcrumb = append(breadcrumb, types.BreadcrumbItem{
			Name: userName,
			URL:  fmt.Sprintf("/u/%s", userID),
		})
	}

	// TODO: 从路径或数据中提取日记本和文章信息
	// 例如: /p/123 -> 日记本名称 > 文章标题
	// 这里需要从数据库或缓存中获取相关信息

	return breadcrumb
}

// ===================== 策略管理器 =====================

// BreadcrumbStrategyManager 面包屑策略管理器
type BreadcrumbStrategyManager struct {
	strategies []BreadcrumbStrategy
}

// NewBreadcrumbStrategyManager 创建面包屑策略管理器
func NewBreadcrumbStrategyManager() *BreadcrumbStrategyManager {
	return &BreadcrumbStrategyManager{
		strategies: []BreadcrumbStrategy{
			NewUserBreadcrumbStrategy(),
			NewDiaryBreadcrumbStrategy(),
			NewPostBreadcrumbStrategy(),
		},
	}
}

// AddStrategy 添加策略
func (bm *BreadcrumbStrategyManager) AddStrategy(strategy BreadcrumbStrategy) {
	bm.strategies = append(bm.strategies, strategy)
}

// GenerateBreadcrumb 生成面包屑
func (bm *BreadcrumbStrategyManager) GenerateBreadcrumb(c *gin.Context) []types.BreadcrumbItem {
	path := c.Request.URL.Path

	// 查找匹配的策略
	for _, strategy := range bm.strategies {
		if strategy.CanHandle(path) {
			return strategy.Generate(c, path)
		}
	}

	// 默认返回首页面包屑
	return []types.BreadcrumbItem{
		{Name: "首页", URL: "/"},
	}
}
