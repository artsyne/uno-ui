package bootstrap

import (
	"fmt"
	"log/slog"

	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/core/database"
	"github.com/biuren/hebe/internal/core/di"
)

// Initializer 核心基础设施初始化器
type Initializer struct {
	config    *config.Config
	container *di.Container
}

// NewInitializer 创建核心初始化器
func NewInitializer() *Initializer {
	return &Initializer{}
}

// InitializeCore 初始化核心基础设施
// 由于配置和日志已在main.go中初始化，这里只初始化数据库和DI容器
func (i *Initializer) InitializeCore() error {
	// 1. 获取已初始化的配置
	i.config = config.GetConfig()
	if i.config == nil {
		return fmt.Errorf("配置未初始化")
	}

	// 2. 初始化数据库
	if err := i.initDatabase(); err != nil {
		return fmt.Errorf("数据库初始化失败: %w", err)
	}

	// 3. 初始化DI容器
	if err := i.initContainer(); err != nil {
		return fmt.Errorf("DI容器初始化失败: %w", err)
	}

	return nil
}

// GetContainer 获取DI容器
func (i *Initializer) GetContainer() *di.Container {
	return i.container
}

// GetConfig 获取配置
func (i *Initializer) GetConfig() *config.Config {
	return i.config
}

// initDatabase 初始化数据库连接
func (i *Initializer) initDatabase() error {
	if err := database.Connect(); err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	slog.Info("数据库连接成功")
	return nil
}

// initContainer 初始化DI容器
func (i *Initializer) initContainer() error {
	i.container = di.New()

	// 注册配置到DI容器
	if err := di.Value(i.container, i.config); err != nil {
		return fmt.Errorf("注册配置到DI容器失败: %w", err)
	}

	slog.Info("成功初始化DI容器并注册配置")
	return nil
}
