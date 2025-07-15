# iDiary.AI 项目文档中心

欢迎来到iDiary.AI项目文档中心！本文档中心包含了系统架构、开发指南、使用手册等全面的技术文档。

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

### 系统功能

- [日志系统增强功能指南](./logging-enhancement-guide.md) - 日志系统的高级功能和使用方式

## 🏗️ 架构特点

### 分层架构

- **App Layer**: 应用编排、HTTP处理、渲染服务
- **Module Layer**: 业务逻辑实现（User、Diary、Post、Task）
- **Core Layer**: 基础设施服务（Config、DI、Database、Error、Middleware）
- **Pkg Layer**: 纯函数工具库

### 核心技术栈

- **后端**: Go + Gin框架
- **数据库**: MySQL
- **缓存**: Redis (Session管理)
- **依赖注入**: Go 1.24泛型DI容器 (零反射)
- **日志**: slog结构化日志
- **模板**: HTML模板引擎

### 设计原则

- **DDD分层架构**: 清晰的领域驱动设计
- **依赖注入**: 类型安全的泛型DI容器
- **统一错误处理**: 4位错误码体系
- **模板函数**: 安全的模板数据访问
- **中间件分层**: 核心中间件 + 业务中间件

## 🚀 快速导航

### 新手入门

1. 查看 [系统架构概览](./architecture-overview.md) 了解整体设计
2. 阅读 [系统启动过程详解](./system-startup-guide.md) 理解启动流程
3. 学习 [Go 1.24泛型DI容器使用指南](./di-container-guide.md) 掌握现代DI容器
4. 参考 [中间件系统详解](./middleware-guide.md) 学习中间件使用

### 开发人员

1. **前端开发**: 查看 [模板函数使用指南](./template-functions-guide.md)
2. **后端开发**: 查看 [Go 1.24泛型DI容器使用指南](./di-container-guide.md) 和 [错误处理和响应服务示例](./error-handling-examples.md)
3. **运维人员**: 查看 [日志系统增强功能指南](./logging-enhancement-guide.md)

### 架构师

1. 深入 [系统架构概览](./architecture-overview.md) 了解设计模式
2. 查看 [存储架构设计](./storage-architecture-design.md) 了解扩展规划
3. 研究各模块的详细实现

## 📚 模块文档

以下模块包含独立的README文档：

- [DI容器文档](../internal/core/di/README.md) - 依赖注入容器详细文档
- [配置系统文档](../internal/core/config/README.md) - 配置管理系统
- [Kiko UI框架文档](../kiko/README.md) - 前端UI组件框架
- [测试文档](../tests/README.md) - 测试框架和策略

## 🔧 开发规范

### 代码结构规范

- **App层**: 只依赖Module和Core层
- **Module层**: 采用Handler-Service-Repository模式
- **Core层**: 提供基础设施服务
- **Pkg层**: 纯函数，无依赖

### 错误处理规范

- 使用4位错误码体系
- 通过AppError统一错误格式
- 中间件统一错误响应处理

### 模板开发规范

- 用户信息通过模板函数访问
- 业务数据通过`.Data`传递
- 所有文本使用国际化函数

### 中间件开发规范

- 核心中间件在Core层
- 业务中间件在Module层
- 遵循执行顺序最佳实践

## 📝 文档维护

### 文档分类

- **架构文档**: `xxx-overview.md`、`xxx-architecture-design.md`
- **指南文档**: `xxx-guide.md`
- **示例文档**: `xxx-examples.md`
- **模块文档**: 各模块的`README.md`

### 更新原则

- 重大功能更新后同步更新文档
- 保持文档结构清晰统一
- 及时清理过时内容
- 添加实际使用示例

### 贡献文档

1. 新功能开发时同步编写文档
2. 发现文档问题及时反馈
3. 遵循统一的文档格式规范
4. 提供清晰的代码示例

## 🔍 问题排查

### 常见问题

- **启动失败**: 查看 [系统启动过程详解](./system-startup-guide.md)
- **错误处理**: 查看 [错误处理和响应服务示例](./error-handling-examples.md)
- **模板问题**: 查看 [模板函数使用指南](./template-functions-guide.md)
- **中间件问题**: 查看 [中间件系统详解](./middleware-guide.md)

### 调试技巧

- 查看结构化日志输出
- 使用DI容器的调试功能
- 检查中间件执行顺序
- 验证错误码和响应格式

## 🎯 项目特色

### 高性能DI容器

- 19.89纳秒解析速度
- 零内存分配
- 类型安全的泛型接口

### 统一渲染服务

- 整合response、view、template功能
- 统一的页面数据结构
- 丰富的模板函数支持

### 分层中间件

- 核心中间件处理系统级功能
- 业务中间件处理业务逻辑
- 清晰的职责分离

### 结构化日志

- slog库支持
- 自动错误追踪
- 性能监控集成

---

**文档版本**: v2.0  
**最后更新**: 2024年  
**维护团队**: iDiary开发团队

如有任何文档问题或建议，请通过项目Issue或讨论区反馈。
