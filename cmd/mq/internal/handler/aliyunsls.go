package handler

import (
	"axis/cmd/mq/internal/svc"
	"axis/common/tpl"
	"axis/common/utils"
	"axis/model"
	"context"
	"encoding/json"
	"errors"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/fx"
	"github.com/zeromicro/go-zero/core/logx"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
)

type Handler func(ctx context.Context, t *asynq.Task) error

type AK struct {
	AccessKey    string
	AccessSecret string
	Account      string
}

type LogStore struct {
	Owner        string
	Endpoint     string
	LogStoreName string
	ProjectName  string
}

type SyncSlsDataTask struct {
	AK
	EntryPoint string
}

type LogStoreLinesLog struct {
	*LogStore
	Lines []map[string]string
}

type BeginLogStore struct {
	EndPoint string
	AK
	LogS []*LogStore
}

type GCount struct {
	Count int64
	mux   *sync.RWMutex
}

// 定义不需要的日志库，过滤掉
var ignoreLogStore = []string{
	"k8s-event",
	"nginx-ingress",
	"config-operation-log",
	"traefik",
	"admin-ui",
	"sas-log",
	"internal-alert-history",
	"internal-operation_log",
	"waf-logstore",
	"function-log",
	"internal-diagnostic_log",
	"redis_audit_log",
	"internal-alert-center-log",
	"internal-etl-log",
	"redis_slow_run_log",
	"nginx-ingress-metrics",
	"nginx-ingress-metrics-result",
	"oss-log-store",
	"redis_audit_log_standard",
	"rds_log",
	"internal-etl-log",
}

/**
处理日志敏感数据
*/

func DealLogMqHandler(svcContext *svc.ServiceContext) Handler {
	return func(ctx context.Context, t *asynq.Task) error {
		defer func() {
			if err := recover(); err != nil {
				var p tpl.SyncSlsPayload
				if err := json.Unmarshal(t.Payload(), &p); err != nil {
					logx.WithContext(ctx).Errorf("解析payload失败！错误信息：%v", err)
					return
				}
				// todo: 发送邮件
				logx.WithContext(ctx).Infof("payload：%v", p)
			}
		}()

		var (
			owner    string
			endpoint string
			once     *sync.Once
			mapAks   = make(map[string]AK)
			//beginCh = make(chan interface{})
			//logStoreDneCh = syncx.NewDoneChan()

			// 全局计数器，计算处理的日志条数
			gCount = &GCount{
				Count: 0,
				mux:   &sync.RWMutex{},
			}
		)

		// 需要处理的logStore时间范围
		yesTimeStr := time.Now().UTC().Add(-24 * time.Hour).Format("2006-01-02")
		tm2, _ := time.Parse("2006-01-02", yesTimeStr)
		beginTime := tm2.Unix()
		endTime := tm2.Unix() + 23*60*60 + 59*60 + 59

		// 手机号正则对象，后面要用
		phoneReg := `(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\\d{8}`
		reg := regexp.MustCompile(phoneReg)
		if reg == nil {
			logx.WithContext(ctx).Error("MustCompile 失败！")
			return errors.New("MustCompile failed")
		}

		// 获取AK
		AKs, err := getAk(ctx, svcContext)
		if err != nil {
			return err
		}
		for _, AK := range AKs {
			mapAks[AK.Account] = AK
		}

		fx.From(func(source chan<- interface{}) {
			// 获取所有logstore列表
			endLogStore, err := getLogStore(ctx, svcContext)
			if err != nil {
				logx.WithContext(ctx).Errorf("获取logStore失败，错误信息：%v", err)
				return
			}
			for k, v := range endLogStore {
				owner = strings.Split(k, "__")[0]
				endpoint = strings.Split(k, "__")[1]
				source <- &BeginLogStore{
					EndPoint: endpoint,
					AK:       mapAks[owner],
					LogS:     v,
				}
			}
		}).Walk(func(item interface{}, pipe chan<- interface{}) {
			begin := item.(*BeginLogStore)
			FilterLogStore(ctx, beginTime, endTime, begin, pipe, gCount)
		}).Walk(func(item interface{}, pipe chan<- interface{}) {
			logs := item.(*LogStoreLinesLog)
			dealLineLog(ctx, svcContext, logs, reg, once, tm2)
		}).Done()

		logx.WithContext(ctx).Infof("本次任务结束, 共处理数据%d条", gCount.Count)

		return nil

	}
}

