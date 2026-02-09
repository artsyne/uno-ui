package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/biuren/hebe/internal/core/database"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/entry"
	"github.com/biuren/hebe/internal/domain/task"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/utils/i18n"
	core "github.com/biuren/hebe/tests/framework/core"
	"github.com/biuren/hebe/tests/framework/helpers"
	. "github.com/biuren/hebe/tests/framework/types"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// TestContext 测试上下文实现
type TestContext struct {
	testID       string
	config       TestConfigInterface
	database     DatabaseManagerInterface
	dataGenerator DataGeneratorInterface
	isolator     DataIsolatorInterface
	cleanup      CleanupManagerInterface
	services     ServiceRegistryInterface
	http         HTTPTestHelperInterface
	ctx          context.Context
	cancel       context.CancelFunc
}

// DatabaseManager 数据库管理器实现
// DatabaseManager 已移至 core 包中

// DataIsolator 数据隔离器实现
type DataIsolator struct {
	prefix      string
	testID      string
	createdData map[string][]interface{}
	mu          sync.RWMutex
}



// ServiceRegistry 服务注册表实现
type ServiceRegistry struct {
	userService  *user.Service
	diaryService *diary.Service
	entryService *entry.Service
	mu           sync.RWMutex
}

// HTTPTestHelper HTTP测试助手实现
type HTTPTestHelper struct {
	router *gin.Engine
	mu     sync.RWMutex
}

// DataGenerator 数据生成器实现
type DataGenerator struct {
	isolator DataIsolatorInterface
	mu       sync.RWMutex
}

// NewTestContext 创建新的测试上下文
func NewTestContext(options ...TestOption) TestContextInterface {
	// 重置配置单例状态，确保每次测试都能重新加载配置
	core.ResetGlobalConfig()
	core.ResetTestConfig()

	// 生成唯一的测试ID
	testID := generateTestID()
	
	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	
	// 获取测试配置
	testConfig := getDefaultTestConfig()
	
	// 创建数据隔离器
	isolator := &DataIsolator{
		prefix:      generatePrefix(testID),
		testID:      testID,
		createdData: make(map[string][]interface{}),
	}
	
	// 创建数据库管理器
	databaseManager := core.GetDatabaseManager()
	
	// 初始化数据库连接（如果失败，数据库管理器仍然可用，但连接为nil）
	if err := databaseManager.Initialize(); err != nil {
		// 数据库初始化失败时记录错误，但不阻止测试上下文创建
		fmt.Printf("Warning: Database initialization failed: %v\n", err)
	}

	// 创建测试上下文
	testCtx := &TestContext{
		testID:       testID,
		config:       testConfig,
		database:     databaseManager,
		isolator:     isolator,
		dataGenerator: newDataGenerator(isolator),
		ctx:          ctx,
		cancel:       cancel,
	}
	
	// 应用选项
	for _, option := range options {
		if optionFunc, ok := option.(func(*TestContext)); ok {
			optionFunc(testCtx)
		}
	}
	
	return testCtx
}

// WithDatabase 数据库选项
func WithDatabase() TestOption {
	return func(tc *TestContext) {
		// 首先初始化全局配置，确保config包的全局配置可用
		globalConfig := core.GetGlobalConfig()
		config.SetConfig(globalConfig)
		
		// 创建数据库管理器
		dbManager := core.GetDatabaseManager()
		
		// 初始化数据库连接
		if err := dbManager.Initialize(); err != nil {
			panic(fmt.Sprintf("Failed to initialize database: %v", err))
		}
		
		// 设置全局数据库连接，让 repository 构造函数可以使用
		if err := setGlobalDatabaseConnection(dbManager.GetDB(), tc.config); err != nil {
			panic(fmt.Sprintf("Failed to set global database connection: %v", err))
		}
		
		tc.database = dbManager
		
		// 创建清理管理器
		tc.cleanup = helpers.NewCleanupManager(dbManager.GetDB().DB, tc.config.(*TestConfig), tc.isolator)
		
		// 创建数据生成器
		tc.dataGenerator = newDataGenerator(tc.isolator)
		
		// 创建服务注册表
		serviceRegistry := &ServiceRegistry{}
		if err := serviceRegistry.InitializeServices(dbManager.GetDB(), tc.config); err != nil {
			panic(fmt.Sprintf("Failed to initialize services: %v", err))
		}
		tc.services = serviceRegistry
		
		// 创建HTTP测试助手
		tc.http = &HTTPTestHelper{}
	}
}

