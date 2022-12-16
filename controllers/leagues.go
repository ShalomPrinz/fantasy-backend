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
	leagueId, appError := lib.InsertItem(ctx, "leagues", entities.InsertLeague{
		Members: memberRef, Name: input.Name,
	})
	if appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	if appError = signUserToLeague(ctx, UID, leagueId); appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"addedLeague": true})
}

func GetLeagueInfo(ctx *gin.Context) {
	UID := ctx.MustGet("UID").(string)

	leagueId := ctx.Query("id")
	if leagueId == "" {
		appError := lib.Error(http.StatusBadRequest, "No League ID Supplied")
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	league, appError := lib.GetSingle[entities.League](ctx, "leagues", leagueId)
	if appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}
	if !entities.LeagueContainsMember(league, UID) {
		appError := lib.Error(http.StatusUnauthorized, "You are not a member of this league")
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	accounts, appError := lib.GetByIds[entities.Account](ctx, "accounts", league.Members)
	if appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}
	var mapError lib.AppError
	mapFunc := func(account entities.Account) entities.Member {
		if mapError.HasError() {
			return entities.Member{}
		}

		member, appError := mapAccountToMember(ctx, account)
		mapError = appError
		return member
	}
	members := utils.Map(accounts, mapFunc)
	if mapError.HasError() {
		return
	}

	detailedLeague := entities.DetailedLeague{
		Entity:  league.Entity,
		Members: members,
		Name:    league.Name,
	}
	ctx.JSON(http.StatusOK, gin.H{"league": detailedLeague})
}

func mapAccountToMember(ctx *gin.Context, account entities.Account) (entities.Member, lib.AppError) {
	team, appError := lib.GetByIds[entities.Player](ctx, "players", account.Team)
	if appError.HasError() {
		return entities.Member{}, appError
	}

	return entities.Member{
		Entity:   account.Entity,
		Nickname: account.Nickname,
		Team:     team,
	}, lib.EmptyError
}