/**
阿里云日志服务定时任务处理
*/
func SyncSlsDataMqHandler(svcContext *svc.ServiceContext) Handler {

	return func(ctx context.Context, t *asynq.Task) error {
		defer func() {
			if err := recover(); err != nil {
				var p tpl.SyncSlsPayload
				if err := json.Unmarshal(t.Payload(), &p); err != nil {
					logx.WithContext(ctx).Errorf("解析payload失败！错误信息：%v", err)
					return
				}
				//todo: 发送邮件
				logx.WithContext(ctx).Infof("payload：%v", p)
			}
		}()

		var (
			endChan  = make(chan struct{})
			taskChan = make(chan *SyncSlsDataTask)
		)

		// 1. 获取AK信息
		AKs, err := getAk(ctx, svcContext)
		if err != nil {
			return err
		}

		// 2. 获取各个地区日志服务外网访问地址
		entryPoints, err := getEntryPoint(ctx, svcContext)
		if err != nil {
			return err
		}

		// 3. 拼装任务
		go func() {
			defer func() {
				close(taskChan)
				endChan <- struct{}{}
			}()
			for _, ak := range AKs {
				for _, entryPoint := range entryPoints {
					taskChan <- &SyncSlsDataTask{
						AK:         ak,
						EntryPoint: entryPoint,
					}
				}
			}
		}()

		// 4. 消费任务
		for i := 0; i < 10; i++ {
			go func() {
				dealProjectAndLogStore(ctx, svcContext, taskChan, endChan)
			}()
		}

		for i := 0; i < 11; i++ {
			<-endChan
		}
		logx.WithContext(ctx).Infof("本次任务结束！")
		return nil
	}
}

// 获取AK信息
func getAk(ctx context.Context, svcContext *svc.ServiceContext) (AKs []AK, err error) {
	aliyunAks, err := svcContext.AliyunAkModel.FindAll()
	if err != nil {
		return nil, err
	}

	for _, aliyunAk := range aliyunAks {
		if aliyunAk.AccessKeyID.String == "LTAI5tHsikvn6jAeqn6Atc5H" {
			continue
		}
		ak := AK{
			AccessKey:    aliyunAk.AccessKeyID.String,
			AccessSecret: aliyunAk.AccessKey.String,
			Account:      aliyunAk.Account.String,
		}
		AKs = append(AKs, ak)
	}
	return
}

// 获取阿里云日志服务各地域公网地址
func getEntryPoint(ctx context.Context, svcContext *svc.ServiceContext) (entryPoints []string, err error) {
	aliyunLogEntrances, err := svcContext.AliyunLogEntranceModel.FindAll()
	if err != nil {
		//logx.WithContext(ctx).Errorf("获取日志服务公网地址失败，错误信息:%v", err)
		return nil, err
	}

	for _, entryPoint := range aliyunLogEntrances {
		entryPoints = append(entryPoints, entryPoint.InternetEntrance)
	}
	return
}

// 获取账号所属的logStore
func getLogStore(ctx context.Context, svcContext *svc.ServiceContext) (endLogStore map[string][]*LogStore, err error) {
	aliyunLogStores, err := svcContext.AliyunLogStoreModel.FindAll()
	if err != nil {
		//logx.WithContext(ctx).Errorf("获取日志服务公网地址失败，错误信息:%v", err)
		return nil, err
	}

	endLogStore = make(map[string][]*LogStore)

	for _, aliyunLogStore := range aliyunLogStores {
		// 过滤不需要的日志
		var flag = false
		for _, ignore := range ignoreLogStore {
			if aliyunLogStore.Name == ignore || strings.HasPrefix(aliyunLogStore.Name, "audit-") {
				flag = true
				break
			}
		}
		if flag == true {
			continue
		}

		idx := aliyunLogStore.Owner + "__" + aliyunLogStore.Endpoint
		if _, ok := endLogStore[idx]; !ok {
			endLogStore[idx] = make([]*LogStore, 0, 100)
		}

		endLogStore[idx] = append(endLogStore[idx], &LogStore{
			Owner:        aliyunLogStore.Owner,
			Endpoint:     aliyunLogStore.Endpoint,
			LogStoreName: aliyunLogStore.Name,
			ProjectName:  aliyunLogStore.ProjectName,
		})
	}
	return
}

