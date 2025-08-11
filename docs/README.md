# iDiary.ai Hebe 项目文档

## 📖 项目概述

iDiary.ai Hebe 是一个基于 Go 1.24 和 Gin 框架的现代化智能日记管理系统，采用领域驱动设计（DDD）的分层架构，通过泛型依赖注入容器实现松耦合设计，具备良好的可维护性和可扩展性。

### 核心特性

- **智能日记管理**: 支持多种内容类型（文字、图片、视频、链接）
- **多云盘存储**: 用户可将日记存储到个人云盘
- **现代化架构**: 基于 Go 1.24 泛型特性的零反射 DI 容器
- **国际化支持**: 完整的多语言支持系统（Kiko插件 + Website应用）
- **响应式设计**: 基于 Tailwind CSS 4.x 的现代 UI
- **高性能**: 优化的数据库查询和缓存策略
- **完整测试**: 单元测试、集成测试、性能测试框架

## 🏗️ 系统架构

### 分层架构概览

```
┌─────────────────────────────────────────────────────────────────┐
│                        表现层 (Website)                         │
├─────────────────────────────────────────────────────────────────┤
│   视图模板     │   静态资源     │   Kiko框架     │   样式文件      │
│  (Views)      │  (Static)     │ (JavaScript)  │ (Tailwind CSS) │
└─────────────────────────────────────────────────────────────────┘
                                   │
┌─────────────────────────────────────────────────────────────────┐
│                        应用层 (App)                             │
├─────────────────────────────────────────────────────────────────┤
│   启动器       │   路由管理     │   视图服务     │   控制器        │
│  (Startup)    │   (Route)     │   (View)      │ (Controller)   │
│               │               │   模板函数     │   中间件        │
│               │               │ (Templates)   │ (Middleware)   │
└─────────────────────────────────────────────────────────────────┘
                                   │
┌─────────────────────────────────────────────────────────────────┐
│                        领域层 (Domain)                          │
├─────────────────────────────────────────────────────────────────┤
│   用户模块     │   日记模块     │   文章模块     │   任务模块      │
│  (User)       │  (Diary)      │  (Post)       │  (Task)        │
│  Service      │  Service      │  Service      │  Service       │
│  Repository   │  Repository   │  Repository   │  Repository    │
│  Types        │  Types        │  Types        │  Types         │
└─────────────────────────────────────────────────────────────────┘
                                   │
┌─────────────────────────────────────────────────────────────────┐
│                      基础设施层 (Infra)                         │
├─────────────────────────────────────────────────────────────────┤
│                        仓储实现 (Repository)                     │
│   用户仓储     │   日记仓储     │   文章仓储     │   任务仓储      │
└─────────────────────────────────────────────────────────────────┘
                                   │
┌─────────────────────────────────────────────────────────────────┐
│                         核心层 (Core)                            │
├─────────────────────────────────────────────────────────────────┤
│   数据库       │   依赖注入     │   初始化器     │                │
│  (Database)   │    (DI)       │ (Initializer) │                │
└─────────────────────────────────────────────────────────────────┘
                                   │
┌─────────────────────────────────────────────────────────────────┐
│                      工具和契约层                               │
├─────────────────────────────────────────────────────────────────┤
│   契约接口     │   配置管理     │   错误处理     │   国际化       │
│ (Contracts)   │  (Config)     │  (Error)      │   (I18n)       │
│   中间件       │   通用工具     │                │                │
│ (Middleware)  │    (Pkg)      │                │                │
└─────────────────────────────────────────────────────────────────┘
```

### 架构层次说明

- **Website Layer**: 前端展示层（模板、静态资源、Kiko框架、Tailwind CSS）
- **App Layer**: 应用服务层（启动器、路由、视图服务、控制器、模板函数、中间件）
- **Domain Layer**: 业务领域层（用户、日记、文章、任务等业务模块，包含服务、仓储接口、类型定义）
- **Infra Layer**: 基础设施层（仓储实现，数据访问具体实现）
- **Core Layer**: 核心基础设施（数据库、DI容器、初始化器）
- **Utils & Contracts**: 工具和契约层（契约接口、配置、错误处理、国际化、中间件、通用工具）

## 📖 核心文档

### 系统架构与设计

- [系统架构概览](./architecture-overview.md) - 完整的系统分层架构和设计原理
- [系统启动过程详解](./system-startup-guide.md) - 详细的启动流程和组件初始化
- [存储架构设计](./storage-architecture-design.md) - 多云盘存储架构设计方案（未来规划）

### 开发指南

