package entities

// My users collection
type Account struct {
	Entity   `mapstructure:",squash"`
	Inbox    []Message `json:"inbox"`
	Leagues  []string  `json:"leagues"`
	Username string    `json:"username"`
	Team     []string  `json:"team"`
}

type AddAccount struct {
	Username string   `json:"username"`
	Team     []string `json:"team"`
}

type InsertAccount struct {
	Inbox    []Message `json:"inbox"`
	Leagues  []string  `json:"leagues"`
	Username string    `json:"username"`
	Team     []string  `json:"team"`
}

type DetailedAccount struct {
	Entity   `mapstructure:",squash"`
	Inbox    []Message    `json:"inbox"`
	Leagues  []LeagueInfo `json:"leagues"`
	Username string       `json:"username"`
	Team     []Player     `json:"team"`
}

// Firebase Auth database, following auth.UserToCreate struct
type AddUser struct {
	FullName string `binding:"required"`
	Username string `binding:"required"`
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

type LoginUser struct {
	IdToken string
}
