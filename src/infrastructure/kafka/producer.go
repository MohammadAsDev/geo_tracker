package kafka

import (
	"context"

	"github.com/MohammadAsDev/geo_tracker/src/config"
	"github.com/MohammadAsDev/geo_tracker/src/entities/events"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	_KafkaProducer *kafka.Producer
	_Topic         string
	_Ctx           context.Context
	_Err           error
}

func NewKafkaProducer(ctx context.Context, config *config.KafkaConfig) *KafkaProducer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.Addr,
		"group.id":          config.GroupId,
		"auto.offset.reset": config.OffsetReset,
	})

	return &KafkaProducer{
		_KafkaProducer: producer,
		_Err:           err,
		_Ctx:           ctx,
		_Topic:         config.Topic,
	}
}

func (producer *KafkaProducer) Produce(command events.Event) error {
	return nil
}
