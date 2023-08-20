package database

import (
	"fmt"
	"strings"

	"github.com/BreD1810/tea-selector/api/internal/models"
)

func (db *Database) createTeaTypeTable(types []string) error {
	creationString := `CREATE TABLE types (
							id INTEGER PRIMARY KEY AUTOINCREMENT,
							name TEXT NOT NULL UNIQUE
					   );`
	if _, err := db.DB.Exec(creationString); err != nil {
		return fmt.Errorf("error creating types table: %w", err)
	}

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

		if _, err := db.DB.Exec(insertString.String()); err != nil {
			return fmt.Errorf("error inserting types into the database", err)
		}
	}

	return nil
}

// GetAllTeaTypesFromDatabase retrieves all the tea types available in the database.
func (db *Database) GetAllTeaTypesFromDatabase() ([]models.TeaType, error) {
	rows, err := db.DB.Query("SELECT * FROM types;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teaTypes := make([]models.TeaType, 0)
	for rows.Next() {
		teaType := new(models.TeaType)
		if err := rows.Scan(&teaType.ID, &teaType.Name); err != nil {
			return nil, err
		}
		teaTypes = append(teaTypes, *teaType)
	}
	return teaTypes, nil
}

// GetTeaTypeFromDatabase retrieves a tea type from the database.
func (db *Database) GetTeaTypeFromDatabase(teaType *models.TeaType) error {
	row := db.DB.QueryRow("SELECT name FROM types WHERE id=$1;", teaType.ID)

	return row.Scan(&teaType.Name)
}

// GetAllTypesTeasFromDatabase gets all teas by types.
func (db *Database) GetAllTypesTeasFromDatabase() ([]models.TypeWithTeas, error) {
	rows, err := db.DB.Query("SELECT * FROM types;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	typesWithTeas := make([]models.TypeWithTeas, 0)
	for rows.Next() {
		typeWithTeas := new(models.TypeWithTeas)
		if err := rows.Scan(&typeWithTeas.Type.ID, &typeWithTeas.Type.Name); err != nil {
			return nil, err
		}

		typesWithTeas = append(typesWithTeas, *typeWithTeas)
	}

	for i := range typesWithTeas {
		teaRows, err := db.DB.Query("SELECT tea.id, tea.name FROM tea WHERE tea.teaType = $1;", typesWithTeas[i].Type.ID)
		if err != nil {
			return nil, err
		}

		for teaRows.Next() {
			tea := new(models.Tea)
			if err := teaRows.Scan(&tea.ID, &tea.Name); err != nil {
				return nil, err
			}
			tea.TeaType.ID = typesWithTeas[i].Type.ID
			tea.TeaType.Name = typesWithTeas[i].Type.Name

			typesWithTeas[i].Teas = append(typesWithTeas[i].Teas, *tea)
		}
	}

	return typesWithTeas, nil
}

// CreateTeaTypeInDatabase adds a new tea type to the database
func (db *Database) CreateTeaTypeInDatabase(teaType *models.TeaType) error {
	if _, err := db.DB.Exec("INSERT INTO types (name) VALUES ($1);", teaType.Name); err != nil {
		return err
	}

	rows, err := db.DB.Query("SELECT ID FROM types WHERE name = ($1);", teaType.Name)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	return rows.Scan(&teaType.ID)
}

// DeleteTeaTypeInDatabase deletes a tea type from the database.
func (db *Database) DeleteTeaTypeInDatabase(teaType *models.TeaType) error {
	rows, err := db.DB.Query("SELECT name FROM types WHERE id=$1;", teaType.ID)
	if err != nil {
		return err
	}

	rows.Next()
	if err = rows.Scan(&teaType.Name); err != nil {
		return err
	}
	rows.Close()

	_, err = db.DB.Exec("DELETE FROM types WHERE id = $1;", teaType.ID)
	return err
}
