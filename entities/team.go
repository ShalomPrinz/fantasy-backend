package entities

type Team struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AddTeam struct {
	Name string `json:"name"`
}

var TeamAttributes = [...]string{"Name"}
