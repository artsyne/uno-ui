# iDiary.ai Hebe - 智能日记平台

> **记录人生，伴你前行** | Record Life, Accompany You Forward

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Gin Framework](https://img.shields.io/badge/Gin-1.10.0-00ADD8?style=flat)](https://gin-gonic.com/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Test Coverage](https://img.shields.io/badge/Coverage-91%25+-brightgreen.svg)](#)

</div>

## 📖 项目概述

**iDiary.ai Hebe** 是一个基于 **Go 1.25** 和 **Gin** 框架的现代化智能日记管理系统，采用**领域驱动设计（DDD）**的分层架构，通过**Go 1.25泛型依赖注入容器**实现零反射、编译时类型安全的高性能设计。

### 🚀 核心特性

- **🤖 智能日记管理**: 支持多种内容类型（文字、图片、视频、链接）的日记创作和管理
- **☁️ 多云存储**: 支持 Dropbox、Google Drive、阿里云盘、WebDAV 等多种云存储服务
- **🏗️ 现代化架构**: 基于 DDD 分层设计，Go 1.25 泛型 DI 容器，零反射高性能
- **🌍 国际化支持**: 智能语言检测，预加载缓存，微秒级翻译查找
- **🎨 响应式设计**: Tailwind CSS 4.x + Kiko.js 组件库，现代化 UI/UX，支持 Light/Dark 主题切换
- **⚡ 高性能**: 优化的数据库查询、Redis 缓存、异步任务处理
- **🧪 完整测试**: 91%+ 覆盖率，真实环境测试框架，性能基准测试
- **🔒 统一错误处理**: 4位错误码体系，结构化日志，优雅的错误降级

### 🛠️ 技术栈详解

**🔧 后端技术**:

- **Go 1.25**: 泛型特性、embed包、现代并发模式
- **Gin**: 高性能 Web 框架，中间件链式设计
- **MySQL 8.0+**: 主数据库存储，优化查询性能
- **Redis**: 会话管理和高速缓存
- **SQLx**: 增强型数据库操作，类型安全
- **gocron**: 异步任务调度（邮件任务、定时清理）

**🎨 前端技术**:

- **Tailwind CSS 4.x**: 原子化 CSS 框架，响应式设计
- **jquery 4.x**: 前端库，简化 DOM 操作
- **Kiko.js**: 自研前端组件库 (Grid、Tab、Notify、Valid 等)
- **JavaScript ES2023**: 现代 JS 特性，模块化设计
- **国际化 i18n**: 多语言支持，智能语言切换

**🏗️ 架构特性**:

- **DDD 分层架构**: Website、App、Domain、Infra、Core、Utils 层次分明
- **Go 1.25 泛型 DI**: 零反射、编译时类型安全、~30ns 解析速度
- **ViewService 模式**: 统一视图渲染、数据绑定、响应格式化
- **中间件链**: Panic 恢复、错误处理、用户认证、数据注入
- **slog 结构化日志**: 分级日志、结构化输出、性能优化
- **Viper 配置管理**: 环境变量、配置文件、动态加载

## 🌍 智能国际化系统

**iDiary.ai Hebe** 支持智能语言检测和多语言 slogan "记录人生，伴你前行"：

| 语言         | Slogan                                      | 状态        |
| ------------ | ------------------------------------------- | ----------- |
| **中文**     | 记录人生，伴你前行                          | ✅ 完整支持 |
| **English**  | Record Life, Accompany You Forward          | ✅ 完整支持 |
| **Deutsch**  | Das Leben aufzeichnen, Dich begleiten       | 🔄 部分支持 |
| **Español**  | Registrar la vida, Acompañarte en el camino | 🔄 部分支持 |
| **Français** | Enregistrer la vie, Vous accompagner        | 🔄 部分支持 |
| **日本語**   | 人生を記録し、あなたと共に前進する          | 🔄 部分支持 |
| **한국어**   | 인생 기록, 함께 나아가다                    | 🔄 部分支持 |
| **Русский**  | Записывайте жизнь, Сопровождайте вас        | 🔄 部分支持 |

### 国际化特性

- **🧠 智能检测**: 用户设置 → Cookie → 浏览器 → 默认语言
- **⚡ 高性能**: 预加载所有语言包，微秒级翻译查找
- **🔄 热切换**: 无需刷新页面即可切换语言
- **📱 自适应**: 支持 RTL 语言和本地化时间格式

## 🏗️ 系统架构

### DDD 分层架构

```
hebe/
├── cmd/                        # 命令行入口
│   ├── idiary/                # 主应用程序
│   └── cleandb/               # 数据库清理工具
├── internal/                  # 内部应用代码
│   ├── app/                   # 应用层 (Application Layer)
│   │   ├── controller/        # HTTP控制器
│   │   ├── middleware/        # 业务中间件
│   │   ├── templates/         # 模板函数
│   │   ├── view/              # ViewService视图服务
│   │   ├── route/             # 路由配置
│   │   └── startup/           # 应用启动管理
│   ├── domain/                # 领域层 (Domain Layer)
│   │   ├── user/              # 用户领域
│   │   ├── diary/             # 日记领域
│   │   ├── post/              # 文章领域
│   │   └── task/              # 任务领域
│   ├── infra/                 # 基础设施层 (Infrastructure Layer)
│   │   └── repository/        # 数据仓储实现
│   ├── core/                  # 核心层 (Core Layer)
│   │   ├── database/          # 数据库连接管理
│   │   ├── di/                # Go 1.25泛型DI容器
│   │   └── initializer/       # 核心初始化器
│   ├── utils/                 # 工具层 (Utils Layer)
│   │   ├── config/            # 配置管理
│   │   ├── error/             # 错误处理
│   │   ├── i18n/              # 国际化服务
│   │   └── middleware/        # 系统级中间件
│   └── contracts/             # 合约层 (Contracts Layer)
│       └── structs.go         # 共享数据结构
├── pkg/                       # 公共工具包
├── website/                   # 前端资源
│   ├── static/                # 静态资源 (CSS, JS, 图片)
│   ├── views/                 # HTML模板
│   └── locales/               # 国际化翻译文件
├── kiko/                      # Kiko.js 前端组件库
├── tests/                     # 测试代码
│   ├── framework/             # 测试框架
│   ├── unit/                  # 单元测试
│   ├── integration/           # 集成测试
│   └── benchmark/             # 性能测试
└── docs/                      # 项目文档
```

### 🎯 核心设计原则

- **🎯 单一职责**: 每个层次和组件都有明确的职责边界
- **🔄 依赖倒置**: 高层模块不依赖低层模块，通过接口实现解耦
- **📈 开闭原则**: 对扩展开放，对修改封闭
- **⚙️ 泛型DI**: 零反射、编译时类型安全的依赖注入
- **🔄 事件驱动**: 异步任务处理和系统解耦（gocron）
- **🧪 数据隔离**: 测试环境完全隔离，支持并发测试

### 📊 架构层次分析

```
依赖方向: Website → App → Domain → Infra → Core → Utils
严格禁止: 反向依赖和跨层直接依赖

层次职责:
├── Website Layer: 表现层（模板、静态资源、UI组件）
├── App Layer: 应用层（控制器、视图服务、路由、中间件）
├── Domain Layer: 领域层（业务逻辑、领域服务、实体）
├── Infra Layer: 基础设施层（仓储实现、数据访问）
├── Core Layer: 核心层（数据库、DI容器、初始化）
└── Utils Layer: 工具层（配置、错误处理、国际化、中间件）
```

## 🚀 快速开始

### 🔧 环境要求

| 组件        | 版本要求 | 用途                      |
| ----------- | -------- | ------------------------- |
| **Go**      | 1.25+    | 支持泛型、embed等现代特性 |
| **MySQL**   | 8.0+     | 主数据库存储              |
| **Redis**   | 6.0+     | 会话管理和高速缓存        |
| **Node.js** | 18+      | 前端构建工具（可选）      |

### ⚡ 快速启动

```bash
# 1. 克隆项目
git clone https://github.com/biuren/hebe.git
cd hebe

# 2. 安装依赖
go mod download

# 3. 配置数据库
# 创建MySQL数据库: CREATE DATABASE idiary_hebe;
# 导入数据结构: mysql < sql/schema.sql

# 4. 配置环境变量
# 复制配置文件并修改数据库连接信息

# 5. 启动应用
go run cmd/idiary/main.go

# 6. 访问应用 🌐
# http://localhost:8080
```

### 🛠️ 开发命令

```bash
# 开发服务器（热重载）
./scripts/run_debug.sh

# 完整测试套件
go test ./tests/...

# 性能基准测试
go test ./tests/benchmark/... -bench=. -benchmem

# 测试覆盖率报告
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

# 数据库清理工具
go run cmd/cleandb/main.go

# 生产构建
go build -ldflags="-s -w" -o bin/idiary ./cmd/idiary/main.go
```

### 🐳 Docker 部署（推荐）

```bash
# 一键启动完整环境（MySQL + Redis + App）
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f idiary

# 访问应用: http://localhost:8080
```

## 📚 完整文档导航

### 🌟 新手入门路径

1. **📖 [项目文档总览](./docs/README.md)** - 开始这里！完整技术文档导航
2. **🏗️ [系统架构概览](./docs/architecture-overview.md)** - 理解DDD分层设计
3. **🚀 [系统启动流程](./docs/system-startup-guide.md)** - 掌握应用启动过程
4. **⚙️ [Go 1.25泛型DI容器](./docs/di-container-guide.md)** - 学习现代依赖注入

### 🔧 核心技术文档

| 技术领域     | 文档链接                                               | 描述                      |
| ------------ | ------------------------------------------------------ | ------------------------- |
| **依赖注入** | [Go 1.25泛型DI容器](./docs/di-container-guide.md)      | 零反射、高性能DI系统      |
| **国际化**   | [国际化系统指南](./docs/i18n-guide.md) | 智能语言检测和翻译        |
| **主题系统** | [主题系统指南](./docs/THEME_SYSTEM_GUIDE.md)           | Light/Dark主题切换系统    |
| **中间件**   | [中间件系统详解](./docs/middleware-guide.md)           | 分层中间件架构            |
| **错误处理** | [错误处理机制](./docs/error-handling-examples.md)      | 统一错误处理和ViewService |
| **日志系统** | [日志系统指南](./docs/logging-enhancement-guide.md)    | slog结构化日志            |

### 🎨 前端开发文档

| 开发领域     | 文档链接                                             | 描述                   |
| ------------ | ---------------------------------------------------- | ---------------------- |
| **模板开发** | [模板函数指南](./docs/template-functions-guide.md)   | 模板函数系统和数据访问 |
| **最佳实践** | [模板开发指南](./docs/template-usage-guide.md)       | 模板架构模式和实践     |
| **国际化**   | [模板翻译指南](./docs/template-translation-guide.md) | 模板国际化实现         |

### 🧪 测试与架构文档

| 领域         | 文档链接                                              | 描述                   |
| ------------ | ----------------------------------------------------- | ---------------------- |
| **测试框架** | [测试框架指南](./docs/testing-framework-guide.md)     | 完整测试体系和数据隔离 |
| **架构分析** | [依赖分析报告](./docs/dependency-analysis-report.md)  | 模块依赖和架构评估     |
| **存储设计** | [存储架构设计](./docs/storage-architecture-design.md) | 多云存储方案           |

## 📊 性能指标

### 核心性能

```
启动性能:
├── 应用冷启动: < 2s
├── 核心初始化: < 100ms
├── 应用层初始化: < 200ms
└── 路由注册: < 50ms

运行时性能:
├── DI容器解析: ~30ns/op
├── 错误处理: < 1μs
├── 国际化翻译: < 1μs
├── 模板渲染: < 5ms
└── 数据库操作: < 10ms

内存使用:
├── 基础内存: ~15MB
├── DI容器: < 1MB
├── 国际化缓存: < 2MB
└── 模板缓存: < 5MB
```

### 测试覆盖率

```
测试统计:
├── 总体覆盖率: 91%+
├── 核心模块: 95%+
├── 领域层: 90%+
├── 应用层: 85%+
└── 测试文件数: 15+
```

### 代码质量

- **Go 1.25泛型**: 编译时类型安全，零运行时开销
- **零反射设计**: 所有核心组件避免反射，确保高性能
- **DDD架构**: 清晰的层次划分，高内聚低耦合
- **测试驱动**: 91%+覆盖率，真实环境测试

## 🛠️ 开发工具

```bash
# 代码质量检查
go vet ./...
golangci-lint run

# 性能分析
go test -bench=. -benchmem -cpuprofile=cpu.prof
go tool pprof cpu.prof

# 测试覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 依赖更新
go mod tidy
go mod download
```

## 🎯 技术亮点

### 🚀 Go 1.25泛型DI容器

- **零反射设计**: 基于泛型特性，编译时类型安全
- **高性能**: ~30ns解析速度，零运行时开销
- **三种生命周期**: Bind瞬态、Singleton单例、Value实例

### 🏗️ DDD分层架构

- **清晰分层**: Website→App→Domain→Infra→Core→Utils
- **依赖倒置**: 高层不依赖低层，接口与实现分离
- **业务内聚**: 领域驱动设计，技术关注点分离

### 🌍 智能国际化

- **预加载缓存**: 启动时加载所有语言文件到内存
- **智能检测**: 用户设置→Cookie→浏览器→默认
- **微秒级翻译**: 高效缓存map，翻译查找<1μs

### 🧪 真实环境测试

- **91%+覆盖率**: 高质量测试保障
- **数据隔离**: 前缀隔离策略，支持并发测试
- **性能基准**: 内置性能测试和回归检测

## 🤝 贡献指南

### 🛠️ 开发环境准备

```bash
# 1. 环境要求
# Go 1.25+, MySQL 8.0+, Redis 6.0+

# 2. 克隆项目
git clone https://github.com/biuren/hebe.git
cd hebe

# 3. 安装依赖
go mod download

# 4. 运行测试
go test ./tests/...

# 5. 启动开发服务器
./scripts/run_debug.sh
```

### 📝 参与流程

1. **Fork项目** 并创建特性分支
2. **编写代码** 并保持91%+测试覆盖率
3. **遵循规范** (golangci-lint, 语义化提交)
4. **更新文档** 同步变更相关文档
5. **提交PR** 并详细描述变更内容

### 🎯 贡献类型

| 类型            | 描述           | 优先级 |
| --------------- | -------------- | ------ |
| 🐛 **Bug修复**  | 发现并修复问题 | 高     |
| ✨ **新功能**   | 添加功能或改进 | 中     |
| 📚 **文档完善** | 改进文档和示例 | 中     |
| 🧪 **测试增强** | 提高测试覆盖率 | 中     |
| 🌍 **国际化**   | 新增语言支持   | 低     |
| 🎨 **UI优化**   | 前端界面改进   | 低     |

## 📄 许可证

本项目采用 **[MIT 许可证](LICENSE)** - 自由使用、修改和分发。

---

<div align="center">

### 🌟 如果这个项目对你有帮助，请给我们一个 Star！

## **记录人生，伴你前行**

## **Record Life, Accompany You Forward**

**Made with ❤️ by iDiary.ai Hebe Team**

[![GitHub stars](https://img.shields.io/github/stars/biuren/hebe.svg?style=social&label=Star)](https://github.com/biuren/hebe/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/biuren/hebe.svg?style=social&label=Fork)](https://github.com/biuren/hebe/network)

</div>
