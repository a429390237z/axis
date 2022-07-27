package main

import (
	"axis/cmd/cron/gitlabcron"
	"github.com/hibiken/asynq"
	"log"
	"time"
)

// 配置文件内容比较少，直接定义
const redisAddr = "192.168.19.98:6379"
const redisPwd = "tiantong99.c0m"

func main() {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	// 周期性任务
	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{
			Addr:     redisAddr,
			Password: redisPwd,
		}, &asynq.SchedulerOpts{
			Location: loc,
		})

	defer scheduler.Shutdown()

	// 配置阿里云sls定时任务
	//slsCron := &aliyunsls.SlsCron{}
	//slsCron.SyncSlsData(scheduler)
	//slsCron.DealLogMq(scheduler)

	// 配置gitlab定时任务
	gitCron := &gitlabcron.GitCron{}
	gitCron.UpdatePubData(scheduler)

	if err := scheduler.Run(); err != nil {
		log.Printf("scheduler err:%v\n", err)
	}

}