// TestContext接口实现
func (tc *TestContext) GetTestID() string {
	return tc.testID
}

func (tc *TestContext) GetConfig() TestConfigInterface {
	return tc.config
}

func (tc *TestContext) GetDatabase() DatabaseManagerInterface {
	return tc.database
}

func (tc *TestContext) GetDataGenerator() DataGeneratorInterface {
	return tc.dataGenerator
}

func (tc *TestContext) GetIsolator() DataIsolatorInterface {
	return tc.isolator
}

func (tc *TestContext) GetCleanup() CleanupManagerInterface {
	return tc.cleanup
}

func (tc *TestContext) GetServices() ServiceRegistryInterface {
	return tc.services
}

func (tc *TestContext) GetHTTP() HTTPTestHelperInterface {
	return tc.http
}

func (tc *TestContext) Context() context.Context {
	return tc.ctx
}

func (tc *TestContext) Cleanup() error {
	defer tc.cancel()
	if tc.cleanup != nil {
		return tc.cleanup.Cleanup("prefix")
	}
	return nil
}

func (tc *TestContext) ForceCleanupDatabase() error {
	if tc.cleanup != nil {
		return tc.cleanup.ForceCleanupDatabase()
	}
	return nil
}

// DatabaseManager 方法实现已移至 core 包中

// DataIsolator接口实现
func (di *DataIsolator) IsolateString(value string) string {
	return fmt.Sprintf("%s_%s", di.prefix, value)
}

func (di *DataIsolator) IsolateEmail(email string) string {
	return fmt.Sprintf("%s_%s", di.prefix, email)
}

func (di *DataIsolator) IsolateHandle(handle string) string {
	isolated := fmt.Sprintf("%s_%s", di.prefix, handle)
	if len(isolated) > 16 {
		return isolated[:16]
	}
	return isolated
}

func (di *DataIsolator) IsolateTitle(title string) string {
	isolated := fmt.Sprintf("%s_%s", di.prefix, title)
	if len(isolated) > 180 {
		return isolated[:180]
	}
	return isolated
}

func (di *DataIsolator) TrackCreated(table string, id interface{}) {
	di.mu.Lock()
	defer di.mu.Unlock()
	di.createdData[table] = append(di.createdData[table], id)
}

func (di *DataIsolator) GetCreatedData() map[string][]interface{} {
	di.mu.RLock()
	defer di.mu.RUnlock()
	result := make(map[string][]interface{})
	for k, v := range di.createdData {
		result[k] = make([]interface{}, len(v))
		copy(result[k], v)
	}
	return result
}

func (di *DataIsolator) GetPrefix() string {
	return di.prefix
}

func (di *DataIsolator) GetTestID() string {
	return di.testID
}

func (di *DataIsolator) GetDataCount() int {
	di.mu.RLock()
	defer di.mu.RUnlock()
	count := 0
	for _, ids := range di.createdData {
		count += len(ids)
	}
	return count
}

func (di *DataIsolator) GetTableCount() int {
	di.mu.RLock()
	defer di.mu.RUnlock()
	return len(di.createdData)
}

// ServiceRegistry接口实现
func (sr *ServiceRegistry) GetUserService() *user.Service {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	return sr.userService
}

func (sr *ServiceRegistry) GetDiaryService() *diary.Service {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	return sr.diaryService
}

func (sr *ServiceRegistry) GetEntryService() *entry.Service {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	return sr.entryService
}

