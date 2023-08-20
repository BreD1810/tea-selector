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
	handler := http.HandlerFunc(GetAllTeaTypesHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET /types returned wrong status code:\n got: %v\n want: %v", status, http.StatusOK)
	}

	expected := `[{"id":1,"name":"Black Tea"},{"id":2,"name":"Green Tea"}]`
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("GET /types returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func allTeaTypeResponseMock() ([]models.TeaType, error) {
	tea1 := models.TeaType{ID: 1, Name: "Black Tea"}
	tea2 := models.TeaType{ID: 2, Name: "Green Tea"}
	return []models.TeaType{tea1, tea2}, nil
}

func TestGetTeaTypeHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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
		handler := http.HandlerFunc(GetTeaTypeHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("GET /type/1 returned wrong status code:\n got: %v\n want: %v", status, http.StatusOK)
		}

		expected := `{"id":1,"name":"Black Tea"}`
		if actual := rr.Body.String(); actual != expected {
			t.Errorf("GET /type/1 returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
		}
	})

	t.Run("error", func(t *testing.T) {
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
		handler := http.HandlerFunc(GetTeaTypeHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("GET /type/10 returned wrong status code:\n got: %v\n want: %v", status, http.StatusInternalServerError)
		}

		expected := `{"error":"ID does not exist in database"}`
		if actual := rr.Body.String(); actual != expected {
			t.Errorf("GET /type/10 returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
		}
	})
}

func getTeaTypeResponseMock(teaType *models.TeaType) error {
	teaType.Name = "Black Tea"
	return nil
}

func getTeaTypeErrorResponseMock(teaType *models.TeaType) error {
	return sql.ErrNoRows
}

func TestCreateTeaTypeHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/type", strings.NewReader(`{"name": "Black Tea"}`))
		if err != nil {
			t.Fatal(err)
		}

		// Mock the response from the database
		oldFunc := CreateTeaTypeFunc
		defer func() { CreateTeaTypeFunc = oldFunc }()
		CreateTeaTypeFunc = createTeaTypeResponseMock

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CreateTeaTypeHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("POST /type returned wrong status code:\n got: %v\n want: %v", status, http.StatusCreated)
		}

		expected := `{"id":10,"name":"Black Tea"}`
		if actual := rr.Body.String(); actual != expected {
			t.Errorf("POST /type returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
		}
	})

	t.Run("error", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/type", strings.NewReader(`{"name": "Black Tea"}`))
		if err != nil {
			t.Fatal(err)
		}

		// Mock the response from the database
		oldFunc := CreateTeaTypeFunc
		defer func() { CreateTeaTypeFunc = oldFunc }()
		CreateTeaTypeFunc = createTeaTypeResponseErrorMock

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CreateTeaTypeHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("POST /type returned wrong status code:\n got: %v\n want: %v", status, http.StatusInternalServerError)
		}

		expected := `{"error":"Error"}`
		if actual := rr.Body.String(); actual != expected {
			t.Errorf("POST /type returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
		}
	})
}

func createTeaTypeResponseMock(teaType *models.TeaType) error {
	teaType.ID = 10
	return nil
}

func createTeaTypeResponseErrorMock(teaType *models.TeaType) error {
	return errors.New("Error")
}

func TestDeleteTeaTypeHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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
		handler := http.HandlerFunc(DeleteTeaTypeHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("DELETE /type returned wrong status code:\n got: %v\n want: %v", status, http.StatusOK)
		}

		expected := `{"name":"Black Tea","result":"success"}`
		if actual := rr.Body.String(); actual != expected {
			t.Errorf("DELETE /type returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
		}
	})

	t.Run("error", func(t *testing.T) {
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
		handler := http.HandlerFunc(DeleteTeaTypeHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("DELETE /type returned wrong status code:\n got: %v\n want: %v", status, http.StatusInternalServerError)
		}

		expected := `{"error":"ID does not exist in database"}`
		if actual := rr.Body.String(); actual != expected {
			t.Errorf("DELETE /type returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
		}
	})
}

func deleteTeaTypeResponseMock(teaType *models.TeaType) error {
	teaType.Name = "Black Tea"
	return nil
}

func deleteTeaTypeResponseErrorMock(teaType *models.TeaType) error {
	return errors.New("sql: Rows are closed")
}