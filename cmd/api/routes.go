package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) routes() *gin.Engine {
	router := gin.Default()

	rg := router.Group("/v1")
	rg.GET("/healthcheck", app.healthcheckHandler)
	rg.POST("/podcasts", app.createPodcastHandler)
	rg.GET("/podcasts/:id", app.getPodcastsHandler)
	rg.PUT("/podcasts/:id", app.updatePodcastHandler)
	rg.DELETE("/podcasts/:id", app.deletePodcastHandler)

	return router
}