func (sr *ServiceRegistry) InitializeServices(db *sqlx.DB, testConfig TestConfigInterface) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	
	// 获取全局配置而不是测试配置
	appConfig := config.GetConfig()
	if appConfig == nil {
		return fmt.Errorf("全局配置未初始化")
	}
	
	// 创建 repository 实例
	userRepo := repository.NewUserRepository()
	taskRepo := repository.NewTaskRepository()
	diaryRepo := repository.NewDiaryRepository()
	entryRepo := repository.NewEntryRepository()
	
	// 创建服务实例
	taskService := task.NewService(appConfig, taskRepo)
	// 创建一个简单的i18n管理器用于测试
	i18nManager := i18n.NewLocalizationManager(appConfig)
	sr.userService = user.NewService(appConfig, userRepo, taskService, i18nManager, nil)
	sr.diaryService = diary.NewService(appConfig, diaryRepo, sr.userService, taskService, i18nManager)
	sr.entryService = entry.NewService(appConfig, sr.userService, taskService, diaryRepo, entryRepo, i18nManager)
	
	return nil
}

// HTTPTestHelper接口实现
func (hth *HTTPTestHelper) SetupRouter() error {
	hth.mu.Lock()
	defer hth.mu.Unlock()
	hth.router = gin.New()
	return nil
}

// request 通用请求方法
func (hth *HTTPTestHelper) request(method, path string, body interface{}, headers map[string]string) *HTTPResponse {
	hth.mu.RLock()
	router := hth.router
	hth.mu.RUnlock()
	
	if router == nil {
		return &HTTPResponse{StatusCode: http.StatusInternalServerError}
	}
	
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			return &HTTPResponse{StatusCode: http.StatusInternalServerError}
		}
	}
	
	req := httptest.NewRequest(method, path, bytes.NewBuffer(reqBody))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	return &HTTPResponse{
		StatusCode: w.Code,
		Body:       w.Body.Bytes(),
		Headers:    map[string]string{"Content-Type": w.Header().Get("Content-Type")},
	}
}

func (hth *HTTPTestHelper) GET(path string, headers map[string]string) *HTTPResponse {
	return hth.request("GET", path, nil, headers)
}

func (hth *HTTPTestHelper) POST(path string, body interface{}, headers map[string]string) *HTTPResponse {
	return hth.request("POST", path, body, headers)
}

func (hth *HTTPTestHelper) PUT(path string, body interface{}, headers map[string]string) *HTTPResponse {
	return hth.request("PUT", path, body, headers)
}

func (hth *HTTPTestHelper) DELETE(path string, headers map[string]string) *HTTPResponse {
	return hth.request("DELETE", path, nil, headers)
}

// 辅助函数
func generateTestID() string {
	return fmt.Sprintf("test_%d_%d", time.Now().Unix(), rand.Intn(10000))
}

func generatePrefix(testID string) string {
	// Shorten prefix to fit in handle (max 16 chars)
	// t + 4 digits + 4 chars = 9 chars. 
	// handle = prefix + "_" + user (max 6) = 16 chars.
	return fmt.Sprintf("t%d%s", time.Now().Unix()%10000, testID[len(testID)-4:])
}

// getDefaultTestConfig 返回默认测试配置
func getDefaultTestConfig() TestConfigInterface {
	// 使用core包中的统一配置加载
	return core.GetTestConfig()
}

// newDataGenerator 创建数据生成器
func newDataGenerator(isolator DataIsolatorInterface) DataGeneratorInterface {
	return &DataGenerator{
		isolator: isolator,
	}
}

// CreateUser 创建测试用户
func (dg *DataGenerator) CreateUser(options ...UserOption) *user.User {
	dg.mu.Lock()
	defer dg.mu.Unlock()
	
	// 设置默认选项
	opts := &UserOptions{
		Handle:      dg.isolator.IsolateHandle("user"),
		Email:       dg.isolator.IsolateEmail("test@example.com"),
		Password:    "password123",
		DisplayName: "Test User",
		TimeZone:    "UTC",
		Lang:        "en",
	}
	
	// 应用用户提供的选项
	for _, option := range options {
		option(opts)
	}
	
	// 创建用户对象
	return &user.User{
		Handle:      opts.Handle,
		Email:       opts.Email,
		Pwd:         opts.Password,
		DisplayName: opts.DisplayName,
		TimeZone:    opts.TimeZone,
		Lang:        opts.Lang,
		Status:      3, // 正常状态
	}
}