// 处理project和logStore
func dealProjectAndLogStore(ctx context.Context, svxContext *svc.ServiceContext, taskChan chan *SyncSlsDataTask, endChan chan struct{}) {

	// 执行完成，发送信号
	defer func() {
		endChan <- struct{}{}
	}()

	for task := range taskChan {
		var (
			offset      = 0
			logProjects = make([]sls.LogProject, 0, 100)
			err         error
		)
		client := sls.CreateNormalInterface(task.EntryPoint, task.AccessKey, task.AccessSecret, "")

		// 获取project
		for {
			projects, count, total, err := client.ListProjectV2(offset, 100)
			if err != nil {
				logx.WithContext(ctx).Errorf("获取project失败，错误信息:%v", err)
				break
			}

			//logx.WithContext(ctx).Infof("地域：%s,获取到的projects数量%d,为：%#v", task.EntryPoint, len(projects), projects)

			for _, project := range projects {
				logProjects = append(logProjects, project)
			}
			if offset+count >= total {
				break
			}
			offset += count
		}

		//logx.WithContext(ctx).Infof("地域：%s, 有%d个project, 为：%#v", task.EntryPoint, len(logProjects), logProjects)

		// 处理project
		for _, project := range logProjects {
			var (
				status int64

				enableTracking int64
				autoSplit      int64
				appendMeta     int64
				telemetryType  int64
			)

			if project.Status == "Normal" {
				status = 1
			} else {
				status = 0
			}

			iCreateTime, _ := strconv.Atoi(project.CreateTime)
			iLastModifyTime, _ := strconv.Atoi(project.LastModifyTime)

			createTime := time.Unix(int64(iCreateTime), 0)
			lastModifyTime := time.Unix(int64(iLastModifyTime), 0)

			aliyunLogProject := &model.AliyunLogProject{
				Name:           project.Name,
				Description:    project.Description,
				Region:         project.Region,
				Status:         status,
				Owner:          task.Account,
				Maintainer:     "",
				CreateTime:     createTime,
				LastModifyTime: lastModifyTime,
				Endpoint:       task.EntryPoint,
			}

			err = svxContext.AliyunLogProjectModel.InsertOrUpdate(nil, aliyunLogProject)
			if err != nil {
				logx.WithContext(ctx).Errorf("写入数据失败, 错误信息:%v", err)
				return
			}

			// 列出当前project下所有的logStore
			logStoreStrs, err := client.ListLogStore(project.Name)
			if err != nil {
				logx.WithContext(ctx).Errorf("获取project下logStore失败，错误信息：%v", err)
				return
			}

			//logx.WithContext(ctx).Infof("地域：%s,当前project为：%s, 有%d个logStore", task.EntryPoint, project.Name, len(logStoreStrs))

			for _, logStoreStr := range logStoreStrs {

				logStore, err := client.GetLogStore(project.Name, logStoreStr)
				if err != nil {
					logx.WithContext(ctx).Errorf("获取logStore失败, 错误信息: %v", err)
					return
				}

				if logStore.WebTracking {
					enableTracking = 1
				} else {
					enableTracking = 0
				}

				if logStore.AutoSplit {
					autoSplit = 1
				} else {
					autoSplit = 0
				}

				if logStore.AppendMeta {
					appendMeta = 1
				} else {
					appendMeta = 0
				}

				if logStore.TelemetryType == "None" {
					telemetryType = 0
				} else {
					telemetryType = 1
				}

				storeCreateTime := time.Unix(int64(logStore.CreateTime), 0)
				storeLastModifyTime := time.Unix(int64(logStore.LastModifyTime), 0)

				aliyunLogstore := &model.AliyunLogstore{
					Name:           logStore.Name,
					ProjectName:    project.Name,
					Endpoint:       task.EntryPoint,
					Ttl:            int64(logStore.TTL),
					ShardCount:     int64(logStore.ShardCount),
					EnableTracking: enableTracking,
					AutoSplit:      autoSplit,
					MaxSplitShard:  int64(logStore.MaxSplitShard),
					AppendMeta:     appendMeta,
					TelemetryType:  telemetryType,
					CreateTime:     storeCreateTime,
					LastModifyTime: storeLastModifyTime,
					Owner:          task.Account,
					Maintainer:     "",
				}

				err = svxContext.AliyunLogStoreModel.InsertOrUpdate(nil, aliyunLogstore)
				if err != nil {
					logx.WithContext(ctx).Errorf("写入logStore失败，错误信息：%v", err)
					return
				}
			}
		}

		if err = client.Close(); err != nil {
			logx.WithContext(ctx).Errorf("关闭日志服务客户端失败，错误信息：%v", err)
		}
	}
}

