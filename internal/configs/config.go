package configs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

// 构造默认返回值对象
const (
	OPS_SUCCESS    int = 0  // 操作成功
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

// ConfigSet contains application configurations for all gin modes
type ConfigSet struct {
	Debug   Config
	Release Config
	Test    Config
}

// 配置文件结构定义
type Config struct {
	Domain          string `yaml:"domain"`
	SitePort        string `yaml:"site_port"`
	SummaryMax      int    `yaml:"summary_max"`
	Suffix          string `yaml:"suffix"`
	PageSize        int    `yaml:"page_size"`
	MaxDaybookLimit int    `yaml:"max_daybook_limit"`
	HashIDSalt      string `yaml:"hash_id_salt"`
	Database        struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	}
	Cookie struct {
		Domain     string `yaml:"domain"`
		SecretKey  string `yaml:"secret_key"`
		ExpireDays int    `yaml:"expire_days"`
		HttpOnly   bool   `yaml:"http_only"`
	}
	Email struct {
		SMTPServer string `yaml:"smtp_server"`
		SMTPUser   string `yaml:"smtp_user"`
		SMTPPwd    string `yaml:"smtp_pwd"`
		SMTPPort   int    `yaml:"smtp_port"`
	}
	Redis struct {
		RedisServer string `yaml:"redis_server"`
		RedisPort   int    `yaml:"redis_port"`
		RedisPwd    string `yaml:"redis_pwd"`
	}
	Upload struct {
		ImgSavePath string `yaml:"img_save_path"`
		ImgType     string `yaml:"img_type"`
		ImgMaxSize  int    `yaml:"img_max_size"`
	}
}

type Param struct {
	Key string `yaml:"key"`
	Val string `yaml:"val"`
}

// 定义支持的语言
type SupportLang struct {
	Code      string `yaml:"code"`
	Label     string `yaml:"label"`
	LocalText string `yaml:"local-text"`
}

// 配置对象
var config *Config

// key-value对象
var params map[string]string

// 时区对象
var timeZone []string

// 支持语言
var supportLangs []SupportLang

// 配置文件目录
var configDir string
var viewsDir string
var saticDir string

// 加载文件
func LoadConfig() {
    // 获取环境变量
    configDir = GetEnv("CONFIG_PATH")
    viewsDir = GetEnv("VIEW_PATH")
    saticDir = GetEnv("STATIC_PATH")

	// 构建文件路径
	configPath := filepath.Join(configDir, "config.yaml")
	paramsPath := filepath.Join(configDir, "params.yaml")
	timeZonePath := filepath.Join(configDir, "timezone.yaml")
	supportLangsPath := filepath.Join(configDir, "support_langs.yaml")

	// 读取 config.yaml 文件
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	// 读取 timezone.yaml 文件
	timeZoneBytes, err := os.ReadFile(timeZonePath)
	if err != nil {
		panic(err)
	}

	// 读取 params.yaml 文件
	paramsBytes, err := os.ReadFile(paramsPath)
	if err != nil {
		panic(err)
	}

	// 读取 support_langs.yaml 文件
	supportLangsBytes, err := os.ReadFile(supportLangsPath)
	if err != nil {
		panic(err)
	}

	// 加载配置和参数
	loadConfig(configBytes)
	loadParams(paramsBytes)
	loadTimeZone(timeZoneBytes)
	loadSupportLangs(supportLangsBytes)
}

// GetConfig returns actual config
func GetConfig() *Config {
	return config
}

// 获取key-val参数
func GetParams() map[string]string {
	return params
}

// 获取时区
func GetTimeZone() []string {
	return timeZone
}

// 获取支持的语言
func GetSupportLangs() []SupportLang {
	return supportLangs
}

// 获取 views 目录
func GetViewsPath() string {
    return viewsDir
}

// GetStaticPath returns path to application public folder
func GetStaticPath() string {
	return saticDir
}

// UploadsPath returns path to public/uploads folder
func UploadsPath() string {
	return filepath.Join(saticDir, "uploads")
}

// 读环境变量，如果为空则 panic
func GetEnv(key string) string {
    val := os.Getenv(key)
    if key == "" {
        panic(key + " environment variable is not set")
    }
    return val
}

// #region 配置加载处理
// 加载key-value格式的参数
func loadParams(data []byte) {
	ps := []Param{}
	err := yaml.Unmarshal(data, &ps)
	if err != nil {
		panic(err)
	}
	params = make(map[string]string)
	for _, val := range ps {
		params[val.Key] = val.Val
	}
}

// 加载时区 timezone.yaml
func loadTimeZone(data []byte) {
	tz := &timeZone
	err := yaml.Unmarshal(data, tz)
	if err != nil {
		panic(err)
	}
}

// 加载支持的语言 support_langs.yaml
func loadSupportLangs(data []byte) {
	supportLangs = make([]SupportLang, 0)
	err := yaml.Unmarshal(data, &supportLangs)
	if err != nil {
		panic(err)
	}
}

// LoadConfig unmarshals config for current GIN_MODE
func loadConfig(data []byte) {
	configs := &ConfigSet{}
	err := yaml.Unmarshal(data, configs)
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
}

// #endregion
