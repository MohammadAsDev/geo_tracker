package usecases

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/MohammadAsDev/geo_tracker/src/config"
	"github.com/MohammadAsDev/geo_tracker/src/entities"
	"github.com/MohammadAsDev/geo_tracker/src/infrastructure/logger/default_logger"
	"github.com/MohammadAsDev/geo_tracker/src/interfaces"
	"github.com/MohammadAsDev/geo_tracker/src/usecases/handlers"
	"github.com/MohammadAsDev/geo_tracker/src/usecases/tracker"
)

type TrackingSystem struct {
	_Tracker tracker.Tracker
	_Logger  interfaces.Logger
}

func NewTrackingSystem(ctx context.Context, trackerType tracker.TrackerType, config *config.Config) *TrackingSystem {
	var tracker_obj tracker.Tracker
	switch trackerType {
	case tracker.DEFAULT:
		tracker_obj = tracker.NewDefaultTracker(ctx, config)
	default:
		tracker_obj = tracker.NopTracker{}
	}
	return &TrackingSystem{
		_Tracker: tracker_obj,
		_Logger:  default_logger.NewDefaultLogger(),
	}
}
func (system *TrackingSystem) _CalculateDistance(start entities.GeoPos, end entities.GeoPos) float64 {
	var r float64 = 6371
	dlat := end.Latitude - start.Latitude
	dlon := end.Longitude - start.Longitude

	start_lat_rad := start.Latitude * math.Pi / 180

	end_lat_rand := end.Latitude * math.Pi / 180

	dlat_rad := dlat * math.Pi / 180
	dlon_rand := dlon * math.Pi / 180

	a :=
		math.Sin(dlat_rad/2)*math.Sin(dlat_rad/2) +
			math.Cos(start_lat_rad)*math.Cos(end_lat_rand)*
				math.Sin(dlon_rand/2)*math.Sin(dlon_rand/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return r * c
}

func (system *TrackingSystem) _CompletedTrip(driver_id uint64, loc entities.GeoPos) bool {
	trips_mem := system._Tracker.Trips()

	const FIXED_DIST_M = 200

	dist := system._CalculateDistance(trips_mem[driver_id].Dest, loc)

	return dist*1000 <= FIXED_DIST_M
}

func (system *TrackingSystem) _UpdateVehicleLocation(driver_id uint64, loc entities.GeoPos) error {

	trips_mem := system._Tracker.Trips()

	if system._Tracker.Cache() == nil {
		system._Logger.PrintErrorLogMessage("tracking system has no cache")
		panic("")
	}
	if !loc.IsValid() {
		return errors.New("invalid location")
	}

	if trip, found := trips_mem[driver_id]; found && system._CompletedTrip(driver_id, loc) {
		// send command to timer to stop
		system._Tracker.Cache().FreeDriver(trip.DriverId)
		system._Tracker.Cache().FreeRider(trip.RiderId)
		fmt.Printf("driver(%d) has finished his trip\n", driver_id)
		return nil
	}

	return system._Tracker.Cache().UpdateVehicleCoords(fmt.Sprintf("%d", driver_id), loc.Longitude, loc.Latitude)
}

func (system *TrackingSystem) StartTrackingSystem() {

	if system._Tracker.ConnServer() == nil {
		system._Logger.PrintErrorLogMessage("no connection server is attached with tracking server")
		panic("")
	}

	msgs_chan, connErrs_chan := system._Tracker.ConnServer().GetMessagesChannel(), system._Tracker.ConnServer().GetErrsChan()
	events_chan, consumerErrs_chan := system._Tracker.Consumer().GetCommandsChannel(), system._Tracker.Consumer().GetErrorsChannel()

	go func() {
		for {
			select {
			case msg := <-msgs_chan:
				system._UpdateVehicleLocation(msg.DriverId, msg.GeoPos)
				system._Logger.PrintInfoLogMessage(fmt.Sprintf("A new message is coming: %v", msg))
			case err := <-connErrs_chan:
				system._Logger.PrintErrorLogMessage(err.Error())
			}
		}
	}()

	go func() {
		for {
			select {
			case event := <-events_chan:
				system._Logger.PrintInfoLogMessage(fmt.Sprintf("Consuming: %v", event))
				if err := handlers.BuildHandler(event, system._Tracker).Hanlde(); err != nil {
					system._Logger.PrintErrorLogMessage(err.Error())
				}
			case err := <-consumerErrs_chan:
				system._Logger.PrintErrorLogMessage(err.Error())
			}
		}
	}()

	if err := system._Tracker.StartTracker(); err != nil {
		panic(err)
	}

}
