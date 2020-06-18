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

	// Tea Types
	router.HandleFunc("/types", getAllTeaTypesHandler).Methods(http.MethodGet)
	router.HandleFunc("/type/{id:[0-9]+}", getTeaTypeHandler).Methods(http.MethodGet)
	router.HandleFunc("/type", createTeaTypeHandler).Methods(http.MethodPost)
	router.HandleFunc("/type/{id:[0-9]+}", deleteTeaTypeHandler).Methods(http.MethodDelete)

	// Tea Owners
	router.HandleFunc("/owners", getAllOwnersHandler).Methods(http.MethodGet)
	router.HandleFunc("/owner/{id:[0-9]+}", getOwnerHandler).Methods(http.MethodGet)

	addr := ":" + cfg.Server.Port
	log.Fatal(http.ListenAndServe(addr, router))
}
