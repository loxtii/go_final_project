package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s Storage) InitDatabase() {
	querySQL := `
		CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY,
		date VARCHAR(8) NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat VARCHAR(128));
	    CREATE INDEX IF NOT EXISTS indexdate ON scheduler (date);`

	log.Println("[INFO] Creating new table...")
	_, err := s.db.Exec(querySQL)
	checkError(err, "table")

}

func checkError(err error, s string) {
	if err != nil {
		log.Println("[Error] Failed: " + s)
		log.Fatal(err)
	}
	log.Println("[Info] Success: " + s)
}

func (s Storage) InsertTask(task Task) (int, error) {
	querySQL := `INSERT INTO scheduler (date, title, comment, repeat) 
	             VALUES (:date, :title, :comment, :repeat)`
	res, err := s.db.Exec(querySQL,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	if err != nil {
		return 0, fmt.Errorf("add query error: %w", err)
	}
	log.Println("[Info] Success: add query executed ")

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("last insertion id error: %w", err)
	}

	return int(id), nil
}

func (s Storage) SelectTasks() ([]Task, error) {
	// log.Println("[WE IN STORAGE] ")
	var tasks []Task
	querySQL := `SELECT id, date, title, comment, repeat 
				 FROM scheduler 
				 ORDER BY date ASC
				 LIMIT 10`
	rows, err := s.db.Query(querySQL)

	// err := s.db.Select(&tasks, querySQL, time.Now().Format("20060102"))  needed   github.com/jmoiron/sqlx
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	// log.Println("[Query Select OK] ")
	defer rows.Close()

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		log.Println("[Scan OK] ")
		if err != nil {
			return nil, err
		}

		log.Println("[DATE] " + task.Date)

		tasks = append(tasks, task)
	}

	return tasks, nil

}
