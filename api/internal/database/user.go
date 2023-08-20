package database

import (
	"github.com/BreD1810/tea-selector/api/internal/models"
)

func createUserTable() {
	creationString := `CREATE TABLE user (
							username TEXT NOT NULL UNIQUE PRIMARY KEY,
							password TEXT NOT NULL
						);`
	_, err := DB.Exec(creationString)
	checkError("creating user table", err)
}

// GetPasswordFromDatabase retrieves a users hashed password from the database.
func GetPasswordFromDatabase(user string) (string, error) {
	row := DB.QueryRow("SELECT password FROM user WHERE username=$1;", user)

	var password string
	if err := row.Scan(&password); err != nil {
		return "", err
	}

	return password, nil
}

// ChangePasswordInDatabase updates a user's password.
func ChangePasswordInDatabase(username string, password string) error {
	if _, err := DB.Exec("UPDATE user SET password=$1 WHERE username=$2;", password, username); err != nil {
		return err
	}
	return nil
}

// CreateUserInDatabase creates a user in the database using a username and pre-hashed password
func CreateUserInDatabase(user models.UserLogin) error {
	if _, err := DB.Exec("INSERT INTO user(username, password) VALUES ($1, $2)", user.Username, user.Password); err != nil {
		return err
	}

	return nil
}
