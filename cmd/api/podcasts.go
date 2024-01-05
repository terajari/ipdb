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

	err := app.models.Podcast.Insert(&podcast)
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
		app.logger.Info(err.Error())
		app.badRequestResponse(ctx, err)
		return
	}

	podcast, err := app.models.Podcast.FindById(path.Id)
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

	podcast, err := app.models.Podcast.FindById(path.Id)

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

	if input.Title != "" {
		podcast.Title = input.Title
	}
	if input.Platform != "" {
		podcast.Platform = input.Platform
	}
	if input.Url != "" {
		podcast.Url = input.Url
	}
	if input.Host != "" {
		podcast.Host = input.Host
	}
	if input.Program != "" {
		podcast.Program = input.Program
	}
	if input.GuestSpeakers != nil {
		podcast.GuestSpeakers = input.GuestSpeakers
	}
	if input.Year != 0 {
		podcast.Year = input.Year
	}
	if input.Language != "" {
		podcast.Language = input.Language
	}
	if input.Tags != nil {
		podcast.Tags = input.Tags
	}

	v := validator.New()

	data.ValidatePodcast(v, &data.Podcast{
		Title:         podcast.Title,
		Platform:      podcast.Platform,
		Url:           podcast.Url,
		Host:          podcast.Host,
		Program:       podcast.Program,
		GuestSpeakers: podcast.GuestSpeakers,
		Year:          podcast.Year,
		Language:      podcast.Language,
		Tags:          podcast.Tags,
	})

	if !v.Valid() {
		app.failedValidationResponse(ctx, v.Errors)
		return
	}

	err = app.models.Podcast.UpdatePodcast(podcast)
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

	err := app.models.Podcast.DeleteById(path.Id)
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

func (app *application) listPodcastHandler(ctx *gin.Context) {
	var input struct {
		Platform string   `form:"platform"`
		Tags     []string `form:"tags"`
		data.Filters
	}

	input.Filters.SortSafelist = []string{
		"title",
		"platform",
		"host",
		"program",
		"guest_speakers",
		"year",
		"language",
		"tags",
	}

	input.Filters = *data.DefaultsFilters(input.Filters)

	if err := ctx.ShouldBindQuery(&input); err != nil {
		app.badRequestResponse(ctx, err)
		return
	}

	v := validator.New()

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(ctx, v.Errors)
		return
	}

	podcasts, metadata, err := app.models.Podcast.GetAll(input.Platform, input.Tags, input.Filters)
	if err != nil {
		app.serverErrorResponse(ctx, err)
		return
	}

	app.logger.Info(input.Platform)

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": podcasts, "metadata": metadata})
}
