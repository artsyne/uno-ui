# JavaScript 文件结构

此目录采用模块化结构组织JavaScript文件，分为三个主要类别:

## 目录结构

```
website/static/src/js/
  ├── tools/     # 通用工具脚本和UI组件
  ├── app/       # 应用相关的业务逻辑
  └── vendors/   # 第三方库和依赖
```

## 目录说明

### 1. tools/

包含与特定业务逻辑无关的通用工具和UI组件：

- `box.js` - 对话框和弹窗组件
- `tab.js` - 标签页组件
- `valid.js` - 表单验证工具
- `grid.js` - 网格组件
- `util.js` - 通用工具函数
- `ajaxLoader.js` - AJAX加载器

### 2. app/

包含应用特定的业务逻辑：

- `app.js` - 应用初始化和全局功能
- `diary.js` - 日记相关功能
- `passport.js` - 用户认证相关功能
- `uploadimg.js` - 图片上传功能

### 3. vendors/

包含第三方库和依赖：

- `imagesloaded.pkgd.min.js` - 图片加载库

## 构建流程

JavaScript文件按照以下顺序构建和合并：

1. vendors/ (第三方库最先加载)
2. tools/ (工具脚本其次)
3. app/ (应用脚本最后)

构建命令：

```bash
npm run build:js
```

构建脚本位于 `scripts/build/build-js.js`，它会按上述顺序处理所有JS文件，并使用Terser进行压缩和混淆。

生成的文件保存在 `website/static/dist/js/idiary.min.js`。

## 跨平台兼容性

构建脚本已优化为同时兼容Mac和Windows系统：

- 使用Node.js的path模块处理路径分隔符
- 提供跨平台的目录创建脚本
- 使用shell:true确保命令在不同平台正确执行
- 对路径中可能包含空格的文件名添加引号

创建目录结构：

```bash
npm run create:dirs
```

## 开发指南

1. 添加新的通用组件或工具应放在 `tools/` 目录
2. 添加新的业务逻辑应放在 `app/` 目录
3. 添加新的第三方依赖应放在 `vendors/` 目录
4. 构建后请检查控制台确保没有错误
