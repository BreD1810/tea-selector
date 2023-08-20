package database

import (
	"errors"
	"strings"

	"github.com/BreD1810/tea-selector/api/internal/models"
)

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

// GetTeaOwnersFromDatabase gets all owners of a tea using the tea's ID.
func GetTeaOwnersFromDatabase(tea *models.Tea) ([]models.Owner, error) {
	rows, err := DB.Query("SELECT owner.id, owner.name FROM teaOwners INNER JOIN owner ON teaOwners.ownerID = owner.id WHERE teaOwners.teaID = $1;", tea.ID)
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

// GetAllTeaOwnersFromDatabase gets all owners for all teas.
func GetAllTeaOwnersFromDatabase() ([]models.TeaWithOwners, error) {
	teaRows, err := DB.Query("SELECT tea.id, tea.name, tea.teaType, types.name FROM tea INNER JOIN types on types.id = tea.teaType;")
	if err != nil {
		return nil, err
	}
	defer teaRows.Close()

	teas := make([]models.Tea, 0)
	for teaRows.Next() {
		tea := new(models.Tea)

		err := teaRows.Scan(&tea.ID, &tea.Name, &tea.TeaType.ID, &tea.TeaType.Name)
		if err != nil {
			return nil, err
		}

		teas = append(teas, *tea)
	}

	teasWithOwners := make([]models.TeaWithOwners, 0)
	for _, tea := range teas {
		owners, err := GetTeaOwnersFromDatabase(&tea)
		if err != nil {
			return nil, err
		}

		teasWithOwners = append(teasWithOwners, models.TeaWithOwners{Tea: tea, Owners: owners})
	}

	return teasWithOwners, nil
}

// CreateTeaOwnerInDatabase adds an owner to a tea in the database.
func CreateTeaOwnerInDatabase(teaID int, owner *models.Owner) (models.Tea, error) {
	tea := new(models.Tea)

	_, err := DB.Exec("INSERT INTO teaOwners VALUES ($1, $2);", teaID, owner.ID)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return *tea, errors.New("This relationship already exists")
		}
		if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			return *tea, errors.New("Either the tea or owner ID do not exist in the database")
		}
		return *tea, err
	}

	row := DB.QueryRow("SELECT tea.id, tea.name, types.id, types.name FROM tea INNER JOIN types ON tea.teaType = types.id WHERE tea.id = $1;", teaID)
	err = row.Scan(&tea.ID, &tea.Name, &tea.TeaType.ID, &tea.TeaType.Name)
	if err != nil {
		return *tea, errors.New("Tea ID not found after insert")
	}

	return *tea, nil
}

// DeleteTeaOwnerFromDatabase deletes an owner of a tea from the database.
func DeleteTeaOwnerFromDatabase(tea *models.Tea, owner *models.Owner) error {
	_, err := DB.Exec("DELETE FROM teaOwners WHERE teaID = $1 AND ownerID = $2;", tea.ID, owner.ID)
	return err
}
