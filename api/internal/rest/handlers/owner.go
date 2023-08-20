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

type OwnerHandler interface {
	GetAllOwners(w http.ResponseWriter, r *http.Request)
	GetOwner(w http.ResponseWriter, r *http.Request)
	GetAllOwnersTeas(w http.ResponseWriter, r *http.Request)
	CreateOwner(w http.ResponseWriter, r *http.Request)
	DeleteOwner(w http.ResponseWriter, r *http.Request)
}

type handlerOfOwner struct {
	db *database.Database
}

func NewOwnerHandler(db *database.Database) OwnerHandler {
	return &handlerOfOwner{db: db}
}

func (h *handlerOfOwner) GetAllOwners(w http.ResponseWriter, r *http.Request) {
	log.Println(`Received request "GET /owners"`)

	owners, err := h.db.GetAllOwnersFromDatabase()
	if err != nil {
		log.Printf("Error retrieving all owners: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println("Successfully handled request to see all owners")
	respondWithJSON(w, http.StatusOK, owners)
}

func (h *handlerOfOwner) GetOwner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Failed to get owner with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusBadRequest, "Invalid owner ID")
		return
	}
	log.Printf("Received request \"GET /owner/%d\"\n", id)

	owner := models.Owner{ID: id}

	if err := h.db.GetOwnerFromDatabase(&owner); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Failed to get owner as ID didn't exist. ID: %d\n", id)
			respondWithError(w, http.StatusInternalServerError, "ID does not exist in database")
			return
		}
		log.Printf("Failed to get owner with id: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Got owner with ID: %d\n", id)
	respondWithJSON(w, http.StatusOK, owner)
}

func (h *handlerOfOwner) GetAllOwnersTeas(w http.ResponseWriter, r *http.Request) {
	log.Println(`Receieved request "GET /owners/teas"`)

	ownersWithTeas, err := h.db.GetAllOwnersTeasFromDatabase()
	if err != nil {
		log.Printf("Error retrieving all teas for all owners: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println("Successfully handled request to see all teas for each owner")
	respondWithJSON(w, http.StatusOK, ownersWithTeas)
}

func (h *handlerOfOwner) CreateOwner(w http.ResponseWriter, r *http.Request) {
	log.Println(`Received request "POST /owner"`)

	var owner models.Owner
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&owner); err != nil {
		log.Printf("Failed to create new owner: %s\n", owner.Name)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := h.db.CreateOwnerInDatabase(&owner); err != nil {
		log.Printf("Error creating owner: %s\n\t Error: %s\n", owner.Name, err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Created new owner. ID: %d, Name: %s\n", owner.ID, owner.Name)
	respondWithJSON(w, http.StatusCreated, owner)
}

func (h *handlerOfOwner) DeleteOwner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Failed to delete owner with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusBadRequest, "Invalid owner ID")
		return
	}
	log.Printf("Received request \"DELETE /owner/%d\"\n", id)

	owner := models.Owner{ID: id}

	if err := h.db.DeleteOwnerFromDatabase(&owner); err != nil {
		if err.Error() == "sql: Rows are closed" {
			log.Printf("Failed to delete owner as ID didn't exist. ID: %d\n", id)
			respondWithError(w, http.StatusInternalServerError, "ID does not exist in database")
			return
		}
		log.Printf("Failed to delete owner with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Deleted owner with ID: %d\n", id)
	respondWithJSON(w, http.StatusOK, map[string]string{"name": owner.Name, "result": "success"})
}
