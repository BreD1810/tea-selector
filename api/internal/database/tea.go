package database

import (
	"errors"

	"github.com/BreD1810/tea-selector/api/internal/models"
)

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

// GetAllTeasFromDatabase gets all the teas from the database.
func GetAllTeasFromDatabase() ([]models.Tea, error) {
	rows, err := DB.Query("SELECT tea.id, tea.name, types.id, types.name FROM tea INNER JOIN types ON types.ID = tea.teaType;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teas := make([]models.Tea, 0)
	for rows.Next() {
		tea := new(models.Tea)
		err := rows.Scan(&tea.ID, &tea.Name, &tea.TeaType.ID, &tea.TeaType.Name)
		if err != nil {
			return nil, err
		}

		teas = append(teas, *tea)
	}
	return teas, nil
}

// GetTeaFromDatabase gets information about a tea from the database using it's ID
func GetTeaFromDatabase(tea *models.Tea) error {
	row := DB.QueryRow("SELECT tea.name, types.id, types.name FROM tea INNER JOIN types ON tea.teaType=types.id WHERE tea.teaType=$1", tea.ID)

	err := row.Scan(&tea.Name, &tea.TeaType.ID, &tea.TeaType.Name)
	if err != nil {
		return err
	}

	return nil
}

// CreateTeaInDatabase creates a new tea in the database. Uses the type ID to do so.
func CreateTeaInDatabase(tea *models.Tea) error {
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
func DeleteTeaFromDatabase(tea *models.Tea) error {
	rows, err := DB.Query("SELECT name FROM tea WHERE id=$1;", tea.ID)
	if err != nil {
		return err
	}

	rows.Next()
	err = rows.Scan(&tea.Name)
	if err != nil {
		return err
	}
	rows.Close()

	_, err = DB.Exec("DELETE FROM tea WHERE id = $1;", tea.ID)
	return err
}
