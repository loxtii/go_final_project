package main

import (
	"log"
	"net/http"
)

func main() {

	log.Println("[INFO] starting task-manager")

	// initializing sqlite storage
	// InitDataBase(cfg.StoragePath)./
	InitDatabase("scheduler.db")

	// intializing handlers for web-server
	http.Handle("/", http.FileServer(http.Dir("./web/")))
	http.Handle("/api/nextdate", http.HandlerFunc(NextDate))

	// starting web-server
	log.Println("[INFO] Starting server on port 7540...")
	log.Fatal(http.ListenAndServe(":7540", nil))
	defer log.Println("[INFO] Server stopped")
}
