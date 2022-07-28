package configs

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/naim6246/Server-stat-aggregrator/models"
)

type AppConfig struct {
	ListenPort     int          `json:"listenPort"`
	Secret         string       `json:"secret"`
	AccesInTTLHour int          `json:"accessTTLInHour"`
	VMs            []*VMConfig  `json:"vms"`
	User           *models.User `json:"user"`
}

const path = "/config.json"

var appConfig *AppConfig

func loadConfig() {
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filePath := fmt.Sprintf("%s%s", workingDir, path)
	stream, err := os.Open(filePath)
	if err != nil {

		panic(err)
	}
	parseErr := json.NewDecoder(stream).Decode(&appConfig)
	if parseErr != nil {
		panic(parseErr)
	}
}

//LoadConfig from json file
func init() {
	var loadOnce sync.Once
	loadOnce.Do(loadConfig)
	// return appConfig
}

func GetAppConfig() *AppConfig {
	return appConfig
}
