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
-- Table structure for table `br_comment`
--

DROP TABLE IF EXISTS `br_comment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_comment` (
  `comment_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `entry_id` bigint unsigned NOT NULL,
  `content` varchar(5000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `add_utc_time` bigint unsigned NOT NULL DEFAULT '0',
  `add_user_id` bigint unsigned NOT NULL,
  `add_user_ip` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`comment_id`),
  KEY `uk_entry_time` (`entry_id`,`add_utc_time` DESC)
) ENGINE=InnoDB AUTO_INCREMENT=817 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_diary`
--

DROP TABLE IF EXISTS `br_diary`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_diary` (
  `diary_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `remark` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `access_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `cover_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `cover_mid_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `like_num` int unsigned NOT NULL DEFAULT '0',
  `member_num` int unsigned NOT NULL DEFAULT '0',
  `view_num` int unsigned NOT NULL DEFAULT '0',
  `view_setting` tinyint unsigned NOT NULL DEFAULT '1',
  `access_pwd` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `social_setting` tinyint unsigned NOT NULL DEFAULT '1',
  `add_utc_time` bigint unsigned NOT NULL DEFAULT '0',
  `add_user_id` bigint unsigned NOT NULL,
  `info_num` int unsigned NOT NULL DEFAULT '0',
  `theme_id` tinyint unsigned NOT NULL DEFAULT '16',
  `status` tinyint unsigned NOT NULL,
  PRIMARY KEY (`diary_id`),
  KEY `idx_add_user_id` (`add_user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1761 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_diary_member`
--

DROP TABLE IF EXISTS `br_diary_member`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_diary_member` (
  `diary_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `is_owner` tinyint unsigned NOT NULL DEFAULT '0',
  `join_utc_time` bigint unsigned NOT NULL DEFAULT '0',
  `post_num` int unsigned NOT NULL DEFAULT '0',
  `status` tinyint unsigned NOT NULL DEFAULT '2',
  PRIMARY KEY (`diary_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_entry`
--

DROP TABLE IF EXISTS `br_entry`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_entry` (
  `entry_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `diary_id` bigint unsigned NOT NULL DEFAULT '0',
  `is_top` tinyint unsigned NOT NULL DEFAULT '0',
  `info_type` tinyint unsigned NOT NULL DEFAULT '0',
  `title` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `summary` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `body` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `view_num` int unsigned NOT NULL DEFAULT '0',
  `like_num` int unsigned NOT NULL DEFAULT '0',
  `comment_num` int unsigned NOT NULL DEFAULT '0',
  `add_user_id` bigint unsigned NOT NULL DEFAULT '0',
  `add_user_ip` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `add_utc_time` bigint unsigned NOT NULL DEFAULT '0',
  `update_utc_time` bigint unsigned NOT NULL DEFAULT '0',
  `status` tinyint unsigned NOT NULL DEFAULT '1',
  PRIMARY KEY (`entry_id`),
  KEY `uk_diary_time` (`diary_id`,`add_utc_time` DESC)
) ENGINE=InnoDB AUTO_INCREMENT=2333 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_job_history`
--

DROP TABLE IF EXISTS `br_job_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_job_history` (
  `job_id` bigint unsigned NOT NULL,
  `job_type` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `body` varchar(4000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `try_num` tinyint unsigned NOT NULL,
  `result` tinyint unsigned NOT NULL,
  `plan_utc_time` bigint unsigned NOT NULL,
  `done_utc_time` bigint unsigned NOT NULL,
  PRIMARY KEY (`job_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_job_queue`
--

DROP TABLE IF EXISTS `br_job_queue`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_job_queue` (
  `job_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `job_type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `body` varchar(5000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `try_num` tinyint unsigned NOT NULL,
  `result` tinyint unsigned NOT NULL,
  `add_utc_time` bigint unsigned NOT NULL,
  `plan_utc_time` bigint unsigned NOT NULL,
  PRIMARY KEY (`job_id`),
  KEY `job_type` (`job_type`)
) ENGINE=InnoDB AUTO_INCREMENT=2239 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_like`
--

DROP TABLE IF EXISTS `br_like`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_like` (
  `like_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `owner_id` bigint unsigned NOT NULL,
  `like_type` tinyint unsigned NOT NULL,
  PRIMARY KEY (`like_id`),
  UNIQUE KEY `uk_user_owner_type` (`user_id`,`owner_id`,`like_type`)
) ENGINE=InnoDB AUTO_INCREMENT=165 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_user`
--

DROP TABLE IF EXISTS `br_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_user` (
  `user_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `handle` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `display_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '',
  `user_group` tinyint unsigned NOT NULL DEFAULT '1',
  `pwd` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '',
  `profile_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `face_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `email_verified` tinyint NOT NULL DEFAULT '0',
  `reg_ip` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `created_utc_time` bigint NOT NULL DEFAULT '0',
  `modified_utc_time` bigint NOT NULL DEFAULT '0',
  `remark` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `last_ip` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `last_login_utc_time` bigint NOT NULL DEFAULT '0',
  `time_zone` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `lang` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `status` tinyint unsigned NOT NULL DEFAULT '1',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `handle_UNIQUE` (`handle`),
  UNIQUE KEY `email_UNIQUE` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=2373 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_user_oauth`
--

DROP TABLE IF EXISTS `br_user_oauth`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_user_oauth` (
  `oauth_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `provider` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `provider_user_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `created_utc_time` bigint NOT NULL,
  PRIMARY KEY (`oauth_id`),
  UNIQUE KEY `uk_provider_user` (`provider`,`provider_user_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_valid_code`
--

DROP TABLE IF EXISTS `br_valid_code`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_valid_code` (
  `valid_code_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `valid_code` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `valid_type_id` tinyint unsigned NOT NULL,
  `expire_utc_time` bigint unsigned NOT NULL,
  `create_utc_time` bigint NOT NULL,
  PRIMARY KEY (`valid_code_id`),
  KEY `idx_user` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2153 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;