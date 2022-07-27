package svc

import (
	"axis/cmd/mq/internal/config"
	"axis/model"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                 config.Config
	Redis                  *redis.Redis
	AliyunAkModel          model.AliyunAkModel
	AliyunLogEntranceModel model.AliyunLogEntranceModel
	AliyunLogProjectModel  model.AliyunLogProjectModel
	AliyunLogStoreModel    model.AliyunLogstoreModel
	AliyunLogRecordModel   model.AliyunLogRecordModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                 c,
		AliyunAkModel:          model.NewAliyunAkModel(sqlx.NewMysql(c.Mysql.DataSource), c.Cache),
		AliyunLogEntranceModel: model.NewAliyunLogEntranceModel(sqlx.NewMysql(c.Mysql.DataSource), c.Cache),
		AliyunLogProjectModel:  model.NewAliyunLogProjectModel(sqlx.NewMysql(c.Mysql.DataSource), c.Cache),
		AliyunLogStoreModel:    model.NewAliyunLogstoreModel(sqlx.NewMysql(c.Mysql.DataSource), c.Cache),
		AliyunLogRecordModel:   model.NewAliyunLogRecordModel(sqlx.NewMysql(c.Mysql.DataSource), c.Cache),
	}
}
