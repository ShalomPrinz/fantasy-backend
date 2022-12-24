package lib

import (
	"fantasy/database/entities"
	"log"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

func CreateUser(ctx *gin.Context, props entities.AddUser) (string, AppError) {
	params := (&auth.UserToCreate{}).
		DisplayName(props.FullName).
		Email(props.Email).
		Password(props.Password)

	user, err := Auth.CreateUser(ctx, params)
	if err != nil {
		log.Printf("error creating user: %v\n", err)
		return "", CreateUserError(err)
	}
	return user.UID, EmptyError
}

func GetUidByToken(ctx *gin.Context, idToken string) (string, AppError) {
	decoded, err := Auth.VerifyIDTokenAndCheckRevoked(ctx, idToken)
	if err != nil {
		log.Printf("ID Token %v is invalid. %v", idToken[:10], err)
		return "", VerifyTokenError(err)
	}
	return decoded.UID, EmptyError
}
