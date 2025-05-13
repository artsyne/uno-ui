# 前端资源构建系统

这个目录包含了网站的所有前端资源，包括JavaScript、CSS和静态资源文件。

## 目录结构

```
website/static/
  ├── src/           # 源文件
  │   ├── js/        # JavaScript源文件
  │   │   ├── app/     # 应用业务逻辑
  │   │   ├── tools/   # 工具组件
  │   │   └── vendors/ # 第三方库
  │   ├── css/       # CSS源文件
  │   └── img/       # 图片资源
  └── dist/         # 构建输出
      ├── js/        # 编译后的JS
      ├── css/       # 编译后的CSS
      └── img/       # 优化后的图片
```

## 构建流程

整个前端资源构建系统支持以下命令：

```bash
# 构建JavaScript
npm run build:js

# 构建CSS
npm run build:css

# 构建所有前端资源
npm run build

# 创建必要的目录结构
npm run create:dirs
```

## 跨平台兼容性

所有构建脚本都已优化为同时兼容Mac、Linux和Windows系统，主要特点包括：

1. 使用Node.js的path模块处理跨平台路径问题
2. 避免使用平台特定的Shell命令
3. 使用文件读写API替代cat/type等命令
4. 所有文件路径使用引号，处理路径中的空格
5. 提供详细的错误处理和日志输出

## 技术栈

- **CSS框架**: Tailwind CSS 4
- **CSS处理**: PostCSS、Autoprefixer、CSSnano
- **JS处理**: Terser (用于压缩和混淆)
- **构建工具**: Node.js脚本

## 开发指南

1. 所有源文件修改应在`src`目录进行
2. 不要直接修改`dist`目录，它会在构建时被覆盖
3. 运行完整构建命令后检查控制台日志
4. 查看各子目录下的README.md了解详细信息

## 详细文档

- [JavaScript构建文档](src/js/README.md)
- [CSS构建文档](src/css/README.md)
