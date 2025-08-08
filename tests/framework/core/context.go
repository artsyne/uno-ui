package framework

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/biuren/hebe/internal/core/database"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/post"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/config"
	framework "github.com/biuren/hebe/tests/framework/helpers"
	. "github.com/biuren/hebe/tests/framework/types"
	"github.com/jmoiron/sqlx"
)

// 定义本地类型，避免循环导入
type DataIsolatorImpl struct {
	testID      string
	prefix      string
	config      *TestConfig
	createdData map[string][]interface{}
	mu          sync.RWMutex
}

type DataGeneratorImpl struct {
	isolator DataIsolatorInterface
	random   string
}

// CleanupManagerImpl 实现在 helpers/cleanup.go 中

// TestContext 统一的测试上下文
// 替代当前的4种不同测试模式，提供一致的测试环境
type TestContext struct {
	// 基础信息
	testID  string
	name    string
	started time.Time
	config  *TestConfig

	// 数据管理
	database DatabaseManagerInterface
	dataGen  DataGeneratorInterface
	isolator DataIsolatorInterface
	cleanup  CleanupManagerInterface

	// 服务管理
	services *ServiceRegistry

	// 生命周期管理
	ctx          context.Context
	cancel       context.CancelFunc
	cleanupFuncs []func() error
	mu           sync.RWMutex
	stats        *TestStats
}

// TestContextOption 测试上下文选项函数类型
type TestContextOption func(*TestContext)

// NewTestContext 创建新的测试上下文
// 这是唯一的测试环境入口点，统一替代所有现有模式
func NewTestContext(testName string, options ...TestContextOption) *TestContext {
	testID := generateTestID(testName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	tc := &TestContext{
		testID:       testID,
		name:         testName,
		started:      time.Now(),
		ctx:          ctx,
		cancel:       cancel,
		cleanupFuncs: make([]func() error, 0),
		stats: &TestStats{
			TestID:    testID,
			TestName:  testName,
			StartTime: time.Now(),
			Success:   true,
		},
	}

	// 应用配置选项
	for _, opt := range options {
		opt(tc)
	}

	// 初始化核心组件
	if err := tc.initializeCore(); err != nil {
		tc.stats.Success = false
		tc.stats.ErrorMessage = err.Error()
		panic(fmt.Sprintf("初始化测试上下文失败: %v", err))
	}

	slog.Debug("测试上下文创建成功",
		"test_id", testID,
		"test_name", testName,
		"prefix", tc.isolator.GetPrefix())

	return tc
}

// initializeCore 初始化核心组件
func (tc *TestContext) initializeCore() error {
	// 初始化配置
	tc.config = GetTestConfig()

	// 初始化全局配置包（用于服务层）
	if err := initGlobalConfig(); err != nil {
		return fmt.Errorf("初始化全局配置失败: %w", err)
	}

	// 初始化数据库管理器
	tc.database = GetDatabaseManager()
	if err := tc.database.Initialize(); err != nil {
		return fmt.Errorf("初始化数据库管理器失败: %w", err)
	}

	// 初始化数据隔离器
	tc.isolator = newDataIsolator(tc.testID, tc.config)

	// 初始化数据生成器
	tc.dataGen = newDataGenerator(tc.isolator)

	// 初始化清理管理器
	tc.cleanup = newCleanupManager(tc.database.GetRawDB(), tc.config, tc.isolator)

	// 初始化服务注册表
	tc.services = NewServiceRegistry()
	if err := tc.services.InitializeServices(tc.database.GetDB()); err != nil {
		return fmt.Errorf("初始化服务注册表失败: %w", err)
	}

	return nil
}

// generateTestID 生成测试ID
func generateTestID(testName string) string {
	// 如果测试名称为空，生成默认ID
	if testName == "" {
		return "default_test"
	}
	// 在测试环境中，为了测试的一致性，直接返回测试名称
	// 在实际使用中，可以根据需要添加时间戳以确保唯一性
	return testName
}

// 接口实现

// GetTestID 获取测试ID
func (tc *TestContext) GetTestID() string {
	return tc.testID
}

// GetConfig 获取配置
func (tc *TestContext) GetConfig() TestConfigInterface {
	return tc.config
}

// GetDatabase 获取数据库管理器
func (tc *TestContext) GetDatabase() DatabaseManagerInterface {
	return tc.database
}

// GetDataGenerator 获取数据生成器
func (tc *TestContext) GetDataGenerator() DataGeneratorInterface {
	return tc.dataGen
}

// GetIsolator 获取数据隔离器
func (tc *TestContext) GetIsolator() DataIsolatorInterface {
	return tc.isolator
}

// GetCleanup 获取清理管理器
func (tc *TestContext) GetCleanup() CleanupManagerInterface {
	return tc.cleanup
}

// GetServices 获取服务注册表
func (tc *TestContext) GetServices() ServiceRegistryInterface {
	return tc.services
}

// GetHTTP 获取HTTP测试助手（暂未实现）
func (tc *TestContext) GetHTTP() HTTPTestHelperInterface {
	// TODO: 实现HTTP测试助手
	return nil
}

// Context 获取上下文
func (tc *TestContext) Context() context.Context {
	return tc.ctx
}

// Cleanup 执行清理
func (tc *TestContext) Cleanup() error {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	// 取消上下文
	if tc.cancel != nil {
		tc.cancel()
	}

	// 更新统计信息
	tc.stats.EndTime = time.Now()
	tc.stats.Duration = tc.stats.EndTime.Sub(tc.stats.StartTime)
	tc.stats.DataCreated = tc.isolator.GetDataCount()

	// 执行自定义清理函数
	var lastErr error
	for i, cleanupFunc := range tc.cleanupFuncs {
		if err := cleanupFunc(); err != nil {
			slog.Error("执行自定义清理函数失败", "index", i, "error", err)
			lastErr = err
		}
	}

	// 执行数据清理
	if tc.cleanup != nil {
		if err := tc.cleanup.Cleanup(tc.config.Cleanup.Strategy); err != nil {
			slog.Error("数据清理失败", "error", err)
			lastErr = err
		}
		tc.stats.DataCleaned = tc.isolator.GetDataCount() // 清理后的数据量
	}

	// 输出统计信息
	tc.logStats()

	if lastErr != nil {
		tc.stats.Success = false
		tc.stats.ErrorMessage = lastErr.Error()
	}

	return lastErr
}

// AddCleanupFunc 添加自定义清理函数
func (tc *TestContext) AddCleanupFunc(fn func() error) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.cleanupFuncs = append(tc.cleanupFuncs, fn)
}

