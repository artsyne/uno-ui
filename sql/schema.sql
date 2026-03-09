-- MySQL dump 10.13  Distrib 8.0.40, for macos14 (arm64)
--
-- Host: localhost    Database: cd
-- ------------------------------------------------------
-- Server version	9.2.0

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `br_app`
--

DROP TABLE IF EXISTS `br_app`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_app` (
  `app_id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键ID',
  `app_key` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '插件唯一标识',
  `app_name` varchar(100) COLLATE utf8mb4_general_ci NOT NULL COMMENT '插件显示名称',
  `app_scope` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '作用域: 1-日记本维度(Diary), 2-用户维度(User)',
  `app_category` tinyint unsigned NOT NULL COMMENT '分类: 1-AI类, 2-导出类, 3-备份类, 4-数据处理',
  `config_json` text COLLATE utf8mb4_general_ci COMMENT '通用配置(JSON格式)',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态(1:启用 2:禁用)',
  `add_utc_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建时间(UTC时间戳)',
  PRIMARY KEY (`app_id`),
  UNIQUE KEY `uk_app_key` (`app_key`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='扩展应用注册主表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_app_ai`
--

DROP TABLE IF EXISTS `br_app_ai`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_app_ai` (
  `app_key` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '关联 br_app.app_key',
  `model_name` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'gpt-4o' COMMENT '调用的AI模型名称',
  `system_prompt` text COLLATE utf8mb4_general_ci COMMENT '系统提示词: 定义AI的角色和全局规则',
  `user_prompt_template` text COLLATE utf8mb4_general_ci COMMENT '用户提示词模板: 包含占位符(如 {{.Content}})的指令',
  `response_format` text COLLATE utf8mb4_general_ci COMMENT '期望的返回JSON格式示例, 用于约束AI输出',
  `temperature` float NOT NULL DEFAULT '0.7' COMMENT '发散度(0-2.0): 值越高越随机, 值越低越严谨',
  `max_tokens` int unsigned NOT NULL DEFAULT '2000' COMMENT '单次调用最大消耗Token数',
  PRIMARY KEY (`app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='AI类插件专用配置表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_comment`
--

