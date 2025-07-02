package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// 全局配置实例
var (
	globalConfig *Config
	v            *viper.Viper
	staticDir    string
)

// InitConfig 初始化配置
func InitConfig() error {
	v = viper.New()

	// 环境变量支持
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("IDIARY")

	// 解析JSON格式的环境变量
	var supportLangs []SupportLang
	if err := parseJSON(getEnvOrDefault("SUPPORT_LANGS", v.GetString("SUPPORT_LANGS")), &supportLangs); err != nil {
		return fmt.Errorf("解析SUPPORT_LANGS失败: %w", err)
	}

	var timezones []string
	if err := parseJSON(getEnvOrDefault("TIMEZONES", v.GetString("TIMEZONES")), &timezones); err != nil {
		return fmt.Errorf("解析TIMEZONES失败: %w", err)
	}

	var params []Param
	if err := parseJSON(getEnvOrDefault("PARAMS", v.GetString("PARAMS")), &params); err != nil {
		return fmt.Errorf("解析PARAMS失败: %w", err)
	}

	// 初始化配置结构体
	globalConfig = &Config{}

	// 手动设置值
	globalConfig.Env = getEnvOrDefault("ENV", v.GetString("ENV"))
	globalConfig.App = AppConfig{
		GinMode:         getEnvOrDefault("GIN_MODE", v.GetString("GIN_MODE")),
		Domain:          getEnvOrDefault("DOMAIN", v.GetString("DOMAIN")),
		Port:            getEnvOrDefault("PORT", v.GetString("PORT")),
		SummaryMax:      getEnvIntOrDefault("SUMMARY_MAX", v.GetInt("SUMMARY_MAX")),
		Suffix:          getEnvOrDefault("SUFFIX", v.GetString("SUFFIX")),
		PageSize:        getEnvIntOrDefault("PAGE_SIZE", v.GetInt("PAGE_SIZE")),
		MaxDaybookLimit: getEnvIntOrDefault("MAX_DAYBOOK_LIMIT", v.GetInt("MAX_DAYBOOK_LIMIT")),
		HashSalt:        getEnvOrDefault("HASH_SALT", v.GetString("HASH_SALT")),
		StaticDir:       getEnvOrDefault("STATIC_DIR", v.GetString("STATIC_DIR")),
	}
	globalConfig.Database = DatabaseConfig{
		Host:     getEnvOrDefault("DB_HOST", v.GetString("DB_HOST")),
		Port:     getEnvIntOrDefault("DB_PORT", v.GetInt("DB_PORT")),
		Name:     getEnvOrDefault("DB_NAME", v.GetString("DB_NAME")),
		User:     getEnvOrDefault("DB_USER", v.GetString("DB_USER")),
		Password: getEnvOrDefault("DB_PASSWORD", v.GetString("DB_PASSWORD")),
		Conn:     getEnvOrDefault("DB_CONN", v.GetString("DB_CONN")),
	}
	globalConfig.Redis = RedisConfig{
		Host:     getEnvOrDefault("REDIS_HOST", v.GetString("REDIS_HOST")),
		Port:     getEnvIntOrDefault("REDIS_PORT", v.GetInt("REDIS_PORT")),
		Password: getEnvOrDefault("REDIS_PASSWORD", v.GetString("REDIS_PASSWORD")),
	}
	globalConfig.Email = EmailConfig{
		SMTPServer:   getEnvOrDefault("SMTP_SERVER", v.GetString("SMTP_SERVER")),
		SMTPUser:     getEnvOrDefault("SMTP_USER", v.GetString("SMTP_USER")),
		SMTPPassword: getEnvOrDefault("SMTP_PASSWORD", v.GetString("SMTP_PASSWORD")),
		SMTPPort:     getEnvIntOrDefault("SMTP_PORT", v.GetInt("SMTP_PORT")),
		SMTPFrom:     getEnvOrDefault("SMTP_FROM", v.GetString("SMTP_FROM")),
		MaxRetries:   getEnvIntOrDefault("EMAIL_MAX_RETRIES", v.GetInt("EMAIL_MAX_RETRIES")),
		BatchSize:    getEnvIntOrDefault("EMAIL_BATCH_SIZE", v.GetInt("EMAIL_BATCH_SIZE")),
	}
	globalConfig.Cookie = CookieConfig{
		Domain:     getEnvOrDefault("COOKIE_DOMAIN", v.GetString("COOKIE_DOMAIN")),
		Secret:     getEnvOrDefault("COOKIE_SECRET", v.GetString("COOKIE_SECRET")),
		ExpireDays: getEnvIntOrDefault("COOKIE_EXPIRE_DAYS", v.GetInt("COOKIE_EXPIRE_DAYS")),
		HttpOnly:   getEnvBoolOrDefault("COOKIE_HTTP_ONLY", v.GetBool("COOKIE_HTTP_ONLY")),
	}
	globalConfig.Upload = UploadConfig{
		Path:    getEnvOrDefault("UPLOAD_PATH", v.GetString("UPLOAD_PATH")),
		Types:   getEnvOrDefault("UPLOAD_TYPES", v.GetString("UPLOAD_TYPES")),
		MaxSize: getEnvIntOrDefault("UPLOAD_MAX_SIZE", v.GetInt("UPLOAD_MAX_SIZE")),
	}
	globalConfig.SupportLangs = supportLangs
	globalConfig.Timezones = timezones
	globalConfig.Params = params
	globalConfig.Log = LogConfig{
		Level:      getEnvOrDefault("LOG_LEVEL", v.GetString("LOG_LEVEL")),
		Format:     getEnvOrDefault("LOG_FORMAT", v.GetString("LOG_FORMAT")),
		Path:       getEnvOrDefault("LOG_PATH", v.GetString("LOG_PATH")),
		Filename:   getEnvOrDefault("LOG_FILENAME", v.GetString("LOG_FILENAME")),
		MaxSize:    getEnvIntOrDefault("LOG_MAX_SIZE", v.GetInt("LOG_MAX_SIZE")),
		MaxAge:     getEnvIntOrDefault("LOG_MAX_AGE", v.GetInt("LOG_MAX_AGE")),
		MaxBackups: getEnvIntOrDefault("LOG_MAX_BACKUPS", v.GetInt("LOG_MAX_BACKUPS")),
		Compress:   getEnvBoolOrDefault("LOG_COMPRESS", v.GetBool("LOG_COMPRESS")),
		LocalTime:  getEnvBoolOrDefault("LOG_LOCAL_TIME", v.GetBool("LOG_LOCAL_TIME")),
	}

	// 设置 static 目录
	cwd, _ := os.Getwd()
	staticDir = filepath.Join(cwd, getEnvOrDefault("STATIC_DIR", v.GetString("STATIC_DIR")))

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
}

