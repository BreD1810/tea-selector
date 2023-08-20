package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/BreD1810/tea-selector/api/internal/database"
	"github.com/BreD1810/tea-selector/api/internal/models"
	"github.com/BreD1810/tea-selector/api/internal/rest/middleware"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// A NewPasswordRequest stores the old and new password for a reset request.
type NewPasswordRequest struct {
	OldPassword string `json:"old"`
	NewPassword string `json:"new"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(`Received request "POST /login"`)

	var userLogin models.UserLogin
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userLogin); err != nil {
		log.Println("Failed to extract username and password")
		respondWithError(w, http.StatusBadRequest, "Bad request body")
		return
	}
	defer r.Body.Close()

	userLogin.Username = strings.ToLower(userLogin.Username)

	// Retrieve from DB
	storedPassword, err := database.GetPasswordFromDatabase(userLogin.Username)
	if err != nil {
		log.Printf("Failed to get password from database for user %q\n", userLogin.Username)
		respondWithError(w, http.StatusBadRequest, "User doesn't exist")
		return
	}

	// Compare hash with sent password.
	storedPasswordBytes := []byte(storedPassword)
	passwordBytes := []byte(userLogin.Password)
	if err := bcrypt.CompareHashAndPassword(storedPasswordBytes, passwordBytes); err != nil {
		log.Printf("Password incorrect for user %q\n", userLogin.Username)
		respondWithError(w, http.StatusBadRequest, "Incorrect password")
		return
	}

	validToken, err := middleware.GenerateJWT(userLogin.Username)
	if err != nil {
		log.Printf("Error generating token: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	log.Printf("Successfully logged in %q\n", userLogin.Username)
	respondWithJSON(w, http.StatusOK, map[string]string{"token": validToken})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(`Received request "POST /register"`)

	var userLogin models.UserLogin
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userLogin); err != nil {
		log.Println(`Failed to extract username and password`)
		respondWithError(w, http.StatusBadRequest, "Bad request body")
		return
	}
	defer r.Body.Close()

	userLogin.Username = strings.ToLower(userLogin.Username)
	passwordBytes := []byte(userLogin.Password)

	hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	if err != nil {
		log.Println(`Error hashing password`)
		respondWithError(w, http.StatusInternalServerError, "Unable to create user")
		return
	}

	userLogin.Password = string(hash)
	if err := database.CreateUserInDatabase(userLogin); err != nil {
		log.Printf("Error creating user: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	validToken, err := middleware.GenerateJWT(userLogin.Username)
	if err != nil {
		log.Printf("Error generating token: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	log.Printf("Successfully registered and logged in %q\n", userLogin.Username)
	respondWithJSON(w, http.StatusOK, map[string]string{"token": validToken})
}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(`Received request "POST /changepassword"`)

	var newPasswordBody NewPasswordRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newPasswordBody); err != nil {
		log.Println("Failed to extract old and new password")
		respondWithError(w, http.StatusBadRequest, "Bad request body")
		return
	}
	defer r.Body.Close()

	// Get the username of the JWT
	username, err := middleware.GetJWTUser(r.Header["Token"][0])
	if err != nil {
		log.Println("Unable to change password")
		respondWithError(w, http.StatusInternalServerError, "Error changing user password")
		return
	}

	// Retrieve old password from DB
	storedPassword, err := database.GetPasswordFromDatabase(username)
	if err != nil {
		log.Printf("Failed to get password from database for user %q\n", username)
		respondWithError(w, http.StatusBadRequest, "User doesn't exist")
		return
	}

	// Compare old hash with sent old password.
	storedPasswordBytes := []byte(storedPassword)
	passwordBytes := []byte(newPasswordBody.OldPassword)
	if err := bcrypt.CompareHashAndPassword(storedPasswordBytes, passwordBytes); err != nil {
		log.Printf("Password incorrect for user %q\n", username)
		respondWithError(w, http.StatusBadRequest, "Incorrect password")
		return
	}

	// Hash new password
	newPasswordBytes := []byte(newPasswordBody.NewPassword)
	hash, err := bcrypt.GenerateFromPassword(newPasswordBytes, bcrypt.MinCost)
	if err != nil {
		log.Println(`Error hashing password`)
		respondWithError(w, http.StatusInternalServerError, "Unable to create user")
		return
	}
	newPassword := string(hash)

	if err := database.ChangePasswordInDatabase(username, newPassword); err != nil {
		log.Printf("Error changing password: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Error changing user password")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// GetAllTeaTypesFunc points to the function to get all tea typevalues. Useful for mocking
var GetAllTeaTypesFunc = database.GetAllTeaTypesFromDatabase

func GetAllTeaTypesHandler(w http.ResponseWriter, r *http.Request) {
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
var GetTeaTypeFunc = database.GetTeaTypeFromDatabase

func GetTeaTypeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Failed to get tea type with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusBadRequest, "Invalid Tea Type ID")
		return
	}
	log.Printf("Received request \"GET /type/%d\"\n", id)

	teaType := models.TeaType{ID: id}

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
var CreateTeaTypeFunc = database.CreateTeaTypeInDatabase

func CreateTeaTypeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(`Received request "POST /type"`)

	var teaType models.TeaType
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
var DeleteTeaTypeFunc = database.DeleteTeaTypeInDatabase

func DeleteTeaTypeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Failed to delete tea type with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusBadRequest, "Invalid Tea Type ID")
		return
	}
	log.Printf("Received request \"DELETE /type/%d\"\n", id)

	teaType := models.TeaType{ID: id}

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

// GetTeaOwnersFunc points to a function to get owners of a tea from the database. Useful for mocking.
var GetTeaOwnersFunc = database.GetTeaOwnersFromDatabase

func GetTeaOwnersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Failed to get owners of tea with ID: %d\n Error: %v\n", id, err)
		respondWithError(w, http.StatusBadRequest, "Invalid tea ID")
		return
	}
	log.Printf("Received request \"GET /tea/%d/owners\"\n", id)

	tea := models.Tea{ID: id}

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
var GetAllTeaOwnersFunc = database.GetAllTeaOwnersFromDatabase

func GetAllTeaOwnersHandler(w http.ResponseWriter, r *http.Request) {
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
var CreateTeaOwnerFunc = database.CreateTeaOwnerInDatabase

func CreateTeaOwnerHandler(w http.ResponseWriter, r *http.Request) {
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
var DeleteTeaOwnerFunc = database.DeleteTeaOwnerFromDatabase

func DeleteTeaOwnerHandler(w http.ResponseWriter, r *http.Request) {
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
var GetAllTypesTeasFunc = database.GetAllTypesTeasFromDatabase

func GetAllTeasTypesHandler(w http.ResponseWriter, r *http.Request) {
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
