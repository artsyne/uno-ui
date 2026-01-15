package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

// 全局配置实例
var (
	globalConfig *Config
	v            *viper.Viper
	staticDir    string
)

// InitConfig 初始化配置
func InitConfig() error {
	// 加载 .env.local 文件，正式发布是删除，环境变量会在外部注入
	if err := gotenv.Load(".env.local"); err != nil {
		// .env.local 文件不存在时不报错，继续使用环境变量
		fmt.Printf("Warning: .env.local file not found: %v\n", err)
	}

	v = viper.New()

	// 环境变量支持
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("IDIARY")

	// 解析JSON格式的环境变量
	var supportLangs []SupportLang
	supportLangsStr := v.GetString("SUPPORT_LANGS")
	if supportLangsStr == "" {
		return fmt.Errorf("SUPPORT_LANGS配置不能为空")
	}
	if err := parseJSON(supportLangsStr, &supportLangs); err != nil {
		return fmt.Errorf("解析SUPPORT_LANGS失败: %w", err)
	}

	var timezones []string
	timezonesStr := v.GetString("TIMEZONES")
	if timezonesStr == "" {
		return fmt.Errorf("TIMEZONES配置不能为空")
	}
	if err := parseJSON(timezonesStr, &timezones); err != nil {
		return fmt.Errorf("解析TIMEZONES失败: %w", err)
	}

	var params []Param
	paramsStr := v.GetString("PARAMS")
	if paramsStr == "" {
		return fmt.Errorf("PARAMS配置不能为空")
	}
	if err := parseJSON(paramsStr, &params); err != nil {
		return fmt.Errorf("解析PARAMS失败: %w", err)
	}

	// 初始化配置结构体
	globalConfig = &Config{}

	// 验证并设置基本配置
	if env := v.GetString("ENV"); env == "" {
		return fmt.Errorf("ENV配置不能为空")
	} else {
		globalConfig.Env = env
	}

	globalConfig.App = AppConfig{
		GinMode:              requireString(v, "GIN_MODE"),
		Domain:               requireString(v, "DOMAIN"),
		Port:                 requireString(v, "PORT"),
		SummaryMax:           requireInt(v, "SUMMARY_MAX"),
		Suffix:               requireString(v, "SUFFIX"),
		PageSize:             requireInt(v, "PAGE_SIZE"),
		MaxDaybookLimit:      requireInt(v, "MAX_DAYBOOK_LIMIT"),
		MaxPostLength:        requireInt(v, "MAX_POST_LENGTH"),
		MaxCommentLength:     requireInt(v, "MAX_COMMENT_LENGTH"),
		MaxDescriptionLength: requireInt(v, "MAX_DESCRIPTION_LENGTH"),
		StaticDir:            requireString(v, "STATIC_DIR"),
		LocalesDir:           requireString(v, "LOCALES_DIR"),
		DefaultLang:          requireString(v, "DEFAULT_LANG"),
	}

	// 用户配置
	globalConfig.User = UserConfig{
		DefaultGroup: int8(requireInt(v, "USER_DEFAULT_GROUP")),
		PwdMaxLen:    requireInt(v, "PWD_MAX_LEN"),
		PwdMinLen:    requireInt(v, "PWD_MIN_LEN"),
		StatusDeleted:    int8(requireInt(v, "USER_STATUS_DEL")),
		StatusOK:         int8(requireInt(v, "USER_STATUS_OK")),
		StatusWaitActive: int8(requireInt(v, "USER_STATUS_WAIT_ACTIVE")),
	}

	// 验证码相关配置
	globalConfig.ValidCode = ValidCodeConfig{
		TypeActiveUser:     int8(requireInt(v, "VALID_CODE_ACTIVE_USER")),
		TypeResetPwd:       int8(requireInt(v, "VALID_CODE_RESET_PWD")),
		ExpireMinutes:      requireInt(v, "VALID_CODE_EXPIRE_MINUTES"),
		Length:             requireInt(v, "VALID_CODE_LEN"),
	}

	// Database配置
	globalConfig.Database = DatabaseConfig{
		Host:     requireString(v, "DB_HOST"),
		Port:     requireInt(v, "DB_PORT"),
		Name:     requireString(v, "DB_NAME"),
		User:     requireString(v, "DB_USER"),
		Password: requireString(v, "DB_PASSWORD"),
		Conn:     requireString(v, "DB_CONN"),
	}

	// Redis配置
	globalConfig.Redis = RedisConfig{
		Host:     requireString(v, "REDIS_HOST"),
		Port:     requireInt(v, "REDIS_PORT"),
		Password: v.GetString("REDIS_PASSWORD"), // Redis密码可以为空
	}

	// Email配置
	globalConfig.Email = EmailConfig{
		SMTPServer:   requireString(v, "SMTP_SERVER"),
		SMTPUser:     requireString(v, "SMTP_USER"),
		SMTPPassword: requireString(v, "SMTP_PASSWORD"),
		SMTPPort:     requireInt(v, "SMTP_PORT"),
		SMTPFrom:     requireString(v, "SMTP_FROM"),
		MaxRetries:   requireInt(v, "EMAIL_MAX_RETRIES"),
		BatchSize:    requireInt(v, "EMAIL_BATCH_SIZE"),
	}

	// Cookie配置
	globalConfig.Cookie = CookieConfig{
		Domain:     requireString(v, "COOKIE_DOMAIN"),
		Secret:     requireString(v, "COOKIE_SECRET"),
		ExpireDays: requireInt(v, "COOKIE_EXPIRE_DAYS"),
		HttpOnly:   requireBool(v, "COOKIE_HTTP_ONLY"),
	}

	// Upload配置
	globalConfig.Upload = UploadConfig{
		Path:    requireString(v, "UPLOAD_PATH"),
		Types:   requireString(v, "UPLOAD_TYPES"),
		MaxSize: requireInt(v, "UPLOAD_MAX_SIZE"),
	}

	// 设置数组配置
	globalConfig.SupportLangs = supportLangs
	globalConfig.Timezones = timezones
	globalConfig.Params = params

	// Log配置
	globalConfig.Log = LogConfig{
		Level:      requireString(v, "LOG_LEVEL"),
		Format:     requireString(v, "LOG_FORMAT"),
		Path:       requireString(v, "LOG_PATH"),
		Filename:   requireString(v, "LOG_FILENAME"),
		MaxSize:    requireInt(v, "LOG_MAX_SIZE"),
		MaxAge:     requireInt(v, "LOG_MAX_AGE"),
		MaxBackups: requireInt(v, "LOG_MAX_BACKUPS"),
		Compress:   requireBool(v, "LOG_COMPRESS"),
		LocalTime:  requireBool(v, "LOG_LOCAL_TIME"),
	}

	// 设置 static 目录
	cwd, _ := os.Getwd()
	staticDir = filepath.Join(cwd, requireString(v, "STATIC_DIR"))

	// 验证配置
	if err := globalConfig.ValidateConfig(); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}

	return nil
}

