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
	apiVersion.DELETE("/remove-subscriber/:queue/:id", handler.RemoveSubscriberHandler)
	apiVersion.POST("/publish-by-id", handler.PublishHandler)
	apiVersion.POST("/publish-to-all", handler.PublishToAllHandler)
	apiVersion.GET("/subscribe", handler.SubscribeHandler)
	apiVersion.GET("/subscribe-by-id/:queue/:id", handler.SubscribeByIdHandler)

	router.Run(":" + config.LocalConfig.Port)
}
