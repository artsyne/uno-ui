package app

import (
	"html/template"
	"strings"
	"testing"

	"github.com/biuren/hebe/internal/app/templates"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/utils/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTemplateLocalizationFunctions 测试模板本地化函数
func TestTemplateLocalizationFunctions(t *testing.T) {
	// 创建测试配置
	cfg := &config.Config{
		App: config.AppConfig{
			LocalesDir:  "../../../website/locales",
			DefaultLang: "zh-cn",
		},
		SupportLangs: []config.SupportLang{
			{Code: "zh-cn", Label: "Chinese Simplified", LocalText: "简体中文"},
			{Code: "en-us", Label: "English", LocalText: "English"},
		},
	}

	// 创建国际化管理器并设置为全局管理器
	manager := i18n.NewLocalizationManager(cfg)
	templates.SetGlobalManager(manager)

	// 获取本地化函数
	funcMap := templates.GetLocalizationFuncs()
	require.NotNil(t, funcMap, "函数映射不应该为空")

	// 测试 getMessage 函数
	t.Run("getMessage函数", func(t *testing.T) {
		getMessageFunc, exists := funcMap["getMessage"]
		require.True(t, exists, "getMessage 函数应该存在")

		// 测试存在的消息
		result := getMessageFunc.(func(string, string) string)("common.save", "zh-cn")
		assert.Equal(t, "保存", result, "应该返回正确的中文翻译")

		result = getMessageFunc.(func(string, string) string)("common.save", "en-us")
		assert.Equal(t, "Save", result, "应该返回正确的英文翻译")

		// 测试不存在的消息
		result = getMessageFunc.(func(string, string) string)("non.existent.key", "zh-cn")
		assert.Equal(t, "non.existent.key", result, "不存在的键应该返回原始键名")
	})

	// 测试 getMessageWithDefault 函数
	t.Run("getMessageWithDefault函数", func(t *testing.T) {
		getMessageWithDefaultFunc, exists := funcMap["getMessageWithDefault"]
		require.True(t, exists, "getMessageWithDefault 函数应该存在")

		// 测试存在的消息
		result := getMessageWithDefaultFunc.(func(string, string, string) string)("common.save", "zh-cn", "默认值")
		assert.Equal(t, "保存", result, "存在的消息应该返回翻译")

		// 测试不存在的消息，有默认值
		result = getMessageWithDefaultFunc.(func(string, string, string) string)("non.existent.key", "zh-cn", "默认值")
		assert.Equal(t, "默认值", result, "不存在的消息应该返回默认值")

		// 测试不存在的消息，无默认值
		result = getMessageWithDefaultFunc.(func(string, string, string) string)("non.existent.key", "zh-cn", "")
		assert.Equal(t, "non.existent.key", result, "无默认值时应该返回原始键名")
	})

	// 测试 getMessages 函数
	t.Run("getMessages函数", func(t *testing.T) {
		getMessagesFunc, exists := funcMap["getMessages"]
		require.True(t, exists, "getMessages 函数应该存在")

		keys := []string{"common.save", "common.edit", "non.existent.key"}
		result := getMessagesFunc.(func([]string, string) map[string]string)(keys, "zh-cn")

		assert.Len(t, result, 3, "应该返回3个消息")
		assert.Equal(t, "保存", result["common.save"], "应该包含正确的翻译")
		assert.Equal(t, "编辑", result["common.edit"], "应该包含正确的翻译")
		assert.Equal(t, "non.existent.key", result["non.existent.key"], "不存在的键应该返回原始键名")
	})

	// 测试 getSupportedLanguages 函数
	t.Run("getSupportedLanguages函数", func(t *testing.T) {
		getSupportedLanguagesFunc, exists := funcMap["getSupportedLanguages"]
		require.True(t, exists, "getSupportedLanguages 函数应该存在")

		result := getSupportedLanguagesFunc.(func() []config.SupportLang)()
		assert.Len(t, result, 2, "应该返回2种支持的语言")
		assert.Equal(t, "zh-cn", result[0].Code, "第一种语言应该是中文")
		assert.Equal(t, "en-us", result[1].Code, "第二种语言应该是英文")
	})

	// 测试 isLanguageSupported 函数
	t.Run("isLanguageSupported函数", func(t *testing.T) {
		isLanguageSupportedFunc, exists := funcMap["isLanguageSupported"]
		require.True(t, exists, "isLanguageSupported 函数应该存在")

		result := isLanguageSupportedFunc.(func(string) bool)("zh-cn")
		assert.True(t, result, "中文应该被支持")

		result = isLanguageSupportedFunc.(func(string) bool)("en-us")
		assert.True(t, result, "英文应该被支持")

		result = isLanguageSupportedFunc.(func(string) bool)("en-us")
		assert.True(t, result, "英文（规范化后）应该被支持")

		result = isLanguageSupportedFunc.(func(string) bool)("fr")
		assert.False(t, result, "法文不应该被支持")
	})
}

// TestTemplateIntegration 测试模板集成
func TestTemplateIntegration(t *testing.T) {
	// 创建测试配置
	cfg := &config.Config{
		App: config.AppConfig{
			LocalesDir:  "../../../website/locales",
			DefaultLang: "zh-cn",
		},
		SupportLangs: []config.SupportLang{
			{Code: "zh-cn", Label: "Chinese Simplified", LocalText: "简体中文"},
			{Code: "en-us", Label: "English", LocalText: "English"},
		},
	}

	// 创建国际化管理器并设置为全局管理器
	manager := i18n.NewLocalizationManager(cfg)
	templates.SetGlobalManager(manager)

	// 测试在实际模板中使用本地化函数
	t.Run("模板中使用本地化函数", func(t *testing.T) {
		// 创建模板
		tmplText := `
{{- $locale := "zh-cn" -}}
保存按钮: {{getMessage "common.save" $locale}}
编辑按钮: {{getMessage "common.edit" $locale}}
带默认值: {{getMessageWithDefault "non.existent" $locale "默认文本"}}
支持的语言: {{range getSupportedLanguages}}{{.LocalText}} {{end}}
是否支持中文: {{isLanguageSupported "zh-cn"}}
`

		// 创建模板并添加本地化函数
		tmpl := template.New("test")
		tmpl.Funcs(templates.GetLocalizationFuncs())
		tmpl, err := tmpl.Parse(tmplText)
		require.NoError(t, err, "模板解析不应该出错")

		// 执行模板
		var buf strings.Builder
		err = tmpl.Execute(&buf, nil)
		require.NoError(t, err, "模板执行不应该出错")

		result := buf.String()
		t.Logf("模板输出:\n%s", result)

		// 验证输出
		assert.Contains(t, result, "保存按钮: 保存", "应该包含正确的保存翻译")
		assert.Contains(t, result, "编辑按钮: 编辑", "应该包含正确的编辑翻译")
		assert.Contains(t, result, "带默认值: 默认文本", "应该包含默认值")
		assert.Contains(t, result, "简体中文", "应该包含中文语言名")
		assert.Contains(t, result, "English", "应该包含英文语言名")
		assert.Contains(t, result, "是否支持中文: true", "应该显示支持中文")
	})

	// 测试多语言模板
	t.Run("多语言模板", func(t *testing.T) {
		tmplText := `
{{- $locale := .Locale -}}
{{getMessage "common.home" $locale}} | {{getMessage "common.settings" $locale}}
`

		tmpl := template.New("multilang")
		tmpl.Funcs(templates.GetLocalizationFuncs())
		tmpl, err := tmpl.Parse(tmplText)
		require.NoError(t, err, "模板解析不应该出错")

		// 测试中文
		var zhBuf strings.Builder
		err = tmpl.Execute(&zhBuf, map[string]string{"Locale": "zh-cn"})
		require.NoError(t, err, "中文模板执行不应该出错")
		zhResult := strings.TrimSpace(zhBuf.String())
		assert.Equal(t, "首页 | 设置", zhResult, "中文输出应该正确")

		// 测试英文
		var enBuf strings.Builder
		err = tmpl.Execute(&enBuf, map[string]string{"Locale": "en-us"})
		require.NoError(t, err, "英文模板执行不应该出错")
		enResult := strings.TrimSpace(enBuf.String())
		assert.Equal(t, "Home | Settings", enResult, "英文输出应该正确")
	})
}

// TestTemplateWithoutGlobalManager 测试没有全局管理器时的行为
func TestTemplateWithoutGlobalManager(t *testing.T) {
	// 清除全局管理器
	templates.SetGlobalManager(nil)

	// 获取本地化函数
	funcMap := templates.GetLocalizationFuncs()

	// 测试 getSupportedLanguages 在没有管理器时的行为
	t.Run("无管理器时getSupportedLanguages", func(t *testing.T) {
		getSupportedLanguagesFunc := funcMap["getSupportedLanguages"]
		result := getSupportedLanguagesFunc.(func() []config.SupportLang)()
		assert.Nil(t, result, "没有管理器时应该返回nil")
	})

	// 测试 isLanguageSupported 在没有管理器时的行为
	t.Run("无管理器时isLanguageSupported", func(t *testing.T) {
		isLanguageSupportedFunc := funcMap["isLanguageSupported"]
		result := isLanguageSupportedFunc.(func(string) bool)("zh-cn")
		assert.False(t, result, "没有管理器时应该返回false")
	})
}

// TestTemplateErrorHandling 测试模板函数的错误处理
func TestTemplateErrorHandling(t *testing.T) {
	// 创建测试配置
	cfg := &config.Config{
		App: config.AppConfig{
			LocalesDir:  "../../../website/locales",
			DefaultLang: "zh-cn",
		},
		SupportLangs: []config.SupportLang{
			{Code: "zh-cn", Label: "Chinese Simplified", LocalText: "简体中文"},
			{Code: "en-us", Label: "English", LocalText: "English"},
		},
	}

	manager := i18n.NewLocalizationManager(cfg)
	templates.SetGlobalManager(manager)
	funcMap := templates.GetLocalizationFuncs()

	// 测试空键处理
	t.Run("空键处理", func(t *testing.T) {
		getMessageFunc := funcMap["getMessage"]
		result := getMessageFunc.(func(string, string) string)("", "zh-cn")
		assert.Equal(t, "", result, "空键应该返回空字符串")
	})

	// 测试空语言处理
	t.Run("空语言处理", func(t *testing.T) {
		getMessageFunc := funcMap["getMessage"]
		result := getMessageFunc.(func(string, string) string)("common.save", "")
		// 应该返回键名或者有合理的降级处理
		assert.NotEmpty(t, result, "空语言时应该有降级处理")
	})

	// 测试空键列表处理
	t.Run("空键列表处理", func(t *testing.T) {
		getMessagesFunc := funcMap["getMessages"]
		result := getMessagesFunc.(func([]string, string) map[string]string)([]string{}, "zh-cn")
		assert.Empty(t, result, "空键列表应该返回空映射")

		// 测试nil键列表
		result = getMessagesFunc.(func([]string, string) map[string]string)(nil, "zh-cn")
		assert.Empty(t, result, "nil键列表应该返回空映射")
	})
}

// TestTemplatePerformance 测试模板函数性能
func TestTemplatePerformance(t *testing.T) {
	// 创建测试配置
	cfg := &config.Config{
		App: config.AppConfig{
			LocalesDir:  "../../../website/locales",
			DefaultLang: "zh-cn",
		},
		SupportLangs: []config.SupportLang{
			{Code: "zh-cn", Label: "Chinese Simplified", LocalText: "简体中文"},
			{Code: "en-us", Label: "English", LocalText: "English"},
		},
	}

	manager := i18n.NewLocalizationManager(cfg)
	templates.SetGlobalManager(manager)
	funcMap := templates.GetLocalizationFuncs()

	// 性能测试：大量消息获取
	t.Run("大量消息获取性能", func(t *testing.T) {
		getMessageFunc := funcMap["getMessage"]
		
		// 预热
		for i := 0; i < 100; i++ {
			getMessageFunc.(func(string, string) string)("common.save", "zh-cn")
		}

		// 性能测试
		for i := 0; i < 1000; i++ {
			result := getMessageFunc.(func(string, string) string)("common.save", "zh-cn")
			assert.Equal(t, "保存", result, "性能测试中消息应该正确")
		}
	})

	// 性能测试：批量消息获取
	t.Run("批量消息获取性能", func(t *testing.T) {
		getMessagesFunc := funcMap["getMessages"]
		keys := []string{"common.save", "common.edit", "common.delete", "common.home", "common.settings"}

		// 性能测试
		for i := 0; i < 200; i++ {
			result := getMessagesFunc.(func([]string, string) map[string]string)(keys, "zh-cn")
			assert.Len(t, result, 5, "批量获取应该返回正确数量的消息")
		}
	})
}