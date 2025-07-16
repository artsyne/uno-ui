# 测试框架使用指南

本项目的测试框架提供了统一的测试环境配置、模拟服务和HTTP测试工具，旨在简化测试编写和提高测试质量。

## 目录结构

```
tests/
├── helpers/                    # 测试助手工具
│   ├── config.go              # 配置助手
│   ├── suite.go               # 测试套件
│   ├── http.go                # HTTP测试助手
│   └── database.go            # 数据库助手（现有）
├── testdata/                  # 测试数据
│   ├── mock/                  # 模拟服务
│   │   ├── database.go        # 数据库模拟器
│   │   └── services.go        # 服务模拟器
│   └── fixtures/              # 测试数据文件
├── unit/                      # 单元测试
│   ├── core/                  # 核心模块测试
│   ├── app/                   # 应用层测试
│   ├── modules/               # 业务模块测试
│   └── pkg/                   # 工具包测试
├── integration/               # 集成测试
│   ├── api/                   # API集成测试
│   └── database/              # 数据库集成测试
└── e2e/                       # 端到端测试
```

## 快速开始

### 1. 基本测试结构

使用 `testhelpers.WithCleanup` 创建自动清理的测试环境：

```go
package mypackage

import (
    "testing"
    testhelpers "github.com/biuren/hebe/tests/helpers"
)

func TestMyFunction(t *testing.T) {
    testhelpers.WithCleanup(t, func(suite *testhelpers.TestSuite) {
        // 在这里编写测试代码
        // 测试结束后会自动清理资源
    })
}
```

### 2. 测试套件组件

测试套件提供以下组件：

- `suite.Config`: 测试配置助手
- `suite.Container`: DI容器
- `suite.Database`: 数据库模拟器
- `suite.User`: 用户服务模拟器
- `suite.Diary`: 日记服务模拟器
- `suite.Post`: 文章服务模拟器
- `suite.Task`: 任务服务模拟器
  ├── testdata/ # 测试数据
  │ ├── fixtures/ # 测试固定数据
  │ └── mock/ # 模拟数据
  └── helpers/ # 测试辅助工具
  ├── database.go # 数据库测试工具
  └── http.go # HTTP测试工具

````

## 测试辅助工具

### DatabaseTestHelper

提供数据库测试的便利方法：

- 事务管理（自动回滚）
- 数据准备和清理
- 断言方法

### HTTPTestHelper

提供HTTP测试的便利方法：

- 请求构建和发送
- 响应验证
- 认证处理

## 运行测试

### 运行所有测试

```bash
cd tests
go test ./... -v
````

### 运行特定模块测试

```bash
# 用户模块测试
go test ./unit/modules/user/ -v

# 任务模块测试
go test ./unit/modules/task/ -v

