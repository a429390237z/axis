package listen

import (
	"axis/cmd/mq/internal/config"
	"axis/cmd/mq/internal/svc"
	"context"
	"github.com/zeromicro/go-zero/core/service"
)

//返回所有消费者
func Mqs(c config.Config) []service.Service {

	svcContext := svc.NewServiceContext(c)
	ctx := context.Background()

	var services []service.Service

	//asynq ： 定时任务
	services = append(services, CronMqs(c, ctx, svcContext)...)
	//other mq ....

	return services
}
