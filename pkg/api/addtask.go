package api

import (
	"encoding/json"
	"final_project/pkg/db"
	"fmt"
	"net/http"
	"time"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, fmt.Sprintf("JSON is wrong: %s", err.Error()))
		return
	}
	if task.Title == "" {
		writeError(w, "Empty title")
		return
	}
	if err := checkDate(&task); err != nil {
		writeError(w, err.Error())
		return
	}
	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}

func checkDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}
	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("date have incorrect format, must be in %s; %w", dateFormat, err)
	}
	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("incorrect repeat rule: %w", err)
		}
	}
	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = now.Format(dateFormat)
		} else {
			task.Date = next
		}
	}
	return nil
}
