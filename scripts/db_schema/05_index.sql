SET NAMES utf8;

DROP TABLE if exists `endpoint`;
CREATE TABLE `endpoint` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `endpoint` varchar(255) NOT NULL DEFAULT '',
  `ts` int(11) DEFAULT NULL,
  `t_create` DATETIME NOT NULL COMMENT 'create time',
  `t_modify` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'last modify time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_endpoint` (`endpoint`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE if exists `endpoint_counter`;
CREATE TABLE `endpoint_counter` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `endpoint_id` int(10) unsigned NOT NULL,
  `counter` varchar(255) NOT NULL DEFAULT '' COMMENT 'metric/tags/typ',
  `ts` int(11) DEFAULT NULL,
  `t_create` DATETIME NOT NULL COMMENT 'create time',
  `t_modify` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'last modify time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_endpoint_id_counter` (`endpoint_id`, `counter`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE if exists `tag_endpoint`;
CREATE TABLE `tag_endpoint` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `tag` varchar(255) NOT NULL DEFAULT '' COMMENT 'srv=tv',
  `endpoint_id` int(10) unsigned NOT NULL,
  `ts` int(11) DEFAULT NULL,
  `t_create` DATETIME NOT NULL COMMENT 'create time',
  `t_modify` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'last modify time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tag_endpoint_id` (`tag`, `endpoint_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

