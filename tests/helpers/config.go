package testhelpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/biuren/hebe/internal/utils/config"
)

// ConfigTestHelper 配置测试辅助工具
type ConfigTestHelper struct {
	originalConfig *config.Config
}

// SetupTestConfig 使用.env.local文件初始化测试环境
// 直接使用项目根目录的.env.local配置，但修改数据库名为测试数据库
func SetupTestConfig(t *testing.T) *ConfigTestHelper {
	t.Helper()

	// 使用统一的环境变量设置
	SetupTestEnvironment()

	// 确保必要的目录存在
	ensureTestDirectories(t)

	// 备份原始配置
	originalConfig := config.GetConfig()

	// 初始化配置（会自动读取.env.local）
	if err := config.InitConfig(); err != nil {
		t.Fatalf("初始化配置失败: %v", err)
	}

	// 获取配置并修改数据库名为测试数据库
	cfg := config.GetConfig()
	if cfg != nil {
		// 修改数据库名为测试数据库
		cfg.Database.Name = "hebe_test"
	}

	return &ConfigTestHelper{
		originalConfig: originalConfig,
	}
}

// ensureTestDirectories 确保测试所需的目录存在
func ensureTestDirectories(t *testing.T) {
	t.Helper()

	projectRoot := getProjectRoot()
	directories := []string{
		filepath.Join(projectRoot, "website/locales"),
		filepath.Join(projectRoot, "website/static"),
		filepath.Join(projectRoot, "tests/test_uploads"),
		filepath.Join(projectRoot, "tests/test_logs"),
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Logf("Warning: Failed to create test directory %s: %v", dir, err)
		}
	}

	// 创建基本的本地化文件（如果不存在）
	createBasicLocalizationFiles(t)
}

// createBasicLocalizationFiles 创建基本的本地化文件（如果不存在）
func createBasicLocalizationFiles(t *testing.T) {
	t.Helper()

	projectRoot := getProjectRoot()

	// 检查并创建中文本地化文件
	zhCNFile := filepath.Join(projectRoot, "website/locales/zh-cn.json")
	if _, err := os.Stat(zhCNFile); os.IsNotExist(err) {
		zhCNContent := `{
	"nav": {
		"home": "首页"
	},
	"common": {
		"default_page_desc": "默认页面描述"
	}
}`
		if err := os.WriteFile(zhCNFile, []byte(zhCNContent), 0644); err != nil {
			t.Logf("Warning: Failed to create %s: %v", zhCNFile, err)
		}
	}

	// 检查并创建英文本地化文件
	enFile := filepath.Join(projectRoot, "website/locales/en.json")
	if _, err := os.Stat(enFile); os.IsNotExist(err) {
		enContent := `{
	"nav": {
		"home": "Home"
	},
	"common": {
		"default_page_desc": "Default page description"
	}
}`
		if err := os.WriteFile(enFile, []byte(enContent), 0644); err != nil {
			t.Logf("Warning: Failed to create %s: %v", enFile, err)
		}
	}
}

// Cleanup 清理测试配置，恢复原始状态
func (h *ConfigTestHelper) Cleanup(t *testing.T) {
	t.Helper()
	if h.originalConfig != nil {
		config.SetConfig(h.originalConfig)
	}
}

// GetTestConfig 获取当前测试配置
func (h *ConfigTestHelper) GetTestConfig() *config.Config {
	return config.GetConfig()
}
