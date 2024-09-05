#!/bin/bash
#chmod +x scripts/run_debug.sh

# 获取脚本所在目录
SCRIPT_DIR=$(dirname "$0")
PROJECT_DIR=$(cd "$SCRIPT_DIR/.." && pwd)

# 显示当前项目目录
echo "Current project directory: $PROJECT_DIR"

# 设置环境变量
export CONFIG_PATH="$PROJECT_DIR/configs"
export STATIC_PATH="$PROJECT_DIR/website/static"
export VIEW_PATH="$PROJECT_DIR/website/views"

# 启动应用程序
if [[ "$OSTYPE" == "msys" ]]; then
    # Windows (Git Bash)
    "$PROJECT_DIR/__debug_bin/debug.idiary.site.exe"
else
    # macOS/Linux
    "$PROJECT_DIR/__debug_bin/debug.idiary.site"
fi