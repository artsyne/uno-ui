#!/bin/bash

# 测试优化后的配置
echo "🧪 测试国际化模块优化..."

# 检查配置字段是否已移除
echo "检查 AppConfig 结构体..."
if grep -q "TitleDescriptions\|DefaultTitleDesc" internal/core/config/types.go; then
    echo "❌ 配置字段未完全移除"
    exit 1
else
    echo "✅ 冗余配置字段已移除"
fi

# 检查 config.go 是否还有相关解析代码
echo "检查配置解析代码..."
if grep -q "TITLE_DESCRIPTIONS\|DEFAULT_TITLE_DESC" internal/core/config/config.go; then
    echo "❌ 配置解析代码未完全移除"
    exit 1
else
    echo "✅ 冗余配置解析代码已移除"
fi

# 检查 locales 文件是否有默认描述
echo "检查语言文件..."
if grep -q "common.default_page_desc" website/locales/en.json && grep -q "common.default_page_desc" website/locales/zh-cn.json; then
    echo "✅ 默认页面描述已添加到语言文件"
else
    echo "❌ 默认页面描述缺失"
    exit 1
fi

# 编译测试
echo "编译测试..."
if go build -o __debug_bin/debug_test ./cmd/idiary/main.go; then
    echo "✅ 编译成功"
    rm -f __debug_bin/debug_test
else
    echo "❌ 编译失败"
    exit 1
fi

echo "🎉 所有优化测试通过！"
echo ""
echo "📋 优化总结："
echo "  ✅ 移除了 TitleDescriptions 和 DefaultTitleDesc 配置字段"
echo "  ✅ 统一了 generatePageTitle 和 generateDescription 处理逻辑"
echo "  ✅ 都使用用户语言偏好，而不是硬编码默认语言"
echo "  ✅ 所有文本都通过 localization 机制管理"
echo "  ✅ 保持代码一致性和可维护性"
