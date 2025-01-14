package broker

import (
	"errors"
	"net"
	"sync"

	"github.com/saadrupai/go-message-broker/app/models"
	"github.com/saadrupai/go-message-broker/app/queue"
)

type Broker struct {
	Queues map[string]*queue.Queue
	Mutex  sync.Mutex
}

func NewBroker() *Broker {
	return &Broker{
		Queues: make(map[string]*queue.Queue),
	}
}

func (b *Broker) AddSubscriber(subscriberReq models.AddSubscriber, connection net.Conn) error {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	if queue, exits := b.Queues[subscriberReq.QueueName]; exits {
		queue.AddSubscriber(subscriberReq, connection)
	} else {
		return errors.New("queue does not exist")
	}

	return nil
}

func (b *Broker) RemoveSubscriber(subscriberId uint, queueName string) error {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	if queue, exits := b.Queues[queueName]; exits {
		err := queue.RemoveSubscriber(subscriberId)
		if err != nil {
			return err
		}
	} else {
		return errors.New("queue does not exist")
	}

	return nil
}

func (b *Broker) CreateQueue(name string, bufferSize int) error {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	if _, exits := b.Queues[name]; exits {
		return errors.New("queue already exists with this name")
	}

	b.Queues[name] = queue.NewQueue(name, bufferSize)
	return nil
}

func (b *Broker) PublishToAll(queueName, message string) error {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	queue, exists := b.Queues[queueName]
	if !exists {
		return errors.New("queue does not exist")
	}

	if err := queue.PublishToAll(message); err != nil {
		return err
	}

	return nil
}

func (b *Broker) PublishById(publishReq models.PublishReq) error {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	queue, exists := b.Queues[publishReq.QueueName]
	if !exists {
		return errors.New("queue does not exist")
	}

	if _, exists := queue.Subscribers[publishReq.SubscriberId]; exists {
		if err := queue.PublishById(publishReq.Message, publishReq.SubscriberId); err != nil {
			return err
		}
	}

	return nil
}

func (b *Broker) SubscribeById(queueName string, subscriberId uint) (string, error) {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	queue, exists := b.Queues[queueName]
	if !exists {
		return "", errors.New("queue does not exist")
	}

	if _, exists := queue.Subscribers[subscriberId]; exists {
		message, err := queue.SubscribeById(subscriberId)
		if err != nil {
			return "", err
		}
		return message, nil
	}

	return "no message available", nil
}

func (b *Broker) Subscribe(queueName string) (string, error) {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	queue, exists := b.Queues[queueName]
	if !exists {
		return "", errors.New("queue does not exist")
	}

	message, err := queue.Subscribe()
	if err != nil {
		return "", err
	}

	return message, nil
}
