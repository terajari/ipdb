package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrBadRequest        = "IPDB-001 - Use correct request format"
	ErrValidation        = "IPDB-002 - Validation error"
	ErrNotFound          = "IPDB-003 - Resource not found"
	ErrServer            = "IPDB-004 - Server error"
	ErrRateLimitExceeded = "IPDB-005 - Rate limit exceeded"
)

func (app *application) badRequestResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"status":  http.StatusBadRequest,
		"message": ErrBadRequest,
		"error":   err.Error(),
	})
}

func (app *application) failedValidationResponse(ctx *gin.Context, errs map[string]string) {
	ctx.JSON(http.StatusUnprocessableEntity, gin.H{
		"status":  http.StatusUnprocessableEntity,
		"message": ErrValidation,
		"errors":  errs,
	})
}

func (app *application) serverErrorResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"status":  http.StatusInternalServerError,
		"message": ErrServer,
		"error":   err.Error(),
	})
}

func (app *application) notFoundResponse(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"status":  http.StatusNotFound,
		"message": ErrNotFound,
	})
}

func (app *application) rateLimitExceededResponse(ctx *gin.Context) {
	ctx.JSON(http.StatusTooManyRequests, gin.H{
		"status":  http.StatusTooManyRequests,
		"message": ErrRateLimitExceeded,
	})
}
