package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type types struct {
	id   int
	name string
}

type tea struct {
	id      int
	name    string
	teaType int
}

type owner struct {
	id   int
	name string
}

type teaOwners struct {
	teaID   int
	ownerID int
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
		db.Exec("PRAGMA foreign_keys = ON;") // Enable foreign key checks
	}
	log.Println("Database initialised.")
}

func createDatabase(loc string) {
	database, err := sql.Open("sqlite3", loc)
	checkError("creating database", err)
	db = database
	db.SetMaxOpenConns(1)
	db.Exec("PRAGMA foreign_keys = ON;") // Enable foreign key checks

	createTeaTypeTable(db)
	createTeaTable(db)
	createOwnerTable(db)
	createTeaOwnersTable(db)
}

func createTeaTypeTable(db *sql.DB) {
	creationString := `CREATE TABLE types (
							id INTEGER PRIMARY KEY AUTOINCREMENT,
							name TEXT NOT NULL
					   );`
	res, err := db.Exec(creationString)
	checkError("creating types table", err)
	log.Println(res)
}

func createTeaTable(db *sql.DB) {
	creationString := `CREATE TABLE tea (
							id INTEGER PRIMARY KEY AUTOINCREMENT,
							name TEXT NOT NULL,
							teaType INTEGER,
							FOREIGN KEY (teaType) REFERENCES types (id)
								ON UPDATE CASCADE
								ON DELETE RESTRICT
					   );`
	res, err := db.Exec(creationString)
	checkError("creating tea table", err)
	log.Println(res)
}

func createOwnerTable(db *sql.DB) {
	creationString := `CREATE TABLE owner (
							id INTEGER PRIMARY KEY AUTOINCREMENT,
							name TEXT NOT NULL
					   );`
	res, err := db.Exec(creationString)
	checkError("creating owner table", err)
	log.Println(res)
}

func createTeaOwnersTable(db *sql.DB) {
	creationString := `CREATE TABLE teaOwners (
							teaID INTEGER,
							ownerID INTEGER,
							FOREIGN KEY (teaID) REFERENCES tea (id)
								ON UPDATE CASCADE
								ON DELETE RESTRICT,
							FOREIGN KEY (ownerID) REFERENCES owner (id)
								ON UPDATE CASCADE
								ON DELETE RESTRICT
					   );`
	res, err := db.Exec(creationString)
	checkError("creating owner table", err)
	log.Println(res)
}

func checkError(s string, e error) {
	if e != nil {
		log.Fatalf("Error "+s+": %v", e)
	}
}
