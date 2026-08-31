package server

import (
	"fmt"
	"net/http"
)

func StartServ() {
	webDir := http.FileServer(http.Dir("./web"))
	http.Handle("/", webDir)
	fmt.Println("Запуск веб-сервера")
	err := http.ListenAndServe(":7540", nil)
	if err != nil {
		panic(err)
	}
}
