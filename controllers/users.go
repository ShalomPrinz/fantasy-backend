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

func GetUserInfo(ctx *gin.Context) {
	UID := ctx.MustGet("UID").(string)

	user := lib.GetSingle[entities.Account](ctx, "accounts", UID)
	team := lib.GetByIds[entities.Player](ctx, "players", user.Team)
	detailed := entities.DetailedAccount{
		Entity:   user.Entity,
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

	lib.InsertItemCustomID(ctx, "accounts", UID, entities.AddAccount{
		Nickname: input.Nickname,
		Team:     []string{},
	})

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