# pkg包测试
go test ./unit/pkg/ -v
```

### 运行集成测试

```bash
go test ./integration/... -v
```

## 测试环境配置

### 数据库配置

测试需要连接到测试数据库。确保以下环境变量已设置：

```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_NAME=hebe_test
export DB_USER=test_user
export DB_PASSWORD=test_password
```

### 邮件配置（可选）

某些测试需要邮件配置：

```bash
export SMTP_SERVER=smtp.example.com
export SMTP_USER=test@example.com
export SMTP_PASSWORD=test_password
```

## 测试覆盖率

查看测试覆盖率：

```bash
go test ./... -cover
```

生成覆盖率报告：

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## 测试原则

1. **隔离性**: 每个测试都应该独立运行，不依赖其他测试
2. **可重复性**: 测试结果应该一致且可重复
3. **快速性**: 单元测试应该快速执行
4. **可读性**: 测试名称和内容应该清晰易懂
5. **全面性**: 覆盖正常情况、边界情况和异常情况

## 持续集成

测试应该在以下情况下自动运行：

- 代码提交时
- 合并请求时
- 发布构建时

目标测试覆盖率: **80%以上**

## 数据库测试

框架支持真实数据库测试，使用 `sql/schema.sql` 自动创建测试数据库结构。

### 1. 真实数据库测试

使用 `WithDatabaseTest` 进行独立的数据库测试：

```go
func TestDatabaseOperation(t *testing.T) {
    testhelpers.WithDatabaseTest(t, func(db *testhelpers.DatabaseTestHelper) {
        // 开始事务
        db.BeginTransaction(t)
        defer db.RollbackTransaction(t)

        // 执行数据库操作
        db.ExecuteSQL(t, "INSERT INTO br_user (user_name, email, ...) VALUES (?, ?, ...)",
                     "testuser", "test@example.com")

        // 查询和验证
        row := db.QueryRow(t, "SELECT user_name FROM br_user WHERE email = ?", "test@example.com")
        var userName string
        row.Scan(&userName)

        assert.Equal(t, "testuser", userName)
    })
}
```

### 2. 集成测试中使用数据库

在测试套件中启用真实数据库：

```go
func TestWithRealDatabase(t *testing.T) {
    // 方法1：设置环境变量
    os.Setenv("USE_REAL_DB", "true")
    defer os.Unsetenv("USE_REAL_DB")

    testhelpers.WithCleanup(t, func(suite *testhelpers.TestSuite) {
        // suite.Database 现在是真实数据库连接
        db := suite.Database

        // 使用数据库进行测试
        db.BeginTransaction(t)
        defer db.RollbackTransaction(t)

        // 测试逻辑...
    })
}

// 方法2：在测试名称中包含 "Database"
func TestDatabaseIntegration(t *testing.T) {
    // 自动检测到 "Database" 关键字，使用真实数据库
    testhelpers.WithCleanup(t, func(suite *testhelpers.TestSuite) {
        // suite.Database 是真实数据库连接
    })
}
```

### 3. 数据库助手方法

```go
// 事务管理
db.BeginTransaction(t)
db.CommitTransaction(t)
db.RollbackTransaction(t)

// SQL执行
result := db.ExecuteSQL(t, "INSERT INTO ...", args...)
row := db.QueryRow(t, "SELECT ...", args...)
rows := db.Query(t, "SELECT ...", args...)

// 数据验证
count := db.CountRows(t, "table_name")
db.AssertRowCount(t, "table_name", expectedCount)
db.TruncateTable(t, "table_name")

// 直接访问数据库连接
sqlDB := db.GetDB()
```

### 4. 测试数据库管理

- **单一测试数据库**: 所有测试共享一个测试数据库（`hebe_test`）
- **Schema加载**: 自动从 `sql/schema.sql` 加载数据库结构（仅在首次使用时）
- **事务隔离**: 每个测试都在独立事务中运行，确保测试数据不互相影响
- **自动回滚**: 测试结束后自动回滚事务，无需手动清理数据

优势：

- **性能优化**: 避免重复创建和删除数据库
- **资源节约**: 减少数据库连接和磁盘操作
- **测试隔离**: 通过事务确保测试间完全隔离

### 5. 数据库配置

测试数据库使用以下默认配置（从.env.local读取）：

- Host: 127.0.0.1
- Port: 3306
- User: root
- Password: 123123
- Database: hebe_test（固定名称，所有测试共享）

如需修改，在测试中设置相应的环境变量：

```go
os.Setenv("IDIARY_DB_HOST", "custom_host")
os.Setenv("IDIARY_DB_USER", "custom_user")
os.Setenv("IDIARY_DB_PASSWORD", "custom_password")
// 然后调用测试
```

### 6. 数据清理策略

```go
// 自动事务回滚（推荐）
func TestWithTransaction(t *testing.T) {
    testhelpers.WithDatabaseTest(t, func(db *testhelpers.DatabaseTestHelper) {
        db.BeginTransaction(t)
        defer db.RollbackTransaction(t) // 自动清理

        // 执行测试...
    })
}

// 手动清理所有表（慎用）
func TestWithCleanup(t *testing.T) {
    testhelpers.WithDatabaseTest(t, func(db *testhelpers.DatabaseTestHelper) {
        defer db.CleanupAllTables(t) // 清空所有表

        // 执行测试...
    })
}
```
