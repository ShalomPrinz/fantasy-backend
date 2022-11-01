package entities

type Player struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	Team string `json:"team"`
}

type AddPlayer struct {
	Name string `json:"name"`
	Role string `json:"role"`
	Team string `json:"team"`
}

var PlayerAttributes = [...]string{"Name", "Role", "Team"}
