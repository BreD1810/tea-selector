package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

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

func initialiseDatabase(cfg Config) {
	log.Println("Initialising the database...")

	// Check if the database doesn't exists
	if _, err := os.Stat(cfg.Database.Location); os.IsNotExist(err) {
		log.Println("Database doesn't exist. Creating...")
		createDatabase(cfg)
		log.Println("Database created.")
	} else {
		database, err := sql.Open("sqlite3", cfg.Database.Location)
		checkError("opening database", err)
		db = database
		db.SetMaxOpenConns(1)
		db.Exec("PRAGMA foreign_keys = ON;") // Enable foreign key checks
	}
	log.Println("Database initialised.")
}

func createDatabase(cfg Config) {
	database, err := sql.Open("sqlite3", cfg.Database.Location)
	checkError("creating database", err)
	db = database
	db.SetMaxOpenConns(1)
	db.Exec("PRAGMA foreign_keys = ON;") // Enable foreign key checks

	createTeaTypeTable(db, cfg.Database.TeaTypes)
	createTeaTable(db)
	createOwnerTable(db, cfg.Database.Owners)
	createTeaOwnersTable(db)
}

func createTeaTypeTable(db *sql.DB, types []string) {
	creationString := `CREATE TABLE types (
							id INTEGER PRIMARY KEY AUTOINCREMENT,
							name TEXT NOT NULL
					   );`
	_, err := db.Exec(creationString)
	checkError("creating types table", err)

	if len(types) > 0 {
		var insertString strings.Builder
		insertString.WriteString("INSERT INTO types (name) VALUES ")

		for i, teaType := range types {
			if i != 0 {
				insertString.WriteString(", ")
			}
			insertString.WriteString("('" + teaType + "')")
		}

		insertString.WriteString(";")

		_, err = db.Exec(insertString.String())
		checkError("inserting types into the database", err)
	}
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
	_, err := db.Exec(creationString)
	checkError("creating tea table", err)
}

func createOwnerTable(db *sql.DB, owners []string) {
	creationString := `CREATE TABLE owner (
							id INTEGER PRIMARY KEY AUTOINCREMENT,
							name TEXT NOT NULL
					   );`
	_, err := db.Exec(creationString)
	checkError("creating owner table", err)

	if len(owners) > 0 {
		var insertString strings.Builder
		insertString.WriteString("INSERT INTO owner (name) VALUES ")

		for i, name := range owners {
			if i != 0 {
				insertString.WriteString(", ")
			}
			insertString.WriteString("('" + name + "')")
		}

		insertString.WriteString(";")

		_, err = db.Exec(insertString.String())
		checkError("inserting owners into the database", err)
	}
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
	_, err := db.Exec(creationString)
	checkError("creating owner table", err)
}

func checkError(s string, e error) {
	if e != nil {
		log.Fatalf("Error "+s+": %v", e)
	}
}
