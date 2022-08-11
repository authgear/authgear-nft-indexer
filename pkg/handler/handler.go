package handler

import "github.com/gin-gonic/gin"

func HandleError(ctx *gin.Context, code int, err error) {
	ctx.AbortWithStatusJSON(code, gin.H{ // nolint: errcheck
		"error": err.Error(),
	})
}
