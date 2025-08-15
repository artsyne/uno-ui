#!/bin/bash

# iDiary.ai Hebe 多语言优化自动化脚本
# 版本: v1.0
# 作者: AI Assistant
# 用途: 自动化执行多语言优化的各个步骤

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKUP_DIR="$PROJECT_ROOT/backup/i18n/$(date +%Y%m%d_%H%M%S)"
LOCALES_DIR="$PROJECT_ROOT/website/locales"
VIEWS_DIR="$PROJECT_ROOT/website/views"

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示帮助信息
show_help() {
    cat << EOF
iDiary.ai Hebe 多语言优化自动化脚本

用法: $0 [选项] [命令]

命令:
  backup              创建备份
  extract             提取硬编码中文文本
  validate            验证语言文件
  check-syntax        检查模板语法
  test                运行多语言测试
  performance         运行性能测试
  replace-navbar      替换导航栏文本
  replace-auth        替换认证页面文本
  replace-diary       替换日记页面文本
  replace-forms       替换表单文本
  replace-all         替换所有硬编码文本
  full-optimization   执行完整优化流程
  help                显示此帮助信息

选项:
  -v, --verbose       详细输出
  -y, --yes           自动确认所有操作
  --dry-run           仅显示将要执行的操作，不实际执行

示例:
  $0 backup                    # 创建备份
  $0 extract                   # 提取硬编码文本
  $0 replace-all --dry-run     # 预览替换操作
  $0 full-optimization -y      # 执行完整优化（自动确认）

EOF
}

# 检查依赖
check_dependencies() {
    log_info "检查依赖工具..."
    
    local missing_tools=()
    
    # 检查 jq
    if ! command -v jq &> /dev/null; then
        missing_tools+=("jq")
    fi
    
    # 检查 grep
    if ! command -v grep &> /dev/null; then
        missing_tools+=("grep")
    fi
    
    # 检查 sed
    if ! command -v sed &> /dev/null; then
        missing_tools+=("sed")
    fi
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        log_error "缺少必要工具: ${missing_tools[*]}"
        log_info "请安装缺少的工具后重试"
        exit 1
    fi
    
    log_success "依赖检查通过"
}

