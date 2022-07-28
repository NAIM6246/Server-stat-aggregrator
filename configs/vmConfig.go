package configs

type VMConfig struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	IP   string `json:"ip"`
	Host string `json:"host"`
	Port int    `json:"port"`
}
