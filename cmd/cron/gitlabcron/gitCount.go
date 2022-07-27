package gitlabcron

import (
	"axis/common/tpl"
	"encoding/json"
	"github.com/hibiken/asynq"
	"log"
	"time"
)

type GitCron struct {
}

func (s *GitCron) UpdatePubData(scheduler *asynq.Scheduler) {

	payload, err := json.Marshal(tpl.GitlabPayload{Email: "429390237@qq.com", Content: "发邮件呀"})
	if err != nil {
		log.Fatal(err)
	}

	task := asynq.NewTask(tpl.DealGitCountMqTpl, payload)
	// 每隔1个小时同步一次
	entryID, err := scheduler.Register("0 */1 * * *", task, asynq.Timeout(time.Minute*30))

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registered an entry: %q\n", entryID)
}
