package settings

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

//构造默认返回值对象
const (
	OPS_SUCCESS         int = 0  // 操作成功
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

	//日记本阅读权限
	DAYBOOK_VIEW_PUBLIC  int8 = 1 //公开
	DAYBOOK_VIEW_PRIVATE int8 = 2 //私有

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
	SummaryMax int `json:"summary_max"`
	Suffix string `json:"suffix"`
	SupportLangs []SupportLangs `json:"support_langs"`
	PageSize int `json:"page_size"`
	MaxDaybookLimit int `json:"max_daybook_limit"`
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

// 全局配置对象
var config *Config

// 全局key-value对象
var params map[string]string

// 全局css对象，存储 tailwind css类名
var css map[string]string

//加载文件
func LoadConfig() {
    // 获取当前工作目录
    dir, _ := os.Getwd()

    // 构建文件路径
    configPath := filepath.Join(dir, "./settings/config.json")
    paramsPath := filepath.Join(dir, "./settings/params.xml")
	cssPath := filepath.Join(dir, "./settings/css.json")

    // 读取 config.json 文件
    configBytes, _ := os.ReadFile(configPath)

    // 读取 params.xml 文件
    paramsBytes, _ := os.ReadFile(paramsPath)

	// 读取 css.json 文件
	cssBytes, _ := os.ReadFile(cssPath)

    // 加载配置和参数
    loadConfig(configBytes)
    loadParams(paramsBytes)
	loadCSS(cssBytes)
}

//GetConfig returns actual config
func GetConfig() *Config {
	return config
}

//获取key-val参数
func GetParams() map[string]string {
	return params
}

//获取css代码
func GetCSS() map[string]string {
	return css
}

//PublicPath returns path to application public folder
func PublicPath() string {
	return config.Public
}

//UploadsPath returns path to public/uploads folder
func UploadsPath() string {
	return path.Join(config.Public, "uploads")
}

//加载css文件
func loadCSS(data []byte) {
	err := json.Unmarshal(data, &css)
	if err != nil {
		panic(err)
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