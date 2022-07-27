set NAMES utf8mb4;
set foreign_key_checks = 0;

DROP TABLE  IF EXISTS `aliyun_logstore`;
CREATE TABLE `aliyun_logstore` (
    `id`  bigint  NOT NULL AUTO_INCREMENT,
    `name` varchar(128)  NOT NULL DEFAULT '' COMMENT '日志库名称',
    `project_name` varchar(128)  NOT NULL DEFAULT ''  COMMENT 'project名称',
    `endpoint`      varchar(32)   NOT NULL DEFAULT '' COMMENT '地域',
    `ttl`           int    NOT NULL  DEFAULT 0  COMMENT '数据的保存时间',
    `shardCount`    int    NOT NULL  DEFAULT 0  COMMENT 'shard分区数',
    `enable_tracking` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否开启Webtracking功能: 1:开启 0：关闭',
    `auto_split`      tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否自动分裂shard: 1:自动分裂 0：不自动分裂',
    `max_split_shard` int  NOT NULL DEFAULT 0  COMMENT '自动分裂时最大的shard个数，最小值为1，最大值为64',
    `appendMeta`      tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否记录外网IP地址的功能：1：记录 0：不记录',
    `telemetry_type`  tinyint(1)    NOT NULL DEFAULT 0 COMMENT '要查询的日志类型：1：Metrics(时序数据）0：None：非时序存储',
    `create_time`   datetime  NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '日志库创建时间',
    `last_modify_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '日志库更新时间',
    `owner`         varchar(100) NOT NULL DEFAULT '' COMMENT '日志库拥有者',
    `maintainer`    varchar(100) NOT NULL DEFAULT '' COMMENT '日志库维护者',
    PRIMARY KEY (`id`),
    UNIQUE KEY idx_name (`owner`, `endpoint`, `project_name`, `name`)
    /*FOREIGN KEY (`project_id`) REFERENCES `aliyun_log_project`(`id`)*/
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='日志库表';

DROP TABLE IF EXISTS `aliyun_log_project`;
CREATE TABLE `aliyun_log_project` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `name` varchar(128) NOT NULL DEFAULT '' COMMENT 'project名称：作为Host的一部分，project名称在阿里云地域内全局唯一,创建后不可修改',
    `description`  varchar(256) NOT NULL  DEFAULT '' COMMENT 'Project描述',
    `region`  varchar(32) NOT NULL DEFAULT '' COMMENT 'project所有地域',
    `status`  tinyint(1) NOT NULL DEFAULT 0 COMMENT 'project状态：1：Normal(正常）0：Disable(禁用)',
    `owner`   varchar(100) NOT NULL DEFAULT '' COMMENT '日志项目拥有者',
    `maintainer` varchar(100) NOT NULL DEFAULT '' COMMENT '日志项目维护者',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'project创建时间',
    `last_modify_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后一次更新project时间',
    `endpoint` varchar(32) NOT NULL DEFAULT '' COMMENT '地域',
    PRIMARY KEY (`id`),
    UNIQUE KEY idx_region_name (`owner`, `endpoint`, `name`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='日志项目表';

set foreign_key_checks  = 1;
