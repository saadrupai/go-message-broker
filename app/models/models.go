package models

type QueueCreateReq struct {
	Name       string `json:"name"`
	BufferSize int    `json:"buffer_size"`
}

type PublishReq struct {
	QueueName string `json:"queue_name"`
	Message   string `json:"message"`
}

type SubscribeReq struct {
	QueueName string `json:"queue_name"`
}