// 解析JSON字符串
func parseJSON(jsonStr string, target interface{}) error {
	if jsonStr == "" {
		return fmt.Errorf("JSON字符串为空")
	}
	return json.Unmarshal([]byte(jsonStr), target)
}

// GetConfig 获取配置
func GetConfig() *Config {
	return globalConfig
}

// SetConfig 设置配置（仅用于测试）
func SetConfig(config *Config) {
	globalConfig = config
	// 如果 viper 实例为空，创建一个新的实例
	if v == nil {
		v = viper.New()
	}
}

// GetViper 获取viper实例
func GetViper() *viper.Viper {
	return v
}

// GetParams 获取参数配置
func GetParams() map[string]string {
	params := make(map[string]string)
	if globalConfig != nil {
		for _, p := range globalConfig.Params {
			params[p.Key] = p.Val
		}
	}
	return params
}

// GetTimeZone 获取时区配置
func GetTimeZone() []string {
	if globalConfig == nil {
		return []string{}
	}
	return globalConfig.Timezones
}

// GetSupportLangs 获取支持的语言
func GetSupportLangs() []SupportLang {
	if globalConfig == nil {
		return []SupportLang{}
	}
	return globalConfig.SupportLangs
}

// GetSupportLangs 获取支持的语言 (Config 方法)
func (c *Config) GetSupportLangs() []SupportLang {
	return c.SupportLangs
}

// GetTimeZone 获取时区配置 (Config 方法)
func (c *Config) GetTimeZone() []string {
	return c.Timezones
}

// GetStaticPath 获取 static 目录
func GetStaticPath() string {
	return staticDir
}

// UploadsPath returns path to public/uploads folder
func UploadsPath() string {
	return filepath.Join(staticDir, "uploads")
}

// GetViewsPath 获取views目录路径
func GetViewsPath() string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, "website/views")
}

// requireString 获取必需的字符串配置值
func requireString(v *viper.Viper, key string) string {
	value := v.GetString(key)
	if value == "" {
		panic(fmt.Sprintf("配置项 %s 不能为空", key))
	}
	return value
}

// requireInt 获取必需的整数配置值
func requireInt(v *viper.Viper, key string) int {
	if !v.IsSet(key) {
		panic(fmt.Sprintf("配置项 %s 不能为空", key))
	}
	return v.GetInt(key)
}

// requireBool 获取必需的布尔配置值
func requireBool(v *viper.Viper, key string) bool {
	// 对于布尔值，viper.IsSet可能对false值返回false
	// 所以我们需要检查是否在环境变量中明确设置了
	if v.IsSet(key) {
		return v.GetBool(key)
	}
	// 检查环境变量是否存在
	envKey := "IDIARY_" + key
	if _, exists := os.LookupEnv(envKey); exists {
		return v.GetBool(key)
	}
	panic(fmt.Sprintf("配置项 %s 不能为空", key))
}
