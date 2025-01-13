package models

type QueueCreateReq struct {
	Name       string `json:"name"`
	BufferSize int    `json:"buffer_size"`
}

type PublishReq struct {
	QueueName    string `json:"queue_name"`
	Message      string `json:"message"`
	SubscriberId uint   `json:"subscriber_id"`
}

type SubscribeReq struct {
	QueueName    string `json:"queue_name"`
	SubscriberId uint   `json:"subscriber_id"`
}

type AddSubscriber struct {
	QueueName      string `json:"queue_name"`
	SubscriberId   uint   `json:"subscriber_id"`
	SubscriberName string `json:"subscriber_name"`
	BufferSize     int    `json:"buffer_size"`
}

type Subscriber struct {
	Id             uint   `json:"id"`
	SubscriberName string `json:"subscriber_name"`
	Channel        chan string
	Webhook        string
}
