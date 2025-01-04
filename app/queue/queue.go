package queue

import (
	"errors"
	"sync"
)

type Queue struct {
	Name    string
	Channel chan string
	Mutex   sync.Mutex
}

func NewQueue(name string, bufferSize int) *Queue {
	return &Queue{
		Name:    name,
		Channel: make(chan string, bufferSize),
	}
}

func (q *Queue) Publish(message string) error {
	select {
	case q.Channel <- message:
		return nil
	default:
		return errors.New("there is no space in queue :" + q.Name)
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
