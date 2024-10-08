package controllers

import (
	"fantasy/database/entities"
	"fantasy/database/lib"
	"fantasy/database/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPlayer(ctx *gin.Context) {
	player, appError := lib.GetSingle[entities.Player](ctx, "players", ctx.Param("id"))
	if appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"player": player})
}

func NewPlayer(ctx *gin.Context) {
	var input entities.AddPlayer
	if appError := lib.BindRequestJSON(ctx, &input); appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	playerId, appError := lib.InsertItem(ctx, "players", entities.GetInsertPlayer(
		input.Name, input.Role, input.Team,
	))
	if appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"playerId": playerId})
}

func isPlayerExists(ctx *gin.Context, id string) bool {
	return lib.IsExists(ctx, "players/"+id)
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
