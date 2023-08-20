package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/BreD1810/tea-selector/api/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase(cfg config.Config) (*Database, error) {
	log.Println("Initialising the database...")

	database := &Database{}

	// Check if the database doesn't exists
	if _, err := os.Stat(cfg.Database.Location); os.IsNotExist(err) {
		log.Println("Database doesn't exist. Creating...")
		db, err := createDatabase(cfg)
		if err != nil {
			return nil, err
		}
		database = db
	} else {
		db, err := sql.Open("sqlite3", cfg.Database.Location)
		if err != nil {
			return nil, fmt.Errorf("error opening db: %w", err)
		}
		database.DB = db
	}

	database.DB.SetMaxOpenConns(1)
	if _, err := database.DB.Exec("PRAGMA foreign_keys = ON;"); err != nil { // Enable foreign key checks
		fmt.Errorf("error enabling foreign key checks: %w", err)
	}

	return database, nil
}

func createDatabase(cfg config.Config) (*Database, error) {
	db, err := sql.Open("sqlite3", cfg.Database.Location)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}
	database := &Database{DB: db}

	if err := database.createTeaTypeTable(cfg.Database.TeaTypes); err != nil {
		return nil, err
	}
	if err := database.createTeaTable(); err != nil {
		return nil, err
	}
	if err := database.createOwnerTable(cfg.Database.Owners); err != nil {
		return nil, err
	}
	if err := database.createTeaOwnersTable(); err != nil {
		return nil, err
	}
	if err := database.createUserTable(); err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}
