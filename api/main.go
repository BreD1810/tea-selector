package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Fatalln(err)
	}
	SetSigningKey(cfg.Server.SigningKey)
	initialiseDatabase(*cfg)

	router := mux.NewRouter().StrictSlash(true)

	// Account functions
	router.HandleFunc("/login", loginHandler).Methods(http.MethodPost)
	router.Handle("/changepassword", isAuthorized(changePasswordHandler)).Methods(http.MethodPost)
	if cfg.Server.RegisterEnabled {
		router.HandleFunc("/register", registerHandler).Methods(http.MethodPost)
	}

	// Tea Types
	router.Handle("/types", isAuthorized(getAllTeaTypesHandler)).Methods(http.MethodGet)
	router.Handle("/types/teas", isAuthorized(getAllTeasTypesHandler)).Methods(http.MethodGet)
	router.Handle("/type/{id:[0-9]+}", isAuthorized(getTeaTypeHandler)).Methods(http.MethodGet)
	router.Handle("/type", isAuthorized(createTeaTypeHandler)).Methods(http.MethodPost)
	router.Handle("/type/{id:[0-9]+}", isAuthorized(deleteTeaTypeHandler)).Methods(http.MethodDelete)

	// Tea Owners
	router.Handle("/owners", isAuthorized(getAllOwnersHandler)).Methods(http.MethodGet)
	router.Handle("/owners/teas", isAuthorized(getAllOwnersTeasHandler)).Methods(http.MethodGet)
	router.Handle("/owner/{id:[0-9]+}", isAuthorized(getOwnerHandler)).Methods(http.MethodGet)
	router.Handle("/owner", isAuthorized(createOwnerHandler)).Methods(http.MethodPost)
	router.Handle("/owner/{id:[0-9]+}", isAuthorized(deleteOwnerHandler)).Methods(http.MethodDelete)

	// Tea
	router.Handle("/teas", isAuthorized(getAllTeasHandler)).Methods(http.MethodGet)
	router.Handle("/teas/owners", isAuthorized(getAllTeaOwnersHandler)).Methods(http.MethodGet)
	router.Handle("/tea/{id:[0-9]+}", isAuthorized(getTeaHandler)).Methods(http.MethodGet)
	router.Handle("/tea", isAuthorized(createTeaHandler)).Methods(http.MethodPost)
	router.Handle("/tea/{id:[0-9]+}", isAuthorized(deleteTeaHandler)).Methods(http.MethodDelete)
	router.Handle("/tea/{id:[0-9]+}/owners", isAuthorized(getTeaOwnersHandler)).Methods(http.MethodGet)
	router.Handle("/tea/{id:[0-9]+}/owner", isAuthorized(createTeaOwnerHandler)).Methods(http.MethodPost)
	router.Handle("/tea/{teaID:[0-9]+}/owner/{ownerID:[0-9]+}", isAuthorized(deleteTeaOwnerHandler)).Methods(http.MethodDelete)

	addr := ":" + cfg.Server.Port
	log.Fatal(http.ListenAndServe(addr, router))
}
