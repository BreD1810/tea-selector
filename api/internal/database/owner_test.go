package database

import (
	"database/sql"
	"testing"

	"github.com/BreD1810/tea-selector/api/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateOwnerTable(t *testing.T) {
	createOwnerString := "CREATE TABLE owner"
	owners := []string{"John", "Jane"}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v", err)
	}
	defer db.Close()
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	mock.ExpectExec(createOwnerString).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("INSERT INTO owner \\(name\\) VALUES \\('" + owners[0] + "'\\), \\('" + owners[1] + "'\\);").WillReturnResult(sqlmock.NewResult(2, 2))

	createOwnerTable(owners)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expections: %s", err)
	}
}

func TestCreateEmptyOwnerTable(t *testing.T) {
	createOwnerString := "CREATE TABLE owner"
	owners := []string{}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v", err)
	}
	defer db.Close()
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	mock.ExpectExec(createOwnerString).WillReturnResult(sqlmock.NewResult(0, 0))

	createOwnerTable(owners)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expections: %s", err)
	}
}

func TestGetAllOwnersFromDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v", err)
	}
	defer db.Close()
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	rows := mock.NewRows([]string{"id", "name"})
	rows.AddRow("1", "John")
	rows.AddRow("2", "Jane")

	mock.ExpectQuery("SELECT \\* FROM owner;").WillReturnRows(rows)

	owners, err := GetAllOwnersFromDatabase()
	if err != nil {
		t.Errorf("Database returned unexpected error: %v\n", err)
	}
	expected := models.Owner{1, "John"}
	if owners[0] != expected {
		t.Errorf("Database returned unexpected result:\n got: %v\n wanted: %v\n", owners[0], expected)
	}
	expected = models.Owner{2, "Jane"}
	if owners[1] != expected {
		t.Errorf("Database returned unexpected result:\n got: %q\n wanted: %q\n", owners[1], expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}

func TestGetOwnerFromDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n", err)
	}
	defer db.Close()
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	expected := "John"
	rows := mock.NewRows([]string{"name"})
	rows.AddRow(expected)
	owner := models.Owner{ID: 1}

	mock.ExpectQuery("SELECT name FROM owner").WithArgs(1).WillReturnRows(rows)

	err = GetOwnerFromDatabase(&owner)
	if err != nil {

	}
	if owner.Name != expected {
		t.Errorf("Database returned unexpected result:\n got: %q\n wanted: %q\n", owner.Name, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}

func TestGetNonExistentOwnerFromDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n", err)
	}
	defer db.Close()
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	expected := ""
	rows := mock.NewRows([]string{"name"})
	rows.AddRow(expected)
	owner := models.Owner{ID: 1}

	mock.ExpectQuery("SELECT name FROM owner").WithArgs(1).WillReturnError(sql.ErrNoRows)

	err = GetOwnerFromDatabase(&owner)
	if err != sql.ErrNoRows {
		t.Errorf("Method returned unexpected error:\n got: %v\n wanted: %v\n", err, sql.ErrNoRows)
	}
	if owner.Name != expected {
		t.Errorf("Database returned unexpected result:\n got: %q\n wanted: %q\n", owner.Name, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}

func TestCreateOwnerInDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n", err)
	}
	defer db.Close()
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	ownerName := "John"
	rows := mock.NewRows([]string{"id"})
	rows.AddRow("1")
	mock.ExpectExec("INSERT INTO owner").WithArgs(ownerName).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT ID FROM owner").WillReturnRows(rows)

	owner := models.Owner{ID: 1, Name: ownerName}
	err = CreateOwnerInDatabase(&owner)
	if err != nil {
		t.Errorf("Error whilst trying to insert tea type into database: %v\n", err)
	}
	if owner.ID != 1 {
		t.Errorf("Tea type ID not updated:\n Got: %d\n Expected: %v\n", owner.ID, 1)
	}
	if owner.Name != ownerName {
		t.Errorf("Tea type Name not as expected:\n Got: %q\n Expected: %q\n", owner.Name, ownerName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}

func TestDeleteOwnerFromDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n", err)
	}
	defer db.Close()
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	ownerName := "John"
	ownerID := 1
	rows := mock.NewRows([]string{"name"})
	rows.AddRow(ownerName)

	mock.ExpectQuery("SELECT name FROM owner").WithArgs(1).WillReturnRows(rows)
	mock.ExpectExec("DELETE FROM owner").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	owner := models.Owner{ID: ownerID}
	err = DeleteOwnerFromDatabase(&owner)
	if err != nil {
		t.Errorf("Error whilst trying to delete owner from database: %v\n", err)
	}
	if owner.ID != ownerID {
		t.Errorf("Owner ID changed:\n Got: %d\n Expected: %v\n", owner.ID, ownerID)
	}
	if owner.Name != ownerName {
		t.Errorf("Owner Name not as expected:\n Got: %q\n Expected: %q\n", owner.Name, ownerName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n\n", err)
	}
}

func TestDeleteNonExistantOwnerFromDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n\n", err)
	}
	defer db.Close()
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	// teaName := "Black Tea"
	ownerID := 1
	rows := mock.NewRows([]string{"name"})

	mock.ExpectQuery("SELECT name FROM owner").WithArgs(1).WillReturnRows(rows)

	owner := models.Owner{ID: ownerID}
	err = DeleteOwnerFromDatabase(&owner)
	if err.Error() != "sql: Rows are closed" {
		t.Errorf("Error whilst trying to delete owner from database: %v\n", err)
	}
	if owner.ID != ownerID {
		t.Errorf("Owner ID changed:\n Got: %d\n Expected: %v\n", owner.ID, ownerID)
	}
	if owner.Name != "" {
		t.Errorf("Owner Name was unexpectedly updated:\n Got: %q\n Expected: \"\"\n", owner.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}