// ForceCleanupDatabase 强制清理数据库
// 无论配置如何，都会立即清理测试数据
func (tc *TestContext) ForceCleanupDatabase() error {
	if tc.cleanup == nil {
		return fmt.Errorf("清理管理器未初始化")
	}
	return tc.cleanup.ForceCleanupDatabase()
}

// logStats 输出统计信息
func (tc *TestContext) logStats() {
	slog.Debug("测试统计信息",
		"test_id", tc.stats.TestID,
		"test_name", tc.stats.TestName,
		"duration", tc.stats.Duration,
		"data_created", tc.stats.DataCreated,
		"data_cleaned", tc.stats.DataCleaned,
		"success", tc.stats.Success)
}

// GetStats 获取统计信息
func (tc *TestContext) GetStats() *TestStats {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	return tc.stats
}

// 选项函数

// WithDatabase 数据库测试选项
func WithDatabase() TestContextOption {
	return func(tc *TestContext) {
		// 数据库总是初始化，此选项为兼容性保留
	}
}

// WithServices 服务测试选项
func WithServices(services ...string) TestContextOption {
	return func(tc *TestContext) {
		// 服务总是初始化，此选项为兼容性保留
	}
}

// WithCleanupEnabled 启用清理选项
func WithCleanupEnabled(enabled bool) TestContextOption {
	return func(tc *TestContext) {
		if tc.config != nil {
			tc.config.Cleanup.Enabled = enabled
		}
	}
}

// WithTimeout 设置超时选项
func WithTimeout(timeout time.Duration) TestContextOption {
	return func(tc *TestContext) {
		if tc.ctx != nil && tc.cancel != nil {
			tc.cancel()
		}
		tc.ctx, tc.cancel = context.WithTimeout(context.Background(), timeout)
	}
}

