-- Copyright 2016 Xiaomi, Inc.
--
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
--     http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

-- -----------------------------------------------------
-- Schema falcon
-- -----------------------------------------------------
-- CREATE SCHEMA IF NOT EXISTS `falcon` DEFAULT CHARACTER SET utf8 ;
-- USE `falcon` ;

-- -----------------------------------------------------
-- Table `kv`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `kv`;
CREATE TABLE `kv` (
  `key` VARCHAR(128) NOT NULL,
  `value` BLOB NOT NULL,
  `type_id` TINYINT(4) NOT NULL,
  PRIMARY KEY (`key`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci;


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
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `index_status` (`status`),
  INDEX `index_type` (`type`),
  UNIQUE INDEX `index_name` (`name`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci COMMENT = '机器';


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
  UNIQUE INDEX `index_name` (`name`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
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
  UNIQUE INDEX `index_name` (`name`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci COMMENT = '角色';

-- -----------------------------------------------------
-- Table `session`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `session`;
CREATE TABLE `session` (
  `session_key` CHAR(64) NOT NULL,
  `session_data` BLOB NULL DEFAULT NULL,
  `session_expiry` INT(11) UNSIGNED NOT NULL,
  PRIMARY KEY (`session_key`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci;


-- -----------------------------------------------------
-- Table `tag`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `tag`;
CREATE TABLE `tag` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL DEFAULT '',
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_name` (`name`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci;

DROP TABLE IF EXISTS `tag_rel`;
CREATE TABLE `tag_rel` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `tag_id` INT(11) UNSIGNED NOT NULL DEFAULT 0,
  `sup_tag_id` INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Superior/Self tag id',
  `type_id` TINYINT(4) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'relation type',
  PRIMARY KEY (`id`),
  INDEX `index_tag_id` (`tag_id`),
  INDEX `index_sup_tag_id` (`sup_tag_id`),
  INDEX `index_type_id` (`type_id`),
  CONSTRAINT `tag_rel_ibfk_1`
    FOREIGN KEY (`tag_id`)
    REFERENCES `tag` (`id`)
    ON DELETE CASCADE,
  CONSTRAINT `tag_rel_ibfk_2`
    FOREIGN KEY (`sup_tag_id`)
    REFERENCES `tag` (`id`)
    ON DELETE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci;


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
  CONSTRAINT `tag_host_rel_ibfk_1`
    FOREIGN KEY (`tag_id`)
    REFERENCES `tag` (`id`)
    ON DELETE CASCADE,
  CONSTRAINT `tag_host_rel_ibfk_2`
    FOREIGN KEY (`host_id`)
    REFERENCES `host` (`id`)
    ON DELETE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci;


-- -----------------------------------------------------
-- Table `user`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `uuid` VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'external/global id',
  `name` VARCHAR(32) NOT NULL,
  `cname` VARCHAR(64) NOT NULL DEFAULT '',
  `email` VARCHAR(128) NOT NULL DEFAULT '',
  `phone` VARCHAR(16) NOT NULL DEFAULT '',
  `im` VARCHAR(32) NOT NULL DEFAULT '',
  `qq` VARCHAR(16) NOT NULL DEFAULT '',
  `disabled` TINYINT(4) NOT NULL DEFAULT '0',
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_name` (`name`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci;


-- -----------------------------------------------------
-- Table `log`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `log`;
CREATE TABLE `log` (
  `id`        INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `module`    VARCHAR(64) NOT NULL DEFAULT '',
  `module_id` INT(11) UNSIGNED NOT NULL DEFAULT 0,
  `user_id`   INT(11) UNSIGNED NOT NULL DEFAULT 0,
  `action`    TINYINT(4) UNSIGNED NOT NULL DEFAULT 0,
  `data`      BLOB NULL DEFAULT NULL,
  `time`      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci;

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
  INDEX `index_type_id` (`type_id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci COMMENT = '节点上的模板关联(tag,tpl,sub_meta)';


DROP VIEW IF EXISTS `tag_role_user_token` ;
DROP TABLE IF EXISTS `tag_role_user_token`;
CREATE OR REPLACE 
VIEW `tag_role_user_token` AS
    SELECT 
        `a`.`tag_id` AS `user_tag_id`,
        `b`.`tag_id` AS `token_tag_id`,
        `a`.`tpl_id` AS `role_id`,
        `a`.`sub_id` AS `user_id`,
        `b`.`sub_id` AS `token_id`
    FROM `tpl_rel` `a` JOIN `tpl_rel` `b`
	ON `a`.`type_id` = 0 AND `b`.`type_id` = 1
	AND `a`.`tpl_id` = `b`.`tpl_id`;


DROP VIEW IF EXISTS `acl`;
DROP TABLE IF EXISTS `acl`;
CREATE OR REPLACE 
VIEW `acl` AS
    SELECT 
        `a`.`user_id` AS `user_id`,
        `a`.`token_id` AS `token_id`,
        `a`.`user_tag_id` AS `user_tag_id`,
        `a`.`token_tag_id` AS `token_tag_id`,
        `a`.`role_id` AS `role_id`
    FROM
        (`tag_role_user_token` `a`
        JOIN `tag_rel` `b` ON (((`a`.`user_tag_id` = `b`.`tag_id`)
            AND (`a`.`token_tag_id` = `b`.`sup_tag_id`))))
    GROUP BY `a`.`user_tag_id`,`a`.`token_id`,`a`.`user_id`
    HAVING (`a`.`token_tag_id` = MAX(`a`.`token_tag_id`));



SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;

INSERT INTO `tag` (`name`) VALUES ('');
INSERT INTO `user` (`uuid`, `name`, `cname`, `email`, `phone`, `im`, `qq`) VALUES
    ('root@localhost', 'system', 'system', 'root@localhost', '', '', '');
