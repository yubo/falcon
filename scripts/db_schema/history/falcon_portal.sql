-- MySQL dump 10.13  Distrib 5.1.66, for redhat-linux-gnu (x86_64)
--
-- Database: falcon_portal
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
-- Table structure for table `action2`
--

DROP TABLE IF EXISTS `action2`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `action2` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uic` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `url` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `send_sms` tinyint(4) NOT NULL DEFAULT '0',
  `send_mail` tinyint(4) NOT NULL DEFAULT '0',
  `callback` tinyint(4) NOT NULL DEFAULT '0',
  `before_callback_sms` tinyint(4) NOT NULL DEFAULT '0',
  `before_callback_mail` tinyint(4) NOT NULL DEFAULT '0',
  `after_callback_sms` tinyint(4) NOT NULL DEFAULT '0',
  `after_callback_mail` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4705 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `agents_info`
--

DROP TABLE IF EXISTS `agents_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `agents_info` (
  `hostname` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `version` varchar(16) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `plugin_version` varchar(32) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `last_update` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `meta` varchar(512) COLLATE utf8_unicode_ci DEFAULT '',
  PRIMARY KEY (`hostname`),
  KEY `idx_time` (`last_update`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `audit`
--

DROP TABLE IF EXISTS `audit`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `audit` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `op_source` varchar(15) NOT NULL COMMENT 'template/strategy/cluster...',
  `op_type` varchar(10) NOT NULL COMMENT 'delete/insert/update',
  `old` varchar(4096) NOT NULL COMMENT 'old info',
  `new` varchar(4096) NOT NULL COMMENT 'new info',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `user` varchar(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=19468 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `cluster`
--

DROP TABLE IF EXISTS `cluster`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cluster` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `node` varchar(255) NOT NULL COMMENT 'xbox tag string',
  `numerator` varchar(9240) NOT NULL,
  `denominator` varchar(9240) NOT NULL,
  `endpoint` varchar(255) NOT NULL,
  `metric` varchar(255) NOT NULL,
  `tags` varchar(255) NOT NULL,
  `ds_type` varchar(32) NOT NULL,
  `step` int(11) NOT NULL,
  `last_update` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `creator` varchar(64) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_cluster_endpoint_metric_tags` (`endpoint`,`metric`,`tags`),
  KEY `idx_cluster_node` (`node`)
) ENGINE=InnoDB AUTO_INCREMENT=7966 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `expression`
--

DROP TABLE IF EXISTS `expression`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `expression` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `expression` varchar(1024) COLLATE utf8_unicode_ci NOT NULL,
  `op` varchar(8) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `condition` varchar(16) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `max_step` int(11) NOT NULL DEFAULT '1',
  `priority` tinyint(4) NOT NULL DEFAULT '0',
  `msg` varchar(1024) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `action_threshold` varchar(16) COLLATE utf8_unicode_ci NOT NULL DEFAULT 'last(#1)',
  `action_id` int(10) unsigned NOT NULL DEFAULT '0',
  `create_user` varchar(64) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `pause` tinyint(1) NOT NULL DEFAULT '0',
  `expr_name` varchar(128) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_expression_expr_name` (`expr_name`)
) ENGINE=InnoDB AUTO_INCREMENT=2125 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `grp`
--

DROP TABLE IF EXISTS `grp`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `grp` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `grp_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `create_user` varchar(64) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `come_from` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_host_grp_grp_name` (`grp_name`)
) ENGINE=InnoDB AUTO_INCREMENT=2910 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `grp_host`
--

DROP TABLE IF EXISTS `grp_host`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `grp_host` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `grp_id` int(10) unsigned NOT NULL,
  `host_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_grp_host_grp_id` (`grp_id`),
  KEY `idx_grp_host_host_id` (`host_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1179307 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `grp_tpl`
--

DROP TABLE IF EXISTS `grp_tpl`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `grp_tpl` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `grp_id` int(10) unsigned NOT NULL,
  `tpl_id` int(10) unsigned NOT NULL,
  `bind_user` varchar(64) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `idx_grp_tpl_grp_id` (`grp_id`),
  KEY `idx_grp_tpl_tpl_id` (`tpl_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4384 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `host`
--

DROP TABLE IF EXISTS `host`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `host` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `hostname` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `pause` tinyint(1) NOT NULL DEFAULT '0',
  `uuid` varchar(16) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `maintain_begin` int(11) NOT NULL DEFAULT '0',
  `maintain_end` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_host_hostname` (`hostname`)
) ENGINE=InnoDB AUTO_INCREMENT=198019 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `mockcfg`
--

DROP TABLE IF EXISTS `mockcfg`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `mockcfg` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT 'name of mockcfg, used for uuid',
  `obj` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT 'desc of object',
  `obj_type` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT 'type of object, host or group or other',
  `metric` varchar(128) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `tags` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `dstype` varchar(32) COLLATE utf8_unicode_ci NOT NULL DEFAULT 'GAUGE',
  `step` int(11) unsigned NOT NULL DEFAULT '60',
  `mock` double NOT NULL DEFAULT '0' COMMENT 'mocked value when nodata occurs',
  `creator` varchar(64) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `t_create` datetime NOT NULL COMMENT 'create time',
  `t_modify` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'last modify time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_name` (`name`),
  UNIQUE KEY `uniq_o_m_t` (`obj`,`obj_type`,`metric`,`tags`)
) ENGINE=InnoDB AUTO_INCREMENT=36 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `plugin_dir`
--

DROP TABLE IF EXISTS `plugin_dir`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `plugin_dir` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `grp_id` int(10) unsigned NOT NULL,
  `dir` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT 'relative to git project root',
  `create_user` varchar(64) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_plugin_dir_grp_id` (`grp_id`)
) ENGINE=InnoDB AUTO_INCREMENT=105 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `strategy`
--

DROP TABLE IF EXISTS `strategy`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `strategy` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `metric_id` int(10) unsigned NOT NULL DEFAULT '0',
  `tags` varchar(2048) COLLATE utf8_unicode_ci DEFAULT NULL,
  `max_step` int(11) NOT NULL DEFAULT '1',
  `priority` tinyint(4) NOT NULL DEFAULT '0',
  `func` varchar(16) COLLATE utf8_unicode_ci NOT NULL DEFAULT 'last(#1)',
  `op` varchar(8) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `condition` varchar(64) COLLATE utf8_unicode_ci NOT NULL,
  `note` varchar(128) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `metric` varchar(1024) COLLATE utf8_unicode_ci DEFAULT NULL,
  `run_begin` varchar(16) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `run_end` varchar(16) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `tpl_id` int(10) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_strategy_tpl_id` (`tpl_id`)
) ENGINE=InnoDB AUTO_INCREMENT=14582 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tpl`
--

DROP TABLE IF EXISTS `tpl`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tpl` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `tpl_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `parent_id` int(10) unsigned NOT NULL DEFAULT '0',
  `action_id` int(10) unsigned NOT NULL DEFAULT '0',
  `create_user` varchar(64) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tpl_tpl_name` (`tpl_name`),
  KEY `idx_tpl_create_user` (`create_user`)
) ENGINE=InnoDB AUTO_INCREMENT=2700 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2016-12-02 14:54:43
