package database

import (
	"errors"
	"fmt"

	"github.com/BreD1810/tea-selector/api/internal/models"
)

func (db *Database) createTeaTable() error {
	creationString := `CREATE TABLE tea (
							id INTEGER PRIMARY KEY AUTOINCREMENT,
							name TEXT NOT NULL UNIQUE,
							teaType INTEGER,
							FOREIGN KEY (teaType) REFERENCES types (id)
								ON UPDATE CASCADE
								ON DELETE RESTRICT
					   );`
	if _, err := db.DB.Exec(creationString); err != nil {
		return fmt.Errorf("error creating tea table: %w", err)
	}

	return nil
}

type TeaStorer interface {
	GetAllTeasFromDatabase() ([]models.Tea, error)
	GetTeaFromDatabase(tea *models.Tea) error
	CreateTeaInDatabase(tea *models.Tea) error
	DeleteTeaFromDatabase(tea *models.Tea) error
}

// GetAllTeasFromDatabase gets all the teas from the database.
func (db *Database) GetAllTeasFromDatabase() ([]models.Tea, error) {
	rows, err := db.DB.Query("SELECT tea.id, tea.name, types.id, types.name FROM tea INNER JOIN types ON types.ID = tea.teaType;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teas := make([]models.Tea, 0)
	for rows.Next() {
		tea := new(models.Tea)
		if err := rows.Scan(&tea.ID, &tea.Name, &tea.TeaType.ID, &tea.TeaType.Name); err != nil {
			return nil, err
		}

		teas = append(teas, *tea)
	}
	return teas, nil
}

// GetTeaFromDatabase gets information about a tea from the database using it's ID
func (db *Database) GetTeaFromDatabase(tea *models.Tea) error {
	row := db.DB.QueryRow("SELECT tea.name, types.id, types.name FROM tea INNER JOIN types ON tea.teaType=types.id WHERE tea.teaType=$1", tea.ID)

	err := row.Scan(&tea.Name, &tea.TeaType.ID, &tea.TeaType.Name)
	return err
}

// CreateTeaInDatabase creates a new tea in the database. Uses the type ID to do so.
func (db *Database) CreateTeaInDatabase(tea *models.Tea) error {
	row := db.DB.QueryRow("SELECT name FROM types WHERE id = $1;", tea.TeaType.ID)
	if err := row.Scan(&tea.TeaType.Name); err != nil {
		return errors.New("Tea type does not exist or is missing")
	}

	if _, err := db.DB.Exec("INSERT INTO tea (name, teaType) VALUES ($1, $2);", tea.Name, tea.TeaType.ID); err != nil {
		return err
	}

	row = db.DB.QueryRow("SELECT id FROM tea WHERE name = $1;", tea.Name)
	if err := row.Scan(&tea.ID); err != nil {
		return errors.New("Tea ID not found after insert")
	}

	return nil
}

// DeleteTeaFromDatabase deletes a tea from the database using it's ID.
func (db *Database) DeleteTeaFromDatabase(tea *models.Tea) error {
	rows, err := db.DB.Query("SELECT name FROM tea WHERE id=$1;", tea.ID)
	if err != nil {
		return err
	}

	rows.Next()
	if err = rows.Scan(&tea.Name); err != nil {
		return err
	}
	rows.Close()

	_, err = db.DB.Exec("DELETE FROM tea WHERE id = $1;", tea.ID)
	return err
}
