package controllers

import (
	"fantasy/database/lib"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func VerifyIdToken(ctx *gin.Context) {
	idToken := ctx.GetHeader(os.Getenv("AUTHHEADER"))
	if idToken == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No logged in user"})
		return
	}

	UID, err := lib.GetUidByToken(ctx, idToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.Set("UID", UID)
	ctx.Next()
}
