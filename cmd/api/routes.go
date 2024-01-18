package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) routes() *gin.Engine {
	router := gin.Default()

	rg := router.Group("/v1")
	rg.Use(app.recoverPanic(), app.rateLimit())

	rg.GET("/podcasts", app.listPodcastHandler)
	rg.GET("/healthcheck", app.healthcheckHandler)
	rg.POST("/podcasts", app.createPodcastHandler)
	rg.GET("/podcasts/:id", app.getPodcastsHandler)
	rg.PUT("/podcasts/:id", app.updatePodcastHandler)
	rg.DELETE("/podcasts/:id", app.deletePodcastHandler)

	rg.POST("/users", app.createUserHandler)
	rg.PUT("/users/activated", app.activateUserHandler)

	rg.POST("/tokens/authentication", app.authTokenHandler)

	return router
}
