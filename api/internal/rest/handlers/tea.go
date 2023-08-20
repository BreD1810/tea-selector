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

// GetAllTeasFunc points to a function to get all teas in the database. Useful for mocking.
var GetAllTeasFunc = database.GetAllTeasFromDatabase

func GetAllTeasHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(`Received request "GET /teas"`)

	teas, err := GetAllTeasFunc()
	if err != nil {
		log.Printf("Error retrieving all teas: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println("Successfully handled request to see all teas")
	respondWithJSON(w, http.StatusOK, teas)
}

// GetTeaFunc points to a function to get information about a tea from the database. Useful for mocking.
var GetTeaFunc = database.GetTeaFromDatabase

func GetTeaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Failed to get tea with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusBadRequest, "Invalid tea ID")
		return
	}
	log.Printf("Received request \"GET /tea/%d\"\n", id)

	tea := models.Tea{ID: id}

	if err := GetTeaFunc(&tea); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Failed to get tea as ID didn't exist. ID: %d\n", id)
			respondWithError(w, http.StatusInternalServerError, "ID does not exist in database")
			return
		}
		log.Printf("Failed to get tea with id: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Got tea with ID: %d\n", id)
	respondWithJSON(w, http.StatusOK, tea)
}

// CreateTeaFunc points to a function to create a tea in the database. Useful for mocking.
var CreateTeaFunc = database.CreateTeaInDatabase

func CreateTeaHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(`Received request "POST /tea"`)

	var tea models.Tea
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&tea); err != nil {
		log.Printf("Failed to create new tea: %s\n", tea.Name)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := CreateTeaFunc(&tea); err != nil {
		log.Printf("Error creating tea: %s\n\t Error: %s\n", tea.Name, err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Created new tea. ID: %d, Name: %q, Type: %q\n", tea.ID, tea.Name, tea.TeaType.Name)
	respondWithJSON(w, http.StatusCreated, tea)
}

// DeleteTeaFunc points to a function to delete a tea from the database. Useful for mocking.
var DeleteTeaFunc = database.DeleteTeaFromDatabase

func DeleteTeaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Failed to delete tea with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusBadRequest, "Invalid tea ID")
		return
	}
	log.Printf("Received request \"DELETE /tea/%d\"\n", id)

	tea := models.Tea{ID: id}

	if err := DeleteTeaFunc(&tea); err != nil {
		if err.Error() == "sql: Rows are closed" {
			log.Printf("Failed to delete tea as ID didn't exist. ID: %d\n", id)
			respondWithError(w, http.StatusInternalServerError, "ID does not exist in database")
			return
		}
		log.Printf("Failed to delete tea with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Deleted tea with ID: %d\n", id)
	respondWithJSON(w, http.StatusOK, map[string]string{"name": tea.Name, "result": "success"})
}
