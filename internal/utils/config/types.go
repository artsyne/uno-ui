package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// Config 配置结构体
type Config struct {
	Env          string           `json:"env" mapstructure:"ENV" validate:"required,oneof=development production test"`
	App          AppConfig        `json:"app" validate:"required"`
	AutoSave     AutoSaveConfig   `json:"auto_save"`
	User         UserConfig       `json:"user" validate:"required"`
	ValidCode    ValidCodeConfig  `json:"valid_code" validate:"required"`
	Database     DatabaseConfig   `json:"database" validate:"required"`
	Redis        RedisConfig      `json:"redis" validate:"required"`
	RateLimit    RateLimitConfig  `json:"rate_limit" validate:"required"`
	BruteForce   BruteForceConfig `json:"brute_force" validate:"required"`
	Email        EmailConfig      `json:"email" validate:"required"`
	Cookie       CookieConfig     `json:"cookie" validate:"required"`
	Upload       UploadConfig     `json:"upload" validate:"required"`
	SupportLangs []SupportLang    `json:"support_langs" validate:"required,dive"`
	Timezones    []string         `json:"timezones" validate:"required,dive,timezone"`
	Params       []Param          `json:"params" validate:"required,dive"`
	Log          LogConfig        `json:"log" validate:"required"`
}

// AppConfig 应用配置
type AppConfig struct {
	GinMode              string `json:"gin_mode" mapstructure:"GIN_MODE" validate:"required,oneof=debug release"`
	Domain               string `json:"domain" mapstructure:"DOMAIN" validate:"required,url"`
	Port                 string `json:"port" mapstructure:"PORT" validate:"required,numeric"`
	SummaryMax           int    `json:"summary_max" mapstructure:"SUMMARY_MAX" validate:"required,min=300,max=2048"`
	Suffix               string `json:"suffix" mapstructure:"SUFFIX" validate:"required"`
	PageSize             int    `json:"page_size" mapstructure:"PAGE_SIZE" validate:"required,min=1,max=200"`
	MaxDiaryLimit        int    `json:"max_diary_limit" mapstructure:"MAX_DIARY_LIMIT" validate:"required,min=1"`
	MaxEntryLength       int    `json:"max_entry_length" mapstructure:"MAX_ENTRY_LENGTH" validate:"required,min=1"`
	MaxCommentLength     int    `json:"max_comment_length" mapstructure:"MAX_COMMENT_LENGTH" validate:"required,min=1"`
	MaxDescriptionLength int    `json:"max_description_length" mapstructure:"MAX_DESCRIPTION_LENGTH" validate:"required,min=1"`
	StaticDir            string `json:"static_dir" mapstructure:"STATIC_DIR" validate:"required"`
	LocalesDir           string `json:"locales_dir" mapstructure:"LOCALES_DIR" validate:"required"`
	DefaultLang          string `json:"default_lang" mapstructure:"DEFAULT_LANG" validate:"required"`
}

// AutoSaveConfig 自动保存配置
type AutoSaveConfig struct {
	Enabled  bool `json:"enabled" mapstructure:"AUTO_SAVE_ENABLED"`
	Interval int  `json:"interval" mapstructure:"AUTO_SAVE_INTERVAL"`
}

// UserConfig 用户相关配置
type UserConfig struct {
	// 密码规则
	PwdMinLen int `json:"pwd_min_len" mapstructure:"PWD_MIN_LEN" validate:"required,min=1,max=6"`
	PwdMaxLen int `json:"pwd_max_len" mapstructure:"PWD_MAX_LEN" validate:"required,min=6,max=100"`

	// 用户状态
	StatusOK         int8 `json:"status_ok" mapstructure:"USER_STATUS_OK" validate:"required"`
	StatusWaitActive int8 `json:"status_wait_active" mapstructure:"USER_STATUS_WAIT_ACTIVE" validate:"required"`
	StatusDeleted    int8 `json:"status_deleted" mapstructure:"USER_STATUS_DEL" validate:"gte=0"`

	// 用户分组
	DefaultGroup int8 `json:"default_group" mapstructure:"DEFAULT_USER_GROUP" validate:"required"`

	// Handle 修改冷却时间（天）
	HandleCooldownDays int `json:"handle_cooldown_days" mapstructure:"HANDLE_COOLDOWN_DAYS" validate:"required,min=0"`
}

