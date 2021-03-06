package main

import (
	"database/sql"
	"errors"
	"testing"

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
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

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
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	mock.ExpectExec(createTeaTypeString).WillReturnResult(sqlmock.NewResult(0, 0))

	createTeaTypeTable(teaTypes)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expections: %s", err)
	}
}

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

func TestCreateTeaOwnersTable(t *testing.T) {
	createTeaOwnersString := "CREATE TABLE teaOwners"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v", err)
	}
	defer db.Close()
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	mock.ExpectExec(createTeaOwnersString).WillReturnResult(sqlmock.NewResult(0, 0))

	createTeaOwnersTable()

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
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	rows := mock.NewRows([]string{"id", "name"})
	rows.AddRow("1", "Black Tea")
	rows.AddRow("2", "Green Tea")

	mock.ExpectQuery("SELECT \\* FROM types;").WillReturnRows(rows)

	teaTypes, err := GetAllTeaTypesFromDatabase()
	if err != nil {
		t.Errorf("Database returned unexpected error: %v\n", err)
	}
	expected := TeaType{1, "Black Tea"}
	if teaTypes[0] != expected {
		t.Errorf("Database returned unexpected result:\n got: %v\n wanted: %v\n", teaTypes[0], expected)
	}
	expected = TeaType{2, "Green Tea"}
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
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	expected := "Black Tea"
	rows := mock.NewRows([]string{"name"})
	rows.AddRow(expected)
	teaType := TeaType{ID: 1}

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
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	expected := ""
	rows := mock.NewRows([]string{"name"})
	rows.AddRow(expected)
	teaType := TeaType{ID: 1}

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
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	teaName := "Black Tea"
	rows := mock.NewRows([]string{"id"})
	rows.AddRow("1")
	mock.ExpectExec("INSERT INTO types").WithArgs("Black Tea").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT ID FROM types").WillReturnRows(rows)

	teaType := TeaType{ID: 1, Name: teaName}
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
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	teaName := "Black Tea"
	teaID := 1
	rows := mock.NewRows([]string{"name"})
	rows.AddRow(teaName)

	mock.ExpectQuery("SELECT name FROM types").WithArgs(1).WillReturnRows(rows)
	mock.ExpectExec("DELETE FROM types").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	teaType := TeaType{ID: teaID}
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
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	// teaName := "Black Tea"
	teaID := 1
	rows := mock.NewRows([]string{"name"})

	mock.ExpectQuery("SELECT name FROM types").WithArgs(1).WillReturnRows(rows)

	teaType := TeaType{ID: teaID}
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
	expected := Owner{1, "John"}
	if owners[0] != expected {
		t.Errorf("Database returned unexpected result:\n got: %v\n wanted: %v\n", owners[0], expected)
	}
	expected = Owner{2, "Jane"}
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
	owner := Owner{ID: 1}

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
	owner := Owner{ID: 1}

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

	owner := Owner{ID: 1, Name: ownerName}
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

	owner := Owner{ID: ownerID}
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

	owner := Owner{ID: ownerID}
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
	expected := Tea{1, "Snowball", TeaType{1, "Black Tea"}}
	if teas[0] != expected {
		t.Errorf("Database returned unexpected result:\n got: %v\n wanted: %v\n", teas[0], expected)
	}
	expected = Tea{2, "Nearly Nirvana", TeaType{2, "White Tea"}}
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
	tea := Tea{ID: expectedTeaID}
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
	tea := Tea{ID: expectedTeaID}
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

	tea := Tea{Name: teaName, TeaType: TeaType{ID: typeID}}
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

	tea := Tea{Name: teaName, TeaType: TeaType{ID: teaTypeID}}
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

	tea := Tea{Name: teaName, TeaType: TeaType{ID: teaTypeID}}
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

	tea := Tea{ID: teaID}
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

	tea := Tea{ID: teaID}
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

func TestGetTeaOwnerFromDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v\n", err)
	}
	defer db.Close()
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	expectedTeaID := 1
	tea := Tea{ID: expectedTeaID}
	expectedOwners := []Owner{{1, "John"}, {2, "Jane"}}
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
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	expectedTeaID := 1
	tea := Tea{ID: expectedTeaID}

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
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	teaID := 1
	owner := Owner{ID: 1}
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
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	teaID := 1
	owner := Owner{ID: 1}

	mock.ExpectExec("INSERT INTO teaOwners").WithArgs(teaID, owner.ID).WillReturnError(errors.New("UNIQUE constraint failed"))

	if _,err := CreateTeaOwnerInDatabase(teaID, &owner); err.Error() != "This relationship already exists" {
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
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	teaID := 1
	owner := Owner{ID: 1}

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
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	tea := Tea{ID: 1}
	owner := Owner{ID: 1}

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
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	tea := Tea{ID: 1}
	owner := Owner{ID: 1}

	mock.ExpectExec("DELETE FROM tea").WithArgs(tea.ID, owner.ID).WillReturnError(sql.ErrNoRows)

	err = DeleteTeaOwnerFromDatabase(&tea, &owner)
	if err != sql.ErrNoRows {
		t.Errorf("Unexpected error whilst trying to delete tea from database: %v\n", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s\n", err)
	}
}
