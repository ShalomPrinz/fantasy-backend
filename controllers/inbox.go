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

	message := &entities.InsertLeagueInvitation{
		From:     UID,
		LeagueId: input.LeagueId,
	}

	inboxRef := lib.SubCollectionRef("accounts", input.To, "inbox")
	messageId, appError := lib.InsertItemToCollection(ctx, inboxRef, message)
	if appError.HasError() {
		ctx.JSON(appError.Code, appError.Json)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"messageId": messageId})
}
