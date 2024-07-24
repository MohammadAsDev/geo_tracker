package entities

type Trip struct {
	TripHash string
	DriverId uint64
	RiderId  uint64
	Dest     GeoPos
}
