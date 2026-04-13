-- MySQL dump 10.13  Distrib 8.4.8, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: hmimg_db
-- ------------------------------------------------------
-- Server version	8.4.8

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
-- Table structure for table `hmimg_albums`
--

DROP TABLE IF EXISTS `hmimg_albums`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hmimg_albums` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci,
  `created_by` int unsigned DEFAULT NULL,
  `cover_image` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_hmimg_albums_created_by` (`created_by`),
  CONSTRAINT `hmimg_albums_created_by_foreign` FOREIGN KEY (`created_by`) REFERENCES `hmimg_users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `hmimg_albums`
--

LOCK TABLES `hmimg_albums` WRITE;
/*!40000 ALTER TABLE `hmimg_albums` DISABLE KEYS */;

/*!40000 ALTER TABLE `hmimg_albums` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `hmimg_collection_items`
--

DROP TABLE IF EXISTS `hmimg_collection_items`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hmimg_collection_items` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `collection_id` int unsigned NOT NULL,
  `item_type` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `item_id` int unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `hmimg_collection_items_collection_id_item_type_item_id_unique` (`collection_id`,`item_type`,`item_id`),
  UNIQUE KEY `uq_collection_item` (`collection_id`,`item_type`,`item_id`),
  KEY `idx_hmimg_collection_items_collection_id` (`collection_id`),
  CONSTRAINT `hmimg_collection_items_collection_id_foreign` FOREIGN KEY (`collection_id`) REFERENCES `hmimg_collections` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `hmimg_collection_items`
--

LOCK TABLES `hmimg_collection_items` WRITE;
/*!40000 ALTER TABLE `hmimg_collection_items` DISABLE KEYS */;

/*!40000 ALTER TABLE `hmimg_collection_items` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `hmimg_collections`
--

DROP TABLE IF EXISTS `hmimg_collections`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hmimg_collections` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci,
  `created_by` int unsigned DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_hmimg_collections_created_by` (`created_by`),
  CONSTRAINT `hmimg_collections_created_by_foreign` FOREIGN KEY (`created_by`) REFERENCES `hmimg_users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `hmimg_collections`
--

LOCK TABLES `hmimg_collections` WRITE;
/*!40000 ALTER TABLE `hmimg_collections` DISABLE KEYS */;

/*!40000 ALTER TABLE `hmimg_collections` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `hmimg_images`
--

DROP TABLE IF EXISTS `hmimg_images`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hmimg_images` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `filename` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `original_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `path` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `size` bigint NOT NULL,
  `mimetype` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `album_id` int unsigned NOT NULL,
  `uploaded_by` int unsigned DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_hmimg_images_album_id` (`album_id`),
  KEY `idx_hmimg_images_uploaded_by` (`uploaded_by`),
  CONSTRAINT `hmimg_images_album_id_foreign` FOREIGN KEY (`album_id`) REFERENCES `hmimg_albums` (`id`) ON DELETE CASCADE,
  CONSTRAINT `hmimg_images_uploaded_by_foreign` FOREIGN KEY (`uploaded_by`) REFERENCES `hmimg_users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `hmimg_images`
--

LOCK TABLES `hmimg_images` WRITE;
/*!40000 ALTER TABLE `hmimg_images` DISABLE KEYS */;

/*!40000 ALTER TABLE `hmimg_images` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `hmimg_settings`
--

DROP TABLE IF EXISTS `hmimg_settings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hmimg_settings` (
  `key` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `value` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `hmimg_settings`
--

LOCK TABLES `hmimg_settings` WRITE;
/*!40000 ALTER TABLE `hmimg_settings` DISABLE KEYS */;
INSERT INTO `hmimg_settings` VALUES ('allow_registration','false','2026-01-27 22:06:25.000','2026-01-28 22:34:33.436'),('default_language','zh','2026-01-28 21:36:45.818','2026-01-28 22:34:33.451'),('max_users','3','2026-01-27 22:06:25.000','2026-01-28 22:34:33.428'),('website_title','Mica Images','2026-01-28 12:18:21.000','2026-01-28 22:34:33.443');
/*!40000 ALTER TABLE `hmimg_settings` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `hmimg_users`
--

DROP TABLE IF EXISTS `hmimg_users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hmimg_users` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `email` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `phone` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `role` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT 'user',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_hmimg_users_username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `hmimg_users`
--

LOCK TABLES `hmimg_users` WRITE;
/*!40000 ALTER TABLE `hmimg_users` DISABLE KEYS */;
INSERT INTO `hmimg_users` VALUES (1,'admin','$2b$10$DFj1EyzpBgXe3h9tgTHTY.KrXYnr4eBchYRqdZldd9Qe3jnGPJ3uu','admin@yourdomaname.com','+8613200000000','admin',NULL,NULL);
/*!40000 ALTER TABLE `hmimg_users` ENABLE KEYS */;
UNLOCK TABLES;
