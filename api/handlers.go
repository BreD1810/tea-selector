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

// A TeaWithOwners details a relationship between a tea Owner and a Tea.
type TeaWithOwners struct {
	Tea    Tea     `json:"tea"`
	Owners []Owner `json:"owners"`
}

// A TypeWithTeas details all the teas of a single type.
type TypeWithTeas struct {
	Type TeaType `json:"type"`
	Teas []Tea   `json:"teas"`
}

// A OwnerWithTeas details all the teas for an owner.
type OwnerWithTeas struct {
	Owner Owner `json:"owner"`
	Teas  []Tea `json:"teas"`
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

// GetAllTeasFunc points to a function to get all teas in the database. Useful for mocking.
var GetAllTeasFunc = GetAllTeasFromDatabase

func getAllTeasHandler(w http.ResponseWriter, r *http.Request) {
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
var GetTeaFunc = GetTeaFromDatabase

func getTeaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Failed to get tea with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusBadRequest, "Invalid tea ID")
		return
	}
	log.Printf("Received request \"GET /tea/%d\"\n", id)

	tea := Tea{ID: id}

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
var CreateTeaFunc = CreateTeaInDatabase

func createTeaHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(`Received request "POST /tea"`)

	var tea Tea
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
var DeleteTeaFunc = DeleteTeaFromDatabase

func deleteTeaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Failed to delete tea with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusBadRequest, "Invalid tea ID")
		return
	}
	log.Printf("Received request \"DELETE /tea/%d\"\n", id)

	tea := Tea{ID: id}

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

// GetTeaOwnersFunc points to a function to get owners of a tea from the database. Useful for mocking.
var GetTeaOwnersFunc = GetTeaOwnersFromDatabase

func getTeaOwnersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Failed to get owners of tea with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusBadRequest, "Invalid tea ID")
		return
	}
	log.Printf("Received request \"GET /tea/%d/owners\"\n", id)

	tea := Tea{ID: id}

	owners, err := GetTeaOwnersFunc(&tea)
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

// GetAllTeaOwnersFunc points to a function to get a list of teas and their owners. Useful for mocking.
var GetAllTeaOwnersFunc = GetAllTeaOwnersFromDatabase

func getAllTeaOwnersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(`Received request "GET /tea/owners"`)

	teasWithOwners, err := GetAllTeaOwnersFunc()
	if err != nil {
		log.Printf("Error retrieving all teas with owners: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println("Successfully handled request to see all teas with owners")
	respondWithJSON(w, http.StatusOK, teasWithOwners)
}

// CreateTeaOwnerFunc points to a function to add an owner to a tea in the database. Useful for mocking.
var CreateTeaOwnerFunc = CreateTeaOwnerInDatabase

func createTeaOwnerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Failed to get owners of tea with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusBadRequest, "Invalid tea ID")
		return
	}
	log.Printf("Received request \"POST /tea/%d/owner\n\"", id)

	var owner Owner
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&owner); err != nil {
		log.Printf("Failed to create new owner of tea with ID: %d\n", id)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	tea, err := CreateTeaOwnerFunc(id, &owner)
	if err != nil {
		log.Printf("Error creating owner for tea with ID: %d\n\t Error: %s\n", id, err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Created new owner for tea. teaID: %d, ownerID: %d\n", id, owner.ID)
	respondWithJSON(w, http.StatusCreated, tea)
}

// DeleteTeaOwnerFunc points to a function to delete an owner from a tea in the database. Useful for mocking.
var DeleteTeaOwnerFunc = DeleteTeaOwnerFromDatabase

func deleteTeaOwnerHandler(w http.ResponseWriter, r *http.Request) {
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

	tea := Tea{ID: teaID}
	owner := Owner{ID: ownerID}

	if err := DeleteTeaOwnerFunc(&tea, &owner); err != nil {
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

// GetAllTypesTeasFunc Func gets all teas by type from the database. Useful for mocking.
var GetAllTypesTeasFunc = GetAllTypesTeasFromDatabase

func getAllTeasTypesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(`Received request "GET /types/teas"`)

	typesWithTeas, err := GetAllTypesTeasFunc()
	if err != nil {
		log.Printf("Error retrieving all types with teas: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println("Successfully handled request to see all teas by type")
	respondWithJSON(w, http.StatusOK, typesWithTeas)
}

// GetAllOwnersTeasFunc gets a list of owners, and their teas. Useful for mocking
var GetAllOwnersTeasFunc = GetAllOwnersTeasFromDatabase

func getAllOwnersTeasHandler(w http.ResponseWriter, r *http.Request) {
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
