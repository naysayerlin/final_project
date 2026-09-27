package api

import (
	"final_project/pkg/db"
	"log"
	"net/http"
	"time"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		log.Printf("SERVER: DONE - ID is empty!")
		writeError(w, http.StatusInternalServerError, "id not set")
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		log.Printf("SERVER: DONE - Cannot get task! ID = %s.", id)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			log.Printf("SERVER: DONE - Cannot delete task! ID = %s.", id)
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		next, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			log.Printf("SERVER: DONE - Repeat rule is incorrect! ID = %s.", id)
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err := db.UpdateDate(next, id); err != nil {
			log.Printf("SERVER: DONE - Cannot update date! ID = %s.", id)
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	log.Printf("SERVER: DONE - Sucsessfull! ID = %s.", id)
	writeJSON(w, http.StatusOK, map[string]any{})
}
