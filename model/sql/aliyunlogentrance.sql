set NAMES utf8mb4;
set foreign_key_checks = 0;

DROP TABLE IF EXISTS `aliyun_log_entrance`;
CREATE TABLE `aliyun_log_entrance` (
                                       `id` int NOT NULL AUTO_INCREMENT,
                                       `region` varchar(32) NOT NULL DEFAULT '' COMMENT '地域(英文）',
                                       `region_cn` varchar(32) NOT NULL DEFAULT '' COMMENT  '地域(中文)',
                                       `internet_entrance` varchar(64) NOT NULL DEFAULT '' COMMENT  '公网入口',
                                       `intranet_entrance` varchar(64) NOT NULL DEFAULT '' COMMENT '私网入口',
                                       PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='阿里云日志入口表';

insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('cn-hangzhou', '华东1（杭州）', 'cn-hangzhou.log.aliyuncs.com', 'cn-hangzhou-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('cn-hangzhou-finance', '华东1（杭州-金融云）', 'cn-hangzhou-finance.log.aliyuncs.com', 'cn-hangzhou-finance-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('cn-shanghai', '华东2（上海）', 'cn-shanghai.log.aliyuncs.com', 'cn-shanghai-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('cn-shanghai-finance-1', '华东2（上海-金融云）', 'cn-shanghai-finance-1.log.aliyuncs.com', 'cn-shanghai-finance-1-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('cn-qingdao', '华北1（青岛）', 'cn-qingdao.log.aliyuncs.com', 'cn-qingdao-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('cn-beijing', '华北2（北京）', 'cn-beijing.log.aliyuncs.com', 'cn-beijing-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('cn-north-2-gov-1', '华北2 阿里政务云1', 'cn-north-2-gov-1.log.aliyuncs.com', 'cn-north-2-gov-1-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('cn-zhangjiakou', '华北3（张家口）', 'cn-zhangjiakou.log.aliyuncs.com', 'cn-zhangjiakou-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('cn-huhehaote', '华北5（呼和浩特）', 'cn-huhehaote.log.aliyuncs.com', 'cn-huhehaote-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('cn-wulanchabu', '华北6（乌兰察布）', 'cn-wulanchabu.log.aliyuncs.com', 'cn-wulanchabu-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('cn-shenzhen', '华南1（深圳）', 'cn-shenzhen.log.aliyuncs.com', 'cn-shenzhen-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('cn-shenzhen-finance', '华南1（深圳-金融云）', 'cn-shenzhen-finance.log.aliyuncs.com', 'cn-shenzhen-finance-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('cn-heyuan', '华南2（河源）', 'cn-heyuan.log.aliyuncs.com', 'cn-heyuan-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('cn-guangzhou', '华南3（广州）', 'cn-guangzhou.log.aliyuncs.com', 'cn-guangzhou-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('cn-chengdu', '西南1（成都）', 'cn-chengdu.log.aliyuncs.com', 'cn-chengdu-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('cn-hongkong', '中国（香港）', 'cn-hongkong.log.aliyuncs.com', 'cn-hongkong-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('ap-northeast-1', '日本（东京）', 'ap-northeast-1.log.aliyuncs.com', 'ap-northeast-1-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('ap-southeast-1', '新加坡', 'ap-southeast-1.log.aliyuncs.com', 'ap-southeast-1-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('ap-southeast-2', '澳大利亚（悉尼）', 'ap-southeast-2.log.aliyuncs.com', 'ap-southeast-2-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('ap-southeast-3', '马来西亚（吉隆坡）', 'ap-southeast-3.log.aliyuncs.com', 'ap-southeast-3-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('ap-southeast-6', '菲律宾（马尼拉）', 'ap-southeast-6.log.aliyuncs.com', 'ap-southeast-6-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('ap-southeast-5', '印度尼西亚（雅加达）', 'ap-southeast-5.log.aliyuncs.com', 'ap-southeast-5-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('me-east-1', '阿联酋（迪拜）', 'me-east-1.log.aliyuncs.com', 'me-east-1-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('us-west-1', '美国（硅谷）', 'us-west-1.log.aliyuncs.com', 'us-west-1-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('eu-central-1', '德国（法兰克福）', 'eu-central-1.log.aliyuncs.com', 'eu-central-1-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('us-east-1', '美国（弗吉尼亚）', 'us-east-1.log.aliyuncs.com', 'us-east-1-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('ap-south-1', '印度（孟买）', 'ap-south-1.log.aliyuncs.com', 'ap-south-1-intranet.log.aliyuncs.com');
insert into `aliyun_log_entrance` (`region`, `region_cn`, `internet_entrance`, `intranet_entrance`) values ('eu-west-1', '英国（伦敦）', 'eu-west-1.log.aliyuncs.com', 'eu-west-1-intranet.log.aliyuncs.com');

set foreign_key_checks = 1;