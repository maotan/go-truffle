package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maotan/go-truffle/logger"
	"github.com/maotan/go-truffle/truffle"
	log "github.com/sirupsen/logrus"
)

var (
	router = gin.Default()
)

// Run will start the server
func Run() error{
	router.Use(logger.LogerMiddleware())
	router.Use(truffle.Recover)
	getRoutes()
	log.Info("app run...")
	return router.Run(":5000")
}

// getRoutes will create our routes of our entire application
// this way every group of routes can be defined in their own file
// so this one won't be so messy
func getRoutes() {
	router.GET("/actuator/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1 := router.Group("/v1")
	addUserRoutes(v1)
	addPingRoutes(v1)

	v2 := router.Group("/v2")
	addPingRoutes(v2)
}