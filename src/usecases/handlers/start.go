package handlers

import (
	"encoding/json"
	"errors"

	"github.com/MohammadAsDev/geo_tracker/src/entities"
	"github.com/MohammadAsDev/geo_tracker/src/entities/events"
)

type StartCommand struct {
	_Handler
}

func (handler StartCommand) Hanlde() error {
	trips_mem := handler.Tracker.Trips()
	payload := handler.Event.Payload

	var start_payload events.StartPayload
	if err := json.Unmarshal([]byte(payload), &start_payload); err != nil {
		return errors.New("error while unmarshling start event payload")
	}

	trips_mem[start_payload.DriverId] = entities.Trip{
		DriverId: start_payload.DriverId,
		RiderId:  start_payload.RiderId,
		TripHash: start_payload.TripId,
		Dest:     start_payload.Destination,
	}
	return nil
}
