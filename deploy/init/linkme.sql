-- MySQL dump 10.13  Distrib 8.2.0, for macos13 (arm64)
--
-- Host: localhost    Database: linkme
-- ------------------------------------------------------
-- Server version	8.2.0

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `casbin_rule`
--

DROP TABLE IF EXISTS `casbin_rule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `casbin_rule` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) DEFAULT NULL,
  `v0` varchar(100) DEFAULT NULL,
  `v1` varchar(100) DEFAULT NULL,
  `v2` varchar(100) DEFAULT NULL,
  `v3` varchar(100) DEFAULT NULL,
  `v4` varchar(100) DEFAULT NULL,
  `v5` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_casbin_rule` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `casbin_rule`
--

LOCK TABLES `casbin_rule` WRITE;
/*!40000 ALTER TABLE `casbin_rule` DISABLE KEYS */;
INSERT INTO `casbin_rule` VALUES (25,'p','10308038636343296','/api/activity/recent','GET','','',''),(2,'p','10308038636343296','/api/checks/approve','POST','','',''),(4,'p','10308038636343296','/api/checks/detail','GET','','',''),(1,'p','10308038636343296','/api/checks/list','POST','','',''),(3,'p','10308038636343296','/api/checks/reject','POST','','',''),(24,'p','10308038636343296','/api/checks/stats','GET','','',''),(5,'p','10308038636343296','/api/permissions/assign','POST','','',''),(13,'p','10308038636343296','/api/permissions/assign_role','POST','','',''),(6,'p','10308038636343296','/api/permissions/list','GET','','',''),(7,'p','10308038636343296','/api/permissions/remove','DELETE','','',''),(16,'p','10308038636343296','/api/permissions/remove_role','DELETE','','',''),(8,'p','10308038636343296','/api/plate/create','POST','','',''),(10,'p','10308038636343296','/api/plate/delete/:plateId','DELETE','','',''),(9,'p','10308038636343296','/api/plate/list','POST','','',''),(26,'p','10308038636343296','/api/plate/update','POST','','',''),(21,'p','10308038636343296','/api/posts/detail_post/:postId','GET','','',''),(19,'p','10308038636343296','/api/posts/list_post','POST','','',''),(23,'p','10308038636343296','/api/posts/stats','GET','','',''),(11,'p','10308038636343296','/api/users/list','POST','','',''),(22,'p','10308038636343296','/api/users/stats','GET','','','');
/*!40000 ALTER TABLE `casbin_rule` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `checks`
--

DROP TABLE IF EXISTS `checks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `checks` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `post_id` bigint NOT NULL,
  `content` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `title` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `author_id` bigint DEFAULT NULL,
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'Pending',
  `remark` text COLLATE utf8mb4_unicode_ci,
  `created_at` bigint NOT NULL,
  `updated_at` bigint NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_checks_author` (`author_id`),
  KEY `idx_checks_updated_at` (`updated_at`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `checks`
--

LOCK TABLES `checks` WRITE;
/*!40000 ALTER TABLE `checks` DISABLE KEYS */;
/*!40000 ALTER TABLE `checks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `interactives`
--

DROP TABLE IF EXISTS `interactives`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `interactives` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `biz_id` bigint DEFAULT NULL,
  `biz_name` varchar(128) DEFAULT NULL,
  `read_count` bigint DEFAULT NULL,
  `like_count` bigint DEFAULT NULL,
  `collect_count` bigint DEFAULT NULL,
  `updated_at` bigint NOT NULL,
  `created_at` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `biz_type_id` (`biz_id`,`biz_name`),
  KEY `idx_interactives_update_time` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `interactives`
--

LOCK TABLES `interactives` WRITE;
/*!40000 ALTER TABLE `interactives` DISABLE KEYS */;
/*!40000 ALTER TABLE `interactives` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `plates`
--

DROP TABLE IF EXISTS `plates`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `plates` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci,
  `created_at` bigint DEFAULT NULL,
  `updated_at` bigint DEFAULT NULL,
  `deleted_at` bigint DEFAULT NULL,
  `deleted` tinyint(1) DEFAULT '0',
  `uid` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_plates_name` (`name`),
  KEY `idx_plates_uid` (`uid`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `plates`
--

LOCK TABLES `plates` WRITE;
/*!40000 ALTER TABLE `plates` DISABLE KEYS */;
INSERT INTO `plates` VALUES (1,'golang板块','golang学习...',1719917266377,1719917266377,0,0,10308038636343296),(2,'123','123',1720351398637,1720352037362,1720352037362,1,10308038636343296),(3,'1234','1234',1720352217856,1720352229176,1720352229176,1,10308038636343296),(4,'java板块123','java学习...',1720352234187,1720352666400,0,0,10308038636343296);
/*!40000 ALTER TABLE `plates` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `posts`
--

DROP TABLE IF EXISTS `posts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `posts` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `title` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `content` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` bigint NOT NULL,
  `updated_at` bigint NOT NULL,
  `deleted_at` bigint DEFAULT NULL,
  `deleted` tinyint(1) DEFAULT '0',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT 'draft',
  `author_id` bigint DEFAULT NULL,
  `slug` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `category_id` bigint DEFAULT NULL,
  `plate_id` bigint DEFAULT NULL,
  `tags` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `comment_count` bigint DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_posts_slug` (`slug`),
  KEY `idx_posts_plate_id` (`plate_id`),
  KEY `idx_posts_updated_time` (`updated_at`),
  KEY `idx_posts_deleted_time` (`deleted_at`),
  KEY `idx_posts_author` (`author_id`),
  KEY `idx_posts_category_id` (`category_id`),
  CONSTRAINT `fk_plates_posts` FOREIGN KEY (`plate_id`) REFERENCES `plates` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `posts`
--

LOCK TABLES `posts` WRITE;
/*!40000 ALTER TABLE `posts` DISABLE KEYS */;
/*!40000 ALTER TABLE `posts` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `profiles`
--

DROP TABLE IF EXISTS `profiles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `profiles` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `nick_name` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `avatar` text COLLATE utf8mb4_unicode_ci,
  `about` text COLLATE utf8mb4_unicode_ci,
  `birthday` varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_profiles_user_id` (`user_id`),
  CONSTRAINT `fk_users_profile` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `profiles`
--

LOCK TABLES `profiles` WRITE;
/*!40000 ALTER TABLE `profiles` DISABLE KEYS */;
INSERT INTO `profiles` VALUES (1,10308038636343296,'admin','admin12341','admin','2020-02-02');
/*!40000 ALTER TABLE `profiles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `recent_activities`
--

DROP TABLE IF EXISTS `recent_activities`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `recent_activities` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `description` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `time` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `recent_activities`
--

LOCK TABLES `recent_activities` WRITE;
/*!40000 ALTER TABLE `recent_activities` DISABLE KEYS */;
/*!40000 ALTER TABLE `recent_activities` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_collection_bizs`
--

DROP TABLE IF EXISTS `user_collection_bizs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_collection_bizs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `uid` bigint DEFAULT NULL,
  `biz_id` bigint DEFAULT NULL,
  `biz_name` varchar(255) DEFAULT NULL,
  `status` bigint DEFAULT NULL,
  `collection_id` bigint DEFAULT NULL,
  `updated_at` bigint NOT NULL,
  `created_at` bigint DEFAULT NULL,
  `deleted` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_user_collection_bizs_collection_id` (`collection_id`),
  KEY `idx_user_collection_bizs_update_time` (`updated_at`),
  KEY `idx_user_collection_bizs_uid` (`uid`),
  KEY `idx_user_collection_bizs_biz_id` (`biz_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_collection_bizs`
--

LOCK TABLES `user_collection_bizs` WRITE;
/*!40000 ALTER TABLE `user_collection_bizs` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_collection_bizs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_like_bizs`
--

DROP TABLE IF EXISTS `user_like_bizs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_like_bizs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `uid` bigint DEFAULT NULL,
  `biz_id` bigint DEFAULT NULL,
  `biz_name` varchar(255) DEFAULT NULL,
  `status` bigint DEFAULT NULL,
  `updated_at` bigint NOT NULL,
  `created_at` bigint DEFAULT NULL,
  `deleted` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_user_like_bizs_uid` (`uid`),
  KEY `idx_user_like_bizs_biz_id` (`biz_id`),
  KEY `idx_user_like_bizs_update_time` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_like_bizs`
--

LOCK TABLES `user_like_bizs` WRITE;
/*!40000 ALTER TABLE `user_like_bizs` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_like_bizs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` bigint DEFAULT NULL,
  `updated_at` bigint DEFAULT NULL,
  `deleted_at` bigint DEFAULT NULL,
  `nickname` varchar(50) DEFAULT NULL,
  `password_hash` longtext NOT NULL,
  `birthday` datetime DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL,
  `phone` varchar(15) DEFAULT NULL,
  `about` longtext,
  `deleted` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_email` (`email`),
  UNIQUE KEY `idx_users_phone` (`phone`),
  KEY `idx_users_deleted_time` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=18281565326938146 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (10308038636343296,1718025027925,1718025027925,0,'Bamboo','$2a$10$LKLuaqEF38fA..XHCfRI.OtCxo58vYub/Fq.40X4CaaV.IqemrUd2','2024-07-03 00:00:00','admin','123456','this is LinkMe',0),(18281565326938132,NULL,1,NULL,'aaa','123',NULL,'adsfasdf','1','asdfasd',0),(18281565326938133,NULL,2,NULL,'bbb','123',NULL,'asdfasdfadsfasdf','2','tgsadgasdf',0),(18281565326938134,NULL,3,NULL,'ccc','123',NULL,'qwereqwr','3','asdfasdfasd',0);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `v_code_sms_logs`
--

DROP TABLE IF EXISTS `v_code_sms_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `v_code_sms_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `sms_id` bigint DEFAULT NULL,
  `sms_type` longtext,
  `mobile` longtext,
  `v_code` longtext,
  `driver` longtext,
  `status` bigint DEFAULT NULL,
  `status_code` longtext,
  `created_at` bigint NOT NULL,
  `updated_at` bigint NOT NULL,
  `deleted_at` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_v_code_sms_logs_updated_time` (`updated_at`),
  KEY `idx_v_code_sms_logs_deleted_time` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `v_code_sms_logs`
--

LOCK TABLES `v_code_sms_logs` WRITE;
/*!40000 ALTER TABLE `v_code_sms_logs` DISABLE KEYS */;
/*!40000 ALTER TABLE `v_code_sms_logs` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-07-07 19:46:58
