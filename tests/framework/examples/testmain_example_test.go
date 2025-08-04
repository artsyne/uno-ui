package examples

import (
	"testing"

	core "github.com/biuren/hebe/tests/framework/core"
	"github.com/stretchr/testify/assert"
)

// TestMain 演示使用统一入口支持标准 flag 语法
// 使用方式：
//   go test ./tests/framework/examples/ -v -cleandb
//   go test ./tests/framework/examples/ -v -cleandb -timeout=60
func TestMain(m *testing.M) {
	core.RunTestsWithFlags(m)
}

// TestStandardFlagSyntax 演示标准 flag 语法的使用
func TestStandardFlagSyntax(t *testing.T) {
	t.Parallel()

	// 检查命令行参数状态
	flags := core.GetTestFlags()
	t.Logf("标准flag语法 - 命令行参数状态: cleandb=%v, verbose=%v, parallel=%v, timeout=%d",
		flags.ClearDB, flags.Verbose, flags.Parallel, flags.Timeout)

	// 检查清理配置
	shouldCleanup := core.ShouldCleanupDatabase()
	t.Logf("是否清理数据库: %v", shouldCleanup)

	ctx := core.NewTestContext("standard_flag_syntax_test")
	defer ctx.Cleanup()

	// 创建测试数据
	generator := ctx.GetDataGenerator()
	user := generator.CreateUser()

	// 验证数据创建成功
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.UserName)
	t.Logf("创建测试用户: %s", user.UserName)

	// 如果启用了清理，验证清理功能
	if shouldCleanup {
		t.Log("测试结束后将执行数据库清理")
	}
}

// TestParameterConfiguration 演示参数配置功能
func TestParameterConfiguration(t *testing.T) {
	t.Parallel()

	flags := core.GetTestFlags()
	t.Logf("参数配置 - cleandb=%v, verbose=%v, parallel=%v, timeout=%d",
		flags.ClearDB, flags.Verbose, flags.Parallel, flags.Timeout)

	// 验证配置状态
	if flags.ClearDB {
		t.Log("✅ 清理数据库已启用")
	} else {
		t.Log("ℹ️ 清理数据库未启用，使用默认配置")
	}

	if flags.Verbose {
		t.Log("✅ 详细模式已启用")
	}

	if !flags.Parallel {
		t.Log("⚠️ 并行测试已禁用")
	}

	if flags.Timeout != 30 {
		t.Logf("⏱️ 自定义超时时间: %d秒", flags.Timeout)
	}

	ctx := core.NewTestContext("parameter_configuration_test")
	defer ctx.Cleanup()

	// 简单的功能验证
	generator := ctx.GetDataGenerator()
	user := generator.CreateUser()
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.UserName)
	t.Logf("创建测试用户: %s", user.UserName)
}