package database

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/BreD1810/tea-selector/api/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateTeaOwnersTable(t *testing.T) {
	createTeaOwnersString := "CREATE TABLE teaOwners"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v", err)
	}
	defer db.Close()
	oldDB := Database
	defer func() { Database = oldDB }()
	Database = db

	mock.ExpectExec(createTeaOwnersString).WillReturnResult(sqlmock.NewResult(0, 0))

	createTeaOwnersTable()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expections: %s", err)
	}
}

func TestGetTeaOwnerFromDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n", err)
	}
	defer db.Close()
	oldDB := Database
	defer func() { Database = oldDB }()
	Database = db

	expectedTeaID := 1
	tea := models.Tea{ID: expectedTeaID}
	expectedOwners := []models.Owner{{1, "John"}, {2, "Jane"}}
	rows := mock.NewRows([]string{"id", "name"})
	rows.AddRow(expectedOwners[0].ID, expectedOwners[0].Name)
	rows.AddRow(expectedOwners[1].ID, expectedOwners[1].Name)

	mock.ExpectQuery("SELECT (.)+ FROM teaOwners").WithArgs(expectedTeaID).WillReturnRows(rows)

	owners, err := GetTeaOwnersFromDatabase(&tea)
	if err != nil {
		t.Errorf("Database returned unexpected error: %v\n", err)
	}
	if owners[0].ID != expectedOwners[0].ID {
		t.Errorf("Database returned unexpected result:\n got: %q\n wanted: %q\n", owners[0].ID, expectedOwners[0].ID)
	}
	if owners[1].ID != expectedOwners[1].ID {
		t.Errorf("Database returned unexpected result:\n got: %q\n wanted: %q\n", owners[1].ID, expectedOwners[1].ID)
	}
	if owners[0].Name != expectedOwners[0].Name {
		t.Errorf("Database returned unexpected result:\n got: %q\n wanted: %q\n", owners[0].Name, expectedOwners[0].Name)
	}
	if owners[1].Name != expectedOwners[1].Name {
		t.Errorf("Database returned unexpected result:\n got: %q\n wanted: %q\n", owners[1].Name, expectedOwners[1].Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}

func TestGetNonExistentTeaOwnerFromDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n", err)
	}
	defer db.Close()
	oldDB := Database
	defer func() { Database = oldDB }()
	Database = db

	expectedTeaID := 1
	tea := models.Tea{ID: expectedTeaID}

	mock.ExpectQuery("SELECT (.)+ FROM teaOwners").WithArgs(expectedTeaID).WillReturnError(sql.ErrNoRows)

	_, err = GetTeaOwnersFromDatabase(&tea)
	if err != sql.ErrNoRows {
		t.Errorf("Database returned unexpected error: %v\n", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}

func TestCreateTeaOwnerFromDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n", err)
	}
	defer db.Close()
	oldDB := Database
	defer func() { Database = oldDB }()
	Database = db

	teaID := 1
	owner := models.Owner{ID: 1}
	rows := mock.NewRows([]string{"id", "name", "typeID", "typeName"})
	rows.AddRow(teaID, "Snowball", 1, "Black Tea")

	mock.ExpectExec("INSERT INTO teaOwners").WithArgs(teaID, owner.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT *(.)+ FROM tea").WithArgs(teaID).WillReturnRows(rows)

	tea, err := CreateTeaOwnerInDatabase(teaID, &owner)
	if err != nil {
		t.Errorf("Database returned unexpected error: %v\n", err)
	}
	if tea.ID != 1 {
		t.Errorf("Database returned unexpected tea ID: %d\n", tea.ID)
	}
	if tea.Name != "Snowball" {
		t.Errorf("Database returned unexpected tea name: %q\n", tea.Name)
	}
	if tea.TeaType.ID != 1 {
		t.Errorf("Database returned unexpected tea type id: %q\n", tea.TeaType.ID)
	}
	if tea.TeaType.Name != "Black Tea" {
		t.Errorf("Database returned unexpected tea type name: %q\n", tea.TeaType.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}

func TestCreateTeaOwnerFromDatabaseRelationshipExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n", err)
	}
	defer db.Close()
	oldDB := Database
	defer func() { Database = oldDB }()
	Database = db

	teaID := 1
	owner := models.Owner{ID: 1}

	mock.ExpectExec("INSERT INTO teaOwners").WithArgs(teaID, owner.ID).WillReturnError(errors.New("UNIQUE constraint failed"))

	if _, err := CreateTeaOwnerInDatabase(teaID, &owner); err.Error() != "This relationship already exists" {
		t.Errorf("Database returned unexpected error: %v\n", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}

func TestCreateTeaOwnerFromDatabaseDoesntExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n", err)
	}
	defer db.Close()
	oldDB := Database
	defer func() { Database = oldDB }()
	Database = db

	teaID := 1
	owner := models.Owner{ID: 1}

	mock.ExpectExec("INSERT INTO teaOwners").WithArgs(teaID, owner.ID).WillReturnError(errors.New("FOREIGN KEY constraint failed"))

	if _, err := CreateTeaOwnerInDatabase(teaID, &owner); err.Error() != "Either the tea or owner ID do not exist in the database" {
		t.Errorf("Database returned unexpected error: %q\n", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}

func TestDeleteTeaOwnerFromDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n", err)
	}
	defer db.Close()
	oldDB := Database
	defer func() { Database = oldDB }()
	Database = db

	tea := models.Tea{ID: 1}
	owner := models.Owner{ID: 1}

	mock.ExpectExec("DELETE FROM tea").WithArgs(tea.ID, owner.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	err = DeleteTeaOwnerFromDatabase(&tea, &owner)
	if err != nil {
		t.Errorf("Unexpected error whilst trying to delete tea owner from database: %v\n", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n\n", err)
	}
}

func TestDeleteNonExistantTeaOwnerFromDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n\n", err)
	}
	defer db.Close()
	oldDB := Database
	defer func() { Database = oldDB }()
	Database = db

	tea := models.Tea{ID: 1}
	owner := models.Owner{ID: 1}

	mock.ExpectExec("DELETE FROM tea").WithArgs(tea.ID, owner.ID).WillReturnError(sql.ErrNoRows)

	err = DeleteTeaOwnerFromDatabase(&tea, &owner)
	if err != sql.ErrNoRows {
		t.Errorf("Unexpected error whilst trying to delete tea from database: %v\n", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}