// ValidCodeConfig 验证码相关配置
type ValidCodeConfig struct {
	// 验证码长度
	Length int `json:"length" mapstructure:"VALID_CODE_LEN" validate:"required,min=6,max=128"`
	// 验证码类型
	TypeActiveUser int8 `json:"type_active_user" mapstructure:"VALID_CODE_ACTIVE_USER" validate:"required"`
	TypeResetPwd   int8 `json:"type_reset_pwd" mapstructure:"VALID_CODE_RESET_PWD" validate:"required"`
	// 验证码过期分钟
	ExpireMinutes int `json:"valid_code_expire_minutes" mapstructure:"VALID_CODE_EXPIRE_MINUTES" validate:"required"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `json:"host" mapstructure:"DB_HOST" validate:"required"`
	Port     int    `json:"port" mapstructure:"DB_PORT" validate:"required,min=1,max=65535"`
	Name     string `json:"name" mapstructure:"DB_NAME" validate:"required"`
	User     string `json:"user" mapstructure:"DB_USER" validate:"required"`
	Password string `json:"password" mapstructure:"DB_PASSWORD"` // 在非生产环境可以为空
	Conn     string `json:"conn" mapstructure:"DB_CONN" validate:"required"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `json:"host" mapstructure:"REDIS_HOST" validate:"required"`
	Port     int    `json:"port" mapstructure:"REDIS_PORT" validate:"required,min=1,max=65535"`
	Password string `json:"password" mapstructure:"REDIS_PASSWORD"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	GeneralLimit float64 `json:"general_limit" mapstructure:"RATE_LIMIT_GENERAL" validate:"required,min=0.1"`
	GeneralBurst int     `json:"general_burst" mapstructure:"RATE_LIMIT_GENERAL_BURST" validate:"required,min=1"`
	AuthLimit    float64 `json:"auth_limit" mapstructure:"RATE_LIMIT_AUTH" validate:"required,min=0.1"`
	AuthBurst    int     `json:"auth_burst" mapstructure:"RATE_LIMIT_AUTH_BURST" validate:"required,min=1"`
}

// BruteForceConfig 暴力破解防护配置
type BruteForceConfig struct {
	MaxFailures    int `json:"max_failures" mapstructure:"BRUTE_FORCE_MAX_FAILURES" validate:"required,min=3"`
	LockoutMinutes int `json:"lockout_minutes" mapstructure:"BRUTE_FORCE_LOCKOUT_MINUTES" validate:"required,min=1"`
}

// EmailConfig 邮件配置
type EmailConfig struct {
	SMTPServer   string `json:"smtp_server" mapstructure:"SMTP_SERVER" validate:"required"`
	SMTPUser     string `json:"smtp_user" mapstructure:"SMTP_USER" validate:"required,email"`
	SMTPPassword string `json:"smtp_password" mapstructure:"SMTP_PASSWORD"` // 在非生产环境可以为空
	SMTPPort     int    `json:"smtp_port" mapstructure:"SMTP_PORT" validate:"required,min=1,max=65535"`
	SMTPFrom     string `json:"smtp_from" mapstructure:"SMTP_FROM" validate:"required,email"`
	MaxRetries   int    `json:"max_retries" mapstructure:"EMAIL_MAX_RETRIES" validate:"required,min=1"`
	BatchSize    int    `json:"batch_size" mapstructure:"EMAIL_BATCH_SIZE" validate:"required,min=1"`
}

// CookieConfig Cookie配置
type CookieConfig struct {
	Domain     string `json:"domain" mapstructure:"COOKIE_DOMAIN" validate:"required"`
	Secret     string `json:"secret" mapstructure:"COOKIE_SECRET" validate:"required,min=16"`
	ExpireDays int    `json:"expire_days" mapstructure:"COOKIE_EXPIRE_DAYS" validate:"required,min=1,max=365"`
	HttpOnly   bool   `json:"http_only" mapstructure:"COOKIE_HTTP_ONLY"`
}

// UploadConfig 上传配置
type UploadConfig struct {
	Path    string `json:"path" mapstructure:"UPLOAD_PATH" validate:"required"`
	Types   string `json:"types" mapstructure:"UPLOAD_TYPES" validate:"required"`
	MaxSize int    `json:"max_size" mapstructure:"UPLOAD_MAX_SIZE" validate:"required,min=1,max=10485760"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `json:"level" mapstructure:"LOG_LEVEL" validate:"required,oneof=debug info warn error"`
	Format     string `json:"format" mapstructure:"LOG_FORMAT" validate:"required,oneof=text json"`
	Path       string `json:"path" mapstructure:"LOG_PATH" validate:"required"`
	Filename   string `json:"filename" mapstructure:"LOG_FILENAME" validate:"required"`
	MaxSize    int    `json:"max_size" mapstructure:"LOG_MAX_SIZE" validate:"required,min=1,max=500"`
	MaxAge     int    `json:"max_age" mapstructure:"LOG_MAX_AGE" validate:"required,min=1"`
	MaxBackups int    `json:"max_backups" mapstructure:"LOG_MAX_BACKUPS" validate:"required,min=1"`
	Compress   bool   `json:"compress" mapstructure:"LOG_COMPRESS"`
	LocalTime  bool   `json:"local_time" mapstructure:"LOG_LOCAL_TIME"`
}

