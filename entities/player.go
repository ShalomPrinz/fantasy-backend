package entities

import "strings"

type Player struct {
	Entity    `mapstructure:",squash"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      string `json:"role"`
	Team      string `json:"team"`
}

type AddPlayer struct {
	Name string `json:"name" binding:"required"`
	Role string `json:"role" binding:"required"`
	Team string `json:"team" binding:"required"`
}

func GetPlayerEntity(Name string, Role string, Team string) Player {
	player := Player{
		Role: Role,
		Team: Team,
	}

	firstName := Name
	lastName := ""

	if strings.Contains(firstName, " ") {
		split := strings.SplitN(firstName, " ", 2)
		firstName = split[0]
		lastName = split[1]
		player.LastName = lastName
	}

	player.FirstName = firstName
	return player
}
