-- USE falcon;
SET NAMES utf8;
SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

-- -----------------------------------------------------
-- Schema falcon
-- -----------------------------------------------------
-- CREATE SCHEMA IF NOT EXISTS `falcon` DEFAULT CHARACTER SET utf8 ;
-- USE `falcon` ;


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
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci;
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
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;



/**
 * nodata mock config
 */
DROP TABLE IF EXISTS `mockcfg`;
CREATE TABLE `mockcfg` (
  `id`       BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name`     VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'name of mockcfg, used for uuid',
  `obj`      VARCHAR(10240) NOT NULL DEFAULT '' COMMENT 'desc of object',
  `obj_type` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'type of object, host or group or other',
  `metric`   VARCHAR(128) NOT NULL DEFAULT '',
  `tags`     VARCHAR(1024) NOT NULL DEFAULT '',
  `dstype`   VARCHAR(32)  NOT NULL DEFAULT 'GAUGE',
  `step`     INT(11) UNSIGNED  NOT NULL DEFAULT 60,
  `mock`     DOUBLE  NOT NULL DEFAULT 0  COMMENT 'mocked value when nodata occurs',
  `creator`  VARCHAR(64)  NOT NULL DEFAULT '',
  `t_create` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_name` (`name`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci;

/**
 *  aggregator aggreator metric config table
 */
DROP TABLE IF EXISTS `aggreator`;
CREATE TABLE `aggreator` (
  `id`          INT UNSIGNED   NOT NULL AUTO_INCREMENT,
  `tag_id`      INT            NOT NULL,
  `numerator`   VARCHAR(10240) NOT NULL,
  `denominator` VARCHAR(10240) NOT NULL,
  `endpoint`    VARCHAR(255)   NOT NULL,
  `metric`      VARCHAR(255)   NOT NULL,
  `tags`        VARCHAR(255)   NOT NULL,
  `ds_type`     VARCHAR(255)   NOT NULL,
  `step`        INT            NOT NULL,
  `last_update` TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `creator`     VARCHAR(255)   NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci;


DROP TABLE IF EXISTS `dashboard_graph`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `dashboard_graph` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `title` char(128) NOT NULL,
  `hosts` varchar(10240) NOT NULL DEFAULT '',
  `counters` varchar(1024) NOT NULL DEFAULT '',
  `screen_id` int(11) unsigned NOT NULL,
  `timespan` int(11) unsigned NOT NULL DEFAULT '3600',
  `graph_type` char(2) NOT NULL DEFAULT 'h',
  `method` char(8) DEFAULT '',
  `position` int(11) unsigned NOT NULL DEFAULT '0',
  `falcon_tags` varchar(512) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `idx_sid` (`screen_id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `dashboard_screen`
--

DROP TABLE IF EXISTS `dashboard_screen`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `dashboard_screen` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `pid` int(11) unsigned NOT NULL DEFAULT '0',
  `name` char(128) NOT NULL,
  `time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_pid` (`pid`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tmp_graph`
--

DROP TABLE IF EXISTS `tmp_graph`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tmp_graph` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `endpoints` varchar(10240) NOT NULL DEFAULT '',
  `counters` varchar(10240) NOT NULL DEFAULT '',
  `ck` varchar(32) NOT NULL,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_ck` (`ck`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;



-- -----------------------------------------------------
-- Table `kv`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `kv`;
CREATE TABLE `kv` (
  `section` VARCHAR(128) NOT NULL,
  `key` VARCHAR(128) NOT NULL,
  `value` BLOB NOT NULL,
  INDEX `index_section` (`section`),
  INDEX `index_key` (`key`),
  UNIQUE INDEX `index_section_key` (`section`, `key`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci;


-- -----------------------------------------------------
-- Table `host`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `host`;
CREATE TABLE `host` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `uuid` VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'external/global id',
  `name` VARCHAR(128) NOT NULL DEFAULT '',
  `type` VARCHAR(64) NOT NULL DEFAULT '',
  `status` VARCHAR(64) NOT NULL DEFAULT '',
  `loc` VARCHAR(128) NOT NULL DEFAULT '',
  `idc` VARCHAR(128) NOT NULL DEFAULT '',
  `disabled` TINYINT(4) NOT NULL DEFAULT '0',
  `pause` TINYINT(4) NOT NULL DEFAULT '0',
  `maintain_begin` INT(11) NOT NULL DEFAULT '0',
  `maintain_end` INT(11) NOT NULL DEFAULT '0',
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `index_status` (`status`),
  INDEX `index_type` (`type`),
  UNIQUE INDEX `index_name` (`name`)
) ENGINE = InnoDB AUTO_INCREMENT=1000 DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci;


-- -----------------------------------------------------
-- Table `token`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `token`;
CREATE TABLE `token` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(32) NOT NULL,
  `cname` VARCHAR(64) NOT NULL,
  `note` VARCHAR(255) NOT NULL DEFAULT '',
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_name` (`name`)
) ENGINE = InnoDB AUTO_INCREMENT=100 DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci COMMENT = '权限点';


-- -----------------------------------------------------
-- Table `role`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(32) NOT NULL,
  `cname` VARCHAR(64) NOT NULL,
  `note` VARCHAR(255) NOT NULL DEFAULT '',
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_name` (`name`)
) ENGINE = InnoDB AUTO_INCREMENT=1000 DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci COMMENT = '角色';

-- -----------------------------------------------------
-- Table `session`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `session`;
CREATE TABLE `session` (
  `session_key` CHAR(64) NOT NULL,
  `session_data` BLOB NULL DEFAULT NULL,
  `session_expiry` INT(11) UNSIGNED NOT NULL,
  PRIMARY KEY (`session_key`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci;


-- -----------------------------------------------------
-- Table `tag`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `tag`;
CREATE TABLE `tag` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL DEFAULT '',
  `type` TINYINT(4) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'tag type',
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_name` (`name`)
) ENGINE = InnoDB AUTO_INCREMENT=1000 DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci;
-- ALTER TABLE `tag`   
-- ADD COLUMN `type` TINYINT(4) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'tag type' AFTER `name`;

DROP TABLE IF EXISTS `tag_rel`;
CREATE TABLE `tag_rel` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `tag_id` INT(11) UNSIGNED NOT NULL DEFAULT 0,
  `sup_tag_id` INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Superior/Self tag id',
  `offset` TINYINT(4) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'relation type',
  PRIMARY KEY (`id`),
  INDEX `index_tag_id` (`tag_id`),
  INDEX `index_sup_tag_id` (`sup_tag_id`),
  INDEX `index_offset` (`offset`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci;


-- -----------------------------------------------------
-- Table `tag_host`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `tag_host`;
CREATE TABLE `tag_host` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `tag_id` INT(11) UNSIGNED NOT NULL DEFAULT 0,
  `host_id` INT(11) UNSIGNED NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `index_tag_id` (`tag_id`),
  INDEX `index_host_id` (`host_id`),
  UNIQUE INDEX `index_tag_host` (`tag_id`, `host_id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci;


-- -----------------------------------------------------
-- Table `user`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `muid` INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'master uid/bind to',
  `uuid` VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'external/global id',
  `name` VARCHAR(128) NOT NULL,
  `cname` VARCHAR(64) NOT NULL DEFAULT '',
  `email` VARCHAR(128) NOT NULL DEFAULT '',
  `phone` VARCHAR(16) NOT NULL DEFAULT '',
  `qq` VARCHAR(16) NOT NULL DEFAULT '',
  `disabled` TINYINT(4) NOT NULL DEFAULT '0',
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `avatarurl` VARCHAR(256) NOT NULL DEFAULT '',
  `extra` BLOB NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_name` (`name`)
) ENGINE = InnoDB AUTO_INCREMENT=1000 DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci;

-- -----------------------------------------------------
-- Table `team`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `team`;
CREATE TABLE `team` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(32) NOT NULL,
  `note` VARCHAR(255) NOT NULL DEFAULT '',
  `creator` INT(11) UNSIGNED NOT NULL,
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_name` (`name`)
) ENGINE = InnoDB AUTO_INCREMENT=1000 DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci;

-- -----------------------------------------------------
-- Table `team_user`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `team_user`;
CREATE TABLE `team_user` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `team_id` INT(11) UNSIGNED NOT NULL,
  `user_id` INT(11) UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `index_team_id` (`team_id`),
  INDEX `index_user_id` (`user_id`),
  UNIQUE INDEX `index_team_user` (`team_id`, `user_id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci;

-- -----------------------------------------------------
-- Table `log`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `log`;
CREATE TABLE `log` (
  `id`        INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `module`    TINYINT(4) UNSIGNED NOT NULL DEFAULT 0,
  `module_id` INT(11) UNSIGNED NOT NULL DEFAULT 0,
  `user_id`   INT(11) UNSIGNED NOT NULL DEFAULT 0,
  `action`    TINYINT(4) UNSIGNED NOT NULL DEFAULT 0,
  `data`      BLOB NULL DEFAULT NULL,
  `time`      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci;

-- -----------------------------------------------------
-- Table `tpl_rel`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `tpl_rel`;
CREATE TABLE `tpl_rel` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `tpl_id` INT(11) UNSIGNED NOT NULL,
  `tag_id` INT(11) UNSIGNED NOT NULL,
  `sub_id` INT(11) UNSIGNED NOT NULL,
  `creator` INT(11) UNSIGNED NOT NULL,
  `type_id` TINYINT(4) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'relation type',
  PRIMARY KEY (`id`),
  INDEX `index_tpl_id` (`tpl_id`),
  INDEX `index_tag_id` (`tag_id`),
  INDEX `index_sub_id` (`sub_id`),
  INDEX `index_creator` (`creator`),
  INDEX `index_type_id` (`type_id`),
  UNIQUE INDEX `index_tpl_rel` (`tpl_id`, `tag_id`, `sub_id`, `type_id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci COMMENT = '节点上的模板关联(tag,tpl,sub_meta)';

-- -----------------------------------------------------
-- Table `tag_tpl`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `tag_tpl`;
CREATE TABLE `tag_tpl` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `tpl_id` INT(11) UNSIGNED NOT NULL,
  `tag_id` INT(11) UNSIGNED NOT NULL,
  `creator` INT(11) UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `index_tpl_id` (`tpl_id`),
  INDEX `index_tag_id` (`tag_id`),
  INDEX `index_creator` (`creator`),
  UNIQUE INDEX `index_tag_tpl` (`tpl_id`, `tag_id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci COMMENT = '节点上的策略模板';

--
-- Table structure for table `action`
--
DROP TABLE IF EXISTS `action`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `action` (
  `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `sendto` VARCHAR(255) NOT NULL DEFAULT '',
  `url` VARCHAR(255) NOT NULL DEFAULT '',
  `send_flag` INT(10) UNSIGNED NOT NULL DEFAULT '0',
  `cb_falg` INT(10) UNSIGNED NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `expression`
--
DROP TABLE IF EXISTS `expression`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `expression` (
  `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) DEFAULT NULL,
  `expression` VARCHAR(1024) NOT NULL,
  `op` VARCHAR(8) NOT NULL DEFAULT '',
  `condition` VARCHAR(16) NOT NULL DEFAULT '',
  `max_step` INT(11) NOT NULL DEFAULT '1',
  `priority` TINYINT(4) NOT NULL DEFAULT '0',
  `msg` VARCHAR(1024) NOT NULL DEFAULT '',
  `action_threshold` VARCHAR(16) NOT NULL DEFAULT 'last(#1)',
  `action_id` INT(10) UNSIGNED NOT NULL DEFAULT '0',
  `create_user_id` INT(10) UNSIGNED NOT NULL DEFAULT '0',
  `pause` TINYINT(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_expression_name` (`name`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `action`
--

DROP TABLE IF EXISTS `action`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `action` (
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
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci;
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
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `template`
--

DROP TABLE IF EXISTS `template`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `template` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '',
  `parent_id` int(10) unsigned NOT NULL DEFAULT '0',
  `action_id` int(10) unsigned NOT NULL DEFAULT '0',
  `create_user_id` int(10) unsigned NOT NULL DEFAULT '0',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`),
  KEY `idx_tpl_create_user` (`create_user_id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;

INSERT INTO `tag` (`id`, `name`) VALUES (1, '');
INSERT INTO `tag_rel` (`tag_id`, `sup_tag_id`) VALUES (1, 1);
INSERT INTO `user` (`id`, `uuid`, `name`, `cname`, `email`, `phone`, `im`, `qq`)
VALUES (1, 'root@localhost', 'system', 'system', 'root@localhost', '', '', '');

INSERT INTO `token` (`id`, `name`, `cname`, `note`) VALUES
    (1, 'falcon_read', 'read', 'read'),
    (2, 'falcon_operate', 'operate', 'operate'),
    (3, 'falcon_admin', 'admin', 'admin');
