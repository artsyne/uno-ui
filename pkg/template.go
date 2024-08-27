package pkg

import (
	"fmt"
	"time"

	//	"../models"
	"github.com/biuren/hebe/internal/configs"
	"github.com/gin-gonic/gin"
)

// TODO：文件需要重构

//IsActive checks uri against currently active (uri, or nil) and returns "active" if they are equal
func IsActive(active interface{}, uri string) string {
	if s, ok := active.(string); ok {
		if s == uri {
			return "active"
		}
	}
	return ""
}

//DateTime prints timestamp in human format
func DateTime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

//DefaultH returns common to all pages template data
func DefaultH(c *gin.Context) gin.H {
	userData, _ := c.Get("UserData")
	isLogin, _ := c.Get("IsLogin")
	return gin.H{
		"UserData": userData,
		"Title":    "",
		"IsLogin":  isLogin,
	}
}

// ajax 操作返回JSON串
func ResponseH(c *gin.Context) gin.H {
	return gin.H{
		"Status": configs.OPS_SUCCESS,
		"Data":   "",
	}
}