- [Go 1.24泛型DI容器使用指南](./di-container-guide.md) - 现代化高性能DI容器完整指南
- [中间件系统详解](./middleware-guide.md) - 完整的中间件架构和使用指南
- [模板函数使用指南](./template-functions-guide.md) - 丰富的模板函数和实际使用示例
- [模板开发最佳实践](./template-usage-guide.md) - 模板开发的架构模式和最佳实践
- [错误处理和响应服务示例](./error-handling-examples.md) - 统一错误处理体系的使用方法
- [国际化系统指南](./internationalization-guide.md) - 多语言支持和智能语言检测


### 系统功能

- [日志系统增强功能指南](./logging-enhancement-guide.md) - 日志系统的高级功能和使用方式

## 🏗️ 架构特点

### 分层架构

- **App Layer**: 应用编排、HTTP处理、渲染服务
- **Domain Layer**: 业务逻辑实现（User、Diary、Post、Task）
- **Core Layer**: 基础设施服务（Config、DI、Database、Error、Middleware）
- **Pkg Layer**: 纯函数工具库（通用工具和辅助功能）

### 核心技术栈

**后端技术**:
- **Go 1.24**: 使用最新泛型特性、embed包、Go Modules、Go Work
- **Gin 1.10.0**: 现代Web框架，支持Context、中间件、路由组、数据绑定
- **MySQL**: 主数据库存储
- **Redis**: 会话管理和缓存
- **SQLx**: 数据库操作增强

**前端技术**:
- **JavaScript ES2023**: 现代JS特性（async/await、箭头函数、模板字符串等）
- **jQuery 4.x**: 现代API支持
- **Tailwind CSS 4.x**: 现代CSS框架
- **Kiko框架**: 自研前端组件库

**基础设施**:
- **Go 1.24泛型DI容器**: 零反射、编译时类型安全
- **slog结构化日志**: 现代日志系统
- **Viper配置管理**: 灵活的配置系统
- **完整测试框架**: 单元测试、集成测试、性能测试

### 设计原则

- **DDD分层架构**: 清晰的领域驱动设计，依赖方向严格控制
- **泛型依赖注入**: 类型安全的零反射DI容器
- **统一错误处理**: 4位错误码体系，AppError结构
- **模板函数系统**: 安全的模板数据访问和国际化支持
- **中间件分层**: 核心中间件（Panic恢复、错误处理）+ 业务中间件
- **仓储模式**: 接口与实现分离，符合依赖倒置原则

## 🚀 快速导航

### 新手入门

1. **架构理解**: 查看 [系统架构概览](./architecture-overview.md) 了解DDD分层设计
2. **启动流程**: 阅读 [系统启动过程详解](./system-startup-guide.md) 理解应用启动
3. **依赖注入**: 学习 [Go 1.24泛型DI容器使用指南](./di-container-guide.md) 掌握现代DI
4. **中间件**: 参考 [中间件系统详解](./middleware-guide.md) 学习请求处理

### 开发人员

**前端开发**:
- [模板函数使用指南](./template-functions-guide.md) - 模板数据访问和渲染
- [国际化系统指南](./internationalization-guide.md) - 多语言支持实现
- [模板翻译指南](./template-translation-guide.md) - 模板国际化最佳实践

**后端开发**:
- [Go 1.24泛型DI容器使用指南](./di-container-guide.md) - 现代依赖注入
- [错误处理和响应服务示例](./error-handling-examples.md) - 统一错误处理
- [存储架构设计](./storage-architecture-design.md) - 多云盘存储方案

**测试开发**:
- [测试框架指南](./testing-framework-guide.md) - 完整测试体系

**运维人员**:
- [日志系统增强功能指南](./logging-enhancement-guide.md) - 结构化日志管理

### 架构师和技术负责人

1. **架构设计**: 深入 [系统架构概览](./architecture-overview.md) 了解DDD设计模式
2. **依赖分析**: 查看 [依赖分析报告](./dependency-analysis-report.md) 了解模块关系
3. **存储设计**: 研究 [存储架构设计](./storage-architecture-design.md) 了解多云盘方案

## 📚 完整文档列表

### 🏗️ 架构设计

- [系统架构概览](./architecture-overview.md) - DDD分层架构设计和模块说明
- [依赖分析报告](./dependency-analysis-report.md) - 模块依赖关系深度分析
- [存储架构设计](./storage-architecture-design.md) - 多云盘存储架构方案

### 🔧 核心组件

- [Go 1.24泛型DI容器使用指南](./di-container-guide.md) - 零反射依赖注入容器
- [中间件系统详解](./middleware-guide.md) - 中间件架构和最佳实践
- [错误处理和响应服务示例](./error-handling-examples.md) - 统一错误处理机制
- [日志系统增强功能指南](./logging-enhancement-guide.md) - slog结构化日志系统

