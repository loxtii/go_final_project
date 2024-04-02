package storage

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitDatabase(dataSource string) {
	log.Println("[INFO] Conecting database...")
	db, err := sql.Open("sqlite", dataSource)
	checkErr(err, "connection")
	defer db.Close()

	if databaseNotExists() {
		querySQL := `
		CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date DATE,
		title TEXT,
		comment TEXT,
		repeat VARCHAR(128));
	    CREATE INDEX IF NOT EXISTS indexdate ON scheduler (date);`

		log.Println("[INFO] Creating new database")
		_, err = db.Exec(querySQL)
		checkErr(err, "table")
	}
}

func databaseNotExists() bool {
	appPath, err := os.Executable()
	checkErr(err, "finding the current path")

	dbFile := filepath.Join(filepath.Dir(appPath), "./scheduler.db")
	_, err = os.Stat(dbFile)

	return err != nil
}

func checkErr(err error, s string) {
	if err != nil {
		log.Println("[Error] Failed: " + s)
		log.Fatal(err)
	}
	log.Println("[Info] Success: " + s)
}
