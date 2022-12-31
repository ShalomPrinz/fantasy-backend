package lib

import (
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
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

func isStatusNotFound(err error) bool {
	return status.Code(err) == codes.NotFound
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
	code, message := serverErrorCode, "create-user-error"

	if err.Error() == "password must be a string at least 6 characters long" {
		code, message = http.StatusBadRequest, "password-too-short"
	} else if strings.Contains(err.Error(), "malformed email string: ") {
		code, message = http.StatusBadRequest, "invalid-email"
	} else if auth.IsEmailAlreadyExists(err) {
		code, message = http.StatusBadRequest, "email-already-exists"
	}

	return Error(code, message)
}

func GetDocumentError(err error) AppError {
	code, message := serverErrorCode, "get-document-error"

	if isStatusNotFound(err) {
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
	code, message := serverErrorCode, "verify-token-error"

	if strings.Contains(err.Error(), "must be a non-empty string") {
		code, message = http.StatusUnauthorized, "id-token-missing"
	} else if strings.Contains(err.Error(), "incorrect number of segments") {
		code, message = http.StatusUnauthorized, "id-token-invalid"
	} else if strings.Contains(err.Error(), "ID token has expired at:") {
		code, message = http.StatusUnauthorized, "id-token-expired"
	} else if strings.Contains(err.Error(), "ID token issued at future timestamp:") {
		code, message = http.StatusUnauthorized, "id-token-expired"
	}

	return Error(code, message)
}

func JsonBindingError(err error) AppError {
	code, message := serverErrorCode, "data-binding-failure"

	if strings.Contains(err.Error(), "Error:Field validation for") &&
		strings.Contains(err.Error(), "failed on the 'required' tag") {
		code, message = http.StatusBadRequest, "missing-request-data"
	}

	return Error(code, message)
}
