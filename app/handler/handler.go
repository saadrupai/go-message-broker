package client

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saadrupai/go-message-broker/app/broker"
	"github.com/saadrupai/go-message-broker/app/models"
)

type Handler struct {
	Broker *broker.Broker
}

func NewHandler(broker *broker.Broker) *Handler {
	return &Handler{
		Broker: broker,
	}
}

func (c *Handler) QueueHandler(ctx *gin.Context) {
	var queueReq *models.QueueCreateReq

	if err := ctx.ShouldBindJSON(&queueReq); err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid request")
		return
	}

	if err := c.Broker.CreateQueue(queueReq.Name, queueReq.BufferSize); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusCreated, "queue created successfully")
}

func (c *Handler) PublishHandler(ctx *gin.Context) {
	var publishReq *models.PublishSubscribeReq

	if err := ctx.ShouldBindJSON(&publishReq); err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid request")
		return
	}

	if err := c.Broker.Publish(publishReq.QueueName, publishReq.Message); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusCreated, "message published successfully")
}

func (c *Handler) SubscribeHandler(ctx *gin.Context) {
	var subscribeReq *models.PublishSubscribeReq

	if err := ctx.ShouldBindJSON(&subscribeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid request")
		return
	}

	message, err := c.Broker.Subscribe(subscribeReq.QueueName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusCreated, message)
}
