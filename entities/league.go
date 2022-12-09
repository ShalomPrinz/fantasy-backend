package entities

import (
	"fantasy/database/utils"
	"log"
	"strings"
)

type League struct {
	Entity  `mapstructure:",squash"`
	Members []string `json:"members"`
	Name    string   `json:"name"`
}

type AddLeague struct {
	Members []string `json:"members"`
	Name    string   `json:"name"`
}

type DetailedLeague struct {
	Entity  `mapstructure:",squash"`
	Members []Account `json:"members"`
	Name    string    `json:"name"`
}

func getLeagueMemberId(member any) string {
	id, properCast := member.(string)
	if !properCast {
		log.Fatal("Error in member cast to account")
	}
	return strings.Replace(id, "accounts/", "", 1)
}

func LeagueContainsMember(league League, memberId string) bool {
	return utils.ArrayContainsString(league.Members, memberId, getLeagueMemberId)
}
