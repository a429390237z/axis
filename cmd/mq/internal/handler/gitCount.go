package handler

import (
	"axis/cmd/mq/internal/svc"
	"axis/common/tpl"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/xanzy/go-gitlab"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"net/http"
	"sync"
	"time"
)

var startEndTimes []StartEndTime

var backend = []map[int]string{
	{5367: "main"},
	{5368: "main"},
	{5384: "main"},
	{5387: "main"},
	{5389: "main"},
	{5390: "main"},
	{6830: "main"},
	{6899: "main"},
	{7242: "master"},
	{8113: "master"},
	{8118: "master"},
	{8119: "master"},
	{8120: "master"},
	{8214: "master"},
	{8957: "master"},
	{9096: "master"},
	{9156: "production/shouchuang"},
	{9156: "production/shouchuang"},
	{9438: "master"},
	{9889: "master"},
	{9995: "master"},
	{10233: "master"},
	{10822: "master"},
	{10848: "master"},
	{10905: "master"},
	{10959: "master"},
	{11060: "master"},
	{11116: "master"},
	{11118: "master"},
	{11422: "master"},
	{11451: "master"},
	{11720: "master"},
	{11793: "master"},
	{11794: "master"},
	{11923: "master"},
}

type optFunc func(options *gitlab.ListProjectPipelinesOptions) error

type StartEndTime struct {
	startTime *time.Time
	endTime   *time.Time
}

type GitClient struct {
	git         *gitlab.Client
	projectChan chan *ProjectB
	timeChan    chan *StartEndTime
	wg          *sync.WaitGroup
}

type ProjectB struct {
	*gitlab.Project
	branch string
}

func DealGitCountHandler(svcContext *svc.ServiceContext) Handler {
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

		git, err := GetGitClient("S5yz1sNARp3jH9DWpAS9", "http://gitlab.yintech.net/api/v4/")

		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}

		nowTime := time.Now().UTC()
		timeBeforeOneHour := nowTime.Add(-time.Hour)

		git.wg.Add(2)
		go git.GetProjects()
		go git.listPipelineByTimeV2(&timeBeforeOneHour, &nowTime)
		git.wg.Wait()

		return nil
	}
}

func GetGitClient(token, baseUrl string) (*GitClient, error) {
	git, err := gitlab.NewClient(token, gitlab.WithBaseURL(baseUrl))
	if err != nil {
		return nil, err
	}

	return &GitClient{
		git:         git,
		projectChan: make(chan *ProjectB),
		timeChan:    make(chan *StartEndTime),
		wg:          &sync.WaitGroup{},
	}, nil
}

func (g *GitClient) GetProjects() {
	defer func() {
		close(g.projectChan)
		g.wg.Done()
	}()

	opt := gitlab.ListOptions{
		Page:    1,
		PerPage: 50,
	}

	for {
		projects, response, err := g.git.Projects.ListProjects(&gitlab.ListProjectsOptions{
			ListOptions: opt,
			Archived:    gitlab.Bool(false),
			//LastActivityAfter: gitlabcron.Time(lTime),
			Membership: gitlab.Bool(true),
			Owned:      gitlab.Bool(true),
			Sort:       gitlab.String("asc"),
			OrderBy:    gitlab.String("id"),
			Visibility: gitlab.Visibility("private"),
		})
		if err != nil {
			log.Fatal(err)
		}
		for _, project := range projects {
			for _, pj := range backend {
				if branch, ok := pj[project.ID]; ok {
					//fmt.Println(project.ID, project.Name, project.Description, project.Namespace.Name)
					g.projectChan <- &ProjectB{
						Project: project,
						branch:  branch,
					}
					break
				}
			}
			//fmt.Println(project.ID, project.Name, project.Description, project.Namespace.Name)
			//g.projectChan<- project
		}
		if response.NextPage == 0 {
			break
		}
		opt.Page = response.NextPage
	}
}

func (g *GitClient) getPipeline(project *ProjectB, opt *gitlab.ListProjectPipelinesOptions, options ...optFunc) {
	for _, fn := range options {
		if fn == nil {
			continue
		}
		if err := fn(opt); err != nil {
			log.Fatal(err)
		}
	}

	var l int

	for {
		pipelines, response, err := g.git.Pipelines.ListProjectPipelines(project.ID, opt)

		if err != nil {
			log.Fatal(err)
		}

		if len(pipelines) == 0 {
			break
		}

		l += len(pipelines)

		//for _, pipeline := range pipelines {
		//	fmt.Println(project.ID, project.Name, pipeline.UpdatedAt, pipeline.Status)
		//}

		if response.NextPage == 0 {
			break
		}

		opt.ListOptions.Page = response.NextPage
	}
	date := (*opt.UpdatedAfter).Add(16 * time.Hour).Local().Format("2006-01-02 15:04:00")
	if l != 0 {
		data := make(map[string]interface{})
		data["department"] = "理财师"
		data["job_name"] = project.Name
		data["count"] = l
		data["user"] = "王飞"
		data["date"] = date
		bytesData, _ := json.Marshal(data)
		req, err := http.NewRequest("POST", "http://test-tcops.yintech.net/dataapi/api/v1/deployAddData", bytes.NewReader(bytesData))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("token", "WYNQ8Sm7lhmSFNJnmw+JwhPF4kdzFtZPLiAFcmSM0ePy5D5MO5jaL5I1E0bYh7wl")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		//resp, err := http.Post("http://test-tcops.yintech.net/dataapi/api/v1/deployAddData", "application/json", bytes.NewReader(bytesData))
		//if err != nil {
		//	log.Fatal(err)
		//}
		fmt.Println(project.ID, project.Name, l, date, resp.Status)
	}

}

