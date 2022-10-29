package entities

type Player struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type AddPlayer struct {
	Name string `json:"name"`
	Role string `json:"role"`
}
