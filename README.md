# iDiary.AI - 智能日记平台

## 架构概述

iDiary.AI 是一个基于 Go 语言开发的现代化智能日记平台，采用分层架构设计，支持多语言国际化，具备完善的渲染系统和用户管理功能。

## 系统架构

```
┌─────────────────────────────────────────────────────────────────┐
│                           前端层 (Frontend)                      │
├─────────────────────────────────────────────────────────────────┤
│  Kiko UI框架  │  模板引擎  │  静态资源  │  多语言支持              │
└─────────────────────────────────────────────────────────────────┘
                                   │
┌─────────────────────────────────────────────────────────────────┐
│                         应用层 (Application)                     │
├─────────────────────────────────────────────────────────────────┤
│           路由系统        │        中间件        │    控制器      │
│         (Router)         │     (Middleware)     │  (Controller)  │
└─────────────────────────────────────────────────────────────────┘
                                   │
┌─────────────────────────────────────────────────────────────────┐
│                         渲染层 (Render)                          │
├─────────────────────────────────────────────────────────────────┤
│  RenderService  │  PageBuilder  │  I18nManager  │  TemplateFuncs │
│       │         │       │       │       │       │       │        │
│  ResponseHandler│  ContextBuilder│  MessageService│  RenderUtils │
└─────────────────────────────────────────────────────────────────┘
                                   │
┌─────────────────────────────────────────────────────────────────┐
│                        业务层 (Business)                         │
├─────────────────────────────────────────────────────────────────┤
│   用户模块     │   日记模块     │   文章模块     │   任务模块      │
│  (User)       │  (Diary)      │  (Post)       │  (Task)        │
│  Service      │  Service      │  Service      │  Service       │
│  Handler      │  Handler      │  Handler      │  Handler       │
│  Repository   │  Repository   │  Repository   │  Repository    │
└─────────────────────────────────────────────────────────────────┘
                                   │
┌─────────────────────────────────────────────────────────────────┐
│                         核心层 (Core)                            │
├─────────────────────────────────────────────────────────────────┤
│   配置管理     │   数据库      │   依赖注入     │   错误处理      │
│  (Config)     │  (Database)   │    (DI)       │  (Error)       │
└─────────────────────────────────────────────────────────────────┘
                                   │
┌─────────────────────────────────────────────────────────────────┐
│                      基础设施层 (Infrastructure)                   │
├─────────────────────────────────────────────────────────────────┤
│    MySQL      │    Redis      │    文件存储     │    日志系统     │
│   数据库       │    缓存       │   (Storage)    │   (Logging)    │
└─────────────────────────────────────────────────────────────────┘
```

## 域名

iDiary.AI

## 日记网功能相关

### slogan

英文："Record Life, Accompany You Forward"
德文："Das Leben aufzeichnen, Dich begleiten"
西班牙文："Registrar la vida, Acompañarte en el camino"
法文："Enregistrer la vie, Vous accompagner dans votre parcours"
日文："人生を記録し、あなたと共に前進する"
韩文："인생 기록, 함께 나아가다"
俄文："Записывайте жизнь, Сопровождайте вас вперед"
中文（简体）："记录人生，伴你前行"
中文（繁体）："記錄人生，伴你前行"
阿拉伯文："تسجيل الحياة، رفقتك في المضي قدمًا"

### chatgpt4推荐的功能

以下是20个AI相关的功能建议，可以提升日记网的竞争力，并更好地体现 slogan "记录人生，伴你前行" 的精髓：

1. 情感分析：通过AI技术分析用户日记中的情感，并提供情绪指数和情感趋势分析，帮助用户更好地了解自己的情绪变化。

2. 智能标签：AI可以根据日记内容自动添加标签，方便用户对日记进行分类和搜索。

3. 语音识别：通过语音识别技术，用户可以通过语音输入日记内容，方便快捷地记录生活中的点滴。

4. 图像识别：AI可以识别用户上传的照片中的物体、人物等内容，并自动关联到相应的日记条目中，丰富日记的多样性。

5. 智能回忆：AI可以根据用户的历史记录和日记内容，提供个性化的回忆推送，让用户能够回顾过去的美好时光。

6. 智能日历：将用户的日记内容与日历相结合，AI可以提供个性化的日程安排和活动推荐，帮助用户更好地规划生活。

7. 自动摘要：AI可以根据日记内容生成自动摘要，让用户能够快速浏览过去的记录，节省时间和精力。

8. 智能情绪支持：通过自然语言处理和情感识别技术，AI可以提供情绪支持和鼓励，帮助用户在困难时刻保持积极心态。

9. 智能分享：AI可以根据用户的喜好和兴趣，自动推荐适合分享的日记内容，并提供多种分享渠道，让用户能够与他人共享人生点滴。

10. 个性化写作指导：AI可以根据用户的写作习惯和偏好，提供个性化的写作指导和建议，帮助用户提升写作水平。

11. 智能翻译：AI可以提供实时翻译功能，让用户能够记录和阅读多种语言的日记。

12. 情感日历：AI可以根据用户的情感状态，生成情感日历，帮助用户更好地了解自己的情绪变化趋势。

