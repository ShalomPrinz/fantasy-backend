package lib

import "github.com/gin-gonic/gin"

func BindRequestJSON(ctx *gin.Context, json any) AppError {
	if err := ctx.ShouldBindJSON(&json); err != nil {
		return JsonBindingError(err)
	}
	return EmptyError
}
