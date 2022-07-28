package models

type VMInfo struct {
	Name   string  `json:"name"`
	VMStat *VMStat `json:"vmStat"`
}

type VMStat struct {
	Cpu          float64          `json:"cpu"`
	Memory       float64          `json:"memory"`
	Disk         *DiskInformation `json:"disk"`
	ResponseTime int64            `json:"responseTime"`
}

type DiskInformation struct {
	TotalSize       float64 `json:"totalSize"`
	TotalUsed       float64 `json:"totalUsed"`
	TotalAvailable  float64 `json:"totalAvailable"`
	UsagePercentage float64 `json:"usagePercentage"`
}
