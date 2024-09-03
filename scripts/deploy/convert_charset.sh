#!/bin/bash

# 函数：显示脚本用法
usage() {
    echo "用法: $0 -db <数据库名称>"
    echo ""
    echo "此脚本将指定 MySQL 数据库中所有表的字符集和排序规则修改为 utf8mb4 和 utf8mb4_general_ci。"
    echo ""
    echo "选项:"
    echo "  -db   数据库名称"
    echo "  --help  显示此帮助信息"
    echo ""
    echo "配置 MySQL 客户端以避免在命令行上传递密码："
    echo "  使用 mysql_config_editor 工具来安全地存储数据库凭据："
    echo "  mysql_config_editor set --login-path=local --host=localhost --user=your_username --password"
    echo ""
    echo "赋予脚本可执行权限的命令行："
    echo "  macOS/Linux: chmod +x convert_charset.sh"
    echo "  Windows (使用 Git Bash 或 WSL): chmod +x convert_charset.sh"
    echo "  Windows (使用 PowerShell): Set-ExecutionPolicy RemoteSigned -Scope CurrentUser"
    exit 1
}

# 使用 getopts 解析命令行参数
while [[ "$#" -gt 0 ]]; do
    case $1 in
        -db) DB_NAME="$2"; shift ;;
        --help) usage ;;
        *) echo "无效的选项: $1" >&2; usage ;;
    esac
    shift
done

# 检查是否提供了足够的参数
if [ -z "$DB_NAME" ]; then
    usage
fi

echo "正在获取数据库 '$DB_NAME' 中的所有表名..."
# 获取所有表名
TABLES=$(mysql --login-path=local -D"$DB_NAME" -e 'SHOW TABLES;' | awk '{ print $1}' | grep -v '^Tables' )

if [ -z "$TABLES" ]; then
    echo "未找到任何表。请检查数据库名称是否正确。"
    exit 1
fi

echo "生成修改字符集的 SQL 脚本..."
# 生成修改字符集的 SQL 脚本
SQL_FILE="alter_tables.sql"
echo "" > $SQL_FILE
for TABLE in $TABLES; do
    echo "ALTER TABLE $TABLE CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;" >> $SQL_FILE
done

echo "执行生成的 SQL 脚本..."
# 执行生成的 SQL 脚本
mysql --login-path=local -D"$DB_NAME" < $SQL_FILE

if [ $? -eq 0 ]; then
    echo "字符集和排序规则修改成功。"
else
    echo "字符集和排序规则修改失败。请检查错误信息。"
    exit 1
fi

# 清理生成的 SQL 文件
rm $SQL_FILE
echo "清理完成。"