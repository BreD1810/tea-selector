package database

import (
	"fmt"
	"strings"

	"github.com/BreD1810/tea-selector/api/internal/models"
)

func (db *Database) createOwnerTable(owners []string) error {
	creationString := `CREATE TABLE owner (
							id INTEGER PRIMARY KEY AUTOINCREMENT,
							name TEXT NOT NULL UNIQUE
					   );`
	if _, err := db.DB.Exec(creationString); err != nil {
		return fmt.Errorf("error creating owner table: %w", err)
	}

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

		if _, err := db.DB.Exec(insertString.String()); err != nil {
			return fmt.Errorf("error inserting owners in to the table: %w", err)
		}
	}
	return nil
}

// GetAllOwnersFromDatabase gets all the owners from the database.
func (db *Database) GetAllOwnersFromDatabase() ([]models.Owner, error) {
	rows, err := db.DB.Query("SELECT * FROM owner;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	owners := make([]models.Owner, 0)
	for rows.Next() {
		owner := new(models.Owner)
		if err := rows.Scan(&owner.ID, &owner.Name); err != nil {
			return nil, err
		}
		owners = append(owners, *owner)
	}
	return owners, nil
}

// GetOwnerFromDatabase gets an owner from the database by their ID.
func (db *Database) GetOwnerFromDatabase(owner *models.Owner) error {
	row := db.DB.QueryRow("SELECT name FROM owner WHERE id=$1;", owner.ID)

	err := row.Scan(&owner.Name)
	return err
}

// CreateOwnerInDatabase adds a new owner to the database
func (db *Database) CreateOwnerInDatabase(owner *models.Owner) error {
	if _, err := db.DB.Exec("INSERT INTO owner (name) VALUES ($1);", owner.Name); err != nil {
		return err
	}

	rows, err := db.DB.Query("SELECT ID FROM owner WHERE name = ($1);", owner.Name)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&owner.ID)
	return err
}

// DeleteOwnerFromDatabase deletes an owner from the database.
func (db *Database) DeleteOwnerFromDatabase(owner *models.Owner) error {
	rows, err := db.DB.Query("SELECT name FROM owner WHERE id=$1;", owner.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	if err := rows.Scan(&owner.Name); err != nil {
		return err
	}

	_, err = db.DB.Exec("DELETE FROM owner WHERE id = $1;", owner.ID)
	return err
}

// GetAllOwnersTeasFromDatabase gets all teas for each owner.
func (db *Database) GetAllOwnersTeasFromDatabase() ([]models.OwnerWithTeas, error) {
	rows, err := db.DB.Query("SELECT * FROM owner;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ownersWithTeas := make([]models.OwnerWithTeas, 0)
	for rows.Next() {
		ownerWithTeas := new(models.OwnerWithTeas)
		if err := rows.Scan(&ownerWithTeas.Owner.ID, &ownerWithTeas.Owner.Name); err != nil {
			return nil, err
		}

		ownersWithTeas = append(ownersWithTeas, *ownerWithTeas)
	}

	for i := range ownersWithTeas {
		teaRows, err := db.DB.Query("SELECT tea.id, tea.name, types.id, types.name FROM teaOwners INNER JOIN tea ON teaOwners.teaID = tea.id INNER JOIN types ON types.id = tea.teaType WHERE teaOwners.ownerID = $1;", ownersWithTeas[i].Owner.ID)
		if err != nil {
			return nil, err
		}

		for teaRows.Next() {
			tea := new(models.Tea)
			if err := teaRows.Scan(&tea.ID, &tea.Name, &tea.TeaType.ID, &tea.TeaType.Name); err != nil {
				return nil, err
			}

			ownersWithTeas[i].Teas = append(ownersWithTeas[i].Teas, *tea)
		}
	}

	return ownersWithTeas, nil
}
