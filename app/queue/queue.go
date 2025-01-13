package queue

import (
	"errors"
	"sync"

	"github.com/google/logger"
	"github.com/saadrupai/go-message-broker/app/models"
	"github.com/saadrupai/go-message-broker/utils"
)

type Queue struct {
	Name        string
	Subscribers map[uint]models.Subscriber
	Channel     chan string
	Mutex       sync.Mutex
}

func NewQueue(name string, bufferSize int) *Queue {
	return &Queue{
		Name:        name,
		Channel:     make(chan string, bufferSize),
		Subscribers: make(map[uint]models.Subscriber, 0),
	}
}

func (q *Queue) AddSubscriber(subscriberReq models.AddSubscriber) error {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	if _, exists := q.Subscribers[subscriberReq.SubscriberId]; exists {
		return errors.New("subscriber already exists")
	}

	newSubscriber := models.Subscriber{
		Id:             subscriberReq.SubscriberId,
		SubscriberName: subscriberReq.SubscriberName,
		Channel:        make(chan string, subscriberReq.BufferSize),
	}

	newSubscriber.Webhook = utils.GenerateWebhook(subscriberReq.SubscriberId)

	go func(subscriber models.Subscriber) {
		for message := range subscriber.Channel {
			logger.Infof("Sending message to %s: %s", subscriber.SubscriberName, message)
			_, err := subscriber.Connection.Write([]byte(message + "\n"))
			if err != nil {
				logger.Error("failed to write message to subscriber")
				return
			}
		}

	}(newSubscriber)

	q.Subscribers[subscriberReq.SubscriberId] = newSubscriber

	logger.Infof("Subscriber %d added successfully", subscriberReq.SubscriberId)

	return nil
}

func (q *Queue) RemoveSubscriber(subscriberId uint) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	if subscriber, exists := q.Subscribers[subscriberId]; exists {
		close(subscriber.Channel)
		delete(q.Subscribers, subscriberId)
		logger.Info("Subscriber %d removed\n", subscriberId)
	}
}

func (q *Queue) PublishToAll(message string) error {
	select {
	case q.Channel <- message:
		return nil
	default:
		return errors.New("there is no space in queue :" + q.Name)
	}
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

	select {
	case message := <-q.Channel:
		return message, nil
	default:
		return "", errors.New("there is no message available in queue")
	}
}
