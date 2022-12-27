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
	if appError := lib.BindRequestJSON(ctx, &input); appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
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

	ctx.JSON(http.StatusOK, gin.H{"leagueId": leagueId})
}

func GetLeagueInfo(ctx *gin.Context) {
	UID := ctx.MustGet("UID").(string)

	leagueId := ctx.Query("id")
	if leagueId == "" {
		appError := lib.Error(http.StatusBadRequest, "missing-request-data")
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	league, appError := lib.GetSingle[entities.League](ctx, "leagues", leagueId)
	if appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}
	if !entities.LeagueContainsMember(league, UID) {
		appError := lib.Error(http.StatusUnauthorized, "not-league-member")
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
