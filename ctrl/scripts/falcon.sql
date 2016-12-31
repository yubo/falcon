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
-- Table `config`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `config`;
CREATE TABLE `config` (
  `key` VARCHAR(128) NOT NULL,
  `value` BLOB NOT NULL,
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
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `index_status` (`status`),
  INDEX `index_type` (`type`),
  UNIQUE INDEX `index_name` (`name`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci COMMENT = '机器';


-- -----------------------------------------------------
-- Table `system`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `system`;
CREATE TABLE `system` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` VARCHAR(32) NOT NULL COMMENT '系统名',
  `cname` VARCHAR(64) NOT NULL COMMENT '中文名',
  `developers` VARCHAR(255) NOT NULL COMMENT '系统开发人员',
  `email` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '组邮箱',
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_name` (`name`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_unicode_ci COMMENT = '系统模块';


-- -----------------------------------------------------
-- Table `token`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `token`;
CREATE TABLE `token` (
  `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(32) NOT NULL,
  `system_id` INT(11) UNSIGNED NOT NULL,
  `cname` VARCHAR(64) NOT NULL,
  `note` VARCHAR(255) NOT NULL DEFAULT '',
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index_name` (`name`),
  INDEX `index_system_id` (`system_id`),
  CONSTRAINT `token_rel_ibfk_1`
    FOREIGN KEY (`system_id`)
    REFERENCES `system` (`id`)
    ON DELETE CASCADE)
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
  `rel_type` TINYINT(4) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'relation type',
  PRIMARY KEY (`id`),
  INDEX `index_tag_id` (`tag_id`),
  INDEX `index_sup_tag_id` (`sup_tag_id`),
  INDEX `index_rel_type` (`rel_type`),
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
  `rel_type` TINYINT(4) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'relation type',
  PRIMARY KEY (`id`),
  INDEX `index_tpl_id` (`tpl_id`),
  INDEX `index_tag_id` (`tag_id`),
  INDEX `index_sub_id` (`sub_id`),
  INDEX `index_rel_type` (`rel_type`))
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
	ON `a`.`rel_type` = 0 AND `b`.`rel_type` = 1
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
    ('admin@falcon', 'admin', 'admin', 'root@localhost', '', '', '');
-- 
-- 
-- INSERT INTO `role` (`name`, `cname`, `note`) VALUES
--     ('admin', '超级管理员', '配置所有选项'),
--     ('manager', '管理员', 'user'),
--     ('sre', '工程师', 'Site Reliability Engineering'),
--     ('user', '普通用户', 'user');
-- 
-- 
-- INSERT INTO `system` (`name`, `cname`, `developers`, `email`) VALUES
--     ('falcon', '机器管理', 'yubo', 'yubo@xiaomi.com');
-- 
-- 
-- INSERT INTO `token` (`name`, `system_id`, `cname`, `note`) VALUES
--     ('falcon-tag-edit', 1, '节点修改', '允许添加删除节点'),
--     ('falcon-host-operate', 1, '机器操作', '重启，改名，关机'),
--     ('falcon-host-bind', 1, '机器挂载', '允许挂载，删除机器与节点的对应关系'),
--     ('falcon-tag-read', 1, '节点读', '查看节点及节点下相关内容');
-- 
-- 
-- INSERT INTO `tag_host` (`tag_id`, `host_id`) VALUES
--     (2, 1),
--     (2, 2),
--     (2, 3),
--     (3, 4),
--     (3, 5),
--     (3, 6),
--     (4, 7),
--     (4, 8),

-- INSERT INTO `tag_rel` (`tag_id`, `sup_tag_id`) VALUES
--     (2, 1),
--     (2, 2),
--     (3, 1),
--     (3, 2),
--     (3, 3),
--     (4, 1),
--     (4, 2),
--     (4, 3),
--     (4, 4),
--     (5, 1),
--     (5, 2),
--     (5, 5),
--     (6, 1),
--     (6, 2),
--     (6, 5),
--     (6, 6),
--     (7, 1),
--     (7, 2),
--     (7, 5),
--     (7, 6),
--     (7, 7),
--     (8, 1),
--     (8, 2),
--     (8, 5),
--     (8, 6),
--     (8, 8);
-- 
-- -- 1
-- INSERT INTO `host` (`uuid`, `name`, `type`, `status`, `loc`, `idc`) VALUES
--     ('xiaomi1', 'c3-mi01.bj', 'machine', 'online', 'bj', 'lg'),
--     ('xiaomi2', 'c3-mi02.bj', 'machine', 'online', 'bj', 'lg'),
--     ('xiaomi3', 'c3-mi03.bj', 'machine', 'online', 'bj', 'lg'),
--     ('inf1', 'c3-inf01.bj', 'machine', 'online', 'bj', 'lg'),
--     ('inf2', 'c3-inf02.bj', 'machine', 'online', 'bj', 'lg'),
--     ('inf3', 'c3-inf03.bj', 'machine', 'online', 'bj', 'lg'),
--     ('falcon1', 'c3-inf-falcon01.bj', '21vianet', 'online', 'bj', 'c3'),
--     ('falcon2', 'c3-inf-falcon02.bj', '21vianet', 'online', 'bj', 'c3'),
--     ('falcon3', 'c3-inf-falcon03.bj', '21vianet', 'online', 'bj', 'c3'),
--     ('miliao1', 'c3-miliao01.bj', 'machine', 'online', 'bj', 'lg'),
--     ('miliao2', 'c3-miliao02.bj', 'machine', 'online', 'bj', 'lg'),
--     ('miliao3', 'c3-miliao03.bj', 'machine', 'online', 'bj', 'lg'),
--     ('op1', 'c3-op01.bj', 'machine', 'online', 'bj', 'lg'),
--     ('op2', 'c3-op02.bj', 'machine', 'online', 'bj', 'lg'),
--     ('op3', 'c3-op03.bj', 'machine', 'online', 'bj', 'lg'),
--     ('dpdk1', 'c3-op-mon-dpdk01.bj', 'machine', 'online', 'bj', 'lg'),
--     ('dpdk2', 'c3-op-mon-dpdk02.bj', 'machine', 'online', 'bj', 'lg'),
--     ('dpdk3', 'c3-op-mon-dpdk03.bj', 'machine', 'online', 'bj', 'lg'),
--     ('uaq1', 'c3-op-mon-uaq01.bj', 'machine', 'online', 'bj', 'lg'),
--     ('uaq2', 'c3-op-mon-uaq02.bj', 'machine', 'online', 'bj', 'lg'),
--     ('uaq3', 'c3-op-mon-uaq03.bj', 'machine', 'online', 'bj', 'lg');
-- 
-- INSERT INTO `user` (`uuid`, `name`, `cname`, `email`, `phone`, `im`, `qq`) VALUES
--     ('cn=yubo,ou=users,dc=yubo,dc=org@ldap', 'yubo', 'yubo', 'yubo@yubo.org', '110', 'x80386', '20507'),
--     ('cn=tom,ou=users,dc=yubo,dc=org@ldap', 'tom', 'tom', 'tom@yubo.org', '111', '4567', '1234');
-- 
-- 
-- INSERT INTO `role` (`name`, `cname`, `note`) VALUES
--     ('admin', '超级管理员', '配置所有选项'),
--     ('manager', '管理员', 'user'),
--     ('sre', '工程师', 'Site Reliability Engineering'),
--     ('user', '普通用户', 'user');
-- 
-- 
-- INSERT INTO `system` (`name`, `cname`, `developers`, `email`) VALUES
--     ('falcon', '机器管理', 'yubo', 'yubo@xiaomi.com');
-- 
-- 
-- INSERT INTO `token` (`name`, `system_id`, `cname`, `note`) VALUES
--     ('falcon-tag-edit', 1, '节点修改', '允许添加删除节点'),
--     ('falcon-host-operate', 1, '机器操作', '重启，改名，关机'),
--     ('falcon-host-bind', 1, '机器挂载', '允许挂载，删除机器与节点的对应关系'),
--     ('falcon-tag-read', 1, '节点读', '查看节点及节点下相关内容');
-- 
-- 
-- INSERT INTO `tag_host` (`tag_id`, `host_id`) VALUES
--     (2, 1),
--     (2, 2),
--     (2, 3),
--     (3, 4),
--     (3, 5),
--     (3, 6),
--     (4, 7),
--     (4, 8),
--     (4, 9),
--     (5, 10),
--     (5, 11),
--     (5, 12),
--     (6, 13),
--     (6, 14),
--     (6, 15),
--     (7, 16),
--     (7, 17),
--     (7, 18),
--     (8, 19),
--     (8, 20),
--     (8, 21);
-- 
-- INSERT INTO `tag_role_token` (`tag_id`, `role_id`, `token_id`) VALUES
--     (1, 1, 1),
--     (1, 1, 2),
--     (1, 1, 3),
--     (1, 1, 4),
--     (1, 2, 1),
--     (1, 2, 2),
--     (1, 2, 3),
--     (1, 2, 4),
--     (1, 3, 2),
--     (1, 3, 3),
--     (1, 3, 4),
--     (1, 4, 4);
-- 
-- 
-- INSERT INTO `tag_role_user` (`tag_id`, `role_id`, `user_id`) VALUES
--     (4, 3, 1);
-- 
