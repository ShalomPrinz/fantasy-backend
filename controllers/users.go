package controllers

import (
	"fantasy/database/entities"
	"fantasy/database/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddTeamPlayer(ctx *gin.Context) {
	UID := ctx.MustGet("UID").(string)

	var input entities.Entity
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isPlayerExists(ctx, input.ID) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no-such-player"})
		return
	}

	path := "players/" + input.ID
	if appError := lib.InsertItemIntoArray(ctx, "accounts", UID, "Team", path); appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func RemoveTeamPlayer(ctx *gin.Context) {
	UID := ctx.MustGet("UID").(string)

	var input entities.Entity
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isPlayerExists(ctx, input.ID) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no-such-player"})
		return
	}

	path := "players/" + input.ID
	if appError := lib.RemoveItemFromArray(ctx, "accounts", UID, "Team", path); appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func GetUserInfo(ctx *gin.Context) {
	UID := ctx.MustGet("UID").(string)

	user, appError := lib.GetSingle[entities.Account](ctx, "accounts", UID)
	if appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	team, appError := lib.GetByIds[entities.Player](ctx, "players", user.Team)
	if appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	leagues, appError := lib.GetByIds[entities.League](ctx, "leagues", user.Leagues)
	if appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	detailed := entities.DetailedAccount{
		Entity:   user.Entity,
		Leagues:  entities.LeaguesToLeaguesInfo(leagues),
		Nickname: user.Nickname,
		Team:     team,
	}
	ctx.JSON(http.StatusOK, gin.H{"user": detailed})
}

func NewUser(ctx *gin.Context) {
	var input entities.AddUser
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	UID, appError := lib.CreateUser(ctx, input)
	if appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	account := entities.InsertAccount{
		Leagues:  []string{},
		Nickname: input.Nickname,
		Team:     []string{},
	}
	if appError := lib.InsertItemCustomID(ctx, "accounts", UID, account); appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

// For internal use only
func signUserToLeague(ctx *gin.Context, UID string, leagueId string) lib.AppError {
	path := "leagues/" + leagueId
	return lib.InsertItemIntoArray(ctx, "accounts", UID, "Leagues", path)
}
