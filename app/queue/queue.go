package queue

import (
	"errors"
	"sync"

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

func (q *Queue) AddSubscriber(subscriberReq models.AddSubscriber) models.Subscriber {
	newSubscriber := models.Subscriber{
		Id:             subscriberReq.SubscriberId,
		SubscriberName: subscriberReq.SubscriberName,
		Channel:        make(chan string, subscriberReq.BufferSize),
	}
	q.Subscribers[subscriberReq.SubscriberId] = newSubscriber

	return newSubscriber
}

func (q *Queue) PublishToAll(message string) error {
	select {
	case q.Channel <- message:
		return nil
	default:
		return errors.New("there is no space in queue :" + q.Name)
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
	select {
	case message := <-q.Channel:
		return message, nil
	default:
		return "", errors.New("there is no message available in queue")
	}
}
