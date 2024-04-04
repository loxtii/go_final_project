package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

func main() {

	log.Println("[INFO] starting task-manager")

	// opening sqlite database
	// InitDataBase(cfg.StoragePath)./
	log.Println("[INFO] Connecting to database...")
	db, err := sql.Open("sqlite", "scheduler.db")
	checkError(err, "connection")
	defer db.Close()

	// initializing new storage
	s := NewStorage(db)
	// creating new table
	s.InitDatabase()
	// creating new task-service
	ts := NewTaskService(s)

	// intializing handlers for web-server
	http.Handle("/", http.FileServer(http.Dir("./web/")))
	http.Handle("/api/nextdate", http.HandlerFunc(NextDate))

	http.Handle("/api/task", http.HandlerFunc(ts.TaskHandler))       // POST GET PUT DELETE
	http.Handle("GET /api/tasks", http.HandlerFunc(ts.TasksHandler)) // GET
	//http.HandleFunc("/api/task/done", http.HandlerFunc(DoneHandler)) // POST    "GET /api/tasks"

	// starting web-server
	log.Println("[INFO] Starting server on port 7540...")
	log.Fatal(http.ListenAndServe(":7540", nil))
	defer log.Println("[INFO] Server stopped")
}
