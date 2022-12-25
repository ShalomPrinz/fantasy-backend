package controllers

import (
	"fantasy/database/lib"
	"os"

	"github.com/gin-gonic/gin"
)

func VerifyIdToken(ctx *gin.Context) {
	idToken := ctx.GetHeader(os.Getenv("AUTHHEADER"))

	UID, appError := lib.GetUidByToken(ctx, idToken)
	if appError.HasError() {
		ctx.AbortWithStatusJSON(appError.Code, appError.Json)
		return
	}

	ctx.Set("UID", UID)
	ctx.Next()
}
