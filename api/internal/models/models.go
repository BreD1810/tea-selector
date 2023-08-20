package models

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

// A UserLogin stores the username and password for a user
type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
