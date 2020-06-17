package main

import (
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

	teaTypes, _ := GetAllTeaTypesFunc()
	expected := TeaType{1, "Black Tea"}
	if teaTypes[0] != expected {
		t.Errorf("Database returned unexpected result:\n got: %v\n wanted: %v", teaTypes[0], expected)
	}
	expected = TeaType{2, "Green Tea"}
	if teaTypes[1] != expected {
		t.Errorf("Database returned unexpected result:\n got: %v\n wanted: %v", teaTypes[1], expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestCreateTeaTypeInDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v", err)
	}
	defer db.Close()
	oldDB := DB
	defer func() { DB = oldDB }()
	DB = db

	teaName := "Black Tea"
	rows := mock.NewRows([]string{"id"})
	rows.AddRow("1")
	mock.ExpectExec("INSERT INTO types \\(name\\) VALUES \\('" + teaName + "'\\)").WillReturnResult(sqlmock.NewResult(1, 1))
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
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestDeleteTeaTypeInDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v", err)
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
	err = DeleteTeaTypeFunc(&teaType)
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
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestDeleteNonExistantTeaTypeInDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v", err)
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
	err = DeleteTeaTypeFunc(&teaType)
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
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}
