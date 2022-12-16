package lib

import (
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppError struct {
	Code int
	Json gin.H
}

var EmptyError AppError

func (err AppError) HasError() bool {
	return err.Code != 0
}

func Error(code int, message string) AppError {
	return AppError{
		Code: code,
		Json: gin.H{"error": message},
	}
}

const (
	serverErrorCode    = http.StatusInternalServerError
	serverErrorMessage = "internal-server-error"
)

var (
	genericServerError = AppError{
		Code: serverErrorCode,
		Json: gin.H{"error": serverErrorMessage},
	}
)

func CreateUserError(err error) AppError {
	code, message := serverErrorCode, serverErrorMessage

	if auth.IsEmailAlreadyExists(err) {
		code, message = http.StatusBadRequest, "email-already-exists"
	} else if auth.IsInvalidEmail(err) {
		code, message = http.StatusUnprocessableEntity, "invalid-email"
	}

	return Error(code, message)
}

func GetDocumentError(err error) AppError {
	code, message := serverErrorCode, serverErrorMessage

	if status.Code(err) == codes.NotFound {
		code, message = http.StatusNotFound, "not-found"
	}

	return Error(code, message)
}

func InsertItemError(err error) AppError {
	return genericServerError
}

func RemoveItemError(err error) AppError {
	return genericServerError
}

func VerifyTokenError(err error) AppError {
	code, message := serverErrorCode, serverErrorMessage

	if strings.Contains(err.Error(), "ID token has expired at:") {
		code, message = http.StatusUnauthorized, "id-token-expired"
	} else if strings.Contains(err.Error(), "ID token issued at future timestamp:") {
		code, message = http.StatusUnauthorized, "id-token-expired"
	} else if auth.IsIDTokenRevoked(err) {
		code, message = http.StatusUnauthorized, "id-token-revoked"
	} else if auth.IsSessionCookieRevoked(err) {
		code, message = http.StatusUnauthorized, "session-cookie-revoked"
	} else if auth.IsUserNotFound(err) {
		code, message = http.StatusNotFound, "user-not-found"
	}

	return Error(code, message)
}
