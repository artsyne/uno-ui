package helpers

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	. "github.com/biuren/hebe/tests/framework/types"
)

// CleanupManagerImpl 清理管理器实现
// 替代事务回滚，提供灵活的数据清理策略
type CleanupManagerImpl struct {
	db           *sql.DB
	config       *TestConfig
	isolator     DataIsolatorInterface
	cleanupFuncs []func() error
	mu           sync.RWMutex
}

// NewCleanupManager 创建清理管理器
func NewCleanupManager(db *sql.DB, config *TestConfig, isolator DataIsolatorInterface) *CleanupManagerImpl {
	return &CleanupManagerImpl{
		db:           db,
		config:       config,
		isolator:     isolator,
		cleanupFuncs: make([]func() error, 0),
	}
}

// Cleanup 执行清理
func (cm *CleanupManagerImpl) Cleanup(strategy string) error {
	if !cm.config.Cleanup.Enabled {
		slog.Debug("清理已禁用，数据将保留用于调试",
			"test_id", cm.isolator.GetTestID(),
			"prefix", cm.isolator.GetPrefix())
		return nil
	}

	slog.Debug("开始数据清理",
		"strategy", strategy,
		"test_id", cm.isolator.GetTestID(),
		"prefix", cm.isolator.GetPrefix())

	// 首先执行自定义清理函数
	if err := cm.executeCleanupFuncs(); err != nil {
		slog.Error("执行自定义清理函数失败", "error", err)
	}

	// 然后执行数据库清理
	switch strategy {
	case "prefix":
		return cm.CleanupByPrefix()
	case "tracked":
		return cm.CleanupTrackedData()
	case "immediate":
		return cm.cleanupImmediate()
	default:
		return cm.CleanupByPrefix()
	}
}

// CleanupByPrefix 按前缀清理
// 清理所有包含测试前缀的数据
func (cm *CleanupManagerImpl) CleanupByPrefix() error {
	// 按依赖关系顺序清理表（避免外键约束错误）
	tables := cm.getTablesInCleanupOrder()

	totalCleaned := 0
	for _, table := range tables {
		cleaned, err := cm.cleanupTableByPrefix(table)
		if err != nil {
			slog.Error("清理表失败", "table", table, "error", err)
			continue // 继续清理其他表
		}
		totalCleaned += cleaned
	}

	slog.Debug("按前缀清理完成",
		"total_cleaned", totalCleaned,
		"prefix", cm.isolator.GetPrefix())

	return nil
}

// CleanupTrackedData 清理跟踪的数据
// 只清理通过isolator跟踪的数据
func (cm *CleanupManagerImpl) CleanupTrackedData() error {
	createdData := cm.isolator.GetCreatedData()
	totalCleaned := 0

	// 按依赖关系顺序清理
	tables := cm.getTablesInCleanupOrder()

	for _, table := range tables {
		if ids, exists := createdData[table]; exists && len(ids) > 0 {
			cleaned, err := cm.cleanupTableByIDs(table, ids)
			if err != nil {
				slog.Error("清理跟踪数据失败", "table", table, "error", err)
				continue
			}
			totalCleaned += cleaned
		}
	}

	slog.Debug("跟踪数据清理完成",
		"total_cleaned", totalCleaned,
		"tables", len(createdData))

	return nil
}

// cleanupImmediate 立即清理
// 强制清理所有相关数据
func (cm *CleanupManagerImpl) cleanupImmediate() error {
	// 先尝试按跟踪数据清理
	if err := cm.CleanupTrackedData(); err != nil {
		slog.Warn("跟踪数据清理失败，尝试按前缀清理", "error", err)
	}

	// 再按前缀清理确保没有遗漏
	return cm.CleanupByPrefix()
}

