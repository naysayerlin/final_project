package main

import (
	"final_project/config"
	"final_project/pkg/db"
	"final_project/server"
	"log"
	"os"
)

func main() {
	cfg := config.Load()
	if err := db.Init(cfg.DBFile); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer db.Close()
	server.StartServ(cfg)
}
