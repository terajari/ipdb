package main

import (
	"errors"
	"net/http"

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

	err = app.Models.User.Insert(user)

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

	ctx.JSON(http.StatusCreated, gin.H{
		"user": user,
	})
}
