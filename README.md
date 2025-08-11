# iDiary.AI - 智能日记平台

> **记录人生，伴你前行** | Record Life, Accompany You Forward

## 项目概述

iDiary.AI 是一个基于 **Go 1.24** 和 **Gin** 框架的现代化智能日记管理系统，采用**领域驱动设计（DDD）**的分层架构，通过**泛型依赖注入容器**实现松耦合设计，具备良好的可维护性和可扩展性。

### 🚀 核心特性

- **🤖 智能日记管理**: 支持多种内容类型（文字、图片、视频、链接）
- **☁️ 多云盘存储**: 用户可将日记存储到个人云盘
- **🏗️ 现代化架构**: 基于 Go 1.24 泛型特性的零反射 DI 容器
- **🌍 国际化支持**: 完整的多语言支持系统（Kiko插件 + Website应用）
- **🎨 响应式设计**: 基于 Tailwind CSS 4.x 的现代 UI
- **⚡ 高性能**: 优化的数据库查询和缓存策略
- **🧪 完整测试**: 单元测试、集成测试、性能测试框架

### 🛠️ 技术栈

**后端技术**:
- **Go 1.24**: 泛型特性、embed包、Go Modules
- **Gin 1.10.0**: 现代Web框架
- **MySQL 8.0+**: 主数据库存储
- **Redis**: 会话管理和缓存
- **SQLx**: 数据库操作增强

**前端技术**:
- **JavaScript ES2023**: 现代JS特性
- **jQuery 4.x**: 现代API支持
- **Tailwind CSS 4.x**: 现代CSS框架
- **Kiko框架**: 自研前端组件库

**基础设施**:
- **Go 1.24泛型DI容器**: 零反射、编译时类型安全
- **slog结构化日志**: 现代日志系统
- **Viper配置管理**: 灵活的配置系统

## 🌍 多语言支持

项目支持以下语言的 slogan "记录人生，伴你前行"：

- **中文**: 记录人生，伴你前行
- **English**: Record Life, Accompany You Forward
- **Deutsch**: Das Leben aufzeichnen, Dich begleiten
- **Español**: Registrar la vida, Acompañarte en el camino
- **Français**: Enregistrer la vie, Vous accompagner dans votre parcours
- **日本語**: 人生を記録し、あなたと共に前進する
- **한국어**: 인생 기록, 함께 나아가다
- **Русский**: Записывайте жизнь, Сопровождайте вас вперед
- **العربية**: تسجيل الحياة، رفقتك في المضي قدمًا

## 🤖 AI 增强功能（规划中）

基于 "记录人生，伴你前行" 的理念，计划集成以下AI功能：

- **情感分析**: 分析日记情感，提供情绪趋势
- **智能标签**: 自动为日记添加分类标签
- **智能回忆**: 个性化的历史回顾推送
- **情感支持**: AI驱动的情绪支持和建议
- **智能摘要**: 自动生成日记摘要和时间线
- **语音识别**: 语音输入日记内容
- **图像识别**: 自动识别和关联照片内容

## 🚀 快速开始

### 环境要求

- **Go**: 1.24+ (支持泛型、embed等现代特性)
- **Node.js**: 18+ (ES2023支持)
- **MySQL**: 8.0+
- **Redis**: 6.0+ (可选)

### 安装和运行

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

### 构建脚本

```bash
# 编译 CSS 文件
npm run build:css

# 编译 JavaScript 文件
npm run build:js

# 同时构建 CSS 和 JS
npm run build

# 构建 Go 程序
go build -o bin/idiary ./cmd/idiary/main.go
```

## 📚 文档导航

### 🏗️ 架构设计
- [📖 完整项目文档](./docs/README.md) - 详细的技术文档和开发指南
- [🏗️ 系统架构概览](./docs/architecture-overview.md) - DDD分层架构设计
- [📊 依赖分析报告](./docs/dependency-analysis-report.md) - 模块依赖关系分析

### 🔧 开发指南
- [🚀 Go 1.24泛型DI容器](./docs/di-container-guide.md) - 现代化依赖注入
- [⚙️ 中间件系统详解](./docs/middleware-guide.md) - 中间件架构和使用
- [🎨 模板函数使用指南](./docs/template-functions-guide.md) - 模板开发最佳实践
- [🌍 国际化系统指南](./docs/internationalization-guide.md) - 多语言支持方案
- [🧪 测试框架指南](./docs/testing-framework-guide.md) - 完整测试体系

### 🎨 前端框架
- [Kiko UI框架文档](./kiko/README.md) - 自研前端组件库

## 📁 项目结构

```plaintext
.
├── cmd/idiary/            # 🚀 主程序入口
├── internal/              # 🏗️ 内部应用逻辑
│   ├── app/               # 📱 应用层（控制器、路由、视图、中间件）
│   ├── domain/            # 🏢 业务领域层（用户、日记、文章、任务）
│   ├── infra/             # 🔧 基础设施层（仓储实现）
│   ├── core/              # ⚙️ 核心层（数据库、DI、初始化器）
│   ├── contracts/         # 📋 契约接口
│   └── utils/             # 🛠️ 工具和辅助功能
├── website/               # 🎨 前端资源
│   ├── views/             # 📄 HTML模板
│   ├── static/            # 📦 静态资源
│   └── locales/           # 🌍 多语言文件
├── kiko/                  # 🎯 Kiko UI框架
├── tests/                 # 🧪 测试代码
│   ├── unit/              # 单元测试
│   ├── integration/       # 集成测试
│   └── benchmark/         # 性能测试
├── docs/                  # 📚 项目文档
├── scripts/               # 📜 构建和部署脚本
├── sql/                   # 🗄️ 数据库脚本
└── pkg/                   # 📦 可重用工具包
```

## 🤝 贡献指南

我们欢迎所有形式的贡献！请查看 [贡献指南](./docs/CONTRIBUTING.md) 了解详细信息。

### 开发流程

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 📞 联系我们

- 项目主页：[iDiary.ai](https://idiary.ai)
- 问题反馈：[GitHub Issues](https://github.com/idiary-ai/hebe/issues)
- 邮箱：support@idiary.ai

---

<div align="center">
  <p>🌟 如果这个项目对你有帮助，请给我们一个 Star！</p>
  <p>Made with ❤️ by iDiary.ai Team</p>
</div>
