package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/terajari/ipdb/internal/data"
)

func (app *application) authTokenHandler(ctx *gin.Context) {

	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		app.badRequestResponse(ctx, err)
		return
	}

	user, err := app.models.User.GetByEmail(input.Email)
	if err != nil {
		app.serverErrorResponse(ctx, err)
		return
	}

	match, err := app.models.User.Matches(user, input.Password)
	if err != nil {
		app.serverErrorResponse(ctx, err)
	}

	if !match {
		app.invalidCredentialResponse(ctx)
		return
	}

	token, err := app.models.Token.New(user.Id, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"authorization_token": token,
	})

}
