package controllers

import (
	"fantasy/database/entities"
	"fantasy/database/lib"
	"fantasy/database/utils"
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

	lib.InsertItem(ctx, "players", entities.GetPlayerEntity(
		input.Name, input.Role, input.Team,
	))

	ctx.JSON(http.StatusOK, gin.H{"addedPlayer": true})
}

func QueryPlayersByName(ctx *gin.Context) {
	term := ctx.Query("term")
	if term != "" {
		term = utils.Capitalize(term)
	}
	queryLimit := 10

	players := lib.QueryTermInField[entities.Player](ctx, "players", lib.Query{
		Field: "FirstName",
		Limit: queryLimit,
		Term:  term,
	})

	if len(players) < queryLimit {
		byLastName := lib.QueryTermInField[entities.Player](ctx, "players", lib.Query{
			Field: "LastName",
			Limit: queryLimit,
			Term:  term,
		})
		players = utils.ConcatDeduplicate(players, byLastName)
	}

	if len(players) > queryLimit {
		players = players[:queryLimit]
	}

	ctx.JSON(http.StatusOK, gin.H{"players": players})
}
