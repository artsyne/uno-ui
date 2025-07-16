package mock

import (
	"database/sql"
	"testing"

	"github.com/biuren/hebe/internal/core/di"
	_ "github.com/go-sql-driver/mysql"
)

// DatabaseMock 数据库模拟器
type DatabaseMock struct {
	DB        *sql.DB
	Container *di.Container
}

// SetupMockDatabase 设置模拟数据库
// 注意：这里使用nil作为模拟数据库，在实际测试中应该使用内存数据库或真实测试数据库
func SetupMockDatabase(t *testing.T, container *di.Container) *DatabaseMock {
	t.Helper()

	// 在实际项目中，这里应该创建测试数据库连接
	// 例如：sqlite内存数据库或测试MySQL实例
	var mockDB *sql.DB = nil

	// 绑定到指定的DI容器
	err := di.Bind(container, func() *sql.DB {
		return mockDB
	})
	if err != nil {
		t.Fatalf("绑定数据库到DI容器失败: %v", err)
	}

	return &DatabaseMock{
		DB:        mockDB,
		Container: container,
	}
}

// Cleanup 清理模拟数据库
func (m *DatabaseMock) Cleanup(t *testing.T) {
	t.Helper()
	if m.DB != nil {
		m.DB.Close()
	}
}
