# 测试架构说明

## 目录结构

```
tests/
├── unit/                     # 单元测试
│   ├── pkg/                 # pkg包测试
│   ├── internal/            # internal包测试
│   │   ├── configs/         # 配置模块测试
│   │   └── errors/          # 错误模块测试
│   └── modules/             # 业务模块测试
│       ├── user/            # 用户模块测试
│       ├── diary/           # 日记模块测试
│       ├── post/            # 日记内容模块测试
│       └── task/            # 任务模块测试
├── integration/             # 集成测试
│   ├── api/                 # API集成测试
│   └── database/            # 数据库集成测试
├── e2e/                     # 端到端测试
├── testdata/                # 测试数据
│   ├── fixtures/            # 测试固定数据
│   └── mock/                # 模拟数据
└── helpers/                 # 测试辅助工具
    ├── database.go          # 数据库测试工具
    └── http.go              # HTTP测试工具
```

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
```

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
