package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/BreD1810/tea-selector/api/internal/database"
	"github.com/BreD1810/tea-selector/api/internal/models"
	"github.com/BreD1810/tea-selector/api/internal/rest/middleware"
	"golang.org/x/crypto/bcrypt"
)

// A NewPasswordRequest stores the old and new password for a reset request.
type NewPasswordRequest struct {
	OldPassword string `json:"old"`
	NewPassword string `json:"new"`
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
