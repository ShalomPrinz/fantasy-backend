package controllers

import (
	"fantasy/database/entities"
	"fantasy/database/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewLeague(ctx *gin.Context) {
	var input entities.AddLeague
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var membersRefs []string
	for _, memberId := range input.Members {
		membersRefs = append(membersRefs, "accounts/"+memberId)
	}

	leagueId := lib.InsertItem(ctx, "leagues", entities.AddLeague{
		Members: membersRefs, Name: input.Name,
	})

	for _, memberId := range input.Members {
		signUserToLeague(ctx, memberId, leagueId)
	}

	ctx.JSON(http.StatusOK, gin.H{"addedLeague": true})
}
