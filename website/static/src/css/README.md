# CSS 文件结构

此目录包含网站的CSS样式文件，采用Tailwind CSS框架。

## 文件说明

- `tailwind.css` - Tailwind CSS基础文件，包含@tailwind指令
- `diary-theme.css` - 自定义主题样式，会被合并到最终CSS中

## 构建流程

CSS文件按照以下步骤构建：

1. 使用Tailwind CLI处理tailwind.css，生成完整的Tailwind样式
2. 合并diary-theme.css和生成的Tailwind样式
3. 使用PostCSS处理合并后的文件（添加前缀、压缩等）

构建命令：

```bash
npm run build:css
```

构建脚本位于 `scripts/build/build-css.js`，最终生成的文件保存在 `website/static/dist/css/idiary.min.css`。

## 跨平台兼容性

构建脚本已优化为同时兼容Mac和Windows系统：

- 使用Node.js的path模块处理路径
- 使用文件读写替代cat/type命令合并文件
- 提供完善的错误处理

## 开发指南

1. 添加新的样式应放在`diary-theme.css`中
2. 修改Tailwind配置请编辑项目根目录下的`tailwind.config.js`
3. 通过`class`属性应用Tailwind类名
4. 使用PostCSS插件可在`postcss.config.js`中配置
