package api

import (
	"final_project/config"
	"net/http"
	"time"
)

var appConfig config.Config

func Init(cfg config.Config) {
	appConfig = cfg
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", doneTaskHandler)
	http.HandleFunc("/api/signin", signInHandler)
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var now time.Time
	nowParameter := r.FormValue("now")
	dateParameter := r.FormValue("date")
	repeatParameter := r.FormValue("repeat")
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if nowParameter == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateFormat, nowParameter)
		if err != nil {
			http.Error(w, "parameter 'now' is invalid!", http.StatusBadRequest)
			return
		}
	}
	nextDate, err := NextDate(now, dateParameter, repeatParameter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(nextDate))
}
