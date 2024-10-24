package libs

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

func NewKafkaDriver(servers []string, consumerGroup string, cfg *sarama.Config) *kafkaDriver {
	return &kafkaDriver{
		handlers:      make(map[string]KafkaHandler),
		cfg:           cfg,
		servers:       servers,
		consumerGroup: consumerGroup,
	}
}

type kafkaDriver struct {
	cfg           *sarama.Config
	handlers      map[string]KafkaHandler
	servers       []string
	consumerGroup string
}

func (k *kafkaDriver) BindToTopic(topic string, handler KafkaHandler) {
	k.handlers[topic] = handler
}

type KafkaHandler func(msg *sarama.ConsumerMessage)

func (k *kafkaDriver) Start() (err error) {
	ctx := context.Background()
	consumerGroup, err := sarama.NewConsumerGroup(k.servers, k.consumerGroup, k.cfg)
	if err != nil {
		return err
	}
	defer consumerGroup.Close()

	for topic, handler := range k.handlers {
		go func(topic string, handler KafkaHandler) {
			consumer := NewConsumer(handler)
			for {
				err := consumerGroup.Consume(ctx, []string{topic}, consumer)
				if err != nil {
					log.Error().Err(err).Msg("Error from consumer")
				}
			}
		}(topic, handler)
	}

	return nil
}
func (k *kafkaDriver) Stop() (err error) {
	return nil
}

func (k *kafkaDriver) Info() string {
	topics := make([]string, 0, len(k.handlers))
	for topic := range k.handlers {
		topics = append(topics, topic)
	}
	return fmt.Sprintf("kafka consumer: Consumer group: %s, topics: %v", k.consumerGroup, topics)
}

func NewConsumer(handler KafkaHandler) sarama.ConsumerGroupHandler {
	return &messageConsumer{}
}

type messageConsumer struct {
	handler KafkaHandler
}

func (m *messageConsumer) Setup(sarama.ConsumerGroupSession) error {
	// not implemented
	return nil
}

func (m *messageConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	// not implemented
	return nil
}

func (m *messageConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		m.handler(message)
		session.MarkMessage(message, "")
	}
	return nil
}
