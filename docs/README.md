# iDiary.AI 项目文档索引

本目录包含iDiary.AI项目的所有技术文档。

## 📚 文档分类

### 架构设计文档

- [internal-architecture-design.md](./internal-architecture-design.md) - 内部架构设计
- [storage-architecture-design.md](./storage-architecture-design.md) - 存储架构设计

### 重构文档

- [REFACTOR_SUMMARY.md](./REFACTOR_SUMMARY.md) - 总体重构总结
- [refactor-render-summary.md](./refactor-render-summary.md) - 渲染系统重构总结
- [pagedata-update-summary.md](./pagedata-update-summary.md) - PageData结构更新总结
- [middleware-refactor-summary.md](./middleware-refactor-summary.md) - 中间件重构总结
- [main-refactor-summary.md](./main-refactor-summary.md) - Main.go重构总结

### 使用指南

- [template-usage-guide.md](./template-usage-guide.md) - 模板使用指南
- [template-pagedata-usage.md](./template-pagedata-usage.md) - PageData模板访问更新指南
- [error-handling-examples.md](./error-handling-examples.md) - 错误处理示例

## 🔧 模块文档

以下模块包含独立的README文档：

- [internal/core/di/README.md](../internal/core/di/README.md) - 依赖注入容器文档
- [internal/core/config/README.md](../internal/core/config/README.md) - 配置系统文档
- [kiko/README.md](../kiko/README.md) - Kiko UI框架文档
- [tests/README.md](../tests/README.md) - 测试文档

## 📝 开发指南

### 添加新文档

1. **架构文档**: 放在 `docs/` 目录下
2. **模块文档**: 放在对应模块的 `README.md` 中
3. **API文档**: 使用代码注释生成
4. **用户文档**: 放在项目根目录的 `README.md` 中

### 文档命名规范

- 使用小写字母和连字符
- 架构文档: `xxx-architecture-design.md`
- 重构文档: `xxx-refactor-summary.md`
- 使用指南: `xxx-usage-guide.md`
- 示例文档: `xxx-examples.md`

### 文档维护

- 每次重大更新后同步更新相关文档
- 保持文档结构清晰，使用统一的格式
- 及时清理过时的文档内容

## 🚀 快速导航

- **新手入门**: 查看 [README.md](../README.md)
- **架构了解**: 查看 [internal-architecture-design.md](./internal-architecture-design.md)
- **开发指南**: 查看各模块的README文档
- **问题排查**: 查看 [error-handling-examples.md](./error-handling-examples.md)
