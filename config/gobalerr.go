package gobalerr

import (
	"github.com/gin-gonic/gin"

	"net/http"
)

//错误对象
type ErrCode struct {
	//错误编号
	Num int
	//错误描述
	ErrDesc string
}

//全局错误码定义
//1xxxxxx输入参数错误，2XXXXX 业务逻辑错误，3xxxxxx数据库类错误，4xxxxx外部依赖错误，5xxxxx定时任务类错误
const (
	INPUT_PARAM_ERR            = 100000
	REGISTER_USER_NAME_EXISTS  = 100001
	REGISTER_USER_EMAIL_EXISTS = 100002
	PARSE_INPUT_PARAM_ERR      = 100003
	UPLOAD_FILE_TYPE_ERR       = 100004
	UPLOAD_FILE_TOO_LARGE      = 100005
	UPLOAD_FILE_ERR            = 100006

	BIZ_ERR         = 200000
	USER_NOT_EXISTS = 200001
	USER_PWD_ERR    = 200002
	PARSE_JSON_ERR  = 200003

	DB_OPERATION_ERR = 300000
)

var EnumErrDesc = map[int]string{
	INPUT_PARAM_ERR:            "输入参数错误",
	REGISTER_USER_NAME_EXISTS:  "注册用户名已存在",
	REGISTER_USER_EMAIL_EXISTS: "注册Email已存在",
	PARSE_INPUT_PARAM_ERR:      "解析输入参数错误",
	UPLOAD_FILE_TYPE_ERR:       "不支持上传文件类型",
	UPLOAD_FILE_TOO_LARGE:      "上传文件太大",
	UPLOAD_FILE_ERR:            "上传文件错误",

	BIZ_ERR:         "业务逻辑错误",
	USER_NOT_EXISTS: "用户不存在",
	USER_PWD_ERR:    "登录密码不正确",
	PARSE_JSON_ERR:  "JSON字符串解析错误",

	DB_OPERATION_ERR: "数据库操作错误",

	400000: "外部依赖错误",

	500000: "数据库操作错误",
}

//显示错误Json并返回
func ResponseErr(c *gin.Context, errCode int) {
	c.JSON(http.StatusBadRequest, gin.H{
		"Status":  errCode,
		"ErrDesc": EnumErrDesc[errCode],
	})
}

func ResponseErrWithDesc(c *gin.Context, errCode int, errDesc string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"Status":  errCode,
		"ErrDesc": errDesc,
	})
}
