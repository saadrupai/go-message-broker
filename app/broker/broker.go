package broker

import (
	"errors"
	"sync"

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

func (b *Broker) CreateQueue(name string, bufferSize int) error {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	if _, exits := b.Queues[name]; exits {
		return errors.New("queue already exists with this name")
	}

	b.Queues[name] = queue.NewQueue(name, bufferSize)
	return nil
}

func (b *Broker) Publish(queueName, message string) error {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	queue, exists := b.Queues[queueName]
	if !exists {
		return errors.New("queue does not exist")
	}

	if err := queue.Publish(message); err != nil {
		return err
	}

	return nil
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
