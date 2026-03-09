# iDiary 扩展 App 架构设计方案

> 插件化架构，支持 AI 分析、导出备份等扩展功能

---

## 一、计费策略与套餐设计

### 1.1 套餐计划

系统提供三个订阅计划，满足不同用户需求：

| 计划 | 价格 | 目标用户 | 核心权益 |
|------|------|----------|----------|
| **Free** | 免费 | 轻度用户 | 基础日记功能 + 有限 AI 体验 |
| **Pro** | ¥30/月 | 重度用户 | 无限 AI 次数 + 高级导出 |
| **Max** | ¥99/月 | 高级用户 | 全功能 + 优先队列 + 云备份 |

### 1.2 套餐功能对比

| 功能 | Free | Pro | Max |
|------|------|-----|-----|
| 日记本数量 | 3 个 | 无限 | 无限 |
| AI 摘要/日 | 5 次 | 无限 | 无限 |
| AI 情感分析/日 | 5 次 | 无限 | 无限 |
| AI 写作润色/日 | 2 次 | 50 次 | 无限 |
| PDF 导出/月 | 2 次 | 20 次 | 无限 |
| 云端备份 | ❌ | ❌ | ✓ |
| 优先队列 | ❌ | ❌ | ✓ |
| API 调用 | ❌ | ❌ | ✓ |
| 客服支持 | ❌ | 邮件 | 优先 |

### 1.3 套餐表结构 (br_plan)

```sql
CREATE TABLE `br_plan` (
  `plan_id` tinyint unsigned NOT NULL AUTO_INCREMENT,
  `plan_name` varchar(50) NOT NULL COMMENT 'Free, Pro, Max',
  `price_monthly` int unsigned NOT NULL DEFAULT '0' COMMENT '月付价格(分)',
  `price_yearly` int unsigned NOT NULL DEFAULT '0' COMMENT '年付价格(分)',
  `features` text COMMENT '功能配置(JSON)',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态(0:禁用 1:启用)',
  PRIMARY KEY (`plan_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**Free 计划 features 配置**：
```json
{
  "diary_limit": 3,
  "ai": {
    "ai_summary": { "daily_limit": 5, "enabled": true },
    "ai_sentiment": { "daily_limit": 5, "enabled": true },
    "ai_polish": { "daily_limit": 2, "enabled": true }
  },
  "export": {
    "export_pdf": { "monthly_limit": 2, "enabled": true },
    "export_markdown": { "monthly_limit": 5, "enabled": true }
  },
  "backup": { "enabled": false },
  "priority_queue": false
}
```

**Pro 计划 features 配置**：
```json
{
  "diary_limit": -1,
  "ai": {
    "ai_summary": { "daily_limit": -1, "enabled": true },
    "ai_sentiment": { "daily_limit": -1, "enabled": true },
    "ai_polish": { "daily_limit": 50, "enabled": true },
    "ai_continue": { "daily_limit": 20, "enabled": true }
  },
  "export": {
    "export_pdf": { "monthly_limit": 20, "enabled": true },
    "export_markdown": { "monthly_limit": -1, "enabled": true },
    "export_epub": { "monthly_limit": 5, "enabled": true }
  },
  "backup": { "enabled": false },
  "priority_queue": false
}
```

**Max 计划 features 配置**：
```json
{
  "diary_limit": -1,
  "ai": {
    "ai_summary": { "daily_limit": -1, "enabled": true },
    "ai_sentiment": { "daily_limit": -1, "enabled": true },
    "ai_polish": { "daily_limit": -1, "enabled": true },
    "ai_continue": { "daily_limit": -1, "enabled": true },
    "ai_translate": { "daily_limit": -1, "enabled": true }
  },
  "export": {
    "export_pdf": { "monthly_limit": -1, "enabled": true },
    "export_markdown": { "monthly_limit": -1, "enabled": true },
    "export_epub": { "monthly_limit": -1, "enabled": true }
  },
  "backup": { "enabled": true },
  "priority_queue": true
}
```