// ServiceRegistry 服务注册表
type ServiceRegistry struct {
	userService  user.ServiceInterface
	diaryService diary.ServiceInterface
	postService  post.ServiceInterface
	initialized  bool
	mu           sync.RWMutex
}

// NewServiceRegistry 创建服务注册表
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{}
}

// InitializeServices 初始化服务
func (sr *ServiceRegistry) InitializeServices(db *sqlx.DB) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	if sr.initialized {
		return nil
	}

	// 创建基础设施层仓储
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	
	// 创建任务服务
	config := config.GetConfig()
	taskService := task.NewService(taskRepo, config)
	
	// 初始化各个服务，注入仓储依赖
	sr.userService = user.NewService(userRepo, taskService, config)
	diaryRepo := repository.NewDiaryRepository()
	sr.diaryService = diary.NewService(diaryRepo)
	postRepo := repository.NewPostRepository(db)
	sr.postService = post.NewService(userRepo, taskService, diaryRepo, postRepo)

	sr.initialized = true

	slog.Debug("服务注册表初始化完成")

	return nil
}

// GetUserService 获取用户服务
func (sr *ServiceRegistry) GetUserService() user.ServiceInterface {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	return sr.userService
}

// GetDiaryService 获取日记服务
func (sr *ServiceRegistry) GetDiaryService() diary.ServiceInterface {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	return sr.diaryService
}

// GetPostService 获取文章服务
func (sr *ServiceRegistry) GetPostService() post.ServiceInterface {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	return sr.postService
}



// 以下是本地实现的构造函数，避免循环导入

// generatePrefix 生成唯一前缀
func generatePrefix(testID string, config IsolationConfig) string {
	now := time.Now()

	// 生成随机字符串
	randomBytes := make([]byte, 3)
	rand.Read(randomBytes)
	random := hex.EncodeToString(randomBytes)

	// 根据配置格式生成前缀
	prefix := config.PrefixFormat
	prefix = strings.ReplaceAll(prefix, "{date}", now.Format("0102")) // MMDD
	prefix = strings.ReplaceAll(prefix, "{time}", now.Format("1504")) // HHMM
	prefix = strings.ReplaceAll(prefix, "{random}", random)           // 随机字符串
	prefix = strings.ReplaceAll(prefix, "{test_id}", testID)          // 测试ID

	// 确保前缀长度不超过限制
	if len(prefix) > config.MaxPrefixLen {
		prefix = prefix[:config.MaxPrefixLen]
	}

	return prefix
}

// newDataIsolator 创建数据隔离器
func newDataIsolator(testID string, config *TestConfig) DataIsolatorInterface {
	prefix := generatePrefix(testID, config.Isolation)

	return &DataIsolatorImpl{
		testID:      testID,
		prefix:      prefix,
		config:      config,
		createdData: make(map[string][]interface{}),
	}
}

// newDataGenerator 创建数据生成器
func newDataGenerator(isolator DataIsolatorInterface) DataGeneratorInterface {
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)

	return &DataGeneratorImpl{
		isolator: isolator,
		random:   hex.EncodeToString(randomBytes),
	}
}

// newCleanupManager 创建清理管理器
func newCleanupManager(db *sql.DB, config *TestConfig, isolator DataIsolatorInterface) CleanupManagerInterface {
	return framework.NewCleanupManager(db, config, isolator)
}

// DataIsolatorImpl 接口实现
func (di *DataIsolatorImpl) IsolateString(value string) string {
	return fmt.Sprintf("%s_%s", di.prefix, value)
}

func (di *DataIsolatorImpl) IsolateEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return di.IsolateString(email)
	}
	return fmt.Sprintf("%s_%s@%s", di.prefix, parts[0], parts[1])
}

func (di *DataIsolatorImpl) IsolateUsername(username string) string {
	isolated := di.IsolateString(username)
	if len(isolated) > 90 {
		isolated = isolated[:90]
	}
	return isolated
}

func (di *DataIsolatorImpl) IsolateTitle(title string) string {
	isolated := di.IsolateString(title)
	if len(isolated) > 180 {
		isolated = isolated[:180]
	}
	return isolated
}

func (di *DataIsolatorImpl) TrackCreated(table string, id interface{}) {
	di.mu.Lock()
	defer di.mu.Unlock()
	di.createdData[table] = append(di.createdData[table], id)
}

