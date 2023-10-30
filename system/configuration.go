package system

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path"

	"github.com/GeertJohan/go.rice"
	"github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

//构造默认返回值对象
const (
	OPS_OK         int = 0  // 操作成功
	VALID_CODE_LEN int = 64 //确认码长度
	PWD_MIN_LEN    int = 6  //最小密码长度
	PWD_MAX_LEN    int = 50 //最大密码长度

	JOB_TYPE_EMAIL string = "email" //email任务
	JOB_TYPE_DB    string = "db"    // 数据库操作任务

	USER_STATUS_OK          int8 = 1 //用户正常状态
	USER_STATUS_WAIT_ACTIVE int8 = 2 //用户待激活
	USER_STATUS_DEL         int8 = 0 //用户已删除状态
	USER_NOT_VIP            int8 = 0 //用户不是vip
	USER_IS_VIP             int8 = 1 //用户是vip

	REG_EMAIL_VAILID_CODE_EXPIRE_DAY      int = 1 // 注册email过期天数
	FIND_PWD_EMAIL_VAILID_CODE_EXPIRE_DAY int = 1 // 重设密码邮件过期天数

	DEFAULT_USER_GROUP int8 = 1 // 用户注册默认分组

	VALID_CODE_ACTIVE_USER_EMAIL        int8 = 1 //激活用户
	VALID_CODE_RESET_PWD_EMAIL          int8 = 2 //重设密码
	VALID_CODE_INVITE_USER_JOIN_DAYBOOK int8 = 3 //邀请用户加入日记本
)

//定义支持的语言
type SupportLangs struct {
    Code     string `json:"code"`
    Label      string `json:"label"`
    LocalText string `json:"local-text"`
}

//Configs contains application configurations for all gin modes
type Configs struct {
	Debug   Config
	Release Config
	Test    Config
}

// 配置文件结构定义
type Config struct {
	Public   string `json:"public"`
	Domain   string `json:"domain"`
	SitePort string `json:"site_port"`
	SupportLangs []SupportLangs `json:"support_langs"`
	Database struct {
		Host     string	`json:"host"`
		Port     int `json:"port"`
		Name     string `json:"name"`
		User     string `json:"user"`
		Password string `json:"password"`
	}
	Cookie struct {
		Domain     string `json:"domain"`
		SecretKey  string `json:"secret_key"`
		ExpireDays int    `json:"expire_days"`
		HttpOnly   bool   `json:"http_only"`
	}
	Email struct {
		SMTPServer string `json:"smtp_server"`
		SMTPUser   string `json:"smtp_user"`
		SMTPPwd    string `json:"smtp_pwd"`
		SMTPPort   int    `json:"smtp_port"`
	}
	Redis struct{
		RedisServer string `json:"redis_server"`
		RedisPort int `json:"redis_port"`
		RedisPwd string `json:"redis_pwd"`
	}
	Upload struct{
		ImgSavePath string `json:"img_save_path"`
		ImgType string `json:"img_type"`
		ImgMaxSize int `json:"img_max_size"`
	}
}

//配置参数对象
type Params struct {
	Param []Param `xml:"param"`
}

type Param struct {
	Key string `xml:"key"`
	Val string `xml:"val"`
}

//current loaded config
var config *Config

//存储key-value对象
var params map[string]string

//加载文件
func LoadConfig() {
	box := rice.MustFindBox("../config")
	loadConfig(box.MustBytes("config.json"))
	loadParams(box.MustBytes("params.xml"))
	logrus.Info("config,params 加载成功！")
}

//设置日志
func SetLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stderr)
	if gin.Mode() == gin.DebugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

//加载key-value格式的参数
func loadParams(data []byte) {
	ps := &Params{}
	err := xml.Unmarshal(data, ps)
	if err != nil {
		panic(err)
	}
	params = make(map[string]string)
	for _, val := range ps.Param {
		params[val.Key] = val.Val
	}
}

//LoadConfig unmarshals config for current GIN_MODE
func loadConfig(data []byte) {
	configs := &Configs{}
	err := json.Unmarshal(data, configs)
	if err != nil {
		panic(err)
	}
	switch gin.Mode() {
	case gin.DebugMode:
		config = &configs.Debug
	case gin.ReleaseMode:
		config = &configs.Release
	case gin.TestMode:
		config = &configs.Test
	default:
		panic(fmt.Sprintf("Unknown gin mode %s", gin.Mode()))
	}
	if !path.IsAbs(config.Public) {
		workingDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		config.Public = path.Join(workingDir, config.Public)
	}
}

//GetConfig returns actual config
func GetConfig() *Config {
	return config
}

//获取key-val参数
func GetParams() map[string]string {
	return params
}

//PublicPath returns path to application public folder
func PublicPath() string {
	return config.Public
}

//UploadsPath returns path to public/uploads folder
func UploadsPath() string {
	return path.Join(config.Public, "uploads")
}
