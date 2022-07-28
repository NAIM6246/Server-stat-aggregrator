package models

type ServerStatResponse struct {
	VMs []*VMInfo `json:"vms"`
}
