package middlewares

import (
	"net/http"
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/biuren/hebe/internal/errors"
)

// 捕获error 错误
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		 // 先调用c.Next()执行后面的中间件，最后再处理错误信息
		c.Next()

		// 所有中间件及router处理完毕后从这里开始执行
		for _, e := range c.Errors {
			err := e.Err
			// 若是自定义的错误则将code、msg返回
			if e, ok := err.(*errors.AppError); ok {
				c.JSON(http.StatusOK, gin.H{
					"code": e.Code,
					"msgkey":  e.MsgKey,
					"msg":e.MsgKey+",todo",
					"data": err.Error(),
				})
			} else {
				// 若非自定义错误则返回详细错误信息err.Error()
				c.JSON(http.StatusOK, gin.H{
					"code": errors.InternalErr.Code,
					"msgkey":  errors.InternalErr.MsgKey,
					"msg":errors.InternalErr.MsgKey+",todo",
					"data": err.Error(),
				})
			}
			return // 检查一个错误就退出
		}
	}
}

// 捕获 panic 错误
func PanicHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // 记录错误日志
                logrus.WithFields(logrus.Fields{
                    "error": err,
                    "stack": string(debug.Stack()),
                }).Error("A panic occurred")

                // 返回错误响应,todo: 后面需要跳转到错误页面
                c.JSON(500, gin.H{"message": "Internal server error"})
            }
        }()
        c.Next()
    }
}