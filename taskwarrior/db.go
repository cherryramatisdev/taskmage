package taskwarrior

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type Task struct {
	Status      string    `json:"status"`
	Entry       time.Time `json:"entry"`
	Modified    time.Time `json:"modified"`
	Description string    `json:"description"`
	Due         time.Time `json:"due"`
	End         time.Time `json:"end"`
}

func FindTaskByDueDate(tasks []*Task, due time.Time) []*Task {
	var output []*Task

	for _, task := range tasks {
		if task.Due.Day() == due.Day() && task.Due.Month() == task.Due.Month() && task.Due.Year() == due.Year() && task.Due.Hour() == due.Hour() {
			output = append(output, task)
		}
	}

	if len(output) == 0 {
		return nil
	}

	return output
}

type taskFromDb struct {
	Status      string `json:"status"`
	Entry       string `json:"entry"`
	Modified    string `json:"modified"`
	Description string `json:"description"`
	Due         string `json:"due"`
	End         string `json:"end"`
}

func Connect() (*sql.DB, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not find user home dir %w", err)
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/.task/taskchampion.sqlite3", home))

	return db, nil
}

type Status string

const (
	Pending   Status = "pending"
	Completed Status = "completed"
	Deleted   Status = "deleted"
)

func GetTasksByStatus(db *sql.DB, status Status) ([]*Task, error) {
	rows, err := db.Query(fmt.Sprintf("select data from tasks where json_extract(data, '$.status') = '%s'", status))

	var tasks []*Task

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var data string

		err = rows.Scan(&data)

		if err != nil {
			return nil, err
		}

		var dbTask taskFromDb

		err = json.Unmarshal([]byte(data), &dbTask)

		if err != nil {
			return nil, err
		}

		entry, _ := strconv.ParseInt(dbTask.Entry, 10, 64)
		modified, _ := strconv.ParseInt(dbTask.Modified, 10, 64)
		due, _ := strconv.ParseInt(dbTask.Due, 10, 64)
		end, _ := strconv.ParseInt(dbTask.End, 10, 64)

		task := Task{
			Status:      dbTask.Status,
			Description: dbTask.Description,
			Entry:       time.Unix(entry, 0),
			Modified:    time.Unix(modified, 0),
			Due:         time.Unix(due, 0),
			End:         time.Unix(end, 0),
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}
