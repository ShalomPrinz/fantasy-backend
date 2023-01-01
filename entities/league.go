package entities

import (
	"fantasy/database/utils"
	"log"
)

type role string

const (
	Admin   role = "admin"
	Regular role = "member"
)

type Member struct {
	Entity   `mapstructure:",squash"`
	Username string   `json:"username"`
	Team     []Player `json:"team"`
	Role     role     `json:"role"`
}

type MemberInfo struct {
	ID   string `json:"memberId"`
	Role role   `json:"role"`
}

type League struct {
	Entity  `mapstructure:",squash"`
	Members []MemberInfo `json:"members"`
	Name    string       `json:"name"`
}

type LeagueInfo struct {
	Entity       `mapstructure:",squash"`
	MembersCount int    `json:"membersCount"`
	Name         string `json:"name"`
}

type AddLeague struct {
	Name string `json:"name" binding:"required"`
}

type InsertLeague struct {
	Members []MemberInfo `json:"members"`
	Name    string       `json:"name"`
}

type DetailedLeague struct {
	Entity  `mapstructure:",squash"`
	Members []Member `json:"members"`
	Name    string   `json:"name"`
}

func getLeagueMemberId(member any) string {
	memberInfo, properCast := member.(MemberInfo)
	if !properCast {
		log.Fatal("Error in member cast to account")
	}
	return memberInfo.ID
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
