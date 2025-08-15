# 构建系统文档

本文档描述了 iDiary.ai Hebe 项目的构建系统，特别是对新主题架构的支持。

## 概述

构建系统现在完全支持基于 Kiko UI 的主题架构，包括：
- Kiko UI 框架构建
- CSS 主题系统构建
- JavaScript 应用构建
- 主题验证系统

## 构建脚本

### 1. CSS 构建 (`scripts/build/build-css.js`)

**功能**：
- 使用 Tailwind CSS 4 生成基础样式
- 支持开发和生产模式
- 合并自定义主题 CSS
- 使用 PostCSS 处理和压缩

**命令**：
```bash
# 开发模式（保留调试信息）
npm run build:css:dev

# 生产模式（压缩优化）
npm run build:css:prod
```

**输出**：
- `website/static/dist/css/idiary.min.css`

### 2. JavaScript 构建 (`scripts/build/build-js.js`)

**功能**：
- 按指定顺序合并 JavaScript 文件
- 生成两个独立文件：
  - `kiko.min.js`: Kiko UI 框架
  - `idiary.min.js`: 业务逻辑
- 支持开发和生产模式
- 使用 Terser 压缩和混淆

**命令**：
```bash
# 开发模式
npm run build:js:dev

# 生产模式
npm run build:js:prod
```

**输出**：
- `website/static/dist/js/kiko.min.js`
- `website/static/dist/js/idiary.min.js`

### 3. Kiko UI 构建 (`kiko/build.js`)

**功能**：
- 合并所有 Kiko UI 组件
- 包含主题管理器
- 支持国际化
- 生成压缩版本

**命令**：
```bash
# 开发模式
npm run build:kiko

# 生产模式
npm run build:kiko:prod
```

**输出**：
- `kiko/dist/kiko.js`
- `kiko/dist/kiko.min.js`

### 4. 主题验证 (`scripts/build/verify-theme.js`)

**功能**：
- 验证 Kiko UI 主题支持
- 检查 CSS 主题类
- 验证 JavaScript 主题功能
- 确认构建完整性

**命令**：
```bash
npm run verify:theme
```

## 构建模式

### 开发模式
- 保留调试信息和注释
- 不压缩 CSS（便于调试）
- 保留 console.log 语句
- 使用时间戳避免缓存

### 生产模式
- 移除所有调试信息
- 压缩和混淆代码
- 移除 console.log 语句
- 优化文件大小

## 完整构建流程

### 开发构建
```bash
npm run build:dev
```

执行顺序：
1. `npm run build:kiko` - 构建 Kiko UI 框架
2. `npm run build:css:dev` - 构建 CSS（开发模式）
3. `npm run build:js:dev` - 构建 JavaScript（开发模式）

### 生产构建
```bash
npm run build:prod
```

执行顺序：
1. `npm run build:kiko:prod` - 构建 Kiko UI 框架（生产模式）
2. `npm run build:css:prod` - 构建 CSS（生产模式）
3. `npm run build:js:prod` - 构建 JavaScript（生产模式）

### 构建测试
```bash
npm run test:build
```

执行顺序：
1. `npm run build:prod` - 生产构建
2. `npm run verify:theme` - 主题验证

## 主题架构支持

### Kiko UI 主题管理
- 集成主题管理器 (`theme.js`)
- 支持 light/dark 主题切换
- 提供主题感知组件
- 统一的主题 API

### CSS 主题系统
- 基于 Tailwind CSS 4
- 使用 CSS 变量实现主题
- 预定义主题类（`.card`, `.btn-*` 等）
- 响应式主题支持

### JavaScript 主题集成
- 主题状态管理
- 自动主题检测
- 主题切换动画
- 本地存储持久化

## 文件结构

```
hebe/
├── scripts/build/
│   ├── build-css.js          # CSS 构建脚本
│   ├── build-js.js           # JavaScript 构建脚本
│   ├── verify-theme.js       # 主题验证脚本
│   └── create-dirs.js        # 目录创建脚本
├── kiko/
│   ├── src/                  # Kiko UI 源码
│   ├── dist/                 # Kiko UI 构建输出
│   └── build.js              # Kiko UI 构建脚本
├── website/static/
│   ├── src/                  # 源文件
│   │   ├── css/
│   │   └── js/
│   └── dist/                 # 构建输出
│       ├── css/
│       └── js/
├── postcss.config.js         # PostCSS 配置
├── tailwind.config.js        # Tailwind 配置
└── package.json              # 构建脚本定义
```

## 配置文件

### PostCSS 配置 (`postcss.config.js`)
- 支持环境变量配置
- 开发模式：不压缩
- 生产模式：完全压缩

### Tailwind 配置 (`tailwind.config.js`)
- 主题系统配置
- 自定义颜色变量
- 响应式断点
- 组件类定义

## 性能优化

### 文件大小优化
- Kiko UI: ~51KB (压缩后)
- CSS: ~58KB (包含 Tailwind + 自定义样式)
- JavaScript: ~18KB (业务逻辑)

### 构建速度优化
- 并行构建支持
- 增量构建（开发模式）
- 缓存优化

## 故障排除

### 常见问题

1. **构建失败**
   ```bash
   npm run verify:theme
   ```
   检查构建输出和主题支持

2. **主题不工作**
   - 检查 Kiko UI 是否正确加载
   - 验证 CSS 变量是否定义
   - 确认 JavaScript 主题管理器已初始化

3. **文件缺失**
   ```bash
   npm run create:dirs
   npm run build:dev
   ```

### 调试模式

开发模式下，所有文件都保留调试信息：
```bash
npm run build:dev
```

## 更新日志

### v2.0.0 - 主题架构支持
- ✅ 集成 Kiko UI 主题系统
- ✅ 支持开发/生产构建模式
- ✅ 添加主题验证脚本
- ✅ 优化构建性能
- ✅ 完善错误处理

### v1.0.0 - 基础构建系统
- ✅ CSS/JavaScript 构建
- ✅ Tailwind CSS 集成
- ✅ 文件压缩和优化