DROP TABLE IF EXISTS `br_comment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_comment` (
  `comment_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '评论ID',
  `entry_id` bigint unsigned NOT NULL COMMENT '关联日记ID',
  `content` varchar(5000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '评论内容',
  `add_utc_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建时间(UTC毫秒)',
  `add_user_id` bigint unsigned NOT NULL COMMENT '创建用户ID',
  `add_user_ip` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '创建用户IP地址',
  PRIMARY KEY (`comment_id`),
  KEY `uk_entry_time` (`entry_id`,`add_utc_time` DESC),
  KEY `idx_entry_comment` (`entry_id`,`comment_id` DESC)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='评论表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_diary`
--

DROP TABLE IF EXISTS `br_diary`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_diary` (
  `diary_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '日记本ID',
  `title` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '日记本标题',
  `remark` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '日记本描述/备注',
  `access_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '访问URL地址',
  `cover_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '封面图片URL',
  `cover_mid_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '中等尺寸封面图片URL',
  `like_num` int unsigned NOT NULL DEFAULT '0' COMMENT '点赞数量',
  `member_num` int unsigned NOT NULL DEFAULT '0' COMMENT '成员数量',
  `view_num` int unsigned NOT NULL DEFAULT '0' COMMENT '浏览次数',
  `view_setting` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '查看权限设置',
  `access_pwd` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '访问密码',
  `social_setting` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '社交功能设置',
  `add_utc_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建时间(UTC时间戳)',
  `add_user_id` bigint unsigned NOT NULL COMMENT '创建者用户ID',
  `info_num` int unsigned NOT NULL DEFAULT '0' COMMENT '日记数量',
  `theme_id` tinyint unsigned NOT NULL DEFAULT '16' COMMENT '主题ID',
  `status` tinyint unsigned NOT NULL COMMENT '状态(0:删除 1:正常)',
  PRIMARY KEY (`diary_id`),
  KEY `idx_user_diary` (`add_user_id`,`diary_id` DESC)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='日记本表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_diary_app`
--

DROP TABLE IF EXISTS `br_diary_app`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_diary_app` (
  `unique_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '唯一自增ID',
  `diary_id` bigint unsigned NOT NULL COMMENT '关联 br_diary.diary_id',
  `app_key` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '关联 br_app.app_key',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态(1:开启 2:关闭)',
  `add_utc_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '为此日记本开启应用的时间(UTC时间戳)',
  PRIMARY KEY (`unique_id`),
  UNIQUE KEY `uk_diary_app` (`diary_id`,`app_key`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='日记本维度应用开启状态表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_diary_author`
--

DROP TABLE IF EXISTS `br_diary_author`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_diary_author` (
  `diary_id` bigint unsigned NOT NULL COMMENT '关联日记本ID',
  `user_id` bigint unsigned NOT NULL COMMENT '关联用户ID',
  `is_admin` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否管理员(0:否 1:是)',
  `join_utc_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '加入时间(UTC时间戳)',
  `post_num` int unsigned NOT NULL DEFAULT '0' COMMENT '发表文章数',
  `status` tinyint unsigned NOT NULL DEFAULT '2' COMMENT '状态(1:活跃 2:待确认)',
  PRIMARY KEY (`diary_id`,`user_id`),
  KEY `idx_user_status_join` (`user_id`,`status`,`join_utc_time` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='日记本作者表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_entry`
--

DROP TABLE IF EXISTS `br_entry`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_entry` (
  `entry_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '日记ID',
  `diary_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '关联日记本ID',
  `is_top` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否置顶(0:否 1:是)',
  `info_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '信息类型(1:文本 2:图片 3:视频 4:音频 5:文件 6:URL)',
  `title` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '日记标题',
  `summary` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '日记摘要',
  `body` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '日记内容',
  `view_num` int unsigned NOT NULL DEFAULT '0' COMMENT '浏览次数',
  `like_num` int unsigned NOT NULL DEFAULT '0' COMMENT '点赞次数',
  `comment_num` int unsigned NOT NULL DEFAULT '0' COMMENT '评论次数',
  `add_user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建用户ID',
  `add_user_ip` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '创建用户IP',
  `add_utc_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建时间(UTC时间戳)',
  `update_utc_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '更新时间(UTC时间戳)',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态(0:删除 1:正常 2:草稿)',
  PRIMARY KEY (`entry_id`),
  KEY `uk_diary_time` (`diary_id`,`add_utc_time` DESC),
  KEY `idx_diary_status_entry` (`diary_id`,`status`,`entry_id` DESC),
  KEY `idx_user_status_update` (`add_user_id`,`status`,`update_utc_time` DESC)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='日记表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_job_history`
--

DROP TABLE IF EXISTS `br_job_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_job_history` (
  `job_id` bigint unsigned NOT NULL COMMENT '任务ID',
  `job_type` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '任务类型',
  `body` varchar(4000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务内容(JSON格式)',
  `try_num` tinyint unsigned NOT NULL COMMENT '尝试次数',
  `result` tinyint unsigned NOT NULL COMMENT '执行结果(0:失败 1:成功)',
  `plan_utc_time` bigint unsigned NOT NULL COMMENT '计划执行时间(UTC时间戳)',
  `done_utc_time` bigint unsigned NOT NULL COMMENT '实际完成时间(UTC时间戳)',
  PRIMARY KEY (`job_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='任务历史记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_job_queue`
--

DROP TABLE IF EXISTS `br_job_queue`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_job_queue` (
  `job_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '任务ID',
  `job_type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '任务类型',
  `body` varchar(5000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '任务内容(JSON格式)',
  `try_num` tinyint unsigned NOT NULL COMMENT '尝试次数',
  `result` tinyint unsigned NOT NULL COMMENT '执行结果(0:失败 1:成功)',
  `add_utc_time` bigint unsigned NOT NULL COMMENT '添加时间(UTC时间戳)',
  `plan_utc_time` bigint unsigned NOT NULL COMMENT '计划执行时间(UTC时间戳)',
  PRIMARY KEY (`job_id`),
  KEY `idx_job_type_id` (`job_type`,`job_id` DESC)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='任务队列表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_like`
--

DROP TABLE IF EXISTS `br_like`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_like` (
  `like_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '点赞ID',
  `user_id` bigint unsigned NOT NULL COMMENT '点赞用户ID',
  `owner_id` bigint unsigned NOT NULL COMMENT '被点赞对象ID(日记或日记本)',
  `like_type` tinyint unsigned NOT NULL COMMENT '点赞类型(1:日记 2:日记本)',
  PRIMARY KEY (`like_id`),
  UNIQUE KEY `uk_user_owner_type` (`user_id`,`owner_id`,`like_type`),
  KEY `idx_owner_type` (`owner_id`,`like_type`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='点赞表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_plan`
--

DROP TABLE IF EXISTS `br_plan`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_plan` (
  `plan_id` tinyint unsigned NOT NULL AUTO_INCREMENT COMMENT '套餐计划ID',
  `plan_name` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '套餐名称',
  `price_monthly` int unsigned NOT NULL DEFAULT '0' COMMENT '月付价格',
  `features` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '套餐功能(JSON格式)',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态(0:禁用 1:启用)',
  PRIMARY KEY (`plan_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='套餐计划表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_user`
--

DROP TABLE IF EXISTS `br_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_user` (
  `user_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `handle` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户唯一标识符(句柄)',
  `display_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '用户显示名称',
  `user_group` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '用户组标识',
  `pwd` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '密码(加密后)',
  `profile_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '个人主页URL',
  `face_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '头像URL',
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '邮箱地址',
  `email_verified` tinyint NOT NULL DEFAULT '0' COMMENT '邮箱是否已验证(0:未验证 1:已验证)',
  `reg_ip` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '注册IP地址',
  `created_utc_time` bigint NOT NULL DEFAULT '0' COMMENT '创建时间(UTC时间戳)',
  `modified_utc_time` bigint NOT NULL DEFAULT '0' COMMENT '最后修改时间(UTC时间戳)',
  `remark` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '用户备注',
  `last_ip` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '最后登录IP地址',
  `last_login_utc_time` bigint NOT NULL DEFAULT '0' COMMENT '最后登录时间(UTC时间戳)',
  `time_zone` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户时区设置',
  `lang` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户语言偏好',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态(0:删除 1:正常)',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `handle_UNIQUE` (`handle`),
  UNIQUE KEY `email_UNIQUE` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_user_app`
--

DROP TABLE IF EXISTS `br_user_app`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_user_app` (
  `unique_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '唯一自增ID',
  `user_id` bigint unsigned NOT NULL COMMENT '关联 br_user.user_id',
  `app_key` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '关联 br_app.app_key',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态(1:开启 2:关闭)',
  `add_utc_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '开启时间(UTC时间戳)',
  PRIMARY KEY (`unique_id`),
  UNIQUE KEY `uk_user_app` (`user_id`,`app_key`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户应用开启状态表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_user_oauth`
--

DROP TABLE IF EXISTS `br_user_oauth`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_user_oauth` (
  `oauth_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'OAuth绑定ID',
  `user_id` bigint unsigned NOT NULL COMMENT '关联用户ID',
  `provider` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'OAuth提供商(如:google, github)',
  `provider_user_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '提供商端的用户ID',
  `created_utc_time` bigint NOT NULL COMMENT '绑定时间(UTC时间戳)',
  PRIMARY KEY (`oauth_id`),
  UNIQUE KEY `uk_provider_user` (`provider`,`provider_user_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户OAuth绑定表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_user_subscription`
--

DROP TABLE IF EXISTS `br_user_subscription`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_user_subscription` (
  `user_id` bigint unsigned NOT NULL COMMENT '关联用户ID',
  `plan_id` tinyint unsigned NOT NULL COMMENT '关联套餐计划ID',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '订阅状态(0:取消 1:激活)',
  `start_utc_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '订阅开始时间(UTC时间戳)',
  `expire_utc_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '订阅到期时间(UTC时间戳)',
  `last_payment_utc_time` bigint unsigned NOT NULL DEFAULT '0' COMMENT '最后支付时间(UTC时间戳)',
  PRIMARY KEY (`user_id`),
  KEY `idx_expire_status` (`expire_utc_time`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户订阅表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_valid_code`
--

DROP TABLE IF EXISTS `br_valid_code`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_valid_code` (
  `valid_code_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '验证码ID',
  `user_id` bigint unsigned NOT NULL COMMENT '关联用户ID',
  `valid_code` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '验证码',
  `valid_type_id` tinyint unsigned NOT NULL COMMENT '验证码类型(1:注册 2:找回密码 3:绑定邮箱)',
  `expire_utc_time` bigint unsigned NOT NULL COMMENT '过期时间(UTC时间戳)',
  `create_utc_time` bigint NOT NULL COMMENT '创建时间(UTC时间戳)',
  PRIMARY KEY (`valid_code_id`),
  KEY `idx_user` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='验证码表';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-02-03 17:49:10
