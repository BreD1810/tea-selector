package database

import (
	"log"
	"strings"

	"github.com/BreD1810/tea-selector/api/internal/models"
)

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

// GetAllTeaTypesFromDatabase retrieves all the tea types available in the database.
func GetAllTeaTypesFromDatabase() ([]models.TeaType, error) {
	rows, err := DB.Query("SELECT * FROM types;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teaTypes := make([]models.TeaType, 0)
	for rows.Next() {
		teaType := new(models.TeaType)
		err := rows.Scan(&teaType.ID, &teaType.Name)
		if err != nil {
			return nil, err
		}
		teaTypes = append(teaTypes, *teaType)
	}
	return teaTypes, nil
}

// GetTeaTypeFromDatabase retrieves a tea type from the database.
func GetTeaTypeFromDatabase(teaType *models.TeaType) error {
	row := DB.QueryRow("SELECT name FROM types WHERE id=$1;", teaType.ID)

	err := row.Scan(&teaType.Name)
	if err != nil {
		return err
	}

	return nil
}

// GetAllTypesTeasFromDatabase gets all teas by types.
func GetAllTypesTeasFromDatabase() ([]models.TypeWithTeas, error) {
	rows, err := DB.Query("SELECT * FROM types;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	typesWithTeas := make([]models.TypeWithTeas, 0)
	for rows.Next() {
		typeWithTeas := new(models.TypeWithTeas)
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

		for teaRows.Next() {
			tea := new(models.Tea)
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

// CreateTeaTypeInDatabase adds a new tea type to the database
func CreateTeaTypeInDatabase(teaType *models.TeaType) error {
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
func DeleteTeaTypeInDatabase(teaType *models.TeaType) error {
	rows, err := DB.Query("SELECT name FROM types WHERE id=$1;", teaType.ID)
	if err != nil {
		return err
	}

	rows.Next()
	err = rows.Scan(&teaType.Name)
	if err != nil {
		return err
	}
	rows.Close()

	_, err = DB.Exec("DELETE FROM types WHERE id = $1;", teaType.ID)
	return err
}
