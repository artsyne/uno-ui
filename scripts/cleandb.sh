#!/bin/bash

# 测试数据库清理脚本
# 用法: ./scripts/cleandb.sh [选项]

set -e

# 获取脚本所在目录的父目录（项目根目录）
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CLEANDB_CMD="$PROJECT_ROOT/bin/cleandb"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示帮助信息
show_help() {
    cat << EOF
测试数据库清理脚本

用法:
  ./scripts/cleandb.sh [选项]

选项:
  -h, --help     显示此帮助信息
  -f, --force    强制清理所有测试数据，不进行确认
  -p, --prefix   只清理指定前缀的测试数据
  -d, --detail   显示详细输出
  --build        重新构建清理工具
  --no-build     跳过构建步骤（如果二进制文件已存在）

示例:
  ./scripts/cleandb.sh                    # 交互式清理所有测试数据
  ./scripts/cleandb.sh -f                 # 强制清理所有测试数据
  ./scripts/cleandb.sh -p test_user_       # 只清理前缀为 test_user_ 的数据
  ./scripts/cleandb.sh -f -d              # 强制清理并显示详细输出
  ./scripts/cleandb.sh --build -f         # 重新构建工具并强制清理

注意:
  - 此工具只清理测试数据库中的数据
  - 强制清理会删除所有测试相关的数据，请谨慎使用
  - 按前缀清理功能需要数据在创建时使用了相应的前缀
EOF
}

# 构建清理工具
build_cleandb() {
    print_info "构建测试数据库清理工具..."
    
    # 确保 bin 目录存在
    mkdir -p "$PROJECT_ROOT/bin"
    
    # 构建清理工具
    cd "$PROJECT_ROOT"
    if go build -o "$CLEANDB_CMD" "./cmd/cleandb"; then
        print_success "清理工具构建完成: $CLEANDB_CMD"
    else
        print_error "构建清理工具失败"
        exit 1
    fi
}

# 检查清理工具是否存在
check_cleandb_exists() {
    if [[ ! -f "$CLEANDB_CMD" ]]; then
        print_warning "清理工具不存在，开始构建..."
        build_cleandb
        return 0
    fi
    
    # 检查源代码是否比二进制文件新
    if [[ "$PROJECT_ROOT/cmd/cleandb/main.go" -nt "$CLEANDB_CMD" ]]; then
        print_warning "源代码已更新，重新构建清理工具..."
        build_cleandb
        return 0
    fi
    
    print_info "使用现有的清理工具: $CLEANDB_CMD"
}

# 解析命令行参数
BUILD_REQUIRED=false
SKIP_BUILD=false
CLEANDB_ARGS=()

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        --build)
            BUILD_REQUIRED=true
            shift
            ;;
        --no-build)
            SKIP_BUILD=true
            shift
            ;;
        *)
            # 将其他参数传递给 cleandb 工具
            CLEANDB_ARGS+=("$1")
            shift
            ;;
    esac
done

# 主逻辑
main() {
    print_info "开始测试数据库清理流程..."
    
    # 检查是否在项目根目录
    if [[ ! -f "$PROJECT_ROOT/go.mod" ]]; then
        print_error "未找到 go.mod 文件，请确保在项目根目录运行此脚本"
        exit 1
    fi
    
    # 构建或检查清理工具
    if [[ "$BUILD_REQUIRED" == "true" ]]; then
        build_cleandb
    elif [[ "$SKIP_BUILD" == "false" ]]; then
        check_cleandb_exists
    fi
    
    # 检查清理工具是否存在
    if [[ ! -f "$CLEANDB_CMD" ]]; then
        print_error "清理工具不存在且构建失败"
        exit 1
    fi
    
    # 运行清理工具
    print_info "执行数据库清理操作..."
    if "$CLEANDB_CMD" "${CLEANDB_ARGS[@]}"; then
        print_success "数据库清理操作完成"
    else
        print_error "数据库清理操作失败"
        exit 1
    fi
}

# 运行主函数
main "$@"