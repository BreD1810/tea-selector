package main

import (
	"log"
	"net/http"

	"github.com/BreD1810/tea-selector/api/internal/config"
	"github.com/BreD1810/tea-selector/api/internal/database"
	"github.com/BreD1810/tea-selector/api/internal/rest/handlers"
	"github.com/BreD1810/tea-selector/api/internal/rest/middleware"
	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalln(err)
	}
	middleware.SetSigningKey(cfg.Server.SigningKey)
	database.InitialiseDatabase(*cfg)

	router := mux.NewRouter().StrictSlash(true)

	// Account functions
	router.HandleFunc("/login", handlers.LoginHandler).Methods(http.MethodPost)
	router.Handle("/changepassword", middleware.IsAuthorized(handlers.ChangePasswordHandler)).Methods(http.MethodPost)
	if cfg.Server.RegisterEnabled {
		router.HandleFunc("/register", handlers.RegisterHandler).Methods(http.MethodPost)
	}

	// Tea Types
	router.Handle("/types", middleware.IsAuthorized(handlers.GetAllTeaTypesHandler)).Methods(http.MethodGet)
	router.Handle("/types/teas", middleware.IsAuthorized(handlers.GetAllTeasTypesHandler)).Methods(http.MethodGet)
	router.Handle("/type/{id:[0-9]+}", middleware.IsAuthorized(handlers.GetTeaTypeHandler)).Methods(http.MethodGet)
	router.Handle("/type", middleware.IsAuthorized(handlers.CreateTeaTypeHandler)).Methods(http.MethodPost)
	router.Handle("/type/{id:[0-9]+}", middleware.IsAuthorized(handlers.DeleteTeaTypeHandler)).Methods(http.MethodDelete)

	// Tea Owners
	router.Handle("/owners", middleware.IsAuthorized(handlers.GetAllOwnersHandler)).Methods(http.MethodGet)
	router.Handle("/owners/teas", middleware.IsAuthorized(handlers.GetAllOwnersTeasHandler)).Methods(http.MethodGet)
	router.Handle("/owner/{id:[0-9]+}", middleware.IsAuthorized(handlers.GetOwnerHandler)).Methods(http.MethodGet)
	router.Handle("/owner", middleware.IsAuthorized(handlers.CreateOwnerHandler)).Methods(http.MethodPost)
	router.Handle("/owner/{id:[0-9]+}", middleware.IsAuthorized(handlers.DeleteOwnerHandler)).Methods(http.MethodDelete)

	// Tea
	router.Handle("/teas", middleware.IsAuthorized(handlers.GetAllTeasHandler)).Methods(http.MethodGet)
	router.Handle("/teas/owners", middleware.IsAuthorized(handlers.GetAllTeaOwnersHandler)).Methods(http.MethodGet)
	router.Handle("/tea/{id:[0-9]+}", middleware.IsAuthorized(handlers.GetTeaHandler)).Methods(http.MethodGet)
	router.Handle("/tea", middleware.IsAuthorized(handlers.CreateTeaHandler)).Methods(http.MethodPost)
	router.Handle("/tea/{id:[0-9]+}", middleware.IsAuthorized(handlers.DeleteTeaHandler)).Methods(http.MethodDelete)
	router.Handle("/tea/{id:[0-9]+}/owners", middleware.IsAuthorized(handlers.GetTeaOwnersHandler)).Methods(http.MethodGet)
	router.Handle("/tea/{id:[0-9]+}/owner", middleware.IsAuthorized(handlers.CreateTeaOwnerHandler)).Methods(http.MethodPost)
	router.Handle("/tea/{teaID:[0-9]+}/owner/{ownerID:[0-9]+}", middleware.IsAuthorized(handlers.DeleteTeaOwnerHandler)).Methods(http.MethodDelete)

	addr := ":" + cfg.Server.Port
	log.Fatal(http.ListenAndServe(addr, router))
}
