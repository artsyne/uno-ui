package response

import "github.com/gin-gonic/gin"

// ResponseServiceInterface 统一响应服务接口
type ResponseServiceInterface interface {
	// JSON响应
	JSON(c *gin.Context, code int, data interface{})
	JSONSuccess(c *gin.Context, data interface{})
	JSONError(c *gin.Context, err error, data ...interface{})

	// HTML响应
	HTML(c *gin.Context, name string, data interface{})
	HTMLError(c *gin.Context, title, message string)

	// 其他响应
	Redirect(c *gin.Context, url string)
}
