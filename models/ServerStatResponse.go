package models

type ServerStatResponse struct {
	VMs []*VMInfo `json:"vms"`
}

type LoggedServerStatsResponse struct {
	Name    string    `json:"name"`
	Serial  int       `json:"serial"`
	VMStats []*VMStat `json:"stats"`
}
