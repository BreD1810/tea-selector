package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var signingKey []byte

// SetSigningKey lets you set the signing key
func SetSigningKey(key string) {
	signingKey = []byte(key)
}

// GenerateJWT generates a JWT token
func GenerateJWT(user string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = user
	claims["exp"] = time.Now().AddDate(0, 0, 7).Unix() // 7 days expiry

	tokenString, err := token.SignedString(signingKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetJWTUser gets the user from a JWT token
func GetJWTUser(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Error parsing JWT")
		}
		return signingKey, nil
	})
	if err != nil {
		return "", nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && token.Valid {
		return "", errors.New("Couldn't extract claims from token")
	}

	username := fmt.Sprintf("%v", claims["user"])

	return username, nil
}

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("Error parsing JWT")
				}
				return signingKey, nil
			})

			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Not Authorized")
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			log.Printf("Not authorized")
			respondWithError(w, http.StatusBadRequest, "Not Authorized")
		}
	})
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
