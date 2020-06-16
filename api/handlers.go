package main

import (
	"encoding/json"
	"net/http"
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

// GetAllTeaTypesFunc points to the function to get all tea typevalues. Useful for mocking
var GetAllTeaTypesFunc = GetAllTeaTypesFromDatabase

func getAllTeaTypesHandler(w http.ResponseWriter, r *http.Request) {
	types := GetAllTeaTypesFunc()
	json.NewEncoder(w).Encode(types)
}
