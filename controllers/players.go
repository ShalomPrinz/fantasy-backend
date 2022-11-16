package controllers

import (
	"fantasy/database/entities"
	"fantasy/database/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPlayers(ctx *gin.Context) {
	players := lib.GetAll[entities.Player](ctx, "players")
	ctx.JSON(http.StatusOK, gin.H{"players": players})
}

func GetPlayer(ctx *gin.Context) {
	player := lib.GetSingle[entities.Player](ctx, "players", ctx.Param("id"))
	ctx.JSON(http.StatusOK, gin.H{"player": player})
}

func NewPlayer(ctx *gin.Context) {
	var input entities.AddPlayer
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lib.InsertItem(ctx, "players", entities.AddPlayer{
		Name: input.Name,
		Role: input.Role,
	})

	ctx.JSON(http.StatusOK, gin.H{"addedPlayer": true})
}
