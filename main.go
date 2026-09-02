package main

import (
	"final_project/pkg/db"
	"final_project/server"
	"fmt"
)

func main() {

	if err := db.Init("scheduler.db"); err != nil {
		fmt.Println(err)
		return
	}
	server.StartServ()
}