### 💻 开发指南

- [系统启动过程详解](./system-startup-guide.md) - 应用启动流程和组件初始化
- [模板函数使用指南](./template-functions-guide.md) - 模板函数系统和数据访问
- [模板使用指南](./template-usage-guide.md) - 模板开发最佳实践
- [模板翻译指南](./template-translation-guide.md) - 模板国际化实现
- [国际化系统指南](./internationalization-guide.md) - 多语言支持完整方案
- [测试框架指南](./testing-framework-guide.md) - 单元测试、集成测试、性能测试
- [配置系统文档](../internal/core/config/README.md) - 配置管理系统
- [Kiko UI框架文档](../kiko/README.md) - 前端UI组件框架
## 🔧 开发环境要求

### 系统要求

- **Go**: 1.24+ (支持泛型、embed等现代特性)
- **Node.js**: 18+ (ES2023支持)
- **MySQL**: 8.0+
- **Redis**: 6.0+

### 浏览器支持

- Chrome 110+
- Edge 110+
- Firefox 128 ESR+

### 快速开始

```bash
# 克隆项目
git clone <repository-url>
cd hebe

# 安装Go依赖
go mod download

# 安装前端依赖
npm install

# 配置环境变量
cp .env.example .env.local
# 编辑 .env.local 配置数据库、Redis等信息

# 构建前端资源
npm run build

# 运行项目
go run cmd/idiary/main.go
```

## 🧪 测试

```bash
# 运行所有测试
go test ./...

# 运行单元测试
go test ./tests/unit/...

# 运行集成测试
go test ./tests/integration/...

# 运行性能测试
go test ./tests/benchmark/...

# 生成测试覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🔧 开发规范

### 架构规范

- **依赖方向**: Website → App → Domain → Infra → Core → Utils (严禁反向依赖)
- **Domain层**: 纯业务逻辑，Service-Repository-Types模式
- **Infra层**: 仓储实现，数据访问具体实现
- **Core层**: 基础设施服务（数据库、DI、初始化器）

### 代码规范

- **Go**: 遵循 golangci-lint 规范，最大函数长度60行，最大参数5个
- **JavaScript**: 遵循 ESLint + Prettier 规范，最大函数长度50行，最大参数4个
- **错误处理**: 使用4位错误码体系和AppError结构
- **模板开发**: 用户信息通过模板函数访问，业务数据通过`.Data`传递

## 📝 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 遵循代码规范（ESLint、golangci-lint）
4. 编写测试用例
5. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
6. 推送到分支 (`git push origin feature/AmazingFeature`)
7. 打开 Pull Request

### 提交规范

- 使用语义化提交规范
- 新功能开发时同步编写文档
- 保持测试覆盖率

## 🎯 项目特色

### 🚀 高性能DI容器

- **零反射设计**: 基于Go 1.24泛型特性
- **编译时类型安全**: 避免运行时错误
- **高并发性能**: ~57纳秒解析速度
- **三种生命周期**: 瞬态/单例/实例

### 🎨 现代化前端

- **Kiko框架**: 自研轻量级UI组件库
- **ES2023特性**: 现代JavaScript语法支持
- **Tailwind CSS 4.x**: 现代CSS框架
- **完整国际化**: 插件级和应用级多语言支持

### 🏗️ DDD架构

- **分层清晰**: 严格的依赖方向控制
- **领域驱动**: 业务逻辑内聚，技术关注点分离
- **仓储模式**: 接口与实现分离
- **统一错误处理**: 4位错误码体系

## 🔍 常见问题

### 快速排查

- **启动失败**: 查看 [系统启动过程详解](./system-startup-guide.md)
- **错误处理**: 查看 [错误处理和响应服务示例](./error-handling-examples.md)
- **模板问题**: 查看 [模板函数使用指南](./template-functions-guide.md)
- **中间件问题**: 查看 [中间件系统详解](./middleware-guide.md)
- **DI容器**: 查看 [Go 1.24泛型DI容器使用指南](./di-container-guide.md)

### 调试技巧

- 查看slog结构化日志输出
- 使用DI容器的类型安全特性
- 检查中间件执行顺序
- 验证错误码和AppError响应格式

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](../LICENSE) 文件了解详情。

## 🤝 支持与反馈

如果您在使用过程中遇到问题，请：

1. 查看相关文档和示例
2. 搜索已有的 Issues
3. 创建新的 Issue 详细描述问题
4. 参与社区讨论

---

*文档版本: v2.0 | 最后更新: 2024年 | 基于最新代码重构更新*
