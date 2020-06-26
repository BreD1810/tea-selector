package main

import (
	"errors"
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
	claims["users"] = user
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(signingKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
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
