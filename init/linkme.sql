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
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `casbin_rule`
--

LOCK TABLES `casbin_rule` WRITE;
/*!40000 ALTER TABLE `casbin_rule` DISABLE KEYS */;
INSERT INTO `casbin_rule` VALUES (2,'p','10308038636343296','/checks/approve','POST','','',''),(4,'p','10308038636343296','/checks/detail','GET','','',''),(1,'p','10308038636343296','/checks/list','GET','','',''),(3,'p','10308038636343296','/checks/reject','POST','','','');
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
  `content` text NOT NULL,
  `title` varchar(255) NOT NULL,
  `author_id` bigint DEFAULT NULL,
  `status` varchar(20) NOT NULL DEFAULT 'Pending',
  `remark` text,
  `created_at` bigint NOT NULL,
  `updated_at` bigint NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_checks_author` (`author_id`),
  KEY `idx_checks_updated_at` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
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
-- Table structure for table `posts`
--

DROP TABLE IF EXISTS `posts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `posts` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `content` text NOT NULL,
  `created_at` bigint NOT NULL,
  `updated_at` bigint NOT NULL,
  `deleted_at` bigint DEFAULT NULL,
  `deleted` tinyint(1) DEFAULT '0',
  `status` varchar(20) DEFAULT 'draft',
  `author_id` bigint DEFAULT NULL,
  `slug` varchar(100) DEFAULT NULL,
  `category_id` bigint DEFAULT NULL,
  `tags` varchar(255) DEFAULT '',
  `comment_count` bigint DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_posts_slug` (`slug`),
  KEY `idx_posts_category_id` (`category_id`),
  KEY `idx_posts_updated_time` (`updated_at`),
  KEY `idx_posts_deleted_time` (`deleted_at`),
  KEY `idx_posts_author` (`author_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `posts`
--

LOCK TABLES `posts` WRITE;
/*!40000 ALTER TABLE `posts` DISABLE KEYS */;
/*!40000 ALTER TABLE `posts` ENABLE KEYS */;
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
  `biz_name` longtext,
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
  `biz_name` longtext,
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
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_email` (`email`),
  UNIQUE KEY `idx_users_phone` (`phone`),
  KEY `idx_users_deleted_time` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=10308038636343297 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (10308038636343296,1718025027925,1718025027925,0,'','$2a$10$LKLuaqEF38fA..XHCfRI.OtCxo58vYub/Fq.40X4CaaV.IqemrUd2',NULL,'admin',NULL,'');
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

-- Dump completed on 2024-06-10 21:15:52
