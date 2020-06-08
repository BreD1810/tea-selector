package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
)


func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", testResponse).Methods("GET")

	log.Fatal(http.ListenAndServe(":7344", router))
}
