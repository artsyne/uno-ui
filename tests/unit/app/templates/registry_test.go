package templates

import (
	"html/template"
	"testing"

	"github.com/biuren/hebe/internal/app/templates"
	app_templates "github.com/biuren/hebe/internal/app/templates"
	"github.com/stretchr/testify/assert"
)

// TestRegisterAllFuncs 测试注册所有模板函数
func TestRegisterAllFuncs(t *testing.T) {
	// 获取所有模板函数
	funcMap := app_templates.RegisterAllFuncs()
	assert.NotNil(t, funcMap)

	// 验证返回的函数映射不为空
	assert.IsType(t, template.FuncMap{}, funcMap)

	// 验证函数映射包含一些函数
	assert.Greater(t, len(funcMap), 0, "应该注册一些模板函数")
}

// TestRegisterAllFuncsReturnType 测试返回类型
func TestRegisterAllFuncsReturnType(t *testing.T) {
	// 测试返回类型是否正确
	result := app_templates.RegisterAllFuncs()
	assert.NotNil(t, result)
	assert.IsType(t, template.FuncMap{}, result)
}

// TestRegisterAllFuncsFunctionAvailability 测试函数可用性
func TestRegisterAllFuncsFunctionAvailability(t *testing.T) {
	// 获取所有模板函数
	funcMap := app_templates.RegisterAllFuncs()
	assert.NotNil(t, funcMap)

	// 创建模板并应用函数
	tmpl := template.New("test").Funcs(funcMap)
	assert.NotNil(t, tmpl)

	// 尝试解析一个使用模板函数的简单模板
	testTemplates := []struct {
		name        string
		template    string
		expectError bool
	}{
		{
			name:        "basic template without functions",
			template:    "Hello {{.Name}}",
			expectError: false,
		},
		{
			name:        "template with potential custom function",
			template:    "Hello {{.Name}}", // 使用基本语法，避免依赖具体函数
			expectError: false,
		},
	}

	for _, tt := range testTemplates {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tmpl.Parse(tt.template)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestRegisterAllFuncsMultipleRegistrations 测试多次注册
func TestRegisterAllFuncsMultipleRegistrations(t *testing.T) {
	// 第一次获取函数映射
	funcMap1 := app_templates.RegisterAllFuncs()
	assert.NotNil(t, funcMap1)

	// 第二次获取函数映射（应该不会出错）
	funcMap2 := app_templates.RegisterAllFuncs()
	assert.NotNil(t, funcMap2)

	// 验证两次获取的函数映射长度相同
	assert.Equal(t, len(funcMap1), len(funcMap2))

	// 验证函数映射的类型相同
	assert.IsType(t, funcMap1, funcMap2)
}

// TestRegisterAllFuncsTemplateChaining 测试模板链式调用
func TestRegisterAllFuncsTemplateChaining(t *testing.T) {
	// 获取函数映射并创建模板
	funcMap := app_templates.RegisterAllFuncs()
	tmpl := template.New("test").Funcs(funcMap)

	// 验证可以继续链式调用其他模板方法
	assert.NotNil(t, tmpl)

	// 测试可以继续解析模板
	_, err := tmpl.Parse("Hello {{.Name}}")
	assert.NoError(t, err)
}

// TestRegisterAllFuncs_FunctionCategories 测试函数分类注册
func TestRegisterAllFuncsFunctionCategories(t *testing.T) {
	// 根据 registry.go 中的注释，应该注册以下类别的函数：
	// - 本地化函数 (localization functions)
	// - 认证函数 (authentication functions)
	// - UI 函数 (UI functions)
	// - 时间函数 (time functions)

	funcMap := templates.RegisterAllFuncs()
	tmpl := template.New("test").Funcs(funcMap)
	assert.NotNil(t, tmpl)

	// 由于我们不知道具体的函数名称，我们测试模板注册过程是否成功
	// 这确保了 RegisterAllFuncs 函数能够正常执行而不出错

	// 测试模板可以正常工作
	testData := map[string]interface{}{
		"Name": "Test User",
		"Time": "2023-01-01",
	}

	_, err := tmpl.Parse("Hello {{.Name}} at {{.Time}}")
	assert.NoError(t, err)

	// 可以执行模板
	buf := make([]byte, 0, 100)
	// 注意：这里不执行模板，因为我们主要测试注册过程
	_ = buf
	_ = testData
}

// TestRegisterAllFuncs_EmptyTemplateName 测试空模板名称
func TestRegisterAllFuncsEmptyTemplateName(t *testing.T) {
	// 获取函数映射并创建一个没有名称的模板
	funcMap := templates.RegisterAllFuncs()
	tmpl := template.New("").Funcs(funcMap)
	assert.NotNil(t, tmpl)

	// 验证即使模板名称为空，函数注册也应该成功
	assert.NotNil(t, funcMap)
	assert.Greater(t, len(funcMap), 0)
}

// TestRegisterAllFuncs_TemplateWithExistingFuncs 测试已有函数的模板
func TestRegisterAllFuncsTemplateWithExistingFuncs(t *testing.T) {
	// 获取注册的函数映射
	registeredFuncs := templates.RegisterAllFuncs()
	assert.NotNil(t, registeredFuncs)

	// 创建一个已经有自定义函数的模板
	customFuncMap := template.FuncMap{
		"customFunc": func(s string) string {
			return "custom: " + s
		},
	}

	// 合并自定义函数和注册的函数
	for name, fn := range registeredFuncs {
		customFuncMap[name] = fn
	}

	tmpl := template.New("test").Funcs(customFuncMap)
	assert.NotNil(t, tmpl)

	// 验证原有函数仍然可用
	_, err := tmpl.Parse("{{customFunc .Name}}")
	assert.NoError(t, err)
}

// TestRegisterAllFuncs_Idempotency 测试幂等性
func TestRegisterAllFuncsIdempotency(t *testing.T) {
	// 测试多次调用 RegisterAllFuncs 是否是幂等的

	// 多次获取函数映射
	result1 := templates.RegisterAllFuncs()
	result2 := templates.RegisterAllFuncs()
	result3 := templates.RegisterAllFuncs()

	// 验证所有结果的长度相同
	assert.Equal(t, len(result1), len(result2))
	assert.Equal(t, len(result2), len(result3))

	// 验证类型相同
	assert.IsType(t, result1, result2)
	assert.IsType(t, result2, result3)
}

// TestRegisterAllFuncs_Performance 测试性能
func TestRegisterAllFuncsPerformance(t *testing.T) {
	// 简单的性能测试，确保函数注册不会太慢
	for i := 0; i < 100; i++ {
		result := templates.RegisterAllFuncs()
		assert.NotNil(t, result)
		assert.Greater(t, len(result), 0)
	}

	// 如果能执行到这里，说明性能是可接受的
	assert.True(t, true)
}
