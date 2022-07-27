package tpl

const (
	SyncSlsDataTpl    = "schedule:syncslsdata"
	DealLogMqTpl      = "schedule:deallogmq"
	DealGitCountMqTpl = "schedule:dealgitcount"
)

type SyncSlsPayload struct {
	Email   string
	Content string
}

type GitlabPayload struct {
	Email   string
	Content string
}
