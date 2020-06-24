package main

import (
	"database/sql"
	"errors"
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
							PRIMARY KEY(teaID, ownerID),
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
	_, err := DB.Exec("INSERT INTO types (name) VALUES ($1);", teaType.Name)
	if err != nil {
		return err
	}

	rows, err := DB.Query("SELECT ID FROM types WHERE name = ($1);", teaType.Name)
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
	_, err := DB.Exec("INSERT INTO owner (name) VALUES ($1);", owner.Name)
	if err != nil {
		return err
	}

	rows, err := DB.Query("SELECT ID FROM owner WHERE name = ($1);", owner.Name)
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

// DeleteOwnerFromDatabase deletes an owner from the database.
func DeleteOwnerFromDatabase(owner *Owner) error {
	rows, err := DB.Query("SELECT name FROM owner WHERE id=$1;", owner.ID)

	rows.Next()
	err = rows.Scan(&owner.Name)
	if err != nil {
		return err
	}
	rows.Close()

	_, err = DB.Exec("DELETE FROM owner WHERE id = $1;", owner.ID)
	return err
}

// GetAllTeasFromDatabase gets all the teas from the database.
func GetAllTeasFromDatabase() ([]Tea, error) {
	rows, err := DB.Query("SELECT tea.id, tea.name, types.id, types.name FROM tea INNER JOIN types ON types.ID = tea.teaType;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teas := make([]Tea, 0)
	for rows.Next() {
		tea := new(Tea)
		err := rows.Scan(&tea.ID, &tea.Name, &tea.TeaType.ID, &tea.TeaType.Name)
		if err != nil {
			return nil, err
		}

		teas = append(teas, *tea)
	}
	return teas, nil
}

// GetTeaFromDatabase gets information about a tea from the database using it's ID
func GetTeaFromDatabase(tea *Tea) error {
	row := DB.QueryRow("SELECT tea.name, types.id, types.name FROM tea INNER JOIN types ON tea.teaType=types.id WHERE tea.teaType=$1", tea.ID)

	err := row.Scan(&tea.Name, &tea.TeaType.ID, &tea.TeaType.Name)
	if err != nil {
		return err
	}

	return nil
}

// CreateTeaInDatabase creates a new tea in the database. Uses the type ID to do so.
func CreateTeaInDatabase(tea *Tea) error {
	row := DB.QueryRow("SELECT name FROM types WHERE id = $1;", tea.TeaType.ID)
	err := row.Scan(&tea.TeaType.Name)
	if err != nil {
		return errors.New("Tea type does not exist or is missing")
	}

	_, err = DB.Exec("INSERT INTO tea (name, teaType) VALUES ($1, $2);", tea.Name, tea.TeaType.ID)
	if err != nil {
		return err
	}

	row = DB.QueryRow("SELECT id FROM tea WHERE name = $1;", tea.Name)
	err = row.Scan(&tea.ID)
	if err != nil {
		return errors.New("Tea ID not found after insert")
	}

	return nil
}

// DeleteTeaFromDatabase deletes a tea from the database using it's ID.
func DeleteTeaFromDatabase(tea *Tea) error {
	rows, err := DB.Query("SELECT name FROM tea WHERE id=$1;", tea.ID)

	rows.Next()
	err = rows.Scan(&tea.Name)
	if err != nil {
		return err
	}
	rows.Close()

	_, err = DB.Exec("DELETE FROM tea WHERE id = $1;", tea.ID)
	return err
}

// GetTeaOwnersFromDatabase gets all owners of a tea using the tea's ID.
func GetTeaOwnersFromDatabase(tea *Tea) ([]Owner, error) {
	rows, err := DB.Query("SELECT owner.id, owner.name FROM teaOwners INNER JOIN owner ON teaOwners.ownerID = owner.id WHERE teaOwners.teaID = $1;", tea.ID)
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

// GetAllTeaOwnersFromDatabase gets all owners for all teas.
func GetAllTeaOwnersFromDatabase() ([]TeaWithOwners, error) {
	teaRows, err := DB.Query("SELECT tea.id, tea.name, tea.teaType, types.name FROM tea INNER JOIN types on types.id = tea.teaType;")
	if err != nil {
		return nil, err
	}
	defer teaRows.Close()

	teas := make([]Tea, 0)
	for teaRows.Next() {
		tea := new(Tea)

		err := teaRows.Scan(&tea.ID, &tea.Name, &tea.TeaType.ID, &tea.TeaType.Name)
		if err != nil {
			return nil, err
		}

		teas = append(teas, *tea)
	}

	teasWithOwners := make([]TeaWithOwners, 0)
	for _, tea := range teas {
		owners, err := GetTeaOwnersFromDatabase(&tea)
		if err != nil {
			return nil, err
		}

		teasWithOwners = append(teasWithOwners, TeaWithOwners{Tea: tea, Owners: owners})
	}

	return teasWithOwners, nil
}

// CreateTeaOwnerInDatabase adds an owner to a tea in the database.
func CreateTeaOwnerInDatabase(teaID int, owner *Owner) error {
	_, err := DB.Exec("INSERT INTO teaOwners VALUES ($1, $2);", teaID, owner.ID)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.New("This relationship already exists")
		}
		if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			return errors.New("Either the tea or owner ID do not exist in the database")
		}
		return err
	}

	return nil
}

// DeleteTeaOwnerFromDatabase deletes an owner of a tea from the database.
func DeleteTeaOwnerFromDatabase(tea *Tea, owner *Owner) error {
	_, err := DB.Exec("DELETE FROM teaOwners WHERE teaID = $1 AND ownerID = $2;", tea.ID, owner.ID)
	return err
}

// GetAllTypesTeasFromDatabase gets all teas by types.
func GetAllTypesTeasFromDatabase() ([]TypeWithTeas, error) {
	rows, err := DB.Query("SELECT * FROM types;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	typesWithTeas := make([]TypeWithTeas, 0)
	for rows.Next() {
		typeWithTeas := new(TypeWithTeas)
		err := rows.Scan(&typeWithTeas.Type.ID, &typeWithTeas.Type.Name)
		if err != nil {
			return nil, err
		}

		typesWithTeas = append(typesWithTeas, *typeWithTeas)
	}

	for i := range typesWithTeas {
		teaRows, err := DB.Query("SELECT tea.id, tea.name FROM tea WHERE tea.teaType = $1;", typesWithTeas[i].Type.ID)
		if err != nil {
			return nil, err
		}

		// typeWithTeas.Teas = make([]Tea, 0)
		for teaRows.Next() {
			tea := new(Tea)
			err := teaRows.Scan(&tea.ID, &tea.Name)
			if err != nil {
				return nil, err
			}
			tea.TeaType.ID = typesWithTeas[i].Type.ID
			tea.TeaType.Name = typesWithTeas[i].Type.Name

			typesWithTeas[i].Teas = append(typesWithTeas[i].Teas, *tea)
		}
	}

	return typesWithTeas, nil
}
