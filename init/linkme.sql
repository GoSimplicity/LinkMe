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
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `casbin_rule`
--

LOCK TABLES `casbin_rule` WRITE;
/*!40000 ALTER TABLE `casbin_rule` DISABLE KEYS */;
INSERT INTO `casbin_rule` VALUES (25,'p','10308038636343296','/api/activity/recent','GET','','',''),(2,'p','10308038636343296','/api/checks/approve','POST','','',''),(4,'p','10308038636343296','/api/checks/detail','GET','','',''),(1,'p','10308038636343296','/api/checks/list','GET','','',''),(3,'p','10308038636343296','/api/checks/reject','POST','','',''),(24,'p','10308038636343296','/api/checks/stats','GET','','',''),(5,'p','10308038636343296','/api/permissions/assign','POST','','',''),(13,'p','10308038636343296','/api/permissions/assign_role','POST','','',''),(6,'p','10308038636343296','/api/permissions/list','GET','','',''),(7,'p','10308038636343296','/api/permissions/remove','DELETE','','',''),(16,'p','10308038636343296','/api/permissions/remove_role','DELETE','','',''),(8,'p','10308038636343296','/api/plate/create','POST','','',''),(10,'p','10308038636343296','/api/plate/delete','DELETE','','',''),(9,'p','10308038636343296','/api/plate/list','GET','','',''),(21,'p','10308038636343296','/api/posts/detail_post/:postId','GET','','',''),(19,'p','10308038636343296','/api/posts/list_post','POST','','',''),(23,'p','10308038636343296','/api/posts/stats','GET','','',''),(11,'p','10308038636343296','/api/users/list','POST','','',''),(22,'p','10308038636343296','/api/users/stats','GET','','','');
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
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `checks`
--

LOCK TABLES `checks` WRITE;
/*!40000 ALTER TABLE `checks` DISABLE KEYS */;
INSERT INTO `checks` VALUES (1,1,'测试内容123','测试标题123',10308038636343296,'Approved','',1719917324148,1719917355),(2,2,'测试内容123','测试标题123',10308038636343296,'Approved','',1719917333741,1719917359),(3,3,'测试内容123','测试标题123',10308038636343296,'Approved','',1719917336331,1719917362),(4,4,'测试内容123','测试标题123',10308038636343296,'Approved','',1719917339910,1719917367),(5,5,'测试内容123','测试标题123',10308038636343296,'Approved','',1719917342213,1719917369),(6,6,'测试内容123','测试标题123',10308038636343296,'Approved','',1719917345003,1719917372),(7,7,'测试内容123','测试标题123',10308038636343296,'Approved','',1719919432168,1719919464),(8,8,'测试内容123','测试标题123',10308038636343296,'Approved','',1719919436217,1719919468),(9,9,'测试内容123','测试标题123',10308038636343296,'Approved','',1719919439007,1719919471),(10,10,'测试内容123','测试标题123',10308038636343296,'Approved','',1719919442178,1719919475),(11,11,'测试内容123','测试标题123',10308038636343296,'Approved','',1719919444471,1719919478),(12,12,'测试内容123','测试标题123',10308038636343296,'Approved','',1719919446993,1719919481),(13,13,'测试内容123','测试标题123',10308038636343296,'Approved','',1719919449239,1719919484),(14,14,'测试内容123','测试标题123',10308038636343296,'Approved','',1719919451482,1719919486),(15,15,'测试内容123','测试标题123',10308038636343296,'Approved','',1719919454600,1719919489),(16,16,'测试内容123','测试标题123',10308038636343296,'Approved','',1719919892218,1719919900),(17,17,'测试内容123','测试标题123',10308038636343296,'Approved','',1719920339457,1719920349),(18,18,'测试内容123','测试标题123',10308038636343296,'Approved','',1720092013962,1720092057),(19,19,'测试内容123','测试标题123',10308038636343296,'Under review','',1720092116596,1720092116596),(20,20,'测试内容123','测试标题123',10308038636343296,'Approved','',1720187591029,1720187598);
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
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `plates`
--