### 1.4 用户订阅表 (br_user_subscription)

```sql
CREATE TABLE `br_user_subscription` (
  `user_id` bigint unsigned NOT NULL COMMENT '关联用户ID',
  `plan_id` tinyint unsigned NOT NULL COMMENT '关联套餐计划ID',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '订阅状态(0:取消 1:激活)',
  `start_utc_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '订阅开始时间',
  `expire_utc_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '订阅到期时间',
  `last_payment_utc_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '最后支付时间',
  PRIMARY KEY (`user_id`),
  KEY `idx_expire_status` (`expire_utc_time`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 1.5 使用量统计表 (br_app_usage)

```sql
CREATE TABLE `br_app_usage` (
  `usage_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `app_key` varchar(50) NOT NULL,
  `usage_date` date NOT NULL COMMENT '使用日期',
  `usage_count` int unsigned NOT NULL DEFAULT '0' COMMENT '使用次数',
  `token_usage` bigint unsigned NOT NULL DEFAULT '0' COMMENT '消耗Token数',
  PRIMARY KEY (`usage_id`),
  UNIQUE KEY `uk_user_app_date` (`user_id`, `app_key`, `usage_date`),
  KEY `idx_month` (`user_id`, `app_key`, `usage_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

---

## 二、架构核心思路

- **微内核设计**：日记本核心逻辑保持纯粹，所有增强功能（AI、备份、导出）均以 App 形式注册
- **配置解耦**：AI 的提示词（Prompt）与应用注册分离，支持热更新。**无需修改代码**即可新增或调整 AI 应用
- **双层路由**：应用分为日记本维度（Diary Scope）和用户维度（User Scope），满足不同隐私需求
- **权限分级**：通过套餐（Plan）控制用户可使用的 App 集合，实现商业化变现

---

## 三、数据库设计

### 3.1 应用注册表 (br_app)

```sql
CREATE TABLE `br_app` (
  `app_id` int unsigned NOT NULL AUTO_INCREMENT,
  `app_key` varchar(50) NOT NULL COMMENT '唯一标识: ai_sentiment, export_pdf',
  `app_name` varchar(100) NOT NULL COMMENT '显示名称',
  `app_scope` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '1:日记本维度, 2:用户维度',
  `app_category` tinyint unsigned NOT NULL COMMENT '1:AI, 2:导出, 3:备份, 4:数据处理',
  `config_json` text COMMENT '通用配置(JSON)',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '1:启用, 2:禁用',
  `add_utc_time` bigint unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`app_id`),
  UNIQUE KEY `uk_app_key` (`app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 3.2 AI 专用配置表 (br_app_ai)

每个 AI App 绑定特定的 AI 模型和提示词配置：

```sql
CREATE TABLE `br_app_ai` (
  `app_key` varchar(50) NOT NULL COMMENT '关联 br_app.app_key',
  `provider` varchar(20) NOT NULL DEFAULT 'minimax' COMMENT 'AI提供商: minimax, openai, glm',
  `model_name` varchar(100) NOT NULL DEFAULT 'abab6.5s-chat' COMMENT '模型名称',
  `system_prompt` text COMMENT 'AI 角色设定',
  `user_prompt_template` text COMMENT '指令模板, 使用 {{.Content}} 占位',
  `response_format` text COMMENT 'JSON 格式要求',
  `temperature` float NOT NULL DEFAULT '0.7',
  `max_tokens` int unsigned NOT NULL DEFAULT '2000',
  PRIMARY KEY (`app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 3.3 应用开启状态表

日记本维度（br_diary_app）：
```sql
CREATE TABLE `br_diary_app` (
  `unique_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `diary_id` bigint unsigned NOT NULL,
  `app_key` varchar(50) NOT NULL,
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '1:开启, 2:关闭',
  `add_utc_time` bigint unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`unique_id`),
  UNIQUE KEY `uk_diary_app` (`diary_id`, `app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

用户维度（br_user_app）：
```sql
CREATE TABLE `br_user_app` (
  `unique_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `app_key` varchar(50) NOT NULL,
  `status` tinyint unsigned NOT NULL DEFAULT '1',
  `add_utc_time` bigint unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`unique_id`),
  UNIQUE KEY `uk_user_app` (`user_id`, `app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 3.4 扩展应用数据存储表 (br_app_data)

通用扩展数据存储，存储 AI 结果、导出配置等：

```sql
CREATE TABLE `br_app_data` (
  `data_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `app_key` varchar(50) NOT NULL COMMENT '关联应用',
  `target_id` bigint unsigned NOT NULL COMMENT '关联目标ID(如entry_id, diary_id)',
  `target_type` tinyint unsigned NOT NULL COMMENT '目标类型(1:日记 2:日记本 3:用户)',
  `data_json` json NOT NULL COMMENT '扩展存储的数据',
  `add_utc_time` bigint unsigned NOT NULL,
  PRIMARY KEY (`data_id`),
  UNIQUE KEY `uk_app_target` (`app_key`, `target_id`, `target_type`)
) ENGINE=InnoDB COMMENT='扩展应用数据存储表';
```

---

## 四、AI 模型与扩展设计

### 4.1 AI 模型是独立扩展

**核心设计**：AI 模型作为独立扩展资源，与 AI App 绑定。每个 AI App 可以配置不同的模型，实现按需选择。

```
br_app (AI App) ←→ br_app_ai (配置) ←→ AI Provider (模型)
```

- 不同的 AI App 可以使用不同的模型（如情感分析用轻量模型，深度分析用增强模型）
- 模型可以独立升级或切换，不影响 App 逻辑
- 支持多模型提供商：MiniMax、OpenAI、GLM、Claude

### 4.2 模型选择指南

| 模型 | 提供商 | 适用场景 | 特点 |
|------|--------|----------|------|
| abab6.5s-chat | MiniMax | 情感分析、写作润色 | 主力模型，性价比高 |
| abab6.5g-chat | MiniMax | 深度分析、多日记综合 | 推理能力强 |
| abab4-chat | MiniMax | 简单分类、标签生成 | 轻量快速 |
| gpt-4o | OpenAI | 通用对话、高质量输出 | 国际首选 |
| glm-4 | GLM | 中文场景、国产替代 | 国内合规 |

### 4.3 预设 AI App 清单

| app_key | app_name | 模型 | 功能 |
|---------|-----------|------|------|
| `ai_summary` | AI 摘要 | abab6.5s-chat | 自动生成日记摘要 |
| `ai_sentiment` | 情感分析 | abab6.5s-chat | 分析情感倾向和关键词 |
| `ai_polish` | AI 润色 | abab6.5s-chat | 优化文字表达 |
| `ai_continue` | 续写建议 | abab6.5g-chat | 提供写作续写灵感 |
| `ai_title` | 标题生成 | abab4-chat | 自动生成吸引人的标题 |
| `ai_translate` | 智能翻译 | gpt-4o | 多语言翻译 |
| `ai_tags` | 智能标签 | abab4-chat | 自动提取内容标签 |

---

## 五、预设扩展 App 清单

### 5.1 AI 类应用

| app_key | app_name | app_scope | 功能说明 |
|---------|-----------|-----------|----------|
| `ai_summary` | AI 摘要 | Diary | 自动生成日记摘要 |
| `ai_sentiment` | 情感分析 | Diary | 分析日记情感倾向 |
| `ai_polish` | AI 润色 | Diary | 优化文字表达 |
| `ai_continue` | 续写建议 | Diary | 提供续写灵感 |
| `ai_title` | 标题生成 | Diary | 自动生成标题 |
| `ai_translate` | 智能翻译 | Diary | 多语言翻译 |
| `ai_tags` | 智能标签 | Diary | 自动提取标签 |

### 5.2 导出类应用

| app_key | app_name | app_scope | 功能说明 |
|---------|-----------|-----------|----------|
| `export_pdf` | 导出 PDF | Diary | 导出为 PDF 电子书 |
| `export_markdown` | 导出 Markdown | Diary | 导出为 Markdown 文件 |
| `export_epub` | 导出 EPUB | Diary | 导出为电子书 |

### 5.3 备份类应用

| app_key | app_name | app_scope | 功能说明 |
|---------|-----------|-----------|----------|
| `backup_cloud` | 云端备份 | User | 定时备份到云存储 |
| `backup_email` | 邮件备份 | User | 定期发送备份到邮箱 |
| `sync_dropbox` | Dropbox 同步 | User | 实时同步到 Dropbox |

### 5.4 数据处理类应用

| app_key | app_name | app_scope | 功能说明 |
|---------|-----------|-----------|----------|
| `import_word` | Word 导入 | User | 批量导入 Word 文档 |
| `tag_extract` | 标签提取 | Diary | 自动提取日记标签 |

---

## 六、核心业务流程

### 6.1 权限校验流

当用户使用某个 App 时，后端执行以下校验：

```
1. 查询 br_user_subscription 获取 plan_id
2. 查询 br_plan.features 检查 App 是否在套餐中
3. 检查 App 的使用限额 (br_app_usage)
4. 校验作用域状态:
   - app_scope=1: 检查 br_diary_app 是否开启
   - app_scope=2: 检查 br_user_app 是否开启
5. 执行 App
```

### 6.2 AI 任务处理流

```
1. 触发: 日记保存成功，且 br_diary_app 开启了 AI App
2. 入队: 插入任务到 br_job_queue
3. 消费: Worker 读取 br_app_ai 配置 → 渲染提示词 → 调用 AI API
4. 存储: 解析结果，存入 br_app_data
5. 通知: 通过 WebSocket/轮询通知前端
```

### 6.3 导出/备份任务处理流

```
1. 触发: 用户点击导出/备份按钮
2. 入队: 写入 br_job_queue，包含格式、范围等参数
3. 消费: Worker 生成文件 → 上传到 OSS/S3
4. 通知: 生成下载链接，站内通知/邮件告知
```

---

## 七、扩展性说明

### 7.1 增加新 AI App

在 `br_app` 注册应用，在 `br_app_ai` 配置 Prompt。**无需修改代码**。

```sql
-- 注册 App
INSERT INTO br_app (app_key, app_name, app_scope, app_category) 
VALUES ('ai_weekly_report', 'AI 周报生成', 1, 1);

-- 配置 AI (绑定模型和提示词)
INSERT INTO br_app_ai (app_key, provider, model_name, system_prompt, user_prompt_template) 
VALUES ('ai_weekly_report', 'minimax', 'abab6.5s-chat', 
  '你是一个周报助手，帮助用户总结一周的日记内容。',
  '请根据以下日记内容生成周报总结：\n\n{{.Content}}');
```

### 7.2 避险设计

- **流量控制**：通过 `max_tokens` 和套餐限额限制 AI 成本
- **隔离性**：即使 AI 服务宕机，不影响日记本基础读写
- **数据安全**：备份 App 独立运行，支持用户配置非美国节点的存储
- **内容审核**：AI 输入输出需经过安全过滤

---

## 八、总结

本方案通过 **「功能注册 + 维度路由 + 权限映射」** 的模式，实现：

| 设计要点 | 实现方式 |
|----------|----------|
| 套餐计费 | Free/Pro/Max 三档，满足不同用户 |
| 扩展注册 | 所有扩展通过 br_app 表管理 |
| AI 模型 | 作为独立扩展，与 App 绑定，支持多提供商 |
| 双层作用域 | 日记本级和用户级，满足不同场景 |
| 异步任务 | 通过 br_job_queue 处理耗时操作 |
| 商业变现 | 套餐权限控制 + 使用量统计 |

**下一步**：生成 SQL 初始脚本，或编写 Go 语言核心函数代码。
