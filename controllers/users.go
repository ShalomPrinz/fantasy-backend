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
	if appError := lib.BindRequestJSON(ctx, &input); appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
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
	if appError := lib.BindRequestJSON(ctx, &input); appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
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

// For internal use only
func getUserTeam(ctx *gin.Context, uid string) (entities.Account, []entities.Player, lib.AppError) {
	user, appError := lib.GetSingle[entities.Account](ctx, "accounts", uid)
	if appError.HasError() {
		return entities.Account{}, nil, appError
	}

	team, appError := lib.GetByIds[entities.Player](ctx, user.Team)
	if appError.HasError() {
		return entities.Account{}, nil, appError
	}

	return user, team, lib.EmptyError
}

func GetUserInfo(ctx *gin.Context) {
	UID := ctx.MustGet("UID").(string)

	user, team, appError := getUserTeam(ctx, UID)
	if appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	inboxRef := lib.SubCollectionRef("accounts", UID, "inbox")
	inbox, appError := lib.GetAll[entities.Message](ctx, inboxRef)
	if appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	leagues, appError := lib.GetByIds[entities.League](ctx, user.Leagues)
	if appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	detailed := entities.DetailedAccount{
		Entity:   user.Entity,
		Inbox:    inbox,
		Leagues:  entities.LeaguesToLeaguesInfo(leagues),
		Username: user.Username,
		Team:     team,
	}
	ctx.JSON(http.StatusOK, gin.H{"user": detailed})
}

func NewUser(ctx *gin.Context) {
	var input entities.AddUser
	if appError := lib.BindRequestJSON(ctx, &input); appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	UID, appError := lib.CreateUser(ctx, input)
	if appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	account := entities.InsertAccount{
		Leagues:  []string{},
		Username: input.Username,
		Team:     []string{},
	}
	if appError := lib.InsertItemCustomID(ctx, "accounts", UID, account); appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"userId": UID})
}

// For internal use only
func signUserToLeague(ctx *gin.Context, UID string, leagueId string) lib.AppError {
	path := "leagues/" + leagueId
	return lib.InsertItemIntoArray(ctx, "accounts", UID, "Leagues", path)
}
