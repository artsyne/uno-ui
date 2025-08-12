package main

import (
	"flag"
	"fmt"
	"log"

	framework "github.com/biuren/hebe/tests/framework/core"
)

func main() {
	// 定义命令行参数
	var (
		help = flag.Bool("help", false, "显示帮助信息")
		h    = flag.Bool("h", false, "显示帮助信息")
		force = flag.Bool("force", false, "强制清理所有测试数据，不进行确认")
		f     = flag.Bool("f", false, "强制清理所有测试数据，不进行确认")
		prefix = flag.String("prefix", "", "只清理指定前缀的测试数据")
		p      = flag.String("p", "", "只清理指定前缀的测试数据")
		detail = flag.Bool("detail", false, "显示详细输出")
		d      = flag.Bool("d", false, "显示详细输出")
	)

	flag.Parse()

	// 显示帮助信息
	if *help || *h {
		showHelp()
		return
	}

	// 合并短参数和长参数
	isForce := *force || *f
	isDetail := *detail || *d
	cleanPrefix := *prefix
	if cleanPrefix == "" {
		cleanPrefix = *p
	}
	if isDetail {
		log.Println("开始清理测试数据库...")
	}
	// 获取数据库管理器
	dbManager := framework.GetDatabaseManager()
	if err := dbManager.Initialize(); err != nil {
		log.Fatalf("初始化数据库管理器失败: %v", err)
	}
	defer dbManager.Close()

	// 检查数据库健康状态
	if !dbManager.IsHealthy() {
		log.Fatal("数据库连接不健康，无法执行清理操作")
	}



	// 执行清理操作
	var err error
	if cleanPrefix != "" {
		log.Fatalf("按前缀清理功能已移除，请使用全量清理。")
	} else {
		// 全量清理
		if !isForce {
			// 询问用户确认
			dbName := dbManager.GetDatabaseName()
			fmt.Printf("警告：这将清理数据库 '%s' 中的所有测试数据！是否继续？(y/N): ", dbName)
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" && response != "yes" && response != "YES" {
				fmt.Println("操作已取消")
				return
			}
		}

		if isDetail {
			log.Println("强制清理所有测试数据...")
		}
		err = truncateAllTables(dbManager, isDetail)
	}

	if err != nil {
		log.Fatalf("清理操作失败: %v", err)
	}

	if isDetail {
		log.Println("测试数据库清理完成")
	} else {
		fmt.Println("清理完成")
	}
}

func showHelp() {
	fmt.Println(`测试数据库清理工具

用法:
  cleandb [选项]

选项:
  -h, --help     显示此帮助信息
  -f, --force    强制清理所有测试数据，不进行确认
  -p, --prefix   只清理指定前缀的测试数据
  -d, --detail   显示详细输出

示例:
  cleandb                    # 交互式清理所有测试数据
  cleandb -f                 # 强制清理所有测试数据
  cleandb -p test_user_      # 只清理前缀为 test_user_ 的数据
  cleandb -f -d              # 强制清理并显示详细输出

注意:
  - 此工具只清理测试数据库中的数据
  - 强制清理会删除所有测试相关的数据，请谨慎使用
  - 按前缀清理功能需要数据在创建时使用了相应的前缀`)
}

// truncateAllTables 清理数据库中所有表的测试数据
func truncateAllTables(dbManager *framework.DatabaseManager, isDetail bool) error {
	db := dbManager.GetRawDB()

	// 获取所有表名
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return fmt.Errorf("获取表名失败: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName);
		err != nil {
			return fmt.Errorf("扫描表名失败: %w", err)
		}
		tables = append(tables, tableName)
	}

	// 禁用外键检查
	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		return fmt.Errorf("禁用外键检查失败: %w", err)
	}

	// 遍历并截断所有表
	for _, table := range tables {
		if isDetail {
			log.Printf("截断表: %s", table)
		}
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table))
		if err != nil {
			return fmt.Errorf("截断表 %s 失败: %w", table, err)
		}
	}

	// 重新启用外键检查
	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	if err != nil {
		return fmt.Errorf("启用外键检查失败: %w", err)
	}

	return nil
}