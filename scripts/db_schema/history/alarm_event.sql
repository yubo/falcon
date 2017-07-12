-- MySQL dump 10.13  Distrib 5.1.66, for redhat-linux-gnu (x86_64)
--
-- Database: alarm_event
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
-- Table structure for table `current`
--

DROP TABLE IF EXISTS `current`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `current` (
  `event_id` varchar(64) NOT NULL DEFAULT '',
  `status` varchar(32) DEFAULT NULL,
  `endpoint` varchar(255) DEFAULT NULL,
  `current_step` tinyint(4) DEFAULT NULL,
  `event_ts` varchar(32) DEFAULT NULL,
  `metric` varchar(64) DEFAULT NULL,
  `pushed_tags` varchar(255) DEFAULT NULL,
  `left_val` double DEFAULT NULL,
  `func` varchar(32) DEFAULT NULL,
  `op` varchar(4) DEFAULT NULL,
  `right_val` double DEFAULT NULL,
  `max_step` int(11) DEFAULT NULL,
  `priority` int(11) DEFAULT NULL,
  `note` varchar(255) DEFAULT NULL,
  `tpl_id` int(11) DEFAULT NULL,
  `action_id` int(11) DEFAULT NULL,
  `expression_id` int(11) DEFAULT NULL,
  `strategy_id` int(11) DEFAULT NULL,
  `last_values` varchar(1024) DEFAULT NULL,
  `db_ts` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`event_id`),
  KEY `endpoint` (`endpoint`),
  KEY `metric` (`metric`,`pushed_tags`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `event`
--

DROP TABLE IF EXISTS `event`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `event` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `event_id` varchar(64) DEFAULT NULL,
  `status` varchar(32) DEFAULT NULL,
  `endpoint` varchar(255) DEFAULT NULL,
  `hostname` varchar(64) NOT NULL DEFAULT '',
  `is_customize` int(1) NOT NULL DEFAULT '1',
  `pdl` varchar(255) NOT NULL DEFAULT '',
  `idc` varchar(16) NOT NULL DEFAULT '',
  `current_step` tinyint(4) DEFAULT NULL,
  `event_ts` varchar(32) DEFAULT NULL,
  `metric` varchar(64) DEFAULT NULL,
  `pushed_tags` varchar(255) DEFAULT NULL,
  `left_val` double DEFAULT NULL,
  `func` varchar(32) DEFAULT NULL,
  `op` varchar(4) DEFAULT NULL,
  `right_val` double DEFAULT NULL,
  `max_step` int(11) DEFAULT NULL,
  `priority` int(11) DEFAULT NULL,
  `note` varchar(255) DEFAULT NULL,
  `tpl_id` int(11) DEFAULT NULL,
  `action_id` int(11) DEFAULT NULL,
  `expression_id` int(11) DEFAULT NULL,
  `strategy_id` int(11) DEFAULT NULL,
  `last_values` varchar(1024) DEFAULT NULL,
  `db_ts` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_ts` (`db_ts`),
  KEY `idx_hostname` (`hostname`),
  KEY `idx_pdl` (`pdl`),
  KEY `idx_is_customize` (`is_customize`),
  KEY `idx_idc` (`idc`)
) ENGINE=InnoDB AUTO_INCREMENT=13818641 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `pdl_event_count`
--

DROP TABLE IF EXISTS `pdl_event_count`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pdl_event_count` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `pdl` char(96) NOT NULL DEFAULT '',
  `count` int(11) NOT NULL DEFAULT '0',
  `ts` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_pdl` (`pdl`)
) ENGINE=InnoDB AUTO_INCREMENT=1650049 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2016-12-12 15:04:20