// getTablesInCleanupOrder 获取清理顺序
// 考虑外键约束，子表优先清理
func (cm *CleanupManagerImpl) getTablesInCleanupOrder() []string {
	if len(cm.config.Cleanup.TablesOrder) > 0 {
		return cm.config.Cleanup.TablesOrder
	}

	// 使用配置中的默认顺序（在config.go中定义）
	return []string{
		"br_comment",      // 评论表（子表）
		"br_like",         // 点赞表（关联表）
		"br_follow",       // 关注表（关联表）
		"br_job_history",  // 任务历史表
		"br_job_queue",    // 任务队列表
		"br_diary_author", // 日记成员表（依赖日记和用户）
		"br_entry",        // 文章表（依赖日记和用户）
		"br_diary",        // 日记表（依赖用户）
		"br_user_stat",    // 用户统计表（依赖用户）
		"br_valid_code",   // 验证码表（依赖用户）
		"br_user",         // 用户表（主表）
	}
}

// cleanupTableByPrefix 按前缀清理表
func (cm *CleanupManagerImpl) cleanupTableByPrefix(table string) (int, error) {
	// 使用白名单验证表名，防止 SQL 注入
	validTables := map[string]bool{
		"br_user":         true,
		"br_diary":        true,
		"br_entry":        true,
		"br_comment":      true,
		"br_like":         true,
		"br_follow":       true,
		"br_diary_author": true,
		"br_job_queue":    true,
		"br_job_history":  true,
	}

	if !validTables[table] {
		slog.Debug("未知表，跳过清理", "table", table)
		return 0, nil
	}

	prefix := cm.isolator.GetPrefix() + "%"
	var query string
	var args []interface{}

	switch table {
	case "br_user":
		// User table uses handle (formerly user_name)
		query = "DELETE FROM br_user WHERE handle LIKE ? OR email LIKE ?"
		args = []interface{}{prefix, prefix}
	case "br_diary":
		// Clean by title prefix OR by users with prefix
		query = "DELETE FROM br_diary WHERE title LIKE ? OR add_user_id IN (SELECT user_id FROM br_user WHERE handle LIKE ?)"
		args = []interface{}{prefix, prefix}
	case "br_entry":
		// Clean by title prefix OR by users with prefix
		query = "DELETE FROM br_entry WHERE title LIKE ? OR add_user_id IN (SELECT user_id FROM br_user WHERE handle LIKE ?)"
		args = []interface{}{prefix, prefix}
	case "br_comment":
		// Clean comments by users with prefix
		query = "DELETE FROM br_comment WHERE add_user_id IN (SELECT user_id FROM br_user WHERE handle LIKE ?)"
		args = []interface{}{prefix}
	case "br_like":
		// Clean likes by users with prefix
		query = "DELETE FROM br_like WHERE user_id IN (SELECT user_id FROM br_user WHERE handle LIKE ?)"
		args = []interface{}{prefix}
	case "br_follow":
		// Clean follows by users with prefix
		query = "DELETE FROM br_follow WHERE user_id IN (SELECT user_id FROM br_user WHERE handle LIKE ?)"
		args = []interface{}{prefix}
	case "br_diary_author":
		// Clean members by users with prefix
		query = "DELETE FROM br_diary_author WHERE user_id IN (SELECT user_id FROM br_user WHERE handle LIKE ?)"
		args = []interface{}{prefix}
	case "br_job_queue", "br_job_history":
		// 对于 body 字段的 LIKE 查询需要特殊处理
		bodyPattern := "%" + cm.isolator.GetPrefix() + "%"
		query = "DELETE FROM " + table + " WHERE body LIKE ?"
		args = []interface{}{bodyPattern}
	}

	result, err := cm.db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("执行清理失败: %w", err)
	}

	affected, _ := result.RowsAffected()
	if affected > 0 {
		slog.Debug("表清理完成", "table", table, "affected", affected)
	}

	return int(affected), nil
}