// CreateDiary 创建测试日记
func (dg *DataGenerator) CreateDiary(userID int64, options ...DiaryOption) *diary.Diary {
	dg.mu.Lock()
	defer dg.mu.Unlock()
	
	// 设置默认选项
	opts := &DiaryOptions{
		Title:         dg.isolator.IsolateTitle("Test Diary"),
		Description:   "Test diary description",
		ViewSetting:   1,
		SocialSetting: 1,
		ThemeID:       1,
	}
	
	// 应用用户提供的选项
	for _, option := range options {
		option(opts)
	}
	
	// 创建日记对象
	return &diary.Diary{
		AddUserID:     userID,
		Title:         opts.Title,
		Remark:        opts.Description,
		ViewSetting:   int8(opts.ViewSetting),
		SocialSetting: int8(opts.SocialSetting),
		ThemeID:       int8(opts.ThemeID),
		Status:        1,
	}
}

// CreateEntry 创建测试文章
func (dg *DataGenerator) CreateEntry(userID, diaryID int64, options ...EntryOption) *entry.Entry {
	dg.mu.Lock()
	defer dg.mu.Unlock()
	
	// 设置默认选项
	opts := &EntryOptions{
		Title:    dg.isolator.IsolateTitle("Test Entry"),
		Body:     "Test entry content",
		Summary:  "Test entry summary",
		InfoType: 1,
	}
	
	// 应用用户提供的选项
	for _, option := range options {
		option(opts)
	}
	
	// 创建文章对象
	return &entry.Entry{
		AddUserID: userID,
		DiaryID:   diaryID,
		Title:     opts.Title,
		Body:      opts.Body,
		Summary:   opts.Summary,
		InfoType:  entry.InfoType(opts.InfoType),
		Status:    1,
	}
}

// GenerateUniqueString 生成唯一字符串
func (dg *DataGenerator) GenerateUniqueString(prefix string) string {
	return dg.isolator.IsolateString(prefix)
}

// GenerateUniqueEmail 生成唯一邮箱
func (dg *DataGenerator) GenerateUniqueEmail(prefix string) string {
	return dg.isolator.IsolateEmail(prefix + "@example.com")
}

// GenerateUniqueHandle 生成唯一Handle
func (dg *DataGenerator) GenerateUniqueHandle(prefix string) string {
	return dg.isolator.IsolateHandle(prefix)
}

// GenerateUniqueTitle 生成唯一标题
func (dg *DataGenerator) GenerateUniqueTitle(prefix string) string {
	return dg.isolator.IsolateTitle(prefix)
}

// GetTestFlags 获取测试标志 - 使用 core 包中的实现
func GetTestFlags() *TestFlags {
	return core.GetTestFlags()
}

// ShouldCleanupDatabase 检查是否应该清理数据库 - 使用 core 包中的实现
func ShouldCleanupDatabase() bool {
	return core.ShouldCleanupDatabase()
}

// RunTestsWithFlags 使用标志运行测试 - 使用 core 包中的实现
func RunTestsWithFlags(m *testing.M) {
	core.RunTestsWithFlags(m)
}

// WithCleanupEnabled 清理选项
func WithCleanupEnabled(enabled bool) TestOption {
	return func(tc *TestContext) {
		if config, ok := tc.config.(*TestConfig); ok {
			config.Cleanup.Enabled = enabled
		}
	}
}

// GetGlobalConfig 导出全局配置函数
func GetGlobalConfig() *config.Config {
	return core.GetGlobalConfig()
}

// setGlobalDatabaseConnection 设置全局数据库连接
// 这个函数用于在测试环境中设置全局数据库连接，让 repository 构造函数可以正常工作
func setGlobalDatabaseConnection(db *sqlx.DB, testConfig TestConfigInterface) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	
	// 构建 DSN 字符串，用于设置全局连接
	// 注意：这里我们直接使用已经建立的连接，而不是重新连接
	// 但是 database.SetDB 需要 DSN 字符串，所以我们需要从配置中构建
	dbConfig := testConfig.GetDatabaseConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
		dbConfig.Charset,
	)
	
	return database.SetDB(dsn)
}
