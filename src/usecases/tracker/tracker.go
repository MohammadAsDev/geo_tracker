package tracker

import (
	"github.com/MohammadAsDev/geo_tracker/src/entities"
	"github.com/MohammadAsDev/geo_tracker/src/interfaces"
)

type TripsMem map[uint64]entities.Trip

type Tracker interface {
	StartTracker() error
	Cache() interfaces.Cache
	ConnServer() interfaces.ConnServer
	Consumer() interfaces.Consumer
	Trips() TripsMem
}

type TrackerType uint8

const (
	NOP     TrackerType = 0
	DEFAULT TrackerType = 1
)
