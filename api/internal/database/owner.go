package database

import (
	"log"
	"strings"

	"github.com/BreD1810/tea-selector/api/internal/models"
)

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

// GetAllOwnersFromDatabase gets all the owners from the database.
func GetAllOwnersFromDatabase() ([]models.Owner, error) {
	rows, err := DB.Query("SELECT * FROM owner;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	owners := make([]models.Owner, 0)
	for rows.Next() {
		owner := new(models.Owner)
		err := rows.Scan(&owner.ID, &owner.Name)
		if err != nil {
			return nil, err
		}
		owners = append(owners, *owner)
	}
	return owners, nil
}

// GetOwnerFromDatabase gets an owner from the database by their ID.
func GetOwnerFromDatabase(owner *models.Owner) error {
	row := DB.QueryRow("SELECT name FROM owner WHERE id=$1;", owner.ID)

	err := row.Scan(&owner.Name)
	if err != nil {
		return err
	}

	return nil
}

// CreateOwnerInDatabase adds a new owner to the database
func CreateOwnerInDatabase(owner *models.Owner) error {
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
func DeleteOwnerFromDatabase(owner *models.Owner) error {
	rows, err := DB.Query("SELECT name FROM owner WHERE id=$1;", owner.ID)
	if err != nil {
		return err
	}

	rows.Next()
	err = rows.Scan(&owner.Name)
	if err != nil {
		return err
	}
	rows.Close()

	_, err = DB.Exec("DELETE FROM owner WHERE id = $1;", owner.ID)
	return err
}

// GetAllOwnersTeasFromDatabase gets all teas for each owner.
func GetAllOwnersTeasFromDatabase() ([]models.OwnerWithTeas, error) {
	rows, err := DB.Query("SELECT * FROM owner;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ownersWithTeas := make([]models.OwnerWithTeas, 0)
	for rows.Next() {
		ownerWithTeas := new(models.OwnerWithTeas)
		err := rows.Scan(&ownerWithTeas.Owner.ID, &ownerWithTeas.Owner.Name)
		if err != nil {
			return nil, err
		}

		ownersWithTeas = append(ownersWithTeas, *ownerWithTeas)
	}

	for i := range ownersWithTeas {
		teaRows, err := DB.Query("SELECT tea.id, tea.name, types.id, types.name FROM teaOwners INNER JOIN tea ON teaOwners.teaID = tea.id INNER JOIN types ON types.id = tea.teaType WHERE teaOwners.ownerID = $1;", ownersWithTeas[i].Owner.ID)
		if err != nil {
			return nil, err
		}

		for teaRows.Next() {
			tea := new(models.Tea)
			err := teaRows.Scan(&tea.ID, &tea.Name, &tea.TeaType.ID, &tea.TeaType.Name)
			if err != nil {
				return nil, err
			}

			ownersWithTeas[i].Teas = append(ownersWithTeas[i].Teas, *tea)
		}
	}

	return ownersWithTeas, nil
}
