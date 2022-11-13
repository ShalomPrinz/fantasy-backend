package controllers

import (
	"fantasy/database/entities"
	"fantasy/database/lib"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func GetUserInfo(ctx *gin.Context) {
	idToken := ctx.GetHeader(os.Getenv("AUTHHEADER"))
	UID, err := lib.GetUidByToken(ctx, idToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user := lib.GetSingle(ctx, "accounts", UID, entities.AccountAttributes[:])
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func NewUser(ctx *gin.Context) {
	var input entities.AddUser
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	UID, err := lib.CreateUser(ctx, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	lib.InsertItemCustomID(ctx, "accounts", UID, entities.AddAccount{
		Nickname: input.Nickname,
		Team:     []string{},
	})

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
