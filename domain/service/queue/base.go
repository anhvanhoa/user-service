package queue

import "auth-service/constants"

type Payload struct {
	Provider string
	Tos      *[]string
	Template string
	Data     map[string]any
}

type QueueClient interface {
	EnqueueMail(payload Payload) (string, error)
	EnqueueAnyTask(taskType constants.QueueType, payload Payload) (string, error)
}
