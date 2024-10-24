package events

import (
	"auth/libs/listeners"
	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

func NewBarEventContract(eventDriver eventsDriver) listeners.BackgroundWorker {
	// Define your topic handlers here
	eventDriver.BindToTopic("foo", func(_ *sarama.ConsumerMessage) {
		log.Info().Msg("foo")
	})

	return eventDriver
}
