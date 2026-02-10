-- MySQL dump 10.13  Distrib 9.6.0, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: flushproject_logs
-- ------------------------------------------------------
-- Server version	9.6.0

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
SET @MYSQLDUMP_TEMP_LOG_BIN = @@SESSION.SQL_LOG_BIN;
SET @@SESSION.SQL_LOG_BIN= 0;

--
-- Table structure for table `admin_common`
--

DROP TABLE IF EXISTS `admin_common`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admin_common` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'Primary Key',
  `date` datetime NOT NULL,
  `admin_name` varchar(255) NOT NULL DEFAULT 'None',
  `admin_id` int NOT NULL DEFAULT '0',
  `user_name` varchar(255) DEFAULT NULL,
  `user_id` int DEFAULT NULL,
  `action` varchar(255) NOT NULL,
  `additional_information` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `admin_common`
--

/*!40000 ALTER TABLE `admin_common` DISABLE KEYS */;
INSERT INTO `admin_common` VALUES (1,'2026-02-10 03:57:59','',0,NULL,NULL,'searched all users',NULL);
/*!40000 ALTER TABLE `admin_common` ENABLE KEYS */;

--
-- Table structure for table `admin_punishments`
--

DROP TABLE IF EXISTS `admin_punishments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admin_punishments` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'Primary Key',
  `date` datetime NOT NULL,
  `admin_name` varchar(255) NOT NULL DEFAULT 'None',
  `admin_id` int NOT NULL DEFAULT '0',
  `user_name` varchar(255) DEFAULT NULL,
  `user_id` int DEFAULT NULL,
  `action` varchar(255) NOT NULL,
  `additional_information` varchar(255) DEFAULT NULL,
  `issued_until` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `admin_punishments`
--

/*!40000 ALTER TABLE `admin_punishments` DISABLE KEYS */;
/*!40000 ALTER TABLE `admin_punishments` ENABLE KEYS */;

--
-- Dumping routines for database 'flushproject_logs'
--
SET @@SESSION.SQL_LOG_BIN = @MYSQLDUMP_TEMP_LOG_BIN;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-02-10  4:36:36
