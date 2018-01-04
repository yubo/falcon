-- USE falcon;
SET NAMES utf8;
SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
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
  `hostname`		varchar(255)	DEFAULT ''	NOT NULL,
  `version`		varchar(16)	DEFAULT ''	NOT NULL,
  `plugin_version`	varchar(32)	DEFAULT ''	NOT NULL,
  `meta`		varchar(512)	DEFAULT ''	NOT NULL,
  PRIMARY KEY (`hostname`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci
COMMENT='falcon agent 统计';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `plugin_dir`
--

DROP TABLE IF EXISTS `plugin_dir`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `plugin_dir` (
  `id`			bigint unsigned			NOT NULL AUTO_INCREMENT,
  `grp_id`		bigint unsigned DEFAULT '0'	NOT NULL,
  `dir`			varchar(255)	DEFAULT ''	NOT NULL,
  `create_user`		varchar(64)	DEFAULT ''	NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_plugin_dir_grp_id` (`grp_id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci
COMMENT='配置 agent 插件';
/*!40101 SET character_set_client = @saved_cs_client */;




DROP TABLE IF EXISTS `dashboard_graph`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `dashboard_graph` (
  `id`			bigint unsigned			NOT NULL AUTO_INCREMENT,
  `title`		char(128)	DEFAULT ''	NOT NULL,
  `hosts`		varchar(10240)	DEFAULT ''	NOT NULL,
  `counters`		varchar(1024)	DEFAULT ''	NOT NULL,
  `screen_id`		bigint unsigned	DEFAULT '0'	NOT NULL,
  `timespan`		integer		DEFAULT '3600'	NOT NULL,
  `graph_type`		char(2)		DEFAULT 'h'	NOT NULL,
  `method`		char(8)		DEFAULT ''	NOT NULL,
  `position`		bigint unsigned	DEFAULT '0'	NOT NULL,
  `falcon_tags`		varchar(512)	DEFAULT ''	NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_sid` (`screen_id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci
COMMENT='图片配置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `dashboard_screen`
--

DROP TABLE IF EXISTS `dashboard_screen`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `dashboard_screen` (
  `id`			bigint unsigned			NOT NULL AUTO_INCREMENT,
  `pid`			bigint unsigned	DEFAULT '0'	NOT NULL,
  `name`		char(128)	DEFAULT ''	NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_pid` (`pid`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci
COMMENT='视图配置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tmp_graph`
--

DROP TABLE IF EXISTS `tmp_graph`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tmp_graph` (
  `id`			bigint unsigned			NOT NULL AUTO_INCREMENT,
  `endpoints`		varchar(10240)	DEFAULT ''	NOT NULL,
  `counters`		varchar(10240)	DEFAULT ''	NOT NULL,
  `ck`			varchar(32)	DEFAULT ''	NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_ck` (`ck`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci
COMMENT='临时图片配置，已废弃';
/*!40101 SET character_set_client = @saved_cs_client */;



-- -----------------------------------------------------
-- Table `kv`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `kv`;
CREATE TABLE `kv` (
  `section`		varchar(128)	DEFAULT ''	NOT NULL,
  `key`			varchar(128)	DEFAULT ''	NOT NULL,
  `value`		BLOB				NOT NULL,
  INDEX `index_section` (`section`),
  INDEX `index_key` (`key`),
  UNIQUE INDEX `index_section_key` (`section`, `key`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci
COMMENT='key value 存储';


-- -----------------------------------------------------
-- Table `host`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `host`;
CREATE TABLE `host` (
  `id`			bigint unsigned			NOT NULL AUTO_INCREMENT,
  `uuid`		varchar(128)	DEFAULT ''	NOT NULL COMMENT 'external/global id',
  `name`		varchar(128)	DEFAULT ''	NOT NULL,
  `type`		varchar(64)	DEFAULT ''	NOT NULL,
  `status`		varchar(64)	DEFAULT ''	NOT NULL,
  `loc`			varchar(128)	DEFAULT ''	NOT NULL,
  `idc`			varchar(128)	DEFAULT ''	NOT NULL,
  `disabled`		integer		DEFAULT '0'	NOT NULL,
  `pause`		integer		DEFAULT '0'	NOT NULL,
  `maintain_begin`	integer		DEFAULT '0'	NOT NULL,
  `maintain_end`	integer		DEFAULT '0'	NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `index_status` (`status`),
  INDEX `index_type` (`type`),
  UNIQUE INDEX `index_name` (`name`)
) ENGINE = InnoDB AUTO_INCREMENT=1000 DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci
COMMENT='机器表';


-- -----------------------------------------------------
-- Table `token`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `token`;
CREATE TABLE `token` (
  `id`			bigint unsigned			NOT NULL AUTO_INCREMENT,
  `name`		varchar(32)	DEFAULT ''	NOT NULL,
  `cname`		varchar(64)	DEFAULT ''	NOT NULL,
  `note`		varchar(255)	DEFAULT ''	NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_name` (`name`)
) ENGINE = InnoDB AUTO_INCREMENT=1000 DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci COMMENT = '权限点';


-- -----------------------------------------------------
-- Table `role`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role` (
  `id`			bigint unsigned	NOT NULL AUTO_INCREMENT,
  `name`		varchar(32)	DEFAULT ''	NOT NULL,
  `cname`		varchar(64)	DEFAULT ''	NOT NULL,
  `note`		varchar(255)	DEFAULT ''	NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_name` (`name`)
) ENGINE = InnoDB AUTO_INCREMENT=1000 DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci
COMMENT = '角色';

-- -----------------------------------------------------
-- Table `session`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `session`;
CREATE TABLE `session` (
  `session_key`		CHAR(64)			NOT NULL,
  `session_data`	BLOB				NULL,
  `session_expiry`	bigint unsigned			NOT NULL,
  PRIMARY KEY (`session_key`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci
COMMENT = '会话';


-- -----------------------------------------------------
-- Table `tag`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `tag`;
CREATE TABLE `tag` (
  `id`			bigint unsigned			NOT NULL AUTO_INCREMENT,
  `name`		varchar(255)	DEFAULT ''	NOT NULL,
  `type`		integer		DEFAULT '0'	NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_name` (`name`)
) ENGINE = InnoDB AUTO_INCREMENT=1000 DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci
COMMENT = '服务树节点';

DROP TABLE IF EXISTS `tag_rel`;
CREATE TABLE `tag_rel` (
  `id`			bigint unsigned			NOT NULL AUTO_INCREMENT,
  `tag_id`		bigint unsigned	DEFAULT '0'	NOT NULL,
  `sup_tag_id`		bigint unsigned	DEFAULT '0'	NOT NULL COMMENT 'Superior/Self tag id',
  `offset`		integer		DEFAULT '0'	NOT NULL COMMENT 'relation type',
  PRIMARY KEY (`id`),
  INDEX `index_tag_id` (`tag_id`),
  INDEX `index_sup_tag_id` (`sup_tag_id`),
  INDEX `index_offset` (`offset`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci
COMMENT = '服务树节点关系';


-- -----------------------------------------------------
-- Table `tag_host`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `tag_host`;
CREATE TABLE `tag_host` (
  `id`			bigint unsigned			NOT NULL AUTO_INCREMENT,
  `tag_id`		bigint unsigned DEFAULT '0'	NOT NULL,
  `host_id`		bigint unsigned DEFAULT '0'	NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `index_tag_id` (`tag_id`),
  INDEX `index_host_id` (`host_id`),
  UNIQUE INDEX `index_tag_host` (`tag_id`, `host_id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci
COMMENT = '服务树节点与机器的绑定关系';


-- -----------------------------------------------------
-- Table `user`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id`			bigint unsigned			NOT NULL AUTO_INCREMENT,
  `muid`		bigint unsigned	DEFAULT 0	NOT NULL COMMENT 'master uid/bind to',
  `uuid`		varchar(128)	DEFAULT ''	NOT NULL COMMENT 'external/global id',
  `name`		varchar(128)	DEFAULT ''	NOT NULL,
  `cname`		varchar(64)	DEFAULT ''	NOT NULL,
  `email`		varchar(128)	DEFAULT ''	NOT NULL,
  `phone`		varchar(16)	DEFAULT ''	NOT NULL,
  `qq`			varchar(16)	DEFAULT ''	NOT NULL,
  `disabled`		integer		DEFAULT '0'	NOT NULL,
  `avatarurl`		varchar(256)	DEFAULT ''	NOT NULL,
  `extra`		BLOB				NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_name` (`name`)
) ENGINE = InnoDB AUTO_INCREMENT=1000
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_unicode_ci
COMMENT = '用户';

-- -----------------------------------------------------
-- Table `log`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `log`;
CREATE TABLE `log` (
  `id`			bigint unsigned			NOT NULL AUTO_INCREMENT,
  `module`		bigint unsigned	DEFAULT '0'	NOT NULL,
  `module_id`		bigint unsigned	DEFAULT '0'	NOT NULL,
  `user_id`		bigint unsigned	DEFAULT '0'	NOT NULL,
  `action`		bigint unsigned	DEFAULT '0'	NOT NULL,
  `data`		BLOB				NULL,
  `time`		integer		DEFAULT '0'	NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci
COMMENT = 'log';





-- -----------------------------------------------------
-- Table `tpl_rel`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `tpl_rel`;
CREATE TABLE `tpl_rel` (
  `id`			bigint unsigned			NOT NULL AUTO_INCREMENT,
  `tpl_id`		bigint unsigned			NOT NULL,
  `tag_id`		bigint unsigned			NOT NULL,
  `sub_id`		bigint unsigned			NOT NULL,
  `creator`		bigint unsigned			NOT NULL,
  `type_id`		bigint unsigned	DEFAULT '0'	NOT NULL COMMENT 'relation type',
  PRIMARY KEY (`id`),
  INDEX `index_tpl_id` (`tpl_id`),
  INDEX `index_tag_id` (`tag_id`),
  INDEX `index_sub_id` (`sub_id`),
  INDEX `index_creator` (`creator`),
  INDEX `index_type_id` (`type_id`),
  UNIQUE INDEX `index_tpl_rel` (`tpl_id`, `tag_id`, `sub_id`, `type_id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci COMMENT = '节点上的模板关联(tag,tpl,sub_meta)';

--
-- Table structure for table `action`
--

DROP TABLE IF EXISTS `action`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `action` (
  `id`			bigint unsigned			NOT NULL AUTO_INCREMENT,
  `uic`			varchar(255)	DEFAULT ''	NOT NULL ,
  `url`			varchar(255)	DEFAULT ''	NOT NULL ,
  `send_sms`		integer		DEFAULT '0'	NOT NULL ,
  `send_mail`		integer		DEFAULT '0'	NOT NULL ,
  `callback`		integer		DEFAULT '0'	NOT NULL ,
  `before_callback_sms`	integer		DEFAULT '0'	NOT NULL ,
  `before_callback_mail` integer	DEFAULT '0'	NOT NULL ,
  `after_callback_sms`	integer		DEFAULT '0'	NOT NULL ,
  `after_callback_mail`	integer		DEFAULT '0'	NOT NULL ,
  PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_unicode_ci
COMMENT='事件行为';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `event_trigger`
--

DROP TABLE IF EXISTS `event_trigger`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `event_trigger` (
  `id`			bigint unsigned			NOT NULL AUTO_INCREMENT,
  `parent_id`		bigint unsigned DEFAULT '0'	NOT NULL COMMENT '父触发器 id(group)',
  `tpl_id`		bigint unsigned	DEFAULT '0'	NOT NULL COMMENT '模板 id',
  `tag_id`		bigint unsigned	DEFAULT '0'	NOT NULL COMMENT '绑定的 tag id',
  `refcnt`		integer		DEFAULT '0'	NOT NULL COMMENT 'reference number',
  `version`		integer		DEFAULT '0'	NOT NULL,
  `name`		varchar(128)			NULL,
  `metric`		varchar(128)	DEFAULT ''	NOT NULL,
  `tags`		varchar(512)	DEFAULT ''	NOT NULL,
  `priority`		integer		DEFAULT '0'	NOT NULL,
  `expr`		varchar(512)	DEFAULT ''	NOT NULL COMMENT 'expr=func+op+value',
  `msg`			varchar(512)	DEFAULT ''	NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `index_parent_id` (`parent_id`),
  INDEX `index_tpl_id` (`tpl_id`),
  INDEX `index_tag_id` (`tag_id`),
  UNIQUE INDEX `index_event_trigger_tag_name` (`tag_id`, `name`)
) ENGINE = InnoDB AUTO_INCREMENT=10000 DEFAULT CHARACTER SET=utf8 COLLATE=utf8_unicode_ci
COMMENT='事件策略';
/*!40101 SET character_set_client = @saved_cs_client */;


SET SQL_MODE=@OLD_SQL_MODE;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;

INSERT INTO `tag` (`id`, `name`) VALUES (1, '');
INSERT INTO `tag_rel` (`tag_id`, `sup_tag_id`) VALUES (1, 1);
INSERT INTO `user` (`id`, `uuid`, `name`, `cname`, `email`, `phone`, `qq`)
VALUES (1, 'root@localhost', 'system', 'system', 'root@localhost', '', '');

INSERT INTO `token` (`id`, `name`, `cname`, `note`) VALUES
    (1, 'falcon_read', 'read', 'read'),
    (2, 'falcon_operate', 'operate', 'operate'),
    (3, 'falcon_admin', 'admin', 'admin');
