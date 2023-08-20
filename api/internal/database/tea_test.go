package database

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/BreD1810/tea-selector/api/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateTeaTable(t *testing.T) {
	createTeaString := "CREATE TABLE tea"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v", err)
	}
	defer db.Close()
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	mock.ExpectExec(createTeaString).WillReturnResult(sqlmock.NewResult(0, 0))

	createTeaTable()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expections: %s", err)
	}
}

func TestGetAllTeasFromDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v", err)
	}
	defer db.Close()
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	rows := mock.NewRows([]string{"id", "name", "id", "name"})
	rows.AddRow("1", "Snowball", "1", "Black Tea")
	rows.AddRow("2", "Nearly Nirvana", "2", "White Tea")

	mock.ExpectQuery("SELECT (.)+ FROM tea").WillReturnRows(rows)

	teas, err := GetAllTeasFromDatabase()
	if err != nil {
		t.Errorf("Database returned unexpected error: %v\n", err)
	}
	expected := models.Tea{1, "Snowball", models.TeaType{1, "Black Tea"}}
	if teas[0] != expected {
		t.Errorf("Database returned unexpected result:\n got: %v\n wanted: %v\n", teas[0], expected)
	}
	expected = models.Tea{2, "Nearly Nirvana", models.TeaType{2, "White Tea"}}
	if teas[1] != expected {
		t.Errorf("Database returned unexpected result:\n got: %q\n wanted: %q\n", teas[1], expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}

func TestGetTeaFromDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n", err)
	}
	defer db.Close()
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	expectedTeaID := 1
	expectedTeaName := "Snowball"
	expectedTypeID := 1
	expectedTypeName := "Black Tea"
	tea := models.Tea{ID: expectedTeaID}
	rows := mock.NewRows([]string{"name", "id", "name"})
	rows.AddRow(expectedTeaName, expectedTypeID, expectedTypeName)

	mock.ExpectQuery("SELECT (.)+ FROM tea").WithArgs(1).WillReturnRows(rows)

	err = GetTeaFromDatabase(&tea)
	if err != nil {
		t.Errorf("Database returned unexpected error: %v\n", err)
	}
	if tea.ID != expectedTeaID {
		t.Errorf("Database returned unexpected result:\n got: %q\n wanted: %q\n", tea.ID, expectedTeaID)
	}
	if tea.Name != expectedTeaName {
		t.Errorf("Database returned unexpected result:\n got: %q\n wanted: %q\n", tea.Name, expectedTeaName)
	}
	if tea.TeaType.ID != expectedTypeID {
		t.Errorf("Database returned unexpected result:\n got: %q\n wanted: %q\n", tea.TeaType.ID, expectedTypeID)
	}
	if tea.TeaType.Name != expectedTypeName {
		t.Errorf("Database returned unexpected result:\n got: %q\n wanted: %q\n", tea.TeaType.Name, expectedTypeName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}

func TestGetNonExistentTeaFromDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n", err)
	}
	defer db.Close()
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	mock.ExpectQuery("SELECT (.)+ FROM tea").WithArgs(10).WillReturnError(sql.ErrNoRows)

	expectedTeaID := 10
	tea := models.Tea{ID: expectedTeaID}
	err = GetTeaFromDatabase(&tea)
	if err != sql.ErrNoRows {
		t.Errorf("Method returned unexpected error:\n got: %v\n wanted: %v\n", err, sql.ErrNoRows)
	}
	if tea.ID != expectedTeaID {
		t.Errorf("Database returned unexpected result:\n got: %q\n wanted: %q\n", tea.ID, expectedTeaID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}

func TestCreateTeaInDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n", err)
	}
	defer db.Close()
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	teaID := 1
	teaName := "Snowball"
	typeID := 1
	typeName := "Black Tea"
	typeRows := mock.NewRows([]string{"name"})
	typeRows.AddRow(typeName)
	teaRows := mock.NewRows([]string{"id"})
	teaRows.AddRow("1")

	mock.ExpectQuery("SELECT name FROM types").WithArgs(typeID).WillReturnRows(typeRows)
	mock.ExpectExec("INSERT INTO tea").WithArgs(teaName, typeID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT id FROM tea").WillReturnRows(teaRows)

	tea := models.Tea{Name: teaName, TeaType: models.TeaType{ID: typeID}}
	err = CreateTeaInDatabase(&tea)
	if err != nil {
		t.Errorf("Error whilst trying to insert tea into database: %s\n", err)
	}
	if tea.ID != teaID {
		t.Errorf("Tea ID unexpectedly updated:\n Got: %d\n Expected: %d\n", tea.ID, teaID)
	}
	if tea.Name != teaName {
		t.Errorf("Tea Name not as expected:\n Got: %q\n Expected: %q\n", tea.Name, teaName)
	}
	if tea.TeaType.ID != typeID {
		t.Errorf("teaType ID unexpectedly updated:\n Got: %d\n Expected: %d\n", tea.TeaType.ID, typeID)
	}
	if tea.TeaType.Name != typeName {
		t.Errorf("teaName not as expected:\n Got: %s\n Expected: %s\n", tea.TeaType.Name, typeName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}

func TestCreateTeaInDatabaseBadType(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n", err)
	}
	defer db.Close()
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	teaName := "Snowball"
	teaTypeID := 1
	expectedError := "Tea type does not exist or is missing"

	mock.ExpectQuery("SELECT name FROM types").WithArgs(1).WillReturnError(sql.ErrNoRows)

	tea := models.Tea{Name: teaName, TeaType: models.TeaType{ID: teaTypeID}}
	err = CreateTeaInDatabase(&tea)
	if err.Error() != expectedError {
		t.Errorf("Wrong error returned:\n Got: %s\n Expected: %s\n", err, expectedError)
	}
	if tea.Name != teaName {
		t.Errorf("Tea name unexpectedly updated:\n Got: %q\n Expected: %q\n", tea.Name, teaName)
	}
	if tea.TeaType.ID != teaTypeID {
		t.Errorf("teaType ID unexpectedly updated:\n Got: %d\n Expected: %d\n", tea.TeaType.ID, teaTypeID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}

func TestCreateTeaInDatabaseInsertError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n", err)
	}
	defer db.Close()
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	teaName := "Snowball"
	teaTypeID := 1
	typeName := "Black Tea"
	typeRows := mock.NewRows([]string{"name"})
	typeRows.AddRow(typeName)
	expectedError := "UNIQUE constraint not met"

	mock.ExpectQuery("SELECT name FROM types").WithArgs(1).WillReturnRows(typeRows)
	mock.ExpectExec("INSERT INTO tea").WithArgs(teaName, teaTypeID).WillReturnError(errors.New(expectedError))

	tea := models.Tea{Name: teaName, TeaType: models.TeaType{ID: teaTypeID}}
	err = CreateTeaInDatabase(&tea)
	if err.Error() != expectedError {
		t.Errorf("Wrong error returned:\n Got: %s\n Expected: %s\n", err, expectedError)
	}
	if tea.Name != teaName {
		t.Errorf("Tea name unexpectedly updated:\n Got: %q\n Expected: %q\n", tea.Name, teaName)
	}
	if tea.TeaType.ID != teaTypeID {
		t.Errorf("teaType ID unexpectedly updated:\n Got: %d\n Expected: %d\n", tea.TeaType.ID, teaTypeID)
	}
	if tea.TeaType.Name != typeName {
		t.Errorf("teaTypeName not correct:\n Got: %q\n Expected: %q\n", tea.TeaType.Name, typeName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}

func TestDeleteTeaFromDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n", err)
	}
	defer db.Close()
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	teaName := "Snowball"
	teaID := 1
	rows := mock.NewRows([]string{"name"})
	rows.AddRow(teaName)

	mock.ExpectQuery("SELECT name FROM tea").WithArgs(teaID).WillReturnRows(rows)
	mock.ExpectExec("DELETE FROM tea").WithArgs(teaID).WillReturnResult(sqlmock.NewResult(1, 1))

	tea := models.Tea{ID: teaID}
	err = DeleteTeaFromDatabase(&tea)
	if err != nil {
		t.Errorf("Error whilst trying to delete tea from database: %v\n", err)
	}
	if tea.ID != teaID {
		t.Errorf("Tea ID changed:\n Got: %d\n Expected: %v\n", tea.ID, teaID)
	}
	if tea.Name != teaName {
		t.Errorf("Tea Name not as expected:\n Got: %q\n Expected: %q\n", tea.Name, teaName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n\n", err)
	}
}

func TestDeleteNonExistantTeaFromDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n\n", err)
	}
	defer db.Close()
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	teaID := 1
	rows := mock.NewRows([]string{"name"})

	mock.ExpectQuery("SELECT name FROM tea").WithArgs(teaID).WillReturnRows(rows)

	tea := models.Tea{ID: teaID}
	err = DeleteTeaFromDatabase(&tea)
	if err.Error() != "sql: Rows are closed" {
		t.Errorf("Error whilst trying to delete tea from database: %v\n", err)
	}
	if tea.ID != teaID {
		t.Errorf("Tea ID changed:\n Got: %d\n Expected: %v\n", tea.ID, teaID)
	}
	if tea.Name != "" {
		t.Errorf("Tea Name was unexpectedly updated:\n Got: %q\n Expected: \"\"\n", tea.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}
