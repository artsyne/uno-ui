package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	. "github.com/biuren/hebe/tests/framework/types"
)

// DataIsolatorImpl 数据隔离器实现
// 使用唯一前缀替代事务，实现测试数据隔离，这是新框架的核心创新
type DataIsolatorImpl struct {
	testID      string
	prefix      string
	config      *TestConfig
	createdData map[string][]interface{} // 记录创建的数据
	mu          sync.RWMutex
}

// NewDataIsolator 创建数据隔离器
func NewDataIsolator(testID string, config *TestConfig) *DataIsolatorImpl {
	prefix := generatePrefix(testID, config.Isolation)

	return &DataIsolatorImpl{
		testID:      testID,
		prefix:      prefix,
		config:      config,
		createdData: make(map[string][]interface{}),
	}
}

// generatePrefix 生成唯一前缀
// 格式: test_0801_1430_abc123 (最大20字符)
// 确保每个测试都有唯一的数据前缀，避免测试间数据冲突
func generatePrefix(testID string, config IsolationConfig) string {
	now := time.Now()

	// 生成随机字符串
	randomBytes := make([]byte, 3)
	rand.Read(randomBytes)
	random := hex.EncodeToString(randomBytes)

	// 根据配置格式生成前缀
	prefix := config.PrefixFormat
	prefix = strings.ReplaceAll(prefix, "{date}", now.Format("0102")) // MMDD
	prefix = strings.ReplaceAll(prefix, "{time}", now.Format("1504")) // HHMM
	prefix = strings.ReplaceAll(prefix, "{random}", random)           // 随机字符串
	prefix = strings.ReplaceAll(prefix, "{test_id}", testID)          // 测试ID

	// 确保前缀长度不超过限制
	if len(prefix) > config.MaxPrefixLen {
		prefix = prefix[:config.MaxPrefixLen]
	}

	return prefix
}

// IsolateString 隔离字符串字段
func (di *DataIsolatorImpl) IsolateString(value string) string {
	return fmt.Sprintf("%s_%s", di.prefix, value)
}

// IsolateEmail 隔离邮箱字段
// 特殊处理邮箱格式，确保邮箱格式的有效性
func (di *DataIsolatorImpl) IsolateEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return di.IsolateString(email)
	}
	return fmt.Sprintf("%s_%s@%s", di.prefix, parts[0], parts[1])
}

// IsolateUsername 隔离用户名
// 确保用户名的唯一性和格式正确性
func (di *DataIsolatorImpl) IsolateUsername(username string) string {
	isolated := di.IsolateString(username)
	// 确保用户名长度不超过数据库限制（通常100字符）
	if len(isolated) > 90 { // 留出一些余量
		isolated = isolated[:90]
	}
	return isolated
}

// IsolateTitle 隔离标题字段
func (di *DataIsolatorImpl) IsolateTitle(title string) string {
	isolated := di.IsolateString(title)
	// 确保标题长度不超过数据库限制（通常200字符）
	if len(isolated) > 180 { // 留出一些余量
		isolated = isolated[:180]
	}
	return isolated
}

// TrackCreated 记录创建的数据
// 用于后续的数据清理
func (di *DataIsolatorImpl) TrackCreated(table string, id interface{}) {
	di.mu.Lock()
	defer di.mu.Unlock()

	di.createdData[table] = append(di.createdData[table], id)
}

// TrackCreatedBatch 批量记录创建的数据
func (di *DataIsolatorImpl) TrackCreatedBatch(table string, ids []interface{}) {
	di.mu.Lock()
	defer di.mu.Unlock()

	di.createdData[table] = append(di.createdData[table], ids...)
}

// GetCreatedData 获取创建的数据
func (di *DataIsolatorImpl) GetCreatedData() map[string][]interface{} {
	di.mu.RLock()
	defer di.mu.RUnlock()

	result := make(map[string][]interface{})
	for table, ids := range di.createdData {
		result[table] = make([]interface{}, len(ids))
		copy(result[table], ids)
	}
	return result
}

// GetCreatedDataForTable 获取指定表的创建数据
func (di *DataIsolatorImpl) GetCreatedDataForTable(table string) []interface{} {
	di.mu.RLock()
	defer di.mu.RUnlock()

	if ids, exists := di.createdData[table]; exists {
		result := make([]interface{}, len(ids))
		copy(result, ids)
		return result
	}
	return nil
}

// GetPrefix 获取隔离前缀
func (di *DataIsolatorImpl) GetPrefix() string {
	return di.prefix
}

// GetTestID 获取测试ID
func (di *DataIsolatorImpl) GetTestID() string {
	return di.testID
}

// ClearTrackedData 清除跟踪的数据
func (di *DataIsolatorImpl) ClearTrackedData() {
	di.mu.Lock()
	defer di.mu.Unlock()

	di.createdData = make(map[string][]interface{})
}

// GetDataCount 获取创建的数据总数
func (di *DataIsolatorImpl) GetDataCount() int {
	di.mu.RLock()
	defer di.mu.RUnlock()

	count := 0
	for _, ids := range di.createdData {
		count += len(ids)
	}
	return count
}

// GetTableCount 获取涉及的表数量
func (di *DataIsolatorImpl) GetTableCount() int {
	di.mu.RLock()
	defer di.mu.RUnlock()

	return len(di.createdData)
}

// BuildWhereLikeClause 构建LIKE查询条件
// 用于按前缀查找数据
func (di *DataIsolatorImpl) BuildWhereLikeClause(column string) string {
	return fmt.Sprintf("%s LIKE '%s%%'", column, di.prefix)
}

// BuildWhereLikeClauseMultiple 构建多字段LIKE查询条件
func (di *DataIsolatorImpl) BuildWhereLikeClauseMultiple(columns []string) string {
	var conditions []string
	for _, column := range columns {
		conditions = append(conditions, di.BuildWhereLikeClause(column))
	}
	return strings.Join(conditions, " OR ")
}

// IsIsolated 检查字符串是否已被隔离
func (di *DataIsolatorImpl) IsIsolated(value string) bool {
	return strings.HasPrefix(value, di.prefix+"_")
}

// GetStats 获取隔离器统计信息
func (di *DataIsolatorImpl) GetStats() map[string]interface{} {
	di.mu.RLock()
	defer di.mu.RUnlock()

	stats := map[string]interface{}{
		"test_id":     di.testID,
		"prefix":      di.prefix,
		"data_count":  di.GetDataCount(),
		"table_count": len(di.createdData),
		"tables":      make(map[string]int),
	}

	for table, ids := range di.createdData {
		stats["tables"].(map[string]int)[table] = len(ids)
	}

	return stats
}
