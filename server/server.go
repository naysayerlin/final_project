package server

import (
	"final_project/config"
	"final_project/pkg/api"
	"log"
	"net/http"
)

const defaultPort = "7540"
const webDir = "./web"

func StartServ(cfg config.Config) {
	api.Init(cfg)
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	log.Printf("Server running on port %s", cfg.Port)
	err := http.ListenAndServe(":"+cfg.Port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