LOCK TABLES `plates` WRITE;
/*!40000 ALTER TABLE `plates` DISABLE KEYS */;
INSERT INTO `plates` VALUES (1,'golang板块','golang学习...',1719917266377,1719917266377,0,0,10308038636343296);
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
  `plate_id` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_posts_slug` (`slug`),
  KEY `idx_posts_category_id` (`category_id`),
  KEY `idx_posts_updated_time` (`updated_at`),
  KEY `idx_posts_deleted_time` (`deleted_at`),
  KEY `idx_posts_author` (`author_id`),
  KEY `idx_posts_plate_id` (`plate_id`),
  CONSTRAINT `fk_plates_posts` FOREIGN KEY (`plate_id`) REFERENCES `plates` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `posts`
--

LOCK TABLES `posts` WRITE;
/*!40000 ALTER TABLE `posts` DISABLE KEYS */;
INSERT INTO `posts` VALUES (1,'测试标题1','测试内容1',1719917275465,1719917356008,0,0,'Published',10308038636343296,'9f6cb49a-76bd-4f54-aae8-af72b63a781a',0,'',0,1),(2,'测试标题2','测试内容2',1719917276431,1719917359132,0,0,'Published',10308038636343296,'d014df61-57f7-4462-97ac-c5d5fee0883e',0,'',0,1),(3,'测试标题3','测试内容3',1719917276849,1719917362018,0,0,'Published',10308038636343296,'81152351-60cb-42aa-b8f3-11b8b3f2e62b',0,'',0,1),(4,'测试标题4','测试内容4',1719917277143,1719917367013,0,0,'Published',10308038636343296,'b14ceef2-78b0-4b5c-836c-25b8a4a10c91',0,'',0,1),(5,'测试标题5','测试内容5',1719917277482,1719917369664,0,0,'Published',10308038636343296,'6061c930-c285-4e6d-9a08-0da3ab2bb969',0,'',0,1),(6,'测试标题6','测试内容6',1719917277780,1719917372418,0,0,'Published',10308038636343296,'949772b8-600f-40f9-b432-f5a0d6d92c25',0,'',0,1),(7,'测试标题7','测试内容7',1719919420735,1719919464655,0,0,'Published',10308038636343296,'ebf10c8b-b9e1-4e5c-9640-a48c4b6d2980',0,'',0,1),(8,'测试标题8','测试内容8',1719919421235,1719919468934,0,0,'Published',10308038636343296,'ce3bc95f-64c6-4f4a-9abf-df8d646c0068',0,'',0,1),(9,'测试标题9','测试内容9',1719919421694,1719919471814,0,0,'Published',10308038636343296,'cfab4f7d-62de-4c24-b3d1-964f2f436a5e',0,'',0,1),(10,'测试标题10','测试内容10',1719919422004,1719919475332,0,0,'Published',10308038636343296,'388df334-2c2b-4d76-9ca7-2d860b8c92bf',0,'',0,1),(11,'测试标题11','测试内容11',1719919422361,1719919478185,0,0,'Published',10308038636343296,'0ba4c399-cb16-4e67-8693-a43542660d14',0,'',0,1),(12,'测试标题12','测试内容12',1719919422694,1719919481367,0,0,'Published',10308038636343296,'330cd427-f0b2-4fd2-af02-ff706a4a2725',0,'',0,1),(13,'测试标题13','测试内容13',1719919423044,1719919484257,0,0,'Published',10308038636343296,'70c17b0f-1f13-44cb-ba3c-c8a62614c2b9',0,'',0,1),(14,'测试标题14','测试内容14',1719919423515,1719919486986,0,0,'Published',10308038636343296,'3cc3e3f2-5355-4cd8-9ac9-153bbecf80d1',0,'',0,1),(15,'测试标题15','测试内容15',1719919424245,1719919489433,0,0,'Published',10308038636343296,'b02374e6-387f-48a7-8f3a-ca2749d6f394',0,'',0,1),(16,'测试标题16','测试内容16',1719919882907,1719919900804,0,0,'Published',10308038636343296,'e28beb93-f604-433f-8a92-da40a72ece7c',0,'',0,1),(17,'测试标题17','测试内容17',1719920331872,1719920357597,0,0,'Published',10308038636343296,'2b5342ce-6c11-416e-a9d3-7efa8149b5df',0,'',0,1),(18,'测试标题123','测试内容123',1720092006447,1720092057457,0,0,'Published',10308038636343296,'3a0bb72f-06c9-400d-b716-31f84015892e',0,'',0,1),(19,'测试标题123','测试内容123',1720092108541,1720092108541,0,0,'Draft',10308038636343296,'5c1a6e1d-1938-4367-960c-d4e873b0750f',0,'',0,1),(20,'测试标题123','测试内容123',1720187583269,1720187598886,0,0,'Published',10308038636343296,'da7007a2-4b20-42cb-9cb8-a015396c91b3',0,'',0,1);
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
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `recent_activities`
--

LOCK TABLES `recent_activities` WRITE;
/*!40000 ALTER TABLE `recent_activities` DISABLE KEYS */;
INSERT INTO `recent_activities` VALUES (1,10308038636343296,'审核通过','1720187598');
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
INSERT INTO `users` VALUES (10308038636343296,1718025027925,1718025027925,0,'Bamboo','$2a$10$LKLuaqEF38fA..XHCfRI.OtCxo58vYub/Fq.40X4CaaV.IqemrUd2','2024-07-03 00:00:00','admin','123456','this is LinkMe',0),(18281565326938132,NULL,1,NULL,'aaa','123',NULL,'adsfasdf','1','asdfasd',0),(18281565326938133,NULL,2,NULL,'bbb','123',NULL,'asdfasdfadsfasdf','2','tgsadgasdf',0),(18281565326938134,NULL,3,NULL,'ccc','123',NULL,'qwereqwr','3','asdfasdfasd',0),(18281565326938135,NULL,4,NULL,'sss','123',NULL,'adsgxzv','4','fasdfa',0),(18281565326938136,NULL,5,NULL,'ddd','123',NULL,'bvxzbadg','5','dsfasd',0),(18281565326938137,NULL,6,NULL,'ggg','123',NULL,'afghreha','6','fasdfasdf',0),(18281565326938138,NULL,7,NULL,'hhh','123',NULL,'bvzxnhh','7','afda',0),(18281565326938139,NULL,8,NULL,'yyy','123',NULL,'gfdshfdgj','8','fas',0),(18281565326938140,NULL,9,NULL,'uuu','123',NULL,'tweygsgf','9','fasdf',0),(18281565326938141,NULL,10,NULL,'iii','123',NULL,'cxvbxczdf','11','asdfasdf',0),(18281565326938142,NULL,11,NULL,'ttt','123',NULL,'retweygdx','22','asdf',0),(18281565326938143,NULL,12,NULL,'rrr','123',NULL,'sdfhdsth','33','asdfa',0),(18281565326938144,NULL,13,NULL,'eee','123',NULL,'erwvzb','44','sdfasdf',0),(18281565326938145,NULL,14,NULL,'www','123',NULL,'sdfhsdfhsre','55','asdfasdfadsf',0);
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

-- Dump completed on 2024-07-06 21:29:21
