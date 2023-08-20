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

type mockTeaStore struct {
	teas []models.Tea
	err  error
}

func (ts *mockTeaStore) GetAllTeasFromDatabase() ([]models.Tea, error) {
	if ts.err != nil {
		return nil, ts.err
	}
	return ts.teas, nil
}

func (ts *mockTeaStore) GetTeaFromDatabase(tea *models.Tea) error {
	return ts.err
}

func (ts *mockTeaStore) CreateTeaInDatabase(tea *models.Tea) error {
	return ts.err
}

func (ts *mockTeaStore) DeleteTeaFromDatabase(tea *models.Tea) error {
	return ts.err
}

func TestGetAllTeas(t *testing.T) {
	testCases := []struct {
		name               string
		teaStoreFn         func() *mockTeaStore
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "success getting all teas",
			teaStoreFn: func() *mockTeaStore {
				teas := []models.Tea{
					models.Tea{ID: 1, Name: "Snowball", TeaType: models.TeaType{ID: 1, Name: "Black Tea"}},
					models.Tea{ID: 2, Name: "Nearly Nirvana", TeaType: models.TeaType{ID: 2, Name: "White Tea"}},
				}
				return &mockTeaStore{teas: teas}
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `[{"id":1,"name":"Snowball","type":{"id":1,"name":"Black Tea"}},{"id":2,"name":"Nearly Nirvana","type":{"id":2,"name":"White Tea"}}]`,
		},
		{
			name: "db error getting all teas",
			teaStoreFn: func() *mockTeaStore {
				return &mockTeaStore{err: errors.New("rigged tea store error")}
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"rigged tea store error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			teaStore := tc.teaStoreFn()
			teaHandler := NewTeaHandler(teaStore)
			handler := http.HandlerFunc(teaHandler.GetAllTeas)
			rr := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, "/teas", nil)
			if err != nil {
				t.Fatal(err)
			}
			handler.ServeHTTP(rr, req)
			if rr.Code != tc.expectedStatusCode {
				t.Errorf("GET /teas returned wrong status code:\n got: %v\n want: %v", rr.Code, tc.expectedStatusCode)
			}

			if actual := rr.Body.String(); actual != tc.expectedBody {
				t.Errorf("GET /teas returned unexpected body:\n got: %v\n wanted: %v", actual, tc.expectedBody)
			}
		})
	}
}

