package utils_test

import (
	"testing"

	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/utils/i18n"
	framework "github.com/biuren/hebe/tests/framework"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	framework.RunTestsWithFlags(m)
}

// TestI18nMessageLoading 测试国际化消息加载功能
func TestI18nMessageLoading(t *testing.T) {
	// 创建测试配置
	cfg := &config.Config{
		App: config.AppConfig{
			LocalesDir:  "../../../website/locales",
			DefaultLang: "zh-cn",
		},
		SupportLangs: []config.SupportLang{
			{Code: "zh-cn", Label: "Chinese Simplified", LocalText: "简体中文"},
			{Code: "en", Label: "English", LocalText: "English"},
		},
	}

	// 创建国际化管理器
	manager := i18n.NewLocalizationManager(cfg)

	// 测试基本消息加载
	t.Run("基本消息加载", func(t *testing.T) {
		// 测试英文消息
		enMsg := manager.GetMessage("common.save", "en")
		assert.NotEqual(t, "common.save", enMsg, "英文消息应该被正确加载")

		// 测试中文消息
		zhMsg := manager.GetMessage("common.save", "zh-cn")
		assert.NotEqual(t, "common.save", zhMsg, "中文消息应该被正确加载")
	})

	// 测试 user.profile 相关消息
	t.Run("user.profile消息加载", func(t *testing.T) {
		// 测试 entries_count
		entriesCountMsg := manager.GetMessage("user.profile.entries_count", "zh-cn")
		assert.NotEqual(t, "user.profile.entries_count", entriesCountMsg, "entries_count 应该被正确加载")
		assert.Equal(t, "篇日记", entriesCountMsg, "entries_count 应该翻译为'篇日记'")

		// 测试其他 user.profile 键
		testCases := []struct {
			key      string
			expected string
		}{
			{"user.profile.introduction", "个人介绍"},
			{"user.profile.menu", "目录"},
			{"user.profile.latest_public_diary", "最新公开日记"},
			{"user.profile.joined_diaries", "参与的日记本"},
			{"user.profile.my_apps", "我的应用"},
			{"user.profile.app_list", "应用列表"},
			{"user.profile.diary", "日记本"},
		}

		for _, tc := range testCases {
			t.Run(tc.key, func(t *testing.T) {
				msg := manager.GetMessage(tc.key, "zh-cn")
				assert.NotEqual(t, tc.key, msg, "消息 %s 应该被正确加载", tc.key)
				assert.Equal(t, tc.expected, msg, "消息 %s 应该翻译为 '%s'", tc.key, tc.expected)
			})
		}
	})

	// 测试消息计数
	t.Run("消息计数验证", func(t *testing.T) {
		count := manager.GetMessageCount()
		assert.Greater(t, count, 300, "应该加载超过300条消息")
		t.Logf("总共加载了 %d 条消息", count)
	})

	// 测试语言回退机制
	t.Run("语言回退机制", func(t *testing.T) {
		// 测试不存在的语言，应该回退到默认语言
		msg := manager.GetMessage("common.save", "fr")
		assert.Equal(t, "保存", msg, "不存在的语言应该回退到默认语言")
	})

	// 测试不存在的键
	t.Run("不存在的键处理", func(t *testing.T) {
		msg := manager.GetMessage("non.existent.key", "zh-cn")
		assert.Equal(t, "non.existent.key", msg, "不存在的键应该返回原始键名")
	})
}

// TestI18nCacheKeyBuilding 测试缓存键构建
func TestI18nCacheKeyBuilding(t *testing.T) {
	t.Run("基本缓存键构建", func(t *testing.T) {
		cacheKey := i18n.BuildCacheKey("zh-cn", "user.profile.entries_count")
		expected := "zh-cn:user.profile.entries_count"
		assert.Equal(t, expected, cacheKey, "缓存键格式应该正确")
	})

	t.Run("大小写统一测试", func(t *testing.T) {
		// 测试大写key被转换为小写
		cacheKey1 := i18n.BuildCacheKey("zh-cn", "USER.LABEL.SETTINGS")
		cacheKey2 := i18n.BuildCacheKey("zh-cn", "user.label.settings")
		assert.Equal(t, cacheKey1, cacheKey2, "大写和小写key应该生成相同的缓存键")
		assert.Equal(t, "zh-cn:user.label.settings", cacheKey1, "key应该被转换为小写")

		// 测试locale大小写转换
		cacheKey3 := i18n.BuildCacheKey("ZH-CN", "user.label.settings")
		cacheKey4 := i18n.BuildCacheKey("zh-cn", "user.label.settings")
		assert.Equal(t, cacheKey3, cacheKey4, "大写和小写locale应该生成相同的缓存键")
		assert.Equal(t, "zh-cn:user.label.settings", cacheKey3, "locale应该被转换为小写")

		// 测试locale和key都是混合大小写
		cacheKey5 := i18n.BuildCacheKey("EN", "Common.Save")
		expected := "en:common.save"
		assert.Equal(t, expected, cacheKey5, "locale和key都应该被转换为小写")
	})
}

// TestI18nMessageRetrieval 测试消息检索功能
func TestI18nMessageRetrieval(t *testing.T) {
	// 设置测试消息
	testKey := "test:message"
	testValue := "测试消息"
	i18n.SetMessage(testKey, testValue)

	// 检索消息
	message, exists := i18n.GetMessage(testKey)
	assert.True(t, exists, "消息应该存在")
	assert.Equal(t, testValue, message, "消息内容应该正确")

	// 检索不存在的消息
	_, exists = i18n.GetMessage("non:existent")
	assert.False(t, exists, "不存在的消息应该返回false")
}
