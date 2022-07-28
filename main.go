package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/naim6246/Server-stat-aggregrator/auth"
	"github.com/naim6246/Server-stat-aggregrator/configs"
	"github.com/naim6246/Server-stat-aggregrator/models"
)

func main() {
	config := configs.GetAppConfig()
	router := chi.NewRouter()

	//cors
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	HandleRequest(router, config)

	fmt.Println("serving on port: ", config.ListenPort)
	http.ListenAndServe(fmt.Sprintf(":%d", config.ListenPort), router)
}

//handle req
func HandleRequest(router chi.Router, config *configs.AppConfig) {

	//login
	router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		user := &models.User{}
		parsingErr := json.NewDecoder(r.Body).Decode(&user)
		if parsingErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(parsingErr)
			return
		}
		if user.AccesID != config.User.AccesID && user.Password != config.User.Password {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("access denied")
			return
		}
		token, err := auth.GenerateToken(user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(parsingErr)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(token)
	})

	// get vms name
	router.With(auth.Authenticate).Get("/vm-list", func(w http.ResponseWriter, r *http.Request) {
		vmList := make([]*models.VMDto, 0)
		for _, vm := range config.VMs {
			vmList = append(vmList, &models.VMDto{
				Name:   vm.Name,
				Serial: vm.Serial,
			})
		}
		w.Header().Add("Cotent-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(vmList)
	})

	//get stats of vm
	router.With(auth.Authenticate).Get("/server-stats", func(w http.ResponseWriter, r *http.Request) {

		response := &models.ServerStatResponse{}

		for _, vm := range config.VMs {
			vmInfo := &models.VMInfo{
				Name:   vm.Name,
				Serial: vm.Serial,
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

}
