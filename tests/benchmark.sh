#!/bin/bash

# 性能基准测试脚本
# 用法: ./scripts/benchmark.sh [选项]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 获取脚本所在目录
SCRIPT_DIR=$(dirname "$0")
PROJECT_DIR=$(cd "$SCRIPT_DIR/../.." && pwd)

echo -e "${BLUE}🚀 Hebe 项目性能基准测试${NC}"
echo "项目目录: $PROJECT_DIR"
echo ""

# 切换到项目目录
cd "$PROJECT_DIR"

# 检查 Go 环境
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go 未安装或不在 PATH 中${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Go 版本: $(go version)${NC}"
echo ""

# 运行基准测试的函数
run_benchmark() {
    local test_name="$1"
    local bench_pattern="$2"
    local description="$3"
    
    echo -e "${YELLOW}📊 运行 $description${NC}"
    echo "测试模式: $bench_pattern"
    echo ""
    
    # 创建日志目录和文件名
    LOG_DIR="logs/benchmark"
    mkdir -p "$LOG_DIR"
    LOG_FILE="$LOG_DIR/benchmark_${test_name}_$(date +%Y%m%d_%H%M%S).log"
    
    # 运行基准测试
    go test -bench="$bench_pattern" -benchmem -count=3 ./tests/benchmark/... | tee "$LOG_FILE"
    
    echo ""
    echo -e "${GREEN}✅ $description 完成${NC}"
    echo "----------------------------------------"
    echo ""
}

# 解析命令行参数
case "${1:-all}" in
    "user")
        run_benchmark "用户相关性能测试" "BenchmarkUser.*" "user"
        ;;
    "diary")
        run_benchmark "日记相关性能测试" "BenchmarkDiary.*" "diary"
        ;;
    "post")
        run_benchmark "文章相关性能测试" "BenchmarkPost.*" "post"
        ;;
    "database"|"db")
        run_benchmark "数据库操作性能测试" "BenchmarkDatabase.*" "database"
        ;;
    "concurrent")
        run_benchmark "并发性能测试" "BenchmarkConcurrent.*" "concurrent"
        ;;
    "all")
        echo -e "${BLUE}🔄 运行所有性能基准测试${NC}"
        echo ""
        
        run_benchmark "用户注册性能测试" "BenchmarkUserRegistration" "user"
        run_benchmark "日记本创建性能测试" "BenchmarkDiaryCreation" "diary"
        run_benchmark "文章创建性能测试" "BenchmarkPostCreation" "post"
        run_benchmark "数据库操作性能测试" "BenchmarkDatabaseOperations" "database"
        run_benchmark "并发用户注册性能测试" "BenchmarkConcurrentUserRegistration" "concurrent"
        
        echo -e "${GREEN}🎉 所有性能基准测试完成！${NC}"
        echo ""
        echo "📈 性能报告已保存到 logs/benchmark/ 目录中"
        ;;
    "help")
        echo "用法: $0 [选项]"
        echo ""
        echo "选项:"
        echo "  all        运行所有基准测试 (默认)"
        echo "  user       运行用户相关性能测试"
        echo "  diary      运行日记本相关性能测试"
        echo "  post       运行文章相关性能测试"
        echo "  database   运行数据库操作性能测试"
        echo "  concurrent 运行并发性能测试"
        echo "  help       显示此帮助信息"
        echo ""
        echo "示例:"
        echo "  $0              # 运行所有测试"
        echo "  $0 user         # 只运行用户相关测试"
        echo "  $0 concurrent   # 只运行并发测试"
        ;;
    *)
        echo -e "${RED}❌ 未知选项: $1${NC}"
        echo "使用 '$0 help' 查看可用选项"
        exit 1
        ;;
esac

echo -e "${BLUE}📋 基准测试提示:${NC}"
echo "• 基准测试结果会受到系统负载影响"
echo "• 建议在系统空闲时运行以获得准确结果"
echo "• 可以多次运行取平均值来提高准确性"
echo "• 日志文件包含详细的性能数据和内存使用情况"
echo ""