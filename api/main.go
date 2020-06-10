package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	cfg := getConfig()
	initialiseDatabase()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", testResponse).Methods("GET")

	addr := ":" + cfg.Server.Port
	log.Fatal(http.ListenAndServe(addr, router))
}
