package api

import (
	"net/http"
	"time"
)

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var now time.Time
	nowParameter := r.FormValue("now")
	dateParameter := r.FormValue("date")
	repeatParameter := r.FormValue("repeat")
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
