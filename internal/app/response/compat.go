package response

import (
	"github.com/biuren/hebe/internal/core/config"
	"github.com/gin-gonic/gin"
)

// LegacyResponseH 兼容旧版本的ResponseH函数
// 用于从pkg/template.go平滑迁移
// TODO: 在所有模块迁移完成后移除此函数
func LegacyResponseH(c *gin.Context) gin.H {
	return gin.H{
		"Status": config.OPS_SUCCESS,
		"Data":   "",
	}
}

// LegacyDefaultH 兼容旧版本的DefaultH函数
// 用于从pkg/template.go平滑迁移
// TODO: 在所有模块迁移完成后移除此函数
func LegacyDefaultH(c *gin.Context) gin.H {
	userData, _ := c.Get("UserData")
	isLogin, _ := c.Get("IsLogin")
	return gin.H{
		"UserData": userData,
		"Title":    "",
		"IsLogin":  isLogin,
	}
}
