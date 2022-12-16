package controllers

import (
	"fantasy/database/entities"
	"fantasy/database/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTeams(ctx *gin.Context) {
	teams, appError := lib.GetAll[entities.Team](ctx, "teams")
	if appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"teams": teams})
}

func GetTeam(ctx *gin.Context) {
	team, appError := lib.GetSingle[entities.Team](ctx, "teams", ctx.Param("id"))
	if appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"team": team})
}

func NewTeam(ctx *gin.Context) {
	var input entities.AddTeam
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if appError := lib.InsertItemCustomID(ctx, "teams", input.ID, map[string]any{}); appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"addedTeam": true})
}
