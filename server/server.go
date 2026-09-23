package server

import (
	"final_project/pkg/api"
	"fmt"
	"net/http"
)

func StartServ() {
	webDir := http.FileServer(http.Dir("./web"))
	api.Init()
	http.Handle("/", webDir)
	fmt.Println("Запуск веб-сервера")
	err := http.ListenAndServe(":7540", nil)
	if err != nil {
		panic(err)
	}
}
