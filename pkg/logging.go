package pkg

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 函数为应用程序配置日志记录器。它将日志格式设置为 JSON，
// 根据当前的 gin 模式（Debug 或 Release）调整日志级别，
// 并设置一个带有日志轮转的日志文件（最大大小 20MB，最多 10 个备份，最大年龄 28 天，带有压缩）。
// 日志输出被写入标准错误和日志文件
func SetLogger() {
	// json 格式化输出日志
	logrus.SetFormatter(&logrus.JSONFormatter{})
	
	if gin.Mode() == gin.DebugMode {
		logrus.SetLevel(logrus.DebugLevel)
	} else if gin.Mode() == gin.ReleaseMode {
        logrus.SetLevel(logrus.ErrorLevel)
    }

	logFile := &lumberjack.Logger{
        Filename:   "logs/app.log", // 日志文件的位置
        MaxSize:    20,             // megabytes
        MaxBackups: 10,
        MaxAge:     28,   //days
        Compress:   true, // disabled by default
    }

    logrus.SetOutput(io.MultiWriter(os.Stderr, logFile))
}