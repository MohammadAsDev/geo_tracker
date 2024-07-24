package tracker

import (
	"context"

	"github.com/MohammadAsDev/geo_tracker/src/config"
	"github.com/MohammadAsDev/geo_tracker/src/entities"
	"github.com/MohammadAsDev/geo_tracker/src/infrastructure/kafka"
	"github.com/MohammadAsDev/geo_tracker/src/infrastructure/redis"
	"github.com/MohammadAsDev/geo_tracker/src/infrastructure/ws"
	"github.com/MohammadAsDev/geo_tracker/src/interfaces"
)

type DefaultTracker struct {
	_Ctx context.Context

	_ConnServer interfaces.ConnServer
	_Cache      interfaces.Cache
	_Consumer   interfaces.Consumer

	_RunningTrips TripsMem
}

func NewDefaultTracker(ctx context.Context, config *config.Config) *DefaultTracker {
	redis_cache := redis.NewRedisCache(ctx, &config.RedisConfig)
	ws_server := ws.NewWsBuilder(ctx, &config.WsConfig).
		WithDefaultCheckOrigin().
		WithSystemCache(redis_cache).
		Build()

	return &DefaultTracker{
		_Ctx: ctx,

		_ConnServer: ws_server,
		_Cache:      redis_cache,
		_Consumer:   kafka.NewKafkaConsumer(ctx, &config.TrackingBrokerConfig.KafkaBroker),

		_RunningTrips: make(map[uint64]entities.Trip),
	}

}

func (tracker *DefaultTracker) StartTracker() error {
	if err := tracker._Consumer.Start(); err != nil {
		return err
	}

	if err := tracker._ConnServer.StartServer(); err != nil {
		return err
	}

	return nil
}

func (tracker *DefaultTracker) Cache() interfaces.Cache {
	return tracker._Cache
}

func (tracker *DefaultTracker) ConnServer() interfaces.ConnServer {
	return tracker._ConnServer
}

func (tracker *DefaultTracker) Consumer() interfaces.Consumer {
	return tracker._Consumer
}

func (tracker *DefaultTracker) Trips() TripsMem {
	return tracker._RunningTrips
}
