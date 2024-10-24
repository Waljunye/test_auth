package events

import (
	"auth/libs"
	"auth/libs/listeners"
)

type eventsDriver interface {
	listeners.BackgroundWorker
	BindToTopic(topic string, handler libs.KafkaHandler)
}