// 处理logStore
func FilterLogStore(ctx context.Context, beginTime int64, endTime int64, begin *BeginLogStore, pipe chan<- interface{}, gCount *GCount) {

	client := sls.CreateNormalInterface(begin.EndPoint, begin.AccessKey, begin.AccessSecret, "")

	for _, logStore := range begin.LogS {

		// GetLogs API Ref: https://intl.aliyun.com/help/doc-detail/29029.htm
		glResp, err := client.GetLogs(logStore.ProjectName, logStore.LogStoreName, "", beginTime, endTime, "*|SELECT  * LIMIT 100000", 0, 0, false)
		if err != nil {
			//logx.WithContext(ctx).Errorf("获取日志失败, 错误信息: %v", err)
			continue
			//time.Sleep(10 * time.Millisecond)
			//continue
		}
		if glResp.Count > 0 {
			pipe <- &LogStoreLinesLog{
				LogStore: logStore,
				Lines:    glResp.Logs,
			}
		}

		logx.WithContext(ctx).Infof("当前logStore: %s, 日志条数：%d", logStore.LogStoreName, glResp.Count)

		gCount.mux.Lock()
		gCount.Count += glResp.Count
		gCount.mux.Unlock()

	}

	if err := client.Close(); err != nil {
		logx.WithContext(ctx).Errorf("关闭日志服务客户端失败，错误信息：%v", err)
	}

}

func dealLineLog(ctx context.Context, svcContext *svc.ServiceContext, logs *LogStoreLinesLog, reg *regexp.Regexp, once *sync.Once, dt time.Time) {

	var (
		// 记录一条敏感数据
		resData string
		// 统计敏感数据次数
		count int64
		err   error
	)

	for _, line := range map2Str(logs.Lines) {
		if res := reg.FindString(line); res != "" {
			once.Do(func() {
				// 超过长度截取
				if len(res) > 512 {
					res = utils.SubStrRuneIndexInString(res, 509) + "..."
				}
				resData = res
			})
			count += 1
			logx.WithContext(ctx).Errorf("Get phone info: %v", res)
		}
	}

	if count > 0 {
		// 拼接数据
		logRecord := &model.AliyunLogRecord{
			Date:         dt,
			ProjectName:  logs.ProjectName,
			LogstoreName: logs.LogStoreName,
			Tp:           1,
			Count:        count,
			Info:         resData,
			ModifyTime:   time.Now(),
		}

		logx.WithContext(ctx).Infof("当前logStore: %s,敏感数据条数: %d", logs.LogStoreName, count)

		err = svcContext.AliyunLogRecordModel.InsertOrUpdate(nil, logRecord)
		if err != nil {
			logx.WithContext(ctx).Errorf("写入aliyunLogRecord失败，错误信息：%v", err)
			return
		}
	}
}

func map2Str(m []map[string]string) []string {
	m2 := make([]string, len(m))
	for i, data := range m {
		var build strings.Builder
		for k, v := range data {
			if strings.HasPrefix(k, "__") {
				continue
			}
			build.WriteString(v)
			build.WriteString(",")
		}
		m2[i] = strings.TrimRight(build.String(), ",")
	}
	return m2
}
