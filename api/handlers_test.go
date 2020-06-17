package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

	expected := "[{\"id\":1,\"name\":\"Black Tea\"},{\"id\":2,\"name\":\"Green Tea\"}]\n"
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("GET /types returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func allTeaResponseMock() []TeaType {
	tea1 := TeaType{ID: 1, Name: "Black Tea"}
	tea2 := TeaType{ID: 2, Name: "Green Tea"}
	return []TeaType{tea1, tea2}
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

	expected := "{\"id\":10,\"name\":\"Black Tea\"}"
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("PUT /type returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func createTeaTypeResponseMock(teaType *TeaType) error {
	teaType.ID = 10
	return nil
}
