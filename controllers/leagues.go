package controllers

import (
	"fantasy/database/entities"
	"fantasy/database/lib"
	"fantasy/database/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewLeague(ctx *gin.Context) {
	UID := ctx.MustGet("UID").(string)

	var input entities.AddLeague
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	memberRef := []string{"accounts/" + UID}
	leagueId := lib.InsertItem(ctx, "leagues", entities.InsertLeague{
		Members: memberRef, Name: input.Name,
	})
	signUserToLeague(ctx, UID, leagueId)

	ctx.JSON(http.StatusOK, gin.H{"addedLeague": true})
}

func GetLeagueInfo(ctx *gin.Context) {
	UID := ctx.MustGet("UID").(string)

	leagueId := ctx.Query("id")
	if leagueId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No League ID Supplied"})
		return
	}

	league := lib.GetSingle[entities.League](ctx, "leagues", leagueId)
	if !entities.LeagueContainsMember(league, UID) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "You are not a member of this league"})
		return
	}

	accounts := lib.GetByIds[entities.Account](ctx, "accounts", league.Members)

	mapFunc := func(account entities.Account) entities.Member {
		return mapAccountToMember(ctx, account)
	}
	members := utils.Map(accounts, mapFunc)

	detailedLeague := entities.DetailedLeague{
		Entity:  league.Entity,
		Members: members,
		Name:    league.Name,
	}
	ctx.JSON(http.StatusOK, gin.H{"league": detailedLeague})
}

func mapAccountToMember(ctx *gin.Context, account entities.Account) entities.Member {
	team := lib.GetByIds[entities.Player](ctx, "players", account.Team)
	return entities.Member{
		Entity:   account.Entity,
		Nickname: account.Nickname,
		Team:     team,
	}
}
