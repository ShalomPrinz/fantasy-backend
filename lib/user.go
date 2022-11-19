package lib

import (
	"fantasy/database/entities"
	"log"
	"time"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

func CreateCookieSession(ctx *gin.Context, idToken string, expiresIn time.Duration) (string, error) {
	cookie, err := Auth.SessionCookie(ctx, idToken, expiresIn)
	if err != nil {
		log.Printf("Failed to create a session cookie: %v", err)
		return "", err
	}
	return cookie, nil
}

func CreateUser(ctx *gin.Context, props entities.AddUser) (string, error) {
	params := (&auth.UserToCreate{}).
		DisplayName(props.FullName).
		Email(props.Email).
		Password(props.Password)

	user, err := Auth.CreateUser(ctx, params)
	if err != nil {
		log.Printf("error creating user: %v\n", err)
		return "", err
	}
	return user.UID, nil
}

func GetUidByToken(ctx *gin.Context, idToken string) (string, error) {
	decoded, err := Auth.VerifyIDTokenAndCheckRevoked(ctx, idToken)
	if err != nil {
		log.Printf("ID Token %v is invalid. %v", idToken[:10], err)
		return "", err
	}

	return decoded.UID, nil
}
