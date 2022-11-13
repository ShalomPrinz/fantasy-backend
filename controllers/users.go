package controllers

import (
	"fantasy/database/entities"
	"fantasy/database/lib"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func GetUserInfo(ctx *gin.Context) {
	idToken, err := ctx.Cookie(os.Getenv("AUTHCOOKIE"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	UID, err := lib.GetUidByToken(ctx, idToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user := lib.GetSingle(ctx, "accounts", UID, entities.AccountAttributes[:])
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func LoginUser(ctx *gin.Context) {
	var input entities.LoginUser
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expiresIn := time.Minute * 20
	cookie, err := lib.CreateCookieSession(ctx, input.IdToken, expiresIn)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ctx.SetCookie(os.Getenv("AUTHCOOKIE"), cookie, int(expiresIn.Seconds()), "/", "localhost", true, true)
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     os.Getenv("AUTHCOOKIE"),
		Value:    url.QueryEscape(cookie),
		MaxAge:   int(expiresIn.Seconds()),
		Path:     "/",
		Domain:   "localhost",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		HttpOnly: true,
	})

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
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
