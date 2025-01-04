package container

import (
	"github.com/gin-gonic/gin"
	"github.com/saadrupai/go-message-broker/app/broker"
	"github.com/saadrupai/go-message-broker/app/config"
)

func Serve(router *gin.Engine) {

	apiVersion := router.Group("/api/v1")

	broker := broker.NewBroker()

	handler := handler.NewHandler(broker)

	apiVersion.POST("/create-queue", handler.QueueHandler)
	apiVersion.POST("/publish", handler.PublishHandler)
	apiVersion.POST("/subscribe", handler.SubscribeHandler)

	router.Run(":" + config.LocalConfig.Port)
}
