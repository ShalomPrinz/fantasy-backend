package entities

type Team struct {
	ID string `json:"id"`
}

type AddTeam struct {
	ID string `json:"id"`
}

var TeamAttributes = [...]string{}
