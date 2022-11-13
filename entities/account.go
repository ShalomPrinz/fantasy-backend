package entities

// My users collection
type Account struct {
	ID       string   `json:"id"`
	Nickname string   `json:"nickname"`
	Team     []string `json:"team"`
}

type AddAccount struct {
	Nickname string   `json:"nickname"`
	Team     []string `json:"team"`
}

var AccountAttributes = [...]string{"Nickname", "Team"}

// Firebase Auth database, following auth.UserToCreate struct
type AddUser struct {
	FullName string
	Nickname string
	Email    string
	Password string
}

type LoginUser struct {
	IdToken string
}
