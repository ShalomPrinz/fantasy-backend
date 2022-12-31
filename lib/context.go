package lib

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BindRequestJSON(ctx *gin.Context, json any) AppError {
	if err := ctx.ShouldBindJSON(&json); err != nil {
		return JsonBindingError(err)
	}
	return EmptyError
}

func GetRequiredParam(ctx *gin.Context, paramKey string) (string, AppError) {
	value := ctx.Query(paramKey)
	if value == "" {
		return "", Error(http.StatusBadRequest, "missing-request-data")
	}
	return value, EmptyError
}