13. 智能引导问题：AI可以根据用户的日记内容，提供相关的引导问题，帮助用户深入思考和回顾过去的经历。

14. 智能目标设定：AI可以帮助用户设定目标，并提供跟踪和反馈功能，帮助用户更好地实现自己的目标。

15. 智能时间线：AI可以根据用户的日记内容和时间信息，生成个人时间线，展示用户的成长和变化历程。

16. 智能情感分析：AI可以分析用户的情感倾向和偏好，提供个性化的情感支持和建议。

17. 智能回顾视频：AI可以根据用户的日记内容和照片，自动生成回顾视频，让用户能够以视觉形式回顾过去的时光。

18. 智能旅行规划：AI可以根据用户的日记内容和旅行偏好，提供个性化的旅行规划和推荐，帮助用户更好地体验旅行。

19. 智能情感日记：AI可以根据用户的情感状态和需求，提供个性化的情感日记模板和写作指导，帮助用户表达情感和记录心情。

20. 智能社交互动：AI可以提供用户之间的社交互动平台，让用户能够分享和交流日记内容，获得更多的反馈和支持。

## Kiko UI框架优化

Kiko UI是一个轻量级的UI框架，为iDiary提供界面交互功能。最近的优化包括：

### 1. 框架结构优化

- 移除了auth.js和app.js中的向后兼容代码
- 创建了00-kiko-init.js统一管理框架初始化
- 简化了pager.js实现，支持无target参数的API请求
- 重写了AJAX请求处理，直接使用jQuery 4.0原生方法
- 增强了错误日志记录和调试功能

### 2. 表单提交修复

- 修复了"Cannot read properties of undefined (reading 'id')"错误
- 优化了getDataFromSource实现，确保表单数据正确序列化
- 解决了POST请求表单数据未传递的问题（content-length: 0）
- 简化了表单处理逻辑，删除冗余代码

### 3. 脚本加载与缓存优化

- 修改构建脚本，合并kiko.min.js和idiary.min.js为idiary.bundle.js
- 实现开发模式使用时间戳，生产模式使用版本号的缓存控制
- 优化脚本加载顺序，确保jQuery在其他脚本之前加载
- 从package.json读取版本号，统一版本管理

### 4. 页面交互优化

- 将滚动监听功能移至初始化文件
- 支持登录后重定向和表单禁用/启用
- 统一错误处理机制
- 简化组件初始化流程

# 项目结构

```plaintext
.
├── cmd/                    # 主程序入口
│   └── idiary/            # iDiary 主程序
│       └── main.go        # 程序入口点
├── internal/              # 内部应用逻辑
│   ├── app/               # 应用层
│   │   ├── bootstrap/     # 应用启动
│   │   ├── controller/    # 控制器
│   │   ├── middleware/    # 中间件
│   │   ├── route/         # 路由配置
│   │   ├── templates/     # 模板文件
│   │   └── view/          # 视图层
│   ├── core/              # 核心系统
│   │   ├── bootstrap/     # 核心启动
│   │   ├── config/        # 配置管理
│   │   ├── database/      # 数据库连接
│   │   ├── di/            # 依赖注入
│   │   ├── error/         # 错误处理
│   │   ├── middleware/    # 核心中间件
│   │   └── validation/    # 数据验证
│   ├── domain/            # 业务领域
│   │   ├── diary/         # 日记模块
│   │   ├── post/          # 文章模块
│   │   ├── task/          # 任务模块
│   │   └── user/          # 用户模块
│   └── shared/            # 共享组件
│       ├── contracts/     # 接口契约
│       ├── http/          # HTTP工具
│       └── i18n/          # 国际化
├── pkg/                   # 可重用的包
│   ├── cookie.go          # Cookie 工具
│   ├── crypto.go          # 加密工具
│   ├── helpers.go         # 辅助函数
│   ├── logging.go         # 日志工具
│   └── ...               # 其他工具包
├── kiko/                  # Kiko UI 框架
│   ├── src/               # 源代码
│   │   ├── core.js        # 核心功能
│   │   ├── i18n/          # 国际化
│   │   ├── ui/            # UI组件
│   │   └── utils/         # 工具函数
│   └── build.js           # 构建脚本
├── website/               # 前端资源
│   ├── static/            # 静态文件
│   │   ├── img/           # 图片资源
│   │   └── src/           # 前端源码
│   └── views/             # 模板文件
├── scripts/               # 构建和部署脚本
│   ├── build/             # 构建脚本
│   └── deploy/            # 部署脚本
├── sql/                   # 数据库脚本
├── tests/                 # 测试代码
│   ├── COMPREHENSIVE_TEST_GUIDE.md # 测试指南
│   ├── coverage.out       # 测试覆盖率报告
│   ├── helpers/           # 测试辅助工具
│   ├── integration/       # 集成测试
│   │   ├── app/           # 应用层集成测试
│   │   ├── database/      # 数据库集成测试
│   │   └── services/      # 服务集成测试
│   └── unit/              # 单元测试
│       ├── app/           # 应用层单元测试
│       ├── core/          # 核心层单元测试
│       ├── domain/        # 业务领域单元测试
│       └── pkg/           # 工具包单元测试
├── logs/                  # 日志文件
└── __debug_bin/           # 调试构建输出
```

