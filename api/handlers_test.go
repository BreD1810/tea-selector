package main

import (
	"database/sql"
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
	GetAllTeaTypesFunc = allTeaTypeResponseMock

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

func allTeaTypeResponseMock() ([]TeaType, error) {
	tea1 := TeaType{ID: 1, Name: "Black Tea"}
	tea2 := TeaType{ID: 2, Name: "Green Tea"}
	return []TeaType{tea1, tea2}, nil
}

func TestGetTeaTypeHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/type", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	// Mock the response from the database
	oldFunc := GetTeaTypeFunc
	defer func() { GetTeaTypeFunc = oldFunc }()
	GetTeaTypeFunc = getTeaTypeResponseMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getTeaTypeHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET /type/1 returned wrong status code:\n got: %v\n want: %v", status, http.StatusOK)
	}

	expected := `{"id":1,"name":"Black Tea"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("GET /type/1 returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func getTeaTypeResponseMock(teaType *TeaType) error {
	teaType.Name = "Black Tea"
	return nil
}

func TestGetTeaTypeHandlerError(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/type", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"id": "10"}
	req = mux.SetURLVars(req, vars)

	// Mock the response from the database
	oldFunc := GetTeaTypeFunc
	defer func() { GetTeaTypeFunc = oldFunc }()
	GetTeaTypeFunc = getTeaTypeErrorResponseMock

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getTeaTypeHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("GET /type/10 returned wrong status code:\n got: %v\n want: %v", status, http.StatusInternalServerError)
	}

	expected := `{"error":"ID does not exist in database"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("GET /type/10 returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func getTeaTypeErrorResponseMock(teaType *TeaType) error {
	return sql.ErrNoRows
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
		t.Errorf("POST /type returned wrong status code:\n got: %v\n want: %v", status, http.StatusCreated)
	}

	expected := `{"id":10,"name":"Black Tea"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("POST /type returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
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
		t.Errorf("POST /type returned wrong status code:\n got: %v\n want: %v", status, http.StatusInternalServerError)
	}

	expected := `{"error":"Error"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("POST /type returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
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
	handler := http.HandlerFunc(getAllOwnersHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET /owners returned wrong status code:\n got: %v\n want: %v", status, http.StatusOK)
	}

	expected := `[{"id":1,"name":"John"},{"id":2,"name":"Jane"}]`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("GET /ownerss returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func allTeaOwnersResponseMock() ([]Owner, error) {
	owner1 := Owner{ID: 1, Name: "John"}
	owner2 := Owner{ID: 2, Name: "Jane"}
	return []Owner{owner1, owner2}, nil
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
	handler := http.HandlerFunc(getOwnerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET /owner/1 returned wrong status code:\n got: %v\n want: %v", status, http.StatusOK)
	}

	expected := `{"id":1,"name":"John"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("GET /owner/1 returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func getOwnerResponseMock(owner *Owner) error {
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
	handler := http.HandlerFunc(getOwnerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("GET /owner/10 returned wrong status code:\n got: %v\n want: %v", status, http.StatusInternalServerError)
	}

	expected := `{"error":"ID does not exist in database"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("GET /owner/10 returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func getHandlerErrorResponseMock(owner *Owner) error {
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
	handler := http.HandlerFunc(createOwnerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("POST /owner returned wrong status code:\n got: %v\n want: %v", status, http.StatusCreated)
	}

	expected := `{"id":10,"name":"John"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("POST /owner returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func createOwnerResponseMock(owner *Owner) error {
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
	handler := http.HandlerFunc(createOwnerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("POST /owner returned wrong status code:\n got: %v\n want: %v", status, http.StatusInternalServerError)
	}

	expected := `{"error":"Error"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("POST /owner returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func createOwnerResponseErrorMock(owner *Owner) error {
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
	handler := http.HandlerFunc(deleteOwnerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("DELETE /owner returned wrong status code:\n got: %v\n want: %v", status, http.StatusOK)
	}

	expected := `{"name":"John","result":"success"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("DELETE /owner returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func deleteOwnerResponseMock(owner *Owner) error {
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
	handler := http.HandlerFunc(deleteOwnerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("DELETE /owner returned wrong status code:\n got: %v\n want: %v", status, http.StatusInternalServerError)
	}

	expected := `{"error":"ID does not exist in database"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("DELETE /owner returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func deleteOwnerResponseErrorMock(owner *Owner) error {
	return errors.New("sql: Rows are closed")
}

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
	handler := http.HandlerFunc(getAllTeasHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET /teas returned wrong status code:\n got: %v\n want: %v", status, http.StatusOK)
	}

	expected := `[{"id":1,"name":"Snowball","type":{"id":1,"name":"Black Tea"}},{"id":2,"name":"Nearly Nirvana","type":{"id":2,"name":"White Tea"}}]`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("GET /teas returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func allTeasResponseMock() ([]Tea, error) {
	tea1 := Tea{ID: 1, Name: "Snowball", TeaType: TeaType{ID: 1, Name: "Black Tea"}}
	tea2 := Tea{ID: 2, Name: "Nearly Nirvana", TeaType: TeaType{ID: 2, Name: "White Tea"}}
	return []Tea{tea1, tea2}, nil
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
	handler := http.HandlerFunc(getTeaHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET /tea/1 returned wrong status code:\n got: %v\n want: %v", status, http.StatusOK)
	}

	expected := `{"id":1,"name":"Snowball","type":{"id":1,"name":"Black Tea"}}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("GET /owner/1 returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func getTeaResponseMock(tea *Tea) error {
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
	handler := http.HandlerFunc(getTeaHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("GET /tea/10 returned wrong status code:\n got: %v\n want: %v", status, http.StatusInternalServerError)
	}

	expected := `{"error":"ID does not exist in database"}`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("GET /tea/10 returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func getTeaResponseErrorMock(tea *Tea) error {
	return sql.ErrNoRows
}
