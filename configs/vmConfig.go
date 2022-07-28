package configs

type VMConfig struct {
	Name   string `json:"name"`
	Serial int    `json:"serial"`
	IP     string `json:"ip"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
}
