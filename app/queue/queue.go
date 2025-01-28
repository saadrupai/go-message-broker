package queue

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/saadrupai/go-message-broker/app/config"
	"github.com/saadrupai/go-message-broker/app/consts"
	"github.com/saadrupai/go-message-broker/app/models"
)

type Queue struct {
	Name        string
	Subscribers map[uint]models.Subscriber
	Mutex       sync.Mutex
}

func NewQueue(name string, bufferSize int) *Queue {
	return &Queue{
		Name:        name,
		Subscribers: make(map[uint]models.Subscriber, 0),
	}
}

func (q *Queue) AddSubscriber(subscriberReq models.AddSubscriber) models.Subscriber {
	newSubscriber := models.Subscriber{
		Id:             subscriberReq.SubscriberId,
		SubscriberName: subscriberReq.SubscriberName,
		Channel:        make(chan string, subscriberReq.BufferSize),
	}
	q.Subscribers[subscriberReq.SubscriberId] = newSubscriber

	return newSubscriber
}

func (q *Queue) RemoveSubscriber(subscriberId uint) {
	delete(q.Subscribers, subscriberId)
}

func (q *Queue) PublishToAll(message string) error {

	err := config.LocalConfig.RedisCLient.Set(context.Background(), consts.PublishToAll, message, 0).Err()
	if err != nil {
		log.Fatal("failed to store data in redis")
		return err
	}
	return nil
}

func (q *Queue) PublishById(message string, subscriberId uint) error {
	select {
	case q.Subscribers[subscriberId].Channel <- message:
		return nil
	default:
		return errors.New("there is no space in queue :" + q.Name)
	}
}

func (q *Queue) SubscribeById(subscriberId uint) (string, error) {
	select {
	case message := <-q.Subscribers[subscriberId].Channel:
		return message, nil
	default:
		return "", errors.New("there is no message available in queue")
	}
}

func (q *Queue) Subscribe() (string, error) {
	message, err := config.LocalConfig.RedisCLient.Get(context.Background(), consts.PublishToAll).Result()
	if err != nil {
		log.Fatal("failed to get data from redis")
		return "", err
	}

	return message, nil
}
