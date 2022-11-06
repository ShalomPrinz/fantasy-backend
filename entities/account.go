package entities

type Account struct {
	ID       string   `json:"id"`
	Nickname string   `json:"nickname"`
	Team     []string `json:"team"`
}

// My users collection
type AddAccount struct {
	Nickname string   `json:"nickname"`
	Team     []string `json:"team"`
}

// Firebase Auth database, following auth.UserToCreate struct
type AddUser struct {
	FullName string
	Nickname string
	Email    string
	Password string
}