func TestGetTea(t *testing.T) {
	testCases := []struct {
		name               string
		teaStoreFn         func() *mockTeaStore
		requestVars        map[string]string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "success getting tea",
			teaStoreFn: func() *mockTeaStore {
				return &mockTeaStore{}
			},
			requestVars:        map[string]string{"id": "1"},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"id":1,"name":"","type":{"id":0,"name":""}}`,
		},
		{
			name: "no tea found",
			teaStoreFn: func() *mockTeaStore {
				return &mockTeaStore{err: sql.ErrNoRows}
			},
			requestVars:        map[string]string{"id": "10"},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"ID does not exist in database"}`,
		},
		{
			name: "db error getting all teas",
			teaStoreFn: func() *mockTeaStore {
				return &mockTeaStore{err: errors.New("rigged tea store error")}
			},
			requestVars:        map[string]string{"id": "1"},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"rigged tea store error"}`,
		},
		{
			name: "error parsing request vars",
			teaStoreFn: func() *mockTeaStore {
				return &mockTeaStore{}
			},
			requestVars:        map[string]string{"id": "bad value"},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"Invalid tea ID"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			teaStore := tc.teaStoreFn()
			teaHandler := NewTeaHandler(teaStore)
			handler := http.HandlerFunc(teaHandler.GetTea)
			rr := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, "/tea", nil)
			if err != nil {
				t.Fatal(err)
			}
			req = mux.SetURLVars(req, tc.requestVars)

			handler.ServeHTTP(rr, req)
			if rr.Code != tc.expectedStatusCode {
				t.Errorf("GET /tea/1 returned wrong status code:\n got: %v\n want: %v", rr.Code, tc.expectedStatusCode)
			}

			if actual := rr.Body.String(); actual != tc.expectedBody {
				t.Errorf("GET /tea/1 returned unexpected body:\n got: %v\n wanted: %v", actual, tc.expectedBody)
			}
		})
	}
}

func TestCreateTea(t *testing.T) {
	testCases := []struct {
		name               string
		teaStoreFn         func() *mockTeaStore
		requestBody        string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:        "success creating tea",
			requestBody: `{"id":1,"name":"Snowball","type":{"id":1,"name":"Black Tea"}}`,
			teaStoreFn: func() *mockTeaStore {
				return &mockTeaStore{}
			},
			expectedStatusCode: http.StatusCreated,
			expectedBody:       `{"id":1,"name":"Snowball","type":{"id":1,"name":"Black Tea"}}`,
		},
		{
			name:        "db error creating tea",
			requestBody: `{"id":1,"name":"Snowball","type":{"id":1,"name":"Black Tea"}}`,
			teaStoreFn: func() *mockTeaStore {
				return &mockTeaStore{err: errors.New("rigged tea store error")}
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"rigged tea store error"}`,
		},
		{
			name:        "error parsing request body",
			requestBody: ``,
			teaStoreFn: func() *mockTeaStore {
				return &mockTeaStore{}
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"Invalid request payload"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			teaStore := tc.teaStoreFn()
			teaHandler := NewTeaHandler(teaStore)
			handler := http.HandlerFunc(teaHandler.CreateTea)
			rr := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodPost, "/tea", strings.NewReader(tc.requestBody))
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatusCode {
				t.Errorf("POST /tea returned wrong status code:\n got: %v\n want: %v", rr.Code, tc.expectedStatusCode)
			}

			if actual := rr.Body.String(); actual != tc.expectedBody {
				t.Errorf("POST /tea returned unexpected body:\n got: %v\n wanted: %v", actual, tc.expectedBody)
			}
		})
	}
}

func TestDeleteTea(t *testing.T) {
	testCases := []struct {
		name               string
		teaStoreFn         func() *mockTeaStore
		requestVars        map[string]string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "success deleting tea",
			teaStoreFn: func() *mockTeaStore {
				return &mockTeaStore{}
			},
			requestVars:        map[string]string{"id": "1"},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"name":"","result":"success"}`,
		},
		{
			name: "db error creating tea",
			teaStoreFn: func() *mockTeaStore {
				return &mockTeaStore{err: errors.New("rigged tea store error")}
			},
			requestVars:        map[string]string{"id": "1"},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"rigged tea store error"}`,
		},
		{
			name: "tea does not exist",
			teaStoreFn: func() *mockTeaStore {
				return &mockTeaStore{err: errors.New("sql: Rows are closed")}
			},
			requestVars:        map[string]string{"id": "1"},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"ID does not exist in database"}`,
		},
		{
			name: "error parsing request vars",
			teaStoreFn: func() *mockTeaStore {
				return &mockTeaStore{}
			},
			requestVars:        map[string]string{"id": "bad value"},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"Invalid tea ID"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			teaStore := tc.teaStoreFn()
			teaHandler := NewTeaHandler(teaStore)
			handler := http.HandlerFunc(teaHandler.DeleteTea)
			rr := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodDelete, "/tea", nil)
			if err != nil {
				t.Fatal(err)
			}
			req = mux.SetURLVars(req, tc.requestVars)

			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatusCode {
				t.Errorf("DELETE /tea returned wrong status code:\n got: %v\n want: %v", rr.Code, tc.expectedStatusCode)
			}

			if actual := rr.Body.String(); actual != tc.expectedBody {
				t.Errorf("DELETE /tea returned unexpected body:\n got: %v\n wanted: %v", actual, tc.expectedBody)
			}
		})
	}
}
