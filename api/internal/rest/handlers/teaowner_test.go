package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BreD1810/tea-selector/api/internal/models"
	"github.com/gorilla/mux"
)

func TestGetTeaOwnersHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/tea/1/owners", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	// Mock the response from the database
	oldFunc := GetTeaOwnersFunc
	defer func() { GetTeaOwnersFunc = oldFunc }()
	GetTeaOwnersFunc = getTeaOwnersResponseMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTeaOwnersHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET /tea/{id}/owners returned wrong status code:\n got: %v\n want: %v", status, http.StatusOK)
	}

	expected := `[{"id":1,"name":"John"},{"id":2,"name":"Jane"}]`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("GET /tea/{id}/owners returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func getTeaOwnersResponseMock(tea *models.Tea) ([]models.Owner, error) {
	owner1 := models.Owner{ID: 1, Name: "John"}
	owner2 := models.Owner{ID: 2, Name: "Jane"}
	return []models.Owner{owner1, owner2}, nil
}

func TestGetTeaOwnersErrorHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/tea/10/owners", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"id": "10"}
	req = mux.SetURLVars(req, vars)

	// Mock the response from the database
	oldFunc := GetTeaOwnersFunc
	defer func() { GetTeaOwnersFunc = oldFunc }()
	GetTeaOwnersFunc = getTeaOwnersErrorResponseMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTeaOwnersHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("GET /tea/{id}/owners returned wrong status code:\n got: %v\n want: %v", status, http.StatusInternalServerError)
	}

	expected := `{"error":"Error"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("GET /tea/{id}/owners returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func getTeaOwnersErrorResponseMock(tea *models.Tea) ([]models.Owner, error) {
	return nil, errors.New("Error")
}

func TestGetAllTeaOwnersHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/tea/owners", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Mock the response from the database
	oldFunc := GetAllTeaOwnersFunc
	defer func() { GetAllTeaOwnersFunc = oldFunc }()
	GetAllTeaOwnersFunc = getAllTeaOwnersResponseMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllTeaOwnersHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET /tea/owners returned wrong status code:\n got: %v\n want: %v", status, http.StatusOK)
	}

	expected := `[` +
		`{"tea":{"id":1,"name":"Snowball","type":{"id":1,"name":"Black Tea"}},"owners":[{"id":1,"name":"John"}]},` +
		`{"tea":{"id":2,"name":"Nearly Nirvana","type":{"id":2,"name":"White Tea"}},"owners":[{"id":1,"name":"John"},{"id":2,"name":"Jane"}]},` +
		`{"tea":{"id":3,"name":"Earl Grey","type":{"id":1,"name":"Black Tea"}},"owners":[]}` +
		`]`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("GET /tea/owners returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func getAllTeaOwnersResponseMock() ([]models.TeaWithOwners, error) {
	tea1 := models.Tea{ID: 1, Name: "Snowball", TeaType: models.TeaType{ID: 1, Name: "Black Tea"}}
	tea2 := models.Tea{ID: 2, Name: "Nearly Nirvana", TeaType: models.TeaType{ID: 2, Name: "White Tea"}}
	tea3 := models.Tea{ID: 3, Name: "Earl Grey", TeaType: models.TeaType{ID: 1, Name: "Black Tea"}}
	owner1 := models.Owner{ID: 1, Name: "John"}
	owner2 := models.Owner{ID: 2, Name: "Jane"}
	teaWithOwners1 := models.TeaWithOwners{Tea: tea1, Owners: []models.Owner{owner1}}
	teaWithOwners2 := models.TeaWithOwners{Tea: tea2, Owners: []models.Owner{owner1, owner2}}
	teaWithOwners3 := models.TeaWithOwners{Tea: tea3, Owners: []models.Owner{}}
	return []models.TeaWithOwners{teaWithOwners1, teaWithOwners2, teaWithOwners3}, nil
}

func TestCreateTeaOwnerHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/tea/1/owner", strings.NewReader(`{"id": 1}`))
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	// Mock the response from the database
	oldFunc := CreateTeaOwnerFunc
	defer func() { CreateTeaOwnerFunc = oldFunc }()
	CreateTeaOwnerFunc = createTeaOwnerResponseMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateTeaOwnerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("POST /tea/{id}/owner returned wrong status code:\n got: %v\n want: %v", status, http.StatusCreated)
	}

	expected := `{"id":1,"name":"Snowball","type":{"id":1,"name":"Black Tea"}}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("POST /tea/{id}/owner returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func createTeaOwnerResponseMock(teaID int, owner *models.Owner) (models.Tea, error) {
	return models.Tea{ID: 1, Name: "Snowball", TeaType: models.TeaType{ID: 1, Name: "Black Tea"}}, nil
}

func TestCreateTeaOwnerErrorHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/tea/10/owner", strings.NewReader(`{"id": 10}`))
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"id": "10"}
	req = mux.SetURLVars(req, vars)

	// Mock the response from the database
	oldFunc := CreateTeaOwnerFunc
	defer func() { CreateTeaOwnerFunc = oldFunc }()
	CreateTeaOwnerFunc = createTeaOwnerResponseErrorMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateTeaOwnerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("POST /tea/{id}/owner returned wrong status code:\n got: %v\n want: %v", status, http.StatusInternalServerError)
	}

	expected := `{"error":"Error"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("POST /tea/{id}/owner returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func createTeaOwnerResponseErrorMock(teaID int, owner *models.Owner) (models.Tea, error) {
	tea := new(models.Tea)
	return *tea, errors.New("Error")
}

func TestDeleteTeaOwnerHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/tea/1/owner/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"teaID": "1", "ownerID": "1"}
	req = mux.SetURLVars(req, vars)

	// Mock the response from the database
	oldFunc := DeleteTeaOwnerFunc
	defer func() { DeleteTeaOwnerFunc = oldFunc }()
	DeleteTeaOwnerFunc = deleteTeaOwnerResponseMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteTeaOwnerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("DELETE /tea/{id}/owner/{id} returned wrong status code:\n got: %v\n want: %v", status, http.StatusOK)
	}

	expected := `{"result":"success"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("DELETE /tea/{id}/owner{id} returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func deleteTeaOwnerResponseMock(tea *models.Tea, owner *models.Owner) error {
	return nil
}

func TestDeleteTeaOwnerErrorHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/tea/10/owner/10", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"teaID": "10", "ownerID": "10"}
	req = mux.SetURLVars(req, vars)

	// Mock the response from the database
	oldFunc := DeleteTeaOwnerFunc
	defer func() { DeleteTeaOwnerFunc = oldFunc }()
	DeleteTeaOwnerFunc = deleteTeaOwnerResponseErrorMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteTeaOwnerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("DELETE /tea/{id}/owner/{id} returned wrong status code:\n got: %v\n want: %v", status, http.StatusInternalServerError)
	}

	expected := `{"error":"Relationship does not exist in database"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("DELETE /tea/{id}/owner/{id} returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func deleteTeaOwnerResponseErrorMock(tea *models.Tea, owner *models.Owner) error {
	return errors.New("sql: Rows are closed")
}
