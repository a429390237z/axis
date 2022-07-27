DROP TABLE  IF EXISTS `aliyun_log_record`;
CREATE TABLE `aliyun_log_record` (
       `id`  bigint  NOT NULL AUTO_INCREMENT,
       `date` date  NOT NULL DEFAULT '0000' COMMENT '统计时间',
       `project_name` varchar(128)  NOT NULL DEFAULT ''  COMMENT 'project名称',
       `logstore_name` varchar(128) NOT NULL DEFAULT ''  COMMENT '日志库名称',
       `type`  int NOT NULL  DEFAULT 0  COMMENT '敏感信息类型：1. 日志存在电话号码',
       `count` int  NOT NULL DEFAULT 0 COMMENT '敏感信息数量',
       `info`  varchar(512) NOT NULL DEFAULT '' COMMENT '单条敏感信息例子：如日志存在电话号码',
       `create_time`   datetime  NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
       `modify_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
       PRIMARY KEY (`id`),
       UNIQUE KEY idx_name (`date`, `project_name`, `logstore_name`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='日志记录表（敏感信息分析）';