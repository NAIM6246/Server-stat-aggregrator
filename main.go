package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/naim6246/aggregrator/configs"
	"github.com/naim6246/aggregrator/models"
)

func main() {
	config := configs.GetAppConfig()
	router := chi.NewRouter()

	router.Get("/server-stat", func(w http.ResponseWriter, r *http.Request) {

		response := &models.ServerStatResponse{}

		for _, vm := range config.VMs {
			vmInfo := &models.VMInfo{
				Name: vm.Name,
			}

			resp, err := http.Get(fmt.Sprintf("http://%s:%d/server-stat", vm.Host, vm.Port))
			if err != nil {
				log.Fatalln(err)
			}
			json.NewDecoder(resp.Body).Decode(&vmInfo.VMStat)
			response.VMs = append(response.VMs, vmInfo)
		}
		w.Header().Add("Cotent-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("serving on port: ", config.ListenPort)
	http.ListenAndServe(fmt.Sprintf(":%d", config.ListenPort), router)
}
