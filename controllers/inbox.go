package controllers

import (
	"fantasy/database/entities"
	"fantasy/database/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewLeagueInvitation(ctx *gin.Context) {
	UID := ctx.MustGet("UID").(string)

	var input entities.AddLeagueInvitation
	if appError := lib.BindRequestJSON(ctx, &input); appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	message := &entities.Message{
		From:     UID,
		LeagueId: input.LeagueId,
	}

	if appError := lib.InsertItemIntoArray(ctx, "accounts", input.To, "Inbox", message); appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
