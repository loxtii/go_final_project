package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitDatabase(dataSource string) {
	log.Println("[INFO] Connecting to database...")
	db, err := sql.Open("sqlite", dataSource)
	checkError(err, "connection")
	defer db.Close()

	if databaseNotExists() {
		querySQL := `
		CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date VARCHAR(8),
		title TEXT,
		comment TEXT,
		repeat VARCHAR(128));
	    CREATE INDEX IF NOT EXISTS indexdate ON scheduler (date);`

		log.Println("[INFO] Creating new table...")
		_, err = db.Exec(querySQL)
		checkError(err, "table")
	}
}

func databaseNotExists() bool {
	appPath, err := os.Executable()
	checkError(err, "finding the current path")

	dbFile := filepath.Join(filepath.Dir(appPath), "./scheduler.db")
	_, err = os.Stat(dbFile)

	return err != nil
}

func checkError(err error, s string) {
	if err != nil {
		log.Println("[Error] Failed: " + s)
		log.Fatal(err)
	}
	log.Println("[Info] Success: " + s)
}
