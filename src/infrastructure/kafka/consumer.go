package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/MohammadAsDev/geo_tracker/src/config"
	"github.com/MohammadAsDev/geo_tracker/src/entities/events"
	"github.com/MohammadAsDev/geo_tracker/src/interfaces"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumer struct {
	_KafkaConsumer *kafka.Consumer
	_Err           error
	_Topic         string
	_Ctx           context.Context
	_TimeOutSecs   int

	_ErrsChan   chan interfaces.ConsumingErr
	_EventsChan chan events.Event
}

const TIME_OUT_SECS = 120

func NewKafkaConsumer(ctx context.Context, config *config.KafkaConfig) interfaces.Consumer {

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.Addr,
		"group.id":          config.GroupId,
	})

	return &KafkaConsumer{
		_KafkaConsumer: consumer,
		_Err:           err,
		_Topic:         config.Topic,
		_TimeOutSecs:   TIME_OUT_SECS,
		_Ctx:           ctx,

		_ErrsChan:   make(chan interfaces.ConsumingErr),
		_EventsChan: make(chan events.Event),
	}
}

func (consumer *KafkaConsumer) GetCommandsChannel() chan events.Event {
	return consumer._EventsChan
}

func (consumer *KafkaConsumer) GetErrorsChannel() chan interfaces.ConsumingErr {
	return consumer._ErrsChan
}

func (consumer *KafkaConsumer) Start() error {
	if err := consumer._KafkaConsumer.Subscribe(consumer._Topic, nil); err != nil {
		return err
	}

	go func() {
		for {
			msg, err := consumer._KafkaConsumer.ReadMessage(time.Second * time.Duration(consumer._TimeOutSecs))
			if err != nil {
				consumer._ErrsChan <- interfaces.ConsumingErr{
					Err:  err,
					Code: interfaces.CONSUMING_ERR,
				}
				continue
			}

			var event events.Event
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				consumer._ErrsChan <- interfaces.ConsumingErr{
					Err:  errors.New("can't read incoming command"),
					Code: interfaces.PARSING_ERR,
				}
				continue
			}
			consumer._EventsChan <- event
		}
	}()
	return nil
}
