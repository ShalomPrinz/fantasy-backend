package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HasIdToken(ctx *gin.Context) {
	idToken := ctx.GetHeader(os.Getenv("AUTHHEADER"))
	if idToken == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No logged in user"})
	}
	ctx.Next()
}
