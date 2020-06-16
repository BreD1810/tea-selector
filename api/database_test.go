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

	mock.ExpectExec(createTeaTypeString).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("INSERT INTO types \\(name\\) VALUES \\('" + teaTypes[0] + "'\\), \\('" + teaTypes[1] + "'\\);").WillReturnResult(sqlmock.NewResult(2, 2))

	createTeaTypeTable(db, teaTypes)

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

	mock.ExpectExec(createTeaTypeString).WillReturnResult(sqlmock.NewResult(0, 0))

	createTeaTypeTable(db, teaTypes)

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

	mock.ExpectExec(createTeaString).WillReturnResult(sqlmock.NewResult(0, 0))

	createTeaTable(db)

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

	mock.ExpectExec(createOwnerString).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("INSERT INTO owner \\(name\\) VALUES \\('" + owners[0] + "'\\), \\('" + owners[1] + "'\\);").WillReturnResult(sqlmock.NewResult(2, 2))

	createOwnerTable(db, owners)

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

	mock.ExpectExec(createOwnerString).WillReturnResult(sqlmock.NewResult(0, 0))

	createOwnerTable(db, owners)

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

	mock.ExpectExec(createTeaOwnersString).WillReturnResult(sqlmock.NewResult(0, 0))

	createTeaOwnersTable(db)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expections: %s", err)
	}
}
