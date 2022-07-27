package cronMq

import (
	"axis/cmd/mq/internal/handler"
	"axis/cmd/mq/internal/svc"
	"axis/common/tpl"
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
)

/**
定时任务
*/
type AsynqTask struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAsynqTask(ctx context.Context, svcCtx *svc.ServiceContext) *AsynqTask {
	return &AsynqTask{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AsynqTask) Start() {

	fmt.Println("AsynqTask start ")

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: l.svcCtx.Config.Redis.Host, Password: l.svcCtx.Config.Redis.Pass},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()

	// 注册定时任务处理函数
	// 注册阿里云日志服务定时同步数据
	//mux.HandleFunc(tpl.SyncSlsDataTpl, handler.SyncSlsDataMqHandler(l.svcCtx))
	//mux.HandleFunc(tpl.DealLogMqTpl, handler.DealLogMqHandler(l.svcCtx))

	// 注册gitlab更新发布任务
	mux.HandleFunc(tpl.DealGitCountMqTpl, handler.DealGitCountHandler(l.svcCtx))

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func (l *AsynqTask) Stop() {
	fmt.Println("AsynqTask stop")

}
