package lib

import (
	"fantasy/database/entities"
	"log"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

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
