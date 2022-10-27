package entities

type Player struct {
	ID   uint   `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
	Role string `json:"role"`
	Team string `json:"team"`
}

type AddPlayer struct {
	Name string `json:"name"`
	Role string `json:"role"`
	Team string `json:"team"`
}
