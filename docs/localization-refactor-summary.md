# 国际化模块重构总结

## 重构目标

- 面向未来的设计，去掉 MessageService 的兼容接口
- 采用预加载机制，提高运行时性能
- 配置驱动，消除硬编码
- 使用更准确的命名约定

## 主要变更

### 1. 文件结构变更

- **删除**: `internal/app/render/i18n.go` (旧的国际化实现)
- **新增**: `internal/app/render/localization.go` (新的本地化管理器)
- **新增**: `website/locales/en.json` (英语资源文件)
- **新增**: `website/locales/zh-cn.json` (简体中文资源文件)

### 2. 核心类重构

- `I18nManager` → `LocalizationManager`
- `NewI18nManager` → `NewLocalizationManager`
- 变量命名: `lang` → `locale` (更准确的术语)

### 3. 接口简化

- **删除**: `MessageService` 接口及其所有实现
- **简化**: `LocalizationManager` 只保留核心功能
- **预加载**: 所有语言文件在应用启动时加载

### 4. 配置驱动改进

- **新增配置字段**:
  - `LocalesDir`: 语言文件目录路径
  - `DefaultLang`: 默认语言
  - `TitleDescriptions`: 页面标题描述映射
  - `DefaultTitleDesc`: 默认页面描述

### 5. 方法签名更新

- `GetMessage(key, locale string) string` - 清晰的参数类型
- `GetMessageWithParams(key, locale string, params map[string]string) string`
- `GetUserLocale(user RenderUser) string`

## 核心优势

### 1. 性能提升

- **预加载机制**: 所有语言文件在启动时一次性加载，消除运行时 I/O
- **内存缓存**: 使用 `sync.RWMutex` 保证并发安全的快速访问
- **快速失败**: 启动时发现配置错误，而不是运行时

### 2. 代码清洁

- **消除硬编码**: 页面描述从硬编码 map 移至配置驱动
- **单一职责**: `LocalizationManager` 专注于本地化，不再承担额外的接口适配责任
- **类型安全**: 明确的参数类型，减少运行时错误

### 3. 可维护性

- **标准命名**: 使用行业标准的 `locale` 术语
- **配置集中**: 所有本地化相关配置都在 `.env.local` 中管理
- **易于扩展**: 添加新语言只需添加对应的 JSON 文件

## 文件内容示例

### 语言资源文件格式

```json
{
  "common.home": "首页",
  "page.title.diary_list": "我的日记",
  "diary.content_placeholder": "{{username}}，今天有什么新鲜事？"
}
```

### 使用方式

```go
// 获取本地化消息
message := localizationManager.GetMessage("common.home", "zh-cn")

// 带参数的消息
params := map[string]string{"username": "张三"}
message := localizationManager.GetMessageWithParams("diary.content_placeholder", "zh-cn", params)
```

## 兼容性说明

- **不向后兼容**: 完全删除了 MessageService 接口
- **依赖注入更新**: PageBuilder 构造函数参数变更
- **配置要求**: 需要在 `.env.local` 中添加新的配置项

## 测试建议

1. 验证所有语言文件能正确加载
2. 测试参数替换功能
3. 验证回退到默认语言的机制
4. 测试并发访问的安全性

## 后续计划

- [ ] 添加更多语言支持
- [ ] 实现热重载功能（可选）
- [ ] 添加语言文件校验工具
- [ ] 优化参数替换性能
