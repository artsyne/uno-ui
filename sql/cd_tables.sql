-- MySQL dump 10.13  Distrib 8.0.38, for macos14 (arm64)
--
-- Host: localhost    Database: cd
-- ------------------------------------------------------
-- Server version	8.3.0

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
  `comment_id` bigint unsigned NOT NULL,
  `comment_type` tinyint unsigned NOT NULL,
  `owner_id` bigint unsigned NOT NULL,
  `content` varchar(5000) NOT NULL DEFAULT '',
  `gmt_create` datetime NOT NULL,
  `add_user_id` bigint unsigned NOT NULL,
  `add_user_name` varchar(100) NOT NULL DEFAULT '',
  `add_user_ip` varchar(128) NOT NULL DEFAULT '',
  PRIMARY KEY (`comment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_diary`
--

DROP TABLE IF EXISTS `br_diary`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_diary` (
  `diary_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(100) NOT NULL DEFAULT '',
  `remark` varchar(500) NOT NULL DEFAULT '',
  `access_url` varchar(255) NOT NULL DEFAULT '',
  `cover_url` varchar(255) NOT NULL DEFAULT '',
  `cover_mid_url` varchar(255) NOT NULL DEFAULT '',
  `like_num` int unsigned NOT NULL DEFAULT '0',
  `member_num` int unsigned NOT NULL DEFAULT '0',
  `view_num` int unsigned NOT NULL DEFAULT '0',
  `view_setting` tinyint unsigned NOT NULL DEFAULT '1',
  `access_pwd` varchar(255) NOT NULL DEFAULT '',
  `social_setting` tinyint unsigned NOT NULL DEFAULT '1',
  `add_utc_time` bigint unsigned NOT NULL DEFAULT '0',
  `add_user_id` bigint unsigned NOT NULL,
  `add_user_name` varchar(100) NOT NULL DEFAULT '',
  `info_num` int unsigned NOT NULL DEFAULT '0',
  `last_action_utc_time` bigint unsigned NOT NULL DEFAULT '0',
  `last_action_user_id` bigint unsigned NOT NULL DEFAULT '0',
  `last_action_user_name` varchar(100) NOT NULL DEFAULT '',
  `theme_id` tinyint unsigned NOT NULL DEFAULT '16',
  `status` tinyint unsigned NOT NULL,
  PRIMARY KEY (`diary_id`)
) ENGINE=InnoDB AUTO_INCREMENT=44 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_feed`
--

DROP TABLE IF EXISTS `br_feed`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_feed` (
  `feed_id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `user_name` varchar(100) NOT NULL DEFAULT '',
  `feed_type_id` int NOT NULL,
  `body` varchar(4000) NOT NULL DEFAULT '',
  `add_time` datetime NOT NULL,
  `add_ip` varchar(39) NOT NULL DEFAULT '',
  PRIMARY KEY (`feed_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_inbox_detail`
--

DROP TABLE IF EXISTS `br_inbox_detail`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_inbox_detail` (
  `msg_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `session_id` bigint unsigned NOT NULL,
  `sender_id` bigint unsigned NOT NULL,
  `sender` varchar(100) NOT NULL DEFAULT '',
  `receiver_id` bigint unsigned NOT NULL,
  `receiver` varchar(100) NOT NULL DEFAULT '',
  `msg` varchar(2000) NOT NULL DEFAULT '',
  `gmt_create` datetime NOT NULL,
  `is_html` tinyint unsigned NOT NULL DEFAULT '0',
  `is_sdel` tinyint unsigned NOT NULL DEFAULT '0',
  `is_rdel` tinyint unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`msg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_inbox_header`
--

DROP TABLE IF EXISTS `br_inbox_header`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_inbox_header` (
  `session_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `last_sender_id` bigint unsigned NOT NULL,
  `last_sender` varchar(100) NOT NULL DEFAULT '',
  `last_receiver_id` bigint unsigned NOT NULL,
  `last_receiver` varchar(100) NOT NULL DEFAULT '',
  `gmt_create` datetime NOT NULL,
  `gmt_modified` datetime NOT NULL,
  `last_msg` varchar(2000) NOT NULL DEFAULT '',
  `msg_num` int unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`session_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_inbox_stat`
--

DROP TABLE IF EXISTS `br_inbox_stat`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_inbox_stat` (
  `session_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `with_user_id` bigint unsigned NOT NULL,
  `unread_num` int unsigned NOT NULL DEFAULT '0',
  `is_del` tinyint unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`session_id`,`user_id`,`with_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_job_history`
--

DROP TABLE IF EXISTS `br_job_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_job_history` (
  `job_id` bigint unsigned NOT NULL,
  `job_type` varchar(200) NOT NULL,
  `body` varchar(4000) NOT NULL DEFAULT '',
  `try_num` tinyint unsigned NOT NULL,
  `result` tinyint unsigned NOT NULL,
  `plan_utc_time` bigint unsigned NOT NULL,
  `done_utc_time` bigint unsigned NOT NULL,
  PRIMARY KEY (`job_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_job_queue`
--

DROP TABLE IF EXISTS `br_job_queue`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_job_queue` (
  `job_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `job_type` varchar(255) NOT NULL,
  `body` varchar(5000) NOT NULL,
  `try_num` tinyint unsigned NOT NULL,
  `result` tinyint unsigned NOT NULL,
  `add_utc_time` bigint unsigned NOT NULL,
  `plan_utc_time` bigint unsigned NOT NULL,
  PRIMARY KEY (`job_id`),
  KEY `job_type` (`job_type`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_like_history`
--

DROP TABLE IF EXISTS `br_like_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_like_history` (
  `user_id` bigint unsigned NOT NULL,
  `like_id` bigint unsigned NOT NULL,
  `like_type` tinyint unsigned NOT NULL,
  PRIMARY KEY (`user_id`,`like_id`,`like_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_my_follow`
--

DROP TABLE IF EXISTS `br_my_follow`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_my_follow` (
  `user_id` bigint unsigned NOT NULL,
  `follow_user_id` bigint unsigned NOT NULL,
  `follow_user_name` varchar(100) NOT NULL DEFAULT '',
  `remark` varchar(500) NOT NULL DEFAULT '',
  `gmt_time` datetime NOT NULL,
  PRIMARY KEY (`user_id`,`follow_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_online`
--

DROP TABLE IF EXISTS `br_online`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_online` (
  `user_id` bigint unsigned NOT NULL,
  `first_activity_ip` varchar(39) NOT NULL DEFAULT '',
  `first_activity_time` datetime NOT NULL,
  `last_activity_ip` varchar(39) NOT NULL DEFAULT '',
  `last_activity_time` datetime NOT NULL,
  `longitude` decimal(10,0) NOT NULL DEFAULT '0',
  `latitude` decimal(10,0) NOT NULL DEFAULT '0',
  `session_id` varchar(100) NOT NULL DEFAULT '',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_post`
--

DROP TABLE IF EXISTS `br_post`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_post` (
  `post_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `diary_id` bigint unsigned NOT NULL DEFAULT '0',
  `is_top` tinyint unsigned NOT NULL DEFAULT '0',
  `info_type` tinyint unsigned NOT NULL DEFAULT '0',
  `title` varchar(200) NOT NULL DEFAULT '',
  `summary` varchar(1024) NOT NULL DEFAULT '',
  `body` text NOT NULL,
  `view_num` int unsigned NOT NULL DEFAULT '0',
  `like_num` int unsigned NOT NULL DEFAULT '0',
  `comment_num` int unsigned NOT NULL DEFAULT '0',
  `add_user_id` bigint unsigned NOT NULL DEFAULT '0',
  `add_user_name` varchar(100) NOT NULL DEFAULT '',
  `add_user_ip` varchar(128) NOT NULL DEFAULT '',
  `add_utc_time` bigint unsigned NOT NULL DEFAULT '0',
  `update_utc_time` bigint unsigned NOT NULL DEFAULT '0',
  `status` tinyint unsigned NOT NULL DEFAULT '1',
  PRIMARY KEY (`post_id`)
) ENGINE=InnoDB AUTO_INCREMENT=94 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_space_member`
--

DROP TABLE IF EXISTS `br_space_member`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_space_member` (
  `space_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `user_name` varchar(100) NOT NULL DEFAULT '',
  `is_admin` tinyint unsigned NOT NULL,
  `join_time` datetime NOT NULL,
  `post_num` int unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`space_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_sys_param`
--

DROP TABLE IF EXISTS `br_sys_param`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_sys_param` (
  `param_name` varchar(50) NOT NULL DEFAULT '',
  `param_value` varchar(2000) NOT NULL DEFAULT '',
  `remark` varchar(500) NOT NULL DEFAULT '',
  PRIMARY KEY (`param_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_user`
--

DROP TABLE IF EXISTS `br_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_user` (
  `user_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varchar(100) NOT NULL DEFAULT '',
  `full_name` varchar(128) DEFAULT '',
  `user_group` tinyint unsigned NOT NULL DEFAULT '1',
  `pwd` varchar(255) NOT NULL DEFAULT '',
  `pwd_salt` varchar(64) NOT NULL DEFAULT '',
  `access_url` varchar(255) NOT NULL DEFAULT '',
  `face_big_url` varchar(255) NOT NULL DEFAULT '',
  `face_mid_url` varchar(255) NOT NULL DEFAULT '',
  `face_small_url` varchar(255) NOT NULL DEFAULT '',
  `email` varchar(128) NOT NULL DEFAULT '',
  `reg_ip` varchar(128) NOT NULL DEFAULT '',
  `reg_utc_time` bigint NOT NULL DEFAULT '0',
  `remark` varchar(1000) NOT NULL DEFAULT '',
  `last_ip` varchar(128) NOT NULL DEFAULT '',
  `last_login_utc_time` bigint NOT NULL DEFAULT '0',
  `time_zone` varchar(64) NOT NULL DEFAULT '',
  `lang` varchar(8) NOT NULL DEFAULT '',
  `status` tinyint unsigned NOT NULL DEFAULT '1',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `br_user_stat`
--

DROP TABLE IF EXISTS `br_user_stat`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `br_user_stat` (
  `user_id` bigint unsigned NOT NULL,
  `follow_num` int unsigned NOT NULL DEFAULT '0',
  `info_num` int unsigned NOT NULL DEFAULT '0',
  `score` int unsigned NOT NULL DEFAULT '0',
  `fans_num` int unsigned NOT NULL DEFAULT '0',
  `like_num` int unsigned NOT NULL DEFAULT '0',
  `gmt_modified` datetime NOT NULL,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
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
  `valid_code` varchar(200) NOT NULL DEFAULT '',
  `valid_type_id` tinyint unsigned NOT NULL,
  `expire_utc_time` bigint unsigned NOT NULL,
  `create_utc_time` bigint NOT NULL,
  PRIMARY KEY (`valid_code_id`)
) ENGINE=InnoDB AUTO_INCREMENT=121 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-08-16 14:50:50
