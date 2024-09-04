-- create schema
CREATE DATABASE idiary
CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;

use idiary;

--
-- Table structure for table `br_app`
--

DROP TABLE IF EXISTS `br_app`;
CREATE TABLE `br_app` (
  `app_id` int unsigned NOT NULL AUTO_INCREMENT,
  `app_name` varchar(200) NOT NULL DEFAULT '',
  `intro` longtext NOT NULL,
  `url` varchar(255) NOT NULL DEFAULT '',
  `author` varchar(200) NOT NULL DEFAULT '',
  `cover_url` varchar(255) NOT NULL DEFAULT '',
  `client_id` varchar(100) NOT NULL DEFAULT '',
  `client_secret` varchar(255) NOT NULL DEFAULT '',
  `created_utc_time` bigint unsigned NOT NULL,
  `updated_utc_time` bigint unsigned NOT NULL,
  `status` tinyint unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`app_id`)
) ;


--
-- Table structure for table `br_app_user`
--

DROP TABLE IF EXISTS `br_app_user`;
CREATE TABLE `br_app_user` (
  `user_id` bigint unsigned NOT NULL,
  `app_id` bigint unsigned NOT NULL,
  `access_token` varchar(255) NOT NULL DEFAULT '',
  `refresh_token` varchar(255) NOT NULL DEFAULT '',
  `scope` varchar(255) NOT NULL DEFAULT '',
  `expires_utc_time` bigint NOT NULL,
  `created_utc_time` bigint NOT NULL,
  `updated_utc_time` bigint NOT NULL,
  `status` tinyint unsigned NOT NULL,
  PRIMARY KEY (`user_id`,`app_id`)
) ;


--
-- Table structure for table `br_comment`
--

DROP TABLE IF EXISTS `br_comment`;
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
) ;


--
-- Table structure for table `br_diary`
--

DROP TABLE IF EXISTS `br_diary`;
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
) ;


--
-- Table structure for table `br_feed`
--

DROP TABLE IF EXISTS `br_feed`;
CREATE TABLE `br_feed` (
  `feed_id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `user_name` varchar(100) NOT NULL DEFAULT '',
  `feed_type_id` int NOT NULL,
  `body` varchar(4000) NOT NULL DEFAULT '',
  `add_time` datetime NOT NULL,
  `add_ip` varchar(39) NOT NULL DEFAULT '',
  PRIMARY KEY (`feed_id`)
) ;


--
-- Table structure for table `br_follow`
--

DROP TABLE IF EXISTS `br_follow`;
CREATE TABLE `br_follow` (
  `user_id` bigint unsigned NOT NULL,
  `follow_user_id` bigint unsigned NOT NULL,
  `follow_user_name` varchar(100) NOT NULL DEFAULT '',
  `remark` varchar(500) NOT NULL DEFAULT '',
  `gmt_time` datetime NOT NULL,
  PRIMARY KEY (`user_id`,`follow_user_id`)
) ;


--
-- Table structure for table `br_job_history`
--

DROP TABLE IF EXISTS `br_job_history`;
CREATE TABLE `br_job_history` (
  `job_id` bigint unsigned NOT NULL,
  `job_type` varchar(200) NOT NULL,
  `body` varchar(4000) NOT NULL DEFAULT '',
  `try_num` tinyint unsigned NOT NULL,
  `result` tinyint unsigned NOT NULL,
  `plan_utc_time` bigint unsigned NOT NULL,
  `done_utc_time` bigint unsigned NOT NULL,
  PRIMARY KEY (`job_id`)
) ;


--
-- Table structure for table `br_job_queue`
--

DROP TABLE IF EXISTS `br_job_queue`;
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
) ;


--
-- Table structure for table `br_like`
--

DROP TABLE IF EXISTS `br_like`;
CREATE TABLE `br_like` (
  `user_id` bigint unsigned NOT NULL,
  `like_id` bigint unsigned NOT NULL,
  `like_type` tinyint unsigned NOT NULL,
  PRIMARY KEY (`user_id`,`like_id`,`like_type`)
) ;


--
-- Table structure for table `br_post`
--

DROP TABLE IF EXISTS `br_post`;
CREATE TABLE `br_post` (
  `post_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `diary_id` bigint unsigned NOT NULL DEFAULT '0',
  `is_top` tinyint unsigned NOT NULL DEFAULT '0',
  `info_type` tinyint unsigned NOT NULL DEFAULT '0',
  `title` varchar(200) NOT NULL DEFAULT '',
  `summary` varchar(1024) NOT NULL DEFAULT '',
  `body` mediumtext NOT NULL,
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
) ;


--
-- Table structure for table `br_space_member`
--

DROP TABLE IF EXISTS `br_space_member`;
CREATE TABLE `br_space_member` (
  `space_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `user_name` varchar(100) NOT NULL DEFAULT '',
  `is_admin` tinyint unsigned NOT NULL,
  `join_time` datetime NOT NULL,
  `post_num` int unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`space_id`,`user_id`)
) ;


--
-- Table structure for table `br_user`
--

DROP TABLE IF EXISTS `br_user`;
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
) ;


--
-- Table structure for table `br_user_stat`
--

DROP TABLE IF EXISTS `br_user_stat`;
CREATE TABLE `br_user_stat` (
  `user_id` bigint unsigned NOT NULL,
  `follow_num` int unsigned NOT NULL DEFAULT '0',
  `info_num` int unsigned NOT NULL DEFAULT '0',
  `fans_num` int unsigned NOT NULL DEFAULT '0',
  `like_num` int unsigned NOT NULL DEFAULT '0',
  `gmt_modified` datetime NOT NULL,
  PRIMARY KEY (`user_id`)
) ;


--
-- Table structure for table `br_valid_code`
--

DROP TABLE IF EXISTS `br_valid_code`;
CREATE TABLE `br_valid_code` (
  `valid_code_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `valid_code` varchar(200) NOT NULL DEFAULT '',
  `valid_type_id` tinyint unsigned NOT NULL,
  `expire_utc_time` bigint unsigned NOT NULL,
  `create_utc_time` bigint NOT NULL,
  PRIMARY KEY (`valid_code_id`)
) ;