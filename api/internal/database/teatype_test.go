package database

import (
	"database/sql"
	"testing"

	"github.com/BreD1810/tea-selector/api/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateTeaTypeTable(t *testing.T) {
	createTeaTypeString := "CREATE TABLE types"
	teaTypes := []string{"Black Tea", "Green Tea"}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v", err)
	}
	defer db.Close()
	oldDB := Database
	defer func() { Database = oldDB }()
	Database = db

	mock.ExpectExec(createTeaTypeString).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("INSERT INTO types \\(name\\) VALUES \\('" + teaTypes[0] + "'\\), \\('" + teaTypes[1] + "'\\);").WillReturnResult(sqlmock.NewResult(2, 2))

	createTeaTypeTable(teaTypes)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expections: %s", err)
	}
}

func TestCreateEmptyTeaTypeTable(t *testing.T) {
	createTeaTypeString := "CREATE TABLE types"
	teaTypes := []string{}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v", err)
	}
	defer db.Close()
	oldDB := Database
	defer func() { Database = oldDB }()
	Database = db

	mock.ExpectExec(createTeaTypeString).WillReturnResult(sqlmock.NewResult(0, 0))

	createTeaTypeTable(teaTypes)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expections: %s", err)
	}
}

func TestGetAllTeaTypesFromDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v", err)
	}
	defer db.Close()
	oldDB := Database
	defer func() { Database = oldDB }()
	Database = db

	rows := mock.NewRows([]string{"id", "name"})
	rows.AddRow("1", "Black Tea")
	rows.AddRow("2", "Green Tea")

	mock.ExpectQuery("SELECT \\* FROM types;").WillReturnRows(rows)

	teaTypes, err := GetAllTeaTypesFromDatabase()
	if err != nil {
		t.Errorf("Database returned unexpected error: %v\n", err)
	}
	expected := models.TeaType{1, "Black Tea"}
	if teaTypes[0] != expected {
		t.Errorf("Database returned unexpected result:\n got: %v\n wanted: %v\n", teaTypes[0], expected)
	}
	expected = models.TeaType{2, "Green Tea"}
	if teaTypes[1] != expected {
		t.Errorf("Database returned unexpected result:\n got: %q\n wanted: %q\n", teaTypes[1], expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}

func TestGetTeaTypeFromDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n", err)
	}
	defer db.Close()
	oldDB := Database
	defer func() { Database = oldDB }()
	Database = db

	expected := "Black Tea"
	rows := mock.NewRows([]string{"name"})
	rows.AddRow(expected)
	teaType := models.TeaType{ID: 1}

	mock.ExpectQuery("SELECT name FROM types").WithArgs(1).WillReturnRows(rows)

	err = GetTeaTypeFromDatabase(&teaType)
	if err != nil {

	}
	if teaType.Name != expected {
		t.Errorf("Database returned unexpected result:\n got: %q\n wanted: %q\n", teaType.Name, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}

func TestGetNonExistantTeaTypeFromDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n", err)
	}
	defer db.Close()
	oldDB := Database
	defer func() { Database = oldDB }()
	Database = db

	expected := ""
	rows := mock.NewRows([]string{"name"})
	rows.AddRow(expected)
	teaType := models.TeaType{ID: 1}

	mock.ExpectQuery("SELECT name FROM types").WithArgs(1).WillReturnError(sql.ErrNoRows)

	err = GetTeaTypeFromDatabase(&teaType)
	if err != sql.ErrNoRows {
		t.Errorf("Method returned unexpected error:\n got: %v\n wanted: %v\n", err, sql.ErrNoRows)
	}
	if teaType.Name != expected {
		t.Errorf("Database returned unexpected result:\n got: %q\n wanted: %q\n", teaType.Name, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}

func TestCreateTeaTypeInDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n", err)
	}
	defer db.Close()
	oldDB := Database
	defer func() { Database = oldDB }()
	Database = db

	teaName := "Black Tea"
	rows := mock.NewRows([]string{"id"})
	rows.AddRow("1")
	mock.ExpectExec("INSERT INTO types").WithArgs("Black Tea").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT ID FROM types").WillReturnRows(rows)

	teaType := models.TeaType{ID: 1, Name: teaName}
	err = CreateTeaTypeInDatabase(&teaType)
	if err != nil {
		t.Errorf("Error whilst trying to insert tea type into database: %v\n", err)
	}
	if teaType.ID != 1 {
		t.Errorf("Tea type ID not updated:\n Got: %d\n Expected: %v\n", teaType.ID, 1)
	}
	if teaType.Name != teaName {
		t.Errorf("Tea type Name not as expected:\n Got: %q\n Expected: %q\n", teaType.Name, teaName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}

func TestDeleteTeaTypeInDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n", err)
	}
	defer db.Close()
	oldDB := Database
	defer func() { Database = oldDB }()
	Database = db

	teaName := "Black Tea"
	teaID := 1
	rows := mock.NewRows([]string{"name"})
	rows.AddRow(teaName)

	mock.ExpectQuery("SELECT name FROM types").WithArgs(1).WillReturnRows(rows)
	mock.ExpectExec("DELETE FROM types").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	teaType := models.TeaType{ID: teaID}
	err = DeleteTeaTypeInDatabase(&teaType)
	if err != nil {
		t.Errorf("Error whilst trying to delete tea type from database: %v\n", err)
	}
	if teaType.ID != teaID {
		t.Errorf("Tea type ID changed:\n Got: %d\n Expected: %v\n", teaType.ID, teaID)
	}
	if teaType.Name != teaName {
		t.Errorf("Tea type Name not as expected:\n Got: %q\n Expected: %q\n", teaType.Name, teaName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n\n", err)
	}
}

func TestDeleteNonExistantTeaTypeInDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n\n", err)
	}
	defer db.Close()
	oldDB := Database
	defer func() { Database = oldDB }()
	Database = db

	// teaName := "Black Tea"
	teaID := 1
	rows := mock.NewRows([]string{"name"})

	mock.ExpectQuery("SELECT name FROM types").WithArgs(1).WillReturnRows(rows)

	teaType := models.TeaType{ID: teaID}
	err = DeleteTeaTypeInDatabase(&teaType)
	if err.Error() != "sql: Rows are closed" {
		t.Errorf("Error whilst trying to delete tea type from database: %v\n", err)
	}
	if teaType.ID != teaID {
		t.Errorf("Tea type ID changed:\n Got: %d\n Expected: %v\n", teaType.ID, teaID)
	}
	if teaType.Name != "" {
		t.Errorf("Tea type Name was unexpectedly updated:\n Got: %q\n Expected: \"\"\n", teaType.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}
