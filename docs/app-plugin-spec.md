
# idiary.ai 插件化架构与套餐权限技术方案

## 1. 架构核心思路

* **微内核设计**：日记本核心逻辑保持纯粹，所有增强功能（AI、备份、导出）均以 App 形式注册。
* **配置解耦**：AI 的提示词（Prompt）与应用注册分离，支持热更新。
* **双层路由**：应用分为 **用户维度 (User Scope)** 和 **日记本维度 (Diary Scope)**，满足不同隐私需求。
* **权限分级**：通过套餐（Plan）控制用户可使用的 App 集合。

---

## 2. 数据库设计 (Schema)

### 2.1 应用注册与 AI 配置

#### `br_app` (应用注册表)

所有功能（包括 AI 和非 AI）的元数据。

```sql
CREATE TABLE `br_app` (
  `app_id` int unsigned NOT NULL AUTO_INCREMENT,
  `app_key` varchar(50) NOT NULL COMMENT '唯一标识: sentiment, export_pdf, backup_s3',
  `app_name` varchar(100) NOT NULL,
  `app_scope` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '1:日记本维度, 2:用户维度',
  `app_category` tinyint unsigned NOT NULL COMMENT '1:AI, 2:导出, 3:备份, 4:导入',
  `config_json` text COMMENT '非AI类的通用配置, 如备份路径模板',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '1:启用, 2:禁用',
  `add_utc_time` bigint unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`app_id`),
  UNIQUE KEY `uk_app_key` (`app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

```

#### `br_app_ai` (AI 专用配置表)

仅针对 `app_category = 1` 的应用。

```sql
CREATE TABLE `br_app_ai` (
  `app_key` varchar(50) NOT NULL COMMENT '关联 br_app.app_key',
  `model_name` varchar(100) NOT NULL DEFAULT 'gpt-4o',
  `system_prompt` text COMMENT 'AI 角色设定',
  `user_prompt_template` text COMMENT '指令模板, 使用 {{.Content}} 占位',
  `response_format` text COMMENT 'JSON 格式要求',
  `temperature` float NOT NULL DEFAULT '0.7',
  `max_tokens` int unsigned NOT NULL DEFAULT '2000',
  PRIMARY KEY (`app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

```

### 2.2 套餐与权限控制

#### `br_plan` (套餐定义表)

```sql
CREATE TABLE `br_plan` (
  `plan_id` int unsigned NOT NULL AUTO_INCREMENT,
  `plan_name` varchar(50) NOT NULL COMMENT 'Free, Pro, Lifetime',
  `app_limit_json` text COMMENT '该套餐包含的 app_key 列表',
  `price_monthly` int unsigned DEFAULT '0',
  `status` tinyint unsigned NOT NULL DEFAULT '1',
  PRIMARY KEY (`plan_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

```

#### `br_user_subscription` (用户订阅表)

```sql
CREATE TABLE `br_user_subscription` (
  `unique_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `plan_id` int unsigned NOT NULL,
  `start_utc_time` bigint unsigned NOT NULL,
  `expire_utc_time` bigint unsigned NOT NULL,
  PRIMARY KEY (`unique_id`),
  KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

```

### 2.3 应用开启状态 (用户决策层)

#### `br_diary_app` (日记本维度开启状态)

```sql
CREATE TABLE `br_diary_app` (
  `unique_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `diary_id` bigint unsigned NOT NULL,
  `app_key` varchar(50) NOT NULL,
  `status` tinyint unsigned NOT NULL DEFAULT '1',
  `add_utc_time` bigint unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`unique_id`),
  UNIQUE KEY `uk_diary_app` (`diary_id`, `app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

```

---

## 3. 核心业务流程

### 3.1 权限校验流 (Middleware/Service)

当用户尝试使用或配置某个 App 时，Go 后端执行以下校验：

1. **套餐校验**：查询 `br_user_subscription` 获取当前 `plan_id`，检查 `br_plan.app_limit_json` 是否包含该 `app_key`。
2. **作用域校验**：
* 如果是 `app_scope = 1` (Diary)，校验 `br_diary_app` 是否为开启状态。
* 如果是 `app_scope = 2` (User)，校验 `br_user_app` (或在用户配置中校验)。



### 3.2 AI 任务处理流 (Async Job)

1. **触发**：日记保存成功，且 `br_diary_app` 中开启了 AI 插件。
2. **入队**：插入任务到 `br_job_queue`。
3. **消费**：Worker 读取 `br_app_ai` 的 Prompt 模板，渲染日记内容，调用 AI 接口。
4. **存储**：分析结果存入 `br_app_result` (原 `br_ai_result`)。

---

## 4. 扩展性说明

* **增加新 AI 插件**：在 `br_app` 注册，在 `br_app_ai` 写 Prompt。**无需改代码**。
* **增加备份/导出功能**：在 `br_app` 注册，`app_scope` 设为 `2` (User)，在 Go 中实现具体的备份逻辑。
* **避险设计**：
* **流量控制**：通过 `max_tokens` 限制 AI 成本。
* **隔离性**：即使 AI 服务宕机，不影响日记本基础读写。
* **数据安全**：备份 App 独立运行，支持用户配置非美国节点的存储（如避开监管风险）。



---

## 5. 结论

本方案通过 **“功能注册 + 维度路由 + 权限映射”** 的模式，既保证了系统初期的极简，又预留了未来成为“日记 App 平台”的空间。

**下一步：**
如果你认可这个方案，我可以为你生成这一整套表的 **SQL 初始脚本**，或者帮你写出 Go 语言中负责 **“App 开启状态校验”** 的核心函数代码。你想先看哪一个？