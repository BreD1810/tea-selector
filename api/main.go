package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func testResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", testResponse).Methods("GET")

	log.Fatal(http.ListenAndServe(":7344", router))
}