// GetViper 获取viper实例
func GetViper() *viper.Viper {
	return v
}

// GetParams 获取参数配置
func GetParams() map[string]string {
	params := make(map[string]string)
	for _, p := range globalConfig.Params {
		params[p.Key] = p.Val
	}
	return params
}

// GetTimeZone 获取时区配置
func GetTimeZone() []string {
	return globalConfig.Timezones
}

// GetSupportLangs 获取支持的语言
func GetSupportLangs() []SupportLang {
	return globalConfig.SupportLangs
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

// getEnvOrDefault 从环境变量获取字符串值，如果不存在则返回默认值
func getEnvOrDefault(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

// getEnvIntOrDefault 从环境变量获取整数值，如果不存在或无法转换则返回默认值
func getEnvIntOrDefault(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}

	return val
}

// getEnvBoolOrDefault 从环境变量获取布尔值，如果不存在或无法转换则返回默认值
func getEnvBoolOrDefault(key string, defaultVal bool) bool {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}

	// 处理常见的布尔值表示
	valStr = strings.ToLower(valStr)
	if valStr == "true" || valStr == "1" || valStr == "yes" || valStr == "y" {
		return true
	}
	if valStr == "false" || valStr == "0" || valStr == "no" || valStr == "n" {
		return false
	}

	return defaultVal
}
