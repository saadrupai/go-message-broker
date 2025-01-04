package container

import (
	"github.com/gin-gonic/gin"
	"github.com/saadrupai/go-message-broker/app/config"
)

func Serve(router *gin.Engine) {

	apiVersion := router.Group("/api/v1")

	apiVersion.POST("/create-queue")
	apiVersion.POST("/publish")
	apiVersion.POST("/subscribe")

	router.Run(":" + config.LocalConfig.Port)
}
