CREATE TABLE `aliyun_ak` (
         `ID` int NOT NULL AUTO_INCREMENT,
         `Account` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '账号ID',
         `PrimaryAccount` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '主账号',
         `SecondaryAccount` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '子账号',
         `AccessKeyID` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT 'AccessKeyID',
         `AccessKey` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL COMMENT 'AccessKey',
         `Permission` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '权限',
         `Info` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '说明',
         `CreateTime` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
         PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=57 DEFAULT CHARSET=utf8mb3 COLLATE=utf8_unicode_ci COMMENT='阿里云账号信息';