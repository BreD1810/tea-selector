package main

import (
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
	types, err := GetAllTeaTypesFunc()
	if err != nil {
		log.Printf("Error retrieving all tea types: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println("Successfully handled request to see all tea types")
	respondWithJSON(w, http.StatusOK, types)
}

// CreateTeaTypeFunc points to the function to create a new type of tea in the database. Useful for mocking.
var CreateTeaTypeFunc = CreateTeaTypeInDatabase

func createTeaTypeHandler(w http.ResponseWriter, r *http.Request) {
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

	teaType := TeaType{ID: id}

	if err := DeleteTeaTypeFunc(&teaType); err != nil {
		if err.Error() == "sql: Rows are closed" {
			log.Printf("Failed to delete tea type as ID didn't exist. ID: %d\n", id)
			respondWithError(w, http.StatusInternalServerError, "ID does not exist in database")
			return
		}
		log.Printf("Failed to delete tea type with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("Delete tea type with ID: %d\n", id)
	respondWithJSON(w, http.StatusOK, map[string]string{"name": teaType.Name, "result": "success"})
}
