package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// A TeaType gives the ID and name for a type of tea.
type TeaType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// A Tea details a tea within the system, with an ID, name and type of the tea.
type Tea struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	TeaType TeaType `json:"type"`
}

// An Owner is someone who has some tea that is in the system.
type Owner struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// A TeaOwners details a relationship between a tea Owner and a Tea.
type TeaOwners struct {
	TeaID   int `json:"teaID"`
	OwnerID int `json:"ownerID"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// GetAllTeaTypesFunc points to the function to get all tea typevalues. Useful for mocking
var GetAllTeaTypesFunc = GetAllTeaTypesFromDatabase

func getAllTeaTypesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(`Received request "GET /types"`)

	types, err := GetAllTeaTypesFunc()
	if err != nil {
		log.Printf("Error retrieving all tea types: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println("Successfully handled request to see all tea types")
	respondWithJSON(w, http.StatusOK, types)
}

// GetTeaTypeFunc points to the function to get information about a tea type. Useful for mocking.
var GetTeaTypeFunc = GetTeaTypeFromDatabase

func getTeaTypeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Failed to get tea type with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusBadRequest, "Invalid Tea Type ID")
		return
	}
	log.Printf("Received request \"GET /type/%d\"\n", id)

	teaType := TeaType{ID: id}

	if err := GetTeaTypeFunc(&teaType); err != nil {
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

// CreateTeaTypeFunc points to the function to create a new type of tea in the database. Useful for mocking.
var CreateTeaTypeFunc = CreateTeaTypeInDatabase

func createTeaTypeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(`Received request "POST /type"`)

	var teaType TeaType
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&teaType); err != nil {
		log.Printf("Failed to create new tea type: %s\n", teaType.Name)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := CreateTeaTypeFunc(&teaType); err != nil {
		log.Printf("Error creating tea type: %s\n\t Error: %s\n", teaType.Name, err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Created new tea type. ID: %d, Name: %s\n", teaType.ID, teaType.Name)
	respondWithJSON(w, http.StatusCreated, teaType)
}

// DeleteTeaTypeFunc points to the function to delete a type of tea in the database. Useful for mocking.
var DeleteTeaTypeFunc = DeleteTeaTypeInDatabase

func deleteTeaTypeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Failed to delete tea type with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusBadRequest, "Invalid Tea Type ID")
		return
	}
	log.Printf("Received request \"DELETE /type/%d\"\n", id)

	teaType := TeaType{ID: id}

	if err := DeleteTeaTypeFunc(&teaType); err != nil {
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

// GetAllOwnersFunc points to a function to get all owners in the database. Useful for mocking.
var GetAllOwnersFunc = GetAllOwnersFromDatabase

func getAllOwnersHandler(w http.ResponseWriter, r *http.Request) {
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
var GetOwnerFunc = GetOwnerFromDatabase

func getOwnerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Failed to get owner with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusBadRequest, "Invalid owner ID")
		return
	}
	log.Printf("Received request \"GET /owner/%d\"\n", id)

	owner := Owner{ID: id}

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

// CreateOwnerFunc points to a function that creates an owner in the database. Useful for mocking.
var CreateOwnerFunc = CreateOwnerInDatabase

func createOwnerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(`Received request "POST /owner"`)

	var owner Owner
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
var DeleteOwnerFunc = DeleteOwnerFromDatabase

func deleteOwnerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Failed to delete owner with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusBadRequest, "Invalid owner ID")
		return
	}
	log.Printf("Received request \"DELETE /owner/%d\"\n", id)

	owner := Owner{ID: id}

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
