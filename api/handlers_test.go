package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

// https://blog.questionable.services/article/testing-http-handlers-go/
func TestGetAllTeaTypesHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/types", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Mock the response from the database
	oldFunc := GetAllTeaTypesFunc
	defer func() { GetAllTeaTypesFunc = oldFunc }()
	GetAllTeaTypesFunc = allTeaResponseMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getAllTeaTypesHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET /types returned wrong status code:\n got: %v\n want: %v", status, http.StatusOK)
	}

	expected := `[{"id":1,"name":"Black Tea"},{"id":2,"name":"Green Tea"}]`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("GET /types returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func allTeaResponseMock() ([]TeaType, error) {
	tea1 := TeaType{ID: 1, Name: "Black Tea"}
	tea2 := TeaType{ID: 2, Name: "Green Tea"}
	return []TeaType{tea1, tea2}, nil
}

func TestCreateTeaTypeHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/type", strings.NewReader(`{"name": "Black Tea"}`))
	if err != nil {
		t.Fatal(err)
	}

	// Mock the response from the database
	oldFunc := CreateTeaTypeFunc
	defer func() { CreateTeaTypeFunc = oldFunc }()
	CreateTeaTypeFunc = createTeaTypeResponseMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createTeaTypeHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("PUT /type returned wrong status code:\n got: %v\n want: %v", status, http.StatusCreated)
	}

	expected := `{"id":10,"name":"Black Tea"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("PUT /type returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func createTeaTypeResponseMock(teaType *TeaType) error {
	teaType.ID = 10
	return nil
}

func TestErrorCreateTeaTypeHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/type", strings.NewReader(`{"name": "Black Tea"}`))
	if err != nil {
		t.Fatal(err)
	}

	// Mock the response from the database
	oldFunc := CreateTeaTypeFunc
	defer func() { CreateTeaTypeFunc = oldFunc }()
	CreateTeaTypeFunc = createTeaTypeResponseErrorMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createTeaTypeHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PUT /type returned wrong status code:\n got: %v\n want: %v", status, http.StatusInternalServerError)
	}

	expected := `{"error":"Error"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("PUT /type returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func createTeaTypeResponseErrorMock(teaType *TeaType) error {
	return errors.New("Error")
}

func TestDeleteTeaTypeHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/type", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	// Mock the response from the database
	oldFunc := DeleteTeaTypeFunc
	defer func() { DeleteTeaTypeFunc = oldFunc }()
	DeleteTeaTypeFunc = deleteTeaTypeResponseMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteTeaTypeHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("DELETE /type returned wrong status code:\n got: %v\n want: %v", status, http.StatusOK)
	}

	expected := `{"name":"Black Tea","result":"success"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("DELETE /type returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func deleteTeaTypeResponseMock(teaType *TeaType) error {
	teaType.Name = "Black Tea"
	return nil
}

func TestDeleteTeaTypeErrorHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/type", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	// Mock the response from the database
	oldFunc := DeleteTeaTypeFunc
	defer func() { DeleteTeaTypeFunc = oldFunc }()
	DeleteTeaTypeFunc = deleteTeaTypeResponseErrorMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteTeaTypeHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("DELETE /type returned wrong status code:\n got: %v\n want: %v", status, http.StatusInternalServerError)
	}

	expected := `{"error":"ID does not exist in database"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("DELETE /type returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func deleteTeaTypeResponseErrorMock(teaType *TeaType) error {
	return errors.New("sql: Rows are closed")
}
