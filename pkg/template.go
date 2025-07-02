package pkg

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// TODO：文件需要重构

// IsActive checks uri against currently active (uri, or nil) and returns "active" if they are equal
func IsActive(active interface{}, uri string) string {
	if s, ok := active.(string); ok {
		if s == uri {
			return "active"
		}
	}
	return ""
}

// DateTime prints timestamp in human format
func DateTime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

// DefaultH returns common to all pages template data
func DefaultH(c *gin.Context) gin.H {
	userData, _ := c.Get("UserData")
	isLogin, _ := c.Get("IsLogin")
	return gin.H{
		"UserData": userData,
		"Title":    "",
		"IsLogin":  isLogin,
	}
}

// ResponseH ajax 操作返回JSON串
// DEPRECATED: 此函数已废弃，请使用 internal/app/response.ResponseService
// 保留此函数仅为向后兼容，将在后续版本中移除
func ResponseH(c *gin.Context) gin.H {
	return gin.H{
		"Status": 0, // 使用硬编码的成功码，避免依赖internal包
		"Data":   "",
	}
}
