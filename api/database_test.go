package main

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateTeaTable(t *testing.T) {
	createTeaString := "CREATE TABLE tea"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectPrepare(createTeaString)
	mock.ExpectExec(createTeaString).WillReturnResult(sqlmock.NewResult(0, 0))

	createTeaTable(db)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expections: %s", err)
	}
}

// Check that no data will be deleted if the table already exists
func TestCreateTeaTableExists(t *testing.T) {
	createTeaString := "CREATE TABLE IF NOT EXISTS tea"
	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "test")
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error occurred setting up mock database: %v", err)
	}
	db.Exec("CREATE TABLE tea (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT);")
	db.Exec("INSERT 'test' INTO tea;")
	defer db.Close()

	mock.ExpectPrepare(createTeaString)
	mock.ExpectExec(createTeaString).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	createTeaTable(db)
	_, err = db.Query("SELECT * FROM tea;")
	if err != nil {
		t.Fatalf("Error reading rows from tea mock: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}
