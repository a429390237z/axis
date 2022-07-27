package config

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Config struct {
	service.ServiceConf

	Cache cache.CacheConf

	Redis redis.RedisConf

	Mysql struct {
		DataSource string
	}
}