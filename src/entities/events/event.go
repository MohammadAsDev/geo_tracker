package events

import "github.com/MohammadAsDev/geo_tracker/src/entities"

const (
	DEACTIVATE_EVENT = 0
	START_EVNET      = 1
)

type EventId uint8

type Event struct {
	EventId EventId `json:"event_id"`
	Payload string  `json:"payload"`
}

type StartPayload struct {
	RiderId     uint64          `json:"rider_id"`
	DriverId    uint64          `json:"driver_id"`
	TripId      string          `json:"trip_id"`
	Destination entities.GeoPos `json:"destination"`
}

type DeactivatePayload struct {
	DriverId uint64 `json:"driver_id"`
}
