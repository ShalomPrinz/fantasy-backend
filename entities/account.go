package entities

type Account struct {
	ID   string   `json:"id"`
	Team []Player `json:"team"`
}

// My users collection
type AddAccount struct {
	Team []Player `json:"team"`
}

// Firebase Auth database, following auth.UserToCreate struct
type AddUser struct {
	DisplayName string
	Email       string
	Password    string
}
