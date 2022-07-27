package listen

import (
	"axis/cmd/mq/internal/config"
	"axis/cmd/mq/internal/mqs/cronMq"
	"axis/cmd/mq/internal/svc"
	"context"
	"github.com/zeromicro/go-zero/core/service"
)

// asynq
// 定时任务
func CronMqs(c config.Config, ctx context.Context, svcContext *svc.ServiceContext) []service.Service {

	return []service.Service{
		// 监听定时任务
		cronMq.NewAsynqTask(ctx, svcContext),
	}
}
