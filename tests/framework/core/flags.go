package framework

import (
	"os"
	"strconv"
	"strings"
	"sync"
)

// TestFlags 测试命令行参数
type TestFlags struct {
	ClearDB   bool // -cleardb 强制清理数据库
	Verbose   bool // -v 详细输出
	Parallel  bool // -parallel 并行执行
	Timeout   int  // -timeout 超时时间(秒)
}

var (
	testFlags     *TestFlags
	flagsInitOnce sync.Once
)

// GetTestFlags 获取测试命令行参数
func GetTestFlags() *TestFlags {
	flagsInitOnce.Do(func() {
		testFlags = parseTestFlags()
	})
	return testFlags
}

// parseTestFlags 解析测试命令行参数
func parseTestFlags() *TestFlags {
	flags := &TestFlags{
		Parallel: true, // 默认值
		Timeout:  30,   // 默认值
	}
	
	// 直接解析os.Args，避免与Go测试框架的flag包冲突
	for _, arg := range os.Args {
		switch {
		case arg == "-cleardb":
			flags.ClearDB = true
		case arg == "-verbose":
			flags.Verbose = true
		case arg == "-parallel=false":
			flags.Parallel = false
		case strings.HasPrefix(arg, "-timeout="):
			if timeoutStr := strings.TrimPrefix(arg, "-timeout="); timeoutStr != "" {
				if timeout, err := strconv.Atoi(timeoutStr); err == nil && timeout > 0 {
					flags.Timeout = timeout
				}
			}
		}
	}
	
	return flags
}



// ShouldCleanupDatabase 检查是否应该清理数据库
// 优先级: 命令行参数 > 环境变量 > 默认配置
func ShouldCleanupDatabase() bool {
	flags := GetTestFlags()
	
	// 如果命令行指定了 -cleardb，则强制清理
	if flags.ClearDB {
		return true
	}
	
	// 否则检查环境变量
	if envVal := os.Getenv("IDIARY_TEST_CLEANUP_ENABLED"); envVal != "" {
		return envVal == "true"
	}
	
	// 默认不清理
	return false
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