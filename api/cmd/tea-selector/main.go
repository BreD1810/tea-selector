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
	database, err := database.NewDatabase(*cfg)
	if err != nil {
		log.Fatalln(err)
	}

	router := mux.NewRouter().StrictSlash(true)

	// Account functions
	userHandler := handlers.NewUserHandler(database)
	router.HandleFunc("/login", userHandler.Login).Methods(http.MethodPost)
	router.Handle("/changepassword", middleware.IsAuthorized(userHandler.ChangePassword)).Methods(http.MethodPost)
	if cfg.Server.RegisterEnabled {
		router.HandleFunc("/register", userHandler.Register).Methods(http.MethodPost)
	}

	// Tea Types
	teaTypeHandler := handlers.NewTeaTypeHandler(database)
	router.Handle("/types", middleware.IsAuthorized(teaTypeHandler.GetAllTeaTypes)).Methods(http.MethodGet)
	router.Handle("/types/teas", middleware.IsAuthorized(teaTypeHandler.GetAllTypesOfTea)).Methods(http.MethodGet)
	router.Handle("/type/{id:[0-9]+}", middleware.IsAuthorized(teaTypeHandler.GetTeaType)).Methods(http.MethodGet)
	router.Handle("/type", middleware.IsAuthorized(teaTypeHandler.CreateTeaType)).Methods(http.MethodPost)
	router.Handle("/type/{id:[0-9]+}", middleware.IsAuthorized(teaTypeHandler.DeleteTeaType)).Methods(http.MethodDelete)

	// Owners
	ownerHandler := handlers.NewOwnerHandler(database)
	router.Handle("/owners", middleware.IsAuthorized(ownerHandler.GetAllOwners)).Methods(http.MethodGet)
	router.Handle("/owners/teas", middleware.IsAuthorized(ownerHandler.GetAllOwnersTeas)).Methods(http.MethodGet)
	router.Handle("/owner/{id:[0-9]+}", middleware.IsAuthorized(ownerHandler.GetOwner)).Methods(http.MethodGet)
	router.Handle("/owner", middleware.IsAuthorized(ownerHandler.CreateOwner)).Methods(http.MethodPost)
	router.Handle("/owner/{id:[0-9]+}", middleware.IsAuthorized(ownerHandler.DeleteOwner)).Methods(http.MethodDelete)

	// Tea
	teaHandler := handlers.NewTeaHandler(database)
	router.Handle("/teas", middleware.IsAuthorized(teaHandler.GetAllTeas)).Methods(http.MethodGet)
	router.Handle("/tea/{id:[0-9]+}", middleware.IsAuthorized(teaHandler.GetTea)).Methods(http.MethodGet)
	router.Handle("/tea", middleware.IsAuthorized(teaHandler.CreateTea)).Methods(http.MethodPost)
	router.Handle("/tea/{id:[0-9]+}", middleware.IsAuthorized(teaHandler.DeleteTea)).Methods(http.MethodDelete)

	// Tea owners
	teaOwnerHandler := handlers.NewTeaOwnerHandler(database)
	router.Handle("/teas/owners", middleware.IsAuthorized(teaOwnerHandler.GetAllTeaOwners)).Methods(http.MethodGet)
	router.Handle("/tea/{id:[0-9]+}/owners", middleware.IsAuthorized(teaOwnerHandler.GetTeaOwners)).Methods(http.MethodGet)
	router.Handle("/tea/{id:[0-9]+}/owner", middleware.IsAuthorized(teaOwnerHandler.CreateTeaOwner)).Methods(http.MethodPost)
	router.Handle("/tea/{teaID:[0-9]+}/owner/{ownerID:[0-9]+}", middleware.IsAuthorized(teaOwnerHandler.DeleteTeaOwner)).Methods(http.MethodDelete)

	addr := ":" + cfg.Server.Port
	log.Fatal(http.ListenAndServe(addr, router))
}