func (di *DataIsolatorImpl) GetCreatedData() map[string][]interface{} {
	di.mu.RLock()
	defer di.mu.RUnlock()
	result := make(map[string][]interface{})
	for table, ids := range di.createdData {
		result[table] = make([]interface{}, len(ids))
		copy(result[table], ids)
	}
	return result
}

func (di *DataIsolatorImpl) GetPrefix() string {
	return di.prefix
}

func (di *DataIsolatorImpl) GetTestID() string {
	return di.testID
}

func (di *DataIsolatorImpl) GetDataCount() int {
	di.mu.RLock()
	defer di.mu.RUnlock()
	count := 0
	for _, ids := range di.createdData {
		count += len(ids)
	}
	return count
}

func (di *DataIsolatorImpl) GetTableCount() int {
	di.mu.RLock()
	defer di.mu.RUnlock()
	return len(di.createdData)
}

// DataGeneratorImpl 接口实现
func (dg *DataGeneratorImpl) CreateUser(options ...UserOption) *user.User {
	opts := &UserOptions{
		Username: dg.isolator.IsolateUsername("user"),
		Email:    dg.isolator.IsolateEmail("user@test.com"),
		Password: "password123",
		FullName: "Test User",
		TimeZone: "Asia/Shanghai",
		Lang:     "zh-cn",
	}

	for _, opt := range options {
		opt(opts)
	}

	return &user.User{
		UserName:         opts.Username,
		Email:            opts.Email,
		Pwd:              opts.Password,
		FullName:         opts.FullName,
		TimeZone:         opts.TimeZone,
		Lang:             opts.Lang,
		CreatedUTCTime:   time.Now().Unix(),
		LastLoginUTCTime: time.Now().Unix(),
		Status:           int8(3),
	}
}

func (dg *DataGeneratorImpl) CreateDiary(userID int64, options ...DiaryOption) *diary.Diary {
	opts := &DiaryOptions{
		Title:         dg.isolator.IsolateTitle("diary"),
		Description:   "Test diary description",
		ViewSetting:   1,
		SocialSetting: 1,
		ThemeID:       1,
	}

	for _, opt := range options {
		opt(opts)
	}

	now := time.Now().Unix()
	return &diary.Diary{
		Title:             opts.Title,
		Remark:            opts.Description,
		ViewSetting:       int8(opts.ViewSetting),
		SocialSetting:     int8(opts.SocialSetting),
		ThemeID:           int8(opts.ThemeID),
		AddUTCTime:        now,
		AddUserID:         userID,
		LastActionUTCTime: now,
		LastActionUserID:  userID,
		Status:            int8(1),
	}
}

func (dg *DataGeneratorImpl) CreatePost(userID, diaryID int64, options ...PostOption) *post.Post {
	opts := &PostOptions{
		Title:    dg.isolator.IsolateTitle("post"),
		Body:     "Test post body content",
		Summary:  "Test post summary",
		InfoType: 1,
	}

	for _, opt := range options {
		opt(opts)
	}

	now := time.Now().Unix()
	return &post.Post{
		Title:         opts.Title,
		Body:          opts.Body,
		Summary:       opts.Summary,
		DiaryID:       diaryID,
		AddUserID:     userID,
		AddUTCTime:    now,
		UpdateUTCTime: now,
		Status:        int8(1),
	}
}

func (dg *DataGeneratorImpl) GenerateUniqueString(prefix string) string {
	return dg.isolator.IsolateString(prefix)
}

func (dg *DataGeneratorImpl) GenerateUniqueEmail(prefix string) string {
	email := fmt.Sprintf("%s@test.com", prefix)
	return dg.isolator.IsolateEmail(email)
}

// CleanupManagerImpl 的所有方法实现都在 helpers/cleanup.go 中

// initGlobalConfig 初始化全局配置包
func initGlobalConfig() error {
	// 获取全局配置对象（会自动设置测试环境变量）
	globalConfig := GetGlobalConfig()
	
	// 验证配置
	if err := globalConfig.ValidateConfig(); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}
	
	// 设置全局配置对象，供其他模块使用
	config.SetConfig(globalConfig)
	
	// 初始化全局数据库连接
	if err := database.Connect(); err != nil {
		return fmt.Errorf("初始化数据库连接失败: %w", err)
	}
	
	return nil
}
