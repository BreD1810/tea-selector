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

	router.HandleFunc("/login", loginHandler).Methods(http.MethodPost)

	// Tea Types
	router.HandleFunc("/types", getAllTeaTypesHandler).Methods(http.MethodGet)
	router.HandleFunc("/types/teas", getAllTeasTypesHandler).Methods(http.MethodGet)
	router.HandleFunc("/type/{id:[0-9]+}", getTeaTypeHandler).Methods(http.MethodGet)
	router.HandleFunc("/type", createTeaTypeHandler).Methods(http.MethodPost)
	router.HandleFunc("/type/{id:[0-9]+}", deleteTeaTypeHandler).Methods(http.MethodDelete)

	// Tea Owners
	router.HandleFunc("/owners", getAllOwnersHandler).Methods(http.MethodGet)
	router.HandleFunc("/owners/teas", getAllOwnersTeasHandler).Methods(http.MethodGet)
	router.HandleFunc("/owner/{id:[0-9]+}", getOwnerHandler).Methods(http.MethodGet)
	router.HandleFunc("/owner", createOwnerHandler).Methods(http.MethodPost)
	router.HandleFunc("/owner/{id:[0-9]+}", deleteOwnerHandler).Methods(http.MethodDelete)

	// Tea
	router.HandleFunc("/teas", getAllTeasHandler).Methods(http.MethodGet)
	router.HandleFunc("/teas/owners", getAllTeaOwnersHandler).Methods(http.MethodGet)
	router.HandleFunc("/tea/{id:[0-9]+}", getTeaHandler).Methods(http.MethodGet)
	router.HandleFunc("/tea", createTeaHandler).Methods(http.MethodPost)
	router.HandleFunc("/tea/{id:[0-9]+}", deleteTeaHandler).Methods(http.MethodDelete)
	router.HandleFunc("/tea/{id:[0-9]+}/owners", getTeaOwnersHandler).Methods(http.MethodGet)
	router.HandleFunc("/tea/{id:[0-9]+}/owner", createTeaOwnerHandler).Methods(http.MethodPost)
	router.HandleFunc("/tea/{teaID:[0-9]+}/owner/{ownerID:[0-9]+}", deleteTeaOwnerHandler).Methods(http.MethodDelete)

	addr := ":" + cfg.Server.Port
	log.Fatal(http.ListenAndServe(addr, router))
}
