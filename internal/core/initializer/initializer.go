package bootstrap

import (
	"log/slog"

	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/core/database"
	"github.com/biuren/hebe/internal/core/di"
	"github.com/biuren/hebe/internal/utils/errors"
)

// Initializer 核心基础设施初始化器
// 负责初始化应用程序的核心基础设施，包括数据库连接和依赖注入容器
type Initializer struct {
	config    *config.Config  // 应用程序配置
	container *di.Container   // 依赖注入容器
}

// NewInitializer 创建核心初始化器
// 使用提供的配置创建一个新的初始化器实例
//
// 参数:
//   - config: 应用程序配置对象
//
// 返回:
//   - *Initializer: 初始化器实例
func NewInitializer(config *config.Config) *Initializer {
	return &Initializer{
		config: config,
	}
}

// InitializeCore 初始化核心基础设施
// 按顺序初始化数据库连接和依赖注入容器
// 由于配置和日志已在main.go中初始化，这里只初始化数据库和DI容器
//
// 返回:
//   - error: 初始化失败时返回错误，成功时返回nil
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
// 返回已初始化的依赖注入容器实例
//
// 返回:
//   - *di.Container: DI容器实例，如果未初始化则返回nil
func (i *Initializer) GetContainer() *di.Container {
	return i.container
}

// initDatabase 初始化数据库连接
// 使用配置中的数据库信息建立数据库连接
//
// 返回:
//   - error: 连接失败时返回错误，成功时返回nil
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
// 创建新的DI容器并注册基础服务（如配置对象）
//
// 返回:
//   - error: 初始化失败时返回错误，成功时返回nil
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
