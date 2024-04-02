package main

import (
	"go_final_project/storage"
	"log"
	"net/http"
)

func main() {
	// initializing config
	//cfg := config.MustLoad()

	// initializing logger
	//log := setupLogger(cfg.Env)
	log.Println("[INFO] starting task-manager")

	// initializing sqlite storage
	// InitDataBase(cfg.StoragePath)
	storage.InitDatabase("./scheduler.db")

	// intializing handler for front-server
	http.Handle("/", http.FileServer(http.Dir("./web/")))
	// starting front-server
	log.Println("[INFO] Starting client on port 7540...")
	log.Fatal(http.ListenAndServe(":7540", nil))
}
