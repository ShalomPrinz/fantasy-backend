package lib

import (
	"context"
	"fantasy/database/entities"
	"log"

	"firebase.google.com/go/auth"
)

func CreateUser(props entities.AddUser) (string, error) {
	ctx := context.Background()
	params := (&auth.UserToCreate{}).
		DisplayName(props.DisplayName).
		Email(props.Email).
		Password(props.Password)

	user, err := Auth.CreateUser(ctx, params)
	if err != nil {
		log.Printf("error creating user: %v\n", err)
		return "", err
	}
	return user.UID, nil
}
