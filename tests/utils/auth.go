package utils

import (
	"fantasy/database/utils"
	"fmt"
)

// Firebase Auth Sign In Response: https://firebase.google.com/docs/reference/rest/auth#section-sign-in-email-password
type IdTokenResponse struct {
	Email        string `json:"email"`
	IdToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalId      string `json:"localId"`
	Registered   bool   `json:"registered"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GenerateIdToken(loginDetails LoginUser) (string, error) {
	path := "identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=testkey"

	var response map[string]any
	if err := PostAuth(path, loginDetails, &response); err != nil {
		fmt.Println("Failed generating id token", err)
		return "", err
	}

	token := utils.MapToStruct[IdTokenResponse](response)
	return token.IdToken, nil
}
