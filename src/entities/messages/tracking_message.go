package messages

import "github.com/MohammadAsDev/geo_tracker/src/entities"

type TrackingMessage struct {
	DriverId uint64
	entities.GeoPos
}
