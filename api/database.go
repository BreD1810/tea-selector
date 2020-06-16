package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type tea struct {
	id   int
	name string
}

var db *sql.DB

func initialiseDatabase(loc string) {
	log.Println("Initialising the database...")

	// Check if the database doesn't exists
	if _, err := os.Stat(loc); os.IsNotExist(err) {
		log.Println("Database doesn't exist. Creating...")
		createDatabase(loc)
		log.Println("Database created.")
	} else {
		database, err := sql.Open("sqlite3", loc)
		checkError("opening database", err)
		db = database
		db.SetMaxOpenConns(1)
	}
	log.Println("Database initialised.")
}

func createDatabase(loc string) {
	database, err := sql.Open("sqlite3", loc)
	checkError("creating database", err)
	db = database
	db.SetMaxOpenConns(1)

	createTeaTable(db)
}

func createTeaTable(db *sql.DB) {
	stmt, err := db.Prepare("CREATE TABLE tea (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT);")
	checkError("creating tea table", err)

	res, err := stmt.Exec()
	checkError("creating tea table", err)
	log.Println(res)
}

func checkError(s string, e error) {
	if e != nil {
		log.Fatalf("Error "+s+": %v", e)
	}
}
