package main

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

type tea struct {
	id   int
	name string
}

var db *sql.DB

func initialiseDatabase() {
	log.Println("Initialising the database.")
	database, err := sql.Open("sqlite3", "tea-store.db")
	checkError("opening database", err)
	db = database

	db.SetMaxOpenConns(1)

	setupTeaTable()
}

func setupTeaTable() {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS tea (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT);")
	checkError("creating tea table", err)

	res, err := stmt.Exec()
	checkError("creating tea table", err)
	log.Println(res)
}

func checkError(s string, e error) {
	if e != nil {
		log.Fatal("Error "+s+": %v", e)
	}
}
