package tracker

import (
	"errors"

	"github.com/MohammadAsDev/geo_tracker/src/interfaces"
)

type NopTracker struct { // Just No Operation Tracker
	// Empty tracker
}

func (tracker NopTracker) StartTracker() error {
	return errors.New("nop tracker")
}

func (tracker NopTracker) Cache() interfaces.Cache {
	return nil
}

func (Tracker NopTracker) ConnServer() interfaces.ConnServer {
	return nil
}

func (Tracker NopTracker) Consumer() interfaces.Consumer {
	return nil
}

func (Trakcker NopTracker) Trips() TripsMem {
	return nil
}
