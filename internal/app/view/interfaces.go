package view

import "github.com/gin-gonic/gin"

// ViewServiceInterface 视图服务接口
type ViewServiceInterface interface {
	// 页面数据创建
	PageData(c *gin.Context, titleKey string) *PageData

	// 页面标题生成
	PageTitle(content, author, pageType string) string

	// 国际化消息
	Message(key string, args ...interface{}) string
}
