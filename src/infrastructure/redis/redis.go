package redis

import (
	"context"
	"fmt"

	"github.com/MohammadAsDev/geo_tracker/src/config"
	"github.com/MohammadAsDev/geo_tracker/src/entities"
	"github.com/MohammadAsDev/geo_tracker/src/interfaces"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	_RedisClient *redis.Client
	_Ctx         context.Context

	_TrackingKey            string
	_ActiveDriversKeyFormat string
	_TripsKeyFormat         string
}

func NewRedisCache(ctx context.Context, config *config.RedisConfig) interfaces.Cache {

	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.Db,
	})

	return &RedisCache{
		_RedisClient: client,
		_Ctx:         ctx,

		_ActiveDriversKeyFormat: config.ActiveDriversKeyFormat,
		_TrackingKey:            config.TrackingKey,
		_TripsKeyFormat:         config.TripsKeyFormat,
	}
}

func (cache *RedisCache) Close() error {
	return cache._RedisClient.Close()
}

func (cache *RedisCache) UpdateVehicleCoords(id string, long float64, lat float64) error {
	_, err := cache._RedisClient.GeoAdd(cache._Ctx, cache._TrackingKey, &redis.GeoLocation{
		Longitude: long,
		Latitude:  lat,
		Name:      fmt.Sprintf("vehicle:%s", id),
	}).Result()
	return err
}

func (cache *RedisCache) GetDriverToken(id string) (string, error) {
	s := fmt.Sprintf(cache._ActiveDriversKeyFormat, id)
	return cache._RedisClient.HGet(cache._Ctx, s, "token").Result()

}

func (cache *RedisCache) GetVehicleDestination(trip_id string) (entities.GeoPos, error) {
	end_coords, err := cache._RedisClient.HMGet(cache._Ctx, fmt.Sprintf(cache._ActiveDriversKeyFormat, trip_id), "end_xcoord", "end_ycoord").Result()
	return entities.GeoPos{
		Longitude: end_coords[1].(float64),
		Latitude:  end_coords[0].(float64),
	}, err
}

func (cache *RedisCache) GetVehicleCoords(id string) (entities.GeoPos, error) {
	coords, err := cache._RedisClient.GeoPos(cache._Ctx, cache._TrackingKey, fmt.Sprintf("vehicle:%s", id)).Result()
	return entities.GeoPos{
		Longitude: coords[0].Longitude,
		Latitude:  coords[0].Latitude,
	}, err
}

func (cache *RedisCache) FreeRider(rider_id uint64) error {
	busy_riders_cache_format := "cache:busy_riders:%d"
	return cache._RedisClient.Del(cache._Ctx, fmt.Sprintf(busy_riders_cache_format, rider_id)).Err()
}

func (cache *RedisCache) FreeDriver(driver_id uint64) error {
	busy_drivers_cache_format := "cache:busy_drivers:%d"
	return cache._RedisClient.Del(cache._Ctx, fmt.Sprintf(busy_drivers_cache_format, driver_id)).Err()

}
