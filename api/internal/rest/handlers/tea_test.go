package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BreD1810/tea-selector/api/internal/models"
	"github.com/gorilla/mux"
)

func TestGetAllTeasHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/teas", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Mock the response from the database
	oldFunc := GetAllTeasFunc
	defer func() { GetAllTeasFunc = oldFunc }()
	GetAllTeasFunc = allTeasResponseMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllTeasHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET /teas returned wrong status code:\n got: %v\n want: %v", status, http.StatusOK)
	}

	expected := `[{"id":1,"name":"Snowball","type":{"id":1,"name":"Black Tea"}},{"id":2,"name":"Nearly Nirvana","type":{"id":2,"name":"White Tea"}}]`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("GET /teas returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func allTeasResponseMock() ([]models.Tea, error) {
	tea1 := models.Tea{ID: 1, Name: "Snowball", TeaType: models.TeaType{ID: 1, Name: "Black Tea"}}
	tea2 := models.Tea{ID: 2, Name: "Nearly Nirvana", TeaType: models.TeaType{ID: 2, Name: "White Tea"}}
	return []models.Tea{tea1, tea2}, nil
}

func TestGetTeaHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/tea", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	// Mock the response from the database
	oldFunc := GetTeaFunc
	defer func() { GetTeaFunc = oldFunc }()
	GetTeaFunc = getTeaResponseMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTeaHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET /tea/1 returned wrong status code:\n got: %v\n want: %v", status, http.StatusOK)
	}

	expected := `{"id":1,"name":"Snowball","type":{"id":1,"name":"Black Tea"}}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("GET /owner/1 returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func getTeaResponseMock(tea *models.Tea) error {
	tea.Name = "Snowball"
	tea.TeaType.ID = 1
	tea.TeaType.Name = "Black Tea"
	return nil
}

func TestGetTeaHandlerError(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/tea", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"id": "10"}
	req = mux.SetURLVars(req, vars)

	// Mock the response from the database
	oldFunc := GetTeaFunc
	defer func() { GetTeaFunc = oldFunc }()
	GetTeaFunc = getTeaResponseErrorMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTeaHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("GET /tea/10 returned wrong status code:\n got: %v\n want: %v", status, http.StatusInternalServerError)
	}

	expected := `{"error":"ID does not exist in database"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("GET /tea/10 returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func getTeaResponseErrorMock(tea *models.Tea) error {
	return sql.ErrNoRows
}

func TestCreateTeaHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/tea", strings.NewReader(`{"name": "Snowball","type":{"id":1}}`))
	if err != nil {
		t.Fatal(err)
	}

	// Mock the response from the database
	oldFunc := CreateTeaFunc
	defer func() { CreateTeaFunc = oldFunc }()
	CreateTeaFunc = createTeaResponseMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateTeaHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("POST /tea returned wrong status code:\n got: %v\n want: %v", status, http.StatusCreated)
	}

	expected := `{"id":1,"name":"Snowball","type":{"id":1,"name":"Black Tea"}}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("POST /tea returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func createTeaResponseMock(tea *models.Tea) error {
	tea.ID = 1
	tea.TeaType.Name = "Black Tea"
	return nil
}

func TestErrorCreateTeaHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/tea", strings.NewReader(`{"name": "Snowball"}`))
	if err != nil {
		t.Fatal(err)
	}

	// Mock the response from the database
	oldFunc := CreateTeaFunc
	defer func() { CreateTeaFunc = oldFunc }()
	CreateTeaFunc = createTeaResponseErrorMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateTeaHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("POST /tea returned wrong status code:\n got: %v\n want: %v", status, http.StatusInternalServerError)
	}

	expected := `{"error":"Error"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("POST /tea returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func createTeaResponseErrorMock(tea *models.Tea) error {
	return errors.New("Error")
}

func TestDeleteTeaHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/tea", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	// Mock the response from the database
	oldFunc := DeleteTeaFunc
	defer func() { DeleteTeaFunc = oldFunc }()
	DeleteTeaFunc = deleteTeaResponseMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteTeaHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("DELETE /tea returned wrong status code:\n got: %v\n want: %v", status, http.StatusOK)
	}

	expected := `{"name":"Snowball","result":"success"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("DELETE /tea returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func deleteTeaResponseMock(tea *models.Tea) error {
	tea.Name = "Snowball"
	return nil
}

func TestDeleteTeaErrorHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/tea", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	// Mock the response from the database
	oldFunc := DeleteTeaFunc
	defer func() { DeleteTeaFunc = oldFunc }()
	DeleteTeaFunc = deleteTeaResponseErrorMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteTeaHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("DELETE /tea returned wrong status code:\n got: %v\n want: %v", status, http.StatusInternalServerError)
	}

	expected := `{"error":"ID does not exist in database"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("DELETE /tea returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func deleteTeaResponseErrorMock(tea *models.Tea) error {
	return errors.New("sql: Rows are closed")
}
