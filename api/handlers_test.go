package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// https://blog.questionable.services/article/testing-http-handlers-go/
func TestTestResponse(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testResponse)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("testResponse returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Hello, world!"
	if actual := rr.Body.String(); actual != expected {
		t.Errorf("testResponse returned unexpected body:\n got: %v\n wanted: %v", actual, expected)
	}
}
