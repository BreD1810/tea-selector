package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg := getConfig()
	initialiseDatabase(cfg)

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/types", getAllTeaTypesHandler).Methods("GET")
	router.HandleFunc("/type", createTeaTypeHandler).Methods("POST")

	addr := ":" + cfg.Server.Port
	log.Fatal(http.ListenAndServe(addr, router))
}
