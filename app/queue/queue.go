package queue

import (
	"errors"
	"sync"

	"github.com/google/logger"
	"github.com/saadrupai/go-message-broker/app/config"
	"github.com/saadrupai/go-message-broker/app/models"
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

	connection, err := config.LocalConfig.Listener.Accept()
	if err != nil {
		logger.Error("failed to accept listener")
	}

	newSubscriber := models.Subscriber{
		Id:             subscriberReq.SubscriberId,
		SubscriberName: subscriberReq.SubscriberName,
		Channel:        make(chan string, subscriberReq.BufferSize),
		Connection:     connection,
	}

	q.Subscribers[subscriberReq.SubscriberId] = newSubscriber

	logger.Infof("Subscriber %d added successfully", subscriberReq.SubscriberId)

	return nil
}

func (q *Queue) RemoveSubscriber(subscriberId uint) error {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	if subscriber, exists := q.Subscribers[subscriberId]; exists {
		close(subscriber.Channel)
		delete(q.Subscribers, subscriberId)
	} else {
		return errors.New("subscriber does not exist")
	}
	return nil
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
	subscriber := q.Subscribers[subscriberId]
	select {
	case subscriber.Channel <- message:
		logger.Infof("Message queued for subscriber %d: %s", subscriberId, message)
		return nil
	default:
		return errors.New("there is no space in queue :" + q.Name)
	}
}

func (q *Queue) SubscribeById(subscriberId uint) {
	subscriber := q.Subscribers[subscriberId]
	logger.Info("starting go routine")

	go func(subscriber models.Subscriber) {

		for {

			logger.Info("listening.........")

			select {
			case message := <-subscriber.Channel:

				logger.Info("writing msg.........")

				_, err := subscriber.Connection.Write([]byte(message + "\n"))
				if err != nil {
					logger.Error("failed to write message to subscriber")
				}
			}
		}
	}(subscriber)
}

func (q *Queue) Subscribe() (string, error) {

	select {
	case message := <-q.Channel:
		return message, nil
	default:
		return "", errors.New("there is no message available in queue")
	}
}

func (q *Queue) SubscribeTest(subscriberId uint) {

	select {
	case message := <-q.Subscribers[subscriberId].Channel:
		_, err := q.Subscribers[subscriberId].Connection.Write([]byte(message + "\n"))
		if err != nil {
			logger.Error("failed to write message to subscriber")
			return
		}
	}
}
