package entities

type Player struct {
	Entity `mapstructure:",squash"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Team   string `json:"team"`
}

type AddPlayer struct {
	Name string `json:"name"`
	Role string `json:"role"`
	Team string `json:"team"`
}
