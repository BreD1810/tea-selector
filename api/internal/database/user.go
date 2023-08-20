package database

import (
	"fmt"

	"github.com/BreD1810/tea-selector/api/internal/models"
)

func (db *Database) createUserTable() error {
	creationString := `CREATE TABLE user (
							username TEXT NOT NULL UNIQUE PRIMARY KEY,
							password TEXT NOT NULL
						);`
	if _, err := db.DB.Exec(creationString); err != nil {
		return fmt.Errorf("error creating user table: %w", err)
	}

	return nil
}

// GetPasswordFromDatabase retrieves a users hashed password from the database.
func (db *Database) GetPasswordFromDatabase(user string) (string, error) {
	row := db.DB.QueryRow("SELECT password FROM user WHERE username=$1;", user)

	var password string
	if err := row.Scan(&password); err != nil {
		return "", err
	}

	return password, nil
}

// ChangePasswordInDatabase updates a user's password.
func (db *Database) ChangePasswordInDatabase(username string, password string) error {
	_, err := db.DB.Exec("UPDATE user SET password=$1 WHERE username=$2;", password, username)
	return err
}

// CreateUserInDatabase creates a user in the database using a username and pre-hashed password
func (db *Database) CreateUserInDatabase(user models.UserLogin) error {
	_, err := db.DB.Exec("INSERT INTO user(username, password) VALUES ($1, $2)", user.Username, user.Password)

	return err
}