# 创建备份
create_backup() {
    log_info "创建备份到: $BACKUP_DIR"
    
    mkdir -p "$BACKUP_DIR"
    
    # 备份视图文件
    if [ -d "$VIEWS_DIR" ]; then
        cp -r "$VIEWS_DIR" "$BACKUP_DIR/views"
        log_success "视图文件备份完成"
    fi
    
    # 备份语言文件
    if [ -d "$LOCALES_DIR" ]; then
        cp -r "$LOCALES_DIR" "$BACKUP_DIR/locales"
        log_success "语言文件备份完成"
    fi
    
    # 创建备份信息文件
    cat > "$BACKUP_DIR/backup_info.txt" << EOF
备份创建时间: $(date)
项目路径: $PROJECT_ROOT
备份内容:
- 视图文件 (website/views)
- 语言文件 (website/locales)

恢复命令:
cp -r $BACKUP_DIR/views/* $VIEWS_DIR/
cp -r $BACKUP_DIR/locales/* $LOCALES_DIR/
EOF
    
    log_success "备份创建完成: $BACKUP_DIR"
}

# 提取硬编码中文文本
extract_chinese_text() {
    log_info "提取硬编码中文文本..."
    
    local output_file="$PROJECT_ROOT/hardcoded_chinese.txt"
    
    # 查找包含中文的HTML文件
    find "$VIEWS_DIR" -name "*.html" -exec grep -Hn "[\u4e00-\u9fa5]" {} \; > "$output_file" 2>/dev/null || true
    
    if [ -s "$output_file" ]; then
        log_success "硬编码中文文本已提取到: $output_file"
        
        # 统计信息
        local total_lines=$(wc -l < "$output_file")
        local total_files=$(awk -F: '{print $1}' "$output_file" | sort | uniq | wc -l)
        
        log_info "统计信息:"
        log_info "  - 包含中文的行数: $total_lines"
        log_info "  - 涉及文件数量: $total_files"
        
        # 显示前10行作为示例
        log_info "前10行示例:"
        head -10 "$output_file" | while read line; do
            echo "    $line"
        done
    else
        log_success "未发现硬编码中文文本"
    fi
}

# 验证语言文件
validate_language_files() {
    log_info "验证语言文件..."
    
    local en_file="$LOCALES_DIR/en.json"
    local zh_file="$LOCALES_DIR/zh-cn.json"
    local errors=0
    
    # 检查文件是否存在
    if [ ! -f "$en_file" ]; then
        log_error "英文语言文件不存在: $en_file"
        ((errors++))
    fi
    
    if [ ! -f "$zh_file" ]; then
        log_error "中文语言文件不存在: $zh_file"
        ((errors++))
    fi
    
    if [ $errors -gt 0 ]; then
        return 1
    fi
    
    # 验证JSON格式
    if ! jq . "$en_file" > /dev/null 2>&1; then
        log_error "英文语言文件JSON格式错误"
        ((errors++))
    else
        log_success "英文语言文件格式正确"
    fi
    
    if ! jq . "$zh_file" > /dev/null 2>&1; then
        log_error "中文语言文件JSON格式错误"
        ((errors++))
    else
        log_success "中文语言文件格式正确"
    fi
    
    if [ $errors -gt 0 ]; then
        return 1
    fi
    
    # 验证键值对一致性
    local en_keys=$(jq -r 'paths(scalars) as $p | $p | join(".")' "$en_file" | sort)
    local zh_keys=$(jq -r 'paths(scalars) as $p | $p | join(".")' "$zh_file" | sort)
    
    if [ "$en_keys" = "$zh_keys" ]; then
        local key_count=$(echo "$en_keys" | wc -l)
        log_success "键值对一致性检查通过，共 $key_count 个键值对"
    else
        log_error "键值对不一致"
        
        # 显示差异
        local en_only=$(comm -23 <(echo "$en_keys") <(echo "$zh_keys"))
        local zh_only=$(comm -13 <(echo "$en_keys") <(echo "$zh_keys"))
        
        if [ -n "$en_only" ]; then
            log_warning "英文独有的键:"
            echo "$en_only" | while read key; do
                echo "    $key"
            done
        fi
        
        if [ -n "$zh_only" ]; then
            log_warning "中文独有的键:"
            echo "$zh_only" | while read key; do
                echo "    $key"
            done
        fi
        
        ((errors++))
    fi
    
    if [ $errors -eq 0 ]; then
        log_success "语言文件验证通过"
        return 0
    else
        log_error "语言文件验证失败"
        return 1
    fi
}

# 检查模板语法
check_template_syntax() {
    log_info "检查模板语法..."
    
    local errors=0
    
    find "$VIEWS_DIR" -name "*.html" | while read file; do
        local relative_path=${file#$PROJECT_ROOT/}
        
        # 检查未闭合的 getMessage 调用
        local unclosed_getmessage=$(grep -n "getMessage.*[^}]$" "$file" || true)
        if [ -n "$unclosed_getmessage" ]; then
            log_error "$relative_path: 发现未闭合的 getMessage 调用"
            echo "$unclosed_getmessage" | while read line; do
                echo "    $line"
            done
            ((errors++))
        fi
        
        # 检查未闭合的模板标签
        local unclosed_tags=$(grep -n "{{[^}]*[^}]$" "$file" || true)
        if [ -n "$unclosed_tags" ]; then
            log_error "$relative_path: 发现未闭合的模板标签"
            echo "$unclosed_tags" | while read line; do
                echo "    $line"
            done
            ((errors++))
        fi
        
        # 检查可能的硬编码中文（排除已使用 getMessage 的行）
        local hardcoded_chinese=$(grep -n "[\u4e00-\u9fa5]" "$file" | grep -v "getMessage" || true)
        if [ -n "$hardcoded_chinese" ]; then
            log_warning "$relative_path: 发现可能的硬编码中文"
            echo "$hardcoded_chinese" | while read line; do
                echo "    $line"
            done
        fi
    done
    
    if [ $errors -eq 0 ]; then
        log_success "模板语法检查通过"
        return 0
    else
        log_error "模板语法检查发现 $errors 个错误"
        return 1
    fi
}

# 替换导航栏文本
replace_navbar_text() {
    log_info "替换导航栏文本..."
    
    local navbar_file="$VIEWS_DIR/navbar.html"
    
    if [ ! -f "$navbar_file" ]; then
        log_error "导航栏文件不存在: $navbar_file"
        return 1
    fi
    
    # 创建临时文件
    local temp_file=$(mktemp)
    
    # 执行替换
    sed -e 's/免费注册/{{getMessage "auth.button.register_free" .Context.Lang}}/g' \
        -e 's/登录/{{getMessage "auth.button.login" .Context.Lang}}/g' \
        -e 's/我的主页/{{getMessage "user.label.my_home" .Context.Lang}}/g' \
        -e 's/用户设置/{{getMessage "user.label.settings" .Context.Lang}}/g' \
        -e 's/退出登录/{{getMessage "auth.button.logout" .Context.Lang}}/g' \
        "$navbar_file" > "$temp_file"
    
    # 检查替换结果
    if [ -s "$temp_file" ]; then
        mv "$temp_file" "$navbar_file"
        log_success "导航栏文本替换完成"
    else
        rm -f "$temp_file"
        log_error "导航栏文本替换失败"
        return 1
    fi
}

# 替换认证页面文本
replace_auth_text() {
    log_info "替换认证页面文本..."
    
    local login_file="$VIEWS_DIR/passport/login.html"
    local register_file="$VIEWS_DIR/passport/register.html"
    
    # 替换登录页面
    if [ -f "$login_file" ]; then
        local temp_file=$(mktemp)
        
        sed -e 's/用户名或邮箱/{{getMessage "auth.label.username_or_email" .Context.Lang}}/g' \
            -e 's/密码/{{getMessage "auth.label.password" .Context.Lang}}/g' \
            -e 's/>登录</>{{getMessage "auth.button.login" .Context.Lang}}</g' \
            -e 's/忘记密码？/{{getMessage "auth.link.forgot_password" .Context.Lang}}/g' \
            -e 's/免费注册/{{getMessage "auth.button.register_free" .Context.Lang}}/g' \
            "$login_file" > "$temp_file"
        
        if [ -s "$temp_file" ]; then
            mv "$temp_file" "$login_file"
            log_success "登录页面文本替换完成"
        else
            rm -f "$temp_file"
            log_error "登录页面文本替换失败"
        fi
    fi
    
    # 替换注册页面
    if [ -f "$register_file" ]; then
        local temp_file=$(mktemp)
        
        sed -e 's/用户名/{{getMessage "auth.label.username" .Context.Lang}}/g' \
            -e 's/密码/{{getMessage "auth.label.password" .Context.Lang}}/g' \
            -e 's/电子邮箱地址/{{getMessage "auth.label.email_address" .Context.Lang}}/g' \
            -e 's/已经有了账户？/{{getMessage "auth.link.already_have_account" .Context.Lang}}/g' \
            "$register_file" > "$temp_file"
        
        if [ -s "$temp_file" ]; then
            mv "$temp_file" "$register_file"
            log_success "注册页面文本替换完成"
        else
            rm -f "$temp_file"
            log_error "注册页面文本替换失败"
        fi
    fi
}

# 替换日记页面文本
replace_diary_text() {
    log_info "替换日记页面文本..."
    
    local diary_index="$VIEWS_DIR/diary/index.html"
    
    if [ -f "$diary_index" ]; then
        local temp_file=$(mktemp)
        
        sed -e 's/日记本设置/{{getMessage "diary.button.settings" .Context.Lang}}/g' \
            -e 's/点赞/{{getMessage "common.like" .Context.Lang}}/g' \
            -e 's/日记本介绍/{{getMessage "diary.label.introduction" .Context.Lang}}/g' \
            -e 's/目录/{{getMessage "common.menu" .Context.Lang}}/g' \
            -e 's/作者列表/{{getMessage "diary.label.authors_list" .Context.Lang}}/g' \
            -e 's/点赞会员列表/{{getMessage "diary.label.liked_members_list" .Context.Lang}}/g' \
            -e 's/写日记/{{getMessage "diary.button.write_entry" .Context.Lang}}/g' \
            -e 's/编辑/{{getMessage "common.edit" .Context.Lang}}/g' \
            "$diary_index" > "$temp_file"
        
        if [ -s "$temp_file" ]; then
            mv "$temp_file" "$diary_index"
            log_success "日记页面文本替换完成"
        else
            rm -f "$temp_file"
            log_error "日记页面文本替换失败"
        fi
    fi
}

# 替换表单文本
replace_forms_text() {
    log_info "替换表单文本..."
    
    local diary_form="$VIEWS_DIR/form/diary.html"
    local post_form="$VIEWS_DIR/form/post.html"
    
    # 替换日记表单
    if [ -f "$diary_form" ]; then
        local temp_file=$(mktemp)
        
        sed -e 's/日记本名称/{{getMessage "diary.form.name_required" .Context.Lang}}/g' \
            -e 's/日记本介绍/{{getMessage "diary.form.description" .Context.Lang}}/g' \
            -e 's/私人日记本/{{getMessage "diary.form.privacy_private" .Context.Lang}}/g' \
            -e 's/公开日记本/{{getMessage "diary.form.privacy_public" .Context.Lang}}/g' \
            -e 's/关闭社交功能/{{getMessage "diary.form.social_disabled" .Context.Lang}}/g' \
            -e 's/开启社交功能/{{getMessage "diary.form.social_enabled" .Context.Lang}}/g' \
            "$diary_form" > "$temp_file"
        
        if [ -s "$temp_file" ]; then
            mv "$temp_file" "$diary_form"
            log_success "日记表单文本替换完成"
        else
            rm -f "$temp_file"
            log_error "日记表单文本替换失败"
        fi
    fi
    
    # 替换文章表单
    if [ -f "$post_form" ]; then
        local temp_file=$(mktemp)
        
        sed -e 's/标题(非必填)/{{getMessage "post.form.title_optional" .Context.Lang}}/g' \
            -e 's/日记内容/{{getMessage "post.form.content_required" .Context.Lang}}/g' \
            -e 's/请输入日记内容.../{{getMessage "post.form.content_placeholder" .Context.Lang}}/g' \
            "$post_form" > "$temp_file"
        
        if [ -s "$temp_file" ]; then
            mv "$temp_file" "$post_form"
            log_success "文章表单文本替换完成"
        else
            rm -f "$temp_file"
            log_error "文章表单文本替换失败"
        fi
    fi
}

# 替换所有硬编码文本
replace_all_text() {
    log_info "开始替换所有硬编码文本..."
    
    replace_navbar_text
    replace_auth_text
    replace_diary_text
    replace_forms_text
    
    log_success "所有文本替换完成"
}

# 运行多语言测试
run_i18n_test() {
    log_info "运行多语言测试..."
    
    # 这里需要根据实际项目的测试命令进行调整
    if [ -f "$PROJECT_ROOT/go.mod" ]; then
        log_info "运行Go测试..."
        cd "$PROJECT_ROOT"
        if go test ./internal/utils/i18n/... -v; then
            log_success "Go测试通过"
        else
            log_error "Go测试失败"
            return 1
        fi
    fi
    
    # 检查服务器是否运行
    if curl -s http://localhost:8080 > /dev/null 2>&1; then
        log_info "测试多语言页面..."
        
        # 测试中文页面
        if curl -s -H "Accept-Language: zh-CN" http://localhost:8080/ | grep -q "getMessage\|登录\|注册"; then
            log_success "中文页面测试通过"
        else
            log_warning "中文页面可能存在问题"
        fi
        
        # 测试英文页面
        if curl -s -H "Accept-Language: en" http://localhost:8080/ | grep -q "getMessage\|Login\|Register"; then
            log_success "英文页面测试通过"
        else
            log_warning "英文页面可能存在问题"
        fi
    else
        log_warning "服务器未运行，跳过页面测试"
    fi
}

# 运行性能测试
run_performance_test() {
    log_info "运行性能测试..."
    
    if ! curl -s http://localhost:8080 > /dev/null 2>&1; then
        log_warning "服务器未运行，跳过性能测试"
        return 0
    fi
    
    for lang in "zh-CN" "en"; do
        log_info "测试语言: $lang"
        
        local start_time=$(date +%s%N)
        curl -s -H "Accept-Language: $lang" http://localhost:8080/ > /dev/null
        local end_time=$(date +%s%N)
        
        local duration=$(( (end_time - start_time) / 1000000 ))
        log_info "加载时间: ${duration}ms"
        
        if [ $duration -lt 2000 ]; then
            log_success "性能测试通过"
        else
            log_warning "加载时间较长: ${duration}ms"
        fi
    done
}

# 执行完整优化流程
full_optimization() {
    log_info "开始执行完整多语言优化流程..."
    
    # 步骤1: 检查依赖
    check_dependencies
    
    # 步骤2: 创建备份
    if [ "$AUTO_YES" != "true" ]; then
        read -p "是否创建备份？(y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            create_backup
        fi
    else
        create_backup
    fi
    
    # 步骤3: 提取硬编码文本
    extract_chinese_text
    
    # 步骤4: 验证语言文件
    if ! validate_language_files; then
        log_error "语言文件验证失败，请先修复语言文件"
        return 1
    fi
    
    # 步骤5: 替换硬编码文本
    if [ "$DRY_RUN" = "true" ]; then
        log_info "[DRY RUN] 将执行文本替换操作"
    else
        if [ "$AUTO_YES" != "true" ]; then
            read -p "是否开始替换硬编码文本？(y/N): " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                replace_all_text
            else
                log_info "跳过文本替换"
            fi
        else
            replace_all_text
        fi
    fi
    
    # 步骤6: 检查模板语法
    check_template_syntax
    
    # 步骤7: 运行测试
    if [ "$AUTO_YES" != "true" ]; then
        read -p "是否运行测试？(y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            run_i18n_test
        fi
    else
        run_i18n_test
    fi
    
    log_success "完整优化流程执行完成！"
    log_info "下一步建议:"
    log_info "  1. 检查替换结果"
    log_info "  2. 运行完整测试套件"
    log_info "  3. 提交代码变更"
}

# 解析命令行参数
VERBOSE=false
AUTO_YES=false
DRY_RUN=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -y|--yes)
            AUTO_YES=true
            shift
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        backup)
            create_backup
            exit 0
            ;;
        extract)
            extract_chinese_text
            exit 0
            ;;
        validate)
            validate_language_files
            exit $?
            ;;
        check-syntax)
            check_template_syntax
            exit $?
            ;;
        test)
            run_i18n_test
            exit $?
            ;;
        performance)
            run_performance_test
            exit $?
            ;;
        replace-navbar)
            replace_navbar_text
            exit $?
            ;;
        replace-auth)
            replace_auth_text
            exit $?
            ;;
        replace-diary)
            replace_diary_text
            exit $?
            ;;
        replace-forms)
            replace_forms_text
            exit $?
            ;;
        replace-all)
            replace_all_text
            exit $?
            ;;
        full-optimization)
            full_optimization
            exit $?
            ;;
        help|--help|-h)
            show_help
            exit 0
            ;;
        *)
            log_error "未知命令: $1"
            show_help
            exit 1
            ;;
    esac
done

# 如果没有提供命令，显示帮助
show_help