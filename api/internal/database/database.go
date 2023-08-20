package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/BreD1810/tea-selector/api/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

// DB is the database being used.
var DB *sql.DB

func checkError(s string, e error) {
	if e != nil {
		log.Fatalf("Error "+s+": %v", e)
	}
}

func InitialiseDatabase(cfg config.Config) {
	log.Println("Initialising the database...")

	// Check if the database doesn't exists
	if _, err := os.Stat(cfg.Database.Location); os.IsNotExist(err) {
		log.Println("Database doesn't exist. Creating...")
		createDatabase(cfg)
		log.Println("Database created.")
	} else {
		database, err := sql.Open("sqlite3", cfg.Database.Location)
		checkError("opening database", err)
		DB = database
		DB.SetMaxOpenConns(1)
		_, _ = DB.Exec("PRAGMA foreign_keys = ON;") // Enable foreign key checks
	}
	log.Println("Database initialised.")
}

func createDatabase(cfg config.Config) {
	database, err := sql.Open("sqlite3", cfg.Database.Location)
	checkError("creating database", err)
	DB = database
	DB.SetMaxOpenConns(1)
	_, _ = DB.Exec("PRAGMA foreign_keys = ON;") // Enable foreign key checks

	createTeaTypeTable(cfg.Database.TeaTypes)
	createTeaTable()
	createOwnerTable(cfg.Database.Owners)
	createTeaOwnersTable()
	createUserTable()
}