// cleanupTableByIDs 按ID清理表数据
func (cm *CleanupManagerImpl) cleanupTableByIDs(table string, ids []interface{}) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	// 定义表和主键的映射关系（白名单）
	tablePrimaryKeys := map[string]string{
		"br_user":        "user_id",
		"br_diary":       "diary_id",
		"br_entry":       "entry_id",
		"br_comment":     "comment_id",
		"br_job_queue":   "job_id",
		"br_job_history": "job_id",
		"br_valid_code":  "valid_code_id",
		"br_like":        "like_id",
	}

	primaryKey, ok := tablePrimaryKeys[table]
	if !ok {
		return 0, fmt.Errorf("未知表: %s", table)
	}

	// 分批清理（避免IN子句过长）
	batchSize := cm.config.Cleanup.BatchSize
	totalCleaned := 0

	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize
		if end > len(ids) {
			end = len(ids)
		}

		batch := ids[i:end]

		// 构建参数化查询
		placeholders := make([]string, len(batch))
		for j := range placeholders {
			placeholders[j] = "?"
		}

		// 使用参数化查询，防止 SQL 注入
		query := fmt.Sprintf("DELETE FROM %s WHERE %s IN (%s)",
			table, primaryKey, strings.Join(placeholders, ","))

		result, err := cm.db.Exec(query, batch...)
		if err != nil {
			slog.Error("批量清理失败", "table", table, "batch", len(batch), "error", err)
			continue
		}

		affected, _ := result.RowsAffected()
		totalCleaned += int(affected)
	}

	if totalCleaned > 0 {
		slog.Debug("表ID清理完成", "table", table, "cleaned", totalCleaned)
	}

	return totalCleaned, nil
}

// AddCleanupFunc 添加自定义清理函数
func (cm *CleanupManagerImpl) AddCleanupFunc(fn func() error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.cleanupFuncs = append(cm.cleanupFuncs, fn)
}

// executeCleanupFuncs 执行自定义清理函数
func (cm *CleanupManagerImpl) executeCleanupFuncs() error {
	cm.mu.RLock()
	funcs := make([]func() error, len(cm.cleanupFuncs))
	copy(funcs, cm.cleanupFuncs)
	cm.mu.RUnlock()

	var lastErr error
	for i, fn := range funcs {
		if err := fn(); err != nil {
			slog.Error("执行清理函数失败", "index", i, "error", err)
			lastErr = err
		}
	}

	return lastErr
}

// GetStats 获取清理统计信息
func (cm *CleanupManagerImpl) GetStats() map[string]interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return map[string]interface{}{
		"cleanup_enabled":    cm.config.Cleanup.Enabled,
		"cleanup_strategy":   cm.config.Cleanup.Strategy,
		"batch_size":         cm.config.Cleanup.BatchSize,
		"custom_funcs_count": len(cm.cleanupFuncs),
		"tables_order":       cm.config.Cleanup.TablesOrder,
	}
}

// ForceCleanupDatabase 强制清理数据库
// 无论配置如何，都会执行数据库清理操作
func (cm *CleanupManagerImpl) ForceCleanupDatabase() error {
	slog.Info("强制清理数据库",
		"test_id", cm.isolator.GetTestID(),
		"prefix", cm.isolator.GetPrefix())

	// 强制执行按前缀清理，忽略配置中的Enabled设置
	tables := cm.getTablesInCleanupOrder()
	totalCleaned := 0

	for _, table := range tables {
		cleaned, err := cm.cleanupTableByPrefix(table)
		if err != nil {
			slog.Error("强制清理表失败", "table", table, "error", err)
			continue // 继续清理其他表
		}
		totalCleaned += cleaned
	}

	slog.Info("强制清理数据库完成",
		"total_cleaned", totalCleaned,
		"prefix", cm.isolator.GetPrefix())

	return nil
}

// ClearCustomFuncs 清除自定义清理函数
func (cm *CleanupManagerImpl) ClearCustomFuncs() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.cleanupFuncs = make([]func() error, 0)
}
