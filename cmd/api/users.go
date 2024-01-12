package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/terajari/ipdb/internal/data"
	"github.com/terajari/ipdb/internal/validator"
)

func (app *application) createUserHandler(ctx *gin.Context) {

	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		app.badRequestResponse(ctx, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err := user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(ctx, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(ctx, v.Errors)
		return
	}

	err = app.models.User.Insert(user)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(ctx, v.Errors)
		default:
			app.serverErrorResponse(ctx, err)
		}
		return
	}

	token, err := app.models.Token.New(user.Id, 3*(24*time.Hour), data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(ctx, err)
		return
	}

	app.background(func() {

		data := map[string]any{
			"Name":           user.Name,
			"Id":             user.Id,
			"TokenPlainText": token.Plaintext,
		}

		err = app.mailler.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.serverErrorResponse(ctx, err)
			return
		}
	})

	ctx.JSON(http.StatusCreated, gin.H{
		"user": user,
	})
}

func (app *application) activateUserHandler(ctx *gin.Context) {

	var input struct {
		Token string `json:"token"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		app.badRequestResponse(ctx, err)
	}

	v := validator.New()

	if data.ValidatePlaintext(v, input.Token); !v.Valid() {
		app.failedValidationResponse(ctx, v.Errors)
	}

	user, err := app.models.User.GetForToken(data.ScopeActivation, input.Token)
	if err != nil {
		app.serverErrorResponse(ctx, err)
		return
	}

	user.Activated = true

	err = app.models.User.ActivateUser(user.Id)
	if err != nil {
		app.serverErrorResponse(ctx, err)
		return
	}

	err = app.models.Token.DeleteAll(data.ScopeActivation, user.Id)
	if err != nil {
		app.serverErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
