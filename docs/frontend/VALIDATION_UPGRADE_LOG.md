# Uno Valid v0.8.0 升级记录

**日期**: 2026-01-22
**执行人**: Trae AI

## 1. 概述

本项目已完成 Uno Valid 插件从旧版本到 v0.8.0 的升级适配。此次升级主要涉及验证规则的配置方式变更（从 HTML5 属性迁移到 `x-uno-valid` 数组配置）以及自定义验证器的异步返回机制优化。

## 2. 变更详情

### 2.1 模板文件 (Templates)

所有使用 `x-uno-valid` 的 HTML 模板已更新为数组配置模式，移除了 `required`, `pattern`, `minlength`, `maxlength` 等原生属性。

| 文件路径 | 修改内容 |
| :--- | :--- |
| `website/views/passport/register.html` | 用户名、邮箱、密码字段规则迁移，添加 `debounce` |
| `website/views/passport/login.html` | 用户名、密码字段规则迁移 |
| `website/views/passport/findpwd.html` | 邮箱字段规则迁移 |
| `website/views/passport/resetpwd.html` | 密码、确认密码字段规则迁移，`matchPwd` 规则显式指定 `target` |
| `website/views/user/settings.html` | 昵称、旧密码、新密码字段规则迁移 |
| `website/views/form/entry.html` | 标题字段规则迁移 |
| `website/views/form/diary.html` | 标题、备注字段规则迁移 |
| `website/views/entry/detail.html` | 评论框字段规则迁移 |

### 2.2 JavaScript 逻辑 (JS)

**`website/static/src/js/modules/auth.js`**:

1. **`checkUsername` & `checkEmail`**:
   - 移除了手动调用 `validStore.showError` 的逻辑。
   - 改为直接返回错误消息字符串（验证失败）或 `true`（验证通过）。
   - 添加了 `try-catch` 处理，确保网络异常时返回通用错误提示。
   - 利用 Uno v0.8.0 的原生异步支持，简化了代码。

2. **`matchPwd`**:
   - 添加了对 `rule.target` 的支持，允许通过配置指定要比对的目标输入框选择器。
   - 保持了对 `#new_pwd` 的默认回退支持，确保向后兼容。

### 2.3 文档 (Docs)

- 更新了 `docs/frontend/COMPONENT_API.md`，新增了 "4. 表单验证" 章节，说明新的配置方式和自定义验证器 API。

## 3. 测试验证

- **单元测试**: 创建了 `tests/js/auth_validation_test.mjs`，对 `auth.js` 中的自定义验证器进行了模拟测试，覆盖了异步验证成功/失败和同步密码匹配场景。
- **集成验证**: 建议运行项目并手动测试注册、登录、修改密码等流程，确保 UI 交互正常。

## 4. 注意事项

- 前端构建: 修改了 `auth.js` 后，需要运行 `npm run build` (或 `npm run build:frontend`) 重新生成 `bundle.js`。
- 缓存: 部署后请清理浏览器缓存以加载最新的 JS 文件。
