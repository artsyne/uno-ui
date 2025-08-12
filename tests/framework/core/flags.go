package framework

import (
	"flag"
	"os"
	"sync"
	"testing"

	. "github.com/biuren/hebe/tests/framework/types"
)

// TestFlags 类型已在 types/structs.go 中定义

var (
	testFlags     *TestFlags
	flagsInitOnce sync.Once
	
	// 标准 flag 定义（用于支持 TestMain）
	cleanDBFlag  *bool
	verboseFlag  *bool
	parallelFlag *bool
	timeoutFlag  *int
)

func init() {
	// 定义 flags，但不立即解析
	cleanDBFlag = flag.Bool("cleandb", false, "强制清理测试数据库")
	verboseFlag = flag.Bool("verbose", false, "详细输出模式")
	parallelFlag = flag.Bool("parallel", true, "启用并行测试")
	timeoutFlag = flag.Int("timeout", 30, "测试超时时间(秒)")
}

// GetTestFlags 获取测试命令行参数
func GetTestFlags() *TestFlags {
	flagsInitOnce.Do(func() {
		testFlags = parseTestFlags()
	})
	return testFlags
}

// parseTestFlags 解析测试命令行参数
// 支持多种方式：标准flag包、环境变量、直接解析os.Args
func parseTestFlags() *TestFlags {
	flags := &TestFlags{
		Parallel: true, // 默认值
		Timeout:  30,   // 默认值
	}
	
	// 使用标准 flag 包解析参数
	if flag.Parsed() {
		flags.ClearDB = *cleanDBFlag
		flags.Verbose = *verboseFlag
		flags.Parallel = *parallelFlag
		flags.Timeout = *timeoutFlag
	}
	
	return flags
}



// ShouldCleanupDatabase 检查是否应该清理数据库
// 优先级: 命令行参数 > 环境变量 > 默认配置
func ShouldCleanupDatabase() bool {
	flags := GetTestFlags()
	
	// 直接返回解析后的标志，因为parseTestFlags已经处理了所有优先级
	return flags.ClearDB
}

// IsVerboseMode 检查是否为详细模式
func IsVerboseMode() bool {
	return GetTestFlags().Verbose
}

// ShouldRunParallel 检查是否应该并行执行
func ShouldRunParallel() bool {
	return GetTestFlags().Parallel
}

// GetTestTimeout 获取测试超时时间
func GetTestTimeout() int {
	return GetTestFlags().Timeout
}

// RunTestsWithFlags 为支持 TestMain 的包提供统一入口
// 使用方式：
//   func TestMain(m *testing.M) {
//       framework.RunTestsWithFlags(m)
//   }
func RunTestsWithFlags(m *testing.M) {
	if !flag.Parsed() {
		flag.Parse()
	}
	
	// 强制重新解析 flags
	flagsInitOnce = sync.Once{}
	testFlags = nil
	
	code := m.Run()
	
	// 全局清理
	if GetTestFlags().ClearDB {
		performGlobalCleanup()
	}
	
	os.Exit(code)
}

// performGlobalCleanup 执行全局清理操作
func performGlobalCleanup() {
	// 这里可以添加全局清理逻辑
	// 例如：清理所有测试数据库、临时文件等
	if IsVerboseMode() {
		println("执行全局清理操作...")
	}
}