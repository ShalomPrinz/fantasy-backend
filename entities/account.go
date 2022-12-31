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

type DetailedAccount struct {
	Entity   `mapstructure:",squash"`
	Leagues  []LeagueInfo `json:"leagues"`
	Nickname string       `json:"nickname"`
	Team     []Player     `json:"team"`
}

// Firebase Auth database, following auth.UserToCreate struct
type AddUser struct {
	FullName string `binding:"required"`
	Nickname string `binding:"required"`
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

type LoginUser struct {
	IdToken string
}
