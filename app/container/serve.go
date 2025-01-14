package container

import (
	"github.com/gin-gonic/gin"
	"github.com/saadrupai/go-message-broker/app/broker"
	"github.com/saadrupai/go-message-broker/app/config"
	"github.com/saadrupai/go-message-broker/app/handler"
)

func Serve(router *gin.Engine) {

	apiVersion := router.Group("/api/v1")

	broker := broker.NewBroker()

	handler := handler.NewHandler(broker)

	apiVersion.POST("/create-queue", handler.QueueHandler)
	apiVersion.POST("/add-subscriber", handler.AddSubscriberHandler)
	apiVersion.DELETE("/remove-subscriber/:id", handler.RemoveSubscriberHandler)
	apiVersion.POST("/publish-by-id", handler.PublishbyIdHandler)
	apiVersion.POST("/publish-to-all", handler.PublishToAllHandler)
	apiVersion.POST("/subscribe", handler.SubscribeHandler)
	apiVersion.POST("/subscribe-by-id", handler.SubscribeByIdHandler)

	router.Run(":" + config.LocalConfig.Port)
}
