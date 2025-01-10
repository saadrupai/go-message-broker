package handler

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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	if err := c.Broker.CreateQueue(queueReq.Name, queueReq.BufferSize); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create queue", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "message published successfully"})
}

func (c *Handler) PublishToAllHandler(ctx *gin.Context) {
	var publishReq *models.PublishReq

	if err := ctx.ShouldBindJSON(&publishReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	if err := c.Broker.PublishToAll(publishReq.QueueName, publishReq.Message); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish messsage", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "message published successfully"})
}

func (c *Handler) PublishHandler(ctx *gin.Context) {
	var publishReq models.PublishReq

	if err := ctx.ShouldBindJSON(&publishReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	if err := c.Broker.Publish(publishReq); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish messsage", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "message published successfully"})
}

func (c *Handler) AddSubscriberHandler(ctx *gin.Context) {
	var addSubscriberReq models.AddSubscriber

	if err := ctx.ShouldBindJSON(&addSubscriberReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	if err := c.Broker.AddSubscriber(addSubscriberReq); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish messsage", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "message published successfully"})
}

func (c *Handler) SubscribeHandler(ctx *gin.Context) {
	var subscribeReq *models.SubscribeReq

	if err := ctx.BindJSON(&subscribeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	message, err := c.Broker.Subscribe(subscribeReq.QueueName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get message", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": message})
}

func (c *Handler) SubscribeByIdHandler(ctx *gin.Context) {
	var subscribeReq *models.SubscribeReq

	if err := ctx.ShouldBindJSON(&subscribeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	message, err := c.Broker.SubscribeById(subscribeReq.QueueName, subscribeReq.SubscriberId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get message", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": message})
}
