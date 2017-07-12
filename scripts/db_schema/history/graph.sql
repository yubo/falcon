-- MySQL dump 10.13  Distrib 5.1.66, for redhat-linux-gnu (x86_64)
--
-- Database: graph
-- ------------------------------------------------------
-- Server version	5.5.31-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `counter_num`
--

DROP TABLE IF EXISTS `counter_num`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `counter_num` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `tag` varchar(1024) NOT NULL DEFAULT '',
  `count` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6030 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `endpoint`
--

DROP TABLE IF EXISTS `endpoint`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `endpoint` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `endpoint` varchar(255) COLLATE latin1_bin NOT NULL DEFAULT '',
  `hostname` varchar(255) COLLATE latin1_bin NOT NULL DEFAULT '',
  `ts` int(11) DEFAULT NULL,
  `t_create` datetime NOT NULL COMMENT 'create time',
  `t_modify` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'last modify time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_endpoint` (`endpoint`)
) ENGINE=InnoDB AUTO_INCREMENT=27297833848 DEFAULT CHARSET=latin1 COLLATE=latin1_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `endpoint_counter`
--

DROP TABLE IF EXISTS `endpoint_counter`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `endpoint_counter` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `endpoint_id` bigint(20) unsigned NOT NULL,
  `counter` varchar(255) COLLATE latin1_bin NOT NULL DEFAULT '',
  `step` int(11) NOT NULL DEFAULT '60' COMMENT 'in second',
  `type` varchar(16) COLLATE latin1_bin NOT NULL COMMENT 'GAUGE|COUNTER|DERIVE',
  `ts` int(11) DEFAULT NULL,
  `t_create` datetime NOT NULL COMMENT 'create time',
  `t_modify` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'last modify time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_endpoint_id_counter` (`endpoint_id`,`counter`)
) ENGINE=InnoDB AUTO_INCREMENT=28828725869 DEFAULT CHARSET=latin1 COLLATE=latin1_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `pdl_quota`
--

DROP TABLE IF EXISTS `pdl_quota`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pdl_quota` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `pdl` varchar(96) NOT NULL DEFAULT '',
  `quota` int(11) NOT NULL DEFAULT '500000',
  PRIMARY KEY (`id`),
  KEY `idx_pdl` (`pdl`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tag_endpoint`
--

DROP TABLE IF EXISTS `tag_endpoint`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tag_endpoint` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `tag` varchar(255) COLLATE latin1_bin NOT NULL DEFAULT '' COMMENT 'srv=tv',
  `endpoint_id` bigint(20) unsigned NOT NULL,
  `ts` int(11) DEFAULT NULL,
  `t_create` datetime NOT NULL COMMENT 'create time',
  `t_modify` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'last modify time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tag_endpoint_id` (`tag`,`endpoint_id`)
) ENGINE=InnoDB AUTO_INCREMENT=153925686147 DEFAULT CHARSET=latin1 COLLATE=latin1_bin;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2016-12-02 14:28:14
