package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/terajari/ipdb/internal/data"
	"github.com/terajari/ipdb/internal/validator"
)

func (app *application) createPodcastHandler(ctx *gin.Context) {
	var input struct {
		Title         string   `json:"title"`
		Platform      string   `json:"platform"`
		Url           string   `json:"url"`
		Host          string   `json:"host"`
		Program       string   `json:"program"`
		GuestSpeakers []string `json:"guest_speakers"`
		Year          int64    `json:"year"`
		Language      string   `json:"language"`
		Tags          []string `json:"tags"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		app.badRequestResponse(ctx, err)
		return
	}

	v := validator.New()
	data.ValidatePodcast(v, &data.Podcast{
		Title:         input.Title,
		Platform:      input.Platform,
		Url:           input.Url,
		Host:          input.Host,
		Program:       input.Program,
		GuestSpeakers: input.GuestSpeakers,
		Year:          input.Year,
		Language:      input.Language,
		Tags:          input.Tags,
	})

	if !v.Valid() {
		app.failedValidationResponse(ctx, v.Errors)
		return
	}

	podcast := data.Podcast{
		Title:         input.Title,
		Platform:      input.Platform,
		Url:           input.Url,
		Host:          input.Host,
		Program:       input.Program,
		GuestSpeakers: input.GuestSpeakers,
		Year:          input.Year,
		Language:      input.Language,
		Tags:          input.Tags,
	}

	err := app.Models.Podcast.Insert(&podcast)
	if err != nil {
		app.serverErrorResponse(ctx, err)
		return
	}

	ctx.Header("Location", fmt.Sprintf("%s/%d", ctx.Request.URL.Path, podcast.Id))

	ctx.JSON(http.StatusCreated, gin.H{"status": http.StatusOK, "data": podcast})
}

func (app *application) getPodcastsHandler(ctx *gin.Context) {

	var path struct {
		Id int64 `uri:"id" binding:"required,gt=0"`
	}

	if err := ctx.ShouldBindUri(&path); err != nil {
		app.logger.Println(err)
		app.badRequestResponse(ctx, err)
		return
	}

	podcast, err := app.Models.Podcast.FindById(path.Id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(ctx)
			return
		default:
			app.serverErrorResponse(ctx, err)
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": podcast})
}

func (app *application) updatePodcastHandler(ctx *gin.Context) {
	var path struct {
		Id int64 `uri:"id" binding:"required,gt=0"`
	}

	if err := ctx.ShouldBindUri(&path); err != nil {
		app.badRequestResponse(ctx, err)
		return
	}

	var input struct {
		Title         string   `json:"title"`
		Platform      string   `json:"platform"`
		Url           string   `json:"url"`
		Host          string   `json:"host"`
		Program       string   `json:"program"`
		GuestSpeakers []string `json:"guest_speakers"`
		Year          int64    `json:"year"`
		Language      string   `json:"language"`
		Tags          []string `json:"tags"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		app.badRequestResponse(ctx, err)
		return
	}

	v := validator.New()

	data.ValidatePodcast(v, &data.Podcast{
		Title:         input.Title,
		Platform:      input.Platform,
		Url:           input.Url,
		Host:          input.Host,
		Program:       input.Program,
		GuestSpeakers: input.GuestSpeakers,
		Year:          input.Year,
		Language:      input.Language,
		Tags:          input.Tags,
	})

	if !v.Valid() {
		app.failedValidationResponse(ctx, v.Errors)
		return
	}

	podcast := data.Podcast{
		Title:         input.Title,
		Platform:      input.Platform,
		Url:           input.Url,
		Host:          input.Host,
		Program:       input.Program,
		GuestSpeakers: input.GuestSpeakers,
		Year:          input.Year,
		Language:      input.Language,
		Tags:          input.Tags,
	}

	err := app.Models.Podcast.UpdatePodcast(&podcast)
	if err != nil {
		app.serverErrorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": podcast})
}

func (app *application) deletePodcastHandler(ctx *gin.Context) {

	var path struct {
		Id int64 `uri:"id" binding:"required,gt=0"`
	}

	if err := ctx.ShouldBindUri(&path); err != nil {
		app.badRequestResponse(ctx, err)
		return
	}

	err := app.Models.Podcast.DeleteById(path.Id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(ctx)
			return
		default:
			app.serverErrorResponse(ctx, err)
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "podcast successfully deleted"})
}
