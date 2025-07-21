package page

import (
	"fmt"

	"github.com/biuren/hebe/internal/app/render/localization"
	"github.com/biuren/hebe/internal/app/render/types"
	"github.com/biuren/hebe/internal/core/config"
	"github.com/gin-gonic/gin"
)

// PageDataBuilder 页面数据构建器接口
type PageDataBuilder interface {
	BuildPageData(title, description, keywords string, data map[string]interface{}, breadcrumb []types.BreadcrumbItem, context *types.RenderContext, localizationManager localization.ManagerInterface) interface{}
}

// TitleGenerator 标题生成器接口
type TitleGenerator interface {
	GenerateTitle(c *gin.Context, titleKey string) string
}

// DescriptionGenerator 描述生成器接口
type DescriptionGenerator interface {
	GenerateDescription(c *gin.Context, titleKey string) string
}

// BreadcrumbGenerator 面包屑生成器接口
type BreadcrumbGenerator interface {
	GenerateBreadcrumb(c *gin.Context) []types.BreadcrumbItem
}

// PageBuilder 页面构建器
type PageBuilder struct {
	config               *config.Config
	contextBuilder       *ContextBuilder
	localizationManager  localization.ManagerInterface
	pageDataBuilder      PageDataBuilder
	titleGenerator       TitleGenerator
	descriptionGenerator DescriptionGenerator
	breadcrumbGenerator  BreadcrumbGenerator
}

// NewPageBuilder 创建页面构建器
func NewPageBuilder(cfg *config.Config, contextBuilder *ContextBuilder, localizationManager localization.ManagerInterface, pageDataBuilder PageDataBuilder) *PageBuilder {
	pb := &PageBuilder{
		config:              cfg,
		contextBuilder:      contextBuilder,
		localizationManager: localizationManager,
		pageDataBuilder:     pageDataBuilder,
	}

	// 初始化各个生成器
	pb.titleGenerator = NewDefaultTitleGenerator(cfg, localizationManager)
	pb.descriptionGenerator = NewDefaultDescriptionGenerator(localizationManager)
	pb.breadcrumbGenerator = NewBreadcrumbStrategyManager()

	return pb
}

// BuildPageData 构建页面数据
func (pb *PageBuilder) BuildPageData(c *gin.Context, titleKey string) interface{} {
	return pb.pageDataBuilder.BuildPageData(
		pb.titleGenerator.GenerateTitle(c, titleKey),
		pb.descriptionGenerator.GenerateDescription(c, titleKey),
		"", // keywords
		make(map[string]interface{}),
		pb.breadcrumbGenerator.GenerateBreadcrumb(c),
		pb.contextBuilder.BuildContext(c),
		pb.localizationManager,
	)
}

// ===================== 标题生成器 =====================

// DefaultTitleGenerator 默认标题生成器
type DefaultTitleGenerator struct {
	config              *config.Config
	localizationManager localization.ManagerInterface
}

// NewDefaultTitleGenerator 创建默认标题生成器
func NewDefaultTitleGenerator(cfg *config.Config, localizationManager localization.ManagerInterface) TitleGenerator {
	return &DefaultTitleGenerator{
		config:              cfg,
		localizationManager: localizationManager,
	}
}

// GenerateTitle 生成页面标题
func (tg *DefaultTitleGenerator) GenerateTitle(c *gin.Context, titleKey string) string {
	locale := tg.localizationManager.GetLocaleFromContext(c)
	title := tg.localizationManager.GetMessage(titleKey, locale)
	siteName := tg.getSiteName()

	if title == titleKey {
		// 如果没有找到对应的标题，使用key本身
		return fmt.Sprintf("%s - %s", titleKey, siteName)
	}

	return fmt.Sprintf("%s - %s", title, siteName)
}

// getSiteName 获取网站名称
func (tg *DefaultTitleGenerator) getSiteName() string {
	if tg.config != nil {
		params := config.GetParams()
		if name, exists := params["site_name"]; exists {
			return name
		}
	}
	return "iDiary"
}

// ===================== 描述生成器 =====================

// DefaultDescriptionGenerator 默认描述生成器
type DefaultDescriptionGenerator struct {
	localizationManager localization.ManagerInterface
}

// NewDefaultDescriptionGenerator 创建默认描述生成器
func NewDefaultDescriptionGenerator(localizationManager localization.ManagerInterface) DescriptionGenerator {
	return &DefaultDescriptionGenerator{
		localizationManager: localizationManager,
	}
}

// GenerateDescription 生成页面描述
func (dg *DefaultDescriptionGenerator) GenerateDescription(c *gin.Context, titleKey string) string {
	locale := dg.localizationManager.GetLocaleFromContext(c)
	descKey := titleKey + ".desc"
	desc := dg.localizationManager.GetMessage(descKey, locale)

	// 如果没有找到对应的描述，使用通用默认描述
	if desc == descKey {
		return dg.localizationManager.GetMessage("common.default_page_desc", locale)
	}

	return desc
}

// ===================== 面包屑生成器 =====================
// 注意：DefaultBreadcrumbGenerator 已被删除，现在使用 BreadcrumbStrategyManager
// 具体实现在 breadcrumb_strategies.go 中