// Param 参数配置
type Param struct {
	Key string `json:"key" validate:"required"`
	Val string `json:"val" validate:"required"`
}

// SupportLang 支持的语言
type SupportLang struct {
	Code      string `json:"code" validate:"required"`
	Label     string `json:"label" validate:"required"`
	LocalText string `json:"local-text" validate:"required"`
}

// 自定义验证器
var validate *validator.Validate

func init() {
	validate = validator.New()

	// 注册自定义验证规则
	_ = validate.RegisterValidation("dir", validateDir)
	_ = validate.RegisterValidation("timezone", validateTimezone)
}

// ValidateConfig 验证配置
func (c *Config) ValidateConfig() error {
	if err := validate.Struct(c); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		var errMsgs []string
		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, formatValidationError(err))
		}
		return fmt.Errorf("config validation failed: %s", strings.Join(errMsgs, "; "))
	}
	return nil
}

// validateDir 验证目录是否存在或可创建
func validateDir(fl validator.FieldLevel) bool {
	dir := fl.Field().String()

	// 获取绝对路径
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return false
	}

	// 检查目录是否存在
	info, err := os.Stat(absPath)
	if err == nil {
		// 目录存在，检查是否是目录
		return info.IsDir()
	}

	// 如果错误不是"不存在"，则返回 false
	if !os.IsNotExist(err) {
		return false
	}

	// 目录不存在，尝试创建
	err = os.MkdirAll(absPath, 0755)
	if err != nil { // 创建失败
		return false
	}

	// 再次检查目录是否创建成功
	info, err = os.Stat(absPath)
	return err == nil && info.IsDir()
}

// validateTimezone 验证时区数组中的每个时区是否有效
func validateTimezone(fl validator.FieldLevel) bool {
	tz := fl.Field().String()

	if tz == "UTC" {
		return true
	}
	if !strings.Contains(tz, "/") && tz != "UTC" {
		return false
	}

	// 尝试加载时区
	_, err := time.LoadLocation(tz)
	return err == nil
}

// 格式化验证错误信息
func formatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", err.Field(), err.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", err.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", err.Field())
	default:
		return fmt.Sprintf("%s failed %s validation", err.Field(), err.Tag())
	}
}