func (g *GitClient) getPipelineV2(project *ProjectB, opt *gitlab.ListProjectPipelinesOptions, options ...optFunc) {
	for _, fn := range options {
		if fn == nil {
			continue
		}
		if err := fn(opt); err != nil {
			log.Fatal(err)
		}
	}

	data := make(map[string]interface{})
	data["department"] = "理财师"
	data["job_name"] = project.Name
	data["user"] = "王飞"
	data["count"] = 1

	for {
		pipelines, response, err := g.git.Pipelines.ListProjectPipelines(project.ID, opt)

		if err != nil {
			log.Fatal(err)
		}

		if len(pipelines) == 0 {
			break
		}

		for _, pipe := range pipelines {
			data["date"] = pipe.UpdatedAt.Local().Format("2006-01-02 15:04:00")
			fmt.Printf("%#v, %s\n", pipe, data["date"])

			//data["date"] = pipe.UpdatedAt.Local().Format("2006-01-02 15:04:00")
			//bytesData, _ := json.Marshal(data)
			//req, err := http.NewRequest("POST", "http://test-tcops.yintech.net/dataapi/api/v1/deployAddData", bytes.NewReader(bytesData))
			//if err != nil {
			//	log.Fatal(err)
			//}
			//req.Header.Set("Content-Type", "application/json")
			//req.Header.Add("token", "WYNQ8Sm7lhmSFNJnmw+JwhPF4kdzFtZPLiAFcmSM0ePy5D5MO5jaL5I1E0bYh7wl")
			//
			//client := &http.Client{}
			//resp, err := client.Do(req)
			//if err != nil {
			//	log.Fatal(err)
			//}
			//fmt.Println(project.ID, project.Name, data["count"], data["date"], resp.Status)
		}

		if response.NextPage == 0 {
			break
		}

		opt.ListOptions.Page = response.NextPage
	}
}

func (g *GitClient) listPipelineByTime() {
	defer g.wg.Done()

	opt := &gitlab.ListProjectPipelinesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: 50,
		},
		Ref:     gitlab.String("master"),
		Status:  gitlab.BuildState("success"),
		OrderBy: gitlab.String("updated_at"),
		Sort:    gitlab.String("desc"),
	}

	for project := range g.projectChan {
		for _, startEndTime := range startEndTimes {
			startTime := startEndTime.startTime
			endTime := startEndTime.endTime
			g.getPipelineV2(project, opt, func(options *gitlab.ListProjectPipelinesOptions) error {
				options.UpdatedAfter = startTime
				options.UpdatedBefore = endTime
				options.Ref = gitlab.String(project.branch)
				return nil
			})
		}
	}
}

func (g *GitClient) listPipelineByTimeV2(starTime, endTime *time.Time) {
	defer g.wg.Done()

	opt := &gitlab.ListProjectPipelinesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: 50,
		},
		Ref:     gitlab.String("master"),
		Status:  gitlab.BuildState("success"),
		OrderBy: gitlab.String("updated_at"),
		Sort:    gitlab.String("desc"),
	}

	for project := range g.projectChan {
		g.getPipelineV2(project, opt, func(options *gitlab.ListProjectPipelinesOptions) error {
			options.UpdatedAfter = starTime
			options.UpdatedBefore = endTime
			options.Ref = gitlab.String(project.branch)
			return nil
		})
	}
}

func GetTimer(startTime time.Time, endTime time.Time) []StartEndTime {
	dur := endTime.Sub(startTime).Hours() / 24
	var startEndTimes = make([]StartEndTime, 0, int(dur+1))
	var (
		start = startTime
		end   time.Time
		i     int
	)
	for {
		end = start.Add(time.Hour * 24)
		if end.After(endTime) {
			break
		}
		startEndTimes = append(startEndTimes, StartEndTime{
			startTime: gitlab.Time(start),
			endTime:   gitlab.Time(end),
		})
		start = end
		i++
	}
	return startEndTimes
}
