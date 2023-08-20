package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/BreD1810/tea-selector/api/internal/database"
	"github.com/BreD1810/tea-selector/api/internal/models"
	"github.com/gorilla/mux"
)

type TeaOwnerHandler interface {
	GetAllTeaOwners(w http.ResponseWriter, r *http.Request)
	GetTeaOwners(w http.ResponseWriter, r *http.Request)
	CreateTeaOwner(w http.ResponseWriter, r *http.Request)
	DeleteTeaOwner(w http.ResponseWriter, r *http.Request)
}

type handlerOfTeaOwner struct {
	db *database.Database
}

func NewTeaOwnerHandler(db *database.Database) TeaOwnerHandler {
	return &handlerOfTeaOwner{db: db}
}

func (h *handlerOfTeaOwner) GetTeaOwners(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Failed to get owners of tea with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusBadRequest, "Invalid tea ID")
		return
	}
	log.Printf("Received request \"GET /tea/%d/owners\"\n", id)

	tea := models.Tea{ID: id}

	owners, err := h.db.GetTeaOwnersFromDatabase(&tea)
	if err != nil {
		if err.Error() == "sql: Rows are closed" {
			log.Printf("Failed to get tea owners as tea ID didn't exist. ID: %d\n", id)
			respondWithError(w, http.StatusInternalServerError, "ID does not exist in database")
			return
		}
		log.Printf("Failed to get tea owners with tea ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Got owners of tea with ID: %d\n", id)
	respondWithJSON(w, http.StatusOK, owners)
}

func (h *handlerOfTeaOwner) GetAllTeaOwners(w http.ResponseWriter, r *http.Request) {
	log.Println(`Received request "GET /tea/owners"`)

	teasWithOwners, err := h.db.GetAllTeaOwnersFromDatabase()
	if err != nil {
		log.Printf("Error retrieving all teas with owners: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println("Successfully handled request to see all teas with owners")
	respondWithJSON(w, http.StatusOK, teasWithOwners)
}

func (h *handlerOfTeaOwner) CreateTeaOwner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Failed to get owners of tea with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusBadRequest, "Invalid tea ID")
		return
	}
	log.Printf("Received request \"POST /tea/%d/owner\n\"", id)

	var owner models.Owner
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&owner); err != nil {
		log.Printf("Failed to create new owner of tea with ID: %d\n", id)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	tea, err := h.db.CreateTeaOwnerInDatabase(id, &owner)
	if err != nil {
		log.Printf("Error creating owner for tea with ID: %d\n\t Error: %s\n", id, err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Created new owner for tea. teaID: %d, ownerID: %d\n", id, owner.ID)
	respondWithJSON(w, http.StatusCreated, tea)
}

func (h *handlerOfTeaOwner) DeleteTeaOwner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teaID, err := strconv.Atoi(vars["teaID"])
	if err != nil {
		log.Printf("Failed to delete owner of tea with teaID: %d\n Error: %v\n", teaID, err)
		respondWithError(w, http.StatusBadRequest, "Invalid tea ID")
		return
	}
	ownerID, err := strconv.Atoi(vars["ownerID"])
	if err != nil {
		log.Printf("Failed to delete owner of tea with ownerID: %d\n Error: %v\n", ownerID, err)
		respondWithError(w, http.StatusBadRequest, "Invalid owner ID")
		return
	}
	log.Printf("Received request \"DELETE /tea/%d/owner/%d\"\n", teaID, ownerID)

	tea := models.Tea{ID: teaID}
	owner := models.Owner{ID: ownerID}

	if err := h.db.DeleteTeaOwnerFromDatabase(&tea, &owner); err != nil {
		if err.Error() == "sql: Rows are closed" {
			log.Println("Failed to delete tea owner as relationship doesn't exist")
			respondWithError(w, http.StatusInternalServerError, "Relationship does not exist in database")
			return
		}
		log.Printf("Failed to delete tea owner. Error: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Deleted tea owner. teaID: %d \t ownerID: %d\n", teaID, ownerID)
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
