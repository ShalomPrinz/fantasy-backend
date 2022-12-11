package entities

// My users collection
type Account struct {
	Entity   `mapstructure:",squash"`
	Leagues  []string `json:"leagues"`
	Nickname string   `json:"nickname"`
	Team     []string `json:"team"`
}

type AddAccount struct {
	Nickname string   `json:"nickname"`
	Team     []string `json:"team"`
}

type InsertAccount struct {
	Leagues  []string `json:"leagues"`
	Nickname string   `json:"nickname"`
	Team     []string `json:"team"`
}

type Member struct {
	Entity   `mapstructure:",squash"`
	Nickname string   `json:"nickname"`
	Team     []Player `json:"team"`
}

type DetailedAccount struct {
	Entity   `mapstructure:",squash"`
	Leagues  []League `json:"leagues"`
	Nickname string   `json:"nickname"`
	Team     []Player `json:"team"`
}

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
