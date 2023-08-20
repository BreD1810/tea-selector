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

// GetAllOwnersFunc points to a function to get all owners in the database. Useful for mocking.
var GetAllOwnersFunc = database.GetAllOwnersFromDatabase

func GetAllOwnersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(`Received request "GET /owners"`)

	owners, err := GetAllOwnersFunc()
	if err != nil {
		log.Printf("Error retrieving all owners: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println("Successfully handled request to see all owners")
	respondWithJSON(w, http.StatusOK, owners)
}

// GetOwnerFunc points to a function to get information about an owner from the database. Useful for mocking
var GetOwnerFunc = database.GetOwnerFromDatabase

func GetOwnerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Failed to get owner with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusBadRequest, "Invalid owner ID")
		return
	}
	log.Printf("Received request \"GET /owner/%d\"\n", id)

	owner := models.Owner{ID: id}

	if err := GetOwnerFunc(&owner); err != nil {
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

// GetAllOwnersTeasFunc gets a list of owners, and their teas. Useful for mocking
var GetAllOwnersTeasFunc = database.GetAllOwnersTeasFromDatabase

func GetAllOwnersTeasHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(`Receieved request "GET /owners/teas"`)

	ownersWithTeas, err := GetAllOwnersTeasFunc()
	if err != nil {
		log.Printf("Error retrieving all teas for all owners: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println("Successfully handled request to see all teas for each owner")
	respondWithJSON(w, http.StatusOK, ownersWithTeas)
}

// CreateOwnerFunc points to a function that creates an owner in the database. Useful for mocking.
var CreateOwnerFunc = database.CreateOwnerInDatabase

func CreateOwnerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(`Received request "POST /owner"`)

	var owner models.Owner
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&owner); err != nil {
		log.Printf("Failed to create new owner: %s\n", owner.Name)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := CreateOwnerFunc(&owner); err != nil {
		log.Printf("Error creating owner: %s\n\t Error: %s\n", owner.Name, err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Created new owner. ID: %d, Name: %s\n", owner.ID, owner.Name)
	respondWithJSON(w, http.StatusCreated, owner)
}

// DeleteOwnerFunc points to a function to delete an owner from the database. Useful for mocking.
var DeleteOwnerFunc = database.DeleteOwnerFromDatabase

func DeleteOwnerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Failed to delete owner with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusBadRequest, "Invalid owner ID")
		return
	}
	log.Printf("Received request \"DELETE /owner/%d\"\n", id)

	owner := models.Owner{ID: id}

	if err := DeleteOwnerFunc(&owner); err != nil {
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
