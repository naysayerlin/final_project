package api

import (
	"encoding/json"
	"final_project/pkg/db"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Println("SERVER: ADD - error decoding input JSON:", err)
		writeError(w, http.StatusBadRequest, fmt.Sprintf("JSON is wrong: %s", err.Error()))
		return
	}
	if task.Title == "" {
		log.Println("SERVER: ADD - Empty Title")
		writeError(w, http.StatusBadRequest, "Empty title")
		return
	}
	if err := checkDate(&task); err != nil {
		log.Println("SERVER: ADD - Date validation error:", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	id, err := db.AddTask(&task)
	if err != nil {
		log.Println("SERVER: ADD - Error adding to the database:", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("SERVER: ADD - Successfully added! ID=%d.", id)
	writeJSON(w, http.StatusOK, map[string]string{"id": fmt.Sprintf("%d", id)})
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

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		log.Println("SERVER: GET - Empty ID")
		writeError(w, http.StatusBadRequest, "Empty id")
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		log.Println("SERVER: GET - Error importing from the database:", err)
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	log.Printf("SERVER: GET - Method Get sucsessful! ID = %s.", id)
	writeJSON(w, http.StatusOK, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Println("SERVER: UPDATE - Error decoding input JSON:", err)
		writeError(w, http.StatusBadRequest, fmt.Sprintf("JSON is wrong: %s", err.Error()))
		return
	}
	if strconv.FormatInt(task.ID, 10) == "" {
		log.Printf("SERVER: UPDATE - Empty ID = %d.", task.ID)
		writeError(w, http.StatusBadRequest, "Empty id!")
		return
	}
	if task.Title == "" {
		log.Printf("SERVER: UPDATE - Empty Title! ID = %d.", task.ID)
		writeError(w, http.StatusBadRequest, "Empty title")
		return
	}
	if err := checkDate(&task); err != nil {
		log.Printf("SERVER: UPDATE - Wrong Date! ID = %d.", task.ID)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := db.UpdateTask(&task); err != nil {
		log.Printf("SERVER: UPDATE - failed to update task! ID = %d.", task.ID)
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	log.Printf("SERVER: UPDATE - Successfully updated! ID = %d.", task.ID)
	writeJSON(w, http.StatusOK, map[string]any{})

}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		log.Printf("SERVER: DELETE - Empty ID = %s.", id)
		writeError(w, http.StatusBadRequest, "Empty ID")
		return
	}
	if err := db.DeleteTask(id); err != nil {
		log.Printf("SERVER: DELETE - Failed to delete ID = %s.", id)
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	log.Printf("SERVER: DELETE - Successfully deleted! ID = %s.", id)
	writeJSON(w, http.StatusOK, map[string]any{})
}
