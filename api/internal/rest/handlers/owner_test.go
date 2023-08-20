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

func TestGetAllOwnersHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/owners", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Mock the response from the database
	oldFunc := GetAllOwnersFunc
	defer func() { GetAllOwnersFunc = oldFunc }()
	GetAllOwnersFunc = allTeaOwnersResponseMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllOwnersHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET /owners returned wrong status code:\n got: %v\n want: %v", status, http.StatusOK)
	}

	expected := `[{"id":1,"name":"John"},{"id":2,"name":"Jane"}]`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("GET /ownerss returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func allTeaOwnersResponseMock() ([]models.Owner, error) {
	owner1 := models.Owner{ID: 1, Name: "John"}
	owner2 := models.Owner{ID: 2, Name: "Jane"}
	return []models.Owner{owner1, owner2}, nil
}

func TestGetOwnerHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/owner", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	// Mock the response from the database
	oldFunc := GetOwnerFunc
	defer func() { GetOwnerFunc = oldFunc }()
	GetOwnerFunc = getOwnerResponseMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetOwnerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET /owner/1 returned wrong status code:\n got: %v\n want: %v", status, http.StatusOK)
	}

	expected := `{"id":1,"name":"John"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("GET /owner/1 returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func getOwnerResponseMock(owner *models.Owner) error {
	owner.Name = "John"
	return nil
}

func TestGetOwnerHandlerError(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/owner", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"id": "10"}
	req = mux.SetURLVars(req, vars)

	// Mock the response from the database
	oldFunc := GetOwnerFunc
	defer func() { GetOwnerFunc = oldFunc }()
	GetOwnerFunc = getHandlerErrorResponseMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetOwnerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("GET /owner/10 returned wrong status code:\n got: %v\n want: %v", status, http.StatusInternalServerError)
	}

	expected := `{"error":"ID does not exist in database"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("GET /owner/10 returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func getHandlerErrorResponseMock(owner *models.Owner) error {
	return sql.ErrNoRows
}

func TestCreateOwnerHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/owner", strings.NewReader(`{"name": "John"}`))
	if err != nil {
		t.Fatal(err)
	}

	// Mock the response from the database
	oldFunc := CreateOwnerFunc
	defer func() { CreateOwnerFunc = oldFunc }()
	CreateOwnerFunc = createOwnerResponseMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateOwnerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("POST /owner returned wrong status code:\n got: %v\n want: %v", status, http.StatusCreated)
	}

	expected := `{"id":10,"name":"John"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("POST /owner returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func createOwnerResponseMock(owner *models.Owner) error {
	owner.ID = 10
	return nil
}

func TestErrorCreateOwnerHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/owner", strings.NewReader(`{"name": "John"}`))
	if err != nil {
		t.Fatal(err)
	}

	// Mock the response from the database
	oldFunc := CreateOwnerFunc
	defer func() { CreateOwnerFunc = oldFunc }()
	CreateOwnerFunc = createOwnerResponseErrorMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateOwnerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("POST /owner returned wrong status code:\n got: %v\n want: %v", status, http.StatusInternalServerError)
	}

	expected := `{"error":"Error"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("POST /owner returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func createOwnerResponseErrorMock(owner *models.Owner) error {
	return errors.New("Error")
}

func TestDeleteOwnerHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/owner", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	// Mock the response from the database
	oldFunc := DeleteOwnerFunc
	defer func() { DeleteOwnerFunc = oldFunc }()
	DeleteOwnerFunc = deleteOwnerResponseMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteOwnerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("DELETE /owner returned wrong status code:\n got: %v\n want: %v", status, http.StatusOK)
	}

	expected := `{"name":"John","result":"success"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("DELETE /owner returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func deleteOwnerResponseMock(owner *models.Owner) error {
	owner.Name = "John"
	return nil
}

func TestDeleteOwnerErrorHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/owner", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	// Mock the response from the database
	oldFunc := DeleteOwnerFunc
	defer func() { DeleteOwnerFunc = oldFunc }()
	DeleteOwnerFunc = deleteOwnerResponseErrorMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteOwnerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("DELETE /owner returned wrong status code:\n got: %v\n want: %v", status, http.StatusInternalServerError)
	}

	expected := `{"error":"ID does not exist in database"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("DELETE /owner returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func deleteOwnerResponseErrorMock(owner *models.Owner) error {
	return errors.New("sql: Rows are closed")
}
