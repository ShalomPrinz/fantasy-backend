package controllers

import (
	"fantasy/database/entities"
	"fantasy/database/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTeams(ctx *gin.Context) {
	teams := lib.GetAll[entities.Team](ctx, "teams")
	ctx.JSON(http.StatusOK, gin.H{"teams": teams})
}

func GetTeam(ctx *gin.Context) {
	team := lib.GetSingle[entities.Team](ctx, "teams", ctx.Param("id"))
	ctx.JSON(http.StatusOK, gin.H{"team": team})
}

func NewTeam(ctx *gin.Context) {
	var input entities.AddTeam
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For now Team only has ID. Later I will replace this function call
	lib.InsertItemCustomID(ctx, "teams", input.ID, map[string]interface{}{})

	ctx.JSON(http.StatusOK, gin.H{"addedTeam": true})
}