# 构建和开发

## 开发环境要求

- Go 1.21+
- Node.js 16+
- MySQL 8.0+
- Redis (可选)

## 构建脚本

### 前端构建

```bash
# 编译 CSS 文件
npm run build:css

# 编译 JavaScript 文件
npm run build:js

# 同时构建 CSS 和 JS
npm run build
```

### 后端构建

```bash
# 构建主程序
go build -o bin/idiary ./cmd/idiary/main.go

# 使用 VS Code 任务构建
# - build:css: 构建CSS文件
# - build:js: 构建JavaScript文件
# - build:go: 构建Go程序
```

## 开发指南

### 添加新的业务模块

1. 在 `internal/domain/` 下创建新模块目录
2. 实现 `service.go`, `handler.go`, `repository.go`, `types.go`
3. 在 `internal/app/route/` 中注册路由
4. 在 `tests/unit/domain/` 和 `tests/integration/` 中添加相应测试

### 运行测试

```bash
# 运行所有测试
go test ./tests/...

# 运行特定模块的单元测试
go test ./tests/unit/domain/post/

# 运行集成测试
go test ./tests/integration/...

# 生成测试覆盖率报告
go test -coverprofile=tests/coverage.out ./tests/...
```

### 添加国际化消息

1. 在 `website/locales/` 目录下的语言文件中添加消息
2. 在模板中使用国际化函数访问消息
3. 前端JavaScript使用Kiko框架的i18n功能

## 最新架构特性

### 领域驱动设计 (DDD)

- **领域分离**: 将业务逻辑从 `internal/module/` 迁移到 `internal/domain/`
- **依赖注入**: 使用DI容器管理组件依赖关系
- **接口抽象**: 通过接口实现松耦合设计
- **测试驱动**: 完善的单元测试和集成测试覆盖

### 验证系统优化

- **分离验证逻辑**: 针对不同操作使用不同的验证方法
- **用户ID验证**: `AddPost` 操作包含用户ID验证，`UpdatePost` 操作不包含
- **错误处理**: 统一的错误处理机制和错误码管理

## icon copyright

<a target="_blank" href="https://icons8.com/icon/jmkHXlIGZKw5/diary">Diary</a> icon by <a target="_blank" href="https://pictogrammers.github.io/@mdi/font/6.5.95/">Material Design Icons</a>

# 数据库结构更新总结

## 更新概述

本次更新使用Node.js脚本成功将cd数据库的表结构更新为最新的schema。

## 更新过程

1. 创建了包含最新schema的SQL文件 `update_schema.sql`
2. 开发了Node.js脚本 `update_db.js` 执行SQL文件
3. 创建了验证脚本 `verify_db_tables.js` 和 `db_update_summary.js` 检查数据库结构

## 更新结果

- 所有预期的12个表都已成功创建：

  - br_comment
  - br_diary
  - br_feed
  - br_follow
  - br_job_history
  - br_job_queue
  - br_like
  - br_post
  - br_diary_member
  - br_user
  - br_user_stat
  - br_valid_code

- 同时保留了一些数据库中已有的表（未在schema中定义）：
  - br_inbox_detail
  - br_inbox_header
  - br_inbox_stat
  - br_like_history
  - br_my_follow
  - br_online
  - br_space_member
  - br_sys_param

## 数据库连接信息

### 开发环境
- 主机: 127.0.0.1
- 端口: 3306
- 用户名: root
- 密码: 123123
- 数据库名: cd

### 测试环境
- 主机: 127.0.0.1
- 端口: 3306
- 用户名: root
- 密码: 123123
- 数据库名: hebe_test

### 🎯 中间件系统 (Middleware System)

**位置**: `internal/app/middleware/`

**功能**:

- 用户认证中间件
- 错误处理中间件
- 日志记录中间件
- 权限验证中间件

### 🎨 前端框架 (Kiko UI)

**位置**: `kiko/`

**特点**:

- 轻量级UI框架
- 基于jQuery 4.0
- 支持AJAX表单提交
- 响应式设计和组件化

## 技术栈

- **后端**: Go 1.21+, Gin Web Framework
- **数据库**: MySQL 8.0+
- **缓存**: Redis (可选)
- **前端**: Kiko UI, jQuery 4.0, HTML5/CSS3
- **模板**: Go HTML Template
- **构建**: npm, Go modules

## 特色功能

### 🤖 AI 增强功能 (规划中)

基于 slogan "记录人生，伴你前行" 的AI功能：

1. **情感分析**: 分析日记情感，提供情绪趋势
2. **智能标签**: 自动为日记添加分类标签
3. **智能回忆**: 个性化的历史回顾推送
4. **情感支持**: AI驱动的情绪支持和建议
5. **智能摘要**: 自动生成日记摘要和时间线

### 🌍 多语言支持

支持以下语言：

- 中文 (简体/繁体)
- English
- Deutsch
- Español
- Français
- 日本語
- 한국어
- Русский
- العربية
