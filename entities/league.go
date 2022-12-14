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

type LeagueInfo struct {
	Entity       `mapstructure:",squash"`
	MembersCount int    `json:"membersCount"`
	Name         string `json:"name"`
}

type AddLeague struct {
	Name string `json:"name"`
}

type InsertLeague struct {
	Members []string `json:"members"`
	Name    string   `json:"name"`
}

type DetailedLeague struct {
	Entity  `mapstructure:",squash"`
	Members []Member `json:"members"`
	Name    string   `json:"name"`
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

func LeaguesToLeaguesInfo(leagues []League) []LeagueInfo {
	return utils.Map(leagues, func(l League) LeagueInfo {
		return LeagueInfo{
			Entity:       l.Entity,
			MembersCount: len(l.Members),
			Name:         l.Name,
		}
	})
}
