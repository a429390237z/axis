package aliyunsls

import (
	"axis/common/tpl"
	"encoding/json"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

type SlsCron struct {
}

func (s *SlsCron) SyncSlsData(scheduler *asynq.Scheduler) {

	payload, err := json.Marshal(tpl.SyncSlsPayload{Email: "429390237@qq.com", Content: "发邮件呀"})
	if err != nil {
		log.Fatal(err)
	}

	task := asynq.NewTask(tpl.SyncSlsDataTpl, payload)
	// 每隔1分钟同步一次
	entryID, err := scheduler.Register("30 16 * * *", task, asynq.Timeout(time.Hour*5))

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registered an entry: %q\n", entryID)
}

func (s *SlsCron) DealLogMq(scheduler *asynq.Scheduler) {

	payload, err := json.Marshal(tpl.SyncSlsPayload{Email: "429390237@qq.com", Content: "发邮件呀"})
	if err != nil {
		log.Fatal(err)
	}

	task := asynq.NewTask(tpl.DealLogMqTpl, payload)
	// 每隔1分钟同步一次
	entryID, err := scheduler.Register("40 16 * * *", task, asynq.Timeout(time.Hour*5))

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registered an entry: %q\n", entryID)
}
