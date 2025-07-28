package view

import (
	"github.com/gin-gonic/gin"

	"github.com/biuren/hebe/internal/shared/http"
)

// responseService 响应服务实现
type responseService struct {
	handler *response.Handler
}

// NewResponseService 创建响应服务
func NewResponseService(handler *response.Handler) ResponseService {
	return &responseService{
		handler: handler,
	}
}

// JSON 返回JSON响应
func (s *responseService) JSON(c *gin.Context, code int, data interface{}) {
	s.handler.JSON(c, code, data)
}

// JSONSuccess 返回成功JSON响应
func (s *responseService) JSONSuccess(c *gin.Context, data interface{}) {
	s.handler.JSONSuccess(c, data)
}

// JSONError 返回错误JSON响应
func (s *responseService) JSONError(c *gin.Context, err error, data ...interface{}) {
	s.handler.JSONError(c, err, data...)
}

// HTML 渲染HTML模板
func (s *responseService) HTML(c *gin.Context, name string, data interface{}) {
	s.handler.HTML(c, name, data)
}

// HTMLError 渲染错误页面
func (s *responseService) HTMLError(c *gin.Context, title, message string) {
	s.handler.HTMLError(c, title, message)
}

// Redirect 重定向响应
func (s *responseService) Redirect(c *gin.Context, url string) {
	s.handler.Redirect(c, url)
}