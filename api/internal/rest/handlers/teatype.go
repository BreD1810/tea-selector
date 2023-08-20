package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/BreD1810/tea-selector/api/internal/database"
	"github.com/BreD1810/tea-selector/api/internal/models"
	"github.com/gorilla/mux"
)

type TeaTypeHandler interface {
	GetAllTeaTypes(w http.ResponseWriter, r *http.Request)
	GetAllTypesOfTea(w http.ResponseWriter, r *http.Request)
	GetTeaType(w http.ResponseWriter, r *http.Request)
	CreateTeaType(w http.ResponseWriter, r *http.Request)
	DeleteTeaType(w http.ResponseWriter, r *http.Request)
}

type handlerOfTeaType struct {
	db *database.Database
}

func NewTeaTypeHandler(db *database.Database) TeaTypeHandler {
	return &handlerOfTeaType{db: db}
}

func (h *handlerOfTeaType) GetAllTeaTypes(w http.ResponseWriter, r *http.Request) {
	log.Println(`Received request "GET /types"`)

	types, err := h.db.GetAllTeaTypesFromDatabase()
	if err != nil {
		log.Printf("Error retrieving all tea types: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println("Successfully handled request to see all tea types")
	respondWithJSON(w, http.StatusOK, types)
}

func (h *handlerOfTeaType) GetTeaType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Failed to get tea type with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusBadRequest, "Invalid Tea Type ID")
		return
	}
	log.Printf("Received request \"GET /type/%d\"\n", id)

	teaType := models.TeaType{ID: id}

	if err := h.db.GetTeaTypeFromDatabase(&teaType); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Failed to get tea type as ID didn't exist. ID: %d\n", id)
			respondWithError(w, http.StatusInternalServerError, "ID does not exist in database")
			return
		}
		log.Printf("Failed to get tea type with id: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Got tea with type with ID: %d\n", id)
	respondWithJSON(w, http.StatusOK, teaType)
}

func (h *handlerOfTeaType) CreateTeaType(w http.ResponseWriter, r *http.Request) {
	log.Println(`Received request "POST /type"`)

	var teaType models.TeaType
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&teaType); err != nil {
		log.Printf("Failed to create new tea type: %s\n", teaType.Name)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := h.db.CreateTeaTypeInDatabase(&teaType); err != nil {
		log.Printf("Error creating tea type: %s\n\t Error: %s\n", teaType.Name, err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Created new tea type. ID: %d, Name: %s\n", teaType.ID, teaType.Name)
	respondWithJSON(w, http.StatusCreated, teaType)
}

func (h *handlerOfTeaType) DeleteTeaType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Failed to delete tea type with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusBadRequest, "Invalid Tea Type ID")
		return
	}
	log.Printf("Received request \"DELETE /type/%d\"\n", id)

	teaType := models.TeaType{ID: id}

	if err := h.db.DeleteTeaTypeInDatabase(&teaType); err != nil {
		if err.Error() == "sql: Rows are closed" {
			log.Printf("Failed to delete tea type as ID didn't exist. ID: %d\n", id)
			respondWithError(w, http.StatusInternalServerError, "ID does not exist in database")
			return
		}
		log.Printf("Failed to delete tea type with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Delete tea type with ID: %d\n", id)
	respondWithJSON(w, http.StatusOK, map[string]string{"name": teaType.Name, "result": "success"})
}

func (h *handlerOfTeaType) GetAllTypesOfTea(w http.ResponseWriter, r *http.Request) {
	log.Println(`Received request "GET /types/teas"`)

	typesWithTeas, err := h.db.GetAllTypesTeasFromDatabase()
	if err != nil {
		log.Printf("Error retrieving all types with teas: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println("Successfully handled request to see all teas by type")
	respondWithJSON(w, http.StatusOK, typesWithTeas)
}
