package bootstrap

import (
	"log/slog"

	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/core/database"
	"github.com/biuren/hebe/internal/core/di"
	"github.com/biuren/hebe/internal/utils/error"
)

// Initializer 核心基础设施初始化器
type Initializer struct {
	config    *config.Config
	container *di.Container
}

// NewInitializer 创建核心初始化器
func NewInitializer(config *config.Config) *Initializer {
	return &Initializer{
		config: config,

	}
}

// InitializeCore 初始化核心基础设施
// 由于配置和日志已在main.go中初始化，这里只初始化数据库和DI容器
func (i *Initializer) InitializeCore() error {
	// 1. 初始化数据库
	if err := i.initDatabase(); err != nil {
		return err // 已经是统一错误格式
	}

	// 2. 初始化DI容器
	if err := i.initContainer(); err != nil {
		return err // 已经是统一错误格式
	}

	return nil
}

// GetContainer 获取DI容器
func (i *Initializer) GetContainer() *di.Container {
	return i.container
}

// initDatabase 初始化数据库连接
func (i *Initializer) initDatabase() error {
	// 检查配置是否存在
	if i.config == nil {
		return errors.DatabaseError.WithDetails("配置为空，无法初始化数据库")
	}
	
	// 使用配置中的数据库信息直接连接
	if err := database.ConnectWithConfig(i.config.Database); err != nil {
		return errors.DatabaseError.WithDetails("数据库连接失败: " + err.Error())
	}

	slog.Info("数据库连接成功")
	return nil
}

// initContainer 初始化DI容器
func (i *Initializer) initContainer() error {

	// 初始化DI容器
	i.container = di.New()

	// 注册配置到DI容器，如果已存在则跳过
	if !di.Has[*config.Config](i.container) {
		if err := di.Value(i.container, i.config); err != nil {
			return errors.DIServiceAlreadyRegistered.WithDetails("注册配置到DI容器失败: " + err.Error())
		}
	}

	slog.Info("成功初始化DI容器并注册配置")
	return nil
}
