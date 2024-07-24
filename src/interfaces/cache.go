package interfaces

import "github.com/MohammadAsDev/geo_tracker/src/entities"

type Cache interface {
	UpdateVehicleCoords(id string, long float64, lat float64) error
	GetDriverToken(id string) (string, error)
	GetVehicleDestination(id string) (entities.GeoPos, error)
	GetVehicleCoords(id string) (entities.GeoPos, error)
	FreeRider(rider_id uint64) error
	FreeDriver(driver_id uint64) error
}
