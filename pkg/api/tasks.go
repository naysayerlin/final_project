package api

import (
	"net/http"

	"final_project/pkg/db"
)

const tasksLimit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	search := r.FormValue("search")

	tasks := make([]*db.Task, 0)
	var err error

	switch {
	case search == "":
		tasks, err = db.Tasks(tasksLimit)
	default:
		if date, ok := db.ParseSearchDate(search); ok {
			tasks, err = db.TasksByDate(date, tasksLimit)
		} else {
			tasks, err = db.TasksBySearch(search, tasksLimit)
		}
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, TasksResp{Tasks: tasks})
}
