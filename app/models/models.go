package models

type QueueCreateReq struct {
	Name       string `json:"name"`
	BufferSize int    `json:"buffer_size"`
}

type PublishSubscribeReq struct {
	QueueName string `json:"queue_name"`
	Message   string `json:"message"`
}
