package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Task struct {
	ID      int64  `json:"id,string"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

const searchDateFormat = "02.01.2006"

func AddTask(task *Task) (int64, error) {
	if db == nil {
		return 0, fmt.Errorf("database must be initialized")
	}
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("not able to insert task: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("not able to get last insert id: %w", err)
	}
	return id, nil
}

func Tasks(limit int) ([]*Task, error) {
	if db == nil {
		return nil, fmt.Errorf("database must be initialized")
	}
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("not able to query tasks: %w", err)
	}
	defer rows.Close()
	return scanTasks(rows)
}

func TasksBySearch(search string, limit int) ([]*Task, error) {
	if db == nil {
		return nil, fmt.Errorf("database must be initialized")
	}
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("not able to query tasks: %w", err)
	}
	defer rows.Close()
	all, err := scanTasks(rows)
	if err != nil {
		return nil, err
	}
	needle := strings.ToLower(search)
	result := []*Task{}
	for _, t := range all {
		if strings.Contains(strings.ToLower(t.Title), needle) ||
			strings.Contains(strings.ToLower(t.Comment), needle) {
			result = append(result, t)
			if len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func TasksByDate(date string, limit int) ([]*Task, error) {
	if db == nil {
		return nil, fmt.Errorf("database must be initialized")
	}
	query := `SELECT id, date, title, comment, repeat FROM scheduler
	          WHERE date = ? ORDER BY date LIMIT ?`
	rows, err := db.Query(query, date, limit)
	if err != nil {
		return nil, fmt.Errorf("not able to query tasks: %w", err)
	}
	defer rows.Close()
	return scanTasks(rows)
}

func scanTasks(rows *sql.Rows) ([]*Task, error) {
	tasks := []*Task{}
	for rows.Next() {
		task := &Task{}
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, fmt.Errorf("not able to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return tasks, nil
}

func ParseSearchDate(search string) (string, bool) {
	t, err := time.Parse(searchDateFormat, search)
	if err != nil {
		return "", false
	}
	return t.Format("20060102"), true
}

func GetTask(id string) (*Task, error) {
	if db == nil {
		return nil, fmt.Errorf("database must be initialized")
	}
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	row := db.QueryRow(query, id)
	task := &Task{}
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cant find task")
		}
		return nil, fmt.Errorf("not able to get task: %w", err)
	}
	return task, nil
}

func UpdateTask(task *Task) error {
	if db == nil {
		return fmt.Errorf("database must be initialized")
	}
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("not able to update task: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("not able to get affected rows: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("cant find task")
	}
	return nil
}

func DeleteTask(id string) error {
	if db == nil {
		return fmt.Errorf("database must be initialized")
	}
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("not able to delete task: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("not able to get affected rows: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("cant find task")
	}
	return nil
}

func UpdateDate(next string, id string) error {
	if db == nil {
		return fmt.Errorf("database must be initialized")
	}
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := db.Exec(query, next, id)
	if err != nil {
		return fmt.Errorf("not able to update date: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("not able to get affected rows : %w", err)
	}
	if count == 0 {
		return fmt.Errorf("cant find task")
	}
	return nil
}
