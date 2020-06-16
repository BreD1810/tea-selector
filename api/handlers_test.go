package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// https://blog.questionable.services/article/testing-http-handlers-go/
func TestGetAllTeaTypesHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/types", nil)
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
		t.Errorf("testResponse returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}

func allTeaResponseMock() []TeaType {
	tea1 := TeaType{ID: 1, Name: "Black Tea"}
	tea2 := TeaType{ID: 2, Name: "Green Tea"}
	return []TeaType{tea1, tea2}
}
