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

	path := "players/" + input.ID
	lib.InsertItemIntoArray(ctx, "accounts", UID, "Team", path)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func RemoveTeamPlayer(ctx *gin.Context) {
	UID := ctx.MustGet("UID").(string)

	var input entities.Entity
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	path := "players/" + input.ID
	lib.RemoveItemFromArray(ctx, "accounts", UID, "Team", path)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func GetUserInfo(ctx *gin.Context) {
	UID := ctx.MustGet("UID").(string)

	user := lib.GetSingle[entities.Account](ctx, "accounts", UID)
	team := lib.GetByIds[entities.Player](ctx, "players", user.Team)
	leagues := lib.GetByIds[entities.League](ctx, "leagues", user.Leagues)

	detailed := entities.DetailedAccount{
		Entity:   user.Entity,
		Leagues:  leagues,
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

	UID, err := lib.CreateUser(ctx, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	lib.InsertItemCustomID(ctx, "accounts", UID, entities.InsertAccount{
		Leagues:  []string{},
		Nickname: input.Nickname,
		Team:     []string{},
	})

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

// For internal use only
func signUserToLeague(ctx *gin.Context, UID string, leagueId string) {
	path := "leagues/" + leagueId
	lib.InsertItemIntoArray(ctx, "accounts", UID, "Leagues", path)
}
