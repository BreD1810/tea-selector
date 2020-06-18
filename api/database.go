package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// DB is the database being used.
var DB *sql.DB

func checkError(s string, e error) {
	if e != nil {
		log.Fatalf("Error "+s+": %v", e)
	}
}

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
		DB = database
		DB.SetMaxOpenConns(1)
		DB.Exec("PRAGMA foreign_keys = ON;") // Enable foreign key checks
	}
	log.Println("Database initialised.")
}

func createDatabase(cfg Config) {
	database, err := sql.Open("sqlite3", cfg.Database.Location)
	checkError("creating database", err)
	DB = database
	DB.SetMaxOpenConns(1)
	DB.Exec("PRAGMA foreign_keys = ON;") // Enable foreign key checks

	createTeaTypeTable(cfg.Database.TeaTypes)
	createTeaTable()
	createOwnerTable(cfg.Database.Owners)
	createTeaOwnersTable()
}

func createTeaTypeTable(types []string) {
	creationString := `CREATE TABLE types (
							id INTEGER PRIMARY KEY AUTOINCREMENT,
							name TEXT NOT NULL UNIQUE
					   );`
	_, err := DB.Exec(creationString)
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

		_, err = DB.Exec(insertString.String())
		checkError("inserting types into the database", err)
	}
}

func createTeaTable() {
	creationString := `CREATE TABLE tea (
							id INTEGER PRIMARY KEY AUTOINCREMENT,
							name TEXT NOT NULL UNIQUE,
							teaType INTEGER,
							FOREIGN KEY (teaType) REFERENCES types (id)
								ON UPDATE CASCADE
								ON DELETE RESTRICT
					   );`
	_, err := DB.Exec(creationString)
	checkError("creating tea table", err)
}

func createOwnerTable(owners []string) {
	creationString := `CREATE TABLE owner (
							id INTEGER PRIMARY KEY AUTOINCREMENT,
							name TEXT NOT NULL UNIQUE
					   );`
	_, err := DB.Exec(creationString)
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

		_, err = DB.Exec(insertString.String())
		checkError("inserting owners into the database", err)
	}
}

func createTeaOwnersTable() {
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
	_, err := DB.Exec(creationString)
	checkError("creating owner table", err)
}

// GetAllTeaTypesFromDatabase retrieves all the tea types available in the database.
func GetAllTeaTypesFromDatabase() ([]TeaType, error) {
	rows, err := DB.Query("SELECT * FROM types;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teaTypes := make([]TeaType, 0)
	for rows.Next() {
		teaType := new(TeaType)
		err := rows.Scan(&teaType.ID, &teaType.Name)
		if err != nil {
			return nil, err
		}
		teaTypes = append(teaTypes, *teaType)
	}
	return teaTypes, nil
}

// GetTeaTypeFromDatabase retrieves a tea type from the database.
func GetTeaTypeFromDatabase(teaType *TeaType) error {
	row := DB.QueryRow("SELECT name FROM types WHERE id=$1;", teaType.ID)

	err := row.Scan(&teaType.Name)
	if err != nil {
		return err
	}

	return nil
}

// CreateTeaTypeInDatabase adds a new tea type to the database
func CreateTeaTypeInDatabase(teaType *TeaType) error {
	_, err := DB.Exec("INSERT INTO types (name) VALUES ('$1');", teaType.Name)
	if err != nil {
		return err
	}

	rows, err := DB.Query("SELECT ID FROM types WHERE name = ('$1');", teaType.Name)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&teaType.ID)
	if err != nil {
		return err
	}
	log.Printf("New ID: %d\n", teaType.ID)

	return nil
}

// DeleteTeaTypeInDatabase deletes a tea type from the database.
func DeleteTeaTypeInDatabase(teaType *TeaType) error {
	rows, err := DB.Query("SELECT name FROM types WHERE id=$1;", teaType.ID)

	rows.Next()
	err = rows.Scan(&teaType.Name)
	if err != nil {
		return err
	}
	rows.Close()

	_, err = DB.Exec("DELETE FROM types WHERE id = $1;", teaType.ID)
	return err
}

// GetAllOwnersFromDatabase gets all the owners from the database.
func GetAllOwnersFromDatabase() ([]Owner, error) {
	rows, err := DB.Query("SELECT * FROM owner;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	owners := make([]Owner, 0)
	for rows.Next() {
		owner := new(Owner)
		err := rows.Scan(&owner.ID, &owner.Name)
		if err != nil {
			return nil, err
		}
		owners = append(owners, *owner)
	}
	return owners, nil
}

// GetOwnerFromDatabase gets an owner from the database by their ID.
func GetOwnerFromDatabase(owner *Owner) error {
	row := DB.QueryRow("SELECT name FROM owner WHERE id=$1;", owner.ID)

	err := row.Scan(&owner.Name)
	if err != nil {
		return err
	}

	return nil
}

// CreateOwnerInDatabase adds a new owner to the database
func CreateOwnerInDatabase(owner *Owner) error {
	_, err := DB.Exec("INSERT INTO owner (name) VALUES ('$1');", owner.Name)
	if err != nil {
		return err
	}

	rows, err := DB.Query("SELECT ID FROM owner WHERE name = ('$1');", owner.Name)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&owner.ID)
	if err != nil {
		return err
	}
	log.Printf("New ID: %d\n", owner.ID)

	return nil
}
